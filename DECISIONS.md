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
7. **Consentimiento informado antes de actuar (con consecuencias).** Ante una acción **consecuente** (irreversible, con impacto fiscal/económico, o que modifica datos), antes de proceder o de recomendar que el usuario la haga:
   1. Explica la acción y por qué es necesaria.
   2. **Advierte de la consecuencia de NO hacerla** (qué riesgo o error se mantiene), de forma **veraz y proporcionada** — sin exagerar ni inventar consecuencias.
   3. Pregunta y espera la decisión del usuario.
   **Alcance:** solo acciones consecuentes. En acciones triviales o de solo lectura **no** se aplica (evitar la fatiga de confirmación, que lleva a aprobar sin leer).

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
7. **Verificación de remediaciones por lote, no una a una.** Al guiar al usuario para corregir varios hallazgos en la web de CoinTracking, **no llames al MCP después de cada corrección individual** para comprobarla. Guía el lote completo de correcciones aplicables primero (confirmando por chat, sin consultar la API entre medias); solo cuando el usuario indique que ha terminado la ronda, invalida la caché **una vez** y verifica **todos** los hallazgos corregidos con una consulta agregada. Si el usuario prefiere ir uno a uno, respétalo pero explica el coste extra de cuota (límite 60/hora).

**Consecuencias:**

- ✅ Menos tokens y menos llamadas a la API; más rápido y barato
- ✅ Compatible con el rigor: el cálculo determinista sobre datos volcados es más fiable y trazable (ADR-006, ADR-009)
- ⚠️ La caché contiene datos reales → **gitignored**; tratarla como sensible
- ⚠️ Requiere gestionar la frescura de la caché (marca de tiempo, invalidación al cambiar datos)

---

## ADR-011: Persistencia y trazabilidad del flujo (nada sin dejar rastro)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

Se detectó un fallo grave: tras auditar Coinbase y **aplicar cambios reales** (borrado de 2 filas que afectan a datos fiscales), **no quedó ningún registro persistente** — ni informe, ni bitácora de cambios, ni se guardó la ruta de la fuente de verdad que aportó el usuario. Si se cierra el chat, se pierde todo el contexto. En un agente crítico (ADR-009) esto es inaceptable: los cambios que van a Hacienda deben ser trazables y el trabajo debe sobrevivir entre sesiones.

**Decisión — todo lo importante del flujo se persiste:**

1. **Informe por auditoría.** Toda auditoría o preparación fiscal genera un **informe persistente** en `reports/output/` (nombre con fecha), con hallazgos, acciones, verificación y pendientes. No se deja solo en el chat.
2. **Registro de cambios (append-only).** Todo cambio aplicado en CoinTracking se anota en `reports/output/REGISTRO-CAMBIOS.md`: qué se cambió, por qué, **evidencia**, estado **antes → después** y **verificación** en vivo. Nunca se borran entradas.
3. **Contexto durable en memoria.** Las rutas de fuentes de datos del usuario, el **estado de la auditoría** (cuentas hechas/pendientes) y las decisiones tomadas por chat se guardan en la **memoria** del proyecto (sobrevive entre sesiones). `CLAUDE.md` indica dónde vive todo.
4. **Ningún cambio consecuente sin rastro.** Si por chat se pide algo que altera datos o decisiones: código/decisiones → git (commit/ADR); datos de la cuenta → `reports/output/`; contexto → memoria. Al **retomar** en una sesión nueva, leer primero la memoria y `reports/output/`.

**Consecuencias:**

- ✅ Trazabilidad completa de los cambios que acaban en Hacienda (refuerza ADR-009)
- ✅ Continuidad entre sesiones: un chat nuevo recupera el estado
- ✅ El asesor puede reconstruir qué se tocó y por qué
- ⚠️ Informes y registro contienen datos reales → viven en `reports/output/` (gitignored); la memoria es privada (`~/.claude`)

---

## ADR-012: División de responsabilidades (Claude Code gestiona, Copilot explota)

**Estado:** Decidido

**Fecha:** 2026-07-02

**Contexto:**

