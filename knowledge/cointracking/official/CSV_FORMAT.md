---
id: KB-A2-001
title: "Formato CSV de CoinTracking (Trade Table)"
level: A
domain: cointracking
source: "CoinTracking — exportaciones de producción, datos reales del usuario"
authority: official
last_verified: 2026-07-03
valid_from: 2026-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-032
  - ADR-004
  - ADR-003

related_docs:
  - COST_BASIS_AND_VALIDATION.md
  - MCP_API.md
  - tools/ct_audit.py

tags:
  - cointracking
  - csv
  - format
  - import
  - trade-table

notes: "Validado contra datos reales (3.649 operaciones, 2 variantes). CoinTracking puede cambiar formato — reversificar si hay cambios de columnas/tipos/tickers."
---

# Formato CSV de CoinTracking (Trade Table)

Este documento describe el formato **real** de las exportaciones de operaciones de CoinTracking, verificado contra exportaciones de producción. Conforme a ADR-004, es la referencia autoritativa para la capa de importación: ninguna suposición sobre el formato debe cerrarse sin contrastarla aquí. `tools/ct_audit.py` detecta automáticamente cuál de las dos variantes recibe (§12 "Detección automática").

> ⚠️ **Aviso de alcance:** §1-10 documentan la variante con configuración regional en español (botón **CSV**); §12 documenta la variante en inglés con más columnas de metadatos (botón **CSV (Exportación Completa)**). CoinTracking permite otros idiomas y conjuntos de columnas aún no muestreados. Las peculiaridades dependientes de configuración están marcadas como **[config]**.

---

## 1. Estructura general

- **Codificación:** UTF-8 con BOM (leer con `utf-8-sig`)
- **Delimitador:** coma (`,`)
- **Todos los campos entrecomillados** con comillas dobles (`"..."`), incluidos los numéricos
- **16 columnas**
- Primera fila = cabecera

### Columnas (en orden)

| # | Nombre en CSV | Significado | Notas |
|---|---------------|-------------|-------|
| 0 | `Tipo` | Tipo de transacción | **[config]** en español; ver §3 |
| 1 | `Compra` | Cantidad recibida/entrante | Vacío si no aplica |
| 2 | `Cur.` | Moneda de "Compra" | |
| 3 | `Venta` | Cantidad enviada/saliente | Vacío si no aplica |
| 4 | `Cur.` | Moneda de "Venta" | **nombre duplicado** |
| 5 | `Comisión` | Importe de la comisión | Puede estar vacío |
| 6 | `Cur.` | Moneda de la comisión | **nombre duplicado**; puede diferir (§5) |
| 7 | `Intercambio` | Cuenta / exchange | Trato como identidad de cuenta (§4) |
| 8 | `Grupo` | Agrupación de CoinTracking | Casi siempre vacío |
| 9 | `Comentario` | Texto libre | Datos sucios; a veces contiene el tx hash (§6) |
| 10 | `Fecha` | Marca temporal | `DD.MM.YYYY HH:MM:SS`, **sin zona horaria** (§2) |
| 11 | `From Address` | Dirección origen | Frecuentemente vacío |
| 12 | `To Address` | Dirección destino | Frecuentemente vacío |
| 13 | `Tx Hash` | Hash de transacción on-chain | Presente solo en ~16-24% de transferencias (§7) |
| 14 | `Sell From Address` | Dirección origen (lado venta) | Frecuentemente vacío |
| 15 | `Sell To Address` | Dirección destino (lado venta) | Frecuentemente vacío |

> 🔧 **Implicación de implementación:** la cabecera tiene **tres columnas con el mismo nombre `Cur.`**. No se puede parsear a un dict por nombre de columna; hay que parsear **por posición**.

---

## 2. Formato de fecha — riesgo de determinismo

- Formato: `DD.MM.YYYY HH:MM:SS` (europeo). Ejemplo: `01.07.2026 08:48:26`
- Parseo: `datetime.strptime(valor, "%d.%m.%Y %H:%M:%S")`
- **No incluye zona horaria.** La exportación usa la zona configurada en la cuenta de CoinTracking, que el CSV no revela.

