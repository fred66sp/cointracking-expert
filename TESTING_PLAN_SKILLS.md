# Plan de Testing — Skills `/audit-cointracking` y `/spanish-tax-return`

**Documento:** Plan para validar que los skills funcionan end-to-end  
**Fecha:** 2026-07-05  
**Estado:** Listo para ejecutar

---

## 🎯 Propósito

Validar que las dos skills principales funcionan correctamente con datos reales:

1. **`/audit-cointracking`** — Auditar una cuenta (reconciliación completa)
2. **`/spanish-tax-return`** — Preparar declaración IRPF de un ejercicio

---

## 📋 Prerequisitos

✅ **Completado:**
- Sistema operacional (P0-P3)
- MCP funcional (proyecto `agp` activo)
- Base de conocimiento validada (111+ documentos)
- Auditoría manual completada (P4.1)

**Antes de ejecutar los skills:**
1. Asegúrate de que el MCP está activo: `cointracking_get_balance()`
2. Proyecto activo debe ser `agp`: `cointracking_switch_project("agp")`
3. API keys configuradas en variables de entorno

---

## 🧪 Caso de Prueba: Proyecto `agp`

Datos conocidos (de P4.1):
- **Balance total:** 19,229.35 EUR
- **Activos:** 39 monedas
- **Transacciones:** 500 registradas
- **Ganancias realizadas (FIFO):** +473.94 EUR
- **Ganancias no realizadas:** +2,092.04 EUR (XRP)
- **Pérdidas no realizadas:** -9,314.45 EUR

---

## ✅ SKILL 1: `/audit-cointracking`

### Paso 1: Ejecutar

```
/audit-cointracking
```

**Contexto a proporcionar:**
- Proyecto: `agp`
- Rango: últimos 2 años (2024-2026)
- Detalle: completo (incluye todos los chequeos)

### Paso 2: Verificar Output

El skill debe devolver (en orden):

#### 2a. Cobertura de Importación ✅
- [ ] Detecta 6 exchanges (Binance, BingX, Coinbase, Ledger, MetaMask)
- [ ] Detecta ~500 transacciones
- [ ] Balance total coincide: ~19,229 EUR

#### 2b. Duplicados ⚠️
- [ ] Detecta Trade IDs únicos (Binance Earn rewards con IDs como `SDL_RW_RDNT...`)
- [ ] NO marca falsos positivos (batch operations como Binance rewards)
- [ ] Si encuentra duplicados: explica qué hace duplicado
- [ ] Pide confirmación antes de borrar

#### 2c. Transferencias Huérfanas
- [ ] Detecta si hay withdrawals sin deposits coincidentes
- [ ] Verifica timestamp y monto (menos fees)
- [ ] Esperado: 0 huérfanas (proyecto está cerrado)

#### 2d. Ventas sin Cost Basis
- [ ] Revisa todas las ventas (Trades tipo SELL)
- [ ] Verifica que tienen compras previas (FIFO)
- [ ] Esperado: todas las ventas tienen cost basis

#### 2e. Saldos Negativos
- [ ] Verifica balance de cada moneda
- [ ] Detecta si alguna está en negativo
- [ ] Esperado: 0 negativos

#### 2f. Coherencia Fiscal Cualitativa
- [ ] Detecta Rewards (279 operaciones)
- [ ] Detecta Staking (35 operaciones)
- [ ] Marca fiscalidad incierta (Rewards/Staking en España)
- [ ] Advierte sobre operación "Lost"

### Paso 3: Verificar Recomendaciones

El skill debe recomendar:
- [ ] Verificar operación "Lost" (1 operación)
- [ ] Consultar asesor fiscal para Rewards (incierto en España)
- [ ] Opcionalmente: harvesting fiscal (monedas en rojo -99%)

### Paso 4: Comparar con P4.1

**El output de `/audit-cointracking` debe ser consistente con [AUDIT_REPORT_COMPLETE_2026-07-05.md](reports/output/agp/AUDIT_REPORT_COMPLETE_2026-07-05.md):**

- Mismos hallazgos positivos ✅
- Mismas advertencias ⚠️
- Mismas recomendaciones

---

## ✅ SKILL 2: `/spanish-tax-return`

### Paso 1: Ejecutar

```
/spanish-tax-return

Ejercicio: 2025
Proyecto: agp
Método de coste: FIFO (obligatorio en España)
```

### Paso 2: Verificar Output — Paso 0 (Reconciliación Previa)

El skill debe **reconciliar primero** (no saltarse):

- [ ] Detecta 500 transacciones
- [ ] Valida balance
- [ ] Detecta duplicados (0 esperados)
- [ ] Detecta transferencias huérfanas (0 esperadas)

### Paso 3: Verificar Output — Paso 1 (Ganancias Patrimoniales)

#### 1a. Modelo 100 (IRPF)
- [ ] Calcula ganancias con FIFO
- [ ] Esperado: +473.94 EUR neto
- [ ] Detalles: ganancias por moneda (BTC +492.87, OM +1,027.49, USDC +553.93)
- [ ] Detalles: pérdidas por moneda (ALGO -64.49, ETH -42.51, FET -341.57, etc.)

#### 1b. Cálculo Correcto
```
Ganancia BTC: +492.87 EUR (entrada con comisión incluida)
Ganancia OM: +1,027.49 EUR (staking reward)
Ganancia USDC: +553.93 EUR (conversion)
... (más monedas)
Total Ganancias: +1,561.95 EUR

Pérdida ALGO: -64.49 EUR
Pérdida ETH: -42.51 EUR
... (más monedas)
Total Pérdidas: -1,088.01 EUR

NETO: +473.94 EUR ✅
```

