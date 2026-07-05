---
id: KB-A1-003
title: "Rendimientos y otras rentas por criptomonedas (IRPF España 2025)"
level: A
domain: taxation
source: "DGT — Consultas V1766-22, 0018-23, V1948-21; LIRPF arts. 25, 27, 33-34"
authority: official
last_verified: 2026-07-02
valid_from: 2025-01-01
valid_until: 2025-12-31
confidence: medium
version: 1.0

related_adr:
  - ADR-032
  - ADR-031
  - ADR-028

related_docs:
  - CAPITAL_GAINS.md
  - INFORMATIVE_OBLIGATIONS.md

tags:
  - taxation
  - capital-income
  - staking
  - airdrops
  - rewards
  - spain
  - 2025

notes: "Ejercicio 2025. Criterios DGT interpretados. Requiere reverificación si consultas DGT se actualizan o para otros ejercicios."
---




# Rendimientos y otras rentas por criptomonedas (IRPF España)

Cubre las rentas de cripto que **no** son ganancias/pérdidas por transmisión (esas están en `CAPITAL_GAINS.md`): staking, lending/intereses, airdrops, recompensas/bonificaciones y minería.

> ⚠️ Base técnica, **no asesoramiento fiscal**. La calificación exacta depende de las circunstancias de cada operación y debe validarla un profesional. Ver disclaimer en `INDEX.md`.

---

## 1. Cuadro resumen

| Tipo de renta | Calificación | Base imponible | Valoración |
|---------------|--------------|----------------|------------|
| **Staking** (validación/creación de bloques) | Rendimiento del capital mobiliario (RCM), art. 25.2 | **Ahorro** | Valor de mercado en EUR al percibir |
| **Lending / intereses** (earn, préstamo de cripto) | RCM (cesión de capitales), art. 25.2 | **Ahorro** `[VERIFICAR consulta específica]` | Valor de mercado en EUR al percibir |
| **Airdrops** (entrega gratuita promocional) | Ganancia patrimonial **no derivada de transmisión** | **General** | Valor de mercado en EUR al recibir |
| **Recompensas / bonificaciones / referidos** | Depende del contexto (ganancia patrimonial o act. económica) | Normalmente **General** `[VERIFICAR caso]` | Valor de mercado en EUR al recibir |
| **Minería** (con medios propios) | Rendimiento de **actividad económica**, art. 27 | **General** | Valor de mercado en EUR al generar |

> 🔑 **Diferencia clave para el usuario:** lo que va a la **base del ahorro** (staking, intereses) tributa a tipos del 19–30 % (2025). Lo que va a la **base general** (airdrops, referidos, minería) tributa al tipo marginal de la escala general (hasta ~45–47 % según la comunidad autónoma). No es lo mismo.

---

## 2. Staking — RCM, base del ahorro (DGT V1766-22)

- La DGT (**consulta vinculante V1766-22**) califica las criptomonedas obtenidas por **validación y creación de bloques** (staking) como **rendimiento del capital mobiliario**, concretamente **cesión de capitales propios a terceros** (**art. 25.2 LIRPF**) — asimilable a un depósito retribuido.
- **Integración:** **base imponible del ahorro** (tramos de `CAPITAL_GAINS.md` §6).
- **Valoración:** valor de mercado en EUR el **día de la percepción** (se admite el tipo de cambio medio del día).
- **Sin retención** (no hay pagador obligado a retener).
- **Sin gastos deducibles** de administración/depósito (las criptos no son valores negociables).

> ⚠️ `[VERIFICAR]` **Staking delegado** (a través de pools o proveedores de servicios): V1766-22 analiza el staking activo/directo; **no se pronuncia** específicamente sobre el delegado. Confirmar con profesional para ese caso.

---

## 3. Lending / intereses — RCM, base del ahorro

Los rendimientos por **prestar** criptomonedas o por productos tipo "earn"/interés se asimilan a **cesión de capitales a terceros** (RCM, art. 25.2) → **base del ahorro**, valorados a mercado en EUR al percibirse. Es el tratamiento coherente con el del staking.

> `[VERIFICAR]` No se cita aquí una consulta DGT específica para lending; el encaje es por analogía con V1766-22 y el concepto de cesión de capitales. Confirmar antes de darlo por cerrado.

---

## 4. Airdrops — ganancia patrimonial (base general) (DGT 0018-23)

- La entrega gratuita de criptomonedas en un **airdrop** (distribución promocional para darlas a conocer) constituye una **ganancia patrimonial en especie NO derivada de una transmisión** (**consulta DGT 0018-23**).
- **Integración:** al no derivar de una transmisión, va a la **base imponible GENERAL** (no a la del ahorro), en el período en que se reciben.
- **Valoración:** valor de mercado en EUR al recibir.

> Ejemplo: recibir 1.000 ARB valorados en 1,20 €/ARB = **1.200 € de ganancia patrimonial** en base general, en el ejercicio de recepción.

---

## 5. Recompensas / bonificaciones / referidos — según contexto (DGT V1948-21)

- La DGT (**V1948-21**) trata los criptoactivos recibidos como **contraprestación de determinadas actividades comerciales en línea** (p. ej. promociones, referidos).
- **Calificación dependiente del contexto:** puede ser **ganancia patrimonial no derivada de transmisión** (base general, como el airdrop) o **rendimiento de actividad económica** si hay ordenación por cuenta propia. No hay una respuesta única.
- **Valoración:** valor de mercado en EUR al recibir.

