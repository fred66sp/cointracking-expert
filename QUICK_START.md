# Quick Start — CoinTracking Expert

**¿Primera vez aquí? Lee esto en 5 minutos.**

---

## ¿Qué es esto?

Un **auditor especializado en criptomonedas** que vive en Claude Code. Reconcilia tus operaciones contra datos reales (exchange, blockchain, banco) y detecta errores antes de que afecten a tu declaración de impuestos.

**No es un asistente auxiliar.** Es un auditor que verifica tus datos y te explica paso a paso qué está mal y cómo arreglarlo.

---

## ¿Qué puedo hacer?

### `/audit-cointracking`

**"Audita mi cuenta de CoinTracking"**

- Conecta automáticamente a tu cuenta (via API)
- Reconcilia operaciones (transferencias, compras, ventas)
- Detecta: duplicados, saldos negativos, compras sin coste, transferencias huérfanas
- Genera informe detallado con cada problema y cómo solucionarlo

**Tiempo:** 10-30 min (depende de volumen)  
**Salida:** Informe en `reports/output/`

---

### `/spanish-tax-return`

**"Prepara mi declaración de la renta 2025"**

- Audita primero (datos validados)
- Calcula ganancias/pérdidas (método FIFO, España)
- Clasifica eventos imponibles
- Genera informe fiscal con cifras estimadas
- Verifica si aplica Modelo 721 (declaración bienes extranjero)

**Precondición:** Auditoría limpia (sin bloqueantes)  
**Salida:** Informe fiscal en `reports/output/`

---

## Cómo Empezar

### Paso 1: Conectar CoinTracking (una sola vez)

```
/mcp
```

Esto abre el panel de conexión MCP. Sigue las instrucciones para autorizar acceso a tu cuenta de CoinTracking.

**¿Sin MCP?** Puedes usar un CSV en lugar de conexión automática (más lento, pero funciona).

---

### Paso 2: Crear un Proyecto

Un proyecto aísla tus datos (auditoría + reportes separados). Ejemplo:

```
Tu cuenta CoinTracking
  ↓
  Proyecto "2025" ← tus operaciones 2025
  Proyecto "2024" ← tus operaciones 2024 (separado)
```

Cuando invocas un skill, te pide: "¿Qué proyecto auditar?"

**Tu primer proyecto:**
- Nombre corto: `agp` (ej. tus iniciales + año)
- Carpeta creada: `USER_INPUT/agp/`
- Reportes en: `reports/output/agp/`

---

### Paso 3: Auditar

```
/audit-cointracking
```

Responde preguntas:
- "¿Qué proyecto?" → `agp`
- "¿Quieres CSV para validación cruzada?" → Sí/No

Luego espera (10-30 min). El informe se guarda en `reports/output/agp/`.

---

### Paso 4 (Opcional): Declaración Fiscal

```
/spanish-tax-return
```

Reutiliza la auditoría que acabas de hacer. Genera informe fiscal.

---

## Conceptos Clave (2 minutos)

### Base de Coste

**¿Qué es?** Lo que pagaste por cada cripto. Necesario para calcular ganancias/pérdidas.

**¿Por qué importa?** Si falta base de coste, las ganancias salen falsas. El auditor verifica que exista.

**Ejemplo:**
```
Compré 1 BTC por 30.000€ (base = 30.000€)
Vendí 1 BTC por 40.000€ (ganancia = 10.000€)
```

---

### Método FIFO

**¿Qué es?** Regla de cálculo de ganancias en España. "First In, First Out" = la primera cripto que entró es la primera que salió.

**¿Por qué?** Porque es lo que la ley española requiere (no hay opción).

**El auditor lo hace automáticamente.**

---

### Transacciones Imponibles

**Tributan:**
- Vender cripto por EUR
- Intercambiar cripto por cripto (ej. BTC → ETH)
- Gastar cripto en compras

**No tributan:**
- Comprar (es inversión, no ingreso)
- Transferir entre tus propias cuentas (es tu dinero moviéndose)
- Holding (simplemente esperar)

---

## Troubleshooting Rápido

### "No encuentro mis reportes"

Mira en:
```
reports/output/<tu-proyecto>/
```

Ejemplo: `reports/output/agp/2026-07-05_auditoria_...md`

---

### "Dice que hay un duplicado pero sé que es legítimo"

El auditor lo lista y pide confirmación antes de eliminar. Verifica el **Trade ID** en tu exchange (Binance app, etc.). Si los Trade IDs son distintos → no es duplicado.

Ver: `knowledge/cointracking/ADR-014` para regla completa.

---

### "¿Cuánto cuesta todo esto?"

- Auditoría: ~6.000 tokens (caché reutiliza, más barato si repites)
- Declaración: ~1.300 tokens
- **Total:** ~7.300 tokens/ciclo

Con caché: **75% ahorro** en ciclos iterativos (audita → corrige → re-audita).

---

### "¿Mis datos son privados?"

Sí. Todo se almacena localmente:
```
USER_INPUT/         ← TÚ lo dejas aquí (ignorado por git)
reports/output/     ← Reportes generados (ignorado por git)
.cache/             ← Datos cacheados (ignorado por git)
```

Nada se sube a internet. El MCP solo accede a tu cuenta si das permiso.

---

## Siguiente Paso

1. **Conecta MCP** (`/mcp`)
2. **Crea tu proyecto** (responde a las preguntas)
3. **Audita** (`/audit-cointracking`)
4. **Lee el informe** (está en `reports/output/`)

**¿Problemas?** Pregunta. El auditor explica cada hallazgo.

---

## Conceptos Avanzados (Si tienes tiempo)

| Concepto | Para qué | Dónde |
|----------|----------|-------|
| ADRs (Architectural Decision Records) | Decisiones del sistema | `adr/` |
| Niveles de Conocimiento (A-F) | Cómo está organizada la base de datos | `knowledge/INDEX.md` |
| CacheManager | Por qué es tan rápido | `tools/cache_manager.py` |
| Pre-commit Hooks | Validación automática | `.git/hooks/` |

Pero para empezar: **no los necesitas.**

---

**¿Listo?** Empieza: `/audit-cointracking`

**¿Dudas?** El auditor está aquí para explicar. No hay pregunta tonta.

---

*Versión:* 1.0  
*Última actualización:* 2026-07-05
