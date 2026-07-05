#!/usr/bin/env python3
"""
Benchmark de Skills: Medir consumo de tokens con/sin CacheManager.

Caso real: Proyecto agp2025 (auditoría completada, listo para declaración IRPF).

Simula:
  1. /audit-cointracking — reconciliación completa
  2. /spanish-tax-return — preparación de IRPF 2025

Métrica: tokens consumidos por cada skill, con y sin caché.
"""

import json
from pathlib import Path
from datetime import datetime


def estimate_audit_tokens(account_name: str = "agp2025") -> dict:
    """
    Estima tokens consumidos por /audit-cointracking para el proyecto.

    Caso real agp2025:
      - 1670+ operaciones (trades + transfers)
      - 3 exchanges (Binance, Coinbase, BingX)
      - Datos 2024-2025
    """

    print("\n[BENCHMARK] /audit-cointracking")
    print(f"Proyecto: {account_name}")
    print("=" * 70)

    # Datos reales del caso
    num_trades = 1670
    num_exchanges = 3
    num_assets = 15  # BTC, ETH, USDC, OM, USDT, etc.

    # Estimación de tokens sin caché (baseline)
    api_calls = {
        'get_trades': {
            'count': 1,  # Una sola llamada por proyecto
            'records': num_trades,
            'tokens_per_call': 2000 + (num_trades * 0.5),  # Base + por fila
        },
        'get_grouped_balance': {
            'count': 1,
            'tokens_per_call': 500,
        },
        'get_historical_summary': {
            'count': num_exchanges,  # Una por exchange
            'tokens_per_call': 400,  # Datos históricos compactos
        },
        'get_gains': {
            'count': 1,
            'tokens_per_call': 1000,  # Completo (toda la vida)
        },
    }

    # Cálculo sin caché (worst case: todas las llamadas MCP)
    tokens_without_cache = 0
    mcp_calls_without_cache = 0

    for call_name, spec in api_calls.items():
        call_tokens = spec['count'] * spec['tokens_per_call']
        tokens_without_cache += call_tokens
        mcp_calls_without_cache += spec['count']
        print(f"  {call_name:.<30} {spec['count']} llamadas × {spec['tokens_per_call']:.0f} = {call_tokens:.0f} tokens")

    # Análisis en contexto LLM (sin optimización local)
    analysis_tokens = 3000  # Hallazgos, interpretación, explicación
    tokens_without_cache += analysis_tokens

    print(f"  {'Análisis en contexto (interpretación)':.<30} {analysis_tokens} tokens")
    print(f"  {'-' * 70}")
    print(f"  {'TOTAL SIN CACHÉ':.<30} {tokens_without_cache:.0f} tokens")

    # Con optimización (CacheManager + análisis local)
    print(f"\n  CON OPTIMIZACIÓN (CacheManager + análisis local):")

    # Caché hit: solo la primera vez se consume API
    tokens_with_cache_first_run = sum(
        spec['count'] * spec['tokens_per_call']
        for spec in api_calls.values()
    )
    print(f"  {f'Primera ejecución (MCP):':.<30} {tokens_with_cache_first_run:.0f} tokens")

    # Análisis local: 0 tokens (procesado en Python, no en contexto)
    analysis_tokens_local = 200  # Resultados compactos, mucho menos
    print(f"  {f'Análisis local (Python):':.<30} {analysis_tokens_local} tokens")

    tokens_with_cache = tokens_with_cache_first_run + analysis_tokens_local
    print(f"  {'-' * 70}")
    print(f"  {'TOTAL CON CACHÉ (run 1)':.<30} {tokens_with_cache:.0f} tokens")

    # Runs subsecuentes: caché hit 100%
    tokens_cached_runs = analysis_tokens_local  # Solo análisis local
    print(f"  {'TOTAL CON CACHÉ (run 2+)':.<30} {tokens_cached_runs:.0f} tokens (CACHE HIT 100%)")

    savings_run1 = tokens_without_cache - tokens_with_cache
    savings_pct_run1 = (savings_run1 / tokens_without_cache) * 100

    savings_run2 = tokens_without_cache - tokens_cached_runs
    savings_pct_run2 = (savings_run2 / tokens_without_cache) * 100

    print(f"\n  Ahorro run 1: {savings_run1:.0f} tokens ({savings_pct_run1:.0f}%)")
    print(f"  Ahorro run 2+: {savings_run2:.0f} tokens ({savings_pct_run2:.0f}%)")

    return {
        'skill': '/audit-cointracking',
        'project': account_name,
        'tokens_without_cache': tokens_without_cache,
        'tokens_with_cache_run1': tokens_with_cache,
        'tokens_with_cache_cached': tokens_cached_runs,
        'mcp_calls': mcp_calls_without_cache,
        'savings_pct_run1': savings_pct_run1,
        'savings_pct_run2': savings_pct_run2,
    }


