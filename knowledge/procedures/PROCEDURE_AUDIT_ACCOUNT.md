---
id: KB-C3-001
title: "Procedimiento: Auditar una Cuenta Completa (6 fases)"
level: C
domain: cointracking
source: "Playbook operativo + ADR-017"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-017
  - ADR-004
  - ADR-009

related_docs:
  - knowledge/patterns/PATTERN_BALANCE_RECONCILIATION.md
  - knowledge/patterns/PATTERN_DUPLICATE_DETECTION.md
  - knowledge/checklists/

tags:
  - procedure
  - audit
  - reconciliation
  - step-by-step

notes: "Orden fijo de 6 fases. No saltar pasos. Cada paso reduce falsos positivos del siguiente."
---

# Procedimiento: Auditar una Cuenta Completa

## Resumen

**6 fases en orden fijo** para auditar un historial completo de CoinTracking contra datos reales.

Duración: 1-4 horas según complejidad (número de exchanges, años, operaciones).

---

## Fase 1: Cobertura de Fuentes (20% del tiempo)

**Objetivo:** Verificar que CoinTracking tiene TODOS los datos del usuario.

### Pasos

1. **Listar exchanges en CoinTracking:** Accounts → Accounts & Wallets
2. **Listar exchanges reales del usuario:** Preguntar "¿Dónde tienes cripto?"
3. **Comparar:** ¿Faltan exchanges en CoinTracking?
   - SÍ → Importar los faltantes (CSV o API)
   - NO → Continuar
4. **Verificar rango temporal:** ¿CoinTracking cubre desde el primer depósito del usuario?
   - NO → Reimportar desde fecha anterior
   - SÍ → Fase 2

**Herramientas:** MCP `get_balance()` para totales, CSV export para detalles.

---

## Fase 2: Duplicados (25% del tiempo)

**Objetivo:** Detectar y validar duplicados.

### Pasos

1. **Ejecutar `ct_audit.py` con flag `--check duplicates`**
   - Salida: lista de transacciones sospechosas
2. **Para cada duplicado sospechoso:**
   - Verificar Trade ID en el exchange real (ADR-014)
   - Trade ID distinto → Legítimas, no eliminar
   - Trade ID igual → Posible duplicado
3. **Verificar fuente de importación:**
   - ¿Una de API y otra de CSV? → Duplicado real
   - Ambas de la misma fuente → Investigar más
4. **Pedir confirmación explícita antes de eliminar** (ADR-026: Categoría B)

**Herramientas:** Binance/Kraken web → Historial de transacciones → Buscar Trade ID

---

## Fase 3: Transferencias (20% del tiempo)

**Objetivo:** Emparejar withdrawal ↔ deposit.

### Pasos

1. **Ejecutar `ct_audit.py` con flag `--check transfers`**
   - Salida: lista de transferencias huérfanas
2. **Para cada transferencia huérfana:**
   - Buscar en CoinTracking el otro lado (withdrawal ↔ deposit)
   - Fecha ±2 horas, importe exacto
   - Si no existe: importar del exchange destino
3. **Validar contra blockchain (si es on-chain):**
   - Tx Hash en CoinTracking ↔ Blockchain explorer
   - Confirma fecha y dirección
4. **Si sigue huérfana:** Documentar en REGISTRO-CAMBIOS

**Herramientas:** Blockchain explorer (etherscan, btcscan), CSV del exchange destino

---

## Fase 4: Tipos y Base de Coste (15% del tiempo)

**Objetivo:** Verificar clasificación de operaciones.

### Pasos

1. **Ejecutar `ct_audit.py` con flag `--check types`**
   - Salida: operaciones sin base de coste (Cost = 0)
2. **Para cada operación problemática:**
   - ¿Es venta sin compra previa? → Importar historial
   - ¿Es Reward/Airdrop mal clasificado? → Cambiar tipo a Deposit
   - ¿Es Fee en tercera moneda mal tratado? → Verificar cálculo
3. **Recalcular ganancias tras cambios**

**Herramientas:** CoinTracking → Editar operación → Cambiar tipo

---

## Fase 5: Purchase Pool (10% del tiempo)

**Objetivo:** Verificar que hay compras suficientes para todas las ventas.

### Pasos

1. **Revisar warnings de CoinTracking:**
   - "All purchasing pools consumed" → Missing historial
   - "No hay una compra adecuada" → Operación específica
2. **Para cada warning:**
   - Diagnosticar causa (historial incompleto, zona horaria, tipo erróneo)
   - Resolver según patrón PATTERN_PURCHASE_POOL_EXHAUSTION
3. **Generar Tax Report (España, FIFO):**
   - CoinTracking → Reports → Tax Report → País: Spain → Método: FIFO
   - Verificar: cero warnings, ganancias ≠ 0

**Herramientas:** CoinTracking → Tax Report

---

## Fase 6: Cierre y Documentación (10% del tiempo)

**Objetivo:** Verificar integridad final y documentar cambios.

### Pasos

1. **Verificación final:**
   - [ ] ¿Cero saldos negativos?
   - [ ] ¿Cero warnings no resueltos?
   - [ ] ¿Tax Report genera sin errores?
2. **Documentar cambios en REGISTRO-CAMBIOS.md:**
   - Qué se cambió (operación, fecha, antes/después)
   - Por qué se cambió (base de coste, duplicado, zona horaria)
   - Evidencia (Trade ID, Tx Hash, mensaje de error)
3. **Guardar Tax Report:**
   - Carpeta: `reports/output/<proyecto>/`
   - Nombre: `YYYY-MM-DD_tax_report_AAAA.csv` (donde AAAA = ejercicio fiscal)

**Herramientas:** CoinTracking → Export → Tax Report

---

## Checklist Final

- [ ] Todos los exchanges importados
- [ ] Rango temporal: desde primer depósito
- [ ] Cero duplicados sin verificación
- [ ] Cero transferencias huérfanas
- [ ] Cero tipos mal clasificados
- [ ] Cero warnings de purchase pool
- [ ] Cero saldos negativos
- [ ] Tax Report generado sin errores
- [ ] Cambios documentados en REGISTRO-CAMBIOS

---

## Tiempo Esperado

| Fase | Tiempo | Si hay problemas |
|------|--------|-----------------|
| 1. Cobertura | 10min | +10-30min |
| 2. Duplicados | 15min | +20-60min |
| 3. Transfers | 15min | +20-60min |
| 4. Tipos | 10min | +10-30min |
| 5. Pool | 10min | +10-30min |
| 6. Cierre | 10min | +5-10min |
| **TOTAL** | **70min** | **+75-220min** |

---

## Integración

- **ADR-017:** Orden fijo de fases
- **ADR-004:** Verificar contra exchange real
- **ADR-009:** Protocolo crítico en cada fase
