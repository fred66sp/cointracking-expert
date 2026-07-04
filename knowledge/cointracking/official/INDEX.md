# Documentación Oficial de CoinTracking (Nivel A2)

**Ubicación:** `knowledge/cointracking/official/`

**Característica:** Documentos que referencian **directamente** a CoinTracking oficial, sin interpretación.

**Autoridad:** `official` — provienen de fuente verificable (CoinTracking o datos reales de usuario)

---

## Documentos

- **[CSV_FORMAT.md](CSV_FORMAT.md)** — Formato de la exportación "Trade Table" (Trade List, simple CSV, completo). Validado contra exportaciones de producción.
  - Columnas, tipos, formato de fecha
  - Variantes (ES, EN)
  - Peculiaridades por configuración
  - Validación automática

- **[COST_BASIS_AND_VALIDATION.md](COST_BASIS_AND_VALIDATION.md)** — Cómo CoinTracking calcula base de coste (purchase pool) y detecta inconsistencias.
  - Algoritmo de FIFO
  - Negativos y missing purchase history
  - Permutas y tercera moneda
  - Casos problemáticos reales

- **[MCP_API.md](MCP_API.md)** — Referencia de la API de CoinTracking (vía MCP).
  - Endpoints disponibles
  - Parámetros y límites
  - Semántica temporal
  - Peculiaridades observadas

- **[WEB_APP_GUIDE.md](WEB_APP_GUIDE.md)** — Guía operativa: cómo usar la web de CoinTracking para corregir problemas.
  - Dónde hacer clic
  - Qué escribir
  - Cómo generar el Tax Report de España (FIFO)
  - Referencias a artículos oficiales

---

## Documentos Relacionados

- Casos reales: `knowledge/cases/` — CT-002 (FLOKI), CT-003 (Missing PH), etc.
- Patrones operativos: `knowledge/cointracking/behavioral/` (próxima sesión)
- Procedimientos: `knowledge/procedures/` (próxima sesión)
- Catálogo de ayuda: `reference/CATALOG.md` — índice de 205 artículos del centro de ayuda

---

## Status de Migración (Fase 2)

- [x] Crear INDEX.md
- [x] Mover CSV_FORMAT.md
- [x] Mover COST_BASIS_AND_VALIDATION.md
- [ ] Mover MCP_API.md (próxima revisión)
- [ ] Mover WEB_APP_GUIDE.md (próxima revisión)

---

## Próximas Sesiones

**Fase 3:** Agregar metadatos YAML (`id: KB-A2-001`, etc.) a cada documento.
