---
id: KB-B4-004
title: "Integración Trust Wallet con CoinTracking"
level: B
domain: cointracking
source: "Trust Wallet official docs + análisis DeFi casos"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-05
confidence: high
version: 1.0

related_adr:
  - ADR-010
  - ADR-003

related_docs:
  - knowledge/wallets/METAMASK_INTEGRATION.md
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md
  - knowledge/cointracking/behavioral/DEFI_SWAPS_MECHANICS.md

tags:
  - wallet
  - trust-wallet
  - hot-wallet
  - mobile
  - defi

notes: "Trust Wallet es hot wallet móvil (Android/iOS) con DEX integrado (Uniswap, Pancakeswap, etc). Propiedad de Binance. Popular en traders activos y mobile-first users."
---

# Integración Trust Wallet

## ¿Qué es Trust Wallet?

**Trust Wallet** es un **hot wallet móvil** (no-custodial):
- App iOS/Android
- Claves privadas almacenadas **localmente** en teléfono
- DEX integrado (Uniswap, PancakeSwap, Curve, etc.)
- Conecta con dApps (Aave, OpenSea, etc.)
- Soporte para 70+ blockchains
- Propiedad de Binance (desde 2018)

**Principal para:** DeFi móvil, trading quick, multi-chain farming

**NO es exchange.** Interactúas directamente con protocolos DeFi.

---

## Redes Soportadas (Principales)

| Red | Blockchain | Uso | Complejidad |
|-----|-----------|-----|-----------|
| **Ethereum** | ETH | Tokens principales, DeFi | Alta (gas caro) |
| **Polygon** | MATIC | DeFi escalado, quick trades | Media |
| **Binance Smart Chain** | BNB | PancakeSwap, Venus | Media |
| **Solana** | SOL | Serum, Orca, Magic Eden | Media |
| **Arbitrum** | ARB | GMX, Camelot, Uniswap V3 | Media |
| **Optimism** | OP | Curve, Aave, Uniswap | Media |
| **Fantom** | FTM | SpookySwap, Aave | Baja |

---

## Operaciones Comunes

### 1. Swap (DEX Integrado)

**Ejemplo: Uniswap V3 en Ethereum (desde Trust Wallet)**

```
Date: 2024-06-15 14:23:00
Action: Swap
From: 1.0 ETH
To: 1,850 USDC
Protocol: Uniswap V3 (integrado en Trust Wallet)
Gas: 0.003 ETH (~6 USD)
Slippage: 0.2%
Status: Confirmed
TX Hash: 0xabcd...
```

**En CoinTracking:**

Aparece como DOS operaciones:
1. **Sell:** 1.0 ETH @ 2,000 USD
2. **Buy:** 1,850 USDC @ 1.00 USD

**Gas** reducido del ETH enviado: (1.0 - 0.003) = 0.997 ETH

**Fiscalidad (España):**
- Ganancia = 1,850 (USDC recibido) - 2,000 (ETH vendido) - 6 (gas)
- En este caso: pérdida de -156 USD

### 2. Liquidity Provider (LP) & Yield Farming

**Ejemplo: PancakeSwap (BSC)**

```
Date: 2024-06-01
Action: Add Liquidity
Amount A: 1.0 CAKE (500 USD)
Amount B: 500 BUSD (500 USD)
LP Token: 1000 CAKE-BUSD LP (certificado)
Gas: 0.001 BNB (~0.30 USD)

Rewards (diario):
Fee Share: 0.01 CAKE + 10 BUSD (cada 24h)
```

**En CoinTracking:**

**Problema:** Trust Wallet NO detecta automáticamente LP en CoinTracking. Requiere:
1. Registrar manualmente Deposit de ambas monedas (o como Swap)
2. Registrar Claim de fees (Income)
3. Registrar Remove Liquidity (cuando sales)

**Alternativa:** Importar dirección en blockchain explorer (mejor).

### 3. Yield Farming (Aave, Compound)

**Ejemplo: Aave Lending en Polygon**

```
Date: 2024-06-01
Action: Deposit
Amount: 10,000 USDC
To: Aave Lending Pool (Polygon)
Receive: aPolUSDC (token de receipt)
APY: 4% annual

Rewards (periódico):
Interest: 38 USDC (0.38% monthly)
Governance Token: 0.5 AAVE (si hay incentivos)
```

