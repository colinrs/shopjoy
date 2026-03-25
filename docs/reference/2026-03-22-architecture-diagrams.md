# Architecture Diagrams

> Visual documentation of the ShopJoy system architecture

## System Overview

```mermaid
graph TB
    subgraph "Client Layer"
        AdminUI[Admin Dashboard<br/>Vue 3 + Element Plus]
        ShopUI[Shop Frontend<br/>Vue 3 + Tailwind]
        MobileApp[Mobile App]
    end

    subgraph "API Gateway Layer"
        Gateway[Kong / nginx<br/>Rate Limiting, SSL]
    end

    subgraph "Application Layer"
        AdminAPI[Admin API<br/>Go + go-zero<br/>:8888]
        ShopAPI[Shop API<br/>Go + go-zero<br/>:8889]
    end

    subgraph "Domain Layer"
        ProductDomain[Product Domain<br/>Entities, Value Objects]
        OrderDomain[Order Domain<br/>Cart, Payment]
        InventoryDomain[Inventory Domain<br/>Warehouses, Stock]
    end

    subgraph "Infrastructure Layer"
        MySQL[(MySQL 8.0<br/>Primary + Replicas)]
        Redis[(Redis 7.0<br/>Cache + Sessions)]
        S3[S3 Storage<br/>Images, Assets]
    end

    AdminUI --> Gateway
    ShopUI --> Gateway
    MobileApp --> Gateway
    Gateway --> AdminAPI
    Gateway --> ShopAPI
    AdminAPI --> ProductDomain
    AdminAPI --> InventoryDomain
    ShopAPI --> OrderDomain
    ShopAPI --> ProductDomain
    ProductDomain --> MySQL
    InventoryDomain --> MySQL
    OrderDomain --> MySQL
    AdminAPI --> Redis
    ShopAPI --> Redis
    AdminAPI --> S3
    ShopAPI --> S3
```

---

## DDD Layered Architecture

```mermaid
graph TB
    subgraph "Interface Layer"
        Handler[HTTP Handlers<br/>Request/Response]
        Middleware[Middleware<br/>Auth, Tenant, Logging]
    end

    subgraph "Application Layer"
        Logic[Business Logic<br/>Use Case Orchestration]
        Service[Application Services<br/>DTO Conversion]
    end

    subgraph "Domain Layer"
        Entity[Entities<br/>Product, SKU, Order]
        VO[Value Objects<br/>Money, Status]
        Repo[Repository Interfaces<br/>Contracts]
        Event[Domain Events<br/>EventBus]
    end

    subgraph "Infrastructure Layer"
        RepoImpl[Repository Impl<br/>GORM, SQL]
        External[External Services<br/>Payment, Shipping]
        Cache[Cache Layer<br/>Redis]
    end

    Handler --> Logic
    Middleware --> Handler
    Logic --> Service
    Service --> Entity
    Service --> VO
    Entity --> Repo
    RepoImpl --> Repo
    RepoImpl --> MySQL[(MySQL)]
    Cache --> Redis[(Redis)]
```

---

## Multi-Tenant Architecture

```mermaid
graph TB
    subgraph "Request Flow"
        Request[HTTP Request]
        JWT[JWT Token<br/>with tenant_id]
        Context[Tenant Context<br/>context.Value]
        Query[DB Query<br/>WHERE tenant_id = ?]
    end

    subgraph "Data Isolation"
        TenantA[Tenant A Data]
        TenantB[Tenant B Data]
        TenantC[Tenant C Data]
    end

    subgraph "Shared Infrastructure"
        DB[(Shared Database)]
        Cache[(Shared Cache<br/>Namespaced Keys)]
    end

    Request --> JWT
    JWT --> Context
    Context --> Query
    Query --> TenantA
    Query --> TenantB
    Query --> TenantC
    TenantA --> DB
    TenantB --> DB
    TenantC --> DB
    TenantA --> Cache
    TenantB --> Cache
    TenantC --> Cache
```

