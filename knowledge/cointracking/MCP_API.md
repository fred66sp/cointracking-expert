# Integración con la API de CoinTracking vía MCP

**Tipo:** Documentación de integración (obra propia)
**Servidor MCP:** `cointracking-mcp` — implementación propia en Go, en `cointracking-mcp/` (spec en `cointracking-mcp/SPEC/`). Sustituye a partir de 2026-07-03 al servidor JS de terceros usado antes (`cointracking-mcp-main/`, de Alessandro Melazzini, MIT); ese repo sigue en el filesystem como referencia de puerto pero ya no es el que arranca `.mcp.json`.
**Última verificación:** 2026-07-03
**Vigencia:** herramientas y parámetros del MCP y límites de la API observados a 2026-07. Pueden cambiar con nuevas versiones del servidor/API — reverificar contra las herramientas `cointracking_*` realmente disponibles en la sesión si esta fecha es antigua (ADR-008).
**Estado:** Operativo; verificado contra la API real (`getBalance`, `getGroupedBalance`, `getGains`, `getTrades`, `getHistoricalCurrency`)

Esta es la vía de **datos en vivo** del agente auditor (ADR-006), complementaria al CSV export. El servidor MCP es **solo lectura** y expone la API de usuario de CoinTracking como herramientas `cointracking_*`, con caché propia (memoria + disco) por proyecto.

> El servidor (`cointracking-mcp/`) es código propio de este repositorio (Go), gobernado por las mismas normas de ADR-006/009. El binario compilado (`cointracking-mcp/dist/`) y la caché (`.cache/cointracking/`) están en `.gitignore` — no se versionan. La configuración de arranque está en `.mcp.json` (y `.vscode/mcp.json` para el uso desde VS Code).

---

## Herramientas disponibles

### Datos de CoinTracking (solo lectura)

| Herramienta | Para qué | Parámetros |
|-------------|----------|------------|
| `cointracking_get_trades` | Todas las operaciones (trades, depósitos, retiradas, staking, minería, airdrops, DeFi…). Devuelve compra/venta, monedas, comisiones, tipo, exchange, grupo, comentario, `imported_from`, `time`, `trade_id`. | `limit`, `order` (ASC/DESC), `start`, `end` (**UNIX segundos**), `trade_prices` (0/1) |
| `cointracking_get_balance` | Balance actual por moneda, con valor en BTC y fiat. | — |
| `cointracking_get_grouped_balance` | Balances agrupados por `exchange`, `type` o `currency`. | `group` (**obligatorio**: `exchange`\|`type`\|`currency`), `exclude_dep_with` (0/1), `type` (filtro opcional) |
| `cointracking_get_gains` | Ganancias realizadas y no realizadas, con **método de coste seleccionable**. | `price`: `best`/`worst`/`oldest`/`newest`, `btc` (0/1) |
| `cointracking_get_historical_summary` | Resumen histórico de cartera agregado por año/mes. | `btc` (0/1), `start`, `end`, `fiat_currency` |
| `cointracking_get_historical_currency` | Balance y valor histórico de una moneda concreta. | `currency` (**obligatorio**), `start`, `end`, `fiat_currency` |

### Control de caché y proyecto (propias de esta implementación, no existían en el servidor JS)

| Herramienta | Para qué | Parámetros |
|-------------|----------|------------|
| `cointracking_invalidate_cache` | Invalida entradas cacheadas por patrón (`getTrades*`, `*`, lista separada por comas). **Úsala siempre que el usuario haya modificado datos en la web de CoinTracking**, antes de volver a consultar — si no, se puede mezclar caché obsoleta con datos frescos (ADR-010). | `pattern` (opcional; vacío o `*` = todo) |
| `cointracking_cache_stats` | Tamaño de caché, hit rate, llamadas a la API por método, y cuota horaria consumida. Útil para decidir si merece la pena volver a consultar. | — |
| `cointracking_close_project` | Señala el fin de una sesión de auditoría: vacía la caché de memoria y confirma que está persistida en disco. Llamarla al terminar. | `project_name` (opcional; debe coincidir con el proyecto configurado) |

### 🔑 Método de coste y España (crítico)

En `cointracking_get_gains`, el parámetro `price` mapea así:

- `oldest` = **FIFO** ← **el que exige España** (ver `../taxation/spain/CAPITAL_GAINS.md` §4)
- `newest` = LIFO
- `best` = menor coste · `worst` = mayor coste

> Para cualquier estimación fiscal española, el agente debe pedir `price: "oldest"`. Aun así, conforme a ADR-006, el resultado es **estimación no vinculante** (lo calcula CoinTracking, no un cálculo determinista propio: `tools/ct_audit.py` solo cubre chequeos mecánicos, no FIFO).

