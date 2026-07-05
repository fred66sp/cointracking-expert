---
id: KB-B1-018
title: "Pendientes de fundamentar / verificar"
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

# Pendientes de fundamentar / verificar

Backlog único de cuestiones abiertas que afectan a la fiabilidad fiscal (ADR-008/009). El agente **no debe afirmar** estos puntos sin cerrarlos contra fuente oficial. Ordenados por impacto.

## 🔴 Alto impacto (afectan a las cifras que van a Hacienda)

0. **🆕 (2026-07-04) El propio FIFO está judicialmente cuestionado — sigue abierto.** El TSJPV (STSJPV 37/2025 y 41/2025, mismo litigio, recurso 75/2024; régimen foral de Bizkaia, análogo a la LIRPF estatal) ha declarado que las criptomonedas NO son "valores homogéneos" y que el FIFO "no es aplicable", proponiendo el coste real de las unidades vendidas. La DGT no ha rectificado (V0525-25, posterior, mantiene FIFO). Sin unificación por el Tribunal Supremo; no consta si hay recurso de casación. No vincula a la AEAT en territorio común, pero es la primera línea jurisprudencial real en sentido contrario a la doctrina DGT. **Antes de cualquier declaración con cifras significativas, advertir de esta incertidumbre y recomendar validación con un profesional.** *(CAPITAL_GAINS.md §4)*
1. ~~**Ámbito del FIFO: ¿global o por cuenta?**~~ — **Resuelto 2026-07-04** (asumiendo que el FIFO se aplique pese al punto 0): V0525-25 confirma expresamente el ámbito **global**, "con independencia del lugar de custodia" (verificado por triangulación de fuentes secundarias; pendiente confirmar contra texto literal si se logra acceder al buscador oficial de la DGT, que dio error de certificado). Coherente con la opción "Depot/Lot separation" de CoinTracking (desactivada = global). *(CAPITAL_GAINS.md §4)*
2. ~~**Fuente de precios históricos en EUR**~~ — **Resuelto 2026-07-04:** CoinTracking usa promedio ponderado (CMC/Coingecko/WorldCoinIndex), no configurable para cripto; sí hay opción limitada para EUR/USD. Sin fuente oficial homologada por la AEAT — criterio: razonable, consistente y documentable; conservar evidencia extra en importes elevados. *(cointracking/COST_BASIS_AND_VALIDATION.md §4.5)*
3. ~~**Compensación de pérdidas**~~ — **Resuelto 2026-07-04:** orden de compensación, límite del 25 % (vigente sin cambios desde 2018) y arrastre a 4 ejercicios, verificado contra el manual práctico de la AEAT. *(CAPITAL_GAINS.md §7)*
4. **Regla de recompra (Art. 33.5 LIRPF)** — Sin consulta DGT que confirme su aplicación a cripto; la propia V0525-25 (28/03/2025) niega que sean "valores homogéneos" a efectos reglamentarios, pero Renta Web aplicaría en la práctica un plazo de 12 meses igualmente. Contradicción sin zanjar. *(CAPITAL_GAINS.md §7)*

## 🟠 Medio impacto (calificación de rentas)

5. **Staking delegado** (pools/proveedores) — V1766-22 cubre el staking activo; no se pronuncia sobre el delegado. *(CAPITAL_INCOME.md §2)*
6. **Lending / intereses** — Falta consulta DGT específica; encaje por analogía con V1766-22. *(CAPITAL_INCOME.md §3)*
7. **Recompensa / Bonificación** — Clasificación caso a caso (RCM base ahorro vs ganancia patrimonial base general). Airdrop simple ya confirmado (V1948-21/V0648-24); pendiente sin doctrina DGT: bonus de bienvenida, cashback, referidos (ocasional vs profesional/actividad económica), Learn & Earn. *(CAPITAL_INCOME.md §5.1)*

## 🟡 Obligaciones informativas

8. ~~**Modelo 721 y autocustodia**~~ — **Resuelto 2026-07-04:** la autocustodia (Ledger/MetaMask, control propio de claves) queda excluida del 721, confirmado contra FAQ oficial de la AEAT. *(INFORMATIVE_OBLIGATIONS.md §1)*
9. ~~**Norma sancionadora exacta**~~ — **Resuelto 2026-07-04:** régimen general arts. 198/199 LGT (20 €/dato, mín. 300 €, máx. 20.000 €; mitad si es voluntario fuera de plazo). *(INFORMATIVE_OBLIGATIONS.md §2)*

## 📊 Datos del usuario a resolver (no normativa)

10. **USDC: ganancia realizada implausible (+555 €)** — probable base de coste distorsionada; investigar con el método auditar → verdad de origen → corregir. *(auditoría global)*
11. **OM (+1.027 €) y BTC (+495 €)** — confirmar si son ventas reales o base faltante. *(auditoría global)*
12. **Depósitos cripto huérfanos** detectados por `tools/ct_audit.py` (p. ej. AGIX/USDT "Ethereum 1", RDNT internos Binance↔Earn) — revisar contra fuentes de origen.

---

> **Regla:** al tocar cualquiera de estos temas, cerrar contra fuente oficial (AEAT/BOE/DGT o datos reales del usuario) y actualizar el documento correspondiente + este backlog.
