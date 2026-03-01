# MongoDB Distributed System Design Guide
## Database Design, Data Modeling at Scale & Performance Tuning

---

## Table of Contents

1. [Introduction to Distributed MongoDB](#1-introduction-to-distributed-mongodb)
2. [MongoDB Architecture for Scale](#2-mongodb-architecture-for-scale)
3. [Data Modeling Principles](#3-data-modeling-principles)
4. [Schema Design Patterns](#4-schema-design-patterns)
5. [Sharding Strategies](#5-sharding-strategies)
6. [Indexing for Distributed Systems](#6-indexing-for-distributed-systems)
7. [Query Optimization](#7-query-optimization)
8. [Aggregation Pipeline Optimization](#8-aggregation-pipeline-optimization)
9. [Write Performance Tuning](#9-write-performance-tuning)
10. [Read Performance Tuning](#10-read-performance-tuning)
11. [Replication & High Availability](#11-replication--high-availability)
12. [Connection Management](#12-connection-management)
13. [Capacity Planning](#13-capacity-planning)
14. [Monitoring & Observability](#14-monitoring--observability)
15. [Security at Scale](#15-security-at-scale)
16. [Multi-Tenancy Design](#16-multi-tenancy-design)
17. [Time-Series Data Design](#17-time-series-data-design)
18. [Event-Driven Architecture](#18-event-driven-architecture)
19. [Migration & Evolution Strategies](#19-migration--evolution-strategies)
20. [Real-World Case Studies](#20-real-world-case-studies)

---

## 1. Introduction to Distributed MongoDB

### Theory: Understanding Distributed Database Systems

**What is a Distributed Database?**
A distributed database is a collection of logically interconnected databases spread across multiple physical locations, connected via a network. Unlike centralized databases where all data resides on a single server, distributed databases partition and/or replicate data across multiple nodes, providing benefits in scalability, availability, and performance.

**The Need for Distribution**
As applications grow, single-server databases face fundamental limitations:
- **Vertical Scaling Limits**: There's a physical ceiling to how much CPU, RAM, and storage you can add to a single machine
- **Single Point of Failure**: One server crash means complete downtime
- **Geographic Latency**: Users far from the server experience slow response times
- **Compliance Requirements**: Data sovereignty laws may require data to reside in specific regions

**MongoDB's Distributed Architecture Philosophy**
MongoDB was designed from the ground up for horizontal scaling. Its document model, combined with automatic sharding and replication, allows applications to:
1. **Scale horizontally** by adding commodity servers rather than expensive specialized hardware
2. **Achieve high availability** through automatic failover with replica sets
3. **Optimize global access** by placing data geographically close to users
4. **Handle diverse workloads** by scaling reads and writes independently

**The CAP Theorem Context**
The CAP theorem states that a distributed system can only guarantee two of three properties: Consistency, Availability, and Partition Tolerance. MongoDB is classified as a **CP system** (Consistency + Partition Tolerance), prioritizing data consistency during network partitions. However, MongoDB provides tunable consistency through read/write concerns, allowing developers to trade consistency for availability when appropriate for their use case.

### Why Distributed MongoDB?

```
┌─────────────────────────────────────────────────────────────┐
│              DISTRIBUTED MONGODB BENEFITS                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Horizontal  │  │    High     │  │   Data      │         │
│  │  Scaling    │  │ Availability│  │  Locality   │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
│                                                             │
│  • Handle billions of documents                             │
│  • Scale reads and writes independently                     │
│  • Survive node/datacenter failures                         │
│  • Place data close to users globally                       │
│  • Meet compliance requirements (data residency)            │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Concepts

| Concept | Description | Use Case |
|---------|-------------|----------|
| **Replica Set** | Group of mongod instances maintaining same data | High availability, read scaling |
| **Sharding** | Horizontal partitioning across servers | Write scaling, large datasets |
| **Zones** | Associate shards with data ranges | Data locality, compliance |
| **Read Preference** | Which nodes to read from | Latency optimization |
| **Write Concern** | Acknowledgment level for writes | Durability guarantees |

### MongoDB Deployment Topology

```
┌─────────────────────────────────────────────────────────────────┐
│                    SHARDED CLUSTER ARCHITECTURE                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│                        ┌─────────────┐                          │
│                        │   mongos    │                          │
│                        │   Router    │                          │
│                        └──────┬──────┘                          │
│                               │                                  │
│            ┌──────────────────┼──────────────────┐              │
│            │                  │                  │              │
│            ▼                  ▼                  ▼              │
│     ┌────────────┐    ┌────────────┐    ┌────────────┐         │
│     │  Shard 1   │    │  Shard 2   │    │  Shard 3   │         │
│     │ (Replica   │    │ (Replica   │    │ (Replica   │         │
│     │   Set)     │    │   Set)     │    │   Set)     │         │
│     └────────────┘    └────────────┘    └────────────┘         │
│     P   S   S         P   S   S         P   S   S              │
│                                                                  │
│                    ┌─────────────────┐                          │
│                    │  Config Servers │                          │
│                    │  (Replica Set)  │                          │
│                    └─────────────────┘                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

P = Primary, S = Secondary
```

### CAP Theorem Considerations

```javascript
// MongoDB is a CP system (Consistency + Partition Tolerance)
// But can be tuned for different consistency/availability tradeoffs

// Strong consistency (default)
db.orders.insertOne(
    { orderId: "12345", total: 99.99 },
    { writeConcern: { w: "majority" } }
);

// Eventual consistency (for higher availability)
db.orders.find({ status: "pending" })
    .readPreference("secondaryPreferred")
    .readConcern("local");
```

---

## 2. MongoDB Architecture for Scale

### Theory: Architectural Foundations

**Shared-Nothing Architecture**
MongoDB employs a "shared-nothing" architecture where each shard operates independently with its own CPU, memory, and storage. This contrasts with "shared-disk" architectures where nodes share storage. The shared-nothing approach provides:
- **Linear scalability**: Performance scales proportionally with added nodes
- **Fault isolation**: Failures in one shard don't directly impact others
- **Independent optimization**: Each shard can be tuned for its specific workload

**The Role of Routing and Coordination**
In MongoDB's architecture, three distinct components work together:

1. **Query Routers (mongos)**: Stateless processes that route client requests to appropriate shards. They maintain no persistent state, allowing horizontal scaling of the routing layer itself.

2. **Config Servers**: Store cluster metadata including shard locations, chunk ranges, and authentication data. They form a replica set for high availability.

3. **Shards**: Store actual data. Each shard is a replica set, providing both horizontal scaling and high availability.

**Data Distribution Model**
MongoDB divides collections into **chunks** (contiguous ranges of shard key values). The balancer process automatically migrates chunks between shards to maintain even distribution. This automatic balancing is crucial because:
- Uneven distribution creates "hot spots" that limit throughput
- Manual balancing doesn't scale with growing data and changing access patterns
- The system can adapt to hardware changes (adding/removing shards)

**Consistency Model**
MongoDB uses a **primary-based replication** model within each shard. All writes go to the primary, which replicates to secondaries via the **oplog** (operation log). This provides:
- Strong consistency for reads from primary
- Tunable consistency for secondary reads
- Automatic leader election during failovers

### Shard Architecture Deep Dive

```
┌─────────────────────────────────────────────────────────────────┐
│                      SHARD INTERNALS                            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Each Shard is a Replica Set:                                   │
│  ┌───────────────────────────────────────────────────────┐     │
│  │                    REPLICA SET                         │     │
│  │                                                        │     │
│  │  ┌─────────┐    ┌─────────┐    ┌─────────┐           │     │
│  │  │ PRIMARY │    │SECONDARY│    │SECONDARY│           │     │
│  │  │         │◄──►│         │◄──►│         │           │     │
│  │  │ Reads   │    │  Reads  │    │  Reads  │           │     │
│  │  │ Writes  │    │  +Oplog │    │  +Oplog │           │     │
│  │  └─────────┘    └─────────┘    └─────────┘           │     │
│  │                                                        │     │
│  │  Optional:                                             │     │
│  │  ┌─────────┐    ┌─────────┐                          │     │
│  │  │ ARBITER │    │ HIDDEN  │ (backup, analytics)      │     │
│  │  └─────────┘    └─────────┘                          │     │
│  └───────────────────────────────────────────────────────┘     │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Config Server Architecture

```javascript
// Config servers store cluster metadata:
// - Shard locations
// - Chunk ranges
// - Database/collection configs
// - Authentication data

// Check config server status
sh.status()

// Config database collections
use config
show collections
// chunks, databases, mongos, shards, tags, version, etc.

// View chunk distribution
db.chunks.find({ ns: "mydb.orders" }).pretty()
```

### Mongos Router Behavior

```javascript
// Mongos routes queries to appropriate shards

// Targeted query (uses shard key) - hits single shard
db.orders.find({ customerId: "cust123" })  // customerId is shard key

// Scatter-gather query - hits ALL shards
db.orders.find({ status: "pending" })  // status is not shard key

// Broadcast operations
db.orders.createIndex({ status: 1 })  // Creates on all shards
```

### WiredTiger Storage Engine

```javascript
// WiredTiger configuration for distributed systems

// mongod.conf
storage:
  engine: wiredTiger
  wiredTiger:
    engineConfig:
      cacheSizeGB: 8  // 50% of RAM - 1GB, max(256MB, 50% RAM - 1GB)
      journalCompressor: snappy
      directoryForIndexes: true  // Separate directories for indexes
    collectionConfig:
      blockCompressor: snappy  // or zstd for better compression
    indexConfig:
      prefixCompression: true

// Monitor cache usage
db.serverStatus().wiredTiger.cache
```

---

## 3. Data Modeling Principles

### Theory: Document-Oriented Data Modeling

**The Paradigm Shift from Relational Modeling**
Relational databases normalize data to eliminate redundancy, following normal forms (1NF, 2NF, 3NF, BCNF). Document databases like MongoDB embrace a different philosophy:

- **Denormalization is acceptable**: Duplicating data is often the right choice when it improves read performance
- **Schema follows access patterns**: Design documents around how data is queried, not how it's logically organized
- **Atomicity at document level**: MongoDB guarantees atomic operations on a single document, so embedding related data provides transactional guarantees without distributed transactions

**The Fundamental Trade-off: Embedding vs Referencing**

| Aspect | Embedding | Referencing |
|--------|-----------|-------------|
| **Read Performance** | Single query retrieves all data | Requires multiple queries or $lookup |
| **Write Performance** | Updating embedded data may require full document rewrite | Updates are targeted to specific documents |
| **Data Consistency** | Denormalized data can become inconsistent | Normalized data has single source of truth |
| **Document Size** | Risk of exceeding 16MB limit | Documents stay small |
| **Atomicity** | Automatic for embedded changes | Requires multi-document transactions |

**Cardinality Analysis**
Understanding relationship cardinality is crucial for modeling decisions:

- **One-to-One**: Almost always embed unless the child document is very large or accessed separately
- **One-to-Few** (<100): Embed. The overhead of referencing outweighs any benefits
- **One-to-Many** (100-1000): Consider embedding with careful size monitoring, or reference with caching
- **One-to-Squillions** (unlimited): Always reference. Embedding creates unbounded document growth

**The Working Set Concept**
MongoDB performs best when the "working set" (frequently accessed data + indexes) fits in RAM. Document design impacts working set size:
- Smaller documents = more documents fit in memory
- Embedding everything creates larger documents, potentially evicting other data from cache
- Strategic embedding keeps related data together, improving cache locality

### Embedding vs Referencing Decision Matrix

```
┌─────────────────────────────────────────────────────────────────┐
│                 EMBEDDING VS REFERENCING                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  EMBED WHEN:                      REFERENCE WHEN:               │
│  ├── Data accessed together       ├── Data accessed separately  │
│  ├── One-to-few relationship      ├── One-to-many (unbounded)   │
│  ├── Data rarely changes          ├── Data frequently updated   │
│  ├── Document size < 16MB         ├── Large sub-documents       │
│  ├── Atomic operations needed     ├── Many-to-many relationship │
│  └── No duplication concerns      └── Avoid duplication         │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Embedding Pattern

```javascript
// GOOD: Embed when data is accessed together
// Order with embedded items - single read retrieves everything
{
    _id: ObjectId("..."),
    orderId: "ORD-2024-001",
    customerId: "CUST-123",
    orderDate: ISODate("2024-03-15"),
    status: "shipped",
    // Embedded - always accessed with order
    items: [
        { 
            productId: "PROD-001",
            name: "Laptop",
            quantity: 1,
            price: 999.99
        },
        {
            productId: "PROD-002", 
            name: "Mouse",
            quantity: 2,
            price: 29.99
        }
    ],
    // Embedded - rarely changes
    shippingAddress: {
        street: "123 Main St",
        city: "New York",
        zip: "10001",
        country: "USA"
    },
    totals: {
        subtotal: 1059.97,
        tax: 84.80,
        shipping: 10.00,
        total: 1154.77
    }
}
```

### Referencing Pattern

```javascript
// Users collection
{
    _id: ObjectId("user123"),
    email: "john@example.com",
    name: "John Doe",
    createdAt: ISODate("2024-01-01")
}

// Orders collection - reference user
{
    _id: ObjectId("order456"),
    userId: ObjectId("user123"),  // Reference
    orderDate: ISODate("2024-03-15"),
    items: [...],
    total: 1154.77
}

// Reviews collection - reference both
{
    _id: ObjectId("review789"),
    userId: ObjectId("user123"),
    productId: ObjectId("product001"),
    rating: 5,
    comment: "Great product!",
    createdAt: ISODate("2024-03-20")
}

// Query with $lookup (join)
db.orders.aggregate([
    { $match: { _id: ObjectId("order456") } },
    {
        $lookup: {
            from: "users",
            localField: "userId",
            foreignField: "_id",
            as: "user"
        }
    },
    { $unwind: "$user" }
]);
```

### Hybrid Pattern (Best of Both Worlds)

```javascript
// Store frequently accessed data embedded, reference for full details
{
    _id: ObjectId("order123"),
    // Embedded summary for display
    customer: {
        _id: ObjectId("user123"),
        name: "John Doe",        // Denormalized for quick access
        email: "john@example.com"
    },
    items: [
        {
            productId: ObjectId("prod001"),
            name: "Laptop",      // Denormalized
            sku: "LPT-001",      // Denormalized
            quantity: 1,
            unitPrice: 999.99,
            // Reference for full product details
            _productRef: ObjectId("prod001")
        }
    ],
    totals: {
        subtotal: 999.99,
        tax: 80.00,
        total: 1079.99
    }
}
```

### Document Size Considerations

```javascript
// MongoDB document limit: 16MB
// Best practice: Keep documents under 1MB for performance

// BAD: Unbounded array growth
{
    userId: "user123",
    // This array will grow forever!
    activityLog: [
        { action: "login", timestamp: ISODate("...") },
        // ... thousands more entries
    ]
}

// GOOD: Bucket pattern for time-series data
{
    userId: "user123",
    date: ISODate("2024-03-15"),
    activities: [
        { action: "login", timestamp: ISODate("..."), minute: 0 },
        { action: "view_product", timestamp: ISODate("..."), minute: 5 },
        // Max ~1000 entries per bucket (one day)
    ],
    count: 2
}
```

---

## 4. Schema Design Patterns

### Theory: Design Patterns for Document Databases

**Why Design Patterns Matter**
Just as software engineering has established patterns (Factory, Singleton, Observer), document database design has evolved patterns that solve recurring problems. These patterns represent distilled wisdom from thousands of production deployments.

**The Pattern Categories**

1. **Representation Patterns**: Address how to structure data within documents
   - Attribute Pattern: Handles polymorphic attributes
   - Schema Versioning: Manages evolving schemas

2. **Grouping Patterns**: Optimize for data that's accessed together
   - Bucket Pattern: Groups time-series data
   - Subset Pattern: Keeps frequently accessed subset embedded

3. **Optimization Patterns**: Improve performance characteristics
   - Computed Pattern: Pre-calculates expensive aggregations
   - Extended Reference: Denormalizes frequently accessed fields

4. **Edge Case Patterns**: Handle exceptional scenarios
   - Outlier Pattern: Manages documents with extreme characteristics

**Choosing Patterns Based on Query Patterns**
The key to selecting appropriate patterns is understanding your access patterns:

- **Read-heavy workloads**: Favor denormalization (Extended Reference, Computed, Subset)
- **Write-heavy workloads**: Favor normalization to reduce document update sizes
- **Variable schemas**: Use Attribute Pattern for flexibility without sparse indexes
- **Time-series data**: Bucket Pattern reduces document count and enables pre-aggregation

**Pattern Composition**
Patterns aren't mutually exclusive. Production systems often combine multiple patterns:
- E-commerce products: Attribute Pattern + Subset Pattern (for reviews) + Computed Pattern (for ratings)
- Social networks: Extended Reference + Outlier Pattern (for celebrity accounts)
- IoT systems: Bucket Pattern + Computed Pattern (for aggregations)

### Attribute Pattern

```javascript
// Problem: Products have varying attributes
// BAD: Sparse fields
{
    _id: "laptop1",
    name: "Gaming Laptop",
    brand: "Dell",
    cpu: "Intel i9",           // Only for laptops
    ram: "32GB",               // Only for laptops
    screenSize: "15.6 inch",   // Only for laptops/monitors
    resolution: "4K",          // Only for some electronics
    color: "Black",            // Most products
    weight: "2.5kg",           // Most products
    batteryLife: "8 hours"     // Only for portable devices
}

// GOOD: Attribute pattern
{
    _id: "laptop1",
    name: "Gaming Laptop",
    brand: "Dell",
    category: "Electronics/Computers/Laptops",
    attributes: [
        { k: "cpu", v: "Intel i9", u: "processor" },
        { k: "ram", v: "32", u: "GB" },
        { k: "screenSize", v: "15.6", u: "inch" },
        { k: "resolution", v: "4K", u: null },
        { k: "color", v: "Black", u: null },
        { k: "weight", v: "2.5", u: "kg" },
        { k: "batteryLife", v: "8", u: "hours" }
    ]
}

// Index for querying attributes
db.products.createIndex({ "attributes.k": 1, "attributes.v": 1 })

// Query products by attribute
db.products.find({ 
    attributes: { 
        $elemMatch: { k: "ram", v: { $gte: "16" } } 
    }
})
```

### Bucket Pattern (Time-Series)

```javascript
// Problem: High-volume time-series data (IoT, logs, metrics)
// BAD: One document per measurement
{
    sensorId: "sensor001",
    timestamp: ISODate("2024-03-15T10:00:00Z"),
    temperature: 23.5
}
// Creates millions of tiny documents!

// GOOD: Bucket pattern
{
    sensorId: "sensor001",
    bucketStart: ISODate("2024-03-15T10:00:00Z"),
    bucketEnd: ISODate("2024-03-15T11:00:00Z"),
    measurements: [
        { t: ISODate("2024-03-15T10:00:00Z"), temp: 23.5, humidity: 45 },
        { t: ISODate("2024-03-15T10:01:00Z"), temp: 23.6, humidity: 44 },
        // ... up to 60 measurements per hour
    ],
    count: 60,
    sum: { temp: 1416.0, humidity: 2700 },  // Pre-aggregated
    avg: { temp: 23.6, humidity: 45 },
    min: { temp: 23.1, humidity: 42 },
    max: { temp: 24.2, humidity: 48 }
}

// Benefits:
// - Fewer documents (60x reduction)
// - Pre-computed aggregates
// - Better index efficiency
// - Reduced storage

// Insert with bucket management
db.sensorData.updateOne(
    {
        sensorId: "sensor001",
        count: { $lt: 60 },
        bucketStart: {
            $gte: ISODate("2024-03-15T10:00:00Z"),
            $lt: ISODate("2024-03-15T11:00:00Z")
        }
    },
    {
        $push: {
            measurements: {
                t: ISODate("2024-03-15T10:30:00Z"),
                temp: 23.8,
                humidity: 46
            }
        },
        $inc: { 
            count: 1,
            "sum.temp": 23.8,
            "sum.humidity": 46
        },
        $min: { "min.temp": 23.8, "min.humidity": 46 },
        $max: { "max.temp": 23.8, "max.humidity": 46 },
        $setOnInsert: {
            sensorId: "sensor001",
            bucketStart: ISODate("2024-03-15T10:00:00Z"),
            bucketEnd: ISODate("2024-03-15T11:00:00Z")
        }
    },
    { upsert: true }
);
```

### Computed Pattern

```javascript
// Problem: Expensive computations on read
// BAD: Calculate totals on every read
db.orders.aggregate([
    { $match: { customerId: "cust123" } },
    { $group: {
        _id: "$customerId",
        totalOrders: { $sum: 1 },
        totalSpent: { $sum: "$total" },
        avgOrderValue: { $avg: "$total" }
    }}
]);

// GOOD: Pre-compute and store
{
    _id: "cust123",
    email: "customer@example.com",
    name: "John Doe",
    // Computed fields - updated on order changes
    stats: {
        totalOrders: 47,
        totalSpent: 5234.56,
        avgOrderValue: 111.37,
        lastOrderDate: ISODate("2024-03-15"),
        firstOrderDate: ISODate("2023-01-10")
    },
    tier: "gold",  // Computed from totalSpent
    computedAt: ISODate("2024-03-15T12:00:00Z")
}

// Update computed fields when order is placed
db.customers.updateOne(
    { _id: "cust123" },
    {
        $inc: {
            "stats.totalOrders": 1,
            "stats.totalSpent": 89.99
        },
        $set: {
            "stats.lastOrderDate": new Date(),
            "stats.avgOrderValue": { 
                $divide: [
                    { $add: ["$stats.totalSpent", 89.99] },
                    { $add: ["$stats.totalOrders", 1] }
                ]
            }
        }
    }
);
```

### Extended Reference Pattern

```javascript
// Problem: Frequently need subset of related document data
// BAD: Always $lookup
db.orders.aggregate([
    { $match: { orderId: "ORD-001" } },
    { $lookup: { from: "products", ... } },  // Expensive join
    { $lookup: { from: "customers", ... } }  // Another expensive join
]);

// GOOD: Extended reference - embed frequently accessed fields
{
    _id: ObjectId("order123"),
    orderId: "ORD-2024-001",
    
    // Extended reference - not just ID
    customer: {
        _id: ObjectId("cust123"),
        name: "John Doe",           // Display name
        email: "john@example.com",  // For notifications
        tier: "gold"                // For discounts
        // Full customer doc has much more data
    },
    
    items: [
        {
            product: {
                _id: ObjectId("prod001"),
                name: "Laptop",       // For display
                sku: "LPT-DELL-001", // For inventory
                imageUrl: "/images/laptop.jpg"  // For UI
                // Full product doc has specs, reviews, etc.
            },
            quantity: 1,
            unitPrice: 999.99
        }
    ]
}

// Update strategy: Background job or triggers
// When product name changes:
db.orders.updateMany(
    { "items.product._id": ObjectId("prod001") },
    { $set: { "items.$[elem].product.name": "New Laptop Name" } },
    { arrayFilters: [{ "elem.product._id": ObjectId("prod001") }] }
);
```

### Outlier Pattern

```javascript
// Problem: Most documents are small, few are huge (power law distribution)
// Example: Most users have 10 followers, celebrities have millions

// GOOD: Detect and handle outliers differently
{
    _id: "user123",
    username: "regular_user",
    followerCount: 150,
    followers: [                    // Embedded for normal users
        "user456", "user789", ...
    ],
    hasOverflow: false
}

{
    _id: "celebrity001",
    username: "famous_person",
    followerCount: 5000000,
    followers: [],                  // Empty - stored separately
    hasOverflow: true               // Flag for overflow
}

// Overflow collection for outliers
{
    _id: ObjectId("..."),
    userId: "celebrity001",
    followers: ["user001", "user002", ...],  // Batch of 10000
    batchNumber: 1
}

// Query handling
async function getFollowers(userId) {
    const user = await db.users.findOne({ _id: userId });
    
    if (!user.hasOverflow) {
        return user.followers;
    }
    
    // For outliers, query overflow collection
    return db.userFollowers
        .find({ userId: userId })
        .toArray()
        .flatMap(doc => doc.followers);
}
```

### Schema Versioning Pattern

```javascript
// Problem: Schema evolves over time
// GOOD: Version field for migration handling
{
    _id: ObjectId("..."),
    schemaVersion: 2,
    email: "user@example.com",
    // Version 2: Split name into first/last
    firstName: "John",
    lastName: "Doe"
}

// Legacy document (version 1)
{
    _id: ObjectId("..."),
    schemaVersion: 1,
    email: "user@example.com",
    name: "John Doe"  // Version 1: Combined name
}

// Application-level migration
function normalizeUser(doc) {
    if (doc.schemaVersion === 1) {
        const [firstName, ...lastParts] = doc.name.split(' ');
        return {
            ...doc,
            firstName,
            lastName: lastParts.join(' '),
            schemaVersion: 2
        };
    }
    return doc;
}

// Lazy migration on read
async function getUser(userId) {
    let user = await db.users.findOne({ _id: userId });
    
    if (user.schemaVersion < CURRENT_VERSION) {
        user = normalizeUser(user);
        await db.users.updateOne(
            { _id: userId },
            { $set: user }
        );
    }
    
    return user;
}

// Batch migration script
db.users.find({ schemaVersion: 1 }).forEach(doc => {
    const [firstName, ...lastParts] = doc.name.split(' ');
    db.users.updateOne(
        { _id: doc._id },
        {
            $set: {
                firstName,
                lastName: lastParts.join(' '),
                schemaVersion: 2
            },
            $unset: { name: "" }
        }
    );
});
```

### Subset Pattern

```javascript
// Problem: Documents contain large arrays, but often only need recent items
// Example: Product reviews - thousands exist, display only top 10

// GOOD: Embed subset, reference full list
{
    _id: "product001",
    name: "Amazing Product",
    price: 99.99,
    
    // Embedded subset - most recent/relevant
    topReviews: [
        {
            _id: "review001",
            userId: "user123",
            userName: "John D.",
            rating: 5,
            title: "Best purchase ever!",
            createdAt: ISODate("2024-03-15")
        },
        // ... top 10 reviews only
    ],
    
    reviewStats: {
        totalCount: 1547,
        averageRating: 4.6,
        ratingDistribution: {
            5: 890, 4: 412, 3: 156, 2: 54, 1: 35
        }
    }
}

// Full reviews in separate collection
{
    _id: "review001",
    productId: "product001",
    userId: "user123",
    userName: "John D.",
    rating: 5,
    title: "Best purchase ever!",
    body: "Full review text...",
    helpful: 234,
    images: [...],
    createdAt: ISODate("2024-03-15")
}

// Update subset when new review is added
db.products.updateOne(
    { _id: "product001" },
    [
        {
            $set: {
                topReviews: {
                    $slice: [
                        { $concatArrays: [
                            [newReview],
                            "$topReviews"
                        ]},
                        10
                    ]
                }
            }
        }
    ]
);
```

---

## 5. Sharding Strategies

### Theory: Fundamentals of Data Sharding

**What is Sharding?**
Sharding (also called horizontal partitioning) divides a dataset across multiple servers. Each server (shard) holds a subset of the total data. Unlike replication (which copies all data to multiple servers), sharding distributes unique chunks of data.

**Why Sharding is Necessary**
1. **Dataset exceeds single server capacity**: When data grows beyond what one server can store
2. **Throughput exceeds single server capability**: When read/write operations exceed what one server can handle
3. **Working set exceeds available RAM**: When active data can't fit in memory, causing disk I/O bottlenecks

**The Shard Key: Your Most Important Decision**
The shard key is the field(s) MongoDB uses to distribute documents across shards. This decision is **permanent** (without resharding the entire collection) and affects:

- **Write distribution**: Whether writes spread evenly or create hot spots
- **Query efficiency**: Whether queries hit one shard (targeted) or all shards (scatter-gather)
- **Data locality**: Whether related data resides together

**Sharding Strategies Explained**

| Strategy | How it Works | Best For | Avoid When |
|----------|--------------|----------|------------|
| **Hashed** | Hash function distributes documents randomly | Even write distribution | Range queries needed |
| **Ranged** | Documents with similar keys on same shard | Range queries, data locality | Monotonic keys (timestamps) |
| **Zoned** | Specific data ranges pinned to specific shards | Geographic locality, compliance | Uniform access patterns |

**The Chunk Concept**
MongoDB organizes sharded data into chunks (default 128MB). Understanding chunks is crucial:

- **Chunk splitting**: When a chunk exceeds max size, MongoDB splits it
- **Chunk migration**: The balancer moves chunks between shards for even distribution
- **Jumbo chunks**: Chunks that can't be split (all documents have same shard key) - a sign of poor shard key choice

**Cardinality and Frequency**
- **Cardinality**: Number of unique shard key values. Low cardinality limits distribution (only N shards can have data)
- **Frequency**: How often each value appears. High frequency values create "jumbo chunks" that can't be split

Ideal shard keys have **high cardinality** AND **low frequency** - many unique values, each appearing in few documents.

### Choosing a Shard Key

```
┌─────────────────────────────────────────────────────────────────┐
│                    SHARD KEY SELECTION                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  IDEAL SHARD KEY PROPERTIES:                                    │
│                                                                  │
│  1. HIGH CARDINALITY                                            │
│     └── Many unique values for even distribution                │
│                                                                  │
│  2. LOW FREQUENCY                                               │
│     └── No single value dominates (avoid hot spots)             │
│                                                                  │
│  3. NON-MONOTONIC                                               │
│     └── Values don't always increase (avoid last-shard writes)  │
│                                                                  │
│  4. QUERY ISOLATION                                             │
│     └── Most queries include shard key (avoid scatter-gather)   │
│                                                                  │
│  5. WRITE DISTRIBUTION                                          │
│     └── Writes spread evenly across shards                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Hashed Sharding

```javascript
// Best for: Even write distribution, no range queries needed
// Creates random distribution based on hash of field value

// Enable sharding on database
sh.enableSharding("mydb")

// Create hashed shard key
sh.shardCollection("mydb.users", { _id: "hashed" })

// Or on custom field
sh.shardCollection("mydb.orders", { orderId: "hashed" })

// Compound hashed key (MongoDB 4.4+)
sh.shardCollection(
    "mydb.events",
    { userId: 1, timestamp: "hashed" }
)

// Pros:
// + Even data distribution
// + Good write scalability
// + No hot spots

// Cons:
// - No efficient range queries
// - Scatter-gather for non-shard-key queries
// - Cannot do targeted queries without exact match
```

### Ranged Sharding

```javascript
// Best for: Range queries, data locality
// Documents with similar keys stored together

sh.shardCollection("mydb.logs", { timestamp: 1 })

// Compound ranged key (recommended)
sh.shardCollection(
    "mydb.orders",
    { customerId: 1, orderDate: 1 }
)

// Pros:
// + Efficient range queries
// + Targeted queries possible
// + Good for time-series with compound key

// Cons:
// - Can create hot spots
// - Monotonic keys cause all writes to one shard
// - Uneven distribution possible
```

### Zone Sharding (Data Locality)

```javascript
// Use case: Geographic data locality, compliance requirements

// Add shard to zone
sh.addShardTag("shard0001", "US")
sh.addShardTag("shard0002", "EU")
sh.addShardTag("shard0003", "APAC")

// Define zone ranges
sh.updateZoneKeyRange(
    "mydb.users",
    { region: "US", _id: MinKey },
    { region: "US", _id: MaxKey },
    "US"
)

sh.updateZoneKeyRange(
    "mydb.users",
    { region: "EU", _id: MinKey },
    { region: "EU", _id: MaxKey },
    "EU"
)

// Shard collection with zone-aware key
sh.shardCollection("mydb.users", { region: 1, _id: 1 })

// Document automatically routes to correct zone
db.users.insertOne({
    _id: ObjectId(),
    region: "EU",
    email: "user@example.de",
    name: "Hans Schmidt"
})
// → Routes to EU shard
```

### Shard Key Anti-Patterns

```javascript
// ❌ BAD: Monotonically increasing key (all writes to last shard)
sh.shardCollection("mydb.logs", { timestamp: 1 })
sh.shardCollection("mydb.orders", { _id: 1 })  // ObjectId is monotonic

// ❌ BAD: Low cardinality (limited distribution)
sh.shardCollection("mydb.orders", { status: 1 })
// Only ~5 unique values (pending, processing, shipped, etc.)

// ❌ BAD: Frequently changing field
sh.shardCollection("mydb.products", { price: 1 })
// Price changes require document migration between shards!

// ✅ GOOD: High cardinality, query-aligned compound key
sh.shardCollection("mydb.orders", { customerId: 1, orderDate: 1 })

// ✅ GOOD: Hashed for write-heavy, no range queries
sh.shardCollection("mydb.events", { eventId: "hashed" })

// ✅ GOOD: Compound with hashed for distribution + targeting
sh.shardCollection("mydb.logs", { tenantId: 1, timestamp: "hashed" })
```

### Pre-splitting Chunks

```javascript
// Pre-split to avoid balancing during initial load
// Hashed shard key: Use numInitialChunks
sh.shardCollection(
    "mydb.events",
    { eventId: "hashed" },
    false,  // unique
    { numInitialChunks: 100 }  // Pre-create 100 chunks
)

// Ranged shard key: Manual split
sh.splitAt("mydb.orders", { customerId: "A" })
sh.splitAt("mydb.orders", { customerId: "M" })
sh.splitAt("mydb.orders", { customerId: "Z" })

// Move chunks to specific shards
sh.moveChunk(
    "mydb.orders",
    { customerId: "A" },
    "shard0001"
)
```

---

## 6. Indexing for Distributed Systems

### Theory: The Science of Database Indexing

**What is an Index?**
An index is a separate data structure that maintains sorted references to documents, enabling rapid lookups without scanning every document. Think of it like a book's index - instead of reading every page to find "MongoDB", you look up "MongoDB" in the index and go directly to the relevant pages.

**How B-Tree Indexes Work**
MongoDB's default index type uses **B-Tree** (Balanced Tree) data structures:

1. **Tree Structure**: Data organized in a hierarchy where each node contains multiple keys and pointers
2. **Balanced Property**: All leaf nodes are at the same depth, ensuring consistent O(log n) lookup time
3. **Sorted Order**: Keys are maintained in sorted order, enabling efficient range queries
4. **Block-Oriented**: Node sizes match disk/memory block sizes for efficient I/O

**Index Selectivity and Cardinality**
- **Selectivity**: The percentage of documents that match a query condition. High selectivity (few matches) benefits most from indexing
- **Cardinality**: Number of unique values in a field. High cardinality fields make better index candidates
- **Rule of Thumb**: Indexes that eliminate 90%+ of documents are highly effective; indexes eliminating less than 50% may not be beneficial

**The Cost of Indexes**
Indexes aren't free - they have trade-offs:

| Benefit | Cost |
|---------|------|
| Faster reads | Slower writes (index updates) |
| Efficient sorting | Storage space |
| Query optimization | Memory consumption |

**Index Selection in Distributed Systems**
In sharded clusters, indexes exist independently on each shard. Important considerations:
- **Local indexes**: Standard indexes on shard-local data
- **Global indexes**: Not supported - queries without shard key must scatter to all shards
- **Index consistency**: Indexes must be identical across all shards

**The Working Set and Index Memory**
For optimal performance, frequently-used indexes should fit in RAM. When indexes exceed available memory:
- Index traversal causes disk I/O
- Query latency increases dramatically
- Consider partial indexes to reduce size

### Index Types and Use Cases

```javascript
// 1. Single Field Index
db.users.createIndex({ email: 1 })
db.orders.createIndex({ createdAt: -1 })  // Descending for recent-first

// 2. Compound Index (most common for complex queries)
db.orders.createIndex({ customerId: 1, orderDate: -1, status: 1 })

// 3. Multikey Index (arrays)
db.products.createIndex({ tags: 1 })
db.products.createIndex({ "variants.sku": 1 })

// 4. Text Index (full-text search)
db.products.createIndex(
    { name: "text", description: "text" },
    { weights: { name: 10, description: 5 } }
)

// 5. Hashed Index (equality only, even distribution)
db.users.createIndex({ sessionToken: "hashed" })

// 6. Geospatial Index
db.locations.createIndex({ coordinates: "2dsphere" })

// 7. Wildcard Index (dynamic schemas)
db.products.createIndex({ "attributes.$**": 1 })

// 8. Partial Index (conditional)
db.orders.createIndex(
    { createdAt: 1 },
    { partialFilterExpression: { status: "pending" } }
)

// 9. TTL Index (auto-expiration)
db.sessions.createIndex(
    { lastAccess: 1 },
    { expireAfterSeconds: 3600 }
)

// 10. Unique Index
db.users.createIndex(
    { email: 1 },
    { unique: true }
)

// 11. Sparse Index
db.users.createIndex(
    { phoneNumber: 1 },
    { sparse: true }  // Only index documents with this field
)
```

### Compound Index Optimization (ESR Rule)

```javascript
// ESR: Equality, Sort, Range
// Order compound index fields following ESR rule

// Query pattern:
db.orders.find({
    customerId: "cust123",    // Equality
    status: "shipped",        // Equality
    orderDate: { $gte: ISODate("2024-01-01") }  // Range
}).sort({ total: -1 })        // Sort

// Optimal index:
db.orders.createIndex({
    customerId: 1,    // E: Equality fields first
    status: 1,        // E: More equality fields
    total: -1,        // S: Sort field
    orderDate: 1      // R: Range field last
})

// Index prefixes are usable:
// - { customerId: 1 }
// - { customerId: 1, status: 1 }
// - { customerId: 1, status: 1, total: -1 }
// - { customerId: 1, status: 1, total: -1, orderDate: 1 }
```

### Covered Queries

```javascript
// Covered query: Answered entirely from index, no document fetch
// Include all queried and projected fields in index

// Index
db.orders.createIndex({
    customerId: 1,
    orderDate: 1,
    total: 1,
    status: 1
})

// Covered query (only index fields, no _id)
db.orders.find(
    { customerId: "cust123", orderDate: { $gte: ISODate("2024-01-01") } },
    { customerId: 1, orderDate: 1, total: 1, status: 1, _id: 0 }
)

// Check with explain
db.orders.explain("executionStats").find(...)
// Look for: totalDocsExamined: 0
```

### Index Intersection

```javascript
// MongoDB can combine multiple indexes for a query
// But usually less efficient than compound index

// Two separate indexes
db.orders.createIndex({ customerId: 1 })
db.orders.createIndex({ status: 1 })

// Query might use index intersection
db.orders.find({ customerId: "cust123", status: "pending" })

// Better: Single compound index
db.orders.createIndex({ customerId: 1, status: 1 })
```

### Index Management for Distributed Systems

```javascript
// Create index in background (deprecated in 4.2+, now default)
db.orders.createIndex(
    { customerId: 1 },
    { background: true }  // Deprecated, indexes build in background by default
)

// Hide index (test before dropping)
db.orders.hideIndex("customerId_1")
// Test application
db.orders.unhideIndex("customerId_1")

// Rolling index build on replica set
// Indexes build on secondaries first, then primary

// Check index usage
db.orders.aggregate([
    { $indexStats: {} }
])

// Drop unused indexes
db.orders.getIndexes()
db.orders.dropIndex("unused_index_name")

// Index size
db.orders.stats().indexSizes
```

---

## 7. Query Optimization

### Theory: Query Processing and Optimization

**The Query Execution Pipeline**
When MongoDB receives a query, it goes through several stages:

1. **Parsing**: Query syntax is validated and converted to internal representation
2. **Planning**: Query planner generates candidate execution plans
3. **Optimization**: Cost-based optimizer selects the best plan
4. **Execution**: Selected plan is executed and results returned

**The Query Planner's Decision Process**
MongoDB's query planner evaluates multiple execution strategies:

- **Index Selection**: Which index (if any) best serves this query?
- **Index Intersection**: Can multiple indexes be combined?
- **Sort Strategy**: Can sorting use an index or require in-memory sort?
- **Projection Pushdown**: Can projection be pushed to index scan level?

**Understanding Query Shapes**
A "query shape" is the combination of:
- Query filter structure (which fields, which operators)
- Sort specification
- Projection fields

MongoDB caches query plans by shape. Once a plan is chosen for a shape, it's reused until:
- Index changes occur
- Collection statistics change significantly
- Server restarts

**The Concept of Selectivity**
Query optimization relies heavily on **selectivity estimation**:
- How many documents will this condition match?
- MongoDB samples data to estimate selectivity
- Poor estimates lead to suboptimal plan selection

**Scatter-Gather vs Targeted Queries**
In sharded clusters:
- **Targeted Query**: Contains shard key, routes to single shard - O(1) shards
- **Scatter-Gather**: No shard key, queries all shards, merges results - O(n) shards

Scatter-gather queries don't scale - doubling shards doubles query cost. Design queries to be targeted whenever possible.

**Query Complexity and Resources**
Query execution consumes resources proportionally:
- **Documents examined**: CPU time, disk I/O
- **Index keys examined**: Memory access, potential disk I/O
- **Documents returned**: Network bandwidth
- **In-memory sorting**: RAM usage (100MB limit without allowDiskUse)

### Query Execution Analysis

```javascript
// Basic explain
db.orders.find({ customerId: "cust123" }).explain()

// Execution stats
db.orders.find({ customerId: "cust123" }).explain("executionStats")

// All plans considered
db.orders.find({ customerId: "cust123" }).explain("allPlansExecution")

// Key metrics to analyze:
{
    "executionStats": {
        "executionSuccess": true,
        "nReturned": 100,           // Documents returned
        "executionTimeMillis": 15,  // Total time
        "totalKeysExamined": 100,   // Index entries scanned
        "totalDocsExamined": 100,   // Documents scanned
        // Goal: nReturned ≈ totalKeysExamined ≈ totalDocsExamined
    },
    "queryPlanner": {
        "winningPlan": {
            "stage": "FETCH",        // FETCH = needed doc lookup
            "inputStage": {
                "stage": "IXSCAN",   // Index Scan (good!)
                "indexName": "customerId_1"
            }
        }
    }
}

// Bad stages to watch for:
// COLLSCAN - Full collection scan (no index used)
// SORT - In-memory sort (index could cover)
// SORT_KEY_GENERATOR - Building sort keys
```

### Query Patterns and Optimization

```javascript
// ✅ GOOD: Targeted query with index
db.orders.find({ customerId: "cust123" })  // Uses customerId index

// ❌ BAD: Regex with leading wildcard
db.products.find({ name: /.*laptop.*/i })  // Full scan

// ✅ GOOD: Text search instead
db.products.find({ $text: { $search: "laptop" } })

// ❌ BAD: $or without indexes on all fields
db.orders.find({
    $or: [
        { customerId: "cust123" },
        { email: "user@example.com" }  // No index!
    ]
})

// ✅ GOOD: Ensure indexes on all $or branches
db.orders.createIndex({ email: 1 })

// ❌ BAD: Negation queries (can't use index efficiently)
db.orders.find({ status: { $ne: "cancelled" } })

// ✅ GOOD: Use $in with expected values
db.orders.find({ status: { $in: ["pending", "processing", "shipped"] } })

// ❌ BAD: Transforming indexed field
db.users.find({ 
    $expr: { $eq: [{ $toLower: "$email" }, "user@example.com"] }
})

// ✅ GOOD: Store normalized value
db.users.createIndex({ emailLower: 1 })
db.users.find({ emailLower: "user@example.com" })
```

### Pagination Strategies

```javascript
// ❌ BAD: Skip/Limit for large offsets
db.orders.find().skip(100000).limit(20)
// Scans 100,000 documents to skip them!

// ✅ GOOD: Keyset pagination (range queries)
// First page
const firstPage = await db.orders
    .find({ customerId: "cust123" })
    .sort({ orderDate: -1, _id: -1 })
    .limit(20)
    .toArray();

// Next page - use last document's values
const lastDoc = firstPage[firstPage.length - 1];
const nextPage = await db.orders
    .find({
        customerId: "cust123",
        $or: [
            { orderDate: { $lt: lastDoc.orderDate } },
            { 
                orderDate: lastDoc.orderDate,
                _id: { $lt: lastDoc._id }
            }
        ]
    })
    .sort({ orderDate: -1, _id: -1 })
    .limit(20)
    .toArray();

// Index for keyset pagination
db.orders.createIndex({ customerId: 1, orderDate: -1, _id: -1 })
```

### Query Hints and Plan Control

```javascript
// Force specific index
db.orders.find({ customerId: "cust123", status: "pending" })
    .hint({ customerId: 1, status: 1 })

// Force collection scan (testing)
db.orders.find({ customerId: "cust123" })
    .hint({ $natural: 1 })

// Set max time for query
db.orders.find({ customerId: "cust123" })
    .maxTimeMS(5000)

// Limit documents to examine
db.orders.find({ customerId: "cust123" })
    .maxScan(10000)  // Deprecated, use maxTimeMS
```

---

## 8. Aggregation Pipeline Optimization

### Theory: Data Transformation and Analysis

**What is the Aggregation Framework?**
The aggregation framework is MongoDB's answer to SQL's GROUP BY, JOIN, and analytical functions. It processes documents through a **pipeline** of stages, where each stage transforms the documents and passes results to the next stage.

**Pipeline vs Set-Based Processing**
Unlike SQL's declarative set-based processing, MongoDB's aggregation uses imperative pipeline processing:

| SQL Approach | Aggregation Approach |
|--------------|---------------------|
| Optimizer chooses execution order | Developer specifies stage order |
| Query describes "what" | Pipeline describes "how" |
| Optimization is automatic | Manual optimization is critical |

**The Stage Model**
Each pipeline stage is an independent operator:

- **Input**: Receives documents from previous stage (or collection)
- **Process**: Applies transformation, filtering, or computation
- **Output**: Passes documents to next stage

**Memory and Performance Characteristics**
Different stages have different resource profiles:

| Stage | Memory Usage | Blocking? | Notes |
|-------|-------------|-----------|-------|
| $match | Low | No | Streaming, can use indexes |
| $project | Low | No | Streaming |
| $group | High | Yes | Must see all documents |
| $sort | High | Yes | 100MB limit without disk |
| $lookup | Variable | Depends | Can be very expensive |
| $unwind | Low | No | Can multiply documents |

**Optimization Principles**

1. **Push Down Filters**: Move $match stages as early as possible - filtering before expensive operations reduces work
2. **Minimize Document Size**: Use $project early to remove unneeded fields, reducing memory and network usage
3. **Leverage Indexes**: Only initial $match and $sort can use indexes - sequence matters
4. **Avoid Unnecessary Stages**: Each stage has overhead - combine operations when possible
5. **Consider $merge for Materialization**: Pre-compute expensive aggregations into collections

**Pipeline Optimization by MongoDB**
MongoDB automatically optimizes some pipeline sequences:
- **Sequence Optimization**: Reorders stages when safe (e.g., moving $match before $project)
- **Coalescence**: Combines adjacent stages (e.g., consecutive $match stages)
- **Projection Optimization**: Removes fields from documents if not used by subsequent stages

### Pipeline Stage Order

```javascript
// ✅ GOOD: Filter early with $match
db.orders.aggregate([
    { $match: { status: "completed", orderDate: { $gte: ISODate("2024-01-01") } } },
    { $group: { _id: "$customerId", total: { $sum: "$amount" } } },
    { $sort: { total: -1 } },
    { $limit: 10 }
])

// ❌ BAD: Filtering after expensive operations
db.orders.aggregate([
    { $group: { _id: "$customerId", total: { $sum: "$amount" } } },  // All docs!
    { $sort: { total: -1 } },
    { $match: { total: { $gte: 1000 } } }  // Too late!
])

// Optimization order:
// 1. $match (filter as early as possible)
// 2. $project or $addFields (reduce document size)
// 3. $group, $lookup, etc.
// 4. $sort (after reducing data volume)
// 5. $limit
```

### Using Indexes in Aggregation

```javascript
// Pipeline uses index only for initial stages:
// $match, $sort (at beginning), $geoNear, $sample

// ✅ Index-supported aggregation
db.orders.createIndex({ customerId: 1, orderDate: -1 })

db.orders.aggregate([
    { $match: { customerId: "cust123" } },  // Uses index
    { $sort: { orderDate: -1 } },           // Uses same index
    { $limit: 100 },
    { $group: { ... } }
])

// ❌ Index can't help after $group
db.orders.aggregate([
    { $group: { _id: "$customerId", total: { $sum: 1 } } },
    { $match: { total: { $gte: 10 } } }  // No index help
])
```

### Optimizing $lookup

```javascript
// ❌ BAD: $lookup without foreign field index
db.orders.aggregate([
    {
        $lookup: {
            from: "customers",
            localField: "customerId",
            foreignField: "_id",  // Needs index!
            as: "customer"
        }
    }
])
// Ensure customers collection has _id index (automatic)
// Or custom index for non-_id foreign fields

// ✅ GOOD: Pipeline $lookup with filtering
db.orders.aggregate([
    { $match: { status: "pending" } },
    {
        $lookup: {
            from: "customers",
            let: { custId: "$customerId" },
            pipeline: [
                { $match: { 
                    $expr: { $eq: ["$_id", "$$custId"] },
                    isActive: true  // Filter in subpipeline
                }},
                { $project: { name: 1, email: 1 } }  // Only needed fields
            ],
            as: "customer"
        }
    }
])

// ✅ GOOD: Limit lookups with $limit first
db.orders.aggregate([
    { $match: { status: "pending" } },
    { $sort: { orderDate: -1 } },
    { $limit: 100 },  // Reduce before $lookup
    { $lookup: { ... } }
])
```

### Memory Management

```javascript
// Default memory limit: 100MB per pipeline stage
// Exceeding throws error

// Allow disk use for large aggregations
db.orders.aggregate([
    { $group: { _id: "$customerId", orders: { $push: "$$ROOT" } } },
    { $sort: { "orders.0.total": -1 } }
], { allowDiskUse: true })

// ✅ BETTER: Reduce memory needs by design
db.orders.aggregate([
    // Project only needed fields early
    { $project: { customerId: 1, total: 1, orderDate: 1 } },
    { $group: { 
        _id: "$customerId", 
        orderCount: { $sum: 1 },
        totalAmount: { $sum: "$total" },
        // Don't accumulate full documents
    }},
    { $sort: { totalAmount: -1 } }
])
```

### Aggregation Expressions

```javascript
// Use $expr for complex conditions
db.products.aggregate([
    {
        $match: {
            $expr: {
                $and: [
                    { $gte: ["$stock", "$minStock"] },
                    { $lt: ["$price", { $multiply: ["$cost", 2] }] }
                ]
            }
        }
    }
])

// Accumulator operators
db.orders.aggregate([
    {
        $group: {
            _id: "$customerId",
            orderCount: { $sum: 1 },
            totalAmount: { $sum: "$total" },
            avgAmount: { $avg: "$total" },
            minOrder: { $min: "$total" },
            maxOrder: { $max: "$total" },
            firstOrder: { $first: "$orderDate" },
            lastOrder: { $last: "$orderDate" },
            orders: { $push: "$orderId" },
            uniqueStatuses: { $addToSet: "$status" }
        }
    }
])

// Window functions (MongoDB 5.0+)
db.sales.aggregate([
    {
        $setWindowFields: {
            partitionBy: "$region",
            sortBy: { date: 1 },
            output: {
                runningTotal: {
                    $sum: "$amount",
                    window: { documents: ["unbounded", "current"] }
                },
                movingAvg: {
                    $avg: "$amount",
                    window: { documents: [-6, 0] }  // 7-day moving avg
                }
            }
        }
    }
])
```

---

## 9. Write Performance Tuning

### Theory: Write Path and Durability

**The Write Path in MongoDB**
When a write operation occurs, it follows this path:

1. **Driver sends operation** to mongos (sharded) or mongod directly
2. **Router determines shard** (for sharded collections using shard key)
3. **Primary receives write** and applies to in-memory data
4. **Journal write** (WAL) for durability before acknowledging
5. **Oplog entry created** for replication
6. **Replication to secondaries** via oplog tailing
7. **Acknowledgment returned** based on write concern

**Write Concern: The Durability Spectrum**
Write concern determines when MongoDB considers a write "successful":

```
   Fast ◄──────────────────────────────────────────────────► Safe
   
   w:0          w:1          w:majority        w:majority+j:true
   (fire &      (primary     (majority of      (majority with
   forget)      ack)         replicas)         journal commit)
```

**Understanding the Durability Guarantee**
- **w:0**: No guarantee - write may be lost even if client thinks it succeeded
- **w:1**: Primary has the data in memory - lost if primary crashes before journal/replication
- **w:majority**: Majority of replicas have data - survives single node failure
- **j:true**: Data is in journal (write-ahead log) - survives crash without data loss

**Write Amplification**
Each logical write causes multiple physical operations:
1. Journal write (WAL)
2. Data file write
3. Index updates (per index)
4. Oplog write
5. Network transfer to secondaries
6. Secondary applies same operations

More indexes = more write amplification = slower writes.

**Batching and Efficiency**
Bulk writes are more efficient because:
- **Reduced round trips**: One network call vs many
- **Batch journal commits**: Multiple operations share journal sync
- **Connection efficiency**: Better utilization of connection pool
- **Lock efficiency**: Fewer lock acquisitions

**Write Distribution in Sharded Clusters**
Write performance in sharded clusters depends on:
- **Targeted writes**: Include shard key, route to single shard
- **Broadcast writes**: Update/delete without shard key hits all shards
- **Chunk migrations**: Background balancing competes for I/O

### Write Concern Levels

```javascript
// Write concern: How many nodes must acknowledge write

// w: 0 - Fire and forget (fastest, least safe)
db.logs.insertOne(
    { message: "debug info" },
    { writeConcern: { w: 0 } }
)

// w: 1 - Primary acknowledges (default)
db.orders.insertOne(
    { orderId: "ORD-001" },
    { writeConcern: { w: 1 } }
)

// w: "majority" - Majority of replica set acknowledges
db.payments.insertOne(
    { paymentId: "PAY-001", amount: 100 },
    { writeConcern: { w: "majority" } }
)

// w: <number> - Specific number of nodes
db.orders.insertOne(
    { orderId: "ORD-001" },
    { writeConcern: { w: 2 } }  // Primary + 1 secondary
)

// j: true - Wait for journal commit
db.payments.insertOne(
    { paymentId: "PAY-001" },
    { writeConcern: { w: "majority", j: true } }
)

// wtimeout - Max wait time
db.orders.insertOne(
    { orderId: "ORD-001" },
    { writeConcern: { w: "majority", wtimeout: 5000 } }
)
```

### Bulk Write Operations

```javascript
// ✅ GOOD: Bulk writes for multiple operations
const bulk = db.orders.initializeUnorderedBulkOp();

for (let i = 0; i < 10000; i++) {
    bulk.insert({
        orderId: `ORD-${i}`,
        customerId: `CUST-${i % 100}`,
        amount: Math.random() * 1000
    });
}

const result = await bulk.execute();

// Ordered bulk (stops on first error)
const orderedBulk = db.orders.initializeOrderedBulkOp();

// Mixed operations
const bulkOps = [
    { insertOne: { document: { orderId: "ORD-001" } } },
    { updateOne: { 
        filter: { orderId: "ORD-002" }, 
        update: { $set: { status: "shipped" } }
    }},
    { deleteOne: { filter: { orderId: "ORD-OLD" } } },
    { replaceOne: {
        filter: { orderId: "ORD-003" },
        replacement: { orderId: "ORD-003", status: "new" }
    }}
];

db.orders.bulkWrite(bulkOps, { ordered: false });

// Batch size for very large operations
const BATCH_SIZE = 1000;
for (let i = 0; i < documents.length; i += BATCH_SIZE) {
    const batch = documents.slice(i, i + BATCH_SIZE);
    await db.collection.insertMany(batch, { ordered: false });
}
```

### Update Optimization

```javascript
// ✅ GOOD: Use update operators, not replace
db.products.updateOne(
    { _id: productId },
    { 
        $set: { price: 99.99 },
        $inc: { viewCount: 1 },
        $currentDate: { lastModified: true }
    }
)

// ❌ BAD: Replace entire document
const product = await db.products.findOne({ _id: productId });
product.price = 99.99;
product.viewCount++;
product.lastModified = new Date();
await db.products.replaceOne({ _id: productId }, product);
// Sends entire document over network!

// Atomic operations
db.inventory.updateOne(
    { _id: productId, stock: { $gte: quantity } },  // Check in query
    { $inc: { stock: -quantity } }                  // Atomic decrement
)

// Array operations
db.orders.updateOne(
    { _id: orderId },
    { 
        $push: { 
            items: { 
                $each: [newItem1, newItem2],
                $position: 0,
                $slice: 100  // Keep max 100 items
            }
        }
    }
)

// Conditional updates
db.products.updateOne(
    { _id: productId },
    [
        {
            $set: {
                discount: {
                    $cond: {
                        if: { $gte: ["$stock", 100] },
                        then: 0.2,
                        else: 0
                    }
                }
            }
        }
    ]
)
```

### Write Distribution in Sharded Clusters

```javascript
// Targeted writes (single shard) - most efficient
db.orders.insertOne({
    customerId: "cust123",  // Shard key
    orderId: "ORD-001",
    amount: 100
})

// Broadcast writes (all shards) - avoid if possible
db.orders.updateMany(
    { status: "pending" },  // No shard key!
    { $set: { status: "processing" } }
)

// Multi-document transactions write concern
const session = client.startSession();
try {
    session.startTransaction({
        writeConcern: { w: "majority" },
        readConcern: { level: "snapshot" }
    });
    
    await db.orders.insertOne({ ... }, { session });
    await db.inventory.updateOne({ ... }, { session });
    
    await session.commitTransaction();
} catch (error) {
    await session.abortTransaction();
    throw error;
} finally {
    session.endSession();
}
```

---

## 10. Read Performance Tuning

### Theory: Read Path and Consistency

**The Read Path in MongoDB**
Read operations follow this path:

1. **Driver sends query** to mongos (sharded) or mongod
2. **Router determines targets** (which shard(s) for sharded collections)
3. **Query planning** identifies optimal execution strategy
4. **Index traversal** locates matching document references
5. **Document fetch** retrieves documents from storage (if not covered query)
6. **Result streaming** returns documents to client
7. **Cross-shard merge** (for scatter-gather queries)

**Read Preference: Availability vs Consistency**
Read preference controls which replica set member serves reads:

| Preference | Consistency | Availability | Latency | Use Case |
|------------|-------------|--------------|---------|----------|
| primary | Strong | Lower | Variable | Transactions, critical reads |
| primaryPreferred | Strong (usually) | Higher | Variable | Default for most apps |
| secondary | Eventual | Higher | Often lower | Analytics, reporting |
| secondaryPreferred | Eventual | Highest | Often lower | Read scaling |
| nearest | Eventual | Highest | Lowest | Latency-sensitive apps |

**Read Concern: Point-in-Time Consistency**
Read concern determines what data a query sees:

- **local**: Returns whatever data exists on the queried node (may be rolled back)
- **available**: Like local, optimized for sharded clusters
- **majority**: Returns data acknowledged by majority (won't be rolled back)
- **linearizable**: Returns data reflecting all writes completed before read started
- **snapshot**: Returns data from a consistent snapshot (for transactions)

**The Consistency-Availability Trade-off**
Stronger consistency = potentially slower reads:

```
Fast ◄────────────────────────────────────────────► Consistent

local/           majority          linearizable
available        (waits for        (blocks until
(immediate)      replication)      confirmed)
```

**Caching Architecture**
MongoDB doesn't have a built-in query cache, but uses:

1. **WiredTiger Cache**: Stores recently accessed data in memory
2. **OS Page Cache**: Additional caching layer at OS level
3. **Application Caching**: External caches (Redis, Memcached) for frequently accessed data

**Projection and Network Efficiency**
Projection isn't just about returning less data - it affects:
- Network bandwidth (smaller documents = faster transfer)
- Client memory (fewer fields = less parsing overhead)
- Covered queries (projection matching index = no document fetch)

### Read Preference

```javascript
// Read preference: Which replica set members to read from

// primary (default) - Only read from primary
db.orders.find({ customerId: "cust123" })
    .readPref("primary")

// primaryPreferred - Primary if available, else secondary
db.orders.find({ customerId: "cust123" })
    .readPref("primaryPreferred")

// secondary - Only read from secondaries
db.analytics.find({ date: "2024-03-15" })
    .readPref("secondary")

// secondaryPreferred - Secondary if available, else primary
db.reports.find({ type: "daily" })
    .readPref("secondaryPreferred")

// nearest - Lowest network latency member
db.products.find({ category: "electronics" })
    .readPref("nearest")

// With tags (read from specific datacenter)
db.orders.find({ customerId: "cust123" })
    .readPref("secondary", [{ dc: "east" }])

// Connection string options
const uri = "mongodb://host1,host2,host3/mydb?readPreference=secondaryPreferred&readPreferenceTags=dc:east";
```

### Read Concern Levels

```javascript
// Read concern: What data visibility guarantees

// "local" - Returns most recent data on node (default)
db.orders.find({ customerId: "cust123" })
    .readConcern("local")

// "available" - Like local, but for sharded clusters
db.orders.find({ customerId: "cust123" })
    .readConcern("available")

// "majority" - Returns data acknowledged by majority
db.orders.find({ customerId: "cust123" })
    .readConcern("majority")

// "linearizable" - Strongest consistency (blocks until confirmed)
db.orders.findOne({ orderId: "ORD-001" })
    .readConcern("linearizable")

// "snapshot" - For multi-document transactions
const session = client.startSession();
session.startTransaction({ readConcern: { level: "snapshot" } });
```

### Projection Optimization

```javascript
// ✅ GOOD: Project only needed fields
db.users.find(
    { status: "active" },
    { name: 1, email: 1, _id: 0 }  // Fewer bytes transferred
)

// ✅ GOOD: Exclude large fields
db.articles.find(
    { category: "tech" },
    { content: 0, comments: 0 }  // Exclude large fields
)

// Covered query with projection
db.orders.createIndex({ customerId: 1, orderDate: 1, status: 1, total: 1 })

db.orders.find(
    { customerId: "cust123" },
    { customerId: 1, orderDate: 1, status: 1, total: 1, _id: 0 }
)
// Entirely from index, no document fetch!

// ❌ BAD: Returning full documents when not needed
db.orders.find({ customerId: "cust123" })
// Returns all fields including large arrays, embedded docs
```

### Query Caching and Optimization

```javascript
// MongoDB doesn't have traditional query cache
// Use application-level caching

// Redis caching pattern
async function getUser(userId) {
    // Try cache first
    const cached = await redis.get(`user:${userId}`);
    if (cached) {
        return JSON.parse(cached);
    }
    
    // Query database
    const user = await db.users.findOne({ _id: userId });
    
    // Cache result (5 minute TTL)
    await redis.setex(`user:${userId}`, 300, JSON.stringify(user));
    
    return user;
}

// In-memory aggregation caching
const CACHE_TTL = 60000; // 1 minute
const cache = new Map();

async function getDailyStats(date) {
    const cacheKey = `stats:${date}`;
    const cached = cache.get(cacheKey);
    
    if (cached && Date.now() - cached.timestamp < CACHE_TTL) {
        return cached.data;
    }
    
    const stats = await db.orders.aggregate([
        { $match: { orderDate: { $gte: startOfDay, $lt: endOfDay } } },
        { $group: { _id: null, total: { $sum: "$amount" }, count: { $sum: 1 } } }
    ]).toArray();
    
    cache.set(cacheKey, { data: stats, timestamp: Date.now() });
    return stats;
}
```

---

## 11. Replication & High Availability

### Theory: Replication and Fault Tolerance

**Why Replication Matters**
Replication is the foundation of high availability in distributed systems. Without replication, a single server failure means complete data loss and service unavailability. MongoDB's replication provides:

- **Data Redundancy**: Multiple copies of data survive hardware failures
- **High Availability**: Automatic failover maintains service continuity
- **Read Scaling**: Distribute read operations across replicas
- **Disaster Recovery**: Geographic distribution protects against site failures

**The Replica Set Model**
MongoDB uses a **primary-secondary replication** model:

- **Primary**: Single node that accepts all write operations
- **Secondaries**: Replicate data from primary asynchronously
- **Arbiter**: Votes in elections but holds no data (for odd member count)

**The Oplog: Heart of Replication**
The **operation log (oplog)** is a capped collection that records all write operations:

1. Primary writes operation to oplog
2. Secondaries tail the oplog and apply operations
3. Oplog entries are idempotent (can be safely replayed)
4. Oplog size determines recovery window

**Election Protocol**
When primary becomes unavailable, an election occurs:

1. Secondaries detect primary failure (heartbeat timeout)
2. Eligible members (priority > 0) request votes
3. Member with highest priority and most recent oplog wins
4. New primary starts accepting writes

Election typically completes in **10-12 seconds** with default settings.

**Consistency vs Availability Trade-off**
| Configuration | Behavior | Risk |
|--------------|----------|------|
| w:1, read from primary | Strong consistency | Data loss if primary fails before replication |
| w:majority | Durable writes | Slower writes, unavailable if majority offline |
| Read from secondary | Higher availability | Stale reads possible |

**Split Brain Prevention**
A "split brain" occurs when network partitions create multiple primaries. MongoDB prevents this through:

- **Majority requirement**: A node needs votes from majority to become primary
- **Term numbers**: Each election increments term; higher terms supersede
- **Write concern majority**: Writes require majority acknowledgment

### Replica Set Configuration

```javascript
// Initialize replica set
rs.initiate({
    _id: "myReplicaSet",
    members: [
        { _id: 0, host: "mongo1:27017", priority: 2 },  // Primary preference
        { _id: 1, host: "mongo2:27017", priority: 1 },
        { _id: 2, host: "mongo3:27017", priority: 1 },
        { _id: 3, host: "mongo4:27017", priority: 0, hidden: true },  // Hidden for backup
        { _id: 4, host: "mongo5:27017", arbiterOnly: true }  // Arbiter
    ],
    settings: {
        chainingAllowed: true,
        heartbeatTimeoutSecs: 10,
        electionTimeoutMillis: 10000
    }
})

// Member configuration options:
// priority: Election priority (0 = can never be primary)
// hidden: Hidden from clients (use for analytics/backup)
// slaveDelay/secondaryDelaySecs: Delayed replication
// votes: Voting member (0 or 1)
// arbiterOnly: Arbiter (votes but holds no data)
// buildIndexes: Whether to build indexes

// Delayed replica (for point-in-time recovery)
rs.add({
    host: "mongo-delayed:27017",
    priority: 0,
    hidden: true,
    secondaryDelaySecs: 3600  // 1 hour delay
})
```

### Replica Set Operations

```javascript
// Check replica set status
rs.status()

// Check replication info
rs.printReplicationInfo()
rs.printSecondaryReplicationInfo()

// Step down primary (for maintenance)
rs.stepDown(60)  // Don't re-elect for 60 seconds

// Reconfigure replica set
const config = rs.conf();
config.members[1].priority = 2;
rs.reconfig(config)

// Force reconfiguration (when majority unavailable)
rs.reconfig(config, { force: true })

// Add/remove members
rs.add("mongo-new:27017")
rs.remove("mongo-old:27017")

// Freeze member from becoming primary
rs.freeze(120)  // 120 seconds
```

### Failover Handling

```javascript
// Application retry logic
const client = new MongoClient(uri, {
    retryWrites: true,
    retryReads: true,
    serverSelectionTimeoutMS: 30000,
    socketTimeoutMS: 360000
});

// Handle failover errors
async function resilientOperation() {
    const maxRetries = 3;
    let retries = 0;
    
    while (retries < maxRetries) {
        try {
            return await db.orders.insertOne({ orderId: "ORD-001" });
        } catch (error) {
            if (isRetryableError(error)) {
                retries++;
                await sleep(Math.pow(2, retries) * 100);  // Exponential backoff
                continue;
            }
            throw error;
        }
    }
}

function isRetryableError(error) {
    const retryableCodes = [
        11600, // InterruptedAtShutdown
        11602, // InterruptedDueToReplStateChange
        10107, // NotWritablePrimary
        13436, // NotPrimaryNoSecondaryOk
        189,   // PrimarySteppedDown
        91,    // ShutdownInProgress
    ];
    return retryableCodes.includes(error.code);
}
```

### Read Scaling with Secondaries

```javascript
// Configure secondary reads for read-heavy workloads

// Read preference per collection
db.getMongo().setReadPref("secondaryPreferred")

// Read preference per query
db.analytics.find({ date: "2024-03-15" })
    .readPref("secondaryPreferred", [
        { dc: "local" },      // First try local datacenter
        { dc: "backup" },     // Then backup datacenter
        {}                    // Finally, any secondary
    ])

// Read preference in connection string
const uri = "mongodb://mongo1,mongo2,mongo3/mydb" +
    "?replicaSet=myReplicaSet" +
    "&readPreference=secondaryPreferred" +
    "&readPreferenceTags=dc:east,usage:analytics" +
    "&readPreferenceTags=dc:east" +
    "&readPreferenceTags=";
```

---

## 12. Connection Management

### Theory: Connection Architecture and Pooling

**Why Connection Pooling Matters**
Database connections are expensive resources:

- **TCP Handshake**: 3-way handshake adds latency (1.5 RTT)
- **TLS Negotiation**: Key exchange adds significant overhead
- **Authentication**: Connection-level auth requires round trips
- **Memory Allocation**: Server allocates memory per connection

Without pooling, creating a new connection for each operation adds **50-200ms** overhead.

**The Connection Pool Model**
A connection pool maintains a set of reusable connections:

```
┌─────────────────┐       ┌───────────────────────┐
│   Application   │       │   Connection Pool     │
│                 │──────►│  ┌─┐ ┌─┐ ┌─┐ ┌─┐     │────► MongoDB
│   Thread 1      │       │  │C│ │C│ │C│ │C│     │
│   Thread 2      │       │  │1│ │2│ │3│ │4│     │
│   Thread 3      │       │  └─┘ └─┘ └─┘ └─┘     │
└─────────────────┘       └───────────────────────┘
```

**Pool Sizing Principles**

1. **Min Pool Size**: Connections to maintain even when idle. Prevents cold start latency.
2. **Max Pool Size**: Upper limit on connections. Prevents overwhelming server.
3. **Optimal Size**: Depends on:
   - Number of concurrent operations
   - Average operation duration
   - Server capacity (connections per mongod)

**Formula**: `Pool Size = Concurrent Operations × Avg Operation Duration + Buffer`

**Connection Lifecycle**
1. **Creation**: Pool creates minPoolSize connections at startup
2. **Checkout**: Thread requests connection; pool provides available one or creates new
3. **Use**: Thread uses connection for operations
4. **Return**: Thread returns connection to pool (not closed)
5. **Idle Timeout**: Connections idle > maxIdleTimeMS are closed
6. **Health Check**: Heartbeats verify connection validity

**Connection Considerations for Distributed Systems**
- **Per-Server Pools**: Driver maintains separate pool per mongod
- **Mongos Routing**: Sharded clusters may need larger pools (mongos fans out)
- **Serverless Considerations**: Lambda/Cloud Functions need smaller pools (maxPoolSize: 1-5)
- **Connection Storms**: Startup spike when all connections created simultaneously

**Timeout Configuration**
| Timeout | Purpose | Typical Value |
|---------|---------|---------------|
| connectTimeoutMS | Initial connection establishment | 30000ms |
| socketTimeoutMS | Max wait for response | 360000ms |
| serverSelectionTimeoutMS | Time to find suitable server | 30000ms |
| waitQueueTimeoutMS | Time to wait for available connection | 10000ms |

### Connection Pooling

```javascript
// Node.js driver connection pool settings
const client = new MongoClient(uri, {
    // Pool settings
    maxPoolSize: 100,           // Max connections per server
    minPoolSize: 10,            // Min connections maintained
    maxIdleTimeMS: 60000,       // Close idle connections after 60s
    waitQueueTimeoutMS: 10000,  // Max wait for connection
    
    // Timeouts
    connectTimeoutMS: 30000,    // Initial connection timeout
    socketTimeoutMS: 360000,    // Socket timeout
    serverSelectionTimeoutMS: 30000,  // Server selection timeout
    
    // Heartbeat
    heartbeatFrequencyMS: 10000,  // Check server health every 10s
    
    // Compression
    compressors: ['snappy', 'zstd', 'zlib'],
    
    // TLS
    tls: true,
    tlsCAFile: '/path/to/ca.pem',
    tlsCertificateKeyFile: '/path/to/client.pem'
});

// Connection events
client.on('connectionPoolCreated', (event) => {
    console.log(`Pool created for ${event.address}`);
});

client.on('connectionPoolCleared', (event) => {
    console.log(`Pool cleared for ${event.address}`);
});

client.on('connectionCheckedOut', (event) => {
    console.log('Connection checked out');
});
```

### Connection String Best Practices

```javascript
// Full connection string with all recommended options
const uri = "mongodb+srv://user:password@cluster.example.com/mydb" +
    // Replica set
    "?replicaSet=myReplicaSet" +
    
    // Connection pool
    "&maxPoolSize=100" +
    "&minPoolSize=10" +
    "&maxIdleTimeMS=60000" +
    "&waitQueueTimeoutMS=10000" +
    
    // Timeouts
    "&connectTimeoutMS=30000" +
    "&socketTimeoutMS=360000" +
    "&serverSelectionTimeoutMS=30000" +
    
    // Write concern
    "&w=majority" +
    "&wtimeoutMS=10000" +
    "&journal=true" +
    
    // Read preference
    "&readPreference=secondaryPreferred" +
    "&readPreferenceTags=dc:east" +
    
    // Read concern
    "&readConcernLevel=majority" +
    
    // Compression
    "&compressors=snappy,zstd" +
    
    // Retries
    "&retryWrites=true" +
    "&retryReads=true" +
    
    // App name (for monitoring)
    "&appName=MyApplication";
```

### Connection Lifecycle Management

```javascript
// Singleton connection pattern
let client = null;
let db = null;

async function connectToDatabase() {
    if (client && client.isConnected()) {
        return db;
    }
    
    client = new MongoClient(uri, options);
    await client.connect();
    db = client.db('mydb');
    
    // Handle process shutdown
    process.on('SIGINT', closeConnection);
    process.on('SIGTERM', closeConnection);
    
    return db;
}

async function closeConnection() {
    if (client) {
        await client.close();
        client = null;
        db = null;
    }
    process.exit(0);
}

// Express.js middleware pattern
app.use(async (req, res, next) => {
    try {
        req.db = await connectToDatabase();
        next();
    } catch (error) {
        next(error);
    }
});

// Lambda/serverless pattern (connection reuse)
let cachedClient = null;

async function handler(event) {
    if (!cachedClient) {
        cachedClient = new MongoClient(uri, {
            maxPoolSize: 1,  // Limit for serverless
            serverSelectionTimeoutMS: 5000
        });
        await cachedClient.connect();
    }
    
    const db = cachedClient.db('mydb');
    // Use db...
}
```

---

## 13. Capacity Planning

### Theory: Sizing and Growth Planning

**Why Capacity Planning is Critical**
Capacity planning prevents two costly scenarios:

1. **Under-provisioning**: System crashes or slows under load, causing outages
2. **Over-provisioning**: Wasted infrastructure spend on unused resources

**The Working Set Concept**
The "working set" is the most frequently accessed data + indexes. Understanding it is crucial:

- **Definition**: Data accessed in a typical time window (often last hour/day)
- **Ideal State**: Working set fits entirely in RAM
- **Degraded State**: Working set exceeds RAM, causing disk I/O
- **Impact**: 10-100x performance difference between RAM and disk access

**Sizing Dimensions**

| Dimension | Key Metric | Scaling Path |
|-----------|------------|---------------|
| **Storage** | Data + Index size | Add shards or disk |
| **Memory** | Working set size | Vertical scaling or more shards |
| **CPU** | Operations per second | Vertical or horizontal scaling |
| **Network** | Bandwidth consumption | Better compression, add capacity |
| **IOPS** | Disk operations/second | SSD, RAID, more shards |

**Growth Projection Methodology**

1. **Baseline Measurement**: Current size, growth rate, peak usage
2. **Growth Modeling**: Linear, exponential, or seasonal patterns
3. **Projection**: Size at 6mo, 1yr, 2yr horizons
4. **Buffer Planning**: 20-40% headroom for unexpected growth
5. **Trigger Points**: When to scale (70-80% utilization)

**The Rule of Thirds for Memory**
```
Total RAM Allocation:
├── WiredTiger Cache (50% of RAM - 1GB)
├── Indexes (should fit in remaining memory)
└── OS/System (leave 1-2GB minimum)
```

**Storage Overhead Factors**
- **Replication Factor**: 3x for standard replica sets
- **Index Overhead**: Typically 10-30% of data size
- **Compression**: WiredTiger typically achieves 50-70% compression
- **Oplog**: Size based on write rate and recovery window needs
- **Journal**: Typically 1-2GB
- **Pre-allocated Space**: WiredTiger pre-allocates files

**When to Shard**
Consider sharding when:
- Single server storage is insufficient
- Working set exceeds available RAM
- Write throughput saturates single primary
- Read scaling beyond replica set capacity

### Storage Estimation

```javascript
// Estimate document size
const sampleDoc = {
    _id: ObjectId(),                    // 12 bytes
    orderId: "ORD-2024-001234567",      // ~20 bytes
    customerId: ObjectId(),             // 12 bytes
    status: "pending",                  // ~8 bytes
    items: [                            // Array overhead + items
        { productId: ObjectId(), quantity: 1, price: 99.99 }  // ~40 bytes each
    ],
    totals: { subtotal: 99.99, tax: 8.00, total: 107.99 },  // ~50 bytes
    createdAt: ISODate(),               // 8 bytes
    updatedAt: ISODate()                // 8 bytes
};

// Check actual document size
Object.bsonsize(sampleDoc)  // Returns bytes

// Collection statistics
db.orders.stats()
// avgObjSize: Average document size
// storageSize: Disk space used
// totalIndexSize: Index space used

// Capacity formula:
// Total Storage = (Avg Doc Size × Doc Count × Replication Factor) 
//                 + (Index Size × Replication Factor)
//                 + Overhead (20-30%)
```

### Scaling Calculations

```
┌─────────────────────────────────────────────────────────────────┐
│                    CAPACITY PLANNING WORKSHEET                  │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  STORAGE REQUIREMENTS:                                          │
│  ├── Documents per day: 1,000,000                               │
│  ├── Avg document size: 2 KB                                    │
│  ├── Daily data growth: 2 GB                                    │
│  ├── Retention period: 365 days                                 │
│  ├── Total data: 730 GB                                         │
│  ├── Indexes (~30%): 220 GB                                     │
│  ├── Replication factor: 3                                      │
│  └── Total storage: (730 + 220) × 3 = 2.85 TB per replica set  │
│                                                                  │
│  MEMORY REQUIREMENTS:                                           │
│  ├── Working set (active data): 50 GB                           │
│  ├── Index memory: 220 GB (ideally fit in RAM)                  │
│  ├── WiredTiger cache: 50% of RAM                               │
│  └── Recommended RAM: 128 GB minimum                            │
│                                                                  │
│  THROUGHPUT REQUIREMENTS:                                       │
│  ├── Peak reads/sec: 50,000                                     │
│  ├── Peak writes/sec: 10,000                                    │
│  ├── Read latency target: < 10ms                                │
│  ├── Write latency target: < 50ms                               │
│  └── Shards needed: 3-5 (based on load testing)                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Monitoring Growth

```javascript
// Track collection growth
db.adminCommand({
    aggregate: 1,
    pipeline: [
        { $currentOp: { allUsers: true, idleConnections: true } }
    ],
    cursor: {}
})

// Database size over time
db.stats()

// Collection size projection
const dailyGrowth = db.orders.stats().size / daysSinceCreation;
const projectedSize30Days = currentSize + (dailyGrowth * 30);

// Shard distribution
db.orders.getShardDistribution()

// Chunk statistics
use config
db.chunks.aggregate([
    { $match: { ns: "mydb.orders" } },
    { $group: {
        _id: "$shard",
        count: { $sum: 1 },
        totalSize: { $sum: "$size" }
    }}
])
```

---

## 14. Monitoring & Observability

### Theory: The Three Pillars of Observability

**Observability vs Monitoring**
- **Monitoring**: Tracking known metrics and alerting on thresholds
- **Observability**: Ability to understand system state from external outputs

Modern systems require both: monitoring catches known issues, observability helps debug unknown problems.

**The Three Pillars**

1. **Metrics**: Numeric measurements over time
   - CPU utilization, memory usage, operation counts
   - Aggregatable, efficient storage, trend analysis
   - Example: "Average query latency was 15ms"

2. **Logs**: Discrete events with context
   - Detailed information about what happened
   - High cardinality, harder to aggregate
   - Example: "Query X took 500ms because index Y was missing"

3. **Traces**: Request flow across services
   - End-to-end request path visualization
   - Critical for distributed systems debugging
   - Example: "Request spent 200ms in MongoDB, 50ms in API"

**Key MongoDB Metrics Categories**

| Category | Metrics | What It Tells You |
|----------|---------|-------------------|
| **Throughput** | opcounters, document metrics | How much work system is doing |
| **Latency** | operation latencies | How fast operations complete |
| **Saturation** | connections, queue depth | How close to capacity |
| **Errors** | asserts, exceptions | System health problems |
| **Replication** | lag, oplog window | Data consistency status |
| **Resource** | CPU, memory, disk | Infrastructure utilization |

**Golden Signals (SRE Framework)**
Google's SRE practices define four golden signals:

1. **Latency**: Time to service a request
2. **Traffic**: Demand on the system (ops/second)
3. **Errors**: Rate of failed requests
4. **Saturation**: How "full" the service is

**Alerting Philosophy**
- **Alert on symptoms, not causes**: Alert when users are impacted, not just when CPU is high
- **Actionable alerts**: Each alert should have a documented response
- **Avoid alert fatigue**: Too many alerts = ignored alerts
- **Use severity levels**: Critical (wake someone up) vs Warning (review tomorrow)

**Baseline and Anomaly Detection**
Static thresholds often don't work well. Better approaches:
- Establish baseline behavior over time
- Alert on deviations from baseline
- Consider time-of-day patterns (peak vs off-peak)
- Use machine learning for anomaly detection

### MongoDB Metrics to Monitor

```javascript
// Server status
db.serverStatus()

// Key metrics:
{
    // Connections
    "connections": {
        "current": 50,        // Active connections
        "available": 950,     // Available connections
        "totalCreated": 1000  // Total created over time
    },
    
    // Operations
    "opcounters": {
        "insert": 10000,
        "query": 50000,
        "update": 15000,
        "delete": 1000,
        "getmore": 5000,
        "command": 100000
    },
    
    // Memory
    "mem": {
        "resident": 8192,     // Resident memory (MB)
        "virtual": 16384,     // Virtual memory (MB)
        "mapped": 0
    },
    
    // WiredTiger cache
    "wiredTiger": {
        "cache": {
            "bytes currently in the cache": 4294967296,
            "maximum bytes configured": 8589934592,
            "pages read into cache": 100000,
            "pages written from cache": 50000
        }
    },
    
    // Replication
    "repl": {
        "ismaster": true,
        "secondary": false,
        "rbid": 1
    }
}

// Replica set lag
db.printSecondaryReplicationInfo()

// Current operations
db.currentOp({ "secs_running": { "$gt": 5 } })

// Collection stats
db.orders.stats()

// Index usage
db.orders.aggregate([{ $indexStats: {} }])
```

### Profiling Slow Queries

```javascript
// Enable profiler
db.setProfilingLevel(1, { slowms: 100 })  // Log queries > 100ms

// Profiler levels:
// 0 - Off
// 1 - Slow operations only
// 2 - All operations

// Query profiler data
db.system.profile.find().sort({ ts: -1 }).limit(10)

// Find slow queries
db.system.profile.find({
    millis: { $gt: 100 },
    op: { $in: ["query", "update", "remove"] }
}).sort({ millis: -1 })

// Aggregation analysis
db.system.profile.aggregate([
    { $match: { millis: { $gt: 100 } } },
    { $group: {
        _id: "$ns",
        count: { $sum: 1 },
        avgMs: { $avg: "$millis" },
        maxMs: { $max: "$millis" }
    }},
    { $sort: { count: -1 } }
])
```

### Alerting Thresholds

```yaml
# Prometheus alerting rules example
groups:
  - name: mongodb_alerts
    rules:
      # High connection usage
      - alert: MongoDBHighConnections
        expr: mongodb_connections_current / mongodb_connections_available > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "MongoDB connection usage > 80%"
      
      # Replication lag
      - alert: MongoDBReplicationLag
        expr: mongodb_replset_member_replication_lag > 10
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "MongoDB replication lag > 10 seconds"
      
      # Cache pressure
      - alert: MongoDBCachePressure
        expr: mongodb_wiredtiger_cache_bytes / mongodb_wiredtiger_cache_max_bytes > 0.95
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "MongoDB cache usage > 95%"
      
      # Slow queries
      - alert: MongoDBSlowQueries
        expr: rate(mongodb_op_latencies_latency_total[5m]) > 100
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "MongoDB average query latency > 100ms"
```

### MongoDB Atlas Monitoring (Managed)

```javascript
// Atlas provides built-in monitoring:
// - Real-time Performance Panel
// - Query Profiler
// - Performance Advisor (index recommendations)
// - Schema Advisor
// - Alert configurations

// Using Atlas Data API for monitoring
const response = await fetch(
    'https://data.mongodb-api.com/app/data-xxx/endpoint/data/v1/action/aggregate',
    {
        method: 'POST',
        headers: {
            'api-key': process.env.ATLAS_API_KEY,
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            dataSource: 'Cluster0',
            database: 'admin',
            collection: '$cmd',
            pipeline: [
                { $currentOp: { allUsers: true } }
            ]
        })
    }
);
```

---

## 15. Security at Scale

### Theory: Defense in Depth for Database Security

**The Security Triad: CIA**
- **Confidentiality**: Only authorized users can access data
- **Integrity**: Data cannot be modified by unauthorized parties
- **Availability**: Authorized users can always access data

**Defense in Depth Layers**
Security should be implemented at multiple layers:

```
┌─────────────────────────────────────────────┐
│         Application Security                │
│  (input validation, parameterized queries)  │
├─────────────────────────────────────────────┤
│           Access Control (RBAC)             │
│   (users, roles, collection-level perms)    │
├─────────────────────────────────────────────┤
│         Authentication & Identity           │
│    (SCRAM, LDAP, X.509, Kerberos)          │
├─────────────────────────────────────────────┤
│            Network Security                 │
│    (TLS, IP whitelisting, VPC, firewall)    │
├─────────────────────────────────────────────┤
│         Encryption (Data Protection)        │
│   (at-rest, in-transit, field-level)        │
├─────────────────────────────────────────────┤
│         Auditing & Compliance               │
│    (logging, monitoring, alerting)          │
└─────────────────────────────────────────────┘
```

**The Principle of Least Privilege**
Every user and service should have only the minimum permissions needed:

- Read-only for reporting applications
- Collection-specific access for microservices
- No admin access for application accounts
- Separate credentials for different environments

**Authentication Methods Comparison**

| Method | Security Level | Management | Use Case |
|--------|---------------|------------|----------|
| SCRAM-SHA-256 | Good | Simple | Small deployments |
| LDAP | Good | Centralized | Enterprise with existing LDAP |
| Kerberos | High | Complex | Enterprise Windows environments |
| X.509 | Highest | Complex | Zero-trust architectures |

**Encryption Layers**

1. **In-Transit (TLS)**: Encrypts data between client and server
   - Prevents eavesdropping and man-in-the-middle attacks
   - Should always be enabled in production

2. **At-Rest**: Encrypts data on disk
   - Protects against physical theft and unauthorized disk access
   - Transparent to applications

3. **Field-Level (CSFLE)**: Encrypts specific fields in documents
   - Client-side encryption - server never sees plaintext
   - Protects even against database administrator access
   - Required for certain compliance (PCI-DSS, HIPAA)

**Injection Prevention**
MongoDB isn't immune to injection attacks:
- Always use parameterized queries
- Never construct queries from user input strings
- Validate and sanitize all input
- Use schema validation to enforce data types

### Authentication

```javascript
// SCRAM-SHA-256 authentication (default)
mongod --auth

// Create admin user
use admin
db.createUser({
    user: "admin",
    pwd: "securePassword123!",
    roles: [
        { role: "userAdminAnyDatabase", db: "admin" },
        { role: "readWriteAnyDatabase", db: "admin" }
    ]
})

// LDAP authentication
security:
  authorization: enabled
  ldap:
    servers: "ldap.example.com"
    bind:
      method: "simple"
      queryUser: "cn=ldap-reader,ou=users,dc=example,dc=com"
      queryPassword: "password"
    userToDNMapping:
      '[{ match: "(.+)", substitution: "uid={0},ou=users,dc=example,dc=com" }]'

// X.509 certificate authentication
db.getSiblingDB("$external").createUser({
    user: "CN=myClient,OU=clients,O=MyOrg,L=City,ST=State,C=US",
    roles: [
        { role: "readWrite", db: "mydb" }
    ]
})
```

### Authorization (RBAC)

```javascript
// Built-in roles
// Database roles: read, readWrite, dbAdmin, dbOwner, userAdmin
// Cluster roles: clusterAdmin, clusterManager, clusterMonitor, hostManager
// All-database roles: readAnyDatabase, readWriteAnyDatabase, userAdminAnyDatabase, dbAdminAnyDatabase
// Superuser: root

// Create application-specific roles
use mydb
db.createRole({
    role: "orderProcessor",
    privileges: [
        {
            resource: { db: "mydb", collection: "orders" },
            actions: ["find", "insert", "update"]
        },
        {
            resource: { db: "mydb", collection: "inventory" },
            actions: ["find", "update"]
        }
    ],
    roles: []
})

// Create user with custom role
db.createUser({
    user: "orderService",
    pwd: "servicePassword!",
    roles: [
        { role: "orderProcessor", db: "mydb" }
    ]
})

// Field-level redaction
db.createView(
    "orders_public",
    "orders",
    [
        {
            $project: {
                orderId: 1,
                status: 1,
                total: 1,
                // Exclude sensitive fields
                customerEmail: 0,
                paymentDetails: 0
            }
        }
    ]
)
```

### Encryption

```javascript
// Encryption at rest (Enterprise)
security:
  enableEncryption: true
  encryptionCipherMode: AES256-CBC
  encryptionKeyFile: /path/to/mongodb-keyfile

// Client-side field level encryption (CSFLE)
const client = new MongoClient(uri, {
    autoEncryption: {
        keyVaultNamespace: "encryption.__keyVault",
        kmsProviders: {
            aws: {
                accessKeyId: process.env.AWS_ACCESS_KEY,
                secretAccessKey: process.env.AWS_SECRET_KEY
            }
        },
        schemaMap: {
            "mydb.users": {
                bsonType: "object",
                encryptMetadata: {
                    keyId: [UUID("...")]
                },
                properties: {
                    ssn: {
                        encrypt: {
                            bsonType: "string",
                            algorithm: "AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic"
                        }
                    },
                    medicalRecords: {
                        encrypt: {
                            bsonType: "array",
                            algorithm: "AEAD_AES_256_CBC_HMAC_SHA_512-Random"
                        }
                    }
                }
            }
        }
    }
});

// TLS/SSL configuration
net:
  tls:
    mode: requireTLS
    certificateKeyFile: /path/to/server.pem
    CAFile: /path/to/ca.pem
    clusterFile: /path/to/cluster.pem
```

### Network Security

```javascript
// IP whitelist (Atlas)
// Configure in Atlas UI or API

// Bind to specific interface
net:
  bindIp: 127.0.0.1,192.168.1.100
  port: 27017

// Firewall rules (Linux iptables)
// iptables -A INPUT -p tcp --dport 27017 -s 192.168.1.0/24 -j ACCEPT
// iptables -A INPUT -p tcp --dport 27017 -j DROP

// VPC peering (Atlas)
// Configure in Atlas Network Access settings
```

---

## 16. Multi-Tenancy Design

### Theory: Multi-Tenant Architecture Patterns

**What is Multi-Tenancy?**
Multi-tenancy is an architecture where a single instance of software serves multiple tenants (customers/organizations). Each tenant's data is isolated and appears as if they have a dedicated system.

**The Isolation vs Efficiency Trade-off**
Multi-tenancy involves balancing two competing concerns:

```
Full Isolation ◄────────────────────────────────────────► Maximum Efficiency
(dedicated per      (shared everything
 tenant - costly)   - noisy neighbor risk)
```

**The Three Multi-Tenancy Models**

| Model | Isolation Level | Operational Cost | Best For |
|-------|----------------|------------------|----------|
| **Database-per-tenant** | Complete | Highest | Compliance-heavy, large tenants |
| **Collection-per-tenant** | High | Medium | Medium tenants, some isolation needed |
| **Shared collection** | Logical | Lowest | Many small tenants, cost-sensitive |

**Key Challenges in Multi-Tenant Systems**

1. **Noisy Neighbor Problem**: One tenant's heavy usage affects others
   - Solutions: Resource quotas, rate limiting, tenant-aware sharding

2. **Data Security**: Ensuring complete tenant isolation
   - Solutions: Always include tenantId in queries, row-level security

3. **Schema Evolution**: Different tenants may need different schemas
   - Solutions: Schema versioning, feature flags per tenant

4. **Tenant Lifecycle**: Creating, suspending, and deleting tenants
   - Solutions: Automated provisioning, clear data deletion procedures

5. **Compliance**: Different tenants have different regulatory requirements
   - Solutions: Zone sharding for data residency, audit logging

**The TenantId Anti-Pattern**
A common mistake is forgetting to include tenantId in queries:

```javascript
// DANGEROUS: Query without tenant filter
db.documents.find({ category: "finance" })  // Returns ALL tenants' data!

// CORRECT: Always include tenantId
db.documents.find({ tenantId: currentTenant, category: "finance" })
```

**Indexing Strategy for Multi-Tenancy**
Compound indexes MUST include tenantId as the first field:
- `{ tenantId: 1, category: 1 }` ✓ Efficient
- `{ category: 1, tenantId: 1 }` ✗ Full scan for tenant-specific queries

### Database-per-Tenant

```javascript
// Each tenant has own database
// Pros: Complete isolation, easy to delete tenant, compliance-friendly
// Cons: More databases to manage, connection pool per tenant

async function getTenantDatabase(tenantId) {
    const dbName = `tenant_${tenantId}`;
    return client.db(dbName);
}

// Create tenant
async function createTenant(tenantId, tenantName) {
    const db = client.db(`tenant_${tenantId}`);
    
    // Create collections with schemas
    await db.createCollection('users');
    await db.createCollection('orders');
    
    // Create indexes
    await db.collection('users').createIndex({ email: 1 }, { unique: true });
    await db.collection('orders').createIndex({ userId: 1, createdAt: -1 });
    
    // Store tenant metadata
    await client.db('platform').collection('tenants').insertOne({
        _id: tenantId,
        name: tenantName,
        database: `tenant_${tenantId}`,
        createdAt: new Date()
    });
}

// Delete tenant (complete data removal)
async function deleteTenant(tenantId) {
    await client.db(`tenant_${tenantId}`).dropDatabase();
    await client.db('platform').collection('tenants').deleteOne({ _id: tenantId });
}
```

### Collection-per-Tenant

```javascript
// Each tenant has own collections within shared database
// Pros: Easier management, shared connection pool
// Cons: Namespace management, harder to isolate

async function getTenantCollection(tenantId, collectionName) {
    const db = client.db('myapp');
    return db.collection(`${tenantId}_${collectionName}`);
}

// Create tenant collections
async function createTenant(tenantId) {
    const db = client.db('myapp');
    
    await db.createCollection(`${tenantId}_users`);
    await db.createCollection(`${tenantId}_orders`);
    
    // Create indexes per tenant collection
    await db.collection(`${tenantId}_users`).createIndex({ email: 1 });
    await db.collection(`${tenantId}_orders`).createIndex({ userId: 1 });
}
```

### Shared Collection with Tenant ID

```javascript
// All tenants share collections, filtered by tenantId
// Pros: Simple, efficient for many small tenants
// Cons: Noisy neighbor risk, requires careful indexing

// Document structure
{
    _id: ObjectId(),
    tenantId: "tenant123",  // Discriminator field
    email: "user@example.com",
    name: "John Doe"
}

// CRITICAL: Include tenantId in all indexes
db.users.createIndex({ tenantId: 1, email: 1 }, { unique: true })
db.users.createIndex({ tenantId: 1, createdAt: -1 })
db.orders.createIndex({ tenantId: 1, userId: 1, orderDate: -1 })

// Always filter by tenantId
async function getUsers(tenantId, filter = {}) {
    return db.users.find({
        tenantId: tenantId,  // Always include!
        ...filter
    }).toArray();
}

// Middleware to inject tenantId
function tenantMiddleware(req, res, next) {
    const tenantId = req.headers['x-tenant-id'];
    if (!tenantId) {
        return res.status(400).json({ error: 'Tenant ID required' });
    }
    req.tenantId = tenantId;
    next();
}

// Use in queries
app.get('/api/users', tenantMiddleware, async (req, res) => {
    const users = await db.users.find({ tenantId: req.tenantId }).toArray();
    res.json(users);
});
```

### Hybrid Multi-Tenancy

```javascript
// Combine approaches based on tenant size/requirements

class TenantManager {
    async getTenantDb(tenantId) {
        const tenant = await this.getTenantConfig(tenantId);
        
        switch (tenant.tier) {
            case 'enterprise':
                // Dedicated database
                return client.db(`enterprise_${tenantId}`);
                
            case 'business':
                // Shared database, dedicated collections
                return {
                    users: client.db('business').collection(`${tenantId}_users`),
                    orders: client.db('business').collection(`${tenantId}_orders`)
                };
                
            default:
                // Shared everything
                return {
                    users: client.db('shared').collection('users'),
                    orders: client.db('shared').collection('orders'),
                    filter: { tenantId }  // Must include in all queries
                };
        }
    }
}

// Shard key for multi-tenant
// Include tenantId in shard key for isolation
sh.shardCollection(
    "myapp.orders",
    { tenantId: 1, orderId: "hashed" }
)

// Zone sharding for tenant isolation
sh.addShardTag("shard0001", "enterprise_tenant1")
sh.updateZoneKeyRange(
    "myapp.orders",
    { tenantId: "tenant1", orderId: MinKey },
    { tenantId: "tenant1", orderId: MaxKey },
    "enterprise_tenant1"
)
```

---

## 17. Time-Series Data Design

### Theory: Time-Series Database Concepts

**What is Time-Series Data?**
Time-series data is a sequence of data points collected at successive, equally spaced points in time. Common examples:
- IoT sensor readings
- Application metrics
- Financial market data
- Log events
- User activity tracking

**Characteristics of Time-Series Workloads**

| Characteristic | Time-Series | General Purpose |
|---------------|-------------|------------------|
| Write Pattern | Append-only, high volume | CRUD operations |
| Update Pattern | Rare to never | Frequent |
| Query Pattern | Range scans by time | Point lookups |
| Data Lifecycle | Often time-bounded | Retained indefinitely |
| Aggregation | Common (avg, sum, min, max) | Variable |

**Why Traditional Document Storage Struggles**

Storing one document per measurement creates problems:
- **Index bloat**: Billions of index entries
- **Storage overhead**: BSON overhead per document significant at scale
- **Query inefficiency**: Many small documents scattered across storage
- **Memory pressure**: Each document has metadata overhead

**The Bucketing Solution**
Grouping multiple measurements into buckets addresses these issues:

```
Without bucketing: 1 document per measurement = 86,400 docs/day/sensor
With bucketing:    1 document per hour = 24 docs/day/sensor
Improvement:       3,600x fewer documents
```

**Time-Series Storage Optimizations**
MongoDB 5.0+ time-series collections implement optimizations:

1. **Automatic bucketing**: Documents grouped by time range
2. **Compression**: Delta encoding for timestamps, Gorilla for floats
3. **Columnar storage**: Better compression for homogeneous data
4. **Sorted storage**: Data physically ordered by time
5. **Metadata extraction**: Common fields stored once per bucket

**Downsampling Strategy**
Retain high resolution for recent data, aggregate older data:

```
Raw data (1 sec) ──► Minute aggregates ──► Hourly ──► Daily
   (7 days)           (30 days)           (1 year)   (forever)
```

This balances query granularity against storage costs.

**Query Patterns for Time-Series**
- **Time range**: Most common - "last hour", "yesterday"
- **Aggregation**: Sum, average, min, max over time windows
- **Downsampling**: Reduce resolution for visualization
- **Anomaly detection**: Compare current to historical

### Time-Series Collection (MongoDB 5.0+)

```javascript
// Native time-series collection
db.createCollection("metrics", {
    timeseries: {
        timeField: "timestamp",
        metaField: "metadata",
        granularity: "seconds"  // seconds, minutes, hours
    },
    expireAfterSeconds: 2592000  // 30 days TTL
})

// Document structure
{
    timestamp: ISODate("2024-03-15T10:30:00Z"),
    metadata: {
        sensorId: "sensor001",
        location: "datacenter-1",
        type: "temperature"
    },
    value: 23.5,
    unit: "celsius"
}

// Automatic bucketing handled by MongoDB
// Optimized storage and queries

// Query time-series data
db.metrics.aggregate([
    {
        $match: {
            timestamp: {
                $gte: ISODate("2024-03-15T00:00:00Z"),
                $lt: ISODate("2024-03-16T00:00:00Z")
            },
            "metadata.sensorId": "sensor001"
        }
    },
    {
        $group: {
            _id: {
                $dateTrunc: {
                    date: "$timestamp",
                    unit: "hour"
                }
            },
            avgValue: { $avg: "$value" },
            maxValue: { $max: "$value" },
            minValue: { $min: "$value" },
            count: { $sum: 1 }
        }
    },
    { $sort: { _id: 1 } }
])
```

### Manual Bucket Pattern

```javascript
// For MongoDB < 5.0 or custom bucketing needs
{
    _id: ObjectId(),
    sensorId: "sensor001",
    bucketStart: ISODate("2024-03-15T10:00:00Z"),
    bucketEnd: ISODate("2024-03-15T11:00:00Z"),
    granularity: "hour",
    measurements: [
        { ts: ISODate("2024-03-15T10:00:00Z"), value: 23.5 },
        { ts: ISODate("2024-03-15T10:01:00Z"), value: 23.6 },
        // ... more measurements
    ],
    stats: {
        count: 60,
        sum: 1416.0,
        avg: 23.6,
        min: 23.1,
        max: 24.2
    }
}

// Upsert with bucket management
async function insertMeasurement(sensorId, timestamp, value) {
    const bucketStart = new Date(timestamp);
    bucketStart.setMinutes(0, 0, 0);
    
    const bucketEnd = new Date(bucketStart);
    bucketEnd.setHours(bucketEnd.getHours() + 1);
    
    await db.sensorData.updateOne(
        {
            sensorId,
            bucketStart,
            "stats.count": { $lt: 60 }  // Max measurements per bucket
        },
        {
            $push: { measurements: { ts: timestamp, value } },
            $inc: { "stats.count": 1, "stats.sum": value },
            $min: { "stats.min": value },
            $max: { "stats.max": value },
            $setOnInsert: {
                sensorId,
                bucketStart,
                bucketEnd,
                granularity: "hour"
            }
        },
        { upsert: true }
    );
    
    // Update average separately
    await db.sensorData.updateOne(
        { sensorId, bucketStart },
        [
            {
                $set: {
                    "stats.avg": { $divide: ["$stats.sum", "$stats.count"] }
                }
            }
        ]
    );
}
```

### Downsampling Strategy

```javascript
// Keep high-resolution data for recent period
// Downsample older data for storage efficiency

// Raw data: Keep 7 days at 1-second granularity
// Minute aggregates: Keep 30 days
// Hourly aggregates: Keep 1 year
// Daily aggregates: Keep forever

async function downsampleData() {
    const now = new Date();
    
    // Create minute aggregates from raw data older than 7 days
    const weekAgo = new Date(now - 7 * 24 * 60 * 60 * 1000);
    
    await db.metrics_raw.aggregate([
        {
            $match: {
                timestamp: { $lt: weekAgo }
            }
        },
        {
            $group: {
                _id: {
                    sensorId: "$metadata.sensorId",
                    minute: {
                        $dateTrunc: { date: "$timestamp", unit: "minute" }
                    }
                },
                avgValue: { $avg: "$value" },
                minValue: { $min: "$value" },
                maxValue: { $max: "$value" },
                count: { $sum: 1 }
            }
        },
        {
            $merge: {
                into: "metrics_minute",
                on: "_id",
                whenMatched: "replace",
                whenNotMatched: "insert"
            }
        }
    ]).toArray();
    
    // Delete raw data after downsampling
    await db.metrics_raw.deleteMany({
        timestamp: { $lt: weekAgo }
    });
}

// Schedule downsampling job
// Run daily via cron or MongoDB scheduled triggers
```

---

## 18. Event-Driven Architecture

### Theory: Event-Driven Design Principles

**What is Event-Driven Architecture?**
Event-Driven Architecture (EDA) is a design paradigm where system behavior is determined by events - significant changes in state. Components communicate by producing and consuming events rather than direct method calls.

**Traditional vs Event-Driven**

| Aspect | Request-Response | Event-Driven |
|--------|-----------------|---------------|
| Coupling | Tight (direct calls) | Loose (events) |
| Scalability | Limited by slowest component | Independent scaling |
| Failure Handling | Cascading failures | Isolated failures |
| Real-time | Polling required | Push-based |

**Core Event-Driven Patterns**

1. **Event Notification**: Notify that something happened (minimal data)
   - Example: "Order 123 was created"

2. **Event-Carried State Transfer**: Include full state in event
   - Example: "Order 123 was created with these items..."

3. **Event Sourcing**: Store events as the source of truth
   - Current state = replaying all events

4. **CQRS**: Separate read and write models
   - Write: Events captured
   - Read: Materialized views optimized for queries

**MongoDB Change Streams: The Foundation**
Change streams provide real-time notification of data changes:

- Subscribe to insert, update, replace, delete, invalidate events
- Resume from a token after failure (exactly-once processing)
- Filter events on server side (reducing network traffic)
- Available at collection, database, or deployment level

**The Outbox Pattern: Why It Matters**
In distributed systems, updating database AND publishing event must be atomic:

```
// PROBLEM: Two separate operations
db.orders.insert(order);     // Succeeds
messageQueue.publish(event); // Fails - inconsistent state!

// SOLUTION: Outbox pattern
transaction {
    db.orders.insert(order);
    db.outbox.insert(event);  // Same transaction
}
// Later: Process outbox and publish
```

**Event Sourcing Benefits and Costs**

| Benefits | Costs |
|----------|-------|
| Complete audit trail | More complex queries |
| Temporal queries ("what was state at X?") | Storage growth |
| Easy debugging (replay events) | Learning curve |
| Supports event replay | Event versioning challenges |

**Eventual Consistency**
Event-driven systems embrace eventual consistency:
- Updates propagate asynchronously
- Different views may be temporarily inconsistent
- Design UI to handle stale data gracefully
- Use version numbers for conflict detection

### Change Streams

```javascript
// Watch collection for changes
const changeStream = db.orders.watch([
    { $match: { "fullDocument.status": "pending" } }
]);

changeStream.on('change', async (change) => {
    console.log('Change detected:', change.operationType);
    
    switch (change.operationType) {
        case 'insert':
            await processNewOrder(change.fullDocument);
            break;
        case 'update':
            await processOrderUpdate(change.documentKey._id, change.updateDescription);
            break;
        case 'delete':
            await processOrderDeletion(change.documentKey._id);
            break;
    }
});

// Watch with resume token for fault tolerance
let resumeToken = await loadResumeToken();  // From persistent storage

const changeStream = db.orders.watch([], {
    resumeAfter: resumeToken,
    fullDocument: 'updateLookup'  // Include full document for updates
});

changeStream.on('change', async (change) => {
    // Process change
    await processChange(change);
    
    // Store resume token for recovery
    await saveResumeToken(change._id);
});

// Watch entire database or cluster
const dbChangeStream = db.watch();
const clusterChangeStream = client.watch();
```

### Event Sourcing Pattern

```javascript
// Events collection stores all state changes
{
    _id: ObjectId(),
    aggregateId: "order-123",
    aggregateType: "Order",
    eventType: "OrderCreated",
    eventData: {
        customerId: "cust-456",
        items: [
            { productId: "prod-001", quantity: 2, price: 49.99 }
        ]
    },
    metadata: {
        timestamp: ISODate("2024-03-15T10:00:00Z"),
        version: 1,
        userId: "user-789",
        correlationId: "req-abc123"
    }
}

// Event store operations
class EventStore {
    async appendEvent(aggregateId, eventType, eventData, metadata) {
        // Get current version
        const lastEvent = await db.events
            .find({ aggregateId })
            .sort({ "metadata.version": -1 })
            .limit(1)
            .toArray();
        
        const version = lastEvent.length > 0 
            ? lastEvent[0].metadata.version + 1 
            : 1;
        
        // Append event with optimistic concurrency
        const result = await db.events.insertOne({
            aggregateId,
            aggregateType: this.getAggregateType(eventType),
            eventType,
            eventData,
            metadata: {
                ...metadata,
                timestamp: new Date(),
                version
            }
        });
        
        return result;
    }
    
    async getEvents(aggregateId, fromVersion = 0) {
        return db.events
            .find({
                aggregateId,
                "metadata.version": { $gt: fromVersion }
            })
            .sort({ "metadata.version": 1 })
            .toArray();
    }
    
    async rebuildAggregate(aggregateId) {
        const events = await this.getEvents(aggregateId);
        let state = {};
        
        for (const event of events) {
            state = this.applyEvent(state, event);
        }
        
        return state;
    }
}

// Index for event store
db.events.createIndex({ aggregateId: 1, "metadata.version": 1 }, { unique: true })
db.events.createIndex({ "metadata.timestamp": 1 })
db.events.createIndex({ eventType: 1 })
```

### Outbox Pattern

```javascript
// Ensures atomic writes with reliable event publishing
// Useful for microservices communication

// Transaction writes both data and outbox event
const session = client.startSession();
try {
    session.startTransaction();
    
    // Write business data
    const order = {
        _id: new ObjectId(),
        customerId: "cust-123",
        items: [...],
        status: "created",
        createdAt: new Date()
    };
    await db.orders.insertOne(order, { session });
    
    // Write outbox event (same transaction)
    await db.outbox.insertOne({
        aggregateType: "Order",
        aggregateId: order._id.toString(),
        eventType: "OrderCreated",
        payload: order,
        createdAt: new Date(),
        published: false
    }, { session });
    
    await session.commitTransaction();
} catch (error) {
    await session.abortTransaction();
    throw error;
} finally {
    session.endSession();
}

// Outbox processor (separate service)
async function processOutbox() {
    while (true) {
        const events = await db.outbox
            .find({ published: false })
            .sort({ createdAt: 1 })
            .limit(100)
            .toArray();
        
        for (const event of events) {
            try {
                // Publish to message queue
                await messageQueue.publish(event.eventType, event.payload);
                
                // Mark as published
                await db.outbox.updateOne(
                    { _id: event._id },
                    { 
                        $set: { 
                            published: true, 
                            publishedAt: new Date() 
                        }
                    }
                );
            } catch (error) {
                console.error('Failed to publish event:', error);
            }
        }
        
        await sleep(1000);  // Poll interval
    }
}

// Cleanup old events
db.outbox.createIndex(
    { publishedAt: 1 },
    { expireAfterSeconds: 604800 }  // 7 days TTL
)
```

---

## 19. Migration & Evolution Strategies

### Theory: Schema Evolution in Document Databases

**The Schema Evolution Challenge**
Unlike relational databases with rigid schemas, document databases are "schema-flexible." This flexibility is both a feature and a responsibility:

- **Advantage**: Easier to evolve schemas without downtime
- **Challenge**: Multiple schema versions may coexist in production

**Types of Schema Changes**

| Change Type | Complexity | Example |
|------------|------------|----------|
| **Additive** | Low | Adding a new field |
| **Rename** | Medium | Changing field name |
| **Restructure** | High | Splitting or combining fields |
| **Type Change** | High | String to Number |
| **Removal** | Medium | Deleting obsolete fields |

**Migration Strategies Compared**

1. **Big Bang Migration**
   - Migrate all documents at once
   - Requires downtime or heavy locking
   - Use for: Small datasets, infrequent changes

2. **Lazy Migration**
   - Migrate documents when accessed
   - No downtime, gradual transition
   - Use for: Large datasets, background migration

3. **Dual-Write Migration**
   - Write to both old and new format during transition
   - Zero downtime, controlled rollout
   - Use for: Critical systems, risk-averse environments

**The Version Field Pattern**
Track schema version in each document:

```javascript
{
  schemaVersion: 3,
  // ... fields in v3 format
}
```

Benefits:
- Application code handles multiple versions
- Clear audit of document state
- Enable targeted migration queries

**Backward Compatibility Principles**

1. **Be liberal in what you accept**: Handle old formats gracefully
2. **Be conservative in what you produce**: Always write latest format
3. **Never remove fields immediately**: Deprecate first, remove later
4. **Default values**: New required fields should have sensible defaults

**Migration Best Practices**

- **Test migrations** on production data copies
- **Use checkpoints** to resume failed migrations
- **Throttle operations** to avoid impacting production
- **Monitor performance** during migration
- **Have rollback plans** for each migration step
- **Validate data** after migration completes

**Index Considerations During Migration**
Schema changes often require index changes:
- Build new indexes before deploying new code
- Remove old indexes after migration completes
- Use `hideIndex()` to test impact before dropping

### Schema Migration Strategies

```javascript
// 1. Lazy Migration (on read/write)
async function getUser(userId) {
    let user = await db.users.findOne({ _id: userId });
    
    if (user.schemaVersion < CURRENT_SCHEMA_VERSION) {
        user = migrateUserDocument(user);
        await db.users.updateOne(
            { _id: userId },
            { $set: user }
        );
    }
    
    return user;
}

function migrateUserDocument(doc) {
    let migrated = { ...doc };
    
    // Version 1 → 2: Split name into firstName/lastName
    if (doc.schemaVersion === 1) {
        const [firstName, ...rest] = doc.name.split(' ');
        migrated.firstName = firstName;
        migrated.lastName = rest.join(' ');
        delete migrated.name;
        migrated.schemaVersion = 2;
    }
    
    // Version 2 → 3: Add default preferences
    if (migrated.schemaVersion === 2) {
        migrated.preferences = {
            notifications: true,
            theme: 'light'
        };
        migrated.schemaVersion = 3;
    }
    
    return migrated;
}

// 2. Batch Migration (background job)
async function batchMigrate(batchSize = 1000) {
    let processed = 0;
    
    while (true) {
        const oldDocs = await db.users
            .find({ schemaVersion: { $lt: CURRENT_SCHEMA_VERSION } })
            .limit(batchSize)
            .toArray();
        
        if (oldDocs.length === 0) break;
        
        const bulkOps = oldDocs.map(doc => ({
            updateOne: {
                filter: { _id: doc._id },
                update: { $set: migrateUserDocument(doc) }
            }
        }));
        
        await db.users.bulkWrite(bulkOps);
        processed += oldDocs.length;
        
        console.log(`Migrated ${processed} documents`);
        await sleep(100);  // Throttle to reduce load
    }
    
    return processed;
}

// 3. Dual-Write Migration (zero downtime)
// Phase 1: Write to both old and new
// Phase 2: Migrate existing data
// Phase 3: Read from new
// Phase 4: Remove old
```

### Data Migration Scripts

```javascript
// Safe migration with checkpoints
async function migrateWithCheckpoint(collectionName, migrateFn) {
    const checkpointKey = `migration_${collectionName}_${Date.now()}`;
    
    // Get or create checkpoint
    let checkpoint = await db.migrations.findOne({ _id: checkpointKey });
    if (!checkpoint) {
        checkpoint = {
            _id: checkpointKey,
            lastProcessedId: null,
            processedCount: 0,
            startedAt: new Date(),
            status: 'running'
        };
        await db.migrations.insertOne(checkpoint);
    }
    
    const BATCH_SIZE = 1000;
    let query = {};
    
    if (checkpoint.lastProcessedId) {
        query._id = { $gt: checkpoint.lastProcessedId };
    }
    
    while (true) {
        const docs = await db[collectionName]
            .find(query)
            .sort({ _id: 1 })
            .limit(BATCH_SIZE)
            .toArray();
        
        if (docs.length === 0) break;
        
        // Process batch
        for (const doc of docs) {
            await migrateFn(doc);
        }
        
        // Update checkpoint
        const lastId = docs[docs.length - 1]._id;
        await db.migrations.updateOne(
            { _id: checkpointKey },
            {
                $set: { lastProcessedId: lastId },
                $inc: { processedCount: docs.length }
            }
        );
        
        query._id = { $gt: lastId };
    }
    
    // Mark complete
    await db.migrations.updateOne(
        { _id: checkpointKey },
        { $set: { status: 'completed', completedAt: new Date() } }
    );
}

// Usage
await migrateWithCheckpoint('users', async (doc) => {
    const migrated = migrateUserDocument(doc);
    await db.users.updateOne(
        { _id: doc._id },
        { $set: migrated }
    );
});
```

### Index Management During Migration

```javascript
// Create index in background without blocking
db.users.createIndex(
    { newField: 1 },
    { 
        background: true,  // Deprecated in 4.2+, now default
        name: "idx_newField"
    }
);

// Monitor index build progress
db.currentOp({ "command.createIndexes": { $exists: true } })

// Rolling index build for replica sets
// 1. Build on secondaries first
// 2. Step down primary
// 3. Build on old primary

// Drop unused indexes
const indexUsage = await db.users.aggregate([
    { $indexStats: {} }
]).toArray();

for (const idx of indexUsage) {
    if (idx.accesses.ops === 0 && idx.name !== '_id_') {
        console.log(`Unused index: ${idx.name}`);
        // db.users.dropIndex(idx.name);  // Uncomment to drop
    }
}
```

---

## 20. Real-World Case Studies

### Theory: Applying Patterns in Production

**Why Case Studies Matter**
Theoretical knowledge needs practical context. Real-world case studies demonstrate:
- How multiple patterns combine in practice
- Trade-offs made under real constraints
- Performance characteristics at scale
- Lessons learned from production incidents

**Common Architectural Patterns**

1. **E-Commerce Systems**
   - High write volume (orders, inventory)
   - Complex queries (search, recommendations)
   - Global distribution requirements
   - Patterns: Extended Reference, Computed, Zone Sharding

2. **IoT/Telemetry Platforms**
   - Massive ingestion rates (millions of events/second)
   - Time-range queries dominant
   - Long-term storage requirements
   - Patterns: Bucket, Time-Series Collections, Downsampling

3. **SaaS Multi-Tenant Applications**
   - Thousands of tenants with varying sizes
   - Strict data isolation requirements
   - Cost efficiency critical
   - Patterns: Tenant ID discrimination, Hybrid Multi-Tenancy

4. **Social/Content Platforms**
   - Highly variable document sizes
   - Fan-out on read and write
   - Real-time engagement metrics
   - Patterns: Outlier, Subset, Computed, Event-Driven

**Scaling Decision Framework**

When facing scale challenges, consider in order:

1. **Index Optimization**: Often 10-100x improvement potential
2. **Schema Redesign**: Apply appropriate patterns
3. **Read Scaling**: Add secondaries, read preference tuning
4. **Sharding**: Horizontal scale for writes and data volume
5. **Caching Layer**: External cache for hot data

**Performance Benchmarking Principles**

- Test with **production-like data volumes**
- Include **realistic query patterns**
- Measure at **expected scale** (not just current)
- Test **failure scenarios** (node loss, network partition)
- Benchmark **write and read paths separately**

**Common Pitfalls to Avoid**

| Pitfall | Consequence | Prevention |
|---------|-------------|------------|
| Wrong shard key | Hot spots, poor distribution | Analyze query patterns before choosing |
| Unbounded arrays | Document size limit, slow updates | Use bucketing or referencing |
| Missing indexes | Full collection scans | Monitor slow query log |
| Over-sharding | Operational complexity | Start with fewer shards, add as needed |
| Ignoring working set | Performance degradation | Monitor cache hit rates |

### Case Study 1: E-Commerce Platform (100M Orders/Year)

```javascript
// Challenge: High write volume, complex queries, global distribution

// Solution Architecture:

// 1. Shard Key Selection
sh.shardCollection(
    "ecommerce.orders",
    { customerId: "hashed" }  // Even distribution by customer
)

// 2. Document Design
{
    _id: ObjectId(),
    orderId: "ORD-2024-001",
    customerId: "cust-123",  // Shard key
    
    // Extended reference for common queries
    customer: {
        _id: "cust-123",
        name: "John Doe",
        email: "john@example.com"
    },
    
    // Embedded items (always accessed together)
    items: [
        {
            productId: "prod-001",
            sku: "SKU-001",
            name: "Product Name",
            quantity: 2,
            unitPrice: 49.99,
            subtotal: 99.98
        }
    ],
    
    // Pre-computed totals
    totals: {
        itemCount: 2,
        subtotal: 99.98,
        tax: 8.00,
        shipping: 5.00,
        total: 112.98
    },
    
    // Status tracking with timestamps
    status: "delivered",
    statusHistory: [
        { status: "pending", timestamp: ISODate("..."), actor: "system" },
        { status: "processing", timestamp: ISODate("..."), actor: "user-123" },
        { status: "shipped", timestamp: ISODate("..."), actor: "user-456" },
        { status: "delivered", timestamp: ISODate("..."), actor: "system" }
    ],
    
    createdAt: ISODate("2024-03-15T10:00:00Z"),
    updatedAt: ISODate("2024-03-16T14:30:00Z")
}

// 3. Indexes
db.orders.createIndex({ customerId: 1, createdAt: -1 })  // Customer order history
db.orders.createIndex({ status: 1, createdAt: -1 })       // Status filtering
db.orders.createIndex({ "items.productId": 1 })           // Product lookup
db.orders.createIndex({ createdAt: -1 })                  // Time-based queries

// 4. Read/Write Split
// Writes: Primary only with majority write concern
// Reads: secondaryPreferred for dashboards/reports

// 5. Results
// - 50,000 orders/hour at peak
// - 10ms average query latency
// - 99.99% availability
```

### Case Study 2: IoT Sensor Platform (1B Events/Day)

```javascript
// Challenge: Massive time-series data, real-time analytics

// Solution Architecture:

// 1. Time-Series Collection
db.createCollection("sensor_events", {
    timeseries: {
        timeField: "timestamp",
        metaField: "sensor",
        granularity: "seconds"
    },
    expireAfterSeconds: 2592000  // 30 days
})

// 2. Document Design
{
    timestamp: ISODate("2024-03-15T10:30:00Z"),
    sensor: {
        id: "sensor-001",
        type: "temperature",
        location: "building-a-floor-3",
        customer: "customer-xyz"
    },
    value: 23.5,
    unit: "celsius",
    quality: "good"
}

// 3. Downsampling Pipeline
// Raw data → Minute aggregates → Hour aggregates → Day aggregates

// 4. Sharding by sensor and time
sh.shardCollection(
    "iot.sensor_events",
    { "sensor.customer": 1, timestamp: 1 }
)

// 5. Zone sharding for data locality
sh.addShardTag("shard-us", "US")
sh.addShardTag("shard-eu", "EU")

// 6. Real-time Aggregation
db.sensor_events.aggregate([
    {
        $match: {
            "sensor.id": "sensor-001",
            timestamp: { $gte: new Date(Date.now() - 3600000) }
        }
    },
    {
        $group: {
            _id: {
                $dateTrunc: { date: "$timestamp", unit: "minute" }
            },
            avgValue: { $avg: "$value" },
            maxValue: { $max: "$value" },
            minValue: { $min: "$value" }
        }
    },
    { $sort: { _id: -1 } }
])

// 7. Results
// - 10,000 events/second write throughput
// - Sub-second query latency for recent data
// - 70% storage reduction with time-series collections
```

### Case Study 3: Multi-Tenant SaaS Application

```javascript
// Challenge: 10,000 tenants, varying sizes, data isolation

// Solution Architecture:

// 1. Hybrid Multi-Tenancy
// - Small tenants: Shared collection with tenantId
// - Large tenants: Dedicated collections
// - Enterprise tenants: Dedicated databases

// 2. Tenant Configuration
{
    _id: "tenant-123",
    name: "Acme Corp",
    tier: "business",  // free, business, enterprise
    isolation: "shared",  // shared, collection, database
    config: {
        maxUsers: 100,
        maxStorage: "10GB",
        features: ["analytics", "api"]
    },
    shard: "shard-us-east-1",  // Assigned shard
    createdAt: ISODate("2024-01-01")
}

// 3. Shard Key Design
sh.shardCollection(
    "saas.documents",
    { tenantId: 1, _id: "hashed" }
)

// 4. Application-Level Routing
class TenantRouter {
    async getCollection(tenantId, collectionName) {
        const tenant = await this.getTenant(tenantId);
        
        switch (tenant.isolation) {
            case 'database':
                return client.db(`tenant_${tenantId}`).collection(collectionName);
            case 'collection':
                return client.db('saas').collection(`${tenantId}_${collectionName}`);
            default:
                return {
                    collection: client.db('saas').collection(collectionName),
                    filter: { tenantId }
                };
        }
    }
}

// 5. Resource Limits
async function checkTenantLimits(tenantId) {
    const tenant = await getTenant(tenantId);
    const usage = await db.documents.aggregate([
        { $match: { tenantId } },
        { $group: {
            _id: null,
            count: { $sum: 1 },
            size: { $sum: { $bsonSize: "$$ROOT" } }
        }}
    ]).toArray();
    
    return {
        withinLimits: usage[0].size < parseSize(tenant.config.maxStorage),
        currentUsage: usage[0]
    };
}

// 6. Results
// - Linear scaling with tenant count
// - Sub-10ms queries with tenant isolation
// - Zero cross-tenant data leakage
```

---

## Summary: Key Takeaways

### Data Modeling Checklist

```
□ Choose embedding vs referencing based on access patterns
□ Consider document size limits (16MB max, <1MB recommended)
□ Apply appropriate schema design patterns
□ Plan for schema evolution with versioning
□ Avoid unbounded arrays
```

### Performance Checklist

```
□ Design indexes following ESR rule
□ Include shard key in all queries (sharded clusters)
□ Use covered queries where possible
□ Implement pagination with keyset pattern
□ Monitor and tune WiredTiger cache
□ Set appropriate read/write concerns
```

### Scalability Checklist

```
□ Choose shard key based on query patterns and cardinality
□ Pre-split chunks before large data loads
□ Implement connection pooling
□ Use read replicas for read scaling
□ Plan capacity based on working set size
```

### Reliability Checklist

```
□ Use replica sets for high availability
□ Configure appropriate write concerns
□ Implement retry logic for transient failures
□ Monitor replication lag
□ Test failover procedures regularly
□ Implement backup and recovery procedures
```

---

## Resources

### Official Documentation
- [MongoDB Manual](https://docs.mongodb.com/manual/)
- [MongoDB University](https://university.mongodb.com/)
- [Performance Best Practices](https://www.mongodb.com/docs/manual/administration/analyzing-mongodb-performance/)

### Tools
- **MongoDB Compass** - GUI for exploration and optimization
- **MongoDB Atlas** - Managed cloud service
- **mongosh** - Modern MongoDB shell
- **Percona PMM** - Monitoring and management

### Community
- [MongoDB Community Forums](https://www.mongodb.com/community/forums/)
- [MongoDB Blog](https://www.mongodb.com/blog)
- [Stack Overflow - MongoDB](https://stackoverflow.com/questions/tagged/mongodb)

---

*Last Updated: February 2026*
