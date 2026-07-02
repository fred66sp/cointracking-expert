# Arquitectura del sistema

**Diseño del Framework CoinTracking Expert**

Este documento describe la arquitectura general del framework CoinTracking Expert, incluyendo responsabilidades de componentes, flujo de datos, interfaces y principios de diseño. La arquitectura enfatiza modularidad, reproducibilidad y separación limpia de responsabilidades.

## Descripción general de alto nivel

El framework está organizado como un pipeline de motores independientes que procesan datos de transacciones a través de etapas progresivas de validación y análisis. Cada motor tiene entradas, salidas y límites de responsabilidad bien definidos.

```
Exportación de CoinTracking → Normalización → Motor de auditoría → Motores de validación → Generación de reportes
                                                     ↓
                                           ├─ Motor de duplicados
                                           ├─ Motor de transferencias
                                           ├─ Motor de libro mayor
                                           ├─ Motor de tenencias
                                           ├─ Motor FIFO
                                           └─ Motor de impuestos
```

## Principios principales

1. **Modularidad**: Cada motor es independiente y testeable
2. **Reproducibilidad**: La misma entrada siempre produce la misma salida
3. **Basado en evidencia**: Todas las conclusiones respaldadas por datos de transacciones
4. **Transparencia**: Cada problema incluye causa, impacto y evidencia
5. **Intervención mínima**: Nunca modificar sin justificación

## Estructura de componentes

### Capa de importación y normalización

Responsable de leer datos de varias fuentes (CSV de CoinTracking, API, entrada manual) y normalizarlos a representación canónica. Maneja conversión de formato, limpieza de datos y validación de esquema.

### Motor de auditoría

Orquesta el proceso de auditoría, administra el flujo de trabajo y coordina entre motores especializados. Detecta inconsistencias y produce informes de auditoría.

### Motores especializados

- **Motor de duplicados**: Identifica duplicados exactos y probabilísticos de transacciones
- **Motor de transferencias**: Empareja depósitos y retiros entre cuentas
- **Motor de libro mayor**: Reconstruye balances cronológicamente
- **Motor de tenencias**: Reconstruye tenencias esperadas desde historial de transacciones
- **Motor FIFO**: Calcula lotes de adquisición e historial de compras faltante
- **Motor de impuestos**: Valida cálculos fiscales y genera informes

### Generación de reportes

Produce reportes de auditoría en múltiples formatos (Markdown, HTML, Excel, JSON) con niveles de detalle configurables.

## Estructuras de datos

Todos los componentes principales utilizan estructuras de datos estandarizadas definidas en schemas/. Estas garantizan consistencia en importaciones, motores y exportaciones.

## Objetivos de calidad

- Cero falsos positivos cuando sea razonablemente posible
- Resultados deterministas y reproducibles
- Trazabilidad y registros de auditoría completos
- Salida legible por humanos y por máquinas
