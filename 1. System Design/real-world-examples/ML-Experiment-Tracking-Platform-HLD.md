# ML Experiment Tracking & Analysis Platform — High-Level Design

---

## 1. Problem Statement

Design a scalable **ML Experiment Tracking and Analysis Platform** (similar to MLflow, Weights & Biases, Neptune.ai) that enables data scientists and ML engineers to **log, compare, visualize, and reproduce** machine learning experiments — including hyperparameters, metrics, artifacts, datasets, and model versions — across teams and projects at enterprise scale.

**Example:**
```
Data Scientist runs 200 training experiments for a fraud detection model:

Platform captures:
 - Hyperparameters: learning_rate=0.001, batch_size=64, epochs=50, optimizer=adam
 - Metrics (per step): loss=0.032, accuracy=0.974, f1=0.968, AUC=0.991
 - Artifacts: model.pth (230MB), confusion_matrix.png, SHAP plots
 - Environment: Python 3.11, PyTorch 2.1, CUDA 12.1, requirements.txt
 - Dataset: fraud_v3.parquet (hash: a3f8c1...), 2.1M rows
 - Git commit: a1b2c3d
 - Duration: 4h 23m on 4× A100 GPUs

Team can: compare runs side-by-side, find best model, promote to staging registry
```

---

## 2. Functional Requirements

| # | Requirement | Description |
|---|-------------|-------------|
| FR-1 | Experiment & Run Management | Create experiments (logical groups), start/stop runs, tag and annotate |
| FR-2 | Parameter Logging | Log hyperparameters, config files, and environment variables per run |
| FR-3 | Metric Logging | Log scalar metrics (loss, accuracy) at step/epoch granularity with time series |
| FR-4 | Artifact Storage | Upload/download artifacts: model weights, plots, datasets, notebooks |
| FR-5 | Run Comparison | Side-by-side comparison of runs — parallel coordinates, scatter plots, tables |
| FR-6 | Model Registry | Version, stage (staging/production/archived), and promote trained models |
| FR-7 | Search & Filter | Query runs by parameters, metrics, tags, status, duration, user |
| FR-8 | Collaboration | Shared workspaces, team projects, comments on runs, @mentions |
| FR-9 | Reproducibility | Capture git commit, environment snapshot, random seeds, data lineage |
| FR-10 | Visualization Dashboard | Real-time training curves, custom charts, system resource utilization |
| FR-11 | Alerting & Notifications | Alert on training completion, metric thresholds (e.g., loss > 1.0), failures |
| FR-12 | SDK / Client Libraries | Python SDK, REST API, CLI — seamless integration with PyTorch, TF, scikit-learn |
| FR-13 | Dataset Versioning | Track dataset versions, lineage, and association with experiment runs |
| FR-14 | Access Control | Project-level RBAC — viewer, contributor, admin roles |

---

## 3. Non-Functional Requirements

| # | Requirement | Target |
|---|-------------|--------|
| NFR-1 | Availability | 99.9% uptime (platform should not block training pipelines) |
| NFR-2 | Metric Ingestion Latency | < 100ms for log calls (non-blocking to training loop) |
| NFR-3 | Query Latency | p99 < 500ms for run search and metric retrieval |
| NFR-4 | Scalability | Support 10K+ concurrent training jobs, millions of runs, billions of metric points |
| NFR-5 | Artifact Storage | Handle artifacts up to 10 GB per run, petabytes total |
| NFR-6 | Durability | Zero data loss for logged metrics, parameters, and artifacts |
| NFR-7 | Consistency | Eventual consistency for dashboards; strong consistency for model registry transitions |
| NFR-8 | Multi-tenancy | Isolated data per organization/team with shared infrastructure |

---

## 4. Capacity Estimation

### 4.1 Traffic

| Metric | Value |
|--------|-------|
| Organizations | ~5,000 |
| Data scientists per org (avg) | ~20 |
| Total active users | ~100,000 |
| Concurrent training runs (peak) | ~10,000 |
| Metric log calls per run per second | ~10 (1 per step × batches) |
| **Peak metric ingestion QPS** | **~100,000** |
| API read requests (dashboards, search) | ~20,000 QPS |
| Artifact uploads per day | ~50,000 (avg 100 MB each) |
| New runs per day | ~200,000 |

### 4.2 Storage

| Data | Estimate |
|------|----------|
| Metric data points per run (avg) | ~50,000 (1000 steps × 50 metrics) |
| Metric data point size | ~50 bytes (run_id, metric_name, step, value, timestamp) |
| Daily metric data | 200K runs × 50K points × 50B = **~500 GB/day** |
| Monthly metric data | ~15 TB |
| Artifact storage per day | 50K × 100MB = **~5 TB/day** |
| Monthly artifact storage | ~150 TB |
| Run metadata per month | ~5 GB |
| Total storage (1 year) | ~2 PB (mostly artifacts) |

### 4.3 Bandwidth

| Direction | Calculation | Result |
|-----------|-------------|--------|
| Metric ingestion (inbound) | 100K QPS × 200 bytes/call | ~20 MB/s |
| Artifact upload (inbound) | ~5 TB/day | ~60 MB/s avg, ~500 MB/s peak |
| Dashboard reads (outbound) | 20K QPS × 2 KB avg response | ~40 MB/s |
| Artifact download (outbound) | ~2 TB/day | ~25 MB/s avg |

---

## 5. Core Data Model

### 5.1 Entity Relationships

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│ Organization │──1:N──│   Project    │──1:N──│  Experiment  │
│              │       │              │       │              │
│ org_id       │       │ project_id   │       │ experiment_id│
│ name         │       │ org_id (FK)  │       │ project_id   │
│ plan         │       │ name         │       │ name         │
│ created_at   │       │ description  │       │ description  │
└──────────────┘       │ visibility   │       │ tags[]       │
                       └──────────────┘       └──────┬───────┘
                                                     │
                                                    1:N
                                                     │
                                              ┌──────▼───────┐
                                              │     Run      │
                                              │              │
                                              │ run_id (UUID)│
                                              │ experiment_id│
                                              │ user_id      │
                                              │ status       │ ← RUNNING|COMPLETED|FAILED|KILLED
                                              │ start_time   │
                                              │ end_time     │
                                              │ git_commit   │
                                              │ source_name  │
                                              │ tags[]       │
                                              └──┬───┬───┬───┘
                                                 │   │   │
                              ┌───────────────┬──┘   │   └──┬─────────────────┐
                              │               │      │      │                 │
                       ┌──────▼──────┐ ┌──────▼──┐ ┌─▼────┐ │  ┌──────────────▼──┐
                       │  Parameter  │ │ Metric  │ │Metric│ │  │   Artifact      │
                       │             │ │(Summary)│ │(Step)│ │  │                 │
                       │ run_id      │ │         │ │      │ │  │ run_id          │
                       │ key         │ │ run_id  │ │run_id│ │  │ path            │
                       │ value       │ │ key     │ │key   │ │  │ storage_uri     │
                       │ type        │ │ min     │ │step  │ │  │ size_bytes      │
                       └─────────────┘ │ max     │ │value │ │  │ content_type    │
                                       │ last    │ │ts    │ │  │ hash (SHA-256)  │
                                       │ count   │ └──────┘ │  └─────────────────┘
                                       └─────────┘          │
                                                     ┌──────▼────────┐
                                                     │  Environment  │
                                                     │               │
                                                     │ run_id        │
                                                     │ python_version│
                                                     │ packages[]    │
                                                     │ docker_image  │
                                                     │ gpu_info      │
                                                     │ os_info       │
                                                     └───────────────┘
