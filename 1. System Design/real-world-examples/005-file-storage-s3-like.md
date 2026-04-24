# Complete System Design: File/Object Storage System (S3-like)

> **Complexity Level:** Advanced  
> **Estimated Time:** 60-90 minutes in interview  
> **Real-World Examples:** Amazon S3, Google Cloud Storage, Azure Blob Storage, MinIO

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Storage Engine](#7-deep-dive-storage-engine)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a scalable object storage system like Amazon S3 that can store and retrieve files of any size with high durability."

**Interviewer's Perspective:**  
This is a premier system design question for Staff+ level roles. It assesses:
- Distributed storage fundamentals (chunking, replication, erasure coding)
- Data durability guarantees and mathematical reasoning
- Consistency models (strong vs. eventual) in a distributed system
- Large-scale data management (petabyte-level)
- Multi-tenancy and access control
- Understanding of storage trade-offs (cost vs. performance vs. durability)

**Why This Problem Is Hard:**  
Unlike typical CRUD systems, an object storage system must handle objects ranging from 1 byte to 5 TB, guarantee 11 nines of durability (losing at most 1 object per 10 billion per year), and serve 100K+ requests/sec—all while keeping costs low enough to store petabytes of data economically.

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I start designing, I'd like to clarify the requirements. Object storage systems can vary widely in their capabilities, so I want to make sure I'm solving the right problem."

**Interviewer:** "Go ahead."

### 2.1 Functional Requirements

**Candidate:** "Let me walk through the functional requirements:
1. Can users upload objects of any size? Is there a maximum?
2. Should we support bucket management—creating and deleting namespaces?
3. Do we need object versioning?
4. What kind of access control is required?
5. Do we need object listing with prefix-based filtering?
6. Should deletes be immediate or soft-deletes?"

**Interviewer:** "Good questions. Here's what we need:
- Upload objects up to 5 TB in size
- Download objects by bucket + key
- Delete objects
- List objects within a bucket with prefix filtering
- Bucket management: create and delete buckets
- Support object versioning (optional, toggled per bucket)
- Access control: public/private buckets, per-object ACLs"

**Candidate:** "Got it. So the core API surface is:
1. ✅ `PUT /{bucket}/{key}` — Upload an object (up to 5 TB)
2. ✅ `GET /{bucket}/{key}` — Download an object
3. ✅ `DELETE /{bucket}/{key}` — Delete an object
4. ✅ `GET /{bucket}?prefix=&marker=&max-keys=` — List objects
5. ✅ `PUT /{bucket}` — Create a bucket
6. ✅ `DELETE /{bucket}` — Delete a bucket (must be empty)
7. ✅ Object versioning (per-bucket toggle)
8. ✅ ACL-based access control"

### 2.2 Non-Functional Requirements

**Candidate:** "Now for non-functional requirements:
1. What durability guarantee do we need?
2. What's the availability target?
3. What's the consistency model—strong or eventual?
4. What's the expected request rate and storage volume?
5. What latency targets should I design for?"

**Interviewer:** "We need:
- 99.999999999% durability (11 nines)—this is non-negotiable
- 99.99% availability (roughly 52 minutes downtime per year)
- Strong read-after-write consistency for PUTs
- Eventual consistency for listings and deletes
- Low latency for small objects (<100ms for objects under 1 MB)
- High throughput for large objects (saturate network bandwidth)"

**Candidate:** "Let me also confirm the scale:
- Total objects: ~1 billion
- Total storage: ~100 PB across data centers
- Request rate: ~100K requests/sec at peak

That's a massive system. The durability requirement alone means we can lose at most 1 object per 10 billion objects per year. I'll need erasure coding, not just replication, to hit that economically."

**Interviewer:** "Exactly. That's the kind of reasoning I'm looking for."

### 2.3 Out of Scope (for initial design)

**Candidate:** "I'll defer these for the advanced features section:
- Server-side encryption
- Cross-region replication
- Lifecycle policies
- Event notifications
- Object lock / WORM compliance"

---

## 3. Scale Estimation

### 3.1 Traffic Estimates

| Metric | Value |
|---|---|
| Total requests/sec | 100,000 |
| Read requests (60%) | 60,000/sec |
| Write requests (30%) | 30,000/sec |
| Delete requests (10%) | 10,000/sec |
| Average object size | 100 KB |
| Daily new objects | 30,000 × 86,400 ≈ 2.6 billion/day |

### 3.2 Storage Estimates

| Metric | Value |
|---|---|
| Total objects | 1 billion |
| Total raw storage | 100 PB |
| Average object size | 100 KB (median much smaller, long tail of large objects) |
| Metadata per object | ~1 KB |
| Total metadata | 1 billion × 1 KB = 1 TB |
| Erasure coding overhead (1.5×) | 100 PB × 1.5 = 150 PB physical storage |

### 3.3 Bandwidth Estimates

| Metric | Calculation | Value |
|---|---|---|
| Read bandwidth | 60,000 × 100 KB | 6 GB/sec |
| Write bandwidth | 30,000 × 100 KB | 3 GB/sec |
| Total bandwidth | | 9 GB/sec ≈ 72 Gbps |
| Network (with overhead) | 72 Gbps × 1.2 | ~86 Gbps |

### 3.4 Infrastructure Estimates

```
Storage nodes:
  - Assume 100 TB usable per node (10 × 16TB HDDs with RAID)
  - 150 PB / 100 TB = 1,500 storage nodes

Metadata nodes:
  - 1 TB metadata fits in ~20 nodes with replication factor 3
  - Need fast SSDs for low-latency lookups

API servers:
  - Assume 2,000 req/sec per server
  - 100,000 / 2,000 = 50 API servers (with headroom: 80)
```

---

## 4. High-Level Design

### 4.1 Architecture Overview

```
┌──────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                     │
│          SDK (aws-sdk)  |  CLI (s3cmd)  |  REST API  |  Web Console     │
└───────────────────────────────┬──────────────────────────────────────────┘
                                │
                        ┌───────▼────────┐
                        │   DNS / Load   │
                        │   Balancer     │  (s3.example.com)
                        └───────┬────────┘
                                │
                    ┌───────────▼───────────┐
                    │     API Gateway        │  Auth, Rate Limiting,
                    │  (Stateless, N nodes)  │  Request Routing
                    └─────┬─────────┬───────┘
                          │         │
              ┌───────────▼──┐  ┌───▼──────────────┐
              │  Metadata    │  │  Data Service     │
              │  Service     │  │  (Chunk Routing)  │
              └──────┬───────┘  └───┬───────────────┘
                     │              │
          ┌──────────▼──────────┐   │
          │  Metadata Store     │   │
          │  (Distributed KV)   │   │
          │  FoundationDB /     │   │
          │  DynamoDB-like      │   │
          └─────────────────────┘   │
                                    │
              ┌─────────────────────▼──────────────────────┐
              │            Placement Service                │
              │  (Consistent Hashing + Virtual Nodes)       │
              │  Decides which data nodes store each chunk  │
              └────────┬───────────────┬───────────────────┘
                       │               │
        ┌──────────────▼──┐    ┌───────▼──────────────┐
        │  Data Nodes      │    │  Data Nodes          │
        │  (Zone A)        │    │  (Zone B)            │
        │  ┌────┐ ┌────┐  │    │  ┌────┐ ┌────┐      │
        │  │HDD │ │HDD │  │    │  │HDD │ │HDD │      │
        │  │Farm│ │Farm│  │    │  │Farm│ │Farm│      │
        │  └────┘ └────┘  │    │  └────┘ └────┘      │
        └─────────────────┘    └──────────────────────┘

Background Services:
┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐
│ Replication Mgr   │  │ Garbage Collector  │  │ Integrity Checker │
│ (repair under-    │  │ (delete orphan     │  │ (scrub all data   │
│  replicated data) │  │  chunks, expired)  │  │  periodically)    │
└───────────────────┘  └───────────────────┘  └───────────────────┘
```

### 4.2 API Design

```
# Object Operations
PUT    /{bucket}/{key}                    # Upload object
GET    /{bucket}/{key}                    # Download object
DELETE /{bucket}/{key}                    # Delete object
HEAD   /{bucket}/{key}                    # Get object metadata

# Bucket Operations
PUT    /{bucket}                          # Create bucket
DELETE /{bucket}                          # Delete bucket (must be empty)
GET    /{bucket}?prefix=&marker=&max-keys= # List objects

# Multipart Upload
POST   /{bucket}/{key}?uploads           # Initiate multipart upload
PUT    /{bucket}/{key}?partNumber=N&uploadId=X  # Upload part
POST   /{bucket}/{key}?uploadId=X        # Complete multipart upload
DELETE /{bucket}/{key}?uploadId=X        # Abort multipart upload
```

### 4.3 Data Flow: Object Upload

```
Client                API GW          Metadata Svc       Placement Svc       Data Nodes
  │                     │                  │                   │                  │
  │─── PUT /b/key ─────▶│                  │                   │                  │
  │                     │── auth + validate▶│                   │                  │
  │                     │                  │── allocate obj ID─▶│                  │
  │                     │                  │                   │── select nodes ──▶│
  │                     │                  │◀─ node list ──────│                  │
  │                     │◀─ upload targets─│                   │                  │
  │                     │                  │                   │                  │
  │─── stream data ────▶│─── chunk + encode ───────────────────────────────────▶│
  │                     │                  │                   │    (parallel to  │
  │                     │                  │                   │     N nodes)     │
  │                     │◀────────────────── ack from all nodes ───────────────│
  │                     │── commit metadata▶│                   │                  │
  │                     │                  │── write metadata ─▶│  (to KV store)  │
  │◀── 200 OK ─────────│                  │                   │                  │
```

### 4.4 Data Flow: Object Download

```
Client                API GW          Metadata Svc       Data Nodes
  │                     │                  │                  │
  │─── GET /b/key ─────▶│                  │                  │
  │                     │── lookup metadata▶│                  │
  │                     │◀── chunk map ────│                  │
  │                     │                  │                  │
  │                     │── parallel fetch chunks ───────────▶│
  │                     │◀── chunk data (streamed) ──────────│
  │                     │── decode + reassemble                │
  │◀── stream response─│                  │                  │
```

---

## 5. Deep Dive: Core Components

### 5.1 API Gateway

**Candidate:** "The API gateway is the single entry point for all client requests. It handles several cross-cutting concerns."

```
┌─────────────────────────────────────────────────┐
│                  API Gateway                     │
│                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │
│  │  Auth &   │  │  Rate    │  │  Request     │  │
│  │  ACL      │  │  Limiter │  │  Router      │  │
│  └──────────┘  └──────────┘  └──────────────┘  │
│                                                  │
│  ┌──────────┐  ┌──────────┐  ┌──────────────┐  │
│  │ Pre-sign │  │ Request  │  │  Multipart   │  │
│  │ URL Gen  │  │ Logging  │  │  Coordinator │  │
│  └──────────┘  └──────────┘  └──────────────┘  │
└─────────────────────────────────────────────────┘
```

**Key Responsibilities:**

1. **Authentication:** Validate HMAC-signed requests (AWS Signature V4-style)
2. **Authorization:** Check bucket and object ACLs
3. **Rate Limiting:** Per-account, per-bucket throttling
4. **Pre-signed URLs:** Generate time-limited signed URLs for direct access
5. **Request Routing:** Route to metadata or data service based on operation

```javascript
// Pre-signed URL generation (JavaScript)
const crypto = require('crypto');

function generatePresignedUrl(bucket, key, expireSeconds, secretKey) {
  const expires = Math.floor(Date.now() / 1000) + expireSeconds;
  const stringToSign = `GET\n/${bucket}/${key}\n${expires}`;
  const signature = crypto
    .createHmac('sha256', secretKey)
    .update(stringToSign)
    .digest('base64url');

  return `https://s3.example.com/${bucket}/${key}` +
    `?expires=${expires}&signature=${signature}`;
}

