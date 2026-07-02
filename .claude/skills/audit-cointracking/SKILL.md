---
name: audit-cointracking
description: Audita una cuenta o exportación de CoinTracking. Reconcilia los datos y detecta transferencias huérfanas, ventas sin base de coste, duplicados, saldos imposibles e incoherencias fiscales españolas, explicando cada hallazgo con evidencia. Usa el MCP de la API de CoinTracking si está disponible y/o el CSV export.
---

# Auditoría de CoinTracking

Ejecuta una auditoría de reconciliación sobre los datos de CoinTracking del usuario, siguiendo el playbook de abajo. Trabaja en español y aplica la base de conocimiento del repo (ver el subagente `cointracking-auditor` y `knowledge/`).

## Paso 0 — Diálogo de arranque y preparación

**Conversa antes de ejecutar, en lenguaje llano** (el usuario no domina CoinTracking; evita "API/MCP/cotejo"). Anuncia qué vas a hacer y **ofrece una comprobación extra con un CSV**, esperando respuesta. Por ejemplo:
> "Voy a revisar tu cuenta leyendo tus datos directamente de CoinTracking. Como comprobación adicional opcional, puedo compararlos con un archivo que descargues tú mismo; así detecto si algo no cuadra entre ambos. ¿Quieres hacer esa comprobación extra (te guío para descargar el archivo) o sigo solo con la conexión automática?"
- Si acepta y no sabe cómo, **guíalo paso a paso** para exportar la lista de operaciones a CSV; consulta los pasos exactos en `knowledge/cointracking/reference/CATALOG.md` (artículo de exportación/backup) y no inventes rutas de menú.

Después:

1. **Carga el conocimiento**: lee `knowledge/cointracking/CSV_FORMAT.md`, `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md`, `knowledge/taxation/spain/CAPITAL_GAINS.md`. Son las reglas que aplicarás y citarás.
2. **Localiza los datos** (ADR-006, ambas vías):
   - **MCP de CoinTracking**: si hay herramientas MCP de CoinTracking disponibles en la sesión, úsalas para datos en vivo. Si el usuario menciona un MCP pero no está conectado, pídele que lo conecte (`/mcp`).
   - **CSV export**: busca una exportación "Trade Table" (p. ej. en la raíz del proyecto). El CSV está en `.gitignore` por privacidad; puede existir localmente.
   - Si hay ambos, usa el CSV como validación cruzada del MCP.
   - Si no hay ninguno, dilo y detente: no inventes datos.
3. **Normaliza mentalmente** según el conocimiento: fechas a UTC desde `Europe/Madrid` con DST (ADR-005), importes con precisión decimal, ticker completo con sufijo (`SOL` ≠ `SOL2`), parsear por posición (3 columnas `Cur.`).

## Paso 1 — Ejecuta el playbook de chequeos

Para cada activo/cuenta según proceda. Cada chequeo cita su regla de conocimiento:

1. **Completitud de importación** — ¿faltan cuentas/exchanges (seguimiento de un solo lado)? Provoca base de coste 0 y ganancias infladas. *(COST_BASIS §3.2, "READ FIRST")*
2. **Transferencias huérfanas** — retiradas sin depósito emparejado y viceversa. Empareja por **Tx Hash** (nivel 1, fuerte) y, si falta, por **moneda + importe ≈ retirada − comisión + ventana temporal + cuentas distintas** (nivel 2, heurístico con confianza). *(CSV_FORMAT §7)*
3. **Ventas sin base de coste** — activos vendidos/permutados sin compra previa registrada. Distingue "compra importada como depósito" de "compra realmente ausente". Nunca asumas base 0 en silencio. *(COST_BASIS §3)*
4. **Duplicados** — filas idénticas: separa duplicado por reimportación (error) de repetición legítima (comisiones/recompensas recurrentes). *(CSV_FORMAT §9, COST_BASIS §4.2)*
5. **Saldos negativos imposibles** — vender/enviar más de lo disponible de un activo. Distingue del **FIAT negativo**, que suele ser artefacto de no importar depósitos fiat, no una imposibilidad. *(COST_BASIS §4.1)*
6. **Orden temporal de transferencias** — un depósito anterior a su retirada rompe la transferencia de base de coste. Verifica `retirada ≤ depósito` tras normalizar a UTC. *(COST_BASIS §2, ADR-005)*
7. **Coherencia fiscal (cualitativa)** — ¿se tratan las permutas cripto-cripto como hecho imponible? ¿se usaría el método **FIFO** (España) y no el pool promediado? Señala riesgos; **no** des la cifra final vinculante. *(CAPITAL_GAINS §3-4, COST_BASIS §1)*
8. **Colisión de tickers** — no fusionar activos con sufijo (`SOL2`, `WLD3`, `ID2`). *(CSV_FORMAT §8)*

Para análisis profundo puedes delegar en el subagente `cointracking-auditor`.

## Paso 2 — Informe

Usa la plantilla `templates/AUDIT_REPORT.md`. Ordena los hallazgos por severidad y, para cada uno: **causa, evidencia, impacto, recomendación**. Cierra con un resumen y con lo **no verificado** (datos faltantes, reglas pendientes de fundamentar como la fiscalidad de staking).

## Límite de determinismo (ADR-006)

Recuerda: encuentras y explicas; **no** produces cifras fiscales vinculantes. Toda cantidad exacta es «estimación no vinculante» salvo que provenga de un cálculo determinista.
