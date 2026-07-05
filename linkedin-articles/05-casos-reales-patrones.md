# Casos reales: los patrones que encontré auditando criptos

Si leíste los artículos anteriores, ya sabes que construí un agente auditor, lo hice confiable, y lo optimicé sin romperlo. Pero hay una pregunta que falta: **¿qué encontré realmente?** ¿Cuáles son los patrones que reaparecen una y otra vez?

Aquí es donde la teoría se convierte en realidad sangrienta de hojas de cálculo.

Audité una cuenta real compleja (6 meses, 3 exchanges, 1.670+ operaciones) y encontré tres patrones que, estoy convencido, tú también encontrarás. Especialmente si tienes cripto en Binance, exchanges derivados, o haces staking.

## Patrón 1: BingX Copy Trading — el dinero que desaparece

Esto es específico, pero es un horror real.

BingX es un exchange chino con una característica llamada "Copy Trading" — básicamente contratas a un trader y ella automatiza tus operaciones copiando las suyas. Suena bien en marketing. En práctica, es un lío.

Cuando audité la cuenta, encontré esta línea:

```
Type: Lost
Amount: -694.67 USDT
Date: (vacío)
Description: (vacío)
```

Sencillo: **no había depósito, no había explicación, desapareció dinero de la nada.**

Investigué. Resulta que BingX tiene sub-cuentas separadas (Spot, Margin, Futures, Copy Trading). El import de CoinTracking captura las tres primeras automáticamente por API. **Copy Trading no.** BingX simplemente no lo exporta.

Así que el dinero se perdió en el limbo entre "confirmado en BingX app" y "no visto en CoinTracking".

### Por qué importa fiscalmente

Aquí viene lo malo: si fue una pérdida legítima, probablemente **se puede deducir en el IRPF**. Pérdidas patrimoniales se compensan contra ganancias.

Pero si no tienes justificante oficial de BingX... ¿cómo lo demuestras a Hacienda?

**La solución:** Contacté a BingX, exporté el historial del Copy Trading account manualmente (screenshots), y documenté todo. Guardé evidencia. Porque si Hacienda pregunta, tengo prueba.

Nota: esto **no es único a BingX**. Cualquier exchange derivado (Bybit, OKX) con sub-cuentas puede tener el mismo problema.

## Patrón 2: Staking — la clasificación fiscal que nadie entiende

Este es más común, pero nadie lo sabe hacer bien.

En el 2024, la cuenta hizo staking de ETH en Binance. 10 ETH bloqueados, recompensas mensuales. Resultado: 0.42 ETH en recompensas a lo largo del año.

En CoinTracking aparecía así:

```
Type: Income
Amount: 0.05 ETH (aprox.)
Date: 2024-06-15
Description: "Binance Earn: ETH Staking Reward"
```

Hasta aquí, CoinTracking lo capturó bien. El problema es qué haces después.

### El error clásico

El 90% de los usuarios que audité pensaban esto: "Gané 0.42 ETH de staking. Eso es una ganancia. IRPF."

Incorrecto. **Completamente incorrecto.**

En España, staking se clasifica como **RCM (Rendimiento de Capital Mobiliario)**. Es decir, intereses. No ganancias patrimoniales. La diferencia importa para el tramo impositivo y cómo se declara.

Peor aún: el valor a declarar es el **EUR a fecha de recepción de cada recompensa**, no el EUR de hoy.

Ejemplo:
- 15 junio 2024: recibo 0.05 ETH, precio 2.500€/ETH → declaro 125€
- Hoy (julio 2026): ese ETH vale 3.200€ → pero **no declaro 160€**, declaro 125€

Eso es el EUR que recibiste, punto.

### Por qué importa

Si cometes este error, sobrepasas o infradeclaras ganancias. Hacienda lo nota. Y en dominios fiscales, "no sabía la regla" no es defensa válida.

Con el agente, es simple: cada recompensa se marca como RCM, se valora a fecha de recepción, y se pasa al asesor con etiqueta clara.

## Patrón 3: Transferencias huérfanas — el detective work invisible

Este es quizás el patrón más común.

