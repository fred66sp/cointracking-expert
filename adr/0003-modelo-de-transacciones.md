---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-003: Modelo de transacciones

**Status:** Accepted

**Date:** 2026-07-04

## Context

El agente audita operaciones de criptomonedas. CoinTracking soporta múltiples tipos de transacciones, pero **no todas tienen el mismo peso fiscal ni el mismo tratamiento contable**. Sin un modelo claro:

1. El agente podría confundir "Reward" con "Income" (impactos fiscales distintos)
2. Una "Transfer" entre exchanges propios se auditaría igual que una venta
3. Comisiones se tratarían como operaciones en lugar de ajustes de coste
4. Operaciones de DeFi (Staking, LP Tokens) no tendrían reglas claras

Un **modelo de transacciones** es el esquema que define:
- Qué tipos existen
- Qué campos son obligatorios
- Qué validaciones aplican
- Cómo se relacionan fiscalmente

## Decision

Se establece un modelo canónico de transacciones. Cada tipo tiene:
- **Definición:** qué es
- **Campos obligatorios:** qué datos debe tener
- **Validaciones:** qué debe cumplirse
- **Tratamiento fiscal:** cómo afecta a ganancias/impuestos
- **Ejemplo real**

### Tipos fundamentales (todas las auditorías)

#### 1. Buy (Compra)

**Definición:** Adquirir un activo a cambio de dinero fiat o crypto.

**Campos obligatorios:**
- `Date` (fecha exacta, con hora)
- `In Asset` (qué compras: BTC, ETH, etc.)
- `In Amount` (cuánto)
- `Cost Per Unit` (precio por unidad)
- `Fee` (comisión en la moneda base o en el activo)
- `Exchange` (dónde compraste)

**Validaciones:**
- Cost Per Unit > 0
- Amount > 0
- Fee ≥ 0
- Date dentro del rango de actividad del usuario

**Tratamiento fiscal:**
- Entra en el pool de coste FIFO
- Base de coste = (Amount × Cost Per Unit) + Fee
- Fecha de adquisición para el cálculo de ganancias patrimoniales

**Ejemplo:**
```
2023-06-15 10:32:15 | Buy | Binance | 1 BTC @ 30.000€ + 10€ fee = 30.010€ base coste
```

---

#### 2. Sell (Venta)

**Definición:** Vender un activo crypto por dinero fiat o recibir otro crypto.

**Campos obligatorios:**
- `Date`
- `Out Asset` (qué vendes)
- `Out Amount`
- `Cost` (cuánto recibes por ello)
- `Fee` (comisión)
- `Exchange`

**Validaciones:**
- Out Amount > 0 (no vendes 0)
- Out Amount ≤ Balance actual del activo (no vendes lo que no tienes)
- Cost ≥ 0 (vendiste por dinero, no regalaste)
- Existe una compra previa con saldo suficiente (para FIFO)

**Tratamiento fiscal:**
- Calcula ganancia/pérdida usando FIFO
- Ganancia = Cost − (Base coste FIFO + Fee)
- Es una ganancia patrimonial, afecta al IRPF

**Ejemplo:**
```
2024-02-20 14:15:00 | Sell | Kraken | 0.5 BTC @ 45.000€ (22.500€ recibidos) − 50€ fee = +15.490€ ganancia (si base coste era 7.010€)
```

---

#### 3. Transfer (Transferencia)

**Definición:** Mover un activo entre dos direcciones que controlas (exchange → exchange, exchange → billetera, billetera → billetera).

**Campos obligatorios:**
- `Date`
- `Asset`
- `Amount`
- `From` (exchange/billetera origen)
- `To` (exchange/billetera destino)
- `Fee` (comisión de red o exchange, si aplica)

**Validaciones:**
- Amount > 0
- From ≠ To
- Existe una dirección conocida en ambos extremos
- La dirección `To` es controlada por el usuario (no envíes a un tercero)

**Tratamiento fiscal:**
- **No es una operación fiscal.** No afecta a ganancias ni a IRPF
- Es un movimiento de custodia
- Fee reduce el saldo (p. ej. BTC enviados − comisión de red = BTC recibido)

**Ejemplo:**
```
2023-09-10 22:40:00 | Transfer | Binance → MetaMask | 2 BTC − 0.0005 BTC fee = 1.9995 BTC recibido
```

---

#### 4. Deposit (Depósito)

**Definición:** Recibir activo en tu cuenta desde una fuente externa (banco fiat, otro exchange, billetera, airdrop, etc.).

**Campos obligatorios:**
- `Date`
- `Asset`
- `Amount`
- `To` (tu cuenta/exchange)
- `From` (origen: banco SEPA, exchange anterior, wallet externa, etc.)

