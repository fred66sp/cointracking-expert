---
id: KB-B1-XXX
title: "Untitled Document"
level: B
domain: cointracking
source: "Internal documentation"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - todo
  - needs-review

notes: "Metadatos agregados automáticamente. Verificar y actualizar conforme ADR-032."
---

---
id: KB-B4-002
title: "Integración MetaMask con CoinTracking"
level: B
domain: cointracking
source: "MetaMask official docs + análisis DeFi casos"
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
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md
  - knowledge/cointracking/behavioral/DEFI_SWAPS_MECHANICS.md

tags:
  - wallet
  - metamask
  - defi
  - ethereum
  - integration

notes: "MetaMask es hot wallet (no-custodial). Principal para DeFi en Ethereum, Polygon, BSC, etc. Operaciones on-chain complejas (swaps, LP, farming). Requiere importación manual de dirección o explorador blockchain."
---

# Integración MetaMask

## ¿Qué es MetaMask?

**MetaMask** es un **hot wallet** (conexión directa internet):
- Browser extension / Mobile app
- Claves privadas almacenadas **localmente** (no en nube, a diferencia de email login)
- Interfaz para DeFi (conectar a Uniswap, Aave, etc)
- Soporta múltiples redes (Ethereum, Polygon, BSC, Arbitrum, Optimism, etc)

**Principal para:** DeFi, farming, LP (liquidity provider), yield

**NO es exchange.** Interactúas directamente con protocolos DeFi.

---

## Redes Soportadas

| Red | Chain ID | Uso | Complejidad |
|-----|----------|-----|-------------|
| **Ethereum Mainnet** | 1 | Token principales, Aave, Curve | Alta (gas caro) |
| **Polygon** | 137 | DeFi escalado, Uniswap V3 | Media (gas barato) |
| **BSC (Binance)** | 56 | Pancakeswap, Venus | Media |
| **Arbitrum** | 42161 | GMX, Camelot | Media |
| **Optimism** | 10 | Curve, Aave | Media |
| **Base** | 8453 | Uniswap V4, Aave | Baja (nueva red) |

---

## Operaciones Comunes

### 1. Swap (DEX)

**Ejemplo: Uniswap V3 en Ethereum**

```
Date: 2024-06-15 14:23:00
Action: Swap
From: 1.0 ETH
To: 1,850 USDC
Protocol: Uniswap V3
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

### 2. Liquidity Provider (LP)

**Ejemplo: Uniswap V2 Pool ETH-USDC**

```
Date: 2024-06-01
Action: Add Liquidity
Amount A: 1.0 ETH (2,000 USD)
Amount B: 2,000 USDC
LP Token: 2,000 UNI-V2 (certificado)
Gas: 0.005 ETH

Rewards (periódico):
Fee Share: 0.01 ETH + 20 USDC (cada semana)
```

**En CoinTracking:**

**Problema:** CoinTracking NO detecta automáticamente LP. Requiere:
1. Buscar manualmente LP token en blockchain explorer
2. Registrar ENTRADA: Deposit de ambas monedas (o como Swap)
3. Registrar SALIDA: Claim de fees (Income)
4. Registrar RETIRO: Remove Liquidity (cuando sales del pool)

**Alternativa:** Importar dirección en blockchain explorer y CoinTracking importa todo automáticamente (mejor).

### 3. Yield Farming

**Ejemplo: Aave Lending**

```
Date: 2024-06-01
Action: Deposit
Amount: 10,000 USDC
To: Aave Lending Pool
Receive: aUSDC (token de receipt)
APY: 5% annual