Transferiste 5 BTC de Binance a Coinbase. En Binance sale: "Withdraw 5 BTC, 2026-07-01 10:00". En Coinbase debería entrar: "Deposit 5 BTC, 2026-07-01 15:30" (4 horas después, blockchain confirmation).

Pero CoinTracking solo vio una de las dos. Resultado: **BTC negativo en Binance**.

Eso es una "transferencia huérfana" — la retirada existe, pero el depósito no se conecta.

### Las 4 causas que encontré

**Causa 1: Blockchain delay** (más común). La transferencia tarda más de 24 horas, y CoinTracking importa datos por fechas. Hoy reimportas y aparece.

**Causa 2: Importación parcial.** Importaste Binance 2024-2025, pero Coinbase solo 2025. La retirada de 2024 está en Binance, pero el depósito de 2024 no está en Coinbase (porque no lo importaste).

**Causa 3: Depósito rechazado.** Rarísimo, pero ocurre. La transferencia es legítima (confirmada en blockchain), pero Coinbase la rechazó por algún motivo. El dinero... existe en blockchain pero no en Coinbase.

**Causa 4: Migración incompleta.** Iniciaste retirada, pero nunca completaste. Dinero en limbo.

### Cómo lo resolví

Para cada caso, seguí un procedimiento de dos niveles:

**Nivel 1: Tx Hash** (lo más confiable). En Binance, copias el Tx Hash de la retirada, lo verificas en Etherscan (blockchain explorer), ves si la dirección es correcta y si el estado es "Success". Luego buscas ese depósito en Coinbase por dirección. Si coincide, emparejado.

**Nivel 2: Heurística** (si no tienes Tx Hash). ¿Misma moneda? ¿Importe ≈ retirada − comisión? ¿Dentro de 24 horas? ¿Cuentas verificadas tuyas? Si todo sí, probablemente es la misma transferencia.

Así resolví 4 transferencias huérfanas en la auditoría. Todas legítimas, todas emparejadas, saldos corregidos.

## Por qué estos patrones importan

Cada uno de estos casos tiene impacto fiscal o de auditoría:

- **BingX:** Pérdidas no documentadas = no deducibles (hasta obtener evidencia). Error potencial: -€150-500 en impuestos si Hacienda lo rechaza.
- **Staking:** Clasificación fiscal equivocada = tramo impositivo incorrecto. Error potencial: overhaul completo del IRPF.
- **Transferencias huérfanas:** Saldo negativo = base de coste incorrecto = ganancias infladas. Error potencial: €0-€5.000+ en impuestos dependiendo del volumen.

**No son "problemas técnicos que se ignoran". Son déficits de auditoría que terminan en errores en la renta.**

## La lección

Cuando construí el agente, asumí que la IA podía encontrar estos patrones con razonamiento puro. Resultó que sí, pero necesitaba:

1. **Dominio específico documentado** — saber qué es staking y cómo se clasifica
2. **Procedimientos claros** — "cómo resolver una transferencia huérfana"
3. **Casos reales verificados** — no asumir, haber visto el patrón en la realidad

La verdad es que estos patrones no son únicos a mi auditoría. Si tienes cripto en múltiples exchanges, es probabilístico que encuentres al menos uno.

Y si no tienes un procedimiento para resolverlo, es fácil pasar por alto o cometer un error que te cueste dinero en Hacienda.

## ¿Y ahora?

Si reconoces alguno de estos patrones en tu cuenta, aquí está lo que tienes que hacer:

**BingX Copy Trading:** Exporta manualmente desde la app. Documenta con screenshots. Contacta al exchange si es necesario.

**Staking:** Cada recompensa es RCM. Valora a fecha de recepción. Pásalo a tu asesor con etiqueta clara ("Rendimiento de Capital").

**Transferencias huérfanas:** Verifica el Tx Hash en blockchain. Si no existe, investiga. Si existe, empareja en CoinTracking manualmente.

Todo esto es manual, un poco tedioso, pero es la diferencia entre "creo que estoy correcto" y "sé que estoy correcto".

Y cuando se trata de números que van a Hacienda, el segundo es lo único que importa.

---

*Y eso cierra la serie de 5 artículos sobre construir un agente auditor de criptos. Si quieres el código, la base de conocimiento, o cualquier detalle técnico, escríbeme.*
