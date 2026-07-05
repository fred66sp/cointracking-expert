# Hoja de Ruta de Cobertura — Qué Falta Documentar

**Documento:** Plan de expansión de la base de conocimiento  
**Fecha:** 2026-07-05  
**Estado:** Priorizado y listo para implementar

---

## 🎯 Propósito

Identificar qué áreas quedan pendientes en la base de conocimiento para expandir el sistema más allá del caso base (Binance, CoinTracking, España IRPF).

---

## 📊 Estado Actual de Cobertura

### ✅ CUBIERTO (100%)

**Exchanges:**
- Binance (Spot, Margin, Futures, Earn, Convert)
- Kraken (Staking)
- Coinbase (Advanced Trade, Staking)
- CoinTracking (formato CSV, integración MCP)

**Blockchains:**
- Ethereum (TX types, gas, smart contracts)
- Bitcoin (UTXO, consolidations)
- Other chains (Polygon, BSC, Solana, Arbitrum)
- Bridges & Wrapping

**Operaciones:**
- Trades (compra/venta)
- Transfers (depósitos/retiradas)
- Staking (depósitos + rewards)
- Airdrops (importación, clasificación)
- DeFi Swaps (slippage, fees)
- Lending (depósitos, interest)
- Fees (comisiones en tercera moneda)
- Mining (rewards)

**Fiscalidad España:**
- IRPF (Modelo 100: ganancias patrimoniales)
- Modelo 721 (patrimonio)
- FIFO (método de coste)
- Capital gains vs income classification
- Staking classification (incierto, marcado)

**Metodología:**
- Auditoría (6 fases)
- Duplicate detection (Trade ID verification)
- Missing purchase history (5 causas)
- Balance reconciliation
- Transfer matching

---

## ❌ NO CUBIERTO (Prioritario)

### Tier 1: Exchanges Adicionales

| Exchange | Prioridad | Razón | Esfuerzo |
|----------|-----------|-------|----------|
| **BingX** | ALTA | Usado en `agp`, no documentado específicamente | 1-2h |
| **Bybit** | MEDIA | Exchange creciente, estructura similar a Binance | 2-3h |
| **OKX** | MEDIA | Exchange grande en Asia, modelo complejo | 3-4h |
| **Crypto.com** | BAJA | Ecosistema específico (Earn, Visa) | 2-3h |
| **XT.com** | BAJA | Exchange pequeña, API estándar | 1-2h |

**Documentos a crear:**
- `knowledge/cointracking/behavioral/BINGX_MECHANICS.md` (KB-B2-010)
- `knowledge/cointracking/behavioral/BYBIT_MECHANICS.md` (KB-B2-011)
- `knowledge/cointracking/behavioral/OKX_MECHANICS.md` (KB-B2-012)
- etc.

### Tier 2: Wallets y Dispositivos Hardware

| Wallet | Prioridad | Razón | Esfuerzo |
|--------|-----------|-------|----------|
| **Ledger Live** | ALTA | Usado en `agp`, parcialmente documentado | 2-3h |
| **Trezor** | MEDIA | Popular, similar a Ledger | 2-3h |
| **MetaMask** | MEDIA | DeFi principal, importación compleja | 3-4h |
| **Trust Wallet** | BAJA | Ecosistema Binance | 2-3h |
| **Multisig (Gnosis Safe)** | BAJA | Casos avanzados | 3-4h |

**Documentos a crear:**
- `knowledge/wallets/LEDGER_LIVE_INTEGRATION.md` (KB-B4-001)
- `knowledge/wallets/TREZOR_INTEGRATION.md` (KB-B4-002)
- `knowledge/wallets/METAMASK_INTEGRATION.md` (KB-B4-003)
- etc.

### Tier 3: Altcoins y Casos Especiales

