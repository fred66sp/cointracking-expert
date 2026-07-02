# Especificación del motor de transferencias

**Emparejamiento de depósito y retiro**

El motor de transferencias empareja transacciones de retiro de una cuenta con transacciones de depósito a otra cuenta, identificando transferencias huérfanas y pares desemparejados.

## Propósito

Emparejar transferencias entre cuentas, detectando casos donde un lado de una transferencia falta u es huérfano, y validando consistencia de transferencia.

## Entradas

- Todas las transacciones de retiro y depósito
- Historial de transacciones multi-cuenta
- Mappings de dirección de exchange y billetera

## Salidas

- Pares de transferencia emparejados
- Transferencias huérfanas (desemparejadas)
- Timeline de transferencia y visualización de flujo

## Responsabilidades

1. Identificar matches potenciales de transferencia (mismo activo, cantidad similar, fechas cercanas)
2. Emparejar retiros a depósitos
3. Detectar transferencias desemparejadas u huérfanas
4. Manejar transferencias de puente y tokens envueltos
5. Reportar sobre timing de transferencia e impacto de comisión

## Algoritmos clave

- Algoritmo de emparejamiento de transferencia (cantidad, activo, timestamp)
- Resolución de ambigüedad para múltiples candidatos
- Detección de puente y token envuelto

## Casos extremos

- Comisiones reduciendo cantidad de transferencia
- Rellenos parciales o transferencias divididas
- Demoras entre retiro y depósito
- Variaciones de formato de dirección
- Transferencias de puente (convirtiendo entre cadenas)
- Activos con múltiples versiones (envuelto, puenteado, etc.)
