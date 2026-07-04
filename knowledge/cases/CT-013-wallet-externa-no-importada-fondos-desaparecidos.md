---
id: KB-C1-013
title: 'CT-013: Wallet externa no importada (fondos "desaparecidos")'
level: C
domain: cointracking
source: Caso auditado en proyecto real
authority: verified
last_verified: '2026-07-03'
valid_from: '2024-01-01'
valid_until: null
confidence: high
version: '1.0'
related_adr:
- ADR-014
- ADR-026
- ADR-004
tags:
- case-study
- transferencias_huerfanas
- cointracking
notes: 'Categoria: transferencias_huerfanas'
---

# CT-013: Wallet externa no importada (fondos "desaparecidos")

**Categoria:** transferencias_huerfanas | **Confianza:** verificado | **Riesgo:** alto

## Sintomas

- Fondos que salieron de un exchange y no aparecen en ningún otro lugar de CoinTracking

## Solucion Recomendada

- Registrar la wallet externa en CoinTracking (manual o vía dirección pública) y añadir el depósito correspondiente.

**Impacto fiscal:** Venta ficticia si se interpreta la retirada como transmisión sin comprobar el destino.
