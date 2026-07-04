# Cheat Sheet — Referencia Rápida del Auditor

**Para:** Cuando conoces el sistema y solo necesitas recordar un detalle rápido  
**Tiempo:** 1 min por búsqueda

---

## ⚡ Operaciones Comunes

### Depósito (Fiat o Cripto)

```
Qué es:     Entrada de dinero a tu exchange/wallet
En CT:      Type = "Deposit"
Cost basis: SÍ (valor fiat si es fiat; precio en exchange si cripto)
Fiscal:     Base de coste (origen del dinero)
Riesgo:     Verificar que el banco lo confirma
```

**Checklist:**
- [ ] Amount cuadra con banco (fiat) o blockchain (cripto)
- [ ] Fecha coincide
- [ ] Si hay fee (Binance cobra 0.0005 BTC en retirada), está registrada

---

### Compra (Trade Buy)

```
Qué es:     Compras cripto con fiat o otra cripto
En CT:      Type = "Trade", Coin_from = fiat/cripto, Coin_to = cripto
Cost basis: SÍ (el precio que pagaste)
Fiscal:     Ganancia patrimonial (si la vendes después)
Riesgo:     Cost basis incorrecto = Declaración falsa
```

**Checklist:**
- [ ] Amount coincide con exchange
- [ ] Price coincide con exchange (o blockchain si DeFi)
- [ ] Fee está incluida o separada (varía por exchange)
- [ ] Cost basis = (Amount × Price) + Fee

---

### Venta (Trade Sell)

```
Qué es:     Vendes cripto por fiat o cripto
En CT:      Type = "Trade", Coin_from = cripto, Coin_to = fiat/cripto
Cost basis: SE CALCULA (FIFO) — CRÍTICO
Fiscal:     Ganancia/pérdida patrimonial = Precio_venta - Cost_basis
Riesgo:     Si no hay compra previa → VENTA SIN ORIGEN (CT-002, CT-017)
```

**Checklist:**
- [ ] ¿Hay compras previas de este activo? → Si NO: ERROR CRÍTICO
- [ ] Cost basis se calculó con FIFO (oldest first)
- [ ] Ganancia/pérdida = Precio_venta - Cost_basis
- [ ] La ganancia se registró en fiscal

---

### Staking (Deposit + Rewards)

```
Qué es:     Bloqueas cripto, recibes rewards
En CT:      Deposit (principal) + Interest (rewards)
Cost basis: Principal SÍ, Interest NO (es ingreso)
Fiscal:     Principal = base de coste. Interest = ingresos del capital (Modelo 721)
Riesgo:     Confundir Interest con Deposit (CT-005)
```

**Checklist:**
- [ ] Principal aparece como "Deposit" o "Staked"
- [ ] Interest aparece como "Interest" o "Reward"
- [ ] Cost basis del principal = cantidad × precio en momento de stake
- [ ] Interest se declara en ingresos (no es ganancia patrimonial)

---

### Airdrop (Income)

```
Qué es:     Recibes cripto gratis (token distribution)
En CT:      Type = "Income" o "Airdrop" (idealmente)
Cost basis: DEPENDE (¿cuándo entra el control fiscal? Ver AEAT/DGT)
Fiscal:     Potencialmente ingresos del capital (Modelo 721) o ganancia patrimonial
Riesgo:     Clasificación incierta en España (CT-010) — [VERIFICAR]
```

**Checklist:**
- [ ] ¿Está bien clasificado como "Airdrop" o "Income"? Si no → Corregir
- [ ] ¿Hay documentación de cuándo y cómo lo recibiste?
- [ ] Marca como `[VERIFICAR]` si la fisicalidad no está documentada en conocimiento
- [ ] Si tienes duda, consulta AEAT: "Tratamiento fiscal de airdrop de criptomonedas"

---

### Lending (Deposit + Interest)

```
Qué es:     Prestas cripto en un protocolo, recibes interest
En CT:      Deposit (loan) + Interest (returns)
Cost basis: Principal SÍ, Interest NO (es ingreso)
Fiscal:     Principal = base de coste. Interest = ingresos
Riesgo:     Confundir con transferencia interna (CT-011)
```

**Checklist:**
- [ ] Principal está marcado como "Lending" o similar
- [ ] Interest está separado
- [ ] Verificar que withdrawal devuelve el principal + interest

---

### DeFi Swap (On-chain)

```
Qué es:     Cambias cripto A por cripto B en un DEX
En CT:      Puede aparecer como 1 operación o 2+ (depende del DEX)
Cost basis: Cripto A sale al precio spot. Cripto B entra al precio spot
Fiscal:     Ganancia/pérdida en Cripto A
Riesgo:     Fragmentación (CT-015), slippage, fees en gas
```

**Checklist:**
- [ ] Amount A que sales coincide con blockchain
- [ ] Amount B que recibes coincide con blockchain
- [ ] Fee de gas está registrado o incluido
- [ ] Price_out / Price_in valida contra blockchain explorer
- [ ] Si está fragmentado: agrupa por TX hash antes de auditar

---

### Transferencia (Interna o Externa)

