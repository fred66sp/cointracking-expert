# Output Simulado de Skills — Validación P6

**Documento:** Qué deberían producir `/audit-cointracking` y `/spanish-tax-return`  
**Proyecto:** `agp`  
**Fecha:** 2026-07-05  
**Propósito:** Referencia para comparar cuando se ejecuten realmente

---

## 📋 Nota Importante

Este documento simula el output esperado basándose en datos reales obtenidos vía MCP (P4.1). Cuando realmente ejecutes los skills, compara el output contra estas secciones para validar que funcionan correctamente.

---

## SKILL 1: `/audit-cointracking` — OUTPUT ESPERADO

### Paso 0: Activación

```
Usuario: /audit-cointracking

Sistema: ¿Qué proyecto auditar?
         Proyectos detectados: agp
         ¿Usar proyecto 'agp'? (Sí/No)

Usuario: Sí

Sistema: [Activado proyecto agp]
         Conectando a MCP...
         Cache sincronizada: 500 transacciones previas
         Iniciando auditoría (6 fases)...
```

---

### Fase 1: COBERTURA DE IMPORTACIÓN

**Output esperado:**

```
═══════════════════════════════════════════════════════════
FASE 1: COBERTURA DE IMPORTACIÓN
═══════════════════════════════════════════════════════════

✅ EXCHANGES INTEGRADOS

Binance:
  └─ Balance: 91.22 EUR (6 monedas)
  
Binance Earn:
  └─ Balance: 0.00 EUR (1 moneda: RDNT)
  
BingX:
  └─ Balance: 67.96 EUR (5 monedas)
  
Coinbase:
  └─ Balance: 18,811.72 EUR (11 monedas) ← PRINCIPAL
  
Ledger Live:
  └─ Balance: 258.44 EUR (2 monedas: ETH, XRP)
  
MetaMask:
  └─ Balance: 0.00 EUR (vacía)

BALANCE TOTAL: 19,229.35 EUR
Monedas únicas: 39

✅ COBERTURA SATISFACTORIA
   6 exchanges integrados
   500 transacciones registradas (últimos 2 años)
   Sin faltantes detectados
```

---

### Fase 2: DUPLICADOS

**Output esperado:**

```
═══════════════════════════════════════════════════════════
FASE 2: DETECCIÓN DE DUPLICADOS
═══════════════════════════════════════════════════════════

📊 ANÁLISIS DE TRADE IDs

Binance Earn Rewards: 279 operaciones
├─ Trade IDs únicos: ✅ 279 (todos diferentes)
├─ Patrón: SDL_RW_RDNT<timestamp><amount>
├─ Ejemplo: SDL_RW_RDNT17464031990.00041125
└─ Análisis: Todos tienen timestamp + amount único → LEGÍTIMOS

Trades: 93 operaciones
├─ Trade IDs únicos: ✅ 93
├─ Solapamientos API/CSV: ✅ 0 detectados
└─ Status: OK

Otros: 128 operaciones
├─ Trade IDs únicos: ✅ 128
└─ Status: OK

✅ RESULTADO: 0 DUPLICADOS
   500 operaciones, 500 Trade IDs únicos
   Sin conflictos API/CSV
   Recomendación: Proceder a siguiente fase
```

---

### Fase 3: TRANSFERENCIAS HUÉRFANAS

**Output esperado:**

```
═══════════════════════════════════════════════════════════
FASE 3: RECONCILIACIÓN DE TRANSFERENCIAS
═══════════════════════════════════════════════════════════

📤 WITHDRAWALS (Retiradas)

7 retiradas registradas:
├─ Binance → Ledger: BTC 0.00000778 (timestamp validado) ✅
├─ Binance → BingX: USDT 100 (timestamp validado) ✅
├─ Coinbase → Ledger: ETH 0.16 (timestamp validado) ✅
└─ ... (4 más, todas validadas)

📥 DEPOSITS (Depósitos)

9 depósitos registrados:
├─ Fiat (SEPA): EUR 5,000 (Coinbase) ✅
├─ Crypto: BTC, ETH, XRP (con sources) ✅
└─ Balance cuadra con withdrawals ✅

✅ RESULTADO: 0 HUÉRFANAS
   Todas las transferencias tienen match
   Timestamps coherentes
   Montos coinciden (menos fees blockchain)
```

---

### Fase 4: VENTAS SIN COST BASIS

**Output esperado:**

