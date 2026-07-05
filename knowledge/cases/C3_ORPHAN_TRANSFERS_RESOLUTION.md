---
id: KB-C3-001
title: "Transferencias Huérfanas: Detección, diagnóstico y resolución"
level: C
domain: cointracking
source: "Auditoría estándar + AUDIT_EXCHANGE_MIGRATION.md"
authority: verified
last_verified: 2026-07-05
valid_until: 2027-07-05
confidence: high
version: 1.0

tags:
  - orphan-transfers
  - unmatched-deposits
  - unmatched-withdrawals
  - blockchain-delays
  - cross-exchange

---

# Transferencias Huérfanas: Diagnóstico y Resolución

**Tipo:** Patrón de auditoría (Nivel C — muy común)  
**Cuando aparece:** Usuario migra entre exchanges, hay delays, o importación parcial  
**Severidad:** Media (detecta inconsistencias) a Alta (saldo negativo)

---

## 1. ¿Qué es una Transferencia Huérfana?

### Definición

Retirada en un exchange **sin depósito correspondiente** en el destino, o al revés.

### Ejemplo

```
BINANCE:
  2026-07-01 12:00 — Withdraw 1 BTC to address 0x123...
  Status: Success

COINBASE:
  (no hay depósito de 1 BTC el 2026-07-01)

RESULTADO EN COINTRACKING:
  Binance: BTC -1 (retirada)
  Coinbase: BTC = sin cambios
  
  → BTC negativo en Binance (HUÉRFANA)
```

---

## 2. Causas Raíz

### A. Blockchain Delay (Más Común)

**¿Qué pasó?**
- Binance: Retirada el 2026-07-01 12:00
- Blockchain: Tarda 2-4 horas en confirmar
- Coinbase: Recibe el 2026-07-01 15:30 (3.5h después)
- Importación: CoinTracking trae datos por fecha, los ve en distinto día

**Síntoma:**
```
Binance (2026-07-01): Withdraw BTC -1
Coinbase (2026-07-02): Deposit BTC +1 (un día después)
```

**Diagnóstico:**
- ✅ Es huérfana temporal (se resuelve sola)
- ⏳ Mañana cuando reimportes, desaparece

---

### B. Importación Parcial (Muy Común)

**¿Qué pasó?**
- Usuario importó Binance 2024-2025
- Pero importó Coinbase solo 2025
- 2024 tiene retiradas de Binance → sin depósitos en Coinbase (porque no se importaron)

**Síntoma:**
```
Binance: Withdraw BTC 2024-06-01 (importado)
Coinbase: No hay depósitos 2024 (no importados)

→ BTC negativo en Binance aunque la transferencia ocurrió
```

**Diagnóstico:**
- ❌ Es huérfana real (hasta importar Coinbase 2024)
- ✅ Se resuelve importando el período faltante

---

### C. Depósito No Acreditado (Raro)

**¿Qué pasó?**
- Usuario retiró de Binance
- Dinero confirmado en blockchain (visible en Etherscan)
- Coinbase **rechazó** o perdió el depósito (bug raro)

**Síntoma:**
```
Blockchain: TX confirmada, fondos llegaron
Coinbase: No hay registro, balance no aumentó
```

**Diagnóstico:**
- ❌ Es huérfana real (dinero perdido)
- ⚠️ Raro, requiere contactar a Coinbase

---

### D. Migración Incompleta (MiCA, Binance EU)

**¿Qué pasó?**
- Usuario migraba Binance → Coinbase por MiCA
- Retirada iniciada pero usuario se fue
- Dinero en limbo (nunca llegó)

**Síntoma:**
```
Binance: Withdraw 2026-07-01 (completado)
Coinbase: Sin depósito (ni siquiera fallido)
Blockchain: TX no existe o está pendiente
```

**Diagnóstico:**
- ❌ Es huérfana real (falta acción)
- ✅ Se resuelve verificando dirección y reintentando

---

## 3. Cómo Detectar

### En CoinTracking

```
1. Ver saldo por activo/exchange
2. Buscar negativos: BTC -0.1, ETH -1.5, etc.
3. Filtrar por fecha (cuando fue la retirada)
4. Correlacionar con depósito en otro exchange
```

### Script de Auditoría

```python
# Buscar retiradas sin depósito correspondiente
for withdraw in binance_withdrawals:
    matching_deposit = find_deposit(
        asset=withdraw.asset,
        amount ≈ withdraw.amount - 0.005,  # tolerancia por comisión
        timestamp ≤ withdraw.timestamp + 24h
    )
    if not matching_deposit:
        print(f"ORPHAN: {withdraw}")
```

---

## 4. Resolución Paso a Paso

### Paso 1: Identificar

```
Retirada: BTC -1.00 de Binance (2026-07-01 12:00)
Búsqueda: ¿Hay +0.999 BTC en Coinbase?
  - 2026-07-01: NO
  - 2026-07-02: SÍ (15:30)
  
DIAGNÓSTICO: Delay de 27.5 horas
CAUSA: Blockchain confirmation time
```

