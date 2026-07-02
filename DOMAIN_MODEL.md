# Modelo de Dominio - CoinTracking Expert Framework

**Especificación del Modelo de Dominio Completo**

---

## Propósito

Este documento define el modelo de dominio del framework CoinTracking Expert, proporcionando una especificación precisa de todas las entidades, sus responsabilidades, relaciones, restricciones y comportamientos. El modelo es la fuente de verdad para todas las decisiones arquitectónicas y de implementación.

El modelo de dominio describe cómo el sistema:
- Representa transacciones de criptomonedas
- Reconstruye saldos y tenencias
- Detecciona duplicados y transferencias huérfanas
- Calcula costos de adquisición usando FIFO
- Valida consistencia de datos
- Genera reportes de auditoría

---

## Principios de Diseño

### 1. **Separación de Responsabilidades**

Cada entidad tiene una única responsabilidad bien definida. Las transacciones representan eventos; los ledgers reconstruyen estado; los audits detectan inconsistencias.

### 2. **Independencia de Fuente**

El modelo es agnóstico respecto a fuentes de datos (CSV, API, manual). La normalización convierte todo a representación canónica.

### 3. **Reproducibilidad Determinista**

Idénticas entradas producen idénticas salidas. No hay estado no-determinista. Todas las decisiones están basadas en datos observables.

### 4. **Trazabilidad Completa**

Toda conclusión es rastreable hasta sus datos de origen. Cada finding documenta evidencia, causa e impacto.

### 5. **Inmutabilidad de Datos Históricos**

Las transacciones y eventos históricos nunca se modifican. Solo se pueden agregar nuevas transacciones o correcciones explícitas.

### 6. **Validación en Límites**

Las validaciones ocurren en los límites del sistema (import). El sistema interno confía en datos ya validados.

### 7. **Modelado Explícito de Incertidumbre**

La incertidumbre se modela explícitamente como findings. El silencio nunca implica corrección.

---

## Descripción General del Dominio

### Contexto de Negocio

CoinTracking Expert valida bases de datos de CoinTracking verificando que sean:
- **Completas**: Todas las transacciones están registradas
- **Consistentes**: Los saldos son matemáticamente correctos
- **Reconciliables**: Los transfers se pueden emparejar
- **Auditables**: Cada conclusión se puede reproducir

### Flujo de Procesos Principales

#### 1. Flujo de Auditoría

```
Exportación CoinTracking
           ↓
    Importación & Normalización
           ↓
    Motor de Reconciliación
           ↓
    Motores Especializados ────┬─→ Motor de Duplicados
                               ├─→ Motor de Transfers
                               ├─→ Motor de Ledger
                               ├─→ Motor de Tenencias
                               ├─→ Motor FIFO
                               └─→ Motor de Impuestos
           ↓
    Síntesis de Findings
           ↓
    Generación de Reportes
```

#### 2. Flujo de Reconciliación

```
Ledger Completo
       ↓
Procesamiento Cronológico
       ↓
Cálculo de Saldos
       ↓
Detección de Balances Negativos
       ↓
Validación de Consistencia
```

#### 3. Flujo FIFO

```
Transacciones Compra/Venta
       ↓
Ordenamiento Cronológico
       ↓
Asignación de Lotes
       ↓
Cálculo de Costo Base
       ↓
Detección de Compras Faltantes
```

#### 4. Flujo de Generación de Reportes

```
Findings Consolidados
       ↓
Deduplicación & Ranking
       ↓
Síntesis de Narrativa
       ↓
Formateo Multi-Formato
       ↓
Reportes Finales
```

---

## Entidades del Núcleo

### 1. Transaction (Transacción)

#### Propósito

Representa un evento atómico que afecta el saldo de un activo en una cuenta. Es la unidad fundamental del análisis.

#### Responsabilidades

- Registrar exactamente qué sucedió, cuándo y dónde
- Mantener trazabilidad con sistemas de origen
- Representar el evento en forma normalizada
- Participar en cálculos de saldo y FIFO

#### Atributos Requeridos

```
id: TransactionId
    - Identificador único único dentro del dominio
    - Generado por el sistema, no por fuente de datos
    - Inmutable

sourceId: SourceTransactionId
    - ID en sistema de origen (CoinTracking, exchange, blockchain)
    - Puede no ser único entre fuentes
    - Requerido para trazabilidad

timestamp: Timestamp
    - Momento exacto cuando ocurrió
    - Precisión: segundo o milisegundo
    - Zona horaria: UTC normalizado

transactionType: TransactionType
    - Enumeración: BUY, SELL, DEPOSIT, WITHDRAWAL, TRANSFER, 
                   STAKING_REWARD, AIRDROP, FEE, DUST, OTHER
    - No nulo

account: Account (referencia)
    - Cuenta donde ocurrió la transacción
    - Relación: muchas transacciones por cuenta

asset: Asset (referencia)
    - Qué se está moviendo
    - Bitcoin, Ethereum, USDT, etc.

quantity: Quantity
    - Cantidad transada
    - Puede ser positivo (entrada) o negativo (salida)
    - Precisión: hasta 18 decimales

price: Money (opcional)
    - Precio unitario al momento de la transacción
    - Moneda: típicamente USD

totalValue: Money (calculado)
    - quantity × price
    - Se calcula si price está disponible

fee: Fee (opcional)
    - Comisión pagada
    - Puede expresarse como cantidad o como dinero

feeAsset: Asset (opcional)
    - En qué activo se pagó la comisión
    - Puede diferir del asset principal

source: DataSource
    - Enumeración: COINTRACKING_CSV, COINTRACKING_API, 
                   EXCHANGE_API, BLOCKCHAIN, MANUAL
    - Indica de dónde vino el dato

status: TransactionStatus
    - Enumeración: CONFIRMED, PENDING, FAILED, DISPUTED
    - Por defecto: CONFIRMED

notes: String (opcional)
    - Anotaciones del usuario o del sistema
    - Para recordar contexto o razones

externalReference: Hash (opcional)
    - Hash de blockchain, ID de exchange, etc.
    - Para trazabilidad completa

```

#### Atributos Opcionales

```
description: String
    - Descripción legible
    - Ejemplo: "Buy ETH on Binance"

tags: List<String>
    - Para categorización flexible
    - Ejemplo: ["margin", "leveraged", "problematic"]

metadata: Map<String, Any>
    - Datos arbitrarios de origen
    - Preserva información específica de fuente

relatedTransactions: List<TransactionId>
    - IDs de transacciones relacionadas
    - Ejemplo: una venta relacionada con una compra

```

#### Restricciones

1. **Cantidad debe ser consistente con tipo**
   - BUY: quantity > 0
   - SELL: quantity < 0
   - DEPOSIT: quantity > 0
   - WITHDRAWAL: quantity < 0
   - TRANSFER: puede ser positivo o negativo según perspectiva
   - STAKING_REWARD: quantity > 0
   - AIRDROP: quantity > 0
   - FEE: quantity < 0

2. **Timestamp debe ser válido**
   - No puede ser en el futuro
   - No puede ser anterior a genesis de blockchain (genesis de Bitcoin: 2009-01-03)

3. **Account y Asset requeridas**
   - Ninguna transacción puede estar "huérfana"

4. **Precisión de Cantidad**
   - Máximo 18 decimales (estándar ERC-20)
   - No debe ser cero

5. **Fee no puede ser negativo**
   - Siempre costo para el usuario

#### Ciclo de Vida

```
Created → Imported → Validated → Processed → Finalized
                         ↓
                    On Error: Disputed
```

**Created**: Sistema crea instancia
**Imported**: Datos del origen
**Validated**: Pasa validaciones de esquema
**Processed**: Incluido en cálculos
**Finalized**: Inmutable para reportes
**Disputed**: Marcado como problemático

#### Relaciones

```
Transaction
    ├─→ Account (many-to-one)
    │      └─ Exactamente una cuenta
    │
    ├─→ Asset (many-to-one)
    │      └─ Exactamente un activo transado
    │
    ├─→ Asset (many-to-one, opcional)
    │      └─ Activo de comisión
    │
    ├─→ LedgerEntry (one-to-one)
    │      └─ Cada transacción tiene entrada en ledger
    │
    ├─→ Transfer (one-to-one, si aplica)
    │      └─ Si es TRANSFER, es parte de Transfer
    │
    ├─→ Trade (one-to-one, si aplica)
    │      └─ Si es BUY o SELL, es parte de Trade
    │
    └─→ List<Finding> (many-to-many)
           └─ Puede ser evidencia de múltiples findings
```

#### Invariantes

```
Invariante 1: Identificador único
    ∀ t1, t2 ∈ Transaction: t1.id = t2.id → t1 = t2

Invariante 2: Integridad de referencia
    ∀ t ∈ Transaction: t.account ≠ null ∧ t.asset ≠ null

Invariante 3: Timestamp consistente
    timestamp(t) ≤ now()

Invariante 4: Suma de comisión válida
    fee > 0 ∨ fee = null

Invariante 5: No puede modificarse después de validación
    status ∈ {CONFIRMED, FINALIZED} → inmutable
```

---

### 2. Asset (Activo)

#### Propósito

Representa una clase de criptomoneda o token. Es el "qué" de cada transacción.

#### Responsabilidades

- Identificar únicamente cada activo
- Mantener metadatos esenciales
- Facilitar identificación cruzada entre fuentes
- Participar en cálculos de cantidad y valor

#### Atributos Requeridos

```
symbol: AssetSymbol (identificador de dominio)
    - Ticker universal
    - Ejemplos: BTC, ETH, USDT, SOL, MATIC
    - Máximo 10 caracteres
    - Upper case siempre

name: String
    - Nombre legible
    - Ejemplos: "Bitcoin", "Ethereum", "Tether USD"

type: AssetType
    - Enumeración: CRYPTOCURRENCY, STABLECOIN, FIAT, TOKEN, NFT
    - Por defecto: CRYPTOCURRENCY

network: Network (opcional)
    - Red blockchain donde existe
    - Ejemplos: ETHEREUM, SOLANA, POLYGON
    - Nulo para activos sin red específica (BTC)

contractAddress: Address (opcional)
    - Dirección de contrato inteligente
    - Para tokens ERC-20, BEP-20, etc.
    - Específico de red

decimals: Integer
    - Decimales standard
    - Bitcoin: 8, Ethereum: 18, USDC: 6
    - Por defecto: 18

```

#### Atributos Opcionales

```
coingeckoId: String
    - ID en CoinGecko para obtener precios
    
coinmarketcapId: String
    - ID en CoinMarketCap

issuedDate: Timestamp
    - Cuándo se creó/lanzó

totalSupply: Quantity
    - Suministro total (si conocido)

website: String
    - URL oficial del proyecto

category: String
    - Categorización: DeFi, Layer-2, Stablecoin, etc.

deprecated: Boolean
    - Si el activo ya no es relevante

aliases: List<String>
    - Símbolos alternativos
    - Ejemplo: "Tether USD" también es "USDT.e" en Avalanche

```

#### Restricciones

