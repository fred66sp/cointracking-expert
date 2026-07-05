# Mapa de Navegación — Busca lo que Necesitas

**¿Necesitas algo? Empieza aquí, no en README.md.**

---

## 🎯 Busco por Necesidad

### "Quiero auditar mi cuenta"

```
1. Lee: QUICK_START.md (2 min)
2. Invoca: /audit-cointracking
3. Salida: reports/output/<proyecto>/AAAA-MM-DD_auditoria_...md
```

**Documentación profunda:**
- `knowledge/cointracking/` — Formato CSV, modelo de coste, importación
- `adr/0037-validacion-obligatoria-en-desarrollo.md` — Cómo valida el sistema
- `tools/ct_audit.py` — Script de chequeos deterministas

---

### "Quiero preparar mi declaración de la renta"

```
1. Lee: QUICK_START.md (2 min)
2. Invoca: /spanish-tax-return
3. Salida: reports/output/<proyecto>/AAAA-MM-DD_declaracion_...md
```

**Documentación profunda:**
- `knowledge/taxation/spain/` — Fiscalidad IRPF (ganancias, rendimientos)
- `knowledge/cointracking/CAPITAL_GAINS.md` — Método FIFO
- `templates/TAX_SUMMARY_ES.md` — Plantilla de informe

---

### "Tengo un problema específico con una operación"

**Busca por síntoma:**
- `knowledge/cointracking/TROUBLESHOOTING_INDEX.md` — Árbol de decisión
  - "Mi balance es negativo"
  - "Tengo un duplicado"
  - "Falta coste en una venta"
  - etc.

---

### "¿Cómo optimizar tokens/coste?"

```
✅ Ya está hecho (ADR-039 implementado)

Si quieres entender cómo funciona:
  1. Lee: adr/0039-optimizacion-tokens-y-cache.md (principios)
  2. Ve: docs/performance/TOKEN_BENCHMARKS.md (cifras)
  3. Ve: implementation/CACHE_ROADMAP.md (roadmap)

Si quieres usar la caché en tu código:
  1. Lee: tools/cache_manager.py (documentación)
  2. Ejemplo: .claude/skills/audit-cointracking/SKILL.md (Paso 0)
```

---

### "¿Cómo funciona el sistema?"

**Arquitectura:**
- `README.md` — Qué es, cómo funciona en general
- `adr/` — Decisiones arquitectónicas (lee índice)
- `docs/` — Documentación operativa

**De lo más importante a menos:**
1. `adr/0037-validacion-obligatoria-en-desarrollo.md` — Gobernanza
2. `adr/0036-convencion-de-ids-de-documentos.md` — Organización
3. `adr/0039-optimizacion-tokens-y-cache.md` — Eficiencia
4. `adr/0038-criterio-auditoria-lotes-no-iterativa.md` — Metodología

---

### "Tengo que migrar de un exchange a otro"

```
Mi cuenta: Binance → Coinbase

1. Lee: knowledge/reference/context/EXCHANGE_REGULATORY_UPDATES_2026.md
   (entiende qué está pasando regulatoriamente)

2. Lee: knowledge/cointracking/AUDIT_EXCHANGE_MIGRATION.md
   (procedimiento paso a paso)

3. Cuando migración esté completa:
   - Crea nuevo proyecto: agp2026
   - Invoca: /audit-cointracking
   - El auditor verifica que retiradas/depósitos estén emparejados
```

---

### "Quiero entender el código"

**Skills (lo que invocas):**
- `.claude/skills/audit-cointracking/SKILL.md` — Playbook de auditoría
- `.claude/skills/spanish-tax-return/SKILL.md` — Playbook de declaración
- `.claude/agents/cointracking-auditor.md` — Rol del subagente

**Tools (scripts que usan los skills):**
- `tools/ct_audit.py` — Validación determinista
- `tools/cache_manager.py` — Caché persistente
- `tools/benchmark_skills.py` — Test de rendimiento

**Configuración:**
- `.mcp.json` — Arranque del servidor MCP
- `.git/hooks/pre-commit` — Validación pre-commit
- `.github/workflows/audit-mega.yml` — CI/CD remoto

---

## 🗂️ Busco por Carpeta

### `knowledge/` — Base de Conocimiento (La Verdad)

```
knowledge/
├── cointracking/          ← Todo sobre CoinTracking
│   ├── CSV_FORMAT.md      ← Cómo interpreta CoinTracking
│   ├── COST_BASIS.md      ← Cómo calcula base de coste
│   └── AUDIT_EXCHANGE_MIGRATION.md ← Procedimiento (NUEVO)
├── taxation/spain/        ← Fiscalidad española
│   ├── CAPITAL_GAINS.md   ← Ganancias patrimoniales (FIFO)
│   ├── CAPITAL_INCOME.md  ← Staking, airdrops, intereses
│   └── INFORMATIVE_OBLIGATIONS.md ← Modelo 721
├── exchanges/             ← Detalles por exchange
│   ├── INDEX.md           ← Listado
│   └── official/BINANCE.md ← Particularidades Binance
└── reference/context/     ← Contexto regulatorio
    ├── BINANCE_EU_MICA_EXIT.md
    └── EXCHANGE_REGULATORY_UPDATES_2026.md (NUEVO)
```

