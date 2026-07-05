---
id: KB-B1-011
title: "Stablecoins y Tokens Especiales en CoinTracking"
level: B
domain: cointracking
source: "CoinTracking + análisis blockchain + casos reales"
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
  - knowledge/taxation/spain/CAPITAL_GAINS.md
  - knowledge/cointracking/behavioral/DEFI_SWAPS_MECHANICS.md
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md

tags:
  - stablecoins
  - altcoins
  - special-tokens
  - bridges
  - wrapped-tokens

notes: "Stablecoins son activos con precio fijo (idealmente 1 USD = 1 token). Pueden estar en múltiples redes. Wrapped tokens son derivados de activos base, requieren cuidado fiscal."
---

# Stablecoins y Tokens Especiales

## ¿Qué son los Stablecoins?

**Stablecoins** son criptomonedas cuyo precio está **vinculado a un activo de referencia** (normalmente USD, EUR, o una cesta).

### Tipos principales

| Tipo | Ejemplo | Respaldado | Riesgo |
|------|---------|-----------|--------|
| **Collateralized (Fiat)** | USDC, USDT | 1:1 con USD en banco | Bajo (si auditoría es verificable) |
| **Collateralized (Crypto)** | DAI (MakerDAO) | Exceso de cripto bloqueado | Medio (dependencia de Oracle) |
| **Algorithmic** | UST (Luna) | Algoritmo + incentivos | Alto (fallidos históricamente) |
| **Commodity-backed** | PAXG (oro digital) | 1 onza de oro | Bajo (si auditoría verificable) |

---

## Stablecoins Más Comunes en CoinTracking

### USDT (Tether)

**Características:**
- Mayor volumen de todos los stablecoins
- Emitido en múltiples blockchains (Ethereum, Tron, BSC, Polygon, etc.)
- Polémico: auditorías incompletas, cuestionamiento sobre respaldo 100%

**En CoinTracking:**
```
Ticker: USDT
Símbolo: Aparece igual en todas las redes, pero TX Hash varía
Precio: Debe ser ~1.00 USD siempre
Riesgo: Si importas desde múltiples redes sin verificar, puede parecer "ganancia" 
        (USDT Ethereum → USDT Tron, aunque sean el mismo valor)
```

**Caso límite:** Puente entre Tron y Ethereum
```
Retirada: 1000 USDT (Ethereum) — costo: 0.01 ETH gas
Depósito: 1000 USDT (Tron) — depósito sin comisión
En CoinTracking: Sale 1000, entra 1000 → sin ganancia/pérdida
Fiscalidad: Transferencia (no tributa)
```

### USDC (USD Coin)

**Características:**
- Más nuevo, mejor auditoría que USDT
- Respaldado por Circle + Coinbase (institucionales confían)
- Menos volumen que USDT, pero creciendo
- Emitido en Ethereum, Polygon, Solana, Arbitrum, Optimism, etc.

**En CoinTracking:**
```
Ticker: USDC
Precio: Debe ser ~1.00 USD
Riesgo: Mismo que USDT (múltiples redes, confusión de símbolo)
```

**Caso real — USDC.e vs USDC:**
```
USDC (nativo en Ethereum): ticker "USDC"
USDC.e (wrapped en Arbitrum): ticker "USDC.e" (distinto)
En CoinTracking: APARECEN COMO ACTIVOS DISTINTOS (correcto)
```

### DAI (MakerDAO)

**Características:**
- Stablecoin descentralizado (no tiene autoridad central)
- Respaldado por cripto (ETH, USDC, etc.)
- Puede perder peg (1 DAI ≠ 1 USD si hay crisis)
- Genera interés en ciertos protocolos

**En CoinTracking:**
```
Ticker: DAI
Precio: Debe ser ~1.00 USD, pero puede fluctuar ±5% en crisis
Riesgo: Posible pérdida si el peg se rompe (raro, pero sucedió 2020-2023)
```

**Caso especial — Earning DAI:**
```
Depósito: 1000 DAI en Aave Lending
Reward: 50 DAI interest (después de 6 meses, APY 5%)
En CoinTracking: 50 DAI = ingresos del capital
Fiscalidad: Sí tributa (RCM, base general)
```

### Otros Stablecoins (Menos Comunes)

| Stablecoin | Emisor | Riesgos |
|---|---|---|
| **PAXUSD** | Paxos | Auditado, buena reputación |
| **BUSD** | Binance + Paxos | Binance dejó de emitir (2023) |
| **TUSD** | TrueUSD | Menor volumen |
| **EUR (sEUR, euroc, etc.)** | Varios | Menos común, solo usuarios europeos |

