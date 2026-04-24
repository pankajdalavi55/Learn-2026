# Complete System Design: Metrics / Monitoring System (Production-Ready)

> **Complexity Level:** Advanced  
> **Estimated Time:** 45-60 minutes in interview  
> **Real-World Examples:** Datadog, Prometheus + Grafana, New Relic, CloudWatch, InfluxDB

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Time-Series Database Engine](#7-deep-dive-time-series-database-engine)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design a metrics collection and monitoring system like Datadog that can ingest, store, query, and alert on time-series data from thousands of services."

**Interviewer's Perspective:**  
This is a top-tier infrastructure design problem used to assess:
- **Time-series database design** — purpose-built storage engines, compression, write optimization
- **Data pipeline at scale** — handling millions of data points per second with sub-second latency
- **Aggregation strategies** — downsampling, rollups, pre-aggregation for fast queries
- **Alerting engine design** — reliable rule evaluation, state machines, avoiding alert fatigue
- **High cardinality handling** — the #1 operational challenge in production monitoring systems
- **Write-heavy system design** — 10,000:1 write-to-read ratio with durability guarantees

This problem reveals whether a candidate truly understands distributed systems internals or just knows surface-level patterns.

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Before I start designing, I'd like to clarify the scope and constraints. Can I ask a few questions?"

**Interviewer:** "Of course, go ahead."

**Candidate:** "What types of metrics should the system support? Just infrastructure metrics like CPU and memory, or also custom application metrics?"

**Interviewer:** "Both. Infrastructure metrics from agents running on hosts, and custom application metrics pushed by services. Think of it like Datadog — the system should handle CPU, memory, disk, network from agents, plus application-level metrics like request latency, error rates, and business KPIs."

**Candidate:** "What metric types do we need? Just simple gauges or also counters, histograms, and summaries?"

**Interviewer:** "All four: counters (monotonically increasing values like total requests), gauges (point-in-time values like CPU usage), histograms (distribution of values in buckets), and summaries (pre-computed quantiles). This is similar to what Prometheus supports."

**Candidate:** "For the query side — should we support a flexible query language, or just simple key-value lookups?"

**Interviewer:** "A flexible query language similar to PromQL. Users should be able to filter by labels, apply aggregations like avg, sum, max, min, percentiles, and group by dimensions. Think: `avg(cpu_usage{service='api', region='us-east'}) by (host)` over the last hour."

**Candidate:** "What about alerting? Simple threshold alerts, or do we also need anomaly detection?"

**Interviewer:** "Start with threshold-based alerts — for example, 'alert when avg CPU > 90% for 5 minutes.' Design the alerting engine to be extensible so we can add anomaly detection later."

**Candidate:** "For data retention — how long do we keep data at full resolution vs. downsampled?"

**Interviewer:** "15 days at full resolution, then downsample to 1-minute intervals for up to 1 year. After 1 year, data can be archived or deleted."

**Candidate:** "What's the collection model — push-based where services send metrics, or pull-based where we scrape endpoints?"

**Interviewer:** "Support both. Push via agents and SDKs, pull via scraping endpoints like Prometheus does. The ingestion layer should be protocol-agnostic."

**Candidate:** "Perfect. Let me summarize the requirements."

### 2.1 Functional Requirements

| # | Requirement | Description |
|---|------------|-------------|
| FR-1 | Metric Ingestion | Ingest metrics from agents and services (CPU, memory, custom metrics) |
| FR-2 | Time-Series Storage | Store data with labels/tags, support metric types: counter, gauge, histogram, summary |
| FR-3 | Flexible Querying | Query with aggregation functions (avg, sum, max, min, percentiles) over time windows |
| FR-4 | Dashboard Visualization | Expose APIs for dashboard rendering with customizable graphs and panels |
| FR-5 | Alerting Rules | Threshold-based alerts with configurable conditions and notification routing |
| FR-6 | Label-Based Filtering | Filter and group metrics by arbitrary key-value labels |
| FR-7 | Downsampling | Automatically roll up fine-grained data into coarser resolutions over time |

### 2.2 Non-Functional Requirements

| # | Requirement | Target |
|---|------------|--------|
| NFR-1 | Ingestion Throughput | Handle 10M data points/sec |
| NFR-2 | Query Latency (Recent) | < 1s for last 1 hour of data |
| NFR-3 | Query Latency (Historical) | < 10s for last 30 days of data |
| NFR-4 | Alerting Availability | 99.9% — missed alerts are unacceptable |
| NFR-5 | Data Retention | 15 days full resolution, 1 year downsampled |
| NFR-6 | Durability | No data loss once acknowledged by ingestion pipeline |
| NFR-7 | Availability | 99.95% for ingestion and query paths |

### 2.3 Scale Parameters

| Parameter | Value |
|-----------|-------|
| Services emitting metrics | 10,000 |
| Metrics per service | ~1,000 |
| Unique time series | ~10 million |
| Data points ingested per second | 10 million |
| Dashboard queries per second | ~1,000 |
| Write-to-read ratio | ~10,000:1 |
| Labels per metric | 5-10 on average |

---

## 3. Scale Estimation

### 3.1 Traffic Estimation

**Candidate:** "Let me work through the numbers to size our infrastructure."

```
Ingestion:
  10,000 services × 1,000 metrics/service = 10M unique time series
  Each series reports every 10 seconds → 10M points / 10s = 1M points/sec baseline
  With sub-second reporting + burst → target 10M points/sec peak

  10M points/sec × 86,400 sec/day = 864 billion data points/day

Queries:
  ~1,000 dashboard queries/sec
  Each query scans 10-100 time series over a time range
  Peak during incidents: 5× → 5,000 queries/sec

Alerting:
  ~50,000 active alert rules
  Each rule evaluates every 15-60 seconds
  ~2,000 rule evaluations/sec
```

### 3.2 Storage Estimation

**Candidate:** "Storage is the critical constraint for a time-series system."

```
Raw data per point:
  timestamp:  8 bytes
  value:      8 bytes (float64)
  overhead:   ~8 bytes (series reference, alignment)
  Total raw:  ~24 bytes per data point (before compression)

Daily raw storage:
  10M points/sec × 24 bytes × 86,400 sec/day = ~20.7 TB/day (uncompressed)

With Gorilla compression (~12:1 ratio):
  ~1.7 TB/day compressed

15-day full resolution:
  1.7 TB/day × 15 days = ~25.5 TB

Downsampled data (1-minute resolution for 1 year):
  10M series × 1 point/min × 60 min/hr × 24 hr × 365 days × 2 bytes
  = ~10.5 TB/year (heavily compressed rollups)

Label index storage:
  10M series × ~200 bytes/series (label pairs + index entries) = ~2 TB

Total active storage: ~40 TB
```

### 3.3 Bandwidth Estimation

```
Ingestion bandwidth:
  10M points/sec × 100 bytes/point (with labels in wire format) = 1 GB/sec inbound
  With batching and compression: ~200 MB/sec actual network

Query bandwidth:
  1,000 queries/sec × 50 KB avg response = 50 MB/sec outbound
```

### 3.4 Infrastructure Summary

| Resource | Estimate |
|----------|----------|
| Ingestion bandwidth | ~200 MB/sec (compressed) |
| Daily storage (compressed) | ~1.7 TB |
| Total active storage | ~40 TB |
| Kafka partitions | 256-512 |
| TSDB nodes | 30-50 (with replication) |
| Ingestion workers | 50-100 |
| Query nodes | 10-20 |

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          METRICS COLLECTION LAYER                          │
│                                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  Host     │  │  App     │  │ Container│  │ Cloud    │  │ Custom   │    │
│  │  Agent    │  │  SDK     │  │  Agent   │  │ Integr.  │  │ Exporter │    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘    │
│       │              │              │              │              │          │
│       │   StatsD     │   HTTP/gRPC  │  Prometheus  │    OTLP     │  Push    │
└───────┼──────────────┼──────────────┼──────────────┼──────────────┼──────────┘
        │              │              │              │              │
        ▼              ▼              ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       INGESTION GATEWAY (Load Balanced)                     │
│                                                                             │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │  Protocol Parsers │ Validation │ Rate Limiting │ Label Normalization  │  │
│  └───────────────────────────────────────────────────────────────────────┘  │
│                              │                                              │
│               Batch + Compress + Partition by metric hash                   │
└──────────────────────────────┼──────────────────────────────────────────────┘
                               ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                        KAFKA (Write Buffer / Decoupler)                     │
│                                                                             │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐        256-512          │
│  │ Part 0  │ │ Part 1  │ │ Part 2  │ │ Part N  │ ◄──── partitions        │
│  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘                          │
└───────┼──────────┼──────────┼──────────┼────────────────────────────────────┘
        │          │          │          │
        ▼          ▼          ▼          ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                         INGESTION WORKERS (50-100)                          │
│                                                                             │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Deserialize │ Resolve Series ID │ Write to WAL │ Buffer in Memory    │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────┼──────────────────────────────────────────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         ▼                     ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   TSDB Node 1   │  │   TSDB Node 2   │  │   TSDB Node N   │
│                 │  │                 │  │                 │
│ ┌─────────────┐ │  │ ┌─────────────┐ │  │ ┌─────────────┐ │
│ │ WAL         │ │  │ │ WAL         │ │  │ │ WAL         │ │
│ ├─────────────┤ │  │ ├─────────────┤ │  │ ├─────────────┤ │
│ │ Head Block  │ │  │ │ Head Block  │ │  │ │ Head Block  │ │
│ │ (in-memory) │ │  │ │ (in-memory) │ │  │ │ (in-memory) │ │
│ ├─────────────┤ │  │ ├─────────────┤ │  │ ├─────────────┤ │
│ │ Chunk Files │ │  │ │ Chunk Files │ │  │ │ Chunk Files │ │
│ │ (on-disk)   │ │  │ │ (on-disk)   │ │  │ │ (on-disk)   │ │
│ ├─────────────┤ │  │ ├─────────────┤ │  │ ├─────────────┤ │
│ │ Inverted    │ │  │ │ Inverted    │ │  │ │ Inverted    │ │
│ │ Index       │ │  │ │ Index       │ │  │ │ Index       │ │
│ └─────────────┘ │  │ └─────────────┘ │  │ └─────────────┘ │
└─────────────────┘  └─────────────────┘  └─────────────────┘
         │                     │                     │
         └─────────┬───────────┘                     │
                   │                                 │
         ┌─────────▼─────────────────────────────────▼──────┐
         │               QUERY ENGINE                       │
         │                                                   │
         │  ┌──────────┐  ┌──────────┐  ┌───────────────┐   │
         │  │ Query    │  │ Shard    │  │ Aggregation   │   │
         │  │ Parser   │──│ Router   │──│ & Merge       │   │
         │  └──────────┘  └──────────┘  └───────────────┘   │
         └───────────┬──────────────────────────────────────┘
                     │
    ┌────────────────┼────────────────┐
    ▼                ▼                ▼
┌────────┐   ┌────────────┐   ┌────────────────┐
│  API   │   │ Dashboard  │   │   ALERTING     │
│ Server │   │ Service    │   │   ENGINE        │
│        │   │ (Grafana)  │   │                │
│        │   │            │   │ ┌────────────┐ │
│        │   │            │   │ │ Rule Eval  │ │
│        │   │            │   │ │ Loop       │ │
│        │   │            │   │ ├────────────┤ │
│        │   │            │   │ │ State      │ │
│        │   │            │   │ │ Machine    │ │
│        │   │            │   │ ├────────────┤ │
│        │   │            │   │ │ Notifier   │ │
│        │   │            │   │ └────────────┘ │
└────────┘   └────────────┘   └────────────────┘
                                      │
                     ┌────────────────┼────────────────┐
                     ▼                ▼                ▼
               ┌──────────┐   ┌──────────┐   ┌──────────┐
               │ PagerDuty│   │  Slack   │   │  Email   │
               └──────────┘   └──────────┘   └──────────┘

Background Services:
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│  Downsampling    │  │   Compaction     │  │  Retention       │
│  Service         │  │   Service        │  │  Manager         │
│                  │  │                  │  │                  │
│ 10s → 1min      │  │ Merge small      │  │ Delete expired   │
│ 1min → 5min     │  │ chunks into      │  │ blocks beyond    │
│ 5min → 1hr      │  │ larger blocks    │  │ retention window │
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

### 4.2 API Design

**Candidate:** "Here are the core APIs the system exposes."

#### Metric Ingestion API

```
POST /api/v1/metrics
Content-Type: application/json

{
  "metrics": [
    {
      "name": "http_requests_total",
      "type": "counter",
      "value": 1,
      "timestamp": 1714000000,
      "labels": {
        "service": "api-gateway",
        "method": "GET",
        "status": "200",
        "region": "us-east-1"
      }
    },
    {
      "name": "cpu_usage_percent",
      "type": "gauge",
      "value": 73.5,
      "timestamp": 1714000000,
      "labels": {
        "host": "web-server-042",
        "region": "us-east-1"
      }
    }
  ]
}

Response: 202 Accepted
{
  "accepted": 2,
  "rejected": 0
}
```

#### Query API

```
GET /api/v1/query?expr=avg(cpu_usage_percent{service="api"})+by+(host)&start=1714000000&end=1714003600&step=60

Response: 200 OK
{
  "status": "success",
  "data": {
    "result_type": "matrix",
    "result": [
      {
        "labels": { "host": "web-01" },
        "values": [[1714000000, 72.3], [1714000060, 74.1], ...]
      },
      {
        "labels": { "host": "web-02" },
        "values": [[1714000000, 65.8], [1714000060, 67.2], ...]
      }
    ]
  }
}
```

#### Alert Rules API

```
POST /api/v1/alerts/rules
Content-Type: application/json

{
  "name": "HighCPU",
  "expr": "avg(cpu_usage_percent{service='api'}) > 90",
  "for": "5m",
  "severity": "critical",
  "annotations": {
    "summary": "CPU usage above 90% for 5 minutes",
    "runbook": "https://wiki.internal/runbooks/high-cpu"
  },
  "notifications": [
    { "channel": "pagerduty", "routing_key": "P1-infra" },
    { "channel": "slack", "webhook": "#alerts-critical" }
  ]
}
```

### 4.3 Data Flow

**Candidate:** "The data flows through three primary paths."

```
WRITE PATH (hot path — latency critical):
  Service → Agent → Gateway → Kafka → Worker → WAL + Memory Buffer → Disk Chunk

QUERY PATH (warm path — sub-second for recent data):
  Dashboard → API → Query Engine → [Memory Buffer + Disk Chunks] → Aggregate → Response

ALERT PATH (critical path — must never miss):
  Rule Evaluator → Query Engine → Compare threshold → State Machine → Notification
```

---

## 5. Deep Dive: Core Components

### 5.1 Ingestion Gateway

**Candidate:** "The ingestion gateway is our front door. It needs to handle multiple protocols and protect the system from overload."

```
┌─────────────────────────────────────────────────────────────────┐
│                      INGESTION GATEWAY                          │
│                                                                 │
│  ┌────────────┐  ┌────────────┐  ┌─────────────┐  ┌─────────┐ │
│  │  StatsD    │  │ Prometheus │  │ OpenTelemetry│  │  HTTP   │ │
│  │  UDP/TCP   │  │ Remote     │  │  OTLP/gRPC  │  │  JSON   │ │
│  │  Parser    │  │ Write      │  │  Receiver    │  │  API    │ │
│  └─────┬──────┘  └─────┬──────┘  └──────┬──────┘  └────┬────┘ │
│        │               │                │              │       │
│        └───────────┬───┴────────────────┴──────┬───────┘       │
│                    ▼                           ▼               │
│          ┌─────────────────┐        ┌──────────────────┐       │
│          │ Label           │        │  Rate Limiter    │       │
│          │ Normalization   │        │  (per tenant)    │       │
│          │ & Validation    │        │                  │       │
│          └────────┬────────┘        └────────┬─────────┘       │
│                   │                          │                 │
│                   └──────────┬───────────────┘                 │
│                              ▼                                 │
│                   ┌──────────────────┐                         │
│                   │  Batch Builder   │                         │
│                   │  (group by       │                         │
│                   │   partition key) │                         │
│                   └────────┬─────────┘                         │
│                            ▼                                   │
│                   ┌──────────────────┐                         │
│                   │  Kafka Producer  │                         │
│                   │  (async batched) │                         │
│                   └──────────────────┘                         │
└─────────────────────────────────────────────────────────────────┘
```

```javascript
// Gateway: Protocol-agnostic metric ingestion handler
class IngestionGateway {
  constructor(kafkaProducer, rateLimiter, validator) {
    this.producer = kafkaProducer;
    this.rateLimiter = rateLimiter;
    this.validator = validator;
    this.batchBuffer = new Map(); // partition → batch
    this.flushIntervalMs = 100;
  }

  async handleIngest(tenantId, rawMetrics, protocol) {
    if (!this.rateLimiter.allow(tenantId)) {
      return { status: 429, rejected: rawMetrics.length };
    }

    const parser = this.getParser(protocol);
    const metrics = parser.parse(rawMetrics);

    let accepted = 0, rejected = 0;
    for (const metric of metrics) {
      if (!this.validator.validate(metric)) {
        rejected++;
        continue;
      }
      this.normalizeLabels(metric);
      const partitionKey = this.computePartition(metric);
      this.addToBatch(partitionKey, metric);
      accepted++;
    }

    return { status: 202, accepted, rejected };
  }

  normalizeLabels(metric) {
    // Sort label keys for consistent series ID hashing
    const sorted = {};
    for (const key of Object.keys(metric.labels).sort()) {
      sorted[key] = metric.labels[key].trim().toLowerCase();
    }
    metric.labels = sorted;
    metric.seriesId = this.hashSeries(metric.name, sorted);
  }

  computePartition(metric) {
    // Consistent hashing ensures same series always goes to same TSDB shard
    return murmurhash3(metric.seriesId) % this.numPartitions;
  }

  hashSeries(name, labels) {
    const labelStr = Object.entries(labels)
      .map(([k, v]) => `${k}=${v}`)
      .join(',');
    return fnv1aHash(`${name}{${labelStr}}`);
  }
}
```

### 5.2 Kafka as Write Buffer

**Candidate:** "Kafka sits between ingestion and storage. It absorbs burst traffic, provides durability, and decouples producers from consumers."

```
Why Kafka:
┌─────────────────────────────────────────────────────────────┐
│ Problem                    │ How Kafka Solves It            │
├────────────────────────────┼────────────────────────────────┤
│ Burst traffic (10× spike)  │ Buffer in partitions, consume  │
│                            │ at steady rate                 │
├────────────────────────────┼────────────────────────────────┤
│ TSDB node failure          │ Messages retained until        │
│                            │ consumer catches up            │
├────────────────────────────┼────────────────────────────────┤
│ Backpressure               │ Consumer lag increases, but    │
│                            │ no data loss                   │
├────────────────────────────┼────────────────────────────────┤
│ Ordering guarantee         │ Per-partition ordering ensures │
│                            │ series data arrives in order   │
└─────────────────────────────────────────────────────────────┘

Kafka Configuration:
  - Topic: metrics-ingest
  - Partitions: 512 (partition by series_id hash)
  - Replication factor: 3
  - Retention: 24 hours (buffer for consumer failures)
  - Compression: LZ4 (fast, good ratio for numeric data)
  - Batch size: 64 KB
  - Linger: 10ms
```

### 5.3 Ingestion Workers

**Candidate:** "Workers consume from Kafka, resolve series IDs, and write to the TSDB."

```python
# Ingestion worker: consumes from Kafka, writes to TSDB
class IngestionWorker:
    def __init__(self, kafka_consumer, tsdb_client, series_cache):
        self.consumer = kafka_consumer
        self.tsdb = tsdb_client
        # LRU cache: series_id → shard assignment (avoids repeated lookups)
        self.series_cache = series_cache
        self.write_buffer = {}  # shard_id → list of (series_id, timestamp, value)
        self.buffer_flush_size = 10_000

    def run(self):
        for batch in self.consumer.poll_batches(max_records=5000, timeout_ms=100):
            for record in batch:
                metric = deserialize(record.value)
                series_id = metric['seriesId']

                shard = self.series_cache.get(series_id)
                if shard is None:
                    shard = self.tsdb.resolve_or_create_series(
                        series_id, metric['name'], metric['labels']
                    )
                    self.series_cache.put(series_id, shard)

                if shard not in self.write_buffer:
                    self.write_buffer[shard] = []
                self.write_buffer[shard].append(
                    (series_id, metric['timestamp'], metric['value'])
                )

            self.flush_if_needed()

    def flush_if_needed(self):
        for shard_id, points in self.write_buffer.items():
            if len(points) >= self.buffer_flush_size:
                self.tsdb.batch_write(shard_id, points)
                self.write_buffer[shard_id] = []
```

### 5.4 Query Engine

**Candidate:** "The query engine parses a PromQL-like expression, fans out to relevant shards, and merges results."

```
Query Execution Pipeline:
┌──────────┐    ┌──────────┐    ┌───────────┐    ┌──────────┐    ┌─────────┐
│  Parse   │───▶│  Plan    │───▶│  Fan Out  │───▶│  Execute │───▶│  Merge  │
│  PromQL  │    │  (which  │    │  (to      │    │  (per    │    │  &      │
│          │    │  shards, │    │  shards)  │    │  shard)  │    │ Return  │
│          │    │  time)   │    │           │    │          │    │         │
└──────────┘    └──────────┘    └───────────┘    └──────────┘    └─────────┘
```

```javascript
class QueryEngine {
  async execute(query) {
    const ast = this.parser.parse(query.expr);
    const timeRange = { start: query.start, end: query.end, step: query.step };
    const plan = this.planner.plan(ast, timeRange);

    // Identify which shards hold the relevant series
    const shardQueries = this.router.route(plan);

    // Fan out to shards in parallel
    const shardResults = await Promise.all(
      shardQueries.map(sq =>
        this.executeOnShard(sq.shardId, sq.seriesSelectors, timeRange)
      )
    );

    // Merge and apply top-level aggregation
    return this.merger.merge(shardResults, plan.aggregation);
  }

  async executeOnShard(shardId, selectors, timeRange) {
    const tsdb = this.getShardConnection(shardId);

    // Step 1: Use inverted index to find matching series IDs
    const seriesIds = await tsdb.lookupSeries(selectors);

    // Step 2: Scan chunks in the time range for each series
    const rawData = await tsdb.scanChunks(seriesIds, timeRange);

    // Step 3: Apply per-shard aggregation to reduce data transfer
    return this.localAggregate(rawData, timeRange.step);
  }
}
```

---

## 6. Deep Dive: Database Design

### 6.1 Time-Series Data Model

**Candidate:** "A time series is uniquely identified by a metric name plus a set of labels."

```
Time Series Identity:
  metric_name + sorted(labels) = unique time series

Examples:
  http_requests_total{service="api", method="GET", status="200"}   → Series A
  http_requests_total{service="api", method="GET", status="500"}   → Series B
  http_requests_total{service="api", method="POST", status="200"}  → Series C
  cpu_usage_percent{host="web-01", region="us-east"}               → Series D

Each series is an append-only stream of (timestamp, value) pairs:
  Series A: [(t1, 1042), (t2, 1043), (t3, 1045), (t4, 1048), ...]
  Series D: [(t1, 72.3), (t2, 74.1), (t3, 73.8), (t4, 71.2), ...]
```

### 6.2 Storage Schema

```
┌─────────────────────────────────────────────────────────────────┐
│                     LOGICAL DATA MODEL                          │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ series table                                              │  │
│  │                                                           │  │
│  │ series_id (uint64)  │ metric_name │ labels (sorted kv)   │  │
│  │ ─────────────────── │ ─────────── │ ──────────────────── │  │
│  │ 0xABCD1234          │ http_req    │ {method=GET,svc=api} │  │
│  │ 0xEFGH5678          │ cpu_usage   │ {host=web-01}        │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ samples table (append-only)                               │  │
│  │                                                           │  │
│  │ series_id (uint64)  │ timestamp (int64) │ value (float64) │  │
│  │ ─────────────────── │ ───────────────── │ ─────────────── │  │
│  │ 0xABCD1234          │ 1714000000        │ 1042.0          │  │
│  │ 0xABCD1234          │ 1714000010        │ 1043.0          │  │
│  │ 0xABCD1234          │ 1714000020        │ 1045.0          │  │
│  └───────────────────────────────────────────────────────────┘  │
│                                                                 │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │ inverted_index                                            │  │
│  │                                                           │  │
│  │ label_name │ label_value │ series_ids (posting list)      │  │
│  │ ────────── │ ─────────── │ ──────────────────────────     │  │
│  │ service    │ api         │ [0xABCD1234, 0x1111AAAA, ...]  │  │
│  │ method     │ GET         │ [0xABCD1234, 0x2222BBBB, ...]  │  │
│  │ host       │ web-01      │ [0xEFGH5678, 0x3333CCCC, ...]  │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

### 6.3 Why General-Purpose Databases Fail

**Candidate:** "Time-series data has unique access patterns that break traditional databases."

```
┌─────────────────────────┬─────────────────────┬────────────────────────────┐
│ Challenge               │ PostgreSQL/MySQL     │ Custom TSDB                │
├─────────────────────────┼─────────────────────┼────────────────────────────┤
│ Write throughput         │ ~50K inserts/sec     │ 10M+ inserts/sec           │
│                         │ (B-tree overhead)    │ (append-only, sequential)  │
├─────────────────────────┼─────────────────────┼────────────────────────────┤
│ Compression             │ ~24 bytes/point      │ ~1.4 bytes/point           │
│                         │ (row overhead)       │ (Gorilla compression)      │
├─────────────────────────┼─────────────────────┼────────────────────────────┤
│ Range scan efficiency   │ Random I/O           │ Sequential I/O             │
│                         │ (B-tree traversal)   │ (columnar chunks)          │
├─────────────────────────┼─────────────────────┼────────────────────────────┤
│ Deletion                │ Per-row delete       │ Drop entire block (O(1))   │
│                         │ (vacuum overhead)    │                            │
├─────────────────────────┼─────────────────────┼────────────────────────────┤
│ Label queries           │ Expensive JOINs      │ Inverted index (like       │
│                         │ or EAV anti-pattern  │ search engine)             │
├─────────────────────────┼─────────────────────┼────────────────────────────┤
│ Downsampling            │ Manual aggregation   │ Built-in rollup engine     │
│                         │ jobs, lock contention│                            │
└─────────────────────────┴─────────────────────┴────────────────────────────┘
```

### 6.4 SQL Comparison (For Context)

**Candidate:** "For reference, here's how the data would look in a relational model — and why it's problematic at our scale."

```sql
-- Relational model (BAD for 10M points/sec)
CREATE TABLE metrics (
    series_id   BIGINT NOT NULL,
    timestamp   BIGINT NOT NULL,
    value       DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (series_id, timestamp)
);

CREATE TABLE series (
    series_id   BIGINT PRIMARY KEY,
    metric_name VARCHAR(256) NOT NULL,
    labels      JSONB NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);

CREATE TABLE labels_index (
    label_name  VARCHAR(128) NOT NULL,
    label_value VARCHAR(256) NOT NULL,
    series_id   BIGINT NOT NULL,
    PRIMARY KEY (label_name, label_value, series_id)
);

-- Query: avg CPU for service "api" over last hour, grouped by host
-- This requires JOIN + GROUP BY + sequential scan — doesn't scale
SELECT
    s.labels->>'host' AS host,
    date_trunc('minute', to_timestamp(m.timestamp)) AS minute,
    AVG(m.value) AS avg_cpu
FROM metrics m
JOIN series s ON m.series_id = s.series_id
JOIN labels_index l ON s.series_id = l.series_id
WHERE l.label_name = 'service' AND l.label_value = 'api'
  AND m.timestamp BETWEEN EXTRACT(EPOCH FROM NOW() - INTERVAL '1 hour')
                       AND EXTRACT(EPOCH FROM NOW())
GROUP BY host, minute
ORDER BY minute;
-- At 10M points/sec, this table grows by 864B rows/day — unsustainable
```

---

## 7. Deep Dive: Time-Series Database Engine

> **This is the KEY section** — the custom TSDB engine is what makes or breaks a monitoring system.

### 7.1 Write Path

**Candidate:** "The write path is optimized for high throughput and durability. Every write goes through three stages."

```
Write Path:
                                                    
  Data Point ──▶ WAL (disk) ──▶ Head Block (memory) ──▶ Chunk File (disk)
                  │                  │                       │
                  │ Sequential       │ Fast lookups          │ Immutable,
                  │ append, crash    │ for recent            │ compressed,
                  │ recovery         │ queries               │ read-optimized
                  │                  │                       │
                  ▼                  ▼                       ▼
              Durability         Low Latency              Long-term
              Guarantee          Reads                    Storage

Timeline:
  ◄──── 2-hour head block ────►◄──── flushed chunk files ──────────────►
  [in-memory, mutable]         [on-disk, immutable, compressed]
```

```python
class TSDBWritePath:
    def __init__(self, wal_dir, chunk_dir, head_block_duration=7200):
        self.wal = WriteAheadLog(wal_dir)
        self.head = HeadBlock(duration_sec=head_block_duration)
        self.chunk_dir = chunk_dir

    def write(self, series_id, timestamp, value):
        # Step 1: Write to WAL for durability (sequential append — fast)
        self.wal.append(series_id, timestamp, value)

        # Step 2: Write to in-memory head block
        self.head.add_sample(series_id, timestamp, value)

        # Step 3: Check if head block is full (2-hour window elapsed)
        if self.head.is_full():
            self.flush_head_to_chunk()

    def flush_head_to_chunk(self):
        """Convert in-memory head block to compressed on-disk chunk."""
        chunk = ChunkBuilder()
        for series_id, samples in self.head.iter_series():
            timestamps = [s.timestamp for s in samples]
            values = [s.value for s in samples]

            # Compress using Gorilla encoding
            compressed_ts = gorilla_encode_timestamps(timestamps)
            compressed_vals = gorilla_encode_values(values)

            chunk.add_series(series_id, compressed_ts, compressed_vals)

        chunk_path = f"{self.chunk_dir}/chunk_{self.head.min_time}_{self.head.max_time}.bin"
        chunk.write_to_disk(chunk_path)

        # Truncate WAL up to the flushed point
        self.wal.truncate_before(self.head.max_time)
        self.head.reset()


class HeadBlock:
    """In-memory buffer for recent data. Optimized for writes and reads."""
    def __init__(self, duration_sec):
        self.duration = duration_sec
        self.series = {}       # series_id → list of (timestamp, value)
        self.min_time = None
        self.max_time = None

    def add_sample(self, series_id, timestamp, value):
        if self.min_time is None:
            self.min_time = timestamp
        self.max_time = max(self.max_time or 0, timestamp)

        if series_id not in self.series:
            self.series[series_id] = []
        self.series[series_id].append((timestamp, value))

    def is_full(self):
        if self.min_time is None:
            return False
        return (self.max_time - self.min_time) >= self.duration
```

### 7.2 Compression Techniques

**Candidate:** "Compression is critical. Without it, we'd need ~20 TB/day. With Gorilla compression, we achieve roughly 1.37 bytes per data point — a 12× reduction."

#### Delta-of-Delta Encoding for Timestamps

```
Intuition:
  Metrics arrive at regular intervals (e.g., every 10 seconds).
  Instead of storing absolute timestamps, store the difference-of-differences.

Example:
  Raw timestamps:   [1000, 1010, 1020, 1030, 1040, 1050]
  Deltas:           [–,    10,   10,   10,   10,   10  ]
  Delta-of-deltas:  [–,    –,    0,    0,    0,    0   ]

  For regular intervals, delta-of-delta = 0 → encode as 1 bit!
  
Encoding scheme (from Facebook Gorilla paper):
  ┌───────────────────────────────────────────────────────┐
  │ Delta-of-Delta Value    │ Encoding                    │
  ├─────────────────────────┼─────────────────────────────┤
  │ 0                       │ '0' (1 bit)                 │
  │ [-63, 64]               │ '10' + 7 bits (9 bits)      │
  │ [-255, 256]             │ '110' + 9 bits (12 bits)    │
  │ [-2047, 2048]           │ '1110' + 12 bits (16 bits)  │
  │ everything else         │ '1111' + 32 bits (36 bits)  │
  └───────────────────────────────────────────────────────┘
  
  For regular 10-second intervals: each timestamp costs just 1 bit!
```

```javascript
// Delta-of-delta encoding for timestamps (Gorilla compression)
class TimestampEncoder {
  constructor() {
    this.bitstream = new BitStream();
    this.prevTimestamp = 0;
    this.prevDelta = 0;
    this.count = 0;
  }

  encode(timestamp) {
    if (this.count === 0) {
      // First timestamp: store full value (header)
      this.bitstream.writeBits(timestamp, 64);
    } else if (this.count === 1) {
      // Second timestamp: store delta
      const delta = timestamp - this.prevTimestamp;
      this.bitstream.writeBits(delta, 14); // 14 bits for initial delta
      this.prevDelta = delta;
    } else {
      // Subsequent: store delta-of-delta
      const delta = timestamp - this.prevTimestamp;
      const deltaOfDelta = delta - this.prevDelta;

      if (deltaOfDelta === 0) {
        this.bitstream.writeBit(0);                       // 1 bit
      } else if (deltaOfDelta >= -63 && deltaOfDelta <= 64) {
        this.bitstream.writeBits(0b10, 2);
        this.bitstream.writeBitsSigned(deltaOfDelta, 7);  // 9 bits total
      } else if (deltaOfDelta >= -255 && deltaOfDelta <= 256) {
        this.bitstream.writeBits(0b110, 3);
        this.bitstream.writeBitsSigned(deltaOfDelta, 9);  // 12 bits total
      } else if (deltaOfDelta >= -2047 && deltaOfDelta <= 2048) {
        this.bitstream.writeBits(0b1110, 4);
        this.bitstream.writeBitsSigned(deltaOfDelta, 12); // 16 bits total
      } else {
        this.bitstream.writeBits(0b1111, 4);
        this.bitstream.writeBitsSigned(deltaOfDelta, 32); // 36 bits total
      }
      this.prevDelta = delta;
    }
    this.prevTimestamp = timestamp;
    this.count++;
  }
}
```

#### XOR Encoding for Float Values

```
Intuition:
  Consecutive metric values are often similar (cpu goes 72.3 → 72.5 → 72.1).
  XOR of similar IEEE 754 floats has long runs of zeros.

Example:
  value1 = 72.3  →  binary: 0 10000000101 0010000100110011001100110011001100110011001100110100
  value2 = 72.5  →  binary: 0 10000000101 0010001000000000000000000000000000000000000000000000
  XOR            →          0 00000000000 0000001100110011001100110011001100110011001100110100
                            ^^^^^^^^^^^^^^^^ leading zeros     ^^^^^^^^^^^^^^^^ trailing zeros

  Store: number of leading zeros + length of meaningful bits + meaningful bits
  For similar consecutive values: typically 20-30 bits instead of 64 bits
```

```python
import struct

def float_to_bits(f):
    """Convert float64 to uint64 bit representation."""
    return struct.unpack('Q', struct.pack('d', f))[0]

def bits_to_float(b):
    """Convert uint64 bit representation to float64."""
    return struct.unpack('d', struct.pack('Q', b))[0]

class ValueEncoder:
    def __init__(self):
        self.bitstream = BitStream()
        self.prev_bits = 0
        self.prev_leading = 0
        self.prev_trailing = 0
        self.count = 0

    def encode(self, value):
        bits = float_to_bits(value)

        if self.count == 0:
            self.bitstream.write_bits(bits, 64)
        else:
            xor = bits ^ self.prev_bits

            if xor == 0:
                # Values are identical — encode as single 0 bit
                self.bitstream.write_bit(0)
            else:
                self.bitstream.write_bit(1)
                leading = count_leading_zeros(xor)
                trailing = count_trailing_zeros(xor)
                meaningful_bits = 64 - leading - trailing

                if (leading >= self.prev_leading and
                    trailing >= self.prev_trailing):
                    # Reuse previous leading/trailing counts
                    self.bitstream.write_bit(0)
                    meaningful = 64 - self.prev_leading - self.prev_trailing
                    self.bitstream.write_bits(
                        xor >> self.prev_trailing, meaningful
                    )
                else:
                    # Write new leading/trailing metadata
                    self.bitstream.write_bit(1)
                    self.bitstream.write_bits(leading, 5)
                    self.bitstream.write_bits(meaningful_bits, 6)
                    self.bitstream.write_bits(xor >> trailing, meaningful_bits)

                    self.prev_leading = leading
                    self.prev_trailing = trailing

        self.prev_bits = bits
        self.count += 1
```

#### Compression Results

```
Compression Effectiveness:
┌──────────────────────┬──────────┬────────────────┬──────────┐
│ Component            │ Raw Size │ Compressed     │ Ratio    │
├──────────────────────┼──────────┼────────────────┼──────────┤
│ Timestamp (8 bytes)  │ 64 bits  │ ~1 bit avg     │ ~64×     │
│ Value (8 bytes)      │ 64 bits  │ ~10 bits avg   │ ~6.4×    │
│ Total per point      │ 16 bytes │ ~1.37 bytes    │ ~11.7×   │
└──────────────────────┴──────────┴────────────────┴──────────┘

Storage savings:
  Before: 10M pts/sec × 16 bytes × 86,400 sec = 13.8 TB/day
  After:  10M pts/sec × 1.37 bytes × 86,400 sec = 1.18 TB/day
  Savings: ~12 TB/day
```

### 7.3 Read Path

**Candidate:** "The read path resolves which series match the query, then scans the appropriate chunks."

```
Read Path:
                                                    
  Query: avg(cpu{service="api"}) over [now-1h, now]
    │
    ▼
  ┌──────────────────────────────────────────────────────┐
  │ Step 1: INVERTED INDEX LOOKUP                        │
  │                                                      │
  │  "service" = "api"  →  posting list: [S1, S4, S7]   │
  │  "__name__" = "cpu"  →  posting list: [S1, S2, S7]   │
  │                                                      │
  │  Intersection: [S1, S7]  (2 matching series)         │
  └──────────────────────┬───────────────────────────────┘
                         │
                         ▼
  ┌──────────────────────────────────────────────────────┐
  │ Step 2: IDENTIFY RELEVANT CHUNKS                     │
  │                                                      │
  │  Time range: [now-1h, now]                           │
  │  For S1: head block (in-memory) has last 2 hours     │
  │  For S7: head block (in-memory) has last 2 hours     │
  │                                                      │
  │  → Only need to read from head block (fast!)         │
  │  → For 30-day queries, scan multiple chunk files     │
  └──────────────────────┬───────────────────────────────┘
                         │
                         ▼
  ┌──────────────────────────────────────────────────────┐
  │ Step 3: DECOMPRESS & AGGREGATE                       │
  │                                                      │
  │  S1: [(t1,72.3), (t2,74.1), (t3,73.8), ...]        │
  │  S7: [(t1,68.9), (t2,70.2), (t3,69.5), ...]        │
  │                                                      │
  │  avg() at each step:                                 │
  │  t1: (72.3 + 68.9) / 2 = 70.6                      │
  │  t2: (74.1 + 70.2) / 2 = 72.15                     │
  │  t3: (73.8 + 69.5) / 2 = 71.65                     │
  └──────────────────────────────────────────────────────┘
```

### 7.4 Inverted Index

**Candidate:** "The inverted index is fundamental — it's how we turn label queries into series lookups in O(log n) time."

```
Inverted Index Structure (like a search engine):

  Label name → Label value → Sorted posting list of series IDs

  ┌──────────────────────────────────────────────────────────────┐
  │  __name__                                                    │
  │    ├── "cpu_usage"     → [S001, S002, S003, S004, ...]      │
  │    ├── "http_requests" → [S100, S101, S102, S103, ...]      │
  │    └── "memory_used"   → [S200, S201, S202, ...]            │
  │                                                              │
  │  service                                                     │
  │    ├── "api"           → [S001, S100, S200, S301, ...]      │
  │    ├── "auth"          → [S002, S101, S201, S302, ...]      │
  │    └── "payments"      → [S003, S102, S202, ...]            │
  │                                                              │
  │  region                                                      │
  │    ├── "us-east-1"     → [S001, S002, S100, S101, ...]      │
  │    └── "eu-west-1"     → [S003, S004, S102, S103, ...]      │
  └──────────────────────────────────────────────────────────────┘

Query: cpu_usage{service="api", region="us-east-1"}

  Step 1: __name__ = "cpu_usage"   → [S001, S002, S003, S004]
  Step 2: service  = "api"          → [S001, S100, S200, S301]
  Step 3: region   = "us-east-1"    → [S001, S002, S100, S101]

  Intersect all three posting lists → [S001]
  (Intersection uses merge-join on sorted lists — O(n) per list)
```

```python
class InvertedIndex:
    def __init__(self):
        # label_name → label_value → sorted list of series_ids
        self.postings = {}

    def add_series(self, series_id, labels):
        for name, value in labels.items():
            if name not in self.postings:
                self.postings[name] = {}
            if value not in self.postings[name]:
                self.postings[name][value] = SortedList()
            self.postings[name][value].add(series_id)

    def lookup(self, matchers):
        """Find series matching all label matchers (AND semantics)."""
        posting_lists = []
        for matcher in matchers:
            if matcher.name in self.postings:
                values = self.postings[matcher.name]
                if matcher.type == 'equal':
                    if matcher.value in values:
                        posting_lists.append(values[matcher.value])
                elif matcher.type == 'regex':
                    merged = SortedList()
                    for val, sids in values.items():
                        if matcher.regex.match(val):
                            merged = merged.union(sids)
                    posting_lists.append(merged)

        if not posting_lists:
            return []
        return self.intersect_sorted(posting_lists)

    def intersect_sorted(self, lists):
        """Intersect multiple sorted posting lists efficiently."""
        # Start with smallest list for efficiency
        lists.sort(key=len)
        result = lists[0]
        for other in lists[1:]:
            result = self.sorted_intersect_two(result, other)
        return result
```

### 7.5 Downsampling

**Candidate:** "We can't keep raw 10-second data forever. Downsampling reduces storage while preserving queryable aggregates."

```
Downsampling Pipeline:
                                                    
  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐    ┌──────────────┐
  │  Raw Data    │    │  1-Minute    │    │  5-Minute    │    │  1-Hour      │
  │  (10-sec)    │───▶│  Rollup     │───▶│  Rollup     │───▶│  Rollup      │
  │              │    │              │    │              │    │              │
  │  Keep 15     │    │  Keep 90     │    │  Keep 180    │    │  Keep 365    │
  │  days        │    │  days        │    │  days        │    │  days        │
  └──────────────┘    └──────────────┘    └──────────────┘    └──────────────┘

What gets stored in each rollup window:
  ┌─────────────────────────────────────────────────────┐
  │ For each (series_id, window_start):                 │
  │                                                     │
  │   min:   minimum value in window                    │
  │   max:   maximum value in window                    │
  │   sum:   sum of all values in window                │
  │   count: number of data points in window            │
  │   avg:   sum / count (derived, not stored)          │
  │                                                     │
  │ This allows accurate re-aggregation at query time:  │
  │   avg = sum(sums) / sum(counts) across windows      │
  │   max = max(maxes) across windows                   │
  └─────────────────────────────────────────────────────┘
```

```python
class DownsamplingService:
    """Periodically rolls up fine-grained data into coarser resolutions."""

    ROLLUP_TIERS = [
        {'source': '10s',  'target': '1m',  'window': 60,    'retention_days': 90},
        {'source': '1m',   'target': '5m',  'window': 300,   'retention_days': 180},
        {'source': '5m',   'target': '1h',  'window': 3600,  'retention_days': 365},
    ]

    def downsample_tier(self, tsdb, tier):
        """Process one rollup tier for all series."""
        window = tier['window']
        pending_blocks = tsdb.get_blocks_needing_rollup(tier['source'], tier['target'])

        for block in pending_blocks:
            for series_id in block.series_ids():
                samples = block.read_series(series_id)
                rollups = self.compute_rollups(samples, window)

                for rollup in rollups:
                    tsdb.write_rollup(
                        tier['target'], series_id,
                        rollup['window_start'],
                        rollup['min'], rollup['max'],
                        rollup['sum'], rollup['count']
                    )

    def compute_rollups(self, samples, window_sec):
        rollups = []
        current_window = None
        buf = []

        for ts, val in samples:
            w_start = (ts // window_sec) * window_sec
            if current_window is None:
                current_window = w_start

            if w_start != current_window:
                rollups.append(self.aggregate_window(current_window, buf))
                buf = []
                current_window = w_start

            buf.append(val)

        if buf:
            rollups.append(self.aggregate_window(current_window, buf))
        return rollups

    def aggregate_window(self, window_start, values):
        return {
            'window_start': window_start,
            'min': min(values),
            'max': max(values),
            'sum': sum(values),
            'count': len(values),
        }
```

### 7.6 Compaction

```
Compaction: merge small chunk files into larger, more efficient blocks.

Before compaction:
  ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
  │ 2-hr   │ │ 2-hr   │ │ 2-hr   │ │ 2-hr   │ │ 2-hr   │ │ 2-hr   │
  │ chunk  │ │ chunk  │ │ chunk  │ │ chunk  │ │ chunk  │ │ chunk  │
  └────────┘ └────────┘ └────────┘ └────────┘ └────────┘ └────────┘

After compaction (level 1):
  ┌─────────────────────────────┐ ┌─────────────────────────────┐
  │       6-hour block          │ │       6-hour block          │
  │  (merged + re-compressed)   │ │  (merged + re-compressed)   │
  └─────────────────────────────┘ └─────────────────────────────┘

After compaction (level 2):
  ┌─────────────────────────────────────────────────────────────┐
  │                    12-hour block                             │
  └─────────────────────────────────────────────────────────────┘

Benefits:
  ✓ Fewer files to scan during queries
  ✓ Better compression (more data = more patterns)
  ✓ Tombstoned (deleted) series are physically removed
  ✓ Index is rebuilt more compactly
```

### 7.7 High Cardinality Handling

**Candidate:** "High cardinality is the #1 operational problem in monitoring systems. It happens when a label has too many unique values."

```
The Problem:
  http_requests_total{user_id="abc123"}     ← 100M users = 100M series!
  http_requests_total{request_id="..."}     ← infinite cardinality
  http_requests_total{pod_name="pod-xxx"}   ← Kubernetes churn

Why it's dangerous:
  - 10M → 100M series = 10× memory for inverted index
  - Write amplification (more series to track)
  - Query performance degrades (larger posting lists)
  - Compaction takes longer
  - A single bad metric can take down the entire cluster

Real-world horror stories:
  - A team added user_id as a label → created 50M new series overnight
  - Kubernetes pod name churn → millions of "zombie" series that never get new data
```

```
Solutions:
┌────────────────────────┬───────────────────────────────────────────────────┐
│ Strategy               │ Implementation                                    │
├────────────────────────┼───────────────────────────────────────────────────┤
│ Cardinality Limits     │ Reject metrics that would create new series      │
│                        │ beyond a per-metric threshold (e.g., 10K series) │
├────────────────────────┼───────────────────────────────────────────────────┤
│ Pre-Aggregation        │ Aggregate at the agent before sending:           │
│                        │ count by (service, status) instead of per-user   │
├────────────────────────┼───────────────────────────────────────────────────┤
│ Bloom Filters          │ Probabilistic check: "does this series exist?"   │
│                        │ before hitting the full inverted index            │
├────────────────────────┼───────────────────────────────────────────────────┤
│ Stale Series Cleanup   │ Mark series with no data for 15 minutes as stale │
│                        │ Remove from head block to free memory             │
├────────────────────────┼───────────────────────────────────────────────────┤
│ Label Value Allow-List │ Only accept known values for specific labels     │
│                        │ (e.g., region must be in known set)              │
├────────────────────────┼───────────────────────────────────────────────────┤
│ Tenant Quotas          │ Each tenant/team has a max active series count   │
│                        │ Reject with 429 when exceeded                    │
└────────────────────────┴───────────────────────────────────────────────────┘
```

```javascript
class CardinalityEnforcer {
  constructor(maxSeriesPerMetric = 10000, maxTotalSeries = 10_000_000) {
    this.maxPerMetric = maxSeriesPerMetric;
    this.maxTotal = maxTotalSeries;
    this.metricSeriesCount = new Map();  // metric_name → count
    this.totalSeries = 0;
    this.bloomFilter = new BloomFilter(10_000_000, 0.01);
  }

  canCreateSeries(metricName, seriesId) {
    // Fast path: series already exists (bloom filter check)
    if (this.bloomFilter.mightContain(seriesId)) {
      return true;  // Existing series — always allow writes
    }

    // New series — enforce limits
    const metricCount = this.metricSeriesCount.get(metricName) || 0;
    if (metricCount >= this.maxPerMetric) {
      this.emitAlert('cardinality_limit', metricName, metricCount);
      return false;  // Per-metric limit exceeded
    }
    if (this.totalSeries >= this.maxTotal) {
      this.emitAlert('total_series_limit', metricName, this.totalSeries);
      return false;  // Global limit exceeded
    }

    // Allow and track
    this.bloomFilter.add(seriesId);
    this.metricSeriesCount.set(metricName, metricCount + 1);
    this.totalSeries++;
    return true;
  }
}
```

### 7.8 Alerting Engine

**Candidate:** "The alerting engine evaluates rules against the TSDB at regular intervals and manages alert lifecycle."

```
Alert State Machine:
                                                    
         ┌──────────────────────────────────────────────────┐
         │                                                  │
         ▼                                                  │
     ┌────────┐    condition     ┌─────────┐    "for"      │
     │   OK   │────true for ───▶│ PENDING  │───duration───▶│
     │        │    1 eval cycle  │          │   elapsed     │
     └────┬───┘                  └────┬─────┘               │
          │                           │                     │
          │                    condition                  ┌────────────┐
          │                    becomes false              │  FIRING    │
          │                           │                   │            │
          │                           ▼                   │ → Send     │
          │                      ┌────────┐               │   notif.   │
          │                      │   OK   │               │ → Repeat   │
          │                      └────────┘               │   interval │
          │                                               └──────┬─────┘
          │                                                      │
          │              condition becomes false                  │
          │◀─────────────────────────────────────────────────────┘
          │                                                      
          ▼                                                      
     ┌──────────┐                                                
     │ RESOLVED │── sends resolution notification                
     └──────────┘                                                
```

```python
import time
from enum import Enum

class AlertState(Enum):
    OK = "ok"
    PENDING = "pending"
    FIRING = "firing"
    RESOLVED = "resolved"

class AlertRule:
    def __init__(self, name, expr, for_duration, severity, notifications):
        self.name = name
        self.expr = expr                      # PromQL expression
        self.for_duration = for_duration      # Seconds to wait before firing
        self.severity = severity
        self.notifications = notifications
        self.state = AlertState.OK
        self.pending_since = None
        self.firing_since = None
        self.last_notification = None
        self.repeat_interval = 3600           # Re-notify every hour while firing

class AlertingEngine:
    def __init__(self, query_engine, notifier, eval_interval=15):
        self.query_engine = query_engine
        self.notifier = notifier
        self.eval_interval = eval_interval    # Seconds between evaluations
        self.rules = []

    def run_evaluation_loop(self):
        """Main loop: evaluate all rules every eval_interval seconds."""
        while True:
            start = time.time()
            for rule in self.rules:
                try:
                    self.evaluate_rule(rule)
                except Exception as e:
                    self.record_eval_failure(rule, e)
            elapsed = time.time() - start
            sleep_time = max(0, self.eval_interval - elapsed)
            time.sleep(sleep_time)

    def evaluate_rule(self, rule):
        result = self.query_engine.instant_query(rule.expr)
        condition_met = len(result) > 0 and any(r['value'] > 0 for r in result)

        if condition_met:
            if rule.state == AlertState.OK:
                rule.state = AlertState.PENDING
                rule.pending_since = time.time()

            elif rule.state == AlertState.PENDING:
                if (time.time() - rule.pending_since) >= rule.for_duration:
                    rule.state = AlertState.FIRING
                    rule.firing_since = time.time()
                    self.send_notification(rule, "FIRING")

            elif rule.state == AlertState.FIRING:
                if self.should_repeat_notification(rule):
                    self.send_notification(rule, "STILL FIRING")

        else:
            if rule.state in (AlertState.PENDING, AlertState.FIRING):
                was_firing = rule.state == AlertState.FIRING
                rule.state = AlertState.OK
                rule.pending_since = None
                if was_firing:
                    self.send_notification(rule, "RESOLVED")
                    rule.firing_since = None

    def should_repeat_notification(self, rule):
        if rule.last_notification is None:
            return True
        return (time.time() - rule.last_notification) >= rule.repeat_interval

    def send_notification(self, rule, status):
        for target in rule.notifications:
            self.notifier.send(
                channel=target['channel'],
                destination=target.get('routing_key') or target.get('webhook'),
                alert_name=rule.name,
                severity=rule.severity,
                status=status,
                fired_at=rule.firing_since,
            )
        rule.last_notification = time.time()
```

---

## 8. Scaling Strategies

### 8.1 Horizontal Sharding

**Candidate:** "We shard the TSDB by consistent hashing on the series ID."

```
Sharding Strategy:
                                                    
  series_id = hash(metric_name + sorted_labels)
  shard = consistent_hash(series_id) % num_shards

  ┌────────────────────────────────────────────────────────────┐
  │ series_id range 0x0000–0x3FFF  →  TSDB Shard 0 (primary)  │
  │                                   TSDB Shard 0 (replica)   │
  ├────────────────────────────────────────────────────────────┤
  │ series_id range 0x4000–0x7FFF  →  TSDB Shard 1 (primary)  │
  │                                   TSDB Shard 1 (replica)   │
  ├────────────────────────────────────────────────────────────┤
  │ series_id range 0x8000–0xBFFF  →  TSDB Shard 2 (primary)  │
  │                                   TSDB Shard 2 (replica)   │
  ├────────────────────────────────────────────────────────────┤
  │ series_id range 0xC000–0xFFFF  →  TSDB Shard 3 (primary)  │
  │                                   TSDB Shard 3 (replica)   │
  └────────────────────────────────────────────────────────────┘

Key property: all data points for a given time series land on the same shard.
This means per-shard queries avoid cross-shard joins for single-series lookups.
```

### 8.2 Separated Read and Write Paths

```
┌──────────────────────────────────────────────────────────────┐
│                    WRITE PATH                                │
│                                                              │
│  Gateway → Kafka → Workers → TSDB Primary Nodes             │
│                                                              │
│  Optimized for:                                              │
│  ✓ Sequential writes (append-only)                           │
│  ✓ High throughput (batch + async)                           │
│  ✓ WAL for durability                                        │
└──────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────┐
│                    READ PATH                                 │
│                                                              │
│  API → Query Engine → TSDB Read Replicas + Query Cache       │
│                                                              │
│  Optimized for:                                              │
│  ✓ Parallel shard scanning                                   │
│  ✓ Result caching (same dashboard query repeated)            │
│  ✓ Read replicas for isolation from writes                   │
│  ✓ Pre-computed rollups for historical queries               │
└──────────────────────────────────────────────────────────────┘
```

### 8.3 Query Caching

```
Cache layers:
  L1: Query result cache (Redis) — exact query + time range → result
      TTL: 10-60 seconds (short, because data changes constantly)
      Hit rate: 60-80% (dashboards auto-refresh same queries)

  L2: Chunk cache (local SSD) — frequently accessed chunks in fast storage
      Avoids reading from slower networked storage

  L3: Metadata cache (in-memory) — series ID → labels mapping, inverted index
      Always warm, critical for query planning
```

### 8.4 Federation (Multi-Cluster)

```
Multi-Region Architecture:
                                                    
  ┌───────────────┐         ┌───────────────┐         ┌───────────────┐
  │  US-East      │         │  EU-West      │         │  AP-South     │
  │  Cluster      │         │  Cluster      │         │  Cluster      │
  │               │         │               │         │               │
  │  Local TSDB   │         │  Local TSDB   │         │  Local TSDB   │
  │  Local Alert  │         │  Local Alert  │         │  Local Alert  │
  └───────┬───────┘         └───────┬───────┘         └───────┬───────┘
          │                         │                         │
          │    Remote Write/Read    │                         │
          └────────────┬────────────┘                         │
                       │                                      │
               ┌───────▼───────┐                              │
               │  Global       │◀─────────────────────────────┘
               │  Query Layer  │
               │               │
               │  Aggregates   │
               │  across       │
               │  clusters     │
               └───────────────┘

  - Each region is self-contained for writes and local queries
  - Global query layer fans out to all regions for cross-region dashboards
  - Remote write: replicate critical metrics to a central long-term store
```

---

## 9. Failure Scenarios & Mitigation

### 9.1 Failure Analysis Table

```
┌──────────────────────────┬──────────────────────────┬───────────────────────────────┐
│ Failure                  │ Impact                   │ Mitigation                    │
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Ingestion gateway crash  │ Metrics dropped during   │ Multiple gateway instances    │
│                          │ failover                 │ behind load balancer; agents  │
│                          │                          │ retry with exponential backoff│
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Kafka broker failure     │ Partition unavailable     │ Replication factor=3; auto   │
│                          │ temporarily              │ leader election               │
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Ingestion worker crash   │ Consumer lag increases    │ Consumer group rebalance;    │
│                          │                          │ auto-scaling based on lag     │
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ TSDB node failure        │ Shard unavailable for    │ Replica promotion; write to  │
│                          │ writes and reads         │ WAL on replica during failover│
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Alerting engine crash    │ Missed alerts!           │ Redundant evaluators (active/ │
│                          │ CRITICAL                 │ active with dedup); watchdog  │
│                          │                          │ alert: "I'm alive" heartbeat │
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Query of death           │ Single query consumes    │ Per-query timeout (30s);     │
│                          │ all resources            │ circuit breaker; limit series │
│                          │                          │ per query (max 10K series)   │
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Clock skew               │ Timestamps in the future │ Accept within ±5 min window; │
│                          │ or far past              │ NTP enforcement on agents;   │
│                          │                          │ server-side timestamp option  │
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Disk full on TSDB node   │ Writes fail, data loss   │ Monitoring on disk usage;    │
│                          │                          │ auto-delete oldest blocks;   │
│                          │                          │ WAL on separate volume       │
├──────────────────────────┼──────────────────────────┼───────────────────────────────┤
│ Cardinality explosion    │ OOM on TSDB nodes        │ Cardinality enforcer; tenant │
│                          │                          │ quotas; circuit breaker on   │
│                          │                          │ series creation              │
└──────────────────────────┴──────────────────────────┴───────────────────────────────┘
```

### 9.2 Handling Ingestion Pipeline Backup

```
Normal Operation:
  Gateway → Kafka → Workers → TSDB
  Consumer lag: ~0 (real-time)

During TSDB slowdown:
  Gateway → Kafka (messages accumulate) → Workers (slower consumption)
  Consumer lag: growing (minutes of data buffered)

Recovery strategy:
  1. Kafka retains data for 24 hours (configured retention)
  2. Workers auto-scale based on consumer lag metric
  3. If lag exceeds 1 hour:
     - Alert ops team
     - Consider dropping lower-priority metrics (debug-level)
     - Enable batch-mode ingestion (larger writes, less frequent)
  4. Once TSDB recovers, workers catch up at accelerated rate
```

### 9.3 Alerting Engine Redundancy

```
┌─────────────────────────────────────────────────────────────────┐
│                 ALERTING HIGH AVAILABILITY                       │
│                                                                 │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │ Evaluator A  │    │ Evaluator B  │    │ Evaluator C  │      │
│  │ (active)     │    │ (active)     │    │ (active)     │      │
│  │              │    │              │    │              │      │
│  │ Rules 0-999  │    │ Rules 0-999  │    │ Rules 0-999  │      │
│  └──────┬───────┘    └──────┬───────┘    └──────┬───────┘      │
│         │                   │                   │              │
│         └───────────┬───────┘                   │              │
│                     ▼                           │              │
│         ┌──────────────────────┐                │              │
│         │  Deduplication Layer │◀───────────────┘              │
│         │  (only send 1 notif │                               │
│         │   per alert)         │                               │
│         └──────────┬───────────┘                               │
│                    ▼                                           │
│         ┌──────────────────────┐                               │
│         │  Notification Router │                               │
│         └──────────────────────┘                               │
│                                                                 │
│  Watchdog alert:                                                │
│  A special "DeadMansSwitch" alert is ALWAYS firing.            │
│  If it stops firing → alerting system itself is down.          │
└─────────────────────────────────────────────────────────────────┘
```

---

## 10. Monitoring & Observability

### 10.1 Meta-Monitoring (Monitor the Monitor)

**Candidate:** "The most ironic failure is when your monitoring system goes down and you don't know about it. We need meta-monitoring."

```
Critical self-monitoring metrics:
┌─────────────────────────────────┬───────────────┬────────────────────────────┐
│ Metric                          │ Alert Threshold│ Why It Matters             │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ ingestion_rate_points_per_sec   │ < 5M (50%     │ Pipeline may be broken     │
│                                 │ of normal)    │                            │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ kafka_consumer_lag_seconds      │ > 300 sec     │ Workers can't keep up      │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ query_latency_p99_seconds       │ > 5 sec       │ Dashboards degraded        │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ alert_eval_duration_seconds     │ > eval_interval│ Alert evaluations falling  │
│                                 │               │ behind                     │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ active_series_count             │ > 12M (120%   │ Cardinality growing        │
│                                 │ of expected)  │ unexpectedly               │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ tsdb_disk_utilization_percent   │ > 80%         │ Storage filling up         │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ wal_corruption_events           │ > 0           │ Data integrity at risk     │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ compaction_duration_seconds     │ > 1800 sec    │ Compaction falling behind  │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ dropped_samples_total           │ > 0           │ Data being lost            │
├─────────────────────────────────┼───────────────┼────────────────────────────┤
│ dead_mans_switch                │ absent        │ Alerting system is down    │
└─────────────────────────────────┴───────────────┴────────────────────────────┘
```

### 10.2 Self-Monitoring Architecture

```
Strategy: Use a SEPARATE, SIMPLE monitoring system to monitor the primary one.

  ┌───────────────────────────────────────────────────────────────┐
  │                  PRIMARY MONITORING SYSTEM                    │
  │                  (the system we're designing)                 │
  │                                                               │
  │   Emits its own internal metrics ──────┐                     │
  └────────────────────────────────────────┼─────────────────────┘
                                           │
                                           ▼
  ┌───────────────────────────────────────────────────────────────┐
  │            SECONDARY META-MONITOR (simple, minimal)           │
  │                                                               │
  │   - Lightweight Prometheus instance                           │
  │   - Scrapes health endpoints of primary system                │
  │   - 5-10 critical alert rules only                            │
  │   - Routes to PagerDuty via independent path                  │
  │   - Runs on separate infrastructure                           │
  └───────────────────────────────────────────────────────────────┘

  Key principle: the meta-monitor must be SIMPLER and have
  FEWER dependencies than the primary system.
```

### 10.3 Operational Dashboard

```
Self-Monitoring Dashboard Panels:
┌────────────────────────────────────────────────────────────────────────┐
│                    SYSTEM HEALTH DASHBOARD                             │
├───────────────────────────────┬────────────────────────────────────────┤
│  Ingestion Rate (pts/sec)     │  Active Series Count                  │
│  ┌──────────────────────┐     │  ┌──────────────────────┐             │
│  │  ▁▁▂▃▅▇██▇▅▃▂▁▁▁▂▃  │     │  │  ────────────────── │ 10.2M      │
│  │  10M avg             │     │  │  (stable)           │             │
│  └──────────────────────┘     │  └──────────────────────┘             │
├───────────────────────────────┼────────────────────────────────────────┤
│  Query Latency P99            │  Kafka Consumer Lag                   │
│  ┌──────────────────────┐     │  ┌──────────────────────┐             │
│  │  ▁▁▁▁▁▁▂▅▂▁▁▁▁▁▁▁  │     │  │  ▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁ │ 0 sec      │
│  │  0.8s avg            │     │  │  (healthy)          │             │
│  └──────────────────────┘     │  └──────────────────────┘             │
├───────────────────────────────┼────────────────────────────────────────┤
│  Alert Eval Duration          │  Disk Utilization                     │
│  ┌──────────────────────┐     │  ┌──────────────────────┐             │
│  │  ▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁▁  │     │  │  ████████████░░░░░░ │ 62%        │
│  │  2.1s avg            │     │  │  25 TB / 40 TB      │             │
│  └──────────────────────┘     │  └──────────────────────┘             │
└───────────────────────────────┴────────────────────────────────────────┘
```

---

## 11. Advanced Features

### 11.1 Service Level Objectives (SLO) Tracking

```
SLO Definition:
  "99.9% of API requests should complete in < 500ms over a 30-day window"

Implementation:
  1. Define SLI (indicator): http_request_duration_seconds{service="api"}
  2. Define good events: requests where duration < 0.5
  3. Track error budget:
     - Total requests in 30 days: ~100M
     - Allowed failures: 100M × 0.001 = 100K
     - Remaining budget: 100K - actual_failures

  Recording rule (pre-computed):
    slo:api_latency:good_rate = 
      sum(rate(http_requests_total{status!~"5.."}[5m])) /
      sum(rate(http_requests_total[5m]))

  Alert when error budget burn rate is too high:
    alert: SLOBurnRateHigh
    expr: slo:api_latency:error_rate > 14.4 * 0.001  (14.4× burn rate)
    for: 1h
```

### 11.2 Recording Rules (Pre-Compute Expensive Queries)

```python
# Recording rules pre-compute expensive queries on a schedule
# Results are stored as new time series for fast dashboard loading

RECORDING_RULES = [
    {
        'name': 'job:http_requests:rate5m',
        'expr': 'sum(rate(http_requests_total[5m])) by (service)',
        'interval': '15s',
    },
    {
        'name': 'job:http_latency:p99_5m',
        'expr': 'histogram_quantile(0.99, sum(rate(http_request_duration_bucket[5m])) by (le, service))',
        'interval': '15s',
    },
    {
        'name': 'job:error_rate:ratio_5m',
        'expr': 'sum(rate(http_requests_total{status=~"5.."}[5m])) by (service) / sum(rate(http_requests_total[5m])) by (service)',
        'interval': '15s',
    },
]

class RecordingRuleEngine:
    def __init__(self, query_engine, tsdb_writer):
        self.query = query_engine
        self.writer = tsdb_writer

    def evaluate_rule(self, rule):
        result = self.query.instant_query(rule['expr'])
        timestamp = int(time.time())
        for series in result:
            self.writer.write(
                metric_name=rule['name'],
                labels=series['labels'],
                timestamp=timestamp,
                value=series['value']
            )
```

### 11.3 Distributed Tracing Correlation

```
Linking metrics to traces:
  When a metric spike occurs, automatically find related traces.

  1. Metric: http_latency_p99{service="checkout"} spikes to 2s
  2. Query trace store: find traces where service=checkout AND duration > 1.5s
     in the same time window
  3. Display correlated traces alongside the metric graph
  
  Implementation: exemplars — sample trace IDs attached to metric data points
  
  Data model extension:
    (series_id, timestamp, value, exemplar_trace_id)
    
  Not every point has an exemplar — typically 1 in 100 to save space.
```

### 11.4 Anomaly Detection

```
Approaches:
┌─────────────────────┬────────────────────────────────────────────────────┐
│ Method              │ Description                                        │
├─────────────────────┼────────────────────────────────────────────────────┤
│ Static Threshold    │ alert if value > X                                 │
│                     │ Simple but requires manual tuning                  │
├─────────────────────┼────────────────────────────────────────────────────┤
│ Z-Score             │ alert if value > mean + 3σ                         │
│                     │ Works for normally distributed metrics              │
├─────────────────────┼────────────────────────────────────────────────────┤
│ Seasonal Decompose  │ Learn daily/weekly patterns                        │
│ (STL)               │ Alert on deviation from expected seasonal value    │
├─────────────────────┼────────────────────────────────────────────────────┤
│ EWMA                │ Exponentially weighted moving average              │
│                     │ Good for detecting gradual drift                   │
├─────────────────────┼────────────────────────────────────────────────────┤
│ Isolation Forest    │ ML-based multivariate anomaly detection            │
│                     │ Detects anomalies across multiple metrics together │
└─────────────────────┴────────────────────────────────────────────────────┘
```

```python
class AnomalyDetector:
    """Simple Z-score based anomaly detection with seasonal baseline."""

    def __init__(self, query_engine, lookback_weeks=4):
        self.query = query_engine
        self.lookback_weeks = lookback_weeks

    def is_anomalous(self, metric_expr, current_value, timestamp):
        # Get historical values at the same time-of-day, same day-of-week
        historical = self.get_seasonal_baseline(metric_expr, timestamp)

        if len(historical) < 10:
            return False  # Not enough data for statistical significance

        mean = sum(historical) / len(historical)
        variance = sum((x - mean) ** 2 for x in historical) / len(historical)
        std_dev = variance ** 0.5

        if std_dev == 0:
            return current_value != mean

        z_score = abs(current_value - mean) / std_dev
        return z_score > 3.0  # 3-sigma threshold

    def get_seasonal_baseline(self, expr, timestamp):
        """Get values from same time-of-week for past N weeks."""
        values = []
        for week in range(1, self.lookback_weeks + 1):
            past_ts = timestamp - (week * 7 * 86400)
            result = self.query.instant_query(expr, at=past_ts)
            if result:
                values.append(result[0]['value'])
        return values
```

### 11.5 Custom Dashboards with Template Variables

```
Dashboard JSON model (Grafana-style):
{
  "title": "Service Overview",
  "variables": [
    { "name": "service", "query": "label_values(http_requests_total, service)" },
    { "name": "region",  "query": "label_values(http_requests_total, region)" }
  ],
  "panels": [
    {
      "title": "Request Rate",
      "expr": "sum(rate(http_requests_total{service='$service', region='$region'}[5m]))",
      "type": "timeseries"
    },
    {
      "title": "Error Rate %",
      "expr": "sum(rate(http_requests_total{service='$service', status=~'5..'}[5m])) / sum(rate(http_requests_total{service='$service'}[5m])) * 100",
      "type": "gauge",
      "thresholds": [1, 5]
    }
  ]
}
```

### 11.6 Log-Metrics Correlation

```
When a metric alert fires, automatically surface related logs:

  Metric alert: error_rate{service="payments"} > 5%
    │
    ▼
  Query logs index: service=payments AND level=ERROR AND time in [alert_start, now]
    │
    ▼
  Return top log patterns:
    - "Payment gateway timeout" (230 occurrences)
    - "Database connection pool exhausted" (45 occurrences)
    
  Display in alert notification alongside the metric graph.
```

---

## 12. Interview Q&A

### Q1: How do you handle 10M data points/sec ingestion?

**Candidate:**  
"The key is a multi-stage pipeline that separates concerns:

1. **Ingestion gateways** are stateless — we can horizontally scale to dozens of instances behind a load balancer. They validate, normalize labels, and batch data.
2. **Kafka** acts as a durable buffer between ingestion and storage. With 512 partitions and LZ4 compression, it handles burst traffic without backpressuring the gateways.
3. **Ingestion workers** consume from Kafka and batch-write to the TSDB. Each worker handles ~100K–200K points/sec, so 50-100 workers handle the full load.
4. **The TSDB write path** is append-only: write to WAL (sequential disk I/O), then to an in-memory head block. No random I/O, no B-tree updates.
5. **Batch writes** — we don't write one point at a time. Workers accumulate points per shard and flush in batches of 10K+ points, amortizing the per-write overhead.

The combination of Kafka buffering + append-only writes + batching makes 10M points/sec achievable on commodity hardware."

---

### Q2: Explain Gorilla compression for time-series data.

**Candidate:**  
"Gorilla compression, from the Facebook paper, exploits two properties of time-series data:

**For timestamps:** Metrics report at regular intervals (e.g., every 10 seconds). Instead of storing absolute timestamps (8 bytes each), we store the delta-of-delta. For perfectly regular intervals, the delta-of-delta is 0, which encodes as a single bit. Even with jitter, most deltas are small and encode in 1-12 bits.

**For values:** Consecutive values of the same metric tend to be similar (CPU doesn't jump from 72% to 3000%). We XOR consecutive IEEE 754 float representations. Similar values produce XORs with long runs of leading and trailing zeros. We encode just the meaningful bits plus a small header.

**Result:** From 16 bytes raw (8B timestamp + 8B value), we achieve approximately 1.37 bytes per point on average — a ~12× compression ratio. This is what makes it feasible to store billions of data points per day on a reasonable number of nodes."

---

### Q3: How do you handle high-cardinality labels?

**Candidate:**  
"High cardinality — labels with many unique values like user_id or request_id — is the most common operational failure mode in monitoring systems. Each unique label combination creates a new time series, and the inverted index and in-memory structures grow linearly with series count.

Our defense is layered:

1. **Per-metric cardinality limits** — if a single metric already has 10K unique series, reject new label combinations. This catches accidental cardinality bombs (someone adds user_id to a metric).
2. **Tenant-level quotas** — each team has a maximum active series count (e.g., 500K). This prevents one team from affecting the entire cluster.
3. **Pre-aggregation at the agent** — instead of `requests{user_id=X}`, the agent computes `sum(requests) by (service, status)` locally before sending. This reduces cardinality by orders of magnitude.
4. **Bloom filters** for fast 'does this series exist?' checks, avoiding full index lookups for every incoming point.
5. **Stale series detection** — series with no data for 15+ minutes are marked stale and evicted from the head block to free memory.
6. **Cardinality dashboards and alerts** — we monitor the top-N metrics by series count and alert when growth rate is abnormal."

---

### Q4: How would you design the alerting engine to avoid false positives?

**Candidate:**  
"False positives erode trust in the alerting system faster than anything. Our design addresses this at multiple levels:

1. **'for' duration** — alerts must be in violation for a configurable duration (e.g., 5 minutes) before firing. This filters transient spikes. The state machine goes OK → PENDING → FIRING, and resets to OK if the condition clears during the pending window.

2. **Evaluation at multiple intervals** — we don't fire on a single evaluation. We evaluate every 15 seconds, and the 'for' clause requires sustained violation across multiple consecutive evaluations.

3. **Hysteresis** — the threshold to resolve an alert can be different from the threshold to trigger it. For example, fire when CPU > 90%, but only resolve when CPU < 80%. This prevents flapping.

4. **Aggregation over time** — alerting on `avg(cpu[5m]) > 90%` rather than instantaneous values smooths out noise.

5. **Anomaly-aware thresholds** — for metrics with known daily patterns (e.g., traffic), we use seasonal baselines so that expected traffic drops don't trigger alerts.

6. **Alert grouping and deduplication** — multiple related alerts are grouped into a single notification. If 50 pods alert simultaneously, the engineer gets one message, not 50.

7. **Runbook and context** — every alert includes annotations linking to a runbook, related dashboards, and recent changes. This doesn't prevent the alert but helps resolve it quickly, reducing 'alert fatigue.'"

---

### Q5: Push vs pull model for metric collection — trade-offs?

**Candidate:**

```
┌─────────────────┬────────────────────────────┬────────────────────────────┐
│ Aspect          │ Push Model (StatsD/OTLP)   │ Pull Model (Prometheus)    │
├─────────────────┼────────────────────────────┼────────────────────────────┤
│ Discovery       │ Services register/push on  │ Central config defines     │
│                 │ startup — no discovery     │ targets; needs service     │
│                 │ needed                     │ discovery                  │
├─────────────────┼────────────────────────────┼────────────────────────────┤
│ Firewall/NAT    │ Works behind NAT — service │ Requires network access    │
│                 │ initiates connection       │ FROM monitor TO target     │
├─────────────────┼────────────────────────────┼────────────────────────────┤
│ Short-lived     │ Can emit before exit       │ May miss short-lived jobs  │
│ processes       │ (push on shutdown)         │ between scrape intervals   │
├─────────────────┼────────────────────────────┼────────────────────────────┤
│ Control over    │ Harder — each service      │ Central control — scrape   │
│ load            │ decides when/how much to   │ interval set centrally     │
│                 │ push                       │                            │
├─────────────────┼────────────────────────────┼────────────────────────────┤
│ Liveness        │ Can't distinguish 'no data'│ Missing scrape = target is │
│ detection       │ from 'target is down'      │ confirmed down             │
├─────────────────┼────────────────────────────┼────────────────────────────┤
│ Scale           │ More natural for high-     │ Scraper can become bottle- │
│                 │ cardinality / high-volume  │ neck at extreme scale      │
├─────────────────┼────────────────────────────┼────────────────────────────┤
│ Best for        │ Cloud-native, serverless,  │ Kubernetes, infrastructure │
│                 │ high-throughput apps        │ with known topology        │
└─────────────────┴────────────────────────────┴────────────────────────────┘

In practice, we support both. Push is the primary path for application metrics at scale,
while pull is used for infrastructure monitoring where Prometheus-style scraping is standard.
```

---

### Q6: How do you query across multiple TSDB shards efficiently?

**Candidate:**  
"Cross-shard queries use a scatter-gather pattern:

1. **Query planning** — the query engine parses the PromQL expression and identifies which labels are used. Using the global label index (a lightweight mapping of label values to shard ranges), it determines which shards might hold matching series.

2. **Scatter** — the query engine sends sub-queries to each relevant shard in parallel. Critically, each shard performs local aggregation before returning results. For `avg(cpu) by (host)`, each shard returns its local sum and count per host, not every raw data point.

3. **Gather** — the query engine merges results from all shards. For `avg`, it computes `sum(shard_sums) / sum(shard_counts)` per group. For `max`, it takes `max(shard_maxes)`.

4. **Optimization** — if the query's label matchers constrain the series to a single shard (e.g., `cpu{host='web-01'}` where we shard by series hash), we skip the scatter phase entirely and query that shard directly.

5. **Timeouts and partial results** — if one shard is slow, we return partial results after a timeout (with a warning header) rather than failing the entire query. Dashboards can show 'partial data' indicators."

---

### Q7: How do you handle late-arriving data points?

**Candidate:**  
"Late data is inevitable in distributed systems — network delays, clock skew, batch processing pipelines sending metrics late.

Our approach:

1. **Accept within a window** — we accept data points with timestamps up to 1 hour in the past. Points older than that are rejected (they'd require reopening compacted blocks).

2. **Head block design** — the in-memory head block accepts out-of-order writes within its time range. We maintain a sorted structure (skip list or similar) rather than assuming strict ordering.

3. **Out-of-order ingestion** — for the on-disk path, recent versions of Prometheus-style TSDBs support an 'out-of-order' head block that can accept samples with timestamps earlier than the last appended sample. This adds some write overhead but handles the common case.

4. **Server-side timestamping option** — for clients that can't be trusted to report accurate timestamps, we offer a mode where the ingestion gateway assigns timestamps on receipt. This trades precision for reliability.

5. **Backfill API** — for batch pipelines that need to write historical data, we provide a separate backfill endpoint that writes directly to chunk files (bypassing the head block), with careful locking to avoid conflicts with compaction."

---

### Q8: How would you implement anomaly detection on metrics?

**Candidate:**  
"Anomaly detection on metrics has unique challenges compared to general ML problems: metrics have strong seasonality (daily/weekly patterns), the system needs to handle millions of series, and false positives are costly.

My approach is tiered:

1. **Statistical methods first** — for most metrics, a Z-score against a seasonal baseline works well. We compute the mean and standard deviation of the same metric at the same time-of-week for the past 4 weeks. If the current value deviates by more than 3σ, it's anomalous. This is cheap to compute and interpretable.

2. **EWMA for trend detection** — Exponentially Weighted Moving Average detects gradual drift that Z-score might miss. If the smoothed average of error_rate is trending upward over the past hour, we alert even if no single point exceeds a threshold.

3. **Pre-computed baselines** — we use recording rules to compute hourly/daily baselines during off-peak times. These are stored as regular time series and referenced during alert evaluation, avoiding expensive real-time computation.

4. **ML for multivariate anomalies** — for advanced use cases, an Isolation Forest model trained on multiple related metrics (CPU + memory + request rate + error rate) can detect anomalies that single-metric detection misses. This runs as a batch job, producing anomaly scores stored as metrics.

5. **Human feedback loop** — when users dismiss or acknowledge anomaly alerts, we feed that back to adjust sensitivity. This is critical for reducing false positives over time."

---

## 13. Production Checklist

### Pre-Launch

- [ ] Load test ingestion at 2× expected peak (20M points/sec)
- [ ] Verify Gorilla compression ratio is within expected range (10-12×)
- [ ] Validate WAL recovery: kill TSDB node, restart, verify no data loss
- [ ] Test alerting end-to-end: rule → evaluation → notification delivery
- [ ] Verify cardinality limits reject correctly and return clear error messages
- [ ] Configure dead man's switch alert (always-firing heartbeat)
- [ ] Set up meta-monitoring on separate infrastructure
- [ ] Test Kafka retention and consumer lag recovery
- [ ] Validate query timeouts and circuit breakers
- [ ] Security: TLS for agent-to-gateway, authentication for query API

### Day 1

- [ ] Monitor ingestion rate vs. expected baseline
- [ ] Watch Kafka consumer lag — should be near zero
- [ ] Verify first compaction cycle completes without errors
- [ ] Check active series count matches expectations
- [ ] Confirm alerting latency (time from condition met to notification received)
- [ ] Monitor TSDB disk usage growth rate
- [ ] Validate dashboard query latencies (p50, p99)

### Week 1

- [ ] First downsampling cycle completes (10s → 1min rollups)
- [ ] Verify historical queries work correctly against rollup data
- [ ] Tune alert thresholds based on initial false positive rate
- [ ] Review cardinality reports — identify top metrics by series count
- [ ] Test TSDB node failure and replica promotion
- [ ] Validate backup and restore procedures
- [ ] Review query patterns and add recording rules for expensive queries

### Month 1

- [ ] Review storage growth vs. estimates — adjust retention if needed
- [ ] Optimize shard distribution — rebalance if hotspots emerge
- [ ] Implement recording rules for the top 20 most expensive dashboard queries
- [ ] Set up SLO tracking for the monitoring system itself
- [ ] Conduct chaos engineering: kill Kafka broker, TSDB node, alerting engine
- [ ] Review alert noise — tune or remove alerts with low signal-to-noise ratio
- [ ] Plan capacity for next quarter based on growth trends

---

## Summary

| Aspect | Decision | Rationale |
|--------|----------|-----------|
| Ingestion protocol | Multi-protocol (StatsD, Prometheus, OTLP, HTTP) | Meet services where they are |
| Write buffer | Kafka (512 partitions, RF=3) | Absorb bursts, decouple ingestion from storage |
| Storage engine | Custom TSDB with Gorilla compression | 12× compression, append-only writes, O(1) deletion |
| Time-series index | Inverted index (label → series IDs) | O(log n) label lookups, fast intersection |
| Query language | PromQL-compatible | Industry standard, expressive aggregations |
| Downsampling | Multi-tier (10s → 1m → 5m → 1h) | Balance storage cost with query granularity |
| Alerting | Rule evaluation loop + state machine | Reliable, avoidance of false positives via 'for' duration |
| Sharding | Consistent hash on series_id | Even distribution, all points for a series on same shard |
| High availability | Replication + redundant alert evaluators | No single point of failure for critical paths |
| Cardinality control | Per-metric limits + tenant quotas + bloom filters | Prevent cardinality bombs from impacting cluster |

### Scalability Path

```
Stage 1 (MVP):         1M points/sec, 1M series, single cluster
Stage 2 (Growth):      10M points/sec, 10M series, sharded TSDB
Stage 3 (Scale):       100M points/sec, 100M series, multi-cluster federation
Stage 4 (Enterprise):  1B+ points/sec, multi-region, tiered storage (SSD + S3)
```

---

> **Key Takeaway:** A monitoring system is a write-heavy, time-series-optimized data platform. The core differentiators are (1) Gorilla compression for storage efficiency, (2) inverted index for fast label queries, (3) a reliable alerting engine with a well-designed state machine, and (4) cardinality management to prevent the system from eating itself. In an interview, demonstrate that you understand *why* general-purpose databases fail at this workload and how each component is tuned for the time-series access pattern.
