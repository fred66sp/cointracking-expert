# Formatos CT-List — visualización compacta de listas en la conversación

**Ámbito (importante):** estos formatos son para **la conversación interactiva con el usuario** (chat), cuando hay que mostrarle listas de operaciones, hallazgos, balances o el recorrido de fondos. **No sustituyen** los informes formales de `reports/output/<proyecto>/` (plantilla `templates/AUDIT_REPORT.md`), que siguen necesitando tablas Markdown y el formato evidencia → causa → impacto → recomendación por trazabilidad ante el asesor (ADR-009/ADR-011). Son dos audiencias distintas: el usuario en el chat (necesita lectura rápida) y el asesor en el informe (necesita trazabilidad completa).

Adaptado 2026-07-04 de una propuesta del usuario ("CT-List"), corrigiendo dos imprecisiones antes de fijarlo: el ejemplo original usaba "Recepción" para un depósito (no es un tipo real de CoinTracking — el verificado es `Depósito`) y citaba el aviso de coste faltante en inglés ("Missing Purchase History") en vez de su forma en español ya documentada (`COST_BASIS_AND_VALIDATION.md` §3.1).

> 🔑 **Regla de uso (igual que el CT-Task, ADR-024):** estos bloques **nunca sustituyen** la explicación en lenguaje llano de `CLAUDE.md` — van justo debajo de ella, o intercalados hallazgo por hallazgo. Ante cualquier ⚠/✗, sigue traduciendo el hallazgo a **qué significa / por qué le importa / qué hacer ahora** (regla ya existente en `CLAUDE.md` §"Usuario objetivo y estilo de guía"). Usa siempre los **tipos de operación ya verificados** contra los datos reales del proyecto (`CSV_FORMAT.md` §3/§12): `Operación`, `Depósito`, `Retirada`, `Ingresos`, `Ingresos por intereses`, `Gasto`, `Staking`, `Recompensa / Bonificación`, `Otras comisiones`. Trata cualquier otro tipo como `[VERIFICAR]`.

---

## CT-Timeline — revisar un historial o seguir el orden cronológico

```
① Fecha ─ Tipo ─ Datos principales ─ Exchange
```

Ejemplo:

```
① 04.07.2026 09:00:00 ─ Depósito ─ 0.50 BTC ─ Ledger
② 04.07.2026 09:15:20 ─ Retirada ─ 0.50 BTC ─ Ledger
③ 04.07.2026 09:15:40 ─ Depósito ─ 0.50 BTC ─ Binance
④ 04.07.2026 10:35:58 ─ Operación ─ Compra 8 ETH | Venta 0.25 BTC | Comisión 0.0005 BTC ─ Binance
```

> Nota: una transferencia entre cuentas propias sigue siendo **Retirada + Depósito** (②③ arriba), nunca un tipo "Transferencia" único (mismo criterio que el CT-Task, ADR-024).

## CT-Audit — detectar y señalar errores

Marca cada línea con `✓` (correcto), `⚠` (advertencia) o `✗` (error). Nunca ocultes un problema; explica el motivo justo debajo, en lenguaje llano.

```
✓ 04.07.2026 Depósito 0.50 BTC
✓ 04.07.2026 Retirada 0.50 BTC (Ledger → Binance)
⚠ 07.07.2026 Venta 0.25 BTC
   "No hay una compra adecuada para esta venta" — no consta ninguna compra previa de BTC con coste registrado.
   Qué significa: CoinTracking no sabe cuánto pagaste por ese BTC, así que calcula la ganancia como si el coste fuera 0 → pagarías más impuestos de la cuenta.
   Qué hacer: si esos BTC vinieron de otra cuenta tuya, falta registrar el depósito de origen; si los compraste, falta dar de alta esa compra.
✓ 12.07.2026 Compra 8 ETH
```

## CT-Balance — balances y descuadres, agrupado por moneda

```
BTC
  +0.50  Depósito
  −0.25  Venta
  ────────────
  Saldo: 0.25 BTC

ETH
  +8  Compra
  −8  Retirada
  ────────────
  Saldo: 0 ETH
```

## CT-Exchange — varias plataformas

```
Ledger
  • Depósito 0.50 BTC
  • Retirada 0.50 BTC (→ Binance)

Binance
  • Depósito 0.50 BTC (← Ledger)
  • Compra 8 ETH
  • Retirada 8 ETH (→ MetaMask)

MetaMask
  • Depósito 8 ETH (← Binance)
```

## CT-Asset — todas las operaciones de una moneda concreta

```
BTC
  04.07  Compra   +0.50
  07.07  Venta    −0.25
  12.07  Compra   +0.10
  ────────────
  Saldo: 0.35 BTC
```

## CT-Flow — recorrido de fondos entre cuentas (útil para detectar transferencias mal registradas)

```
Ledger
  ↓ Retirada 0.50 BTC
Binance
  ↓ (Compra 8 ETH)
  ↓ Retirada 8 ETH
MetaMask
```

## CT-Task — dar de alta o corregir una operación manual

Ya documentado en `WEB_APP_GUIDE.md` §4bis (ADR-024). No se repite aquí; usar ese formato tal cual.

---

## Selección automática del formato

| El usuario quiere... | Formato |
|---|---|
| Introducir/corregir una operación manual | CT-Task (`WEB_APP_GUIDE.md` §4bis) |
| Entender qué ocurrió / seguir el orden temporal | CT-Timeline |
| Revisar errores de una auditoría | CT-Audit |
| Comprobar balances o descuadres | CT-Balance |
| Revisar un exchange/wallet concreto | CT-Exchange |
| Revisar una moneda concreta | CT-Asset |
| Entender el recorrido de una transferencia | CT-Flow |

## Reglas generales

- Una operación = una línea; mantén el orden cronológico salvo que agrupes por activo o exchange.
- Fechas siempre en `DD.MM.AAAA HH:MM:SS`.
- Usa los nombres de tipo de operación **verificados** (ver arriba); no inventes ni traduzcas libremente.
- En la **conversación**, prefiere estos formatos a tablas Markdown o párrafos largos cuando haya varias operaciones que revisar. En los **informes formales** de `reports/output/`, sigue usando tablas y el formato evidencia/causa/impacto/recomendación (no aplica esta regla ahí).
- Todo hallazgo (`⚠`/`✗`) lleva, debajo, la traducción a qué significa / por qué importa / qué hacer — nunca lo dejes solo con el símbolo.
