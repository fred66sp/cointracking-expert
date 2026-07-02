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

## ADR-004: Estrategia de desarrollo (híbrido pragmático)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

Existía una contradicción entre dos documentos del repositorio:

- `ROADMAP.md` define un enfoque **documentación primero**: completar especificaciones y base de conocimiento antes de escribir código (sin implementación hasta la Fase 4).
- `ARCHITECTURE_REVIEW.md` recomienda lo contrario: **implementar primero**, validar con datos reales de CoinTracking y refinar las especificaciones de forma iterativa, para evitar la divergencia especificación-realidad.

El riesgo central que motiva esta decisión: **nadie conoce los datos reales de CoinTracking hasta que los mira**. Una especificación de import, duplicados o transferencias escrita sobre suposiciones puede resultar incorrecta al enfrentarse a un CSV real (comisiones que descuadran cantidades, zonas horarias distintas, movimientos en el mismo segundo). Documentar mucho sobre datos no vistos genera trabajo que luego hay que descartar.

Al mismo tiempo, hay partes del dominio que **sí** están definidas por fuentes externas estables (reglas fiscales españolas, principios de arquitectura) y se benefician de especificarse por completo antes de programar.

**Opciones consideradas:**

1. **Documentación primero (puro)**: todas las specs completas antes de cualquier código. Predecible, pero con alto riesgo de specs de datos no validadas contra la realidad.
2. **Implementación temprana (iterativa)**: código del núcleo cuanto antes, validado con datos reales. Menor riesgo de divergencia; menos énfasis documental.
3. **Híbrido pragmático**: documentación primero para lo estable; validación con datos reales antes de cerrar las specs que dependen de datos desordenados.

**Decisión:**

Se adopta el **híbrido pragmático**:

- **Documentación primero** para el dominio estable y de fuente externa:
  - Reglas de tributación (definidas por normativa; no se "descubren" programando)
  - Principios, arquitectura y contratos entre motores
  - Metodología de auditoría
- **Validación con datos reales antes de cerrar la spec** para el dominio de datos desordenados:
  - Formato CSV de CoinTracking, importación y normalización
  - Detección de duplicados, emparejamiento de transferencias, reconstrucción de libro mayor
  - Peculiaridades por exchange
  - → Estas specs se redactan en borrador, se contrastan contra **exportaciones reales de CoinTracking** y solo entonces se dan por cerradas.
- **Especificar cada motor justo antes de implementarlo**, no los nueve por adelantado, para evitar el agotamiento de especificación.
- **Las specs son documentos vivos**: se refinan si la implementación o los datos reales revelan supuestos incorrectos.

`ARCHITECTURE_REVIEW.md` queda como una revisión asesora (una instantánea de opinión), no como estrategia vinculante.

**Consecuencias:**

- ✅ El repositorio deja de contradecirse: hay una única estrategia vinculante
- ✅ Se preserva la fortaleza del proyecto (disciplina de documentación) donde aporta valor
- ✅ Se neutraliza el riesgo de divergencia especificación-realidad en las partes sensibles a datos
- ✅ Se evita el agotamiento de especificación (specs por motor, justo a tiempo)
- ⚠️ Requiere conseguir exportaciones reales de CoinTracking pronto — es una dependencia crítica, no opcional
- ⚠️ Exige disciplina para clasificar cada pieza como "estable" vs "sensible a datos"

---

## ADR-005: Zona horaria de importación y normalización a UTC

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

La exportación CSV de CoinTracking ("Trade Table") contiene marcas temporales **sin zona horaria** (formato `DD.MM.YYYY HH:MM:SS`). CoinTracking almacena internamente en UTC y exporta en la zona configurada por el usuario en su cuenta. Sin conocer esa zona, el mismo instante puede interpretarse de formas distintas, rompiendo la reproducibilidad (riesgo señalado en `ARCHITECTURE_REVIEW.md` §7.4 y en `knowledge/cointracking/CSV_FORMAT.md` §2).

