---
id: KB-B1-XXX
title: "Untitled Document"
level: B
domain: cointracking
source: "Internal documentation"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - todo
  - needs-review

notes: "Metadatos agregados automáticamente. Verificar y actualizar conforme ADR-032."
---

---
id: "KB-C1-007"
title: "Caso CT-007: Transferencia interna confundida con venta"
level: "C"
domain: "cointracking"
source: "Análisis de casos reales auditados"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2024-01-01"
valid_until:
confidence: "high"
version: "1.0"
related_adr:
  - ADR-003
  - ADR-009
  - ADR-010
related_docs:
  - knowledge/patterns/INDEX.md
  - knowledge/cointracking/COST_BASIS_AND_VALIDATION.md
tags:
  - case
  - transferencias_huerfanas
  - verified
  - operativo
---

# CT-007: Transferencia interna confundida con venta

## Síntomas

- Ganancia inesperada sin motivo aparente

## Causa Probable

**Hecho:** Una salida de fondos está clasificada como "Trade" (venta) en lugar de "Transfer".

**Hipótesis:** Error de clasificación manual al editar la operación, o importación ambigua.

**Supuesto:** La operación era en realidad un movimiento entre monederos/exchanges propios del mismo usuario.

## Evidencia Mínima

- Historial del exchange/wallet de origen y de destino mostrando el mismo importe y activo
- Titularidad de ambas cuentas confirmada como del mismo usuario

## Pasos de Diagnóstico

1. Revisar el tipo asignado a la operación.
1. Comparar importes y fechas entre origen y destino para confirmar que es el mismo movimiento de fondos.

## Solución Recomendada

- Cambiar el tipo de la operación a "Transfer" (ver knowledge/cointracking/WEB_APP_GUIDE.md §2).

## Anti-patrón

Asumir que toda salida del exchange implica una venta o transmisión.

## Por qué Falso Positivo

Una transferencia entre monederos propios del mismo titular no es una transmisión a efectos fiscales; solo lo es si cambia la titularidad económica del activo.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** alto
- **Impacto fiscal:** Ganancia patrimonial ficticia si no se corrige antes de generar el informe fiscal.

## Señales Tempranas

- Beneficio registrado sin una venta real correspondiente en el exchange

## Validación Antes/Después

**Antes:**
- Clasificado como Trade

**Después:**
- Clasificado como Transfer

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en los tipos de importación soportados o en cómo se infiere el tipo por defecto.
- **Fuente para revalidar:** knowledge/cointracking/WEB_APP_GUIDE.md §2
