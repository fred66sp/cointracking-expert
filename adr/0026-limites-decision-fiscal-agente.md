---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-026: Límites de decisión fiscal — qué decide el agente solo vs qué requiere intervención humana

**Status:** Accepted

**Date:** 2026-07-04
**Accepted:** 2026-07-05 — la matriz A/B/C ya está en uso de facto (protocolo de consentimiento informado en CLAUDE.md, reglas de confirmación explícita antes de borrar duplicados en la skill `audit-cointracking`). Los 6 pendientes de refinamiento (criterios automáticos A vs B, protocolo uniforme de confirmación, etc.) quedan como trabajo futuro, no bloquean la aceptación.

## Context

El proyecto define en ADR-006 un "límite de determinismo": el agente es un LLM explicable que **encuentra y explica** problemas (análisis cualitativo), pero **no produce cifras fiscales vinculantes** por sí mismo. Sin embargo, la práctica operativa ha revelado un hueco crítico: no está claro *en qué decisiones concretas* el agente actúa de forma independiente vs. cuándo debe parar y pedir confirmación humana.

### Caso real (2026-07-03): El incidente FLOKI

Durante la auditoría de la cuenta `agp2025`, el agente detectó 29 operaciones FLOKI con el mismo instante temporal, precio y volumen registrados en CoinTracking. El patrón coincidía con un indicador habitual de posibles transacciones duplicadas.

Al no existir una regla explícita de autorización previa para actuaciones potencialmente irreversibles, el usuario eliminó dichas operaciones directamente en CoinTracking.

**Resultado:** ~1,6 millones de FLOKI desaparecieron del saldo.

**Investigación posterior (ADR-014, ADR-019):** Las operaciones eran **legítimas**. Binance había ejecutado órdenes paralelas en el mismo segundo (batching), diferenciadas mediante `Trade ID` independientes. El patrón coincidía con duplicados reales, pero la similitud era únicamente coincidencia, no evidencia de error.

**Raíz del problema:** 
- Detectar un patrón sospechoso no implica demostrar que existe un error
- Una recomendación técnica desencadenó una modificación irreversible de datos
- El agente no tenía un gate que detuviera la cadena de decisiones antes de la acción
- El proyecto necesita separar explícitamente el **diagnóstico técnico** de la **toma de decisiones con impacto contable o fiscal**

**Riesgo general:** Este proyecto emite **informes que van a Hacienda**. Cada cifra o decisión que se declare incorrectamente tiene consecuencias fiscales reales. Necesitamos una línea clara: hasta dónde puede llegar el agente sin intervención humana, y dónde **debe detenerse y esperar confirmación**.

## Decision

Se establece una **matriz de decisiones obligatoria** para todas las actuaciones del agente relacionadas con auditoría, reconciliación y fiscalidad. Las decisiones se clasifican en tres categorías según su reversibilidad e impacto fiscal.

### Categoría A: Decisiones informativas (el agente decide de forma autónoma)

Corresponden a actividades que **no modifican datos** ni generan consecuencias fiscales directas. El agente puede ejecutarlas automáticamente siempre que diferencie claramente entre hechos observados e hipótesis de diagnóstico.

**Ejemplos:**
- Detectar balances negativos en un activo
- Identificar posibles transposiciones de activos (BTC comprado, pero registrado como venta)
- Detectar posibles operaciones duplicadas **como hipótesis** (no como hecho)
- Localizar advertencias ("warnings") de CoinTracking
- Informar sobre posibles inconsistencias de holdings
- Explicar el funcionamiento de FIFO, staking, transferencias o comisiones
- Resumir resultados de una auditoría
- Clasificar riesgos (bajo, medio, alto) indicando la evidencia utilizada

### Categoría B: Diagnóstico con confirmación explícita del usuario

Corresponden a recomendaciones cuya ejecución **pueda modificar información contable**, afectar la reconciliación o alterar el resultado fiscal. **Antes de continuar, el agente debe solicitar confirmación explícita del usuario.**

La confirmación debe realizarse después de presentar:
- La evidencia disponible
- Las limitaciones del diagnóstico
- Los riesgos conocidos
- Las comprobaciones recomendadas (p. ej. revisión de `Trade ID`, `Order ID`, información del exchange)

**Ejemplos:**

1. **Eliminar operaciones consideradas duplicadas:** El caso FLOKI demuestra que la coincidencia temporal, de precio y volumen constituye únicamente un indicio, no una prueba suficiente para eliminar registros. Pedir verificación de `Trade ID` primero.

2. **Fusionar operaciones:** Cuando dos operaciones parecen ser componentes de una sola (p. ej. dos mitades de una venta).

3. **Reclasificar tipos de transacción:** Si el tipo actual es incorrecto (p. ej. "Venta" que debería ser "Permuta"). Incluir referencia a la regla fiscal.

4. **Modificar transferencias:** Cuando faltan datos de origen o destino.

5. **Sustituir datos importados:** Reemplazar información que provenía de una importación anterior.

6. **Repetir una importación eliminando datos existentes:** Potencialmente destructivo si hay inconsistencias no resueltas.

7. **Reconstruir base de coste:** Cuando existan varias alternativas válidas o discrepancias FIFO vs. cálculo manual.

8. **Aceptar una reconciliación basada en hipótesis no verificadas completamente:** Si depende de datos externos no confirmados.

