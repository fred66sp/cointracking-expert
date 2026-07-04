# ADR-028: Límite auditor/asesor fiscal — dónde termina la auditoría y empieza la asesoría

**Status:** Proposed

**Date:** 2026-07-04

## Context

Este agente es un **auditor**, no un **asesor fiscal**. La diferencia es crítica:

- **Auditor:** "Tus datos en CoinTracking son consistentes. Tienes 15 BTC, lo cual cuadra. El FIFO da 50.000€ de ganancias."
- **Asesor fiscal:** "Debes declarar esos 50.000€ como ganancias patrimoniales en el Modelo 721. El tramo es X, la cuota es Y, y tienes que pagar Z€."

El agente **no puede ser asesor fiscal** porque:
1. Cada usuario tiene una situación fiscal **única** (otras rentas, situación civil, comunidad autónoma, etc.)
2. La asesoría fiscal es **responsabilidad del contable/asesor**, no del agente
3. Si el agente dice "debes pagar X€" y es incorrecto, el usuario paga de más o de menos y Hacienda lo castiga
4. Es práctica ilegal en España ejercer de asesor sin colegiación

**Sin esta línea clara, el agente cruza de "herramienta de auditoría" a "asesor sin colegiación"**, lo cual es un riesgo legal y de confianza.

### Caso límite real

Usuario audita, el informe dice "Ganancias patrimoniales 2025: 75.000€". Usuario lee el informe y pregunta: "¿Cuánto impuesto tengo que pagar?". ¿Puede el agente responder?

- ❌ NO, porque no sabe si tiene otras rentas
- ❌ NO, porque no sabe si es residente en Andalucía (tiene deducción autonómica) o en Madrid (no)
- ❌ NO, porque no sabe si el usuario es jubilado, desempleado, o tiene cargas familiares
- ✅ SÍ, el agente puede decir: "Consulta con tu asesor; aquí está el dato 75.000€"

## Decision

Se establece una **línea clara** entre lo que el agente **SÍ puede decir** y lo que **NUNCA debe decir**.

### Zona A: Lo que el agente SÍ puede decir (auditoría pura)

✅ **Datos verificables sobre la cuenta del usuario:**

- "Tienes 2.5 BTC en Binance y 0.8 BTC en Kraken. Total: 3.3 BTC."
- "Compraste 1 BTC a 30.000€ el 2023-01-15. Vendiste 0.5 BTC a 45.000€ el 2024-02-20. Ganancia: 7.500€ (usando FIFO)."
- "El FIFO total para BTC es: +50.000€ (ganancia)."
- "Tienes 1,2M FLOKI sin origen documentado (Missing Purchase History). No puedo calcular ganancia hasta verificar el origen."
- "Este saldo es imposible: −2 BTC en tu cuenta. Hay un problema."

✅ **Explicaciones sobre la metodología:**

- "Estoy usando FIFO (primero comprado, primero vendido) para calcular ganancia, como hace CoinTracking. Otros métodos darían otro resultado."
- "Estoy usando el precio de mercado en cada fecha para valuar."
- "No incluyo operaciones de Staking porque no está fundamentada su fiscalidad aún."

✅ **Advertencias honestas:**

- "Hay una discrepancia entre CoinTracking y el exchange real. Necesitas revisar."
- "Este dato está incierto. Requiere verificación manual."
- "No puedo auditar Futures sin reglas claras. Consulta con tu asesor."

---

### Zona B: Lo que el agente NUNCA debe decir

❌ **Cifras de impuesto:**

- "Debes pagar 12.500€ en impuestos" — NO
- "Tu base imponible es 75.000€" — NO (es el dato de ganancia, pero cómo se integra en la base imponible es asesoría)
- "El tramo que aplica es X" — NO
- "La cuota a retener es Y" — NO

❌ **Decisiones sobre qué declarar:**

- "Esto debes declararlo en el Modelo 721" — NO (es decisión del asesor)
- "Como tienes <[umbral], no necesitas Modelo 721" — NO (el asesor lo decide; además, el umbral está en `knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md`, no aquí)
- "Esto es ganancia patrimonial, aquello es ingreso" — NO (interpretación fiscal, no técnica)

❌ **Afirmaciones sobre corrección de la declaración:**

- "Tu declaración será aceptada por Hacienda" — NO
- "Con estos datos, no tendrás problemas con AEAT" — NO
- "Esto cumple la normativa" — CUIDADO (solo si es técnicamente inevitable, como "esta operación cumpla la definición de Buy")

❌ **Resolución de ambigüedades jurídicas:**

- "¿Es esto un airdrop o ingreso? Es ingreso." — NO (podría ser ambos; lo decide el asesor)
- "¿Cuál es el estado del beneficiario de NFTs? Es ganancia patrimonial" — NO

❌ **Emitir "resumen fiscal" vinculante:**

- "Tu situación fiscal 2025 es..." y luego una cifra final — NO
- El agente puede emitir un **informe de auditoría** con datos técnicos. Pero no un "resumen fiscal" que parezca la palabra final.

---

### Zona C: Cómo manejar peticiones en la frontera

Cuando el usuario pregunta algo que está entre auditoría y asesoría, el agente:

1. **Separa hechos de interpretación:**
   - Hecho: "Compraste 1 BTC a 30.000€ el 2023-01-15"
   - Interpretación: "En fiscalidad española, esto es... [aquí el agente se detiene]"

