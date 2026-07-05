---
id: KB-B1-017
title: "Cambios regulatorios de exchanges (2026) — Impacto en reconciliación"
level: B
domain: cointracking
source: "Investigación web + auditoría de agp2025 (caso real del usuario)"
authority: verified
last_verified: 2026-07-05
valid_from: 2026-01-01
valid_until: 2026-12-31
confidence: medium
version: 1.0

tags:
  - regulatory
  - operational
  - exchanges
  - verification

notes: "Resumen de cambios 2026 que afectan auditoría de CoinTracking. Usuario actual (agp2025) impactado por MiCA-Binance y conversión USDT→USDC Q1 2025."
---

# Cambios Regulatorios de Exchanges (2026) — Impacto en Reconciliación

**Fecha:** 2026-07-05  
**Alcance:** Cambios confirmados que afectan a la auditoría de datos de CoinTracking  
**Usuario relevante:** agp2025 (Binance, Coinbase, BingX, 2024-2025)

---

## Resumen Ejecutivo

| Exchange | Cambio | Fecha | Impacto | Estado |
|----------|--------|-------|--------|--------|
| **Binance** | Salida de UE (MiCA) | 2026-07-01 | Migración forzosa a otros exchanges | ✅ Verificado: agp2025 migrando a Coinbase |
| **Binance** | Conversión USDT→USDC | Q1 2025 | Hecho imponible: permuta cripto-cripto | ✅ Auditado: sin discrepancias (agp2025) |
| **Coinbase** | Expansión EU (MiCA-ready) | 2024-2026 | Destino principal de migraciones desde Binance | ✅ Operativo |
| **BingX** | Cambios en producto Futures | 2025 | Migración de derivados, pérdidas en Copy Trading | ✅ Documentado: agp2025 cierre exacto (72,94 USDT) |

---

## 1. Binance — Salida de la UE (MiCA, 2026-07)

### Qué pasó

Binance no aseguró licencia MiCA antes del **1 de julio de 2026**. Desde esa fecha:
- ❌ Nuevas órdenes spot: bloqueadas para usuarios EU
- ❌ Depósitos nuevos: bloqueados
- ❌ Productos Earn/Staking: cerrados
- ✅ Retiradas: abiertas (fondos accesibles)

**Fuente:** Comunicado oficial Binance + cobertura prensa (CoinDesk, Euronews, crypto.news).

### Impacto en Reconciliación

**Migración de activos:** Si el usuario retiró fondos de Binance en julio 2026:
- ✅ Es una **transferencia entre cuentas propias** (no tributa)
- ⚠️ Pero **debe estar bien emparejada en CoinTracking** (retirada Binance + depósito en destino)
- 🔍 **Auditar como transferencia huérfana potencial** si los datos muestran inconsistencias (ver `audit-cointracking/SKILL.md`, Paso 1)

**Casos especiales:**
- Si Binance forzó liquidación de algún activo antes de permitir retirada → **permuta cripto-cripto** (tributa, Art. 37.1.h LIRPF)
- Precedente real: conversión USDT→USDC Q1 2025 (ya en agp2025, sin impacto fiscal negativo)

### Qué hacer

En la próxima auditoría (proyecto agp2026 esperado, cuando usuario finalice la migración):

```
1. Buscar período julio 2026 (cuando MiCA entró en vigor)
2. Listar retiradas Binance + depósitos en destino
3. Verificar emparejamiento: 
   - ¿Tx Hash coincide?
   - ¿Importe ≈ retirada − comisión?
   - ¿Tiempo razonable (< 24h)?
4. Si hay conversion Trade antes de retirada → clasificar como hecho imponible
5. Contrastar contra histórico real de Binance (antes de que degrade acceso)
```

**Referencia:** `knowledge/reference/context/BINANCE_EU_MICA_EXIT.md`

---

## 2. Binance — Conversión USDT→USDC (Q1 2025)

### Qué pasó

Binance obligó a convertir USDT (Tether) a USDC (Circle) a principios de 2025, probablemente por cambios regulatorios. Impactó a usuarios con holdings USDT.

### Impacto en Auditoría de agp2025

**Estado:** ✅ **Verificado sin discrepancias**

- **Dato:** agp2025 registra una conversión USDT→USDC el 13 de enero de 2025
- **Clasificación:** Permuta cripto-cripto (tributa por Art. 37.1.h LIRPF)
- **Resultado:** Ganancia nominal ~0€ (mismo valor en EUR)
- **Discrepancia:** Ninguna detectada entre CoinTracking y datos reales

**Conclusión:** Si el usuario tiene agp2026, buscar si hay operaciones similares (conversiones forzosas) y clasificarlas igual.

---

## 3. Coinbase — Expansión en la UE (MiCA-ready)

### Qué pasó

Coinbase obtuvo licencia MiCA y es uno de los principales destinos de migraciones desde Binance (Q2-Q3 2026).

### Impacto en Auditoría

**Estado:** ✅ **Operativo en agp2025**

- **Dato:** agp2025 usa Coinbase desde 2024 (sin MiCA)
- **En 2026:** Será destino de migraciones desde Binance (transacciones nuevas esperadas en agp2026)

