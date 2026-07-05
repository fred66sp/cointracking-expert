---

## Resumen Ejecutivo

**Fase 1 (HOY):** Crear estructura y formalizar arquitectura  
**Fase 2 (próxima sesión):** Reorganizar directorios y crear índices  
**Fase 3 (iterativo):** Agregar metadatos y crear documentos nuevos  

**Riesgo:** Bajo (la Fase 1 no rompe nada)

---


# Plan de Migración — Sistema de Conocimiento Jerárquico

**Documento:** Plan operativo para implementar ADR-033 en 3 fases  
**Status:** Fase 1 completa, Fase 2-3 pendientes  
**Última actualización:** 2026-07-05  
**Dificultad:** Media (reestructuración sin romper nada)



## Fase 1 ✅ COMPLETADA (2026-07-05)

### Tareas

- [x] Crear **ADR-033** — Sistema de Conocimiento Jerárquico
- [x] Crear **plantilla YAML** — `knowledge/.metadata/METADATA_TEMPLATE.md`
- [x] Crear **INDEX_MASTER.md** — Mapa completo del sistema
- [x] Crear **estructura vacía** de directorios (sin mover archivos)
- [x] Documentar **MIGRATION_PLAN.md** (este documento)

### Commits

**Commit 1 (próximamente):**
- Crear ADR-033
- Crear `knowledge/.metadata/METADATA_TEMPLATE.md`
- Crear `knowledge/INDEX_MASTER.md`
- Crear directorios vacíos
- Crear este plan

**Mensaje:**
```
feat(knowledge): implementar ADR-033 — Sistema de Conocimiento Jerárquico (Fase 1)

- Crear ADR-033: formaliza arquitectura de 6 niveles (A-F)
- Crear plantilla YAML estándar para metadatos (validada contra ADR-032)
- Crear INDEX_MASTER.md: mapa navegable de todo el sistema
- Crear estructura vacía: directorios de niveles A-F (sin mover archivos)
- Documentar plan de migración Fase 2 y 3

No hay cambios en archivos existentes. Fase 1 = preparación.
```

---

## Fase 2 ⏭️ PRÓXIMA SESIÓN

### Tarea: Reorganizar directorios

**Objetivo:** Mover archivos existentes a sus nuevas ubicaciones según ADR-033, crear índices cruzados

### Subtareas

#### Subtarea 2.1: Reorganizar `cointracking/`

**Mover:**
```
cointracking/CSV_FORMAT.md
cointracking/COST_BASIS_AND_VALIDATION.md
→ cointracking/official/
```

**Crear referencias:**
- En `cointracking/INDEX.md` — actualizar referencias (ahora puntan a `official/`)
- En skills (`/audit-cointracking`, `/spanish-tax-return`) — actualizar imports (si las hay)
- En otros ADRs — actualizar referencias (si las hay)

**Validar:**
- ✅ Links en INDEX.md funcionan
- ✅ Links en skills funcionan
- ✅ Links en docs funcionan

---

#### Subtarea 2.2: Crear `cointracking/behavioral/`

**Crear índice:** `cointracking/behavioral/INDEX.md`

**Contenido del índice:**
```markdown
# Conocimiento Operativo de CoinTracking (Nivel B1)

Esta carpeta documenta cómo funciona *realmente* CoinTracking, validado contra datos reales.

**Documentos pendientes:**
- BALANCE_CALCULATION_ALGORITHM.md
- PURCHASE_POOL_MECHANICS.md
- MISSING_PURCHASE_HISTORY_CAUSES.md
- DUPLICATE_DETECTION_HEURISTICS.md
- API_VS_CSV_OVERLAP.md
- FEE_HANDLING.md

**Status:** Fase 3 (crear cuando se identifique necesidad)
```

---

#### Subtarea 2.3: Reorganizar `exchanges/`

**Renombrar y mover:**
```
exchanges/BINANCE.md
exchanges/BINANCE_EU_MICA_EXIT.md
→ exchanges/official/BINANCE.md
→ exchanges/reference/BINANCE_EU_MICA_EXIT.md
```

**Crear índices:**
- `exchanges/official/INDEX.md` — Índice de exchanges oficiales (A3)
- `exchanges/behavioral/INDEX.md` — Índice de operativas (B2)

**Validar:**
- ✅ Links en docs actualizados
- ✅ Links en skills actualizados

---

#### Subtarea 2.4: Migrar `patterns/` → `cases/`

**Renombrar:**
```
patterns/cointracking_casos_v2.yaml
→ cases/cointracking_casos_v2.yaml (TODAVÍA EN YAML, no convertir a .md aún)
```