**Validaciones:**
- Amount > 0
- `To` es controlada por ti
- Si es fiat (EUR, USD), debe tener corresponsal en banco
- Si es crypto, debe tener un origen verificable o estar documentado como `[PENDIENTE]`

**Tratamiento fiscal:**
- **No es operación fiscal** (es entrada de capital)
- Pero establece el origen de fondos (importante para la Ley de Blanqueo)
- Si es airdrop, puede ser ingreso (ADR-030, pendiente)

**Ejemplo:**
```
2023-01-15 09:00:00 | Deposit | Bank SEPA → Binance | 10.000€
2023-06-20 14:30:00 | Deposit | Wallet (airdrop) → Ethereum | 100 FLOKI
```

---

#### 5. Withdrawal (Retirada)

**Definición:** Enviar activo desde tu cuenta a una dirección externa (banco fiat, otro exchange, billetera, etc.).

**Campos obligatorios:**
- `Date`
- `Asset`
- `Amount`
- `From` (tu cuenta)
- `To` (destino)
- `Fee` (comisión de red o exchange)

**Validaciones:**
- Amount > 0
- Amount + Fee ≤ Balance actual
- `From` es controlada por ti
- `To` es conocida (banco, otra wallet, etc.)

**Tratamiento fiscal:**
- **No es operación fiscal** (es retiro de capital)
- Reduce el saldo
- El fee afecta el coste realizado

**Ejemplo:**
```
2024-03-05 16:20:00 | Withdrawal | Kraken → Bank SEPA | 5.000€ − 15€ fee = 4.985€ recibido
```

---

### Tipos avanzados (requieren validación especial)

#### 6. Staking / Rewards (ingresos por validación)

**Definición:** Recibir remuneración por bloquear crypto o validar en red.

**Campos obligatorios:**
- `Date`
- `Asset` (la moneda que stackeas)
- `Reward Amount` (cuánto ganas)
- `Staking Pool / Exchange` (Lido, Binance Staking, Ethereum, etc.)

**Validaciones:**
- Reward Amount > 0
- Existe evidencia: pool, transacción en blockchain, o confirmación de exchange

**Tratamiento fiscal:**
- Es **ingreso del ejercicio** (rendimiento del capital)
- Valora al precio de mercado en el momento de recepción
- Entra como base del ahorro (sección C del IRPF, tramo específico)
- `[PENDIENTE FUNDAMENTAR]` en ADR-030 (fiscalidad staking)

**Ejemplo:**
```
2024-01-10 | Staking Reward | Lido (ETH) | +0.05 ETH @ 2.500€/ETH = 125€ ingreso
```

---

#### 7. Airdrop (distribución gratuita)

**Definición:** Recibir activo de forma gratuita por poseer otra moneda o participar en una red.

**Campos obligatorios:**
- `Date`
- `Asset` (nueva moneda)
- `Amount`
- `Source` (qué fork, qué campaña, qué red)

**Validaciones:**
- Amount > 0
- Existe un origen documentable (blockchain transaction, exchange notification, etc.)

**Tratamiento fiscal:**
- Posiblemente **ingreso del ejercicio** (depende de Hacienda y del tipo de airdrop)
- `[PENDIENTE FUNDAMENTAR]` en ADR-030

**Ejemplo:**
```
2024-02-15 | Airdrop | Base Chain drop | +100 BASE tokens
```

---

#### 8. Fee / Commission (comisión o coste de transacción)

**Definición:** Costo pagado por ejecutar una operación (no es una operación en sí, es un ajuste).

**Campos obligatorios:**
- `Date`
- `Asset` (en qué se paga: EUR, BTC, USDT, etc.)
- `Amount`
- `Related Transaction` (la operación a la que afecta)

**Validaciones:**
- Amount > 0
- La moneda de la fee debe estar disponible
- No puede haber fee sin una operación asociada

**Tratamiento fiscal:**
- **No es operación fiscal** (es costo)
- Se suma a la base de coste de la operación principal o reduce el ingreso
- Ejemplo: si compras 1 BTC a 30.000€ + 10€ de fee, el coste base es 30.010€

**Nota:** En CoinTracking, las comisiones generalmente van **incluidas** en la operación (Buy tiene Fee field). Este tipo es para cuando están **separadas**.

---

#### 9. Convert / Swap (intercambio de cripto por cripto)

**Definición:** Cambiar un activo por otro criptográficamente (p. ej. USDT → BTC en un DEX o exchange).

**Campos obligatorios:**
- `Date`
- `Out Asset` (das)
- `Out Amount`
- `In Asset` (recibes)
- `In Amount`
- `Fee` (si aplica)
- `Platform` (Uniswap, Kraken Spot, Binance, etc.)