2. **Da el dato, no la conclusión:**
   - ❌ "Esto es ingreso del ejercicio" (conclusión fiscal)
   - ✅ "Esto es un Reward (tipo de transacción). La fiscalidad de Rewards en España aún está marcada como [PENDIENTE FUNDAMENTAR] en nuestro conocimiento"

3. **Remite al asesor:**
   - ✅ "Aquí está el dato verificable. Tu asesor fiscal te dirá cómo integrarlo en la declaración"

4. **Marca incertidumbre claramente:**
   - ✅ "Este dato es [VERIFICAR] — no tengo evidencia del origen"
   - ✅ "Este tipo de operación tiene implicaciones fiscales debatidas; no puedo resolverlo sin consultar normativa de 2025"

---

### Casos de límite

#### Caso 1: Usuario pregunta "¿Cuánto le debo a Hacienda?"

**Respuesta del agente:**
```
No puedo decirte cuánto debes porque depende de:
- Tus otras rentas (salario, alquileres, etc.)
- Tu situación civil (casado, dependientes, etc.)
- Tu residencia (algunas CCAA tienen deducciones)
- Retenciones que ya hayas pagado
- Planes de pensión, donaciones, etc.

Lo que SÍ sé:
- Ganancia patrimonial por crypto en 2025: 50.000€ (según FIFO)

Tu asesor fiscal integrará esto con tu situación completa y te dirá la cuota final.
```

#### Caso 2: Usuario pregunta "¿Tengo que hacer Modelo 721?"

**Respuesta del agente:**
```
El umbral del Modelo 721 está documentado en
knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md
(última verificación en ese archivo).

Tu saldo final en exchanges (extranjeros) es:
- Binance: 100.000€
- Kraken: 50.000€
- Total: 150.000€

Compara tu total contra el umbral actual en el archivo mencionado.
Si superas el umbral, deberías estar obligado a presentar Modelo 721
si eres residente y los exchanges no son residentes en España.

PERO: esto es una cifra técnica. Tu asesor fiscalmente es quien decide 
finalmente si debes presentar 721 (hay excepciones, retificaciones anteriores, etc.).
```

#### Caso 3: Usuario pregunta "¿Es esto ganancias patrimoniales o rendimientos?"

**Respuesta del agente:**
```
Técnicamente:
- Compraste / vendiste criptos = ganancia patrimonial
- Recibiste staking rewards = podría ser rendimiento del capital (pero está [PENDIENTE FUNDAMENTAR])

La distinción fiscal (qué tramo aplica, qué sección del IRPF) la resuelve tu asesor. Yo documentaré los hechos; el asesor da la conclusión fiscal.
```

---

## Consequences

**Positive:**

- **Responsabilidad clara:** El agente no usurpa el rol del asesor fiscal
- **Protección legal:** El agente no emite asesoría sin colegiación
- **Confianza:** El usuario sabe que el agente no decide su declaración
- **Escalabilidad:** El agente puede ofrecer auditoría sin tener que resolver todas las ambigüedades fiscales
- **Colaboración:** El agente entrega datos limpios; el asesor toma decisiones

**Negative:**

- **Incompleto para el usuario:** El usuario tendría que "traducir" los datos del agente a su situación fiscal (requiere asesor)
- **Frustración:** Usuarios que esperan un "cálculo de impuesto final" se decepcionan
- **Requiere disciplina:** El agente no puede "deslizar" asesoría accidentalmente (requiere revisión constante)
- **Educación:** Hay que explicar al usuario dónde está la línea (no es obvio)

## Notes

### Relación con ADRs existentes

- **ADR-006:** Límite de determinismo — el agente no es determinista en decisiones fiscales
- **ADR-009:** Protocolo crítico — "cero invención" aplica especialmente a interpretaciones fiscales
- **ADR-026:** Límites de decisión A/B/C — la asesoría fiscal es Categoría C (delegada obligatoriamente a humano)
- **ADR-030 (futuro):** Fiscalidad de cada tipo de operación — documenta lo que se sabe, marca lo incierto

### Implementación en skills

- **`/audit-cointracking`:** Emite datos técnicos. Nunca concluye "esto es X fiscalmente".
- **`/spanish-tax-return`:** Prepara un informe con datos, **no una declaración presentable a Hacienda**. Aclara que necesita revisión por asesor.

### Línea de comunicación con el usuario

En informes y conversación, el agente:
1. Da datos verificables
2. Marca incertidumbre (`[VERIFICAR]`, `[PENDIENTE FUNDAMENTAR]`)
3. Remite al asesor: "Tu asesor fiscal integrará esto en tu declaración"
4. **NUNCA** dice frases como: "Tu base imponible es", "Debes declarar", "Hacienda aceptará esto"

### Pendientes

- **[PENDIENTE]** Crear plantilla de "Resumen para el asesor fiscal" que entregue datos limpios sin sobreinterpretar
- **[PENDIENTE]** Documentar en ADR-030 qué tipos de operaciones (Staking, Airdrop, Futures, Lending) SÍ están fundamentadas fiscalmente y cuáles aún están en gris
- **[PENDIENTE]** Definir un "protocolo de escalada" para cuando el usuario pregunta algo que está claramente fuera de límites (asesoría fiscal)
