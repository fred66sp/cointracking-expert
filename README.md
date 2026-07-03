# CoinTracking Expert

**Agente de IA auditor de CoinTracking** para reconciliación de criptomonedas y fiscalidad española (IRPF).

El agente vive en Claude Code, se apoya en una base de conocimiento propia y audita los datos de una cuenta de CoinTracking —accediendo por la API (vía MCP) o por el CSV export— para detectar y **explicar** problemas (transferencias huérfanas, ventas sin base de coste, duplicados, saldos imposibles, incoherencias fiscales), **guiar al usuario paso a paso a corregirlos en la web de CoinTracking** y preparar la declaración.

> ⚠️ Herramienta de reconciliación y diagnóstico, **no asesoramiento fiscal**. El agente encuentra y explica; no produce cifras fiscales vinculantes (ver `DECISIONS.md`, ADR-006).

## Para usuarios

👉 **¿Quieres usar el agente?** Lee [USER_GUIDE.md](USER_GUIDE.md) — guía paso a paso para auditar tus datos y preparar la declaración fiscal. Hay secciones para usuarios nuevos y experimentados.

---

## Cómo funciona

Le dices lo que quieres y el agente lo enruta:

- **Reconciliar / auditar** los datos → skill **`/audit-cointracking`**.
- **Preparar la declaración de la renta** (IRPF) de un ejercicio → skill **`/spanish-tax-return`** (reconcilia primero y luego prepara lo fiscal).

El agente carga su conocimiento (`knowledge/`), obtiene los datos (MCP en vivo o CSV) y devuelve un informe con formato **evidencia → causa → impacto → recomendación**, citando la regla aplicada.

## Estructura

```
.claude/
  agents/cointracking-auditor.md      # El subagente auditor (rol y principios)
  skills/audit-cointracking/          # Playbook de reconciliación (/audit-cointracking)
  skills/spanish-tax-return/          # Preparación de la declaración IRPF (/spanish-tax-return)
cointracking-mcp/                     # Servidor MCP propio (Go): API de CoinTracking + caché + multi-proyecto
knowledge/                            # El "cerebro" del agente (fuente de verdad)
  cointracking/                       # Formato CSV, coste, integración MCP, guía web (remediación), catálogo
  taxation/spain/                     # Fiscalidad IRPF: ganancias, FIFO, Modelo 721
  exchanges/                          # Contexto regulatorio/operativo de exchanges (p. ej. MiCA)
  patterns/                           # Casos de reconciliación curados (cointracking_casos_v2.yaml)
  blockchains/ · wallets/ · faq/      # Pendientes de poblar (solo INDEX.md por ahora)
docs/GLOSSARY.md                      # Glosario de términos
templates/                            # Plantillas de informe (auditoría, declaración)
tools/ct_audit.py                     # Chequeos deterministas vetados (saldos, transferencias, duplicados…)
tests/fixtures/                       # Caso de prueba de oro (sintético) para regresión del tool
USER_INPUT/<proyecto>/                # Aquí deja el usuario sus CSV/fuentes, por proyecto (ignorado por git, ADR-013)
reports/output/<proyecto>/            # Informes generados, por proyecto (ignorado por git, ADR-013)
AGENT_CHANGE_REQUESTS.md              # Bandeja de peticiones de mejora desde el uso real (Copilot → Claude Code, ADR-012)
DECISIONS.md                          # Registro de decisiones (ADR-001…022 y siguientes)
FOUNDATION.md                         # Principios de ingeniería del proyecto
CLAUDE.md                             # Instrucciones para Claude Code
.github/copilot-instructions.md       # Instrucciones equivalentes para GitHub Copilot (explotación)
.mcp.json / .vscode/mcp.json          # Arranque del servidor MCP propio (cointracking-mcp/dist/)
```

## Acceso a los datos

Dos vías (ADR-006):

- **MCP de la API de CoinTracking** (datos en vivo, solo lectura). Servidor propio en Go (`cointracking-mcp/`, ADR-016), compilado localmente (`dist/cointracking-mcp.exe`); credenciales solo por variables de entorno, nunca en el repo. Ver `knowledge/cointracking/MCP_API.md`.
- **CSV export** ("Trade Table"). Ver el formato validado en `knowledge/cointracking/CSV_FORMAT.md`.

## Privacidad y seguridad

- Los datos financieros reales (CSV, informes en `reports/output/`) y las credenciales de la API **nunca** se versionan (excluidos en `.gitignore`).
- El servidor MCP es de solo lectura.

## Estado

El agente está en uso real: reconciliación completa y declaración de IRPF preparada de principio a fin sobre una cuenta real multi-exchange (proyecto `agp2025`). La base de conocimiento cubre el formato de CoinTracking, su modelo de coste, la fiscalidad española y el contexto regulatorio de exchanges (p. ej. MiCA); quedan puntos marcados como `[PENDIENTE DE FUNDAMENTAR]` o `[VERIFICAR]` según van apareciendo en casos reales. El flujo Claude Code (gestión) / GitHub Copilot (explotación) retroalimenta el conocimiento y las reglas del agente con cada caso auditado (ver `AGENT_CHANGE_REQUESTS.md` y `DECISIONS.md`).

## Quién mantiene y quién usa (ADR-012)

- **Claude Code** mantiene el agente (código, conocimiento, reglas, ADRs, skills, tool). Instrucciones en `CLAUDE.md`.
- **GitHub Copilot** lo explota (auditar, declarar, generar informes) **sin modificarlo**. Instrucciones en `.github/copilot-instructions.md`; MCP en `.vscode/mcp.json`; prompts en `.github/prompts/`. Peticiones de cambio → `AGENT_CHANGE_REQUESTS.md`.

## Licencia

Ver [LICENSE](LICENSE).
