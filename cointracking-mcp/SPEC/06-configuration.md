# 06 - Configuration

## Parámetros y Defaults

**CLI flags para TODO + envvars solo para credenciales:**
- **Credenciales:** CLI flags > variables de entorno > error fatal
- **Resto de parámetros:** CLI flags únicamente (con defaults sensatos)
- **Sin `.env`**

### Credenciales (Obligatorios, sin default)

| Parámetro | Flag CLI | Envvar | Descripción |
|-----------|----------|--------|-------------|
| API Key | `--api-key` | `COINTRACKING_API_KEY` | Tu API key de CoinTracking |
| API Secret | `--api-secret` | `COINTRACKING_API_SECRET` | Tu API secret de CoinTracking |

**Precedencia (en orden):**
1. CLI flag (`--api-key xxx --api-secret yyy`)
2. Variable de entorno (`COINTRACKING_API_KEY`, `COINTRACKING_API_SECRET`)
3. Si no hay nada: **ERROR fatal al startup**

**Ejemplos:**

```bash
# Option A: CLI flags (máxima prioridad)
./cointracking-mcp --api-key abc123 --api-secret xyz789

# Option B: Variables de entorno
export COINTRACKING_API_KEY=abc123
export COINTRACKING_API_SECRET=xyz789
./cointracking-mcp

# Option C: CLI override (precedence sobre envvar)
COINTRACKING_API_KEY=old_key ./cointracking-mcp --api-key new_key
# Usa new_key

# Option D: Sin credenciales
./cointracking-mcp
# ERROR: --api-key y --api-secret requeridos o variables de entorno.
```

**Seguridad:** Credenciales nunca en logs, nunca en archivos versionados.

---

### Opcionales (Con defaults)

#### Configuración de Cuenta

| Parámetro | Flag CLI | Default | Descripción |
|-----------|----------|---------|-------------|
| Account Tier | `--tier` | `pro` | Tu plan de CoinTracking (`pro`, `expert`, `unlimited`) |
| Project Name | `--project` | `default` | Nombre del proyecto (para aislar caché) |

**Ejemplos:**
```bash
./cointracking-mcp --api-key xxx --api-secret yyy --tier unlimited --project agp
```

**Mapeo de límites:**
- `--tier pro` o `--tier expert` → 20 llamadas/hora
- `--tier unlimited` → 60 llamadas/hora

**Proyectos:**
- Cada proyecto tiene caché aislada
- Útil para auditar múltiples cuentas de CT
- Ejemplo: `--project agp` vs `--project cliente_b`
- Si no se especifica: usa `default`

---

#### Caché

| Parámetro | Flag CLI | Default | Descripción |
|-----------|----------|---------|-------------|
| Cache Max Size | `--cache-max-size` | Según tier | Máximo de entradas en caché (LRU) |
| Cache Persist Dir | `--cache-dir` | `./cache` | Ruta para persistencia en disco |
| Project Env Dir | `--project-env-dir` | *(vacío = desactivado)* | Credenciales por proyecto (ADR-040): si `{dir}/{proyecto}.env` existe, ese proyecto usa su propia cuenta de CoinTracking (`COINTRACKING_API_KEY`/`COINTRACKING_API_SECRET`, `COINTRACKING_TIER` opcional). Sin fichero → credenciales del proceso. Fichero incompleto → error, nunca fallback silencioso. El directorio no se versiona jamás. |

**Defaults de CACHE_MAX_SIZE por ACCOUNT_TIER:**
- `--tier pro` → 5000 entradas (cubre 3.500 registros)
- `--tier expert` → 50000 entradas (cubre 20–100k registros)
- `--tier unlimited` → 150000 entradas (sin límite)

**TTLs (hardcoded, sin parámetro CLI):**
- getTrades: 3600s (1h)
- getBalance: 600s (10 min)
- Históricos: 7200s (2h)
- getGains: 3600s (1h)
| `CACHE_TTL_TRADES` | `3600` | TTL de getTrades (segundos) | 60–86400 |
| `CACHE_TTL_BALANCE` | `600` | TTL de getBalance (segundos) | 60–86400 |
| `CACHE_TTL_GROUPED_BALANCE` | `600` | TTL de getGroupedBalance (segundos) | 60–86400 |
| `CACHE_TTL_SUMMARY` | `7200` | TTL de getHistoricalSummary (segundos) | 60–86400 |
| `CACHE_TTL_CURRENCY` | `7200` | TTL de getHistoricalCurrency (segundos) | 60–86400 |
| `CACHE_TTL_GAINS` | `3600` | TTL de getGains (segundos) | 60–86400 |
| `CACHE_PERSIST_DIR` | `./cache` | Ruta para persistencia de caché (SQLite) | ruta absoluta o relativa |

