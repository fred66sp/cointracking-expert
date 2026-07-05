---
id: KB-B2-011
title: "Mecánicas de Bybit: Trading Spot y Derivados"
level: B
domain: cointracking
source: "Bybit official docs + análisis de casos reales"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-05
confidence: high
version: 1.0

related_adr:
  - ADR-003
  - ADR-010

related_docs:
  - knowledge/cointracking/behavioral/BINANCE_SPOT_MECHANICS.md
  - knowledge/cointracking/behavioral/BINGX_MECHANICS.md
  - knowledge/cointracking/behavioral/API_VS_CSV_OVERLAP.md

tags:
  - exchange
  - bybit
  - behavioral
  - spot
  - derivados

notes: "Bybit es exchange asiático similar a Binance/BingX pero con énfasis en derivados. Crecimiento rápido, buena API, integración en CoinTracking disponible."
---

# Bybit Mechanics

## Características Principales

**Bybit** es un exchange descentralizado y centralizado híbrido, con fuerte presencia en Asia y Europa, que soporta:
- **Spot Trading:** Compra/venta directa de criptomonedas
- **Margin Trading:** Trading apalancado (con comisión de interés)
- **Perpetual Futures:** Contratos perpetuos con apalancamiento (similar a Binance)
- **Linear Perpetuals:** Contratos perpetuos denominados en USD/USDC (no inversas)
- **Options:** Opciones (calls/puts) — caso avanzado

---

## Operaciones en CoinTracking

### Spot Trading

| Campo | Valor | Notas |
|-------|-------|-------|
| **Type** | Trade | Compra/venta estándar |
| **Exchange** | Bybit | Identificador único |
| **Currency In/Out** | USDT, BTC, ETH, etc | Pares comunes (más que BingX) |
| **Fee Currency** | USDT (por defecto) o BIT (Bybit Token) | Reducción de comisión con BIT |
| **Trade ID** | Presente | Único por operación |

**Ejemplo real:**
```
Date: 2024-06-15 14:23:00
Type: Trade (BUY)
Buy: 1.0 BTC @ 65000 USDT
Fee: 130 USDT (0.1%, reducible a 0.075% con BIT)
Exchange: Bybit
```

### Margin Trading

| Campo | Diferencia |
|-------|-----------|
| **Type** | Trade (igual) |
| **Grupo** | "Bybit Margin" (opcional) |
| **Fee** | Interés de préstamo (variable, típico 5-15% anual) |
| **Risk** | Liquidación automática si ratio cae |

**⚠️ Advertencia:** El trading apalancado amplifica pérdidas. La fiscalidad es tratamiento de capital (igual que Spot), pero el riesgo de liquidación es alto.

### Perpetual Futures (Linear & Inverse)

| Campo | Valor |
|-------|-------|
| **Type** | "Perpetual Futures" o "Derivatives" |
| **Funding Fee** | Pagos cada 8 horas (como Binance) |
| **Position Type** | LONG / SHORT |
| **Liquidation Risk** | Sí (pérdida total del colateral posible) |
| **Notation** | Linear (USDT) o Inverse (BTC) |

**Diferencia clave:**
- **Linear:** PnL en USDT (menos arriesgado, más común)
- **Inverse:** PnL en BTC/ETH (más volátil, requiere cuidado en coste base)

**Fiscalidad:** Igual debate que Binance (patrimonial vs rendimiento). En España, probablemente ganancia patrimonial (ADR-006 límite de determinismo).

---

## Integración con CoinTracking

### ✅ Métodos Soportados

1. **API Connection (Recomendado)**
   - CoinTracking soporta Bybit API (validado 2026-07-05)
   - Datos en vivo, actualizados automáticamente
   - Comisiones incluidas
   - Soporta Spot, Margin, Futures

2. **CSV Import**
   - Bybit permite exportar Trade History (Spot + Margin separado)
   - Formato: estándar (Date, Type, Pair, Amount, Price, Fee)
   - Limitación: Manual, requiere actualización periódica
   - Futures requieren descarga separada

### Importación vía API

