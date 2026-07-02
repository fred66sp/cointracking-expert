# CLAUDE.md

Instrucciones de proyecto para Claude Code. Trabaja siempre en **español**.

## Qué es este proyecto

Un **agente de IA auditor de CoinTracking** que vive en Claude Code (ADR-006). Se apoya en una base de conocimiento propia y accede a los datos del usuario por dos vías: el **MCP de la API de CoinTracking** y/o el **CSV export** ("Trade Table"). Encuentra y explica problemas de reconciliación y fiscalidad española.

No es un SDK ni una librería de motores deterministas: eso se descartó (ver ADR-006). El "producto" es el agente + su conocimiento.

## Estructura

- `.claude/agents/cointracking-auditor.md` — el subagente auditor (rol, principios, límite de determinismo).
- `.claude/skills/audit-cointracking/` — el playbook invocable con `/audit-cointracking`.
- `knowledge/` — el **cerebro** del agente (fuente de verdad):
  - `cointracking/` — formato CSV, modelo de coste, integración MCP, catálogo de referencia.
  - `taxation/spain/` — fiscalidad IRPF (ganancias patrimoniales, FIFO, Modelo 721).
- `DECISIONS.md` — registro de decisiones (ADR-001…007). Gobernanza vinculante.
- `templates/AUDIT_REPORT.md` — plantilla de informe.
- `reports/output/` — informes generados (ignorado por git: datos sensibles).
- `.mcp.json` — arranque del servidor MCP (credenciales por `--env-file`, sin secretos en el repo).

## Cómo se usa

Invoca `/audit-cointracking`. El agente carga el conocimiento, obtiene los datos (MCP si está conectado; si no, el CSV) y produce un informe con formato **evidencia → causa → impacto → recomendación**.

## Principios (de FOUNDATION.md)

- **Basado en evidencia:** cada conclusión, respaldada por datos concretos.
- **Explicabilidad:** causa, evidencia, impacto y recomendación en cada hallazgo.
- **El silencio no es aceptable:** declara la incertidumbre; no la ocultes.
- **Nunca inventes reglas fiscales.** Si el conocimiento no cubre un caso, márcalo como `[PENDIENTE DE FUNDAMENTAR]`; no improvises. Cita la fuente (documento y sección).

## Límite de determinismo (ADR-006) — crítico

El agente es un LLM: **no es determinista**. Encuentra y explica problemas cualitativos, pero **no produce cifras fiscales vinculantes**. Toda cantidad exacta es «estimación no vinculante» salvo que provenga de un cálculo determinista.

## Convenciones

- **Idioma (ADR-001):** contenido, docstrings, comentarios y mensajes en **español**; nombres de archivos, carpetas e identificadores en **inglés**.
- **Zonas horarias (ADR-005):** fechas de CoinTracking en `Europe/Madrid` (con DST) → normalizar a UTC.
- **Privacidad:** datos reales del usuario (CSV, informes) y credenciales **nunca** se versionan (ya en `.gitignore`).
- **Commits:** en español, uno por cambio lógico; decisiones importantes → nuevo ADR en `DECISIONS.md`.
