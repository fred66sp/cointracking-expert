# 03 - Cache Strategy

## Objetivo

Reducir llamadas a la API de CoinTracking (según el límite de tu plan: 20 o 60 por hora) y acelerar respuestas, manteniendo datos frescos y permitiendo invalidación manual cuando el usuario modifique datos en la web.

## Arquitectura

### En Memoria (L1, Primaria - Fast)

- **Tipo:** LRU (Least Recently Used) + TTL por entrada
- **Tamaño:** configurable, default según `ACCOUNT_TIER`:
  - `pro`: 5000 entradas (cubre 3.500 registros de CT)
  - `expert`: 50000 entradas (cubre 20–100k registros de CT)
  - `unlimited`: 150000 entradas (sin límite de registros)
- **Evicción:** cuando se alcanza el límite, eliminar entrada menos usada recientemente
- **Expiración:** cada entrada tiene TTL; transcurrido, se considera expirada

### En Disco (L2, Secundaria - Persistent)

- **Ubicación:** `${CACHE_PERSIST_DIR}/${PROJECT_NAME}/cointracking-mcp.db` (por proyecto)
  - Default: `./cache/{PROJECT_NAME}/`
  - Si `PROJECT_NAME=default` → `./cache/default/`
  - Si `PROJECT_NAME=agp` → `./cache/agp/`
- **Propósito:** persistir caché entre reinicios, aislada **por proyecto**
- **Estrategia:** escribir a disco **automáticamente** (asincronía, no bloquea respuestas)
- **Limpieza:** eliminar entradas expiradas al startup y periódicamente

## Estrategia de Lectura (Hit/Miss)

```
Agente pide: getTrades(limit=100)
↓
1. ¿En memoria (L1)?
   Sí → retornar datos + marcar "FROM_MEMORY"
   No → paso 2
   
2. ¿En disco (L2)?
   Sí → cargar a memoria + retornar datos + marcar "FROM_DISK"
   No → paso 3
   
3. ¿Llamar a API?
   Sí → obtener datos + guardar en memoria + guardar en disco + retornar
```

**Validación de TTL:**
- Aplicada en todos los niveles (memoria y disco)
- Si expirado → eliminar y tratar como miss
- Server valida timestamp actual contra `expires_at` de cada entrada

## Aislamiento por Proyecto

Cada instancia del servidor tiene un `PROJECT_NAME` (default: `default`). La caché es completamente aislada por proyecto:

**En memoria:**
- Caché LRU separada por proyecto
- Si cambias `PROJECT_NAME`, cambias de LRU (o cargas una nueva si es la primera vez)

**En disco:**
- Cada proyecto tiene su propio directorio: `./cache/{PROJECT_NAME}/`
- Los datos de un proyecto nunca se mezclan con otro
- Útil para auditar múltiples cuentas de CT en paralelo

**Ejemplo:**
```
./cache/
├── default/              ← PROJECT_NAME=default (default)
│   └── cointracking-mcp.db
├── agp/                  ← PROJECT_NAME=agp
│   └── cointracking-mcp.db
└── cliente_b/            ← PROJECT_NAME=cliente_b
    └── cointracking-mcp.db
```

**Cambiar de proyecto:**
```bash
# Terminal 1: auditar proyecto A
PROJECT_NAME=agp ./cointracking-mcp

# Terminal 2: auditar proyecto B (caché aislada)
PROJECT_NAME=cliente_b ./cointracking-mcp
```

**Nota:** Si usas el mismo `PROJECT_NAME` en múltiples instancias del server simultáneamente, la caché en disco puede sufrir race conditions. Cada proyecto debería ser auditado por una sola instancia a la vez (o usar solo caché en memoria sin persistencia).

## TTL por Tipo de Dato

Los tiempos de expiración reflejan qué tan frecuentemente cambian los datos:

| Tipo | TTL | Razón |
|------|-----|-------|
| `getTrades` | 3600 segundos (1h) | Usuario modifica trades en web; cambios poco frecuentes |
| `getBalance` | 600 segundos (10 min) | Más volátil; precios cambian, depósitos/retiros pueden llegar |
| `getGroupedBalance` | 600 segundos (10 min) | Similar a balance |
| `getHistoricalSummary` | 7200 segundos (2h) | Datos históricos; estables mientras no se modifiquen trades |
| `getHistoricalCurrency` | 7200 segundos (2h) | Precios históricos; no cambian |
| `getGains` | 3600 segundos (1h) | Recalculado desde trades; expira igual que trades |

**Configurable por variables de entorno:**
```bash
export CACHE_TTL_TRADES=3600
export CACHE_TTL_BALANCE=600
export CACHE_TTL_SUMMARY=7200
export CACHE_TTL_GAINS=3600
export CACHE_TTL_HISTORICAL=7200
./cointracking-mcp
```