```

### 5.2 Model Registry Entities

```
┌──────────────────┐       ┌──────────────────┐       ┌──────────────────┐
│ RegisteredModel  │──1:N──│  ModelVersion    │──N:1──│      Run         │
│                  │       │                  │       │                  │
│ name (unique)    │       │ model_name       │       │ run_id           │
│ description      │       │ version (auto)   │       │ (source of       │
│ tags[]           │       │ run_id (FK)      │       │  model artifact) │
│ created_at       │       │ artifact_path    │       └──────────────────┘
│ last_updated     │       │ stage            │ ← NONE|STAGING|PRODUCTION|ARCHIVED
└──────────────────┘       │ status           │ ← PENDING_REGISTRATION|READY|FAILED
                           │ description      │
                           │ created_at       │
                           └──────────────────┘
```

---

## 6. High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              CLIENTS                                        │
│   Python SDK  |  REST API  |  CLI  |  Web Dashboard  |  Jupyter Extension   │
│                                                                             │
│   sdk.log_param("lr", 0.001)                                                │
│   sdk.log_metric("loss", 0.032, step=100)                                   │
│   sdk.log_artifact("model.pth")                                             │
└──────────────────────────────┬──────────────────────────────────────────────┘
                               │
                          ┌────▼─────┐
                          │   CDN    │  (Dashboard static assets, cached charts)
                          └────┬─────┘
                               │
                       ┌───────▼────────┐
                       │  API Gateway   │  Auth, rate limiting, routing
                       │ (Kong / Envoy) │  SDK write path → Ingestion API
                       └───────┬────────┘  Dashboard reads → Query API
                               │
            ┌──────────────────┼────────────────────────┐
            │                  │                         │
     ┌──────▼───────┐  ┌──────▼───────┐  ┌─────────────▼──────────┐
     │  Ingestion   │  │  Query       │  │   Artifact             │
     │  Service     │  │  Service     │  │   Service              │
     │              │  │              │  │                        │
     │ - log params │  │ - search     │  │ - upload (presigned)   │
     │ - log metrics│  │   runs       │  │ - download             │
     │ - log tags   │  │ - get run    │  │ - list artifacts       │
     │ - start/end  │  │ - compare    │  │ - delete               │
     │   runs       │  │ - dashboards │  │                        │
     └──────┬───────┘  └──────┬───────┘  └─────────────┬──────────┘
            │                  │                         │
            │          ┌───────┘                         │
            ▼          ▼                                 ▼
     ┌─────────────────────┐                   ┌────────────────┐
     │  Message Queue      │                   │  Object Store  │
     │  (Kafka / Pulsar)   │                   │  (S3 / GCS /   │
     │                     │                   │   MinIO)        │
     │  Topics:            │                   │                │
     │  - metric-events    │                   │  /org/project/ │
     │  - run-events       │                   │   experiment/  │
     │  - model-events     │                   │    run/        │
     └──────┬──────────────┘                   │     artifacts/ │
            │                                  └────────────────┘
            │
     ┌──────▼───────────────────────────────────────────────────┐
     │                    PROCESSING LAYER                       │
     │                                                          │
     │  ┌────────────────┐  ┌────────────────┐  ┌────────────┐ │
     │  │ Metric Writer  │  │ Aggregation    │  │  Alert     │ │
     │  │ (Consumer)     │  │ Service        │  │  Engine    │ │
     │  │                │  │                │  │            │ │
     │  │ Batch writes   │  │ Pre-compute    │  │ Threshold  │ │
     │  │ to TimeSeries  │  │ metric summary │  │ & anomaly  │ │
     │  │ DB             │  │ (min/max/avg)  │  │ detection  │ │
     │  └────────┬───────┘  └────────┬───────┘  └─────┬──────┘ │
     └───────────┼──────────────────┼──────────────────┼────────┘
                 │                  │                   │
                 ▼                  ▼                   ▼
     ┌───────────────────────────────────────────────────────────┐
     │                     DATA LAYER                            │
     │                                                           │
     │  ┌──────────────┐  ┌──────────────┐  ┌────────────────┐  │
     │  │  PostgreSQL   │  │ TimeSeries DB│  │    Redis       │  │
     │  │  (Metadata)   │  │ (InfluxDB /  │  │   (Cache)      │  │
     │  │               │  │  TimescaleDB)│  │                │  │
     │  │ - runs        │  │              │  │ - hot metrics  │  │
     │  │ - params      │  │ - step-level │  │ - run status   │  │
     │  │ - experiments │  │   metrics    │  │ - search cache │  │
     │  │ - model reg.  │  │ - system     │  │ - session      │  │
     │  │ - users/orgs  │  │   resource   │  │                │  │
     │  └──────────────┘  │   metrics    │  └────────────────┘  │
     │                     └──────────────┘                      │
     │  ┌──────────────┐  ┌──────────────┐                      │
     │  │Elasticsearch │  │  S3 / GCS    │                      │
     │  │ (Run Search) │  │ (Artifacts)  │                      │
     │  │              │  │              │                      │
     │  │ - full-text  │  │ - models     │                      │
     │  │ - param/     │  │ - plots      │                      │
     │  │   metric     │  │ - datasets   │                      │
     │  │   filtering  │  │ - checkpoints│                      │
     │  └──────────────┘  └──────────────┘                      │
     └───────────────────────────────────────────────────────────┘
            │
     ┌──────▼───────────────────────────────────────────────────┐
     │                  ADDITIONAL SERVICES                      │
     │                                                          │
     │  ┌────────────────┐  ┌────────────────┐  ┌────────────┐ │
     │  │ Model Registry │  │ Notification   │  │  Auth      │ │
     │  │ Service        │  │ Service        │  │  Service   │ │
     │  │                │  │                │  │            │ │
     │  │ - register     │  │ - email        │  │ - JWT/OAuth│ │
     │  │ - version      │  │ - Slack        │  │ - RBAC     │ │
     │  │ - stage        │  │ - webhook      │  │ - API keys │ │
     │  │   transition   │  │ - PagerDuty    │  │ - teams    │ │
     │  └────────────────┘  └────────────────┘  └────────────┘ │
     └──────────────────────────────────────────────────────────┘
```

