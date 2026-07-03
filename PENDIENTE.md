
## TODO:
- Pendiente poder trabajar por proyectos. (HECHO OK, fase 1 — DECISIONS.md#ADR-013 corregido: proyecto obligatorio, puerta de entrada al arrancar chat, `USER_INPUT/<proyecto>/` y `reports/output/<proyecto>/`, caso existente migrado a `agp`. MCP pospuesto explícitamente, no aislado por proyecto todavía — ver limitación conocida en el ADR.)

- MCP aislado por proyecto en caliente. (HECHO OK, 2026-07-03 — DECISIONS.md#ADR-016. Se descartó la idea original de un `.env_cointracking_mcp` reescrito al cambiar de proyecto: el binario Go no leía ficheros de config y, además, el servidor solo arranca una vez por sesión, así que reescribir un fichero no evitaba el reinicio. En su lugar: nuevo tool MCP `cointracking_switch_project(project_name)` que cambia el proyecto activo del proceso ya arrancado — flush+close de la caché saliente, abre/crea la del proyecto entrante, sin tocar `.mcp.json` ni reiniciar Claude Code. Ambas skills lo llaman en su Paso -1. Código: `cointracking-mcp/internal/tools/{app,switch_project}.go`, test `TestSwitchProject`.)



- revisar total fiat recibido 32.000,00 € (HECHO OK)

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




