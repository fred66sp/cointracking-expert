# Optimización sin sacrificar confianza: cómo escalar un agente auditor

Si leíste el artículo anterior, ya sabes que construí un agente de IA con arquitectura seria: MCP propio, separación dev/explotación, base de conocimiento auditable. Todo pensado para que confíes en él cuando audite dinero real.

Pero aquí viene el problema que nadie te menciona cuando construyes IA para producción: **funcionar es fácil. Escalar sin romperse es duro.**

El agente funcionaba. Auditaba cuentas. Pero empecé a notar que:
- Cada auditoría quemaba 8.000+ tokens de CoinTracking API (caras, y con límite de 60/hora)
- El contexto del agente crecía con cada operación auditada
- Las auditorías tardaban más conforme los datos aumentaban

Tenía dos caminos: 
1. Optimizar a ciegas. Cachear todo, cortar contexto, ir más rápido. El riesgo es obvio: pierdo trazabilidad.
2. Optimizar dentro de marcos documentados. Más lento de escribir, pero duermes tranquilo.

Elegí opción 2. Y eso requería algo que probablemente no esperabas: **Architectural Decision Records (ADRs)**.

## ¿Qué es un ADR y por qué importa para optimizar?

Un ADR es un documento que dice: "Aquí hay un problema, aquí están las opciones que consideré, aquí está la decisión que tomé, y aquí está el porqué."

En el software normal, los ADRs son documentación bonita que escribes después. En dominios de confianza (criptos, impuestos, finanzas), los ADRs son **barandillas de seguridad**. Son lo que permite que un dev optimice sin que todo se derrumbe.

¿Por qué? Porque cada decisión de arquitectura está escrita. Verificable. Si algo sale mal después, no es "el agente decidió", es "el agente aplicó ADR-039, que dice esto".

Cuando empecé a optimizar, tenía ADRs establecidos:

- **ADR-009:** Protocolo crítico del auditor (cero invención, todo trazable)
- **ADR-036:** Protocolo de desarrollo (cómo se aprueban cambios sin romper confianza)
- **ADR-037:** Pre-commit hooks (validación automática de cambios)
- **ADR-038:** MEGA audits (auditorías completas, no iterativas — evita ciclos de error)

Con eso como red de seguridad, pude escribir:

- **ADR-039:** Optimización de tokens y caché

## El núcleo: ADR-039

ADR-039 dice esto: "Necesitamos optimizar sin comprometer la integridad de auditoría." Y para eso, tres cosas:

**1. Integridad de auditoría:** La optimización jamás cambia el resultado. Si sin caché obtengo "ganancia = 45k€", con caché también. La optimización es invisible.

**2. No cachear conclusiones:** Solo cacheo datos brutos e intermedios reproducibles. Si cacheo "conclusión: esta operación es un duplicado", estoy arruinado. Cacheo "datos de Binance API a las 14:47", y el agente reasona sobre eso.

**3. Minimización de contexto:** Menos información en la conversación = tokens ahorrados sin perder análisis. En lugar de "aquí están tus 1.500 operaciones", envío agregados: "523 compras, 312 ventas, 5 transferencias huérfanas".

Con eso definido, implementé tres capas:

### Capa 1: Persistent Cache con Versionado (VersionTracker)

Cacheo las respuestas de la API de CoinTracking, pero con un truco: cada documento de la base de conocimiento y cada ADR tiene un número de versión. Si el conocimiento cambia (porque actualicé una regla fiscal), el caché se invalida automáticamente.

¿Por qué importa? Porque si CoinTracking cambia su API, o España cambia las reglas de Modelo 721, el agente lo sabe y refresca. No te quedas con conclusiones viejas.

### Capa 2: TTL Dinámico (CacheTTLManager)

No todo envejece igual. Las operaciones de trading (histórico de compra/venta) son permanentes — no cambian, así que cacheo para siempre. Los saldos envejecen cada 15 minutos — si hiciste una retirada hace 10 minutos, necesito el saldo fresco.

