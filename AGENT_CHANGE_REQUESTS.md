# Peticiones de cambio al agente

Bandeja de entrada para mejoras del **agente** (código, conocimiento, reglas, tool, skills).

- **Copilot (explotación):** si durante el uso detectas un bug, un hueco de conocimiento o una regla a cambiar, **NO lo edites** — **añade una entrada aquí** (append) y sigue. (ADR-012)
- **Claude Code (gestión):** procesa estas peticiones, aplica el cambio con gobernanza (ADR/commit) y marca la entrada como ✅ hecha.

Formato de entrada:
```
## [PENDIENTE] AAAA-MM-DD — Título breve
- **Qué:** qué falla o falta.
- **Dónde:** fichero/regla afectada (p. ej. knowledge/…, tools/ct_audit.py).
- **Evidencia:** dato/caso que lo motiva.
- **Propuesta:** (opcional) qué cambio se sugiere.
```

---

<!-- Añade nuevas peticiones debajo de esta línea -->

## [✅ HECHO] 2026-07-05 — Arquitectura jerárquica del conocimiento (ADR-033 + metadatos YAML)

- **Qué:** Formalizar el sistema de conocimiento con arquitectura jerárquica (6 niveles A-F) y metadatos YAML obligatorios, permitiendo que ADR-032 (Knowledge with Temporal Validity) sea operacional.
- **Dónde:** `adr/0033-*.md`, `knowledge/.metadata/METADATA_TEMPLATE.md`, `knowledge/INDEX_MASTER.md`, `knowledge/.metadata/MIGRATION_PLAN.md`, estructura de directorios.
- **Evidencia:** Propuesta de usuario (2026-07-05) solicitando "Sistema de Conocimiento Jerárquico" con niveles de autoridad (A=oficial, B=operativo, C=casos, D=auxiliar, E=referencia, F=governance).
- **Propuesta:** Implementar en 3 fases: Fase 1 (estructura + formalización), Fase 2 (reorganizar directorios), Fase 3 (metadatos + documentos nuevos).
- **Resuelto (2026-07-05, Fase 1):** 
  - Crear ADR-033: operacionaliza arquitectura de 6 niveles + metadatos YAML (validado contra ADR-032, ADR-031, CLAUDE.md)
  - Crear `knowledge/.metadata/METADATA_TEMPLATE.md`: esquema YAML estándar con 11 campos obligatorios
  - Crear `knowledge/INDEX_MASTER.md`: mapa navegable de todos los niveles (estado actual + brechas)
  - Crear `knowledge/.metadata/MIGRATION_PLAN.md`: plan detallado de Fase 2-3 (no bloqueante, iterativo)
  - Crear estructura vacía de directorios: `authorities/`, `official/`, `behavioral/`, `cases/`, `procedures/`, `checklists/`, `decision-trees/`, `reference/` (sin mover archivos aún)
  - Commit: 491b93f
- **Próximos pasos:**
  - Fase 2 (próxima sesión): Reorganizar directorios, crear índices cruzados, actualizar referencias
  - Fase 3 (iterativo): Agregar metadatos YAML, convertir YAML a `.md`, crear procedimientos/patrones/checklists

## [✅ HECHO] 2026-07-03 — Precondición explícita de artefacto para cerrar la cifra anual exacta en `spanish-tax-return`
- **Qué:** el playbook permite avanzar con clasificación fiscal aunque no exista en el workspace el artefacto mínimo para cerrar la cifra anual exacta de base del ahorro (Tax Report oficial del ejercicio, hoja `Resumen`), dejando el cierre bloqueado al final.
- **Dónde:** `.claude/skills/spanish-tax-return/SKILL.md` (Paso 3 y Paso 6) y checklist operativo de salida del informe.
- **Evidencia:** caso real `agp2025` (03.07.2026): se pudieron cerrar eventos, recompensas y derivados, pero no la cifra total exacta anual porque no estaba adjunto el Tax Report 2025 completo; solo había cierre parcial BTC/USDC/OM en `REGISTRO-CAMBIOS.md`.
- **Propuesta:** añadir gate explícito: para marcar el informe como "listo para presentar", exigir presencia del Tax Report oficial del año (o su cifra `Resumen` documentada con evidencia) antes de cerrar la sección de base del ahorro.
- **Resuelto (2026-07-03):** añadido el gate explícito en el Paso 3 (no cerrar la cifra de base del ahorro sin el artefacto, aunque sí se puede avanzar con el resto del informe) y recordatorio en el Paso 6 (no marcar "listo para presentar" sin ese gate satisfecho). En este mismo caso, el usuario aportó el Tax Report 2025 en la misma sesión y el informe quedó cerrado (`2026-07-03_declaracion_2025.md` §7). Decisión registrada en **DECISIONS.md#ADR-021**.