**Pasos:**
1. Bybit: Account → API Management → Create API Key
2. Elegir permisos: "Query only" (lectura)
3. CoinTracking: Settings → Exchanges → Bybit → Connect
4. Pegar API Key & Secret
5. Seleccionar período inicial (últimos 6-12 meses)
6. Sincronizar

**Validación:**
- [ ] Balance en CoinTracking coincide con Bybit (refresh)
- [ ] Operaciones mostradas: count correcto
- [ ] Comisiones incluidas (USDT o BIT)
- [ ] Trade IDs únicos (sin duplicados)
- [ ] Margin interest visible (si aplica)
- [ ] Funding fees capturados (si hay futures)

---

## Casos Límite y Peculiaridades

### 1. Comisión en BIT (Bybit Token)

**Problema:** Bybit incentiva usar BIT token para reducir comisión (0.1% → 0.075%).

**Solución:** 
- CoinTracking detecta automáticamente si la comisión está en BIT
- El fee se registra en la operación correspondiente
- **Verificar:** si hay BIT holdings no explicados (son restos de comisiones), documentarlos como reducción de cost basis

**Fiscalidad:** La comisión sigue siendo gasto, sea en USDT o BIT (el ticker del fee no cambia el tratamiento).

### 2. Liquidación de Posición

Si margin/futures es liquidado por caída de precio:
- Type: "Liquidation" o "Forced Closure"
- Pérdida es total del colateral (en ese contrato)

**Fiscalidad:** Pérdida patrimonial (deducible).

**En CoinTracking:** Aparecerá como operación única con balance negativo antes/después. Verificar que la pérdida coincide con el colateral bloqueado.

### 3. Funding Fees Negativos (Ganancias)

Cada 8 horas, si la tasa de funding es negativa, **tú recibes dinero** (como ingreso):
- Pago positivo → tú pagas (es fee)
- Pago negativo → tú recibes (es ingreso)

**En CoinTracking:**
- Aparece como "Funding Fee" o "Rewards"
- Si es negativo → registrarlo como Income (no Fee)
- Si es positivo → registrarlo como Fee

**Fiscalidad:** Funding negativo = ingresos del capital (Modelo 721), si aplica.

### 4. Cross-Collateral (Múltiples Monedas como Margen)

Bybit permite usar múltiples activos como colateral para una posición:
- ETH + BTC → margin pool
- Liquidación ocurre si valor total cae

**Problema:** CoinTracking puede no capturar bien la composición del colateral.

**Solución:** 
- Si usas cross-collateral, documenta manualmente qué activos están bloqueados
- Regístraloe en el comentario de la operación
- Verifica que el saldo disponible en CoinTracking coincide con (total - bloqueado)

### 5. Conversión Rápida (Similar a Binance Convert)

Bybit ofrece "Quick Exchange" para cambiar activos sin comisión (o con comisión mínima):
- Operación: Swap BTC → USDT
- Fee: Incluido en el tipo de cambio

**En CoinTracking:** Aparece como Trade (venta BTC, compra USDT). Fiscalidad: ganancia patrimonial (igual que manual swap en DEX).

---

## Validación en CoinTracking

### Checklist: Spot vs Margin vs Futures

```
[ ] API conectado en CoinTracking
[ ] Balance actual coincide (refresh)
[ ] Últimas 10 operaciones Spot visibles
[ ] Últimas 10 operaciones Margin visibles (si aplica)
[ ] Últimas 10 operaciones Futures visibles (si aplica)
[ ] Comisiones incluidas (USDT o BIT)
[ ] Trade IDs únicos (sin duplicados)
[ ] Funding fees capturados (si hay futures)
[ ] Margin interest visible (si hay margin)
[ ] No hay solapamiento API+CSV
[ ] Saldos negativos: 0 (excepto fiat, si normal)
```

### Detección de Problemas Comunes

