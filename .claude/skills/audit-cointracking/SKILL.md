---
name: audit-cointracking
description: Audita una cuenta o exportación de CoinTracking. Reconcilia los datos y detecta transferencias huérfanas, ventas sin base de coste, duplicados, saldos imposibles e incoherencias fiscales españolas, explicando cada hallazgo con evidencia. Usa el MCP de la API de CoinTracking si está disponible y/o el CSV export.
---

# Auditoría de CoinTracking

Ejecuta una auditoría de reconciliación sobre los datos de CoinTracking del usuario, siguiendo el playbook de abajo. Trabaja en español y aplica la base de conocimiento del repo (ver el subagente `cointracking-auditor` y `knowledge/`).

## Paso -1 — Proyecto activo obligatorio (ADR-013)

**Antes de cualquier otra cosa**, si esta conversación todavía no tiene un proyecto activo fijado: lista las subcarpetas de `USER_INPUT/` (los proyectos existentes) y pregunta al usuario con cuál quiere trabajar, o si quiere crear uno nuevo (pide un nombre). Ver `CLAUDE.md` §"Proyecto activo obligatorio". Una vez fijado, reutilízalo el resto de la conversación.

En cuanto quede fijado (o si ya lo estaba de una skill anterior en la misma conversación pero el MCP no se ha sincronizado aún), y si hay herramientas `cointracking_*` disponibles: llama a `cointracking_switch_project(project_name=<proyecto activo>)` antes de cualquier otra tool `cointracking_*`. Alinea en caliente el proyecto del MCP con el proyecto de datos activo, sin reiniciar el servidor (ADR-016). Si la respuesta trae `already_active: true`, no hace falta avisar de nada; si cambia de proyecto, ya no hace falta la advertencia de aislamiento que existía antes de ADR-016.

## Paso 0 — Diálogo de arranque y preparación

**Conversa antes de ejecutar, en lenguaje llano** (el usuario no domina CoinTracking; evita "API/MCP/cotejo"). Anuncia qué vas a hacer y **ofrece una comprobación extra con un CSV**, esperando respuesta. Por ejemplo:
> "Voy a revisar tu cuenta leyendo tus datos directamente de CoinTracking. Como comprobación adicional opcional, puedo compararlos con un archivo que descargues tú mismo; así detecto si algo no cuadra entre ambos. ¿Quieres hacer esa comprobación extra (te guío para descargar el archivo) o sigo solo con la conexión automática?"
- Si acepta y no sabe cómo, **guíalo paso a paso** para exportar la lista de operaciones a CSV. **Consulta primero `knowledge/cointracking/WEB_APP_GUIDE.md` §"La página de Transacciones"** — ya tiene la ruta verificada sobre la interfaz real (botón **Export** → **CSV (Exportación Completa)**, la que usa el agente como Trade Table). Solo si esa guía no cubre el caso, acude a `knowledge/cointracking/reference/CATALOG.md` (artículo oficial de exportación/backup) y **verifícalo en la sesión** antes de instruir — no inventes rutas de menú ni prefieras una fuente externa genérica a una interna ya verificada. Pídele que **guarde el archivo en la carpeta `USER_INPUT/<proyecto>/`** (el proyecto activo fijado en el Paso -1).

Después:

1. **Carga el conocimiento**: lee `knowledge/cointracking/CSV_FORMAT.md`, `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md`, `knowledge/taxation/spain/CAPITAL_GAINS.md`. Son las reglas que aplicarás y citarás.
2. **Localiza los datos** (ADR-006, ambas vías):
   - **MCP de CoinTracking**: si hay herramientas MCP de CoinTracking disponibles en la sesión, úsalas para datos en vivo. Si el usuario menciona un MCP pero no está conectado, pídele que lo conecte (`/mcp`).
   - **CSV export**: búscalo en la carpeta **`USER_INPUT/<proyecto>/`** (ahí deja el usuario los archivos que le pedimos; ver `USER_INPUT/README.md`). Está en `.gitignore` por privacidad.
   - Si hay ambos, usa el CSV como validación cruzada del MCP.
   - Si no hay ninguno, dilo y detente: no inventes datos.
   - **Documentación adicional bajo demanda:** si un hallazgo concreto lo justifica (duplicados dudosos, transferencia huérfana, discrepancia de balance…), consulta `knowledge/cointracking/DOCUMENT_CHECKLIST.md` para saber qué informe extra de CoinTracking pedir (§A) o qué dato del propio exchange solicitar al usuario (§B, p. ej. trade_id para ADR-014). No lo pidas todo de entrada; solo lo que el hallazgo requiera.
3. **Normaliza mentalmente** según el conocimiento: fechas a UTC desde `Europe/Madrid` con DST (ADR-005), importes con precisión decimal, ticker completo con sufijo (`SOL` ≠ `SOL2`), parsear por posición (3 columnas `Cur.`).
4. **Sé económico (ADR-010).** Cachea lo obtenido en `.cache/cointracking/` (con marca de tiempo) y reutilízalo; no vuelvas a llamar si ya lo tienes fresco. Consulta lo mínimo (rango, `limit`, agregados). Para volúmenes grandes (historial completo), vuelca a fichero y **procésalo con un script**, subiendo al contexto solo el resultado. No pegues JSON crudo.

## Paso 1 — Ejecuta el playbook de chequeos

