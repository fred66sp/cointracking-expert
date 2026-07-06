# Flujo de Governance — Cómo Registrar Decisiones (ADRs)

**Documento:** Cómo crear y aprobar Decisiones Arquitectónicas (ADRs)  
**Audiencia:** Desarrolladores, arquitectos, responsables del proyecto  
**Última actualización:** 2026-07-05

---

## 🎯 Propósito

Un **ADR (Architectural Decision Record)** es un documento que registra:

- ¿**Qué** decisión se tomó?
- ¿**Por qué** (contexto, alternativas consideradas)?
- ¿**Cuáles** son las consecuencias?
- ¿**Cuándo** entra en vigor?

**No es:**
- Un comentario de código
- Un commit message largo
- Una documentación de cómo usar algo (eso va en `knowledge/`)

**Es:**
- Un registro **permanente** de "por qué existe X así"
- Justificación para que otros entiendan el razonamiento
- Trazabilidad de decisiones pasadas

---

## 📋 Cuándo Crear un ADR

### ✅ Crea un ADR si:

- Cambias la **arquitectura** (p. ej. de SDK a MCP)
- Introduces un **nuevo nivel** de conocimiento (p. ej. Nivel C)
- Tomas una **decisión importante** con impacto duradero (p. ej. CoinTracking CSV vs API)
- Decides **deprecar** algo existente (p. ej. un método de auditoría)
- Resuelves un **dilema técnico** (p. ej. FIFO vs LIFO para España)
- Estableces una **política** (p. ej. no invertir skills, solo auditar)

### ❌ NO creas un ADR si:

- Es una **corrección** (errata, bug fix)
- Es una **adición menor** (un doc nuevo sin afectar arquitectura)
- Es una **operación puntual** (cambio de credenciales, actualización de metadatos)
- Es un **experimento temporal** (prueba rápida que quizá se descarta)

---

## 🏗️ Estructura de un ADR (MADR v2)

Usa la plantilla **MADR 2.0** (Markdown Architecture Decision Records):

```markdown
---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-NXX: (Título)

**Status:** Accepted  <!-- Proposed | Accepted | Deprecated | Superseded -->

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** AAAA-MM-DD

## Context

¿Cuál es la situación que requiere una decisión?
¿Qué fuerzas o restricciones nos afectan?
¿Qué hemos intentado antes?

## Decision

¿Qué decidimos hacer? (Breve, claro, activo.)

## Rationale

¿Por qué esta es la mejor opción?
¿Qué alternativas consideramos?
¿Cuál es el trade-off principal?

## Consequences

¿Qué beneficios trae?
¿Qué costos o riesgos tiene?
¿Qué necesita cambiar como resultado?

## Related

- [ADR-XXX](ADR-XXX.md) (Decisión anterior relacionada)
- [knowledge/...](../knowledge/.../...) (Documentación técnica)
```

---

## 1️⃣ Ejemplo: ADR Completo

### El Problema

Después de auditar varias cuentas, notamos que el sistema necesita un **nuevo nivel de conocimiento** para registrar casos verificados. Los niveles actuales (A-E) no cubren "un caso real auditado que se vuelve patrón".

### El ADR

