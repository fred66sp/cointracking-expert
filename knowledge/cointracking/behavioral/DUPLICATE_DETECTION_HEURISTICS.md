---
id: KB-B1-006
title: "Heurísticas de Detección de Duplicados en CoinTracking"
level: B
domain: cointracking
source: "Análisis de casos CT-001, CT-002 + patrones reales"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-03
confidence: medium
version: 1.0

related_adr:
  - ADR-014
  - ADR-032

related_docs:
  - knowledge/patterns/PATTERN_DUPLICATE_DETECTION.md
  - knowledge/cases/CT-001-duplicate-same-timestamp.md
  - knowledge/cases/CT-002-floki-batching.md

tags:
  - cointracking
  - duplicates
  - detection
  - behavioral
  - critical

notes: "Crítico para auditoría: cómo CoinTracking detecta duplicados automáticamente."
---




# Heurísticas de Detección de Duplicados en CoinTracking

## Definición: Duplicado

**Duplicado** = Misma transacción registrada DOS VECES en CoinTracking.

**Causa común:** Importación múltiple (API + CSV, o API dos veces).

```
Ejemplo:
  Venta de 100 FLOKI @ 0.10 USDT

  CoinTracking ve (si hay duplicado):
    TX 1: 100 FLOKI → 10 USDT
    TX 2: 100 FLOKI → 10 USDT
  
  Balance calculado:
    -200 FLOKI (¡FALSO! Vendiste solo 100)
    +20 USDT (¡FALSO! Recibiste solo 10)
```

---

## Heurística 1: Timestamp Idéntico

**Si dos transacciones tienen exactamente el MISMO timestamp = PROBABLE DUPLICADO**

```
Indicador fuerte:
  TX1: 2024-03-15 14:23:45 | 100 FLOKI | 10 USDT
  TX2: 2024-03-15 14:23:45 | 100 FLOKI | 10 USDT
  
¿Son duplicados?
  SÍ, muy probable (99%)
  
¿Cómo verificar?
  → Ver la fuente de ambas TX
  → Si son de la misma API/CSV → Duplicado
  → Si son de sources diferentes → Menos probable, pero aún sospechoso
  
Acción:
  → Eliminar una (la más reciente, que es el duplicado de reimportación)
```

---

## Heurística 2: Cantidad + Precio Idénticos (mismo segundo)

**Si dos transacciones tienen cantidad, precio Y fecha/hora en el MISMO SEGUNDO = MUY PROBABLE DUPLICADO**

```
Indicador muy fuerte:
  TX1: 2024-03-15 14:23:45 | 100 FLOKI @ 0.10 USDT
  TX2: 2024-03-15 14:23:45 | 100 FLOKI @ 0.10 USDT
  
¿Son duplicados?
  SÍ, casi seguro (95%+)
  
EXCEPTO (ADR-014):
  Si tienen Trade IDs diferentes en Binance
    → NO son duplicados, es Binance batching
```

---

## Heurística 3: Binance Batching (FALSO POSITIVO)

**CRÍTICO: No todos los duplicados aparentes son duplicados reales**

```
Problema (CT-002 — FLOKI):
  29 operaciones idénticas el mismo segundo:
    TX1: 2024-03-17 18:39:11 | 4570 FLOKI
    TX2: 2024-03-17 18:39:11 | 4570 FLOKI
    ...
    TX29: 2024-03-17 18:39:11 | 4570 FLOKI
  
  ¿Son duplicados?
    NO, son operaciones legítimas de Binance batching
  
  ¿Cómo distinguir?
    → Ver Trade ID en Binance API
    → Si Trade IDs son DIFERENTES → NO son duplicados
    → Si Trade IDs son IDÉNTICOS → SÍ son duplicados
  
  Verificación (MANDATORIA antes de eliminar):
    Binance → Herramientas → Historial
    Buscar fecha/operación
    Comparar Trade IDs:
      - FLOKI1: Trade ID = 100369243
      - FLOKI2: Trade ID = 100369244 (diferente)
      - ...
    
    Si todos distintos → SON LEGÍTIMOS, NO ELIMINAR
```

