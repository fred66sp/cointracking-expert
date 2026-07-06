# ADRs — Índice de Decisiones Arquitectónicas

**Ubicación:** `adr/`

**Característica:** Registro de decisiones arquitectónicas (MADR format).

**Total:** 42 ADRs (001-042, todos Accepted salvo notas puntuales — ver `adr/README.md` para el índice secuencial completo, este documento los agrupa por nivel de conocimiento). Corregido 2026-07-05: este índice llevaba desactualizado desde la creación de ADR-033, con conteo, temas y sección "Pendientes" incorrectos.

---

## Por Nivel de Conocimiento (ADR-033)

### Nivel A — Fuentes Oficiales

**Goviernan:** Qué fuentes son autoritativas

- **ADR-002:** Fuente de verdad (jerarquía de autoridades)
  - Blockchain > API > CSV > CoinTracking > Inferencia
  - Cómo resolver conflictos entre fuentes

---

### Nivel B — Operativo

**Goviernan:** Cómo funciona realmente

- **ADR-003:** Modelo de transacciones canónico
  - 10 tipos: Buy, Sell, Transfer, Deposit, Withdrawal, Staking, Airdrop, Fee, Convert, Futures
  - Campos obligatorios, validaciones

- **ADR-004:** Reconciliación con datos reales
  - Verificar contra exchange/blockchain, no solo CT
  - Cotejar depósitos, retiradas, saldos

---

### Nivel C — Casos

**Goviernan:** Qué hacemos ante casos específicos

- **ADR-014:** Validación de duplicados con Trade ID
  - Regla: misma fecha ≠ duplicado (caso FLOKI)
  - Exigir confirmación explícita antes de borrar

- **ADR-026:** Matriz de decisiones A/B/C
  - Categoría A: agente decide solo
  - Categoría B: requiere confirmación usuario
  - Categoría C: delega a humano

---

### Nivel D — Auxiliar

**Goviernan:** Herramientas y estructuras

- **ADR-017:** Protocolo de diagnóstico en orden fijo
  - 6 fases: cobertura → duplicados → transfers → tipos → pool → fiscal
  - Reduce falsos positivos

- **ADR-024:** Formato bloque resumen (CT-Task)
  - Guía de correcciones manuales en CoinTracking
  - Estructura: Tipo | Fecha | Campos | Exchange | Grupo | Comentario

- **ADR-025:** Migración de DECISIONS.md a MADR individual
  - Un ADR = un archivo
  - Mejora navegabilidad

---

### Nivel E — Referencia

**Goviernan:** Contexto y convenciones

- **ADR-001:** Principios de auditoría
  - Cero invención, trazabilidad, consentimiento

- **ADR-005:** Zona horaria (Europe/Madrid con DST, UTC interno)

- **ADR-008:** Vigencia del conocimiento
  - Documentos envejecen → revisar periódicamente

---

### Nivel F — Governance

**Goviernan:** Cómo funciona el sistema completo

- **ADR-006:** Límites del determinismo (LLM, no motor exacto)
  - Cualitativo: sí. Cifras fiscales vinculantes: no.

- **ADR-009:** Protocolo crítico (10 prohibiciones)
  - Nunca borrar sin triple confirmación
  - Nunca ocultar incertidumbre
  - Nunca inferir origen de fondos

- **ADR-010:** Eficiencia de tokens y caché
  - Cachear a disco, reutilizar
  - Procesar grandes volúmenes con scripts

- **ADR-011:** Persistencia y trazabilidad
  - Informes en `reports/output/`
  - REGISTRO-CAMBIOS append-only
  - Memoria durable entre sesiones

- **ADR-012:** División de responsabilidades (Claude Code vs. Copilot)
  - Claude Code: gestiona agente
  - Copilot: explota agente
  - AGENT_CHANGE_REQUESTS.md: bandeja de peticiones

