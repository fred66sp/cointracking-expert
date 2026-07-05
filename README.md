# CoinTracking Expert

**Agente de IA auditor de CoinTracking** para reconciliación de criptomonedas y fiscalidad española (IRPF).

El agente vive en Claude Code, se apoya en una base de conocimiento propia y audita los datos de una cuenta de CoinTracking —accediendo por la API (vía MCP) o por el CSV export— para detectar y **explicar** problemas (transferencias huérfanas, ventas sin base de coste, duplicados, saldos imposibles, incoherencias fiscales), **guiar al usuario paso a paso a corregirlos en la web de CoinTracking** y preparar la declaración.

> ⚠️ Herramienta de reconciliación y diagnóstico, **no asesoramiento fiscal**. El agente encuentra y explica; no produce cifras fiscales vinculantes (ver `DECISIONS.md`, ADR-006).

## 🚀 Inicio Rápido (P0-P3 Completado)

**El sistema está 100% operacional.** Si es tu primera vez:

1. **5 minutos:** Lee [knowledge/QUICK_START.md](knowledge/QUICK_START.md)
2. **Necesitas algo:** [knowledge/NAVIGATION_MAP.md](knowledge/NAVIGATION_MAP.md) (busca por función)
3. **Tienes un problema:** [knowledge/TROUBLESHOOTING_INDEX.md](knowledge/TROUBLESHOOTING_INDEX.md) (busca por síntoma)
4. **Referencia rápida:** [knowledge/CHEAT_SHEET.md](knowledge/CHEAT_SHEET.md) (1 página, operaciones comunes)

**Documentación operativa para desarrolladores:**
- 📖 [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) — compilar/arrancar servidor MCP
- 📖 [knowledge/KNOWLEDGE_MAINTENANCE.md](knowledge/KNOWLEDGE_MAINTENANCE.md) — crear/mantener documentos
- 📖 [GOVERNANCE_WORKFLOW.md](GOVERNANCE_WORKFLOW.md) — registrar decisiones (ADRs)

**Detalles técnicos:**
- 🗺️ [knowledge/INDEX_MASTER.md](knowledge/INDEX_MASTER.md) — estructura de 6 niveles A-F
- 📋 [knowledge/CHEAT_SHEET.md](knowledge/CHEAT_SHEET.md) — operaciones comunes + fórmulas

---

## Para desarrolladores / Contribuyentes

### Protocolo de Desarrollo (ADR-036/37/38)

**Antes de modificar la base de conocimiento, lee:**
1. 📋 [adr/0036-convencion-de-ids-de-documentos.md](adr/0036-convencion-de-ids-de-documentos.md) — Convención de IDs (KB-[NIVEL][SUB]-NNN)
2. 📋 [adr/0037-validacion-obligatoria-en-desarrollo.md](adr/0037-validacion-obligatoria-en-desarrollo.md) — Validación DURANTE creación, no después
3. 📋 [adr/0038-criterio-auditoria-lotes-no-iterativa.md](adr/0038-criterio-auditoria-lotes-no-iterativa.md) — Auditoría MEGA completa

**Workflow para crear/modificar documentos:**
```bash
# 1. Crear documento con YAML completo (ver ADR-036)
# 2. Escribir contenido
# 3. VALIDAR ANTES de commit:
python tools/audit_mega_complete.py

# 4. Si hay errores: FIX y re-valida
# 5. Si está limpio: git commit (hook valida automáticamente)
```

**Instalar pre-commit hook (validación local):**
```bash
powershell -ExecutionPolicy Bypass .\tools\install_hooks.ps1
# O en bash:
chmod +x .git/hooks/pre-commit
```

**CI/CD automático:** GitHub Actions valida en cada push (ver `.github/workflows/audit-mega.yml`)

---

## Para usuarios

