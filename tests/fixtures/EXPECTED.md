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