**Validaciones:**
- Out Amount > 0 e In Amount > 0
- Ambos activos existen
- Out Amount ≤ Balance actual

**Tratamiento fiscal:**
- Se trata como una **venta + compra en dos transacciones**:
  - Venta de Out Asset al precio de mercado en ese momento
  - Compra de In Asset al precio de mercado en ese momento
- Ambas pueden generar ganancias/pérdidas

**Ejemplo:**
```
2024-01-20 15:00:00 | Convert (Uniswap) | Venden 1.000 USDT (1€/USDT) = 1.000€ + Compran 0.025 BTC @ 40.000€ = 1.000€
→ Compra: 0.025 BTC @ 40.000€, venta: 1.000 USDT @ 1€/USDT
```

---

#### 10. Futures / Margin (operaciones apalancadas)

**Definición:** Operaciones con apalancamiento (futures, margin trading, short selling).

**Campos obligatorios:**
- `Date (Open)` y `Date (Close)`
- `Asset` (subyacente)
- `Leverage` (2x, 5x, 10x, etc.)
- `Entry Price` y `Exit Price`
- `PnL` (ganancia/pérdida realizada)
- `Funding Fee` (si aplica, costo de llevar abierta la posición)

**Validaciones:**
- Leverage ≥ 1
- Entry Price > 0 y Exit Price > 0
- PnL calculable = (Exit − Entry) × Amount × Leverage − Funding Fees
- Date Close ≥ Date Open

**Tratamiento fiscal:**
- Son operaciones especulativas
- Pueden ser ganancias patrimoniales, pero hay debate sobre si son ingresos o si cotizan diferente
- `[PENDIENTE FUNDAMENTAR]` en ADR-030

**Nota:** Futures **nunca se mezclan** con operaciones spot sin una regla clara (ADR-???). Requieren análisis separado.

---

#### Otros tipos (parcialmente soportados)

- **Lending / Yield Farming:** Prestar crypto para ganar intereses. `[PENDIENTE FUNDAMENTAR]`
- **NFT / Colectibles:** Compra/venta de NFTs. Tratamiento fiscal unclear.
- **Permuta / Exchange:** Intercambiar crypto con un tercero no en exchange (p. ej. vender 1 BTC por 100.000€ en privado). Requiere documentación manual.
- **Herencia / Donación:** Recibir crypto sin pagar. Implicaciones fiscales especiales.

---

## Consequences

**Positive:**

- **Consistencia:** Todos los ADRs y skills usan la misma definición de "Buy", "Transfer", etc.
- **Validación automática:** El agente puede detectar operaciones malformadas (Buy sin precio, Transfer entre direcciones ajenas, etc.)
- **Trazabilidad fiscal:** Cada tipo tiene claro cómo afecta a impuestos
- **Escalabilidad:** Si aparece un nuevo tipo (p. ej. Synthetic Positions), hay un lugar para documentarlo

**Negative:**

- **Rigidez inicial:** Si el usuario tiene operaciones raras que no caben en estos tipos, hay que extender el modelo
- **Complejidad:** Tipos avanzados (Futures, Staking, Lending) son complejos fiscalmente; el agente no puede resolverlos solo
- **Documentación dual:** Este ADR + `knowledge/cointracking/CATALOG.md` pueden desincronizarse

## Notes

### Relación con ADRs existentes

- **ADR-001:** Convenciones — este es el modelo de datos del proyecto
- **ADR-002:** Fuente de verdad — cada transacción debe validarse contra esta jerarquía
- **ADR-004:** Reconciliación — usa este modelo para validar coherencia
- **ADR-017:** Diagnóstico en orden fijo — sigue este modelo al clasificar operaciones
- **Futuro ADR-030:** Fiscalidad de cada tipo (Staking, Airdrop, Futures, etc.)

### Fuentes

- CoinTracking soporta estos tipos (más otros 10+); ver `knowledge/cointracking/CATALOG.md`
- La blockchain distingue entre transacciones tipo: transfer, contract call, etc. Mapeo a este modelo: no perfecto, depende del tipo de activo
- Fiscalidad española (AEAT): cada tipo tiene tratamiento distinto según el Modelo 721, tramos IRPF, etc.

### Pendientes

- **[PENDIENTE]** Definir `Convert` como operación dual (Sell + Buy) en el código del auditor
- **[PENDIENTE]** Documentar el mapeo entre tipos de CoinTracking y tipos del modelo (algunos overlaps, algunas omisiones)
- **[PENDIENTE]** Crear ADR-030 (Fiscalidad de cada tipo: Staking, Airdrop, Futures, Lending, etc.)
- **[PENDIENTE]** Definir cómo manejar operaciones "hibridas" (p. ej. "Deposit" que es un airdrop, o "Transfer" que es una venta disfrazada)