---

## Wrapped Tokens (Derivados de Base)

**Wrapped token** = versión de un activo en una blockchain distinta a su original.

### Ejemplos Clave

**Wrapped Bitcoin (wBTC)**
```
Original: BTC (Bitcoin blockchain)
Wrapped: wBTC (Ethereum ERC-20 token)
Relación: 1 wBTC = 1 BTC (teóricamente)
Mecanismo: Custodia 1:1 (bridge)
```

**En CoinTracking:**
- `BTC` y `wBTC` son **activos distintos**
- Cambiar BTC → wBTC es un **Trade** (ganancia/pérdida posible si hay fee)
- Vender wBTC = venta de wBTC, no de BTC

**Fiscalidad:**
- Compra 1 BTC @ 60,000 EUR (cost basis: 60,000)
- Wrap a wBTC (0 coste, solo transferencia)
- Vende 1 wBTC @ 59,900 EUR (pérdida: -100 EUR patrimonial)

**Caso límite — Descuento de wrap:**
```
Compra: 1 BTC @ 60,000 EUR
Wrap a wBTC: Fee 100 EUR (custodia)
Vende wBTC @ 60,000 EUR
PnL: -100 EUR (fee) [ganancia 0, pero gastó]
Fiscalidad: ¿Es el fee deducible? [VERIFICAR con asesor]
```

### Wrapped Ether (wETH)

```
Original: ETH (Ethereum blockchain nativo)
Wrapped: wETH (token ERC-20 en Ethereum)
Relación: 1 wETH = 1 ETH
Razón común: Necesitar ETH como token para contratos DEX
```

**En CoinTracking:**
- Aparecen como activos distintos (ETH vs wETH)
- Conversión ETH → wETH es **sin coste** (solo gas)
- Pero a veces hay slippage si usas DEX para convertir

### Otros Wrapped Tokens

| Original | Wrapped | Chain | Uso |
|----------|---------|-------|-----|
| BTC | wBTC, renBTC, sBTC | Ethereum, Polygon | DeFi en Ethereum |
| ETH | wETH | Casi toda cadena | Compatibilidad ERC-20 |
| SOL | wSOL | Ethereum, Polygon | Cross-chain DeFi |
| AVAX | wAVAX | Ethereum, Polygon | Liquidez cruzada |

---

## Token Splits y Token Migrations

### Token Split (Ejemplo: SHIB)

Algunos proyectos **dividen el token en 1 millón de fragmentos** sin cambiar el valor.

```
Antes: 1 SHIB = 0.00001 USD (1 millón de SHIB = 10 USD)
Split 1:1000000: 1.000.000 SHIB = 0.00001 USD (1 millón de SHIB = 10 USD)
```

**Efecto en CoinTracking:**
- La cantidad cambia (×1 millón)
- El precio unitario baja (÷1 millón)
- El valor total es el mismo
- **Riesgo:** Si no importas el split, ves "pérdida masiva" en balance

**Cómo verificar:** Buscar en CoinGecko si hubo "Split" en la historia del token.

### Token Migration

Algunos proyectos **cambian de blockchain** o **reemiten un nuevo token**.

```
Ejemplo: Luna Classic (LUNC) vs Luna 2.0 (LUNA)
Antes (2022): 1 LUNA = 80 USD (alto, antes del crash)
Crash (2022): 1 LUNA = 0.00001 USD (colapso)
Reemisión (2023): Nuevo LUNA lanzado, holders de LUNA antigua reciben 1:1 conversion
```

**En CoinTracking:**
- Viejo token: LUNA (anterior al crash)
- Nuevo token: LUNA (reemitido)
- Pueden aparecer **como el mismo ticker** (peligro de confusión)
- Solución: CoinTracking debería detectar migración automáticamente

**Caso límite — Sin migración en CoinTracking:**
```
Tenías 10,000 LUNA @ $80 = $800,000 (2022)
Crash: LUNA → $0.00001 = $0.10 (pérdida total)
Reemisión: Recibes 10,000 LUNA nuevo @ $0.01 = $100
En CoinTracking: Si no registra migration, ves balance "malo" sin explicación
```

---

## Bridges y Multi-Chain Stablecoins

**Bridge** = mecanismo para transferir un token de una blockchain a otra.

### Ejemplo: USDC en múltiples cadenas

```
USDC Ethereum: Address 0xA0b86991c6218b36...
USDC Polygon: Address 0x2791Bca1f2de4661...
USDC Arbitrum: Address 0xFF970A61A04b1cA5...
```

Cada uno es un **contrato distinto** pero representan el **mismo activo** (USDC).

