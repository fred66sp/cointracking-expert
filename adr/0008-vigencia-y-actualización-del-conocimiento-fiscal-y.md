---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-008: Vigencia y actualización del conocimiento (fiscal y CoinTracking)

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-02

## Context

El conocimiento del agente tiene **dos patas** y **ambas caducan**:

- **Fiscal:** la normativa cambia cada año — tramos de la base del ahorro (el tramo alto pasó del 28 % al 30 % en 2025), umbrales de obligaciones informativas, criterios de la DGT (consultas matizadas o superadas) y plazos de campaña.
- **CoinTracking:** la plataforma evoluciona — formato del CSV export, nuevos tickers y sufijos de colisión, herramientas y parámetros del MCP/API, límites de tasa, y peculiaridades por exchange.

Un conocimiento fijado en una fecha puede quedar **obsoleto** y hacer que el agente dé cifras o supuestos incorrectos.

## Decision

**Decisión:**

1. **Metadatos de vigencia.** Todo documento de conocimiento sensible al tiempo declara en su cabecera: **Última verificación** (fecha) y **Vigencia** (ejercicios/versión a los que aplica).
2. **Comprobación de vigencia obligatoria.** Antes de apoyarse en un dato que puede haber cambiado (tramos, tipos, umbral del Modelo 721, criterios DGT; o formato CSV, tickers, herramientas MCP), el agente **compara** el contexto (ejercicio solicitado, fecha de hoy) con la "Última verificación"/"Vigencia" del documento.
3. **Ante posible desfase**, el agente **avisa al usuario** y **reverifica contra la fuente autorizada** antes de afirmar:
   - Fiscal → AEAT / BOE / DGT (búsqueda web).
   - CoinTracking → centro de ayuda oficial (URLs en `knowledge/cointracking/reference/CATALOG.md`) y **los datos reales del usuario** (el CSV/MCP son la verdad sobre el formato actual).
   Nunca presenta como vigente un dato sin confirmar que aplica.
4. **Checklist de revisión** en los índices de cada pata (`knowledge/taxation/spain/INDEX.md` y `knowledge/cointracking/INDEX.md`) con lo que cambia y con qué periodicidad.

## Consequences

- ✅ El agente no arrastra información caducada (ni fiscal ni de plataforma); se autoactualiza cuando hace falta
- ✅ Transparencia: distingue "verificado para 2025 / esta versión" de "asumido"
- ⚠️ Requiere disciplina de metadatos y, en su caso, una verificación extra (web o contra los datos reales)
- ⚠️ La reverificación depende de tener acceso a la fuente en la sesión
