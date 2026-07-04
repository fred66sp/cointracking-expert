#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Script de migracion: casos YAML → archivos MD individuales con metadatos

Uso: python scripts/migrate_cases_yaml_to_md.py

Lee cointracking_casos_v2.yaml y genera un archivo .md por cada caso
en knowledge/cases/CT-XXX-*.md con frontmatter YAML completo.
"""

import yaml
import os
from pathlib import Path
from datetime import datetime

YAML_PATH = Path("knowledge/cases/cointracking_casos_v2.yaml")
OUTPUT_DIR = Path("knowledge/cases")

def slugify(title):
    """Convertir titulo a slug"""
    return (
        title.lower()
        .replace("á", "a")
        .replace("é", "e")
        .replace("í", "i")
        .replace("ó", "o")
        .replace("ú", "u")
        .replace(" ", "-")
        .replace("(", "")
        .replace(")", "")
        .replace(",", "")
        .replace(".", "")
        .replace('"', "")
        .replace("'", "")
    )

def migrate_cases():
    """Migrar todos los casos"""

    with open(YAML_PATH, 'r', encoding='utf-8') as f:
        data = yaml.safe_load(f)

    cases = data  # Es una lista de casos

    print(f"Migrando {len(cases)} casos...\n")

    for idx, case in enumerate(cases, 1):
        case_id = case.get('id', f'CT-{idx:03d}')
        titulo = case.get('titulo', 'Sin titulo')

        # Generar nombre de archivo
        slug = slugify(titulo)[:40]  # Limitar a 40 caracteres
        filename = f"{case_id.lower()}-{slug}.md"
        filepath = OUTPUT_DIR / filename

        # Generar KB-ID
        kb_id = f"KB-C1-{idx:03d}"

        # Metadatos
        frontmatter = {
            "id": kb_id,
            "title": f"Caso {case_id}: {titulo}",
            "level": "C",
            "domain": "cointracking",
            "source": "Análisis de casos reales auditados",
            "authority": "verified",
            "last_verified": "2026-07-05",
            "valid_from": "2024-01-01",
            "valid_until": None,
            "confidence": case.get("nivel_confianza", "medium"),
            "version": "1.0",
            "related_adr": ["ADR-003", "ADR-009", "ADR-010"],
            "related_docs": [
                "knowledge/patterns/INDEX.md",
                "knowledge/cointracking/COST_BASIS_AND_VALIDATION.md"
            ],
            "tags": [
                "case",
                case.get("categoria", "audit"),
                "verified",
                "operativo"
            ]
        }

        # Generar contenido MD
        lines = []

        # Metadatos YAML
        lines.append("---")
        for key, value in frontmatter.items():
            if value is None:
                lines.append(f"{key}:")
            elif isinstance(value, list):
                lines.append(f"{key}:")
                for item in value:
                    lines.append(f"  - {item}")
            else:
                lines.append(f'{key}: "{value}"')
        lines.append("---\n")

        # Titulo
        lines.append(f"# {case_id}: {titulo}\n")

        # Síntomas
        if case.get("sintomas"):
            lines.append("## Síntomas\n")
            for sintoma in case["sintomas"]:
                lines.append(f"- {sintoma}")
            lines.append("")

        # Causa probable
        if case.get("causa_probable"):
            lines.append("## Causa Probable\n")
            cp = case["causa_probable"]
            if isinstance(cp, dict):
                if cp.get("hecho"):
                    lines.append(f"**Hecho:** {cp['hecho']}\n")
                if cp.get("hipotesis"):
                    lines.append(f"**Hipótesis:** {cp['hipotesis']}\n")
                if cp.get("supuesto"):
                    lines.append(f"**Supuesto:** {cp['supuesto']}\n")
            else:
                lines.append(f"{cp}\n")

        # Evidencia mínima
        if case.get("evidencia_minima"):
            lines.append("## Evidencia Mínima\n")
            for evidencia in case["evidencia_minima"]:
                lines.append(f"- {evidencia}")
            lines.append("")

        # Pasos diagnóstico
        if case.get("pasos_diagnostico"):
            lines.append("## Pasos de Diagnóstico\n")
            for paso in case["pasos_diagnostico"]:
                lines.append(f"1. {paso}")
            lines.append("")

        # Solución recomendada
        if case.get("solucion_recomendada"):
            lines.append("## Solución Recomendada\n")
            for solucion in case["solucion_recomendada"]:
                lines.append(f"- {solucion}")
            lines.append("")

        # Anti-patrón
        if case.get("anti_patron"):
            lines.append("## Anti-patrón\n")
            lines.append(f"{case['anti_patron']}\n")

        # Falso positivo
        if case.get("por_que_falso_positivo"):
            lines.append("## Por qué Falso Positivo\n")
            lines.append(f"{case['por_que_falso_positivo']}\n")

        # Riesgo e impacto
        lines.append("## Evaluación\n")
        lines.append(f"- **Confianza:** {case.get('nivel_confianza', 'medio')}")
        lines.append(f"- **Riesgo:** {case.get('nivel_riesgo', 'medio')}")
        lines.append(f"- **Impacto fiscal:** {case.get('impacto_fiscal_potencial', 'variable')}\n")

        # Señales tempranas
        if case.get("senales_tempranas"):
            lines.append("## Señales Tempranas\n")
            for senal in case["senales_tempranas"]:
                lines.append(f"- {senal}")
            lines.append("")

        # Validación antes/después
        if case.get("validacion_antes_despues"):
            lines.append("## Validación Antes/Después\n")
            vad = case["validacion_antes_despues"]
            if vad.get("antes"):
                lines.append("**Antes:**")
                for item in vad["antes"]:
                    lines.append(f"- {item}\n")
            if vad.get("despues"):
                lines.append("**Después:**")
                for item in vad["despues"]:
                    lines.append(f"- {item}\n")

        # Vigencia
        if case.get("vigencia"):
            lines.append("## Vigencia\n")
            v = case["vigencia"]
            if v.get("fecha_revision"):
                lines.append(f"- **Última revisión:** {v['fecha_revision']}")
            if v.get("motivo_caducidad_potencial"):
                lines.append(f"- **Riesgo de caducidad:** {v['motivo_caducidad_potencial']}")
            if v.get("fuente_recomendada_para_revalidar"):
                lines.append(f"- **Fuente para revalidar:** {v['fuente_recomendada_para_revalidar']}\n")

        # Guardar archivo
        content = "\n".join(lines)

        try:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"[OK] {filename}")
        except Exception as e:
            print(f"[ERROR] {filename}: {e}")

    print(f"\n[DONE] {len(cases)} casos migrados a knowledge/cases/")

if __name__ == "__main__":
    migrate_cases()
