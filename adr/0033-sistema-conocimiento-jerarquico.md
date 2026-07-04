# ADR-033: Sistema de Conocimiento Jerárquico

**Status:** Accepted

**Date:** 2026-07-05

## Context

El agente audita criptos y prepara declaraciones fiscales. Necesita conocimiento que:

1. **Venga de fuentes verificables** — no puedo usar "de memoria", Wikipedia, o sospechas
2. **Tenga vigencia clara** — un dato de 2024 puede estar obsoleto en 2026
3. **Declare autoridad** — qué es "ley", qué es "cómo funciona CoinTracking", qué es "patrón visto"
4. **Permita trazabilidad** — cada conclusión remite a su origen (AEAT, CoinTracking, casos reales)

**Problema sin este ADR:**

```
Usuario pregunta: "¿Por qué no se declara esta operación de staking?"

Puedo responder basándome en:
- Normativa (fuente: DGT)
- Documentación de CoinTracking (fuente: centro de ayuda)
- Casos que he visto (patrón empírico)

¿Cuál es más autoritativa? ¿Cuándo cambió? ¿Quién lo verifica?

Hay confusión.
```

**ADR-032 (Knowledge with Temporal Validity) operacionaliza "vigencia".**

**Este ADR operacionaliza "estructura jerárquica" y "niveles de autoridad".**

---

## Decision

Crear un **Sistema de Conocimiento con 6 Niveles Jerárquicos**, donde cada documento declara:

1. **Qué nivel es** (A=oficial, B=operativo, C=casos, D=auxiliar, E=referencia, F=governance)
2. **Cuándo fue verificado**
3. **Cuándo caduca**
4. **Quién lo verifica**

### Arquitectura de 6 Niveles

#### Nivel A — Fuentes Oficiales (Authoritative)

Estas son las únicas fuentes sobre las que el agente puede fundamentar conclusiones sin advertencias adicionales.

**A1: España (AEAT, BOE, DGT)**

Documentos que referencian directamente normativa oficial:
- `IRPF_CAPITAL_GAINS.md` — Ley del IRPF, reglamento, consultas DGT
- `MODELO_721.md` — Obligación informativa, umbrales, plazos
- `FILING_DEADLINES.md` — Plazos ordinario/extraordinario (ADR-031)
- `STAKING_CLASSIFICATION.md` — Rendimiento vs ganancia patrimonial (DGT)
- `AIRDROPS_CLASSIFICATION.md` — Valuación, momento exigible
- etc.

**A2: CoinTracking Oficial**

Documentos que referencian directo a CoinTracking:
- `CSV_FORMAT.md`
- `COST_BASIS_AND_VALIDATION.md`
- `API_REFERENCE.md` (nuevo)
- `IMPORT_MECHANISMS.md` (nuevo)
- etc.

**A3: Exchanges Oficiales**

Documentación oficial de cada exchange:
- `BINANCE_OFFICIAL.md`
- `KRAKEN_OFFICIAL.md`
- etc.

**Ubicación propuesta:** `knowledge/authorities/` y subdirectorios específicos

---

#### Nivel B — Conocimiento Operativo (Behavioral)

Explica **cómo funciona realmente** CoinTracking, exchanges y blockchain. Validado contra datos reales pero no es normativa.

**B1: CómoFunciona CoinTracking**

- `BALANCE_CALCULATION_ALGORITHM.md`
- `PURCHASE_POOL_MECHANICS.md`
- `MISSING_PURCHASE_HISTORY_CAUSES.md`
- `DUPLICATE_DETECTION_HEURISTICS.md`
- `API_VS_CSV_OVERLAP.md`
- etc.

**B2: Operativas de Exchanges**

- `BINANCE_MODULES.md` — Spot, Convert, Earn, Futures, Margin, etc.
- `BINANCE_IMPORT_LIMITATIONS.md`
- `KRAKEN_STAKING_MECHANICS.md`
- etc.

**B3: Blockchain**

- `ETHEREUM.md` — Tipos de tx, fees, gas, bridges
- `BITCOIN.md` — UTXO model
- etc.

**Ubicación propuesta:** `knowledge/cointracking/behavioral/`, `knowledge/exchanges/behavioral/`, `knowledge/blockchains/`

---

#### Nivel C — Empirismo Verificado (Verified Cases)

