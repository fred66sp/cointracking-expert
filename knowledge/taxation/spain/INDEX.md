# Fiscalidad de criptomonedas en España

**Jurisdicción:** España (IRPF — personas físicas residentes)
**Última verificación de fuentes:** 2026-07-02
**Estado:** En construcción

Esta sección documenta el tratamiento fiscal de las criptomonedas ("monedas virtuales") para **personas físicas residentes** en España, según la normativa del IRPF y la doctrina de la AEAT/DGT.

> ⚠️ **Regla del proyecto (FOUNDATION.md):** *nunca se inventan reglas fiscales*. Todo contenido de esta sección debe estar respaldado por una fuente oficial citada (BOE/LIRPF, manuales AEAT, consultas DGT). Lo no verificado se marca explícitamente como **[PENDIENTE DE FUNDAMENTAR]**.
>
> ⚠️ **Aviso:** Este material es una base de conocimiento técnica para reconciliación y cálculo, **no asesoramiento fiscal**. La calificación final de cada operación puede depender de circunstancias particulares y debe validarla un profesional. El framework *genera datos y cálculos trazables*; no presenta declaraciones (ver PROJECT_CHARTER §5, fuera de alcance).

## Alcance

- Personas físicas residentes fiscales en España
- Tributación en el IRPF (no IVA, no IS, no no-residentes en esta primera fase)

## Documentos

- **[CAPITAL_GAINS.md](CAPITAL_GAINS.md)** — Ganancias y pérdidas patrimoniales: venta por fiat, permuta cripto-cripto, valoración, método FIFO, integración en base del ahorro y tramos. *(Fundamentado en AEAT/LIRPF)*
- **[INFORMATIVE_OBLIGATIONS.md](INFORMATIVE_OBLIGATIONS.md)** — Obligaciones informativas: Modelo 721 (criptos en el extranjero) y contexto de los Modelos 172/173. *(Fundamentado en AEAT)*
- **CAPITAL_INCOME.md** — *(pendiente)* Rendimientos del capital mobiliario y otras rentas: staking, lending, intereses, airdrops, minería. Requiere fundamentar la calificación con consultas DGT específicas.

## Relación con los motores

- El **motor FIFO** implementa el criterio de identificación de unidades exigido por la AEAT (ver CAPITAL_GAINS §4).
- El **motor fiscal** consume estas reglas para calcular ganancias/pérdidas y su integración en la base del ahorro.
- El **motor de reportes** produce el detalle trazable por operación (evidencia), no la declaración.

## Fuentes principales

- Agencia Tributaria — Manual práctico IRPF 2025, cap. 11 (ganancias y pérdidas patrimoniales), sección monedas virtuales
- Ley 35/2006 del IRPF (LIRPF): arts. 14, 33–37, 46, 49
- AEAT — Modelo 721 (declaración informativa sobre monedas virtuales situadas en el extranjero)
