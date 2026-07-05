---
id: KB-B1-019
title: "Cómo CoinTracking maneja Comisiones (Fees) en diferentes monedas"
level: B
domain: cointracking
source: "Casos reales de comisiones complejas"
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
  - knowledge/blockchains/GAS_FEE_HANDLING.md
  - knowledge/cointracking/behavioral/BINANCE_SPOT_MECHANICS.md

tags:
  - cointracking
  - fees
  - commissions
  - behavioral

notes: "Operativo: cómo registrar fees cuando están en moneda diferente."
---




# Cómo CoinTracking Maneja Comisiones (Fees)

## Tipos de Comisiones

### Tipo 1: Comisión en la Misma Moneda

**Más simple: La comisión se resta del activo negociado**

```
Ejemplo: Venta de 100 USDC

Binance:
  Vendo: 100 USDC
  Comisión: 0.1 USDC (0.1%)
  Recibo neto: 99.9 USDC
  
CoinTracking:
  Tipo: Sell
  Cantidad: 100 USDC
  Precio: 1 USDT/USDC (sin cambio)
  Fee: 0.1 USDC
  
¿Cómo aparece?
  Opción A (automático): 
    "Sell 100 USDC" (el fee se incluye en el cálculo)
  Opción B (manual):
    "Sell 100 USDC, Fee 0.1 USDC"
```

---

### Tipo 2: Comisión en Moneda Diferente

**Complicado: La comisión está en otro token**

```
Ejemplo: Compra de 1 BTC pagada en USDT

Binance:
  Compro: 1 BTC
  Precio: 50.000 USDT
  Comisión: 25 USDT (0.05%, paga comisión reducida con BNB)
  Total pagado: 50.025 USDT
  
¿Cómo aparece en CoinTracking?

Ideal (automático):
  TX1: Buy 1 BTC @ 50.025 USDT (fee incluido)
  
Realidad (frecuente):
  TX1: Buy 1 BTC @ 50.000 USDT
  TX2: Fee 25 USDT (operación separada)
  
Problema:
  - CoinTracking ve dos operaciones separadas
  - Cost basis puede quedar confuso
  - O no incluye el fee en el costo total
```

---

### Tipo 3: Comisión en BNB (Binance Coin)

**Caso especial: Descuento de Binance por pagar en BNB**

```
Beneficio:
  Binance: Si pagas comisión en BNB, es 25% más barata
  
Ejemplo:
  Compro 1 BTC @ 50.000 USDT
  
  Comisión si pago en USDT: 25 USDT
  Comisión si pago en BNB: 18.75 BNB (equivalente a 18.75 USDT)
  
¿Cómo aparece?

Ideal:
  TX: Buy 1 BTC @ 50.000 USDT, Fee 18.75 BNB
  
Realidad:
  CoinTracking puede ver:
    TX1: Buy 1 BTC @ 50.000 USDT (no ve el fee)
    TX2: BNB sent -18.75 (parece pérdida de BNB)
  
Problema:
  - La pérdida de BNB no se conecta con la compra de BTC
  - Parece que "perdiste BNB"
  - Cost basis de BTC está incompleto
```

---

## Problema: Cost Basis Incorrecto por Fee No Incluido

```
Escenario:
  Compro 1 BTC @ 50.000 USDT
  Comisión: 50 USDT (pagada en BNB)
  
¿Cuál es el cost basis correcto de 1 BTC?

Respuesta:
  Cost basis = 50.000 USDT + 50 USDT = 50.050 USDT
  (La comisión es parte del costo total)
  
¿Qué pasa si CoinTracking registra solo 50.000 USDT?
  
  Vendo después a 60.000 USDT
  
  Con cost correcto:
    Ganancia = 60.000 - 50.050 = 9.950 USDT
    
  Con cost incorrecto:
    Ganancia = 60.000 - 50.000 = 10.000 USDT (FALSO, 50 USDT demasiado)
```

---

## Validación en CoinTracking

### Detectar Fees Faltantes

