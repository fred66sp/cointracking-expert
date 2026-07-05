---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-032: Knowledge with Temporal Validity — validación de vigencia de conocimiento

**Status:** Accepted

**Date:** 2026-07-05

**Accepted:** 2026-07-05

## Context

El proyecto maneja **conocimiento que envejece**. Ejemplos:

- **Fiscalidad:** Tramo IRPF 28% → 30% en 2025. Modelo 721 umbral 50.000€ (puede cambiar).
- **Regulación:** MiCA entra en vigor 2024. Salida de Binance UE en 2026-07. DAC8 cambia plazos.
- **CoinTracking:** Formato CSV, tickers, API endpoints, parámetros, límites de tasa.
- **Plazos:** Campaña fiscal (31.03, 30.06, 31.10), períodos de prescripción.

Cada dato tiene una **fecha de vigencia**. Si el agente usa un dato envejecido:
- Resultado: cifra incorrecta
- Impacto: usuario declara mal
- Riesgo: sanción de Hacienda

**Problema sin este ADR:**

```
knowledge/taxation/spain/CAPITAL_GAINS.md

"Última verificación: 2025-06-15
El tramo de la base del ahorro es 28%"

Hoy es: 2026-08-01

¿Es vigente ese 28%? ¿Cambió a 29%?
El agente **no sabe**.
```

**ADR-008** declara que hay "vigencia", pero **no operacionaliza** cómo el agente actúa cuando detecta desfase.

Este ADR lo hace.

## Decision

Todo documento de conocimiento que tenga vigencia temporal declara **metadatos de vigencia obligatorios**. El agente valida esos metadatos antes de usar el conocimiento. Si no puede validar, **no usa el dato**.

### Estructura de metadatos

Cada archivo en `knowledge/` que sea temporal declara en su cabecera (YAML frontmatter):

```yaml
---
title: Nombre del documento
vigencia:
  valid_from: YYYY-MM-DD (fecha desde la que el conocimiento es válido)
  valid_until: YYYY-MM-DD (fecha hasta la que es válido; null = indefinido)
  last_verified: YYYY-MM-DD (fecha de última verificación contra fuente)
  source: "AEAT" / "BOE" / "CoinTracking" / "DGT" / "URL" (origen de verdad)
  confidence: "high" / "medium" / "low" (cuán fiable es)
  notes: "p. ej. 'Cambio anual. Reverificar si es enero siguiente'"
---
```

**Ejemplo real:**

```yaml
---
title: "Tramos de la base del ahorro (IRPF 2026)"
vigencia:
  valid_from: "2026-01-01"
  valid_until: "2026-12-31"
  last_verified: "2026-07-05"
  source: "AEAT - Consultas DGT y normativa fiscal 2026"
  confidence: "high"
  notes: "Cambio anual. Valores específicos para 2026. Reverificar enero 2027."
---
```

### Protocolo de validación

Antes de usar cualquier dato de `knowledge/`, el agente ejecuta:

#### Paso 1: Leer metadatos de vigencia

```
Doc = Load("knowledge/taxation/spain/CAPITAL_GAINS.md")
Vigencia = Doc.metadata.vigencia
Hoy = Today()
```

#### Paso 2: Validar vigencia

```
¿Hoy >= Vigencia.valid_from?
  SI → continúa
  NO → PARAR. El conocimiento aún no es válido.

¿Vigencia.valid_until != null?
  SI → ¿Hoy <= Vigencia.valid_until?
         SI → continúa
         NO → PARAR. El conocimiento está envejecido.
  NO → continúa (indefinido)
```

#### Paso 3: Validar recency (edad de la verificación)

```
DiasDesdeVerificacion = Hoy - Vigencia.last_verified

¿Vigencia.confidence == "high"?
  SI → MaxDías = 365 (1 año)
¿Vigencia.confidence == "medium"?
  SI → MaxDías = 180 (6 meses)
¿Vigencia.confidence == "low"?
  SI → MaxDías = 30 (1 mes)

¿DiasDesdeVerificacion > MaxDías?
  SI → AVISO: "El conocimiento no ha sido verificado en [días]. Reverificar contra fuente oficial."
  NO → continúa
```

