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
id: KB-B2-007
title: "Cómo CoinTracking maneja Coinbase Advanced Trade y Staking"
level: B
domain: cointracking
source: "Análisis de casos Coinbase"
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
  - knowledge/cointracking/behavioral/BINANCE_SPOT_MECHANICS.md
  - knowledge/cointracking/behavioral/STAKING_MECHANICS.md

tags:
  - coinbase
  - advanced-trade
  - staking
  - behavioral

notes: "Operativo: cómo CoinTracking importa Coinbase Advanced Trading y rewards."
---

# Cómo CoinTracking Maneja Coinbase Advanced Trade y Staking

## Coinbase: Dos Plataformas

**IMPORTANTE:** Coinbase tiene dos interfases diferentes:

| Plataforma | Uso | Comisión | En CT |
|---|---|---|---|
| **Coinbase.com** | Principiantes | 1.5% | Básico |
| **Coinbase Advanced** | Traders | 0.1-0.6% | Recomendado |

CoinTracking detecta ambas, pero importa comisiones diferentes.

---

## Coinbase Advanced Trade

### Características

```
¿Qué es?
  - Interfaz profesional de Coinbase
  - Órdenes límite, market orders
  - Comisiones más bajas (0.4% típico)
  - API más robusta

¿Cómo aparece?

Ideal:
  - Compra: 1 BTC @ 50.000 USDC (comisión incluida)
  - Venta: 1 BTC @ 60.000 USDC (comisión deducida)

Realidad:
  CoinTracking puede ver:
    - La operación principal
    - La comisión separada (BTC o USDC)
    - O ambas combinadas
```

### Importación en CoinTracking

```
Método 1: API (Recomendado)
  CoinTracking → Settings → Exchanges → Add Coinbase
  Autorizar acceso a API
  
  Resultado:
    - Todas las operaciones importadas
    - Comisiones incluidas
    - En tiempo real

Método 2: CSV (Manual)
  Coinbase → Account → Download account data
  Incluye: trades, deposits, rewards
  
  Problema:
    - CSV puede no incluir comisiones claras
    - Requiere parseo manual
```

---

## Coinbase Staking (Ethereum, SOL, etc.)

### Cómo Funciona

```
Coinbase ofrece staking directo de activos

Ejemplo: Staking de ETH

Paso 1: Convertir ETH → stETH (Ethereum staking token)
  Automático en Coinbase (no visible como operación)
  
Paso 2: Recibir rewards
  Acumulan diariamente
  Aparecen como "Staking rewards"
  
Paso 3: Convertir de vuelta (opcional)
  stETH → ETH
  Puedes hacerlo en cualquier momento (liquid staking)

¿Cómo aparece en CoinTracking?

Ideal:
  TX1: Staking 32 ETH
  TX2: Income +0.5 ETH/día
  TX3: Unstaking si quieres retirar

Realidad:
  CoinTracking puede ver:
    - Cambio de saldo ETH → stETH
    - Rewards como "Income"
    - Pero conexión puede ser confusa
```

### Problema: Convertir ETH → stETH

```
¿Qué es stETH?

En Coinbase:
  - 1 ETH = 1 stETH (aproximadamente)
  - Pero son activos diferentes
  - stETH acumula valor (incluye rewards)
  
En CoinTracking:
  ETH y stETH aparecen como activos separados
  
¿Impacto?

Si CoinTracking ve:
  TX: Exchange 32 ETH → 32 stETH
  
¿Es una venta?
  NO (mismo valor)
  Pero CoinTracking puede clasificarlo como Trade
  
Cálculo de ganancia:
  Si ET y stETH tienen precios diferentes:
    Ganancia = 32 stETH × precio_stETH - 32 ETH × precio_ETH
    
  Pero después de staking, stETH vale más:
    Ganancia real: Los rewards acumulados
```

---

## Validación en CoinTracking

### Verificar Importación

