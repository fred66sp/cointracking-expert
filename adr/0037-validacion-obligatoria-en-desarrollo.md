# ADR-037: Validación Obligatoria en Desarrollo

**Status:** Accepted  
**Proposed:** 2026-07-05  
**Accepted:** 2026-07-05  
**Last Updated:** 2026-07-05

---

## Problema

Auditorías iterativas encuentran problemas nuevos cada vez, síntoma de que:
- Documentos se crean sin validación DURANTE desarrollo
- Validadores se ejecutan POST-creación (demasiado tarde)
- Cambios pueden romper el sistema sin detección inmediata
- No hay "puerta de entrada" que rechace trabajo incompleto

**Ejemplo:** 91 "errores" reportados por auditoría inicial que resultaron ser:
- Convención de IDs no documentada (no error, patrón intencional)
- YAML incompleto (error real, debería rechazarse en el commit)
- Links rotos (error real, debería detectarse antes)

---

## Decisión

**La validación es obligatoria DURANTE desarrollo, no después.**

### Nivel 1: Pre-Commit Hook (Local)

Cada developer ejecuta ANTES de `git commit`:

```bash
python tools/audit_mega_complete.py
```

Si algún error es crítico:
- ❌ **RECHAZA EL COMMIT**
- El developer **FIX antes de commitear** (no después)

### Nivel 2: CI/CD (Remoto)

Al push a rama de desarrollo:
- Ejecuta auditoría MEGA automáticamente
- Si hay errores críticos: **RECHAZA EL PUSH**
- Mostrar cuál es el error específico

### Nivel 3: Pre-Merge

Antes de mergear a `main`:
- **Debe pasar 3 auditorías sin hallazgos**
- No se mergea nada que tenga errores críticos pendientes

---

## Definición de "Hecho"

Un documento está **HECHO** solo si:

1. ✅ **Frontmatter YAML completo:**
   - `id`, `title`, `level`, `domain`, `source`, `authority`
   - `last_verified`, `valid_from`, `valid_until`
   - `confidence`, `version`

2. ✅ **valid_until NUNCA null para Nivel A/B**

3. ✅ **ID único** (no duplicados en repo)

4. ✅ **ID sigue convención** (KB-[A-F][1-9]-NNN o KB-X-NNN para ejemplos)

5. ✅ **Referencias resolubles:**
   - ADR-XXX citados existen
   - Links markdown apuntan a archivos que existen

6. ✅ **Coherencia:**
   - Documentos relacionados usan misma vigencia (valid_until)
   - No hay conflictos semánticos con docs relacionados

7. ✅ **Pasa auditoría MEGA sin errores críticos**

---

## Proceso de Creación (Con Validación)

1. **Crear documento** con template YAML completo (plantilla in `adr/0036-convencion-de-ids-de-documentos.md`)
2. **Escribir contenido**
3. **Ejecutar validador:**
   ```bash
   python tools/audit_mega_complete.py
   ```
4. **Si hay errores:** FIX
5. **Si está limpio:** `git commit`

---

## Implicaciones

| Antes | Ahora |
|-------|-------|
| ❌ Crear sin validación | ✅ Validar DURANTE creación |
| ❌ Validar POST-creación | ✅ Rechaza PRs con errores |
| ❌ Problemas en auditoría | ✅ Problemas antes de commit |
| ❌ Múltiples iteraciones de fix | ✅ Una pasada, limpio |

---

## Referencias

- ADR-036: Convención de IDs
- ADR-038: Criterio de Auditoría (auditoría completa, no iterativa)
- `tools/audit_mega_complete.py`: Validador automático

---

**Decisión:** ACEPTADA (implementación inmediata tras ADR-036)