## [✅ HECHO] 2026-07-03 — Verificar/normalizar el comportamiento de `cointracking_get_historical_summary` con `start/end`
- **Qué:** en preparación de renta 2025 (`agp2025`), la llamada `cointracking_get_historical_summary(start=1735686000, end=1767221999)` devolvió serie histórica diaria de 2025, pero también un punto final en fecha actual (2026-07-03), lo que sugiere que el parámetro `end` podría no aplicarse de forma estricta o que la semántica temporal no está documentada.
- **Dónde:** conocimiento `knowledge/cointracking/MCP_API.md` (documentar semántica real) y/o servidor `cointracking-mcp` (si es bug real de filtrado).
- **Evidencia:** respuesta MCP del 03.07.2026 en proyecto `agp2025`: últimos puntos incluyen `2025-12-30T23:00:00Z` y además `2026-07-03T15:31:59Z` pese a `end=1767221999`.
- **Propuesta:** añadir test de integración y especificar contrato exacto de fechas (`inclusive/exclusive`, timezone y granularidad diaria) para evitar interpretaciones erróneas en chequeos de Modelo 721.
- **Resuelto (2026-07-03):** revisado `cointracking-mcp/internal/tools/historical_summary.go` y `cached.go` — el servidor reenvía `start`/`end` sin modificar y la clave de caché los incluye, así que **no es un bug de nuestro código**; lo más probable es que la propia API de CoinTracking añada un punto "actual" adicional (no confirmado contra documentación oficial, no hay artículo público). Documentado como advertencia empírica en `MCP_API.md` con mitigación práctica: filtrar la serie por fecha en el consumidor para cualquier corte exacto (Modelo 721). Añadida la misma advertencia al Paso 5 de `spanish-tax-return/SKILL.md`. No se tocó el servidor Go (causa no confirmada). Decisión registrada en **DECISIONS.md#ADR-020**.

## [✅ HECHO] 2026-07-03 — Chequeo automático de discrepancia FIFO por asimetría de valoración (`trade_prices`)
- **Qué:** falta un control específico que detecte cuando la ganancia de `get_gains(price:"oldest")` diverge fuertemente de una reconstrucción FIFO desde `get_trades(trade_prices=1)` por usar lados de valoración distintos (`buy_value_in_cur` vs `sell_value_in_cur`) en permutas cripto-cripto.
- **Dónde:** `tools/ct_audit.py` (nuevo chequeo opcional de coherencia de ganancias) y playbook `.claude/skills/audit-cointracking/SKILL.md` (paso explícito de diagnóstico cuando hay brecha material).
- **Evidencia:** caso real `agp2025` (2026-07-03): BTC `get_gains` = `+492,87 EUR` vs reconstrucción FIFO manual ~`+94,71 EUR`; diferencia ~`398 EUR`. En los 37 trades BTC, la suma de `buy_value_in_cur - sell_value_in_cur` = `+397,72 EUR` (mismo orden de magnitud). Comisiones, duplicados y tipado quedaron descartados como causa principal.
- **Propuesta:** añadir al auditor un "delta bridge" reproducible por activo que descomponga la diferencia en: (a) efecto de comisiones, (b) efecto de transferencias, (c) efecto de elegir buy/sell value por operación; y que marque automáticamente `[VERIFICAR]` si la brecha supera umbral configurable.
- **Resuelto (2026-07-03) — alcance reducido conscientemente:** documentado el hallazgo como hipótesis empírica (no confirmada por CoinTracking) en `COST_BASIS_AND_VALIDATION.md` §4.4, con el recipe manual de diagnóstico, y añadido un sub-paso explícito en la fase 6 del Paso 1 de `audit-cointracking/SKILL.md` para aplicarlo ante brechas materiales. **No se automatizó en `tools/ct_audit.py`**: ese tool es determinista sobre el CSV (esquema verificado), y este chequeo depende de campos de la respuesta JSON de `get_trades(trade_prices=1)` cuyo esquema exacto no está verificado en `MCP_API.md` — automatizarlo ahora habría fijado nombres de campo sin confirmar (ADR-009). Queda como pendiente real verificar ese esquema antes de automatizar. Decisión registrada en **DECISIONS.md#ADR-018**.
- **⚠️ Corrección posterior (2026-07-03, mismo día):** el usuario contrastó BTC/USDC/OM contra el Tax Report oficial de CoinTracking (España, FIFO) de 2024 y 2025. Resultado: **el Tax Report oficial coincide casi al céntimo con `get_gains`**, y **la reconstrucción FIFO manual estaba mal en los tres activos** (no arrastraba bien la base de coste en cadenas de permutas). La hipótesis de "asimetría de valoración" queda **descartada como causa raíz**. `COST_BASIS_AND_VALIDATION.md` §4.4 se actualizó de "hipótesis abierta" a "resuelto"; nueva regla operativa: confiar por defecto en `get_gains`/Tax Report oficial, no en reconstrucciones manuales. Ver **DECISIONS.md#ADR-019** (corrige ADR-018).