> ⚠️ **Riesgo de determinismo:** sin zona horaria explícita, dos exportaciones con distinta configuración regional producen instantes distintos para la misma operación → rompe la reproducibilidad.
>
> ✅ **Resuelto en ADR-005:** la zona de origen es un **parámetro obligatorio de importación**; la marca temporal se interpreta como hora local en la zona IANA declarada y se normaliza a **UTC**. Para la cuenta de referencia, la zona es `Europe/Madrid` (CET/CEST, **con horario de verano**). Se usa `zoneinfo` (DST-aware); **nunca** un offset fijo `+01:00`, que desplazaría 1 h las operaciones de verano.
>
> ⚠️ **Verificación pendiente:** confirmar que CoinTracking aplica DST al exportar (lo esperado) y no un offset fijo. Método definitivo: comparar una transferencia con `Tx Hash` contra su marca temporal on-chain (UTC). Ej.: retirada FARM del `01.07.2026 08:47:54` local (verano) debería ser `06:47:54 UTC` si hay DST, o `07:47:54 UTC` si offset fijo.

---

## 3. Tipos de transacción (columna `Tipo`)

Valores observados en datos reales **[config: español]** y su frecuencia:

| Valor CSV | Frecuencia | Categoría canónica propuesta |
|-----------|-----------:|------------------------------|
| `Operación` | 666 | Trade (compra↔venta) |
| `Otras comisiones` | 404 | Fee |
| `Recompensa / Bonificación` | 324 | Ingreso (reward) |
| `Staking` | 150 | Ingreso (staking) |
| `Pérdidas por Derivados / Futuros` | 114 | Derivados (pérdida) |
| `Beneficio de Derivados / Futuros` | 82 | Derivados (beneficio) |
| `Depósito` | 44 | Entrada (transferencia o fiat) |
| `Retirada` | 29 | Salida (transferencia) |
| `Ingresos por intereses` | 9 | Ingreso (interés) |
| `Ingresos` | 5 | Ingreso (genérico) |
| `Gasto` | 1 | Gasto |

> 🔧 La capa de normalización debe mapear estos literales a un enum canónico **interno en inglés** (PEP 8, ADR-001), tolerando valores desconocidos sin fallar (registrar y marcar para revisión).

### Qué columnas se rellenan según el tipo

- **Operación (trade):** `Compra`+`Cur.` (recibido) y `Venta`+`Cur.` (entregado). Comisión opcional.
- **Depósito:** `Compra`+`Cur.` (entrante). `Venta` vacío.
- **Retirada:** `Venta`+`Cur.` (saliente) + `Comisión`. `Compra` vacío.
- **Otras comisiones / Staking / Recompensa / Ingresos:** típicamente solo un lado relleno.

---

## 4. Cuentas / exchanges (columna `Intercambio`)

Valores reales: `Binance` (1191), `Binance Earn` (330), `BingX` (243), `Coinbase` (35), `Ledger Live` (21), `Metamask` (8).

> 🔑 **Hallazgo clave:** `Binance Earn` aparece como **cuenta distinta** de `Binance`. Los traspasos entre productos Earn y spot se registran como pares Depósito/Retirada entre "Binance" y "Binance Earn". La detección de transferencias debe tratar cada valor de `Intercambio` como una identidad de cuenta independiente.

---

## 5. Comisiones

- En **Operación**: 246 de 666 llevan comisión; 420 no. La moneda de comisión más común es la del lado vendido, pero no siempre.
- **La moneda de la comisión puede diferir de ambos activos operados.** Casos reales:
  - Retirada `300.62905200 USDC` con comisión `0.42 USDT`
  - Operación `BTC/USDC` con comisión `19.18 EUR`
  - Operación `AGIX/USDT` con comisión `0.00585300 ETH`

> 🔧 Debe tratarse la comisión como un `(importe, moneda)` **independiente** (así lo hace `tools/ct_audit.py`), nunca asumir que coincide con el activo comprado o vendido.

---

## 6. Columna `Comentario` — datos sucios reales

Texto libre, altamente heterogéneo. Patrones observados en depósitos:
- `from an external account <X>` (7)
- `BingX fund deposit` (8), `BingX manual deposit from note` (1)
- `Banco Transfer (SEPA)` (6) — **nótese la inconsistencia**: aparece también como `Bank Transfer (SEPA)`
- `Depósito PayPal Europe S.A.` (5)
- `Binance Savings Redemption`, `Cierre Binance Earn - traspaso a Spot`, `Binance Savings Purchase`
- A veces el **comentario contiene el propio tx hash** (ej. en retiradas on-chain)
- `[PENDIENTE_EVIDENCIA_API]` — anotación manual del usuario, no generada por CoinTracking

> 🔧 El comentario **no es fiable como dato estructurado**. Úsese solo como evidencia/pista secundaria, nunca como clave de reconciliación. No hacer coincidencias exactas de texto (`Banco` vs `Bank`).

---

## 7. Emparejamiento de transferencias — realidad medida

