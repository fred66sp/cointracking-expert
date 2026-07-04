# Sistema de Conocimiento — Mapa Maestro

**Documento:** Descripción general de la arquitectura jerárquica de conocimiento del agente  
**Validez:** Permanente (es estructura, no contenido)  
**Última actualización:** 2026-07-05  
**Arquitectura:** Definida en ADR-033

---

## 🎯 Propósito

Este documento es un **índice navegable** de todos los niveles de conocimiento (A-F). Permite:

1. Entender **dónde vive cada tipo de conocimiento**
2. **Navegar rápidamente** por tipo de información
3. **Identificar brechas** (qué documentos faltan)
4. **Validar referencias** (cada documento sabe dónde buscar más)

---

## Arquitectura de 6 Niveles

```
┌──────────────────────────────────────────────────────────┐
│ NIVEL A — FUENTES OFICIALES (Authoritative)              │
│  - A1: España (AEAT, BOE, DGT)                            │
│  - A2: CoinTracking Oficial                              │
│  - A3: Exchanges Oficiales                               │
└──────────────────────────────────────────────────────────┘
         ↓ Fundamento
┌──────────────────────────────────────────────────────────┐
│ NIVEL B — OPERATIVO (Behavioral)                         │
│  - B1: Cómo funciona CoinTracking                        │
│  - B2: Operativas de Exchanges                           │
│  - B3: Blockchain                                        │
└──────────────────────────────────────────────────────────┘
         ↓ Validado contra
┌──────────────────────────────────────────────────────────┐
│ NIVEL C — CASOS REALES (Verified Cases)                 │
│  - C1: Casos auditados                                   │
│  - C2: Patrones recurrentes                              │
│  - C3: Procedimientos operativos                         │
└──────────────────────────────────────────────────────────┘
         ↓ Organiza
┌──────────────────────────────────────────────────────────┐
│ NIVEL D — AUXILIAR (Supporting)                          │
│  - D1: Checklists                                        │
│  - D2: Árboles de decisión                               │
│  - D3: Índices y referencias                             │
└──────────────────────────────────────────────────────────┘
         ↓ Contextualiza
┌──────────────────────────────────────────────────────────┐
│ NIVEL E — REFERENCIA (Reference)                         │
│  - E1: Glosario                                          │
│  - E2: Contexto e historiadores                          │
└──────────────────────────────────────────────────────────┘
         ↓ Gobernado por
┌──────────────────────────────────────────────────────────┐
│ NIVEL F — GOVERNANCE (Governance)                        │
│  - F1: ADRs (decisiones arquitectónicas)                 │
│  - F2: Índices maestros                                  │
│  - F3: Metadatos del sistema                             │
└──────────────────────────────────────────────────────────┘
```

---

## NIVEL A — Fuentes Oficiales

**Ubicación:** `knowledge/authorities/`, `knowledge/cointracking/official/`, `knowledge/exchanges/official/`

**Características:**
- ✅ Vienen de fuentes verificables (AEAT, BOE, CoinTracking, exchanges)
- ✅ `authority: official`
- ✅ **NUNCA tienen `valid_until: null`** (siempre especificar caducidad)
- ✅ `confidence: high`
- ✅ Se pueden citar directamente en conclusiones fiscales/técnicas

### A1 — España (AEAT, BOE, DGT)

**Ubicación propuesta:** `knowledge/authorities/spain/`

**Estado actual:**
- ✅ `taxation/spain/INFORMATIVE_OBLIGATIONS.md` (Modelo 721)
- ✅ `taxation/spain/CAPITAL_GAINS.md` (Ganancias patrimoniales)
- ✅ `taxation/spain/CAPITAL_INCOME.md` (Ingresos del capital)
- ❌ `IRPF_SUMMARY.md` (pendiente)
- ❌ `FILING_DEADLINES.md` (creado para ADR-031, pendiente aquí)
- ❌ `STAKING_CLASSIFICATION.md` (pendiente)
- ❌ `AIRDROPS_CLASSIFICATION.md` (pendiente)

