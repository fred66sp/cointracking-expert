# Especificación del motor de tenencias

**Reconstrucción de tenencias actuales e históricas**

El motor de tenencias reconstruye tenencias actuales e históricas desde el historial de transacciones y las compara con tenencias reportadas por CoinTracking, detectando discrepancias.

## Propósito

Reconstruir tenencias actuales desde el historial de transacciones y compararlas con tenencias esperadas, detectando transacciones no reportadas o inconsistencias de datos.

## Entradas

- Libro mayor de transacciones completo
- Tenencias actuales como reportadas por CoinTracking
- Snapshots de tenencias históricas (opcional)

## Salidas

- Tenencias reconstruidas
- Detección y reporte de discrepancias
- Timeline de tenencias

## Responsabilidades

1. Calcular tenencias actuales desde historial de transacciones completo
2. Comparar con tenencias reportadas
3. Detectar y explicar discrepancias
4. Construir timeline de tenencias históricas
5. Reportar cualquier activo con tenencias cero o transacciones solo-transferencia

## Algoritmos clave

- Cálculo de tenencias desde libro mayor
- Detección de discrepancia y atribución
- Construcción de timeline histórico

## Casos extremos

- Cantidades de polvo (tenencias muy pequeñas)
- Recompensas de staking y tokens bloqueados
- Tokens envueltos y puentes
- Activos con múltiples instancias (ej. diferentes versiones DEX)
- Transacciones recientes afectando tenencias actuales
