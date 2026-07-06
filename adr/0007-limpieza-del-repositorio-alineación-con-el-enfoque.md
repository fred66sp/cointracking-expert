---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-007: Limpieza del repositorio (alineación con el enfoque agente)

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-02

## Context

Tras ADR-006 (el producto es un agente de IA en Claude Code, no un SDK de motores deterministas), el repositorio seguía conteniendo el andamiaje de la visión anterior: paquetes Python vacíos, especificaciones de motores, dependencias y CI de Python, y documentos de la visión de framework. Ese material ya no describe lo que se construye y genera ruido.

## Decision

**Decisión:**

Se eliminan los artefactos que solo servían al SDK descartado:

- `src/` (paquetes Python vacíos), `requirements.txt`, `requirements-dev.txt`
- `.github/workflows/ci.yml` (CI de pytest/flake8/mypy) y `.github/ISSUE_TEMPLATE/`
- `engines/` (9 specs de motores deterministas → sustituidos por el playbook del agente en `.claude/skills/`)
- `ARCHITECTURE.md`, `ARCHITECTURE_REVIEW.md`, `DOMAIN_MODEL.md`, `ROADMAP.md`, `PROJECT_CHARTER.md`
- `CONTRIBUTING.md`, `docs/DEVELOPMENT_GUIDE.md`, `docs/INDEX.md`, `docs/PROJECT_MANIFESTO.md`
- Carpetas vacías de scaffolding: `cases/`, `examples/`, `prompts/`, `schemas/`, `scripts/`, `tests/`
- `COPILOT.md` → sustituido por `CLAUDE.md` (lo carga Claude Code)

Se conservan y adaptan: `.claude/` (agente + skill), `.mcp.json`, `knowledge/`, `DECISIONS.md`, `FOUNDATION.md`, `templates/`, `docs/GLOSSARY.md`, `LICENSE`, `CHANGELOG.md`, y `README.md` (reescrito para el agente).

## Consequences

- ✅ El repositorio refleja lo que es: un agente + su base de conocimiento
- ✅ Menos ruido; navegación y mantenimiento más simples
- ✅ Todo lo eliminado permanece en el historial de git si se necesita recuperar
- ⚠️ **ADRs anteriores (002, 003, 004, 006) referencian documentos ya eliminados** (`ROADMAP.md`, `ARCHITECTURE_REVIEW.md`, `DOMAIN_MODEL.md`, `PROJECT_CHARTER.md`). Se conservan sin reescribir: son **registro histórico** de las decisiones tal como se tomaron. Este ADR-007 es el contexto que explica esas referencias.
