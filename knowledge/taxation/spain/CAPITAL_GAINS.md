---
id: KB-A1-002
title: "Ganancias y pérdidas patrimoniales por criptomonedas (IRPF España 2025-2026)"
level: A
domain: taxation
source: "AEAT — Manual práctico IRPF 2025, LIRPF arts. 14, 33, 35, 37.1.h, 46.b, 49 (vigente 2025-2026)"
authority: official
last_verified: 2026-07-05
valid_from: 2025-01-01
valid_until: 2026-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-032
  - ADR-031
  - ADR-028

related_docs:
  - INFORMATIVE_OBLIGATIONS.md
  - knowledge/taxation/spain/CAPITAL_INCOME.md

tags:
  - taxation
  - capital-gains
  - irpf
  - spain
  - 2025

notes: "Ejercicio 2025. Normativa cambia cada año. Requiere reverificación enero 2026 para siguiente ejercicio."
---




# Ganancias y pérdidas patrimoniales por criptomonedas (IRPF España)

> ⚠️ Base técnica para cálculo y reconciliación, **no asesoramiento fiscal**. Ver disclaimer en INDEX.md.

---

## 1. Hecho imponible: qué operaciones generan ganancia/pérdida patrimonial

Según la AEAT, generan **ganancia o pérdida patrimonial** (que se imputa en el momento de la entrega de las monedas):

| Operación | ¿Tributa como ganancia/pérdida patrimonial? | Referencia |
|-----------|---------------------------------------------|------------|
| Venta de cripto por moneda fiduciaria (EUR, USD…) | **Sí** | Art. 35, 14, 46.b LIRPF |
| Permuta de una cripto por otra (cripto-cripto) | **Sí** — es una permuta | Art. 37.1.h LIRPF |
| Pago de bienes/servicios con cripto | **Sí** (transmisión a valor de mercado) | Art. 34-35 LIRPF |

**No** generan ganancia/pérdida patrimonial (no hay alteración patrimonial o hay diferimiento):

- **Compra** de cripto con moneda fiduciaria (solo fija el valor de adquisición; no hay tributación hasta transmitir)
- **Mantener** (holding) sin transmitir
- **Transferencias entre cuentas/monederos del propio titular** (no hay transmisión a un tercero; solo mueven el activo)

> 🔑 **Implicación para la clasificación de eventos:** solo las **transmisiones** (venta, permuta, pago) son eventos fiscales de ganancia/pérdida. Las transferencias internas ya emparejadas en la auditoría (`tools/ct_audit.py`, `knowledge/cointracking/CSV_FORMAT.md` §7) **no** son hechos imponibles — deben excluirse del cálculo de ganancias, aunque sí trasladan el valor de adquisición (coste) entre cuentas.

> ⚠️ Otras rentas (staking, lending, intereses, airdrops, minería) **no** son ganancias patrimoniales por transmisión, sino que se califican como rendimientos u otras rentas. Su tratamiento está en **[CAPITAL_INCOME.md](CAPITAL_INCOME.md)**. Regla clave (CAPITAL_INCOME §7): el valor por el que tributan al percibirse pasa a ser su **coste de adquisición** aquí.

---

## 2. Valoración — venta por moneda fiduciaria

**Ganancia/pérdida = valor de transmisión − valor de adquisición** (Art. 35 LIRPF).

- **Valor de adquisición:** importe real de compra en euros (aplicando el tipo de cambio de la fecha de compra si se pagó en otra divisa) **+ gastos y comisiones** inherentes a la adquisición.
- **Valor de transmisión:** importe real de la venta en euros **− gastos y comisiones** inherentes a la transmisión.

> 🔧 Las **comisiones** (columna `Comisión` del CSV) forman parte de la valoración: suman al coste de adquisición y restan del valor de transmisión. Al preparar la declaración, el agente debe convertir cada comisión a EUR a la fecha correspondiente (ver §5 sobre divisas).

---

## 3. Valoración — permuta cripto-cripto (Art. 37.1.h LIRPF)

El intercambio de una cripto por otra es una **permuta**. La ganancia/pérdida es la diferencia entre:

