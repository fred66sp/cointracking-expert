# Guía de Optimización de Tokens (ADR-039)

**Implementación de caché y procesamiento local para reducir consumo de tokens en 80-90%.**

---

## Estado

- ✅ ADR-039 documentado (estrategia de 3 capas)
- ✅ `cache_manager.py` implementado (Fase 1)
- ⏳ Integración en skills (Fase 2)
- ⏳ Validación (Fase 3)

---

## Cómo Usar (Para Desarrolladores)

### Paso 1: Usar CacheManager en Skills

**En `/audit-cointracking` skill, reemplaza:**

```python
# ANTES (sin caché):
from mcp import cointracking
balance = cointracking.get_grouped_balance(project='agp')

# DESPUÉS (con caché):
from tools.cache_manager import CacheManager

mgr = CacheManager('agp')
balance = mgr.get_or_fetch(
    'get_grouped_balance',
    {'project': 'agp'},
    mcp_call_fn=lambda call, params: cointracking[call](**params),
    max_age_hours=24
)
```

### Paso 2: Procesar Localmente

**En skill, reemplaza:**

```python
# ANTES: Análisis en contexto LLM (~3000 tokens)
trades = mcp.get_trades()
# → Pasar trades al contexto para análisis

# DESPUÉS: Análisis local (~0 tokens)
trades = mgr.get_or_fetch('get_trades', {...}, mcp_call_fn=..., max_age_hours=1)
duplicates = ct_audit.detect_duplicates(trades)  # Python puro
orphans = ct_audit.detect_orphan_transfers(trades)
# → Pasar solo hallazgos al contexto
```

### Paso 3: Invalidar Cuando Sea Necesario

```python
# Usuario hizo cambios en CoinTracking
mgr.invalidate_pattern('get_trades')  # Limpia caché de trades
mgr.invalidate_all()  # Limpia TODO el caché del proyecto
```

---

## Cálculo de Ahorro

### Auditoría Típica

| Paso | Hoy | Optimizado | Ahorro |
|------|-----|-----------|--------|
| 1. get_trades() | 2000 tokens | 0 (caché) | - |
| 2. get_grouped_balance() | 500 tokens | 0 (caché) | - |
| 3. Análisis en contexto | 3000 tokens | 500 (local) | 83% |
| **Total** | **5500** | **500** | **91%** |

### Declaración Fiscal

| Paso | Hoy | Optimizado |
|------|-----|-----------|
| Auditoría | 5500 | 500 (reutiliza caché) |
| get_gains() | 1000 | 200 (caché + local) |
| Contexto IRPF | 5000 | 1000 (template + agregados) |
| **Total** | **11500** | **1700** | **85%** |

---

## Roadmap de Implementación

### Fase 1: CacheManager ✅

- `tools/cache_manager.py` implementado
- Funcionalidad básica: get_or_fetch, invalidate
- Testing manual: `python tools/cache_manager.py`

### Fase 2: Integración en Skills (Próxima)

**Timeline:** 1-2 horas  
**Tasks:**
1. Actualizar `/audit-cointracking` skill para usar CacheManager
2. Actualizar `/spanish-tax-return` skill
3. Documentar en skill docstring

### Fase 3: Validación (Próxima)

**Timeline:** 1 hora  
**Tasks:**
1. Ejecutar auditoría típica: medir tokens antes/después
2. Comparar con baseline (5500 → ~500)
3. Documentar resultados en ADR-039 como "Accepted"

---

## Archivos

| Archivo | Propósito |
|---------|-----------|
| `adr/0039-optimizacion-tokens-y-cache.md` | ADR con estrategia completa |
| `tools/cache_manager.py` | Implementación de CacheManager |
| `OPTIMIZATION_GUIDE.md` | Este archivo — guía de integración |

---

## Próximas Tareas

- [ ] Integrar CacheManager en `/audit-cointracking`
- [ ] Integrar CacheManager en `/spanish-tax-return`
- [ ] Validar ahorro (medir tokens)
- [ ] Documentar resultados en ADR-039
- [ ] Actualizar skills docstring con instrucciones

---

**Estado:** Fase 1 completa ✅ | Esperando integración en skills (Fase 2)
