#!/usr/bin/env python3
"""
Fijar IDs duplicados en documentos de conocimiento.

Estrategia:
1. Detectar todos los IDs
2. Identificar duplicados y genéricos (KB-X-XXX)
3. Asignar nuevos IDs secuenciales sin conflictos

Uso:
  python tools/fix_duplicate_ids.py --dry-run
  python tools/fix_duplicate_ids.py --apply
"""

import re
import sys
from pathlib import Path
from collections import defaultdict

def extract_yaml_and_id(content: str):
    """Extrae YAML y el ID."""
    match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
    if not match:
        return None, None, content

    yaml_block = match.group(1)
    rest = content[match.end():]

    # Extrae ID
    id_match = re.search(r'id:\s*([KB\-\d]+)', yaml_block)
    doc_id = id_match.group(1) if id_match else None

    return yaml_block, doc_id, rest

def is_generic_id(doc_id: str) -> bool:
    """Detecta IDs genéricos (KB-X-XXX)."""
    return bool(re.match(r'^KB-[A-F]-XXX$', doc_id)) if doc_id else False

def parse_id(doc_id: str) -> tuple:
    """Parsea KB-[A-F]-NNN. Returns (level, number) o (None, None) si inválido."""
    match = re.match(r'^KB-([A-F])-(\d+)$', doc_id)
    if match:
        return match.group(1), int(match.group(2))
    return None, None

def get_next_id(level: str, used_numbers: set) -> str:
    """Obtiene el siguiente ID disponible para un nivel."""
    num = 1
    while num in used_numbers:
        num += 1
    return f"KB-{level}-{str(num).zfill(3)}"

def main():
    dry_run = '--apply' not in sys.argv

    print(f"{'='*70}")
    print(f"FIJAR IDs DUPLICADOS")
    print(f"{'='*70}")
    print(f"Modo: {'DRY RUN' if dry_run else 'APPLY'}")
    print()

    knowledge_dir = Path('h:/cointracking-expert/knowledge')
    md_files = sorted(knowledge_dir.rglob('*.md'))

    # Recopilar todos los IDs y archivos
    id_to_files = defaultdict(list)
    all_files = {}

    for file_path in md_files:
        try:
            content = file_path.read_text(encoding='utf-8')
            yaml_block, doc_id, rest = extract_yaml_and_id(content)
            if doc_id:
                id_to_files[doc_id].append(file_path)
                all_files[file_path] = (yaml_block, doc_id, rest)
        except Exception as e:
            print(f"ERROR {file_path}: {e}")

    # Identificar duplicados y genéricos
    to_fix = []
    for doc_id, files in id_to_files.items():
        if len(files) > 1 or is_generic_id(doc_id):
            to_fix.extend(files)

    # Agrupar by nivel
    by_level = defaultdict(lambda: {'used': set(), 'files': []})

    for file_path in all_files:
        yaml_block, doc_id, rest = all_files[file_path]
        if doc_id:
            level, num = parse_id(doc_id)
            if level:
                by_level[level]['used'].add(num)
                by_level[level]['files'].append((file_path, doc_id, yaml_block, rest))

    # Reasignar IDs a duplicados
    replacements = {}
    for doc_id, files in id_to_files.items():
        if len(files) > 1 or is_generic_id(doc_id):
            level, _ = parse_id(doc_id)
            if level:
                for i, file_path in enumerate(files):
                    if i == 0:
                        # First one keeps (or gets fixed ID if generic)
                        if is_generic_id(doc_id):
                            new_id = get_next_id(level, by_level[level]['used'])
                            by_level[level]['used'].add(int(new_id.split('-')[-1]))
                            replacements[file_path] = new_id
                    else:
                        # Resto gets new IDs
                        new_id = get_next_id(level, by_level[level]['used'])
                        by_level[level]['used'].add(int(new_id.split('-')[-1]))
                        replacements[file_path] = new_id

    # Aplicar cambios
    fixed = 0
    for file_path, new_id in replacements.items():
        if file_path not in all_files:
            continue

        yaml_block, old_id, rest = all_files[file_path]
        new_yaml = re.sub(r"id:\s*[KB\-\d]+", f"id: {new_id}", yaml_block)

        if dry_run:
            print(f"  {file_path.name}: {old_id} → {new_id}")
        else:
            new_content = f"---\n{new_yaml}\n---\n{rest}"
            file_path.write_text(new_content, encoding='utf-8')
            print(f"  FIXED {file_path.name}: {old_id} → {new_id}")
            fixed += 1

    print()
    print(f"{'='*70}")
    print(f"Total a reasignar: {len(replacements)}")
    if not dry_run:
        print(f"Total fijados: {fixed}")
    print(f"{'='*70}")

if __name__ == '__main__':
    main()
