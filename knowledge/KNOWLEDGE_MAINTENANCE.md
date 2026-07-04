# Mantenimiento de la Base de Conocimiento

**Documento:** Cómo agregar, actualizar y deprecar documentos de conocimiento  
**Audiencia:** Desarrolladores, especialistas de criptomonedas  
**Última actualización:** 2026-07-05

---

## 🎯 Propósito

La base de conocimiento es el **cerebro auditor**. Debe mantenerse:

- ✅ **Íntegra:** Sin documentos rotos o referencias perdidas
- ✅ **Vigente:** Metadatos actualizados, nada expirado sin aviso
- ✅ **Trazable:** Cada cambio justificado en un commit + ADR si es importante
- ✅ **Navegable:** Índices actualizados, sin brechas

---

## 📂 Estructura de Carpetas

```
knowledge/
├── cointracking/
│   ├── official/           ← NIVEL A (fuentes oficiales)
│   │   ├── CSV_FORMAT.md
│   │   ├── COST_BASIS_AND_VALIDATION.md
│   │   └── INDEX.md
│   ├── behavioral/         ← NIVEL B (operativo)
│   │   ├── STAKING_MECHANICS.md
│   │   ├── PURCHASE_POOL_MECHANICS.md
│   │   └── ...
│   └── reference/          ← Referencias/catálogos
│       ├── CATALOG.md
│       └── INDEX.md
├── exchanges/
│   ├── official/           ← NIVEL A
│   └── behavioral/         ← NIVEL B
├── blockchains/            ← NIVEL B
├── taxation/spain/         ← NIVEL A (fiscal española)
├── cases/                  ← NIVEL C1 (casos verificados)
│   ├── ct-001-*.md
│   └── INDEX.md
├── patterns/               ← NIVEL C2 (patrones)
├── procedures/             ← NIVEL C3 (procedimientos)
├── checklists/             ← NIVEL D (auxiliar)
├── decision-trees/         ← NIVEL D (auxiliar)
├── reference/              ← NIVEL E (glosario, contexto)
├── INDEX_MASTER.md         ← Mapa maestro
├── QUICK_START.md          ← Navegación (nuevo)
├── NAVIGATION_MAP.md       ← Navegación (nuevo)
├── TROUBLESHOOTING_INDEX.md ← Navegación (nuevo)
└── CHEAT_SHEET.md          ← Navegación (nuevo)
```

---

## 1️⃣ Crear un Documento Nuevo

### Paso 1: Decidir el Nivel

| Nivel | Cuándo | Ejemplo |
|-------|--------|---------|
| **A** | Fuente oficial (AEAT, CoinTracking, exchange) | CSV_FORMAT.md |
| **B** | Cómo funciona algo (comportamiento verificado) | STAKING_MECHANICS.md |
| **C** | Caso real auditado / Patrón / Procedimiento | ct-001.md, PATTERN_DUPLICATE_DETECTION.md |
| **D** | Checklist o árbol de decisión | CHECKLIST_DUPLICATES.md |
| **E** | Glosario, contexto, historiadores | GLOSSARY.md |
| **F** | Decisión arquitectónica (ADR) | adr/ADR-033.md |

### Paso 2: Crear el Archivo

**Formato de nombre:**
- **Nivel A-E:** `TEMA_O_TITULO.md` (snake_case, descriptivo)
- **Nivel C (casos):** `ct-NXX-descripcion-corta.md` (CT = CoinTracking case, NXX = número)
- **Nivel F (ADRs):** `adr/ADR-NXX.md` (ADR = Architectural Decision Record)

**Ejemplo:**
```bash
touch knowledge/cointracking/behavioral/EXCHANGES_CUSTOM_MAPPINGS.md
```

### Paso 3: Escribir Frontmatter YAML

Copia esta plantilla al inicio del archivo:

