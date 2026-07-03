
## TODO:
- Pendiente poder trabajar por proyectos. (HECHO OK, fase 1 — DECISIONS.md#ADR-013 corregido: proyecto obligatorio, puerta de entrada al arrancar chat, `USER_INPUT/<proyecto>/` y `reports/output/<proyecto>/`, caso existente migrado a `agp`. MCP pospuesto explícitamente, no aislado por proyecto todavía — ver limitación conocida en el ADR.)

- MCP aislado por proyecto en caliente. (HECHO OK, 2026-07-03 — DECISIONS.md#ADR-016. Se descartó la idea original de un `.env_cointracking_mcp` reescrito al cambiar de proyecto: el binario Go no leía ficheros de config y, además, el servidor solo arranca una vez por sesión, así que reescribir un fichero no evitaba el reinicio. En su lugar: nuevo tool MCP `cointracking_switch_project(project_name)` que cambia el proyecto activo del proceso ya arrancado — flush+close de la caché saliente, abre/crea la del proyecto entrante, sin tocar `.mcp.json` ni reiniciar Claude Code. Ambas skills lo llaman en su Paso -1. Código: `cointracking-mcp/internal/tools/{app,switch_project}.go`, test `TestSwitchProject`.)



- revisar total fiat recibido 31.000,00 € (HECHO OK, 2026-07-03 — cotejado contra `H:\cripto-agp\crypto-project` y, al detectar una nota incompleta, contra evidencia primaria: justificantes bancarios PDF, extracto BBK completo (`movimientos.xls`) e historial de BingX (`Transaction_History_*.xlsx`). Los 31.000,00 € que mostraba CoinTracking eran incompletos: faltaban **3.000,00 €** de 3 transferencias SEPA a BingX (14.08, 21.08 y 01.09.2025 — no 2, como decía la nota `flujo-fiat-y-depositos.md`) que solo estaban registradas como llegada de USDT (2 de ellas) o ni eso (la 3ª), sin su lado EUR. Corregido en CoinTracking (Depósito EUR + Trade EUR→USDT por cada una, mismo patrón que ya usa Binance) y verificado en vivo por MCP tras invalidar caché: `Total Fiat Recibido` = **34.000,00 €** = 24.500 Binance SEPA + 6.500 Coinbase [2.500 PayPal + 4.000 SEPA dic-2024] + 3.000 BingX SEPA ago-sep 2025. Detalle en `reports/output/agp/REGISTRO-CAMBIOS.md` §"2026-07-03 — Hallazgo de una 3ª transferencia BingX y cierre definitivo (34.000,00 €)". No afecta a ninguna declaración ya presentada. Efecto colateral detectado y **no resuelto todavía**: el saldo mostrado de BingX subió ~1.000 € por encima de lo real al añadir la 3ª operación — ver nueva tarea abajo.)

- Reimportación completa y fresca de BingX en CoinTracking. (HECHO OK, 2026-07-03 — sesión larga con varios intentos fallidos documentados en `REGISTRO-CAMBIOS.md` §"Reconstrucción completa de BingX", hasta dar con el método bueno: borrar BingX, reimportar **solo** Order History 2024+2025 —Spot+Standard_Futures+USDM, nunca junto con Transaction History, que duplica— y añadir a mano los 8 depósitos on-chain + 3 pares Depósito EUR/Trade EUR→USDT vía un único CSV de importación genérica. Resultado: cuenta sin duplicados, todas las monedas exactas salvo **USDT con ~324 de diferencia sin origen identificado** (real 72,94 vs calculado 397,40) — pendiente aparte, con 3 vías concretas anotadas en `audit_state.md` para resolverlo con calma.)

- Cerrar el descuadre de ~324 USDT en el saldo de BingX. (HECHO OK, 2026-07-03 — resuelto reconstruyendo desde los saldos acumulados del ledger real de BingX. BingX cuadra ahora al céntimo: **72,93993932 USDT = saldo real 72,94**. Causa raíz: Futuros importados 2 veces + ~695 USDT perdidos en sub-cuenta Copy Trading que BingX no exporta. Solución: borrar todos los Futuros, conservar depósitos/fiat/spot, y añadir el resultado neto de Futuros por año (del ledger real) + una entrada de reconciliación de −694,67 marcada como "Lost" NO deducible. Detalle en `REGISTRO-CAMBIOS.md` §"CIERRE DEFINITIVO de BingX". BingX queda estático — el usuario no volverá a operar ahí.)
  - ⚠️ **Pendiente menor (opcional):** si se quisiera usar fiscalmente la pérdida de ~695 USDT de Copy Trading, pedir a BingX ese export para documentarla; mientras, queda como "Lost" no deducible (lado seguro).
  - **Futuro:** cuando el usuario traspase el residual de BingX a otra wallet, añadir esa única retirada a mano (trivial).

- Minimiza llamadas MCP: agrupa verificaciones al final por lote, evita validaciones una a una. (HECHO OK — DECISIONS.md#ADR-010 punto 7, knowledge/cointracking/MCP_API.md, .claude/skills/audit-cointracking/SKILL.md Paso 3)
- Usa y respeta caché de CoinTracking según ADR-010; invalídala cuando haya cambios de datos. (HECHO OK — ya cubierto por ADR-010 y el servidor MCP propio con `cointracking_invalidate_cache`/`cointracking_cache_stats`; skill actualizada para llamarlas explícitamente)

- Define checklist de documentación mínima a pedir al usuario (Trade Table, Missing Transactions, Double-entry-detailed, Realized/Unrealized, saldo por exchange). (HECHO OK — knowledge/cointracking/DOCUMENT_CHECKLIST.md)
- Genera documentación de usuario para ese flujo de como hacer para descargar el csv de cualquiera de los análisis. (HECHO OK — knowledge/cointracking/WEB_APP_GUIDE.md §7bis, verificado contra artículos oficiales; lo no confirmable queda marcado para reverificar antes de instruir clic a clic)

- Revisa completamente la documentación pues veo mucho desactualizado y referencias a framework y no a agente. (HECHO OK — FOUNDATION.md reescrito; docs/GLOSSARY.md, knowledge/faq/INDEX.md, knowledge/taxation/spain/INDEX.md, CAPITAL_GAINS.md, CAPITAL_INCOME.md, INFORMATIVE_OBLIGATIONS.md, COST_BASIS_AND_VALIDATION.md, CSV_FORMAT.md, cointracking/INDEX.md, CATALOG.md y cointracking-auditor.md actualizados de "motor/framework" a "agente/skill/tools/ct_audit.py". SONNET5_MCP_PROMPT.md eliminado (ya cumplió su propósito, el MCP en Go está construido). README.md y CLAUDE.md ya estaban al día.)

- Programar nuestro propio MCP en golang. (HECHO OK)




## prompt base de conocimiento — ✅ HECHO (2026-07-03)

Petición "Integrar casos ChatGPT como base curada v2" ejecutada por fases A-D.
Resultado: `knowledge/patterns/cointracking_casos_v2.yaml` (20 casos, esquema canónico).
Legacy/deprecación documentada en `knowledge/patterns/INDEX.md` y `DECISIONS.md#ADR-015`.
`AGENT_CHANGE_REQUESTS.md` marcado como resuelto. `LEEME.md` y `PROMPT_CHATGPT_AGENTE.md` eliminados (auxiliares ya absorbidos).