**Crear índices:**
- `cases/INDEX.md` — Índice de 20 casos (C1)
  
**No hacer en Fase 2:**
- ❌ NO convertir YAML a .md (eso es Fase 3)
- ❌ NO agregar metadatos YAML a cada caso (eso es Fase 3)

**Razonamiento:** Separar reorganización de refactoring de formato

**Validar:**
- ✅ El YAML sigue siendo válido en nueva ubicación
- ✅ Enlaces desde `knowledge/` actualizados
- ✅ Enlaces desde skills actualizados

---

#### Subtarea 2.5: Crear carpeta `procedures/`

**Crear:**
- `knowledge/procedures/INDEX.md` — Plantilla vacía

**Contenido:**
```markdown
# Procedimientos Operativos (Nivel C3)

Guías paso a paso validadas en proyecto real.

**Pendientes (Fase 3):**
- PROCEDURE_AUDIT_ACCOUNT.md
- PROCEDURE_RECONCILE_TRANSFERS.md
- PROCEDURE_IMPORT_CSV.md
- PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md
- PROCEDURE_VALIDATE_FIFO.md
```

---

#### Subtarea 2.6: Crear `adr/INDEX.md`

**Ubicación:** `adr/INDEX.md`

**Contenido:** Mapeo de ADRs a niveles de conocimiento

**Estructura:**
```markdown
# ADRs — Índice de Decisiones Arquitectónicas

## Por nivel de conocimiento

### ADRs del Nivel A (Fuentes oficiales)
- ADR-031 (Validación de plazos) → usa conocimiento fiscal de A1

### ADRs del Nivel B (Operativo)
- ADR-004 (Reconciliación con datos reales)

### ADRs del Nivel C (Casos)
- ADR-026 (Matriz decisiones) → usa caso CT-002 (FLOKI)
- ADR-014 (Validación de duplicados) → usa caso CT-002

### ADRs del Nivel D (Auxiliar)
- ADR-029 (Protocolo de no-hacer) → define checklists (D1)

### ADRs del Nivel F (Governance)
- ADR-033 (Sistema de Conocimiento Jerárquico)
- ADR-032 (Knowledge with Temporal Validity)

[... más detalle]
```

---

#### Subtarea 2.7: Migrar `docs/GLOSSARY.md` → `knowledge/reference/GLOSSARY.md`

**Mover:**
```
docs/GLOSSARY.md
→ knowledge/reference/GLOSSARY.md
```

**Crear referencia en docs:**
```markdown
# Glosario

Ver `knowledge/reference/GLOSSARY.md`
```

**Validar:**
- ✅ Links actualizados en docs/
- ✅ Links actualizados en skills/

---

#### Subtarea 2.8: Crear índices cruzados

**Crear:**
- `knowledge/.metadata/AUTHORITY_MATRIX.md` — Quién verifica cada nivel
- `knowledge/taxation/spain/INDEX.md` — Si no existe

**Ejemplo de AUTHORITY_MATRIX.md:**
```markdown
# Matriz de Autoridad — Quién Verifica Qué

| Nivel | Responsable | Frecuencia | Metodología |
|-------|---|---|---|
| A1 (España) | Alfredo (revisión anual) | Enero | Contrastar AEAT/BOE/DGT |
| A2 (CoinTracking) | Alfredo + casos reales | Cuando CT actualiza | Casos reales vs CT behavior |
| A3 (Exchanges) | Alfredo + documentación | Cuando exchange actualiza | Documentación oficial |
| B1-B3 | Alfredo + casos | Siempre que se audite | Datos reales vs comportamiento |
| C1 | Alfredo (auditoría) | Siempre que se audite caso | Proyecto real |
| C2-C3 | Alfredo (después de casos) | Cada 2-3 casos | Generalización |
| D1-D3 | Alfredo (derivado de C) | Cuando cambien C | Síntesis |
| F | Alfredo (governance) | Cuando decida | Consenso |

[... más detalle]
```

---

### Commits (Fase 2)

**Commit 1:** Reorganizar `cointracking/`
```
refactor(knowledge): reorganizar documentos de CoinTracking (Nivel A2)

- Mover CSV_FORMAT.md, COST_BASIS_AND_VALIDATION.md → cointracking/official/
- Crear cointracking/behavioral/ (vacío, para Fase 3)
- Actualizar referencias en INDEX.md
- Actualizar references en skills
```

**Commit 2:** Reorganizar `exchanges/`
```
refactor(knowledge): reorganizar documentos de exchanges

- Mover BINANCE.md → exchanges/official/
- Mover BINANCE_EU_MICA_EXIT.md → knowledge/reference/context/
- Crear exchanges/behavioral/ (vacío, para Fase 3)
- Actualizar referencias
```