```yaml
---
id: "KB-B1-013"              # Único, sigue el patrón: KB-[NIVEL][SUBDOMINIO]-[NUM]
title: "Nombre descriptivo del documento"
level: "B"                   # A, B, C, D, E, o F
domain: "cointracking"       # cointracking, taxation, exchanges, blockchains, general
source: "Fuente de verdad"   # URL, documento oficial, análisis propio, etc
authority: "official"        # official (A), verified (B-D), reference (E)
last_verified: "2026-07-05"  # YYYY-MM-DD (hoy si es nuevo)
valid_from: "2024-01-01"     # Cuándo entra en vigor
valid_until: "2027-07-05"    # Cuándo expira (OBLIGATORIO para Nivel A)
confidence: "high"           # high, medium, low
version: "1.0"               # Semver (1.0 para nuevo)

related_adr:
  - ADR-033
  - ADR-010

related_docs:
  - knowledge/cointracking/CSV_FORMAT.md
  - knowledge/cases/ct-001.md

tags:
  - operativo
  - exchange
  - behavioral
---
```

### Paso 4: Escribir el Contenido

**Estructura mínima:**

```markdown
# Título Descriptivo

## Propósito (Sección Obligatoria)

Una frase: "Este documento explica cómo..."

## Contenido Técnico

(Tu contenido aquí)

## Señales de Alerta / Casos Límite

Qué puede salir mal.

## Referencias

Apunta a otros documentos relacionados.
```

### Paso 5: Validar Frontmatter

Ejecuta el validador:

```bash
python scripts/validate_knowledge_metadata.py
```

**Debe mostrar:** `[OK] 1 archivos válidos`

---

## 2️⃣ Actualizar un Documento Existente

### Si el Contenido Cambió (Pero No Es Versionado)

```yaml
last_verified: "2026-07-05"  # Hoy
version: "1.1"               # Incrementa minor
```

**Ejemplo:** Corriges una errata, clarifica una frase, añades un ejemplo.

```bash
git commit -m "Clarificar sección X en DOCUMENTO.md"
```

### Si el Contenido Cambió Significativamente (Es Versionado)

```yaml
last_verified: "2026-07-05"
version: "2.0"               # Incrementa major si es cambio importante
```

**Ejemplo:** Cambias la regla de cálculo, añades nueva sección, reescribes completamente.

```bash
git commit -m "Actualizar DOCUMENTO.md: cambio en definición de X (v2.0)"
```

### Si Expira/Envejece

```yaml
valid_until: "2026-12-31"    # Acércate a la fecha actual
# O
valid_until: null            # Si nunca expira (solo Nivel E)
```

Marca como `[PENDIENTE DE VERIFICAR]` en el documento si ya pasó la fecha.

---

## 3️⃣ Deprecar un Documento

### Paso 1: Crear Documento de Deprecación

Crea un archivo nuevo: `knowledge/.metadata/DEPRECATIONS.md`

```markdown
# Documentos Deprecados

## 2026-07-05

- **EXCHANGE_FEES_CALCULATION.md** (KB-B2-010)
  - Razón: Fusionado con FEE_HANDLING.md
  - Reemplazado por: [FEE_HANDLING.md](../cointracking/behavioral/FEE_HANDLING.md)
  - Acción: Buscar referencias y actualizar
```

### Paso 2: Reemplazar Referencias

```bash
# Busca dónde aparece el doc
grep -r "EXCHANGE_FEES_CALCULATION" knowledge/

# Actualiza los enlaces
# Antes: [EXCHANGE_FEES_CALCULATION.md](...)
# Después: [FEE_HANDLING.md](...)
```

### Paso 3: Mover o Eliminar

```bash
# Opción A: Mover a carpeta de archive
mv knowledge/cointracking/behavioral/EXCHANGE_FEES_CALCULATION.md \
   knowledge/.deprecated/EXCHANGE_FEES_CALCULATION.md

# Opción B: Eliminar si no es necesario
rm knowledge/cointracking/behavioral/EXCHANGE_FEES_CALCULATION.md
```

