# Complete System Design: Distributed Cache System (Redis-like)

> **Complexity Level:** Advanced  
> **Estimated Time:** 60-90 minutes in interview  
> **Real-World Examples:** Redis, Memcached, Amazon ElastiCache, Hazelcast, Apache Ignite

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Distributed Cache Architecture](#7-deep-dive-distributed-cache-architecture)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a distributed in-memory cache system like Redis that provides low-latency data access for web applications."

**Interviewer's Perspective:**  
This is a foundational system design question for Staff+ level roles. It assesses:
- **Distributed systems fundamentals** — consistent hashing, data partitioning, replication
- **In-memory data structure design** — hash tables, skip lists, linked lists
- **Eviction policies** — LRU, LFU, TTL-based, and their trade-offs
- **Consistency vs availability** — CAP theorem applied to caching
- **Failure handling** — node failures, split-brain, cache stampede
- **Performance engineering** — achieving sub-millisecond latency at scale

**Why This Problem Is Hard:**  
A distributed cache must serve millions of operations per second with sub-millisecond latency, handle node failures without data loss, scale horizontally across hundreds of nodes, and support rich data structures — all while keeping memory usage efficient. The real challenge is balancing consistency, availability, and partition tolerance in an in-memory distributed system where every microsecond matters.

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I design the system, I'd like to clarify the scope. Caching systems range from simple key-value stores to full-featured data structure servers. Can I ask a few questions?"

**Interviewer:** "Go ahead."

**Candidate:** "What data operations do we need to support? Just basic GET/SET, or richer data types like lists and sorted sets?"

**Interviewer:** "We need a rich data model. Think Redis — strings, hashes, lists, sets, and sorted sets. Support GET, SET, DELETE as the basics, and type-specific commands like HSET/HGET, LPUSH/LPOP, SADD/SMEMBERS."

**Candidate:** "Should keys support time-to-live (TTL)? And do we need atomic operations like INCR/DECR?"

**Interviewer:** "Yes to both. TTL per key is essential, and atomic increment/decrement is critical for use cases like rate limiting and counters."

**Candidate:** "What about pub/sub messaging and Lua scripting?"

**Interviewer:** "Pub/sub is required. Lua scripting is nice-to-have — include it in your advanced features section."

**Candidate:** "For non-functional requirements, what latency targets should I design for?"

**Interviewer:** "Sub-millisecond for GET and SET operations — under 1ms at p99. Each node should handle at least 1 million operations per second."

**Candidate:** "What's the availability target, and should I prioritize availability over consistency?"

**Interviewer:** "99.99% availability. For a cache, we lean toward AP in CAP theorem — availability over strict consistency. Eventual consistency for replication is fine."

**Candidate:** "Should the system support persistence, or is it purely in-memory?"

**Interviewer:** "Both. Users should be able to choose between pure in-memory mode and persistence options — RDB snapshots and append-only file (AOF) logs."

**Candidate:** "Let me confirm the scale: how many nodes, total keys, and aggregate throughput?"

**Interviewer:** "Design for a 100-node cluster storing 100 million keys across 500GB of total memory, handling 10 million operations per second in aggregate."

**Candidate:** "Perfect. Let me summarize what I've gathered."

### 2.1 Functional Requirements

| # | Requirement | Description |
|---|------------|-------------|
| FR-1 | GET/SET/DELETE | Basic key-value CRUD operations |
| FR-2 | Multiple Data Types | String, hash, list, set, sorted set |
| FR-3 | TTL per Key | Time-to-live expiration on any key |
| FR-4 | Atomic Operations | INCR, DECR, SETNX (set-if-not-exists) |
| FR-5 | Pub/Sub Messaging | Publish and subscribe to channels |
| FR-6 | Pipelining | Batch multiple commands in a single round-trip |
| FR-7 | Lua Scripting | Server-side scripts for complex atomic operations |

### 2.2 Non-Functional Requirements

| # | Requirement | Target |
|---|------------|--------|
| NFR-1 | Latency | < 1ms p99 for GET/SET operations |
| NFR-2 | Throughput | 1M+ ops/sec per node |
| NFR-3 | Availability | 99.99% uptime |
| NFR-4 | Horizontal Scaling | Add/remove nodes without downtime |
| NFR-5 | Persistence | Optional RDB snapshots and AOF logs |
| NFR-6 | Replication | Leader-follower with automatic failover |
| NFR-7 | Memory Efficiency | < 10% overhead beyond raw data size |

### 2.3 Scale Parameters

| Parameter | Value |
|-----------|-------|
| Cluster Nodes | 100 |
| Total Keys | 100 million |
| Total Memory | 500 GB |
| Aggregate Ops/sec | 10 million |
| Ops/sec per Node | ~100K |
| Avg Key Size | 50 bytes |
| Avg Value Size | 500 bytes |
| Replication Factor | 3 (1 leader + 2 replicas) |

---

## 3. Scale Estimation

### 3.1 Traffic Estimates

**Candidate:** "Let me work through the numbers."

```
Aggregate throughput:
  10M ops/sec across 100 nodes
  = 100K ops/sec per node

Read:Write ratio (typical cache):
  80:20 read-heavy
  Reads:  8M ops/sec aggregate  → 80K reads/sec per node
  Writes: 2M ops/sec aggregate  → 20K writes/sec per node
```

### 3.2 Memory Estimates

```
Per key-value pair:
  Key:   50 bytes (average)
  Value: 500 bytes (average)
  Metadata (TTL, type, pointers): ~80 bytes
  Total per entry: ~630 bytes

100M keys × 630 bytes = 63 GB raw data
With hash table overhead (~1.5x): ~95 GB
With replication (3x): ~285 GB

Per node (100 nodes):
  Raw data: ~950 MB
  With overhead: ~1.5 GB
  Recommended RAM per node: 8 GB (headroom for spikes, fragmentation)
```

### 3.3 Network Bandwidth

```
Client traffic:
  10M ops/sec × 550 bytes avg (key + value) = 5.5 GB/sec cluster-wide
  Per node: 55 MB/sec client traffic

Replication traffic:
  2M writes/sec × 550 bytes × 2 replicas = 2.2 GB/sec cluster-wide
  Per node: 22 MB/sec replication traffic

Total per node: ~77 MB/sec → 10 Gbps NIC sufficient
```

### 3.4 Capacity Summary

| Metric | Per Node | Cluster-Wide |
|--------|----------|--------------|
| Operations/sec | 100K | 10M |
| Memory (data) | ~1.5 GB | ~150 GB |
| Memory (with overhead) | ~5 GB | ~500 GB |
| Network (client) | 55 MB/s | 5.5 GB/s |
| Network (replication) | 22 MB/s | 2.2 GB/s |
| Keys | 1M | 100M |

---

## 4. High-Level Design

### 4.1 Architecture Overview

**Candidate:** "Let me sketch the high-level architecture."

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CLIENT APPLICATIONS                         │
│  (Web Servers, API Servers, Microservices)                         │
└────────────────────────────┬────────────────────────────────────────┘
                             │
                    ┌────────▼────────┐
                    │  Client Library  │
                    │ ┌──────────────┐ │
                    │ │  Consistent  │ │
                    │ │   Hashing /  │ │
                    │ │ Slot Mapping │ │
                    │ ├──────────────┤ │
                    │ │  Connection  │ │
                    │ │    Pooling   │ │
                    │ ├──────────────┤ │
                    │ │  Pipelining  │ │
                    │ └──────────────┘ │
                    └────────┬────────┘
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
  ┌───────▼───────┐  ┌──────▼──────┐  ┌───────▼───────┐
  │  Cache Node 1 │  │ Cache Node 2│  │ Cache Node N  │
  │ ┌───────────┐ │  │ ┌─────────┐ │  │ ┌───────────┐ │
  │ │ Event Loop│ │  │ │Event    │ │  │ │ Event Loop│ │
  │ │ (single   │ │  │ │Loop     │ │  │ │ (single   │ │
  │ │  thread)  │ │  │ │         │ │  │ │  thread)  │ │
  │ ├───────────┤ │  │ ├─────────┤ │  │ ├───────────┤ │
  │ │ In-Memory │ │  │ │In-Memory│ │  │ │ In-Memory │ │
  │ │ Data Store│ │  │ │Data     │ │  │ │ Data Store│ │
  │ │ (Hash     │ │  │ │Store    │ │  │ │ (Hash     │ │
  │ │  Table)   │ │  │ │         │ │  │ │  Table)   │ │
  │ ├───────────┤ │  │ ├─────────┤ │  │ ├───────────┤ │
  │ │Persistence│ │  │ │Persist. │ │  │ │Persistence│ │
  │ │ RDB / AOF │ │  │ │RDB/AOF  │ │  │ │ RDB / AOF │ │
  │ └───────────┘ │  │ └─────────┘ │  │ └───────────┘ │
  │       │       │  │      │      │  │       │       │
  │  ┌────▼────┐  │  │ ┌────▼───┐  │  │  ┌────▼────┐  │
  │  │Replicas │  │  │ │Replicas│  │  │  │Replicas │  │
  │  │ R1, R2  │  │  │ │ R1, R2 │  │  │  │ R1, R2  │  │
  │  └─────────┘  │  │ └────────┘  │  │  └─────────┘  │
  └───────────────┘  └─────────────┘  └───────────────┘
          │                  │                  │
          └──────────────────┼──────────────────┘
                             │
              ┌──────────────▼──────────────┐
              │      Cluster Manager        │
              │  ┌────────────────────────┐  │
              │  │   Gossip Protocol      │  │
              │  │   (node discovery,     │  │
              │  │    health checks)      │  │
              │  ├────────────────────────┤  │
              │  │   Configuration        │  │
              │  │   Service (ZooKeeper)  │  │
              │  ├────────────────────────┤  │
              │  │   Slot Assignment      │  │
              │  │   & Rebalancing        │  │
              │  └────────────────────────┘  │
              └─────────────────────────────┘
```

### 4.2 API Design

**Candidate:** "Here's the command interface, modeled after Redis protocol."

```
# String Operations
SET key value [EX seconds] [PX milliseconds] [NX|XX]
GET key
DEL key [key ...]
INCR key
DECR key
MGET key [key ...]
MSET key value [key value ...]

# Hash Operations
HSET key field value [field value ...]
HGET key field
HGETALL key
HDEL key field [field ...]

# List Operations
LPUSH key value [value ...]
RPUSH key value [value ...]
LPOP key
RPOP key
LRANGE key start stop

# Set Operations
SADD key member [member ...]
SMEMBERS key
SISMEMBER key member
SINTER key [key ...]

# Sorted Set Operations
ZADD key score member [score member ...]
ZRANGE key start stop [WITHSCORES]
ZRANK key member
ZRANGEBYSCORE key min max

# Key Operations
EXPIRE key seconds
TTL key
EXISTS key
KEYS pattern

# Pub/Sub
PUBLISH channel message
SUBSCRIBE channel [channel ...]
UNSUBSCRIBE channel [channel ...]
```

### 4.3 Data Flow

**Write Path (SET):**

```
Client ──SET foo bar──▶ Client Library
                            │
                   hash("foo") mod 16384
                        = slot 4521
                            │
                     slot 4521 → Node 3
                            │
                        ┌───▼───┐
                        │Node 3 │
                        │(Leader)│
                        └───┬───┘
                            │
              ┌─────────────┼─────────────┐
              │             │             │
         ┌────▼────┐  ┌────▼────┐  ┌─────▼─────┐
         │ Write   │  │ Replicate│  │ Replicate │
         │ to mem  │  │ to R1   │  │ to R2     │
         └────┬────┘  └─────────┘  └───────────┘
              │
         ┌────▼────┐
         │ AOF     │ (if persistence enabled)
         │ append  │
         └────┬────┘
              │
         ◄──OK──
```

**Read Path (GET):**

```
Client ──GET foo──▶ Client Library
                        │
               hash("foo") mod 16384
                    = slot 4521
                        │
                 slot 4521 → Node 3
                        │
                    ┌───▼───┐
                    │Node 3 │──▶ Hash table lookup
                    │(Leader)│      O(1) average
                    └───┬───┘
                        │
                   ◄──"bar"──
```

---

## 5. Deep Dive: Core Components

### 5.1 Cache Node Architecture

**Candidate:** "Each cache node uses a single-threaded event loop, just like Redis. This eliminates lock contention and context-switching overhead."

```
┌──────────────────────────────────────────────────────┐
│                    CACHE NODE                        │
│                                                      │
│  ┌────────────────────────────────────────────────┐  │
│  │           Event Loop (Single Thread)           │  │
│  │                                                │  │
│  │   ┌─────────┐    ┌──────────┐    ┌─────────┐  │  │
│  │   │  epoll  │───▶│ Command  │───▶│ Execute │  │  │
│  │   │  wait   │    │  Parser  │    │ Command │  │  │
│  │   └─────────┘    └──────────┘    └────┬────┘  │  │
│  │        ▲                              │       │  │
│  │        │         ┌──────────┐         │       │  │
│  │        └─────────│  Send    │◄────────┘       │  │
│  │                  │ Response │                  │  │
│  │                  └──────────┘                  │  │
│  └────────────────────────────────────────────────┘  │
│                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐  │
│  │ Hash Table   │  │ Skip List    │  │ Linked    │  │
│  │ (key→value)  │  │ (sorted sets)│  │ List      │  │
│  └──────────────┘  └──────────────┘  └───────────┘  │
│                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐  │
│  │ Expiry Table │  │ AOF Writer   │  │ RDB       │  │
│  │ (key→TTL)    │  │ (async)      │  │ Snapshot  │  │
│  └──────────────┘  └──────────────┘  └───────────┘  │
└──────────────────────────────────────────────────────┘
```

**Why single-threaded?**

```python
# The bottleneck is NOT CPU — it's network I/O and memory access.
# A single thread avoids:
#   1. Lock contention on shared data structures
#   2. Context switching overhead (thousands of μs per switch)
#   3. Cache line bouncing across CPU cores
#
# With I/O multiplexing (epoll), one thread handles 100K+ connections:

import selectors

class CacheEventLoop:
    def __init__(self):
        self.selector = selectors.DefaultSelector()
        self.data_store = {}

    def run(self):
        while True:
            events = self.selector.select(timeout=0.001)
            for key, mask in events:
                callback = key.data
                callback(key.fileobj, mask)

            self._process_time_events()  # TTL expiry, background tasks

    def _process_time_events(self):
        self._expire_keys_lazily()
        self._run_background_persistence()
```

### 5.2 Client Library

**Candidate:** "The client library is cluster-aware and handles routing, connection pooling, and pipelining."

```javascript
class CacheClient {
  constructor(clusterNodes) {
    this.slotMap = new Map();       // slot → node address
    this.connectionPool = new Map(); // node → [connections]
    this.initClusterTopology(clusterNodes);
  }

  async set(key, value, options = {}) {
    const slot = this.hashSlot(key);
    const node = this.slotMap.get(slot);
    const conn = await this.getConnection(node);

    let cmd = `SET ${key} ${value}`;
    if (options.ex) cmd += ` EX ${options.ex}`;
    if (options.nx) cmd += ` NX`;

    const result = await conn.execute(cmd);

    if (result === 'MOVED') {
      this.refreshTopology();
      return this.set(key, value, options); // retry with updated topology
    }
    return result;
  }

  hashSlot(key) {
    // Support hash tags: {user:1000}.name → hash on "user:1000"
    const tagMatch = key.match(/\{(.+?)\}/);
    const hashKey = tagMatch ? tagMatch[1] : key;
    return crc16(hashKey) % 16384;
  }

  async pipeline(commands) {
    // Group commands by target node
    const grouped = new Map();
    for (const cmd of commands) {
      const slot = this.hashSlot(cmd.key);
      const node = this.slotMap.get(slot);
      if (!grouped.has(node)) grouped.set(node, []);
      grouped.get(node).push(cmd);
    }

    // Send all commands to each node in a single round-trip
    const promises = [...grouped.entries()].map(([node, cmds]) =>
      this.sendPipeline(node, cmds)
    );
    return Promise.all(promises);
  }
}
```

### 5.3 Cluster Manager

**Candidate:** "The cluster manager uses a gossip protocol for node discovery and health checking."

```
Gossip Protocol — Node Communication:

  Node A ──PING──▶ Node B
  Node B ──PONG──▶ Node A
         (includes B's view of cluster state)

  Every 100ms, each node:
    1. Picks a random node
    2. Sends PING with its cluster state
    3. Receives PONG with the other node's state
    4. Merges states (crdt-style conflict resolution)

  Failure detection:
    - If Node B doesn't PONG within 300ms → mark as PFAIL (possible failure)
    - If majority of nodes mark B as PFAIL → promote to FAIL
    - Trigger failover: promote B's replica to leader
```

```python
class ClusterManager:
    def __init__(self, node_id, nodes):
        self.node_id = node_id
        self.nodes = {n.id: n for n in nodes}
        self.slot_assignment = {}  # slot → node_id
        self.epoch = 0

    def gossip_tick(self):
        """Called every 100ms."""
        target = self._random_node()
        my_state = self._build_cluster_state()
        response = target.send_ping(my_state)
        self._merge_state(response)

    def detect_failure(self, node_id):
        node = self.nodes[node_id]
        if time.time() - node.last_pong > 0.3:
            node.status = 'PFAIL'
            pfail_count = sum(
                1 for n in self.nodes.values()
                if n.pfail_reports.get(node_id)
            )
            if pfail_count > len(self.nodes) // 2:
                node.status = 'FAIL'
                self._trigger_failover(node_id)

    def _trigger_failover(self, failed_node_id):
        replicas = self._get_replicas(failed_node_id)
        best_replica = max(replicas, key=lambda r: r.replication_offset)
        best_replica.promote_to_leader()
        self._reassign_slots(failed_node_id, best_replica.id)
        self.epoch += 1
```

### 5.4 Persistence Manager

```
Two persistence strategies:

┌─────────────────────────────────────────────────────┐
│                  RDB Snapshots                       │
│                                                      │
│  - Point-in-time binary snapshot of entire dataset   │
│  - Created via fork() — child writes while parent    │
│    continues serving (copy-on-write)                 │
│  - Compact binary format, fast to load               │
│  - Data loss window: since last snapshot              │
│  - Typical interval: every 5-15 minutes              │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│             Append-Only File (AOF)                   │
│                                                      │
│  - Logs every write command to disk                  │
│  - Three fsync policies:                             │
│      always:    fsync after every write (safest)     │
│      everysec:  fsync once per second (recommended)  │
│      no:        let OS decide (fastest)              │
│  - AOF rewrite: compact log by reading current state │
│  - Data loss window: up to 1 second (everysec)       │
└─────────────────────────────────────────────────────┘
```

```python
class PersistenceManager:
    def __init__(self, data_store, config):
        self.data_store = data_store
        self.aof_file = open(config.aof_path, 'a')
        self.aof_policy = config.aof_fsync  # 'always', 'everysec', 'no'
        self.rdb_interval = config.rdb_save_interval

    def append_aof(self, command):
        """Append command to AOF log."""
        self.aof_file.write(f"*{len(command)}\r\n")
        for arg in command:
            self.aof_file.write(f"${len(arg)}\r\n{arg}\r\n")

        if self.aof_policy == 'always':
            os.fsync(self.aof_file.fileno())

    def create_rdb_snapshot(self):
        """Fork and write snapshot in child process."""
        pid = os.fork()
        if pid == 0:
            # Child: write binary snapshot
            with open('dump.rdb', 'wb') as f:
                f.write(b'REDIS0011')  # magic + version
                for key, value in self.data_store.items():
                    self._write_kv_pair(f, key, value)
                f.write(b'\xFF')  # EOF marker
                f.write(self._checksum())
            os._exit(0)
        # Parent continues serving requests (copy-on-write)
```

---

## 6. Deep Dive: Database Design

### 6.1 In-Memory Data Structures

**Candidate:** "The core of the cache is a hash table mapping keys to typed values. Each data type uses its own optimized structure."

```
Main Key-Value Hash Table:
┌───────────────────────────────────────────────────────┐
│ dict (hash table with incremental rehashing)          │
│                                                       │
│  Bucket 0: key_a → {type: STRING, value: "hello"}    │
│  Bucket 1: key_b → {type: HASH, value: {f1:v1,...}}  │
│  Bucket 2: NULL                                       │
│  Bucket 3: key_c → {type: ZSET, value: skiplist}     │
│  ...                                                  │
│  Bucket N: key_z → {type: LIST, value: quicklist}     │
│                                                       │
│  Load factor threshold: 1.0 (triggers rehash)         │
│  Rehash: incremental — 1 bucket per operation         │
└───────────────────────────────────────────────────────┘
```

**Hash Table with Incremental Rehashing:**

```python
class IncrementalHashTable:
    """
    Two tables: ht[0] (active) and ht[1] (rehash target).
    During rehash, every operation migrates one bucket from ht[0] to ht[1].
    """
    def __init__(self, initial_size=4):
        self.ht = [
            [None] * initial_size,  # ht[0]: active table
            None                     # ht[1]: rehash target (None if not rehashing)
        ]
        self.sizes = [initial_size, 0]
        self.used = [0, 0]
        self.rehash_idx = -1  # -1 means not rehashing

    def get(self, key):
        h = hash(key)
        for table_idx in (0, 1):
            if self.ht[table_idx] is None:
                continue
            idx = h % self.sizes[table_idx]
            entry = self.ht[table_idx][idx]
            while entry:
                if entry.key == key:
                    return entry.value
                entry = entry.next
        return None

    def set(self, key, value):
        if self.rehash_idx >= 0:
            self._rehash_step()  # migrate one bucket per operation

        load_factor = self.used[0] / self.sizes[0]
        if load_factor >= 1.0 and self.rehash_idx < 0:
            self._start_rehash()

        table_idx = 1 if self.rehash_idx >= 0 else 0
        idx = hash(key) % self.sizes[table_idx]
        # Insert into appropriate table
        entry = Entry(key, value)
        entry.next = self.ht[table_idx][idx]
        self.ht[table_idx][idx] = entry
        self.used[table_idx] += 1

    def _rehash_step(self):
        """Migrate one bucket from ht[0] to ht[1]."""
        while self.ht[0][self.rehash_idx] is None:
            self.rehash_idx += 1
            if self.rehash_idx >= self.sizes[0]:
                self._finish_rehash()
                return

        entry = self.ht[0][self.rehash_idx]
        while entry:
            next_entry = entry.next
            idx = hash(entry.key) % self.sizes[1]
            entry.next = self.ht[1][idx]
            self.ht[1][idx] = entry
            self.used[0] -= 1
            self.used[1] += 1
            entry = next_entry

        self.ht[0][self.rehash_idx] = None
        self.rehash_idx += 1
```

**Skip List for Sorted Sets:**

```
Skip List Structure (for ZSET operations):

Level 3:  HEAD ──────────────────────────────────▶ 67 ──────────────▶ NIL
Level 2:  HEAD ──────────────▶ 23 ──────────────▶ 67 ──────────────▶ NIL
Level 1:  HEAD ──────▶ 12 ──▶ 23 ──▶ 34 ──────▶ 67 ──▶ 89 ──────▶ NIL
Level 0:  HEAD ──▶ 7 ▶ 12 ──▶ 23 ──▶ 34 ──▶ 45 ▶ 67 ──▶ 89 ──▶ 92 ▶ NIL

- Average O(log N) for insert, delete, search
- Simpler than balanced trees (no rotations)
- Cache-friendly sequential access for ZRANGE
```

### 6.2 Key Metadata

```sql
-- Logical representation of key metadata (stored in-memory, not SQL)
-- Shown as SQL for clarity

CREATE TABLE key_metadata (
    key           VARCHAR(512) PRIMARY KEY,
    type          ENUM('string', 'hash', 'list', 'set', 'zset'),
    encoding      ENUM('raw', 'int', 'ziplist', 'hashtable', 'skiplist', 'quicklist'),
    ttl_ms        BIGINT DEFAULT -1,        -- -1 means no expiry
    created_at    BIGINT,                    -- unix timestamp ms
    last_accessed BIGINT,                    -- for LRU tracking
    access_count  INT DEFAULT 0,             -- for LFU tracking
    memory_bytes  INT,                       -- approximate memory usage
    slot          INT                        -- hash slot (0-16383)
);

-- Expiry index (for active expiry scanning)
CREATE INDEX idx_expiry ON key_metadata(ttl_ms) WHERE ttl_ms > 0;
```

### 6.3 Memory Encoding Optimizations

**Candidate:** "Redis uses compact encodings for small data to save memory."

| Data Type | Small Encoding | Large Encoding | Threshold |
|-----------|---------------|----------------|-----------|
| String | int (if numeric) | raw SDS | N/A |
| Hash | ziplist | hashtable | 64 fields or 128 bytes per value |
| List | ziplist | quicklist | 128 entries or 64 bytes per entry |
| Set | intset (if all ints) | hashtable | 512 members |
| Sorted Set | ziplist | skiplist + hashtable | 128 members or 64 bytes per member |

```
Ziplist (compact sequential memory layout):
┌──────┬─────┬───────┬───────┬───────┬─────┐
│zlbytes│zltail│zllen │entry1 │entry2 │zlend│
│4 bytes│4 bytes│2 byte│ ...   │ ...   │1 byte│
└──────┴─────┴───────┴───────┴───────┴─────┘

Each entry:
┌──────────────┬──────────┬───────┐
│prev_entry_len│ encoding │ data  │
└──────────────┴──────────┴───────┘
```

---

## 7. Deep Dive: Distributed Cache Architecture

**Interviewer:** "This is the most critical section. Walk me through how you'd distribute data across nodes."

### 7.1 Consistent Hashing

**Candidate:** "Consistent hashing minimizes key redistribution when nodes are added or removed."

```
Hash Ring with Virtual Nodes:

                        0 (top)
                    ╱           ╲
                  N1v1          N2v1
                ╱                  ╲
             N3v2                   N1v2
            │                         │
         N2v3                       N3v1
            │                         │
             N1v3                   N2v2
                ╲                  ╱
                  N3v3          N1v4
                    ╲           ╱
                     2^32 - 1

  Key "user:1001" → hash = 0x3A2B...
  Walk clockwise → lands on N2v1 → routed to Node 2

  Virtual nodes (vnodes): each physical node owns multiple
  points on the ring. With 150 vnodes per node, key
  distribution variance drops below 5%.
```

**Consistent Hashing Implementation:**

```javascript
const crypto = require('crypto');

class ConsistentHash {
  constructor(virtualNodesPerNode = 150) {
    this.ring = new Map();        // hash position → node
    this.sortedKeys = [];         // sorted hash positions
    this.vnodeCount = virtualNodesPerNode;
    this.nodes = new Set();
  }

  addNode(node) {
    this.nodes.add(node);
    for (let i = 0; i < this.vnodeCount; i++) {
      const vkey = `${node}:vnode:${i}`;
      const hash = this._hash(vkey);
      this.ring.set(hash, node);
      this.sortedKeys.push(hash);
    }
    this.sortedKeys.sort((a, b) => a - b);
  }

  removeNode(node) {
    this.nodes.delete(node);
    for (let i = 0; i < this.vnodeCount; i++) {
      const vkey = `${node}:vnode:${i}`;
      const hash = this._hash(vkey);
      this.ring.delete(hash);
    }
    this.sortedKeys = this.sortedKeys.filter(k => this.ring.has(k));
  }

  getNode(key) {
    if (this.sortedKeys.length === 0) return null;
    const hash = this._hash(key);
    // Binary search for the first position >= hash
    let lo = 0, hi = this.sortedKeys.length;
    while (lo < hi) {
      const mid = (lo + hi) >>> 1;
      if (this.sortedKeys[mid] < hash) lo = mid + 1;
      else hi = mid;
    }
    // Wrap around
    const idx = lo % this.sortedKeys.length;
    return this.ring.get(this.sortedKeys[idx]);
  }

  _hash(key) {
    return parseInt(
      crypto.createHash('md5').update(key).digest('hex').slice(0, 8),
      16
    );
  }

  // Impact analysis: how many keys move when adding a node
  getRedistributionRatio() {
    // Theoretical: only 1/N keys need to move (N = number of nodes)
    return 1 / this.nodes.size;
  }
}

// Usage
const ring = new ConsistentHash(150);
ring.addNode('cache-node-1');
ring.addNode('cache-node-2');
ring.addNode('cache-node-3');

console.log(ring.getNode('user:1001'));  // → 'cache-node-2'
console.log(ring.getNode('session:abc')); // → 'cache-node-1'

// Adding a node only moves ~1/4 of keys (not 3/4 like naive modulo)
ring.addNode('cache-node-4');
```

### 7.2 Data Partitioning (Hash Slots)

**Candidate:** "Redis Cluster uses a fixed number of hash slots (16384) instead of a continuous hash ring. This makes slot assignment explicit and manageable."

```
Hash Slot Assignment (16384 slots):

  Node 1: slots [0      - 5460]      → 5461 slots
  Node 2: slots [5461   - 10922]     → 5462 slots
  Node 3: slots [10923  - 16383]     → 5461 slots

  Key → Slot mapping:
    slot = CRC16(key) mod 16384

  Example:
    "user:1001" → CRC16("user:1001") = 7342 → slot 7342 → Node 2
    "order:555" → CRC16("order:555") = 3201 → slot 3201 → Node 1

  Hash Tags (force keys to same slot):
    "{user:1001}.profile" → CRC16("user:1001") → same slot
    "{user:1001}.orders"  → CRC16("user:1001") → same slot
    Enables multi-key operations on co-located keys
```

```python
def hash_slot(key: str) -> int:
    """Calculate the hash slot for a key, supporting hash tags."""
    # Extract hash tag if present: {tag}rest → hash on "tag"
    start = key.find('{')
    if start != -1:
        end = key.find('}', start + 1)
        if end != -1 and end != start + 1:
            key = key[start + 1:end]

    return crc16(key.encode()) % 16384

# Slot migration: moving slots between nodes during resharding
class SlotMigrator:
    def __init__(self, source_node, target_node, slot):
        self.source = source_node
        self.target = target_node
        self.slot = slot

    def migrate(self):
        """Live migration of a slot's keys from source to target."""
        # 1. Mark slot as MIGRATING on source, IMPORTING on target
        self.source.set_slot_state(self.slot, 'MIGRATING', self.target.id)
        self.target.set_slot_state(self.slot, 'IMPORTING', self.source.id)

        # 2. Move keys one by one (or in batches)
        keys = self.source.get_keys_in_slot(self.slot, batch_size=100)
        while keys:
            for key in keys:
                # Atomic MIGRATE: serializes key, sends to target, deletes from source
                self.source.migrate_key(key, self.target.address)
            keys = self.source.get_keys_in_slot(self.slot, batch_size=100)

        # 3. Update slot ownership in cluster state
        self.target.set_slot_owner(self.slot, self.target.id)
        self.source.set_slot_owner(self.slot, self.target.id)
```

### 7.3 Replication

**Candidate:** "Each leader node replicates to one or more followers for fault tolerance."

```
Leader-Follower Replication:

  ┌──────────┐     async replication     ┌──────────┐
  │  Leader   │──────────────────────────▶│ Replica 1│
  │  Node 3   │                           │          │
  │           │     async replication     ├──────────┤
  │  Writes   │──────────────────────────▶│ Replica 2│
  │  + Reads  │                           │          │
  └──────────┘                           └──────────┘
                                          (reads only)

  Replication Process:
  1. Replica sends PSYNC to leader
  2. Leader starts RDB snapshot (background)
  3. Leader streams RDB to replica
  4. Leader sends buffered writes that occurred during snapshot
  5. Ongoing: leader streams every write command to replicas
```

| Replication Mode | Latency Impact | Durability | Use Case |
|-----------------|---------------|------------|----------|
| Async (default) | None | May lose recent writes on failover | General caching |
| Semi-sync (WAIT) | +1-5ms | Guaranteed N replicas received | Important data |
| Sync | +5-50ms | Zero data loss | Rarely used in caches |

**Failover with Raft-style Leader Election:**

```python
class FailoverManager:
    def __init__(self, node, cluster_state):
        self.node = node
        self.cluster = cluster_state
        self.current_epoch = 0

    def initiate_failover(self, failed_leader_id):
        """Replica promotes itself when leader fails."""
        # 1. Increment epoch (like Raft term)
        self.current_epoch += 1

        # 2. Request votes from all leaders
        votes = 0
        for node in self.cluster.get_leader_nodes():
            vote = node.request_vote(
                epoch=self.current_epoch,
                candidate=self.node.id,
                replication_offset=self.node.replication_offset
            )
            if vote.granted:
                votes += 1

        # 3. Need majority of leaders to agree
        quorum = len(self.cluster.get_leader_nodes()) // 2 + 1
        if votes >= quorum:
            self._promote_to_leader(failed_leader_id)

    def _promote_to_leader(self, old_leader_id):
        # Take over the failed leader's slot assignments
        slots = self.cluster.get_slots_for_node(old_leader_id)
        for slot in slots:
            self.cluster.assign_slot(slot, self.node.id)

        self.node.role = 'leader'
        self.cluster.broadcast_topology_change()
```

### 7.4 Eviction Policies

**Candidate:** "When memory is full, the cache must decide which keys to evict. There are several strategies."

#### Approximated LRU (Least Recently Used)

```
Standard LRU requires a doubly-linked list across ALL keys — too expensive.
Redis uses approximated LRU: sample N random keys, evict the oldest.

  Standard LRU:                    Approximated LRU:
  ┌───┐  ┌───┐  ┌───┐  ┌───┐     Sample 5 random keys:
  │ A │─▶│ B │─▶│ C │─▶│ D │       key_x: last_access = 10:00:01
  └───┘  └───┘  └───┘  └───┘       key_m: last_access = 09:58:30  ← evict this
  most                 least        key_q: last_access = 10:00:05
  recent               recent       key_r: last_access = 09:59:45
                                    key_j: last_access = 10:00:02
  O(1) but needs linked list
  across millions of keys       Evict key_m (oldest in sample)
```

```python
import time
import random

class ApproximatedLRU:
    def __init__(self, max_memory, sample_size=5):
        self.max_memory = max_memory
        self.sample_size = sample_size
        self.store = {}          # key → value
        self.access_time = {}    # key → last_access_timestamp
        self.memory_used = 0

    def get(self, key):
        if key not in self.store:
            return None
        self.access_time[key] = time.monotonic()
        return self.store[key]

    def set(self, key, value):
        size = len(key) + len(str(value)) + 80  # approximate

        while self.memory_used + size > self.max_memory and self.store:
            self._evict()

        if key in self.store:
            old_size = len(key) + len(str(self.store[key])) + 80
            self.memory_used -= old_size

        self.store[key] = value
        self.access_time[key] = time.monotonic()
        self.memory_used += size

    def _evict(self):
        """Evict the least-recently-used key from a random sample."""
        if not self.store:
            return

        keys = list(self.store.keys())
        sample = random.sample(keys, min(self.sample_size, len(keys)))

        # Find the key with the oldest access time in the sample
        victim = min(sample, key=lambda k: self.access_time.get(k, 0))

        size = len(victim) + len(str(self.store[victim])) + 80
        del self.store[victim]
        del self.access_time[victim]
        self.memory_used -= size
```

#### LFU (Least Frequently Used) with Morris Counter

```python
import random
import math

class MorrisCounter:
    """
    Probabilistic counter: stores log2(count) in 8 bits.
    Actual count can be billions, but only uses 1 byte.
    Used by Redis for LFU eviction.
    """
    def __init__(self):
        self.log_counter = 5  # start at 5 (not 0, to give new keys a chance)

    def increment(self):
        """Probabilistically increment: P(increment) = 1 / (counter - 4)."""
        if self.log_counter >= 255:
            return
        probability = 1.0 / (self.log_counter - 4)  # decay factor
        if random.random() < probability:
            self.log_counter += 1

    def decrement_over_time(self, minutes_since_last_access, decay_time=1):
        """Halve the counter periodically so old popular keys decay."""
        decrement = minutes_since_last_access // decay_time
        self.log_counter = max(0, self.log_counter - decrement)

    @property
    def approximate_count(self):
        return 2 ** self.log_counter
```

#### Eviction Policy Comparison

| Policy | Best For | Drawback | How It Works |
|--------|----------|----------|--------------|
| **LRU** | General purpose, recency matters | Ignores frequency — a rarely-used key accessed once beats a popular key | Evict least recently accessed |
| **LFU** | Frequency matters (popular items stay) | Slow to adapt — old popular keys stick around | Evict least frequently accessed |
| **TTL** | Time-bound data (sessions, tokens) | Doesn't help with non-TTL keys | Evict keys closest to expiry |
| **Random** | Simple workloads, uniform access | May evict hot keys | Random eviction |
| **noeviction** | When data loss is unacceptable | Returns errors when full | Rejects new writes |
| **allkeys-lru** | Cache use case, any key evictable | N/A | LRU across all keys |
| **volatile-lru** | Only expire keys with TTL set | Keys without TTL never evicted | LRU only among keys with TTL |
| **allkeys-lfu** | Workloads with clear hot/cold split | Needs tuning for decay | LFU across all keys |

### 7.5 Cache Patterns

**Candidate:** "How the application interacts with the cache determines consistency and performance."

#### Cache-Aside (Lazy Loading)

```
Application ──GET──▶ Cache ──MISS──▶ Application ──READ──▶ Database
                                         │
                                    Write to cache
                                         │
                                    ◄──response──
```

```javascript
async function getUserCacheAside(userId) {
  // 1. Check cache first
  const cached = await cache.get(`user:${userId}`);
  if (cached) return JSON.parse(cached);

  // 2. Cache miss — read from database
  const user = await db.query('SELECT * FROM users WHERE id = ?', [userId]);

  // 3. Populate cache for next time
  await cache.set(`user:${userId}`, JSON.stringify(user), { ex: 3600 });

  return user;
}
```

**Pros:** Only requested data is cached; cache failures don't break reads.  
**Cons:** Cache miss = 3 round trips; stale data until TTL expires.

#### Write-Through

```
Application ──WRITE──▶ Cache ──WRITE──▶ Database
                  │
             ◄──OK──
```

```javascript
async function updateUserWriteThrough(userId, data) {
  // 1. Write to cache AND database atomically
  const user = { ...data, updated_at: Date.now() };

  await Promise.all([
    cache.set(`user:${userId}`, JSON.stringify(user), { ex: 3600 }),
    db.query('UPDATE users SET ? WHERE id = ?', [data, userId])
  ]);

  return user;
}
```

**Pros:** Cache is always up-to-date; no stale reads.  
**Cons:** Higher write latency; cache may hold data that's never read.

#### Write-Behind (Write-Back)

```
Application ──WRITE──▶ Cache ──OK──▶ Application
                          │
                   (async, batched)
                          │
                     ┌────▼────┐
                     │Database │
                     └─────────┘
```

```python
class WriteBehindCache:
    def __init__(self, cache, db, flush_interval=1.0):
        self.cache = cache
        self.db = db
        self.write_buffer = {}
        self.flush_interval = flush_interval
        self._start_flush_thread()

    def set(self, key, value):
        self.cache.set(key, value)
        self.write_buffer[key] = value  # buffer for async DB write

    def _flush_to_db(self):
        """Periodically flush buffered writes to database in batch."""
        if not self.write_buffer:
            return
        batch = dict(self.write_buffer)
        self.write_buffer.clear()
        self.db.batch_upsert(batch)  # single batch write
```

**Pros:** Lowest write latency; batches reduce DB load.  
**Cons:** Data loss risk if cache node fails before flush.

#### Read-Through

```
Application ──GET──▶ Cache ──MISS──▶ Cache ──READ──▶ Database
                                       │
                                  Cache stores it
                                       │
                                  ◄──response──
```

```javascript
class ReadThroughCache {
  constructor(cache, db) {
    this.cache = cache;
    this.db = db;
  }

  async get(key, loader) {
    const cached = await this.cache.get(key);
    if (cached) return JSON.parse(cached);

    // Cache handles the DB read internally
    const value = await loader(key);
    await this.cache.set(key, JSON.stringify(value), { ex: 3600 });
    return value;
  }
}

// Usage
const userCache = new ReadThroughCache(cache, db);
const user = await userCache.get(`user:${userId}`, (key) =>
  db.query('SELECT * FROM users WHERE id = ?', [key.split(':')[1]])
);
```

**Pros:** Simpler application code — cache manages its own population.  
**Cons:** Adds coupling; first request still slow.

#### Pattern Comparison

| Pattern | Read Latency | Write Latency | Consistency | Complexity |
|---------|-------------|---------------|-------------|------------|
| Cache-aside | Miss: high, Hit: low | N/A (app writes DB) | Eventual (stale until TTL) | Low |
| Write-through | Always low | Higher (sync DB write) | Strong | Medium |
| Write-behind | Always low | Lowest | Eventual (buffered) | High |
| Read-through | Miss: high, Hit: low | N/A | Eventual | Medium |

### 7.6 Hot Key Problem

**Candidate:** "A hot key is a key that receives disproportionate traffic, potentially overwhelming a single node."

```
Problem visualization:

  Normal distribution:          Hot key scenario:
  ┌─────┐ ┌─────┐ ┌─────┐     ┌─────┐ ┌██████┐ ┌─────┐
  │ 33K │ │ 33K │ │ 34K │     │ 10K │ │ 80K  │ │ 10K │
  │ops/s│ │ops/s│ │ops/s│     │ops/s│ │ops/s │ │ops/s│
  └──┬──┘ └──┬──┘ └──┬──┘     └──┬──┘ └──┬───┘ └──┬──┘
   Node1   Node2   Node3       Node1   Node2   Node3

  Example hot keys:
  - Celebrity tweet going viral → "tweet:12345"
  - Flash sale product → "product:iphone15"
  - Global config → "feature_flags"
```

**Detection:**

```python
class HotKeyDetector:
    def __init__(self, threshold_ops_per_sec=10000):
        self.threshold = threshold_ops_per_sec
        self.key_counters = {}
        self.window_start = time.time()

    def record_access(self, key):
        self.key_counters[key] = self.key_counters.get(key, 0) + 1

        if time.time() - self.window_start >= 1.0:
            hot_keys = {
                k: v for k, v in self.key_counters.items()
                if v > self.threshold
            }
            if hot_keys:
                self._alert_hot_keys(hot_keys)
            self.key_counters.clear()
            self.window_start = time.time()
```

**Solutions:**

```
1. Local Cache (L1 cache in application)
   App ──▶ Local HashMap (100ms TTL) ──miss──▶ Redis
   Reduces Redis load by 90%+ for hot keys

2. Key Replication (read from replicas)
   hot_key → read from random(leader, replica1, replica2)
   Spreads load across 3 nodes instead of 1

3. Key Splitting
   "product:123" → "product:123:shard:{0..9}"
   Writes: update all 10 shards
   Reads:  pick random shard → spreads across nodes
```

```javascript
class HotKeyMitigator {
  constructor(cache, localCacheTtlMs = 100) {
    this.cache = cache;
    this.localCache = new Map();
    this.localTtl = localCacheTtlMs;
  }

  async get(key) {
    // L1: Check local in-process cache
    const local = this.localCache.get(key);
    if (local && Date.now() - local.time < this.localTtl) {
      return local.value;
    }

    // L2: Fetch from distributed cache
    const value = await this.cache.get(key);
    this.localCache.set(key, { value, time: Date.now() });

    return value;
  }

  async getWithSplitting(key, shardCount = 10) {
    // Read from a random shard to spread load
    const shard = Math.floor(Math.random() * shardCount);
    return this.cache.get(`${key}:shard:${shard}`);
  }

  async setWithSplitting(key, value, shardCount = 10) {
    // Write to all shards
    const promises = [];
    for (let i = 0; i < shardCount; i++) {
      promises.push(this.cache.set(`${key}:shard:${i}`, value));
    }
    await Promise.all(promises);
  }
}
```

---

## 8. Scaling Strategies

### 8.1 Adding Nodes with Minimal Data Movement

**Candidate:** "When we add a new node to the cluster, we need to redistribute hash slots without downtime."

```
Before: 3 nodes, 16384 slots
  Node 1: [0-5460]      5461 slots
  Node 2: [5461-10922]  5462 slots
  Node 3: [10923-16383] 5461 slots

After: 4 nodes (Node 4 added)
  Node 1: [0-4095]       4096 slots  (-1365 slots)
  Node 2: [4096-8191]    4096 slots  (-1366 slots)
  Node 3: [8192-12287]   4096 slots  (-1365 slots)
  Node 4: [12288-16383]  4096 slots  (+4096 slots, from all 3)

Only ~25% of keys move (1/N where N=4)
```

### 8.2 Live Resharding (Slot Migration)

```
Slot Migration Process (zero downtime):

  Phase 1: Mark slot as MIGRATING/IMPORTING
  ┌────────┐                    ┌────────┐
  │Source   │  slot 1000:       │Target  │
  │Node     │  MIGRATING ──────▶│Node    │
  │         │                   │IMPORTING│
  └────────┘                    └────────┘

  Phase 2: Migrate keys (one by one or batched)
  ┌────────┐    MIGRATE key     ┌────────┐
  │Source   │──────────────────▶│Target  │
  │         │   (atomic: dump   │        │
  │         │    + restore +    │        │
  │         │    delete)        │        │
  └────────┘                    └────────┘

  Phase 3: During migration, handle in-flight requests
  Client ──GET key──▶ Source
    If key exists on Source → return value
    If key already migrated → return ASK redirect to Target

  Phase 4: Update slot ownership, broadcast to cluster
```

### 8.3 Client-Side Caching (Redis 6+ Tracking)

```
Server-Assisted Client Caching:

  Client                           Server
    │                                │
    │──GET user:1001────────────────▶│
    │◀─"Alice" (server tracks that  ─│
    │   client cached user:1001)     │
    │                                │
    │  (client caches locally)       │
    │                                │
    │  ... another client updates    │
    │  user:1001 ...                 │
    │                                │
    │◀─INVALIDATE user:1001─────────│
    │                                │
    │  (client evicts local copy)    │
    │                                │
    │──GET user:1001────────────────▶│
    │◀─"Bob"────────────────────────│

  Reduces network round-trips by 90%+ for read-heavy keys.
  Server maintains an "invalidation table" tracking which
  clients cached which keys.
```

### 8.4 Read Replicas for Read-Heavy Workloads

```
Read scaling with READONLY replicas:

  Writes: ──────────────────▶ Leader only
  Reads:  ──────────────────▶ Leader OR any Replica

  ┌──────────┐
  │  Leader   │ ◄── all writes
  │  Node 1   │
  └─────┬────┘
        │ replication
  ┌─────┼─────────────┐
  │     │             │
  ┌─────▼────┐ ┌──────▼───┐
  │ Replica  │ │ Replica  │
  │ 1a       │ │ 1b       │  ◄── reads distributed
  └──────────┘ └──────────┘

  Trade-off: reads from replicas may return slightly stale data
  (replication lag typically < 1ms under normal conditions)
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Node Failure

```
Timeline of automatic failover:

  T+0s:    Leader Node 3 crashes
  T+0.3s:  Replicas detect missed heartbeats → mark PFAIL
  T+1s:    Majority of leaders agree → mark FAIL
  T+1.5s:  Best replica (highest replication offset) starts election
  T+2s:    Majority of leaders vote → replica promoted to leader
  T+2.5s:  Cluster topology updated, clients redirect

  Data loss window:
  - Async replication: may lose writes in the last ~1 second
  - WAIT command (semi-sync): guaranteed N replicas acknowledged
```

| Failure Type | Detection Time | Recovery Time | Data Loss |
|-------------|---------------|---------------|-----------|
| Leader crash | 1-2 seconds | 2-5 seconds | Last ~1 second of writes |
| Replica crash | 1-2 seconds | Immediate (fewer replicas) | None |
| Network partition | 1-2 seconds | Varies | Depends on partition side |

### 9.2 Network Partition (Split-Brain Prevention)

```
Split-brain scenario:

  Partition A (minority)      │     Partition B (majority)
  ┌────────┐  ┌────────┐     │     ┌────────┐  ┌────────┐  ┌────────┐
  │Node 1  │  │Node 2  │     │     │Node 3  │  │Node 4  │  │Node 5  │
  │(leader)│  │(replica)│    █│█    │(replica)│  │(leader)│  │(leader)│
  └────────┘  └────────┘     │     └────────┘  └────────┘  └────────┘

  Prevention:
  1. Minority side: Node 1 can't reach majority → stops accepting writes
     (cluster-node-timeout + cluster-require-full-coverage)
  2. Majority side: Node 3 (replica of Node 1) gets promoted to leader
  3. When partition heals: Node 1 discovers higher epoch, demotes itself
     to replica, syncs from new leader (some writes may be lost)
```

```python
class SplitBrainProtection:
    def __init__(self, node, cluster):
        self.node = node
        self.cluster = cluster
        self.min_replicas_to_write = 1  # require at least 1 replica ACK

    def can_accept_writes(self):
        reachable_leaders = self.cluster.count_reachable_leaders()
        total_leaders = self.cluster.total_leader_count()

        # Refuse writes if we can't reach majority
        if reachable_leaders < total_leaders // 2 + 1:
            return False

        # Refuse writes if no replicas are reachable (data safety)
        reachable_replicas = self.cluster.count_reachable_replicas(self.node.id)
        if reachable_replicas < self.min_replicas_to_write:
            return False

        return True
```

### 9.3 Memory Exhaustion

```
Memory pressure cascade:

  Used Memory:  ████████████████████░░░░  85%  → Normal
  Used Memory:  ██████████████████████░░  92%  → Warning alert
  Used Memory:  ████████████████████████  100% → Eviction kicks in

  Eviction process:
  1. maxmemory reached → trigger eviction policy
  2. Sample keys → evict according to policy (LRU/LFU/TTL)
  3. If noeviction policy → return OOM error to clients
  4. Continue serving after freeing memory

  Memory fragmentation:
  - Ratio > 1.5 means significant fragmentation
  - Redis 4.0+ has active defragmentation (jemalloc)
  - Monitor: INFO MEMORY → mem_fragmentation_ratio
```

### 9.4 Cache Stampede (Thundering Herd)

**Candidate:** "When a popular key expires, hundreds of concurrent requests all miss the cache and hit the database simultaneously."

```
The problem:

  T=0: Key "product:123" expires (TTL reached)

  T=0.001: Request 1 ──cache miss──▶ DB query ──────────┐
  T=0.002: Request 2 ──cache miss──▶ DB query ───────┐  │
  T=0.003: Request 3 ──cache miss──▶ DB query ────┐  │  │
  T=0.004: Request 4 ──cache miss──▶ DB query ──┐ │  │  │
  ...                                            │ │  │  │
  T=0.050: Request 100 ──cache miss──▶ DB ──┐   │ │  │  │
                                             │   │ │  │  │
                                             ▼   ▼ ▼  ▼  ▼
                                          DATABASE OVERLOADED
```

**Solutions:**

```javascript
// Solution 1: Distributed lock (only one request populates cache)
async function getWithLock(key) {
  let value = await cache.get(key);
  if (value) return JSON.parse(value);

  const lockKey = `lock:${key}`;
  const acquired = await cache.set(lockKey, '1', { nx: true, ex: 10 });

  if (acquired) {
    // Winner: fetch from DB and populate cache
    const data = await db.query(key);
    await cache.set(key, JSON.stringify(data), { ex: 3600 });
    await cache.del(lockKey);
    return data;
  } else {
    // Loser: wait and retry from cache
    await sleep(50);
    return getWithLock(key); // retry
  }
}

// Solution 2: Probabilistic early expiration
async function getWithEarlyRefresh(key, ttl = 3600, beta = 1.0) {
  const result = await cache.get(key);
  if (!result) return fetchAndCache(key, ttl);

  const { value, expiry, computeTime } = JSON.parse(result);
  const now = Date.now();

  // XFetch: probabilistically refresh before actual TTL
  // P(refresh) increases as we approach expiry
  const remaining = (expiry - now) / 1000;
  const randomThreshold = computeTime * beta * Math.log(Math.random());

  if (remaining + randomThreshold <= 0) {
    // Proactively refresh (non-blocking)
    fetchAndCache(key, ttl); // fire-and-forget
  }

  return value;
}

// Solution 3: Stale-while-revalidate
async function getStaleWhileRevalidate(key) {
  const result = await cache.get(key);
  if (!result) return fetchAndCache(key, 3600);

  const { value, softExpiry } = JSON.parse(result);

  if (Date.now() > softExpiry) {
    // Serve stale, refresh in background
    fetchAndCache(key, 3600); // non-blocking
  }

  return value; // always returns immediately
}
```

### 9.5 Hot Key Melting a Node

```
Symptoms:
  - Single node CPU at 100%
  - Latency spike on that node (p99 > 100ms)
  - Other nodes healthy

Emergency response:
  1. Identify: SLOWLOG GET 10, redis-cli --hotkeys
  2. Immediate: Add local caching in application (100ms TTL)
  3. Short-term: Replicate hot key to all read replicas, spread reads
  4. Long-term: Implement key splitting, client-side caching
```

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

```
┌──────────────────────────────────────────────────────────────────┐
│                    CACHE MONITORING DASHBOARD                    │
├──────────────────────────────┬───────────────────────────────────┤
│  Hit Rate          98.5%     │  Memory Usage     4.2 / 8.0 GB   │
│  Miss Rate          1.5%     │  Memory Frag.     1.12 ratio      │
│  Eviction Rate    120/sec    │  Connected Clients  2,340         │
│  Ops/sec          95,420     │  Blocked Clients      3           │
├──────────────────────────────┼───────────────────────────────────┤
│  Keyspace                    │  Replication                      │
│  Total Keys    1,023,456     │  Role: leader                     │
│  Expires       456,789       │  Replicas: 2                      │
│  Avg TTL       1,842s        │  Repl. Lag: 0.2ms                │
├──────────────────────────────┼───────────────────────────────────┤
│  Latency (ms)                │  Network                          │
│  GET p50: 0.12  p99: 0.45   │  Input:  12.3 MB/s               │
│  SET p50: 0.15  p99: 0.52   │  Output: 45.6 MB/s               │
│  DEL p50: 0.10  p99: 0.38   │  Rejected Connections: 0          │
└──────────────────────────────┴───────────────────────────────────┘
```

### 10.2 Alerting Rules

| Metric | Warning | Critical | Action |
|--------|---------|----------|--------|
| Hit rate | < 90% | < 80% | Review eviction policy, increase memory |
| Memory usage | > 80% | > 90% | Scale up or add nodes |
| Eviction rate | > 1K/sec | > 10K/sec | Increase memory, review TTLs |
| Replication lag | > 100ms | > 1 second | Check network, reduce write rate |
| Connected clients | > 80% max | > 95% max | Increase maxclients, check leaks |
| Command latency p99 | > 5ms | > 50ms | Check slow log, hot keys |
| Fragmentation ratio | > 1.5 | > 2.0 | Enable active defrag, restart |

### 10.3 Slow Log Analysis

```python
class SlowLogAnalyzer:
    def __init__(self, cache_client, threshold_us=10000):
        self.client = cache_client
        self.threshold = threshold_us

    def analyze(self):
        slow_entries = self.client.execute('SLOWLOG', 'GET', 100)

        # Group by command type
        by_command = {}
        for entry in slow_entries:
            log_id, timestamp, duration_us, command = entry[:4]
            cmd_name = command[0]
            if cmd_name not in by_command:
                by_command[cmd_name] = {'count': 0, 'total_us': 0, 'max_us': 0}
            by_command[cmd_name]['count'] += 1
            by_command[cmd_name]['total_us'] += duration_us
            by_command[cmd_name]['max_us'] = max(
                by_command[cmd_name]['max_us'], duration_us
            )

        return {
            cmd: {
                **stats,
                'avg_us': stats['total_us'] / stats['count']
            }
            for cmd, stats in sorted(
                by_command.items(),
                key=lambda x: x[1]['total_us'],
                reverse=True
            )
        }
```

### 10.4 Health Check Endpoint

```javascript
app.get('/health/cache', async (req, res) => {
  const checks = {};

  // Connectivity
  try {
    const start = process.hrtime.bigint();
    await cache.ping();
    const latencyNs = Number(process.hrtime.bigint() - start);
    checks.connectivity = {
      status: latencyNs < 5_000_000 ? 'healthy' : 'degraded',
      latency_ms: latencyNs / 1_000_000
    };
  } catch (e) {
    checks.connectivity = { status: 'unhealthy', error: e.message };
  }

  // Memory
  const info = await cache.info('memory');
  const usedPct = info.used_memory / info.maxmemory * 100;
  checks.memory = {
    status: usedPct < 80 ? 'healthy' : usedPct < 90 ? 'warning' : 'critical',
    used_pct: usedPct.toFixed(1),
    fragmentation_ratio: info.mem_fragmentation_ratio
  };

  // Replication
  const replInfo = await cache.info('replication');
  checks.replication = {
    status: replInfo.connected_slaves > 0 ? 'healthy' : 'degraded',
    connected_replicas: replInfo.connected_slaves,
    repl_lag_ms: replInfo.master_repl_offset - replInfo.slave_repl_offset
  };

  const overallStatus = Object.values(checks).every(c => c.status === 'healthy')
    ? 'healthy' : 'degraded';
  res.json({ status: overallStatus, checks });
});
```

---

## 11. Advanced Features

### 11.1 Distributed Locks (Redlock Algorithm)

**Candidate:** "Distributed locks in a cache cluster must handle node failures. The Redlock algorithm acquires locks across a majority of independent nodes."

```
Redlock Algorithm (lock across N=5 independent nodes):

  Client                    Redis 1  Redis 2  Redis 3  Redis 4  Redis 5
    │──SET lock:res NX EX 10──▶ ✓
    │──SET lock:res NX EX 10──────────▶ ✓
    │──SET lock:res NX EX 10──────────────────▶ ✗ (timeout)
    │──SET lock:res NX EX 10──────────────────────────▶ ✓
    │──SET lock:res NX EX 10──────────────────────────────────▶ ✓
    │
    │  Acquired on 4/5 nodes (> N/2 + 1 = 3) → LOCK ACQUIRED
    │  Lock validity = TTL - elapsed_time
    │
    │  ... do critical section work ...
    │
    │──DEL lock:res──▶ all 5 nodes (release)
```

```python
import time
import uuid

class RedLock:
    def __init__(self, redis_instances, ttl_ms=10000):
        self.instances = redis_instances
        self.ttl = ttl_ms
        self.quorum = len(redis_instances) // 2 + 1

    def acquire(self, resource):
        lock_value = str(uuid.uuid4())
        start = time.monotonic_ns()
        acquired = 0

        for instance in self.instances:
            try:
                if instance.set(f"lock:{resource}", lock_value,
                               nx=True, px=self.ttl):
                    acquired += 1
            except Exception:
                pass

        elapsed_ms = (time.monotonic_ns() - start) / 1_000_000
        validity_ms = self.ttl - elapsed_ms

        if acquired >= self.quorum and validity_ms > 0:
            return Lock(resource, lock_value, validity_ms)
        else:
            # Failed: release any acquired locks
            self._release_all(resource, lock_value)
            return None

    def release(self, lock):
        self._release_all(lock.resource, lock.value)

    def _release_all(self, resource, value):
        # Only delete if we still hold the lock (compare value)
        lua_script = """
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        else
            return 0
        end
        """
        for instance in self.instances:
            try:
                instance.eval(lua_script, 1, f"lock:{resource}", value)
            except Exception:
                pass
```

### 11.2 Rate Limiting

```python
class SlidingWindowRateLimiter:
    def __init__(self, cache, max_requests, window_seconds):
        self.cache = cache
        self.max_requests = max_requests
        self.window = window_seconds

    def is_allowed(self, client_id):
        key = f"ratelimit:{client_id}"
        now = time.time()
        window_start = now - self.window

        # Lua script for atomic sliding window
        lua = """
        local key = KEYS[1]
        local now = tonumber(ARGV[1])
        local window_start = tonumber(ARGV[2])
        local max_requests = tonumber(ARGV[3])
        local window = tonumber(ARGV[4])

        -- Remove old entries outside window
        redis.call('ZREMRANGEBYSCORE', key, '-inf', window_start)

        -- Count current entries
        local count = redis.call('ZCARD', key)

        if count < max_requests then
            -- Add current request
            redis.call('ZADD', key, now, now .. ':' .. math.random())
            redis.call('EXPIRE', key, window)
            return 1
        else
            return 0
        end
        """
        return bool(self.cache.eval(lua, 1, key, now, window_start,
                                     self.max_requests, self.window))
```

### 11.3 Session Storage

```javascript
class SessionStore {
  constructor(cache, ttlSeconds = 1800) {
    this.cache = cache;
    this.ttl = ttlSeconds;
  }

  async createSession(userId, metadata = {}) {
    const sessionId = crypto.randomUUID();
    const session = {
      userId,
      createdAt: Date.now(),
      ...metadata
    };
    await this.cache.set(
      `session:${sessionId}`,
      JSON.stringify(session),
      { ex: this.ttl }
    );
    return sessionId;
  }

  async getSession(sessionId) {
    const data = await this.cache.get(`session:${sessionId}`);
    if (!data) return null;

    // Sliding expiration: renew TTL on every access
    await this.cache.expire(`session:${sessionId}`, this.ttl);
    return JSON.parse(data);
  }

  async destroySession(sessionId) {
    await this.cache.del(`session:${sessionId}`);
  }
}
```

### 11.4 Pub/Sub for Real-Time Events

```
Pub/Sub Architecture:

  Publisher 1 ──PUBLISH "orders" msg──▶ ┌──────────┐
  Publisher 2 ──PUBLISH "orders" msg──▶ │  Cache    │
                                        │  Node     │──▶ Subscriber A
                                        │           │──▶ Subscriber B
  Publisher 3 ──PUBLISH "alerts" msg──▶ │ Channel   │──▶ Subscriber C
                                        │ Registry  │
                                        └──────────┘

  Note: Pub/Sub is fire-and-forget. If subscriber is disconnected,
  messages are lost. For durability, use Streams instead.
```

### 11.5 Redis Streams (Event Sourcing)

```python
class EventStream:
    """Redis Streams for durable, ordered event processing."""

    def __init__(self, cache, stream_name):
        self.cache = cache
        self.stream = stream_name

    def publish(self, event_type, data):
        """Append event to stream. Returns auto-generated ID."""
        return self.cache.xadd(self.stream, {
            'type': event_type,
            'data': json.dumps(data),
            'timestamp': str(time.time())
        })

    def consume(self, group, consumer, count=10):
        """Read events as part of a consumer group (exactly-once)."""
        messages = self.cache.xreadgroup(
            group, consumer, {self.stream: '>'}, count=count
        )
        results = []
        for stream_name, entries in messages:
            for msg_id, fields in entries:
                results.append((msg_id, fields))
                self.cache.xack(self.stream, group, msg_id)
        return results
```

### 11.6 Specialized Data Structures

```
HyperLogLog (cardinality estimation):
  PFADD visitors "user:1001"    # Add element
  PFADD visitors "user:1002"
  PFCOUNT visitors              # → 2 (approximate unique count)
  Memory: only 12 KB regardless of cardinality!
  Use case: unique visitors, unique search queries

Bloom Filter (membership test):
  BF.ADD filter "user:1001"     # Add element
  BF.EXISTS filter "user:1001"  # → 1 (might be present)
  BF.EXISTS filter "user:9999"  # → 0 (definitely not present)
  Use case: avoid expensive DB lookups for non-existent keys

Geospatial Indexing:
  GEOADD restaurants -73.856 40.848 "pizza_place"
  GEOADD restaurants -73.961 40.776 "sushi_bar"
  GEORADIUS restaurants -73.9 40.8 5 km  # find within 5km
  Use case: nearby search, delivery radius
```

---

## 12. Interview Q&A

### Q1: How does consistent hashing work and why is it important for caching?

**Candidate:** "Consistent hashing maps both keys and nodes onto a circular hash space (ring). Each key is assigned to the first node encountered when walking clockwise around the ring.

The critical advantage is during scaling: when adding or removing a node, only `K/N` keys need to move (K = total keys, N = total nodes), compared to nearly all keys with naive modulo hashing.

Virtual nodes (vnodes) solve the problem of uneven distribution — each physical node gets 100-200 positions on the ring, ensuring balanced load. Redis Cluster takes a slightly different approach with 16384 fixed hash slots, which gives explicit control over slot-to-node mapping and makes rebalancing more predictable."

**Interviewer:** "What happens if a node goes down mid-operation?"

**Candidate:** "The client library maintains a cached cluster topology. On connection failure or MOVED redirect, it refreshes the topology and retries. The key previously mapped to the failed node is now served by its replica (which was promoted). Ongoing requests see at most one redirect."

---

### Q2: Compare LRU vs LFU eviction — when would you use each?

**Candidate:** "LRU evicts the key that hasn't been accessed for the longest time. It's great for workloads with temporal locality — recent data is likely to be accessed again. But LRU has a critical weakness: a one-time scan of infrequently-used keys can evict popular keys.

LFU evicts the key accessed least frequently. It protects popular keys from being evicted by one-time access patterns. Redis uses a probabilistic counter (Morris counter) to track frequency in just 8 bits per key.

| Scenario | Best Policy |
|----------|-------------|
| General web caching, sessions | LRU |
| CDN/content with clear hot/cold | LFU |
| Mixed: some keys always popular | LFU |
| Workload changes rapidly | LRU |

In practice, Redis's `allkeys-lfu` policy works well for most production systems because it naturally adapts — the decay mechanism ensures old popular keys eventually get evicted if they're no longer accessed."

---

### Q3: How do you handle cache stampede (thundering herd)?

**Candidate:** "Cache stampede happens when a popular key expires and hundreds of concurrent requests simultaneously hit the database. Three main solutions:

1. **Distributed locking**: Use SETNX to let only one request rebuild the cache. Others wait and retry from cache. Simple but adds latency for waiting requests.

2. **Probabilistic early refresh (XFetch)**: Before TTL expires, each request has an increasing probability of refreshing the cache proactively. The formula `remaining_ttl + compute_time × β × ln(random())` ensures exactly one request typically refreshes before expiry.

3. **Stale-while-revalidate**: Set a soft TTL shorter than the hard TTL. After soft expiry, serve stale data while refreshing in the background. Users never see a cache miss.

I'd choose stale-while-revalidate for most production systems because it provides the best user experience — no request ever waits for a DB query."

---

### Q4: How does Redis achieve sub-millisecond latency with a single thread?

**Candidate:** "Redis's single-threaded model is actually an advantage, not a limitation:

1. **No lock contention**: All data structures are accessed without locks, eliminating the overhead of mutex/spinlock acquisition that multi-threaded systems face.

2. **I/O multiplexing (epoll/kqueue)**: A single thread handles thousands of connections using kernel event notification. The thread never blocks waiting for one client — it processes whichever connections have data ready.

3. **In-memory data**: All operations are memory-bound. A hash table lookup is ~100ns, far faster than any network I/O. The bottleneck is the network, not the CPU.

4. **Efficient data structures**: Purpose-built structures like SDS (Simple Dynamic Strings), ziplists, and skip lists are cache-line friendly and minimize memory allocations.

5. **Pipelining**: Clients can send multiple commands without waiting for responses, amortizing network round-trip overhead.

Since Redis 6.0, I/O threading handles network read/write in parallel threads while command execution remains single-threaded — this gives the best of both worlds."

---

### Q5: How would you implement distributed locks using a cache?

**Candidate:** "A naive `SETNX key value EX ttl` works for single-node Redis but fails in distributed scenarios. The Redlock algorithm solves this:

1. Generate a unique lock value (UUID).
2. Try to acquire the lock on N independent Redis instances (e.g., 5) with a short timeout.
3. If acquired on majority (N/2 + 1 = 3), the lock is valid. Lock validity = TTL minus elapsed acquisition time.
4. On release, delete the lock on all instances — but only if the value matches (preventing accidental release of someone else's lock).

Critical implementation details:
- **Clock drift**: Add a small clock drift factor to validity time.
- **Fencing tokens**: Use an incrementing counter alongside the lock to prevent race conditions during long pauses (GC, network delays).
- **Auto-release TTL**: Always set a TTL so locks don't persist forever if the holder crashes.

The main controversy around Redlock (raised by Martin Kleppmann) is that it relies on timing assumptions. For truly critical distributed locks, a consensus-based system like ZooKeeper or etcd is more reliable, though significantly slower."

---

### Q6: Cache-aside vs write-through — trade-offs?

**Candidate:** "

| Aspect | Cache-Aside | Write-Through |
|--------|------------|---------------|
| Read path | App checks cache → on miss, reads DB, populates cache | App reads cache (always populated) |
| Write path | App writes DB only. Cache populated on next read | App writes cache AND DB synchronously |
| Consistency | Stale until TTL expires or invalidation | Always consistent (cache = DB) |
| Write latency | Lower (DB only) | Higher (cache + DB) |
| Cache pollution | Only caches requested data | May cache data that's never read |
| Implementation | Simple, most popular | Needs middleware or proxy layer |

I'd recommend cache-aside for most systems because:
- It's simple and well-understood
- It naturally avoids caching unused data
- Cache failures are graceful (fall back to DB)

Write-through makes sense when you have strict read-after-write consistency requirements and can tolerate the extra write latency."

---

### Q7: How do you handle hot keys that overload a single node?

**Candidate:** "Hot keys are one of the most common production issues with distributed caches. My approach has three tiers:

**Detection**: Monitor per-key access frequency. Redis's `--hotkeys` flag, custom counters (INCR on a shadow key), or proxy-layer tracking.

**Immediate mitigation**:
- **L1 local cache**: Add a small in-process cache (HashMap with 50-200ms TTL) in the application. This absorbs 90%+ of hot key reads before they hit Redis.

**Structural solutions**:
- **Read from replicas**: Enable READONLY mode on replicas and distribute reads across leader + replicas. Spreads load across 3 nodes.
- **Key splitting**: Shard the hot key into N sub-keys (`product:123:shard:0` through `product:123:shard:9`). Reads pick a random shard; writes update all shards. This distributes the key across N different hash slots (and nodes).

**Redis 6+ solution**:
- **Client-side caching with server-assisted invalidation**: The server tracks which clients cached which keys and sends invalidation messages on changes. This reduces Redis traffic by 90%+ for hot keys."

---

### Q8: How do you ensure data consistency between cache and database?

**Candidate:** "This is the hardest problem in caching. Common strategies, ranked by consistency guarantee:

1. **TTL-based expiry** (weakest): Set a TTL and accept stale data for that duration. Simple, works for most read-heavy workloads where slight staleness is acceptable.

2. **Cache invalidation on write**: After writing to DB, immediately delete the cache key. The next read triggers a cache miss and repopulates from DB. The race condition: if read happens between DB write and cache invalidation.

3. **Double-delete pattern**: Delete cache → write DB → wait small delay → delete cache again. The second delete catches the race condition where a stale read repopulated the cache.

```
DELETE cache(key)  →  UPDATE db(key)  →  sleep(500ms)  →  DELETE cache(key)
```

4. **CDC (Change Data Capture)**: Use database binlog (MySQL) or WAL (PostgreSQL) to stream changes to a consumer that invalidates the cache. This decouples the application from cache invalidation and handles all write paths including direct DB modifications.

5. **Write-through** (strongest): Write to cache and DB atomically (or at least synchronously). Guarantees cache is always current but has higher write latency.

For most production systems, I recommend cache invalidation on write with a short TTL as a safety net, plus CDC for comprehensive coverage of all write paths."

---

## 13. Production Checklist

### Pre-Launch

- [ ] Memory limit set (`maxmemory`) with appropriate eviction policy
- [ ] Persistence configured (RDB + AOF for important data, disabled for pure cache)
- [ ] Replication set up (1 leader + 2 replicas per shard)
- [ ] Connection pool configured in client (min/max connections, timeout)
- [ ] Timeout settings: command timeout, connection timeout, read timeout
- [ ] Disable dangerous commands in production (`KEYS *`, `FLUSHALL`, `DEBUG`)
- [ ] TLS encryption enabled for data in transit
- [ ] AUTH password / ACL configured
- [ ] Benchmark with realistic data and traffic patterns
- [ ] Failover tested: kill a leader, verify automatic promotion

### Day-1

- [ ] Monitor hit rate (target > 95%)
- [ ] Monitor eviction rate (should be near zero initially)
- [ ] Verify replication lag < 1ms under load
- [ ] Check memory fragmentation ratio < 1.5
- [ ] Slow log enabled (`slowlog-log-slower-than 10000`)
- [ ] Alerting configured for all critical metrics
- [ ] Verify client retry logic handles MOVED/ASK redirects
- [ ] Connection pool health verified (no leaks)

### Week-1

- [ ] Review slow log for unexpected patterns
- [ ] Analyze key size distribution (`DEBUG OBJECT` sampling)
- [ ] Validate TTL distribution (no mass-expiry at same time)
- [ ] Test scaling: add a node, verify live slot migration
- [ ] Verify backup/restore procedure (RDB snapshot + AOF replay)
- [ ] Load test to 2× expected peak traffic
- [ ] Document runbook for common failure scenarios

### Month-1

- [ ] Capacity planning review: memory trend, ops/sec growth
- [ ] Optimize: identify and split hot keys
- [ ] Review eviction policy effectiveness (LRU vs LFU tuning)
- [ ] Evaluate client-side caching for frequently-read keys
- [ ] Security audit: review ACL rules, network access
- [ ] Update to latest stable version (if applicable)
- [ ] Disaster recovery drill: simulate full cluster failure + restore

---

## Summary

| Aspect | Design Decision | Rationale |
|--------|----------------|-----------|
| Partitioning | 16384 hash slots | Fixed slots enable explicit assignment and predictable migration |
| Replication | Async leader-follower | Sub-ms writes; WAIT for stronger guarantees when needed |
| Eviction | Approximated LRU/LFU | Near-optimal eviction without maintaining expensive linked lists |
| Consistency | AP (availability > consistency) | Cache tolerates staleness; DB is source of truth |
| Persistence | RDB + AOF hybrid | RDB for fast restarts; AOF for minimal data loss |
| Thread Model | Single-threaded event loop | Eliminates locking overhead; network is the bottleneck |
| Failure Detection | Gossip protocol + quorum | Distributed detection without single point of failure |
| Client Routing | Cluster-aware client library | Smart routing avoids proxy bottleneck |
| Hot Key Mitigation | L1 local cache + key splitting | Multi-layer defense against traffic concentration |
| Stampede Prevention | Stale-while-revalidate + locks | Zero-latency reads with background refresh |

### Scalability Path

```
Phase 1: Single Node (MVP)
  1 node, 10K ops/sec, 5GB data
  Simple key-value, TTL, basic persistence

Phase 2: Replicated (High Availability)
  1 leader + 2 replicas
  Automatic failover, read scaling
  50K ops/sec, 5GB data

Phase 3: Clustered (Horizontal Scale)
  10 nodes × 3 (leader + replicas) = 30 instances
  1M ops/sec, 50GB data
  Hash slot partitioning, live resharding

Phase 4: Global (Multi-Region)
  100 nodes across 3 regions
  10M ops/sec, 500GB data
  Cross-region replication, local reads
  Client-side caching, L1 local cache

Phase 5: Platform (Cache-as-a-Service)
  1000+ nodes, multi-tenant
  Auto-scaling based on traffic
  Tiered storage (hot in-memory, warm on SSD)
  Full observability platform
```

---

> **Interview Tip:** When discussing a distributed cache in a system design interview, always start with the single-node architecture (event loop, data structures, eviction) before scaling to distributed concerns (hashing, replication, partitioning). This shows depth of understanding and a methodical approach to complexity.
