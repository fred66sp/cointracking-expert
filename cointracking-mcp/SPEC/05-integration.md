# 05 - Integration

## Cómo Funciona con el Agente

El agente (auditor de Claude Code) llama al MCP server a través del protocolo estándar de Model Context Protocol.

### Flujo Típico de Auditoría

```
1. Usuario: "audita mi cuenta de CoinTracking"
2. Agente (/audit-cointracking): 
   - Obtiene credenciales del usuario o MCP env
   - Llama cointracking_get_trades(limit=500, start=..., end=...)
3. MCP Server:
   - Valida parámetros
   - Busca en caché
   - (Si miss) Llama API de CT
   - Retorna JSON de trades
4. Agente:
   - Recibe 500 trades
   - Parsea JSON
   - Ejecuta lógica de auditoría (búsqueda de inconsistencias)
   - Reporta hallazgos al usuario
5. Usuario (si hay problemas):
   - Modifica datos en web de CoinTracking
   - Vuelve al agente: "ya lo arreglé"
6. Agente:
   - Llama cointracking_invalidate_cache("getTrades*")
   - Re-ejecuta auditoría
   - Confirma que se arregló
```

## Configuración en Claude Code

### Archivo: `.claude/settings.json` del Agente

```json
{
  "mcp_servers": {
    "cointracking": {
      "command": "H:\\cointracking-mcp\\dist\\cointracking-mcp",
      "env": {
        "COINTRACKING_API_KEY": "${COINTRACKING_API_KEY}",
        "COINTRACKING_API_SECRET": "${COINTRACKING_API_SECRET}"
      }
    }
  }
}
```

**Notas:**
- Las credenciales vienen de variables de entorno (sin `.env`)
- El path del binario debe ser absoluto
- El agente hereda las variables de entorno del harness de Claude Code
- Todos los parámetros opcionales tienen defaults sensatos

### Variables de Entorno

Las credenciales se pasan como variables de entorno (sin `.env`):

```bash
export COINTRACKING_API_KEY=your_api_key_here
export COINTRACKING_API_SECRET=your_api_secret_here
./cointracking-mcp
```

Parámetros opcionales (si quieres desviarte de defaults):
```bash
export PROJECT_NAME=agp
export ACCOUNT_TIER=unlimited
export CACHE_PERSIST_DIR=./cache
export LOG_LEVEL=debug
./cointracking-mcp
```

## Protocolo MCP (Detalles Técnicos)

### Request → Response

El agente envía (stdin):
```json
{
  "jsonrpc": "2.0",
  "id": "req-123",
  "method": "tools/call",
  "params": {
    "name": "cointracking_get_trades",
    "arguments": {
      "limit": 100,
      "start": 1234567890,
      "end": 1234654290
    }
  }
}
```

El server responde (stdout):
```json
{
  "jsonrpc": "2.0",
  "id": "req-123",
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{...JSON de trades...}"
      }
    ]
  }
}
```

**Error:**
```json
{
  "jsonrpc": "2.0",
  "id": "req-123",
  "error": {
    "code": -32603,
    "message": "Internal error",
    "data": {
      "code": "RATE_LIMIT",
      "message": "Rate limit excedido (configurable por tipo de cuenta)"
    }
  }
}
```

### Herramientas Disponibles (Tool Definitions)

El server publica estas herramientas al startup:

```json
{
  "name": "cointracking_get_trades",
  "description": "Returns all trades and transactions from CoinTracking...",
  "inputSchema": {
    "type": "object",
    "properties": {
      "limit": {
        "type": "integer",
        "description": "Maximum number of trades to return..."
      },
      "start": {
        "type": "integer",
        "description": "Start time as UNIX timestamp in SECONDS..."
      },
      ...
    }
  }
}
```

## Decisiones de Integración

### 1. ¿Dónde Vive el MCP?

**Opción A (Recomendada):** Repo separado
- Ciclo de vida independiente
- Reutilizable en otros contextos
- Versionado por separado

**Opción B:** Subdir del agente (`cointracking-expert/tools/cointracking-mcp`)
- Todo junto
- Más fácil para desarrollo local
- Menos limpio arquitectónicamente

**Decisión:** Opción A (como indica el usuario).

### 2. ¿Caché En Memoria o Persistido?

**Fase 0:** Solo en memoria (rápido, simple)
**Fase 1:** Opcional persistencia a SQLite

Razón: caché en memoria es suficiente durante una sesión de auditoría (típicamente minutos), y reiniciar el server es seguro (solo re-consulta a la API si es necesario).

### 3. ¿Cómo se Pasan Credenciales?

**Variables de entorno** + `.env` (estándar de MCP)
- `COINTRACKING_API_KEY`
- `COINTRACKING_API_SECRET`

No hardcodear en código.

### 4. ¿El Server Hace Validaciones Fiscales?

**No.** El server:
- Proporciona datos limpiable
- Detecta inconsistencias técnicas (duplicados, saldos imposibles)

El agente auditor:
- Interpreta los hallazgos
- Explica impacto fiscal
- Guía correcciones

### 5. ¿Reintento en Rate Limit?

**No automático.** El server reporta 429 (rate limit) y el agente decide:
- Esperar y reintentar
- Usar datos en caché más antiguos
- Acotar la consulta

## Testing de Integración

Probar con fixtures:

```bash
# 1. Fixture de trades (sample_trades.json)
# 2. Fixture de balance (sample_balance.json)
# 3. Lanzar el server en modo test (sin credenciales reales)
# 4. Agente hace llamadas
# 5. Verificar resultados contra expected.json
```

Ver `tests/fixtures/` para más.

## Límites Conocidos

1. **Solo lectura** — API de CT no permite creación/edición vía API
2. **Rate limit según plan** — 20 (PRO/EXPERT, default) o 60 (UNLIMITED); sin reintento automático
   - Por defecto seguro: 20 llamadas/hora
   - Si recibe 429, informa al agente para que ajuste `ACCOUNT_TIER` si es necesario
3. **Caché en memoria** — si el server crashea, se pierde (datos en CT permanecen)
4. **Sin sincronización multi-instancia** — si el usuario lanza varios servers con el mismo API key, el caché no se sincroniza (futura mejora)

## Roadmap de Integración

- [ ] **Fase 0:** Servidor base + 6 tools + caché en memoria
- [ ] **Prueba manual:** Agente llama server, obtiene datos
- [ ] **Fixtures y tests:** sample_trades.csv, validación determinista
- [ ] **Fase 1:** Caché persistido, validaciones, agregados
- [ ] **Fase 2:** Detección de anomalías, UI de stats
