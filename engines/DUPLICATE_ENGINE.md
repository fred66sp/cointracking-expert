# Especificación del motor de duplicados

**Detección de duplicados de transacciones**

El motor de duplicados identifica transacciones duplicadas resultantes de múltiples fuentes de importación, fallos de API, re-entrada manual o corrupción de datos. Soporta matching exacto y detección probabilística.

## Propósito

Identificar y reportar transacciones duplicadas con niveles variados de confianza, soportando revisión manual y remediación.

## Entradas

- Dataset de transacciones completo
- Reglas de detección de duplicados y umbrales
- Metadatos de fuente de transacción (CSV, API, manual, etc.)

## Salidas

- Lista de duplicados exactos
- Lista de duplicados probables con puntuaciones de confianza
- Análisis de impacto de duplicados (qué transacciones causan problemas)

## Responsabilidades

1. Detectar duplicados exactos (idénticos todos los campos)
2. Detectar duplicados probables (similares pero no idénticos)
3. Detectar duplicados entre fuentes (CSV + API + manual)
4. Cuantificar impacto de duplicados en balances y tenencias
5. Sugerir remediación (cuál remover)

## Algoritmos clave

- Matching exacto en características de transacción
- Matching difuso para detección probabilística
- Detección de duplicados entre fuentes
- Análisis de impacto

## Casos extremos

- Duplicados parciales (algunos campos difieren)
- Comisiones registradas separadamente
- Diferencias de redondeo
- Diferencias de timestamp
- Transacciones divididas vs. transacciones única
- Depósitos vs. transacciones "compra" (misma acción, diferente etiqueta)
