# USER_INPUT — deja aquí tus archivos, por proyecto

**Desde ADR-013 (2026-07-03), todo trabajo con el agente ocurre dentro de un proyecto.** Cada proyecto tiene su propia subcarpeta aquí: `USER_INPUT/<nombre_proyecto>/`. Los datos de un proyecto **nunca se mezclan** con los de otro.

## Cómo usarla

1. Al empezar a pedirle algo sobre CoinTracking (auditar, declarar…), el agente te preguntará **con qué proyecto quieres trabajar** (o si quieres crear uno nuevo). Esto ocurre siempre al arrancar, antes de hacer nada más.
2. Cuando el agente te pida un archivo, guárdalo en la subcarpeta de **ese proyecto**: `USER_INPUT/<nombre_proyecto>/`. Normalmente será el **CSV export de CoinTracking** ("Trade Table"), pero también cualquier otra fuente que se te solicite (ver `knowledge/cointracking/DOCUMENT_CHECKLIST.md`).
3. Dile al agente que ya está; él lo buscará ahí.

Si no sabes cómo exportar el CSV de CoinTracking, **pídeselo al agente** y te guiará paso a paso.

## Proyectos existentes

- **`agp/`** — caso migrado a esta estructura el 2026-07-03 (antes vivía en `USER_INPUT/` sin subcarpeta). Coincide con el nombre `agp` que ya usa `.mcp.json --project agp` para el MCP, aunque por ahora (fase 1 de ADR-013) el proyecto de datos de usuario y el `--project` del MCP **no están enlazados todavía** — ver limitación conocida en `DECISIONS.md#ADR-013`.

## 🔒 Privacidad (importante)

- Todo lo que pongas aquí son **tus datos reales** (historial financiero). **Nunca se sube al repositorio** — esta carpeta está excluida en `.gitignore`; solo se versiona este `README.md`.
- No compartas estos archivos con nadie salvo tu asesor fiscal.

## Formatos habituales

- `CoinTracking · Trade Table.csv` — exportación de operaciones de CoinTracking.
- Otros que el agente te indique según el caso (ver `knowledge/cointracking/DOCUMENT_CHECKLIST.md`).
