# Guía de uso de la web de CoinTracking (remediación y operación)

**Tipo:** Conocimiento operativo destilado de fuentes oficiales de CoinTracking
**Fuentes:** Centro de ayuda oficial (URLs al pie y en `reference/CATALOG.md`)
**Última verificación:** 2026-07-03
**Vigencia:** la **interfaz de CoinTracking cambia**. Los pasos concretos (nombres de menús/botones) pueden variar — **antes de dar instrucciones clic a clic, verifica los pasos vigentes** en el artículo oficial citado (ADR-008). No inventes rutas de menú (ADR-009).
**Estado:** Destilado y reelaborado (no copia verbatim)

Esta guía permite al agente **acompañar a un usuario novato** a operar la web de CoinTracking, sobre todo para **corregir** los problemas que detecta la auditoría y para **generar los informes**. Explica en pasos sencillos, uno a uno, y en lenguaje llano (ver estilo de guía en `CLAUDE.md`).

> 🔒 **Regla (ADR-009):** cuando vayas a guiar una acción en la web, **abre y lee el artículo oficial citado** para confirmar los pasos actuales antes de instruir. Si no puedes verificarlos, dilo y no improvises la ruta.

---

## La página de Transacciones (`enter_coins.php`)

*(Verificado sobre la interfaz real, 2026-07. La UI puede cambiar; ADR-008.)*

Es donde se ven, filtran, editan, fusionan y borran operaciones. Se accede desde el menú lateral **Transacciones**. Pestañas superiores: **Saldo por Exchange**, **Saldo actual**, **Lista de operaciones**.

**Botones de acción** (sobre la tabla):
- **Nueva** — añadir una operación manual.
- **Editar** — editar la(s) fila(s) marcada(s).
- **Duplicar** — duplicar la marcada.
- **Fusionar para Operación** — combinar un **depósito + una retirada** seleccionados en una sola operación (trade). Útil al reconstruir transferencias.
- **Eliminar** — borrar la(s) marcada(s).
- **Borrado / Edición masivos** — acciones en lote.

**Filtrar (Búsqueda avanzada):** caja de búsqueda libre + filtros por columna (Type, Buy, Cur., Sell, Cur., Fee, Cur., **Exchange**, Group, **Comment**, **Date**), cada uno con conmutador **Smart / RegEx**. Ejemplo usado en la auditoría de Coinbase: Exchange = `Coinbase` y búsqueda `16.03.2024`.

**Seleccionar filas:** casilla a la izquierda de cada fila (Ctrl para varias, Shift para un rango). El pie indica *"Mostrando 1 a N de un total de N (filtradas de un total de 1.828 entradas)"* — sirve para confirmar el alcance del filtro antes de borrar.

**Exportar** (botón *Export*): Copiar al Portapapeles, Imprimir, **CSV**, Excel, PDF, **PDF / CSV (Exportación Completa)**, JSON (limitado), XML (limitado), HTML. → El **CSV** es la "Trade Table" que usa el agente (ver `CSV_FORMAT.md`).

**Opciones de la tabla (pie):** *Vista de tabla* (Extendida), *Edición integrada* (edición en línea sobre la propia tabla), *Autocompletar*, *Modo CSP*.

**Columnas:** las mismas del CSV export — Tipo, Compra, Cur., Venta, Cur., Comisión, Cur., Intercambio, Grupo, Comentario, Fecha (+ columnas de direcciones y Tx Hash). Todas ordenables.

> 🔧 **Para guiar un borrado/edición:** filtra por Exchange + fecha (formato `DD.MM.AAAA`), confirma con el pie cuántas filas quedan, marca las casillas correctas y usa **Eliminar** (o **Editar**). Recomienda **copia de seguridad** antes.

---

## Mapa de remediación: hallazgo de auditoría → cómo arreglarlo

