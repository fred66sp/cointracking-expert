---
id: KB-B2-004
title: "Cómo CoinTracking maneja Binance Earn (productos de rendimiento)"
level: B
domain: cointracking
source: "Casos reales + análisis"
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
  - knowledge/cointracking/behavioral/LENDING_MECHANICS.md
  - knowledge/taxation/spain/CAPITAL_INCOME.md

tags:
  - cointracking
  - binance
  - earn
  - staking
  - behavioral

notes: "Operativo: cómo CoinTracking registra productos Earn de Binance."
---




# Cómo CoinTracking Maneja Binance Earn

## Definición

**Binance Earn** = Familia de productos de rendimiento:

```
- Flexible Earn (staking flexible)
- Locked Staking (staking bloqueado, mejor APY)
- Dual Investment (renta fija + opción sobre cripto)
- BETH (recibe interés en forma de tokens)
```

---

## Problema Común

**Síntoma:** Earn registrado como "Depósito genérico" sin los intereses.

```
CoinTracking ve:
  1. Deposit: 5 ETH (depósito inicial)
  2. ???: No aparecen los intereses
  
Usuario piensa: "¿Dónde están mis ganancias?"
```

**Causa:** Binance no exporta los intereses automáticamente; hay que importarlos manualmente.

---

## Cómo CoinTracking Registra Earn

### Flujo ideal

```
Binance Earn → CoinTracking API:
  1. Depósito inicial: 5 ETH
     Tipo: "Deposit" (a Earn)
  
  2. Cada día/semana: Interés acumulado
     Tipo: "Income" o "Reward"
     Ejemplo: +0.01 ETH
  
  3. Final del período (locked): Retiro
     Tipo: "Withdrawal"
     Recibiste: 5 + (intereses) ETH
```

### Flujo problemático

```
Si no conectas API a Binance:
  CoinTracking ve solo:
    1. Depósito
    2. (NADA)
    3. Retiro
  
Los intereses desaparecen de CoinTracking
  → Tax Report incompleto
```

---

## Validación y Corrección

### Identificar Earn mal registrado

```
CoinTracking → Transacciones:
  Busca depósitos que dicen "Earn" o "Staking"
  
¿Tiene intereses asociados después?
  SÍ → OK
  NO → Hay que añadirlos manualmente
```

### Corregir

```
Opción 1: Conectar API de Binance
  - CoinTracking → Settings → Exchanges
  - Agregar Binance (con API keys)
  - Reimportar histórico
  - Los intereses se cargan automáticamente
  
Opción 2: Importar manualmente desde Binance
  - Binance → Earn → Historial de recompensas
  - Exportar como CSV
  - Importar en CoinTracking
  
Opción 3: Añadir manualmente (tedioso)
  - Sumar intereses de cada período
  - Crear manualmente operaciones "Income"
```

---

## Estructura de un Producto Earn

### Ejemplo: Flexible Earn de ETH

```
Inicio: 1 enero 2025
  - Depósito: 5 ETH (precio: 1.500€/ETH = 7.500€)
  
Acumulación (diaria):
  - 1-7 enero: +0.001 ETH = 1,50€
  - 8-14 enero: +0.001 ETH = 1,50€
  - ...

Retiro: 31 enero 2025
  - Retiras: 5.03 ETH (precio: 2.000€/ETH)
  - Valor: 5.03 × 2.000€ = 10.060€
  - Ganancia: +0.03 ETH = 60€ (en el precio de retiro)

CoinTracking debe registrar:
  - Depósito: 5 ETH @ 1.500€
  - Income: 0.03 ETH @ 2.000€ = 60€ (aproximado)
```

---

## Tratamiento Fiscal (España, IRPF)

**Earn = Rendimiento del capital**

```
Regla (DGT):
  - Los intereses se acumulan diariamente
  - Momento exigible: cuando se acredita el interés
  - Valuación: precio del activo en el momento
  
CRÍTICO: Cada interés pequeño es un "ingreso" separado
  - Si ganas 0.03 ETH durante el año
  - Son 365 ingresos pequeños (1 por día, aproximadamente)
  - PERO fiscalmente se suman en el IRPF total
```

**Ejemplo simplificado:**
```
Earn de 0.03 ETH durante 2025 @ 1.500€/ETH = 45€
  (Los 0.03 ETH se acumularon cuando el precio era ~1.500€)

IRPF 2025:
  - Rendimiento del capital: 45€
  - Se suma a otros ingresos/ganancias
```

---

## Caso Especial: BETH (Staking de ETH a través de Binance)

**BETH** = Token que representa ETH bloqueado en Ethereum 2.0 staking.

```
Característico:
  - Tienes BETH (1 BETH ≈ 1 ETH)
  - BETH acumula interés (el precio sube)
  - Si quieres retirar: BETH → ETH (1:1)
  
CoinTracking:
  - Trata BETH como un activo separado
  - Registra la compra de BETH
  - Los intereses aparecen como apreciación (no como Income)
  
PROBLEMA:
  - La apreciación de BETH es "rendimiento del capital"
  - Pero CoinTracking lo trata como "ganancia patrimonial"
  - Diferencia fiscal importante
```

---

## Validación en CoinTracking

```
Reports → Gains:
  ¿Los ingresos de Earn aparecen?
    SÍ → OK
    NO → Verificar que estén importados/añadidos
    
¿El total es sensato?
  Esperado ≈ (Capital × APY × días) / 365
  Ejemplo: 5 ETH × 5% APY × 365 días / 365 ≈ 0.25 ETH
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Earn es Deposit + Income
- **LENDING_MECHANICS.md:** Similitudes y diferencias
- **CAPITAL_INCOME.md:** Tratamiento fiscal
