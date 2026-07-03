# Fiscalidad de criptomonedas en España

**Jurisdicción:** España (IRPF — personas físicas residentes)
**Última verificación de fuentes:** 2026-07-02
**Estado:** En construcción

Esta sección documenta el tratamiento fiscal de las criptomonedas ("monedas virtuales") para **personas físicas residentes** en España, según la normativa del IRPF y la doctrina de la AEAT/DGT.

> ⚠️ **Regla del proyecto (`FOUNDATION.md`, ADR-009):** *nunca se inventan reglas fiscales*. Todo contenido de esta sección debe estar respaldado por una fuente oficial citada (BOE/LIRPF, manuales AEAT, consultas DGT). Lo no verificado se marca explícitamente como **[PENDIENTE DE FUNDAMENTAR]**.
>
> ⚠️ **Aviso:** Este material es una base de conocimiento técnica para reconciliación y cálculo, **no asesoramiento fiscal**. La calificación final de cada operación puede depender de circunstancias particulares y debe validarla un profesional. El agente *genera datos y análisis trazables*; no presenta declaraciones ni produce cifras vinculantes (ver ADR-006, límite de determinismo).

## Alcance

- Personas físicas residentes fiscales en España
- Tributación en el IRPF (no IVA, no IS, no no-residentes en esta primera fase)

## Documentos

- **[CAPITAL_GAINS.md](CAPITAL_GAINS.md)** — Ganancias y pérdidas patrimoniales: venta por fiat, permuta cripto-cripto, valoración, método FIFO, integración en base del ahorro y tramos. *(Fundamentado en AEAT/LIRPF)*
- **[INFORMATIVE_OBLIGATIONS.md](INFORMATIVE_OBLIGATIONS.md)** — Obligaciones informativas: Modelo 721 (criptos en el extranjero) y contexto de los Modelos 172/173. *(Fundamentado en AEAT)*
- **[CAPITAL_INCOME.md](CAPITAL_INCOME.md)** — Rendimientos y otras rentas: staking (RCM, base del ahorro, DGT V1766-22), lending/intereses, airdrops (ganancia patrimonial, base general, DGT 0018-23), recompensas/referidos (DGT V1948-21) y minería (actividad económica). Incluye la regla clave: el valor al percibir = coste de adquisición futuro. *(Fundamentado en DGT/LIRPF)*
- **[PENDIENTES.md](PENDIENTES.md)** — Backlog único de cuestiones abiertas (FIFO global vs cuenta, precios históricos EUR, staking delegado, autocustodia 721…) que el agente no debe afirmar sin cerrar contra fuente oficial.

## Relación con el agente y el cálculo fiscal (ADR-006)

El agente **no calcula** por sí mismo la cifra fiscal vinculante; aplica estas reglas de forma cualitativa y delega el cálculo exacto:

- **Identificación de unidades FIFO:** exigida por la AEAT (ver `CAPITAL_GAINS.md` §4). El cálculo determinista lo hace el **Informe de Impuestos de CoinTracking** (método FIFO + España) o `cointracking_get_gains(price:"oldest")` como verificación de coherencia (nunca como cifra vinculante).
- **Clasificación de eventos imponibles:** la skill `spanish-tax-return` consume estas reglas para separar ganancias/pérdidas patrimoniales (base del ahorro) de rendimientos (`CAPITAL_INCOME.md`), citando la fuente de cada regla.
- **Detalle trazable por operación:** lo produce el informe de la skill (`templates/TAX_SUMMARY_ES.md`), con evidencia de cada cifra — no sustituye a la declaración ni al Informe de Impuestos oficial.

## Mantenimiento y vigencia (ADR-008)

La normativa fiscal cambia cada año. **Revisar anualmente** (idealmente antes de la campaña de renta) y actualizar la "Última verificación"/"Vigencia" de cada documento:

- **Tramos y tipos** de la base del ahorro y de la base general (`CAPITAL_GAINS.md` §6).
- **Modelo 721:** umbral (50.000 €), regla de repetición (+20.000 €) y plazos (`INFORMATIVE_OBLIGATIONS.md`).
- **Criterios DGT** sobre staking, airdrops, lending, recompensas (`CAPITAL_INCOME.md`): nuevas consultas o matices.
- **Método FIFO** y reglas de compensación de pérdidas: cambios legislativos.

Ante cualquier dato dependiente del ejercicio, el agente debe comprobar la vigencia y, si procede, reverificar contra AEAT/BOE/DGT antes de afirmar.

## Fuentes principales

- Agencia Tributaria — Manual práctico IRPF 2025, cap. 11 (ganancias y pérdidas patrimoniales), sección monedas virtuales
- Ley 35/2006 del IRPF (LIRPF): arts. 14, 33–37, 46, 49
- AEAT — Modelo 721 (declaración informativa sobre monedas virtuales situadas en el extranjero)
