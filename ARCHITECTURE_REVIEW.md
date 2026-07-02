# Revisión de arquitectura: CoinTracking Expert

**Fecha:** 2026-07-02  
**Revisor:** Arquitecto jefe de software  
**Fase del proyecto:** Fase 1 (Fundación del proyecto)  
**Estado:** Etapa temprana (v0.1.0)  

> **Nota (2026-07-02):** Este es un documento **asesor**, no vinculante. Su recomendación de "implementar primero" fue evaluada y **no adoptada**: el proyecto mantiene la estrategia de documentación primero. Ver **ADR-004** en `DECISIONS.md`, que también recoge las mitigaciones adoptadas frente a los riesgos aquí señalados. Las demás observaciones de esta revisión siguen siendo válidas.

---

## Resumen ejecutivo

CoinTracking Expert es **un proyecto ambicioso y bien intencionado con excelente documentación y alcance peligrosamente ambicioso**. La fundación es sólida—los principios guía, modelo de dominio y visión arquitectónica son ejemplares. Sin embargo, el proyecto sufre de brechas críticas entre su especificación elaborada y su implementación nula.

**El riesgo central:** Agotamiento de especificación sin validación de implementación. El modelo de dominio de 3000 líneas es integral pero sin probar contra requisitos reales. La hoja de ruta lista 7 fases en 3+ años sin evidencia de viabilidad.

**Veredicto:** Este proyecto aún no está listo para implementación. Necesita enfoque reiniciado en los primeros 2-3 hitos antes de proceder a la visión ambiciosa a largo plazo.

---

## 1. Evaluación general

### Estado actual
- **Implementación:** 0% (solo stubs de paquetes vacíos)
- **Documentación:** 85% completada para Fase 1
- **Especificación:** Altamente detallada pero sin validar
- **Testabilidad:** No existen pruebas; CI/CD fallará en primer run
- **Riesgo de entrega:** **CRÍTICO** (brecha especificación-realidad)

### Lo que funciona bien
1. **Manifiesto y principios** son claros y bien articulados
2. **Modelo de dominio** es integral y alineado con DDD
3. **Estructura del repositorio** es lógica y escalable
4. **Visión arquitectónica** sigue patrones de arquitectura limpia
5. **Organización de base de conocimiento** es pragmática
6. **Directrices de contribución** existen y son razonables
7. **Modelo de gobernanza** con ADRs está establecido (aunque incompleto)

### Lo que no funciona
1. **Sin código funcional** - 100% implementaciones stub
2. **Sin decisiones de stack tecnológico** - Python elegido, pero ¿qué librerías?
3. **Sin estrategia de persistencia de datos** - ¿Cómo se almacenarán los datos?
4. **Sin suite de pruebas actual** - CI/CD referencias requirements no existent
5. **Sin ejemplos de integración** - Casos de estudio del mundo real están faltando
6. **Especificaciones incompletas** - Specs de motores son 1-2 páginas cuando deberían ser 10+
7. **Inconsistencia de idioma** - Modelo de dominio mezcla español y pseudocódigo Kotlin

---

## 2. Fortalezas del repositorio

### 2.1 Calidad de documentación
El proyecto excela explicando **intención**. Cada documento principal responde el "por qué":
- **PROJECT_CHARTER.md** articula visión y no-objetivos claramente
- **PROJECT_MANIFESTO.md** proporciona fundamento filosófico
- **FOUNDATION.md** define la expectativa de cultura de ingeniería
- **DOMAIN_MODEL.md** es exhaustivo para servir como especificación

Esto es raro y valioso.

### 2.2 Pensamiento arquitectónico
El enfoque de Domain-Driven Design es sólido:
- **Contextos acotados** claramente identificados (Importación, Transacción, Libro mayor, Auditoría, Impuestos, Reportes)
- **Raíces de agregado** bien identificadas (Cuenta, Auditoría)
- **Objetos de valor** modelados con restricciones
- **Invariantes** explícitamente establecidos
- **Relaciones de entidad** documentadas

La arquitectura en capas (Importación → Normalización → Motores de auditoría → Reporte) es limpia.

### 2.3 Principios guía
Los principios centrales son excelentes y diferencian este proyecto:
- **Basado en evidencia:** Cada conclusión respaldada por datos
- **Reproducibilidad:** Misma entrada → misma salida siempre
- **Explicabilidad:** Cada hallazgo incluye causa, impacto, recomendación
- **Intervención mínima:** Nunca modificar sin justificación
- **El silencio no es aceptable:** La incertidumbre debe ser reportada

Estos principios, si se mantienen, producirán software confiable.

### 2.4 Organización de base de conocimiento
La estructura del directorio `/knowledge` está bien diseñada:
```
knowledge/
  ├── cointracking/      # Conocimiento específico de plataforma
  ├── exchanges/         # APIs y comportamiento de exchanges
  ├── wallets/          # Formatos y estándares de billeteras
  ├── blockchains/      # Hechos específicos de blockchain
  ├── taxation/         # Reglas fiscales por jurisdicción
  ├── patterns/         # Patrones comunes de reconciliación
  └── faq/             # Preguntas frecuentes
```

Esta organización soporta la visión de conocimiento versionado como un activo de primera clase.

### 2.5 Estructura de gobernanza
Los ADRs (Registros de decisión de arquitectura) están planeados. Esta es una práctica buena para documentar decisiones.

