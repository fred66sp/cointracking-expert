# Binance — particularidades de importación en CoinTracking

**Fuente:** centro de ayuda oficial de CoinTracking — [Binance FAQ](https://cointracking.freshdesk.com/en/support/solutions/articles/29000049344-binance-faq) y [Binance Import Restrictions](https://cointracking.freshdesk.com/en/support/solutions/articles/29000039887-binance-import-restrictions).
**Verificado:** 2026-07-04.

> ⚠️ Relevante para la auditoría: varias limitaciones de aquí explican por qué pueden faltar operaciones, aparecer sin base de coste o requerir reclasificación manual — ver `../cointracking/COST_BASIS_AND_VALIDATION.md` §3.

## 1. Límite histórico de la API de Binance

Binance **discontinuó el acceso por API a datos anteriores a septiembre de 2022**. Para historial previo a esa fecha, solo sirve el CSV exportado directamente desde Binance.

## 2. Productos NO soportados (o solo parcialmente) por la API de CoinTracking

| Producto/tipo | Soporte por API | Recomendación |
|---|---|---|
| Margin PnL, Leveraged Tokens, Battle | No soportado | Importar vía CSV |
| Futuros | Solo **últimos 3 meses** vía API | Datos anteriores requieren CSV |
| Savings, Staking, Liquid Swap, Pool Savings, Finance (Earn) | No soportado / parcial | Importar vía CSV; revisar tipo asignado tras importar |
| Binance Visa Card, Crypto Loans, Recurring Buy | No soportado | Importar vía CSV o entrada manual |
| Launchpad | Soporte parcial (no incluye el gasto en BNB del ticket) | Revisar manualmente |
| NFT (IGO, Marketplace) | No soportado | Entrada manual si es relevante |
| Auto-Invest | No importado | Entrada manual o CSV si Binance lo exporta |

## 3. Casos especiales de mapeo

- **Swaps:** se importan como **"gasto e ingreso (no gravable)"**, no como operación de `Trade`, por falta de Trade ID emparejable. Puede requerir reclasificación manual si se quiere tratar como permuta a efectos fiscales.
- **Dust → BNB (conversión de polvo):** la API **solo trae las últimas 100 conversiones**; el CSV las importa todas, pero puede haber pequeñas diferencias de cálculo entre ambas vías.

## 4. Recomendación oficial de importación

CoinTracking recomienda: **importar todo el histórico vía CSV** y **activar la API solo para operaciones futuras** (con fecha de inicio = hoy), para evitar los huecos anteriores. Si se usan ambos métodos (API + CSV), **revisar duplicados**.

## 5. Implicación para la auditoría (`audit-cointracking`)

- Si una cuenta Binance se importó **solo por API**, sospechar de entrada de: historial anterior a sept-2022 ausente, Earn/Staking/Savings mal clasificados o ausentes, futuros de más de 3 meses de antigüedad ausentes.
- Ante una venta sin base de coste o un saldo negativo en un activo relacionado con Earn/Staking/Swap, revisar primero si el origen fue uno de estos productos con soporte parcial antes de asumir un problema de datos distinto.
