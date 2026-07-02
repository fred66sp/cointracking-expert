# Development Guide

**Setup and Development Workflow for CoinTracking Expert**

This guide provides instructions for setting up a development environment, running tests, building documentation, and contributing code to the CoinTracking Expert project.

## Prerequisites

- Python 3.9 or higher
- Git
- A code editor (VS Code, PyCharm, etc.)
- pip and venv for package management

## Initial Setup

1. Clone the repository: `git clone https://github.com/cointracking-expert/cointracking-expert.git`
2. Navigate to project directory: `cd cointracking-expert`
3. Create virtual environment: `python -m venv venv`
4. Activate virtual environment: `source venv/bin/activate` (Linux/Mac) or `venv\Scripts\activate` (Windows)
5. Install dependencies: `pip install -r requirements.txt`
6. Install development dependencies: `pip install -r requirements-dev.txt`

## Project Structure

- `docs/` - Documentation and guides
- `knowledge/` - Domain knowledge base
- `engines/` - Engine specifications
- `src/cointracking_expert/` - Python source code
- `tests/` - Test suites
- `schemas/` - Data models
- `examples/` - Usage examples
- `cases/` - Audit cases
- `scripts/` - Utility scripts
- `reports/` - Generated reports
- `prompts/` - AI prompt templates

## Running Tests

Execute tests with pytest: `pytest tests/`

For coverage report: `pytest --cov=src tests/`

## Building Documentation

Generate HTML documentation: `mkdocs build`

Serve documentation locally: `mkdocs serve`

## Code Standards

Follow PEP 8. Run formatter: `black src/ tests/`

Run linter: `flake8 src/ tests/`

Run type checker: `mypy src/`

## Creating a Pull Request

1. Create feature branch: `git checkout -b feature/your-feature-name`
2. Make changes and commit: `git commit -am "Description of changes"`
3. Push to remote: `git push origin feature/your-feature-name`
4. Submit pull request on GitHub
5. Address review feedback
6. Merge after approval

## Documentation Guidelines

- Use clear, concrete language
- Include examples where helpful
- Update docs when code changes
- Keep docs synchronized with implementation