```markdown
---
id: ADR-033
title: "Sistema de Conocimiento Jerárquico (6 Niveles)"
status: "Accepted"
date: "2026-06-15"
relates_to:
  - ADR-006 (Determinar arquitectura del agente)
  - ADR-010 (Eficiencia de caché)
---

# ADR-033: Sistema de Conocimiento Jerárquico de 6 Niveles

## Context

El agente auditor necesita una **base de conocimiento estructurada**:

- Actualmente tienen "docs sueltos" en varias carpetas
- No hay diferenciación clara entre "fuente oficial" vs "caso verificado"
- Metadatos inconsistentes (algunas docs tienen vencimiento, otras no)
- Usuarios no saben dónde buscar información
- Imposible saber si un documento está actualizado

**Alternativas consideradas:**

1. **Flat structure** (todos los docs al mismo nivel)
   - ✅ Simple
   - ❌ Confuso, sin jerarquía

2. **Hierarchical by domain** (carpetas por tema: CoinTracking, Taxation)
   - ✅ Organizado por tema
   - ❌ Mezcla fuentes oficiales con análisis propio

3. **Hierarchical by authority** (6 niveles A-F: oficial → operativo → verificado → auxiliar → referencia → governance)
   - ✅ Claro dónde buscar cada tipo de info
   - ✅ Fácil de validar (cada nivel tiene reglas)
   - ✅ Trazable (qué es oficial vs análisis propio)
   - ✅ Escalable (nuevos docs van a nivel correcto automáticamente)

**Elegimos opción 3.**

## Decision

Implementamos un **sistema de 6 niveles jerárquicos**:

- **Nivel A:** Fuentes Oficiales (AEAT, CoinTracking, exchanges)
- **Nivel B:** Operativo (Cómo funciona, comportamientos verificados)
- **Nivel C:** Casos Reales (Auditados, patrones, procedimientos)
- **Nivel D:** Auxiliar (Checklists, árboles de decisión)
- **Nivel E:** Referencia (Glosario, contexto, historiadores)
- **Nivel F:** Governance (ADRs, decisiones registradas)

Cada nivel tiene:
- **Ubicación clara** (carpeta, patrón de nombre)
- **Metadatos obligatorios** (YAML frontmatter con 11 campos)
- **Authority level** (official, verified, reference)
- **Confidence score** (high, medium, low)
- **Vencimiento** (valid_until; obligatorio para Nivel A)

## Rationale

**¿Por qué esta estructura?**

1. **Claridad:** Usuario nuevo no se pierde
2. **Confianza:** Sabe qué documentos son "ley" vs "análisis"
3. **Mantenibilidad:** Cada nivel tiene reglas claras
4. **Escalabilidad:** Fácil agregar nuevos documentos
5. **Trazabilidad:** Metadata permite saber si algo está vigente

**Trade-off principal:**

- ❌ Requiere más estructura (frontmatter YAML obligatorio)
- ❌ Validación más estricta
- ✅ Pero: Valor a largo plazo supera el costo inicial

## Consequences

**Beneficios:**

- ✅ Sistema escalable y mantenible
- ✅ Usuarios pueden confiar en el conocimiento
- ✅ Fácil detectar documentos vencidos
- ✅ Patrones reutilizables (otros auditorías similares)

**Costos:**

- ❌ Reorganizar documentos existentes (~40 docs)
- ❌ Escribir scripts de validación
- ❌ Documentar la arquitectura (este ADR + guías)

**Cambios requeridos:**

- Crear carpetas para Nivel D, E (antes no existían)
- Migrar 20 YAML cases a markdown individual con KB-C1-XXX IDs
- Crear script de validación metadata
- Actualizar índices maestros

## Related

- [ADR-006](ADR-006.md) — Determinar arquitectura del agente (SDK vs LLM)
- [ADR-010](ADR-010.md) — Eficiencia y caché de tokens
- [ADR-032](ADR-032.md) — Temporal validity pattern (vencimiento de docs)
- [knowledge/INDEX_MASTER.md](knowledge/INDEX_MASTER.md) — Implementación completa
```

---

## 2️⃣ Crear un ADR Nuevo

### Paso 1: Determinar el Número

```bash
# Busca el ADR más reciente
ls -la adr/ADR-*.md | tail -1
# ADR-033.md

# Próximo: ADR-034
```

### Paso 2: Crear el Archivo

```bash
touch adr/ADR-034.md
```

### Paso 3: Llenar la Plantilla

Copia y personaliza:

```markdown
---
id: ADR-034
title: "Tu Decisión Aquí"
status: "Proposed"  # Comienza como "Proposed"
date: "2026-07-05"
relates_to:
  - ADR-033
  - ADR-006
---

# ADR-034: (Tu Decisión)

## Context

(¿Cuál es el problema?)

## Decision

(¿Qué decidimos?)

## Rationale

(¿Por qué?)

## Consequences

(¿Qué cambia?)

## Related

- [ADR-XXX](ADR-XXX.md)
```

### Paso 4: Revisar y Aprobar

- **Tú solo:** Marca como `status: "Accepted"` directamente
- **Equipo:** Deja como `status: "Proposed"` y pide feedback primero

### Paso 5: Commit

```bash
git commit -m "ADR-034: (Título)

Context: (1-2 líneas)
Decision: (1 línea)
"
```

---

## 3️⃣ Actualizar un ADR

### Si Fue Reemplazado

Marca el viejo como `Superseded`:

```markdown
---
id: ADR-032
title: "..."
status: "Superseded by ADR-033"
---
```

Actualiza el nuevo para referenciar:

```markdown
---
id: ADR-033
title: "..."
relates_to:
  - ADR-032  # Reemplaza esto
---
```

### Si Cambió de Dirección

```markdown
status: "Accepted (revised 2026-07-05)"
```

