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
id: KB-B1-007
title: "Causas de Missing Purchase History en CoinTracking"
level: B
domain: cointracking
source: "Casos CT-003 + análisis"
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
  - knowledge/cases/CT-003-missing-purchase-history.md
  - knowledge/cointracking/behavioral/PURCHASE_POOL_MECHANICS.md

tags:
  - cointracking
  - purchase-history
  - cost-basis
  - diagnostic
  - behavioral

notes: "Operativo: cómo y por qué CoinTracking muestra 'missing purchase history'."
---

# Causas de Missing Purchase History en CoinTracking

## Síntoma: Missing Purchase History

**¿Cuándo aparece?**

```
CoinTracking → Transacciones → Filter: "Missing Purchase History"

O en Reports → Gains:
  "No purchase history found for this transaction"
  "Cost basis: Unknown"
```

**¿Qué significa?**

```
CoinTracking detectó una VENTA pero NO encontró
una COMPRA anterior que la respalda.

Ejemplo:
  Vendiste: 100 USDC @ 1 USDT (ganancia 0€)
  Pero: CoinTracking no vio dónde compraste ese USDC
  
Resultado:
  - Cost basis desconocido
  - Ganancia no puede calcularse
  - Audit fallida
```

---

## Causa 1: Depósito Inicial (Fiat → Cripto) No Documentado

**La causa más común**

```
Situación:
  Abriste cuenta en Binance
  Depositaste 1.000€ (fiat) en euros
  Compraste 1 BTC
  Vendiste después
  
¿Dónde está el problema?
  - Importaste solo la VENTA de BTC
  - Pero NO importaste la COMPRA original
  
¿Por qué?
  - Binance no exporta "depósitos fiat" como operaciones
  - CSV de CoinTracking solo muestra cripto → cripto
  - El BTC "aparece" sin origen
  
Resultado: Missing Purchase History
```

---

## Causa 2: Importación Incompleta (CSV Parcial)

```
Situación:
  Descargaste CSV de Binance del rango 1 enero - 31 enero
  Pero compraste BTC el 15 diciembre (FUERA del rango)
  Vendiste BTC el 15 enero (DENTRO del rango)
  
¿Dónde está el problema?
  - La VENTA está en el CSV (15 enero)
  - La COMPRA NO está en el CSV (15 diciembre)
  
Resultado: Missing Purchase History
```

---

## Causa 3: Transferencias Entre Wallets (Depósitos Externos)

```
Situación:
  Recibiste 0.5 BTC en tu wallet Binance
  Viniendo de otro wallet (Kraken, Coinbase, etc.)
  
¿Dónde está el problema?
  - CoinTracking solo ve: "Deposit 0.5 BTC" (sin origen)
  - NO ve la venta anterior en Kraken
  
¿Cómo apareció?
  - Como "Deposit" o "Transfer in" genérico
  - Sin documentación de base de coste
  
Resultado: Si vendes después, Missing Purchase History
```

---

## Causa 4: Airdrop, Staking, o Mining (Ingresos Gratuitos)

```
Situación:
  Recibiste 100 TOKEN gratis (airdrop)
  O ganancias de staking
  O recompensas de mining
  
CoinTracking los registra como:
  - "Airdrop" (si lo detecta)
  - "Income" (si lo importa bien)
  - "Deposit" (si no sabe)
  
Resultado:
  - Si vendes después: "Cost basis = 0"
  - NO es "Missing" pero es "Unknown"
```

---

## Causa 5: Transacción Blockchain Manual o No Importada

```
Situación:
  Hiciste swap en Uniswap (DeFi, on-chain)
  El CSV de Etherscan muestra dos transfers separados
  Importaste solo el "transfer out", no el "transfer in"
  
CoinTracking ve:
  - Sale: 1.000 USDC (OK)
  - No entrada de ETH (FALTA)
  
Resultado: Missing Purchase History para el ETH "fantasma"
```

---

## Diagnóstico: Árbol de Decisión

