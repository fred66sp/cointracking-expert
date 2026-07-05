#!/usr/bin/env python3
"""
Script para agregar metadatos YAML template a documentos sin frontmatter.

Uso:
  python tools/add_yaml_to_docs.py --dry-run    # Ver qué se cambiaría
  python tools/add_yaml_to_docs.py --apply      # Aplicar cambios reales
"""

import os
import re
import sys
from pathlib import Path
from datetime import datetime
from typing import List, Tuple

# Documentos que ya tienen YAML y no tocar
SKIP_FILES = {
    'INDEX_MASTER.md', 'README.md', 'index.md', 'INDEX.md',
    'FOUNDATION.md', '.gitignore', 'LICENSE',
}

# Mapeo de directorios a nivel de conocimiento
DIR_TO_LEVEL = {
    'knowledge/authorities': 'A',
    'knowledge/taxation': 'A',
    'knowledge/exchanges/official': 'A',
    'knowledge/cointracking/official': 'A',
    'knowledge/cointracking/behavioral': 'B',
    'knowledge/exchanges/behavioral': 'B',
    'knowledge/blockchains': 'B',
    'knowledge/wallets': 'B',
    'knowledge/cases': 'C',
    'knowledge/patterns': 'C',
    'knowledge/procedures': 'C',
    'knowledge/checklists': 'D',
    'knowledge/flows': 'D',
    'knowledge/reference': 'E',
}

YAML_TEMPLATE = '''---
id: {id}
title: "{title}"
level: {level}
domain: cointracking
source: "Internal documentation"
authority: verified
last_verified: {date}
valid_from: 2024-01-01
valid_until: {valid_until}
confidence: medium
version: 1.0

tags:
  - todo
  - needs-review

notes: "Metadatos agregados automáticamente. Verificar y actualizar conforme ADR-032."
---

'''

def has_yaml_frontmatter(content: str) -> bool:
    """Verificar si el documento ya tiene YAML frontmatter."""
    return content.startswith('---') and '---' in content[3:100]

def get_level_from_path(file_path: str) -> str:
    """Determinar nivel de conocimiento desde la ruta."""
    for prefix, level in DIR_TO_LEVEL.items():
        if prefix in file_path:
            return level
    return 'B'  # Default

def get_next_id(level: str, directory: str) -> str:
    """Generar ID basado en nivel y directorio."""
    level_map = {
        'A': 'KB-A1',
        'B': 'KB-B1',
        'C': 'KB-C1',
        'D': 'KB-D1',
        'E': 'KB-E1',
        'F': 'KB-F1',
    }
    return f"{level_map.get(level, 'KB-B1')}-XXX"  # XXX debe editarse manualmente

def extract_title(content: str) -> str:
    """Extraer título de la primera línea que parece un título."""
    lines = content.split('\n')[:10]
    for line in lines:
        if line.startswith('#'):
            return line.strip('# ').strip()
    return "Untitled Document"

def get_valid_until(level: str) -> str:
    """Determinar valid_until según el nivel."""
    if level == 'A':
        return '2027-07-05'  # Nivel A: 1 año
    elif level == 'B':
        return '2027-12-31'  # Nivel B: 18 meses
    else:
        return 'null'  # Otros: indefinido

def process_file(file_path: str, apply: bool = False) -> Tuple[bool, str]:
    """
    Procesar un archivo para agregar YAML si no lo tiene.

    Returns: (changed, message)
    """
    if file_path.endswith(tuple(SKIP_FILES)):
        return False, f"SKIP: {file_path} (lista de exclusión)"

    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        return False, f"ERROR: {file_path} (no se puede leer: {e})"

    # Si ya tiene YAML, saltar
    if has_yaml_frontmatter(content):
        return False, f"SKIP: {file_path} (ya tiene YAML)"

    # Generar metadatos
    level = get_level_from_path(file_path)
    title = extract_title(content)
    doc_id = get_next_id(level, os.path.dirname(file_path))
    valid_until = get_valid_until(level)
    date = datetime.now().strftime('%Y-%m-%d')

    yaml_content = YAML_TEMPLATE.format(
        id=doc_id,
        title=title,
        level=level,
        date=date,
        valid_until=valid_until,
    )

    new_content = yaml_content + content

    if apply:
        try:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            return True, f"UPDATED: {file_path} (added YAML, level {level})"
        except Exception as e:
            return False, f"ERROR: {file_path} (no se puede escribir: {e})"
    else:
        return True, f"WOULD UPDATE: {file_path} (add YAML, level {level})"

def main():
    dry_run = '--dry-run' in sys.argv or not ('--apply' in sys.argv)

    print(f"{'='*70}")
    print(f"AGREGAR YAML FRONTMATTER A DOCUMENTOS")
    print(f"{'='*70}")
    print(f"Modo: {'DRY RUN' if dry_run else 'APPLY'}")
    print()

    # Buscar documentos markdown
    knowledge_dir = Path('h:/cointracking-expert/knowledge')
    md_files = list(knowledge_dir.rglob('*.md'))

    total = len(md_files)
    updated = 0
    errors = 0
    skipped = 0

    changes = []

    for file_path in sorted(md_files):
        changed, message = process_file(str(file_path), apply=not dry_run)

        if changed and not message.startswith('SKIP'):
            updated += 1
            changes.append(message)
        elif message.startswith('SKIP'):
            skipped += 1
        elif message.startswith('ERROR'):
            errors += 1
            print(f"  {message}")

    # Resumen
    print()
    print(f"{'='*70}")
    print(f"RESUMEN")
    print(f"{'='*70}")
    print(f"Total archivos .md encontrados: {total}")
    print(f"Serían actualizados: {updated}")
    print(f"Ignorados (ya tienen YAML): {skipped}")
    print(f"Errores: {errors}")
    print()

    if updated > 0 and updated <= 10:
        print("Cambios que se aplicarían:")
        for change in changes[:10]:
            print(f"  {change}")
        if len(changes) > 10:
            print(f"  ... y {len(changes) - 10} más")
    elif updated > 10:
        print(f"Primeros 10 cambios:")
        for change in changes[:10]:
            print(f"  {change}")
        print(f"  ... y {updated - 10} más")

    print()
    print(f"{'='*70}")
    if dry_run:
        print("Para aplicar los cambios, ejecuta:")
        print("  python tools/add_yaml_to_docs.py --apply")
    else:
        print(f"✅ {updated} archivos actualizados")
    print(f"{'='*70}")

if __name__ == '__main__':
    main()
