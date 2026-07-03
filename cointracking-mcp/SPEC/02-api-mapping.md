# 02 - API Mapping

## Herramientas Expuestas

| Tool | Método CT API | Descripción | Parámetros Clave |
|------|---------------|-------------|------------------|
| `cointracking_get_trades` | `getTrades` | Todas las operaciones (trades, depósitos, retiros, staking, minería, etc.) | `limit`, `start`, `end`, `order`, `trade_prices` |
| `cointracking_get_balance` | `getBalance` | Saldo actual por moneda | (ninguno) |
| `cointracking_get_grouped_balance` | `getGroupedBalance` | Saldo agrupado por ubicación (exchange, wallet) | (ninguno) |
| `cointracking_get_historical_summary` | `getHistoricalSummary` | Resumen histórico de la cartera (por fecha) | `start`, `end`, `fiat_currency` |
| `cointracking_get_historical_currency` | `getHistoricalCurrency` | Precio histórico de una moneda | `currency`, `start`, `end`, `fiat_currency` |
| `cointracking_get_gains` | `getGains` | Resumen de ganancias (por par de monedas) | (ninguno) |

## Esquema de Parámetros

### getTrades
```go
type GetTradesInput struct {
    Limit       int    `json:"limit,omitempty"`       // Máximo de resultados (recomendado para cuentas antiguas)
    Order       string `json:"order,omitempty"`       // "ASC" o "DESC" (default: DESC = más nuevos primero)
    Start       int64  `json:"start,omitempty"`       // UNIX timestamp en SEGUNDOS (no milisegundos)
    End         int64  `json:"end,omitempty"`         // UNIX timestamp en SEGUNDOS
    TradePrices int    `json:"trade_prices,omitempty"` // 0 o 1; si 1, incluir valores en BTC y fiat
}
```

**Nota crítica:** `start` y `end` son en **segundos**, no milisegundos. El agente debe convertir.

### getBalance
Sin parámetros. Siempre retorna saldo actual.

### getGroupedBalance
Sin parámetros. Retorna agrupación por exchange/wallet.

### getHistoricalSummary
```go
type GetHistoricalSummaryInput struct {
    Start        int64  `json:"start,omitempty"`        // UNIX timestamp en SEGUNDOS
    End          int64  `json:"end,omitempty"`
    FiatCurrency string `json:"fiat_currency,omitempty"` // Ej: "EUR", "USD"
}
```

### getHistoricalCurrency
```go
type GetHistoricalCurrencyInput struct {
    Currency     string `json:"currency"`              // Ej: "BTC", "ETH" (REQUERIDO)
    Start        int64  `json:"start,omitempty"`
    End          int64  `json:"end,omitempty"`
    FiatCurrency string `json:"fiat_currency,omitempty"`
}
```

### getGroupedBalance
Sin parámetros.

### getGains
Sin parámetros. Retorna resumen de ganancias por par (compra/venta).

## Autenticación (CoinTracking API)

**Método:** HMAC-SHA512 en POST

```
POST https://cointracking.info/api/v1/
Content-Type: application/x-www-form-urlencoded

Headers:
  Key: <API_KEY>
  Sign: <HMAC_SHA512_HEX(body, API_SECRET)>

Body (URL-encoded):
  method=<método>
  nonce=<nonce_estrictamente_creciente>
  [otros parámetros]
```

**Nonce:** debe ser estrictamente creciente por API key. Usar `time.Now().UnixMilli()` con contador interno.

## Formatos de Respuesta

Todas las respuestas vienen en JSON desde la API. El server debe:
1. Validar estructura con esquemas (Go structs con etiquetas JSON)
2. Cachear como JSON crudo
3. Retornar como JSON al cliente MCP (agente lo parsea)

## Rate Limits

El servidor aplica límites según tu plan de CoinTracking:

| Plan | Límite | Configuración |
|------|--------|----------------|
| PRO / EXPERT | 20 llamadas/hora | `COINTRACKING_ACCOUNT_TIER=pro` (default) |
| UNLIMITED | 60 llamadas/hora | `COINTRACKING_ACCOUNT_TIER=unlimited` |

**Comportamiento del server:**

1. **Por defecto:** comienza con `ACCOUNT_TIER=pro` (20 llamadas/hora) — conservador y seguro
2. **Track de llamadas en la ventana horaria:** mantiene conteo interno
3. **Si recibe 429 (rate limit excedido):**
   - Informa al agente: "Límite alcanzado. Si tu cuenta es UNLIMITED (60 llamadas/hora), configura `COINTRACKING_ACCOUNT_TIER=unlimited` como variable de entorno y reinicia."
   - El agente puede guiar al usuario a hacer el cambio
4. **Reporte de consumo actual en stats** — útil para debugging
5. **Sin reintento automático** — el agente decide qué hacer (esperar, cambiar tier, acotar consultas)

**Flujo típico si tiene UNLIMITED:**
```
1. User: "audita mi cuenta"
2. Server: comienza con 20 llamadas/hora
3. Agente pide muchos datos → Server alcanza 20 → 429
4. Server al agente: "Rate limit. Parece que tienes más límite. Si es UNLIMITED, configura ACCOUNT_TIER=unlimited"
5. User: export COINTRACKING_ACCOUNT_TIER=unlimited && ./cointracking-mcp
6. Agente reintenta → ahora 60 llamadas/hora
```

## Error Handling

La API de CT retorna errores en dos formas:

1. **HTTP 429:** Rate limit
2. **HTTP ≠ 2xx:** Otros errores (credenciales inválidas, parámetros malos, etc.)
3. **`{ "success": 0, "error": "..." }`:** Error a nivel de API (aunque HTTP 200)

El server debe:
- Diferenciar entre errores de red, autenticación, parámetros, rate limit
- Retornar errores de forma clara al agente
- No cachear respuestas de error (excepto quizá algunos)

## Caché — Claves

Clave en caché = hash determinista de `(método + parámetros normalizados)`

```
getTrades + {limit: 100, start: 1234567890} 
  → key = sha256("getTrades|end=0|limit=100|order=|start=1234567890|trade_prices=0")
```

Normalización:
- Ordenar parámetros alfabéticamente
- Omitir valores vacíos/nil
- Convertir números a string con formato fijo

## Integración con Auditor

El agente llama a estas herramientas para:
1. **Auditoría:** `getTrades` + validaciones deterministas (duplicados, huérfanos, saldos)
2. **Preparación fiscal:** `getGains`, históricos para reconciliación
3. **Investigación:** `getBalance` para estado actual, `getHistoricalCurrency` para precios

Típicamente el auditor hace:
1. Obtener trades con `getTrades` (completo, sin limit)
2. Obtener balance actual con `getBalance`
3. Validar consistencia (el server puede ayudar con caché + validaciones)
4. Si encuentra inconsistencias, pide al usuario que corrija en CT web
5. Invalida caché y re-ejecuta
