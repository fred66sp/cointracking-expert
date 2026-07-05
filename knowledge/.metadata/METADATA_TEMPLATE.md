---
id: KB-B1-XXX
title: "Plantilla de Metadatos YAML para Documentos de Conocimiento"
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

# Plantilla de Metadatos YAML para Documentos de Conocimiento

## Uso

Cada documento en `knowledge/` comienza con un **frontmatter YAML** estructurado. Este archivo documenta el esquema exacto.

---

## Frontmatter YAML Canónico

```yaml
---
id: KB-[NIVEL]-[NÚMERO]
title: Título descriptivo del documento
level: A|B|C|D|E|F
domain: cointracking|taxation|exchanges|blockchain|other
source: "AEAT/BOE/DGT/CoinTracking/Binance/Caso real/Patrón/Documentación"
authority: official|verified|empirical|reference
last_verified: YYYY-MM-DD
valid_from: YYYY-MM-DD
valid_until: YYYY-MM-DD o null (null = indefinido)
confidence: high|medium|low
version: 1.0

related_adr:
  - ADR-002
  - ADR-032

related_docs:
  - COST_BASIS_AND_VALIDATION.md
  - cointracking_casos_v2.yaml

tags:
  - reconciliation
  - balances
  - cointracking

notes: "Próxima revisión: enero 2027"
---

[Contenido del documento aquí]
```

---

## Guía de Cada Campo

### `id`

**Formato:** `KB-[NIVEL]-[NÚMERO]`

**Niveles:** A, B, C, D, E, F

**Número:** Secuencial dentro de cada nivel

**Ejemplos:**
- `KB-A1-001` — Primer documento del Nivel A1
- `KB-B2-003` — Tercer documento del Nivel B2
- `KB-C1-042` — Caso real #42

**Propósito:** Identificación única para referencias cruzadas

---

### `title`

**Formato:** Cadena descriptiva (no código)

**Ejemplos:**
- "Formato CSV de CoinTracking (Trade Table)"
- "FLOKI: 29 transacciones idénticas no son duplicadas"
- "Plazos de presentación de IRPF (España 2026)"

---

### `level`

**Valores permitidos:** A, B, C, D, E, F

**Mapeo a función:**
- **A** — Fuentes oficiales (AEAT, BOE, DGT, CoinTracking oficial, exchanges oficiales)
- **B** — Conocimiento operativo (cómo funciona CoinTracking, comportamiento de exchanges)
- **C** — Empirismo verificado (casos reales auditados, patrones, procedimientos)
- **D** — Auxiliar (checklists, árboles de decisión, índices)
- **E** — Referencia (glosario, historiadores, contexto)
- **F** — Governance (ADRs, metadatos del sistema)

**Propósito:** Define autoridad implícita del documento

---

### `domain`

**Valores sugeridos:**
- `cointracking` — Cómo funciona CoinTracking
- `taxation` — Normativa tributaria (IRPF, Modelo 721, etc.)
- `exchanges` — Particularidades de exchanges (Binance, Kraken, etc.)
- `blockchain` — Tecnología blockchain
- `procedures` — Pasos operativos
- `patterns` — Patrones recurrentes
- `other` — Otros dominios

**Propósito:** Categorización para búsqueda y filtrado

---

### `source`

**Formato:** Cadena descriptiva (cita la fuente primaria)

**Ejemplos:**
- `"AEAT — Modelo 721, preguntas frecuentes"`
- `"CoinTracking — Centro de ayuda oficial + datos reales proyecto agp"`
- `"Binance — Documentación oficial de API"`
- `"Caso auditado en proyecto agp, 2024-03"`
- `"Generalización de 20 casos en cointracking_casos_v2.yaml"`

**Propósito:** Permite al usuario verificar la fuente original

---

### `authority`

**Valores permitidos:** `official`, `verified`, `empirical`, `reference`

**Significado:**

| Valor | Significa | Uso permitido |
|-------|-----------|--------------|
| `official` | Viene de una fuente oficial sin interpretación | Puedo fundamentar conclusiones fiscales/técnicas. No requiero advertencia. |
| `verified` | He auditado esto en proyecto real y funciona | Puedo apoyar diagnósticos; cito fuente si hay duda |
| `empirical` | Patrón observado en casos; no es certeza | Propongo hipótesis; no es fundamento único |
| `reference` | Solo contexto, definición, historiador | Solo para explicar conceptos |

