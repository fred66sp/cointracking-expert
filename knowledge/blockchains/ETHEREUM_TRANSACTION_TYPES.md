---
id: KB-B3-001
title: "Tipos de transacciones Ethereum y cómo afectan a CoinTracking"
level: B
domain: blockchains
source: "Ethereum Yellow Paper + casos reales"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/cointracking/behavioral/DEFI_SWAPS_MECHANICS.md
  - knowledge/cointracking/behavioral/LENDING_MECHANICS.md

tags:
  - ethereum
  - blockchain
  - transactions
  - behavioral

notes: "Operativo: tipos de tx en Ethereum y cómo aparecen en CoinTracking."
---

# Tipos de Transacciones Ethereum y su Impacto en CoinTracking

## Definición: Tipos de Transacción

En Ethereum, una "transacción" (tx) es un movimiento de datos/valor. CoinTracking lo ve como un evento de blockchain.

**Tipos principales:**
1. **ETH transfers** — Movimiento simple de ETH
2. **ERC-20 transfers** — Movimiento de tokens (USDC, DAI, etc.)
3. **Contract interactions** — Llamadas a smart contracts (swaps, lending, etc.)
4. **Internal transfers** — Movimientos dentro de un contrato

---

## Tipo 1: ETH Transfers (Movimiento Simple)

```
En blockchain:
  De: 0xAlice...
  A: 0xBob...
  Cantidad: 10 ETH
  Fee (gas): 0.001 ETH
  
CoinTracking ve:
  - ETH Sale: 10.001 ETH (incluye gas)
  O
  - ETH Withdrawal: 10 ETH
  - Gas expense: 0.001 ETH (si está disponible)
  
Fácil de auditar: ✓ Único movimiento
```

---

## Tipo 2: ERC-20 Transfers (Token Movements)

```
En blockchain:
  Token: USDC
  De: 0xAlice...
  A: 0xBob...
  Cantidad: 1.000 USDC
  Fee (gas en ETH): 0.005 ETH
  
CoinTracking ve:
  - USDC Withdrawal: 1.000 USDC
  - Gas expense: 0.005 ETH
  
Fácil de auditar: ✓ Dos líneas (token + gas)
```

---

## Tipo 3: Contract Interactions (Swaps, Lending, etc.)

**CRÍTICA: Este es el más complejo**

```
En blockchain (Uniswap swap, ejemplo):
  1. Llamada a Uniswap contract
  2. "Envío" 1.000 USDC al contrato
  3. Contrato "Devuelve" 1 ETH a mi wallet
  4. Gas: 0.02 ETH
  
CoinTracking ve:
  - USDC sent: 1.000 USDC
  - ETH received: 1 ETH
  - Gas: 0.02 ETH
  
PERO: ¿Los relaciona como un swap o como 2 operaciones?
  → SI hay API de Uniswap: Trade (1 operación)
  → SI NO hay API: Dos transfers separadas
  
Fácil de auditar: ✗ Necesita contexto
```

---

## Tipo 4: Internal Transfers (Movimientos dentro de Contratos)

```
En blockchain:
  - Un contrato mueve fondos a otro contrato
  - No hay transferencia de gas (lo paga quien llama)
  - CoinTracking frecuentemente NO los ve
  
Ejemplo:
  1. Yo llamo Aave deposit
  2. Aave mueve USDC → reserve interna
  3. Me envía aToken (token de recibo)
  
CoinTracking ve:
  - Transfer out: USDC
  - Transfer in: aToken
  
¿Son dos operaciones o una?
  → Depende de la importación (API vs Etherscan CSV)
```

---

## Impacto en Auditoría

### Gas as Expense

```
Cada transaction en Ethereum gasta gas (ETH).

¿Cómo aparece?
  - Si CoinTracking lo detecta: "Gas Expense" (reduce ganancia)
  - Si NO lo detecta: ETH "desaparece" (error de saldo)
  
¿Es una ganancia/pérdida fiscal?
  - Técnicamente: SÍ (es dinero que gastas)
  - Fiscalmente (España): NO tenemos guía oficial
  - Buena práctica: Deducir como gasto de operación
```

### Contract Calls vs Transfers

```
¿Problema?: CoinTracking a veces ve:
  - Transfer out: 1.000 USDC (parece simple)
  
¿Realidad blockchain?
  - Llamada a Uniswap que incluye swap
  - Debería ser: 1.000 USDC → 1 ETH (Trade)
  
Validación:
  → Ver etherscan TX original
  → Confirmar si fue swap, lending, etc.
  → Reclasificar en CoinTracking si es necesario
```

---

## Validación en CoinTracking

```
Reports → Transacciones:
  Para cada TX de Ethereum:
    1. ¿Qué tipo es (transfer simple o contract call)?
    2. ¿El gas está registrado?
    3. ¿Las operaciones están relacionadas?
    
¿Sospecha?
  → Ir a Etherscan
  → Ver "To Address" (quién fue el destinatario)
  → Si es un contrato conocido (Uniswap, Aave) → Es contract call
  → Si es una dirección normal (wallet) → Es transfer simple
```

---

## Herramientas de Validación

```
1. Etherscan.io — Ver transacción on-chain
   https://etherscan.io/tx/<TXID>
   
   Información:
     - To: ¿Contrato o wallet?
     - Input data: ¿Qué función se llamó?
     - Gas used: Cuánto fue el gas real
   
2. CoinTracking → Transactions
   Enlace directo: [Ver en Etherscan]
   
3. Python script (local)
   - Validar que CoinTracking vea todos los transfers
   - Comparar contra etherscan export
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Tipos canónicos
- **DEFI_SWAPS_MECHANICS.md:** Swaps on-chain son contract calls
- **LENDING_MECHANICS.md:** Lending es interacción con contrato
