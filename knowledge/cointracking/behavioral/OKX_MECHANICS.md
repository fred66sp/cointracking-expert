---
id: KB-B2-012
title: "Mecánicas de OKX: Trading Completo + Web3"
level: B
domain: cointracking
source: "OKX official docs + análisis de casos reales"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-05
confidence: high
version: 1.0

related_adr:
  - ADR-003
  - ADR-010
  - ADR-022

related_docs:
  - knowledge/cointracking/behavioral/BINANCE_SPOT_MECHANICS.md
  - knowledge/cointracking/behavioral/BYBIT_MECHANICS.md
  - knowledge/wallets/METAMASK_INTEGRATION.md

tags:
  - exchange
  - okx
  - behavioral
  - spot
  - derivados
  - web3

notes: "OKX es exchange asiático complejo con ecosistema Web3 integrado (wallet, DEX, staking). Mayor que Bybit en volumen. Integración CoinTracking disponible pero requiere atención a API scope."
---

# OKX Mechanics

## Características Principales

**OKX** (ex-OKEx) es el tercer exchange más grande del mundo, con fuerte presencia en Asia y crecimiento en occidente. Ofrece un ecosistema híbrido:

- **Trading Centralizado:** Spot, Margin, Perpetuals (muy completo)
- **Wallet Web3 Nativa:** Almacenamiento de cripto on-chain (similar a MetaMask pero integrado)
- **DEX Integrado:** Swap en-cadena (Uniswap/Curve) accesible desde la app
- **Staking & Earn:** Productos de rendimiento (delegado, yield farming)
- **Mining & NFTs:** Soporte para minería y operaciones NFT

---

## Operaciones en CoinTracking

### 1. Spot Trading

| Campo | Valor | Notas |
|-------|-------|-------|
| **Type** | Trade | Compra/venta estándar |
| **Exchange** | OKX | Identificador único |
| **Currency In/Out** | USDT, BTC, ETH, USDC, etc | Pares muy amplios |
| **Fee Currency** | OKB (OKX Token) o USDT | Reducción si usas OKB |
| **Trade ID** | Presente | Único por operación |

**Ejemplo real:**
```
Date: 2024-06-15 14:23:00
Type: Trade (BUY)
Buy: 2.0 ETH @ 2500 USDT c/u
Fee: 5 USDT (0.1%, reducible con OKB)
Exchange: OKX
```

### 2. Margin Trading (Cross & Isolated)

**OKX ofrece dos modos:**
- **Cross Margin:** Un solo colateral para múltiples posiciones
- **Isolated Margin:** Colateral dedicado por par (más seguro)

| Campo | Diferencia |
|-------|-----------|
| **Type** | Trade (igual) |
| **Grupo** | "OKX Cross" o "OKX Isolated" (recomendado) |
| **Fee** | Interés de préstamo (variable, 5-15% anual) |
| **Risk** | Liquidación automática si ratio cae |

**⚠️ Crítico:** Cross margin es más arriesgado (liquidación domino). Documentar el modo en el comentario.

### 3. Perpetual Futures

**OKX ofrece múltiples tipos:**

| Tipo | Notación | Tamaño Contrato | Apalancamiento |
|------|----------|---|---|
| **Standard Perpetual** | Linear (USDT) | 1 Cripto | Hasta 75x |
| **Inverse Perpetual** | Inverse (Cripto) | 100 Cripto | Hasta 125x |
| **Quarterly Futures** | Vence cada 3 meses | 1 Cripto | Hasta 75x |

**En CoinTracking:**
```
Type: "Perpetual Futures" o "Derivatives"
Funding Fee: Cada 8 horas (como otros exchanges)
Position Type: LONG / SHORT
Apalancamiento: Documentar en comentario (crítico para auditoría)
```

### 4. Wallet Web3 & DEX Integrado

**Operación especial:** Swap en-cadena (Uniswap V3, Curve, etc) accesible desde OKX app.

```
Type: Swap
From: 1.0 USDC
To: 3000 USDT
Gas: 0.001 ETH (Ethereum) o variable (otra cadena)
Protocol: Uniswap V3 (ejemplo)
Chain: Ethereum, Polygon, Arbitrum, etc
```

