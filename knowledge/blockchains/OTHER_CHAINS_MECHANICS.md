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
id: KB-B3-003
title: "Cómo CoinTracking maneja otras blockchains (Polygon, BSC, Solana, etc.)"
level: B
domain: blockchains
source: "Análisis de casos reales"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: medium
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md
  - knowledge/blockchains/BITCOIN_TRANSACTION_TYPES.md

tags:
  - blockchain
  - polygon
  - bsc
  - solana
  - multichain
  - behavioral

notes: "Operativo: cómo CoinTracking trata blockchains L2 y alternativas."
---

# Cómo CoinTracking Maneja Otras Blockchains

## Panorama General

**Blockchains principales** (aparte de Ethereum y Bitcoin):

| Chain | Modelo | Fee Gas | Problema en CT | Solución |
|-------|--------|---------|---|---|
| **Polygon** | EVM (Ethereum clone) | MATIC (muy bajo) | Transfers duplicados | Reimportar via Polygon RPC |
| **BSC** (BNB Smart Chain) | EVM | BNB (bajo) | Confunde con Binance | Etiquetar "Binance Chain" |
| **Solana** | Unique (no EVM) | SOL (variable) | No detecta SPL tokens bien | Import manual |
| **Arbitrum** | EVM | ARB + ETH | Nuevo, soporte limitado | Usar Etherscan + manual |
| **Optimism** | EVM | OP + ETH | Igual que Arbitrum | Igual |

---

## Tipo 1: Polygon (Layer 2 de Ethereum)

```
Característica clave:
  - Blockchain EVM (compatible con Ethereum)
  - Fees muy bajos (MATIC cost ~0.001€)
  - Tokens "wrapped" (USDT.e, USDC.e en lugar de USDT, USDC)
  
¿Cómo aparece en CoinTracking?
  
  Ideal:
    - API de CoinTracking detecta "Polygon network"
    - Clasifica transacciones como "Polygon USDT.e"
  
  Problema:
    - Si importas CSV de Polygon: CoinTracking ve "USDT.e" como activo diferente
    - Puede mezclar con "USDT" de Ethereum
    - Balance confuso
```

---

## Validación en Polygon

```
CoinTracking → Transacciones:
  ¿Aparece "Polygon" o "MATIC" en la descripción?
    SÍ → OK
    NO → Verificar que no sea duplicado de Ethereum
    
Chequeo:
  - Polygon USDC ≠ Ethereum USDC
  - Cada uno es un activo separado
  - Si tienes ambos, deben aparecer como dos saldos diferentes
```

---

## Tipo 2: BSC (BNB Smart Chain)

**PELIGRO CRÍTICO: Binance Chain vs Binance Exchange**

```
¿Confusión común?
  - Transacción en BSC (blockchain)
  - Binance Exchange (plataforma de trading)
  
Son COMPLETAMENTE DIFERENTES:
  1. BSC = blockchain (como Ethereum)
  2. Binance Exchange = plataforma (como Kraken)
  
¿Cómo aparece en CoinTracking?
  
  Ideal:
    - API detecta "BSC network"
    - Clasifica como "BNB Smart Chain"
  
  Problema:
    - Si importas CSV de BSC: CoinTracking puede confundirlo con Spot Binance
    - Resultado: duplicados o confusión
```

---

## Validación en BSC

```
Si tienes transacciones en BSC:
  
  1. Verificar que sean BNB Smart Chain, no Binance Spot
     → Ir a Bscscan.com (explorador BSC)
     → Buscar tu transaction hash
     → Si aparece en bscscan (no en binance.com) → Es BSC
  
  2. En CoinTracking:
     → Etiquetar como "BNB Smart Chain" o "Binance Chain"
     → NO confundir con Binance Spot trading
```

---

## Tipo 3: Solana

**CRÍTICO: Solana es TOTALMENTE DIFERENTE a EVM**

```
Característica clave:
  - No es EVM (no es compatible con Ethereum)
  - Tokens SPL (Solana Program Library)
  - Fees en SOL (micro fees, muy bajos)
  - Transacciones más complejas (cuenta sistema)
  
¿Cómo aparece en CoinTracking?
  
  Problema:
    - CoinTracking tiene soporte LIMITADO para Solana
    - SPL tokens a menudo no se importan correctamente
    - Balance puede ser incompleto
  
  Solución:
    - Importar manualmente desde Solscan.io (explorador)
    - Crear operaciones manualmente si falta algo
```

---

## Validación en Solana

```
Si tienes fondos en Solana:
  
  1. Verificar saldo en Solscan.io
     → Buscar tu wallet (dirección pública)
     → Anotar todos los SPL tokens
  
  2. Comparar contra CoinTracking
     → ¿Aparecen todos los tokens?
       SÍ → OK
       NO → Hay que importar manualmente
  
  3. Añadir manualmente en CoinTracking
     → Si falta algún token SPL
     → Crear como "Deposit" en la fecha que recibiste
```

---

## Tipo 4: Arbitrum y Optimism (Layer 2 de Ethereum)

```
Característica clave:
  - Son "Ethereum L2" (más baratos, más rápidos)
  - Transacciones eventualmente se liquidan en Ethereum mainnet
  
¿Cómo aparece en CoinTracking?
  
  Estado (2026):
    - Soporte limitado en CoinTracking
    - Muchas transacciones se ven como "Ethereum" incorrectamente
    - Confusión con gas fees
  
  Solución (hoy):
    - Verificar en Arbiscan.io (Arbitrum) o Optimismscan.io
    - Si CoinTracking no las detecta bien, importar manualmente
    - Etiquetar claramente "Arbitrum" o "Optimism"
```

---

## Consolidación: Multichain Balance

**COMPLICACIÓN: Un mismo token en varias chains**

```
Ejemplo:
  - USDC en Ethereum (dirección A)
  - USDC en Polygon (dirección B)
  - USDC en Arbitrum (dirección C)
  - USDC en BSC (dirección D)
  
¿Son el mismo USDC?
  - NO. Cada uno es un token separado en su blockchain
  - No son intercambiables (sin bridge)
  
¿Cómo CoinTracking lo ve?
  
  Ideal:
    - 4 activos diferentes: USDC, USDC.e, USDC (Arbitrum), USDC (BSC)
    - Saldos separados
  
  Problema:
    - CoinTracking los agrupa como "USDC" sin distinguir chain
    - Balance total parece correcto, pero origen es confuso
  
  Solución:
    - Buscar en CoinTracking si hay "etiqueta de chain"
    - Si no existe, etiquetar manualmente
    - Documentar: "USDC Ethereum 1000€" + "USDC Polygon 500€" ≠ "USDC Total 1500€" si están en chains diferentes
```

---

## Tratamiento Fiscal

**IMPORTANTE: Todos los blockchains son iguales fiscalmente**

```
No importa si es Ethereum, Polygon, Solana, etc.

Regla fiscal española:
  - Venta de cripto = ganancia patrimonial
  - Es lo mismo en cualquier blockchain
  - El token USDC en Polygon o Ethereum vale lo mismo en €
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Extendible a cualquier blockchain
- **ETHEREUM_TRANSACTION_TYPES.md:** Base de conceptos EVM
- **BITCOIN_TRANSACTION_TYPES.md:** Base de conceptos UTXO