| Caso | Prioridad | Razón | Esfuerzo |
|------|-----------|-------|----------|
| **Stablecoins (USDC, USDT, DAI)** | ALTA | Frecuentes, tratamiento especial | 1-2h |
| **Wrapped tokens (wBTC, wETH)** | MEDIA | DeFi, bridging | 1-2h |
| **Tokens con splits (SHIB, etc)** | MEDIA | Cambios en supply | 1-2h |
| **Altcoins sin soporte CT** | BAJA | Importación manual | 2-3h |
| **Mega-cap alts (SOL, ADA, XRP)** | BAJA | Documentación de casos reales | 1-2h per |

**Documentos a crear:**
- `knowledge/cointracking/reference/STABLECOINS_GUIDE.md` (KB-E3-001)
- `knowledge/cointracking/reference/TOKEN_SPLITS_HANDLING.md` (KB-E3-002)
- `knowledge/cointracking/reference/WRAPPED_TOKENS.md` (KB-E3-003)

### Tier 4: Fiscalidad de Otros Países

| País | Prioridad | Razón | Esfuerzo |
|------|-----------|-------|----------|
| **🇬🇧 UK (CGT, Income Tax)** | MEDIA | Brexit cambió régimen, HMRC guidance compleja | 4-5h |
| **🇺🇸 USA (IRS, Wash Sale)** | MEDIA | Diferentes estados, complicado | 5-6h |
| **🇩🇪 Alemania (Spekulationsfrist)** | BAJA | Régimen diferente (1 año hold) | 3-4h |
| **🇨🇭 Suiza (Cantonal tax)** | BAJA | Muy variable por cantón | 4-5h |
| **🇸🇬 Singapur (Gain vs Income)** | BAJA | Régimen simple pero diferente | 2-3h |

**Documentos a crear:**
- `knowledge/taxation/uk/CAPITAL_GAINS_TAX.md` (KB-A4-001)
- `knowledge/taxation/usa/IRS_REPORTING.md` (KB-A4-002)
- `knowledge/taxation/germany/SPEKULATIONSFRIST.md` (KB-A4-003)
- etc.

### Tier 5: Casos Avanzados

| Caso | Prioridad | Razón | Esfuerzo |
|------|-----------|-------|----------|
| **Herencias (inherited crypto)** | BAJA | Régimen especial, base del ahorro | 2-3h |
| **Empresas (business crypto)** | BAJA | Contabilidad diferente (balance sheet) | 3-4h |
| **Institutional holdings** | BAJA | Fondos, custodias | 2-3h |
| **DeFi yield farming** | BAJA | LP tokens, impermanent loss | 3-4h |
| **Options & Derivatives** | BAJA | Futuros, opciones (muy complejo) | 5-6h |

---

## 📈 Priorización Recomendada

### Fase 5 (Próximas 2 Semanas)

**Tier 1 (Exchanges):**
1. BingX (2h) — usado en `agp`, impacto inmediato
2. Bybit (3h) — estructura similar a Binance, reutilizable

**Tier 2 (Wallets):**
3. Ledger Live (3h) — usado en `agp`, parcialmente documentado
4. MetaMask (4h) — DeFi, high demand

**Total: ~12 horas**

### Fase 6 (Semanas 3-4)

**Tier 2 (Wallets):**
5. Trezor (3h)
6. Trust Wallet (3h)

**Tier 3 (Altcoins):**
7. Stablecoins guide (2h)
8. Wrapped tokens (2h)

**Total: ~10 horas**

### Fase 7+ (Después de Semana 4)

**Tier 3-4:**
- Token splits handling
- Altcoins sin soporte CT
- Fiscalidad otros países (UK, USA)

**Tier 5:**
- Casos avanzados (herencias, empresas, DeFi farming)

---

## 🛠️ Estructura para Nuevos Documentos

### Patrón para Exchange

```markdown
---
id: KB-B2-010  (siguiente disponible)
title: "[Exchange Name] Mechanics"
level: B
domain: cointracking
source: [Official docs / CoinTracking support]
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-05
confidence: medium
version: 1.0
---

# [Exchange] Mechanics

## Características Principales
...

## Integración con CoinTracking
...

## Casos Límite
...

## Referencias
```

### Patrón para Wallet

