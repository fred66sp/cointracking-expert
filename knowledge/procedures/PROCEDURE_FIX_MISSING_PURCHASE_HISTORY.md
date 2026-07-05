---
id: KB-C3-002
title: "Procedimiento: Resolver Missing Purchase History (origen de activo)"
level: C
domain: cointracking
source: "Caso CT-002 + ADR-004"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-004
  - ADR-009

related_docs:
  - CT-002-venta-sin-historial-de-compra-previo-missing-purchase-history.md

tags:
  - procedure
  - missing-purchase-history
  - step-by-step

notes: "Paso a paso para resolver operaciones sin historial de compra."
---




# Procedimiento: Resolver Missing Purchase History

## Síntoma

CoinTracking warning: "No hay una compra adecuada para esta venta"

O visualizar: Ganancia = 0€ (Cost base = 0)

---

## Paso 1: Identificar el Activo y la Venta

```
CoinTracking → Reports → Gains
  ¿Encuentra warnings?
    SÍ → Anotar: Activo, Fecha, Cantidad
    NO → Continuar
```

---

## Paso 2: Buscar Compras Previas en CoinTracking

```
CoinTracking → Transacciones
  Filtrar: Tipo = "Buy" + Activo + Fecha < venta
  
¿Existe al menos una compra?
  SÍ → Paso 4 (problema cronológico)
  NO → Paso 3 (historial incompleto)
```

---

## Paso 3A: Importar Historial Faltante (Caso común)

Si no hay compra registrada en CoinTracking:

```
¿Dónde compró el usuario ese activo originalmente?
  a) En el mismo exchange (pero antes de la fecha actual de importación)
  b) En otro exchange (Coinbase, Kraken, etc.)
  c) Regalo/Airdrop/Mining (no hubo compra)

Caso a) Mismo exchange, historial anterior:
  → Reimportar desde fecha anterior en CoinTracking
  → Cuentas → Editar → Rango de fecha: mover inicio a 1 año antes

Caso b) Otro exchange:
  → Importar ese exchange en CoinTracking
  → Asegurar que el rango cubre la compra original

Caso c) Sin compra:
  → Ir a Paso 5 (crear entrada manual)
```

Después de reimportar:
1. Esperar a que CoinTracking procese (2-5 min)
2. Regenerar Tax Report
3. Verificar que el warning desaparece

---

## Paso 3B: Si Sigue Faltando (poco común)

```
¿El exchange real tiene registro de la compra original?
  NO → El usuario nunca la registró. Ir a Paso 5.
  SÍ → Exportar CSV de esa compra desde el exchange.
       → Importar CSV manualmente en CoinTracking.
```

---

## Paso 4: Corregir Zona Horaria (si hay compra pero aparece DESPUÉS)

```
CoinTracking → Settings → Timezone
  ¿Está en Europe/Madrid (con DST)?
    NO → Cambiar a Europe/Madrid
         → Reimportar (reiniciar importación API/CSV)
    SÍ → Las fechas en CoinTracking son correctas.
         Problema está en el exchange original.
         Verificar Tx Hash para confirmar fecha real.
```

---

## Paso 5: Crear Entrada Manual (último recurso)

Si tras importar sigue sin haber compra:

**IMPORTANTE:** Solo si existe evidencia documental (captura, Tx Hash, etc.)

```
1. Obtener evidencia:
   - Extracto del exchange (pantalla/PDF)
   - Tx Hash (si es on-chain)
   - Confirmación de email
   
2. Crear entrada en CoinTracking:
   - Tipo: Buy
   - Fecha/Hora: la del documento
   - Cantidad: la real
   - Precio: el que se pagó (o estimar si no está claro)
   - Comisión: incluir si aplica
   - Moneda: la del exchange origen
   - Intercambio: donde se compró
   
3. Regenerar Tax Report
4. Verificar warning desaparece
5. Documentar en REGISTRO-CAMBIOS qué se creó manualmente y por qué
```

---

## Paso 6: Verif ication Final

```
CoinTracking → Reports → Gains
  ¿El warning de "No hay compra" desaparece?
    SÍ → ✅ Resuelto
    NO → Repetir Paso 2-5 para esa venta específica
```

---

## Checklist

- [ ] Identificar activo, fecha, cantidad de la venta
- [ ] Buscar compras previas en CoinTracking
- [ ] Si falta: importar historial o exchange anterior
- [ ] Si aún falta: verificar zona horaria
- [ ] Si aún falta: crear entrada manual (con evidencia)
- [ ] Regenerar Tax Report
- [ ] Verificar warning desaparece
- [ ] Documentar cambios

---

## Integración

- **ADR-004:** Verificar contra exchange real
- **ADR-009:** No inventar origen sin evidencia