**Notas:**
- TTL en segundos. Defaults reflejan volatilidad: balance (10 min, más volátil), trades (1h, estable), históricos (2h, fijos)
- `CACHE_PERSIST_DIR` SIEMPRE se persiste (default: `./cache`)
- Caché por proyecto en subdirectorios: `./cache/{PROJECT_NAME}/`
- La persistencia es **automática** (sin interferencia del usuario)

---

#### Logging

| Parámetro | Flag CLI | Default | Descripción |
|-----------|----------|---------|-------------|
| Log Level | `--log-level` | `info` | Nivel de logging (`debug`, `info`, `warn`, `error`) |
| Log Format | `--log-format` | `text` | Formato de logs (`text`, `json`)

**Usos:**
- `debug` — todo detalle (caché hits/misses, parámetros, tiempos)
- `info` — operaciones normales (API calls, stats, cambios caché)
- `warn` — advertencias (aproximación a rate limit, caché expiraciones frecuentes)
- `error` — solo errores

---

#### Miscelánea

| Parámetro | Flag CLI | Default | Descripción |
|-----------|----------|---------|-------------|
| Timezone | `--timezone` | `UTC` | Zona horaria para logs (IANA tz: `Europe/Madrid`, `America/New_York`) |

---

## Mínimo para Funcionar (Auditoría de 1 sesión)

**Opción A: CLI flags (recomendado)**
```bash
./cointracking-mcp \
  --api-key your_key_here \
  --api-secret your_secret_here
```

**Opción B: Variables de entorno + CLI flags**
```bash
export COINTRACKING_API_KEY=your_key_here
export COINTRACKING_API_SECRET=your_secret_here
./cointracking-mcp
```

✅ **Suficiente.** El resto usa defaults sensatos (PROJECT_NAME=default, TIER=pro, LOG_LEVEL=info, etc.).

---

## Personalización Completa (CLI flags)

```bash
./cointracking-mcp \
  --api-key your_key_here \
  --api-secret your_secret_here \
  --project agp \
  --tier unlimited \
  --cache-max-size 50000 \
  --cache-dir ./cache \
  --log-level debug \
  --log-format json \
  --timezone Europe/Madrid
```

O si prefieres envvars solo para credenciales:
```bash
export COINTRACKING_API_KEY=your_key_here
export COINTRACKING_API_SECRET=your_secret_here

./cointracking-mcp \
  --project agp \
  --tier unlimited \
  --cache-max-size 50000 \
  --log-level debug
```

---

## Validación de Configuración

Al startup, el server:

1. ✅ Lee variables de entorno del sistema (CLI override > shell exports > defaults)
2. ✅ Valida que `COINTRACKING_API_KEY` y `COINTRACKING_API_SECRET` existan
   - Si no → error fatal, salida 1
3. ✅ Lee `PROJECT_NAME` (default: `default` si no se especifica)
   - Valida: alfanumérico + _ + - (sin espacios, sin caracteres especiales)
   - Si no válido → error fatal, salida 1
4. ✅ Valida que los valores numéricos estén en rango
   - Si no → warning, usa default
5. ✅ Valida que `ACCOUNT_TIER` sea válido
   - Si no → error fatal, salida 1
6. ✅ Crea `CACHE_PERSIST_DIR/{PROJECT_NAME}` si se especifica y no existe
   - Las cachés de diferentes proyectos se aíslan en subdirectorios
7. ✅ Carga caché en memoria del proyecto actual (si existe en disco)
8. ✅ Loguea la configuración final (ofuscando credenciales, mostrando resto)

**Ejemplo de log al startup:**
```
[2026-07-02 10:30:45] cointracking-mcp v0.1.0 starting
[2026-07-02 10:30:45] API Key: abc123***xyz789 (ofuscado)
[2026-07-02 10:30:45] Project: agp
[2026-07-02 10:30:45] Tier: pro (20 calls/hour)
[2026-07-02 10:30:45] Cache: 5000 entries, ./cache/agp
[2026-07-02 10:30:45] Log level: info
[2026-07-02 10:30:45] Loaded 42 cached entries from disk
[2026-07-02 10:30:45] Listening on stdio (MCP)
```

### Ofuscación de Credenciales en Logs

