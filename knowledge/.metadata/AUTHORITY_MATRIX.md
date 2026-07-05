---

## Matriz por Nivel

| Nivel | Responsable Principal | Frecuencia | Metodología | Status |
|-------|---|---|---|---|
| **A1** — España (AEAT, BOE, DGT) | Alfredo (revisión anual) | Enero | Contrastar AEAT/BOE/DGT + consultas DGT | ✅ |
| **A2** — CoinTracking (oficial) | Alfredo + casos reales | Cuando CT actualiza | Verificar en web/API + validar contra datos reales | ✅ |
| **A3** — Exchanges (oficial) | Alfredo + doc. oficial | Cuando exchange actualiza | Documentación oficial del exchange | ⚠️ (solo Binance) |
| **B1** — Operativo CT | Alfredo + casos reales | Cada auditoría | Contrastar casos reales vs. comportamiento CT | ✅ (casos v2) |
| **B2** — Operativo exchanges | Alfredo + doc. oficial | Cuando exchange actualiza | Documentación oficial + casos reales | ⚠️ (Binance) |
| **B3** — Blockchain | Alfredo + casos | Raramente (blockchain estable) | Documentación oficial blockchain | ⚠️ (vacío) |
| **C1** — Casos auditados | Alfredo (auditor) | Cada caso nuevo | Proyecto real + validación verificable | ✅ (20 casos v2) |
| **C2** — Patrones | Alfredo (después de casos) | Cada 2-3 casos | Generalización de casos | ⚠️ (vacío) |
| **C3** — Procedimientos | Alfredo (operativo) | Cuando cambian skills | Experiencia operativa probada | ⚠️ (vacío) |
| **D1** — Checklists | Alfredo (síntesis de C) | Cuando cambian C | Compilación de casos/patrones | ⚠️ (vacío) |
| **D2** — Árboles decisión | Alfredo (síntesis de C) | Cuando cambian C | Formalización de lógica operativa | ⚠️ (vacío) |
| **D3** — Índices | Alfredo (mantenimiento) | Cada cambio estructura | Revisión manual | ⚠️ (parcial) |
| **E1** — Glosario | Alfredo (referencia) | Raramente | Definiciones estándar | ✅ |
| **E2** — Historiadores | Alfredo (contexto) | Rara vez | Documentación histórica | ⚠️ (parcial) |
| **F1** — ADRs | Alfredo (decisiones) | Cuando hay decisión | Consenso + documentación | ✅ (33 ADRs) |
| **F2** — Índices maestros | Alfredo (navegación) | Cada cambio estructura | Actualización manual | ✅ |
| **F3** — Metadatos sistema | Alfredo (gobernanza) | Cuando cambian reglas | Especificación YAML | ✅ (ADR-032) |

**Leyenda:**
- ✅ = Completo / Vigente
- ⚠️ = Parcial / Pendiente
- ❌ = Ausente

---


# Matriz de Autoridad — Quién Verifica Qué

**Documento:** Define responsabilidades de verificación por nivel de conocimiento (ADR-033)



## Política de Actualización

### Niveles A (Fuentes Oficiales)

**Frecuencia:** Anual (enero)

**Proceso:**
1. Revisar `valid_until` de cada documento A
2. Contrastar contra fuente oficial (AEAT, BOE, DGT, CoinTracking, exchange)
3. Si cambió: actualizar metadatos + aumentar `version`
4. Commit: `chore(knowledge): update A-level documents for YYYY`

**Responsable:** Alfredo (usuario principal auditor)

---

### Niveles B (Operativo)

**Frecuencia:** Cuando CoinTracking/exchange actualiza + cada auditoría (validar empíricamente)

**Proceso:**
1. Cuando sea relevante para el caso en auditoría: validar B contra datos reales
2. Si hay discrepancia: documentar en AGENT_CHANGE_REQUESTS.md
3. En Fase 3, formalizar como documento B o caso C1

**Responsable:** Alfredo (auditoría empírica)

---

### Niveles C (Casos)

**Frecuencia:** Cada vez que se audita un caso nuevo

**Proceso:**
1. Auditoría completa de un caso del usuario
2. Documentar: síntomas → diagnóstico → solución
3. Cotejar contra conocimiento existente (A, B)
4. Si caso nuevo: agregar a C1 (ahora `cases/cointracking_casos_v2.yaml`)
5. Si detecta patrón recurrente: considerar C2 (patrones generalizados)
6. Si crea procedimiento: considerar C3 (procedimientos)

**Responsable:** Alfredo (auditor especializado)

---

### Niveles D (Auxiliar)

**Frecuencia:** Después de cada grupo de casos/cambios en C

**Proceso:**
1. Síntesis de casos (C1) → Patrones (C2)
2. Formalización de lógica operativa → Árboles de decisión (D2)
3. Compilación de checklists desde procedimientos (D3)
4. Validación: ¿los D son coherentes con C?

**Responsable:** Alfredo (después de auditorías)

---

### Niveles E (Referencia)

**Frecuencia:** Rara vez (referencias raramente envejecen)

**Proceso:**
1. E1 (Glosario): actualizar definiciones si cambia terminología
2. E2 (Historiadores): actualizar solo si nuevo evento regulatorio

**Responsable:** Alfredo (mantenimiento)

---

### Nivel F (Governance)

**Frecuencia:** Cuando hay decisiones o cambios estructura

**Proceso:**
1. **F1 (ADRs):** Crear ADR para decisiones arquitectónicas
2. **F2 (Índices):** Actualizar índices cuando cambia estructura
3. **F3 (Metadatos):** Aplicar SCHEMA.yaml cuando cambian reglas

**Responsable:** Alfredo (decisiones) + Claude Code (formalización)

---

## Alertas de Vigencia

### Automáticas (Fase 3+)

Script `vigencia-alerts.py` genera informe diario:

- Documentos A: próximos a `valid_until` dentro de 30 días
- Documentos B: `last_verified` > 180 días (confidence: medium) o > 30 días (confidence: low)
- Documentos C: `last_verified` > 1 año

Actualizar `knowledge/.metadata/VIGENCIA_ALERTS.md`.

---

## Escalabilidad

Cuando se agregue un nuevo nivel o dominio:

1. Definir en `INDEX_MASTER.md` (ADR-033)
2. Agregar fila a esta matriz
3. Especificar: responsable, frecuencia, metodología
4. Crear documento de proceso similar a arriba

---

## Vigencia de Esta Matriz

```yaml
id: F3-AUTHORITY-001
title: "Matriz de Autoridad — Quién Verifica Qué"
level: F
source: "ADR-033"
authority: reference
last_verified: 2026-07-05
valid_from: 2026-07-05
valid_until: null
confidence: high
version: 1.0
related_adr:
  - ADR-033
  - ADR-032
  - ADR-008
notes: "Documento vivo: se actualiza cuando cambian responsabilidades o frecuencias."
```
