---
id: KB-A2-002
title: "Base de coste, purchase pool y validación en CoinTracking"
level: A
domain: cointracking
source: "CoinTracking centro de ayuda (destilado), validado contra datos reales"
authority: official
last_verified: 2026-07-02
valid_from: 2026-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-032
  - ADR-004
  - ADR-019
  - ADR-018

related_docs:
  - CSV_FORMAT.md
  - MCP_API.md
  - knowledge/patterns/cointracking_casos_v2.yaml
  - tools/ct_audit.py

tags:
  - cointracking
  - cost-basis
  - purchase-pool
  - fifo
  - validation

notes: "Destilado de documentación oficial CT. Comportamiento puede cambiar — reverificar contra reference/CATALOG.md si es antiguo."
---

# Base de coste, "purchase pool" y validación en CoinTracking

Este documento explica, **con nuestras propias palabras**, cómo CoinTracking calcula la base de coste y detecta inconsistencias. Es conocimiento crítico para la auditoría (`tools/ct_audit.py`, skill `audit-cointracking`) y para la preparación fiscal (skill `spanish-tax-return`), porque define cómo se comporta la plataforma cuyos datos auditamos.

---

## 1. El "purchase pool" (pool de compras)

CoinTracking mantiene, por activo, un **pool de coste** que se alimenta y se consume según el tipo de operación:

| Alimentan el pool (añaden base de coste) | Consumen el pool (deducen) | No afectan al pool |
|------------------------------------------|----------------------------|--------------------|
| Compras (Trade), regalos, recompensas, ingresos, staking, airdrops | Ventas (Trade), gastos, donaciones | **Depósitos y retiradas (transferencias internas)** |

> 🔑 **Clave para la auditoría:** las **transferencias entre cuentas propias no alteran la base de coste** — son movimientos internos. Esto coincide con la fiscalidad española (ver `../taxation/spain/CAPITAL_GAINS.md` §1: las transferencias internas no son hecho imponible) y con el emparejamiento de transferencias (`CSV_FORMAT.md` §7).

El pool se divide además entre tenencias a **corto y largo plazo**, y la vista de *Gains* puede agrupar por día ("Group by Day"), lo que genera entradas segmentadas.

> ⚠️ **Diferencia de método a vigilar:** el "purchase pool" por defecto de la página *Gains* es un mecanismo de **agrupación/promedio**, no necesariamente FIFO. CoinTracking permite seleccionar el **método fiscal** (FIFO, etc.) en el informe de impuestos. Para España, que **exige FIFO** (AEAT, ver CAPITAL_GAINS §4), la skill `spanish-tax-return` debe asegurar que se compara contra el método FIFO (`price:"oldest"` en el MCP, o el Informe de Impuestos de CoinTracking), no contra la vista de pool promediado. Además, la opción de precios **"Best Prices"** en la página *Gains* refleja los valores realmente pagados (Trade Price report).

---

## 2. Base de coste y transferencias (regla del orden temporal)

CoinTracking transfiere la base de coste entre cuentas mediante el par retirada→depósito:

- Una **retirada genera** la base de coste que "viaja".
- Un **depósito debe casar** con esa retirada para **recibir** la base.
- **Si el depósito tiene una marca temporal *anterior* a la retirada, la base de coste NO se transfiere** y se genera una advertencia.

> 🔑 **Implicación crítica para ADR-005 (zonas horarias):** un orden temporal incorrecto —causado, por ejemplo, por una **mala interpretación de la zona horaria** o por relojes distintos entre exchanges— **rompe la transferencia de base de coste** y produce ganancias infladas. Esto refuerza por qué normalizamos todo a UTC de forma determinista. La auditoría (`tools/ct_audit.py`) debe verificar `retirada ≤ depósito` tras la normalización.

Regla de uso de CoinTracking: **usar depósito/retirada solo para activos que ya se poseen con base de coste**. Para monedas nuevas sin base (regalos, minería, staking, airdrops) debe usarse **el tipo de transacción correcto**, no un simple depósito.

---

## 3. Advertencias del informe de ganancias (y su causa raíz)

Estas advertencias de CoinTracking son, de hecho, **hallazgos de auditoría** que el agente debe saber interpretar y explicar (citando esta sección como fuente):

