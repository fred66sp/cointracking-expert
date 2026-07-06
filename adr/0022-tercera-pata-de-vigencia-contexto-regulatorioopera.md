---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-022: Tercera pata de vigencia — contexto regulatorio/operativo de exchanges (extiende ADR-008)

**Status:** Accepted

**Deciders:** Alfredo González P. (propietario, aprueba) · Claude Code (agente, propone)

**Date:** 2026-07-03

## Context

ADR-008 definía la vigencia del conocimiento en dos patas: normativa fiscal española y formato/plataforma de CoinTracking. Durante esta sesión apareció un caso real que no encaja en ninguna de las dos: **Binance perdió la licencia MiCA** en la UE (deadline 1 de julio de 2026) y cortó servicios a usuarios de la UE, forzando al usuario a migrar activos a Coinbase (documentado en `knowledge/exchanges/BINANCE_EU_MICA_EXIT.md`). El usuario señaló, acertadamente, que este tipo de cambio regulatorio/operativo de los exchanges **cambia de año en año** y conviene revisarlo explícitamente antes de preparar cada declaración, no solo descubrirlo por casualidad.

MiCA en sí no cambia la fiscalidad española (sigue rigiendo `CAPITAL_GAINS.md`), pero genera eventos operativos reales (migraciones forzosas, posibles conversiones) que si no se revisan pueden dejar huecos de reconciliación sin detectar.

## Decision

**Decisión:**

Añadir una **tercera pata explícita de vigencia** (ADR-008 pasa a tener 3 patas, no 2): antes de preparar una declaración, además de revisar normativa fiscal y formato CoinTracking, **hacer una búsqueda web breve sobre el contexto regulatorio/operativo de los exchanges relevantes del ejercicio** (licencias, cierres, migraciones forzosas, restricciones de producto) y avisar si aparece algo material. Añadido en `.claude/skills/spanish-tax-return/SKILL.md` (nota de vigencia al principio del documento). No se exige si el usuario confirma que no ha habido cambio de exchange en el ejercicio (evita búsquedas innecesarias, ADR-010).

## Consequences

- ✅ Institucionaliza el hallazgo de esta sesión (Binance/MiCA) como parte del protocolo recurrente, no como un hallazgo puntual que se olvida al año siguiente.
- ✅ Sigue sin inventar reglas fiscales (ADR-009): MiCA se trata como contexto operativo, no como norma tributaria nueva.
- ⚠️ Coste: una búsqueda web adicional por declaración cuando hay exchanges involucrados — aceptable dado el riesgo de perder un hecho imponible real (conversión forzosa) o de auditar una migración incompleta.