1. **Symbol único en dominio**
   - Un BTC representa todos los bitcoins
   - Pero se puede tener WBTC (wrapped Bitcoin)

2. **Symbol válido**
   - Solo caracteres alfanuméricos
   - No espacios

3. **ContractAddress requiere Network**
   - Si está presente, Network también debe estarlo

4. **Decimals entre 0 y 18**
   - Prácticamente ningún asset tiene más de 18

#### Relaciones

```
Asset
    ├─→ Network (many-to-one, opcional)
    │
    ├─→ List<Transaction> (one-to-many)
    │      └─ Todos los movimientos de este activo
    │
    ├─→ List<Holding> (one-to-many)
    │      └─ Tenencias actuales
    │
    ├─→ List<AcquisitionLot> (one-to-many)
    │      └─ Lotes de adquisición
    │
    └─→ List<Price> (one-to-many)
           └─ Histórico de precios
```

#### Identidad

```
Asset se identifica por:
    - symbol: AssetSymbol
    - network: Network (si aplica)
    
Clave única compuesta: (symbol, network)

Ejemplos:
    - (BTC, null) = Bitcoin
    - (USDT, ETHEREUM) = Tether en Ethereum
    - (USDT, POLYGON) = Tether en Polygon
```

---

### 3. Exchange (Intercambio)

#### Propósito

Representa una plataforma de intercambio de criptomonedas. Agrupa todas las cuentas en esa plataforma.

#### Responsabilidades

- Identificar plataforma de origen
- Mantener metadatos de integración
- Facilitar auditoría de datos de origen
- Participar en matching de transferencias

#### Atributos Requeridos

```
id: ExchangeId (identificador de dominio)
    - Único, generado por sistema
    - Ejemplos: BINANCE, COINBASE, KRAKEN

name: String
    - Nombre legible
    - Ejemplos: "Binance", "Coinbase", "Kraken"

type: ExchangeType
    - Enumeración: CENTRALIZED, DECENTRALIZED, DEX, BRIDGE
    - Por defecto: CENTRALIZED

supportedAssets: List<Asset>
    - Activos que soporta
    - Para validación cruzada

supportedNetworks: List<Network>
    - Redes blockchain soportadas

```

#### Atributos Opcionales

```
apiEndpoint: String
    - URL de API (si integrado)

website: String
    - Sitio web oficial

country: String
    - País de regulación

regulatedStatus: RegulatoryStatus
    - Enumeración: REGULATED, UNREGULATED, PENDING

supportedDataFormats: List<DataFormat>
    - CSV, JSON, XML, etc.
    - Qué formatos de export ofrece

```

#### Relaciones

```
Exchange
    ├─→ List<Account> (one-to-many)
    │      └─ Todas las cuentas en este exchange
    │
    ├─→ List<Asset> (many-to-many)
    │      └─ Activos que soporta
    │
    └─→ List<Network> (many-to-many)
           └─ Redes que soporta
```

#### Identidad

Exchange se identifica por:
```
- id: ExchangeId (único en dominio)
```

---

### 4. Wallet (Billetera)

#### Propósito

Representa una billetera blockchain (hardware, software, multi-sig). Es una forma de custodiar activos.

#### Responsabilidades

- Identificar billetera específica
- Mantener direcciones blockchain
- Facilitar matching de transacciones blockchain
- Participar en auditoría de tenencias

#### Atributos Requeridos

```
id: WalletId (identificador de dominio)
    - Único en dominio
    - Generado por sistema

name: String
    - Nombre legible
    - Ejemplo: "Hardware Wallet 1", "Ledger Live"

type: WalletType
    - Enumeración: HARDWARE, SOFTWARE, MULTISIG, 
                   HARDWARE_ABSTRACTION_LAYER, OTHER
    
network: Network
    - Red blockchain primaria
    - Ejemplo: ETHEREUM, SOLANA

addresses: List<Address>
    - Direcciones blockchain de esta billetera
    - Puede tener múltiples direcciones

```

#### Atributos Opcionales

```
provider: String
    - Manufacturer o software
    - Ejemplo: "Ledger", "Trezor", "MetaMask"

isMultisig: Boolean
    - Si requiere múltiples firmas

requiredSignatures: Integer
    - Si es multisig, cuántas se requieren

totalAddresses: Integer
    - Cantidad de direcciones potenciales

notes: String
    - Anotaciones del usuario

```

#### Relaciones

```
Wallet
    ├─→ Network (many-to-one)
    │      └─ Red blockchain principal
    │
    ├─→ List<Address> (one-to-many)
    │      └─ Direcciones de esta billetera
    │
    └─→ List<Account> (one-to-many)
           └─ Cuentas asociadas (en CoinTracking)
```

---

### 5. Account (Cuenta)

#### Propósito

Representa un contenedor de transacciones. Típicamente es una cuenta en un exchange, billetera, o aggregador como CoinTracking.

Es el agregado raíz principal para transacciones.

#### Responsabilidades

- Agrupar transacciones coherentemente
- Mantener trazabilidad de origen
- Participar en cálculos de saldo
- Facilitar auditoría segregada por cuenta

#### Atributos Requeridos

```
id: AccountId (identificador de dominio)
    - Único en dominio
    - Generado por sistema

name: String
    - Nombre legible
    - Ejemplo: "Binance Main", "Ledger Wallet 1"

source: AccountSource
    - Enumeración: EXCHANGE, WALLET, COINTRACKING, AGGREGATOR
    - Tipo de origen

exchange: Exchange (referencia, si aplicable)
    - Si es EXCHANGE, referencia a Exchange
    
wallet: Wallet (referencia, si aplicable)
    - Si es WALLET, referencia a Wallet

transactions: List<Transaction> (composición)
    - Todas las transacciones de esta cuenta
    - Responsabilidad exclusiva de Account

ledger: Ledger
    - Ledger de esta cuenta
    - Reconstruye saldos de transacciones

```

#### Atributos Opcionales

```
externalId: String
    - ID en sistema de origen
    - Ejemplo: usuario de exchange

status: AccountStatus
    - Enumeración: ACTIVE, ARCHIVED, SUSPENDED, CLOSED
    - Por defecto: ACTIVE

lastSyncTime: Timestamp
    - Última sincronización de datos

importedFrom: String
    - Formato de importación original
    - "CoinTracking CSV", "Binance API", etc.

```

#### Restricciones

1. **Exactamente uno de Exchange o Wallet**
   - Una Account es O de un exchange O de una billetera
   - O ambos nulos para COINTRACKING/AGGREGATOR

2. **Transactions no pueden ser nulas**
   - Lista puede estar vacía, pero nunca null

3. **Ledger sincronizado con Transactions**
   - Invariante: ledger.entries proviene de transactions

#### Ciclo de Vida

```
Created → Imported → Synced → Active → Archived
                       ↓
                 (periódicamente updated)
```

#### Relaciones

```
Account (Agregado Raíz)
    ├─→ Exchange (many-to-one, opcional)
    │
    ├─→ Wallet (many-to-one, opcional)
    │
    ├─→ List<Transaction> (one-to-many, composición)
    │      └─ Propiedad exclusiva
    │
    ├─→ Ledger (one-to-one, composición)
    │      └─ Computado desde transactions
    │
    ├─→ List<Holding> (one-to-many)
    │      └─ Tenencias actuales
    │
    └─→ List<Transfer> (one-to-many)
           └─ Transfers a/desde esta cuenta
```

#### Invariantes

```
Invariante 1: Identidad
    ∀ a1, a2 ∈ Account: a1.id = a2.id → a1 = a2

Invariante 2: Exactamente una fuente
    (exchange ≠ null XOR wallet ≠ null) ∨ (exchange = null ∧ wallet = null)

Invariante 3: Transacciones inmutables históricamente
    Para t ∈ transactions donde t.timestamp < now - 90 días:
        t es inmutable

Invariante 4: Ledger derivado
    ledger.entries es función pura de transactions
```

---

### 6. Ledger (Libro Mayor)

#### Propósito

Reconstruye el saldo de cada activo en una cuenta a lo largo del tiempo. Es el estado derivado de todas las transacciones.

#### Responsabilidades

- Procesar transacciones cronológicamente
- Mantener saldos corrientes para cada activo
- Detectar estados imposibles (balances negativos)
- Facilitar validación de consistencia

#### Atributos Requeridos

```
id: LedgerId
    - Único, típicamente derivado de Account.id

account: Account (referencia)
    - A qué cuenta pertenece

entries: List<LedgerEntry> (composición)
    - Entrada para cada transacción
    - Ordenado cronológicamente
    - Nunca vacío si account tiene transacciones

assetBalances: Map<Asset, Quantity>
    - Saldo actual para cada activo
    - Derivado del último entry de cada activo
    - Actualizado incrementalmente

```

#### Atributos Opcionales

```
lastUpdated: Timestamp
    - Cuándo se calculó por última vez

isValid: Boolean
    - Si no hay estados imposibles
    - Por defecto true, se marca false si hay negatives

```

#### Restricciones

1. **Entries ordenadas cronológicamente**
   - entries[i].timestamp ≤ entries[i+1].timestamp

2. **Entries corresponden a Transactions**
   - Cada transaction tiene exactamente un entry
   - No hay entries huérfanas

3. **assetBalances consistentes**
   - Para cada asset: balance = sum de entries de ese asset

#### Ciclo de Vida

```
Created → Populated → Updated → Finalized
                        ↓
                  (on each new transaction)
```

#### Relaciones

```
Ledger (Agregado)
    ├─→ Account (many-to-one)
    │      └─ Exactamente una account
    │
    └─→ List<LedgerEntry> (one-to-many, composición)
           └─ Propiedad exclusiva
```

---

### 7. LedgerEntry (Entrada de Ledger)

#### Propósito

Representa una línea en el ledger: el impacto de una transacción en un saldo.

#### Responsabilidades

- Registrar saldo antes y después de transacción
- Facilitar trazabilidad de cambios de saldo
- Detectar imposibilidades (negativos)
- Participar en validación de consistencia

#### Atributos Requeridos

```
id: LedgerEntryId
    - Único, típicamente: account.id + transaction.id

ledger: Ledger (referencia)
    - A qué ledger pertenece

transaction: Transaction (referencia)
    - Qué transacción causó este entry

timestamp: Timestamp
    - Copiad de transaction.timestamp

asset: Asset (referencia)
    - Qué activo se vio afectado

quantity: Quantity
    - Cantidad movida (positivo o negativo)

balanceBefore: Quantity
    - Saldo antes de esta transacción

balanceAfter: Quantity
    - Saldo después de esta transacción
    - = balanceBefore + quantity

isoDate: Date
    - Fecha en UTC (para bucketing)

```

#### Atributos Opcionales

```
notes: String
    - Anotaciones de auditoria

```

#### Restricciones

1. **balanceAfter = balanceBefore + quantity**
   - Invariante de cálculo

2. **balanceBefore consistente con entry anterior**
   - Si existe entry anterior para mismo asset:
     entries[i-1].balanceAfter = entries[i].balanceBefore

