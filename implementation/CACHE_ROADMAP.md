# Roadmap de Implementación — ADR-039 (Caché y Optimización)

**Documento de soporte para ADR-039**  
**Estado:** Fase 3 completada (2026-07-05)  
**⚠️ Este documento puede cambiar; el ADR-039 es estable.**

---

## Fases (Completadas)

### Fase 1: Caché Persistente ✅

**Timeline:** Completada 2026-07-04  
**Entregables:**
- `tools/cache_manager.py` — gestor de caché con TTL, invalidación, manifest
- `.cache/cointracking/<proyecto>/` — estructura de almacenamiento
- Integración MCP: métodos `get_or_fetch()`, `invalidate_all()`, `invalidate_pattern()`

**Resultado:** Caché funcionando, almacenando trades/balances con versionado.

---

### Fase 2: Integración en Skills ✅

**Timeline:** Completada 2026-07-05  
**Entregables:**
- `/audit-cointracking` SKILL: Paso 0 con instrucciones CacheManager
- `/spanish-tax-return` SKILL: Paso 1 con reutilización de caché
- Documentación de uso en skills (ejemplos Python)

**Resultado:** Skills usan CacheManager automáticamente; caché transparente al usuario.

---

### Fase 3: Validación ✅

**Timeline:** Completada 2026-07-05  
**Entregables:**
- `tools/benchmark_skills.py` — test automatizado (3 auditorías + IRPF)
- `tools/test_cache_savings.py` — demostración de ahorro de caché
- `reports/SKILLS_BENCHMARK_REPORT.md` — informe de resultados en producción

**Resultado:** Validado en agp2025 (caso real, 1.670+ operaciones): 47-75% ahorro comprobado.

---

## Fases Futuras (Opcionales)

### Fase 4: Versionado de Caché

**Objetivo:** Detectar automáticamente si la caché es vigente.

**Tareas:**
- [ ] Extender manifest para incluir versiones:
  ```json
  {
    "cache_key": "get_trades_...",
    "generated_at": "2026-07-05T10:42:00Z",
    "adr_versions": {"ADR-039": "1.0"},
    "knowledge_versions": {"CAPITAL_GAINS": "2.1"},
    "mcp_version": "1.3.2"
  }
  ```
- [ ] Añadir método `is_valid(adr_versions, knowledge_versions)` a CacheManager
- [ ] Invalidar si algún ADR o KB cambió de versión

**Impacto:** Menor invalidación innecesaria; más confianza en caché.

---

### Fase 5: Niveles de Caché Dinámicos

**Objetivo:** TTL distintos por tipo de información.

**Tareas:**
- [ ] Definir estrategia de TTL por tipo:
  ```python
  TTL = {
    'trades': 'permanent_until_import_change',
    'balance': 15 * 60,  # 15 min
    'holdings': 15 * 60,
    'tax_report': 24 * 60 * 60,  # 24h
  }
  ```
- [ ] Implementar en CacheManager: `get_or_fetch(..., ttl_type='trades')`
- [ ] Auto-invalidación basada en cambios detectados (not just by age)

**Impacto:** Mejor balance entre frescura y ahorro.

---

### Fase 6: Dashboard de Caché

**Objetivo:** Visibilidad de ahorros de caché en CLI.

**Tareas:**
- [ ] Integrar `rtk gain --cache` → mostrar estadísticas:
  - Hits/misses
  - Tokens ahorrados por operación
  - Histórico de ahorro

**Impacto:** Usuario ve directamente el ROI de la optimización.

---

## Decisiones de Arquitectura (Referencia)

| Decisión | Implementado | Ubicación |
|----------|-------------|-----------|
| Caché persistente vs. en memoria | Persistente | Phase 1 ✅ |
| TTL por defecto | 24h (configurable) | Phase 1 ✅ |
| Invalidación por patrón | Sí (pattern matching) | Phase 1 ✅ |
| Procesamiento local | Sí (Python, no contexto LLM) | Phase 2 ✅ |
| Versionado | No aún (Phase 4) | TODO |
| TTL dinámico | No aún (Phase 5) | TODO |
| Never cache conclusions | Sí (principio, no técnica) | ADR-039 ✅ |

---

## Riesgos y Mitigación

### Riesgo: Caché desincronizada

**Síntoma:** Usuario ve datos viejos en auditoría.  
**Mitigación:** Phase 4 (versionado) detecta automáticamente.

### Riesgo: Hit rate bajo al inicio

**Síntoma:** Primeras auditorías sin ahorro.  
**Mitigación:** Aceptable; ahorro acumula en auditorías posteriores.

### Riesgo: Overhead de I/O

**Síntoma:** Lectura/escritura de disco más lenta que MCP.  
**Mitigación:** `.cache/` local, JSON compacto; benchmarks confirman OK.

### Riesgo: Crecimiento descontrolado de caché

**Síntoma:** `.cache/` ocupa GB tras meses de uso.  
**Mitigación:** Añadir limpieza automática (LRU) en Phase 4.

---

## Dependencias Externas

- [ ] Cambios en API MCP: revisar formatos de respuesta
- [ ] Cambios en ADR-037 (validación): actualizar manifest
- [ ] Cambios en knowledge base: revisar invalidación

---

## Métricas de Éxito

- ✅ Fase 1: CacheManager implementado y funcional
- ✅ Fase 2: Skills usando caché automáticamente
- ✅ Fase 3: 47%+ ahorro validado en producción
- ⏳ Fase 4: Caché versioned (en espera)
- ⏳ Fase 5: TTL dinámico (en espera)
- ⏳ Fase 6: Dashboard (en espera)

---

## Próximos Hitos

| Fecha | Hito | Estado |
|------|------|--------|
| 2026-07-05 | Fase 3 completada + validación | ✅ DONE |
| 2026-09-05 | Fase 4: Versionado (planificado) | ⏳ TODO |
| 2026-11-05 | Fase 5: TTL dinámico (planificado) | ⏳ TODO |
| 2026-12-31 | Revisión anual: Impacto real vs. estimado | ⏳ TODO |

---

**Este documento describe la cronología de implementación.**  
**Para los principios y decisiones arquitectónicas, ver ADR-039.**
