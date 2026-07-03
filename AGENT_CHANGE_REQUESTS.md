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