La cuenta de referencia usada para validar el formato tiene configurada la zona **"(GMT+01:00) Brussels, Copenhagen, Madrid, Paris"**, que corresponde a la zona IANA **`Europe/Madrid`** (equivalente en reglas a `Europe/Paris`/`Europe/Brussels`). Esta zona **observa horario de verano**: CET (`+01:00`) en invierno y CEST (`+02:00`) de finales de marzo a finales de octubre.

**Opciones consideradas:**

1. **Offset fijo `+01:00`**: simple, pero **incorrecto en verano** (desplaza 1 h todas las operaciones de CEST).
2. **Asumir UTC**: incorrecto; las fechas son hora local del usuario.
3. **Zona IANA `Europe/Madrid` (DST-aware) → UTC**: interpreta la hora local respetando el horario de verano y normaliza a UTC.

**Decisión:**

- La capa de importación interpreta cada marca temporal como **hora local en la zona IANA declarada** y la convierte a **UTC** para almacenamiento y cálculo interno.
- La **zona de origen es un parámetro obligatorio de importación** (no se asume silenciosamente). Para la cuenta de referencia el valor es `Europe/Madrid`.
- Todos los timestamps internos, comparaciones, ordenación de libro mayor y fronteras de año fiscal operan en **UTC**.
- Se usa una librería con base de datos de zonas horarias (`zoneinfo` de la biblioteca estándar de Python 3.9+) para gestionar el DST automáticamente; **nunca** un offset fijo.

**Consecuencias:**

- ✅ Reproducibilidad y consistencia entre plataformas (todo en UTC)
- ✅ Horario de verano gestionado correctamente (sin desfase de 1 h en verano)
- ✅ Fronteras de año fiscal y cruces on-chain (UTC) correctos
- ⚠️ La importación **debe exigir** que el usuario declare su zona; un valor incorrecto desplaza los datos
- ⚠️ **Riesgo residual a verificar:** queda por confirmar si CoinTracking exporta hora local *con* DST (lo esperado) o con offset fijo. Verificación definitiva: cruzar una transferencia con `Tx Hash` contra la marca temporal on-chain (siempre UTC). Ver `CSV_FORMAT.md` §2/§11.
- ⚠️ Casos ambiguos del cambio de hora (hora repetida/inexistente en la transición DST): política a definir si aparecen en datos reales

**Notas adicionales:**

Esta decisión resuelve la cuestión abierta n.º 1 de `knowledge/cointracking/CSV_FORMAT.md` §11.

---

## ADR-006: El producto es un agente de IA auditor (Claude Code) sobre la base de conocimiento

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

El `ROADMAP.md` y el `PROJECT_CHARTER.md` describen un framework/SDK grande: motores deterministas en Python (ledger, FIFO, fiscal, tenencias, reportes), luego CLI, API, y la IA relegada a la Fase 7. Es un plan de meses. Al revisar el objetivo real, se constata que lo que aporta valor **ahora** —y que reaprovecha todo el conocimiento ya documentado— es un **agente de IA que audita los datos de CoinTracking del usuario**, no el SDK completo.

**Decisión:**

El **producto principal a corto plazo** es un **agente auditor de IA** que:
- Vive en **Claude Code** como **subagente + skill** (sin infraestructura de código propia).
- Usa como "cerebro" la base de conocimiento del repo (`knowledge/cointracking/*`, `knowledge/taxation/spain/*`).
- Accede a los datos por **dos vías**: el **MCP de la API de CoinTracking** (datos en vivo, cuando esté conectado) y el **CSV export** (Trade Table) como fuente/validación cruzada.
- Detecta y **explica** problemas de auditoría citando las reglas documentadas, con el formato evidencia → causa → impacto → recomendación.

**Límite de determinismo (reconciliación con FOUNDATION):**

FOUNDATION establece "la IA explica; los motores calculan" y exige reproducibilidad. Un agente LLM no es determinista. Por tanto:
- El agente **encuentra y explica** problemas (análisis cualitativo): transferencias huérfanas, ventas sin base de coste, duplicados, saldos imposibles, incoherencias fiscales. Esto es justo lo que FOUNDATION autoriza para la IA ("explicar, guiar, diagnosticar, resumir, asistir").
- El agente **no** produce cifras fiscales vinculantes por sí mismo. Las cantidades exactas (FIFO, base imponible) se marcan como **estimación no vinculante** o se delegan a un **cálculo determinista** (helper/función), nunca al criterio libre del LLM.

