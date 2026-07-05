---
id: KB-B4-001
title: "Integración Ledger Live con CoinTracking"
level: B
domain: cointracking
source: "Ledger official docs + análisis casos reales"
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
  - knowledge/cointracking/behavioral/BINANCE_SPOT_MECHANICS.md
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md

tags:
  - wallet
  - ledger
  - hardware
  - integration
  - blockchain

notes: "Ledger Live es wallet hardware (cold storage). No es exchange. Operaciones en blockchain (on-chain). Requiere importación manual o dirección pública en CoinTracking."
---

# Integración Ledger Live

## ¿Qué es Ledger Live?

**Ledger** es un **hardware wallet** (cold storage):
- Dispositivo físico (Ledger Nano X, Nano S, Stax)
- Claves privadas NUNCA tocan internet
- Operaciones firma en dispositivo
- Ledger Live = interfaz para ver saldos y hacer transacciones

**NO es exchange.** No hay comisiones de plataforma (solo gas blockchain).

---

## Operaciones Soportadas

### ✅ On-Chain Transactions

| Tipo | Descripción | Ejemplo |
|------|-------------|---------|
| **Send** | Enviar cripto a otra dirección | Transferencia a exchange |
| **Receive** | Recibir cripto desde otra dirección | Depósito desde exchange |
| **Swap** | Cambiar una moneda por otra (DEX) | ETH → USDC en 1inch |
| **Staking** | Bloquear monedas para rewards | Ethereum staking (ETH 2.0) |
| **Claim Rewards** | Recoger rewards de staking | Periodic ETH staking rewards |

### ❌ NO Soportado

- Trading (compra/venta por fiat) → Usa exchange
- Margin/Futures → Usa exchange
- Lending → Usa protocolo DeFi (no Ledger)

---

## Integración con CoinTracking

### Opción 1: Importación Manual (Dirección Pública)

**Pasos:**

1. **Ledger Live:** Account → Copy Address
2. **CoinTracking:** Add Account → Blockchain → [moneda]
3. Pegar dirección pública
4. CoinTracking escanea blockchain automáticamente
5. Transacciones importadas

**Ventaja:** 
- Seguro (no necesita claves)
- Automático (CoinTracking actualiza al escanear)

**Limitación:** 
- Solo cripto on-chain
- No captura saldos internos si usa Ledger Earn

### Opción 2: API Ledger (No Disponible Directamente)

**Nota:** Ledger Live NO expone API pública directamente. CoinTracking importa vía blockchain explorer (Etherscan, etc).

---

## Caso Real: Proyecto `agp`

**Ledger Live en `agp`:**
```
Balance actual:
├─ ETH: 0.15962661 (248.32 EUR)
├─ XRP: 10.00000000 (10.12 EUR)
└─ Total: 258.44 EUR (1.3% cartera)
```

**Transacciones importadas:** ~15 (depósitos, retiradas)

### Validación en `agp`

✅ **Positivos:**
- Balance sincronizado
- Transacciones on-chain verificadas
- Sin duplicados (cada TX blockchain es única)
- Gas fees capturados automáticamente

⚠️ **Puntos de atención:**
- XRP no es EVM chain (requiere importación específica)
- Gas fees en ETH pueden distorsionar cost basis

---

## Operaciones Específicas

### 1. Send (Transferencia Saliente)

**En Ledger Live:**
```
Date: 2024-06-15 14:23:00
Type: Send
From: 0x1234... (tu address)
To: 0xabcd... (exchange/otra wallet)
Amount: 1.5 ETH
Gas: 0.001 ETH
Status: Confirmed
```

**En CoinTracking:**
- Type: **Withdrawal** (si va a exchange)
- Type: **Transfer** (si va a otra wallet tuya)
- Monto: 1.5 ETH
- Fee: 0.001 ETH (incluida automáticamente)

### 2. Receive (Transferencia Entrante)

**En Ledger Live:**
```
Date: 2024-06-15 14:00:00
Type: Receive
From: 0xabcd... (exchange)
To: 0x1234... (tu address)
Amount: 1.5 ETH
Gas: 0 ETH (pagado por remitente)
Status: Confirmed
```

**En CoinTracking:**
- Type: **Deposit** (si viene de exchange)
- Type: **Transfer** (si viene de otra wallet tuya)
- Monto: 1.5 ETH
- Fee: 0 (remitente pagó)

### 3. Swap (On-Chain, Ejemplo 1inch)

