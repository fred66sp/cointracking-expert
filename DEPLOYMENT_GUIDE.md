# Guía de Deployment — Servidor MCP de CoinTracking

**Documento:** Cómo compilar, arrancar y operar el servidor MCP  
**Audiencia:** Desarrolladores, DevOps  
**Última actualización:** 2026-07-05

---

## 🎯 Propósito

El servidor MCP (`cointracking-mcp/`) es un **proxy Go de solo lectura** a la API de CoinTracking. Expone herramientas (`cointracking_get_trades`, `cointracking_get_balance`, etc.) a Claude Code vía MCP protocol.

**No es:**
- Un SDK de CoinTracking
- Una base de datos de historiales
- Un almacenamiento de datos (temporal: caché en memory + SQLite)

**Es:**
- Un servidor que traduce llamadas de Claude a la API de CT
- Cachea resultados por proyecto (memory + disco)
- Maneja rate limiting (60 ó 20 llamadas/hora según plan)

---

## 📋 Requisitos Previos

### Software
- **Go 1.21+** (para compilar)
- **Claude Code** (IDE + MCP support)
- **PowerShell** o **Bash** (terminal)

### Credenciales de CoinTracking
- Cuenta CoinTracking con API habilitada
- API Key (obtener en Account → API)
- API Secret (guardar en lugar seguro)

### Configuración del Sistema
- Variables de entorno persistentes (Windows: Propiedades del Sistema)
- `.mcp.json` y `.vscode/mcp.json` en el repo (ya existen)

---

## 1️⃣ Compilar el Servidor (Primera Vez)

### Paso 1: Navegar a la carpeta

```bash
cd cointracking-mcp
```

### Paso 2: Compilar

**Windows (PowerShell):**
```powershell
go build -o dist/cointracking-mcp.exe ./cmd/cointracking-mcp
```

**macOS/Linux (Bash):**
```bash
go build -o dist/cointracking-mcp ./cmd/cointracking-mcp
```

### Paso 3: Verificar

```bash
./dist/cointracking-mcp --version
# Debería mostrar: cointracking-mcp version X.Y.Z
```

### ⚠️ Notas

- **Primera compilación:** puede tomar 30-60 seg (descarga dependencias)
- **Compilaciones futuras:** 5-10 seg
- **El binario NO se versiona** (está en `.gitignore`). Después de clonar, hay que compilar.
- **Cross-compile:** si necesitas binary para otra OS, usa `GOOS=linux GOARCH=amd64 go build ...`

---

## 2️⃣ Configurar Credenciales

### Opción A: Variables de Entorno Persistentes (Recomendado)

**Windows:**
1. Panel de Control → Sistema → Propiedades Avanzadas → Variables de Entorno
2. Añade dos variables de usuario (NOT sistema):
   - `COINTRACKING_API_KEY` = `<tu_key>`
   - `COINTRACKING_API_SECRET` = `<tu_secret>`
3. Cierra y reabre Claude Code (hereda las variables)

**macOS/Linux:**
```bash
export COINTRACKING_API_KEY="<tu_key>"
export COINTRACKING_API_SECRET="<tu_secret>"
# Persistente: añade a ~/.bash_profile o ~/.zprofile
```

### Opción B: Flags al Arrancar (Menos Seguro)

```bash
./dist/cointracking-mcp --api-key <key> --api-secret <secret>
```

### ⚠️ Seguridad

- **NUNCA commitear credenciales** a git
- **NUNCA ponerlas en `.mcp.json`** (se versiona)
- Usar SOLO variables de entorno o flags (`--env-file` puede funcionar)

---

## 3️⃣ Arrancar el Servidor

El servidor se arranca automáticamente cuando Claude Code carga `.mcp.json`:

```json
{
  "name": "cointracking-mcp",
  "command": "cointracking-mcp/dist/cointracking-mcp.exe",  // Windows
  // "command": "cointracking-mcp/dist/cointracking-mcp",   // macOS/Linux
  "args": ["--tier", "pro"],  // o "unlimited" si tienes plan
  "disabled": false
}
```

### ✅ El servidor está corriendo si:

1. **En Claude Code:** Aparecen las herramientas `cointracking_*` en el MCP menu
2. **En terminal:** Puedes ver logs con `--debug` flag