// Signature verification middleware
function verifySignature(req, res, next) {
  const { expires, signature } = req.query;
  if (Date.now() / 1000 > parseInt(expires)) {
    return res.status(403).json({ error: 'URL expired' });
  }

  const stringToSign = `${req.method}\n${req.path}\n${expires}`;
  const expectedSig = crypto
    .createHmac('sha256', getSecretKey(req))
    .update(stringToSign)
    .digest('base64url');

  if (signature !== expectedSig) {
    return res.status(403).json({ error: 'Invalid signature' });
  }
  next();
}
```

### 5.2 Metadata Service

**Candidate:** "The metadata service is the brain of the system. It maps every bucket + key to the physical locations of the object's chunks."

```
┌──────────────────────────────────────────────┐
│              Metadata Service                 │
│                                               │
│  ┌───────────────┐  ┌────────────────────┐   │
│  │  Namespace     │  │  Object Lookup     │   │
│  │  Manager       │  │  (bucket + key     │   │
│  │  (buckets)     │  │   → chunk map)     │   │
│  └───────────────┘  └────────────────────┘   │
│                                               │
│  ┌───────────────┐  ┌────────────────────┐   │
│  │  Version       │  │  Listing Index     │   │
│  │  Manager       │  │  (prefix queries)  │   │
│  └───────────────┘  └────────────────────┘   │
│                                               │
│           ┌──────────────────┐                │
│           │ Distributed KV   │                │
│           │ Store Backend    │                │
│           └──────────────────┘                │
└──────────────────────────────────────────────┘
```

**Key Operations:**

```python
# Metadata service — object lookup (Python pseudocode)

class MetadataService:
    def __init__(self, kv_store):
        self.kv = kv_store  # FoundationDB / DynamoDB client

    def put_object_metadata(self, bucket, key, metadata):
        """Atomic write of object metadata after data is committed."""
        composite_key = f"{bucket}/{key}"

        record = {
            "bucket": bucket,
            "key": key,
            "version_id": metadata["version_id"],
            "size": metadata["size"],
            "content_type": metadata["content_type"],
            "checksum_sha256": metadata["checksum"],
            "created_at": metadata["created_at"],
            "storage_class": metadata.get("storage_class", "STANDARD"),
            "chunk_locations": metadata["chunk_locations"],
            # Each chunk: { chunk_id, nodes: [node1, node2, ...], offset, length }
        }

        # Atomic compare-and-swap for strong consistency
        self.kv.atomic_put(composite_key, record)

    def get_object_metadata(self, bucket, key, version_id=None):
        """Lookup object metadata. Returns chunk locations for data retrieval."""
        composite_key = f"{bucket}/{key}"

        if version_id:
            composite_key += f"#{version_id}"

        record = self.kv.get(composite_key)
        if not record:
            raise ObjectNotFoundError(bucket, key)

        return record

    def list_objects(self, bucket, prefix="", marker="", max_keys=1000):
        """Range scan for listing objects. Eventually consistent."""
        start_key = f"{bucket}/{prefix}{marker}"
        end_key = f"{bucket}/{prefix}\xff"

        results = self.kv.range_scan(
            start=start_key,
            end=end_key,
            limit=max_keys + 1
        )

        is_truncated = len(results) > max_keys
        objects = results[:max_keys]
        next_marker = objects[-1]["key"] if is_truncated else None

        return {
            "objects": objects,
            "is_truncated": is_truncated,
            "next_marker": next_marker,
        }
```

### 5.3 Data Service

**Candidate:** "The data service handles the actual bytes. It's responsible for chunking objects, encoding them, and routing chunks to the right data nodes."

**Responsibilities:**
- Split incoming objects into fixed-size chunks (64 MB)
- Apply erasure coding to each chunk group
- Stream chunks to selected data nodes in parallel
- On reads, fetch chunks in parallel and reassemble

```javascript
// Data service — chunking and upload coordination
const CHUNK_SIZE = 64 * 1024 * 1024; // 64 MB

async function uploadObject(stream, objectSize, placementPlan) {
  const chunks = [];
  let chunkIndex = 0;
  let offset = 0;

  for await (const chunkData of splitStream(stream, CHUNK_SIZE)) {
    const chunkId = generateChunkId();
    const checksum = computeSHA256(chunkData);
    const targetNodes = placementPlan.getNodesForChunk(chunkIndex);

    // Erasure code: 8 data shards + 4 parity shards
    const shards = erasureEncode(chunkData, { dataShards: 8, parityShards: 4 });

    // Write shards to nodes in parallel
    const writePromises = shards.map((shard, i) =>
      targetNodes[i].writeShard(chunkId, i, shard, checksum)
    );
    await Promise.all(writePromises);

    chunks.push({
      chunk_id: chunkId,
      index: chunkIndex,
      offset: offset,
      length: chunkData.length,
      checksum: checksum,
      shard_locations: targetNodes.map((node, i) => ({
        node_id: node.id,
        shard_index: i,
      })),
    });

    offset += chunkData.length;
    chunkIndex++;
  }

  return chunks;
}
```

### 5.4 Placement Service

**Candidate:** "The placement service decides where to store each chunk. It uses consistent hashing with virtual nodes and respects failure domain constraints."

```
Placement Constraints:
  1. Shards of the same chunk must be on different physical servers
  2. At least 2 failure zones (racks / availability zones) must be used
  3. Nodes with more free capacity get more virtual nodes
  4. Respect storage class requirements (SSD vs HDD)

Hash Ring (simplified):
  ┌────────────────────────────────────────┐
  │                                        │
  │     Node-A(v1)     Node-B(v1)          │
  │        ●───────────────●               │
  │       /                 \              │
  │      /                   \             │
  │  Node-D(v2)●           ●Node-C(v1)    │
  │      \                   /             │
  │       \                 /              │
  │        ●───────────────●               │
  │     Node-A(v2)     Node-B(v2)          │
  │                                        │
  └────────────────────────────────────────┘
```

```python
# Placement service — consistent hashing with zone awareness
import hashlib
from collections import defaultdict

class PlacementService:
    def __init__(self, nodes, virtual_nodes_per_node=150):
        self.ring = {}
        self.sorted_keys = []
        self.node_zone_map = {}

        for node in nodes:
            self.node_zone_map[node.id] = node.zone
            for i in range(virtual_nodes_per_node):
                vnode_key = f"{node.id}:vnode:{i}"
                hash_val = self._hash(vnode_key)
                self.ring[hash_val] = node
                self.sorted_keys.append(hash_val)

        self.sorted_keys.sort()

    def _hash(self, key):
        return int(hashlib.sha256(key.encode()).hexdigest(), 16)

    def get_nodes_for_chunk(self, chunk_id, num_shards=12):
        """Select nodes for all shards, ensuring zone diversity."""
        hash_val = self._hash(chunk_id)
        selected = []
        seen_nodes = set()
        zones_used = set()
        idx = self._find_start_index(hash_val)

        while len(selected) < num_shards:
            node = self.ring[self.sorted_keys[idx % len(self.sorted_keys)]]

            if node.id not in seen_nodes:
                zone = self.node_zone_map[node.id]
                # Enforce: no more than ceil(num_shards/num_zones) per zone
                zone_count = sum(
                    1 for n in selected
                    if self.node_zone_map[n.id] == zone
                )
                max_per_zone = (num_shards + 2) // 3  # 3 zones assumed
                if zone_count < max_per_zone:
                    selected.append(node)
                    seen_nodes.add(node.id)
                    zones_used.add(zone)

            idx += 1

        return selected

    def _find_start_index(self, hash_val):
        import bisect
        idx = bisect.bisect_left(self.sorted_keys, hash_val)
        return idx % len(self.sorted_keys)
```

---

## 6. Deep Dive: Database Design

### 6.1 Metadata Schema

**Candidate:** "The metadata layer must store object metadata, bucket configuration, and support fast key lookups and prefix-based listings."

#### Object Metadata

```sql
-- Logical schema (stored in distributed KV store as serialized records)

