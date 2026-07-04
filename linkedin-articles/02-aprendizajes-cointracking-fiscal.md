# Lo que aprendí: CoinTracking, fiscalidad y reconciliación

En el artículo anterior conté por qué decidí construir un agente de IA para auditar mis criptos. Ahora voy a compartir lo que descubrí en el camino. Adelanto: el problema no es tan simple como "CoinTracking está roto".

## Descubrimiento 1: CoinTracking funciona, pero necesita contexto

Lo primero que aprendí es que CoinTracking no es una herramienta mala. Funciona bien si sabes qué estás haciendo. El problema es que **nadie explica realmente qué estás haciendo**.

CoinTracking importa tus datos de los exchanges usando API o CSV, y luego aplica un modelo de coste para calcular ganancias. Suena simple. Pero aquí vienen los matices:

El modelo de coste que usa es **FIFO** (Primero en Entrar, Primero en Salir): la primera crypto que compraste es la primera que vendiste. Eso es correcto fiscalmente en España, pero CoinTracking asume que sabes que así funciona. Si tú, en tu cabeza, esperabas que te calcule los costes de otra forma (por ejemplo, identificando específicamente qué monedas vendiste), te va a parecer que está mal cuando en realidad es solo... diferente.

Segundo problema: **los datos que importa son tan buenos como la fuente**. Si Binance reporta una operación con fecha y hora incorrecta, CoinTracking la importa así. Si hiciste una transferencia entre billeteras y CoinTracking no detecta que es tuya (porque no está documentada explícitamente), la marca como "transferencia huérfana". No es fallo de CoinTracking; es que tu exchange no envió la información correcta, o no hubo manera de conectar los puntos.

El verdadero descubrimiento: **CoinTracking es un espejo de tus datos. Si tus datos están desordenados, el reflejo está desordenado.**

## Descubrimiento 2: Fiscalidad española hace esto exponencialmente más difícil

Y aquí es donde entra Hacienda a fastidiarte la tarde.

España tiene reglas complicadas para criptos: cada vez que vendes o permuttas (eso es intercambiar una cripto por otra), generaste una ganancia patrimonial que **hay que declarar en el IRPF**.

Y luego está el Modelo 721. Si el valor conjunto de tus criptomonedas custodiadas en entidades situadas en el extranjero supera los 50.000€ a 31 de diciembre, puede que tengas que presentarlo. Es una declaración informativa (no es un impuesto), pero es obligatoria si se supera el límite. Y aquí viene lo importante: **solo cuentan las que están custodiadas por terceros**. Si tus criptos están en una billetera física (donde controlas tú las claves), no entran en el cómputo.

Porque una cosa es que CoinTracking te calcule las ganancias; otra es entender qué tienes que declarar dónde.

Pero hay más:

- **Ganancias patrimoniales**: se calculan como (Precio de venta - Precio de coste). El precio de coste es donde todo se vuelve manual. CoinTracking te lo calcula automáticamente con FIFO, pero ¿y si vendiste en pérdida? ¿Y si cambiaste de exchange y perdiste el histórico de una operación?
- **Modelo 721**: hay que reportar holdings finales, no solo ganancias. ¿Cuál es tu saldo en BTC al 31 de diciembre? Si no coincide entre CoinTracking y tu exchange porque falta una transferencia, te hace falta detective work.
- **Transferencias entre billeteras**: si transferiste 1 BTC de Binance a tu billetera física, eso es una transferencia, no una venta. CoinTracking debería detectarlo automáticamente, pero si no tuvo visibilidad de ambos lados, la registra como "retirada" en un lado y falta la visibilidad del otro.

Entonces llega tu asesor fiscal o contable, ve las números, levanta una ceja y te pregunta: "¿Por qué CoinTracking dice que tienes 2 BTC y el año pasado declaraste 1.8?" 

Y tú te quedas en blanco. Porque no tienes ni idea qué pasó en noviembre.

## Descubrimiento 3: Reconciliación es arqueología de datos

Esto es lo que más me sorprendió.

Reconciliar tus criptos no es sentarse a mirar números. Es trabajo real de detective:

- Ves una transferencia que falta en CoinTracking. Vas a Binance, hurgas en el historial, la encuentras. Salió en noviembre a las 14:47.
- Rastreas dónde llegó. Tu billetera física. ¿La registraste en CoinTracking? No. ¿Por qué? Porque asumiste que solo necesitabas los exchanges, no tus billeteras propias.
- Resultado: CoinTracking tiene "Retirada 1 BTC" pero ningún "Depósito" que lo explique. 1 BTC fantasma.

O este otro caso real que encontré: Binance procesó 29 compras de la misma moneda en el mismo segundo. Idénticas. Misma fecha, misma hora, mismo precio, mismo volumen. En CoinTracking parecían duplicados obvios.

Así que probé a eliminar algunos. Saldo: -5 BTC. Error fatal.

Luego investigué en la API de Binance: cada operación tenía su propio identificador único. Eran 29 operaciones legítimas, no un error. Binance las procesó en paralelo en el mismo segundo. Sin esa API, habría eliminado datos reales y arruinado todo.

## Descubrimiento 4: Por qué un agente de IA es mejor que un script burdo

Un script determinista es rápido: "¿Dos operaciones con el mismo precio y hora? Duplicado. Elimina." Pero es peligrosamente frágil. En el caso de Binance, habría matado datos reales.

Un agente de IA puede razonar:

- "Ves dos operaciones iguales. Espera, voy a revisar la API de Binance. Ah, cada una tiene un ID único. No son duplicadas."
- "Falta un depósito que se corresponde con una retirada. ¿Qué pasó? Dame contexto. ¿Confirmado en blockchain? ¿Llegó a tu billetera?"
- "Tú dices que ganaste 50k€. CoinTracking me muestra 45k€. ¿Qué 5 operaciones no encajan? Aquí están. Verifica antes de declarar."

El agente de IA no es infalible, pero es **sincero**. No solo te dice "esto está mal"; te explica **por qué** lo cree y qué información necesita para estar seguro. Eso importa cuando auditas dinero real.

## Resumen

CoinTracking funciona. La fiscalidad española te complica la vida. Y reconciliar requiere ser detective, detective y más detective.

Pero aquí viene lo interesante: ¿cómo construyes un agente que audite dinero real sin que sea una caja negra? ¿Cómo le das seguridad al usuario de que no va a arruinarle los números?

Eso es lo que cuento en el siguiente artículo. Y spoiler: la arquitectura importa.

---

*Próximo: "Arquitectura: cómo construí un agente confiable para dinero real"*