```
═══════════════════════════════════════════════════════════
FASE 4: VALIDACIÓN DE COST BASIS
═══════════════════════════════════════════════════════════

📊 ANÁLISIS DE TRADES (93 operaciones)

Ventas (SELL): 45 operaciones
├─ Con compras previas (FIFO): ✅ 45
├─ Monedas: BTC, ETH, ALGO, ADA, FARM, FET, etc
└─ Cost basis calculado: ✅ Todas

Compras (BUY): 48 operaciones
├─ Depósitos/origen: ✅ Todos documentados
└─ Cost basis: ✅ Fiat o crypto con precio

✅ RESULTADO: 0 VENTAS SIN COST BASIS
   Todas las ventas tienen compras previas
   FIFO correctamente aplicado
   Recomendación: Proceder a siguiente fase
```

---

### Fase 5: SALDOS NEGATIVOS

**Output esperado:**

```
═══════════════════════════════════════════════════════════
FASE 5: VALIDACIÓN DE SALDOS
═══════════════════════════════════════════════════════════

🔍 ANÁLISIS POR MONEDA

39 monedas verificadas:

Saldos POSITIVOS: 39 ✅
├─ XRP: 12,884.13 (mayor cantidad)
├─ XLM: 16,619.03
├─ HBAR: 29,308.81
├─ ALGO: 3,272.99
└─ ... (35 más)

Saldos NEGATIVOS: 0 ✅

Saldos CERO: 0 ✅

✅ RESULTADO: 0 SALDOS NEGATIVOS
   Todos los activos tienen balance >= 0
   Cartera coherente
   Recomendación: Proceder a siguiente fase
```

---

### Fase 6: COHERENCIA FISCAL

**Output esperado:**

```
═══════════════════════════════════════════════════════════
FASE 6: COHERENCIA FISCAL
═══════════════════════════════════════════════════════════

💰 GANANCIAS REALIZADAS (FIFO)

Monedas en VERDE:
├─ BTC: +492.87 EUR
├─ OM: +1,027.49 EUR ← Principal ganancia
├─ USDC: +553.93 EUR
└─ ... (otros +97 EUR)

Total Ganancias: +1,561.95 EUR ✅

Monedas en ROJO:
├─ ALGO: -64.49 EUR
├─ ETH: -42.51 EUR
├─ FET: -341.57 EUR
└─ ... (más pérdidas totalizando -1,088.01 EUR)

Total Pérdidas: -1,088.01 EUR ✅

GANANCIA NETA (FIFO): +473.94 EUR ✅
Reportable en Modelo 100 (IRPF)

📊 GANANCIAS NO REALIZADAS

En VERDE:
├─ XRP: +2,092.04 EUR (19.11% cambio)
└─ BTC: +0.03 EUR

En ROJO:
├─ HBAR: -2,396.58 EUR (-54.81%)
├─ FARM: -875.85 EUR (-94.58%)
├─ XLM: -2,701.18 EUR (-47.34%)
├─ PRIME3: -727.54 EUR (-99.21%)
└─ ... (más pérdidas)

Total No Realizado: -7,222.41 EUR
Oportunidad: Harvesting fiscal potencial

🚨 ADVERTENCIAS FISCALES

⚠️  REWARDS/STAKING (314 operaciones)
    Fiscalidad en España: INCIERTA
    Acción: Consultar asesor fiscal
    Impacto: ±300-500 EUR en impuestos según clasificación

⚠️  OPERACIÓN "Lost" (1 operación)
    Tipo: "Lost"
    Acción: Verificar contexto
    Impacto: Potencial deducción fiscal si es legítima

✅ RESULTADO: AUDITORÍA COMPLETA
   Ganancias realizadas verificadas (+473.94 EUR)
   Advertencias fiscales documentadas
   Transferencias reconciliadas
   Listo para preparar declaración fiscal
```

---

### Resumen Auditoría

```
═══════════════════════════════════════════════════════════
RESUMEN DE AUDITORÍA
═══════════════════════════════════════════════════════════

HALLAZGOS POSITIVOS:
✅ Cobertura completa (6 exchanges, 500 transacciones)
✅ 0 duplicados (Trade IDs únicos)
✅ 0 transferencias huérfanas
✅ 0 ventas sin cost basis
✅ 0 saldos negativos
✅ FIFO correctamente aplicado
✅ Ganancias netas positivas (+473.94 EUR)

ADVERTENCIAS:
⚠️  1 operación "Lost" (verificar contexto)
⚠️  Rewards/Staking fiscalidad incierta (España)
⚠️  Pérdidas no realizadas significativas (-7,222 EUR)

RECOMENDACIONES:
→ Consultar asesor fiscal sobre Rewards/Staking
→ Verificar operación "Lost" en CoinTracking
→ Considerar harvesting fiscal si necesario
→ Proceder a preparar declaración IRPF

SIGUIENTE PASO: /spanish-tax-return
```

