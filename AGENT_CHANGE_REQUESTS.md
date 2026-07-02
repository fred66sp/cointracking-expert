# Peticiones de cambio al agente

Bandeja de entrada para mejoras del **agente** (código, conocimiento, reglas, tool, skills).

- **Copilot (explotación):** si durante el uso detectas un bug, un hueco de conocimiento o una regla a cambiar, **NO lo edites** — **añade una entrada aquí** (append) y sigue. (ADR-012)
- **Claude Code (gestión):** procesa estas peticiones, aplica el cambio con gobernanza (ADR/commit) y marca la entrada como ✅ hecha.

Formato de entrada:
```
## [PENDIENTE] AAAA-MM-DD — Título breve
- **Qué:** qué falla o falta.
- **Dónde:** fichero/regla afectada (p. ej. knowledge/…, tools/ct_audit.py).
- **Evidencia:** dato/caso que lo motiva.
- **Propuesta:** (opcional) qué cambio se sugiere.
```

---

<!-- Añade nuevas peticiones debajo de esta línea -->