3. **Cantidad no puede ser cero**
   - Cada entry representa cambio real

#### Relaciones

```
LedgerEntry
    ├─→ Ledger (many-to-one)
    │
    ├─→ Transaction (many-to-one)
    │      └─ Exactamente una transacción origen
    │
    └─→ Asset (many-to-one)
           └─ Exactamente un activo afectado
```

---

### 8. Transfer (Transferencia)

#### Propósito

Representa movimiento de activos entre dos cuentas. Modela la intención de matching entre withdrawal y deposit.

#### Responsabilidades

- Emparejar withdrawal/deposit correspondientes
- Detectar transferencias huérfanas
- Calcular tiempos de transferencia
- Participar en validación de consistencia

#### Atributos Requeridos

```
id: TransferId
    - Único en dominio

sourceTransaction: Transaction
    - Withdrawal desde cuenta origen

destinationTransaction: Transaction (opcional)
    - Deposit a cuenta destino
    - Puede ser null si no está emparejada

sourceAccount: Account
    - Cuenta que envía

destinationAccount: Account (opcional)
    - Cuenta que recibe
    - Nula si no se conoce destino

asset: Asset
    - Qué se está transfiriendo

quantity: Quantity
    - Cantidad enviada

status: TransferStatus
    - Enumeración: MATCHED, ORPHANED, PENDING, SUSPICIOUS

isMatched: Boolean
    - Si source y destination están emparejados

sourceTimestamp: Timestamp
    - Cuándo se envió (from sourceTransaction)

destinationTimestamp: Timestamp (opcional)
    - Cuándo se recibió (from destinationTransaction)

transferDuration: Duration (opcional)
    - Diferencia entre source y destination
    - Calculado si ambas timestamps existen

matchConfidence: Float
    - 0.0 a 1.0
    - Qué tan seguro es el matching

```

#### Atributos Opcionales

```
fee: Money
    - Comisión de transferencia

feeAsset: Asset
    - En qué activo se cobró

expectedQuantity: Quantity
    - Cantidad esperada a recibir
    - Puede diferir de quantity debido a fees

notes: String
    - Notas sobre matching

```

#### Restricciones

1. **Source y destination no pueden ser misma cuenta**
   - Transfer requiere dos cuentas diferentes

2. **Asset consistente**
   - sourceTransaction.asset = destinationTransaction.asset

3. **Quantity consistencia**
   - Si MATCHED:
     sourceTransaction.quantity = -destinationTransaction.quantity (approx)
   - Diferencia ≤ fee razonable

4. **Timestamps consistentes**
   - Si MATCHED:
     sourceTimestamp ≤ destinationTimestamp
   - destinationTimestamp - sourceTimestamp ≤ X horas
       (típicamente 24-48 horas)

#### Ciclo de Vida

```
Created → Candidate → Matched/Orphaned → Verified
                           ↓
                      On Investigation: Disputed
```

#### Relaciones

```
Transfer
    ├─→ Transaction (many-to-one)
    │      └─ sourceTransaction
    │
    ├─→ Transaction (many-to-one, opcional)
    │      └─ destinationTransaction
    │
    ├─→ Account (many-to-one)
    │      └─ sourceAccount
    │
    ├─→ Account (many-to-one, opcional)
    │      └─ destinationAccount
    │
    └─→ Asset (many-to-one)
           └─ Asset transferido
```

#### Invariantes

```
Invariante 1: Origen != Destino
    sourceAccount.id ≠ destinationAccount.id

Invariante 2: Transacciones válidas
    sourceTransaction.transactionType ∈ {WITHDRAWAL, TRANSFER}
    destinationTransaction.transactionType ∈ {DEPOSIT, TRANSFER}

Invariante 3: Si MATCHED, consistency
    quantity(source) + fee ≈ quantity(destination)
```

---

### 9. Trade (Comercio)

#### Propósito

Modela una transacción de compra-venta. Vincula una compra con su venta correspondiente para FIFO.

#### Responsabilidades

- Vincular compra y venta relacionadas
- Facilitar cálculo de ganancia/pérdida
- Participar en auditoría de lotes FIFO
- Calcular costo base

#### Atributos Requeridos

```
id: TradeId
    - Único en dominio

buyTransaction: Transaction
    - Transacción de compra
    - transactionType = BUY

sellTransaction: Transaction (opcional)
    - Transacción de venta
    - transactionType = SELL
    - Puede ser null (compra no vendida aún)

asset: Asset
    - Qué activo se compró

buyQuantity: Quantity
    - Cantidad comprada (positivo)

buyPrice: Money
    - Precio por unidad en compra

buyCost: Money
    - Costo total = buyQuantity × buyPrice

sellQuantity: Quantity (opcional)
    - Cantidad vendida (positivo)

sellPrice: Money (opcional)
    - Precio por unidad en venta

sellProceeds: Money (opcional)
    - Ingresos totales = sellQuantity × sellPrice

realizationGain: Money (opcional)
    - Ganancia realizada = sellProceeds - buyCost
    - Si existe sellTransaction

holding: Holding (referencia)
    - A qué tenencia actual corresponde
    - Si no vendida completamente

```

#### Atributos Opcionales

```
holdingQuantity: Quantity
    - Parte de la compra aún siendo held
    - = buyQuantity - sellQuantity (si parcial)

partialSales: List<Trade>
    - Otros Trades si fue venta parcial

exchangeUsed: Exchange
    - En dónde se hizo el trade

notes: String

```

#### Ciclos de Vida

```
Created → Paired (matched buy+sell) → Closed → Taxable
         ↓
      Holding (sin sell)
```

#### Relaciones

```
Trade
    ├─→ Transaction (many-to-one)
    │      └─ buyTransaction
    │
    ├─→ Transaction (many-to-one, opcional)
    │      └─ sellTransaction
    │
    ├─→ Asset (many-to-one)
    │
    └─→ Holding (many-to-one, opcional)
           └─ Si aún en holdings
```

---

### 10. Fee (Comisión)

#### Propósito

Modela una comisión aplicada a una transacción.

#### Responsabilidades

- Registrar costo de comisión
- Facilitar cálculo de costo total
- Participar en reconciliación
- Afectar cálculo de saldo

#### Atributos Requeridos

```
id: FeeId

transaction: Transaction (referencia)
    - A qué transacción corresponde

feeAsset: Asset
    - En qué activo se cobra
    - Ejemplo: USDT, BNB

feeQuantity: Quantity
    - Cantidad cobrada
    - Siempre positivo

feeInCurrency: Money (opcional)
    - Equivalente en moneda fiat
    - USD típicamente

feeType: FeeType
    - Enumeración: NETWORK, EXCHANGE, BRIDGE, OTHER

```

#### Atributos Opcionales

```
exchangeRate: Float
    - Tipo de cambio usado para conversión

feePercentage: Float
    - Si es porcentaje de cantidad

```

#### Relaciones

```
Fee
    ├─→ Transaction (many-to-one)
    │
    └─→ Asset (many-to-one)
           └─ feeAsset
```

---

### 11. Holding (Tenencia)

#### Propósito

Representa cantidad actual de un activo en una cuenta.

Es derivado, calculado desde Ledger.

#### Responsabilidades

- Representar estado actual
- Facilitar comparación con holdings reportados
- Participar en detección de discrepancias
- Servir como punto de referencia para auditoría

#### Atributos Requeridos

```
id: HoldingId
    - Tipicamente: account.id + asset.id

account: Account (referencia)

asset: Asset (referencia)

quantity: Quantity
    - Cantidad actual de este activo en esta account
    - Derivado del Ledger

cost: Money (opcional)
    - Costo de adquisición total (FIFO)

unrealizedGain: Money (opcional)
    - Ganancia no realizada = current_value - cost

status: HoldingStatus
    - Enumeración: ACTIVE, LIQUIDATED, DUST, WATCHED

lastUpdated: Timestamp
    - Cuándo se calculó por última vez

acquisitionLots: List<AcquisitionLot>
    - Lotes FIFO que componen esta tenencia

```

#### Restricciones

1. **Quantity derivada de Ledger**
   - = ledger.assetBalances[asset] para esta account

2. **Quantity no negativa**
   - quantity ≥ 0
   - Si < 0, es un error en ledger

3. **Cost solo si quantity > 0**

#### Relaciones

```
Holding
    ├─→ Account (many-to-one)
    │
    ├─→ Asset (many-to-one)
    │
    └─→ List<AcquisitionLot> (one-to-many)
           └─ Lotes FIFO
```

---

### 12. AcquisitionLot (Lote de Adquisición)

#### Propósito

Representa una compra específica que forma parte de una tenencia. Modela FIFO a nivel de lote.

#### Responsabilidades

- Registrar una adquisición específica
- Facilitar tracking de costo base
- Participar en cálculo de ganancia/pérdida
- Mantener trazabilidad de origen

#### Atributos Requeridos

```
id: AcquisitionLotId

holding: Holding (referencia)
    - A qué tenencia pertenece

buyTransaction: Transaction
    - Transacción original de compra
    - Inmutable

quantity: Quantity
    - Cantidad en este lote
    - Positivo

unitCost: Money
    - Costo por unidad

totalCost: Money
    - quantity × unitCost

acquisitionDate: Timestamp
    - Cuándo se adquirió

remainingQuantity: Quantity
    - Parte no vendida aún
    - ≤ quantity

soldQuantity: Quantity
    - Parte vendida
    - = quantity - remainingQuantity

```

#### Atributos Opcionales

```
sales: List<Disposal>
    - Disposiciones (ventas) de este lote

averageUnitPrice: Money
    - Si fueron múltiples compras agregadas

notes: String

```

#### Ciclo de Vida

```
Created → Acquired → Partially Sold → Fully Sold → Taxable
```

#### Relaciones

```
AcquisitionLot
    ├─→ Holding (many-to-one)
    │
    ├─→ Transaction (many-to-one)
    │      └─ buyTransaction
    │
    └─→ List<Disposal> (one-to-many, opcional)
           └─ Ventas de este lote
```

#### Invariantes

```
Invariante: Lotes no se superponen cronológicamente en mismo activo
    Para lotes en mismo holding, ordenados por acquisitionDate

Invariante: remainingQuantity + soldQuantity = quantity
```

---

### 13. Disposal (Disposición)

#### Propósito

Modela la venta o disposición de una parte de un AcquisitionLot.

#### Responsabilidades

- Vincular venta específica con lote FIFO
- Facilitar cálculo de ganancia/pérdida específica
- Participar en auditoría de lotes

#### Atributos Requeridos

```
id: DisposalId

acquisitionLot: AcquisitionLot (referencia)
    - Lote que se está vendiendo

sellTransaction: Transaction
    - Transacción de venta

quantity: Quantity
    - Cantidad vendida de este lote
    - ≤ acquisitionLot.quantity

unitPrice: Money
    - Precio por unidad en venta

totalProceeds: Money
    - quantity × unitPrice

gainLoss: Money
    - totalProceeds - (quantity × acquisitionLot.unitCost)

gainLossPercentage: Float
    - gainLoss / (quantity × acquisitionLot.unitCost)

disposalDate: Timestamp
    - Cuándo se vendió

```