---

## 3. Debilidades críticas

### 3.1 Brecha implementación-especificación (CRÍTICO)

**El problema:** El proyecto tiene especificación elaborada pero implementación nula.

```
Especificación (85% hecha)    ↔    Implementación (0% hecha)    → RIESGO: Desalineación
Modelo de dominio (3000 líneas)         Paquetes Python vacíos
Specs de motores (9 documentos)         Sin motores reales
Manifiesto (elocuente)                  Sin código que lo enforce
```

**Por qué importa:** Las especificaciones divergen de la realidad una vez que empiezas a codificar. El modelo de dominio hace suposiciones que pueden no funcionar:
- Suposiciones de inmutabilidad pueden ser impracticables con bases de datos
- Precisión de cantidad (18 decimales) puede causar problemas de rendimiento
- Heurísticas de emparejamiento de transferencias pueden ser demasiado estrictas
- Cálculo FIFO puede tener casos extremos no cubiertos

**Evidencia:** El pipeline de CI/CD references `requirements.txt` y `requirements-dev.txt` que no existen. Fallará en el primer push.

**Nivel de riesgo:** 🔴 **CRÍTICO**

### 3.2 Decisiones de stack tecnológico faltantes (CRÍTICO)

**El problema:** El modelo de dominio usa pseudocódigo Kotlin, pero el proyecto es Python. Librerías clave no están definidas:

```python
# ¿Cómo debería verse esto?
class Transaction:
    id: TransactionId
    quantity: Quantity
    price: Money
    fee: Fee
    
# Opciones sin resolver:
# - ¿dataclasses vs pydantic v2 vs attrs?
# - ¿BigDecimal vs Decimal vs float?
# - ¿namedtuples inmutables vs frozen dataclasses?
# - ¿Base de datos: SQLite, PostgreSQL, o en memoria?
# - ¿ORM: SQLAlchemy, Tortoise, Piccolo, o manual?
```

**Por qué importa:** Estas opciones afectan:
- Complejidad de serialización (pydantic vs JSON personalizado)
- Seguridad de tipos (dataclasses la pierden; pydantic la preserva)
- Rendimiento (opción de base de datos afecta patrones de consulta)
- Testabilidad (en memoria vs persistente)

**Evidencia:** No existe `setup.py`, `pyproject.toml` o `requirements.txt`. Configurar el ambiente de desarrollo es imposible.

**Nivel de riesgo:** 🔴 **CRÍTICO**

### 3.3 Especificaciones de motores incompletas (ALTO)

Comparando longitud de spec con complejidad:

| Motor | Estado | Estimado |
|--------|--------|----------|
| Motor de auditoría | 45 líneas | Debería ser 50+ |
| Motor de libro mayor | 45 líneas | Debería ser 30+ |
| Motor FIFO | 45 líneas | Debería ser 40+ |
| **Motor de duplicados** | **Solo stub (1 línea)** | **Debería ser 35+** |
| **Motor de transferencias** | **Solo stub (1 línea)** | **Debería ser 40+** |
| **Motor de tenencias** | **Solo stub (1 línea)** | **Debería ser 25+** |
| **Motor de reportes** | **Solo stub (1 línea)** | **Debería ser 35+** |
| **Motor fiscal** | **Solo stub (1 línea)** | **Debería ser 50+** |
| **Motor de reconciliación** | **Solo stub (1 línea)** | **Debería ser 30+** |

**Faltando de todas las specs:**
- Casos extremos y manejo de errores
- Análisis de complejidad de algoritmo
- Estructuras de datos e índices necesarios
- Consideraciones de rendimiento
- Puntos de integración con otros motores
- Escenarios de fallo (¿qué si la fusión falla?)
- Estrategia de concurrencia y bloqueo

**Ejemplo:** Spec FIFO no cubre:
- ¿Qué si el mismo activo se compra dos veces en el mismo segundo?
- ¿Cómo manejar ventas parciales (0.5 de 1.0 lote)?
- ¿Deberían redondearse cantidades fraccionarias?
- ¿Cómo manejar airdrops/forks que crean nuevos activos?

**Nivel de riesgo:** 🔴 **ALTO**

### 3.4 Base de conocimiento vacía (ALTO)

Los directorios `/knowledge` existen pero contienen solo stubs de `INDEX.md`:

```
knowledge/
  ├── cointracking/INDEX.md    ✓ (existe)
  │   ├── csv_format.md        ✗ (falta)
  │   ├── api_methods.md        ✗ (falta)
  │   └── ...                  ✗ (falta)
  ├── exchanges/INDEX.md       ✓ (existe)
  │   ├── binance.md           ✗ (falta)
  │   ├── coinbase.md          ✗ (falta)
  │   └── ...                  ✗ (falta)
```

**Por qué importa:** El modelo de dominio referencias conocimiento que no existe:
- Mapeos de redes de activos (USDT en Ethereum vs Polygon se ven idénticos)
- Estructuras de comisiones de exchanges y cómo se registran
- Conocimiento específico de blockchain (fees de gas, activos envueltos)
- Reglas de jurisdicción fiscal (prometido para cumplimiento español)

Estos detalles son críticos para:
1. Detección de duplicados (mismo activo en diferentes redes)
2. Emparejamiento de transferencias (contabilizar comisiones)
3. Cálculo de impuestos (reglas específicas de jurisdicción)

