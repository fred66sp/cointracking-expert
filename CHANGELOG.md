# Registro de cambios

**Historial de lanzamientos de CoinTracking Expert**

Todos los cambios notables en el proyecto CoinTracking Expert se documentan en este archivo. Este proyecto sigue [Versionado Semántico](https://semver.org/): MAYOR.MENOR.PARCHE para números de versión.

## [No lanzado]

### Agregado
- Agente auditor de CoinTracking en Claude Code (subagente + skill `/audit-cointracking`)
- Skill `/spanish-tax-return` para preparar la declaración de IRPF de un ejercicio, reconciliando primero (ADR-006)
- Base de conocimiento: formato CSV, modelo de coste, integración MCP y fiscalidad española (IRPF)
- Servidor MCP propio en Go (`cointracking-mcp/`), sustituyendo al servidor JS de terceros usado antes (`cointracking-mcp-main/`); incluye tools propios `cointracking_invalidate_cache`, `cointracking_cache_stats`, `cointracking_close_project` y `cointracking_switch_project` (cambio de proyecto activo en caliente, ADR-016)
- Estructura multi-proyecto (`USER_INPUT/<proyecto>/`, `reports/output/<proyecto>/`) para aislar datos entre casos (ADR-013)
- Persistencia y trazabilidad del flujo: informes en `reports/output/`, `REGISTRO-CAMBIOS.md` append-only, memoria durable entre sesiones (ADR-011)
- División de responsabilidades: Claude Code gestiona el agente, GitHub Copilot lo explota vía `.github/copilot-instructions.md`, con `AGENT_CHANGE_REQUESTS.md` como bandeja de peticiones de mejora desde el uso real (ADR-012)
- Base de casos/patrones de reconciliación curada (`knowledge/patterns/cointracking_casos_v2.yaml`, 20 casos, esquema canónico) reemplazando el YAML legacy (ADR-015)
- Conocimiento sobre contexto regulatorio/operativo de exchanges (`knowledge/exchanges/`), p. ej. la salida de Binance de la UE por MiCA (2026-07) y su impacto en reconciliación
- Registro de decisiones ampliado a ADR-022 (ver `DECISIONS.md` para el índice completo)
- Protocolo de diagnóstico en orden fijo para la auditoría (6 fases: cobertura → duplicados → transferencias → tipos/base de coste → purchase pool → cierre fiscal), endurecido contra falsos positivos (ADR-017)
- Validación de duplicados con `trade_id`/`Tx ID` y consentimiento explícito antes de cualquier borrado (ADR-014)
- Regla de reconciliar siempre depósitos/retiradas/saldos contra la fuente externa real (banco/exchange), no solo contra la coherencia interna de CoinTracking

### Cambiado
- Giro de alcance: de framework/SDK de motores deterministas a agente de IA (ADR-006)
- Validado con un caso real completo (proyecto `agp2025`): reconciliación cerrada de Coinbase, Binance, BingX y Ledger, depósitos fiat verificados (34.000 €), y una declaración de IRPF 2025 preparada de principio a fin

### Eliminado
- Andamiaje del SDK descartado: paquetes Python vacíos, specs de motores, CI de pytest y documentos de la visión de framework (ADR-007)
- YAML legacy de casos de reconciliación (`cointracking_casos_base.yaml`), deprecado en favor de la v2 curada (ADR-015)

### Arreglado
- N/A

### Seguridad
- N/A

---

## Formato de versión

Cada lanzamiento incluye:
- Número de versión (versionado semántico)
- Fecha de lanzamiento
- Características agregadas (nueva funcionalidad)
- Características cambiadas (modificaciones a funcionalidad existente)
- Características deprecadas (a ser eliminadas en versiones futuras)
- Características eliminadas (funcionalidad previamente deprecada)
- Bugs arreglados y problemas
- Actualizaciones de seguridad

---

## Lanzamientos futuros

La dirección del proyecto se registra en las decisiones arquitectónicas de [DECISIONS.md](DECISIONS.md).
