# SQL Database Design Interview Questions
## For 4-10 Years Experience

---

## Table of Contents
1. [Normalization & Denormalization](#1-normalization--denormalization)
2. [Schema Design Principles](#2-schema-design-principles)
3. [Indexing Strategies](#3-indexing-strategies)
4. [Keys & Constraints](#4-keys--constraints)
5. [Relationships & Cardinality](#5-relationships--cardinality)
6. [Performance Optimization](#6-performance-optimization)
7. [Data Modeling Scenarios](#7-data-modeling-scenarios)
8. [Transactions & ACID](#8-transactions--acid)
9. [Partitioning & Sharding](#9-partitioning--sharding)
10. [Real-World Design Problems](#10-real-world-design-problems)
11. [SQL Database Selection](#11-sql-database-selection)
12. [Situation-Based Interview Questions](#12-situation-based-interview-questions)

---

## 1. Normalization & Denormalization

### Q1: Explain the different Normal Forms with examples.

**Answer:**

| Normal Form | Rule | Example Violation |
|-------------|------|-------------------|
| **1NF** | Atomic values, no repeating groups | Column with comma-separated values |
| **2NF** | 1NF + No partial dependencies | Non-key column depends on part of composite key |
| **3NF** | 2NF + No transitive dependencies | Column A → Column B → Column C |
| **BCNF** | Every determinant is a candidate key | Non-candidate key determines another column |
| **4NF** | No multi-valued dependencies | Independent multi-valued facts in same table |
| **5NF** | No join dependencies | Table can be decomposed into smaller tables |

**Example - Normalizing an Order Table:**

```sql
-- Unnormalized (violates 1NF)
CREATE TABLE Orders_Bad (
    order_id INT,
    customer_name VARCHAR(100),
    products VARCHAR(500)  -- "Laptop, Mouse, Keyboard"
);

-- Normalized to 3NF
CREATE TABLE Customers (
    customer_id INT PRIMARY KEY,
    customer_name VARCHAR(100),
    email VARCHAR(100)
);

CREATE TABLE Orders (
    order_id INT PRIMARY KEY,
    customer_id INT REFERENCES Customers(customer_id),
    order_date DATE
);

CREATE TABLE Order_Items (
    order_item_id INT PRIMARY KEY,
    order_id INT REFERENCES Orders(order_id),
    product_id INT REFERENCES Products(product_id),
    quantity INT,
    unit_price DECIMAL(10,2)
);
```

---

### Q2: When would you choose to denormalize a database?

**Answer:**

**Denormalization Scenarios:**

1. **Read-Heavy Workloads** - When reads vastly outnumber writes
2. **Reporting/Analytics** - Dashboard queries needing aggregated data
3. **Performance Optimization** - Eliminate expensive JOINs
4. **Caching Frequently Accessed Data** - Store computed values

**Denormalization Techniques:**

```sql
-- Technique 1: Adding Redundant Columns
CREATE TABLE Orders (
    order_id INT PRIMARY KEY,
    customer_id INT,
    customer_name VARCHAR(100),  -- Redundant but avoids JOIN
    total_amount DECIMAL(10,2),  -- Pre-computed aggregate
    item_count INT               -- Pre-computed count
);

-- Technique 2: Materialized Views
CREATE MATERIALIZED VIEW monthly_sales AS
SELECT 
    DATE_TRUNC('month', order_date) as month,
    SUM(total_amount) as revenue,
    COUNT(*) as order_count
FROM Orders
GROUP BY DATE_TRUNC('month', order_date);

-- Technique 3: Summary Tables
CREATE TABLE Daily_Sales_Summary (
    summary_date DATE PRIMARY KEY,
    total_orders INT,
    total_revenue DECIMAL(15,2),
    avg_order_value DECIMAL(10,2),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Trade-offs:**
| Aspect | Normalized | Denormalized |
|--------|------------|--------------|
| Write Performance | Better | Worse (update anomalies) |
| Read Performance | More JOINs | Faster queries |
| Storage | Less | More |
| Data Consistency | Easier | Harder to maintain |
| Flexibility | High | Lower |

---

### Q3: What is the difference between BCNF and 3NF? When does it matter?

**Answer:**

**3NF allows:** Non-prime attributes can be determined by candidate keys
**BCNF requires:** EVERY determinant must be a candidate key

```sql
-- Example where 3NF ≠ BCNF
-- Table: Student_Course_Instructor
-- student_id, course, instructor

-- Functional Dependencies:
-- {student_id, course} → instructor
-- instructor → course (Each instructor teaches only one course)

-- This is in 3NF but NOT in BCNF because:
-- 'instructor' determines 'course' but 'instructor' is not a candidate key

-- BCNF Decomposition:
CREATE TABLE Instructor_Course (
    instructor_id INT PRIMARY KEY,
    course VARCHAR(100)
);

CREATE TABLE Student_Instructor (
    student_id INT,
    instructor_id INT,
    PRIMARY KEY (student_id, instructor_id)
);
```

**When it matters:** Large tables with complex dependencies where anomalies can cause data corruption.

---

## 2. Schema Design Principles

### Q4: How do you approach designing a database schema from scratch?

**Answer:**

**Step-by-Step Approach:**

```
1. Requirements Gathering
   ├── Identify entities (nouns)
   ├── Identify relationships (verbs)
   ├── Determine cardinality
   └── List business rules/constraints

2. Conceptual Design (ER Diagram)
   ├── Define entities and attributes
   ├── Identify primary keys
   └── Map relationships

3. Logical Design
   ├── Convert ER to tables
   ├── Apply normalization
   ├── Define data types
   └── Add constraints

4. Physical Design
   ├── Index strategy
   ├── Partitioning decisions
   ├── Storage considerations
   └── Performance tuning
```

**Example - E-commerce System:**

```sql
-- Step 1: Identify Entities
-- Users, Products, Orders, Categories, Reviews, Addresses

-- Step 2: Define Relationships
-- User (1) ────< (M) Orders
-- Order (1) ────< (M) Order_Items
-- Product (M) ────< (M) Categories (via junction table)

-- Step 3: Create Schema
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    
    INDEX idx_email (email),
    INDEX idx_created_at (created_at)
);

CREATE TABLE products (
    product_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    cost_price DECIMAL(10,2),
    stock_quantity INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_sku (sku),
    INDEX idx_price (price),
    FULLTEXT INDEX idx_name_desc (name, description)
);
```

---

### Q5: What are surrogate keys vs natural keys? When to use each?

**Answer:**

| Aspect | Natural Key | Surrogate Key |
|--------|-------------|---------------|
| **Definition** | Business-meaningful column | System-generated identifier |
| **Example** | SSN, Email, ISBN | AUTO_INCREMENT, UUID |
| **Stability** | Can change | Never changes |
| **Size** | Often larger | Typically INT/BIGINT |
| **JOINs** | Slower with composite | Faster single-column |

**When to Use Natural Keys:**
- Lookup tables (country_code, currency_code)
- When business identifier is stable and unique
- When external systems reference by natural key

**When to Use Surrogate Keys:**
- Primary entities (users, orders, products)
- When natural key can change (email, phone)
- When natural key is composite or large

```sql
-- Natural Key Example
CREATE TABLE countries (
    country_code CHAR(2) PRIMARY KEY,  -- ISO code as natural key
    country_name VARCHAR(100) NOT NULL
);

-- Surrogate Key Example
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY AUTO_INCREMENT,  -- Surrogate
    email VARCHAR(255) UNIQUE NOT NULL           -- Natural key as alternate
);

-- Best Practice: Use both
CREATE TABLE products (
    product_id BIGINT PRIMARY KEY AUTO_INCREMENT,  -- Surrogate for JOINs
    sku VARCHAR(50) UNIQUE NOT NULL,               -- Natural for business queries
    upc VARCHAR(12) UNIQUE                         -- Another natural identifier
);
```

---

### Q6: How do you handle soft deletes vs hard deletes?

**Answer:**

```sql
-- Soft Delete Implementation
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL,
    deleted_by BIGINT NULL,
    
    -- Unique constraint only on active records
    UNIQUE INDEX idx_email_active (email, is_deleted)
);

-- Query patterns with soft delete
-- Always exclude deleted records
SELECT * FROM users WHERE is_deleted = FALSE;

-- Create a view for convenience
CREATE VIEW active_users AS
SELECT * FROM users WHERE is_deleted = FALSE;

-- Soft delete operation
UPDATE users 
SET is_deleted = TRUE, 
    deleted_at = CURRENT_TIMESTAMP,
    deleted_by = @current_user_id
WHERE user_id = @target_user_id;
```

**Advanced Pattern - Temporal Tables:**

```sql
-- History table for audit trail
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    valid_from TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    valid_to TIMESTAMP NULL,
    is_current BOOLEAN DEFAULT TRUE
);

-- Archive deleted/modified records
CREATE TABLE users_history (
    history_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,
    email VARCHAR(255),
    name VARCHAR(100),
    valid_from TIMESTAMP,
    valid_to TIMESTAMP,
    operation ENUM('INSERT', 'UPDATE', 'DELETE'),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**Comparison:**

| Aspect | Soft Delete | Hard Delete |
|--------|-------------|-------------|
| Recovery | Easy | Requires backup |
| Storage | Increases over time | Constant |
| Query Complexity | Higher | Lower |
| Compliance | Better for audits | Needs separate audit |
| Referential Integrity | Maintained | Cascades needed |

---

## 3. Indexing Strategies

### Q7: Explain different types of indexes and when to use each.

**Answer:**

```sql
-- 1. B-Tree Index (Default) - Range queries, equality, sorting
CREATE INDEX idx_created_at ON orders(created_at);
-- Use: WHERE created_at BETWEEN '2024-01-01' AND '2024-12-31'
-- Use: ORDER BY created_at

-- 2. Hash Index - Exact match only (Memory tables)
CREATE INDEX idx_email USING HASH ON users(email);
-- Use: WHERE email = 'test@example.com'
-- NOT for: WHERE email LIKE 'test%'

-- 3. Composite Index - Multiple columns
CREATE INDEX idx_customer_date ON orders(customer_id, order_date);
-- Use: WHERE customer_id = 1 AND order_date > '2024-01-01'
-- Use: WHERE customer_id = 1 (leftmost prefix)
-- NOT efficient: WHERE order_date > '2024-01-01' (skips first column)

-- 4. Covering Index - Include all queried columns
CREATE INDEX idx_covering ON orders(customer_id, order_date, total_amount);
-- Covers: SELECT order_date, total_amount WHERE customer_id = 1
-- No table lookup needed

-- 5. Partial/Filtered Index (PostgreSQL)
CREATE INDEX idx_active_orders ON orders(order_date)
WHERE status = 'PENDING';
-- Only indexes pending orders - smaller, faster

-- 6. Full-Text Index
CREATE FULLTEXT INDEX idx_search ON products(name, description);
-- Use: WHERE MATCH(name, description) AGAINST('laptop gaming')

-- 7. Spatial Index (for geographic data)
CREATE SPATIAL INDEX idx_location ON stores(coordinates);
```

**Index Selection Matrix:**

| Query Pattern | Recommended Index |
|--------------|-------------------|
| `WHERE col = value` | B-Tree or Hash |
| `WHERE col > value` | B-Tree |
| `WHERE col LIKE 'prefix%'` | B-Tree |
| `WHERE col1 = x AND col2 = y` | Composite (col1, col2) |
| `ORDER BY col` | B-Tree |
| `GROUP BY col` | B-Tree |
| Text search | Full-Text |

---

### Q8: What is index selectivity and how does it affect performance?

**Answer:**

**Selectivity Formula:**
```
Selectivity = Number of Distinct Values / Total Rows
```

```sql
-- High Selectivity (Good for indexing)
-- email column: 1,000,000 unique / 1,000,000 rows = 1.0
CREATE INDEX idx_email ON users(email);  -- Excellent choice

-- Low Selectivity (Poor for indexing)
-- gender column: 3 unique / 1,000,000 rows = 0.000003
CREATE INDEX idx_gender ON users(gender);  -- Waste of resources

-- Check selectivity
SELECT 
    COUNT(DISTINCT email) / COUNT(*) as email_selectivity,
    COUNT(DISTINCT gender) / COUNT(*) as gender_selectivity,
    COUNT(DISTINCT status) / COUNT(*) as status_selectivity
FROM users;
```

**Cardinality Analysis:**

```sql
-- MySQL: Check index cardinality
SHOW INDEX FROM users;

-- PostgreSQL: Check statistics
SELECT 
    attname,
    n_distinct,
    correlation
FROM pg_stats 
WHERE tablename = 'users';

-- Rule of Thumb:
-- Selectivity > 0.1 → Good candidate for index
-- Selectivity < 0.01 → Poor candidate (unless used in combination)
```

**Composite Index Ordering:**

```sql
-- Order columns by selectivity (highest first)
-- If user_id has high selectivity and status has low:

-- GOOD: High selectivity first
CREATE INDEX idx_user_status ON orders(user_id, status);

-- LESS OPTIMAL: Low selectivity first
CREATE INDEX idx_status_user ON orders(status, user_id);
```

---

### Q9: How do you identify and fix slow queries using EXPLAIN?

**Answer:**

```sql
-- Basic EXPLAIN
EXPLAIN SELECT * FROM orders WHERE customer_id = 100;

-- Detailed EXPLAIN (MySQL)
EXPLAIN FORMAT=JSON SELECT * FROM orders WHERE customer_id = 100;

-- EXPLAIN ANALYZE (PostgreSQL) - Actually executes
EXPLAIN ANALYZE SELECT * FROM orders WHERE customer_id = 100;
```

**Key Metrics to Watch:**

| Metric | Good | Bad |
|--------|------|-----|
| type | const, eq_ref, ref | ALL (full scan) |
| rows | Low numbers | High numbers |
| Extra | Using index | Using filesort, Using temporary |
| key | Shows index name | NULL |

**Common Problems & Solutions:**

```sql
-- Problem 1: Full Table Scan
EXPLAIN SELECT * FROM orders WHERE YEAR(created_at) = 2024;
-- "type: ALL" - Function on column prevents index use

-- Solution: Sargable query
EXPLAIN SELECT * FROM orders 
WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';
-- "type: range" - Uses index

-- Problem 2: Filesort
EXPLAIN SELECT * FROM orders WHERE customer_id = 1 ORDER BY created_at;
-- "Extra: Using filesort" if no composite index

-- Solution: Composite index
CREATE INDEX idx_customer_created ON orders(customer_id, created_at);

-- Problem 3: Large IN clause
EXPLAIN SELECT * FROM products WHERE category_id IN (1,2,3,...100);
-- Can cause issues with large lists

-- Solution: Use JOIN instead
SELECT p.* FROM products p
INNER JOIN categories c ON p.category_id = c.category_id
WHERE c.parent_category_id = 5;

-- Problem 4: SELECT *
EXPLAIN SELECT * FROM orders WHERE customer_id = 100;
-- May not use covering index

-- Solution: Select only needed columns
EXPLAIN SELECT order_id, total_amount FROM orders WHERE customer_id = 100;
```

---

## 4. Keys & Constraints

### Q10: Explain different types of constraints with implementation examples.

**Answer:**

```sql
-- 1. PRIMARY KEY - Uniquely identifies each row
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    -- OR
    user_id BIGINT,
    CONSTRAINT pk_users PRIMARY KEY (user_id)
);

-- Composite Primary Key
CREATE TABLE order_items (
    order_id BIGINT,
    product_id BIGINT,
    quantity INT,
    PRIMARY KEY (order_id, product_id)
);

-- 2. FOREIGN KEY - Referential integrity
CREATE TABLE orders (
    order_id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    
    CONSTRAINT fk_orders_customer 
        FOREIGN KEY (customer_id) 
        REFERENCES customers(customer_id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

-- 3. UNIQUE - Prevent duplicates
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY,
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    
    CONSTRAINT uq_phone UNIQUE (phone)
);

-- 4. CHECK - Validate data
CREATE TABLE products (
    product_id BIGINT PRIMARY KEY,
    price DECIMAL(10,2),
    quantity INT,
    
    CONSTRAINT chk_price CHECK (price > 0),
    CONSTRAINT chk_quantity CHECK (quantity >= 0)
);

-- 5. NOT NULL - Prevent nulls
CREATE TABLE orders (
    order_id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,
    order_date DATE NOT NULL DEFAULT CURRENT_DATE
);

-- 6. DEFAULT - Provide default values
CREATE TABLE audit_log (
    log_id BIGINT PRIMARY KEY,
    action VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_processed BOOLEAN DEFAULT FALSE
);
```

**Foreign Key Actions:**

| Action | ON DELETE | ON UPDATE |
|--------|-----------|-----------|
| CASCADE | Delete child rows | Update child FKs |
| RESTRICT | Prevent deletion | Prevent update |
| SET NULL | Set FK to NULL | Set FK to NULL |
| SET DEFAULT | Set FK to default | Set FK to default |
| NO ACTION | Same as RESTRICT | Same as RESTRICT |

---

### Q11: How do you design self-referencing relationships?

**Answer:**

```sql
-- Example 1: Employee-Manager Hierarchy
CREATE TABLE employees (
    employee_id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    manager_id BIGINT NULL,
    
    CONSTRAINT fk_manager 
        FOREIGN KEY (manager_id) 
        REFERENCES employees(employee_id)
        ON DELETE SET NULL
);

-- Query: Get employee with manager name
SELECT 
    e.name as employee_name,
    m.name as manager_name
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.employee_id;

-- Example 2: Category Hierarchy
CREATE TABLE categories (
    category_id BIGINT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    parent_category_id BIGINT NULL,
    level INT DEFAULT 0,
    path VARCHAR(500),  -- Materialized path: "/1/5/12/"
    
    CONSTRAINT fk_parent_category
        FOREIGN KEY (parent_category_id)
        REFERENCES categories(category_id)
        ON DELETE CASCADE
);

-- Recursive CTE to get full hierarchy
WITH RECURSIVE category_tree AS (
    -- Base case: root categories
    SELECT 
        category_id, 
        name, 
        parent_category_id,
        0 as level,
        CAST(category_id AS CHAR(500)) as path
    FROM categories
    WHERE parent_category_id IS NULL
    
    UNION ALL
    
    -- Recursive case
    SELECT 
        c.category_id,
        c.name,
        c.parent_category_id,
        ct.level + 1,
        CONCAT(ct.path, '/', c.category_id)
    FROM categories c
    INNER JOIN category_tree ct ON c.parent_category_id = ct.category_id
)
SELECT * FROM category_tree;
```

**Hierarchy Design Patterns:**

| Pattern | Pros | Cons |
|---------|------|------|
| Adjacency List | Simple, easy updates | Slow tree queries |
| Materialized Path | Fast subtree queries | Path updates complex |
| Nested Set | Very fast reads | Complex writes |
| Closure Table | Flexible queries | Extra table needed |

---

## 5. Relationships & Cardinality

### Q12: How do you implement Many-to-Many relationships correctly?

**Answer:**

```sql
-- Junction/Bridge Table Pattern
CREATE TABLE students (
    student_id BIGINT PRIMARY KEY,
    name VARCHAR(100)
);

CREATE TABLE courses (
    course_id BIGINT PRIMARY KEY,
    title VARCHAR(200)
);

-- Junction table with composite primary key
CREATE TABLE student_courses (
    student_id BIGINT,
    course_id BIGINT,
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    grade VARCHAR(2) NULL,
    
    PRIMARY KEY (student_id, course_id),
    FOREIGN KEY (student_id) REFERENCES students(student_id) ON DELETE CASCADE,
    FOREIGN KEY (course_id) REFERENCES courses(course_id) ON DELETE CASCADE
);

-- With surrogate key (when junction has its own identity)
CREATE TABLE enrollments (
    enrollment_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    student_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP NULL,
    grade VARCHAR(2),
    
    UNIQUE KEY uq_student_course (student_id, course_id),
    FOREIGN KEY (student_id) REFERENCES students(student_id),
    FOREIGN KEY (course_id) REFERENCES courses(course_id)
);
```

**Many-to-Many with Attributes:**

```sql
-- Product Tags with weight/relevance
CREATE TABLE product_tags (
    product_id BIGINT,
    tag_id BIGINT,
    relevance_score DECIMAL(3,2) DEFAULT 1.00,
    added_by BIGINT,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (product_id, tag_id),
    INDEX idx_tag_relevance (tag_id, relevance_score DESC)
);

-- User Roles with validity period
CREATE TABLE user_roles (
    user_id BIGINT,
    role_id BIGINT,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NULL,
    granted_by BIGINT,
    
    PRIMARY KEY (user_id, role_id),
    INDEX idx_expires (expires_at)
);
```

---

### Q13: How do you handle optional vs mandatory relationships?

**Answer:**

```sql
-- Mandatory Relationship (NOT NULL foreign key)
CREATE TABLE orders (
    order_id BIGINT PRIMARY KEY,
    customer_id BIGINT NOT NULL,  -- Every order MUST have a customer
    
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id)
);

-- Optional Relationship (NULL allowed)
CREATE TABLE employees (
    employee_id BIGINT PRIMARY KEY,
    manager_id BIGINT NULL,  -- Not everyone has a manager (CEO)
    department_id BIGINT NULL,  -- Can be temporarily unassigned
    
    FOREIGN KEY (manager_id) REFERENCES employees(employee_id),
    FOREIGN KEY (department_id) REFERENCES departments(department_id)
);

-- One-to-One Optional (Profile may not exist)
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY,
    email VARCHAR(255) NOT NULL
);

CREATE TABLE user_profiles (
    profile_id BIGINT PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL,  -- 1:1 enforced by UNIQUE
    bio TEXT,
    avatar_url VARCHAR(500),
    
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- One-to-One Mandatory (Split table for performance)
CREATE TABLE products (
    product_id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL
);

CREATE TABLE product_details (
    product_id BIGINT PRIMARY KEY,  -- Same PK enforces 1:1
    full_description TEXT,
    specifications JSON,
    
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE CASCADE
);
```

---

## 6. Performance Optimization

### Q14: How do you optimize a database for high-write workloads?

**Answer:**

```sql
-- 1. Minimize Indexes on Write-Heavy Tables
-- Only keep essential indexes
CREATE TABLE events (
    event_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT,
    event_type VARCHAR(50),
    payload JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    
    -- Minimal indexes - only what's needed for reads
    INDEX idx_user_created (user_id, created_at)
);

-- 2. Use Batch Inserts
INSERT INTO events (user_id, event_type, payload)
VALUES 
    (1, 'click', '{}'),
    (2, 'view', '{}'),
    (3, 'purchase', '{}');
-- Much faster than individual inserts

-- 3. Partition by Time (for time-series data)
CREATE TABLE events (
    event_id BIGINT,
    user_id BIGINT,
    created_at TIMESTAMP,
    PRIMARY KEY (event_id, created_at)
)
PARTITION BY RANGE (UNIX_TIMESTAMP(created_at)) (
    PARTITION p_2024_01 VALUES LESS THAN (UNIX_TIMESTAMP('2024-02-01')),
    PARTITION p_2024_02 VALUES LESS THAN (UNIX_TIMESTAMP('2024-03-01')),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- 4. Use Appropriate Data Types
-- BIGINT (8 bytes) vs INT (4 bytes) matters at scale
-- VARCHAR(255) vs TEXT affects buffer efficiency
-- TIMESTAMP (4 bytes) vs DATETIME (8 bytes)

-- 5. Async Processing with Queue Table
CREATE TABLE write_queue (
    queue_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    target_table VARCHAR(50),
    operation ENUM('INSERT', 'UPDATE', 'DELETE'),
    payload JSON,
    status ENUM('PENDING', 'PROCESSING', 'COMPLETED', 'FAILED') DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP NULL,
    
    INDEX idx_status_created (status, created_at)
);
```

**Configuration Tuning (MySQL/InnoDB):**

```sql
-- Buffer pool size (70-80% of RAM for dedicated DB server)
innodb_buffer_pool_size = 12G

-- Log file size (larger = fewer checkpoints, better write performance)
innodb_log_file_size = 2G

-- Flush method (O_DIRECT bypasses OS cache)
innodb_flush_method = O_DIRECT

-- Disable query cache (deprecated, causes contention)
query_cache_type = 0
```

---

### Q15: Explain query optimization techniques with examples.

**Answer:**

```sql
-- 1. Use Covering Indexes
-- Bad: Requires table lookup
SELECT name, email FROM users WHERE status = 'active';

-- Add covering index
CREATE INDEX idx_status_covering ON users(status, name, email);
-- Now query uses "Using index" - no table access

-- 2. Avoid Functions on Indexed Columns
-- Bad: Can't use index
SELECT * FROM orders WHERE YEAR(created_at) = 2024;

-- Good: Sargable
SELECT * FROM orders 
WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';

-- 3. Use EXISTS Instead of IN for Subqueries
-- Slower (processes all results)
SELECT * FROM orders 
WHERE customer_id IN (SELECT customer_id FROM vip_customers);

-- Faster (stops at first match)
SELECT * FROM orders o
WHERE EXISTS (
    SELECT 1 FROM vip_customers v 
    WHERE v.customer_id = o.customer_id
);

-- 4. Optimize JOINs
-- Ensure JOIN columns are indexed
CREATE INDEX idx_customer_id ON orders(customer_id);

-- Use proper JOIN type
SELECT o.*, c.name
FROM orders o
INNER JOIN customers c ON o.customer_id = c.customer_id;  -- Use INNER when possible

-- 5. Limit Results Early
-- Bad: Fetches all, then limits
SELECT * FROM large_table ORDER BY created_at LIMIT 10;

-- Better: Use indexed column
SELECT * FROM large_table 
WHERE created_at > '2024-01-01' 
ORDER BY created_at 
LIMIT 10;

-- 6. Use UNION ALL Instead of UNION
-- UNION removes duplicates (sorts entire result)
SELECT name FROM customers UNION SELECT name FROM suppliers;

-- UNION ALL keeps all rows (no sort needed)
SELECT name FROM customers UNION ALL SELECT name FROM suppliers;

-- 7. Optimize COUNT Queries
-- Slow: Full table scan
SELECT COUNT(*) FROM orders WHERE status = 'pending';

-- Alternative: Maintain counter table
CREATE TABLE order_counts (
    status VARCHAR(20) PRIMARY KEY,
    count INT DEFAULT 0
);
-- Update via triggers or application logic
```

---

### Q16: How do you handle database connection pooling?

**Answer:**

**Why Connection Pooling:**
- Creating connections is expensive (TCP handshake, authentication)
- Database has limited connections
- Reusing connections improves performance

**Pool Sizing Formula:**
```
connections = ((core_count * 2) + effective_spindle_count)

For SSD: connections ≈ (CPU cores * 2) + 1
Example: 8 core server → ~17 connections per instance
```

**Configuration Example (HikariCP - Java):**

```java
HikariConfig config = new HikariConfig();
config.setMaximumPoolSize(20);          // Max connections
config.setMinimumIdle(5);               // Min idle connections
config.setIdleTimeout(300000);          // 5 minutes
config.setConnectionTimeout(30000);     // 30 seconds
config.setMaxLifetime(1800000);         // 30 minutes
```

**Database-Side Configuration (PostgreSQL):**

```sql
-- Check current connections
SELECT count(*) FROM pg_stat_activity;

-- Set max connections
ALTER SYSTEM SET max_connections = 200;

-- Use PgBouncer for connection pooling
-- pgbouncer.ini
[databases]
mydb = host=localhost dbname=mydb

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
```

---

## 7. Data Modeling Scenarios

### Q17: Design a schema for an audit logging system.

**Answer:**

```sql
-- Generic Audit Table
CREATE TABLE audit_log (
    audit_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    
    -- What changed
    table_name VARCHAR(100) NOT NULL,
    record_id BIGINT NOT NULL,
    operation ENUM('INSERT', 'UPDATE', 'DELETE') NOT NULL,
    
    -- Change details
    old_values JSON NULL,
    new_values JSON NULL,
    changed_columns JSON NULL,  -- ["status", "price"]
    
    -- Who/When/Where
    changed_by BIGINT NULL,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45) NULL,
    user_agent VARCHAR(500) NULL,
    
    -- Indexes
    INDEX idx_table_record (table_name, record_id),
    INDEX idx_changed_at (changed_at),
    INDEX idx_changed_by (changed_by)
) ENGINE=InnoDB
PARTITION BY RANGE (UNIX_TIMESTAMP(changed_at)) (
    PARTITION p_2024_q1 VALUES LESS THAN (UNIX_TIMESTAMP('2024-04-01')),
    PARTITION p_2024_q2 VALUES LESS THAN (UNIX_TIMESTAMP('2024-07-01')),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- Trigger for automatic auditing
DELIMITER //
CREATE TRIGGER audit_users_update
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    INSERT INTO audit_log (table_name, record_id, operation, old_values, new_values)
    VALUES (
        'users',
        OLD.user_id,
        'UPDATE',
        JSON_OBJECT(
            'email', OLD.email,
            'name', OLD.name,
            'status', OLD.status
        ),
        JSON_OBJECT(
            'email', NEW.email,
            'name', NEW.name,
            'status', NEW.status
        )
    );
END//
DELIMITER ;

-- Query: Get all changes to a specific record
SELECT * FROM audit_log 
WHERE table_name = 'users' AND record_id = 123
ORDER BY changed_at DESC;

-- Query: Get all changes by a specific user
SELECT * FROM audit_log 
WHERE changed_by = 456
ORDER BY changed_at DESC
LIMIT 100;
```

---

### Q18: Design a schema for a multi-tenant SaaS application.

**Answer:**

**Approach 1: Shared Schema with tenant_id**

```sql
-- Add tenant_id to all tables
CREATE TABLE tenants (
    tenant_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL,
    subdomain VARCHAR(50) UNIQUE NOT NULL,
    plan ENUM('FREE', 'BASIC', 'PRO', 'ENTERPRISE') DEFAULT 'FREE',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE TABLE users (
    user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    
    -- Unique email per tenant
    UNIQUE KEY uq_tenant_email (tenant_id, email),
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
    INDEX idx_tenant_id (tenant_id)
);

CREATE TABLE projects (
    project_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_id BIGINT NOT NULL,
    name VARCHAR(200) NOT NULL,
    
    FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
    INDEX idx_tenant_id (tenant_id)
);

-- ALWAYS filter by tenant_id
SELECT * FROM users WHERE tenant_id = @current_tenant_id;

-- Row Level Security (PostgreSQL)
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON users
    USING (tenant_id = current_setting('app.current_tenant')::BIGINT);
```

**Approach 2: Schema per Tenant**

```sql
-- Create schema for each tenant
CREATE SCHEMA tenant_acme;
CREATE SCHEMA tenant_globex;

-- Create tables in each schema
CREATE TABLE tenant_acme.users (...);
CREATE TABLE tenant_globex.users (...);

-- Application sets search_path based on tenant
SET search_path TO tenant_acme;
SELECT * FROM users;  -- Queries tenant_acme.users
```

**Comparison:**

| Aspect | Shared Schema | Schema per Tenant | Database per Tenant |
|--------|---------------|-------------------|---------------------|
| Isolation | Low | Medium | High |
| Resource Usage | Efficient | Moderate | High |
| Maintenance | Easier | Moderate | Complex |
| Scaling | Harder | Easier | Easiest |
| Cross-tenant Queries | Easy | Possible | Difficult |

---

### Q19: Design a schema for a real-time notification system.

**Answer:**

```sql
-- Notification Types/Templates
CREATE TABLE notification_types (
    type_id INT PRIMARY KEY AUTO_INCREMENT,
    type_code VARCHAR(50) UNIQUE NOT NULL,
    title_template VARCHAR(200) NOT NULL,
    body_template TEXT,
    default_channels JSON,  -- ["email", "push", "in_app"]
    is_active BOOLEAN DEFAULT TRUE
);

-- User notification preferences
CREATE TABLE user_notification_settings (
    user_id BIGINT,
    type_id INT,
    channel VARCHAR(20),  -- 'email', 'push', 'in_app', 'sms'
    is_enabled BOOLEAN DEFAULT TRUE,
    
    PRIMARY KEY (user_id, type_id, channel),
    FOREIGN KEY (type_id) REFERENCES notification_types(type_id)
);

-- Notifications queue
CREATE TABLE notifications (
    notification_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    type_id INT NOT NULL,
    
    -- Content
    title VARCHAR(200) NOT NULL,
    body TEXT,
    data JSON,  -- Additional payload
    action_url VARCHAR(500),
    
    -- Status tracking per channel
    email_status ENUM('PENDING', 'SENT', 'FAILED', 'SKIPPED') DEFAULT 'PENDING',
    push_status ENUM('PENDING', 'SENT', 'FAILED', 'SKIPPED') DEFAULT 'PENDING',
    in_app_status ENUM('PENDING', 'SENT', 'FAILED', 'SKIPPED') DEFAULT 'PENDING',
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    read_at TIMESTAMP NULL,
    
    -- Indexes
    INDEX idx_user_created (user_id, created_at DESC),
    INDEX idx_user_read (user_id, read_at),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (type_id) REFERENCES notification_types(type_id)
);

-- Real-time delivery tracking
CREATE TABLE notification_deliveries (
    delivery_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    notification_id BIGINT NOT NULL,
    channel VARCHAR(20) NOT NULL,
    status ENUM('PENDING', 'SENT', 'DELIVERED', 'READ', 'FAILED') DEFAULT 'PENDING',
    sent_at TIMESTAMP NULL,
    delivered_at TIMESTAMP NULL,
    read_at TIMESTAMP NULL,
    error_message TEXT NULL,
    retry_count INT DEFAULT 0,
    
    INDEX idx_notification (notification_id),
    INDEX idx_status_retry (status, retry_count),
    FOREIGN KEY (notification_id) REFERENCES notifications(notification_id)
);

-- Query: Get unread notifications for user
SELECT * FROM notifications 
WHERE user_id = @user_id AND read_at IS NULL
ORDER BY created_at DESC
LIMIT 50;

-- Query: Mark as read
UPDATE notifications 
SET read_at = CURRENT_TIMESTAMP
WHERE notification_id = @notification_id AND user_id = @user_id;
```

---

## 8. Transactions & ACID

### Q20: Explain ACID properties with practical examples.

**Answer:**

```sql
-- ATOMICITY: All or nothing
START TRANSACTION;

-- Transfer $100 from Account A to Account B
UPDATE accounts SET balance = balance - 100 WHERE account_id = 'A';
UPDATE accounts SET balance = balance + 100 WHERE account_id = 'B';

-- If any statement fails, both are rolled back
COMMIT;  -- Or ROLLBACK on error

-- CONSISTENCY: Database rules are maintained
CREATE TABLE accounts (
    account_id VARCHAR(20) PRIMARY KEY,
    balance DECIMAL(15,2) NOT NULL,
    
    CONSTRAINT chk_positive_balance CHECK (balance >= 0)
);

-- This will fail if it violates constraint
UPDATE accounts SET balance = balance - 1000 WHERE account_id = 'A';
-- Error if balance would go negative

-- ISOLATION: Concurrent transactions don't interfere
-- Session 1
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
START TRANSACTION;
SELECT balance FROM accounts WHERE account_id = 'A';  -- Returns 1000
-- Session 2 tries to modify same row - blocked until Session 1 commits

-- DURABILITY: Committed data survives crashes
COMMIT;  -- Data is persisted to disk
-- Even if server crashes, data is safe
```

**Isolation Levels:**

```sql
-- READ UNCOMMITTED (fastest, least safe)
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
-- Can see uncommitted changes (dirty reads)

-- READ COMMITTED (default in PostgreSQL, Oracle)
SET TRANSACTION ISOLATION LEVEL READ COMMITTED;
-- Only sees committed data, but may see different data on re-read

-- REPEATABLE READ (default in MySQL)
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ;
-- Same query returns same results within transaction

-- SERIALIZABLE (slowest, safest)
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
-- Full isolation, transactions appear to run sequentially
```

| Issue | Read Uncommitted | Read Committed | Repeatable Read | Serializable |
|-------|------------------|----------------|-----------------|--------------|
| Dirty Read | ✓ | ✗ | ✗ | ✗ |
| Non-Repeatable Read | ✓ | ✓ | ✗ | ✗ |
| Phantom Read | ✓ | ✓ | ✓ (MySQL: ✗) | ✗ |

---

### Q21: How do you handle deadlocks in database transactions?

**Answer:**

```sql
-- Deadlock Example:
-- Transaction 1:
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;  -- Locks row 1
UPDATE accounts SET balance = balance + 100 WHERE id = 2;  -- Waits for row 2

-- Transaction 2 (concurrent):
BEGIN;
UPDATE accounts SET balance = balance - 50 WHERE id = 2;   -- Locks row 2
UPDATE accounts SET balance = balance + 50 WHERE id = 1;   -- Waits for row 1
-- DEADLOCK!

-- Prevention Strategies:

-- 1. Consistent Lock Ordering
-- Always lock resources in the same order
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;  -- Always lock lower ID first
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- 2. Use SELECT FOR UPDATE for explicit locking
BEGIN;
SELECT * FROM accounts WHERE id IN (1, 2) ORDER BY id FOR UPDATE;
-- Now update in any order - locks acquired in consistent order
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;

-- 3. Keep transactions short
-- Bad: Long transaction
BEGIN;
SELECT * FROM orders WHERE status = 'pending' FOR UPDATE;
-- ... process each order (may take minutes)
COMMIT;

-- Good: Process in batches
BEGIN;
SELECT * FROM orders WHERE status = 'pending' ORDER BY id LIMIT 100 FOR UPDATE;
-- ... process batch quickly
COMMIT;

-- 4. Use lock timeouts
SET innodb_lock_wait_timeout = 5;  -- Wait max 5 seconds

-- 5. Retry pattern in application code
-- Pseudocode:
-- for retry in range(3):
--     try:
--         execute_transaction()
--         break
--     except DeadlockError:
--         sleep(random() * 0.1)
--         continue
```

**Detecting Deadlocks:**

```sql
-- MySQL: Show recent deadlock
SHOW ENGINE INNODB STATUS;  -- Check "LATEST DETECTED DEADLOCK" section

-- PostgreSQL: Check for locks
SELECT * FROM pg_locks WHERE NOT granted;

-- Monitor deadlock frequency
SELECT * FROM information_schema.innodb_metrics 
WHERE name = 'lock_deadlocks';
```

---

## 9. Partitioning & Sharding

### Q22: Explain database partitioning strategies with examples.

**Answer:**

```sql
-- 1. RANGE Partitioning (by date ranges)
CREATE TABLE orders (
    order_id BIGINT,
    customer_id BIGINT,
    order_date DATE,
    total DECIMAL(10,2),
    PRIMARY KEY (order_id, order_date)
)
PARTITION BY RANGE (YEAR(order_date)) (
    PARTITION p_2022 VALUES LESS THAN (2023),
    PARTITION p_2023 VALUES LESS THAN (2024),
    PARTITION p_2024 VALUES LESS THAN (2025),
    PARTITION p_2025 VALUES LESS THAN (2026),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- 2. LIST Partitioning (by category)
CREATE TABLE customers (
    customer_id BIGINT,
    name VARCHAR(100),
    region VARCHAR(20),
    PRIMARY KEY (customer_id, region)
)
PARTITION BY LIST COLUMNS (region) (
    PARTITION p_north VALUES IN ('NY', 'MA', 'CT'),
    PARTITION p_south VALUES IN ('FL', 'GA', 'TX'),
    PARTITION p_west VALUES IN ('CA', 'WA', 'OR'),
    PARTITION p_other VALUES IN (DEFAULT)
);

-- 3. HASH Partitioning (distribute evenly)
CREATE TABLE sessions (
    session_id BIGINT PRIMARY KEY,
    user_id BIGINT,
    data JSON
)
PARTITION BY HASH(user_id)
PARTITIONS 8;  -- Creates 8 partitions

-- 4. KEY Partitioning (MySQL-specific hash)
CREATE TABLE logs (
    log_id BIGINT,
    created_at TIMESTAMP,
    message TEXT,
    PRIMARY KEY (log_id, created_at)
)
PARTITION BY KEY(log_id)
PARTITIONS 4;

-- 5. Composite Partitioning (Range + Hash)
CREATE TABLE events (
    event_id BIGINT,
    user_id BIGINT,
    event_date DATE,
    event_type VARCHAR(50),
    PRIMARY KEY (event_id, event_date, user_id)
)
PARTITION BY RANGE (YEAR(event_date))
SUBPARTITION BY HASH(user_id)
SUBPARTITIONS 4 (
    PARTITION p_2024 VALUES LESS THAN (2025),
    PARTITION p_2025 VALUES LESS THAN (2026)
);
```

**Partition Management:**

```sql
-- Add new partition
ALTER TABLE orders ADD PARTITION (
    PARTITION p_2026 VALUES LESS THAN (2027)
);

-- Drop old partition (fast data deletion)
ALTER TABLE orders DROP PARTITION p_2022;

-- Reorganize partitions
ALTER TABLE orders REORGANIZE PARTITION p_future INTO (
    PARTITION p_2026 VALUES LESS THAN (2027),
    PARTITION p_future VALUES LESS THAN MAXVALUE
);

-- Query specific partition
SELECT * FROM orders PARTITION (p_2024)
WHERE customer_id = 123;
```

---

### Q23: When and how would you implement database sharding?

**Answer:**

**When to Shard:**
- Single database can't handle write load
- Data too large for single server
- Geographic distribution needed
- Regulatory/compliance requirements (data residency)

**Sharding Strategies:**

```sql
-- 1. Key-Based (Hash) Sharding
-- shard_id = hash(user_id) % num_shards

-- Shard 0: user_ids where hash(id) % 4 = 0
CREATE TABLE users_shard_0 (
    user_id BIGINT PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(255)
);

-- Application routing logic:
-- def get_shard(user_id):
--     return hash(user_id) % 4

-- 2. Range-Based Sharding
-- Shard by ID ranges
-- Shard 0: user_id 1 - 1,000,000
-- Shard 1: user_id 1,000,001 - 2,000,000

-- 3. Geographic Sharding
-- Shard by region
-- Shard US: All US customers
-- Shard EU: All EU customers

-- 4. Tenant-Based Sharding (Multi-tenant)
-- Each large tenant gets own shard
-- Small tenants share shards
```

**Cross-Shard Query Pattern:**

```sql
-- Global lookup table (replicated to all shards)
CREATE TABLE user_shard_mapping (
    user_id BIGINT PRIMARY KEY,
    shard_id INT NOT NULL,
    
    INDEX idx_shard (shard_id)
);

-- Scatter-gather query (application level)
-- Pseudocode:
-- results = []
-- for shard in all_shards:
--     partial_result = shard.query("SELECT * FROM orders WHERE status = 'pending'")
--     results.extend(partial_result)
-- return merge_and_sort(results)
```

**Sharding Comparison:**

| Strategy | Pros | Cons |
|----------|------|------|
| Hash | Even distribution | Difficult range queries |
| Range | Easy range queries | Hot spots possible |
| Geographic | Low latency | Uneven data distribution |
| Directory | Flexible | Single point of failure |

---

## 10. Real-World Design Problems

### Q24: Design a database schema for an e-commerce order management system.

**Answer:**

```sql
-- Core Tables
CREATE TABLE customers (
    customer_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_email (email)
);

CREATE TABLE addresses (
    address_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    customer_id BIGINT NOT NULL,
    address_type ENUM('BILLING', 'SHIPPING') NOT NULL,
    street_address VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100),
    postal_code VARCHAR(20) NOT NULL,
    country_code CHAR(2) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
    INDEX idx_customer (customer_id)
);

CREATE TABLE products (
    product_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    base_price DECIMAL(10,2) NOT NULL,
    category_id BIGINT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_sku (sku),
    INDEX idx_category (category_id),
    FULLTEXT INDEX idx_search (name, description)
);

CREATE TABLE product_inventory (
    inventory_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    product_id BIGINT NOT NULL,
    warehouse_id BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    reserved_quantity INT NOT NULL DEFAULT 0,
    
    UNIQUE KEY uq_product_warehouse (product_id, warehouse_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id),
    
    CHECK (quantity >= 0),
    CHECK (reserved_quantity >= 0),
    CHECK (reserved_quantity <= quantity)
);

-- Order Management
CREATE TABLE orders (
    order_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_number VARCHAR(50) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    
    -- Status tracking
    status ENUM('PENDING', 'CONFIRMED', 'PROCESSING', 'SHIPPED', 'DELIVERED', 'CANCELLED', 'REFUNDED') DEFAULT 'PENDING',
    
    -- Addresses (denormalized snapshot)
    shipping_address JSON NOT NULL,
    billing_address JSON NOT NULL,
    
    -- Amounts
    subtotal DECIMAL(12,2) NOT NULL,
    tax_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    shipping_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    total_amount DECIMAL(12,2) NOT NULL,
    
    -- Payment
    payment_status ENUM('PENDING', 'AUTHORIZED', 'CAPTURED', 'FAILED', 'REFUNDED') DEFAULT 'PENDING',
    payment_method VARCHAR(50),
    
    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    shipped_at TIMESTAMP NULL,
    delivered_at TIMESTAMP NULL,
    
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
    INDEX idx_customer (customer_id),
    INDEX idx_status (status),
    INDEX idx_created (created_at),
    INDEX idx_order_number (order_number)
);

CREATE TABLE order_items (
    order_item_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    
    -- Snapshot of product at order time
    product_name VARCHAR(255) NOT NULL,
    product_sku VARCHAR(50) NOT NULL,
    
    quantity INT NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    FOREIGN KEY (product_id) REFERENCES products(product_id),
    INDEX idx_order (order_id),
    
    CHECK (quantity > 0)
);

CREATE TABLE order_status_history (
    history_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    order_id BIGINT NOT NULL,
    old_status VARCHAR(20),
    new_status VARCHAR(20) NOT NULL,
    changed_by BIGINT,
    notes TEXT,
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (order_id) REFERENCES orders(order_id),
    INDEX idx_order_changed (order_id, changed_at)
);

-- Useful Queries
-- Get order with items
SELECT 
    o.*,
    JSON_ARRAYAGG(
        JSON_OBJECT(
            'product_name', oi.product_name,
            'quantity', oi.quantity,
            'unit_price', oi.unit_price
        )
    ) as items
FROM orders o
JOIN order_items oi ON o.order_id = oi.order_id
WHERE o.order_id = @order_id
GROUP BY o.order_id;

-- Daily sales report
SELECT 
    DATE(created_at) as order_date,
    COUNT(*) as order_count,
    SUM(total_amount) as total_revenue,
    AVG(total_amount) as avg_order_value
FROM orders
WHERE status NOT IN ('CANCELLED', 'REFUNDED')
AND created_at >= DATE_SUB(CURRENT_DATE, INTERVAL 30 DAY)
GROUP BY DATE(created_at)
ORDER BY order_date DESC;
```

---

### Q25: Design a database schema for a social media feed system.

**Answer:**

```sql
-- Users
CREATE TABLE users (
    user_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    bio VARCHAR(500),
    avatar_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_username (username)
);

-- Follow relationships
CREATE TABLE follows (
    follower_id BIGINT NOT NULL,
    following_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users(user_id),
    FOREIGN KEY (following_id) REFERENCES users(user_id),
    
    INDEX idx_following (following_id),
    
    CHECK (follower_id != following_id)
);

-- Posts
CREATE TABLE posts (
    post_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    content TEXT,
    media_urls JSON,  -- Array of image/video URLs
    
    -- Engagement counters (denormalized for performance)
    likes_count INT DEFAULT 0,
    comments_count INT DEFAULT 0,
    shares_count INT DEFAULT 0,
    
    -- Visibility
    visibility ENUM('PUBLIC', 'FOLLOWERS', 'PRIVATE') DEFAULT 'PUBLIC',
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    INDEX idx_user_created (user_id, created_at DESC),
    INDEX idx_created (created_at DESC),
    FULLTEXT INDEX idx_content (content)
);

-- Likes
CREATE TABLE post_likes (
    user_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id, post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (post_id) REFERENCES posts(post_id),
    
    INDEX idx_post (post_id)
);

-- Comments
CREATE TABLE comments (
    comment_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    post_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    parent_comment_id BIGINT NULL,  -- For replies
    content TEXT NOT NULL,
    likes_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_deleted BOOLEAN DEFAULT FALSE,
    
    FOREIGN KEY (post_id) REFERENCES posts(post_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (parent_comment_id) REFERENCES comments(comment_id),
    
    INDEX idx_post_created (post_id, created_at)
);

-- Feed cache (pre-computed feeds for performance)
CREATE TABLE user_feed_cache (
    user_id BIGINT NOT NULL,
    post_id BIGINT NOT NULL,
    post_score DECIMAL(10,4),  -- For ranking
    created_at TIMESTAMP,
    
    PRIMARY KEY (user_id, post_id),
    INDEX idx_user_score (user_id, post_score DESC)
);

-- Feed Generation Query (Pull model - simpler but slower)
SELECT p.*
FROM posts p
INNER JOIN follows f ON p.user_id = f.following_id
WHERE f.follower_id = @current_user_id
  AND p.is_deleted = FALSE
  AND p.visibility IN ('PUBLIC', 'FOLLOWERS')
ORDER BY p.created_at DESC
LIMIT 20;

-- With cursor-based pagination
SELECT p.*
FROM posts p
INNER JOIN follows f ON p.user_id = f.following_id
WHERE f.follower_id = @current_user_id
  AND p.is_deleted = FALSE
  AND p.created_at < @cursor_timestamp
ORDER BY p.created_at DESC
LIMIT 20;

-- Optimized: Use feed cache (Push model - faster reads)
SELECT p.*
FROM posts p
INNER JOIN user_feed_cache fc ON p.post_id = fc.post_id
WHERE fc.user_id = @current_user_id
ORDER BY fc.post_score DESC
LIMIT 20;
```

---

### Q26: Design a schema for a booking/reservation system (hotels, appointments, etc.).

**Answer:**

```sql
-- Resources (rooms, tables, time slots)
CREATE TABLE resources (
    resource_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_type VARCHAR(50) NOT NULL,  -- 'ROOM', 'TABLE', 'APPOINTMENT_SLOT'
    name VARCHAR(200) NOT NULL,
    description TEXT,
    capacity INT,
    price_per_unit DECIMAL(10,2),
    attributes JSON,  -- Flexible attributes
    is_active BOOLEAN DEFAULT TRUE,
    
    INDEX idx_type (resource_type)
);

-- Availability slots
CREATE TABLE availability (
    availability_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    resource_id BIGINT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status ENUM('AVAILABLE', 'BLOCKED', 'MAINTENANCE') DEFAULT 'AVAILABLE',
    price_override DECIMAL(10,2) NULL,  -- For dynamic pricing
    
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    INDEX idx_resource_time (resource_id, start_time, end_time),
    
    CHECK (end_time > start_time)
);

-- Bookings
CREATE TABLE bookings (
    booking_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    booking_reference VARCHAR(20) UNIQUE NOT NULL,
    customer_id BIGINT NOT NULL,
    resource_id BIGINT NOT NULL,
    
    -- Time slot
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    
    -- Status
    status ENUM('PENDING', 'CONFIRMED', 'CHECKED_IN', 'COMPLETED', 'CANCELLED', 'NO_SHOW') DEFAULT 'PENDING',
    
    -- Pricing
    base_price DECIMAL(10,2) NOT NULL,
    taxes DECIMAL(10,2) DEFAULT 0,
    fees DECIMAL(10,2) DEFAULT 0,
    discount DECIMAL(10,2) DEFAULT 0,
    total_price DECIMAL(10,2) NOT NULL,
    
    -- Payment
    payment_status ENUM('PENDING', 'PARTIAL', 'PAID', 'REFUNDED') DEFAULT 'PENDING',
    
    -- Metadata
    guests INT DEFAULT 1,
    special_requests TEXT,
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id),
    INDEX idx_reference (booking_reference),
    INDEX idx_resource_time (resource_id, start_time, end_time),
    INDEX idx_customer (customer_id),
    INDEX idx_status (status),
    
    CHECK (end_time > start_time)
);

-- Prevent double booking with check constraint or trigger
-- MySQL doesn't support exclusion constraints, use application logic or trigger

DELIMITER //
CREATE TRIGGER prevent_double_booking
BEFORE INSERT ON bookings
FOR EACH ROW
BEGIN
    DECLARE conflict_count INT;
    
    SELECT COUNT(*) INTO conflict_count
    FROM bookings
    WHERE resource_id = NEW.resource_id
      AND status NOT IN ('CANCELLED', 'NO_SHOW')
      AND (
          (NEW.start_time >= start_time AND NEW.start_time < end_time)
          OR (NEW.end_time > start_time AND NEW.end_time <= end_time)
          OR (NEW.start_time <= start_time AND NEW.end_time >= end_time)
      );
    
    IF conflict_count > 0 THEN
        SIGNAL SQLSTATE '45000' 
        SET MESSAGE_TEXT = 'Time slot already booked';
    END IF;
END//
DELIMITER ;

-- Check availability query
SELECT r.*, a.*
FROM resources r
LEFT JOIN availability a ON r.resource_id = a.resource_id
    AND a.start_time <= @requested_start
    AND a.end_time >= @requested_end
    AND a.status = 'AVAILABLE'
WHERE r.resource_type = 'ROOM'
  AND r.is_active = TRUE
  AND NOT EXISTS (
      SELECT 1 FROM bookings b
      WHERE b.resource_id = r.resource_id
        AND b.status NOT IN ('CANCELLED', 'NO_SHOW')
        AND b.start_time < @requested_end
        AND b.end_time > @requested_start
  );

-- Booking calendar view
SELECT 
    DATE(start_time) as date,
    resource_id,
    COUNT(*) as bookings_count,
    SUM(CASE WHEN status = 'CONFIRMED' THEN 1 ELSE 0 END) as confirmed,
    SUM(CASE WHEN status = 'PENDING' THEN 1 ELSE 0 END) as pending
FROM bookings
WHERE start_time BETWEEN @start_date AND @end_date
GROUP BY DATE(start_time), resource_id
ORDER BY date, resource_id;
```

---

## Bonus: Quick Reference Cheat Sheet

### Naming Conventions

```sql
-- Tables: plural, snake_case
users, order_items, product_categories

-- Columns: singular, snake_case
user_id, created_at, is_active

-- Primary keys: table_singular_id
user_id, order_id, product_id

-- Foreign keys: referenced_table_id
customer_id (in orders table)

-- Indexes: idx_table_columns
idx_users_email, idx_orders_customer_created

-- Constraints: type_table_columns
pk_users, uq_users_email, fk_orders_customer
```

### Common Data Types Guide

| Use Case | MySQL | PostgreSQL |
|----------|-------|------------|
| Auto-increment ID | BIGINT AUTO_INCREMENT | BIGSERIAL or BIGINT |
| UUID | CHAR(36) or BINARY(16) | UUID |
| Money | DECIMAL(10,2) | NUMERIC(10,2) or MONEY |
| Boolean | TINYINT(1) or BOOLEAN | BOOLEAN |
| Timestamp | TIMESTAMP | TIMESTAMPTZ |
| JSON | JSON | JSONB (better) |
| Variable string | VARCHAR(n) | VARCHAR(n) or TEXT |
| Fixed string | CHAR(n) | CHAR(n) |
| Large text | TEXT | TEXT |
| Enum | ENUM('A','B','C') | Custom TYPE |

### Index Checklist

```
□ Primary key columns (automatic)
□ Foreign key columns
□ Columns in WHERE clauses
□ Columns in JOIN conditions
□ Columns in ORDER BY
□ Columns in GROUP BY
□ High selectivity columns
□ Composite indexes for common query patterns
□ Covering indexes for read-heavy queries
```

---

## 11. SQL Database Selection

### Q27: Compare the major SQL databases and when to use each.

**Answer:**

#### Overview of Major SQL Databases

| Database | Type | Best For | License |
|----------|------|----------|--------|
| **PostgreSQL** | Open Source | Complex queries, GIS, JSON, extensibility | PostgreSQL License (Free) |
| **MySQL** | Open Source | Web applications, read-heavy workloads | GPL / Commercial |
| **MariaDB** | Open Source | MySQL replacement, enterprise features | GPL |
| **SQL Server** | Commercial | Enterprise, .NET integration, BI | Commercial (Express free) |
| **Oracle** | Commercial | Large enterprise, mission-critical | Commercial |
| **SQLite** | Embedded | Mobile apps, embedded systems, testing | Public Domain |

---

### Q28: PostgreSQL vs MySQL - When to choose which?

**Answer:**

#### Feature Comparison

| Feature | PostgreSQL | MySQL (InnoDB) |
|---------|------------|----------------|
| **ACID Compliance** | Full | Full |
| **JSON Support** | JSONB (binary, indexed) | JSON (text-based) |
| **Full-Text Search** | Built-in, powerful | Basic |
| **Geospatial (GIS)** | PostGIS (excellent) | Limited |
| **Replication** | Streaming, logical | Master-slave, Group |
| **Partitioning** | Declarative, flexible | Range, List, Hash |
| **Window Functions** | Full support | Full support (8.0+) |
| **CTEs (WITH clause)** | Full, recursive | Full (8.0+) |
| **Materialized Views** | Native support | Not supported |
| **Custom Types** | Excellent | Limited |
| **Extensions** | Rich ecosystem | Limited |
| **MVCC** | True MVCC | MVCC with gaps |

#### Performance Characteristics

```sql
-- PostgreSQL Strengths:
-- 1. Complex queries with many JOINs
-- 2. Write-heavy workloads (better MVCC)
-- 3. JSONB operations
-- 4. Large analytical queries

-- MySQL Strengths:
-- 1. Simple read-heavy queries
-- 2. High-concurrency web applications
-- 3. Replication setup simplicity
-- 4. Lower memory footprint
```

#### Choose PostgreSQL When:

```
✓ Complex data relationships and queries
✓ Need advanced data types (arrays, hstore, JSONB)
✓ Geospatial/GIS requirements
✓ Data integrity is critical (stricter by default)
✓ Need materialized views or advanced analytics
✓ Custom functions/extensions needed
✓ Write-heavy workloads
✓ Financial/scientific applications
```

#### Choose MySQL When:

```
✓ Simple CRUD web applications
✓ Read-heavy workloads (80%+ reads)
✓ Existing team expertise in MySQL
✓ Simple replication requirements
✓ WordPress, Drupal, or PHP applications
✓ Need minimal resource usage
✓ Simpler operational requirements
```

---

### Q29: When should you consider SQL Server or Oracle?

**Answer:**

#### Microsoft SQL Server

**Best Suited For:**

| Scenario | Why SQL Server |
|----------|----------------|
| .NET/C# Stack | Native integration, Entity Framework |
| Windows Environment | Active Directory, SSRS, SSIS |
| Business Intelligence | SSAS, Power BI integration |
| Enterprise with Support | Microsoft backing, SLAs |
| Hybrid Cloud | Azure SQL seamless migration |

**Key Features:**

```sql
-- SQL Server Unique Features:

-- 1. Temporal Tables (built-in history)
CREATE TABLE Products (
    ProductID INT PRIMARY KEY,
    Name NVARCHAR(100),
    Price DECIMAL(10,2),
    ValidFrom DATETIME2 GENERATED ALWAYS AS ROW START,
    ValidTo DATETIME2 GENERATED ALWAYS AS ROW END,
    PERIOD FOR SYSTEM_TIME (ValidFrom, ValidTo)
)
WITH (SYSTEM_VERSIONING = ON);

-- Query historical data automatically
SELECT * FROM Products
FOR SYSTEM_TIME AS OF '2024-01-01';

-- 2. Always Encrypted (column-level encryption)
CREATE COLUMN MASTER KEY MyCMK
WITH (KEY_STORE_PROVIDER_NAME = 'AZURE_KEY_VAULT');

-- 3. In-Memory OLTP
CREATE TABLE HotData (
    ID INT PRIMARY KEY NONCLUSTERED,
    Data NVARCHAR(100)
) WITH (MEMORY_OPTIMIZED = ON);

-- 4. PolyBase (query external data)
SELECT * FROM OPENROWSET(
    BULK 'https://storage.blob.core.windows.net/data/*.parquet',
    FORMAT = 'PARQUET'
) AS rows;
```

**Editions:**

| Edition | Use Case | Cost |
|---------|----------|------|
| Express | Development, small apps | Free (10GB limit) |
| Standard | SMB workloads | ~$4,000/core |
| Enterprise | Mission-critical | ~$15,000/core |
| Azure SQL | Cloud-native | Pay-per-use |

---

#### Oracle Database

**Best Suited For:**

| Scenario | Why Oracle |
|----------|------------|
| Large Enterprise | Proven at massive scale |
| Mission-Critical | RAC, Data Guard, zero downtime |
| Complex Transactions | Advanced locking, read consistency |
| Existing Investment | Oracle ERP, applications |
| Regulatory Compliance | Audit, security features |

**Key Features:**

```sql
-- Oracle Unique Features:

-- 1. Real Application Clusters (RAC)
-- Multiple instances, single database
-- Automatic failover, load balancing

-- 2. Flashback Technology
-- Query past data
SELECT * FROM orders AS OF TIMESTAMP 
    (SYSTIMESTAMP - INTERVAL '1' HOUR);

-- Restore dropped table
FLASHBACK TABLE orders TO BEFORE DROP;

-- 3. Partitioning (most advanced)
CREATE TABLE sales (
    sale_id NUMBER,
    sale_date DATE,
    amount NUMBER
)
PARTITION BY RANGE (sale_date)
INTERVAL (NUMTOYMINTERVAL(1, 'MONTH'))
(
    PARTITION p_initial VALUES LESS THAN (DATE '2024-01-01')
);
-- Auto-creates partitions!

-- 4. Advanced Compression
ALTER TABLE orders COMPRESS FOR OLTP;

-- 5. Automatic Indexing (19c+)
EXEC DBMS_AUTO_INDEX.CONFIGURE('AUTO_INDEX_MODE', 'IMPLEMENT');
```

**When to Avoid Oracle:**
- Startup/small company (cost prohibitive)
- Simple web applications (overkill)
- Tight budget constraints
- Cloud-native architecture preference

---

### Q30: How do you select a database for a new project?

**Answer:**

#### Decision Framework

```
┌─────────────────────────────────────────────────────────────┐
│                 DATABASE SELECTION MATRIX                    │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. WORKLOAD TYPE                                           │
│     ├── OLTP (transactions) → MySQL, PostgreSQL, SQL Server │
│     ├── OLAP (analytics) → PostgreSQL, SQL Server, Oracle   │
│     └── Mixed → PostgreSQL, SQL Server                      │
│                                                             │
│  2. SCALE REQUIREMENTS                                      │
│     ├── Small (<100GB) → Any                                │
│     ├── Medium (100GB-1TB) → PostgreSQL, MySQL, SQL Server  │
│     └── Large (1TB+) → PostgreSQL, Oracle, SQL Server       │
│                                                             │
│  3. BUDGET                                                  │
│     ├── Zero → PostgreSQL, MySQL, MariaDB                   │
│     ├── Moderate → SQL Server Standard, Cloud managed       │
│     └── Enterprise → Oracle, SQL Server Enterprise          │
│                                                             │
│  4. TEAM EXPERTISE                                          │
│     └── Leverage existing skills when possible              │
│                                                             │
│  5. ECOSYSTEM                                               │
│     ├── .NET → SQL Server                                   │
│     ├── PHP/WordPress → MySQL                               │
│     ├── Python/Ruby/Node → PostgreSQL                       │
│     └── Java/Enterprise → Oracle, PostgreSQL                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Selection Criteria Checklist

```sql
-- TECHNICAL REQUIREMENTS
□ Data volume (current and projected)
□ Transaction rate (reads/writes per second)
□ Query complexity (simple CRUD vs complex analytics)
□ Data types needed (JSON, GIS, arrays, full-text)
□ Consistency requirements (ACID strictness)
□ Availability requirements (99.9%, 99.99%, 99.999%)
□ Backup and recovery needs
□ Replication requirements

-- OPERATIONAL REQUIREMENTS
□ Team expertise and learning curve
□ Monitoring and management tools
□ Community/vendor support
□ Documentation quality
□ Hiring availability (DBA talent pool)

-- BUSINESS REQUIREMENTS
□ Budget (licensing, hardware, cloud costs)
□ Compliance requirements (HIPAA, PCI, GDPR)
□ Vendor lock-in concerns
□ Long-term viability
□ Integration with existing systems
```

---

### Q31: Compare cloud-managed SQL database options.

**Answer:**

#### Major Cloud Database Services

| Service | Provider | Based On | Best For |
|---------|----------|----------|----------|
| **Amazon RDS** | AWS | MySQL, PostgreSQL, SQL Server, Oracle, MariaDB | Multi-engine flexibility |
| **Amazon Aurora** | AWS | MySQL/PostgreSQL compatible | High performance, scalability |
| **Azure SQL Database** | Azure | SQL Server | .NET apps, Microsoft ecosystem |
| **Azure Database for PostgreSQL** | Azure | PostgreSQL | PostgreSQL on Azure |
| **Cloud SQL** | GCP | MySQL, PostgreSQL, SQL Server | GCP workloads |
| **AlloyDB** | GCP | PostgreSQL compatible | High-performance PostgreSQL |
| **CockroachDB** | Multi-cloud | PostgreSQL compatible | Global distribution |
| **PlanetScale** | Multi-cloud | MySQL compatible | Serverless MySQL |

---

#### Amazon Aurora Deep Dive

```sql
-- Aurora Advantages:
-- 1. 5x throughput of MySQL, 3x of PostgreSQL
-- 2. Storage auto-scales up to 128TB
-- 3. 6-way replication across 3 AZs
-- 4. Automatic failover in <30 seconds
-- 5. Read replicas for read scaling

-- Aurora Serverless v2
-- Auto-scales compute based on demand
-- Pay per ACU (Aurora Capacity Unit)
-- Ideal for variable workloads

-- When to choose Aurora:
✓ Need MySQL/PostgreSQL compatibility
✓ High availability is critical
✓ Unpredictable workloads (Serverless)
✓ Global applications (Aurora Global Database)
✓ Already on AWS
```

#### Azure SQL Database Deep Dive

```sql
-- Deployment Options:

-- 1. Single Database (isolated)
CREATE DATABASE mydb
(
    EDITION = 'GeneralPurpose',
    SERVICE_OBJECTIVE = 'GP_Gen5_2',  -- 2 vCores
    MAXSIZE = 32 GB
);

-- 2. Elastic Pool (shared resources)
-- Multiple databases share compute
-- Cost-effective for SaaS with variable per-tenant load

-- 3. Managed Instance (full SQL Server)
-- Near 100% compatibility with on-premises
-- VNet integration, SQL Agent, CLR

-- 4. Hyperscale
-- Up to 100TB databases
-- Instant backups regardless of size
-- Fast scale-out read replicas

-- Unique Features:
-- - Automatic tuning
-- - Intelligent Query Processing
-- - Built-in threat detection
-- - Geo-replication with one click
```

#### Cost Comparison (Approximate Monthly)

| Workload | AWS RDS | Aurora | Azure SQL | Cloud SQL |
|----------|---------|--------|-----------|----------|
| Small (2 vCPU, 8GB) | ~$100 | ~$150 | ~$150 | ~$100 |
| Medium (8 vCPU, 32GB) | ~$400 | ~$600 | ~$600 | ~$400 |
| Large (32 vCPU, 128GB) | ~$1,600 | ~$2,400 | ~$2,400 | ~$1,600 |
| Storage (per GB/month) | ~$0.12 | ~$0.10 | ~$0.12 | ~$0.17 |

*Prices vary by region and configuration*

---

### Q32: Self-hosted vs Cloud-managed databases - How to decide?

**Answer:**

#### Comparison Matrix

| Aspect | Self-Hosted | Cloud-Managed |
|--------|-------------|---------------|
| **Setup Time** | Days/Weeks | Minutes |
| **Maintenance** | Your team | Provider |
| **Patching** | Manual | Automatic |
| **Backups** | Configure yourself | Built-in |
| **HA/DR** | Complex setup | Toggle/Click |
| **Scaling** | Hardware provisioning | API call |
| **Cost (small)** | Higher | Lower |
| **Cost (large)** | Lower | Higher |
| **Control** | Full | Limited |
| **Performance Tuning** | Full access | Some restrictions |
| **Compliance** | Full control | Certifications available |

#### Choose Self-Hosted When:

```
✓ Strict data residency requirements
✓ Need full control over configuration
✓ Very large scale (cost optimization)
✓ Specialized hardware requirements
✓ Already have DBA expertise
✓ Predictable, steady workloads
✓ Air-gapped/offline requirements
```

#### Choose Cloud-Managed When:

```
✓ Speed of deployment is critical
✓ Limited DBA resources
✓ Variable/unpredictable workloads
✓ Need global distribution
✓ High availability without complexity
✓ Startup/growing company
✓ Want to focus on application, not infrastructure
```

#### Hybrid Approach

```sql
-- Production: Cloud-managed for reliability
-- Development: Self-hosted containers for cost

-- Docker Compose for local development
-- docker-compose.yml
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_PASSWORD: devpassword
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"

-- Terraform for cloud deployment
resource "aws_rds_cluster" "production" {
  cluster_identifier = "prod-cluster"
  engine            = "aurora-postgresql"
  engine_version    = "15.4"
  master_username   = "admin"
  master_password   = var.db_password
  
  serverlessv2_scaling_configuration {
    min_capacity = 0.5
    max_capacity = 16
  }
}
```

---

### Q33: What factors affect SQL database performance across different engines?

**Answer:**

#### Storage Engine Comparison

```sql
-- MySQL Storage Engines:

-- InnoDB (Default, recommended)
-- + ACID compliant
-- + Row-level locking
-- + Foreign key support
-- + Crash recovery
CREATE TABLE users (
    id INT PRIMARY KEY
) ENGINE=InnoDB;

-- MyISAM (Legacy)
-- + Faster for read-only
-- - Table-level locking
-- - No transactions
-- - No foreign keys
-- Avoid for new applications

-- PostgreSQL: Single engine with table access methods
-- All tables support full ACID, MVCC
```

#### Connection Handling

| Database | Connection Model | Pooling Recommendation |
|----------|------------------|----------------------|
| PostgreSQL | Process per connection | PgBouncer essential at scale |
| MySQL | Thread per connection | Built-in thread pool (Enterprise) |
| SQL Server | Thread pool | Built-in, generally sufficient |
| Oracle | Dedicated/Shared servers | Connection Manager for scale |

```sql
-- PostgreSQL: Check connections
SELECT count(*) FROM pg_stat_activity;

-- Max connections (default 100)
SHOW max_connections;

-- MySQL: Check connections
SHOW STATUS LIKE 'Threads_connected';
SHOW VARIABLES LIKE 'max_connections';

-- SQL Server
SELECT COUNT(*) FROM sys.dm_exec_connections;
```

#### Query Optimizer Differences

```sql
-- PostgreSQL: Cost-based optimizer
-- Very sophisticated, handles complex queries well
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT * FROM orders WHERE customer_id = 100;

-- Output includes:
-- - Actual time
-- - Buffer usage
-- - Row estimates vs actual

-- MySQL: Cost-based (improved in 8.0)
EXPLAIN ANALYZE
SELECT * FROM orders WHERE customer_id = 100;

-- SQL Server: Very advanced optimizer
SET STATISTICS IO ON;
SET STATISTICS TIME ON;
SELECT * FROM orders WHERE customer_id = 100;

-- View execution plan
SET SHOWPLAN_XML ON;
```

#### Memory Configuration

```sql
-- PostgreSQL Key Settings:
shared_buffers = 25% of RAM          -- Shared memory for caching
effective_cache_size = 75% of RAM    -- OS cache estimate
work_mem = 256MB                      -- Per-operation memory
maintenance_work_mem = 1GB            -- For VACUUM, CREATE INDEX

-- MySQL/InnoDB Key Settings:
innodb_buffer_pool_size = 70% of RAM  -- Primary cache
innodb_log_file_size = 1-2GB          -- Redo log size
innodb_flush_log_at_trx_commit = 1    -- Durability setting

-- SQL Server:
-- Usually auto-configured
-- max server memory = Total RAM - OS needs
EXEC sp_configure 'max server memory', 12288;  -- 12GB
RECONFIGURE;
```

---

### Q34: How do you migrate between different SQL databases?

**Answer:**

#### Migration Strategy

```
┌─────────────────────────────────────────────────────────────┐
│                    MIGRATION PHASES                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. ASSESSMENT                                              │
│     ├── Schema compatibility analysis                       │
│     ├── Data type mapping                                   │
│     ├── Feature gap analysis                                │
│     └── Performance baseline                                │
│                                                             │
│  2. SCHEMA CONVERSION                                       │
│     ├── Data types                                          │
│     ├── Stored procedures                                   │
│     ├── Functions and triggers                              │
│     └── Indexes and constraints                             │
│                                                             │
│  3. DATA MIGRATION                                          │
│     ├── Initial bulk load                                   │
│     ├── Incremental sync (CDC)                              │
│     └── Validation                                          │
│                                                             │
│  4. APPLICATION CHANGES                                     │
│     ├── Connection strings                                  │
│     ├── SQL syntax differences                              │
│     └── ORM/Driver updates                                  │
│                                                             │
│  5. TESTING & CUTOVER                                       │
│     ├── Performance testing                                 │
│     ├── Parallel run                                        │
│     └── Rollback plan                                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Common Data Type Mappings

| MySQL | PostgreSQL | SQL Server | Oracle |
|-------|------------|------------|--------|
| INT | INTEGER | INT | NUMBER(10) |
| BIGINT | BIGINT | BIGINT | NUMBER(19) |
| TINYINT | SMALLINT | TINYINT | NUMBER(3) |
| FLOAT | REAL | FLOAT | FLOAT |
| DOUBLE | DOUBLE PRECISION | FLOAT(53) | BINARY_DOUBLE |
| DECIMAL(p,s) | NUMERIC(p,s) | DECIMAL(p,s) | NUMBER(p,s) |
| VARCHAR(n) | VARCHAR(n) | NVARCHAR(n) | VARCHAR2(n) |
| TEXT | TEXT | NVARCHAR(MAX) | CLOB |
| BLOB | BYTEA | VARBINARY(MAX) | BLOB |
| DATETIME | TIMESTAMP | DATETIME2 | TIMESTAMP |
| DATE | DATE | DATE | DATE |
| BOOLEAN | BOOLEAN | BIT | NUMBER(1) |
| JSON | JSONB | NVARCHAR(MAX) | JSON (21c+) |
| AUTO_INCREMENT | SERIAL/IDENTITY | IDENTITY | SEQUENCE |

#### Syntax Differences

```sql
-- LIMIT/OFFSET
-- MySQL/PostgreSQL:
SELECT * FROM users LIMIT 10 OFFSET 20;

-- SQL Server:
SELECT * FROM users ORDER BY id OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY;

-- Oracle:
SELECT * FROM users OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY;  -- 12c+

-- String Concatenation
-- MySQL:
SELECT CONCAT(first_name, ' ', last_name) FROM users;

-- PostgreSQL:
SELECT first_name || ' ' || last_name FROM users;

-- SQL Server:
SELECT first_name + ' ' + last_name FROM users;
-- Or: CONCAT(first_name, ' ', last_name)

-- Oracle:
SELECT first_name || ' ' || last_name FROM users;

-- Current Timestamp
-- MySQL: NOW(), CURRENT_TIMESTAMP
-- PostgreSQL: NOW(), CURRENT_TIMESTAMP
-- SQL Server: GETDATE(), SYSDATETIME()
-- Oracle: SYSDATE, SYSTIMESTAMP

-- Auto-increment
-- MySQL:
CREATE TABLE users (id INT AUTO_INCREMENT PRIMARY KEY);

-- PostgreSQL:
CREATE TABLE users (id SERIAL PRIMARY KEY);
-- Or (modern): id INT GENERATED ALWAYS AS IDENTITY

-- SQL Server:
CREATE TABLE users (id INT IDENTITY(1,1) PRIMARY KEY);

-- Oracle:
CREATE TABLE users (id NUMBER GENERATED ALWAYS AS IDENTITY PRIMARY KEY);
```

#### Migration Tools

| Tool | From | To | Type |
|------|------|-----|------|
| **AWS DMS** | Any | Any (AWS) | Cloud service |
| **Azure DMA** | SQL Server, MySQL, Oracle | Azure SQL | Cloud service |
| **pgLoader** | MySQL, SQLite | PostgreSQL | Open source |
| **ora2pg** | Oracle | PostgreSQL | Open source |
| **SQLines** | Any | Any | Commercial |
| **Flyway/Liquibase** | N/A | Any | Schema versioning |

```bash
# pgLoader example: MySQL to PostgreSQL
pgloader mysql://user:pass@localhost/sourcedb \
         postgresql://user:pass@localhost/targetdb

# AWS DMS via CLI
aws dms create-replication-task \
  --replication-task-identifier mysql-to-postgres \
  --source-endpoint-arn arn:aws:dms:...:source \
  --target-endpoint-arn arn:aws:dms:...:target \
  --migration-type full-load-and-cdc
```

---

### Q35: What are the licensing considerations for SQL databases?

**Answer:**

#### License Comparison

| Database | License Type | Cost Model | Key Restrictions |
|----------|--------------|------------|------------------|
| **PostgreSQL** | PostgreSQL License | Free | None |
| **MySQL Community** | GPL v2 | Free | GPL copyleft |
| **MySQL Enterprise** | Commercial | Per socket/year | Proprietary features |
| **MariaDB Community** | GPL v2 | Free | GPL copyleft |
| **MariaDB Enterprise** | BSL/Commercial | Subscription | Enterprise features |
| **SQL Server Express** | Free | Free | 10GB DB, 1GB RAM, 4 cores |
| **SQL Server Standard** | Per core | ~$4,000/core | Some features limited |
| **SQL Server Enterprise** | Per core | ~$15,000/core | Full features |
| **Oracle Standard** | Per processor | ~$17,500/proc | Feature limited |
| **Oracle Enterprise** | Per processor | ~$47,500/proc | Add-ons extra |

#### Hidden Costs to Consider

```
┌─────────────────────────────────────────────────────────────┐
│                    TOTAL COST OF OWNERSHIP                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  LICENSING                                                  │
│  ├── Base license                                           │
│  ├── Support/maintenance (15-22% annually)                  │
│  ├── Additional features/options                            │
│  └── Development/test environments                          │
│                                                             │
│  INFRASTRUCTURE                                             │
│  ├── Servers/VMs                                            │
│  ├── Storage                                                │
│  ├── Network                                                │
│  └── Backup infrastructure                                  │
│                                                             │
│  OPERATIONS                                                 │
│  ├── DBA salary/training                                    │
│  ├── Monitoring tools                                       │
│  ├── Security/compliance                                    │
│  └── Disaster recovery                                      │
│                                                             │
│  MIGRATION (one-time)                                       │
│  ├── Assessment and planning                                │
│  ├── Schema/code conversion                                 │
│  ├── Testing                                                │
│  └── Downtime costs                                         │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Cost Optimization Strategies

```sql
-- 1. Use open-source where possible
-- PostgreSQL for 90% of use cases

-- 2. Cloud managed for operational savings
-- Compare: DBA salary vs RDS premium

-- 3. Right-size instances
-- Use auto-scaling (Aurora Serverless, Azure Hyperscale)

-- 4. Reserved capacity for predictable workloads
-- AWS: Reserved Instances (up to 60% savings)
-- Azure: Reserved Capacity (up to 65% savings)

-- 5. Separate OLTP and OLAP
-- Expensive database for transactions
-- Cheaper analytics solution (Redshift, BigQuery)
```

---

### Database Selection Quick Reference

```
┌─────────────────────────────────────────────────────────────┐
│                 QUICK SELECTION GUIDE                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  "I need a general-purpose database"                        │
│  → PostgreSQL                                               │
│                                                             │
│  "I'm building a WordPress/PHP site"                        │
│  → MySQL                                                    │
│                                                             │
│  "I'm in a Microsoft/.NET shop"                             │
│  → SQL Server                                               │
│                                                             │
│  "I need maximum reliability and have budget"               │
│  → Oracle or SQL Server Enterprise                          │
│                                                             │
│  "I need a database for mobile/embedded"                    │
│  → SQLite                                                   │
│                                                             │
│  "I need global distribution"                               │
│  → CockroachDB, Azure SQL Hyperscale                        │
│                                                             │
│  "I want serverless/auto-scaling"                           │
│  → Aurora Serverless, PlanetScale, Neon                     │
│                                                             │
│  "I need strong GIS support"                                │
│  → PostgreSQL + PostGIS                                     │
│                                                             │
│  "I need time-series data"                                  │
│  → TimescaleDB (PostgreSQL extension)                       │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 12. Situation-Based Interview Questions

### Q36: Your production database suddenly slows down during peak hours. How do you diagnose and fix it?

**Situation:** You receive alerts that API response times have increased from 200ms to 5+ seconds during peak traffic. Users are complaining about slow page loads.

**Answer:**

#### Step 1: Immediate Diagnosis

```sql
-- 1. Check active connections and running queries
-- PostgreSQL:
SELECT 
    pid,
    now() - pg_stat_activity.query_start AS duration,
    query,
    state,
    wait_event_type,
    wait_event
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY duration DESC
LIMIT 20;

-- MySQL:
SHOW FULL PROCESSLIST;
-- Or:
SELECT * FROM information_schema.processlist 
WHERE command != 'Sleep' 
ORDER BY time DESC;

-- 2. Check for lock contention
-- PostgreSQL:
SELECT 
    blocked_locks.pid AS blocked_pid,
    blocked_activity.usename AS blocked_user,
    blocking_locks.pid AS blocking_pid,
    blocking_activity.usename AS blocking_user,
    blocked_activity.query AS blocked_statement
FROM pg_catalog.pg_locks blocked_locks
JOIN pg_catalog.pg_stat_activity blocked_activity ON blocked_activity.pid = blocked_locks.pid
JOIN pg_catalog.pg_locks blocking_locks ON blocking_locks.locktype = blocked_locks.locktype
JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
WHERE NOT blocked_locks.granted;

-- MySQL:
SELECT * FROM information_schema.innodb_lock_waits;
```

#### Step 2: Identify the Root Cause

```sql
-- 3. Check slow query log
-- Enable if not already:
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 1;  -- Log queries > 1 second

-- 4. Check table statistics and missing indexes
-- PostgreSQL:
SELECT 
    schemaname,
    tablename,
    seq_scan,
    seq_tup_read,
    idx_scan,
    idx_tup_fetch,
    n_tup_ins,
    n_tup_upd,
    n_tup_del
FROM pg_stat_user_tables
WHERE seq_scan > 1000  -- Tables with many sequential scans
ORDER BY seq_scan DESC;

-- 5. Check cache hit ratio
-- PostgreSQL:
SELECT 
    sum(heap_blks_read) as heap_read,
    sum(heap_blks_hit) as heap_hit,
    sum(heap_blks_hit) / (sum(heap_blks_hit) + sum(heap_blks_read)) as ratio
FROM pg_statio_user_tables;
-- Should be > 0.99 (99%)
```

#### Step 3: Common Fixes

```sql
-- Issue: Missing index on frequently filtered column
CREATE INDEX CONCURRENTLY idx_orders_status_date 
ON orders(status, created_at);

-- Issue: Table bloat (PostgreSQL)
VACUUM ANALYZE orders;

-- Issue: Statistics out of date
ANALYZE orders;

-- Issue: Long-running transaction blocking others
-- Kill the blocking query (carefully!)
SELECT pg_terminate_backend(pid);

-- Issue: Connection pool exhausted
-- Increase pool size or add connection pooler (PgBouncer)
```

#### Step 4: Preventive Measures

```
✓ Set up query performance monitoring (pg_stat_statements)
✓ Configure slow query alerts
✓ Regular VACUUM and ANALYZE schedules
✓ Index usage analysis dashboard
✓ Connection pool monitoring
✓ Load testing before peak periods
```

---

### Q37: A critical table has grown to 500 million rows and queries are timing out. How do you handle this?

**Situation:** The `events` table started with 1 million rows but has grown to 500 million over 3 years. Simple queries now take 30+ seconds.

**Answer:**

#### Immediate Assessment

```sql
-- Check table size
-- PostgreSQL:
SELECT 
    pg_size_pretty(pg_total_relation_size('events')) as total_size,
    pg_size_pretty(pg_relation_size('events')) as table_size,
    pg_size_pretty(pg_indexes_size('events')) as index_size;

-- Check row count and estimate
SELECT reltuples::bigint AS row_estimate
FROM pg_class WHERE relname = 'events';

-- Analyze query patterns
SELECT query, calls, mean_time, total_time
FROM pg_stat_statements
WHERE query LIKE '%events%'
ORDER BY total_time DESC
LIMIT 10;
```

#### Short-Term Solutions (Hours)

```sql
-- 1. Add covering indexes for common queries
CREATE INDEX CONCURRENTLY idx_events_user_date 
ON events(user_id, event_date DESC)
INCLUDE (event_type, metadata);  -- Covering index

-- 2. Optimize the most problematic query
-- Before:
SELECT * FROM events WHERE user_id = 123 ORDER BY event_date DESC;

-- After:
SELECT event_id, event_type, event_date, metadata 
FROM events 
WHERE user_id = 123 
  AND event_date > CURRENT_DATE - INTERVAL '30 days'
ORDER BY event_date DESC
LIMIT 100;

-- 3. Add query hints/force index usage
-- MySQL:
SELECT * FROM events FORCE INDEX (idx_events_user_date) 
WHERE user_id = 123;
```

#### Medium-Term Solutions (Days)

```sql
-- 1. Implement table partitioning
-- Create new partitioned table
CREATE TABLE events_partitioned (
    event_id BIGINT,
    user_id BIGINT,
    event_type VARCHAR(50),
    event_date TIMESTAMP,
    metadata JSONB,
    PRIMARY KEY (event_id, event_date)
) PARTITION BY RANGE (event_date);

-- Create partitions
CREATE TABLE events_2024_q1 PARTITION OF events_partitioned
    FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');
CREATE TABLE events_2024_q2 PARTITION OF events_partitioned
    FOR VALUES FROM ('2024-04-01') TO ('2024-07-01');
-- ... more partitions

-- Migrate data in batches (during low traffic)
INSERT INTO events_partitioned
SELECT * FROM events 
WHERE event_date >= '2024-01-01' AND event_date < '2024-04-01';

-- 2. Archive old data
CREATE TABLE events_archive AS 
SELECT * FROM events WHERE event_date < '2023-01-01';

DELETE FROM events WHERE event_date < '2023-01-01';
```

#### Long-Term Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   RECOMMENDED ARCHITECTURE                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  HOT DATA (Recent 90 days)                                  │
│  ├── Partitioned table in primary database                  │
│  ├── Full indexing                                          │
│  └── Fast SSD storage                                       │
│                                                             │
│  WARM DATA (90 days - 1 year)                               │
│  ├── Separate partitions                                    │
│  ├── Minimal indexes                                        │
│  └── Standard storage                                       │
│                                                             │
│  COLD DATA (1+ years)                                       │
│  ├── Archive tables or data warehouse                       │
│  ├── Compressed storage                                     │
│  └── Accessed via analytics queries only                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

### Q38: Two microservices need to update the same data. How do you prevent race conditions?

**Situation:** Order Service and Inventory Service both need to update product stock. Sometimes overselling occurs because both services read the same stock value before either updates it.

**Answer:**

#### Understanding the Problem

```sql
-- Race condition scenario:
-- Time T1: Order Service reads stock = 10
-- Time T2: Inventory Service reads stock = 10
-- Time T3: Order Service sets stock = 10 - 3 = 7
-- Time T4: Inventory Service sets stock = 10 - 5 = 5
-- Result: Lost update! Should be 2, but got 5
```

#### Solution 1: Optimistic Locking (Recommended)

```sql
-- Add version column to table
ALTER TABLE products ADD COLUMN version INT DEFAULT 1;

-- Application reads current state
SELECT product_id, stock_quantity, version 
FROM products WHERE product_id = 100;
-- Returns: stock=10, version=5

-- Update with version check
UPDATE products 
SET stock_quantity = stock_quantity - 3,
    version = version + 1
WHERE product_id = 100 
  AND version = 5;  -- Only succeed if version unchanged

-- Check if update succeeded
-- If rows_affected = 0, another process modified it
-- Retry: Read again and retry the operation
```

**Application Code Pattern:**

```python
def decrement_stock(product_id, quantity, max_retries=3):
    for attempt in range(max_retries):
        # Read current state
        product = db.query(
            "SELECT stock_quantity, version FROM products WHERE product_id = %s",
            product_id
        )
        
        if product.stock_quantity < quantity:
            raise InsufficientStockError()
        
        # Attempt optimistic update
        result = db.execute(
            """UPDATE products 
               SET stock_quantity = stock_quantity - %s, version = version + 1
               WHERE product_id = %s AND version = %s""",
            quantity, product_id, product.version
        )
        
        if result.rows_affected == 1:
            return True  # Success!
        
        # Version changed, retry
        time.sleep(0.01 * (2 ** attempt))  # Exponential backoff
    
    raise ConcurrencyError("Could not update after retries")
```

#### Solution 2: Pessimistic Locking

```sql
-- Lock the row before reading
BEGIN;

SELECT stock_quantity 
FROM products 
WHERE product_id = 100 
FOR UPDATE;  -- Locks the row

-- Other transactions trying to select FOR UPDATE will wait

UPDATE products 
SET stock_quantity = stock_quantity - 3
WHERE product_id = 100;

COMMIT;
-- Lock released
```

#### Solution 3: Atomic Operations

```sql
-- Single atomic statement - no read-then-write
UPDATE products 
SET stock_quantity = stock_quantity - 3
WHERE product_id = 100 
  AND stock_quantity >= 3;  -- Prevent negative

-- Check result
-- rows_affected = 1: Success
-- rows_affected = 0: Insufficient stock
```

#### Solution 4: Database Serialization

```sql
-- Use SERIALIZABLE isolation level
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
BEGIN;

SELECT stock_quantity FROM products WHERE product_id = 100;
-- If another transaction modifies this row and commits,
-- this transaction will fail on commit

UPDATE products SET stock_quantity = stock_quantity - 3 WHERE product_id = 100;

COMMIT;  -- May fail with serialization error - must retry
```

#### Comparison

| Approach | Pros | Cons | Best For |
|----------|------|------|----------|
| Optimistic | No locks, high concurrency | Retries under contention | Low contention |
| Pessimistic | No retries needed | Blocks other transactions | High contention |
| Atomic | Simple, no race possible | Limited to simple operations | Simple decrements |
| Serializable | Full isolation | Performance penalty | Financial/critical |

---

### Q39: You need to implement multi-tenancy. A large customer wants data isolation. What do you recommend?

**Situation:** Your SaaS platform uses shared tables with `tenant_id`. A large enterprise customer requires stricter data isolation for compliance. How do you accommodate them without rebuilding the entire system?

**Answer:**

#### Hybrid Multi-Tenancy Approach

```sql
-- Current architecture: Shared schema
-- tenant_id in every table
SELECT * FROM orders WHERE tenant_id = @current_tenant;

-- New architecture: Hybrid
-- Small tenants: Continue using shared tables
-- Large/Enterprise tenants: Dedicated schema or database
```

#### Implementation Strategy

```sql
-- Step 1: Create tenant configuration table
CREATE TABLE tenant_config (
    tenant_id BIGINT PRIMARY KEY,
    tenant_name VARCHAR(200),
    isolation_level ENUM('SHARED', 'SCHEMA', 'DATABASE') DEFAULT 'SHARED',
    database_host VARCHAR(255) NULL,
    schema_name VARCHAR(100) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Regular tenant
INSERT INTO tenant_config (tenant_id, tenant_name, isolation_level)
VALUES (1, 'Small Corp', 'SHARED');

-- Enterprise tenant with schema isolation
INSERT INTO tenant_config (tenant_id, tenant_name, isolation_level, schema_name)
VALUES (100, 'Enterprise Inc', 'SCHEMA', 'tenant_enterprise');

-- VIP tenant with dedicated database
INSERT INTO tenant_config (tenant_id, tenant_name, isolation_level, database_host)
VALUES (500, 'Global Bank', 'DATABASE', 'db-globalbank.region.rds.amazonaws.com');
```

#### Schema-Level Isolation

```sql
-- Create dedicated schema for enterprise tenant
CREATE SCHEMA tenant_enterprise;

-- Create identical tables in tenant schema
CREATE TABLE tenant_enterprise.orders (
    order_id BIGINT PRIMARY KEY,
    -- Note: No tenant_id needed!
    customer_name VARCHAR(100),
    total_amount DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Migrate existing data
INSERT INTO tenant_enterprise.orders 
SELECT order_id, customer_name, total_amount, created_at
FROM public.orders 
WHERE tenant_id = 100;

-- Delete from shared table
DELETE FROM public.orders WHERE tenant_id = 100;

-- Application routing logic:
-- if tenant.isolation_level == 'SCHEMA':
--     SET search_path TO tenant_enterprise;
-- Query: SELECT * FROM orders; (no tenant_id filter needed)
```

#### Row-Level Security (RLS) for Shared Tenants

```sql
-- PostgreSQL: Row Level Security
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

-- Create policy
CREATE POLICY tenant_isolation ON orders
    USING (tenant_id = current_setting('app.current_tenant')::BIGINT);

-- Application sets tenant context
SET app.current_tenant = '123';

-- All queries automatically filtered
SELECT * FROM orders;  -- Only sees tenant 123's orders
```

#### Connection Router Pattern

```python
class TenantConnectionRouter:
    def get_connection(self, tenant_id):
        config = self.get_tenant_config(tenant_id)
        
        if config.isolation_level == 'DATABASE':
            # Return connection to dedicated database
            return create_connection(host=config.database_host)
        
        elif config.isolation_level == 'SCHEMA':
            # Return connection with schema set
            conn = self.shared_pool.get_connection()
            conn.execute(f"SET search_path TO {config.schema_name}")
            return conn
        
        else:  # SHARED
            # Return shared connection with RLS
            conn = self.shared_pool.get_connection()
            conn.execute(f"SET app.current_tenant = '{tenant_id}'")
            return conn
```

#### Data Migration for New Enterprise Customer

```sql
-- 1. Create schema
CREATE SCHEMA tenant_newcustomer;

-- 2. Copy table structures
CREATE TABLE tenant_newcustomer.orders (LIKE public.orders INCLUDING ALL);
CREATE TABLE tenant_newcustomer.customers (LIKE public.customers INCLUDING ALL);

-- 3. Migrate data with minimal downtime
BEGIN;
-- Lock briefly to prevent concurrent writes
LOCK TABLE public.orders IN EXCLUSIVE MODE;

-- Copy data
INSERT INTO tenant_newcustomer.orders 
SELECT * FROM public.orders WHERE tenant_id = @enterprise_tenant_id;

-- Update tenant config
UPDATE tenant_config 
SET isolation_level = 'SCHEMA', schema_name = 'tenant_newcustomer'
WHERE tenant_id = @enterprise_tenant_id;

COMMIT;

-- 4. Clean up shared tables (can be done async)
DELETE FROM public.orders WHERE tenant_id = @enterprise_tenant_id;
```

---

### Q40: A report that used to take 2 seconds now takes 2 minutes. What changed?

**Situation:** A daily sales report suddenly degraded from 2 seconds to 2+ minutes. No code changes were deployed. What's your investigation approach?

**Answer:**

#### Systematic Investigation

```sql
-- Step 1: Get the actual execution plan
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT)
SELECT 
    date_trunc('day', order_date) as day,
    COUNT(*) as orders,
    SUM(total_amount) as revenue
FROM orders
WHERE order_date >= '2024-01-01'
GROUP BY date_trunc('day', order_date)
ORDER BY day;

-- Compare with expected plan
-- Look for:
-- - Seq Scan instead of Index Scan
-- - High actual rows vs estimated rows
-- - Many buffer reads
```

#### Common Causes and Detection

```sql
-- Cause 1: Statistics are stale
SELECT 
    schemaname,
    tablename,
    last_analyze,
    last_autoanalyze,
    n_live_tup,
    n_dead_tup
FROM pg_stat_user_tables
WHERE tablename = 'orders';

-- Fix:
ANALYZE orders;

-- Cause 2: Index was dropped or corrupted
SELECT 
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'orders';

-- Check if index is being used:
SELECT 
    indexrelname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
WHERE relname = 'orders';

-- Cause 3: Table bloat (many dead tuples)
SELECT 
    n_live_tup,
    n_dead_tup,
    round(100.0 * n_dead_tup / NULLIF(n_live_tup + n_dead_tup, 0), 2) as dead_pct
FROM pg_stat_user_tables
WHERE tablename = 'orders';

-- Fix:
VACUUM FULL orders;  -- Locks table!
-- Or:
VACUUM orders;  -- Less aggressive

-- Cause 4: Data volume crossed a threshold
SELECT 
    pg_size_pretty(pg_total_relation_size('orders')),
    (SELECT count(*) FROM orders) as row_count;
```

#### Comparing Query Plans

```sql
-- Save current execution plan
EXPLAIN (FORMAT JSON)
SELECT ... -- Your query
\g /tmp/current_plan.json

-- Compare key metrics:
```

| Metric | Before (Fast) | After (Slow) | Issue |
|--------|---------------|--------------|-------|
| Scan Type | Index Scan | Seq Scan | Index not used |
| Rows Estimated | 10,000 | 10,000 | - |
| Rows Actual | 10,000 | 5,000,000 | Table grew |
| Buffers | 500 | 50,000 | More data to read |
| Sort Method | quicksort | external merge | Memory overflow |

#### Query Plan Changed - Why?

```sql
-- PostgreSQL planner makes cost-based decisions
-- Factors that change decisions:

-- 1. Table statistics changed
SELECT attname, n_distinct, correlation
FROM pg_stats WHERE tablename = 'orders';

-- 2. Table grew beyond planner threshold
-- Seq scan cheaper when reading > ~20% of table
SELECT 
    seq_scan,
    idx_scan,
    seq_tup_read,
    idx_tup_fetch
FROM pg_stat_user_tables
WHERE tablename = 'orders';

-- 3. Index correlation degraded
-- Physical vs logical order diverged
REINDEX INDEX idx_orders_date;

-- 4. work_mem too small for new data volume
SHOW work_mem;  -- Default 4MB
SET work_mem = '256MB';  -- Try larger
-- Then re-test query
```

#### Resolution Steps

```sql
-- 1. Update statistics
ANALYZE orders;

-- 2. If index not being used, hint or recreate
DROP INDEX idx_orders_date;
CREATE INDEX idx_orders_date ON orders(order_date);

-- 3. If table bloated
VACUUM ANALYZE orders;

-- 4. If data volume issue, consider partitioning
-- See Q37 for partitioning strategy

-- 5. Optimize the query itself
-- Add date range limit
WHERE order_date >= CURRENT_DATE - INTERVAL '90 days'

-- 6. Create covering index
CREATE INDEX idx_orders_date_covering 
ON orders(order_date) 
INCLUDE (total_amount);
```

---

### Q41: You discover the database is running without backups. How do you fix this and prevent data loss?

**Situation:** After joining a new team, you discover the production database has no backup strategy. How do you implement one urgently while minimizing risk?

**Answer:**

#### Immediate Actions (Day 1)

```sql
-- 1. Take immediate manual backup
-- PostgreSQL:
pg_dump -h localhost -U username -F c -f /backup/emergency_backup.dump dbname

-- MySQL:
mysqldump -u username -p --single-transaction --routines dbname > /backup/emergency_backup.sql

-- SQL Server:
BACKUP DATABASE dbname TO DISK = '/backup/emergency_backup.bak' WITH COMPRESSION;

-- 2. Copy to external location immediately
-- AWS S3:
aws s3 cp /backup/emergency_backup.dump s3://company-backups/emergency/
```

#### Backup Strategy Implementation

```
┌─────────────────────────────────────────────────────────────┐
│                    BACKUP STRATEGY                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  TIER 1: Continuous (Point-in-Time Recovery)                │
│  ├── WAL archiving (PostgreSQL)                             │
│  ├── Binary log (MySQL)                                     │
│  └── Transaction log backup (SQL Server)                    │
│  → RPO: Minutes                                             │
│                                                             │
│  TIER 2: Daily Full Backups                                 │
│  ├── Automated nightly backups                              │
│  ├── Stored in separate region                              │
│  └── Retained for 30 days                                   │
│  → RTO: Hours                                               │
│                                                             │
│  TIER 3: Weekly/Monthly Archives                            │
│  ├── Compressed, encrypted                                  │
│  ├── Cold storage (S3 Glacier)                              │
│  └── Retained for 1-7 years (compliance)                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### PostgreSQL Point-in-Time Recovery Setup

```sql
-- postgresql.conf changes:
wal_level = replica
archive_mode = on
archive_command = 'aws s3 cp %p s3://db-backups/wal/%f'
archive_timeout = 60  -- Archive at least every minute

-- Automated backup script
#!/bin/bash
set -e
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="/tmp/backup_${DATE}.dump"

# Full backup
pg_dump -h localhost -U postgres -F c -f $BACKUP_FILE dbname

# Upload to S3
aws s3 cp $BACKUP_FILE s3://db-backups/daily/$BACKUP_FILE

# Cleanup local
rm $BACKUP_FILE

# Notify
echo "Backup completed: $BACKUP_FILE" | slack-notify
```

#### Cloud-Managed Backup (Recommended)

```sql
-- AWS RDS: Enable automated backups
aws rds modify-db-instance \
  --db-instance-identifier mydb \
  --backup-retention-period 30 \
  --preferred-backup-window "03:00-04:00" \
  --apply-immediately

-- Azure SQL: Automatic backups included
-- Configure long-term retention:
az sql db ltr-policy set \
  --resource-group mygroup \
  --server myserver \
  --database mydb \
  --weekly-retention P4W \
  --monthly-retention P12M \
  --yearly-retention P5Y \
  --week-of-year 1
```

#### Backup Verification (Critical!)

```sql
-- Schedule regular restore tests
-- Monthly restore to test environment:

-- 1. Restore backup
pg_restore -h test-db -U postgres -d test_restore /backup/latest.dump

-- 2. Verify data integrity
SELECT COUNT(*) FROM orders;  -- Compare with production
SELECT MAX(created_at) FROM orders;  -- Check recent data

-- 3. Run sample queries
SELECT * FROM orders WHERE order_id = 12345;  -- Spot check

-- 4. Automated verification script
#!/bin/bash
PROD_COUNT=$(psql -h prod-db -t -c "SELECT COUNT(*) FROM orders")
TEST_COUNT=$(psql -h test-db -t -c "SELECT COUNT(*) FROM orders")

if [ "$PROD_COUNT" != "$TEST_COUNT" ]; then
    echo "ALERT: Backup verification failed!" | slack-notify
    exit 1
fi
```

#### Documentation Required

```markdown
## Backup Runbook

### Daily Backup Schedule
- Time: 03:00 UTC
- Retention: 30 days
- Location: s3://company-backups/daily/

### Recovery Procedures
1. Point-in-time recovery: [link to doc]
2. Full restore: [link to doc]
3. Partial table restore: [link to doc]

### Contacts
- DBA on-call: +1-xxx-xxx-xxxx
- AWS Support: Case #xxxxx

### Last Successful Restore Test
- Date: 2024-01-15
- Duration: 45 minutes
- Result: SUCCESS
```

---

### Q42: The application needs to scale globally. How do you design the database architecture?

**Situation:** Your application currently serves US users from a single database in us-east-1. You need to expand to Europe and Asia with <100ms latency requirements.

**Answer:**

#### Architecture Options

```
┌─────────────────────────────────────────────────────────────┐
│           GLOBAL DATABASE ARCHITECTURE OPTIONS              │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Option 1: Read Replicas per Region                         │
│  ┌─────────┐    Async     ┌─────────┐                       │
│  │ Primary │ ──────────── │ Replica │                       │
│  │ US-East │              │ EU-West │                       │
│  └─────────┘              └─────────┘                       │
│       │                        ▲                            │
│       │        Async           ╎                            │
│       └────────────────── ┌─────────┐                       │
│                           │ Replica │                       │
│                           │AP-Tokyo │                       │
│                           └─────────┘                       │
│  Pros: Simple, low cost                                     │
│  Cons: Writes go to US (high latency)                       │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Option 2: Active-Active Multi-Master                       │
│  ┌─────────┐   Sync      ┌─────────┐                        │
│  │ Primary │ ◄────────── │ Primary │                        │
│  │ US-East │ ──────────► │ EU-West │                        │
│  └─────────┘             └─────────┘                        │
│       ▲                        ▲                            │
│       │         Sync           │                            │
│       └──────── ┌─────────┐ ───┘                            │
│                 │ Primary │                                 │
│                 │AP-Tokyo │                                 │
│                 └─────────┘                                 │
│  Pros: Local reads AND writes                               │
│  Cons: Conflict resolution, complexity, cost                │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Option 3: Regional Databases with Eventual Sync            │
│  ┌─────────┐             ┌─────────┐                        │
│  │  US DB  │             │  EU DB  │                        │
│  │  Users  │    Event    │  Users  │                        │
│  └─────────┘ ◄── Bus ──► └─────────┘                        │
│       ▲                        ▲                            │
│       │                        │                            │
│       └────────┬───────────────┘                            │
│           ┌─────────┐                                       │
│           │ Asia DB │                                       │
│           │  Users  │                                       │
│           └─────────┘                                       │
│  Pros: Full local performance, isolation                    │
│  Cons: Complex sync, eventual consistency                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Recommended: Aurora Global Database (AWS)

```sql
-- Setup Aurora Global Database
-- Primary: us-east-1
-- Secondary: eu-west-1, ap-northeast-1

-- Characteristics:
-- - Replication lag: <1 second typically
-- - Promotes secondary in <1 minute
-- - Up to 5 secondary regions

-- Application connection routing:
-- Read: Connect to local region endpoint
-- Write: Connect to primary region endpoint

-- Write-forwarding (Aurora 3.0+):
-- Local writes forwarded to primary
-- Added latency but simpler app logic
```

#### Data Partitioning by Region

```sql
-- Option 3 implementation: Regional user data

-- User routing table (replicated globally)
CREATE TABLE user_regions (
    user_id BIGINT PRIMARY KEY,
    home_region VARCHAR(20) NOT NULL,  -- 'US', 'EU', 'APAC'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Regional tables (in each region's database)
CREATE TABLE orders (
    order_id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    -- No tenant_id needed - regional isolation
    total_amount DECIMAL(10,2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cross-region sync via events
CREATE TABLE sync_events (
    event_id BIGINT PRIMARY KEY,
    entity_type VARCHAR(50),  -- 'user_profile', 'shared_config'
    entity_id BIGINT,
    operation VARCHAR(20),  -- 'INSERT', 'UPDATE', 'DELETE'
    payload JSONB,
    source_region VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    replicated_at TIMESTAMP NULL,
    
    INDEX idx_pending (replicated_at) WHERE replicated_at IS NULL
);
```

#### Application Routing Logic

```python
class GlobalDatabaseRouter:
    def get_connection(self, user_id, operation='read'):
        # Determine user's home region
        user_region = cache.get(f'user_region:{user_id}')
        if not user_region:
            user_region = self.lookup_user_region(user_id)
        
        if operation == 'write':
            # Writes go to user's home region
            return self.get_regional_connection(user_region, 'primary')
        else:
            # Reads from closest region
            closest = self.get_closest_region()
            
            if user_region == closest:
                # User data is local
                return self.get_regional_connection(closest, 'replica')
            else:
                # May need to route to user's home region
                return self.get_regional_connection(user_region, 'replica')
    
    def get_closest_region(self):
        # Based on request origin
        client_region = request.headers.get('X-Client-Region', 'US')
        return self.region_mapping.get(client_region, 'US')
```

#### Handling Global vs Regional Data

```sql
-- Global data (read everywhere, write to primary)
-- - Product catalog
-- - Configuration
-- - Reference data
CREATE TABLE products (
    product_id BIGINT PRIMARY KEY,
    name VARCHAR(200),
    price DECIMAL(10,2)
);  -- Replicated to all regions

-- Regional data (read/write in home region)
-- - User profiles
-- - Orders
-- - User preferences
CREATE TABLE orders (
    order_id BIGINT PRIMARY KEY,
    user_id BIGINT,
    -- Lives only in user's home region
);

-- Shared data needing sync
-- - User activity across regions
-- - Analytics events
-- Use event-based eventual consistency
```

---

### Q43: A developer accidentally deleted production data. How do you recover and prevent this?

**Situation:** A developer ran `DELETE FROM users WHERE status = 'inactive'` without a WHERE clause test, deleting 50,000 active users at 2 PM on a business day.

**Answer:**

#### Immediate Recovery

```sql
-- Step 1: Assess the damage
SELECT COUNT(*) FROM users;  -- How many remain?

-- Step 2: Check if transaction is still open
-- If caught immediately:
ROLLBACK;  -- If in transaction

-- Step 3: Point-in-Time Recovery (if available)
-- PostgreSQL with WAL archiving:
pg_restore --target-time="2024-03-15 13:59:00" \
  --target-action=promote \
  -d recovered_db /backup/base_backup

-- Extract deleted data
CREATE TABLE users_recovered AS 
SELECT * FROM dblink(
    'dbname=recovered_db',
    'SELECT * FROM users'
) AS t(user_id BIGINT, email VARCHAR, ...);

-- Step 4: Restore deleted records
INSERT INTO users 
SELECT * FROM users_recovered ur
WHERE NOT EXISTS (
    SELECT 1 FROM users u WHERE u.user_id = ur.user_id
);
```

#### If Using Soft Deletes (Better Position)

```sql
-- If soft delete was implemented:
-- Just update the flag
UPDATE users 
SET is_deleted = FALSE, deleted_at = NULL
WHERE deleted_at >= '2024-03-15 14:00:00';
```

#### AWS RDS Recovery

```bash
# Restore to point in time
aws rds restore-db-instance-to-point-in-time \
  --source-db-instance-identifier production-db \
  --target-db-instance-identifier recovery-db \
  --restore-time "2024-03-15T13:59:00Z" \
  --db-subnet-group-name prod-subnets

# Wait for restore
aws rds wait db-instance-available --db-instance-identifier recovery-db

# Connect and extract data
pg_dump -h recovery-db.xxx.rds.amazonaws.com -t users | \
  psql -h production-db.xxx.rds.amazonaws.com
```

#### Prevention Measures

```sql
-- 1. Require WHERE clause for DELETE/UPDATE (PostgreSQL)
-- Use pg_safeupdate extension
CREATE EXTENSION pg_safeupdate;
SET safeupdate.enabled = true;

-- This will fail:
DELETE FROM users;  -- ERROR: DELETE requires WHERE clause

-- 2. Use transaction wrappers
-- Force explicit transaction confirmation
BEGIN;
DELETE FROM users WHERE status = 'inactive';
-- Shows: 50,000 rows affected
-- Developer sees this is wrong:
ROLLBACK;

-- 3. Implement soft deletes
ALTER TABLE users ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP;

-- Application "delete" becomes:
UPDATE users SET is_deleted = TRUE, deleted_at = CURRENT_TIMESTAMP
WHERE user_id = @id;

-- 4. Row-level audit trailing
CREATE TABLE users_audit (
    audit_id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    operation VARCHAR(10),
    old_data JSONB,
    new_data JSONB,
    changed_by VARCHAR(100),
    changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trigger to capture all changes
CREATE OR REPLACE FUNCTION audit_users()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO users_audit (user_id, operation, old_data, new_data, changed_by)
    VALUES (
        COALESCE(OLD.user_id, NEW.user_id),
        TG_OP,
        row_to_json(OLD),
        row_to_json(NEW),
        current_user
    );
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_audit_trigger
AFTER INSERT OR UPDATE OR DELETE ON users
FOR EACH ROW EXECUTE FUNCTION audit_users();
```

#### Access Control Improvements

```sql
-- 5. Separate read-only and read-write roles
CREATE ROLE app_readonly;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO app_readonly;

CREATE ROLE app_readwrite;
GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO app_readwrite;
-- Note: No DELETE permission!

CREATE ROLE app_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO app_admin;

-- 6. Require approval for destructive queries
-- Use a review workflow:
-- - Developer submits DELETE/UPDATE query
-- - DBA reviews and executes
-- - All changes logged

-- 7. Delayed deletion
CREATE TABLE deletion_queue (
    queue_id BIGSERIAL PRIMARY KEY,
    table_name VARCHAR(100),
    where_clause TEXT,
    requested_by VARCHAR(100),
    requested_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    approved_by VARCHAR(100),
    approved_at TIMESTAMP,
    executed_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'PENDING'
);

-- Process deletions after approval + waiting period
```

---

### Q44: Your team wants to add a new column to a table with 100 million rows. How do you do it safely?

**Situation:** The `orders` table has 100M rows. You need to add a `discount_code VARCHAR(50)` column. The application cannot have downtime.

**Answer:**

#### Risk Assessment

```sql
-- Check current table size and activity
SELECT 
    pg_size_pretty(pg_total_relation_size('orders')) as total_size,
    n_live_tup as row_count
FROM pg_stat_user_tables WHERE tablename = 'orders';

-- Check for locks
SELECT * FROM pg_locks WHERE relation = 'orders'::regclass;

-- Check replication lag if applicable
SELECT client_addr, state, sent_lsn, write_lsn, flush_lsn, replay_lsn
FROM pg_stat_replication;
```

#### Safe Approach: PostgreSQL

```sql
-- PostgreSQL: Adding nullable column is instant (metadata only)
-- This is SAFE and fast:
ALTER TABLE orders ADD COLUMN discount_code VARCHAR(50);
-- Completes in milliseconds, no table rewrite

-- Adding with DEFAULT (PostgreSQL 11+) is also fast:
ALTER TABLE orders ADD COLUMN discount_code VARCHAR(50) DEFAULT '';
-- Also instant - default stored in catalog, not written to rows

-- DANGEROUS: Adding NOT NULL with DEFAULT (pre-PG11)
-- This rewrites entire table:
ALTER TABLE orders ADD COLUMN discount_code VARCHAR(50) NOT NULL DEFAULT '';
-- Would lock table for hours!
```

#### Safe Approach: MySQL

```sql
-- MySQL 8.0+: Instant ADD COLUMN (usually)
ALTER TABLE orders ADD COLUMN discount_code VARCHAR(50), ALGORITHM=INSTANT;
-- If INSTANT not possible, falls back to INPLACE

-- Check if operation will be instant:
SELECT * FROM information_schema.innodb_tables 
WHERE name = 'mydb/orders';

-- For older MySQL or if INSTANT fails:
-- Use pt-online-schema-change (Percona)
pt-online-schema-change \
  --alter "ADD COLUMN discount_code VARCHAR(50)" \
  --execute \
  D=mydb,t=orders

-- This:
-- 1. Creates new table with new schema
-- 2. Creates triggers to sync changes
-- 3. Copies data in batches
-- 4. Swaps tables atomically
```

#### Safe Approach: Adding NOT NULL Column

```sql
-- Step 1: Add nullable column (fast)
ALTER TABLE orders ADD COLUMN discount_code VARCHAR(50);

-- Step 2: Backfill in batches (during low traffic)
DO $$
DECLARE
    batch_size INT := 10000;
    updated INT;
BEGIN
    LOOP
        UPDATE orders 
        SET discount_code = ''
        WHERE discount_code IS NULL
          AND order_id IN (
              SELECT order_id FROM orders 
              WHERE discount_code IS NULL 
              LIMIT batch_size
          );
        
        GET DIAGNOSTICS updated = ROW_COUNT;
        RAISE NOTICE 'Updated % rows', updated;
        
        EXIT WHEN updated = 0;
        
        -- Brief pause to reduce load
        PERFORM pg_sleep(0.1);
    END LOOP;
END $$;

-- Step 3: Add NOT NULL constraint (fast, just validation)
ALTER TABLE orders ALTER COLUMN discount_code SET NOT NULL;

-- Alternative Step 3: Add constraint without validation first
ALTER TABLE orders ADD CONSTRAINT chk_discount_code_not_null 
    CHECK (discount_code IS NOT NULL) NOT VALID;

-- Then validate separately (can be interrupted)
ALTER TABLE orders VALIDATE CONSTRAINT chk_discount_code_not_null;
```

#### Adding Index Safely

```sql
-- NEVER do this on production:
CREATE INDEX idx_orders_discount ON orders(discount_code);
-- Locks table for entire duration!

-- ALWAYS use CONCURRENTLY (PostgreSQL):
CREATE INDEX CONCURRENTLY idx_orders_discount ON orders(discount_code);
-- Takes longer but doesn't block reads/writes

-- MySQL equivalent:
ALTER TABLE orders ADD INDEX idx_discount (discount_code), ALGORITHM=INPLACE, LOCK=NONE;
```

#### Rollback Plan

```sql
-- If something goes wrong, have rollback ready:

-- Simple rollback:
ALTER TABLE orders DROP COLUMN discount_code;

-- If data was migrated to new column:
ALTER TABLE orders DROP COLUMN discount_code;
-- Application uses old column again
```

---

### Q45: How would you debug a database connection leak in your application?

**Situation:** Production alerts show database connections climbing steadily until max_connections is reached, then the application crashes. Restarting temporarily fixes it.

**Answer:**

#### Immediate Diagnosis

```sql
-- 1. Check current connection status
-- PostgreSQL:
SELECT 
    count(*) as total,
    count(*) FILTER (WHERE state = 'idle') as idle,
    count(*) FILTER (WHERE state = 'active') as active,
    count(*) FILTER (WHERE state = 'idle in transaction') as idle_in_txn,
    max(now() - backend_start) as oldest_connection
FROM pg_stat_activity
WHERE datname = current_database();

-- MySQL:
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';
SHOW PROCESSLIST;

-- 2. Identify connection sources
-- PostgreSQL:
SELECT 
    client_addr,
    usename,
    application_name,
    count(*) as connections,
    count(*) FILTER (WHERE state = 'idle') as idle
FROM pg_stat_activity
WHERE datname = current_database()
GROUP BY client_addr, usename, application_name
ORDER BY connections DESC;

-- 3. Find long-running idle connections
SELECT 
    pid,
    usename,
    application_name,
    client_addr,
    state,
    now() - state_change as idle_duration,
    now() - backend_start as connection_age,
    query
FROM pg_stat_activity
WHERE state = 'idle'
  AND now() - state_change > interval '10 minutes'
ORDER BY idle_duration DESC;

-- 4. Find idle transactions (major red flag!)
SELECT 
    pid,
    now() - xact_start as transaction_duration,
    query
FROM pg_stat_activity
WHERE state = 'idle in transaction'
ORDER BY xact_start;
```

#### Common Causes & Fixes

```python
# Cause 1: Not closing connections after use
# BAD:
def get_user(user_id):
    conn = database.connect()  # Opens connection
    result = conn.execute("SELECT * FROM users WHERE id = %s", user_id)
    return result  # Connection never closed!

# GOOD:
def get_user(user_id):
    conn = database.connect()
    try:
        result = conn.execute("SELECT * FROM users WHERE id = %s", user_id)
        return result.fetchone()
    finally:
        conn.close()  # Always close

# BETTER: Use context manager
def get_user(user_id):
    with database.connect() as conn:
        result = conn.execute("SELECT * FROM users WHERE id = %s", user_id)
        return result.fetchone()
    # Automatically closed
```

```python
# Cause 2: Not closing connections on exceptions
# BAD:
def process_order(order_id):
    conn = pool.get_connection()
    conn.begin()
    
    # If this raises an exception...
    process_payment(order_id)  # Exception here!
    
    conn.commit()
    pool.return_connection(conn)  # Never reached!

# GOOD:
def process_order(order_id):
    conn = pool.get_connection()
    try:
        conn.begin()
        process_payment(order_id)
        conn.commit()
    except Exception:
        conn.rollback()
        raise
    finally:
        pool.return_connection(conn)  # Always returned
```

```python
# Cause 3: Connection pool misconfiguration
# Check pool settings:
pool_config = {
    'min_connections': 5,      # Minimum idle connections
    'max_connections': 20,     # Maximum connections
    'max_idle_time': 300,      # Connections idle > 5 min are closed
    'connection_timeout': 30,   # Wait max 30s for connection
    'max_lifetime': 3600,      # Recycle connections after 1 hour
}

# HikariCP (Java) settings:
hikari:
  maximum-pool-size: 20
  minimum-idle: 5
  idle-timeout: 300000
  max-lifetime: 1800000
  connection-timeout: 30000
  leak-detection-threshold: 60000  # Alert if connection held > 60s
```

#### Emergency Mitigation

```sql
-- Kill idle connections older than 30 minutes
-- PostgreSQL:
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'idle'
  AND now() - state_change > interval '30 minutes'
  AND pid != pg_backend_pid();

-- Set statement timeout globally
ALTER DATABASE mydb SET statement_timeout = '30s';

-- Set idle session timeout (PostgreSQL 14+)
ALTER DATABASE mydb SET idle_session_timeout = '10min';

-- Kill idle-in-transaction connections
SELECT pg_terminate_backend(pid)
FROM pg_stat_activity
WHERE state = 'idle in transaction'
  AND now() - xact_start > interval '5 minutes';
```

#### Monitoring Implementation

```sql
-- Create monitoring view
CREATE VIEW connection_monitor AS
SELECT 
    now() as check_time,
    (SELECT setting::int FROM pg_settings WHERE name = 'max_connections') as max_conn,
    count(*) as total_conn,
    count(*) FILTER (WHERE state = 'active') as active,
    count(*) FILTER (WHERE state = 'idle') as idle,
    count(*) FILTER (WHERE state = 'idle in transaction') as idle_in_txn,
    round(100.0 * count(*) / 
        (SELECT setting::int FROM pg_settings WHERE name = 'max_connections'), 2) as pct_used
FROM pg_stat_activity;

-- Alert query for monitoring system
SELECT 
    CASE 
        WHEN pct_used > 80 THEN 'CRITICAL'
        WHEN pct_used > 60 THEN 'WARNING'
        ELSE 'OK'
    END as status,
    *
FROM connection_monitor;
```

---

### Situation-Based Questions Quick Reference

| Situation | Key Investigation | Common Fix |
|-----------|-------------------|------------|
| Sudden slowdown | Check running queries, locks, stats | ANALYZE, add index, kill blocking query |
| Table too large | Partition strategy, archive policy | Implement partitioning, archive old data |
| Race conditions | Identify concurrent access patterns | Optimistic locking, atomic operations |
| Multi-tenant isolation | Assess compliance needs | Hybrid: shared + dedicated schemas |
| Report degradation | Compare execution plans | Update stats, adjust indexes |
| No backups | Immediate manual backup | Implement PITR + automated backups |
| Global scale | Latency requirements per region | Read replicas or regional databases |
| Accidental deletion | Check backup/PITR availability | Restore from PITR, implement soft delete |
| Schema migration | Table size, lock requirements | Online schema change tools |
| Connection leak | Connection state distribution | Fix code, configure pool timeouts |

---

## Summary

This guide covers essential database design concepts for experienced engineers:

1. **Normalization** - Understand trade-offs between normalized and denormalized designs
2. **Schema Design** - Apply proper keys, constraints, and relationships
3. **Indexing** - Strategic index selection for query optimization
4. **Transactions** - ACID compliance and isolation levels
5. **Partitioning/Sharding** - Scaling strategies for large datasets
6. **Real-World Design** - Practical schema patterns for common systems
7. **Database Selection** - Choose the right SQL database for your needs
8. **Situation-Based Problems** - Handle real-world production scenarios

**Key Interview Tips:**
- Always explain trade-offs in your design decisions
- Consider write vs read workload when designing
- Think about scale and future requirements
- Mention monitoring and maintenance strategies
- Discuss data consistency requirements

---

*Last Updated: February 2026*