---

## 7. Component Deep Dive

### 7.1 SDK / Client Library (Python)

The SDK is the primary interface for data scientists. It must be **non-blocking** — logging calls should never slow down the training loop.

```
┌──────────────────────────────────────────────────┐
│                 Python SDK                         │
│                                                    │
│  ┌──────────────┐   ┌──────────────────────────┐  │
│  │  Public API   │   │  Background Thread Pool  │  │
│  │               │   │                          │  │
│  │  log_param()  │──►│  ┌──────────────────┐    │  │
│  │  log_metric() │   │  │  In-Memory Buffer │    │  │
│  │  log_artifact│   │  │  (batch queue)     │    │  │
│  │  start_run() │   │  └────────┬───────────┘    │  │
│  │  end_run()   │   │           │                │  │
│  │  set_tag()   │   │  ┌────────▼───────────┐    │  │
│  └──────────────┘   │  │  Batch Sender      │    │  │
│                      │  │  (flush every 1s   │    │  │
│                      │  │   or 100 items)    │    │  │
│                      │  └────────┬───────────┘    │  │
│                      └───────────┼────────────────┘  │
│                                  │                    │
│   ┌──────────────────────────────▼─────────────────┐ │
│   │  Retry Logic (exp backoff, local disk fallback)│ │
│   └────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

**Key design decisions:**
| Decision | Choice | Rationale |
|----------|--------|-----------|
| Async logging | Background thread + buffer | Never blocks training loop |
| Batching | Flush every 1s or 100 items | Reduces HTTP overhead by 100× |
| Failure handling | Retry → local disk fallback | Training must not fail if platform is down |
| Artifact upload | Presigned URL → direct to S3 | Bypass API server for large files |
| Auto-logging | Framework integrations | `sdk.autolog()` for PyTorch/TF/sklearn |

### 7.2 Ingestion Service

Handles the write-heavy path — all parameter, metric, and run state logging.

**Request flow:**
1. SDK sends batched log request (params + metrics + tags)
2. Ingestion service validates and authenticates
3. Writes run metadata (params, tags, status) to PostgreSQL
4. Publishes metric events to Kafka topic `metric-events`
5. Returns acknowledgment immediately (async processing)

```
SDK Batch Request
    │
    ▼
Ingestion Service
    │
    ├──► PostgreSQL (params, tags, run status) — synchronous
    │
    ├──► Kafka: metric-events (step-level metrics) — async, fire-and-forget
    │
    └──► Redis (update run status cache, latest metrics) — async
```

**Why Kafka for metrics?** At 100K QPS, writing each metric point directly to the time-series DB would overwhelm it. Kafka acts as a buffer, and the Metric Writer consumer does batched inserts.

### 7.3 Metric Writer (Kafka Consumer)

Consumes from `metric-events` topic, buffers, and performs bulk writes to the time-series database.

```
Kafka: metric-events
    │
    ▼
┌───────────────────────────────────┐
│        Metric Writer              │
│                                   │
│  Consumer Group (N partitions)    │
│  ┌──────────────────────────┐     │
│  │ Micro-batch accumulator  │     │
│  │ (5s window or 10K points)│     │
│  └────────────┬─────────────┘     │
│               │                   │
│  ┌────────────▼─────────────┐     │
│  │ Bulk INSERT into         │     │
│  │ TimescaleDB/InfluxDB     │     │
│  │ (10K points per batch)   │     │
│  └──────────────────────────┘     │
└───────────────────────────────────┘
```

**Partitioning strategy:** Kafka topic partitioned by `run_id` — ensures all metrics for a run are processed in order and by the same consumer.

### 7.4 Query Service

Serves the read path — dashboard rendering, run comparison, search.

```
Dashboard / API Request
    │
    ▼
Query Service
    │
    ├── Run metadata, params → PostgreSQL
    │
    ├── Step-level metrics (training curves) → TimeSeries DB
    │
    ├── Full-text search, filtered queries → Elasticsearch
    │
    └── Hot data (live run metrics, status) → Redis Cache
```

**Query optimization patterns:**

| Query Type | Source | Optimization |
|------------|--------|-------------|
| Get run details | PostgreSQL | Indexed on run_id (PK) |
| Training curve (loss over steps) | TimescaleDB | Hypertable partitioned by time; downsampling for long runs |
| Compare N runs | TimescaleDB + PostgreSQL | Parallel fan-out; result merge |
| Search "runs where accuracy > 0.95" | Elasticsearch | Denormalized run doc with params + summary metrics |
| Live training metrics | Redis | Pub/Sub for real-time chart updates |
| Run table (paginated list) | Elasticsearch | Sorted by start_time, filtered by experiment |

### 7.5 Artifact Service

Handles large binary files (models, plots, datasets) without routing them through the API server.

```
SDK: sdk.log_artifact("model.pth", 230MB)
    │
    │  Step 1: Request presigned upload URL
    ▼
┌─────────────────┐
│ Artifact Service │──► Generate presigned PUT URL (S3, 15 min TTL)
└────────┬────────┘
         │
         │  Return presigned URL
         ▼
SDK uploads directly to S3
    │
    │  Step 2: Confirm upload
    ▼
┌─────────────────┐
│ Artifact Service │──► Record artifact metadata in PostgreSQL
└─────────────────┘     (run_id, path, size, hash, storage_uri)
```

**Why presigned URLs?**
- API servers never handle large file bytes — no memory/bandwidth bottleneck
- SDK uploads directly to object store at full network speed
- Scales horizontally without API server being the bottleneck

**Artifact organization in object store:**
```
s3://ml-platform-artifacts/
  └── {org_id}/
      └── {project_id}/
          └── {experiment_id}/
              └── {run_id}/
                  ├── model/
                  │   ├── model.pth
                  │   └── model_config.json
                  ├── plots/
                  │   ├── confusion_matrix.png
                  │   └── training_curves.html
                  ├── data/
                  │   └── dataset_manifest.json
                  └── environment/
                      ├── requirements.txt
                      └── conda.yaml