#### Relaciones

```
Disposal
    ├─→ AcquisitionLot (many-to-one)
    │
    └─→ Transaction (many-to-one)
           └─ sellTransaction
```

---

### 14. TaxEvent (Evento Impositivo)

#### Propósito

Modela un evento que genera obligación tributaria.

#### Responsabilidades

- Registrar eventos taxables
- Facilitar cálculo de impuestos
- Participar en auditoría de consistencia fiscal
- Generar reportes de impuestos

#### Atributos Requeridos

```
id: TaxEventId

transaction: Transaction (referencia)
    - O vinculada a transacción

eventType: TaxEventType
    - Enumeración: CAPITAL_GAIN, CAPITAL_LOSS, INCOME, 
                   STAKING_REWARD, AIRDROP, HARD_FORK, EXPENSE

asset: Asset

quantity: Quantity
    - Cantidad involucrada

baseAmount: Money
    - Monto base para cálculo de impuesto

taxableAmount: Money
    - Monto sujeto a impuesto

jurisdiction: Jurisdiction
    - Enumeración: SPAIN, EU, USA, etc.

eventDate: Timestamp

frequency: TaxFrequency
    - Enumeración: ONCE, RECURRING_DAILY, RECURRING_MONTHLY, etc.

```

#### Relaciones

```
TaxEvent
    ├─→ Transaction (many-to-one)
    │
    ├─→ Asset (many-to-one)
    │
    └─→ Disposal (many-to-one, si aplica)
           └─ Para capital gains
```

---

### 15. Audit (Auditoría)

#### Propósito

Agregado raíz que representa una auditoría completa. Orquesta todo el proceso de validación.

#### Responsabilidades

- Coordinar todos los motores de auditoría
- Mantener findings consolidados
- Generar reportes
- Facilitar seguimiento de remediación

#### Atributos Requeridos

```
id: AuditId
    - Único en dominio

account: Account (referencia)
    - Qué cuenta se está auditando

createdAt: Timestamp
    - Cuándo se creó la auditoría

startedAt: Timestamp
    - Cuándo comenzó

completedAt: Timestamp (opcional)
    - Cuándo se completó

status: AuditStatus
    - Enumeración: CREATED, RUNNING, COMPLETED, FAILED

findings: List<Finding> (composición)
    - Todos los findings consolidados

rules: List<Rule> (referencia)
    - Qué reglas se aplicaron

engineResults: Map<String, EngineResult>
    - Resultados de cada motor

summary: AuditSummary
    - Resumen ejecutivo

```

#### Atributos Opcionales

```
configurationApplied: AuditConfiguration
    - Qué configuración se usó

notes: String

reportFormat: ReportFormat
    - En qué formato generar reporte

executedBy: String
    - Usuario o sistema que ejecutó

```

#### Ciclo de Vida

```
Created → Running → Completed → Reported → Archived
            ↓
         On Error: Failed → Investigation → Retried
```

#### Relaciones

```
Audit (Agregado Raíz)
    ├─→ Account (many-to-one)
    │
    ├─→ List<Finding> (one-to-many, composición)
    │      └─ Propiedad exclusiva
    │
    ├─→ List<Rule> (many-to-many)
    │      └─ Qué reglas aplicó
    │
    └─→ List<Report> (one-to-many)
           └─ Reportes generados
```

---

### 16. Finding (Hallazgo)

#### Propósito

Representa un problema, inconsistencia o hallazgo detectado durante auditoría.

#### Responsabilidades

- Registrar hallazgos específicos
- Documentar evidencia
- Facilitar trazabilidad
- Participar en síntesis de reportes

#### Atributos Requeridos

```
id: FindingId
    - Único en dominio

audit: Audit (referencia)
    - A qué auditoría pertenece

findingType: FindingType
    - Enumeración: DUPLICATE, ORPHANED_TRANSFER, 
                   NEGATIVE_BALANCE, MISSING_PURCHASE_HISTORY,
                   INCONSISTENCY, MISMATCH, WARNING, INFO

severity: Severity
    - Enumeración: CRITICAL, HIGH, MEDIUM, LOW, INFO

title: String
    - Título conciso del finding
    - Ejemplo: "Duplicate transaction detected"

description: String
    - Descripción detallada

affectedTransactions: List<Transaction>
    - Transacciones involucradas
    - Evidencia

evidence: List<Evidence>
    - Datos que sustentan el finding

cause: String
    - Por qué ocurrió

impact: String
    - Qué efecto tiene en auditoría

recommendedAction: String
    - Qué se debería hacer

detectedAt: Timestamp
    - Cuándo se detectó

relatedFindings: List<FindingId>
    - Otros findings relacionados

status: FindingStatus
    - Enumeración: OPEN, ACKNOWLEDGED, RESOLVED, DISMISSED

```

#### Atributos Opcionales

```
resolution: String
    - Cómo se resolvió (si RESOLVED)

resolvedAt: Timestamp
    - Cuándo se resolvió

riskLevel: RiskLevel
    - Enumeración: CRITICAL, HIGH, MEDIUM, LOW, NONE

financialImpact: Money
    - Impacto monetario estimado

investigationRequired: Boolean
    - Si necesita investigación manual

```

#### Restricciones

1. **Severity coherente con FindingType**
   - DUPLICATE → HIGH/MEDIUM
   - MISSING_PURCHASE_HISTORY → CRITICAL
   - INFO → debe ser INFO severity

2. **Evidence no vacía**
   - Siempre debe haber evidencia

3. **RecommendedAction no vacío**
   - Siempre se sugiere acción

#### Relaciones

```
Finding
    ├─→ Audit (many-to-one)
    │
    ├─→ List<Transaction> (many-to-many)
    │      └─ affectedTransactions
    │
    ├─→ List<Evidence> (one-to-many)
    │
    └─→ List<Rule> (many-to-many, opcional)
           └─ Qué regla(s) violó
```

---

### 17. Rule (Regla)

#### Propósito

Define una regla de validación que se aplica durante auditoría.

#### Responsabilidades

- Codificar lógica de validación
- Facilitar detección de problemas
- Permitir configuración de auditoría
- Documentar qué se espera

#### Atributos Requeridos

```
id: RuleId
    - Identificador único

name: String
    - Nombre descriptivo
    - Ejemplo: "No negative balances"

description: String
    - Qué valida

ruleExpression: RuleExpression
    - Expresión/código de validación

ruleType: RuleType
    - Enumeración: DATA_QUALITY, CONSISTENCY, RECONCILIATION, 
                   COMPLETENESS, PLAUSIBILITY

applicableTransactionTypes: List<TransactionType>
    - A qué tipos aplica
    - Vacío = aplica a todos

severity: RuleSeverity
    - Enumeración: CRITICAL, HIGH, MEDIUM, LOW, INFO
    - Severidad por defecto del finding

enabled: Boolean
    - Si se aplica actualmente

version: String
    - Versión de la regla
    - Para tracking de cambios

```

#### Atributos Opcionales

```
relatedRules: List<RuleId>
    - Reglas relacionadas

prerequisites: List<RuleId>
    - Reglas que deben ejecutarse primero

documentationUrl: String
    - Referencia a documentación

examples: List<RuleExample>
    - Ejemplos de qué detecta

```

#### Relaciones

```
Rule
    ├─→ List<Audit> (many-to-many)
    │      └─ Qué audits la aplicaron
    │
    └─→ List<Finding> (many-to-many)
           └─ Qué findings generó
```

---

### 18. Report (Reporte)

#### Propósito

Documento generado que sintetiza resultados de auditoría en forma presentable.

#### Responsabilidades

- Documentar hallazgos
- Presentar recomendaciones
- Facilitar seguimiento
- Generar en múltiples formatos

#### Atributos Requeridos

```
id: ReportId

audit: Audit (referencia)
    - De qué auditoría

generatedAt: Timestamp

reportFormat: ReportFormat
    - Enumeración: MARKDOWN, HTML, EXCEL, JSON

title: String

executiveSummary: String

findings: List<FindingEntry>
    - Findings resumidos

recommendations: List<Recommendation>
    - Acciones recomendadas

statistics: ReportStatistics
    - Estadísticas agregadas

content: String
    - Contenido en format especificado

```

#### Atributos Opcionales

```
generatedBy: String
    - Usuario/sistema que generó

version: String

metadata: Map<String, Any>
    - Metadata de generación

```

#### Relaciones

```
Report
    └─→ Audit (many-to-one)
           └─ Exactamente una auditoría origen
```

---

## Value Objects

Los Value Objects son inmutables y se identifican por sus valores, no por identidad.

### Money

```python
from dataclasses import dataclass
from decimal import Decimal
from enum import Enum

@dataclass(frozen=True)
class Money:
    amount: Decimal
    currency: str
    
    def __post_init__(self):
        if self.amount < 0:
            raise ValueError("Amount must be non-negative")
        try:
            float(self.amount)
        except (ValueError, OverflowError):
            raise ValueError("Amount must be a number")
    
    def plus(self, other: 'Money') -> 'Money':
        if self.currency != other.currency:
            raise ValueError(f"Cannot add {self.currency} to {other.currency}")
        return Money(self.amount + other.amount, self.currency)
    
    def times(self, factor: Decimal) -> 'Money':
        return Money(self.amount * factor, self.currency)
```

### Quantity

```python
from dataclasses import dataclass
from decimal import Decimal, ROUND_HALF_UP

@dataclass(frozen=True)
class Quantity:
    amount: Decimal
    decimals: int = 18
    
    def __post_init__(self):
        if not (0 <= self.decimals <= 18):
            raise ValueError("Decimals must be between 0 and 18")
        try:
            float(self.amount)
        except (ValueError, OverflowError):
            raise ValueError("Amount must be a number")
    
    @property
    def normalized_amount(self) -> Decimal:
        return self.amount.quantize(
            Decimal(10) ** -self.decimals,
            rounding=ROUND_HALF_UP
        )
```

### Timestamp

```python
from dataclasses import dataclass
from datetime import datetime, timezone

@dataclass(frozen=True)
class Timestamp:
    instant: datetime
    
    def __post_init__(self):
        if self.instant > datetime.now(timezone.utc):
            raise ValueError("Cannot be in future")
    
    @classmethod
    def now(cls) -> 'Timestamp':
        return cls(datetime.now(timezone.utc))
```

### CurrencyPair

```python
from dataclasses import dataclass

@dataclass(frozen=True)
class CurrencyPair:
    from_currency: str
    to_currency: str
    
    def __post_init__(self):
        if self.from_currency == self.to_currency:
            raise ValueError("Currencies must differ")
    
    def inverse(self) -> 'CurrencyPair':
        return CurrencyPair(self.to_currency, self.from_currency)
```

### Hash