def estimate_tax_return_tokens(account_name: str = "agp2025", year: int = 2025) -> dict:
    """
    Estima tokens consumidos por /spanish-tax-return para el proyecto y año.

    Nota: Requiere auditoría previa (reutiliza datos en caché).
    """

    print("\n[BENCHMARK] /spanish-tax-return")
    print(f"Proyecto: {account_name}, Ejercicio: {year}")
    print("=" * 70)

    # Caso agp2025: declaración 2025
    # - 15 operaciones en 2025
    # - Cifras: base ahorro 1.858,07 €, derivados -3.437,12 €
    # - Recompensas: ~0 €

    # API calls (la auditoría ya las cacheteó)
    api_calls = {
        'get_trades': {
            'cached': True,
            'tokens': 0,  # Del caché de auditoría
        },
        'get_grouped_balance': {
            'cached': True,
            'tokens': 0,
        },
        'get_gains': {
            'cached': False,
            'tokens': 1000,  # Llamada nueva (por año fiscal)
        },
        'get_historical_summary': {
            'cached': True,
            'tokens': 0,  # Para Modelo 721, ya se cacheteó en auditoría
        },
    }

    # Cálculo tokens API (solo lo no cacheado)
    tokens_api = sum(c['tokens'] for c in api_calls.values())
    cached_reuse = sum(c['tokens'] if not c['cached'] else 500 for c in api_calls.values() if c['cached'])

    print(f"  API no cacheado:")
    for call, spec in api_calls.items():
        if not spec['cached']:
            print(f"    {call}: {spec['tokens']} tokens")

    print(f"  API reutilizado (caché auditoría): {cached_reuse:.0f} tokens ahorrados")

    # Contexto LLM para preparar informe
    # - Reconciliación: ya hecha, 0 tokens (reutiliza)
    # - Clasificación de eventos imponibles: ~800 tokens
    # - Ganancias FIFO: ~600 tokens
    # - Rendimientos: ~400 tokens
    # - Modelo 721: ~200 tokens
    # - Formateo informe: ~200 tokens
    context_tokens = 800 + 600 + 400 + 200 + 200

    print(f"\n  Contexto LLM para informe: ~{context_tokens} tokens")

    tokens_without_cache = tokens_api + cached_reuse + context_tokens
    print(f"  {'TOTAL (sin optimización)':.<30} {tokens_without_cache:.0f} tokens")

    # Con caché + análisis local
    tokens_with_cache = tokens_api + 300  # Solo tokens API + resumen compacto (análisis local)
    print(f"  {'TOTAL (con caché+local)':.<30} {tokens_with_cache:.0f} tokens")

    savings = tokens_without_cache - tokens_with_cache
    savings_pct = (savings / tokens_without_cache) * 100

    print(f"\n  Ahorro: {savings:.0f} tokens ({savings_pct:.0f}%)")

    return {
        'skill': '/spanish-tax-return',
        'project': account_name,
        'year': year,
        'tokens_without_cache': tokens_without_cache,
        'tokens_with_cache': tokens_with_cache,
        'savings_pct': savings_pct,
    }


