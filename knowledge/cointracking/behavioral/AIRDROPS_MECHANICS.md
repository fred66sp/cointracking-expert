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
id: KB-B1-002
title: "Cómo CoinTracking maneja Airdrops (regalos de tokens)"
level: B
domain: cointracking
source: "Casos reales CT-010 + análisis"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: medium
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/taxation/spain/CAPITAL_INCOME.md
  - CT-010-airdrop-registrado-como-compra-con-coste-artificial.md

tags:
  - cointracking
  - airdrops
  - behavioral
  - mechanical

notes: "Operativo: cómo CoinTracking importa y clasifica airdrops; errores comunes."
---

# Cómo CoinTracking Maneja Airdrops

## Definición

**Airdrop** = Distribución gratuita de tokens a holders de otro activo o dirección.

**Ejemplo:** Recibir 100 TOKEN gratis por tener 1 ETH en tu wallet el 1 de enero.

---

## Problema Común: CT-010

**Síntoma:** Airdrop registrado como "Compra" con coste = 0€.

**Causa:** CoinTracking no sabe cómo clasificar el airdrop en la importación.

**Impacto:**
- Puede distorsionar ganancias (si aparece como compra con cost=0 y luego se vende)
- Puede ser confuso en el Tax Report

---

## Cómo CoinTracking Registra Airdrops

### Importación automática (mejor caso)

```
Si importas vía API / CSV y el airdrop está documentado:
  CoinTracking lo clasifica como "Airdrop" (tipo específico)
  
Resultado:
  - Tipo = "Airdrop"
  - Cost base = 0 (fue regalo)
  - Fecha = día del airdrop
  - Cantidad = tokens recibidos
```

### Importación manual o CSV deficiente (peor caso)

```
Si el CSV no especifica que fue airdrop:
  CoinTracking puede clasificar como:
    a) "Deposit" → OK (aunque no sea 100% exacto)
    b) "Buy" → MALO (implica que pagaste por ellos)
    c) "Income" → Aceptable (tratado como ingreso)
```

---

## Validación y Corrección

### Identificar airdrops mal clasificados

```
CoinTracking → Transacciones:
  Filtrar: Tipo = "Buy" + Cantidad pequeña + Sin descripción de precio
  
¿Viste entradas con precio = 0 y descripción vaga?
    SÍ → Probablemente es un airdrop mal clasificado
```

### Corregir

```
Editar operación:
  Tipo: Cambiar de "Buy" a "Airdrop" (si existe)
  Si no existe "Airdrop", cambiar a "Deposit"
  
Precio/Cost: Dejar en 0 (fue regalo)
Comisión: 0 (no hubo comisión)
```

---

## Tratamiento Fiscal (España, IRPF)

**Airdrops = Rendimiento del capital en el momento de recepción**

```
Regla (DGT):
  - Momento exigible: fecha del airdrop (cuando se acredita en tu wallet)
  - Valuación: precio de mercado del token EN ESE MOMENTO
  - Impacto: se suma a "Rendimientos del capital" en IRPF
  - No es ganancia patrimonial (eso es venta posterior)
```

**Ejemplo:**
```
1 enero 2025: Airdrop de 100 TOKEN (precio ese día: 10€/token)
  → IRPF 2025: Rendimiento = 100 × 10€ = 1.000€

15 junio 2025: Vendes 100 TOKEN a 50€/token
  → Ganancia patrimonial = (50€ - 10€) × 100 = 4.000€
  → IRPF 2025: Rendimiento (1.000€) + Ganancia (4.000€) = 5.000€ total
```

---

## Validación en CoinTracking

```
Reports → Gains:
  ¿Los airdrops aparecen como "Income" o "Airdrop"?
    SÍ → OK
    NO → Verificar y corregir tipos
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Airdrop es tipo específico
- **CAPITAL_INCOME.md:** Tratamiento fiscal como rendimiento del capital
