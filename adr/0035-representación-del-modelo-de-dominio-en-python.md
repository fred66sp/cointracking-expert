---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-003: Representación del modelo de dominio en Python

**Status:** Accepted

**Date:** 2026-07-02

## Context

El `DOMAIN_MODEL.md` originalmente usaba pseudocódigo Kotlin. Ya fue traducido a pseudocódigo Python. Queda decidir la tecnología concreta con la que se materializará el modelo cuando comience la implementación (Fase 4).

## Decision

**Decisión:**

El modelo de dominio se implementará con **Pydantic v2**, en coherencia con ADR-002:

- Entidades y objetos de valor como `BaseModel`
- Inmutabilidad mediante `model_config = ConfigDict(frozen=True)`
- Validación de invariantes con validadores de Pydantic (`@field_validator`, `@model_validator`)
- Cantidades, precios y comisiones tipados como `Decimal`
- Identificadores como tipos dedicados (p. ej. `TransactionId`) para seguridad de tipos
- Nomenclatura de atributos en `snake_case` (PEP 8), según ADR-001

## Consequences

- ✅ Coherencia total con el stack de ADR-002
- ✅ Las invariantes del modelo de dominio quedan enforced en tiempo de construcción
- ⚠️ El pseudocódigo Python actual de `DOMAIN_MODEL.md` es orientativo; al implementar puede requerir ajustes menores hacia la sintaxis real de Pydantic v2

**Próximos pasos:**

- Al llegar a la Fase 4, materializar los objetos de valor primero (`Quantity`, `Money`, `Timestamp`)
- Validar el modelo contra exportaciones reales de CoinTracking