| Problema | Síntoma | Solución |
|----------|---------|----------|
| **API desconectado** | Balance no actualiza | Reconectar API en Settings |
| **CSV duplica** | Operaciones aparecen 2x | Eliminar CSV, usar solo API |
| **BIT fee no registrado** | Comisión falta o en USDT | Verificar que CoinTracking detecta BIT automáticamente |
| **Liquidación oculta** | Pérdida sin operación clara | Buscar "Liquidation" o balance negativo |
| **Funding fees negativos** | Ingresos sin operación | Buscar "Funding Fee" (valor negativo = ingreso) |
| **Cross-collateral confuso** | Saldo disponible ≠ esperado | Verificar colateral bloqueado manualmente |
| **Margin interest no visible** | Balance negativo solo en fiat | Revisar si es interés de préstamo (normal) |

---

## Casos de Uso Reales

### Caso 1: Spot + Futures Paralelos

**Usuario:** Hace spot trading (BTC/USDT) y abre un perpetual corto simultáneamente para hedge.

```
Spot: BUY 1.0 BTC @ 65000 USDT (2024-06-15)
Futures: SHORT 1.0 BTC perpetual @ 65500 USDT (mismo día)
```

**En CoinTracking:**
- Spot: dos operaciones (compra/venta si cierra) — registra con FIFO
- Futures: una operación (posición abierta + funding fees periódicos)

**Fiscalidad:** Spot es ganancia patrimonial (FIFO); Futures es derivado (consultar asesor).

### Caso 2: Margin Liquidación

**Usuario:** Abre margin long con 5x, pero el precio cae 25%.

```
Margin Long: BUY 5 BTC @ 65000 USDT (colateral: 1 BTC)
Liquidación: LOSS 1 BTC (colateral perdido)
```

**En CoinTracking:**
- Aparecerá como Liquidation
- Pérdida es total de colateral

**Fiscalidad:** Pérdida patrimonial (deducible).

### Caso 3: Funding Fees Positivos (Pagando)

**Usuario:** Abre perpetual long, pero la tasa de funding es positiva (paga cada 8h).

```
Position: LONG 1.0 BTC perpetual (abierto)
Funding (8h): -50 USDT (pago)
Funding (8h): -50 USDT (pago)
Funding (8h): -50 USDT (pago)
...
```

**En CoinTracking:**
- Cada pago aparece como Fee en la línea de la posición o como operación separada
- Total: -150 USDT (ejemplo) → aumenta cost basis del BTC

**Fiscalidad:** Gastos de operación (deducibles).

---

## Diferencias vs Binance

| Aspecto | Binance | Bybit |
|---------|---------|-------|
| **Spot comisión** | 0.1% (0.075% con BNB) | 0.1% (0.075% con BIT) |
| **Futures notation** | Linear + Inverse | Linear + Inverse |
| **Funding rate** | Cada 8h | Cada 8h |
| **Margin interest** | Variable (5-15%) | Variable (5-15%) |
| **Quick exchange** | Sí (Convert) | Sí (Quick Exchange) |
| **Options** | Sí (Vanilla + Perpetual) | Sí (Vanilla) |
| **API stability** | Muy alta | Alta |
| **CoinTracking support** | Completo | Completo |

---

## Referencias y Recursos

- [Bybit Official API Docs](https://bybit-exchange.github.io/) (inglés)
- [CoinTracking Bybit Integration](https://www.cointracking.info/en/api_keys.php) (en plataforma)
- [Bybit Help Center](https://www.bybit.com/en-US/help-center/) (español disponible)

---

## Notas Operativas

**Para auditoría:**
- Bybit es más grande que BingX, comparable a Binance en volumen
- Excelente soporte API y estabilidad
- Margin liquidation automática es rápida (requiere atención)
- Futures son muy populares (vigilar funding fees)

**Para fiscalidad:**
- Spot: ganancias patrimoniales (FIFO)
- Margin: igual que Spot (+ riesgo de liquidación)
- Futures: puede ser capital o rendimiento (consultar asesor)
- Funding fees positivos (pagos): gastos de operación (deducibles)
- Funding fees negativos (ingresos): ingresos del capital (Modelo 721)

---

**Documento:** Bybit Mechanics  
**Nivel:** B2-011  
**Status:** Operacional  
**Creado:** 2026-07-05
