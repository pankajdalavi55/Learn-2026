# SQL Database Architecture & Design Guide

## Table of Contents
1. [Architectural Foundations](#architectural-foundations)
2. [SQL Database Selection](#sql-database-selection)
3. [Database Design Principles](#database-design-principles)
4. [Schema Design Patterns](#schema-design-patterns)
5. [Scaling Strategies](#scaling-strategies)
6. [High Availability & Disaster Recovery](#high-availability--disaster-recovery)
7. [Performance Optimization](#performance-optimization)
8. [Microservices Database Patterns](#microservices-database-patterns)
9. [Caching & Query Optimization](#caching--query-optimization)
10. [Operational Considerations](#operational-considerations)

---

## Architectural Foundations

### Core Architecture Principles

**1. ACID Properties**
```
Atomicity: All-or-nothing transactions
Consistency: Data integrity enforced by constraints
Isolation: Concurrent transaction safety
Durability: Data persistence after commit
```

**2. CAP Theorem Application**
- SQL databases: **Consistency + Availability**
- Trade-off: Partition tolerance limited in traditional ACID setup
- Modern solutions: Add replication for partition tolerance

**3. Architectural Layers**
```
Application Layer
    ↓
ORM Layer (Hibernate, JPA, SQLAlchemy)
    ↓
Connection Pool (HikariCP, C3P0)
    ↓
Query Optimization Layer
    ↓
Storage Engine
    ↓
Physical Disk
```

### Key Architectural Decisions

| Decision | Factor | Consideration |
|----------|--------|---------------|
| **Vertical vs Horizontal** | Read/Write patterns | Sharding complexity vs hardware costs |
| **Replication Strategy** | Availability needs | Master-slave, multi-master trade-offs |
| **Normalization Level** | Query performance | 3NF vs denormalization for analytics |
| **Data Consistency** | Business requirements | Strong vs eventual consistency |

---

## SQL Database Selection

### PostgreSQL vs MySQL vs SQL Server vs Oracle

#### PostgreSQL: Advanced ACID Guarantees
**Best for:** Complex queries, advanced features, open-source
```yaml
Pros:
  - Full ACID compliance
  - Advanced data types (JSON, arrays, ranges)
  - Powerful full-text search
  - Multi-version concurrency control (MVCC)
  - Extensibility through custom functions
  - Strong data integrity constraints
Cons:
  - Slightly slower for simple read-heavy workloads
  - Larger memory footprint
  - Smaller ecosystem than MySQL
Ideal Use Cases:
  - Complex reporting systems
  - Data warehousing
  - Applications requiring advanced SQL features
  - Systems with complex business logic in database
```

#### MySQL 8.0+: Performance & Simplicity
**Best for:** Web applications, read-heavy systems, distributed setup
```yaml
Pros:
  - High throughput for simple queries
  - Excellent replication support
  - Lower memory footprint
  - Wide hosting support
  - Fast for read-heavy workloads
Cons:
  - Limited to basic data types
  - Less advanced query optimization
  - Window functions added recently
Ideal Use Cases:
  - Social media platforms
  - E-commerce systems
  - Content management systems
  - High-volume read scenarios
```

#### SQL Server: Enterprise Grade
**Best for:** Enterprise applications, Windows-heavy environments
```yaml
Pros:
  - Advanced analytics and reporting (SSIS, SSRS)
  - Excellent performance monitoring
  - Strong security features
  - Great integration with .NET ecosystem
Cons:
  - Expensive licensing
  - Windows/cloud-focused
  - Vendor lock-in
Ideal Use Cases:
  - Enterprise applications
  - Business intelligence systems
  - Systems with complex compliance needs
```

#### Oracle: Mission-Critical Systems
**Best for:** Highly demanding applications with unlimited budget
```yaml
Pros:
  - Extreme scalability and reliability
  - Advanced partitioning and compression
  - Sophisticated security model
  - Unmatched performance tuning capabilities
Cons:
  - Extremely expensive
  - Steep learning curve
  - Complex administration
Ideal Use Cases:
  - Banking and financial systems
  - Large multinational enterprises
  - Systems processing millions of transactions daily
```

### Selection Matrix

```
               Volume  Complexity  Cost  Speed  Open-Source
PostgreSQL      High     Very High  Low   Good       ✓
MySQL           Very High Medium    Low   Excellent  ✓
SQLite          Low      Low        Free  Good       ✓
SQL Server      High      High      High  Excellent  ✗
Oracle          Very High Very High  Very High Excellent ✗
```

---

## Database Design Principles

### 1. Normalization vs Denormalization

#### Normalization Levels
```
1NF (First Normal Form)
  - Atomic values only
  - No repeating groups
  Example: Customer(id, name, email) ✓
  Bad: Customer(id, name, emails[])

2NF (Second Normal Form)
  - Must be in 1NF
  - No partial dependencies (non-key attributes depend on all keys)
  Example: Order(order_id, product_id, quantity, price)
  Bad: Order(order_id, product_id, manufacturer) when multiple products per manufacturer

3NF (Third Normal Form)
  - Must be in 2NF
  - No transitive dependencies
  Example: Employee(emp_id, name, dept_id) + Department(dept_id, dept_name)
  Bad: Employee(emp_id, name, dept_name, location) - location is determined by dept_name
```

#### Normalization vs Denormalization Trade-offs
```
Normalized (3NF):
  ✓ Update efficiency (single source of truth)
  ✓ Storage efficiency
  ✓ Referential integrity maintained
  ✗ Complex joins (performance overhead)
  ✗ Higher query latency

Denormalized:
  ✓ Faster reads (fewer joins)
  ✓ Simpler queries
  ✗ Data duplication
  ✗ Update anomalies (multiple places to update)
  ✗ Higher risk of inconsistency
```

#### Strategic Approach
```
Rule: Normalize by default, denormalize for performance

Denormalization Strategy:
1. Design normalized schema
2. Identify query bottlenecks through profiling
3. Selectively denormalize (not entire schema)
4. Use materialized views or caching instead when possible
5. Implement consistency mechanisms (triggers, application logic)

Example: Product Reviews System
Normalized:
  - Product(product_id, name, price)
  - Review(review_id, product_id, rating, text)
  - Average rating: SELECT AVG(rating) FROM Review WHERE product_id=?

Denormalized for performance:
  - Product(product_id, name, price, avg_rating, review_count)
  - Review(review_id, product_id, rating, text)
  - Maintain avg_rating with triggers or batch updates
```

### 2. Domain-Driven Design (DDD) for Databases

```
Aggregate Design:
- Entity (has identity): Product, Order, Customer
- Value Object (no identity): Address, Money, Rating
- Aggregate Root: Entity responsible for consistency
  Example: Order (root) contains OrderItems, shipment tracking

Repository Pattern:
- Interface: OrderRepository
  save(order): void
  findById(id): Order
  findByCustomerId(customerId): List<Order>
- Implementation: SQL queries encapsulated

Bounded Contexts:
- Order Service DB: Orders, OrderItems, Shipments
- Inventory Service DB: Products, Stock, Reservations
- Separate databases for different bounded contexts
```

### 3. Entity-Relationship Model

**Core Concepts:**
```
Entity: Real-world thing (Customer, Product, Order)
Attribute: Property of entity (name, email, price)
Relationship: Association between entities
  - One-to-One: Person ↔ Passport
  - One-to-Many: Customer ← Orders
  - Many-to-Many: Students ↔ Courses (junction table)
```

---

## Schema Design Patterns

### 1. Time-Series Pattern
**Use case:** Monitoring, metrics, events, analytics
```sql
-- Core time-series table
CREATE TABLE metrics (
  metric_id BIGINT PRIMARY KEY,
  metric_name VARCHAR(255) NOT NULL,
  timestamp BIGINT NOT NULL,
  value DOUBLE NOT NULL,
  tags JSONB,  -- PostgreSQL JSON for flexible tagging
  INDEX (metric_name, timestamp)
);

-- Time-partitioned approach
CREATE TABLE metrics_2024_01 PARTITION OF metrics
  FOR VALUES FROM ('2024-01-01') TO ('2024-02-01');

-- Bucketing for time aggregation
SELECT 
  metric_name,
  DATE_TRUNC('hour', to_timestamp(timestamp)) as hour,
  AVG(value) as avg_value
FROM metrics
GROUP BY metric_name, hour;
```

### 2. Soft Delete Pattern
**Use case:** Audit trail, data recovery, compliance
```sql
CREATE TABLE users (
  user_id BIGINT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  deleted BOOLEAN DEFAULT FALSE,
  deleted_at TIMESTAMP,
  INDEX (deleted)  -- Important for filtering
);

-- View for active records
CREATE VIEW active_users AS
SELECT * FROM users WHERE deleted = FALSE;

-- Restore mechanism
UPDATE users SET deleted = FALSE WHERE user_id = ?;
```

### 3. Audit Log Pattern
**Use case:** Compliance, debugging, version history
```sql
CREATE TABLE user_audit (
  audit_id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  action VARCHAR(50) NOT NULL,  -- INSERT, UPDATE, DELETE
  old_values JSON,
  new_values JSON,
  changed_by VARCHAR(255),
  changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  INDEX (user_id, changed_at)
);

-- Trigger (PostgreSQL)
CREATE TRIGGER user_audit_trigger
AFTER UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION audit_user_changes();

FUNCTION audit_user_changes() RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO user_audit (user_id, action, old_values, new_values, changed_at)
  VALUES (NEW.user_id, 'UPDATE', row_to_json(OLD), row_to_json(NEW), NOW());
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```

### 4. Hierarchical Data Pattern
**Use case:** Organizational structures, file systems, categories
```sql
-- Adjacency List (simple but poor recursion performance)
CREATE TABLE categories (
  category_id INT PRIMARY KEY,
  parent_id INT,
  name VARCHAR(255),
  FOREIGN KEY (parent_id) REFERENCES categories(category_id)
);

-- Nested Set Model (efficient reads, complex writes)
CREATE TABLE categories_nested (
  category_id INT PRIMARY KEY,
  lft INT NOT NULL,
  rgt INT NOT NULL,
  name VARCHAR(255),
  UNIQUE (lft, rgt)
);

-- Query all descendants of Electronics
SELECT * FROM categories_nested
WHERE lft > @left AND rgt < @right;

-- Materialized Path (best for most cases)
CREATE TABLE categories_path (
  category_id INT PRIMARY KEY,
  path VARCHAR(255),  -- '1/5/12/'
  name VARCHAR(255)
);

-- Query all descendants
SELECT * FROM categories_path WHERE path LIKE '/1/5/%';
```

### 5. JSONB Storage Pattern (PostgreSQL)
**Use case:** Semi-structured data, flexible schemas
```sql
CREATE TABLE products (
  product_id BIGINT PRIMARY KEY,
  name VARCHAR(255),
  metadata JSONB NOT NULL,  -- Flexible attribute storage
  created_at TIMESTAMP,
  CONSTRAINT valid_json CHECK (jsonb_typeof(metadata) = 'object')
);

-- Insert flexible data
INSERT INTO products VALUES
(1, 'Laptop', '{"brand": "Dell", "cpu": "i9", "ram": 32}'),
(2, 'Phone', '{"brand": "Apple", "model": "Pro", "storage": 256}');

-- Query JSON data
SELECT * FROM products WHERE metadata ->> 'brand' = 'Dell';

-- GIN Index for fast JSONB queries
CREATE INDEX idx_products_metadata ON products USING GIN (metadata);
```

---

## Scaling Strategies

### 1. Vertical Scaling
**Approach:** Increase hardware resources (CPU, RAM, Storage)
```
Pros:
  - Simpler to implement
  - No application changes
  - Works well until hardware limits

Cons:
  - Expensive per unit performance gain
  - Hardware limits (can't scale infinitely)
  - Single point of failure
  - Downtime for upgrades

When to use:
  - Initial phase of growth
  - Temporary spikes
  - When horizontal scaling isn't feasible
```

### 2. Read Replicas (Horizontal Replication)
**Approach:** Multiple read-only copies for distributed reads
```
Master ──(write)──→ [Primary Database]
                    ↓ (replication)
                  [Read Replica 1]
                  [Read Replica 2]
                  [Read Replica 3]

Architecture:
Pros:
  - Distributes read load
  - Improved read throughput
  - Enables geographic distribution
  - Provides failover capability

Cons:
  - Replication lag (eventual consistency)
  - Increased storage usage
  - Complex failover management
  - Write bottleneck at master

Configuration (MySQL):
-- Master
SET GLOBAL binlog_format = 'ROW';
CREATE USER 'replication'@'slave.ip' IDENTIFIED BY 'password';
GRANT REPLICATION SLAVE ON *.* TO 'replication'@'slave.ip';

-- Slave
CHANGE MASTER TO MASTER_HOST = 'master.ip',
                   MASTER_USER = 'replication',
                   MASTER_PASSWORD = 'password',
                   MASTER_LOG_FILE = 'mysql-bin.000001',
                   MASTER_LOG_POS = 154;
START SLAVE;
```

### 3. Sharding (Horizontal Partitioning)
**Approach:** Distribute data across multiple databases by key
```
Sharding Strategies:

1. Range-based Sharding:
   user_id 1-1000000    → Shard 1
   user_id 1000001-2000000 → Shard 2
   user_id 2000001-3000000 → Shard 3
   Problem: Uneven distribution (hotspots)

2. Hash-based Sharding:
   shard_id = hash(user_id) % number_of_shards
   Pros: Even distribution
   Cons: Difficult resharding

3. Directory-based Sharding:
   UserShardMap(user_id → shard_id)
   Pros: Flexible, easy to reshard
   Cons: Additional lookup overhead

4. Geographic Sharding:
   Users in APAC → Shard in Singapore
   Users in EMEA → Shard in London
   Users in Americas → Shard in New York

Architecture:
[Application] → [Shard Router] → [Shard 1: DB1]
                                [Shard 2: DB2]
                                [Shard 3: DB3]

Pros:
  - Horizontal scalability
  - Improved performance through data isolation
  - Enables geographic distribution

Cons:
  - Complex application logic
  - Cross-shard queries expensive
  - Distributed transactions difficult
  - Resharding is painful

Best Practices:
- Choose shard key wisely (immutable, high-cardinality)
- Use consistent hashing for flexibility
- Implement shard-aware ORM layer
- Plan for resharding from day 1
- Monitor shard balance

Code Example:
class ShardRouter:
    def __init__(self, num_shards):
        self.num_shards = num_shards
    
    def get_shard_id(self, user_id):
        return hash(user_id) % self.num_shards
    
    def get_shard_connection(self, user_id):
        shard_id = self.get_shard_id(user_id)
        return self.shard_connections[shard_id]
    
    def query(self, user_id, sql):
        connection = self.get_shard_connection(user_id)
        return connection.execute(sql)
```

### 4. Partitioning (Within Single Database)
**Approach:** Divide table into smaller chunks based on key range
```sql
-- Range Partitioning (time-based)
CREATE TABLE sales (
  sale_id BIGINT,
  sale_date DATE,
  amount DECIMAL(10,2)
) PARTITION BY RANGE (YEAR(sale_date)) (
  PARTITION p2023 VALUES LESS THAN (2024),
  PARTITION p2024 VALUES LESS THAN (2025),
  PARTITION p2025 VALUES LESS THAN (2026),
  PARTITION pmax VALUES LESS THAN MAXVALUE
);

-- Hash Partitioning (even distribution)
CREATE TABLE logs (
  log_id BIGINT,
  message TEXT
) PARTITION BY HASH(log_id) PARTITIONS 4;

-- List Partitioning (categorical)
CREATE TABLE regions (
  region_id INT,
  region_name VARCHAR(255)
) PARTITION BY LIST(region_id) (
  PARTITION pna VALUES IN (1,2,3),
  PARTITION peurope VALUES IN (4,5,6),
  PARTITION papac VALUES IN (7,8,9)
);

Benefits:
- Improved query performance (partition pruning)
- Easier maintenance (drop old partitions)
- Better index management
- Parallel query execution
```

### 5. Scaling Decision Tree

```
Read Optimization Needed?
├─ Yes, single database size < 1TB
│  └─ Use Read Replicas
├─ Yes, need geographic distribution
│  └─ Use Read Replicas + CDN for data
└─ No, or read replicas insufficient

Write Optimization Needed?
├─ Yes, high volume of writes
│  └─ Consider Sharding or Queue + Async writes
├─ No, current master sufficient
│  └─ Optimize schema, indexes, queries

Data Size?
├─ < 1TB, < 100K ops/sec
│  └─ Vertical scaling + optimization
├─ 1-10TB, 100K-1M ops/sec
│  └─ Read replicas + partitioning
└─ > 10TB, > 1M ops/sec
   └─ Sharding required
```

---

## High Availability & Disaster Recovery

### 1. Master-Slave Replication with Failover
```
Active-Passive Setup:
┌─────────────────┐
│  Primary (Master)│ ← All writes
└────────┬────────┘
         │ (replication)
         ↓
┌─────────────────┐
│ Secondary (Slave)│ ← Failover target
└─────────────────┘

Failover Detection (Sentinel/Orchestration):
1. Primary health check fails (heartbeat, connection attempt)
2. Promote Secondary to Primary
3. Redirect application writes to new Primary
4. Bring down old Primary for investigation
5. Resync old Primary from new Primary

Configuration (MySQL + Orchestrator):
- Uses orchestrator for automatic failover detection
- No manual intervention needed
- Preserves data consistency
```

### 2. Multi-Master Replication
```
Circular Replication:
Master 1 ←→ Master 2 ←→ Master 3 ↻

Pros:
  - Active-active writes (any node accepts writes)
  - Higher availability
  - Geographic distribution

Cons:
  - Complex conflict resolution
  - Write conflicts possible
  - Higher latency
  - Operational complexity

Conflict Resolution Strategies:
1. Last-write-wins: Timestamp-based
2. Custom application logic: Domain knowledge
3. Version vectors: Causality tracking
4. Quorum-based: Majority decides

Use Cases:
- Distributed systems across data centers
- Always-on geographically distributed systems
```

### 3. Backup & Recovery Strategy

**Backup Types:**
```
Full Backup:
- Complete database copy
- Largest size
- Fastest restore
- Frequency: Weekly

Incremental Backup:
- Changes since last backup
- Smaller size
- Slower restore (need full + incrementals)
- Frequency: Daily

Differential Backup:
- Changes since last full backup
- Medium size
- Faster restore than incremental
- Frequency: Daily

Transaction Log Backup (PITR):
- Enable point-in-time recovery
- Frequency: Every 15 minutes or continuous
```

**Backup Strategy:**
```
Daily Schedule:
02:00 - Full backup (weekly)
04:00 - Incremental backup (M-S)
06:00-23:00 - Transaction log backups (every 15 min)

Retention Policy:
- Daily: Keep 7 days
- Weekly: Keep 4 weeks
- Monthly: Keep 12 months

Testing:
- Monthly restore test to non-production
- Measure RTO (Recovery Time Objective)
- Measure RPO (Recovery Point Objective)

Implementation (PostgreSQL):
-- Enable WAL archiving
wal_level = replica
archive_mode = on
archive_command = 'cp %p /backup/wal_archive/%f'

-- Backup command
pg_basebackup -h localhost -U backup_user -D /backup/full_backup -X stream -P
```

### 4. Disaster Recovery Plan

**RTO & RPO Definition:**
```
RTO (Recovery Time Objective): Maximum allowable downtime
- Critical systems: 15 minutes
- Important systems: 1 hour
- Non-critical: 4 hours

RPO (Recovery Point Objective): Maximum data loss acceptable
- Financial systems: None (0 minutes)
- Transactional: 5 minutes
- Analytical: 1 hour

DR Tier Selection:
Tier 1 - Hot Standby: Fully operational duplicate (RTO: 0-15 min)
Tier 2 - Warm Standby: Partially configured (RTO: 15 min - 1 hour)
Tier 3 - Cold Standby: Configuration available offsite (RTO: 4+ hours)
```

---

## Performance Optimization

### 1. Indexing Strategy

**Index Types:**
```
Single Column Index:
CREATE INDEX idx_users_email ON users(email);
- Best for: Equality searches, WHERE email = ?

Composite Index (Multi-column):
CREATE INDEX idx_orders_customer_date ON orders(customer_id, order_date DESC);
- Best for: Multiple column filters, range queries
- Query: WHERE customer_id = ? AND order_date > ?

Partial Index:
CREATE INDEX idx_active_users ON users(email) WHERE deleted = FALSE;
- Only indexes active records
- Smaller index size
- Faster for active record queries

Full-Text Index (PostgreSQL):
CREATE INDEX idx_products_search ON products USING GIN(to_tsvector('english', description));
- Full-text search queries
- Language-aware stemming

JSONB Index (PostgreSQL):
CREATE INDEX idx_metadata ON products USING GIN(metadata);
- JSONB column queries
- Path-specific index: CREATE INDEX idx_metadata_brand ON products((metadata->'brand'))

Spatial Index:
CREATE INDEX idx_location ON restaurants USING GIST(location);
- Geographic queries (within radius, etc.)

Covering Index:
CREATE INDEX idx_orders_customer_id_total ON orders(customer_id) INCLUDE (order_total);
- Index includes additional columns (no table lookup needed)
```

**Index Design Rules:**
```
1. Index foreign keys
   CREATE INDEX idx_orders_customer_id ON orders(customer_id);

2. Index columns in WHERE clauses
   For query: WHERE status = ? AND created_at > ?
   CREATE INDEX idx_status_created ON orders(status, created_at);

3. Order columns wisely (equality before range)
   Good: INDEX (status, created_at) for WHERE status = ? AND created_at > ?
   Bad:  INDEX (created_at, status)

4. Avoid indexing low-cardinality columns
   Bad: INDEX ON gender (only M/F/Other)
   OK: INDEX ON user_type (more variety)

5. Monitor index usage
   PostgreSQL: pg_stat_user_indexes
   MySQL: INFORMATION_SCHEMA.STATISTICS

6. Drop unused indexes
   Unused indexes slow down writes without benefiting reads
```

**Index Cost Analysis:**
```
Pros:
  - 10-100x faster reads (for good indexes)
  - Enable efficient joins
  - Support range queries

Cons:
  - 5-30% slower writes (maintain index)
  - Insert/update latency increased
  - Higher memory usage
  - Storage overhead

Trade-off: Balance read vs write performance
- Read-heavy: More indexes
- Write-heavy: Fewer indexes
- Mixed: Strategic indexes on important queries
```

### 2. Query Optimization

**Execution Plan Analysis:**
```sql
-- PostgreSQL EXPLAIN
EXPLAIN (ANALYZE, BUFFERS) 
SELECT o.*, c.name
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
WHERE o.order_date > '2024-01-01';

Output Analysis:
- Seq Scan (sequential table scan): Slow, full scan
- Index Scan: Fast, uses index
- Bitmap Index Scan: Medium, efficient for range queries
- Nested Loop: Good for small tables
- Hash Join: Good for joining larger sets
- Sort: Check if ORDER BY can use index

Cost Interpretation:
(actual time=2.5..15.3 rows=1000)
- First time to first row: 2.5ms
- Total time: 15.3ms
- Returned rows: 1000
- Higher numbers indicate worse performance

MySQL EXPLAIN:
EXPLAIN 
SELECT o.*, c.name
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
WHERE o.order_date > '2024-01-01';

Key findings:
- type: ALL (full scan), index, range
- possible_keys: Which indexes could be used
- key: Which index was actually used
- rows: Estimated rows examined (should be small)
```

**Query Optimization Techniques:**

```sql
1. Selective Projection (fetch only needed columns)
Bad:
SELECT * FROM orders;  -- Fetches all columns

Good:
SELECT order_id, customer_id, total 
FROM orders WHERE status = 'completed';

2. Use LIMIT and PAGINATION
Bad:
SELECT * FROM products;  -- Millions of rows

Good:
SELECT * FROM products LIMIT 50 OFFSET 100;  -- Page 3 with 50 per page

3. Filter early (push WHERE down)
Bad:
SELECT SUM(amount) FROM orders 
INNER JOIN order_items ON orders.order_id = order_items.order_id;

Good:
SELECT SUM(amount) FROM order_items
WHERE order_id IN (SELECT order_id FROM orders WHERE status = 'completed');

4. Avoid OR in complex queries
Bad:
SELECT * FROM orders 
WHERE customer_id = 1 OR customer_id = 2 OR customer_id = 3;

Good:
SELECT * FROM orders WHERE customer_id IN (1, 2, 3);

5. Avoid functions on indexed columns
Bad:
SELECT * FROM users WHERE LOWER(email) = 'test@example.com';
-- Can't use index on email

Good:
SELECT * FROM users WHERE email = 'test@example.com';
-- Uses index

6. Batch operations instead of individual queries
Bad:
for user_id in user_ids:
    INSERT INTO user_stats (user_id, score) VALUES (user_id, score);

Good:
INSERT INTO user_stats (user_id, score) 
VALUES (1, 100), (2, 200), (3, 300), (4, 400);
-- Single round trip, faster

7. Use UNION ALL instead of UNION
Bad:
SELECT email FROM users WHERE status = 'active'
UNION  -- Removes duplicates (slower)
SELECT email FROM users WHERE premium = true;

Good:
SELECT email FROM users WHERE status = 'active' OR premium = true;
-- Single query, no duplicate elimination overhead
```

### 3. Connection Pooling

**Concept:** Maintain pool of database connections for reuse
```
Without Pooling:
App Request → Create Connection → Execute Query → Close Connection → Response
(Slow: connection overhead for each request)

With Pooling:
Connection Pool (pre-created connections)
    ↓
App Request → Get from Pool → Execute Query → Return to Pool → Response
(Fast: reuse existing connections)

Configuration (HikariCP - Java):
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/db
    username: user
    password: pass
    hikari:
      maximum-pool-size: 10      # Max connections
      minimum-idle: 5            # Idle connections to maintain
      connection-timeout: 30000  # Timeout to get connection
      idle-timeout: 600000       # 10 min before closing idle
      max-lifetime: 1800000      # 30 min connection max age

Pool Sizing Formula:
pool_size = (core_count × 2) + (disk_spindle_count)
For 8 cores + SSD: pool_size = 16-20

Monitoring:
- Connection utilization (% of active vs pool size)
- Wait times for connection
- Connection churn rate
```

---

## Microservices Database Patterns

### 1. Database Per Service Pattern
```
User Service
  ├─ Database: user_db
  └─ Tables: users, user_profiles, permissions
  
Order Service
  ├─ Database: order_db
  └─ Tables: orders, order_items, shipments

Product Service
  ├─ Database: product_db
  └─ Tables: products, categories, inventory

Benefits:
  - Schema flexibility per service
  - Technology choice per service
  - Scalability independence
  - Fault isolation

Challenges:
  - No distributed transactions
  - Eventual consistency needed
  - Cross-service queries difficult
  - Data duplication across services
```

### 2. Event Sourcing Pattern
```
Instead of storing current state, store events that led to state

Traditional:
User Table: (user_id=1, name='John', email='john@test.com', status='active')

Event Sourcing:
Event Log:
  UserCreated(user_id=1, name='John', email='john@test.com') @t1
  UserEmailChanged(user_id=1, email='john.doe@test.com') @t2
  UserStatusChanged(user_id=1, status='active') @t3

Benefits:
  - Complete audit trail
  - Easy to debug (replay events)
  - Natural fit with event streaming
  - Temporal queries (state at time T)

Implementation:
```python
class Event:
    def __init__(self, event_type, aggregate_id, data, timestamp):
        self.event_type = event_type
        self.aggregate_id = aggregate_id
        self.data = data
        self.timestamp = timestamp

class EventStore:
    def append(self, event):
        # Store event in event_log table
        db.insert('event_log', {
            'aggregate_id': event.aggregate_id,
            'event_type': event.event_type,
            'data': json.dumps(event.data),
            'timestamp': event.timestamp
        })
    
    def get_events(self, aggregate_id):
        # Retrieve all events for aggregation
        return db.select('event_log', where={'aggregate_id': aggregate_id})

# Replay events to reconstruct current state
def reconstruct_user(user_id):
    events = event_store.get_events(user_id)
    user = User(user_id)
    for event in events:
        if event.event_type == 'UserCreated':
            user.name = event.data['name']
        elif event.event_type == 'UserEmailChanged':
            user.email = event.data['new_email']
    return user
```

### 3. CQRS (Command Query Responsibility Segregation)
```
Separate read and write models

Commands (Writes):
  CreateOrder → Write Model → Event Store → Projections

Queries (Reads):
  GetOrder ← Read Model (denormalized views)

Architecture:
Commands
  ├─ Create Order
  ├─ UpdateOrder
  └─ CancelOrder
     ↓
Event Store (append-only log)
     ↓
Projections (denormalized read models)
     ↓
Queries (fast reads)

Example:
-- Write Model (normalized)
CREATE TABLE orders (
  order_id BIGINT PRIMARY KEY,
  customer_id BIGINT,
  status VARCHAR(50),
  created_at TIMESTAMP
);

-- Read Model (denormalized projection)
CREATE TABLE order_projections (
  order_id BIGINT PRIMARY KEY,
  customer_name VARCHAR(255),
  customer_email VARCHAR(255),
  status VARCHAR(50),
  items_count INT,
  total_amount DECIMAL(10,2),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- When OrderCreated event occurs:
INSERT INTO order_projections
SELECT o.order_id, c.name, c.email, o.status, COUNT(*), SUM(i.price), o.created_at, NOW()
FROM orders o
JOIN customers c ON o.customer_id = c.customer_id
LEFT JOIN order_items i ON o.order_id = i.order_id
GROUP BY o.order_id;
```

### 4. Saga Pattern for Distributed Transactions
```
Orchestrated Saga:
OrderCreated Event
  ├─ Order Service: Create Order
  ├─ Payment Service: Process Payment (if order created)
  ├─ Inventory Service: Reserve Items (if payment successful)
  └─ Notification Service: Send email (if all successful)

If Payment fails:
  OrderCreated → Create Order → Payment Fails → Compensate: Cancel Order

Orchestrator coordinates across services:

    OrderSaga (Orchestrator)
       ↓
    OrderService: createOrder() ✓
       ↓
    PaymentService: processPayment() ✓
       ↓
    InventoryService: reserveItems() ✗ FAILED
       ↓
    Compensate:
       PaymentService.refund()
       OrderService.cancelOrder()

Choreography Saga:
Services listen to events and emit compensating events

OrderCreated → OrderService
OrderCreated → PaymentService processes, emits PaymentProcessed
PaymentProcessed → InventoryService, emits InventoryReserved
InventoryReserved → NotificationService, sends email

If InventoryReserved fails:
  InventoryFailed → PaymentService (compensate: refund)
  InventoryFailed → OrderService (compensate: cancel order)
```

---

## Caching & Query Optimization

### 1. Caching Layers

**Application-Level Caching:**
```java
// In-memory cache with TTL (minutes)
@Cacheable("products", key = "#productId")
public Product getProduct(Long productId) {
    return productRepository.findById(productId);  // Cache hit: skipped
}

// Cache eviction after update
@CacheEvict("products", key = "#productId")
public void updateProduct(Long productId, Product product) {
    productRepository.save(product);
}

// Cache configuration
@Configuration
@EnableCaching
public class CacheConfig {
    @Bean
    public CacheManager cacheManager() {
        return new ConcurrentMapCacheManager("products", "users", "orders");
    }
}
```

**Redis Caching:**
```
Architecture:
Request
  ├─ Check Redis Cache (fast, 1-10ms)
  ├─ If hit: Return cached value
  └─ If miss: Query Database → Cache result in Redis → Return

Redis Advantages:
  - In-memory (extremely fast: <1ms)
  - Distributed (shared across services)
  - TTL support (auto expiration)
  - Rich data structures (strings, sets, lists, sorted sets)

Configuration:
```yaml
spring:
  redis:
    host: redis.example.com
    port: 6379
    timeout: 2000ms
    jedis:
      pool:
        max-active: 8
        max-idle: 8
        min-idle: 0

# Cache with TTL
product.cache.ttl=3600 # 1 hour
```

**Database Query Cache:**
```sql
-- MySQL Query Cache (deprecated in MySQL 5.7.20)
-- PostgreSQL doesn't have built-in query cache
-- Solution: Use Redis or application caching

-- Materialized Views (PostgreSQL)
CREATE MATERIALIZED VIEW product_summary AS
SELECT 
  category_id,
  COUNT(*) as total_products,
  AVG(price) as avg_price,
  MAX(price) as max_price
FROM products
GROUP BY category_id;

-- Refresh periodically
REFRESH MATERIALIZED VIEW product_summary;

-- Query
SELECT * FROM product_summary WHERE category_id = 5;
```

### 2. Cache Invalidation Strategies

**Time-Based (TTL):**
```
Cache expires after fixed duration
- Simple to implement
- Risk: Stale data (until expiration)
- Use for: Non-critical, slowly changing data

Example: Product prices cached for 10 minutes
When product updates: Wait max 10 minutes for consistency (slow)
```

**Event-Based Invalidation:**
```
Cache evicted when data changes
- Immediate consistency
- More complex logic
- Use for: Critical data, frequently changing

Example: When product updates
  ProductUpdated Event
    ├─ Update database
    └─ Invalidate product cache
       Cache expires immediately, next access rebuilds

Implementation:
@Service
public class ProductService {
  @Autowired
  private CacheManager cacheManager;
  
  public void updateProduct(Product product) {
    productRepository.save(product);
    // Invalidate cache immediately
    cacheManager.getCache("products").evict(product.getId());
  }
}
```

**LRU Cache Eviction:**
```
Remove least recently used items when capacity full
- Balanced approach
- Good for bounded caches

Configuration:
```yaml
spring:
  cache:
    type: simple
    cache-names: products, users
    # JVM heap memory limit
    max-size: 1000
```

**Cache Warming:**
```
Pre-populate cache at startup

@ApplicationReady
public void warmCache() {
  // Load frequently accessed data on startup
  List<Product> topProducts = productRepository.findTopSelling(100);
  topProducts.forEach(p -> cache.put(p.getId(), p));
}

Benefits:
  - No cold start (no slow first requests)
  - Predictable performance
  - Useful for critical data
```

### 3. Statistical Query Optimization

```sql
-- Pre-compute aggregate tables (OLAP)
CREATE TABLE daily_sales_summary (
  date DATE PRIMARY KEY,
  total_sales DECIMAL(15,2),
  order_count INT,
  customer_count INT,
);

-- Scheduled job (nightly)
INSERT INTO daily_sales_summary
SELECT 
  DATE(order_date) as date,
  SUM(total) as total_sales,
  COUNT(*) as order_count,
  COUNT(DISTINCT customer_id) as customer_count
FROM orders
WHERE order_date = CURDATE();

-- Fast query (milliseconds)
SELECT SUM(total_sales) FROM daily_sales_summary 
WHERE date BETWEEN '2024-01-01' AND '2024-12-31';
-- vs scanning billions of order rows

Window Functions (analytical queries):
SELECT 
  order_id,
  customer_id,
  amount,
  SUM(amount) OVER (PARTITION BY customer_id ORDER BY order_date) as running_total,
  ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY order_date) as order_seq
FROM orders;
```

---

## Operational Considerations

### 1. Monitoring & Observability

**Key Metrics:**
```
Performance Metrics:
  - Query latency (p50, p95, p99)
  - Throughput (queries per second)
  - Connection count and utilization
  - Lock contention
  - Replication lag (for replicas)

Health Metrics:
  - CPU utilization
  - Memory usage
  - Disk I/O and space
  - Network bandwidth
  - Error rates

Queries to Monitor:
  - Slow queries (>100ms)
  - Cache hit ratio
  - Table scan frequency
  - Long-running transactions
```

**Monitoring Tools:**
```
PostgreSQL:
  - pg_stat_statements (tracks query statistics)
  - pgAdmin (web-based admin)
  - pg_stat_user_indexes (index usage)

MySQL:
  - Performance Schema
  - MySQL Workbench
  - Percona Monitoring and Management (PMM)

Query Examples:
-- PostgreSQL: Slow queries
SELECT query, mean_exec_time, calls FROM pg_stat_statements
ORDER BY mean_exec_time DESC LIMIT 10;

-- MySQL: Slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL slow_query_log_file = '/var/log/mysql/slow.log';
SET GLOBAL long_query_time = 1;  -- 1 second

-- Index utilization check
SELECT schemaname, tablename, indexname, idx_scan
FROM pg_stat_user_indexes
ORDER BY idx_scan;
```

### 2. Migration & Schema Evolution

**Online Schema Changes:**
```
Challenge: Changing schema on live database without downtime

Strategy 1: Expand-Contract Pattern
1. Expand: Add new column/table
   ALTER TABLE users ADD COLUMN email_v2 VARCHAR(255);

2. Dual-write: Application writes to both old and new
   INSERT INTO users (email, email_v2) VALUES (?, ?);

3. Backfill: Copy existing data to new column
   UPDATE users SET email_v2 = email WHERE email_v2 IS NULL;

4. Validate: Verify data consistency
   SELECT COUNT(*) FROM users WHERE email != email_v2 AND email_v2 IS NOT NULL;

5. Switch: Read from new column
   Application changes: SELECT email_v2 FROM users;

6. Contract: Remove old column
   ALTER TABLE users DROP COLUMN email;

Strategy 2: Use tools like Pt-online-schema-change
pt-online-schema-change \
  --alter "ADD COLUMN age INT" \
  D=mydb,t=users \
  --execute

Benefits: Minimal locking, online changes
```

**Backward-Compatible Changes:**
```
Good Changes (no downtime):
  - Add nullable column
  - Add column with default value
  - Add new table
  - Rename column (with backward compatibility layer)
  - Add index (takes time but non-blocking)

Bad Changes (require downtime):
  - Drop column/table
  - Change column type
  - Add NOT NULL column without default
  - Add highly selective index during peak traffic

Change Deployment Order:
1. Deploy new application version (supports old and new schema)
2. Make schema changes
3. Old application still works (uses old schema)
4. Cutover: Remove backward compatibility code

Versioning Pattern:
-- Schema version tracking
CREATE TABLE schema_version (
  version INT PRIMARY KEY,
  applied_at TIMESTAMP,
  description VARCHAR(255)
);

-- Versioned columns for large changes
ALTER TABLE users ADD COLUMN data_v2 JSONB;  -- New data format
-- Old application uses data column
-- New application uses data_v2 column
-- After cutover: drop data column
```

### 3. Capacity Planning

**Growth Forecasting:**
```
Data Growth:
  Current: 100GB
  Monthly growth: 5%
  Forecast: 100 * (1.05)^12 = 180GB after 1 year
  
Estimated capacity needed: 200GB (with 10% buffer)

Query Volume:
  Current: 1M queries/day
  Growth rate: 20% YoY
  
Resource Requirements:
  - Storage: Planning for 2-3 years ahead
  - Memory: Increase cache pool for growing dataset
  - Connection pool: Scale with application instances
  - Network: Monitor bandwidth utilization
```

**Capacity Timeline:**
```
Immediate (0-3 months):
  - Monitor current metrics
  - Identify optimization opportunities
  - Add indexes on slow queries

Short-term (3-12 months):
  - Upgrade hardware if needed (vertical scaling)
  - Implement caching strategy
  - Archive old data if applicable

Long-term (1-2 years):
  - Plan database sharding if needed
  - Setup read replicas for geographic distribution
  - Migrate to managed database service if operational burden too high
```

### 4. Database Security

**Access Control:**
```sql
-- Principle of least privilege
-- Create roles with minimal required permissions

-- Admin role (full access, for maintenance)
CREATE ROLE admin_user LOGIN PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE maindb TO admin_user;

-- Application role (limited access)
CREATE ROLE app_user LOGIN PASSWORD 'app_password';
GRANT CONNECT ON DATABASE maindb TO app_user;
GRANT USAGE ON SCHEMA public TO app_user;
GRANT SELECT, INSERT, UPDATE ON orders TO app_user;
GRANT SELECT ON products TO app_user;

-- Read-only role
CREATE ROLE readonly_user LOGIN PASSWORD 'readonly_password';
GRANT SELECT ON products, orders TO readonly_user;

-- Analytics role (aggregates only, no individual records)
CREATE ROLE analytics_user LOGIN PASSWORD 'analytics_password';
GRANT USAGE ON SCHEMA public TO analytics_user;
-- No direct table access, only views with aggregated data
```

**Encryption:**
```
At Rest:
  - Full disk encryption (OS level: BitLocker, dm-crypt)
  - Transparent Data Encryption (TDE) for sensitive columns
  - Backup encryption

In Transit:
  - TLS/SSL for database connections
  - Encrypted connections from application to database
  - VPN for replication traffic

Example (PostgreSQL with SSL):
ssl = on
ssl_cert_file = '/etc/postgresql/server.crt'
ssl_key_file = '/etc/postgresql/server.key'

Connection string: postgres://user:pass@host:5432/db?sslmode=require

Application-Level Encryption:
- Sensitive data (SSN, credit cards) encrypted before storage
- Encryption keys stored separately
- Decrypt only when needed
```

**Audit & Compliance:**
```
Enable logging:
PostgreSQL:
  log_connections = on
  log_disconnections = on
  log_statement = 'all'  -- or 'mod' for selective logging

Log Analysis:
  - Failed authentication attempts
  - Operations by sensitive tables
  - Unusual access patterns
  - Data exports

GDPR/CCPA Compliance:
  - Right to be forgotten: Soft delete + secure purge
  - Data portability: Export functionality
  - Audit trail of who accessed what
  - Data minimization: Don't store unnecessary data
```

---

## Summary: Architecture Decision Matrix

| Scenario | Recommended Approach | Database | Pattern |
|----------|---------------------|----------|---------|
| **Startup, <10M records** | Single PostgreSQL | PostgreSQL | Normalized, vertical scaling |
| **Social media, read-heavy** | Read replicas + Cache | MySQL + Redis | Denormalized, read replicas |
| **E-commerce, high writes** | Sharded MySQL | MySQL | Database per service, sharding |
| **Financial system, must not lose data** | Multi-master + PITR | PostgreSQL/Oracle | Event sourcing, multi-region |
| **Analytics, massive data** | Data warehouse | DuckDB/Snowflake | Columnar, materialized views |
| **Microservices architecture** | Separate DB per service | Varies | CQRS, event sourcing, sagas |
| **Real-time analytics** | Streaming DB | ClickHouse/Druid | Time-series partitioning |

---

## References & Resources

### Learning Resources
- PostgreSQL Official Documentation: https://www.postgresql.org/docs/
- MySQL 8.0 Reference: https://dev.mysql.com/doc/
- "Designing Data-Intensive Applications" by Martin Kleppmann
- SQL Performance Explained by SQL University
- "Database Design for Humans" by Emily Kausalya

### Tools & Platforms
- DBeaver (Visual database management)
- DataGrip (JetBrains database IDE)
- pgAdmin (PostgreSQL admin web interface)
- Percona Toolkit (MySQL toolset)
- Liquibase (Database migration and versioning)

---

**Last Updated:** March 4, 2026
**Version:** 1.0
**Status:** Complete Production Guide
