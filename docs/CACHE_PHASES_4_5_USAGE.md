# Uso de Caché Fase 4-5 en Skills

**Cómo usar versionado automático (Fase 4) y TTL dinámico (Fase 5)**

---

## Fase 4: Versionado Automático

### Problema que Resuelve

```
Escenario: Usuario ejecuta auditoría, luego ADR-039 cambia versión

Sin versionado:
  - Caché sigue usando datos con ADR-039 v1.0
  - Conclusiones basadas en reglas v1.0
  - Posible inconsistencia si v1.1 cambió lógica

Con versionado (Fase 4):
  - Se detecta: ADR-039 v1.0 → v1.1
  - Caché se invalida automáticamente
  - Se refetcha con v1.1 (conclusiones correctas)
```

### Cómo Usar

#### En Skills (Recomendado)

```python
from tools.cache_manager import CacheManager

mgr = CacheManager('agp2025')

# Usar método con versionado
trades = mgr.get_or_fetch_with_version_check(
    'get_trades',
    {'limit': None, 'start': unix_ts, 'end': unix_ts},
    mcp_call_fn=lambda call, params: mcp.cointracking_get_trades(**params),
    max_age_hours=24
)

# Si ADR-039 cambió versión:
#   [CACHE INVALIDATED] Cambios detectados:
#   - adr_0039: 1.0 → 1.1
# (automáticamente refetcha)
```

#### Verificar Manualmente si Caché es Válida

```python
mgr = CacheManager('agp2025')

# Versiones cuando se guardó el caché
cached_versions = {'adr_0039': '1.0', 'kb_capital_gains': '2.1'}

# Versiones actuales
current_versions = mgr.current_versions

# ¿Sigue siendo válida?
valid = mgr.is_cache_valid_by_version(cached_versions)

if not valid:
    diff = mgr.version_tracker.get_version_diff(cached_versions, current_versions)
    print(mgr.version_tracker.explain_invalidation(diff))
    # Invalidar y refetch
    mgr.invalidate_all()
```

### Qué Versiones Se Rastrean

La clase `VersionTracker` automáticamente extrae `version:` de:

1. **ADRs:** `adr/*.md` → `adr_0039`, `adr_0037`, etc.
2. **Knowledge Base:** `knowledge/**/*.md` → `kb_capital_gains`, `kb_cost_basis`, etc.
3. **MCP:** `.mcp.json` → `mcp` (si existe)

**Fronmatter YAML requerido:**
```yaml
---
version: 1.0
---
```

---

## Fase 5: TTL Dinámico

### Problema que Resuelve

```
Escenario: Usuario modifica un trade (venta), luego audita

Sin TTL dinámico (TTL = 24h):
  - Trades en caché 15 min después de cambio
  - Auditoría usa trades viejos
  - Ganancias incorrectas

Con TTL dinámico (Fase 5):
  - Trades: TTL permanente (hasta user_import)
  - Balance: TTL 15 min
  - get_gains: invalida si trades cambiaron
  - Auditoría siempre usa datos correctos
```

### Cómo Usar

#### En Skills (Simple)

```python
from tools.cache_ttl_manager import CacheTTLManager

mgr = CacheTTLManager('agp2025')

# En lugar de get_or_fetch(), usar get_or_fetch_dynamic()
# TTL se aplica automáticamente según tipo

trades = mgr.get_or_fetch_dynamic(
    'get_trades',
    {'limit': None},
    mcp_call_fn=lambda call, params: mcp.cointracking_get_trades(**params)
)
# → TTL permanente (hasta reimportación)

balance = mgr.get_or_fetch_dynamic(
    'get_grouped_balance',
    {},
    mcp_call_fn=lambda call, params: mcp.cointracking_get_grouped_balance(**params)
)
# → TTL 15 min (estado vivo)
```

#### Personalizar TTL

```python
# Si necesitas cambiar TTL para un caso específico
mgr.TTL_STRATEGIES['get_grouped_balance']['ttl_hours'] = 5  # 5h en lugar de 15 min

balance = mgr.get_or_fetch_dynamic(...)  # Usa 5h
```

#### Ver la Estrategia de TTL

```python
print(mgr.explain_ttl_strategy())
# Output:
# === Estrategia de TTL Dinámico ===
#
# get_trades........................ Permanente* (invalida por: user_import, version_change)
# get_grouped_balance............... 15 minutos (invalida por: user_operation, time_based)
# get_balance....................... 15 minutos (invalida por: user_operation, time_based)
# get_gains......................... Permanente* (invalida por: trades_change, version_change)
# ...
```

### TTL Predeterminados

