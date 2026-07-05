# Registros de Decisiones Arquitectonicas (ADRs)

Este directorio contiene las decisiones arquitectonicas del proyecto, documentadas en formato MADR (Markdown Any Decision Record).

## Indice

- [ADR-001: Idioma del repositorio (contenido en español, identificadores en inglés)](./0001-idioma-del-repositorio-contenido-en-español,-identificadores-en-inglés.md)
- [ADR-002: Fuente de verdad — jerarquía de autoridades](./0002-fuente-de-verdad.md)
- [ADR-003: Modelo de transacciones canónico](./0003-modelo-de-transacciones.md)
- [ADR-004: Estrategia de desarrollo (híbrido pragmático)](./0004-estrategia-de-desarrollo-híbrido-pragmático.md)
- [ADR-005: Zona horaria de importación y normalización a UTC](./0005-zona-horaria-de-importación-y-normalización-a-utc.md)
- [ADR-006: El producto es un agente de IA auditor (Claude Code) sobre la base de conocimiento](./0006-el-producto-es-un-agente-de-ia-auditor-claude-code-sobre-la-base-de-conocimiento.md)
- [ADR-007: Limpieza del repositorio (alineación con el enfoque agente)](./0007-limpieza-del-repositorio-alineación-con-el-enfoque-agente.md)
- [ADR-008: Vigencia y actualización del conocimiento (fiscal y CoinTracking)](./0008-vigencia-y-actualización-del-conocimiento-fiscal-y-cointracking.md)
- [ADR-009: Protocolo de agente crítico (cero invención, máxima cautela)](./0009-protocolo-de-agente-crítico-cero-invención,-máxima-cautela.md)
- [ADR-010: Eficiencia de tokens y caché de datos de CoinTracking](./0010-eficiencia-de-tokens-y-caché-de-datos-de-cointracking.md)
- [ADR-011: Persistencia y trazabilidad del flujo (nada sin dejar rastro)](./0011-persistencia-y-trazabilidad-del-flujo-nada-sin-dejar-rastro.md)
- [ADR-012: División de responsabilidades (Claude Code gestiona, Copilot explota)](./0012-división-de-responsabilidades-claude-code-gestiona,-copilot-explota.md)
- [ADR-013: Estructura multi-proyecto obligatoria (datos de usuario y estado; MCP pospuesto)](./0013-estructura-multi-proyecto-obligatoria-datos-de-usuario-y-estado;-mcp-pospuesto.md)
- [ADR-014: Validación de duplicados con trade_id y consentimiento explícito](./0014-validación-de-duplicados-con-trade_id-y-consentimiento-explícito.md)
- [ADR-015: Integración de la base de casos ChatGPT como v2 curada (patrones de reconciliación)](./0015-integración-de-la-base-de-casos-chatgpt-como-v2-curada-patrones-de-reconciliación.md)
- [ADR-016: Cambio de proyecto activo en caliente en el MCP (`cointracking_switch_project`)](./0016-cambio-de-proyecto-activo-en-caliente-en-el-mcp-`cointracking_switch_project`.md)
- [ADR-017: Protocolo de diagnóstico en orden fijo para la auditoría (endurecer falsos positivos)](./0017-protocolo-de-diagnóstico-en-orden-fijo-para-la-auditoría-endurecer-falsos-positivos.md)
- [ADR-018: Discrepancia `get_gains` vs FIFO manual — documentar como hipótesis, no automatizar en `ct_audit.py` (aún)](./0018-discrepancia-`get_gains`-vs-fifo-manual-—-documentar-como-hipótesis,-no-automatizar-en-`ct_audit.py`-aún.md)
- [ADR-019: Cierre y corrección de ADR-018 — `get_gains` confirmado fiable, la reconstrucción FIFO manual era la que fallaba](./0019-cierre-y-corrección-de-adr-018-—-`get_gains`-confirmado-fiable,-la-reconstrucción-fifo-manual-era-la-que-fallaba.md)
- [ADR-020: `get_historical_summary` puede devolver un punto fuera del rango `end` pedido — filtrar por fecha en el consumidor](./0020-`get_historical_summary`-puede-devolver-un-punto-fuera-del-rango-`end`-pedido-—-filtrar-por-fecha-en-el-consumidor.md)
- [ADR-021: Gate explícito de artefacto antes de cerrar la cifra anual exacta en `spanish-tax-return`](./0021-gate-explícito-de-artefacto-antes-de-cerrar-la-cifra-anual-exacta-en-`spanish-tax-return`.md)
- [ADR-022: Tercera pata de vigencia — contexto regulatorio/operativo de exchanges (extiende ADR-008)](./0022-tercera-pata-de-vigencia-—-contexto-regulatorio/operativo-de-exchanges-extiende-adr-008.md)
- [ADR-023: El MCP es dueño del ciclo de vida de sus archivos de caché (`cointracking_delete_project`)](./0023-el-mcp-es-dueño-del-ciclo-de-vida-de-sus-archivos-de-caché-`cointracking_delete_project`.md)
- [ADR-024: Formato "bloque-resumen" obligatorio al guiar altas/correcciones manuales en CoinTracking](./0024-formato-"bloque-resumen"-obligatorio-al-guiar-altas/correcciones-manuales-en-cointracking.md)
- [ADR-025: Formatos "CT-List" para mostrar listas de operaciones/hallazgos en la conversación](./0025-formatos-"ct-list"-para-mostrar-listas-de-operaciones/hallazgos-en-la-conversación.md)
- [ADR-026: Límites de decisión fiscal (matriz A/B/C)](./0026-limites-decision-fiscal-agente.md)
- [ADR-027: Integración de nuevos exchanges (protocolo de 4 fases)](./0027-integracion-nuevos-exchanges.md)
- [ADR-028: Límite auditor / asesor fiscal](./0028-limite-auditor-asesor-fiscal.md)
- [ADR-029: Protocolo de no-hacer (10 prohibiciones explícitas)](./0029-protocolo-de-no-hacer.md)
- [ADR-030: Validación y verificación de ADRs críticos](./0030-validacion-verificacion-adrs-criticos.md)
- [ADR-031: Plazos y períodos de declaración (Hacienda)](./0031-plazos-periodos-declaracion-hacienda.md)
- [ADR-032: Knowledge with Temporal Validity (metadatos YAML)](./0032-knowledge-temporal-validity.md)
- [ADR-033: Sistema de conocimiento jerárquico (6 niveles A-F)](./0033-sistema-conocimiento-jerarquico.md)
- [ADR-034: Stack de tecnología Python](./0034-stack-de-tecnología-python.md)
- [ADR-035: Representación del modelo de dominio en Python](./0035-representación-del-modelo-de-dominio-en-python.md)

## Proceso de ADR

Toda decision arquitectonica importante debe:

1. Ser propuesta en una rama nueva con un ADR borrador
2. Ser discutida en revision de codigo
3. Ser revisada por arquitecto del proyecto
4. Ser aprobada por el equipo
5. Ser completada con la decision final

Las decisiones menores pueden ser documentadas informalmente en CHANGELOG.md.
