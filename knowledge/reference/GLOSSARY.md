---
id: KB-E1-001
title: "Glosario: Terminología y Definiciones de CoinTracking Expert"
level: E
domain: cointracking
source: "Corpus documental + ADR-001 (convenciones)"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: null
confidence: high
version: 1.0

related_adr:
  - ADR-001
  - ADR-003
  - ADR-025

related_docs:
  - knowledge/INDEX_MASTER.md
  - knowledge/cointracking/CT_LIST_FORMATS.md

tags:
  - glossary
  - reference
  - terminology
  - español

notes: "Glosario completo de términos técnicos y operativos del proyecto. Referencia viva, actualizar con nuevos términos."
---




# Glosario

**Terminología y definiciones para CoinTracking Expert**

Este glosario define términos clave y conceptos usados a lo largo del proyecto CoinTracking Expert (un **agente de IA auditor** en Claude Code, no un framework de motores; ver `DECISIONS.md#ADR-006`). Entender estos términos es esencial para trabajar con el agente y su base de conocimiento (`knowledge/`).

## Activo

Criptomoneda, token o moneda fiduciaria registrada en CoinTracking. También llamado **Coin**.

## Airdrop

Distribución gratuita de un activo digital que normalmente se registra como un ingreso.

## API

Interfaz que permite a CoinTracking importar automáticamente transacciones desde un exchange o servicio (**importación automática**). Generalmente proporciona comodidad, aunque puede tener limitaciones históricas o de permisos.

## Auditoría

El proceso completo de validar una base de datos de CoinTracking para completitud, consistencia y cumplimiento. Una auditoría examina transacciones, reconstruye balances, detecta problemas y produce un informe detallado.

## Balance

La cantidad de un activo específico mantenido en una cuenta, billetera o exchange en un punto dado en el tiempo. Se reconstruye desde el historial de transacciones (**reconstrucción de balances**).

## Balance negativo

Situación donde el saldo reconstruido de una moneda es inferior a cero en algún momento del historial. Habitualmente indica transacciones faltantes, duplicadas o mal clasificadas.

## Base de datos

Conjunto completo de transacciones almacenadas en CoinTracking.

## Binance Convert

Servicio de Binance que permite intercambiar activos sin utilizar el mercado Spot tradicional. Es una fuente habitual de problemas de importación.

## Binance Earn

Conjunto de productos de inversión de Binance (Simple Earn, Flexible, Locked, etc.) que generan recompensas o intereses.

## Cartera

Nombre en español de **Wallet**: dirección o cuenta donde se guardan activos digitales.

## Comentario

Campo opcional de una transacción utilizado para añadir información adicional.

## Comisión

Coste asociado a una operación, transferencia o retirada. También llamado **Fee**.

## Coste de adquisición

Valor utilizado para calcular el beneficio o pérdida de una venta. También llamado **Cost Basis**.

## CSV

Archivo utilizado para importar o exportar transacciones de CoinTracking o de un exchange.

## Exchange

Plataforma donde se compran, venden o intercambian activos digitales. En algunos campos de CoinTracking aparece como **Intercambio**.

## Exportación

Proceso de obtener los datos de CoinTracking en formato CSV, Excel, PDF u otros formatos.

## FIFO (First-In-First-Out / Primero-en-entrar-primero-en-salir)

Método de contabilidad que asigna lotes de adquisición a tenencias basado en orden cronológico. Los primeros activos comprados son los primeros vendidos.

## Futures

Operaciones con contratos derivados que normalmente requieren un tratamiento específico en CoinTracking.

## Grupo

Campo opcional que permite clasificar operaciones dentro de CoinTracking.

## Importación

Proceso mediante el cual se añaden nuevas transacciones a CoinTracking, ya sea de forma **automática** (vía API) o **manual** (introducidas directamente desde la interfaz).

## Informe fiscal

Documento generado por CoinTracking que calcula beneficios, pérdidas y otros datos fiscales.

## Ingreso

Transacción que incrementa el patrimonio sin tratarse de una compra tradicional (p. ej. airdrop, recompensa).

## Integridad del historial

Propiedad por la que todas las entradas y salidas necesarias para reconstruir correctamente los balances están presentes en la base de datos.

## Libro mayor

Registro completo de todas las transacciones para una cuenta, organizado cronológicamente (**Ledger**). Se usa para reconstruir balances y validar consistencia.

## Lote

Conjunto de monedas adquirido en una fecha concreta que puede utilizarse posteriormente para el cálculo FIFO o LIFO.

## Moneda fiduciaria

Divisa tradicional como EUR, USD o GBP.