- **ADR-013:** Proyecto activo obligatorio (multi-proyecto)
  - `USER_INPUT/<proyecto>/` aísla datos
  - MCP sincronizado en caliente

- **ADR-016:** MCP sincronizado en caliente (ADR-013)
  - `cointracking_switch_project()` sin reiniciar

- **ADR-015:** Integración de casos ChatGPT (v2 curada)
  - Casos v2: esquema canónico, 20 casos
  - Legacy deprecado

- **ADR-018:** Discrepancia FIFO manual vs. `get_gains`
  - Hipótesis abierta (ADR-019 la corrige)

- **ADR-019:** Cierre de ADR-018 (confiar en `get_gains` oficial)
  - Tax Report oficial = fuente de verdad
  - Reconstrucciones manuales → errores

- **ADR-020:** Verificación/normalización `historical_summary` con parámetros temporales
  - API añade punto "actual" adicional
  - Mitigación: filtrar serie en consumidor

- **ADR-021:** Precondición de Tax Report para cerrar cifra anual
  - No cerrar base del ahorro sin Tax Report oficial
  - Gate explícito en skill

- **ADR-022:** Vigencia del conocimiento (ADR-008, operacionalizado)
  - Fiscal, CoinTracking, regulatorio: fechas vigencia verificables

- **ADR-023:** El MCP es dueño del ciclo de vida de sus archivos de caché
  - Tool `cointracking_delete_project` — borra la caché de un proyecto sin colisionar con el proceso del servidor

- **ADR-024:** Formato CT-Task para operaciones manuales
  - Bloque-resumen obligatorio tras guía de correcciones
  - Estructura: Tipo | Fecha | Campos | Exchange | etc.

- **ADR-025:** Migración DECISIONS.md → MADR individual
  - 25 ADRs migrados (001-025)
  - Después: 026-033+

- **ADR-027:** Protocolo de integración de nuevos exchanges
  - 4 fases: preintegración, importación, validación, documentación

- **ADR-028:** Límite auditor / asesor fiscal (Zona A/B/C)
  - Qué puede decir el agente
  - Qué nunca puede decir
  - Qué delega a humanos

- **ADR-029:** Protocolo de no-hacer (10 prohibiciones)
  - Nunca recomendar borrado sin confirmación triple
  - Nunca modificar sin registro explícito
  - [9 más]

- **ADR-030:** Validación y verificación de ADRs críticos
  - Checklist para ADRs sobre cifras/hechos mutables
  - Criticidad: alto/medio/bajo

- **ADR-031:** Validación temporal previa de obligaciones fiscales
  - Máquina de estados: ORDINARIO | LATE | FUTURE | UNKNOWN | EXPIRED
  - Sin hardcoding de fechas (ADR-032)

- **ADR-032:** Knowledge with Temporal Validity (metadatos YAML)
  - Frontmatter obligatorio (id, authority, vigencia, confidence, etc.)
  - Protocolo de validación antes de usar dato

- **ADR-033:** Sistema de Conocimiento Jerárquico
  - 6 niveles (A-F) con autoridad clara
  - Estructura de directorios + INDEX_MASTER
  - Metadatos YAML operacionalizan ADR-032

- **ADR-034/035:** históricos del framework Python/SDK descartado (ADR-006/007) — títulos internos "ADR-002"/"ADR-003" preservados tal cual quedaron, no confundir con los ADR-002/003 vigentes (Fuente de verdad / Modelo de transacciones). Se conservan por trazabilidad, no son decisiones activas.

- **ADR-036:** Convención de IDs de documentos de conocimiento (`KB-[NIVEL][SUBSECCIÓN]-[NÚMERO]`)

- **ADR-037:** Validación obligatoria en desarrollo — pre-commit hook ejecuta `tools/audit_mega_complete.py`

- **ADR-038:** Criterio de auditoría en lotes (no iterativa) — 3 pasadas limpias consecutivas antes de declarar "sistema limpio"

