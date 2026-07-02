# Integración con la API de CoinTracking vía MCP

**Tipo:** Documentación de integración (obra propia)
**Servidor MCP:** `cointracking-mcp` (proyecto externo de terceros, MIT, de Alessandro Melazzini)
**Última verificación:** 2026-07-02
**Vigencia:** herramientas y parámetros del MCP y límites de la API observados a 2026-07. Pueden cambiar con nuevas versiones del servidor/API — reverificar contra las herramientas `cointracking_*` realmente disponibles en la sesión si esta fecha es antigua (ADR-008).
**Estado:** Documentado; requiere credenciales para operar

Esta es la vía de **datos en vivo** del agente auditor (ADR-006), complementaria al CSV export. El servidor MCP es **solo lectura** y expone la API de usuario de CoinTracking como herramientas `cointracking_*`.

> ⚠️ El servidor es un proyecto **independiente y no oficial** (no afiliado a CoinTracking), con licencia MIT. Se trata como **dependencia externa local**: no se versiona en este repositorio (está en `.gitignore`). La configuración de arranque está en `.mcp.json`.

---

## Herramientas disponibles (todas de solo lectura)

| Herramienta | Para qué | Parámetros |
|-------------|----------|------------|
| `cointracking_get_trades` | Todas las operaciones (trades, depósitos, retiradas, staking, minería, airdrops, DeFi…). Devuelve compra/venta, monedas, comisiones, tipo, exchange, grupo, comentario, `imported_from`, `time`, `trade_id`. | `limit`, `order` (ASC/DESC), `start`, `end` (**UNIX segundos**), `trade_prices` (0/1) |
| `cointracking_get_balance` | Balance actual por moneda, con valor en BTC y fiat. | — |
| `cointracking_get_grouped_balance` | Balances agrupados por `exchange`, `type` o `currency`. | grupo |
| `cointracking_get_gains` | Ganancias realizadas y no realizadas, con **método de coste seleccionable**. | `price`: `best`/`worst`/`oldest`/`newest`, `btc` (0/1) |
| `cointracking_get_historical_summary` | Resumen histórico de cartera agregado por año/mes. | — |
| `cointracking_get_historical_currency` | Balance y valor histórico de una moneda concreta. | moneda |

### 🔑 Método de coste y España (crítico)

En `cointracking_get_gains`, el parámetro `price` mapea así:

- `oldest` = **FIFO** ← **el que exige España** (ver `../taxation/spain/CAPITAL_GAINS.md` §4)
- `newest` = LIFO
- `best` = menor coste · `worst` = mayor coste

> Para cualquier estimación fiscal española, el agente debe pedir `price: "oldest"`. Aun así, conforme a ADR-006, el resultado es **estimación no vinculante** (lo calcula CoinTracking, no nuestro motor determinista).

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

## Límites y buenas prácticas

- **Límite de tasa:** 60 llamadas/hora (cuentas Unlimited). Ante HTTP 429, esperar y **acotar** las consultas. No hacer barridos innecesarios.
- **Tiempos:** `start`/`end` en **segundos UNIX** (no ms). Coherente con ADR-005, convertir desde fechas UTC.
- **Historiales grandes:** en `get_trades`, pasar siempre `limit` y `start`/`end`.
- **Validación cruzada:** cuando exista también el CSV, comparar ambos para detectar discrepancias de importación (ADR-006, doble vía).

---

## Cómo conectarlo (una vez)

1. Instalar el servidor localmente (repo externo) en `cointracking-mcp-main/` en la raíz del proyecto (`npm install && npm run build`; el `dist/` ya viene compilado).
2. Obtener credenciales en CoinTracking → **Account → API** (basta **solo lectura**).
3. Proporcionar las credenciales **sin commitearlas nunca**: variables de entorno `COINTRACKING_API_KEY` y `COINTRACKING_API_SECRET` (las lee `.mcp.json` por expansión `${...}`), o el `.env` del propio servidor.
4. Reiniciar Claude Code para que cargue el servidor; entonces aparecen las herramientas `cointracking_*`.

> 🔒 **Seguridad:** la clave/secreto de la API son sensibles. Nunca deben aparecer en el repositorio ni en `.mcp.json`. Usar variables de entorno o el `.env` local (ya ignorado por git).
