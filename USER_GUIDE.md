# Guía de usuario: CoinTracking Expert

## 1. Introducción

### Qué es este agente

Un **asistente de IA** que audita tus datos de CoinTracking (la plataforma de seguimiento de criptomonedas) para encontrar problemas de reconciliación y ayudarte a preparar la declaración de la renta española (IRPF).

El agente vive en Claude Code y accede a tus datos de dos formas:
- **MCP API** (datos en vivo, directo desde CoinTracking)
- **CSV export** (exportas el archivo tú; el agente lo analiza localmente)

### Qué hace

✅ **Auditar** tus movimientos cripto y encontrar:
- Transferencias sin origen ni destino (huérfanas)
- Ventas que no tienen base de coste
- Duplicados
- Saldos que no cuadran
- Incoherencias fiscales españolas

✅ **Explicar** cada hallazgo: qué significa, por qué le importa, qué hacer

✅ **Guiarte paso a paso** a corregir los datos en la web de CoinTracking

✅ **Preparar la declaración fiscal** del año: calcular ganancias patrimoniales, identificar información para el Modelo 721, etc.

### Qué NO es

❌ **No es asesoramiento fiscal.** El agente encuentra e interpreta; no produce cifras fiscales vinculantes. Lleva los resultados a tu asesor fiscal.

❌ **No es determinista.** Es un LLM: explica bien, pero puede haber matices. Siempre verifica los hallazgos.

❌ **No modifica nada.** Solo lee tus datos y analiza. Tú decides qué cambiar en CoinTracking.

### En qué puede ayudarte

**Si tienes un CSV o credenciales MCP:** El agente audita en minutos y explica qué está mal.

**Si quieres hacer la declaración de la renta:** El agente reconcilia primero, luego prepara lo que necesitas (ganancias, modelo 721, etc.) para llevar al asesor.

**Si no sabes si hay problemas:** El agente lo revela. Muchos usuarios descubren transferencias huérfanas o ventas mal registradas solo corriendo la auditoría.

---

## 2. ¿Cuál es tu perfil?

Elige uno. Te ahorrará tiempo.

### Perfil A: Ya conozco CoinTracking

✅ Has usado la plataforma al menos una vez  
✅ Sabes qué es un CSV export  
✅ Entiendes "base de coste", "FIFO" o "transferencia interna"  

