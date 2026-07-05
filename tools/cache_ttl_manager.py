#!/usr/bin/env python3
"""
Gestor de caché con TTL dinámico por tipo de dato.

Implementa Fase 5: Diferentes tiempos de vida para distintos tipos de datos.

Niveles de TTL:

  - Trades (histórico):     Permanente* (no cambia salvo reimportación)
  - Holdings:               15 minutos (estado vivo)
  - Balance:                15 minutos (estado vivo)
  - Tax Report (anual):     24 horas
  - Gains (FIFO):           Permanente* (determinista si trades no cambian)
  - Knowledge base:         Versionado (cambios explícitos)

*Permanente = hasta detectar cambios en origen (por versioning o reimportación)
"""

import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).parent))

from cache_manager import CacheManager
from cache_metrics import CacheMetrics
from typing import Any, Dict, Optional


class CacheTTLManager(CacheManager):
    """Extiende CacheManager con TTL dinámico por tipo de dato + métricas (Fase 6)."""

    def __init__(self, project_name: str, cache_dir: str = '.cache/cointracking'):
        """Inicializar con TTL dinámico y rastreador de métricas."""
        super().__init__(project_name, cache_dir)
        # Fase 6: Rastreador de métricas
        self.metrics = CacheMetrics(project_name, cache_dir)

    # Definir TTL por tipo de llamada (en horas)
    TTL_STRATEGIES = {
        # Datos históricos (no cambian salvo reimportación)
        'get_trades': {
            'ttl_hours': 999999,  # "Permanente"
            'invalidate_on': ['user_import', 'version_change'],
            'description': 'Trades históricos'
        },

        # Estado actual (cambia entre operaciones)
        'get_grouped_balance': {
            'ttl_hours': 0.25,  # 15 min
            'invalidate_on': ['user_operation', 'time_based'],
            'description': 'Balance actual (por activo)'
        },

        'get_balance': {
            'ttl_hours': 0.25,  # 15 min
            'invalidate_on': ['user_operation', 'time_based'],
            'description': 'Balance actual'
        },

        # Datos calculados (FIFO, depende de trades)
        'get_gains': {
            'ttl_hours': 999999,  # Permanente si trades no cambian
            'invalidate_on': ['trades_change', 'version_change'],
            'description': 'Ganancias FIFO (determinista)'
        },

        # Histórico de precios/saldos
        'get_historical_summary': {
            'ttl_hours': 24,  # Un día
            'invalidate_on': ['time_based', 'version_change'],
            'description': 'Histórico de balance/precios'
        },

        'get_historical_currency': {
            'ttl_hours': 24,
            'invalidate_on': ['time_based', 'version_change'],
            'description': 'Histórico de conversión EUR'
        },

        # Informe fiscal (anual)
        'get_tax_report': {
            'ttl_hours': 24,
            'invalidate_on': ['time_based', 'version_change'],
            'description': 'Informe de impuestos'
        }
    }

    def get_ttl_for_call(self, call_name: str) -> Dict[str, Any]:
        """
        Obtener configuración de TTL para una llamada.

        Returns:
            {'ttl_hours': 24, 'invalidate_on': [...], 'description': '...'}
        """
        if call_name in self.TTL_STRATEGIES:
            return self.TTL_STRATEGIES[call_name]

        # Default: 24 horas (ser conservador)
        return {
            'ttl_hours': 24,
            'invalidate_on': ['time_based', 'version_change'],
            'description': f'Default (24h)'
        }

    def get_or_fetch_dynamic(
        self,
        call_name: str,
        params: Dict,
        mcp_call_fn,
        force_refresh: bool = False,
        track_metrics: bool = True
    ) -> Any:
        """
        Get or fetch con TTL dinámico según tipo de dato.

        (Fase 5: TTL dinámico + Fase 6: Métricas)

        Args:
            call_name: Nombre de la llamada (ej. 'get_trades')
            params: Parámetros
            mcp_call_fn: Función para llamar MCP
            force_refresh: Forzar refetch
            track_metrics: Registrar hits/misses en métricas (Fase 6)

        Returns:
            Datos del caché o MCP
        """
        ttl_config = self.get_ttl_for_call(call_name)
        ttl_hours = ttl_config['ttl_hours']
        description = ttl_config['description']

        # Determinar si va a ser cache hit o miss
        cache_key = self._cache_key(call_name, params)

        # Verificar si está en caché (antes de llamar)
        if not force_refresh:
            cached_data, age_hours = self._load_cache(cache_key)
            if cached_data is not None and age_hours < ttl_hours:
                # CACHE HIT - registrar en métricas
                if track_metrics:
                    tokens_saved = self._estimate_tokens_for_call(call_name)
                    self.metrics.record_cache_hit(call_name, tokens_saved=tokens_saved, age_hours=age_hours)

                return cached_data

        # CACHE MISS - será una llamada MCP
        result = self.get_or_fetch_with_version_check(
            call_name=call_name,
            params=params,
            mcp_call_fn=mcp_call_fn,
            max_age_hours=ttl_hours,
            force_refresh=force_refresh
        )

        # Registrar en métricas
        if track_metrics:
            tokens_cost = self._estimate_tokens_for_call(call_name)
            self.metrics.record_mcp_call(call_name, tokens_cost=tokens_cost)

        return result

    def _estimate_tokens_for_call(self, call_name: str) -> int:
        """Estimar tokens para una llamada (desde TOKEN_BENCHMARKS)."""
        # Estimaciones desde docs/performance/TOKEN_BENCHMARKS.md
        estimates = {
            'get_trades': 2835,
            'get_grouped_balance': 500,
            'get_balance': 300,
            'get_gains': 1000,
            'get_historical_summary': 400,
            'get_historical_currency': 400,
            'get_tax_report': 800,
        }
        return estimates.get(call_name, 500)  # Default 500 si no está en lista

    def explain_ttl_strategy(self) -> str:
        """Explicar la estrategia de TTL completa."""
        lines = ["=== Estrategia de TTL Dinámico ===\n"]

        for call_name, config in self.TTL_STRATEGIES.items():
            ttl = config['ttl_hours']
            if ttl >= 999999:
                ttl_str = "Permanente*"
            elif ttl < 1:
                ttl_str = f"{int(ttl * 60)} minutos"
            else:
                ttl_str = f"{ttl}h"

            invalidate = ", ".join(config['invalidate_on'])
            lines.append(
                f"{call_name:.<30} {ttl_str:>12} "
                f"(invalida por: {invalidate})"
            )

        lines.append("\n* Permanente = hasta detectar cambios por versionado o reimportación")
        return "\n".join(lines)

    def report_cache_strategy(self) -> Dict[str, Any]:
        """
        Generar reporte de estrategia de caché para documentación.

        Útil para validar que TTLs son sensatos.
        """
        report = {
            'cache_strategy': 'Dynamic TTL per data type (Fase 5)',
            'strategies': self.TTL_STRATEGIES,
            'current_versions': self.current_versions,
            'cache_location': str(self.cache_dir),
            'manifest_path': str(self.manifest_path)
        }

        return report


# Ejemplo de uso
if __name__ == '__main__':
    mgr = CacheTTLManager('test_project')

    print(mgr.explain_ttl_strategy())

    # Simular uso
    print("\n\n=== Simulación: Usar caché con TTL dinámico ===\n")

    def mock_mcp(call_name: str, params: Dict):
        return {'data': 'ejemplo', 'call': call_name}

    # Trades: TTL permanente (hasta cambios)
    trades = mgr.get_or_fetch_dynamic('get_trades', {}, mock_mcp)
    print("Trades cacheados (permanente hasta reimportación)")

    # Balance: TTL 15 min
    balance = mgr.get_or_fetch_dynamic('get_grouped_balance', {}, mock_mcp)
    print("Balance cacheado (15 min)")

    # Stats
    stats = mgr.stats()
    print(f"\nCaché stats: {stats['total_entries']} entradas, {stats['total_size_kb']} KB")
