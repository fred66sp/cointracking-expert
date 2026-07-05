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
id: KB-B1-008
title: "Cómo CoinTracking calcula saldos (algoritmo de balance)"
level: B
domain: cointracking
source: "Análisis de casos + documentación técnica"
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
  - knowledge/cointracking/behavioral/PURCHASE_POOL_MECHANICS.md
  - knowledge/cointracking/official/COST_BASIS_AND_VALIDATION.md

tags:
  - cointracking
  - balance
  - algorithm
  - behavioral

notes: "Operativo: cómo CoinTracking mantiene saldos correctos mediante algoritmo."
---

# Cómo CoinTracking Calcula Saldos (Algoritmo de Balance)

## Concepto Base

**Balance = Cantidad total acumulada de un activo**

```
Ejemplo: BTC

Operación 1: Buy 1 BTC → Balance = 1 BTC
Operación 2: Buy 0.5 BTC → Balance = 1.5 BTC
Operación 3: Sell 0.3 BTC → Balance = 1.2 BTC
Operación 4: Deposit 0.1 BTC → Balance = 1.3 BTC

Balance final: 1.3 BTC
```

---

## Algoritmo de CoinTracking (Simplificado)

```
Balance(activo, fecha) = Σ(todas las operaciones hasta fecha)

Donde cada operación suma/resta según su tipo:

  BUY: +cantidad
  SELL: -cantidad
  DEPOSIT: +cantidad
  WITHDRAWAL: -cantidad
  INCOME: +cantidad
  EXPENSE: -cantidad
  TRANSFER OUT: -cantidad
  TRANSFER IN: +cantidad
  TRADE: -salida, +entrada (dos activos)
  
Fórmula general:
  Balance(X, fecha) = Balance(X, fecha anterior) + Δ(X, fecha)
  
Donde Δ(X, fecha) = suma de cambios de X en esa fecha
```

---

## Cálculo Paso a Paso

### Ejemplo Completo: Auditar Balance de BTC

```
Fecha | Operación | Cantidad | Balance | Nota
------|-----------|----------|---------|------
1 ene | Buy | +1 BTC | 1.0 BTC | Compra inicial
5 ene | Buy | +0.5 BTC | 1.5 BTC | Segunda compra
10 ene| Deposit | +0.1 BTC | 1.6 BTC | Transfer desde Kraken
15 ene| Sell | -0.3 BTC | 1.3 BTC | Venta parcial
20 ene| Withdrawal | -0.1 BTC | 1.2 BTC | Envío a cold wallet
25 ene| Trade | -0.5 ETH, +0.05 BTC | 1.25 BTC | Swap en Uniswap

Balance final (25 ene): 1.25 BTC
```

---

## Orden de Operaciones: FIFO para Ganancia

**IMPORTANTE:** El balance se calcula SUMA simple, pero la ganancia usa FIFO.

```
No es lo mismo:

1. BALANCE CÁLCULO:
   - Suma todas las operaciones de compra/venta
   - Resultado: balance de cantidad
   
2. GANANCIA CÁLCULO (para tax):
   - Usa FIFO para determinar cuál compra se vende
   - Resultado: ganancia/pérdida para IRPF
   
Ejemplo:
  Compra 1: 1 BTC @ 40.000€
  Compra 2: 1 BTC @ 50.000€
  Venta: 1 BTC @ 60.000€
  
  Balance: 1 BTC (correcto, tienes 1)
  Ganancia: 60.000€ - 40.000€ = 20.000€ (FIFO vende la primera compra)
  Si fuera promedio: 60.000€ - 45.000€ = 15.000€ (INCORRECTO en España)
```

---

## Problema: Balance Negativo

**¿Cuándo CoinTracking calcula balance negativo?**

```
Síntoma:
  Balance de BTC: -0.5 BTC (¡IMPOSIBLE!)
  
¿Cómo sucede?

Causa: Operación sin origen documentado

Ejemplo:
  Venta: -1 BTC (pero no hay compra documentada)
  Balance: -1 BTC (rojo)
  
¿Qué pasó?
  CoinTracking ve una venta de BTC que "no tienes"
  
Razones:
  1. Depósito fiat → BTC no importado
  2. Transfer externo sin documentación
  3. Airdrop/mining no registrado
  4. Importación parcial (CSV incompleto)
```

