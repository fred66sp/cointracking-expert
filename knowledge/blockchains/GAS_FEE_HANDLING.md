---
id: KB-B3-005
title: "Cómo manejar Gas Fees en CoinTracking"
level: B
domain: blockchains
source: "Análisis de casos reales blockchain"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-03
confidence: medium
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/blockchains/ETHEREUM_TRANSACTION_TYPES.md
  - knowledge/blockchains/BITCOIN_TRANSACTION_TYPES.md

tags:
  - gas-fees
  - blockchain
  - expenses
  - behavioral

notes: "Operativo: cómo registrar y validar gas fees en auditoría."
---




# Cómo Manejar Gas Fees en CoinTracking

## Definición: Gas Fee

**Gas Fee** = Comisión de blockchain por procesar una transacción.

```
Ejemplos:
  - Ethereum: Fee en ETH (puede ser mucho en congestion)
  - Bitcoin: Fee en BTC (miners lo reciben)
  - Polygon: Fee en MATIC (muy bajo, <1 centavo)
  - Solana: Fee en SOL (muy bajo)
  
Varía según:
  - Congestión de la red (Ethereum: 0.001 a 0.1 ETH)
  - Tamaño de la transacción (Bitcoin: más inputs = más fee)
  - Configuración manual del gas (Ethereum: "slow", "fast", "custom")
```

---

## Problema: Gas Fees en CoinTracking

### Problema 1: Gas no se registra automáticamente

```
¿Qué ve CoinTracking?

Ideal:
  TX: Enviar 1 ETH a 0xBob
    - Transfer: 1 ETH out
    - Gas: 0.001 ETH out (PERO CoinTracking puede no verlo)
    
Realidad:
  CoinTracking registra:
    - Transfer: 1 ETH out
    - Gas: ??? (puede no estar en el CSV/API)
    
Resultado:
  - Balance parece correcto (1.001 ETH desaparecidos)
  - PERO el valor de 0.001 ETH puede no estar contabilizado
  - Auditoría: "Pérdida de 0.001 ETH no explicada"
```

### Problema 2: Gas en cadena distinta

```
Caso: Swap en Uniswap (Ethereum)

TX sale:
  - Envías 1.000 USDC (dirección A)
  - Recibes 1 ETH (dirección A)
  - Gas: 0.02 ETH
  
¿Cómo aparece?

Esperado:
  TX1: USDC out -1.000
  TX2: ETH in +1
  TX3: Gas -0.02 ETH
  
Realidad (frecuente):
  TX1: USDC out -1.000
  TX2: ETH in +1
  TX3: Gas??? (no se ve, o está en otra TX)
  
Resultado:
  - USDC y ETH se registran
  - Gas puede estar en otra TX, o no registrado
  - Balance de ETH parece más alto de lo que es
```

---

## Cómo CoinTracking Registra Gas

### Ethereum (EVM chains)

```
Si importas por API:
  ✓ CoinTracking detecta automaticamente gas usado
  ✓ Lo registra como "Gas Expense" o "Fee"
  
Si importas por CSV (Etherscan):
  ⚠️ Depende del formato del CSV
  ⚠️ Si el CSV incluye "Gas Used" → Importa bien
  ⚠️ Si no → Tienes que añadir manualmente

Verificación:
  CoinTracking → Transacciones → Buscar por "Gas"
  ¿Aparecen gastos de gas?
    SÍ → OK
    NO → Hay que importar manualmente
```

### Bitcoin

```
Gas = "Fees" (comisiones a miners)

¿Cómo aparece?

Si importas por API (blockchain explorer):
  ✓ Bitcoin fees se registran
  
Si importas por CSV de wallet:
  ⚠️ Depende de la wallet y su exportación
  ⚠️ Algunos wallets NO incluyen fees

Validación:
  - Blockchain.com: Shows fee en cada TX
  - Comparar contra CoinTracking
```

### Solana, Polygon, BSC (bajo gas)

```
Gas = muy bajo (0.00001 a 0.001 en su token)

¿Por qué es problemático?
  - CoinTracking a veces lo ignora (es tan pequeño)
  - Pero suma si haces muchas TX
  
Ejemplo:
  100 TX en Polygon @ 0.001 MATIC cada una
  = 0.1 MATIC ($0.02) (prácticamente ignorable)
  
Validación:
  - Si auditas Polygon intensamente → necesitas precisión
  - Si solo tienes algunas TX → puede ser ignorable
```