**Propósito:** Permite al agente graduarse en confianza

---

### `last_verified`

**Formato:** YYYY-MM-DD

**Ejemplos:**
- `2026-07-05` — Fresco
- `2025-06-15` — 1 año de antigüedad

**Propósito:** ADR-032 usa esto para detectar documentos envejecidos

**Regla operativa:** El agente alerta si `Hoy - last_verified > [días según confidence]`:
- `confidence: high` → máx 365 días sin verificar
- `confidence: medium` → máx 180 días
- `confidence: low` → máx 30 días

---

### `valid_from`

**Formato:** YYYY-MM-DD

**Ejemplos:**
- `2026-01-01` — Válido desde enero
- `2024-09-01` — Cambio de normativa

**Propósito:** Marca cuándo el conocimiento comienza a ser válido

**Regla:** Si `Hoy < valid_from` → documento aún no aplicable, no usar

---

### `valid_until`

**Formato:** YYYY-MM-DD o `null`

**Ejemplos:**
- `2026-12-31` — Válido solo en 2026
- `null` — Indefinido (raramente cambia)

**Propósito:** Marca cuándo caduca

**Regla crítica:** Documentos de `authority: official` (Nivel A) **NUNCA deben tener `valid_until: null`** — siempre especificar fecha

**Validación:** Si `Hoy > valid_until` → envejecido, no usar sin reverificar

---

### `confidence`

**Valores permitidos:** `high`, `medium`, `low`

**Significado:**

| Valor | Significa | Cuando usar |
|-------|-----------|------------|
| `high` | Verificado contra múltiples fuentes o casos reales | Nivel A o múltiples casos C |
| `medium` | Verificado pero con incertidumbres/cambios potenciales | Nivel B (pueden cambiar con updates) |
| `low` | Hipótesis, patrón visto pocas veces, fuente no oficial | Documentos en validación |

**Propósito:** El agente decide si necesita reverificar

---

### `version`

**Formato:** Semver: MAYOR.MENOR.PARCHE

**Ejemplos:**
- `1.0` — Versión inicial estable
- `1.1` — Actualización menor
- `2.0` — Cambio mayor

**Propósito:** Trackear evolución

---

### `related_adr`

**Formato:** Lista de ADRs relacionados

**Ejemplos:**
```yaml
related_adr:
  - ADR-032  # Define cómo este documento debe tener metadatos
  - ADR-031  # Usa este documento para validar plazos
  - ADR-009  # Protocolo crítico que aplica aquí
```

**Propósito:** Entiender gobernanza (si ADR cambia, ¿afecta este documento?)

---

### `related_docs`

**Formato:** Lista de otros documentos relacionados

**Ejemplos:**
```yaml
related_docs:
  - COST_BASIS_AND_VALIDATION.md
  - cointracking_casos_v2.yaml
  - knowledge/checklists/CHECKLIST_DUPLICATES.md
```

**Propósito:** Navegación entre documentos

---

### `tags`

**Formato:** Lista de palabras clave en minúsculas

**Ejemplos:**
```yaml
tags:
  - reconciliation
  - duplicates
  - cointracking
  - auditing
```

**Propósito:** Búsqueda y categorización

---

### `notes`

**Formato:** Texto libre

**Ejemplos:**
```yaml
notes: "Próxima revisión: enero 2027 (cambio anual de tramos IRPF)"
notes: "Status: patrón observado 3+ veces, pendiente ADR"
```

**Propósito:** Notas operativas para mantenimiento

---

## Ejemplos Completos

### Ejemplo 1: Nivel A (Oficial)

```yaml
---
id: KB-A1-001
title: "Tramos de la base del ahorro (IRPF 2026)"
level: A
domain: taxation
source: "AEAT — Normativa fiscal 2026, BOE, DGT"
authority: official
last_verified: 2026-07-05
valid_from: 2026-01-01
valid_until: 2026-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-032
  - ADR-031

tags:
  - taxation
  - irpf
  - capital-gains

notes: "Cambio anual. Requiere reverificación cada enero."
---

# Tramos de la base del ahorro (IRPF 2026)

[Contenido...]
```

---

### Ejemplo 2: Nivel B (Operativo)

