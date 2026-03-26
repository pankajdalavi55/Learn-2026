# Database Selection and Design Guide
## Comprehensive Guide to Database Types, Selection Criteria, and Design Considerations

---

## Table of Contents

1. [Introduction to Database Systems](#1-introduction-to-database-systems)
2. [Database Categories and Types](#2-database-categories-and-types)
3. [Relational Databases (RDBMS)](#3-relational-databases-rdbms)
4. [NoSQL Databases](#4-nosql-databases)
5. [NewSQL and Distributed SQL](#5-newsql-and-distributed-sql)
6. [Specialized Databases](#6-specialized-databases)
7. [ACID vs BASE: Consistency Models](#7-acid-vs-base-consistency-models)
8. [CAP Theorem and Trade-offs](#8-cap-theorem-and-trade-offs)
9. [Database Selection Criteria](#9-database-selection-criteria)
10. [Data Modeling Approaches](#10-data-modeling-approaches)
11. [Scaling Strategies](#11-scaling-strategies)
12. [Performance Considerations](#12-performance-considerations)
13. [High Availability and Disaster Recovery](#13-high-availability-and-disaster-recovery)
14. [Security Considerations](#14-security-considerations)
15. [Cost Analysis](#15-cost-analysis)
16. [Migration Strategies](#16-migration-strategies)
17. [Polyglot Persistence](#17-polyglot-persistence)
18. [Cloud vs On-Premise Databases](#18-cloud-vs-on-premise-databases)
19. [Database Selection Decision Framework](#19-database-selection-decision-framework)
20. [Use Case Scenarios and Recommendations](#20-use-case-scenarios-and-recommendations)

---

## 1. Introduction to Database Systems

### Theory: What is a Database?

A **database** is an organized collection of structured data stored electronically. A **Database Management System (DBMS)** is the software that interacts with end users, applications, and the database itself to capture, store, and analyze data.

**Why Databases Matter**
- **Data Persistence**: Store data beyond application lifecycle
- **Data Integrity**: Ensure accuracy and consistency
- **Concurrent Access**: Multiple users/applications access data simultaneously
- **Data Security**: Control who can access what data
- **Query Capability**: Efficiently retrieve specific data
- **Scalability**: Handle growing data volumes

### Evolution of Database Systems

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATABASE EVOLUTION TIMELINE                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1960s          1970s          1980s-90s       2000s          2010s+        │
│    │              │               │              │              │            │
│    ▼              ▼               ▼              ▼              ▼            │
│ ┌──────┐     ┌──────┐       ┌──────────┐   ┌─────────┐   ┌───────────┐     │
│ │Flat  │     │Relat-│       │Object-   │   │ NoSQL   │   │ NewSQL/   │     │
│ │Files │────►│ional │──────►│Relational│──►│ DBs     │──►│ Cloud     │     │
│ │      │     │(SQL) │       │Databases │   │         │   │ Native    │     │
│ └──────┘     └──────┘       └──────────┘   └─────────┘   └───────────┘     │
│                                                                              │
│  Sequential   Structured    Complex Data   Scale &       Best of           │
│  Access       Queries       + Inheritance  Flexibility   Both Worlds       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Key Database Concepts

| Concept | Definition | Example |
|---------|------------|---------|
| **Schema** | Structure definition of data organization | Table definitions with columns and types |
| **Index** | Data structure improving query speed | B-tree index on customer_id |
| **Transaction** | Unit of work that is atomic | Transfer money between accounts |
| **Normalization** | Organizing data to reduce redundancy | 3NF database design |
| **Query** | Request for data from database | SELECT * FROM users WHERE active = true |
| **Replication** | Copying data to multiple nodes | Master-slave replication |
| **Sharding** | Horizontal partitioning across servers | User data split by region |

---

## 2. Database Categories and Types

### Theory: Database Taxonomy

Databases can be classified by multiple dimensions: data model, consistency model, deployment model, and use case optimization.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        DATABASE TAXONOMY                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  BY DATA MODEL:                                                              │
│  ├── Relational (SQL)         Tables with rows and columns                  │
│  ├── Document                 JSON/BSON documents                            │
│  ├── Key-Value               Simple key-value pairs                          │
│  ├── Wide-Column             Column families                                 │
│  ├── Graph                   Nodes and edges                                 │
│  ├── Time-Series             Timestamped data points                         │
│  └── Vector                  High-dimensional vectors                        │
│                                                                              │
│  BY CONSISTENCY MODEL:                                                       │
│  ├── ACID                    Strong consistency                              │
│  └── BASE                    Eventual consistency                            │
│                                                                              │
│  BY DEPLOYMENT:                                                              │
│  ├── Single-node             Traditional deployment                          │
│  ├── Distributed             Multiple nodes                                  │
│  └── Cloud-native            Designed for cloud                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Database Types Overview

| Type | Data Model | Best For | Examples |
|------|------------|----------|----------|
| **Relational** | Tables, rows, columns | Structured data, complex queries, ACID | PostgreSQL, MySQL, Oracle, SQL Server |
| **Document** | JSON/BSON documents | Semi-structured, flexible schema | MongoDB, CouchDB, Firestore |
| **Key-Value** | Key-value pairs | Caching, sessions, simple lookups | Redis, DynamoDB, Memcached |
| **Wide-Column** | Column families | Time-series, write-heavy workloads | Cassandra, HBase, ScyllaDB |
| **Graph** | Nodes and edges | Relationships, social networks | Neo4j, Amazon Neptune, JanusGraph |
| **Time-Series** | Timestamped data | IoT, metrics, monitoring | InfluxDB, TimescaleDB, Prometheus |
| **Vector** | High-dimensional vectors | AI/ML, similarity search | Pinecone, Milvus, Weaviate |
| **Search** | Inverted indexes | Full-text search, logging | Elasticsearch, Solr, OpenSearch |

### When to Use Which Type

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATABASE SELECTION QUICK GUIDE                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Need ACID transactions + complex queries?                                   │
│  └──► Relational (PostgreSQL, MySQL)                                        │
│                                                                              │
│  Need flexible schema + document-oriented?                                   │
│  └──► Document (MongoDB, CouchDB)                                           │
│                                                                              │
│  Need ultra-fast reads/writes for simple data?                              │
│  └──► Key-Value (Redis, DynamoDB)                                           │
│                                                                              │
│  Need high write throughput at scale?                                        │
│  └──► Wide-Column (Cassandra, ScyllaDB)                                     │
│                                                                              │
│  Need to traverse relationships efficiently?                                 │
│  └──► Graph (Neo4j, Neptune)                                                │
│                                                                              │
│  Need time-based data with aggregations?                                     │
│  └──► Time-Series (InfluxDB, TimescaleDB)                                   │
│                                                                              │
│  Need full-text search + analytics?                                          │
│  └──► Search (Elasticsearch)                                                │
│                                                                              │
│  Need similarity search for AI/ML?                                           │
│  └──► Vector (Pinecone, Milvus)                                             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Relational Databases (RDBMS)

### Theory: The Relational Model

The **relational model**, introduced by Edgar F. Codd in 1970, organizes data into **relations** (tables) consisting of **tuples** (rows) and **attributes** (columns). This model is based on set theory and first-order predicate logic.

**Core Principles:**
1. **Data Independence**: Physical storage is separate from logical structure
2. **Structural Independence**: Changes to schema don't require application changes
3. **Data Integrity**: Constraints ensure data validity
4. **Relational Algebra**: Mathematical foundation for queries

### RDBMS Characteristics

| Characteristic | Description |
|---------------|-------------|
| **Schema-on-Write** | Structure defined before data insertion |
| **ACID Compliance** | Atomicity, Consistency, Isolation, Durability |
| **SQL Interface** | Standardized query language |
| **Joins** | Combine data from multiple tables |
| **Indexes** | B-tree, hash, and specialized indexes |
| **Transactions** | Multi-statement atomic operations |
| **Constraints** | Primary key, foreign key, unique, check |

### Normalization Theory

**Why Normalize?**
- Eliminate redundant data
- Ensure data dependencies make sense
- Reduce update anomalies

**Normal Forms:**

| Form | Rule | Example Violation |
|------|------|-------------------|
| **1NF** | Atomic values, no repeating groups | Comma-separated tags in one field |
| **2NF** | 1NF + no partial dependencies | Non-key column depends on part of composite key |
| **3NF** | 2NF + no transitive dependencies | Column A → Column B → Column C |
| **BCNF** | 3NF + every determinant is a candidate key | Stricter than 3NF |

### Popular Relational Databases Comparison

| Feature | PostgreSQL | MySQL | SQL Server | Oracle |
|---------|------------|-------|------------|--------|
| **Open Source** | Yes | Yes | No | No |
| **JSON Support** | Excellent (JSONB) | Good | Good | Good |
| **Full-Text Search** | Built-in | Basic | Full-text | Oracle Text |
| **Replication** | Streaming, Logical | Master-Slave | Always On | Data Guard |
| **Partitioning** | Range, List, Hash | Range, List, Hash, Key | Range, Hash, List | All types |
| **Extensions** | PostGIS, pg_trgm | Limited | CLR Integration | Many |
| **Best For** | Complex queries, Analytics | Web apps, Simple scaling | Enterprise, .NET | Enterprise, Complex |

### When to Use Relational Databases

**Ideal Use Cases:**
- Financial transactions requiring ACID
- Complex reporting and analytics
- Applications with well-defined schemas
- Multi-table joins are common
- Strong data integrity requirements

**Avoid When:**
- Schema changes frequently
- Horizontal scaling is primary concern
- Simple key-value access patterns
- Hierarchical or graph-like data

```sql
-- Example: E-commerce schema (3NF)
CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_id INT REFERENCES customers(id),
    order_date TIMESTAMP DEFAULT NOW(),
    status VARCHAR(50) NOT NULL,
    total DECIMAL(10,2) NOT NULL
);

CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(id),
    product_id INT REFERENCES products(id),
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL
);

-- Complex query example
SELECT c.name, COUNT(o.id) as order_count, SUM(o.total) as total_spent
FROM customers c
JOIN orders o ON c.id = o.customer_id
WHERE o.order_date >= '2024-01-01'
GROUP BY c.id, c.name
HAVING SUM(o.total) > 1000
ORDER BY total_spent DESC;
```

---

## 4. NoSQL Databases

### Theory: The NoSQL Movement

**NoSQL** (Not Only SQL) emerged to address limitations of relational databases at web scale:
- Need for horizontal scaling
- Flexible/dynamic schemas
- Distributed architecture by design
- High availability requirements

**Key Insight:** Different data access patterns require different data models. One size doesn't fit all.

### Document Databases

**Theory:**
Document databases store data as self-contained documents (typically JSON/BSON). Each document contains all related data, reducing the need for joins.

**Characteristics:**
- Schema-flexible (schema-on-read)
- Documents can have different structures
- Embedded data for frequently accessed together
- References for many-to-many relationships

**Popular Options:**

| Database | Strengths | Considerations |
|----------|-----------|----------------|
| **MongoDB** | Rich queries, aggregation, horizontal scaling | Memory-intensive for large datasets |
| **CouchDB** | Multi-master replication, offline-first | Slower queries than MongoDB |
| **Firestore** | Real-time sync, mobile-friendly | Google Cloud lock-in |
| **DocumentDB** | MongoDB-compatible, AWS managed | Not fully MongoDB compatible |

```javascript
// Document model example
{
    "_id": "order_123",
    "customer": {
        "id": "cust_456",
        "name": "John Doe",
        "email": "john@example.com"
    },
    "items": [
        {"product": "Laptop", "quantity": 1, "price": 999.99},
        {"product": "Mouse", "quantity": 2, "price": 29.99}
    ],
    "total": 1059.97,
    "status": "shipped",
    "created_at": "2024-03-15T10:00:00Z"
}
```

### Key-Value Databases

**Theory:**
The simplest NoSQL model - data stored as key-value pairs. Optimized for simple CRUD operations with O(1) lookups.

**Characteristics:**
- Extremely fast reads/writes
- Limited query capability (by key only)
- Often used for caching
- Horizontal scaling via consistent hashing

**Popular Options:**

| Database | Strengths | Best For |
|----------|-----------|----------|
| **Redis** | Rich data structures, pub/sub, Lua scripting | Caching, sessions, real-time |
| **Memcached** | Simple, multi-threaded, memory-efficient | Pure caching |
| **DynamoDB** | Fully managed, auto-scaling, global tables | Serverless, AWS ecosystem |
| **etcd** | Strong consistency, distributed | Configuration, service discovery |

```python
# Key-value patterns
# Session storage
SET session:abc123 '{"user_id": 456, "expires": "2024-03-15T12:00:00Z"}'
GET session:abc123

# Rate limiting
INCR rate_limit:user:456:minute
EXPIRE rate_limit:user:456:minute 60

# Caching
SETEX cache:user:456 3600 '{"name": "John", "email": "john@example.com"}'
```

### Wide-Column Databases

**Theory:**
Wide-column stores organize data by column families rather than rows. Optimized for write-heavy workloads and time-series data.

**Characteristics:**
- Column-oriented storage
- Sparse data handling (no null storage)
- Configurable consistency levels
- High write throughput

**Popular Options:**

| Database | Strengths | Best For |
|----------|-----------|----------|
| **Cassandra** | Linear scalability, tunable consistency | Time-series, write-heavy |
| **HBase** | Hadoop integration, strong consistency | Big data analytics |
| **ScyllaDB** | Cassandra-compatible, higher performance | Drop-in Cassandra replacement |

```cql
-- Cassandra data model example
CREATE TABLE sensor_data (
    sensor_id UUID,
    timestamp TIMESTAMP,
    temperature FLOAT,
    humidity FLOAT,
    PRIMARY KEY (sensor_id, timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);

-- Time-series query
SELECT * FROM sensor_data 
WHERE sensor_id = uuid() 
AND timestamp >= '2024-03-14' 
AND timestamp < '2024-03-15';
```

### Graph Databases

**Theory:**
Graph databases model data as nodes (entities) and edges (relationships). Optimized for traversing relationships, which is expensive in relational joins.

**Core Concepts:**
- **Nodes**: Entities (Person, Product, Location)
- **Edges**: Relationships (KNOWS, PURCHASED, LOCATED_IN)
- **Properties**: Key-value pairs on nodes and edges

**When Graphs Excel:**
- Relationships are first-class citizens
- Need to traverse multiple relationship hops
- Relationship patterns are dynamic
- Questions like "friends of friends who bought X"

**Popular Options:**

| Database | Query Language | Best For |
|----------|---------------|----------|
| **Neo4j** | Cypher | General graph, knowledge graphs |
| **Amazon Neptune** | Gremlin, SPARQL | AWS managed, RDF support |
| **JanusGraph** | Gremlin | Distributed, large-scale graphs |
| **TigerGraph** | GSQL | Real-time analytics |

```cypher
// Neo4j Cypher examples
// Find friends of friends
MATCH (user:Person {name: 'John'})-[:KNOWS]->(friend)-[:KNOWS]->(fof)
WHERE NOT (user)-[:KNOWS]->(fof) AND user <> fof
RETURN DISTINCT fof.name;

// Product recommendations
MATCH (user:Person {id: 'user123'})-[:PURCHASED]->(product)<-[:PURCHASED]-(other)
      -[:PURCHASED]->(recommendation)
WHERE NOT (user)-[:PURCHASED]->(recommendation)
RETURN recommendation.name, COUNT(*) AS score
ORDER BY score DESC
LIMIT 10;
```

---

## 5. NewSQL and Distributed SQL

### Theory: The Best of Both Worlds

**NewSQL** databases combine:
- **SQL interface** and relational model
- **ACID transactions** across distributed nodes
- **Horizontal scalability** of NoSQL
- **High availability** through replication

**Why NewSQL Emerged:**
Traditional RDBMS struggled with horizontal scaling. NoSQL sacrificed SQL and ACID. NewSQL provides a middle ground.

### How Distributed SQL Works

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DISTRIBUTED SQL ARCHITECTURE                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                         ┌─────────────┐                                     │
│                         │ SQL Query   │                                     │
│                         │   Layer     │                                     │
│                         └──────┬──────┘                                     │
│                                │                                            │
│                    ┌───────────┼───────────┐                                │
│                    │           │           │                                │
│                    ▼           ▼           ▼                                │
│              ┌──────────┐ ┌──────────┐ ┌──────────┐                        │
│              │  Node 1  │ │  Node 2  │ │  Node 3  │                        │
│              │ ┌──────┐ │ │ ┌──────┐ │ │ ┌──────┐ │                        │
│              │ │Data  │ │ │ │Data  │ │ │ │Data  │ │                        │
│              │ │Shard │ │ │ │Shard │ │ │ │Shard │ │                        │
│              │ └──────┘ │ │ └──────┘ │ │ └──────┘ │                        │
│              │ ┌──────┐ │ │ ┌──────┐ │ │ ┌──────┐ │                        │
│              │ │Raft  │ │ │ │Raft  │ │ │ │Raft  │ │                        │
│              │ │Group │ │ │ │Group │ │ │ │Group │ │                        │
│              │ └──────┘ │ │ └──────┘ │ │ └──────┘ │                        │
│              └──────────┘ └──────────┘ └──────────┘                        │
│                                                                              │
│  Key Technologies:                                                          │
│  • Distributed consensus (Raft/Paxos)                                       │
│  • Automatic sharding                                                       │
│  • Distributed transactions (2PC with optimizations)                        │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Popular NewSQL Databases

| Database | Open Source | Best For | Notable Features |
|----------|-------------|----------|------------------|
| **CockroachDB** | Yes | Global distribution, PostgreSQL compatible | Serializable isolation, geo-partitioning |
| **TiDB** | Yes | MySQL compatible, HTAP | TiKV storage, TiFlash for analytics |
| **YugabyteDB** | Yes | PostgreSQL/Cassandra compatible | Raft consensus, YSQL + YCQL |
| **Spanner** | No (GCP) | Global consistency, Google scale | TrueTime, external consistency |
| **PlanetScale** | No (managed) | MySQL, serverless | Vitess-based, schema changes |

### When to Consider NewSQL

**Ideal Use Cases:**
- Need SQL but outgrowing single-node RDBMS
- Global application requiring distributed data
- ACID transactions across regions
- Want relational model without NoSQL trade-offs

**Considerations:**
- Higher operational complexity than single-node
- Latency for cross-region transactions
- Cost typically higher than single-node
- Team needs distributed systems knowledge

```sql
-- CockroachDB example: Geo-partitioned table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email STRING NOT NULL,
    region STRING NOT NULL,
    data JSONB
) PARTITION BY LIST (region) (
    PARTITION us_west VALUES IN ('us-west'),
    PARTITION us_east VALUES IN ('us-east'),
    PARTITION europe VALUES IN ('eu-west', 'eu-central')
);

-- Assign partitions to specific regions
ALTER PARTITION us_west OF TABLE users CONFIGURE ZONE USING
    constraints='[+region=us-west1]';
ALTER PARTITION europe OF TABLE users CONFIGURE ZONE USING
    constraints='[+region=europe-west1]';
```

---

## 6. Specialized Databases

### Time-Series Databases

**Theory:**
Optimized for timestamped data points. Use cases: IoT sensors, application metrics, financial data.

**Optimizations:**
- Time-based partitioning
- Compression algorithms for sequential data
- Automatic rollups/downsampling
- Time-range query optimization

| Database | Best For | Key Features |
|----------|----------|--------------|
| **InfluxDB** | Metrics, IoT | InfluxQL/Flux, retention policies |
| **TimescaleDB** | PostgreSQL users | Full SQL, hypertables |
| **Prometheus** | Monitoring | Pull model, PromQL, alerting |
| **QuestDB** | High-frequency data | Sub-millisecond queries |

### Vector Databases

**Theory:**
Store and query high-dimensional vectors for AI/ML applications. Enable similarity search using distance metrics (cosine, Euclidean).

**Use Cases:**
- Semantic search
- Recommendation systems
- Image/audio similarity
- RAG (Retrieval Augmented Generation)

| Database | Best For | Key Features |
|----------|----------|--------------|
| **Pinecone** | Production AI apps | Managed, hybrid search |
| **Milvus** | Self-hosted, large scale | Multiple indexes, GPU support |
| **Weaviate** | Knowledge graphs + vectors | GraphQL API, modules |
| **Chroma** | Local development | Simple API, embeddings |
| **pgvector** | PostgreSQL integration | HNSW, IVFFlat indexes |

### Search Engines

**Theory:**
Inverted index-based systems optimized for full-text search and log analytics.

| Database | Best For | Key Features |
|----------|----------|--------------|
| **Elasticsearch** | Search + analytics | Full-text, aggregations, ML |
| **OpenSearch** | AWS ecosystem | Elasticsearch fork |
| **Solr** | Enterprise search | Mature, Lucene-based |
| **Meilisearch** | Simple search | Typo-tolerant, instant |

### In-Memory Databases

**Theory:**
Data stored primarily in RAM for sub-millisecond access times.

| Database | Type | Best For |
|----------|------|----------|
| **Redis** | Key-Value + Structures | Caching, sessions, pub/sub |
| **Memcached** | Key-Value | Simple caching |
| **VoltDB** | Relational | ACID in-memory |
| **SAP HANA** | Multi-model | Enterprise analytics |

---

## 7. ACID vs BASE: Consistency Models

### Theory: The Consistency Spectrum

Database systems face fundamental trade-offs between consistency and availability. ACID and BASE represent two philosophies on this spectrum.

### ACID Properties

**ACID** ensures reliable transactions:

| Property | Meaning | Example |
|----------|---------|---------|
| **Atomicity** | All or nothing | Transfer fails = both accounts unchanged |
| **Consistency** | Valid state to valid state | Constraints always satisfied |
| **Isolation** | Concurrent transactions don't interfere | Serializable reads |
| **Durability** | Committed data persists | Survives system crash |

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ACID TRANSACTION EXAMPLE                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  Bank Transfer: $100 from Account A to Account B                            │
│                                                                              │
│  BEGIN TRANSACTION;                                                          │
│    UPDATE accounts SET balance = balance - 100 WHERE id = 'A';              │
│    UPDATE accounts SET balance = balance + 100 WHERE id = 'B';              │
│  COMMIT;                                                                     │
│                                                                              │
│  Atomicity: Both updates happen or neither                                   │
│  Consistency: Total money in system unchanged                                │
│  Isolation: Other transactions see either before or after, not during       │
│  Durability: Once committed, survives power failure                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### BASE Properties

**BASE** prioritizes availability over consistency:

| Property | Meaning | Implication |
|----------|---------|-------------|
| **Basically Available** | System appears to work | May return stale data |
| **Soft state** | State may change without input | Due to eventual consistency |
| **Eventually consistent** | Given time, all nodes converge | Not immediately consistent |

### Isolation Levels (SQL Standard)

| Level | Dirty Read | Non-Repeatable Read | Phantom Read | Performance |
|-------|------------|---------------------|--------------|-------------|
| **Read Uncommitted** | Possible | Possible | Possible | Fastest |
| **Read Committed** | No | Possible | Possible | Fast |
| **Repeatable Read** | No | No | Possible | Medium |
| **Serializable** | No | No | No | Slowest |

### When to Use Each

| ACID (Strong Consistency) | BASE (Eventual Consistency) |
|--------------------------|----------------------------|
| Financial transactions | Social media feeds |
| Inventory management | Shopping cart (temporary) |
| Healthcare records | Analytics dashboards |
| Legal documents | Session data |
| Banking | Content caching |

---

## 8. CAP Theorem and Trade-offs

### Theory: The CAP Theorem

The **CAP theorem** (Brewer's theorem) states: A distributed system can provide at most **two of three** guarantees simultaneously:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           CAP THEOREM                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                            Consistency                                       │
│                                 /\                                           │
│                                /  \                                          │
│                               /    \                                         │
│                              /  CP  \                                        │
│                             /________\                                       │
│                            /          \                                      │
│                           /     CA     \                                     │
│                          /              \                                    │
│                         /________________\                                   │
│                  Availability            Partition                          │
│                                          Tolerance                          │
│                                                                              │
│  CP: Consistency + Partition Tolerance (sacrifice Availability)             │
│      Examples: MongoDB, HBase, Redis Cluster                                │
│                                                                              │
│  AP: Availability + Partition Tolerance (sacrifice Consistency)             │
│      Examples: Cassandra, DynamoDB, CouchDB                                 │
│                                                                              │
│  CA: Consistency + Availability (sacrifice Partition Tolerance)             │
│      Examples: Single-node RDBMS (not truly distributed)                    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Understanding the Trade-offs

**Consistency (C):**
Every read receives the most recent write or an error.

**Availability (A):**
Every request receives a (non-error) response, without guarantee of most recent data.

**Partition Tolerance (P):**
System continues operating despite network partitions between nodes.

### Reality: PACELC Theorem

The CAP theorem only describes behavior during partitions. **PACELC** extends this:

```
If there is a Partition:
    Choose between Availability and Consistency
Else (normal operation):
    Choose between Latency and Consistency
```

| Category | During Partition | Normal Operation | Examples |
|----------|-----------------|------------------|----------|
| **PA/EL** | Availability | Latency | Cassandra, DynamoDB |
| **PA/EC** | Availability | Consistency | PNUTS |
| **PC/EL** | Consistency | Latency | MongoDB, HBase |
| **PC/EC** | Consistency | Consistency | VoltDB, BigTable |

### Database CAP Classification

| Database | CAP | Behavior |
|----------|-----|----------|
| PostgreSQL (single) | CA | No partition tolerance |
| MongoDB | CP | Primary unavailable during partition |
| Cassandra | AP | Always writable, eventual consistency |
| CockroachDB | CP | Prefers consistency over availability |
| DynamoDB | AP (tunable) | Configurable per operation |
| Redis Cluster | CP | Minority partition unavailable |
| Neo4j (cluster) | CP | Consistency prioritized |

---

## 9. Database Selection Criteria

### Theory: Systematic Selection Process

Choosing a database is a multi-factor decision. Consider these dimensions:

### Primary Selection Factors

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATABASE SELECTION FACTORS                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. DATA MODEL FIT                                                          │
│     └── Does the data model match your domain?                              │
│                                                                              │
│  2. QUERY PATTERNS                                                          │
│     └── What queries will you run frequently?                               │
│                                                                              │
│  3. SCALABILITY REQUIREMENTS                                                │
│     └── Vertical (bigger machine) vs Horizontal (more machines)?            │
│                                                                              │
│  4. CONSISTENCY REQUIREMENTS                                                │
│     └── Strong (ACID) vs Eventual (BASE)?                                   │
│                                                                              │
│  5. AVAILABILITY REQUIREMENTS                                               │
│     └── Uptime SLA? Disaster recovery needs?                                │
│                                                                              │
│  6. OPERATIONAL COMPLEXITY                                                  │
│     └── Team expertise? Managed vs self-hosted?                             │
│                                                                              │
│  7. ECOSYSTEM & TOOLING                                                     │
│     └── Libraries, ORM support, monitoring tools?                           │
│                                                                              │
│  8. COST                                                                    │
│     └── Licensing, infrastructure, operations?                              │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Decision Matrix Template

| Criterion | Weight | Option A | Option B | Option C |
|-----------|--------|----------|----------|----------|
| Data model fit | 25% | Score 1-5 | Score 1-5 | Score 1-5 |
| Query support | 20% | Score 1-5 | Score 1-5 | Score 1-5 |
| Scalability | 15% | Score 1-5 | Score 1-5 | Score 1-5 |
| Consistency | 15% | Score 1-5 | Score 1-5 | Score 1-5 |
| Operational ease | 10% | Score 1-5 | Score 1-5 | Score 1-5 |
| Ecosystem | 10% | Score 1-5 | Score 1-5 | Score 1-5 |
| Cost | 5% | Score 1-5 | Score 1-5 | Score 1-5 |
| **Weighted Total** | 100% | **Sum** | **Sum** | **Sum** |

### Questions to Ask

**Data Characteristics:**
- What is the data structure (tabular, hierarchical, graph)?
- What is the expected data volume (GB, TB, PB)?
- What is the growth rate?
- Is data frequently updated or append-only?

**Access Patterns:**
- Read-heavy, write-heavy, or balanced?
- Point lookups vs complex queries vs aggregations?
- Real-time requirements vs batch processing?
- Geographic distribution of users?

**Consistency Needs:**
- Can application tolerate stale reads?
- Are multi-document transactions required?
- What is the cost of inconsistency?

**Operational:**
- What is the team's expertise?
- Self-hosted or managed service preferred?
- What monitoring/tooling exists?
- Disaster recovery requirements?

---

## 10. Data Modeling Approaches

### Theory: Different Paradigms

Data modeling varies significantly across database types. Understanding these differences is crucial for optimal design.

### Relational Data Modeling

**Principles:**
- Normalize to reduce redundancy
- Use foreign keys for relationships
- Design for data integrity first
- Optimize with denormalization if needed

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    RELATIONAL MODEL (3NF)                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────┐                │
│  │   CUSTOMERS  │     │    ORDERS    │     │ ORDER_ITEMS  │                │
│  ├──────────────┤     ├──────────────┤     ├──────────────┤                │
│  │ id (PK)      │     │ id (PK)      │     │ id (PK)      │                │
│  │ name         │◄────│ customer_id  │     │ order_id     │────┐           │
│  │ email        │     │ order_date   │◄────│ product_id   │    │           │
│  └──────────────┘     │ status       │     │ quantity     │    │           │
│                       └──────────────┘     │ price        │    │           │
│                                            └──────────────┘    │           │
│                                                   ┌─────────────┘           │
│                                                   ▼                         │
│                                            ┌──────────────┐                │
│                                            │   PRODUCTS   │                │
│                                            ├──────────────┤                │
│                                            │ id (PK)      │                │
│                                            │ name         │                │
│                                            │ price        │                │
│                                            └──────────────┘                │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Document Data Modeling

**Principles:**
- Model for queries, not relationships
- Embed frequently accessed together data
- Reference when data is accessed separately
- Accept some denormalization

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DOCUMENT MODEL                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  {                                                                           │
│    "_id": "order_123",                                                       │
│    "customer": {                    // Embedded - accessed together          │
│      "id": "cust_456",                                                       │
│      "name": "John Doe",                                                     │
│      "email": "john@example.com"                                             │
│    },                                                                        │
│    "items": [                       // Embedded array                        │
│      {                                                                       │
│        "product_id": "prod_789",    // Reference for full product details    │
│        "name": "Laptop",            // Denormalized for display              │
│        "quantity": 1,                                                        │
│        "price": 999.99                                                       │
│      }                                                                       │
│    ],                                                                        │
│    "total": 999.99,                 // Pre-computed                          │
│    "status": "shipped",                                                      │
│    "created_at": "2024-03-15"                                                │
│  }                                                                           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Key-Value Data Modeling

**Principles:**
- Design keys for access patterns
- Use prefixes for namespacing
- Consider serialization format
- Plan for key expiration

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    KEY-VALUE PATTERNS                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  User Sessions:                                                              │
│    Key: session:{session_id}                                                │
│    Value: {"user_id": 123, "created_at": "...", "expires_at": "..."}        │
│                                                                              │
│  Caching:                                                                    │
│    Key: cache:user:{user_id}                                                │
│    Value: {full user object}                                                │
│    TTL: 3600 seconds                                                        │
│                                                                              │
│  Counters:                                                                   │
│    Key: counter:page_views:{page_id}:{date}                                 │
│    Value: 12345 (integer)                                                   │
│                                                                              │
│  Rate Limiting:                                                              │
│    Key: ratelimit:{user_id}:{endpoint}:{minute}                             │
│    Value: request count                                                     │
│    TTL: 60 seconds                                                          │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Graph Data Modeling

**Principles:**
- Entities become nodes with labels
- Relationships become edges with types
- Properties on both nodes and edges
- Design for traversal patterns

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    GRAPH MODEL                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│         ┌─────────────────────────────────────────────────┐                 │
│         │                   FOLLOWS                        │                 │
│         ▼                                                 │                 │
│    ┌─────────┐   KNOWS   ┌─────────┐   WORKS_AT   ┌──────┴────┐            │
│    │ Person  │◄─────────►│ Person  │─────────────►│  Company  │            │
│    │         │           │         │              │           │            │
│    │ name:   │           │ name:   │              │ name:     │            │
│    │ "Alice" │           │ "Bob"   │              │ "Acme"    │            │
│    └─────────┘           └────┬────┘              └───────────┘            │
│                               │                                            │
│                          PURCHASED                                         │
│                          {date: "..."}                                     │
│                               │                                            │
│                               ▼                                            │
│                          ┌─────────┐                                       │
│                          │ Product │                                       │
│                          │         │                                       │
│                          │ name:   │                                       │
│                          │ "Laptop"│                                       │
│                          └─────────┘                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 11. Scaling Strategies

### Theory: Vertical vs Horizontal Scaling

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    SCALING STRATEGIES                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  VERTICAL SCALING (Scale Up)          HORIZONTAL SCALING (Scale Out)        │
│                                                                              │
│  ┌─────────────────┐                  ┌───┐ ┌───┐ ┌───┐ ┌───┐              │
│  │                 │                  │   │ │   │ │   │ │   │              │
│  │                 │                  │DB1│ │DB2│ │DB3│ │DB4│              │
│  │  BIGGER         │                  │   │ │   │ │   │ │   │              │
│  │  SERVER         │                  └───┘ └───┘ └───┘ └───┘              │
│  │                 │                                                        │
│  │  More CPU       │                  Add more servers                      │
│  │  More RAM       │                  Data distributed across nodes         │
│  │  More Disk      │                  Each node handles subset              │
│  │                 │                                                        │
│  └─────────────────┘                                                        │
│                                                                              │
│  Pros:                                Pros:                                 │
│  • Simple                             • Near-linear scalability             │
│  • No app changes                     • High availability                   │
│  • ACID preserved                     • Cost-effective at scale             │
│                                                                              │
│  Cons:                                Cons:                                 │
│  • Hardware limits                    • Complexity                          │
│  • Single point of failure            • Distributed transactions hard       │
│  • Expensive at high end              • App changes may be needed           │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Read Scaling Patterns

**1. Read Replicas:**
```
             Writes                    Reads
                │                    ┌───┴───┐
                ▼                    ▼       ▼
           ┌────────┐          ┌────────┐ ┌────────┐
           │PRIMARY │─────────►│REPLICA │ │REPLICA │
           └────────┘          └────────┘ └────────┘
```

**2. Caching Layer:**
```
    Request ──► Cache Hit? ──► Yes ──► Return cached
                   │
                   No
                   │
                   ▼
              Query DB ──► Cache result ──► Return
```

### Write Scaling Patterns

**1. Sharding (Horizontal Partitioning):**
```
    All Data
        │
        ├── Shard 1: Users A-H
        ├── Shard 2: Users I-P
        └── Shard 3: Users Q-Z
```

**2. Sharding Strategies:**

| Strategy | Description | Best For |
|----------|-------------|----------|
| **Hash-based** | Hash function determines shard | Even distribution |
| **Range-based** | Key ranges map to shards | Range queries |
| **Directory-based** | Lookup table maps keys | Flexible, complex |
| **Geographic** | Location determines shard | Data locality |

### Database-Specific Scaling

| Database | Read Scaling | Write Scaling |
|----------|--------------|---------------|
| PostgreSQL | Read replicas, PgBouncer | Citus, manual sharding |
| MySQL | Read replicas, ProxySQL | Vitess, manual sharding |
| MongoDB | Replica sets | Auto-sharding |
| Cassandra | Any node reads | Linear write scaling |
| DynamoDB | Auto-scaling | Auto-scaling (pay per use) |
| CockroachDB | Distributed reads | Distributed writes |

---

## 12. Performance Considerations

### Theory: Performance Dimensions

Database performance is multi-dimensional:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    PERFORMANCE DIMENSIONS                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  LATENCY                          THROUGHPUT                                │
│  ├── Query response time          ├── Operations per second                 │
│  ├── P50, P95, P99 percentiles    ├── Reads/sec, Writes/sec                │
│  └── Network round trips          └── Concurrent connections               │
│                                                                              │
│  RESOURCE UTILIZATION             SCALABILITY                               │
│  ├── CPU usage                    ├── Linear vs logarithmic                │
│  ├── Memory consumption           ├── Scale-up headroom                    │
│  ├── Disk I/O                     └── Scale-out efficiency                 │
│  └── Network bandwidth                                                      │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Key Performance Factors

**1. Indexing Strategy:**
```
                    Without Index          With Index
Query Time:         O(n) - scan all        O(log n) - tree traversal
1M records:         ~1000ms                ~1ms
10M records:        ~10000ms               ~1.3ms
```

**2. Query Optimization:**
- Use EXPLAIN/EXPLAIN ANALYZE to understand execution plans
- Avoid SELECT * when not needed
- Use pagination for large result sets
- Batch operations when possible

**3. Connection Management:**
```
Connection Pooling Benefits:
├── Avoid connection overhead (50-200ms per connection)
├── Reuse existing connections
├── Limit server-side resource usage
└── Handle connection failures gracefully
```

**4. Caching Strategies:**

| Cache Type | Location | Use Case |
|------------|----------|----------|
| Query cache | Database | Repeated identical queries |
| Result cache | Application | Computed results |
| Object cache | Application/External | Frequently accessed entities |
| CDN cache | Edge | Static/semi-static content |

### Performance Anti-Patterns

| Anti-Pattern | Problem | Solution |
|--------------|---------|----------|
| N+1 queries | Multiple round trips | Use JOINs or batch fetching |
| Missing indexes | Full table scans | Analyze queries, add indexes |
| Over-indexing | Slow writes, storage waste | Remove unused indexes |
| Large transactions | Lock contention | Break into smaller units |
| SELECT * | Excess data transfer | Select only needed columns |
| No connection pooling | Connection overhead | Implement pooling |

### Benchmarking Guidelines

```
Performance Benchmarking Checklist:
□ Use production-like data volumes
□ Simulate realistic query patterns
□ Test at expected peak load
□ Measure latency percentiles (P50, P95, P99)
□ Monitor resource utilization
□ Test failure scenarios
□ Warm up caches before measuring
□ Run multiple iterations for statistical significance
```

---

## 13. High Availability and Disaster Recovery

### Theory: Availability Fundamentals

**Availability = Uptime / Total Time**

| Availability | Downtime/Year | Downtime/Month | Downtime/Day |
|--------------|---------------|----------------|--------------|
| 99% (two 9s) | 3.65 days | 7.31 hours | 14.4 minutes |
| 99.9% (three 9s) | 8.76 hours | 43.8 minutes | 1.44 minutes |
| 99.99% (four 9s) | 52.6 minutes | 4.38 minutes | 8.64 seconds |
| 99.999% (five 9s) | 5.26 minutes | 26.3 seconds | 864 milliseconds |

### High Availability Patterns

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    HA ARCHITECTURE PATTERNS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ACTIVE-PASSIVE (Failover)                                                  │
│                                                                              │
│  ┌────────┐    ┌────────┐                                                   │
│  │ PRIMARY│────│STANDBY │   Standby takes over if primary fails            │
│  │ Active │    │Passive │   RPO: Nearly zero with sync replication         │
│  └────────┘    └────────┘   RTO: Minutes (failover time)                   │
│                                                                              │
│  ACTIVE-ACTIVE (Multi-Primary)                                              │
│                                                                              │
│  ┌────────┐    ┌────────┐                                                   │
│  │PRIMARY │◄──►│PRIMARY │   Both accept writes, sync bidirectionally       │
│  │   1    │    │   2    │   Conflict resolution needed                      │
│  └────────┘    └────────┘   Higher complexity, lower RTO                    │
│                                                                              │
│  REPLICA SET (MongoDB, etc.)                                                │
│                                                                              │
│  ┌────────┐    ┌────────┐    ┌────────┐                                    │
│  │PRIMARY │────│SECOND- │────│SECOND- │   Automatic election               │
│  │        │    │  ARY   │    │  ARY   │   Self-healing                     │
│  └────────┘    └────────┘    └────────┘   RTO: 10-30 seconds               │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Recovery Objectives

| Metric | Definition | Example |
|--------|------------|---------|
| **RPO** (Recovery Point Objective) | Maximum acceptable data loss | 1 hour: Can lose up to 1 hour of data |
| **RTO** (Recovery Time Objective) | Maximum acceptable downtime | 15 minutes: Must recover within 15 min |

### Disaster Recovery Strategies

| Strategy | RPO | RTO | Cost |
|----------|-----|-----|------|
| Backup & Restore | Hours-Days | Hours-Days | Low |
| Pilot Light | Minutes-Hours | Hours | Medium |
| Warm Standby | Minutes | Minutes-Hour | High |
| Multi-Site Active-Active | Near Zero | Near Zero | Highest |

### Database-Specific HA Solutions

| Database | HA Solution | Failover Time |
|----------|-------------|---------------|
| PostgreSQL | Patroni, pg_auto_failover | 10-30 seconds |
| MySQL | Group Replication, InnoDB Cluster | 10-30 seconds |
| MongoDB | Replica Sets | 10-12 seconds |
| Cassandra | Built-in (any node) | Near instant |
| CockroachDB | Built-in (Raft) | Near instant |

---

## 14. Security Considerations

### Theory: Database Security Layers

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATABASE SECURITY LAYERS                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      APPLICATION LAYER                               │   │
│  │  • Input validation    • Parameterized queries    • Error handling   │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      ACCESS CONTROL LAYER                            │   │
│  │  • Authentication      • Authorization (RBAC)     • Row-level security│  │
│  └─────────────────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      NETWORK LAYER                                   │   │
│  │  • TLS encryption      • IP whitelisting         • VPC/Private network│  │
│  └─────────────────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      DATA LAYER                                      │   │
│  │  • Encryption at rest  • Field-level encryption  • Data masking      │   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────────────────────┐   │
│  │                      AUDIT & COMPLIANCE                              │   │
│  │  • Activity logging    • Access auditing         • Compliance reports│   │
│  └─────────────────────────────────────────────────────────────────────┘   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Authentication Methods

| Method | Security Level | Complexity | Best For |
|--------|---------------|------------|----------|
| Password (SCRAM) | Medium | Low | Development, simple deployments |
| LDAP/AD | Medium-High | Medium | Enterprise environments |
| Kerberos | High | High | Windows-heavy environments |
| X.509 Certificates | High | High | Zero-trust, service-to-service |
| IAM Integration | High | Medium | Cloud-native applications |

### Encryption Types

| Type | Protects Against | Implementation |
|------|-----------------|----------------|
| **In-Transit (TLS)** | Eavesdropping, MITM | TLS certificates |
| **At-Rest** | Physical theft, disk access | Database encryption features |
| **Field-Level** | Admin access, backup exposure | Application-level (CSFLE) |

### Security Best Practices Checklist

```
□ Enable authentication (never run without auth in production)
□ Use TLS for all connections
□ Implement least-privilege access
□ Enable encryption at rest
□ Use parameterized queries (prevent injection)
□ Regular security patching
□ Enable audit logging
□ Regular backup testing
□ Network isolation (VPC, firewall)
□ Regular security assessments
```

---

## 15. Cost Analysis

### Theory: Total Cost of Ownership (TCO)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATABASE TCO COMPONENTS                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  INFRASTRUCTURE COSTS                                                       │
│  ├── Compute (CPU, Memory)                                                  │
│  ├── Storage (SSD, HDD, IOPS)                                               │
│  ├── Network (Bandwidth, egress)                                            │
│  └── Backup storage                                                         │
│                                                                              │
│  LICENSING COSTS                                                            │
│  ├── Database license (if proprietary)                                      │
│  ├── Enterprise features                                                    │
│  └── Support contracts                                                      │
│                                                                              │
│  OPERATIONAL COSTS                                                          │
│  ├── DBA time and expertise                                                 │
│  ├── Monitoring tools                                                       │
│  ├── Training                                                               │
│  └── Incident response                                                      │
│                                                                              │
│  OPPORTUNITY COSTS                                                          │
│  ├── Developer productivity                                                 │
│  ├── Time to market                                                         │
│  └── Technical debt                                                         │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Cost Comparison: Self-Hosted vs Managed

| Factor | Self-Hosted | Managed Service |
|--------|-------------|-----------------|
| Infrastructure | Lower base cost | Higher base cost |
| Operations | High (DBA needed) | Low (included) |
| Scaling | Manual, complex | Often automatic |
| Availability | DIY setup | Built-in HA |
| Updates/Patches | Manual | Automatic |
| Expertise needed | High | Lower |
| **Best for** | Large scale, specific needs | Most applications |

### Managed Database Pricing Models

| Pricing Model | Best For | Examples |
|---------------|----------|----------|
| **Instance-based** | Predictable workloads | RDS, Cloud SQL |
| **Serverless** | Variable workloads | Aurora Serverless, DynamoDB |
| **Throughput-based** | High-volume apps | DynamoDB, Cosmos DB |
| **Storage-based** | Data-heavy, low query | S3 + Athena |

### Cost Optimization Strategies

```
Cost Reduction Tactics:
├── Right-size instances (don't over-provision)
├── Use reserved instances for predictable workloads
├── Implement caching to reduce database load
├── Archive cold data to cheaper storage
├── Optimize queries to reduce compute
├── Use read replicas judiciously
├── Consider serverless for variable workloads
└── Regular cost reviews and cleanup
```

---

## 16. Migration Strategies

### Theory: Migration Approaches

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    MIGRATION STRATEGIES                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  BIG BANG MIGRATION                                                         │
│  ┌──────────┐                    ┌──────────┐                               │
│  │ Old DB   │ ──── Downtime ────►│ New DB   │                               │
│  └──────────┘                    └──────────┘                               │
│  • Simple but risky                                                         │
│  • Requires maintenance window                                              │
│  • All-or-nothing                                                           │
│                                                                              │
│  PARALLEL RUN                                                               │
│  ┌──────────┐           ┌──────────┐                                        │
│  │ Old DB   │◄──Writes─►│ New DB   │                                        │
│  └──────────┘           └──────────┘                                        │
│  • Dual writes during transition                                            │
│  • Compare results for validation                                           │
│  • Higher operational cost                                                  │
│                                                                              │
│  STRANGLER FIG                                                              │
│  ┌──────────────────────────────────┐                                       │
│  │   Feature A ────► New DB         │                                       │
│  │   Feature B ────► New DB         │  Gradual migration by feature         │
│  │   Feature C ────► Old DB (still) │                                       │
│  └──────────────────────────────────┘                                       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Migration Planning Checklist

```
Pre-Migration:
□ Inventory current schema and data
□ Document all access patterns
□ Identify data transformation needs
□ Plan for downtime (if any)
□ Set success criteria
□ Create rollback plan
□ Test migration on staging environment

During Migration:
□ Monitor progress and errors
□ Validate data integrity
□ Test application functionality
□ Document any issues

Post-Migration:
□ Verify all data migrated correctly
□ Performance testing
□ Update documentation
□ Retire old database (after confirmation)
```

### Data Transformation Considerations

| Challenge | Strategy |
|-----------|----------|
| Schema differences | ETL transformation rules |
| Data type changes | Conversion functions |
| Null handling | Default values or explicit handling |
| Encoding changes | Character set conversion |
| Reference integrity | Order of table migration |

---

## 17. Polyglot Persistence

### Theory: Right Tool for Each Job

**Polyglot persistence** means using different databases for different use cases within the same application, choosing the best fit for each data access pattern.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    POLYGLOT PERSISTENCE ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                          ┌─────────────────┐                                │
│                          │   Application   │                                │
│                          │     Layer       │                                │
│                          └────────┬────────┘                                │
│                                   │                                         │
│         ┌───────────┬─────────────┼─────────────┬───────────┐              │
│         │           │             │             │           │              │
│         ▼           ▼             ▼             ▼           ▼              │
│    ┌─────────┐ ┌─────────┐ ┌───────────┐ ┌─────────┐ ┌───────────┐        │
│    │ Postgre │ │ MongoDB │ │   Redis   │ │  Neo4j  │ │Elasticsearch│       │
│    │   SQL   │ │         │ │           │ │         │ │           │        │
│    └─────────┘ └─────────┘ └───────────┘ └─────────┘ └───────────┘        │
│                                                                              │
│    Transactions Content    Caching &    Social      Full-text              │
│    & Reporting  Mgmt       Sessions     Graph       Search                 │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### When to Use Polyglot Persistence

**Benefits:**
- Optimal performance for each use case
- Better fit for different data models
- Scale different components independently
- Avoid forcing square peg into round hole

**Challenges:**
- Increased operational complexity
- Data synchronization between systems
- Distributed transactions are hard
- Team needs broader expertise

### Common Polyglot Combinations

| Primary DB | + Complementary | Use Case |
|------------|-----------------|----------|
| PostgreSQL | + Redis | Relational data + Caching |
| PostgreSQL | + Elasticsearch | Relational + Full-text search |
| MongoDB | + Redis | Documents + Sessions/Cache |
| Any RDBMS | + Neo4j | Core data + Graph relationships |
| Any DB | + InfluxDB | Business data + Metrics/Monitoring |

### Data Synchronization Patterns

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATA SYNC PATTERNS                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  DUAL WRITE (Application writes to both)                                    │
│  ├── Simple but risky (consistency issues)                                  │
│  └── One write can fail, leaving inconsistent state                         │
│                                                                              │
│  CHANGE DATA CAPTURE (CDC)                                                  │
│  ├── Capture changes from primary DB                                        │
│  ├── Stream to secondary systems                                            │
│  └── Tools: Debezium, AWS DMS, Kafka Connect                                │
│                                                                              │
│  EVENT SOURCING                                                             │
│  ├── Events are source of truth                                             │
│  ├── Materialize views in each database                                     │
│  └── Eventually consistent but reliable                                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 18. Cloud vs On-Premise Databases

### Theory: Deployment Model Trade-offs

```
┌─────────────────────────────────────────────────────────────────────────────┐
│              CLOUD VS ON-PREMISE COMPARISON                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│                    ON-PREMISE                  CLOUD MANAGED                │
│                                                                              │
│  Control         Full control                 Limited customization        │
│  Scaling         Manual, hardware lead time   Minutes, API-driven          │
│  Cost Model      CapEx (upfront)              OpEx (pay-as-you-go)         │
│  Operations      You manage everything        Provider manages infra       │
│  Expertise       Need DBA, SysAdmin           Less specialized staff       │
│  Security        You control all layers       Shared responsibility        │
│  Latency         Predictable                  Variable (can optimize)      │
│  Compliance      Easier for some requirements May have limitations         │
│  Vendor Lock-in  Low                          Moderate to High             │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Cloud Database Categories

**1. Lift-and-Shift (IaaS):**
- Run traditional DB on cloud VMs
- Full control, full responsibility
- Examples: PostgreSQL on EC2/GCE

**2. Managed Databases (DBaaS):**
- Provider manages infrastructure
- You manage data and queries
- Examples: RDS, Cloud SQL, Atlas

**3. Cloud-Native:**
- Built for cloud from ground up
- Auto-scaling, serverless options
- Examples: Aurora, Spanner, Cosmos DB, DynamoDB

### Major Cloud Database Offerings

| Provider | Relational | Document | Key-Value | Graph |
|----------|------------|----------|-----------|-------|
| **AWS** | RDS, Aurora | DocumentDB | DynamoDB, ElastiCache | Neptune |
| **GCP** | Cloud SQL, Spanner | Firestore | Memorystore | - |
| **Azure** | SQL Database | Cosmos DB | Cosmos DB, Cache for Redis | Cosmos DB |

### Decision Framework: Cloud vs On-Premise

**Choose Cloud When:**
- Variable/unpredictable workloads
- Speed to market is critical
- Limited DBA resources
- Global distribution needed
- Elasticity is required

**Choose On-Premise When:**
- Regulatory requirements mandate it
- Extremely predictable high-volume workloads
- Data sovereignty concerns
- Already have infrastructure and expertise
- Specific hardware requirements

---

## 19. Database Selection Decision Framework

### The Decision Tree

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                    DATABASE SELECTION DECISION TREE                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  START: What is your data structure?                                        │
│         │                                                                   │
│         ├── Tabular with relationships ──► Need ACID?                       │
│         │                                     │                             │
│         │                                     ├── Yes ──► Scale needs?      │
│         │                                     │           │                 │
│         │                                     │           ├── Single node OK │
│         │                                     │           │   ──► PostgreSQL │
│         │                                     │           │                 │
│         │                                     │           └── Distributed    │
│         │                                     │               ──► CockroachDB│
│         │                                     │                              │
│         │                                     └── No ──► Cassandra (writes) │
│         │                                                                   │
│         ├── Hierarchical/Nested ──► MongoDB, CouchDB                        │
│         │                                                                   │
│         ├── Simple Key-Value ──► Redis, DynamoDB                            │
│         │                                                                   │
│         ├── Highly connected (graph) ──► Neo4j, Neptune                     │
│         │                                                                   │
│         ├── Time-stamped sequences ──► InfluxDB, TimescaleDB                │
│         │                                                                   │
│         ├── Vector embeddings ──► Pinecone, Milvus, pgvector                │
│         │                                                                   │
│         └── Full-text search dominant ──► Elasticsearch                     │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Quick Reference Matrix

| Use Case | Primary Recommendation | Alternative |
|----------|----------------------|-------------|
| Web application (general) | PostgreSQL | MySQL |
| Content management | MongoDB | PostgreSQL + JSONB |
| E-commerce | PostgreSQL + Redis | MongoDB |
| Analytics/Reporting | PostgreSQL | ClickHouse |
| Real-time gaming | Redis | DynamoDB |
| Social network | Neo4j + PostgreSQL | MongoDB |
| IoT data | TimescaleDB | InfluxDB, Cassandra |
| Session storage | Redis | DynamoDB |
| Full-text search | Elasticsearch | PostgreSQL (small scale) |
| ML/AI applications | pgvector / Pinecone | Milvus |
| Financial systems | PostgreSQL | CockroachDB |
| Global SaaS | CockroachDB | Spanner |

### Evaluation Scorecard

| Criterion | Weight | DB 1 | DB 2 | DB 3 |
|-----------|--------|------|------|------|
| **Functional Fit** | | | | |
| Data model match | 20% | /5 | /5 | /5 |
| Query capability | 15% | /5 | /5 | /5 |
| Consistency model | 10% | /5 | /5 | /5 |
| **Non-Functional** | | | | |
| Scalability | 15% | /5 | /5 | /5 |
| Performance | 10% | /5 | /5 | /5 |
| Availability | 10% | /5 | /5 | /5 |
| **Operational** | | | | |
| Team expertise | 5% | /5 | /5 | /5 |
| Ecosystem/Tooling | 5% | /5 | /5 | /5 |
| Managed options | 5% | /5 | /5 | /5 |
| Cost | 5% | /5 | /5 | /5 |
| **Weighted Score** | **100%** | | | |

---

## 20. Use Case Scenarios and Recommendations

### Scenario 1: E-Commerce Platform

**Requirements:**
- Product catalog (variable attributes)
- Order management (ACID required)
- User sessions
- Search functionality

**Recommendation:**
```
┌────────────────────────────────────────────────────────────────┐
│                E-Commerce Architecture                          │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PostgreSQL           │ Orders, customers, inventory (ACID)    │
│  MongoDB/Elasticsearch│ Product catalog, search                 │
│  Redis                │ Sessions, caching, cart                 │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

### Scenario 2: Social Media Platform

**Requirements:**
- User profiles
- Posts/content
- Social graph (follows, friends)
- Activity feeds
- Real-time notifications

**Recommendation:**
```
┌────────────────────────────────────────────────────────────────┐
│                Social Media Architecture                        │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PostgreSQL     │ User accounts, core data                      │
│  Neo4j          │ Social graph (relationships)                  │
│  MongoDB        │ Posts, comments, media metadata               │
│  Redis          │ Feeds, notifications, caching                 │
│  Elasticsearch  │ Search                                        │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

### Scenario 3: IoT Platform

**Requirements:**
- High-volume sensor data ingestion
- Time-series queries and aggregations
- Device management
- Real-time alerting

**Recommendation:**
```
┌────────────────────────────────────────────────────────────────┐
│                IoT Platform Architecture                        │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  TimescaleDB/InfluxDB│ Time-series sensor data                  │
│  PostgreSQL          │ Device registry, configuration           │
│  Redis               │ Real-time state, pub/sub                 │
│  Kafka               │ Event streaming                          │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

### Scenario 4: Financial Trading System

**Requirements:**
- Strict ACID compliance
- High-frequency transactions
- Audit logging
- Zero data loss tolerance

**Recommendation:**
```
┌────────────────────────────────────────────────────────────────┐
│                Financial System Architecture                    │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PostgreSQL     │ Transactions, accounts (primary)              │
│  (with Patroni) │ High availability                             │
│  Redis Cluster  │ Rate limiting, session management             │
│  Kafka          │ Transaction event streaming                   │
│  ClickHouse     │ Analytics and reporting                       │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

### Scenario 5: Content Management System

**Requirements:**
- Flexible content schemas
- Rich media management
- Full-text search
- Multi-language support

**Recommendation:**
```
┌────────────────────────────────────────────────────────────────┐
│                CMS Architecture                                 │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  MongoDB        │ Content storage (flexible schema)             │
│  Elasticsearch  │ Full-text search                              │
│  Redis          │ Caching, session management                   │
│  S3-compatible  │ Media file storage                            │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

### Scenario 6: AI/ML Application

**Requirements:**
- Vector embeddings storage
- Similarity search
- Model metadata
- Experiment tracking

**Recommendation:**
```
┌────────────────────────────────────────────────────────────────┐
│                AI/ML Application Architecture                   │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│  PostgreSQL +   │ Vector search, metadata                       │
│  pgvector       │                                               │
│  OR Pinecone    │ Managed vector database (production)          │
│  Redis          │ Feature store caching                         │
│  S3             │ Model artifacts, training data                │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

---

## Summary: Key Takeaways

### Database Selection Principles

```
1. MATCH DATA MODEL TO DOMAIN
   └── Don't force relational for graph data or vice versa

2. DESIGN FOR QUERY PATTERNS
   └── Choose database optimized for your access patterns

3. CONSIDER SCALABILITY EARLY
   └── Changing databases later is expensive

4. START SIMPLE, EVOLVE
   └── Don't over-engineer; add complexity as needed

5. OPERATIONAL REALITY MATTERS
   └── Factor in team expertise and operational burden

6. PLAN FOR FAILURE
   └── High availability and disaster recovery from day one

7. SECURITY IS NON-NEGOTIABLE
   └── Defense in depth across all layers

8. MONITOR AND OPTIMIZE
   └── Performance is an ongoing concern, not one-time
```

### Quick Reference: Database Types

| Type | When to Use | When to Avoid |
|------|-------------|---------------|
| **Relational** | Complex queries, ACID, well-defined schema | Rapid schema changes, extreme scale |
| **Document** | Flexible schema, hierarchical data | Many-to-many relationships |
| **Key-Value** | Simple lookups, caching | Complex queries |
| **Wide-Column** | Write-heavy, time-series | Ad-hoc queries |
| **Graph** | Relationship traversal | Simple CRUD |
| **Time-Series** | Timestamped data | Non-temporal data |
| **Vector** | AI/ML similarity search | Traditional queries |

---

## Resources

### Books
- *Designing Data-Intensive Applications* - Martin Kleppmann
- *Database Internals* - Alex Petrov
- *NoSQL Distilled* - Pramod Sadalage & Martin Fowler

### Websites
- [DB-Engines Ranking](https://db-engines.com/en/ranking)
- [Use The Index, Luke](https://use-the-index-luke.com/)
- [Database of Databases](https://dbdb.io/)

### Tools
- **Schema Design**: dbdiagram.io, ERDPlus
- **Benchmarking**: sysbench, YCSB, pgbench
- **Monitoring**: DataDog, Prometheus, pgAdmin
- **Migration**: Flyway, Liquibase, AWS DMS

---

*Last Updated: February 2026*
