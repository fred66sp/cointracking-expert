---
id: KB-B1-009
title: "Duplicados por Importar API y CSV Simultáneamente (Overlap)"
level: B
domain: cointracking
source: "Casos reales de reimportación"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-014
  - ADR-032

related_docs:
  - knowledge/cointracking/behavioral/DUPLICATE_DETECTION_HEURISTICS.md
  - knowledge/patterns/PATTERN_DUPLICATE_DETECTION.md

tags:
  - cointracking
  - api
  - csv
  - duplicates
  - overlap
  - behavioral

notes: "Crítico: cómo evitar/limpiar duplicados de fuentes mixtas."
---

# Duplicados por Importar API y CSV Simultáneamente (Overlap)

## Problema: Mezclar Fuentes

**Síntoma más común en auditoría:**

```
Usuario:
  - Descargó CSV de Binance (histórico 2024)
  - Conectó API de Binance (desde hoy)
  - Importó ambos en CoinTracking
  
Resultado:
  - Operaciones de 2024 aparecen DUPLICADAS
  - CSV: 100 USDC (source: "Binance CSV")
  - API: 100 USDC (source: "Binance API")
  - Mismo timestamp, mismo importe
  
Balance:
  - USDC en CT: 200 (¡FALSO!)
  - USDC en realidad: 100
```

---

## ¿Cuándo Sucede el Overlap?

### Caso 1: Reimportar CSV Histórico con API Activa

```
Timeline:
  Enero-Junio 2024: Operaciones en Binance
  15 julio 2024: Usuario conecta API a CoinTracking
  20 julio 2024: Usuario descarga CSV (enero-junio)
  21 julio 2024: Usuario importa CSV en CT
  
Resultado:
  - API trajo datos enero-junio (histórico + hoy)
  - CSV también trae enero-junio
  - OVERLAP: enero-junio duplicados
```

### Caso 2: Importar CSV Dos Veces Accidentalmente

```
Timeline:
  1 febrero: Importa CSV de Binance (enero-febrero)
  15 febrero: Importa MISMO CSV again (olvida que ya lo hizo)
  
Resultado:
  - Todas las operaciones duplicadas
  - Balance es 2x
```

### Caso 3: Cambiar de Source y Reimportar

```
Timeline:
  Enero: Importa CSV de Binance
  Marzo: Cambia a API de Binance (API es más preciso)
  Abril: API trae historial desde enero
  
Resultado:
  - Enero-marzo: Duplicados
  - Abril+: Solo API (OK)
```

---

## Detección de Overlap

### Indicadores de Overlap

```
¿Cómo saber si tienes overlap?

Signo 1: Balance duplicado
  Balance en CT: 2 BTC
  Balance real (wallet): 1 BTC
  → Probable overlap

Signo 2: Report duplicado
  CoinTracking → Reports → Gains
  Ves MISMA venta DOS VECES
  → Overlap confirmado

Signo 3: Source mixto
  CoinTracking → Transacciones
  Filtra por Binance
  ¿Ves dos filas idénticas con source diferente?
    - TX1: "Binance API"
    - TX2: "Binance CSV"
  → Overlap confirmado
```

### Script de Detección

```
CoinTracking → Tools → Duplicates:
  Automáticamente busca "posibles duplicados"
  
¿Qué ve?
  - Mismo timestamp
  - Mismo activo
  - Misma cantidad
  
PERO: No siempre detecta overlap
  (si hay pequeñas diferencias en timestamp)
```

---

## Limpieza de Overlap

### Opción 1: Eliminar CSV (Mejor si Tienes API)

```
Razonamiento:
  - API es más preciso (datos en tiempo real)
  - API tiene más metadata (fee, hash, etc.)
  - CSV es estático y propenso a errores
  
Pasos:
  1. Conectar API de Binance (si no lo hizo)
  2. Reimportar histórico vía API
  3. Eliminar TODAS las operaciones importadas vía CSV
  4. Verificar que no falte nada
  
Verificación:
  [ ] ¿Balance es ahora correcto?
  [ ] ¿Ganancias coinciden con expectations?
  [ ] ¿Hay operaciones "raras" sin source API?
      SÍ → Completar manualmente
```

### Opción 2: Eliminar API (Si Prefieres CSV)