**En CoinTracking:**
- All importados como "USDC"
- Precio = ~1.00 USD (en todas)
- Transferencia entre redes = "transferencia" (no tributa)

**Caso límite — Bridge Fee:**
```
Envío: 1000 USDC Ethereum → Polygon via Stargate Bridge
Fee: 1 USDC (bridge + slippage)
Recibido: 999 USDC Polygon
En CoinTracking: Sale 1000, entra 999 → pérdida de 1 USDC (registrada)
Fiscalidad: Fee es costo de operación (probablemente deducible)
```

---

## Importación de Stablecoins en CoinTracking

### Validación de Precio

**Checklist:**
```
[ ] USDT: Precio ~1.00 USD (alerta si ±5%)
[ ] USDC: Precio ~1.00 USD (alerta si ±5%)
[ ] DAI: Precio ~1.00 USD (alerta si ±10%, crisis)
[ ] wBTC: Precio ≈ BTC precio (alerta si ±5%, fee de custodia)
[ ] wETH: Precio ≈ ETH precio (alerta si ±5%)
```

### Validación de Red

```
Stablecoin XYZ en múltiples redes:
[ ] USDC Ethereum (original)
[ ] USDC Polygon (bridged)
[ ] USDC Solana (bridged)
[ ] Verificar que CoinTracking los importa en la "cuenta" correcta
[ ] Verificar que los saldos sumados coinciden con la suma real
```

### Validación de Símbolo

**Problema:** Múltiples tokens pueden tener mismo símbolo si vienen de exchanges distintos.

```
Ejemplo: "USDT" puede ser:
- USDT Ethereum (ERC-20)
- USDT Tron (TRC-20)
- USDT BSC (BEP-20)
CoinTracking debería diferenciarlos, pero verificar.
```

---

## Fiscalidad de Stablecoins

### Stablecoins como "Cash Proxy"

**Concepto:** Si compras USDT @ 1.00 USD y vendes @ 1.00 USD, ¿hay ganancia?

**Respuesta en España:**
- Técnicamente SÍ (es un activo, no dinero fiat)
- Pero en práctica: **cero ganancia** (si la venta es @ 1.00)
- Implicación: No requiere declaración si no hay diferencia de precio

**Regla fiscal:**
```
Ganancia = Venta - Costo - Comisión
USDT: Costo 1000 USD, Venta 1000 USD = Ganancia 0 → no tributa
```

### Wrapped Tokens

**Regla:** El wrapped token es un **activo distinto** del base.

```
Compra 1 BTC @ 60,000 EUR
Wrap a wBTC (sin fee): Cost basis sigue siendo 60,000 EUR en wBTC
Vende 1 wBTC @ 59,900 EUR
PnL: -100 EUR (pérdida patrimonial, deducible)
```

### Stablecoins en Yield Farming

**Caso:** Depositas USDC en Aave, recibes interest.

```
Deposit: 1000 USDC (base de coste)
Interest: 50 USDC (después de 6 meses)
Fiscalidad: 50 USDC = ingresos del capital (RCM, base general)
Valor: 50 USD × 1.00 = 50 EUR
```

---

## Mejores Prácticas en CoinTracking

### ✅ Hacer

1. **Verificar precio de stablecoins** (debe ser ~1.00)
2. **Documentar red** en el comentario (USDC Ethereum vs Polygon)
3. **Registrar bridges** como transferencias (no trades)
4. **Rastrear wrapped tokens** separado del base
5. **Incluir fee de bridge** en cost basis
6. **Diferenciar nuevo token** del viejo en caso de migration

### ❌ No hacer

1. ❌ **Asumir precio 1.00** sin verificar (algunos perden peg)
2. ❌ **Confundir BTC con wBTC** (son activos distintos)
3. ❌ **Registrar bridge como Trade** (es transferencia)
4. ❌ **Perder track de wrapped tokens** (aparecen como "pérdida" si no se ven)
5. ❌ **Ignorar token splits** (aparecen como "pérdida masiva")

---

## Referencias

- [Circle USDC (oficial)](https://www.circle.com/usdc)
- [Tether USDT (oficial)](https://tether.to/)
- [MakerDAO DAI (oficial)](https://makerdao.com/)
- [Wrapped Bitcoin (oficial)](https://wbtc.network/)
- [Stargate Finance (bridges)](https://stargate.finance/)
- [CoinGecko (historial de tokens)](https://www.coingecko.com/)

---

**Documento:** Stablecoins and Special Tokens  
**Nivel:** B1-011  
**Status:** Operacional  
**Creado:** 2026-07-05