👉 **¿Quieres usar el agente?** Lee [knowledge/QUICK_START.md](knowledge/QUICK_START.md) — guía paso a paso para auditar tus datos y preparar la declaración fiscal. Hay secciones para usuarios nuevos y experimentados.

---

## Cómo funciona

Le dices lo que quieres y el agente lo enruta:

- **Reconciliar / auditar** los datos → skill **`/audit-cointracking`**.
- **Preparar la declaración de la renta** (IRPF) de un ejercicio → skill **`/spanish-tax-return`** (reconcilia primero y luego prepara lo fiscal).

El agente carga su conocimiento (`knowledge/`), obtiene los datos (MCP en vivo o CSV) y devuelve un informe con formato **evidencia → causa → impacto → recomendación**, citando la regla aplicada.

## ⚡ Próximos Pasos Opcionales

### P4.2: Testear Skills
```
1. Auditar cuenta con: /audit-cointracking
2. Preparar IRPF con: /spanish-tax-return
3. Documentar cualquier inconveniente o mejora
```

### P4.3: Ampliar Conocimiento
```
- Wallets específicas (Ledger, Trezor, MetaMask)
- Altcoins menos comunes
- Fiscalidad de otros países (UK, US)
- Casos de uso especiales (herencias, empresas)
```

---

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

## Estado (2026-07-05)

**🟢 SISTEMA 100% OPERACIONAL**

### Validación (P0)
- ✅ 68 documentos YAML validados (metadatos íntegros)
- ✅ Estructura A-F completamente documentada (111+ archivos)
- ✅ Archivos críticos verificados

### Navegabilidad (P1)
- ✅ QUICK_START.md — entrada para usuarios nuevos (5 min)
- ✅ NAVIGATION_MAP.md — búsqueda por función
- ✅ TROUBLESHOOTING_INDEX.md — búsqueda por síntoma
- ✅ CHEAT_SHEET.md — referencia rápida

### Infraestructura (P2)
- ✅ DEPLOYMENT_GUIDE.md — compilar/arrancar MCP
- ✅ knowledge/KNOWLEDGE_MAINTENANCE.md — mantener conocimiento
- ✅ GOVERNANCE_WORKFLOW.md — registrar decisiones (ADRs)

### Integración (P3)
- ✅ MCP funcional (servidor Go compilado)
- ✅ Proyecto `agp` activo y auditado
- ✅ 19,229.35 EUR en 39 activos, 500 transacciones
- ✅ +473.94 EUR de ganancia neta (FIFO) verificada

### Auditoría Real (P4.1)
- ✅ Proyecto `agp` completamente auditado
- ✅ Reporte generado: [reports/output/agp/AUDIT_REPORT_COMPLETE_2026-07-05.md](reports/output/agp/AUDIT_REPORT_COMPLETE_2026-07-05.md)
- ✅ Listo para preparar IRPF 2025

**El agente está en uso real:** reconciliación completa y declaración de IRPF preparadas sobre cuentas reales multi-exchange. La base de conocimiento cubre el formato de CoinTracking, su modelo de coste, la fiscalidad española y el contexto regulatorio de exchanges (p. ej. MiCA). El flujo Claude Code (gestión) / GitHub Copilot (explotación) retroalimenta el conocimiento con cada caso auditado (ver `AGENT_CHANGE_REQUESTS.md`, `adr/`, y `CHANGELOG.md`).

## Quién mantiene y quién usa (ADR-012)

- **Claude Code** mantiene el agente (código, conocimiento, reglas, ADRs, skills, tool). Instrucciones en `CLAUDE.md`.
- **GitHub Copilot** lo explota (auditar, declarar, generar informes) **sin modificarlo**. Instrucciones en `.github/copilot-instructions.md`; MCP en `.vscode/mcp.json`; prompts en `.github/prompts/`. Peticiones de cambio → `AGENT_CHANGE_REQUESTS.md`.

## Licencia

Ver [LICENSE](LICENSE).