El agente decide automáticamente cuánto cachea cada cosa, basándose en qué tipo de dato es. Esto ahorra llamadas a la API sin sacrificar frescura donde importa.

### Capa 3: Métricas (CacheMetrics)

Aquí es donde se pone visual. Cada auditoría, el agente registra:
- Cuántos hits de caché tuvo
- Cuántos tokens ahorró
- Qué llamadas fueron más caras

El resultado: un dashboard que te muestra "esta semana ahorré 4.3M tokens, hit rate del 91%".

Eso **es** la optimización. Visible. Medible. No una caja negra.

## Los ADRs que lo posibilitaron

Aquí está lo importante: sin los ADRs anteriores (009, 036, 037, 038), **no habría podido hacer esto de forma segura**.

**ADR-036** (Protocolo de desarrollo) dice cómo se aprueban cambios. No es "hago un commit y listo". Es: cambio → validación automática → test → revisión → commit registrado. Cada paso auditable.

**ADR-037** (Pre-commit hooks) valida automáticamente antes de que un cambio llegue a main. Si intento cachear algo que no debería, el hook lo rechaza.

**ADR-038** (MEGA audits) dice que la auditoría debe ser **completa**, no iterativa. No es: audito, encuentro error, corrijo, audito de nuevo, encuentro otro error... Eso es ciclos infinitos. MEGA audit significa: verifico todo a la vez, listo.

Con eso en lugar, cuando optimicé con ADR-039, la confianza no bajó. Solo mejoró la eficiencia.

## Números reales

Audité el proyecto real (agp2025: 1.670+ operaciones, 3 exchanges, 6 meses de histórico):

- **Sin caché:** 8.535 tokens por auditoría completa
- **Con caché (hit 1):** 5.735 tokens (33% ahorro)
- **Con caché (hit 2):** 200 tokens (98% ahorro)
- **En casos iterativos** (usuario corrige errores, re-audita): 75% ahorro promedio

Anualizado: si 50 usuarios auditasen sus cuentas una vez al mes, eso es 620K tokens ahorrados. No es revolucionario, pero es **real, medible, y verificable**.

## El cambio de paradigma

Aquí está lo que me sorprendió:

Muchos devs piensan que "optimizar un agente de IA" significa "hazlo más rápido, gasta menos contexto". Fin. Una carrera de velocidad.

Pero cuando auditas dinero real, la optimización tiene que ser diferente:

**Antes:** "Optimización = velocidad". (Riesgo: pierdo trazabilidad.)

**Después:** "Optimización = eficiencia dentro de marcos documentados". (Resultado: más rápido Y más confiable.)

Los ADRs son lo que convierte la optimización de "magia negra" a "decisión transparente".

## Lección para CTOs

Si construyes IA escalable en dominios donde los números importan (finanzas, impuestos, salud), aquí está el order correcto:

1. **Primero:** Gobernanza. ADRs. Decisiones documentadas.
2. **Segundo:** Validación. Pre-commit hooks. Tests automáticos.
3. **Tercero:** Optimización. Ahora sí, optimiza. Sabes en qué marcos lo haces.

Optimizar primero = desastre. He visto devs hacer benchmarks increíbles, ahorrar 80% de tokens, y luego descubrir que cachean conclusiones. O que validan datos una sola vez y nunca refresque.

Con gobernanza primero, optimizar es seguro. Es invisible. Es confiable.

## ¿Y ahora qué?

El agente es rápido. Es eficiente. Pero sigue siendo honesto: no es determinista. No te dirá "tienes que declarar exactamente €45.000,42". Te dirá "parece que tienes una ganancia de ~45k€, aquí están los números que lo sustentan, verifica antes de declarar".

Eso no cambió con la optimización. Lo que cambió es que ahora lo hace sin quemar 8.500 tokens cada vez.

Y eso, para un dev que escaló un proyecto personal a herramienta de producción, es bastante.

---

*Próximo: Casos reales — lo que descubrí auditando criptos en la práctica (Copy Trading losses, peculiaridades de Binance, tratamiento de staking).*