```

### 7.6 Model Registry Service

Manages the lifecycle of trained models from experimentation to production.

```
┌────────────────────────────────────────────────────────────────┐
│                    MODEL LIFECYCLE                              │
│                                                                │
│   Training Run                                                 │
│       │                                                        │
│       ▼                                                        │
│   ┌──────────┐    ┌───────────┐    ┌────────────┐    ┌──────┐ │
│   │   NONE   │───►│  STAGING  │───►│ PRODUCTION │───►│ARCHVD│ │
│   │          │    │           │    │            │    │      │ │
│   │ Logged   │    │ Validation│    │ Serving    │    │ Old  │ │
│   │ artifact │    │ & testing │    │ traffic    │    │ ver. │ │
│   └──────────┘    └───────────┘    └────────────┘    └──────┘ │
│                                                                │
│   Stage transitions require:                                   │
│   - RBAC permission (admin/ML lead)                            │
│   - Approval workflow (optional)                               │
│   - Validation checks (metrics thresholds, bias tests)         │
│   - Audit log entry                                            │
└────────────────────────────────────────────────────────────────┘
```

**Registry features:**

| Feature | Description |
|---------|-------------|
| Auto-versioning | Each registration auto-increments version number |
| Stage gates | Configurable checks before stage promotion |
| Model lineage | Links to source run, dataset, git commit |
| Serving integration | Webhook on PRODUCTION promotion → trigger deployment |
| Rollback | One-click revert to previous production version |
| Model card | Auto-generated documentation (metrics, data, limitations) |

### 7.7 Alert Engine

Monitors active training runs and triggers notifications.

```
Kafka: metric-events + run-events
    │
    ▼
┌─────────────────────────────────────────────┐
│              Alert Engine                     │
│                                              │
│  ┌──────────────────────┐                    │
│  │  Rule Evaluator      │                    │
│  │                      │                    │
│  │  - loss > threshold  │                    │
│  │  - NaN detected      │  ┌───────────────┐ │
│  │  - run exceeded      │──► Notification  │ │
│  │    max duration      │  │ Service       │ │
│  │  - GPU util < 10%    │  │ (Slack/Email/ │ │
│  │  - run completed     │  │  Webhook)     │ │
│  │  - model promoted    │  └───────────────┘ │
│  └──────────────────────┘                    │
│                                              │
│  ┌──────────────────────┐                    │
│  │  Anomaly Detector    │                    │
│  │                      │                    │
│  │  - loss spike        │                    │
│  │  - metric plateau    │                    │
│  │  - resource anomaly  │                    │
│  └──────────────────────┘                    │
└─────────────────────────────────────────────┘
```

---

## 8. Data Flow — End to End

### 8.1 Experiment Run Lifecycle (Write Path)

```
Data Scientist                          Platform
    │                                      │
    │  1. sdk.init(project="fraud")        │
    │─────────────────────────────────────►│  Auth → create/get experiment
    │  ◄─── experiment_id ─────────────────│
    │                                      │
    │  2. run = sdk.start_run()            │
    │─────────────────────────────────────►│  Create run in PostgreSQL (status=RUNNING)
    │  ◄─── run_id (UUID) ────────────────│   Capture git_commit, environment
    │                                      │
    │  3. sdk.log_params({lr: 0.001,       │
    │         batch: 64, optimizer: adam})  │
    │─────────────────────────────────────►│  Batch insert into PostgreSQL (params table)
    │                                      │
    │  4. [Training loop]                  │
    │     for epoch in range(50):          │
    │       sdk.log_metric("loss", val,    │
    │                      step=epoch)     │
    │─────── (buffered, batch every 1s) ──►│  → Kafka: metric-events
    │                                      │  → Metric Writer → TimescaleDB
    │                                      │  → Redis (latest values for live dashboard)
    │                                      │
    │  5. sdk.log_artifact("model.pth")    │
    │─────────────────────────────────────►│  Presigned URL → SDK uploads to S3
    │      [SDK uploads 230MB to S3]       │  Artifact metadata saved to PostgreSQL
    │                                      │
    │  6. sdk.end_run(status="COMPLETED")  │
    │─────────────────────────────────────►│  Update run status in PostgreSQL
    │                                      │  Flush remaining buffered metrics
    │                                      │  Kafka: run-events (for alerting)
    │                                      │  Sync to Elasticsearch (run doc)
    │                                      │
```

### 8.2 Dashboard Query Path (Read)

```
Dashboard User
    │
    │  1. Open experiment "fraud-detection"
    │─────────────────────────────────────►  Query Service
    │                                        → Elasticsearch: get runs for experiment
    │  ◄─── run list (paginated) ─────────   (params, summary metrics, status)
    │
    │  2. Select 3 runs to compare
    │─────────────────────────────────────►  Query Service
    │                                        → TimescaleDB: get loss/accuracy curves
    │  ◄─── metric time series ───────────   → PostgreSQL: get full params
    │                                        → Redis: live metrics for RUNNING runs
    │
    │  3. View training curves
    │   [Browser renders interactive charts]
    │   [WebSocket for live run updates]
    │
    │  4. Search: "runs where f1 > 0.96    
    │             AND lr < 0.01"            
    │─────────────────────────────────────►  Elasticsearch: filtered query
    │  ◄─── matching runs ────────────────
    │
```

---

## 9. Database Strategy

### 9.1 Database Selection

| Database | Use Case | Why This Choice |
|----------|----------|-----------------|
| **PostgreSQL** | Run metadata, params, tags, experiments, users, model registry | ACID for registry transitions; relational integrity; rich querying |
| **TimescaleDB** (or InfluxDB) | Step-level metric time series | Purpose-built for time-series; hypertables with auto-partitioning; 10× compression; downsampling |
| **Elasticsearch** | Run search and filtering | Complex queries across params+metrics+tags; full-text search on run names/notes |
| **Redis** | Live run cache, session, rate limiting | Sub-ms reads for dashboard live updates; Pub/Sub for real-time charts |
| **S3 / GCS** | Artifact storage (models, plots, data) | Virtually unlimited storage; cost-effective; lifecycle policies |
| **Kafka** | Metric event streaming | Decouple ingestion from storage; buffer for burst traffic; exactly-once semantics |

### 9.2 Key Schema Decisions

**PostgreSQL — Runs table (partitioned by created month):**
```
runs
├── run_id          UUID (PK)
├── experiment_id   UUID (FK → experiments)
├── user_id         UUID (FK → users)
├── status          ENUM (RUNNING, COMPLETED, FAILED, KILLED)
├── start_time      TIMESTAMPTZ
├── end_time        TIMESTAMPTZ
├── git_commit_hash VARCHAR(40)
├── source_name     TEXT
├── created_at      TIMESTAMPTZ
└── INDEX on (experiment_id, created_at DESC)