### Paso 4: Verificar Output — Paso 2 (Patrimonio)

#### 2a. Modelo 721 (Patrimonio)
- [ ] Balance total a 31/12/2025: necesita precios históricos (no disponibles hoy)
- [ ] Desglose por exchange (Coinbase 70%, Ledger 1.3%, BingX 0.4%, Binance 0.5%)
- [ ] Desglose por moneda (XRP principal, HBAR, XLM, ALGO, etc.)

**Nota:** El skill debe **advertir** que necesita precios a 31/12/2025, no precios de hoy (2026-07-05)

### Paso 5: Verificar Output — Paso 3 (Ingresos del Capital)

#### 3a. Rewards & Staking
- [ ] Detecta 279 Rewards de Binance Earn
- [ ] Detecta 35 Staking operations
- [ ] Total: 314 operaciones de ingreso
- [ ] Fiscalidad: **MARCA COMO INCIERTA** (España no tiene jurisprudencia consolidada)
- [ ] Advierte: "Consulte con asesor fiscal antes de declarar"

### Paso 6: Verificar Output — Paso Final (Informe)

El skill debe generar un informe con:

- [ ] **Detalles Modelo 100:** ganancias/pérdidas por moneda, FIFO, neto
- [ ] **Detalles Modelo 721:** patrimonio a 31/12/2025 (advertencia: precios históricos)
- [ ] **Advertencias Fiscales:** Rewards/Staking incierto, operación "Lost" a verificar
- [ ] **Archivos generados:** informe en `reports/output/agp/`
- [ ] **Siguiente paso:** "Presente este informe a su asesor fiscal"

---

## 📊 Checklist de Testing

### SKILL `/audit-cointracking`

- [ ] Ejecuta sin errores
- [ ] Detecta cobertura (6 exchanges, 500 transacciones)
- [ ] Detecta 0 duplicados (Trade IDs únicos)
- [ ] Detecta 0 transferencias huérfanas
- [ ] Detecta 0 saldos negativos
- [ ] Detecta 0 ventas sin cost basis
- [ ] Advierte sobre 1 operación "Lost"
- [ ] Advierte sobre Rewards/Staking (incierto fiscal)
- [ ] Output coincide con AUDIT_REPORT_COMPLETE_2026-07-05.md
- [ ] Informe generado en `reports/output/agp/`

### SKILL `/spanish-tax-return`

- [ ] Ejecuta sin errores
- [ ] Reconcilia primero (0 problemas detectados)
- [ ] Calcula Modelo 100: +473.94 EUR (FIFO)
- [ ] Detecta 279 Rewards + 35 Staking
- [ ] Marca Rewards/Staking como inciertos (España)
- [ ] Genera detalles Modelo 721 (con advertencia: precios a 31/12)
- [ ] Genera informe completo
- [ ] Informe guardado en `reports/output/agp/`
- [ ] Informe listo para presentar a asesor fiscal

---

## 🚨 Criterios de Fallo

El testing **FALLA** si:

1. ❌ El skill `/audit-cointracking` detecta duplicados falsos (Rewards de Binance)
2. ❌ El skill `/spanish-tax-return` calcula ganancia diferente a +473.94 EUR
3. ❌ El skill no marca Rewards/Staking como inciertos (España)
4. ❌ El skill no genera informe en `reports/output/agp/`
5. ❌ El informe no es presentable a un asesor fiscal

---

## ✅ Criterios de Éxito

El testing **PASA** si:

1. ✅ `/audit-cointracking` detecta todos los hallazgos (sin falsos positivos)
2. ✅ `/spanish-tax-return` calcula +473.94 EUR (FIFO)
3. ✅ Ambos skills generan informes presentables
4. ✅ Los informes son consistentes entre sí
5. ✅ Las advertencias fiscales están claras

---

## 📝 Cómo Ejecutar

### Prerequisito: Proyecto Activo

```
cointracking_switch_project("agp")
```

### Test 1: Auditoría

```
/audit-cointracking
```

Verificar output contra checklist "SKILL `/audit-cointracking`" arriba.

### Test 2: Declaración IRPF

```
/spanish-tax-return
```

Proporcionar:
- Ejercicio: 2025
- Método: FIFO

Verificar output contra checklist "SKILL `/spanish-tax-return`" arriba.

---

## 📊 Reporte de Testing

Después de ejecutar ambos skills, completar:

**Test Auditoría:**
- [ ] Pasó/Falló
- [ ] Hallazgos: _______________
- [ ] Problemas: _______________

**Test IRPF:**
- [ ] Pasó/Falló
- [ ] Modelo 100: +473.94 EUR / diferente: _______________
- [ ] Problemas: _______________

**Conclusión:**
- [ ] Sistema 100% operacional
- [ ] Necesita fixes: _______________

---

## 🚪 Siguiente

Si ambos tests pasan:
- ✅ Sistema completamente validado
- ✅ Listo para usar con cuentas reales
- ✅ Documentación completa (P0-P3-P4)

Si hay problemas:
- Documentar en este archivo
- Crear issues en `AGENT_CHANGE_REQUESTS.md`
- Priorizar fixes antes de usar en producción

---

**Documento:** Plan de Testing Skills  
**Creado:** 2026-07-05  
**Estado:** Pendiente ejecución  
**Responsable:** Usuario (ejecuta skills) + Sistema (verifica output)
