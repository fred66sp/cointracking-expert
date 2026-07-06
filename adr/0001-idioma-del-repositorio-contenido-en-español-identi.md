---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-001: Idioma del repositorio (contenido en español, identificadores en inglés)

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-02

## Context

El proyecto tiene como objetivo servir principalmente a usuarios hispanohablantes con enfoque en cumplimiento fiscal español. El equipo de desarrollo también es hispanohablante. Al mismo tiempo, el código Python debe seguir convenciones universales (PEP 8) para mantenerse legible, buscable e interoperable con el ecosistema.

## Decision

*(Restaurada 2026-07-05 desde `DECISIONS.md` §ADR-001 — la migración automática a MADR, ADR-025, dejó esta sección vacía.)*

Se adopta el modelo **híbrido**:

- **En español (contenido para humanos):**
  - Contenido de toda la documentación (`.md`)
  - Docstrings
  - Comentarios de código
  - Mensajes de error y de log dirigidos al usuario
- **En inglés (identificadores técnicos):**
  - Nombres de archivos y carpetas (`README.md`, `src/`, `engines/`)
  - Nombres de clases, funciones, métodos y variables (PEP 8)

## Consequences

- ✅ Documentación accesible para usuarios y equipo hispanohablante
- ✅ Código que respeta PEP 8 y es interoperable con el ecosistema Python
- ✅ Nombres de archivo estables y buscables (identificadores técnicos universales)
- ⚠️ Requiere disciplina para mantener la separación (contenido vs identificador)
- ⚠️ Menor comunidad potencial de contribuidores globales por la documentación en español

## Notes

**Notas adicionales:**

Esta decisión es **permanente** para el proyecto y prevalece sobre cualquier documento que sugiera "todo en español".
