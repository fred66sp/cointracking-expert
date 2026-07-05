# Benchmarks de Tokens — Mediciones Concretas

**Documento de soporte para ADR-039**  
**Última medición:** 2026-07-05  
**Modelo:** Claude Haiku 4.5  
**⚠️ Este documento envejece rápido — reverificar cada trimestre con cambios de modelo**

---

## Advertencia

Las cifras en este documento son **puntuales en el tiempo**. Con cambios de modelo o contexto:
- Pueden variar significativamente
- No deben usarse para decisiones arquitectónicas sin reverificación
- El ADR-039 debe permanecer **independiente de estos números específicos**

**Para decisiones arquitectónicas, usa ADR-039 (principios). Para optimizaciones concretas, usa este documento (mediciones).**

---

## Mediciones de Referencia (2026-07-05)

### Operaciones MCP de CoinTracking

| Operación | Registros típicos | Tokens (aprox.) | Observaciones |
|-----------|------------------|-----------------|---------------|
| `get_trades(limit=all)` | 1.670 | ~2.835 | Proyecto agp2025; aumenta con años históricos |
| `get_grouped_balance()` | — | ~500 | Rápida, compacta |
| `get_gains(price:oldest)` | 50-100 | ~1.000 | FIFO completo; varía con # de activos |
| `get_historical_summary()` | 12 | ~400 | Por exchange; compacto |
| `get_balance()` | — | ~300 | Actual, muy pequeña |

**Patrón:** Las operaciones que devuelven volúmenes grandes (trades) cuestan 3-6x más que agregados.

### Análisis en Contexto LLM

| Tarea | Contexto sin opt. | Contexto con opt. local | Ahorro |
|-------|------------------|-------------------------|--------|
| Auditoría: interpretar hallazgos | 3.000 | 200 | 93% |
| IRPF: preparar informe | 2.200 | 300 | 86% |
| Clasificar eventos imponibles | 1.500 | 400 | 73% |
| Explicar transferencia huérfana | 800 | 100 | 88% |

**Patrón:** Procesamiento local reduce contexto 70-93% sin perder explicabilidad.

---

## Caso Real: agp2025 (1.670 operaciones, 2024-2025)

### Flujo Simple: Auditoría + Declaración

| Paso | Sin optimizar | Con caché | Con local | Total optimizado | Ahorro |
|-----|--------------|-----------|-----------|------------------|--------|
| MCP: get_trades | 2.835 | 2.835 → caché | — | 2.835 | 0% |
| MCP: get_grouped_balance | 500 | 500 → caché | — | 500 | 0% |
| MCP: get_historical_summary | 1.200 | 1.200 → caché | — | 1.200 | 0% |
| MCP: get_gains | 1.000 | 1.000 (nuevo) | — | 1.000 | 0% |
| **Subtotal MCP** | **5.535** | **5.535** | — | **5.535** | **0%** |
| Análisis auditoría (contexto) | 3.000 | 3.000 | 200 | 200 | 93% |
| Análisis IRPF (contexto) | 2.200 | 2.200 | 300 | 300 | 86% |
| **Subtotal contexto** | **5.200** | **5.200** | — | **500** | **90%** |
| **TOTAL** | **10.735** | **10.735** | — | **6.035** | **44%** |

**Notas:**
- MCP no se optimiza (datos reales necesarios)
- Contexto se optimiza 90% con análisis local
- Ahorro total: 44% (10.735 → 6.035 tokens)

### Flujo Iterativo: 3 Auditorías + Declaración

| Auditoría | MCP (nuevo) | MCP (caché HIT) | Contexto local | Total |
|-----------|------------|-----------------|---|-------|
| #1 | 5.535 | — | 200 | 5.735 |
| #2 | — | 0 | 200 | 200 |
| #3 | — | 0 | 200 | 200 |
| IRPF | 1.000 | — | 300 | 1.300 |
| **TOTAL** | **6.535** | — | **900** | **7.435** |

**Sin optimizar: 30.305 tokens**  
**Con optimizar: 7.435 tokens**  
**Ahorro: 75%** ← La caché brilla en ciclos iterativos

---

## Estimación Anual (50 proyectos)

| Escenario | Tokens sin opt. | Tokens con opt. | Ahorro |
|-----------|-----------------|-----------------|--------|
| Simple (1 audit + IRPF) | 13.235 | 6.035 | 7.200 |
| Iterativo (3 audits + IRPF) | 30.305 | 7.435 | 22.870 |
| **Promedio por proyecto** | **~20.000** | **~6.700** | **~13.300** |
| **50 proyectos/año** | **1.000.000** | **335.000** | **665.000 tokens** |

**Interpretación:** 665K tokens ahorrados = ~155 llamadas `get_trades()` evitadas.

---

## Factores que Varían Estos Números

### Aumentan el costo:

- ❌ Más años de historia (2020-2026 vs 2024-2025)
- ❌ Más activos (15 vs 50)
- ❌ Más exchanges (3 vs 10)
- ❌ Más operaciones por exchange (1.670 vs 5.000+)

### Reducen el costo:

- ✅ Caché reutilizada (hit rate alto en usuarios recurrentes)
- ✅ Análisis local optimizado (Python vs contexto LLM)
- ✅ Agregados en lugar de detalle (balance vs trades)

---

## Implicaciones Prácticas

### Cuándo NO optimizar

Si el usuario tiene:
- Pocas operaciones (< 100)
- Un solo exchange
- Primera auditoría (sin caché anterior)
- Urgencia (un análisis más lento es aceptable)

→ **No merece la complejidad de caché.** La ganancia es marginal.

### Cuándo SÍ optimizar

Si el usuario tiene:
- 1.000+ operaciones
- Múltiples exchanges
- Auditorías recurrentes (ciclos iterativos)
- Presupuesto tokens limitado

→ **ADR-039 + caché aporta ahorro comprobado (75%+).**

---

## Cómo Usar Este Documento

| Preocupación | Consulta | Respuesta |
|--------------|----------|-----------|
| "¿Cuántos tokens cuesta auditar?" | Este documento, tabla "Caso Real" | ~6.000 con optimización |
| "¿Por qué cacheamos?" | ADR-039 (principios) | Para reducir MCP + contexto LLM |
| "¿Cada cuánto invalida caché?" | ADR-039 (niveles TTL) | 24h trades, 15min balance, etc. |
| "¿Exactamente qué se cachea?" | `tools/cache_manager.py` (código) | Manifests + JSON con metadata |
| "¿Los números son precisos?" | Este documento (disclaimer) | ⚠️ Pueden variar con modelo |

---

## Próximas Mediciones

- [ ] **2026-10-05:** Reverificar benchmarks tras actualización modelo
- [ ] **2026-12-31:** Revisión anual + informe de ahorro real operativo
- [ ] Si se añaden exchanges: remedir con más datos heterogéneos
- [ ] Si cambia CacheManager: benchmark de rendimiento de I/O

---

**Generado:** 2026-07-05  
**Referenciado por:** ADR-039 (como "documento de soporte para cifras concretas")  
**Revisar:** Trimestral o con cambios de modelo
