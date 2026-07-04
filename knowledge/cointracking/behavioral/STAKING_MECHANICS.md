---
id: KB-B1-001
title: "Cómo CoinTracking maneja Staking y Rewards"
level: B
domain: cointracking
source: "Análisis de casos reales + centro de ayuda CT"
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
  - knowledge/taxonomy/spain/CAPITAL_INCOME.md
  - CT-005-recompensas-de-staking-clasificadas-como-deposito-generico.md

tags:
  - cointracking
  - staking
  - rewards
  - behavioral

notes: "Operativo: cómo CoinTracking registra y classifica staking vs rewards."
---

# Cómo CoinTracking Maneja Staking y Rewards

## Diferencia Crítica en CoinTracking

**Staking** y **Rewards** se comportan diferente en CoinTracking dependiendo de cómo se importen:

### Staking (depósito de fondos)

```
Depósito en Binance Earn:
  5 ETH → Flexible Earn → CoinTracking registra como "Deposit"
  
Resultado:
  - No aparece como "venta" (no consume base de coste)
  - El ETH sigue en balance (aunque esté bloqueado en earn)
  - Cuando se retira: aparece como "Withdrawal"
```

### Rewards (ingresos generados)

```
Después de N días:
  +0.1 ETH recompensa → CoinTracking registra como "Income" o "Deposit"
  
Resultado:
  - Puede aparecer como "Deposita" (correcto) o como "Income" (confuso)
  - Afecta al cálculo de ganancias si se clasifica mal
```

---

## Problema Común (CT-005 patrón)

**Síntoma:** Recompensa de staking clasificada como "Deposit (generic)" en lugar de Income.

**Causa:** Importación vía CSV o API con tipo incorrecto.

**Impacto fiscal:** Si se trata como "Deposit" (depósito fiat), distorsiona la base del ahorro. Si es "Income" (rendimiento), es la clasificación correcta para IRPF.

---

## Cómo Diferenciar en CoinTracking

### En la importación (mejor momento)

```
Binance → Earn → Flexible Earn:
  - Depósito inicial: Tipo = "Deposit" (correcto)
  - Recompensa: Tipo = "Income" o "Reward" (correcto)
  
Si importas vía CSV y todo aparece como "Deposit":
  → Cambiar las recompensas a "Income" manualmente
```

### Identificar recompensas

```
Busca en CoinTracking:
  - Cantidad pequeña (0.001 a 0.1 del activo)
  - Fecha: después del depósito inicial
  - Tipo: puede aparecer como "Deposit" pero el nombre sugiere "reward"
  - Descripción: si dice "staking reward", "earn interest", es recompensa
```

---

## Tratamiento Fiscal (España, IRPF)

**Recompensas de staking = Rendimiento del capital (Tipo: Income)**

```
Regla (DGT):
  - Momento exigible: cuando se acredita la recompensa
  - Valuación: precio del activo en ese momento
  - Impacto: se suma a "Rendimientos del capital" en IRPF
  - No es ganancia patrimonial (eso es venta posterior)
```

**Ejemplo:**
```
1 junio: Depósito 5 ETH en Binance Earn
1 julio: Recompensa +0.1 ETH (precio: 1500€)

CoinTracking debe registrar:
  - 1 junio: Deposit 5 ETH
  - 1 julio: Income 0.1 ETH (valuado 150€)

IRPF 2025:
  - 150€ va a "Rendimientos del capital"
  - Si después vendes los 5.1 ETH: ganancia patrimonial por la diferencia
```

---

## Validación en CoinTracking

```
Reports → Gains:
  ¿Aparecen "rendimientos del capital" en la sección correspondiente?
    SÍ → OK
    NO → Verificar que staking rewards estén clasificados como "Income"
```

---

## Caso Referencia: CT-005

**Síntoma:** Recompensa de staking registrada como "Deposit (generic)"

**Problema:** CoinTracking la suma al pool de depósitos, no a rendimientos.

**Solución:**
1. Editar operación en CoinTracking
2. Cambiar Tipo: "Deposit (generic)" → "Income"
3. Regenerar Tax Report

**Resultado:** La recompensa aparece correctamente en IRPF como rendimiento del capital.

---

## Integración

- **ADR-003:** Modelo de transacciones — Staking deposit ≠ Income reward
- **CAPITAL_INCOME.md:** Tratamiento fiscal de los rendimientos