### Tenant Context Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant Middleware
    participant Handler
    participant Logic
    participant Repository
    participant Database

    Client->>Gateway: Request + JWT
    Gateway->>Middleware: Forward Request
    Middleware->>Middleware: Validate JWT
    Middleware->>Middleware: Extract tenant_id
    Middleware->>Handler: Context with tenant_id
    Handler->>Logic: Execute Logic
    Logic->>Repository: Query with context
    Repository->>Repository: Get tenant_id from context
    Repository->>Database: SELECT ... WHERE tenant_id = ?
    Database-->>Repository: Results
    Repository-->>Logic: Entities
    Logic-->>Handler: Response
    Handler-->>Client: HTTP Response
```

---

## Product Domain Model

```mermaid
classDiagram
    class Product {
        +int64 ID
        +TenantID TenantID
        +string Name
        +Money Price
        +Status Status
        +int Stock
        +int64 CategoryID
        +PutOnSale() error
        +TakeOffSale() error
        +UpdateStock(qty int) error
    }

    class SKU {
        +int64 ID
        +int64 ProductID
        +string Code
        +Money Price
        +int Stock
        +int AvailableStock
        +int LockedStock
        +map Attributes
        +LockStock(qty int) error
        +DeductStock(qty int) error
    }

    class Category {
        +int64 ID
        +int64 ParentID
        +string Name
        +int Level
        +int Sort
        +Enable() void
        +Disable() void
    }

    class Brand {
        +int64 ID
        +string Name
        +string Logo
        +bool EnablePage
        +Enable() void
        +Disable() void
    }

    class Money {
        +int64 Amount
        +string Currency
        +Add(other Money) Money
        +Subtract(other Money) Money
    }

    class Status {
        <<enumeration>>
        DRAFT
        ON_SALE
        OFF_SALE
        DELETED
    }

    Product "1" --> "*" SKU : has variants
    Product --> Status
    Product --> Money : price
    Product --> Category : belongs to
    Product --> Brand : has brand
    SKU --> Money : price
```

---

## Product Status State Machine

```mermaid
stateDiagram-v2
    [*] --> Draft: Create Product
    Draft --> OnSale: PutOnSale()<br/>requires stock > 0
    Draft --> Deleted: SoftDelete()
    OnSale --> OffSale: TakeOffSale()
    OnSale --> Deleted: SoftDelete()
    OffSale --> OnSale: PutOnSale()<br/>requires stock > 0
    OffSale --> Deleted: SoftDelete()
    Deleted --> [*]

    note right of Draft
        Initial state
        Not visible to customers
    end note

    note right of OnSale
        Visible in shop
        Purchasable
    end note

    note right of OffSale
        Hidden from shop
        Data preserved
    end note

    note right of Deleted
        Terminal state
        Cannot be restored
    end note
```

---

## Inventory Domain Model

```mermaid
classDiagram
    class Warehouse {
        +int64 ID
        +string Code
        +string Name
        +string Country
        +bool IsDefault
        +Enable() void
        +Disable() void
    }

    class WarehouseInventory {
        +int64 ID
        +string SKUCode
        +int64 WarehouseID
        +int AvailableStock
        +int LockedStock
        +TotalStock() int
    }

    class InventoryLog {
        +int64 ID
        +string SKUCode
        +int64 WarehouseID
        +string ChangeType
        +int ChangeQuantity
        +int BeforeStock
        +int AfterStock
        +string Remark
    }

    class SKU {
        +int Stock
        +int AvailableStock
        +int LockedStock
        +int SafetyStock
        +LockStock(qty) error
        +UnlockStock(qty) error
        +DeductLockedStock(qty) error
        +RestoreStock(qty) void
        +IsLowStock() bool
    }

    Warehouse "1" --> "*" WarehouseInventory : contains
    WarehouseInventory --> SKU : references
    InventoryLog --> WarehouseInventory : logs
    SKU --> InventoryLog : tracks