---

## Registrar Gas en CoinTracking

### Opción 1: Importación Automática (Mejor)

```
Busca que CoinTracking importe gas automáticamente:

Paso 1: Conectar API de Ethereum (Infura, Alchemy, etc.)
  CoinTracking → Settings → Exchanges → Add Ethereum wallet
  (Usa tu dirección pública, no private key)
  
Resultado:
  - CoinTracking importa todas las TX
  - Incluye gas automáticamente
  - Balance es preciso

O: Reimportar desde Etherscan CSV completo
  - Etherscan permite exportar CSV con "Gas Used"
  - Asegúrate de que incluya esa columna
```

### Opción 2: Añadir Gas Manualmente

```
Si el gas no se importó:

CoinTracking → Add Transaction:
  Tipo: "Expense" o "Fee"
  Cantidad: Gas en ETH/BTC/etc
  Precio: 0 (fue gasto, no compra)
  Fecha: MISMA FECHA que la TX original
  
Descripción: "Gas fee for TX: 0x123456..."
  
Resultado:
  - Se resta del balance de ETH
  - Cuenta como "gasto" en auditoría
```

### Opción 3: Incluir Gas en la TX Principal

```
Si tienes una TX de swap:
  Sale: 1.000 USDC
  Entra: 1 ETH
  Gas: 0.02 ETH
  
Opción:
  Editar la TX en CoinTracking
  → Añadir "Fee" o "Comisión"
  → Cantidad: 0.02 ETH
  
Resultado:
  - Una sola TX (más limpio)
  - Contabiliza gas en la ganancia/pérdida del swap
```

---

## Tratamiento Fiscal (España, IRPF)

**Gas fees = Gasto de operación (deducible teóricamente)**

```
Regla (no oficial, pero lógica):
  - Gas fee es dinero que gasta para operar
  - Debería reducir ganancia patrimonial
  
Ejemplo:
  Vendes 1 BTC a 60.000€ (compraste a 40.000€)
  Gas fee: 0.001 BTC = 60€
  
  Ganancia nominal: 60.000€ - 40.000€ = 20.000€
  Ganancia neta (si deduces gas): 20.000€ - 60€ = 19.940€
  
Realidad fiscal (España):
  - NO hay guía oficial de DGT sobre deducción de gas
  - Algunos asesores lo permiten, otros no
  - Enfoque conservador: NO deducir gas
  
RECOMENDACIÓN:
  - Documentar TODOS los gas fees
  - En la auditoría, mostrarlos separados
  - El asesor fiscal decide si deducir o no
```

---

## Validación en CoinTracking

```
Para auditar gas fees:

Reports → Transactions:
  [ ] Buscar "Gas" o "Fee"
  
  Para CADA transacción blockchain:
    [ ] ¿Hay gasto de gas registrado?
    [ ] ¿El monto es sensato?
    [ ] ¿La fecha coincide con la TX?
    
  Total de gas fees:
    [ ] Sumar todos los gastos
    [ ] Comparar contra blockchain explorer
    [ ] ¿Son similares?
      SÍ → OK
      NO → Hay gas fees faltantes
    
Corrección:
  [ ] Si faltan gas fees → Añadir manualmente
  [ ] Si hay duplicados → Eliminar
```

---

## Herramientas de Verificación

```
Ethereum:
  https://etherscan.io/tx/<TX_HASH>
  → Buscar "Gas Used" y "Gas Price"
  → Multiplicar: Gas Used × Gas Price = Fee en Wei
  → Convertir a ETH (1 ETH = 10^18 Wei)

Bitcoin:
  https://blockchain.com/btc/tx/<TX_HASH>
  → Buscar "Fee" en BTC
  → Comparar con CoinTracking

Polygon/BSC:
  https://polygonscan.com/tx/<TX_HASH>
  https://bscscan.com/tx/<TX_HASH>
  → Buscar "Gas Used" y "Gas Price"
  → Mismo cálculo que Ethereum
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Gas fees son expenses
- **ETHEREUM_TRANSACTION_TYPES.md:** Gas en Ethereum
- **BITCOIN_TRANSACTION_TYPES.md:** Fees en Bitcoin
