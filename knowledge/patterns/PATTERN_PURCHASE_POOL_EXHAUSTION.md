---
id: KB-C2-004
title: "Patrón: Agotamiento de Purchase Pool (warning catálogo)"
level: C
domain: cointracking
source: "Casos CT-002, CT-017 + COST_BASIS_AND_VALIDATION.md"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-009

related_docs:
  - CT-002-venta-sin-historial-de-compra-previo-missing-purchase-history.md
  - CT-017-coste-cero-por-compra-omitida-de-ejercicios-anteriores.md
  - knowledge/cointracking/official/COST_BASIS_AND_VALIDATION.md

tags:
  - pattern
  - purchase-pool
  - warnings
  - fifo

notes: "Warning 'All purchasing pools consumed' significa que hay más ventas que compras. Nunca ignorar."
---

# Patrón: Agotamiento de Purchase Pool

## Síntoma

```
CoinTracking warning: "All purchasing pools consumed"
```

O visualmente:
- Ganancia de 0€ en una venta (cost = 0)
- "No hay una compra adecuada para esta venta"
- Purchase pool vacío pero hay más ventas registradas

---

## Significado

**Purchase pool vacío = No hay historial de compra para las ventas.**

Causas posibles:
1. Historial incompleto (años/exchanges anteriores no importados)
2. Orden cronológico incorrecto (zona horaria hace que venta aparezca antes de compra)
3. Tipo de operación clasificado erróneamente (Reward/Airdrop debería ser Deposit)
4. Transferencia sin emparejar (depósito faltante de wallet externa)

---

## Diagnóstico Rápido

```
¿Warning "All purchasing pools consumed"?
  SÍ → ¿Existe al menos una compra previa (Buy/Deposit)?
       NO → Missing Purchase History. Paso 1.
       SÍ → ¿La fecha de la compra es ANTES que la venta?
            NO → Zona horaria. Cambiar a Europe/Madrid.
            SÍ → Historial incompleto. Importar años anteriores.
  NO → OK, continuar auditoría
```

---

## Paso 1: Verificar Completitud del Historial

```
¿El primer depósito/compra de CoinTracking coincide con el 
primer movimiento del usuario en el exchange real?
  NO → Importar desde fecha anterior
  SÍ → Paso 2
```

---

## Paso 2: Verificar Zona Horaria

```
CoinTracking: Settings → Timezone
  ¿Europe/Madrid (con DST)?
    NO → Cambiar y reimportar
    SÍ → Paso 3
```

---

## Paso 3: Verificar Tipos de Operación

```
¿Hay Airdrops/Rewards clasificados como "Expense" o "Fee"?
  SÍ → Cambiar a "Deposit" o "Buy"
  NO → Historico verdaderamente incompleto
```

---

## Regla de Oro

> **"All purchasing pools consumed" NUNCA es normal.**
>
> Siempre indica datos incompletos, nunca ignorar.

---

## Casos de Referencia

| Caso | Warning | Causa | Solución |
|------|---------|-------|----------|
| **CT-002** | All pools consumed | Venta sin compra previa | Importar histórico 2023 |
| **CT-017** | Ganancia 0€ | Compra en año anterior | Importar desde 2022 |
| **CT-020** | "Missing Purchase History" | Airdrop clasificado mal | Cambiar tipo a Deposit |

---

## Integración

- **ADR-009:** Protocolo crítico — este warning no puede ignorarse
- **knowledge/cointracking/official/COST_BASIS_AND_VALIDATION.md:** Explica cálculo del pool
