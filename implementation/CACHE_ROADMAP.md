# Roadmap de Implementación — ADR-039 (Caché y Optimización)

**Documento de soporte para ADR-039**
**Estado:** ver tabla de fases abajo — corregido 2026-07-05 tras encontrar contradicciones internas (partes del documento decían "completada 2026-07-05" y otras "planificada 2026-09/11" para las mismas fases).
**⚠️ Este documento puede cambiar; el ADR-039 es estable.**

---

## Estado Real por Fase

| Fase | Código | ¿Conectado a las skills en producción? | Estado |
|------|--------|------------------------------------------|--------|
| **1. Caché persistente** | `tools/cache_manager.py` (`get_or_fetch`) | Base de la que heredan las demás fases | Completa |
| **2. Integración en skills** | `.claude/skills/*/SKILL.md` | ✅ Sí — actualizado 2026-07-05 a `CacheTTLManager.get_or_fetch_dynamic()` (antes usaban `CacheManager.get_or_fetch()` básico, sin versionado ni TTL dinámico) | Completa |
| **3. Validación** | `tools/benchmark_skills.py`, `tools/test_cache_savings.py` | N/A (son scripts de test) | Completa — 47-75% ahorro reproducible |
| **4. Versionado automático** | `tools/version_tracker.py`, `CacheManager.get_or_fetch_with_version_check()` | ✅ Sí, vía `CacheTTLManager.get_or_fetch_dynamic()` — bug de manifest y crash de encoding corregidos 2026-07-05 | Completa y conectada |
| **5. TTL dinámico** | `tools/cache_ttl_manager.py` (`CacheTTLManager.get_or_fetch_dynamic()`) | ✅ Sí — ambas skills lo usan desde 2026-07-05 | Completa y conectada |
| **6. Métricas/dashboard** | `tools/cache_metrics.py`, `tools/cache_cli.py` | ✅ Sí, automático dentro de `get_or_fetch_dynamic()` | Completa y conectada |

**Historial de la corrección (2026-07-05):** hasta esta fecha, las skills usaban `CacheManager.get_or_fetch()` (Fase 1) con `max_age_hours` fijo pasado a mano — las Fases 4-6 existían como código correcto pero standalone, nunca invocado desde el flujo real de auditoría. El CHANGELOG previo a esta corrección decía "Integrado en skills" para las Fases 4-6, lo cual era impreciso. Se corrigió cambiando el import y la llamada en ambos `SKILL.md` a `CacheTTLManager`/`get_or_fetch_dynamic()` (misma interfaz, cambio de bajo riesgo). Además se corrigió un bug real que habría hecho el versionado inoperante incluso conectado (ver `CHANGELOG.md`: acceso incorrecto al manifest + crash de encoding Unicode en Windows).

**Limitación que persiste, documentada honestamente:** `VersionTracker` solo detecta `version:` en frontmatter YAML. Los ADRs (`adr/*.md`, formato MADR plano) no tienen frontmatter, así que cambiar un ADR no invalida ningún caché — solo cambios en `knowledge/` lo hacen. No es un bug; es un límite del formato actual de los ADRs.

---

## Detalle de Cada Fase

### Fase 1: Caché Persistente ✅ (completa, en uso)

- `tools/cache_manager.py` — gestor de caché con TTL fijo, invalidación, manifest
- `.cache/cointracking/<proyecto>/` — estructura de almacenamiento
- Métodos: `get_or_fetch()`, `invalidate_all()`, `invalidate_pattern()`

### Fase 2: Integración en Skills ✅ (completa)

- `/audit-cointracking` SKILL: Paso 0 usa `CacheTTLManager.get_or_fetch_dynamic()`
- `/spanish-tax-return` SKILL: Paso 1 usa `CacheTTLManager.get_or_fetch_dynamic()`
- Hasta 2026-07-05 usaban `CacheManager.get_or_fetch()` básico (Fase 1 sin versionado ni TTL dinámico) — corregido, ver historial arriba

### Fase 3: Validación ✅ (completa)

