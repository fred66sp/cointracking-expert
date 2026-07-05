#!/usr/bin/env python3
"""
CLI para mostrar estadísticas de caché.

Equivalente a: rtk gain --cache

Uso:
  python cache_cli.py [proyecto] [periodo]

Ejemplos:
  python cache_cli.py agp2025 session
  python cache_cli.py agp2025 lifetime
  python cache_cli.py agp2025 detailed
"""

import sys
import argparse
from pathlib import Path

# Agregar tools a path
sys.path.insert(0, str(Path(__file__).parent))

from cache_metrics import CacheMetrics


def main():
    parser = argparse.ArgumentParser(
        description='Mostrar estadísticas de caché - Fase 6'
    )
    parser.add_argument('project', help='Nombre del proyecto (ej. agp2025)')
    parser.add_argument(
        'period',
        nargs='?',
        default='session',
        choices=['session', 'today', 'week', 'month', 'lifetime', 'detailed'],
        help='Período a reportar (default: session)'
    )

    args = parser.parse_args()

    try:
        metrics = CacheMetrics(args.project)

        if args.period == 'detailed':
            print(metrics.get_report_detailed())
        else:
            print(metrics.get_report(args.period))

    except FileNotFoundError:
        print(f"Error: No hay datos de caché para proyecto '{args.project}'")
        print("Ejecuta una auditoría primero: /audit-cointracking")
        sys.exit(1)

    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)


if __name__ == '__main__':
    main()