### 3.1 "No hay una compra adecuada para esta venta" (pools de compra agotados)
Significa que se vende un activo del que **no consta compra con base de coste**. Causas típicas:
- El activo aparece **solo como depósito**, sin operación previa de compra.
- Una **compra con tarjeta** se importó como *depósito* en lugar de *Trade*.
- Por el método fiscal (p. ej. FIFO), se consumió una fuente sin base de coste.
- Con **separación de lotes (Depot Separation)**, el depósito es anterior a la retirada, rompiendo la transferencia de base.

> 🔧 **Qué es "Depot/Lot separation"** (verificado 2026-07-04 contra el [centro de ayuda oficial](https://cointracking.freshdesk.com/en/support/solutions/articles/29000038566-depot-lot-separation)): opción global de la cuenta que decide si el cálculo de coste (FIFO/LIFO/HIFO) trata cada exchange/wallet como un "depot" fiscal independiente.
> - **Desactivada (valor por defecto):** CoinTracking calcula de forma **global**, combinando todos los exchanges/wallets en una sola cola FIFO por activo.
> - **Activada:** cada exchange/wallet se convierte en su propio depot, con cálculo independiente — pensado para jurisdicciones que lo exigen (p. ej. determinados casos en EE. UU.), **no la práctica habitual asumida para España** (ver `knowledge/taxation/spain/CAPITAL_GAINS.md` §4).
> - Si está activada, **requiere que toda transferencia entre cuentas del usuario esté completa** (retirada + depósito correspondiente emparejados); si falta una pata, genera exactamente el aviso "no hay una compra adecuada para esta venta" de este apartado.
> - **Acción para el auditor:** antes de fiar un cálculo de ganancias del informe de CoinTracking, comprobar en Configuración de la cuenta si esta opción está activada o desactivada, y hacerlo explícito en el informe.

### 3.2 Base de coste = 0 (o coste por unidad irreal)
Si falta la compra, CoinTracking asigna una **base de coste de 0** → **ganancia inflada** e impuestos mayores. Causas: compras registradas como *depósitos* (sin base), importaciones parciales, historiales incompletos, seguimiento de un solo lado.

**Ejemplo canónico (importación de un solo lado):**
```
Compras 1 BTC por 15.000 $ en Coinbase
Lo envías a Kraken y lo vendes por 30.000 $
Si solo se importa el depósito en Kraken (no la retirada de Coinbase):
  CoinTracking asume compra a 0 $ y ganancia de 30.000 $ (en vez de 15.000 $)
```

> 🔧 **Para la auditoría de "compras faltantes":** detectar activos vendidos sin base de coste, distinguir "compra importada como depósito" de "compra realmente ausente", y **nunca** asumir base 0 en silencio — reportarlo con evidencia (principio "el silencio no es aceptable").

### 3.3 Advertencias con FIAT extranjero
Solo la **moneda principal** de la cuenta tiene soporte completo; otras divisas FIAT pueden disparar advertencias de base de coste. Solución de CoinTracking: registrar las FIAT secundarias como *Income (no imponible)*.

---

## 4. Otros comportamientos relevantes

### 4.1 FIAT negativo (no es un error)
Un saldo FIAT negativo representa el **acumulado de fiat gastado en cripto menos las ganancias fiat realizadas**, **no** el efectivo real en el exchange. Aparece cuando se compra cripto sin haber registrado el depósito FIAT previo.

> 🔧 La auditoría **no debe** tratar un FIAT negativo como imposibilidad contable automática si su origen es la ausencia de depósitos FIAT; debe distinguir "imposible" (p. ej. vender más cripto de la que se tiene) de "artefacto de no importar depósitos FIAT".

### 4.2 Duplicados recurrentes
CoinTracking reconoce el patrón de **transacciones duplicadas recurrentes** (típico de reimportaciones de API/CSV). Coincide con nuestro hallazgo de 88 filas idénticas (`CSV_FORMAT.md` §9): la detección de duplicados (`tools/ct_audit.py`, ADR-014) debe distinguir duplicado por reimportación de repetición legítima.

### 4.3 Metodología de validación de cuentas (CoinTracking "READ FIRST")
El procedimiento oficial de saneamiento sigue este orden, que es una buena guía para la **auditoría de reconciliación** (skill `audit-cointracking`):
1. Importar **todos** los datos de **todos** los exchanges y wallets (evitar seguimiento de un solo lado).
2. Comparar con los **saldos reales**.
3. Detectar y eliminar **duplicados**.
4. Resolver pequeñas discrepancias de saldo.
5. Corregir monedas faltantes / cálculos erróneos.
6. Revisar la página de **transacciones faltantes** (depósitos sin retirada emparejada).
7. Considerar **forks e ICOs**.

---

## 4.4 `get_gains` es fiable; desconfía de reconstrucciones FIFO manuales (RESUELTO — hipótesis descartada)

**Estado: cerrado y confirmado contra el Tax Report oficial (caso real `agp2025`, 2026-07-03).**

En un caso real se detectó una brecha grande entre `cointracking_get_gains(price:"oldest")` (que debería ser FIFO) y una **reconstrucción FIFO hecha a mano** sobre `cointracking_get_trades(trade_prices=1)`: BTC +492,87 € (API) vs. +94,71 € (manual) — una diferencia de ~398 €; USDC y OM mostraban el mismo patrón a menor escala. Se llegó a sospechar una "asimetría de valoración" en permutas (qué lado de la operación se usa para valorar en EUR).

**Se contrastaron ambas cifras contra el Tax Report oficial de CoinTracking** (España, FIFO), descargado en Excel para los ejercicios 2024 y 2025 (los activos en disputa se vendieron en 2024). Resultado:

| Activo | Tax Report oficial | `get_gains(price:"oldest")` | Reconstrucción FIFO manual |
|---|---|---|---|
| BTC | 503,50 € | 492,87 € | 94,71 € |
| USDC | 554,61 € | 553,93 € | 635,61 € |
| OM | 1.027,49 € | 1.027,49 € | 1.114,89 € |

> 🔑 **Conclusión (regla a aplicar de ahora en adelante):** el Tax Report oficial coincide **casi al céntimo** con `get_gains(price:"oldest")` en los tres activos (diferencias de ~10 € compatibles con redondeo/corte de fecha). **La reconstrucción FIFO manual estaba mal en los tres casos**, y de forma más marcada en BTC. La causa más probable: un FIFO reconstruido a mano operación por operación no arrastra correctamente la base de coste a través de **cadenas de permutas cripto-cripto** (un activo comprado con otro activo, no con EUR/fiat directamente); `get_gains` de CoinTracking sí lo gestiona bien internamente.
>
> **Por tanto:** ante una discrepancia entre `get_gains(price:"oldest")` y un cálculo FIFO propio, **el valor por defecto a confiar es `get_gains`/el Tax Report oficial**, no la reconstrucción manual — a menos que se repita este mismo contraste contra el Tax Report oficial y dé un resultado distinto en ese caso concreto. La hipótesis de "asimetría de valoración por lado de permuta" (documentada antes en esta sección y en `DECISIONS.md#ADR-018`) queda **descartada como causa raíz**: coincidía en magnitud por casualidad, no explicaba el fenómeno real.

**Si vuelve a aparecer una brecha así:** antes de sospechar de `get_gains`, sospecha primero de la reconstrucción manual (sobre todo si hay cadenas de permutas cripto-cripto en el activo). Contrastar contra el Tax Report oficial (país España, método FIFO) del ejercicio correspondiente sigue siendo la forma correcta de resolver la duda — ver `WEB_APP_GUIDE.md` §7.

---

## 4.5 Fuentes de precios históricos y su fiabilidad para Hacienda

Verificado 2026-07-04 contra el [centro de ayuda oficial de CoinTracking](https://cointracking.freshdesk.com/en/support/solutions/articles/29000007214-price-sources-and-currency-settings) y cruzado con criterio profesional.

- **Precio de las criptomonedas:** CoinTracking usa un **promedio ponderado por volumen** entre varias fuentes (CoinMarketCap, Coingecko, WorldCoinIndex). **No es configurable** — no existe una opción para decir "usar solo CoinGecko" o "usar solo Binance" como fuente global.
- **Conversión de divisa (EUR/USD):** aquí **sí hay una opción configurable**, en *Account Settings*: para EUR se puede elegir entre el promedio ponderado o Bitstamp/Kraken; para USD entre el promedio ponderado o Bitstamp/Kraken/Coinbase Pro. Para el resto de divisas no hay selección.
- **Materias primas** (oro/plata, si aplica): de xe.com.
- **Fiabilidad jurídica:** la AEAT no ha homologado ninguna fuente de precios concreta (ni CoinTracking ni otra); la normativa exige un **valor de mercado razonable, consistente entre ejercicios y documentable**, no una fuente específica. Aplica tanto a ventas/permutas (`CAPITAL_GAINS.md` §5) como a la valoración a 31/12 del Modelo 721 (`INFORMATIVE_OBLIGATIONS.md` §1) — la Orden HFP/886/2023 tampoco impone una fuente de cotización.

> 🔑 **Recomendación práctica para el agente:** usar los precios de CoinTracking como base por defecto (técnicamente sólidos y ampliamente usados). Para operaciones de **importe elevado** o activos **muy ilíquidos** (memecoins, tokens recién lanzados, DeFi de bajo volumen), recomendar al usuario conservar evidencia adicional (captura del precio en el exchange donde ejecutó la operación, CSV del exchange) — no porque el precio de CoinTracking sea incorrecto, sino para reforzar la trazabilidad ante una comprobación. Mantener siempre el **mismo criterio** de un ejercicio a otro.

---

## 5. Implicaciones consolidadas para la auditoría y la preparación fiscal

| Chequeo (skill/tool) | Qué extraemos de aquí |
|-------|------------------------|
| Transferencias (`ct_audit.py`, `audit-cointracking`) | Emparejar ambos lados; garantizar `retirada ≤ depósito` tras normalizar a UTC; transferencias no tocan base de coste |
| Compras faltantes (`audit-cointracking`) | Detectar ventas sin base de coste; distinguir "compra como depósito" vs ausente; nunca asumir base 0 en silencio |
| FIFO / Fiscal (`spanish-tax-return`) | Usar FIFO (España), no el pool promediado; las transferencias trasladan base pero no tributan |
| Balances (`ct_audit.py`) | FIAT negativo puede ser artefacto de no importar depósitos FIAT, no siempre imposibilidad |
| Duplicados (`ct_audit.py`, ADR-014) | Duplicado por reimportación ≠ repetición legítima |
| Reconciliación (`audit-cointracking`) | Seguir el orden de validación: importación completa → comparar saldos → duplicados → faltantes |
| Brecha `get_gains` vs FIFO manual (`audit-cointracking`, `spanish-tax-return`) | Confía en `get_gains`/Tax Report oficial por defecto (§4.4); si dudas, contrasta contra el Tax Report oficial del ejercicio, no repitas una reconstrucción manual sin verificarla |

---

## Fuentes (CoinTracking, contenido propietario — enlaces públicos)

- [How does the CoinTracking purchase pool work?](https://cointracking.freshdesk.com/en/support/solutions/articles/29000031793-how-does-the-cointracking-purchase-pool-work-)
- [Purchase Pool, Warnings in the Purchase Pool and Margin Loss Δ](https://cointracking.freshdesk.com/en/support/solutions/articles/29000018239-purchase-pool-warnings-in-the-purchase-pool-and-margin-loss--CE-B4-delta)
- [Warnings in the tax report ("no suitable purchase to this sale")](https://cointracking.freshdesk.com/en/support/solutions/articles/29000021912-warnings-in-the-tax-report-there-is-no-suitable-purchase-to-this-sale-all-purchasing-pools-consume)
- [Warnings In The Capital Gains Report](https://cointracking.freshdesk.com/en/support/solutions/articles/29000007206-warnings-in-the-capital-gains-report)
- [Why "Cost per Unit" on Gains Page Differs from the Real Price you Paid?](https://cointracking.freshdesk.com/en/support/solutions/articles/29000029902-why-cost-per-unit-on-gains-page-differs-from-the-real-price-you-paid)
- [Entering Deposits, Withdrawals, and Transfers Between Exchanges / Wallets](https://cointracking.freshdesk.com/en/support/solutions/articles/29000007201-entering-deposits-withdrawals-and-transfers-between-exchanges-wallets)
- [Transaction type examples](https://cointracking.freshdesk.com/en/support/solutions/articles/29000018272-transaction-type-examples)
- [READ FIRST: General account imbalances](https://cointracking.freshdesk.com/en/support/solutions/articles/29000018817-read-first-general-account-imbalances)
- [My account is showing negative FIAT values](https://cointracking.freshdesk.com/en/support/solutions/articles/29000018210-my-account-is-showing-negative-fiat-values)
- [Reoccurring Duplicate Transactions](https://cointracking.freshdesk.com/en/support/solutions/articles/29000018219-reoccurring-duplicate-transactions)