**En CoinTracking:**

1. **Deposit:** 10,000 USDC
   - Type: Deposit o Staking
   - Cost Basis: 10,000 USDC

2. **Income:** 38 USDC interest (cada periodo)
   - Type: Interest
   - Fiscalidad: Ingresos del capital
   - Value: 38 USDC @ precio date

3. **Income:** 0.5 AAVE reward
   - Type: Income
   - Fiscalidad: Ingresos del capital
   - Value: 0.5 AAVE @ precio date

4. **Withdraw:** 10,038 USDC + 0.5 AAVE
   - Operación simple (retirada)

---

## Integración con CoinTracking

### Opción 1: Dirección Pública (Recomendado)

**Pasos:**

1. **Trust Wallet:** Account → Copy Address (o QR)
2. **CoinTracking:** Add Account → Blockchain → [moneda]
3. Pegar dirección pública
4. CoinTracking escanea blockchain explorer automáticamente
5. Todas las transacciones importadas

**Ventaja:**
- Automático
- Captura TODA actividad (swaps, LP, farming, bridges)
- Sin riesgo (solo lectura)

**Limitación:**
- Requiere esperar a que CoinTracking parsee transacciones complejas (DEX, LP)

### Opción 2: CSV Manual (NO DISPONIBLE)

**Nota:** Trust Wallet no exporta CSV. Necesitas usar blockchain explorer (Etherscan, PolygonScan, etc.) para exportar.

---

## Casos Límite y Peculiaridades

### 1. Multi-Chain Swaps (Cross-Chain Bridge)

**Ejemplo: Transferir USDC de Ethereum a Polygon via Stargate**

```
Send: 1,000 USDC (Ethereum)
Receive: 1,000 USDC (Polygon) [después de ~15 min]
Gas: 0.002 ETH (Ethereum)
Bridge Fee: 1 USDC
```

**En CoinTracking:**

Aparece como:
1. **Withdraw:** 1,001 USDC (Ethereum) - incluye fee
2. **Deposit:** 1,000 USDC (Polygon) - después de delay

**Problema:** Puede parecer "huérfano" si CoinTracking no conecta automáticamente.

**Solución:** Verificar manualmente que monto enviado - fee = monto recibido.

### 2. Smart Contract Interaction Failures

**Riesgo:** Interactúas con contrato que falla, pero aún pagas gas.

```
Status: Failed
Gas: 0.002 ETH (pagado pero tx rechazado)
Action: Nada (transacción fue reverted)
```

**En CoinTracking:**
- Aparecerá como **Fee** (la pérdida de gas)
- No hay operación correspondiente (swap/LP no ocurrió)
- Es una **pérdida real**

**Fiscalidad:** Probablemente deducible (costo operativo), pero requiere verificación.

### 3. Slippage Handling (Trade Vs Expected)

**Ejemplo:**
```
Esperaba: 1.0 ETH → 2,000 USDC
Recibí: 1.0 ETH → 1,980 USDC (slippage 1%)
```

**En CoinTracking:**
- Sell: 1.0 ETH @ precio spot
- Buy: 1,980 USDC @ precio spot
- Pérdida por slippage: 20 USDC (incluida en el cálculo de ganancia)

### 4. Multiple Token Standards (BEP-20, ERC-20, SPL, etc.)

**Problema:** Mismo símbolo, diferentes blockchains.

```
USDT en BSC (BEP-20): 0x55d398326f99059fF775485246999027B3197955
USDT en Ethereum (ERC-20): 0xdAC17F958D2ee523a2206206994597C13D831ec7
USDT en Solana (SPL): EPjFWdd5Au...
```

**En CoinTracking:** Deberían aparecer como "USDT" por separado según la red.

### 5. NFT & Token Interactions (OpenSea, Magic Eden)

**Trust Wallet soporta compra/venta de NFTs.**

```
Date: 2024-06-15
Action: Buy NFT
NFT: "Cool Cat #1234"
Price: 2.5 ETH
Gas: 0.03 ETH
Marketplace: OpenSea
```

