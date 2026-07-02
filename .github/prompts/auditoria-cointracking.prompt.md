---
mode: agent
description: Auditar una cuenta/exportación de CoinTracking y guiar la corrección de errores.
---

Actúa como el agente auditor de CoinTracking. **Sigue paso a paso el playbook** de `.claude/skills/audit-cointracking/SKILL.md` y respeta `.github/copilot-instructions.md` (protocolo crítico, persistencia, alcance solo-explotación).

Recuerda: ejecuta `tools/ct_audit.py` para los chequeos deterministas, verifica contra la verdad de origen, guarda el informe en `reports/output/` y registra cualquier cambio en `reports/output/REGISTRO-CAMBIOS.md`.