El proyecto del agente lo **construye y mantiene Claude Code**; la **explotación** diaria (auditar cuentas, guiar correcciones, generar informes) la hará el usuario con **GitHub Copilot (Sonnet)**. Que dos herramientas modifiquen el "cerebro" del agente sin gobernanza rompería la trazabilidad y fiabilidad (ADR-009/011). Hay que fijar una frontera clara.

**Decisión — frontera de modificación:**

- **Claude Code = gestor del agente.** Único autorizado a modificar el "agente": `tools/`, `knowledge/`, `CLAUDE.md`, `DECISIONS.md` (ADRs), `.claude/`, `.github/`, `.vscode/`, `templates/`, `tests/`, `.mcp.json`. Todo cambio pasa por gobernanza (ADR si es relevante + commit).
- **Copilot = explotador.** **Lee** todo y **usa** el agente siguiendo los playbooks. **NO modifica** el agente. Sus únicas escrituras permitidas:
  1. **Outputs** en `reports/output/` (informes y `REGISTRO-CAMBIOS.md`).
  2. **Append** a `AGENT_CHANGE_REQUESTS.md` (bandeja de peticiones): si detecta un bug en el tool, un hueco de conocimiento, una regla a cambiar, etc., **lo anota ahí** en vez de editarlo; Claude Code lo procesa.
- Copilot **guía** al usuario a cambiar datos en CoinTracking (eso es acción del usuario, no del agente), y lo registra según ADR-011.

**Consecuencias:**

- ✅ Un único responsable del agente → cambios gobernados y trazables
- ✅ Copilot aporta mejoras sin romper nada: las canaliza como peticiones
- ✅ Separación limpia: "cerebro/reglas" (Claude) vs "operación/outputs" (Copilot)
- ⚠️ Requiere que Copilot respete la frontera (reforzado en `.github/copilot-instructions.md`); no es un candado técnico, es una norma

---

## ADR-013: Estructura multi-proyecto obligatoria (datos de usuario y estado; MCP pospuesto)

**Estado:** ✅ Decidido para la **fase 1** (aislamiento de `USER_INPUT/`, `reports/output/` y estado); el aislamiento del **MCP por proyecto queda pospuesto** (ver "Cuestión abierta" abajo).

**Fecha:** 2026-07-02 (propuesta inicial) — **revisado y redecidido 2026-07-03** tras corrección del usuario sobre el alcance real.

**Contexto:**

La propuesta original (v1, 2026-07-02) planteaba esto como una mejora "nice to have" a implementar cuando hubiera un segundo caso. El usuario corrigió el enfoque el 2026-07-03: **no es opcional ni futuro** — todo trabajo del agente sobre CoinTracking (auditar, declarar, lo que sea) debe ocurrir **siempre** dentro de un **proyecto activo**, porque eso es lo que aísla qué CSV y qué datos se usan. Sin esto, el agente puede mezclar sin querer datos de casos distintos.

**Decisión (fase 1 — sin tocar el MCP):**

1. **Estructura por proyecto**, dentro del repo (gitignored, salvo los `README.md`):
   - `USER_INPUT/<nombre_proyecto>/` — CSV y otras fuentes del caso (sustituye el uso plano anterior de `USER_INPUT/`).
   - `reports/output/<nombre_proyecto>/` — informes y `REGISTRO-CAMBIOS.md` del caso (sustituye el uso plano anterior de `reports/output/`).
   - Estado del proyecto: por ahora se sigue usando la memoria global (`audit_state` en `~/.claude`), pero **prefijada por proyecto**; migrar a un `estado.md` por proyecto queda como mejora futura si hace falta que Copilot lo lea directamente sin memoria (ADR-011/012).
2. **Puerta de entrada obligatoria (lo nuevo de la corrección):** en cualquier conversación, en cuanto el usuario pida algo relacionado con CoinTracking y **todavía no haya un proyecto activo fijado en esa conversación**, el agente debe, antes de ejecutar nada más:
   1. Listar los proyectos existentes (subcarpetas de `USER_INPUT/`).
   2. Si hay uno o más, **preguntar** con cuál trabajar, o si se quiere crear uno nuevo.
   3. Si no hay ninguno, ofrecer crear el primero (pidiendo un nombre).
   - Una vez fijado el proyecto activo en la conversación, se reutiliza para el resto de la sesión (no se vuelve a preguntar salvo que el usuario pida cambiar de proyecto).