**Evidencia:** Las specs de exchange reclaman soporte inicial para Binance, Coinbase, Kraken, Bybit, OKX, KuCoin, BingX—pero no hay documentación de sus formatos de datos, APIs o peculiaridades.

**Nivel de riesgo:** 🟠 **ALTO**

---

## 4. Documentos faltantes

Documentos esenciales que deberían existir antes de que comience la implementación:

### 4.1 Arquitectura de implementación (CRÍTICO)
**Qué:** Decisiones de stack tecnológico y justificación
**Debería incluir:**
- Versión de Python y por qué (3.9+ mencionada pero no justificada)
- Librería de validación (pydantic, dataclasses, attrs—cuál y por qué)
- Persistencia de base de datos (en memoria para MVP, cuál para escala?)
- Serialización (JSON schema, Protocol Buffers, o personalizado?)
- Framework web (FastAPI, Flask, Django—si API planeada)
- Sistema de tipos (configuración de mypy, nivel de strictness)

**Por qué:** Sin esto, los desarrolladores comienzan a codificar y discuten sobre fundamentos después.

### 4.2 Especificación de API (ALTO - si REST API es planificada)
**Qué:** Contrato OpenAPI/AsyncAPI para API REST
**Debería incluir:**
- Endpoints `/import` (upload CSV, procesamiento en lotes)
- Endpoints `/audit` (iniciar auditoría, obtener resultados)
- Endpoints `/reports` (descargar en diferentes formatos)
- Estrategia de autenticación
- Limitación de velocidad
- Respuestas de error

**Por qué:** El diseño de API afecta toda la arquitectura. Diseña antes de implementar.

### 4.3 Schema de persistencia de datos (CRÍTICO)
**Qué:** Cómo se almacenan los datos en reposo
**Debería incluir:**
- Diagrama entidad-relación
- Schema SQL (si base de datos) o JSON schema (si basado en archivos)
- Estrategia de indexación
- Estrategia de migración (cómo manejar cambios de versión)
- Estrategia de backup y recuperación

**Por qué:** Los cambios de base de datos son costosos de retrorregular. Diseña primero.

### 4.4 Estrategia de testing (CRÍTICO)
**Qué:** Cómo se probará el proyecto
**Debería incluir:**
- Enfoque de prueba unitaria (qué probar, qué mockear)
- Estrategia de prueba integración (datos reales vs falsos)
- Conjuntos de datos de prueba (pequeño, mediano, large datasets)
- Casos de estudio del mundo real (exportaciones reales de CoinTracking)
- Baselines de rendimiento (throughput esperado)
- Tolerancia a inestabilidad (¿aceptar fallos ocasionales o cero-tolerancia?)

**Por qué:** Sin esto, la suite de pruebas se vuelve inmantenible. Define expectativas upfront.

### 4.5 Diseño de CLI (ALTO - si CLI es planificada)
**Qué:** Estructura de comando y opciones
**Debería incluir:**
- Jerarquía de comandos (`coin-tracker audit`, `coin-tracker report`, etc.)
- Archivos de configuración (.cointracker.yml, .env)
- Formatos de salida (markdown, html, json, excel)
- Reportaje de progreso (para auditorías largas)
- Manejo de errores y códigos de salida

**Por qué:** La UX de CLI molda cómo los usuarios interactúan con el proyecto.

### 4.6 Modelo de seguridad (MEDIO-ALTO)
**Qué:** Cómo el sistema maneja datos sensibles
**Debería incluir:**
- Estrategia de autenticación (si aplica)
- Modelo de autorización (quién puede auditar qué)
- Encriptación de datos (en reposo, en tránsito)
- Registro de auditoría (quién hizo qué y cuándo)
- Gestión de secretos (.env, config)
- Escaneo de dependencias (para vulnerabilidades de seguridad)

**Por qué:** Los datos cripto son sensibles. Los usuarios esperan seguridad.

### 4.7 Baseline de rendimiento (MEDIO)
**Qué:** Throughput y latencia esperados
**Debería incluir:**
- Tamaños de dataset (pequeño: 1k txs, mediano: 100k txs, grande: 1M txs)
- Latencia esperada (auditoría completada en < 10 segundos para 100k txs?)
- Uso de memoria (¿caber en 512MB? ¿2GB?)
- Uso de disco (para caching, reportes)
- Estrategia de escalabilidad (horizontal, vertical, o ambas?)

**Por qué:** Sin baselines, "lento" es subjetivo. Haz que sea objetivo.

### 4.8 Logging y observabilidad (MEDIO)
**Qué:** Cómo se diagnosticarán problemas en producción
**Debería incluir:**
- Niveles de log y qué loguear en cada
- Formato de logging estructurado (JSON con campos)
- Métricas de rendimiento (duración de auditoría, tiempos de motor)
- Reporte de errores (stack traces, contexto)
- Umbrales de alerting (si aplica)

**Por qué:** El debugging en producción es difícil sin buen logging.

### 4.9 Estrategia de deployment (MEDIO)
**Qué:** Cómo se empaquetará e lanzará el proyecto
**Debería incluir:**
- Formatos de paquete (PyPI, Docker, etc.)
- Numeración de versión (versionado semántico confirmado)
- Proceso de lanzamiento (rama main → tag → publish)
- Política de cambios rotos (cómo deprecate)
- Estrategia de rollback (si algo se rompe en producción)

