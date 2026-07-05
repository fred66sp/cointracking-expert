#!/usr/bin/env python3
"""
Gestor de caché persistente para MCP de CoinTracking.

Implementa ADR-039: Optimización de tokens mediante caché local y procesamiento
en Python en lugar de contexto LLM.

FILOSOFÍA:
  - Reutilizar datos MCP que no han cambiado (caché)
  - Procesar localmente (Python), no en contexto LLM
  - Pasar al LLM solo hallazgos resumidos, nunca JSON crudo

EJEMPLO DE USO EN SKILL:

    from tools.cache_manager import CacheManager

    # Inicializar para proyecto activo
    mgr = CacheManager('agp2025')

    # Get trades (primera vez = MCP call, después = caché)
    trades = mgr.get_or_fetch(
        'get_trades',
        {'limit': None, 'start': unix_ts, 'end': unix_ts},
        mcp_call_fn=lambda call, params: mcp.cointracking_get_trades(**params),
        max_age_hours=24
    )

    # Procesar localmente (no en contexto LLM)
    duplicates = detect_duplicates(trades)  # Python puro
    orphans = detect_orphan_transfers(trades)
    summary = f"Encontrado: {len(duplicates)} duplicados, {len(orphans)} huérfanas"

    # Pasar al LLM solo el resumen (~200 tokens)
    # NO pasar trades completo (~3000 tokens)

INVALIDACIÓN:

    # Si usuario cambió datos en CoinTracking
    mgr.invalidate_all()  # Borrar TODO

    # Si solo queremos refrescar trades
    mgr.invalidate_pattern('get_trades')

ESTADÍSTICAS:

    stats = mgr.stats()
    print(f"Caché: {stats['total_entries']} entradas, {stats['total_size_kb']} KB")

REFERENCIA:
  - ADR-039: adr/0039-optimizacion-tokens-y-cache.md
  - Benchmark: docs/performance/TOKEN_BENCHMARKS.md
  - Roadmap: implementation/CACHE_ROADMAP.md
"""

import json
import hashlib
import time
from pathlib import Path
from datetime import datetime, timedelta
from typing import Any, Dict, Optional

# Importar version tracker
import sys
sys.path.insert(0, str(Path(__file__).parent))
from version_tracker import VersionTracker

