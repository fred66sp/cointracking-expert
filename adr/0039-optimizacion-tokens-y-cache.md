# ADR-039: Optimización de Recursos y Estrategia de Caché

**Status:** Accepted  
**Proposed:** 2026-07-05  
**Accepted:** 2026-07-05  
**Last Updated:** 2026-07-05

**Clasificación:** Arquitectónico (afecta a todo el sistema: rendimiento, coste, escalabilidad)

---

## Contexto

El sistema audita operaciones cripto usando dos recursos finitos:

1. **Llamadas MCP:** Límite 60/hora al servidor de CoinTracking
2. **Tokens de contexto LLM:** Recursivo del modelo; costo operativo

Sin estrategia de caché y procesamiento local, ambos se consumen innecesariamente en datos ya disponibles o análisis que no requieren LLM.

---

## Problema

El sistema incurre en **costos redundantes** sin obtener mejor calidad:

- Llamadas MCP duplicadas (mismo proyecto, mismo rango de fechas)
- JSON crudo en contexto LLM (5.000 transacciones cuando bastaban 10 hallazgos)
- Recompilación de análisis idénticos tras cambios menores
- Escalabilidad comprometida (50+ proyectos/año → token budget insostenible)

**Principio violado:** El contexto LLM debe contener el mínimo de información necesario para tomar decisiones trazables.

---

## Decisión

**Optimizar consumo de recursos mediante caché distribuida y procesamiento local, sin comprometer trazabilidad ni reproducibilidad.**

La optimización se implementa en **tres capas**:

### Capa 1: Caché Persistente

**¿Qué?** Almacenar respuestas MCP en disco con versioning.

**¿Por qué?** Las respuestas MCP no cambian a menos que el usuario importe de nuevo. Reutilizarlas evita llamadas redundantes.

**Regla de invalidación:**
```
Invalidar si:
  - Usuario importó datos nuevos en CoinTracking
  - Cambió proyecto activo
  - ADR relacionados con auditoría cambiaron de versión
  - Knowledge base de reglas fiscales se actualizó
  - Versionado de caché reporta obsolescencia
  
No invalidar por:
  - Paso del tiempo (si los datos no cambiaron en origen)
  - Cambio en other projects (cada proyecto tiene caché aislada)
```

**Niveles de TTL (configurables):**

| Tipo | TTL | Justificación |
|------|-----|---------------|
| Trades (histórico) | Permanente* | No cambia salvo reimportación |
| Holdings | 15 min | Pueden cambiar entre operaciones |
| Balance actual | 15 min | Estado vivo del exchange |
| Tax Report anual | 24h | Generado una vez/año |
| Gains (FIFO) | Permanente* | Determinista si trades no cambian |
| Knowledge base | Versionado | Cambios ex plícitos |

*Permanente hasta detectar cambios en origen.

### Capa 2: Agregados antes que Detalle

**¿Qué?** Preferir `get_grouped_balance()` sobre `get_trades()` cuando sea posible.

**¿Por qué?** Agregados devuelven el mismo contexto en 20% del tamaño.

**Regla:**
- Si necesitas **cifras**: usa agregados.
- Si necesitas **estructura/justificación de cifras**: usa detalle (ya en caché).

### Capa 3: Procesamiento Local sin Contexto LLM

**¿Qué?** Delegar análisis de datos a scripts Python (no a LLM).

**¿Por qué?** El LLM no necesita procesar 5.000 transacciones para concluir:
```
- 3 balances negativos
- 2 posibles duplicados  
- 1 transferencia huérfana
```

Basta con recibir esos hallazgos resumidos.

**Regla de oro:**
```
LLM recibe:

NUNCA:   JSON crudo, listas completas, datos sin filtrar
SIEMPRE: Evidencias resumidas, métricas, anomalías, resultados intermedios reproducibles
```

**Arquitectura:**
```
Datos origen (MCP/CSV)
    ↓
Análisis local (Python, determinista)
    ↓
Hallazgos compactos (< 500 tokens)
    ↓
Contexto LLM (interpretación, explicación)
```

---

## Principios Arquitectónicos

### Principio 1: Integridad de Auditoría

> **La optimización nunca debe alterar el resultado de una auditoría; únicamente reducir el coste computacional necesario para obtenerlo.**

Contraejemplo: No omitir una transferencia huérfana solo porque su detección se optimizó.

### Principio 2: No Cachear Conclusiones

> **Solo se cachean datos y resultados intermedios reproducibles. Nunca conclusiones.**

- ✅ Cachear: trades (datos), gains (resultado FIFO determinista), balance (hecho, no interpretación)
- ❌ Cachear: "es una venta imponible" (conclusión depende de ADRs), "usuario tiene riesgo fiscal" (depende de escenarios)

Las conclusiones dependen de:
- ADRs vigentes (cambian)
- Reglas fiscales (actualizan cada año)
- Versión del agente (evoluciona)
- Algoritmos de reconciliación (se mejoran)

Si algo de esto cambia, las conclusiones previas invalidan.

### Principio 3: Minimización de Contexto

> **El agente solo incorporará al contexto del LLM la cantidad mínima de información necesaria para tomar una decisión trazable, reproducible y explicable.**

Este principio es transversal: afecta a caché, a formato de respuestas, a templates, a todo.

**Aplicación:**
- No incluyas 500 transacciones en contexto; incluye el resumen.
- No pases JSON crudo; pasa tablas procesadas.
- No describa cada operación; describe hallazgos y solicita detalle si es necesario.

---

## Implicaciones Operativas

### Auditoría

