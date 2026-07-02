# COPILOT.md

# Marco de trabajo CoinTracking Expert

Este archivo contiene instrucciones permanentes para asistentes de IA codificación (GitHub Copilot, Claude Code, Cursor, ChatGPT y herramientas similares).

Siempre lee este archivo antes de proponer cambios.

---

# Misión del proyecto

Este repositorio busca convertirse en el marco de referencia de código abierto para auditar bases de datos de CoinTracking, reconciliación de criptomonedas y análisis fiscal de criptomonedas español.

El proyecto prioriza la exactitud, la explicabilidad y la mantenibilidad sobre la velocidad de implementación.

---

# Fuente de verdad

Siempre considera el siguiente orden de autoridad:

1. FOUNDATION.md
2. Especificaciones del proyecto
3. Documentos de arquitectura
4. Código fuente existente

Si la implementación entra en conflicto con la documentación, asume que la documentación es correcta e informa de la inconsistencia.

Nunca cambies el comportamiento empresarial silenciosamente.

---

# Principios de trabajo

Antes de implementar cualquier cosa:

- Lee la especificación relevante.
- Busca implementaciones existentes.
- Reutiliza conceptos existentes.
- Evita duplicación.

Nunca introduzcas reglas de negocio sin documentación.

---

# Flujo de trabajo de desarrollo

Para cada tarea:

1. Entiende la solicitud.
2. Identifica módulos afectados.
3. Revisa las especificaciones.
4. Sugiere cambios de arquitectura si es necesario.
5. Implementa solo la funcionalidad solicitada.
6. Añade o actualiza pruebas.
7. Actualiza documentación.
8. Resume los cambios.

---

# Filosofía del repositorio

El repositorio es impulsado por documentación.

La documentación es parte del producto.

Las reglas de negocio son independientes de la implementación.

La IA asiste a los desarrolladores.

La IA no define el comportamiento contable.

---

# Estilo de código

Prioriza:

- legibilidad
- mantenibilidad
- comportamiento determinista
- funciones pequeñas
- bajo acoplamiento
- alta cohesión

Evita optimizaciones prematuras.

---

# Documentación

Cada característica importante debe tener:

- especificación
- ejemplos
- implementación
- pruebas

Nunca dejes comportamiento sin documentar.

---

# Arquitectura

Respeta los límites de módulos.

No mezcles:

- contabilidad
- tributación
- lógica específica del exchange
- reportes
- IA

Cada módulo debe tener una única responsabilidad.

---

# Pruebas

Cada bug debe generar:

- una prueba de regresión
- documentación actualizada si es necesaria

Nunca corrijas un bug sin prevenir su recurrencia.

---

# Comportamiento de IA

Cuando los requisitos sean ambiguos:

Pregunta.

No asumas.

Cuando la documentación sea incompleta:

Propón mejoras antes de la implementación.

Cuando la arquitectura parezca inconsistente:

Detente y explica el problema.

No implementes soluciones alternativas a problemas arquitectónicos.

---

# Lista de verificación de calidad

Antes de considerar una tarea como completada, verifica:

- Documentación actualizada
- Pruebas actualizadas
- Especificaciones aún válidas
- Sin lógica duplicada
- Sin violaciones arquitectónicas
- Nomenclatura consistente
- Interfaces públicas coherentes

---

# Filosofía de commits

Prefiere muchos commits pequeños sobre commits grandes.

Cada commit debe representar un cambio lógico.

---

# Pull Requests

Cada pull request debe explicar:

- Qué cambió
- Por qué cambió
- Impacto arquitectónico
- Consideraciones futuras

---

# Objetivo a largo plazo

Cada cambio debe mover el repositorio hacia convertirse en:

- un SDK de Python reutilizable
- un motor de reconciliación profesional
- un motor contable determinista
- un marco fiscal de criptomonedas español
- una plataforma de auditoría asistida por IA

Nunca optimices por conveniencia a corto plazo a expensas de la mantenibilidad a largo plazo.

---

# Regla final

Piensa como el ingeniero jefe de un proyecto de código abierto a largo plazo.

Cada decisión debe hacer que el repositorio sea más fácil de entender, más fácil de mantener y más fácil de extender.