**Próximo paso:** Migrar a A1 en Fase 2, añadir metadatos YAML

---

### A2 — CoinTracking Oficial

**Ubicación propuesta:** `knowledge/cointracking/official/`

**Estado actual:**
- ✅ `cointracking/CSV_FORMAT.md`
- ✅ `cointracking/COST_BASIS_AND_VALIDATION.md`
- ✅ `cointracking/reference/CATALOG.md` (índice de 205 artículos oficiales)
- ✅ `cointracking/MCP_API.md`
- ✅ `cointracking/WEB_APP_GUIDE.md`
- ❌ `API_REFERENCE.md` (especificar endpoints formalmente)
- ❌ `IMPORT_MECHANISMS.md` (CSV vs API vs manual)
- ❌ `TRANSACTION_TYPES.md` (10 tipos canonicales + mapeo CT)
- ❌ `BALANCE_REPORTS.md` (estructura de reportes)
- ❌ `TAX_REPORT_GENERATION.md` (cómo genera el informe fiscal)

**Próximo paso:** Migrar a A2 en Fase 2, reorganizar en subcarpetas, añadir metadatos

---

### A3 — Exchanges Oficiales

**Ubicación propuesta:** `knowledge/exchanges/official/`

**Estado actual:**
- ✅ `exchanges/BINANCE.md` (limitaciones de import, módulos)
- ✅ `exchanges/BINANCE_EU_MICA_EXIT.md` (contexto regulatorio)
- ❌ `BINANCE_API_REFERENCE.md` (formal, con endpoints)
- ❌ `KRAKEN_OFFICIAL.md`
- ❌ `COINBASE_OFFICIAL.md`
- ❌ `BYBIT_OFFICIAL.md`

**Próximo paso:** Crear en Fase 3 según demanda de nuevos exchanges

---

## NIVEL B — Operativo

**Ubicación:** `knowledge/cointracking/behavioral/`, `knowledge/exchanges/behavioral/`, `knowledge/blockchains/`

**Características:**
- ✅ Explica **cómo funciona realmente** (no es teoría, es práctica)
- ✅ `authority: verified` (contrastado contra datos reales)
- ✅ Puede cambiar si CoinTracking/exchanges actualizan
- ✅ Se pueden usar para apoyar diagnósticos

### B1 — Cómo Funciona CoinTracking

**Ubicación propuesta:** `knowledge/cointracking/behavioral/`

**Estado actual:**
- ✅ `cointracking/COST_BASIS_AND_VALIDATION.md` (purchase pool, negativos, etc.) — *nota: parcialmente B1, parcialmente A2*
- ✅ `cointracking/TROUBLESHOOTING.md` (índice de síntomas)
- ❌ `BALANCE_CALCULATION_ALGORITHM.md` (cómo evoluciona el saldo)
- ❌ `PURCHASE_POOL_MECHANICS.md` (formalizar el algoritmo)
- ❌ `MISSING_PURCHASE_HISTORY_CAUSES.md` (por qué aparece)
- ❌ `DUPLICATE_DETECTION_HEURISTICS.md` (qué detecta CT automáticamente)
- ❌ `API_VS_CSV_OVERLAP.md` (cuándo hay duplicados entre fuentes)
- ❌ `FEE_HANDLING.md` (comisiones en tercera moneda)

**Próximo paso:** Crear en Fase 3 después de casos C1

---

### B2 — Operativas de Exchanges

**Ubicación propuesta:** `knowledge/exchanges/behavioral/`

**Estado actual:**
- ✅ `exchanges/BINANCE.md` (módulos, limitaciones) — *nota: parcialmente A3, parcialmente B2*
- ❌ `BINANCE_MODULES.md` (Spot, Convert, Earn, Futures, Margin en detalle)
- ❌ `BINANCE_IMPORT_WORKFLOW.md` (cómo importar correctamente)
- ❌ `KRAKEN_STAKING_MECHANICS.md`
- ❌ `COINBASE_ADVANCED_TRADE.md`

