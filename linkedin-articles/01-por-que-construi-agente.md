# Por qué construí un agente de IA para auditar mis criptos

Tengo posiciones de cripto en 5 exchanges diferentes. Sé que suena como el inicio de un chiste sobre diversificación, pero la realidad es más deprimente: cada plataforma tiene su propia interfaz, formatos de datos y, aparentemente, su propia versión de cuántos euros gasté en noviembre de 2023.

Cuando llegó abril de 2024 y mi asesor fiscal me pidió "simplemente" que consolidara todo para la declaración de la renta, pensé: "Fácil, metemos todo en CoinTracking y listo". CoinTracking promete exactamente eso en su marketing: importa tus trades de tus exchanges, te calcula ganancias, Modelo 721, todo resuelto.

Si estás leyendo esto ahora que acabamos de cerrar 2025 y tienes la renta encima, probablemente reconozcas exactamente esta situación.

La realidad fue diferente.

## El problema que probablemente reconocerás

Los números no cuadraban. Una operación aparecía duplicada en un exchange, otra no aparecía en CoinTracking, una tercera estaba con el precio incorrecto. Y lo peor: no sabía cuál era la fuente de verdad. ¿Era fallo del exchange? ¿De CoinTracking? ¿Mío, por no saber cómo usar la herramienta?

Pasé semanas buscando en documentación de ayuda, intentando entender qué significa que una "transferencia sea huérfana", por qué CoinTracking calculaba un coste base diferente al que yo esperaba, o cómo es posible que dos importaciones del mismo archivo CSV tuvieran resultados distintos.

Y aquí viene lo realmente frustrante: **nada de esto era culpa mía por ser principiante**. Resulta que es un problema estructural. CoinTracking es una herramienta compleja —y no mala, en el fondo— pero está diseñada para contadores que ya entienden fiscalidad, no para devs como yo que solo queremos que sus números cuadren.

## Un descubrimiento incómodo

Un día googleo "CoinTracking problems" y caigo en un agujero de conejo. Resulta que decenas de empresas, desde Shopify apps hasta empresas de consultoría fiscal, ofrecen exactamente **un único servicio**: reconciliar tus datos de CoinTracking. Revisar duplicados, buscar transferencias que no aparecen, validar saldos, explicarte por qué tu ganancia en CoinTracking es diferente a la que esperabas.

Y cobran. Bastante. He visto precios desde €200 hasta €5000+ dependiendo de la complejidad. Para un problema que debería ser... automático.

En ese momento tuve dos opciones:

1. Pagar a alguien para que pase 4 horas mirando mis datos en una hoja de cálculo.
2. Invertir esas 4 horas en construir algo que lo hiciera automáticamente.

Soy dev. Opción 2 era inevitable.

## El giro

Empecé a pensar: si CoinTracking es complicado y la API es accesible, ¿por qué no construir un agente que literalmente haga lo que hacen esos consultores? Que importe mis datos, cruce referencias con los exchanges originales, encuentre inconsistencias y me *explique* qué está mal y por qué.

Pero aquí es donde la cosa se pone interesante. No quería un script determinista que dijera "Error: saldo negativo en BTC". Necesitaba algo que pudiera razonar sobre el dominio, explicar por qué dos operaciones que parecen iguales no son duplicados, o por qué una transferencia desapareció en el camino.

Eso solo se resuelve con IA. Y eso requiere arquitectura seria. Porque si voy a auditar dinero real —lo que voy a declarar a Hacienda— no puedo confiar a ciegas en una caja negra. Necesitaba trazabilidad total, gobernanza clara, y separación nítida entre lo que es determinista y lo que es razonamiento.

Así que construí un agente de IA especializado en reconciliar criptos.

## Lo que descubrí en el camino

En el proceso aprendí cosas interesantes:
- CoinTracking no es el culpable; la verdadera complejidad está en cómo los exchanges reportan los datos.
- La fiscalidad española lo complica todo (Modelo 721, ganancias patrimoniales, reglas de coste FIFO...).
- Un LLM que explique sus razonamientos es *mucho* mejor para reconciliación que un algoritmo determinista (y explico por qué en el siguiente artículo).
- Para auditar dinero real, necesitas arquitectura seria. Traza, gobernanza, confianza.

## ¿Por qué cuento esto aquí?

Si tienes cripto en múltiples exchanges, esto te va a sonar familiar. Y construir esto me enseñó un montón sobre IA en problemas reales, no en sandbox.

Voy a compartir dos artículos más:
1. **Lo que descubrí** sobre CoinTracking, fiscalidad y reconciliación
2. **Cómo lo construí** con confianza cuando los errores salen caros

Si trabajas con criptos, datos financieros o IA en ámbitos donde los números importan, algo de esto te puede servir.

Una pregunta: **¿tu experiencia con CoinTracking fue parecida? ¿Cómo lo resolviste?** Cuéntame en comentarios.

---

*Próximo: "Lo que aprendí sobre CoinTracking, fiscalidad y reconciliación"*