- el **valor de adquisición** de la moneda entregada, y
- **el mayor** de estos dos: el **valor de mercado** de la moneda entregada **o** el **valor de mercado** de la moneda recibida.

> 🔑 **Consecuencia clave:** en una permuta cripto-cripto **se realiza la ganancia/pérdida acumulada** de la cripto entregada, valorada en EUR a la fecha de la permuta, aunque no se haya pasado por dinero fiat. Es el error fiscal más común de los inversores. El agente debe tratar toda `Operación` cripto↔cripto como evento imponible, no solo las ventas a EUR (ver skill `spanish-tax-return` Paso 2).

> 🔧 Requiere **precio de mercado en EUR** de ambos activos en la fecha/hora de la operación. Esto crea una dependencia de datos de precio histórico (fuente a definir; candidata a spec propia y a conocimiento en `knowledge/`).

---

## 4. Identificación de unidades: método FIFO

En transmisiones parciales de monedas virtuales **homogéneas** adquiridas en distintos momentos y precios, la **doctrina administrativa de la DGT** establece el criterio **FIFO** (*First In, First Out*): se considera que **las transmitidas son las adquiridas en primer lugar**. Fundamento citado: consultas vinculantes **V1604-18, V0975-22, V2520-22, V0648-24 y V0525-25** (28/03/2025 — esta última clasifica las criptomonedas como "bienes inmateriales de naturaleza fungible" y aplica el FIFO por la vía del art. 37.2 LIRPF, no por el concepto de "valores homogéneos") (verificado 2026-07-04; reverificar vigencia según ADR-008 antes de aplicar a un ejercicio distinto).

> ⚠️ **Cita frecuente pero incorrecta:** la consulta **V0999-18** (2018), citada a veces como origen del criterio FIFO, en realidad solo **clasifica las monedas virtuales como activos intangibles** — no se pronuncia sobre FIFO. No usarla como fundamento de este criterio.

> 🔴 **CRÍTICO — el propio FIFO está judicialmente cuestionado (hallazgo 2026-07-04, actualizado 2026-07-04):** el **Tribunal Superior de Justicia del País Vasco** (STSJPV n.º 37/2025 y n.º 41/2025, ambas de 9/01/2025, recurso 75/2024 — es un único litigio, no tres sentencias distintas pese a que algunas fuentes secundarias las citan por separado; aplican la Norma Foral 13/2013 de Bizkaia, de contenido análogo a la LIRPF estatal) ha **rechazado expresamente** que las criptomonedas sean "valores homogéneos" (por reserva de ley: no cumplen los requisitos del art. 47.1 del reglamento foral — emisor único, misma operación financiera, derechos frente a un emisor) y ha declarado que el método FIFO propio de esos valores **"no es aplicable"**, debiendo determinarse la ganancia por el **régimen general** (art. 44.1 del reglamento foral) — en la práctica, el coste real de las unidades efectivamente vendidas.
>
> **Estado procesal (verificado 2026-07-04):** no consta ninguna sentencia del Tribunal Supremo que unifique esta cuestión, ni en territorio foral ni en territorio común, y no se ha podido confirmar públicamente si estas sentencias han sido recurridas en casación. **La DGT no ha cambiado su doctrina**: la consulta V0525-25 (28/03/2025), posterior a las sentencias, mantiene el criterio FIFO sin alteración. Hoy conviven dos planos: la **doctrina administrativa** (DGT, vinculante para la Administración, previsiblemente la que seguirá aplicando la AEAT en territorio común) y una **línea judicial foral** en sentido contrario, sin unificación. Aplicabilidad directa de las sentencias limitada a Bizkaia; **no vincula a la AEAT en territorio común**, pero introduce **incertidumbre genuina**. Como reacción, Bizkaia ha propuesto una norma foral expresa que impone el FIFO desde el 1/1/2025 — señal de que el criterio no estaba tan asentado como parecía. **No presentar el FIFO como criterio jurídicamente cerrado e indiscutido**; advertir de esta línea jurisprudencial al usuario/asesor antes de aplicarlo sin más.
>
> **Precisión importante:** el litigio consistió en que la Inspección regularizó aplicando **FIFO global** (todas las unidades del contribuyente, en cualquier wallet/exchange) mientras los contribuyentes habían aplicado **FIFO por exchange**. El TSJPV **no valida el FIFO por exchange como criterio correcto** — simplemente anula la regularización de la Inspección por falta de cobertura legal para tratar las criptomonedas como "valores homogéneos" en absoluto (con cualquier ámbito). No establece una doctrina alternativa sobre el ámbito del cálculo.