3. **Migración del caso existente:** los datos que vivían en `USER_INPUT/` y `reports/output/` planos se migraron a `USER_INPUT/agp/` y `reports/output/agp/` el 2026-07-03 (nombre elegido por coincidir con el `--project agp` ya usado en `.mcp.json`).

**Cuestión abierta — MCP no aislado por proyecto todavía:**

El servidor MCP (`cointracking-mcp/`) solo admite el proyecto como **flag de arranque del proceso** (`--project`, ver `cointracking-mcp/SPEC/06-configuration.md`); ninguna tool acepta hoy un parámetro `project_name` en tiempo de ejecución. Cambiar de proyecto en el MCP implicaría reiniciar el servidor (y por tanto Claude Code). El usuario decidió explícitamente **posponer esto**: por ahora el proyecto de datos de usuario (CSV) y el `--project` del MCP **no están enlazados** — el MCP sigue usando el valor fijo de `.mcp.json` (`agp`) con independencia del proyecto de datos activo en la conversación. Si en el futuro se trabaja con un segundo proyecto real, resolver esto (opción más probable: añadir `project_name` como parámetro a las tools existentes del servidor Go) antes de confiar en el MCP para ese segundo proyecto.

**✅ Resuelta 2026-07-03 — ver ADR-016.** El MCP ya expone `cointracking_switch_project` para cambiar de proyecto activo en caliente, sin reiniciar el servidor.

**Consecuencias:**

- ✅ Aislamiento real entre casos en la capa de datos de usuario (CSV) e informes/estado
- ✅ Flujo predecible: nunca se opera "a ciegas" sin saber en qué proyecto se está
- ✅ El MCP ya está aislado por proyecto en caliente desde ADR-016 (antes de ADR-016, sus datos en vivo seguían siendo los de `--project agp` fijo con independencia del proyecto de datos activo)
- ⚠️ Requirió migrar rutas existentes (`USER_INPUT/`, `reports/output/`) y actualizar ambas skills y `CLAUDE.md` para aplicar la puerta de entrada

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

## ADR-014: Validación de duplicados con trade_id y consentimiento explícito

**Estado:** Decidido

**Fecha:** 2026-07-03

**Contexto:**

El 2026-07-03, al auditar CoinTracking, se detectaron como "duplicados exactos" operaciones FLOKI del 17.03.2024 que eran legítimas. El algoritmo comparaba solo campos del CSV (`(tipo, buy_amount, buy_currency, sell_amount, sell_currency, fee, exchange, fecha)`); como múltiples operaciones reales de Binance ocurrieron en el mismo segundo (batching), todas parecían idénticas. Basándose en ello, el usuario eliminó 29 copias que en realidad eran transacciones distintas (con trade_ids distintos en Binance API). El resultado: ~1,6 millones de FLOKI se perdieron del saldo hasta restaurar de backup.

**Raíz:** ct_audit.py no disponía de trade_id para distinguir operaciones; el CSV de CoinTracking tampoco lo incluye en todas las filas. La lógica de "duplicado = campos 100% idénticos" falló en presencia de transacciones legítimas separadas pero aparentemente idénticas.

**Decisión:**

1. **Usar trade_id como identificador único cuando esté disponible (fortaleza):**
   - Si el CSV incluye trade_id, dos operaciones con trade_ids distintos **nunca son duplicados**, aunque todos los demás campos sean idénticos.
   - Si trade_id está vacío, caer a la heurística siguiente.

2. **Heurística cautelosa para duplicados sin trade_id:**
   - Si hay **exactamente 2 copias idénticas** → probable duplicado de reimportación; marcar para revisión.
   - Si hay **3-10 copias** → **ADVERTENCIA** (posibles operaciones legítimas en el mismo segundo); reportar pero **no recomendar eliminar**.
   - Si hay **más de 10 copias idénticas** → **muy probablemente legítimas** (batching de Binance, transacciones FIAT repetidas, recompensas recurrentes); reportar como "INFORMACIÓN" e **indicar que requiere confirmación manual en Binance API** antes de eliminar.

