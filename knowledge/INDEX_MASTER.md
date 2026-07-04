# Sistema de Conocimiento — Mapa Maestro

**Documento:** Descripción general de la arquitectura jerárquica de conocimiento del agente  
**Validez:** Permanente (es estructura, no contenido)  
**Última actualización:** 2026-07-05 (Validación P0 + Navegabilidad P1 completadas)  
**Arquitectura:** Definida en ADR-033

---

## 🧭 ATAJOS DE NAVEGACIÓN (Lee Primero)

**Si es tu primer viaje:**
- 👉 [QUICK_START.md](QUICK_START.md) (5 min)

**Si buscas algo específico:**
- 👉 [NAVIGATION_MAP.md](NAVIGATION_MAP.md) (por función/pregunta)
- 👉 [TROUBLESHOOTING_INDEX.md](TROUBLESHOOTING_INDEX.md) (por síntoma)
- 👉 [CHEAT_SHEET.md](CHEAT_SHEET.md) (referencia rápida)

**Si quieres el mapa completo:**
- 👇 Continúa leyendo este documento

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

**Ubicación:** `knowledge/cointracking/behavioral/`

**Estado actual (10/12 documentos) ✅ 83%**
- ✅ `cointracking/COST_BASIS_AND_VALIDATION.md` (purchase pool, negativos, etc.) — *nota: parcialmente B1, parcialmente A2*
- ✅ `cointracking/TROUBLESHOOTING.md` (índice de síntomas)
- ✅ `STAKING_MECHANICS.md` (KB-B1-001) — Depósitos vs rewards, clasificación fiscal
- ✅ `AIRDROPS_MECHANICS.md` (KB-B1-002) — Importación, clasificación, tratamiento fiscal
- ✅ `LENDING_MECHANICS.md` (KB-B1-003) — Depósitos, intereses, diferencia principal
- ✅ `DEFI_SWAPS_MECHANICS.md` (KB-B1-004) — On-chain swaps, slippage, fees
- ✅ `PURCHASE_POOL_MECHANICS.md` (KB-B1-005) — FIFO, pool exhaustion, validación
- ✅ `DUPLICATE_DETECTION_HEURISTICS.md` (KB-B1-006) — Matriz de decisión, batching, Trade ID
- ✅ `MISSING_PURCHASE_HISTORY_CAUSES.md` (KB-B1-007) — 5 causas, diagnóstico, soluciones
- ✅ `BALANCE_CALCULATION_ALGORITHM.md` (KB-B1-008) — Suma acumulada, evolución, validación
- ✅ `API_VS_CSV_OVERLAP.md` (KB-B1-009) — Duplicados mixtos, limpieza, prevención
- ✅ `FEE_HANDLING.md` (KB-B1-010) — Comisiones en tercera moneda, cost basis, fiscal

**Próximo paso (Fase 4+):** Pasar a Nivel C (casos verificados)

---

### B2 — Operativas de Exchanges

**Ubicación:** `knowledge/cointracking/behavioral/`

**Estado actual (7/9 documentos) ✅ 78%**
- ✅ `cointracking/behavioral/BINANCE_SPOT_MECHANICS.md` (KB-B2-001) — Compras/ventas normales, API, importación
- ✅ `cointracking/behavioral/BINANCE_MARGIN_MECHANICS.md` (KB-B2-002) — Trading apalancado, borrow, repay, riesgos fiscales
- ✅ `cointracking/behavioral/BINANCE_FUTURES_MECHANICS.md` (KB-B2-003) — Perpetuos, PnL, funding fees, ⚠️ fiscalidad crítica
- ✅ `cointracking/behavioral/BINANCE_EARN_MECHANICS.md` (KB-B2-004) — Staking, APY, BETH, productos Earn
- ✅ `cointracking/behavioral/BINANCE_CONVERT_MECHANICS.md` (KB-B2-005) — Intercambios rápidos, fee integrado
- ✅ `cointracking/behavioral/KRAKEN_STAKING_MECHANICS.md` (KB-B2-006) — On-chain staking, rewards, kETH
- ✅ `cointracking/behavioral/COINBASE_ADVANCED_TRADE.md` (KB-B2-007) — Advanced Trade, staking, stETH, stUSDC
- ✅ `exchanges/BINANCE.md` (módulos, limitaciones) — *nota: parcialmente A3, parcialmente B2*
- ❌ `BINANCE_IMPORT_WORKFLOW.md` (cómo importar correctamente)