**Salta directamente a [Tema 1](#tema-1-preparar-y-cargar-tus-datos).** Todas las secciones asumen que ya navegas CoinTracking.

### Perfil B: Primera vez con CoinTracking

❌ No sé qué es CoinTracking  
❌ No he usado la plataforma nunca  
❌ No conozco términos como "base de coste" o "FIFO"  

**Lee la sección [Guía para novatos](#guía-para-novatos-antes-de-empezar) primero.** Luego sigue el mismo flujo que el Perfil A.

---

## Guía para novatos (antes de empezar)

### ¿Qué es CoinTracking?

Una plataforma web (cointracking.info) donde registras todos tus movimientos cripto: compras, ventas, intercambios, transferencias. Sirve para:
1. **Saber cuánto tienes** (saldo)
2. **Calcular ganancias/pérdidas** (cuánto ganaste/perdiste en cada operación)
3. **Cumplir con Hacienda** (tener los datos listos para la declaración de la renta)

### Términos que verás aquí

| Término | Significa |
|---------|-----------|
| **CSV** | Un archivo de texto con tus datos en formato tabla. Puedes exportarlo desde CoinTracking. |
| **Base de coste** | Lo que pagaste por esa cripto. Si compraste 1 BTC a 20.000€, la base es 20.000€. |
| **FIFO** | Forma de calcular ganancias: "primera entrada, primera salida". Si compras 1 BTC a 20K y otro a 25K, y vendes 1, supones que vendiste el primero. |
| **Transferencia interna** | Mover cripto entre tus propias billeteras. No es venta (no genera ganancia/pérdida). |
| **Saldo** | Cuánta cripto tienes ahora. |
| **Modelo 721** | Formulario que Hacienda usa para registrar activos financieros en el extranjero (si tienes cripto, entra aquí). |

Ver [Glosario completo](#8-glosario-breve).

### Cómo obtener tus datos de CoinTracking

**CSV export (recomendado para auditoría):**
1. Entra en cointracking.info
2. Ve a **Reports** → **Trade Table**
3. Haz clic en **Export** (arriba a la derecha)
4. Elige **CSV** y descarga el archivo
5. Guárdalo en una carpeta (p. ej., `H:\cripto-datos\trades.csv`)

**MCP API (en vivo, si tienes credenciales de CoinTracking):**
- Requiere crear una "API Key" en CoinTracking (en Settings → API)
- El agente accede en directo sin necesidad de CSV
- Más rápido, pero necesita permisos

Usa **CSV si tienes dudas**; es más sencillo.

---

## Tema 1: Preparar y cargar tus datos

### Opción 1: Usar CSV export

**Paso 0:** Al pedirle al agente que audite o declare, lo primero que hará es preguntarte **con qué proyecto quieres trabajar** (o si quieres crear uno nuevo). Cada proyecto aísla tus datos de otros casos que pudieras tener.

**Paso 1:** Exporta el CSV desde CoinTracking (ve a [Guía para novatos](#cómo-obtener-tus-datos-de-cointracking) si no sabes cómo).

**Paso 2:** Mueve el archivo a la carpeta de **tu proyecto**:
```
H:\cointracking-expert\USER_INPUT\<nombre_de_tu_proyecto>\trades.csv
```

**Paso 3:** Cuando corras el agente, le dices: "Tengo un CSV export aquí" (o "audita mi CSV"). El agente lo lee.

**Ventajas:**
- Privado (el archivo está en tu PC)
- Funciona sin conexión
- Fácil de repetir (guardas el mismo archivo)

**Desventajas:**
- Manual: tienes que actualizar el CSV cada vez

---

### Opción 2: Usar MCP API (datos en vivo)

**Paso 1:** Entra en cointracking.info → **Settings** → **API**.

**Paso 2:** Crea una **API Key** (copia la clave).

**Paso 3:** Guarda la clave en un lugar seguro (no la compartas, no la commits).

**Paso 4:** Cuando corras el agente, le dices: "Usa la API" (o "audita con MCP"). El agente accede en directo.

**Ventajas:**
- Automático: siempre tienes los datos actuales
- Sin sincronización manual

**Desventajas:**
- Requiere credenciales (un poco más de setup)
- Requiere conexión a internet

---

### Privacidad y seguridad

✅ **Con CSV:** El archivo está en tu PC, nunca sale (solo el agente lo lee localmente).

✅ **Con API:** El agente accede en solo lectura (no modifica nada). Las credenciales nunca se guardan en el repositorio (excluidas en `.gitignore`).

---

## Tema 2: Ejecutar la auditoría

El agente audita con el comando **`/audit-cointracking`**.

### Para el Perfil A (conoces CoinTracking)

**Qué esperar en 5 minutos:**

1. Dices: `/audit-cointracking`
2. El agente carga tu CSV o conecta con la API
3. Analiza 5 tipos de problemas comunes:
   - Transferencias huérfanas (origen/destino faltante)
   - Ventas sin base de coste
   - Saldos negativos imposibles
   - Duplicados
   - Incoherencias fiscales españolas

4. **Resultado:** Un informe con cada hallazgo en formato:
   - **Qué es:** descripción
   - **Por qué le importa:** impacto (fiscal, contable)
   - **Qué hacer:** pasos concretos

---

### Para el Perfil B (novato)

**Paso a paso:**

1. **Prepara los datos** (sigue [Tema 1](#tema-1-preparar-y-cargar-tus-datos)).

2. **Corre la auditoría:**
   - Abre Claude Code
   - Escribe: `/audit-cointracking`
   - Presiona Enter

3. **El agente te pide confirmar:**
   - Qué datos usar (CSV o API)
   - Rango de fechas (opcional: "solo 2024", o "todo")

4. **Espera unos segundos.** El agente analiza.

5. **Lee el informe que sale:** verás una lista de hallazgos (si los hay) con explicaciones claras.

---

### Qué significan los hallazgos más comunes

#### Transferencia huérfana
**Qué es:** Un movimiento cripto sin origen claro (entra dinero, pero no sabes de dónde) o sin destino (sale, pero no sabes a dónde).

**Por qué importa:** Hacienda quiere saber de dónde vino. Si no aparece, parece renta no declarada.

**Qué hacer:** Busca en tus registros, en los exchanges que usabas, en tus carteras. Actualiza el registro en CoinTracking (Editar la operación, añadir descripción/referencia).

---

#### Venta sin base de coste
**Qué es:** Vendiste cripto, pero el sistema no encuentra la compra original (o cuánto pagaste).

**Por qué importa:** Para calcular la ganancia, necesitas: `ganancia = precio_venta - base_coste`. Sin base, no hay ganancia calculada; fiscal asume ganancia = 100% (lo peor).

**Qué hacer:** Busca la compra original en tu historial. Si no está, añádela manualmente en CoinTracking (Transactions → Add / Edit).

---

#### Saldo negativo
**Qué es:** El sistema dice que tienes -0.5 BTC (imposible; no puedes tener negativo).

**Por qué importa:** Error lógico. Alguien vendió sin tener. Puede ser un duplicado o un registro mal ordenado.

**Qué hacer:** Busca ventas sin compra correspondiente. Edita o borra el duplicado en CoinTracking.

---

#### Duplicado
**Qué es:** La misma operación aparece dos veces en el historial.

**Por qué importa:** Hacienda verá el doble de ganancia/pérdida. Impuestos incorrectos.

**Qué hacer:** Borra uno de los dos en CoinTracking (confirma antes que sean realmente iguales).

---

#### Incoherencia fiscal española
**Qué es:** Algo que según la ley española no cuadra (p. ej., un intercambio cripto-cripto registrado como venta; una compra en EUR pero registrada en otra divisa sin conversión correcta).

**Por qué importa:** DGT (Hacienda) podría rechazarlo o pedir aclaraciones.

**Qué hacer:** El agente sugiere la corrección. Edita el registro en CoinTracking.

---

## Tema 3: Actuar sobre los resultados

### Paso 1: Lee el informe

El agente genera un **informe** con formato:
```
Hallazgo #1: [descripción corta]
Qué es: [explicación]
Por qué importa: [impacto]
Recomendación: [pasos concretos]
```

Anota los **números de hallazgo** que necesites corregir.

---

### Paso 2: Corrige en CoinTracking

**Para editar una operación:**
1. Entra en cointracking.info
2. Ve a **Transactions** (o **Trade Table**)
3. Busca la operación por fecha/cantidad
4. Haz clic en el icono de **lápiz** (Edit)
5. Modifica lo necesario (fecha, cantidad, precio, descrición, etc.)
6. Guarda (botón **Save** o **Update**)

**Para borrar:**
- Haz clic en el icono de **basura** (Delete)
- CoinTracking pedirá confirmación

**Para añadir:**
- Botón **Add Transaction** o **Add Trade**
- Rellena: fecha, tipo (Buy/Sell/Transfer), cantidad, moneda, precio, descripción

---

### Paso 3: Verifica que cambió

Una vez hayas corregido en CoinTracking web:

1. **Espera 1-2 minutos** (CoinTracking actualiza lentamente)
2. **Re-audita:** vuelve a correr `/audit-cointracking`
3. **Comprueba:** el hallazgo desapareció o se resolvió

Si sigue ahí, pide ayuda en el chat (describe qué cambiaste y por qué sigue fallando).

---

### Cuándo re-auditar

- **Después de cada corrección importante:** verifica que funcionó
- **Una sola vez si hay muchos hallazgos:** audita, lista todos, corrígelos todos, luego re-audita
- **Ante cambios grandes:** si importaste más operaciones, re-audita

---

## Tema 4: Preparar la declaración fiscal (IRPF)

> ⚠️ **Esto no es asesoramiento fiscal.** Lleva estos números a tu asesor fiscal. Él decide qué declaras.

### Requisitos previos

✅ **Tus datos deben estar limpios** (audita primero, corrige los hallazgos)

✅ **Debes saber el año fiscal** (p. ej., 2024 para el IRPF que presentarás en 2025)

---

### Cómo funciona

**Comando:** `/spanish-tax-return`

El agente te preguntará:
1. **¿Qué año fiscal?** (2024, 2023, etc.)
2. **CSV o API?** (mismo que para la auditoría)

Luego:
1. **Audita automáticamente** tus datos (si hay problemas, los marca)
2. **Calcula ganancias patrimoniales** usando el método FIFO (primero que compres, primero que vendas)
3. **Identifica lo que va en el Modelo 721** (si tienes más de 50.000€ en cripto al 31 de diciembre)
4. **Genera un informe** con resumen fiscal: ganancias, pérdidas, saldos finales, etc.

**Resultado:** Un documento que llevas a tu asesor fiscal (él traduce esto al IRPF oficial).

---

### Qué sale del agente

- **Ganancias patrimoniales:** suma de ganancias en cada venta (aplicando FIFO)
- **Pérdidas patrimoniales:** suma de pérdidas (para compensar ganancias)
- **Saldo final de criptos:** cuánto tienes el 31 de diciembre
- **Criptos para Modelo 721:** si superas 50.000€, cuáles son
- **Notas y dudas:** cosas que tu asesor debe revisar manualmente

---

### Qué llevas al asesor fiscal

Descarga el **informe PDF/Excel** que genera el agente (o cópialo a mano). Dale al asesor:

1. **El resumen de ganancias/pérdidas**
2. **El saldo final** (qué tienes el 31 de diciembre)
3. **Lista de criptos para Modelo 721** (si aplica)
4. **Todas las dudas o correcciones** que el agente marcó

El asesor **rellena el IRPF oficial** con estos datos.

---

## 7. Preguntas frecuentes + Troubleshooting

### ¿Por qué el saldo no cuadra?

**Causa común:** Transferencias internas no marcadas correctamente (cuando mueves cripto entre tus propias billeteras, CoinTracking puede malinterpretarla).

**Qué hacer:**
1. Revisa las transferencias internas en CoinTracking
2. Asegúrate que estén marcadas como "Transfer" (no "Sell/Buy")
3. Edita si es necesario
4. Re-audita

---

### ¿Qué pasa si el CSV es muy grande?

Si tienes miles de operaciones, el CSV puede ser lento.

**Qué hacer:**
- Usa **rango de fechas** en la auditoría (p. ej., "solo 2024")
- O usa **MCP API** (es más rápido para volúmenes grandes)

---

### ¿Puedo usar solo la API, sin CSV?

Sí. Si tienes credenciales MCP de CoinTracking, el agente accede directamente. No necesitas CSV.

---

### El agente dice "transferencia huérfana" pero sé que es correcta

Puede pasar. El agente es automático; a veces pierde contexto.

**Qué hacer:**
1. Verifica manualmente: ¿de dónde vino esa cripto realmente?
2. Edita el registro en CoinTracking (añade descripción, referencia, o actualiza la billetera origen/destino)
3. Re-audita
4. Si sigue marcándolo, avísale al agente en el chat: "Esta transferencia es de Exchange X a mi billetera Y, no es huérfana"

---

### ¿Cuánto tarda una auditoría?

- **CSV pequeño (<500 operaciones):** 1-2 minutos
- **CSV grande (>5000 operaciones):** 5-10 minutos
- **API:** suele ser más rápido que CSV

---

### ¿Puedo auditar solo un año?

Sí. Dile al agente: "Audita solo 2024" (o el año que quieras). Se filtra automáticamente.

---

### ¿Qué pasa con las comisiones/fees?

Las comisiones de transacción **restan de la ganancia** (o suman a la base de coste). CoinTracking las suele registrar automáticamente. Si no aparecen, añádelas manualmente.

---

### El agente me dice que hay un duplicado pero no lo veo

Los duplicados pueden ser sutiles: misma operación, fechas muy cercanas, pero cantidades ligeramente distintas (p. ej., 1.0 BTC vs 1.00001 BTC).

**Qué hacer:**
1. Usa la herramienta "Find" (Ctrl+F) en CoinTracking
2. Busca por cantidad exacta
3. Si encuentras dos casi iguales, una es probable duplicado
4. Bórrala (confirma antes con tu exchange o billetera)

---

### ¿Dónde están mis informes?

Los informes se guardan en la carpeta de tu proyecto:
```
H:\cointracking-expert\reports\output\<nombre_de_tu_proyecto>\
```

Llevan fecha: `AUDIT_2026-07-03.md`, `TAX_RETURN_2024_2026-07-03.md`, etc.

Búscalos ahí para descargarlos o compartirlos con tu asesor.

---

## 8. Glosario breve

| Término | Definición |
|---------|-----------|
| **API** | Interfaz que permite que programas hablen con CoinTracking automáticamente (sin meterte en la web a mano). |
| **Base de coste** | Lo que pagaste originalmente por esa cripto. Necesario para calcular ganancias. |
| **CSV** | Archivo de texto con datos en formato tabla (filas y columnas). CoinTracking te lo deja exportar. |
| **DGT** | Dirección General de Tributos (la oficina de Hacienda española que gestiona impuestos). |
| **Duplicado** | Misma operación registrada dos veces por error. Cuesta impuestos extras si no se borra. |
| **FIFO** | Método de cálculo: "Primera Entrada, Primera Salida". Si compras dos lotes, vendes el más antiguo primero. |
| **Ganancia patrimonial** | Dinero que ganas al vender cripto a precio más alto que lo que pagaste. Se declara en el IRPF. |
| **IRPF** | Impuesto sobre la Renta de las Personas Físicas. La declaración que presentas cada año a Hacienda. |
| **Modelo 721** | Formulario para declarar activos en el extranjero (si tienes cripto, va aquí). Obligatorio si superas 50.000€. |
| **MCP** | En este proyecto, es la API de CoinTracking integrada (los datos en vivo desde la web). |
| **Transferencia interna** | Mover cripto entre tus propias billeteras. No es venta (no genera impuestos directos). |
| **Transferencia huérfana** | Operación sin origen o destino claro. El agente las marca como problema. |

**¿Necesitas más detalle?** Ve a [docs/GLOSSARY.md](docs/GLOSSARY.md) (si existe en tu proyecto).

---

## Siguientes pasos

✅ **Tienes los datos listos?** → Ve a [Tema 1](#tema-1-preparar-y-cargar-tus-datos)

✅ **¿Quieres auditar?** → Escribe `/audit-cointracking`

✅ **¿Quieres preparar la declaración?** → Escribe `/spanish-tax-return`

✅ **¿Tienes dudas?** → Lee las [Preguntas frecuentes](#7-preguntas-frecuentes--troubleshooting)

---

**Última actualización:** 3 de julio de 2026  
**Versión:** 1.0
