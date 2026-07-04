# Casos Reales Auditados (Nivel C1)

**Ubicación:** `knowledge/cases/`

**Característica:** Documentación de **casos reales auditados** en proyectos del usuario.

**Autoridad:** `verified` — casos reales = máxima confianza, no es teoría

---

## Estado Actual (Fase 3)

Los 20 casos han sido **convertidos a archivos `.md` individuales** con metadatos YAML.

**Archivo original (legacy):** [`cointracking_casos_v2.yaml`](cointracking_casos_v2.yaml) — mantenido como backup

---

## Casos Documentados (20 archivos .md)

**Cada caso tiene su propio archivo:**

1. [CT-001-transferencia-entre-exchanges-importada-solo-en-origen.md](CT-001-transferencia-entre-exchanges-importada-solo-en-origen.md) — Transferencias huérfanas
2. [CT-002-venta-sin-historial-de-compra-previo-missing-purchase-history.md](CT-002-venta-sin-historial-de-compra-previo-missing-purchase-history.md) — Missing Purchase History
3. [CT-003-api-y-csv-importados-simultaneamente-duplicado-por-doble-fuente.md](CT-003-api-y-csv-importados-simultaneamente-duplicado-por-doble-fuente.md) — API + CSV solapados
4. [CT-004-balance-negativo-por-orden-cronologico-incorrecto-zona-horaria.md](CT-004-balance-negativo-por-orden-cronologico-incorrecto-zona-horaria.md) — Zona horaria
5. [CT-005-recompensas-de-staking-clasificadas-como-deposito-generico.md](CT-005-recompensas-de-staking-clasificadas-como-deposito-generico.md) — Staking
6. [CT-006-binance-convert-importado-como-venta-y-compra-independientes.md](CT-006-binance-convert-importado-como-venta-y-compra-independientes.md) — Binance Convert
7. [CT-007-transferencia-interna-confundida-con-venta.md](CT-007-transferencia-interna-confundida-con-venta.md) — Transferencias internas
8. [CT-008-duplicados-aparentes-por-ejecucion-parcial-de-una-orden.md](CT-008-duplicados-aparentes-por-ejecucion-parcial-de-una-orden.md) — Órdenes parciales
9. [CT-009-comision-fee-omitida-en-la-importacion.md](CT-009-comision-fee-omitida-en-la-importacion.md) — Comisiones
10. [CT-010-airdrop-registrado-como-compra-con-coste-artificial.md](CT-010-airdrop-registrado-como-compra-con-coste-artificial.md) — Airdrops
11. [CT-011-lending-tratado-como-transferencia-generica.md](CT-011-lending-tratado-como-transferencia-generica.md) — Lending
12. [CT-012-balance-negativo-por-importacion-parcial-via-api.md](CT-012-balance-negativo-por-importacion-parcial-via-api.md) — Importación parcial
13. [CT-013-wallet-externa-no-importada-fondos-desaparecidos.md](CT-013-wallet-externa-no-importada-fondos-desaparecidos.md) — Wallets externas
14. [CT-014-recompensas-de-mineria-mining-registradas-como-deposito.md](CT-014-recompensas-de-mineria-mining-registradas-como-deposito.md) — Mining
15. [CT-015-swap-defi-fragmentado-en-varias-operaciones-on-chain.md](CT-015-swap-defi-fragmentado-en-varias-operaciones-on-chain.md) — Swaps DeFi
16. [CT-016-duplicados-por-reimportacion-completa-del-mismo-periodo.md](CT-016-duplicados-por-reimportacion-completa-del-mismo-periodo.md) — Reimportación
17. [CT-017-coste-cero-por-compra-omitida-de-ejercicios-anteriores.md](CT-017-coste-cero-por-compra-omitida-de-ejercicios-anteriores.md) — Años anteriores
18. [CT-018-token-renombrado-interpretado-como-un-activo-distinto.md](CT-018-token-renombrado-interpretado-como-un-activo-distinto.md) — Tokens renombrados
19. [CT-019-balance-negativo-tras-eliminar-una-compra-confundida-con-duplicado.md](CT-019-balance-negativo-tras-eliminar-una-compra-confundida-con-duplicado.md) — Falso positivo
20. [CT-020-advertencia-tecnica-interpretada-como-error-fiscal-definitivo.md](CT-020-advertencia-tecnica-interpretada-como-error-fiscal-definitivo.md) — Advertencias

**Resumen:** 20 casos, cada uno con metadatos YAML (KB-C1-001 a KB-C1-020)

---

## Próximas Sesiones

**Fase 3:**
1. Convertir YAML → archivos `.md` individuales (`CT-001-*.md`, `CT-002-*.md`, etc.)
2. Insertar metadatos YAML en cada archivo (`id: KB-C1-001`, etc.)
3. Actualizar INDEX.md con referencias a archivos individuales

---

## Relación con otros Niveles

- **Nivel B1-B3:** Estos casos **fundamentan** el conocimiento operativo
- **Nivel C2:** Patrones se **derivan** de estos casos
- **Nivel C3:** Procedimientos se **basan en** estos casos
- **Nivel D1-D2:** Checklists y árboles de decisión se **construyen a partir de** estos casos

---

## Esquema Actual (YAML)

Cada caso tiene 16 campos obligatorios:

```yaml
- id: CT-002
  titulo: "FLOKI: 29 transacciones idénticas no son duplicadas"
  categoria: duplicados
  sintomas: [síntomas observados]
  causa_probable:
    hecho: [la causa raíz verificada]
  evidencia_minima: [qué se necesita para confirmar]
  pasos_diagnostico: [cómo diagnosticar]
  solucion_recomendada: [qué hacer]
  anti_patron: [qué NO hacer]
  por_que_falso_positivo: [por qué se confunde]
  nivel_confianza: verificado
  nivel_riesgo: alto
  impacto_fiscal_potencial: [pérdida de dinero / multa]
  senales_tempranas: [indicios iniciales]
  validacion_antes_despues: [cómo validar]
  vigencia:
    fecha_revision: 2026-07-03
    motivo_caducidad_potencial: [cuándo revalidar]
    fuente_recomendada_para_revalidar: [dónde verificar]
```

---

## Política de Vigencia

- Cada caso declara `vigencia` (fecha_revision, motivo_caducidad, fuente para revalidar)
- Antes de citar un caso, verificar que no esté desfasado
- Casos con `nivel_confianza: hipotesis` o `pendiente_verificar` requieren reverificación antes de usar

---

## Acceso Rápido

- **Por síntoma:** Ver `knowledge/cointracking/TROUBLESHOOTING.md` (índice por síntoma que remite a casos)
- **Por patrón:** Ver `knowledge/patterns/` (patrones generalizados desde casos)
- **Por categoría:** Ver tabla de arriba

---

## Próxima Migración

En Fase 3, cada caso tendrá su propio archivo:

```
knowledge/cases/
├── CT-001-duplicate-same-timestamp.md
├── CT-002-floki-batching.md
├── CT-003-missing-purchase-history.md
└── ... (CT-004 a CT-020)
```

Cada archivo tendrá:
- Metadatos YAML (`id: KB-C1-002`, etc.)
- Contenido del caso (antes en YAML, ahora markdown)
- Enlaces a ADRs relacionados
- Enlaces a patrones derivados
