# Nivel C: Casos Específicos y Patrones de Auditoría

**Ubicación:** `knowledge/cases/`

**Característica:** Documentación de **patrones reales auditados** y procedimientos específicos para resolver problemas comunes en auditorías de CoinTracking.

**Autoridad:** `verified` — casos reales del proyecto `agp2025` = máxima confianza, no es teoría

---

## Estructura Jerárquica (Semilla 1.0)

Nivel C se divide en **3 bloques temáticos** organizados por complejidad y relación:

### C1: Operaciones Especiales por Exchange

**Cómo importan y qué peculiaridades tiene cada exchange**

- [**C1-001: Binance Spot Mechanics**](C1_BINANCE_SPOT_MECHANICS.md) (KB-C1-001)
  - Dust → BNB auto-conversión
  - Binance Convert (cambios internos)
  - Swaps (integraciones DeFi)
  - Binance Earn/Staking (no se importa bien)
  - **Caso real:** agp2025 (1.000+ ops, limpio)

- [**C1-002: BingX Copy Trading Losses**](C1_BINGX_COPY_TRADING_LOSSES.md) (KB-C1-002)
  - **Problema crítico:** Copy Trading no se exporta
  - Cómo detectar las pérdidas "Lost"
  - **Tratamiento fiscal:** ¿Se deduce o no? (TBD, requiere asesor)
  - **Caso real:** agp2025 (-694,67 USDT, no deducible hasta verificar)

### C2: Rendimientos y Eventos de Ingresos

**Staking, Earn, Rewards — clasificación y auditoría fiscal**

- [**C2-001: Staking and Rewards**](C2_STAKING_AND_REWARDS.md) (KB-C2-001)
  - Tipos de staking (simple, bloqueado, DeFi)
  - **Clasificación fiscal:** RCM (Rendimiento de Capital Mobiliario)
  - Valor a declarar: EUR a fecha de recepción (no a hoy)
  - Liquid staking (stETH, xSOL) — casos complejos
  - Auditoría paso a paso
  - **Caso típico:** ETH Staking en Binance

### C3: Transferencias y Problemas de Importación

**Cómo empareja transferencias, qué son las "huérfanas" y cómo resolverlas**

- [**C3-001: Orphan Transfers Resolution**](C3_ORPHAN_TRANSFERS_RESOLUTION.md) (KB-C3-001)
  - Qué es una transferencia huérfana
  - 4 causas raíz (delay, importación parcial, no acreditado, migración incompleta)
  - Detección en CoinTracking
  - Resolución: Nivel 1 (Tx Hash) y Nivel 2 (Heurística)
  - Emparejamiento manual y verificación
  - **Caso real:** MiCA Migration (Binance → Coinbase)

---

## Cómo Usar Nivel C

### Para Auditor/Usuario

Si encuentras un patrón específico durante la auditoría:

1. **Identifica el patrón:** "Copy Trading losses", "Staking rewards", "Orphan transfer"
2. **Busca en Nivel C** el documento correspondiente
3. **Sigue el checklist y procedimiento** (paso a paso)
4. **Verifica contra la fuente real** (exchange, blockchain, asesor si aplica)

### Para Documentación

**Cuándo agregar a Nivel C:**
- Encontraste un patrón real en auditoría
- Es suficientemente común (>1 usuario lo tiene)
- Tiene procedimiento/checklist claro y verificado
- Tiene implicación fiscal o de auditoría material

**Formato de documento:**
```
---
id: KB-C[1-3]-[NUM]
title: descriptivo + nombre caso real
level: C
domain: cointracking
source: verificable (usuario, fecha)
authority: verified
last_verified: AAAA-MM-DD
valid_until: AAAA-MM-DD
confidence: high
version: 1.0
tags: [palabras clave]
---

# [Título]
1. ¿Qué es?
2. Cómo aparece/se detecta
3. Auditoría paso a paso
4. Tratamiento (fiscal/técnico)
5. Caso real verificado
6. Checklist completo
```

---

## Estadísticas de Cobertura

| Bloque | Documentos | Estado | Verificación |
|--------|-----------|--------|-----------|
| **C1** | 2 de ∞ | ✅ En progreso | agp2025 |
| **C2** | 1 de ∞ | ✅ En progreso | Binance Earn |
| **C3** | 1 de ∞ | ✅ En progreso | MiCA Migration |
| **TOTAL** | 4 | ✅ Semilla 1.0 | Verified |

---

## Próximas Adiciones (No Solicitadas Aún)

**Exchange Mechanics (C1):**
- C1-003: Bybit/OKX Futures mechanics
- C1-004: Kraken staking
- C1-005: Coinbase Earn/Staking

**Income & Rewards (C2):**
- C2-002: Airdrops y airdrops condicionales
- C2-003: Lending rewards (Celsius, Aave)
- C2-004: Mining income

**Transfer Patterns (C3):**
- C3-002: Bridge operations (Polygon, Arbitrum, Optimism)
- C3-003: Wrapped tokens (wBTC, wETH, xSOL)
- C3-004: Token migration events (USDT→USDC forzoso, etc.)
- C3-005: Smart contract failures (failed TX, revert)

---

## Relación con Otros Niveles

| Relación | Conexión |
|----------|----------|
| **Nivel B** | Nivel C **ejemplifica** los casos reales de Nivel B |
| **Nivel A** | Nivel C **verifica y valida** los principios de Nivel A |
| **ADRs** | Nivel C **operacionaliza** lo decidido en ADRs (validación, vigencia) |
| **Skills** | Nivel C **guía el diagnostico** en `/audit-cointracking` y `/spanish-tax-return` |

---

## Política de Vigencia (ADR-008/ADR-022)

Cada documento declara:
- `last_verified`: Cuándo se verificó por última vez contra caso real
- `valid_until`: Fecha caducidad de validez
- `source`: Dónde verificar si caduca

**Antes de citar un caso en auditoría:**
1. Comprobar que no esté expirado
2. Si expira pronto (< 30 días), reverificar contra dato real
3. Si expira por regulación (fiscal/exchange), buscar web breve

---

**Última actualización:** 2026-07-05  
**Versión:** 1.0 (semilla inicial de 4 documentos)
