# Árboles de Decisión (Nivel D2)

**Ubicación:** `knowledge/decision-trees/`

**Característica:** Máquinas de estado y flujos de decisión.

**Autoridad:** `verified` — basadas en lógica operativa probada

---

## Status

Esta carpeta está **vacía por ahora** (Fase 3).

---

## Árboles Pendientes (Fase 3)

- `FLOW_AUDIT.md` — Diagrama completo de auditoría (6 fases)
- `FLOW_DUPLICATE_DETECTION.md` — ¿Duplicado? (matriz de heurísticas)
- `FLOW_TRANSFER_MATCHING.md` — ¿Transfer huérfana? (emparejar)
- `FLOW_MISSING_PURCHASE_HISTORY.md` — ¿Falta origen? (diagnosticar tipo)
- `FLOW_FISCAL_DECISION.md` — ¿Conflicto fiscal? (A/B/C, ADR-026)
- `FLOW_RESOLVE_NEGATIVE_BALANCE.md` — ¿Saldo negativo? (resolver)

---

## Formato

Cada árbol:
- Diagramas ASCII o Mermaid
- Máquinas de estado: pregunta → rama → acción
- Estados terminales: "OK", "Requiere acción manual", "Escalar"

---

## Relación con Checklists (D1)

- **D2 (Árboles):** Lógica de decisión ("si X entonces Y")
- **D1 (Checklists):** Validación paso a paso ("verificar items 1-5")

Los árboles **deciden**, los checklists **verifican**.

---

## Próximas Sesiones

**Fase 3:** Crear árboles formalizando lógica de skills y playbooks.