**Consecuencias:**

- ✅ Resultado utilizable de inmediato, reaprovechando todo el conocimiento y los principios ya escritos
- ✅ Coherente con FOUNDATION si se respeta el límite de determinismo
- ✅ Los "motores" del charter pasan a ser un **playbook de auditoría** (procedimientos del agente); pueden materializarse como helpers deterministas si se necesita rigor numérico
- ⚠️ **Supera al ROADMAP/charter en el corto plazo:** el SDK completo queda como visión futura/opcional, no como camino inmediato. Ver nota en `ROADMAP.md`.
- ⚠️ La calidad del agente depende de la cobertura del conocimiento; huecos conocidos (p. ej. fiscalidad de staking) limitan su precisión hasta cerrarse

**Próximos pasos:**

1. Definir el subagente auditor (`.claude/agents/`) con rol, principios y límite de determinismo.
2. Escribir el playbook de auditoría como skill invocable (`.claude/skills/`).
3. Conectar el MCP de CoinTracking; usar el CSV como alternativa.

---

## ADR-007: Limpieza del repositorio (alineación con el enfoque agente)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

Tras ADR-006 (el producto es un agente de IA en Claude Code, no un SDK de motores deterministas), el repositorio seguía conteniendo el andamiaje de la visión anterior: paquetes Python vacíos, especificaciones de motores, dependencias y CI de Python, y documentos de la visión de framework. Ese material ya no describe lo que se construye y genera ruido.

**Decisión:**

Se eliminan los artefactos que solo servían al SDK descartado:

- `src/` (paquetes Python vacíos), `requirements.txt`, `requirements-dev.txt`
- `.github/workflows/ci.yml` (CI de pytest/flake8/mypy) y `.github/ISSUE_TEMPLATE/`
- `engines/` (9 specs de motores deterministas → sustituidos por el playbook del agente en `.claude/skills/`)
- `ARCHITECTURE.md`, `ARCHITECTURE_REVIEW.md`, `DOMAIN_MODEL.md`, `ROADMAP.md`, `PROJECT_CHARTER.md`
- `CONTRIBUTING.md`, `docs/DEVELOPMENT_GUIDE.md`, `docs/INDEX.md`, `docs/PROJECT_MANIFESTO.md`
- Carpetas vacías de scaffolding: `cases/`, `examples/`, `prompts/`, `schemas/`, `scripts/`, `tests/`
- `COPILOT.md` → sustituido por `CLAUDE.md` (lo carga Claude Code)

Se conservan y adaptan: `.claude/` (agente + skill), `.mcp.json`, `knowledge/`, `DECISIONS.md`, `FOUNDATION.md`, `templates/`, `docs/GLOSSARY.md`, `LICENSE`, `CHANGELOG.md`, y `README.md` (reescrito para el agente).

**Consecuencias:**

- ✅ El repositorio refleja lo que es: un agente + su base de conocimiento
- ✅ Menos ruido; navegación y mantenimiento más simples
- ✅ Todo lo eliminado permanece en el historial de git si se necesita recuperar
- ⚠️ **ADRs anteriores (002, 003, 004, 006) referencian documentos ya eliminados** (`ROADMAP.md`, `ARCHITECTURE_REVIEW.md`, `DOMAIN_MODEL.md`, `PROJECT_CHARTER.md`). Se conservan sin reescribir: son **registro histórico** de las decisiones tal como se tomaron. Este ADR-007 es el contexto que explica esas referencias.

---

## ADR-008: Vigencia y actualización del conocimiento (fiscal y CoinTracking)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

El conocimiento del agente tiene **dos patas** y **ambas caducan**:

- **Fiscal:** la normativa cambia cada año — tramos de la base del ahorro (el tramo alto pasó del 28 % al 30 % en 2025), umbrales de obligaciones informativas, criterios de la DGT (consultas matizadas o superadas) y plazos de campaña.
- **CoinTracking:** la plataforma evoluciona — formato del CSV export, nuevos tickers y sufijos de colisión, herramientas y parámetros del MCP/API, límites de tasa, y peculiaridades por exchange.