Conocimiento derivado de **casos reales auditados** — no es teoría, es práctica.

**C1: Casos Reales Auditados**

Actualmente: `patterns/cointracking_casos_v2.yaml` (20 casos en un YAML)

Propuesta: Reorganizar a `knowledge/cases/` con archivos individuales:

```
knowledge/cases/
├── CT-001-duplicate-same-timestamp.md
├── CT-002-floki-batching.md
├── CT-003-missing-purchase-history.md
└── INDEX.md
```

**C2: Patrones Recurrentes**

Generalización desde casos (no son casos concretos, son patrones):
- `PATTERN_DUPLICATE_DETECTION.md`
- `PATTERN_BALANCE_RECONCILIATION.md`
- `PATTERN_TRANSFER_MATCHING.md`
- etc.

**C3: Procedimientos Operativos**

Paso a paso validado en proyecto real:
- `PROCEDURE_AUDIT_ACCOUNT.md`
- `PROCEDURE_RECONCILE_TRANSFERS.md`
- `PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md`
- etc.

**Ubicación propuesta:** `knowledge/cases/`, `knowledge/patterns/`, `knowledge/procedures/`

---

#### Nivel D — Auxiliar (Supporting)

Herramientas para el agente: no aportan conocimiento nuevo, organizan el existente.

**D1: Checklists**

Listas estructuradas:
- `CHECKLIST_NEGATIVE_BALANCES.md`
- `CHECKLIST_DUPLICATES.md`
- `CHECKLIST_WARNINGS.md`
- `CHECKLIST_FIFO_VALIDATION.md`

**D2: Árboles de Decisión**

Máquinas de estado:
- `FLOW_AUDIT.md` — ¿balance negativo? → investigar → resolver
- `FLOW_DUPLICATE_DETECTION.md` — ¿misma fecha? → ¿misma TX ID?
- `FLOW_FISCAL_DECISION.md` — ¿hay conflicto fiscal? → A/B/C (ADR-026)

**D3: Índices y Referencias Cruzadas**

- `knowledge/INDEX_MASTER.md` — Mapa completo del sistema
- Índices de búsqueda por dominio, tag, ADR relacionado

**Ubicación propuesta:** `knowledge/checklists/`, `knowledge/decision-trees/`, `knowledge/reference/`

---

#### Nivel E — Referencia (Reference)

Contexto, definiciones, historiadores. No envejece rápido.

**E1: Glosario**

- Definiciones (FIFO, purchase pool, base de coste, ganancia patrimonial, etc.)

**E2: Contexto e Historiadores**

- `REGULATORY_TIMELINE.md` — Hitos regulatorios (DAC8, MiCA, etc.)
- `BINANCE_EU_MICA_EXIT.md` — Salida de Binance UE 2026-07
- etc.

**Ubicación propuesta:** `knowledge/reference/`

---

#### Nivel F — Governance (Governance)

Cómo funciona el sistema de conocimiento.

**F1: ADRs (Decisiones Arquitectónicas)**

- Los ADRs existentes (ya en `adr/`)
- Nuevo: `adr/INDEX.md` — Mapea ADR → Nivel de conocimiento que usa

**F2: Índices Maestros**

- `knowledge/INDEX_MASTER.md` — Arquitectura completa

**F3: Metadatos del Sistema**

- `knowledge/.metadata/SCHEMA.yaml` — Esquema YAML oficial
- `knowledge/.metadata/AUTHORITY_MATRIX.md` — Quién verifica cada nivel
- `knowledge/.metadata/DEPRECATION_LOG.md` — Documentos obsoletos
- `knowledge/.metadata/VIGENCIA_ALERTS.md` — Próximos a expirar

**Ubicación propuesta:** `adr/`, `knowledge/`, `knowledge/.metadata/`

---

### Metadatos Obligatorios (validado contra ADR-032)

Todo documento declara **frontmatter YAML** con:

```yaml
id: KB-[NIVEL]-[NÚMERO]
title: Título descriptivo
level: A|B|C|D|E|F
domain: cointracking|taxation|exchanges|blockchain|other
source: "AEAT/BOE/DGT/CoinTracking/Binance/Caso/Patrón"
authority: official|verified|empirical|reference
last_verified: YYYY-MM-DD
valid_from: YYYY-MM-DD
valid_until: YYYY-MM-DD (null = indefinido)
confidence: high|medium|low
version: 1.0
related_adr: [ADR-###, ...]
related_docs: [archivo.md, ...]
tags: [tag1, tag2, ...]
notes: "Próxima revisión: ..."
```

