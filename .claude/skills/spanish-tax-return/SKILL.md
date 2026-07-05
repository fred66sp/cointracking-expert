---
name: spanish-tax-return
description: Prepara la declaración de la renta española (IRPF) de un ejercicio con criptomonedas de CoinTracking. Úsalo cuando el usuario quiera "hacer la declaración de la renta", "impuestos cripto", "IRPF criptomonedas", "renta 2024/2025…", declarar ganancias patrimoniales, o el Modelo 721. Reconcilia los datos primero y luego prepara y explica lo relevante para la declaración; no produce cifras fiscales vinculantes.
---

# Declaración de la renta española (IRPF) con criptomonedas

Prepara y explica lo necesario para la declaración de IRPF de un ejercicio, a partir de los datos de CoinTracking. Trabaja en español y aplica la base de conocimiento (`knowledge/taxation/spain/`, `knowledge/cointracking/`).

> ⚠️ **No es asesoramiento fiscal.** El agente reconcilia, prepara y explica; **no calcula tu cuota ni produce cifras vinculantes** (ADR-006). Las cantidades exactas provienen del Informe de Impuestos de CoinTracking (método FIFO + España) o de un asesor.

> 🗓️ **Vigencia (ADR-008/ADR-022) — tres patas, no solo la fiscal.** Antes de dar por buena cualquier cifra:
> 1. **Normativa fiscal:** contrasta la "Vigencia" de los documentos de `knowledge/taxation/spain/` con el año y la fecha de hoy. Si el ejercicio no está cubierto o la verificación es antigua, **avísalo y reverifica** tramos, umbrales y criterios contra AEAT/BOE/DGT antes de darlos.
> 2. **Contexto regulatorio/operativo de los exchanges que usa el usuario** (ADR-022): la normativa cripto cambia de año en año y puede generar eventos reales que hay que revisar — licencias perdidas o ganadas (p. ej. MiCA en la UE), cierres, migraciones forzosas, restricciones de producto. Antes de preparar la declaración, **haz una búsqueda web breve** ("[exchange] regulación/licencia [año]", "[exchange] MiCA [año]") para los exchanges relevantes de la cuenta del ejercicio en curso, y si aparece algo material, avísalo y guía al usuario a revisar esas operaciones (ver `knowledge/reference/context/BINANCE_EU_MICA_EXIT.md` como ejemplo ya documentado). No es necesario si el usuario ya confirma que no ha habido cambios de exchange en el ejercicio.
> 3. **Formato/plataforma CoinTracking:** ver el resto de `knowledge/cointracking/`.

## Paso -1 — Proyecto activo obligatorio (ADR-013)

**Antes de cualquier otra cosa**, si esta conversación todavía no tiene un proyecto activo fijado: lista las subcarpetas de `USER_INPUT/` (los proyectos existentes) y **pregunta siempre con cuál quiere trabajar, incluso si solo hay uno** — nunca lo asumas ni lo anuncies como hecho ("trabajo sobre X salvo que digas lo contrario" no vale). Si no hay ninguno, ofrece crear el primero (pide un nombre). **Usa `AskUserQuestion` si está disponible** (es justo el caso: una decisión que solo puede tomar el usuario); si no, pregunta en texto plano. ⛔ **La pregunta cierra el turno: no añadas nada más después** (ni listas de pendientes, ni resumen de la sesión anterior), aunque tengas memoria de qué estaba pendiente. Ver `CLAUDE.md` §"Proyecto activo obligatorio". Una vez fijado (con confirmación explícita del usuario), reutilízalo el resto de la conversación (y si `/audit-cointracking` ya fijó uno antes en la misma conversación, reutiliza ese, no vuelvas a preguntar).

En cuanto quede fijado (o si ya lo estaba de una skill anterior en la misma conversación pero el MCP no se ha sincronizado aún), y si hay herramientas `cointracking_*` disponibles: llama a `cointracking_switch_project(project_name=<proyecto activo>)` antes de cualquier otra tool `cointracking_*`. Alinea en caliente el proyecto del MCP con el proyecto de datos activo, sin reiniciar el servidor (ADR-016). Si la respuesta trae `already_active: true`, no hace falta avisar de nada; si cambia de proyecto, ya no hace falta la advertencia de aislamiento que existía antes de ADR-016.

## Paso 0 — Diálogo de arranque (conversa antes de ejecutar)

Trata cada solicitud como si fuera un usuario nuevo (no asumas contexto previo). **No arranques en silencio: explica el plan y ofrece opciones.**

1. **Anuncia el plan** en lenguaje natural. Por ejemplo:
   > "Para preparar tu declaración de {AÑO}, primero haré una **auditoría** de tus datos para asegurarme de que las cifras son fiables, conectándome a tu cuenta de CoinTracking por la **API**. Luego prepararé el resumen fiscal del ejercicio."

2. **Confirma lo mínimo imprescindible:**
   - **Ejercicio fiscal** (año natural; "renta 2025" = 1 ene–31 dic 2025, se presenta en 2026). Si no lo dicen, pregúntalo.
   - **Perfil:** persona física **residente fiscal en España**, moneda EUR. Si no encaja, dilo (esta skill solo cubre ese caso).

