---

## §A — Informes de CoinTracking (además de MCP + Trade Table)

### A.0 — Mínimo de arranque (ofrecer siempre, como comprobación opcional)

Ya cubierto en el Paso 0 de `audit-cointracking`: el **CSV "Trade Table"** como cruce opcional del MCP. No se requiere nada más para empezar.

### A.1 — Pedir bajo demanda, según el hallazgo

| Informe de CoinTracking | Qué aporta que MCP/CSV no dan directamente | Pídelo cuando… | Artículo oficial |
|---|---|---|---|
| **Missing Transactions Report** | Detector propio de CoinTracking de depósitos/retiradas sin pareja (usa su lógica interna, no solo nuestra heurística de `ct_audit.py`) | La auditoría señala transferencias huérfanas y quieres una segunda fuente antes de recomendar borrar/crear nada | [Missing Transactions Report](https://cointracking.freshdesk.com/en/support/solutions/articles/29000048812-missing-transactions-report) |
| **Double-Entry List** | Lista de transferencias que CoinTracking ya emparejó como "doble entrada" (dos lados de un mismo movimiento), con saldo antes/después por lado | Quieres confirmar si una transferencia sospechosa ya está bien emparejada en el sistema antes de tocarla | [Double-Entry List](https://cointracking.freshdesk.com/en/support/solutions/articles/29000049939-double-entry-list) |
| **Duplicate Transactions (informe nativo)** | Detector de duplicados propio de CoinTracking, con su propio criterio (puede diferir del de `ct_audit.py`) | Hay sospecha de duplicados y quieres contrastar antes de aplicar el protocolo de ADR-014 | [Duplicate Transactions](https://cointracking.freshdesk.com/en/support/solutions/articles/29000048918-duplicate-transactions) |
| **Realized and Unrealized Gains Report** | Ganancia/pérdida realizada y no realizada **por activo** (agregada, no por operación), con coste medio por unidad y precio actual — buena vista rápida de qué activos concentran el riesgo fiscal, pero **no trae los warnings textuales** (esos están en el Tax Report/`get_gains`, no en este CSV) | Quieres una vista agregada por activo antes de profundizar operación a operación con `cointracking_get_gains` | [Realized and Unrealized Gains Report](https://cointracking.freshdesk.com/en/support/solutions/articles/29000043166-realized-and-unrealized-gains-report) |
| **Balance by Exchange** | Balance por exchange y por moneda calculado por CoinTracking desde el histórico importado (distinto del balance en vivo de la API del exchange) | Sospechas de discrepancia entre `cointracking_get_grouped_balance` y lo que el usuario ve en el propio exchange — es un problema documentado ("Balance by Exchange is different from API Live Data") | [Balance by Exchange](https://cointracking.freshdesk.com/en/support/solutions/articles/29000044124-balance-by-exchange) |
| **Balance by Currency / Current Balance** | Vista alternativa del balance actual, útil como segunda lectura rápida | Quieres una comprobación visual rápida sin llamar al MCP | [Balance by Currency](https://cointracking.freshdesk.com/en/support/solutions/articles/29000049929-balance-by-currency) |
| **Trade Statistics / Number of Trades** | Conteo de operaciones por exchange/periodo — útil para detectar de un vistazo un hueco de importación (p. ej. "cero operaciones en 2024") | Sospechas de completitud de importación (COST_BASIS §3.2) y quieres una cifra rápida antes de pedir el histórico completo por MCP | [Trade Statistics](https://cointracking.freshdesk.com/en/support/solutions/articles/29000049934-trade-statistics) |
| **Validate Transactions** ("Validate my account") | Ejecuta las validaciones internas de CoinTracking sobre la cuenta completa (más allá de lo que expone el MCP) | Al arrancar una auditoría completa (no puntual), como primera pasada antes de tus propios chequeos | [Validate Transactions](https://cointracking.freshdesk.com/en/support/solutions/articles/29000049924-validate-transactions) · [How to validate my account?](https://cointracking.freshdesk.com/en/support/solutions/articles/29000035339-how-to-validate-my-account-) |
| **Roll Forward / Audit Report** | Concilia saldo inicial + entradas − salidas = saldo final de un periodo — el informe más parecido a una auditoría contable clásica | El usuario pide explícitamente "una auditoría" formal o necesita justificar el saldo de un ejercicio ante el asesor | [Roll Forward / Audit Report](https://cointracking.freshdesk.com/en/support/solutions/articles/29000049968-roll-forward-audit-report) |
| **Value at Transaction** | Valor de mercado registrado por CoinTracking en el momento exacto de una operación concreta | Hay disputa sobre el coste de adquisición de una operación puntual (p. ej. CT-006 permutas, CT-010 airdrops) | [Value at Transaction](https://cointracking.freshdesk.com/en/support/solutions/articles/29000048905-value-at-transaction) |

**No pedir de forma rutinaria** (bajo demanda solo si el caso lo requiere específicamente): Transaction Flow Report, NFT Center, Tax-privileged Coins, Depot/Lot separation — son de nicho y solo aportan valor en casos concretos (DeFi complejo, NFTs, jurisdicciones con holding period).

> **Nota de vigencia:** el formato exacto (columnas, separador, nombre de fichero) de cada informe no está aún verificado contra una exportación real en este repositorio. Antes de asumir columnas concretas, contrastar contra el export que aporte el usuario o contra el artículo oficial correspondiente (ADR-008/009).

---


# Checklist de documentación a solicitar al usuario

**Tipo:** Conocimiento operativo (obra propia), basado en el catálogo oficial de CoinTracking
**Fuente:** `reference/CATALOG.md` (categorías "Analysis", "Tax-Reports + Realized & Unrealized Gains", "Data Validation")
**Última verificación:** 2026-07-03
**Vigencia:** el catálogo de informes de CoinTracking puede ampliarse o renombrarse; reverificar los nombres de informe en la web antes de pedirlos si esta fecha es antigua (ADR-008).

Esta guía resuelve una pregunta operativa: además del **MCP** (datos en vivo) y del **CSV "Trade Table"**, ¿qué más merece la pena pedir al usuario? Hay dos fuentes distintas, con propósitos distintos:

- **§A — Informes propios de CoinTracking** que el MCP/CSV no cubren (o cubren de forma indirecta) y que sirven de **validación cruzada** frente al mismo motor de CoinTracking.
- **§B — Información del exchange original**, a la que el agente **no tiene acceso** (ni por MCP ni por CSV de CoinTracking), y que solo el usuario puede aportar directamente desde la web/app del exchange.

No pidas todo esto de golpe (ADR-009, evitar fatiga de confirmación): usa **§A-mínimo** al arrancar una auditoría, y el resto **bajo demanda**, cuando un hallazgo concreto lo justifique.

**Cómo guiar la descarga:** cuando pidas al usuario uno de los informes de §A, explícale cómo descargarlo con `knowledge/cointracking/WEB_APP_GUIDE.md` §7bis (patrón general de exportación + ubicación conocida de cada informe, y cuándo hace falta reverificar el artículo oficial antes de dar el paso clic a clic).



## §B — Información del exchange original (fuera del alcance de CoinTracking)

El agente **no tiene acceso** al exchange en sí (ni por MCP ni por el CSV de CoinTracking, que ya es una capa derivada). Cuando un hallazgo requiere confirmar algo que solo el exchange sabe con certeza, pide al usuario que lo consulte **directamente en ese exchange** (el que corresponda al hallazgo, no siempre el mismo) y te traiga el dato (captura, extracto o export oficial):

| Qué pedir al exchange | Para qué hallazgo | Por qué CoinTracking no basta |
|---|---|---|
| **Trade ID / Order ID de la operación** (historial de operaciones/*trade history* de la cuenta, o vía su API) | Duplicados (CT-003, CT-008, CT-016, CT-019) — distinguir operaciones legítimas en el mismo segundo de reimportaciones reales, per **ADR-014** | El CSV de CoinTracking no siempre incluye `trade_id`; sin él, dos filas idénticas son indistinguibles |
| **Extracto/historial oficial de la operación** (fecha, importe, comisión exactos) | Comisiones ausentes (CT-009), coste incorrecto puntual | El fee puede faltar en la exportación de CoinTracking si el exchange no lo incluyó en su propio export |
| **Historial de "Convert" / swap interno** del exchange | Permutas complejas (CT-006) | CoinTracking puede recibir el Convert como dos filas sin vínculo; el exchange sabe que fue una sola acción |
| **Confirmación de titularidad de una dirección** (para retiradas a wallet propia) | Transferencias huérfanas / wallet externa (CT-001, CT-013) | CoinTracking no puede verificar de quién es una dirección; solo el usuario (o el explorador on-chain) lo confirma |
| **Hash de transacción on-chain** (para cruzar contra un explorador público) | Transferencias huérfanas, swaps DeFi fragmentados (CT-015) | El explorador on-chain es la única fuente que muestra *todos* los eventos técnicos de una transacción |
| **Anuncio oficial del proyecto** (rebrand, migración de red, airdrop) | Token renombrado (CT-018), airdrops (CT-010) | Ni CoinTracking ni el exchange certifican la naturaleza del evento; solo el proyecto emisor |

**Regla general:** solo pide esto cuando un hallazgo concreto de la auditoría lo exija (no de forma preventiva); explica **por qué** lo necesitas y qué falso positivo evita, en lenguaje llano (ver estilo de guía en `CLAUDE.md`). Adapta la petición al exchange concreto del hallazgo (Binance, Coinbase, BingX…): cada uno tiene su propio nombre de export y su propia web/app, y el formato exacto solo se confirma cuando el usuario lo aporte.

⚠️ **Cuidado con el tipo de export al pedir `trade_id`/Order ID (ADR-014):** muchos exchanges ofrecen más de un tipo de export (p. ej. un "extracto de cuenta"/"historial de transacciones" general por un lado, y un "historial de órdenes/trades" por otro). El primero suele **no incluir** `trade_id`; para verificar duplicados hace falta específicamente el que sí lo incluya. No asumas que un "historial de transacciones" cualquiera lo trae — confírmalo con el usuario o con la documentación del exchange en cuestión.

---

## Dónde dejar los archivos

Igual que el CSV Trade Table: en `USER_INPUT/<proyecto>/` (el proyecto activo de la conversación, ADR-013; ver `USER_INPUT/README.md`). Indícale al usuario el nombre de archivo esperado (p. ej. `CoinTracking · Missing Transactions.csv`) para que no se confunda con la Trade Table.

## Vigencia

Los nombres e IDs de artículo de este documento están verificados contra `reference/CATALOG.md` (2026-07-02/03). Si CoinTracking renombra o retira un informe, actualizar esta tabla y el catálogo en el mismo commit (ADR-008).
