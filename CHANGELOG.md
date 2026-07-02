# Registro de cambios

**Historial de lanzamientos de CoinTracking Expert**

Todos los cambios notables en el proyecto CoinTracking Expert se documentan en este archivo. Este proyecto sigue [Versionado Semántico](https://semver.org/): MAYOR.MENOR.PARCHE para números de versión.

## [No lanzado]

### Agregado
- Agente auditor de CoinTracking en Claude Code (subagente + skill `/audit-cointracking`)
- Base de conocimiento: formato CSV, modelo de coste, integración MCP y fiscalidad española (IRPF)
- Integración con la API de CoinTracking vía MCP (solo lectura)
- Registro de decisiones (ADR-001…007)

### Cambiado
- Giro de alcance: de framework/SDK de motores deterministas a agente de IA (ADR-006)

### Eliminado
- Andamiaje del SDK descartado: paquetes Python vacíos, specs de motores, CI de pytest y documentos de la visión de framework (ADR-007)

### Cambiado
- N/A

### Deprecado
- N/A

### Eliminado
- N/A

### Arreglado
- N/A

### Seguridad
- N/A

---

## Formato de versión

Cada lanzamiento incluye:
- Número de versión (versionado semántico)
- Fecha de lanzamiento
- Características agregadas (nueva funcionalidad)
- Características cambiadas (modificaciones a funcionalidad existente)
- Características deprecadas (a ser eliminadas en versiones futuras)
- Características eliminadas (funcionalidad previamente deprecada)
- Bugs arreglados y problemas
- Actualizaciones de seguridad

---

## Lanzamientos futuros

La dirección del proyecto se registra en las decisiones arquitectónicas de [DECISIONS.md](DECISIONS.md).
