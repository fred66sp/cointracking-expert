# ADR-021: Gate explícito de artefacto antes de cerrar la cifra anual exacta en `spanish-tax-return`

**Status:** Accepted

**Date:** 2026-07-03

## Context

Petición de Copilot (explotación, ADR-012) tras preparar la renta 2025 de `agp2025`: pudo cerrar la clasificación de eventos, rendimientos y derivados sin el Tax Report oficial del ejercicio, pero llegó al final del informe sin ese artefacto y sin una regla explícita que le impidiera marcar el documento como "listo para presentar" pese a faltar la cifra anual exacta de la base del ahorro. El playbook dejaba el cierre "bloqueado" implícitamente (por falta de datos), pero no exigía comprobar la presencia del artefacto antes de declarar el informe cerrado.

**Decisión:**

Añadir un gate explícito en `.claude/skills/spanish-tax-return/SKILL.md`:

1. **Paso 3:** si no está en el workspace el Tax Report oficial del ejercicio (o su cifra `Resumen` ya documentada con evidencia), no cerrar la cifra de base del ahorro ni marcarla como definitiva — declararla `[VERIFICAR]` y pedir el artefacto al usuario. El resto de secciones (eventos, rendimientos, Modelo 721) sí pueden avanzar sin ese artefacto.
2. **Paso 6:** recordatorio explícito de que el informe solo se marca "listo para presentar" si ese gate está satisfecho.

No se añade un chequeo de código (`tools/ct_audit.py` no interviene aquí, es un gate de proceso/redacción, no un chequeo mecánico sobre datos).

## Decision

[Decision not found]

## Consequences

- ✅ Cierra el caso real: en la misma sesión, el usuario aportó el Tax Report 2025 y se completó el informe (`reports/output/agp2025/2026-07-03_declaracion_2025.md` §7, "listo para presentar").
- ✅ Evita que un futuro informe se declare cerrado por omisión cuando en realidad falta el artefacto determinante.
- ✅ Entrada `AGENT_CHANGE_REQUESTS.md` del 2026-07-03 ("Precondición explícita de artefacto...") marcada como hecha.
