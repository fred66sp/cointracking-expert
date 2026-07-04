# ADR-012: División de responsabilidades (Claude Code gestiona, Copilot explota)

**Status:** Accepted

**Date:** 2026-07-02

## Context

El proyecto del agente lo **construye y mantiene Claude Code**; la **explotación** diaria (auditar cuentas, guiar correcciones, generar informes) la hará el usuario con **GitHub Copilot (Sonnet)**. Que dos herramientas modifiquen el "cerebro" del agente sin gobernanza rompería la trazabilidad y fiabilidad (ADR-009/011). Hay que fijar una frontera clara.

**Decisión — frontera de modificación:**

- **Claude Code = gestor del agente.** Único autorizado a modificar el "agente": `tools/`, `knowledge/`, `CLAUDE.md`, `DECISIONS.md` (ADRs), `.claude/`, `.github/`, `.vscode/`, `templates/`, `tests/`, `.mcp.json`. Todo cambio pasa por gobernanza (ADR si es relevante + commit).
- **Copilot = explotador.** **Lee** todo y **usa** el agente siguiendo los playbooks. **NO modifica** el agente. Sus únicas escrituras permitidas:
  1. **Outputs** en `reports/output/` (informes y `REGISTRO-CAMBIOS.md`).
  2. **Append** a `AGENT_CHANGE_REQUESTS.md` (bandeja de peticiones): si detecta un bug en el tool, un hueco de conocimiento, una regla a cambiar, etc., **lo anota ahí** en vez de editarlo; Claude Code lo procesa.
- Copilot **guía** al usuario a cambiar datos en CoinTracking (eso es acción del usuario, no del agente), y lo registra según ADR-011.

## Decision

[Decision not found]

## Consequences

- ✅ Un único responsable del agente → cambios gobernados y trazables
- ✅ Copilot aporta mejoras sin romper nada: las canaliza como peticiones
- ✅ Separación limpia: "cerebro/reglas" (Claude) vs "operación/outputs" (Copilot)
- ⚠️ Requiere que Copilot respete la frontera (reforzado en `.github/copilot-instructions.md`); no es un candado técnico, es una norma
