---
id: KB-C2-001
title: "Patrón: Detección de Duplicados (matriz de heurísticas)"
level: C
domain: cointracking
source: "Generalización de 20 casos reales + ADR-014, ADR-026"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-014
  - ADR-026
  - ADR-003

related_docs:
  - CT-002-venta-sin-historial-de-compra-previo-missing-purchase-history.md
  - CT-003-api-y-csv-importados-simultaneamente-duplicado-por-doble-fuente.md
  - CT-008-duplicados-aparentes-por-ejecucion-parcial-de-una-orden.md
  - CT-016-duplicados-por-reimportacion-completa-del-mismo-periodo.md

tags:
  - pattern
  - duplicates
  - cointracking
  - auditing

notes: "Patrón derivado de 4+ casos reales. Regla de oro: Trade ID distinto = legítimas."
---




# Patrón: Detección de Duplicados

## Heurísticas (Matriz)

**Antes de marcar algo como duplicado, verificar en orden:**

| Heurística | Presente | NO Presente | Conclusión |
|---|---|---|---|
| **Fecha+Hora idéntica** | ✓ | | Sospechoso |
| Precio idéntico | ✓ | | Muy sospechoso |
| Volumen idéntico | ✓ | | Muy sospechoso |
| **Trade ID distinto en Binance** | ✓ | ❌ | **LEGÍTIMAS** (no son duplicados) |
| Comisión idéntica | ✓ | | Sospechoso |
| **Misma fuente de importación** (API + API, CSV + CSV) | ✓ | | Probables duplicados |
| Fuentes distintas (API + CSV) | ✓ | | **DUPLICADOS REALES** |

---

## Caso de Referencia: CT-002 (FLOKI)

**Real:** 29 transacciones Buy FLOKI el 17.03.2024 18:39:11, todas con valores idénticos.

**Apariencia:** Parecían duplicadas (misma fecha, hora, precio, volumen).

**Realidad:** Cada una tenía **Trade ID distinto en Binance API**:
- FLOKIUSDT22086512
- FLOKIUSDT100369243
- ... (27 más)

**Conclusión:** **Legítimas, no duplicados** (Binance hizo batching de múltiples órdenes en el mismo segundo).

**Lección:** Trade ID es la fuente de verdad.

---

## Protocolo de Detección (Paso a Paso)

### Paso 1: Detectar candidatos

```
¿Misma fecha + hora?
  SÍ → Paso 2
  NO → No es duplicado probable
```

### Paso 2: Verificar precio + volumen

```
¿Misma fecha + precio + volumen?
  SÍ → Sospechoso. Paso 3
  NO → Probablemente NO sea duplicado
```

### Paso 3: CRÍTICO — Verificar Trade ID en exchange

```
Ir a Binance/Kraken/Coinbase (según exchange):
  Buscar Trade ID de ambas transacciones en historial
  
¿Trade ID distinto?
  SÍ → ✅ LEGÍTIMAS. NO ELIMINAR.
  NO (Trade ID igual o ambos vacíos) → Paso 4
```

### Paso 4: Verificar fuente de importación

```
¿Una vino de API y otra de CSV?
  SÍ → Casi seguro DUPLICADO (reimportación)
  NO (ambas de la misma fuente) → Paso 5
```

### Paso 5: Buscar ejecución parcial

```
¿Estos son parte de una orden que se ejecutó en múltiples tramos?
  SÍ → Probablemente LEGÍTIMAS (partes de la misma orden)
  NO → Probable duplicado
```

### Paso 6: Pedir confirmación antes de eliminar

```
Explicar al usuario:
  - Qué lo hace sospechoso (fecha/precio/volumen idénticos)
  - Qué lo hace potencialmente legítimo (Trade ID, fuente, lógica)
  
NUNCA eliminar sin confirmación explícita del usuario (ADR-026: Categoría B).
```

---

## Falsos Positivos Comunes

| Situación | Parece Duplicado | Es en Realidad |
|-----------|-----------------|---|
| Múltiples pequeñas órdenes en el mismo segundo (Binance batching) | ✓ | Legítimas (Trade IDs distintos) |
| Reimportación accidental (API + CSV mismo periodo) | ✓ | **Duplicado real** |
| Ejecución parcial de una orden grande (p. ej. 10 pequeñas compras = 1 orden grande) | ✓ | Legítimas (partes de lo mismo) |
| Swap en DeFi fragmentado en varias tx on-chain | ✓ | Legítimas (pasos del swap) |

---

## Regla de Oro

> **Si hay duda, el Trade ID en el exchange es la fuente de verdad.**
> 
> Trade ID distinto en el exchange = Operaciones legítimas, no duplicados.
> 
> Trade ID igual (o ambos vacíos) = Probable duplicado.

---

## Integración con ADRs

- **ADR-014:** Validación de duplicados con Trade ID → Este patrón lo operacionaliza
- **ADR-026:** Matriz A/B/C → Duplicados son Categoría B (requieren confirmación usuario)
- **ADR-003:** Modelo de transacciones → Tipos influyen en lo que es/no es duplicado

---

## Próximas Mejoras

- Automatizar búsqueda de Trade ID si MCP conectado
- Crear checklist de verificación (D1)
- Crear árbol de decisión (D2)