## Clave de Caché

Clave determinista = hash de `(método + parámetros normalizados)`:

```
getTrades(limit=100, start=1234567890)
→ "getTrades|limit=100|start=1234567890|end=0|order=|trade_prices=0"
→ hash_key = sha256(...)
```

**Normalización:**
1. Ordenar parámetros alfabéticamente
2. Incluir todos los parámetros conocidos (con valores default si no se especifican)
3. Números en formato fijo (no decimales innecesarios)
4. Strings en minúscula (ej: "ASC" → "asc" para uniformidad)

## Hit/Miss Logic

```
1. Agente llama cointracking_get_trades(limit=100, start=...)
2. Server computa clave = sha256(...)
3. Buscar en LRU[clave]
   a. Si existe Y no expirado → retornar data + marcar como "FROM_CACHE"
   b. Si existe Y expirado → eliminar entrada, continuar al paso 4
4. Miss en caché → llamar API
5. Si API éxito → cachear(clave, data, ttl) + retornar data
6. Si API error → retornar error (NO cachear errores)
```

## Invalidación

### Manual (Agente Iniciada)

Agente puede pedir invalidación cuando el usuario modifique datos en CT web:

```
Tool: cointracking_invalidate_cache
Input: {
  pattern?: string  // Ej: "getTrades*", "*" (todas)
}
Output: {
  invalidated: int  // Entradas eliminadas
}
```

**Patrones útiles:**
- `getTrades*` — invalidar todas las variantes de getTrades
- `getBalance` — invalidar solo balance
- `*` — limpiar todo el caché
- `getGains,getHistoricalSummary` — múltiples

**Típicamente, el agente guía así al usuario:**
1. "He encontrado inconsistencias. Vamos a arreglarlo en CoinTracking."
2. [Usuario modifica datos en web]
3. Agente: "Listo. Ahora invalido la caché para que tengamos datos frescos."
4. Llamada a `cointracking_invalidate_cache("getTrades*")`
5. Agente re-ejecuta auditoría

### Automática

No hay invalidación automática basada en tiempo. El agente es responsable de detectar cuándo los datos han cambiado y pedir invalidación.

## Logging de Caché

Cada consulta se loguea (nivel `debug`):

```
[2026-07-02 10:30:45] getTrades(limit=100, start=1234567890) → CACHE HIT (expires at 11:30:45)
[2026-07-02 10:35:12] getBalance() → CACHE MISS → API CALL (50ms) → cached
[2026-07-02 10:40:00] invalidateCache("getTrades*") → 5 entries removed
```

El agente puede solicitar logs via herramienta especial `cointracking_cache_stats`:

```
Output: {
  size: 42,
  hits: 150,
  misses: 23,
  total_calls: 173,
  hit_rate: 0.867,
  calls_to_api: [
    { method: "getTrades", count: 10, last_call: "2026-07-02T10:35:12Z" },
    { method: "getBalance", count: 5, last_call: "2026-07-02T10:40:00Z" }
  ]
}
```

## Estrategia de Persistencia (Fase 1+)

Si `CACHE_PERSIST_DIR` se configura:

1. **Al startup:** cargar caché de disco, eliminar entradas expiradas
2. **En tiempo real:** escribir nuevas entradas a disco (asincronía, con batching cada 5s)
3. **Al shutdown:** volcar caché a disco (graceful shutdown)

Formato: SQLite con tabla:
```
CREATE TABLE cache (
  key TEXT PRIMARY KEY,
  method TEXT,
  value BLOB,              -- JSON comprimido con gzip
  cached_at TIMESTAMP,
  expires_at TIMESTAMP
)
```

## Comportamiento Bajo Rate Limit

Si la API retorna 429 (rate limit):

1. **Aceptar:** API está saturada, avisar al agente
2. **No cachear:** los errores 429 no se cachean
3. **No reintentar automáticamente:** el agente es responsable de reintentarlo después
4. **Marcar en logs:** "RATE_LIMIT_EXCEEDED"

## Validación

El servidor valida parámetros **antes** de cachear/consultar API:

- ¿Parámetros conocidos? Si no, error
- ¿Valores válidos? (ej: `order` debe ser "ASC" o "DESC")
- ¿Timestamps en segundos? (convertir de milisegundos si es necesario)
- ¿Currencies válidas? (contrachequear contra catálogo de CT)

Si falla validación → error 400, no cachear.

## Estadísticas (Fase 1)

Recopilar y exponer:

- Cache size (actual / máximo)
- Hit rate (hits / total_calls)
- Calls by method
- Última llamada a API (cuándo fue, qué método)
- Tasa de consumo del límite de rate limit configurado (llamadas usadas / límite horario)

Útil para debugging y decisiones del agente ("¿vale la pena pedir esto de nuevo o cachea?").