**Commit 3:** Migrar casos + crear índices
```
refactor(knowledge): migrar casos y crear índices (Fase 2)

- Mover cointracking_casos_v2.yaml → knowledge/cases/
- Crear adr/INDEX.md
- Crear knowledge/.metadata/AUTHORITY_MATRIX.md
- Crear índices en procedures/, checklists/, decision-trees/
- Migrar docs/GLOSSARY.md → knowledge/reference/
```

**Total:** 3 commits (uno por subsistema reorganizado)

---

## Fase 3 ⏭️ ITERATIVO (Semanas 2-3)

### Tarea: Agregar metadatos y crear documentos nuevos

**Objetivo:** Completar la transición, automatizar vigencia, crear contenido nuevo

### Subtareas

#### Subtarea 3.1: Agregar metadatos YAML a documentos existentes

**Archivos a actualizar (prioridad alta):**

1. `knowledge/cointracking/official/CSV_FORMAT.md`
   - Actualizar con frontmatter YAML
   - `id: KB-A2-001`
   - `authority: official`

2. `knowledge/cointracking/official/COST_BASIS_AND_VALIDATION.md`
   - `id: KB-A2-002`
   - `authority: verified` (contrastado contra datos reales)

3. `knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md`
   - `id: KB-A1-001`
   - `authority: official`

4. `knowledge/taxation/spain/CAPITAL_GAINS.md`
   - `id: KB-A1-002`
   - `authority: official`

5. Todos los documentos existentes (iterativamente)

**Script:** Crear `scripts/add-metadata-skeleton.py`
- Lee documento sin metadatos
- Propone metadatos basados en ubicación + nombre
- Usuario revisa y completa (especialmente `source`, `notes`)

---

#### Subtarea 3.2: Convertir `cointracking_casos_v2.yaml` → archivos `.md` individuales

**Proceso:**
1. Parse YAML
2. Para cada caso: crear archivo `.md` en `knowledge/cases/`
3. Insertar metadatos YAML en cada archivo
4. Crear `knowledge/cases/INDEX.md` actualizado

**Ejemplo:**
```
knowledge/cases/CT-002-floki-batching.md
---
id: KB-C1-002
title: "CT-002: FLOKI — 29 transacciones idénticas no son duplicadas"
level: C
...
---

# CT-002: FLOKI...
[contenido]
```

**Script:** Crear `scripts/yaml-to-md-cases.py` (Python)

**Validar:**
- ✅ Todos los casos convertidos
- ✅ No se pierde información del YAML
- ✅ Metadatos son válidos

**Nota:** Mantener `cointracking_casos_v2.yaml` como backup (solo para esta sesión, luego borrar)

---

#### Subtarea 3.3: Crear documentos de Nivel C2 (Patrones)

**Patrones a derivar de los 20 casos:**

- `PATTERN_DUPLICATE_DETECTION.md` — Matriz: qué hace/no hace duplicado
- `PATTERN_BALANCE_RECONCILIATION.md` — Cómo reconocer inconsistencias
- `PATTERN_TRANSFER_MATCHING.md` — Heurísticas de matching
- `PATTERN_PURCHASE_POOL_EXHAUSTION.md` — Síntomas de agotamiento

**Cada patrón:**
- Explica el patrón (no es un caso concreto)
- Cita los casos que lo documentan
- Da la heurística
- Explica excepciones

**Commits:** Uno por cada patrón

---

#### Subtarea 3.4: Crear documentos de Nivel C3 (Procedimientos)

**Extraer de skills:**

- De `/audit-cointracking` Paso 1-6:
  - `PROCEDURE_AUDIT_ACCOUNT.md`

- De `/spanish-tax-return` Paso 0-7:
  - `PROCEDURE_PREPARE_TAX_DECLARATION.md`

- De guías existentes:
  - `PROCEDURE_RECONCILE_TRANSFERS.md`
  - `PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md`

**Cada procedimiento:**
- Paso a paso
- Con imágenes/diagramas si aplica
- Cita los ADRs relacionados
- Enlaza a los casos de referencia

**Commits:** Uno por procedimiento

---

#### Subtarea 3.5: Crear documentos de Nivel D1 (Checklists)

**De los procedimientos + casos:**

- `CHECKLIST_NEGATIVE_BALANCES.md`
- `CHECKLIST_DUPLICATES.md`
- `CHECKLIST_WARNINGS.md`
- `CHECKLIST_FIFO_VALIDATION.md`
- `CHECKLIST_AUDIT_COMPLETE.md`
- `CHECKLIST_TAX_DECLARATION.md`