## Normalización

Proceso de convertir datos de transacciones de varias fuentes (CSV, API, manual) a una representación canónica.

## Operación

Compra, venta o intercambio entre activos. También llamada **Trade**.

## Precio de adquisición

Valor pagado originalmente por un activo.

## Precio de mercado

Valor actual estimado de un activo.

## Reconciliación

Proceso de emparejar transacciones entre dos fuentes (ej. registros de exchange vs base de datos de CoinTracking) o verificar consistencia dentro de una fuente única. También llamada **Conciliación**: verificar que los holdings reales del usuario coinciden con los holdings reconstruidos por CoinTracking.

## Recompensa

Ingreso procedente de staking, cashback, promociones o programas similares. También llamada **Rewards**.

## Retirada

Salida de activos desde un exchange o cartera.

## Spot

Mercado tradicional donde se compran y venden activos inmediatamente.

## Staking

Bloqueo temporal de activos para obtener recompensas.

## Tenencia

La cantidad de un activo de criptomoneda específico mantenido en un punto específico en el tiempo (**Holding**). Se reconstruye desde el historial de transacciones.

## Token

Activo digital emitido normalmente sobre una blockchain existente.

## Transacción duplicada

Una transacción que aparece más de una vez en el dataset. Puede ser un **duplicado exacto** (idéntica en todos sus campos) o un **duplicado probable** (muy similar, probablemente la misma operación registrada dos veces). Ver ADR-014 para el caso conocido de falsos positivos por *batching* de Binance en el mismo segundo.

## Transferencia

Movimiento de activos entre dos cuentas, billeteras o exchanges. Incluye depósitos y retiradas.

## Validación

Proceso de verificar datos para consistencia, completitud y cumplimiento con reglas definidas.

## Warning

Advertencia generada por CoinTracking indicando una posible inconsistencia.

## Fuente de verdad

Conjunto de datos considerado como referencia principal para una auditoría (CoinTracking, CSV original, API del exchange, extracto oficial, etc.). Ver [[reconcile_against_real_exchange]]: los datos de CoinTracking nunca se dan por buenos solo por ser internamente consistentes.

## Historial de compras faltante

Situación donde un activo muestra un balance negativo en algún punto, indicando que transacciones estaban faltando del dataset. CoinTracking lo señala como *Missing Purchase History*.

## Workflow de auditoría

Secuencia recomendada para revisar una cuenta:

1. Comprobar balances negativos.
2. Revisar Missing Purchase History.
3. Detectar duplicados.
4. Revisar transfers.
5. Validar holdings.
6. Revisar warnings.
7. Generar informe fiscal.

## CoinTracking

Plataforma tercera de contabilidad de criptomonedas y tracking de portfolio. Este agente audita y valida los datos de una cuenta de CoinTracking.

---

# Formatos y modos del auditor

Términos propios de este proyecto: formatos de presentación (`knowledge/cointracking/CT_LIST_FORMATS.md`, `WEB_APP_GUIDE.md` §4bis) y modos de trabajo del agente.

## CT-Task

Formato estándar (bloque-resumen) para indicar al usuario cómo introducir manualmente una transacción en CoinTracking, siguiendo exactamente el orden del formulario. Se usa siempre al cierre de una guía de alta o corrección manual (ADR-024).

## CT-Timeline

Formato cronológico utilizado para visualizar el historial de operaciones en la conversación (ADR-025).

## CT-Audit

Formato utilizado para mostrar hallazgos de auditoría con marcas ✓/⚠/✗ (ADR-025).

## CT-Balance

Formato utilizado para reconstruir y visualizar balances por activo (ADR-025).

## CT-Exchange

Formato utilizado para agrupar operaciones por exchange o wallet (ADR-025).

## CT-Asset

Formato utilizado para agrupar operaciones por criptomoneda (ADR-025).

## CT-Flow

Formato utilizado para representar el recorrido de fondos entre wallets y exchanges (ADR-025).

## Modo Diagnóstico

Modo de trabajo donde el asistente recopila primero toda la información necesaria antes de proponer soluciones.

## Modo Auditoría

Modo de revisión completa de una cuenta CoinTracking para detectar inconsistencias técnicas y fiscales.

## Riesgo Bajo

No se detectan problemas relevantes que afecten a balances o informes fiscales.

## Riesgo Medio

Existen inconsistencias que pueden afectar parcialmente a balances o cálculos fiscales.

## Riesgo Alto

Se detectan errores importantes (balances negativos, compras faltantes, duplicados, APIs incompletas, etc.) que invalidan los informes hasta ser corregidos.