---

## Mapeo herramienta → chequeo de auditoría

| Chequeo (skill `audit-cointracking`) | Herramienta(s) MCP |
|--------------------------------------|--------------------|
| Completitud de importación | `get_grouped_balance` (por exchange) + `get_balance` |
| Transferencias huérfanas | `get_trades` (filtrar depósitos/retiradas) |
| Ventas sin base de coste | `get_gains` (avisos) + `get_trades` |
| Duplicados | `get_trades` |
| Saldos negativos imposibles | `get_balance`, `get_grouped_balance` |
| Orden temporal de transferencias | `get_trades` (campo `time`) |
| Coherencia fiscal (cualitativa) | `get_gains` con `price: "oldest"` |
| Tenencias a 31/12 (Modelo 721) | `get_historical_summary`, `get_historical_currency` |

---

## Caché y eficiencia de tokens (ADR-010)

El servidor **ya cachea internamente** (memoria L1 + disco SQLite L2, por proyecto, en `.cache/cointracking/agp/`), con TTL por tipo de dato (trades/gains: 1h; balance: 10min; históricos: 2h). Esto es transparente: el agente llama a las tools normalmente y el servidor decide caché vs. API. Aun así:

- **Invalida tras cambios:** si guías al usuario a modificar algo en la web de CoinTracking, llama a `cointracking_invalidate_cache` antes de volver a consultar (si no, puede devolver datos cacheados obsoletos hasta que expire el TTL).
- **Consultas dirigidas:** acota `start`/`end` y `limit`; prioriza agregados (`get_grouped_balance`, `get_gains`) sobre `get_trades` completo — reduce tanto el consumo de cuota como el tamaño de la respuesta.
- **Datos grandes → código:** para el historial de operaciones, vuelca la respuesta a fichero y **procésalo con un script** (filtrar por año, detectar huérfanas/duplicados, sumar); sube al contexto solo el resultado. No pegues el JSON crudo completo.
- **Cierre de sesión:** al terminar una auditoría, llama a `cointracking_close_project` para liberar memoria y confirmar que la caché quedó persistida en disco.

## Límites y buenas prácticas

- **Límite de tasa:** 20 llamadas/hora (plan pro/expert) o 60/hora (unlimited) — configurado al arrancar el servidor con `--tier`. Consulta `cointracking_cache_stats` para ver la cuota consumida en la hora en curso antes de lanzar consultas grandes. Ante HTTP 429, esperar y **acotar** las consultas; el servidor no reintenta automáticamente.
- **Tiempos:** `start`/`end` en **segundos UNIX** (no ms). Coherente con ADR-005, convertir desde fechas UTC.
- **Historiales grandes:** en `get_trades`, pasar siempre `limit` y `start`/`end`.
- **Validación cruzada:** cuando exista también el CSV, comparar ambos para detectar discrepancias de importación (ADR-006, doble vía).
- **Verificación de remediaciones por lote (ADR-010 punto 7):** durante la remediación guiada, no llames a ninguna herramienta `cointracking_*` para comprobar cada corrección individual. Guía primero todo el lote de correcciones que el usuario vaya a aplicar en esta ronda; solo cuando confirme que ha terminado, llama **una vez** a `cointracking_invalidate_cache` y luego a las herramientas de lectura necesarias para verificar **todos** los hallazgos corregidos en la menor cantidad de consultas posible (agregados como `get_grouped_balance`/`get_gains` antes que `get_trades` completo).

---

## Cómo conectarlo (una vez)

1. Compilar el servidor: en `cointracking-mcp/`, `go build -o dist/cointracking-mcp.exe ./cmd/cointracking-mcp` (el binario no se versiona, hay que compilarlo tras clonar).
2. Obtener credenciales en CoinTracking → **Account → API** (basta **solo lectura**).
3. Proporcionar las credenciales **sin commitearlas nunca**: variables de entorno `COINTRACKING_API_KEY` y `COINTRACKING_API_SECRET`. En esta máquina ya están definidas como variables de entorno **de usuario** persistentes (no hace falta re-exportarlas cada sesión); el proceso del servidor las hereda automáticamente. En otra máquina, exportarlas antes de lanzar Claude Code, o pasarlas por `--api-key`/`--api-secret` (nunca committeadas).
4. Reiniciar Claude Code para que cargue el servidor (arranca vía `.mcp.json`/`.vscode/mcp.json`, proyecto `agp`); entonces aparecen las herramientas `cointracking_*`.

> 🔒 **Seguridad:** la clave/secreto de la API son sensibles. Nunca deben aparecer en el repositorio, en `.mcp.json` ni en `.vscode/mcp.json` (ambos se versionan). Usar solo variables de entorno.