---

## Heurística 4: Exchange/Source Diferentes

**Si dos transacciones vienen de SOURCE DIFERENTE = PROBABLE DUPLICADO PARCIAL**

```
Ejemplo:
  TX1: 100 USDC | Source: Binance API
  TX2: 100 USDC | Source: Manual import
  
¿Son duplicados?
  Probable sí, pero no seguro
  
¿Cómo verificar?
  → ¿Cuando importé manualmente, también estaba en API?
  → ¿O fue una corrección que añadí después?
  
Acción:
  → Si duplicado: eliminar la más antigua o menos confiable
  → Si no duplicado: mantener ambas con etiquetas claras
```

---

## Heurística 5: Descripción Idéntica (pero Difícil de Usar)

**Si la descripción es 100% idéntica = PODRÍA ser duplicado, pero es débil**

```
Menos confiable porque:
  - Usuarios pueden copiar descripción manualmente
  - CoinTracking auto-genera descripciones iguales
  
Ejemplo (débil):
  TX1: "Binance USDC withdrawal"
  TX2: "Binance USDC withdrawal"
  
¿Son duplicados?
  Posible, pero no concluyente
  
Necesita más verificación (timestamp, cantidad, fecha)
```

---

## Matriz de Decisión

```
Timestamp | Cantidad | Precio | Trade ID | Conclusión
----------|----------|--------|----------|-------------------
Idéntico  | Idéntico | Idéntico | Distinto | NO duplicado (batching)
Idéntico  | Idéntico | Idéntico | Idéntico | DUPLICADO (eliminar)
Idéntico  | Idéntico | Idéntico | ???      | PROBABLE duplicado
Distinto  | Idéntico | Idéntico | -        | Poco probable dup
Idéntico  | Distinto | -        | -        | NO duplicado
Distinto  | Distinto | -        | -        | NO duplicado
```

---

## Validación en CoinTracking

```
CoinTracking tiene detección automática de duplicados:
  
  Transacciones → Herramientas → Duplicates
  
¿Qué hace?
  - Busca transacciones con timestamp/cantidad idénticos
  - Las marca como posibles duplicados
  - PERO: NO distingue entre batching legítimo y duplicados reales
  
¿Qué debes hacer?
  1. Revisar cada "posible duplicado" de la lista
  2. Verificar Trade ID en Binance (si aplica)
  3. Si Trade ID distinto → MANTENER (no es duplicado)
  4. Si Trade ID idéntico → ELIMINAR (es duplicado real)
  5. Si no hay Trade ID disponible → PREGUNTA al usuario antes de eliminar
```

---

## Checklist de Verificación (ADR-014)

```
Ante CADA duplicado detectado:

[ ] ¿El timestamp es idéntico?
[ ] ¿La cantidad es idéntica?
[ ] ¿El precio es idéntico?
[ ] ¿Vienen de la misma fuente? (API, CSV, manual)

Si SÍ a todo:
  [ ] BUSCAR EN BINANCE: ¿Trade IDs distintos?
      SÍ → NO ES DUPLICADO, MANTENER
      NO → ES DUPLICADO, ELIMINAR
  [ ] Si no hay Binance info → PEDIR CONFIRMACIÓN al usuario

[ ] Tras eliminación → Regenerar Tax Report
[ ] Verificar que balance sea consistente
```

---

## Caso Especial: Transacciones Parcialmente Duplicadas

```
¿Qué pasa si?
  - TX1 importada: 100 USDC
  - TX2 importada: 50 USDC (parcial)
  
¿Son duplicados?
  NO, son parciales (puede ser que TX1 se dividió)
  
Acción:
  - Investigar el origen (¿es una corrección?)
  - Si es error: eliminar el más dudoso
  - Si es real: mantener ambos con notas
```

---

## Integración

- **ADR-014:** Protocolo de verificación de Trade ID antes de eliminar
- **PATTERN_DUPLICATE_DETECTION.md:** Casos documentados
- **CT-001, CT-002:** Casos de referencia
