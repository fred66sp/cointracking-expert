---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-016: Cambio de proyecto activo en caliente en el MCP (`cointracking_switch_project`)

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-03

## Context

ADR-013 dejó abierta la cuestión de que el MCP (`cointracking-mcp/`) solo admite el proyecto como flag de arranque del proceso (`--project` en `.mcp.json`) y ninguna tool aceptaba `project_name` en tiempo de ejecución — cambiar de proyecto exigía reiniciar el servidor. El usuario propuso primero escribir un `.env` que el agente reescribiera al cambiar de proyecto; al evaluarlo, se identificaron dos problemas: (1) el binario Go no lee ningún fichero hoy, habría que implementar parseo de `.env`, y (2) el servidor MCP se arranca una vez por sesión, así que reescribir un fichero de config no evita el reinicio — solo lo hace más seguro de tocar que `.mcp.json`. Se optó por una alternativa que resuelve el problema de raíz: un tool MCP nuevo.

## Decision

**Decisión:**

Añadir el tool `cointracking_switch_project(project_name)` a `cointracking-mcp/internal/tools/switch_project.go`, hermano de `cointracking_close_project` ya existente:

1. Valida `project_name` con la misma regla que `--project` al arrancar (`config.ValidateProjectName`, alfanumérico + `_` + `-`).
2. Si coincide con el proyecto ya activo, es un no-op (`already_active: true`) que no toca la caché.
3. Si no, hace flush + close de la caché del proyecto saliente (igual que `close_project`) y abre/crea la caché SQLite del proyecto entrante bajo `{cache-dir}/{project}` (misma lógica que `NewApp` al arrancar, extraída a `openProjectCache` para no duplicarla).
4. Credenciales, `--tier` y el limitador de tasa son del proceso (una cuenta de CoinTracking), no del proyecto: no cambian.

`App` (en `app.go`) pasa a guardar `cfg`/`cache`/`store` bajo un `sync.RWMutex`, con accesores (`Project()`, `CacheManager()`, `Store()`, `CacheDir()`) que los demás tools usan en vez de leer los campos directamente — así una llamada a `switch_project` no puede dejar a otra tool leyendo un puntero a medio reemplazar.

## Consequences

- ✅ Resuelve la cuestión abierta de ADR-013: el proyecto de datos activo en la conversación y el proyecto del MCP pueden mantenerse enlazados sin reiniciar Claude Code ni editar `.mcp.json`.
- ✅ Más simple que la alternativa del `.env`: no requiere que el binario Go parsee ficheros de config ni que el agente escriba en disco antes de cada cambio de proyecto — un tool call basta.
- ✅ Verificado con test de integración (`TestSwitchProject`): nombre inválido rechazado sin tocar estado, no-op al re-seleccionar el mismo proyecto, aislamiento entre proyectos, y recarga desde disco (sin llamadas nuevas a la API) al volver a un proyecto ya visitado en el proceso.
- ✅ Ambas skills (`audit-cointracking`, `spanish-tax-return`) y `CLAUDE.md` §"Proyecto activo obligatorio" actualizados para llamar a `cointracking_switch_project` en la puerta de entrada (Paso -1) en cuanto se fija el proyecto activo.
