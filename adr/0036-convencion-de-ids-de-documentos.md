---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-036: Convención de IDs de Documentos de Conocimiento

**Status:** Accepted  
**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)  
**Proposed:** 2026-07-05  
**Accepted:** 2026-07-05  
**Last Updated:** 2026-07-05

---

## Problema

No hay una convención clara y documentada para los IDs de documentos en la base de conocimiento. Esto causa:
- Validadores inciertos sobre qué patrón aceptar
- Documentos creados con IDs inconsistentes
- Auditorías que reportan "inválidos" documentos que están en el formato correcto

**Ejemplo:** `KB-B4-001` (Trezor Integration) vs `KB-B-NNN` (genérico esperado)
- Script de auditoría rechaza como "INVALID"
- Pero 91 documentos usan este patrón
- Pattern es claramente intencional, no error

---

## Decisión

**La convención oficial de IDs es:**

```
KB-[NIVEL][SUBSECCIÓN]-[NÚMERO]
```

Donde:

### NIVEL (obligatorio)
Uno de: `A`, `B`, `C`, `D`, `E`, `F` (ver ADR-033)

### SUBSECCIÓN (obligatorio)
Número entero `1–9` que agrupa documentos por tema dentro del nivel:

#### Nivel A (Official / Oficial)
- **A1:** AEAT/DGT/BOE (fiscal España oficial)
- **A2:** CoinTracking official (guía oficial CT, formatos, APIs)

#### Nivel B (Operational / Operacional Verificado)
- **B1:** CoinTracking internals (CSV, cost basis, importación, validación)
- **B2:** Exchanges (Binance, Kraken, Coinbase, Bybit, OKX, BingX, etc.)
- **B3:** Blockchains (Ethereum, Bitcoin, transacciones, fees, bridges)
- **B4:** Wallets (Ledger, MetaMask, Trezor, Trust Wallet)

#### Nivel C (Cases / Casos y Patrones)
- **C1:** Auditoría cases reales (ct-001 a ct-020)
- **C2:** Patrones de reconciliación (balance, duplicates, transfers, pools)
- **C3:** Procedimientos (audit account, fix missing history, reconcile transfers)

#### Nivel D (Auxiliary / Auxiliar)
- **D1:** (reservado)
- **D2:** Decision trees (audit flow, duplicate detection, negative balance)

#### Nivel E (Reference / Referencia)
- **E1:** Reference (glossary, catalogs)
- **E2:** Reference (project history, navigation maps)

#### Nivel F (Governance / Gobernanza)
- **F1:** ADRs, decisions, governance
- **F2:** Status, roadmap, changelog

### NÚMERO (obligatorio)
Número secuencial `001–999` dentro de la subsección:
- Ejemplo: `KB-B1-010` = décimo documento de B1
- Ejemplo: `KB-A1-001` = primer documento de A1

---

## Ejemplos Válidos

| ID | Documento | Nivel | Subsección | Número |
|---|---|---|---|---|
| `KB-A1-001` | INFORMATIVE_OBLIGATIONS.md | A | 1 | 001 |
| `KB-A2-001` | CSV_FORMAT.md | A | 2 | 001 |
| `KB-B1-010` | DUPLICATE_EDGE_CASES.md | B | 1 | 010 |
| `KB-B4-003` | TREZOR_INTEGRATION.md | B | 4 | 003 |
| `KB-C1-001` | ct-001-transferencia...md | C | 1 | 001 |
| `KB-C2-001` | PATTERN_DUPLICATE_DETECTION.md | C | 2 | 001 |
| `KB-C3-001` | PROCEDURE_AUDIT_ACCOUNT.md | C | 3 | 001 |
| `KB-D2-001` | FLOW_DUPLICATE_DETECTION.md | D | 2 | 001 |
| `KB-E1-001` | GLOSSARY.md | E | 1 | 001 |

---

## Ejemplos Inválidos

| ID | Razón |
|---|---|
| `KB-B-001` | Falta SUBSECCIÓN |
| `KB-B1` | Falta NÚMERO |
| `KB-B15-001` | SUBSECCIÓN fuera de rango (>9) |
| `KB-B1-AAA` | NÚMERO no numérico |
| `KB-X1-001` | NIVEL inválido (no es A–F) |

---

## Genéricos para Ejemplos/Templates

Para archivos que son **documentación sobre cómo crear documentos** (no documentos reales), usar:
```
KB-[NIVEL]-NNN
```

Ejemplo:
- `METADATA_TEMPLATE.md` → `KB-A-NNN`, `KB-B-NNN` (ejemplos genéricos)
- `MIGRATION_PLAN.md` → `KB-C-NNN` (ejemplo genérico)

---

## Implicaciones

1. **Validador de auditoría** debe aceptar patrón `KB-[A-F][1-9]-\d{3}`
2. **Documentos con IDs genéricos** (KB-X-NNN) son válidos solo en `.metadata/` o plantillas
3. **Ningún documento de conocimiento real** debe tener ID genérico
4. **Subsección determina agrupación temática**, no convención de creación

---

## Referencias

- ADR-033: Sistema de conocimiento jerárquico (6 niveles A-F)
- ADR-032: Knowledge with Temporal Validity
- `knowledge/INDEX_MASTER.md`: Índice que agrupa por subsección

---

**Decisión:** ACEPTADA (deduce de uso actual, documentada para futuro)
