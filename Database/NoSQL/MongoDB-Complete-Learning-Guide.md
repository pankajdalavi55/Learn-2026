# MongoDB Complete Learning Guide
## From Fundamentals to Advanced Concepts

---

## Table of Contents

1. [Introduction to MongoDB](#1-introduction-to-mongodb)
2. [Installation & Setup](#2-installation--setup)
3. [Database Fundamentals](#3-database-fundamentals)
4. [CRUD Operations](#4-crud-operations)
5. [Query Operators](#5-query-operators)
6. [Aggregation Framework](#6-aggregation-framework)
7. [Indexes & Performance](#7-indexes--performance)
8. [Data Modeling](#8-data-modeling)
9. [Schema Validation](#9-schema-validation)
10. [Transactions](#10-transactions)
11. [Replication](#11-replication)
12. [Sharding](#12-sharding)
13. [Security](#13-security)
14. [Backup & Recovery](#14-backup--recovery)
15. [Performance Tuning](#15-performance-tuning)
16. [MongoDB with Applications](#16-mongodb-with-applications)
17. [Atlas (Cloud MongoDB)](#17-atlas-cloud-mongodb)
18. [Best Practices](#18-best-practices)
19. [Common Patterns](#19-common-patterns)
20. [Troubleshooting](#20-troubleshooting)

---

## 1. Introduction to MongoDB

### What is MongoDB?

MongoDB is a document-oriented NoSQL database designed for scalability, flexibility, and high performance. Instead of storing data in tables with fixed schemas, MongoDB stores data as flexible JSON-like documents (BSON - Binary JSON).

### Key Features

| Feature | Description |
|---------|-------------|
| **Document Model** | Flexible, JSON-like documents |
| **Schema Flexibility** | Dynamic schemas, no migrations needed |
| **Scalability** | Horizontal scaling via sharding |
| **High Availability** | Built-in replication |
| **Rich Queries** | Powerful query language with aggregation |
| **Indexing** | Multiple index types including geospatial |
| **GridFS** | Store large files |
| **ACID Transactions** | Multi-document transactions (4.0+) |

### MongoDB vs Relational Databases

| Concept | RDBMS | MongoDB |
|---------|-------|---------|
| Database | Database | Database |
| Table | Table | Collection |
| Row | Row | Document |
| Column | Column | Field |
| Primary Key | Primary Key | _id |
| Index | Index | Index |
| Join | JOIN | $lookup / Embedding |
| Foreign Key | Foreign Key | Reference |
| Partition | Partition | Shard |

### When to Use MongoDB

**Best For:**
- Applications with rapidly changing requirements
- Content management systems
- Real-time analytics
- IoT and time-series data
- Mobile applications
- Catalog and inventory management
- User profiles and personalization

**Consider Alternatives When:**
- Heavy relational operations needed
- Complex multi-table transactions required
- Strong schema enforcement is critical
- Existing RDBMS works well

### MongoDB Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT APPLICATIONS                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  mongosh │  │  Node.js │  │  Python  │  │   Java   │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                     MONGODB DRIVERS                         │
│              (Connection Pooling, Wire Protocol)            │
└─────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                      MONGOS (Router)                        │
│                  (For Sharded Clusters)                     │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│   Shard 1    │   │   Shard 2    │   │   Shard 3    │
│  (Replica    │   │  (Replica    │   │  (Replica    │
│    Set)      │   │    Set)      │   │    Set)      │
└──────────────┘   └──────────────┘   └──────────────┘
```

### BSON Data Types

```javascript
// BSON supports rich data types
{
    // String
    name: "John Doe",
    
    // Numbers
    age: 30,                    // Int32
    salary: NumberLong(75000),  // Int64
    price: 19.99,               // Double
    amount: NumberDecimal("19.99"),  // Decimal128
    
    // Boolean
    isActive: true,
    
    // Date
    createdAt: new Date(),
    birthDate: ISODate("1990-05-15"),
    
    // ObjectId (12-byte unique identifier)
    _id: ObjectId("507f1f77bcf86cd799439011"),
    
    // Array
    tags: ["mongodb", "database", "nosql"],
    
    // Embedded Document
    address: {
        street: "123 Main St",
        city: "New York",
        zip: "10001"
    },
    
    // Null
    middleName: null,
    
    // Binary Data
    fileData: BinData(0, "base64encodeddata"),
    
    // Regular Expression
    pattern: /^[a-z]+$/i,
    
    // Timestamp (for internal use)
    ts: Timestamp(1234567890, 1),
    
    // UUID
    uuid: UUID("550e8400-e29b-41d4-a716-446655440000")
}
```

---

## 2. Installation & Setup

### Installing MongoDB on Windows

```powershell
# Using Chocolatey
choco install mongodb

# Or download from MongoDB website
# https://www.mongodb.com/try/download/community

# Start MongoDB service
net start MongoDB

# Connect using mongosh
mongosh
```

### Installing MongoDB on Linux (Ubuntu/Debian)

```bash
# Import MongoDB public GPG Key
curl -fsSL https://pgp.mongodb.com/server-7.0.asc | \
   sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg --dearmor

# Add MongoDB repository
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list

# Update and install
sudo apt update
sudo apt install -y mongodb-org

# Start MongoDB
sudo systemctl start mongod
sudo systemctl enable mongod

# Connect
mongosh
```

### Installing MongoDB on macOS

```bash
# Using Homebrew
brew tap mongodb/brew
brew install mongodb-community@7.0

# Start MongoDB
brew services start mongodb-community@7.0

# Connect
mongosh
```

### Docker Installation (Recommended for Development)

```yaml
# docker-compose.yml
version: '3.8'
services:
  mongodb:
    image: mongo:7.0
    container_name: mongodb_dev
    environment:
      MONGO_INITDB_ROOT_USERNAME: admin
      MONGO_INITDB_ROOT_PASSWORD: adminpassword
      MONGO_INITDB_DATABASE: myapp
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
      - ./init-mongo.js:/docker-entrypoint-initdb.d/init-mongo.js

  mongo-express:
    image: mongo-express
    container_name: mongo_express
    environment:
      ME_CONFIG_MONGODB_ADMINUSERNAME: admin
      ME_CONFIG_MONGODB_ADMINPASSWORD: adminpassword
      ME_CONFIG_MONGODB_URL: mongodb://admin:adminpassword@mongodb:27017/
    ports:
      - "8081:8081"
    depends_on:
      - mongodb

volumes:
  mongodb_data:
```

```bash
# Start
docker-compose up -d

# Connect
docker exec -it mongodb_dev mongosh -u admin -p adminpassword
```

### MongoDB Configuration (mongod.conf)

```yaml
# mongod.conf

# Storage
storage:
  dbPath: /var/lib/mongodb
  journal:
    enabled: true
  wiredTiger:
    engineConfig:
      cacheSizeGB: 2

# Logging
systemLog:
  destination: file
  logAppend: true
  path: /var/log/mongodb/mongod.log

# Network
net:
  port: 27017
  bindIp: 127.0.0.1

# Security
security:
  authorization: enabled

# Replication (for replica sets)
replication:
  replSetName: rs0

# Sharding (for sharded clusters)
# sharding:
#   clusterRole: shardsvr

# Operation Profiling
operationProfiling:
  mode: slowOp
  slowOpThresholdMs: 100
```

### mongosh Commands Reference

```javascript
// Connect to MongoDB
mongosh                                    // Local connection
mongosh "mongodb://localhost:27017"        // Connection string
mongosh --host localhost --port 27017      // Explicit host/port
mongosh -u admin -p password --authenticationDatabase admin

// Database commands
show dbs                    // List databases
use myapp                   // Switch/create database
db                          // Show current database
db.dropDatabase()           // Drop current database

// Collection commands
show collections            // List collections
db.createCollection("users") // Create collection
db.users.drop()             // Drop collection

// Help
help                        // General help
db.help()                   // Database help
db.collection.help()        // Collection help

// Shell utilities
cls                         // Clear screen
exit                        // Exit shell
load("script.js")           // Load JavaScript file

// Pretty print
db.users.find().pretty()
```

---

## 3. Database Fundamentals

### Creating Databases and Collections

```javascript
// Databases are created implicitly when you first store data
use myapp

// Create collection explicitly
db.createCollection("users")

// Create collection with options
db.createCollection("logs", {
    capped: true,
    size: 10485760,      // 10MB
    max: 5000            // Max 5000 documents
})

// Create collection with validation
db.createCollection("products", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name", "price"],
            properties: {
                name: { bsonType: "string" },
                price: { bsonType: "number", minimum: 0 }
            }
        }
    }
})

// List all collections
db.getCollectionNames()
show collections

// Collection info
db.users.stats()
db.users.storageSize()
db.users.totalSize()
db.users.count()
```

### Document Structure

```javascript
// Basic document
{
    _id: ObjectId("507f1f77bcf86cd799439011"),  // Auto-generated if not provided
    name: "John Doe",
    email: "john@example.com",
    age: 30,
    tags: ["developer", "mongodb"],
    address: {
        street: "123 Main St",
        city: "New York",
        state: "NY",
        zip: "10001"
    },
    createdAt: new Date(),
    updatedAt: new Date()
}

// Nested documents (embedded)
{
    _id: ObjectId(),
    title: "MongoDB Guide",
    author: {
        name: "Jane Smith",
        email: "jane@example.com",
        bio: "Database expert"
    },
    chapters: [
        { number: 1, title: "Introduction", pages: 20 },
        { number: 2, title: "Installation", pages: 15 },
        { number: 3, title: "CRUD Operations", pages: 30 }
    ]
}

// Document with reference (normalized)
// users collection
{
    _id: ObjectId("user123"),
    name: "John Doe",
    email: "john@example.com"
}

// orders collection
{
    _id: ObjectId("order123"),
    userId: ObjectId("user123"),  // Reference to users
    items: [
        { productId: ObjectId("prod1"), quantity: 2 },
        { productId: ObjectId("prod2"), quantity: 1 }
    ],
    total: 99.99
}
```

### ObjectId Structure

```javascript
// ObjectId is a 12-byte identifier
// 4 bytes: timestamp (seconds since Unix epoch)
// 5 bytes: random value (unique to machine and process)
// 3 bytes: incrementing counter

// Create ObjectId
let id = new ObjectId()
let specificId = ObjectId("507f1f77bcf86cd799439011")

// Extract timestamp
id.getTimestamp()  // Returns Date

// ObjectId methods
id.toString()      // "507f1f77bcf86cd799439011"
id.valueOf()       // Same as toString()

// Generate ObjectId from timestamp
ObjectId.createFromTime(Date.now() / 1000)
```

---

## 4. CRUD Operations

### INSERT - Create Documents

```javascript
// Insert one document
db.users.insertOne({
    name: "John Doe",
    email: "john@example.com",
    age: 30,
    createdAt: new Date()
})

// Insert with specific _id
db.users.insertOne({
    _id: "user_123",
    name: "Jane Doe",
    email: "jane@example.com"
})

// Insert multiple documents
db.users.insertMany([
    { name: "Alice", email: "alice@example.com", age: 25 },
    { name: "Bob", email: "bob@example.com", age: 35 },
    { name: "Charlie", email: "charlie@example.com", age: 28 }
])

// Insert with ordered: false (continue on error)
db.users.insertMany(
    [
        { _id: 1, name: "User 1" },
        { _id: 1, name: "Duplicate" },  // Error, but continue
        { _id: 2, name: "User 2" }
    ],
    { ordered: false }
)

// Bulk write for mixed operations
db.users.bulkWrite([
    { insertOne: { document: { name: "New User" } } },
    { updateOne: { filter: { name: "Alice" }, update: { $set: { age: 26 } } } },
    { deleteOne: { filter: { name: "Bob" } } }
])
```

### FIND - Read Documents

```javascript
// Find all documents
db.users.find()

// Find with condition
db.users.find({ age: 30 })

// Find with multiple conditions (implicit AND)
db.users.find({ age: 30, isActive: true })

// Find one document
db.users.findOne({ email: "john@example.com" })

// Find by _id
db.users.findOne({ _id: ObjectId("507f1f77bcf86cd799439011") })

// Projection (select specific fields)
db.users.find({}, { name: 1, email: 1 })           // Include only these
db.users.find({}, { password: 0, __v: 0 })         // Exclude these
db.users.find({}, { name: 1, _id: 0 })             // Include name, exclude _id

// Projection with nested fields
db.users.find({}, { "address.city": 1, name: 1 })

// Sorting
db.users.find().sort({ name: 1 })                   // Ascending
db.users.find().sort({ age: -1, name: 1 })         // Descending age, ascending name

// Limit and skip (pagination)
db.users.find().limit(10)                          // First 10
db.users.find().skip(20).limit(10)                 // Skip 20, get 10

// Count documents
db.users.countDocuments({ isActive: true })
db.users.estimatedDocumentCount()                  // Faster, uses metadata

// Distinct values
db.users.distinct("city")
db.users.distinct("city", { isActive: true })

// Cursor methods
let cursor = db.users.find()
cursor.hasNext()      // Check if more documents
cursor.next()         // Get next document
cursor.toArray()      // Convert to array
cursor.forEach(doc => print(doc.name))
```

### UPDATE - Modify Documents

```javascript
// Update one document
db.users.updateOne(
    { email: "john@example.com" },      // Filter
    { $set: { age: 31, updatedAt: new Date() } }  // Update
)

// Update multiple documents
db.users.updateMany(
    { isActive: false },
    { $set: { archived: true } }
)

// Replace entire document (except _id)
db.users.replaceOne(
    { email: "john@example.com" },
    { 
        name: "John Smith",
        email: "john@example.com",
        age: 32
    }
)

// Upsert (insert if not exists)
db.users.updateOne(
    { email: "new@example.com" },
    { $set: { name: "New User", email: "new@example.com" } },
    { upsert: true }
)

// Update with array filters
db.orders.updateOne(
    { _id: ObjectId("order123") },
    { $set: { "items.$[elem].price": 25 } },
    { arrayFilters: [{ "elem.productId": ObjectId("prod1") }] }
)

// findOneAndUpdate (returns document)
db.users.findOneAndUpdate(
    { email: "john@example.com" },
    { $inc: { loginCount: 1 } },
    { returnDocument: "after" }         // Return updated document
)

// Update operators
db.users.updateOne(
    { _id: ObjectId("...") },
    {
        $set: { name: "New Name" },           // Set field value
        $unset: { tempField: "" },            // Remove field
        $inc: { age: 1 },                     // Increment
        $mul: { salary: 1.1 },                // Multiply
        $rename: { oldName: "newName" },      // Rename field
        $min: { lowestScore: 50 },            // Update if less
        $max: { highestScore: 100 },          // Update if greater
        $currentDate: { lastModified: true }  // Set to current date
    }
)

// Array update operators
db.users.updateOne(
    { _id: ObjectId("...") },
    {
        $push: { tags: "new-tag" },           // Add to array
        $addToSet: { tags: "unique-tag" },    // Add if not exists
        $pop: { tags: 1 },                    // Remove last (-1 for first)
        $pull: { tags: "remove-this" },       // Remove matching
        $pullAll: { tags: ["a", "b"] }        // Remove all matching
    }
)

// Push with modifiers
db.users.updateOne(
    { _id: ObjectId("...") },
    {
        $push: {
            scores: {
                $each: [90, 85, 95],
                $sort: -1,                    // Sort descending
                $slice: 10                    // Keep only top 10
            }
        }
    }
)

// Update nested array element
db.users.updateOne(
    { "addresses.type": "home" },
    { $set: { "addresses.$.city": "New York" } }  // $ = first match
)
```

### DELETE - Remove Documents

```javascript
// Delete one document
db.users.deleteOne({ email: "john@example.com" })

// Delete multiple documents
db.users.deleteMany({ isActive: false })

// Delete all documents
db.users.deleteMany({})

// findOneAndDelete (returns deleted document)
db.users.findOneAndDelete({ email: "john@example.com" })

// Drop collection (faster than deleteMany for all docs)
db.users.drop()
```

---

## 5. Query Operators

### Comparison Operators

```javascript
// $eq - Equal
db.users.find({ age: { $eq: 30 } })
db.users.find({ age: 30 })  // Shorthand

// $ne - Not equal
db.users.find({ status: { $ne: "deleted" } })

// $gt, $gte - Greater than, greater than or equal
db.users.find({ age: { $gt: 25 } })
db.users.find({ age: { $gte: 25 } })

// $lt, $lte - Less than, less than or equal
db.users.find({ age: { $lt: 30 } })
db.users.find({ age: { $lte: 30 } })

// $in - Matches any value in array
db.users.find({ status: { $in: ["active", "pending"] } })

// $nin - Matches none of the values
db.users.find({ status: { $nin: ["deleted", "archived"] } })

// Combined range query
db.users.find({ age: { $gte: 25, $lte: 35 } })
```

### Logical Operators

```javascript
// $and - All conditions must match
db.users.find({
    $and: [
        { age: { $gte: 25 } },
        { age: { $lte: 35 } }
    ]
})

// Implicit $and (shorthand)
db.users.find({ age: { $gte: 25 }, status: "active" })

// $or - At least one must match
db.users.find({
    $or: [
        { age: { $lt: 25 } },
        { age: { $gt: 60 } }
    ]
})

// $nor - None must match
db.users.find({
    $nor: [
        { status: "deleted" },
        { status: "archived" }
    ]
})

// $not - Negates expression
db.users.find({
    age: { $not: { $gt: 30 } }
})

// Combined logical operators
db.users.find({
    $and: [
        { $or: [{ city: "NYC" }, { city: "LA" }] },
        { status: "active" }
    ]
})
```

### Element Operators

```javascript
// $exists - Field exists
db.users.find({ phone: { $exists: true } })
db.users.find({ phone: { $exists: false } })

// $type - Field is of specific type
db.users.find({ age: { $type: "int" } })
db.users.find({ age: { $type: "number" } })  // Any numeric type
db.users.find({ data: { $type: "array" } })

// BSON type numbers or aliases
// "double" (1), "string" (2), "object" (3), "array" (4),
// "binData" (5), "objectId" (7), "bool" (8), "date" (9),
// "null" (10), "regex" (11), "int" (16), "long" (18), "decimal" (19)
```

### Array Operators

```javascript
// $all - Array contains all elements
db.products.find({
    tags: { $all: ["electronics", "sale"] }
})

// $elemMatch - Array element matches all conditions
db.orders.find({
    items: {
        $elemMatch: {
            productId: ObjectId("..."),
            quantity: { $gte: 2 }
        }
    }
})

// $size - Array has exact size
db.users.find({ tags: { $size: 3 } })

// Array by index
db.users.find({ "scores.0": { $gt: 90 } })  // First element > 90
```

### Evaluation Operators

```javascript
// $regex - Regular expression match
db.users.find({ name: { $regex: /^john/i } })
db.users.find({ email: { $regex: "@gmail\\.com$", $options: "i" } })

// $expr - Use aggregation expressions
db.products.find({
    $expr: { $gt: ["$price", "$cost"] }  // price > cost
})

db.orders.find({
    $expr: {
        $gt: [
            { $multiply: ["$quantity", "$price"] },
            100
        ]
    }
})

// $mod - Modulo
db.items.find({ quantity: { $mod: [4, 0] } })  // quantity % 4 == 0

// $text - Text search (requires text index)
db.articles.find({ $text: { $search: "mongodb tutorial" } })
db.articles.find({ $text: { $search: "\"exact phrase\"" } })
db.articles.find({ $text: { $search: "mongodb -mysql" } })  // Exclude mysql

// $where - JavaScript expression (avoid if possible, slow)
db.users.find({
    $where: function() {
        return this.name.length > 10
    }
})

// $jsonSchema - Validate document structure
db.users.find({
    $jsonSchema: {
        required: ["email", "name"],
        properties: {
            email: { bsonType: "string" },
            age: { bsonType: "int", minimum: 0 }
        }
    }
})
```

### Geospatial Operators

```javascript
// Create 2dsphere index first
db.places.createIndex({ location: "2dsphere" })

// $near - Find nearest
db.places.find({
    location: {
        $near: {
            $geometry: {
                type: "Point",
                coordinates: [-73.97, 40.77]
            },
            $maxDistance: 5000  // meters
        }
    }
})

// $geoWithin - Within shape
db.places.find({
    location: {
        $geoWithin: {
            $geometry: {
                type: "Polygon",
                coordinates: [[
                    [-73.99, 40.75],
                    [-73.98, 40.75],
                    [-73.98, 40.76],
                    [-73.99, 40.76],
                    [-73.99, 40.75]
                ]]
            }
        }
    }
})

// $geoIntersects - Intersects with geometry
db.areas.find({
    region: {
        $geoIntersects: {
            $geometry: {
                type: "Point",
                coordinates: [-73.97, 40.77]
            }
        }
    }
})
```

### Projection Operators

```javascript
// $slice - Limit array elements
db.posts.find({}, { comments: { $slice: 5 } })         // First 5
db.posts.find({}, { comments: { $slice: -5 } })        // Last 5
db.posts.find({}, { comments: { $slice: [10, 5] } })   // Skip 10, get 5

// $elemMatch - Project matching array element
db.orders.find(
    {},
    { items: { $elemMatch: { productId: ObjectId("...") } } }
)

// $ - Project first matching array element
db.orders.find(
    { "items.productId": ObjectId("...") },
    { "items.$": 1 }
)

// $meta - Include text search score
db.articles.find(
    { $text: { $search: "mongodb" } },
    { score: { $meta: "textScore" } }
).sort({ score: { $meta: "textScore" } })
```

---

## 6. Aggregation Framework

### Pipeline Basics

```javascript
// Basic aggregation structure
db.collection.aggregate([
    { $stage1: { ... } },
    { $stage2: { ... } },
    { $stage3: { ... } }
])

// Example: Sales report
db.orders.aggregate([
    { $match: { status: "completed" } },
    { $group: {
        _id: "$customerId",
        totalSpent: { $sum: "$total" },
        orderCount: { $sum: 1 }
    }},
    { $sort: { totalSpent: -1 } },
    { $limit: 10 }
])
```

### Common Pipeline Stages

#### $match - Filter Documents

```javascript
// Filter at the start for better performance
db.orders.aggregate([
    { $match: {
        status: "completed",
        orderDate: { $gte: ISODate("2024-01-01") }
    }}
])
```

#### $project - Reshape Documents

```javascript
db.users.aggregate([
    { $project: {
        fullName: { $concat: ["$firstName", " ", "$lastName"] },
        email: 1,
        age: 1,
        isAdult: { $gte: ["$age", 18] },
        _id: 0
    }}
])

// Computed fields
db.orders.aggregate([
    { $project: {
        _id: 1,
        items: 1,
        subtotal: { $sum: "$items.price" },
        tax: { $multiply: [{ $sum: "$items.price" }, 0.1] },
        total: { 
            $multiply: [
                { $sum: "$items.price" },
                1.1
            ]
        }
    }}
])
```

#### $group - Group and Aggregate

```javascript
// Basic grouping
db.orders.aggregate([
    { $group: {
        _id: "$customerId",
        totalAmount: { $sum: "$total" },
        averageOrder: { $avg: "$total" },
        maxOrder: { $max: "$total" },
        minOrder: { $min: "$total" },
        orderCount: { $sum: 1 },
        orders: { $push: "$_id" },
        firstOrder: { $first: "$orderDate" },
        lastOrder: { $last: "$orderDate" }
    }}
])

// Group by multiple fields
db.orders.aggregate([
    { $group: {
        _id: {
            year: { $year: "$orderDate" },
            month: { $month: "$orderDate" }
        },
        revenue: { $sum: "$total" },
        orders: { $sum: 1 }
    }},
    { $sort: { "_id.year": 1, "_id.month": 1 } }
])

// Group all (no _id)
db.products.aggregate([
    { $group: {
        _id: null,
        totalProducts: { $sum: 1 },
        avgPrice: { $avg: "$price" },
        priceRange: { $push: "$price" }
    }}
])

// $addToSet (unique values only)
db.orders.aggregate([
    { $group: {
        _id: "$customerId",
        uniqueProducts: { $addToSet: "$productId" }
    }}
])
```

#### $sort, $limit, $skip

```javascript
db.users.aggregate([
    { $sort: { score: -1, name: 1 } },
    { $skip: 20 },
    { $limit: 10 }
])
```

#### $unwind - Deconstruct Arrays

```javascript
// Expand array into separate documents
db.orders.aggregate([
    { $unwind: "$items" },
    { $group: {
        _id: "$items.productId",
        totalSold: { $sum: "$items.quantity" },
        revenue: { $sum: { $multiply: ["$items.quantity", "$items.price"] } }
    }}
])

// Preserve empty arrays
db.orders.aggregate([
    { $unwind: {
        path: "$items",
        preserveNullAndEmptyArrays: true
    }}
])

// Include array index
db.orders.aggregate([
    { $unwind: {
        path: "$items",
        includeArrayIndex: "itemIndex"
    }}
])
```

#### $lookup - Join Collections

```javascript
// Basic lookup (left outer join)
db.orders.aggregate([
    { $lookup: {
        from: "users",
        localField: "customerId",
        foreignField: "_id",
        as: "customer"
    }},
    { $unwind: "$customer" }
])

// Lookup with pipeline (correlated subquery)
db.orders.aggregate([
    { $lookup: {
        from: "products",
        let: { orderItems: "$items" },
        pipeline: [
            { $match: {
                $expr: { $in: ["$_id", "$$orderItems.productId"] }
            }},
            { $project: { name: 1, price: 1 } }
        ],
        as: "productDetails"
    }}
])

// Multiple lookups for complex joins
db.orders.aggregate([
    { $lookup: {
        from: "users",
        localField: "customerId",
        foreignField: "_id",
        as: "customer"
    }},
    { $lookup: {
        from: "products",
        localField: "items.productId",
        foreignField: "_id",
        as: "products"
    }}
])
```

#### $addFields / $set - Add New Fields

```javascript
db.users.aggregate([
    { $addFields: {
        fullName: { $concat: ["$firstName", " ", "$lastName"] },
        isAdult: { $gte: ["$age", 18] },
        ageCategory: {
            $switch: {
                branches: [
                    { case: { $lt: ["$age", 18] }, then: "minor" },
                    { case: { $lt: ["$age", 65] }, then: "adult" },
                ],
                default: "senior"
            }
        }
    }}
])

// $set is alias for $addFields (MongoDB 4.2+)
db.users.aggregate([
    { $set: { updatedAt: "$$NOW" } }
])
```

#### $facet - Multiple Pipelines

```javascript
// Run multiple aggregations in parallel
db.products.aggregate([
    { $facet: {
        "categoryCounts": [
            { $group: { _id: "$category", count: { $sum: 1 } } }
        ],
        "priceStats": [
            { $group: {
                _id: null,
                avgPrice: { $avg: "$price" },
                maxPrice: { $max: "$price" },
                minPrice: { $min: "$price" }
            }}
        ],
        "topProducts": [
            { $sort: { sales: -1 } },
            { $limit: 5 },
            { $project: { name: 1, sales: 1 } }
        ]
    }}
])
```

#### $bucket / $bucketAuto - Histogram

```javascript
// Manual bucket boundaries
db.users.aggregate([
    { $bucket: {
        groupBy: "$age",
        boundaries: [0, 18, 30, 50, 70, 100],
        default: "Unknown",
        output: {
            count: { $sum: 1 },
            users: { $push: "$name" }
        }
    }}
])

// Auto-bucket (MongoDB determines boundaries)
db.products.aggregate([
    { $bucketAuto: {
        groupBy: "$price",
        buckets: 5,
        output: {
            count: { $sum: 1 },
            avgPrice: { $avg: "$price" }
        }
    }}
])
```

#### $out / $merge - Write Results

```javascript
// Write to new collection (replaces)
db.orders.aggregate([
    { $group: { _id: "$customerId", total: { $sum: "$amount" } } },
    { $out: "customer_totals" }
])

// Merge into existing collection (MongoDB 4.2+)
db.orders.aggregate([
    { $group: { _id: "$customerId", total: { $sum: "$amount" } } },
    { $merge: {
        into: "customer_totals",
        on: "_id",
        whenMatched: "merge",
        whenNotMatched: "insert"
    }}
])
```

### Aggregation Expressions

```javascript
// Arithmetic
{ $add: [expr1, expr2] }
{ $subtract: [expr1, expr2] }
{ $multiply: [expr1, expr2] }
{ $divide: [expr1, expr2] }
{ $mod: [expr1, expr2] }
{ $abs: expr }
{ $ceil: expr }
{ $floor: expr }
{ $round: [expr, places] }

// String
{ $concat: [str1, str2, ...] }
{ $substr: [string, start, length] }
{ $toLower: expr }
{ $toUpper: expr }
{ $trim: { input: string } }
{ $split: [string, delimiter] }
{ $regexMatch: { input: string, regex: /pattern/ } }

// Date
{ $year: dateExpr }
{ $month: dateExpr }
{ $dayOfMonth: dateExpr }
{ $hour: dateExpr }
{ $minute: dateExpr }
{ $dateToString: { format: "%Y-%m-%d", date: dateExpr } }
{ $dateFromString: { dateString: "2024-03-15" } }
{ $dateDiff: { startDate: date1, endDate: date2, unit: "day" } }

// Conditional
{ $cond: { if: boolExpr, then: trueExpr, else: falseExpr } }
{ $cond: [boolExpr, trueExpr, falseExpr] }  // Array syntax
{ $ifNull: [expr, replacement] }
{ $switch: {
    branches: [
        { case: expr1, then: result1 },
        { case: expr2, then: result2 }
    ],
    default: defaultResult
}}

// Array
{ $size: arrayExpr }
{ $arrayElemAt: [array, index] }
{ $first: arrayExpr }
{ $last: arrayExpr }
{ $slice: [array, n] }
{ $filter: { input: array, as: "item", cond: { $gt: ["$$item.price", 100] } } }
{ $map: { input: array, as: "item", in: { $multiply: ["$$item.price", 1.1] } } }
{ $reduce: { input: array, initialValue: 0, in: { $add: ["$$value", "$$this"] } } }
{ $in: [expr, array] }

// Type conversion
{ $toString: expr }
{ $toInt: expr }
{ $toDouble: expr }
{ $toDate: expr }
{ $toObjectId: expr }
{ $type: expr }  // Returns type name

// Comparison
{ $eq: [expr1, expr2] }
{ $ne: [expr1, expr2] }
{ $gt: [expr1, expr2] }
{ $gte: [expr1, expr2] }
{ $lt: [expr1, expr2] }
{ $lte: [expr1, expr2] }
{ $cmp: [expr1, expr2] }  // Returns -1, 0, or 1
```

### Real-World Aggregation Examples

```javascript
// E-commerce: Top customers with order details
db.orders.aggregate([
    { $match: { status: "completed", orderDate: { $gte: ISODate("2024-01-01") } } },
    { $lookup: {
        from: "users",
        localField: "customerId",
        foreignField: "_id",
        as: "customer"
    }},
    { $unwind: "$customer" },
    { $group: {
        _id: "$customerId",
        customerName: { $first: "$customer.name" },
        totalSpent: { $sum: "$total" },
        orderCount: { $sum: 1 },
        avgOrderValue: { $avg: "$total" },
        lastOrder: { $max: "$orderDate" }
    }},
    { $sort: { totalSpent: -1 } },
    { $limit: 10 },
    { $project: {
        _id: 0,
        customerId: "$_id",
        customerName: 1,
        totalSpent: { $round: ["$totalSpent", 2] },
        orderCount: 1,
        avgOrderValue: { $round: ["$avgOrderValue", 2] },
        lastOrder: 1
    }}
])

// Time-series: Daily active users
db.events.aggregate([
    { $match: { eventType: "login", timestamp: { $gte: ISODate("2024-01-01") } } },
    { $group: {
        _id: {
            date: { $dateToString: { format: "%Y-%m-%d", date: "$timestamp" } },
            userId: "$userId"
        }
    }},
    { $group: {
        _id: "$_id.date",
        activeUsers: { $sum: 1 }
    }},
    { $sort: { _id: 1 } }
])

// Product analytics: Category performance
db.products.aggregate([
    { $lookup: {
        from: "orders",
        let: { productId: "$_id" },
        pipeline: [
            { $unwind: "$items" },
            { $match: { $expr: { $eq: ["$items.productId", "$$productId"] } } },
            { $group: {
                _id: null,
                totalSold: { $sum: "$items.quantity" },
                revenue: { $sum: { $multiply: ["$items.quantity", "$items.price"] } }
            }}
        ],
        as: "sales"
    }},
    { $unwind: { path: "$sales", preserveNullAndEmptyArrays: true } },
    { $group: {
        _id: "$category",
        productCount: { $sum: 1 },
        totalSold: { $sum: { $ifNull: ["$sales.totalSold", 0] } },
        revenue: { $sum: { $ifNull: ["$sales.revenue", 0] } },
        avgPrice: { $avg: "$price" }
    }},
    { $sort: { revenue: -1 } }
])
```

---

## 7. Indexes & Performance

### Index Types

```javascript
// Single field index
db.users.createIndex({ email: 1 })           // Ascending
db.users.createIndex({ createdAt: -1 })      // Descending

// Compound index
db.orders.createIndex({ customerId: 1, orderDate: -1 })

// Unique index
db.users.createIndex({ email: 1 }, { unique: true })

// Sparse index (only index documents with the field)
db.users.createIndex({ phone: 1 }, { sparse: true })

// Partial index (only index matching documents)
db.orders.createIndex(
    { orderDate: 1 },
    { partialFilterExpression: { status: "pending" } }
)

// TTL index (auto-delete documents)
db.sessions.createIndex(
    { createdAt: 1 },
    { expireAfterSeconds: 3600 }  // Delete after 1 hour
)

// Text index (for full-text search)
db.articles.createIndex({ title: "text", content: "text" })
db.articles.createIndex(
    { title: "text", content: "text", tags: "text" },
    { weights: { title: 10, content: 5, tags: 1 } }
)

// Geospatial 2dsphere index
db.places.createIndex({ location: "2dsphere" })

// Geospatial 2d index (legacy, flat geometry)
db.legacy.createIndex({ coords: "2d" })

// Hashed index (for sharding)
db.users.createIndex({ email: "hashed" })

// Wildcard index (MongoDB 4.2+)
db.products.createIndex({ "attributes.$**": 1 })
db.logs.createIndex({ "$**": 1 })  // All fields

// Compound wildcard (MongoDB 7.0+)
db.products.createIndex({ category: 1, "specs.$**": 1 })
```

### Managing Indexes

```javascript
// List all indexes
db.users.getIndexes()

// Index usage statistics
db.users.aggregate([{ $indexStats: {} }])

// Index size
db.users.stats().indexSizes

// Drop index
db.users.dropIndex("email_1")
db.users.dropIndex({ email: 1 })

// Drop all indexes (except _id)
db.users.dropIndexes()

// Rebuild indexes
db.users.reIndex()

// Create index in background (non-blocking)
db.users.createIndex({ email: 1 }, { background: true })

// Hide index (test impact of removing)
db.users.hideIndex("email_1")
db.users.unhideIndex("email_1")
```

### Query Optimization

```javascript
// Use explain() to analyze queries
db.users.find({ email: "test@example.com" }).explain()
db.users.find({ email: "test@example.com" }).explain("executionStats")
db.users.find({ email: "test@example.com" }).explain("allPlansExecution")

// Key explain fields:
// - queryPlanner.winningPlan: Chosen execution plan
// - executionStats.totalDocsExamined: Documents scanned
// - executionStats.totalKeysExamined: Index keys scanned
// - executionStats.executionTimeMillis: Query time

// Ideal: totalDocsExamined ≈ documents returned
// Using index: "stage": "IXSCAN" or "FETCH + IXSCAN"
// Collection scan: "stage": "COLLSCAN" (usually bad)

// Aggregation explain
db.orders.explain("executionStats").aggregate([
    { $match: { status: "pending" } },
    { $group: { _id: "$customerId", total: { $sum: "$amount" } } }
])
```

### Index Strategy Best Practices

```javascript
// 1. Create indexes for common queries
db.orders.createIndex({ customerId: 1, status: 1 })

// 2. Order matters in compound indexes
// Supports: { customerId: 1 }, { customerId: 1, status: 1 }
// Does NOT support: { status: 1 } alone

// 3. ESR Rule (Equality, Sort, Range)
// Put equality matches first, then sort fields, then range queries
db.orders.createIndex({
    status: 1,       // Equality
    orderDate: -1,   // Sort
    amount: 1        // Range
})

// 4. Covered queries (index-only, no document fetch)
db.users.createIndex({ email: 1, name: 1 })
db.users.find(
    { email: "test@example.com" },
    { email: 1, name: 1, _id: 0 }  // All fields in index
)

// 5. Use partial indexes for filtered queries
db.orders.createIndex(
    { orderDate: 1 },
    { partialFilterExpression: { status: { $ne: "archived" } } }
)

// 6. Monitor slow queries
db.setProfilingLevel(1, { slowms: 100 })
db.system.profile.find().sort({ ts: -1 }).limit(10)
```

---

## 8. Data Modeling

### Embedding vs Referencing

```javascript
// EMBEDDING (Denormalized)
// Good for: 1:1, 1:few relationships, data accessed together
{
    _id: ObjectId("..."),
    name: "John Doe",
    email: "john@example.com",
    address: {                    // Embedded document
        street: "123 Main St",
        city: "New York",
        zip: "10001"
    },
    orders: [                     // Embedded array
        { productId: ObjectId(), quantity: 2, price: 29.99 },
        { productId: ObjectId(), quantity: 1, price: 49.99 }
    ]
}

// REFERENCING (Normalized)
// Good for: 1:many, many:many, large/unbounded arrays, independent entities

// users collection
{
    _id: ObjectId("user123"),
    name: "John Doe",
    email: "john@example.com"
}

// orders collection
{
    _id: ObjectId("order123"),
    userId: ObjectId("user123"),  // Reference
    items: [
        { productId: ObjectId("prod1"), quantity: 2 }
    ]
}

// products collection
{
    _id: ObjectId("prod1"),
    name: "Widget",
    price: 29.99
}
```

### Common Patterns

#### One-to-One

```javascript
// Embedded (preferred for tightly coupled data)
{
    _id: ObjectId("user123"),
    name: "John Doe",
    profile: {
        bio: "Developer",
        avatar: "url...",
        settings: { theme: "dark" }
    }
}

// Referenced (for large or independent data)
// users collection
{ _id: ObjectId("user123"), name: "John Doe", profileId: ObjectId("profile123") }

// profiles collection
{ _id: ObjectId("profile123"), bio: "Developer", avatar: "url..." }
```

#### One-to-Few

```javascript
// Embed the "few" side
{
    _id: ObjectId("user123"),
    name: "John Doe",
    addresses: [
        { type: "home", street: "123 Main St", city: "NYC" },
        { type: "work", street: "456 Office Blvd", city: "NYC" }
    ]
}
```

#### One-to-Many

```javascript
// Reference from "many" side
// blog posts
{
    _id: ObjectId("post123"),
    title: "MongoDB Guide",
    authorId: ObjectId("user123"),  // Reference to author
    content: "..."
}

// users
{
    _id: ObjectId("user123"),
    name: "John Doe"
}

// Query with $lookup
db.posts.aggregate([
    { $lookup: {
        from: "users",
        localField: "authorId",
        foreignField: "_id",
        as: "author"
    }}
])
```

#### One-to-Squillions (Unbounded)

```javascript
// Reference from parent
// host document
{
    _id: ObjectId("host123"),
    hostname: "server1.example.com",
    ip: "192.168.1.100"
}

// millions of log entries
{
    _id: ObjectId("log123"),
    hostId: ObjectId("host123"),  // Reference to host
    message: "Error occurred",
    timestamp: ISODate("2024-03-15T10:30:00Z")
}

// NEVER embed unbounded arrays!
```

#### Many-to-Many

```javascript
// Two-way referencing
// authors collection
{
    _id: ObjectId("author1"),
    name: "Author One",
    bookIds: [ObjectId("book1"), ObjectId("book2")]
}

// books collection
{
    _id: ObjectId("book1"),
    title: "Book One",
    authorIds: [ObjectId("author1"), ObjectId("author2")]
}

// Or junction collection (for metadata)
// book_authors collection
{
    bookId: ObjectId("book1"),
    authorId: ObjectId("author1"),
    role: "primary",
    royaltyPercent: 60
}
```

### Tree Structures

```javascript
// Parent Reference
{
    _id: "MongoDB",
    parent: "Databases",
    name: "MongoDB"
}

// Child Reference
{
    _id: "Databases",
    children: ["MongoDB", "PostgreSQL", "MySQL"],
    name: "Databases"
}

// Materialized Path
{
    _id: "MongoDB",
    path: ",Databases,NoSQL,MongoDB,",
    name: "MongoDB"
}
// Query all ancestors: /,Databases,/
// Query all descendants: /,Databases,NoSQL,MongoDB,/

// Nested Sets
{
    _id: "MongoDB",
    name: "MongoDB",
    left: 3,
    right: 4
}
// All descendants: left > parent.left AND right < parent.right

// Array of Ancestors
{
    _id: "MongoDB",
    name: "MongoDB",
    ancestors: [
        { _id: "Databases", name: "Databases" },
        { _id: "NoSQL", name: "NoSQL" }
    ]
}
```

### Polymorphic Pattern

```javascript
// Different types in same collection
// Products with varying attributes
{
    _id: ObjectId(),
    type: "book",
    name: "MongoDB Guide",
    price: 49.99,
    // Book-specific fields
    author: "John Smith",
    isbn: "978-0-xxx",
    pages: 350
}

{
    _id: ObjectId(),
    type: "electronics",
    name: "Laptop",
    price: 999.99,
    // Electronics-specific fields
    brand: "Dell",
    specs: {
        ram: "16GB",
        storage: "512GB SSD"
    },
    warranty: "2 years"
}

// Query by type
db.products.find({ type: "book" })
db.products.find({ type: "electronics", "specs.ram": "16GB" })
```

### Bucket Pattern (Time-Series)

```javascript
// Instead of one document per measurement
// Bucket multiple measurements together
{
    _id: ObjectId(),
    sensorId: "sensor_001",
    date: ISODate("2024-03-15"),
    // Array of hourly readings
    measurements: [
        { hour: 0, temp: 22.5, humidity: 45 },
        { hour: 1, temp: 22.3, humidity: 46 },
        { hour: 2, temp: 22.1, humidity: 47 },
        // ... up to hour 23
    ],
    count: 24,
    avgTemp: 22.5,
    avgHumidity: 46
}

// Benefits:
// - Fewer documents
// - Pre-aggregated statistics
// - Better index efficiency
```

### Computed Pattern

```javascript
// Store computed/derived data
{
    _id: ObjectId("movie123"),
    title: "Great Movie",
    ratings: [5, 4, 5, 3, 5, 4, 4, 5],
    // Pre-computed values
    ratingCount: 8,
    ratingSum: 35,
    ratingAvg: 4.375
}

// Update with $inc for efficiency
db.movies.updateOne(
    { _id: ObjectId("movie123") },
    {
        $push: { ratings: 5 },
        $inc: { ratingCount: 1, ratingSum: 5 }
    }
)

// Periodic recalculation if needed
db.movies.updateOne(
    { _id: ObjectId("movie123") },
    [{
        $set: {
            ratingAvg: { $divide: ["$ratingSum", "$ratingCount"] }
        }
    }]
)
```

---

## 9. Schema Validation

### Creating Validation Rules

```javascript
// Create collection with validation
db.createCollection("users", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            title: "User Validation",
            required: ["name", "email", "age"],
            properties: {
                name: {
                    bsonType: "string",
                    description: "must be a string and is required",
                    minLength: 2,
                    maxLength: 100
                },
                email: {
                    bsonType: "string",
                    pattern: "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
                    description: "must be a valid email"
                },
                age: {
                    bsonType: "int",
                    minimum: 0,
                    maximum: 150,
                    description: "must be an integer between 0 and 150"
                },
                status: {
                    enum: ["active", "inactive", "pending"],
                    description: "must be one of the allowed values"
                },
                address: {
                    bsonType: "object",
                    required: ["city"],
                    properties: {
                        street: { bsonType: "string" },
                        city: { bsonType: "string" },
                        zip: { bsonType: "string" }
                    }
                },
                tags: {
                    bsonType: "array",
                    items: { bsonType: "string" },
                    minItems: 0,
                    maxItems: 10,
                    uniqueItems: true
                },
                createdAt: {
                    bsonType: "date"
                }
            },
            additionalProperties: false  // Reject unknown fields
        }
    },
    validationLevel: "strict",     // "off", "moderate", "strict"
    validationAction: "error"      // "error" or "warn"
})
```

### Modifying Validation

```javascript
// Add/modify validation on existing collection
db.runCommand({
    collMod: "users",
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["name", "email"],
            properties: {
                // Updated schema
            }
        }
    },
    validationLevel: "moderate",  // Only validate new/updated docs
    validationAction: "warn"      // Log warnings instead of rejecting
})

// Get current validation
db.getCollectionInfos({ name: "users" })[0].options.validator
```

### Validation with Query Operators

```javascript
// Using query operators for validation
db.createCollection("products", {
    validator: {
        $and: [
            { price: { $type: "number", $gt: 0 } },
            { quantity: { $type: "int", $gte: 0 } },
            { category: { $in: ["electronics", "clothing", "food"] } },
            { $or: [
                { discountPrice: { $exists: false } },
                { $expr: { $lt: ["$discountPrice", "$price"] } }
            ]}
        ]
    }
})
```

### Validation Best Practices

```javascript
// 1. Start with moderate validation during migration
db.runCommand({
    collMod: "users",
    validationLevel: "moderate",  // Validate only inserts/updates
    validationAction: "warn"      // Log, don't reject
})

// 2. Gradually tighten validation
// After fixing existing documents:
db.runCommand({
    collMod: "users",
    validationLevel: "strict",
    validationAction: "error"
})

// 3. Use bypass for admin operations
db.users.insertOne(
    { name: "Admin", email: "admin@example.com" },
    { bypassDocumentValidation: true }
)

// 4. Check validation errors
db.users.insertOne({ name: "Invalid" })
// Error: Document failed validation
```

---

## 10. Transactions

### Single Document Transactions

```javascript
// MongoDB guarantees atomicity for single document operations
// All fields in one document update atomically
db.accounts.updateOne(
    { _id: "account123" },
    {
        $inc: { balance: -100 },
        $push: {
            transactions: {
                type: "debit",
                amount: 100,
                date: new Date()
            }
        }
    }
)
```

### Multi-Document Transactions (4.0+)

```javascript
// Start a session
const session = db.getMongo().startSession()

// Start transaction
session.startTransaction({
    readConcern: { level: "snapshot" },
    writeConcern: { w: "majority" }
})

try {
    const accounts = session.getDatabase("bank").accounts
    
    // Transfer money between accounts
    accounts.updateOne(
        { _id: "account1" },
        { $inc: { balance: -100 } },
        { session }
    )
    
    accounts.updateOne(
        { _id: "account2" },
        { $inc: { balance: 100 } },
        { session }
    )
    
    // Commit transaction
    session.commitTransaction()
    print("Transaction committed")
} catch (error) {
    // Abort on error
    session.abortTransaction()
    print("Transaction aborted: " + error)
} finally {
    session.endSession()
}
```

### Transaction with Retry Logic

```javascript
// Recommended pattern with retry
function runTransactionWithRetry(txnFunc, session) {
    while (true) {
        try {
            txnFunc(session)
            break
        } catch (error) {
            if (
                error.hasOwnProperty("errorLabels") &&
                error.errorLabels.includes("TransientTransactionError")
            ) {
                print("TransientTransactionError, retrying...")
                continue
            } else {
                throw error
            }
        }
    }
}

function commitWithRetry(session) {
    while (true) {
        try {
            session.commitTransaction()
            print("Transaction committed")
            break
        } catch (error) {
            if (
                error.hasOwnProperty("errorLabels") &&
                error.errorLabels.includes("UnknownTransactionCommitResult")
            ) {
                print("UnknownTransactionCommitResult, retrying...")
                continue
            } else {
                throw error
            }
        }
    }
}

// Usage
const session = db.getMongo().startSession()
try {
    runTransactionWithRetry(() => {
        session.startTransaction()
        
        // Your operations here
        db.accounts.updateOne({ _id: "acc1" }, { $inc: { balance: -100 } }, { session })
        db.accounts.updateOne({ _id: "acc2" }, { $inc: { balance: 100 } }, { session })
        
        commitWithRetry(session)
    }, session)
} finally {
    session.endSession()
}
```

### Read Concern Levels

```javascript
// local - Returns most recent data (may not be majority committed)
db.collection.find().readConcern("local")

// available - Same as local for non-sharded collections
db.collection.find().readConcern("available")

// majority - Returns data acknowledged by majority of replica set
db.collection.find().readConcern("majority")

// linearizable - Reflects all successful majority writes
db.collection.find().readConcern("linearizable")

// snapshot - For multi-document transactions
session.startTransaction({ readConcern: { level: "snapshot" } })
```

### Write Concern Levels

```javascript
// w: 1 - Acknowledge write on primary only
db.collection.insertOne(doc, { writeConcern: { w: 1 } })

// w: "majority" - Acknowledge when majority of nodes have written
db.collection.insertOne(doc, { writeConcern: { w: "majority" } })

// w: 0 - No acknowledgment (fire and forget)
db.collection.insertOne(doc, { writeConcern: { w: 0 } })

// j: true - Wait for journal commit
db.collection.insertOne(doc, { writeConcern: { w: 1, j: true } })

// wtimeout - Timeout for write concern
db.collection.insertOne(doc, { writeConcern: { w: "majority", wtimeout: 5000 } })
```

---

## 11. Replication

### Replica Set Overview

```
┌─────────────────────────────────────────────────────────────┐
│                      REPLICA SET                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────┐    Sync    ┌─────────────┐                │
│  │   PRIMARY   │ ────────►  │  SECONDARY  │                │
│  │  (Read/Write)│           │  (Read Only) │                │
│  └─────────────┘            └─────────────┘                │
│         │                          ▲                        │
│         │         Sync             │                        │
│         └──────────────────────────┤                        │
│                                    │                        │
│                          ┌─────────────┐                   │
│                          │  SECONDARY  │                   │
│                          │  (Read Only) │                   │
│                          └─────────────┘                   │
│                                                             │
│  ┌─────────────┐  (Optional - voting only)                 │
│  │   ARBITER   │                                           │
│  └─────────────┘                                           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Setting Up Replica Set

```javascript
// Start mongod with replica set name
// mongod --replSet rs0 --port 27017 --dbpath /data/rs0-0
// mongod --replSet rs0 --port 27018 --dbpath /data/rs0-1
// mongod --replSet rs0 --port 27019 --dbpath /data/rs0-2

// Connect to one member and initiate
rs.initiate({
    _id: "rs0",
    members: [
        { _id: 0, host: "localhost:27017" },
        { _id: 1, host: "localhost:27018" },
        { _id: 2, host: "localhost:27019" }
    ]
})

// Check status
rs.status()

// Check configuration
rs.conf()

// View current primary
rs.isMaster()
db.hello()  // MongoDB 5.0+
```

### Managing Replica Set

```javascript
// Add member
rs.add("localhost:27020")

// Add with specific configuration
rs.add({
    host: "localhost:27020",
    priority: 0.5,
    votes: 1,
    hidden: false
})

// Add arbiter
rs.addArb("localhost:27021")

// Remove member
rs.remove("localhost:27020")

// Change configuration
cfg = rs.conf()
cfg.members[1].priority = 2
rs.reconfig(cfg)

// Force reconfiguration
rs.reconfig(cfg, { force: true })

// Step down primary (trigger election)
rs.stepDown(60)  // Step down for 60 seconds
```

### Member Configuration Options

```javascript
rs.reconfig({
    _id: "rs0",
    members: [
        {
            _id: 0,
            host: "mongo1:27017",
            priority: 2,           // Higher = more likely to be primary
            votes: 1
        },
        {
            _id: 1,
            host: "mongo2:27017",
            priority: 1,
            votes: 1,
            secondaryDelaySecs: 3600  // 1 hour delayed (for recovery)
        },
        {
            _id: 2,
            host: "mongo3:27017",
            priority: 0,           // Can never be primary
            votes: 1,
            hidden: true           // Hidden from clients
        },
        {
            _id: 3,
            host: "arbiter:27017",
            arbiterOnly: true      // Only for voting
        }
    ],
    settings: {
        electionTimeoutMillis: 10000,
        heartbeatIntervalMillis: 2000
    }
})
```

### Read Preference

```javascript
// primary - Always read from primary (default)
db.collection.find().readPref("primary")

// primaryPreferred - Primary if available, else secondary
db.collection.find().readPref("primaryPreferred")

// secondary - Always read from secondary
db.collection.find().readPref("secondary")

// secondaryPreferred - Secondary if available, else primary
db.collection.find().readPref("secondaryPreferred")

// nearest - Read from nearest member (lowest latency)
db.collection.find().readPref("nearest")

// With tags
db.collection.find().readPref("secondary", [{ dc: "east" }])
```

### Monitoring Replication

```javascript
// Replica set status
rs.status()

// Replication lag
rs.printSecondaryReplicationInfo()

// Oplog info
rs.printReplicationInfo()

// Oplog size
use local
db.oplog.rs.stats().maxSize

// Change oplog size
db.adminCommand({ replSetResizeOplog: 1, size: 16384 })  // MB
```

---

## 12. Sharding

### Sharding Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     SHARDED CLUSTER                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│                      ┌──────────┐                          │
│                      │  mongos  │  (Query Router)          │
│                      │  Router  │                          │
│                      └──────────┘                          │
│                            │                                │
│              ┌─────────────┼─────────────┐                 │
│              ▼             ▼             ▼                 │
│        ┌──────────┐  ┌──────────┐  ┌──────────┐          │
│        │  Shard 1 │  │  Shard 2 │  │  Shard 3 │          │
│        │(Replica  │  │(Replica  │  │(Replica  │          │
│        │   Set)   │  │   Set)   │  │   Set)   │          │
│        └──────────┘  └──────────┘  └──────────┘          │
│              │             │             │                 │
│              └─────────────┼─────────────┘                 │
│                            ▼                                │
│                    ┌──────────────┐                        │
│                    │Config Servers│                        │
│                    │ (Replica Set)│                        │
│                    └──────────────┘                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Setting Up Sharded Cluster

```javascript
// 1. Start config servers (replica set)
// mongod --configsvr --replSet configRS --port 27019 --dbpath /data/config

// 2. Start shard servers (replica sets)
// mongod --shardsvr --replSet shard1RS --port 27018 --dbpath /data/shard1
// mongod --shardsvr --replSet shard2RS --port 27028 --dbpath /data/shard2

// 3. Start mongos router
// mongos --configdb configRS/localhost:27019 --port 27017

// 4. Connect to mongos and add shards
sh.addShard("shard1RS/localhost:27018")
sh.addShard("shard2RS/localhost:27028")

// 5. Enable sharding on database
sh.enableSharding("myapp")

// 6. Shard a collection
sh.shardCollection("myapp.users", { region: 1 })  // Range sharding
sh.shardCollection("myapp.logs", { _id: "hashed" })  // Hashed sharding
```

### Shard Key Strategies

```javascript
// Ranged Sharding (good for range queries)
sh.shardCollection("myapp.orders", { orderDate: 1 })
// Data distributed by date ranges

// Hashed Sharding (even distribution)
sh.shardCollection("myapp.events", { userId: "hashed" })
// Hash of userId determines shard

// Compound Shard Key
sh.shardCollection("myapp.products", { category: 1, productId: 1 })
// Queries on category hit fewer shards

// Zone Sharding (direct data to specific shards)
sh.addShardTag("shard1", "US")
sh.addShardTag("shard2", "EU")

sh.updateZoneKeyRange(
    "myapp.users",
    { region: "US" },
    { region: "USS" },  // Range end
    "US"
)

sh.updateZoneKeyRange(
    "myapp.users",
    { region: "EU" },
    { region: "EUS" },
    "EU"
)
```

### Managing Sharded Cluster

```javascript
// Cluster status
sh.status()

// Shard distribution
db.collection.getShardDistribution()

// Chunk operations
sh.splitAt("myapp.users", { region: "M" })  // Split chunk at key
sh.moveChunk("myapp.users", { region: "US" }, "shard2")  // Move chunk

// Balancer control
sh.getBalancerState()
sh.startBalancer()
sh.stopBalancer()
sh.setBalancerState(true)

// Balancer window
use config
db.settings.updateOne(
    { _id: "balancer" },
    { $set: { activeWindow: { start: "02:00", stop: "06:00" } } },
    { upsert: true }
)
```

### Shard Key Selection Best Practices

```javascript
// Good shard key characteristics:
// 1. High cardinality (many unique values)
// 2. Distributed writes (avoid hotspots)
// 3. Query isolation (queries target few shards)

// GOOD: Compound key with high cardinality
{ tenantId: 1, timestamp: 1 }

// GOOD: Hashed key for write distribution
{ orderId: "hashed" }

// BAD: Low cardinality (only few values)
{ status: 1 }  // Creates hotspots

// BAD: Monotonically increasing (timestamp, auto-increment)
{ _id: 1 }  // All writes go to one shard
{ timestamp: 1 }  // Same problem

// Consider: Hashed for monotonic keys
{ timestamp: "hashed" }
```

---

## 13. Security

### Authentication

```javascript
// Enable authentication in mongod.conf
// security:
//   authorization: enabled

// Create admin user
use admin
db.createUser({
    user: "admin",
    pwd: "securePassword123",
    roles: [
        { role: "userAdminAnyDatabase", db: "admin" },
        { role: "readWriteAnyDatabase", db: "admin" },
        { role: "clusterAdmin", db: "admin" }
    ]
})

// Create application user
use myapp
db.createUser({
    user: "appUser",
    pwd: "appPassword123",
    roles: [
        { role: "readWrite", db: "myapp" }
    ]
})

// Create read-only user
db.createUser({
    user: "reader",
    pwd: "readerPassword",
    roles: [{ role: "read", db: "myapp" }]
})

// Authenticate
db.auth("admin", "securePassword123")

// Connect with authentication
// mongosh -u admin -p securePassword123 --authenticationDatabase admin
```

### Role-Based Access Control (RBAC)

```javascript
// Built-in roles:
// Database: read, readWrite, dbAdmin, userAdmin, dbOwner
// Cluster: clusterManager, clusterMonitor, clusterAdmin
// All databases: readAnyDatabase, readWriteAnyDatabase, userAdminAnyDatabase, dbAdminAnyDatabase
// Superuser: root

// Create custom role
use myapp
db.createRole({
    role: "orderManager",
    privileges: [
        {
            resource: { db: "myapp", collection: "orders" },
            actions: ["find", "update", "insert"]
        },
        {
            resource: { db: "myapp", collection: "products" },
            actions: ["find"]
        }
    ],
    roles: []  // Can inherit from other roles
})

// Create user with custom role
db.createUser({
    user: "orderProcessor",
    pwd: "orderPass123",
    roles: [{ role: "orderManager", db: "myapp" }]
})

// Grant additional role
db.grantRolesToUser("orderProcessor", [{ role: "read", db: "analytics" }])

// Revoke role
db.revokeRolesFromUser("orderProcessor", [{ role: "read", db: "analytics" }])

// View user info
db.getUser("orderProcessor")

// View role info
db.getRole("orderManager", { showPrivileges: true })
```

### Encryption

```javascript
// 1. TLS/SSL in Transit
// mongod.conf
// net:
//   tls:
//     mode: requireTLS
//     certificateKeyFile: /etc/mongodb/server.pem
//     CAFile: /etc/mongodb/ca.pem

// Connect with TLS
// mongosh --tls --host localhost --tlsCertificateKeyFile client.pem --tlsCAFile ca.pem

// 2. Encryption at Rest (Enterprise)
// mongod.conf
// security:
//   enableEncryption: true
//   encryptionKeyFile: /etc/mongodb/encryption-key

// 3. Client-Side Field Level Encryption (CSFLE)
// Encrypt sensitive fields before sending to server
// Available in drivers (Node.js, Python, Java, etc.)
```

### Network Security

```javascript
// Bind to specific IP
// net:
//   bindIp: 127.0.0.1,192.168.1.100
//   port: 27017

// IP Whitelist (recommended for production)
// net:
//   bindIp: 0.0.0.0
// + Use firewall rules

// Disable HTTP interface
// net:
//   http:
//     enabled: false
```

### Auditing (Enterprise)

```javascript
// Enable auditing in mongod.conf
// auditLog:
//   destination: file
//   format: JSON
//   path: /var/log/mongodb/audit.json
//   filter: '{ atype: { $in: ["authCheck", "authenticate"] } }'

// Audit all auth events
// auditLog:
//   destination: file
//   format: JSON
//   path: /var/log/mongodb/audit.json
```

---

## 14. Backup & Recovery

### mongodump / mongorestore

```bash
# Backup entire database
mongodump --db myapp --out /backup/$(date +%Y%m%d)

# Backup with authentication
mongodump --uri "mongodb://admin:password@localhost:27017" --db myapp --out /backup

# Backup specific collection
mongodump --db myapp --collection users --out /backup

# Backup with query
mongodump --db myapp --collection logs --query '{"date": {"$gte": {"$date": "2024-01-01T00:00:00Z"}}}' --out /backup

# Backup with compression
mongodump --db myapp --gzip --archive=/backup/myapp.archive.gz

# Backup sharded cluster (from mongos)
mongodump --host mongos.example.com --port 27017 --out /backup

# Restore entire backup
mongorestore --db myapp /backup/myapp

# Restore with drop (replace existing)
mongorestore --db myapp --drop /backup/myapp

# Restore specific collection
mongorestore --db myapp --collection users /backup/myapp/users.bson

# Restore from archive
mongorestore --gzip --archive=/backup/myapp.archive.gz

# Restore to different database
mongorestore --db myapp_restored /backup/myapp
```

### Point-in-Time Recovery with Oplog

```bash
# Backup with oplog (for PITR)
mongodump --db myapp --oplog --out /backup

# Restore with oplog replay
mongorestore --oplogReplay /backup

# Replay oplog to specific point
mongorestore --oplogReplay --oplogLimit "1234567890:1" /backup
```

### Filesystem Snapshots

```bash
# For WiredTiger (default storage engine)
# 1. Lock writes
mongosh --eval "db.fsyncLock()"

# 2. Take filesystem snapshot (LVM, EBS, etc.)
lvcreate --snapshot --name mongodb-snapshot --size 10G /dev/vg/mongodb

# 3. Unlock
mongosh --eval "db.fsyncUnlock()"

# 4. Copy data from snapshot
mount /dev/vg/mongodb-snapshot /mnt/snapshot
cp -r /mnt/snapshot/mongodb /backup/
```

### MongoDB Atlas Backup

```javascript
// Atlas provides automated backups:
// - Continuous backups with point-in-time recovery
// - Scheduled snapshots
// - Cross-region restore

// Restore via Atlas UI or API
// Download snapshots for local restore
```

### Backup Best Practices

```bash
#!/bin/bash
# backup_mongodb.sh

BACKUP_DIR="/backup/mongodb"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=30
MONGODB_URI="mongodb://backup_user:password@localhost:27017"

# Create dated backup directory
mkdir -p $BACKUP_DIR/$DATE

# Backup with oplog for PITR capability
mongodump \
    --uri="$MONGODB_URI" \
    --oplog \
    --gzip \
    --out $BACKUP_DIR/$DATE

# Check backup success
if [ $? -eq 0 ]; then
    echo "Backup completed: $BACKUP_DIR/$DATE"
    
    # Upload to S3
    aws s3 sync $BACKUP_DIR/$DATE s3://my-bucket/mongodb-backups/$DATE/
    
    # Clean old backups
    find $BACKUP_DIR -type d -mtime +$RETENTION_DAYS -exec rm -rf {} +
else
    echo "Backup failed!"
    exit 1
fi
```

---

## 15. Performance Tuning

### Hardware Considerations

```
┌─────────────────────────────────────────────────────────────┐
│                  HARDWARE REQUIREMENTS                      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  RAM:                                                       │
│  └── Working set should fit in RAM                         │
│  └── More RAM = less disk I/O                              │
│  └── Rule: 50% for WiredTiger cache                        │
│                                                             │
│  Storage:                                                   │
│  └── SSD strongly recommended                              │
│  └── Separate disk for journal                             │
│  └── RAID 10 for redundancy + performance                  │
│                                                             │
│  CPU:                                                       │
│  └── More cores for aggregation, indexing                  │
│  └── Fast single-thread for single queries                 │
│                                                             │
│  Network:                                                   │
│  └── Low latency for replica set communication             │
│  └── High bandwidth for replication, sharding              │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### WiredTiger Configuration

```yaml
# mongod.conf
storage:
  wiredTiger:
    engineConfig:
      cacheSizeGB: 4          # Default: 50% of (RAM - 1GB)
      journalCompressor: snappy
      directoryForIndexes: true  # Separate indexes on different disk
    collectionConfig:
      blockCompressor: snappy  # snappy, zlib, zstd, none
    indexConfig:
      prefixCompression: true
```

### Connection Pooling

```javascript
// Connection string with pool settings
// mongodb://host:27017/?maxPoolSize=100&minPoolSize=10&maxIdleTimeMS=30000

// Node.js driver example
const client = new MongoClient(uri, {
    maxPoolSize: 100,
    minPoolSize: 10,
    maxIdleTimeMS: 30000,
    waitQueueTimeoutMS: 5000
})

// Monitor connections
db.serverStatus().connections
// {
//   current: 50,
//   available: 51150,
//   totalCreated: 1234
// }
```

### Query Profiling

```javascript
// Enable profiler
db.setProfilingLevel(1, { slowms: 100 })  // Log queries > 100ms
db.setProfilingLevel(2)  // Log all queries
db.setProfilingLevel(0)  // Disable

// Check profiler status
db.getProfilingStatus()

// Query profiler data
db.system.profile.find({ millis: { $gt: 100 } }).sort({ ts: -1 }).limit(10)

// Common slow query patterns
db.system.profile.find({
    op: "query",
    "command.filter": { $exists: true }
}).sort({ millis: -1 }).limit(10)
```

### Monitoring Queries

```javascript
// Current operations
db.currentOp()
db.currentOp({ "active": true, "secs_running": { $gt: 5 } })

// Kill long-running operation
db.killOp(opId)

// Server statistics
db.serverStatus()
db.serverStatus().opcounters  // Operation counts
db.serverStatus().metrics     // Detailed metrics
db.serverStatus().wiredTiger  // Storage engine stats

// Collection statistics
db.collection.stats()

// Aggregation for analysis
db.collection.aggregate([{ $collStats: { storageStats: {} } }])
```

### Index Optimization

```javascript
// Find unused indexes
db.collection.aggregate([{ $indexStats: {} }])
// Look for indexes with ops: 0

// Analyze query plans
db.collection.find({ field: value }).explain("executionStats")

// Key metrics:
// - totalDocsExamined vs nReturned (should be close)
// - executionTimeMillis
// - indexesUsed

// Hint specific index
db.collection.find({ field: value }).hint({ field: 1 })

// Build indexes in background
db.collection.createIndex({ field: 1 }, { background: true })
```

### Memory Management

```javascript
// Check memory usage
db.serverStatus().mem
// {
//   bits: 64,
//   resident: 4096,  // MB
//   virtual: 8192,   // MB
//   supported: true
// }

// WiredTiger cache status
db.serverStatus().wiredTiger.cache

// Clear cache (restart required)
// Reduce working set or add RAM
```

---

## 16. MongoDB with Applications

### Node.js Driver

```javascript
// Installation
// npm install mongodb

const { MongoClient, ObjectId } = require('mongodb')

const uri = "mongodb://localhost:27017"
const client = new MongoClient(uri)

async function main() {
    try {
        await client.connect()
        const db = client.db("myapp")
        const users = db.collection("users")
        
        // Insert
        const result = await users.insertOne({
            name: "John Doe",
            email: "john@example.com",
            createdAt: new Date()
        })
        console.log("Inserted:", result.insertedId)
        
        // Find
        const user = await users.findOne({ email: "john@example.com" })
        console.log("Found:", user)
        
        // Find with options
        const allUsers = await users.find({})
            .sort({ createdAt: -1 })
            .limit(10)
            .toArray()
        
        // Update
        await users.updateOne(
            { _id: result.insertedId },
            { $set: { name: "John Smith" } }
        )
        
        // Delete
        await users.deleteOne({ _id: result.insertedId })
        
        // Aggregation
        const stats = await users.aggregate([
            { $group: { _id: "$status", count: { $sum: 1 } } }
        ]).toArray()
        
        // Transaction
        const session = client.startSession()
        try {
            session.startTransaction()
            await users.updateOne({ _id: id1 }, { $inc: { balance: -100 } }, { session })
            await users.updateOne({ _id: id2 }, { $inc: { balance: 100 } }, { session })
            await session.commitTransaction()
        } catch (error) {
            await session.abortTransaction()
            throw error
        } finally {
            session.endSession()
        }
        
    } finally {
        await client.close()
    }
}

main().catch(console.error)
```

### Mongoose (ODM for Node.js)

```javascript
// Installation
// npm install mongoose

const mongoose = require('mongoose')

// Connect
mongoose.connect('mongodb://localhost:27017/myapp')

// Define schema
const userSchema = new mongoose.Schema({
    name: { type: String, required: true },
    email: { type: String, required: true, unique: true, lowercase: true },
    age: { type: Number, min: 0, max: 150 },
    status: { type: String, enum: ['active', 'inactive'], default: 'active' },
    tags: [String],
    address: {
        street: String,
        city: String,
        zip: String
    },
    createdAt: { type: Date, default: Date.now }
}, {
    timestamps: true  // Adds createdAt and updatedAt
})

// Add index
userSchema.index({ email: 1 })
userSchema.index({ name: 'text' })

// Add method
userSchema.methods.getFullName = function() {
    return `${this.firstName} ${this.lastName}`
}

// Add static
userSchema.statics.findByEmail = function(email) {
    return this.findOne({ email: email.toLowerCase() })
}

// Add virtual
userSchema.virtual('isAdult').get(function() {
    return this.age >= 18
})

// Pre-save hook
userSchema.pre('save', function(next) {
    this.updatedAt = new Date()
    next()
})

// Create model
const User = mongoose.model('User', userSchema)

// CRUD operations
async function examples() {
    // Create
    const user = new User({ name: 'John', email: 'john@example.com' })
    await user.save()
    
    // Or
    const user2 = await User.create({ name: 'Jane', email: 'jane@example.com' })
    
    // Find
    const foundUser = await User.findById(user._id)
    const byEmail = await User.findOne({ email: 'john@example.com' })
    const allActive = await User.find({ status: 'active' })
    
    // Update
    await User.updateOne({ _id: user._id }, { name: 'John Smith' })
    await User.findByIdAndUpdate(user._id, { name: 'John Doe' }, { new: true })
    
    // Delete
    await User.deleteOne({ _id: user._id })
    await User.findByIdAndDelete(user._id)
    
    // Query builder
    const results = await User.find()
        .where('age').gte(18).lte(65)
        .where('status').equals('active')
        .select('name email')
        .sort('-createdAt')
        .limit(10)
        .exec()
    
    // Populate (join)
    const order = await Order.findById(orderId).populate('userId')
}
```

### Python (PyMongo)

```python
# Installation
# pip install pymongo

from pymongo import MongoClient
from bson.objectid import ObjectId
from datetime import datetime

# Connect
client = MongoClient('mongodb://localhost:27017/')
db = client['myapp']
users = db['users']

# Insert
result = users.insert_one({
    'name': 'John Doe',
    'email': 'john@example.com',
    'created_at': datetime.utcnow()
})
print(f"Inserted: {result.inserted_id}")

# Find
user = users.find_one({'email': 'john@example.com'})
all_users = list(users.find({'status': 'active'}).sort('created_at', -1).limit(10))

# Update
users.update_one(
    {'_id': result.inserted_id},
    {'$set': {'name': 'John Smith'}}
)

# Delete
users.delete_one({'_id': result.inserted_id})

# Aggregation
pipeline = [
    {'$match': {'status': 'active'}},
    {'$group': {'_id': '$category', 'count': {'$sum': 1}}}
]
results = list(users.aggregate(pipeline))

# Transaction
with client.start_session() as session:
    with session.start_transaction():
        users.update_one({'_id': id1}, {'$inc': {'balance': -100}}, session=session)
        users.update_one({'_id': id2}, {'$inc': {'balance': 100}}, session=session)
```

### Java Driver

```java
// Maven dependency
// <dependency>
//     <groupId>org.mongodb</groupId>
//     <artifactId>mongodb-driver-sync</artifactId>
//     <version>4.9.0</version>
// </dependency>

import com.mongodb.client.*;
import org.bson.Document;
import static com.mongodb.client.model.Filters.*;
import static com.mongodb.client.model.Updates.*;

public class MongoExample {
    public static void main(String[] args) {
        try (MongoClient client = MongoClients.create("mongodb://localhost:27017")) {
            MongoDatabase db = client.getDatabase("myapp");
            MongoCollection<Document> users = db.getCollection("users");
            
            // Insert
            Document user = new Document("name", "John Doe")
                .append("email", "john@example.com")
                .append("createdAt", new Date());
            users.insertOne(user);
            
            // Find
            Document found = users.find(eq("email", "john@example.com")).first();
            
            // Update
            users.updateOne(
                eq("email", "john@example.com"),
                set("name", "John Smith")
            );
            
            // Delete
            users.deleteOne(eq("email", "john@example.com"));
            
            // Aggregation
            List<Document> results = users.aggregate(Arrays.asList(
                match(eq("status", "active")),
                group("$category", sum("count", 1))
            )).into(new ArrayList<>());
        }
    }
}
```

---

## 17. Atlas (Cloud MongoDB)

### Getting Started with Atlas

```
1. Create Atlas Account: https://cloud.mongodb.com
2. Create Organization and Project
3. Build a Cluster:
   - Choose cloud provider (AWS, GCP, Azure)
   - Select region
   - Choose cluster tier (M0 Free, M10+, etc.)
4. Configure Security:
   - Add IP whitelist
   - Create database user
5. Get connection string
```

### Connection String

```javascript
// Atlas connection string format
const uri = "mongodb+srv://<username>:<password>@cluster0.xxxxx.mongodb.net/<database>?retryWrites=true&w=majority"

// Node.js
const { MongoClient } = require('mongodb')
const client = new MongoClient(uri)

// Python
from pymongo import MongoClient
client = MongoClient(uri)

// mongosh
// mongosh "mongodb+srv://cluster0.xxxxx.mongodb.net" --username admin
```

### Atlas Features

```javascript
// Atlas Search (Full-text search)
// Create search index in Atlas UI, then query:
db.products.aggregate([
    {
        $search: {
            index: "default",
            text: {
                query: "laptop gaming",
                path: ["name", "description"],
                fuzzy: { maxEdits: 1 }
            }
        }
    },
    { $limit: 10 },
    { $project: { name: 1, score: { $meta: "searchScore" } } }
])

// Atlas Charts - Embedded dashboards
// Create visualizations in Atlas UI and embed:
// <iframe src="https://charts.mongodb.com/..." />

// Atlas Data Lake - Query data in S3
// db.s3data.find({ year: 2024 })

// Atlas App Services (formerly Realm)
// - Serverless functions
// - GraphQL API
// - Data API (REST)
// - Device Sync
```

### Atlas Administration

```javascript
// Using Atlas Admin API
// curl --user "public_key:private_key" \
//   "https://cloud.mongodb.com/api/atlas/v1.0/groups/{projectId}/clusters"

// Atlas CLI
// atlas clusters list
// atlas clusters create myCluster --provider AWS --region US_EAST_1 --tier M10

// Monitoring
// - Real-time performance metrics in Atlas UI
// - Alerts configuration
// - Query Profiler
// - Performance Advisor (index recommendations)
```

---

## 18. Best Practices

### Schema Design

```javascript
// 1. Design for your queries, not data relationships
// Think: "What questions will I ask?"

// 2. Embed data that is accessed together
{
    _id: ObjectId(),
    title: "Product",
    reviews: [  // Embed if always needed with product
        { user: "john", rating: 5, comment: "Great!" }
    ]
}

// 3. Reference data accessed independently or unbounded
{
    _id: ObjectId("product123"),
    title: "Product"
}
{
    _id: ObjectId(),
    productId: ObjectId("product123"),  // Reference
    user: "john",
    rating: 5
}

// 4. Avoid unbounded arrays
// BAD: Array can grow indefinitely
{ followers: [userId1, userId2, ...millionsMore] }

// GOOD: Reverse the reference
{ userId: X, followsUser: Y }

// 5. Use appropriate data types
{
    price: NumberDecimal("19.99"),  // Not string or float
    createdAt: new Date(),           // Not string
    userId: ObjectId("...")          // Not string
}
```

### Query Optimization

```javascript
// 1. Create indexes for common queries
db.orders.createIndex({ customerId: 1, orderDate: -1 })

// 2. Use covered queries when possible
db.users.createIndex({ email: 1, name: 1 })
db.users.find({ email: "x" }, { email: 1, name: 1, _id: 0 })

// 3. Limit returned fields
db.users.find({}, { name: 1, email: 1 })  // Not find({})

// 4. Use $match early in aggregations
db.orders.aggregate([
    { $match: { status: "pending" } },  // Filter first!
    { $group: { ... } }
])

// 5. Avoid $where and JavaScript execution
// BAD:
db.users.find({ $where: "this.age > 18" })
// GOOD:
db.users.find({ age: { $gt: 18 } })

// 6. Use explain() to analyze
db.orders.find({ status: "pending" }).explain("executionStats")
```

### Application Best Practices

```javascript
// 1. Use connection pooling
const client = new MongoClient(uri, {
    maxPoolSize: 50,
    minPoolSize: 10
})

// 2. Handle connection errors
client.on('error', (error) => {
    console.error('MongoDB connection error:', error)
})

// 3. Use retryable writes
const uri = "mongodb://...?retryWrites=true"

// 4. Set appropriate timeouts
const client = new MongoClient(uri, {
    serverSelectionTimeoutMS: 5000,
    socketTimeoutMS: 45000
})

// 5. Use bulk operations for multiple writes
const bulk = db.users.initializeUnorderedBulkOp()
bulk.insert({ name: "User1" })
bulk.insert({ name: "User2" })
bulk.find({ status: "inactive" }).update({ $set: { archived: true } })
await bulk.execute()

// 6. Handle ObjectId properly
const { ObjectId } = require('mongodb')
const id = new ObjectId(stringId)
```

### Security Best Practices

```javascript
// 1. Enable authentication
// security:
//   authorization: enabled

// 2. Use least privilege
db.createUser({
    user: "appUser",
    pwd: "strongPassword",
    roles: [{ role: "readWrite", db: "myapp" }]  // Only what's needed
})

// 3. Enable TLS
// net:
//   tls:
//     mode: requireTLS

// 4. Sanitize user input (prevent injection)
// BAD:
const query = { $where: `this.name == '${userInput}'` }

// GOOD:
const query = { name: userInput }

// 5. Use IP whitelisting in Atlas

// 6. Rotate credentials regularly
```

### Operations Best Practices

```javascript
// 1. Monitor key metrics
db.serverStatus().opcounters
db.serverStatus().connections
db.serverStatus().mem

// 2. Set up alerts for:
// - High CPU/memory usage
// - Replication lag
// - Connection pool exhaustion
// - Slow queries

// 3. Regular backups with testing
// Test restore process quarterly

// 4. Use replica sets for HA
// Minimum 3 members (or 2 + arbiter)

// 5. Plan for capacity
// Monitor growth trends
// Scale before hitting limits

// 6. Document everything
// Schema design decisions
// Index strategy
// Backup procedures
```

---

## 19. Common Patterns

### Caching Pattern

```javascript
// Cache frequently accessed data with TTL
db.cache.createIndex({ expireAt: 1 }, { expireAfterSeconds: 0 })

// Store cache entry
db.cache.insertOne({
    key: "user:123:profile",
    data: { name: "John", email: "john@example.com" },
    expireAt: new Date(Date.now() + 3600000)  // 1 hour
})

// Read from cache
const cached = await db.cache.findOne({ key: "user:123:profile" })
if (cached) {
    return cached.data
}
// Cache miss - fetch from source and cache
```

### Event Sourcing

```javascript
// Store events, not state
db.events.insertOne({
    aggregateId: "order-123",
    type: "OrderCreated",
    data: { items: [...], total: 99.99 },
    timestamp: new Date(),
    version: 1
})

db.events.insertOne({
    aggregateId: "order-123",
    type: "PaymentReceived",
    data: { amount: 99.99, method: "card" },
    timestamp: new Date(),
    version: 2
})

// Rebuild state from events
const events = await db.events.find({ aggregateId: "order-123" })
    .sort({ version: 1 })
    .toArray()

let state = {}
events.forEach(event => {
    state = applyEvent(state, event)
})
```

### Queue Pattern

```javascript
// Simple job queue
db.jobs.createIndex({ status: 1, priority: -1, createdAt: 1 })

// Add job
db.jobs.insertOne({
    type: "send_email",
    data: { to: "user@example.com", subject: "Hello" },
    status: "pending",
    priority: 1,
    createdAt: new Date(),
    attempts: 0
})

// Claim and process job (atomic)
const job = await db.jobs.findOneAndUpdate(
    { status: "pending" },
    { 
        $set: { status: "processing", startedAt: new Date() },
        $inc: { attempts: 1 }
    },
    { sort: { priority: -1, createdAt: 1 }, returnDocument: "after" }
)

// Complete or fail
await db.jobs.updateOne(
    { _id: job._id },
    { $set: { status: "completed", completedAt: new Date() } }
)
```

### Singleton Pattern

```javascript
// Ensure only one document of a type exists
db.config.createIndex({ type: 1 }, { unique: true })

// Upsert configuration
db.config.updateOne(
    { type: "app_settings" },
    { 
        $set: { 
            type: "app_settings",
            maintenance: false,
            version: "1.0.0"
        }
    },
    { upsert: true }
)
```

### Rate Limiting

```javascript
// Track requests per user
db.rateLimit.createIndex({ userId: 1, windowStart: 1 })
db.rateLimit.createIndex({ windowStart: 1 }, { expireAfterSeconds: 3600 })

async function checkRateLimit(userId, limit = 100) {
    const windowStart = new Date()
    windowStart.setMinutes(0, 0, 0)  // Start of hour
    
    const result = await db.rateLimit.findOneAndUpdate(
        { userId, windowStart },
        { $inc: { count: 1 } },
        { upsert: true, returnDocument: "after" }
    )
    
    return result.value.count <= limit
}
```

---

## 20. Troubleshooting

### Common Issues

```javascript
// 1. Connection refused
// Check: mongod running, port correct, firewall rules
// mongosh --host localhost --port 27017

// 2. Authentication failed
// Check: username, password, authSource
// mongosh -u user -p pass --authenticationDatabase admin

// 3. Write concern error
// Solution: Check replica set health
rs.status()

// 4. Timeout errors
// Solution: Check network, increase timeout, optimize query
const client = new MongoClient(uri, { serverSelectionTimeoutMS: 10000 })

// 5. Duplicate key error
// Solution: Check unique indexes, handle upserts properly
db.collection.createIndex({ email: 1 }, { unique: true })
```

### Diagnostic Commands

```javascript
// Server information
db.serverStatus()
db.hostInfo()
db.version()

// Current operations
db.currentOp()
db.currentOp({ "active": true })

// Kill operation
db.killOp(opId)

// Profiler
db.setProfilingLevel(2)
db.system.profile.find().sort({ ts: -1 }).limit(10)

// Lock information
db.serverStatus().locks
db.serverStatus().globalLock

// Replication status
rs.status()
rs.printReplicationInfo()
rs.printSecondaryReplicationInfo()

// Sharding status
sh.status()
db.printShardingStatus()
```

### Log Analysis

```bash
# MongoDB log location (default)
# /var/log/mongodb/mongod.log

# Common log patterns to watch:
# - "Slow query" - Performance issues
# - "Connection accepted/ended" - Connection issues
# - "replSet" - Replication events
# - "WiredTiger" - Storage engine issues

# Enable verbose logging
db.setLogLevel(1)  # 0-5, higher = more verbose
db.setLogLevel(2, "query")  # Specific component

# Get log
db.adminCommand({ getLog: "global" })
```

### Performance Issues

```javascript
// 1. Slow queries
// - Check explain output
db.collection.find({ field: value }).explain("executionStats")
// - Add appropriate indexes
// - Reduce returned fields

// 2. High CPU
// - Check for collection scans
db.currentOp({ "command.filter": { $exists: true } })
// - Check for JavaScript execution ($where)
// - Check aggregation pipelines

// 3. High memory
// - Check working set size vs RAM
db.serverStatus().wiredTiger.cache
// - Check for large sorts without index
// - Reduce batch sizes

// 4. Slow writes
// - Check write concern
// - Check for index overhead
// - Consider bulk operations

// 5. Replication lag
rs.printSecondaryReplicationInfo()
// - Check oplog size
// - Check network latency
// - Check secondary resources
```

---

## Quick Reference

### mongosh Commands

```javascript
// Database
use dbname              // Switch database
show dbs                // List databases
db.dropDatabase()       // Drop current database

// Collections
show collections        // List collections
db.createCollection()   // Create collection
db.collection.drop()    // Drop collection

// CRUD
db.col.insertOne()      // Insert one
db.col.insertMany()     // Insert many
db.col.find()           // Query
db.col.findOne()        // Query one
db.col.updateOne()      // Update one
db.col.updateMany()     // Update many
db.col.deleteOne()      // Delete one
db.col.deleteMany()     // Delete many

// Indexes
db.col.createIndex()    // Create index
db.col.getIndexes()     // List indexes
db.col.dropIndex()      // Drop index

// Aggregation
db.col.aggregate([])    // Aggregation pipeline
db.col.countDocuments() // Count documents
db.col.distinct()       // Distinct values
```

### Common Operators

```javascript
// Comparison
$eq, $ne, $gt, $gte, $lt, $lte, $in, $nin

// Logical
$and, $or, $not, $nor

// Element
$exists, $type

// Array
$all, $elemMatch, $size

// Update
$set, $unset, $inc, $push, $pull, $addToSet

// Aggregation stages
$match, $group, $project, $sort, $limit, $skip, $lookup, $unwind
```

---

## Resources

### Official Documentation
- [MongoDB Manual](https://docs.mongodb.com/manual/)
- [MongoDB University](https://university.mongodb.com/)
- [MongoDB Drivers](https://docs.mongodb.com/drivers/)

### Tools
- **MongoDB Compass** - Official GUI
- **mongosh** - MongoDB Shell
- **Studio 3T** - Third-party GUI
- **Robo 3T** - Lightweight GUI

### Community
- [MongoDB Community Forums](https://www.mongodb.com/community/forums/)
- [Stack Overflow - mongodb tag](https://stackoverflow.com/questions/tagged/mongodb)
- [MongoDB Blog](https://www.mongodb.com/blog)

---

*Last Updated: February 2026*