```
CoinTracking → Transacciones:
  
Para CADA compra:
  [ ] ¿Hay una línea de "Fee" asociada?
      SÍ → OK
      NO → Verificar si el fee fue en otra moneda
          
  [ ] Si el fee está en otra moneda (BNB, etc.):
      ¿Aparece como "Transfer out" o "Expense"?
      SÍ → OK (aunque no conectado visualmente)
      NO → Fee faltante, necesita añadirse
      
Para CADA venta:
  [ ] ¿El precio incluye comisión o es neto?
      Verificar contra Binance:
        Binance: "Total received" (ya con fee deducido)
        vs
        CoinTracking: precio que registraste
      ¿Iguales?
      SÍ → OK
      NO → Ajustar el precio o añadir fee
```

### Verificar Fee vs Binance

```
Binance → Historial:
  Busca la transacción
  Anota: "Comisión" o "Fee"
  
Compara con CoinTracking:
  ¿El fee en CT coincide con Binance?
  
¿Y la moneda de la comisión?
  Binance: "Comisión en BNB"
  CoinTracking: ¿Tiene entrada de "BNB sent"?
  
Mismatch → Ajustar en CT
```

---

## Registrar Fees Correctamente

### Opción 1: Incluir Fee en el Precio (Mejor para Auditoría)

```
Compro 1 BTC @ 50.000 USDT con fee de 50 USDT en BNB

En CoinTracking:
  Tipo: Buy
  Cantidad: 1 BTC
  Precio: 50.050 USDT/BTC (INCLUYE fee)
  
Resultado:
  - Cost basis: 50.050 USDT
  - Una sola línea (más limpio)
  - Ganancia calculada correcta

Nota: El fee en BNB aparecerá como "BNB spent" separado
  Pero eso es OK, audita BNB por separado
```

### Opción 2: Registrar Fee como Línea Separada

```
Compro 1 BTC @ 50.000 USDT con fee de 50 USDT en BNB

En CoinTracking:
  TX1: Buy 1 BTC @ 50.000 USDT
  TX2: Expense 50 USDT (o 50 BNB convertido)
  
Problema:
  - Cost basis de BTC: 50.000 USDT (FALSO, sin fee)
  - Fee aparece separado (menos claro)
  
Solución:
  Editar TX1 después:
    Añadir "Commission/Fee": 50 USDT
    CoinTracking automáticamente ajusta cost basis
```

### Opción 3: BNB Comisión (Especial)

```
Si la comisión fue en BNB:
  
CoinTracking → Add Transaction:
  Tipo: Expense (o Fee)
  Activo: BNB
  Cantidad: 18.75 BNB
  Precio: Precio BNB ese día (ej. 200€/BNB = 3.750€)
  Descripción: "Trading fee for BTC buy"
  
Resultado:
  - Se registra como gasto de BNB
  - Se resta del balance de BNB
  - Aparece en auditoría
```

---

## Tratamiento Fiscal (España, IRPF)

**Fees = Gasto de operación (teóricamente deducible)**

```
Regla lógica:
  - Fee es dinero que gastas para operar
  - Debería reducir la ganancia
  
Ejemplo:
  Vendo 1 BTC a 60.000€ (compraste a 50.000€)
  Fee: 50€
  
  Ganancia nominal: 10.000€
  Ganancia neta: 10.000€ - 50€ = 9.950€
  
Realidad fiscal (España):
  - NO hay criterio oficial de DGT
  - Algunos asesores lo permiten, otros no
  - Enfoque conservador: NO deducir fee
  
RECOMENDACIÓN:
  - Documentar TODOS los fees
  - Incluirlos en cost basis (opción 1 arriba)
  - El asesor fiscal verifica después
```

---

## Checklist: Auditar Fees

```
Antes de auditar:

[ ] ¿Todos los fees están registrados?
    Comparar CoinTracking vs Binance historial
    
[ ] ¿Los fees están en la moneda correcta?
    Ej. Si fue en BNB, ¿aparece BNB sent?
    
[ ] ¿El cost basis incluye fees?
    [ ] Check: precio en CT = (total pagado / cantidad)
    
[ ] ¿Hay fees "huérfanos" (sin transacción)?
    Ver si hay BNB/comisión sin transacción asociada
    
[ ] ¿Se pueden deducir fiscalmente?
    Documentar para consulta con asesor
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Fees son parte del cost basis
- **GAS_FEE_HANDLING.md:** Fees en blockchain (gas)
- **BINANCE_SPOT_MECHANICS.md:** Fees en Binance Spot
