#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script de alertas: documentos de conocimiento próximos a expirar.

Uso: python scripts/check_knowledge_vigencia.py

Encuentra documentos cuya vigencia está próxima a expirar o ha expirado,
y genera un reporte.

Umbrales:
- CRÍTICO (rojo): valid_until ha pasado o pasa en <7 días
- ADVERTENCIA (amarillo): last_verified > 365 días (confidence=high)
- INFORMACIÓN (azul): próximos a expirar (>30 días)
"""

import os
import re
from pathlib import Path
from datetime import datetime, timedelta
import yaml

def extract_frontmatter(content):
    """Extract YAML frontmatter from markdown"""
    match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
    if not match:
        return None
    try:
        return yaml.safe_load(match.group(1))
    except yaml.YAMLError:
        return None

def days_until(date_str):
    """Calcular días hasta una fecha"""
    try:
        target = datetime.strptime(date_str, "%Y-%m-%d")
        delta = (target - datetime.now()).days
        return delta
    except:
        return None

def check_vigencia():
    """Check all knowledge files for vigencia issues"""

    knowledge_path = Path("knowledge")
    md_files = sorted(knowledge_path.glob("**/*.md"))

    criticos = []
    advertencias = []
    informacion = []

    for filepath in md_files:
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()
        except:
            continue

        frontmatter = extract_frontmatter(content)
        if not frontmatter:
            continue

        doc_id = frontmatter.get("id", "UNKNOWN")
        title = frontmatter.get("title", filepath.name)
        level = frontmatter.get("level", "?")
        valid_until = frontmatter.get("valid_until")
        last_verified = frontmatter.get("last_verified")
        confidence = frontmatter.get("confidence", "?")

        # Check valid_until
        if valid_until:
            days = days_until(valid_until)
            if days is not None:
                if days < 0:
                    criticos.append({
                        "id": doc_id,
                        "title": title,
                        "reason": f"EXPIRADO hace {abs(days)} días",
                        "date": valid_until,
                        "severity": "CRÍTICO"
                    })
                elif days < 7:
                    criticos.append({
                        "id": doc_id,
                        "title": title,
                        "reason": f"Vence en {days} días",
                        "date": valid_until,
                        "severity": "CRÍTICO"
                    })
                elif days < 30:
                    informacion.append({
                        "id": doc_id,
                        "title": title,
                        "reason": f"Vence en {days} días",
                        "date": valid_until,
                        "severity": "INFO"
                    })

        # Check last_verified
        if last_verified:
            days = days_until(last_verified)
            if days is not None:
                days_old = abs(days)
                if confidence == "high" and days_old > 365:
                    advertencias.append({
                        "id": doc_id,
                        "title": title,
                        "reason": f"last_verified hace {days_old} días (confidence=high)",
                        "date": last_verified,
                        "severity": "ADVERTENCIA"
                    })
                elif confidence == "medium" and days_old > 180:
                    advertencias.append({
                        "id": doc_id,
                        "title": title,
                        "reason": f"last_verified hace {days_old} días (confidence=medium)",
                        "date": last_verified,
                        "severity": "ADVERTENCIA"
                    })

    # Report
    print("\n" + "="*70)
    print("ALERTAS DE VIGENCIA — Sistema de Conocimiento")
    print("="*70 + "\n")

    if criticos:
        print(f"🔴 CRÍTICO ({len(criticos)}):\n")
        for item in sorted(criticos, key=lambda x: x["date"]):
            print(f"  {item['id']:12} | {item['reason']}")
            print(f"  {' '*12} | {item['title']}\n")

    if advertencias:
        print(f"🟡 ADVERTENCIA ({len(advertencias)}):\n")
        for item in sorted(advertencias, key=lambda x: x["date"]):
            print(f"  {item['id']:12} | {item['reason']}")
            print(f"  {' '*12} | {item['title']}\n")

    if informacion:
        print(f"ℹ️ INFORMACIÓN ({len(informacion)}):\n")
        for item in sorted(informacion, key=lambda x: x["date"]):
            print(f"  {item['id']:12} | {item['reason']}")
            print(f"  {' '*12} | {item['title']}\n")

    if not (criticos or advertencias or informacion):
        print("✅ Todos los documentos están vigentes\n")

    print("="*70 + "\n")

if __name__ == "__main__":
    check_vigencia()
