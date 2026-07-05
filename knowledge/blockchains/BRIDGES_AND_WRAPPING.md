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
id: KB-B3-004
title: "Bridges y Wrapped Tokens en CoinTracking"
level: B
domain: blockchains
source: "Análisis de casos multichain"
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
  - knowledge/blockchains/OTHER_CHAINS_MECHANICS.md
  - knowledge/cointracking/behavioral/MISSING_PURCHASE_HISTORY_CAUSES.md

tags:
  - bridges
  - wrapped-tokens
  - multichain
  - behavioral

notes: "Operativo: cómo manejar bridges y wrapped tokens en auditoría."
---

# Bridges y Wrapped Tokens en CoinTracking

## Definición: Bridge

**Bridge** = Mecanismo para transferir un token de una blockchain a otra.

```
Ejemplo: Enviar USDC de Ethereum a Polygon

Proceso:
  1. Bloqueas 100 USDC en Ethereum (smart contract)
  2. El bridge mina 100 USDC.e en Polygon (wrapped)
  3. Ahora tienes 100 USDC.e en Polygon
  4. Si quieres volver: quema USDC.e en Polygon → desbloquea en Ethereum

Resultado: MISMO DINERO en dos blockchains diferentes
```

---

## Definición: Wrapped Token

**Wrapped Token** = Representación de un token en una blockchain diferente.

```
Ejemplos:
  - USDC en Ethereum → USDC.e en Polygon
  - USDC en Ethereum → USDC en Arbitrum (versión oficial)
  - ETH en Bitcoin (raro) → WETH en Ethereum
  - BTC en Ethereum → WBTC en Ethereum

Características:
  - 1 wrapped = 1 original (1:1)
  - NO son el mismo activo técnicamente (blockchains diferentes)
  - Pero SÍ representan el mismo valor
```

---

## Problema en CoinTracking

```
¿Cómo CoinTracking ve un bridge?

Ideal:
  Transacción: "Bridge USDC Ethereum → Polygon"
  Entrada: 100 USDC.e
  Salida: 100 USDC
  
Realidad:
  TX1: Withdraw USDC (Ethereum) -100
  TX2: Deposit USDC.e (Polygon) +100
  
Problema:
  - CoinTracking ve DOS transacciones separadas
  - No conecta que es un bridge
  - Balance parece confuso
  - "Costo" de USDC vs USDC.e puede diferir
```

---

## Validación en CoinTracking

### Identificar bridges mal registrados

```
Síntoma:
  - USDC desaparece en una fecha
  - USDC.e aparece en otra fecha (casi igual)
  - Timestamps diferentes, pero cercanos (segundos)
  
Verificación:
  1. ¿Ambas transacciones tienen la MISMA cantidad?
     SÍ → Probable bridge
     NO → Probable venta/compra, no bridge
     
  2. ¿La fecha/hora es cercana (mismo día, pocas horas)?
     SÍ → Probable bridge
     NO → Operaciones separadas
     
  3. ¿Vienen de la misma dirección (tu wallet)?
     SÍ → Probable bridge
     NO → Transferencia a otro, no bridge
```

### Corregir el bridge en CoinTracking

```
Opción 1: Editar manualmente para relacionar
  - Editar "USDC Withdraw"
  - Cambiar tipo: "Transfer"
  - Añadir nota: "Bridge a Polygon"
  
  CoinTracking automáticamente puede:
    - Vincular con el "USDC.e Deposit"
    - O dejarlas separadas pero etiquetadas

Opción 2: Usar CoinTracking Bridge Tool (si existe)
  - CoinTracking → Tools → Bridges
  - Especificar: "100 USDC Ethereum → 100 USDC.e Polygon"
  - Resultado: se relacionan automáticamente

Opción 3: Crear una transacción "Exchange"
  - Eliminar las dos separadas
  - Crear una: "Exchange 100 USDC → 100 USDC.e"
```

---

## Caso Especial: Wrapped Tokens Permanentes

**Situación: Tienes WBTC (BTC wrapped en Ethereum)**

```
¿Qué es WBTC?
  - 1 WBTC = 1 BTC en valor
  - Pero es un token ERC-20 (no es BTC real)
  - Vive en Ethereum
  
¿Cómo aparece en CoinTracking?
  
  Forma correcta:
    - WBTC es un activo SEPARADO de BTC
    - No confundir: WBTC ≠ BTC
    - Si vendes WBTC, vendes el wrapped, no BTC real
    
  Problema común:
    - CoinTracking a veces agrupa WBTC con BTC
    - Balance total muestra BTC + WBTC mezclados
    - Es incorrecto
```

---

## Tratamiento Fiscal (España, IRPF)

**Bridges NO crean ganancia/pérdida (son transferencias internas)**

```
Regla:
  - Bridge USDC Ethereum → USDC.e Polygon = NO venta
  - Cost basis se mantiene igual
  - NO hay ganancia patrimonial

Excepción:
  - Si el precio del token varió entre Ethereum y Polygon
  - Ejemplo: USDC vale 0.98€ en Polygon (por arbitraje)
  - Entonces SÍ hay ganancia técnicamente
  - Pero en auditoría práctica: se trata como transferencia sin ganancia

Wrapped tokens (WBTC, WETH):
  - Tienen su propio precio en Ethereum
  - Si WBTC vale menos de BTC → pérdida al wrappear (raro)
  - Si WBTC vale igual a BTC → sin ganancia
```

---

## Multichain Complexity: Mismo Token en Múltiples Chains

```
Escenario complejo:
  - Tengo 1.000 USDC en Ethereum
  - Tengo 500 USDC.e en Polygon
  - Tengo 200 USDC en Arbitrum
  - ¿Total USDC?

En CoinTracking:
  Ideal: Muestra TRES activos diferentes
    - USDC (1.000€)
    - USDC.e (500€)
    - USDC Arbitrum (200€)
    - Total: 1.700€ (pero NO son directamente intercambiables)
  
  Problema: CoinTracking agrupa como "USDC total 1.700"
    - Correcto en valor total
    - Incorrecto en liquidez (no puedo usar USDC.e en Ethereum directamente)

Validación:
  → Verificar que cada token esté etiquetado con su chain
  → No asumir que "total USDC" es correcta sin desglosar por chain
```

---

## Validación en CoinTracking

```
Si tienes multichain/wrapped:

Reports → Transactions:
  [ ] Cada token tiene etiqueta de chain?
      ✓ USDC (Ethereum)
      ✓ USDC.e (Polygon)
      ✓ WBTC (Ethereum)
      SÍ → OK
      NO → Etiquetar manualmente
  
  [ ] Los bridges están claramente marcados?
      SÍ → OK
      NO → Añadir notas

  [ ] Cost basis es consistente por chain?
      SÍ → OK
      NO → Revisar base de coste
```

---

## Herramientas Externas

```
Para verificar bridges:
  - Etherscan: Bridge transactions en Ethereum
  - Polygonscan: Depósitos de bridge en Polygon
  - Arbitrumone.org: Transacciones en Arbitrum
  
Buscar:
  - Bridge contract address
  - Timestamps coinciden?
  - Cantidades son iguales (menos fees)?
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Bridges son transferencias (sin venta)
- **OTHER_CHAINS_MECHANICS.md:** Multichain overview
- **MISSING_PURCHASE_HISTORY_CAUSES.md:** Bridges pueden causar cost basis confuso