| Llamada | TTL | Razón | Invalida Por |
|---------|-----|-------|--------------|
| `get_trades` | ♾️ permanente | Histórico (no cambia) | reimportación, version |
| `get_grouped_balance` | 15 min | Estado vivo | operación, tiempo |
| `get_balance` | 15 min | Estado vivo | operación, tiempo |
| `get_gains` | ♾️ permanente | Determinista (FIFO) | trade_change, version |
| `get_historical_summary` | 24h | Cambia diariamente | tiempo, version |
| `get_historical_currency` | 24h | Cambios EUR diarios | tiempo, version |
| `get_tax_report` | 24h | Reportes anuales | tiempo, version |

---

## Combinando Fase 4 + Fase 5

La forma **recomendada** en skills es combinar ambas:

```python
from tools.cache_ttl_manager import CacheTTLManager

class AuditSkill:
    def __init__(self, project_name: str):
        # Usa CacheTTLManager (que hereda de CacheManager)
        # Combina versionado + TTL dinámico automáticamente
        self.cache = CacheTTLManager(project_name)

    def reconcile(self):
        """Auditoría con caché inteligente."""

        # Trades: permanente (user_import) + versionado (ADRs)
        trades = self.cache.get_or_fetch_dynamic(
            'get_trades',
            {},
            mcp_call_fn=...
        )
        # Si usuario reimportó → invalida
        # Si ADR-037 cambió → invalida
        # Si TTL pasó → usa caché (permanente)

        # Balance: 15 min + versionado
        balance = self.cache.get_or_fetch_dynamic(
            'get_grouped_balance',
            {},
            mcp_call_fn=...
        )
        # Si 15 min pasaron → fetch nuevo
        # Si KB cambió → fetch nuevo
        # Si usuario hizo operación → fetch nuevo

        # Gains: permanente si trades no cambiaron
        gains = self.cache.get_or_fetch_dynamic(
            'get_gains',
            {},
            mcp_call_fn=...
        )
        # Si trades invalidos → invalida también
        # Si ADR-039 cambió → invalida

        return {
            'trades': trades,
            'balance': balance,
            'gains': gains
        }
```

---

## Flujo Completo: Ejemplo Real

### Escenario

```
1. Usuario: "Audita mi cuenta"
   → Se cachean datos (trades, gains, balance)
   → Se guardan versiones (adr_0039: 1.0, kb_capital_gains: 2.1)

2. Usuario: "Crea una transferencia manual en CoinTracking"
   → Datos de CoinTracking cambian
   → Usuario: "Audita de nuevo"

3. Sistema detecta:
   → Balance actualizado (15 min pasados) → fetch
   → Trades no cambiaron (sin reimportación) → reutiliza caché
   → get_gains invalida automáticamente (depende de trades) → reutiliza

4. Auditoría usa: balance nuevo + trades viejos (correctos) + gains viejos (correctas)
   → Resultado consistente sin llamadas MCP innecesarias
```

### Llamadas MCP Evitadas

```
Sin Fase 4-5:
  - Auditoría #1: 3 llamadas MCP (trades, balance, gains)
  - Balance actualizado en web
  - Auditoría #2: 3 llamadas MCP (todas de nuevo)
  - Total: 6 llamadas

Con Fase 4-5:
  - Auditoría #1: 3 llamadas MCP
  - Balance actualizado en web
  - Auditoría #2: 1 llamada MCP (balance) + 0 (trades/gains del caché)
  - Total: 4 llamadas (33% ahorro)
```

---

## Integración en Skills Existentes

### Cambio Mínimo (Recomendado)

En `/audit-cointracking` SKILL.md:

**Antes:**
```python
mgr = CacheManager(project_name)
balance = mgr.get_or_fetch('get_grouped_balance', {...}, ...)
```

**Después:**
```python
from tools.cache_ttl_manager import CacheTTLManager
mgr = CacheTTLManager(project_name)  # Misma interfaz, mejor inteligencia
balance = mgr.get_or_fetch_dynamic('get_grouped_balance', {...}, ...)  # TTL automático
```

**Beneficio:** Cero cambios en lógica, caché más inteligente automáticamente.

---

## Monitoreo

### Ver Estadísticas

```python
stats = mgr.stats()
print(f"Caché: {stats['total_entries']} entradas")
print(f"Tamaño: {stats['total_size_kb']} KB")
print(f"Versiones actuales: {stats['current_versions']}")
```

### Reportar Estrategia

```python
report = mgr.report_cache_strategy()
# Useful para documentación / auditoría
# → {'cache_strategy': 'Dynamic TTL per data type', 'strategies': {...}, ...}
```

---

**Fecha:** 2026-07-05  
**Versión:** Fase 4-5 Completadas  
**Próxima:** Fase 6 Dashboard (indefinido)