### Paso 4: Commit

```bash
git commit -m "Deprecar EXCHANGE_FEES_CALCULATION.md, fusionado en FEE_HANDLING.md"
```

---

## 4️⃣ Agregar un Caso Nuevo (Nivel C1)

### Paso 1: Auditar y Documentar el Caso

Durante una auditoría real, si encuentras un problema **nuevo** no cubierto:

```bash
# Número siguiente: grep -r "^id: \"KB-C1-" knowledge/cases/ | tail -1
# (Busca el número más alto, suma 1)

# Crear archivo
touch knowledge/cases/ct-021-descripcion-corta.md
```

### Paso 2: Llenar con la Evidencia

```yaml
---
id: "KB-C1-021"
title: "Caso CT-021: (Descripción del problema)"
level: "C"
domain: "cointracking"
source: "Auditoría real (proyecto: agp, fecha: 2026-07-05)"
authority: "verified"
last_verified: "2026-07-05"
valid_from: "2026-07-05"
valid_until: null
confidence: "high"
version: "1.0"

related_adr:
  - ADR-003 (si es sobre duplicados)
  - ADR-009 (si es sobre auditoría)

related_docs:
  - knowledge/cointracking/behavioral/...
---

# CT-021: (Descripción)

## Síntomas
- Descripción observable

## Causa Probable
- Deducción lógica

## Evidencia Mínima
- Qué datos confirman el caso

## Pasos de Diagnóstico
1. ...
2. ...

## Solución Recomendada
- Cómo arreglarlo

## Validación Antes/Después
- Antes: Estado inicial
- Después: Estado correcto
```

### Paso 3: Commit

```bash
git commit -m "Agregar caso CT-021: descripción (verificado en auditoría agp)"
```

---

## 5️⃣ Mantener Metadatos Vigentes

### Diaria: Revisar Alertas de Expiración

```bash
python scripts/check_knowledge_vigencia.py
```

**Salida esperada:**
```
[WARN] 3 documentos caducados o próximos a caducar:
  - IRPF_2024.md (vence 2026-04-30) [7 meses atrás]
  - BINANCE_FEES.md (vence 2026-09-30) [2 meses]
```

### Semanal: Actualizar `last_verified`

Si verificaste un documento contra una fuente oficial:

```yaml
last_verified: "2026-07-05"  # Hoy
version: "1.2"
```

**Qué significa:** "Verifiqué esto contra [fuente] el 2026-07-05 y es correcto."

### Mensual: Revisar Vigencia General

```bash
# Lista documentos que vencen en 30 días
grep -r "valid_until:" knowledge/ | grep -E "2026-08"
```

Actualiza si es necesario:

```yaml
# Antes:
valid_until: "2026-08-05"

# Después (después de verificar):
valid_until: "2027-08-05"
```

---

## 6️⃣ Validación Automática

### Ejecutar el Validador

```bash
python scripts/validate_knowledge_metadata.py
```

**Verificaciones que hace:**
- ✅ Frontmatter YAML válido
- ✅ ID único (no hay duplicados)
- ✅ Campos obligatorios presentes
- ✅ Tipos de dato correctos
- ✅ `valid_until` no es null para Nivel A

### Interpretar Errores

| Error | Solución |
|-------|----------|
| `No tiene frontmatter YAML válido` | Revisa la sintaxis YAML (espacios, comillas) |
| `ID duplicado: KB-B1-005` | Cambia el ID a uno único |
| `Campo ausente: 'confidence'` | Añade `confidence: high/medium/low` |
| `valid_until no puede ser null (Nivel A)` | Cambia a fecha específica, ej: `2027-07-05` |

---

## 7️⃣ Referencias Cruzadas

### Actualizar Índices Cuando Añades Docs

Tres lugares que necesitan actualización:

1. **INDEX_MASTER.md** — Mapa de alto nivel
2. **NAVIGATION_MAP.md** — Búsquedas por función
3. **Índices locales** — `INDEX.md` en cada carpeta