-- Primary key: (bucket_name, object_key, version_id)
CREATE TABLE object_metadata (
    bucket_name     VARCHAR(63)     NOT NULL,
    object_key      VARCHAR(1024)   NOT NULL,
    version_id      VARCHAR(64)     NOT NULL DEFAULT 'null',
    
    -- Object properties
    size            BIGINT          NOT NULL,
    content_type    VARCHAR(256),
    content_encoding VARCHAR(64),
    checksum_sha256 CHAR(64)        NOT NULL,
    etag            CHAR(32)        NOT NULL,     -- MD5 or multipart etag
    
    -- Timestamps
    created_at      TIMESTAMP       NOT NULL,
    last_modified   TIMESTAMP       NOT NULL,
    delete_marker   BOOLEAN         DEFAULT FALSE,
    
    -- Storage info
    storage_class   VARCHAR(32)     DEFAULT 'STANDARD',
    -- STANDARD | INFREQUENT_ACCESS | ARCHIVE | DEEP_ARCHIVE
    
    -- Chunk locations (stored as JSON array)
    -- Each entry: { chunk_id, shard_locations: [{node_id, shard_idx}], size, checksum }
    chunk_manifest  JSONB           NOT NULL,
    
    -- Access control
    owner_id        VARCHAR(64)     NOT NULL,
    acl             JSONB,
    
    -- User metadata (x-amz-meta-*)
    user_metadata   JSONB,
    
    PRIMARY KEY (bucket_name, object_key, version_id)
);

-- Index for listing objects by prefix
CREATE INDEX idx_listing ON object_metadata (bucket_name, object_key);
```

#### Bucket Metadata

```sql
CREATE TABLE bucket_metadata (
    bucket_name         VARCHAR(63)     PRIMARY KEY,
    owner_id            VARCHAR(64)     NOT NULL,
    region              VARCHAR(32)     NOT NULL,
    created_at          TIMESTAMP       NOT NULL,
    
    -- Configuration
    versioning_status   VARCHAR(16)     DEFAULT 'Disabled',
    -- Disabled | Enabled | Suspended
    
    acl                 JSONB           NOT NULL DEFAULT '{"public": false}',
    
    -- Lifecycle rules (JSON array of rules)
    lifecycle_rules     JSONB           DEFAULT '[]',
    
    -- Encryption
    default_encryption  JSONB,
    -- { "algorithm": "AES256", "kms_key_id": "..." }
    
    -- Replication
    replication_config  JSONB,
    
    -- Notifications
    notification_config JSONB,
    
    -- Quotas
    max_objects         BIGINT,
    max_storage_bytes   BIGINT
);
```

#### Multipart Upload Tracking

```sql
CREATE TABLE multipart_uploads (
    bucket_name     VARCHAR(63)     NOT NULL,
    object_key      VARCHAR(1024)   NOT NULL,
    upload_id       VARCHAR(64)     NOT NULL,
    
    initiated_at    TIMESTAMP       NOT NULL,
    owner_id        VARCHAR(64)     NOT NULL,
    storage_class   VARCHAR(32)     DEFAULT 'STANDARD',
    
    -- Status: IN_PROGRESS | COMPLETED | ABORTED
    status          VARCHAR(16)     DEFAULT 'IN_PROGRESS',
    
    PRIMARY KEY (upload_id)
);

CREATE TABLE multipart_parts (
    upload_id       VARCHAR(64)     NOT NULL,
    part_number     INTEGER         NOT NULL,     -- 1 to 10,000
    size            BIGINT          NOT NULL,
    checksum_sha256 CHAR(64)        NOT NULL,
    etag            CHAR(32)        NOT NULL,
    chunk_manifest  JSONB           NOT NULL,
    uploaded_at     TIMESTAMP       NOT NULL,
    
    PRIMARY KEY (upload_id, part_number)
);
```

### 6.2 Why Distributed KV Store over SQL?

**Candidate:** "For the metadata store, I'd choose a distributed KV store like FoundationDB or a DynamoDB-like system over traditional SQL. Here's why:"

| Criteria | Distributed KV Store | Traditional SQL (PostgreSQL) |
|---|---|---|
| **Scale** | Horizontally scalable to 1000s of nodes | Vertical scaling limits; sharding is complex |
| **Key-range scans** | Native support (for prefix listings) | B-tree index works but less scalable |
| **Latency** | Single-digit ms for point lookups | Comparable, but harder to maintain at scale |
| **Availability** | Built-in replication + failover | Requires external HA setup (Patroni, etc.) |
| **Schema flexibility** | Schema-less; easy to evolve | Schema migrations on 1 billion rows are costly |
| **Transactions** | FoundationDB: full ACID; DynamoDB: conditional writes | Full ACID |
| **Operational cost** | Higher learning curve, but scales better | Well-understood, but ops burden at 100PB |

**Interviewer:** "Would you ever use SQL for any part of this?"

**Candidate:** "Yes—for the account/billing service and bucket-level configuration. Those are low-volume, relational data that benefits from SQL's query flexibility. But for the object metadata that must handle 100K lookups/sec across a billion keys, a distributed KV store is the right choice."

---

## 7. Deep Dive: Storage Engine

> This is the most critical section for the interview. The storage engine determines durability, cost, and performance.

### 7.1 Object Chunking

**Candidate:** "Large objects are split into fixed-size chunks of 64 MB. Small objects (under the chunk size) are stored as a single chunk."

```
Object (350 MB)
┌─────────────────────────────────────────────────────┐
│                  Original Object                     │
└─────────────────────────────────────────────────────┘
                        │
                  Split into chunks
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
  ┌──────────┐   ┌──────────┐   ┌──────────┐
  │ Chunk 0  │   │ Chunk 1  │   │ Chunk 2  │
  │  64 MB   │   │  64 MB   │   │  64 MB   │
  └──────────┘   └──────────┘   └──────────┘
        ▼               ▼               ▼
  ┌──────────┐   ┌──────────┐   ┌──────────┐
  │ Chunk 3  │   │ Chunk 4  │   │ Chunk 5  │
  │  64 MB   │   │  64 MB   │   │  30 MB   │  (last chunk, partial)
  └──────────┘   └──────────┘   └──────────┘

Why 64 MB?
  - Small enough to parallelize reads/writes across nodes
  - Large enough to amortize metadata overhead
  - Matches typical disk I/O sweet spot for sequential writes
  - Same chunk size used by GFS, HDFS (adjustable per workload)
```

### 7.2 Erasure Coding vs. Replication

**Candidate:** "This is the most important trade-off in the entire system. Let me compare both approaches."

#### Three-Way Replication

```
Chunk → Copy 1 (Node A) + Copy 2 (Node B) + Copy 3 (Node C)

Storage overhead: 3× (store 3 copies of everything)
For 100 PB raw data → 300 PB physical storage

Durability (simplified):
  P(single disk failure/year) ≈ 2%
  P(losing all 3 copies) = (0.02)^3 = 8 × 10^-6
  That's only 5 nines of durability — NOT ENOUGH for our 11-nine target
```

#### Erasure Coding (Reed-Solomon 8+4)

```
                         Chunk (64 MB)
                              │
                    ┌─────────▼─────────┐
                    │  Erasure Encoder   │
                    │  Reed-Solomon      │
                    │  (8 data + 4 par)  │
                    └─────────┬─────────┘
                              │
    ┌──────┬──────┬──────┬───┼───┬──────┬──────┬──────┐
    ▼      ▼      ▼      ▼   ▼   ▼      ▼      ▼      │
  ┌────┐┌────┐┌────┐┌────┐┌────┐┌────┐┌────┐┌────┐    │
  │ D0 ││ D1 ││ D2 ││ D3 ││ D4 ││ D5 ││ D6 ││ D7 │    │
  │8MB ││8MB ││8MB ││8MB ││8MB ││8MB ││8MB ││8MB │    │
  └────┘└────┘└────┘└────┘└────┘└────┘└────┘└────┘    │
  ┌────┐┌────┐┌────┐┌────┐                             │
  │ P0 ││ P1 ││ P2 ││ P3 │  ← Parity shards           │
  │8MB ││8MB ││8MB ││8MB │                             │
  └────┘└────┘└────┘└────┘                             │
                                                        │
  Total: 12 shards × 8 MB = 96 MB for 64 MB data       │
  Storage overhead: 1.5× (vs 3× for replication)        │
                                                        ▼
  Can tolerate loss of ANY 4 shards out of 12
```

**Durability Math (11 Nines):**

```
Reed-Solomon (8, 4) durability calculation:

Given:
  - Annual Disk Failure Rate (AFR) = 2%
  - p = probability a single shard is lost = 0.02
  - n = 12 total shards
  - k = 8 data shards (need any 8 of 12 to reconstruct)
  - Can tolerate t = 4 failures

P(data loss) = P(5 or more shards fail simultaneously before repair)

Using binomial distribution:
  P(data loss) = Σ(i=5 to 12) C(12,i) × p^i × (1-p)^(12-i)

  P(5 failures) = C(12,5) × (0.02)^5 × (0.98)^7
                 = 792 × 3.2×10^-9 × 0.868
                 = 2.2 × 10^-6

  P(6+ failures) is negligible (< 10^-10)

But with proactive repair (detect + re-replicate within hours):
  Effective p during repair window ≈ 0.0001 (not full year)

  P(data loss during repair) ≈ C(12,5) × (0.0001)^5 × (0.9999)^7
                              ≈ 7.92 × 10^-17

  That's approximately 16 nines — exceeding our 11-nine target.
