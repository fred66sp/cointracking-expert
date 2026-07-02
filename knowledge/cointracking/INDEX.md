# Base de conocimiento de CoinTracking

**Documentación de referencia para la plataforma CoinTracking**

Esta sección contiene documentación exhaustiva sobre la plataforma CoinTracking, sus características, estructuras de datos, mecanismos de importación, formatos de exportación y limitaciones conocidas. Este conocimiento es esencial para entender cómo validar bases de datos de CoinTracking y reconciliar sus datos con fuentes externas.

## Documentos

- **[CSV_FORMAT.md](CSV_FORMAT.md)** — Formato de la exportación "Trade Table", validado contra datos reales. Referencia autoritativa para la capa de importación (columnas, tipos, fechas, comisiones, emparejamiento de transferencias, colisión de tickers, duplicados).
- **[COST_BASIS_AND_VALIDATION.md](COST_BASIS_AND_VALIDATION.md)** — Cómo CoinTracking calcula la base de coste ("purchase pool") y detecta inconsistencias. Conocimiento destilado de fuentes oficiales, base para los motores fiscal, FIFO, de transferencias y reconciliación.
- **[reference/CATALOG.md](reference/CATALOG.md)** — Índice de los 205 artículos oficiales del centro de ayuda de CoinTracking (título, URL pública, categoría, relevancia). Fuente para destilar más conocimiento propio.

## Contenidos

Este directorio contendrá documentación detallada cubriendo:

- Características y capacidades de CoinTracking
- Mecanismos de importación (CSV, API, manual)
- Formatos y estructura de exportación ✅ (Trade Table — ver CSV_FORMAT.md)
- Tipos de transacciones y manejo ✅ (ver CSV_FORMAT.md §3)
- Reporte de tenencias y balances
- Generación de reportes fiscales
- Limitaciones conocidas y peculiaridades ✅ (ver CSV_FORMAT.md §7-9)
- Documentación de API y ejemplos
- Problemas comunes de datos y resoluciones