**Política:** Credenciales (API key, API secret) **siempre en logs, ofuscadas pero identificables**:

| Original | Log |
|----------|-----|
| `abc123def456xyz789uvw` | `abc123***uvw` |
| `secret_key_xyz_12345` | `secret_***45` |
| `a1b2c3d4e5` | `a1b2***e5` |

**Formato:** Mostrar primeros 6 caracteres + `***` + últimos 2–3 caracteres.

**Dónde se aplica:**
- ✅ Log al startup (credenciales parseadas)
- ✅ Log de llamadas a la API (identificar cuál credential se usó)
- ✅ Log de cualquier error relacionado con autenticación
- ✅ Cualquier log relevante donde se necesite identificar la credential

**Nunca loguear completo:**
- ❌ `[INFO] Using API key: abc123def456xyz789uvw`
- ✅ `[INFO] Using API key: abc123***uvw`

**Beneficio:** Auditoría clara (saber qué credential se usó) sin exponer el secret.

---

## Precedencia de Configuración

### Para Credenciales (SOLO para `--api-key` y `--api-secret`)

De mayor a menor prioridad:

1. **CLI flag** (`--api-key xxx --api-secret yyy`) — máxima prioridad
2. **Variable de entorno del sistema** (`COINTRACKING_API_KEY`, `COINTRACKING_API_SECRET`)
3. **Sin credenciales** → ERROR fatal

**Ejemplos:**

```bash
# 1. CLI flags (máxima prioridad)
./cointracking-mcp --api-key abc123 --api-secret xyz789

# 2. Variables de entorno del sistema
export COINTRACKING_API_KEY=abc123
export COINTRACKING_API_SECRET=xyz789
./cointracking-mcp

# 3. CLI override (precedence sobre envvar)
export COINTRACKING_API_KEY=old_key
./cointracking-mcp --api-key new_key
# Usa new_key (ignora envvar)

# 4. Sin credenciales → ERROR
./cointracking-mcp --project agp
# ERROR: --api-key y --api-secret requeridos o variable de entorno.
```

### Para Todo Lo Demás (CLI flags solamente)

```bash
./cointracking-mcp \
  --api-key xxx \
  --api-secret yyy \
  --project agp \              # Si no: default
  --tier unlimited \           # Si no: pro
  --cache-max-size 50000 \     # Si no: según tier
  --cache-dir ./cache \        # Si no: ./cache
  --log-level debug \          # Si no: info
  --timezone Europe/Madrid     # Si no: UTC
```

**Orden exacto:**
```
1. ¿CLI flag?
   ./cointracking-mcp --project agp
2. ¿Default hardcoded?
   Sí → usa default (agp → default, tier → pro, etc.)
```

---

## Ejemplos de Uso Típico

**Auditoría de 1 proyecto (simplest):**
```bash
./cointracking-mcp --api-key xxx --api-secret yyy
# Usa defaults: project=default, tier=pro, log=info, cache=./cache
```

**Múltiples proyectos:**
```bash
# Proyecto A
./cointracking-mcp --api-key xxx --api-secret yyy --project cliente_a

# Proyecto B (otra terminal)
./cointracking-mcp --api-key xxx --api-secret yyy --project cliente_b

# Cache aislada:
# ./cache/cliente_a/  ← datos del cliente A
# ./cache/cliente_b/  ← datos del cliente B
```

**Con credenciales en envvars:**
```bash
export COINTRACKING_API_KEY=xxx
export COINTRACKING_API_SECRET=yyy

./cointracking-mcp --project agp --tier unlimited
```

**Para debugging:**
```bash
./cointracking-mcp \
  --api-key xxx \
  --api-secret yyy \
  --project agp \
  --log-level debug \
  --log-format json \
  --cache-dir ./cache
```

**Una línea (CLI flags):**
```bash
./cointracking-mcp --api-key xxx --api-secret yyy --project agp --tier unlimited --log-level debug
```

**Una línea (envvar + CLI flags):**
```bash
COINTRACKING_API_KEY=xxx COINTRACKING_API_SECRET=yyy ./cointracking-mcp --project agp --tier unlimited
```

---

## Validación de Defaults (Tests)

Cada default se valida en tests:
- ✅ Server arranca solo con credenciales
- ✅ Cache funciona con max_size=1000
- ✅ TTLs sensatos (no ni muy altos ni muy bajos)
- ✅ Log levels funcionan correctamente
- ✅ Persist_dir se crea si no existe