---

## Validación en CoinTracking

### Ver Balance por Activo

```
CoinTracking → Home (Dashboard):
  Muestra balance total en € por cada activo
  
¿Dónde verificar?
  CoinTracking → Holdings (o Portfolio)
  Columna: "Units" o "Quantity"
  
Busca:
  [ ] ¿Algún balance es negativo?
      SÍ → Hay un problema
      NO → OK
      
  [ ] ¿Los balances son consistentes con tus wallets?
      Compara con:
        - Saldo en Binance
        - Saldo en tu cold wallet
        - Etc.
      SÍ → Datos sincronizados
      NO → Falta importación o hay error
```

### Ver Balance Histórico

```
CoinTracking → Reports → Balance History:
  Muestra cómo evolucionó el balance en el tiempo
  
Gráfico esperado:
  - Balance sube cuando compras/recibes
  - Balance baja cuando vendes/envías
  - Nunca debe ser negativo (excepto error)
  
Si ves dip negativo:
  → Hay operación sin origen documentado
  → Necesita corrección antes de auditar
```

---

## Caso Especial: Múltiples Exchanges Simultáneos

**¿Qué pasa si tienes fondos en varios lugares?**

```
Situación:
  - Binance: 1 BTC
  - Kraken: 0.5 BTC
  - Cold wallet: 0.3 BTC
  - Total real: 1.8 BTC
  
¿Cómo CoinTracking suma?

Ideal:
  CoinTracking suma TODOS los depósitos/fondos
  Balance total: 1.8 BTC
  
Realidad:
  Si no importas Kraken:
    Balance en CT: 1.3 BTC (solo Binance + Cold)
    Balance real: 1.8 BTC
    DIFERENCIA: 0.5 BTC faltantes
    
¿Cómo detectar?
  [ ] Importa TODAS las fuentes (Binance API, Kraken API, etc.)
  [ ] Suma los balances en CT
  [ ] Compara con realidad (suma manual de wallets)
  [ ] ¿Son iguales?
      SÍ → Datos completos
      NO → Falta una fuente
```

---

## Algoritmo de Balance por Fecha

**CoinTracking recalcula balance cada vez que añades una operación:**

```
Cuando editas/añades una transacción:
  
  1. CoinTracking busca la fecha
  2. Reordena las operaciones de esa fecha
  3. Recalcula balance desde esa fecha EN ADELANTE
  4. Actualiza todos los balances posteriores
  
Impacto:
  - Si añades una compra retroactiva (ej. del 1 enero)
  - CoinTracking recalcula balance de 1 enero hasta hoy
  - Las ganancias pueden cambiar (porque FIFO cambió)
  
Ejemplo:
  Tenías 10 ganancias en enero
  Añades compra el 15 enero
  Las ganancias de enero NO cambian (pero enero 15+ sí)
```

---

## Validación: Checklist de Balance

```
Antes de auditar:

[ ] ¿Todos los saldos son positivos?
    NO → Corregir operaciones sin origen
    
[ ] ¿Los saldos coinciden con realidad?
    Comparar:
      - CoinTracking BTC balance
      - Suma de tus direcciones BTC reales
    ¿Iguales?
      NO → Falta importación de fuente
      
[ ] ¿El balance evoluciona lógicamente?
    Ver gráfico de balance histórico
    ¿Tiene "dips" raros (caídas inesperadas)?
      SÍ → Hay transacción sin origen
      
[ ] ¿Se recalculó el balance correctamente?
    Después de editar transacciones:
      Regenera Reports
      ¿Los datos cambiaron correctamente?
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Balance es suma acumulada
- **PURCHASE_POOL_MECHANICS.md:** Pool se extrae del balance para ganancias
- **COST_BASIS_AND_VALIDATION.md:** Balance debe ser consistente con cost basis