**Regla:** Si algo no está aquí, no es cierto oficialmente. Ask before claiming.

---

### `adr/` — Decisiones Arquitectónicas

```
adr/
├── 0036-convencion-de-ids.md           ← Nombres únicos
├── 0037-validacion-obligatoria.md      ← Gobernanza
├── 0038-criterio-auditoria.md          ← Metodología
├── 0039-optimizacion-tokens-y-cache.md ← Eficiencia (REFACTORIZADO)
└── README.md → índice completo
```

**Leer cuando:** Necesites entender por qué algo se hace así.

---

### `.claude/skills/` — Puntos de Entrada

```
.claude/skills/
├── audit-cointracking/       ← /audit-cointracking
│   └── SKILL.md (playbook)
└── spanish-tax-return/       ← /spanish-tax-return
    └── SKILL.md (playbook)
```

**Invoca desde chat:** `/audit-cointracking` o `/spanish-tax-return`

---

### `tools/` — Scripts Deterministas

```
tools/
├── ct_audit.py               ← Chequeos de validación
├── cache_manager.py          ← Gestor de caché
├── benchmark_skills.py       ← Test de rendimiento
└── test_cache_savings.py     ← Validación de ahorros
```

**Úsalos:** No los invokes directamente; los skills los usan.

---

### `docs/` — Documentación Operativa

```
docs/
├── QUICK_START.md            ← EMPIEZA AQUÍ (nuevo)
├── NAVIGATION_MAP.md         ← Este archivo
├── performance/
│   └── TOKEN_BENCHMARKS.md   ← Cifras concretas (nuevo)
└── ...
```

---

### `reports/output/<proyecto>/` — Tus Reportes

```
reports/output/agp2025/
├── AAAA-MM-DD_auditoria_...md      ← Auditoría
├── AAAA-MM-DD_declaracion_...md    ← IRPF
└── REGISTRO-CAMBIOS.md             ← Log de cambios aplicados
```

**Permisos:** Ignorado por git (datos reales, privado).

---

### `USER_INPUT/<proyecto>/` — Tu Espacio

```
USER_INPUT/agp2025/
├── trades.csv                      ← CSV que descargaste de CoinTracking
└── ... (otros archivos que dejes)
```

**Permisos:** Ignorado por git (datos reales, privado).

---

## 📊 Matriz: "Quiero hacer X"

| Quiero... | Invoca | Lee primero | Salida |
|----------|--------|-------------|--------|
| Auditar | `/audit-cointracking` | QUICK_START.md | report/output/.../auditoria_... |
| Declarar | `/spanish-tax-return` | QUICK_START.md | report/output/.../declaracion_... |
| Entender arquitectura | — | adr/README.md | Entendimiento |
| Optimizar tokens | — | adr/0039-... | Automático |
| Migrar exchange | `/audit-cointracking` | AUDIT_EXCHANGE_MIGRATION.md | Report (verificado) |
| Encontrar error | — | TROUBLESHOOTING_INDEX.md | Solución |
| Contribuir código | — | ADRs + CLAUDE.md | PR con gobernanza |

---

## 🔀 Flujos Típicos

### Auditoría Iterativa (Usuario tiene problema)

```
1. /audit-cointracking
   → Lee: Informe de hallazgos
2. Usuario corrige en CoinTracking (paso a paso guiado)
3. Confirma: "Listo"
4. /audit-cointracking (de nuevo)
   → Verifica que desaparecieron los problemas
5. Si todo OK: "Tu cuenta está limpia"
```

---

### Declaración de la Renta

```
1. /audit-cointracking
   → Asegura datos limpios
2. /spanish-tax-return
   → Prepara informe fiscal
3. Usuario revisa
4. Si hay dudas: "¿Qué significa esto?"
   → Auditor explica
5. Usuario valida con asesor (recomendado)
6. Usuario declara
```

---

## ⚡ Atajos

| Busco | Atajo |
|-------|-------|
| Base de coste | `knowledge/cointracking/COST_BASIS_AND_VALIDATION.md` |
| Método FIFO | `knowledge/taxation/spain/CAPITAL_GAINS.md` §4 |
| Modelo 721 | `knowledge/taxation/spain/INFORMATIVE_OBLIGATIONS.md` |
| Cómo cachear | `tools/cache_manager.py` (docstring) |
| Validación sistema | `adr/0037-validacion-obligatoria-en-desarrollo.md` |
| Troubleshooting | `knowledge/cointracking/TROUBLESHOOTING_INDEX.md` |
| Cambios exchange 2026 | `knowledge/reference/context/EXCHANGE_REGULATORY_UPDATES_2026.md` |

---

## 🎓 Aprende Progresivamente

**Nivel 1 (5 min):** QUICK_START.md → Invoca skill  
**Nivel 2 (30 min):** Lee reportes → Entiende hallazgos  
**Nivel 3 (1h):** Explora `knowledge/cointracking/` → Modelo interno  
**Nivel 4 (2h):** Lee `adr/` → Decisiones arquitectónicas  
**Nivel 5 (avanzado):** Lee `tools/` → Implementación  

---

**¿No encuentras algo?** Busca en el índice maestro: `adr/README.md` (el más completo).

---

*Versión:* 1.0  
*Última actualización:* 2026-07-05
