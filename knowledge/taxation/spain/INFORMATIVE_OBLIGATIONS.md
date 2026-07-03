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

> 🔑 **Implicación para el agente:** el Modelo 721 aplica a saldos en **exchanges/custodios no residentes**. Los activos en **autocustodia** (wallets propias tipo Ledger/MetaMask, donde el usuario tiene sus claves) quedan fuera del 721 según el requisito (1) — punto a matizar y **verificar** por su relevancia (parte de los datos de la cuenta de referencia están en Ledger Live/MetaMask). **[VERIFICAR alcance exacto de autocustodia]**

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

> ⚠️ Sanciones mencionadas en fuentes secundarias (p. ej. 200 € por no presentar, 150 € por presentación incorrecta): **[PENDIENTE DE FUNDAMENTAR con la norma sancionadora exacta]** antes de afirmarlas como definitivas.

---

## 3. Relación con el resto del conocimiento del agente

- El Modelo 721 depende de **tenencias valoradas a 31/12** → datos del MCP (`get_historical_summary`/`get_historical_currency`) + fuente de precios EUR.
- Requiere metadato por cuenta: **jurisdicción del custodio** (residente/no residente) y **tipo** (custodia/autocustodia) → candidato a enriquecer `knowledge/exchanges/`.
- La declaración de IRPF (ganancias/pérdidas) es independiente y se cubre en CAPITAL_GAINS.md.

## 4. Cuestiones abiertas

1. Alcance exacto del 721 respecto a **autocustodia** (§1) — relevante para Ledger/MetaMask de la cuenta de referencia.
2. Norma sancionadora exacta de 172/173/721 (§2).
3. Determinación del **valor a 31/12** (fuente y método de precio) — compartida con la preparación fiscal de `CAPITAL_GAINS.md`.
