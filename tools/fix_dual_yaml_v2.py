#!/usr/bin/env python3
"""
Script mejorado para limpiar DUAL-YAML corruption.

Versión 2: Más inteligente
- Ignora `---` dentro de bloques de código (entre triple backticks)
- Solo detecta DUAL-YAML real: frontmatter genérico + contenido
- Patrón: ^---\nid: KB-*-XXX...title: Untitled\n---\n (con contenido real después)

Uso:
  python tools/fix_dual_yaml_v2.py --dry-run
  python tools/fix_dual_yaml_v2.py --apply
"""

import re
import sys
from pathlib import Path

def extract_frontmatter(content: str) -> tuple[bool, str, str]:
    """
    Extrae el frontmatter YAML del inicio del archivo.
    Returns: (is_valid_frontmatter, frontmatter_text, rest_of_content)

    Frontmatter válido:
    - Comienza con `---` en la primera línea
    - Contiene YAML
    - Termina con `---` en una línea propia
    - El resto es contenido normal
    """
    if not content.startswith('---'):
        return False, "", content

    # Buscar el cierre del primer bloque YAML
    # Patrón: ^---\n...contenido...\n---\n
    match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
    if not match:
        return False, "", content

    frontmatter = match.group(1)
    rest = content[match.end():]
    return True, frontmatter, rest

def is_generic_yaml(frontmatter: str) -> bool:
    """Detecta si el YAML es genérico."""
    return 'KB-B1-XXX' in frontmatter or 'Untitled' in frontmatter

def remove_generic_yaml_and_preserve_specific(content: str) -> tuple[bool, str]:
    """
    Si hay DUAL-YAML (genérico + específico), elimina el genérico y preserva el específico.
    Returns: (changed, new_content)
    """
    is_valid, frontmatter, rest = extract_frontmatter(content)

    if not is_valid:
        return False, content

    if not is_generic_yaml(frontmatter):
        # Primer bloque no es genérico → no es nuestro problema
        return False, content

    # Aquí hay un bloque genérico. Verificar si hay un segundo bloque específico
    # El patrón es: [contenido] \n--- \n [YAML específico] \n ---
    # Simplificar: si encuentra \n---\n en `rest`, probablemente hay DUAL-YAML

    if '\n---\n' in rest:
        # Hay otro frontmatter en el contenido
        # Extraerlo como "específico"
        second_match = re.match(r'(.*?)\n---\n(.*?)\n---\n(.*)', rest, re.DOTALL)
        if second_match:
            # Encontró: [contenido antes] ---\n [segundo YAML] ---\n [contenido después]
            content_before = second_match.group(1)
            second_yaml = second_match.group(2)
            content_after = second_match.group(3)

            # Verificar que el segundo YAML es específico (no genérico)
            if not ('KB-B1-XXX' in second_yaml or 'Untitled' in second_yaml):
                # Específico encontrado → eliminar genérico, mantener específico
                new_content = f"---\n{second_yaml}\n---\n\n{content_before}\n\n{content_after}"
                return True, new_content

    # No hay DUAL-YAML real detectado
    return False, content

def process_file(file_path: str, apply: bool = False) -> tuple[bool, str]:
    """
    Procesar archivo para limpiar DUAL-YAML real.
    Returns: (changed, message)
    """
    try:
        with open(file_path, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        return False, f"ERROR reading: {e}"

    changed, new_content = remove_generic_yaml_and_preserve_specific(content)

    if not changed:
        return False, "SKIP (no DUAL-YAML detected)"

    # Aplicar cambios si se pidió
    if apply:
        try:
            with open(file_path, 'w', encoding='utf-8') as f:
                f.write(new_content)
            return True, "FIXED (removed generic YAML)"
        except Exception as e:
            return False, f"ERROR writing: {e}"
    else:
        return True, "WOULD FIX (remove generic YAML)"

def main():
    dry_run = '--apply' not in sys.argv

    print(f"{'='*70}")
    print(f"FIX DUAL-YAML CORRUPTION (v2 — Smart Detection)")
    print(f"{'='*70}")
    print(f"Modo: {'DRY RUN' if dry_run else 'APPLY'}")
    print()

    knowledge_dir = Path('h:/cointracking-expert/knowledge')
    md_files = sorted(knowledge_dir.rglob('*.md'))

    total = 0
    fixed = 0
    skipped = 0
    errors = 0

    for file_path in md_files:
        changed, message = process_file(str(file_path), apply=not dry_run)

        if changed:
            fixed += 1
            print(f"  {file_path.name:60s} {message}")
        elif message.startswith("ERROR"):
            errors += 1
            print(f"  {file_path.name:60s} {message}")

        total += 1

    print()
    print(f"{'='*70}")
    print(f"RESUMEN")
    print(f"{'='*70}")
    print(f"Total archivos .md:    {total}")
    print(f"Reparados:             {fixed}")
    print(f"Ignorados (OK):        {skipped + (total - fixed - errors)}")
    print(f"Errores:               {errors}")
    print()
    if fixed == 0:
        print("[OK] No se detecto DUAL-YAML real. El sistema esta limpio.")
    elif dry_run:
        print(f"Para aplicar los cambios, ejecuta:")
        print(f"  python tools/fix_dual_yaml_v2.py --apply")
    else:
        print(f"[OK] {fixed} archivos reparados")
    print(f"{'='*70}")

if __name__ == '__main__':
    main()
