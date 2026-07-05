---
id: KB-B2-009
title: "Mecánicas Genéricas de Exchanges: Cómo Manejar Exchanges No Documentados"
level: B
domain: cointracking
source: "Patrón genérico + análisis"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-03
confidence: medium
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/cointracking/behavioral/BINANCE_IMPORT_WORKFLOW.md
  - knowledge/cointracking/behavioral/API_VS_CSV_OVERLAP.md

tags:
  - exchange
  - generic
  - import
  - behavioral
  - advanced

notes: "Cómo auditar exchanges no documentados específicamente (Bybit, OKX, Huobi, etc)."
---




# Mecánicas Genéricas de Exchanges: Exchanges No Documentados

## El Problema: Exchanges Menos Comunes

**Exchanges comunes documentados:**
- Binance (B2-001)
- Kraken (B2-006)
- Coinbase (B2-007)

**Exchanges menos comunes (sin doc específica):**
- Bybit
- OKX (antiguo OKCoin)
- Huobi
- Crypto.com
- XT.com
- Otros

**Solución:** Usar patrón genérico.

---

## Patrones de Todos los Exchanges

Aunque cada exchange es distinto, todos siguen patrones comunes:

```
1. Depósitos (fiat o cripto) → Saldo aumenta
2. Retiros (fiat o cripto) → Saldo disminuye
3. Compras/Ventas → Cambio de activos
4. Transfers internos → Movimiento entre wallets (del mismo usuario)
5. Comisiones → Reducen el saldo
6. Rewards → Aumentan el saldo (si hay staking/earn)
```

---

## Importación Genérica (3 Opciones)

### Opción 1: API Directa (Si Existe)

**CoinTracking soporta API de 200+ exchanges.**

**Pasos:**
1. CoinTracking → Settings → Exchanges
2. Busca el exchange (p. ej. "OKX")
3. Si aparece → conectar con API keys
4. Si no aparece → usa CSV

**Ventaja:** Datos en vivo, completo, actualizado

**Limitación:** No todos los exchanges están soportados

### Opción 2: CSV Genérico

**Pasos:**
1. Accede al exchange
2. Busca: "Download History", "Export CSV", "Trade History"
3. Descarga en formato CSV
4. Asegúrate que incluya:
   - Fecha/hora
   - Tipo (Buy/Sell/Deposit/Withdrawal/Fee)
   - Cantidad
   - Precio
   - Comisión

5. Importa en CoinTracking → Settings → Import from CSV

**Ventaja:** Funciona con cualquier exchange

**Limitación:** Manual, requiere actualizar periódicamente

### Opción 3: Importación Manual (Último Recurso)

**Si el exchange no tiene export y no está en CoinTracking:**

1. Anotar manualmente cada operación (tedioso)
2. O usar herramienta tercera:
   - Koinly (importa de +700 exchanges)
   - Zapper (para DeFi)
   - Etherscan (para Ethereum)

---

## Validación Genérica (4 Pasos)

### Paso 1: Verificar Completitud

```
CoinTracking muestra: N operaciones
Exchange real muestra: M operaciones

¿N == M?
  SÍ → Completo, OK
  NO → Faltantes = M - N operaciones
       → Buscar si están duplicadas o filtradas
```

### Paso 2: Verificar Saldo

```
CoinTracking balance (hoy): X BTC
Exchange balance (real, hoy): X BTC

¿Coinciden?
  SÍ → Importación correcta
  NO → Hay operaciones faltantes o duplicadas
```

### Paso 3: Verificar Comisiones

```
¿El CSV incluye comisiones por separado?
  SÍ → OK, aparecen como "Fee" o "Comisión"
  NO → Pueden estar incluidas en el precio
       → Verificar contra el exchange
```

### Paso 4: Verificar Duplicados

```
¿Hay operaciones duplicadas (misma fecha, cantidad, precio)?
  SÍ → Investigar:
       - ¿Mismo Trade ID?
       - ¿Misma API + CSV (overlap)?
  NO → OK
```

---

## Tratamiento de Operaciones Especiales

### Depósitos Fiat

```
Tipo: Deposit (EUR, USD, etc)
Origen: Bank transfer / SEPA / Wire
Tratamiento: Base de coste (cost basis)

Ejemplo:
  Date: 2024-01-15
  Type: Deposit
  Amount: 1000 EUR
  Exchange: OKX
  
En CoinTracking: Base de coste = 1000 EUR
```

### Conversiones Internas

```
Tipo: Trade (USDT → USDC, etc)
Nota: Dentro del mismo exchange, sin movimiento blockchain

Tratamiento:
  - Venta de USDT
  - Compra de USDC
  - Ganancia/pérdida por diferencia de precio

Ejemplo:
  1000 USDT @ 1.00 USD = 1000 USD
  Cambia por: 1000 USDC @ 1.01 USD = 1010 USD
  Ganancia: 10 USD (por diferencia de precio)
```

### Transfers Entre Tus Wallets

```
Tipo: Transfer interno (si el exchange lo permite)
O: Withdrawal + Deposit (en dos exchanges)

NO genera ganancia/pérdida (es transferencia propia)
SÍ requiere validación (deposito = retiro)

Ejemplo:
  Retiras 1 BTC de Binance el 2024-01-15 14:00
  Recibes 1 BTC en OKX el 2024-01-15 14:15
  
  → No es ganancia, es transferencia
  → Verificar que montos cuadren (menos fees)
```

---

## Checklist: Auditar un Exchange Desconocido

```
[ ] Exchange está en CoinTracking (buscar)
    SÍ → Usar API
    NO → Continuar
    
[ ] Exchange permite exportar CSV
    SÍ → Descargar
    NO → Usar herramienta tercera (Koinly, Zapper)
    
[ ] CSV tiene todas las columnas necesarias
    (Date, Type, Amount, Price, Fee)
    SÍ → Importar
    NO → Conseguir mejor export
    
[ ] Validar completitud (count de operaciones)
[ ] Validar saldo (coincide con exchange real)
[ ] Validar comisiones (están registradas)
[ ] Detectar duplicados (si usaste múltiples fuentes)

[ ] ¿Todo OK?
    SÍ → Exchange listo para auditar
    NO → Investigar antes de continuar
```

---

## Exchanges Especiales (Casos Límite)

### DEX (Uniswap, SushiSwap, etc)

```
¿Qué son?
  - No son exchanges centralizados
  - Las operaciones son on-chain
  - CoinTracking NO las importa de DEX directamente
  
Solución:
  1. Ver operaciones en Etherscan / blockchain explorer
  2. Importar manualmente en CoinTracking
  O
  1. Usar herramientas como Zapper
  2. Exportar desde Zapper a CSV
  3. Importar en CoinTracking
```

### Staking Pools (Lido, Rocket, etc)

```
¿Qué son?
  - No son exchanges
  - Protocolos DeFi de staking
  - CoinTracking puede importarlos si hay API
  
Solución:
  1. Ver en etherscan/blockchain explorer
  2. Importar manualmente
  3. O conectar wallet directamente si CoinTracking lo soporta
```

---

## Integración

- **BINANCE_IMPORT_WORKFLOW.md:** Patrón específico (Binance)
- **API_VS_CSV_OVERLAP.md:** Evitar duplicados API+CSV
