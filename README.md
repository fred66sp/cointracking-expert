# CoinTracking Expert

**Agente de IA auditor de CoinTracking** para reconciliación de criptomonedas y fiscalidad española (IRPF).

El agente vive en Claude Code, se apoya en una base de conocimiento propia y audita los datos de una cuenta de CoinTracking —accediendo por la API (vía MCP) o por el CSV export— para detectar y **explicar** problemas: transferencias huérfanas, ventas sin base de coste, duplicados, saldos imposibles e incoherencias fiscales.

> ⚠️ Herramienta de reconciliación y diagnóstico, **no asesoramiento fiscal**. El agente encuentra y explica; no produce cifras fiscales vinculantes (ver `DECISIONS.md`, ADR-006).

## Cómo funciona

1. Invocas la skill **`/audit-cointracking`** en Claude Code.
2. El agente carga su conocimiento (`knowledge/`), obtiene los datos (MCP en vivo o CSV) y aplica su playbook de chequeos.
3. Devuelve un informe con formato **evidencia → causa → impacto → recomendación**, citando la regla aplicada.

## Estructura

```
.claude/
  agents/cointracking-auditor.md      # El subagente auditor (rol y principios)
  skills/audit-cointracking/          # El playbook invocable (/audit-cointracking)
knowledge/                            # El "cerebro" del agente (fuente de verdad)
  cointracking/                       # Formato CSV, modelo de coste, integración MCP, catálogo
  taxation/spain/                     # Fiscalidad IRPF: ganancias, FIFO, Modelo 721
docs/GLOSSARY.md                      # Glosario de términos
templates/AUDIT_REPORT.md             # Plantilla de informe de auditoría
reports/output/                       # Informes generados (ignorado por git)
DECISIONS.md                          # Registro de decisiones (ADR-001…007)
FOUNDATION.md                         # Principios de ingeniería del proyecto
CLAUDE.md                             # Instrucciones para Claude Code
.mcp.json                             # Configuración del servidor MCP de CoinTracking
```

## Acceso a los datos

Dos vías (ADR-006):

- **MCP de la API de CoinTracking** (datos en vivo, solo lectura). Servidor externo instalado localmente; se configura en `.mcp.json` con credenciales cargadas por `node --env-file` (nunca en el repo). Ver `knowledge/cointracking/MCP_API.md`.
- **CSV export** ("Trade Table"). Ver el formato validado en `knowledge/cointracking/CSV_FORMAT.md`.

## Privacidad y seguridad

- Los datos financieros reales (CSV, informes en `reports/output/`) y las credenciales de la API **nunca** se versionan (excluidos en `.gitignore`).
- El servidor MCP es de solo lectura.

## Estado

El agente está construido, conectado a la API y validado con datos reales. La base de conocimiento cubre el formato de CoinTracking, su modelo de coste y la fiscalidad española; quedan puntos marcados como `[PENDIENTE DE FUNDAMENTAR]` (p. ej. fiscalidad de staking).

## Licencia

Ver [LICENSE](LICENSE).
