# Análisis de Brechas en ADRs — Sesión 2026-07-05

**Fecha:** 2026-07-05  
**Autor:** Análisis de Copilot sobre estructura de ADRs  
**Estado:** Identificadas brechas, solicitud de mejora creada

---

## Resumen ejecutivo

El proyecto tiene **32 ADRs bien estructurados** en gobernanza, arquitectura y principios. Sin embargo, **falta la Capa 2 (Conciliación)** — el motor operativo que define qué es una auditoría correcta en CoinTracking.

**Impacto:** Sin estos ADRs, el agente tiene gobernanza excelente pero carece de especificidad operativa en el dominio crítico.

---

## Estado actual (6 Capas propuestas)

### ✅ Capa 1 — Principios (COMPLETA)

ADRs 001, 002, 003, 004, 005, 006, 009, 026

- Principios de auditoría
- Fuente de verdad (Exchange vs CT vs Blockchain)
- Modelo de transacciones
- Reconciliación con datos reales
- Orden del diagnóstico
- Límites del determinismo
- Protocolo crítico
- Matriz de decisiones A/B/C

**Estado:** ✅ **ROBUSTA**

---

### ❌ Capa 2 — Conciliación (INCOMPLETA)

**ADRs faltantes:**

| ADR | Tema | Por qué importa |
|-----|------|-----------------|
| (nuevo) | Flujo de conciliación (pipeline) | Define orden invariante: importación → normalización → balances → transfers → duplicados → warnings → missing PH → holdings → FIFO → tax report |
| (nuevo) | Modelo de balances | Qué significa un balance "correcto". Cuándo es negativo. Cuándo parar. |
| (nuevo) | Missing Purchase History | Causa más común de auditoría falsa. Merece ADR dedicado (no solo mención en ADR-004). |
| (nuevo) | Transfers (emparejar) | Cómo enlazar withdrawal → blockchain → deposit. Tolerancias. |
| (nuevo) | Duplicados (clasificación) | FLOKI demostró que "misma fecha" no basta. Necesita matriz: Trade ID, Order ID, Hash, Cantidad, Precio, etc. |
| (nuevo) | Holdings (validación) | Comparar CT vs Exchange vs Wallet vs Blockchain. Reconciliación final. |
| (nuevo) | Cost Basis / FIFO (operativo) | Cuándo confiar en CoinTracking, cuándo recalcular. Discrepancias aceptables. |
| (nuevo) | Warnings (catálogo) | Cada warning: gravedad, impacto, acción. No todos los warnings pesan igual. |

**Estado:** ❌ **CRÍTICO — Falta**

---

### ⚠️ Capa 3 — Integración (PARCIAL)

ADRs 027 (nuevos exchanges), más faltantes:

| ADR | Tema | Por qué importa |
|-----|------|-----------------|
| ✅ 027 | Nuevos exchanges | ✅ EXISTE |
| ❌ (nuevo) | API de CoinTracking | Importación, límites, paginación, histórico, validación. |
| ❌ (nuevo) | CSV (formato y versionado) | Errores frecuentes, validaciones, cambios de formato. |
| ❌ (nuevo) | API + CSV (solapamientos) | Cómo gestionar cuando se usa ambas fuentes. Prioridad. Detección de duplicados entre fuentes. |
| ❌ (nuevo) | Binance específico | Tiene módulos independientes: Spot, Convert, Earn, Flexible Earn, Locked Earn, Rewards, Dust, Funding, Futures, Margin. Cada uno genera tipos de operaciones distintas. |

**Estado:** ⚠️ **Parcial (27% cubierto)**

---

### ⚠️ Capa 4 — Fiscalidad (PARCIAL)

ADRs 028 (límite auditor/asesor), 031 (plazos), más faltantes:

| ADR | Tema | Por qué importa |
|-----|------|-----------------|
| ✅ 028 | Límite auditor/asesor fiscal | ✅ EXISTE |
| ✅ 031 | Plazos y períodos | ✅ EXISTE |
| ❌ (nuevo) | FIFO (fiscalidad específica) | No solo "usar FIFO", sino cómo tratarlo en IRPF. |
| ❌ (nuevo) | Airdrops | Clasificación fiscal, valuación, registro. |
| ❌ (nuevo) | Staking y rewards | Rendimiento del capital vs ganancias. Valuación momento de recepción. |
| ❌ (nuevo) | Lending / Yield Farming | Tratamiento fiscal poco claro. Requiere fundamentación. |
| ❌ (nuevo) | NFTs | Ganancias patrimoniales vs colección. Valuación. |
| ❌ (nuevo) | LP Tokens, Wrapped Assets, Bridges, Hard Forks | Casos especiales, cada uno con tratamiento fiscal distinto. |

