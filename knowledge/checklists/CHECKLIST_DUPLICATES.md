---
id: KB-D1-001
title: "Checklist: Detección de Duplicados (antes de eliminar)"
level: D
domain: cointracking
source: "PATTERN_DUPLICATE_DETECTION + ADR-014"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-014
  - ADR-026

tags:
  - checklist
  - duplicates

notes: "Verificar cada item antes de marcar como duplicado."
---




# Checklist: Detección de Duplicados

**NUNCA eliminar sin completar TODOS los items y confirmación explícita.**

## Verificación Técnica

- [ ] ¿Misma fecha + hora (mismo segundo)?
- [ ] ¿Mismo precio unitario?
- [ ] ¿Mismo volumen (cantidad)?
- [ ] ¿Misma comisión?
- [ ] ¿Mismo tipo de operación (Buy/Sell)?

## Verificación contra Exchange Real

- [ ] ¿Accediste a Binance/Kraken/Coinbase?
- [ ] ¿Buscaste el Trade ID / Transaction ID de AMBAS operaciones?
  - Trade ID #1: `_________________`
  - Trade ID #2: `_________________`
- [ ] ¿Son distintos?
  - SÍ → **ALTO**: Son legítimas. NO ELIMINAR.
  - NO → Continuar

## Verificación de Fuente

- [ ] ¿Una vino de API y otra de CSV?
  - SÍ → Probable duplicado (reimportación)
  - NO → Continuar
- [ ] ¿Ambas de la misma fuente?
  - SÍ → Más sospechoso

## Consentimiento del Usuario

- [ ] ¿Explicaste al usuario qué lo hace sospechoso?
- [ ] ¿Mencionaste la posibilidad de que sean legítimas?
- [ ] ¿El usuario confirmó explícitamente "Sí, bórralos"?
  - SÍ → Proceder
  - NO → NO eliminar

## Documentación

- [ ] ¿Documentaste en REGISTRO-CAMBIOS.md?
  - Qué se eliminó
  - Por qué se consideró duplicado
  - Confirmación del usuario
  - Fecha y hora