**Próximo paso:** Crear en Fase 3 según demanda

---

### B3 — Blockchain

**Ubicación propuesta:** `knowledge/blockchains/` (ya existe)

**Estado actual:**
- ❌ `ETHEREUM.md` (tipos tx, fees, gas, bridges)
- ❌ `BITCOIN.md` (UTXO, fees, address format)
- ❌ `POLYGON.md`, `BNB_SMART_CHAIN.md`, `SOLANA.md`, etc.

**Próximo paso:** Crear en Fase 3 según demanda (no crítico para auditoría, sí para transfers)

---

## NIVEL C — Empirismo Verificado

**Ubicación:** `knowledge/cases/`, `knowledge/patterns/`, `knowledge/procedures/`

**Características:**
- ✅ Basado en **casos reales auditados**, no teoría
- ✅ `authority: verified` (probado en proyecto real)
- ✅ Confianza alta, pero no es normativa
- ✅ Se pueden usar para diagnosticar problemas

### C1 — Casos Reales Auditados

**Ubicación propuesta:** `knowledge/cases/` (reorganizar desde `patterns/`)

**Estado actual:**
- ✅ `patterns/cointracking_casos_v2.yaml` — 20 casos, esquema canónico

**Casos documentados en v2:**
- CT-001 — Duplicados mismo timestamp
- CT-002 — FLOKI (Trade IDs distintos, no son duplicados)
- CT-003 — Missing Purchase History (compras sin origen)
- CT-004, CT-005, ..., CT-020 (ver `INDEX.md` en patterns/)

**Próximo paso:**
- Fase 2: Migrar de YAML a archivos `.md` individuales en `knowledge/cases/`
- Fase 3: Añadir metadatos YAML a cada archivo

**Ejemplo propuesto:**
```
knowledge/cases/
├── INDEX.md (índice de casos)
├── CT-001-duplicate-same-timestamp.md
├── CT-002-floki-batching.md
├── CT-003-missing-purchase-history.md
└── ... (CT-004 a CT-020)
```

---

### C2 — Patrones Recurrentes

**Ubicación propuesta:** `knowledge/patterns/` (expandir desde `INDEX.md` actual)

**Estado actual:**
- ✅ `patterns/INDEX.md` (existe, documenta esquema)
- ✅ `patterns/cointracking_casos_v2.yaml` (20 casos)
- ❌ Patrones generalizados (archivos separados)

**Patrones a documentar:**
- `PATTERN_DUPLICATE_DETECTION.md` — Matriz: qué hace/no hace duplicado
- `PATTERN_BALANCE_RECONCILIATION.md` — Cómo reconocer saldos inconsistentes
- `PATTERN_TRANSFER_MATCHING.md` — Heurísticas para withdrawal/deposit
- `PATTERN_PURCHASE_POOL_EXHAUSTION.md` — Síntomas de "purchase pool consumed"
- `PATTERN_MISSING_ORIGIN.md` — Dónde viene este activo

**Próximo paso:** Crear en Fase 3 (derivar de los 20 casos)

---

### C3 — Procedimientos Operativos

**Ubicación propuesta:** `knowledge/procedures/`

