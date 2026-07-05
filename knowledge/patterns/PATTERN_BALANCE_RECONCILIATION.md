---
id: KB-B1-XXX
title: "Untitled Document"
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

---
id: KB-C2-002
title: "Patrón: Reconciliación de Balances (diagnóstico de inconsistencias)"
level: C
domain: cointracking
source: "Generalización de 20 casos reales + ADR-004, ADR-017"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-004
  - ADR-017
  - ADR-031

related_docs:
  - CT-001-transferencia-entre-exchanges-importada-solo-en-origen.md
  - CT-004-balance-negativo-por-orden-cronologico-incorrecto-zona-horaria.md
  - CT-012-balance-negativo-por-importacion-parcial-via-api.md
  - CT-013-wallet-externa-no-importada-fondos-desaparecidos.md

tags:
  - pattern
  - balances
  - reconciliation
  - cointracking

notes: "Patrón para diagnosticar balances negativos o inconsistentes. Orden de verificación: cobertura → importación → zona horaria."
---

# Patrón: Reconciliación de Balances

## Síntomas de Inconsistencia

- **Balance negativo** en un activo o exchange
- **Holdings divergentes** entre CoinTracking y el exchange real
- **Advertencias** de transacciones incompletas
- **Saldo no cuadra** entre fuentes

---

## Árbol de Diagnóstico

```
¿Balance negativo?
  SÍ → Paso 1 (Cobertura de fuentes)
  NO → ¿Divergencia respecto a exchange real?
        SÍ → Paso 2 (Comparar con exchange)
        NO → Probablemente OK
```

---

## Paso 1: Cobertura de Fuentes (Antes que nada)

**Regla crítica (ADR-004):** CoinTracking es internamente consistente = No significa correcto.

Verificar **contra la fuente real** (exchange, banco, blockchain):

### 1.1 ¿Están importadas todas las cuentas/exchanges?

```
Listar en CoinTracking: Accounts → Exchanges
Listar en el exchange real: API / Web

¿Faltan exchanges?
  SÍ → Importar el faltante
  NO → Paso 1.2
```

### 1.2 ¿Está cubierto todo el rango temporal?

```
¿La importación cubre desde el primer depósito del usuario?
  NO → Importar el histórico que falta
  SÍ → Paso 1.3
```

### 1.3 ¿El saldo total de CoinTracking coincide con el del exchange?

```
CoinTracking: Accounts → Balance by exchange → Total
Exchange (Binance/Kraken/etc): Portfolio / Wallet

¿Coinciden (tolerancia: comisiones pending)?
  SÍ → OK. Pasar a Paso 2 si aún hay inconsistencia por activo
  NO → Falta importar algo. Paso 1.1
```

---

## Paso 2: Análisis por Activo

Si el balance total cuadra pero hay negativos por activo:

### 2.1 ¿Hay una compra registrada?

```
Buscar operaciones Buy/Deposit del activo:
  ¿Existe al menos una?
    NO → Missing Purchase History (ADR-009 patrón). Importar.
    SÍ → Paso 2.2
```

### 2.2 ¿La suma de compras >= suma de ventas?

```
Pool de compras = SUM(Buy + Deposit + Rewards/Airdrops)
Pool de ventas = SUM(Sell + Withdraw + Conversiones)

¿Pool compras >= Pool ventas?
  SÍ → OK
  NO → Falta una compra. Paso 2.1 (¿Historial incompleto?)
```

### 2.3 ¿Hay una transferencia sin emparejar?

```
Buscar transacciones de tipo "Transfer" del activo:
  Withdrawal sin Deposit correspondiente?
    SÍ → Huérfana (CT-001 patrón). Importar depósito faltante.
    NO → OK
```

---

## Paso 3: Orden Cronológico y Zona Horaria

**Problema común (CT-004):** Zona horaria hace que compra aparezca DESPUÉS de venta.

### 3.1 Verificar zona horaria

```
CoinTracking: Settings → Timezone
  ¿Está en Europe/Madrid con DST?
    NO → Cambiar a Europe/Madrid
    SÍ → Paso 3.2
```

### 3.2 Verificar orden cronológico del usuario

```
¿Los timestamps en CoinTracking coinciden con el exchange original?
  NO → Recalcular (CoinTracking soporta reiniciar importación)
  SÍ → Paso 4
```

---

## Paso 4: Reasignación Manual

Si tras pasos 1-3 sigue habiendo negativos:

```
Revisar cada operación Sell / Withdraw del activo:
  ¿Hay una compra/depósito anterior que lo justifique?
    NO → Crear la compra manualmente (con evidencia)
    SÍ → Reasignar manualmente en CoinTracking si pool order is wrong
```

---

## Casos de Referencia

| Caso | Síntoma | Causa | Solución |
|------|---------|-------|----------|
| **CT-001** | Balance ↓ en destino | Transferencia no importada en destino | Importar CSV/API del destino |
| **CT-004** | Balance ↓ en activo | Zona horaria (compra después venta) | Cambiar a Europe/Madrid + reiniciar import |
| **CT-012** | Balance ↓ después API import | Importación parcial (falta histórico) | Reimportar desde fecha anterior |
| **CT-013** | Balance ↓ inexplicable | Wallet externa sin importar | Crear cuenta en CT + importar CSV |

---

## Checklist de Diagnóstico

- [ ] ¿Están todos los exchanges importados?
- [ ] ¿Cubre el rango temporal desde el primer depósito?
- [ ] ¿El total de CoinTracking = total del exchange?
- [ ] ¿Hay compras registradas para el activo negativo?
- [ ] ¿La suma de compras >= suma de ventas?
- [ ] ¿Hay transferencias sin emparejar?
- [ ] ¿Zona horaria = Europe/Madrid?
- [ ] ¿El orden cronológico es correcto?

---

## Regla de Oro

> **Un balance negativo siempre indica datos incompletos, no un error de cálculo.**
>
> CoinTracking puede estar internamente consistente pero externamente incorrecto.
> Siempre verificar contra el exchange real.

---

## Integración con ADRs

- **ADR-004:** Reconciliación con datos reales → Este patrón lo operacionaliza
- **ADR-017:** Orden de diagnóstico → Este patrón sigue ese orden (cobertura primero)
- **ADR-005:** Zona horaria → Incluye verificación explícita
