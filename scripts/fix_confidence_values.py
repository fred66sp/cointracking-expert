#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script para corregir valores de confidence inválidos en casos migrados

Mapea valores antiguos a valores válidos:
  - "verificado" → "high"
  - "probable" → "medium"
  - "hipotesis" → "medium"
  - "pendiente_verificar" → "low"
  - "critico" → "high"
"""

import re
from pathlib import Path

CASES_DIR = Path("knowledge/cases")

# Mapeo de valores antiguos a nuevos
CONFIDENCE_MAP = {
    "verificado": "high",
    "probable": "medium",
    "hipotesis": "medium",
    "pendiente_verificar": "low",
    "critico": "high"
}

def fix_confidence():
    """Corregir confidence en todos los casos"""

    md_files = sorted(CASES_DIR.glob("ct-*.md"))

    print(f"Corrigiendo {len(md_files)} archivos...\n")

    fixed_count = 0
    for filepath in md_files:
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()

            content_modified = False
            # Buscar y reemplazar confidence
            for old_val, new_val in CONFIDENCE_MAP.items():
                pattern = f'confidence: "{old_val}"'
                if pattern in content:
                    new_pattern = f'confidence: "{new_val}"'
                    content = content.replace(pattern, new_pattern)
                    content_modified = True

            # Guardar solo si cambió
            if content_modified:
                with open(filepath, 'w', encoding='utf-8') as f:
                    f.write(content)
                fixed_count += 1
                print(f"[OK] {filepath.name}")

        except Exception as e:
            print(f"[SKIP] {filepath.name}: {str(e)[:50]}")

    print(f"\n[DONE] {fixed_count} archivos corregidos")

if __name__ == "__main__":
    fix_confidence()
