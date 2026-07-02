# Especificación del motor de auditoría

**Orquestación completa del proceso de auditoría**

El motor de auditoría es el orquestador central responsable de coordinar el flujo de trabajo de auditoría completo. Administra importaciones, normalización y delega a motores de validación especializados, luego produce reportes de auditoría exhaustivos con hallazgos y recomendaciones.

## Propósito

Orquestar una auditoría completa de una base de datos de CoinTracking ejecutando todos los motores de validación y sintetizando resultados en un reporte de auditoría coherente.

## Entradas

- Archivo de exportación de CoinTracking (CSV o exportación de base de datos)
- Configuración especificando qué motores ejecutar
- Opcional: Datos de referencia de exchanges o billeteras
- Opcional: Resultados de auditoría anterior para comparación

## Salidas

- Reporte de auditoría completo (markdown, HTML, JSON)
- Lista detallada de hallazgos con evidencia
- Recomendaciones para remediación
- Estadísticas y métricas de resumen

## Responsabilidades

1. Importar y normalizar datos de transacciones
2. Coordinar ejecución de todos los motores de validación
3. Sintetizar hallazgos de todos los motores
4. Generar reporte de auditoría en múltiples formatos
5. Proporcionar resumen ejecutivo y lista de hallazgos

## Algoritmos clave

- Orquestación de flujo de trabajo
- Deduplicación de hallazgos y ranking de severidad
- Generación de reporte con opciones de formato

## Casos extremos

- Datasets vacíos
- Transacción única
- Datasets muy grandes (rendimiento)
- Fuentes de datos mixtas (API + CSV + manual)