3. **Ofrece la comprobación extra con el CSV** en lenguaje llano (sin decir "API/MCP/cotejo"); pregúntalo y **espera respuesta**:
   > "Voy a leer tus datos directamente de CoinTracking (conexión automática). Como comprobación adicional opcional, puedo compararlos con un archivo que descargues tú mismo desde CoinTracking; así, si algo no cuadra entre ambos, lo detecto. ¿Quieres hacer esa comprobación extra? Si sí, te guío para descargar el archivo."
   - Si acepta y no sabe cómo, **guíalo paso a paso** para exportar la lista de operaciones a CSV. Consulta los pasos exactos en `knowledge/cointracking/reference/CATALOG.md` (artículo de exportación/backup) antes de dárselos; no inventes rutas de menú. Pídele que **deje el archivo en la carpeta `USER_INPUT/<proyecto>/`** (el proyecto activo fijado en el Paso -1).

4. **Comprueba el acceso a datos:** el MCP de CoinTracking debe estar conectado (herramientas `cointracking_*`); y/o busca el CSV en **`USER_INPUT/<proyecto>/`**. Si no hay ninguna fuente, detente y pide al usuario que deje el archivo ahí (o que conecte el MCP).

Solo tras este diálogo, continúa con el Paso 1.

## Paso 1 — Reconciliar PRIMERO (puerta de calidad)

**No des ninguna cifra fiscal sobre datos sin reconciliar.**

- **Reutiliza si ya está hecho.** Si en **esta misma conversación** ya se ejecutó una auditoría (p. ej. el usuario invocó `/audit-cointracking` antes) sobre **la misma fuente y sin cambios desde entonces**, **reutiliza sus hallazgos**; **no la repitas** (evita gastar llamadas a la API, límite 60/h). Usa `CacheTTLManager` para automatizar esto con TTL por tipo y versionado (ver `tools/cache_ttl_manager.py`):
  ```python
  from tools.cache_ttl_manager import CacheTTLManager
  mgr = CacheTTLManager(project_name)
  # TTL automático por tipo de dato + invalida si knowledge/ subió de versión
  balance = mgr.get_or_fetch_dynamic('get_grouped_balance', {...}, mcp_call_fn=...)
  ```
  Dilo: "reutilizo la auditoría que acabamos de hacer (datos en caché)".
- **Re-audita solo si:** no se ha hecho aún, los **datos han cambiado** (el usuario corrigió operaciones), o cambió la fuente.
- Cuando toque auditar, usa la skill `audit-cointracking` o el subagente `cointracking-auditor`.

- Si aparecen hallazgos que **distorsionan la base de coste** (ventas sin base de coste, importación incompleta, transferencias mal emparejadas, ganancias implausibles p. ej. en stablecoins), **DETENTE**: informa de que la declaración no será fiable hasta corregirlos y lista los bloqueantes con su recomendación.
- Solo continúa cuando los datos estén razonablemente limpios (o el usuario acepte el riesgo explícitamente).

## Paso 2 — Eventos imponibles del ejercicio

**Optimización de tokens (ADR-039):** Procesa datos localmente, no en contexto LLM. Usa `tools/ct_audit.py` y `CacheTTLManager` para análisis local antes de pasar resultados compactos al contexto. Ejemplo:
```python
mgr = CacheTTLManager(project_name)
trades = mgr.get_or_fetch_dynamic('get_trades', {...}, mcp_call_fn=...)
# Procesar en Python (sin contexto):
taxable_events = analyze_trades_locally(trades)  # Devuelve resumen, no JSON crudo
# Pasar solo resumen al contexto (~200 tokens vs ~3000 sin procesar)
```

Con `cointracking_get_trades` acotado a `start`/`end` del año (UNIX **segundos**; convierte 1 ene 00:00 y 31 dic 23:59:59 de `Europe/Madrid` a UTC, ADR-005), identifica y clasifica:

- **Tributan como ganancia/pérdida patrimonial (base del ahorro):**
  - Ventas de cripto por fiat (EUR).
  - **Permutas cripto-cripto** (Art. 37.1.h LIRPF) — recuérdalo: *sí* tributan aunque no se pase por euros. Es el error más común.
  - Pagos de bienes/servicios con cripto.
- **No tributan aquí:** compras con fiat, holding, y **transferencias entre cuentas propias** (exclúyelas; solo trasladan coste).
- Usa el ticker completo (`SOL2` ≠ `SOL`) y trata comisiones como parte de la valoración. Reglas: `knowledge/taxation/spain/CAPITAL_GAINS.md` §1–5.

## Paso 3 — Ganancias/pérdidas (base del ahorro, FIFO)

