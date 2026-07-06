---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-042: Proactividad en gobernanza — el agente sugiere el ADR, el usuario aprueba

**Status:** Accepted

**Date:** 2026-07-05

## Context

**Caso que lo motivó (2026-07-05):** durante la limpieza del proyecto legado `demo` apareció un hallazgo operativo con entidad propia (7 servidores MCP huérfanos bloqueando cachés — ADR-041). El agente lo diagnosticó, lo resolvió con un protocolo cuidadoso… y lo dejó como anécdota en el chat. Fue **el usuario** quien tuvo que pedir "crea un ADR para que quede registrado el caso". El conocimiento se salvó, pero por iniciativa equivocada: la gobernanza no puede depender de que el usuario reconozca qué merece registro — para eso está el agente, que es quien vive el hallazgo con todo el detalle.

El patrón se había repetido antes: ADR-036/037/038 nacieron porque el usuario dijo "vamos a tener que crear ADRs también para el desarrollo"; ADR-014 (FLOKI) se documentó tras el incidente, no al detectarlo. La regla existente ("decisiones importantes → nuevo ADR", `CLAUDE.md` §Convenciones) es pasiva: dice *qué* registrar, no *quién* debe darse cuenta.

## Decision

**Detectar que algo merece ADR es responsabilidad del agente, no del usuario.** En cuanto ocurra un evento ADR-digno, el agente lo dice **en el momento** (no al final de la sesión, no "si surge"), pide permiso y, con el visto bueno, lo crea.

### Disparadores (cuándo el agente DEBE sugerir un ADR)

1. **Hallazgo operativo repetible con diagnóstico no obvio** — un síntoma que volverá y cuya causa costó derivar (ej.: huérfanos MCP → ADR-041; `get_historical_summary` fuera de rango → ADR-020).
2. **Decisión de diseño o arquitectura tomada en conversación** — cualquier cambio de comportamiento del sistema acordado por chat que no sea trivial (ej.: credenciales por proyecto → ADR-040; TTL/versionado de caché → ADR-039).
3. **Causa raíz de un bug que revela una regla vinculante** — cuando el fix enseña algo que debe gobernar el futuro, no solo arreglar el presente (ej.: duplicados con mismo timestamp → ADR-014).
4. **Corrección del usuario que redefine cómo trabajar** — cuando el usuario endereza el enfoque del agente y esa corrección debe sobrevivir a la sesión (ej.: proyecto activo obligatorio → ADR-013; este mismo ADR).
5. **Protocolo improvisado con éxito** — si el agente tuvo que inventar un procedimiento paso a paso para salir de una situación y funcionó, ese procedimiento es candidato a quedar escrito.

### Flujo obligatorio

1. **Detectar y decir en el momento:** al cerrar el evento (fix verificado, decisión tomada, protocolo ejecutado), el agente añade a su respuesta una propuesta breve: qué registraría, por qué merece ADR y qué disparador aplica.
2. **Pedir permiso (Categoría B de ADR-026):** crear un ADR es acción con consecuencia de gobernanza — requiere consentimiento explícito. Con `AskUserQuestion` si hay opciones reales (crear ya / anotar para luego / no merece); en texto si es un sí/no simple.
3. **Con el visto bueno, crear el ADR completo en la misma sesión:** contenido según ADR-030 (clasificación de criticidad y su checklist), frontmatter con `version:` (ADR-039), actualización de **ambos** índices (`adr/README.md` y `adr/INDEX.md`) y entrada en `CHANGELOG.md`.
4. **Si el usuario declina, no insistir** — pero dejar una línea en el CHANGELOG o en la memoria de sesión indicando que el evento ocurrió y se decidió no registrarlo, para que la decisión de no documentar también sea trazable.

### Guardarraíl contra el ruido (proporcionalidad)

No todo merece ADR — proponer de más devalúa la gobernanza tanto como proponer de menos. **No** se sugiere ADR para: fixes rutinarios sin lección general, decisiones ya cubiertas por un ADR existente (ahí se propone *actualizar* ese ADR, no crear otro), preferencias de estilo (van a memoria), ni detalles de implementación sin decisión de fondo. Regla práctica: si el evento no encaja en ninguno de los 5 disparadores, no se propone.

## Consequences

- ✅ La gobernanza deja de depender de que el usuario reconozca qué merece registro; el agente, que vive el hallazgo con detalle completo, es quien lo señala.
- ✅ El consentimiento se mantiene (nada se crea sin permiso) — cambia la **iniciativa**, no la autoridad.
- ✅ La decisión de *no* documentar algo también queda trazada.
- ⚠️ Riesgo de sobre-proponer: mitigado por los 5 disparadores cerrados y el guardarraíl de proporcionalidad.
- ⚠️ Exige disciplina del agente en mitad del trabajo (proponer al cerrar el evento, no dejarlo para un "luego" que no llega).

## Notes

- Relación: ADR-026 (proponer es Categoría A — el agente lo hace solo; crear es Categoría B — requiere confirmación), ADR-030 (checklist de validación del ADR resultante), ADR-009 §7 (consentimiento informado).
- Operacionalizado en `CLAUDE.md` §Convenciones (que se carga en cada sesión) y en la memoria durable del agente — este ADR es la fuente; esos dos son los recordatorios activos.