```

**Comparison Summary:**

| Metric | 3× Replication | Erasure Coding (8+4) |
|---|---|---|
| Storage overhead | 3.0× | 1.5× |
| Durability | ~5 nines | ~16 nines (with repair) |
| Read latency | Low (any copy) | Higher (decode needed) |
| Write latency | Low (parallel writes) | Higher (encoding overhead) |
| Repair bandwidth | Full chunk copy | 1 shard + decode |
| Cost at 100 PB | 300 PB physical | 150 PB physical |
| Savings at $0.02/GB/mo | $6M/mo | $3M/mo |

**Candidate:** "For our system, I'd use erasure coding for the standard storage class and replication only for the hot/frequently-accessed tier where read latency matters most."

### 7.3 Data Placement with Consistent Hashing

```
Traditional Hashing:
  node = hash(chunk_id) % N     ← Adding a node remaps ~all chunks!

Consistent Hashing:
  node = first_node_clockwise(hash(chunk_id))  ← Adding a node remaps ~1/N chunks

Virtual Nodes:
  Each physical node → 100-200 virtual nodes on the ring
  Advantages:
    - Better load distribution
    - Weight nodes by capacity (more vnodes = more data)
    - Smooth rebalancing when nodes join/leave

Zone-Aware Placement:
  ┌─────────────────────────────────────────────────┐
  │                   Hash Ring                      │
  │                                                  │
  │   Zone A: Nodes 1-500     (rack A, power A)     │
  │   Zone B: Nodes 501-1000  (rack B, power B)     │
  │   Zone C: Nodes 1001-1500 (rack C, power C)     │
  │                                                  │
  │   Rule: 12 shards spread across ≥ 2 zones       │
  │         (ideally 4 per zone for 3 zones)         │
  └─────────────────────────────────────────────────┘
```

### 7.4 Write Path (Detailed)

```python
# Complete write path implementation

class WriteCoordinator:
    def __init__(self, metadata_svc, placement_svc, erasure_coder):
        self.metadata = metadata_svc
        self.placement = placement_svc
        self.ec = erasure_coder

    async def upload_object(self, bucket, key, data_stream, content_type, size):
        # 1. Validate bucket exists and user has write permission
        bucket_meta = await self.metadata.get_bucket(bucket)
        if not bucket_meta:
            raise BucketNotFoundError(bucket)

        # 2. Generate version ID if versioning enabled
        version_id = generate_uuid() if bucket_meta["versioning"] == "Enabled" else "null"

        # 3. Split into chunks and process each
        chunk_manifest = []
        full_checksum = hashlib.sha256()

        async for chunk_data in chunk_stream(data_stream, CHUNK_SIZE):
            full_checksum.update(chunk_data)
            chunk_id = generate_chunk_id()
            chunk_checksum = hashlib.sha256(chunk_data).hexdigest()

            # 4. Erasure encode the chunk
            shards = self.ec.encode(chunk_data, data_shards=8, parity_shards=4)

            # 5. Get placement plan
            target_nodes = self.placement.get_nodes_for_chunk(chunk_id, num_shards=12)

            # 6. Write shards to data nodes in parallel
            write_tasks = []
            for shard_idx, (shard, node) in enumerate(zip(shards, target_nodes)):
                shard_checksum = hashlib.sha256(shard).hexdigest()
                write_tasks.append(
                    node.write_shard(chunk_id, shard_idx, shard, shard_checksum)
                )

            results = await asyncio.gather(*write_tasks, return_exceptions=True)

            # 7. Verify enough shards were written (need at least 8 of 12)
            successful = [r for r in results if not isinstance(r, Exception)]
            if len(successful) < 8:
                await self._cleanup_partial_writes(chunk_id, target_nodes)
                raise WriteFailureError("Insufficient shard writes")

            # 8. Record chunk location
            chunk_manifest.append({
                "chunk_id": chunk_id,
                "size": len(chunk_data),
                "checksum": chunk_checksum,
                "shard_locations": [
                    {"node_id": node.id, "shard_idx": idx}
                    for idx, (node, r) in enumerate(zip(target_nodes, results))
                    if not isinstance(r, Exception)
                ],
            })

        # 9. Commit metadata atomically
        object_metadata = {
            "bucket": bucket,
            "key": key,
            "version_id": version_id,
            "size": size,
            "content_type": content_type,
            "checksum": full_checksum.hexdigest(),
            "etag": compute_etag(chunk_manifest),
            "created_at": datetime.utcnow().isoformat(),
            "storage_class": "STANDARD",
            "chunk_manifest": chunk_manifest,
        }
        await self.metadata.put_object_metadata(bucket, key, object_metadata)

        return {"etag": object_metadata["etag"], "version_id": version_id}
```

### 7.5 Read Path (Detailed)

```python
class ReadCoordinator:
    def __init__(self, metadata_svc, erasure_coder):
        self.metadata = metadata_svc
        self.ec = erasure_coder

    async def download_object(self, bucket, key, version_id=None, byte_range=None):
        # 1. Fetch metadata
        meta = await self.metadata.get_object_metadata(bucket, key, version_id)

        # 2. Determine which chunks to fetch
        if byte_range:
            chunks_needed = self._resolve_byte_range(meta["chunk_manifest"], byte_range)
        else:
            chunks_needed = meta["chunk_manifest"]

        # 3. Fetch and decode each chunk
        async def fetch_chunk(chunk_info):
            shard_locations = chunk_info["shard_locations"]

            # Fetch 8 shards in parallel (only need 8 of 12 to decode)
            # Try fastest nodes first (sort by latency history)
            sorted_locs = self._sort_by_latency(shard_locations)
            shards = [None] * 12
            fetched = 0

            async for loc in parallel_fetch(sorted_locs[:10]):  # fetch 10, need 8
                shard_data = await loc["node"].read_shard(
                    chunk_info["chunk_id"], loc["shard_idx"]
                )
                shard_checksum = hashlib.sha256(shard_data).hexdigest()

                if self._verify_shard_checksum(shard_checksum, chunk_info, loc["shard_idx"]):
                    shards[loc["shard_idx"]] = shard_data
                    fetched += 1
                    if fetched >= 8:
                        break

            if fetched < 8:
                raise ReadFailureError("Insufficient shards available")

            # 4. Erasure decode to reconstruct original chunk
            return self.ec.decode(shards, data_shards=8, parity_shards=4)

        # 5. Stream chunks back to client
        for chunk_info in chunks_needed:
            chunk_data = await fetch_chunk(chunk_info)
            yield chunk_data
```

### 7.6 Multipart Upload

**Candidate:** "For objects larger than ~100 MB, we use multipart uploads. This is essential for reliability—if a part fails, only that part needs to be retried."

```
Multipart Upload Flow:
                                                        
  1. Initiate       POST /{bucket}/{key}?uploads
                    → Returns upload_id

  2. Upload Parts   PUT /{bucket}/{key}?partNumber=1&uploadId=X
     (parallel)     PUT /{bucket}/{key}?partNumber=2&uploadId=X
                    PUT /{bucket}/{key}?partNumber=3&uploadId=X
                    ...
                    Each part: 5 MB to 5 GB, up to 10,000 parts

  3. Complete        POST /{bucket}/{key}?uploadId=X
                    Body: [{ partNumber: 1, etag: "..." }, ...]
                    → Stitches parts into final object

  4. (or) Abort      DELETE /{bucket}/{key}?uploadId=X
                    → Cleans up all uploaded parts
```

```javascript
// Multipart upload — client-side implementation
const PART_SIZE = 100 * 1024 * 1024; // 100 MB per part

async function multipartUpload(bucket, key, filePath) {
  const fileSize = fs.statSync(filePath).size;
  const numParts = Math.ceil(fileSize / PART_SIZE);

  // Step 1: Initiate
  const { uploadId } = await apiCall('POST', `/${bucket}/${key}?uploads`);

  // Step 2: Upload parts in parallel (with concurrency limit)
  const partResults = [];
  const concurrency = 4;

  for (let batch = 0; batch < numParts; batch += concurrency) {
    const batchPromises = [];
    for (let i = batch; i < Math.min(batch + concurrency, numParts); i++) {
      const start = i * PART_SIZE;
      const end = Math.min(start + PART_SIZE, fileSize);
      const partStream = fs.createReadStream(filePath, { start, end: end - 1 });

      batchPromises.push(
        uploadPartWithRetry(bucket, key, uploadId, i + 1, partStream, 3)
      );
    }
    const results = await Promise.all(batchPromises);
    partResults.push(...results);
  }

  // Step 3: Complete
  const response = await apiCall('POST', `/${bucket}/${key}?uploadId=${uploadId}`, {
    parts: partResults.map(r => ({
      partNumber: r.partNumber,
      etag: r.etag,
    })),
  });

  return response;
}

async function uploadPartWithRetry(bucket, key, uploadId, partNum, stream, maxRetries) {
  for (let attempt = 0; attempt < maxRetries; attempt++) {
    try {
      const response = await apiCall(
        'PUT',
        `/${bucket}/${key}?partNumber=${partNum}&uploadId=${uploadId}`,
        stream
      );
      return { partNumber: partNum, etag: response.etag };
    } catch (err) {
      if (attempt === maxRetries - 1) throw err;
      await sleep(Math.pow(2, attempt) * 1000); // exponential backoff
    }
  }
}
```

### 7.7 Data Integrity: Checksums at Every Level

```
End-to-End Integrity Chain:

Client                 API Gateway            Data Node             Disk
  │                       │                      │                   │
  │── Content-SHA256 ────▶│                      │                   │
  │   (full object hash)  │                      │                   │
  │                       │── Chunk checksum ────▶│                   │
  │                       │   (per 64MB chunk)    │                   │
  │                       │                      │── Shard checksum ─▶│
  │                       │                      │   (per shard)      │
  │                       │                      │                   │── Disk checksum
  │                       │                      │                   │   (filesystem)
  │                       │                      │                   │
  
  Verification points:
  ┌─────────────────────────────────────────────────────────────────┐
  │ 1. Upload:  Client sends Content-SHA256 header                  │
  │             API gateway verifies after receiving full object     │
  │ 2. Chunking: Each chunk gets SHA-256 before erasure coding      │
  │ 3. Shard:   Each shard gets SHA-256 before writing to node      │
  │ 4. Storage: Data node verifies on write, periodic scrubbing     │
  │ 5. Read:    Shards verified on read, chunk verified after decode │
  │ 6. Return:  Full object checksum verified before sending to     │
  │             client (or client verifies with returned ETag)      │
  └─────────────────────────────────────────────────────────────────┘