3. **Implementar consentimiento informado (refuerza ADR-009):**
   - Antes de eliminar duplicados, el agente **lista exactamente cuáles se eliminarán** con ejemplos concretos (monto, tipo, fecha, cantidad de copias).
   - **Advierte:** "Si estas operaciones son legítimas (según Binance API), eliminarlas causará saldo negativo del activo."
   - **Pide confirmación explícita** del usuario antes de proceder.

4. **Usar el MCP como árbitro (cuando esté disponible):**
   - Si el MCP de CoinTracking está conectado, consulta el número de operaciones de ese tipo con trade_ids distintos.
   - Si son más de 1, son legítimas; no eliminar.
   - Si es solo 1, es un duplicado real; OK eliminar.

**Cambios en ct_audit.py:**

- Parsear trade_id cuando esté disponible; usarlo en la clave de duplicados.
- Reportar duplicados con **gravedad escalada** según la cantidad de copias.
- Devolver resultado con campo `confidence` ("DEFINITE_DUPLICATE", "PROBABLE", "LIKELY_LEGITIMATE").

**Cambios en el agente (skill audit-cointracking):**

- Mostrar al usuario una **tabla de duplicados detectados** con cantidad, confianza y ejemplo.
- Si confianza < ALTA, preguntar: "¿Confirmas que quieres eliminarlos? (verificar primero en Binance API con `Tx ID`)"
- No proceder sin confirmación explícita si hay dudas.

**Cambios en CLAUDE.md:**

- Agregar sección ⚠️ sobre falsos positivos en duplicados.
- Instruir: "Antes de eliminar duplicados, verifica en Binance que tengan el MISMO `Tx ID`. Si tienen IDs distintos, son legítimas."

**Consecuencias:**

- ✅ Evita falsos positivos como el del 2026-07-03
- ✅ Refuerza ADR-009 (consentimiento informado antes de actuar)
- ✅ El usuario toma la decisión final, no el agente
- ⚠️ Requiere que el usuario verifique en Binance API si no confía
- ⚠️ ct_audit.py debe ser más conservador; menos automatización

---

## ADR-015: Integración de la base de casos ChatGPT como v2 curada (patrones de reconciliación)

**Estado:** Decidido

**Fecha:** 2026-07-03

**Contexto:**

Copilot (explotación) propuso, vía `AGENT_CHANGE_REQUESTS.md` (petición 2026-07-02), integrar `cointracking_casos_extended.yaml` (20 casos generados con un prompt curado a un agente ChatGPT auxiliar, ver handoff `reports/output/2026-07-02_handoff_integracion_casos_chatgpt.md`) como ampliación de `cointracking_casos_base.yaml` (10 casos, esquema mínimo, ya en el repo). El candidato aportaba más cobertura y anti-patrones, pero con heterogeneidad de estilo (listas inline vs bloque), campos vacíos como `""` en vez de `null`, y profundidad desigual en evidencia/diagnóstico.

**Decisión:**

Se ejecuta la migración por fases definida en el handoff:

- **Fase A (esquema):** se fija un esquema canónico de 16 campos (ver `knowledge/patterns/INDEX.md` §Esquema). Todos los `""` pasan a `null`; todas las listas se homogeneizan en formato bloque.
- **Fase B (curación):** los casos más resumidos del candidato (antiguos CT-004/05/06/09/11-20) se amplían con evidencia mínima accionable y pasos de diagnóstico concretos, y se enlazan con conocimiento ya existente del repo (`COST_BASIS_AND_VALIDATION.md`, `CSV_FORMAT.md`, `WEB_APP_GUIDE.md`) en vez de inventar detalle nuevo sin respaldo. Los casos de duplicados (CT-003, CT-008, CT-016, CT-019) se alinean explícitamente con **ADR-014** (validación por `trade_id` y consentimiento antes de eliminar). El caso de airdrops (antiguo CT-010) mantiene `nivel_confianza: pendiente_verificar` porque el tratamiento fiscal exacto no está cerrado en `knowledge/taxation/spain/PENDIENTES.md`.
- **Fase C (versionado):** se crea `knowledge/patterns/cointracking_casos_v2.yaml` como base **vigente**. `cointracking_casos_base.yaml` (raíz del repo) pasa a **legacy/deprecado**: no se usa en auditorías nuevas, se conserva como respaldo histórico. `cointracking_casos_extended.yaml` (raíz) queda documentado como material de origen ya superado por v2.
- **Fase D (validación):** se verifica sintaxis YAML, 100% de campos del esquema presentes en los 20 casos, y cobertura de las 5 categorías críticas de regresión (transferencias huérfanas, ventas sin base de coste, duplicados, saldos negativos, rendimientos mal clasificados) — todas presentes.

