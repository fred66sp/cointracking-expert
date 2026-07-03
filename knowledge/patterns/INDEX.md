# Base de conocimiento de patrones de reconciliación

**Documentación de referencia para patrones comunes de reconciliación**

Esta sección documenta patrones comunes de reconciliación, problemas de datos y soluciones encontrados en carteras de criptomonedas del mundo real. Estos patrones ayudan al agente a **encontrar y explicar** problemas (ADR-006): son conocimiento de diagnóstico cualitativo, no cifras fiscales vinculantes.

## Documentos

- **[cointracking_casos_v2.yaml](cointracking_casos_v2.yaml)** — base de casos **vigente**, 20 casos con esquema canónico homogéneo (ver §"Esquema" abajo). Úsala como fuente principal de patrones al auditar o guiar correcciones.

## Estado de las bases de casos (histórico e integración)

| Fichero | Estado | Notas |
|---|---|---|
| `knowledge/patterns/cointracking_casos_v2.yaml` | ✅ **Vigente** | Base curada v2. Integra y sustituye a las dos anteriores. |
| `cointracking_casos_extended.yaml` (raíz del repo) | 🗑️ Eliminado (2026-07-03) | Candidato original aportado por el agente ChatGPT auxiliar. Su contenido quedó íntegramente integrado en v2; recuperable del historial de git si hiciera falta. |
| `cointracking_casos_base.yaml` (raíz del repo) | 🗑️ Eliminado (2026-07-03) | Baseline anterior a la integración (esquema mínimo: id/título/categoría/síntomas/causa/diagnóstico/solución/riesgo). Su contenido ya estaba cubierto y ampliado por v2 (CT-001↔CT-002/CT-017, CT-002↔CT-004, CT-003↔CT-003/CT-016, CT-004↔CT-006, CT-005↔CT-005/CT-011/CT-014, CT-006↔CT-007/CT-013, CT-007↔CT-012, CT-008↔CT-009, CT-009↔CT-015, CT-010↔CT-018). Recuperable del historial de git si hiciera falta (ver ADR-015).

## Esquema canónico de `cointracking_casos_v2.yaml`

Cada caso contiene los 16 campos siguientes (todos obligatorios; el valor es `null` cuando el campo no aplica al caso, nunca cadena vacía):

`id`, `titulo`, `categoria`, `sintomas`, `causa_probable` (`hecho`/`hipotesis`/`supuesto`), `evidencia_minima`, `pasos_diagnostico`, `solucion_recomendada`, `anti_patron`, `por_que_falso_positivo`, `nivel_confianza` (`verificado`/`probable`/`hipotesis`/`pendiente_verificar`), `nivel_riesgo` (`bajo`/`medio`/`alto`/`critico`), `impacto_fiscal_potencial`, `senales_tempranas`, `validacion_antes_despues` (`antes`/`despues`), `vigencia` (`fecha_revision`/`motivo_caducidad_potencial`/`fuente_recomendada_para_revalidar`).

Categorías usadas: `transferencias_huerfanas`, `ventas_sin_base_de_coste`, `duplicados`, `saldos_imposibles_o_negativos`, `rendimientos`, `permutas_complejas`, `casos_limite_espana`.

## Política de vigencia y revalidación (ADR-008)

- Cada caso declara su propia `vigencia` (fecha de revisión, motivo de posible caducidad, fuente recomendada). Antes de citar un caso al usuario, comprueba que su `fecha_revision` no está desfasada frente a cambios conocidos de CoinTracking o de la normativa fiscal citada en `impacto_fiscal_potencial`.
- Los casos con `nivel_confianza: pendiente_verificar` o `hipotesis` (p. ej. CT-010, CT-015, CT-018) **no se presentan al usuario como hecho establecido**; requieren reverificación contra la fuente indicada en `fuente_recomendada_para_revalidar` antes de usarse en un informe.
- Varios casos (CT-003, CT-008, CT-016, CT-019) están alineados con **DECISIONS.md#ADR-014** (validación de duplicados con `trade_id` y consentimiento explícito); si ADR-014 cambia, revisar estos casos en el mismo commit.

## Categorías de patrones

- Patrones de calidad de datos
- Patrones de errores de importación
- Patrones de datos faltantes
- Técnicas de reconciliación
- Resoluciones de problemas comunes
- Mejores prácticas para gestión de cartera

## Contenidos

La documentación cubre:

- Patrones de transacciones duplicadas
- Detección de transacciones faltantes
- Patrones de discrepancias de balance
- Patrones de emparejamiento de transferencias
- Variaciones de cálculo de comisión
- Inconsistencias de nomenclatura de activos
- Diferencias de formato de datos de exchange
- Patrones de recuperación de datos corrupto
