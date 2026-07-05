---
id: KB-B2-018
title: "Auditoría de migraciones de exchange — Guía de reconciliación"
level: B
domain: cointracking
source: "Auditoría de agp2025 (Binance EU MiCA, migración a Coinbase en progreso)"
authority: verified
last_verified: 2026-07-05
valid_from: 2026-01-01
valid_until: 2026-12-31
confidence: high
version: 1.0

tags:
  - audit
  - exchange-migration
  - reconciliation
  - procedures

notes: "Procedimiento práctico derivado de auditoría real (agp2025). Usuario: migración Binance→Coinbase por MiCA."
---

# Auditoría de Migraciones de Exchange — Guía de Reconciliación

**Contexto:** Cuando un usuario migra fondos de un exchange a otro (por MiCA, cierre de cuenta, cambio de estrategia, etc.), la auditoría debe verificar que **no aparezcan como venta** (hecho imponible erróneo) ni queden huérfanas (saldo inconsistente).

**Aplicable a:** agp2026 (Binance→Coinbase, julio 2026), y cualquier proyecto futuro con cambios de exchange.

---

## 1. Detección de Migraciones en CoinTracking

### Qué buscar

En la lista de operaciones de CoinTracking, una migración típica aparece así:

**Exchange origen (p. ej. Binance):**
```
Date        | Type       | Cur. | Amount | Fee   | Comments
2026-07-05  | Withdrawal | BTC  | 0.50   | 0.001 | "Transferencia a Coinbase"
2026-07-05  | Withdrawal | ETH  | 10.0   | 0.0   | "Transferencia a Coinbase"
```

**Exchange destino (p. ej. Coinbase):**
```
Date        | Type      | Cur. | Amount | Fee | Comments
2026-07-05  | Deposit   | BTC  | 0.499  | 0.0 | "Transferencia desde Binance"
2026-07-05  | Deposit   | ETH  | 10.0   | 0.0 | "Transferencia desde Binance"
```

### Qué NO buscar (errores comunes)

❌ **Tipo Trade:** Si aparece `Buy` o `Sell` en la operación, NO es una transferencia de cuentas propias → es un hecho imponible real (reexaminar).

❌ **Cantidad inconsistente:** Si los montos no coinciden (p. ej. retira 1 BTC, llega 0.5 BTC sin explicación) → buscar qué pasó en el medio (conversión forzosa, comisión anormal, error de importación).

❌ **Fecha lejana:** Si la retirada es del 1 de julio pero el depósito llega el 15, buscar por qué tardó tanto o si está mal registrada.

---

## 2. Verificación Manual — Checklist de Reconciliación

Para cada operación de retirada/depósito de la migración:

### Paso 1: Emparejar por Tx Hash (Nivel 1 — Fuerte)

**Qué es:** Transaction Hash o ID único de la blockchain, generado por el exchange origen y visible en el destino.

**Cómo verificarlo:**

```
BINANCE:
  1. Login a Binance
  2. Wallet → Withdrawal History → buscar la operación
  3. Copiar "Transaction ID" o "Tx Hash" (hexadecimal largo)

COINBASE:
  1. Login a Coinbase
  2. Movements → Receives → buscar la operación del mismo día/moneda
  3. Copiar "Transaction Hash" (si está visible)

CoinTracking:
  1. En la retirada de Binance: ¿hay un campo "Tx Hash"?
  2. En el depósito de Coinbase: ¿hay un campo "Tx Hash"?
  3. ¿Coinciden?
```

**Resultado:**
- ✅ **Si coinciden:** Transferencia verificada, es seguro emparejarlas.
- ⚠️ **Si no coinciden o falta:** Pasar a Nivel 2 (heurística).

### Paso 2: Emparejar por Heurística (Nivel 2 — Moderado)

Cuando no hay Tx Hash o no coinciden, usar heurística:

**Criterios (deben cumplir TODOS):**