**Baja definitiva (2026-07-03):** confirmado por el usuario que `cointracking_casos_base.yaml` y `cointracking_casos_extended.yaml` ya no aportan nada sobre v2; se eliminan del repositorio (quedan recuperables en el historial de git).

**Materiales auxiliares:** `LEEME.md` y `PROMPT_CHATGPT_AGENTE.md` (raíz del repo) eran documentación de apoyo para preparar el candidato con ChatGPT; su contenido queda absorbido por este ADR y por `knowledge/patterns/INDEX.md`, por lo que se eliminan tras la integración.

**Consecuencias:**

- ✅ El agente dispone de 20 casos con esquema homogéneo, evidencia mínima explícita y trazabilidad de fuente (`fuente_recomendada_para_revalidar`)
- ✅ Los casos de duplicados quedan coherentes con el incidente y la corrección de ADR-014
- ✅ El estado legacy/deprecado de la base anterior queda documentado (cierra la petición de `AGENT_CHANGE_REQUESTS.md` 2026-07-02)
- ⚠️ El contenido sigue siendo conocimiento de patrón (cualitativo); ningún caso constituye una cifra fiscal vinculante (ADR-006/009)
- ⚠️ Los casos `pendiente_verificar`/`hipotesis` requieren reverificación antes de usarse en un informe

---

## ADR-016: Cambio de proyecto activo en caliente en el MCP (`cointracking_switch_project`)

**Estado:** Decidido

**Fecha:** 2026-07-03

**Contexto:**

ADR-013 dejó abierta la cuestión de que el MCP (`cointracking-mcp/`) solo admite el proyecto como flag de arranque del proceso (`--project` en `.mcp.json`) y ninguna tool aceptaba `project_name` en tiempo de ejecución — cambiar de proyecto exigía reiniciar el servidor. El usuario propuso primero escribir un `.env` que el agente reescribiera al cambiar de proyecto; al evaluarlo, se identificaron dos problemas: (1) el binario Go no lee ningún fichero hoy, habría que implementar parseo de `.env`, y (2) el servidor MCP se arranca una vez por sesión, así que reescribir un fichero de config no evita el reinicio — solo lo hace más seguro de tocar que `.mcp.json`. Se optó por una alternativa que resuelve el problema de raíz: un tool MCP nuevo.

**Decisión:**

Añadir el tool `cointracking_switch_project(project_name)` a `cointracking-mcp/internal/tools/switch_project.go`, hermano de `cointracking_close_project` ya existente:

1. Valida `project_name` con la misma regla que `--project` al arrancar (`config.ValidateProjectName`, alfanumérico + `_` + `-`).
2. Si coincide con el proyecto ya activo, es un no-op (`already_active: true`) que no toca la caché.
3. Si no, hace flush + close de la caché del proyecto saliente (igual que `close_project`) y abre/crea la caché SQLite del proyecto entrante bajo `{cache-dir}/{project}` (misma lógica que `NewApp` al arrancar, extraída a `openProjectCache` para no duplicarla).
4. Credenciales, `--tier` y el limitador de tasa son del proceso (una cuenta de CoinTracking), no del proyecto: no cambian.

`App` (en `app.go`) pasa a guardar `cfg`/`cache`/`store` bajo un `sync.RWMutex`, con accesores (`Project()`, `CacheManager()`, `Store()`, `CacheDir()`) que los demás tools usan en vez de leer los campos directamente — así una llamada a `switch_project` no puede dejar a otra tool leyendo un puntero a medio reemplazar.

**Consecuencias:**