**En CoinTracking:**
- Aparece como DOS operaciones: Sell USDC, Buy USDT
- Gas incluido en cost basis de USDC

**⚠️ Especial:** OKX wallet es **custodia del usuario, no de OKX**. Los swaps on-chain generan dirección de transacción verificable en blockchain explorer (importancia para auditoría).

### 5. Staking y Earn

**OKX ofrece múltiples productos:**

| Producto | Tipo | APY | Bloqueo |
|----------|------|-----|---------|
| **Staking delegado** | Proof of Stake | 3-8% | Variable |
| **Earn (flexible)** | Lending | 2-5% | Ninguno |
| **Earn (fixed)** | Lending 30-180 días | 5-12% | Sí |
| **Yield farming** | DEX LP | 5-50% | Sí (pool lock) |

**En CoinTracking:**
```
Operación 1 - Deposit: 10 ETH → Staking Pool
Operación 2 - Reward: 0.3 ETH (periódico)
Operación 3 - Withdraw: 10.3 ETH
```

**Fiscalidad:** Rewards = ingresos del capital (base general, Modelo 721), no base del ahorro.

---

## Integración con CoinTracking

### ✅ Métodos Soportados

1. **API Connection (Recomendado)**
   - CoinTracking soporta OKX API (validado 2026-07-05)
   - ⚠️ **API Scope crítico:** Elegir correctamente qué subsistemas importar (Spot, Margin, Futures, Wallet)
   - Datos en vivo, actualizados automáticamente
   - Comisiones incluidas

2. **CSV Import**
   - OKX permite exportar por subsistema (Spot, Margin, Futures, Earn)
   - Formato: variable según producto
   - Limitación: Manual, requiere actualización periódica
   - **Wallet Web3 no exporta CSV** (solo blockchain address)

### Importación vía API

**Pasos (crítico: API Scope):**

1. OKX: Account → API Management → Create Trading Account API
2. **Permisos: seleccionar solo "Read-Only" (lectura)**
3. **Scope: Eligir qué incluir:**
   - ✅ Spot Trading
   - ✅ Margin Trading
   - ✅ Perpetual Futures
   - ⚠️ Wallet (si usas on-chain swaps; ver más abajo)
   - ⚠️ Staking/Earn (si aplica)
4. CoinTracking: Settings → Exchanges → OKX → Connect
5. Pegar API Key, Secret, Passphrase (OKX requiere 3 campos)
6. Sincronizar

**Validación:**
- [ ] Balance en CoinTracking coincide con OKX
- [ ] Operaciones Spot importadas correctamente
- [ ] Operaciones Margin visible (si aplica)
- [ ] Operaciones Futures visible (si aplica)
- [ ] Comisiones incluidas (OKB o USDT)
- [ ] Trade IDs únicos
- [ ] Rewards/Staking visible (si aplica)

### Wallet Web3 OKX (Caso Especial)

Si usas **OKX Wallet integrado** para on-chain swaps:
- OKX Wallet **NO se importa por API** (es custodia del usuario, no de OKX)
- **Solución:** Importar dirección pública del wallet en CoinTracking como "Blockchain Address" (igual que Metamask o Ledger)
- Transacciones on-chain aparecerán automáticamente

**En CoinTracking:**
```
Add Account → Blockchain → Ethereum (u otra cadena)
Pegar dirección del OKX Wallet
CoinTracking escanea y captura transacciones on-chain
```

---

## Casos Límite y Peculiaridades

### 1. Comisión en OKB (OKX Token)

**Mecanismo:** OKX incentiva usar OKB para reducir comisión (0.1% → 0.075% o mejor).

**Tratamiento:**
- CoinTracking detecta automáticamente comisiones en OKB
- Fee se registra en la moneda especificada (OKB)
- **Verificar:** Saldo de OKB no explicado = restos de comisiones

**Fiscalidad:** Comisión sigue siendo gasto, sea en USDT u OKB.

### 2. Cross Margin Liquidación (Domino)

**Riesgo especial:** En cross margin, liquidación de un par puede caer en cascada.

