# Formato CSV de CoinTracking (Trade Table)

**Estado:** Validado contra datos reales
**Fuente:** Exportación real "Trade Table" de una cuenta CoinTracking (1.828 operaciones, rango 2024-03-01 → 2026-07-01, 6 cuentas)
**Fecha de validación:** 2026-07-02

Este documento describe el formato **real** de la exportación "Trade Table" de CoinTracking, verificado contra una exportación de producción. Conforme a ADR-004, es la referencia autoritativa para la capa de importación: ninguna suposición sobre el formato debe cerrarse sin contrastarla aquí.

> ⚠️ **Aviso de alcance:** Esta es UNA variante de exportación (Trade Table) con la configuración regional en español y separador de fecha europeo. CoinTracking permite otros idiomas, formatos de fecha y conjuntos de columnas. Las peculiaridades dependientes de configuración están marcadas como **[config]**.

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

> ⚠️ **Riesgo (identificado en ARCHITECTURE_REVIEW §7.4):** sin zona horaria explícita, dos exportaciones con distinta configuración regional producen instantes distintos para la misma operación → rompe la reproducibilidad.
>
> **Decisión pendiente (candidata a ADR):** la capa de importación debe requerir que el usuario declare la zona horaria de su exportación (o asumir UTC de forma explícita y documentada), y normalizar internamente a UTC. Ver §10.

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

> 🔑 **Hallazgo clave:** `Binance Earn` aparece como **cuenta distinta** de `Binance`. Los traspasos entre productos Earn y spot se registran como pares Depósito/Retirada entre "Binance" y "Binance Earn". El motor de transferencias debe tratar cada valor de `Intercambio` como una identidad de cuenta independiente.

---

## 5. Comisiones

- En **Operación**: 246 de 666 llevan comisión; 420 no. La moneda de comisión más común es la del lado vendido, pero no siempre.
- **La moneda de la comisión puede diferir de ambos activos operados.** Casos reales:
  - Retirada `300.62905200 USDC` con comisión `0.42 USDT`
  - Operación `BTC/USDC` con comisión `19.18 EUR`
  - Operación `AGIX/USDT` con comisión `0.00585300 ETH`

> 🔧 El modelo de dominio debe tratar la comisión como un `(importe, moneda)` **independiente**, nunca asumir que coincide con el activo comprado o vendido.

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

> 🔧 **Diseño obligado del Motor de Transferencias (en niveles):**
> 1. **Determinista fuerte:** casar por `Tx Hash` cuando exista (alta confianza).
> 2. **Heurístico:** para el resto, casar por `moneda` + `importe ≈ retirada − comisión` + ventana temporal + cuentas distintas, con puntuación de confianza.
> 3. **Reportar huérfanos** explícitamente (principio "el silencio no es aceptable").
>
> Una spec basada solo en Tx Hash habría cubierto <20% de los casos.

---

## 8. Colisión de tickers — sufijo numérico de CoinTracking

CoinTracking **desambigua símbolos repetidos añadiendo un dígito**. Tickers reales observados:
`ICP2`, `ID2`, `PEPE4`, `PRIME3`, `SEI2`, `SOL2`, `THETA2`, `WLD3`.

> 🔑 **Crítico:** `SOL2` **no es** `SOL`; `WLD3` **no es** `WLD`. Son activos distintos que comparten símbolo base. El motor de duplicados y la normalización de activos **no deben** fusionar `SOL` con `SOL2`, ni asumir que el sufijo es un error tipográfico. La identidad de activo de CoinTracking es el ticker completo, sufijo incluido.

---

## 9. Duplicados exactos — no todos son errores

- **88 filas** son 100% idénticas a otra (32 grupos de duplicados), p. ej.:
  - `Otras comisiones · 0.00009792 BNB · Binance · 19.12.2024 18:49:12` ×3
  - `Depósito · 3.07724282 RDNT · Binance · 19.06.2025 00:16:00` ×2

> 🔧 Filas idénticas **no implican** error de importación: comisiones pequeñas recurrentes, recompensas periódicas y micro-movimientos pueden repetirse legítimamente en el mismo segundo. El motor de duplicados debe **señalar para revisión**, no eliminar automáticamente, y considerar contexto (tipo, si hay una operación que las justifique).

---

## 10. Reglas de normalización (resumen para la capa de importación)

1. Leer con `utf-8-sig`, parsear **por posición** (no por nombre de columna).
2. Números: `Decimal(valor)` (ADR-002); hasta 8 decimales; sin separador de miles. Vacío → ausente, no cero.
3. Fecha: `strptime("%d.%m.%Y %H:%M:%S")` → normalizar a UTC (**requiere decidir la zona de origen**, §2).
4. `Tipo`: mapear literales español → enum canónico interno (inglés); tolerar desconocidos.
5. Comisión: `(Decimal, moneda)` independiente; puede faltar o estar en tercera moneda.
6. Activo: usar el ticker **completo** (con sufijo); nunca fusionar por símbolo base.
7. Preservar la **fila original** junto a la normalizada (trazabilidad / procedencia).
8. `Comentario`/`Grupo`: metadato no fiable; no usar como clave.

---

## 11. Decisiones abiertas que este análisis genera

Estas cuestiones deberían resolverse (candidatas a ADR o a spec de motor) antes de cerrar las specs correspondientes:

1. **Zona horaria de las fechas** (§2): ¿exigir declaración del usuario? ¿asumir UTC? → afecta a reproducibilidad. **Bloqueante para el motor de libro mayor.**
2. **Umbrales del emparejamiento heurístico de transferencias** (§7): ventana temporal máxima, tolerancia de importe. → definir con más datos reales.
3. **Política ante duplicados exactos** (§9): criterios para distinguir repetición legítima de error.
4. **Otras variantes de exportación**: este documento cubre "Trade Table" en español. Documentar otras (idiomas, "Full Data Table", API) cuando se disponga de muestras.
