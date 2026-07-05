---
id: KB-B1-012
title: "Altcoins: Casos Especiales y Trampas Comunes"
level: B
domain: cointracking
source: "CoinTracking + casos reales de auditoría"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-05
confidence: high
version: 1.0

related_adr:
  - ADR-003
  - ADR-010

related_docs:
  - knowledge/cointracking/behavioral/STABLECOINS_AND_SPECIAL_TOKENS.md
  - knowledge/taxation/spain/CAPITAL_GAINS.md

tags:
  - altcoins
  - special-cases
  - governance-tokens
  - inflation
  - delisting

notes: "Altcoins incluyen tokens 'no principales' (no BTC/ETH). Casos especiales: token splits, airdrops de governance, delisting, ticker changes, scams."
---

# Altcoins: Casos Especiales y Trampas

## ¿Qué es un "Altcoin"?

**Altcoin** = cualquier criptomoneda que **no es Bitcoin o Ethereum** (por definición, "Alt" = "Alternative").

Incluye:
- DeFi tokens (AAVE, COMP, UNI, CURVE)
- Layer 2 tokens (ARB, OP)
- L1 alternatives (SOL, AVAX, NEAR)
- Governance tokens (MKR, MAKER)
- Memcoins (DOGE, SHIB)
- Proyectos fallidos o en riesgo

---

## Caso 1: Airdrops y Governance Tokens

### Airdrop de Governance Token

**Ejemplo real: Uniswap UNI (2020)**

```
Usuario: Usa Uniswap, no espera nada
Sorpresa: Recibe 400 UNI airdrop gratis (@ $0.93 = $372)
En CoinTracking: 
  - Tipo: Income o Reward
  - Fecha: Día del airdrop
  - Cantidad: 400 UNI
  - Valor: $372 (a precio de mercado en esa fecha)
  - Fiscalidad: Ingresos del capital (RCM, base general) @ $372
```

**Fiscalidad en España:**
- **Valor de recepción:** Precio de mercado en fecha de airdrop (DGT)
- **Tributación:** Ingresos del capital (Modelo 721, base general)
- **Coste de adquisición:** Ese valor es el cost basis para futuro FIFO
- Si vendes después a $1.50: ganancia patrimonial = (400 × $1.50) - (400 × $0.93) = $228

### Airdrop Falso (Scam)

**Problema:** Algunos "airdrops" son **scams** que piden conectar wallet.

```
Aviso: "Conecta MetaMask para reclamar MOONTOKEN airdrop"
Riesgo: Autorizar contrato malicioso que vacía tu wallet
En CoinTracking: Si alguna vez ves MOONTOKEN con valor = 0, posible scam
```

**Defensa:**
- ✅ Airdrops legítimos aparecen **automáticamente** en wallet (sin acción)
- ❌ Airdrops que piden "conectar" o "aprobar" son casi siempre scams
- ✅ Si lo conectas: revocar el contrato malicioso inmediatamente

---

## Caso 2: Token Splits (Cambio de Escala)

### Ejemplo: SHIB Token Split (2021)

**Antes de split:**
```
Tenías: 1,000,000 SHIB @ 0.00001 USD = 10 USD
Precio unitario: muy bajo
```

**Después de split 1:1,000,000:**
```
Cantidad cambió: 1,000,000,000,000 SHIB (un trillón)
Precio unitario: 0.000000001 USD (proporcional)
Valor total: Igual = 10 USD
```

**Efecto en CoinTracking:**
- Si no importas el split, el balance se ve **"dividido"**
- Cantidad original: 1M SHIB
- Nueva cantidad: 1T SHIB
- Precio: 0.00001 → 0.000000001
- **Resultado:** Ves pérdida de 0.001 USD (por unidad)

**Solución:**
1. Buscar "SHIB token split" en CoinGecko o Twitter
2. Registrar manualmente la fecha del split
3. En CoinTracking: buscar la operación de split o contactar soporte