**Próximo paso (Fase 4+):** B3 (blockchains) y otros exchanges (Bybit, Kraken Advanced, etc.)

---

### B3 — Blockchain

**Ubicación:** `knowledge/blockchains/`

**Estado actual (6/6 documentos) ✅ 100% COMPLETO**
- ✅ `ETHEREUM_TRANSACTION_TYPES.md` (KB-B3-001) — Tipos de tx, contract calls, gas
- ✅ `BITCOIN_TRANSACTION_TYPES.md` (KB-B3-002) — UTXO model, consolidations, change outputs
- ✅ `OTHER_CHAINS_MECHANICS.md` (KB-B3-003) — Polygon, BSC, Solana, Arbitrum, multichain balance
- ✅ `BRIDGES_AND_WRAPPING.md` (KB-B3-004) — Bridges, wrapped tokens, multichain complexity
- ✅ `GAS_FEE_HANDLING.md` (KB-B3-005) — Registro, validación, tratamiento fiscal
- ✅ `STAKING_MECHANICS_BLOCKCHAIN.md` (KB-B3-006) — Validadores, rewards on-chain, slashing, PoS/DPoS

**Próximo paso (Fase 4+):** Pasar a Nivel C (casos verificados)

---

## NIVEL C — Empirismo Verificado

**Ubicación:** `knowledge/cases/`, `knowledge/patterns/`, `knowledge/procedures/`

**Características:**
- ✅ Basado en **casos reales auditados**, no teoría
- ✅ `authority: verified` (probado en proyecto real)
- ✅ Confianza alta, pero no es normativa
- ✅ Se pueden usar para diagnosticar problemas

### C1 — Casos Reales Auditados

**Ubicación:** `knowledge/cases/`

**Estado actual (20/20 documentos) ✅ 100% COMPLETO**
- ✅ `patterns/cointracking_casos_v2.yaml` — Fuente original (20 casos)
- ✅ Migración a archivos `.md` individuales completada
- ✅ Cada caso con KB-C1-001 a KB-C1-020
- ✅ Metadatos YAML completos (id, title, level, domain, authority, dates, confidence, version)

**Casos migrados (CT-001 a CT-020):**
- KB-C1-001: CT-001 — Transferencia entre exchanges importada solo en origen
- KB-C1-002: CT-002 — Venta sin historial de compra previo
- KB-C1-003: CT-003 — API y CSV importados simultáneamente (duplicado)
- KB-C1-004: CT-004 — Balance negativo por orden cronológico incorrecto
- KB-C1-005: CT-005 — Recompensas de staking clasificadas como depósito
- KB-C1-006: CT-006 — Binance Convert importado como venta/compra separadas
- KB-C1-007: CT-007 — Transferencia interna confundida con venta
- KB-C1-008: CT-008 — Duplicados aparentes por ejecución parcial
- KB-C1-009: CT-009 — Comisión (fee) omitida en importación
- KB-C1-010: CT-010 — Airdrop registrado como compra con coste artificial
- KB-C1-011: CT-011 — Lending tratado como transferencia genérica
- KB-C1-012: CT-012 — Balance negativo por importación parcial vía API
- KB-C1-013: CT-013 — Wallet externa no importada (fondos desaparecidos)
- KB-C1-014: CT-014 — Recompensas de minería registradas como depósito
- KB-C1-015: CT-015 — Swap DeFi fragmentado en varias operaciones on-chain
- KB-C1-016: CT-016 — Duplicados por reimportación completa del mismo periodo
- KB-C1-017: CT-017 — Coste cero por compra omitida de ejercicios anteriores
- KB-C1-018: CT-018 — Token renombrado interpretado como activo distinto
- KB-C1-019: CT-019 — Balance negativo tras eliminar compra confundida con duplicado
- KB-C1-020: CT-020 — Advertencia técnica interpretada como error fiscal