```

```python
# Background integrity scrubber
class IntegrityScrubber:
    """Runs on each data node, periodically verifies all stored shards."""

    def __init__(self, storage_engine, metadata_client, repair_queue):
        self.storage = storage_engine
        self.metadata = metadata_client
        self.repair_queue = repair_queue
        self.scrub_interval_days = 30  # Full scrub every 30 days

    async def scrub_all_shards(self):
        """Iterate through all local shards and verify checksums."""
        corrupted_count = 0

        async for shard_info in self.storage.iterate_all_shards():
            stored_data = await self.storage.read_raw(shard_info["shard_path"])
            actual_checksum = hashlib.sha256(stored_data).hexdigest()

            if actual_checksum != shard_info["expected_checksum"]:
                corrupted_count += 1
                await self.repair_queue.enqueue({
                    "chunk_id": shard_info["chunk_id"],
                    "shard_idx": shard_info["shard_idx"],
                    "node_id": self.storage.node_id,
                    "type": "CORRUPTION_DETECTED",
                })

        return {"total_scanned": self.storage.shard_count, "corrupted": corrupted_count}
```

### 7.8 Storage Classes

**Candidate:** "Not all data has the same access pattern. We offer multiple storage classes to optimize cost."

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        Storage Class Tiers                               │
│                                                                          │
│  ┌─────────────┐  ┌─────────────────┐  ┌────────────┐  ┌────────────┐  │
│  │  STANDARD   │  │ INFREQUENT      │  │  ARCHIVE   │  │   DEEP     │  │
│  │             │  │ ACCESS (IA)     │  │            │  │  ARCHIVE   │  │
│  │  Hot data   │  │ Warm data       │  │ Cold data  │  │ Frozen     │  │
│  │  SSD+HDD   │  │ HDD only        │  │ Tape-like  │  │ Offline    │  │
│  │             │  │                 │  │            │  │            │  │
│  │ $0.023/GB  │  │ $0.0125/GB      │  │ $0.004/GB  │  │ $0.00099/  │  │
│  │ per month  │  │ per month       │  │ per month  │  │ GB/month   │  │
│  │             │  │                 │  │            │  │            │  │
│  │ Encoding:  │  │ Encoding:       │  │ Encoding:  │  │ Encoding:  │  │
│  │ EC(8,4)    │  │ EC(6,6)         │  │ EC(4,8)    │  │ EC(4,8)+   │  │
│  │ 1.5× over  │  │ 2.0× overhead   │  │ 3.0× over  │  │ tape copy  │  │
│  │            │  │ Higher durabil. │  │            │  │            │  │
│  │ Retrieval: │  │ Retrieval:      │  │ Retrieval: │  │ Retrieval: │  │
│  │ Instant    │  │ Instant         │  │ 1-5 hours  │  │ 12+ hours  │  │
│  └─────────────┘  └─────────────────┘  └────────────┘  └────────────┘  │
│                                                                          │
│  Lifecycle transition: STANDARD → IA (30d) → ARCHIVE (90d) → DEEP (180d)│
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.9 Garbage Collection

**Candidate:** "Deletes in an object storage system are not immediate. We use a mark-and-sweep approach to safely reclaim storage."

```
Delete Flow:

  1. Client sends DELETE /{bucket}/{key}
  2. Metadata service marks object as deleted (or inserts delete marker if versioned)
  3. Object is immediately invisible to GET/LIST operations
  4. Actual data chunks are NOT deleted yet (safety window)

Garbage Collection (Background):

  ┌──────────────────────────────────────────────────────────┐
  │                  GC Pipeline                              │
  │                                                           │
  │  Phase 1: MARK                                           │
  │  ┌───────────────────────────────────────────────────┐   │
  │  │ Scan metadata for:                                 │   │
  │  │  - Delete markers older than 24 hours              │   │
  │  │  - Expired multipart uploads (> 7 days)           │   │
  │  │  - Overwritten versions past retention window      │   │
  │  │ Output: list of chunk_ids to delete                │   │
  │  └───────────────────────────────────────────────────┘   │
  │                                                           │
  │  Phase 2: VERIFY                                         │
  │  ┌───────────────────────────────────────────────────┐   │
  │  │ For each chunk_id:                                 │   │
  │  │  - Confirm NO live object references this chunk    │   │
  │  │  - Cross-check with metadata store                 │   │
  │  │  - Add to "safe to delete" list                    │   │
  │  └───────────────────────────────────────────────────┘   │
  │                                                           │
  │  Phase 3: SWEEP                                          │
  │  ┌───────────────────────────────────────────────────┐   │
  │  │ For each verified chunk_id:                        │   │
  │  │  - Send delete commands to all data nodes          │   │
  │  │  - Nodes delete individual shards from disk        │   │
  │  │  - Remove chunk metadata from KV store             │   │
  │  └───────────────────────────────────────────────────┘   │
  └──────────────────────────────────────────────────────────┘
```

```python
class GarbageCollector:
    def __init__(self, metadata_svc, data_node_client):
        self.metadata = metadata_svc
        self.data_nodes = data_node_client
        self.safety_window = timedelta(hours=24)

    async def run_gc_cycle(self):
        # Phase 1: Mark — find deletable chunks
        candidates = await self._find_gc_candidates()

        # Phase 2: Verify — double-check no live references
        verified = []
        for chunk_id in candidates:
            refs = await self.metadata.count_references(chunk_id)
            if refs == 0:
                verified.append(chunk_id)

        # Phase 3: Sweep — delete from data nodes
        deleted_count = 0
        for chunk_id in verified:
            shard_locations = await self.metadata.get_shard_locations(chunk_id)
            delete_tasks = [
                self.data_nodes.delete_shard(loc["node_id"], chunk_id, loc["shard_idx"])
                for loc in shard_locations
            ]
            await asyncio.gather(*delete_tasks, return_exceptions=True)
            await self.metadata.remove_chunk_metadata(chunk_id)
            deleted_count += 1

        return {"candidates": len(candidates), "deleted": deleted_count}

    async def _find_gc_candidates(self):
        cutoff = datetime.utcnow() - self.safety_window
        candidates = []

        # Deleted objects past safety window
        async for obj in self.metadata.scan_deleted_objects(before=cutoff):
            for chunk in obj["chunk_manifest"]:
                candidates.append(chunk["chunk_id"])

        # Expired multipart uploads
        async for upload in self.metadata.scan_expired_uploads(max_age_days=7):
            for part in upload["parts"]:
                for chunk in part["chunk_manifest"]:
                    candidates.append(chunk["chunk_id"])

        return candidates
```

---

## 8. Scaling Strategies

### 8.1 Horizontal Scaling of Data Nodes

```
Adding New Nodes:
  
  Before (3 nodes):
  ┌──────┐  ┌──────┐  ┌──────┐
  │Node A│  │Node B│  │Node C│
  │ 90%  │  │ 88%  │  │ 92%  │  ← capacity pressure
  └──────┘  └──────┘  └──────┘

  After adding Node D and Node E:
  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐
  │Node A│  │Node B│  │Node C│  │Node D│  │Node E│
  │ 55%  │  │ 54%  │  │ 57%  │  │ 40%  │  │ 42%  │
  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘

  Rebalancing:
    - New nodes get virtual nodes on the hash ring
    - Background process migrates shards from overloaded nodes
    - Throttled to avoid impacting live traffic (e.g., 100 MB/sec per node)
    - Takes hours/days for large clusters — that's acceptable
```

### 8.2 Metadata Sharding

```
Shard Key: hash(bucket_name + "/" + object_key)

┌─────────────────────────────────────────────────────────────┐
│                   Metadata Sharding                          │
│                                                              │
│  Shard 0: keys [0x0000..., 0x3FFF...]                       │
│  Shard 1: keys [0x4000..., 0x7FFF...]                       │
│  Shard 2: keys [0x8000..., 0xBFFF...]                       │
│  Shard 3: keys [0xC000..., 0xFFFF...]                       │
│                                                              │
│  Each shard: 3-node Raft group for consensus                │
│  Total: ~250 million objects per shard = manageable          │
│                                                              │
│  Hot-spot mitigation:                                        │
│    - If a single bucket has billions of objects,             │
│      the hash distributes its keys across ALL shards         │
│    - Bucket-level operations (list) fan out to all shards    │
└─────────────────────────────────────────────────────────────┘
```

### 8.3 Cross-Region Replication

```
                  Region US-East                Region EU-West
               ┌────────────────┐           ┌────────────────┐
               │  Full cluster  │           │  Full cluster  │
  Write ──────▶│  (primary)     │──async──▶│  (replica)     │──▶ Read
               │                │  replicate│                │
               └────────────────┘           └────────────────┘
                       │                            │
                       └─────── Replication Log ────┘
                         (Kafka-based change stream)

  Replication modes:
    1. Async (default): ~seconds latency, no write penalty
    2. Sync: write only acks after both regions confirm (slower, stronger)
```

### 8.4 CDN Integration

```
  Frequently accessed objects ("hot objects"):

  Client ──▶ CDN Edge (CloudFront / Akamai)
              │
              │ Cache HIT? ──▶ Return immediately (< 10ms)
              │
              │ Cache MISS? ──▶ Fetch from origin (our storage)
                                └── Cache for TTL period
                                    (configurable per bucket)
