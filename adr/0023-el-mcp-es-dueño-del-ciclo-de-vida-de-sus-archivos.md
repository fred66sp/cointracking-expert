# ADR-023: El MCP es dueño del ciclo de vida de sus archivos de caché (`cointracking_delete_project`)

**Status:** Accepted

**Date:** 2026-07-03

## Context

Al limpiar proyectos de prueba (`binance2025`, `pruebas`) se borraron a mano las carpetas de `.cache/cointracking/<proyecto>` con `rm -rf` desde Claude Code. El archivo `cointracking-mcp.db` de `pruebas` quedó bloqueado ("Device or resource busy") incluso tras cerrar el proyecto (`cointracking_close_project`) y cambiar a otro (`cointracking_switch_project`), porque `close_project` solo vacía la caché en memoria y hace flush — nunca llama a `Store.Close()` (que sí cierra el `*sql.DB`). El usuario señaló, acertadamente, que el problema de fondo es de diseño: un proceso externo (Claude Code) no debería gestionar el ciclo de vida de archivos que pertenecen a otro proceso (el servidor MCP); eso invita justo a este tipo de carrera con el bloqueo de archivos de Windows.

**Decisión:**

El servidor MCP pasa a ser responsable de crear **y destruir** sus propios archivos de caché:

1. Nueva tool `cointracking_delete_project(project_name)` (`internal/tools/delete_project.go`) que borra permanentemente `cache-dir/<project>`. Rechaza borrar el proyecto actualmente activo (evita cerrarse el archivo bajo los pies); para otros proyectos, reintenta `os.RemoveAll` con backoff corto (bloqueos transitorios de antivirus/indexador en Windows tras cerrar un handle).
2. Claude Code y el usuario **ya no deben hacer `rm -rf` sobre `.cache/cointracking/<proyecto>`** a mano; deben usar esta tool (o, si el MCP no está conectado, aceptar que puede fallar por bloqueo y no forzar el borrado).
3. Pendiente, no bloqueante para esta decisión: revisar por qué `Store.Close()` no libera el handle de forma fiable en Windows (posible causa: comportamiento del VFS de `modernc.org/sqlite` o interacción con antivirus/indexador) — el reintento con backoff en la nueva tool es un paliativo, no una corrección de la causa raíz.

## Decision

[Decision not found]

## Consequences

- ✅ Alinea con ADR-012 (división de responsabilidades): el servidor gestiona sus propios recursos; Claude Code ya no necesita saber cómo está implementada la persistencia de caché para poder limpiarla.
- ✅ Evita el fallo "file busy" que bloqueaba el borrado de proyectos descartados.
- ⚠️ `USER_INPUT/<proyecto>` y `reports/output/<proyecto>` siguen siendo responsabilidad de Claude Code (no son del MCP), así que borrar un proyecto por completo sigue siendo una operación en dos partes: `rm -rf` de esas dos carpetas + `cointracking_delete_project` para la caché.
