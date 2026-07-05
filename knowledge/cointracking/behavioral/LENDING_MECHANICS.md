---
id: KB-B1-003
title: "Cómo CoinTracking maneja Lending y Yield Farming"
level: B
domain: cointracking
source: "Casos reales CT-011 + análisis"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-07-03
confidence: medium
version: 1.0

related_adr:
  - ADR-003
  - ADR-032

related_docs:
  - knowledge/taxonomy/spain/CAPITAL_INCOME.md
  - CT-011-lending-tratado-como-transferencia-generica.md

tags:
  - cointracking
  - lending
  - defi
  - behavioral

notes: "Operativo: cómo CoinTracking registra operaciones de lending; errores comunes."
---




# Cómo CoinTracking Maneja Lending y Yield Farming

## Definición

**Lending** = Prestar cripto a una plataforma (p. ej. Aave, Compound) para ganar interés.

**Yield Farming** = Proporcionar liquidez a pools y ganar comisiones/tokens de recompensa.

---

## Problema Común: CT-011

**Síntoma:** Lending registrado como "Transferencia genérica" en lugar de operación específica.

**Causa:** CoinTracking no tiene tipo específico para lending (depende del tipo de fuente).

**Impacto:**
- Balance parece "desaparecido" (la transferencia a la plataforma lending)
- Intereses no aparecen como ingresos
- Tax Report confuso

---

## Cómo CoinTracking Registra Lending

### Flujo ideal

```
Día 1: Depositas 10 ETH en Aave
  Tipo: "Deposit" (a Aave)
  Cantidad: 10 ETH
  
Día 30: Ganas 0.05 ETH de interés
  Tipo: "Income" o "Reward"
  Cantidad: 0.05 ETH
  
Día 60: Retiras 10.05 ETH
  Tipo: "Withdrawal"
  Cantidad: 10.05 ETH
```

### Flujo problemático (common)

```
Si importas vía CSV desde Aave:
  La plataforma registra:
    - Envío de ETH a Aave como "Transfer out" (genérico)
    - Intereses como "Income" (si está documentado)
    - Retiro como "Transfer in" (genérico)
  
CoinTracking lo importa como:
    - Transfer: 10 ETH → Aave (confuso, parece pérdida)
    - Income: 0.05 ETH (correcto si está)
    - Transfer: 10.05 ETH ← Aave (confuso, parece ganancia)
```

---

## Validación y Corrección

### Identificar lending mal clasificado

```
CoinTracking → Transacciones:
  Busca "Transfer" a plataformas conocidas:
    - Aave, Compound, Lido, dYdX, etc.
  
Si viste "Transfer out" sin "Transfer in" correspondiente:
    → Probablemente es lending incompleto o mal registrado
```

### Corregir

```
Opción 1: Editar manualmente
  Transfer out: 10 ETH a Aave → Cambiar a "Deposit (Aave)"
  Income: 0.05 ETH ← Aave → Mantener como "Income"
  Transfer in: 10.05 ETH ← Aave → Cambiar a "Withdrawal (Aave)"
  
Opción 2: Reimportar desde Aave directamente
  - Aave tiene una exportación específica
  - CoinTracking puede importarla y clasificar mejor
```

---

## Tratamiento Fiscal (España, IRPF)

**Intereses de lending = Rendimiento del capital**

```
Regla (DGT):
  - Momento exigible: cuando se acredita el interés
  - Valuación: en ETH/precio del activo en ese momento
  - Impacto: se suma a "Rendimientos del capital" en IRPF
  
Depósito del principal:
  - NO es venta (no consumes base de coste)
  - El principal sigue siendo base de coste (para cuando vendas)
```

**Ejemplo:**
```
1 enero: Depositas 10 ETH en Aave (precio 1.500€ = 15.000€ total)
  → Base de coste registrada: 15.000€

1 junio: Retiras 10.05 ETH (precio 2.000€ = 20.100€ total)
  → Interés: 0.05 ETH = 0.05 × 2.000€ = 100€
  → Ganancia por apreciación: (2.000€ - 1.500€) × 10 ETH = 5.000€
  
IRPF 2025:
  - Rendimientos: 100€ (interés)
  - Ganancias patrimoniales: 5.000€ (apreciación)
  - Total: 5.100€
```

---

## Validación en CoinTracking

```
Reports → Gains:
  ¿Aparecen los intereses de lending en "Income/Rendimientos"?
    SÍ → OK
    NO → Verificar que income esté registrado y clasificado correctamente
    
Check: ¿La ganancia por apreciación es correcta?
  Debería ser: (precio retiro - precio depósito) × cantidad
```

---

## Integración

- **ADR-003:** Modelo de transacciones — Lending es depósito de principal + ingresos
- **CAPITAL_INCOME.md:** Intereses son rendimiento del capital
