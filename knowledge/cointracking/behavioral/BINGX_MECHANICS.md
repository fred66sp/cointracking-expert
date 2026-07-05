---
id: KB-B2-010
title: "Mecánicas de BingX: Trading Spot y Derivados"
level: B
domain: cointracking
source: "BingX official docs + análisis de casos reales + auditoría agp2025 (2026-07-03)"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-05
confidence: high
version: 1.1

related_adr:
  - ADR-003
  - ADR-010

related_docs:
  - knowledge/cointracking/behavioral/BINANCE_SPOT_MECHANICS.md
  - knowledge/cointracking/behavioral/API_VS_CSV_OVERLAP.md

tags:
  - exchange
  - bingx
  - behavioral
  - spot
  - derivados

notes: "BingX es exchange asiático similar a Binance pero con características propias. Soporta Spot, Margin y Derivados. Integración en CoinTracking es API-based."
---




# BingX Mechanics

## Características Principales

**BingX** es un exchange descentralizado y centralizado híbrido, principal en Asia, que soporta:
- **Spot Trading:** Compra/venta directa de criptomonedas
- **Margin Trading:** Trading apalancado
- **Perpetual Futures:** Contratos perpetuos con apalancamiento
- **Copy Trading:** Réplica de estrategias de otros traders

---

## Operaciones en CoinTracking

### Spot Trading

| Campo | Valor | Notas |
|-------|-------|-------|
| **Type** | Trade | Compra/venta estándar |
| **Exchange** | BingX | Identificador único |
| **Currency In/Out** | USDT, BTC, ETH, etc | Pares comunes |
| **Fee Currency** | USDT (por defecto) | Puede ser otra |
| **Trade ID** | Presente | Único por operación |

**Ejemplo real (proyecto `agp`):**
```
Date: 2024-06-15 14:23:00
Type: Trade (BUY)
Buy: 100 USDT @ 1.00
Sell: —
Fee: 0.1 USDT (0.1%)
Exchange: BingX
```

### Margin Trading

| Campo | Diferencia |
|-------|-----------|
| **Type** | Trade (igual) |
| **Grupo** | "BingX Margin" (opcional) |
| **Fee** | Interés de préstamo (variable) |
| **Risk** | Liquidación si cae debajo de ratio |

**⚠️ Advertencia:** El trading apalancado genera pérdidas/ganancias amplificadas. La fiscalidad es tratamiento de capital (mismo que Spot).

### Perpetual Futures

| Campo | Valor |
|-------|-------|
| **Type** | "Perpetual Futures" o "Derivatives" |
| **Funding Fee** | Pagos periódicos (cada 8h) |
| **Position Type** | LONG / SHORT |
| **Liquidation Risk** | Sí (pérdida total posible) |

**Fiscalidad:** En España, puede ser:
- Opción A: Ganancia patrimonial (como spot)
- Opción B: Rendimiento del capital (si es actividad profesional)
**Recomendación:** Consultar asesor fiscal

---

## Integración con CoinTracking

### ✅ Métodos Soportados

1. **API Connection (Recomendado)**
   - CoinTracking soporta BingX API
   - Datos en vivo, actualizados automáticamente
   - Comisiones incluidas

2. **CSV Import**
   - BingX permite exportar Trade History
   - Formato: estándar (Date, Type, Coin, Amount, Price, Fee)
   - Limitación: Manual, requiere actualización periódica

### Importación via API

**Pasos:**
1. BingX: Account → API Management → Create API Key
2. CoinTracking: Settings → Exchanges → BingX → Conectar
3. Autorizar acceso (solo lectura recomendado)
4. Seleccionar período inicial (últimos 6-12 meses)
5. Sincronizar

**Validación:**
- [ ] Balance en CoinTracking coincide con BingX
- [ ] Operaciones mostradas: count correcto
- [ ] Comisiones incluidas
- [ ] Trade IDs únicos

---

## Casos Límite y Peculiaridades

### 1. Fee en Múltiples Monedas

**Problema:** BingX puede cobrar comisiones en:
- USDT (por defecto)
- BNB (reducción)
- La moneda vendida (si aplica)

**Solución:** CoinTracking registra automáticamente. Verificar que cost basis incluye la comisión.

### 2. Funding Fees (Perpetuos)

Cada 8 horas, BingX cobra/paga funding fee si tienes posiciones abiertas:
- Pago positivo → tú pagas (es fee)
- Pago negativo → tú recibes (es ingreso)

**En CoinTracking:**
- Aparece como "Fee" o "Income" según signo
- Incluirlo en cost basis si es fee