#### Paso 4: Actuar según estado

**Si vigencia es válida y reciente:**
```
✅ Usar el dato.
Registrar: "Usado dato de [fuente] verificado [fecha]"
```

**Si está en el límite de recency:**
```
⚠️ AVISO AL USUARIO:
"Este dato proviene de [fuente], última verificación [fecha].
Si es antiguo, consultar fuente oficial: [URL/source]"
```

**Si está envejecido o fuera de vigencia:**
```
❌ PARAR
"El conocimiento sobre [tema] está envejecido (verificado [fecha], válido hasta [fecha]).
Requiere reverificación contra [fuente oficial].
No continuaré sin confirmar vigencia."

[Esperar que el usuario verifique o que el agente pueda acceder a la fuente]
```

**Si `confidence == "low"`:**
```
⚠️ BAJA CONFIANZA
"Este dato tiene confianza baja. Consultar fuente oficial antes de usarlo.
Fuente: [source]"
```

### Categorización de conocimiento por criticidad

No todo conocimiento requiere vigencia estricta. Se clasifica:

#### Nivel 1: CRÍTICO (vigencia obligatoria)

Cualquier cambio impacta directamente en la declaración fiscal:
- Tramos IRPF
- Umbrales Modelo 721
- Plazos de presentación
- Tipos de ganancias patrimoniales
- Cambios regulatorios (MiCA, DAC8)

**Regla:** `valid_until` NUNCA puede ser null. Siempre tiene fecha.

#### Nivel 2: IMPORTANTE (vigencia recomendada)

Cambios afectan cálculos pero no el gesto fiscal inmediato:
- Formatos CSV de CoinTracking
- Parámetros de API
- Tickers y símbolos
- Límites de tasa (API calls/hora)

**Regla:** Tiene `valid_until`, pero puede ser "mientras esté vigente en CoinTracking" (más flexible).

#### Nivel 3: REFERENCIA (vigencia no crítica)

Información de contexto que no envejece rápido:
- Definiciones (qué es FIFO)
- Explicaciones de normas
- Historiadores (cómo era antes)

**Regla:** Puede no tener `valid_until` (indefinido), pero sí `last_verified`.

### Integración en el agente

#### En skills (`/audit-cointracking`, `/spanish-tax-return`)

Antes de usar cualquier dato de `knowledge/`:

```python
def use_knowledge(doc_path, data_key):
    doc = load_document(doc_path)
    vigencia = doc.metadata.get('vigencia')
    
    if not vigencia:
        # Sin metadatos = no usar
        return ERROR("[PENDIENTE FUNDAMENTAR] El documento no declara vigencia")
    
    if not validate_vigencia(vigencia):
        return ERROR(f"[PENDIENTE DE VERIFICAR] Vigencia envejecida: {vigencia}")
    
    # OK para usar
    return doc[data_key]
```

#### En documentación de `knowledge/`

Todo documento debe tener cabecera YAML con metadatos. Ejemplo:

```
---
title: "Ganancias patrimoniales en IRPF (España 2026)"
vigencia:
  valid_from: "2026-01-01"
  valid_until: "2026-12-31"
  last_verified: "2026-07-05"
  source: "AEAT - Normativa fiscal 2026, BOE"
  confidence: "high"
  notes: "Anual. Requiere reverificación enero 2027."
---

[contenido del documento]
```

#### En copilot-instructions.md (ADR-012)

Copilot debe respetar este protocolo:

> Nunca uses un dato de `knowledge/` sin verificar su metadato `vigencia`.
> Si está envejecido, marca como `[PENDIENTE DE VERIFICAR]`.
> Si confidence es "low", cita la fuente.

---

## Consequences

**Positive:**

