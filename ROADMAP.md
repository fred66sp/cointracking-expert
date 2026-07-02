# Hoja de ruta de desarrollo

**Cronograma del proyecto CoinTracking Expert**

> **Nota (2026-07-02) — ADR-006:** El **corto plazo** ya no sigue este cronograma. El producto inmediato es un **agente de IA auditor en Claude Code** que se apoya en la base de conocimiento (ver `DECISIONS.md`, ADR-006). Las fases de framework/SDK descritas abajo (motores deterministas, CLI, API) pasan a ser **visión futura/opcional**, no el camino actual. Este documento se conserva como referencia de esa visión a largo plazo.

Este documento describe las fases de desarrollo planificadas para el framework CoinTracking Expert, desde la fundación del proyecto hasta diagnósticos asistidos por IA. La hoja de ruta sigue un enfoque de desarrollo impulsado por documentación donde las especificaciones preceden a la implementación.

## Fase 1: Fundación del proyecto (Actual)

Establecer infraestructura del proyecto, gobernanza y organización de la base de conocimiento. Los entregables incluyen carta de proyecto, documentación de arquitectura, directrices de contribución y estructura de conocimiento inicial. Sin código de implementación en esta fase.

## Fase 2: Desarrollo de la base de conocimiento

Construir conocimiento de dominio integral cubriendo CoinTracking, exchanges, billeteras, blockchains, tributación y patrones de reconciliación. Esta fase puebla los directorios de conocimiento con documentación estructurada y casos de auditoría del mundo real.

**Dependencia crítica (ADR-004):** obtener exportaciones reales de CoinTracking cuanto antes. Las especificaciones sensibles a datos (formato CSV, importación, duplicados, transferencias) se redactan en borrador y se validan contra esos datos reales antes de darse por cerradas. No se cierran specs de datos sobre suposiciones.

## Fase 3: Especificaciones de motores

Documentar especificaciones funcionales completas para todos los motores (Auditoría, Reconciliación, Libro mayor, FIFO, Tenencias, Transferencia, Duplicados, Impuestos, Reporte). Cada especificación incluye entradas, salidas, algoritmos, casos extremos y escenarios de prueba.

## Fase 4: Implementación Python

Implementar librería Python central con todos los motores. Incluye modelos de datos, capa de importación/normalización e implementaciones de motores individuales. Pruebas unitarias e integración exhaustivas.

## Fase 5: Interfaz de línea de comandos

Desarrollar herramienta CLI para ejecutar auditorías, generar reportes y consultar resultados. Incluye gestión de configuración, formato de salida y capacidades de procesamiento por lotes.

## Fase 6: API REST

Crear API RESTful para acceso remoto a motores de auditoría. Incluye autenticación, limitación de velocidad, gestión de trabajos y streaming de reportes.

## Fase 7: Integración de agentes de IA

Integrar modelos de IA (Claude, ChatGPT, otros) como capa de explicabilidad. La IA asiste en explicar hallazgos, generar recomendaciones y diagnósticos interactivos sin reemplazar cálculos deterministas.
