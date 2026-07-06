---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.1
---

<!-- v1.1 (2026-07-05): añadido "Deciders" al checklist universal de pre-commit
     (campo MADR que faltaba en los 42 ADRs, detectado por ADR Explorer). -->

# ADR-030: Validación y verificación de ADRs críticos antes de commit

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-04
**Accepted:** 2026-07-05 — protocolo en uso activo (aplicado en la corrección de Nivel C del 2026-07-05); pendientes de implementación futura (sub-checklist fiscal, automatización de detección) no bloquean la aceptación del protocolo en sí.

## Context

Un ADR es una **decisión documentada** del proyecto. Cuando ese ADR toca:
- **Fiscalidad española** (umbrales, tratamientos, normativa)
- **Cifras/umbrales concretos** (50.000€, FIFO, comisiones)
- **Matemática o cálculos** (fórmulas FIFO, valoración, ganancias)
- **Hechos verificables** (API de exchanges, funciones de CoinTracking, cambios normativos)

...un error en el ADR es un error **en la gobernanza del agente**. Y ese error se propaga a:
1. Las skills (`/audit-cointracking`, `/spanish-tax-return`)
2. Las decisiones del agente (qué se audita, cómo se calcula)
3. Los informes que van a Hacienda

**Caso real (2026-07-04):** ADR-028 afirmaba que el umbral del Modelo 721 era 600.000€. Es 50.000€ (vigente 2026-07-04). Mañana podría cambiar a 45.000€. Un usuario que siguiera el ADR declararía mal y Hacienda lo castigaría.

**El problema arquitectónico subyacente:** Los ADRs no deberían ser custodios de cifras mutables. Esas cifras deberían vivir en `knowledge/` con metadatos de vigencia. Los ADRs las referencian, no las repiten.

Sin un **protocolo de validación antes de commit**, este tipo de errores volverán a pasar.

## Decision

Se establece un **protocolo obligatorio de validación para ADRs críticos** antes de commitearlos. El protocolo consta de:

### Nivel 1: Clasificación del ADR

Cada ADR nuevo se clasifica en una de estas categorías:

#### 🔴 CRÍTICO (verificación exhaustiva obligatoria)

Toca cualquiera de estos temas:

- **Fiscalidad española:** umbrales (Modelo 721, IRPF), tratamientos (FIFO, ganancias patrimoniales, rendimientos), normativa (AEAT, DGT, BOE)
- **Cifras concretas:** importe, porcentaje, fecha límite, umbral numérico
- **Cálculos o fórmulas:** FIFO, base de coste, ganancias, tasas, comisiones
- **Hechos verificables sobre exchanges:** límites API, formatos, funciones técnicas de CoinTracking/Binance/Kraken/etc.
- **Cambios regulatorios:** leyes nuevas, cambios de requisitos, suspensiones (p. ej. MiCA y salida de exchanges)

**Ejemplos:** ADR-002, ADR-003, ADR-028, ADR-029, cualquier ADR sobre fiscalidad futura (ADR-031+)

#### 🟡 MEDIO (revisión estándar)

Toca decisiones arquitectónicas pero sin cifras concretas:

- Estructura de datos, integración MCP, caché, persistencia
- Procesos y flujos internos del agente
- Convenciones de nombres, formatos de documentación

**Ejemplos:** ADR-010, ADR-013, ADR-016, ADR-025, ADR-027

#### 🟢 BAJO (lectura rápida)

Documentación, convenciones, organizativos:

- Idioma, formato de fechas, convenciones de archivos
- Índices, referencias, estructura del repo

**Ejemplos:** ADR-001, ADR-005

---

### Nivel 2: Checklist de validación (por categoría)

#### Para ADRs CRÍTICOS (🔴)

**Antes de abrir el PR o hacer commit, ejecutar TODOS estos puntos:**

