---
id: KB-B2-005
title: "Cómo CoinTracking maneja Binance Convert (intercambios rápidos)"
level: B
domain: cointracking
source: "Casos reales + análisis"
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
  - knowledge/taxonomy/spain/CAPITAL_GAINS.md
  - knowledge/cointracking/behavioral/DEFI_SWAPS_MECHANICS.md

tags:
  - cointracking
  - binance
  - convert
  - behavioral

notes: "Operativo: cómo CoinTracking registra Binance Convert."
---

# Cómo CoinTracking Maneja Binance Convert

## Definición

**Binance Convert** = Intercambio rápido entre criptos directamente en Binance.

**Equivalente:** Like a DEX (Uniswap) pero centralizado en Binance.

```
Ejemplo: 1 ETH → 50 USDT (al mejor precio disponible)
  - Es más rápido que Spot (sin orden en el libro)
  - Precio es mejor que Spot (menos spread)
```

---

## Cómo CoinTracking Registra Convert

### Flujo ideal

```
Si importas vía API de Binance:
  CoinTracking detecta el Convert
  
Tipo: "Trade" o "Exchange"
  - Salida: 1 ETH
  - Entrada: 50 USDT
  - Fee: incluido en el precio
  
Resultado: La ganancia/pérdida se calcula correctamente
```

### Flujo problemático

```
Si el CSV no especifica "Convert":
  CoinTracking lo registra como:
    - "Transfer out" (1 ETH desaparece)
    - "Transfer in" (50 USDT aparece)
  
Usuario piensa: "¿Dónde fue el ETH?"
```

---

## Validación y Corrección

### Identificar Convert mal clasificado

```
CoinTracking → Transacciones:
  Busca:
    - Transfer out (token A)
    - Transfer in (token B)
    - MISMA FECHA + MISMA HORA (segundos de diferencia)
    - Descripción menciona "Convert" o "Swap"
  
¿Lo encontraste?
    → Probablemente es un Convert que no fue clasificado
```

### Corregir

```
Opción 1: Editar para relacionar
  - Editar "Transfer out"
  - Cambiar Tipo: "Trade" o "Exchange"
  - Agregar "Received currency": Token B
  - Agregar "Received amount": cantidad
  
  CoinTracking automáticamente:
    - Borra la "Transfer in" duplicada
    - Crea una operación unificada
  
Opción 2: Reimportar desde Binance API
  - Reconectar API de Binance
  - Reimportar histórico
  - Convert se clasifica correctamente
```

---

## Estructura de un Convert

### Ejemplo: ETH a USDT

```
Binance Convert:
  - From: 1 ETH
  - To: 50 USDT (precio actual del mercado)
  - Fee: 0.25% (incluida)
  
Neto recibido: 50 × (1 - 0.0025) = 49.875 USDT

CoinTracking registra:
  - Sale: 1 ETH
  - Entra: 49.875 USDT
  
¿Ganancia o pérdida?
  Si compraste el ETH a 40.000€:
    - Valor venta: 49.875 USDT ≈ 49.875€
    - Costo: 40.000€ (el ETH original)
    - Ganancia: 49.875€ - 40.000€ = 9.875€ (aprox.)
```

---

## Diferencias: Convert vs Spot

```
SPOT:
  ✓ Asientos en el libro de órdenes
  ✓ Precio exacto depende del libro
  ✓ Puede ejecutar a tramos diferentes
  ✗ Más lento (esperar orden)

CONVERT:
  ✓ Precio inmediato (mejor que Spot)
  ✓ Sin orden pendiente
  ✓ Transacción instantánea
  ✗ Fee fija (0.25%)
  ✗ No ves el orden de ejecución
```

---

## Tratamiento Fiscal (España, IRPF)

**Convert = Trade (venta + compra, como Spot)**

```
Regla (igual a Spot):
  - Venta del token A al precio del Convert
  - Compra del token B al precio del Convert
  - Ganancia/pérdida: diferencia entre costo y precio
  
Ejemplo:
  Vendiste: 1 ETH a 50.000€
  Compraste: 50 USDT (≈50€)
  
Si el costo del ETH fue 40.000€:
    Ganancia = 50.000€ - 40.000€ = 10.000€
```

---

## Ventaja de Convert para Auditoría

```
✓ Más limpio que DeFi (no hay on-chain)
✓ Registrado automáticamente por Binance API
✓ No hay "transacciones fantasma" (como en swaps DeFi)
✓ Fee es transparente
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Convert es Trade
- **CAPITAL_GAINS.md:** Cálculo de ganancias (como Spot)
- **DEFI_SWAPS_MECHANICS.md:** Diferencias con swaps on-chain