```bash
./dist/cointracking-mcp --debug
# Muestra: [INFO] Server starting on port 7777
#          [INFO] Cache loaded: /path/to/.cache/cointracking/agp/
```

### ❌ Si no arranca:

```bash
# Prueba manual
./dist/cointracking-mcp --tier pro --debug

# Errores comunes:
# [ERROR] Missing COINTRACKING_API_KEY → Configura variables de entorno
# [ERROR] port 7777 already in use → Mata el proceso anterior: lsof -i :7777
# [ERROR] cache dir not found → Crea: mkdir -p .cache/cointracking/
```

---

## 4️⃣ Verificar Que Funciona

### En Claude Code

Abre un chat y llama a una herramienta:

```
cointracking_get_balance()
```

**Respuesta esperada:**
```json
{
  "balance": [
    {"currency": "BTC", "amount": 1.234, "btc_value": 1.234},
    {"currency": "ETH", "amount": 10.5, "btc_value": 0.456}
  ]
}
```

### En Terminal (si tienes `curl`)

```bash
curl http://localhost:7777/api/balance \
  -H "X-API-Key: <tu_key>" \
  -H "X-API-Secret: <tu_secret>"
```

---

## 5️⃣ Rate Limiting

### Límites por Plan

| Plan | Llamadas/Hora | TTL Caché |
|------|---------------|-----------|
| **Free** | 20 | 30 min |
| **Pro/Expert** | 60 | Configurable |
| **Unlimited** | Ninguno (bucket token) | — |

### Cómo Verificar Consumo

```
cointracking_cache_stats()
```

**Respuesta:**
```json
{
  "cache_hits": 45,
  "cache_misses": 5,
  "api_calls_this_hour": 12,
  "api_calls_limit": 60,
  "remaining": 48
}
```

### Estrategia de Caché

1. **Después de cambios en CT:** `cointracking_invalidate_cache` antes de re-consultar
2. **Consultas grandes:** Usa agregados (`get_grouped_balance`) antes que `get_trades` completo
3. **Monitoreo:** Llama a `cache_stats` antes de operaciones costosas

---

## 6️⃣ Mantenimiento Rutinario

### Diario

- Verifica que el servidor arranca al abrir Claude Code
- Si ves error `HTTP 429` (rate limit), espera 10 min antes de reintentar

### Semanal

```bash
# Limpia caché antiguo (>7 días)
rm .cache/cointracking/*/cached_*.json
```

### Mensual

```bash
# Actualiza dependencias Go
cd cointracking-mcp
go get -u ./...
go mod tidy
```

### Anual (o después de cambios en la API de CT)

- Revisar [cointracking/official/MCP_API.md](knowledge/cointracking/MCP_API.md) para nuevas herramientas
- Compilar versión nueva si hay cambios
- Actualizar `--tier` si cambia tu plan

---

## 7️⃣ Troubleshooting

### "Herramientas cointracking_* no aparecen"

1. ¿Está el servidor corriendo?
   ```bash
   netstat -ano | findstr :7777  # Windows
   lsof -i :7777                  # macOS/Linux
   ```

2. ¿Es la ruta de `.mcp.json` correcta?
   ```bash
   cat .mcp.json | grep "command"
   # Debe ser ruta relativa: cointracking-mcp/dist/cointracking-mcp.exe
   ```

3. ¿Credenciales son válidas?
   ```bash
   echo $COINTRACKING_API_KEY  # Windows: echo %COINTRACKING_API_KEY%
   # Si está vacío: configura variables de entorno
   ```

### "HTTP 429 (Rate Limit)"

```
[ERROR] HTTP 429: Too Many Requests
```

**Causa:** Consumiste tus 20/60 llamadas de la hora

**Solución:**
1. Espera 60 minutos (se resetea)
2. Upgrade a plan `unlimited` si auditas cuentas grandes
3. Usa `get_grouped_balance` en lugar de `get_trades` completo

### "Cache corrupta"

```bash
# Borra caché completa
rm -rf .cache/cointracking/

# Reinicia servidor
# Próxima llamada recrea caché
```

### "Errores de conexión (timeout)"

