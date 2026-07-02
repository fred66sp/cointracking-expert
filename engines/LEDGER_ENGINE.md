# Especificación del motor de libro mayor

**Reconstrucción de balance cronológico**

El motor de libro mayor reconstruye el balance de cada activo a lo largo del tiempo procesando transacciones cronológicamente. Detecta balances negativos (estados imposibles) y valida consistencia.

## Propósito

Reconstruir el historial completo de balance para cada activo procesando todas las transacciones en orden cronológico, y detectar estados imposibles (balances negativos sin datos de fuente suficientes).

## Entradas

- Dataset de transacciones normalizado y validado
- Balances iniciales (si los hay)

## Salidas

- Historial completo de balance para cada activo
- Detección y reporte de balance negativo
- Estado de verificación de balance

## Responsabilidades

1. Ordenar transacciones cronológicamente por timestamp
2. Procesar transacciones en orden, actualizando balances corrientes
3. Detectar estados de balance negativo
4. Reportar historial de transacciones faltante cuando se detecte
5. Verificar balances finales contra datos de referencia conocidos (si están disponibles)

## Algoritmos clave

- Reconstrucción de libro mayor cronológico
- Tracking de estado de balance por activo y cuenta
- Detección y validación de balance negativo

## Casos extremos

- Transacciones con timestamps idénticos
- Conversiones de zona horaria
- Splits y consolidaciones de activos
- Transferencias de puente y activos envueltos
- Cálculos de comisión afectando balances