| Criterio | Aceptable | Rojo |
|----------|-----------|------|
| **Moneda** | Exacta (BTC = BTC) | Distinta (BTC ≠ ETH) |
| **Importe retirada** | Importe depósito ≤ (retirada − comisión razonable) | Desproporcionado |
| **Comisión** | Estándar del exchange (~0.1-0.5% o fija) | Anormalmente alta |
| **Ventana temporal** | Depósito dentro de 24h de retirada | > 48h sin explicación |
| **Moneda destino** | Coinbase recibe la misma que Binance envía | Aparece conversión no explicada |

**Ejemplo de aplicación:**

```
Retirada Binance:
  - Fecha: 2026-07-05 14:30 UTC
  - Moneda: BTC
  - Cantidad: 0.50
  - Comisión: 0.001 BTC
  - Total enviado: 0.50 (cantidad neta)

Depósito Coinbase:
  - Fecha: 2026-07-05 15:45 UTC (75 min después ✅)
  - Moneda: BTC
  - Cantidad recibida: 0.499

Análisis:
  - Moneda: BTC = BTC ✅
  - Importe: 0.499 ≈ 0.50 − 0.001 ✅ (comisión stándar)
  - Ventana: 75 min < 24h ✅
  - Conclusión: EMPAREJADA ✓
```

---

## 3. Registro en CoinTracking — Cómo Arreglarlo si Falta el Emparejamiento

Si CoinTracking **no las ha emparejado automáticamente**, pueden quedar como:
- ❌ Una retirada sin depósito = "transferencia huérfana" (saldo negativo incorrecto en destino)
- ❌ Un depósito sin retirada = "dinero de la nada" (compra no fundamentada)

### Solución: Crear una Transferencia Manual en CoinTracking

**Procedimiento:** Ver `knowledge/cointracking/WEB_APP_GUIDE.md` §4bis — "Alta de operaciones manuales: Transferencia entre cuentas propias".

**Ejemplo:**

```
[ Retirada | 2026-07-05 ] BTC 0.50, Binance → Comisión 0.001, Comentario "Migración a Coinbase"
[ Depósito | 2026-07-05 ] BTC 0.499, Coinbase ← Comisión 0, Comentario "Migración desde Binance"
```

**Luego:** Volver a auditar para verificar que los saldos cuadren (sin negativos, sin huérfanas).

---

## 4. Casos Especiales — Conversiones Forzosas

### Qué es

Algunos exchanges obligan a convertir un activo en otro antes de permitir retirada:
- Ejemplo real agp2025: USDT→USDC (Binance, Q1 2025)
- Probable en futuro: activos no conforme MiCA

### Cómo detectarlas

En la secuencia de operaciones antes de una retirada, buscar:

```
Date        | Type  | Cur.  | Amount
2026-07-05  | Trade | USDT  | -5000   → Vender USDT
2026-07-05  | Trade | USDC  | +5000   → Comprar USDC (mismo día, ~mismo monto)
2026-07-05  | Wdraw | USDC  | -5000   → Retirar USDC
```

### Impacto fiscal

**Esto es una PERMUTA cripto-cripto** (Art. 37.1.h LIRPF), **no una transferencia de cuentas propias**. Tributa aunque no pase por EUR:

- **Ganancia:** si USDC vale más que USDT en ese momento
- **Pérdida:** si vale menos

**Qué hacer en auditoría:**
1. Identificar la conversión forzosa (buscará por comentarios o secuencia de Trade-Trade-Withdrawal)
2. Calcular la ganancia/pérdida en EUR (usar `get_gains` o calcular manualmente)
3. Clasificarla **separada de la transferencia** en el informe fiscal
4. Si es confusa, marcar `[VERIFICAR]` y pedir al usuario confirmación

---

## 5. Procedimiento Completo — Flujo de Auditoría de Migración