### Categoría C: Decisiones delegadas obligatoriamente a intervención humana

El agente **no debe adoptar estas decisiones** ni recomendar una actuación definitiva cuando no exista evidencia suficiente o cuando la decisión corresponda al ámbito profesional o legal. En estos casos, debe detener el proceso y solicitar revisión humana.

**Ejemplos:**

1. **Decidir el tratamiento fiscal correcto entre interpretaciones posibles:** Múltiples lecturas válidas de normativa española requieren criterio profesional.

2. **Validar una declaración tributaria para presentación:** La responsabilidad es del asesor fiscal, no del agente.

3. **Determinar la base imponible ante incertidumbres jurídicas:** Requiere análisis legal especializado.

4. **Emitir asesoramiento fiscal definitivo:** "Tus ganancias patrimoniales son X € y debes declarar Y" (estas cifras vinculantes son responsabilidad del contable).

5. **Asumir que un cálculo fiscal es correcto sin reconciliación suficiente:** Incluso si CoinTracking lo afirma, debe verificarse contra Tax Report oficial.

6. **Confirmar la eliminación de datos cuando la evidencia sea insuficiente:** El caso FLOKI muestra el riesgo de actuar con evidencia incompleta.

7. **Resolver discrepancias cuya única fuente sea una inferencia del agente:** Si no hay respaldo externo, no es decisión del agente.

## Consequences

**Positive:**

- **Claridad operativa:** el agente, Copilot y el usuario saben exactamente dónde está la línea entre diagnosis y decisión
- **Protección contra decisiones irreversibles no validadas:** como el caso FLOKI, donde un patrón sospechoso desencadenó la pérdida de 1,6M FLOKI
- **Separación explícita:** observación → diagnóstico → decisión son fases distintas, no comprimidas en una sola
- **Alineación con ADR-009:** refuerza el principio de "consentimiento informado antes de actuar" con gates concretos
- **Responsabilidad clara:** el agente detecta, explica y sugiere; el usuario y el contable deciden
- **Trazabilidad mejorada:** cada decisión de Categoría B queda registrada (usuario confirmó el riesgo)
- **Escalabilidad:** cuando aparezca un nuevo tipo de operación (staking, DAC8, rewards), hay un criterio para clasificarlo
- **Reducción de fricciones:** el agente no recomendará acciones sin un criterio explícito sobre confirmación

**Negative:**

- **Más lento:** confirmaciones adicionales = más pasos antes de cerrar un informe
- **Requiere disciplina:** del agente (no "saltarse" los gates) y del usuario (no ignorar advertencias)
- **Ambigüedad residual:** en la frontera entre Categoría A y B siempre habrá grises
- **Responsabilidad compartida, no delegada:** el usuario debe verificar a veces, no solo confiar
- **Riesgo de confirmación ciega:** si el usuario aprueba sin leer, el problema persiste

## Notes

### Relación con ADRs existentes:

- **ADR-006 (límite de determinismo):** Este ADR operacionaliza el límite teórico. ADR-006 dice "no es determinista"; ADR-026 dice "por eso, aquí NO decides, aquí PARAS."
- **ADR-009 (protocolo crítico):** Refuerza "consentimiento informado antes de actuar" con matrices concretas. ADR-009 es el principio; ADR-026 es la aplicación operativa.
- **ADR-012 (división de responsabilidades):** Clarifica quién decide en qué. Copilot (explotación) aplica estas matrices; si algo está en Categoría C, escala.
- **ADR-014 (validación de duplicados):** Caso concreto de Categoría B — detección automática, pero confirmación obligatoria antes de borrar (como debería haber sucedido con FLOKI).

### Fundamento del ADR:

Este ADR se basa en el análisis real del incidente FLOKI (ADR-014, ADR-019) donde un patrón sospechoso (coincidencia temporal, precio, volumen) resultó corresponder a operaciones legítimas. Operacionaliza el principio de máxima cautela de ADR-009: un patrón sospechoso nunca constituye evidencia suficiente para ejecutar acciones irreversibles.

### Pendientes abiertos:

1. **[PENDIENTE]** Definir criterios objetivos para clasificar automáticamente un caso entre Categoría A y Categoría B.
2. **[PENDIENTE]** Establecer si determinadas comprobaciones (`Trade ID`, `Order ID`, hash de blockchain) pueden reducir automáticamente el número de confirmaciones necesarias.
3. **[PENDIENTE]** Diseñar un protocolo uniforme de confirmación para acciones irreversibles.
4. **[PENDIENTE]** Definir la responsabilidad residual cuando el usuario autorice una acción basada en evidencia incompleta.
5. **[PENDIENTE]** Evaluar si futuras integraciones mediante MCP permiten elevar determinados diagnósticos de Categoría B a Categoría A gracias a evidencia adicional verificable.
6. **[PENDIENTE]** Automatizar la lista de decisiones Categoría A tomadas al cerrar una auditoría, para que el usuario sepa qué se decidió sin pedirlo.

### Verificaciones:

- Este ADR se alinea con ADR-006 (límite de determinismo) y ADR-009 (protocolo crítico)
- El caso FLOKI documenta el riesgo real de actuar sin confirmación
- Las tres categorías reflejan la experiencia operativa del proyecto y su necesidad de gobernanza fiscal
