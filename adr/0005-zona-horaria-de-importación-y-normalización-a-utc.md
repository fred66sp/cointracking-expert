# ADR-005: Zona horaria de importación y normalización a UTC

**Status:** Accepted

**Date:** 2026-07-02

## Context

La exportación CSV de CoinTracking ("Trade Table") contiene marcas temporales **sin zona horaria** (formato `DD.MM.YYYY HH:MM:SS`). CoinTracking almacena internamente en UTC y exporta en la zona configurada por el usuario en su cuenta. Sin conocer esa zona, el mismo instante puede interpretarse de formas distintas, rompiendo la reproducibilidad (riesgo señalado en `ARCHITECTURE_REVIEW.md` §7.4 y en `knowledge/cointracking/CSV_FORMAT.md` §2).

La cuenta de referencia usada para validar el formato tiene configurada la zona **"(GMT+01:00) Brussels, Copenhagen, Madrid, Paris"**, que corresponde a la zona IANA **`Europe/Madrid`** (equivalente en reglas a `Europe/Paris`/`Europe/Brussels`). Esta zona **observa horario de verano**: CET (`+01:00`) en invierno y CEST (`+02:00`) de finales de marzo a finales de octubre.

## Decision

*(Restaurada 2026-07-05 desde `DECISIONS.md` §ADR-005 — la migración automática a MADR, ADR-025, dejó esta sección vacía.)*

- La capa de importación interpreta cada marca temporal como **hora local en la zona IANA declarada** y la convierte a **UTC** para almacenamiento y cálculo interno.
- La **zona de origen es un parámetro obligatorio de importación** (no se asume silenciosamente). Para la cuenta de referencia el valor es `Europe/Madrid`.
- Todos los timestamps internos, comparaciones, ordenación de libro mayor y fronteras de año fiscal operan en **UTC**.
- Se usa una librería con base de datos de zonas horarias (`zoneinfo` de la biblioteca estándar de Python 3.9+) para gestionar el DST automáticamente; **nunca** un offset fijo.

## Consequences

- ✅ Reproducibilidad y consistencia entre plataformas (todo en UTC)
- ✅ Horario de verano gestionado correctamente (sin desfase de 1 h en verano)
- ✅ Fronteras de año fiscal y cruces on-chain (UTC) correctos
- ⚠️ La importación **debe exigir** que el usuario declare su zona; un valor incorrecto desplaza los datos
- ⚠️ **Riesgo residual a verificar:** queda por confirmar si CoinTracking exporta hora local *con* DST (lo esperado) o con offset fijo. Verificación definitiva: cruzar una transferencia con `Tx Hash` contra la marca temporal on-chain (siempre UTC). Ver `CSV_FORMAT.md` §2/§11.
- ⚠️ Casos ambiguos del cambio de hora (hora repetida/inexistente en la transición DST): política a definir si aparecen en datos reales

## Notes

**Notas adicionales:**

Esta decisión resuelve la cuestión abierta n.º 1 de `knowledge/cointracking/CSV_FORMAT.md` §11.
