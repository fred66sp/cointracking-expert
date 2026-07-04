# Conocimiento Operativo de CoinTracking (Nivel B1)

**Ubicación:** `knowledge/cointracking/behavioral/`

**Característica:** Documentos que explican **cómo funciona realmente** CoinTracking, validado contra datos reales.

**Autoridad:** `verified` — contrastado contra casos reales de auditoría, no es normativa

---

## Status

Esta carpeta está **vacía por ahora** (Fase 3). Los documentos se crearán después de que los casos reales (Nivel C1) estén migrando a archivos individuales.

---

## Documentos Pendientes (Fase 3)

- `BALANCE_CALCULATION_ALGORITHM.md` — Cómo evoluciona el saldo en CT
- `PURCHASE_POOL_MECHANICS.md` — Cómo avanza el pool en permutas
- `MISSING_PURCHASE_HISTORY_CAUSES.md` — Por qué aparece, cómo evitarlo
- `DUPLICATE_DETECTION_HEURISTICS.md` — Qué detecta CT automáticamente
- `API_VS_CSV_OVERLAP.md` — Cuándo hay duplicados entre fuentes
- `FEE_HANDLING.md` — Cómo calcula comisiones en tercera moneda

---

## Relación con Nivel A2 (Oficial)

- **A2** = Qué dice CoinTracking (documentación oficial)
- **B1** = Cómo se comporta realmente (observación + casos)

Los documentos B1 **fundamentan y explican** lo que dice A2.

Ejemplo:
- A2 (CSV_FORMAT.md): "Campo `type` puede ser 'Trade', 'Deposit', etc."
- B1 (DUPLICATE_DETECTION_HEURISTICS.md): "Cuando ves dos 'Trade' con misma fecha, CT los marca como sospechosos, pero Trade ID distinto = legítimo (caso FLOKI)"

---

## Próximas Sesiones

**Fase 3:** Derivar estos documentos de los casos reales en `knowledge/cases/`.
