---
id: KB-F1-001
title: "Guía Rápida (5 minutos)"
level: F
domain: cointracking
source: "Documentación interna"
authority: reference
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - navigation
  - reference

notes: "Documento de navegación/referencia de la base de conocimiento"
---

## 🚀 En 3 Pasos

### 1. ¿Tienes un problema específico?
→ **Ve a [Índice de Troubleshooting](TROUBLESHOOTING_INDEX.md)**

Busca el síntoma (ej: "saldo negativo", "duplicados", "missing cost basis") y sigue el árbol de decisión.

### 2. ¿Necesitas entender cómo funciona CoinTracking?
→ **Ve a [Comportamiento de CoinTracking](cointracking/behavioral/INDEX.md)**

- `DUPLICATE_DETECTION_HEURISTICS.md` — Cómo CoinTracking decide si son duplicados
- `PURCHASE_POOL_MECHANICS.md` — Cómo funciona FIFO
- `API_VS_CSV_OVERLAP.md` — Qué pasa si importas dos veces

### 3. ¿Necesitas verificar algo contra datos reales?
→ **Ve a [Casos Verificados](cases/INDEX.md)**

20 casos auditados de verdad. Busca uno similar al tuyo → diagnostica el problema.

---


# Guía de Inicio Rápido — Sistema de Auditoría de CoinTracking

**Para:** Usuarios nuevos del agente auditor  
**Tiempo:** 5 minutos  
**Objetivo:** Entender dónde buscar información



## 🗺️ El Sistema en 6 Niveles

```
┌─ NIVEL A (Oficial) ──────────────────┐  ← Fuentes: AEAT, CoinTracking, Exchanges
│ ¿Qué dice la ley/herramienta?        │
├─ NIVEL B (Operativo) ────────────────┤  ← Cómo funcionan las cosas
│ ¿Cómo lo hace CoinTracking?          │
├─ NIVEL C (Verificado) ────────────────┤  ← Casos reales
│ ¿Pasó esto a alguien ya?             │
├─ NIVEL D (Auxiliar) ──────────────────┤  ← Listas y árboles
│ ¿Qué checklist sigo? ¿Qué flujo?     │
├─ NIVEL E (Referencia) ────────────────┤  ← Glosario e historia
│ ¿Qué significa este término?         │
└─ NIVEL F (Governance) ────────────────┘  ← Decisiones registradas
  ¿Cómo se decidió esto?
```

**Orden natural de búsqueda:**
1. ¿Hay un caso igual? → Nivel C
2. ¿Cómo funciona? → Nivel B
3. ¿Qué dice la ley/herramienta? → Nivel A

---

## ⚡ Búsquedas Comunes

| Pregunta | Dónde ir | Documento |
|----------|----------|-----------|
| "Tengo duplicados" | Nivel D | [CHECKLIST_DUPLICATES.md](checklists/CHECKLIST_DUPLICATES.md) |
| "Mi balance es negativo" | Nivel D | [FLOW_NEGATIVE_BALANCE.md](decision-trees/FLOW_NEGATIVE_BALANCE.md) |
| "Falta el origen de una compra" | Nivel C | [CT-002](cases/ct-002-venta-sin-historial-de-compra-previo-mis.md) |
| "API y CSV importados juntos" | Nivel B | [API_VS_CSV_OVERLAP.md](cointracking/behavioral/API_VS_CSV_OVERLAP.md) |
| "¿Cómo funciona el FIFO?" | Nivel B | [PURCHASE_POOL_MECHANICS.md](cointracking/behavioral/PURCHASE_POOL_MECHANICS.md) |
| "¿Qué es una ganancia patrimonial?" | Nivel E | [GLOSSARY.md](reference/GLOSSARY.md) |
| "¿Cómo fue la decisión X?" | Nivel F | [adr/INDEX.md](../adr/INDEX.md) |

---

## 📚 Documentos Esenciales (Primero)

Léelos en este orden si es tu primera auditoría:

1. **[CSV_FORMAT.md](cointracking/official/CSV_FORMAT.md)** (10 min)  
   → Entiende la estructura de los datos de CoinTracking

2. **[PURCHASE_POOL_MECHANICS.md](cointracking/behavioral/PURCHASE_POOL_MECHANICS.md)** (15 min)  
   → Entiende cómo CoinTracking calcula cost basis

3. **[PROCEDURE_AUDIT_ACCOUNT.md](procedures/PROCEDURE_AUDIT_ACCOUNT.md)** (20 min)  
   → Pasos de auditoría paso a paso

4. **[CAPITAL_GAINS.md](taxation/spain/CAPITAL_GAINS.md)** (10 min)  
   → Cómo se declara fiscalmente

---

## 🔧 Si Algo No Funciona

### "No encuentro el documento que necesito"
→ Ve a [NAVIGATION_MAP.md](NAVIGATION_MAP.md) (índice por función, no por tema)

### "Tengo un problema que no aparece aquí"
→ Ve a [TROUBLESHOOTING_INDEX.md](TROUBLESHOOTING_INDEX.md) (lista de síntomas)

### "¿Esto está desactualizado?"
→ Busca `last_verified` en el documento (en el encabezado YAML)  
   Si dice 2025 o antes, marca como `[VERIFICAR]` antes de usar

---

## 💾 Una Nota Importante

**Este sistema vive en dos lugares:**

- **`knowledge/`** — Base de conocimiento (autoridad, operativo, casos, referencias)
- **`adr/`** — Decisiones arquitectónicas (por qué se decidió así)

Ambos están versionados en git. Las actualizaciones de conocimiento requieren commit (con evidencia y razonamiento).

---

## ❓ Preguntas Frecuentes

**P: ¿Cuál es la diferencia entre Nivel C1, C2 y C3?**  
R: C1 = casos reales específicos (CT-001 a CT-020). C2 = patrones generalizados (qué hace duplicado). C3 = procedimientos paso a paso (cómo auditar).

**P: ¿Puedo confiar en las cifras del agente?**  
R: Diagnóstico sí (encontrar problemas). Cifras fiscales exactas: se necesita verificación humana contra fuente oficial (AEAT, DGT).

**P: ¿El sistema cubre mi exchange?**  
R: Si está en Binance, Kraken, Coinbase: documentación completa. Otros: [GENERIC_EXCHANGE_MECHANICS.md](cointracking/behavioral/GENERIC_EXCHANGE_MECHANICS.md).

**P: ¿Qué pasa con las criptos que CoinTracking no conoce?**  
R: Búscalas en [CATALOG.md](cointracking/reference/CATALOG.md) (205+ criptos indexadas). Si no aparece, abre issue.

---

## 🚪 Siguiente

- **Para auditar tu cuenta:** Skill `/audit-cointracking`
- **Para preparar impuestos:** Skill `/spanish-tax-return`
- **Para explorar más:** Ve a [Mapa Maestro](INDEX_MASTER.md)
