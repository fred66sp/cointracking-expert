---
id: KB-B2-002
title: "Cómo CoinTracking maneja Binance Margin (trading apalancado)"
level: B
domain: cointracking
source: "Análisis + casos complejos"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-06-30
confidence: medium
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/taxation/spain/CAPITAL_GAINS.md

tags:
  - cointracking
  - binance
  - margin
  - behavioral
  - complex

notes: "Operativo: cómo CoinTracking trata margin trading; fiscalidad compleja."
---




# Cómo CoinTracking Maneja Binance Margin

## Definición

**Binance Margin** = Trading apalancado: pedir dinero prestado para multiplicar posición.

**Equivalente:** Comprar 2 BTC con solo 1 BTC (1x apalancamiento).

---

## Problema: Complejidad de Registros

CoinTracking ve múltiples operaciones:

```
1. Depósito inicial (fiat → Margin)
2. Borrow (préstamo de cripto)
3. Buy (compra con dinero prestado)
4. Sell (venta para cerrar posición)
5. Repay (devolución del préstamo)
6. Fee (interés del préstamo)
7. PnL (ganancia/pérdida realizada)
```

**Síntoma:** CoinTracking las ve como 7 operaciones separadas.

**Impacto fiscal:** Confuso; hay que relacionarlas correctamente.

---

## Cómo CoinTracking Registra Margin

### Flujo

```
1. Depositas 1 BTC en Margin
   Tipo: "Deposit" (a Margin)
   
2. Pides prestado: 1 BTC adicional
   CoinTracking lo ve como "Borrow" (si está disponible)
   O como una entrada confusa
   
3. Compras 2 BTC total (con el dinero prestado)
   Tipo: "Buy" (2 BTC a precio X)
   
4. Vendes 2 BTC cuando baja el precio
   Tipo: "Sell" (2 BTC a precio Y)
   
5. Devuelves 1 BTC prestado
   Tipo: "Repay" o "Transfer"
   
6. Interés del préstamo
   Tipo: "Fee" o "Expense"
   
7. Tu ganancia neta
   = (Precio venta - Precio compra) - Interés
```

---

## Validación en CoinTracking

```
Reports → Margin Account:
  ¿Se muestra el saldo de préstamo?
    SÍ → OK
    NO → Verificar que "Borrow" esté registrado

¿Los intereses aparecen como "Expense" o "Fee"?
    SÍ → OK (se restan de la ganancia)
    NO → Añadir manualmente
    
¿La ganancia neta es sensata?
  Esperado = (Venta - Compra) - Interés
```

---

## Tratamiento Fiscal (España, IRPF)

**CRÍTICO: Margin = Más complejo que Spot**

```
Regla (DGT — INTERPRETACIÓN TÉCNICA):
  
1. La operación de compra es REAL (tienes los BTC)
2. La operación de venta es REAL (vendes los BTC)
3. El interés del préstamo es GASTO (reduce la ganancia)
4. La ganancia se calcula: (Venta - Compra) - Interés

Ejemplo:
  - Compra: 2 BTC a 50.000€/BTC = 100.000€ costo
  - Venta: 2 BTC a 60.000€/BTC = 120.000€
  - Interés: 1.000€
  - Ganancia: 120.000€ - 100.000€ - 1.000€ = 19.000€
```

**ADVERTENCIA:** El tratamiento de margin puede variar según:
- Si consideras que es "inversión" o "actividad económica"
- Si hacienda considera los intereses como gasto deducible
- El número de operaciones (frecuencia = actividad económica)

→ **CONSULTA A UN ASESOR FISCAL ESPECIALIZADO** si haces margin trading intenso.

---

## Integración

- **ADR-003:** Modelo de transacciones
- **CAPITAL_GAINS.md:** Cálculo de ganancias (con intereses)
- ⚠️ **Solicitar revisión fiscal para casos complejos**