**Qué auditar en migraciones Binance→Coinbase:**
- Retiradas Binance → depósitos Coinbase, emparejados
- Saldos verificados contra app real de Coinbase (no solo CoinTracking)
- Sin liquidaciones forzosas (si las hay, clasificarlas como hechos imponibles)

**Referencia:** Ya hay procedimiento establecido en `audit-cointracking/SKILL.md`, Paso 1, punto 3 (transferencias huérfanas).

---

## 4. BingX — Cambios en Derivados (2025)

### Qué pasó

BingX cambió la estructura de productos derivados a lo largo de 2025. Impactó:
- Cierre de posiciones Futuros (Standard, Perpetual)
- Migración de Copy Trading a nueva estructura
- Pérdidas acumuladas en sub-cuenta Copy Trading (no exportadas por BingX)

### Impacto en Auditoría de agp2025

**Estado:** ✅ **Cerrado y cuadrado exacto (2026-07-03)**

**Hallazgos durante auditoría:**

| Línea | Descripción | Importe | Clasificación | Tratamiento |
|------|-------------|---------|---------------|-------------|
| Standard 2024 | Pérdidas Futuros | −3.188,08 USDT | Derivatives Loss | Deducible |
| Standard 2025 | Pérdidas Futuros | −3.910,87 USDT | Derivatives Loss | Deducible |
| Perpetual 2024 | Pérdidas Futuros | −5,92 USDT | Derivatives Loss | Deducible |
| Trial Fund | Ingresos prueba | +10 USDT | Income (no fiscal) | No deducible |
| **Copy Trading** | **Pérdidas no exportadas** | **−694,67 USDT** | **"Lost" (NO deducible)** | **Requiere evidencia** |

**Punto crítico:** La entrada "Lost −694,67 USDT" = pérdidas reales en Copy Trading que BingX no exporta. Marcada como "Lost" (no deducible) a propósito. **No usar como pérdida deducible sin validación del asesor o export oficial de BingX.**

**Qué hacer en futuras auditorías de BingX:**
- Solicitar al usuario export específico de Copy Trading si hay actividad
- Clasificar por separado (hecho imponible real vs. pérdida verificable)
- No asumir que "Lost" = pérdida deducible sin evidencia

**Referencia:** `reports/output/agp2025/REGISTRO-CAMBIOS.md`, sección "CIERRE DEFINITIVO de BingX".

---

## 5. Otros Exchanges — Sin Cambios Material Detectado (2026)

### Kraken, OKX, KuCoin

- ✅ Operativos en MiCA (o planeando)
- ✅ Soportados en CoinTracking
- ⚠️ Sin eventos regulatorios críticos detectados en la sesión

**Verificar si el usuario usa estos exchanges:** Si aparecen en agp2026, buscar qué productos usan y si hay cambios de 2026 sin documentar aquí.

---

## 6. Checklist para Próximas Auditorías (agp2026 y posteriores)

### Antes de Auditar

- [ ] **Buscar web** cambios recientes de cada exchange usado (usar términos: "[exchange] MiCA 2026", "[exchange] regulatory update", "[exchange] product changes")
- [ ] **Actualizar este documento** si aparecen cambios no listados aquí
- [ ] **Revisar secciones anteriores** de la auditoría para contexto (qué activos, qué período)

### Durante la Auditoría

- [ ] **Transferencias Binance→destino:** Verificar emparejamiento (no huérfanas)
- [ ] **Conversiones forzosas:** Buscar operaciones Trade de tipo "exchange directive" o notas del usuario
- [ ] **Derivados BingX:** Separar pérdidas verificables de "Lost" (requiere evidencia)
- [ ] **Saldos finales:** Cotejar contra app real del exchange, no solo CoinTracking
- [ ] **Documentar cambios:** Actualizar `REGISTRO-CAMBIOS.md` con fecha, exchange, impacto

### Después de Auditar

- [ ] **Revisar clasificación fiscal:** ¿conversiones = permutas? ¿transferencias = no imponibles?
- [ ] **Si hay dudas:** Consultar con asesor antes de incluir en declaración

---

## Referencias

| Documento | Propósito |
|-----------|-----------|
| `knowledge/reference/context/BINANCE_EU_MICA_EXIT.md` | Contexto detallado de salida Binance UE |
| `knowledge/exchanges/official/BINANCE.md` | Particularidades técnicas de importación Binance |
| `knowledge/taxation/spain/CAPITAL_GAINS.md` | Tributación de permutas cripto-cripto (Art. 37.1.h LIRPF) |
| `tools/ct_audit.py` | Detección automática de transferencias huérfanas |
| `reports/output/agp2025/REGISTRO-CAMBIOS.md` | Documentación de auditoría y cambios aplicados |

---

**Generado:** 2026-07-05  
**Vigencia:** 2026-12-31 (reverificar cambios regulatorios si es más nuevo)  
**Usuario relevante:** agp2025 (migración Binance→Coinbase en curso), futuros proyectos que usen estos exchanges