Un conocimiento fijado en una fecha puede quedar **obsoleto** y hacer que el agente dé cifras o supuestos incorrectos.

**Decisión:**

1. **Metadatos de vigencia.** Todo documento de conocimiento sensible al tiempo declara en su cabecera: **Última verificación** (fecha) y **Vigencia** (ejercicios/versión a los que aplica).
2. **Comprobación de vigencia obligatoria.** Antes de apoyarse en un dato que puede haber cambiado (tramos, tipos, umbral del Modelo 721, criterios DGT; o formato CSV, tickers, herramientas MCP), el agente **compara** el contexto (ejercicio solicitado, fecha de hoy) con la "Última verificación"/"Vigencia" del documento.
3. **Ante posible desfase**, el agente **avisa al usuario** y **reverifica contra la fuente autorizada** antes de afirmar:
   - Fiscal → AEAT / BOE / DGT (búsqueda web).
   - CoinTracking → centro de ayuda oficial (URLs en `knowledge/cointracking/reference/CATALOG.md`) y **los datos reales del usuario** (el CSV/MCP son la verdad sobre el formato actual).
   Nunca presenta como vigente un dato sin confirmar que aplica.
4. **Checklist de revisión** en los índices de cada pata (`knowledge/taxation/spain/INDEX.md` y `knowledge/cointracking/INDEX.md`) con lo que cambia y con qué periodicidad.

**Consecuencias:**

- ✅ El agente no arrastra información caducada (ni fiscal ni de plataforma); se autoactualiza cuando hace falta
- ✅ Transparencia: distingue "verificado para 2025 / esta versión" de "asumido"
- ⚠️ Requiere disciplina de metadatos y, en su caso, una verificación extra (web o contra los datos reales)
- ⚠️ La reverificación depende de tener acceso a la fuente en la sesión

---

## ADR-009: Protocolo de agente crítico (cero invención, máxima cautela)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

Este agente trata **cifras de inversión en cripto** y produce informes (auditoría y resumen fiscal) que se envían a un **asesor fiscal** para presentar la declaración. Es un **agente crítico**: cualquier error se paga caro ante Hacienda. La corrección prevalece sobre la utilidad, la rapidez o la exhaustividad.

**Decisión — reglas de obligado cumplimiento:**

1. **Cero invención, cero improvisación.** Toda afirmación (dato fiscal, comportamiento de CoinTracking, cifra, clasificación) debe apoyarse en una de tres bases: (a) los **datos reales** del usuario, (b) la **base de conocimiento fundamentada** del repo, o (c) una **fuente oficial verificada** en la sesión. Sin respaldo, no se afirma.
2. **Ante un hueco o duda: parar y resolver, nunca rellenar.** El orden es: buscar en la base de conocimiento → si no está, **buscar en fuente oficial** (AEAT/BOE/DGT; centro de ayuda de CoinTracking) → si sigue sin resolverse, **preguntar al usuario**. Jamás completar con suposiciones para "quedar bien".
3. **Separar hechos de estimaciones.** Todo informe distingue explícitamente: **verificado** (con fuente citada) / **estimación no vinculante** / **supuesto pendiente de confirmar** `[VERIFICAR]` / **no verificable** con los datos disponibles.
4. **Peca de cauto.** Ante la duda, marca, avisa y escala. Es preferible "esto no lo sé con certeza, hay que verificar X" a una cifra que podría ser incorrecta.
5. **Trazabilidad total.** Toda cifra reportada debe poder rastrearse a su origen (operación, fuente, regla). Nada "de memoria".
6. **El informe es para un profesional.** Debe ser transparente y autoconsciente de sus límites: el asesor debe ver de dónde sale cada dato y qué queda por confirmar. El agente **no sustituye** su criterio ni el cálculo determinista.

Este protocolo consolida y prevalece sobre el resto de principios (FOUNDATION, ADR-006 determinismo, ADR-008 vigencia) y gobierna todas las skills y el subagente.

**Consecuencias:**

