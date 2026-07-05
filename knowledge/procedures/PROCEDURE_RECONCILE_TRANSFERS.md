---
id: KB-C3-003
title: "Procedimiento: Emparejar Transferencias Entre Exchanges"
level: C
domain: cointracking
source: "Casos CT-001, CT-015"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-004

related_docs:
  - PATTERN_TRANSFER_MATCHING.md
  - CT-001-transferencia-entre-exchanges-importada-solo-en-origen.md

tags:
  - procedure
  - transfers
  - reconciliation

notes: "Cómo resolver una transferencia huérfana (falta withdrawal o deposit)."
---




# Procedimiento: Emparejar Transferencias

## Causas Raíz Comunes (diagnostica antes de actuar)

Antes de asumir que es un error de importación, identifica cuál de estas 4 causas aplica — cada una tiene una resolución distinta:

| Causa | Síntoma | Resolución |
|-------|---------|-----------|
| **1. Blockchain delay** (más común) | Depósito aparece 1 día después en el exchange destino | Se resuelve solo al reimportar; no requiere acción |
| **2. Importación parcial** | Falta un rango de fechas completo en el exchange destino | Importar el período faltante |
| **3. Depósito no acreditado** (raro) | Blockchain confirma la TX pero el exchange destino no la refleja | Contactar soporte del exchange; posible pérdida real |
| **4. Migración incompleta** | Retirada iniciada pero nunca llegó a ningún sitio (ni blockchain la confirma) | Verificar dirección y reintentar desde el exchange origen |

## Síntoma

- Balance negativo en el exchange destino
- Withdrawal en origen pero sin Deposit en destino
- Advertencia "transacción incompleta"

---

## Pasos

### Paso 1: Identificar la Transferencia Huérfana

```
CoinTracking → Transacciones
  Filtrar: Tipo = "Withdrawal"
  ¿Existe un "Deposit" correspondiente 2h después en otro exchange?
    SÍ → Probablemente ya está emparejada. OK.
    NO → Paso 2
```

### Paso 2: Buscar en el Exchange Destino

```
Exchange destino (Kraken, Coinbase, etc):
  Historial → Buscar Deposit del mismo activo ±2h de la fecha del Withdrawal
  
¿Existe?
  SÍ → Importar ese exchange en CoinTracking (si no está)
  NO → Paso 3
```

### Paso 3: Verificar Blockchain

```
Si es una transferencia on-chain (p. ej. Bitcoin):
  → CoinTracking: obtener Tx Hash
  → Blockchain explorer (btcscan, etherscan, etc): buscar Tx Hash
  → Confirmar: dirección origen, destino, cantidad
  
¿Confirmó la transferencia on-chain?
  SÍ → Paso 4 (crear Deposit manual)
  NO → Posible Tx Hash incorrecto. Investigar.
```

### Paso 4: Crear Entrada Manual (si es necesario)

```
Si el Deposit no está registrado en CoinTracking:
  
1. Crear transacción:
   - Tipo: Deposit
   - Fecha: la que muestra blockchain/exchange destino
   - Cantidad: exacta
   - Activo: el mismo
   - Intercambio: el destino
   - Comisión: si aplica (red fees)

2. Regenerar Tax Report
3. Verificar: balance ya no es negativo
```

### Paso 5: Documentar

```
REGISTRO-CAMBIOS.md:
  - Qué transferencia (origen, destino, cantidad, fecha)
  - Por qué estaba huérfana (import. parcial del destino)
  - Acción tomada (crear Deposit manual o reimportar)
  - Evidencia (Tx Hash, fecha blockchain)
```

---

## Checklist

- [ ] Identificar Withdrawal huérfana
- [ ] Buscar Deposit correspondiente ±2h
- [ ] Si no existe: verificar blockchain
- [ ] Si confirmó: crear Deposit manual en CoinTracking
- [ ] Regenerar Tax Report
- [ ] Documentar

---

## Integración

- **ADR-004:** Verificar contra exchange real y blockchain
