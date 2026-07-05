#!/usr/bin/env python3
"""
Rastreador de versiones de Knowledge Base (y, si existiera, ADRs con frontmatter).

⚠️ LIMITACIÓN REAL (descubierta 2026-07-05, revisión de robustez):
  Este tracker solo detecta un campo `version: X.X` dentro de un bloque de
  frontmatter YAML (delimitado por `---`). Los documentos de `knowledge/`
  SÍ tienen ese frontmatter y por tanto SÍ son rastreables. Los ADRs de
  `adr/` usan formato MADR plano (`**Status:** Accepted`, sin YAML), así
  que get_current_versions() devuelve 0 claves `adr_*` en este repo tal
  como está hoy. En la práctica: cambiar un ADR NO invalida ningún caché
  automáticamente; cambiar un documento de `knowledge/` (subiendo su
  `version:`) sí lo hace. No asumir lo contrario sin añadir frontmatter
  YAML a los ADRs primero.

Implementa detección automática de cambios en:
  - Knowledge Base (YAML frontmatter) — funcional
  - ADRs (Architecture Decision Records) — no funcional con el formato MADR actual
  - MCP (versión del servidor, desde .mcp.json si declara "version")

Uso:
  tracker = VersionTracker()
  current = tracker.get_current_versions()
  # -> {'kb_capital_gains': '2.1', ...}  (no incluye adr_* con el formato actual)

  # Verificar si caché es válida
  valid = tracker.is_cache_valid(cached_versions, current)
  # -> True/False
"""

import json
import hashlib
from pathlib import Path
from typing import Dict, Any
import re


class VersionTracker:
    """Rastrear versiones de ADRs y Knowledge Base."""

    def __init__(self, project_root: str = '.'):
        self.project_root = Path(project_root)
        self.adr_dir = self.project_root / 'adr'
        self.kb_dir = self.project_root / 'knowledge'

    def get_current_versions(self) -> Dict[str, str]:
        """Obtener versiones actuales de ADRs y KB."""
        versions = {}

        # Leer ADRs
        if self.adr_dir.exists():
            for adr_file in self.adr_dir.glob('*.md'):
                if adr_file.name == 'README.md':
                    continue
                version = self._extract_version(adr_file)
                if version:
                    adr_id = adr_file.stem  # ej. 0039-optimizacion-tokens-y-cache
                    versions[f'adr_{adr_id}'] = version

        # Leer Knowledge Base
        if self.kb_dir.exists():
            for kb_file in self.kb_dir.rglob('*.md'):
                if kb_file.name == 'INDEX.md' or kb_file.name == 'README.md':
                    continue
                version = self._extract_version(kb_file)
                if version:
                    # ej. knowledge/taxation/spain/CAPITAL_GAINS.md → kb_capital_gains
                    relative = kb_file.relative_to(self.kb_dir)
                    key = f"kb_{relative.stem.lower()}"
                    versions[key] = version

        # MCP version (si existe .mcp.json)
        mcp_version = self._get_mcp_version()
        if mcp_version:
            versions['mcp'] = mcp_version

        return versions

    def _extract_version(self, file_path: Path) -> str:
        """Extraer 'version: X.X' del frontmatter YAML."""
        try:
            content = file_path.read_text(encoding='utf-8')
            match = re.search(r'^---\n(.*?)\n---', content, re.DOTALL)
            if match:
                frontmatter = match.group(1)
                version_match = re.search(r'^version:\s*(\S+)', frontmatter, re.MULTILINE)
                if version_match:
                    return version_match.group(1)
        except Exception:
            pass
        return None

    def _get_mcp_version(self) -> str:
        """Leer versión de .mcp.json."""
        mcp_file = self.project_root / '.mcp.json'
        if mcp_file.exists():
            try:
                data = json.loads(mcp_file.read_text())
                return data.get('version', None)
            except Exception:
                pass
        return None

    def is_cache_valid(
        self,
        cached_versions: Dict[str, str],
        current_versions: Dict[str, str],
        exclude_keys: list = None
    ) -> bool:
        """
        Verificar si caché sigue siendo válida comparando versiones.

        Args:
            cached_versions: Versiones cuando se guardó el caché
            current_versions: Versiones actuales del sistema
            exclude_keys: Claves a ignorar (ej. ['mcp'] si no queremos invalidar por MCP)

        Returns:
            True si versiones coinciden, False si algo cambió
        """
        if exclude_keys is None:
            exclude_keys = []

        for key, cached_version in cached_versions.items():
            if key in exclude_keys:
                continue

            current_version = current_versions.get(key)
            if current_version != cached_version:
                # Versión cambió → caché inválida
                return False

        return True

    def get_version_diff(
        self,
        cached_versions: Dict[str, str],
        current_versions: Dict[str, str]
    ) -> Dict[str, Dict[str, str]]:
        """
        Identificar qué versiones cambiaron.

        Returns:
            {'changed': {'adr_039': {'was': '1.0', 'now': '1.1'}}, 'new': {...}, 'removed': {...}}
        """
        diff = {'changed': {}, 'new': {}, 'removed': {}}

        # Detectar cambios
        for key, cached_version in cached_versions.items():
            current_version = current_versions.get(key)
            if current_version is None:
                diff['removed'][key] = cached_version
            elif current_version != cached_version:
                diff['changed'][key] = {'was': cached_version, 'now': current_version}

        # Detectar nuevas versiones
        for key, current_version in current_versions.items():
            if key not in cached_versions:
                diff['new'][key] = current_version

        return diff

    def explain_invalidation(self, diff: Dict) -> str:
        """Explicar en texto por qué caché se invalidó."""
        reasons = []

        if diff['changed']:
            for key, change in diff['changed'].items():
                reasons.append(f"  - {key}: {change['was']} -> {change['now']}")

        if diff['new']:
            for key in diff['new'].keys():
                reasons.append(f"  - {key}: [NUEVA]")

        if diff['removed']:
            for key in diff['removed'].keys():
                reasons.append(f"  - {key}: [REMOVIDA]")

        if not reasons:
            return "Sin cambios de versión"

        return "Cambios detectados:\n" + "\n".join(reasons)


# Ejemplo de uso
if __name__ == '__main__':
    tracker = VersionTracker()

    # Versiones actuales
    current = tracker.get_current_versions()
    print("Versiones actuales:")
    for key, version in sorted(current.items()):
        print(f"  {key}: {version}")

    # Simular caché vieja
    cached = {
        'adr_0039': '1.0',
        'kb_capital_gains': '2.1',
        'mcp': '1.3.2'
    }

    # Verificar si sigue siendo válida
    valid = tracker.is_cache_valid(cached, current)
    print(f"\n¿Caché válida? {valid}")

    if not valid:
        diff = tracker.get_version_diff(cached, current)
        print(tracker.explain_invalidation(diff))
