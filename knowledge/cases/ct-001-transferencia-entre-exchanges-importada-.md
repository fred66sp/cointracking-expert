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
id: "KB-C1-001"
title: "Caso CT-001: Transferencia entre exchanges importada solo en origen"
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

# CT-001: Transferencia entre exchanges importada solo en origen

## Síntomas

- Balance negativo en el exchange destino
- Advertencias de transacciones incompletas
- Holdings inferiores a los reales

## Causa Probable

**Hecho:** Solo existe el registro del retiro (withdrawal) en el exchange de origen.

**Hipótesis:** Falta importar el depósito correspondiente en el exchange destino.

**Supuesto:** El exchange destino no fue importado, o su CSV/API no cubre esa fecha.

## Evidencia Mínima

- Tx Hash (si es on-chain) o referencia de la operación
- Fecha y hora del retiro
- Importe y activo exactos
- Confirmación de que el depósito no aparece en el exchange destino tras revisar su historial completo

## Pasos de Diagnóstico

1. Revisar la transacción en Transacciones (enter_coins.php) del exchange de origen.
1. Buscar el mismo importe y activo en el exchange destino, ampliando el rango de fechas de búsqueda.
1. Comprobar si la importación del destino tiene un filtro temporal que excluya esa fecha.
1. Revisar si la API o el CSV del destino cubren el periodo del depósito.

## Solución Recomendada

- Importar el historial faltante del exchange destino (API o CSV).
- Si no existe registro de origen, crear la transferencia manualmente enlazando ambos lados (ver knowledge/cointracking/WEB_APP_GUIDE.md §2).
- Validar que el balance del activo queda no negativo tras recalcular.

## Evaluación

- **Confianza:** verificado
- **Riesgo:** alto
- **Impacto fiscal:** Puede provocar ventas sin base de coste y balances erróneos que arrastran el cálculo FIFO.

## Señales Tempranas

- Balance negativo justo tras la fecha de una transferencia
- Diferencias entre holdings reales del exchange y los de CoinTracking

## Validación Antes/Después

**Antes:**
- Balance negativo

**Después:**
- Balance reconciliado

- Sin advertencias de transferencia relacionadas

## Vigencia

- **Última revisión:** 2026-07-03
- **Riesgo de caducidad:** Cambios en los importadores API/CSV de exchanges que alteren cómo se detectan retiros/depósitos.
- **Fuente para revalidar:** knowledge/cointracking/CSV_FORMAT.md §7 (Emparejamiento de transferencias) y CoinTracking Help Center
