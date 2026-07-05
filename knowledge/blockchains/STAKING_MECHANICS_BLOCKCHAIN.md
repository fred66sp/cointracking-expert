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
id: KB-B3-006
title: "Cómo funcionan Staking y Validadores en Blockchain (on-chain)"
level: B
domain: blockchains
source: "Análisis de Ethereum 2.0, Solana, Cosmos"
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
  - knowledge/cointracking/behavioral/STAKING_MECHANICS.md
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md

tags:
  - staking
  - validators
  - blockchain
  - on-chain
  - behavioral

notes: "Operativo: cómo funcionan validadores y rewards en blockchain."
---

# Cómo Funcionan Staking y Validadores en Blockchain

## Concepto Base

**Staking** = Bloquear cripto en un validador para asegurar la red.

```
Equivalente financiero:
  - Pones dinero en un depósito (bloqueado)
  - El banco te paga interés
  - Si el banco quiebra, pierdes el dinero
  
En blockchain:
  - Bloqueas cripto en un validador
  - Recibes rewards (inflación de la red)
  - Si el validador actúa mal, pierdes cripto ("slashing")
```

---

## Tipos de Staking

### Tipo 1: Proof of Stake (PoS) — Ethereum 2.0

```
Ejemplo: Ethereum 2.0 (post-Merge 2022)

Requisitos:
  - Mínimo: 32 ETH bloqueados
  - Hardware: Nodo completo (o pool)
  - Tiempo: Indefinido (lock-up indefinido)
  
Rewards:
  - ~3-4% APY (variable)
  - Pagado en ETH
  - Se acumula automáticamente
  
Riesgo:
  - "Slashing": Si el validador actúa mal, pierdes parte del stake
  - Raro (~0.1% anual si sigues reglas)
  
¿Cómo aparece?

On-chain:
  TX1: Deposit 32 ETH → Staking contract
  TX2: Reward +0.02 ETH (cada slot, ~12 segundos)
  
En CoinTracking:
  - Depósito: -32 ETH (enviado a contrato)
  - Rewards: +0.02 ETH (cada ~12 seg, puede agruparse)
```

### Tipo 2: Delegated Proof of Stake (DPoS) — Solana, Cosmos

```
Ejemplo: Solana

Requisitos:
  - Mínimo: 0.1 SOL (muy bajo)
  - Hardware: Nada (delegas a un validador)
  - Tiempo: 24-48h para activarse
  
Rewards:
  - ~5-10% APY
  - Pagado en SOL
  - Se acumula automáticamente
  
Ventaja vs PoS:
  - Flexible (no hay lock-up obligatorio)
  - Escalable (muchos validadores)
  
¿Cómo aparece?

On-chain:
  TX1: Delegate 10 SOL → Validator
  TX2: Reward +0.003 SOL (cada época ~2 días)
  
En CoinTracking:
  - Delegación: -10 SOL (enviado a validator account)
  - Rewards: +0.003 SOL (periódicamente)
```

### Tipo 3: Liquid Staking (stETH, stSOL, etc.)

```
Nuevo modelo: Staking sin lock-up

Ejemplo: Lido (stETH en Ethereum)

¿Cómo funciona?
  1. Depositas 32 ETH en Lido
  2. Lido mina stETH (token de recibo)
  3. Lido lo valida por ti (mucho más seguro)
  4. stETH acumula valor (rewards incluidos)
  5. Puedes vender stETH en cualquier momento
  
Ventaja:
  - Liquidity (no está bloqueado como ETH directo)
  - Rewards automáticos (stETH = valor más alto)
  
¿Cómo aparece?

On-chain:
  TX1: Deposit 32 ETH → Lido contract
  TX2: Mint 32 stETH (1:1 inicialmente)
  TX3: stETH value aumenta (rewards incluidos)
  
En CoinTracking:
  - TX1: -32 ETH (enviado a Lido)
  - TX2: +32 stETH (recibido)
  - Rewards: Implícitos (stETH vale más después)
```

---

## Validadores: Slashing y Penalizaciones

**RIESGO IMPORTANTE**

