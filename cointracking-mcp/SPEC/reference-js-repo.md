# Reference - Repo JS Original

## Ubicación

`H:\cointracking-expert\cointracking-mcp-main`

Repo oficial: https://github.com/alemuenchen/cointracking-mcp

## Estructura

```
cointracking-mcp-main/
├── src/
│   ├── index.ts                 # Punto de entrada; registra 6 tools
│   ├── api-client.ts            # Cliente HTTP a API de CT (autenticación HMAC-SHA512)
│   ├── tools/                   # Una herramienta por archivo
│   │   ├── get-trades.ts
│   │   ├── get-balance.ts
│   │   ├── get-grouped-balance.ts
│   │   ├── get-historical-summary.ts
│   │   ├── get-historical-currency.ts
│   │   └── get-gains.ts
│   └── utils/
│       ├── errors.ts            # Clases de error
│       └── formatting.ts        # Formateo JSON
├── tests/
│   ├── api-client.test.ts
│   ├── errors.test.ts
│   ├── formatting.test.ts
│   └── schemas.test.ts
├── package.json                 # Node.js, TypeScript, vitest
├── tsconfig.json
└── node_modules/ (ignorar)
```

## Archivos Clave para Porter

### `src/api-client.ts`

**Qué hace:**
- Construye URL-encoded POST body
- Computa HMAC-SHA512 (Key header)
- Maneja nonce estrictamente creciente
- Retorna parsed JSON o lanza CoinTrackingError

**Porter a Go:**
- Usar `crypto/hmac`, `crypto/sha512` (stdlib)
- `url.Values` para URL encoding
- `net/http.Client` para llamadas
- Nonce: usar `time.Now().UnixMilli()` + contador atómico

### `src/tools/*.ts`

**Estructura:**
```typescript
export const getTradesSchema = { limit: z.number(), ... }
export const getTradesDefinition = { description: "...", inputSchema, annotations }
export async function getTradesHandler(args: { limit?, ... }) { 
  const data = await coinTrackingRequest("getTrades", args)
  return textResult(formatJson(data))
}
```

**Porter a Go:**
- Schema → Go structs con etiquetas JSON
- Handler → función que acepta struct con parámetros
- Retornar MCP ToolResult con contenido formateado

### `src/utils/errors.ts`

**Tipos de error:**
- `CoinTrackingError` — base
- `RATE_LIMIT` — 429 de CT
- `NETWORK` — conectividad
- `BAD_JSON` — respuesta no-JSON
- `API_ERROR` — `{ success: 0, error: "..." }`

**Porter a Go:**
```go
type CoinTrackingError struct {
  Message string
  Code    string // RATE_LIMIT, NETWORK, etc.
}
```

## Decisiones de Diseño del Original (JS)

### ✅ Mantener
- **6 tools separadas** — clara división de responsabilidades
- **Schema validation con Zod** → Go: usar structs + json.Unmarshal + custom validators
- **HMAC-SHA512** — correcto, mantener exacto
- **Nonce monotónico** — crítico, mantener lógica
- **Error types específicos** — mejor debugging

### ⚠️ Mejorar en Go
- **Sin caché** → Agregar LRU + TTL
- **Sin validaciones deterministas** → Agregar tool de validation
- **Sin logging estructurado** → Agregar leveled logging (debug, info, warn, error)
- **Sin persistencia** → Opcionalmente, SQLite o JSON
- **Tests limitados** → Agregar fixtures reales (sample_trades.csv)

## Parámetros API - Exactitud Crítica

### getTrades
- **limit:** limita resultados (recomendado para cuentas antiguas)
- **order:** "ASC" o "DESC" (default: DESC)
- **start/end:** UNIX timestamp en **SEGUNDOS**, no milisegundos ⚠️
- **trade_prices:** 0 o 1; si 1, incluye valores en BTC y fiat

**Nota crítica:** Timestamps en SEGUNDOS. El repo JS no hace conversión; el usuario/agente debe pasar segundos.

### getBalance
Sin parámetros.

### getGroupedBalance
Sin parámetros.

### getHistoricalSummary
- **start/end:** UNIX seconds
- **fiat_currency:** ej "EUR", "USD"

### getHistoricalCurrency
- **currency:** REQUERIDO; ej "BTC", "ETH"
- **start/end:** UNIX seconds
- **fiat_currency:** ej "EUR"

### getGains
Sin parámetros.

## Rate Limit Handling

La API retorna:
- **HTTP 429:** "Rate limit exceeded (60 calls/hour...)"

El repo JS:
```javascript
if (res.status === 429) {
  throw new CoinTrackingError(
    "Rate limit exceeded (60 calls/hour...). Wait before retrying...",
    "RATE_LIMIT"
  );
}
```

**En Go:** Mismo manejo; no reintentar automáticamente.

## Testing Fixtures

No existen fixtures en el repo JS original. Necesitamos crear:
- `tests/fixtures/sample_trades.csv` — export de CoinTracking sintético
- `tests/fixtures/sample_balance.json` — respuesta de getBalance
- `tests/expected.json` — resultado esperado de validaciones

## Publicación / Distribución

El repo JS tiene:
- `package.json` con `bin: { "cointracking-mcp": "dist/index.js" }`
- Se puede instalar vía npm: `npm install -g cointracking-mcp`

**En Go:**
- `main.go` en `cmd/cointracking-mcp/`
- Makefile con targets: `build`, `test`, `install`
- GitHub Releases con binarios pre-compilados (Windows, macOS, Linux)

## Notas de Código

- Variables en camelCase (seguir Go idioms: camelCase para var, PascalCase para exported)
- Comentarios en inglés (Go convention); strings de usuario en español
- Error handling: no panics; retornar errores
- Concurrency: si hay caché en memory, puede haber acceso concurrente (usar sync.RWMutex si es necesario)
