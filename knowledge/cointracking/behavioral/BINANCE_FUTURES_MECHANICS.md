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
id: KB-B2-003
title: "Cómo CoinTracking maneja Binance Futures (contratos perpetuos)"
level: B
domain: cointracking
source: "Análisis + casos especializados"
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
  - futures
  - behavioral
  - complex

notes: "Operativo: cómo CoinTracking trata Futures; fiscalidad muy compleja."
---

# Cómo CoinTracking Maneja Binance Futures

## Definición

**Binance Futures** = Contratos perpetuos (perps): apuestas sobre el precio futuro, sin expiración.

**Equivalente:** Apuesta "el precio de BTC bajará a 40.000€" sin límite de tiempo.

---

## CRÍTICA FISCAL

⚠️ **Este es el tipo de operación más complejo fiscalmente en España.**

```
Razones:
  1. No hay compra/venta REAL de cripto (es un contrato)
  2. El PnL (ganancia/pérdida) es el "efectivo"
  3. Hacienda no tiene criterio claro (como de 2026)
  4. Podría ser "inversión especulativa" o "actividad económica"
```

---

## Cómo CoinTracking Registra Futures

### Flujo

```
1. Depositas USDT en la cuenta Futures
   Tipo: "Deposit" (a Futures)
   
2. Abres posición Long o Short
   Ejemplo: Long 1 BTC @ 50.000€
   
   CoinTracking ve:
     - Apalancamiento (si existe)
     - Collateral (garantía)
     - Entry price
   
3. Cierre de posición
   Exit price: 60.000€
   Ganancia: 60.000€ - 50.000€ = 10.000€ (en USDT)
   
   CoinTracking la registra como "PnL" o "Income"
   
4. Funding fee (cada 8h)
   Pequeñas cantidades de USDT pagadas/recibidas
   
   CoinTracking las ve como transacciones separadas
```

---

## Problemas Comunes

```
❌ CoinTracking no tiene tipo específico para "Futures"
   → Las operaciones aparecen como Income/Expense
   
❌ Los funding fees aparecen como "Transfer" confuso
   → Parecen depósitos/retiros inesperados
   
❌ El PnL se registra como "ganancia" inmediata
   → No diferencia entre realizado y unrealized
```

---

## Validación en CoinTracking

```
CoinTracking → Transacciones:
  Filtra por "Futures" o "USDTM"
  
¿Aparecen operaciones claras de Open/Close?
    NO → Futures está mezclado con Spot
    → Revisar manualmente
    
¿Los PnL aparecen como "Income" o "Expense"?
    SÍ → OK (aunque no es 100% exacto)
    
¿Los funding fees están separados?
    SÍ → Recopilar y documentar
    NO → Buscar en Binance directamente
```

---

## Tratamiento Fiscal (España, IRPF) — ⚠️ COMPLICADO

**Situación actual (2026): NO HAY CRITERIO OFICIAL de la DGT para Futures**

```
Lo que sabemos:
  - NO es compra/venta de cripto (no hay transferencia)
  - ES un "instrumento financiero derivado"
  - La ganancia es REAL (dinero USDT ganado)
  
Posibilidad 1: Ganancia patrimonial
  - Se trata como si vendieras cripto
  - Impuesto normal 19-45% según tramo
  
Posibilidad 2: Rendimiento del capital
  - Se trata como "interés de inversión"
  - Se suma a rendimientos
  
Posibilidad 3: Actividad económica
  - Si haces Futures frecuentemente
  - Hacienda podría considerarlo "profesión"
  - Aplican reglas de autónomo/empresa
```

**EJEMPLO (escenario):**
```
Ganas 10.000€ en Futures en 2025
Hacienda te pregunta: "¿Qué es esto?"

Respuesta segura:
  - "Ganancia patrimonial por derivados"
  - Pago impuesto como ganancia patrimonial (19%)
  - Total: ~1.900€
  
Respuesta arriesgada:
  - "Es solo especulación, no es nada"
  - Hacienda audita y te pide que justifiques
  - Te cobran atrasos + intereses + multa
```

---

## Recomendación

```
🔴 SI HACES TRADING DE FUTURES:
  
  1. DOCUMENTA TODA OPERACIÓN
     - Fecha, hora, entry, exit, PnL
     - Collateral, apalancamiento
  
  2. CONSULTA A UN ASESOR FISCAL
     - Especializado en criptos
     - ANTES de declarar, no después
  
  3. CONSIDERA DECLARAR COMO ACTIVIDAD ECONÓMICA
     - Si haces >50 operaciones/año
     - Mejor declaración + seguridad legal
  
  4. EVITA OCULTAR OPERACIONES
     - Binance informa a Hacienda (plataformas de cripto)
     - El riesgo de auditoría es alto
```

---

## Integración

- **ADR-003:** Modelo de transacciones (limitado para Futures)
- **CAPITAL_GAINS.md:** No aplica 100% (no es compra/venta)
- 🔴 **REQUIERE CONSULTA FISCAL ESPECIALIZADA**
