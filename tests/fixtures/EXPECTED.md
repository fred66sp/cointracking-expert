# Caso de prueba de oro — resultados esperados

Fixture: `sample_trades.csv` (datos **sintéticos anonimizados**, no reales). Sirve de prueba de regresión de `tools/ct_audit.py`: si se cambia el tool o la lógica, debe seguir dando estos resultados.

**Cómo ejecutar:**
```
python tools/ct_audit.py tests/fixtures/sample_trades.csv --check all \
  --expect-balances '{"Binance":{"EUR":"0","BTC":"0","DOGE":"-100","USDT":"50"},"Kraken":{"BTC":"0.019","SOL2":"20","ADA":"99.5","XLM":"4"}}'
```

**Debe cumplirse (`validacion.ok = true`):**

| Chequeo | Resultado esperado |
|---|---|
| **Saldos** | Binance: DOGE −100, USDT 50 (EUR 0, BTC 0). Kraken: BTC 0.019, SOL2 20, ADA 99.5, XLM 4. |
| **Saldos negativos** | 1: DOGE −100 en Binance, tipo **cripto** (imposibilidad). |
| **Depósitos huérfanos** | 4: EUR 1000 (fiat=true, externo legítimo), SOL2 50 (cripto), XLM 2 ×2 (cripto). |
| **Retiradas huérfanas** | 0 (la retirada de 0.02 BTC se empareja con el depósito de 0.019 BTC por **Tx Hash**). |
| **Duplicados exactos** | 1 grupo: Depósito 2 XLM ×2. |
| **Colisión de tickers** | `SOL2`. |

**Qué valida cada fallo plantado:**
- Regla de saldo correcta (Σ Compra − Σ Venta, **sin restar Comisión**; el fee de la retirada BTC y la comisión ADA no se doble-cuentan).
- Emparejamiento por Tx Hash (la transferencia BTC no se marca como huérfana).
- Distinción fiat vs cripto en huérfanos y negativos.
- Detección de duplicado exacto por reimportación.
- Colisión de ticker con sufijo (`SOL2` ≠ `SOL`).

---

## Fixture 2: `sample_trades_double_claim.csv` — emparejamiento exclusivo

Regresión del bug encontrado 2026-07-05 (ver `CHANGELOG.md`): el emparejamiento heurístico de transferencias no era exclusivo — dos depósitos idénticos podían "compartir" la misma retirada única, y ninguno se reportaba como huérfano aunque uno de los dos necesariamente lo era.

**Escenario:** 1 retirada de 1 BTC (Binance) + 2 depósitos idénticos de 1 BTC (Kraken a las 12:30, Coinbase a las 12:45), ambos dentro de la ventana temporal y sin Tx Hash.

**Cómo ejecutar:**
```
python tools/ct_audit.py tests/fixtures/sample_trades_double_claim.csv --check transfers
```

**Debe cumplirse:**

| Chequeo | Resultado esperado |
|---|---|
| **Depósitos huérfanos** | 1: el depósito de Coinbase (12:45) — el más tardío pierde el desempate. |
| **Retiradas huérfanas** | 0 (la retirada se empareja con el depósito de Kraken, el más antiguo). |

**Regla de desempate:** cuando varios depósitos podrían matchear con la misma retirada, gana el más antiguo (orden por fecha, no por posición en el CSV) — el resto queda huérfano de verdad, no oculto por un match compartido.
