# Base de conocimiento de exchanges

**Documentación de referencia para exchanges soportados**

Esta sección contiene documentación detallada para cada exchange soportado, incluyendo estructura de API, formatos de exportación, tipos de transacciones, formatos de dirección y discrepancias conocidas entre datos de exchange y CoinTracking.

## Exchanges soportados

- Binance
- Coinbase
- Kraken
- Bybit
- OKX
- KuCoin
- BingX

## Contenidos por exchange

Para cada exchange, la documentación cubre:

- Endpoints de API y autenticación
- Formatos de exportación CSV
- Tipos de transacciones
- Manejo de comisiones
- Formatos de dirección
- Problemas de datos conocidos
- Patrones de reconciliación
- Errores comunes de importación

## Documentación de exchanges

- **Official:** [`official/BINANCE.md`](official/BINANCE.md) — particularidades de importación de Binance en CoinTracking: límite histórico de la API (sept-2022), productos con soporte parcial/nulo (Earn/Staking/Savings, Futuros >3 meses, Launchpad, Auto-Invest), y casos especiales de mapeo (Swaps, Dust→BNB).

## Contexto regulatorio relevante

- [`reference/context/BINANCE_EU_MICA_EXIT.md`](../reference/context/BINANCE_EU_MICA_EXIT.md) — salida de Binance de la UE por MiCA (2026-07): qué buscar en la reconciliación cuando el usuario migra de exchange por este motivo (transferencias vs. posibles conversiones forzosas imponibles).

- [`reference/context/EXCHANGE_REGULATORY_UPDATES_2026.md`](../reference/context/EXCHANGE_REGULATORY_UPDATES_2026.md) — cambios regulatorios y operativos de 2026 (Binance MiCA, USDT→USDC Q1 2025, Coinbase expansión EU, BingX derivados): impacto en auditoría de CoinTracking y checklist para próximas auditorías.
