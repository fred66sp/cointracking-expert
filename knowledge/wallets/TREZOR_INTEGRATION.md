---
id: KB-B4-003
title: "Integración Trezor con CoinTracking"
level: B
domain: cointracking
source: "Trezor official docs + análisis casos reales"
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
  - knowledge/wallets/LEDGER_INTEGRATION.md
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md

tags:
  - wallet
  - trezor
  - hardware
  - integration
  - blockchain

notes: "Trezor es hardware wallet competidor directo de Ledger. Open-source, buena interfaz web, soporte para 1500+ coins. Integración CoinTracking similar a Ledger (por dirección pública)."
---

# Integración Trezor

## ¿Qué es Trezor?

**Trezor** es un **hardware wallet** (cold storage) similar a Ledger:
- Dispositivo físico (Trezor One, Trezor Model T, Trezor Safe 3)
- Claves privadas NUNCA tocan internet
- Operaciones firma en dispositivo
- Trezor Suite = interfaz web para ver saldos y hacer transacciones

**Diferencias vs Ledger:**
- ✅ **Open-source** (código público, verificable)
- ✅ **Interfaz web más simple** (menos bloatware)
- ✅ **Passphrase support** (protección adicional)
- ⚠️ **Menos rápido que Ledger** en actualizaciones de firmware
- ⚠️ **Menor penetración en institucionales**

**NO es exchange.** No hay comisiones de plataforma (solo gas blockchain).

---

## Operaciones Soportadas

### ✅ On-Chain Transactions

| Tipo | Descripción | Ejemplo |
|------|-------------|---------|
| **Send** | Enviar cripto a otra dirección | Transferencia a exchange |
| **Receive** | Recibir cripto desde otra dirección | Depósito desde exchange |
| **Swap** | Cambiar una moneda por otra (DEX) | ETH → USDC en Uniswap |
| **Staking** | Bloquear monedas para rewards | Ethereum staking (ETH 2.0) |
| **Claim Rewards** | Recoger rewards de staking | Periodic rewards |

### ❌ NO Soportado

- Trading (compra/venta por fiat) → Usa exchange
- Margin/Futures → Usa exchange
- Lending → Usa protocolo DeFi (no Trezor)

---

## Integración con CoinTracking

### Opción 1: Importación Manual (Dirección Pública) — RECOMENDADO

**Pasos:**

1. **Trezor Suite:** Account → Copy Address
2. **CoinTracking:** Add Account → Blockchain → [moneda]
3. Pegar dirección pública
4. CoinTracking escanea blockchain automáticamente
5. Transacciones importadas

**Ventaja:**
- Seguro (no necesita claves)
- Automático (CoinTracking actualiza al escanear)
- Soporta múltiples redes

**Limitación:**
- Solo cripto on-chain
- No captura saldos internos si usa Trezor Earn (si aplica)

### Opción 2: API / Conexión Directa (NO DISPONIBLE)

**Nota:** Trezor Suite NO expone API pública. CoinTracking importa vía blockchain explorer (Etherscan, etc.), igual que Ledger.

---

## Dispositivos Trezor Soportados

| Dispositivo | Año | Características |
|---|---|---|
| **Trezor One** | 2013 | Pantalla pequeña, presupuesto |
| **Trezor Model T** | 2018 | Pantalla color, USB-C, mucha RAM |
| **Trezor Safe 3** | 2023 | Mejor seguridad, display mejorado |

**Para auditoría:** El modelo no importa (la dirección es la misma en todos).

---

## Caso Real: Trezor en Auditoría

**Ejemplo de uso (similar a Ledger):**
```
Balance actual en Trezor:
├─ BTC: 0.05 (€2,000)
├─ ETH: 1.5 (€3,500)
├─ USDC: 5,000 (€5,000)
└─ Total: €10,500 (5-10% de cartera típica)
```

**Importación en CoinTracking:**
1. Copiar dirección de cada activo
2. Agregar como Blockchain Address
3. CoinTracking escanea y captura historial

**Validación:**
- ✅ Balance sincronizado
- ✅ Transacciones on-chain verificadas
- ✅ Gas fees capturados automáticamente
- ✅ Sin duplicados (cada TX blockchain es única)

---

## Operaciones Específicas

### 1. Send (Transferencia Saliente)

**En Trezor Suite:**
```
Date: 2024-06-15 14:23:00
Type: Send
From: xpub123... (tu address derivada)
To: 3J7xK9... (exchange/otra wallet)
Amount: 0.5 BTC
Gas: 0.001 BTC
Status: Confirmed
```

**En CoinTracking:**
- Type: **Withdrawal** (si va a exchange)
- Type: **Transfer** (si va a otra wallet tuya)
- Monto: 0.5 BTC
- Fee: 0.001 BTC (incluida automáticamente)

### 2. Receive (Transferencia Entrante)

**En Trezor Suite:**
```
Date: 2024-06-15 14:00:00
Type: Receive
From: 1A1Z7aD... (exchange)
To: xpub456... (tu address)
Amount: 0.5 BTC
Gas: 0 BTC (pagado por remitente)
Status: Confirmed
```

**En CoinTracking:**
- Type: **Deposit** (si viene de exchange)
- Type: **Transfer** (si viene de otra wallet tuya)
- Monto: 0.5 BTC
- Fee: 0 (remitente pagó)

### 3. Swap On-Chain (1inch, Uniswap)

