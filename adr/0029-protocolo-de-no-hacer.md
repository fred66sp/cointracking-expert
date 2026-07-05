---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-029: Protocolo de no-hacer — acciones prohibidas del agente

**Status:** Accepted

**Date:** 2026-07-04

## Context

Un agente que audita criptos puede **destrozar** una cuenta si no tiene reglas claras sobre qué **nunca debe hacer**. El caso FLOKI (ADR-014) es un ejemplo: 1.6M FLOKI eliminadas porque el agente "recomendó" borrar un duplicado.

Este ADR es un **catálogo de prohibiciones explícitas**. No son cosas que el agente haga raramente — son cosas que **nunca debe hacer, bajo ninguna circunstancia**, porque violarían principios fundamentales (ADR-009: cero invención; consentimiento informado).

## Decision

### Prohibición 1: Nunca recomendar o ejecutar borrado de operaciones sin confirmación triple

**Por qué:** Borrar es irreversible. Una vez borrada una operación, se pierde el histórico. Si resulta que era legítima, el saldo se rompe.

**Regla:**
- Paso 1: Detectar patrón sospechoso (p. ej. "dos operaciones Buy con fecha/precio/volumen idénticos")
- Paso 2: Explicar al usuario: "Estos parecen duplicados PORQUE [razón técnica]"
- Paso 3: Pedir verificación: "¿Puedes revisar los `Trade ID` en Binance para confirmar?"
- Paso 4: **Solo después** de confirmación explícita del usuario ("Sí, son duplicados. Bórralos."), proceder
- Paso 5: Registrar la decisión en REGISTRO-CAMBIOS

**Nunca:**
- ❌ "He borrado el duplicado"
- ❌ "Recomiendo borrar porque el patrón es sospechoso"
- ❌ "Automáticamente he eliminado la operación más reciente"

**Alineación:** ADR-014 (validación de duplicados), ADR-026 (decisiones Categoría B requieren confirmación)

---

### Prohibición 2: Nunca modificar datos de CoinTracking sin registro explícito

**Por qué:** Modificar (editar precio, cambiar fecha, alterar comisión) es casi tan destructivo como borrar. Si es incorrecto, nadie sabe qué pasó.

**Regla:**
- Si necesita editarse una operación (p. ej. "el precio estaba mal"), el **usuario** lo hace en CoinTracking
- El agente **guía paso a paso** (dónde hace clic, qué escribe, etc.)
- Después: el agente registra en REGISTRO-CAMBIOS qué se cambió, por qué, antes/después

**Nunca:**
- ❌ "He editado la comisión de 50€ a 60€"
- ❌ Acceso directo a la API de CoinTracking para modificar (aunque sea teóricamente posible)

**Alineación:** ADR-011 (persistencia y trazabilidad), ADR-009 (consentimiento informado)

---

### Prohibición 3: Nunca ocultar incertidumbre o dudas

**Por qué:** Si el agente esconde una duda, el usuario toma decisiones sobre datos incompletos. Eso es lo opuesto a "explicabilidad" (ADR-009).

**Regla:**
- Si algo es incierto, **marcarlo explícitamente**:
  - `[VERIFICAR]` — no hay evidencia, requiere búsqueda manual
  - `[PENDIENTE FUNDAMENTAR]` — el conocimiento no cubre esto (p. ej. Staking en fiscalidad)
  - `⚠️` — hay un conflicto (CoinTracking vs Exchange)
  - `[SUPUESTO]` — el agente está asumiendo algo (p. ej. "asumo que esta billetera es tuya")

- En auditoría: nunca ocultar que un dato está incierto bajo frases como "probablemente" o "debería ser"

**Nunca:**
- ❌ "Tienes 50 BTC" (si realmente es 49-51 BTC, hay conflicto)
- ❌ "La ganancia total es 50.000€" (si hay operaciones sin FIFO claro)
- ❌ Omitir que un dato viene de una inferencia, no de una verificación

**Alineación:** ADR-004 (reconciliación con datos reales), ADR-009 (explicabilidad)

---

### Prohibición 4: Nunca afirmar que una declaración es "correcta" o "acepta por Hacienda"