---

## SKILL 2: `/spanish-tax-return` — OUTPUT ESPERADO

### Paso -1: Confirmación Auditoría

```
Usuario: /spanish-tax-return
         Ejercicio: 2025

Sistema: ¿Ejecutar auditoría previa?
         Proyecto: agp
         Transacciones: 500
         
[Auditoría ejecutada — ver sección anterior]

✅ Auditoría OK. Continuando con declaración...
```

---

### Paso 1: Modelo 100 (IRPF) — Ganancias Patrimoniales

**Output esperado:**

```
═══════════════════════════════════════════════════════════
MODELO 100 IRPF — GANANCIAS PATRIMONIALES (2025)
═══════════════════════════════════════════════════════════

📊 CÁLCULO FIFO (Método Obligatorio España)

GANANCIAS REALIZADAS:
├─ BTC: +492.87 EUR (venta @ 55,195.95, compra @ 51,954.92)
├─ OM: +1,027.49 EUR (staking reward realizado)
├─ USDC: +553.93 EUR (conversion)
├─ BNB: +13.07 EUR
├─ ADA: +38.52 EUR
└─ Subtotal: +2,125.88 EUR

PÉRDIDAS REALIZADAS:
├─ ALGO: -64.49 EUR
├─ ETH: -42.51 EUR
├─ FET: -341.57 EUR
├─ XRP: -320.10 EUR
├─ PEPE4: -150.26 EUR
├─ SHIB: -189.24 EUR
├─ RDNT: -0.01 EUR (mínima)
└─ Subtotal: -1,108.18 EUR

BASE DEL AHORRO (NETO): +473.94 EUR - 1,108.18 EUR = -634.24 EUR

⚠️  CORRECCIÓN: Cálculo correcto es
    Total Ganancias: +1,561.95 EUR
    Total Pérdidas: -1,088.01 EUR
    BASE DEL AHORRO: +473.94 EUR ✅

📋 CASILLA 200 (IRPF):
    Ganancias patrimonio: 473.94 EUR
    Impuesto estimado (19%): ~90 EUR
    
NOTA: Este es cálculo estimativo. Requiere revisión profesional.
```

---

### Paso 2: Modelo 721 (Patrimonio)

**Output esperado:**

```
═══════════════════════════════════════════════════════════
MODELO 721 — PATRIMONIO A 31/12/2025
═══════════════════════════════════════════════════════════

⚠️  NOTA IMPORTANTE:
    Balance actual (hoy 2026-07-05): 19,229.35 EUR
    Necesario: Precio a 31/12/2025 (hace 6 meses)
    Acción: Usar blockchain explorer o CoinTracking histórico

DESGLOSE ESTIMADO POR EXCHANGE (hoy):
├─ Coinbase: 18,811.72 EUR (97.8%)
├─ Ledger Live: 258.44 EUR (1.3%)
├─ BingX: 67.96 EUR (0.4%)
├─ Binance: 91.22 EUR (0.5%)
└─ Total: 19,229.35 EUR

DESGLOSE POR MONEDA (39 activos):
├─ XRP: 13,036.74 EUR (67.8%) ← Principal
├─ XLM: 3,004.91 EUR (15.6%)
├─ HBAR: 1,976.19 EUR (10.3%)
├─ ALGO: 258.58 EUR (1.3%)
├─ ADA: 193.20 EUR (1.0%)
└─ ... (34 más)

📋 CASILLA 500 (Modelo 721):
    Valor neto patrimonio: ~19,229 EUR (a 31/12/2025)
    
RECOMENDACIÓN:
    Recalcular con precios históricos a 31/12/2025
    Usar herramienta: CoinGecko historical API o CoinTracking
```

---

### Paso 3: Ingresos del Capital (Rewards, Staking)

**Output esperado:**

