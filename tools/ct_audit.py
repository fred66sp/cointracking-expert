#!/usr/bin/env python3
"""
ct_audit.py — Chequeos deterministas sobre una exportación de CoinTracking (dos variantes soportadas).

Objetivo (ADR-006, ADR-009): que el agente ejecute lógica de auditoría **vetada y
reproducible** en vez de re-derivarla cada vez (donde puede equivocarse, p. ej.
doble-contando comisiones). Devuelve resultados compactos (ADR-010): no vuelca el CSV.

Regla de saldo VETADA (ver knowledge/cointracking/CSV_FORMAT.md y COST_BASIS_AND_VALIDATION.md):
    saldo(activo) = Σ Compra[activo] − Σ Venta[activo]
La columna **Comisión NO se resta**: es informativa; el fee real ya está incluido en
Venta (retiradas/gastos) o llega como una fila aparte "Otras comisiones"/"Other Fee".

Formatos soportados (detectados automáticamente por la cabecera; ver CSV_FORMAT.md):
  - "es_trade_table": export "CSV" simple, locale español, 16 columnas, tipos en español,
    fecha DD.MM.YYYY.
  - "en_full_export": export "CSV (Exportación Completa)", locale inglés, 13 columnas
    (con columna LPN, sin columnas de dirección), tipos en inglés, fecha YYYY-MM-DD.

Uso:
    python tools/ct_audit.py <csv> [--exchange NOMBRE] [--check balances|transfers|duplicates|collisions|all]
    python tools/ct_audit.py <csv> --expect-balances '{"Coinbase":{"EUR":"-500"}}'   # validación

Salida: JSON compacto por stdout.
"""
import csv, sys, json, argparse
from decimal import Decimal, getcontext
from datetime import datetime, timedelta

# En Windows, stdout usa por defecto cp1252 y corrompería tildes/ñ del Comentario al
# volcar JSON (ensure_ascii=False). Forzar UTF-8 con independencia de la consola/entorno.
if hasattr(sys.stdout, "reconfigure"):
    sys.stdout.reconfigure(encoding="utf-8")

getcontext().prec = 40
TOL = Decimal("0.00000001")  # tolerancia de cuadre

# Índices de columna comunes a ambos formatos (parsear por POSICIÓN: hay columnas 'Cur.' repetidas)
C_TIPO, C_BUY, C_BUYCUR, C_SELL, C_SELLCUR, C_FEE, C_FEECUR = 0, 1, 2, 3, 4, 5, 6
C_EXCH, C_GROUP, C_COMMENT, C_DATE = 7, 8, 9, 10

FIAT = {"EUR", "USD", "GBP"}

# Formatos conocidos: cabecera esperada (columna 0), índice de Tx Hash/Tx-ID, tipos de
# depósito/retirada y formato de fecha. Ver knowledge/cointracking/CSV_FORMAT.md.
FORMATS = {
    "es_trade_table": {
        "header0": "Tipo",
        "min_cols": 16,
        "c_txhash": 13,
        "deposit_types": {"Depósito"},
        "withdrawal_types": {"Retirada"},
        "date_fmt": "%d.%m.%Y %H:%M:%S",
    },
    "en_full_export": {
        "header0": "Type",
        "min_cols": 13,
        "c_txhash": 12,
        "deposit_types": {"Deposit"},
        "withdrawal_types": {"Withdrawal"},
        "date_fmt": "%Y-%m-%d %H:%M:%S",
    },
}

# Estado del formato activo (fijado por detect_format() antes de auditar)
C_TXHASH = None
DEPOSIT_TYPES = None
WITHDRAWAL_TYPES = None
DATE_FMT = None