```
[ERROR] dial tcp: i/o timeout
```

**Posibles causas:**
- Firewall bloqueando puerto 7777
- API de CoinTracking caída (comprueba en status.cointracking.info)
- Problema de DNS (intenta con IP directa si es posible)

---

## 8️⃣ Logs y Debugging

### Ver Logs en Tiempo Real

**Opción 1:** Arrancar con `--debug`
```bash
./dist/cointracking-mcp --debug 2>&1 | tee server.log
```

**Opción 2:** Inspeccionar archivo de log
```bash
tail -f .cache/cointracking/server.log
```

### Niveles de Log

- `[DEBUG]` — Detalles internos (caché hit/miss, request/response)
- `[INFO]` — Eventos importantes (servidor arrancó, caché invalidada)
- `[WARN]` — Potenciales problemas (rate limit alto, caché vencida)
- `[ERROR]` — Fallos (API error, credenciales inválidas)

### Ejemplo de Troubleshooting

```bash
$ ./dist/cointracking-mcp --debug
[INFO] Starting cointracking-mcp v1.0.0
[INFO] Loading credentials from env
[DEBUG] API Key loaded (length: 32)
[DEBUG] API Secret loaded (length: 64)
[INFO] Cache dir: .cache/cointracking/agp/
[DEBUG] Cache tables: trades (524 rows), balance (8 rows)
[INFO] Server listening on :7777
```

---

## 🔄 Actualizar el Servidor

### Si Cambió el Código (cointracking-mcp/)

```bash
cd cointracking-mcp
git pull origin main
go mod tidy
go build -o dist/cointracking-mcp.exe ./cmd/cointracking-mcp
```

Reinicia Claude Code.

### Si Cambió la API de CoinTracking

1. Revisa [knowledge/cointracking/MCP_API.md](knowledge/cointracking/MCP_API.md)
2. Updatea `cointracking-mcp/internal/tools/*.go` si hay cambios
3. Recompila y testea

---

## 📊 Monitoreo en Producción

Si la auditoría es de una cuenta muy grande (100K+ operaciones):

### Antes de Lanzar

```bash
cointracking_cache_stats()  # ¿Cuántas llamadas quedan?
cointracking_get_trades(limit: 1)  # ¿API responde?
```

### Durante la Auditoría

```bash
# Cada 10 min
cointracking_cache_stats()
# ¿Consumo es proporcional al esperado?
```

### Después

```bash
cointracking_close_project(project_name: "agp")
# Persiste caché a disco, libera memoria
```

---

## 🔐 Seguridad

### Credenciales

- ✅ Variables de entorno
- ❌ Flags en comando (visible en `ps aux`)
- ❌ Hardcodeadas en código
- ❌ Gitignored pero sin encriptación

### Caché

- La caché local (`.cache/cointracking/`) contiene datos **en claro**
- Si otro usuario accede a tu máquina: puede ver la caché
- **Solución:** Usa permisos de carpeta restrictivos (`chmod 700` en Unix)

### Auditoría

- El servidor no registra credenciales en logs
- Los datos de CoinTracking no se suben a ningún lado
- Cada proyecto tiene su propia caché aislada (multitenancy segura)

---

## 📞 Referencia Rápida

| Necesito | Comando |
|----------|---------|
| Compilar | `go build -o dist/cointracking-mcp.exe ./cmd/cointracking-mcp` |
| Arrancar | Automático al abrir Claude Code |
| Verificar que corre | `netstat -ano \| findstr :7777` |
| Ver estadísticas | `cointracking_cache_stats()` |
| Invalidar caché | `cointracking_invalidate_cache(pattern: "*")` |
| Cerrar sesión | `cointracking_close_project(project_name: "agp")` |
| Actualizar | `git pull && go build ...` |
| Debugar | `./dist/cointracking-mcp --debug` |

---

## 🚪 Siguiente

- [knowledge/KNOWLEDGE_MAINTENANCE.md](knowledge/KNOWLEDGE_MAINTENANCE.md) — Cómo mantener la base de conocimiento
- [GOVERNANCE_WORKFLOW.md](GOVERNANCE_WORKFLOW.md) — Cómo registrar decisiones (ADRs)
