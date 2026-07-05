---
id: KB-C1-002
title: "BingX: Copy Trading losses no exportadas — Auditoría y tratamiento fiscal"
level: C
domain: cointracking
source: "Auditoría agp2025 (2026-07-03)"
authority: verified
last_verified: 2026-07-05
valid_until: 2027-07-05
confidence: high
version: 1.0

tags:
  - bingx
  - copy-trading
  - losses
  - non-deductible
  - documentation

---

# BingX: Copy Trading Losses — Caso Crítico

**Tipo:** Patrón operativo con implicación fiscal (Nivel C)  
**Cuando aplica:** Usuario audita BingX con Copy Trading activado  
**Caso base:** agp2025 (cierre 2026-07-03)  
**Impacto:** Pérdidas reales pero sin justificante oficial

---

## 1. El Problema

### A. BingX No Exporta Copy Trading

**Situación:**
- Usuario ejecuta Copy Trading en BingX (seguir trader, operaciones automáticas)
- Las operaciones ocurren (pérdidas reales en sub-cuenta)
- CoinTracking **no las importa automáticamente**
- Solo aparece: una fila "Lost" con importe total

**Ejemplo (agp2025):**
```
Type: Lost (No categorizado)
Amount: -694.67 USDT
Timestamp: (vacío, agregado)
Description: (vacío)
```

### B. Por Qué Ocurre

**Causa técnica:**
- BingX tiene sub-cuentas separadas (Spot, Margin, Futures, Copy Trading)
- El import de CoinTracking captura: Spot, Margin, Futures
- Copy Trading **requiere export manual separado** (BingX no lo ofrece públicamente)

**Por qué BingX lo hace:**
- Copy Trading es "blackbox" (seguir trader, sin tu decisión directa)
- Responsabilidad legal difusa (¿Es tu pérdida o del trader copiado?)
- Falta de API para exportar

---

## 2. Auditoría: Detectar y Verificar

### Paso 1: Detectar "Lost"

En CoinTracking, buscar tipo "Lost" o saldo desconocido:

```sql
SELECT * FROM trades WHERE type = 'Lost' OR type LIKE '%Lost%'
```

O ver en csv:
```
Type,Amount,Fee
Lost,-694.67,0
```

### Paso 2: Verificar Origen

**Conecta a BingX app:**
1. Ir a: Account > Sub-account
2. Ver: "Copy Trading Account" (si existe)
3. Buscar: History > Copier Results o Similar Trader Results
4. Sumar: Todas las pérdidas del período

**Pregunta clave:** ¿El monto en CoinTracking (-694.67 USDT) coincide con pérdidas reales en BingX Copy Trading?

### Paso 3: Documentación

Si es confirmado:
```markdown
- [✓] Copy Trading losses confirmadas en BingX
- [✓] Monto: 694.67 USDT
- [✓] Período: 2024-01-01 a 2026-07-03
- [✓] No exportables por BingX
- [✗] Sin justificante oficial
```

---

## 3. Tratamiento Fiscal (Crítico)

### A. ¿Se Puede Deducir?

**Respuesta oficial:** Requiere asesor, pero considera:

**Tesis pro-deducción:**
- Son pérdidas reales (dinero perdido)
- Tienen origen verificable (sub-cuenta BingX)
- Se pueden demostrar (historial BingX)

**Tesis contra-deducción:**
- Sin justificante oficial (export, confirmación)
- "Lost" es categoría ambigua (¿error? ¿liquidación? ¿pérdida?)
- AEAT podría rechazarlo por falta de documentación

### B. Clasificación Correcta

Si se deducen, ¿dónde?

**Opción 1: Ganancia patrimonial negativa (base del ahorro)**
- Se compensan contra ganancias de ese mismo año
- Resto se puede arrastrar (limitado)

**Opción 2: No deducible (hasta obtener justificante)**
- Registra como "pérdida no documentada"
- Asesor lo revisa y decide

---

## 4. Recomendación Práctica (agp2025)

### En la Auditoría

**Estado:** Marked as [NO DEDUCIBLE UNTIL VERIFIED]

```markdown
## BingX Copy Trading Losses

Detectado: -694.67 USDT
Origen: Copy Trading sub-account (no exportado)
Verificado: Sí (historial BingX consultado)
Justificante: No (BingX no provee)

ACCIÓN PENDIENTE:
  1. Contactar BingX: "Necesito export oficial de Copy Trading 2024-2026"
  2. Si BingX no lo proporciona: Documentar con screenshot del historial
  3. Asesor fiscal: Revisar deductibilidad sin justificante oficial
```

### En la Declaración

**Opción A (conservadora):**
```
Base del ahorro: X (sin incluir Copy Trading loss)
Nota: Copy Trading loss -694.67 USDT pending official documentation
```

**Opción B (agresiva, requiere asesor):**
```
Base del ahorro: X - 694.67 (incluye Copy Trading loss, pending verification)
Documentación: Available upon request (BingX account history)
```

---

## 5. Cómo Obtener Justificante

### Alternativa 1: BingX Export (Ideal)

```
BingX App → Account → Sub-accounts → Copy Trading Account
→ History → Export (si existe)
```

**Esperado:**
- Trader name
- Entry/exit price
- P&L per trade
- Total loss

---

### Alternativa 2: Manual Documentation

Si BingX no exporta:

1. **Screenshot del saldo final:** Copy Trading Account balance (negativo)
2. **Historial de trades:** Copy Trader Results (últimos 50+)
3. **Cálculo:**
   - Starting balance: [X]
   - Ending balance: [Y]
   - Loss: X - Y

3. **Email a BingX:** "Solicito confirmación oficial de pérdidas Copy Trading 2024-2026"

---

## 6. Caso Real: agp2025 (Cierre)

**Hallazgo (2026-07-03):**
```
BingX Copy Trading loss: -694.67 USDT (no exportado)
Verificación: ✓ Confirmado en app BingX
Justificante: ✗ No oficial
Decisión: MARCADO COMO [NO DEDUCIBLE] hasta obtener evidencia
```

**Documentación:**
- REGISTRO-CAMBIOS.md § "CIERRE DEFINITIVO de BingX"
- Memory: audit_state.md

**Para el usuario:**
> "Tienes una pérdida real de 694.67 USDT en Copy Trading, pero BingX no lo exporta. No la deduces en la declaración hasta obtener un justificante oficial de BingX (export o confirmación por email). Guarda tu solicitud a BingX y el historial de la app por si AEAT lo pide."

---

## 7. Checklist: Si Auditas BingX

- [ ] ¿Hay "Lost" o saldo desconocido? → Buscar origen
- [ ] ¿Cuánto es? → Registrar monto exacto
- [ ] ¿Copy Trading activo/inactivo? → Confirmar período
- [ ] ¿Coincide con historial BingX? → Verificar monto
- [ ] ¿Justificante disponible? → SI: guardar | NO: documentar para asesor
- [ ] ¿Se deduce o no? → Decisión asesor + usuario

---

**Próximo:** Si encuentras Copy Trading losses en otro exchange (Bybit, OKX), documenta y avisa para C1-003.