**Regla crítica:** Documentos de Nivel A (`authority: official`) **NUNCA deben tener `valid_until: null`** — siempre especificar fecha de caducidad.

---

### Niveles de Autoridad (atributo `authority`)

| Valor | Significa | Uso permitido |
|-------|-----------|--------------|
| `official` | Viene de fuente oficial (AEAT, BOE, CoinTracking, exchange) sin interpretación | Puedo fundamentar conclusiones fiscales/técnicas. No requiero advertencia. |
| `verified` | Auditado en proyecto real y funciona | Puedo apoyar diagnósticos; cito la fuente si hay duda |
| `empirical` | Patrón observado en casos; tendencia, no certeza | Propongo hipótesis; no es fundamento único |
| `reference` | Solo contexto, definición, historiador | Solo para explicar conceptos |

---

## Consequences

### Positive

✅ **Trazabilidad completa:** Cada conclusión remonta a su origen (oficial, verificado, empírico)

✅ **Vigencia clara:** Sé exactamente cuándo revisar cada documento (validado contra ADR-032)

✅ **Confianza graduada:** No confundo "ley" con "patrón visto"

✅ **Escalabilidad:** Sistema extensible para nuevos exchanges, nuevos casos, nueva normativa

✅ **Responsabilidad:** Queda claro quién verifica qué y cuándo

✅ **Integración con ADR-032:** Los metadatos YAML permiten que ADR-032 (Knowledge with Temporal Validity) sea operacional

### Negative

⚠️ **Overhead:** Más metadatos = más mantenimiento

⚠️ **Riesgo de abandono:** Si nadie actualiza vigencias, el sistema se degrada

⚠️ **Complejidad inicial:** Reestructurar 15+ documentos

---

## Notes

### Relación con ADRs existentes

- **ADR-008:** Vigencia — este ADR operacionaliza la estructura de vigencia
- **ADR-032:** Knowledge with Temporal Validity — define metadatos YAML que este ADR usa
- **ADR-031:** Validación de plazos — usa Nivel A1 (FILING_DEADLINES.md)
- **ADR-026:** Matriz decisiones — usa Nivel C1 (casos FLOKI)
- **ADR-009:** Protocolo crítico — refuerza "cero invención" documentando fuentes

### Plan de Implementación

**Fase 1 (HOY) — Formalización:**

1. ✅ Crear ADR-033 (este documento)
2. ✅ Crear plantilla YAML estándar (`knowledge/.metadata/METADATA_TEMPLATE.md`)
3. ✅ Crear `knowledge/INDEX_MASTER.md` — Mapa maestro de todo el sistema
4. ✅ Crear estructura vacía de directorios (no mover archivos aún)
5. ✅ Documentar plan de migración

**Fase 2 (próxima sesión) — Reorganización:**

1. Reorganizar directorios dentro de `knowledge/` (mover CSV_FORMAT, COST_BASIS, etc.)
2. Actualizar referencias en docs/, skills/, adr/
3. Crear `adr/INDEX.md` — Mapea ADR a nivel de conocimiento

**Fase 3 (iterativo) — Metadatos:**

1. Añadir frontmatter YAML a documentos existentes
2. Crear nuevos documentos (procedures, checklists, decision-trees)
3. Validar metadatos en pre-commit

### Pendientes (no bloqueantes)

- [ ] Crear validador YAML de metadatos (script Python)
- [ ] Automatizar alertas de documentos próximos a expirar
- [ ] Definir "quién verifica qué" por nivel (responsabilidad)
- [ ] Dashboard de vigencia de documentos

### Precedentes

La propuesta de arquitectura jerárquica fue sugerida por el usuario el 2026-07-05 (conversación con Alfredo González P., proyecto h--cointracking-expert).

Se consolida en **6 niveles** (A-F) en lugar de 10, manteniendo toda la funcionalidad requerida, para mejor mantenibilidad.

La estructura está validada contra:
- La propuesta original del usuario
- ADR-008 (vigencia)
- ADR-032 (metadatos temporales)
- CLAUDE.md (protocolo crítico)
- Estructura actual de `knowledge/` (no rompe nada)
