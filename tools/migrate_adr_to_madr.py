#!/usr/bin/env python3
"""
Migra ADRs de DECISIONS.md al formato MADR en archivos individuales.
Uso: python tools/migrate_adr_to_madr.py
"""

import re
import os
from pathlib import Path

DECISIONS_FILE = Path(__file__).parent.parent / "DECISIONS.md"
ADR_DIR = Path(__file__).parent.parent / "adr"
ADR_DIR.mkdir(exist_ok=True)

def extract_adrs(content):
    """Extrae bloques ADR de DECISIONS.md."""
    adrs = []
    parts = content.split("## ADR-")

    for part in parts[1:]:  # Saltar la primera parte (header)
        lines = part.split("\n")

        # Extraer numero y titulo de la primera linea
        first_line = lines[0]
        match = re.match(r"(\d+):(.*)", first_line)
        if not match:
            continue

        adr_num = match.group(1).zfill(3)
        title = match.group(2).strip()

        # Buscar Estado y Fecha
        state = ""
        date = ""
        body_start = 0

        for i, line in enumerate(lines):
            if "**Estado:**" in line:
                state = line.split("**Estado:**")[1].strip()
            elif "**Fecha:**" in line:
                date = line.split("**Fecha:**")[1].strip()
            elif "**Contexto:**" in line:
                body_start = i
                break

        # Extraer body (desde Contexto hasta el siguiente ---  o fin)
        body_lines = []
        for i in range(body_start, len(lines)):
            if lines[i].startswith("---"):
                break
            body_lines.append(lines[i])

        body = "\n".join(body_lines).strip()

        if body:
            adrs.append({
                "num": adr_num,
                "title": title,
                "state": state,
                "date": date,
                "body": body
            })

    return adrs

def parse_adr_body(body):
    """Parsea el cuerpo del ADR en secciones."""
    sections = {}
    current_section = None
    current_content = []

    for line in body.split("\n"):
        if line.startswith("**Contexto:**"):
            if current_section:
                sections[current_section] = "\n".join(current_content).strip()
            current_section = "context"
            current_content = []
        elif line.startswith("**Opciones consideradas:**"):
            if current_section:
                sections[current_section] = "\n".join(current_content).strip()
            current_section = "options"
            current_content = []
        elif line.startswith("**Decision"):
            if current_section:
                sections[current_section] = "\n".join(current_content).strip()
            current_section = "decision"
            current_content = []
        elif line.startswith("**Consecuencias:**"):
            if current_section:
                sections[current_section] = "\n".join(current_content).strip()
            current_section = "consequences"
            current_content = []
        elif any(line.startswith(p) for p in ["**Notas", "**Proximos", "**Cambios", "**Materiales", "**Baja", "**Cuestion abierta"]):
            if current_section:
                sections[current_section] = "\n".join(current_content).strip()
            current_section = "notes"
            current_content = [line]
        else:
            current_content.append(line)

    if current_section:
        sections[current_section] = "\n".join(current_content).strip()

    return sections

def convert_to_madr(adr):
    """Convierte un ADR al formato MADR."""
    sections = parse_adr_body(adr["body"])

    status_map = {
        "Decidido": "Accepted",
        "Propuesto": "Proposed",
        "Pendiente": "Pending",
        "Rechazado": "Rejected",
    }

    status = status_map.get(adr["state"], adr["state"])

    madr = f"""# ADR-{adr['num']}: {adr['title']}

**Status:** {status}

**Date:** {adr['date']}

## Context

{sections.get('context', '[Context not found]')}

## Decision

{sections.get('decision', '[Decision not found]')}

## Consequences

{sections.get('consequences', '[Consequences not found]')}
"""

    if 'notes' in sections and sections['notes'].strip():
        madr += f"\n## Notes\n\n{sections['notes']}\n"

    return madr

def generate_readme(adrs):
    """Genera el README.md del directorio adr/."""
    readme = """# Registros de Decisiones Arquitectonicas (ADRs)

Este directorio contiene las decisiones arquitectonicas del proyecto, documentadas en formato MADR (Markdown Any Decision Record).

## Indice

"""

    for adr in adrs:
        readme += f"- [ADR-{adr['num']}: {adr['title']}](./0{adr['num']}-{adr['title'].lower().replace(' ', '-').replace('(', '').replace(')', '')}.md)\n"

    readme += "\n## Proceso de ADR\n\n"
    readme += "Toda decision arquitectonica importante debe:\n\n"
    readme += "1. Ser propuesta en una rama nueva con un ADR borrador\n"
    readme += "2. Ser discutida en revision de codigo\n"
    readme += "3. Ser revisada por arquitecto del proyecto\n"
    readme += "4. Ser aprobada por el equipo\n"
    readme += "5. Ser completada con la decision final\n\n"
    readme += "Las decisiones menores pueden ser documentadas informalmente en CHANGELOG.md.\n"

    return readme

def generate_filename(adr):
    """Genera el nombre de archivo para un ADR."""
    title = adr["title"].lower()
    title = re.sub(r'[^\w\s-]', '', title)
    title = re.sub(r'\s+', '-', title)
    title = title[:50].rstrip('-')
    return f"0{adr['num']}-{title}.md"

def main():
    print("Leyendo DECISIONS.md...")
    with open(DECISIONS_FILE, 'r', encoding='utf-8') as f:
        content = f.read()

    print("Extrayendo ADRs...")
    adrs = extract_adrs(content)
    print(f"   Encontrados {len(adrs)} ADRs")

    print("Convirtiendo a formato MADR...")
    for adr in adrs:
        madr_content = convert_to_madr(adr)
        filename = generate_filename(adr)
        filepath = ADR_DIR / filename

        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(madr_content)

        print(f"   OK ADR-{adr['num']}")

    print("Generando README.md...")
    readme = generate_readme(adrs)
    readme_path = ADR_DIR / "README.md"
    with open(readme_path, 'w', encoding='utf-8') as f:
        f.write(readme)
    print(f"   OK README.md")

    print(f"\nMigracion completada: {len(adrs)} ADRs en adr/")

if __name__ == "__main__":
    main()
