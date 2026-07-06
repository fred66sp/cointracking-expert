---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-041: Procesos MCP huérfanos en Windows — diagnóstico y protocolo de limpieza

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-05

## Context

**Caso real (2026-07-05):** al intentar borrar la caché del proyecto `demo` con `cointracking_delete_project`, el borrado falló con *"The process cannot access the file because it is being used by another process"* — pese a que el servidor de la sesión activa no tenía `demo` abierto. La investigación (`Get-Process cointracking-mcp`) reveló **7 procesos del servidor MCP corriendo simultáneamente**, arrancados entre el 03-07 y el 04-07 por sesiones de Claude Code ya cerradas. Solo uno pertenecía a la sesión viva; los otros 6 eran **huérfanos**, y uno de ellos (PID 8388) mantenía un handle SQLite abierto sobre la caché de `demo`.

**Causa raíz:** en Windows, al cerrar una ventana/sesión de Claude Code, el proceso hijo del servidor MCP (transporte stdio) **no siempre muere con su padre**. Si el padre termina sin cerrar limpiamente los handles del pipe, el hijo nunca ve EOF en stdin y `server.Run` sigue bloqueado indefinidamente. Cada huérfano consume ~18 MB y, lo importante para este proyecto, **retiene abierto el fichero SQLite del último proyecto que tuvo activo** — bloqueando borrados y, en teoría, compitiendo por escrituras.

**Por qué merece ADR:** el síntoma ("file busy" al borrar una caché) es opaco y reaparecerá; sin protocolo escrito, cada vez habría que re-derivar el diagnóstico. Además, matar procesos a ciegas puede romper la sesión viva.

## Decision

### Diagnóstico (cuándo sospechar de huérfanos)

- `cointracking_delete_project` falla con "file busy"/"being used by another process" sobre un proyecto que la sesión actual **no** tiene activo.
- O directamente: `Get-Process cointracking-mcp | Select-Object Id, StartTime` devuelve **más de un proceso** con una sola ventana de Claude Code abierta.

### Protocolo de limpieza (verificado en el caso real)

1. **Confirmar con el usuario cuántas ventanas/sesiones de Claude Code tiene abiertas.** Cada ventana viva tiene su servidor legítimo — matar el de una sesión abierta rompe su conexión MCP. Con N ventanas abiertas y M procesos, hay M−N huérfanos.
2. **No intentar identificar el proceso propio por PID o fecha de arranque** (no es fiable). En su lugar, **matar de uno en uno, del más viejo al más nuevo**, y tras cada baja **verificar que el MCP de la sesión sigue vivo** con una llamada inocua (`cointracking_cache_stats`). Si la verificación falla, se mató el propio — recuperable reconectando con `/mcp`, pero el orden viejo→nuevo lo hace improbable (el propio suele ser de los más recientes).
3. Tras cada baja, **reintentar el borrado bloqueado**; cuando el holder muere, el borrado pasa. Continuar hasta dejar exactamente un proceso por ventana abierta.
4. Para el borrado de cachés usar siempre `cointracking_delete_project` (ADR-023); el `Remove-Item` manual solo como reintento inmediato dentro de este protocolo, cuando se sabe que el único candidato a holder acaba de morir.

### Prevención (evaluada, no implementada)

Un watchdog en el servidor Go (goroutine que vigile la desaparición del proceso padre vía polling de PPID, o EOF de stdin) haría que los huérfanos se auto-terminaran. **Se pospone**: el mecanismo fiable en Windows tiene matices (el PPID puede reciclarse; el EOF no llega si el padre muere mal), el coste del síntoma es bajo una vez documentado este protocolo, y añadir lógica de auto-terminación mal calibrada podría matar el servidor en situaciones legítimas — peor remedio que enfermedad. Si los huérfanos se vuelven frecuentes o causan corrupción real, reabrir esta decisión.

## Consequences

- ✅ El síntoma "file busy" al borrar cachés tiene diagnóstico y receta escritos; no hay que re-derivarlos.
- ✅ El protocolo protege la sesión viva (verificación de vida tras cada baja) en vez de un `taskkill` masivo a ciegas.
- ✅ Caso real resuelto con él: 6 huérfanos eliminados, caché de `demo` liberada y borrada, servidor propio intacto (verificado).
- ⚠️ La prevención automática queda pendiente a propósito; hasta entonces la limpieza es manual y reactiva.
- ⚠️ Si el usuario tiene varias ventanas abiertas, el protocolo exige su confirmación previa — no es automatizable sin riesgo.

## Notes

- Relación: ADR-023 (el MCP es dueño de sus ficheros de caché — este ADR cubre el caso en que *otro proceso MCP* es el dueño accidental), ADR-016 (el proyecto activo del servidor es el que retiene el handle).
- El guard "no borrar el proyecto activo" de `cointracking_delete_project` funcionó correctamente durante el caso real (obligó a `switch_project` antes de borrar `agp`) — el bloqueo de `demo` era de un proceso ajeno, escenario distinto que ese guard no puede cubrir.
