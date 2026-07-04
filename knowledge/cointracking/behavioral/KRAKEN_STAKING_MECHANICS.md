---
id: KB-B2-006
title: "Cómo CoinTracking maneja Kraken Staking y Rewards"
level: B
domain: cointracking
source: "Análisis de casos Kraken + documentación"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: medium
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/cointracking/behavioral/STAKING_MECHANICS.md
  - knowledge/cointracking/behavioral/LENDING_MECHANICS.md

tags:
  - kraken
  - staking
  - rewards
  - behavioral

notes: "Operativo: cómo CoinTracking importa staking de Kraken."
---

# Cómo CoinTracking Maneja Kraken Staking y Rewards

## Diferencia Kraken vs Binance

| Característica | Binance Earn | Kraken Staking |
|---|---|---|
| **Método** | Depósito → Earn | Staking directo |
| **Retención** | Flexible/Locked | Directo en cuenta |
| **Rewards** | Acumulan diariamente | Acumulan diariamente |
| **Importación** | API + Earn historial | API + CSV |
| **Complejidad** | Media | Media-Alta |

---

## Tipos de Staking en Kraken

### Tipo 1: On-Chain Staking (Ethereum 2.0, Solana, etc.)

```
Ejemplo: Staking de ETH en Kraken

¿Qué sucede?
  1. Depositas 32 ETH en Kraken (mínimo)
  2. Kraken lo valida en Ethereum 2.0
  3. Recibes recompensas (APY ~3-4%)
  4. El ETH está bloqueado (no puedes retirarlo hasta 2025)
  
¿Cómo aparece en CoinTracking?

Ideal:
  - Depósito: 32 ETH en staking
  - Income: +0.5 ETH/día (recompensas)
  - Tipo: "Staking" o "Reward"

Realidad:
  - Kraken muestra "32 ETH staked"
  - Las recompensas aparecen como "Income"
  - Pero pueden no estar separadas visualmente
```

### Tipo 2: Kraken Staking (No-Lock Alternative)

```
Introducido por Kraken como alternativa flexible

¿Qué sucede?
  1. Depositas ETH (incluso 0.1 ETH)
  2. Kraken lo proporciona a validadores
  3. Recibes recompensas en tu cuenta
  4. Puedes retirarlo en cualquier momento

¿Cómo aparece en CoinTracking?
  - Transacción: "Staking reward" o "Income"
  - Monto: pequeño (0.001 ETH por día, típicamente)
  - Frecuencia: diaria o semanal (depende de Kraken)
```

---

## Problema: Rewards Incompletas o No Importadas

**Síntoma común:**

```
Tengo 32 ETH en staking en Kraken desde enero
Debería tener ~100 ETH en recompensas (31 ETH × ~4% APY)
Pero CoinTracking solo muestra 20 ETH

¿Por qué?

Causas:
  1. Kraken API no importa rewards automáticamente
  2. CSV de Kraken no incluye rewards (solo saldo)
  3. Rewards aparecen como "income" pero no están conectadas
  4. Importación parcial (solo últimos meses)
```

---

## Importación de Kraken en CoinTracking

### Vía API (Mejor)

```
Pasos:
  1. CoinTracking → Settings → Exchanges → Add Kraken
  2. Autorizar API key (public + private)
  3. Asegurarse que permissions incluyen "Ledger"
  
Resultado:
  - CoinTracking importa todas las transacciones
  - Incluye staking rewards
  - En tiempo real
  
Verificación:
  [ ] ¿Aparecen operaciones de staking?
  [ ] ¿Las rewards aparecen como "Income"?
  [ ] ¿El balance de ETH es correcto?
```

### Vía CSV (Manual)

```
Kraken → Historial de fondos → Descargar CSV

El CSV incluye:
  - Depósitos/Retiros
  - Trades
  - Staking rewards (si los muestra)
  
Problema:
  - No siempre incluye rewards en formato legible
  - Requiere parseo manual
  
Pasos:
  1. Descargar CSV de Kraken
  2. Filtrar por "staking" o "reward"
  3. Importar en CoinTracking
  4. Verificar que aparezcan como "Income"
```