**Por qué:** La gestión de lanzamiento es importante para sostenibilidad a largo plazo.

### 4.10 Guía de integración (MEDIO)
**Qué:** Cómo conectar con servicios externos
**Debería incluir:**
- Integración de API CoinTracking (si aplica)
- Ejemplos de integración de API de exchanges (Binance, Coinbase, etc.)
- Integración de RPC de blockchain (si verificación on-chain planeada)
- Integración de datos de precio (¿dónde obtener precios históricos?)

**Por qué:** Los usuarios reales necesitan conectar a fuentes de datos externas.

---

## 5. Conceptos faltantes

### 5.1 Recuperación de errores (CRÍTICO)

**Problema:** El modelo de dominio asume datos perfectos. ¿Qué pasa cuando no lo son?

Escenarios de ejemplo:
- Importación CSV es 99% completada cuando servidor se cae—¿reanudar o reiniciar?
- API de exchange devuelve formatos de timestamp inconsistentes—¿normalizar o error?
- Emparejamiento de transferencia encuentra caso ambiguo (múltiples matches posibles)—¿cuál heurística gana?
- Detección de duplicados se queda sin memoria en 1M transacciones—¿fallar o procesar en lotes?

**Faltando:** Una estrategia para:
- Fallos parciales (fallar rápido vs continuar con advertencias)
- Idempotencia (¿puedes re-ejecutar la misma auditoría dos veces de forma segura?)
- Rollback (¿qué si el usuario se da cuenta que la importación fue incorrecta?)
- Logging y recuperación (¿cómo diagnosticar y arreglar?)

### 5.2 Concurrencia (MEDIO)

**Problema:** ¿Puedes ejecutar múltiples auditorías en paralelo? ¿Deberían los motores ejecutarse en paralelo?

**Faltando:** Decisiones sobre:
- Ejecución de proceso vs thread vs async
- Thread safety de objetos de valor (Quantity, Money, etc.)
- Estrategia de bloqueo para recursos compartidos
- Condiciones de carrera en emparejamiento de transferencias (¿dos usuarios emparejando la misma transferencia?)

### 5.3 Compatibilidad de versión (MEDIO)

**Problema:** Las especificaciones cambiarán. ¿Cómo manejas datos viejos?

Ejemplo:
- v1.0 almacena cantidades como floats
- v1.1 cambia a Decimal para precisión
- ¿Puedes cargar datos de v1.0 en v1.1?

**Faltando:** Una estrategia de migración de datos.

### 5.4 Registro de auditoría (MEDIO)

**Problema:** El manifiesto exige reproducibilidad. Pero ¿qué cambió entre ejecuciones de auditoría?

**Faltando:**
- Versionado de reglas (¿cuál versión de regla "sin transacciones duplicadas" se aplicó?)
- Captura de configuración (¿qué settings se usaron?)
- Tracking de versión de motor (¿cuál versión del motor FIFO?)
- Procedencia de datos (¿de dónde vino cada transacción?)

### 5.5 Explicabilidad en código (MEDIO)

El manifiesto dice "cada hallazgo incluye causa, impacto, recomendación." Pero ¿cómo implementa el código esto?

**Faltando:** Patrones para:
- Capturar "por qué" una regla se disparó (no solo "regla 47 se disparó")
- Almacenar evidencia como datos estructurados (no solo texto formateado)
- Generar explicaciones en múltiples idiomas (actualmente se asume inglés)
- Explicaciones asistidas por IA (manifiesto lo menciona pero modelo de dominio no)

### 5.6 Soporte multi-usuario/tenant (MEDIO)

**Problema:** El modelo de dominio muestra Account/Wallet/Exchange pero ninguna entidad User.

Si esto se convierte en servicio web, necesitas:
- Autenticación y autorización de usuario
- Aislamiento de datos entre usuarios
- Modelos de compartición y colaboración
- Logging de auditoría de quién hizo qué

**Faltando:** Todo esto.

### 5.7 Internacionalización (BAJO)

El proyecto promete "Informes fiscales españoles" pero el código y manifiesto están en español.

**Faltando:**
- Cómo soportar múltiples idiomas en reportes
- Cómo localizar terminología (tipos de transacción, valores de estado)
- Manejo de monedas (€ vs $ vs ₽)

### 5.8 Diseño offline-first (BAJO)

**Problema:** ¿Qué si el usuario no tiene conexión a internet?

**Faltando:** Estrategia para:
- Usar datos de exchange cacheados
- Ejecución de auditoría offline
- Sincronización cuando la conexión regresa

---

## 6. Inconsistencias arquitectónicas

### 6.1 Desajuste de idioma del modelo de dominio

**Problema:** El modelo de dominio está en español con pseudocódigo Kotlin; el resto del proyecto es inglés Python.

```kotlin
// Desde DOMAIN_MODEL.md (español, Kotlin)
data class Transaction(
    val id: TransactionId,
    val sourceId: SourceTransactionId,
    val quantity: Quantity
)
```

vs

```python
# Lo que realmente necesitamos (español, Python)
@dataclass
class Transaction:
    id: TransactionId
    source_id: SourceTransactionId  # ¡snake_case!
    quantity: Quantity
```

