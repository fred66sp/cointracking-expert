#!/usr/bin/env python3
"""
Rastreador de métricas de ahorros de caché.

Implementa Fase 6: Dashboard de caché — visualizar ahorros en tiempo real.

Uso:
  metrics = CacheMetrics('agp2025')
  metrics.record_cache_hit('get_trades', tokens_saved=2835, age_hours=2.5)
  print(metrics.get_report('session'))
"""

import json
from pathlib import Path
from datetime import datetime, timedelta
from typing import Dict, Any, List
from collections import defaultdict


class CacheMetrics:
    """Rastrear y reportar métricas de ahorros de caché."""

    def __init__(self, project_name: str, cache_dir: str = '.cache/cointracking'):
        self.project_name = project_name
        self.metrics_dir = Path(cache_dir) / project_name
        self.metrics_dir.mkdir(parents=True, exist_ok=True)
        self.metrics_file = self.metrics_dir / 'metrics.json'

        self.metrics = self._load_metrics()

    def _load_metrics(self) -> Dict:
        """Cargar metrics desde disco."""
        if self.metrics_file.exists():
            try:
                return json.loads(self.metrics_file.read_text())
            except Exception:
                pass

        # Crear nuevo si no existe
        return {
            'project': self.project_name,
            'created': datetime.now().isoformat(),
            'session': {'operations': []},
            'daily': {},
            'weekly': {},
            'monthly': {},
            'lifetime': {
                'hits': 0,
                'misses': 0,
                'tokens_saved': 0,
                'tokens_cost': 0,
                'operations': 0
            }
        }

    def _save_metrics(self):
        """Guardar metrics a disco."""
        self.metrics_file.write_text(json.dumps(self.metrics, indent=2))

    def record_cache_hit(self, call_name: str, tokens_saved: int, age_hours: float):
        """Registrar un cache hit (ahorró tokens)."""
        operation = {
            'call': call_name,
            'type': 'CACHE_HIT',
            'tokens_saved': tokens_saved,
            'age_hours': round(age_hours, 2),
            'timestamp': datetime.now().isoformat()
        }

        # Agregar a sesión
        self.metrics['session']['operations'].append(operation)

        # Actualizar agregados
        self._update_aggregates(tokens_saved=tokens_saved, is_hit=True)
        self._save_metrics()

        print(f"[METRIC] CACHE HIT: {call_name} ({tokens_saved} tokens ahorrados)")

    def record_mcp_call(self, call_name: str, tokens_cost: int):
        """Registrar una llamada MCP (miss, costó tokens)."""
        operation = {
            'call': call_name,
            'type': 'MCP_MISS',
            'tokens_cost': tokens_cost,
            'timestamp': datetime.now().isoformat()
        }

        # Agregar a sesión
        self.metrics['session']['operations'].append(operation)

        # Actualizar agregados
        self._update_aggregates(tokens_cost=tokens_cost, is_hit=False)
        self._save_metrics()

        print(f"[METRIC] MCP CALL: {call_name} ({tokens_cost} tokens gastados)")

    def _update_aggregates(self, tokens_saved: int = 0, tokens_cost: int = 0, is_hit: bool = False):
        """Actualizar estadísticas agregadas (día, semana, mes, vida)."""
        now = datetime.now()
        today = now.date().isoformat()
        week_start = (now - timedelta(days=now.weekday())).date().isoformat()
        month = now.strftime('%Y-%m')

        # Lifetime
        self.metrics['lifetime']['operations'] += 1
        if is_hit:
            self.metrics['lifetime']['hits'] += 1
            self.metrics['lifetime']['tokens_saved'] += tokens_saved
        else:
            self.metrics['lifetime']['misses'] += 1
            self.metrics['lifetime']['tokens_cost'] += tokens_cost

        # Daily
        if today not in self.metrics['daily']:
            self.metrics['daily'][today] = {'hits': 0, 'misses': 0, 'tokens_saved': 0, 'tokens_cost': 0}
        self.metrics['daily'][today]['tokens_saved'] += tokens_saved
        self.metrics['daily'][today]['tokens_cost'] += tokens_cost
        if is_hit:
            self.metrics['daily'][today]['hits'] += 1
        else:
            self.metrics['daily'][today]['misses'] += 1

        # Weekly (aproximado)
        if week_start not in self.metrics['weekly']:
            self.metrics['weekly'][week_start] = {'hits': 0, 'misses': 0, 'tokens_saved': 0, 'tokens_cost': 0}
        self.metrics['weekly'][week_start]['tokens_saved'] += tokens_saved
        self.metrics['weekly'][week_start]['tokens_cost'] += tokens_cost
        if is_hit:
            self.metrics['weekly'][week_start]['hits'] += 1
        else:
            self.metrics['weekly'][week_start]['misses'] += 1

        # Monthly
        if month not in self.metrics['monthly']:
            self.metrics['monthly'][month] = {'hits': 0, 'misses': 0, 'tokens_saved': 0, 'tokens_cost': 0}
        self.metrics['monthly'][month]['tokens_saved'] += tokens_saved
        self.metrics['monthly'][month]['tokens_cost'] += tokens_cost
        if is_hit:
            self.metrics['monthly'][month]['hits'] += 1
        else:
            self.metrics['monthly'][month]['misses'] += 1

    def get_stats(self, period: str = 'lifetime') -> Dict[str, Any]:
        """
        Obtener estadísticas para un período.

        period: 'session', 'today', 'week', 'month', 'lifetime'
        """
        if period == 'session':
            ops = self.metrics['session']['operations']
            hits = sum(1 for op in ops if op['type'] == 'CACHE_HIT')
            misses = sum(1 for op in ops if op['type'] == 'MCP_MISS')
            saved = sum(op.get('tokens_saved', 0) for op in ops if op['type'] == 'CACHE_HIT')
            cost = sum(op.get('tokens_cost', 0) for op in ops if op['type'] == 'MCP_MISS')

            return {
                'period': 'session',
                'hits': hits,
                'misses': misses,
                'total_operations': hits + misses,
                'tokens_saved': saved,
                'tokens_cost': cost,
                'net': saved - cost,
                'hit_rate': round((hits / (hits + misses) * 100), 1) if (hits + misses) > 0 else 0
            }

        elif period == 'today':
            today = datetime.now().date().isoformat()
            daily_stats = self.metrics['daily'].get(today, {})
            return {
                'period': 'today',
                'hits': daily_stats.get('hits', 0),
                'misses': daily_stats.get('misses', 0),
                'total_operations': daily_stats.get('hits', 0) + daily_stats.get('misses', 0),
                'tokens_saved': daily_stats.get('tokens_saved', 0),
                'tokens_cost': daily_stats.get('tokens_cost', 0),
            }

        elif period == 'week':
            week_start = (datetime.now() - timedelta(days=datetime.now().weekday())).date().isoformat()
            weekly_stats = self.metrics['weekly'].get(week_start, {})
            return {
                'period': f'week_of_{week_start}',
                'hits': weekly_stats.get('hits', 0),
                'misses': weekly_stats.get('misses', 0),
                'total_operations': weekly_stats.get('hits', 0) + weekly_stats.get('misses', 0),
                'tokens_saved': weekly_stats.get('tokens_saved', 0),
                'tokens_cost': weekly_stats.get('tokens_cost', 0),
            }

        elif period == 'month':
            month = datetime.now().strftime('%Y-%m')
            monthly_stats = self.metrics['monthly'].get(month, {})
            return {
                'period': month,
                'hits': monthly_stats.get('hits', 0),
                'misses': monthly_stats.get('misses', 0),
                'total_operations': monthly_stats.get('hits', 0) + monthly_stats.get('misses', 0),
                'tokens_saved': monthly_stats.get('tokens_saved', 0),
                'tokens_cost': monthly_stats.get('tokens_cost', 0),
            }

        else:  # lifetime
            return {
                'period': 'lifetime',
                'hits': self.metrics['lifetime']['hits'],
                'misses': self.metrics['lifetime']['misses'],
                'total_operations': self.metrics['lifetime']['operations'],
                'tokens_saved': self.metrics['lifetime']['tokens_saved'],
                'tokens_cost': self.metrics['lifetime']['tokens_cost'],
            }

    def get_report(self, period: str = 'session') -> str:
        """Generar reporte formateado para CLI."""
        stats = self.get_stats(period)

        lines = [
            "\n" + "=" * 70,
            f"ESTADISTICAS DE CACHE - {stats['period'].upper()}",
            "=" * 70,
            ""
        ]

        # Operaciones
        lines.append(f"Operaciones:")
        lines.append(f"  Hits (caché reutilizado):  {stats['hits']:>5}")
        lines.append(f"  Misses (MCP llamado):      {stats['misses']:>5}")
        lines.append(f"  Total:                     {stats['total_operations']:>5}")

        if stats['total_operations'] > 0:
            hit_rate = round((stats['hits'] / stats['total_operations'] * 100), 1)
            lines.append(f"  Hit Rate:                  {hit_rate:>4}%")

        lines.append("")

        # Tokens
        lines.append(f"Tokens:")
        lines.append(f"  Ahorrados (hits):          {stats['tokens_saved']:>6}")
        lines.append(f"  Gastados (misses):         {stats['tokens_cost']:>6}")
        net = stats['tokens_saved'] - stats['tokens_cost']
        lines.append(f"  Neto:                      {net:>6}")

        if stats['tokens_cost'] > 0:
            savings_pct = round((stats['tokens_saved'] / (stats['tokens_saved'] + stats['tokens_cost']) * 100), 1)
            lines.append(f"  Ahorro %:                  {savings_pct:>4}%")

        lines.append("")
        lines.append("=" * 70 + "\n")

        return "\n".join(lines)

    def get_breakdown_by_call(self) -> Dict[str, Dict]:
        """Desglose de ahorros por llamada (qué call más ahorra)."""
        breakdown = defaultdict(lambda: {'hits': 0, 'misses': 0, 'tokens_saved': 0, 'tokens_cost': 0})

        for op in self.metrics['session']['operations']:
            call = op['call']
            if op['type'] == 'CACHE_HIT':
                breakdown[call]['hits'] += 1
                breakdown[call]['tokens_saved'] += op.get('tokens_saved', 0)
            else:
                breakdown[call]['misses'] += 1
                breakdown[call]['tokens_cost'] += op.get('tokens_cost', 0)

        return dict(breakdown)

    def get_report_detailed(self) -> str:
        """Reporte detallado con desglose por call."""
        lines = [self.get_report('session')]

        breakdown = self.get_breakdown_by_call()
        if breakdown:
            lines.append("Desglose por Llamada:")
            lines.append("-" * 70)

            for call, stats in sorted(breakdown.items(), key=lambda x: x[1]['tokens_saved'], reverse=True):
                lines.append(f"\n{call}:")
                lines.append(f"  Hits:    {stats['hits']}")
                lines.append(f"  Misses:  {stats['misses']}")
                lines.append(f"  Ahorrados: {stats['tokens_saved']}")

        return "\n".join(lines)


# Ejemplo de uso
if __name__ == '__main__':
    metrics = CacheMetrics('test_project')

    # Simular hits y misses
    metrics.record_cache_hit('get_trades', tokens_saved=2835, age_hours=2.5)
    metrics.record_cache_hit('get_grouped_balance', tokens_saved=500, age_hours=0.1)
    metrics.record_cache_hit('get_gains', tokens_saved=1000, age_hours=1.5)
    metrics.record_mcp_call('get_historical_summary', tokens_cost=400)

    # Mostrar reportes
    print(metrics.get_report('session'))
    print(metrics.get_report_detailed())
