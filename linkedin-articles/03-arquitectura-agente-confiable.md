# Arquitectura: cómo construí un agente confiable para dinero real

Si leíste los artículos anteriores, ya sabes el problema: CoinTracking es complejo, la fiscalidad española es más compleja todavía, y reconciliar tus criptos requiere ser detective. Ahora viene lo interesante: cómo construyo un agente de IA que audite dinero real sin que sea una caja negra.

## El reto que no es trivial

Aquí está el dilema: un LLM no es determinista. Eso significa que no es una calculadora. No te dirá "la ganancia es €50.000 exactos" de la misma forma que lo hace Excel. Razona, explica, pero no garantiza el mismo resultado dos veces.

Eso es un problema cuando lo que estás auditando va a Hacienda.

¿Cómo confío en un agente que no es 100% predecible? ¿Cómo evito que una alucinación del modelo (dirá algo que no es cierto) destruya un saldo o sugiera eliminar una operación legítima?

La respuesta no es "confía ciegamente". Es arquitectura.

## Solución 1: Fuente de datos bajo control (MCP propio)

Construí un servidor propio en Go que conecta directamente con la API de CoinTracking y con los datos que el usuario aporta (CSV, billeteras, exchanges).

¿Por qué no usar la API directamente? Porque necesito control total:

- **Caché de datos**: las respuestas de CoinTracking son grandes. No quiero hacer 100 llamadas a la API por cada auditoría. Cacheo lo que veo, con timestamp, para saber cuándo fue el último acceso.
- **Multi-proyecto**: si auditas una cuenta, eso aísla qué datos se usan. Los datos financieros reales nunca se mezclan entre proyectos.
- **Validación en la puerta**: antes de que el agente vea los datos, los valido. Formato correcto, fechas coherentes, saldos razonables. Si algo huele mal, lo digo antes de dejar que el LLM haga magia.

El servidor es simple pero es *mío*. Eso significa que si algo sale mal, sé exactamente qué pasó.

## Solución 2: Gobernanza explícita (ADRs)

"ADR" significa "Architectural Decision Record" — es un documento que dice: "Aquí hay un problema, aquí están las opciones, aquí está la decisión que tomé, y aquí está el porqué."

Desde día uno documenté decisiones:

- ¿Cuál es el modelo de coste? FIFO (Primero en Entrar, Primero en Salir).
- ¿Cómo manejo las transferencias huérfanas? Primero pregunto; no asumo.
- ¿Qué es un duplicado y cómo lo detecto? Aquí está la regla, con ejemplos.
- ¿Dónde empieza y dónde termina la confianza del agente? Aquí, en esta línea.

Todo esto está escrito. No en Slack, no en mi cabeza. Documentado. Porque cuando auditas dinero real, la trazabilidad importa. Si alguien pregunta después "¿por qué eliminaste esa operación?", tengo una respuesta.

## Solución 3: Separación clara entre desarrollo y explotación

Aquí está la arquitectura que más me enorgullece:

**Claude Code** (mi CLI para dev) es donde construyo, actualizo, mejoro el agente. Donde cambio reglas, añado conocimiento, tomo decisiones.

**GitHub Copilot** (el editor, same model pero diferente contexto) es donde *se usa* el agente. Donde auditas tus datos y preparas la declaración. Sin poder modificar nada.

¿Por qué esta separación? Porque:

- Los que auditamos con Copilot no queremos que el agente se redefina mientras lo usamos. La confianza requiere estabilidad.
- Los cambios del agente tienen trazabilidad. Cada mejora es un commit documentado.
- Si algo sale mal en una auditoría, sé exactamente qué versión del agente se usó.

Es como tener un abogado que actualiza su expertise en su despacho, pero cuando representa a un cliente, aplica lo que sabe en ese momento. No cambia de opinión a mitad del juicio.

## Solución 4: Base de conocimiento auditable

El agente no es solo código. Tiene una "base de conocimiento" — documentos que explican:

- Cómo funciona CoinTracking (formato CSV, modelo de coste, quirks de la API).
- Cómo funciona la fiscalidad española (Modelo 721, ganancias patrimoniales, tramos de IRPF).
- Contexto de exchanges (qué es MiCA, cuándo Binance cerró en la UE, regulaciones).
- Casos reales curados (esto es lo que encontré, esto es cómo lo resolví).

Todo esto es texto. Legible. Verificable. Si el agente me dice "aquí hay un problema porque es transferencia huérfana", puedo ir a la base de conocimiento y ver exactamente qué regla aplicó.

No es una caja negra. Es razonamiento transparente.

## Trade-offs: qué NO hace el agente

Aquí es donde tengo que ser honesto.

El agente **no es determinista**. Si le pides que audite la misma cuenta dos veces, las explicaciones serán ligeramente distintas. Pero las conclusiones (qué está mal) serán las mismas.

El agente **no produce cifras fiscales vinculantes**. Te dice "parece que tienes una ganancia de 45k€", no "debes declarar 45k€ a Hacienda". Eso es trabajo para un contable.

El agente **no reemplaza un abogado o un asesor fiscal**. Es una herramienta de reconciliación y diagnóstico. Te muestra problemas, te explica qué son, pero la decisión fiscal es tuya.

¿Por qué estos límites? Porque la responsabilidad es tuya. Yo construyo herramientas. Tú usas esas herramientas en tu contexto. Si algo sale mal, quiero que sea transparente dónde fue el límite.

## Por qué esto importa

Si construyes IA para dominios "serios" — finanzas, salud, legal, impuestos — no puedes simplemente meter un LLM en producción y esperar lo mejor.

Necesitas:

- Fuentes de datos bajo control.
- Decisiones documentadas.
- Separación clara entre desarrollo y uso.
- Explicabilidad en cada paso.
- Honestidad sobre límites.

Eso es arquitectura para confianza. Y confianza es lo que permite que un dev audite sus propios impuestos sin pagar a una empresa 500€ por hacerlo.

¿Merece la pena? Sí. Ha tardado más de lo que hubiera imaginado. Pero cuando estés auditando dinero que va a Hacienda, ese esfuerzo en arquitectura es lo que te deja dormir tranquilo.

---

*Y eso es la serie. Si te interesa el código, la base de conocimiento o cómo usarlo, cuéntame.*
