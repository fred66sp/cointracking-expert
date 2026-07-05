---
id: KB-C2-003
title: "Patrón: Emparejamiento de Transferencias (Withdrawal ↔ Deposit)"
level: C
domain: cointracking
source: "Casos CT-001, CT-004, CT-015 + ADR-004"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-004
  - ADR-005

related_docs:
  - CT-001-transferencia-entre-exchanges-importada-solo-en-origen.md
  - CT-015-swap-defi-fragmentado-en-varias-operaciones-on-chain.md

tags:
  - pattern
  - transfers
  - matching
  - cointracking

notes: "Heurísticas para emparejar withdrawal/deposit entre exchanges y wallets."
---




# Patrón: Emparejamiento de Transferencias

## Definición

Una transferencia legítima = Withdrawal (origen) + Deposit (destino) + blockchain confirmation.

Una transferencia **huérfana** = Falta uno de los tres.

---

## Heurísticas de Matching

### Búsqueda por fecha + importe

```
Withdrawal: BTC 1.5 el 2024-03-15 10:30 UTC (Binance)
Buscar Deposit: BTC 1.5 en rango 2024-03-15 10:30 a 2024-03-15 13:00 UTC (Kraken)

¿Coincide fecha ± 2h e importe exacto?
  SÍ → Probablemente la misma transferencia
  NO → Posible huérfana o comisiones de red
```

### Búsqueda con comisiones de red

```
Si no coincide importe exacto:
  Withdrawal: 1.5 BTC
  Deposit: 1.495 BTC (- 0.005 comisión de red)

¿Deposit ≈ Withdrawal - fee_estimado?
  SÍ → Legítima (comisión de red deducida)
  NO → Seguir investigando
```

### Búsqueda por blockchain

```
Si CoinTracking no tiene fecha exacta:
  Usar Tx Hash (en CoinTracking si está disponible)
  Buscar en Blockchain explorer (etherscan, btcscan, etc)
  
¿Confirma fecha de la transacción on-chain?
  SÍ → Usar fecha blockchain como truth
  NO → Verificar Tx Hash (podría ser incorrecto)
```

---

## Tolerancias

| Parámetro | Tolerancia | Razón |
|-----------|-----------|-------|
| Fecha | ±2 horas | Zona horaria, blockchain confirmation delay |
| Importe | Exact (salvo comisión) | Comisiones de red, slippage en bridges |
| Exchange destino | Flexible | Podría llegar a wallet externa |

---

## Caso de Referencia: CT-001

**Síntoma:** Balance negativo en Kraken.

**Investigación:**
- Withdrawal BTC 2.0 desde Binance el 2024-03-01 09:00 UTC
- NO existe Deposit en Kraken

**Solución:**
- Buscar en Kraken historial: encontrar Deposit 2.0 BTC el 2024-03-01 11:30 UTC
- Emparejar manualmente (o importar CSV faltante del destino)

---

## Checklist de Matching

- [ ] ¿Existe Withdrawal en el origen?
- [ ] ¿Existe Deposit en el destino?
- [ ] ¿Fechas coinciden (±2h)?
- [ ] ¿Importes coinciden o explica la diferencia?
- [ ] ¿Hay Tx Hash que confirma la transferencia?
- [ ] ¿Si es bridge/swap, hay registros en DeFi?

---

## Integración con ADRs

- **ADR-004:** Verificar contra exchange real
- **ADR-005:** Zona horaria ±2h en búsqueda