**Ejemplo:** Añadiste nuevo doc en `behavioral/`:

```bash
# 1. Edita knowledge/cointracking/behavioral/INDEX.md
# Añade entrada para el nuevo doc

# 2. Edita knowledge/INDEX_MASTER.md
# Actualiza el estado de "B2 — Operativas de Exchanges" (% completado)

# 3. Edita knowledge/NAVIGATION_MAP.md
# Añade entrada en tabla si aplica a una necesidad común

# 4. Commit
git commit -m "Actualizar índices para nuevo documento XXX.md"
```

---

## 8️⃣ Flujo de Governance

### Para Cambios Pequeños (Correcciones, Erratas)

```bash
# 1. Edita el archivo
# 2. Commit corto
git commit -m "Corregir errata en CSV_FORMAT.md"
# 3. Push a main
git push origin main
```

**No necesita ADR.**

### Para Cambios Medianos (Nuevo documento, nuevo nivel de conocimiento)

```bash
# 1. Edita / crea el archivo
# 2. Commit explicativo
git commit -m "Agregar EXCHANGES_CUSTOM_MAPPINGS.md (Nivel B2-010)

Explica cómo manejar exchanges con mapeos de símbolos personalizados.
Relacionado con ADR-003 (múltiples exchanges).
"
# 3. Push a main
git push origin main
```

**Considera si merece un ADR (↓).**

### Para Cambios Grandes (Nueva arquitectura, decisión importante)

```bash
# 1. Crea ADR en adr/ADR-NXX.md
# Sigue el template MADR (ver adr/README.md)

# 2. Commit del ADR
git commit -m "ADR-034: Introducir Nivel X de conocimiento

Propón nuevo nivel de conocimiento para casos de uso X.
Relacionado con ADR-033 (arquitectura).
"

# 3. Espera feedback (si es en equipo)

# 4. Implementa cambios en conocimiento
# 5. Commit de implementación
git commit -m "Implementar ADR-034: estructura + 5 docs iniciales"

# 6. Push a main
git push origin main
```

---

## 9️⃣ Checklist de Lanzamiento

Cuando hayas terminado un documento o lote:

- [ ] Frontmatter YAML completo y válido
- [ ] `id` es único (KB-X-YYY sin duplicados)
- [ ] `valid_until` tiene fecha si es Nivel A
- [ ] `last_verified` es hoy o reciente
- [ ] Referencias cruzadas actualizadas (INDEX.md, NAVIGATION_MAP.md)
- [ ] Validador `validate_knowledge_metadata.py` pasa sin errores
- [ ] Enlaces internos usan rutas relativas (`../` si es necesario)
- [ ] Sin emojis en YAML (usa ASCII en metadatos)
- [ ] Commit message explica qué y por qué
- [ ] Si es importante: nuevo ADR en `adr/`

---

## 🔟 Referencia Rápida

| Necesito | Comando |
|----------|---------|
| Ver estado general | `python scripts/validate_knowledge_metadata.py` |
| Buscar documentos vencidos | `python scripts/check_knowledge_vigencia.py` |
| Crear documento nuevo | `touch knowledge/.../NOMBRE.md` + copiar plantilla YAML |
| Actualizar índices | Edita `INDEX_MASTER.md`, `NAVIGATION_MAP.md`, índices locales |
| Crear caso nuevo | `touch knowledge/cases/ct-NXX.md` + copiar plantilla |
| Deprecar documento | Crear entrada en `DEPRECATIONS.md`, actualizar referencias |
| Commit | `git commit -m "Descripción clara de qué cambió"` |
| Crear ADR | `touch adr/ADR-NXX.md` + copiar plantilla MADR |

---

## 🚪 Siguiente

- [GOVERNANCE_WORKFLOW.md](GOVERNANCE_WORKFLOW.md) — Cómo registrar decisiones importantes (ADRs)
- [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) — Cómo compilar y arrancar el servidor MCP
