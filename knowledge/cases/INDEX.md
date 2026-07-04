# Casos Reales Auditados (Nivel C1)

**Ubicación:** `knowledge/cases/`

**Característica:** Documentación de **casos reales auditados** en proyectos del usuario.

**Autoridad:** `verified` — casos reales = máxima confianza, no es teoría

---

## Estado Actual (Fase 2)

Los 20 casos se encontraban en `knowledge/patterns/cointracking_casos_v2.yaml` (esquema YAML canónico).

**En esta sesión:** Se han migrado a esta carpeta. En Fase 3 se convertirán a archivos `.md` individuales con metadatos YAML.

**Archivo:** [`cointracking_casos_v2.yaml`](cointracking_casos_v2.yaml)

---

## Casos Documentados (v2)

| ID | Título | Categoría | Gravedad |
|----|--------|-----------|----------|
| CT-001 | Duplicados (mismo timestamp) | duplicados | medio |
| CT-002 | FLOKI: 29 transacciones idénticas NO son duplicadas | duplicados | alto |
| CT-003 | Missing Purchase History: compras sin origen | missing_ph | crítico |
| ... | (17 casos más) | ... | ... |

**Resumen:** 20 casos, esquema canónico, 16 campos obligatorios por caso.

---

## Próximas Sesiones

**Fase 3:**
1. Convertir YAML → archivos `.md` individuales (`CT-001-*.md`, `CT-002-*.md`, etc.)
2. Insertar metadatos YAML en cada archivo (`id: KB-C1-001`, etc.)
3. Actualizar INDEX.md con referencias a archivos individuales

---

## Relación con otros Niveles

- **Nivel B1-B3:** Estos casos **fundamentan** el conocimiento operativo
- **Nivel C2:** Patrones se **derivan** de estos casos
- **Nivel C3:** Procedimientos se **basan en** estos casos
- **Nivel D1-D2:** Checklists y árboles de decisión se **construyen a partir de** estos casos

---

## Esquema Actual (YAML)

Cada caso tiene 16 campos obligatorios:

```yaml
- id: CT-002
  titulo: "FLOKI: 29 transacciones idénticas no son duplicadas"
  categoria: duplicados
  sintomas: [síntomas observados]
  causa_probable:
    hecho: [la causa raíz verificada]
  evidencia_minima: [qué se necesita para confirmar]
  pasos_diagnostico: [cómo diagnosticar]
  solucion_recomendada: [qué hacer]
  anti_patron: [qué NO hacer]
  por_que_falso_positivo: [por qué se confunde]
  nivel_confianza: verificado
  nivel_riesgo: alto
  impacto_fiscal_potencial: [pérdida de dinero / multa]
  senales_tempranas: [indicios iniciales]
  validacion_antes_despues: [cómo validar]
  vigencia:
    fecha_revision: 2026-07-03
    motivo_caducidad_potencial: [cuándo revalidar]
    fuente_recomendada_para_revalidar: [dónde verificar]
```

---

## Política de Vigencia

- Cada caso declara `vigencia` (fecha_revision, motivo_caducidad, fuente para revalidar)
- Antes de citar un caso, verificar que no esté desfasado
- Casos con `nivel_confianza: hipotesis` o `pendiente_verificar` requieren reverificación antes de usar

---

## Acceso Rápido

- **Por síntoma:** Ver `knowledge/cointracking/TROUBLESHOOTING.md` (índice por síntoma que remite a casos)
- **Por patrón:** Ver `knowledge/patterns/` (patrones generalizados desde casos)
- **Por categoría:** Ver tabla de arriba

---

## Próxima Migración

En Fase 3, cada caso tendrá su propio archivo:

```
knowledge/cases/
├── CT-001-duplicate-same-timestamp.md
├── CT-002-floki-batching.md
├── CT-003-missing-purchase-history.md
└── ... (CT-004 a CT-020)
```

Cada archivo tendrá:
- Metadatos YAML (`id: KB-C1-002`, etc.)
- Contenido del caso (antes en YAML, ahora markdown)
- Enlaces a ADRs relacionados
- Enlaces a patrones derivados
