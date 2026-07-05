---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-013: Estructura multi-proyecto obligatoria (datos de usuario y estado; MCP pospuesto)

**Status:** Accepted

**Nota (2026-07-05):** el título dice "MCP pospuesto" porque en la fase 1 (2026-07-02) el aislamiento del MCP por proyecto quedaba fuera de alcance. Esa cuestión **se resolvió el 2026-07-03 por ADR-016** (`cointracking_switch_project`) — ver sección "Cuestión abierta" más abajo, ya marcada como resuelta en el propio cuerpo del documento. El nombre del archivo no se cambia para no romper referencias existentes; el estado real es: ambas fases completas.

**Date:** 2026-07-02 (propuesta inicial) — **revisado y redecidido 2026-07-03** tras corrección del usuario sobre el alcance real.

## Context

La propuesta original (v1, 2026-07-02) planteaba esto como una mejora "nice to have" a implementar cuando hubiera un segundo caso. El usuario corrigió el enfoque el 2026-07-03: **no es opcional ni futuro** — todo trabajo del agente sobre CoinTracking (auditar, declarar, lo que sea) debe ocurrir **siempre** dentro de un **proyecto activo**, porque eso es lo que aísla qué CSV y qué datos se usan. Sin esto, el agente puede mezclar sin querer datos de casos distintos.

## Decision

**Decisión (fase 1 — sin tocar el MCP):**

1. **Estructura por proyecto**, dentro del repo (gitignored, salvo los `README.md`):
   - `USER_INPUT/<nombre_proyecto>/` — CSV y otras fuentes del caso (sustituye el uso plano anterior de `USER_INPUT/`).
   - `reports/output/<nombre_proyecto>/` — informes y `REGISTRO-CAMBIOS.md` del caso (sustituye el uso plano anterior de `reports/output/`).
   - Estado del proyecto: por ahora se sigue usando la memoria global (`audit_state` en `~/.claude`), pero **prefijada por proyecto**; migrar a un `estado.md` por proyecto queda como mejora futura si hace falta que Copilot lo lea directamente sin memoria (ADR-011/012).
2. **Puerta de entrada obligatoria (lo nuevo de la corrección):** en cualquier conversación, en cuanto el usuario pida algo relacionado con CoinTracking y **todavía no haya un proyecto activo fijado en esa conversación**, el agente debe, antes de ejecutar nada más:
   1. Listar los proyectos existentes (subcarpetas de `USER_INPUT/`).
   2. Si hay uno o más, **preguntar** con cuál trabajar, o si se quiere crear uno nuevo.
   3. Si no hay ninguno, ofrecer crear el primero (pidiendo un nombre).
   - Una vez fijado el proyecto activo en la conversación, se reutiliza para el resto de la sesión (no se vuelve a preguntar salvo que el usuario pida cambiar de proyecto).
3. **Migración del caso existente:** los datos que vivían en `USER_INPUT/` y `reports/output/` planos se migraron a `USER_INPUT/agp/` y `reports/output/agp/` el 2026-07-03 (nombre elegido por coincidir con el `--project agp` ya usado en `.mcp.json`).

**Cuestión abierta — MCP no aislado por proyecto todavía:**

El servidor MCP (`cointracking-mcp/`) solo admite el proyecto como **flag de arranque del proceso** (`--project`, ver `cointracking-mcp/SPEC/06-configuration.md`); ninguna tool acepta hoy un parámetro `project_name` en tiempo de ejecución. Cambiar de proyecto en el MCP implicaría reiniciar el servidor (y por tanto Claude Code). El usuario decidió explícitamente **posponer esto**: por ahora el proyecto de datos de usuario (CSV) y el `--project` del MCP **no están enlazados** — el MCP sigue usando el valor fijo de `.mcp.json` (`agp`) con independencia del proyecto de datos activo en la conversación. Si en el futuro se trabaja con un segundo proyecto real, resolver esto (opción más probable: añadir `project_name` como parámetro a las tools existentes del servidor Go) antes de confiar en el MCP para ese segundo proyecto.

**✅ Resuelta 2026-07-03 — ver ADR-016.** El MCP ya expone `cointracking_switch_project` para cambiar de proyecto activo en caliente, sin reiniciar el servidor.

## Consequences

- ✅ Aislamiento real entre casos en la capa de datos de usuario (CSV) e informes/estado
- ✅ Flujo predecible: nunca se opera "a ciegas" sin saber en qué proyecto se está
- ✅ El MCP ya está aislado por proyecto en caliente desde ADR-016 (antes de ADR-016, sus datos en vivo seguían siendo los de `--project agp` fijo con independencia del proyecto de datos activo)
- ⚠️ Requirió migrar rutas existentes (`USER_INPUT/`, `reports/output/`) y actualizar ambas skills y `CLAUDE.md` para aplicar la puerta de entrada
