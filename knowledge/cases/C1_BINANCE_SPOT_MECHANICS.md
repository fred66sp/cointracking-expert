---
id: KB-C1-001
title: "Binance Spot: Mecánicas de importación y peculiaridades"
level: C
domain: cointracking
source: "Auditoría agp2025 + casos reales"
authority: verified
last_verified: 2026-07-05
valid_until: 2027-07-05
confidence: high
version: 1.0

tags:
  - binance
  - spot-trading
  - import-mechanics
  - common-issues

---

# Binance Spot: Importación y Peculiaridades

**Tipo:** Patrón operativo (Nivel C — casos reales)  
**Cuando aplica:** Usuario audita cuenta con Binance Spot  
**Caso base:** agp2025 (auditoría 2026-07-03)

---

## 1. Cómo Importa Binance en CoinTracking

### Fuentes de Datos

Binance tiene **3 APIs/exports distintos** (causa de confusión):

1. **API pública:** Solo histórico desde sept-2022
2. **Transaction History (CSV):** Todas las transacciones (depósitos, retiradas, conversiones)
3. **Order History (CSV):** Solo operaciones spot (compra/venta)

**CoinTracking soporta:**
- ✅ Import vía API (automático)
- ✅ Import manual (CSV descargado)
- ⚠️ A veces hay overlap: API + CSV importados = duplicados

---

## 2. Peculiaridades Comunes

### A. Dust → BNB (Auto-conversión)

**¿Qué es?**  
Si tienes < 10 USDT o pequeños saldos de altcoins, Binance los convierte automáticamente a BNB.

**Cómo aparece en CoinTracking:**
```
Type: Trade
Currency: USDT
Amount: -0.50
Fee Currency: BNB
Fee Amount: 0.00001
To Currency: BNB
To Amount: 0.000008
```

**Auditoría:**
- ✅ Es una venta legítima (USDT → BNB)
- ✅ Genera ganancia/pérdida (aunque sea mínima)
- ✅ Tributa (Art. 37.1.h LIRPF — permuta cripto-cripto)

**Común error:** Pensar que es "gratis" (no lo es, tributa).

---

### B. Binance Convert (Cambio Interno)

**¿Qué es?**  
Herramienta de Binance para convertir entre criptos sin pasar por orderbook.

**Cómo aparece:**
```
Type: Trade
Description: "Binance Convert USDT to USDC"
Currency: USDT
Amount: -1000
To Currency: USDC
To Amount: 999.95
Fee: 0.05 USDC
```

**Auditoría:**
- ✅ Es una permuta (tributa)
- ✅ La comisión (0.05 USDC) es parte del coste
- ⚠️ Puede confundirse con depósito si no se lee bien

**Caso real (agp2025):** USDT→USDC conversión Q1 2025 (MiCA). ✅ Verificado sin discrepancias.

---

### C. Swaps (Dex Integraciones)

**¿Qué es?**  
Binance permite hacer swaps con protocolos DeFi (Uniswap, 1inch) desde la app.

**Cómo aparece:**
```
Type: Trade
Description: "Swap: ETH to MATIC via 1inch"
Currency: ETH
Amount: -1.0
To Currency: MATIC
To Amount: 2500
Fee: 0.005 ETH (slippage)
```

**Auditoría:**
- ✅ Es una permuta (tributa)
- ✅ Fee/slippage es parte del coste
- ⚠️ Verificar que CoinTracking lo importa correctamente (a veces lo pierde)

---

### D. Binance Earn/Staking (No Se Importa Bien)

**⚠️ Problema crítico:**  
Binance Earn (flexible, bloqueado) frecuentemente **no se importa en CoinTracking**.

**Síntomas:**
- Balance de USDT en Earn no aparece en CoinTracking
- Recompensas de staking ausentes
- Saldo en CoinTracking ≠ saldo en app

**Solución:**
1. Exportar manualmente desde Binance: Wallet > Earn History
2. Crear operaciones manuales en CoinTracking:
   ```
   Tipo: Income
   Moneda: USDT
   Cantidad: 150 (recompensa del mes)
   Fecha: 30/06/2025
   Comentario: "Binance Earn flexible recompensa"
   ```

---

## 3. Auditoría Paso a Paso (Binance Spot)

### Checklist

- [ ] **Import completo:**
  - ¿Cuántas operaciones hay?
  - ¿Desde qué fecha?
  - ¿API o CSV?

- [ ] **Duplicados:**
  - ¿Hay operaciones repetidas? (API + CSV importados)
  - Verificar Trade ID (único en Binance)

- [ ] **Dust & Conversiones:**
  - ¿Hay conversiones pequeñas (< 1 USDT)?
  - ¿Se cuenta como venta (tributa)?

- [ ] **Saldo final:**
  - ¿Coincide con app de Binance?
  - ¿Falta algún activo (Earn, staking)?

- [ ] **Base de coste:**
  - ¿Hay ventas sin compra previa?
  - ¿Se arrastra correctamente FIFO?

---

## 4. Errores Comunes y Soluciones

| Error | Síntoma | Causa | Solución |
|-------|---------|-------|----------|
| Duplicados | Mismo trade 2 veces | API + CSV importados | Eliminar duplicado verificando Trade ID |
| Saldo negativo | BTC: -0.5 | Venta sin compra | Importar histórico o crear operación manual |
| Earn ausente | USDT balance bajo | Earn no se importa | Crear manuales desde historial Earn |
| Swap invisible | ETH desaparece | CoinTracking no lo captura | Crear manual o re-importar con CSV |
| Dust ignorado | pequeño saldo | Binance lo convirtió a BNB | Buscar conversión USDT→BNB |

---

## 5. Caso Real: agp2025

**Estado:** ✅ Limpio

- 1.000+ operaciones Binance Spot
- Importadas por API desde 2024
- Verificación: Saldo final = app real ✓
- Duplicados: 0 (Trade IDs únicos) ✓
- Base de coste: Completa (FIFO verifi cado) ✓
- Dust: Presente pero contabilizado ✓

---

**Próximo:** Si encuentras un patrón diferente, avisa y lo documentamos en C1-002.