**Próximo paso (Fase 4+):** Crear C2 (patrones) y C3 (procedimientos) desde borradores existentes

---

### C2 — Patrones Recurrentes

**Ubicación:** `knowledge/patterns/`

**Estado actual (4/4 documentos) ✅ 100% COMPLETO**
- ✅ `PATTERN_DUPLICATE_DETECTION.md` (KB-C2-001) — Matriz: qué hace/no hace duplicado
- ✅ `PATTERN_BALANCE_RECONCILIATION.md` (KB-C2-002) — Cómo reconocer saldos inconsistentes
- ✅ `PATTERN_TRANSFER_MATCHING.md` (KB-C2-003) — Heurísticas para withdrawal/deposit
- ✅ `PATTERN_PURCHASE_POOL_EXHAUSTION.md` (KB-C2-004) — Síntomas de "purchase pool consumed"

**Próximo paso:** Integración con auditor

---

### C3 — Procedimientos Operativos

**Ubicación:** `knowledge/procedures/`

**Estado actual (3/3 documentos) ✅ 100% COMPLETO**
- ✅ `PROCEDURE_AUDIT_ACCOUNT.md` (KB-C3-001) — 6 fases de auditoría
- ✅ `PROCEDURE_RECONCILE_TRANSFERS.md` (KB-C3-002) — Emparejar withdrawal/deposit
- ✅ `PROCEDURE_FIX_MISSING_PURCHASE_HISTORY.md` (KB-C3-003) — Resolver origen

**Próximo paso:** Integración con auditor

---

## NIVEL D — Auxiliar

**Ubicación:** `knowledge/checklists/`, `knowledge/decision-trees/`

**Características:**
- ✅ Organizan conocimiento de Niveles A-C
- ✅ No aportan nuevo conocimiento, lo estructuran
- ✅ Los LLMs los siguen muy bien
- ✅ Reducen variabilidad de respuestas

### D1 — Checklists

**Ubicación:** `knowledge/checklists/`

**Estado actual (3/3 documentos) ✅ 100% COMPLETO**
- ✅ `CHECKLIST_DUPLICATES.md` (KB-D1-001) — 8 heurísticas antes de borrar
- ✅ `CHECKLIST_NEGATIVE_BALANCES.md` (KB-D1-002) — Diagnosticar saldo negativo
- ✅ `CHECKLIST_AUDIT_COMPLETE.md` (KB-D1-003) — ¿Está lista la auditoría para reportar?

**Próximo paso:** Integración con auditor

---

### D2 — Árboles de Decisión

**Ubicación:** `knowledge/decision-trees/`

**Estado actual (3/3 documentos) ✅ 100% COMPLETO**
- ✅ `FLOW_DUPLICATE_DETECTION.md` (KB-D2-001) — ¿Misma fecha? → ¿TX ID?
- ✅ `FLOW_NEGATIVE_BALANCE.md` (KB-D2-002) — ¿Balance negativo? Diagnóstico
- ✅ `FLOW_COMPLETE_AUDIT.md` (KB-D2-003) — Diagrama completo de auditoría (6 fases)

**Próximo paso:** Integración con auditor

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
- ✅ Glosario, contexto, historiadores
- ✅ `authority: verified` (fundamental pero no normativo)
- ✅ Cambia lentamente (generalmente estable)

### E1 — Glosario

**Estado actual (1/1 documentos) ✅ 100% COMPLETO**
- ✅ `GLOSSARY.md` (KB-E1-001) — 50+ términos técnicos y operativos

### E2 — Contexto e Historiadores

**Estado actual (1/1 documentos) ✅ 100% COMPLETO**
- ✅ `PROJECT_HISTORY.md` (KB-E2-001) — Evolución del proyecto, decisiones clave, roadmap

**Próximo paso:** Nivel F (Governance)

---

## NIVEL E — Referencia (Anterior)

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
