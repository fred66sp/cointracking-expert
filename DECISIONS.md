# Registros de decisiones arquitectónicas

**Decisiones arquitectónicas importantes del proyecto CoinTracking Expert**

Este archivo documenta decisiones arquitectónicas significativas usando el formato ADR (Architecture Decision Record). Cada decisión incluye el contexto, opciones consideradas, decisión tomada y consecuencias.

---

## ADR-001: Idioma del repositorio (contenido en español, identificadores en inglés)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

El proyecto tiene como objetivo servir principalmente a usuarios hispanohablantes con enfoque en cumplimiento fiscal español. El equipo de desarrollo también es hispanohablante. Al mismo tiempo, el código Python debe seguir convenciones universales (PEP 8) para mantenerse legible, buscable e interoperable con el ecosistema.

**Opciones consideradas:**

1. **Todo en inglés**: Estándar de la industria, comunidad global más grande
2. **Todo en español (documentos y código)**: Accesibilidad máxima, pero rompe convenciones de programación y dificulta búsquedas técnicas
3. **Híbrido**: Contenido en español, identificadores técnicos en inglés

**Decisión:**

Se adopta el modelo **híbrido**:

- **En español (contenido para humanos):**
  - Contenido de toda la documentación (`.md`)
  - Docstrings
  - Comentarios de código
  - Mensajes de error y de log dirigidos al usuario
- **En inglés (identificadores técnicos):**
  - Nombres de archivos y carpetas (`README.md`, `src/`, `engines/`)
  - Nombres de clases, funciones, métodos y variables (PEP 8)

**Consecuencias:**

- ✅ Documentación accesible para usuarios y equipo hispanohablante
- ✅ Código que respeta PEP 8 y es interoperable con el ecosistema Python
- ✅ Nombres de archivo estables y buscables (identificadores técnicos universales)
- ⚠️ Requiere disciplina para mantener la separación (contenido vs identificador)
- ⚠️ Menor comunidad potencial de contribuidores globales por la documentación en español

**Notas adicionales:**

Esta decisión es **permanente** para el proyecto y prevalece sobre cualquier documento que sugiera "todo en español".

---

## ADR-002: Stack de tecnología Python

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

Se requiere seleccionar un stack de tecnología Python para la implementación. Necesitamos decisiones sobre:
- Versión mínima de Python
- Librería de validación (pydantic, dataclasses, attrs)
- Tipo numérico para cantidades y precios (float vs Decimal)
- Base de datos (SQLite, PostgreSQL, en memoria)
- Framework web (FastAPI, Flask, Django) para API futura

**Opciones consideradas:**

1. **Pydantic v2 + Decimal + SQLite**: Moderno, bien mantenido, estándar de industria; validación y serialización robustas
2. **Dataclasses + Decimal + SQLite**: Más ligero, sin dependencias externas para el modelo; validación manual
3. **Attrs + Decimal + SQLite**: Equilibrio entre características y simplicidad; comunidad más pequeña

**Decisión:**

Se adopta **Pydantic v2 + Decimal + SQLite**:

- **Validación y modelos:** Pydantic v2 (`BaseModel` con `model_config = ConfigDict(frozen=True)` para inmutabilidad)
- **Tipo numérico:** `decimal.Decimal` para todas las cantidades, precios y comisiones — **nunca `float`** (garantiza determinismo y reproducibilidad, mitiga el riesgo de aritmética de punto flotante identificado en la revisión de arquitectura)
- **Persistencia:** SQLite para el MVP, con capa de repositorio que permita migrar a PostgreSQL sin cambiar la lógica de dominio
- **Versión de Python:** 3.11+ (coincide con la matriz de CI; se puede ampliar el rango si es necesario)
- **Framework web:** aplazado hasta la Fase 6 (API REST); candidato preferente FastAPI por su integración nativa con Pydantic

**Consecuencias:**

- ✅ Validación y serialización automáticas y robustas
- ✅ Determinismo garantizado por `Decimal`
- ✅ Migración de persistencia sin tocar el dominio (patrón repositorio)
- ✅ Continuidad natural hacia FastAPI en la fase de API
- ⚠️ Pydantic v2 añade una dependencia externa y una curva de aprendizaje
- ⚠️ `Decimal` es más lento que `float`; aceptable frente al requisito de exactitud

**Notas adicionales:**

