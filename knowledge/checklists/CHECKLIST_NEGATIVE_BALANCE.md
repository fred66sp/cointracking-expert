---
id: KB-D1-002
title: "Checklist: Diagnosticar Saldo Negativo"
level: D
domain: cointracking
source: "PATTERN_BALANCE_RECONCILIATION"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-004

tags:
  - checklist
  - balances
  - negative

notes: "Paso a paso para resolver saldo negativo."
---

# Checklist: Diagnosticar Saldo Negativo

## Cobertura de Fuentes

- [ ] ¿Están importados TODOS los exchanges?
- [ ] ¿Cubre el rango desde el primer depósito?
- [ ] ¿El total de CoinTracking = total del exchange?

## Por Activo

- [ ] ¿Existe al least una compra (Buy/Deposit)?
- [ ] ¿La suma de compras >= suma de ventas?
- [ ] ¿Hay una transferencia sin emparejar?

## Zona Horaria

- [ ] ¿CoinTracking está en Europe/Madrid?
- [ ] ¿El orden cronológico es correcto?

## Resolución

- [ ] Acción tomada: `_______________________`
- [ ] Documentado en REGISTRO-CAMBIOS
- [ ] ¿Saldo ahora positivo o 0?