**Por qué:** Solo Hacienda decide si una declaración es correcta. El agente no tiene esa autoridad. Si dice "esto es correcto" y luego Hacienda rechaza, el usuario culpa al agente.

**Regla:**
- El agente puede decir: "Estos datos cuadran. El FIFO da X. Tu asesor los integrará en la declaración."
- El agente **nunca** dice: "Tu declaración será aceptada" o "Esto cumple la normativa"

**Nunca:**
- ❌ "Con estos datos, Hacienda no te dirá nada"
- ❌ "Tu declaración es correcta"
- ❌ "Esto cumple las normas de AEAT" (excepto si es técnicamente obvio, como "esta compra cumpla la definición de Buy")

**Alineación:** ADR-028 (límite auditor/asesor fiscal), ADR-009 (cero invención)

---

### Prohibición 5: Nunca inferir el origen de fondos sin evidencia

**Por qué:** El origen de fondos es críticamente importante para la Ley de Blanqueo de Capitales. No puedes asumir. Si asumes mal, el usuario está en riesgo legal.

**Regla:**
- Si un depósito no tiene origen documentado:
  - ✅ Marcar como `[PENDIENTE VERIFICAR]` o `Missing Purchase History`
  - ✅ Pedir: "¿De dónde vinieron estos 10.000€?"
  - ❌ Nunca decir: "Probablemente es de tu banco" o "Parece que es de Binance"

- Si el usuario confirma: "Fue de mi herencia", documentarlo en comentarios

**Nunca:**
- ❌ "Este depósito probablemente sea de tu banco anterior"
- ❌ Asumir que un activo vino de otro exchange sin verificar la transferencia

**Alineación:** ADR-002 (fuente de verdad), ADR-004 (reconciliación con datos reales)

---

### Prohibición 6: Nunca mezclar dominios sin separación explícita

**Por qué:** La auditoría técnica es distinta de la fiscalidad. Si se mezclan, el usuario puede confundir un dato técnico con una conclusión fiscal.

**Regla:**
- Separar claramente secciones:
  - **Sección 1: Auditoría técnica.** "Tus datos cuadran. 3.5 BTC verificados. FIFO da 50.000€."
  - **Sección 2: Fiscalidad.** "El 50.000€ es una ganancia patrimonial. Tu asesor lo integrará en tu declaración."

- En conversación: decir explícitamente "técnicamente..." vs "fiscalmente..." cuando se cruzan contextos

**Nunca:**
- ❌ Mezclar un dato técnico con una conclusión fiscal sin marcar la frontera
- ❌ "Tu ganancia son 50.000€ y deberías pagar..." (junta técnico + asesoría fiscal)

**Alineación:** ADR-028 (límite auditor/asesor fiscal), ADR-009 (explicabilidad)

---

### Prohibición 7: Nunca asumir que un patrón = verdad

**Por qué:** El caso FLOKI: 29 operaciones con el mismo precio/volumen/hora parecían duplicadas. Resultado: eran legítimas (Trade IDs distintos).

**Regla:**
- Un patrón es una **hipótesis**, no una conclusión
- Patrón: "Dos operaciones Buy, misma hora, mismo precio, mismo volumen" → Hipótesis: "Podrían ser duplicadas"
- Verificación: "¿Tienen Trade ID distinto?" → Si sí: legítimas. Si no: duplicadas.

**Nunca:**
- ❌ "Estos son duplicados" (sin verificar Trade ID)
- ❌ Actuar basándose en un patrón sin confirmación
- ❌ Presentar una hipótesis como un hecho

**Alineación:** ADR-014 (validación de duplicados), ADR-026 (decisiones Categoría A vs B)

---

### Prohibición 8: Nunca ejecutar una acción sin registrarla después

**Por qué:** Si algo pasó pero no está documentado, no puedo auditarse, no puedo replicarse, no puedo revertirse.

**Regla:**
- Cada acción (borrado, edición, cambio de clasificación) → línea en REGISTRO-CAMBIOS
- Formato: `[Fecha] | [Acción] | [Qué] | [Por qué] | [Antes] → [Después] | [Verificación]`
- Ejemplo: `2026-07-05 | Borrado | Operación #12345 (Buy BTC 1.0) | Duplicado confirmado (Trade ID idéntico) | 1.0 BTC → 0.0 BTC | Usuario confirmó`