Sin optimizar:
- MCP: 5.500 tokens (get_trades, get_balance, get_gains, get_historical_summary)
- Contexto: 3.000 tokens (análisis en LLM)
- **Total: ~8.500 tokens**

Con optimizar (run 1):
- MCP: 5.500 tokens (datos nuevos, necesarios)
- Contexto: 200 tokens (análisis local, hallazgos compactos)
- **Total: ~5.700 tokens** (33% ahorro)

Con caché (run 2+):
- MCP: 0 tokens (CACHE HIT 100%)
- Contexto: 200 tokens
- **Total: ~200 tokens** (98% ahorro)

### Declaración Fiscal

Sin optimizar:
- MCP: 1.000 tokens (get_gains nuevo)
- Contexto: 2.200 tokens (preparar informe)
- **Total: ~3.200 tokens**

Con optimizar:
- MCP: 1.000 tokens (nuevo, necesario)
- Contexto: 300 tokens (template + resultados previos)
- **Total: ~1.300 tokens** (59% ahorro)

---

## Versionado y Reproducibilidad

La caché registra el contexto con el que fue generada:

```json
{
  "project": "agp2025",
  "timestamp": "2026-07-05T10:42:00Z",
  "data_versions": {
    "adr_039": "1.0",
    "capital_gains_rules": "2025-edition",
    "ct_import_format": "2026-q2",
    "knowledge_base": "v2.3.1"
  },
  "validity": "valid_if_versions_match"
}
```

Esto permite:
- Detectar si una auditoría previa sigue siendo válida
- Regenerar con versiones nuevas si algo cambió
- Auditar el linaje de datos (de dónde salió cada cifra)

---

## Criterios de Invalidación (Completos)

La caché se invalida (parcial o total) si:

1. **Usuario hace cambios en CoinTracking** (importa nuevos datos, edita operación)
2. **Cambio de proyecto activo** (cada proyecto tiene caché aislada)
3. **Cambio de exchange** (agregar nuevo, eliminar)
4. **Cambio de rango de fechas** (p. ej. auditor ahora cubre 2023-2026)
5. **Cambio de versión MCP** (formato de respuesta distinto)
6. **Cambio de ADR relacionado** (ADR-036, ADR-037, ADR-038, etc.)
7. **Cambio de base de conocimiento** (reglas fiscales, procedimientos)
8. **Cambio de algoritmo de auditoría** (nuevo método FIFO, nueva validación)
9. **Llamada explícita de usuario** (`cointracking_invalidate_cache`)

**Implementación:** Fase 4 (versionado automático) validará estas condiciones.

---

## Riesgos Conocidos y Mitigación

| Riesgo | Síntoma | Mitigación |
|--------|---------|-----------|
| Caché desincronizada | Auditoría reporta datos viejos | Fase 4: versionado automático |
| Hit rate bajo al inicio | Primeras auditorías sin ahorro | Aceptable; ahorro acumula en run 2+ |
| Overhead I/O | Lectura/escritura más lenta que MCP | Datos locales minimizan latencia |
| Crecimiento caché | `.cache/` → GB tras meses | Fase 4: limpieza LRU automática |
| Cambios regulatorios imprevistos | Auditoría antigua con reglas viejas | Versionado obliga reverificación |

---

## Implementación por Fases

| Fase | Estado | Entrega | Dependencias |
|------|--------|---------|--------------|
| **1** | ✅ Aceptada | `tools/cache_manager.py` | — |
| **2** | ✅ Aceptada | Integración en skills | Fase 1 |
| **3** | ✅ Aceptada | Validación en producción | Fase 1, 2 |
| **4** | ⏳ Planificada | Versionado automático | Fase 1-3 |
| **5** | ⏳ Planificada | TTL dinámico por tipo | Fase 4 |

**Documentos de soporte:**
- `docs/performance/TOKEN_BENCHMARKS.md` — cifras concretas (se actualiza cada trimestre)
- `implementation/CACHE_ROADMAP.md` — cronología de fases (se ajusta según necesidad)

---

## Validación (Aceptación Fase 3)

### Test en Producción

Ejecutado en agp2025 (1.670+ operaciones, 2024-2025):

**Flujo iterativo (3 auditorías + IRPF):**
- Sin optimizar: 30.305 tokens
- Con optimizar: 7.435 tokens
- **Ahorro: 75%**

**Escalabilidad (50 proyectos/año, 2 operaciones/proyecto):**
- Sin optimizar: ~1.000.000 tokens/año
- Con optimizar: ~335.000 tokens/año
- **Ahorro: ~665.000 tokens/año**

### Conclusiones

1. ✅ Caché persistente: funcionando, no degrada calidad
2. ✅ Procesamiento local: reduce contexto 90%+
3. ✅ Escalabilidad: viable para 50+ proyectos/año
4. ✅ Integridad: resultados sin cambios (solo costo)

---

## Referencias

**Principios relacionados:**
- ADR-010: Eficiencia de tokens
- ADR-037: Validación obligatoria en desarrollo
- ADR-038: Criterio de auditoría por lotes

**Documentos técnicos:**
- `tools/cache_manager.py` — implementación del gestor de caché
- `tools/ct_audit.py` — análisis local determinista
- `docs/performance/TOKEN_BENCHMARKS.md` — cifras concretas (actualizar trimestral)
- `implementation/CACHE_ROADMAP.md` — fases de implementación

**Casos de uso:**
- `reports/SKILLS_BENCHMARK_REPORT.md` — validación en caso real (agp2025)

---

**Decisión:** ACCEPTED ✓

**Versionado:** 1.0  
**Próxima revisión:** 2026-12-31 (anual) o si cambia modelo LLM  
**Responsable:** @agente-cointracking