```
¿Qué es slashing?

Si un validador actúa mal (ataca la red), pierde cripto

Ejemplos:
  - Doble firma (firmar dos bloques conflictivos)
  - Ataque de 51%
  - Inactividad prolongada
  
Cantidad perdida:
  - Ethereum: 0.5% a 100% (depende de la falta)
  - Solana: 0.5% a 5% (depende del validador)
  
Impacto fiscal:
  - Si pierdes ETH por slashing: ES PÉRDIDA
  - Reportar en IRPF como "pérdida patrimonial"
  
En CoinTracking:
  - Slashing aparece como "penalty" o "fee"
  - Debe reportarse correctamente
```

---

## Rewards en CoinTracking (On-Chain)

### Cómo se Registran

```
Si importas la blockchain directamente (p. ej., Etherscan):

CoinTracking ve:
  TX1: Sent ETH -32 (depósito en contrato)
  TX2: Received ETH +0.02 (reward, puede haber muchas)
  TX3: Received ETH +0.02
  ...
  TX100: Received ETH +0.02
  
Resultado:
  - 100 transacciones por 1 depósito + rewards
  - Balance correcto al final (32 + 2 = 34 ETH)
  - Pero muy tedioso de auditar

Solución:
  Agrupar rewards manualmente:
    1 TX: Income +100 ETH (total de rewards en período)
    Precio: Promedio de ETH en el período
```

### Agregar Rewards Manualmente

```
Si tienes muchos rewards pequeños:

CoinTracking → Add Transaction:
  Tipo: Income
  Activo: ETH
  Cantidad: Total de rewards en período (ej. 2 ETH)
  Precio: Promedio de ETH (ej. 1.500€)
  Fecha: Fecha promedio (ej. 30 junio)
  Descripción: "Staking rewards Ethereum 2.0 enero-junio"
  
Impacto:
  - Una transacción en lugar de 100
  - Menos tedioso
  - Menos preciso (pero aceptable)
```

---

## Tratamiento Fiscal (España, IRPF)

**Rewards en Blockchain = Rendimiento del Capital**

```
Regla:
  - Momento: Cuando la blockchain acredita el reward
  - Valuación: Precio de la cripto en ese momento
  - Impacto: Se suma a "Rendimientos del capital"
  
Ejemplo: Ethereum Staking
  Depositas: 32 ETH @ 1.500€ = 48.000€ (enero)
  Recibes rewards: 0.1 ETH/mes @ 1.500€ (promedio)
    = 150€/mes × 12 = 1.800€/año
    
  IRPF 2025:
    - Rendimientos: 1.800€
```

**Slashing = Pérdida Patrimonial**

```
Si sufres slashing:
  Pierdes: 0.5 ETH (ejemplo)
  Precio: 2.000€/ETH
  Pérdida: 1.000€
  
IRPF 2025:
  - Pérdida patrimonial: -1.000€
  - Puede compensar otras ganancias
```

---

## Validación: Verificar Staking On-Chain

```
Para auditar staking en blockchain:

Paso 1: Etherscan (Ethereum)
  Busca tu dirección
  Filtra: "To: Staking contract"
  ¿Ves el depósito de 32 ETH?
  
Paso 2: Rewards History
  Etherscan → "Internal Transactions"
  Busca: "Received from Staking contract"
  ¿Cuántos rewards? ¿Cuánto total?
  
Paso 3: Comparar con CoinTracking
  ¿CT muestra el mismo total?
  
Paso 4: Slashing Check
  ¿Hay transacciones de "penalty"?
  Si sí, documentar para IRPF
```

---

## Caso Especial: Staking Pool vs Self-Staking

| Aspecto | Self-Staking (32 ETH) | Pool (Cualquier cantidad) |
|---|---|---|
| **Mínimo** | 32 ETH | 0.1 ETH |
| **Hardware** | Nodo propio | Nada |
| **Rewards** | 100% | Menos (comisión del pool) |
| **Riesgo slashing** | Mayor (tu responsabilidad) | Menor (pool lo maneja) |
| **Liquidez** | 0 (bloqueado) | Puede ser líquido (stETH) |

**En CoinTracking:**
  - Self: Depósitos simples, rewards on-chain
  - Pool: Depósitos a exchange/Lido, rewards más claros

---

## Integración

- **ADR-003:** Modelo de transacciones — Staking son depósitos + rewards
- **STAKING_MECHANICS.md:** Staking en exchanges
- **ETHEREUM_TRANSACTION_TYPES.md:** Contratos y transacciones internas
