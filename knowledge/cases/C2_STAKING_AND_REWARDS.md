---
id: KB-C2-001
title: "Staking y Rewards: Auditoría y clasificación fiscal"
level: C
domain: cointracking
source: "Patrones comunes + documentación oficial"
authority: verified
last_verified: 2026-07-05
valid_until: 2027-07-05
confidence: high
version: 1.0

tags:
  - staking
  - rewards
  - income
  - tax-classification
  - eth-solana-polygon

---

# Staking y Rewards: Auditoría y Fiscalidad

**Tipo:** Patrón operativo con clasificación fiscal (Nivel C)  
**Cuando aplica:** Usuario audita Ethereum staking, Solana rewards, o similares  
**Aplicable a:** Binance Earn, Lido, Solend, cualquier protocolo

---

## 1. Tipos de Staking en CoinTracking

### A. Staking Simple (Binance Earn Flexible)

**¿Qué es?**
Bloqueas cripto, recibes recompensas periódicas (diario/semanal/mensual).

**Cómo aparece:**
```
Type: Income
Description: "Binance Earn: ETH Staking Reward"
Currency: ETH
Amount: +0.00542 ETH
Date: 2026-06-30
```

**Auditoría:**
- ✅ Income registrado correctamente
- ✅ Cantidad matches Binance app
- ⚠️ Valor EUR a fecha de recepción (no a hoy)

### B. Staking Bloqueado (Período Fijo)

**¿Qué es?**
Bloqueas cripto 30/60/90 días, APY más alto que flexible.

**Cómo aparece:**
```
Type: Locked Staking / Income
Description: "Ethereum Staking 30d: Principal + Reward"
Currency: ETH
Amount: +X.XXXX ETH
Date: [fecha de liberación]
```

**Auditoría:**
- ✅ Recompensa = (reward solo, sin principal)
- ⚠️ Fecha = cuando se libera, no cuando se inicia
- ⚠️ Principal devuelto el mismo día: ¿Se cuenta como venta? NO (es retorno de lo tuyo)

### C. DeFi Staking (Lido, Rocket Pool, etc.)

**¿Qué es?**
Protocolo descentralizado, recibes token que acumula recompensas.

**Cómo aparece:**
```
Método 1 (rewards):
Type: Income
Amount: +0.05 stETH (acumulado)

Método 2 (compuesto):
No explícito: stETH/ETH valúa diferente
```

**Auditoría:**
- ⚠️ CoinTracking a veces no lo importa
- ⚠️ Si es token derivado (stETH, xSOL): cuidado con base de coste
- ⚠️ Liquidación de staked: ¿Venta? SÍ (porque vendes el staking token)

---

## 2. Clasificación Fiscal Española (Crítica)

### Staking = Rendimiento de Capital (RCM)

**Base legal:** DGT V1766-22 (consulta DGT)

**Dónde va:**
```
IRPF → Sección II: Rendimiento de Capital Mobiliario
     → Bloque III (Internet): Criptomonedas
     → Sub-base: "Rendimientos de staking"
```

**Impuesto:**
- Tipo: 19% (estado) + CCAA (variable)
- Deducción: No (tributan brutos)
- Complementario: Integración en base imponible general

### Valor a Declarar

**Regla:** Valor EUR a **fecha de recepción**, no a fecha de valoración.

**Ejemplo:**
```
30-06-2026: Recibo 0.05 ETH de staking
Valor EUR ese día: 2000 €/ETH
Declaro: 0.05 × 2000 = 100 €

Hoy (05-07-2026): ETH = 2200 €
Valor hoy: 110 €
Declaro: 100 € (el de recepción, NO 110 €)
```

---

## 3. Casos Especiales

### A. Staking Delegado (Pools)

**¿Qué es?**
No tienes tus monedas, las tiene el pool. Recibe recompensas el pool, te transfieren la tuya.

