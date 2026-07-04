---
id: KB-B3-002
title: "Tipos de transacciones Bitcoin y cómo afectan a CoinTracking"
level: B
domain: blockchains
source: "Bitcoin whitepaper + casos reales"
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
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md

tags:
  - bitcoin
  - blockchain
  - transactions
  - utxo
  - behavioral

notes: "Operativo: tipos de tx en Bitcoin; UTXO model y cómo afecta auditoría."
---

# Tipos de Transacciones Bitcoin y su Impacto en CoinTracking

## Diferencia Crítica: Bitcoin vs Ethereum

```
ETHEREUM:
  - Modelo de "cuenta" (como banco)
  - Una dirección = un saldo
  - Transacción = cambio de saldo
  
BITCOIN:
  - Modelo de "UTXO" (Unspent Transaction Output)
  - Una dirección = lista de UTXO (monedas sin gastar)
  - Transacción = consumir algunos UTXO, crear otros nuevos
```

---

## Modelo UTXO (Unspent Transaction Output)

**Analogía:** Bitcoin es como dinero en efectivo.

```
Tienes en el bolsillo:
  - 1 billete de 10€
  - 1 billete de 5€
  - Total: 15€

Quieres pagar 12€. Opciones:
  a) Dar el billete de 10€ + 2€ cambio
  b) Dar el billete de 5€ + 5€ + 2€ cambio
  
En Bitcoin:
  a) Consumir UTXO de 1 BTC (0.5 BTC), recibir 0.4 BTC cambio
  b) Consumir UTXO de 0.3 BTC + 0.2 BTC, recibir 0.1 BTC cambio
```

---

## Tipo 1: Simple Send (Envío Simple)

```
En blockchain:
  De: 1A2B3C... (mi dirección)
  A: XYZ... (dirección de otro)
  Cantidad: 0.5 BTC
  Fee: 0.001 BTC
  
¿Qué pasó internamente?
  - Bitcoin seleccionó uno o más UTXO de mi cartera
  - Los "gastó" (consumió)
  - Creó un UTXO nuevo de 0.5 BTC para el destinatario
  - Creó un UTXO de cambio (lo que sobró, si lo hay)
  
CoinTracking ve:
  - BTC sent: 0.5 BTC
  - Fee: 0.001 BTC
  
Fácil de auditar: ✓
```

---

## Tipo 2: Consolidation (Juntar UTXO)

**Problema común en auditoría: "¿A quién se envió el BTC?"**

```
En blockchain:
  De: 1A2B3C... (MI DIRECCIÓN)
  A: 1D4E5F... (MI DIRECCIÓN, pero diferente)
  Cantidad: 2 BTC (total de varios UTXO)
  Fee: 0.002 BTC
  
¿Qué pasó?
  - Bitcoin consolidó 5 pequeños UTXO en 1 grande
  - Los envió a una dirección mía (misma wallet, derivada)
  - Razón: Las wallets multi-dirección generan muchos UTXO pequeños
  - Periodicamente hay que consolidarlos
  
CoinTracking ve:
  - BTC withdrawal: 2 BTC (parece que me lo robaron)
  - En algunos casos: BTC deposit: 2 BTC (suerte, recuperé el dinero)
  
PROBLEMA:
  - ¿Ganancia o pérdida? Ninguna (es mi dinero moviéndose)
  - ¿Cómo CoinTracking lo ve? Confuso (dos direcciones diferentes)
  
Fácil de auditar: ✗ Necesita verificación de que ambas direcciones son mías
```

---

## Tipo 3: Change Output (Moneda de Cambio)

**CRÍTICO: Puede parecer "pérdida" en CoinTracking**

```
En blockchain:
  Tengo UTXO de 2 BTC
  Quiero enviar 0.5 BTC
  
Transacción Bitcoin:
  Input: 2 BTC (lo que gasto)
  Output 1: 0.5 BTC (destinatario)
  Output 2: 1.499 BTC (cambio a mi dirección)
  Fee: 0.001 BTC
  
CoinTracking puede ver:
  Opción A (correcta):
    - BTC sent: 0.5 BTC
    - Fee: 0.001 BTC
    
  Opción B (incorrecta):
    - BTC sent: 2 BTC (confunde input con output)
    - BTC received: 1.499 BTC (pero no lo conecta)
    - Resultado: Pérdida de 0.5 BTC (¡FALSO!)
    
PELIGRO: Balance negativo si CoinTracking ve solo output de cambio
```

---

## Validación en CoinTracking

```
Bitcoin transactions son más confusas que Ethereum porque:
  1. Las wallets pueden generar múltiples direcciones
  2. El "cambio" es un output igual a cualquier otro
  3. CoinTracking ve la transacción cruda, no la intención
  
Validación en CoinTracking:
  1. Cada BTC sent → ¿A quién fue?
  2. Si es a mi misma dirección → Consolidation (OK, sin ganancia)
  3. Si es a otra dirección → Verificar que sea un envío real
  
Herramientas:
  - Blockchain.com (explorer)
  - Tu wallet original (puede etiquetar direcciones)
```

---

## Tipo 4: Mining Rewards / Staking (Lightning Network, etc.)

```
En blockchain:
  - Un bloque fue minado/validado por alguien
  - Se crea un "coinbase transaction"
  - El minero recibe la recompensa + fees
  
¿CoinTracking lo ve?
  - Si tu wallet recibió la recompensa: depende de si conectaste API
  - Si fue en cambio de pool mining: depende del pool (pueden estar mezclados)
  
PROBLEMA:
  - Pool mining rewards a menudo llegan en un lote
  - CoinTracking las ve como un único BTC received
  - Pero en realidad eran 100 bloques pequeños
```

---

## Tratamiento Fiscal (España, IRPF)

**Bitcoin = Igual a Ethereum**

```
Ganancias patrimoniales:
  - Venta de BTC a mayor precio que compra → ganancia
  - Fee de transacción → gasto (deducible si está documentado)
  
Consolidations / cambios internos:
  - NO son ventas (no consumes base de coste)
  - NO generan ganancia/pérdida
  - Simplemente reorganizas tu activo
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Bitcoin UTXO model
- **ETHEREUM_TRANSACTION_TYPES.md:** Comparación de modelos