```
¿Ves "Missing Purchase History"?
  
  SÍ → ¿Es un DEPÓSITO que no tiene origen?
    SÍ → Causa 1 o 3 (depósito fiat o transfer externo)
    NO → ¿Es una VENTA?
      SÍ → ¿Hay compra más antigua documentada?
        SÍ → Error raro, contacta soporte
        NO → Importación incompleta (Causa 2)
      NO → ¿Es un AIRDROP/REWARD?
        SÍ → Normal, cost basis = 0 (Causa 4)
        NO → Transfer blockchain faltante (Causa 5)
```

---

## Soluciones por Causa

### Causa 1: Depósito Inicial Fiat

```
Solución: Documentar la compra inicial

Opción A (mejor): Añadir transacción manual
  CoinTracking → Add Transaction
  
  Tipo: Buy
  Divisa: EUR → BTC (o lo que compres)
  Cantidad: 1 BTC
  Precio: 40.000€ (precio histórico de ese día)
  Fecha: cuando hiciste el depósito+compra
  Exchange: Binance
  
Opción B: Editar el depósito existente
  Si el "Deposit" ya existe, editar su "Cost basis"
  Añadir el valor que pagaste

Resultado: CoinTracking ahora calcula ganancia correctamente
```

### Causa 2: Importación Incompleta

```
Solución: Reimportar CSV con rango correcto

Paso 1: En Binance
  Descargar CSV: Rango COMPLETO (desde la compra original)
  Incluir: Buy del 15 dic + Sell del 15 ene
  
Paso 2: En CoinTracking
  Transacciones → Import from CSV
  Seleccionar el nuevo CSV (completo)
  
Paso 3: Verificar
  CoinTracking → Reports → Check para el mismo activo
  ¿Desapareció "Missing Purchase History"?
    SÍ → OK
    NO → Ir a Causa 3/4/5
```

### Causa 3: Transfer Externo

```
Solución: Documentar el origen

Opción A: Si conoces la otra plataforma
  - Descargar CSV de Kraken (dónde vendiste)
  - Importar en CoinTracking como "Kraken sell"
  - Resultado: Se vinculan automáticamente

Opción B: Si es transferencia heredada/regalo
  - Crear transacción manual "Buy"
  - Cantidad: 0.5 BTC
  - Precio: precio de mercado del día que recibiste
  - Exchange: "Transfer externo" (etiqueta personalizada)

Opción C: Si es de blockchain (on-chain transfer)
  - Usar base de coste de la dirección anterior
  - O documentar como "Transfer in" con precio de mercado ese día
```

### Causa 4: Airdrop/Staking/Mining

```
Solución: Clasificar correctamente

Para auditoría:
  CoinTracking → Editar el ingreso
  Tipo: DEBE ser "Income" o "Airdrop"
  Cost basis: 0€ (fue gratuito)
  Precio: precio de mercado en el momento que recibiste
  
Resultado:
  - Aparece como "Rendimiento del capital" en IRPF
  - La venta posterior calcula ganancia correctamente
  - "Missing Purchase History" desaparece
```

### Causa 5: Transacción Blockchain Incompleta

```
Solución: Buscar y completar la transacción

Paso 1: Verificar en Etherscan (o explorer)
  - Buscar el TXN ID
  - Ver TODOS los outputs (pueden ser múltiples)
  
Paso 2: Importar el CSV de Etherscan completo
  - O crear transacciones manuales por cada output
  
Paso 3: En CoinTracking
  - Relacionar "Transfer out" + "Transfer in"
  - Como un único Swap/Trade
  
Resultado: Cost basis definido correctamente
```

---

## Validación en CoinTracking

```
Después de cualquier corrección:

CoinTracking → Reports → Gains:
  ¿Desapareció "Missing Purchase History"?
    SÍ → Ganancia calculada, OK
    NO → Hay otro problema, revisar causas
    
¿La ganancia calculada es sensata?
  Esperado: (Precio venta - Costo) × cantidad
  
Si no cuadra:
  → Verificar cost basis en detalle
  → Comparar contra facturas/extracts de exchange
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Cost basis es obligatorio
- **PURCHASE_POOL_MECHANICS.md:** Cómo se usan las compras para calcular ganancias
- **CT-003:** Caso de referencia
