# ADRs — Índice de Decisiones Arquitectónicas

**Ubicación:** `adr/`

**Característica:** Registro de decisiones arquitectónicas (MADR format).

**Total:** 33 ADRs (001-033, ADR-033 en Fase 1, otros pending)

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

- **ADR-023:** [PENDIENTE — slots para futuras ADRs]

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

---

## Relacionados Cruzados

### Por ADR

Ejemplo: ADR-031 (Validación temporal) → usa **ADR-032** (metadatos) + **Nivel A1** (FILING_DEADLINES.md)

### Por Nivel

Ejemplo: **Nivel C1** (Casos) → referencia **ADR-014** (validación duplicados) + **ADR-026** (matriz A/B/C)

---

## Pendientes (Fase 3+)

Slots reservados para futuras ADRs (especialmente **Capa 2: Conciliación**):

- ADR-034: Flujo de conciliación (pipeline invariante)
- ADR-035: Modelo de balances (qué es "correcto")
- ADR-036: Missing Purchase History (operativo)
- ADR-037: Transfers (emparejar withdrawal/deposit)
- ADR-038: Duplicados (matriz de clasificación)
- ADR-039: Holdings (validación CT vs. exchange)
- ADR-040: Cost Basis / FIFO (cuándo confiar en CT)
- [más...]

---

## Cómo Usar Este Índice

1. **Por problema:** "Tengo duplicados" → ADR-014 / ADR-026 → nivel C/D
2. **Por nivel:** "Necesito fuentes oficiales" → Nivel A → ADRs 002, etc.
3. **Por ADR:** "¿Qué ADRs afectan a la auditoría?" → Niveles B/C/D
4. **Por referencia cruzada:** "¿Qué depende de ADR-032?" → Buscar en la sección

---

## Status

- 33 ADRs creados y formalizados
- Capas 1, 5, 6 robustas
- Capa 2 (Conciliación): ADRs 034-040 pendientes (BLOQUEANTE)
- Capas 3, 4 (Integración, Fiscalidad): parciales

Ver `docs/ADR_GAP_ANALYSIS_2026-07-05.md` para análisis de brechas.