| Hallazgo (ver informe de auditoría) | Acción en la web | Artículo oficial |
|-------------------------------------|------------------|------------------|
| Venta sin base de coste / "compra importada como depósito" | Cambiar el tipo de la operación de **Depósito → Trade**, o editar el valor del activo | Editar/eliminar operaciones · Editar valor del activo · Avisos de ganancias |
| Transferencia huérfana (falta un lado) | Registrar el **lado que falta** (retirada o depósito), con la retirada **antes** que el depósito | Entrada de depósitos/retiradas/transferencias |
| Duplicados | Revisar el informe de **duplicados** y borrarlos (individual o en lote) | Duplicate Transactions · Bulk Delete |
| Depósito de fiat faltante (fiat negativo) | Añadir el **depósito de fiat** (p. ej. ingreso SEPA) que falta | Cómo introducir operaciones |
| Saldo/coste incorrecto puntual | **Editar valor del activo** (precio de compra/venta) de la operación | Editar valor del activo |
| Muchas operaciones a corregir igual | **Edición masiva** (tipo, hora, precio) | Bulk Edit |
| Datos incompletos de un exchange | Revisar la **importación** (API/CSV) y el informe de transacciones faltantes | Importación · Missing Transactions Report |

---

## Tareas frecuentes (en lenguaje de usuario)