```python
from dataclasses import dataclass
from enum import Enum

class HashAlgorithm(Enum):
    SHA256 = "sha256"
    KECCAK256 = "keccak256"
    BLAKE2B = "blake2b"

@dataclass(frozen=True)
class Hash:
    value: str
    algorithm: HashAlgorithm
    
    def __post_init__(self):
        if not self.value or self.value.isspace():
            raise ValueError("Hash value cannot be blank")
        if not (32 <= len(self.value) <= 128):
            raise ValueError(f"Hash length must be between 32 and 128 characters")
```

### Address

```python
from dataclasses import dataclass
import re
from enum import Enum

class Network(Enum):
    BITCOIN = "bitcoin"
    ETHEREUM = "ethereum"
    BNB_CHAIN = "bnb_chain"
    SOLANA = "solana"
    POLYGON = "polygon"
    ARBITRUM = "arbitrum"
    BASE = "base"

@dataclass(frozen=True)
class Address:
    value: str
    network: Network
    
    def __post_init__(self):
        if not self.value or self.value.isspace():
            raise ValueError("Address cannot be blank")
        self._validate(self.network)
    
    def _validate(self, network: Network) -> None:
        if network == Network.ETHEREUM:
            if not self.value.startswith("0x"):
                raise ValueError("Ethereum address must start with 0x")
            if len(self.value) != 42:
                raise ValueError("Ethereum address must be 42 characters long")
        elif network == Network.BITCOIN:
            if not re.match(r"^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$", self.value):
                raise ValueError("Invalid Bitcoin address format")
        # ... validación para otras redes
```

### TransactionId

```python
from dataclasses import dataclass
import uuid

@dataclass(frozen=True)
class TransactionId:
    value: str
    
    def __post_init__(self):
        if not self.value or self.value.isspace():
            raise ValueError("TransactionId cannot be blank")
        if not (1 <= len(self.value) <= 256):
            raise ValueError("TransactionId length must be between 1 and 256")
    
    @classmethod
    def generate(cls) -> 'TransactionId':
        return cls(str(uuid.uuid4()))
```

### TradeId

```python
from dataclasses import dataclass

@dataclass(frozen=True)
class TradeId:
    value: str
    
    def __post_init__(self):
        if not self.value or self.value.isspace():
            raise ValueError("TradeId cannot be blank")
```

### Fee

```python
from dataclasses import dataclass
from enum import Enum

class FeeType(Enum):
    NETWORK = "network"
    EXCHANGE = "exchange"
    BRIDGE = "bridge"
    OTHER = "other"

@dataclass(frozen=True)
class Fee:
    amount: Money
    fee_type: FeeType
    
    def __post_init__(self):
        if self.amount.amount < 0:
            raise ValueError("Fee must be non-negative")
```

### ExchangeId

```python
from dataclasses import dataclass

@dataclass(frozen=True)
class ExchangeId:
    value: str
    
    @classmethod
    def BINANCE(cls) -> 'ExchangeId':
        return cls("BINANCE")
    
    @classmethod
    def COINBASE(cls) -> 'ExchangeId':
        return cls("COINBASE")
    
    @classmethod
    def KRAKEN(cls) -> 'ExchangeId':
        return cls("KRAKEN")
    
    @classmethod
    def BYBIT(cls) -> 'ExchangeId':
        return cls("BYBIT")
    
    @classmethod
    def OKX(cls) -> 'ExchangeId':
        return cls("OKX")
    
    @classmethod
    def KUCOIN(cls) -> 'ExchangeId':
        return cls("KUCOIN")
    
    @classmethod
    def BINGX(cls) -> 'ExchangeId':
        return cls("BINGX")
```

### WalletId

```python
from dataclasses import dataclass
import uuid

@dataclass(frozen=True)
class WalletId:
    value: str
    
    @classmethod
    def generate(cls) -> 'WalletId':
        return cls(str(uuid.uuid4()))
```

### AssetSymbol

```python
from dataclasses import dataclass
import re

@dataclass(frozen=True)
class AssetSymbol:
    value: str
    
    def __post_init__(self):
        if not self.value or self.value.isspace():
            raise ValueError("AssetSymbol cannot be blank")
        if not re.match(r"^[A-Z0-9]{1,10}$", self.value):
            raise ValueError("AssetSymbol must be 1-10 uppercase letters or digits")
```

### Network

```python
from enum import Enum

class Network(Enum):
    BITCOIN = "bitcoin"
    ETHEREUM = "ethereum"
    BNB_CHAIN = "bnb_chain"
    SOLANA = "solana"
    POLYGON = "polygon"
    ARBITRUM = "arbitrum"
    BASE = "base"
    AVALANCHE = "avalanche"
    OPTIMISM = "optimism"
    FANTOM = "fantom"
    CRONOS = "cronos"
    HARMONY = "harmony"
```

### Status

Múltiples enums de status para diferentes entidades.

```python
from enum import Enum

class TransactionStatus(Enum):
    CONFIRMED = "confirmed"
    PENDING = "pending"
    FAILED = "failed"
    DISPUTED = "disputed"

class AccountStatus(Enum):
    ACTIVE = "active"
    ARCHIVED = "archived"
    SUSPENDED = "suspended"
    CLOSED = "closed"

class AuditStatus(Enum):
    CREATED = "created"
    RUNNING = "running"
    COMPLETED = "completed"
    FAILED = "failed"

class HoldingStatus(Enum):
    ACTIVE = "active"
    LIQUIDATED = "liquidated"
    DUST = "dust"
    WATCHED = "watched"

class TransferStatus(Enum):
    MATCHED = "matched"
    ORPHANED = "orphaned"
    PENDING = "pending"
    SUSPICIOUS = "suspicious"

class FindingStatus(Enum):
    OPEN = "open"
    ACKNOWLEDGED = "acknowledged"
    RESOLVED = "resolved"
    DISMISSED = "dismissed"
```

### Severity

```python
from enum import Enum

class Severity(Enum):
    CRITICAL = "critical"   # Requiere acción inmediata
    HIGH = "high"          # Problema significativo
    MEDIUM = "medium"      # Notable pero no urgente
    LOW = "low"            # Problema menor
    INFO = "info"          # Informacional
```

### RiskLevel

```python
from enum import Enum

class RiskLevel(Enum):
    CRITICAL = "critical"
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low"
    NONE = "none"
```

### FindingType

```python
from enum import Enum

class FindingType(Enum):
    DUPLICATE = "duplicate"                           # Transacción duplicada
    ORPHANED_TRANSFER = "orphaned_transfer"          # Transferencia sin par
    NEGATIVE_BALANCE = "negative_balance"            # Estado imposible
    MISSING_PURCHASE_HISTORY = "missing_purchase"    # Vendido más que comprado
    INCONSISTENCY = "inconsistency"                  # Desajuste de datos
    MISMATCH = "mismatch"                            # Esperado vs actual
    WARNING = "warning"                              # Inusual pero posible
    INFO = "info"                                    # Informacional
```

---

## Enumeraciones del Dominio

### TransactionType

```python
from enum import Enum

class TransactionType(Enum):
    BUY = "buy"                               # Compra de activo
    SELL = "sell"                             # Venta de activo
    DEPOSIT = "deposit"                       # Entrada de fondos
    WITHDRAWAL = "withdrawal"                 # Salida de fondos
    TRANSFER = "transfer"                     # Transferencia entre cuentas
    STAKING_REWARD = "staking_reward"         # Recompensa de staking
    AIRDROP = "airdrop"                       # Airdrop recibido
    FEE = "fee"                               # Comisión pagada
    DUST = "dust"                             # Cantidad insignificante
    INTERNAL_TRANSFER = "internal_transfer"   # Transferencia interna (mismo exchange)
    MARGIN_INTEREST = "margin_interest"       # Interés de margin
    DIVIDEND = "dividend"                     # Dividendo
    INFLATION = "inflation"                   # Inflación de token
    MERGE = "merge"                           # Merge de blockchain
    HARD_FORK = "hard_fork"                   # Fork de blockchain
    OTHER = "other"                           # Otro tipo no clasificado
```

### DataSource

```python
from enum import Enum

class DataSource(Enum):
    COINTRACKING_CSV = "cointracking_csv"         # Export CSV de CoinTracking
    COINTRACKING_API = "cointracking_api"         # Import desde API de CoinTracking
    EXCHANGE_API = "exchange_api"                 # API de exchange
    BLOCKCHAIN = "blockchain"                     # Blockchain explorer
    BLOCKCHAIN_INDEXER = "blockchain_indexer"     # Indexer como Etherscan
    WALLET_API = "wallet_api"                     # API de billetera
    MANUAL = "manual"                             # Entrada manual del usuario
    AGGREGATOR = "aggregator"                     # Servicio agregador
    IMPORT_FILE = "import_file"                   # Archivo importado
    UNKNOWN = "unknown"                           # Origen desconocido
```

### ExchangeType

```python
from enum import Enum

class ExchangeType(Enum):
    CENTRALIZED = "centralized"          # CEX (Binance, Coinbase, etc.)
    DECENTRALIZED = "decentralized"      # DEX en blockchain
    HYBRID = "hybrid"                    # Híbrido
    WALLET_PROVIDER = "wallet_provider"  # Proveedor de billetera
    BRIDGE = "bridge"                    # Bridge entre chains
    AGGREGATOR = "aggregator"            # Agregador
    UNKNOWN = "unknown"                  # Desconocido
```

### WalletType

```python
from enum import Enum

class WalletType(Enum):
    HARDWARE = "hardware"                        # Hardware wallet (Ledger, Trezor)
    SOFTWARE = "software"                        # Software wallet (MetaMask)
    MULTISIG = "multisig"                        # Multi-signature
    HARDWARE_ABSTRACTION = "hardware_abstraction"  # HAL (Smart contract wallet)
    EXCHANGE_CUSTODY = "exchange_custody"        # Custodia en exchange
    CUSTODIAL = "custodial"                      # Custodia terceros
    PAPER = "paper"                              # Billetera de papel
    OTHER = "other"                              # Otra
```

### AccountSource

```python
from enum import Enum

class AccountSource(Enum):
    EXCHANGE = "exchange"              # Cuenta en exchange
    WALLET = "wallet"                  # Billetera blockchain
    COINTRACKING = "cointracking"      # Portafolio CoinTracking
    AGGREGATOR = "aggregator"          # Servicio agregador
    VIRTUAL = "virtual"                # Cuenta virtual/imaginaria
    OTHER = "other"                    # Otra
```

### RuleType

```python
from enum import Enum

class RuleType(Enum):
    DATA_QUALITY = "data_quality"          # Validación de calidad
    CONSISTENCY = "consistency"            # Validación de consistencia
    RECONCILIATION = "reconciliation"      # Validación de reconciliación
    COMPLETENESS = "completeness"          # Validación de completitud
    PLAUSIBILITY = "plausibility"          # Validación de plausibilidad
    BUSINESS_LOGIC = "business_logic"      # Lógica de negocio
    COMPLIANCE = "compliance"              # Conformidad
    OTHER = "other"                        # Otra
```

### RuleSeverity

