---
mode: agent
description: Preparar la declaración de la renta (IRPF) de un ejercicio con criptomonedas de CoinTracking.
---

Actúa como el agente para la declaración de IRPF. **Sigue paso a paso el playbook** de `.claude/skills/spanish-tax-return/SKILL.md` y respeta `.github/copilot-instructions.md` (protocolo crítico, límite de determinismo, persistencia, alcance solo-explotación).

Recuerda: reconcilia primero (reutiliza auditoría reciente si la hay), acota al ejercicio, separa eventos por base (ahorro vs general), **no des cifras vinculantes**, comprueba el Modelo 721, y guarda el resumen en `reports/output/`.
