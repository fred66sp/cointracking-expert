# Glosario

**Terminología y definiciones para CoinTracking Expert**

Este glosario define términos clave y conceptos usados a lo largo del proyecto CoinTracking Expert (un **agente de IA auditor** en Claude Code, no un framework de motores; ver `DECISIONS.md#ADR-006`). Entender estos términos es esencial para trabajar con el agente y su base de conocimiento (`knowledge/`).

## Auditoría

El proceso completo de validar una base de datos de CoinTracking para completitud, consistencia y cumplimiento. Una auditoría examina transacciones, reconstruye balances, detecta problemas y produce un informe detallado.

## Balance

La cantidad de un activo específico mantenido en una cuenta, billetera o exchange en un punto dado en el tiempo. Se reconstruye desde el historial de transacciones.

## Transacción duplicada

Una transacción que aparece más de una vez en el dataset, ya sea como un duplicado exacto o un match probabilístico.

## FIFO (First-In-First-Out / Primero-en-entrar-primero-en-salir)

Método de contabilidad que asigna lotes de adquisición a tenencias basado en orden cronológico. Los primeros activos comprados son los primeros vendidos.

## Tenencia

La cantidad de un activo de criptomoneda específico mantenido en un punto específico en el tiempo. Se reconstruye desde el historial de transacciones.

## Libro mayor

Registro completo de todas las transacciones para una cuenta, organizado cronológicamente. Se usa para reconstruir balances y validar consistencia.

## Historial de compras faltante

Situación donde un activo muestra un balance negativo en algún punto, indicando que transacciones estaban faltando del dataset.

## Normalización

Proceso de convertir datos de transacciones de varias fuentes (CSV, API, manual) a una representación canónica.

## Reconciliación

Proceso de emparejar transacciones entre dos fuentes (ej. registros de exchange vs base de datos de CoinTracking) o verificar consistencia dentro de una fuente única.

## Transferencia

Movimiento de activos entre dos cuentas, billeteras o exchanges. Incluye depósitos y retiros.

## Validación

Proceso de verificar datos para consistencia, completitud y cumplimiento con reglas definidas.

## CoinTracking

Plataforma tercera de contabilidad de criptomonedas y tracking de portfolio. Este agente audita y valida los datos de una cuenta de CoinTracking.
