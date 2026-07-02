# Guía de desarrollo

**Configuración y flujo de trabajo de desarrollo para CoinTracking Expert**

Esta guía proporciona instrucciones para configurar un ambiente de desarrollo, ejecutar pruebas, construir documentación y contribuir código al proyecto CoinTracking Expert.

## Requisitos previos

- Python 3.9 o superior
- Git
- Un editor de código (VS Code, PyCharm, etc.)
- pip y venv para gestión de paquetes

## Configuración inicial

1. Clona el repositorio: `git clone https://github.com/cointracking-expert/cointracking-expert.git`
2. Navega al directorio del proyecto: `cd cointracking-expert`
3. Crea ambiente virtual: `python -m venv venv`
4. Activa ambiente virtual: `source venv/bin/activate` (Linux/Mac) o `venv\Scripts\activate` (Windows)
5. Instala dependencias: `pip install -r requirements.txt`
6. Instala dependencias de desarrollo: `pip install -r requirements-dev.txt`

## Estructura del proyecto

- `docs/` - Documentación y guías
- `knowledge/` - Base de conocimiento del dominio
- `engines/` - Especificaciones de motores
- `src/cointracking_expert/` - Código fuente Python
- `tests/` - Suites de pruebas
- `schemas/` - Modelos de datos
- `examples/` - Ejemplos de uso
- `cases/` - Casos de auditoría
- `scripts/` - Scripts de utilidad
- `reports/` - Reportes generados
- `prompts/` - Plantillas de prompts de IA

## Ejecutar pruebas

Ejecuta pruebas con pytest: `pytest tests/`

Para reporte de cobertura: `pytest --cov=src tests/`

## Construir documentación

Genera documentación HTML: `mkdocs build`

Sirve documentación localmente: `mkdocs serve`

## Estándares de código

Sigue PEP 8. Ejecuta formateador: `black src/ tests/`

Ejecuta linter: `flake8 src/ tests/`

Ejecuta verificador de tipos: `mypy src/`

## Crear pull request

1. Crea rama de característica: `git checkout -b feature/tu-nombre-caracteristica`
2. Realiza cambios y commit: `git commit -am "Descripción de cambios"`
3. Push a remoto: `git push origin feature/tu-nombre-caracteristica`
4. Envía pull request en GitHub
5. Aborda feedback de revisión
6. Fusiona después de aprobación

## Directrices de documentación

- Usa lenguaje claro y concreto
- Incluye ejemplos cuando sean útiles
- Actualiza docs cuando el código cambia
- Mantén docs sincronizado con implementación