```
CoinTracking → Transacciones:
  Filtrar por "Coinbase"
  
[ ] ¿Aparecen operaciones de trading?
    SÍ → OK
    NO → Verificar API
    
[ ] ¿Las comisiones están incluidas?
    Comparar con Coinbase:
      CT price ≈ Coinbase "Executed price" + comisión
      
[ ] ¿Staking rewards aparecen?
    Filtrar: "Income" + "Coinbase"
    ¿Tienen valor sensato?
```

### Verificar Staking

```
En Coinbase:
  Account → Staking → Rewards history
  Anotar: Total de rewards
  
En CoinTracking:
  Filtrar: "Income" + "Staking"
  Sumar valores
  
Comparar:
  ¿Coinbase total ≈ CT total?
  
Si no:
  [ ] Buscar rewards faltantes
  [ ] Importar manualmente si faltan
```

---

## Registrar Comisiones Correctamente

### Opción 1: En el Precio (Mejor)

```
Compro 1 BTC @ 50.000 USDC con comisión de 200 USDC

En CoinTracking:
  Tipo: Buy
  Cantidad: 1 BTC
  Precio: 50.200 USDC/BTC (INCLUYE comisión)
  
Resultado:
  - Cost basis: 50.200 USDC
  - Una sola línea (limpio)
  - Ganancia calculada correcta
```

### Opción 2: Como Fee Separado

```
Compro 1 BTC @ 50.000 USDC con comisión de 200 USDC

En CoinTracking:
  TX1: Buy 1 BTC @ 50.000 USDC
  TX2: Fee -200 USDC (o -0.004 BTC equivalente)
  
Problema:
  - Cost basis: 50.000 USDC (incorrecto)
  - Fee aparece separado (confuso)
  
Solución: Editar TX1 después
  Añadir "Commission": 200 USDC
  CoinTracking ajusta cost basis automáticamente
```

---

## Tratamiento Fiscal (España, IRPF)

**Coinbase Trading = Igual que Binance**

```
Regla:
  - Venta de cripto = ganancia patrimonial
  - FIFO para determinar cost basis
  - Comisiones reducen ganancia (teóricamente)
```

**Coinbase Staking Rewards = Rendimiento del Capital**

```
Regla:
  - Cada reward es rendimiento en el día que se acredita
  - Valuación: precio ETH/SOL ese día
  - Se suma a "Rendimientos"
  
Ejemplo:
  Enero: +0.01 ETH reward @ 1.500€ = +15€
  Febrero: +0.01 ETH reward @ 2.000€ = +20€
  Total IRPF 2025: 35€
```

---

## Coinbase Conversion: Automatic Staking

**Coinbase ofrece "Auto-convert" de ciertos activos a staking**

```
Ejemplo: USDC → stUSDC (acumula interés)

¿Qué hace?
  Coinbase automáticamente convierte USDC en stUSDC
  Recibe ~1.5% APY en forma de stUSDC

¿Cómo aparece en CT?

Ideal:
  TX: "Convert USDC → stUSDC for staking"

Realidad:
  CoinTracking ve:
    - Balance de USDC baja
    - Balance de stUSDC sube
    - Puede parecer un trade (perdida/ganancia)

Impacto fiscal:
  - NO es una venta (es solo conversión)
  - Los rewards son ingresos (rendimiento)
```

---

## Validación: Checklist Coinbase

```
Antes de auditar:

[ ] ¿Coinbase Advanced API conectado?
    SÍ → OK
    NO → Conectar
    
[ ] ¿Comisiones incluidas en precio?
    Comparar CT vs Coinbase
    
[ ] ¿Staking rewards importados?
    Comparar: Coinbase vs CT
    
[ ] ¿stETH/stUSDC tratados como activos separados?
    Verificar balance
    
[ ] ¿Conversiones son correctas?
    ETH → stETH no debería ser ganancia/pérdida
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Advanced Trade es igual a Spot
- **BINANCE_SPOT_MECHANICS.md:** Comparación de comisiones
- **STAKING_MECHANICS.md:** Rewards y tratamiento fiscal