### Paso 2: Nivel 1 — Tx Hash

Si tienes Tx Hash, es definitivo:

**De Binance:**
```
Wallet > Withdraw History > Click withdraw
→ Copy "Tx Hash" (ej. 0x1a2b3c...)
```

**Verifica en Etherscan/Blockchain explorer:**
```
Busca 0x1a2b3c...
→ "To Address": dirección de Coinbase (validar)
→ "Status": Success (confirma llegó)
```

**En Coinbase:**
```
Movements > Deposits > Busca dirección
→ ¿Hay depósito de 1 BTC?
```

**Resultado:**
- ✅ Tx Hash coincide → Empareja en CoinTracking
- ❌ Tx no existe → Dinero nunca se envió

---

### Paso 3: Nivel 2 — Heurística (Sin Tx Hash)

Si no tienes Tx Hash, usa heurística:

**Criterios (TODOS deben cumplirse):**
1. **Moneda:** Igual (BTC = BTC)
2. **Importe:** Depósito ≈ Retirada − comisión estándar
   ```
   Retirada: 1.00 BTC
   Comisión normal Binance: 0.0005 BTC
   Depósito esperado: ~0.9995 BTC
   
   Si es: 0.9995 BTC ✅
   Si es: 0.9900 BTC ❓ (comisión anormal)
   Si es: 0.5 BTC ❌ (no es la misma)
   ```
3. **Ventana temporal:** Depósito dentro de 24h de retirada
   ```
   Retirada: 2026-07-01 12:00
   Aceptable: 2026-07-01 13:00 a 2026-07-02 13:00
   Dudoso: 2026-07-03 (48h después, blockchain delay extremo)
   ```
4. **Cuentas:** Usuario confirmó que son sus cuentas
5. **Orden temporal:** Retirada < Depósito (en tiempo real, no importación)

**Decisión:**
- Todos ✅ → Empareja (probabilidad alta)
- Alguno ❌ → No emparejes (es distinto movimiento)

---

## 5. Emparejamiento en CoinTracking

### Opción 1: Manual (Recommendado)

Si CoinTracking no las emparejó:

```
1. Ir a retirada (Binance)
2. Click "Pair Transfer"
3. Seleccionar depósito (Coinbase) correspondiente
4. Confirmar
```

### Opción 2: Crear Manual

Si falta el depósito:

```
Tipo: Deposit
Exchange: Coinbase
Asset: BTC
Amount: 0.9995 BTC
Date: 2026-07-01 15:30 (cuando llegó)
Comment: "Manual: Binance withdrawal (Tx: 0x1a2b3c...)"
Fee: 0 (ya contada en retirada)
```

---

## 6. Caso Real: MiCA Migration (Binance → Coinbase)

### Escenario

```
Julio 2026: Usuario migra 5 BTC de Binance a Coinbase
Retirada: 2026-07-05 10:00 (completada)
Depósito Coinbase: 2026-07-05 14:00 (llegó)

AUDITORÍA:

1. Detectado: BTC -5 en Binance (2026-07-05)
2. Verificado: BTC +4.9975 en Coinbase (2026-07-05)
3. Moneda: BTC = BTC ✓
4. Importe: 5 - 0.0025 (comisión) = 4.9975 ✓
5. Ventana: 4h diferencia ✓
6. Tx Hash: Validado en Etherscan ✓

RESULTADO: Emparejada ✓
```

---

## 7. Errores Comunes

| Error | Síntoma | Causa | Solución |
|-------|---------|-------|----------|
| No importar destino | Retirada sin depósito | Archivo CSV no descargado | Descargar CSV 2024 de Coinbase |
| Confundir comisiones | Importe no coincide | Comisión olvidada en heurística | Restar 0.5-1% al comparar |
| Delay > 24h | Depósito al día siguiente | Blockchain congestionado | Validar Tx Hash, empareja aunque tarde |
| Dirección errónea | TX en explorer pero sin depósito | Retirada a dirección equivocada | Contactar exchange, fondos perdidos |

---

## 8. Checklist Completo

Cuando encuentres huérfana:

- [ ] **Identificar retirada:** Monto, fecha, asset, exchange origen
- [ ] **Buscar depósito:** Mismo asset, ventana ≤24h, destino
- [ ] **Nivel 1 (Tx Hash):**
  - [ ] Copiar Tx Hash de retirada
  - [ ] Verificar en Etherscan (llegó)
  - [ ] Buscar depósito en destino por dirección
- [ ] **Nivel 2 (Heurística, si no hay TxHash):**
  - [ ] Moneda coincide
  - [ ] Importe ≈ (retirda - comisión)
  - [ ] Ventana temporal ≤24h
  - [ ] Orden temporal consistente
- [ ] **Empareja en CoinTracking:** Manual o nuevo depósito
- [ ] **Verifica:** Saldo sin negativos, balance cuadra

---

**Próximo:** Si encuentras patrón distinto (wrapped tokens, bridge operations), crea C3-002.
