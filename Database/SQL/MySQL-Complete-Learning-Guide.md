# MySQL Complete Learning Guide
## From Fundamentals to Advanced Concepts

---

## Table of Contents

1. [Introduction to MySQL](#1-introduction-to-mysql)
2. [Installation & Setup](#2-installation--setup)
3. [Database Fundamentals](#3-database-fundamentals)
4. [Data Types](#4-data-types)
5. [CRUD Operations](#5-crud-operations)
6. [Querying Data](#6-querying-data)
7. [Joins & Relationships](#7-joins--relationships)
8. [Subqueries & CTEs](#8-subqueries--ctes)
9. [Aggregation & Grouping](#9-aggregation--grouping)
10. [Indexes & Performance](#10-indexes--performance)
11. [Views](#11-views)
12. [Stored Procedures & Functions](#12-stored-procedures--functions)
13. [Triggers](#13-triggers)
14. [Transactions](#14-transactions)
15. [User Management & Security](#15-user-management--security)
16. [Backup & Recovery](#16-backup--recovery)
17. [Replication](#17-replication)
18. [Performance Tuning](#18-performance-tuning)
19. [MySQL 8.0+ Features](#19-mysql-80-features)
20. [Best Practices](#20-best-practices)

---

## 1. Introduction to MySQL

### What is MySQL?

MySQL is an open-source relational database management system (RDBMS) that uses Structured Query Language (SQL). It's one of the most popular databases worldwide, powering many web applications, including WordPress, Facebook, and Twitter.

### Key Features

| Feature | Description |
|---------|-------------|
| **Open Source** | Free to use under GPL license |
| **Cross-Platform** | Runs on Windows, Linux, macOS |
| **High Performance** | Optimized for read-heavy workloads |
| **Scalability** | Handles large databases efficiently |
| **Security** | Strong data protection and access control |
| **Replication** | Built-in master-slave replication |
| **ACID Compliant** | Full transaction support with InnoDB |

### MySQL vs Other Databases

| Aspect | MySQL | PostgreSQL | SQL Server |
|--------|-------|------------|------------|
| License | Open Source (GPL) | Open Source | Commercial |
| Best For | Web applications | Complex queries | Enterprise/.NET |
| JSON Support | Good (5.7+) | Excellent (JSONB) | Good |
| Full-Text Search | Built-in | Built-in | Built-in |
| Replication | Easy setup | More complex | Enterprise feature |
| Learning Curve | Easy | Moderate | Moderate |

### MySQL Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT LAYER                            │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ MySQL CLI│  │   JDBC   │  │   PHP    │  │  Python  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     SERVER LAYER                            │
│  ┌─────────────────────────────────────────────────────┐   │
│  │              Connection Pool / Thread Pool           │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                    Query Parser                      │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │                   Query Optimizer                    │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │               Query Execution Engine                 │   │
│  └─────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    STORAGE ENGINE LAYER                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  InnoDB  │  │  MyISAM  │  │  Memory  │  │  Archive │    │
│  │(Default) │  │ (Legacy) │  │  (RAM)   │  │(Compress)│    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### Storage Engines Comparison

| Feature | InnoDB | MyISAM | Memory |
|---------|--------|--------|--------|
| **Transactions** | ✅ Yes | ❌ No | ❌ No |
| **Foreign Keys** | ✅ Yes | ❌ No | ❌ No |
| **Row Locking** | ✅ Yes | ❌ Table only | ✅ Yes |
| **Crash Recovery** | ✅ Yes | ❌ No | ❌ No |
| **Full-Text Index** | ✅ Yes (5.6+) | ✅ Yes | ❌ No |
| **Data Caching** | ✅ Yes | ❌ No | N/A |
| **Best For** | General use | Read-only, legacy | Temporary data |

**Always use InnoDB** unless you have a specific reason not to.

---

## 2. Installation & Setup

### Installing MySQL on Windows

```powershell
# Using Chocolatey
choco install mysql

# Or download MySQL Installer from:
# https://dev.mysql.com/downloads/installer/

# After installation, start MySQL service
net start mysql

# Connect to MySQL
mysql -u root -p
```

### Installing MySQL on Linux (Ubuntu/Debian)

```bash
# Update package list
sudo apt update

# Install MySQL Server
sudo apt install mysql-server

# Start MySQL service
sudo systemctl start mysql
sudo systemctl enable mysql

# Secure installation (set root password, remove test DB, etc.)
sudo mysql_secure_installation

# Connect to MySQL
sudo mysql -u root -p
```

### Installing MySQL on macOS

```bash
# Using Homebrew
brew install mysql

# Start MySQL service
brew services start mysql

# Secure installation
mysql_secure_installation

# Connect
mysql -u root -p
```

### Docker Installation (Recommended for Development)

```yaml
# docker-compose.yml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    container_name: mysql_dev
    environment:
      MYSQL_ROOT_PASSWORD: rootpassword
      MYSQL_DATABASE: myapp
      MYSQL_USER: developer
      MYSQL_PASSWORD: devpassword
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql
    command: --default-authentication-plugin=mysql_native_password

volumes:
  mysql_data:
```

```bash
# Start MySQL container
docker-compose up -d

# Connect to MySQL in container
docker exec -it mysql_dev mysql -u root -p
```

### Initial Configuration

```sql
-- Connect as root
mysql -u root -p

-- Check MySQL version
SELECT VERSION();

-- Show all databases
SHOW DATABASES;

-- Create a new database
CREATE DATABASE myapp;

-- Create a new user
CREATE USER 'developer'@'localhost' IDENTIFIED BY 'password123';

-- Grant privileges
GRANT ALL PRIVILEGES ON myapp.* TO 'developer'@'localhost';
FLUSH PRIVILEGES;

-- Use the new database
USE myapp;
```

### MySQL Configuration File (my.cnf / my.ini)

```ini
[mysqld]
# Basic Settings
port = 3306
datadir = /var/lib/mysql
socket = /var/run/mysqld/mysqld.sock

# Character Set
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci

# InnoDB Settings
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 1
innodb_file_per_table = 1

# Query Cache (deprecated in 8.0)
# query_cache_type = 0

# Connection Settings
max_connections = 200
wait_timeout = 28800

# Logging
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2

[client]
default-character-set = utf8mb4
```

---

## 3. Database Fundamentals

### Creating Databases

```sql
-- Create a simple database
CREATE DATABASE myapp;

-- Create with specific character set
CREATE DATABASE myapp
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

-- Create if not exists
CREATE DATABASE IF NOT EXISTS myapp;

-- Show create statement
SHOW CREATE DATABASE myapp;
```

### Selecting and Dropping Databases

```sql
-- Select a database to use
USE myapp;

-- Show current database
SELECT DATABASE();

-- List all databases
SHOW DATABASES;

-- Drop a database (CAREFUL!)
DROP DATABASE myapp;

-- Drop if exists
DROP DATABASE IF EXISTS myapp;
```

### Creating Tables

```sql
-- Basic table creation
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Table with multiple constraints
CREATE TABLE products (
    product_id INT PRIMARY KEY AUTO_INCREMENT,
    sku VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    stock_quantity INT DEFAULT 0 CHECK (stock_quantity >= 0),
    category_id INT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_category (category_id),
    INDEX idx_name (name),
    FULLTEXT INDEX idx_search (name, description)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Table with foreign keys
CREATE TABLE orders (
    order_id INT PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    order_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    total_amount DECIMAL(12, 2) NOT NULL,
    status ENUM('pending', 'processing', 'shipped', 'delivered', 'cancelled') DEFAULT 'pending',
    
    CONSTRAINT fk_orders_user 
        FOREIGN KEY (user_id) 
        REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
) ENGINE=InnoDB;

-- Junction table for many-to-many
CREATE TABLE order_items (
    order_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit_price DECIMAL(10, 2) NOT NULL,
    
    PRIMARY KEY (order_id, product_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(product_id) ON DELETE RESTRICT
) ENGINE=InnoDB;
```

### Altering Tables

```sql
-- Add a column
ALTER TABLE users ADD COLUMN phone VARCHAR(20) AFTER email;

-- Add column with default
ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT FALSE;

-- Modify column type
ALTER TABLE users MODIFY COLUMN phone VARCHAR(30);

-- Rename column
ALTER TABLE users CHANGE COLUMN phone phone_number VARCHAR(30);

-- Drop column
ALTER TABLE users DROP COLUMN phone_number;

-- Add index
ALTER TABLE users ADD INDEX idx_email (email);

-- Add unique constraint
ALTER TABLE users ADD UNIQUE INDEX uq_username (username);

-- Add foreign key
ALTER TABLE orders 
    ADD CONSTRAINT fk_orders_user 
    FOREIGN KEY (user_id) REFERENCES users(id);

-- Drop foreign key
ALTER TABLE orders DROP FOREIGN KEY fk_orders_user;

-- Drop index
ALTER TABLE users DROP INDEX idx_email;

-- Rename table
ALTER TABLE users RENAME TO app_users;

-- Or
RENAME TABLE users TO app_users;

-- Change engine
ALTER TABLE users ENGINE = InnoDB;
```

### Viewing Table Structure

```sql
-- Describe table structure
DESCRIBE users;
-- Or
DESC users;
-- Or
SHOW COLUMNS FROM users;

-- Show full column information
SHOW FULL COLUMNS FROM users;

-- Show create table statement
SHOW CREATE TABLE users;

-- Show indexes
SHOW INDEX FROM users;

-- Show table status
SHOW TABLE STATUS LIKE 'users';

-- List all tables
SHOW TABLES;

-- Show tables with pattern
SHOW TABLES LIKE 'user%';
```

---

## 4. Data Types

### Numeric Types

```sql
-- Integer Types
CREATE TABLE numeric_examples (
    -- Exact integers
    tiny_col TINYINT,           -- -128 to 127 (1 byte)
    tiny_unsigned TINYINT UNSIGNED,  -- 0 to 255
    small_col SMALLINT,         -- -32,768 to 32,767 (2 bytes)
    medium_col MEDIUMINT,       -- -8,388,608 to 8,388,607 (3 bytes)
    int_col INT,                -- -2.1B to 2.1B (4 bytes)
    big_col BIGINT,             -- -9.2Q to 9.2Q (8 bytes)
    
    -- Decimal/Fixed-point (exact)
    price DECIMAL(10, 2),       -- 10 digits, 2 decimal places
    salary NUMERIC(12, 2),      -- Same as DECIMAL
    
    -- Floating-point (approximate)
    float_col FLOAT,            -- 4 bytes, ~7 decimal digits
    double_col DOUBLE,          -- 8 bytes, ~15 decimal digits
    
    -- Boolean (alias for TINYINT(1))
    is_active BOOLEAN           -- TRUE/FALSE or 1/0
);

-- Best Practices:
-- Use INT for IDs, counts
-- Use BIGINT for large numbers, timestamps in milliseconds
-- Use DECIMAL for money (never FLOAT!)
-- Use BOOLEAN for flags
```

### String Types

```sql
CREATE TABLE string_examples (
    -- Fixed-length string (padded with spaces)
    country_code CHAR(2),       -- Always 2 bytes
    
    -- Variable-length string
    username VARCHAR(50),       -- Max 50 chars, uses only needed space
    email VARCHAR(255),         -- Common for emails
    
    -- Text types (for large content)
    short_text TINYTEXT,        -- Max 255 bytes
    description TEXT,           -- Max 65,535 bytes (~64KB)
    content MEDIUMTEXT,         -- Max 16,777,215 bytes (~16MB)
    full_text LONGTEXT,         -- Max 4,294,967,295 bytes (~4GB)
    
    -- Binary types
    file_hash BINARY(32),       -- Fixed-length binary
    file_data VARBINARY(1000),  -- Variable-length binary
    image BLOB,                 -- Binary large object
    document LONGBLOB,          -- Large binary
    
    -- Enum (one value from list)
    status ENUM('active', 'inactive', 'pending'),
    
    -- Set (multiple values from list)
    permissions SET('read', 'write', 'delete', 'admin')
);

-- Examples
INSERT INTO string_examples (status, permissions) 
VALUES ('active', 'read,write');
```

### Date and Time Types

```sql
CREATE TABLE datetime_examples (
    -- Date only (YYYY-MM-DD)
    birth_date DATE,            -- '2024-03-15'
    
    -- Time only (HH:MM:SS)
    start_time TIME,            -- '14:30:00'
    
    -- Date and time
    created_at DATETIME,        -- '2024-03-15 14:30:00'
    
    -- Timestamp (auto-converts to UTC, range 1970-2038)
    updated_at TIMESTAMP,       -- Auto-updates possible
    
    -- Year only
    graduation_year YEAR        -- 2024
);

-- Timestamp with auto-update
CREATE TABLE articles (
    id INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- Date/Time Functions
SELECT 
    NOW(),                      -- Current date and time
    CURDATE(),                  -- Current date
    CURTIME(),                  -- Current time
    DATE('2024-03-15 14:30:00'), -- Extract date
    TIME('2024-03-15 14:30:00'), -- Extract time
    YEAR(NOW()),                -- Extract year
    MONTH(NOW()),               -- Extract month
    DAY(NOW()),                 -- Extract day
    HOUR(NOW()),                -- Extract hour
    DATEDIFF('2024-03-20', '2024-03-15'),  -- Days difference
    DATE_ADD(NOW(), INTERVAL 7 DAY),       -- Add days
    DATE_SUB(NOW(), INTERVAL 1 MONTH)      -- Subtract month
;
```

### JSON Type (MySQL 5.7+)

```sql
CREATE TABLE products (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200),
    attributes JSON,            -- Flexible schema
    metadata JSON
);

-- Insert JSON data
INSERT INTO products (name, attributes) VALUES 
    ('Laptop', '{"brand": "Dell", "ram": 16, "storage": "512GB SSD"}'),
    ('Phone', JSON_OBJECT('brand', 'Apple', 'model', 'iPhone 15', 'storage', 256));

-- Query JSON data
SELECT 
    name,
    attributes->>'$.brand' AS brand,           -- Extract as text
    attributes->'$.ram' AS ram,                -- Extract as JSON
    JSON_EXTRACT(attributes, '$.storage') AS storage
FROM products;

-- Filter by JSON value
SELECT * FROM products 
WHERE attributes->>'$.brand' = 'Dell';

-- JSON functions
SELECT 
    JSON_OBJECT('key', 'value'),               -- Create object
    JSON_ARRAY(1, 2, 3),                       -- Create array
    JSON_KEYS(attributes),                      -- Get keys
    JSON_LENGTH(attributes),                    -- Count elements
    JSON_CONTAINS(attributes, '"Dell"', '$.brand'),  -- Check contains
    JSON_SEARCH(attributes, 'one', 'Dell')     -- Search value
FROM products WHERE id = 1;

-- Update JSON
UPDATE products 
SET attributes = JSON_SET(attributes, '$.ram', 32)
WHERE id = 1;

-- Add to JSON
UPDATE products 
SET attributes = JSON_INSERT(attributes, '$.color', 'Silver')
WHERE id = 1;

-- Remove from JSON
UPDATE products 
SET attributes = JSON_REMOVE(attributes, '$.color')
WHERE id = 1;
```

### Spatial Types (GIS)

```sql
CREATE TABLE locations (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100),
    coordinates POINT,          -- Single point
    area POLYGON,              -- Closed shape
    route LINESTRING           -- Line/path
);

-- Insert spatial data
INSERT INTO locations (name, coordinates) VALUES 
    ('Office', ST_GeomFromText('POINT(40.7128 -74.0060)'));

-- Query with spatial functions
SELECT 
    name,
    ST_X(coordinates) AS latitude,
    ST_Y(coordinates) AS longitude,
    ST_Distance_Sphere(
        coordinates,
        ST_GeomFromText('POINT(40.7580 -73.9855)')
    ) AS distance_meters
FROM locations;
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

-- Insert with all columns (specify all values in order)
INSERT INTO users 
VALUES (NULL, 'mike', 'mike@example.com', 'hash4', NOW(), NOW());

-- Insert from SELECT
INSERT INTO users_backup (username, email, password_hash)
SELECT username, email, password_hash FROM users WHERE is_active = TRUE;

-- Insert with ON DUPLICATE KEY UPDATE (upsert)
INSERT INTO users (username, email, password_hash)
VALUES ('john_doe', 'john_new@example.com', 'new_hash')
ON DUPLICATE KEY UPDATE 
    email = VALUES(email),
    password_hash = VALUES(password_hash);

-- Insert IGNORE (skip duplicates silently)
INSERT IGNORE INTO users (username, email, password_hash)
VALUES ('john_doe', 'duplicate@example.com', 'hash');

-- REPLACE (delete then insert if exists)
REPLACE INTO users (id, username, email, password_hash)
VALUES (1, 'john_doe', 'john_replaced@example.com', 'new_hash');

-- Get last inserted ID
SELECT LAST_INSERT_ID();
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
    CONCAT(first_name, ' ', last_name) AS full_name
FROM users;

-- Select distinct values
SELECT DISTINCT status FROM orders;

-- Select with WHERE
SELECT * FROM users WHERE is_active = TRUE;

-- Select with multiple conditions
SELECT * FROM products
WHERE price > 100 
  AND category_id = 5
  AND (stock_quantity > 0 OR is_preorder = TRUE);

-- Select with IN
SELECT * FROM users WHERE id IN (1, 2, 3, 4, 5);

-- Select with BETWEEN
SELECT * FROM orders 
WHERE order_date BETWEEN '2024-01-01' AND '2024-12-31';

-- Select with LIKE (pattern matching)
SELECT * FROM users WHERE email LIKE '%@gmail.com';
SELECT * FROM products WHERE name LIKE 'iPhone%';
SELECT * FROM users WHERE username LIKE '_ohn%';  -- Second char is 'o'

-- Select with NULL checks
SELECT * FROM users WHERE phone IS NULL;
SELECT * FROM users WHERE phone IS NOT NULL;

-- Select with ORDER BY
SELECT * FROM products ORDER BY price ASC;
SELECT * FROM products ORDER BY created_at DESC, name ASC;

-- Select with LIMIT and OFFSET
SELECT * FROM products ORDER BY id LIMIT 10;          -- First 10
SELECT * FROM products ORDER BY id LIMIT 10 OFFSET 20; -- Skip 20, get 10
SELECT * FROM products ORDER BY id LIMIT 20, 10;      -- Same as above

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

-- Update with JOIN
UPDATE orders o
INNER JOIN users u ON o.user_id = u.id
SET o.status = 'cancelled'
WHERE u.is_active = FALSE;

-- Update with subquery
UPDATE products 
SET category_id = (SELECT id FROM categories WHERE name = 'Electronics')
WHERE name LIKE '%Phone%';

-- Update with LIMIT (useful for batch updates)
UPDATE users SET is_notified = TRUE 
WHERE is_notified = FALSE 
ORDER BY created_at 
LIMIT 1000;

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

-- Delete with JOIN
DELETE o FROM orders o
INNER JOIN users u ON o.user_id = u.id
WHERE u.is_deleted = TRUE;

-- Delete with subquery
DELETE FROM order_items 
WHERE order_id IN (
    SELECT order_id FROM orders WHERE status = 'cancelled'
);

-- Delete with LIMIT (batch delete)
DELETE FROM logs WHERE created_at < '2023-01-01' LIMIT 10000;

-- Delete all rows (keeps table structure)
DELETE FROM temp_data;

-- TRUNCATE - faster way to delete all rows
TRUNCATE TABLE temp_data;
-- Note: TRUNCATE resets AUTO_INCREMENT, cannot be rolled back, ignores triggers
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
    
    -- IN (list of values)
    AND category_id IN (1, 2, 3)
    
    -- NOT IN
    AND status NOT IN ('deleted', 'archived')
    
    -- NULL checks
    AND deleted_at IS NULL
    AND description IS NOT NULL
    
    -- LIKE patterns
    AND name LIKE 'iPhone%'       -- Starts with
    AND email LIKE '%@gmail.com'  -- Ends with
    AND code LIKE 'A_123'         -- _ = single character
    
    -- REGEXP (regular expressions)
    AND email REGEXP '^[a-z]+@'
;
```

### Logical Operators

```sql
-- AND - all conditions must be true
SELECT * FROM products 
WHERE price > 100 AND stock > 0 AND is_active = TRUE;

-- OR - at least one condition must be true
SELECT * FROM products 
WHERE category_id = 1 OR category_id = 2;

-- NOT - negates condition
SELECT * FROM products WHERE NOT is_deleted;

-- Complex conditions with parentheses
SELECT * FROM products
WHERE (category_id = 1 OR category_id = 2)
  AND price < 500
  AND (stock > 0 OR is_preorder = TRUE);

-- XOR - exclusive or (one or the other, not both)
SELECT * FROM products WHERE is_featured XOR is_sale;
```

### String Functions

```sql
SELECT 
    -- Length
    LENGTH('Hello'),              -- 5 (bytes)
    CHAR_LENGTH('Hello'),         -- 5 (characters)
    
    -- Case conversion
    UPPER('hello'),               -- 'HELLO'
    LOWER('HELLO'),               -- 'hello'
    
    -- Concatenation
    CONCAT('Hello', ' ', 'World'), -- 'Hello World'
    CONCAT_WS(', ', 'a', 'b', 'c'), -- 'a, b, c'
    
    -- Substring
    SUBSTRING('Hello World', 1, 5),  -- 'Hello'
    SUBSTRING('Hello World', 7),     -- 'World'
    LEFT('Hello', 2),               -- 'He'
    RIGHT('Hello', 2),              -- 'lo'
    
    -- Trim
    TRIM('  hello  '),             -- 'hello'
    LTRIM('  hello'),              -- 'hello'
    RTRIM('hello  '),              -- 'hello'
    
    -- Replace
    REPLACE('Hello World', 'World', 'MySQL'), -- 'Hello MySQL'
    
    -- Position
    LOCATE('World', 'Hello World'), -- 7
    INSTR('Hello World', 'World'),  -- 7
    
    -- Padding
    LPAD('42', 5, '0'),            -- '00042'
    RPAD('Hi', 5, '*'),            -- 'Hi***'
    
    -- Reverse
    REVERSE('Hello'),              -- 'olleH'
    
    -- Format
    FORMAT(1234567.891, 2)         -- '1,234,567.89'
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
    TRUNCATE(3.14159, 2),  -- 3.14
    
    -- Absolute and Sign
    ABS(-5),               -- 5
    SIGN(-5),              -- -1
    SIGN(5),               -- 1
    
    -- Power and Root
    POWER(2, 8),           -- 256
    SQRT(16),              -- 4
    
    -- Modulo
    MOD(17, 5),            -- 2
    17 % 5,                -- 2
    
    -- Random
    RAND(),                -- Random 0-1
    FLOOR(RAND() * 100),   -- Random 0-99
    
    -- Min/Max
    GREATEST(1, 5, 3),     -- 5
    LEAST(1, 5, 3),        -- 1
    
    -- Logarithm
    LOG(10),               -- Natural log
    LOG10(100),            -- 2
    LOG2(8)                -- 3
;
```

### Date Functions

```sql
SELECT 
    -- Current date/time
    NOW(),                          -- 2024-03-15 14:30:00
    CURRENT_TIMESTAMP(),            -- Same as NOW()
    CURDATE(),                      -- 2024-03-15
    CURRENT_DATE(),                 -- Same as CURDATE()
    CURTIME(),                      -- 14:30:00
    CURRENT_TIME(),                 -- Same as CURTIME()
    
    -- Extract parts
    YEAR('2024-03-15'),             -- 2024
    MONTH('2024-03-15'),            -- 3
    DAY('2024-03-15'),              -- 15
    HOUR('14:30:45'),               -- 14
    MINUTE('14:30:45'),             -- 30
    SECOND('14:30:45'),             -- 45
    DAYOFWEEK('2024-03-15'),        -- 6 (Friday, 1=Sunday)
    DAYOFYEAR('2024-03-15'),        -- 75
    WEEK('2024-03-15'),             -- 11
    QUARTER('2024-03-15'),          -- 1
    
    -- Date arithmetic
    DATE_ADD('2024-03-15', INTERVAL 7 DAY),      -- 2024-03-22
    DATE_ADD('2024-03-15', INTERVAL 1 MONTH),    -- 2024-04-15
    DATE_SUB('2024-03-15', INTERVAL 1 YEAR),     -- 2023-03-15
    '2024-03-15' + INTERVAL 7 DAY,               -- 2024-03-22
    
    -- Difference
    DATEDIFF('2024-03-20', '2024-03-15'),        -- 5
    TIMESTAMPDIFF(HOUR, '2024-03-15 10:00', '2024-03-15 14:30'),  -- 4
    
    -- Formatting
    DATE_FORMAT('2024-03-15', '%M %d, %Y'),      -- 'March 15, 2024'
    DATE_FORMAT('2024-03-15', '%W'),             -- 'Friday'
    DATE_FORMAT('2024-03-15 14:30:00', '%h:%i %p'), -- '02:30 PM'
    
    -- Parsing
    STR_TO_DATE('15-03-2024', '%d-%m-%Y'),       -- 2024-03-15
    
    -- First/Last day of month
    LAST_DAY('2024-03-15'),         -- 2024-03-31
    DATE_FORMAT('2024-03-15', '%Y-%m-01')        -- 2024-03-01
;

-- Common date format specifiers:
-- %Y - 4-digit year (2024)
-- %y - 2-digit year (24)
-- %M - Month name (March)
-- %m - Month number (03)
-- %D - Day with suffix (15th)
-- %d - Day number (15)
-- %W - Weekday name (Friday)
-- %H - Hour 24h (14)
-- %h - Hour 12h (02)
-- %i - Minutes (30)
-- %s - Seconds (45)
-- %p - AM/PM
```

### Conditional Functions

```sql
-- IF function
SELECT 
    name,
    IF(stock > 0, 'In Stock', 'Out of Stock') AS availability
FROM products;

-- IFNULL - replace NULL with value
SELECT IFNULL(phone, 'No phone') AS phone FROM users;

-- NULLIF - return NULL if values are equal
SELECT NULLIF(discount, 0) AS discount FROM orders;  -- NULL if discount is 0

-- COALESCE - return first non-NULL value
SELECT COALESCE(nickname, username, email) AS display_name FROM users;

-- CASE expression
SELECT 
    name,
    price,
    CASE 
        WHEN price < 50 THEN 'Budget'
        WHEN price < 200 THEN 'Mid-range'
        WHEN price < 1000 THEN 'Premium'
        ELSE 'Luxury'
    END AS tier
FROM products;

-- Simple CASE
SELECT 
    order_id,
    CASE status
        WHEN 'pending' THEN 'Awaiting Payment'
        WHEN 'processing' THEN 'Being Prepared'
        WHEN 'shipped' THEN 'On the Way'
        WHEN 'delivered' THEN 'Completed'
        ELSE 'Unknown'
    END AS status_text
FROM orders;
```

---

## 7. Joins & Relationships

### Types of Joins

```sql
-- Sample tables
CREATE TABLE customers (
    id INT PRIMARY KEY,
    name VARCHAR(100)
);

CREATE TABLE orders (
    id INT PRIMARY KEY,
    customer_id INT,
    total DECIMAL(10, 2)
);

INSERT INTO customers VALUES (1, 'Alice'), (2, 'Bob'), (3, 'Charlie');
INSERT INTO orders VALUES (101, 1, 100.00), (102, 1, 200.00), (103, 2, 150.00), (104, NULL, 50.00);
```

### INNER JOIN

```sql
-- Returns only matching rows from both tables
SELECT 
    c.name AS customer_name,
    o.id AS order_id,
    o.total
FROM customers c
INNER JOIN orders o ON c.id = o.customer_id;

-- Result:
-- | customer_name | order_id | total  |
-- |---------------|----------|--------|
-- | Alice         | 101      | 100.00 |
-- | Alice         | 102      | 200.00 |
-- | Bob           | 103      | 150.00 |
```

### LEFT JOIN (LEFT OUTER JOIN)

```sql
-- Returns all rows from left table, matched rows from right (or NULL)
SELECT 
    c.name AS customer_name,
    o.id AS order_id,
    o.total
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id;

-- Result:
-- | customer_name | order_id | total  |
-- |---------------|----------|--------|
-- | Alice         | 101      | 100.00 |
-- | Alice         | 102      | 200.00 |
-- | Bob           | 103      | 150.00 |
-- | Charlie       | NULL     | NULL   |

-- Find customers with no orders
SELECT c.name
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id
WHERE o.id IS NULL;
```

### RIGHT JOIN (RIGHT OUTER JOIN)

```sql
-- Returns all rows from right table, matched rows from left (or NULL)
SELECT 
    c.name AS customer_name,
    o.id AS order_id,
    o.total
FROM customers c
RIGHT JOIN orders o ON c.id = o.customer_id;

-- Result:
-- | customer_name | order_id | total  |
-- |---------------|----------|--------|
-- | Alice         | 101      | 100.00 |
-- | Alice         | 102      | 200.00 |
-- | Bob           | 103      | 150.00 |
-- | NULL          | 104      | 50.00  |
```

### CROSS JOIN (Cartesian Product)

```sql
-- Returns all combinations of rows
SELECT 
    c.name,
    p.name AS product
FROM customers c
CROSS JOIN products p;

-- Every customer paired with every product
-- Rows = customers × products
```

### Self Join

```sql
-- Table referencing itself
CREATE TABLE employees (
    id INT PRIMARY KEY,
    name VARCHAR(100),
    manager_id INT
);

INSERT INTO employees VALUES 
    (1, 'CEO', NULL),
    (2, 'Manager A', 1),
    (3, 'Manager B', 1),
    (4, 'Employee 1', 2),
    (5, 'Employee 2', 2);

-- Get employee with their manager's name
SELECT 
    e.name AS employee,
    m.name AS manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

### Multiple Joins

```sql
-- Joining multiple tables
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

-- Mix of join types
SELECT 
    c.name AS customer_name,
    COUNT(o.id) AS order_count,
    COALESCE(SUM(o.total), 0) AS total_spent
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id AND o.status != 'cancelled'
GROUP BY c.id, c.name;
```

### Join with Conditions

```sql
-- Join condition in ON vs WHERE
-- ON: Affects the join itself
-- WHERE: Filters results after join

-- These are different for LEFT JOIN:
SELECT *
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id AND o.total > 100;
-- Returns all customers, but only orders > 100 are joined

SELECT *
FROM customers c
LEFT JOIN orders o ON c.id = o.customer_id
WHERE o.total > 100;
-- Filters out customers with no orders > 100

-- For INNER JOIN, both produce same results
```

---

## 8. Subqueries & CTEs

### Subqueries in WHERE

```sql
-- Scalar subquery (returns single value)
SELECT * FROM products 
WHERE price > (SELECT AVG(price) FROM products);

-- IN subquery (returns list)
SELECT * FROM customers 
WHERE id IN (
    SELECT DISTINCT customer_id FROM orders WHERE total > 1000
);

-- NOT IN subquery
SELECT * FROM products 
WHERE id NOT IN (
    SELECT DISTINCT product_id FROM order_items
);

-- EXISTS subquery (checks for existence)
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

-- Comparison with subquery
SELECT * FROM orders
WHERE total >= ALL (SELECT total FROM orders WHERE customer_id = 1);

SELECT * FROM products
WHERE price > ANY (SELECT price FROM products WHERE category_id = 5);
```

### Subqueries in FROM (Derived Tables)

```sql
-- Subquery as a table
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
    GROUP BY c.id, c.name
) AS category_stats
WHERE product_count > 5
ORDER BY avg_price DESC;

-- Join with derived table
SELECT 
    c.name,
    customer_totals.total_spent
FROM customers c
JOIN (
    SELECT customer_id, SUM(total) AS total_spent
    FROM orders
    WHERE status = 'completed'
    GROUP BY customer_id
) AS customer_totals ON c.id = customer_totals.customer_id
WHERE customer_totals.total_spent > 1000;
```

### Subqueries in SELECT

```sql
-- Scalar subquery in SELECT
SELECT 
    name,
    price,
    price - (SELECT AVG(price) FROM products) AS price_vs_avg,
    (SELECT COUNT(*) FROM order_items WHERE product_id = p.id) AS times_ordered
FROM products p;

-- Correlated subquery
SELECT 
    c.name,
    (SELECT MAX(o.total) 
     FROM orders o 
     WHERE o.customer_id = c.id) AS max_order
FROM customers c;
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

-- Recursive CTE (for hierarchical data)
WITH RECURSIVE category_path AS (
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
        CONCAT(cp.path, ' > ', c.name) AS path,
        cp.level + 1
    FROM categories c
    INNER JOIN category_path cp ON c.parent_id = cp.id
)
SELECT * FROM category_path ORDER BY path;

-- Recursive CTE for employee hierarchy
WITH RECURSIVE org_chart AS (
    SELECT 
        id, 
        name, 
        manager_id, 
        1 AS level,
        CAST(name AS CHAR(500)) AS hierarchy
    FROM employees
    WHERE manager_id IS NULL
    
    UNION ALL
    
    SELECT 
        e.id,
        e.name,
        e.manager_id,
        oc.level + 1,
        CONCAT(oc.hierarchy, ' -> ', e.name)
    FROM employees e
    INNER JOIN org_chart oc ON e.manager_id = oc.id
)
SELECT * FROM org_chart ORDER BY hierarchy;
```

---

## 9. Aggregation & Grouping

### Aggregate Functions

```sql
-- Basic aggregates
SELECT 
    COUNT(*) AS total_products,              -- Count all rows
    COUNT(description) AS with_description,  -- Count non-NULL
    COUNT(DISTINCT category_id) AS categories,
    SUM(price) AS total_value,
    AVG(price) AS average_price,
    MIN(price) AS cheapest,
    MAX(price) AS most_expensive,
    STD(price) AS price_std_dev,             -- Standard deviation
    VARIANCE(price) AS price_variance
FROM products;

-- Group concatenation
SELECT 
    category_id,
    GROUP_CONCAT(name ORDER BY name SEPARATOR ', ') AS product_names
FROM products
GROUP BY category_id;

-- With custom separator and limit
SELECT 
    category_id,
    GROUP_CONCAT(DISTINCT name ORDER BY name DESC SEPARATOR ' | ') AS products
FROM products
GROUP BY category_id;
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
    YEAR(order_date) AS year,
    MONTH(order_date) AS month,
    COUNT(*) AS orders,
    SUM(total) AS revenue
FROM orders
GROUP BY YEAR(order_date), MONTH(order_date)
ORDER BY year, month;

-- Group by expression
SELECT 
    CASE 
        WHEN price < 50 THEN 'Budget'
        WHEN price < 200 THEN 'Mid-range'
        ELSE 'Premium'
    END AS price_tier,
    COUNT(*) AS count
FROM products
GROUP BY price_tier;
```

### HAVING Clause

```sql
-- Filter groups (not rows)
SELECT 
    category_id,
    COUNT(*) AS product_count,
    AVG(price) AS avg_price
FROM products
GROUP BY category_id
HAVING COUNT(*) > 5 AND AVG(price) > 100;

-- WHERE vs HAVING
-- WHERE filters rows BEFORE grouping
-- HAVING filters groups AFTER aggregation

SELECT 
    category_id,
    COUNT(*) AS product_count,
    AVG(price) AS avg_price
FROM products
WHERE is_active = TRUE           -- Filter rows first
GROUP BY category_id
HAVING AVG(price) > 100;         -- Then filter groups
```

### GROUP BY with ROLLUP

```sql
-- ROLLUP generates subtotals and grand total
SELECT 
    YEAR(order_date) AS year,
    MONTH(order_date) AS month,
    COUNT(*) AS orders,
    SUM(total) AS revenue
FROM orders
GROUP BY YEAR(order_date), MONTH(order_date) WITH ROLLUP;

-- Result includes:
-- - Per month rows
-- - Per year subtotals (month = NULL)
-- - Grand total (year = NULL, month = NULL)

-- Use GROUPING() to identify rollup rows
SELECT 
    IF(GROUPING(year), 'All Years', year) AS year,
    IF(GROUPING(month), 'All Months', month) AS month,
    COUNT(*) AS orders,
    SUM(total) AS revenue
FROM (
    SELECT 
        YEAR(order_date) AS year,
        MONTH(order_date) AS month,
        total
    FROM orders
) AS dated_orders
GROUP BY year, month WITH ROLLUP;
```

### Window Functions (MySQL 8.0+)

```sql
-- ROW_NUMBER - assign sequential numbers
SELECT 
    name,
    category_id,
    price,
    ROW_NUMBER() OVER (ORDER BY price DESC) AS price_rank,
    ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC) AS category_rank
FROM products;

-- RANK and DENSE_RANK
SELECT 
    name,
    price,
    RANK() OVER (ORDER BY price DESC) AS rank_with_gaps,
    DENSE_RANK() OVER (ORDER BY price DESC) AS rank_no_gaps
FROM products;
-- If two products have same price:
-- RANK: 1, 1, 3, 4, 4, 6
-- DENSE_RANK: 1, 1, 2, 3, 3, 4

-- Running totals
SELECT 
    order_date,
    total,
    SUM(total) OVER (ORDER BY order_date) AS running_total,
    AVG(total) OVER (ORDER BY order_date ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS moving_avg_7day
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
    FIRST_VALUE(name) OVER (PARTITION BY category_id ORDER BY price DESC) AS most_expensive,
    LAST_VALUE(name) OVER (
        PARTITION BY category_id 
        ORDER BY price DESC 
        RANGE BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
    ) AS cheapest
FROM products;

-- NTH_VALUE
SELECT 
    name,
    price,
    NTH_VALUE(name, 2) OVER (ORDER BY price DESC) AS second_most_expensive
FROM products;

-- NTILE - divide into buckets
SELECT 
    name,
    price,
    NTILE(4) OVER (ORDER BY price) AS price_quartile
FROM products;

-- Percentage of total
SELECT 
    category_id,
    name,
    price,
    price / SUM(price) OVER (PARTITION BY category_id) * 100 AS pct_of_category,
    price / SUM(price) OVER () * 100 AS pct_of_total
FROM products;
```

---

## 10. Indexes & Performance

### Understanding Indexes

```sql
-- Index is a data structure that improves query speed
-- Like a book's index - helps find pages quickly

-- Without index: Full table scan (read every row)
-- With index: Direct lookup (much faster)
```

### Types of Indexes

```sql
-- 1. Primary Key Index (automatically created)
CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,  -- Clustered index
    email VARCHAR(100)
);

-- 2. Unique Index
CREATE UNIQUE INDEX idx_email ON users(email);
-- Or
ALTER TABLE users ADD UNIQUE INDEX idx_email (email);

-- 3. Regular (Non-unique) Index
CREATE INDEX idx_created ON users(created_at);

-- 4. Composite Index (multiple columns)
CREATE INDEX idx_name_date ON orders(customer_id, order_date);

-- 5. Full-Text Index (for text search)
CREATE FULLTEXT INDEX idx_search ON products(name, description);

-- 6. Spatial Index (for GIS data)
CREATE SPATIAL INDEX idx_coordinates ON locations(point_column);

-- 7. Prefix Index (for long strings)
CREATE INDEX idx_email_prefix ON users(email(50));
```

### Creating and Managing Indexes

```sql
-- Create index on existing table
CREATE INDEX idx_status ON orders(status);

-- Create unique index
CREATE UNIQUE INDEX idx_sku ON products(sku);

-- Create composite index
CREATE INDEX idx_category_price ON products(category_id, price);

-- Show indexes
SHOW INDEX FROM orders;

-- Show index statistics
SHOW INDEX FROM orders WITH INFORMATION SCHEMA;

-- Drop index
DROP INDEX idx_status ON orders;
-- Or
ALTER TABLE orders DROP INDEX idx_status;

-- Rename index (MySQL 5.7+)
ALTER TABLE orders RENAME INDEX idx_old TO idx_new;
```

### Index Design Best Practices

```sql
-- 1. Index columns used in WHERE clauses
-- Frequently filtered columns
CREATE INDEX idx_status ON orders(status);
CREATE INDEX idx_date ON orders(order_date);

-- 2. Index columns used in JOINs
-- Foreign keys should generally be indexed
CREATE INDEX idx_customer ON orders(customer_id);

-- 3. Index columns used in ORDER BY
CREATE INDEX idx_created ON products(created_at DESC);

-- 4. Composite index - column order matters!
-- Most selective column first OR leftmost prefix rule
CREATE INDEX idx_composite ON orders(customer_id, status, order_date);

-- This index helps with:
-- WHERE customer_id = 1
-- WHERE customer_id = 1 AND status = 'pending'
-- WHERE customer_id = 1 AND status = 'pending' AND order_date > '2024-01-01'

-- NOT helpful for:
-- WHERE status = 'pending'  (doesn't use leftmost column)
-- WHERE order_date > '2024-01-01'  (skips customer_id and status)

-- 5. Covering index (includes all columns needed)
CREATE INDEX idx_covering ON orders(customer_id, order_date, total);
-- Query can be satisfied entirely from index
SELECT order_date, total FROM orders WHERE customer_id = 1;
```

### EXPLAIN - Analyzing Queries

```sql
-- Basic EXPLAIN
EXPLAIN SELECT * FROM orders WHERE customer_id = 100;

-- EXPLAIN output columns:
-- id: Query identifier
-- select_type: SIMPLE, PRIMARY, SUBQUERY, etc.
-- table: Table being accessed
-- type: Join type (from best to worst):
--       system > const > eq_ref > ref > range > index > ALL
-- possible_keys: Indexes that might be used
-- key: Index actually used
-- key_len: Length of index used
-- ref: Columns compared to index
-- rows: Estimated rows to examine
-- Extra: Additional information

-- Detailed EXPLAIN
EXPLAIN FORMAT=JSON SELECT * FROM orders WHERE customer_id = 100;

-- EXPLAIN ANALYZE (MySQL 8.0.18+) - actually executes
EXPLAIN ANALYZE SELECT * FROM orders WHERE customer_id = 100;
```

### Query Optimization Examples

```sql
-- BAD: Function on indexed column (can't use index)
SELECT * FROM users WHERE YEAR(created_at) = 2024;

-- GOOD: Use range comparison
SELECT * FROM users 
WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';

-- BAD: Leading wildcard (can't use index)
SELECT * FROM products WHERE name LIKE '%phone%';

-- GOOD: Use FULLTEXT search
SELECT * FROM products 
WHERE MATCH(name, description) AGAINST('phone' IN NATURAL LANGUAGE MODE);

-- BAD: Implicit type conversion (may skip index)
SELECT * FROM users WHERE phone = 1234567890;  -- phone is VARCHAR

-- GOOD: Use correct type
SELECT * FROM users WHERE phone = '1234567890';

-- BAD: OR on different columns (often full scan)
SELECT * FROM products WHERE name = 'iPhone' OR category_id = 5;

-- GOOD: Use UNION
SELECT * FROM products WHERE name = 'iPhone'
UNION
SELECT * FROM products WHERE category_id = 5;

-- BAD: SELECT * when you only need some columns
SELECT * FROM orders WHERE customer_id = 1;

-- GOOD: Select only needed columns (use covering index)
SELECT order_date, total FROM orders WHERE customer_id = 1;
```

### Index Hints

```sql
-- Force using a specific index
SELECT * FROM orders FORCE INDEX (idx_customer)
WHERE customer_id = 100 AND status = 'pending';

-- Suggest an index
SELECT * FROM orders USE INDEX (idx_customer)
WHERE customer_id = 100;

-- Ignore an index
SELECT * FROM orders IGNORE INDEX (idx_status)
WHERE status = 'pending';
```

---

## 11. Views

### Creating Views

```sql
-- Basic view
CREATE VIEW active_products AS
SELECT id, name, price, stock_quantity
FROM products
WHERE is_active = TRUE AND stock_quantity > 0;

-- Using a view
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
LEFT JOIN products p ON c.id = p.category_id AND p.is_active = TRUE
GROUP BY c.id, c.name;
```

### View Options

```sql
-- View with CHECK OPTION (enforces WHERE condition on INSERT/UPDATE)
CREATE VIEW active_users AS
SELECT * FROM users WHERE is_active = TRUE
WITH CHECK OPTION;

-- Attempting to insert inactive user will fail
INSERT INTO active_users (username, email, is_active) 
VALUES ('test', 'test@example.com', FALSE);  -- Error!

-- View with LOCAL CHECK OPTION (only checks this view)
CREATE VIEW premium_active_users AS
SELECT * FROM active_users WHERE account_type = 'premium'
WITH LOCAL CHECK OPTION;

-- View with CASCADED CHECK OPTION (checks all underlying views)
CREATE VIEW premium_active_users AS
SELECT * FROM active_users WHERE account_type = 'premium'
WITH CASCADED CHECK OPTION;

-- Updatable view (can INSERT, UPDATE, DELETE)
CREATE VIEW simple_products AS
SELECT id, name, price FROM products;
-- Must not contain: DISTINCT, GROUP BY, HAVING, UNION, aggregate functions, subqueries in FROM

-- Read-only view (cannot be updated)
-- Views with JOINs, DISTINCT, GROUP BY, etc.
```

### Managing Views

```sql
-- Show views
SHOW FULL TABLES WHERE Table_type = 'VIEW';

-- Show view definition
SHOW CREATE VIEW active_products;

-- Modify view
CREATE OR REPLACE VIEW active_products AS
SELECT id, name, price, stock_quantity, category_id
FROM products
WHERE is_active = TRUE;

-- Or use ALTER
ALTER VIEW active_products AS
SELECT id, name, price
FROM products
WHERE is_active = TRUE AND price > 0;

-- Drop view
DROP VIEW IF EXISTS active_products;

-- Drop multiple views
DROP VIEW view1, view2, view3;
```

### View Use Cases

```sql
-- 1. Simplify complex queries
CREATE VIEW monthly_sales AS
SELECT 
    DATE_FORMAT(order_date, '%Y-%m') AS month,
    COUNT(*) AS order_count,
    SUM(total) AS revenue,
    AVG(total) AS avg_order_value
FROM orders
WHERE status = 'completed'
GROUP BY DATE_FORMAT(order_date, '%Y-%m');

-- 2. Row-level security
CREATE VIEW user_orders AS
SELECT * FROM orders WHERE customer_id = @current_user_id;

-- 3. Hide complexity from applications
CREATE VIEW product_catalog AS
SELECT 
    p.id,
    p.name,
    p.description,
    p.price,
    c.name AS category,
    COALESCE(AVG(r.rating), 0) AS avg_rating,
    COUNT(r.id) AS review_count
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
LEFT JOIN reviews r ON p.id = r.product_id
WHERE p.is_active = TRUE
GROUP BY p.id, p.name, p.description, p.price, c.name;

-- 4. Backward compatibility after schema changes
-- Old table: users(id, name, email)
-- New table: users(id, first_name, last_name, email)
CREATE VIEW users_legacy AS
SELECT 
    id,
    CONCAT(first_name, ' ', last_name) AS name,
    email
FROM users;
```

---

## 12. Stored Procedures & Functions

### Stored Procedures

```sql
-- Change delimiter (needed to define procedure body)
DELIMITER //

-- Basic stored procedure
CREATE PROCEDURE GetAllProducts()
BEGIN
    SELECT * FROM products WHERE is_active = TRUE;
END //

DELIMITER ;

-- Call procedure
CALL GetAllProducts();

-- Procedure with IN parameters
DELIMITER //

CREATE PROCEDURE GetProductsByCategory(IN category_id_param INT)
BEGIN
    SELECT * FROM products 
    WHERE category_id = category_id_param AND is_active = TRUE;
END //

DELIMITER ;

CALL GetProductsByCategory(5);

-- Procedure with OUT parameters
DELIMITER //

CREATE PROCEDURE GetProductStats(
    IN category_id_param INT,
    OUT product_count INT,
    OUT avg_price DECIMAL(10,2)
)
BEGIN
    SELECT COUNT(*), AVG(price)
    INTO product_count, avg_price
    FROM products
    WHERE category_id = category_id_param;
END //

DELIMITER ;

-- Call with OUT parameters
CALL GetProductStats(5, @count, @avg);
SELECT @count AS product_count, @avg AS average_price;

-- Procedure with INOUT parameters
DELIMITER //

CREATE PROCEDURE DoubleValue(INOUT val INT)
BEGIN
    SET val = val * 2;
END //

DELIMITER ;

SET @num = 5;
CALL DoubleValue(@num);
SELECT @num;  -- Returns 10
```

### Procedure with Logic

```sql
DELIMITER //

CREATE PROCEDURE ProcessOrder(
    IN p_customer_id INT,
    IN p_product_id INT,
    IN p_quantity INT,
    OUT p_result VARCHAR(100)
)
BEGIN
    DECLARE v_stock INT;
    DECLARE v_price DECIMAL(10,2);
    DECLARE v_order_id INT;
    
    -- Start transaction
    START TRANSACTION;
    
    -- Check stock
    SELECT stock_quantity, price INTO v_stock, v_price
    FROM products WHERE id = p_product_id FOR UPDATE;
    
    IF v_stock IS NULL THEN
        SET p_result = 'Product not found';
        ROLLBACK;
    ELSEIF v_stock < p_quantity THEN
        SET p_result = 'Insufficient stock';
        ROLLBACK;
    ELSE
        -- Create order
        INSERT INTO orders (customer_id, total, status)
        VALUES (p_customer_id, v_price * p_quantity, 'pending');
        
        SET v_order_id = LAST_INSERT_ID();
        
        -- Add order item
        INSERT INTO order_items (order_id, product_id, quantity, unit_price)
        VALUES (v_order_id, p_product_id, p_quantity, v_price);
        
        -- Update stock
        UPDATE products 
        SET stock_quantity = stock_quantity - p_quantity
        WHERE id = p_product_id;
        
        SET p_result = CONCAT('Order created: ', v_order_id);
        COMMIT;
    END IF;
END //

DELIMITER ;

CALL ProcessOrder(1, 5, 2, @result);
SELECT @result;
```

### Stored Functions

```sql
DELIMITER //

-- Simple function
CREATE FUNCTION GetProductCount(category_id_param INT)
RETURNS INT
DETERMINISTIC
READS SQL DATA
BEGIN
    DECLARE cnt INT;
    SELECT COUNT(*) INTO cnt FROM products WHERE category_id = category_id_param;
    RETURN cnt;
END //

DELIMITER ;

-- Use function in queries
SELECT 
    c.name,
    GetProductCount(c.id) AS product_count
FROM categories c;

-- Function with calculations
DELIMITER //

CREATE FUNCTION CalculateDiscount(
    price DECIMAL(10,2),
    discount_percent INT
)
RETURNS DECIMAL(10,2)
DETERMINISTIC
NO SQL
BEGIN
    RETURN price * (1 - discount_percent / 100);
END //

DELIMITER ;

SELECT name, price, CalculateDiscount(price, 20) AS discounted_price
FROM products;

-- Function characteristics:
-- DETERMINISTIC: Same input always produces same output
-- NOT DETERMINISTIC: May produce different results
-- READS SQL DATA: Reads but doesn't modify data
-- MODIFIES SQL DATA: May modify data
-- NO SQL: No SQL statements
```

### Managing Procedures and Functions

```sql
-- Show all procedures
SHOW PROCEDURE STATUS WHERE Db = 'myapp';

-- Show procedure definition
SHOW CREATE PROCEDURE GetAllProducts;

-- Show all functions
SHOW FUNCTION STATUS WHERE Db = 'myapp';

-- Show function definition
SHOW CREATE FUNCTION GetProductCount;

-- Drop procedure
DROP PROCEDURE IF EXISTS GetAllProducts;

-- Drop function
DROP FUNCTION IF EXISTS GetProductCount;
```

### Cursors in Stored Procedures

```sql
DELIMITER //

CREATE PROCEDURE ProcessAllPendingOrders()
BEGIN
    DECLARE v_order_id INT;
    DECLARE v_done INT DEFAULT FALSE;
    
    -- Declare cursor
    DECLARE order_cursor CURSOR FOR
        SELECT id FROM orders WHERE status = 'pending';
    
    -- Declare handler for when no more rows
    DECLARE CONTINUE HANDLER FOR NOT FOUND SET v_done = TRUE;
    
    -- Open cursor
    OPEN order_cursor;
    
    -- Loop through rows
    read_loop: LOOP
        FETCH order_cursor INTO v_order_id;
        
        IF v_done THEN
            LEAVE read_loop;
        END IF;
        
        -- Process each order
        UPDATE orders SET status = 'processing' WHERE id = v_order_id;
        
    END LOOP;
    
    -- Close cursor
    CLOSE order_cursor;
END //

DELIMITER ;
```

---

## 13. Triggers

### Creating Triggers

```sql
-- Trigger timing: BEFORE or AFTER
-- Trigger event: INSERT, UPDATE, DELETE

DELIMITER //

-- BEFORE INSERT trigger
CREATE TRIGGER before_user_insert
BEFORE INSERT ON users
FOR EACH ROW
BEGIN
    SET NEW.created_at = NOW();
    SET NEW.updated_at = NOW();
    SET NEW.email = LOWER(TRIM(NEW.email));
END //

-- AFTER INSERT trigger
CREATE TRIGGER after_order_insert
AFTER INSERT ON orders
FOR EACH ROW
BEGIN
    INSERT INTO order_audit (order_id, action, action_time)
    VALUES (NEW.id, 'INSERT', NOW());
END //

-- BEFORE UPDATE trigger
CREATE TRIGGER before_product_update
BEFORE UPDATE ON products
FOR EACH ROW
BEGIN
    SET NEW.updated_at = NOW();
    
    -- Store price history if price changed
    IF OLD.price != NEW.price THEN
        INSERT INTO price_history (product_id, old_price, new_price, changed_at)
        VALUES (OLD.id, OLD.price, NEW.price, NOW());
    END IF;
END //

-- AFTER UPDATE trigger
CREATE TRIGGER after_stock_update
AFTER UPDATE ON products
FOR EACH ROW
BEGIN
    IF NEW.stock_quantity = 0 AND OLD.stock_quantity > 0 THEN
        INSERT INTO notifications (message, created_at)
        VALUES (CONCAT('Product out of stock: ', NEW.name), NOW());
    END IF;
END //

-- BEFORE DELETE trigger
CREATE TRIGGER before_customer_delete
BEFORE DELETE ON customers
FOR EACH ROW
BEGIN
    -- Archive before deleting
    INSERT INTO customers_archive 
    SELECT *, NOW() as archived_at FROM customers WHERE id = OLD.id;
END //

-- AFTER DELETE trigger
CREATE TRIGGER after_order_delete
AFTER DELETE ON orders
FOR EACH ROW
BEGIN
    INSERT INTO order_audit (order_id, action, action_time)
    VALUES (OLD.id, 'DELETE', NOW());
END //

DELIMITER ;
```

### Trigger References

```sql
-- OLD and NEW references:
-- INSERT: Only NEW is available
-- UPDATE: Both OLD and NEW are available
-- DELETE: Only OLD is available

DELIMITER //

CREATE TRIGGER log_user_changes
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    INSERT INTO user_changes_log (
        user_id,
        field_changed,
        old_value,
        new_value,
        changed_at
    )
    SELECT 
        NEW.id,
        'email',
        OLD.email,
        NEW.email,
        NOW()
    WHERE OLD.email != NEW.email
    
    UNION ALL
    
    SELECT 
        NEW.id,
        'username',
        OLD.username,
        NEW.username,
        NOW()
    WHERE OLD.username != NEW.username;
END //

DELIMITER ;
```

### Managing Triggers

```sql
-- Show all triggers
SHOW TRIGGERS;

-- Show triggers for specific table
SHOW TRIGGERS LIKE 'users';

-- Show trigger definition
SHOW CREATE TRIGGER before_user_insert;

-- Drop trigger
DROP TRIGGER IF EXISTS before_user_insert;

-- List triggers from information_schema
SELECT 
    TRIGGER_NAME,
    EVENT_MANIPULATION,
    EVENT_OBJECT_TABLE,
    ACTION_TIMING,
    ACTION_STATEMENT
FROM information_schema.TRIGGERS
WHERE TRIGGER_SCHEMA = 'myapp';
```

### Trigger Best Practices

```sql
-- 1. Keep triggers simple and fast
-- Triggers execute for EVERY affected row

-- 2. Avoid triggers that call other triggers (cascading)

-- 3. Use triggers for:
--    - Audit logging
--    - Data validation
--    - Maintaining derived data
--    - Enforcing business rules

-- 4. Avoid triggers for:
--    - Complex business logic (use procedures)
--    - Sending emails/notifications (use queues)
--    - Time-consuming operations

-- 5. Document triggers thoroughly

-- Example: Audit trigger with good practices
DELIMITER //

CREATE TRIGGER audit_orders
AFTER INSERT OR UPDATE OR DELETE ON orders
FOR EACH ROW
BEGIN
    DECLARE v_action VARCHAR(10);
    DECLARE v_order_id INT;
    DECLARE v_data JSON;
    
    IF NEW.id IS NOT NULL THEN
        SET v_order_id = NEW.id;
        SET v_data = JSON_OBJECT(
            'customer_id', NEW.customer_id,
            'total', NEW.total,
            'status', NEW.status
        );
        IF OLD.id IS NULL THEN
            SET v_action = 'INSERT';
        ELSE
            SET v_action = 'UPDATE';
        END IF;
    ELSE
        SET v_order_id = OLD.id;
        SET v_action = 'DELETE';
        SET v_data = JSON_OBJECT(
            'customer_id', OLD.customer_id,
            'total', OLD.total,
            'status', OLD.status
        );
    END IF;
    
    INSERT INTO audit_log (table_name, record_id, action, data, created_at)
    VALUES ('orders', v_order_id, v_action, v_data, NOW());
END //

DELIMITER ;
```

---

## 14. Transactions

### Transaction Basics

```sql
-- Start a transaction
START TRANSACTION;
-- Or
BEGIN;

-- Execute queries
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- Commit (save changes)
COMMIT;

-- Or rollback (undo changes)
ROLLBACK;
```

### ACID Properties

```
A - Atomicity:    All operations succeed or all fail
C - Consistency:  Database remains in valid state
I - Isolation:    Concurrent transactions don't interfere
D - Durability:   Committed data persists through failures
```

### Transaction Example

```sql
-- Transfer money between accounts
DELIMITER //

CREATE PROCEDURE TransferMoney(
    IN from_account INT,
    IN to_account INT,
    IN amount DECIMAL(10,2),
    OUT result VARCHAR(100)
)
BEGIN
    DECLARE from_balance DECIMAL(10,2);
    DECLARE EXIT HANDLER FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        SET result = 'Transaction failed due to error';
    END;
    
    START TRANSACTION;
    
    -- Lock the rows we'll modify
    SELECT balance INTO from_balance 
    FROM accounts 
    WHERE id = from_account 
    FOR UPDATE;
    
    IF from_balance < amount THEN
        ROLLBACK;
        SET result = 'Insufficient funds';
    ELSE
        -- Debit from source
        UPDATE accounts 
        SET balance = balance - amount 
        WHERE id = from_account;
        
        -- Credit to destination
        UPDATE accounts 
        SET balance = balance + amount 
        WHERE id = to_account;
        
        -- Log the transfer
        INSERT INTO transfers (from_account, to_account, amount, transfer_date)
        VALUES (from_account, to_account, amount, NOW());
        
        COMMIT;
        SET result = 'Transfer successful';
    END IF;
END //

DELIMITER ;

CALL TransferMoney(1, 2, 100.00, @result);
SELECT @result;
```

### Isolation Levels

```sql
-- Check current isolation level
SELECT @@transaction_isolation;
-- Or (older MySQL versions)
SELECT @@tx_isolation;

-- Set isolation level for session
SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;
SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED;
SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ;  -- Default
SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE;

-- Set for next transaction only
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;
```

### Isolation Levels Explained

```sql
-- 1. READ UNCOMMITTED (lowest isolation)
-- Can see uncommitted changes from other transactions (dirty reads)
-- Use: Only when dirty reads are acceptable

-- 2. READ COMMITTED
-- Only sees committed data
-- Same query may return different results (non-repeatable reads)
-- Use: General web applications

-- 3. REPEATABLE READ (MySQL default)
-- Same query returns same results within transaction
-- Phantom reads possible (new rows appear)
-- Use: Most applications, good default

-- 4. SERIALIZABLE (highest isolation)
-- Complete isolation, transactions appear sequential
-- Use: Financial transactions, critical operations
```

### Locking

```sql
-- Shared lock (SELECT ... FOR SHARE)
-- Multiple transactions can read, but none can modify
SELECT * FROM products WHERE id = 1 FOR SHARE;
-- Or older syntax:
SELECT * FROM products WHERE id = 1 LOCK IN SHARE MODE;

-- Exclusive lock (SELECT ... FOR UPDATE)
-- Only one transaction can access
SELECT * FROM products WHERE id = 1 FOR UPDATE;

-- Skip locked rows (MySQL 8.0+)
SELECT * FROM products WHERE status = 'available' 
FOR UPDATE SKIP LOCKED LIMIT 1;

-- Wait with timeout
SELECT * FROM products WHERE id = 1 
FOR UPDATE WAIT 5;  -- Wait max 5 seconds

-- Nowait - fail immediately if locked
SELECT * FROM products WHERE id = 1 
FOR UPDATE NOWAIT;

-- Table lock
LOCK TABLES products WRITE;
-- Do operations
UNLOCK TABLES;

LOCK TABLES products READ, orders WRITE;
UNLOCK TABLES;
```

### Savepoints

```sql
START TRANSACTION;

INSERT INTO orders (customer_id, total) VALUES (1, 100);

SAVEPOINT order_created;

INSERT INTO order_items (order_id, product_id, quantity) VALUES (LAST_INSERT_ID(), 1, 5);

-- Oops, wrong product
ROLLBACK TO SAVEPOINT order_created;

-- Insert correct item
INSERT INTO order_items (order_id, product_id, quantity) VALUES (LAST_INSERT_ID(), 2, 3);

COMMIT;

-- Release savepoint (optional, freed on commit/rollback)
RELEASE SAVEPOINT order_created;
```

### Autocommit

```sql
-- Check autocommit status
SELECT @@autocommit;

-- Disable autocommit
SET autocommit = 0;
-- Now every statement needs explicit COMMIT

-- Enable autocommit (default)
SET autocommit = 1;
-- Each statement is automatically committed
```

---

## 15. User Management & Security

### Creating Users

```sql
-- Create user with password
CREATE USER 'username'@'localhost' IDENTIFIED BY 'password123';

-- Create user allowing connection from any host
CREATE USER 'username'@'%' IDENTIFIED BY 'password123';

-- Create user with specific host
CREATE USER 'username'@'192.168.1.%' IDENTIFIED BY 'password123';

-- Create user with password expiration
CREATE USER 'username'@'localhost' 
IDENTIFIED BY 'password123'
PASSWORD EXPIRE INTERVAL 90 DAY;

-- Create user with account locking
CREATE USER 'username'@'localhost' 
IDENTIFIED BY 'password123'
ACCOUNT LOCK;

-- Create user with resource limits
CREATE USER 'username'@'localhost' 
IDENTIFIED BY 'password123'
WITH 
    MAX_QUERIES_PER_HOUR 1000
    MAX_UPDATES_PER_HOUR 100
    MAX_CONNECTIONS_PER_HOUR 50
    MAX_USER_CONNECTIONS 5;
```

### Granting Privileges

```sql
-- Grant all privileges on database
GRANT ALL PRIVILEGES ON myapp.* TO 'username'@'localhost';

-- Grant specific privileges
GRANT SELECT, INSERT, UPDATE ON myapp.* TO 'username'@'localhost';

-- Grant on specific table
GRANT SELECT, INSERT ON myapp.orders TO 'username'@'localhost';

-- Grant on specific columns
GRANT SELECT (id, name, price), UPDATE (price) 
ON myapp.products 
TO 'username'@'localhost';

-- Grant with ability to grant to others
GRANT SELECT ON myapp.* TO 'username'@'localhost' WITH GRANT OPTION;

-- Common privilege sets
-- Read-only
GRANT SELECT ON myapp.* TO 'readonly'@'localhost';

-- Read-write (no structure changes)
GRANT SELECT, INSERT, UPDATE, DELETE ON myapp.* TO 'readwrite'@'localhost';

-- Application user
GRANT SELECT, INSERT, UPDATE, DELETE, EXECUTE ON myapp.* TO 'app'@'localhost';

-- Developer
GRANT ALL PRIVILEGES ON myapp.* TO 'developer'@'localhost';

-- Apply privilege changes
FLUSH PRIVILEGES;
```

### Revoking Privileges

```sql
-- Revoke specific privileges
REVOKE INSERT, UPDATE ON myapp.* FROM 'username'@'localhost';

-- Revoke all privileges
REVOKE ALL PRIVILEGES ON myapp.* FROM 'username'@'localhost';

-- Revoke grant option
REVOKE GRANT OPTION ON myapp.* FROM 'username'@'localhost';
```

### Managing Users

```sql
-- Show all users
SELECT User, Host FROM mysql.user;

-- Show privileges for user
SHOW GRANTS FOR 'username'@'localhost';

-- Show current user
SELECT CURRENT_USER();

-- Change password
ALTER USER 'username'@'localhost' IDENTIFIED BY 'new_password';

-- Force password change on next login
ALTER USER 'username'@'localhost' PASSWORD EXPIRE;

-- Unlock user
ALTER USER 'username'@'localhost' ACCOUNT UNLOCK;

-- Rename user
RENAME USER 'oldname'@'localhost' TO 'newname'@'localhost';

-- Drop user
DROP USER 'username'@'localhost';
DROP USER IF EXISTS 'username'@'localhost';
```

### Roles (MySQL 8.0+)

```sql
-- Create roles
CREATE ROLE 'app_read', 'app_write', 'app_admin';

-- Grant privileges to roles
GRANT SELECT ON myapp.* TO 'app_read';
GRANT INSERT, UPDATE, DELETE ON myapp.* TO 'app_write';
GRANT ALL PRIVILEGES ON myapp.* TO 'app_admin';

-- Grant roles to users
GRANT 'app_read' TO 'user1'@'localhost';
GRANT 'app_read', 'app_write' TO 'user2'@'localhost';
GRANT 'app_admin' TO 'admin'@'localhost';

-- Set default roles
SET DEFAULT ROLE 'app_read' TO 'user1'@'localhost';
SET DEFAULT ROLE ALL TO 'admin'@'localhost';

-- Activate roles in session
SET ROLE 'app_read';
SET ROLE ALL;
SET ROLE NONE;

-- Show roles
SELECT * FROM mysql.role_edges;

-- Drop role
DROP ROLE 'app_read';
```

### Security Best Practices

```sql
-- 1. Remove anonymous users
DELETE FROM mysql.user WHERE User='';

-- 2. Remove test database
DROP DATABASE IF EXISTS test;

-- 3. Disable remote root access
DELETE FROM mysql.user WHERE User='root' AND Host NOT IN ('localhost', '127.0.0.1', '::1');

-- 4. Use strong passwords
CREATE USER 'app'@'localhost' IDENTIFIED BY 'C0mpl3x!P@ssw0rd#2024';

-- 5. Grant minimum necessary privileges
-- Bad: GRANT ALL
-- Good: GRANT SELECT, INSERT, UPDATE ON specific_tables

-- 6. Use SSL/TLS connections
CREATE USER 'secure'@'%' 
IDENTIFIED BY 'password' 
REQUIRE SSL;

-- Or require specific certificate
CREATE USER 'secure'@'%' 
IDENTIFIED BY 'password' 
REQUIRE X509;

-- 7. Enable audit logging (MySQL Enterprise)
-- Or use triggers/application logging

-- 8. Regular security audits
-- Check for overprivileged users
SELECT User, Host, Super_priv, Grant_priv 
FROM mysql.user 
WHERE Super_priv = 'Y' OR Grant_priv = 'Y';
```

---

## 16. Backup & Recovery

### mysqldump - Logical Backup

```bash
# Backup single database
mysqldump -u root -p myapp > myapp_backup.sql

# Backup with routines and triggers
mysqldump -u root -p --routines --triggers myapp > myapp_full.sql

# Backup specific tables
mysqldump -u root -p myapp users orders > tables_backup.sql

# Backup all databases
mysqldump -u root -p --all-databases > all_databases.sql

# Backup with compression
mysqldump -u root -p myapp | gzip > myapp_backup.sql.gz

# Backup for InnoDB (consistent snapshot)
mysqldump -u root -p --single-transaction --quick myapp > myapp_backup.sql

# Backup structure only (no data)
mysqldump -u root -p --no-data myapp > myapp_structure.sql

# Backup data only (no structure)
mysqldump -u root -p --no-create-info myapp > myapp_data.sql

# Backup with WHERE condition
mysqldump -u root -p myapp orders --where="order_date > '2024-01-01'" > recent_orders.sql
```

### Restoring from Backup

```bash
# Restore from SQL file
mysql -u root -p myapp < myapp_backup.sql

# Restore from compressed backup
gunzip < myapp_backup.sql.gz | mysql -u root -p myapp

# Restore with progress (pv utility)
pv myapp_backup.sql | mysql -u root -p myapp

# Restore specific tables
mysql -u root -p myapp < tables_backup.sql
```

### Binary Log Backup (Point-in-Time Recovery)

```sql
-- Enable binary logging in my.cnf
-- [mysqld]
-- log-bin = mysql-bin
-- binlog_format = ROW
-- server-id = 1

-- Show binary log status
SHOW MASTER STATUS;

-- List binary logs
SHOW BINARY LOGS;

-- View binary log contents
SHOW BINLOG EVENTS IN 'mysql-bin.000001';
```

```bash
# Backup binary logs
mysqlbinlog mysql-bin.000001 > binlog_backup.sql

# Point-in-time recovery
# 1. Restore last full backup
mysql -u root -p myapp < myapp_backup.sql

# 2. Apply binary logs up to specific time
mysqlbinlog --stop-datetime="2024-03-15 14:00:00" mysql-bin.000001 | mysql -u root -p

# 3. Or apply up to specific position
mysqlbinlog --stop-position=107 mysql-bin.000001 | mysql -u root -p
```

### Physical Backup (MySQL Enterprise / Percona XtraBackup)

```bash
# Using Percona XtraBackup
# Full backup
xtrabackup --backup --target-dir=/backup/full

# Prepare backup for restoration
xtrabackup --prepare --target-dir=/backup/full

# Restore
xtrabackup --copy-back --target-dir=/backup/full

# Incremental backup
xtrabackup --backup --target-dir=/backup/inc1 --incremental-basedir=/backup/full
```

### Backup Strategy

```
┌─────────────────────────────────────────────────────────────┐
│                    BACKUP STRATEGY                          │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Daily:                                                     │
│  └── Full mysqldump with --single-transaction               │
│                                                             │
│  Hourly:                                                    │
│  └── Binary log backup                                      │
│                                                             │
│  Retention:                                                 │
│  ├── Daily backups: Keep 30 days                           │
│  ├── Weekly backups: Keep 12 weeks                         │
│  └── Monthly backups: Keep 12 months                       │
│                                                             │
│  Storage:                                                   │
│  ├── Local: Fast recovery                                  │
│  └── Remote: Disaster recovery (S3, GCS, Azure)            │
│                                                             │
│  Testing:                                                   │
│  └── Monthly restore test to verify backup integrity        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Automated Backup Script

```bash
#!/bin/bash
# backup_mysql.sh

# Configuration
DB_USER="backup_user"
DB_PASS="backup_password"
DB_NAME="myapp"
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30

# Create backup directory
mkdir -p $BACKUP_DIR

# Perform backup
mysqldump -u $DB_USER -p$DB_PASS \
    --single-transaction \
    --routines \
    --triggers \
    --quick \
    $DB_NAME | gzip > $BACKUP_DIR/${DB_NAME}_${DATE}.sql.gz

# Check if backup was successful
if [ $? -eq 0 ]; then
    echo "Backup completed successfully: ${DB_NAME}_${DATE}.sql.gz"
    
    # Upload to S3 (optional)
    aws s3 cp $BACKUP_DIR/${DB_NAME}_${DATE}.sql.gz s3://my-bucket/mysql-backups/
    
    # Delete old backups
    find $BACKUP_DIR -name "*.sql.gz" -mtime +$RETENTION_DAYS -delete
else
    echo "Backup failed!"
    exit 1
fi
```

### Import/Export Large Data

```sql
-- Export to CSV
SELECT * FROM orders
INTO OUTFILE '/tmp/orders.csv'
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n';

-- Import from CSV
LOAD DATA INFILE '/tmp/orders.csv'
INTO TABLE orders
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n';

-- Load with column mapping
LOAD DATA INFILE '/tmp/data.csv'
INTO TABLE products
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\n'
IGNORE 1 ROWS  -- Skip header
(name, price, @category_name)
SET category_id = (SELECT id FROM categories WHERE name = @category_name);
```

---

## 17. Replication

### Replication Overview

```
┌─────────────────────────────────────────────────────────────┐
│                    REPLICATION TYPES                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. Asynchronous (Default)                                  │
│     Master ──────────────► Replica                          │
│     - Replica may lag behind                                │
│     - No performance impact on master                       │
│                                                             │
│  2. Semi-synchronous                                        │
│     Master ◄────ACK─────► Replica                          │
│     - At least one replica confirms                         │
│     - Better durability guarantee                           │
│                                                             │
│  3. Group Replication (MySQL 8.0)                          │
│     Master ◄────────────► Master                           │
│     - Multi-master, all nodes equal                         │
│     - Automatic failover                                    │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Setting Up Basic Replication

**Master Configuration (my.cnf):**

```ini
[mysqld]
server-id = 1
log-bin = mysql-bin
binlog_format = ROW
binlog_do_db = myapp  # Optional: replicate specific database
```

**Master Setup:**

```sql
-- Create replication user
CREATE USER 'repl_user'@'%' IDENTIFIED BY 'repl_password';
GRANT REPLICATION SLAVE ON *.* TO 'repl_user'@'%';
FLUSH PRIVILEGES;

-- Get master status
SHOW MASTER STATUS;
-- Note the File and Position values

-- Lock tables for consistent backup
FLUSH TABLES WITH READ LOCK;

-- Backup database
-- mysqldump -u root -p myapp > backup.sql

-- Unlock
UNLOCK TABLES;
```

**Replica Configuration (my.cnf):**

```ini
[mysqld]
server-id = 2
relay-log = relay-bin
read_only = 1
```

**Replica Setup:**

```sql
-- Restore backup
-- mysql -u root -p myapp < backup.sql

-- Configure replication
CHANGE MASTER TO
    MASTER_HOST = '192.168.1.100',
    MASTER_USER = 'repl_user',
    MASTER_PASSWORD = 'repl_password',
    MASTER_LOG_FILE = 'mysql-bin.000001',
    MASTER_LOG_POS = 154;

-- MySQL 8.0.23+ syntax
CHANGE REPLICATION SOURCE TO
    SOURCE_HOST = '192.168.1.100',
    SOURCE_USER = 'repl_user',
    SOURCE_PASSWORD = 'repl_password',
    SOURCE_LOG_FILE = 'mysql-bin.000001',
    SOURCE_LOG_POS = 154;

-- Start replication
START REPLICA;
-- Or older syntax: START SLAVE;

-- Check replication status
SHOW REPLICA STATUS\G
-- Or: SHOW SLAVE STATUS\G
```

### Monitoring Replication

```sql
-- Check replica status
SHOW REPLICA STATUS\G

-- Key fields to monitor:
-- Slave_IO_Running: Yes
-- Slave_SQL_Running: Yes
-- Seconds_Behind_Master: 0 (or low number)
-- Last_Error: (should be empty)

-- Check master status
SHOW MASTER STATUS;

-- Show connected replicas
SHOW PROCESSLIST;

-- Monitor replication lag
SELECT 
    @master_pos := MASTER_LOG_POS,
    @slave_pos := Relay_Log_Pos,
    @master_pos - @slave_pos AS lag_bytes
FROM information_schema.processlist;
```

### GTID-Based Replication (Recommended)

```ini
# my.cnf on both master and replica
[mysqld]
gtid_mode = ON
enforce_gtid_consistency = ON
```

```sql
-- Configure replica with GTID
CHANGE REPLICATION SOURCE TO
    SOURCE_HOST = '192.168.1.100',
    SOURCE_USER = 'repl_user',
    SOURCE_PASSWORD = 'repl_password',
    SOURCE_AUTO_POSITION = 1;

START REPLICA;
```

### Read Replica for Load Balancing

```sql
-- Application logic (pseudocode):
-- Writes: Connect to master
-- Reads: Connect to replica(s)

-- ProxySQL configuration for automatic routing:
-- mysql_servers:
--   - hostname: master, hostgroup_id: 0 (write)
--   - hostname: replica1, hostgroup_id: 1 (read)
--   - hostname: replica2, hostgroup_id: 1 (read)

-- In application:
-- Write query:
UPDATE users SET name = 'New Name' WHERE id = 1;

-- Read query (can go to replica):
SELECT * FROM users WHERE id = 1;
```

---

## 18. Performance Tuning

### Identifying Performance Issues

```sql
-- Enable slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;  -- Log queries > 2 seconds
SET GLOBAL log_queries_not_using_indexes = 'ON';

-- Show slow query log location
SHOW VARIABLES LIKE 'slow_query_log_file';

-- Show status counters
SHOW GLOBAL STATUS LIKE 'Slow_queries';
SHOW GLOBAL STATUS LIKE 'Questions';

-- Show process list
SHOW FULL PROCESSLIST;

-- Kill long-running query
KILL <process_id>;

-- Performance Schema queries
SELECT * FROM performance_schema.events_statements_summary_by_digest
ORDER BY SUM_TIMER_WAIT DESC LIMIT 10;
```

### Key Configuration Parameters

```ini
# my.cnf - Performance Tuning

[mysqld]
# InnoDB Buffer Pool (70-80% of RAM on dedicated server)
innodb_buffer_pool_size = 4G

# Buffer pool instances (1 per GB, max 64)
innodb_buffer_pool_instances = 4

# Log file size (larger = better write performance)
innodb_log_file_size = 512M

# Flush method (O_DIRECT recommended for Linux)
innodb_flush_method = O_DIRECT

# Transaction durability (1 = safest, 2 = faster)
innodb_flush_log_at_trx_commit = 1

# IO capacity (higher for SSDs)
innodb_io_capacity = 2000
innodb_io_capacity_max = 4000

# Thread concurrency (0 = auto)
innodb_thread_concurrency = 0

# Read/write IO threads
innodb_read_io_threads = 4
innodb_write_io_threads = 4

# Connection settings
max_connections = 200
thread_cache_size = 50

# Query cache (DISABLED in MySQL 8.0)
# query_cache_type = 0

# Temp table settings
tmp_table_size = 256M
max_heap_table_size = 256M

# Sort and join buffers
sort_buffer_size = 4M
join_buffer_size = 4M

# Table cache
table_open_cache = 4000
table_definition_cache = 2000
```

### Query Optimization

```sql
-- 1. Use EXPLAIN to analyze queries
EXPLAIN SELECT * FROM orders WHERE customer_id = 100;

-- 2. Optimize SELECT
-- Bad: SELECT *
SELECT * FROM users;
-- Good: Select only needed columns
SELECT id, name, email FROM users;

-- 3. Avoid SELECT DISTINCT when possible
-- Bad:
SELECT DISTINCT category_id FROM products;
-- Better:
SELECT category_id FROM products GROUP BY category_id;

-- 4. Use indexes effectively
-- Ensure indexes exist on WHERE, JOIN, ORDER BY columns
CREATE INDEX idx_customer ON orders(customer_id);
CREATE INDEX idx_status_date ON orders(status, order_date);

-- 5. Avoid functions on indexed columns
-- Bad:
SELECT * FROM users WHERE YEAR(created_at) = 2024;
-- Good:
SELECT * FROM users WHERE created_at >= '2024-01-01' AND created_at < '2025-01-01';

-- 6. Use LIMIT for large result sets
SELECT * FROM logs ORDER BY created_at DESC LIMIT 100;

-- 7. Optimize JOINs
-- Ensure join columns are indexed
-- Join on integers rather than strings
-- Use INNER JOIN when possible (more efficient than OUTER)

-- 8. Use batch operations
-- Bad: Many single inserts
INSERT INTO logs (message) VALUES ('log1');
INSERT INTO logs (message) VALUES ('log2');
-- Good: Batch insert
INSERT INTO logs (message) VALUES ('log1'), ('log2'), ('log3');

-- 9. Use prepared statements (application code)
-- Reduces parsing overhead for repeated queries

-- 10. Avoid correlated subqueries
-- Bad:
SELECT * FROM products p 
WHERE price > (SELECT AVG(price) FROM products WHERE category_id = p.category_id);
-- Good: Use JOIN
SELECT p.* FROM products p
JOIN (SELECT category_id, AVG(price) as avg_price FROM products GROUP BY category_id) c
ON p.category_id = c.category_id AND p.price > c.avg_price;
```

### Monitoring and Profiling

```sql
-- Enable profiling for session
SET profiling = 1;

-- Run query
SELECT * FROM orders WHERE customer_id = 100;

-- Show profiles
SHOW PROFILES;

-- Show specific profile
SHOW PROFILE FOR QUERY 1;
SHOW PROFILE CPU, BLOCK IO FOR QUERY 1;

-- Performance Schema
SELECT * FROM performance_schema.events_statements_history 
ORDER BY TIMER_WAIT DESC LIMIT 5;

-- Table statistics
SELECT 
    TABLE_NAME,
    TABLE_ROWS,
    DATA_LENGTH,
    INDEX_LENGTH,
    DATA_FREE
FROM information_schema.TABLES
WHERE TABLE_SCHEMA = 'myapp'
ORDER BY DATA_LENGTH DESC;

-- Index usage statistics
SELECT 
    TABLE_NAME,
    INDEX_NAME,
    STAT_VALUE
FROM mysql.innodb_index_stats
WHERE database_name = 'myapp';

-- InnoDB status
SHOW ENGINE INNODB STATUS\G
```

### Query Cache (Note: Deprecated in 8.0)

```sql
-- MySQL 5.7 and earlier only
-- Query cache often causes more harm than good due to invalidation

-- Check cache status
SHOW VARIABLES LIKE 'query_cache%';
SHOW STATUS LIKE 'Qcache%';

-- Disable (recommended)
SET GLOBAL query_cache_type = OFF;
SET GLOBAL query_cache_size = 0;
```

---

## 19. MySQL 8.0+ Features

### Window Functions

```sql
-- Row numbering
SELECT 
    name,
    price,
    ROW_NUMBER() OVER (ORDER BY price DESC) AS price_rank
FROM products;

-- Partitioned ranking
SELECT 
    category_id,
    name,
    price,
    RANK() OVER (PARTITION BY category_id ORDER BY price DESC) AS category_rank
FROM products;

-- Running totals
SELECT 
    order_date,
    total,
    SUM(total) OVER (ORDER BY order_date) AS running_total
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
    LAG(total) OVER (ORDER BY order_date) AS prev_day,
    LEAD(total) OVER (ORDER BY order_date) AS next_day
FROM orders;
```

### Common Table Expressions (CTEs)

```sql
-- Simple CTE
WITH high_value_customers AS (
    SELECT customer_id, SUM(total) as total_spent
    FROM orders
    GROUP BY customer_id
    HAVING SUM(total) > 10000
)
SELECT c.name, h.total_spent
FROM customers c
JOIN high_value_customers h ON c.id = h.customer_id;

-- Recursive CTE
WITH RECURSIVE category_tree AS (
    SELECT id, name, parent_id, 0 AS level
    FROM categories WHERE parent_id IS NULL
    
    UNION ALL
    
    SELECT c.id, c.name, c.parent_id, ct.level + 1
    FROM categories c
    JOIN category_tree ct ON c.parent_id = ct.id
)
SELECT * FROM category_tree;
```

### JSON Enhancements

```sql
-- JSON_TABLE - Convert JSON to table
SELECT jt.*
FROM products,
JSON_TABLE(
    attributes,
    '$' COLUMNS (
        brand VARCHAR(100) PATH '$.brand',
        ram INT PATH '$.ram',
        storage VARCHAR(50) PATH '$.storage'
    )
) AS jt;

-- JSON aggregation
SELECT 
    category_id,
    JSON_ARRAYAGG(name) AS products,
    JSON_OBJECTAGG(id, name) AS product_map
FROM products
GROUP BY category_id;

-- JSON path expressions
SELECT 
    attributes->>'$.brand' AS brand,
    attributes->>'$.specs.ram' AS ram
FROM products;
```

### Invisible Indexes

```sql
-- Create invisible index
CREATE INDEX idx_status ON orders(status) INVISIBLE;

-- Make existing index invisible (for testing removal impact)
ALTER TABLE orders ALTER INDEX idx_status INVISIBLE;

-- Make visible again
ALTER TABLE orders ALTER INDEX idx_status VISIBLE;

-- Query uses visible indexes by default
-- Force using invisible indexes
SET optimizer_switch = 'use_invisible_indexes=on';
```

### Descending Indexes

```sql
-- Create descending index
CREATE INDEX idx_date_desc ON orders(order_date DESC);

-- Mixed ascending/descending
CREATE INDEX idx_category_price ON products(category_id ASC, price DESC);

-- Useful for ORDER BY optimization
SELECT * FROM products 
WHERE category_id = 5 
ORDER BY price DESC;  -- Uses index efficiently
```

### Functional Indexes

```sql
-- Index on expression
CREATE INDEX idx_year ON orders((YEAR(order_date)));

-- Now this query uses the index
SELECT * FROM orders WHERE YEAR(order_date) = 2024;

-- Index on JSON field
CREATE INDEX idx_brand ON products((CAST(attributes->>'$.brand' AS CHAR(100))));
```

### Instant ADD COLUMN

```sql
-- MySQL 8.0.12+
-- Adding columns is often instant (no table rebuild)
ALTER TABLE products ADD COLUMN discount_code VARCHAR(50);

-- Check if operation will be instant
ALTER TABLE products ADD COLUMN new_col INT, ALGORITHM=INSTANT;
```

### Roles

```sql
-- Create role
CREATE ROLE 'app_developer';

-- Grant privileges to role
GRANT SELECT, INSERT, UPDATE, DELETE ON myapp.* TO 'app_developer';

-- Grant role to user
GRANT 'app_developer' TO 'john'@'localhost';

-- Set default role
SET DEFAULT ROLE 'app_developer' TO 'john'@'localhost';
```

### Other 8.0 Features

```sql
-- CHECK constraints
CREATE TABLE products (
    id INT PRIMARY KEY,
    price DECIMAL(10,2) CHECK (price >= 0),
    quantity INT CHECK (quantity >= 0)
);

-- DEFAULT expressions
CREATE TABLE orders (
    id INT PRIMARY KEY,
    order_number VARCHAR(20) DEFAULT (UUID()),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Multi-valued indexes for JSON arrays
CREATE INDEX idx_tags ON products((CAST(attributes->'$.tags' AS UNSIGNED ARRAY)));

SELECT * FROM products 
WHERE JSON_CONTAINS(attributes->'$.tags', CAST('[1, 2]' AS JSON));

-- EXPLAIN ANALYZE
EXPLAIN ANALYZE SELECT * FROM orders WHERE customer_id = 100;

-- Histogram statistics
ANALYZE TABLE products UPDATE HISTOGRAM ON price, category_id;

-- Show histograms
SELECT * FROM information_schema.COLUMN_STATISTICS 
WHERE TABLE_NAME = 'products';
```

---

## 20. Best Practices

### Schema Design

```sql
-- 1. Use appropriate data types
-- Bad:
phone VARCHAR(255)
-- Good:
phone VARCHAR(20)

-- Bad:
amount FLOAT
-- Good:
amount DECIMAL(10,2)

-- 2. Always have a PRIMARY KEY
CREATE TABLE orders (
    id INT PRIMARY KEY AUTO_INCREMENT,
    ...
);

-- 3. Use foreign keys for referential integrity
CREATE TABLE order_items (
    id INT PRIMARY KEY,
    order_id INT,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

-- 4. Use consistent naming conventions
-- Tables: plural, snake_case (users, order_items)
-- Columns: singular, snake_case (user_id, created_at)
-- Indexes: idx_table_column (idx_orders_customer_id)

-- 5. Add timestamps
CREATE TABLE articles (
    id INT PRIMARY KEY,
    title VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 6. Use UTF8MB4 for full Unicode support
CREATE DATABASE myapp 
CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

-- 7. Document your schema
-- Add comments to tables and columns
ALTER TABLE users COMMENT = 'Registered application users';
ALTER TABLE users MODIFY email VARCHAR(255) COMMENT 'User login email address';
```

### Query Writing

```sql
-- 1. Select only needed columns
SELECT id, name, email FROM users;  -- Not SELECT *

-- 2. Use parameterized queries (application code)
-- Prevents SQL injection

-- 3. Limit result sets
SELECT * FROM logs ORDER BY created_at DESC LIMIT 100;

-- 4. Use proper JOIN types
SELECT o.* FROM orders o
INNER JOIN customers c ON o.customer_id = c.id;  -- Not implicit join

-- 5. Avoid N+1 queries
-- Bad: Query users, then loop to query orders for each
-- Good: JOIN or use IN clause

-- 6. Use transactions appropriately
START TRANSACTION;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;
UPDATE accounts SET balance = balance + 100 WHERE id = 2;
COMMIT;

-- 7. Use EXPLAIN to understand query execution
EXPLAIN SELECT * FROM orders WHERE customer_id = 100;
```

### Security

```sql
-- 1. Use strong passwords
CREATE USER 'app'@'localhost' IDENTIFIED BY 'Str0ng!P@ssw0rd#2024';

-- 2. Grant minimum necessary privileges
GRANT SELECT, INSERT, UPDATE ON myapp.* TO 'app_user'@'localhost';

-- 3. Never use root for applications

-- 4. Use SSL for connections
CREATE USER 'remote'@'%' IDENTIFIED BY 'password' REQUIRE SSL;

-- 5. Validate and sanitize input in application code

-- 6. Don't store passwords in plain text
-- Use password_hash() in application

-- 7. Regular security audits
SELECT User, Host, authentication_string FROM mysql.user;
```

### Maintenance

```sql
-- 1. Regular backups
-- Schedule daily mysqldump

-- 2. Monitor for issues
-- Check slow query log regularly
-- Monitor replication lag
-- Watch disk space

-- 3. Optimize tables periodically
OPTIMIZE TABLE orders;

-- 4. Update statistics
ANALYZE TABLE orders;

-- 5. Check for table corruption
CHECK TABLE orders;

-- 6. Keep MySQL updated
-- Apply security patches

-- 7. Monitor connections
SHOW STATUS LIKE 'Threads_connected';
SHOW STATUS LIKE 'Max_used_connections';
```

### Development Workflow

```sql
-- 1. Use version control for schema changes
-- Keep all migrations in source control

-- 2. Use migrations for schema changes
-- Tools: Flyway, Liquibase, Laravel Migrations

-- 3. Test on staging before production

-- 4. Have a rollback plan

-- 5. Use transactions for multi-statement operations

-- 6. Log slow queries in development

-- 7. Profile queries before deployment
SET profiling = 1;
-- Run query
SHOW PROFILE;
```

---

## Quick Reference

### Common Commands

```sql
-- Database operations
CREATE DATABASE dbname;
DROP DATABASE dbname;
USE dbname;
SHOW DATABASES;

-- Table operations
CREATE TABLE tablename (...);
DROP TABLE tablename;
ALTER TABLE tablename ADD COLUMN col_name type;
DESCRIBE tablename;
SHOW TABLES;

-- User operations
CREATE USER 'user'@'host' IDENTIFIED BY 'password';
GRANT privileges ON db.* TO 'user'@'host';
DROP USER 'user'@'host';

-- Data operations
INSERT INTO table (...) VALUES (...);
SELECT columns FROM table WHERE condition;
UPDATE table SET column = value WHERE condition;
DELETE FROM table WHERE condition;

-- Transactions
START TRANSACTION;
COMMIT;
ROLLBACK;
```

### Useful System Commands

```sql
-- Server information
SELECT VERSION();
SELECT @@hostname;
SHOW VARIABLES LIKE 'datadir';
SHOW STATUS;

-- Session information
SELECT USER();
SELECT DATABASE();
SHOW PROCESSLIST;

-- Table information
SHOW TABLE STATUS;
SHOW CREATE TABLE tablename;
SHOW INDEX FROM tablename;
```

---

## Resources

### Official Documentation
- [MySQL Documentation](https://dev.mysql.com/doc/)
- [MySQL Reference Manual](https://dev.mysql.com/doc/refman/8.0/en/)

### Learning Resources
- MySQL Tutorial (mysqltutorial.org)
- W3Schools MySQL Tutorial
- MySQL Official Training

### Tools
- **MySQL Workbench** - Official GUI tool
- **phpMyAdmin** - Web-based administration
- **DBeaver** - Universal database tool
- **HeidiSQL** - Lightweight Windows client

---

*Last Updated: February 2026*
