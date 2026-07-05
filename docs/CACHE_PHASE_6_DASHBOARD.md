# Fase 6: Dashboard de Caché

**Visualización de ahorros en tiempo real**

**Status:** ✅ COMPLETADA (2026-07-05)

---

## Qué Es Fase 6

Sistema de rastreamiento y reporteo de ahorros de caché. Muestra en tiempo real cuántos tokens se está ahorrando.

---

## Componentes

### 1. CacheMetrics (`tools/cache_metrics.py`)

Rastreador automático de hits/misses:

```python
from tools.cache_metrics import CacheMetrics

metrics = CacheMetrics('agp2025')

# Automáticamente registrado por CacheTTLManager:
metrics.record_cache_hit('get_trades', tokens_saved=2835, age_hours=2.5)
metrics.record_mcp_call('get_historical_summary', tokens_cost=400)

# Obtener estadísticas
stats = metrics.get_stats('session')
# → {'hits': 3, 'misses': 1, 'tokens_saved': 4335, ...}

# Mostrar reporte
print(metrics.get_report('lifetime'))
```

### 2. CacheTTLManager Integrado

Automáticamente registra en métricas:

```python
from tools.cache_ttl_manager import CacheTTLManager

mgr = CacheTTLManager('agp2025')  # ← Ya incluye metrics

# Uso normal
trades = mgr.get_or_fetch_dynamic('get_trades', {}, mcp_call_fn=...)

# Internamente:
#   Si CACHE_HIT → metrics.record_cache_hit(...) 
#   Si MCP_MISS → metrics.record_mcp_call(...)
```

### 3. CLI (`tools/cache_cli.py`)

Mostrar reportes desde terminal:

```bash
python tools/cache_cli.py agp2025 session
# → Estadísticas de esta sesión

python tools/cache_cli.py agp2025 lifetime
# → Estadísticas totales desde inicio

python tools/cache_cli.py agp2025 detailed
# → Reporte detallado con desglose por call
```

---

## Salida Típica

```
======================================================================
ESTADISTICAS DE CACHE - SESSION
======================================================================

Operaciones:
  Hits (caché reutilizado):      3
  Misses (MCP llamado):          1
  Total:                         4
  Hit Rate:                  75.0%

Tokens:
  Ahorrados (hits):            4335
  Gastados (misses):            400
  Neto:                        3935
  Ahorro %:                  91.6%

======================================================================
```

---

## Períodos Disponibles

| Período | Descripción | Uso |
|---------|-------------|-----|
| `session` | Esta sesión (desde arranque) | Verificar ahorro inmediato |
| `today` | Hoy (desde las 00:00) | Trending diario |
| `week` | Esta semana | Tendencia semanal |
| `month` | Este mes | Tendencia mensual |
| `lifetime` | Total histórico | Ahorro acumulado |
| `detailed` | Con desglose por call | Saber dónde se ahorra más |

---

## Métricas Rastreadas

### Por Operación

```json
{
  "call": "get_trades",
  "type": "CACHE_HIT",
  "tokens_saved": 2835,
  "age_hours": 2.5,
  "timestamp": "2026-07-05T10:42:00Z"
}
```

### Agregadas

- **Hits:** Operaciones reutilizadas del caché
- **Misses:** Operaciones que llamaron MCP
- **Tokens_saved:** Total de tokens ahorrados (hits × estimado)
- **Tokens_cost:** Total de tokens gastados (misses × estimado)
- **Hit_rate:** Porcentaje de hits

---

## Estimaciones de Tokens

Basadas en `docs/performance/TOKEN_BENCHMARKS.md`:

| Llamada | Tokens |
|---------|--------|
| get_trades | 2.835 |
| get_grouped_balance | 500 |
| get_balance | 300 |
| get_gains | 1.000 |
| get_historical_summary | 400 |
| get_historical_currency | 400 |
| get_tax_report | 800 |

---

## Almacenamiento

Persiste en:

```
.cache/cointracking/<proyecto>/
├── manifest.json
├── metrics.json  ← Datos de caché (hits, misses, ahorros)
└── ...
```

Formato:

```json
{
  "project": "agp2025",
  "created": "2026-07-05T09:00:00Z",
  "session": {
    "operations": [
      { "call": "get_trades", "type": "CACHE_HIT", ... },
      { "call": "get_grouped_balance", "type": "CACHE_HIT", ... }
    ]
  },
  "daily": { "2026-07-05": { "hits": 2, "misses": 0, ... } },
  "weekly": { "2026-07-01": { "hits": 12, "misses": 3, ... } },
  "monthly": { "2026-07": { "hits": 50, "misses": 8, ... } },
  "lifetime": { "hits": 127, "misses": 23, "tokens_saved": 445000, ... }
}
```

---

## Integración con Skills

En `/audit-cointracking` o `/spanish-tax-return`:

```python
mgr = CacheTTLManager(project_name)
# ← Ya rastrean automáticamente

# Luego de ejecutar skill:
print(mgr.metrics.get_report('session'))
# Usuario ve: "He ahorrado 4.335 tokens en esta auditoría"
```

---

## Ejemplo: Flujo Completo

```
1. Usuario: /audit-cointracking agp2025
   ├─ Descarga trades (MCP MISS, 2835 tokens)
   ├─ Descarga balance (MCP MISS, 500 tokens)
   ├─ Descarga gains (CACHE HIT, reutiliza 1000 tokens ahorrados)
   └─ Descarga historical (CACHE HIT, reutiliza 400 tokens ahorrados)

2. Sistema registra en metrics:
   ├─ get_trades: MISS (2835 tokens)
   ├─ get_grouped_balance: MISS (500 tokens)
   ├─ get_gains: HIT (1000 tokens ahorrados)
   └─ get_historical_summary: HIT (400 tokens ahorrados)

3. Usuario: python tools/cache_cli.py agp2025 session
   └─ [RESULTADO]
   Hits: 2, Misses: 2, Hit Rate: 50%, Ahorrados: 1400, Gastados: 3335, Neto: -1935

4. Usuario: python tools/cache_cli.py agp2025 lifetime
   └─ [RESULTADO después de varias auditorías]
   Hits: 127, Hit Rate: 84%, Ahorrados: 445000, Gastados: 86000, Neto: +359000
```

---

## Beneficios

✅ **Visualización:** Usuario ve exactamente cuánto ahorra  
✅ **ROI Comprobable:** "He ahorrado 445K tokens este mes"  
✅ **Optimización Basada en Datos:** Saber dónde se ahorra más  
✅ **Feedback Positivo:** "La caché funciona"  
✅ **Historial:** Trending de ahorros a lo largo del tiempo

---

## Roadmap

- ✅ Fase 6a: Rastreador básico (CacheMetrics)
- ✅ Fase 6b: Integración automática (CacheTTLManager)
- ✅ Fase 6c: CLI (cache_cli.py)
- ⏳ Fase 6d (futuro): Dashboard web / `rtk gain --cache` integrado
- ⏳ Fase 6e (futuro): Alertas si hit rate baja

---

**Fecha:** 2026-07-05  
**Versión:** Fase 6 Completada  
**Próximo:** Dashboard web o integración RTK
