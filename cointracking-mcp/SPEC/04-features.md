# 04 - Features

## Herramientas de Control (Fase 0)

### Cierre de Proyecto

```
Tool: cointracking_close_project
Input: {
  project_name?: string  // Si no se especifica, cierra el proyecto actual
}
Output: {
  project: string,
  cached_entries: int,
  entries_persisted: int,
  message: string
}
```

**Propósito:**
- Agente indica que termina la auditoría de un proyecto
- Server **vacía caché de memoria** del proyecto
- Server **asegura que disco está sincronizado** (flush pending writes)
- Server loguea estadísticas finales (cuántas entradas se cachearon, etc.)

**Ejemplo:**
```
Agente: "Auditoría de cliente_a completada. Cerrando proyecto."
→ cointracking_close_project(project_name="cliente_a")
← {
  "project": "cliente_a",
  "cached_entries": 1250,
  "entries_persisted": 1250,
  "message": "Proyecto cerrado. Caché persistida en ./cache/cliente_a/. LRU limpiada de memoria."
}
```

**Comportamiento:**
- Libera memoria (LRU del proyecto)
- Asegura sincronización disco (fsync)
- Loguea estadísticas
- Próxima vez que se pida proyecto_a, recarga desde disco

---

### Cambio de Proyecto en Caliente

```
Tool: cointracking_switch_project
Input: {
  project_name: string  // obligatorio; alfanumérico, "_" y "-"
}
Output: {
  project: string,
  previous_project?: string,
  already_active: bool,
  entries_cleared_previous_project: int,
  entries_loaded_from_disk: int,
  message: string
}
```

**Propósito:**
- Mueve el proyecto activo del proceso MCP ya arrancado, sin reiniciar el servidor ni editar `.mcp.json` (ver ADR-013 en `CLAUDE.md` del repo principal: la caché por proyecto existía, pero cambiar de proyecto exigía reiniciar el proceso).
- Internamente hace lo mismo que `NewApp` al arrancar: vacía y cierra la caché del proyecto saliente (igual que `cointracking_close_project`), y abre/crea la caché SQLite del proyecto entrante bajo `{cache-dir}/{project}` (SPEC 03).
- Credenciales, `--tier` y el limitador de tasa son del proceso (una cuenta de CoinTracking), no del proyecto, así que no cambian.
- Si `project_name` coincide con el proyecto ya activo, es un no-op (`already_active: true`) que no toca la caché.

**Ejemplo:**
```
Agente: "El usuario quiere trabajar ahora con el proyecto cliente_b."
→ cointracking_switch_project(project_name="cliente_b")
← {
  "project": "cliente_b",
  "previous_project": "agp",
  "already_active": false,
  "entries_cleared_previous_project": 1250,
  "entries_loaded_from_disk": 340,
  "message": "Proyecto activo cambiado de \"agp\" a \"cliente_b\". Caché aislada en ./cache/cliente_b."
}
```

**Comportamiento:**
- Llamar a esta tool en cuanto se fije el proyecto activo de la conversación, antes de cualquier `cointracking_get_*`.
- Si el proyecto ya se había visitado en el proceso, su caché en disco se recarga (sin llamadas nuevas a la API para lo que siga vigente).
- Nombre de proyecto inválido → error de validación, sin tocar la caché actual.

---

## Validaciones Deterministas (Fase 1)

El server puede ejecutar validaciones **deterministas** sobre los datos de CoinTracking, reportando inconsistencias sin necesidad de que el agente re-derive lógica:

### Sobre getTrades

```
Tool: cointracking_validate_trades
Input: {
  trades: []Trade,  // Usar getTrades como entrada
  check_duplicates: bool,
  check_orphans: bool,
  check_balances: bool
}
Output: {
  issues: [
    {
      type: "DUPLICATE",
      severity: "error" | "warning" | "info",
      trades: [id1, id2, ...],
      message: "Operación duplicada: BTC comprado 2026-07-02 en Kraken dos veces"
    },
    {
      type: "ORPHAN_TRANSFER",
      severity: "error",
      trade_id: "...",
      message: "Transferencia de entrada sin origen (depósito huérfano)"
    },
    ...
  ]
}
```

**Tipos de validación:**