```
Razonamiento:
  - CSV es "snapshot" (no cambia)
  - API puede tener sincronización fallida
  - Prefiero lo que tengo documentado
  
Pasos:
  1. Desconectar API de Binance
  2. Eliminar TODAS las operaciones importadas vía API
  3. Mantener CSV (en rango específico)
  4. Añadir manualmente lo que falte después del CSV
  
Desventaja:
  - Requiere más trabajo (manual para reciente)
  - Más propenso a errores
```

### Opción 3: Merged Manual (Complejo, No Recomendado)

```
¿Idea?
  Resolver manualmente cada duplicado
  Guardar "la mejor versión"
  
Problema:
  - Muy tedioso (100+ operaciones)
  - Propenso a errores
  
Solo usar si:
  - Overlap es pequeño (<10 operaciones)
  - Las operaciones son complejas (necesitas datos de ambas)
```

---

## Prevención: Mejores Prácticas

### Estrategia 1: API Primaria (Recomendado)

```
Flujo correcto:

Paso 1: Conectar API lo antes posible
  CoinTracking → Settings → Exchanges → Add Binance API
  
Paso 2: Importar histórico completo vía API
  CoinTracking detecta automáticamente todo
  
Paso 3: NUNCA importar CSV de Binance (ya está todo)
  
Resultado:
  - Una fuente única (API)
  - Sin overlaps
  - Datos precisos
```

### Estrategia 2: CSV Única por Rango

```
Si prefieres CSV (sin API):

Paso 1: Descarga CSV una sola vez
  Binance → Tu historial de criptos → Download CSV
  Especificar rango: enero 2024 - diciembre 2024
  
Paso 2: Importa en CoinTracking
  UNA sola vez
  
Paso 3: Para nuevas operaciones (2025)
  Opción A: Descargar nuevo CSV (2025) e importar (sin overlapping)
  Opción B: Conectar API para nuevo período
  
Nunca reimportar CSV antiguo
```

### Estrategia 3: API + CSV Complementaria

```
Híbrida (usado si hay problemas con API):

Paso 1: Conectar API de Binance (primaria)
  Importa todo lo que puede
  
Paso 2: Si algo falta, descargar CSV específico
  CSV de solo lo que falta (ej. un mes)
  Importar en rango diferente
  
Paso 3: Verificar no hay overlap
  Antes de importar CSV, verificar fecha/operaciones
  Asegurarse que no están ya en API
```

---

## Limpiar Overlap: Checklist

```
Paso 1: Diagnosis
  [ ] Confirmar que hay overlap (comparar con realidad)
  [ ] Contar cuántas operaciones duplicadas hay
  
Paso 2: Decidir estrategia
  [ ] ¿Tengo API funcionando?
      SÍ → Opción 1 (eliminar CSV)
      NO → Opción 2 (mantener CSV)
      
Paso 3: Limpiar
  CoinTracking → Transacciones
  Filtrar por source problemática
  [ ] Seleccionar todas las duplicadas
  [ ] Delete (si estoy seguro)
  [ ] O Edit si necesito mantener ciertas operaciones
  
Paso 4: Verificar
  [ ] Balance es ahora correcto?
  [ ] Ganancias calculadas son razonables?
  [ ] Timestamps son correctos?
  [ ] Regenerar Reports
  
Paso 5: Documentar
  [ ] Anotar en audit log qué se limpió y por qué
```

---

## Caso Especial: Partial Overlap

**¿Qué si solo PARTE está duplicada?**

```
Situación:
  CSV (enero-julio 2024): 100 operaciones
  API (marzo 2024 onwards): 80 operaciones
  
Overlap: Marzo-julio (60 operaciones)
No overlap: Enero-febrero (40 operaciones de CSV)

¿Cómo limpiar?

Opción A: Eliminar solo el CSV (más fácil)
  - Elimina CSV completo
  - Verifica que API trajo enero-febrero
  - Si no, añade manualmente

Opción B: Eliminar solo API duplicada
  - Mantener CSV
  - Eliminar API de marzo-julio
  - Mantener API de agosto+ (solo en API)
```

---

## Integración

- **ADR-014:** Protocolo de verificación antes de eliminar
- **DUPLICATE_DETECTION_HEURISTICS.md:** Cómo detectar duplicados
- **PATTERN_DUPLICATE_DETECTION.md:** Casos reales
