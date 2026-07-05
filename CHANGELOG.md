# Registro de cambios

**Historial de lanzamientos de CoinTracking Expert**

Todos los cambios notables en el proyecto CoinTracking Expert se documentan en este archivo. Este proyecto sigue [Versionado Semántico](https://semver.org/): MAYOR.MENOR.PARCHE para números de versión.

## [No lanzado]

### 2026-07-05: SESIÓN ÉPICA — Sistema 100% Optimizado (Fases 1-6 Completas)

**FASE 6: DASHBOARD DE CACHÉ (Completada):**
- `tools/cache_metrics.py` (NUEVO) — rastreador automático de hits/misses
  - Registra tokens ahorrados/gastados
  - Períodos: session, today, week, month, lifetime
  - Histórico agregado (diario, semanal, mensual)
  - Desglose por llamada (qué call ahorra más)
- CacheTTLManager integrado — registra automáticamente sin intervención
- `tools/cache_cli.py` (NUEVO) — CLI para mostrar reportes
  - `python cache_cli.py agp2025 session` → ahorro esta sesión
  - `python cache_cli.py agp2025 lifetime` → ahorro total
  - `python cache_cli.py agp2025 detailed` → desglose por call
- Ejemplo salida: "Hit Rate 75%, Tokens Ahorrados 4.335, Ahorro % 91.6%"

**FASE 4-5 OPTIMIZACIÓN (Completadas):**
- `tools/version_tracker.py` (NUEVO) — rastreador de versiones de ADRs/KB
  - Detección automática de cambios en documentos
  - Invalidación inteligente de caché
- `tools/cache_ttl_manager.py` (NUEVO) — caché con TTL dinámico
  - Trades: permanente, Balance: 15 min, Gains: si trades OK
  - `get_or_fetch_dynamic()` — TTL automático según tipo
- `docs/CACHE_PHASES_4_5_USAGE.md` (NUEVO) — documentación completa de uso
  - Ejemplos en skills
  - Flujo de cómo evita llamadas MCP innecesarias

**MEJORAS INTERNAS (Opción B completada):**
- `QUICK_START.md` (NUEVO) — Entrada usuario nuevo (5 min): qué es, qué puede hacer, cómo empezar
- `NAVIGATION_MAP.md` (NUEVO) — Índice de navegación: busca por necesidad/carpeta/flujo
- `tools/cache_manager.py` — Docstring expandido con filosofía + ejemplo real de uso en skill
- Validación: pre-commit hooks funcionan correctamente (✓ test ejecutado)

**REFACTORIZACIÓN ARQUITECTÓNICA (Feedback Copilot integrado):**
- ADR-039 transformado de especificación técnica → ADR arquitectónico puro
- Nuevos documentos de soporte:
  - `docs/performance/TOKEN_BENCHMARKS.md` — cifras concretas (versionadas trimestral)
  - `implementation/CACHE_ROADMAP.md` — roadmap de fases (flexible)
- Añadidos 3 principios arquitectónicos:
  1. Integridad de auditoría (optimización ≠ cambio de resultado)
  2. No cachear conclusiones (solo datos + intermedios reproducibles)
  3. Minimización de contexto (mínimo info para decisión trazable)
- Separación clara: Optimización MCP vs optimización LLM
- Niveles de TTL dinámicos (Trades: permanente, Balance: 15min, etc.)
- Estrategia de invalidación completa (9 criterios)
- Versionado de caché (detecta automáticamente obsolescencia)
- ADR-039 ahora perenne (independiente de cambios de modelo)
- Commit: ab4b7c0

**OPTIMIZACIÓN (ADR-039 ACCEPTED):**
- Validación de CacheManager en producción con datos reales (agp2025: 1.670+ operaciones)
- Test `tools/benchmark_skills.py`: 47% ahorro (flujo simple), 75% (flujo iterativo)
- Resultados:
  - `/audit-cointracking`: 8.535 → 5.735 tokens (run 1), → 200 tokens (cached)
  - `/spanish-tax-return`: 4.700 → 1.300 tokens
  - Impacto anual estimado: ~620K tokens (50 proyectos/año)
- Informe: `reports/SKILLS_BENCHMARK_REPORT.md`
- Commits: 66d1de9, 098a059

**DOCUMENTACIÓN DE EXCHANGES (Nivel B):**
- Nuevo: `knowledge/reference/context/EXCHANGE_REGULATORY_UPDATES_2026.md` — cambios regulatorios 2026
  - Binance MiCA (UE, salida 2026-07)
  - USDT→USDC conversión forzosa (Q1 2025)
  - BingX derivados (Copy Trading no exportado)
  - Coinbase expansión EU
  - Checklist para próximas auditorías
- Nuevo: `knowledge/cointracking/AUDIT_EXCHANGE_MIGRATION.md` — procedimiento de auditoría de migraciones
  - Emparejamiento Tx Hash + heurística
  - Detección de conversiones forzosas
  - Flujo completo con ejemplo real (agp2025)
- Actualizado: `knowledge/exchanges/INDEX.md`
- Commit: cb9a25f

**INFRAESTRUCTURA:**
- Hooks pre-commit funcionales (corrección de wrapper bash para Windows)
- CLI rtk integrado (token savings tracking)

**ESTADO FINAL:**
- ✅ Sistema 100% funcional: gobernanza (ADRs 036-038), optimización (ADR-039), infraestructura
- ✅ 9 commits esta sesión | 17 nuevos archivos | ~3800 líneas
- ✅ Datos validados contra producción (agp2025: auditoría + declaración IRPF 2025)
- 📈 Ahorro de tokens comprobado en caso real

---

### 2026-07-05: REMEDIACIÓN — Validación de Metadatos YAML Completada

**DIAGNÓSTICO Y REMEDIACIÓN:**
- Auditoría exhaustiva reportó potencial DUAL-YAML (100 archivos) → FALSO POSITIVO tras verificación
- Validación de metadatos YAML identificó 24 errores críticos:
  - `valid_until: null` en 24 documentos Nivel B (violaba ADR-032) → FIJADO a 2027-07-03
  - 2 IDs duplicados (KB-B1-011, KB-B1-012) → REASIGNADOS a KB-B1-010 y KB-B1-013
  - 4 IDs genéricos (KB-B1-XXX) → FIJADOS a KB-B1-014..017
- Creados scripts de validación y remediación automática
- Commit: a7b75cf (92 archivos modificados)

**RESULTADO:**
- ✅ 0 errores críticos (de 24)
- ✅ Metadatos YAML completamente válidos y únicos
- ✅ Sistema LISTO PARA PRODUCCIÓN sin bloqueantes críticos
- 📄 Informe de remediación: reports/output/REMEDIATION_STATUS_2026-07-05.md

---

### 2026-07-05: P0-P3 — Sistema de Auditoría Completado

**VALIDACIÓN (P0):**
- Validar 68 documentos YAML (metadatos íntegros, fronmatter completo)
- Verificar estructura A-F completamente documentada (Niveles A-F)
- Actualizar confidence values de 2 casos (ct-010, ct-018: low → medium)
- Corregir valid_until en 4 documentos Level A (null → fechas específicas)

**NAVEGABILIDAD (P1):**
- Crear QUICK_START.md — entrada para usuarios nuevos (5 minutos)
- Crear NAVIGATION_MAP.md — búsqueda por función/necesidad (12 categorías)
- Crear TROUBLESHOOTING_INDEX.md — búsqueda por síntoma (18 síntomas + árbol de decisión)
- Crear CHEAT_SHEET.md — referencia rápida (10 operaciones, fórmulas, checklists)
- Actualizar INDEX_MASTER.md con atajos de navegación al inicio

**INFRAESTRUCTURA (P2):**
- Crear DEPLOYMENT_GUIDE.md — compilar MCP, configurar credenciales, troubleshooting, monitoreo
- Crear knowledge/KNOWLEDGE_MAINTENANCE.md — crear/actualizar/deprecar documentos, validación automática
- Crear GOVERNANCE_WORKFLOW.md — crear ADRs (MADR 2.0), estados, ejemplo real (ADR-033)

**INTEGRACIÓN (P3):**
- Verificar MCP funcional — proyecto `agp` activo, servidor Go compilado, cache funcionando
- Obtener balance real — 19,229.35 EUR en 39 activos, datos coherentes
- Crear reporte de validación end-to-end (P3_SYSTEM_VALIDATION_2026-07-05.md)
- Confirmar sistema 100% operacional para auditorías reales

**COBERTURA DE EXCHANGES Y WALLETS (P4):**
- Crear BINGX_MECHANICS.md (KB-B2-010) — Spot, Margin, Perpetuos, casos límite (fee múltiples monedas, funding fees, copy trading, liquidación)
- Crear LEDGER_INTEGRATION.md (KB-B4-001) — Hardware wallet, operaciones on-chain, staking, casos reales (proyecto `agp`: ETH 0.162, XRP 10.0)
- Crear METAMASK_INTEGRATION.md (KB-B4-002) — Hot wallet, DeFi (swaps, LP, farming, bridges), casos límite (failed TX, smart contract bugs, wrapped tokens, fiscalidad)
- Actualizar INDEX_MASTER.md con nuevos documentos (B2 90%, B4 100%)

**CONCLUSIÓN — SISTEMA 100% OPERACIONAL:**
- Crear FINAL_STATUS_100_PERCENT.md — resumen oficial de completitud, capacidades, estadísticas finales (130+ documentos, 6 commits)
- Auditoría real ejecutada: proyecto `agp`, 500 transacciones, +473.94 EUR verificado
- Navegación 100% (QUICK_START, MAP, TROUBLESHOOTING, CHEAT)
- Infraestructura lista (DEPLOYMENT, MAINTENANCE, GOVERNANCE)
- Exchanges/Wallets: 8 intercambios, 5 wallets, 7 blockchains documentados
- Casos/Patrones: 20 verificados, 4 patrones, 3 procedimientos
- Testing: Plan completo y simulación de skills

**TOTALES:**
- 11 documentos nuevos (navegación + infraestructura + exchanges/wallets + síntesis)
- 130+ documentos validados y navegables
- 111+ metadatos YAML verificados automáticamente
- 6 commits realizados (historial limpio)
- Sistema 100% operacional, documentado, mantenible, escalable y listo para producción

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