> 🔑 **Homogéneas = mismo activo.** Recuérdese (ver `knowledge/cointracking/CSV_FORMAT.md` §8) que CoinTracking desambigua símbolos repetidos con sufijo (`SOL` vs `SOL2`): son activos **distintos**, cada uno con su propia cola FIFO.

> ✅ **Ámbito del FIFO (asumiendo que se aplique): GLOBAL, no por exchange — resuelto 2026-07-04.** La consulta **V0525-25** (28/03/2025) confirma expresamente que la identificación FIFO se hace considerando **todas las unidades del mismo tipo que posee el contribuyente, con independencia del lugar de custodia** (verificado por triangulación de varias fuentes secundarias independientes citando ese pasaje; no se pudo acceder al texto en el buscador oficial de la DGT por un error de certificado al intentar la conexión directa — **[VERIFICAR]** contra el texto literal cuando sea posible acceder). Es la posición administrativa que la AEAT aplicará en territorio común; coherente con que la Inspección, en el litigio vasco, también partiera de un FIFO global (ver recuadro anterior). Esto **no es lo mismo** que decir que el FIFO en sí esté fuera de duda (ver el recuadro crítico de arriba) — son dos preguntas distintas: *si* se aplica FIFO, y *sobre qué ámbito*.
>
> **Relevancia práctica en CoinTracking (si finalmente se aplica FIFO):** la herramienta tiene una opción llamada **"Depot/Lot separation"** (ver `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` §3.1). **Desactivada (valor por defecto)**, CoinTracking calcula FIFO/LIFO/HIFO de forma **global**, combinando todos los exchanges/wallets — coherente con la doctrina de V0525-25. **Activada**, cada exchange/wallet pasa a ser un "depot" independiente (pensado para otras jurisdicciones, p. ej. EE. UU.), lo que **ya no coincidiría** con el criterio de la DGT. Antes de preparar una declaración, **comprobar en la cuenta del usuario que esta opción está desactivada**.

### Ejemplo (ilustrativo)

```
Compras:   1) 01/02  1.0 BTC a 20.000 € (coste 20.000 €)
           2) 01/05  1.0 BTC a 30.000 € (coste 30.000 €)
Venta:     01/09  1.5 BTC a 40.000 €/BTC  -> valor transmisión 60.000 €

FIFO -> se transmiten:
   1.0 BTC del lote 1 (coste 20.000 €)
   0.5 BTC del lote 2 (coste 15.000 €)
Coste de adquisición aplicado = 35.000 €
Ganancia patrimonial = 60.000 − 35.000 = 25.000 €
Queda en cartera: 0.5 BTC del lote 2 (coste 15.000 €)
```

---

## 5. Divisas y precio en EUR

- Toda valoración se expresa en **EUR**. Operaciones en otras divisas (USDT, USD…) requieren conversión a EUR **a la fecha de cada operación**.
- Coherente con ADR-002, todos los importes se manejan con `Decimal`; coherente con ADR-005, las fechas se normalizan a UTC antes de asociar precios.

> 🔧 **Dependencia:** el cálculo fiscal necesita una fuente de **precios históricos EUR** por activo y fecha. Definir origen (CoinTracking ya incorpora valoraciones; o fuente externa) es una decisión pendiente con impacto en determinismo y reproducibilidad.

---

## 6. Integración en la base imponible del ahorro

Las ganancias/pérdidas por transmisión de criptos son **renta del ahorro** (Art. 46.b LIRPF) e integran la **base imponible del ahorro** (Art. 49).

### Tramos de la base del ahorro — ejercicio 2025

Escala **estatal en su totalidad** (no cedida a CCAA): aplica igual en toda España. Progresiva por tramos.

| Base liquidable del ahorro | Tipo (2025) |
|----------------------------|-------------|
| 0 – 6.000 € | 19 % |
| 6.000 – 50.000 € | 21 % |
| 50.000 – 200.000 € | 23 % |
| 200.000 – 300.000 € | 27 % |
| Más de 300.000 € | **30 %** |