```

---

## Stock Management Flow

```mermaid
sequenceDiagram
    participant User
    participant OrderService
    participant InventoryService
    participant SKU
    participant Database

    Note over User,Database: Order Creation Flow (Lock Stock)

    User->>OrderService: Create Order
    OrderService->>InventoryService: LockStock(skuCode, qty)
    InventoryService->>SKU: LockStock(qty)
    SKU->>SKU: AvailableStock -= qty
    SKU->>SKU: LockedStock += qty
    SKU->>Database: Update SKU
    InventoryService->>Database: Create Log
    OrderService-->>User: Order Created

    Note over User,Database: Payment Success Flow (Deduct Locked)

    User->>OrderService: Payment Success
    OrderService->>InventoryService: DeductLockedStock(skuCode, qty)
    InventoryService->>SKU: DeductLockedStock(qty)
    SKU->>SKU: LockedStock -= qty
    SKU->>Database: Update SKU
    InventoryService->>Database: Create Log
    OrderService-->>User: Order Completed

    Note over User,Database: Order Cancellation Flow (Unlock)

    User->>OrderService: Cancel Order
    OrderService->>InventoryService: UnlockStock(skuCode, qty)
    InventoryService->>SKU: UnlockStock(qty)
    SKU->>SKU: LockedStock -= qty
    SKU->>SKU: AvailableStock += qty
    SKU->>Database: Update SKU
    InventoryService->>Database: Create Log
    OrderService-->>User: Order Cancelled
```

---

## Category Hierarchy

```mermaid
graph TB
    Root1[Electronics<br/>Level 1]
    Root2[Clothing<br/>Level 1]
    Root3[Home & Garden<br/>Level 1]

    Root1 --> Cat1[Phones]
    Root1 --> Cat2[Computers]
    Root1 --> Cat3[Accessories]

    Cat1 --> Cat1a[Phone Cases]
    Cat1 --> Cat1b[Screen Protectors]

    Cat2 --> Cat2a[Laptops]
    Cat2 --> Cat2b[Desktops]

    Root2 --> Cat4[Men]
    Root2 --> Cat5[Women]

    Cat4 --> Cat4a[Shirts]
    Cat4 --> Cat4b[Pants]

    Root3 --> Cat6[Furniture]
    Root3 --> Cat7[Decor]

    style Root1 fill:#e1f5fe
    style Root2 fill:#e8f5e9
    style Root3 fill:#fff3e0
    style Cat1a fill:#fce4ec
    style Cat1b fill:#fce4ec
    style Cat2a fill:#f3e5f5
    style Cat2b fill:#f3e5f5
```

---

## Multi-Market Architecture

```mermaid
graph TB
    subgraph "Global Catalog"
        Product[Product]
        SKU[SKU Variants]
    end

    subgraph "Market Configuration"
        CNMarket[China Market<br/>CNY, VAT 13%]
        USMarket[US Market<br/>USD, No VAT]
        UKMarket[UK Market<br/>GBP, VAT 20%, IOSS]
        DEMarket[Germany Market<br/>EUR, VAT 19%, IOSS]
    end

    subgraph "Market-Specific Data"
        CNPrice[¥999]
        USPrice[$149]
        UKPrice[£119]
        DEPrice[€139]

        CNContent[Chinese Content]
        USContent[English Content]
        DEContent[German Content]
    end

    Product --> CNMarket
    Product --> USMarket
    Product --> UKMarket
    Product --> DEMarket

    CNMarket --> CNPrice
    CNMarket --> CNContent
    USMarket --> USPrice
    USMarket --> USContent
    UKMarket --> UKPrice
    DEMarket --> DEPrice
    DEMarket --> DEContent
```

---

## API Request Flow

```mermaid
sequenceDiagram
    participant Client
    participant Gateway
    participant Auth as AuthMiddleware
    participant Handler
    participant Logic
    participant Repo as Repository
    participant DB as Database
    participant Cache as Redis

    Client->>Gateway: HTTP Request + JWT
    Gateway->>Auth: Forward

    alt Invalid Token
        Auth-->>Client: 401 Unauthorized
    end

    Auth->>Auth: Validate JWT
    Auth->>Auth: Extract tenant_id
    Auth->>Auth: Inject tenant context
    Auth->>Handler: Forward with context

    Handler->>Handler: Parse request
    Handler->>Logic: Execute business logic

    Logic->>Cache: Check cache

    alt Cache Hit
        Cache-->>Logic: Return cached data
    else Cache Miss
        Logic->>Repo: Query data
        Repo->>DB: SELECT ... WHERE tenant_id = ?
        DB-->>Repo: Results
        Repo-->>Logic: Domain entities
        Logic->>Cache: Set cache
    end

    Logic-->>Handler: Response DTO
    Handler-->>Client: HTTP Response
