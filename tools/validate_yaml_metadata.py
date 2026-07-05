#!/usr/bin/env python3
"""
Validar integridad de metadatos YAML en la base de conocimiento.

Chequeos:
1. Todos los documentos tienen YAML frontmatter
2. IDs están en formato correcto (KB-[NIVEL]-NNN)
3. No hay duplicados de ID
4. Levels son válidos (A, B, C, D, E, F)
5. valid_until está definido para Nivel A/B (nunca null)
6. Campos requeridos presentes

Uso:
  python tools/validate_yaml_metadata.py
"""

import re
import sys
from pathlib import Path
from typing import Dict, List, Tuple

import yaml

def extract_yaml(content: str) -> Dict | None:
    """Extrae y parsea YAML frontmatter."""
    if not content.startswith('---'):
        return None

    match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
    if not match:
        return None

    try:
        return yaml.safe_load(match.group(1))
    except yaml.YAMLError:
        return None

def validate_id(doc_id: str) -> Tuple[bool, str]:
    """Valida formato de ID."""
    if not doc_id:
        return False, "ID vacio"

    # Patrón: KB-[A-F]-NNN
    match = re.match(r'^KB-([A-F])-(\d{3,})$', doc_id)
    if not match:
        return False, f"Formato incorrecto: {doc_id} (esperado KB-[A-F]-NNN)"

    return True, doc_id

def validate_level(level: str) -> bool:
    """Valida que level sea A-F."""
    return level in ['A', 'B', 'C', 'D', 'E', 'F']

def validate_date(date_str: str | None, field_name: str) -> Tuple[bool, str]:
    """Valida formato de fecha."""
    if date_str is None:
        return False, f"{field_name}: null (debe tener fecha)"

    if isinstance(date_str, str):
        if re.match(r'^\d{4}-\d{2}-\d{2}$', date_str):
            return True, date_str
        else:
            return False, f"{field_name}: formato incorrecto ({date_str})"

    return False, f"{field_name}: tipo incorrecto ({type(date_str)})"

def validate_document(file_path: Path) -> List[str]:
    """
    Valida un documento.
    Returns: lista de errores encontrados (vacía si OK)
    """
    errors = []

    try:
        content = file_path.read_text(encoding='utf-8')
    except Exception as e:
        return [f"ERROR reading: {e}"]

    yaml_data = extract_yaml(content)
    if not yaml_data:
        return [f"ERROR: no valid YAML frontmatter"]

    # Chequeos
    doc_id = yaml_data.get('id')
    if not doc_id:
        errors.append("MISSING: id")
    else:
        valid, msg = validate_id(doc_id)
        if not valid:
            errors.append(f"INVALID id: {msg}")

    level = yaml_data.get('level')
    if not level:
        errors.append("MISSING: level")
    elif not validate_level(level):
        errors.append(f"INVALID level: {level} (debe ser A-F)")

    # valid_until crítico para Nivel A/B
    if level in ['A', 'B']:
        valid_until = yaml_data.get('valid_until')
        if valid_until is None or valid_until == 'null':
            errors.append(f"CRITICAL: valid_until is null (Nivel {level})")
        else:
            valid, msg = validate_date(valid_until, 'valid_until')
            if not valid:
                errors.append(f"INVALID valid_until: {msg}")

    # valid_from
    valid_from = yaml_data.get('valid_from')
    if valid_from:
        valid, msg = validate_date(valid_from, 'valid_from')
        if not valid:
            errors.append(f"INVALID valid_from: {msg}")

    # Campos opcionales pero comunes
    for field in ['title', 'domain', 'source', 'authority', 'last_verified', 'confidence']:
        if field not in yaml_data:
            errors.append(f"MISSING: {field}")

    return errors

def main():
    print(f"{'='*70}")
    print(f"VALIDAR INTEGRIDAD DE METADATOS YAML")
    print(f"{'='*70}")
    print()

    knowledge_dir = Path('h:/cointracking-expert/knowledge')
    md_files = sorted(knowledge_dir.rglob('*.md'))

    total = 0
    ok_count = 0
    warning_count = 0
    error_count = 0

    ids_seen: Dict[str, str] = {}
    critical_errors: List[str] = []
    warnings: List[str] = []

    for file_path in md_files:
        errors = validate_document(file_path)
        total += 1

        if not errors:
            ok_count += 1
        else:
            # Extraer ID para deduplicación
            content = file_path.read_text(encoding='utf-8')
            yaml_data = extract_yaml(content)
            doc_id = yaml_data.get('id', 'NO_ID') if yaml_data else 'NO_ID'

            has_critical = any('CRITICAL' in e for e in errors)
            if has_critical:
                error_count += 1
                for error in errors:
                    if 'CRITICAL' in error:
                        critical_errors.append(f"  {file_path.name}: {error}")
            else:
                warning_count += 1
                for error in errors:
                    warnings.append(f"  {file_path.name}: {error}")

            # Track ID
            if doc_id != 'NO_ID':
                if doc_id in ids_seen:
                    critical_errors.append(f"  DUPLICATE ID {doc_id}: {ids_seen[doc_id]} y {file_path.name}")
                else:
                    ids_seen[doc_id] = str(file_path.name)

    print(f"RESULTADOS")
    print(f"{'='*70}")
    print(f"Total documentos:       {total}")
    print(f"Sin problemas:          {ok_count}")
    print(f"Con warnings:           {warning_count}")
    print(f"Con errores críticos:   {error_count}")
    print()

    if critical_errors:
        print(f"ERRORES CRITICOS ({len(critical_errors)}):")
        print(f"{'-'*70}")
        for err in critical_errors[:20]:
            print(err)
        if len(critical_errors) > 20:
            print(f"  ... y {len(critical_errors) - 20} más")
        print()

    if warnings:
        print(f"WARNINGS ({len(warnings)}):")
        print(f"{'-'*70}")
        for warn in warnings[:20]:
            print(warn)
        if len(warnings) > 20:
            print(f"  ... y {len(warnings) - 20} más")
        print()

    print(f"{'='*70}")
    if error_count == 0:
        print("[OK] Metadatos YAML validados correctamente")
    else:
        print(f"[ATENCIÓN] {error_count} documentos con errores críticos")
    print(f"{'='*70}")

if __name__ == '__main__':
    main()