params
├── run_id          UUID (FK → runs)
├── key             VARCHAR(255)
├── value           TEXT
├── type            ENUM (STRING, FLOAT, INT, BOOL)
└── UNIQUE INDEX on (run_id, key)
```

**TimescaleDB — Metrics hypertable (auto-partitioned by time):**
```
metrics
├── time            TIMESTAMPTZ
├── run_id          UUID
├── metric_key      VARCHAR(255)
├── step            BIGINT
├── value           DOUBLE PRECISION
└── HYPERTABLE partition by time (1 day chunks)
    INDEX on (run_id, metric_key, step)
```

**Elasticsearch — Denormalized run document:**
```json
{
    "run_id": "uuid",
    "experiment_id": "uuid",
    "experiment_name": "fraud-detection-v3",
    "user": "alice@company.com",
    "status": "COMPLETED",
    "start_time": "2026-03-26T10:00:00Z",
    "duration_seconds": 15780,
    "params": {
        "learning_rate": 0.001,
        "batch_size": 64,
        "optimizer": "adam",
        "model_arch": "transformer"
    },
    "metrics_summary": {
        "loss": { "min": 0.032, "max": 2.341, "last": 0.032 },
        "accuracy": { "min": 0.51, "max": 0.974, "last": 0.974 },
        "f1_score": { "min": 0.48, "max": 0.968, "last": 0.968 }
    },
    "tags": ["production-candidate", "transformer", "v3"],
    "git_commit": "a1b2c3d"
}
```

---

## 10. Real-Time Dashboard Architecture

Live training dashboards require real-time metric streaming to the browser.

```
┌─────────────────────────────────────────────────────────────────┐
│                  REAL-TIME DASHBOARD                             │
│                                                                 │
│  Training Job                                                   │
│      │                                                          │
│      │ sdk.log_metric("loss", 0.032, step=5000)                 │
│      ▼                                                          │
│  Ingestion Service                                              │
│      │                                                          │
│      ├──► Kafka: metric-events                                  │
│      │                                                          │
│      └──► Redis PUBLISH channel: "run:{run_id}:metrics"         │
│               │                                                 │
│               ▼                                                 │
│      ┌─────────────────┐                                        │
│      │  WebSocket       │ ◄──── Redis SUBSCRIBE                 │
│      │  Gateway         │                                       │
│      │  (per dashboard  │                                       │
│      │   connection)    │                                       │
│      └────────┬────────┘                                        │
│               │                                                 │
│               ▼                                                 │
│      Browser (WebSocket client)                                 │
│      ┌──────────────────────────────────┐                       │
│      │  Chart.js / D3 / Plotly          │                       │
│      │  Live-updating training curves   │                       │
│      │  ┌──────────────────────────┐    │                       │
│      │  │  Loss ────╲              │    │                       │
│      │  │           ──╲──╲         │    │                       │
│      │  │               ──╲──────  │    │                       │
│      │  │  0  1000  2000  3000  5K │    │                       │
│      │  │           steps          │    │                       │
│      │  └──────────────────────────┘    │                       │
│      └──────────────────────────────────┘                       │
└─────────────────────────────────────────────────────────────────┘
```

**Downsampling for large runs:** A run with 1M steps would have 1M data points per metric. For rendering, we downsample:

| Steps in Run | Display Strategy |
|-------------|-----------------|
| < 10,000 | Show all points |
| 10K – 100K | LTTB (Largest Triangle Three Buckets) downsample to 5K points |
| > 100K | Progressive loading: coarse view → zoom triggers detail fetch |

---

## 11. Run Comparison Engine

A core differentiator — enabling data scientists to compare runs efficiently.

### 11.1 Comparison Types

```
┌──────────────────────────────────────────────────────────────┐
│                   RUN COMPARISON VIEWS                        │
│                                                              │
│  1. TABLE VIEW                                               │
│  ┌──────────┬──────────┬──────────┬──────────┐               │
│  │ Run      │ Run-A    │ Run-B    │ Run-C    │               │
│  ├──────────┼──────────┼──────────┼──────────┤               │
│  │ lr       │ 0.001    │ 0.01     │ 0.001    │               │
│  │ batch    │ 64       │ 128      │ 256      │               │
│  │ accuracy │ 0.974    │ 0.961    │ 0.980 ★  │               │
│  │ f1       │ 0.968    │ 0.955    │ 0.972 ★  │               │
│  │ duration │ 4h 23m   │ 2h 10m ★ │ 6h 45m   │               │
│  └──────────┴──────────┴──────────┴──────────┘               │
│                                                              │
│  2. PARALLEL COORDINATES                                     │
│     lr ──── batch ──── epochs ──── accuracy ──── f1          │
│     │         │          │           │           │           │
│     ╲─────────╱──────────╲───────────╱───────────╱  Run-A    │
│      ╲───────╱────────────╲─────────╱───────────╱   Run-B    │
│       ╲─────╱──────────────╲───────╱───────────╱    Run-C    │
│                                                              │
│  3. OVERLAY TRAINING CURVES                                  │
│     Loss                                                     │
│     │  ╲  Run-A                                              │
│     │   ╲──╲ Run-B                                           │
│     │     ──╲──── Run-C                                      │
│     │         ──────────                                     │
│     └──────────────────── steps                              │
│                                                              │
│  4. SCATTER PLOT (lr vs accuracy, colored by model_arch)     │
│                                                              │
└──────────────────────────────────────────────────────────────┘
```

### 11.2 Comparison Query Flow

```
Compare(run_ids=[A, B, C])
    │
    ├──► PostgreSQL: SELECT params for run_ids (parallel)
    │    → {run_id → {param_key → value}}
    │
    ├──► TimescaleDB: SELECT metrics for run_ids (parallel, downsampled)
    │    → {run_id → {metric_key → [(step, value), ...]}}
    │
    └──► Merge & return unified comparison response
         ~ 200ms for 10 runs × 5 metrics × 5K steps each
```

---

## 12. Search & Filtering

### 12.1 Query Language

Support a structured query DSL for power users:

```
# Find best performing runs
metrics.accuracy > 0.95 AND params.optimizer = "adam" AND status = "COMPLETED"

# Find runs with specific tags
tags IN ["production-candidate", "v3"] AND params.learning_rate < 0.01

# Time-range filter
created > "2026-03-01" AND experiment.name = "fraud-detection"

# Sort by metric
ORDER BY metrics.f1_score DESC LIMIT 20
```

### 12.2 Search Architecture

```
User Query: "params.lr < 0.01 AND metrics.accuracy > 0.95"
    │
    ▼
Query Parser (DSL → Elasticsearch Query DSL)
    │
    ▼
Elasticsearch
    │
    ├── Index: runs (denormalized run documents)
    │   - params stored as nested fields (typed)
    │   - metric summaries as nested fields
    │   - tags as keyword array
    │
    ├── Query: bool { must: [range(params.lr < 0.01), range(metrics.accuracy.last > 0.95)] }
    │
    └── Return: matching run_ids + highlighted fields