```

---

## Deployment Architecture

```mermaid
graph TB
    subgraph "Production Environment"
        subgraph "Load Balancer"
            LB[nginx / ALB]
        end

        subgraph "Kubernetes Cluster"
            subgraph "API Pods"
                Admin1[Admin API Pod 1]
                Admin2[Admin API Pod 2]
                Shop1[Shop API Pod 1]
                Shop2[Shop API Pod 2]
            end

            subgraph "Services"
                AdminSvc[Admin Service<br/>ClusterIP]
                ShopSvc[Shop Service<br/>ClusterIP]
            end
        end

        subgraph "Data Layer"
            MySQLMaster[(MySQL Primary)]
            MySQLReplica[(MySQL Replica)]
            RedisCluster[(Redis Cluster)]
        end

        subgraph "Storage"
            S3Bucket[S3 Bucket<br/>Images, Assets]
        end
    end

    LB --> AdminSvc
    LB --> ShopSvc
    AdminSvc --> Admin1
    AdminSvc --> Admin2
    ShopSvc --> Shop1
    ShopSvc --> Shop2
    Admin1 --> MySQLMaster
    Admin2 --> MySQLReplica
    Admin1 --> RedisCluster
    Admin2 --> RedisCluster
    Shop1 --> MySQLReplica
    Shop2 --> MySQLReplica
    Shop1 --> RedisCluster
    Shop2 --> RedisCluster
    Admin1 --> S3Bucket
    Shop1 --> S3Bucket
```

---

## Error Handling Flow

```mermaid
flowchart TD
    Request[Incoming Request]
    Validate{Validate Request}
    Auth{Auth Check}
    Business{Business Logic}
    DB{Database Operation}
    Success[Success Response]
    Error[Error Response]

    Request --> Validate
    Validate -->|Invalid| Error400[400 Bad Request]
    Validate -->|Valid| Auth
    Auth -->|Failed| Error401[401 Unauthorized]
    Auth -->|Success| Business
    Business -->|Domain Error| Error4xx[4xx Client Error]
    Business -->|Valid| DB
    DB -->|Not Found| Error404[404 Not Found]
    DB -->|Constraint| Error409[409 Conflict]
    DB -->|Timeout| Error503[503 Service Unavailable]
    DB -->|Success| Success

    style Error400 fill:#ffcdd2
    style Error401 fill:#ffcdd2
    style Error4xx fill:#ffcdd2
    style Error404 fill:#ffcdd2
    style Error409 fill:#ffcdd2
    style Error503 fill:#ffcdd2
    style Success fill:#c8e6c9
```

---

## Caching Strategy

```mermaid
flowchart LR
    subgraph "Cache Layers"
        L1[L1: Application Cache<br/>In-Memory]
        L2[L2: Redis Cache<br/>Shared]
        L3[L3: Database<br/>MySQL]
    end

    subgraph "Cache Keys"
        ProductKey["product:{tenant}:{id}"]
        ListKey["products:{tenant}:list:{hash}"]
        CategoryKey["category:tree:{tenant}"]
    end

    subgraph "Cache Strategies"
        ReadThrough[Read-Through]
        WriteThrough[Write-Through]
        CacheAside[Cache-Aside]
    end

    Request[Request] --> L1
    L1 -->|Miss| L2
    L2 -->|Miss| L3

    L2 --> ProductKey
    L2 --> ListKey
    L2 --> CategoryKey

    Write[Write Operation] --> L3
    Write -->|Invalidate| L2
    Write -->|Invalidate| L1
```

---

## Database Schema Overview

```mermaid
erDiagram
    tenants ||--o{ products : owns
    tenants ||--o{ categories : owns
    tenants ||--o{ brands : owns
    tenants ||--o{ markets : owns
    tenants ||--o{ warehouses : owns

    categories ||--o{ products : contains
    categories ||--o{ categories : "parent-child"

    brands ||--o{ products : produces

    products ||--o{ skus : has
    products ||--o{ product_localizations : has
    products ||--o{ product_markets : "sold in"

    skus ||--o{ warehouse_inventories : "stored in"
    skus ||--o{ inventory_logs : tracks

    warehouses ||--o{ warehouse_inventories : contains
    markets ||--o{ product_markets : configures
```

---

## Legend

| Symbol | Meaning |
|--------|---------|
| `()` | Database / External Service |
| `[]` | Application / Service |
| `-->` | Data Flow |
| `-->>` | Response |
| `<>` | Alternative Path |
| `()` in sequence | Participant |