Documenta el cambio al final:

```markdown
## Revision History

- 2026-07-05: Actualizado para incluir Nivel C (casos)
- 2026-06-15: Original
```

---

## 4️⃣ Estatus de un ADR

| Estado | Significado | Cuándo |
|--------|------------|--------|
| **Proposed** | Idea, esperando feedback | Acabo de crear |
| **Accepted** | Decidido e implementando | Después de revisar |
| **Deprecated** | Decisión antigua, ya no aplica | Tras varios años |
| **Superseded by ADR-NXX** | Reemplazado por otra decisión | Cuando cambiamos de dirección |

---

## 5️⃣ Índice de ADRs

Actualiza `adr/INDEX.md` cuando crees uno nuevo:

```markdown
# ADR Index

## Active Decisions

- [ADR-033: Sistema de Conocimiento Jerárquico](ADR-033.md) — 6 niveles A-F
- [ADR-034: (Tu ADR)](ADR-034.md) — Descripción corta
- ...

## Historical (Deprecated/Superseded)

- [ADR-001: ...](ADR-001.md) — Superseded by ADR-015
```

---

## 6️⃣ Ejemplos de ADRs en Este Proyecto

| ADR | Decisión | Estado |
|-----|----------|--------|
| [ADR-001](adr/ADR-001.md) | Idioma español para contenido | Accepted |
| [ADR-006](adr/ADR-006.md) | Arquitectura: LLM auditor, no SDK | Accepted |
| [ADR-009](adr/ADR-009.md) | Protocolo de auditoría crítico | Accepted |
| [ADR-010](adr/ADR-010.md) | Eficiencia de tokens y caché | Accepted |
| [ADR-014](adr/ADR-014.md) | Falsos positivos en duplicados (Trade ID) | Accepted |
| [ADR-033](adr/ADR-033.md) | Sistema de 6 niveles jerárquicos | Accepted |

---

## 7️⃣ Checklist para Crear un ADR

- [ ] Número único (ADR-XXX sin duplicados)
- [ ] Sección **Context:** Claro ¿cuál es el problema?
- [ ] Sección **Decision:** Breve, activo ("Implementamos X", no "Decidimos considerar X")
- [ ] Sección **Rationale:** Alternativas consideradas + trade-offs
- [ ] Sección **Consequences:** Beneficios + costos + cambios requeridos
- [ ] Sección **Related:** Links a otros ADRs o documentos relevantes
- [ ] Status correcto (`Proposed` o `Accepted`)
- [ ] **Deciders** asignados (campo MADR; sin él, herramientas como ADR Explorer marcan la decisión como huérfana — detectado 2026-07-05)
- [ ] Fecha correcta (YYYY-MM-DD)
- [ ] Frontmatter YAML con `version:` (invalidación de caché, ADR-039)
- [ ] Commit message claro
- [ ] Actualizados **ambos** índices: `adr/README.md` y `adr/INDEX.md`

---

## 8️⃣ Flujo de Decisión: ¿Necesito un ADR?

```
┌─ ¿Es un cambio importante/duradero?
│
├─ SÍ, afecta arquitectura
│  └─ → Crear ADR (status: Proposed)
│     → Revisar si es en equipo
│     → Marcar como Accepted cuando decidas
│     → Implementar cambios
│     → Actualizar adr/INDEX.md
│
└─ NO, es operativo/táctico
   └─ → Commit directamente
      → Descripción clara en commit message
      → No necesita ADR
```

---

## 🔟 Referencia Rápida

| Necesito | Qué Hacer |
|----------|-----------|
| Ver decisiones pasadas | `cat adr/INDEX.md` |
| Crear ADR nuevo | `touch adr/ADR-NXX.md` + plantilla MADR |
| Revisar contexto de decisión X | `cat adr/ADR-NXX.md` |
| Entender por qué existe Nivel B | `cat adr/ADR-033.md` |
| Proponer cambio importante | Crear ADR-NXX con status: Proposed |
| Marcar ADR como decidido | Cambiar status: Accepted |
| Reemplazar ADR viejo | Marcar como "Superseded by ADR-NXX" |

---

## 🚪 Siguiente

- [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) — Cómo compilar y arrancar el servidor MCP
- [knowledge/KNOWLEDGE_MAINTENANCE.md](knowledge/KNOWLEDGE_MAINTENANCE.md) — Cómo mantener la base de conocimiento
- [adr/README.md](adr/README.md) — Guía MADR completa