```yaml
---
id: KB-B1-002
title: "FIFO: Cómo evoluciona el purchase pool en CoinTracking"
level: B
domain: cointracking
source: "CoinTracking — Centro de ayuda + análisis de casos reales"
authority: verified
last_verified: 2026-07-04
valid_from: 2025-01-01
valid_until: null
confidence: high
version: 1.2

related_adr:
  - ADR-003
  - ADR-004

related_docs:
  - COST_BASIS_AND_VALIDATION.md
  - cointracking_casos_v2.yaml

tags:
  - cointracking
  - fifo
  - purchase-pool
  - cost-basis

notes: "CoinTracking puede cambiar; reverificar si output varía."
---

# FIFO: Cómo evoluciona el purchase pool en CoinTracking

[Contenido...]
```

---

### Ejemplo 3: Nivel C (Caso Real)

```yaml
---
id: KB-C1-002
title: "CT-002: FLOKI — 29 transacciones idénticas no son duplicadas"
level: C
domain: cointracking
source: "Caso auditado en proyecto agp, 2024-03"
authority: verified
last_verified: 2026-07-03
valid_from: 2024-03-17
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-014
  - ADR-026
  - ADR-003

tags:
  - duplicates
  - trade-id
  - case-study
  - binance

notes: "Caso de referencia para ADR-014. Trade IDs validados."
---

# CT-002: FLOKI — 29 transacciones idénticas no son duplicadas

[Caso completo...]
```

---

### Ejemplo 4: Nivel D (Checklist)

```yaml
---
id: KB-D1-001
title: "Checklist: Detección de duplicados"
level: D
domain: cointracking
source: "Generalización de 20 casos + ADR-014"
authority: verified
last_verified: 2026-07-05
valid_from: 2026-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-014
  - ADR-026

related_docs:
  - knowledge/cases/CT-002-floki-batching.md
  - knowledge/patterns/PATTERN_DUPLICATE_DETECTION.md

tags:
  - checklist
  - duplicates
  - audit

notes: "Usar antes de recomendar borrado."
---

# Checklist: Detección de duplicados

- [ ] ¿Misma fecha (mismo segundo)?
- [ ] ¿Mismo precio?
- [ ] ¿Mismo volumen?
- [ ] ¿Misma comisión?
- [ ] **✅ Verificar Trade ID en Binance API**
  - Si diferentes → legítimas, NO ELIMINAR
  - Si iguales → posible duplicado, pedir confirmación

[...]
```

---

## Validación de Metadatos

Antes de hacer commit, validar:

1. ✅ YAML válido (usar `yamllint` o equivalente)
2. ✅ `id` único en repositorio
3. ✅ `level` está en {A, B, C, D, E, F}
4. ✅ `authority` está en {official, verified, empirical, reference}
5. ✅ `confidence` está en {high, medium, low}
6. ✅ `last_verified` es YYYY-MM-DD y no es futura
7. ✅ `valid_from` ≤ `valid_until` (si `valid_until` no es null)
8. ✅ Si `authority: official`, entonces `valid_until` NO es null
9. ✅ `related_adr` apuntan a ADRs que existen
10. ✅ `related_docs` apuntan a documentos que existen
11. ✅ Tiene al menos 3 tags

**Script:** `scripts/validate-knowledge-metadata.py` (por crear en Fase 3)

---

## Procedimiento de Actualización Anual

**Cada enero (o cuando cambie normativa):**

1. Revisar todos los documentos con `valid_until <= ahora`
2. Contrastar contra fuente oficial
3. Actualizar `last_verified`, `valid_from`, `valid_until`
4. Aumentar `version` (MENOR si solo update, MAYOR si cambio)
5. Actualizar `notes`
6. Commit: `chore(knowledge): update A-level documents for 2027`

---

## FAQ

**P: ¿Si olvido actualizar `last_verified`?**

R: Es responsabilidad del mantenedor. Automatizar en CI/CD: alertar si no se verifica en X meses.

---

**P: ¿Un documento sin `related_adr`?**

R: Aceptable para Nivel D/E. No recomendable para A/B/C.

---

**P: ¿Cambio `version` de 1.0 a 1.1 = actualizar `last_verified`?**

R: Sí. `last_verified` marca cuándo se verificó/actualizó, aunque sea un cambio menor.

---

## Relación con ADR-032

Este esquema de metadatos es el operacional de **ADR-032: Knowledge with Temporal Validity**.

ADR-032 define **qué hacer** cuando un documento está envejecido (alertar, bloquear, reverificar).

Esta plantilla define **qué documentar** para que ADR-032 pueda actuar.