**Auditoría:**
- ✅ Es staking (recompensas son RCM)
- ⚠️ Riesgo: Pool puede default/robar (pero fiscalmente igual)

---

### B. Liquid Staking (stETH, xSOL, etc.)

**¿Qué es?**
Recibes token que representa tu stake + acumula recompensas.

**Auditoría compleja:**
```
Escenario 1: Depositas 1 ETH, recibes 1 stETH
- ¿Es venta? NO (1:1, cambio de forma)
- ¿Afecta base coste? NO (sigues siendo ETH)

Escenario 2: stETH crece (recompensas compuestas)
- ¿Cómo se reporta? Depende del protocolo
- CoinTracking: A veces no lo importa
- Solución: Manual si aplica

Escenario 3: Vendes stETH
- ¿Es venta? SÍ (vendes el staking token)
- Base coste: valor original EUR + recompensas acumuladas
```

---

## 4. Auditoría Paso a Paso

### Checklist

- [ ] **Detectar staking:**
  - Buscar "Income" tipo staking
  - Buscar "Locked" o "Earn"
  - Buscar derivados (stETH, xSOL)

- [ ] **Verificar cantidad:**
  - Matches con app (Binance, Lido, etc.)
  - Periodos completos (no meses faltantes)

- [ ] **Valoración EUR:**
  - A fecha de recepción, NO hoy
  - Tasa de cambio verificada

- [ ] **Base de coste:**
  - Si es liquid staking: include recompensas acumuladas
  - Si es compuesto: separar principal de recompensas

- [ ] **Venta de staking:**
  - Si vendes stETH/xSOL: ¿Se cuenta como venta? SÍ
  - Base coste: principal + recompensas

---

## 5. Caso Típico: Ethereum Staking en Binance

### Auditoría Ejemplo

```
Usuario: Hizo staking de 10 ETH en Binance desde 2024-06-01
Recompensas recibidas: 0.42 ETH (a lo largo del año)

AUDITORÍA:

1. Detectado:
   - 12 operaciones "Income" (staking mensual)
   - Total: 0.42 ETH

2. Verificado:
   - Matches Binance Earn history ✓
   - Valores EUR a fecha recepción ✓
   - No hay períodos faltantes ✓

3. Valoración:
   Jun 2024: 0.035 ETH × 2500 €/ETH = 87.50 €
   Jul 2024: 0.035 ETH × 2400 €/ETH = 84 €
   ... (12 meses)
   Total RCM 2024: 1.020 € (aproximado)

4. Clasificación fiscal:
   Base: Rendimiento de Capital Mobiliario
   Tramo: 19% estatal
   Importe: 1.020 €

5. Resultado:
   Declarar en IRPF: +1.020 € (RCM)
```

---

## 6. Errores Comunes

| Error | Síntoma | Causa | Solución |
|-------|---------|-------|----------|
| Valor hoy | Declaro 110 € | Usa precio de hoy en lugar de recepción | Recalcular con tasas EUR diarias |
| Confundir con venta | Pienso que vendo ETH | Es staking, no venta | Clasificar como Income (RCM) |
| Falta de recompensas | Saldo staking < esperado | No importadas en CoinTracking | Crear manuales desde Binance |
| Liquid staking confusion | No sé cómo reportar stETH | Protocol no claro | Asesor fiscal: split principal vs recompensas |

---

## 7. Documentación para Asesor

Si tienes dudas, prepara para asesor:

```markdown
## Staking Summary

- Plataforma: Binance
- Asset: ETH
- Período: 2024-06-01 to 2026-06-30
- Principal: 10 ETH
- Recompensas totales: 0.42 ETH
- Valor EUR recompensas: 1.020 € (aprox.)
- Status: Classified as RCM (income)
- Pendiente: Validar tratamiento liquid staking
```

---

**Próximo:** Si auditas Solana staking, Polygon, o DeFi protocolo distinto, documenta y crea C2-002.
