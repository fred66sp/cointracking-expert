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
id: KB-B1-011
title: "Casos Límite de Duplicados: Bots, Flash Loans, Token Splits"
level: B
domain: cointracking
source: "Casos especializados + análisis"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: medium
version: 1.0

related_adr:
  - ADR-014
  - ADR-032

related_docs:
  - knowledge/cointracking/behavioral/DUPLICATE_DETECTION_HEURISTICS.md
  - knowledge/patterns/PATTERN_DUPLICATE_DETECTION.md

tags:
  - cointracking
  - duplicates
  - edge-cases
  - advanced
  - behavioral

notes: "Casos especializados que se parecen a duplicados pero no lo son. Solo para auditorías avanzadas."
---

# Casos Límite de Duplicados: Bots, Flash Loans, Token Splits

## Caso 1: Trading Bots (Compra-Venta en el Mismo Segundo)

**Síntoma:**
```
100 transacciones Buy/Sell de BTC-USDT
Todas en el mismo segundo (p. ej. 2024-03-15 14:23:45)
Cada una con 0.001 BTC
```

**¿Es duplicado?**
❌ **NO.** Son operaciones legítimas de un bot de trading de alta frecuencia.

**¿Cómo verificar?**
- ✓ Cada operación tiene `Order ID` distinto
- ✓ Cada operación tiene `Trade ID` distinto en Binance API
- ✓ Alternancia perfecta: Buy/Sell/Buy/Sell/...
- ✓ Suma de todas = ganancia/pérdida esperada

**Acción en CoinTracking:**
Mantener todas las operaciones. El bot generó cientos de trades reales.

---

## Caso 2: Flash Loans (Préstamo + Devuelve en la Misma Transacción)

**Síntoma:**
```
TX1 (block 12345): Receive 1000 USDC (flash loan)
TX2 (block 12345): Send 1000 USDC + 1 USDC fee (devolución)
Ambas en el mismo bloque (mismo segundo aprox)
Misma cantidad (salvo comisión)
```

**¿Es duplicado?**
❌ **NO.** Son dos transacciones distintas on-chain.

**¿Por qué confunden?**
- Mismo activo
- Mismo monto (casi)
- Mismo timestamp (mismo bloque)
- Parecen reversibles

**Validación:**
- ✓ Direcciones de origen/destino distintas
- ✓ TX IDs distintos en blockchain
- ✓ Una es entrada, otra es salida
- ✓ Comisión de flash loan documenta la intención

**Acción en CoinTracking:**
Mantener ambas. Es una transacción de flash loan legítima (préstamo relámpago).

**Impacto fiscal:**
- Flash loan sin ganancia = no es gravable
- Fee del flash loan = gasto deducible (si se documenta)

---

## Caso 3: Token Split (1:1 → 1:n o 1:n → 1:1)

**Síntoma:**
```
Proyecto XYZ anuncia: "Token split 1:2"
1000 XYZ se convierten en 2000 XYZ automáticamente

CoinTracking muestra:
  TX1: Sell 1000 XYZ
  TX2: Buy 2000 XYZ
  O directamente: "Split 1000 → 2000"
```

**¿Es duplicado?**
❌ **NO.** Es una acción corporativa del token.

**¿Por qué parece raro?**
- La cantidad cambia (1000 → 2000)
- El precio cambia también (mantiene valor total)
- Puede aparecer como dos operaciones

**Validación:**
- ✓ Anuncio oficial del proyecto confirmando el split
- ✓ Ratio 1:1 (1 token antiguo = n tokens nuevos)
- ✓ Precio/token baja proporcionalmente
- ✓ Valor total permanece igual (antes de fees)

**Acción en CoinTracking:**
Mantener ambas TX (si las hay) O registrar como "Split" (si CoinTracking lo soporta).

**Impacto fiscal:**
- Token split = no es venta (no genera ganancia)
- Cost basis se ajusta (1000 @ 100€/token → 2000 @ 50€/token)

---

## Caso 4: Micro-transacciones (Polvo)

**Síntoma:**
```
150 transacciones Buy de 0.00000001 BTC
Todas en el mismo segundo
Mismo precio
```

**¿Es duplicado?**
⚠️ **POSIBLE** (pero raro).

**¿Cuándo es legítimo?**
- Pruebas de dirección (envío de "polvo" para verificar que controlas la wallet)
- Micropagos automáticos
- Cambio fragmentado

**Validación:**
- ¿Trade IDs todos distintos? → Legítimas
- ¿Trade IDs todos iguales? → Probablemente duplicado

**Acción:**
Preguntar al usuario: "¿Reconoces estas micro-transacciones de polvo?"

---

## Caso 5: Reinversión Automática (Staking Compuesto)

**Síntoma:**
```
Cada 24 horas (aprox): Income 0.001 ETH
Automáticamente después: Buy 0.001 ETH (con el ingreso)

100+ pares de (Income + Buy) muy cercanas
```

**¿Es duplicado?**
❌ **NO.** Es reinversión automática (si está documentada).

**Validación:**
- ✓ Income y Buy son operaciones separadas
- ✓ Income = recompensa de staking
- ✓ Buy = reinversión de la recompensa
- ✓ Documentado en el exchange (p. ej. "Auto-compound enabled")

**Acción:**
Mantener ambas. Documentar que es reinversión automática de staking.

---

## Heurística General para Edge Cases

```
Pasos para decidir si es duplicado:

1. ¿Trade ID idéntico en el exchange?
   SÍ → DUPLICADO real → eliminar
   NO → Continuar
   
2. ¿Los campos (excepto timestamp) son EXACTAMENTE idénticos?
   SÍ → Sospechoso, continuar
   NO → Probablemente legítimo → mantener
   
3. ¿Hay contexto que explique por qué aparecen juntas?
   (Bot, flash loan, token split, staking compuesto)
   SÍ → LEGÍTIMAS → mantener
   NO → Probable duplicado → pregunta al usuario
   
4. Ante duda: MANTENER y documentar
   (Mejor tener transacciones de más que perder datos reales)
```

---

## Integración

- **ADR-014:** Protocolo de Trade ID — aplica a todos estos casos
- **DUPLICATE_DETECTION_HEURISTICS.md:** Matriz de decisión principal
