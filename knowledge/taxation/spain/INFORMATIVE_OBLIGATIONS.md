# Obligaciones informativas sobre criptomonedas (España)

**Jurisdicción:** España
**Fuentes:** AEAT — Modelo 721 (declaración informativa sobre monedas virtuales situadas en el extranjero) y preguntas frecuentes; Modelos 172 y 173
**Última verificación:** 2026-07-02
**Vigencia:** umbrales y plazos vigentes a 2026-07. Pueden cambiar por año — reverificar (AEAT) el umbral (50.000 €), la regla de repetición (+20.000 €) y los plazos si esta fecha es antigua (ADR-008).
**Estado:** Fundamentado en fuente oficial

> ⚠️ Base técnica, **no asesoramiento fiscal**. Ver disclaimer en INDEX.md.

Estas son declaraciones **informativas** (no liquidan impuesto), distintas de la tributación en IRPF (ver CAPITAL_GAINS.md).

---

## 1. Modelo 721 — criptomonedas situadas en el extranjero

Declaración informativa **del contribuyente**, en vigor desde el ejercicio **2023**, sobre monedas virtuales custodiadas en el extranjero.

### Quién está obligado

Personas físicas o jurídicas residentes (y otros supuestos como establecimientos permanentes) que ostenten alguna de estas condiciones sobre las criptos: **titular, autorizado, beneficiario, con poder de disposición, o titular real**.

Deben concurrir **dos requisitos**:
1. Las monedas virtuales están custodiadas por **personas o entidades que prestan servicios de salvaguarda de claves criptográficas privadas** en nombre de terceros (es decir, custodios/exchanges), y
2. dichas personas/entidades **no son residentes en España**.

> ✅ **Confirmado 2026-07-04** contra las [preguntas frecuentes oficiales de la AEAT sobre el Modelo 721](https://sede.agenciatributaria.gob.es/Sede/todas-gestiones/impuestos-tasas/declaraciones-informativas/modelo-721-decla-sobre-monedas-extranjero/preguntas-frecuentes-sobre-modelo-721.html): las criptomonedas en **autocustodia** (el propio contribuyente controla las claves privadas — Ledger, Trezor, MetaMask, u otra wallet software/hardware) **NO se declaran en el Modelo 721**, con independencia de su valor y de si es *hot* o *cold wallet*. Lo determinante es **quién controla las claves privadas**, no el dispositivo ni el protocolo. Base normativa: Disposición adicional 18.ª de la Ley 58/2003 (LGT) + art. 42 quater del Reglamento (RD 1065/2007, RGAT) + Orden HFP/886/2023 (aprueba el modelo). El requisito (1) de arriba (custodia por un tercero) es exactamente lo que excluye la autocustodia.
>
> **Implicación para el agente:** los saldos de la cuenta de referencia en Ledger Live/MetaMask (autocustodia) **quedan fuera del ámbito del 721**, aunque superen el umbral de 50.000 €; solo cuentan los saldos en exchanges/custodios no residentes.

### Umbral

Obligación si el **valor conjunto** de esas criptos a **31 de diciembre** supera los **50.000 €**.

### Plazo

Del **1 de enero al 31 de marzo** del año siguiente al ejercicio declarado.

### Periodicidad en años sucesivos

Tras la primera presentación, solo hay que volver a presentar si el **saldo conjunto a 31/12 aumenta en más de 20.000 €** respecto al que determinó la última declaración presentada.

> 🔧 **Implicación para el agente:** para dar soporte al 721 se necesita **valoración de tenencias en EUR a 31/12** por activo y custodio (`cointracking_get_historical_summary`/`get_historical_currency`, ver `knowledge/cointracking/MCP_API.md`), y la clasificación de cada cuenta como **residente / no residente** y **custodia / autocustodia**. La skill `spanish-tax-return` (Paso 5) lo cubre; candidato a enriquecer con un reporte específico si hace falta.

---

## 2. Modelos 172 y 173 — obligaciones de los proveedores (contexto)

**No son obligaciones del inversor particular**, sino de las **entidades residentes en España** que prestan servicios con criptomonedas (exchanges, custodios, monederos, cambistas):

- **Modelo 172** — Declaración informativa sobre **saldos** en monedas virtuales.
- **Modelo 173** — Declaración informativa sobre **operaciones** con monedas virtuales.

Se citan aquí porque explican de dónde obtiene la AEAT información cruzada sobre las operaciones del contribuyente (relevante para la coherencia de datos y trazabilidad), pero **el agente, orientado al inversor persona física, no los genera**.

> ✅ **Confirmado 2026-07-04** — los Modelos 172, 173 y 721 son declaraciones informativas y, desde la reforma de la Ley 5/2022 (que sustituyó el régimen sancionador especial y "confiscatorio" del antiguo Modelo 720, cuestionado por el TJUE), **no tienen sanción específica propia**: se aplica el régimen general de infracciones de la **Ley 58/2003, General Tributaria (LGT)**.
> - **No presentar / fuera de plazo con requerimiento previo (art. 198 LGT):** multa fija de **20 € por cada dato o conjunto de datos** referido a la misma persona/entidad que debiera haberse incluido, con **mínimo 300 € y máximo 20.000 €**.
> - **Presentación voluntaria fuera de plazo, sin requerimiento previo:** la sanción se **reduce a la mitad** → 10 €/dato, mínimo 150 €, máximo 10.000 €.
> - **Datos incorrectos o incompletos:** puede entrar en juego el **art. 199 LGT**, con cuantía según el tipo de incumplimiento; no hay un importe único cerrado — no fijar una cifra concreta sin revisar el caso.

---

## 3. Relación con el resto del conocimiento del agente

- El Modelo 721 depende de **tenencias valoradas a 31/12** → datos del MCP (`get_historical_summary`/`get_historical_currency`) + fuente de precios EUR.
- Requiere metadato por cuenta: **jurisdicción del custodio** (residente/no residente) y **tipo** (custodia/autocustodia) → candidato a enriquecer `knowledge/exchanges/`.
- La declaración de IRPF (ganancias/pérdidas) es independiente y se cubre en CAPITAL_GAINS.md.

## 4. Cuestiones abiertas

1. ~~Alcance exacto del 721 respecto a autocustodia~~ — **Resuelto 2026-07-04** (§1): la autocustodia queda fuera, confirmado contra FAQ oficial de la AEAT.
2. ~~Norma sancionadora exacta de 172/173/721~~ — **Resuelto 2026-07-04** (§2): régimen general arts. 198/199 LGT, cuantías confirmadas.
3. Determinación del **valor a 31/12** (fuente y método de precio) — compartida con la preparación fiscal de `CAPITAL_GAINS.md`, sigue pendiente.
