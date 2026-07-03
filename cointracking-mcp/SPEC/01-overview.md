# 01 - Overview

## Propósito

Servidor MCP que actúa como **puente entre la API de CoinTracking y el agente auditor de Claude Code**, proporcionando:

1. **Interfaz uniforme** a las herramientas de CoinTracking (trades, balances, históricos, ganancias)
2. **Caché inteligente, por proyecto** (LRU + TTL; reduce llamadas a la API; respeta límite de rate limit configurado)
3. **Multi-proyecto** (audita múltiples cuentas con cachés aisladas; útil para diferentes clientes/cuentas)
4. **Validaciones deterministas** (detección de anomalías, inconsistencias, cálculos verificables)
5. **Trazabilidad** (log de todas las llamadas, con timestamps y parámetros)

## Arquitectura

```
┌─────────────────────────────────────────┐
│   Agente Auditor (Claude Code)          │
│   /audit-cointracking                   │
│   /spanish-tax-return                   │
└──────────────────┬──────────────────────┘
                   │ MCP (stdio JSON-RPC)
                   ▼
┌─────────────────────────────────────────┐
│   cointracking-mcp (Go Server)          │
│                                         │
│  ┌─────────────────────────────────┐   │
│  │  Tool Handlers                  │   │
│  │  - getTrades                    │   │
│  │  - getBalance                   │   │
│  │  - getHistoricalSummary         │   │
│  │  - getHistoricalCurrency        │   │
│  │  - getGroupedBalance            │   │
│  │  - getGains                     │   │
│  └─────────────────────────────────┘   │
│           │                             │
│           ▼                             │
│  ┌─────────────────────────────────┐   │
│  │  Cache Layer (in-process)       │   │
│  │  - LRU + TTL                    │   │
│  │  - Invalidación manual          │   │
│  │  - Persistencia opcional a disk │   │
│  └─────────────────────────────────┘   │
│           │                             │
│           ▼                             │
│  ┌─────────────────────────────────┐   │
│  │  API Client (HMAC-SHA512)       │   │
│  │  - Rate limit tracking          │   │
│  │  - Nonce monotónico             │   │
│  │  - Error handling               │   │
│  └─────────────────────────────────┘   │
└──────────────────┬──────────────────────┘
                   │ HTTPS
                   ▼
        ┌──────────────────────┐
        │ CoinTracking API v1  │
        │ https://cointracking │
        │   .info/api/v1/      │
        └──────────────────────┘
```

## Comportamiento

### Startup
1. Leer credenciales (variables de entorno: `COINTRACKING_API_KEY`, `COINTRACKING_API_SECRET`)
2. Inicializar caché en memoria (LRU, configurable en tamaño)
3. Optionally: cargar caché persistida de disco (si existe)
4. Escuchar en stdin con protocolo MCP

### Llamada Típica
1. Agente pide una herramienta (ej: `getTrades`, parámetros)
2. Server busca en caché con clave = (método + parámetros normalizados)
3. Si hit + no expirado → responder desde caché
4. Si miss o expirado → llamar a API, cachear resultado, responder
5. Log de la llamada (fuente, parámetros, tiempo, caché/API)

### Invalidación
- Manual: agente envía una herramienta especial `invalidateCache(pattern)` cuando el usuario modifica datos en CoinTracking
- TTL: cada entrada en caché expira después de cierto tiempo (configurable por tipo)

## Configuración

**CLI flags para TODO + envvars solo para credenciales**

**Mínimo (CLI flags):**
```bash
./cointracking-mcp --api-key xxx --api-secret yyy
```

**Con más opciones:**
```bash
./cointracking-mcp \
  --api-key xxx \
  --api-secret yyy \
  --project agp \
  --tier unlimited \
  --log-level debug
```

**O credenciales en envvars:**
```bash
# Bash/Linux/macOS
export COINTRACKING_API_KEY=xxx
export COINTRACKING_API_SECRET=yyy
./cointracking-mcp --project agp --tier unlimited

# PowerShell (Windows)
$env:COINTRACKING_API_KEY = "xxx"
$env:COINTRACKING_API_SECRET = "yyy"
./cointracking-mcp --project agp --tier unlimited

# CMD.exe (Windows)
set COINTRACKING_API_KEY=xxx
set COINTRACKING_API_SECRET=yyy
cointracking-mcp.exe --project agp --tier unlimited
```

✅ Listo. El rest usa **defaults sensatos**:
- **Proyecto:** `default` (si quieres otro: `PROJECT_NAME=agp`)
- **Caché:** 5000–150000 entradas por proyecto según tier (L1 memoria + L2 disco, TTL inteligentes)
  - Pro: 5000 entradas (cubre 3.500 registros)
  - Expert: 50000 entradas (cubre 20–100k registros)
  - Unlimited: 150000 entradas (sin límite)
  - Búsqueda: primero memoria → luego disco → si no, API
- **Rate limit:** 20 llamadas/hora (PRO; cambia a `ACCOUNT_TIER=unlimited` si tienes ese plan)
- **Persistencia:** **siempre en disco** (`./cache/{PROJECT_NAME}/`) + automática
- **Cierre:** herramienta `cointracking_close_project` para vaciar memoria y sincronizar
- **Log:** info level

**Para customizar** (opcional): ver [06 - Configuration](06-configuration.md).

## Invocación desde Claude Code

1. En settings.json del agente:
```json
{
  "mcp_servers": {
    "cointracking": {
      "command": "/path/to/cointracking-mcp",
      "env": {
        "COINTRACKING_API_KEY": "${COINTRACKING_API_KEY}",
        "COINTRACKING_API_SECRET": "${COINTRACKING_API_SECRET}"
      }
    }
  }
}
```

2. El agente simplemente llama a las herramientas como de costumbre (no sabe que hay caché; es transparente).

## Limitaciones Conocidas (Fase 0)

- Solo lectura (la API de CT no permite modificación por API)
- No sincronización con múltiples instancias (memoria local)
- No reintento automático por rate limit (agente debe manejar)