```mermaid
flowchart TD
    A["Usuario notifica: 'Migré de Binance a Coinbase'"]
    B["Buscar retiradas Binance + depósitos Coinbase <br/>en fecha aproximada"]
    C{"¿Encontradas<br/>ambas operaciones?"}
    D["NO: Buscar más contexto<br/>(fechas alternas, monedas)<br/>Preguntar al usuario"]
    E["SÍ: Intentar emparejar"]
    F{"¿Tx Hash<br/>coincide?"}
    G["Usar heurística:<br/>moneda, importe±comisión,<br/>ventana <24h"]
    H{"¿Heurística<br/>OK?"}
    I["HUÉRFANA: Crear transferencia<br/>manual en CoinTracking"]
    J["Verificar saldos sin negativos"]
    K{"¿Hay Trade<br/>conversión<br/>antes de wdraw?"}
    L["SÍ: Clasificar como permuta<br/>cripto-cripto (hecho imponible)"]
    M["NO: Es transferencia pura<br/>(no imponible)"]
    N["Documentar en REGISTRO-CAMBIOS"]
    O["Cerrado: migración auditada"]

    A → B
    B → C
    C -->|NO| D
    C -->|SÍ| E
    E → F
    F -->|SÍ| J
    F -->|NO| G
    G → H
    H -->|NO| I
    H -->|SÍ| J
    I → J
    J → K
    K -->|SÍ| L
    K -->|NO| M
    L → N
    M → N
    N → O
```

---

## 6. Ejemplo Real: agp2025 (Binance → Coinbase, Julio 2026)

**Estado:** En progreso (usuario iniciando migración por MiCA).

**Lo que esperar en la próxima auditoría (agp2026):**

```
BINANCE (retiradas):
  - Julio 2026: 1+ retiros de activos (BTC, ETH, USDC, etc.)
  - Comisiones normales (~0.1% o estándar por asset)
  - Comentarios posibles: "Salida UE MiCA", o sin comentario

COINBASE (depósitos):
  - Julio 2026: depósitos correspondientes
  - Monedas idénticas, montos cercanos tras comisión
  - Fechas ≤ 24h después de Binance

VERIFICACIÓN:
  1. Buscar Tx Hashes en ambos exchanges
  2. Si faltan, usar heurística (moneda, monto, tiempo)
  3. Si hay conversión USDT→USDC adicional, separarla (es permuta)
  4. Crear transferencias manuales si CoinTracking no las emparejó
  5. Verificar saldos finales en app real (no solo CoinTracking)
```

---

## 7. Checklist Rápido — Antes de Auditar Migraciones

- [ ] ¿El usuario menciona cambio de exchange? → Preguntar fechas aproximadas
- [ ] ¿Hay retiradas sin depósito correspondiente? → Buscar en el segundo exchange
- [ ] ¿Hay depósitos sin retirada clara? → Buscar en exchange origen
- [ ] ¿Montos razonables tras comisiones? → Calcular: depósito ≈ retirada − comisión estándar
- [ ] ¿Hay Trade antes de Withdrawal? → Puede ser conversión forzosa (hecho imponible)
- [ ] ¿Saldos finales coinciden con app real? → Siempre verificar en segunda capa
- [ ] ¿Documentado en REGISTRO-CAMBIOS? → Añadir entrada: qué, por qué, evidencia, antes→después

---

## Referencias

| Documento | Para qué |
|-----------|----------|
| `knowledge/reference/context/EXCHANGE_REGULATORY_UPDATES_2026.md` | Cambios regulatorios (MiCA, conversiones forzosas) |
| `knowledge/cointracking/WEB_APP_GUIDE.md` §4bis | Cómo crear transferencias manuales |
| `tools/ct_audit.py` | Detección automática de transferencias huérfanas |
| `audit-cointracking/SKILL.md`, Paso 1.3 | Procedimiento de verificación de transferencias |
| `reports/output/agp2025/REGISTRO-CAMBIOS.md` | Ejemplo real de migración documentada |

---

**Aplicable a:** agp2026 (cuando usuario finalice migración Binance→Coinbase) y proyectos posteriores.

**Última actualización:** 2026-07-05 (sesión de Fase 3, validación de skills)