1. ✅ **Verificación de hechos contra fuente oficial:**
   - Si menciona un umbral fiscal (50.000€, 600.000€, etc.): verificar contra AEAT, BOE, o `knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md`
   - Si menciona un cambio en normativa: citar el documento oficial (Ley, RD, Orden, Resolución DGT) con número y fecha
   - Si menciona un comportamiento de CoinTracking/API: verificar contra `knowledge/cointracking/reference/CATALOG.md` o documento oficial del exchange
   - **Si no puedes verificarlo: MARCAR CON `[VERIFICAR]` Y NO AFIRMAR COMO HECHO**

2. ✅ **Citas de fuente en cada afirmación crítica:**
   - Cada cifra/umbral tiene una nota al pie o comentario: `(Fuente: AEAT Modelo 721, consultado 2026-07-04)`
   - Cada normativa: número de ley/RD, artículo, fecha
   - Ninguna cifra sin atribución

3. ✅ **Contradicción con conocimiento existente:**
   - ¿Este ADR contradice algo en `knowledge/`?
   - ¿Hay un conflicto con otro ADR?
   - Si sí: resolver antes de commitear (editar ambos documentos, dejar notas de cross-reference)

4. ✅ **Separación clara: hecho vs especulación, y NO hardcoding de cifras mutables:**
   - ❌ "El umbral es 50.000€" — NO (mañana cambia). ✅ "El umbral está en knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md"
   - ❌ "Especulación: Probablemente Hacienda acepte esto" ← NO OK (marcar como `[ESTIMA]` o eliminar)
   - Si hay incertidumbre: marcada explícitamente como `[PENDIENTE FUNDAMENTAR]` o `[SUPUESTO]`
   - **REGLA DE ORO:** Cualquier cifra fiscal/umbral/porcentaje que cambie con la normativa → NO en el ADR, referenciar el archivo de `knowledge/` donde vive (con "última verificación" para que se reverifi que si es antigua)

5. ✅ **Revisión de ejemplos numéricos:**
   - Si hay un ejemplo con dinero/FIFO/cálculo: ¿las matemáticas son correctas?
   - Ejemplo: "Ganancia = 45.000€ - 30.000€ = 15.000€". Verificar: 45k - 30k = 15k ✓
   - Si el ejemplo es ilustrativo (no real): aclarar explícitamente: "Ejemplo hipotético:"

6. ✅ **Referencia a ADRs relacionados:**
   - ¿Hay un ADR anterior que cubre esto parcialmente?
   - ¿Se cita correctamente?
   - ¿Hay duplicación de criterio entre ADRs?

7. ✅ **Vigencia (ADR-008) + PROHIBICIÓN de hardcoding de cifras mutables:**
   - Si es fiscalidad: ¿es vigente para 2026? (cambios anuales)
   - Si es normativa: ¿cambió recientemente?
   - Si hay duda: agregar "Última verificación: YYYY-MM-DD. Reverificar si es antigua."
   - **CRÍTICO:** Si el ADR menciona una cifra fiscal (umbral, tramo, porcentaje, comisión, límite), NO puede decir "es X€". Debe decir "ver knowledge/taxation/spain/[archivo].md (última verificación YYYY-MM-DD)". La cifra vive en `knowledge/`, no en el ADR. Los ADRs referencian, no repiten.

8. ✅ **Lectura completa por quien no lo escribió:**
   - **Obligatorio:** Antes de commitear, otra persona lee el ADR completo
   - Objetivo: detectar errores, ambigüedades, datos incorrectos
   - Resultado: aprobación explícita ("Revisado y OK para commitear") o cambios solicitados

---

#### Para ADRs MEDIO (🟡)

1. ✅ Lectura una vez para coherencia
2. ✅ Verificación rápida de que no contradice otros ADRs
3. ✅ Referencias correctas (links a archivos, otros ADRs)

#### Para ADRs BAJO (🟢)

1. ✅ Lectura rápida de sintaxis/formato

---

### Nivel 3: Checklist de pre-commit (universal)

Antes de `git commit`, TODOS los ADRs cumplen:

```
☐ Formato MADR correcto (Status, Deciders, Date, Context, Decision, Consequences, Notes)
☐ Índice en adr/README.md actualizado
☐ CHANGELOG.md actualizado (si es importante)
☐ No hay typos o errores gramaticales obvios
☐ No hay frases ambiguas ("típicamente", "probablemente") sin marcar
☐ Si es crítico (🔴): checklist de validación 1-8 completado
☐ Si es crítico (🔴): aprobación de otra persona
```

---

### Nivel 4: Responsabilidad

**Quién es responsable:**

- **Autor del ADR:** Completar la validación, solicitar revisión
- **Revisor:** Leer el ADR completo, verificar puntos críticos, aprobar o solicitar cambios
- **Que no commitee sin aprobación de revisor (si es 🔴 o 🟡)**

**Quién puede ser revisor:**
- Claude Code (el agente actual)
- Un segundo Claude/LLM con contexto
- El usuario (si es su especialidad)
- Otro desarrollador (si hay equipo)

---

## Consequences

**Positive:**

- **Prevención de errores fiscales:** Cada cifra/umbral está verificado antes de publicarse
- **Trazabilidad de fuentes:** Un ADR siempre cita dónde vienen los datos
- **Protección legal:** Si hay error, está documentado que se hizo el esfuerzo de verificar
- **Escala:** El protocolo es reproducible — sigue siendo manual, pero es un proceso claro
- **Confianza:** Los usuarios/asesores saben que los ADRs pasaron revisión

**Negative:**

- **Más lento:** ADRs críticos tardan más (verificación + revisión)
- **Bloqueos:** Un ADR crítico no se commitea sin aprobación
- **Requiere disciplina:** El autor tiene que resistir la tentación de commitear sin revisar
- **Requiere expertise:** La revisión de ADRs fiscales necesita alguien que entienda fiscalidad
- **Falsos positivos:** A veces algo que parece "crítico" no lo es (requiere buen juicio)

## Notes

### Relación con ADRs existentes

- **ADR-001:** Convenciones — este ADR es una convención de proceso de desarrollo
- **ADR-006:** Límite de determinismo — los ADRs documentan decisiones; si están mal, el agente actúa mal
- **ADR-008:** Vigencia del conocimiento — este ADR operacionaliza la verificación de vigencia
- **ADR-009:** Protocolo crítico — este ADR refuerza "cero invención" documentando cómo verificar
- **ADR-011:** Persistencia — los ADRs son parte de la persistencia (gobernanza duradera)

### Implementación

1. **En este repositorio:** Aplicar el checklist antes de commitear cada nuevo ADR (empezando con éste)

2. **Para futuros ADRs sobre fiscalidad:** Crear un sub-checklist específico
   - Ej: "Si es sobre Modelo 721, verificar contra `knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md`"

3. **Herramienta opcional:** Crear un script de pre-commit que detecte:
   - Cifras fiscales mencionadas sin fuente
   - Afirmaciones definitivas ("esto ES", "debes") sin `[VERIFICAR]` o cita
   - Referencias a `knowledge/` que no existen

4. **Documentación:** Este checklist va en `CLAUDE.md` como referencia para Copilot también (ADR-012: división de responsabilidades)

### Pendientes

- **[PENDIENTE]** Crear un sub-checklist específico para ADRs de fiscalidad (umbrales, tratamientos, normativa)
- **[PENDIENTE]** Definir "revisor calificado" para ADRs fiscales (¿quién valida la asesoría fiscal?)
- **[PENDIENTE]** Automatizar detección de cifras fiscales sin fuente en pre-commit hook
- **[PENDIENTE]** Documentar "excepciones" al checklist (cuándo un ADR NO necesita revisión exhaustiva)
- **[PENDIENTE]** Crear un archivo de "Errores aprendidos" — registro de errores encontrados en ADRs y cómo se previnieron

### Por qué esto importa

**ADR-028 tenía 600.000€ en lugar de 50.000€.** Un usuario que siguiera ese ADR declararía sin Modelo 721 cuando estaba obligado. Hacienda lo castigaría. El error se originó porque **no había un proceso de validación**. Este ADR lo fija.

**"Desarrollar una herramienta de dominio crítico" significa: no pueden pasar errores como éste.**