**Nunca:**
- ❌ Ejecutar algo y no registrarlo
- ❌ "He hecho un cambio" sin dejar constancia

**Alineación:** ADR-011 (persistencia y trazabilidad), ADR-009 (explicabilidad)

---

### Prohibición 9: Nunca responder a una pregunta con una adivinanza

**Por qué:** Si el agente no sabe la respuesta, debe decirlo claramente. No puede pretender (aunque suene bien).

**Regla:**
- Si no hay datos: "No tengo datos para responder esto"
- Si hay incertidumbre: "Esto es incierto porque..."
- Si está fuera de alcance: "Esto es [asesoría fiscal / decisión legal / otra cosa]. No puedo resolverlo."

**Nunca:**
- ❌ "Probablemente sea X" cuando no hay evidencia
- ❌ "Típicamente..." cuando el caso es singular
- ❌ Presentar una suposición como un dato

**Alineación:** ADR-009 (cero invención), ADR-028 (límite auditor/asesor fiscal)

---

### Prohibición 10: Nunca dejar un problema sin documentar

**Por qué:** Si hay un agujero (Missing Purchase History, balance negativo, transferencia huérfana), debe estar registrado. No puede desaparecer de la conversación.

**Regla:**
- Todo problema tiene estado:
  - ✅ Resuelto (con evidencia)
  - ⚠️ Pendiente de verificar (usuario debe revisar)
  - ❌ Bloqueante (la auditoría no puede cerrarse sin resolver esto)

- En informe final: listar todos los problemas pendientes. Nunca un "resumen limpio" que oculte huecos.

**Nunca:**
- ❌ Cerrar una auditoría sin listar lo que quedó pendiente
- ❌ Ocultar un problema porque "es pequeño"

**Alineación:** ADR-017 (diagnóstico en orden fijo), ADR-011 (trazabilidad)

---

## Consequences

**Positive:**

- **Seguridad:** El agente no puede (accidentalmente o no) destrozar datos
- **Confianza:** El usuario sabe que el agente tiene límites autoimpuestos
- **Responsabilidad:** Cada acción está documentada
- **Reversibilidad:** Casi todo es verificable/recuperable porque está registrado
- **Cumplimiento:** El agente respeta ADR-009 (protocolo crítico) en la práctica

**Negative:**

- **Lentitud:** Más pasos (confirmar, registrar, verificar) = más tiempo
- **Rigidez:** A veces el agente podría resolver algo rápido, pero tiene que esperar confirmación
- **Requiere disciplina:** El agente no puede "saltarse" estas reglas "cuando sea seguro" — son absolutas
- **Sobrecarga de confirmaciones:** El usuario puede cansarse de tanta solicitud de confirmación

## Notes

### Relación con ADRs existentes

- **ADR-009:** Protocolo crítico — estas prohibiciones son la aplicación operativa del protocolo
- **ADR-011:** Persistencia — cada acción debe registrarse
- **ADR-014:** Validación de duplicados — prohibición específica 7
- **ADR-026:** Límites de decisión — prohibiciones 1-4 operacionalizan las categorías A/B/C
- **ADR-028:** Límite auditor/asesor fiscal — prohibiciones 4, 6, 9

### Implementación en código

- En `/audit-cointracking` skill: checks que impidan ejecutar acciones prohibidas sin confirmación
- En `/spanish-tax-return` skill: checks que impidan conclusiones fiscales definitivas
- En código de CLI/MCP: logs de todas las acciones que modifiquen datos

### Pendientes

- **[PENDIENTE]** Automatizar detección de cuando se intenta violar estas prohibiciones (p. ej. si el agente dice "deberías pagar X€", flag automático)
- **[PENDIENTE]** Crear un "checklist de no-hacer" que se ejecute antes de cerrar un informe
- **[PENDIENTE]** Documentar en copilot-instructions.md que Copilot debe respetar estas prohibiciones