```

### 8.5 Tiered Storage Migration

```python
# Automated storage class transition
class LifecyclePolicyEngine:
    async def evaluate_bucket_policies(self, bucket):
        rules = await self.metadata.get_lifecycle_rules(bucket)

        for rule in rules:
            matching_objects = await self.metadata.list_objects(
                bucket, prefix=rule.get("prefix", "")
            )
            for obj in matching_objects:
                age_days = (datetime.utcnow() - obj["created_at"]).days

                for transition in rule.get("transitions", []):
                    if age_days >= transition["days"]:
                        current_class = obj["storage_class"]
                        target_class = transition["storage_class"]
                        if self._is_valid_transition(current_class, target_class):
                            await self._transition_object(
                                bucket, obj["key"], target_class
                            )

    def _is_valid_transition(self, current, target):
        hierarchy = ["STANDARD", "INFREQUENT_ACCESS", "ARCHIVE", "DEEP_ARCHIVE"]
        return hierarchy.index(target) > hierarchy.index(current)
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Data Node Failure

```
Scenario: Node-7 (holding 8,000 shards) goes offline

Detection:
  - Heartbeat timeout (30 seconds)
  - Placement service marks node as SUSPECT
  - After 5 minutes with no recovery → marked DEAD

Response:
  ┌───────────────────────────────────────────────────────┐
  │  Replication Manager Repair Flow                       │
  │                                                        │
  │  1. Query metadata: "Which chunks had shards on       │
  │     Node-7?"                                           │
  │  2. For each affected chunk:                           │
  │     a. Check remaining shard count                     │
  │     b. If < threshold (e.g., < 10 of 12):             │
  │        - Read 8 surviving shards from other nodes      │
  │        - Decode → re-encode → create new shards        │
  │        - Place new shards on healthy nodes              │
  │        - Update metadata with new shard locations       │
  │  3. Prioritize: chunks with fewest surviving shards    │
  │     get repaired first                                 │
  └───────────────────────────────────────────────────────┘

Impact:
  - 8 data + 4 parity = 12 shards per chunk
  - Losing 1 node = losing 1 shard per affected chunk
  - 11 shards remain → still readable (need only 8)
  - Repair bandwidth: 8000 shards × 8MB = 64 GB to re-create
  - At 100 MB/sec repair rate: ~10 minutes to fully repair
```

### 9.2 Metadata Service Failure

```
Scenario: Metadata shard leader crashes

┌──────────────────────────────────────────────────────┐
│          Raft Consensus Failover                      │
│                                                       │
│  Normal:                                              │
│    Leader (Node M1) ←── reads/writes                  │
│    Follower (Node M2)                                 │
│    Follower (Node M3)                                 │
│                                                       │
│  M1 crashes:                                          │
│    1. M2 and M3 detect missing heartbeat (150ms)      │
│    2. Election timeout: 150-300ms                     │
│    3. New leader elected (e.g., M2)                   │
│    4. M2 replays uncommitted log entries              │
│    5. Service restored in < 1 second                  │
│                                                       │
│  Impact:                                              │
│    - Reads: brief stall (~500ms), then resume         │
│    - In-flight writes: client retries (idempotent)    │
│    - Zero data loss (Raft guarantees committed logs)   │
└──────────────────────────────────────────────────────┘
```

### 9.3 Corrupt Data Detection and Repair

```
Detection methods:
  1. Read-time verification: checksum mismatch on shard read
  2. Background scrubbing: periodic full-scan (every 30 days)
  3. Client-reported: client detects ETag mismatch

Repair flow:
  ┌─────────────────────────────────────────────┐
  │ Corrupt shard detected on Node-X, Shard #3  │
  │                                              │
  │ 1. Mark shard as BAD in metadata             │
  │ 2. Read 8 good shards from other nodes       │
  │ 3. Decode original chunk data                │
  │ 4. Re-encode shard #3                        │
  │ 5. Write new shard to Node-Y (replacement)   │
  │ 6. Update metadata: shard #3 → Node-Y        │
  │ 7. Delete bad shard from Node-X              │
  │                                              │
  │ Total time: ~seconds for a single shard      │
  └─────────────────────────────────────────────┘
```

### 9.4 Network Partition Between Data Centers

```
Scenario: Network split between DC-East and DC-West

┌─────────────────┐         ╳         ┌─────────────────┐
│    DC-East      │    partition       │    DC-West      │
│  (primary)      │         ╳         │  (secondary)    │
│                 │                    │                 │
│  Can still      │                    │  Read-only mode │
│  serve reads    │                    │  (stale data OK │
│  AND writes     │                    │  for reads)     │
│                 │                    │  Writes queued  │
└─────────────────┘                    └─────────────────┘

Resolution strategy:
  1. Primary DC continues serving all operations
  2. Secondary DC serves reads from local replicas (eventual consistency)
  3. Secondary DC queues writes locally
  4. On partition heal: reconcile queued writes
     - Conflict resolution: last-writer-wins with vector clocks
     - Or: reject secondary writes (fail-safe mode)
```

### 9.5 Partial Upload Failure Recovery

```
Scenario: Client uploading 5 TB file, connection drops at part 3847 of 5000

Recovery:
  1. Multipart upload state is persisted server-side
  2. Client calls LIST PARTS to see which parts were committed
  3. Client resumes from part 3848
  4. Only ~115 GB needs re-uploading (not 5 TB)
  
  Cleanup:
    - If client never completes: auto-abort after 7 days (lifecycle policy)
    - GC reclaims orphaned parts
```

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

```
┌─────────────────────────────────────────────────────────────────────┐
│                    Monitoring Dashboard                               │
│                                                                      │
│  ┌─── Durability ──────────────┐  ┌─── Availability ─────────────┐  │
│  │ Under-replicated chunks: 42 │  │ Success rate: 99.97%         │  │
│  │ Repair queue depth: 1,230   │  │ 5xx errors: 0.03%           │  │
│  │ Scrub errors (30d): 3       │  │ P50 latency: 12ms           │  │
│  │ Data loss events: 0         │  │ P99 latency: 145ms          │  │
│  └─────────────────────────────┘  └──────────────────────────────┘  │
│                                                                      │
│  ┌─── Storage ─────────────────┐  ┌─── Traffic ──────────────────┐  │
│  │ Total: 142 PB / 200 PB cap │  │ GET: 58K/sec                 │  │
│  │ Node utilization:           │  │ PUT: 31K/sec                 │  │
│  │   min: 62%  avg: 78%       │  │ DELETE: 11K/sec              │  │
│  │   max: 91%  ⚠ rebalance    │  │ Bandwidth: 8.7 GB/sec       │  │
│  └─────────────────────────────┘  └──────────────────────────────┘  │
│                                                                      │
│  ┌─── Per-Node Health ─────────────────────────────────────────┐    │
│  │ Node  │ Status │ Capacity │ Shards │ Disk I/O │ Net I/O    │    │
│  │ N-001 │  ✓ OK  │   78%    │ 12,400 │ 240MB/s  │ 380MB/s   │    │
│  │ N-002 │  ✓ OK  │   82%    │ 13,100 │ 210MB/s  │ 350MB/s   │    │
│  │ N-003 │  ⚠ SLOW│   91%    │ 14,800 │ 490MB/s  │ 420MB/s   │    │
│  │ N-007 │  ✗ DOWN│    --    │   --   │   --     │   --      │    │
│  └──────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────┘
```

### 10.2 Alerting Rules

```yaml
# alerting-rules.yml (Prometheus format)
groups:
  - name: storage-critical
    rules:
      - alert: DurabilityAtRisk
        expr: under_replicated_chunks > 1000
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "{{ $value }} chunks below target replication"

      - alert: DataNodeDown
        expr: up{job="data-node"} == 0
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Data node {{ $labels.instance }} is down"

      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.001
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Error rate {{ $value | humanizePercentage }} exceeds 0.1%"

      - alert: StorageCapacityHigh
        expr: node_storage_used_bytes / node_storage_total_bytes > 0.85
        for: 30m
        labels:
          severity: warning
        annotations:
          summary: "Node {{ $labels.node_id }} at {{ $value | humanizePercentage }} capacity"

      - alert: RepairQueueBacklog
        expr: repair_queue_depth > 5000
        for: 15m
        labels:
          severity: warning
        annotations:
          summary: "Repair queue backlog: {{ $value }} chunks pending"

      - alert: HighP99Latency
        expr: histogram_quantile(0.99, rate(request_duration_seconds_bucket[5m])) > 0.5
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "P99 latency {{ $value }}s exceeds 500ms threshold"
```

### 10.3 Audit Logging

```python
# Compliance-grade audit logging for every data access
class AuditLogger:
    def __init__(self, log_sink):
        self.sink = log_sink  # Kafka topic or immutable log store

    def log_access(self, event_type, request_context, result):
        entry = {
            "timestamp": datetime.utcnow().isoformat(),
            "event_type": event_type,  # GET_OBJECT, PUT_OBJECT, DELETE_OBJECT
            "request_id": request_context["request_id"],
            "source_ip": request_context["source_ip"],
            "user_agent": request_context["user_agent"],
            "account_id": request_context["account_id"],
            "bucket": request_context["bucket"],
            "key": request_context.get("key"),
            "http_method": request_context["method"],
            "http_status": result["status_code"],
            "bytes_transferred": result.get("bytes", 0),
            "latency_ms": result["latency_ms"],
            "tls_version": request_context.get("tls_version"),
            "auth_type": request_context.get("auth_type"),  # HMAC, IAM, pre-signed
        }
        self.sink.send("audit-log", entry)
```

---

## 11. Advanced Features

### 11.1 Object Versioning

**Candidate:** "When versioning is enabled on a bucket, every PUT creates a new version instead of overwriting. Deletes insert a delete marker rather than removing the object."

```
Versioning data model:

Key: photos/cat.jpg

  Version History (newest first):
  ┌──────────────────────────────────────────────────┐
  │ v3 │ DELETE MARKER │ 2026-04-24T10:00:00Z        │ ← current "deleted"
  │ v2 │ 1.2 MB       │ 2026-04-23T15:30:00Z        │ ← restorable
  │ v1 │ 800 KB       │ 2026-04-20T09:00:00Z        │ ← restorable
  └──────────────────────────────────────────────────┘

  GET /bucket/photos/cat.jpg          → 404 (delete marker is latest)
  GET /bucket/photos/cat.jpg?versionId=v2  → returns v2 (1.2 MB)
  DELETE /bucket/photos/cat.jpg?versionId=v3 → removes delete marker
  GET /bucket/photos/cat.jpg          → returns v2 (now latest)
```

