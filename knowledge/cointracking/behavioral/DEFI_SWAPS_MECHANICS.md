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
id: KB-B1-004
title: "Cómo CoinTracking maneja DeFi Swaps y operaciones complejas"
level: B
domain: cointracking
source: "Casos reales CT-009 + análisis"
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
  - knowledge/blockchain/TRANSACTION_TYPES.md
  - CT-009-defi-swap-aparece-como-dos-transferencias-separadas.md

tags:
  - cointracking
  - defi
  - swaps
  - behavioral

notes: "Operativo: cómo CoinTracking registra swaps de DeFi; cómo relacionarlos."
---

# Cómo CoinTracking Maneja DeFi Swaps

## Definición

**DeFi Swap** = Intercambio directo de un token por otro sin intermediario (p. ej. Uniswap: 10 USDC → 1 ETH).

**Características:**
- Ocurre ON-CHAIN (blockchain)
- Usualmente en un mismo bloque (transacción única que se divide en 2)
- CoinTracking lo ve como 2 operaciones separadas (salida + entrada)

---

## Problema Común: CT-009

**Síntoma:** Swap registrado como "2 transferencias independientes" sin relación.

```
CoinTracking muestra:
  - Transfer out: 10 USDC → dirección A (parece pérdida)
  - Transfer in: 1 ETH ← dirección A (parece ganancia inesperada)
  
Usuario piensa: "¿Dónde desapareció el USDC?"
```

**Causa:** CoinTracking importa el swap ON-CHAIN y ve 2 movimientos separados sin contexto de que son un intercambio.

**Impacto fiscal:** Puede parecer un fallo si no se documenta que fue un swap (DeFi trade).

---

## Cómo CoinTracking Registra Swaps

### Importación ideal

```
Si importas vía CoinTracking + API de Uniswap/1Inch:
  - Detecta automáticamente que es un swap
  - Tipo: "Trade" (Uniswap) o "Exchange" (si lo soporta)
  - Sale: 10 USDC
  - Entra: 1 ETH
  - Relacionados entre sí
  
Resultado: La ganancia/pérdida del swap se calcula correctamente
```

### Importación manual o vía CSV

```
Si importas el CSV desde etherscan (manual):
  CoinTracking ve:
    - Transfer out: 10 USDC (parece salida simple)
    - Transfer in: 1 ETH (parece entrada simple)
  
CoinTracking NO sabe que es un swap → no lo relaciona
  → Balance puede verse confuso
  → Ganancias no se calculan correctamente
```

---

## Validación y Corrección

### Identificar swaps mal registrados

```
CoinTracking → Transacciones:
  Busca:
    - Transfer out (token A)
    - Transfer in (token B)
    - MISMA FECHA + MISMA HORA (o segundos de diferencia)
    - MISMA DIRECCIÓN de origen/destino
  
¿Lo encontraste?
    → Probablemente es un swap que no fue clasificado
```

### Corregir

```
Opción 1: Editar para relacionar manualmente
  En CoinTracking:
    - Editar primero: "Transfer out 10 USDC"
    - Cambiar Tipo: "Trade" o "Exchange"
    - Cambiar "Received currency": ETH
    - Cambiar "Received amount": 1
  
  CoinTracking automáticamente:
    - Borra la "Transfer in" duplicada
    - Crea una sola operación: "Trade 10 USDC → 1 ETH"
  
Opción 2: Reimportar desde Uniswap API
  - Uniswap exporta swaps con el formato correcto
  - CoinTracking puede importarla y clasificarla bien
```

---

## Tratamiento Fiscal (España, IRPF)

**DeFi Swap = Trade (intercambio = venta + compra)**

```
Regla (DGT):
  - El momento: cuando se ejecuta el swap (on-chain timestamp)
  - El USDC: se vende al precio de momento del swap
  - El ETH: se compra al precio de momento del swap
  - Ganancia/pérdida: diferencia entre valor salida vs entrada
  
Fórmula:
  Ganancia = (Valor de entrada - Valor de salida)
  Ejemplo: Si vendo 10 USDC (10€) y recibo 1 ETH (2.000€)
    → Ganancia = 2.000€ - 10€ = 1.990€ (ganancia patrimonial)
```

**Ejemplo (CRÍT IMPORTANTE):**
```
15 enero 2025 - Swap en Uniswap:
  - Salida: 10 USDC (precio: 1€/USDC = 10€ costo)
  - Entrada: 1 ETH (precio: 2.000€/ETH)
  
CoinTracking debe calcular:
  - Venta de USDC: -10€ (recuperación de base)
  - Compra de ETH: 2.000€ (nueva base de coste)
  - Ganancia: 2.000€ - 10€ = 1.990€

IRPF 2025:
  - Ganancia patrimonial: 1.990€
```

---

## Validación en CoinTracking

```
Reports → Gains:
  ¿El swap aparece como un intercambio (Trade)?
    SÍ → OK
    NO → Verificar y corregir manualmente
    
Check: ¿La ganancia es sensata?
  - Si entraste 1 ETH a 2.000€ y saliste 10 USDC a 10€
  - Ganancia debe ser ~1.990€
  - Si aparece 0 o muy diferente → error en clasificación
```

---

## Caso Especial: Slippage y Fees

**Slippage:** diferencia entre precio esperado y precio real en el swap.

```
Ejemplo:
  - Esperabas: 10 USDC → 1.01 ETH
  - Recibiste: 10 USDC → 0.99 ETH (por slippage)
  
¿Cómo aparece en CoinTracking?
  - CoinTracking registra lo REAL: 10 USDC → 0.99 ETH
  - La "pérdida" por slippage es: 0.02 ETH × precio_ETH
  - Se calcula como ganancia/pérdida del trade
```

**Fees:** comisión de la DeFi (gas + protocolo).

```
Ejemplo:
  - Swap: 10 USDC → 1 ETH
  - Gas + Protocolo: 0.001 ETH
  - Recibiste realmente: 0.999 ETH
  
¿Cómo registrar?
  - Opción 1: CoinTracking calcula automáticamente (si importa desde DEX)
  - Opción 2: Añadir manualmente: Fee = 0.001 ETH
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Swap es Trade (venta + compra)
- **TRANSACTION_TYPES.md:** Tipos de transacciones blockchain
- **CAPITAL_GAINS.md:** Cálculo de ganancias en swaps
