# Notas y brainstorm para la serie

## Ideas generales
[Aquí anota ideas generales que apliquen a los 3 artículos]

---

## Artículo 1: Por qué construí...

### Anécdotas/historias
[Anotar anécdotas específicas que puedas usar para abrir el artículo. Ejemplos: momento exacto de frustración, pregunta que te hizo clic, etc.]

### Datos/validación
[Notas sobre la validación del problema: cuántas empresas ofrecen este servicio, precios, etc.]

### Humor/tonalidad
[Ideas para frases ingeniosas o momentos con humor]

---

## Artículo 2: Aprendizajes

### Descubrimiento 1 - CoinTracking
[Detalles específicos sobre cómo funciona, formatos CSV, qué sorprendió]

### Descubrimiento 2 - Fiscalidad
[Reglas españolas clave, referencias a AEAT/BOE, ejemplos]

### Descubrimiento 3 - Reconciliación
[Tipos de problemas encontrados. Historias reales (sin datos sensibles)]

### Descubrimiento 4 - LLM vs algoritmo
[Argumentos técnicos: por qué explicabilidad > determinismo aquí]

---

## Artículo 3: Arquitectura

### Reto central
[Cómo plantear "dinero real = no puedo fallar" sin sonar alarmista]

### MCP propio
[Por qué no usar la API directa. Control, caché, multiproyecto]

### ADRs
[Ejemplos concretos de decisiones arquitectónicas documentadas]

### Dev/explotación
[Cómo explicar Claude Code + GitHub Copilot sin que suene raro]

### Knowledge base
[Ejemplos de cómo la KB hace auditable el agente]

### Limitaciones
[Ser explícito: qué NO es determinista, dónde está el límite de confianza]

---

## Referencias

Archivos del proyecto que pueden ser útiles:
- `DECISIONS.md` — ADRs concretos
- `knowledge/` — base de conocimiento
- `tools/ct_audit.py` — lógica determinista
- `README.md` — visión general del proyecto

---

## Checklist por artículo

### Artículo 1
- [ ] Hook convincente
- [ ] Problema claramente relatado
- [ ] Validación de mercado
- [ ] Giro/decisión
- [ ] Cierre invitador

### Artículo 2
- [ ] 3-4 descubrimientos bien desarrollados
- [ ] Accesible para no-devs
- [ ] Ejemplos concretos
- [ ] Transición suave a artículo 3

### Artículo 3
- [ ] Reto central bien planteado
- [ ] 4 soluciones claras
- [ ] Trade-offs honesto
- [ ] Cierre inspirador para devs