### 1. Corregir una "venta sin base de coste"
Suele deberse a que una **compra** se importó como **Depósito** (que no lleva coste) en vez de **Trade**. Dos vías:
- **Cambiar el tipo** de la operación a *Trade* e indicar por cuánto se compró; o
- **Editar el valor del activo** de la operación: en la página de **Transactions/Operaciones**, seleccionar la operación → *Edit* → *Edit Asset Value* → ajustar los valores de compra y venta.
- Fuentes: [Editar y eliminar operaciones](https://cointracking.freshdesk.com/en/support/solutions/articles/29000044546-how-to-edit-and-delete-transactions) · [Editar valor del activo](https://cointracking.freshdesk.com/en/support/solutions/articles/29000033219-edit-asset-value-manually) · [Avisos en el informe de ganancias](https://cointracking.freshdesk.com/en/support/solutions/articles/29000007206-warnings-in-the-capital-gains-report)

### 2. Registrar/corregir una transferencia entre cuentas propias
Una transferencia necesita **los dos lados**: la **retirada** en la cuenta de origen y el **depósito** en la de destino. Importante: la **retirada debe tener fecha anterior o igual** al depósito, o no se traslada la base de coste. Usar los tipos *Depósito*/*Retirada* solo para monedas que ya se poseen con coste (para regalos/recompensas, usar el tipo correcto).
- Fuente: [Entrada de depósitos, retiradas y transferencias](https://cointracking.freshdesk.com/en/support/solutions/articles/29000007201-entering-deposits-withdrawals-and-transfers-between-exchanges-wallets)

### 3. Eliminar duplicados
Usar el informe de **duplicados** para localizarlos y borrarlos; si son muchos, borrado en lote. Recuerda: filas idénticas **no siempre** son error (comisiones/recompensas recurrentes) — revisar antes de borrar.
- Fuentes: [Duplicate Transactions](https://cointracking.freshdesk.com/en/support/solutions/articles/29000048918-duplicate-transactions) · [Transacciones duplicadas recurrentes](https://cointracking.freshdesk.com/en/support/solutions/articles/29000018219-reoccurring-duplicate-transactions) · [Borrado en lote](https://cointracking.freshdesk.com/en/support/solutions/articles/29000043099-bulk-delete-transactions)

### 4. Introducir una operación manual (p. ej. un ingreso fiat que falta)
Se pueden añadir operaciones a mano en la página de introducción de operaciones; hay campos obligatorios según el tipo.
- Fuentes: [Cómo introducir operaciones](https://cointracking.freshdesk.com/en/support/solutions/articles/29000018166-how-to-enter-transactions-into-cointracking) · [Datos obligatorios (custom importer)](https://cointracking.freshdesk.com/en/support/solutions/articles/29000032507-custom-importer-and-mandatory-data-to-add-transactions)

### 4bis. Bloque-resumen para introducir una operación manual ("resumen para copiar")

Cuando la tarea implique **crear, modificar o corregir una operación manual** en CoinTracking, tras explicar en lenguaje llano qué hay que hacer y por qué (regla general de este documento y de `CLAUDE.md`), cierra con un **bloque compacto** que el usuario pueda copiar campo a campo sobre el formulario, en el mismo orden que aparece en la tabla/formulario (`Tipo, Compra, Cur., Venta, Cur., Comisión, Cur., Intercambio, Grupo, Comentario, Fecha` — verificado, ver arriba §"Columnas").

> 🔑 **Regla de uso (importante):** este bloque **nunca sustituye** la explicación en lenguaje llano ni el paso a paso — es un añadido al final, para que quien ya entendió qué va a hacer pueda copiarlo rápido sin releer párrafos. Adáptalo siempre a la tarea concreta que se está pidiendo; no lo fuerces si no aporta claridad (p. ej. en una pregunta puramente conceptual no hace falta).

**Estructura:**

```
[ <Tipo> | <Fecha DD.MM.AAAA HH:MM:SS> ] [ <campos principales según el tipo> ] [ Intercambio: … | Grupo: … | Comentario: … ]
```

- El tercer bloque (Intercambio/Grupo/Comentario) es opcional campo a campo: omite lo que no aplique.
- No inventes valores que el usuario no haya dado.
- Si la operación puede generar un aviso de CoinTracking (balance negativo, "no hay compra adecuada para esta venta", etc.), dilo **después** del bloque, nunca dentro.

**Campos principales por tipo — confirmados contra tus datos reales (`CSV_FORMAT.md` §3/§12):**

| Tipo | Campos principales |
|---|---|
| `Operación` (Trade) | `Compra: X CUR \| Venta: X CUR \| Comisión: X CUR` (comisión opcional) |
| `Depósito` | `Cantidad: X CUR` (entra por `Compra`) |
| `Retirada` | `Cantidad: X CUR` (sale por `Venta`, + comisión si aplica) |
| `Ingresos` (no "Ingreso") | `Cantidad: X CUR` |
| `Ingresos por intereses` | `Cantidad: X CUR` |
| `Gasto` | `Cantidad: X CUR` |
| `Staking` | `Cantidad: X CUR` |
| `Recompensa / Bonificación` | `Cantidad: X CUR` |
| `Otras comisiones` | `Cantidad: X CUR` |

> ⚠️ **`[VERIFICAR]` — no confirmados contra el desplegable real del formulario ni contra tus datos:** `Donación`, `Minería`, `Airdrop`, `Regalo`. Es plausible que existan con esos nombres u otros parecidos, pero antes de usarlos en una instrucción clic a clic, confírmalos en la sesión (abrir la página de introducción manual, o el artículo oficial "Cómo introducir operaciones") en vez de asumirlos.

> 🔴 **Corrección importante — no existe un tipo único "Transferencia":** una transferencia entre cuentas propias del usuario se registra como **dos operaciones separadas**, una `Retirada` en la cuenta de origen y un `Depósito` en la de destino, con la retirada en fecha igual o anterior (§2 de esta guía). Nunca uses un bloque `[ Transferencia | … ]` como si fuera un tipo real de CoinTracking — usa dos bloques, uno de cada tipo.

**Ejemplos:**

```
[ Operación | 04.07.2026 10:35:58 ] [ Compra: 0.25 BTC | Venta: 20000 EUR | Comisión: 10 EUR ] [ Intercambio: Binance | Grupo: Spot | Comentario: Compra BTC ]

[ Depósito | 01.07.2026 09:15:00 ] [ Cantidad: 0.80 BTC ] [ Intercambio: Ledger ]

[ Retirada | 08.07.2026 16:05:12 ] [ Cantidad: 500 USDT ] [ Intercambio: Binance | Comentario: Envío a Ledger ]
```

Transferencia Ledger → Binance de 0,50 ETH (dos tareas, no una):

```
[ Retirada | 08.07.2026 16:10:00 ] [ Cantidad: 0.50 ETH ] [ Intercambio: Ledger ]
[ Depósito | 08.07.2026 16:10:40 ] [ Cantidad: 0.50 ETH ] [ Intercambio: Binance ]
```

### 5. Ediciones masivas (tipo, hora, precio)
Para corregir muchas operaciones a la vez: cambiar tipo, ajustar hora, fijar precio por unidad, etc.
- Fuentes: [Bulk Edit y Delete](https://cointracking.freshdesk.com/en/support/solutions/articles/29000043331-bulk-edit-and-delete) · [Cambiar tipo](https://cointracking.freshdesk.com/en/support/solutions/articles/29000043132-bulk-edit-change-trade-type) · [Ajustar hora](https://cointracking.freshdesk.com/en/support/solutions/articles/29000043229-bulk-edit-adjust-trade-time) · [Fijar precio por unidad](https://cointracking.freshdesk.com/en/support/solutions/articles/29000043231-bulk-edit-set-price-per-unit)

### 6. Validar la cuenta y transacciones faltantes
CoinTracking ofrece herramientas de validación y un informe de transacciones faltantes (depósitos sin su retirada emparejada).
- Fuentes: [Cómo validar mi cuenta](https://cointracking.freshdesk.com/en/support/solutions/articles/29000035339-how-to-validate-my-account-) · [Missing Transactions Report](https://cointracking.freshdesk.com/en/support/solutions/articles/29000048812-missing-transactions-report)

### 7bis. Descargar el CSV de un informe de "Analysis" (Missing Transactions, Double-Entry, Duplicados, Realized/Unrealized, Balance by Exchange…)

Cuando la auditoría necesita uno de los informes adicionales de `DOCUMENT_CHECKLIST.md` §A (no la Trade Table), guía al usuario así:

1. **Patrón general** (visto en la página de Transacciones, §"La página de Transacciones" arriba): casi todos los informes de CoinTracking tienen, sobre la propia tabla del informe, un botón **Export/Exportar** con opciones **CSV**, Excel, PDF. El patrón es: abrir el informe → botón Export → elegir CSV → descargar.
2. **Dónde está cada informe (confirmado contra el artículo oficial a fecha 2026-07-03; la interfaz puede haber cambiado — reverifica antes de instruir clic a clic, ADR-008/009):**
   - **Balance by Exchange** — página propia: `cointracking.info/balance_by_exchange.php` (accesible también desde el menú de balances/holdings). Permite exportar en Excel/CSV/PDF.
   - **Realized and Unrealized Gains Report** — menú **Reports → Realized & Unrealized Gains**.
   - **Double-Entry List** — dentro de **Analysis → Transactions**; el botón de exportación está sobre la tabla del libro (ledger).
   - **Missing Transactions Report** y **Duplicate Transactions** — están en la sección **Analysis** del menú, pero el artículo oficial no detalla el nombre exacto del botón/submenú; **antes de dar el paso clic a clic, abre el artículo oficial citado en `DOCUMENT_CHECKLIST.md` y confírmalo en la sesión** (no lo inventes).
3. **Si no encuentras el botón de export en un informe concreto:** pide al usuario una captura de pantalla del informe (sin acciones que no puedas verificar) o remítelo al artículo oficial correspondiente para que confirme el nombre exacto del botón en su versión de la interfaz.
- Fuentes: artículos oficiales enlazados en `DOCUMENT_CHECKLIST.md` §A.1 (cada informe tiene el suyo).

### 7. Generar el informe fiscal de España (con FIFO)
En la configuración del **Tax Report**: elegir **país = España** y **método de cálculo = FIFO** (España usa FIFO; ver `../taxation/spain/CAPITAL_GAINS.md` §4). El informe muestra: ganancias patrimoniales realizadas (spot, NFT, derivados), rendimientos generales (staking, minería, airdrops), ingresos no gravables, otros pagos, y el detalle de operaciones. Incluye una sección sobre el **Modelo 721**.
- Fuente: [Guía del informe fiscal (España)](https://cointracking.freshdesk.com/en/support/solutions/articles/29000045612-guide-to-navigating-and-understanding-your-tax-report-spain-)
- 🔑 Este informe (FIFO/España) es la **calculadora determinista** de referencia para las cifras fiscales (ADR-006): el agente prepara/limpia los datos y verifica; la cifra vinculante sale de aquí o del asesor.

---

## Fuentes

Centro de ayuda oficial de CoinTracking (contenido propietario; enlaces públicos). Índice completo en `reference/CATALOG.md`. Ante cualquier paso clic a clic, **verificar en el artículo** por si la interfaz cambió.
