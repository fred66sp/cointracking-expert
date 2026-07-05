---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-027: Integración de nuevos exchanges sin perder trazabilidad histórica

**Status:** Accepted

**Date:** 2026-07-04
**Accepted:** 2026-07-05 — protocolo de 4 fases aplicado en la práctica (p. ej. reconstrucción completa de BingX en agp2025, 2026-07-03).

## Context

El proyecto soporta multi-proyecto (ADR-013) y cada proyecto puede auditar múltiples exchanges. Cuando un usuario agrega un nuevo exchange (Kraken, Bybit, Coinbase), el riesgo no está en la importación técnica en sí, sino en cómo esos datos nuevos interactúan con un histórico ya reconciliado.

### El problema real

Imagina que ya tienes una auditoría cerrada de Binance:
- Saldos validados
- FIFO calculado
- Informe listo para el asesor fiscal

Ahora importas Kraken de 2022. ¿Qué pasa?

1. **Transferencias contradictorias:** Ves una "Retirada 5 BTC" en Binance el 1 de marzo, pero no aparece "Depósito 5 BTC" en Kraken. ¿Dónde está el BTC? ¿Se perdió? ¿Fue a una billetera propia no importada?

2. **Duplicados invisibles:** Si importas CSV de Kraken y luego API de Kraken en abril, podrías tener 500 operaciones duplicadas sin saberlo (caso ADR-014: Binance batching, pero a nivel de importación).

3. **FIFO roto:** Si agregaste 10 BTC en Kraken en enero y ahora eso entra en el cálculo FIFO, todas las ventas posteriores pueden tener un coste base distinto.

4. **Histórico en dos versiones:** Una auditoría de Binance diciembre-2025 es válida. Una auditoría de Binance-Kraken conjunta en julio-2026 es completamente diferente. Si el usuario modifica algo en Kraken en agosto, ¿cuál auditoría es la verdadera?

**La experiencia de auditoría demuestra que las inconsistencias no suelen producirse durante la importación en sí, sino por la interacción entre el nuevo histórico y los datos ya reconciliados:** APIs con limitación temporal, importaciones parcialmente solapadas, cambios de nomenclatura de activos, operaciones ya existentes reimportadas, y comisiones desiguales entre fuentes.

**ADR-013** dice que cada proyecto aísla datos. Pero **dentro** de un proyecto, agregar una fuente nueva es operación de alto riesgo si no hay protocolo.

## Decision

**La integración de un nuevo exchange no es una importación, es un evento de reconciliación.** El protocolo consta de cuatro fases obligatorias.

### Fase 1 — Preintegración

Antes de importar cualquier dato, el agente debe verificar el estado actual y **esperar confirmación explícita del usuario** (consentimiento informado, ADR-009):

- Listar la auditoría actual: qué exchanges hay, saldos finales, Missing Purchase History, balances negativos
- Anunciar el nuevo exchange, su rango temporal esperado y motivo de la importación
- Documentar exactamente por qué se agrega (usuario quiere cambiar a Kraken, agregó Bybit como fuente secundaria, etc.)

**Si ya hay problemas sin resolver, PARAR.** No agrandes el caos.

Si alguno de estos elementos no puede verificarse, la integración queda marcada como **[PENDIENTE]**.

### Fase 2 — Importación controlada

La importación debe realizarse minimizando el riesgo de alterar el histórico previamente reconciliado:

- Si usas CSV **y** API, importa una a la vez, valida, luego la otra (evita solapamientos desconocidos)
- Identifica posibles limitaciones temporales de la API
- Documenta exactamente qué información se incorporó (rango de fechas, número de operaciones, fuente)
- Conserva el origen de cada conjunto de datos

**El agente únicamente podrá recomendar la importación.** Nunca asumirá que una importación ha sido correcta únicamente porque CoinTracking no muestre errores técnicos.

### Fase 3 — Validación posterior (auditoría completa)

**Tienes que re-auditar el proyecto entero**, no solo el nuevo exchange. Como mínimo, revisa:

✔ Balances negativos (nuevos o agravados)
✔ Missing Purchase History (transferencias sin origen)
✔ Warnings de CoinTracking
✔ Holdings finales (activos nuevos)
✔ Transferencias entre exchanges (retirada en A debe casar con depósito en B)
✔ Operaciones potencialmente duplicadas
✔ Cambios en el coste FIFO (es esperado, pero documéntalo)
✔ Diferencias frente a la auditoría anterior
✔ Posibles cambios en la nomenclatura de activos o comisiones desiguales

Si algo no cierra, marca como **[PENDIENTE DE VERIFICAR]** y pide al usuario que revise en el exchange real.

Las diferencias sin explicación documentada impedirán considerar finalizada la integración.

### Fase 4 — Documentación

Escribe en REGISTRO-CAMBIOS con estos campos obligatorios:

```
Integración: Kraken (2026-07-04)
- Exchange: Kraken
- Fecha de integración: 2026-07-04
- Fuente: API de Kraken, rango 2022-01-01 a 2026-07-04
- Operaciones importadas: ~2500
- Método: API (evitó CSV duplicado)
- Transferencias detectadas: Binance→Kraken 3 operaciones, todas casadas ✓
- FIFO antes: 15.230€ | después: 15.450€ (impacto: +220€ por nuevas adquisiciones Kraken 2022)
- Incidencias: falta 1 transferencia enero-2023 (usuario confirmó: fue a billetera privada no importada)
- Estado final: validado ✓
```

La documentación debe permitir reproducir posteriormente el proceso completo.

## Consequences

**Positive:**

- **Trazabilidad histórica completa:** Mantiene el registro exacto de cuándo y cómo entró cada fuente de datos
- **Reduce riesgos de duplicados:** Al importar CSV y API por separado, identificas solapamientos antes de que pasen desapercibidos
- **Evita sorpresas fiscales:** Si el FIFO cambió 5000€ inesperadamente, lo ves ahora, no en la declaración a Hacienda
- **Protección contra modificaciones no autorizadas:** Las auditorías previas no quedan invalidadas sin revisión formal
- **Separación clara:** observación técnica ≠ decisión contable. El usuario sabe qué verificar y cuándo
- **Alineado con ADR-009:** consentimiento informado explícito antes de actuar
- **Reproducibilidad:** al año siguiente, sabes exactamente qué sucedió y cómo recuperarlo si es necesario

**Negative:**

- **Más lento:** Re-auditar el proyecto completo toma 30-60 minutos si el proyecto es grande
- **Puede revelar problemas históricos:** A veces descubres que tu auditoría anterior (p. ej. solo Binance) tenía huecos de trazabilidad
- **Requiere rigor manual:** No es automático; exige verificar transferencias contra el exchange real
- **Ambigüedad residual:** En la frontera entre "esperado" y "pendiente de verificar" siempre hay grises
- **Responsabilidad compartida:** El usuario debe verificar a veces, no solo confiar

## Ejemplos

### Ejemplo 1 — Incorporación de Kraken (transferencias casadas)

Situación inicial:
- Binance completamente reconciliado (5 años, 500 operaciones).

Nueva acción:
- Usuario importa Kraken mediante API, rango 2022-01-01 a presente.

El agente deberá comprobar:
- Si existen transferencias Binance → Kraken (retiradas que deberían tener depósito espejo)
- Si todas las retiradas de Binance coinciden con depósitos en Kraken
- Si aparecen activos nuevos en Kraken que no estaban en Binance
- Si surge Missing Purchase History (FLOKI o BTC sin origen rastreable)

**Resultado esperado:** Diferencias documentadas. Si falta una retirada/depósito, se marca como [PENDIENTE DE VERIFICAR].

---

### Ejemplo 2 — Incorporación de Bybit (duplicados entre CSV y API)

Situación inicial:
- Binance reconciliada, usuario decidió agregar Bybit.

