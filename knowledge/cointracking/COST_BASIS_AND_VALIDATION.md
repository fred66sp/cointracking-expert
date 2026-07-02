# Base de coste, "purchase pool" y validación en CoinTracking

**Tipo:** Conocimiento propio destilado de fuentes oficiales de CoinTracking
**Fuentes:** Centro de ayuda de CoinTracking (ver enlaces al pie); catálogo en `reference/CATALOG.md`
**Última verificación:** 2026-07-02
**Vigencia:** destilado del centro de ayuda de CoinTracking a 2026-07. El comportamiento de la plataforma (purchase pool, avisos, métodos) puede cambiar — reverificar contra las URLs de `reference/CATALOG.md` si esta fecha es antigua (ADR-008).
**Estado:** Destilado y reelaborado (no copia verbatim)

Este documento explica, **con nuestras propias palabras**, cómo CoinTracking calcula la base de coste y detecta inconsistencias. Es conocimiento crítico para los motores fiscal, FIFO, de transferencias y de reconciliación, porque define cómo se comporta la plataforma cuyos datos auditamos.

---

## 1. El "purchase pool" (pool de compras)

CoinTracking mantiene, por activo, un **pool de coste** que se alimenta y se consume según el tipo de operación:

| Alimentan el pool (añaden base de coste) | Consumen el pool (deducen) | No afectan al pool |
|------------------------------------------|----------------------------|--------------------|
| Compras (Trade), regalos, recompensas, ingresos, staking, airdrops | Ventas (Trade), gastos, donaciones | **Depósitos y retiradas (transferencias internas)** |

> 🔑 **Clave para nuestros motores:** las **transferencias entre cuentas propias no alteran la base de coste** — son movimientos internos. Esto coincide con la fiscalidad española (ver `../taxation/spain/CAPITAL_GAINS.md` §1: las transferencias internas no son hecho imponible) y con el emparejamiento de transferencias (`CSV_FORMAT.md` §7).

El pool se divide además entre tenencias a **corto y largo plazo**, y la vista de *Gains* puede agrupar por día ("Group by Day"), lo que genera entradas segmentadas.

> ⚠️ **Diferencia de método a vigilar:** el "purchase pool" por defecto de la página *Gains* es un mecanismo de **agrupación/promedio**, no necesariamente FIFO. CoinTracking permite seleccionar el **método fiscal** (FIFO, etc.) en el informe de impuestos. Para España, que **exige FIFO** (AEAT, ver CAPITAL_GAINS §4), el motor fiscal debe asegurar que se compara contra el método FIFO, no contra la vista de pool promediado. Además, la opción de precios **"Best Prices"** en la página *Gains* refleja los valores realmente pagados (Trade Price report).

---

## 2. Base de coste y transferencias (regla del orden temporal)

CoinTracking transfiere la base de coste entre cuentas mediante el par retirada→depósito:

- Una **retirada genera** la base de coste que "viaja".
- Un **depósito debe casar** con esa retirada para **recibir** la base.
- **Si el depósito tiene una marca temporal *anterior* a la retirada, la base de coste NO se transfiere** y se genera una advertencia.

> 🔑 **Implicación crítica para ADR-005 (zonas horarias):** un orden temporal incorrecto —causado, por ejemplo, por una **mala interpretación de la zona horaria** o por relojes distintos entre exchanges— **rompe la transferencia de base de coste** y produce ganancias infladas. Esto refuerza por qué normalizamos todo a UTC de forma determinista. El motor de transferencias debe garantizar `retirada ≤ depósito` tras la normalización.

Regla de uso de CoinTracking: **usar depósito/retirada solo para activos que ya se poseen con base de coste**. Para monedas nuevas sin base (regalos, minería, staking, airdrops) debe usarse **el tipo de transacción correcto**, no un simple depósito.

---

## 3. Advertencias del informe de ganancias (y su causa raíz)

Estas advertencias de CoinTracking son, de hecho, **hallazgos de auditoría** que nuestro framework debe reproducir de forma determinista:

### 3.1 "No hay una compra adecuada para esta venta" (pools de compra agotados)
Significa que se vende un activo del que **no consta compra con base de coste**. Causas típicas:
- El activo aparece **solo como depósito**, sin operación previa de compra.
- Una **compra con tarjeta** se importó como *depósito* en lugar de *Trade*.
- Por el método fiscal (p. ej. FIFO), se consumió una fuente sin base de coste.
- Con **separación de lotes (Depot Separation)**, el depósito es anterior a la retirada, rompiendo la transferencia de base.

### 3.2 Base de coste = 0 (o coste por unidad irreal)
Si falta la compra, CoinTracking asigna una **base de coste de 0** → **ganancia inflada** e impuestos mayores. Causas: compras registradas como *depósitos* (sin base), importaciones parciales, historiales incompletos, seguimiento de un solo lado.

**Ejemplo canónico (importación de un solo lado):**
```
Compras 1 BTC por 15.000 $ en Coinbase
Lo envías a Kraken y lo vendes por 30.000 $
Si solo se importa el depósito en Kraken (no la retirada de Coinbase):
  CoinTracking asume compra a 0 $ y ganancia de 30.000 $ (en vez de 15.000 $)
```

> 🔧 **Para nuestro motor de "compras faltantes":** detectar activos vendidos sin base de coste, distinguir "compra importada como depósito" de "compra realmente ausente", y **nunca** asumir base 0 en silencio — reportarlo con evidencia (principio "el silencio no es aceptable").

### 3.3 Advertencias con FIAT extranjero
Solo la **moneda principal** de la cuenta tiene soporte completo; otras divisas FIAT pueden disparar advertencias de base de coste. Solución de CoinTracking: registrar las FIAT secundarias como *Income (no imponible)*.

---

## 4. Otros comportamientos relevantes

### 4.1 FIAT negativo (no es un error)
Un saldo FIAT negativo representa el **acumulado de fiat gastado en cripto menos las ganancias fiat realizadas**, **no** el efectivo real en el exchange. Aparece cuando se compra cripto sin haber registrado el depósito FIAT previo.

> 🔧 El motor de libro mayor **no debe** tratar un FIAT negativo como imposibilidad contable automática si su origen es la ausencia de depósitos FIAT; debe distinguir "imposible" (p. ej. vender más cripto de la que se tiene) de "artefacto de no importar depósitos FIAT".

### 4.2 Duplicados recurrentes
CoinTracking reconoce el patrón de **transacciones duplicadas recurrentes** (típico de reimportaciones de API/CSV). Coincide con nuestro hallazgo de 88 filas idénticas (`CSV_FORMAT.md` §9): el motor de duplicados debe distinguir duplicado por reimportación de repetición legítima.

### 4.3 Metodología de validación de cuentas (CoinTracking "READ FIRST")
El procedimiento oficial de saneamiento sigue este orden, que es una buena guía para el **motor de reconciliación**:
1. Importar **todos** los datos de **todos** los exchanges y wallets (evitar seguimiento de un solo lado).
2. Comparar con los **saldos reales**.
3. Detectar y eliminar **duplicados**.
4. Resolver pequeñas discrepancias de saldo.
5. Corregir monedas faltantes / cálculos erróneos.
6. Revisar la página de **transacciones faltantes** (depósitos sin retirada emparejada).
7. Considerar **forks e ICOs**.

---

## 5. Implicaciones consolidadas para los motores

| Motor | Qué extraemos de aquí |
|-------|------------------------|
| Transferencias | Emparejar ambos lados; garantizar `retirada ≤ depósito` tras normalizar a UTC; transferencias no tocan base de coste |
| Compras faltantes | Detectar ventas sin base de coste; distinguir "compra como depósito" vs ausente; nunca asumir base 0 en silencio |
| FIFO / Fiscal | Usar FIFO (España), no el pool promediado; las transferencias trasladan base pero no tributan |
| Libro mayor | FIAT negativo puede ser artefacto de no importar depósitos FIAT, no siempre imposibilidad |
| Duplicados | Duplicado por reimportación ≠ repetición legítima |
| Reconciliación | Seguir el orden de validación: importación completa → comparar saldos → duplicados → faltantes |

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