- Método obligatorio en España: **FIFO** (`CAPITAL_GAINS.md` §4). En el MCP, `cointracking_get_gains(price:"oldest")` = FIFO.
- **Limitación clave (decláralo):** `get_gains` devuelve ganancias **de toda la vida**, no del ejercicio, y **no** genera el informe español. Por tanto:
  - El agente **verifica coherencia** y explica los eventos del año, pero **la cifra exacta del ejercicio** debe salir del **Informe de Impuestos de CoinTracking** (ajustes: método FIFO, jurisdicción España, año 2025) o de un asesor. Guía al usuario para generarlo con `knowledge/cointracking/WEB_APP_GUIDE.md` §7.
  - Toda cantidad que ofrezca el agente es **«estimación no vinculante»**.
  - ⛔ **Precondición de artefacto (no un bloqueo total, un gate de honestidad):** si al preparar la declaración **no tienes en el workspace** el Tax Report oficial del ejercicio (o su cifra `Resumen` ya documentada con evidencia en `REGISTRO-CAMBIOS.md`/memoria del proyecto), **no marques la cifra de base del ahorro como cerrada ni el informe como "listo para presentar"** — dilo explícitamente como `[VERIFICAR]` y pide al usuario ese artefacto (mismo procedimiento que ya usamos para BTC/USDC/OM, `DECISIONS.md#ADR-019`). Puedes seguir avanzando con el resto de secciones (eventos, rendimientos, Modelo 721) sin ese artefacto; solo el cierre de la cifra anual exacta depende de él.
- Tramos de la base del ahorro por ejercicio: `CAPITAL_GAINS.md` §6 (versionar por año; 2025 llega al 30 %).
- Compensación de pérdidas: `CAPITAL_GAINS.md` §7 (parcialmente `[PENDIENTE DE FUNDAMENTAR]`; no afirmar porcentajes/plazos sin fuente).

## Paso 4 — Rendimientos: staking, recompensas, airdrops, intereses

Aplica **[CAPITAL_INCOME.md](../../../knowledge/taxation/spain/CAPITAL_INCOME.md)** (fundamentado en consultas DGT). Clasifica cada operación del ejercicio y sepáralas por base:

- **`Staking`** → RCM, **base del ahorro** (DGT V1766-22). `[VERIFICAR]` si es staking delegado (pools).
- **`Ingresos por intereses`** (lending/earn) → RCM, **base del ahorro**.
- **`Recompensa / Bonificación`** → según origen: promocional/referido/airdrop → ganancia patrimonial, **base general** (DGT 0018-23 / V1948-21); recompensa de staking/earn → RCM, base del ahorro. **Inspecciona el origen de cada una.**
- **Minería** (si aplica) → actividad económica, **base general**.

Reglas de oro:
- **Valoración:** valor de mercado en EUR a la fecha de recepción.
- **Regla transversal (CAPITAL_INCOME §7):** ese valor es el **coste de adquisición** de esas monedas para el futuro FIFO. Verifica que CoinTracking lo haya registrado así (si no, aparecerán "ventas sin base de coste").
- **No calcules cifras vinculantes** (ADR-006); da estimaciones y remite al Informe oficial / asesor. En los puntos `[VERIFICAR]`, dilo y recomienda confirmación profesional.

## Paso 5 — Obligación informativa: Modelo 721

Comprueba si a **31/12 del ejercicio** el valor conjunto de criptos en **custodios NO residentes en España** supera **50.000 €** (`INFORMATIVE_OBLIGATIONS.md` §1) — es decir, comprueba el umbral, no asumas que aplica. Usa `get_grouped_balance(exchange)` + valoración e `get_historical_*`. ⚠️ **`get_historical_summary` puede devolver un punto fuera del rango pedido** (`MCP_API.md` — hallazgo empírico, ADR-020): filtra tú mismo por fecha el punto más cercano a 31/12, no confíes en que `end` corte exacto. Marca `[VERIFICAR]` el tratamiento de la autocustodia (Ledger/MetaMask). Recuerda plazo (1 ene–31 mar del año siguiente) y regla de repetición (+20.000 €).

## Paso 6 — Informe

Usa `templates/TAX_SUMMARY_ES.md`. Incluye: ejercicio y perfil; estado de reconciliación (bloqueantes); eventos imponibles del año; ganancia/pérdida **estimada** de base del ahorro (no vinculante, y de dónde sale la cifra exacta); rendimientos cuantificados pero **sin calificar** donde proceda; Modelo 721; y disclaimer + recomendación de validación profesional.

⛔ **Gate de cierre (Paso 3):** solo marca el informe como "listo para presentar" si la cifra anual exacta de base del ahorro está respaldada por el Tax Report oficial del ejercicio (o su cifra `Resumen` ya documentada). Si falta ese artefacto, dilo explícitamente en el informe (`[VERIFICAR]`) — no lo des por cerrado igualmente.

**Guárdalo** en `reports/output/<proyecto>/AAAA-MM-DD_declaracion_<ejercicio>.md` (ADR-011); no lo dejes solo en el chat. Registra en `reports/output/<proyecto>/REGISTRO-CAMBIOS.md` cualquier cambio aplicado y actualiza la memoria (`audit_state`).

## Recordatorio de límite de determinismo (ADR-006)

Encuentra, prepara y explica. **No** produzcas la cifra vinculante de la declaración. La calculadora determinista es el Informe de Impuestos de CoinTracking (FIFO/España) o el asesor.