```

**Sync pipeline (PostgreSQL → Elasticsearch):**

```
PostgreSQL (source of truth)
    │
    │  CDC (Change Data Capture) via Debezium
    ▼
Kafka: run-changes topic
    │
    ▼
ES Sync Consumer
    │
    │  Denormalize: join run + params + metric_summary
    ▼
Elasticsearch Index (eventual consistency, ~2s lag)
```

---

## 13. Reproducibility & Lineage

### 13.1 What Gets Captured Per Run

```
┌───────────────────────────────────────────────────────┐
│                REPRODUCIBILITY SNAPSHOT                 │
│                                                        │
│  ┌─────────────────┐  ┌──────────────────────────────┐ │
│  │  Code            │  │  Environment                 │ │
│  │  - git_repo_url  │  │  - python_version: 3.11.7    │ │
│  │  - git_commit:   │  │  - packages:                 │ │
│  │    a1b2c3d       │  │      torch==2.1.0            │ │
│  │  - git_branch:   │  │      numpy==1.26.0           │ │
│  │    feature/v3    │  │      pandas==2.1.4            │ │
│  │  - dirty: false  │  │  - cuda_version: 12.1        │ │
│  │  - diff (if      │  │  - docker_image: (optional)  │ │
│  │    dirty)        │  │  - os: Linux 6.1              │ │
│  └─────────────────┘  └──────────────────────────────┘ │
│                                                        │
│  ┌─────────────────┐  ┌──────────────────────────────┐ │
│  │  Data             │  │  Compute                    │ │
│  │  - dataset_hash:  │  │  - gpu: 4× A100 80GB       │ │
│  │    sha256:a3f8c1  │  │  - cpu: 64 cores            │ │
│  │  - dataset_uri:   │  │  - ram: 512 GB              │ │
│  │    s3://data/v3   │  │  - cloud: AWS p4d.24xlarge  │ │
│  │  - row_count:     │  │  - duration: 4h 23m         │ │
│  │    2,100,000      │  │  - cost_estimate: $48.50    │ │
│  │  - random_seed:   │  │                             │ │
│  │    42             │  │                             │ │
│  └─────────────────┘  └──────────────────────────────┘ │
└───────────────────────────────────────────────────────┘
```

### 13.2 Lineage Graph

```
Dataset v2          Dataset v3 (cleaned)
    │                    │
    ▼                    ▼
Experiment: fraud-v2  Experiment: fraud-v3
  Run-001 (0.94 acc)    Run-042 (0.968 f1) ★
  Run-002 (0.91 acc)    Run-043 (0.972 f1) ★★
                             │
                             ▼
                     Model Registry
                     "fraud-detector"
                       v1 → ARCHIVED
                       v2 → PRODUCTION ◄── from Run-043
                       v3 → STAGING    ◄── from Run-078
```

---

## 14. Multi-Tenancy & Access Control

### 14.1 Tenant Isolation

```
┌────────────────────────────────────────────────────────┐
│                   MULTI-TENANCY                         │
│                                                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Organization: AcmeCorp                           │  │
│  │  ┌─────────────────┐  ┌─────────────────────┐    │  │
│  │  │ Team: ML-Fraud  │  │ Team: ML-RecSys     │    │  │
│  │  │                 │  │                      │    │  │
│  │  │ Project: fraud  │  │ Project: recommend   │    │  │
│  │  │  - Experiment A │  │  - Experiment X      │    │  │
│  │  │  - Experiment B │  │  - Experiment Y      │    │  │
│  │  │                 │  │                      │    │  │
│  │  │ Members:        │  │ Members:             │    │  │
│  │  │  alice (admin)  │  │  bob (admin)         │    │  │
│  │  │  carol (contrib)│  │  dave (viewer)       │    │  │
│  │  └─────────────────┘  └─────────────────────┘    │  │
│  └──────────────────────────────────────────────────┘  │
│                                                        │
│  Data isolation: all queries include org_id filter      │
│  Storage isolation: S3 prefix per org                   │
│  Network isolation: VPC peering for enterprise          │
└────────────────────────────────────────────────────────┘
```

### 14.2 RBAC Matrix

| Permission | Viewer | Contributor | Admin | Org Owner |
|------------|--------|-------------|-------|-----------|
| View runs & metrics | ✓ | ✓ | ✓ | ✓ |
| Create runs & log data | ✗ | ✓ | ✓ | ✓ |
| Delete runs | ✗ | Own only | ✓ | ✓ |
| Manage experiments | ✗ | ✗ | ✓ | ✓ |
| Model registry (promote) | ✗ | ✗ | ✓ | ✓ |
| Manage team members | ✗ | ✗ | ✓ | ✓ |
| Billing & org settings | ✗ | ✗ | ✗ | ✓ |
| API key management | ✗ | Own only | ✓ | ✓ |

---

## 15. Scalability Strategy

| Concern | Solution |
|---------|----------|
| Metric ingestion (100K QPS) | Kafka buffering → batched writes to TimescaleDB; horizontal scaling of ingestion service |
| Time-series query performance | TimescaleDB hypertables (auto-partition by time); continuous aggregates for summary queries; downsampling for old data |
| Large artifact uploads | Presigned URLs → direct to S3; multipart upload for files > 100MB |
| Run search at scale | Elasticsearch cluster with sharding by org_id; denormalized run documents |
| Dashboard live updates | Redis Pub/Sub → WebSocket gateway; fan-out per active dashboard |
| Database growth | PostgreSQL partitioned by month; TimescaleDB compression (10× for data > 7 days); S3 lifecycle policies (Glacier for artifacts > 1 year) |
| Multi-region | Active-passive for metadata DBs; S3 cross-region replication for artifacts; CDN for dashboard |
| Burst traffic | Auto-scaling ingestion service pods; Kafka absorbs spikes; back-pressure to SDK (retry with backoff) |

---

## 16. Caching Strategy

```
┌──────────────┐     ┌──────────────┐     ┌───────────────┐     ┌──────────────┐
│  Browser     │────►│    CDN       │────►│  Redis L1     │────►│  Database    │
│  Cache       │     │ (Dashboard   │     │  (Hot Data)   │     │  (Source of  │
│              │     │  assets)     │     │               │     │   Truth)     │
│ API response │     │              │     │ - live metrics│     │              │
│ (30s TTL)    │     │ JS/CSS/fonts │     │ - run status  │     │ PostgreSQL   │
│              │     │ (7 day TTL)  │     │ - search cache│     │ TimescaleDB  │
└──────────────┘     └──────────────┘     │ (60s TTL)     │     │ Elasticsearch│
                                          └───────────────┘     └──────────────┘
