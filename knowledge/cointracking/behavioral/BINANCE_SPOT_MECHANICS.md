---
id: KB-B1-XXX
title: "Untitled Document"
level: B
domain: cointracking
source: "Internal documentation"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - todo
  - needs-review

notes: "Metadatos agregados automáticamente. Verificar y actualizar conforme ADR-032."
---

---
id: KB-B2-001
title: "Cómo CoinTracking maneja Binance Spot (compras/ventas normales)"
level: B
domain: cointracking
source: "Casos reales + análisis"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

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

## Integración

- **ADR-003:** Modelo de transacciones — Spot Buy/Sell
- **CAPITAL_GAINS.md:** Cálculo de ganancias patrimoniales