- `tools/benchmark_skills.py` — test automatizado (3 auditorías + IRPF)
- `tools/test_cache_savings.py` — demostración de ahorro de caché
- `reports/SKILLS_BENCHMARK_REPORT.md` — informe de resultados en producción
- Validado en agp2025 (caso real, 1.670+ operaciones): 47-75% ahorro comprobado

### Fase 4: Versionado Automático — completa y conectada

- `tools/version_tracker.py` — rastreador de versiones, extrae `version:` de frontmatter YAML
  - **Limitación real:** solo funciona para `knowledge/` (tiene frontmatter YAML). Los ADRs (`adr/*.md`) usan formato MADR plano sin frontmatter, así que no son rastreables — cambiar un ADR no invalida caché con el formato actual. Cambiar un documento de `knowledge/` sí lo hace.
- `CacheManager.is_cache_valid_by_version()`, `get_or_fetch_with_version_check()` — funcionales tras el fix de 2026-07-05 (bug de acceso a manifest + crash de encoding en Windows, ambos corregidos; ver CHANGELOG)
- **Invocado por ambas skills** desde 2026-07-05, vía `CacheTTLManager.get_or_fetch_dynamic()`.

### Fase 5: TTL Dinámico — completa y conectada

- `tools/cache_ttl_manager.py` (`CacheTTLManager`) — extiende `CacheManager` con TTL por tipo de dato:
  ```python
  {
    'get_trades': {'ttl_hours': 999999, 'invalidate_on': ['user_import', 'version_change']},
    'get_grouped_balance': {'ttl_hours': 0.25, 'invalidate_on': ['user_operation']},  # 15 min
    'get_gains': {'ttl_hours': 999999, 'invalidate_on': ['trades_change']},
    'get_historical_summary': {'ttl_hours': 24},
  }
  ```
- Método `get_or_fetch_dynamic()` funcional (tras el fix de 2026-07-05, ahora también verifica versión antes de servir un hit con TTL permanente)
- **Invocado por ambas skills** desde 2026-07-05.

### Fase 6: Dashboard de Caché — completa y conectada

- `tools/cache_metrics.py` — rastreador de hits/misses/ahorros, persistido en `.cache/cointracking/<proyecto>/metrics.json`
- `tools/cache_cli.py` — CLI para reportes (`session`, `today`, `week`, `month`, `lifetime`, `detailed`)
- Se integra automáticamente **dentro de `CacheTTLManager.get_or_fetch_dynamic()`** — al usarlo ambas skills, las métricas ahora sí se registran en uso real, no solo en scripts de test/benchmark.

---

## Riesgos y Mitigación

### Riesgo: Caché desincronizada por cambio de conocimiento

**Síntoma:** Usuario corrige una regla en `knowledge/`, pero una auditoría con caché de `get_trades`/`get_gains` ya guardado sigue usando conclusiones calculadas con la regla vieja.
**Mitigación real hoy:** automática desde 2026-07-05 — al conectar las skills a `CacheTTLManager`, un cambio de versión en `knowledge/` invalida el caché aunque el TTL sea "permanente". Sigue sin cubrir cambios en `adr/` (ver limitación de Fase 4 arriba); para esos casos, invalidar caché manualmente (`cointracking_invalidate_cache`).

### Riesgo: Hit rate bajo al inicio

**Síntoma:** Primeras auditorías sin ahorro.
**Mitigación:** Aceptable; ahorro acumula en auditorías posteriores.

### Riesgo: Overhead de I/O

**Síntoma:** Lectura/escritura de disco más lenta que MCP.
**Mitigación:** `.cache/` local, JSON compacto; benchmarks confirman OK.

### Riesgo: Crecimiento descontrolado de caché

**Síntoma:** `.cache/` ocupa espacio creciente tras meses de uso.
**Mitigación:** pendiente — no hay limpieza automática (LRU) implementada.

---

## Dependencias Externas

- [ ] Cambios en API MCP: revisar formatos de respuesta
- [ ] Cambios en knowledge base: ya disparan invalidación automática (Fase 4 conectada); verificar tras cada edición de `version:` en frontmatter

---

**Este documento describe la cronología de implementación real.**
**Para los principios y decisiones arquitectónicas, ver ADR-039.**