```
Margen disponible: 1000 USDT
Posición 1: BUY 10 ETH @ 2000 USDT (colateral 20000)
Posición 2: BUY 2 BTC @ 65000 USDT (colateral 130000)
Precio cae: ETH → 1000, BTC → 30000
Liquidación: Ambas posiciones cerradas
Pérdida: Casi total del colateral
```

**En CoinTracking:**
- Aparecerá como Liquidation
- Documentar el modo ("Cross") en comentario

**Fiscalidad:** Pérdida patrimonial total (deducible).

### 3. Quarterly Futures (Vencimiento)

Futures que vencen cada 3 meses (no perpetuos).

```
Contrato: ETHUSDT-230616 (vence 16 Jun 2023)
Posición: LONG 5 ETH
Cierre: Automático al vencimiento
PnL: Liquidado y transferido a cuenta spot
```

**En CoinTracking:**
- Tratamiento igual a perpetuos, pero con fecha de cierre conocida
- Documentar fecha de vencimiento en comentario

**Fiscalidad:** Ganancia/pérdida patrimonial al vencimiento.

### 4. Wallet Web3 Swaps On-Chain

**OKX Wallet puede hacer swaps** sin tocar el exchange:

```
OKX Wallet Address: 0xabc123...
Operación: Swap 1.0 USDC → 3000 USDT en Uniswap V3
Chain: Ethereum
Gas: 0.002 ETH
```

**En CoinTracking:**
- Aparecerá como DOS operaciones importadas vía blockchain address
- Sell: 1.0 USDC
- Buy: 3000 USDT
- Fee: Gas (0.002 ETH)

**Fiscalidad:** Ganancia patrimonial (igual que manual DEX swap).

### 5. Staking Delegado vs Yield Farming

**Diferencia fiscal:**

| Tipo | Producto OKX | Fiscalidad |
|------|---|---|
| **Staking delegado** | "Staking" en OKX | RCM, base del ahorro |
| **Lending** | "Earn Flexible" o "Earn Fixed" | RCM, base del ahorro |
| **Yield Farming** | "DeFi Strategies" | RCM, base del ahorro (pero más riesgo) |

**Regla común:** Todos son ingresos del capital, pero la base inicial (depósito) es el coste de adquisición para futuro FIFO.

### 6. Fee Tiers y Maker/Taker

OKX ofrece comisiones **dinámicas** según volumen y rango:

```
Volumen 30d < $100k: Taker 0.15%, Maker 0.10%
Volumen 30d $100k-$500k: Taker 0.12%, Maker 0.08%
...
Con OKB: -25% en cada tier
```

**En CoinTracking:**
- La comisión importada debe ser la **real pagada**, no la teórica
- Si hay duda, comparar con invoice de OKX

---

## Validación en CoinTracking

### Checklist: OKX Completo

```
[ ] API conectado en CoinTracking
[ ] Balance actual coincide (refresh)
[ ] Spot: últimas 10 operaciones visibles
[ ] Margin: últimas 10 operaciones visibles (si aplica)
[ ] Futures (Perpetual + Quarterly): visible (si aplica)
[ ] Comisiones incluidas (USDT o OKB)
[ ] Trade IDs únicos (sin duplicados)
[ ] Funding fees capturados (si hay perpetuos)
[ ] Margin interest visible (si hay margin)
[ ] Staking rewards visible (si aplica)
[ ] Wallet Web3 importado como Blockchain Address (si aplica)
[ ] Saldos negativos: 0 (excepto fiat si es normal)
[ ] No hay solapamiento API+CSV+Blockchain
```

### Detección de Problemas Comunes

| Problema | Síntoma | Solución |
|----------|---------|----------|
| **API desconectado** | Balance no actualiza | Reconectar API |
| **Scope insuficiente** | Spot visible, Futures no | Volver a crear API key con Futures habilitado |
| **CSV duplica** | Operaciones 2x | Eliminar CSV, usar solo API |
| **OKB fee no detectado** | Comisión falta | Verificar que CoinTracking soporta OKB (debería) |
| **Cross margin liquidación oculta** | Pérdida sin operación | Buscar "Liquidation" |
| **Wallet Web3 no sincroniza** | Swaps on-chain no aparecen | Agregar dirección wallet como Blockchain Address |
| **Quarterly futures vencen sin registro** | PnL desaparece | Verificar fecha de vencimiento + cierre automático |