```javascript
// Version-aware metadata lookup
async function getObject(bucket, key, versionId) {
  if (versionId) {
    const meta = await metadataStore.get(`${bucket}/${key}#${versionId}`);
    if (!meta) throw new NotFoundError();
    if (meta.delete_marker) throw new NotFoundError('Delete marker');
    return meta;
  }

  // Get latest version
  const versions = await metadataStore.rangeScan({
    prefix: `${bucket}/${key}#`,
    order: 'DESC',
    limit: 1,
  });

  if (versions.length === 0) throw new NotFoundError();
  if (versions[0].delete_marker) throw new NotFoundError();
  return versions[0];
}
```

### 11.2 Lifecycle Policies

```json
{
  "rules": [
    {
      "id": "archive-old-logs",
      "prefix": "logs/",
      "enabled": true,
      "transitions": [
        { "days": 30, "storage_class": "INFREQUENT_ACCESS" },
        { "days": 90, "storage_class": "ARCHIVE" },
        { "days": 365, "storage_class": "DEEP_ARCHIVE" }
      ],
      "expiration": {
        "days": 730
      },
      "noncurrent_version_expiration": {
        "days": 30
      },
      "abort_incomplete_multipart_upload": {
        "days_after_initiation": 7
      }
    }
  ]
}
```

### 11.3 Pre-Signed URLs for Temporary Access

```python
# Generate a pre-signed URL allowing a third party to upload directly
def generate_presigned_upload_url(bucket, key, content_type, max_size_bytes,
                                   expire_seconds=3600, secret_key=None):
    expires = int(time.time()) + expire_seconds
    policy = {
        "expiration": datetime.utcfromtimestamp(expires).strftime("%Y-%m-%dT%H:%M:%SZ"),
        "conditions": [
            {"bucket": bucket},
            ["eq", "$key", key],
            ["content-length-range", 1, max_size_bytes],
            {"Content-Type": content_type},
        ],
    }
    policy_b64 = base64.b64encode(json.dumps(policy).encode()).decode()
    signature = hmac.new(
        secret_key.encode(), policy_b64.encode(), hashlib.sha256
    ).hexdigest()

    return {
        "url": f"https://s3.example.com/{bucket}",
        "fields": {
            "key": key,
            "policy": policy_b64,
            "signature": signature,
            "Content-Type": content_type,
        },
    }
```

### 11.4 Event Notifications

```
Event notification system:

  Object Event ──▶ Notification Service ──▶ Fanout
                                              │
                            ┌─────────────────┼────────────────┐
                            ▼                 ▼                ▼
                       ┌─────────┐     ┌───────────┐   ┌────────────┐
                       │  Queue  │     │  Lambda   │   │  Webhook   │
                       │ (SQS)   │     │ Function  │   │  (HTTP)    │
                       └─────────┘     └───────────┘   └────────────┘

Supported events:
  - s3:ObjectCreated:Put
  - s3:ObjectCreated:Post (multipart complete)
  - s3:ObjectCreated:Copy
  - s3:ObjectRemoved:Delete
  - s3:ObjectRemoved:DeleteMarkerCreated
  - s3:LifecycleTransition
```

```javascript
// Event notification dispatcher
class NotificationDispatcher {
  async dispatchEvent(bucket, eventType, objectKey, metadata) {
    const bucketConfig = await this.getBucketNotificationConfig(bucket);
    const event = {
      eventVersion: '2.1',
      eventSource: 'objectstore:s3',
      eventTime: new Date().toISOString(),
      eventName: eventType,
      s3: {
        bucket: { name: bucket },
        object: {
          key: objectKey,
          size: metadata.size,
          eTag: metadata.etag,
          versionId: metadata.version_id,
        },
      },
    };

    const matchingRules = bucketConfig.rules.filter(rule =>
      this.matchesFilter(rule, eventType, objectKey)
    );

    const dispatches = matchingRules.map(rule => {
      switch (rule.destination.type) {
        case 'queue':
          return this.sendToQueue(rule.destination.arn, event);
        case 'function':
          return this.invokeLambda(rule.destination.arn, event);
        case 'webhook':
          return this.callWebhook(rule.destination.url, event);
      }
    });

    await Promise.allSettled(dispatches);
  }
}
```

### 11.5 Server-Side Encryption

```
Three encryption modes:

┌─────────────────────────────────────────────────────────────────┐
│                                                                  │
│  SSE-S3 (Service-managed keys)                                   │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ Object ──▶ AES-256-GCM encrypt ──▶ Store encrypted   │       │
│  │            (per-object key)         data + encrypted  │       │
│  │            wrapped by master key    key envelope      │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                  │
│  SSE-KMS (Customer-managed KMS keys)                             │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ Object ──▶ Request data key from KMS                  │       │
│  │         ──▶ Encrypt with data key                     │       │
│  │         ──▶ Store encrypted data + encrypted data key │       │
│  │ (Customer controls key rotation, audit in KMS)        │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                  │
│  SSE-C (Customer-provided keys)                                  │
│  ┌──────────────────────────────────────────────────────┐       │
│  │ Client sends encryption key in HTTP header            │       │
│  │ Server encrypts, stores data, discards key            │       │
│  │ Client must provide key again on every read           │       │
│  │ (We never store the key — maximum customer control)   │       │
│  └──────────────────────────────────────────────────────┘       │
└─────────────────────────────────────────────────────────────────┘
```

### 11.6 Cross-Region Replication (CRR)

```python
class CrossRegionReplicator:
    def __init__(self, source_region, dest_region, change_stream):
        self.source = source_region
        self.dest = dest_region
        self.stream = change_stream  # Kafka consumer

    async def run(self):
        async for event in self.stream.consume("metadata-changelog"):
            if not self._should_replicate(event):
                continue

            bucket = event["bucket"]
            key = event["key"]

            if event["type"] == "PUT":
                meta = await self.source.metadata.get(bucket, key)
                data_stream = await self.source.data.read_object(meta)
                await self.dest.data.write_object(bucket, key, data_stream, meta)
                await self.dest.metadata.put(bucket, key, meta)

            elif event["type"] == "DELETE":
                await self.dest.metadata.delete(bucket, key)
                # Data GC handles actual chunk cleanup

            self._update_replication_lag_metric(event["timestamp"])

    def _should_replicate(self, event):
        bucket_config = self.replication_configs.get(event["bucket"])
        if not bucket_config:
            return False
        if bucket_config.get("prefix_filter"):
            return event["key"].startswith(bucket_config["prefix_filter"])
        return True
```

### 11.7 Object Lock (WORM Compliance)

```
Object Lock modes (for regulatory compliance):

  GOVERNANCE mode:
    - Protected from deletion by normal users
    - Admin users with special permission CAN override
    - Use case: internal data retention policies

  COMPLIANCE mode:
    - NOBODY can delete, not even root/admin
    - Object is immutable until retention period expires
    - Use case: SEC Rule 17a-4, HIPAA, financial records

  Legal Hold:
    - Prevents deletion regardless of retention period
    - Can be applied/removed independently of retention
    - Use case: litigation hold on specific documents

Implementation:
  - retention_mode: GOVERNANCE | COMPLIANCE
  - retain_until_date: timestamp
  - legal_hold: boolean
  - Enforced at metadata layer: DELETE/overwrite blocked if locked
```

---

## 12. Interview Q&A

### Q1: How do you achieve 11 nines of durability?

**Candidate:** "Eleven nines of durability means we can lose at most 1 object per 10 billion objects per year. Here's how we achieve it:

1. **Erasure coding (Reed-Solomon 8+4):** Each chunk is encoded into 12 shards. We can lose any 4 shards and still reconstruct the data. The mathematical probability of losing 5+ shards simultaneously is approximately 10^-17.

2. **Zone-aware placement:** Shards are distributed across at least 2-3 availability zones (different power, network, cooling). This means a full rack failure or even a zone outage doesn't reach the failure threshold.

3. **Proactive repair:** Background processes continuously monitor shard health. When a disk or node fails, the repair manager immediately re-creates missing shards on healthy nodes. The repair window is hours, not days, which keeps the effective failure probability ultra-low.

4. **Integrity scrubbing:** Every shard is checksummed and verified every 30 days. Silent data corruption (bit rot) is detected and repaired before it compounds.

5. **Multiple layers of checksums:** End-to-end checksums from client through to disk, verified on every read and write."

---

### Q2: Erasure coding vs. replication — trade-offs?

**Candidate:** "This is one of the most important architectural decisions.

**Replication (3× copies):**
- Pros: Simple implementation, fast reads (any replica), fast repair (just copy)
- Cons: 3× storage cost ($6M/month at 100 PB), only ~5 nines durability
- Best for: Hot data, small metadata, caches

**Erasure coding (8+4 RS):**
- Pros: 1.5× storage cost ($3M/month savings), 16+ nines durability, bandwidth-efficient repair
- Cons: CPU cost for encoding/decoding, higher read latency (need 8 shards), more complex implementation
- Best for: Bulk data storage, cold/warm tiers

**Our approach:** Use erasure coding as the default for STANDARD class. For the hot access tier, we can use replication for metadata and frequently accessed small objects where single-digit millisecond latency matters. The storage engine is pluggable—each storage class can use a different encoding scheme."

---

### Q3: How do you handle a 5TB file upload?

**Candidate:** "A 5 TB file absolutely requires multipart upload. Here's the flow:

1. **Client initiates** multipart upload → receives `upload_id`
2. **Client splits** the file into parts (e.g., 100 MB each = 50,000 parts, but max 10,000 parts so we'd use 500 MB parts)
3. **Parallel upload** of parts with configurable concurrency (e.g., 8 parallel streams)
4. Each part is independently chunked (64 MB) and erasure coded
5. Server tracks part ETags and checksums
6. **Client sends complete** request with ordered list of part ETags
7. Server stitches the chunk manifests into a single object metadata entry
8. **Abort path:** If upload fails or is abandoned, lifecycle policy cleans up parts after 7 days

Key resilience features:
- Parts are independently retried on failure
- Upload state is server-side, so client can resume after disconnect
- Checksums per-part prevent corrupt data from being accepted
- Total upload time for 5 TB at 1 Gbps: ~11 hours (parallelism helps saturate bandwidth)"

---

### Q4: How do you ensure consistency for read-after-write?

**Candidate:** "Strong read-after-write consistency is critical for user experience—you upload a file and immediately need to read it back.

**Implementation:**

1. **Write path:** Object data is written to data nodes first. Only after all required shards are confirmed does the metadata service commit the metadata record.

2. **Metadata commit is atomic:** We use a distributed KV store with linearizable writes (FoundationDB provides this natively). The PUT returns 200 only after the metadata is committed.

3. **Read path:** All reads go through the metadata service, which always reads from the leader (not stale replicas). Since the metadata was committed before the PUT response, any subsequent GET will find the metadata.

4. **For listings:** We provide eventual consistency because listing requires range scans across multiple metadata shards. Updating listing indexes across all shards synchronously would be too expensive. Most users accept this trade-off (AWS S3 offered eventual consistency for years before switching to strong consistency in 2020).

The key insight is: strong consistency for point reads (GET by exact key) is achievable at low cost by routing through the Raft leader. Strong consistency for range scans (LIST) is much more expensive and usually not required."

---

### Q5: How would you implement object versioning?

**Candidate:** "Versioning adds a dimension to the key space.

**Data model change:**
- Without versioning: key = `{bucket}/{object_key}`
- With versioning: key = `{bucket}/{object_key}#{version_id}`
- `version_id` is a timestamp-based UUID (sortable, unique)