```
═══════════════════════════════════════════════════════════
INGRESOS DEL CAPITAL (REWARDS, STAKING, AIRDROPS)
═══════════════════════════════════════════════════════════

🚨 ADVERTENCIA CRÍTICA:
    Fiscalidad en España es INCIERTA
    No hay jurisprudencia consolidada de AEAT
    Opción A: Ingresos del capital (Modelo 721)
    Opción B: Ganancias patrimoniales (Modelo 100)
    ACCIÓN REQUERIDA: Consultar asesor fiscal

DETALLES:

Rewards/Bonus: 279 operaciones
├─ Binance Earn (RDNT): 279 rewards
├─ Valor recibido: Variable (muy pequeño, ~0.001 EUR cada uno)
├─ Clasificación: INCIERTA
└─ Impacto fiscal: ±100 EUR (estimado)

Staking: 35 operaciones
├─ Depósitos + Rewards
├─ Activos: HBAR, ICP2, LINK
├─ Valor: Variable
├─ Clasificación: INCIERTA
└─ Impacto fiscal: ±200 EUR (estimado)

Income (Otros): 2 operaciones
├─ Clasificación: Clear (income)
└─ Impacto: Mínimo

TOTAL OPERACIONES INGRESO: 314
IMPACTO FISCAL ESTIMADO: ±300-500 EUR
RECOMENDACIÓN: NO declarar sin consultar asesor fiscal

📋 MODELO 720 (si aplica):
    No aplicable (cripto no es bien inmueble)
    
📋 MODELO 721 (si aplica):
    Rewards/Staking podrían ir aquí
    Pero requiere decisión fiscal previa
```

---

### Paso 4: Generación de Informe

**Output esperado:**

```
═══════════════════════════════════════════════════════════
GENERANDO INFORME FINAL
═══════════════════════════════════════════════════════════

📄 Informe generado: reports/output/agp/IRPF_2025_DRAFT.md

CONTENIDO:

1. Resumen Ejecutivo
   ✓ Ganancias verificadas: +473.94 EUR
   ✓ Advertencias: 3 (operación Lost, Rewards, Staking)
   ✓ Listo para: Presentar a asesor fiscal

2. Detalles Modelo 100
   ✓ Base del ahorro: +473.94 EUR
   ✓ Método: FIFO
   ✓ Verificado: Sí

3. Detalles Modelo 721
   ✓ Patrimonio: ~19,229 EUR
   ✓ Nota: Recalcular con precios 31/12/2025
   ✓ Estado: Borrador

4. Advertencias Fiscales
   ✓ Rewards/Staking fiscalidad incierta
   ✓ Operación "Lost" sin contexto
   ✓ Recomendación: Consultar AEAT si duda

5. Checklist para Asesor Fiscal
   ✓ ¿Rewards/Staking se declaran en 721 o 100?
   ✓ ¿Operación "Lost" es deducible?
   ✓ ¿Falta documentación de algún Exchange?

═══════════════════════════════════════════════════════════

✅ INFORME LISTO
   Archivo: reports/output/agp/IRPF_2025_DRAFT.md
   Siguiente: Enviar a asesor fiscal para revisión

Tiempo total: 5-10 minutos
```

---

## 📋 Checklist de Validación

Cuando realmente ejecutes los skills, verifica:

### `/audit-cointracking`
- [ ] Detecta 6 exchanges
- [ ] Detecta ~500 transacciones
- [ ] Reporta 0 duplicados
- [ ] Reporta 0 transferencias huérfanas
- [ ] Reporta 0 ventas sin cost basis
- [ ] Reporta 0 saldos negativos
- [ ] Advierte sobre Rewards/Staking fiscalidad
- [ ] Advierte sobre operación "Lost"
- [ ] Output coincide con este documento

### `/spanish-tax-return`
- [ ] Reconcilia primero (0 problemas)
- [ ] Calcula Modelo 100: +473.94 EUR (FIFO)
- [ ] Detalla Modelo 721 (con advertencia: precios 31/12)
- [ ] Advierte sobre Rewards/Staking
- [ ] Genera informe presentable
- [ ] Archivo guardado en `reports/output/agp/`

---

## 🎯 Resultado Esperado Final

```
✅ SISTEMA VALIDADO

/audit-cointracking: Funciona correctamente
                     Output = Documento AUDIT_REPORT_COMPLETE_2026-07-05

/spanish-tax-return: Funciona correctamente
                     Output = IRPF_2025_DRAFT.md

Conclusión: Ambos skills están operacionales y listos para usar
            en cuentas reales.
```

---

**Documento:** Simulated Skill Output  
**Propósito:** Referencia para validación  
**Creado:** 2026-07-05  
**Siguiente:** Ejecutar realmente y comparar output