**En CoinTracking:**
- Type: Custom "NFT Purchase"
- Cost Basis: 2.5 ETH (precio de compra)
- Value: 2.5 ETH @ precio spot

**Fiscalidad:** [PENDIENTE DE FUNDAMENTAR] — probablemente ganancia patrimonial al vender (igual que arte físico).

---

## Validación en CoinTracking

### Checklist: Trust Wallet + CoinTracking

```
[ ] Dirección pública agregada (una por blockchain principal)
[ ] Balance actual sincronizado
[ ] Swaps desglosados en venta/compra
[ ] LP tokens visibles
[ ] Yield/interest capturado como Income
[ ] Failed TX aparecen como Fee
[ ] Gas fees incluidos
[ ] Ninguna operación duplicada
[ ] Puentes capturados correctamente
[ ] Saldos no negativos (excepto si es error conocido)
```

### Problemas Comunes

| Problema | Síntoma | Solución |
|----------|---------|----------|
| **Swap no aparece** | TX en Trust Wallet, no en CT | Esperar a que CoinTracking parsee DEX |
| **LP token desaparece** | Balance negativo de LP token | Importar dirección en blockchain explorer |
| **Yield no se captura** | aUSDC balance no actualiza | Verificar que token reward importado |
| **Bridge "huérfano"** | Retiro sin depósito equivalente | Verificar monto - fee = recibido |
| **NFT no aparece** | Compra de NFT sin registro | NFT no criptomoneda; registrar manualmente |
| **Saldo negativo** | Balance < 0 en algún activo | Indicador: datos incompletos o error parsing |

---

## Mejores Prácticas

### ✅ Seguro

- Usar **dirección pública** (solo lectura)
- Verificar TX en **blockchain explorer** (Etherscan, etc.)
- Documentar **interacciones complejas** (LP, farming)
- Incluir **gas spend en decisiones**

### ⚠️ Riesgos

- **Smart Contract Risk:** Código puede tener bugs
- **Slippage:** Recibir menos que esperado
- **Front-running:** Bots pueden adelantarse a tu TX
- **Phone Theft:** Sin passphrase = acceso a fondos perdido
- **Seed Phrase Loss:** Sin recuperación = fondos irrecuperables

### 📋 Para Auditoría

- Documentar todas las direcciones usadas (por blockchain)
- Verificar gas spend vs beneficio de operación
- Incluir failed TX como pérdidas
- Marcar operaciones complejas para revisión fiscal

---

## Fiscalidad (España)

### Resumen

| Operación | Tipo Fiscal |
|-----------|----------|
| **Swap** | Ganancia patrimonial |
| **Deposit a LP/Aave** | Base de coste (no fiscal) |
| **Withdraw de LP/Aave** | Ganancia patrimonial (si subió precio) |
| **Interest/Yield** | Ingresos del capital (Modelo 721) |
| **Governance Token** | Ingresos del capital o regalo (verificar) |
| **Gas Fee** | Incluida en cost basis |
| **Failed TX Gas** | Pérdida (deducible, probablemente) |
| **NFT Purchase** | Ganancia patrimonial al vender [VERIFICAR] |

### Ejemplo Completo

```
2024-01-01: Deposit 10,000 USDC → Aave Polygon
  Cost Basis: 10,000 USDC

2024-06-01: Claim Interest
  Income: 200 USDC (Modelo 721)

2024-12-31: Withdraw 10,200 USDC (+ 5 AAVE governance)
  Ganancia Patrimonial: 10,200 - 10,000 = +200 USD
  Ingreso Capital: 5 AAVE @ precio 31/12

Total Impacto Fiscal:
  Capital Gains: +200 USD
  Income: 200 USD (USDC) + [5 AAVE value]
```

---

## Referencias

- [Trust Wallet Official](https://trustwallet.com/)
- [Trust Wallet Help](https://support.trustwallet.com/)
- [CoinTracking: Import Address](https://www.cointracking.info/en/portfolio_management.php)
- [Etherscan (Ethereum Explorer)](https://etherscan.io/)
- [PolygonScan (Polygon Explorer)](https://polygonscan.com/)

---

**Documento:** Trust Wallet Integration  
**Nivel:** B4-004  
**Status:** Operacional  
**Creado:** 2026-07-05