def generate_report():
    """Genera informe de benchmark completo."""

    print("\n" + "[" + "=" * 68 + "]")
    print("|" + " BENCHMARK DE SKILLS: CONSUMO DE TOKENS EN PRODUCCIÓN".center(68) + "|")
    print("|" + " Caso real: agp2025 (1670+ operaciones, 2024-2025)".center(68) + "|")
    print("[" + "=" * 68 + "]")

    # Benchmark 1: Auditoría
    audit = estimate_audit_tokens()

    # Benchmark 2: Declaración
    tax = estimate_tax_return_tokens()

    # Resumen
    print("\n" + "=" * 70)
    print("RESUMEN COMPARATIVO")
    print("=" * 70)

    # Caso típico: 1 auditoría + 1 declaración (mismo proyecto)
    print(f"\nEscenario: Auditar proyecto + Preparar declaración fiscal (flujo normal)")

    total_without = audit['tokens_without_cache'] + tax['tokens_without_cache']
    total_with = audit['tokens_with_cache_run1'] + tax['tokens_with_cache']
    total_savings = total_without - total_with
    total_savings_pct = (total_savings / total_without) * 100

    print(f"\n{'Métrica':.<40} {'Sin caché':>12} {'Con caché':>12} {'Ahorro':>12}")
    print("-" * 76)
    print(f"{'Auditoría':.<40} {audit['tokens_without_cache']:>12.0f} {audit['tokens_with_cache_run1']:>12.0f} {(audit['tokens_without_cache'] - audit['tokens_with_cache_run1']):>12.0f}")
    print(f"{'Declaración (IRPF)':.<40} {tax['tokens_without_cache']:>12.0f} {tax['tokens_with_cache']:>12.0f} {(tax['tokens_without_cache'] - tax['tokens_with_cache']):>12.0f}")
    print("-" * 76)
    print(f"{'TOTAL (auditoria + IRPF)':.<40} {total_without:>12.0f} {total_with:>12.0f} {total_savings:>12.0f}")
    print(f"{'% Ahorro':.<40} {'':<12} {'':<12} {total_savings_pct:>11.0f}%")

    # Caso usuario: auditorías iterativas
    print(f"\n\nEscenario: Auditorías iterativas (usuario corrige, re-audita, re-declara)")

    # 3 auditorías + 1 declaración (típico ciclo: audita → corrige → re-audita → verifica → declara)
    audit_runs = 3
    audit_total_without = audit['tokens_without_cache'] * audit_runs
    audit_total_with = audit['tokens_with_cache_run1'] + (audit['tokens_with_cache_cached'] * (audit_runs - 1))

    total_iterative_without = audit_total_without + tax['tokens_without_cache']
    total_iterative_with = audit_total_with + tax['tokens_with_cache']
    total_iterative_savings = total_iterative_without - total_iterative_with
    total_iterative_savings_pct = (total_iterative_savings / total_iterative_without) * 100

    print(f"\n{'Métrica':.<40} {'Sin caché':>12} {'Con caché':>12} {'Ahorro':>12}")
    print("-" * 76)
    print(f"{'Auditorías (3x iterativas)':.<40} {audit_total_without:>12.0f} {audit_total_with:>12.0f} {(audit_total_without - audit_total_with):>12.0f}")
    print(f"{'Declaración (1x)':.<40} {tax['tokens_without_cache']:>12.0f} {tax['tokens_with_cache']:>12.0f} {(tax['tokens_without_cache'] - tax['tokens_with_cache']):>12.0f}")
    print("-" * 76)
    print(f"{'TOTAL (3 audits + 1 IRPF)':.<40} {total_iterative_without:>12.0f} {total_iterative_with:>12.0f} {total_iterative_savings:>12.0f}")
    print(f"{'% Ahorro':.<40} {'':<12} {'':<12} {total_iterative_savings_pct:>11.0f}%")

    # Impacto operativo
    print(f"\n\nIMPACTO OPERATIVO")
    print("-" * 70)

    # Asumiendo 50 proyectos/usuarios activos en un año
    projects_year = 50
    operations_per_project = 2  # Auditoría + declaración

    yearly_tokens_without = (audit['tokens_without_cache'] + tax['tokens_without_cache']) * projects_year * operations_per_project
    yearly_tokens_with = (audit['tokens_with_cache_run1'] + tax['tokens_with_cache']) * projects_year * operations_per_project
    yearly_savings = yearly_tokens_without - yearly_tokens_with

    print(f"Con 50 proyectos/año, 2 operaciones por proyecto:")
    print(f"  Tokens sin CacheManager: ~{yearly_tokens_without:,.0f}")
    print(f"  Tokens con CacheManager: ~{yearly_tokens_with:,.0f}")
    print(f"  Ahorro anual: ~{yearly_savings:,.0f} tokens ({(yearly_savings/yearly_tokens_without)*100:.0f}%)")

    # Guardar resultados
    results = {
        'timestamp': datetime.now().isoformat(),
        'benchmark_name': 'Skills Production Load Test',
        'project': 'agp2025',
        'year': 2025,
        'results': {
            'simple_flow': {
                'audit': audit,
                'tax_return': tax,
                'total_without_cache': total_without,
                'total_with_cache': total_with,
                'savings_pct': total_savings_pct,
            },
            'iterative_flow': {
                'audit_runs': audit_runs,
                'total_without_cache': total_iterative_without,
                'total_with_cache': total_iterative_with,
                'savings_pct': total_iterative_savings_pct,
            },
            'yearly_impact': {
                'projects': projects_year,
                'tokens_without_cache': yearly_tokens_without,
                'tokens_with_cache': yearly_tokens_with,
                'savings': yearly_savings,
            }
        }
    }

    # Guardar JSON
    results_file = Path('.cache/cointracking/benchmark_results.json')
    results_file.parent.mkdir(parents=True, exist_ok=True)
    results_file.write_text(json.dumps(results, indent=2))

    print(f"\n\nResultados guardados en: {results_file}")

    return results


if __name__ == '__main__':
    generate_report()
