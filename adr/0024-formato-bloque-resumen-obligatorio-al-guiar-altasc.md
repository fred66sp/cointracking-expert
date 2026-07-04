# ADR-024: Formato "bloque-resumen" obligatorio al guiar altas/correcciones manuales en CoinTracking

**Status:** Accepted

**Date:** 2026-07-04

## Context

Al guiar a un usuario para crear o corregir una operación manual en CoinTracking (fiat faltante, transferencia con un lado sin registrar, corrección de tipo, etc.), el agente daba la instrucción en prosa. El usuario probó una propuesta de formato compacto tipo "CT-Task" (generada externamente con ChatGPT) para normalizar cómo se piden estos datos, y confirmó explícitamente tras verla aplicada en dos ejemplos reales que le habría ahorrado mucho tiempo en las sesiones de corrección masiva ya realizadas sobre su cuenta. Se verificó la propuesta contra los datos reales del proyecto (`knowledge/cointracking/CSV_FORMAT.md`) antes de adoptarla: el orden de campos era correcto, pero contenía errores (`Ingreso` en vez de `Ingresos`; un tipo `Transferencia` inexistente — CoinTracking modela una transferencia como un par `Retirada`+`Depósito`, nunca como una única entrada) y tipos no verificados (`Donación`, `Minería`, `Airdrop`, `Regalo`).

**Decisión:**

1. El agente **siempre** cierra la guía de una tarea de alta/corrección manual en CoinTracking con el bloque-resumen documentado en `knowledge/cointracking/WEB_APP_GUIDE.md` §4bis — nunca en su lugar de la explicación en lenguaje llano, siempre después de ella y adaptado a la tarea concreta.
2. Estructura fija: `[ <Tipo> | <Fecha DD.MM.AAAA HH:MM:SS> ] [ <campos principales> ] [ Intercambio: … | Grupo: … | Comentario: … ]`, con el tercer bloque opcional campo a campo.
3. Usar solo los tipos/etiquetas ya confirmados contra datos reales del proyecto; tratar cualquier tipo no verificado como `[VERIFICAR]` antes de usarlo en una instrucción clic a clic (ADR-009).
4. Una transferencia entre cuentas propias nunca es un bloque único: siempre dos tareas encadenadas (`Retirada` en origen + `Depósito` en destino, con la retirada en fecha igual o anterior).
5. Los avisos que la operación pueda disparar en CoinTracking (balance negativo, "no hay compra adecuada para esta venta"...) se indican **después** del bloque, nunca dentro.
6. Esta norma aplica a **cualquiera que use el agente** (Claude Code y Copilot, ADR-012), por lo que vive en `CLAUDE.md` (estilo de guía) y no solo en la memoria de una sesión.

## Decision

[Decision not found]

## Consequences

- ✅ Comunicación consistente y predecible en toda tarea de alta/corrección manual, mejorando la experiencia del usuario novato sin sacrificar la explicación en lenguaje llano que ya exige `CLAUDE.md`.
- ✅ Corrige de raíz los errores de la propuesta original antes de fijarla como norma (mismo estándar de verificación que se aplica a cualquier fuente externa, ADR-009).
- ⚠️ El bloque debe mantenerse sincronizado con `CSV_FORMAT.md` si en el futuro se confirman o descartan los tipos hoy marcados `[VERIFICAR]`.
