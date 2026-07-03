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
