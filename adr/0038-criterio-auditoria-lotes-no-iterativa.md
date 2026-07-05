---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-038: Criterio de Auditoría (Lotes, No Iterativa)

**Status:** Accepted  
**Proposed:** 2026-07-05  
**Accepted:** 2026-07-05  
**Last Updated:** 2026-07-05

---

## Problema

Las auditorías iterativas (audita → encuentra problema → arregla → audita nuevamente) son síntoma de:
- Falta de visibilidad del estado real del sistema
- Problemas se encuentran en múltiples pasadas, no de una vez
- No hay "punto de confianza" donde podamos afirmar "está limpio"
- Cada auditoría genera más trabajo que la anterior

**Patrón observado:**
1. Auditoría 1: "Sistema 100% limpio"
2. Auditoría 2: "Encontré 35 documentos con problemas"
3. Auditoría 3: "Encontré 91 errores (pero eran falsos positivos)"
4. Auditoría 4: "10 errores reales"
5. Auditoría MEGA: "0 errores críticos"

Este patrón NO es sostenible. Es síntoma de proceso defectuoso.

---

## Decisión

**Las auditorías son COMPLETAS (todas las dimensiones EN UNA SOLA PASADA), no iterativas.**

### Definición: Auditoría Completa

Una auditoría MEGA valida TODO simultáneamente:
1. **YAML:** Frontmatter válido en todos los documentos
2. **IDs:** Únicos, siguen convención (ADR-036)
3. **Vigencia:** valid_until nunca null (A/B), definido para C-F
4. **Campos:** Obligatorios presentes
5. **Referencias:** ADRs, links, coherencia
6. **Coherencia:** Documentos relacionados consistentes
7. **Confidence:** Valores válidos
8. **Authority:** Valores válidos
9. **Fechas:** Formatos correctos, valid_from ≤ valid_until

**Output:** UNA LISTA DE TODOS LOS PROBLEMAS (no progresiva)

### Nunca Afirmar "100% Limpio" Hasta:

1. ✅ Auditoría MEGA encuentra 0 errores críticos
2. ✅ Se arreglan los 0 errores encontrados
3. ✅ **Auditoría MEGA #2 (re-ejecuta) encuentra 0 nuevos errores**
4. ✅ **Auditoría MEGA #3 (re-ejecuta) encuentra 0 nuevos errores**

**Criterio:** 3 auditorías consecutivas sin hallazgos = realmente limpio

### Nunca Decir:

| ❌ Incorrecto | ✅ Correcto |
|---|---|
| "Encontré problema X" | "Encontré 14 problemas totales" |
| "100% limpio (pero..." | "0 errores críticos tras arreglar 14 problemas" |
| "Listo para producción" | "Listo tras 3 auditorías limpias" |

---

## Implicaciones

| Aspecto | Cambio |
|--------|--------|
| **Auditorías** | Completas, no parciales (todas las dimensiones) |
| **Frecuencia** | Máximo 1 por sesión (para no gastar tokens) |
| **Criterio de cierre** | 3 auditorías limpias, no "encontramos 0 problemas" |
| **Visibilidad** | Reporte completo UP FRONT (no descubrimientos graduales) |
| **Confianza** | Real solo después de 3 pasadas limpias |

---

## Ejemplo: Uso Correcto

**Sesión 1:**
```
MEGA Audit #1 → Encontrados 10 problemas
→ Arreglar los 10
MEGA Audit #2 → 0 problemas encontrados
→ Fin sesión
```

**Sesión 2 (próxima):**
```
MEGA Audit #3 → 0 problemas encontrados
→ Ahora sí: "Sistema 100% limpio tras 3 auditorías limpias"
```

**Lo que NO hacer:**
```
Auditar #1 → Limpiar → Auditar #2 → Limpiar → ...
(infinito: "siempre encuentra algo más")
```

---

## Referencias

- ADR-036: Convención de IDs
- ADR-037: Validación obligatoria en desarrollo
- `tools/audit_mega_complete.py`: Implementa auditoría MEGA

---

**Decisión:** ACEPTADA (previene síndrome de "siempre hay algo más")