Este es el hallazgo que más contradice la intuición y que valida la estrategia de ADR-004.

**Datos:** 44 depósitos, 29 retiradas.

### Nivel 1 — por `Tx Hash` (fuerte pero escaso)
- Solo **7/44 depósitos** y **7/29 retiradas** llevan `Tx Hash`.
- De ellos, **6 parejas** casan por hash idéntico en ambos lados.
- En las 6 parejas se cumple **exactamente**: `depósito = retirada − comisión`.
  - Ej: FARM retirada `10.45500000`, comisión `0.06600000` → depósito `10.38900000`.
- Quedan **1 depósito y 1 retirada huérfanos** incluso teniendo hash.

### Nivel 2 — sin hash (la mayoría: 37 depósitos, 22 retiradas)
Son sobre todo movimientos internos (p. ej. Binance ↔ Binance Earn: redenciones de Savings, cierres de Earn).
- El match exacto por `(moneda, importe idéntico, misma fecha-hora)` **solo resolvió 3 de 37**.
- El resto requiere heurística tolerante: misma moneda, importe compatible con `retirada − comisión ≈ depósito`, ventana temporal (los timestamps difieren, p. ej. 32 s entre retirada y depósito), y cuentas distintas.

> 🔧 **Diseño obligado de la detección de transferencias (en niveles), tal como lo aplica `tools/ct_audit.py` y la skill `audit-cointracking`:**
> 1. **Determinista fuerte:** casar por `Tx Hash` cuando exista (alta confianza).
> 2. **Heurístico:** para el resto, casar por `moneda` + `importe ≈ retirada − comisión` + ventana temporal + cuentas distintas, con puntuación de confianza.
> 3. **Reportar huérfanos** explícitamente (principio "el silencio no es aceptable").
>
> Una regla basada solo en Tx Hash habría cubierto <20% de los casos.

---

## 8. Colisión de tickers — sufijo numérico de CoinTracking

CoinTracking **desambigua símbolos repetidos añadiendo un dígito**. Tickers reales observados:
`ICP2`, `ID2`, `PEPE4`, `PRIME3`, `SEI2`, `SOL2`, `THETA2`, `WLD3`.

> 🔑 **Crítico:** `SOL2` **no es** `SOL`; `WLD3` **no es** `WLD`. Son activos distintos que comparten símbolo base. La detección de duplicados y la normalización de activos **no deben** fusionar `SOL` con `SOL2`, ni asumir que el sufijo es un error tipográfico. La identidad de activo de CoinTracking es el ticker completo, sufijo incluido.

---

## 9. Duplicados exactos — no todos son errores

- **88 filas** son 100% idénticas a otra (32 grupos de duplicados), p. ej.:
  - `Otras comisiones · 0.00009792 BNB · Binance · 19.12.2024 18:49:12` ×3
  - `Depósito · 3.07724282 RDNT · Binance · 19.06.2025 00:16:00` ×2

> 🔧 Filas idénticas **no implican** error de importación: comisiones pequeñas recurrentes, recompensas periódicas y micro-movimientos pueden repetirse legítimamente en el mismo segundo. La detección de duplicados (`tools/ct_audit.py`, ADR-014) debe **señalar para revisión**, no eliminar automáticamente, y considerar contexto (tipo, `trade_id`, si hay una operación que las justifique).

---

## 10. Reglas de normalización (resumen para la capa de importación)

1. Leer con `utf-8-sig`, parsear **por posición** (no por nombre de columna).
2. Números: `Decimal(valor)` (ADR-002); hasta 8 decimales; sin separador de miles. Vacío → ausente, no cero.
3. Fecha: `strptime("%d.%m.%Y %H:%M:%S")` → interpretar en la zona IANA declarada (parámetro obligatorio; `Europe/Madrid` para la cuenta de referencia) y normalizar a UTC con `zoneinfo`, DST incluido (ADR-005, §2).
4. `Tipo`: mapear literales español → enum canónico interno (inglés); tolerar desconocidos.
5. Comisión: `(Decimal, moneda)` independiente; puede faltar o estar en tercera moneda.
6. Activo: usar el ticker **completo** (con sufijo); nunca fusionar por símbolo base.
7. Preservar la **fila original** junto a la normalizada (trazabilidad / procedencia).
8. `Comentario`/`Grupo`: metadato no fiable; no usar como clave.

---

## 11. Decisiones abiertas que este análisis genera

Estas cuestiones deberían resolverse (candidatas a ADR o a nueva regla en `tools/ct_audit.py`) antes de darlas por cerradas:

1. ~~**Zona horaria de las fechas**~~ → ✅ **Resuelto en ADR-005**: zona declarada por el usuario (`Europe/Madrid` para la cuenta de referencia), interpretación DST-aware y normalización a UTC. Queda solo verificar el DST contra datos on-chain (§2).
2. **Umbrales del emparejamiento heurístico de transferencias** (§7): ventana temporal máxima, tolerancia de importe. → definir con más datos reales.
3. **Política ante duplicados exactos** (§9): criterios para distinguir repetición legítima de error.
4. ~~**Otras variantes de exportación**~~ → ✅ Ver §12: variante `en_full_export` (botón "CSV (Exportación Completa)", locale inglés) documentada y soportada por `tools/ct_audit.py`.

---

## 12. Variante "en_full_export" — botón "CSV (Exportación Completa)", locale inglés

**Fuente:** exportación real de un usuario (proyecto `binance-2025`, 1.821 operaciones, 2024-04 → 2025), obtenida desde **Transacciones → Export → CSV (Exportación Completa)** (ver `WEB_APP_GUIDE.md`).
**Fecha de validación:** 2026-07-03.

Esta variante es **distinta** de la "Trade Table" simple descrita en §1-10 (esa se obtiene con el botón **CSV** normal). Diferencias:

| | `es_trade_table` (§1-10) | `en_full_export` (esta sección) |
|---|---|---|
| Botón de export | CSV | **CSV (Exportación Completa)** |
| Nº de columnas | 16 | **13** |
| Idioma cabecera/valores | Español (`Tipo`, `Depósito`, `Retirada`…) | **Inglés** (`Type`, `Deposit`, `Withdrawal`…) |
| Formato de fecha | `DD.MM.YYYY HH:MM:SS` | **`YYYY-MM-DD HH:MM:SS`** |
| Columnas de dirección (`From/To/Sell Address`) | Sí (3) | **No** |
| Columna `LPN` | No existe | Sí (presente pero **vacía en el 100%** de la muestra; propósito no confirmado) |
| Cobertura de `Tx-ID` en transferencias | ~16-24% | **~85%** (1.552/1.821 filas en la muestra; más alta también fuera de depósitos/retiradas) |

### Columnas (en orden), 0-indexadas

| # | Nombre en CSV | Equivale a (§1) |
|---|---|---|
| 0 | `Type` | `Tipo` |
| 1 | `Buy` | `Compra` |
| 2 | `Cur.` | `Cur.` (moneda de Buy) |
| 3 | `Sell` | `Venta` |
| 4 | `Cur.` | `Cur.` (moneda de Sell) |
| 5 | `Fee` | `Comisión` |
| 6 | `Cur.` | `Cur.` (moneda de Fee) |
| 7 | `Exchange` | `Intercambio` |
| 8 | `Group` | `Grupo` |
| 9 | `Comment` | `Comentario` |
| 10 | `Date` | `Fecha` |
| 11 | `LPN` | *(sin equivalente; siempre vacía en la muestra)* |
| 12 | `Tx-ID` | `Tx Hash` |

### Tipos de transacción observados (columna `Type`)

| Valor CSV | Frecuencia (muestra) | Equivale a (§3) |
|---|---:|---|
| `Trade` | 667 | `Operación` |
| `Other Fee` | 404 | `Otras comisiones` |
| `Reward / Bonus` | 324 | `Recompensa / Bonificación` |
| `Staking` | 150 | `Staking` |
| `Derivatives / Futures Loss` | 114 | `Pérdidas por Derivados / Futuros` |
| `Derivatives / Futures Profit` | 82 | `Beneficio de Derivados / Futuros` |
| `Deposit` | 40 | `Depósito` |
| `Withdrawal` | 25 | `Retirada` |
| `Interest Income` | 9 | `Ingresos por intereses` |
| `Income` | 5 | `Ingresos` |
| `Spend` | 1 | `Gasto` |

### Detección automática

`tools/ct_audit.py` detecta la variante por la cabecera (`header[0] == "Type"` + ≥13 columnas → `en_full_export`; `header[0] == "Tipo"` + ≥16 columnas → `es_trade_table`) y ajusta índices de columna, tipos de depósito/retirada y formato de fecha en consecuencia (`detect_format()`/`configure_format()`). Si la cabecera no coincide con ninguna variante conocida, **falla explícitamente** en vez de adivinar (ADR-009).

> 🔧 Todas las reglas de §5-9 (comisiones, comentario sucio, colisión de tickers, duplicados) aplican igual en esta variante; solo cambian los literales de columna/tipo y el parseo de fecha.