**Formato:** Lista de control ([ ] items)

**Commits:** Uno por cada checklist o grupo de checklists relacionados

---

#### Subtarea 3.6: Crear documentos de Nivel D2 (Árboles de Decisión)

**Diagramas / máquinas de estado:**

- `FLOW_AUDIT.md` — Diagrama completo
- `FLOW_DUPLICATE_DETECTION.md` — ¿Duplicado?
- `FLOW_TRANSFER_MATCHING.md` — ¿Transfer?
- `FLOW_MISSING_PURCHASE_HISTORY.md` — ¿Missing origen?
- `FLOW_FISCAL_DECISION.md` — ¿Conflicto fiscal?

**Formato:** Markdown con diagramas ASCII o Mermaid

**Commits:** Uno por árbol

---

#### Subtarea 3.7: Crear validador de metadatos

**Script:** `scripts/validate-knowledge-metadata.py`

**Valida:**
1. YAML válido
2. Campos obligatorios presentes
3. Tipos de datos correctos
4. Referencias (related_adr, related_docs) existen
5. Vigencia coherente
6. Nivel declarado coincide con ubicación

**Integración:** Pre-commit hook

**Commits:** Uno para el script

---

#### Subtarea 3.8: Crear alertas de vigencia

**Script:** `scripts/vigencia-alerts.py`

**Ejecutar periódicamente:**
1. Encuentra documentos próximos a `valid_until`
2. Genera informe
3. Actualiza `knowledge/.metadata/VIGENCIA_ALERTS.md`

**Commits:** Uno para el script

---

### Commit Pattern (Fase 3)

**Cada documento/patrón:**
```
feat(knowledge): agregar [documento] (Nivel [X])

- Crear [documento]
- Metadatos YAML con vigencia (ADR-032)
- Enlazar a casos/ADRs relacionados
- Validar con script
```

**Cada patrón:**
```
docs(patterns): generalizar [patrón] desde [casos]

- Crear [PATTERN_*.md]
- Citar casos C1, C2, C3 que fundamentan el patrón
- Dar heurísticas claras
```

---

## Validación General

### Después de Fase 1

- [ ] ADR-033 creado y revisado
- [ ] Plantilla YAML clara
- [ ] INDEX_MASTER es navegable
- [ ] Directorios existen (vacíos)

### Después de Fase 2

- [ ] Todos los archivos reorganizados
- [ ] Índices cruzados creados
- [ ] Références actualizadas (docs, skills, ADRs)
- [ ] NO hay links rotos

### Después de Fase 3

- [ ] Todos los documentos tienen metadatos YAML válidos
- [ ] Todos los casos en archivos `.md` individuales
- [ ] Patrones, procedimientos, checklists, árboles creados
- [ ] Validador y alertas de vigencia funcionan
- [ ] Sistema es 100% operacional

---

## Rollback / Reversibilidad

**Fase 1:** Sin riesgo (solo crea nuevos archivos, no toca existentes)

**Fase 2:** Bajo riesgo (solo reorganiza directorios; links son actualizables)
- Reversible: `git reset --hard [commit-anterior-a-fase2]`

**Fase 3:** Bajo riesgo (agrega información, no la borra)
- Reversible: los originales siguen en git

---

## Timeline Estimado

| Fase | Sesión | Horas | Status |
|------|--------|-------|--------|
| 1 | 2026-07-05 | 2h | ✅ COMPLETA |
| 2 | 2026-07-06 | 1.5h | ⏭️ Próxima |
| 3.1-3.3 | 2026-07-07 | 2h | ⏭️ Próxima |
| 3.4-3.8 | 2026-07-08+ | 3-4h | ⏭️ Iterativo |

**Total:** ~8-9 horas distribuidas (muy bajo riesgo)

---

## Notas

- Este plan NO es bloqueante. Si hay urgencias, cualquier fase es pausable.
- La Fase 1 está completa. Podemos saltar a Fase 2 cuando sea conveniente.
- La Fase 3 puede hacerse en paralelo a otros trabajos (es iterativa).
- Cada commit es pequeño, reversible, sin riesgo.

---

## Validación

```yaml
id: F3-MIGRATION-001
title: "Plan de Migración — Sistema de Conocimiento Jerárquico"
level: F
domain: other
source: "ADR-033"
authority: reference
last_verified: 2026-07-05
valid_from: 2026-07-05
valid_until: null
confidence: high
version: 1.0
related_adr:
  - ADR-033
tags:
  - migration
  - plan
  - governance
notes: "Documento vivo: se actualiza cuando cambian fases."
```