---

## Comparativa: OKX vs Binance vs Bybit

| Aspecto | Binance | Bybit | OKX |
|---------|---------|-------|-----|
| **Spot comisión** | 0.1% (BNB) | 0.1% (BIT) | 0.1% (OKB) |
| **Margin tipos** | Cross, Isolated | Cross, Isolated | Cross, Isolated |
| **Futures** | Perpetual + Quarterly | Perpetual | Perpetual + Quarterly |
| **Max apalancamiento** | 75x | 75x | 75x |
| **Wallet integrada** | No | No | Sí (Web3) |
| **DEX integrado** | No | No | Sí |
| **Staking** | Sí (Earn) | No | Sí (Staking + Earn) |
| **Volumen** | Más alto | Alto | Alto |
| **API estabilidad** | Muy alta | Alta | Alta |
| **CoinTracking support** | Completo | Completo | Completo |
| **Complejidad auditoría** | Media | Media | Alta (Web3) |

---

## Casos de Uso Reales

### Caso 1: Spot + Perpetuos + Wallet Web3

**Usuario:** Opera en Spot OKX (buy/hold), abre perpetuos como hedge, y usa OKX Wallet para yield farming en Uniswap.

```
Operaciones Spot: BUY 2 ETH @ 2500 USDT
Operaciones Futures: SHORT 2 ETH perpetual @ 2600 USDT (hedge)
Wallet Web3: Deposit 1 ETH en Uniswap V3 LP (comisión 0.3%)
```

**En CoinTracking:**
- Spot: importado vía API
- Futures: importado vía API
- Wallet: importado vía blockchain address

**Fiscalidad:** Spot (FIFO), Futures (derivado), LP (ganancias patrimoniales + rewards = RCM).

### Caso 2: Fixed Earn + Vencimiento

**Usuario:** Pone 10000 USDT en OKX Earn Fixed (30 días, 8% APY).

```
Deposit: 10000 USDT (2024-06-01)
Interest accrual: 65 USDT (2024-07-01)
Withdrawal: 10065 USDT
```

**En CoinTracking:**
- Deposit: registrar como Staking/Income base
- Interest: registrar como Income (RCM, base general)
- Withdrawal: registrar como retiro

**Fiscalidad:** 65 USDT = ingresos del capital.

### Caso 3: Liquidación Quarterly Futures

**Usuario:** Abre quarterly futures, pero el precio cae.

```
Contrato: BTCUSDT-230916 (vence 16 Sep 2023)
Posición: SHORT 0.5 BTC @ 30000 USDT
Cierre automático: 29000 USDT (30 días antes del vencimiento)
PnL: +500 USDT ganancia
```

**En CoinTracking:**
- Registrar como Derivatives con fecha de cierre
- PnL: +500 USDT (ganancia patrimonial)

**Fiscalidad:** Ganancia patrimonial (500 USDT).

---

## Referencias y Recursos

- [OKX Official API Docs](https://www.okx.com/docs/en/) (multilingual)
- [CoinTracking OKX Integration](https://www.cointracking.info/en/api_keys.php) (en plataforma)
- [OKX Help Center](https://www.okx.com/help/) (español disponible)
- [OKX Wallet (Web3)](https://www.okx.com/web3/wallet)

---

## Notas Operativas

**Para auditoría:**
- OKX es muy complejo (3 subsistemas trading + Web3)
- API scope debe elegirse con cuidado (no importar sin necesidad)
- Wallet Web3 requiere importación manual por blockchain address
- Quarterly futures tienen fecha de vencimiento fija (importante para cierre de período)

**Para fiscalidad:**
- Spot: ganancias patrimoniales (FIFO)
- Margin: igual que Spot (+ riesgo de liquidación)
- Perpetuos: derivado (consultar asesor)
- Quarterly futures: derivado con vencimiento conocido
- Staking/Earn/Farming: ingresos del capital (RCM, base general)
- Wallet Web3 swaps: ganancias patrimoniales (igual que DEX manual)

---

**Documento:** OKX Mechanics  
**Nivel:** B2-012  
**Status:** Operacional  
**Creado:** 2026-07-05
