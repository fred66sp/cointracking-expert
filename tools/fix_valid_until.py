#!/usr/bin/env python3
"""
Fijar valid_until: null en documentos Nivel B.

Establece valid_until a 2027-07-03 para todos los docs Nivel B que tienen null.
"""

import re
from pathlib import Path

def fix_valid_until(file_path: Path) -> bool:
    """Fixa valid_until en un archivo. Returns True si cambio fue aplicado."""
    content = file_path.read_text(encoding='utf-8')

    # Extraer YAML
    match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
    if not match:
        return False

    yaml_block = match.group(1)
    rest = content[match.end():]

    # Verificar si es Nivel B y tiene valid_until: null o valid_until: 'null'
    if 'level: B' not in yaml_block:
        return False

    if 'valid_until: null' not in yaml_block and "valid_until: 'null'" not in yaml_block:
        return False

    # Fijar valid_until
    new_yaml_block = re.sub(
        r"valid_until:\s*(?:null|'null')",
        "valid_until: 2027-07-03",
        yaml_block
    )

    if new_yaml_block == yaml_block:
        return False

    new_content = f"---\n{new_yaml_block}\n---\n{rest}"
    file_path.write_text(new_content, encoding='utf-8')
    return True

def main():
    print(f"{'='*70}")
    print(f"FIJAR valid_until: null EN NIVEL B")
    print(f"{'='*70}")
    print()

    knowledge_dir = Path('h:/cointracking-expert/knowledge')
    md_files = sorted(knowledge_dir.rglob('*.md'))

    fixed = 0
    for file_path in md_files:
        if fix_valid_until(file_path):
            fixed += 1
            print(f"  FIXED {file_path.name}")

    print()
    print(f"{'='*70}")
    print(f"Total fijados: {fixed}")
    print(f"{'='*70}")

if __name__ == '__main__':
    main()
