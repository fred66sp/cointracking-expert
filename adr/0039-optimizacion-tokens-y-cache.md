# ADR-039: Optimización de Tokens y Estrategia de Caché

**Status:** Proposed  
**Proposed:** 2026-07-05  
**Accepted:** -  
**Last Updated:** 2026-07-05

---

## Problema

El sistema usa MCP (servidor Go) para acceder a CoinTracking, pero:
- Límite: 60 llamadas/hora
- Contexto: Cada respuesta JSON grande consume tokens
- Ineficiencia: Llamadas redundantes sin caché persistente
- No hay estrategia clara de optimización de tokens

**Costo típico hoy:**
- `cointracking_get_trades()` completo: ~2000 tokens
- `cointracking_get_grouped_balance()`: ~500 tokens
- Auditoría MEGA completa: ~5-10K tokens (depende de tamaño)
- Preparar IRPF: ~15-20K tokens (con contexto de conocimiento)

---

## Decisión

**Estrategia de tokens de tres capas:**

### Capa 1: Caché Local (MEJOR AHORRO)

**Implementar caché persistente** en `.cache/cointracking/`:

```
.cache/cointracking/
├── {project}/
│   ├── balance_2026-07-05_12-30.json     (última llamada a get_balance)
│   ├── grouped_balance_2026-07-05.json   (por día)
│   ├── trades_2026-01-01_2026-07-05.json (rango completo)
│   ├── gains_fifo_2026.json              (por año fiscal)
│   └── cache_manifest.json               (metadata)
```

**Reglas:**
- ✅ Reutilizar caché si:
  - Menos de 24h de antigüedad (balances)
  - Mismo rango de fechas (trades, gains)
  - Usuario no pide "refresco"
- ❌ Invalidar caché si:
  - Usuario hizo cambios en CoinTracking (guiar correcciones)
  - Cambio de proyecto activo
  - Usuario ejecuta `cointracking_invalidate_cache`

**Ahorro estimado:** 50-70% de llamadas MCP

### Capa 2: Agregados (No Detalles)

**Preferir llamadas agregadas sobre detalles:**

| Costoso | Económico | Ahorro |
|---------|-----------|--------|
| `get_trades()` completo | `get_grouped_balance()` | 60% |
| Listar 500 operaciones | Contar por categoría | 70% |
| Descargar historial año | `get_gains()` por ejercicio | 80% |

**Regla:** Si solo necesitas cifras, usa agregados.

### Capa 3: Procesamiento Local

**Procesar datos localmente, no en contexto LLM:**

```python
# ❌ COSTOSO: Pasar 5000 trades al LLM para análisis
trades_full = mcp.get_trades(limit=5000)  # ~2000 tokens
# → Luego en contexto: summarize(trades_full)  # +3000 tokens

# ✅ ECONÓMICO: Procesar en Python, pasar solo resultado
trades_full = mcp.get_trades(limit=5000)  # ~2000 tokens
summary = analyze_trades_locally(trades_full)  # Python puro
# → Pasar solo summary al LLM (~200 tokens)
```

**Herramientas:** `tools/ct_audit.py`, scripts de análisis en Python.

---

## Implicaciones

### Para `/audit-cointracking` Skill

```
Hoy (sin caché):
  1. get_trades() → 2000 tokens
  2. get_grouped_balance() → 500 tokens
  3. Análisis en contexto → +3000 tokens
  Total: ~5500 tokens/auditoría

Con optimizaciones:
  1. get_trades() (caché) → reutiliza (0 tokens)
  2. get_grouped_balance() (caché) → reutiliza (0 tokens)
  3. Análisis local (Python) → 0 tokens
  4. Pasar solo hallazgos → 500 tokens
  Total: ~500 tokens/auditoría (90% ahorro)
```

### Para `/spanish-tax-return` Skill

```
Hoy:
  1. Auditoría (reutiliza datos) → 500 tokens
  2. get_gains() → 1000 tokens
  3. Preparar IRPF en contexto → +5000 tokens
  Total: ~6500 tokens/declaración

Con optimizaciones:
  1. Auditoría (caché) → 100 tokens
  2. get_gains() (caché + local process) → 200 tokens
  3. Preparar IRPF con template → 1000 tokens
  Total: ~1300 tokens/declaración (80% ahorro)
```

---

## Implementación (Roadmap)

### Fase 1: Caché Persistente (Week 1)

```python
# tools/cache_manager.py (crear)
class CacheManager:
    def get_or_fetch(self, call_name, params, max_age_hours=24):
        cached = self.load_from_disk(call_name, params)
        if cached and cached['age'] < max_age_hours:
            return cached['data']
        # Si no existe o está viejo, fetch y cache
        data = mcp_call(call_name, params)
        self.save_to_disk(call_name, params, data)
        return data

    def invalidate(self, project):
        # Borrar caché del proyecto
        shutil.rmtree(f'.cache/cointracking/{project}/')
```

**En skills:** Reemplazar llamadas MCP con `cache_mgr.get_or_fetch()`

### Fase 2: Procesamiento Local (Week 2)

```python
# tools/analysis_local.py (crear/mejorar)
def analyze_trades_for_audit(trades):
    """Análisis sin contexto LLM, devuelve hallazgos compactos"""
    duplicates = detect_duplicates(trades)
    orphans = detect_orphan_transfers(trades)
    negative_balances = check_negative_balances(trades)
    return {
        'duplicates': duplicates,
        'orphans': orphans,
        'negative_balances': negative_balances,
        'summary': f"{len(duplicates)} duplicados, {len(orphans)} huérfanas..."
    }
```

**En skill:** Pasar solo hallazgos al contexto.

### Fase 3: Validación (Week 3)

- Medir tokens antes/después
- Comparar con baseline (5500 → ~500 tokens)
- Documentar en ADR-039 como "Accepted"

---

## Estimación de Ahorro

| Componente | Hoy | Optimizado | Ahorro |
|-----------|-----|-----------|--------|
| Auditoría | 5500 | 500 | 91% |
| IRPF | 6500 | 1300 | 80% |
| Auditoría MEGA | 5000 | 1000 | 80% |
| Promedio | 5667 | 933 | **84% ahorro total** |

**Impacto:** Con 50 operaciones/mes, pasa de ~283K tokens a ~47K tokens.

---

## Riesgos / Consideraciones

| Riesgo | Mitigation |
|--------|-----------|
| Caché desincronizado | Invalidar al detectar cambios usuario |
| Hit rate bajo al inicio | Será alto después de 2-3 usos |
| Análisis local incompleto | Usar MCP solo si análisis local reporta ambigüedad |
| Overhead de I/O disco | Mínimo en `.cache/` local (< 1ms) |

---

## Referencias

- ADR-010: Eficiencia de tokens y caché de CoinTracking
- `tools/ct_audit.py`: Ya procesa localmente, no en contexto
- `tools/cache_manager.py`: A crear (Fase 1)

---

**Decisión:** PROPUESTA (requiere implementación en 3 fases)