**Arreglo:** Traduce el modelo de dominio a Python español y usa sintaxis actual de dataclass.

### 6.2 Patrón de repositorio vago

El modelo de dominio define interfaces de repositorio:

```kotlin
interface TransactionRepository {
    fun save(transaction: Transaction): TransactionId
    fun findById(id: TransactionId): Transaction?
    fun findByAccountId(accountId: AccountId): List<Transaction>
}
```

Pero no hay guía sobre:
- Implementación (SQLAlchemy, manual, en memoria?)
- Semántica de transacción (ACID? Consistencia eventual?)
- Optimización de consulta (índices? operaciones en lote?)
- Mecanismo de persistencia (¿base de datos? ¿archivos JSON?)

### 6.3 Contextos acotados poco claros

El modelo de dominio menciona 6 contextos acotados pero no define sus APIs:

- Contexto de importación: ¿Cómo expone transacciones normalizadas?
- Contexto de transacción: ¿Qué operaciones soporta?
- Contexto de libro mayor: ¿Es solo lectura? ¿Puede ser actualizado?
- Contexto de auditoría: ¿Cómo orquesta motores?
- Contexto fiscal: ¿Cuál es el formato de salida?
- Contexto de reportes: ¿Es separado o parte de auditoría?

**Faltando:** Capas anti-corrupción entre contextos. ¿Cómo asegura importación calidad de datos antes de pasar al contexto de transacción?

### 6.4 Suposiciones de inmutabilidad sin probar

El manifiesto enfatiza determinismo y reproducibilidad a través de inmutabilidad:

> Las transacciones de una Account son inmutables después de cierto tiempo

Pero:
- ¿Cuál es "cierto tiempo"? (¿Algunos días? Vago.)
- ¿Cómo se enforce inmutabilidad en una base de datos? (¿Restricciones de BD? ¿Lógica de aplicación?)
- ¿Qué si el usuario quiere corregir una transacción?
- ¿Qué hay de transacciones de corrección vs mutar originales?

### 6.5 Desajuste de severidad de hallazgo vs accionabilidad

El modelo de dominio define:

```kotlin
enum class Severity {
    CRITICAL,  // Requiere acción inmediata
    HIGH,      // Problema significativo
    MEDIUM,    // Notable pero no urgente
    LOW,       // Problema menor
    INFO       // Informacional
}
```

Pero ejemplos en el manifiesto muestran:
- "2 transacciones duplicadas afectando $50k" → ¿Debería ser CRITICAL? ¿O HIGH?
- "Transferencia huérfana de 0.5 BTC" → ¿Es un problema si sabes adónde fue?

**Faltando:** Enlace entre severidad e impacto comercial real.

---

## 7. Riesgos

### 7.1 Divergencia especificación-realidad (PROBABILIDAD: ALTO, IMPACTO: CRÍTICO)

**Riesgo:** El modelo de dominio de 3000 líneas no está probado contra datos reales de CoinTracking.

**Cuándo sucede:**
- Desarrollador comienza a implementar motor de transferencias
- Encuentra caso extremo: transferencias con mismo timestamp, misma cantidad, pero diferentes cuentas
- Las heurísticas de la especificación no funcionan
- Se requieren cambios de especificación
- Se pierden meses de trabajo de arquitectura

**Mitigación:**
1. Implementa PRIMERO, documenta SEGUNDO
2. Comienza con datos reales de exportación de CoinTracking (obtén de usuarios)
3. Construye iterativamente: especificación → implementación → validación → refina especificación
4. No pretendas conocer todos los casos extremos

### 7.2 Scope creep a fallo (PROBABILIDAD: MEDIO, IMPACTO: CRÍTICO)

**Riesgo:** El proyecto se vuelve ambicioso e intenta enviar todo a la vez.

**Cuándo sucede:**
- El equipo se enfoca en "visión a largo plazo" (API, MCP, agentes de IA)
- Nunca envía la librería central
- Después de 2 años, no tiene nada que mostrar a usuarios

**Mitigación:**
1. Envía MVP primero: librería + motor de auditoría simple
2. Obtén feedback de usuarios
3. Deja que los usuarios impulsen lo siguiente (CLI, API, etc.)
4. Corta el scope sin piedad para v1.0

### 7.3 Complejidad de jurisdicción fiscal (PROBABILIDAD: ALTO, IMPACTO: ALTO)

**Riesgo:** Las reglas fiscales son complejas, contradictorias y cambian frecuentemente.

**Cuándo sucede:**
- Usuario presenta declaración de impuestos usando reporte generado
- Autoridad fiscal no está de acuerdo con tus cálculos
- Usuario es responsable por penalizaciones
- El proyecto recibe culpa

**Mitigación:**
1. No reclames generar documentos fiscales; genera datos
2. Deja que contadores y profesionales fiscales interpreten
3. Incluye disclaimers prominentes
4. Construye rastros de auditoría para que los errores puedan ser identificados
5. Versiona reglas fiscales y deja a usuarios seleccionar versión del año

### 7.4 Falso sentido de determinismo (PROBABILIDAD: MEDIO, IMPACTO: MEDIO)

**Riesgo:** El proyecto reclama "misma entrada → misma salida siempre" pero hay muchos no-determinismos acechando:

