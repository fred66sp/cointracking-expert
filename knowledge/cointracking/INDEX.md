# Base de conocimiento de CoinTracking

**Documentación de referencia para la plataforma CoinTracking**

Esta sección contiene documentación exhaustiva sobre la plataforma CoinTracking, sus características, estructuras de datos, mecanismos de importación, formatos de exportación y limitaciones conocidas. Este conocimiento es esencial para entender cómo validar bases de datos de CoinTracking y reconciliar sus datos con fuentes externas.

## Documentos

- **[CSV_FORMAT.md](CSV_FORMAT.md)** — Formato de la exportación "Trade Table", validado contra datos reales. Referencia autoritativa para la capa de importación (columnas, tipos, fechas, comisiones, emparejamiento de transferencias, colisión de tickers, duplicados).
- **[COST_BASIS_AND_VALIDATION.md](COST_BASIS_AND_VALIDATION.md)** — Cómo CoinTracking calcula la base de coste ("purchase pool") y detecta inconsistencias. Conocimiento destilado de fuentes oficiales, base para la auditoría (`tools/ct_audit.py`, skill `audit-cointracking`) y la preparación fiscal (skill `spanish-tax-return`).
- **[WEB_APP_GUIDE.md](WEB_APP_GUIDE.md)** — Guía operativa de la web de CoinTracking: cómo **corregir** los problemas que detecta la auditoría (cambiar tipos, editar valores, transferencias, duplicados, edición masiva) y cómo **generar el informe fiscal de España (FIFO)**. Orientada a guiar al usuario paso a paso, citando el artículo oficial de cada acción.
- **[reference/CATALOG.md](reference/CATALOG.md)** — Índice de los 205 artículos oficiales del centro de ayuda de CoinTracking (título, URL pública, categoría, relevancia). Fuente para destilar más conocimiento propio.
- **[DOCUMENT_CHECKLIST.md](DOCUMENT_CHECKLIST.md)** — qué pedir al usuario más allá del MCP y el CSV Trade Table: informes propios de CoinTracking (Missing Transactions, Double-Entry List, Realized/Unrealized Gains, Balance by Exchange…) para validación cruzada, e información del exchange original (trade_id, extractos, hash on-chain) para los casos que ni CoinTracking sabe certificar.

## Mantenimiento y vigencia (ADR-008)

CoinTracking evoluciona. **Revisar periódicamente** y actualizar la "Vigencia" de cada documento cuando cambien:

- **Formato del CSV export** (columnas, tipos de transacción, formato de fecha) — `CSV_FORMAT.md`.
- **Tickers y sufijos de colisión** (nuevos activos, p. ej. `SOL2`, `WLD3`) — `CSV_FORMAT.md` §8.
- **Herramientas y parámetros del MCP / API** y límites de tasa — `MCP_API.md`.
- **Comportamiento de la plataforma** (purchase pool, avisos, métodos de coste) — `COST_BASIS_AND_VALIDATION.md`.

Fuente autorizada: el **centro de ayuda oficial** (URLs en `reference/CATALOG.md`) y, sobre el formato real, **los propios datos del usuario** (CSV/MCP mandan). Ante la duda, el agente contrasta con los datos reales antes de asumir.

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
