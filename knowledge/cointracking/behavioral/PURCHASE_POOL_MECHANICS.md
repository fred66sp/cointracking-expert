---
id: KB-B1-005
title: "Cómo funciona el Pool de Compras en CoinTracking (FIFO)"
level: B
domain: cointracking
source: "Análisis de CoinTracking + documentación oficial"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-03
confidence: high
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/cointracking/official/COST_BASIS_AND_VALIDATION.md
  - knowledge/taxonomy/spain/CAPITAL_GAINS.md

tags:
  - cointracking
  - purchase-pool
  - fifo
  - cost-basis
  - behavioral

notes: "Crítico: cómo CoinTracking rastrea la base de coste mediante FIFO."
---




# Cómo Funciona el Pool de Compras en CoinTracking (FIFO)

## Concepto: El Pool de Compras

**Pool de Compras** = Lista de todas las compras de un activo, en orden.

```
Ejemplo: BTC

Compra 1: 1 BTC @ 40.000€ (1 enero)
Compra 2: 0.5 BTC @ 50.000€ (15 enero)
Compra 3: 0.2 BTC @ 45.000€ (20 enero)

Pool actual:
  [1 BTC @ 40.000€] [0.5 BTC @ 50.000€] [0.2 BTC @ 45.000€]
  Total: 1.7 BTC, costo total: 77.500€
  Costo medio: 77.500€ / 1.7 ≈ 45.588€/BTC
```

---

## FIFO: First In, First Out

**Regla (España, DGT):** Vendes primero lo que compraste primero.

```
Pool actual:
  [1 BTC @ 40.000€] [0.5 BTC @ 50.000€] [0.2 BTC @ 45.000€]

Vendes 1.2 BTC @ 60.000€:
  1. Consume 1 BTC @ 40.000€ (la más antigua)
  2. Consume 0.2 BTC @ 50.000€ (la segunda)
  
Ganancia calculada:
  - 1 BTC: 60.000€ - 40.000€ = 20.000€
  - 0.2 BTC: 60.000€ - 50.000€ = 2.000€
  - Total: 22.000€

Pool restante después:
  [0.3 BTC @ 50.000€] [0.2 BTC @ 45.000€]
  Total: 0.5 BTC
```

---

## Cómo CoinTracking Gestiona el Pool

### Visualizar el pool

```
CoinTracking → Transactions:
  Cada BUY se añade al pool (en orden de fecha)
  
CoinTracking → Reports → Gains:
  Cuando vendes, CoinTracking automáticamente:
    1. Busca compras no consumidas (FIFO)
    2. Calcula ganancia para cada tramo
    3. Suma la ganancia total
```

### Problema: Pool Vacío (All pools consumed)

```
Síntoma:
  CoinTracking muestra: "All pools consumed"
  
¿Qué significa?
  - Intentaste vender más BTC del que tienes base de coste
  - Ejemplo: Tienes 2 BTC pero solo 1.5 BTC en compras documentadas
  
Causa común:
  - Depósito inicial sin origen documentado
  - BTC recibido de airdrop/reward pero no importado
  - Base de coste incompleta
  
Impacto:
  - Ganancia calculada puede ser 0 o incorrecta
  - Auditoría fallida
```

### Cómo corregir "All pools consumed"

```
Paso 1: Identificar qué BTC faltan
  CoinTracking → Gains → Ver la venta problemática
  Anotar: ¿Cuánto vendiste? ¿Cuánto está documentado en compras?
  Diferencia = lo que falta
  
Paso 2: Buscar el origen del BTC faltante
  - ¿Fue un depósito inicial (fiat → cripto)?
  - ¿Fue recibido de otro wallet?
  - ¿Fue airdrop/reward?
  
Paso 3: Documentar la compra faltante
  CoinTracking → Add Transaction:
    Tipo: Buy
    Cantidad: BTC faltante
    Precio: precio histórico (de ese día)
    Fecha: cuando recibiste el BTC
    Exchange: donde obtuviste (Binance, banco, etc.)
```

---

## Ejemplo Completo: Auditoría del Pool

```
Situación: Alfredo tiene 2 BTC. Compras documentadas:
  - Compra 1: 1 BTC @ 30.000€ (1 enero 2024)
  - Compra 2: 0.5 BTC @ 40.000€ (15 enero 2024)
  
Total documentado: 1.5 BTC (costo: 50.000€)
Tengo: 2 BTC

¿Dónde están los 0.5 BTC restantes?
  - No están documentados
  - Pool está incompleto
  
Vendo 2 BTC a 60.000€:
  CoinTracking calcula:
    - 1 BTC @ 30.000€ → ganancia 30.000€
    - 0.5 BTC @ 40.000€ → ganancia 10.000€
    - 0.5 BTC ??? → Pool vacío, error
  
Solución:
  Buscar dónde vino el 0.5 BTC faltante
  - ¿Fue un depósito fiat directo a exchange?
  - ¿Fue transferencia desde otro wallet (heredado, regalo)?
  - ¿Fue airdrop/staking?
  
Una vez identificado:
  Crear entrada manual en CoinTracking:
    Tipo: Buy (o Deposit, si es fiat)
    Cantidad: 0.5 BTC
    Precio: precio del día que recibiste
    Fecha: fecha exacta
```

---

## Validación en CoinTracking

```
CoinTracking → Reports → Gains:
  Para cada venta:
    1. ¿Se calcula ganancia sin error?
       SÍ → OK
       NO → Pool incompleto
    
    2. ¿La ganancia es sensata?
       Esperado: (Precio venta - Costo promedio) × cantidad
       
    3. ¿Hay "All pools consumed"?
       SÍ → Hay BTC sin origen documentado
       
Acción si hay error:
  → Ir a "Purchase History" de ese activo
  → Verificar que todas las compras estén listadas
  → Si falta, añadir manualmente
```

---

## Caso Especial: Depósito Inicial (Fiat → Cripto)

```
Situación común: "Comencé con 1 BTC que heredé"

¿Cómo registrarlo en CoinTracking?
  
Opción 1: Buy ficticio (mejor para auditoría)
  - Fecha: cuando recibiste (aunque fue regalo)
  - Tipo: Buy
  - Cantidad: 1 BTC
  - Precio: precio de mercado en esa fecha
  - Exchange: "Herencia" o "Regalo" (personalizados)
  - Resultado: base de coste clara
  
Opción 2: Deposit genérico (menos preciso)
  - Tipo: Deposit
  - Cantidad: 1 BTC
  - Resultado: no hay base de coste clara → problemas en ganancia
  
RECOMENDACIÓN: Opción 1 (buy ficticio con precio real del día)
```

---

## Integración

- **ADR-003:** Modelo de transacciones — FIFO es obligatorio en España
- **CAPITAL_GAINS.md:** Cálculo de ganancias basado en pool
- **COST_BASIS_AND_VALIDATION.md:** Validación de base de coste
