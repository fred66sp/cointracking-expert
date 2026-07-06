---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-009: Protocolo de agente crítico (cero invención, máxima cautela)

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-02

## Context

Este agente trata **cifras de inversión en cripto** y produce informes (auditoría y resumen fiscal) que se envían a un **asesor fiscal** para presentar la declaración. Es un **agente crítico**: cualquier error se paga caro ante Hacienda. La corrección prevalece sobre la utilidad, la rapidez o la exhaustividad.

## Decision

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

## Consequences

- ✅ Minimiza el riesgo de error costoso ante Hacienda
- ✅ Informes fiables y auditables por el asesor
- ⚠️ El agente será más lento y preguntará/buscará más a menudo — es intencionado y deseable en este dominio
- ⚠️ Puede negarse a dar una cifra si no puede fundamentarla; correcto por diseño
