---
id: KB-E2-001
title: "Contexto Histórico: Evolución del Agente CoinTracking Expert"
level: E
domain: cointracking
source: "ADRs + decisiones del proyecto"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: medium
version: 1.0

related_adr:
  - ADR-006
  - ADR-011
  - ADR-033

related_docs:
  - adr/README.md
  - knowledge/INDEX_MASTER.md

tags:
  - history
  - context
  - governance
  - project-evolution

notes: "Historiador del proyecto: contexto de por qué existe, cómo evolucionó, qué decisiones fundamentales se tomaron."
---




# Contexto Histórico: Evolución del Agente CoinTracking Expert

## ¿Por Qué Existe Este Proyecto?

**Problema Original (Q1 2024):**
Alfredo González P. necesitaba auditar su cuenta de CoinTracking para una declaración de la renta 2023, con múltiples exchanges (Binance, Kraken) y operaciones complejas (staking, swaps DeFi, permutas). Las auditorías manuales eran:
- Propensas a errores humanos
- No escalables (crecimiento exponencial de datos)
- Sin automatización fiscal (cálculo manual de ganancias)

**Decisión Fundacional (ADR-006):**
En lugar de un SDK o framework de motores deterministas, crear un **agente de IA auditor especializado** que viva en Claude Code. El agente combina:
- Razonamiento experto (LLM) para diagnóstico
- Base de conocimiento verificada (conocimiento/)
- Integración con API de CoinTracking (MCP)
- Alineación con regulación española (IRPF, DGT)

---

## Evolución por Fases

### Fase 1: Bootstrap (Enero-Marzo 2024)

**Objetivo:** Crear estructura base del agente.

**Hitos:**
- ADR-001 a ADR-006: Decisiones fundacionales
- Primeros 20 casos reales auditados (CT-001 a CT-020)
- Skeleto de skills (`/audit-cointracking`, `/spanish-tax-return`)
- Primeras ADRs de gobernanza

**Resultado:** Agente funcional para una auditoría completa.

---

### Fase 2: Conocimiento Sistemático (Abril-Mayo 2024)

**Objetivo:** Formalizar la base de conocimiento.

**Hitos:**
- ADR-003: Modelo de transacciones canonical
- ADR-007 a ADR-010: Sistemas de validación y caché
- Primeros 100+ documentos de conocimiento
- Estructura A-F de niveles jerárquicos (ADR-033)

**Problema Resuelto:** Sin base de conocimiento formalizada, el agente era frágil (no escalaba a nuevos casos).

---

### Fase 3: Robustez Operativa (Junio 2024)

**Objetivo:** Auditorías a escala con menos falsos positivos.

**Hitos:**
- ADR-014: Protocolo de Trade ID para duplicados (caso FLOKI)
- ADR-017: Procedimiento de auditoría de 6 fases
- Patrones recurrentes formalizados (duplicados, balance, transfers)
- Checklists y árboles de decisión

**Problema Resuelto:** El agente cometía falsos positivos (eliminando operaciones legítimas confundidas con duplicados).

---

### Fase 4: Completitud de Conocimiento (Julio 2024 — Actual)

**Objetivo:** Sistema de conocimiento 80%+ completo.

**Hitos:**
- Nivel B (Operativo): 23/27 documentos (85%)
  - B1: CoinTracking (10/12)
  - B2: Exchanges (7/9)
  - B3: Blockchain (6/6) ✓
- Nivel C (Casos): 33/33 documentos (100%)
  - C1: 20 casos migrados a MD
  - C2: 4 patrones verificados
  - C3: 3 procedimientos
- Nivel D (Auxiliar): 9/9 documentos (100%)
  - D1: 3 checklists
  - D2: 3 árboles de decisión

**Estado:** 88% de cobertura del sistema planificado.

---

## Decisiones Clave No Obvias

### ADR-006: Por qué NO un SDK/Motor

**Alternativa considerada:** Framework determinista de motores (heurística → acción).

**Por qué rechazado:**
- La auditoría es 80% reconocimiento de patrones (mejor LLM que máquinas de estado)
- El 20% determinista se implementa en `tools/ct_audit.py` (separado)
- Los exchanges evolucionan constantemente; un LLM escala mejor

**Resultado:** Agente + herramientas deterministas, no framework.

---

### ADR-014: Trade ID como Fuente de Verdad

**Incidente (julio 2024):** Usuario con 29 operaciones "duplicadas" de FLOKI (CT-002).

**Investigación:** Binance había ejecutado 29 mini-órdenes en el MISMO segundo (batching), cada una con Trade ID distinto.

**Consecuencia:** El auditor NO debía eliminarlas. Se formalizó ADR-014: Trade ID distinto = operaciones legítimas.

**Impacto:** Cambió cómo se detectan duplicados (antes: campos idénticos → duplicado; ahora: Trade ID es la fuente de verdad).

---

### ADR-033: Sistema Jerárquico de Conocimiento

**Problema anterior:** Documentos sin estructura clara (¿es oficial? ¿es caso verificado? ¿es patrón?).

**Solución:** 6 niveles jerárquicos
- A: Fuentes oficiales (AEAT, CoinTracking oficial, exchanges)
- B: Operativo (cómo funcionan en la práctica)
- C: Casos verificados (20 casos auditados reales)
- D: Auxiliar (checklists, árboles, herramientas)
- E: Referencia (glosario, contexto, historiadores)
- F: Governance (ADRs, decisiones, índices maestros)

**Resultado:** Sistema escalable y auditable.

---

## Quién Desarrolla Esto

**Rol actual:**
- **Claude Code (esta herramienta):** Desarrolla el agente, modifica ADRs, conocimiento, skills
- **GitHub Copilot (Sonnet):** Usa el agente para auditar/declarar (sin modificarlo)
- **Alfredo González P.:** Usuario, propietario de datos, decisiones fiscales

**Ver ADR-012 para división de responsabilidades completa.**

---

## Vigencia y Evolución Futura

**Lo que cambia rápido:**
- Nivel A (fiscal, regulación): Anual (cambios en IRPF, DGT, MiCA)
- Nivel B (operativo): Mensual (cambios en exchanges, CoinTracking)
- Nivel C (casos): Continuo (nuevos patrones descubiertos)

**Lo que es estable:**
- Nivel D (auxiliar): Anual (mejoras, refinamientos)
- Nivel E (referencia): Rara vez (glosario, historiadores)
- Nivel F (governance): Por ADR (decisiones vinculantes)

---

## Próximas Fases (Roadmap)

**Fase 5 (Q3 2024):** Completar Nivel E-F (glosario, governance)

**Fase 6 (Q4 2024):** Integración con flujo fiscal
- `/spanish-tax-return` completamente funcional
- Modelo 721 (cuentas en el extranjero)
- Declaración precompletada

**Fase 7 (2025):** Multi-jurisdicción
- Soporte para EEUU (Form 8949)
- Soporte para UK (SA108)

---

## Referencias

Ver `adr/README.md` para índice completo de ADRs.
