#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script de validación de metadatos YAML en documentos de conocimiento.

Uso: python scripts/validate_knowledge_metadata.py [--fix]

Valida que todos los documentos en knowledge/ tengan metadatos completos y correctos.

Reglas:
- Todo .md debe tener frontmatter YAML
- Campos obligatorios: id, title, level, domain, source, authority, last_verified, confidence, version
- id debe ser único
- level debe estar en {A, B, C, D, E, F}
- authority debe estar en {official, verified, empirical, reference}
- confidence debe estar en {high, medium, low}
- valid_until: null es solo permitido para Nivel E
"""

import os
import re
import sys
from pathlib import Path
from datetime import datetime
import yaml

REQUIRED_FIELDS = {
    "id", "title", "level", "domain", "source", "authority",
    "last_verified", "valid_from", "valid_until", "confidence", "version"
}

VALID_LEVELS = {"A", "B", "C", "D", "E", "F"}
VALID_AUTHORITIES = {"official", "verified", "empirical", "reference"}
VALID_CONFIDENCES = {"high", "medium", "low"}

errors = []
warnings = []
ids_seen = set()

def extract_frontmatter(content):
    """Extract YAML frontmatter from markdown"""
    match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
    if not match:
        return None
    try:
        return yaml.safe_load(match.group(1))
    except yaml.YAMLError as e:
        return None

def validate_file(filepath):
    """Validate a single markdown file"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
    except Exception as e:
        errors.append(f"{filepath}: Error leyendo archivo: {e}")
        return

    # Extract frontmatter
    frontmatter = extract_frontmatter(content)
    if not frontmatter:
        errors.append(f"{filepath}: No tiene frontmatter YAML válido")
        return

    # Check required fields
    for field in REQUIRED_FIELDS:
        if field not in frontmatter:
            errors.append(f"{filepath}: Falta campo requerido: {field}")

    # Validate id (must be unique)
    doc_id = frontmatter.get("id")
    if doc_id:
        if doc_id in ids_seen:
            errors.append(f"{filepath}: ID duplicado: {doc_id}")
        else:
            ids_seen.add(doc_id)

    # Validate level
    level = frontmatter.get("level")
    if level and level not in VALID_LEVELS:
        errors.append(f"{filepath}: level inválido: {level} (debe estar en {VALID_LEVELS})")

    # Validate authority
    authority = frontmatter.get("authority")
    if authority and authority not in VALID_AUTHORITIES:
        errors.append(f"{filepath}: authority inválido: {authority}")

    # Validate confidence
    confidence = frontmatter.get("confidence")
    if confidence and confidence not in VALID_CONFIDENCES:
        errors.append(f"{filepath}: confidence inválido: {confidence}")

    # Validate valid_until for Nivel A (CRÍTICO)
    valid_until = frontmatter.get("valid_until")
    if level == "A" and valid_until is None:
        errors.append(f"{filepath}: CRÍTICO - Nivel A no debe tener valid_until: null")

    # Validate last_verified is not in the future
    last_verified = frontmatter.get("last_verified")
    if last_verified:
        try:
            # Handle both string and datetime.date objects
            if isinstance(last_verified, str):
                verified_date = datetime.strptime(last_verified, "%Y-%m-%d")
            else:
                verified_date = datetime.combine(last_verified, datetime.min.time())

            if verified_date > datetime.now():
                errors.append(f"{filepath}: last_verified es futuro: {last_verified}")
        except (ValueError, TypeError):
            errors.append(f"{filepath}: last_verified formato inválido: {last_verified}")

    # Warn if confidence is low
    if confidence == "low":
        warnings.append(f"{filepath}: confidence=low (requiere revisión)")

def main():
    knowledge_path = Path("knowledge")

    if not knowledge_path.exists():
        print("❌ No se encontró carpeta knowledge/")
        sys.exit(1)

    # Find all .md files
    md_files = sorted(knowledge_path.glob("**/*.md"))

    if not md_files:
        print("⚠️ No se encontraron archivos .md en knowledge/")
        return

    print(f"Validando {len(md_files)} archivos...\n")

    for filepath in md_files:
        validate_file(filepath)

    # Report
    print(f"\n{'='*60}")
    print(f"RESULTADOS")
    print(f"{'='*60}\n")

    if errors:
        print(f"[ERROR] ERRORES ({len(errors)}):\n")
        for error in errors:
            print(f"  • {error}")
        print()

    if warnings:
        print(f"[WARN] ADVERTENCIAS ({len(warnings)}):\n")
        for warning in warnings:
            print(f"  • {warning}")
        print()

    print(f"[OK] {len(md_files) - len(errors)} archivos válidos")
    print(f"[OK] {len(ids_seen)} IDs únicos")

    if errors:
        sys.exit(1)
    else:
        print("\n[OK] Validación exitosa")
        sys.exit(0)

if __name__ == "__main__":
    main()
