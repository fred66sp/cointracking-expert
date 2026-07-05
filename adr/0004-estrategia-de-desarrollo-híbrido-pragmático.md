---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-004: Estrategia de desarrollo (híbrido pragmático)

**Status:** Accepted

**Date:** 2026-07-02

## Context

Existía una contradicción entre dos documentos del repositorio:

- `ROADMAP.md` define un enfoque **documentación primero**: completar especificaciones y base de conocimiento antes de escribir código (sin implementación hasta la Fase 4).
- `ARCHITECTURE_REVIEW.md` recomienda lo contrario: **implementar primero**, validar con datos reales de CoinTracking y refinar las especificaciones de forma iterativa, para evitar la divergencia especificación-realidad.

El riesgo central que motiva esta decisión: **nadie conoce los datos reales de CoinTracking hasta que los mira**. Una especificación de import, duplicados o transferencias escrita sobre suposiciones puede resultar incorrecta al enfrentarse a un CSV real (comisiones que descuadran cantidades, zonas horarias distintas, movimientos en el mismo segundo). Documentar mucho sobre datos no vistos genera trabajo que luego hay que descartar.

Al mismo tiempo, hay partes del dominio que **sí** están definidas por fuentes externas estables (reglas fiscales españolas, principios de arquitectura) y se benefician de especificarse por completo antes de programar.

## Decision

*(Restaurada 2026-07-05 desde `DECISIONS.md` §ADR-004 — la migración automática a MADR, ADR-025, dejó esta sección vacía.)*

Se adopta el **híbrido pragmático**:

- **Documentación primero** para el dominio estable y de fuente externa:
  - Reglas de tributación (definidas por normativa; no se "descubren" programando)
  - Principios, arquitectura y contratos entre motores
  - Metodología de auditoría
- **Validación con datos reales antes de cerrar la spec** para el dominio de datos desordenados:
  - Formato CSV de CoinTracking, importación y normalización
  - Detección de duplicados, emparejamiento de transferencias, reconstrucción de libro mayor
  - Peculiaridades por exchange
  - → Estas specs se redactan en borrador, se contrastan contra **exportaciones reales de CoinTracking** y solo entonces se dan por cerradas.
- **Especificar cada motor justo antes de implementarlo**, no los nueve por adelantado, para evitar el agotamiento de especificación.
- **Las specs son documentos vivos**: se refinan si la implementación o los datos reales revelan supuestos incorrectos.

`ARCHITECTURE_REVIEW.md` queda como una revisión asesora (una instantánea de opinión), no como estrategia vinculante.

## Consequences

- ✅ El repositorio deja de contradecirse: hay una única estrategia vinculante
- ✅ Se preserva la fortaleza del proyecto (disciplina de documentación) donde aporta valor
- ✅ Se neutraliza el riesgo de divergencia especificación-realidad en las partes sensibles a datos
- ✅ Se evita el agotamiento de especificación (specs por motor, justo a tiempo)
- ⚠️ Requiere conseguir exportaciones reales de CoinTracking pronto — es una dependencia crítica, no opcional
- ⚠️ Exige disciplina para clasificar cada pieza como "estable" vs "sensible a datos"