- Aritmética de float (diferente orden de operaciones → diferentes resultados)
- Ordenamiento de hash (depende de versión de Python y SO)
- Redondeo de punto flotante
- Precisión de timestamp (microsegundos vs milisegundos)

**Cuándo sucede:**
- Usuario ejecuta auditoría en Windows, obtiene diferentes resultados que en Linux
- Usuario ejecuta auditoría en Python 3.9, obtiene diferentes resultados en 3.11
- Timestamp de auditoría capturado con diferente precisión
- Se ve como bug pero es realmente no-determinismo

**Mitigación:**
1. Usa Decimal, no float, para cantidades y precios
2. Prueba a través de versiones de Python
3. Captura y versiona toda configuración
4. Documenta cuál cosas es no-deterministas (ej. timing de generación de reporte)

### 7.5 Complejidad de datos del mundo real (PROBABILIDAD: ALTO, IMPACTO: ALTO)

**Riesgo:** Las exportaciones reales de CoinTracking son más desordenadas que lo que el modelo de dominio asume.

**Cuándo sucede:**
- El campo a veces está vacío, a veces lleno
- La cantidad tiene 18 decimales en una fila, 2 en otra
- Los timestamps a veces son UTC, a veces son locales
- Los tipos de transacción no están en el enum
- La estructura de comisión de exchange es más compleja que la especificación

**Mitigación:**
1. Obtén exportaciones reales de CoinTracking de usuarios temprano
2. Construye validación de importación que reporte problemas
3. Haz la normalización más tolerante (no falles en tipos desconocidos)
4. Almacena datos originales junto con datos normalizados

### 7.6 Mantenimiento de base de conocimiento (PROBABILIDAD: MEDIO, IMPACTO: MEDIO)

**Riesgo:** La base de conocimiento (specs de exchange, datos de blockchain, reglas fiscales) se vuelve obsoleta.

**Cuándo sucede:**
- Binance cambia formato CSV
- New blockchain gana significancia
- Reglas fiscales cambian
- Nadie actualiza documentación
- Los usuarios golpean bugs de conocimiento obsoleto

**Mitigación:**
1. Establece frecuencia de revisión (trimestral para impuestos, mensual para exchanges)
2. Fija versiones de specs de exchange/blockchain
3. Configura monitoreo de cambios (APIs de precio cambian, exchanges añaden tokens)
4. Directrices de contribución comunitaria para actualizaciones de conocimiento

### 7.7 Optimización prematura (PROBABILIDAD: MEDIO, IMPACTO: MEDIO)

**Riesgo:** El modelo de dominio es complejo (objetos de valor, inmutabilidad estricta, event sourcing) que puede ser overkill.

**Cuándo sucede:**
- Implementar Quantity con Decimal de 18-decimales es más lento de lo necesario
- Restricciones de inmutabilidad hacen actualizaciones incómodas
- Event sourcing añade complejidad sin beneficio
- Podría haber sido enviado más rápido con enfoque más simple

**Mitigación:**
1. Simplifica primero: usa float + Decimal, no tipos de dinero fancy
2. Haz cosas mutables hasta que necesites inmutabilidad
3. Mide antes de optimizar
4. YAGNI: You Aren't Gonna Need It (hasta que lo hagas)

### 7.8 Deuda de testing (PROBABILIDAD: ALTO, IMPACTO: ALTO)

**Riesgo:** Sin pruebas ahora → saltarse pruebas después → código inprobable.

**Cuándo sucede:**
- v1.0 se envía sin pruebas
- v1.1 "solo necesita una característica más"
- En v1.5, el código es inprobable
- Los bugs son irreparables sin romper todo
- El proyecto muere

**Mitigación:**
1. Escribe pruebas desde el día uno
2. Apunta a 80%+ cobertura en caminos críticos (motores)
3. Requiere pruebas para cada corrección de bug
4. Haz que CI enforce cobertura

---

## 8. Mejoras sugeridas

### 8.1 Acciones inmediatas (Este mes)

1. **Crea `requirements.txt` y `requirements-dev.txt`**
   - Arregla fallos del pipeline de CI/CD
   - Documenta dependencias

2. **Implementa definiciones de tipo Python**
   - Convierte pseudocódigo Kotlin a dataclasses o pydantic Python
   - Configura mypy
   - Establece type checking como no-opcional

3. **Crea fixtures de prueba**
   - Pequeñas, medianas, grandes exportaciones de muestra de CoinTracking
   - Ejemplos reales de exportación (anonimizados) de usuarios
   - Ejemplos de caso extremo (duplicados, transferencias huérfanas)

4. **Escribe primera prueba**
   - Prueba creación y validación de transacción
   - Haz que pase
   - Prueba que infraestructura de prueba funciona

5. **Completa especificaciones de motor**
   - Escribe specs completas para motores de duplicados, transferencias, tenencias, impuestos, reportes
   - Incluye casos extremos, algoritmos, manejo de errores

6. **Establece disciplina de ADR**
   - Completa ADR-001 (decisión de idioma)
   - Documenta decisiones clave en DECISIONS.md
   - Haz ADRs no-opcional antes de implementación

### 8.2 Acciones a corto plazo (Próximos 3 meses)

1. **Implementa modelo de transacción básico**
   - Crea `src/cointracking_expert/models/transaction.py`
   - Haz inmutable y validable
   - Escribe pruebas exhaustivas