### 3. Copy Trading — ⚠️ Sub-cuenta NO exportada (crítico)

**Corregido 2026-07-05 tras caso real (agp2025):** la afirmación anterior de este documento ("cada operación del trader copiado genera una operación tuya, con Trade ID") **era incorrecta**. La evidencia real muestra lo contrario:

- BingX organiza el capital en **sub-cuentas** (Fund, Standard Futures, Perpetual Futures, Copy Trading, …).
- **CoinTracking (vía API o CSV) NO trae el historial de la sub-cuenta de Copy Trading.** No genera filas de tipo "Trade" ni Trade IDs — simplemente no aparece.
- El dinero transferido internamente desde Fund hacia Copy Trading es indistinguible, en las exportaciones disponibles, de una transferencia a Futuros normal.

**Cómo se detecta (indirectamente, por diferencia de saldos):**
```
Fund Account → transfirió X USDT a sub-cuentas de trading
Standard Futures + Perpetual Futures (exportados) → solo reciben Y USDT (Y < X)
Diferencia (X − Y) → atrapada en Copy Trading (no exportado), asumiendo
                      que no hay retiradas externas sin explicar
```

Esta diferencia **no es una operación individual**: es un ajuste de reconciliación de saldo agregado, registrado como una única fila `Type: Lost` con el importe neto perdido (no un CSV nativo de BingX).

**Caso real verificado (agp2025, cierre 2026-07-03):** ~694,67 USDT perdidos en la sub-cuenta de Copy Trading a lo largo de 2024. Confirmado por triangulación de saldos (Fund → Futuros exportados no cuadraba) y descartadas retiradas externas (ninguna fila "Withdraw" en ningún sheet de BingX). Es "casi seguro" Copy Trading por descarte, no una confirmación directa de BingX (que no expone esa sub-cuenta).

**Tratamiento fiscal:** Al no existir justificante documental de BingX para esa sub-cuenta, la pérdida se registra como `Lost` **no deducible**, no como pérdida de derivados. Solo se reclasificaría si el usuario obtiene de BingX un export específico de Copy Trading, o su asesor fiscal decide otra cosa. El agente no produce cifras fiscales vinculantes (ADR-006).

**Validación:** Si el saldo reconstruido de las sub-cuentas exportadas (Standard + Perpetual Futures) no cuadra contra lo que el Fund Account transfirió, sospecha de Copy Trading no exportado antes de asumir error de importación.

### 4. Liquidación de Posición

Si margin/futures es liquidado por caída de precio:
- Type: "Lost" o "Liquidation"
- Pérdida es total del colateral (en ese par)

**Fiscalidad:** Pérdida patrimonial (deducible).

---

## Validación en CoinTracking

### Checklist: Spot vs API vs CSV

```
[ ] API conectado en CoinTracking
[ ] Balance actual coincide (refresh)
[ ] Últimas 10 operaciones visibles
[ ] Comisiones incluidas
[ ] Trade IDs únicos (sin duplicados)
[ ] No hay solapamiento API+CSV
```

### Detección de Problemas Comunes

| Problema | Síntoma | Solución |
|----------|---------|----------|
| **API desconectado** | Balance no actualiza | Reconectar API en Settings |
| **CSV duplica** | Operaciones aparecen 2x | Eliminar CSV, usar solo API |
| **Fee en moneda ajena** | USDT fee no en BTC pair | Verificar cost basis incluye fee |
| **Margin liquidado** | Balance baja abruptamente | Buscar operación "Lost" |
| **Funding fees no vistos** | Pérdida sin operación clara | Buscar "Fee" o "Income" |

---

## Referencias y Recursos

- [BingX Official API Docs](https://bingx-api.github.io/) (inglés)
- [CoinTracking BingX Integration](https://www.cointracking.info/en/api_keys.php) (en plataforma)
- [CoinTracking FAQ: BingX](https://www.cointracking.info/en/faq/) (buscar BingX)

---

## Notas Operativas

**Para auditoría:**
- BingX es menor en volumen vs Binance pero creciente
- Estructura similar a Binance pero menos features
- Importante validar API vs CSV overlap
- Futures/Margin requieren cuidado con liquidaciones

**Para fiscalidad:**
- Spot: ganancias patrimoniales (FIFO)
- Margin: igual que Spot (+ riesgo de liquidación)
- Futures/Perpetuos: puede ser capital o rendimiento (consultar asesor)

---

**Documento:** BingX Mechanics  
**Nivel:** B2-010  
**Status:** Operacional  
**Creado:** 2026-07-05