def detect_format(header):
    """Identifica la variante de export por la cabecera. Falla explícitamente si no la reconoce
    (ADR-009: no improvisar índices de columna sobre un formato no verificado)."""
    for name, spec in FORMATS.items():
        if len(header) >= spec["min_cols"] and header[0].strip() == spec["header0"]:
            return name, spec
    raise ValueError(
        f"Formato de CSV no reconocido (cabecera: {header}). "
        "No es ninguna de las variantes verificadas en knowledge/cointracking/CSV_FORMAT.md "
        "(es_trade_table, en_full_export). No se audita sin confirmar el formato antes."
    )


def configure_format(spec):
    global C_TXHASH, DEPOSIT_TYPES, WITHDRAWAL_TYPES, DATE_FMT
    C_TXHASH = spec["c_txhash"]
    DEPOSIT_TYPES = spec["deposit_types"]
    WITHDRAWAL_TYPES = spec["withdrawal_types"]
    DATE_FMT = spec["date_fmt"]


def D(x):
    x = (x or "").strip()
    return Decimal(x) if x else Decimal(0)


def parse_date(s):
    s = (s or "").strip()
    return datetime.strptime(s, DATE_FMT) if s else None


def load_rows(path):
    with open(path, encoding="utf-8-sig", newline="") as f:
        rows = list(csv.reader(f))
    header, body = rows[0], rows[1:]
    fmt_name, spec = detect_format(header)
    configure_format(spec)
    return fmt_name, body


def _rows_for(rows, exchange=None):
    return [r for r in rows if (exchange is None or r[C_EXCH] == exchange)]


def reconstruct_balances(rows, exchange=None):
    """{cuenta: {activo: Decimal}} usando Compra(+)/Venta(−), ignorando Comisión."""
    bal = {}
    for r in _rows_for(rows, exchange):
        acc = r[C_EXCH]
        bal.setdefault(acc, {})
        if r[C_BUY] and r[C_BUYCUR]:
            bal[acc][r[C_BUYCUR]] = bal[acc].get(r[C_BUYCUR], Decimal(0)) + D(r[C_BUY])
        if r[C_SELL] and r[C_SELLCUR]:
            bal[acc][r[C_SELLCUR]] = bal[acc].get(r[C_SELLCUR], Decimal(0)) - D(r[C_SELL])
    # limpiar ceros
    for acc in bal:
        bal[acc] = {k: v for k, v in bal[acc].items() if abs(v) > TOL}
    return bal


def negative_balances(rows, exchange=None):
    """Saldos negativos. Distingue FIAT (suele ser artefacto) de cripto (imposibilidad)."""
    bal = reconstruct_balances(rows, exchange)
    out = []
    for acc, assets in bal.items():
        for a, v in assets.items():
            if v < -TOL:
                out.append({"cuenta": acc, "activo": a, "saldo": str(v),
                            "tipo": "fiat" if a in FIAT else "cripto"})
    return out


def find_orphan_transfers(rows, window_min=1440):
    """Depósitos/retiradas sin contraparte. Nivel 1: Tx Hash. Nivel 2: moneda+importe(≈venta−fee)+ventana+cuenta distinta."""
    deps = [r for r in rows if r[C_TIPO] in DEPOSIT_TYPES]
    wds = [r for r in rows if r[C_TIPO] in WITHDRAWAL_TYPES]
    win = timedelta(minutes=window_min)

    def wd_expected(w):  # lo que debería llegar: venta − fee (si fee misma moneda)
        fee = D(w[C_FEE]) if w[C_FEECUR] == w[C_SELLCUR] else Decimal(0)
        return D(w[C_SELL]) - fee

    def match(dep):
        h = dep[C_TXHASH]
        for w in wds:
            if h and w[C_TXHASH] and h == w[C_TXHASH]:
                return ("hash", w)
        dcur, damt, ddate = dep[C_BUYCUR], D(dep[C_BUY]), parse_date(dep[C_DATE])
        best = None
        for w in wds:
            if w[C_SELLCUR] != dcur or w[C_EXCH] == dep[C_EXCH]:
                continue
            if abs(wd_expected(w) - damt) > TOL:
                continue
            wdate = parse_date(w[C_DATE])
            if ddate and wdate and abs((ddate - wdate)) <= win:
                dt = abs((ddate - wdate))
                if best is None or dt < best[1]:
                    best = (w, dt)
        return ("heur", best[0]) if best else (None, None)

    orphan_dep = []
    matched_wd = set()
    for d in deps:
        kind, w = match(d)
        if kind is None:
            orphan_dep.append({"cuenta": d[C_EXCH], "activo": d[C_BUYCUR], "importe": d[C_BUY],
                               "fecha": d[C_DATE], "fiat": d[C_BUYCUR] in FIAT,
                               "comentario": d[C_COMMENT]})
        else:
            matched_wd.add(id(w))
    orphan_wd = [{"cuenta": w[C_EXCH], "activo": w[C_SELLCUR], "importe": w[C_SELL],
                  "fecha": w[C_DATE], "comentario": w[C_COMMENT]}
                 for w in wds if id(w) not in matched_wd]
    return {"depositos_huerfanos": orphan_dep, "retiradas_huerfanas": orphan_wd,
            "nota": "Los depósitos fiat huérfanos suelen ser ingresos externos legítimos (no error)."}