**Ejemplo (Ethereum):**
```
Date: 2024-06-15 14:10:00
Type: Swap
From: 1.0 ETH
To: 1,850 USDC
Gas: 0.003 ETH
Protocol: 1inch
```

**En CoinTracking:**
- Operación 1: **Sell ETH** 1.0 @ precio spot
- Operación 2: **Buy USDC** 1,850 @ precio spot
- Fee: 0.003 ETH (reducida del ETH enviado)

**Fiscalidad:** Ganancia/pérdida = precio USDC recibido - precio ETH vendido - gas

### 4. Staking (Ethereum 2.0)

**En Trezor Suite:**
```
Date: 2024-01-01
Type: Stake
Amount: 2.0 ETH
Status: Staking
APY: 3-4%

Date: Monthly (periódico)
Type: Claim Reward
Amount: 0.006 ETH (aproximado)
```

**En CoinTracking:**
1. **Deposit:** 2.0 ETH
   - Type: Deposit o Staking
   - Cost basis: 2.0 ETH @ precio stake
   
2. **Income:** 0.006 ETH (rewards)
   - Type: Income o Reward
   - Fiscalidad: Ingresos del capital (Modelo 721, España)
   - Value: 0.006 ETH @ precio reward date

---

## Validación en CoinTracking

### Checklist: Trezor + CoinTracking

```
[ ] Dirección pública agregada en CoinTracking
[ ] Balance actual sincronizado
[ ] Últimas 10 transacciones visibles
[ ] Gas fees incluidos
[ ] Swap desgllosado en venta/compra
[ ] Staking rewards capturados
[ ] Ninguna operación duplicada
[ ] Saldo negativo de ningún activo
```

### Problemas Comunes

| Problema | Síntoma | Solución |
|----------|---------|----------|
| **Dirección no sincroniza** | Balance = 0 en CoinTracking | Verificar dirección correcta, esperar 24h |
| **Gas fee duplicado** | Fee aparece 2x | Verificar en blockchain explorer si es correcta |
| **Swap aparece como venta** | No se captura la compra | Esperar a que CoinTracking parsee DEX |
| **Staking no visible** | Solo aparece el claim, no el lock | Buscar "Stake" en tipo operación |
| **Saldo negativo** | Balance < 0 en algún activo | Gas fee excedió saldo (error de importación) |
| **Múltiples direcciones** | Importé 10 direcciones del mismo dispositivo | Combinar en una (CoinTracking soporta consolidación) |

---

## Diferencias vs Ledger Live

| Aspecto | Ledger Live | Trezor Suite |
|---------|-----------|---|
| **Interfaz** | App escritorio + web | Web (trezor.io) |
| **Open Source** | No | Sí |
| **Passphrase** | Soportado | Soportado |
| **Actualizaciones** | Frecuentes | Menos frecuentes |
| **Soporte tokens** | 5,500+ | 1,500+ |
| **Penetración** | Mayor (institucional) | Menor |
| **Costo dispositivo** | €79 (Nano S Plus) | €59 (One) / €149 (Safe 3) |
| **CoinTracking soporte** | Igual (por address) | Igual (por address) |

---

## Seguridad y Mejores Prácticas

### ✅ Seguro

- Claves privadas **nunca** se exponen
- CoinTracking ve **solo** dirección pública
- Importación es **solo lectura**
- Open-source (puedes auditar el código)

### ⚠️ Riesgo

- **Phishing:** Asegúrate que copias dirección de Trezor Suite auténtica (https://trezor.io/)
- **Fake Tokens:** No asumas que un token es legítimo solo por aparecer en balance
- **Lost Device:** Sin recovery phrase = fondos irrecuperables
- **Passphrase Risk:** Si olvidas passphrase = acceso a fondos perdido (aunque device funcione)

---

## Fiscalidad (España)

### Operaciones On-Chain

| Operación | Fiscalidad |
|-----------|----------|
| **Send/Receive** | Transfer (no fiscal) |
| **Swap** | Ganancia patrimonial (FIFO) |
| **Staking (bloqueo)** | Base de coste (no fiscal hasta retiro) |
| **Staking (reward)** | Ingreso del capital (Modelo 721) |
| **Gas Fee** | Incluida en cost basis |

### Ejemplo: Staking ETH en Trezor

```
2024-01-01: Staking 2 ETH @ 2,000 EUR = 4,000 EUR cost basis
Sin fiscal hasta que retires.

2024-03-01: Reward 0.006 ETH @ 2,100 EUR = 12.60 EUR
Ingreso del capital (Modelo 721) [INCIERTO SI APLICA]

2025-12-31: Retiras 2.006 ETH @ 3,000 EUR = 6,018 EUR
Ganancia = 6,018 - 4,000 = 2,018 EUR (patrimonial)
(+ ingresos del capital acumulados)
```

---

## Referencias

- [Trezor Suite (oficial)](https://suite.trezor.io/)
- [Trezor Help (support)](https://trezor.io/support/)
- [CoinTracking: Add Blockchain Address](https://www.cointracking.info/en/portfolio_management.php)
- [Etherscan (Ethereum Explorer)](https://etherscan.io/)
- [Trezor Source Code (GitHub)](https://github.com/trezor/)

---

**Documento:** Trezor Integration  
**Nivel:** B4-003  
**Status:** Operacional  
**Creado:** 2026-07-05