```

| Layer | What's Cached | TTL | Invalidation |
|-------|---------------|-----|-------------|
| Browser | API responses (run list, comparison) | 30s | On user action (refresh, new run) |
| CDN | Dashboard static assets (JS, CSS, images) | 7 days | Cache-busting with content hash |
| Redis L1 | Active run status + latest metrics | Real-time | Overwritten on each metric log |
| Redis L2 | Search results for popular queries | 60s | Invalidated on run completion |
| Redis L3 | Metric summary aggregates | 5 min | Recomputed by aggregation service |

---

## 17. Handling Failure & Durability

| Failure Scenario | Mitigation |
|------------------|------------|
| SDK cannot reach ingestion service | SDK buffers to local disk; retries with exponential backoff; uploads buffered data on reconnect |
| Kafka broker down | Multi-broker cluster (3+ replicas); producer retries; SDK falls back to synchronous PostgreSQL write |
| TimescaleDB overloaded | Kafka acts as buffer (can absorb hours of backlog); Metric Writer pauses and resumes; no data loss |
| S3 upload fails mid-transfer | SDK retries multipart upload (resumable); temporary presigned URLs regenerated |
| PostgreSQL failover | Primary-replica with automatic failover (Patroni/RDS Multi-AZ); <30s recovery |
| Elasticsearch index lag | Dashboard falls back to PostgreSQL for search (slower but correct); background re-sync |
| WebSocket gateway crash | Client auto-reconnects; missed metrics fetched via REST on reconnect |
| Training job crashes | Run marked FAILED after heartbeat timeout (5 min); all metrics up to crash are preserved |

**Data durability guarantees:**
- Kafka: `acks=all`, `min.insync.replicas=2` — no acknowledged metric is lost
- PostgreSQL: synchronous replication for model registry; async for run metadata
- S3: 99.999999999% (11 nines) durability
- TimescaleDB: WAL replication + daily backups to S3

---

## 18. API Design

### 18.1 Run Management

**POST `/api/v2/runs`** — Create a new run
```json
// Request
{
    "experiment_id": "exp-uuid-123",
    "tags": [
        { "key": "model_arch", "value": "transformer" },
        { "key": "team", "value": "fraud" }
    ],
    "environment": {
        "python_version": "3.11.7",
        "git_commit": "a1b2c3d"
    }
}

// Response (201 Created)
{
    "run_id": "run-uuid-456",
    "experiment_id": "exp-uuid-123",
    "status": "RUNNING",
    "start_time": "2026-03-26T10:00:00Z",
    "artifact_uri": "s3://ml-platform/acme/fraud/exp-123/run-456/"
}
```

### 18.2 Batch Logging

**POST `/api/v2/runs/{run_id}/log-batch`** — Log params, metrics, and tags in one call
```json
// Request
{
    "params": [
        { "key": "learning_rate", "value": "0.001" },
        { "key": "batch_size", "value": "64" }
    ],
    "metrics": [
        { "key": "loss", "value": 0.032, "step": 5000, "timestamp": 1711440000 },
        { "key": "accuracy", "value": 0.974, "step": 5000, "timestamp": 1711440000 }
    ],
    "tags": [
        { "key": "phase", "value": "training" }
    ]
}

// Response (200 OK)
{
    "status": "ok",
    "logged": { "params": 2, "metrics": 2, "tags": 1 }
}
```

### 18.3 Run Search

**POST `/api/v2/runs/search`** — Search runs with filter expression
```json
// Request
{
    "experiment_ids": ["exp-uuid-123"],
    "filter": "metrics.accuracy > 0.95 AND params.optimizer = 'adam'",
    "order_by": ["metrics.f1_score DESC"],
    "max_results": 20,
    "page_token": "..."
}

// Response (200 OK)
{
    "runs": [
        {
            "run_id": "run-uuid-456",
            "status": "COMPLETED",
            "start_time": "2026-03-26T10:00:00Z",
            "end_time": "2026-03-26T14:23:00Z",
            "params": { "learning_rate": "0.001", "optimizer": "adam" },
            "metrics": {
                "accuracy": { "last": 0.974, "min": 0.51, "max": 0.974 },
                "f1_score": { "last": 0.968, "min": 0.48, "max": 0.968 }
            },
            "tags": { "model_arch": "transformer" }
        }
    ],
    "next_page_token": "..."
}
```

### 18.4 Artifact Upload

**GET `/api/v2/runs/{run_id}/artifacts/presigned-upload`** — Get presigned URL
```json
// Request query params: ?path=model/model.pth&size_bytes=241172480

// Response (200 OK)
{
    "upload_url": "https://s3.amazonaws.com/ml-platform/...?X-Amz-Signature=...",
    "method": "PUT",
    "expires_in_seconds": 900,
    "headers": {
        "Content-Type": "application/octet-stream"
    }
}
```

### 18.5 Model Registry

**POST `/api/v2/model-registry/{model_name}/versions`** — Register a model version
```json
// Request
{
    "run_id": "run-uuid-456",
    "artifact_path": "model/model.pth",
    "description": "Transformer fraud detector, trained on fraud_v3 dataset"
}

// Response (201 Created)
{
    "model_name": "fraud-detector",
    "version": 3,
    "stage": "NONE",
    "status": "READY",
    "run_id": "run-uuid-456",
    "created_at": "2026-03-26T14:30:00Z"
}
```

**PATCH `/api/v2/model-registry/{model_name}/versions/{version}/stage`** — Promote model
```json
// Request
{
    "stage": "PRODUCTION",
    "archive_existing_production": true
}

