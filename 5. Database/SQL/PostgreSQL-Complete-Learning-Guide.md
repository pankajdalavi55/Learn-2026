# PostgreSQL Complete Learning Guide
## From Fundamentals to Advanced Concepts

---

## Table of Contents

1. [Introduction to PostgreSQL](#1-introduction-to-postgresql)
2. [Installation & Setup](#2-installation--setup)
3. [Database Fundamentals](#3-database-fundamentals)
4. [Data Types](#4-data-types)
5. [CRUD Operations](#5-crud-operations)
6. [Querying Data](#6-querying-data)
7. [Joins & Relationships](#7-joins--relationships)
8. [Subqueries & CTEs](#8-subqueries--ctes)
9. [Aggregation & Grouping](#9-aggregation--grouping)
10. [Indexes & Performance](#10-indexes--performance)
11. [Views & Materialized Views](#11-views--materialized-views)
12. [Functions & Stored Procedures](#12-functions--stored-procedures)
13. [Triggers & Rules](#13-triggers--rules)
14. [Transactions & Concurrency](#14-transactions--concurrency)
15. [User Management & Security](#15-user-management--security)
16. [Backup & Recovery](#16-backup--recovery)
17. [Replication & High Availability](#17-replication--high-availability)
18. [Performance Tuning](#18-performance-tuning)
19. [Advanced Features](#19-advanced-features)
20. [Best Practices](#20-best-practices)

---

## 1. Introduction to PostgreSQL

### What is PostgreSQL?

PostgreSQL (often called "Postgres") is a powerful, open-source object-relational database system with over 35 years of active development. It has earned a strong reputation for reliability, feature robustness, and performance.

### Key Features

| Feature | Description |
|---------|-------------|
| **Open Source** | Free under PostgreSQL License (similar to MIT) |
| **ACID Compliant** | Full transaction support with MVCC |
| **Extensible** | Custom data types, functions, operators |
| **Standards Compliant** | Strong SQL standard compliance |
| **Advanced Data Types** | JSON/JSONB, Arrays, Hstore, Range types |
| **Full-Text Search** | Built-in powerful text search |
| **Geospatial** | PostGIS extension for GIS |
| **Replication** | Built-in streaming replication |

### PostgreSQL vs Other Databases

| Aspect | PostgreSQL | MySQL | SQL Server |
|--------|------------|-------|------------|
| License | PostgreSQL (free) | GPL/Commercial | Commercial |
| SQL Compliance | Excellent | Good | Good |
| JSON Support | Excellent (JSONB) | Good | Good |
| Extensions | 100+ extensions | Limited | Limited |
| Full-Text Search | Excellent | Good | Good |
| Partitioning | Native | Native | Native |
| Window Functions | Full support | Full (8.0+) | Full |
| Best For | Complex queries, analytics | Web apps | Enterprise/.NET |

### PostgreSQL Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT APPLICATIONS                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  psql    │  │  JDBC    │  │   PHP    │  │  Python  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    POSTMASTER PROCESS                       │
│         (Main daemon - spawns backend processes)            │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│   Backend    │   │   Backend    │   │   Backend    │
│   Process    │   │   Process    │   │   Process    │
│  (Per Conn)  │   │  (Per Conn)  │   │  (Per Conn)  │
└──────────────┘   └──────────────┘   └──────────────┘
        │                   │                   │
        └───────────────────┼───────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    SHARED MEMORY                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │Shared Buffer│  │   WAL       │  │  Lock       │         │
│  │   Pool      │  │   Buffers   │  │  Tables     │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    STORAGE LAYER                            │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │ Data Files  │  │  WAL Files  │  │  Temp Files │         │
│  │ (Tables,    │  │  (Write-    │  │             │         │
│  │  Indexes)   │  │   Ahead     │  │             │         │
│  │             │  │   Log)      │  │             │         │
│  └─────────────┘  └─────────────┘  └─────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### Background Processes

```
┌─────────────────────────────────────────────────────────────┐
│                  BACKGROUND PROCESSES                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────┐  Writes dirty buffers to disk         │
│  │ Background Writer│                                       │
│  └─────────────────┘                                        │
│                                                             │
│  ┌─────────────────┐  Writes WAL to archive                │
│  │ WAL Writer      │                                        │
│  └─────────────────┘                                        │
│                                                             │
│  ┌─────────────────┐  Creates periodic checkpoints         │
│  │ Checkpointer    │                                        │
│  └─────────────────┘                                        │
│                                                             │
│  ┌─────────────────┐  VACUUM and ANALYZE                   │
│  │ Autovacuum      │                                        │
│  └─────────────────┘                                        │
│                                                             │
│  ┌─────────────────┐  Collects statistics                  │
│  │ Stats Collector │                                        │
│  └─────────────────┘                                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Installation & Setup

### Installing on Windows

```powershell
# Using Chocolatey
choco install postgresql

# Or download installer from:
# https://www.postgresql.org/download/windows/

# After installation, connect using psql
psql -U postgres

# Or use pgAdmin GUI tool
```

### Installing on Linux (Ubuntu/Debian)

```bash
# Add PostgreSQL repository
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -

# Update and install
sudo apt update
sudo apt install postgresql-16 postgresql-contrib-16

# Start PostgreSQL
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Switch to postgres user and connect
sudo -u postgres psql
```

### Installing on macOS

```bash
# Using Homebrew
brew install postgresql@16

# Start PostgreSQL
brew services start postgresql@16

# Create default database
createdb

# Connect
psql
```

### Docker Installation (Recommended for Development)

```yaml
# docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:16
    container_name: postgres_dev
    environment:
      POSTGRES_USER: developer
      POSTGRES_PASSWORD: devpassword
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql

volumes:
  postgres_data:
```

```bash
# Start PostgreSQL container
docker-compose up -d

# Connect to PostgreSQL in container
docker exec -it postgres_dev psql -U developer -d myapp
```

### Initial Configuration

```sql
-- Connect as postgres superuser
sudo -u postgres psql

-- Create a new database
CREATE DATABASE myapp;

-- Create a new user
CREATE USER developer WITH PASSWORD 'password123';

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE myapp TO developer;

-- Connect to the database
\c myapp

-- Grant schema privileges
GRANT ALL ON SCHEMA public TO developer;
```

### PostgreSQL Configuration (postgresql.conf)

```ini
# Connection Settings
listen_addresses = '*'          # Listen on all interfaces
port = 5432
max_connections = 200

# Memory Settings
shared_buffers = 2GB           # 25% of RAM
effective_cache_size = 6GB     # 75% of RAM
work_mem = 64MB                # Per-operation memory
maintenance_work_mem = 512MB   # For VACUUM, CREATE INDEX

# WAL Settings
wal_level = replica
max_wal_size = 2GB
min_wal_size = 512MB
checkpoint_completion_target = 0.9

# Query Planner
random_page_cost = 1.1         # For SSDs (4.0 for HDD)
effective_io_concurrency = 200 # For SSDs

# Logging
log_destination = 'stderr'
logging_collector = on
log_directory = 'log'
log_filename = 'postgresql-%Y-%m-%d_%H%M%S.log'
log_statement = 'ddl'          # Log DDL statements
log_min_duration_statement = 1000  # Log slow queries (ms)

# Autovacuum
autovacuum = on
autovacuum_max_workers = 3
autovacuum_naptime = 1min
```

### pg_hba.conf (Authentication)

```ini
# TYPE  DATABASE  USER       ADDRESS         METHOD
local   all       postgres                   peer
local   all       all                        peer
host    all       all        127.0.0.1/32    scram-sha-256
host    all       all        ::1/128         scram-sha-256
host    all       all        0.0.0.0/0       scram-sha-256
```

### psql Commands Reference

```sql
-- Connect to database
psql -U username -d database -h hostname -p port

-- Common psql commands
\l              -- List databases
\c dbname       -- Connect to database
\dt             -- List tables
\dt+            -- List tables with size
\d tablename    -- Describe table
\di             -- List indexes
\dv             -- List views
\df             -- List functions
\du             -- List users/roles
\dn             -- List schemas
\dx             -- List extensions
\timing         -- Toggle timing of commands
\x              -- Toggle expanded output
\q              -- Quit
\?              -- Help on psql commands
\h              -- Help on SQL commands

-- Execute SQL file
\i filename.sql

-- Output to file
\o output.txt
SELECT * FROM users;
\o
```

---

## 3. Database Fundamentals

### Creating Databases

```sql
-- Create a simple database
CREATE DATABASE myapp;

-- Create with specific settings
CREATE DATABASE myapp
    WITH 
    OWNER = developer
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1;

-- Create from template
CREATE DATABASE myapp_test TEMPLATE myapp;
```

### Schemas

```sql
-- Create schema (namespace for objects)
CREATE SCHEMA sales;
CREATE SCHEMA hr AUTHORIZATION developer;

-- Create table in schema
CREATE TABLE sales.orders (
    id SERIAL PRIMARY KEY,
    customer_id INT,
    total DECIMAL(10, 2)
);

-- Set search path
SET search_path TO sales, public;

-- Show current search path
SHOW search_path;

-- Permanently set search path for user
ALTER USER developer SET search_path TO sales, public;

-- List schemas
\dn
SELECT schema_name FROM information_schema.schemata;

-- Drop schema
DROP SCHEMA sales;
DROP SCHEMA sales CASCADE;  -- Drop with all objects
```

### Creating Tables

```sql
-- Basic table creation
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Table with multiple constraints
CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    stock_quantity INT DEFAULT 0 CHECK (stock_quantity >= 0),
    category_id INT,
    is_active BOOLEAN DEFAULT TRUE,
    metadata JSONB,
    tags TEXT[],  -- Array type
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_name ON products USING gin(to_tsvector('english', name));
CREATE INDEX idx_products_metadata ON products USING gin(metadata);

-- Table with foreign keys
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(12, 2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' 
        CHECK (status IN ('pending', 'processing', 'shipped', 'delivered', 'cancelled'))
);

-- Junction table for many-to-many
CREATE TABLE order_items (
    order_id INT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES products(product_id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL,
    PRIMARY KEY (order_id, product_id)
);
```

### Altering Tables

```sql
-- Add column
ALTER TABLE users ADD COLUMN phone VARCHAR(20);

-- Add column with default
ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT FALSE;

-- Drop column
ALTER TABLE users DROP COLUMN phone;

-- Rename column
ALTER TABLE users RENAME COLUMN phone TO phone_number;

-- Change column type
ALTER TABLE users ALTER COLUMN phone TYPE VARCHAR(30);

-- Set/drop NOT NULL
ALTER TABLE users ALTER COLUMN email SET NOT NULL;
ALTER TABLE users ALTER COLUMN email DROP NOT NULL;

-- Set/drop default
ALTER TABLE users ALTER COLUMN is_active SET DEFAULT TRUE;
ALTER TABLE users ALTER COLUMN is_active DROP DEFAULT;

-- Add constraint
ALTER TABLE products ADD CONSTRAINT price_positive CHECK (price >= 0);
ALTER TABLE products ADD CONSTRAINT unique_sku UNIQUE (sku);

-- Add foreign key
ALTER TABLE orders ADD CONSTRAINT fk_orders_user 
    FOREIGN KEY (user_id) REFERENCES users(id);

-- Drop constraint
ALTER TABLE products DROP CONSTRAINT price_positive;

-- Rename table
ALTER TABLE users RENAME TO app_users;

-- Change table owner
ALTER TABLE users OWNER TO developer;
```

### Viewing Table Structure

```sql
-- Describe table
\d users
\d+ users  -- More details

-- Get detailed column information
SELECT 
    column_name,
    data_type,
    character_maximum_length,
    is_nullable,
    column_default
FROM information_schema.columns
WHERE table_name = 'users';

-- Show table constraints
SELECT conname, contype, pg_get_constraintdef(oid) 
FROM pg_constraint 
WHERE conrelid = 'users'::regclass;

-- Show indexes
SELECT indexname, indexdef 
FROM pg_indexes 
WHERE tablename = 'users';

-- Table size
SELECT pg_size_pretty(pg_total_relation_size('users'));
```

---

## 4. Data Types

### Numeric Types

```sql
CREATE TABLE numeric_examples (
    -- Integer types
    small_col SMALLINT,         -- -32,768 to 32,767 (2 bytes)
    int_col INTEGER,            -- -2.1B to 2.1B (4 bytes)
    big_col BIGINT,             -- -9.2Q to 9.2Q (8 bytes)
    
    -- Auto-increment
    id SERIAL,                  -- Auto-increment INTEGER
    id_big BIGSERIAL,           -- Auto-increment BIGINT
    
    -- Exact numeric
    price NUMERIC(10, 2),       -- 10 digits, 2 decimal
    salary DECIMAL(12, 2),      -- Same as NUMERIC
    
    -- Floating point (approximate)
    float_col REAL,             -- 4 bytes, 6 decimal digits
    double_col DOUBLE PRECISION, -- 8 bytes, 15 decimal digits
    
    -- Boolean
    is_active BOOLEAN           -- TRUE, FALSE, NULL
);

-- Best Practices:
-- Use INTEGER for IDs, counts
-- Use BIGINT for large numbers
-- Use NUMERIC for money (never REAL/DOUBLE!)
-- Use BOOLEAN for flags
```

### Character Types

```sql
CREATE TABLE text_examples (
    -- Fixed-length (padded with spaces)
    country_code CHAR(2),       -- Always 2 chars
    
    -- Variable-length with limit
    username VARCHAR(50),       -- Max 50 chars
    email VARCHAR(255),
    
    -- Variable-length unlimited
    description TEXT,           -- No length limit
    content TEXT,
    
    -- Note: VARCHAR and TEXT have same performance in PostgreSQL
    -- VARCHAR(n) is just for validation
);
```

### Date and Time Types

```sql
CREATE TABLE datetime_examples (
    -- Date only
    birth_date DATE,            -- '2024-03-15'
    
    -- Time without timezone
    start_time TIME,            -- '14:30:00'
    
    -- Time with timezone
    start_time_tz TIME WITH TIME ZONE,  -- '14:30:00+05:30'
    
    -- Timestamp without timezone
    created_at TIMESTAMP,       -- '2024-03-15 14:30:00'
    
    -- Timestamp with timezone (RECOMMENDED)
    created_at_tz TIMESTAMPTZ,  -- '2024-03-15 14:30:00+00'
    
    -- Interval
    duration INTERVAL           -- '2 hours 30 minutes'
);

-- Date/Time Functions
SELECT 
    CURRENT_DATE,                           -- Today's date
    CURRENT_TIME,                           -- Current time
    CURRENT_TIMESTAMP,                      -- Current timestamp
    NOW(),                                  -- Same as CURRENT_TIMESTAMP
    LOCALTIME,                              -- Time without TZ
    LOCALTIMESTAMP,                         -- Timestamp without TZ
    
    -- Extract parts
    EXTRACT(YEAR FROM NOW()),
    EXTRACT(MONTH FROM NOW()),
    EXTRACT(DAY FROM NOW()),
    EXTRACT(HOUR FROM NOW()),
    EXTRACT(DOW FROM NOW()),                -- Day of week (0=Sunday)
    EXTRACT(EPOCH FROM NOW()),              -- Unix timestamp
    
    -- Date arithmetic
    NOW() + INTERVAL '7 days',
    NOW() - INTERVAL '1 month',
    NOW() + '2 hours'::INTERVAL,
    
    -- Date difference
    AGE('2024-03-15', '2020-01-01'),       -- Interval difference
    '2024-03-15'::DATE - '2024-03-01'::DATE, -- Days difference (integer)
    
    -- Formatting
    TO_CHAR(NOW(), 'YYYY-MM-DD'),
    TO_CHAR(NOW(), 'Month DD, YYYY'),
    TO_CHAR(NOW(), 'HH24:MI:SS'),
    
    -- Parsing
    TO_DATE('15-03-2024', 'DD-MM-YYYY'),
    TO_TIMESTAMP('2024-03-15 14:30:00', 'YYYY-MM-DD HH24:MI:SS'),
    
    -- Truncation
    DATE_TRUNC('month', NOW()),             -- First of month
    DATE_TRUNC('year', NOW()),              -- First of year
    DATE_TRUNC('hour', NOW())               -- Start of hour
;
```

### JSON Types

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200),
    -- JSON: Stores as text, parses on each access
    attributes JSON,
    -- JSONB: Binary format, faster queries, supports indexing (RECOMMENDED)
    metadata JSONB
);

-- Insert JSON data
INSERT INTO products (name, metadata) VALUES 
    ('Laptop', '{"brand": "Dell", "ram": 16, "storage": "512GB"}'),
    ('Phone', '{"brand": "Apple", "model": "iPhone 15", "colors": ["black", "white"]}');

-- Query JSON data
SELECT 
    name,
    metadata->>'brand' AS brand,            -- Extract as text
    metadata->'ram' AS ram,                 -- Extract as JSON
    metadata#>>'{colors,0}' AS first_color, -- Path extraction
    jsonb_array_length(metadata->'colors') AS color_count
FROM products;

-- Filter by JSON value
SELECT * FROM products 
WHERE metadata->>'brand' = 'Dell';

SELECT * FROM products 
WHERE metadata @> '{"brand": "Dell"}';  -- Contains

SELECT * FROM products 
WHERE metadata ? 'colors';              -- Has key

SELECT * FROM products 
WHERE metadata ?| ARRAY['brand', 'model'];  -- Has any key

SELECT * FROM products 
WHERE metadata ?& ARRAY['brand', 'ram'];    -- Has all keys

-- JSONB operators
SELECT 
    metadata || '{"warranty": "2 years"}'::JSONB,  -- Concatenate
    metadata - 'ram',                              -- Remove key
    metadata #- '{colors,0}'                       -- Remove by path
FROM products;

-- JSON functions
SELECT 
    jsonb_pretty(metadata),                 -- Pretty print
    jsonb_typeof(metadata->'ram'),          -- Get type
    jsonb_object_keys(metadata),            -- Get keys
    jsonb_each(metadata),                   -- Key-value pairs
    jsonb_array_elements(metadata->'colors') -- Expand array
FROM products WHERE metadata ? 'colors';

-- Update JSON
UPDATE products 
SET metadata = jsonb_set(metadata, '{ram}', '32')
WHERE id = 1;

-- Create GIN index for JSONB
CREATE INDEX idx_metadata ON products USING gin(metadata);
CREATE INDEX idx_brand ON products USING btree((metadata->>'brand'));
```

### Array Types

```sql
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200),
    tags TEXT[],               -- Array of text
    ratings INTEGER[]          -- Array of integers
);

-- Insert arrays
INSERT INTO articles (title, tags, ratings) VALUES 
    ('PostgreSQL Guide', ARRAY['database', 'sql', 'tutorial'], ARRAY[5, 4, 5]),
    ('Python Basics', '{"programming", "python"}', '{4, 4, 3}');

-- Query arrays
SELECT 
    title,
    tags[1] AS first_tag,           -- 1-indexed!
    tags[1:2] AS first_two_tags,    -- Slice
    array_length(tags, 1) AS tag_count,
    array_to_string(tags, ', ') AS tags_str
FROM articles;

-- Array operators
SELECT * FROM articles WHERE 'sql' = ANY(tags);      -- Contains element
SELECT * FROM articles WHERE tags @> ARRAY['sql'];   -- Contains array
SELECT * FROM articles WHERE tags && ARRAY['sql', 'python'];  -- Overlaps
SELECT * FROM articles WHERE tags <@ ARRAY['database', 'sql', 'tutorial', 'extra'];  -- Contained by

-- Array functions
SELECT 
    array_cat(tags, ARRAY['new']),                   -- Concatenate
    array_append(tags, 'new'),                       -- Append
    array_prepend('new', tags),                      -- Prepend
    array_remove(tags, 'sql'),                       -- Remove element
    array_position(tags, 'sql'),                     -- Find position
    unnest(tags)                                     -- Expand to rows
FROM articles WHERE id = 1;

-- Create GIN index for arrays
CREATE INDEX idx_tags ON articles USING gin(tags);
```

### Range Types

```sql
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    room_id INT,
    during TSTZRANGE,           -- Timestamp with timezone range
    date_range DATERANGE        -- Date range
);

-- Insert ranges
INSERT INTO reservations (room_id, during, date_range) VALUES 
    (1, '[2024-03-15 14:00, 2024-03-15 16:00)', '[2024-03-15, 2024-03-16)'),
    (1, '[2024-03-15 17:00, 2024-03-15 19:00)', '[2024-03-15, 2024-03-16)');

-- Range operators
SELECT * FROM reservations 
WHERE during @> '2024-03-15 15:00'::TIMESTAMPTZ;     -- Contains point

SELECT * FROM reservations 
WHERE during && '[2024-03-15 15:00, 2024-03-15 18:00)'::TSTZRANGE;  -- Overlaps

-- Prevent overlapping reservations
ALTER TABLE reservations 
ADD CONSTRAINT no_overlap EXCLUDE USING gist (room_id WITH =, during WITH &&);

-- Range functions
SELECT 
    lower(during),              -- Start
    upper(during),              -- End
    isempty(during),           -- Is empty
    during * '[2024-03-15 15:00, 2024-03-15 20:00)'::TSTZRANGE  -- Intersection
FROM reservations;
```

### Other Types

```sql
-- UUID
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT,
    data JSONB
);

-- Network types
CREATE TABLE servers (
    id SERIAL PRIMARY KEY,
    ip_address INET,            -- IP address
    mac_address MACADDR,        -- MAC address
    network CIDR                -- Network with mask
);

INSERT INTO servers (ip_address, network) VALUES 
    ('192.168.1.100', '192.168.1.0/24');

SELECT * FROM servers WHERE ip_address << '192.168.1.0/24';  -- Contained in network

-- Geometric types
CREATE TABLE locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    coordinates POINT,
    area POLYGON
);

INSERT INTO locations (name, coordinates) VALUES 
    ('Office', POINT(40.7128, -74.0060));

-- Full geometric support via PostGIS extension
CREATE EXTENSION postgis;

-- Enum type
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'shipped', 'delivered', 'cancelled');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    status order_status DEFAULT 'pending'
);

-- Composite type
CREATE TYPE address AS (
    street VARCHAR(200),
    city VARCHAR(100),
    zip_code VARCHAR(20),
    country VARCHAR(50)
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    shipping_address address,
    billing_address address
);

INSERT INTO customers (name, shipping_address) VALUES 
    ('John', ROW('123 Main St', 'New York', '10001', 'USA'));

SELECT (shipping_address).city FROM customers;
```

---

## 5. CRUD Operations

### INSERT - Creating Data

```sql
-- Insert single row
INSERT INTO users (username, email, password_hash)
VALUES ('john_doe', 'john@example.com', 'hashed_password');

-- Insert multiple rows
INSERT INTO users (username, email, password_hash) VALUES
    ('jane_doe', 'jane@example.com', 'hash1'),
    ('bob_smith', 'bob@example.com', 'hash2'),
    ('alice_wong', 'alice@example.com', 'hash3');

-- Insert from SELECT
INSERT INTO users_backup (username, email, password_hash)
SELECT username, email, password_hash FROM users WHERE is_active = TRUE;

-- Insert with RETURNING (get inserted data)
INSERT INTO users (username, email, password_hash)
VALUES ('new_user', 'new@example.com', 'hash')
RETURNING id, username, created_at;

-- Insert with ON CONFLICT (upsert)
INSERT INTO users (username, email, password_hash)
VALUES ('john_doe', 'john_new@example.com', 'new_hash')
ON CONFLICT (username) DO UPDATE SET 
    email = EXCLUDED.email,
    password_hash = EXCLUDED.password_hash;

-- Insert with ON CONFLICT DO NOTHING
INSERT INTO users (username, email, password_hash)
VALUES ('john_doe', 'duplicate@example.com', 'hash')
ON CONFLICT (username) DO NOTHING;

-- Insert with ON CONFLICT on constraint
INSERT INTO products (sku, name, price)
VALUES ('SKU001', 'Product', 99.99)
ON CONFLICT ON CONSTRAINT products_sku_key DO UPDATE SET 
    price = EXCLUDED.price;

-- Insert with default values
INSERT INTO users (username, email, password_hash)
VALUES ('test', 'test@example.com', 'hash')
RETURNING *;

-- Insert JSON
INSERT INTO products (name, metadata)
VALUES ('Laptop', '{"brand": "Dell", "ram": 16}');

-- Insert array
INSERT INTO articles (title, tags)
VALUES ('PostgreSQL', ARRAY['database', 'sql']);
```

### SELECT - Reading Data

```sql
-- Select all columns
SELECT * FROM users;

-- Select specific columns
SELECT username, email FROM users;

-- Select with alias
SELECT 
    username AS user_name,
    email AS email_address,
    first_name || ' ' || last_name AS full_name
FROM users;

-- Select distinct
SELECT DISTINCT status FROM orders;
SELECT DISTINCT ON (customer_id) * FROM orders ORDER BY customer_id, order_date DESC;

-- Select with WHERE
SELECT * FROM users WHERE is_active = TRUE;

-- Select with multiple conditions
SELECT * FROM products
WHERE price > 100 
  AND category_id = 5
  AND (stock_quantity > 0 OR is_preorder = TRUE);

-- Select with IN
SELECT * FROM users WHERE id IN (1, 2, 3, 4, 5);
SELECT * FROM users WHERE id = ANY(ARRAY[1, 2, 3, 4, 5]);

-- Select with BETWEEN
SELECT * FROM orders 
WHERE order_date BETWEEN '2024-01-01' AND '2024-12-31';

-- Select with pattern matching
SELECT * FROM users WHERE email LIKE '%@gmail.com';
SELECT * FROM users WHERE email ILIKE '%@GMAIL.COM';  -- Case-insensitive
SELECT * FROM users WHERE username ~ '^[a-z]+$';      -- Regex
SELECT * FROM users WHERE username ~* '^[A-Z]+$';    -- Case-insensitive regex

-- Select with NULL
SELECT * FROM users WHERE phone IS NULL;
SELECT * FROM users WHERE phone IS NOT NULL;

-- Select with ORDER BY
SELECT * FROM products ORDER BY price ASC;
SELECT * FROM products ORDER BY NULLS FIRST;
SELECT * FROM products ORDER BY created_at DESC NULLS LAST;

-- Select with LIMIT and OFFSET
SELECT * FROM products ORDER BY id LIMIT 10;
SELECT * FROM products ORDER BY id LIMIT 10 OFFSET 20;

-- Select with FETCH (SQL standard)
SELECT * FROM products ORDER BY id FETCH FIRST 10 ROWS ONLY;
SELECT * FROM products ORDER BY id OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY;

-- Select with CASE
SELECT 
    name,
    price,
    CASE 
        WHEN price < 50 THEN 'Budget'
        WHEN price < 200 THEN 'Mid-range'
        ELSE 'Premium'
    END AS price_category
FROM products;

-- Select with COALESCE and NULLIF
SELECT 
    COALESCE(nickname, username, email) AS display_name,
    NULLIF(discount, 0) AS discount
FROM users;
```

### UPDATE - Modifying Data

```sql
-- Update single column
UPDATE users SET email = 'new_email@example.com' WHERE id = 1;

-- Update multiple columns
UPDATE users 
SET 
    email = 'updated@example.com',
    is_verified = TRUE,
    updated_at = NOW()
WHERE id = 1;

-- Update with calculation
UPDATE products SET price = price * 1.10 WHERE category_id = 5;

-- Update with RETURNING
UPDATE products 
SET price = price * 0.9 
WHERE category_id = 5
RETURNING id, name, price;

-- Update with FROM (join)
UPDATE orders o
SET status = 'cancelled'
FROM users u
WHERE o.user_id = u.id AND u.is_active = FALSE;

-- Update with subquery
UPDATE products 
SET category_id = (SELECT id FROM categories WHERE name = 'Electronics')
WHERE name LIKE '%Phone%';

-- Update JSON
UPDATE products 
SET metadata = metadata || '{"sale": true}'::JSONB
WHERE id = 1;

UPDATE products 
SET metadata = jsonb_set(metadata, '{price}', '999')
WHERE id = 1;

-- Update array
UPDATE articles 
SET tags = array_append(tags, 'featured')
WHERE id = 1;

-- Update with CASE
UPDATE products
SET stock_status = CASE
    WHEN stock_quantity = 0 THEN 'out_of_stock'
    WHEN stock_quantity < 10 THEN 'low_stock'
    ELSE 'in_stock'
END;
```

### DELETE - Removing Data

```sql
-- Delete specific rows
DELETE FROM users WHERE id = 1;

-- Delete with multiple conditions
DELETE FROM orders WHERE status = 'cancelled' AND order_date < '2023-01-01';

-- Delete with RETURNING
DELETE FROM users WHERE is_deleted = TRUE RETURNING *;

-- Delete with USING (join)
DELETE FROM orders o
USING users u
WHERE o.user_id = u.id AND u.is_deleted = TRUE;

-- Delete with subquery
DELETE FROM order_items 
WHERE order_id IN (
    SELECT order_id FROM orders WHERE status = 'cancelled'
);

-- Delete all rows
DELETE FROM temp_data;

-- TRUNCATE - faster way to delete all rows
TRUNCATE TABLE temp_data;
TRUNCATE TABLE temp_data RESTART IDENTITY;  -- Reset sequences
TRUNCATE TABLE orders, order_items CASCADE;  -- Multiple tables with FK
```

---

## 6. Querying Data

### Comparison Operators

```sql
SELECT * FROM products WHERE
    -- Equality
    status = 'active'
    
    -- Not equal
    AND status != 'deleted'
    AND status <> 'archived'
    
    -- Greater/Less than
    AND price > 100
    AND price >= 100
    AND stock < 50
    AND stock <= 50
    
    -- BETWEEN (inclusive)
    AND price BETWEEN 100 AND 500
    
    -- IN
    AND category_id IN (1, 2, 3)
    
    -- NOT IN
    AND status NOT IN ('deleted', 'archived')
    
    -- NULL checks
    AND deleted_at IS NULL
    AND description IS NOT NULL
    
    -- Pattern matching
    AND name LIKE 'iPhone%'         -- Case-sensitive
    AND name ILIKE 'iphone%'        -- Case-insensitive
    AND name SIMILAR TO '%(Pro|Max)%'  -- SQL regex
    
    -- POSIX regex
    AND email ~ '^[a-z]+@'          -- Case-sensitive
    AND email ~* '^[A-Z]+@'         -- Case-insensitive
    AND email !~ '^admin'           -- Does not match
;
```

### String Functions

```sql
SELECT 
    -- Length
    LENGTH('Hello'),              -- 5
    CHAR_LENGTH('Hello'),         -- 5
    OCTET_LENGTH('Hello'),        -- Bytes
    
    -- Case conversion
    UPPER('hello'),               -- 'HELLO'
    LOWER('HELLO'),               -- 'hello'
    INITCAP('hello world'),       -- 'Hello World'
    
    -- Concatenation
    'Hello' || ' ' || 'World',    -- 'Hello World'
    CONCAT('Hello', ' ', 'World'), -- 'Hello World'
    CONCAT_WS(', ', 'a', 'b', 'c'), -- 'a, b, c'
    
    -- Substring
    SUBSTRING('Hello World' FROM 1 FOR 5),  -- 'Hello'
    SUBSTRING('Hello World', 7),            -- 'World'
    LEFT('Hello', 2),                       -- 'He'
    RIGHT('Hello', 2),                      -- 'lo'
    
    -- Trim
    TRIM('  hello  '),            -- 'hello'
    LTRIM('  hello'),             -- 'hello'
    RTRIM('hello  '),             -- 'hello'
    TRIM(BOTH 'x' FROM 'xxhelloxx'), -- 'hello'
    
    -- Replace
    REPLACE('Hello World', 'World', 'PostgreSQL'),
    TRANSLATE('hello', 'el', 'ip'),  -- 'hippo'
    
    -- Position
    POSITION('World' IN 'Hello World'),  -- 7
    STRPOS('Hello World', 'World'),      -- 7
    
    -- Padding
    LPAD('42', 5, '0'),            -- '00042'
    RPAD('Hi', 5, '*'),            -- 'Hi***'
    
    -- Reverse
    REVERSE('Hello'),              -- 'olleH'
    
    -- Split
    SPLIT_PART('a,b,c', ',', 2),   -- 'b'
    STRING_TO_ARRAY('a,b,c', ','), -- {a,b,c}
    
    -- Format
    FORMAT('Hello %s, you have %s messages', 'John', 5),
    TO_CHAR(1234567.89, 'FM9,999,999.99')  -- '1,234,567.89'
;
```

### Numeric Functions

```sql
SELECT 
    -- Rounding
    ROUND(3.14159, 2),     -- 3.14
    CEIL(3.14),            -- 4
    CEILING(3.14),         -- 4
    FLOOR(3.94),           -- 3
    TRUNC(3.14159, 2),     -- 3.14
    
    -- Absolute and Sign
    ABS(-5),               -- 5
    SIGN(-5),              -- -1
    SIGN(5),               -- 1
    
    -- Power and Root
    POWER(2, 8),           -- 256
    SQRT(16),              -- 4
    CBRT(27),              -- 3 (cube root)
    
    -- Modulo
    MOD(17, 5),            -- 2
    17 % 5,                -- 2
    
    -- Random
    RANDOM(),              -- Random 0-1
    FLOOR(RANDOM() * 100), -- Random 0-99
    
    -- Greatest/Least
    GREATEST(1, 5, 3),     -- 5
    LEAST(1, 5, 3),        -- 1
    
    -- Logarithm
    LN(10),                -- Natural log
    LOG(100),              -- Base 10 log
    LOG(2, 8),             -- Log base 2 of 8 = 3
    
    -- Trigonometric
    PI(),
    DEGREES(PI()),         -- 180
    RADIANS(180),          -- PI
    SIN(PI()/2),          -- 1
    COS(0),               -- 1
    
    -- Division
    DIV(17, 5)            -- Integer division = 3
;
```

### Conditional Functions

```sql
-- CASE expression
SELECT 
    name,
    CASE 
        WHEN price < 50 THEN 'Budget'
        WHEN price < 200 THEN 'Mid-range'
        WHEN price < 1000 THEN 'Premium'
        ELSE 'Luxury'
    END AS tier
FROM products;

-- Simple CASE
SELECT 
    CASE status
        WHEN 'pending' THEN 'Awaiting Payment'
        WHEN 'processing' THEN 'Being Prepared'
        WHEN 'shipped' THEN 'On the Way'
        WHEN 'delivered' THEN 'Completed'
        ELSE 'Unknown'
    END AS status_text
FROM orders;

-- COALESCE (first non-NULL)
SELECT COALESCE(nickname, username, email) AS display_name FROM users;

-- NULLIF (return NULL if equal)
SELECT NULLIF(discount, 0) AS discount FROM orders;

-- GREATEST / LEAST
SELECT GREATEST(price, min_price, 0) AS effective_price FROM products;
SELECT LEAST(quantity, max_allowed) AS order_quantity FROM cart;

-- CASE in aggregate
SELECT 
    COUNT(*) AS total,
    COUNT(CASE WHEN status = 'active' THEN 1 END) AS active,
    SUM(CASE WHEN is_premium THEN price ELSE 0 END) AS premium_revenue
FROM products;
```

---

## 7. Joins & Relationships

### INNER JOIN

```sql
-- Returns only matching rows from both tables
SELECT 
    c.name AS customer_name,
    o.id AS order_id,
    o.total
FROM customers c
INNER JOIN orders o ON c.id = o.customer_id;

-- Using USING (when column names match)
SELECT *
FROM orders
INNER JOIN customers USING (customer_id);

-- Natural join (automatically joins on same-named columns)
SELECT * FROM orders NATURAL JOIN customers;
```

### LEFT JOIN

```sql
-- Returns all rows from left table, matching from right
SELECT 
    c.name AS customer_name,
    o.id AS order_id,
    COALESCE(o.total, 0) AS total
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id;

-- Find customers with no orders
SELECT c.*
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id
WHERE o.id IS NULL;
```

### RIGHT JOIN

```sql
-- Returns all rows from right table, matching from left
SELECT 
    c.name AS customer_name,
    o.id AS order_id,
    o.total
FROM customers c
RIGHT JOIN orders o ON c.id = o.customer_id;
```

### FULL OUTER JOIN

```sql
-- Returns all rows from both tables
SELECT 
    c.name AS customer_name,
    o.id AS order_id,
    o.total
FROM customers c
FULL OUTER JOIN orders o ON c.id = o.customer_id;

-- Find unmatched rows from either side
SELECT *
FROM customers c
FULL OUTER JOIN orders o ON c.id = o.customer_id
WHERE c.id IS NULL OR o.id IS NULL;
```

### CROSS JOIN

```sql
-- Cartesian product - all combinations
SELECT 
    c.name AS customer,
    p.name AS product
FROM customers c
CROSS JOIN products p;

-- Equivalent:
SELECT * FROM customers, products;
```

### Self Join

```sql
-- Table referencing itself
SELECT 
    e.name AS employee,
    m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;

-- Find employees with same salary
SELECT 
    e1.name,
    e2.name,
    e1.salary
FROM employees e1
JOIN employees e2 ON e1.salary = e2.salary AND e1.id < e2.id;
```

### LATERAL Join

```sql
-- LATERAL allows subquery to reference columns from preceding FROM items
SELECT 
    c.name,
    recent_orders.total,
    recent_orders.order_date
FROM customers c
CROSS JOIN LATERAL (
    SELECT total, order_date
    FROM orders
    WHERE customer_id = c.id
    ORDER BY order_date DESC
    LIMIT 3
) AS recent_orders;

-- Get top N per group efficiently
SELECT 
    d.name AS department,
    top_employees.name,
    top_employees.salary
FROM departments d
CROSS JOIN LATERAL (
    SELECT name, salary
    FROM employees
    WHERE department_id = d.id
    ORDER BY salary DESC
    LIMIT 3
) AS top_employees;
```

### Multiple Joins

```sql
SELECT 
    o.id AS order_id,
    c.name AS customer_name,
    p.name AS product_name,
    oi.quantity,
    oi.unit_price
FROM orders o
INNER JOIN customers c ON o.customer_id = c.id
INNER JOIN order_items oi ON o.id = oi.order_id
INNER JOIN products p ON oi.product_id = p.id
WHERE o.order_date >= '2024-01-01'
ORDER BY o.order_date DESC;
```

---

## 8. Subqueries & CTEs

### Subqueries in WHERE

```sql
-- Scalar subquery
SELECT * FROM products 
WHERE price > (SELECT AVG(price) FROM products);

-- IN subquery
SELECT * FROM customers 
WHERE id IN (
    SELECT DISTINCT customer_id FROM orders WHERE total > 1000
);

-- EXISTS subquery
SELECT * FROM customers c
WHERE EXISTS (
    SELECT 1 FROM orders o 
    WHERE o.customer_id = c.id AND o.total > 500
);

-- NOT EXISTS
SELECT * FROM customers c
WHERE NOT EXISTS (
    SELECT 1 FROM orders o WHERE o.customer_id = c.id
);

-- ANY / ALL
SELECT * FROM products
WHERE price > ANY (SELECT price FROM products WHERE category_id = 5);

SELECT * FROM products
WHERE price > ALL (SELECT price FROM products WHERE category_id = 5);
```

### Subqueries in FROM

```sql
-- Derived table
SELECT 
    category_name,
    avg_price,
    product_count
FROM (
    SELECT 
        c.name AS category_name,
        AVG(p.price) AS avg_price,
        COUNT(*) AS product_count
    FROM products p
    JOIN categories c ON p.category_id = c.id
    GROUP BY c.name
) AS category_stats
WHERE product_count > 5
ORDER BY avg_price DESC;
```

### Subqueries in SELECT

```sql
-- Scalar subquery
SELECT 
    name,
    price,
    price - (SELECT AVG(price) FROM products) AS price_vs_avg,
    (SELECT COUNT(*) FROM order_items WHERE product_id = p.id) AS times_ordered
FROM products p;
```

### Common Table Expressions (CTEs)

```sql
-- Basic CTE
WITH high_value_orders AS (
    SELECT * FROM orders WHERE total > 500
)
SELECT 
    c.name,
    COUNT(h.id) AS high_value_count
FROM customers c
JOIN high_value_orders h ON c.id = h.customer_id
GROUP BY c.id, c.name;

-- Multiple CTEs
WITH 
order_totals AS (
    SELECT 
        customer_id,
        COUNT(*) AS order_count,
        SUM(total) AS total_spent
    FROM orders
    GROUP BY customer_id
),
customer_tiers AS (
    SELECT 
        customer_id,
        CASE 
            WHEN total_spent >= 10000 THEN 'Gold'
            WHEN total_spent >= 5000 THEN 'Silver'
            ELSE 'Bronze'
        END AS tier
    FROM order_totals
)
SELECT 
    c.name,
    ot.order_count,
    ot.total_spent,
    ct.tier
FROM customers c
JOIN order_totals ot ON c.id = ot.customer_id
JOIN customer_tiers ct ON c.id = ct.customer_id
ORDER BY ot.total_spent DESC;

-- Recursive CTE
WITH RECURSIVE category_tree AS (
    -- Base case: top-level categories
    SELECT 
        id,
        name,
        parent_id,
        name AS path,
        0 AS level
    FROM categories
    WHERE parent_id IS NULL
    
    UNION ALL
    
    -- Recursive case
    SELECT 
        c.id,
        c.name,
        c.parent_id,
        ct.path || ' > ' || c.name,
        ct.level + 1
    FROM categories c
    INNER JOIN category_tree ct ON c.parent_id = ct.id
)
SELECT * FROM category_tree ORDER BY path;

-- Generate series with recursive CTE
WITH RECURSIVE dates AS (
    SELECT DATE '2024-01-01' AS date
    UNION ALL
    SELECT date + 1 FROM dates WHERE date < '2024-01-31'
)
SELECT * FROM dates;

-- Or use generate_series
SELECT generate_series(
    '2024-01-01'::DATE,
    '2024-01-31'::DATE,
    '1 day'::INTERVAL
) AS date;
```

### CTE with Data Modification (Writeable CTEs)

```sql
-- CTE with INSERT/UPDATE/DELETE
WITH deleted_orders AS (
    DELETE FROM orders 
    WHERE status = 'cancelled' AND order_date < '2023-01-01'
    RETURNING *
)
INSERT INTO orders_archive SELECT * FROM deleted_orders;

-- CTE with RETURNING
WITH new_customer AS (
    INSERT INTO customers (name, email)
    VALUES ('New Customer', 'new@example.com')
    RETURNING id
)
INSERT INTO orders (customer_id, total)
SELECT id, 0 FROM new_customer;
```

---

## 9. Aggregation & Grouping

### Aggregate Functions

```sql
SELECT 
    COUNT(*) AS total_products,
    COUNT(description) AS with_description,
    COUNT(DISTINCT category_id) AS categories,
    SUM(price) AS total_value,
    AVG(price) AS average_price,
    MIN(price) AS cheapest,
    MAX(price) AS most_expensive,
    STDDEV(price) AS price_std_dev,
    VARIANCE(price) AS price_variance
FROM products;

-- Array aggregation
SELECT 
    category_id,
    ARRAY_AGG(name ORDER BY name) AS products,
    ARRAY_AGG(DISTINCT name) AS unique_products
FROM products
GROUP BY category_id;

-- String aggregation
SELECT 
    category_id,
    STRING_AGG(name, ', ' ORDER BY name) AS product_list
FROM products
GROUP BY category_id;

-- JSON aggregation
SELECT 
    category_id,
    JSON_AGG(name) AS products_json,
    JSON_OBJECT_AGG(id, name) AS product_map
FROM products
GROUP BY category_id;

-- Boolean aggregation
SELECT 
    BOOL_AND(is_active) AS all_active,
    BOOL_OR(is_featured) AS any_featured
FROM products;
```

### GROUP BY

```sql
-- Basic grouping
SELECT 
    category_id,
    COUNT(*) AS product_count,
    AVG(price) AS avg_price
FROM products
GROUP BY category_id;

-- Group by multiple columns
SELECT 
    DATE_TRUNC('month', order_date) AS month,
    status,
    COUNT(*) AS orders,
    SUM(total) AS revenue
FROM orders
GROUP BY DATE_TRUNC('month', order_date), status
ORDER BY month, status;

-- Group by expression
SELECT 
    CASE 
        WHEN price < 50 THEN 'Budget'
        WHEN price < 200 THEN 'Mid-range'
        ELSE 'Premium'
    END AS price_tier,
    COUNT(*) AS count
FROM products
GROUP BY 1;  -- Reference by position
```

### HAVING

```sql
-- Filter groups
SELECT 
    category_id,
    COUNT(*) AS product_count,
    AVG(price) AS avg_price
FROM products
GROUP BY category_id
HAVING COUNT(*) > 5 AND AVG(price) > 100;
```

### GROUP BY with ROLLUP, CUBE, GROUPING SETS

```sql
-- ROLLUP - hierarchical subtotals
SELECT 
    EXTRACT(YEAR FROM order_date) AS year,
    EXTRACT(MONTH FROM order_date) AS month,
    SUM(total) AS revenue
FROM orders
GROUP BY ROLLUP(
    EXTRACT(YEAR FROM order_date),
    EXTRACT(MONTH FROM order_date)
);

-- CUBE - all combinations
SELECT 
    category_id,
    status,
    COUNT(*) AS count
FROM products
GROUP BY CUBE(category_id, status);

-- GROUPING SETS - specific combinations
SELECT 
    category_id,
    EXTRACT(YEAR FROM created_at) AS year,
    COUNT(*) AS count
FROM products
GROUP BY GROUPING SETS (
    (category_id, EXTRACT(YEAR FROM created_at)),
    (category_id),
    (EXTRACT(YEAR FROM created_at)),
    ()
);

-- Use GROUPING() to identify aggregation level
SELECT 
    CASE WHEN GROUPING(category_id) = 1 THEN 'All Categories' ELSE category_id::TEXT END,
    CASE WHEN GROUPING(status) = 1 THEN 'All Status' ELSE status END,
    COUNT(*)
FROM products
GROUP BY ROLLUP(category_id, status);
```

### Window Functions

```sql
-- ROW_NUMBER
SELECT 
    name,
    category_id,
    price,
    ROW_NUMBER() OVER (ORDER BY price DESC) AS overall_rank,
    ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC) AS category_rank
FROM products;

-- RANK and DENSE_RANK
SELECT 
    name,
    price,
    RANK() OVER (ORDER BY price DESC) AS rank_with_gaps,
    DENSE_RANK() OVER (ORDER BY price DESC) AS rank_no_gaps
FROM products;

-- Running totals
SELECT 
    order_date,
    total,
    SUM(total) OVER (ORDER BY order_date) AS running_total,
    SUM(total) OVER (
        ORDER BY order_date 
        ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW
    ) AS cumulative_sum
FROM orders;

-- Moving averages
SELECT 
    order_date,
    total,
    AVG(total) OVER (
        ORDER BY order_date 
        ROWS BETWEEN 6 PRECEDING AND CURRENT ROW
    ) AS moving_avg_7day
FROM orders;

-- LEAD and LAG
SELECT 
    order_date,
    total,
    LAG(total, 1) OVER (ORDER BY order_date) AS prev_total,
    LEAD(total, 1) OVER (ORDER BY order_date) AS next_total,
    total - LAG(total, 1) OVER (ORDER BY order_date) AS diff_from_prev
FROM orders;

-- FIRST_VALUE and LAST_VALUE
SELECT 
    category_id,
    name,
    price,
    FIRST_VALUE(name) OVER w AS most_expensive,
    LAST_VALUE(name) OVER w AS cheapest
FROM products
WINDOW w AS (
    PARTITION BY category_id 
    ORDER BY price DESC 
    RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
);

-- NTH_VALUE
SELECT 
    name,
    price,
    NTH_VALUE(name, 2) OVER (ORDER BY price DESC) AS second_expensive
FROM products;

-- NTILE - divide into buckets
SELECT 
    name,
    price,
    NTILE(4) OVER (ORDER BY price) AS price_quartile
FROM products;

-- PERCENT_RANK and CUME_DIST
SELECT 
    name,
    price,
    PERCENT_RANK() OVER (ORDER BY price) AS percentile,
    CUME_DIST() OVER (ORDER BY price) AS cumulative_dist
FROM products;

-- Named window
SELECT 
    name,
    category_id,
    price,
    ROW_NUMBER() OVER category_window AS row_num,
    SUM(price) OVER category_window AS category_total,
    AVG(price) OVER category_window AS category_avg
FROM products
WINDOW category_window AS (PARTITION BY category_id ORDER BY price);
```

---

## 10. Indexes & Performance

### Index Types

```sql
-- 1. B-tree (default) - equality, range queries
CREATE INDEX idx_email ON users(email);
CREATE INDEX idx_created ON users(created_at);

-- 2. Hash - equality only (rarely used)
CREATE INDEX idx_email_hash ON users USING hash(email);

-- 3. GIN (Generalized Inverted Index) - arrays, JSONB, full-text
CREATE INDEX idx_tags ON articles USING gin(tags);
CREATE INDEX idx_metadata ON products USING gin(metadata);
CREATE INDEX idx_search ON products USING gin(to_tsvector('english', name || ' ' || description));

-- 4. GiST (Generalized Search Tree) - geometric, range, full-text
CREATE INDEX idx_location ON places USING gist(coordinates);
CREATE INDEX idx_period ON reservations USING gist(during);

-- 5. SP-GiST (Space-Partitioned GiST) - non-balanced structures
CREATE INDEX idx_ip ON servers USING spgist(ip_address);

-- 6. BRIN (Block Range Index) - large ordered tables
CREATE INDEX idx_created_brin ON logs USING brin(created_at);

-- 7. Partial index (conditional)
CREATE INDEX idx_active_products ON products(name) WHERE is_active = TRUE;
CREATE INDEX idx_recent_orders ON orders(order_date) WHERE order_date > '2024-01-01';

-- 8. Expression/Functional index
CREATE INDEX idx_lower_email ON users(LOWER(email));
CREATE INDEX idx_year ON orders(EXTRACT(YEAR FROM order_date));
CREATE INDEX idx_jsonb_brand ON products((metadata->>'brand'));

-- 9. Covering index (INCLUDE)
CREATE INDEX idx_customer_orders ON orders(customer_id) INCLUDE (order_date, total);

-- 10. Unique index
CREATE UNIQUE INDEX idx_unique_email ON users(email);
CREATE UNIQUE INDEX idx_unique_sku ON products(sku) WHERE deleted_at IS NULL;

-- 11. Multi-column index
CREATE INDEX idx_category_price ON products(category_id, price);

-- 12. Concurrent index creation (doesn't lock table)
CREATE INDEX CONCURRENTLY idx_name ON products(name);
```

### Managing Indexes

```sql
-- List indexes
SELECT indexname, indexdef FROM pg_indexes WHERE tablename = 'users';
\di users  -- psql command

-- Index size
SELECT pg_size_pretty(pg_relation_size('idx_email'));

-- All indexes with sizes
SELECT
    schemaname,
    tablename,
    indexname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
ORDER BY pg_relation_size(indexrelid) DESC;

-- Rename index
ALTER INDEX idx_email RENAME TO idx_users_email;

-- Drop index
DROP INDEX idx_email;
DROP INDEX CONCURRENTLY idx_email;  -- Non-blocking
DROP INDEX IF EXISTS idx_email;

-- Rebuild index
REINDEX INDEX idx_email;
REINDEX TABLE users;
REINDEX DATABASE myapp;
```

### EXPLAIN - Query Analysis

```sql
-- Basic EXPLAIN
EXPLAIN SELECT * FROM orders WHERE customer_id = 100;

-- EXPLAIN ANALYZE (actually runs query)
EXPLAIN ANALYZE SELECT * FROM orders WHERE customer_id = 100;

-- Verbose output
EXPLAIN (ANALYZE, BUFFERS, FORMAT TEXT) 
SELECT * FROM orders WHERE customer_id = 100;

-- JSON format (easier to parse)
EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) 
SELECT * FROM orders WHERE customer_id = 100;

-- Understanding EXPLAIN output:
-- Seq Scan: Full table scan (no index)
-- Index Scan: Uses index, fetches rows from table
-- Index Only Scan: Uses covering index only
-- Bitmap Index Scan: Combines multiple indexes
-- Nested Loop: For each row in outer, scan inner
-- Hash Join: Build hash table, probe with other
-- Merge Join: Sort both, merge
```

### Query Optimization

```sql
-- 1. Use indexes effectively
-- Add indexes on WHERE, JOIN, ORDER BY columns

-- 2. Avoid functions on indexed columns
-- Bad:
SELECT * FROM users WHERE LOWER(email) = 'test@example.com';
-- Good: Create functional index
CREATE INDEX idx_lower_email ON users(LOWER(email));

-- 3. Use covering indexes
CREATE INDEX idx_covering ON orders(customer_id) INCLUDE (order_date, total);

-- 4. Use partial indexes
CREATE INDEX idx_active ON products(name) WHERE is_active = TRUE;

-- 5. Limit result sets
SELECT * FROM logs ORDER BY created_at DESC LIMIT 100;

-- 6. Use EXISTS instead of IN for large subqueries
-- Slower:
SELECT * FROM customers WHERE id IN (SELECT customer_id FROM orders);
-- Faster:
SELECT * FROM customers c WHERE EXISTS (SELECT 1 FROM orders o WHERE o.customer_id = c.id);

-- 7. Batch operations
-- Instead of many small inserts:
INSERT INTO logs VALUES (...), (...), (...);

-- 8. Use COPY for bulk loading
COPY products FROM '/tmp/products.csv' WITH CSV HEADER;

-- 9. Analyze tables for statistics
ANALYZE products;

-- 10. Update statistics after major changes
VACUUM ANALYZE products;
```

---

## 11. Views & Materialized Views

### Regular Views

```sql
-- Create view
CREATE VIEW active_products AS
SELECT id, name, price, stock_quantity
FROM products
WHERE is_active = TRUE AND stock_quantity > 0;

-- Use view
SELECT * FROM active_products WHERE price < 100;

-- View with joins
CREATE VIEW order_details AS
SELECT 
    o.id AS order_id,
    o.order_date,
    c.name AS customer_name,
    c.email AS customer_email,
    o.total,
    o.status
FROM orders o
JOIN customers c ON o.customer_id = c.id;

-- View with aggregation
CREATE VIEW category_stats AS
SELECT 
    c.id AS category_id,
    c.name AS category_name,
    COUNT(p.id) AS product_count,
    AVG(p.price) AS avg_price,
    MIN(p.price) AS min_price,
    MAX(p.price) AS max_price
FROM categories c
LEFT JOIN products p ON c.id = p.category_id
GROUP BY c.id, c.name;

-- Updatable view (simple views)
CREATE VIEW simple_products AS
SELECT id, name, price FROM products;

UPDATE simple_products SET price = 99.99 WHERE id = 1;

-- View with security barrier
CREATE VIEW user_data WITH (security_barrier) AS
SELECT id, username, email FROM users WHERE is_public = TRUE;
```

### Materialized Views

```sql
-- Materialized view (cached results)
CREATE MATERIALIZED VIEW monthly_sales AS
SELECT 
    DATE_TRUNC('month', order_date) AS month,
    COUNT(*) AS order_count,
    SUM(total) AS revenue,
    AVG(total) AS avg_order_value
FROM orders
WHERE status = 'completed'
GROUP BY DATE_TRUNC('month', order_date)
WITH DATA;

-- Create without data initially
CREATE MATERIALIZED VIEW sales_summary AS
SELECT ... WITH NO DATA;

-- Query like a regular table
SELECT * FROM monthly_sales WHERE month >= '2024-01-01';

-- Refresh materialized view
REFRESH MATERIALIZED VIEW monthly_sales;

-- Refresh concurrently (needs unique index)
CREATE UNIQUE INDEX idx_monthly_sales_month ON monthly_sales(month);
REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_sales;

-- Schedule refresh (using pg_cron extension)
SELECT cron.schedule('refresh_sales', '0 * * * *', 
    'REFRESH MATERIALIZED VIEW CONCURRENTLY monthly_sales');
```

### Managing Views

```sql
-- List views
\dv
SELECT viewname FROM pg_views WHERE schemaname = 'public';

-- View definition
\d+ active_products
SELECT pg_get_viewdef('active_products', TRUE);

-- Alter view
CREATE OR REPLACE VIEW active_products AS
SELECT id, name, price, stock_quantity, category_id
FROM products
WHERE is_active = TRUE;

-- Rename view
ALTER VIEW active_products RENAME TO available_products;

-- Change owner
ALTER VIEW active_products OWNER TO developer;

-- Drop view
DROP VIEW IF EXISTS active_products;
DROP VIEW active_products CASCADE;

-- Drop materialized view
DROP MATERIALIZED VIEW monthly_sales;
```

---

## 12. Functions & Stored Procedures

### Creating Functions

```sql
-- Simple function
CREATE OR REPLACE FUNCTION get_product_count(category_id_param INT)
RETURNS INT
LANGUAGE plpgsql
AS $$
DECLARE
    cnt INT;
BEGIN
    SELECT COUNT(*) INTO cnt 
    FROM products 
    WHERE category_id = category_id_param;
    RETURN cnt;
END;
$$;

-- Use function
SELECT get_product_count(5);
SELECT c.name, get_product_count(c.id) FROM categories c;

-- Function with multiple parameters
CREATE OR REPLACE FUNCTION calculate_discount(
    price DECIMAL,
    discount_percent INT
)
RETURNS DECIMAL
LANGUAGE plpgsql
IMMUTABLE  -- Same input always returns same output
AS $$
BEGIN
    RETURN price * (1 - discount_percent / 100.0);
END;
$$;

-- Function returning table
CREATE OR REPLACE FUNCTION get_top_products(limit_count INT)
RETURNS TABLE (
    product_id INT,
    product_name VARCHAR,
    total_sold BIGINT
)
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN QUERY
    SELECT 
        p.id,
        p.name,
        COALESCE(SUM(oi.quantity), 0)
    FROM products p
    LEFT JOIN order_items oi ON p.id = oi.product_id
    GROUP BY p.id, p.name
    ORDER BY COALESCE(SUM(oi.quantity), 0) DESC
    LIMIT limit_count;
END;
$$;

SELECT * FROM get_top_products(10);

-- Function with OUT parameters
CREATE OR REPLACE FUNCTION get_order_stats(
    customer_id_param INT,
    OUT total_orders INT,
    OUT total_spent DECIMAL,
    OUT avg_order DECIMAL
)
LANGUAGE plpgsql
AS $$
BEGIN
    SELECT 
        COUNT(*),
        SUM(total),
        AVG(total)
    INTO total_orders, total_spent, avg_order
    FROM orders
    WHERE customer_id = customer_id_param;
END;
$$;

SELECT * FROM get_order_stats(1);

-- SQL function (simpler syntax)
CREATE OR REPLACE FUNCTION active_products_count()
RETURNS BIGINT
LANGUAGE SQL
STABLE
AS $$
    SELECT COUNT(*) FROM products WHERE is_active = TRUE;
$$;
```

### Stored Procedures (PostgreSQL 11+)

```sql
-- Stored procedure (can use COMMIT/ROLLBACK)
CREATE OR REPLACE PROCEDURE process_order(
    p_customer_id INT,
    p_product_id INT,
    p_quantity INT
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_stock INT;
    v_price DECIMAL;
    v_order_id INT;
BEGIN
    -- Check stock
    SELECT stock_quantity, price INTO v_stock, v_price
    FROM products WHERE id = p_product_id FOR UPDATE;
    
    IF v_stock IS NULL THEN
        RAISE EXCEPTION 'Product not found';
    END IF;
    
    IF v_stock < p_quantity THEN
        RAISE EXCEPTION 'Insufficient stock: % available', v_stock;
    END IF;
    
    -- Create order
    INSERT INTO orders (customer_id, total, status)
    VALUES (p_customer_id, v_price * p_quantity, 'pending')
    RETURNING id INTO v_order_id;
    
    -- Add order item
    INSERT INTO order_items (order_id, product_id, quantity, unit_price)
    VALUES (v_order_id, p_product_id, p_quantity, v_price);
    
    -- Update stock
    UPDATE products 
    SET stock_quantity = stock_quantity - p_quantity
    WHERE id = p_product_id;
    
    COMMIT;
    
    RAISE NOTICE 'Order % created successfully', v_order_id;
END;
$$;

-- Call procedure
CALL process_order(1, 5, 2);

-- Procedure with transaction control
CREATE OR REPLACE PROCEDURE batch_archive_orders(batch_size INT)
LANGUAGE plpgsql
AS $$
DECLARE
    archived_count INT;
BEGIN
    LOOP
        WITH deleted AS (
            DELETE FROM orders
            WHERE status = 'completed' 
            AND order_date < NOW() - INTERVAL '1 year'
            LIMIT batch_size
            RETURNING *
        )
        INSERT INTO orders_archive SELECT * FROM deleted;
        
        GET DIAGNOSTICS archived_count = ROW_COUNT;
        
        COMMIT;  -- Commit after each batch
        
        EXIT WHEN archived_count < batch_size;
        
        -- Optional: Add delay between batches
        PERFORM pg_sleep(0.1);
    END LOOP;
END;
$$;
```

### PL/pgSQL Control Structures

```sql
CREATE OR REPLACE FUNCTION complex_example(input_value INT)
RETURNS TEXT
LANGUAGE plpgsql
AS $$
DECLARE
    result TEXT;
    counter INT := 0;
    rec RECORD;
BEGIN
    -- IF-THEN-ELSE
    IF input_value < 0 THEN
        result := 'Negative';
    ELSIF input_value = 0 THEN
        result := 'Zero';
    ELSE
        result := 'Positive';
    END IF;
    
    -- CASE
    result := CASE input_value
        WHEN 1 THEN 'One'
        WHEN 2 THEN 'Two'
        ELSE 'Other'
    END;
    
    -- Simple LOOP
    LOOP
        counter := counter + 1;
        EXIT WHEN counter >= 10;
    END LOOP;
    
    -- WHILE loop
    WHILE counter < 20 LOOP
        counter := counter + 1;
    END LOOP;
    
    -- FOR loop (integer)
    FOR i IN 1..10 LOOP
        counter := counter + i;
    END LOOP;
    
    -- FOR loop (reverse)
    FOR i IN REVERSE 10..1 LOOP
        counter := counter + i;
    END LOOP;
    
    -- FOR loop (query)
    FOR rec IN SELECT * FROM products WHERE price < 100 LOOP
        RAISE NOTICE 'Product: %', rec.name;
    END LOOP;
    
    -- FOREACH (array)
    DECLARE
        arr INT[] := ARRAY[1, 2, 3];
        elem INT;
    BEGIN
        FOREACH elem IN ARRAY arr LOOP
            counter := counter + elem;
        END LOOP;
    END;
    
    RETURN result;
END;
$$;
```

### Error Handling

```sql
CREATE OR REPLACE FUNCTION safe_divide(a DECIMAL, b DECIMAL)
RETURNS DECIMAL
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN a / b;
EXCEPTION
    WHEN division_by_zero THEN
        RAISE WARNING 'Division by zero, returning NULL';
        RETURN NULL;
    WHEN numeric_value_out_of_range THEN
        RAISE WARNING 'Numeric overflow';
        RETURN NULL;
    WHEN OTHERS THEN
        RAISE WARNING 'Unknown error: %', SQLERRM;
        RETURN NULL;
END;
$$;

-- Raise exceptions
CREATE OR REPLACE FUNCTION validate_order(quantity INT)
RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
    IF quantity <= 0 THEN
        RAISE EXCEPTION 'Quantity must be positive'
            USING ERRCODE = 'check_violation',
                  HINT = 'Provide a quantity greater than 0';
    END IF;
    
    IF quantity > 1000 THEN
        RAISE EXCEPTION 'Quantity % exceeds maximum allowed', quantity;
    END IF;
END;
$$;
```

### Managing Functions

```sql
-- List functions
\df
SELECT proname, prosrc FROM pg_proc WHERE pronamespace = 'public'::regnamespace;

-- View function definition
\df+ function_name
SELECT pg_get_functiondef('function_name'::regproc);

-- Drop function
DROP FUNCTION IF EXISTS get_product_count(INT);
DROP FUNCTION get_product_count;  -- If no other overloads

-- Drop procedure
DROP PROCEDURE IF EXISTS process_order(INT, INT, INT);
```

---

## 13. Triggers & Rules

### Creating Triggers

```sql
-- Trigger function
CREATE OR REPLACE FUNCTION update_modified_timestamp()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$;

-- Create trigger
CREATE TRIGGER trigger_users_updated
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_modified_timestamp();

-- BEFORE INSERT trigger
CREATE OR REPLACE FUNCTION normalize_email()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.email = LOWER(TRIM(NEW.email));
    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_normalize_email
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION normalize_email();

-- AFTER trigger for audit logging
CREATE OR REPLACE FUNCTION log_order_changes()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO order_audit (order_id, action, new_data, changed_at)
        VALUES (NEW.id, 'INSERT', row_to_json(NEW), NOW());
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO order_audit (order_id, action, old_data, new_data, changed_at)
        VALUES (NEW.id, 'UPDATE', row_to_json(OLD), row_to_json(NEW), NOW());
    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO order_audit (order_id, action, old_data, changed_at)
        VALUES (OLD.id, 'DELETE', row_to_json(OLD), NOW());
    END IF;
    RETURN COALESCE(NEW, OLD);
END;
$$;

CREATE TRIGGER trigger_order_audit
    AFTER INSERT OR UPDATE OR DELETE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION log_order_changes();

-- INSTEAD OF trigger (for views)
CREATE OR REPLACE FUNCTION update_order_view()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE orders SET status = NEW.status WHERE id = NEW.order_id;
    UPDATE customers SET name = NEW.customer_name WHERE id = NEW.customer_id;
    RETURN NEW;
END;
$$;

CREATE TRIGGER trigger_order_view_update
    INSTEAD OF UPDATE ON order_details
    FOR EACH ROW
    EXECUTE FUNCTION update_order_view();

-- Statement-level trigger
CREATE OR REPLACE FUNCTION log_bulk_operation()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO bulk_operation_log (table_name, operation, timestamp)
    VALUES (TG_TABLE_NAME, TG_OP, NOW());
    RETURN NULL;
END;
$$;

CREATE TRIGGER trigger_log_bulk
    AFTER INSERT OR UPDATE OR DELETE ON products
    FOR EACH STATEMENT
    EXECUTE FUNCTION log_bulk_operation();

-- Conditional trigger (WHEN)
CREATE TRIGGER trigger_price_change
    AFTER UPDATE ON products
    FOR EACH ROW
    WHEN (OLD.price IS DISTINCT FROM NEW.price)
    EXECUTE FUNCTION log_price_change();
```

### Managing Triggers

```sql
-- List triggers
\dS orders
SELECT trigger_name, event_manipulation, action_timing 
FROM information_schema.triggers 
WHERE event_object_table = 'orders';

-- Disable trigger
ALTER TABLE orders DISABLE TRIGGER trigger_order_audit;
ALTER TABLE orders DISABLE TRIGGER ALL;

-- Enable trigger
ALTER TABLE orders ENABLE TRIGGER trigger_order_audit;
ALTER TABLE orders ENABLE TRIGGER ALL;

-- Drop trigger
DROP TRIGGER IF EXISTS trigger_order_audit ON orders;
```

### Rules (Alternative to Triggers)

```sql
-- Rule for logging
CREATE RULE log_product_inserts AS
    ON INSERT TO products
    DO ALSO
    INSERT INTO product_log (product_id, action, timestamp)
    VALUES (NEW.id, 'INSERT', NOW());

-- Rule for redirecting
CREATE RULE redirect_old_orders AS
    ON INSERT TO orders
    WHERE NEW.order_date < '2020-01-01'
    DO INSTEAD
    INSERT INTO orders_archive VALUES (NEW.*);

-- Rule for preventing deletes
CREATE RULE prevent_delete AS
    ON DELETE TO important_data
    DO INSTEAD NOTHING;

-- Drop rule
DROP RULE IF EXISTS log_product_inserts ON products;
```

---

## 14. Transactions & Concurrency

### Transaction Basics

```sql
-- Start transaction
BEGIN;
-- Or
START TRANSACTION;

-- Execute queries
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- Commit
COMMIT;

-- Or rollback
ROLLBACK;
```

### Savepoints

```sql
BEGIN;

INSERT INTO orders (customer_id, total) VALUES (1, 100);

SAVEPOINT order_created;

INSERT INTO order_items (order_id, product_id, quantity) 
VALUES (currval('orders_id_seq'), 1, 5);

-- Oops, wrong product
ROLLBACK TO SAVEPOINT order_created;

-- Insert correct item
INSERT INTO order_items (order_id, product_id, quantity) 
VALUES (currval('orders_id_seq'), 2, 3);

COMMIT;

-- Release savepoint (optional)
RELEASE SAVEPOINT order_created;
```

### Isolation Levels

```sql
-- Check current isolation level
SHOW transaction_isolation;

-- Set isolation level for transaction
BEGIN ISOLATION LEVEL READ UNCOMMITTED;
BEGIN ISOLATION LEVEL READ COMMITTED;  -- Default
BEGIN ISOLATION LEVEL REPEATABLE READ;
BEGIN ISOLATION LEVEL SERIALIZABLE;

-- Or set during transaction
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- Set default for session
SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL SERIALIZABLE;
```

### Isolation Levels Explained

```
┌─────────────────────────────────────────────────────────────┐
│                   ISOLATION LEVELS                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  READ UNCOMMITTED  (Not actually supported - same as RC)    │
│  └── Can see uncommitted changes (dirty reads)              │
│                                                             │
│  READ COMMITTED (Default)                                   │
│  └── Only sees committed data                               │
│  └── Same query may return different results                │
│  └── Ideal for most web applications                        │
│                                                             │
│  REPEATABLE READ                                            │
│  └── Same query returns same results                        │
│  └── Uses MVCC snapshots                                    │
│  └── May serialize differently than serial execution        │
│                                                             │
│  SERIALIZABLE                                               │
│  └── Strictest isolation                                    │
│  └── Transactions appear to run serially                    │
│  └── May throw serialization errors (retry needed)          │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Locking

```sql
-- Row-level locks
SELECT * FROM accounts WHERE id = 1 FOR UPDATE;           -- Exclusive lock
SELECT * FROM accounts WHERE id = 1 FOR NO KEY UPDATE;    -- Weaker exclusive
SELECT * FROM accounts WHERE id = 1 FOR SHARE;            -- Shared lock
SELECT * FROM accounts WHERE id = 1 FOR KEY SHARE;        -- Weakest shared

-- Wait options
SELECT * FROM accounts WHERE id = 1 FOR UPDATE NOWAIT;    -- Fail immediately
SELECT * FROM accounts WHERE id = 1 FOR UPDATE SKIP LOCKED; -- Skip locked rows

-- Table-level locks
LOCK TABLE accounts IN ACCESS SHARE MODE;
LOCK TABLE accounts IN ROW SHARE MODE;
LOCK TABLE accounts IN ROW EXCLUSIVE MODE;
LOCK TABLE accounts IN SHARE UPDATE EXCLUSIVE MODE;
LOCK TABLE accounts IN SHARE MODE;
LOCK TABLE accounts IN SHARE ROW EXCLUSIVE MODE;
LOCK TABLE accounts IN EXCLUSIVE MODE;
LOCK TABLE accounts IN ACCESS EXCLUSIVE MODE;

-- Advisory locks (application-level)
SELECT pg_advisory_lock(12345);              -- Acquire exclusive lock
SELECT pg_advisory_xact_lock(12345);         -- Transaction-level lock
SELECT pg_try_advisory_lock(12345);          -- Non-blocking
SELECT pg_advisory_unlock(12345);            -- Release lock

-- Check locks
SELECT * FROM pg_locks WHERE relation = 'accounts'::regclass;
```

### Handling Deadlocks and Conflicts

```sql
-- PostgreSQL automatically detects and resolves deadlocks
-- One transaction will be aborted with:
-- ERROR: deadlock detected

-- For serializable isolation, handle serialization failures:
-- ERROR: could not serialize access due to concurrent update

-- Application pattern for retry:
-- BEGIN;
-- -- Try operation
-- -- On serialization_failure: ROLLBACK, wait, retry
-- COMMIT;

-- Lock timeout setting
SET lock_timeout = '10s';

-- Statement timeout
SET statement_timeout = '30s';
```

---

## 15. User Management & Security

### Creating Users and Roles

```sql
-- Create user (role with LOGIN)
CREATE USER developer WITH PASSWORD 'password123';

-- Create role without login
CREATE ROLE readonly;

-- Create role with options
CREATE ROLE admin WITH 
    LOGIN 
    PASSWORD 'adminpass'
    SUPERUSER
    CREATEDB
    CREATEROLE
    REPLICATION;

-- Create role with expiration
CREATE ROLE temp_user WITH 
    LOGIN 
    PASSWORD 'temppass'
    VALID UNTIL '2024-12-31';

-- Create role with connection limit
CREATE ROLE app_user WITH 
    LOGIN 
    PASSWORD 'apppass'
    CONNECTION LIMIT 10;
```

### Granting Privileges

```sql
-- Database privileges
GRANT CONNECT ON DATABASE myapp TO developer;
GRANT CREATE ON DATABASE myapp TO developer;

-- Schema privileges
GRANT USAGE ON SCHEMA public TO developer;
GRANT CREATE ON SCHEMA public TO developer;

-- Table privileges
GRANT SELECT ON products TO readonly;
GRANT SELECT, INSERT, UPDATE, DELETE ON products TO developer;
GRANT ALL PRIVILEGES ON products TO admin;

-- All tables in schema
GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO developer;

-- Future tables (default privileges)
ALTER DEFAULT PRIVILEGES IN SCHEMA public 
    GRANT SELECT ON TABLES TO readonly;

-- Column-level privileges
GRANT SELECT (id, name, price) ON products TO limited_user;
GRANT UPDATE (price) ON products TO price_editor;

-- Sequence privileges
GRANT USAGE, SELECT ON SEQUENCE products_id_seq TO developer;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO developer;

-- Function privileges
GRANT EXECUTE ON FUNCTION get_product_count(INT) TO developer;
GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA public TO developer;
```

### Role Membership

```sql
-- Grant role to user
GRANT readonly TO developer;
GRANT admin TO developer WITH ADMIN OPTION;  -- Can grant to others

-- Set role
SET ROLE readonly;
SET ROLE NONE;
RESET ROLE;

-- Inherit privileges automatically
ALTER ROLE developer INHERIT;

-- Don't inherit (must SET ROLE explicitly)
ALTER ROLE developer NOINHERIT;
```

### Revoking Privileges

```sql
-- Revoke specific privileges
REVOKE INSERT, UPDATE ON products FROM developer;

-- Revoke all privileges
REVOKE ALL PRIVILEGES ON products FROM developer;

-- Revoke role membership
REVOKE readonly FROM developer;

-- Revoke with CASCADE
REVOKE ALL ON products FROM developer CASCADE;
```

### Managing Users

```sql
-- List roles
\du
SELECT rolname, rolsuper, rolcreaterole, rolcreatedb 
FROM pg_roles;

-- Show privileges
\dp products
SELECT grantee, privilege_type 
FROM information_schema.role_table_grants 
WHERE table_name = 'products';

-- Change password
ALTER ROLE developer WITH PASSWORD 'newpassword';

-- Rename role
ALTER ROLE developer RENAME TO senior_developer;

-- Modify role options
ALTER ROLE developer CREATEDB;
ALTER ROLE developer NOCREATEDB;

-- Drop role
DROP ROLE IF EXISTS developer;

-- Drop owned objects before dropping role
DROP OWNED BY developer;
REASSIGN OWNED BY developer TO postgres;
DROP ROLE developer;
```

### Row-Level Security (RLS)

```sql
-- Enable RLS on table
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

-- Create policy
CREATE POLICY user_orders ON orders
    FOR ALL
    TO app_users
    USING (user_id = current_user_id());

-- Separate policies for operations
CREATE POLICY select_own_orders ON orders
    FOR SELECT
    USING (user_id = current_user_id());

CREATE POLICY insert_own_orders ON orders
    FOR INSERT
    WITH CHECK (user_id = current_user_id());

CREATE POLICY update_own_orders ON orders
    FOR UPDATE
    USING (user_id = current_user_id())
    WITH CHECK (user_id = current_user_id());

-- Policy for specific role
CREATE POLICY admin_all_orders ON orders
    FOR ALL
    TO admin
    USING (TRUE);

-- Drop policy
DROP POLICY IF EXISTS user_orders ON orders;

-- Force RLS even for table owner
ALTER TABLE orders FORCE ROW LEVEL SECURITY;
```

### Security Best Practices

```sql
-- 1. Use strong passwords
CREATE ROLE app WITH PASSWORD 'C0mpl3x!P@ssw0rd#2024';

-- 2. Grant minimum privileges
GRANT SELECT, INSERT, UPDATE ON app_tables TO app_user;

-- 3. Use schemas for organization
CREATE SCHEMA app_data;
GRANT USAGE ON SCHEMA app_data TO app_user;

-- 4. Use SSL connections
-- In postgresql.conf: ssl = on
-- In pg_hba.conf: hostssl all all 0.0.0.0/0 scram-sha-256

-- 5. Use password hashing
-- Default in PostgreSQL 14+: scram-sha-256

-- 6. Regular audits
SELECT usename, usesuper, usecreatedb 
FROM pg_user 
WHERE usesuper = TRUE;

-- 7. Enable connection logging
-- In postgresql.conf:
-- log_connections = on
-- log_disconnections = on
```

---

## 16. Backup & Recovery

### pg_dump - Logical Backup

```bash
# Backup single database (SQL format)
pg_dump -U postgres myapp > myapp_backup.sql

# Backup with custom format (compressed, can restore individual objects)
pg_dump -U postgres -Fc myapp > myapp_backup.dump

# Backup with directory format (parallel backup)
pg_dump -U postgres -Fd -j 4 myapp -f myapp_backup_dir

# Backup with tar format
pg_dump -U postgres -Ft myapp > myapp_backup.tar

# Backup specific tables
pg_dump -U postgres -t users -t orders myapp > tables_backup.sql

# Backup schema only (no data)
pg_dump -U postgres --schema-only myapp > myapp_schema.sql

# Backup data only (no schema)
pg_dump -U postgres --data-only myapp > myapp_data.sql

# Backup with compression
pg_dump -U postgres myapp | gzip > myapp_backup.sql.gz

# Backup all databases
pg_dumpall -U postgres > all_databases.sql

# Backup globals only (roles, tablespaces)
pg_dumpall -U postgres --globals-only > globals.sql
```

### pg_restore - Restoring Backups

```bash
# Restore SQL format
psql -U postgres -d myapp < myapp_backup.sql

# Restore custom format
pg_restore -U postgres -d myapp myapp_backup.dump

# Restore with parallel processing
pg_restore -U postgres -d myapp -j 4 myapp_backup.dump

# Restore specific table
pg_restore -U postgres -d myapp -t users myapp_backup.dump

# Restore schema only
pg_restore -U postgres -d myapp --schema-only myapp_backup.dump

# Restore data only
pg_restore -U postgres -d myapp --data-only myapp_backup.dump

# List contents of backup
pg_restore -l myapp_backup.dump

# Restore with drop/create
pg_restore -U postgres -d postgres --clean --create myapp_backup.dump
```

### Point-in-Time Recovery (PITR)

```ini
# postgresql.conf settings for WAL archiving
wal_level = replica
archive_mode = on
archive_command = 'cp %p /backup/wal_archive/%f'
```

```bash
# Create base backup
pg_basebackup -U postgres -D /backup/base -Ft -z -P

# For PITR, create recovery.signal and configure:
# In postgresql.conf:
# restore_command = 'cp /backup/wal_archive/%f %p'
# recovery_target_time = '2024-03-15 14:00:00'
# recovery_target_action = 'promote'
```

### Continuous Archiving

```bash
# pgBackRest - popular backup tool
# Install and configure pgbackrest.conf

# Full backup
pgbackrest --stanza=myapp --type=full backup

# Differential backup
pgbackrest --stanza=myapp --type=diff backup

# Incremental backup
pgbackrest --stanza=myapp --type=incr backup

# Restore
pgbackrest --stanza=myapp restore

# Point-in-time recovery
pgbackrest --stanza=myapp --type=time --target='2024-03-15 14:00:00' restore
```

### Backup Script Example

```bash
#!/bin/bash
# backup_postgres.sh

DB_NAME="myapp"
DB_USER="postgres"
BACKUP_DIR="/backup/postgres"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

mkdir -p $BACKUP_DIR

# Create backup
pg_dump -U $DB_USER -Fc $DB_NAME > $BACKUP_DIR/${DB_NAME}_${DATE}.dump

if [ $? -eq 0 ]; then
    echo "Backup completed: ${DB_NAME}_${DATE}.dump"
    
    # Upload to S3 (optional)
    aws s3 cp $BACKUP_DIR/${DB_NAME}_${DATE}.dump s3://my-bucket/postgres-backups/
    
    # Delete old backups
    find $BACKUP_DIR -name "*.dump" -mtime +$RETENTION_DAYS -delete
else
    echo "Backup failed!"
    exit 1
fi
```

### COPY Command

```sql
-- Export to CSV
COPY products TO '/tmp/products.csv' WITH CSV HEADER;

-- Export query results
COPY (SELECT * FROM products WHERE price > 100) TO '/tmp/expensive.csv' WITH CSV HEADER;

-- Import from CSV
COPY products FROM '/tmp/products.csv' WITH CSV HEADER;

-- Import with column mapping
COPY products (name, price, category_id) FROM '/tmp/products.csv' WITH CSV HEADER;

-- Using \copy in psql (client-side)
\copy products TO '/tmp/products.csv' WITH CSV HEADER
\copy products FROM '/tmp/products.csv' WITH CSV HEADER
```

---

## 17. Replication & High Availability

### Streaming Replication Setup

**Primary Server (postgresql.conf):**
```ini
listen_addresses = '*'
wal_level = replica
max_wal_senders = 10
wal_keep_size = 1GB
synchronous_commit = on
```

**Primary Server (pg_hba.conf):**
```ini
host replication repl_user 192.168.1.0/24 scram-sha-256
```

**Primary Server Setup:**
```sql
-- Create replication user
CREATE ROLE repl_user WITH REPLICATION LOGIN PASSWORD 'repl_password';
```

**Standby Server Setup:**
```bash
# Take base backup from primary
pg_basebackup -h primary_host -U repl_user -D /var/lib/postgresql/data -P -R

# The -R flag creates standby.signal and sets primary_conninfo
```

**Standby postgresql.conf:**
```ini
hot_standby = on
primary_conninfo = 'host=primary_host port=5432 user=repl_user password=repl_password'
```

### Monitoring Replication

```sql
-- On primary: Check connected standbys
SELECT 
    client_addr,
    state,
    sent_lsn,
    write_lsn,
    flush_lsn,
    replay_lsn,
    sync_state
FROM pg_stat_replication;

-- On standby: Check replication status
SELECT 
    pg_is_in_recovery(),
    pg_last_wal_receive_lsn(),
    pg_last_wal_replay_lsn(),
    pg_last_xact_replay_timestamp();

-- Replication lag
SELECT 
    EXTRACT(EPOCH FROM (NOW() - pg_last_xact_replay_timestamp())) AS lag_seconds;
```

### Synchronous Replication

```ini
# postgresql.conf on primary
synchronous_standby_names = 'standby1, standby2'
# Or for quorum:
synchronous_standby_names = 'FIRST 2 (standby1, standby2, standby3)'
```

### Logical Replication

```sql
-- On publisher
CREATE PUBLICATION my_publication FOR TABLE users, orders;
-- Or all tables:
CREATE PUBLICATION my_publication FOR ALL TABLES;

-- On subscriber
CREATE SUBSCRIPTION my_subscription
    CONNECTION 'host=publisher_host dbname=myapp user=repl_user password=repl_password'
    PUBLICATION my_publication;

-- Check subscription status
SELECT * FROM pg_stat_subscription;

-- Modify publication
ALTER PUBLICATION my_publication ADD TABLE products;

-- Drop
DROP SUBSCRIPTION my_subscription;
DROP PUBLICATION my_publication;
```

### Failover

```bash
# Promote standby to primary
pg_ctl promote -D /var/lib/postgresql/data

# Or via SQL (PostgreSQL 12+)
SELECT pg_promote();
```

### Patroni - HA Cluster Management

```yaml
# patroni.yml example
scope: postgres-cluster
name: node1

restapi:
  listen: 0.0.0.0:8008
  connect_address: node1:8008

etcd:
  hosts: etcd1:2379,etcd2:2379,etcd3:2379

bootstrap:
  dcs:
    ttl: 30
    loop_wait: 10
    retry_timeout: 10
    maximum_lag_on_failover: 1048576
    postgresql:
      use_pg_rewind: true
      parameters:
        wal_level: replica
        hot_standby: 'on'
        max_wal_senders: 10

postgresql:
  listen: 0.0.0.0:5432
  connect_address: node1:5432
  data_dir: /var/lib/postgresql/data
  authentication:
    replication:
      username: repl_user
      password: repl_password
    superuser:
      username: postgres
      password: postgres_password
```

---

## 18. Performance Tuning

### Key Configuration Parameters

```ini
# Memory Settings
shared_buffers = 4GB              # 25% of RAM
effective_cache_size = 12GB       # 75% of RAM
work_mem = 64MB                   # Per-operation
maintenance_work_mem = 1GB        # For VACUUM, CREATE INDEX
huge_pages = try

# WAL Settings
wal_buffers = 64MB
checkpoint_completion_target = 0.9
max_wal_size = 4GB
min_wal_size = 1GB

# Query Planner
random_page_cost = 1.1            # For SSDs (4.0 for HDD)
effective_io_concurrency = 200   # For SSDs
default_statistics_target = 100   # Increase for complex queries

# Connections
max_connections = 200
superuser_reserved_connections = 3

# Parallel Query
max_parallel_workers_per_gather = 4
max_parallel_workers = 8
max_parallel_maintenance_workers = 4

# Autovacuum
autovacuum_max_workers = 4
autovacuum_naptime = 1min
autovacuum_vacuum_scale_factor = 0.1
autovacuum_analyze_scale_factor = 0.05

# Logging
log_min_duration_statement = 1000  # Log queries > 1 second
log_checkpoints = on
log_lock_waits = on
log_temp_files = 0
```

### Monitoring Queries

```sql
-- Enable pg_stat_statements extension
CREATE EXTENSION pg_stat_statements;

-- Top queries by total time
SELECT 
    query,
    calls,
    total_exec_time / 1000 AS total_seconds,
    mean_exec_time / 1000 AS avg_seconds,
    rows
FROM pg_stat_statements
ORDER BY total_exec_time DESC
LIMIT 10;

-- Top queries by calls
SELECT query, calls, rows
FROM pg_stat_statements
ORDER BY calls DESC
LIMIT 10;

-- Reset statistics
SELECT pg_stat_statements_reset();
```

### Analyzing Performance

```sql
-- Table statistics
SELECT 
    schemaname,
    relname,
    n_live_tup,
    n_dead_tup,
    last_vacuum,
    last_autovacuum,
    last_analyze,
    last_autoanalyze
FROM pg_stat_user_tables;

-- Index usage
SELECT 
    schemaname,
    relname,
    indexrelname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;

-- Unused indexes
SELECT 
    schemaname,
    relname,
    indexrelname,
    pg_size_pretty(pg_relation_size(indexrelid)) AS index_size
FROM pg_stat_user_indexes
WHERE idx_scan = 0
ORDER BY pg_relation_size(indexrelid) DESC;

-- Table bloat estimate
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname || '.' || tablename)) AS total_size,
    pg_size_pretty(pg_relation_size(schemaname || '.' || tablename)) AS table_size,
    pg_size_pretty(pg_indexes_size(schemaname || '.' || tablename)) AS index_size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname || '.' || tablename) DESC;

-- Active queries
SELECT 
    pid,
    age(clock_timestamp(), query_start) AS duration,
    usename,
    query,
    state
FROM pg_stat_activity
WHERE state != 'idle'
ORDER BY query_start;

-- Waiting queries
SELECT 
    pid,
    usename,
    wait_event_type,
    wait_event,
    query
FROM pg_stat_activity
WHERE wait_event IS NOT NULL;

-- Locks
SELECT 
    pg_locks.pid,
    pg_class.relname,
    pg_locks.mode,
    pg_locks.granted
FROM pg_locks
JOIN pg_class ON pg_locks.relation = pg_class.oid
WHERE pg_locks.mode != 'AccessShareLock';
```

### Vacuum and Analyze

```sql
-- Manual vacuum
VACUUM products;
VACUUM FULL products;  -- Reclaims space, locks table
VACUUM ANALYZE products;

-- Analyze for statistics
ANALYZE products;
ANALYZE products(price, category_id);

-- Check autovacuum progress
SELECT * FROM pg_stat_progress_vacuum;

-- Set aggressive autovacuum for specific table
ALTER TABLE high_traffic_table SET (
    autovacuum_vacuum_scale_factor = 0.01,
    autovacuum_analyze_scale_factor = 0.01
);
```

### Connection Pooling (PgBouncer)

```ini
# pgbouncer.ini
[databases]
myapp = host=localhost port=5432 dbname=myapp

[pgbouncer]
listen_addr = *
listen_port = 6432
auth_file = /etc/pgbouncer/userlist.txt
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20
min_pool_size = 5
reserve_pool_size = 5
```

---

## 19. Advanced Features

### Full-Text Search

```sql
-- Add tsvector column
ALTER TABLE articles ADD COLUMN search_vector tsvector;

-- Populate search vector
UPDATE articles SET search_vector = 
    to_tsvector('english', coalesce(title, '') || ' ' || coalesce(content, ''));

-- Create GIN index
CREATE INDEX idx_articles_search ON articles USING gin(search_vector);

-- Create trigger to maintain search vector
CREATE FUNCTION update_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector = to_tsvector('english', 
        coalesce(NEW.title, '') || ' ' || coalesce(NEW.content, ''));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_search_vector
    BEFORE INSERT OR UPDATE ON articles
    FOR EACH ROW
    EXECUTE FUNCTION update_search_vector();

-- Search
SELECT title, content
FROM articles
WHERE search_vector @@ to_tsquery('english', 'postgresql & performance');

-- Ranking
SELECT 
    title,
    ts_rank(search_vector, to_tsquery('english', 'postgresql')) AS rank
FROM articles
WHERE search_vector @@ to_tsquery('english', 'postgresql')
ORDER BY rank DESC;

-- Highlight matched terms
SELECT 
    title,
    ts_headline('english', content, to_tsquery('english', 'postgresql'),
        'StartSel=<b>, StopSel=</b>') AS highlighted
FROM articles
WHERE search_vector @@ to_tsquery('english', 'postgresql');
```

### Partitioning

```sql
-- Range partitioning
CREATE TABLE orders (
    id BIGSERIAL,
    order_date DATE NOT NULL,
    customer_id INT,
    total DECIMAL(10, 2)
) PARTITION BY RANGE (order_date);

-- Create partitions
CREATE TABLE orders_2024_q1 PARTITION OF orders
    FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');
CREATE TABLE orders_2024_q2 PARTITION OF orders
    FOR VALUES FROM ('2024-04-01') TO ('2024-07-01');
CREATE TABLE orders_2024_q3 PARTITION OF orders
    FOR VALUES FROM ('2024-07-01') TO ('2024-10-01');
CREATE TABLE orders_2024_q4 PARTITION OF orders
    FOR VALUES FROM ('2024-10-01') TO ('2025-01-01');

-- Default partition (catches unmatched rows)
CREATE TABLE orders_default PARTITION OF orders DEFAULT;

-- List partitioning
CREATE TABLE customers (
    id SERIAL,
    name VARCHAR(100),
    country VARCHAR(2)
) PARTITION BY LIST (country);

CREATE TABLE customers_us PARTITION OF customers FOR VALUES IN ('US');
CREATE TABLE customers_uk PARTITION OF customers FOR VALUES IN ('UK', 'GB');
CREATE TABLE customers_eu PARTITION OF customers FOR VALUES IN ('DE', 'FR', 'IT', 'ES');

-- Hash partitioning
CREATE TABLE logs (
    id BIGSERIAL,
    message TEXT,
    created_at TIMESTAMP
) PARTITION BY HASH (id);

CREATE TABLE logs_0 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 0);
CREATE TABLE logs_1 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 1);
CREATE TABLE logs_2 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 2);
CREATE TABLE logs_3 PARTITION OF logs FOR VALUES WITH (MODULUS 4, REMAINDER 3);

-- Attach/detach partitions
ALTER TABLE orders DETACH PARTITION orders_2024_q1;
ALTER TABLE orders ATTACH PARTITION orders_2024_q1
    FOR VALUES FROM ('2024-01-01') TO ('2024-04-01');
```

### Extensions

```sql
-- List available extensions
SELECT * FROM pg_available_extensions;

-- List installed extensions
\dx
SELECT * FROM pg_extension;

-- Install extension
CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS uuid-ossp;
CREATE EXTENSION IF NOT EXISTS hstore;
CREATE EXTENSION IF NOT EXISTS PostGIS;

-- Popular extensions:
-- pg_stat_statements - Query performance statistics
-- pgcrypto - Cryptographic functions
-- uuid-ossp - UUID generation
-- hstore - Key-value storage
-- PostGIS - Geospatial support
-- pg_trgm - Trigram similarity
-- btree_gin - GIN support for more types
-- tablefunc - Crosstab functions
-- pg_cron - Job scheduling

-- Drop extension
DROP EXTENSION IF EXISTS pg_trgm;

-- Example: UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
SELECT uuid_generate_v4();

-- Example: Encryption
CREATE EXTENSION IF NOT EXISTS pgcrypto;
SELECT crypt('password', gen_salt('bf'));
SELECT crypt('password', stored_hash) = stored_hash AS password_match;
```

### Foreign Data Wrappers (FDW)

```sql
-- Connect to another PostgreSQL database
CREATE EXTENSION postgres_fdw;

CREATE SERVER remote_server
    FOREIGN DATA WRAPPER postgres_fdw
    OPTIONS (host 'remote.example.com', port '5432', dbname 'remote_db');

CREATE USER MAPPING FOR local_user
    SERVER remote_server
    OPTIONS (user 'remote_user', password 'remote_password');

-- Import foreign tables
IMPORT FOREIGN SCHEMA public
    FROM SERVER remote_server
    INTO local_schema;

-- Or create individual foreign table
CREATE FOREIGN TABLE remote_users (
    id INT,
    name VARCHAR(100),
    email VARCHAR(255)
)
SERVER remote_server
OPTIONS (schema_name 'public', table_name 'users');

-- Query foreign table like local table
SELECT * FROM remote_users WHERE id = 1;
```

### Listen/Notify

```sql
-- Session 1: Listen for events
LISTEN order_created;

-- Session 2: Send notification
NOTIFY order_created, '{"order_id": 12345, "total": 99.99}';

-- In application code, handle notifications asynchronously
-- Useful for real-time features without polling
```

### Table Inheritance

```sql
-- Parent table
CREATE TABLE cities (
    name VARCHAR(100),
    population INT
);

-- Child table inherits columns
CREATE TABLE capitals (
    country VARCHAR(100)
) INHERITS (cities);

-- Query parent includes children
SELECT * FROM cities;  -- Returns cities and capitals

-- Query only parent
SELECT * FROM ONLY cities;  -- Returns only cities

-- Check inheritance
SELECT * FROM capitals;  -- Has name, population, country
```

---

## 20. Best Practices

### Schema Design

```sql
-- 1. Always have a primary key
CREATE TABLE orders (
    id BIGSERIAL PRIMARY KEY,
    ...
);

-- 2. Use appropriate data types
-- Bad: storing money as FLOAT
price FLOAT  -- NO!
-- Good:
price NUMERIC(10, 2)

-- 3. Use UUID for distributed systems
id UUID PRIMARY KEY DEFAULT gen_random_uuid()

-- 4. Add timestamps
created_at TIMESTAMPTZ DEFAULT NOW(),
updated_at TIMESTAMPTZ DEFAULT NOW()

-- 5. Use TIMESTAMPTZ for timestamps
created_at TIMESTAMP WITH TIME ZONE  -- Good
created_at TIMESTAMP               -- Avoid (no timezone info)

-- 6. Use schemas for organization
CREATE SCHEMA app;
CREATE SCHEMA audit;
CREATE SCHEMA archive;

-- 7. Use meaningful names
-- Tables: plural, snake_case
-- Columns: singular, snake_case
-- Indexes: idx_tablename_columns
-- Constraints: pk_, fk_, uq_, ck_

-- 8. Document with comments
COMMENT ON TABLE orders IS 'Customer orders';
COMMENT ON COLUMN orders.status IS 'Order status: pending, processing, shipped, delivered';
```

### Query Best Practices

```sql
-- 1. Select only needed columns
SELECT id, name, email FROM users;  -- Not SELECT *

-- 2. Use parameterized queries (in application)
-- Prevents SQL injection

-- 3. Use EXPLAIN ANALYZE
EXPLAIN ANALYZE SELECT * FROM orders WHERE customer_id = 100;

-- 4. Add appropriate indexes
CREATE INDEX idx_orders_customer ON orders(customer_id);

-- 5. Use CTEs for readability
WITH active_users AS (
    SELECT * FROM users WHERE is_active = TRUE
)
SELECT * FROM active_users WHERE created_at > '2024-01-01';

-- 6. Use transactions appropriately
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;

-- 7. Use LIMIT for large queries
SELECT * FROM logs ORDER BY created_at DESC LIMIT 100;

-- 8. Avoid N+1 queries
-- Bad: Query users, then loop for each user's orders
-- Good: JOIN or use subqueries
```

### Security

```sql
-- 1. Use strong passwords
CREATE ROLE app WITH PASSWORD 'C0mpl3x!P@ssword#2024';

-- 2. Grant minimum privileges
GRANT SELECT, INSERT, UPDATE ON app_tables TO app_user;

-- 3. Use SSL connections
-- In pg_hba.conf: hostssl all all 0.0.0.0/0 scram-sha-256

-- 4. Use Row-Level Security when needed
ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

-- 5. Regular security audits
SELECT rolname, rolsuper FROM pg_roles WHERE rolsuper;

-- 6. Keep PostgreSQL updated
```

### Maintenance

```sql
-- 1. Regular VACUUM and ANALYZE
VACUUM ANALYZE;

-- 2. Monitor bloat
SELECT * FROM pg_stat_user_tables 
WHERE n_dead_tup > n_live_tup * 0.2;

-- 3. Monitor slow queries
SELECT query, calls, mean_exec_time 
FROM pg_stat_statements 
ORDER BY mean_exec_time DESC;

-- 4. Monitor connections
SELECT count(*) FROM pg_stat_activity;

-- 5. Regular backups
-- Schedule pg_dump or use continuous archiving

-- 6. Test restore procedures
-- Regularly verify backups work

-- 7. Monitor replication lag
SELECT * FROM pg_stat_replication;
```

---

## Quick Reference

### Common Commands

```sql
-- Database operations
CREATE DATABASE dbname;
DROP DATABASE dbname;
\c dbname

-- Table operations
CREATE TABLE tablename (...);
DROP TABLE tablename;
ALTER TABLE tablename ADD COLUMN col type;
\d tablename

-- User operations
CREATE USER username WITH PASSWORD 'password';
GRANT privileges ON database/table TO username;
DROP USER username;

-- Data operations
INSERT INTO table (...) VALUES (...);
SELECT columns FROM table WHERE condition;
UPDATE table SET column = value WHERE condition;
DELETE FROM table WHERE condition;

-- Transactions
BEGIN;
COMMIT;
ROLLBACK;
SAVEPOINT name;
```

### psql Commands

```
\l          List databases
\c dbname   Connect to database
\dt         List tables
\d table    Describe table
\di         List indexes
\dv         List views
\df         List functions
\du         List users
\dx         List extensions
\timing     Toggle timing
\x          Toggle expanded output
\q          Quit
```

---

## Resources

### Official Documentation
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [PostgreSQL Wiki](https://wiki.postgresql.org/)

### Learning Resources
- PostgreSQL Tutorial (postgresqltutorial.com)
- PostgreSQL Exercises
- Awesome PostgreSQL (GitHub)

### Tools
- **pgAdmin** - Official GUI administration tool
- **DBeaver** - Universal database tool
- **DataGrip** - JetBrains IDE
- **psql** - Command-line client
- **pgcli** - Enhanced command-line

### Extensions
- **PostGIS** - Geospatial support
- **TimescaleDB** - Time-series data
- **Citus** - Distributed PostgreSQL
- **pgvector** - Vector similarity search

---

*Last Updated: February 2026*