| Tipo | Descripción | Determinista |
|------|-------------|--------------|
| `DUPLICATE` | Operación idéntica (divisa, cantidad, timestamp) | Sí |
| `ORPHAN_TRANSFER` | Depósito/retiro sin contraparte | Sí (si se tiene histórico completo) |
| `IMPOSSIBLE_BALANCE` | Saldo negativo en ningún momento | Sí |
| `MISSING_FEES` | Operación con comisión pero no registrada | No (requiere IA) |
| `CURRENCY_MISMATCH` | Divisa en el trade no existe en CT | Sí |

### Sobre getBalance

```
Tool: cointracking_validate_balance
Input: {
  trades: []Trade,
  current_balance: Balance,
  fiat_currency: string
}
Output: {
  reconciled: bool,
  unreconciled: [
    {
      currency: "BTC",
      expected: 0.5,
      actual: 0.48,
      difference: -0.02,
      possible_cause: "Transferencias internas no registradas"
    }
  ],
  totals: {
    value_in_fiat: 10000,
    fiat_currency: "EUR"
  }
}
```

## Agregados Precomputados (Fase 1+)

El server puede precomputar agregados útiles **desde getTrades**, sin esperar a que el agente lo haga:

```
Tool: cointracking_aggregate_trades
Input: {
  trades: []Trade,
  group_by: "currency" | "exchange" | "type" | "month",
  include_stats: bool  // Desviación estándar, mediana, etc.
}
Output: {
  groups: [
    {
      key: "BTC",
      count: 150,
      total_buy: 2.5,
      total_sell: 1.0,
      net: 1.5,
      fees: 0.01,
      exchanges: ["Kraken", "Binance"],
      types: ["buy", "trade", "transfer"]
    }
  ]
}
```

Muy útil para análisis rápido sin re-procesar.

## Detección de Anomalías (Fase 2+)

Usar patrones heurísticos para detectar problemas que no son estrictamente deterministas pero son sospechosos:

```
Tool: cointracking_detect_anomalies
Input: {
  trades: []Trade,
  sensitivity: "low" | "medium" | "high"
}
Output: {
  anomalies: [
    {
      type: "UNUSUAL_VOLUME",
      severity: "warning",
      details: "Compra de 100 BTC el 2026-07-01 (media histórica: 0.5)"
    },
    {
      type: "TIMESTAMP_CLUSTERING",
      severity: "info",
      details: "50 operaciones en el mismo segundo (posible re-import)"
    }
  ]
}
```

**Nota:** No son errores determinísticos, sino alertas para revisión manual.

## Determinismo Verificable (Arquitectura)

Todo lo que se marca como "determinista" debe:

1. **Ser reproducible:** mismo input → mismo output
2. **Ser auditable:** código simple, sin IA, con lógica explícita
3. **Tener tests:** fixtures con casos conocidos
4. **Tener trazabilidad:** cada resultado explica cómo se llegó a él

Ejemplo de resultado determinista:
```json
{
  "type": "DUPLICATE",
  "trades": ["trade_id_1", "trade_id_2"],
  "hash_1": "sha256(trade_1_canonical_form)",
  "hash_2": "sha256(trade_2_canonical_form)",
  "reason": "Hashes idénticos → duplicado exacto"
}
```

## Logging y Auditoría (Fase 0)

Cada operación se loguea:

```
[2026-07-02 10:30:45] getTrades(limit=100, start=1234567890)
  Source: agent (Claude Code)
  Result: 150 trades, 5000 bytes
  Cached: no (miss)
  API Time: 345ms
  Total Time: 352ms
```

Útil para:
- Entender qué pidió el agente (debugging)
- Estimar consumo de API
- Auditar acceso a datos sensibles

## Escalabilidad y Limitaciones (Fase 0)

**Fase 0 (actual):**
- 6 tools estándar (como el repo JS)
- Caché en memoria
- Validación básica de parámetros
- Logging a stderr

**Fase 1 (siguiente):**
- Validaciones deterministas
- Caché persistido
- Agregados precomputados
- Stats de caché

**Fase 2+ (futuro):**
- Detección de anomalías (ML/heurísticas)
- Sincronización multi-instancia (si es necesario)
- Integración con herramientas externas (auditoría, impuestos)

## Coordinación con el Agente

El agente es responsable de:
- Decidir cuándo llamar qué herramienta
- Validar que los datos tienen sentido (semantics)
- Explicar hallazgos al usuario
- Pedir correcciones en CoinTracking

El server es responsable de:
- Cachear eficientemente
- Validaciones deterministas (duplicados, balances imposibles)
- Trazabilidad total
- Rendimiento y manejo de errores de API