**On PUT (versioning enabled):**
1. Generate new `version_id`
2. Write data chunks as normal
3. Write metadata with composite key `bucket/key#version_id`
4. Update a 'latest pointer' for `bucket/key` → `version_id`

**On GET (no version specified):**
1. Read the 'latest pointer' for `bucket/key`
2. If it points to a delete marker → return 404
3. Otherwise fetch that version's data

**On DELETE (versioning enabled):**
1. Don't remove any data
2. Insert a new version with `delete_marker: true`
3. Update latest pointer to the delete marker

**Space management:**
- Lifecycle policies can auto-expire non-current versions after N days
- GC reclaims chunks only when all referencing versions are deleted

This is exactly how S3 versioning works in production."

---

### Q6: How do you handle hot objects (millions of reads/sec on one object)?

**Candidate:** "A hot object—like a viral video or popular config file—can receive millions of reads per second. The storage cluster can't handle this from data nodes alone.

**Multi-layer strategy:**

1. **CDN Layer:** For publicly accessible or cacheable objects, the CDN handles most reads. Cache TTL can be set per-object. This handles 90%+ of hot-object traffic.

2. **Read replicas / caching tier:** For authenticated or non-CDN-eligible hot objects, we add an in-memory cache layer (Redis/Memcached) in front of the data nodes. The cache stores the fully assembled object.

3. **Replicated reads from multiple shards:** Erasure-coded data can be read from different subsets of shards. If 8 of 12 shards can serve a read, we have 12-choose-8 = 495 different read sets, distributing load across nodes.

4. **Auto-detection and promotion:** Monitor per-key request rates. When an object exceeds a threshold (e.g., 10K reads/sec), automatically:
   - Push to CDN edge caches
   - Create additional read-optimized replicas (full copies on SSD)
   - Update DNS/routing to distribute reads

5. **Range-read optimization:** For large hot objects (like video), encourage byte-range requests so different parts are served by different nodes simultaneously."

---

### Q7: How would you design a lifecycle policy engine?

**Candidate:** "The lifecycle engine runs as a background service that evaluates rules against all objects, typically once per day.

**Architecture:**

1. **Policy storage:** Lifecycle rules stored in bucket metadata (JSON).

2. **Evaluation engine:** A scheduled job that:
   - Iterates through all buckets with lifecycle rules
   - For each rule, scans matching objects (by prefix)
   - Evaluates age-based conditions against current timestamp
   - Enqueues transition or deletion actions

3. **Action execution:**
   - **Transition (e.g., STANDARD → IA):** Re-encode data with the target storage class's erasure coding parameters. Update metadata. This is a data migration that can be batched and throttled.
   - **Expiration:** Mark object for deletion. Follows same GC pipeline as manual deletes.
   - **Noncurrent version expiration:** Delete old versions past retention.

4. **Scale considerations:**
   - 1 billion objects can't be scanned sequentially in a day
   - Partition the scan across workers by metadata shard
   - Use last-evaluated markers for resumability
   - Throttle transitions to avoid I/O storms

5. **Idempotency:** Each evaluation is idempotent. If the engine crashes mid-run, it can resume safely."

---

### Q8: How do you handle cross-region replication with conflict resolution?

**Candidate:** "Cross-region replication (CRR) is fundamentally an asynchronous process with inherent conflict potential.

**Normal flow:**
1. Every metadata change in the source region produces a changelog event (via Kafka or a change data capture stream).
2. The CRR service consumes events and replays them in the destination region.
3. Replication lag is typically seconds, monitored via metrics.

**Conflict scenarios and resolution:**

1. **Same key written in both regions (active-active):**
   - Use last-writer-wins with wall-clock timestamps
   - Each write carries a hybrid logical clock (HLC) timestamp
   - On conflict: the write with the higher HLC wins
   - Losing write is preserved as a non-current version (if versioning enabled)

2. **Delete in source, write in destination:**
   - If timestamps are close, the later operation wins
   - If delete is newer: object is deleted in both regions
   - If write is newer: object exists in both regions

3. **Network partition between regions:**
   - Replication pauses; events queue in Kafka
   - On recovery: events are replayed in order
   - For active-active: vector clocks detect concurrent writes

**Best practice:** Most production systems use active-passive CRR (one primary region for writes) to avoid conflicts entirely. Active-active is only used when write latency requirements demand it."

---

## 13. Production Checklist

### Pre-Launch (Week -2 to -1)

- [ ] Load test: sustain 100K req/sec for 24 hours with mixed workloads
- [ ] Failure injection: kill 10% of data nodes, verify repair completes within SLA
- [ ] Durability validation: upload 1M objects, verify all retrievable with correct checksums
- [ ] Security audit: penetration testing, ACL bypass attempts, encryption verification
- [ ] Capacity planning: confirm 20% headroom on storage and compute
- [ ] Runbook: documented procedures for every alert
- [ ] Backup: metadata store snapshots to separate storage system

### Day 1

- [ ] Gradual traffic ramp: 1% → 10% → 50% → 100% over 4 hours
- [ ] Monitor error rates, latency percentiles, replication lag
- [ ] Verify all storage classes functioning (STANDARD, IA, ARCHIVE)
- [ ] Confirm audit logging is capturing all operations
- [ ] On-call engineers briefed and available

### Week 1

- [ ] Review P99 latency trends — identify slow paths
- [ ] Analyze storage distribution — check for hot nodes
- [ ] Verify garbage collection is running correctly
- [ ] Test multipart upload with 1 TB file end-to-end
- [ ] Validate pre-signed URLs work correctly with SDK
- [ ] Confirm lifecycle policy engine is evaluating rules

### Month 1

- [ ] Full durability audit: verify shard counts match expectations for all objects
- [ ] Cost analysis: actual storage cost vs. projections
- [ ] Performance tuning: optimize erasure coding parameters based on real workload
- [ ] Capacity planning update: project growth for next 6 months
- [ ] Enable cross-region replication for critical buckets
- [ ] Document lessons learned and update architecture docs
- [ ] Plan for next features: event notifications, object lock

---

## Summary: Technical Decisions

| Decision | Choice | Rationale |
|---|---|---|
| **Metadata Store** | Distributed KV (FoundationDB) | Horizontal scale, linearizable reads, range scans |
| **Data Encoding** | Erasure coding (Reed-Solomon 8+4) | 1.5× overhead vs 3× replication; 16+ nines durability |
| **Chunk Size** | 64 MB | Balance between parallelism and metadata overhead |
| **Data Placement** | Consistent hashing + virtual nodes | Minimal data movement on cluster changes |
| **Consistency Model** | Strong read-after-write; eventual for listings | Practical trade-off: point reads via leader, listings fan out |
| **Large File Upload** | Multipart upload (up to 10,000 parts) | Resumable, parallelizable, failure-tolerant |
| **Durability Verification** | Multi-layer checksums + periodic scrubbing | Catches bit rot, network corruption, software bugs |
| **Access Control** | Bucket ACLs + per-object ACLs + pre-signed URLs | Flexible: public, private, and temporary access |
| **Encryption** | SSE-S3 / SSE-KMS / SSE-C | Defense in depth; customer key management options |
| **Garbage Collection** | Mark-and-sweep with 24-hour safety window | Prevents accidental data loss from race conditions |

## Scalability Path

```
Phase 1 (Launch):          Phase 2 (Growth):          Phase 3 (Planet-Scale):
  1 region                   3 regions                  10+ regions
  100 data nodes             500 data nodes             5,000+ data nodes
  10 PB storage              100 PB storage             1+ EB storage
  10K req/sec                100K req/sec               1M+ req/sec
  Basic storage classes      All storage classes        Custom tiers
  Single-region              Cross-region replication   Global namespace
  No CDN                     CDN for public objects     Full CDN integration
  Manual scaling             Auto-scaling               ML-driven placement
```

---

> **Interview Tip:** When discussing object storage in a FAANG interview, always lead with the durability math. Explaining why erasure coding achieves 11 nines while replication only gives 5 nines—and backing it with the binomial probability calculation—immediately signals deep understanding. Follow up with the write path (how data flows from client to disk) and you'll cover 80% of what interviewers want to hear.
