---
id: KB-B2-001
title: "Cómo CoinTracking maneja Binance Spot (compras/ventas normales)"
level: B
domain: cointracking
source: "Casos reales + análisis + auditoría agp2025"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-03
confidence: high
version: 1.1

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/taxonomy/spain/CAPITAL_GAINS.md

tags:
  - cointracking
  - binance
  - spot
  - behavioral

notes: "Operativo: cómo CoinTracking importa y trata operaciones Spot de Binance."
---




# Cómo CoinTracking Maneja Binance Spot

## Definición

**Binance Spot** = Mercado de contado: compras/ventas de criptos al precio actual (delivery inmediata).

**Equivalente en bolsa:** Comprar/vender acciones al precio de mercado sin apalancamiento.

---

## Cómo CoinTracking Registra Spot

### Flujo ideal (con API de Binance)

```
Si conectas CoinTracking a tu API de Binance:
  CoinTracking importa automáticamente:
    - Cada orden completada (Buy o Sell)
    - Precio, cantidad, fee, timestamp exacto
    - Los clasifica como "Buy" o "Sell" correctamente
  
Resultado:
  - Base de coste exacta
  - Ganancias/pérdidas calculadas correctamente
  - Tax Report preciso
```

### Flujo manual (CSV)

```
Si descargas CSV desde Binance manualmente:
  CoinTracking puede importar directamente si:
    - CSV está en formato estándar de Binance
    - Incluye: date, symbol, side (Buy/Sell), qty, price, commission
  
Si el CSV está mal formateado:
    - Tendrás que añadir operaciones manualmente
```

---

## Estructura de un Trade Spot

### Ejemplo: Compra de BTC

```
Binance → Spot Buy:
  - Pair: BTCUSDT
  - Side: Buy
  - Quantity: 0.5 BTC
  - Price: 50.000 USDT
  - Total spent: 25.000 USDT
  - Fee: 25 USDT (0.1%)
  - Net cost: 25.025 USDT
  
CoinTracking registra:
  - Tipo: Buy
  - Cantidad: 0.5 BTC
  - Precio: 50.025 USDT/BTC (incluye fee)
  - Cost base: 25.025 USDT
```

### Ejemplo: Venta de BTC

```
Binance → Spot Sell:
  - Quantity: 0.25 BTC
  - Price: 55.000 USDT/BTC
  - Received: 13.750 USDT
  - Fee: 13.75 USDT (0.1%)
  - Net received: 13.736.25 USDT
  
CoinTracking calcula:
  - Venta de 0.25 BTC
  - Cost base consumida: 0.25 × 50.025 = 12.506.25 USDT
  - Ingresos netos: 13.736.25 USDT
  - Ganancia: 13.736.25 - 12.506.25 = 1.230 USDT
```

---

## Validación en CoinTracking

```
CoinTracking → Transacciones:
  Filtra por "Binance" o "BNANCE"
  
¿Todas las operaciones están como "Buy" o "Sell"?
  SÍ → OK
  NO → Verificar y corregir tipos
  
¿Los precios y fees incluyen comisión?
  SÍ → OK
  NO → Editar y añadir fee
```

---

## Tratamiento Fiscal (España, IRPF)

**Spot Sell = Venta de cripto (ganancia patrimonial)**

```
Regla FIFO (DGT):
  - Se vende primero lo comprado primero
  - Ganancia = Precio venta - Precio compra
  - Se registra como "Ganancia patrimonial" en IRPF
  
Ejemplo:
  Compré: 1 BTC a 40.000€
  Vendí: 1 BTC a 60.000€
  Ganancia: 60.000€ - 40.000€ = 20.000€
```

---

## Casos Límite y Peculiaridades

Binance tiene varios mecanismos internos que generan operaciones Spot "no obvias" — todas tributan aunque no lo parezca:

### Dust → BNB (auto-conversión)

Si tienes saldos residuales (< 10 USDT) de altcoins, Binance los convierte automáticamente a BNB.

```
CoinTracking registra:
  Type: Trade
  Currency: USDT (o el altcoin) — Amount: -0.50
  To Currency: BNB — To Amount: 0.000008
```

**Es una venta legítima y tributa** (Art. 37.1.h LIRPF — permuta cripto-cripto), aunque el importe sea mínimo. Error común: asumir que es "gratis".

### Binance Convert (cambio interno sin orderbook)

Herramienta de Binance para convertir entre criptos sin pasar por el libro de órdenes.

```
Type: Trade — Description: "Binance Convert USDT to USDC"
Currency: USDT (-1000) → To Currency: USDC (999.95) — Fee: 0.05 USDC
```

Es una permuta (tributa); la comisión forma parte del coste. Puede confundirse con un depósito si no se lee la descripción completa.

### Swaps (integraciones DeFi desde la app)

Binance permite swaps con protocolos DeFi (Uniswap, 1inch) desde su propia interfaz.

```
Type: Trade — Description: "Swap: ETH to MATIC via 1inch"
Fee: 0.005 ETH (slippage)
```

Tributa como permuta; el slippage/fee forma parte del coste. Verificar que CoinTracking lo importe — a veces se pierde.

### Binance Earn / Staking (con frecuencia no se importa)

**Problema crítico:** el balance en Earn (flexible o bloqueado) y sus recompensas frecuentemente **no llegan a CoinTracking** vía API. Síntoma: saldo de CoinTracking ≠ saldo real de la app.

**Solución:** exportar manualmente desde Binance (Wallet → Earn History) y crear las operaciones de tipo `Income` en CoinTracking. Ver `STAKING_MECHANICS.md` para el tratamiento fiscal completo.

---

## Checklist de Auditoría (Binance Spot)

- [ ] **Import completo:** ¿cuántas operaciones, desde qué fecha, API o CSV?
- [ ] **Duplicados:** ¿hay operaciones repetidas por solapamiento API+CSV? Verificar Trade ID único.
- [ ] **Dust y conversiones:** ¿hay conversiones pequeñas (< 1 USDT)? ¿se cuentan como venta?
- [ ] **Saldo final:** ¿coincide con la app de Binance? ¿falta Earn/staking?
- [ ] **Base de coste:** ¿hay ventas sin compra previa? ¿FIFO arrastra correctamente?

**Caso real verificado (agp2025):** 1.000+ operaciones Binance Spot importadas por API desde 2024. Saldo final = app real ✓. Duplicados = 0 (Trade IDs únicos) ✓. Base de coste completa (FIFO verificado) ✓. Dust presente y contabilizado ✓.

---

## Integración

- **ADR-003:** Modelo de transacciones — Spot Buy/Sell
- **CAPITAL_GAINS.md:** Cálculo de ganancias patrimoniales