```python
from enum import Enum

class RuleSeverity(Enum):
    CRITICAL = "critical"      # Fallo crítico
    HIGH = "high"              # Fallo importante
    MEDIUM = "medium"          # Advertencia
    LOW = "low"                # Información
    INFO = "info"              # Solo informativo
```

### ReportFormat

```python
from enum import Enum

class ReportFormat(Enum):
    MARKDOWN = "markdown"      # Markdown
    HTML = "html"              # HTML
    EXCEL = "excel"            # Excel/XLSX
    JSON = "json"              # JSON
    PDF = "pdf"                # PDF
    CSV = "csv"                # CSV
    TEXT = "text"              # Texto plano
```

### TaxEventType

```python
from enum import Enum

class TaxEventType(Enum):
    CAPITAL_GAIN = "capital_gain"          # Ganancia de capital
    CAPITAL_LOSS = "capital_loss"          # Pérdida de capital
    INCOME = "income"                      # Ingreso (staking, etc.)
    STAKING_REWARD = "staking_reward"      # Recompensa staking
    AIRDROP = "airdrop"                    # Airdrop
    HARD_FORK = "hard_fork"                # Bifurcación
    DIVIDEND = "dividend"                  # Dividendo
    GIFT = "gift"                          # Regalo (no taxable)
    EXPENSE = "expense"                    # Gasto deducible
    OTHER = "other"                        # Otro
```

### Jurisdiction

```python
from enum import Enum

class Jurisdiction(Enum):
    SPAIN = "spain"
    EU = "eu"
    USA = "usa"
    UK = "uk"
    CANADA = "canada"
    AUSTRALIA = "australia"
    SINGAPORE = "singapore"
    JAPAN = "japan"
    UNKNOWN = "unknown"
```

---

## Agregados y Raíces de Agregado

### Agregado: Account

**Raíz de Agregado**: `Account`

**Entidades internas**:
- `Transaction` (propiedad)
- `Ledger` (propiedad, derivada)
- `LedgerEntry` (propiedad de Ledger)

**Value Objects internos**:
- `AccountId`
- `AccountStatus`
- `AccountSource`

**Invariantes de agregado**:

```
1. Una Account tiene exactamente uno o cero Exchange/Wallet
2. Cada Transaction pertenece a exactamente una Account
3. Cada Ledger pertenece a exactamente una Account
4. El Ledger es función pura de Transactions
5. Las Transactions de una Account son inmutables después de cierto tiempo
```

**Límites de agregado**:

- No modificar Transactions de otras Accounts
- No modificar Ledger directamente (solo recalcular)
- Transacciones de una Account son independientes de otras

**Repositorio**:

```python
from abc import ABC, abstractmethod
from typing import Optional, List

class AccountRepository(ABC):
    """Repositorio para gestión de cuentas."""
    
    @abstractmethod
    def save(self, account: Account) -> AccountId:
        """Guarda una nueva cuenta."""
        pass
    
    @abstractmethod
    def find_by_id(self, account_id: AccountId) -> Optional[Account]:
        """Encuentra una cuenta por ID."""
        pass
    
    @abstractmethod
    def find_by_exchange_id(self, exchange_id: ExchangeId) -> List[Account]:
        """Encuentra todas las cuentas para un exchange."""
        pass
    
    @abstractmethod
    def find_by_wallet_id(self, wallet_id: WalletId) -> List[Account]:
        """Encuentra todas las cuentas para una billetera."""
        pass
    
    @abstractmethod
    def update(self, account: Account) -> None:
        """Actualiza una cuenta existente."""
        pass
    
    @abstractmethod
    def delete(self, account_id: AccountId) -> None:
        """Elimina una cuenta."""
        pass
```

---

### Agregado: Audit

**Raíz de Agregado**: `Audit`

**Entidades internas**:
- `Finding` (propiedad)

**Value Objects internos**:
- `AuditId`
- `AuditStatus`

**Invariantes de agregado**:

```
1. Un Audit es para exactamente una Account
2. Cada Finding pertenece a un Audit
3. Los Findings no se pueden modificar después de COMPLETED
4. Los resultados de motores son inmutables
```

**Repositorio**:

```python
from abc import ABC, abstractmethod
from typing import Optional, List

class AuditRepository(ABC):
    """Repositorio para gestión de auditorías."""
    
    @abstractmethod
    def save(self, audit: Audit) -> AuditId:
        """Guarda una nueva auditoría."""
        pass
    
    @abstractmethod
    def find_by_id(self, audit_id: AuditId) -> Optional[Audit]:
        """Encuentra una auditoría por ID."""
        pass
    
    @abstractmethod
    def find_by_account_id(self, account_id: AccountId) -> List[Audit]:
        """Encuentra todas las auditorías para una cuenta."""
        pass
    
    @abstractmethod
    def update(self, audit: Audit) -> None:
        """Actualiza una auditoría existente."""
        pass
    
    @abstractmethod
    def delete(self, audit_id: AuditId) -> None:
        """Elimina una auditoría."""
        pass
```

---

## Contextos Acotados (Bounded Contexts)

El sistema se divide en varios contextos acotados:

### 1. **Importation Context**

Responsable de leer datos de fuentes externas y normalizarlos.

**Entidades**:
- SourceTransaction
- ImportMapping
- ValidationResult

**Servicios**:
- CsvImporter
- ApiImporter
- DataNormalizer
- SchemaValidator

**Fronteras**:
- Input: Raw data (CSV, JSON, API responses)
- Output: Normalized Transaction objects
- Límite: No debe contener lógica de auditoría

### 2. **Transaction Context**

Modela transacciones individuales.

**Agregados**:
- Transaction
- Transfer
- Trade

**Servicios**:
- TransactionValidator
- TransferMatcher
- TradeAnalyzer

**Fronteras**:
- Input: Normalized transactions
- Output: Validated transaction objects
- Límite: No modifica agregados de Account

### 3. **Ledger Context**

Reconstruye saldos desde transacciones.

**Agregados**:
- Ledger
- Account (parcialmente)

**Servicios**:
- LedgerReconstructor
- BalanceCalculator
- NegativeDetector

**Fronteras**:
- Input: Transactions ordered
- Output: Ledger with balance history
- Límite: Derivado (read-only de Transactions)

### 4. **Audit Context**

Coordina auditoría completa.

**Agregados**:
- Audit
- Finding
- Rule

**Servicios**:
- AuditOrchestrator
- RuleEngine
- FindingDetector

**Fronteras**:
- Input: Account + Configuration
- Output: Audit with Findings
- Límite: Orquestación

### 5. **Tax Context**

Calcula impuestos.

**Agregados**:
- TaxEvent
- Disposal
- AcquisitionLot (parcialmente)

**Servicios**:
- TaxCalculator
- GainLossComputer
- JurisdictionHandler

**Fronteras**:
- Input: Trades + Jurisdiction
- Output: Tax reports
- Límite: Específico de jurisdicción

### 6. **Reporting Context**

Genera reportes.

**Agregados**:
- Report

**Servicios**:
- ReportGenerator
- FindingFormatter
- StatisticsComputer

**Fronteras**:
- Input: Audit results
- Output: Reports in multiple formats
- Límite: Presentación

---

## Servicios del Dominio

Los servicios del dominio encapsulan lógica que cruza múltiples agregados.

### 1. **TransactionReconciliationService**

Valida transacciones y reconcilia con fuentes externas.

```python
from abc import ABC, abstractmethod
from typing import List

class TransactionReconciliationService(ABC):
    """Servicio para reconciliación de transacciones."""
    
    @abstractmethod
    def validate_transaction(self, transaction: Transaction) -> ValidationResult:
        """Valida una transacción individual."""
        pass
    
    @abstractmethod
    def reconcile_with_source(
        self,
        transactions: List[Transaction],
        source: DataSource
    ) -> ReconciliationResult:
        """Reconcilia transacciones con una fuente externa."""
        pass
    
    @abstractmethod
    def detect_duplicates(
        self,
        transactions: List[Transaction]
    ) -> List[DuplicatePair]:
        """Detecta transacciones duplicadas."""
        pass
```

### 2. **TransferMatchingService**

Empareja transfers entre cuentas.

```python
from abc import ABC, abstractmethod
from typing import List

class TransferMatchingService(ABC):
    """Servicio para emparejamiento de transferencias."""
    
    @abstractmethod
    def match_transfers(
        self,
        withdrawals: List[Transaction],
        deposits: List[Transaction]
    ) -> List[Transfer]:
        """Empareja retiros con depósitos."""
        pass
    
    @abstractmethod
    def find_orphaned_transfers(
        self,
        transactions: List[Transaction]
    ) -> List[Transaction]:
        """Encuentra transferencias huérfanas sin emparejar."""
        pass
```

### 3. **LedgerReconstructionService**

Reconstruye ledgers desde transacciones.

```python
from abc import ABC, abstractmethod
from typing import List

class LedgerReconstructionService(ABC):
    """Servicio para reconstrucción de libro mayor."""
    
    @abstractmethod
    def reconstruct_ledger(self, account: Account) -> Ledger:
        """Reconstruye el libro mayor de una cuenta."""
        pass
    
    @abstractmethod
    def detect_negative_balances(
        self,
        ledger: Ledger
    ) -> List[NegativeBalanceViolation]:
        """Detecta balances negativos (estados imposibles)."""
        pass
```

### 4. **FifoCalculationService**

Calcula FIFO y lotes de adquisición.

```python
from abc import ABC, abstractmethod
from typing import List

class FifoCalculationService(ABC):
    """Servicio para cálculo FIFO."""
    
    @abstractmethod
    def calculate_acquisition_lots(
        self,
        account: Account
    ) -> List[AcquisitionLot]:
        """Calcula lotes de adquisición usando FIFO."""
        pass
    
    @abstractmethod
    def calculate_disposals(
        self,
        lots: List[AcquisitionLot],
        sells: List[Transaction]
    ) -> List[Disposal]:
        """Calcula disposiciones (ventas) de lotes."""
        pass
    
    @abstractmethod
    def detect_missing_purchase_history(
        self,
        account: Account
    ) -> List[MissingPurchaseViolation]:
        """Detecta historial de compras faltante."""
        pass
```

### 5. **TaxCalculationService**

Calcula impuestos.

```python
from abc import ABC, abstractmethod
from typing import List

class TaxCalculationService(ABC):
    """Servicio para cálculo de impuestos."""
    
    @abstractmethod
    def calculate_tax_events(
        self,
        account: Account,
        jurisdiction: Jurisdiction
    ) -> List[TaxEvent]:
        """Calcula eventos fiscales para una cuenta en una jurisdicción."""
        pass
    
    @abstractmethod
    def calculate_capital_gains(
        self,
        disposals: List[Disposal]
    ) -> CapitalGainsReport:
        """Calcula ganancias de capital desde disposiciones."""
        pass
    
    @abstractmethod
    def validate_tax_calculations(
        self,
        reported: TaxReport,
        calculated: TaxReport
    ) -> TaxValidationResult:
        """Valida cálculos fiscales reportados vs calculados."""
        pass
```