Nueva acción:
- Primero importa CSV de Bybit (Spot trading, julio-diciembre 2025).
- Luego conecta API de Bybit (intención: obtener datos completos 2023-presente).

Riesgos:
- La API podría repetir las operaciones del CSV.
- Las comisiones podrían diferir (CoinTracking API vs. CSV crudo).
- Órdenes parcialmente rellenadas podrían aparecer duplicadas.

El agente deberá:
- Importar CSV primero, validar, documentar número de operaciones.
- Luego importar API, comprobar solapamiento.
- Detener cualquier propuesta de eliminación automática hasta verificar origen de ambas fuentes.

**Alineación:** ADR-009 (máxima cautela) + ADR-014 (validar duplicados, no asumir que el patrón prueba duplicidad).

---

### Ejemplo 3 — Incorporación de Coinbase (trazabilidad de depósitos externos)

Situación inicial:
- Binance reconciliada desde 2020, usuario importó primero desde banco.

Nueva acción:
- Usuario agrega Coinbase por primera vez (2024-presente, 50 operaciones).

Durante la validación posterior aparecen:
- Depósitos en Coinbase sin origen conocido (7 depósitos de €1000 cada uno, sin transferencia desde Binance ni desde el banco).

El agente deberá:
- Buscar transferencias desde Binance u otros exchanges ya importados.
- Comprobar si existen wallets externas conectadas a Coinbase (importar esos datos si aplica).
- Marcar como **[PENDIENTE]** cualquier movimiento sin trazabilidad documentada.

**No deberá inventar el origen de los fondos.** Si el usuario lo sabe (p. ej. "fue de una herencia"), queda documentado en el comentario de la operación.

---

## Notes

### Relación con ADRs existentes

- **ADR-004:** Reconciliación basada en datos reales antes de cerrar specs
- **ADR-009:** No asumir que una importación es correcta porque CoinTracking no grite; consentimiento informado antes de actuar
- **ADR-010:** La gestión de caché no debe reutilizar información que pueda haber quedado invalidada tras una nueva importación
- **ADR-011:** Toda modificación debe quedar registrada en REGISTRO-CAMBIOS con trazabilidad completa
- **ADR-013:** Multi-proyecto aísla datos, pero dentro de un proyecto hay un único histórico que debe ser consistente
- **ADR-014:** Validar duplicados antes de actuar (aplica aquí: importar 2x el mismo exchange sin darte cuenta)
- **ADR-020:** Datos obtenidos mediante MCP poseen vigencia limitada; pueden requerir nueva validación tras cambios en CoinTracking

### Por qué esto es crítico

**Caso hipotético:** Usuario audita solo Binance, declara ganancias. Un año después, importa Kraken y descubre que tuvo 50 BTC en Kraken desde 2021. Eso cambia TODA la base FIFO y las ganancias que declaró. Asesor fiscal rechaza la declaración y hay que re-hacer todo.

**Con este ADR:** El usuario sabe que agregar Kraken es "operación de reconciliación", no "importación técnica". Ejecuta el protocolo de 4 fases desde el principio y evita sorpresas en Hacienda.

### Pendientes abiertos

- **[PENDIENTE]** Definir un identificador único de "versión de auditoría" tras cada integración, para poder rastrear qué cambió
- **[PENDIENTE]** Automatizar la detección de solapamientos entre CSV y API (overlap en fechas, operaciones, comisiones)
- **[PENDIENTE]** Definir criterios objetivos para aceptar automáticamente transferencias entre exchanges cuando existan coincidencias de importe, activo y ventana temporal
- **[PENDIENTE]** Evaluar integración de hashes de reconciliación para detectar modificaciones posteriores del histórico
- **[PENDIENTE]** Establecer un protocolo de rollback coordinado con el futuro ADR sobre versionado y reapertura de auditorías
- **[PENDIENTE]** Cómo manejar el rollback si descubres que un exchange fue un error agregarlo
- **[PENDIENTE]** Automatizar el bloqueo de cambios en CoinTracking durante una integración en curso (para evitar condiciones de carrera)