---

## Validación en CoinTracking

### Verificar Rewards Importados

```
CoinTracking → Reports → Income:
  ¿Aparecen "Staking rewards" o "Kraken income"?
  
¿Cuánto debería ser?
  Fórmula: Principal × APY × Días / 365
  
  Ejemplo:
    32 ETH × 3.5% × 365 / 365 = 1.12 ETH/año
    En 6 meses: ~0.56 ETH
  
¿Coincide con realidad?
  [ ] En Kraken → Mi historial de recompensas
  [ ] Comparar total de rewards
  [ ] ¿Son iguales?
      SÍ → OK, todas importadas
      NO → Faltan rewards, importar manualmente
```

### Detectar Rewards Faltantes

```
En Kraken:
  Staking → Historial de recompensas
  Anotar: Total de ETH recibido
  
En CoinTracking:
  Filtrar: "Income" + "Kraken"
  Sumar todos los valores
  
Comparar:
  Kraken total: 10 ETH
  CT total: 4 ETH
  Faltantes: 6 ETH
  
¿Qué hacer?
  Si faltan, importarlas manualmente (ver abajo)
```

---

## Añadir Rewards Manualmente

```
Si falta importación automática:

CoinTracking → Add Transaction:
  Para CADA reward de Kraken (o en lote):
  
  Tipo: Income
  Activo: ETH
  Cantidad: X ETH (la recompensa ese día)
  Precio: Precio de ETH ese día
  Fecha: Fecha exacta que Kraken lo acreditó
  Exchange: Kraken
  Descripción: "Kraken staking reward"
  
Repetir para cada día/semana (tedioso si son muchos)

Alternativa (más rápido):
  CoinTracking → Add Transaction:
    Tipo: Income
    Cantidad: Total de rewards (ej. 6 ETH)
    Precio: Precio promedio del período
    Fecha: Fecha de promedio
    
  Nota: Menos preciso (pero más rápido)
```

---

## Tratamiento Fiscal (España, IRPF)

**Kraken Staking Rewards = Rendimiento del Capital**

```
Regla:
  - Día que Kraken acredita la recompensa
  - Valuación: precio ETH ese día
  - Se suma a "Rendimientos del capital"
  
Ejemplo:
  Enero 2025: Recibes 0.1 ETH (precio 1.500€)
    → +150€ en Rendimientos
    
  Junio 2025: Recibes 0.1 ETH (precio 2.000€)
    → +200€ en Rendimientos
    
  Total IRPF 2025: 150€ + 200€ = 350€
```

---

## Caso Especial: kETH (Staked ETH Token)

**Kraken ofrece "kETH", token que representa ETH stakeado**

```
¿Qué es?
  - 1 kETH = 1 ETH stakeado
  - Puedes vender kETH (aunque esté stakeado)
  - Puedes recibir rewards mientras tienes kETH
  
¿Cómo aparece?

Si convertiste ETH → kETH:
  TX1: Exchange 32 ETH → 32 kETH
  
¿Es una venta?
  Técnicamente no (mismo valor)
  Pero CoinTracking lo puede clasificar como Trade
  
Problema:
  - Cost basis de kETH se hereda de ETH
  - Pero son assets diferentes en CT
```

---

## Validación: Checklist Kraken Staking

```
Antes de auditar:

[ ] ¿Kraken API está conectado?
    SÍ → OK
    NO → Conectar o usar CSV
    
[ ] ¿Todos los rewards importados?
    Comparar: Kraken vs CoinTracking
    
[ ] ¿Los rewards son "Income"?
    Verificar tipo de transacción
    
[ ] ¿El tratamiento fiscal es correcto?
    Cada reward como rendimiento del capital
    
[ ] ¿El balance de ETH es correcto?
    Principal + rewards = balance total
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Staking es Deposit + Income
- **STAKING_MECHANICS.md:** Diferencias Binance vs otros
- **CAPITAL_INCOME.md:** Tratamiento fiscal de rewards
