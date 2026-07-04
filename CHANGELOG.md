# Registro de cambios

**Historial de lanzamientos de CoinTracking Expert**

Todos los cambios notables en el proyecto CoinTracking Expert se documentan en este archivo. Este proyecto sigue [Versionado Semántico](https://semver.org/): MAYOR.MENOR.PARCHE para números de versión.

## [No lanzado]

### Agregado
- **ADR-033: Sistema de Conocimiento Jerárquico** — arquitectura de 6 niveles (A-F) con metadatos YAML obligatorios, operacionaliza ADR-032 (Knowledge with Temporal Validity); incluye INDEX_MASTER.md (mapa navegable) y MIGRATION_PLAN.md (Fase 2-3)
- Agente auditor de CoinTracking en Claude Code (subagente + skill `/audit-cointracking`)
- Skill `/spanish-tax-return` para preparar la declaración de IRPF de un ejercicio, reconciliando primero (ADR-006)
- Base de conocimiento: formato CSV, modelo de coste, integración MCP y fiscalidad española (IRPF)
- Servidor MCP propio en Go (`cointracking-mcp/`), sustituyendo al servidor JS de terceros usado antes (`cointracking-mcp-main/`); incluye tools propios `cointracking_invalidate_cache`, `cointracking_cache_stats`, `cointracking_close_project` y `cointracking_switch_project` (cambio de proyecto activo en caliente, ADR-016)
- Estructura multi-proyecto (`USER_INPUT/<proyecto>/`, `reports/output/<proyecto>/`) para aislar datos entre casos (ADR-013)
- Persistencia y trazabilidad del flujo: informes en `reports/output/`, `REGISTRO-CAMBIOS.md` append-only, memoria durable entre sesiones (ADR-011)
- División de responsabilidades: Claude Code gestiona el agente, GitHub Copilot lo explota vía `.github/copilot-instructions.md`, con `AGENT_CHANGE_REQUESTS.md` como bandeja de peticiones de mejora desde el uso real (ADR-012)
- Base de casos/patrones de reconciliación curada (`knowledge/patterns/cointracking_casos_v2.yaml`, 20 casos, esquema canónico) reemplazando el YAML legacy (ADR-015)
- Conocimiento sobre contexto regulatorio/operativo de exchanges (`knowledge/exchanges/`), p. ej. la salida de Binance de la UE por MiCA (2026-07) y su impacto en reconciliación
- Registro de decisiones arquitectónicas (ADRs 001-032) migrado a formato MADR individual en `adr/` (una decisión por archivo), reemplazando el monolítico `DECISIONS.md` (ADR-025)
- **Nivel MVP de auditoría (4 ADRs críticos faltantes) + protocolo de validación:**
  - ADR-002: Jerarquía de fuentes de verdad — Blockchain > API > CSV > CoinTracking, con casos de resolución de conflictos
  - ADR-003: Modelo de transacciones — 10 tipos canonicales (Buy, Sell, Transfer, Deposit, Withdrawal, Staking, Airdrop, Fee, Convert, Futures) con campos obligatorios y validaciones
  - ADR-028: Límite auditor/asesor fiscal — dónde termina la auditoría técnica y empieza la asesoría (Zona A/B/C)
  - ADR-029: Protocolo de no-hacer — 10 prohibiciones explícitas (nunca borrar sin triple confirmación, nunca ocultar incertidumbre, nunca inferir origen, etc.)
  - ADR-030: Validación y verificación de ADRs críticos — checklist exhaustiva para ADRs sobre fiscalidad, cifras o hechos mutables; REGLA DE ORO: cifras fiscales viven en `knowledge/`, ADRs solo referencian
- **Arquitectura transversal de vigencia (2 ADRs):**
  - ADR-032: Knowledge with Temporal Validity — metadatos YAML para todo conocimiento que envejece (valid_from, valid_until, last_verified, source, confidence); protocolo de validación antes de usar dato; 3 niveles de criticidad
  - ADR-031: Validación temporal previa de obligaciones fiscales — máquina de estados (ORDINARIO | LATE | FUTURE | UNKNOWN | EXPIRED) sin hardcoding de fechas; integración en `/spanish-tax-return` Paso 0.5
- ADR-026: Matriz de decisiones explícita — qué decide el agente solo (Categoría A), qué requiere confirmación del usuario (Categoría B), qué delega a humanos (Categoría C), operacionalizando el límite de determinismo de ADR-006
- ADR-027: Protocolo de integración de nuevos exchanges en multi-proyecto, 4 fases obligatorias (preintegración con consentimiento, importación controlada, validación exhaustiva con 9 chequeos, documentación total), con 3 ejemplos prácticos y 7 pendientes de automatización
- Protocolo de diagnóstico en orden fijo para la auditoría (6 fases: cobertura → duplicados → transferencias → tipos/base de coste → purchase pool → cierre fiscal), endurecido contra falsos positivos (ADR-017)
- Validación de duplicados con `trade_id`/`Tx ID` y consentimiento explícito antes de cualquier borrado (ADR-014)
- Regla de reconciliar siempre depósitos/retiradas/saldos contra la fuente externa real (banco/exchange), no solo contra la coherencia interna de CoinTracking
- Glosario (`docs/GLOSSARY.md`) ampliado con terminología de CoinTracking/exchanges y los formatos y modos propios del auditor (CT-Task, CT-List, niveles de riesgo)
- Índice de troubleshooting por síntoma/warning (`knowledge/cointracking/TROUBLESHOOTING.md`), que enruta a los casos de `cointracking_casos_v2.yaml` y a `COST_BASIS_AND_VALIDATION.md` sin duplicar conocimiento

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
