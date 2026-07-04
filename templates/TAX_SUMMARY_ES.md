# Resumen para la declaración de la renta (IRPF) — Ejercicio {AÑO}

- **Perfil:** persona física residente fiscal en España · moneda EUR
- **Ejercicio:** {AÑO} (se presenta en {AÑO+1})
- **Fuente de datos:** {MCP API / CSV Trade Table}
- **Fecha de preparación:** {AAAA-MM-DD}
- **Preparado por:** agente `spanish-tax-return` (Claude Code)

> ⚠️ **No es asesoramiento fiscal ni una declaración.** Documento de preparación y coherencia. Las cifras marcadas como «estimación no vinculante» no sustituyen el Informe de Impuestos de CoinTracking (FIFO/España) ni la revisión de un profesional (ADR-006).

> 🔒 **Nivel de confianza (ADR-009).** Cada dato indica su base: **verificado** (con fuente/datos) · **estimación no vinculante** · **supuesto `[VERIFICAR]`** · **no verificable**. El asesor debe poder rastrear cada cifra a su origen; los huecos se declaran, no se rellenan.

---

## 1. Estado de reconciliación (puerta de calidad)

{✅ Datos reconciliados / ⛔ Bloqueado: hay problemas que distorsionan la base de coste}

- Bloqueantes detectados: {lista con severidad y recomendación, o "ninguno"}
- **Conclusión:** {¿se puede confiar en las cifras del ejercicio, o hay que corregir datos primero?}

## 2. Eventos imponibles del ejercicio {AÑO}

| Tipo                                 | Nº ops | Observaciones                      |
|--------------------------------------|-------:|------------------------------------|
| Ventas a fiat (EUR)                  |        |                                    |
| Permutas cripto-cripto (Art. 37.1.h) |        | tributan aunque no pasen por euros |
| Pagos con cripto                     |        |                                    |

- Excluidos (no imponibles): compras con fiat, holding, transferencias entre cuentas propias.

## 3. Ganancias y pérdidas patrimoniales — base del ahorro (FIFO)

- **Estimación no vinculante:** {resultado del ejercicio, o "no calculable de forma fiable aquí"}
- **Cifra exacta:** generar el **Informe de Impuestos de CoinTracking** con método **FIFO** y jurisdicción **España**, ejercicio {AÑO}.
- Tramos aplicables {AÑO}: ver `knowledge/taxation/spain/CAPITAL_GAINS.md` §6.
- Compensación de pérdidas: {aplicable / pendiente de fundamentar}.

## 4. Rendimientos: staking, recompensas, airdrops, intereses

> 🔴 Calificación fiscal **pendiente de fundamentar** — no calculada.

| Tipo                      | Nº ops | Importe (cuantificado) |
|---------------------------|-------:|------------------------|
| Staking                   |        |                        |
| Recompensa / Bonificación |        |                        |
| Ingresos por intereses    |        |                        |

- **Acción:** fundamentar su tributación (consulta DGT) o validar con profesional antes de declarar.

## 5. Obligación informativa — Modelo 721

- Tenencias en custodios **no residentes** a 31/12/{AÑO}: {valor} €
- ¿Supera 50.000 €?: {sí/no} → {obligado / no obligado}
- Autocustodia (Ledger/MetaMask): `[VERIFICAR]` alcance.
- Plazo: 1 ene–31 mar de {AÑO+1}. Repetición si el saldo sube > 20.000 €.

## 6. Advertencias y próximos pasos

- {Datos a corregir antes de declarar}
- {Puntos pendientes de fundamentar: staking, compensación de pérdidas}
- **Recomendación:** validar con un asesor fiscal y con el Informe de Impuestos oficial de CoinTracking antes de presentar.
