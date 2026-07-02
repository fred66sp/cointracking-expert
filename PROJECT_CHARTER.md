# Carta de Proyecto

**Proyecto:** Framework de CoinTracking Expert

**Estado:** Borrador

**Versión:** 0.1.0

**Propietario:** Proyecto CoinTracking Expert

**Última actualización:** 2026-07-02

---

# 1. Visión

Construir el framework de código abierto más confiable para auditar, reconciliar y validar bases de datos de CoinTracking, con soporte de primera clase para contabilidad de criptomonedas e informes fiscales españoles.

El proyecto tiene como objetivo convertirse en la implementación de referencia para reconciliación de criptomonedas al combinar algoritmos deterministas, conocimiento del dominio y diagnósticos asistidos por IA.

---

# 2. Misión

Permitir a los usuarios verificar que sus datos de CoinTracking sean completos, internamente consistentes y listos para producir informes fiscales confiables.

El framework debe explicar cada problema detectado, identificar su causa raíz y recomendar la acción correctiva mínima respaldada por evidencia.

---

# 3. Objetivos

El proyecto proporciona:

- Reconciliación de transacciones
- Reconstrucción de libro mayor
- Validación de tenencias
- Detección de duplicados
- Emparejamiento de transferencias
- Análisis de historial de compras faltante
- Validación de consistencia fiscal
- Informes de auditoría profesionales
- Diagnósticos asistidos por IA
- Base de conocimiento extensible

---

# 4. Alcance

El proyecto incluye:

## CoinTracking

- Importaciones CSV
- Importaciones API
- Transacciones manuales
- Informes
- Advertencias
- Tenencias
- Informes fiscales

## Exchanges

Inicialmente:

- Binance
- Coinbase
- Kraken
- Bybit
- OKX
- KuCoin
- BingX

Se pueden soportar exchanges adicionales más tarde.

## Billeteras

- Ledger Live
- MetaMask
- Trust Wallet
- Rabby
- Billeteras hardware

## Blockchains

Inicialmente:

- Bitcoin
- Ethereum
- BNB Chain
- Solana
- Polygon
- Arbitrum
- Base

---

# 5. Fuera de alcance

Las siguientes características están intencionalmente excluidas de las primeras versiones:

- Gestión de carteras
- Automatización de trading
- Ejecución en exchange
- Asesoramiento de inversión
- Presentación de impuestos
- Servicios de custodia

El framework valida información.

No ejecuta operaciones financieras.

---

# 6. Principios guía

## Basado en evidencia

Cada conclusión debe estar respaldada por datos observables.

Sin suposiciones.

---

## Reproducibilidad

El mismo dataset siempre debe producir el mismo resultado de auditoría.

---

## Explicabilidad

Cada problema detectado debe incluir:

- Causa
- Evidencia
- Impacto
- Acción recomendada

---

## Intervención mínima

Nunca recomendar eliminar o modificar transacciones sin suficiente evidencia.

---

## Desarrollo impulsado por documentación

Ninguna característica será implementada antes de que su especificación funcional haya sido aprobada.

---

## Arquitectura modular

Cada componente debe ser reemplazable sin afectar al resto del sistema.

---

# 7. Objetivos de calidad

El framework siempre debe apuntar a producir:

- Cero falsos positivos cuando sea razonablemente posible
- Resultados deterministas
- Trazabilidad completa
- Cálculos reproducibles
- Informes legibles para humanos
- Informes legibles para máquinas

---

# 8. Criterios de éxito

Una auditoría se considera completa solo cuando:

- No existen balances negativos inexplicados
- No permanece historial de compras faltante sin resolver
- Todas las transacciones duplicadas son identificadas
- Las transferencias están emparejadas o justificadas explícitamente
- Las tenencias reconstruidas desde el historial coinciden con los balances esperados
- Los informes generados son internamente consistentes

---

# 9. Principios de arquitectura

El framework está compuesto de motores independientes.

Ejemplo:

```
Capa de importación
        │
Capa de normalización
        │
Motor de auditoría
        │
├── Motor de duplicados
├── Motor de transferencias
├── Motor de libro mayor
├── Motor de tenencias
├── Motor FIFO
└── Motor de informes
```

Cada motor debe exponer interfaces bien definidas.

Las reglas de negocio deben permanecer independientes de los modelos de IA.

---

# 10. Estrategia de conocimiento

El conocimiento es un componente de primera clase del proyecto.

La documentación debe organizarse en:

- CoinTracking
- Exchanges
- Billeteras
- Blockchains
- Tributación
- Reconciliación
- Contabilidad
- Metodología de auditoría
- Casos del mundo real

El conocimiento debe estar versionado junto con el código fuente.

---

# 11. Estrategia de IA

La inteligencia artificial es una interfaz, no la fuente de verdad.

El framework debe separar:

- Reglas de negocio
- Cálculos deterministas
- Explicaciones de IA

Los modelos de IA nunca deben reemplazar cálculos deterministas.

Su función es:

- Explicar
- Guiar
- Diagnosticar
- Resumir
- Asistir

---

# 12. Filosofía del proyecto

El framework no intenta "adivinar" resultados contables.

En su lugar, los reconstruye desde el historial de transacciones.

Cuando existe incertidumbre, debe ser reportada explícitamente.

El silencio nunca es preferido a la incertidumbre.

---

# 13. Visión a largo plazo

El proyecto debe evolucionar hacia un ecosistema completo compuesto de:

- Librería Python
- Interfaz de línea de comandos (CLI)
- API REST
- Servidor MCP
- Agentes de IA
- Base de conocimiento
- Motor de auditoría
- Motor de impuestos
- Motor de informes

Todos los componentes deben compartir las mismas reglas de negocio y base de conocimiento.

---

# 14. Gobernanza

Las decisiones arquitectónicas principales deben documentarse como Registros de Decisión de Arquitectura (ADR).

Cada cambio significativo debe:

1. Ser propuesto
2. Ser revisado
3. Ser documentado
4. Ser implementado
5. Ser probado

---

# 15. Regla central

**Una reconciliación nunca se considera completa porque "se vea correcta".**

Una reconciliación está completa solo cuando cada balance, tenencia y cálculo fiscal puede ser reproducido desde el historial completo de transacciones utilizando reglas deterministas.

Este principio tiene prioridad sobre todas las demás decisiones de diseño en el proyecto.
