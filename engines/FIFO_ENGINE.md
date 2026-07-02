# Especificación del motor FIFO

**Reconstrucción de lote de adquisición y cálculo de base de costo**

El motor FIFO reconstruye lotes de adquisición usando contabilidad de primero-en-entrar-primero-en-salir, asignando transacciones de compra a transacciones de venta basado en orden cronológico. Identifica historial de compras faltante y calcula base de costo para propósitos fiscales.

## Propósito

Asignar lotes de adquisición a tenencias usando método FIFO, calcular base de costo, detectar historial de compras faltante y validar cálculos fiscales basados en costo de lote.

## Entradas

- Libro mayor completo con todas las transacciones de compra y venta
- Datos de precio históricos (opcional, para validación)
- Configuración fiscal (método de base de costo)

## Salidas

- Asignación de lote de adquisición para todas las tenencias
- Cálculos de base de costo
- Detección de historial de compras faltante
- Cálculos de ganancia/pérdida fiscal por lote

## Responsabilidades

1. Reconstruir lotes de adquisición cronológicamente
2. Emparejar tenencias a transacciones de compra específicas
3. Detectar situaciones donde más activos fueron vendidos que comprados
4. Calcular base de costo para cada tenencia
5. Calcular ganancias/pérdidas realizadas y no realizadas

## Algoritmos clave

- Emparejamiento de lote FIFO
- Tracking de base de costo
- Cálculo de ganancia/pérdida
- Detección de historial de compras faltante

## Casos extremos

- Cantidades fraccionarias
- Consolidaciones y splits de activos
- Múltiples precios de compra
- Adquisición de costo cero (airdrops, recompensas de staking)
- Ventas de tenencias parciales