**Fiscalidad:** Sin efecto (valor total no cambia).

### Ejemplo: Tokens con decimales especiales

```
Algunos tokens no usan 18 decimales (estándar):
- USDC: 6 decimales (1 USDC = 0.000001 token unitario)
- SHIB: 18 decimales
- TrueUSD (TUSD): 18 decimales

Si CoinTracking no detecta bien, ves cantidades masivas o microscópicas.
```

---

## Caso 3: Delisting y Cambio de Ticker

### Delisting de Binance

**Ejemplo: Algunos altcoins son deslistados de Binance (2023-2024)**

```
Token: XYZ_USDT par en Binance
Anuncio: Delisting en 30 días
Usuario: Tiene 1000 XYZ en Binance Spot
Acción: DEBE vender o transferir ANTES del delisting
```

**En CoinTracking:**
- Si no vendes antes: balance se "congela"
- Valor sigue siendo el del último precio conocido (no actualiza)
- **Riesgo:** Puede ser imposible vender si sólo cotizaba en Binance

**Fiscalidad:**
- Si vendes en delisting: ganancia/pérdida al precio de venta
- Si no vendes (pierde valor): pérdida patrimonial (documentable)

### Cambio de Ticker

**Ejemplo: Coin renueva su identidad**

```
Antes: Ticker "ABC" @ $1.50
Después: Nuevo nombre, ticker "XYZ" @ $1.50 (mismo proyecto, rebrand)
En CoinTracking: Pueden aparecer como DOS activos distintos
```

**Problema:** Saldo de ABC desaparece, saldo de XYZ aparece.

**Solución:** 
- Buscar "ABC rebrand to XYZ" en documentación oficial
- Registrar manualmente la conversión (si no fue automática)
- Verificar que CoinTracking detecta la migración

---

## Caso 4: Governance Tokens y Votación

### AAVE (Protocolo Aave)

```
Deposit 1000 USDC en Aave Lending
Receive: aUSDC (token de recepción)
Después: Ganas derecho a AAVE governance token
Reward: 10 AAVE (después de 6 meses)
```

**En CoinTracking:**
- USDC deposit: base de coste 1000 USDC
- AAVE reward: ingresos del capital @ precio AAVE en fecha recepción
- aUSDC balance: aparece como activo (realmente es derecho a USDC + interés)

**Fiscalidad (España):**
- USDC: Base de coste (1000 USDC)
- AAVE: Ingreso del capital (RCM, base general)
- aUSDC: ¿Es tributario al recibirlo? [PENDIENTE DE FUNDAMENTAR — posible: no, es solo recepción; o sí, es token distinto]

### MakerDAO (MKR)

```
Deposita 10 ETH como colateral en Maker
Genera: 5000 DAI stablecoin
Fee: 6% anual en DAI
Gobernanza: Tenedor de MKR vota sobre parámetros del protocolo
```

**Fiscalidad:**
- ETH depositado: Base de coste 10 ETH
- DAI generado: Deuda (no es ingreso)
- Fee de 6%: Interés sobre deuda (deducible, probablemente)

---

## Caso 5: Tokens Inflacionarios vs Deflacionarios

### Token Inflacionario (Ejemplo: Typical Earn Token)

```
Deposit 100 TOKEN @ 1 USD = 100 USD
APY: 50% (muy alto, común en protocolo nuevo)
Después 1 año: 150 TOKEN @ ??? USD

¿Cuál es el precio de 150 TOKEN después de 1 año?
- Si protocolo crece: Puede ser $1.50 (ganancia)
- Si protocolo muere: Puede ser $0.01 (pérdida)
```

**Riesgo:** Tokens con APY extremadamente altos suelen tener inflación extrema (precio cae).

### Token Deflacionario (Ejemplo: Algunos Memecoins)

Algunos tokens **queman supply** para subir el precio unitario.