## [✅ HECHO] 2026-07-03 — Endurecer auditoría con enfoque conservador (sin tocar cálculo)
- **Qué:** Falta una guía operativa estricta y uniforme para conciliación que minimice falsos positivos y recomendaciones arriesgadas (especialmente en duplicados, transferencias huérfanas y warning de purchase pool).
- **Dónde:** playbooks y documentación operativa del agente (`.claude/skills/audit-cointracking/SKILL.md`, `.claude/skills/spanish-tax-return/SKILL.md`, y/o checklists en `knowledge/cointracking/` según criterio de mantenimiento).
- **Evidencia:** prospección web del 03.07.2026 en artículos oficiales de CoinTracking confirma causas recurrentes: datos incompletos, API+CSV solapados, una sola pata de transferencia, tipados erróneos y costes base ausentes. Fuentes clave: `READ FIRST: General account imbalances`, `Duplicate Transactions`, `Missing Transactions Report`, `Validate Transactions`, `Transaction Flow Report`, `Roll Forward / Audit Report`, `Warnings in the tax report (all purchasing pools consumed)`, `Binance Import Restrictions`.
- **Propuesta:** añadir protocolo de diagnóstico en orden fijo: (1) cobertura de fuentes/periodos, (2) duplicados con verificación TX ID/Trade ID obligatoria, (3) matching de transferencias con tolerancias explícitas, (4) validación de tipos y fees en tercera moneda, (5) análisis purchase pool, (6) cierre con riesgos residuales. Incluir regla explícita de no recomendar borrado masivo sin evidencia y confirmación del usuario.
- **Resuelto (2026-07-03):** el conocimiento de fondo (purchase pool, transferencias, duplicados) ya estaba destilado en `COST_BASIS_AND_VALIDATION.md` y `CSV_FORMAT.md` — no hacía falta prospección nueva, solo reordenar. Reescrito el Paso 1 de `.claude/skills/audit-cointracking/SKILL.md` con las 6 fases en orden fijo propuestas, más la regla explícita de no recomendar borrado masivo sin evidencia y confirmación. `spanish-tax-return/SKILL.md` no se tocó (ya delega la reconciliación en `audit-cointracking`, sin duplicar el playbook). Los umbrales numéricos de tolerancia para el matching de transferencias **no se inventaron** — quedan documentados como heurística abierta (`CSV_FORMAT.md` §11.2), conforme a ADR-009. Decisión registrada en **DECISIONS.md#ADR-017**.

