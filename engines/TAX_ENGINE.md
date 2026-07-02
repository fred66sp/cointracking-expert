# Especificación del motor fiscal

**Cálculo y validación de obligación fiscal**

El motor fiscal calcula obligaciones fiscales basadas en el historial de transacciones y valida cálculos fiscales de CoinTracking. Soporte especializado para reglas de tributación española.

## Propósito

Calcular ganancias/pérdidas realizadas, calcular obligaciones fiscales, validar cálculos fiscales de CoinTracking y generar reportes fiscales para múltiples jurisdicciones con énfasis en cumplimiento español.

## Entradas

- Libro mayor de transacciones completo con costos
- Asignaciones de lote de adquisición FIFO
- Configuración fiscal (jurisdicción, método, exenciones)
- Datos de precio para cálculo de ganancia/pérdida

## Salidas

- Ganancias y pérdidas realizadas por transacción
- Obligación fiscal por año y clase de activo
- Reporte fiscal en formato específico de jurisdicción
- Resultados de validación comparando con cálculos de CoinTracking

## Responsabilidades

1. Calcular ganancias/pérdidas realizadas para todas las disposiciones
2. Identificar eventos taxables
3. Calcular obligación fiscal por año
4. Soportar múltiples métodos de contabilidad (FIFO, etc.)
5. Generar reportes fiscales específicos de jurisdicción
6. Validar contra cálculos de CoinTracking

## Algoritmos clave

- Cálculo de ganancia/pérdida realizada
- Agregación de obligación fiscal
- Aplicación de reglas específicas de jurisdicción
- Generación de reporte fiscal

## Reglas fiscales españolas

- Tributación de ganancias de capital
- Clasificación de criptomoneda
- Reglas de venta ficticia (si aplica)
- Requisitos de reportaje
- Reglas de deducción

## Casos extremos

- Regalos y transferencias personales
- Recompensas de staking y airdrops
- Hard forks y splits de token
- Tokens envueltos y puenteados
- Transacciones entre jurisdicciones
- Transacciones no de largo's length
