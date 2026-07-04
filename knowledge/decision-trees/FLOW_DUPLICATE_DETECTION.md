---
id: KB-D2-001
title: "Árbol de Decisión: ¿Es un Duplicado?"
level: D
domain: cointracking
source: "PATTERN_DUPLICATE_DETECTION + ADR-014"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-014

tags:
  - decision-tree
  - duplicates
  - flowchart

notes: "Máquina de estados para clasificar duplicados potenciales."
---

# Árbol de Decisión: ¿Es un Duplicado?

```
┌──────────────────────────────────────┐
│ Hallazgo: Dos transacciones           │
│ (misma fecha, precio, volumen)        │
└────────────┬─────────────────────────┘
             │
             ▼
    ┌────────────────────┐
    │ ¿Trade ID idéntico │
    │ en Binance/Kraken? │
    └────┬───────────┬───┘
         │ SÍ        │ NO
         ▼           ▼
    ┌──────────┐   ┌─────────────────┐
    │ POSIBLE  │   │ ¿Misma fuente?  │
    │ DUPLICADO│   │ (API+API CSV+CSV)
    └────┬─────┘   └────┬───────┬────┘
         │              │ SÍ    │ NO
         │              ▼       ▼
         │          ┌──────┐  ┌────────────┐
         │          │DUPLICADO│ PROBABLEMENTE
         │          │REAL  │  │LEGÍTIMAS
         │          └──────┘  │(Batching)
         │                    └────────────┘
         │
         ▼
    ┌──────────────────────────────┐
    │ PEDIR CONFIRMACIÓN AL USUARIO │
    │ antes de eliminar (ADR-026)   │
    └──────────────────────────────┘
         │
         ▼
    ┌──────────────────────┐
    │ ¿Usuario confirma?   │
    │ "Sí, bórralos"       │
    └────┬─────────┬───────┘
         │ SÍ      │ NO
         ▼         ▼
    ┌────────┐  ┌───────────┐
    │ELIMINAR│  │MANTENER   │
    │+ DOC   │  │NO HACER   │
    └────────┘  └───────────┘
```

---

## Interpretación

1. **Trade ID idéntico en exchange** → Probablemente duplicado PERO requiere confirmación
2. **Trade ID distinto** → NO son duplicados (legítimas, p. ej. FLOKI batching)
3. **Confirmación explícita del usuario** → Requisito antes de eliminar (ADR-026)
4. **Documentación** → Registrar en REGISTRO-CAMBIOS
