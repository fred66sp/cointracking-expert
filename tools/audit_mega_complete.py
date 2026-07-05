#!/usr/bin/env python3
"""
AUDITORÍA MEGA-EXHAUSTIVA — Una sola pasada, todas las dimensiones.

Valida TODO simultáneamente:
  1. YAML frontmatter válido
  2. IDs únicos (sin duplicados)
  3. Niveles correctos (A-F)
  4. Campos obligatorios
  5. Valid_until nunca null para A/B
  6. Fechas válidas
  7. Confidence/Authority valores válidos
  8. ADRs citados existen
  9. Referencias a docs resolubles
  10. Coherencia entre docs relacionados

Uso:
  python tools/audit_mega_complete.py
"""

import re
import sys
from pathlib import Path
from collections import defaultdict
import yaml

# Raíz del repo derivada de la ubicación del script (tools/ -> raíz), no
# cableada: con la ruta absoluta anterior ('h:/cointracking-expert') el
# workflow de CI en ubuntu-latest no encontraba ningún documento y "pasaba"
# en vacío sin validar nada (corregido 2026-07-05).
REPO_ROOT = Path(__file__).resolve().parent.parent

class MegaAudit:
    def __init__(self):
        self.errors = defaultdict(list)
        self.warnings = defaultdict(list)
        self.docs = {}
        self.ids_seen = {}
        self.adr_files = set()

    def extract_yaml(self, content: str) -> dict | None:
        """Extrae YAML frontmatter."""
        if not content.startswith('---'):
            return None
        match = re.match(r'^---\n(.*?)\n---\n', content, re.DOTALL)
        if not match:
            return None
        try:
            return yaml.safe_load(match.group(1))
        except:
            return None

    def validate_file(self, file_path: Path) -> None:
        """Valida UN archivo contra todos los criterios."""
        try:
            content = file_path.read_text(encoding='utf-8')
        except Exception as e:
            self.errors[str(file_path)].append(f"ERROR reading: {e}")
            return

        yaml_data = self.extract_yaml(content)
        fname = file_path.name
        frel = str(file_path.resolve().relative_to(REPO_ROOT))

        # 1. YAML válido?
        if not yaml_data:
            # Es OK para INDEX.md o templates
            if 'INDEX.md' in fname or '.metadata' in str(file_path):
                return
            self.errors[frel].append("ERROR: No valid YAML frontmatter")
            return

        # 2-5. Campos obligatorios y válidos
        doc_id = yaml_data.get('id')
        level = yaml_data.get('level')
        valid_from = yaml_data.get('valid_from')
        valid_until = yaml_data.get('valid_until')
        confidence = yaml_data.get('confidence')
        authority = yaml_data.get('authority')
        last_verified = yaml_data.get('last_verified')

        # Validar ID
        if not doc_id:
            self.errors[frel].append("MISSING: id")
        elif not re.match(r'^KB-[A-F][1-9]-\d{3}$', str(doc_id)):
            # Aceptar genéricos (KB-X-NNN) en .metadata/
            if not (re.match(r'^KB-[A-F]-NNN$', str(doc_id)) and '.metadata' in frel):
                self.errors[frel].append(f"INVALID id: {doc_id} (esperado KB-[A-F][1-9]-NNN)")
        else:
            # Verificar duplicados
            if doc_id in self.ids_seen:
                self.errors[frel].append(f"DUPLICATE id: {doc_id} (también en {self.ids_seen[doc_id]})")
            else:
                self.ids_seen[doc_id] = frel

        # Validar level
        if not level:
            self.errors[frel].append("MISSING: level")
        elif level not in ['A', 'B', 'C', 'D', 'E', 'F']:
            self.errors[frel].append(f"INVALID level: {level}")

        # Validar valid_until (CRÍTICO)
        if level in ['A', 'B']:
            if valid_until is None or valid_until == 'null' or valid_until == '':
                self.errors[frel].append(f"CRITICAL: valid_until null/empty (Nivel {level})")

        # Validar fechas
        if valid_from:
            if not re.match(r'^\d{4}-\d{2}-\d{2}$', str(valid_from)):
                self.errors[frel].append(f"INVALID valid_from format: {valid_from}")

        if valid_until and valid_until not in ['null', '']:
            if not re.match(r'^\d{4}-\d{2}-\d{2}$', str(valid_until)):
                self.errors[frel].append(f"INVALID valid_until format: {valid_until}")

        # Validar confidence
        if confidence and confidence not in ['high', 'medium', 'low']:
            self.warnings[frel].append(f"INVALID confidence: {confidence}")

        # Validar authority
        if authority and authority not in ['official', 'verified', 'empirical', 'reference']:
            self.warnings[frel].append(f"INVALID authority: {authority}")

        # Validar last_verified (debe estar presente para A/B)
        if level in ['A', 'B']:
            if not last_verified:
                self.warnings[frel].append("MISSING: last_verified (A/B deberían tener)")

        # Almacenar para validaciones posteriores
        self.docs[frel] = {'yaml': yaml_data, 'content': content, 'path': file_path}

    def validate_references(self) -> None:
        """Valida referencias entre documentos."""
        for frel, doc in self.docs.items():
            content = doc['content']

            # Buscar ADR-XXX
            adr_refs = re.findall(r'ADR-(\d+)', content)
            for adr_num in adr_refs:
                adr_num_int = int(adr_num)
                # Verificar que existe algún archivo que matchee
                found = False
                for p in (REPO_ROOT / 'adr').glob(f'0{adr_num_int:03d}-*'):
                    found = True
                    break
                if not found:
                    self.errors[frel].append(f"BROKEN REF: ADR-{adr_num} (no existe)")

            # Buscar [documento](path) links internos
            doc_links = re.findall(r'\[([^\]]+)\]\(([^)]+)\)', content)
            for text, path in doc_links:
                if path.startswith('http') or path.startswith('#'):
                    continue  # URLs externas y anchors son OK
                path_no_anchor = path.split('#', 1)[0]
                if not path_no_anchor:
                    continue  # link tipo "#seccion" en el mismo archivo
                # Los links relativos en Markdown se resuelven contra el
                # directorio del archivo que contiene el link, no contra la
                # raíz del repo (bug encontrado 2026-07-05: generaba ~100
                # falsos positivos "BROKEN LINK" en docs de navegación).
                doc_dir = REPO_ROOT / Path(frel).parent
                full_path = (doc_dir / path_no_anchor).resolve()
                if not full_path.exists():
                    self.warnings[frel].append(f"BROKEN LINK: {path}")

    def validate_coherence(self) -> None:
        """Valida coherencia entre docs relacionados."""
        # Ejemplo: CAPITAL_GAINS y CAPITAL_INCOME deben referenciar plazos
        capital_gains = next((d for d in self.docs if 'CAPITAL_GAINS' in d), None)
        capital_income = next((d for d in self.docs if 'CAPITAL_INCOME' in d), None)

        if capital_gains and capital_income:
            gains_yaml = self.docs[capital_gains]['yaml']
            income_yaml = self.docs[capital_income]['yaml']

            # Deben estar en el mismo nivel y con valid_until similar
            if gains_yaml.get('valid_until') != income_yaml.get('valid_until'):
                self.warnings[capital_gains].append(
                    f"COHERENCE: valid_until ({gains_yaml.get('valid_until')}) != "
                    f"CAPITAL_INCOME ({income_yaml.get('valid_until')})"
                )

    def run(self) -> None:
        """Ejecuta auditoría completa."""
        print("=" * 80)
        print("AUDITORÍA MEGA-EXHAUSTIVA")
        print("=" * 80)
        print()

        # Fase 1: Validar todos los archivos
        print("[1/3] Validando archivos...")
        knowledge_dir = (REPO_ROOT / 'knowledge')
        md_files = list(knowledge_dir.rglob('*.md'))

        for file_path in sorted(md_files):
            self.validate_file(file_path)

        print(f"  [OK] {len(md_files)} archivos procesados")

        # Fase 2: Validar referencias
        print("[2/3] Validando referencias cruzadas...")
        self.validate_references()
        print(f"  [OK] Referencias validadas")

        # Fase 3: Validar coherencia
        print("[3/3] Validando coherencia...")
        self.validate_coherence()
        print(f"  [OK] Coherencia validada")
        print()

        # Reporte
        print("=" * 80)
        print("RESULTADOS")
        print("=" * 80)

        total_errors = sum(len(v) for v in self.errors.values())
        total_warnings = sum(len(v) for v in self.warnings.values())

        print(f"Errores críticos: {total_errors}")
        print(f"Warnings: {total_warnings}")
        print()

        if total_errors > 0:
            print("ERRORES CRÍTICOS:")
            print("-" * 80)
            for frel, errs in sorted(self.errors.items()):
                if errs:
                    print(f"\n{frel}:")
                    for err in errs:
                        print(f"  [ERROR] {err}")

        if total_warnings > 0:
            print("\n\nWARNINGS:")
            print("-" * 80)
            for frel, warns in sorted(self.warnings.items()):
                if warns:
                    print(f"\n{frel}:")
                    for warn in warns:
                        print(f"  [WARN] {warn}")

        print()
        print("=" * 80)
        if total_errors == 0:
            print("[OK] SISTEMA LIMPIO - 0 ERRORES CRITICOS")
        else:
            print(f"[FAIL] {total_errors} ERRORES CRITICOS ENCONTRADOS")
        print("=" * 80)

if __name__ == '__main__':
    audit = MegaAudit()
    audit.run()