- ✅ Minimiza el riesgo de error costoso ante Hacienda
- ✅ Informes fiables y auditables por el asesor
- ⚠️ El agente será más lento y preguntará/buscará más a menudo — es intencionado y deseable en este dominio
- ⚠️ Puede negarse a dar una cifra si no puede fundamentarla; correcto por diseño

---

## ADR-010: Eficiencia de tokens y caché de datos de CoinTracking

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

Las respuestas de CoinTracking (MCP/API) pueden ser muy grandes (historial de miles de operaciones, balances con decenas de activos). Volcar ese JSON al contexto del LLM cada vez consume muchos tokens y, si se repite, malgasta también llamadas a la API (límite 60/h). Hay que trabajar de forma económica sin perder rigor (ADR-009).

**Decisión — protocolo de eficiencia:**

1. **Caché a disco.** Al obtener datos de CoinTracking, guárdalos en `.cache/cointracking/` (ignorado por git: son datos reales) con **marca de tiempo**. Antes de llamar, comprueba si hay un snapshot reutilizable.
2. **Reutilización.** Dentro de una misma conversación, reutiliza siempre lo ya obtenido (no recalcules ni recargues). Entre sesiones, reutiliza el snapshot si está **fresco**; si es antiguo o el usuario cambió datos, **refresca** y avísalo.
3. **Consultas mínimas y dirigidas.** Pide solo lo necesario: acota por **rango de fechas** y `limit`, y usa **agregados** (`get_grouped_balance`, `get_gains`) antes que el detalle completo. No traigas todo el historial si solo hace falta un ejercicio.
4. **Procesa lo grande con código, no en el contexto.** Para volúmenes grandes (p. ej. historial de operaciones), vuelca a un fichero y usa **scripts** (python/bash) para filtrar/agregar; sube al contexto **solo el resultado compacto** (conteos, totales, filas relevantes), nunca el JSON crudo completo. Cuando sea posible, obtén los datos con utilidades que **escriban directamente a disco** para que no pasen por el contexto.
5. **Nada de JSON crudo en salidas.** Informes y respuestas resumen y citan totales/ejemplos; no pegan volcados completos.
6. **Invalidación por cambios (CRÍTICO).** En cuanto pidas al usuario **modificar algo en CoinTracking** (editar/borrar/añadir operaciones, reimportar, corregir tipos), la caché queda **obsoleta**: márcala como inválida y **no la reutilices**. Antes de volver a dar cifras o informes, **confirma con el usuario que hizo el cambio** y **refresca** los datos (nueva consulta/volcado). Nunca mezcles hallazgos calculados con datos antiguos y datos nuevos.

**Consecuencias:**

- ✅ Menos tokens y menos llamadas a la API; más rápido y barato
- ✅ Compatible con el rigor: el cálculo determinista sobre datos volcados es más fiable y trazable (ADR-006, ADR-009)
- ⚠️ La caché contiene datos reales → **gitignored**; tratarla como sensible
- ⚠️ Requiere gestionar la frescura de la caché (marca de tiempo, invalidación al cambiar datos)

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
- ADR-004: Estrategia de desarrollo (híbrido pragmático) ✅ Decidido
- ADR-005: Zona horaria de importación y normalización a UTC ✅ Decidido
- ADR-006: Producto = agente de IA auditor (Claude Code) sobre la base de conocimiento ✅ Decidido
- ADR-007: Limpieza del repositorio (alineación con el enfoque agente) ✅ Decidido
- ADR-008: Vigencia y actualización del conocimiento (fiscal y CoinTracking) ✅ Decidido
- ADR-009: Protocolo de agente crítico (cero invención, máxima cautela) ✅ Decidido
- ADR-010: Eficiencia de tokens y caché de datos de CoinTracking ✅ Decidido

---

## Proceso de ADR

Toda decisión arquitectónica importante debe:

1. Ser propuesta en una rama nueva con un ADR borrador
2. Ser discutida en revisión de código
3. Ser revisada por arquitecto del proyecto
4. Ser aprobada por el equipo
5. Ser completada con la decisión final

Las decisiones menores pueden ser documentadas informalmente en CHANGELOG.md.