### 6. **AuditOrchestrationService**

Orquesta auditoría completa.

```python
from abc import ABC, abstractmethod
from typing import Dict, List

class AuditOrchestrationService(ABC):
    """Servicio para orquestación de auditorías."""
    
    @abstractmethod
    def execute_audit(
        self,
        account: Account,
        configuration: AuditConfiguration
    ) -> Audit:
        """Ejecuta una auditoría completa."""
        pass
    
    @abstractmethod
    def synthesize_findings(
        self,
        engine_results: Dict[str, List[Finding]]
    ) -> List[Finding]:
        """Sintetiza hallazgos de múltiples motores."""
        pass
```

### 7. **ReportGenerationService**

Genera reportes en múltiples formatos.

```python
from abc import ABC, abstractmethod

class ReportGenerationService(ABC):
    """Servicio para generación de reportes."""
    
    @abstractmethod
    def generate_report(
        self,
        audit: Audit,
        format: ReportFormat
    ) -> Report:
        """Genera un reporte en el formato especificado."""
        pass
    
    @abstractmethod
    def generate_executive_summary(self, audit: Audit) -> str:
        """Genera un resumen ejecutivo."""
        pass
```

---

## Interfaces de Repositorio

```python
from abc import ABC, abstractmethod
from typing import Optional, List

class TransactionRepository(ABC):
    """Repositorio para transacciones."""
    
    @abstractmethod
    def save(self, transaction: Transaction) -> TransactionId:
        pass
    
    @abstractmethod
    def find_by_id(self, transaction_id: TransactionId) -> Optional[Transaction]:
        pass
    
    @abstractmethod
    def find_by_account_id(self, account_id: AccountId) -> List[Transaction]:
        pass
    
    @abstractmethod
    def find_by_asset(self, asset: Asset) -> List[Transaction]:
        pass
    
    @abstractmethod
    def find_by_timestamp_range(
        self,
        start: Timestamp,
        end: Timestamp
    ) -> List[Transaction]:
        pass
    
    @abstractmethod
    def update(self, transaction: Transaction) -> None:
        pass
    
    @abstractmethod
    def delete(self, transaction_id: TransactionId) -> None:
        pass


class AssetRepository(ABC):
    """Repositorio para activos."""
    
    @abstractmethod
    def save(self, asset: Asset) -> Asset:
        pass
    
    @abstractmethod
    def find_by_symbol(self, symbol: AssetSymbol) -> Optional[Asset]:
        pass
    
    @abstractmethod
    def find_by_network(self, network: Network) -> List[Asset]:
        pass
    
    @abstractmethod
    def find_all(self) -> List[Asset]:
        pass


class ExchangeRepository(ABC):
    """Repositorio para exchanges."""
    
    @abstractmethod
    def save(self, exchange: Exchange) -> ExchangeId:
        pass
    
    @abstractmethod
    def find_by_id(self, exchange_id: ExchangeId) -> Optional[Exchange]:
        pass
    
    @abstractmethod
    def find_all(self) -> List[Exchange]:
        pass


class WalletRepository(ABC):
    """Repositorio para billeteras."""
    
    @abstractmethod
    def save(self, wallet: Wallet) -> WalletId:
        pass
    
    @abstractmethod
    def find_by_id(self, wallet_id: WalletId) -> Optional[Wallet]:
        pass
    
    @abstractmethod
    def find_by_network(self, network: Network) -> List[Wallet]:
        pass
    
    @abstractmethod
    def find_by_address(self, address: Address) -> Optional[Wallet]:
        pass


class LedgerRepository(ABC):
    """Repositorio para libros mayores."""
    
    @abstractmethod
    def save(self, ledger: Ledger) -> LedgerId:
        pass
    
    @abstractmethod
    def find_by_account_id(self, account_id: AccountId) -> Optional[Ledger]:
        pass
    
    @abstractmethod
    def update(self, ledger: Ledger) -> None:
        pass


class HoldingRepository(ABC):
    """Repositorio para tenencias."""
    
    @abstractmethod
    def save(self, holding: Holding) -> HoldingId:
        pass
    
    @abstractmethod
    def find_by_account_id(self, account_id: AccountId) -> List[Holding]:
        pass
    
    @abstractmethod
    def find_by_asset(self, asset: Asset) -> List[Holding]:
        pass


class FindingRepository(ABC):
    """Repositorio para hallazgos de auditoría."""
    
    @abstractmethod
    def save(self, finding: Finding) -> FindingId:
        pass
    
    @abstractmethod
    def find_by_audit_id(self, audit_id: AuditId) -> List[Finding]:
        pass
    
    @abstractmethod
    def find_by_severity(self, severity: Severity) -> List[Finding]:
        pass
    
    @abstractmethod
    def update(self, finding: Finding) -> None:
        pass


class AcquisitionLotRepository(ABC):
    """Repositorio para lotes de adquisición."""
    
    @abstractmethod
    def save(self, lot: AcquisitionLot) -> AcquisitionLotId:
        pass
    
    @abstractmethod
    def find_by_holding_id(self, holding_id: HoldingId) -> List[AcquisitionLot]:
        pass
    
    @abstractmethod
    def find_by_account_id(self, account_id: AccountId) -> List[AcquisitionLot]:
        pass


class ReportRepository(ABC):
    """Repositorio para reportes."""
    
    @abstractmethod
    def save(self, report: Report) -> ReportId:
        pass
    
    @abstractmethod
    def find_by_audit_id(self, audit_id: AuditId) -> List[Report]:
        pass
    
    @abstractmethod
    def find_by_id(self, report_id: ReportId) -> Optional[Report]:
        pass
```

---

## Eventos del Dominio

Los eventos del dominio representan hechos significativos que ocurrieron.

```python
from dataclasses import dataclass
from abc import ABC, abstractmethod
from typing import Any, List

@dataclass
class DomainEvent(ABC):
    """Evento de dominio base."""
    occurred_at: Timestamp
    aggregate_id: Any
    
    def __post_init__(self):
        if not self.occurred_at:
            self.occurred_at = Timestamp.now()


@dataclass
class TransactionImportedEvent(DomainEvent):
    """Transacción importada."""
    transaction_id: TransactionId
    account_id: AccountId
    source: DataSource


@dataclass
class TransactionValidatedEvent(DomainEvent):
    """Transacción validada."""
    transaction_id: TransactionId
    is_valid: bool


@dataclass
class DuplicateDetectedEvent(DomainEvent):
    """Duplicado detectado."""
    transaction1_id: TransactionId
    transaction2_id: TransactionId
    audit_id: AuditId


@dataclass
class TransferMatchedEvent(DomainEvent):
    """Transferencia emparejada."""
    transfer_id: TransferId
    withdrawal_id: TransactionId
    deposit_id: TransactionId


@dataclass
class TransferOrphanedEvent(DomainEvent):
    """Transferencia huérfana."""
    transaction_id: TransactionId
    account_id: AccountId


@dataclass
class NegativeBalanceDetectedEvent(DomainEvent):
    """Balance negativo detectado."""
    account_id: AccountId
    asset: Asset
    balance: Quantity


@dataclass
class AuditCompletedEvent(DomainEvent):
    """Auditoría completada."""
    audit_id: AuditId
    finding_count: int
    critical_count: int


@dataclass
class ReportGeneratedEvent(DomainEvent):
    """Reporte generado."""
    report_id: ReportId
    audit_id: AuditId
    format: ReportFormat


class DomainEventPublisher(ABC):
    """Publicador de eventos del dominio."""
    
    @abstractmethod
    def publish(self, event: DomainEvent) -> None:
        """Publica un evento único."""
        pass
    
    @abstractmethod
    def publish_all(self, events: List[DomainEvent]) -> None:
        """Publica múltiples eventos."""
        pass
```

---

## Reglas de Validación

### Validaciones de Entidad

#### Transaction

```python
from dataclasses import dataclass
from typing import List, Optional

@dataclass
class ValidationResult:
    """Resultado de validación."""
    errors: List[str]
    
    @property
    def is_valid(self) -> bool:
        """Devuelve True si no hay errores."""
        return len(self.errors) == 0


class TransactionValidator:
    """Validador de transacciones."""
    
    @staticmethod
    def validate(transaction: Transaction) -> ValidationResult:
        """Valida una transacción."""
        errors = []
        
        if transaction.quantity == 0:
            errors.append("La cantidad no puede ser cero")
        
        if not TransactionValidator._is_valid_timestamp(transaction.timestamp):
            errors.append("Timestamp inválido")
        
        if transaction.fee and transaction.fee.amount < 0:
            errors.append("La comisión no puede ser negativa")
        
        quantity_error = TransactionValidator._validate_quantity_consistency(transaction)
        if quantity_error:
            errors.append(quantity_error)
        
        return ValidationResult(errors=errors)
    
    @staticmethod
    def _validate_quantity_consistency(t: Transaction) -> Optional[str]:
        """Valida que la cantidad sea consistente con el tipo."""
        if t.transaction_type == TransactionType.BUY:
            if t.quantity <= 0:
                return "Cantidad de BUY debe ser positiva"
        elif t.transaction_type == TransactionType.SELL:
            if t.quantity >= 0:
                return "Cantidad de SELL debe ser negativa"
        elif t.transaction_type == TransactionType.DEPOSIT:
            if t.quantity <= 0:
                return "Cantidad de DEPOSIT debe ser positiva"
        elif t.transaction_type == TransactionType.WITHDRAWAL:
            if t.quantity >= 0:
                return "Cantidad de WITHDRAWAL debe ser negativa"
        
        return None
    
    @staticmethod
    def _is_valid_timestamp(timestamp: Timestamp) -> bool:
        """Valida que el timestamp sea válido."""
        return timestamp.instant <= Timestamp.now().instant
```

#### Account

```python
class AccountValidator:
    """Validador de cuentas."""
    
    @staticmethod
    def validate(account: Account) -> ValidationResult:
        """Valida una cuenta."""
        errors = []
        
        if not account.transactions and account.ledger.entries:
            errors.append("Transacciones vacías pero libro mayor no")
        
        if account.exchange and account.wallet:
            errors.append("Una cuenta no puede tener ambos exchange y billetera")
        
        consistency_error = AccountValidator._validate_ledger_consistency(account)
        if consistency_error:
            errors.append(consistency_error)
        
        return ValidationResult(errors=errors)
    
    @staticmethod
    def _validate_ledger_consistency(account: Account) -> Optional[str]:
        """Valida que el libro mayor sea consistente con las transacciones."""
        reconstructed = LedgerReconstructionService.reconstruct_ledger(account)
        if len(reconstructed.entries) != len(account.ledger.entries):
            return "Desajuste en número de entradas del libro mayor"
        return None
```

#### Ledger

