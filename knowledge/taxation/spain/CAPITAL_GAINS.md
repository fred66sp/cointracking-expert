# Ganancias y pérdidas patrimoniales por criptomonedas (IRPF España)

**Jurisdicción:** España — IRPF, personas físicas residentes
**Fuentes:** Manual práctico IRPF 2025 (AEAT); LIRPF arts. 14, 33, 35, 37.1.h, 46.b, 49
**Última verificación:** 2026-07-02
**Vigencia:** cifras y criterios del ejercicio **2025**. La normativa cambia cada año — reverificar (AEAT/BOE/DGT) para otros ejercicios o si esta fecha es antigua (ADR-008).
**Estado:** Fundamentado en fuente oficial

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

> 🔑 **Implicación para el motor fiscal:** solo las **transmisiones** (venta, permuta, pago) son eventos fiscales de ganancia/pérdida. Las transferencias internas emparejadas por el motor de transferencias **no** son hechos imponibles — deben excluirse del cálculo de ganancias, aunque sí trasladan el valor de adquisición (coste) entre cuentas.

> ⚠️ Otras rentas (staking, lending, intereses, airdrops, minería) **no** son ganancias patrimoniales por transmisión, sino que se califican como rendimientos u otras rentas. Su tratamiento está en **[CAPITAL_INCOME.md](CAPITAL_INCOME.md)**. Regla clave (CAPITAL_INCOME §7): el valor por el que tributan al percibirse pasa a ser su **coste de adquisición** aquí.

---

## 2. Valoración — venta por moneda fiduciaria

**Ganancia/pérdida = valor de transmisión − valor de adquisición** (Art. 35 LIRPF).

- **Valor de adquisición:** importe real de compra en euros (aplicando el tipo de cambio de la fecha de compra si se pagó en otra divisa) **+ gastos y comisiones** inherentes a la adquisición.
- **Valor de transmisión:** importe real de la venta en euros **− gastos y comisiones** inherentes a la transmisión.

> 🔧 Las **comisiones** (columna `Comisión` del CSV) forman parte de la valoración: suman al coste de adquisición y restan del valor de transmisión. El motor fiscal debe convertir cada comisión a EUR a la fecha correspondiente (ver §5 sobre divisas).

---

## 3. Valoración — permuta cripto-cripto (Art. 37.1.h LIRPF)

El intercambio de una cripto por otra es una **permuta**. La ganancia/pérdida es la diferencia entre:

- el **valor de adquisición** de la moneda entregada, y
- **el mayor** de estos dos: el **valor de mercado** de la moneda entregada **o** el **valor de mercado** de la moneda recibida.

> 🔑 **Consecuencia clave:** en una permuta cripto-cripto **se realiza la ganancia/pérdida acumulada** de la cripto entregada, valorada en EUR a la fecha de la permuta, aunque no se haya pasado por dinero fiat. Es el error fiscal más común de los inversores. El motor fiscal debe tratar toda `Operación` cripto↔cripto como evento imponible, no solo las ventas a EUR.

> 🔧 Requiere **precio de mercado en EUR** de ambos activos en la fecha/hora de la operación. Esto crea una dependencia de datos de precio histórico (fuente a definir; candidata a spec propia y a conocimiento en `knowledge/`).

---

## 4. Identificación de unidades: método FIFO

En transmisiones parciales de monedas virtuales **homogéneas** adquiridas en distintos momentos y precios, la AEAT establece el criterio **FIFO** (*First In, First Out*): se considera que **las transmitidas son las adquiridas en primer lugar**.

> 🔑 **Homogéneas = mismo activo.** Recuérdese (ver `knowledge/cointracking/CSV_FORMAT.md` §8) que CoinTracking desambigua símbolos repetidos con sufijo (`SOL` vs `SOL2`): son activos **distintos**, cada uno con su propia cola FIFO.

> ❓ **Cuestión abierta:** ¿el FIFO se aplica por **cartera global** del contribuyente (todas las cuentas juntas) o **por exchange/cuenta**? El criterio general de la AEAT apunta al conjunto del mismo activo del contribuyente, con independencia de dónde se custodie. **[PENDIENTE DE FUNDAMENTAR con consulta DGT específica]** antes de cerrar la spec del motor FIFO.

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

Las pérdidas patrimoniales de la base del ahorro se compensan con ganancias de la misma base, y existen reglas de compensación con rendimientos del capital mobiliario y de arrastre a ejercicios siguientes.

> **[PENDIENTE DE FUNDAMENTAR]** los porcentajes y plazos exactos de compensación (Art. 49 LIRPF y su evolución por ejercicio) antes de implementarlos. Existe además una posible **regla anti-aplicación de pérdidas por recompra** de activos homogéneos (Art. 33.5 LIRPF); su aplicabilidad a criptomonedas debe verificarse con doctrina específica. No implementar sin fuente.

---

## 8. Resumen para el motor fiscal

1. Clasificar cada operación: transmisión (venta/permuta/pago) vs no imponible (compra, holding, transferencia interna).
2. Para cada transmisión, aplicar **FIFO** por activo (identidad = ticker completo) para obtener el coste.
3. Valorar en **EUR** a la fecha (venta: importe real; permuta: regla del Art. 37.1.h).
4. Incluir comisiones en la valoración.
5. Sumar ganancias/pérdidas del ejercicio e integrarlas en la **base del ahorro** con los **tramos del año correspondiente**.
6. Producir detalle trazable por operación (evidencia), no la declaración.

**Cuestiones abiertas:** ámbito del FIFO (global vs por cuenta) §4; fuente de precios históricos EUR §5; reglas exactas de compensación de pérdidas §7. El tratamiento de staking/lending/airdrops/minería ya está fundamentado en **[CAPITAL_INCOME.md](CAPITAL_INCOME.md)** (con sus propios `[VERIFICAR]`).
