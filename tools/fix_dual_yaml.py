#!/usr/bin/env python3
"""
Script para limpiar DUAL-YAML corruption en documentos.

Problema: Script anterior agregó YAML genérico SIN eliminar el existente.
Solución: Detectar y eliminar el primer bloque (genérico) si hay dos bloques.
Mantener: El segundo bloque (específico y correcto).

Uso:
  python tools/fix_dual_yaml.py --dry-run    # Ver qué haría
  python tools/fix_dual_yaml.py --apply      # Aplicar cambios reales
"""

import os
import re
import sys
from pathlib import Path

def count_yaml_blocks(content: str) -> int:
    """Contar bloques YAML (--- ... ---)."""
    matches = re.findall(r'^---\s*$', content, re.MULTILINE)
    return len(matches)

def has_generic_yaml(content: str) -> bool:
    """Detectar si el PRIMER bloque es genérico (KB-B1-XXX, Untitled)."""
    # Extraer primer bloque YAML
    match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
    if not match:
        return False

    first_block = match.group(1)
    # Es genérico si tiene id: KB-*-XXX o title: "Untitled Document"
    return 'KB-B1-XXX' in first_block or 'Untitled Document' in first_block

def remove_first_yaml_block(content: str) -> str:
    """Remover el primer bloque YAML si hay dos."""
    # Verificar que hay exactamente dos bloques
    if count_yaml_blocks(content) != 2:
        return content

    # Remover el primero (encuentra el segundo ---)
    match = re.match(r'^---\n.*?\n---\n(.*)', content, re.DOTALL)
    if match:
        return match.group(1)
    return content

def process_file(file_path: str, apply: bool = False) -> tuple[bool, str]:
    """
    Procesar archivo para limpiar DUAL-YAML.
    Returns: (changed, message)
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        return False, f"ERROR reading {file_path}: {e}"

    # Contar bloques
    num_blocks = count_yaml_blocks(content)

    if num_blocks < 2:
        return False, f"SKIP {file_path} (blocks={num_blocks})"

    if num_blocks > 2:
        return False, f"ERROR {file_path} (blocks={num_blocks}, >2)"

    # Verificar si el primero es genérico
    if not has_generic_yaml(content):
        return False, f"SKIP {file_path} (first YAML not generic)"

    # Limpiar
    if apply:
        try:
            new_content = remove_first_yaml_block(content)
            if new_content == content:
                return False, f"ERROR {file_path} (failed to remove first block)"

            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            return True, f"FIXED {file_path} (removed generic YAML)"
        except Exception as e:
            return False, f"ERROR writing {file_path}: {e}"
    else:
        return True, f"WOULD FIX {file_path} (remove generic YAML)"

def main():
    dry_run = '--dry-run' in sys.argv or not ('--apply' in sys.argv)

    print(f"{'='*70}")
    print(f"FIX DUAL-YAML CORRUPTION")
    print(f"{'='*70}")
    print(f"Modo: {'DRY RUN' if dry_run else 'APPLY'}")
    print()

    knowledge_dir = Path('h:/cointracking-expert/knowledge')
    md_files = list(knowledge_dir.rglob('*.md'))

    total = 0
    fixed = 0
    skipped = 0
    errors = 0

    for file_path in sorted(md_files):
        changed, message = process_file(str(file_path), apply=not dry_run)

        if changed:
            fixed += 1
            print(f"  {message}")
        elif message.startswith("SKIP"):
            skipped += 1
        elif message.startswith("ERROR"):
            errors += 1
            print(f"  {message}")

        total += 1

    print()
    print(f"{'='*70}")
    print(f"RESUMEN")
    print(f"{'='*70}")
    print(f"Total archivos .md: {total}")
    print(f"Serían reparados: {fixed}")
    print(f"Ignorados (OK): {skipped}")
    print(f"Errores: {errors}")
    print()
    print(f"{'='*70}")
    if dry_run:
        print("Para aplicar los cambios, ejecuta:")
        print("  python tools/fix_dual_yaml.py --apply")
    else:
        print(f"✅ {fixed} archivos reparados")
    print(f"{'='*70}")

if __name__ == '__main__':
    main()