**Estado actual:**
- ❌ No existen (están en skills/*)

**Procedimientos a documentar:**
- `PROCEDURE_AUDIT_ACCOUNT.md` — 6 fases de auditoría
- `PROCEDURE_RECONCILE_TRANSFERS.md` — Emparejar withdrawal/deposit
- `PROCEDURE_IMPORT_CSV.md` — Evitar duplicados
- `PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md` — Resolver origen
- `PROCEDURE_VALIDATE_FIFO.md` — Validar cost basis
- `PROCEDURE_PREPARE_TAX_DECLARATION.md` — IRPF 2026

**Próximo paso:**
- Fase 2: Crear carpeta `knowledge/procedures/`
- Fase 3: Extraer playbooks de skills y documentarlos aquí

---

## NIVEL D — Auxiliar

**Ubicación:** `knowledge/checklists/`, `knowledge/decision-trees/`

**Características:**
- ✅ Organizan conocimiento de Niveles A-C
- ✅ No aportan nuevo conocimiento, lo estructuran
- ✅ Los LLMs los siguen muy bien
- ✅ Reducen variabilidad de respuestas

### D1 — Checklists

**Ubicación propuesta:** `knowledge/checklists/`

**Estado actual:**
- ❌ No existen formalmente (están en skills)

**Checklists a crear:**
- `CHECKLIST_NEGATIVE_BALANCES.md` — Diagnosticar saldo negativo
- `CHECKLIST_DUPLICATES.md` — 8 heurísticas antes de borrar
- `CHECKLIST_WARNINGS.md` — Cada warning de CT: qué significa, qué hacer
- `CHECKLIST_FIFO_VALIDATION.md` — Validar cost basis paso a paso
- `CHECKLIST_AUDIT_COMPLETE.md` — ¿Está lista la auditoría para reportar?
- `CHECKLIST_TAX_DECLARATION.md` — ¿Está lista la declaración para presentar?

**Próximo paso:** Crear en Fase 3 (compilar desde skills)

---

### D2 — Árboles de Decisión

**Ubicación propuesta:** `knowledge/decision-trees/`

**Estado actual:**
- ❌ No existen formalmente

**Árboles a crear:**
- `FLOW_AUDIT.md` — Diagrama completo de auditoría (6 fases)
- `FLOW_DUPLICATE_DETECTION.md` — ¿Misma fecha? → ¿TX ID?
- `FLOW_TRANSFER_MATCHING.md` — ¿Withdrawal? → ¿Blockchain? → ¿Deposit?
- `FLOW_MISSING_PURCHASE_HISTORY.md` — ¿Falta origen? → tipo
- `FLOW_FISCAL_DECISION.md` — ¿Conflicto fiscal? → A/B/C (ADR-026)
- `FLOW_RESOLVE_NEGATIVE_BALANCE.md` — Paso a paso para resolver

**Próximo paso:** Crear en Fase 3 (formalizar lógica de skills)

---

### D3 — Índices y Referencias

**Ubicación propuesta:** `knowledge/reference/`, `knowledge/.metadata/`

**Estado actual:**
- ✅ `cointracking/reference/CATALOG.md` (205 artículos oficiales CT indexados)
- ✅ `docs/GLOSSARY.md` (términos básicos) — *mover a E1*
- ❌ `knowledge/INDEX_MASTER.md` (ESTE DOCUMENTO)
- ❌ `adr/INDEX.md` (mapa de ADRs)
- ❌ `knowledge/.metadata/AUTHORITY_MATRIX.md` (quién verifica qué)
- ❌ `knowledge/.metadata/DEPRECATION_LOG.md` (documentos obsoletos)
- ❌ `knowledge/.metadata/VIGENCIA_ALERTS.md` (próximos a expirar)

**Próximo paso:** Crear en Fase 2

---

## NIVEL E — Referencia

**Ubicación:** `knowledge/reference/`

**Características:**
- ✅ Contexto, definiciones, historiadores
- ✅ Raramente envejecen
- ✅ `valid_until: null` permitido

### E1 — Glosario

**Ubicación propuesta:** `knowledge/reference/GLOSSARY.md`

**Estado actual:**
- ✅ `docs/GLOSSARY.md` (existe)

**Próximo paso:** Migrar a E1 en Fase 2

**Términos cubiertos:**
- FIFO
- Purchase pool
- Base de coste
- Ganancia patrimonial vs rendimiento
- Staking vs rewards
- Permuta
- Modelo 721
- Etc.

---

### E2 — Contexto e Historiadores

**Ubicación propuesta:** `knowledge/reference/CONTEXT/`

**Estado actual:**
- ✅ `exchanges/BINANCE_EU_MICA_EXIT.md` (MiCA, salida de Binance UE 2026-07)
- ❌ `REGULATORY_TIMELINE.md` (hitos: MiCA entrada vigor, DAC8, etc.)
- ❌ `COINTRACKING_EVOLUTION.md` (cambios de CT a lo largo del tiempo)

**Próximo paso:** Crear en Fase 3

---

## NIVEL F — Governance

**Ubicación:** `adr/`, `knowledge/`, `knowledge/.metadata/`

**Características:**
- ✅ Gobiernan cómo funciona todo lo demás
- ✅ Son ADRs, índices, y metadatos del sistema

### F1 — Decisiones Arquitectónicas (ADRs)

**Ubicación:** `adr/0001-*.md` a `adr/0033-*.md`

**ADRs Importantes:**
- ADR-001 — Principios de auditoría
- ADR-002 — Fuente de verdad (jerarquía)
- ADR-003 — Modelo de transacciones
- ADR-008 — Vigencia del conocimiento
- ADR-009 — Protocolo crítico
- ADR-026 — Matriz de decisiones (A/B/C)
- **ADR-031** — Validación de plazos (fecha)
- **ADR-032** — Knowledge with Temporal Validity (metadatos)
- **ADR-033** — Sistema de Conocimiento Jerárquico (este)

**Próximo paso:** Crear `adr/INDEX.md` en Fase 2 (mapea ADR → nivel de conocimiento)

---

### F2 — Índices Maestros

**Ubicación:** Múltiples ubicaciones

**Documentos:**
- `knowledge/INDEX_MASTER.md` — ESTE DOCUMENTO
- `knowledge/cointracking/INDEX.md` — Índice de documentos CT (ya existe)
- `knowledge/patterns/INDEX.md` — Índice de casos (ya existe)
- `knowledge/taxation/spain/INDEX.md` — Índice de fiscal (pendiente)
- `adr/INDEX.md` — Índice de ADRs (pendiente)

**Próximo paso:** Crear `adr/INDEX.md` en Fase 2

---

### F3 — Metadatos del Sistema

**Ubicación:** `knowledge/.metadata/`

**Documentos:**
- ✅ `knowledge/.metadata/METADATA_TEMPLATE.md` (esquema YAML oficial) — CREADO HOY
- ❌ `knowledge/.metadata/SCHEMA.yaml` (esquema formal, validable)
- ❌ `knowledge/.metadata/AUTHORITY_MATRIX.md` (quién verifica cada nivel)
- ❌ `knowledge/.metadata/DEPRECATION_LOG.md` (documentos obsoletos/retirados)
- ❌ `knowledge/.metadata/VIGENCIA_ALERTS.md` (próximos a expirar)

**Próximo paso:** Crear en Fase 3

---

## Brechas Actuales (Estado vs. Propuesto)

| Nivel | Cobertura | Estado | Crítico |
|-------|-----------|--------|---------|
| **A1** | España (oficial) | 60% | ⭐⭐⭐ |
| **A2** | CoinTracking (oficial) | 70% | ⭐⭐⭐ |
| **A3** | Exchanges (oficial) | 20% | ⭐⭐ |
| **B1** | Operativo CT | 50% | ⭐⭐⭐ |
| **B2** | Operativo exchanges | 30% | ⭐⭐ |
| **B3** | Blockchain | 10% | ⭐ |
| **C1** | Casos auditados | 100% (v2 YAML) | ⭐⭐⭐ |
| **C2** | Patrones | 0% | ⭐⭐ |
| **C3** | Procedimientos | 0% | ⭐⭐ |
| **D1** | Checklists | 0% | ⭐⭐ |
| **D2** | Árboles decisión | 0% | ⭐⭐ |
| **D3** | Índices | 60% | ⭐ |
| **E1** | Glosario | 100% | ⭐ |
| **E2** | Historiadores | 20% | ⭐ |
| **F1** | ADRs | 100% | ⭐⭐⭐ |
| **F2** | Índices maestros | 30% | ⭐⭐ |
| **F3** | Metadatos | 20% | ⭐⭐ |

---

## Plan de Migración (3 Fases)

### ✅ Fase 1 (HOY — 2026-07-05)

**Formalización:**
- ✅ ADR-033 creado
- ✅ `knowledge/.metadata/METADATA_TEMPLATE.md` creado
- ✅ `knowledge/INDEX_MASTER.md` creado (ESTE)
- ✅ Estructura vacía de directorios
- ✅ Plan documentado

**Commits:**
- Uno: Crear ADR-033 + plantilla + INDEX_MASTER

---

### ⏭️ Fase 2 (Próxima sesión — 2026-07-06 o después)

**Reorganización de directorios:**
- Mover `cointracking/CSV_FORMAT.md`, `COST_BASIS.md` → `cointracking/official/`
- Mover/crear `BINANCE.md` → `exchanges/official/`
- Crear `adr/INDEX.md`
- Crear `knowledge/.metadata/AUTHORITY_MATRIX.md`

**Actualizar referencias:**
- En skills (`/audit-cointracking`, `/spanish-tax-return`)
- En docs (links a documentos reubicados)
- En ADRs (si alguno referencia rutas)

**Commits:**
- Uno por carpeta principal reorganizada

---

### ⏭️ Fase 3 (Iterativo — 2026-07-07+)

**Agregar metadatos YAML:**
- Actualizar documentos existentes con frontmatter

**Crear documentos nuevos:**
- Procedures (C3)
- Checklists (D1)
- Árboles de decisión (D2)
- Patrones (C2)

**Automatizar:**
- Script de validación de metadatos
- Alertas de vigencia

**Commits:**
- Uno por cada documento o grupo cohesivo

---

## Navegación Rápida

| Necesito... | Voy a... |
|---|---|
| Saber qué es FIFO | E1 - Glosario |
| Entender purchase pool | B1 - Operativo CT o C1 - Casos |
| Resolver saldo negativo | D2 - Árbol de decisión (FLOW) |
| Validar un duplicado | D1 - Checklist |
| Citar fuente oficial AEAT | A1 - Autoridades españolas |
| Ver cómo funciona Binance | B2 - Operativa exchanges |
| Entender un caso real | C1 - Casos auditados |
| Buscar síntoma ("missing PH") | Troubleshooting → C1 |
| Conocer ADRs relacionados | F1 - ADRs + `adr/INDEX.md` |

---

## Próximos Pasos

1. **Hoy (Fase 1):** Revisar y aprobar este documento
2. **Mañana (Fase 2):** Reorganizar directorios + crear índices
3. **Semana 2 (Fase 3):** Agregar metadatos + crear documentos nuevos

---

## Contacto y Mantenimiento

**Responsable actual:** Alfredo González P. + Claude Code (gestión)

**Política de actualización:**
- Cuando cambien normativas fiscales (España): actualizar A1, notificar en `knowledge/.metadata/VIGENCIA_ALERTS.md`
- Cuando cambien features de CT: actualizar A2/B1
- Cuando se audite un nuevo caso: añadir a C1
- Cuando se identifique un patrón nuevo: documentar en C2

---

## Validación

Este documento fue creado como parte de **ADR-033: Sistema de Conocimiento Jerárquico** (2026-07-05).

Métadatos propios:

```yaml
id: F2-MASTER-001
title: "Sistema de Conocimiento — Mapa Maestro"
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
  - ADR-032
  - ADR-031
tags:
  - architecture
  - governance
  - knowledge-system
notes: "Estructura permanente. Actualizaciones: cuando cambien niveles A-E."
```
