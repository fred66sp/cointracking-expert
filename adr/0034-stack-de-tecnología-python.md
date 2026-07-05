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

[Decision not found]

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