- **ADR-039:** Optimización de tokens y caché — arquitectura de 3 capas (persistente, agregados, procesamiento local), TTL dinámico y versionado automático por knowledge/

- **ADR-040:** Credenciales por proyecto en el MCP — multi-cuenta opcional vía `--project-env-dir` + `<proyecto>.env` (fail-closed; los secretos nunca pasan por la conversación)

- **ADR-041:** Procesos MCP huérfanos en Windows — síntoma "file busy" al borrar cachés; protocolo: contar ventanas abiertas, matar viejo→nuevo verificando vida del MCP propio tras cada baja

- **ADR-042:** Proactividad en gobernanza — detectar que algo merece ADR es responsabilidad del agente: proponer en el momento, pedir permiso (Cat. B), crear con visto bueno; 5 disparadores + guardarraíl de proporcionalidad

---

## Relacionados Cruzados

### Por ADR

Ejemplo: ADR-031 (Validación temporal) → usa **ADR-032** (metadatos) + **Nivel A1** (FILING_DEADLINES.md)

### Por Nivel

Ejemplo: **Nivel C1** (Casos) → referencia **ADR-014** (validación duplicados) + **ADR-026** (matriz A/B/C)

---

## Corrección 2026-07-05: los "pendientes" de Fase 3+ ya se resolvieron (por otra vía)

Esta sección listaba temas de "Capa 2: Conciliación" reservados para ADR-034 a 040. Esos números ya se usaron para decisiones no relacionadas (034/035 son históricos del framework descartado; 036-039 son gobernanza de desarrollo y caché — ver arriba). Los temas en sí **no quedaron sin resolver**: se cubrieron como documentos de Nivel C (patrones/procedimientos), no como ADRs nuevos:

| Tema pendiente listado | Resuelto en |
|---|---|
| Transfers (emparejar withdrawal/deposit) | `knowledge/patterns/PATTERN_TRANSFER_MATCHING.md` (KB-C2-003) |
| Duplicados (matriz de clasificación) | `knowledge/patterns/PATTERN_DUPLICATE_DETECTION.md` (KB-C2-001) |
| Missing Purchase History (operativo) | `knowledge/patterns/PATTERN_PURCHASE_POOL_EXHAUSTION.md` (KB-C2-004) |
| Modelo de balances (qué es "correcto") | `knowledge/patterns/PATTERN_BALANCE_RECONCILIATION.md` (KB-C2-002) |
| Flujo de conciliación (pipeline invariante) | `knowledge/procedures/PROCEDURE_AUDIT_ACCOUNT.md` (KB-C3-001) |
| Cost Basis / FIFO (cuándo confiar en CT) | `knowledge/cointracking/official/COST_BASIS_AND_VALIDATION.md` (Nivel A2) |

Si en el futuro se necesita un ADR nuevo, usar el siguiente número libre (043+), nunca reutilizar uno ya asignado.

---

## Cómo Usar Este Índice

1. **Por problema:** "Tengo duplicados" → ADR-014 / ADR-026 → nivel C/D
2. **Por nivel:** "Necesito fuentes oficiales" → Nivel A → ADRs 002, etc.
3. **Por ADR:** "¿Qué ADRs afectan a la auditoría?" → Niveles B/C/D
4. **Por referencia cruzada:** "¿Qué depende de ADR-032?" → Buscar en la sección

---

## Status

- 42 ADRs creados y formalizados (todos Accepted salvo notas puntuales — ver `adr/README.md`)
- Gobernanza de desarrollo (036-038), optimización de caché (039) y Nivel C (patrones/procedimientos) cubren lo que esta sección marcaba como pendiente/bloqueante
- Pendiente real, sin resolver: la discrepancia de declaraciones ya presentadas del proyecto `agp2025` (ver memoria de proyecto, no es un tema de ADR)

Ver `docs/ADR_GAP_ANALYSIS_2026-07-05.md` para el análisis de brechas original (puede estar igualmente desactualizado; verificar antes de citarlo).