```markdown
---
id: KB-B4-001
title: "[Wallet] Integration Guide"
level: B
domain: cointracking
source: [Official docs / case analysis]
authority: verified
...
---

# [Wallet] Integration

## Tipos de Operaciones Soportadas
...

## Importación en CoinTracking
...

## Validación y Limitaciones
...
```

### Patrón para Fiscalidad

```markdown
---
id: KB-A4-001
title: "[País] Capital Gains Tax"
level: A
domain: taxation
source: [HMRC/IRS/official government]
authority: official
...
---

# [País] Taxation Framework

## Base de Coste
...

## Reporting Requirements
...

## Ejemplos
```

---

## 📊 Matriz de Impacto vs Esfuerzo

```
ALTO IMPACTO, BAJO ESFUERZO:
  ✓ BingX (2h) — usado en agp
  ✓ Ledger Live (3h) — usado en agp
  ✓ Stablecoins guide (2h) — frecuente

ALTO IMPACTO, ALTO ESFUERZO:
  ~ Bybit (3h) — popular pero no crítico
  ~ MetaMask (4h) — DeFi pero avanzado
  ~ UK Taxation (5h) — mercado grande

BAJO IMPACTO, BAJO ESFUERZO:
  ~ XT.com (2h) — exchange pequeña
  ~ Token splits (2h) — casos raros

BAJO IMPACTO, ALTO ESFUERZO:
  ✗ Herencias (3h) — muy específico
  ✗ Options (6h) — muy complejo
  ✗ USA Taxation (6h) — muy variable
```

---

## ✅ Checklist de Implementación

Para cada documento nuevo:

- [ ] Crear archivo en carpeta correcta (behavioral, official, etc)
- [ ] Escribir frontmatter YAML (11 campos obligatorios)
- [ ] Validar: `python scripts/validate_knowledge_metadata.py`
- [ ] Actualizar INDEX_MASTER.md (% completado)
- [ ] Actualizar NAVIGATION_MAP.md si aplica
- [ ] Commit con mensaje descriptivo
- [ ] Crear ADR si es decisión importante

---

## 🎯 Métricas de Éxito

**Fase 5 (después):**
- +4 exchanges documentados (BingX, Bybit, y otros)
- +2 wallets documentadas (Ledger, MetaMask)
- +3 altcoins/casos documentados
- +10-15 documentos nuevos
- Sistema cubre 95% de casos de usuarios

**Fase 6+ (visión):**
- 20+ exchanges documentados
- 10+ wallets documentadas
- Fiscalidad de 5+ países
- 150+ documentos totales
- Sistema es referencia estándar para auditoría de cripto

---

## 🚪 Cómo Continuar

### Paso 1: Seleccionar Siguiente Documento

Recomendación: **BingX** (usado en `agp`, no documentado)

```bash
# Crear archivo
touch knowledge/cointracking/behavioral/BINGX_MECHANICS.md

# Copiar plantilla de KB-B2-010
# Llenar con detalles específicos de BingX
# Validar metadatos
# Commit
```

### Paso 2: Validar

```bash
python scripts/validate_knowledge_metadata.py
# Debe mostrar +1 en B2 total
```

### Paso 3: Actualizar Índices

- INDEX_MASTER.md (% B2)
- NAVIGATION_MAP.md (si BingX merece entrada)

### Paso 4: Commit

```bash
git commit -m "Agregar KB-B2-010: BingX Mechanics"
```

---

## 📞 Referencias

- [knowledge/INDEX_MASTER.md](knowledge/INDEX_MASTER.md) — estado actual de cada nivel
- [KNOWLEDGE_MAINTENANCE.md](knowledge/KNOWLEDGE_MAINTENANCE.md) — cómo crear docs
- [ADR-033](adr/ADR-033.md) — arquitectura de 6 niveles
- [GOVERNANCE_WORKFLOW.md](GOVERNANCE_WORKFLOW.md) — cómo registrar decisiones

---

**Documento:** Hoja de Ruta de Cobertura  
**Creado:** 2026-07-05  
**Versión:** 1.0  
**Estado:** Listo para implementar

El sistema base está completo. Esta hoja de ruta amplía la cobertura de forma sistemática.