- ✅ Resuelve la cuestión abierta de ADR-013: el proyecto de datos activo en la conversación y el proyecto del MCP pueden mantenerse enlazados sin reiniciar Claude Code ni editar `.mcp.json`.
- ✅ Más simple que la alternativa del `.env`: no requiere que el binario Go parsee ficheros de config ni que el agente escriba en disco antes de cada cambio de proyecto — un tool call basta.
- ✅ Verificado con test de integración (`TestSwitchProject`): nombre inválido rechazado sin tocar estado, no-op al re-seleccionar el mismo proyecto, aislamiento entre proyectos, y recarga desde disco (sin llamadas nuevas a la API) al volver a un proyecto ya visitado en el proceso.
- ✅ Ambas skills (`audit-cointracking`, `spanish-tax-return`) y `CLAUDE.md` §"Proyecto activo obligatorio" actualizados para llamar a `cointracking_switch_project` en la puerta de entrada (Paso -1) en cuanto se fija el proyecto activo.

---

## ADR-017: Protocolo de diagnóstico en orden fijo para la auditoría (endurecer falsos positivos)

**Estado:** Decidido

**Fecha:** 2026-07-03

**Contexto:**

Petición de Copilot (explotación, ADR-012) en `AGENT_CHANGE_REQUESTS.md` (2026-07-03): el playbook de reconciliación (`Paso 1` de `audit-cointracking/SKILL.md`) listaba los ocho chequeos sin un orden vinculante, lo que abre la puerta a falsos positivos — p. ej. marcar "duplicados" o "huérfanas" cuando la causa real es cobertura incompleta de exchanges, o recomendar un borrado antes de haber verificado el `Trade ID`. La propuesta se apoyó en prospección del centro de ayuda oficial de CoinTracking (`READ FIRST: General account imbalances`, `Duplicate Transactions`, `Missing Transactions Report`, `Validate Transactions`, `Roll Forward / Audit Report`, avisos de "purchase pool agotado", `Binance Import Restrictions`), cuyo contenido relevante ya estaba destilado en `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` y `CSV_FORMAT.md` — no hizo falta destilar conocimiento nuevo, solo **reordenar y explicitar** el playbook existente.

**Decisión:**

Reescribir el Paso 1 de `.claude/skills/audit-cointracking/SKILL.md` con un **orden fijo de 6 fases** (cada una reduce falsos positivos de la siguiente):

1. Cobertura de fuentes/periodo y saldos (incluye saldos negativos).
2. Duplicados, con verificación de Trade ID/Tx ID obligatoria (ADR-014) **antes** de recomendar cualquier eliminación.
3. Transferencias huérfanas y orden temporal, por niveles (Tx Hash fuerte / heurístico con tolerancias — los umbrales exactos siguen abiertos, `CSV_FORMAT.md` §11.2, no se inventan).
4. Tipos, comisiones en tercera moneda, ventas sin base de coste y colisión de tickers.
5. Interpretación de avisos del "purchase pool" agotado.
6. Cierre: coherencia fiscal (FIFO) y riesgos residuales.

Se añade además una regla explícita, generalizando ADR-014: **nunca recomendar un borrado masivo sin evidencia por fila y confirmación explícita del usuario**, aplicable a cualquier hallazgo (no solo duplicados).

No se ha tocado `spanish-tax-return/SKILL.md` porque ya delega la reconciliación en `audit-cointracking` (Paso 1 de esa skill) sin duplicar el playbook. El subagente `.claude/agents/cointracking-auditor.md` (usado para análisis profundo dentro de la misma skill) se alinea con el mismo orden de 6 fases y con la regla de no borrado sin confirmación, para que no diverja del playbook cuando se delega en él.

**Consecuencias:**

- ✅ Reduce el riesgo de diagnósticos apresurados: el agente no declara "duplicado" o "huérfana" antes de confirmar cobertura completa de fuentes.
- ✅ No inventa umbrales numéricos no fundamentados (ventana temporal, tolerancia de importe) — se documentan como heurística abierta, coherente con ADR-009 (cero invención).
- ✅ Entrada `AGENT_CHANGE_REQUESTS.md` del 2026-07-03 marcada como hecha.
- ⚠️ Pendiente real (no cerrado por este ADR): definir umbrales de emparejamiento de transferencias con más datos (`CSV_FORMAT.md` §11.2) y el destilado del "purchase pool" ya existía, así que este ADR es de **reordenación operativa**, no de conocimiento nuevo.

