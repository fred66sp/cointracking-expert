# Serie de 3 artículos LinkedIn: "Construí un agente de IA para auditar mis criptos"

## Visión general

Narración del journey de desarrollo de CoinTracking Expert como agente de IA, desde la frustración personal hasta la solución arquitectónica. Dirigida a devs + profesionales fintech/fiscal.

**Tono:** Profesional/didáctico + humor (accesible, no frío)  
**Perspectiva:** Dev que resolvió un problema propio en gestión cripto  
**Objetivo principal:** Compartir aprendizajes + visibilidad de la herramienta

---

## Estructura de la serie

### 1. "Por qué construí un agente de IA para auditar mis criptos"
**Tema:** Problema + frustración + validación de mercado  
**Objetivo:** Hook emocional, engagement, relatable  
**Tono:** Confesional con humor

- Abre con frustración real: múltiples exchanges, CoinTracking no cuadraba
- Problema relatable: datos contradictorios, incomprensión, declaración = suplicio
- Descubrimiento: empresas cobran €€€ solo por conciliar datos
- Giro: "Soy dev. Voy a construir esto"
- Cierre: invita a reconocer el problema

### 2. "Lo que aprendí: CoinTracking, fiscalidad y reconciliación"
**Tema:** Descubrimientos técnicos y de dominio  
**Objetivo:** Valor educativo, credibilidad en el dominio  
**Tono:** Didáctico con humor

- 3-4 descubrimientos clave (modelo de coste, duplicados, transferencias huérfanas, etc.)
- Por qué fiscalidad española complica todo (umbral Modelo 721: 50.000€)
- La reconciliación no es solo números; es detective work
- Cierre: agente de IA > script determinista para este problema

### 3. "Arquitectura: cómo construí un agente confiable para dinero real"
**Tema:** Decisiones arquitectónicas y gobernanza  
**Objetivo:** Credibilidad técnica, inspirar a otros devs  
**Tono:** Profesional/técnico pero con humanidad

- Reto: un error fiscal = dinero real; LLM no es determinista; ¿cómo confío?
- Soluciones: MCP propio, ADRs, separación dev/explotación, knowledge base auditable
- Trade-offs y por qué importan en dominios de alto riesgo
- Cierre: si construyes IA para dominios serios, estas prácticas importan

---

## Status

- [x] Artículo 1: "Por qué construí..." — completo
- [x] Artículo 2: "Lo que aprendí..." — completo
- [x] Artículo 3: "Arquitectura..." — completo
- [x] Review + ajustes
- [x] Listo para publicar (2026-01-XX)

---

## Posts de LinkedIn (introduction)

### Artículo 1
Comparto este artículo sobre por qué construí un agente de IA para auditar mis criptos. Si tienes posiciones en múltiples exchanges y la declaración de renta se te hizo un suplicio, probablemente reconozcas el problema. → [enlace]

### Artículo 2
Segundo artículo de la serie: lo que descubrí auditando criptos. El problema no es CoinTracking; es que reconciliar tus datos requiere ser detective, entender fiscalidad española y no caer en trampas que te arruinan el saldo. → [enlace]

### Artículo 3
Tercero y final: la arquitectura detrás del agente. Cuando auditas criptos que van a Hacienda, la confianza no es opcional. Aquí está cómo lo hice: MCP propio, ADRs, separación dev/explotación, y transparencia total. → [enlace]

---

## Notas

- Usar ejemplos reales del proyecto (pero sin exposiciones de datos sensibles)
- Enlazar a conocimiento base (`knowledge/`) cuando sea relevante
- Cada artículo tiene ~1200-1400 palabras (LinkedIn sweet spot)
- Serie conectada pero cada artículo funciona de forma independiente