Esta decisión desbloquea ADR-003 (traducción del modelo de dominio) y la creación de `requirements.txt` / `pyproject.toml`.

---

## ADR-003: Representación del modelo de dominio en Python

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

El `DOMAIN_MODEL.md` originalmente usaba pseudocódigo Kotlin. Ya fue traducido a pseudocódigo Python. Queda decidir la tecnología concreta con la que se materializará el modelo cuando comience la implementación (Fase 4).

**Decisión:**

El modelo de dominio se implementará con **Pydantic v2**, en coherencia con ADR-002:

- Entidades y objetos de valor como `BaseModel`
- Inmutabilidad mediante `model_config = ConfigDict(frozen=True)`
- Validación de invariantes con validadores de Pydantic (`@field_validator`, `@model_validator`)
- Cantidades, precios y comisiones tipados como `Decimal`
- Identificadores como tipos dedicados (p. ej. `TransactionId`) para seguridad de tipos
- Nomenclatura de atributos en `snake_case` (PEP 8), según ADR-001

**Consecuencias:**

- ✅ Coherencia total con el stack de ADR-002
- ✅ Las invariantes del modelo de dominio quedan enforced en tiempo de construcción
- ⚠️ El pseudocódigo Python actual de `DOMAIN_MODEL.md` es orientativo; al implementar puede requerir ajustes menores hacia la sintaxis real de Pydantic v2

**Próximos pasos:**

- Al llegar a la Fase 4, materializar los objetos de valor primero (`Quantity`, `Money`, `Timestamp`)
- Validar el modelo contra exportaciones reales de CoinTracking

---

## ADR-004: Estrategia de desarrollo (documentación primero)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

Existía una contradicción entre dos documentos del repositorio:

- `ROADMAP.md` define un enfoque **documentación primero**: completar especificaciones y base de conocimiento antes de escribir código (sin implementación hasta la Fase 4).
- `ARCHITECTURE_REVIEW.md` recomienda lo contrario: **implementar primero**, validar con datos reales de CoinTracking y refinar las especificaciones de forma iterativa, para evitar la divergencia especificación-realidad.

Ambas estrategias no pueden gobernar simultáneamente.

**Opciones consideradas:**

1. **Documentación primero**: especificaciones completas antes de implementar. Más predecible; riesgo de "agotamiento de especificación" y de specs no validadas contra datos reales.
2. **Implementación temprana (iterativa)**: código del núcleo cuanto antes, validado con datos reales. Menor riesgo de divergencia; entrega valor antes.

**Decisión:**

Se mantiene el enfoque **documentación primero**, tal como define `ROADMAP.md`. `ARCHITECTURE_REVIEW.md` queda como una revisión asesora (una instantánea de opinión), no como estrategia vinculante.

**Mitigaciones adoptadas** (para neutralizar los riesgos que señala la revisión):

- Obtener exportaciones reales de CoinTracking **temprano**, para validar las especificaciones antes de darlas por cerradas
- Especificar cada motor **justo antes** de su implementación, no todos por adelantado, para evitar el agotamiento de especificación
- Mantener las specs como documentos vivos: se refinan si la implementación revela supuestos incorrectos

**Consecuencias:**

- ✅ El repositorio deja de contradecirse: hay una única estrategia vinculante
- ✅ Se preserva la fortaleza del proyecto (disciplina de documentación)
- ⚠️ Riesgo de divergencia especificación-realidad: mitigado con datos reales tempranos
- ⚠️ La entrega de software funcional llega más tarde que en el enfoque iterativo

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

- ADR-001: Idioma del repositorio (contenido en español, identificadores en inglés) ✅ Decidido
- ADR-002: Stack de tecnología Python (Pydantic v2 + Decimal + SQLite) ✅ Decidido
- ADR-003: Representación del modelo de dominio en Python (Pydantic v2) ✅ Decidido
- ADR-004: Estrategia de desarrollo (documentación primero) ✅ Decidido

---

## Proceso de ADR

Toda decisión arquitectónica importante debe:

1. Ser propuesta en una rama nueva con un ADR borrador
2. Ser discutida en revisión de código
3. Ser revisada por arquitecto del proyecto
4. Ser aprobada por el equipo
5. Ser completada con la decisión final

Las decisiones menores pueden ser documentadas informalmente en CHANGELOG.md.