**En Ledger Live:**
```
Date: 2024-06-15 14:10:00
Type: Swap
From: 1.0 ETH
To: 1,234.56 USDC
Gas: 0.002 ETH
Protocol: 1inch
Status: Confirmed
```

**En CoinTracking:**
Aparece como DOS operaciones:
1. **Sell ETH:** 1.0 ETH @ ~2000 USD (precio spot)
2. **Buy USDC:** 1,234.56 USDC @ ~1 USD
3. **Fee:** 0.002 ETH (reducida del ETH recibido)

**Fiscalidad:** Ganancia/pérdida = Precio USDC - Precio ETH

### 4. Staking (Ethereum Example)

**En Ledger Live:**
```
Date: 2024-01-01
Type: Stake
Amount: 5.0 ETH
Status: Staking (bloqueado)

Date: Monthly (periódico)
Type: Claim Reward
Amount: 0.05 ETH (reward)
```

**En CoinTracking:**
1. **Deposit:** 5.0 ETH (principal, no es ingreso)
   - Type: Deposit o Staking
   - Cost basis: 5.0 ETH @ precio stake
   
2. **Income:** 0.05 ETH (reward, es ingreso)
   - Type: Income o Reward
   - Fiscalidad: Ingresos del capital (Modelo 721, España)
   - Value: 0.05 ETH @ precio reward date

---

## Validación en CoinTracking

### Checklist: Ledger + CoinTracking

```
[ ] Dirección pública agregada en CoinTracking
[ ] Balance actual sincronizado
[ ] Últimas 10 transacciones visibles
[ ] Gas fees incluidos
[ ] Swap desgllosado en venta/compra
[ ] Staking rewards capturados
[ ] Ninguna operación duplicada
[ ] Saldo negativo de ninguna moneda
```

### Problemas Comunes

| Problema | Síntoma | Solución |
|----------|---------|----------|
| **Dirección no sincroniza** | Balance = 0 en CoinTracking | Verificar dirección correcta, esperar 24h |
| **Gas fee duplicado** | Fee aparece 2x | Verificar en blockchain explorer si es correcta |
| **Swap aparece como venta** | No se captura la compra | Esperar a que CoinTracking parsee DEX |
| **Staking no visible** | Solo aparece el claim, no el lock | Buscar "Stake" o "Lock" en tipo operación |
| **Saldo negativo ETH** | Balance < 0 | Gas fee excedió saldo (error de importación) |

---

## Seguridad y Mejores Prácticas

### ✅ Seguro

- Claves privadas **nunca** se exponen
- CoinTracking ve **solo** dirección pública
- Importación es **solo lectura**

### ⚠️ Riesgo

- **Phishing:** Asegúrate que copias dirección de Ledger Live auténtica
- **Fake Tokens:** No asumas que un token es legítimo solo porque aparece en balance
- **Lost Keys:** Si pierdes dispositivo Ledger sin recovery phrase → fondos irrecuperables

---

## Fiscalidad (España)

### Operaciones On-Chain

| Operación | Fiscalidad |
|-----------|------------|
| **Send/Receive** | Transfer (no fiscal) |
| **Swap** | Ganancia patrimonial (FIFO) |
| **Staking (bloqueo)** | Base de coste (no fiscal hasta) |
| **Staking (reward)** | Ingreso del capital (Modelo 721) |
| **Gas Fee** | Incluida en cost basis |

### Ejemplo: Staking ETH

```
Date: 2024-01-01
Staking 5 ETH @ 2,000 EUR = 10,000 EUR cost basis
No fiscal hasta que retires.

Date: Monthly
Reward 0.05 ETH @ 2,100 EUR = 105 EUR
Ingreso del capital (Modelo 721) [INCIERTO EN ESPAÑA]

Date: 2025-12-31
Retiras 5.05 ETH @ 3,000 EUR = 15,150 EUR
Ganancia = 15,150 - 10,000 = 5,150 EUR (patrimonial)
(+ ingresos del capital acumulados)
```

---

## Referencias

- [Ledger Live Official](https://www.ledger.com/en/ledger-live)
- [CoinTracking: Add Blockchain Address](https://www.cointracking.info/en/portfolio_management.php)
- [Etherscan (Ethereum Explorer)](https://etherscan.io/)

---

**Documento:** Ledger Integration  
**Nivel:** B4-001  
**Status:** Operacional  
**Creado:** 2026-07-05
