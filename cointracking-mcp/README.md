# cointracking-mcp (Go)

MCP server que expone la API de CoinTracking como herramientas de solo lectura (operaciones, saldos, históricos, ganancias) para Claude Code, Claude Desktop y otros clientes MCP.

**Mejoras sobre el original (JS):**
- Caché integrado en el servidor (estrategia inteligente de TTL)
- Binario Go nativo (rápido, sin Node.js)
- Determinismo verificable (validaciones, detecciones de anomalías)
- Integración profunda con el auditor de CoinTracking

## Estado

✅ **Implementado (v0.1.0)** — Fases 1–5 completas: cliente API (HMAC-SHA512), 6 tools de datos + 4 tools de control de caché/proyecto (`invalidate_cache`, `cache_stats`, `close_project`, `switch_project`), caché L1 (LRU+TTL en memoria) + L2 (SQLite en disco, aislada por proyecto), cambio de proyecto activo en caliente sin reiniciar el servidor, CLI flags, logs con credenciales ofuscadas, y tests (unitarios + integración end-to-end sobre el protocolo MCP real).

**Desviaciones respecto a las specs originales** (decisiones de implementación tomadas al portar, no huecos sin resolver):
- MCP SDK: se usa el SDK oficial `github.com/modelcontextprotocol/go-sdk` (el paquete `github.com/anthropics/go-mcp` referenciado en 06-configuration.md no existe en el proxy de Go).
- SQLite: se usa el driver puro Go `modernc.org/sqlite` en vez de `mattn/go-sqlite3` (cgo), porque la máquina de desarrollo no tiene compilador de C instalado.
- `cointracking_get_grouped_balance` sí requiere parámetros (`group` obligatorio), siguiendo el repo JS de referencia real; 02-api-mapping.md lo listaba sin parámetros, lo cual no es correcto frente al repo que se pidió portar.
- `go test -race` no se pudo ejecutar en esta máquina (requiere cgo, sin compilador de C disponible). La concurrencia se revisó manualmente: un único `sync.RWMutex` protege todo el estado de la LRU. Un test de integración (`TestCachePersistsAcrossRestart`) sí encontró y verificó la corrección de una condición de carrera real: las escrituras async a disco no se esperaban al cerrar, así que se añadió un `sync.WaitGroup` (`Store.Flush`) invocado desde `close_project` y desde el cierre del servidor.

## Especificaciones

- [01 - Overview](SPEC/01-overview.md) — qué es, por qué, arquitectura general
- [02 - API Mapping](SPEC/02-api-mapping.md) — tools, métodos de la API, parámetros
- [03 - Cache Strategy](SPEC/03-cache-strategy.md) — almacenamiento, TTL, invalidación
- [04 - Features](SPEC/04-features.md) — validaciones, agregados, anomalías
- [05 - Integration](SPEC/05-integration.md) — cómo se comunica con el agente
- [06 - Configuration](SPEC/06-configuration.md) — **parámetros, defaults, validación** ← **EMPEZAR AQUÍ**
- [Reference](SPEC/reference-js-repo.md) — notas sobre el repo JS original

## Roadmap

- [x] Especificaciones finales
- [x] Implementación base (API client, tools)
- [x] Caché integrado (L1 memoria + L2 disco)
- [ ] Validaciones deterministas (`validate_trades`, `validate_balance` — Fase 1+ de 04-features.md, no implementadas aún)
- [x] Tests (unitarios + integración; sin fixtures de `tests/fixtures/` aún)
- [x] Documentación de uso
- [ ] Release v0.1.0 (pendiente de probar contra la API real y publicar binarios)

## Desarrollo

```bash
cd cointracking-mcp

# Compilar
go build -o dist/cointracking-mcp.exe ./cmd/cointracking-mcp

# Tests
go test ./...

# Ejecutar (mínimo: credenciales por flag o por env var)
./dist/cointracking-mcp.exe --api-key xxx --api-secret yyy
# o
export COINTRACKING_API_KEY=xxx
export COINTRACKING_API_SECRET=yyy
./dist/cointracking-mcp.exe --project agp --tier unlimited
```

### Integración con Claude Code (`.mcp.json`)

```json
{
  "mcpServers": {
    "cointracking": {
      "command": "H:/cointracking-expert/cointracking-mcp/dist/cointracking-mcp.exe",
      "args": ["--project", "agp", "--tier", "unlimited"],
      "env": {
        "COINTRACKING_API_KEY": "...",
        "COINTRACKING_API_SECRET": "..."
      }
    }
  }
}
```

`--project` solo fija el proyecto **inicial**. Para moverse a otro proyecto sin reiniciar el servidor (ni tocar este fichero), el agente llama en tiempo de ejecución a la tool `cointracking_switch_project(project_name=...)` — ver [ADR-016](../DECISIONS.md#adr-016-cambio-de-proyecto-activo-en-caliente-en-el-mcp-cointracking_switch_project) y [04 - Features](SPEC/04-features.md#cambio-de-proyecto-en-caliente).