2. **Implementa capa de importación**
   - Importación CSV de CoinTracking
   - Normalización a forma canónica
   - Validación básica y reporte de errores

3. **Implementa motor de libro mayor**
   - Reconstrucción de balance cronológica
   - Detección de balance negativo
   - Prueba contra datasets de 1k, 10k, 100k transacciones

4. **Implementa motor de duplicados**
   - Detección exacta de duplicados
   - Emparejamiento probabilístico (mismo timestamp, misma cantidad, mismo activo)
   - Escribe casos de prueba para casos extremos

5. **Establece patrón de repositorio**
   - Decide persistencia (SQLite ahora, PostgreSQL después)
   - Implementa repositorio en memoria para testing
   - Implementa repositorio persistente

6. **Documenta contratos de API**
   - Define interfaces para cada motor
   - Escribe spec OpenAPI (aunque no implementada aún)

### 8.3 Acciones a mediano plazo (Meses 4-6)

1. **Implementa motores restantes**
   - Motor de transferencias
   - Motor de tenencias
   - Motor FIFO
   - Motor fiscal (jurisdicción española)
   - Motor de reportes

2. **Implementa CLI**
   - Comando básico `coin-tracker audit`
   - Soporte de archivo de configuración
   - Opciones de formato de salida

3. **Testing exhaustivo**
   - Exportaciones reales de CoinTracking de 10+ usuarios
   - Suite de prueba de caso extremo
   - Testing de rendimiento (auditoría de 1M transacciones en < 30 segundos)

4. **Contenido de documentación**
   - Puebla directorios `/knowledge` con contenido real
   - Specs de exchange (Binance, Coinbase, Kraken, etc.)
   - Specs de blockchain (Bitcoin, Ethereum, etc.)
   - Reglas fiscales (al menos España, idealmente UE)

5. **Auditoría de seguridad**
   - Revisa vulnerabilidades comunes
   - Gestión de secretos (claves de API, etc.)
   - Encriptación de datos

---

## 9. Hoja de ruta propuesta para los próximos 10 hitos

### Hito 1: Configuración del proyecto (Semana 1-2)
- [ ] Crear requirements.txt y requirements-dev.txt
- [ ] Configurar ambiente de desarrollo (venv, mypy, pytest, black)
- [ ] Arreglar pipeline de CI/CD
- [ ] Implementar primera prueba que pase
- [ ] Completar DECISIONS.md con ADRs clave

**Entregable:** Ambiente de desarrollo funcional, CI/CD pasando

### Hito 2: Modelo de datos central (Semana 3-4)
- [ ] Implementar modelos de Python dataclasses/pydantic desde modelo de dominio
- [ ] Crear objetos de valor (Quantity, Money, Timestamp, etc.)
- [ ] Agregar validación exhaustiva
- [ ] Escri

bir 100+ pruebas para modelo de datos
- [ ] Documentar API del modelo de datos

**Entregable:** Estructuras de datos principales con 80%+ cobertura de prueba

### Hito 3: Importación de CSV (Semana 5-6)
- [ ] Implementar parser CSV de CoinTracking
- [ ] Implementar normalización de datos
- [ ] Crear validación de importación y reporte de errores
- [ ] Escribir pruebas para exportaciones de muestra (pequeño, mediano, grande)
- [ ] Documentar API de importación

**Entregable:** Puede cargar y normalizar exportaciones de CoinTracking

### Hito 4: Motor de libro mayor (Semana 7-9)
- [ ] Implementar procesamiento de transacción cronológico
- [ ] Calcular balances corrientes
- [ ] Detectar balances negativos
- [ ] Escribir pruebas para datasets de 100k+ transacciones
- [ ] Documentar casos extremos (transacciones del mismo segundo, etc.)

**Entregable:** Puede reconstruir libros mayores y detectar imposibilidades

### Hito 5: Motor de duplicados (Semana 10-11)
- [ ] Implementar detección exacta de duplicados
- [ ] Implementar emparejamiento probabilístico
- [ ] Afinar heurísticas de matching con datos reales
- [ ] Escribir pruebas para casos extremos
- [ ] Documentar algoritmo de matching

**Entregable:** Puede detectar duplicados con < 5% tasa de falsos positivos

### Hito 6: Motor de transferencias (Semana 12-13)
- [ ] Implementar emparejamiento de depósito/retiro
- [ ] Manejar cálculos de comisión
- [ ] Detectar transferencias huérfanas
- [ ] Escribir pruebas para transferencias entre exchanges
- [ ] Documentar heurísticas de matching

**Entregable:** Puede emparejar transferencias e identificar huérfanas

### Hito 7: Motor FIFO (Semana 14-16)
- [ ] Implementar tracking de lote de adquisición
- [ ] Calcular base de costo
- [ ] Implementar emparejamiento de disposición
- [ ] Detectar historial de compras faltante
- [ ] Escribir pruebas para ventas parciales, manejo de comisión

**Entregable:** Puede calcular lotes FIFO y base de costo

### Hito 8: Motor de tenencias (Semana 17-18)
- [ ] Calcular tenencias actuales desde libro mayor
- [ ] Tracking de lotes de adquisición por tenencia
- [ ] Calcular ganancias no realizadas
- [ ] Validar consistencia con libro mayor
- [ ] Escribir pruebas para escenarios complejos