> ⚠️ **Cambio 2025:** el último tramo (> 300.000 €) subió del 28 % (ejercicios anteriores) al **30 %**. El cálculo fiscal debe **versionar los tramos por ejercicio**: no cablear tipos; seleccionarlos según el año fiscal de la operación.

---

## 7. Compensación de pérdidas

Las pérdidas patrimoniales de la base del ahorro se compensan, por este orden (**Art. 49 LIRPF**; verificado 2026-07-04 contra el manual práctico de Renta de la AEAT):

1. **Primero**, con ganancias patrimoniales de la misma base del ahorro (sin límite).
2. Si queda saldo negativo, con el saldo positivo de **rendimientos del capital mobiliario** de la base del ahorro (intereses, staking…), con el límite del **25 %** de ese saldo positivo.
3. Si aún queda saldo negativo, se **arrastra a los 4 ejercicios siguientes**, en el mismo orden; pasado ese plazo sin utilizar, se pierde.

> El 25 % lleva vigente desde el ejercicio **2018** (fue subiendo progresivamente desde el 10 % en 2015 tras la reforma de la Ley 26/2014: 2015 10 %, 2016 15 %, 2017 20 %, 2018 en adelante 25 %). **No ha cambiado para 2025.** Reverificar igualmente contra la fuente (`sede.agenciatributaria.gob.es`, manual práctico de Renta del ejercicio correspondiente) antes de aplicarlo a otro año, por la política de vigencia de ADR-008.

### Regla anti-aplicación de pérdidas por recompra (Art. 33.5.e/f LIRPF)

- **Lo que dice la ley:** no computan las pérdidas si se recompran valores/participaciones homogéneos dentro de **2 meses** (mercados regulados) o **1 año** (no cotizados). Redactada pensando en acciones/fondos, no en criptoactivos.
- **¿Se aplica a criptomonedas?** **No hay consulta DGT vinculante que lo confirme expresamente.** Además, la propia DGT (consulta **V0525-25**, 28/03/2025) niega que las criptomonedas sean "valores homogéneos" en el sentido del reglamento del IRPF — la misma calificación que la regla de recompra exige — lo que debilita (no descarta) su aplicabilidad literal.
- **Práctica operativa de la AEAT:** según fuentes secundarias (no la propia consulta), Renta Web incluye igualmente una casilla específica que aplica esta regla a criptomonedas con el plazo de **12 meses**, pese al vacío doctrinal — es decir, la Administración la aplica en la práctica aunque no esté confirmada por escrito para cripto.
- **[PENDIENTE DE FUNDAMENTAR]** con el texto literal de una consulta DGT o norma que zanje esta contradicción entre lo que dice la ley/V0525-25 y lo que hace Renta Web. No presentar como "regla de 2 meses/1 año confirmada para cripto" sin ese matiz.

---

## 8. Resumen para la skill `spanish-tax-return`

1. Clasificar cada operación: transmisión (venta/permuta/pago) vs no imponible (compra, holding, transferencia interna).
2. Para cada transmisión, aplicar **FIFO** por activo (identidad = ticker completo) para obtener el coste — cálculo determinista vía el Informe de Impuestos de CoinTracking, no por el LLM (ADR-006).
3. Valorar en **EUR** a la fecha (venta: importe real; permuta: regla del Art. 37.1.h).
4. Incluir comisiones en la valoración.
5. Sumar ganancias/pérdidas del ejercicio e integrarlas en la **base del ahorro** con los **tramos del año correspondiente**.
6. Producir detalle trazable por operación (evidencia) en el informe de la skill, no la declaración.

**Cuestiones abiertas:** ámbito del FIFO (global vs por cuenta) §4 — comprobar entretanto la opción "Depot/Lot separation" de CoinTracking; fuente de precios históricos EUR §5; reglas exactas de compensación de pérdidas §7. El tratamiento de staking/lending/airdrops/minería ya está fundamentado en **[CAPITAL_INCOME.md](CAPITAL_INCOME.md)** (con sus propios `[VERIFICAR]`).
