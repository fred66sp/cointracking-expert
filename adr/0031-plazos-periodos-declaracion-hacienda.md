# ADR-031: Validación temporal previa de obligaciones fiscales

**Status:** Proposed

**Date:** 2026-07-05

## Context

El agente prepara declaraciones fiscales (IRPF, Modelo 721, etc.). Sin validar si es el momento correcto, podría:

1. Preparar una declaración para un ejercicio futuro (p. ej. 2027 cuando es 2024)
2. Avisar tarde al usuario que pasó el plazo de campaña ordinaria
3. No detectar que está en situación fiscal "especial" (fuera de plazo ordinario)

**Caso real:** Usuario llama 1 julio 2026. "Quiero hacer renta 2025". Plazo ordinario fue hasta 30.06.2026 (ayer). El agente debe **antes de auditar nada** detectar que la situación cambió.

## Decision

**Antes de proceder con auditoría o declaración, el agente ejecuta validación temporal** (Paso 0.5 de `/spanish-tax-return`):

### 1. Validar ejercicio

```
¿Ejercicio < Año actual? SÍ → válido (no futuro)
¿Ejercicio >= (Año actual - 4)? SÍ → válido (dentro de prescripción)
SINO → PARAR: Ejercicio fuera de rango
```

### 2. Determinar estado fiscal

```
Deadlines = Load("knowledge/taxation/spain/FILING_DEADLINES.md")
  [Con metadatos vigencia: ADR-032]

Estado = Evaluate(Hoy, Ejercicio, Deadlines[Ejercicio])

Resultado: ORDINARIO | LATE | UNKNOWN | FUTURE | EXPIRED
```

### 3. Máquina de estados

```
ORDINARIO
  → ✅ Campaña abierta, sin penalizaciones
  → Continuar con auditoría

LATE
  → ⚠️ Fuera de plazo ordinario (puede haber penalizaciones)
  → Avisar usuario: "¿Quieres continuar?"
  → Consultar con asesor fiscal (ADR-028: límite auditor/asesor)

FUTURE
  → ❌ Ejercicio aún no finalizado
  → No se puede declarar

UNKNOWN
  → ⚠️ [PENDIENTE DE VERIFICAR]
  → Deadlines no está vigente (ADR-032)
  → Reverificar contra `knowledge/taxation/spain/FILING_DEADLINES.md`

EXPIRED
  → ⚠️ Fuera de todo plazo conocido
  → Requiere revisión especial (recurso, etc.)
  → Remitir a asesor fiscal
```

### 4. Avisos al usuario

**ORDINARIO:**
```
✅ Plazo ordinario vigente.
Continuar sin penalizaciones.
```

**LATE:**
```
⚠️ Pasó el plazo ordinario.
Presenta con posibles consecuencias (consulta asesor).
¿Continuar?
```

**FUTURE:**
```
❌ El ejercicio [año] aún no finaliza (31.12.[año]).
Revienta en [fecha]. Intenta entonces.
```

**UNKNOWN:**
```
⚠️ [PENDIENTE DE VERIFICAR]
Los plazos están envejecidos o no son verificables.
Requiere reverificación manual.
```

**EXPIRED:**
```
❌ Fuera de todo plazo de presentación ordinario/voluntario.
Consulta asesor fiscal sobre opciones (recurso, etc.).
```

---

## Consequences

**Positive:**

- **Prevención:** El usuario sabe antes de empezar si la declaración es oportuna
- **Claridad:** Estados claros (ORDINARIO, LATE, FUTURE, etc.)
- **Responsabilidad:** Avisos explícitos, no sorpresas
- **Arquitectura pura:** Sin hardcoding de fechas (viven en `knowledge/`)
- **Reutilizable:** Mismo protocolo para cualquier obligación con plazo

**Negative:**

- **Requiere mantenimiento:** `knowledge/taxation/spain/FILING_DEADLINES.md` debe actualizarse anualmente
- **Ambigüedades legales:** "LATE" es estado técnico, no legal (el asesor decide si se puede presentar)

---

## Notes

### Relación con ADRs existentes

- **ADR-032:** Knowledge with Temporal Validity — define metadatos y validación de vigencia que este ADR usa
- **ADR-028:** Límite auditor/asesor — estado LATE no prohíbe, solo avisa
- **ADR-008:** Vigencia — este ADR implementa control de vigencia operativo

### Fuente de datos

**Archivo:** `knowledge/taxation/spain/FILING_DEADLINES.md`

Estructura:
```yaml
---
title: "Plazos de declaración (Hacienda 2026)"
vigencia:
  valid_from: "2026-01-01"
  valid_until: "2026-12-31"
  last_verified: "2026-07-05"
  source: "AEAT, BOE"
  confidence: "high"
---

IRPF 2025:
  ordinario_hasta: "2026-06-30"
  late_hasta: "2027-06-30"

Modelo721 2025:
  ordinario_hasta: "2026-03-31"
  late_hasta: "2026-10-31"
```

**CRÍTICO:** Las fechas específicas viven **SOLO en este archivo**, no en ADRs.

### Integración en `/spanish-tax-return`

```
Paso -1 (Pre-flight):
  - Pedir ejercicio
  - Validar temporal (este ADR)
  - Si no OK: PARAR

Paso 0: Auditoría
Paso 1+: Declaración
```

### Pendientes

- **[PENDIENTE]** Crear `knowledge/taxation/spain/FILING_DEADLINES.md` con metadatos ADR-032
- **[PENDIENTE]** Implementar máquina de estados en `/spanish-tax-return`
- **[PENDIENTE]** Protocolo para prórrogas extraordinarias (AEAT a veces amplía plazos)
