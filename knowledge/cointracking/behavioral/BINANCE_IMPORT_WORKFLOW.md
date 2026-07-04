---
id: KB-B2-008
title: "Workflow de Importación de Binance: API vs CSV Paso a Paso"
level: B
domain: cointracking
source: "Procedimiento estándar + análisis"
authority: verified
last_verified: 2026-07-05
valid_from: 2024-01-01
valid_until: 2027-12-31
confidence: high
version: 1.0

related_adr:
  - ADR-003
  - ADR-010

related_docs:
  - knowledge/cointracking/behavioral/BINANCE_SPOT_MECHANICS.md
  - knowledge/cointracking/behavioral/API_VS_CSV_OVERLAP.md

tags:
  - binance
  - import
  - workflow
  - step-by-step
  - behavioral

notes: "Guía paso a paso: cómo importar Binance en CoinTracking (API recomendado)."
---

# Workflow de Importación de Binance: API vs CSV

## Decisión: ¿API o CSV?

| Factor | API | CSV |
|--------|-----|-----|
| **Velocidad** | Automática, continua | Manual, una sola vez |
| **Precisión** | Máxima (datos en vivo) | Depende del export |
| **Completitud** | Toda la historia | Solo lo que exportes |
| **Seguridad** | Lee-solo (sin dinero) | Archivo local |
| **Complejidad** | Requiere API keys | Solo archivo |

**Recomendación:** **API (primario) + CSV (backup)**

---

## Opción 1: Importación por API (RECOMENDADO)

### Paso 1: Crear API Keys en Binance

**Ubicación:** `https://www.binance.com/en/account/api-management`

**Pasos:**
1. Haz login en Binance
2. Perfil → API Management (o Seguridad → Gestión de API)
3. Crea una nueva clave ("Create" o "Crear")

### Paso 2: Configurar Permisos

**Ojo:** Solo lectura de historial, SIN permisos de trading/retirada.

```
Permisos necesarios:
✓ Ver datos de transacciones (Lectura)
✓ Acceso a historial

Permisos PROHIBIDOS:
✗ Trading (comprar/vender)
✗ Transferencias (retirar)
✗ Cambiar configuración
```

**Por qué:** Si alguien roba la API key, no puede gastar tu dinero.

### Paso 3: Whitelist IP (Opcional pero Recomendado)

```
Agregar tu IP actual (o rango):
- Ve a "Restrict access to trusted IPs only"
- Ingresa tu IP (ej. 203.0.113.45)
- Guarda

Resultado: API key solo funciona desde tu red
```

### Paso 4: Configurar en CoinTracking

**Ubicación:** `CoinTracking → Settings → Exchanges`

**Pasos:**
1. Clic en "Add Exchange"
2. Selecciona "Binance"
3. Pega:
   - API Key (de Binance)
   - Secret Key (de Binance)
4. Clic en "Connect"
5. CoinTracking verifica conexión

**Resultado:** ✓ Binance conectado. CoinTracking importará automáticamente.

### Paso 5: Verificar Importación

```
CoinTracking → Home (Dashboard):
  ¿Aparecen operaciones de Binance?
  
  SÍ → ✓ OK, espera 24h para sincronización completa
  NO → ✗ Verificar API keys / permisos
```

**Tiempo:** 5-10 minutos. Datos completos en 24h.

---

## Opción 2: Importación por CSV (Si API No Funciona)

### Paso 1: Exportar CSV desde Binance

**Ubicación:** `Binance → Historial → (Spot, Margin, etc)`

**Pasos:**
1. Login en Binance
2. Menú: "Portfolio" → "Spot" (o Margin/Futures según lo que quieras)
3. Haz scroll → botón "Download" o "Exportar"
4. Selecciona rango de fechas:
   - **Importante:** Desde la fecha de tu primer depósito
   - No es necesario ser exacto (CoinTracking deduplica)
5. Selecciona formato: **CSV**
6. Descarga el archivo

**Resultado:** `trades_YYYY-MM-DD.csv` (u otro nombre)

### Paso 2: Preparar el CSV

**Verificar contenido:** Abre el CSV en Excel/LibreOffice

```
Debe tener columnas como:
- Date
- Market
- Price
- Amount
- Fee

Si falta algo importante → usar API en su lugar
```

### Paso 3: Importar en CoinTracking

**Ubicación:** `CoinTracking → Settings → Import → From File (CSV)`

**Pasos:**
1. Clic en "Upload CSV"
2. Selecciona el archivo descargado
3. CoinTracking detecta formato automáticamente
4. Revisa preview:
   - ¿Columnas están correctas?
   - ¿Fechas están en formato correcto?
5. Clic en "Import"

**Tiempo:** 5 minutos (plus procesamiento, 1-2 minutos)

### Paso 4: Verificar Importación

```
CoinTracking → Home:
  ¿Operaciones de Binance?
  
  SÍ → ✓ OK, revisa que no haya duplicados
  NO → Verificar formato CSV / contactar soporte CT
```

---

## Diferencias en Resultados

### Con API

```
Operaciones mostradas:
- Completo: Desde apertura de cuenta
- Detalle: Incluye fees, comisiones, tickers exactos
- Actualización: Automática (nuevas operaciones cada 24h)
- Riesgo: Cero (lectura solamente)
```

### Con CSV

```
Operaciones mostradas:
- Parcial: Solo lo que exportaste
- Detalle: Depende del formato del CSV
- Actualización: Manual (cada vez que exportes)
- Riesgo: Archivo local (pérdida si borras)
```

---

## Workflow Recomendado (Híbrido)

```
Día 1:
  1. Crear API keys en Binance (permisos read-only)
  2. Conectar API en CoinTracking
  3. Esperar 24h para sincronización inicial

Día 2+:
  1. CoinTracking trae datos automáticamente
  2. Revisar cada semana que no haya problemas
  3. (Opcional) Descargar CSV como backup anual

En caso de problema con API:
  1. Descargar CSV manual
  2. Importar en CoinTracking
  3. Investigar por qué falló API
```

---

## Troubleshooting

### Problema: "API Key invalid"

**Soluciones:**
1. ¿Copiaste la API Key completa (sin espacios)?
2. ¿Copiaste la Secret Key (no la API key dos veces)?
3. ¿Están habilitados los permisos de lectura en Binance?
4. ¿Pasó menos de 1 minuto desde que creaste la key? (espera 5 min)

### Problema: "No hay operaciones importadas"

**Soluciones:**
1. ¿Tienes operaciones en Binance (no solo depósito)?
2. ¿La API key tiene permisos de "Ver datos de transacciones"?
3. ¿Verificaste la IP whitelist (si está activada)?

### Problema: "Duplicados después de importar CSV"

**Ver:** `API_VS_CSV_OVERLAP.md`
- **Solución:** Eliminar CSV y mantener solo API (recomendado)

---

## Checklist Final

```
[ ] API keys creados (solo lectura)
[ ] Permisos verificados (sin trading/retirada)
[ ] IP whitelisted (opcional, recomendado)
[ ] API conectado en CoinTracking
[ ] Datos importados (verificar 24h después)
[ ] Balance coincide con Binance (comparar)
[ ] Ningún duplicado API+CSV (si usaste ambos)
```

---

## Integración

- **BINANCE_SPOT_MECHANICS.md:** Qué tipos de operaciones importar
- **API_VS_CSV_OVERLAP.md:** Cómo evitar duplicados