- **Prevención de datos envejecidos:** El agente detecta automáticamente cuándo un dato es viejo
- **Confianza en la fuente:** Cada dato tiene trazabilidad a su origen (AEAT, BOE, CoinTracking)
- **Escalabilidad:** El mismo protocolo funciona para todo conocimiento (fiscal, técnico, regulatorio)
- **Responsabilidad:** Queda claro quién verificó qué y cuándo
- **Automatizable:** Se puede checkear en pre-commit (¿todos los docs tienen `vigencia`?)
- **Alineado con ADR-008:** Operacionaliza el concepto de "vigencia"

**Negative:**

- **Overhead de mantenimiento:** Cada documento requiere metadatos y actualización anual
- **Riesgo de abandono:** Si nadie actualiza `last_verified`, el sistema se vuelve inútil
- **Falsos positivos:** A veces el conocimiento sigue siendo válido aunque `valid_until` pasó
- **Requiere disciplina:** Alguien tiene que actualizar la cabecera YAML cada año

---

## Notes

### Relación con ADRs existentes

- **ADR-008:** Vigencia y actualización del conocimiento — este ADR operacionaliza ADR-008 con un protocolo concreto
- **ADR-009:** Protocolo crítico — esta estructura refuerza "cero invención" documentando fuentes
- **ADR-030:** Validación de ADRs — este ADR asegura que los datos que citan los ADRs tienen vigencia verificable
- **ADR-031:** Validación de plazos — referencia `knowledge/taxation/spain/FILING_DEADLINES.md` con metadatos de vigencia

### Implementación gradual

**Fase 1 (INMEDIATO):**
- Actualizar `knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md` con metadatos (ya tiene `Última verificación`, convertir a YAML)
- Crear `knowledge/taxation/spain/FILING_DEADLINES.md` con metadatos de plazos anuales

**Fase 2 (SEMANA 1):**
- Auditar todos los documentos de `knowledge/` y añadir metadatos donde falten
- Crear checklist pre-commit: "¿Todos los docs tienen `vigencia`?"

**Fase 3 (SEMANA 2):**
- Implementar en `/audit-cointracking`: validar vigencia antes de usar datos de `knowledge/`
- Implementar en `/spanish-tax-return`: idem

**Fase 4 (SEMANA 3+):**
- Automatizar en CI/CD: alerta si documento está a 30 días de `valid_until`
- Dashboard de "Documentos que requieren reverificación próximamente"

### Pendientes

- **[PENDIENTE]** Definir "quién actualiza `last_verified` cada año" (responsabilidad)
- **[PENDIENTE]** Crear plantilla YAML estándar para cabecera de `vigencia`
- **[PENDIENTE]** Script que valide `vigencia` en pre-commit
- **[PENDIENTE]** Alertas automáticas si documento está próximo a `valid_until`
- **[PENDIENTE]** Dashboard de "Documentos envejecidos" (para auditoría)

### Ejemplo práctico

**Archivo:** `knowledge/taxation/spain/CAPITAL_GAINS.md`

**Antes:**
```
# Ganancias patrimoniales en IRPF (España)

Última verificación: 2026-07-05

Los tramos son...
```

**Después:**
```yaml
---
title: "Ganancias patrimoniales en IRPF (España 2026)"
vigencia:
  valid_from: "2026-01-01"
  valid_until: "2026-12-31"
  last_verified: "2026-07-05"
  source: "AEAT, BOE, Normativa fiscal 2026"
  confidence: "high"
  notes: "Anual. Requiere actualización cada enero."
---

# Ganancias patrimoniales en IRPF (España)

Los tramos son...
```

**En el agente:**
```python
doc = load("knowledge/taxation/spain/CAPITAL_GAINS.md")
if not validate_vigencia(doc.vigencia):
    raise RequiresVerification(
        f"Documento envejecido. Último verificado: {doc.vigencia.last_verified}"
    )
use_data(doc)
```

---

## Por qué este ADR es importante

**Sin ADR-032:** El agente usa datos envejecidos y confía en que "probablemente sigan siendo válidos".

**Con ADR-032:** El agente sabe exactamente cuándo un dato deja de ser válido y **se niega a actuar** sin reverificar.

Eso es la diferencia entre "herramienta casual" y "herramienta de dominio crítico".