> ⚠️ `[VERIFICAR caso a caso]` El tipo `Recompensa / Bonificación` de CoinTracking es ambiguo: si son recompensas de staking/earn → ver §2/§3 (base ahorro); si son promociones/referidos → base general. Hay que inspeccionar el origen de cada una.

### 5.1 Desglose por origen concreto (criterio profesional, sin consulta DGT específica salvo lo indicado)

Verificado 2026-07-04 (búsqueda cruzada de doctrina secundaria; ninguna consulta DGT localizada para estos sub-casos salvo donde se indica):

| Origen concreto | Calificación habitual | ¿Consulta DGT específica? |
|---|---|---|
| Airdrop simple (token gratis por poseer wallet/protocolo) | Ganancia patrimonial, art. 33.1/37.1.l | **Sí** — V1948-21, V0648-24 (ver §4) |
| Bonus de bienvenida (registro en exchange) | Ganancia patrimonial, art. 33.1 | No — analogía con airdrop |
| Cashback promocional (no descuento de comisión) | Ganancia patrimonial, art. 33.1 | No — analogía |
| Referido ocasional (usuario particular) | Ganancia patrimonial, art. 33.1 | No — analogía |
| Referido recurrente/organizado ("afiliado profesional") | Podría calificar como **actividad económica** (art. 27) en vez de ganancia patrimonial, si hay habitualidad y ordenación de medios | No — sin doctrina DGT |
| Learn & Earn (contraprestación: ver curso/responder test) | Sin criterio asentado — al existir una actuación del usuario, se discute entre ganancia patrimonial o prestación retribuida | No — sin doctrina DGT |

> `[VERIFICAR]` Esta tabla es **criterio profesional por analogía**, no doctrina DGT confirmada (salvo el airdrop simple). No afirmar al usuario una calificación cerrada para estos sub-casos sin advertir la falta de consulta específica.

---

## 6. Minería — actividad económica (base general)

- La minería con **medios personales y materiales propios** constituye **actividad económica** (rendimientos de actividades económicas, **art. 27 LIRPF**), según la DGT.
- **Ingreso íntegro:** valor de mercado del criptoactivo minado **en el momento de su generación**.
- **Gastos deducibles:** electricidad, amortización de equipos, local, asesoría, etc. (rendimiento neto = ingresos − gastos).
- **Integración:** **base general**.
- **Obligaciones formales:** alta en IAE (epígrafe **831.9**) y, según dimensión, alta como autónomo. En principio no sujeta a IVA (no hay contraprestación directa).

---

## 7. Regla transversal CRÍTICA: valor al percibir = coste de adquisición futuro

El valor de mercado en EUR por el que tributa una renta al **recibirla** (staking, airdrop, recompensa, minería) se convierte en el **valor de adquisición (base de coste)** de esas monedas para cuando **luego se transmitan** (venta o permuta), a efectos del FIFO y de la ganancia/pérdida patrimonial.

> 🔑 Esto **evita la doble imposición**: primero tributas por la renta al recibir; después, al vender, solo tributa la **variación** desde ese valor. El agente debe usar ese valor como coste en `CAPITAL_GAINS.md` §2–4. Si no se registra, aparecerá el problema de "venta sin base de coste" (ver `../../cointracking/COST_BASIS_AND_VALIDATION.md` §3).

---

## 8. Mapeo a los tipos de CoinTracking (ver `CSV_FORMAT.md` §3)

| Tipo en CoinTracking | Tratamiento probable | Base |
|----------------------|----------------------|------|
| `Staking` | RCM (V1766-22) | Ahorro |
| `Ingresos por intereses` | RCM (lending) | Ahorro `[VERIFICAR]` |
| `Recompensa / Bonificación` | Según origen: airdrop/promoción → GP; staking/earn → RCM | General o Ahorro `[VERIFICAR caso]` |
| `Ingresos` (genérico) | Según origen | `[VERIFICAR caso]` |
| (minería, si aplica) | Actividad económica | General |

---

## 9. Cuestiones abiertas

1. **Staking delegado** (pools/proveedores): sin pronunciamiento específico de la DGT (§2).
2. **Lending/intereses:** falta cita de consulta DGT específica; encaje por analogía (§3).
3. **`Recompensa / Bonificación`:** requiere clasificar el origen de cada operación (§5, §8).
4. Confirmar que CoinTracking asigna como **coste de adquisición** el valor a la fecha de recepción de estas rentas (§7); si no, corregirlo para no inflar ganancias futuras.

---

## Fuentes

- DGT, consulta vinculante **V1766-22** (staking → RCM art. 25.2, base del ahorro)
- DGT, consulta **0018-23** (airdrops → ganancia patrimonial no derivada de transmisión, base general)
- DGT, consulta **V1948-21** (criptoactivos como contraprestación de actividades comerciales en línea)
- LIRPF (Ley 35/2006): arts. 25 (RCM), 27 (actividades económicas), 33-34 (ganancias patrimoniales)
- Cuatrecasas, "Tributación del staking en IRPF"; Uría Menéndez, consultas DGT sobre criptoactivos