---

## ADR-018: Discrepancia `get_gains` vs FIFO manual — documentar como hipótesis, no automatizar en `ct_audit.py` (aún)

**Estado:** Decidido

**Fecha:** 2026-07-03

**Contexto:**

Petición de Copilot (explotación, ADR-012) en `AGENT_CHANGE_REQUESTS.md`, tras investigar en el caso real `agp2025` por qué `cointracking_get_gains(price:"oldest")` (FIFO) da una ganancia de BTC (+492,87 €) muy distinta de una reconstrucción FIFO manual sobre `get_trades(trade_prices=1)` (+94,71 €; brecha ~398 €). Descartó comisiones, duplicados/mal tipado y FIFO-vs-pool como causa; la variable que explica una magnitud del mismo orden es la asimetría de qué lado de la permuta (compra o venta) se usa para valorar en EUR cada operación. Propuso automatizar un "delta bridge" en `tools/ct_audit.py` que descomponga la brecha por activo y avise si supera un umbral.

**Decisión:**

Se documenta el hallazgo, pero **no se automatiza todavía** en `tools/ct_audit.py`:

1. Nueva sección `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` §4.4, etiquetada explícitamente como **hipótesis empírica no confirmada por CoinTracking** (no hay artículo oficial que la respalde), con el caso de referencia, lo que se descartó como causa, y el recipe manual para reproducir el diagnóstico en otro caso.
2. Nuevo sub-paso en la fase 6 ("Cierre") del Paso 1 de `.claude/skills/audit-cointracking/SKILL.md`: si `get_gains` diverge de forma material de una reconstrucción FIFO manual, aplicar el recipe de COST_BASIS §4.4 antes de concluir, y marcar el resultado `[VERIFICAR]`.
3. **No se automatiza en `tools/ct_audit.py`** porque ese tool opera de forma determinista sobre el CSV export (esquema fijo y ya verificado, `CSV_FORMAT.md`), mientras que este diagnóstico depende de la respuesta JSON de `get_trades(trade_prices=1)` del MCP, cuyo esquema exacto de campos de valoración por lado (nombres de campo reales) **no está verificado** en `knowledge/cointracking/MCP_API.md`. Fijar nombres de campo sin verificarlos habría sido inventar contra ADR-009. Queda como trabajo pendiente real: verificar el esquema de `trade_prices=1` contra una llamada real y, solo entonces, decidir si se automatiza (¿en `ct_audit.py` extendido a JSON, o en un script nuevo?) — no se decide la forma aquí para no anticipar sin datos.

**Consecuencias:**

- ✅ El hallazgo del caso real queda capturado como conocimiento reutilizable, sin presentarlo como regla cerrada ni como comportamiento documentado por CoinTracking.
- ✅ La skill guía explícitamente a no dejar sin explicar una brecha material entre `get_gains` y una reconstrucción manual, ni a declararla "correcta" sin pasar por el recipe.
- ✅ No se infla el alcance de `tools/ct_audit.py` con un chequeo basado en un esquema de API sin verificar (ADR-009).
- ⚠️ Pendiente real, no cerrado por este ADR: verificar el esquema exacto de `get_trades(trade_prices=1)` y decidir la forma de automatización cuando haya otro caso real o tiempo para esa verificación.
- ✅ Entrada `AGENT_CHANGE_REQUESTS.md` del 2026-07-03 ("Chequeo automático de discrepancia FIFO...") marcada como hecha, con el alcance real aplicado (documentación + playbook, no automatización de código).

---

## ADR-019: Cierre y corrección de ADR-018 — `get_gains` confirmado fiable, la reconstrucción FIFO manual era la que fallaba

**Estado:** Decidido

**Fecha:** 2026-07-03

**Contexto:**

ADR-018 dejó la brecha BTC/USDC/OM como hipótesis `[VERIFICAR]` ("asimetría de valoración en permutas"), pendiente de contrastar contra el Tax Report oficial de CoinTracking. El usuario descargó los Tax Reports oficiales (España, FIFO) de **2024 y 2025** en Excel y se hizo el contraste real: las 39 operaciones de BTC (y todas las de OM) resultaron ser del ejercicio **2024**, no 2025 — el primer intento de mirar solo el informe de 2025 no encontraba nada porque era el año equivocado, no porque la cifra fuera cero.