```
Qué es:     Mueves cripto de una wallet/exchange a otra tuya
En CT:      Withdrawal + Deposit (2 operaciones)
Cost basis: NO (es movimiento, no compra/venta)
Fiscal:     NO genera ganancia/pérdida
Riesgo:     Desaparecer datos (CT-013), mal cronometraje, no emparejar
```

**Checklist:**
- [ ] Hay UN withdrawal en origen con Amount = A
- [ ] Hay UN deposit en destino con Amount ≈ A - Fee_blockchain
- [ ] Fechas: withdrawal ≤ deposit (blockchain explorer valida)
- [ ] NO aparece como venta/compra

---

### Comisión (Fee)

```
Qué es:     Comisión de transacción (Binance, blockchain, etc)
En CT:      Type = "Fee" o incluida en otra operación
Cost basis: SÍ en la operación asociada (no es deducible aislado)
Fiscal:     Se suma al cost basis de compra; se resta del price en venta
Riesgo:     Fee en moneda ajena (CT-009) — puede no estar en cripto
```

**Checklist:**
- [ ] ¿Fee está por separado o incluido? (varía por exchange)
- [ ] ¿En qué moneda está? (EUR, BTC, BUSD, etc)
- [ ] Si moneda ajena: convierte a cripta target para cost basis
- [ ] Cost basis = Amount + Fee

---

## 🔢 Fórmulas Clave

### Cost Basis (FIFO)

```
Compras:  100 BTC @ 30k = 3M
          50 BTC @ 35k = 1.75M

Vendo 60 BTC:
  - 60 primeros BTC usan las compras más antiguas
  - 60 = 100 (compra 1, parcial) → 60 × 30k = 1.8M

Cost Basis = 1.8M
```

**En CT:** Clic en operación → "Sold Coins" → ver desglose FIFO

---

### Ganancia Patrimonial

```
Venda 1 BTC @ 50k
Cost Basis 35k
Ganancia = 50k - 35k = 15k

Fiscal: +15k en Ganancias de Capital (Modelo 100 del IRPF)
```

---

### Ingresos del Capital (Staking, Interest)

```
Staking: 10 BTC, recibes 0.5 BTC de rewards
Ingresos = 0.5 BTC × precio_spot_cuando_recibiste (p.ej. 50k) = 25k

Fiscal: +25k en Ingresos del Capital (Modelo 721)
```

---

## 📋 Checklists Rápidos

### Antes de Cada Operación

- [ ] ¿Está el tipo correcto? (Deposit, Trade, Interest, Fee, etc)
- [ ] ¿El amount coincide con exchange/blockchain?
- [ ] ¿La fecha coincide?
- [ ] ¿Cost basis está correcto?
- [ ] ¿Fee está incluida o separada?

### Antes de la Auditoría Completa

- [ ] ¿Balance de cada cripto coincide con exchange?
- [ ] ¿Hay duplicados? (Comprobar Trade ID)
- [ ] ¿Hay ventas sin compra? (Cost basis = 0)
- [ ] ¿Hay transferencias desapareadas? (Depósito sin retiro)
- [ ] ¿Hay advertencias técnicas? (Ver CT-020)

### Antes de la Declaración

- [ ] ¿Todas las ganancias están calculadas con FIFO?
- [ ] ¿Todas las ingresos (staking, interest, airdrops) están en Modelo 721?
- [ ] ¿Balance a 31/12 coincide con realidad?
- [ ] ¿Están todas las transferencias reconciliadas?

---

## 🚨 Rojo = Acción Inmediata

| Síntoma | Acción |
|---------|--------|
| Venta sin cost basis | → CT-002, CT-017: buscar compra |
| Balance negativo | → FLOW_NEGATIVE_BALANCE |
| Dos operaciones idénticas | → DUPLICATE_DETECTION_HEURISTICS (verificar Trade ID) |
| Fee en moneda ajena | → FEE_HANDLING |
| Staking como depósito | → STAKING_MECHANICS (reclasificar) |
| Airdrop como compra | → AIRDROPS_MECHANICS (reclasificar) |

---

## 📞 Referencias Rápidas

| Necesito | Voy a |
|----------|-------|
| Entender un tipo de operación | [Tabla anterior en este doc] |
| Ver un caso real | [CT-001](cases/ct-001-transferencia-entre-exchanges-importada-.md) a [CT-020](cases/ct-020-advertencia-tecnica-no-es-error-fiscal.md) |
| Saber qué hacer | [TROUBLESHOOTING_INDEX.md](TROUBLESHOOTING_INDEX.md) |
| Navegar el sistema | [NAVIGATION_MAP.md](NAVIGATION_MAP.md) |
| Definición de términos | [GLOSSARY.md](reference/GLOSSARY.md) |
| Pasos detallados | [PROCEDURE_AUDIT_ACCOUNT.md](procedures/PROCEDURE_AUDIT_ACCOUNT.md) |

---

## 💾 Imprímelo

Este documento es una página. Perfecta para tener a mano mientras auditas.

**Link:** [knowledge/CHEAT_SHEET.md](CHEAT_SHEET.md)
