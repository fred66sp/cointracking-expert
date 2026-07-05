---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-002: Stack de tecnología Python

**Status:** Accepted

**Date:** 2026-07-02

## Context

Se requiere seleccionar un stack de tecnología Python para la implementación. Necesitamos decisiones sobre:
- Versión mínima de Python
- Librería de validación (pydantic, dataclasses, attrs)
- Tipo numérico para cantidades y precios (float vs Decimal)
- Base de datos (SQLite, PostgreSQL, en memoria)
- Framework web (FastAPI, Flask, Django) para API futura

## Decision

*(Restaurada 2026-07-05 desde `DECISIONS.md` §ADR-002 — la migración automática a MADR, ADR-025, dejó esta sección vacía. Nota: este archivo es la decisión histórica "ADR-002: Stack de tecnología Python" del framework/SDK descartado por ADR-006/007; se conserva por trazabilidad y no debe confundirse con el ADR-002 vigente, "Fuente de verdad".)*

Se adopta **Pydantic v2 + Decimal + SQLite**:

- **Validación y modelos:** Pydantic v2 (`BaseModel` con `model_config = ConfigDict(frozen=True)` para inmutabilidad)
- **Tipo numérico:** `decimal.Decimal` para todas las cantidades, precios y comisiones — **nunca `float`** (garantiza determinismo y reproducibilidad, mitiga el riesgo de aritmética de punto flotante identificado en la revisión de arquitectura)
- **Persistencia:** SQLite para el MVP, con capa de repositorio que permita migrar a PostgreSQL sin cambiar la lógica de dominio
- **Versión de Python:** 3.11+ (coincide con la matriz de CI; se puede ampliar el rango si es necesario)
- **Framework web:** aplazado hasta la Fase 6 (API REST); candidato preferente FastAPI por su integración nativa con Pydantic

## Consequences

- ✅ Validación y serialización automáticas y robustas
- ✅ Determinismo garantizado por `Decimal`
- ✅ Migración de persistencia sin tocar el dominio (patrón repositorio)
- ✅ Continuidad natural hacia FastAPI en la fase de API
- ⚠️ Pydantic v2 añade una dependencia externa y una curva de aprendizaje
- ⚠️ `Decimal` es más lento que `float`; aceptable frente al requisito de exactitud

## Notes

**Notas adicionales:**

Esta decisión desbloquea ADR-003 (traducción del modelo de dominio) y la creación de `requirements.txt` / `pyproject.toml`.