## [✅ HECHO] 2026-07-02 — Integrar casos ChatGPT como base curada v2
- **Qué:** Integrar el contenido de `cointracking_casos_extended.yaml` en el conocimiento del agente mediante una versión curada y homogénea (v2), manteniendo `cointracking_casos_base.yaml` como legacy temporal hasta validar la transición.
- **Dónde:** conocimiento de casos del repositorio (`cointracking_casos_base.yaml`, `cointracking_casos_extended.yaml`, y documentación relacionada en `knowledge/` y/o `docs/` según diseño final de Claude).
- **Evidencia:** el fichero extendido aporta cobertura y estructura útiles, pero presenta heterogeneidad de formato (listas inline vs bloque), campos vacíos con string vacío, y variabilidad de detalle en evidencia/diagnóstico. Handoff preparado en `reports/output/2026-07-02_handoff_integracion_casos_chatgpt.md` con proceso cerrado y DoD.
- **Propuesta:** ejecutar migración por fases: (A) normalizar esquema y tipos, (B) curar contenido y confianza, (C) versionar en v2 con convivencia controlada de legacy, (D) validar con criterios de aceptación explícitos. Al cerrar, dejar trazabilidad documental de estado legacy/deprecación.
- **Resuelto (2026-07-03):** ejecutadas las fases A-D. Resultado en `knowledge/patterns/cointracking_casos_v2.yaml` (20 casos, esquema canónico, sin campos vacíos inconsistentes, 5 categorías críticas de regresión cubiertas). Estado legacy/deprecado de `cointracking_casos_base.yaml` documentado en `knowledge/patterns/INDEX.md` y en **DECISIONS.md#ADR-015**. Ficheros auxiliares `LEEME.md` y `PROMPT_CHATGPT_AGENTE.md` eliminados (contenido absorbido por ADR-015 e INDEX.md).

## [PENDIENTE] 2026-07-05 — Crear ADRs de Capa 2 (Conciliación) — Motor operativo del agente

- **Qué:** El proyecto tiene 32 ADRs bien estructurados en gobernanza, arquitectura y principios (Capas 1, 5, 6), pero **falta la Capa 2 (Conciliación)** — el motor operativo que define qué es una auditoría correcta en CoinTracking. Sin estos ADRs, el agente tiene gobernanza excelente pero carece de especificidad operativa en el dominio crítico.
- **Dónde:** `adr/` — crear nuevos ADRs 033-040 aproximadamente para cubrir:
  1. Flujo de conciliación (pipeline invariante: importación → normalización → balances → transfers → duplicados → warnings → missing PH → holdings → FIFO)
  2. Modelo de balances (qué es un balance "correcto", cuándo es negativo, cuándo parar auditoría)
  3. Missing Purchase History (causa + detección + falsas alarmas + impacto fiscal)
  4. Transfers (cómo emparejar withdrawal/deposit/blockchain, tolerancias)
  5. Duplicados (matriz de clasificación: Trade ID, Order ID, Hash, Cantidad, Precio — más allá de "misma fecha")
  6. Holdings (validación CT vs Exchange vs Wallet vs Blockchain)
  7. Cost Basis / FIFO (operativo: cuándo confiar en CT, cuándo recalcular, discrepancias aceptables)
  8. Warnings (catálogo: gravedad, impacto, acción recomendada)
- **Evidencia:** análisis de Copilot (2026-07-05) identificó que de los 15 ADRs imprescindibles propuestos para un agente de auditoría, 8 existen (Capa 1 + principios) pero **7 faltan** (Capa 2 completa). Documentado en `docs/ADR_GAP_ANALYSIS_2026-07-05.md` con matriz de cobertura por capa.
- **Impacto:** CRÍTICO. Sin Capa 2, el agente es "robusto en gobernanza" pero no es "específico a CoinTracking". Es como tener un protocolo de calidad sin definir qué es calidad en el dominio.
- **Propuesta:** Sesión futura dedicada a diseñar e implementar los 8 ADRs de Capa 2, validados contra:
  - Documentación oficial de CoinTracking (centro de ayuda)
  - Comportamiento real de exchanges (Binance, Kraken, etc.)
  - Casos de auditoría reales del proyecto
  - Feedback de Copilot (usuario final del agente)
- **Prioridad:** ⭐⭐⭐⭐⭐ (BLOQUEANTE para que el agente sea productivo)
- **Siguiente paso:** Leer `docs/ADR_GAP_ANALYSIS_2026-07-05.md` y ejecutar en sesión futura (2026-07-06 o después).
