---
id: KB-D2-002
title: "Árbol de Decisión: Saldo Negativo"
level: D
domain: cointracking
source: "PATTERN_BALANCE_RECONCILIATION"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

tags:
  - decision-tree
  - negative-balance

notes: "Máquina de estados para diagnosticar y resolver saldos negativos."
---




# Árbol de Decisión: Saldo Negativo

```
┌──────────────────┐
│ Balance < 0      │
│ en [Activo]      │
└────────┬─────────┘
         │
         ▼
    ┌────────────────────────┐
    │ ¿Están TODOS los       │
    │ exchanges importados?   │
    └────┬─────────┬─────────┘
         │ NO      │ SÍ
         ▼         ▼
    ┌────────┐  ┌──────────────────┐
    │IMPORTAR│  │ ¿Cubre rango     │
    │FALTANTE│  │ desde primer      │
    └────────┘  │ depósito?        │
         │      └────┬─────────┬──┘
         │           │ NO      │ SÍ
         │           ▼         ▼
         │      ┌─────────┐ ┌───────────┐
         │      │REIMPORTAR│ ¿Hay      │
         │      │HISTÓRICO │ compra?   │
         │      └─────────┘ └────┬──┬─┘
         │           │           │  │
         │           └──────┬────┘  │ NO
         │                  │       ▼
         │                  ▼   ┌──────────┐
         │              ┌─────┐ │MISSING PH│
         │              │OK   │ │IMPORTAR  │
         │              └─────┘ │HISTÓRICO │
         │                      └──────────┘
         │                          │
         └──────────────┬───────────┘
                        │
                        ▼
                   ┌────────────┐
                   │ ¿Zona hor. │
                   │ correcta?  │
                   └────┬─┬────┘
                        │ │
                    NO ▼   ▼ SÍ
                   ┌────┐ ┌─────┐
                   │FIX │ │OK   │
                   └────┘ └─────┘
                        │
                        ▼
                   ┌────────────┐
                   │ RESUELTO   │
                   └────────────┘
```

---

## Diagnóstico

**Orden fijo:**
1. Cobertura (¿faltan exchanges?)
2. Rango (¿cubre desde el inicio?)
3. Compras (¿hay al menos una?)
4. Zona horaria (¿Europe/Madrid?)
5. Resuelto (balance ≥ 0)