Rewards (periódico):
Interest: 500 USDC (0.38% monthly)
Governance Token: 10 AAVE (si hay incentivos)
```

**En CoinTracking:**

1. **Deposit:** 10,000 USDC (base de coste)
   - Type: Deposit o Staking
   - Cost Basis: 10,000 USDC

2. **Income:** 500 USDC interest (cada periodo)
   - Type: Interest
   - Fiscalidad: Ingresos del capital (Modelo 721)
   - Value: 500 USDC @ precio date

3. **Income:** 10 AAVE reward (si aplica)
   - Type: Income
   - Fiscalidad: Ingresos del capital
   - Value: 10 AAVE @ precio date

4. **Withdraw:** 10,500 USDC (principal + interés acumulado)
   - Operación simple (retirada)

---

## Integración con CoinTracking

### Opción 1: Dirección Pública (Recomendado)

**Pasos:**

1. **MetaMask:** Copy Account Address
2. **CoinTracking:** Add Account → Blockchain → Ethereum (o red específica)
3. Pegar dirección
4. CoinTracking escanea blockchain explorer automáticamente
5. Todas las transacciones importadas

**Ventaja:** 
- Automático
- Captura TODA actividad (swaps, LP, farming, bridges)
- Sin riesgo (solo lectura)

**Limitación:**
- Requiere esperar a que CoinTracking parsee transacciones complejas

### Opción 2: CSV Manual

**NO DISPONIBLE:** MetaMask no exporta CSV. Necesitas usar blockchain explorer (Etherscan) para exportar historial.

---

## Casos Límite y Peculiaridades

### 1. Failed Transactions

**Problema:** Una transacción falla (ej. slippage excedido), pero aún pagas gas.

```
Status: Failed
Gas: 0.002 ETH (-6 USD perdidos)
Action: Nada (transacción rechazada)
```

**En CoinTracking:**
- Aparecerá como **Fee** (la pérdida de gas)
- No hay operación correspondiente (swap/LP no ocurrió)
- Es una **pérdida real** (gas quemado)

**Fiscalidad:** ¿Es deducible? Probablemente sí (costo operativo), pero requiere verificación con asesor.

### 2. Smart Contract Bugs

**Riesgo:** Interactúas con contrato que tiene bug y pierdes fondos.

```
Example: Fake Uniswap fork, no-exit scam
Gas: Pagaste
Fondos: Perdidos completamente
TX Status: Success (contrato ejecutó, pero no hizo nada)
```

**En CoinTracking:**
- La transacción aparecerá pero NO capturará los movimientos internos
- Balance se verá negativo temporalmente
- Requiere **investigación manual** en blockchain explorer

**Fiscalidad:** Pérdida total (deducible). Documentar con TX hash y descripción.

### 3. Cross-Chain Bridge

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

**Problema:** Puede parecer "huérfano" si CoinTracking no conecta ambas operaciones automáticamente.

**Solución:** Verificar manualmente que monto enviado - fee = monto recibido.

### 4. Wrapped Tokens

**Ejemplo: Wrapped Ether (wETH)**

```
ETH (nativo) convertido a wETH (token ERC-20)
Operación: Wrap ETH → wETH (DEX como Uniswap)
```

**En CoinTracking:**
- Aparece como Swap: ETH → wETH
- Tasa: 1:1
- Sin diferencia fiscal (solo cambio de formato)

---

## Validación en CoinTracking

### Checklist: MetaMask + CoinTracking

```
[ ] Dirección pública agregada
[ ] Balance actual sincronizado
[ ] Swaps desglosados en venta/compra
[ ] LP tokens visibles
[ ] Yield/interest capturado como Income
[ ] Failed TX aparecen como Fee
[ ] Gas fees incluidos
[ ] Ninguna operación duplicada
[ ] Puentes capturados correctamente
```

### Problemas Comunes

| Problema | Síntoma | Solución |
|----------|---------|----------|
| **Swap no aparece** | TX en MetaMask, no en CT | Esperar a que CoinTracking parsee |
| **LP token desaparece** | Balance negativo de UNI-V2 | Importar dirección en blockchain explorer |
| **Yield no se captura** | aUSDC balance no actualiza | Verificar que aUSDC importada |
| **Bridge "huérfano"** | Retiro sin depósito equivalente | Verificar monto - fee = recibido |
| **Saldo negativo** | Balance < 0 en algunas monedas | Indicador: datos incompletos, revisar blockchain explorer |

---

## Mejores Prácticas

### ✅ Seguro

- Usar dirección pública (solo lectura)
- Verificar TX en blockchain explorer
- Documentar contratos complejos

### ⚠️ Riesgos

- **Smart Contract Risk:** Código puede tener bugs
- **Slippage:** Recibir menos que esperado en swap
- **Front-running:** Bots pueden "adelantarse" a tu TX
- **Lost Keys:** Sin recovery phrase = fondos irrecuperables

### 📋 Para Auditoría

- Documentar todas las direcciones usadas
- Verificar gas spend vs beneficio
- Incluir failed TX como pérdidas
- Marcar operaciones complejas para revisión fiscal

---

## Fiscalidad (España)

### Resumen

| Operación | Tipo Fiscal |
|-----------|-------------|
| **Swap** | Ganancia patrimonial |
| **Deposit a LP/Aave** | Base de coste (no fiscal) |
| **Withdraw de LP/Aave** | Ganancia patrimonial (si subió precio) |
| **Interest/Yield** | Ingresos del capital (Modelo 721) |
| **Governance Token** | Ingresos del capital o regalo (verificar) |
| **Gas Fee** | Incluida en cost basis o deducible |
| **Failed TX Gas** | Pérdida (deducible, probablemente) |

### Ejemplo Completo

```
2024-01-01: Deposit 10,000 USDC → Aave
  Cost Basis: 10,000 USDC

2024-06-01: Claim Interest
  Income: 500 USDC (Modelo 721)

2024-12-31: Withdraw 10,500 USDC (+ 100 AAVE governance)
  Ganancia Patrimonial: 10,500 - 10,000 = +500 USD
  Ingreso Capital: 100 AAVE @ precio 31/12

Total Impacto Fiscal:
  Capital Gains: +500 USD
  Income: 500 USD (USDC) + [100 AAVE value]
```

---

## Referencias

- [MetaMask Official](https://metamask.io/)
- [Etherscan (Ethereum Explorer)](https://etherscan.io/)
- [Defiscan (Polygon Explorer)](https://www.defiscan.io/)
- [CoinTracking: Import Address](https://www.cointracking.info/en/portfolio_management.php)

---

**Documento:** MetaMask Integration  
**Nivel:** B4-002  
**Status:** Operacional  
**Creado:** 2026-07-05