```
Inicio: 1 trillón de tokens
Quema: 10% cada año
Año 1: 900 mil millones de tokens (supply desciende)
Efecto: Menos tokens = precio sube (teóricamente)
```

**En CoinTracking:** Sin efecto directo (la cantidad que tengas sigue siendo la misma).

---

## Caso 6: Tokens Fake o Scam

### Fake Token (Impersonación)

```
Real: Uniswap (UNI), oficial, dirección 0x1f9840a85d5af5bf1d1762f925bdaddc4201f984
Fake: Uniswap (UNI), FALSO, dirección 0xabc123...

Usuario compra TOKEN FALSO por error en DEX
Valor: 0 (scam)
```

**En CoinTracking:**
- Si importas dirección y contrato es falso: balance será 0
- Si compras en exchange donde el par existe (raro): CoinTracking mostrará "compra"

**Defensa:**
- Verificar dirección del contrato en CoinGecko / Etherscan ANTES de comprar
- Usar solo pares listados en exchanges establecidos (Uniswap, Curve, etc.)

### Rug Pull

```
Proyecto: Nuevo DEX/Yield Farm "MegaYield"
APY: 1000% (ridículo, pero atrae)
Usuario: Deposita 100,000 USDC
Fundador: Retira todo el dinero (rug pull)
Usuario: Pierde 100,000 USDC
```

**En CoinTracking:**
- Depósito: 100,000 USDC (base de coste)
- Pérdida: -100,000 USDC (deducible)
- Registrar: como Fee o Loss (documentar con evidencia)

---

## Caso 7: Tokens de Layer 2 (Arbitrum, Optimism)

### ARB (Arbitrum)

```
Usuario: Opera en Arbitrum
Sorpresa: Recibe ARB airdrop (governance token)
Cantidad: Basada en actividad histórica en la red
Fiscalidad: Ingresos del capital (RCM)
```

### OP (Optimism)

```
Mismo mecanismo que ARB
Airdrop basado en actividad en Optimism
Fiscalidad: Ingresos del capital
```

**En CoinTracking:**
- Importar como Income/Reward
- Valor = precio de mercado en fecha de airdrop
- Cost basis para futuro FIFO = ese valor

---

## Checklist: Altcoin en CoinTracking

```
[ ] Ticker verificado contra CoinGecko (no es fake)
[ ] Precio en rango esperado (no 0 o infinito)
[ ] Cantidad verificada contra wallet/exchange
[ ] Si hay airdrop: registrado como Income
[ ] Si hay token split: cantidad actualizada
[ ] Si es wrapped/migracion: documentado en comentario
[ ] Saldo final > 0 (sin negativos sin explicación)
[ ] Dirección del contrato verificada (en Etherscan)
```

---

## Mejores Prácticas

### ✅ Hacer

1. **Verificar cualquier token nuevo** en CoinGecko ANTES de importar
2. **Documentar airdrops** como Income (con fecha + cantidad)
3. **Rastrear governance tokens** por separado (AAVE, UNI, MKR, etc.)
4. **Buscar historial de token** (splits, cambios, delisting)
5. **Registrar dirección del contrato** en comentario si es dudoso

### ❌ No hacer

1. ❌ **Asumir que todo token con precio = 0 es error**
2. ❌ **Ignorar airdrops** (son ingresos del capital, tributan)
3. ❌ **Mezclar versiones antiguas/nuevas** del mismo proyecto
4. ❌ **Comprar airdrops "reclamables"** (probables scams)
5. ❌ **Confiar en APY extremo** sin verificar proyecto

---

## Referencias

- [CoinGecko (búsqueda de tokens)](https://www.coingecko.com/)
- [Etherscan (verificar contratos)](https://etherscan.io/)
- [Defi Pulse (auditoría de protocolos)](https://defipulse.com/)
- [Rekt News (incidents/scams)](https://rekt.news/)

---

**Documento:** Altcoins Special Cases  
**Nivel:** B1-012  
**Status:** Operacional  
**Creado:** 2026-07-05
