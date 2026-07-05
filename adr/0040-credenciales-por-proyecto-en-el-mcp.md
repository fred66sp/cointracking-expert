---
# Versionado para invalidación de caché (VersionTracker / ADR-039).
# Incrementar `version` al cambiar la decisión de forma material;
# los cachés calculados con la versión anterior se invalidan solos.
version: 1.0
---

# ADR-040: Credenciales por proyecto en el MCP (multi-cuenta opcional)

**Status:** Accepted

**Date:** 2026-07-05

## Context

ADR-016 decidió que credenciales, tier y limitador de tasa son **del proceso**, no del proyecto: `cointracking_switch_project` cambia la carpeta de caché, pero todas las llamadas frescas siguen consultando la misma cuenta de CoinTracking. Para el caso de uso original (una persona, una cuenta, varios ejercicios fiscales como proyectos) es correcto.

La revisión de robustez de 2026-07-05 lo señaló como limitación con un riesgo real asociado: si algún día dos proyectos correspondieran a **cuentas distintas** (p. ej. auditar la cuenta de un familiar o un cliente — la proyección de uso del propio proyecto contempla decenas de proyectos/año), la caché del proyecto B se llenaría silenciosamente con datos de la cuenta A. Se mitigó primero alineando la documentación (aviso en `app.go`, en la descripción del tool y en `CLAUDE.md`), pero la mitigación era solo un cartel de "no pases por aquí". El usuario pidió cerrar el MCP "del todo, haciendo las cosas bien".

Restricción de seguridad no negociable: las credenciales **no pueden viajar por la conversación** (ni como parámetro de `switch_project` ni en ninguna respuesta) — quedarían registradas en el contexto del LLM y en los logs de la sesión.

## Decision

Soporte **opcional** de credenciales por proyecto vía ficheros `.env` locales, sin tocar la conversación:

1. **Nuevo flag `--project-env-dir`** (por defecto vacío = desactivado = comportamiento actual exacto). Apunta a un directorio local, fuera del repo o ignorado por git, que **nunca** se versiona.
2. **Resolución de credenciales al fijar proyecto** (arranque con `--project` y cada `cointracking_switch_project`):
   - Si `--project-env-dir` está vacío → credenciales del proceso (como hasta ahora).
   - Si `{dir}/{project}.env` **no existe** → credenciales del proceso (fallback explícito, se registra en el log).
   - Si existe → debe contener `COINTRACKING_API_KEY` y `COINTRACKING_API_SECRET` (mismos nombres que las variables de entorno del proceso; `COINTRACKING_TIER` opcional). **Si está incompleto es un error que aborta el switch** — nunca un fallback silencioso a la cuenta del proceso, porque consultar la cuenta equivocada es exactamente el accidente que esta función previene (fail-closed, ADR-009).
3. **Cliente y limitador de tasa pasan a ser intercambiables bajo el mismo mutex** que ya protege cfg/cache/store (accesores `Client()`/`Rate()`); el límite horario de CoinTracking es por API key, así que una cuenta nueva estrena tracker propio. Si el proyecto destino resuelve a las **mismas** credenciales, se reutilizan cliente y tracker actuales (conserva la ventana horaria consumida).
4. **Orden de operaciones en el switch:** las credenciales del destino se resuelven y validan **antes** de cerrar nada (un `.env` malformado = error limpio sin tocar estado); el swap de cliente ocurre solo después de abrir con éxito la caché del destino (mismo patrón de rollback que ya existía).
5. **Trazabilidad sin secretos:** la respuesta de `switch_project` incluye `credentials_source` (`process` | `project-env`) y la API key **ofuscada** (`Obfuscate`, primeros 6 + últimos 3), para que el agente y el usuario siempre sepan qué cuenta está activa sin exponer nada.

## Consequences

- ✅ Dos proyectos pueden corresponder a dos cuentas de CoinTracking distintas sin reiniciar el servidor ni contaminar cachés.
- ✅ Compatibilidad total hacia atrás: sin `--project-env-dir`, el binario se comporta exactamente igual que antes.
- ✅ Los secretos nunca pasan por la conversación ni por el repo: viven en ficheros locales del usuario.
- ✅ Fail-closed: un fichero de credenciales incompleto aborta en vez de degradar a la cuenta equivocada.
- ⚠️ El tamaño de la caché L1 sigue siendo del proceso (se fija al arrancar por tier); un tier por proyecto solo ajusta el límite de tasa. Aceptable: el tamaño de L1 es una optimización, no un dato de corrección.
- ⚠️ El usuario es responsable de proteger el directorio de `.env` (permisos del sistema de archivos); el servidor solo lo lee.

## Notes

- Sustituye la mitigación documental de M6 (2026-07-05) por una solución real; los avisos de "misma cuenta" se actualizan para describir el comportamiento condicional.
- Relación: ADR-013 (multi-proyecto), ADR-016 (switch en caliente — su decisión "credenciales del proceso" pasa a ser el **caso por defecto**, no el único), ADR-009 (fail-closed, cero secretos en conversación).
- Verificación: tests de integración con servidor HTTP falso que asevera la cabecera `Key` por petición (cuenta A vs cuenta B), fallback sin fichero, error con fichero incompleto sin tocar estado, y conservación de la ventana de tasa al reutilizar credenciales.
