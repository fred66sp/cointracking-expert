# Registros de decisiones arquitectónicas

**Decisiones arquitectónicas importantes del proyecto CoinTracking Expert**

Este archivo documenta decisiones arquitectónicas significativas usando el formato ADR (Architecture Decision Record). Cada decisión incluye el contexto, opciones consideradas, decisión tomada y consecuencias.

---

## ADR-001: Idioma del repositorio es español

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

El proyecto tiene como objetivo servir principalmente a usuarios hispanohablantes con enfoque en cumplimiento fiscal español. El equipo de desarrollo también es hispanohablante.

**Opciones consideradas:**

1. **Inglés**: Estándar de la industria, comunidad global más grande
2. **Español**: Accesibilidad para usuarios españoles, claridad para el equipo local

**Decisión:**

Todos los documentos, código y comunicaciones del proyecto serán en español para maximizar la accesibilidad para el público objetivo.

**Consecuencias:**

- ✅ Mejor accesibilidad para usuarios hispanohablantes
- ✅ Mayor claridad para el equipo de desarrollo local
- ⚠️ Menor comunidad potencial de contribuidores globales
- ⚠️ Necesidad de traducción para documentación técnica estándar

---

## ADR-002: Stack de tecnología Python (Pendiente)

**Estado:** Pendiente de completar

**Contexto:**

Se requiere seleccionar un stack de tecnología Python para la implementación. Necesitamos decisiones sobre:
- Versión mínima de Python
- Librería de validación (pydantic, dataclasses, attrs)
- Base de datos (SQLite, PostgreSQL, en memoria)
- Framework web (FastAPI, Flask, Django) para API futura

**Opciones siendo evaluadas:**

1. **Pydantic v2 + SQLAlchemy + FastAPI**: Moderno, bien mantenido, estándar de industria
2. **Dataclasses + SQLite + Flask**: Más ligero, menos dependencias
3. **Attrs + PostgreSQL + Starlette**: Equilibrio entre características y simplicidad

**Decisión:**

Pendiente de consenso del equipo.

**Próximos pasos:**

- Evaluar cada opción con pruebas de prototipo
- Documentar criterios de decisión
- Completar esta ADR

---

## ADR-003: Traducción del modelo de dominio (Pendiente)

**Estado:** Pendiente de completar

**Contexto:**

El DOMAIN_MODEL.md está escrito en español con pseudocódigo Kotlin. Necesita ser traducido a Python real.

**Decisión pendiente:**

- ¿Usar dataclasses de Python?
- ¿Usar pydantic v2?
- ¿Usar attrs?
- ¿Usar protocolos de Python?

**Próximos pasos:**

- Decidir sobre ADR-002 primero
- Implementar modelos Python basado en esa decisión
- Validar contra casos de uso reales

---

## Plantilla para futuros ADRs

```
## ADR-###: Título de la decisión

**Estado:** Propuesto / Pendiente / Decidido / Rechazado

**Fecha:** YYYY-MM-DD

**Contexto:**

[Describe el problema y por qué es importante...]

**Opciones consideradas:**

1. Opción A: [Descripción]
2. Opción B: [Descripción]
3. Opción C: [Descripción]

**Decisión:**

[Describe cuál fue elegida y por qué...]

**Consecuencias:**

- ✅ Beneficios
- ⚠️ Ventajas
- ❌ Riesgos o desventajas

**Notas adicionales:**

[Información relevante adicional...]
```

---

## Índice de ADRs

- ADR-001: Idioma del repositorio es español ✅ Decidido
- ADR-002: Stack de tecnología Python ⏳ Pendiente
- ADR-003: Traducción del modelo de dominio ⏳ Pendiente

---

## Proceso de ADR

Toda decisión arquitectónica importante debe:

1. Ser propuesta en una rama nueva con un ADR borrador
2. Ser discutida en revisión de código
3. Ser revisada por arquitecto del proyecto
4. Ser aprobada por el equipo
5. Ser completada con la decisión final

Las decisiones menores pueden ser documentadas informalmente en CHANGELOG.md.