**Entregable:** Puede calcular y validar tenencias

### Hito 9: CLI y reportes básicos (Semana 19-20)
- [ ] Implementar comando `coin-tracker audit`
- [ ] Implementar generación de reporte markdown
- [ ] Añadir soporte de archivo de configuración
- [ ] Escribir pruebas de integración
- [ ] Documentar uso de CLI

**Entregable:** Puede ejecutar auditoría completa desde línea de comandos, generar reporte markdown

### Hito 10: Documentación e impuestos españoles (Semana 21-24)
- [ ] Puebla directorios `/knowledge` con contenido real
- [ ] Implementar cálculo de eventos fiscales españoles
- [ ] Escribir ejemplos exhaustivos
- [ ] Crear casos de estudio del mundo real
- [ ] Preparar para lanzamiento v0.2.0

**Entregable:** Listo para pruebas beta, candidato de lanzamiento v0.2.0

---

## 10. Orden de implementación recomendado

### Por qué este orden?

1. **Configuración primero** (Hito 1): No puedes construir sobre fundaciones rotas
2. **Modelo de datos segundo** (Hito 2): Todo lo demás depende de esto
3. **Importación tercero** (Hito 3): Necesitas datos para probar
4. **Libro mayor y duplicados cuarto** (Hitos 4-5): Motores de auditoría central
5. **Transferencia y FIFO quinto** (Hitos 6-7): Motores avanzados
6. **Tenencias, CLI, reportes último** (Hitos 8-10): UI y conocimiento

### Qué NO hacer

❌ **No inicies con API** (Fase 6 es demasiado temprano)  
❌ **No inicies con UI** (envía librería primero)  
❌ **No inicies con integración de IA** (necesita motores centrales primero)  
❌ **No optimices prematuramente** (haz que funcione primero)  
❌ **No saltes testing** (prueba desde el día uno)  

---

## 11. Factores críticos de éxito

### 1. Disciplina de documentación-primero
La fortaleza del proyecto es documentación. **Mantén esto.** Cada motor necesita especificación antes de implementación. Cada decisión necesita un ADR.

### 2. Datos de usuario real temprano
No solo pruebas con datos sintéticos. Obtén exportaciones reales de CoinTracking de usuarios reales. Los datos desordenados reales revelarán suposiciones.

### 3. Loop de feedback ajustado
- Implementa → Prueba con datos reales → Ajusta especificación → Itera
- No esperes perfección; envía y aprende

### 4. Enfoque en MVP v1.0
Envía una librería que pueda:
- Cargar exportaciones de CoinTracking
- Detectar duplicados
- Reconstruir libros mayores
- Calcular lotes FIFO
- Generar reportes markdown

Todo lo demás (API, CLI, agentes de IA) es v2.0+.

### 5. Preservación de principios
El manifiesto es excelente. **Haz que se enforce en código:**
- Cada hallazgo debe tener evidencia
- Cada cálculo debe ser reproducible
- Cada cambio debe estar justificado
- El silencio nunca es aceptable

---

## 12. Evaluación final

### Qué mantener

✅ **Fundación de Domain-Driven Design** - Excelente  
✅ **Manifiesto y principios** - Raros y valiosos  
✅ **Estructura de repositorio** - Lógica y escalable  
✅ **Visión de arquitectura** - Capas limpias, buena separación  
✅ **Disciplina de documentación** - Establece el ejemplo

### Qué arreglar

🔧 **Decisiones de stack tecnológico** - Faltan opciones críticas  
🔧 **Especificaciones de motor** - Incompletas para 6/9 motores  
🔧 **Contenido de base de conocimiento** - Existe pero vacío  
🔧 **Infraestructura de testing** - Inexistente  
🔧 **Guía de implementación** - Demasiada especificación, muy poco cómo  

### Qué cuestionar

❓ **Alcance de visión a largo plazo** - ¿Se necesita realmente API? ¿MCP?  
❓ **Enfoque en impuestos españoles** - ¿Limita adopción internacional?  
❓ **Suposiciones de inmutabilidad** - ¿Funcionarán en práctica?  
❓ **Complejidad del modelo de dominio** - ¿Vale el costo de implementación?

### Veredicto

**Este proyecto tiene una fundación excelente pero necesita probar que puede enviar.** Las próximas 10 semanas (a través del hito 4-5) son críticas. Si el equipo puede:
1. Arreglar el pipeline de CI/CD (requirements.txt)
2. Implementar modelo de datos básico
3. Construir importación CSV funcional
4. Obtener un motor de libro mayor que maneje 100k transacciones

...entonces el proyecto ha probado que su enfoque funciona. Entonces puedes construir confiadamente los motores restantes y eventualmente enviar v1.0.

Si el equipo en su lugar persigue la visión a largo plazo mientras lucha con conceptos básicos, el proyecto se estancará.

**Recomendación:** Enfoca sin piedad en enviar v1.0 con solo funcionalidad central. Deja que los usuarios impulsen lo que viene después.

---

**Revisión completada**

**Próximos pasos:**
1. Comparte esta revisión con el equipo
2. Discute hallazgos en reunión de revisión de arquitectura
3. Crea items de acción concretos para hitos 1-3
4. Establece checkpoints de progreso semanal
5. Revisita esta revisión después del hito 5 (Motor de duplicados)