def find_exact_duplicates(rows, exchange=None):
    from collections import Counter
    c = Counter(tuple(r) for r in _rows_for(rows, exchange))
    return [{"veces": v, "tipo": k[C_TIPO], "compra": k[C_BUY], "cur_compra": k[C_BUYCUR],
             "venta": k[C_SELL], "cur_venta": k[C_SELLCUR], "fecha": k[C_DATE],
             "comentario": k[C_COMMENT]}
            for k, v in c.items() if v > 1]


def ticker_collisions(rows, exchange=None):
    cur = set()
    for r in _rows_for(rows, exchange):
        for c in (r[C_BUYCUR], r[C_SELLCUR]):
            if c:
                cur.add(c)
    return sorted(c for c in cur if c and c[-1].isdigit())


def main():
    ap = argparse.ArgumentParser(description="Chequeos deterministas de CoinTracking (Trade Table).")
    ap.add_argument("csv")
    ap.add_argument("--exchange", default=None)
    ap.add_argument("--check", default="all",
                    choices=["balances", "negatives", "transfers", "duplicates", "collisions", "all"])
    ap.add_argument("--window-min", type=int, default=1440)
    ap.add_argument("--expect-balances", default=None,
                    help='JSON {"Cuenta":{"ACTIVO":"saldo"}} para validar la reconstrucción')
    args = ap.parse_args()

    fmt_name, rows = load_rows(args.csv)
    out = {"filas": len(rows), "exchange": args.exchange or "(todas)", "formato_detectado": fmt_name}

    if args.check in ("balances", "all"):
        bal = reconstruct_balances(rows, args.exchange)
        out["balances"] = {a: {k: str(v) for k, v in sorted(d.items())} for a, d in bal.items()}
    if args.check in ("negatives", "all"):
        out["saldos_negativos"] = negative_balances(rows, args.exchange)
    if args.check in ("transfers", "all"):
        out["transferencias"] = find_orphan_transfers(rows, args.window_min)
    if args.check in ("duplicates", "all"):
        out["duplicados_exactos"] = find_exact_duplicates(rows, args.exchange)
    if args.check in ("collisions", "all"):
        out["colision_tickers"] = ticker_collisions(rows, args.exchange)

    if args.expect_balances:
        expect = json.loads(args.expect_balances)
        bal = reconstruct_balances(rows)
        mismatches = []
        for acc, assets in expect.items():
            for a, v in assets.items():
                got = bal.get(acc, {}).get(a, Decimal(0))
                if abs(got - Decimal(str(v))) > TOL:
                    mismatches.append({"cuenta": acc, "activo": a, "esperado": str(v), "obtenido": str(got)})
        out["validacion"] = {"ok": not mismatches, "descuadres": mismatches}

    print(json.dumps(out, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