**Estado:** ⚠️ **Muy parcial (29% cubierto)**

---

### ✅ Capa 5 — Arquitectura (ROBUSTA)

ADRs 010, 013, 016, 020, 025, 027, etc.

- MCP, caché, persistencia, versionado, multi-proyecto, etc.

**Estado:** ✅ **ROBUSTA**

---

### ✅ Capa 6 — Gobernanza (ROBUSTA)

ADRs 030, 032, más principios transversales:

- Validación de ADRs críticos
- Conocimiento con vigencia temporal
- Trazabilidad, evidencias, confianza
- Reversibilidad

**Estado:** ✅ **ROBUSTA** (aunque ADR-032 implementa, faltan ADRs sobre confianza/evidencia específicamente)

---

## Brechas críticas identificadas

### 1. **Capa 2 está debajo del 20% de cobertura**

Sin los 8 ADRs de conciliación, el agente no puede definir qué es una auditoría correcta.

**Impacto:** Alto — esta es la razón de ser del proyecto.

---

### 2. **Capa 3 está debajo del 30% de cobertura**

Sin especificar cómo funciona API, CSV, y el comportamiento de exchanges como Binance, las importaciones son propensas a errores.

**Impacto:** Alto — fuente de duplicados y Missing PH falsos.

---

### 3. **Capa 4 está debajo del 30% de cobertura**

La fiscalidad es transversal. Sin ADRs específicos por tipo de operación, el agente no sabe cómo clasificar.

**Impacto:** Alto — riesgo de cifras incorrectas declaradas.

---

### 4. **Falta matriz de confianza/evidencia**

Copilot propone que cada conclusión del agente debería indicar:

```
Confidence: 95% | 80% | 60% | 20%
Evidencias: [API Binance, CoinTracking, CSV, Trade ID, Hash blockchain]
Reglas: [ADR-004, ADR-017]
¿Por qué?: [Explicación en lenguaje llano]
```

**Impacto:** Medio — mejora trazabilidad pero no es bloqueante.

---

## Los 15 ADRs prioritarios (según Copilot)

Si se tuviera que poner en producción mañana, estos serían imprescindibles:

1. ✅ ADR-001 — Principios de auditoría
2. ✅ ADR-002 — Fuente de verdad
3. ✅ ADR-003 — Modelo de transacciones
4. ❌ **ADR-??** — Flujo de conciliación
5. ✅ ADR-004 — Reconciliación con datos reales
6. ✅ ADR-005 — Orden del diagnóstico
7. ❌ **ADR-??** — Modelo de balances
8. ❌ **ADR-??** — Transfers
9. ❌ **ADR-??** — Missing Purchase History
10. ❌ **ADR-??** — Validación de duplicados
11. ❌ **ADR-??** — Cost Basis / FIFO
12. ✅ ADR-006 — Límites del determinismo
13. ✅ ADR-009 — Protocolo crítico
14. ✅ ADR-026 — Matriz de decisiones A/B/C
15. ❌ **ADR-??** — Trazabilidad de conclusiones

**Status:** 8 de 15 existen. **Faltan 7 críticos de la Capa 2.**

---

## Recomendación

**Prioridad 1 (BLOQUEANTE):**
- Capa 2: Los 8 ADRs de conciliación (flujo, balances, missing PH, transfers, duplicados, holdings, FIFO, warnings)

**Prioridad 2 (IMPORTANTE):**
- Capa 3: API de CoinTracking, CSV, Binance específico
- Capa 4: FIFO fiscal, Airdrops, Staking

**Prioridad 3 (MEJORA):**
- Matriz de confianza/evidencia por conclusión

---

## Próximos pasos

1. ✅ Crear solicitud de mejora en `AGENT_CHANGE_REQUESTS.md`
2. ⏭️ Sesión futura: Diseñar e implementar Capa 2 (conciliación)
3. ⏭️ Sesión futura: Completar Capa 3 (integración específica)
4. ⏭️ Sesión futura: Especificar Capa 4 (fiscalidad por tipo)

---

## Conclusión

El proyecto tiene **gobernanza y arquitectura excelentes**. Lo que falta es **especificidad operativa en el dominio** (Capa 2 — conciliación) que es el corazón del agente auditor.

Sin ella, el agente puede ser robusto, pero no puede ser **específico** a CoinTracking.

Es como tener un protocolo de calidad pero no definir qué es "calidad" en el dominio específico.