```python
from decimal import Decimal
from typing import List

class LedgerValidator:
    """Validador de libro mayor."""
    
    @staticmethod
    def validate(ledger: Ledger) -> ValidationResult:
        """Valida un libro mayor."""
        errors = []
        
        chrono_errors = LedgerValidator._validate_chronological_order(ledger)
        errors.extend(chrono_errors)
        
        calc_errors = LedgerValidator._validate_balance_calculations(ledger)
        errors.extend(calc_errors)
        
        negative_errors = LedgerValidator._detect_negative_balances(ledger)
        errors.extend(negative_errors)
        
        return ValidationResult(errors=errors)
    
    @staticmethod
    def _validate_balance_calculations(ledger: Ledger) -> List[str]:
        """Valida que los cálculos de balance sean correctos."""
        errors = []
        previous_balance = Decimal('0')
        
        for entry in ledger.entries:
            if entry.balance_before != previous_balance:
                errors.append(f"Desajuste de balance en entrada {entry.id}")
            
            expected_balance_after = entry.balance_before + entry.quantity
            if entry.balance_after != expected_balance_after:
                errors.append(f"Error de cálculo de balance en {entry.id}")
            
            previous_balance = entry.balance_after
        
        return errors
    
    @staticmethod
    def _validate_chronological_order(ledger: Ledger) -> List[str]:
        """Valida que los entries estén en orden cronológico."""
        errors = []
        for i in range(1, len(ledger.entries)):
            if ledger.entries[i].timestamp < ledger.entries[i-1].timestamp:
                errors.append(f"Orden cronológico incorrecto en entrada {ledger.entries[i].id}")
        return errors
    
    @staticmethod
    def _detect_negative_balances(ledger: Ledger) -> List[str]:
        """Detecta balances negativos (estados imposibles)."""
        errors = []
        for entry in ledger.entries:
            if entry.balance_after < 0:
                errors.append(f"Balance negativo detectado en entrada {entry.id}: {entry.balance_after}")
        return errors
```

---

## Invariantes del Dominio

### Nivel de Sistema

```
Invariante 1: Identidad Global Única
    ∀ entidades: id es único en todo el dominio

Invariante 2: Integridad Referencial
    Toda referencia a un agregado es válida

Invariante 3: Inmutabilidad Histórica
    Los datos históricos no se modifican

Invariante 4: Trazabilidad Completa
    Todo resultado es rastreable a evidencia

Invariante 5: Reproducibilidad
    Idénticos inputs → Idénticos outputs
```

### Nivel de Agregado

```
Agregado Account:
    - Exactamente una Account por (Exchange, Wallet) pair
    - Transactions de una Account no cruzan a otra
    - Ledger es función pura de Transactions

Agregado Audit:
    - Un Audit por Account per período
    - Findings no se modifican post-completion
    - Todos los Findings tienen evidencia
```

### Nivel de Entidad

```
Transaction:
    - Quantity consistente con TransactionType
    - Timestamp en rango válido
    - References válidas a Account y Asset

Transfer:
    - Si MATCHED: ambos source y destination existen
    - Quantity razonablemente cerca
    - Timestamps en orden

AcquisitionLot:
    - remainingQuantity + soldQuantity = quantity
    - unitCost > 0
    - acquisitionDate anterior a todas las disposiciones
```

---

## Restricciones de Negocio

### Validaciones Temporales

```
1. Transacciones no pueden ser futuras
2. Transferencias: deposit dentro de 48 horas de withdrawal
3. Blockchain genesis date: no anterior
4. Tax events: año fiscal consistente
```

### Validaciones de Cantidad

```
1. Quantities positivas excepto para SELL/WITHDRAWAL
2. Precisión máxima: 18 decimales
3. No se aceptan cantidades infinitas o NaN
4. Fee no negativa
```

### Validaciones de Consistencia

```
1. Ledger reconciliado con transacciones
2. Holdings = suma de acquisition lots
3. Transfers: entrada ≈ salida (within fees)
4. Tax calculations: reproducibles desde transacciones
```

### Validaciones de Integridad

```
1. Asset referenciado debe existir
2. Account referenciada debe existir
3. Exchange/Wallet referenciadas deben existir
4. Ningún dato huérfano
```

---

## Reglas de Identidad

Cada entidad se identifica de forma única dentro del dominio:

```
Transaction:
    - sourceId + account.id + timestamp es único
      (evita duplicados de import)
    - Pero system.id es PK global

Account:
    - (exchange.id, user.id) si exchange
    - (wallet.id, user.id) si wallet
    - Derivado: address si wallet

Asset:
    - symbol + network es key natural
    - Pero system.id es PK global

Transfer:
    - Matched por: (withdrawal, deposit) pair
    - Orphaned por: withdrawal sin deposit matched

Holding:
    - account + asset es key natural
    - Solo una holding per account-asset

AcquisitionLot:
    - buyTransaction + holding es unique
    - No hay duplicados de compra

Finding:
    - audit + type + affectedTransaction(s) es unique
    - Previene findings duplicados en mismo audit
```

---

## Ejemplos de Interacción entre Entidades

### 1. Flujo de Auditoría Completa

```
1. Usuario carga CoinTracking export
   ↓
2. Sistema importa y crea Transactions
   ↓
3. Para cada Account:
   ├─ Reconstruye Ledger desde Transactions
   ├─ Detecta balances negativos
   ├─ Busca duplicados
   ├─ Empareja Transfers
   ├─ Calcula Acquisition Lots (FIFO)
   ├─ Calcula Tax Events
   ├─ Genera Findings
   ↓
4. Consolida Findings en Audit
   ↓
5. Genera Reports en múltiples formatos
```

### 2. Reconciliación de Transferencias

```
Transaction WITHDRAWAL(200 BTC, from Account A)
    → Transfer.sourceTransaction
    ↓
Busca DEPOSIT(~200 BTC, to Account B)
    dentro de 48 horas
    ↓
Si encontrada:
    → Transfer.destinationTransaction
    → Transfer.status = MATCHED
    → Transfer.matchConfidence = 0.95
    
Si no encontrada:
    → Transfer.status = ORPHANED
    → Crea Finding "Orphaned Transfer"
```

### 3. Cálculo FIFO

```
Account.transactions ordenadas por timestamp:
    [BUY 1 BTC @ $30k, BUY 1 BTC @ $35k, SELL 0.5 BTC @ $50k]
    ↓
Para cada BUY crea AcquisitionLot:
    Lot 1: 1 BTC @ $30k, quantity=1
    Lot 2: 1 BTC @ $35k, quantity=1
    ↓
Para SELL crea Disposal (FIFO):
    Toma de Lot 1: 0.5 BTC @ $30k
    Crea Disposal:
        - quantity: 0.5
        - cost: 0.5 * $30k = $15k
        - proceeds: 0.5 * $50k = $25k
        - gain: $25k - $15k = $10k
    ↓
Lot 1.remainingQuantity = 0.5 BTC
    ↓
Crea Holding:
    - asset: BTC
    - quantity: 1.5 BTC
    - acquisitionLots: [Lot1(partial), Lot2(full)]
    - cost: 1.5 * $32.5k = $48.75k
```

### 4. Generación de Reportes

```
Audit completada con Findings
    ↓
Síntesis de Findings:
    ├─ Deduplicación (mismo issue reportado 2x)
    ├─ Ranking por severity
    ├─ Agrupamiento por tipo
    ↓
Para cada ReportFormat:
    ├─ MARKDOWN: Lista formateada, legible
    ├─ HTML: Versión web, tablas interactivas
    ├─ EXCEL: Datos estructurados, filtrable
    ├─ JSON: API consumible
    ↓
Genera ReportStatistics:
    - Total findings: 47
    - CRITICAL: 3
    - HIGH: 12
    - MEDIUM: 18
    - LOW: 14
    ↓
Incluye ExecutiveSummary:
    "Audit found 3 critical issues:
     - 2 duplicate transactions affecting $50k
     - 1 orphaned transfer of 0.5 BTC"
    ↓
Genera Report object con content
```

---

## Diagramas de Clase

### Diagrama Principal de Agregados

```mermaid
graph TB
    subgraph "Account Aggregate"
        Account["Account<br/>─ id: AccountId<br/>─ name: String<br/>─ source: Source<br/>─ exchange: Exchange<br/>─ wallet: Wallet"]
        Transaction["Transaction<br/>─ id: TxId<br/>─ timestamp<br/>─ type: TxType<br/>─ quantity<br/>─ price<br/>─ fee"]
        Ledger["Ledger<br/>─ id: LedgerId<br/>─ entries: List"]
        LedgerEntry["LedgerEntry<br/>─ transaction<br/>─ balanceBefore<br/>─ balanceAfter"]
        
        Account -->|1:N| Transaction
        Account -->|1:1| Ledger
        Ledger -->|1:N| LedgerEntry
        Transaction -->|referenced by| LedgerEntry
    end
    
    subgraph "Audit Aggregate"
        Audit["Audit<br/>─ id: AuditId<br/>─ account<br/>─ status: AuditStatus<br/>─ findings: List"]
        Finding["Finding<br/>─ id: FindingId<br/>─ type: FindingType<br/>─ severity<br/>─ description"]
        Rule["Rule<br/>─ id: RuleId<br/>─ name<br/>─ ruleType<br/>─ enabled"]
        
        Audit -->|1:N| Finding
        Audit -->|N:M| Rule
    end
    
    subgraph "FIFO Aggregate"
        Holding["Holding<br/>─ id: HoldingId<br/>─ account<br/>─ asset<br/>─ quantity"]
        AcqLot["AcquisitionLot<br/>─ buyTx<br/>─ quantity<br/>─ unitCost"]
        Disposal["Disposal<br/>─ sellTx<br/>─ quantity<br/>─ gain"]
        
        Holding -->|1:N| AcqLot
        AcqLot -->|1:N| Disposal
    end
    
    Account -->|validated by| Audit
    Transaction -->|referenced by| Finding
    Holding -->|derived from| Account
    Transaction -->|part of| Holding
```

### Diagrama de Transferencias

```mermaid
graph LR
    AccountA["Account A"]
    AccountB["Account B"]
    WithdrawalTx["Transaction WITHDRAWAL<br/>200 BTC<br/>Timestamp: T1"]
    DepositTx["Transaction DEPOSIT<br/>200 BTC<br/>Timestamp: T2"]
    Transfer["Transfer<br/>─ source: Withdrawal<br/>─ destination: Deposit<br/>─ matched: true"]
    
    AccountA -->|contains| WithdrawalTx
    AccountB -->|contains| DepositTx
    WithdrawalTx -->|source| Transfer
    DepositTx -->|destination| Transfer
    Transfer -->|validation| Finding["Finding: Matched OK"]
```

---

## Conclusión

Este modelo de dominio proporciona la especificación completa para implementar el CoinTracking Expert Framework. Define:

✅ **Entidades** con responsabilidades y restricciones claras
✅ **Agregados** con límites precisos y invariantes
✅ **Value Objects** para mayor seguridad de tipos
✅ **Servicios de Dominio** para lógica que cruza agregados
✅ **Eventos del Dominio** para comunicación entre sistemas
✅ **Repositorios** para persistencia
✅ **Contextos Acotados** para separación de concerns

Un implementador puede usar este documento como especificación completa sin tomar decisiones de negocio adicionales.
