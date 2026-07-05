---
id: KB-D2-003
title: "Árbol de Decisión: Flujo de Auditoría (6 fases)"
level: D
domain: cointracking
source: "PROCEDURE_AUDIT_ACCOUNT + ADR-017"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-017

tags:
  - decision-tree
  - audit-flow
  - main-process

notes: "Máquina de estados del flujo de auditoría completa (6 fases)."
---




# Árbol de Decisión: Flujo de Auditoría

```
┌──────────────────────────────┐
│ START: Auditar Cuenta        │
│                              │
│ [Proyecto] [Ejercicio] [Año] │
└──────────────┬───────────────┘
               │
               ▼
    ┌─────────────────────────────┐
    │ FASE 1: Cobertura de Fuentes │
    │ ✓ ¿Todos los exchanges?      │
    │ ✓ ¿Rango completo?           │
    │ ✓ ¿Total = exchange real?    │
    └──────┬──────────────────┬───┘
           │ OK               │ PROBLEMAS
           │                  ▼
           │             ┌───────────┐
           │             │ IMPORTAR  │
           │             │ FALTANTE  │
           │             └─────┬─────┘
           │                   │
           ▼                   ▼
    ┌─────────────────────────────┐
    │ FASE 2: Duplicados          │
    │ ✓ Detectar sospechosos      │
    │ ✓ Verificar Trade ID        │
    │ ✓ Confirmar con usuario     │
    └──────┬──────────────────┬───┘
           │ OK               │ PROBLEMAS
           │                  ▼
           │             ┌───────────┐
           │             │ ELIMINAR  │
           │             │ + DOC     │
           │             └─────┬─────┘
           │                   │
           ▼                   ▼
    ┌─────────────────────────────┐
    │ FASE 3: Transferencias      │
    │ ✓ Emparejar W/D             │
    │ ✓ Verificar blockchain      │
    │ ✓ Crear manual si falta     │
    └──────┬──────────────────┬───┘
           │ OK               │ PROBLEMAS
           │                  ▼
           │             ┌───────────┐
           │             │ EMPAREJAR │
           │             │ + DOC     │
           │             └─────┬─────┘
           │                   │
           ▼                   ▼
    ┌─────────────────────────────┐
    │ FASE 4: Tipos y Base Coste  │
    │ ✓ Verificar clasificación   │
    │ ✓ Cheque cost basis         │
    │ ✓ Corregir si es necesario  │
    └──────┬──────────────────┬───┘
           │ OK               │ PROBLEMAS
           │                  ▼
           │             ┌───────────┐
           │             │ RECLASIF. │
           │             │ + RECALC  │
           │             └─────┬─────┘
           │                   │
           ▼                   ▼
    ┌─────────────────────────────┐
    │ FASE 5: Purchase Pool       │
    │ ✓ ¿Warnings de pool?        │
    │ ✓ ¿Suficientes compras?     │
    │ ✓ Resolver si falta         │
    └──────┬──────────────────┬───┘
           │ OK               │ PROBLEMAS
           │                  ▼
           │             ┌───────────┐
           │             │ RESOLVER  │
           │             │ MISSING PH│
           │             └─────┬─────┘
           │                   │
           ▼                   ▼
    ┌─────────────────────────────┐
    │ FASE 6: Cierre              │
    │ ✓ Verificar integridad      │
    │ ✓ Documentar cambios        │
    │ ✓ Guardar Tax Report        │
    └──────┬──────────────────┬───┘
           │ OK               │ PROBLEMAS
           │                  ▼
           │             ┌───────────┐
           │             │ REVISAR   │
           │             │ NUEVAMENTE│
           │             └─────┬─────┘
           │                   │
           └─────┬─────────────┘
                 │
                 ▼
          ┌──────────────┐
          │ AUDITORÍA OK │
          │ REPORTABLE   │
          └──────────────┘
```

---

## Regla Clave (ADR-017)

**Orden fijo de fases.** NO saltar pasos. Cada fase reduce falsos positivos del siguiente.

Si hay problema en fase N → resuelve en fase N, recalcula, continúa → fase N+1.
