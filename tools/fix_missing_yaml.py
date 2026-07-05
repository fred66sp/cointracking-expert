#!/usr/bin/env python3
"""
Agregar YAML faltante a documentos de navegación/referencia.
"""

import re
from pathlib import Path

# Mapeo: ruta → (id, title, level)
FIXES = {
    'knowledge/QUICK_START.md': ('KB-F1-001', 'Guía Rápida (5 minutos)', 'F'),
    'knowledge/NAVIGATION_MAP.md': ('KB-F1-002', 'Mapa de Navegación por Función', 'F'),
    'knowledge/CHEAT_SHEET.md': ('KB-F1-003', 'Hoja de Referencia Rápida', 'F'),
    'knowledge/INDEX_MASTER.md': ('KB-F1-004', 'Índice Maestro de la Base de Conocimiento', 'F'),
    'knowledge/cointracking/CT_LIST_FORMATS.md': ('KB-B1-020', 'Formatos CT-List para Listas en Conversación (ADR-025)', 'B'),
    'knowledge/cointracking/DOCUMENT_CHECKLIST.md': ('KB-B1-021', 'Checklist de Documentos CoinTracking', 'B'),
    'knowledge/cointracking/TROUBLESHOOTING.md': ('KB-B1-022', 'Troubleshooting: Síntomas y Soluciones', 'B'),
    'knowledge/KNOWLEDGE_MAINTENANCE.md': ('KB-F1-005', 'Mantenimiento de la Base de Conocimiento', 'F'),
    'knowledge/cases/ct-013-wallet-externa-no-importada-fondos-desap.md': ('KB-C1-013', 'Caso: Wallet Externa No Importada — Fondos Desaparecen', 'C'),
}

def add_yaml(file_path: Path, doc_id: str, title: str, level: str) -> bool:
    """Agrega YAML completo a un archivo."""
    try:
        content = file_path.read_text(encoding='utf-8')
    except Exception as e:
        print(f"ERROR reading {file_path}: {e}")
        return False

    # Si el archivo comienza con --- pero está vacío, reemplaza
    if content.startswith('---\n\n'):
        yaml_block = f"""---
id: {doc_id}
title: "{title}"
level: {level}
domain: cointracking
source: "Documentación interna"
authority: reference
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - navigation
  - reference

notes: "Documento de navegación/referencia de la base de conocimiento"
---

"""
        # Quitar el --- --- inicial vacío
        rest = re.sub(r'^---\n\n', '', content)
        new_content = yaml_block + rest
    else:
        # Agregar YAML al inicio
        yaml_block = f"""---
id: {doc_id}
title: "{title}"
level: {level}
domain: cointracking
source: "Documentación interna"
authority: reference
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: medium
version: 1.0

tags:
  - navigation
  - reference

notes: "Documento de navegación/referencia de la base de conocimiento"
---

"""
        new_content = yaml_block + content

    try:
        file_path.write_text(new_content, encoding='utf-8')
        print(f"FIXED: {file_path.name}")
        return True
    except Exception as e:
        print(f"ERROR writing {file_path}: {e}")
        return False

def main():
    base = Path('h:/cointracking-expert')
    fixed = 0
    for rel_path, (doc_id, title, level) in FIXES.items():
        file_path = base / rel_path
        if file_path.exists():
            if add_yaml(file_path, doc_id, title, level):
                fixed += 1
        else:
            print(f"NOT FOUND: {file_path}")

    print(f"\n[OK] {fixed} archivos fijados")

if __name__ == '__main__':
    main()
