# Pendientes de fundamentar / verificar

Backlog único de cuestiones abiertas que afectan a la fiabilidad fiscal (ADR-008/009). El agente **no debe afirmar** estos puntos sin cerrarlos contra fuente oficial. Ordenados por impacto.

## 🔴 Alto impacto (afectan a las cifras que van a Hacienda)

1. **Ámbito del FIFO: ¿global o por cuenta?** — ¿Se aplica el FIFO al conjunto del mismo activo del contribuyente (todas las cuentas) o por exchange? Cambia el resultado. Fuente: criterio AEAT/DGT. *(CAPITAL_GAINS.md §4)*
2. **Fuente de precios históricos en EUR** — Necesaria para valorar permutas cripto-cripto y tenencias a 31/12 (Modelo 721). Sin ella no hay cálculo. Decidir origen (CoinTracking vs externa) y método. *(CAPITAL_GAINS.md §5)*
3. **Compensación de pérdidas** — Porcentajes y plazos exactos de compensación en la base del ahorro y arrastre a ejercicios siguientes (Art. 49 LIRPF, evolución por año). *(CAPITAL_GAINS.md §7)*
4. **Regla de recompra (Art. 33.5 LIRPF)** — Aplicabilidad a criptomonedas de la no-cómputo de pérdidas por recompra de activos homogéneos (2 meses / 1 año). *(CAPITAL_GAINS.md §7)*

## 🟠 Medio impacto (calificación de rentas)

5. **Staking delegado** (pools/proveedores) — V1766-22 cubre el staking activo; no se pronuncia sobre el delegado. *(CAPITAL_INCOME.md §2)*
6. **Lending / intereses** — Falta consulta DGT específica; encaje por analogía con V1766-22. *(CAPITAL_INCOME.md §3)*
7. **Recompensa / Bonificación** — Clasificación caso a caso (RCM base ahorro vs ganancia patrimonial base general). *(CAPITAL_INCOME.md §5)*

## 🟡 Obligaciones informativas

8. **Modelo 721 y autocustodia** — Alcance exacto respecto a wallets autocustodiadas (Ledger/MetaMask). *(INFORMATIVE_OBLIGATIONS.md §1)*
9. **Norma sancionadora exacta** de los Modelos 172/173/721. *(INFORMATIVE_OBLIGATIONS.md §2)*

## 📊 Datos del usuario a resolver (no normativa)

10. **USDC: ganancia realizada implausible (+555 €)** — probable base de coste distorsionada; investigar con el método auditar → verdad de origen → corregir. *(auditoría global)*
11. **OM (+1.027 €) y BTC (+495 €)** — confirmar si son ventas reales o base faltante. *(auditoría global)*
12. **Depósitos cripto huérfanos** detectados por `tools/ct_audit.py` (p. ej. AGIX/USDT "Ethereum 1", RDNT internos Binance↔Earn) — revisar contra fuentes de origen.

---

> **Regla:** al tocar cualquiera de estos temas, cerrar contra fuente oficial (AEAT/BOE/DGT o datos reales del usuario) y actualizar el documento correspondiente + este backlog.