// Response (200 OK)
{
    "model_name": "fraud-detector",
    "version": 3,
    "stage": "PRODUCTION",
    "previous_production_version": 2,
    "transitioned_at": "2026-03-26T15:00:00Z"
}
```

---

## 19. Tech Stack Summary

| Layer | Technology |
|-------|------------|
| Python SDK | Python (asyncio + threading), requests, boto3 |
| CLI | Python (Click/Typer) |
| Web Dashboard | React.js / Next.js, Plotly.js / D3.js, WebSocket |
| API Gateway | Kong / Envoy / AWS API Gateway |
| Ingestion Service | Go / Java (high throughput, low GC pause) |
| Query Service | Python (FastAPI) / Java (Spring Boot) |
| Artifact Service | Go (efficient presigned URL generation + S3 proxy) |
| Model Registry Service | Python (FastAPI) / Java |
| Alert Engine | Python (Faust / Flink) |
| Primary Database | PostgreSQL 16 (with partitioning) |
| Time-Series Database | TimescaleDB (PostgreSQL extension) or InfluxDB |
| Search Engine | Elasticsearch 8.x |
| Cache | Redis Cluster |
| Message Queue | Apache Kafka (with Schema Registry) |
| Object Storage | AWS S3 / GCS / MinIO (self-hosted) |
| CDC Pipeline | Debezium → Kafka → Elasticsearch |
| Container Orchestration | Kubernetes (EKS / GKE) |
| CI/CD | GitHub Actions + ArgoCD |
| Monitoring | Prometheus + Grafana + Jaeger |
| Auth | Keycloak / Auth0 (OAuth 2.0 + OIDC) |

---

## 20. Deployment Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Kubernetes Cluster                            │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ Ingestion Svc│  │ Query Svc    │  │ Artifact Svc         │  │
│  │ (10 pods,    │  │ (5 pods,     │  │ (3 pods)             │  │
│  │  HPA on QPS) │  │  HPA on CPU) │  │                      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ Model Reg.   │  │ Alert Engine │  │ WebSocket Gateway    │  │
│  │ Svc (3 pods) │  │ (2 pods)     │  │ (5 pods, sticky)     │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │ Metric Writer│  │ ES Sync      │  │ Auth Service         │  │
│  │ (N pods =    │  │ Consumer     │  │ (2 pods)             │  │
│  │  Kafka parts)│  │ (3 pods)     │  │                      │  │
│  └──────────────┘  └──────────────┘  └──────────────────────┘  │
│                                                                 │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                  STATEFUL SERVICES                         │ │
│  │  PostgreSQL (Primary + 2 Read Replicas)                    │ │
│  │  TimescaleDB (Primary + 1 Replica)                         │ │
│  │  Elasticsearch (3-node cluster)                            │ │
│  │  Redis Cluster (6 nodes, 3 primary + 3 replica)            │ │
│  │  Kafka (3 brokers + ZooKeeper/KRaft)                       │ │
│  └────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘

CI/CD: GitHub → Actions → Docker Build → ECR → ArgoCD → K8s
Environments: Dev → Staging → Production (Blue-Green deployments)
```

---

## 21. Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Async metric logging | Kafka buffer + batch writer | 100K QPS would overwhelm direct DB writes; Kafka absorbs bursts |
| Time-series DB for metrics | TimescaleDB over plain PostgreSQL | 10× compression, hypertable auto-partitioning, native downsampling, purpose-built for time-series queries |
| Presigned URLs for artifacts | Direct client-to-S3 upload | API servers never touch large files; scales without bandwidth bottleneck |
| Denormalized search index | Elasticsearch + Debezium CDC | Complex multi-field queries across params+metrics; eventual consistency is acceptable for search |
| Non-blocking SDK | Background threads + local disk fallback | Training loop must never block or fail due to platform issues |
| Model registry in PostgreSQL | Strong consistency (SERIALIZABLE transactions) | Stage transitions (staging→production) must be atomic; no split-brain |
| WebSocket for live dashboards | Redis Pub/Sub → WS Gateway | Sub-second metric updates for active training monitoring |
| Multi-tenancy via row-level filtering | Shared infrastructure, org_id on every table | Cost-effective; simpler ops than per-tenant DBs; RLS in PostgreSQL for enforcement |

---

## 22. Key Interview Discussion Points

| Question | Answer |
|----------|--------|
| **Why not write metrics directly to PostgreSQL?** | At 100K QPS with step-level granularity, PostgreSQL would be overwhelmed. Kafka buffers and batch-writes to a time-series DB handle this at 10× the throughput with better compression. |
| **How do you ensure training never blocks?** | SDK uses background threads, in-memory buffer, and local disk fallback. Log calls return immediately. Even if the platform is completely down, the training job continues unaffected. |
| **How do you handle a 10 GB model artifact?** | Presigned URLs for direct-to-S3 multipart upload. API server only handles metadata (~200 bytes). Client uploads at full bandwidth, with resumable chunks. |
| **How do you keep the search index in sync?** | Debezium captures PostgreSQL WAL changes, publishes to Kafka, and an ES sync consumer updates Elasticsearch. Typical lag is ~2 seconds. |
| **CAP trade-off?** | AP for dashboards and search (eventual consistency is fine — stale by seconds). CP for model registry stage transitions (SERIALIZABLE isolation — no two models in PRODUCTION simultaneously). |
| **How do you scale metric storage long-term?** | TimescaleDB compression (10× for data > 7 days), continuous aggregates for summary queries, and tiered retention: raw data for 90 days, downsampled (1-min averages) for 1 year, then archive to S3. |
| **How does run comparison work for 100 runs?** | Parallel fan-out queries to TimescaleDB (per-run metric curves) + PostgreSQL (params). Downsampled to 2K points per metric per run. Response assembled in ~500ms. |
| **Why Elasticsearch instead of PostgreSQL for search?** | Runs have 50+ dynamic param/metric keys. Elasticsearch handles dynamic mappings, nested field queries, and full-text search far better than relational JOINs across param/metric tables. |
| **How do you enforce multi-tenancy isolation?** | Every table has `org_id`. PostgreSQL Row-Level Security (RLS) policies enforce isolation at the DB level. S3 paths are prefixed by org. API layer validates org membership on every request. |
| **Self-hosted vs cloud consideration?** | MinIO replaces S3, TimescaleDB runs on-prem, Kafka cluster self-managed. Platform is cloud-agnostic by design — all components have open-source equivalents. |

---

## 23. Summary — System at a Glance

```
         ╔═══════════════════════════════════════════════════════════╗
         ║     ML EXPERIMENT TRACKING & ANALYSIS PLATFORM            ║
         ╠═══════════════════════════════════════════════════════════╣
         ║                                                           ║
         ║   WRITE PATH (async, non-blocking):                       ║
         ║   SDK → Ingestion Service → Kafka → Batch Writer          ║
         ║                                     → TimescaleDB         ║
         ║   Artifacts: SDK → Presigned URL → S3 (direct)            ║
         ║                                                           ║
         ║   READ PATH (cached, optimized):                          ║
         ║   Dashboard → Query Service → Redis / TimescaleDB / ES    ║
         ║   Live updates: Redis Pub/Sub → WebSocket → Browser       ║
         ║                                                           ║
         ║   MODEL LIFECYCLE:                                        ║
         ║   Run → Registry → NONE → STAGING → PRODUCTION            ║
         ║   (strong consistency, RBAC-gated transitions)             ║
         ║                                                           ║
         ║   KEY INSIGHTS:                                           ║
         ║   1. Never block the training loop (async SDK + fallback)  ║
         ║   2. Separate write path (Kafka) from read path (cache)    ║
         ║   3. Presigned URLs for artifacts (bypass API servers)     ║
         ║   4. Time-series DB for metrics (compression + speed)      ║
         ║   5. Denormalized search for complex run queries           ║
         ║                                                           ║
         ╚═══════════════════════════════════════════════════════════╝
```

---