class CacheManager:
    def __init__(self, project_name: str, cache_dir: str = '.cache/cointracking'):
        """Inicializa gestor de caché para un proyecto."""
        self.project_name = project_name
        self.cache_dir = Path(cache_dir) / project_name
        self.cache_dir.mkdir(parents=True, exist_ok=True)
        self.manifest_path = self.cache_dir / 'manifest.json'
        self.manifest = self._load_manifest()

        # Version tracking (Fase 4)
        self.version_tracker = VersionTracker()
        self.current_versions = self.version_tracker.get_current_versions()

    def _load_manifest(self) -> Dict:
        """Carga manifest de caché (metadata de archivos)."""
        if self.manifest_path.exists():
            try:
                return json.loads(self.manifest_path.read_text())
            except:
                return {}
        return {}

    def _save_manifest(self) -> None:
        """Guarda manifest actualizado."""
        self.manifest_path.write_text(json.dumps(self.manifest, indent=2))

    def _cache_key(self, call_name: str, params: Dict) -> str:
        """Genera clave única para una llamada + parámetros."""
        params_str = json.dumps(params, sort_keys=True)
        params_hash = hashlib.md5(params_str.encode()).hexdigest()[:8]
        return f"{call_name}_{params_hash}"

    def _cache_file(self, cache_key: str) -> Path:
        """Ruta del archivo de caché para una clave."""
        return self.cache_dir / f"{cache_key}.json"

    def get_or_fetch(
        self,
        call_name: str,
        params: Dict,
        mcp_call_fn,
        max_age_hours: int = 24,
        force_refresh: bool = False
    ) -> Any:
        """
        Obtiene datos de caché o llama MCP si no existe/está viejo.

        Args:
            call_name: Nombre de la función MCP (ej. 'get_balance')
            params: Parámetros de la llamada
            mcp_call_fn: Función que ejecuta la llamada MCP
            max_age_hours: Máximo edad del caché en horas (default 24)
            force_refresh: Si True, ignora caché y refetch

        Returns:
            Datos (del caché o MCP)
        """
        cache_key = self._cache_key(call_name, params)

        # Intentar cargar de caché
        if not force_refresh:
            cached_data, age_hours = self._load_cache(cache_key)
            if cached_data is not None and age_hours < max_age_hours:
                print(f"[CACHE HIT] {call_name} ({age_hours:.1f}h old)")
                return cached_data

        # Caché no disponible o demasiado viejo → fetch
        print(f"[CACHE MISS] {call_name} → MCP call")
        data = mcp_call_fn(call_name, params)

        # Guardar en caché
        self._save_cache(cache_key, call_name, params, data)

        return data

    def _load_cache(self, cache_key: str) -> tuple[Optional[Any], float]:
        """
        Carga datos de caché.

        Returns:
            (data, age_hours) o (None, 999) si no existe
        """
        data, age_hours, _versions = self._load_cache_full(cache_key)
        return data, age_hours

    def _load_cache_full(self, cache_key: str) -> tuple[Optional[Any], float, Optional[Dict[str, str]]]:
        """
        Carga datos de caché incluyendo las versiones de ADR/KB con las que se guardó.

        Returns:
            (data, age_hours, versions) o (None, 999, None) si no existe
        """
        cache_file = self._cache_file(cache_key)
        if not cache_file.exists():
            return None, 999, None

        try:
            cached = json.loads(cache_file.read_text())
            timestamp = cached.get('timestamp', 0)
            age_seconds = time.time() - timestamp
            age_hours = age_seconds / 3600

            return cached.get('data'), age_hours, cached.get('versions')
        except:
            return None, 999, None

    def _save_cache(self, cache_key: str, call_name: str, params: Dict, data: Any) -> None:
        """Guarda datos en caché."""
        cache_file = self._cache_file(cache_key)

        cache_entry = {
            'call_name': call_name,
            'params': params,
            'timestamp': time.time(),
            'data': data
        }

        cache_file.write_text(json.dumps(cache_entry, indent=2, default=str))

        # Actualizar manifest
        self.manifest[cache_key] = {
            'call_name': call_name,
            'created': datetime.now().isoformat(),
            'file_size_kb': cache_file.stat().st_size / 1024
        }
        self._save_manifest()

    def invalidate_all(self) -> None:
        """Invalida TODO el caché del proyecto."""
        import shutil
        if self.cache_dir.exists():
            shutil.rmtree(self.cache_dir)
            self.cache_dir.mkdir(parents=True, exist_ok=True)
            self.manifest = {}
            print(f"[CACHE CLEARED] Proyecto {self.project_name}")

    def invalidate_pattern(self, pattern: str) -> None:
        """Invalida caché que matchee un patrón (ej. 'get_trades')."""
        for key in list(self.manifest.keys()):
            if pattern.lower() in key.lower():
                cache_file = self._cache_file(key)
                if cache_file.exists():
                    cache_file.unlink()
                del self.manifest[key]
        self._save_manifest()
        print(f"[CACHE CLEARED] Pattern '{pattern}' (affected keys: {sum(1 for k in self.manifest if pattern.lower() in k.lower())})")

    def is_cache_valid_by_version(self, cached_versions: Dict[str, str] = None) -> bool:
        """
        Verificar si caché sigue siendo válida comparando versiones.

        (Fase 4: Versionado automático)

        Args:
            cached_versions: Versiones cuando se guardó el caché.
                            Si None, usa las del manifest actual.

        Returns:
            True si versiones coinciden, False si algo cambió
        """
        if cached_versions is None:
            # Leer versiones del manifest
            cached_versions = self.manifest.get('versions', {})
            if not cached_versions:
                # Si no hay versiones en manifest, asumir inválido (forzar refresh)
                return False

        # Comparar con versiones actuales
        valid = self.version_tracker.is_cache_valid(cached_versions, self.current_versions)

        if not valid:
            # Explicar qué cambió
            diff = self.version_tracker.get_version_diff(cached_versions, self.current_versions)
            print(f"[CACHE INVALIDATED] Cambios de versión detectados:")
            print(self.version_tracker.explain_invalidation(diff))

        return valid

    def get_or_fetch_with_version_check(
        self,
        call_name: str,
        params: Dict,
        mcp_call_fn,
        max_age_hours: int = 24,
        force_refresh: bool = False
    ) -> Any:
        """
        Get or fetch con validación automática de versiones.

        (Fase 4: Versionado automático)

        Similar a get_or_fetch(), pero además:
          1. Verifica que las versiones de ADRs/KB no cambiaron
          2. Invalida automáticamente si algo cambió
          3. Refetcha datos sin esperar TTL
        """
        cache_key = self._cache_key(call_name, params)

        # Cargar el archivo de caché (si existe) junto con las versiones con las que se guardó
        if not force_refresh:
            cached_result, age_hours, cached_versions = self._load_cache_full(cache_key)
            if cached_result is not None and age_hours < max_age_hours:
                # Verificar versiones ANTES de servir el hit: un TTL "permanente"
                # no debe blindar datos calculados con conocimiento ya desactualizado.
                if cached_versions and not self.is_cache_valid_by_version(cached_versions):
                    force_refresh = True
                else:
                    print(f"[CACHE HIT] {call_name} ({age_hours:.1f}h old, versions OK)")
                    return cached_result

        # Caché no disponible/viejo/inválido → fetch y guardar con versiones
        print(f"[CACHE MISS] {call_name} - MCP call")
        data = mcp_call_fn(call_name, params)
        self._save_cache_with_versions(cache_key, call_name, params, data)

        return data

    def _save_cache_with_versions(self, cache_key: str, call_name: str, params: Dict, data: Any) -> None:
        """Guardar caché incluyendo versiones de ADRs/KB."""
        cache_file = self._cache_file(cache_key)

        cache_entry = {
            'call_name': call_name,
            'params': params,
            'timestamp': time.time(),
            'data': data,
            'versions': self.current_versions  # Guardar versiones actuales
        }

        cache_file.write_text(json.dumps(cache_entry, indent=2, default=str))

        # Actualizar manifest
        self.manifest[cache_key] = {
            'call_name': call_name,
            'created': datetime.now().isoformat(),
            'file_size_kb': cache_file.stat().st_size / 1024,
            'versions': self.current_versions
        }
        self._save_manifest()

    def stats(self) -> Dict:
        """Estadísticas del caché."""
        total_size = sum(
            (self.cache_dir / f"{k}.json").stat().st_size / 1024
            for k in self.manifest if (self.cache_dir / f"{k}.json").exists()
        )

        return {
            'project': self.project_name,
            'total_entries': len(self.manifest),
            'total_size_kb': round(total_size, 2),
            'entries': self.manifest,
            'current_versions': self.current_versions
        }


# Ejemplo de uso
if __name__ == '__main__':
    def mock_mcp_call(call_name: str, params: Dict):
        """Mock de una llamada MCP (para testing)."""
        return {
            'call': call_name,
            'params': params,
            'data': [1, 2, 3],
            'timestamp': time.time()
        }

    # Test
    mgr = CacheManager('test_project')

    # Primera llamada → MCP
    result1 = mgr.get_or_fetch('get_balance', {'asset': 'BTC'}, mock_mcp_call)
    print(f"Resultado 1: {result1}")

    # Segunda llamada → CACHE
    result2 = mgr.get_or_fetch('get_balance', {'asset': 'BTC'}, mock_mcp_call)
    print(f"Resultado 2 (desde caché): {result2}")

    # Estadísticas
    print(f"\nEstadísticas: {mgr.stats()}")

    # Limpiar
    mgr.invalidate_all()