**Usa `tools/ct_audit.py`** (código vetado y determinista; ADR-006/009) para los chequeos mecánicos sobre el CSV — saldos, saldos negativos, transferencias huérfanas, duplicados, colisiones — en vez de re-derivar la lógica en el momento (evita errores como doble-contar comisiones). Ejecútalo y sube al contexto **solo el resultado** (ADR-010). Ejemplo: `python tools/ct_audit.py "<csv>" --exchange Coinbase --check all`. Luego **interpreta y explica** cada hallazgo con el conocimiento:

1. **Completitud de importación** — ¿faltan cuentas/exchanges (seguimiento de un solo lado)? Provoca base de coste 0 y ganancias infladas. *(COST_BASIS §3.2, "READ FIRST")*
2. **Transferencias huérfanas** — retiradas sin depósito emparejado y viceversa. Empareja por **Tx Hash** (nivel 1, fuerte) y, si falta, por **moneda + importe ≈ retirada − comisión + ventana temporal + cuentas distintas** (nivel 2, heurístico con confianza). *(CSV_FORMAT §7)*
3. **Ventas sin base de coste** — activos vendidos/permutados sin compra previa registrada. Distingue "compra importada como depósito" de "compra realmente ausente". Nunca asumas base 0 en silencio. *(COST_BASIS §3)*
4. **Duplicados** — filas idénticas: separa duplicado por reimportación (error) de repetición legítima (comisiones/recompensas recurrentes). *(CSV_FORMAT §9, COST_BASIS §4.2)*
5. **Saldos negativos imposibles** — vender/enviar más de lo disponible de un activo. Distingue del **FIAT negativo**, que suele ser artefacto de no importar depósitos fiat, no una imposibilidad. *(COST_BASIS §4.1)*
6. **Orden temporal de transferencias** — un depósito anterior a su retirada rompe la transferencia de base de coste. Verifica `retirada ≤ depósito` tras normalizar a UTC. *(COST_BASIS §2, ADR-005)*
7. **Coherencia fiscal (cualitativa)** — ¿se tratan las permutas cripto-cripto como hecho imponible? ¿se usaría el método **FIFO** (España) y no el pool promediado? Señala riesgos; **no** des la cifra final vinculante. *(CAPITAL_GAINS §3-4, COST_BASIS §1)*
8. **Colisión de tickers** — no fusionar activos con sufijo (`SOL2`, `WLD3`, `ID2`). *(CSV_FORMAT §8)*

Para análisis profundo puedes delegar en el subagente `cointracking-auditor`.

## Paso 2 — Informe (persistente, ADR-011)

Usa la plantilla `templates/AUDIT_REPORT.md`. Ordena los hallazgos por severidad y, para cada uno: **causa, evidencia, impacto, recomendación**. Cierra con un resumen y con lo **no verificado**.

**Guárdalo** en `reports/output/<proyecto>/AAAA-MM-DD_auditoria_<cuenta>.md` (no lo dejes solo en el chat). Actualiza la memoria del proyecto (`audit_state`) con lo hecho/pendiente, indicando a qué proyecto corresponde.

## Paso 3 — Remediación guiada (ofrécela)

Tras el informe, **ofrece ayudar a corregir** los hallazgos en la web de CoinTracking. Para cada problema con solución, guía al usuario **paso a paso y en lenguaje llano** usando `knowledge/cointracking/WEB_APP_GUIDE.md` (mapa de remediación: hallazgo → acción → artículo oficial).

- **Antes de dar los pasos clic a clic, abre y lee el artículo oficial citado** para confirmar que siguen vigentes (la interfaz cambia; ADR-008/009). Si no puedes verificarlos, dilo y no improvises la ruta.
- Explica y guía los problemas **de uno en uno en la conversación** (para no abrumar al usuario), pero **no verifiques cada uno contra el MCP según se resuelve**. Ve marcando en una lista interna qué le has pedido corregir, y pide al usuario que confirme por chat cuando lo haya aplicado en la web; sigue con el siguiente sin llamar a ninguna herramienta `cointracking_*` de por medio.
- **Verifica todo el lote al final, no una corrección a la vez (ADR-010 punto 7; minimiza llamadas MCP).** Cuando el usuario confirme que ha terminado la ronda completa de correcciones: llama **una sola vez** a `cointracking_invalidate_cache` y luego re-obtén los datos con el mínimo de consultas (agregados antes que `get_trades` completo) para comprobar de golpe todos los hallazgos de esta ronda. Si alguno persiste, señálalo específicamente. Si el usuario prefiere verificar uno a uno, respétalo, pero antes explícale que consume más cuota de API (límite 60/hora).
- Recuerda al usuario que, tras esta verificación por lote, si queda algún hallazgo pendiente puede corregirse en una nueva ronda con el mismo patrón (corregir todo → verificar todo).
- **Registra cada cambio aplicado (ADR-011)** en `reports/output/<proyecto>/REGISTRO-CAMBIOS.md` (append-only): qué se cambió, por qué, evidencia, estado antes→después y verificación en vivo. Actualiza también la memoria (`audit_state`).

## Límite de determinismo (ADR-006)

Recuerda: encuentras y explicas; **no** produces cifras fiscales vinculantes. Toda cantidad exacta es «estimación no vinculante» salvo que provenga de un cálculo determinista.