**Decisión — la corrección se basa en el contraste real:**

| Activo | Tax Report oficial (2024+2025) | `get_gains(price:"oldest")` | Reconstrucción FIFO manual |
|---|---|---|---|
| BTC | 503,50 € | 492,87 € | 94,71 € |
| USDC | 554,61 € | 553,93 € | 635,61 € |
| OM | 1.027,49 € | 1.027,49 € | 1.114,89 € |

El Tax Report oficial coincide casi al céntimo con `get_gains`; la reconstrucción manual estaba mal en los tres activos. **Se corrige la conclusión de ADR-018:** la hipótesis de "asimetría de valoración por lado de permuta" queda **descartada como causa raíz** (coincidía en magnitud por casualidad). La causa más probable real: la reconstrucción manual, operación por operación, no arrastraba bien la base de coste a través de cadenas de permutas cripto-cripto; `get_gains` sí lo hace.

Se actualiza `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` §4.4 (de "hipótesis abierta" a "resuelto"), el sub-paso de la fase 6 de `audit-cointracking/SKILL.md`, `reports/output/agp2025/REGISTRO-CAMBIOS.md` y la memoria de proyecto (`audit_state`).

**Consecuencias:**

- ✅ Regla operativa nueva y más simple que la de ADR-018: ante una discrepancia `get_gains` vs. reconstrucción propia, **confiar por defecto en `get_gains`/Tax Report oficial**, no en el cálculo manual — salvo que un contraste real diga lo contrario en ese caso.
- ✅ Dato adicional relevante para la declaración: BTC (503,50 €) y OM (1.027,49 €) son ganancias del **ejercicio 2024**, no 2025 — queda pendiente (fuera de este ADR) confirmar con el usuario/asesor si ya se declararon en su ejercicio.
- ⚠️ Sigue sin verificarse el esquema exacto de `get_trades(trade_prices=1)` (ADR-018 punto pendiente) — ya no es urgente porque la recomendación operativa ya no depende de reconstruir manualmente ese cálculo, pero queda abierto si algún día se quiere automatizar un chequeo de coherencia distinto.
- ✅ Ejemplo documentado de por qué ADR-009 (cero invención) importa incluso para el propio agente: una hipótesis con evidencia aparentemente fuerte (coincidencia de magnitud) resultó ser una pista falsa; solo el contraste contra la fuente autorizada (Tax Report oficial) lo confirmó.

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
- ADR-011: Persistencia y trazabilidad del flujo (nada sin dejar rastro) ✅ Decidido
- ADR-012: División de responsabilidades (Claude Code gestiona, Copilot explota) ✅ Decidido
- ADR-013: Estructura multi-proyecto obligatoria (datos de usuario y estado; MCP pospuesto) ✅ Decidido (fase 1)
- ADR-014: Validación de duplicados con trade_id y consentimiento explícito ✅ Decidido
- ADR-015: Integración de la base de casos ChatGPT como v2 curada (patrones de reconciliación) ✅ Decidido
- ADR-016: Cambio de proyecto activo en caliente en el MCP (`cointracking_switch_project`) ✅ Decidido
- ADR-017: Protocolo de diagnóstico en orden fijo para la auditoría (endurecer falsos positivos) ✅ Decidido
- ADR-018: Discrepancia `get_gains` vs FIFO manual — documentar como hipótesis, no automatizar en `ct_audit.py` (aún) ✅ Decidido (corregido por ADR-019)
- ADR-019: Cierre y corrección de ADR-018 — `get_gains` confirmado fiable, la reconstrucción FIFO manual era la que fallaba ✅ Decidido

---

## Proceso de ADR

Toda decisión arquitectónica importante debe:

1. Ser propuesta en una rama nueva con un ADR borrador
2. Ser discutida en revisión de código
3. Ser revisada por arquitecto del proyecto
4. Ser aprobada por el equipo
5. Ser completada con la decisión final

Las decisiones menores pueden ser documentadas informalmente en CHANGELOG.md.
