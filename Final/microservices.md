



## Observability

### Observability in Spring Boot Microservices (Simple Version | 8 Minutes)

#### What Is Observability, and Why It Matters (45 seconds)

Observability means being able to understand **what's actually happening inside your system** from the outside — without needing to guess or add new code every time something breaks. In a monolith, if something goes wrong, you check one app's logs and you're mostly done. In microservices, a single user request might pass through 6+ services — so when something breaks, "which service caused it, and why?" becomes a real challenge without the right tools in place.

**Good line to open with:** "Observability isn't the same as monitoring. Monitoring tells you *something* is wrong (like a dashboard turning red). Observability lets you actually figure out *why* it's wrong, by digging into the data — even for problems you didn't specifically anticipate ahead of time."

#### The Three Pillars of Observability (this is the core structure — 30 seconds)

1. **Logs** — detailed records of individual events ("what happened, in detail, at this exact moment").
2. **Metrics** — numeric measurements over time ("how many requests, how slow, how often failing").
3. **Traces** — the journey of a single request as it flows across multiple services.

Together, these answer three different questions: logs tell you the details of one specific thing, metrics tell you the overall health/trends, and traces tell you the path a request took across your whole system.

### Pillar 1: Logging (1.5 minutes)

#### The Problem With Logging in Microservices

In a monolith, logs live in one place. In microservices, logs are scattered across dozens of services, each with their own log files. If Order Service calls Payment Service and something fails, you need to manually go find and match up log lines across two separate services — painful without the right setup.

#### The Fix: Centralized Logging + Correlation IDs

- **Centralized logging** — all services ship their logs to one central place, typically the **ELK Stack** (Elasticsearch, Logstash, Kibana) or a managed alternative (Datadog, Splunk). Now you search/filter logs from all services in one dashboard, instead of SSH-ing into individual machines.
- **Correlation ID** — a unique ID generated at the very start of a request (usually at the API Gateway), then passed along in headers to every downstream service call. Every log line, across every service, includes this same ID — so you can search for one ID and instantly see the full story across all services involved in that request.

java

```java
// Example: adding a correlation ID via a filter
@Component
public class CorrelationIdFilter extends OncePerRequestFilter {
    protected void doFilterInternal(HttpServletRequest req, HttpServletResponse res, FilterChain chain) {
        String correlationId = req.getHeader("X-Correlation-Id");
        if (correlationId == null) {
            correlationId = UUID.randomUUID().toString();
        }
        MDC.put("correlationId", correlationId); // added to every log line automatically
        res.setHeader("X-Correlation-Id", correlationId);
        chain.doFilter(req, res);
    }
}
```

`MDC` (Mapped Diagnostic Context) is what lets the correlation ID automatically show up in every log statement without manually adding it to each `log.info(...)` call.

**Important production point:** logs should be **structured** (JSON format), not plain text — this makes them searchable/filterable by field (service name, correlation ID, log level) instead of relying on fragile text pattern matching.

### Pillar 2: Metrics (2 minutes)

#### What Metrics Actually Track

Numeric data collected over time — request counts, response times, error rates, CPU/memory usage, queue sizes, thread pool usage, and business-level numbers too (orders placed per minute, payment failures per hour).

#### Spring Boot's Tooling: Micrometer + Actuator + Prometheus + Grafana

- **Spring Boot Actuator** — exposes built-in operational endpoints like `/actuator/health`, `/actuator/metrics`, `/actuator/prometheus` out of the box.
- **Micrometer** — the metrics-collection library baked into Spring Boot; think of it as "SLF4J but for metrics" — one common API, multiple possible backends.
- **Prometheus** — a time-series database that periodically "scrapes" (pulls) metrics from each service's `/actuator/prometheus` endpoint and stores them over time.
- **Grafana** — the dashboard/visualization layer on top of Prometheus's data — the actual graphs and dashboards people look at.

yaml

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health, metrics, prometheus
  metrics:
    tags:
      application: order-service
```

java

```java
// Custom business metric example
Counter orderCounter = Counter.builder("orders.placed")
    .description("Total orders placed")
    .register(meterRegistry);

orderCounter.increment();
```

#### The Four Golden Signals (worth naming explicitly — a well-known framework)

A widely referenced set of metrics to always track for any service (from Google's SRE practices):

1. **Latency** — how long requests take (especially P95/P99, not just average — averages hide bad outliers).
2. **Traffic** — how many requests are coming in.
3. **Errors** — rate of failed requests.
4. **Saturation** — how "full" your resources are (CPU, memory, thread pool usage, queue depth).

**Good line for an interview:** "I'd always track P95/P99 latency, not just average response time — because averages can look perfectly healthy while a meaningful chunk of real users are having a genuinely bad, slow experience."

### Pillar 3: Distributed Tracing (2 minutes)

#### The Problem It Solves

A single user request might flow: API Gateway → Order Service → Payment Service → Inventory Service → back to Order Service. If it's slow or fails, which specific hop caused it? Logs alone don't easily show you the full timeline across services — tracing does.

#### How It Works

Each request gets a **trace ID** (similar idea to the correlation ID, and often the same thing in practice) when it enters the system. As it moves through each service, every step is recorded as a **span** — a single unit of work with a start time, end time, and metadata. All spans sharing the same trace ID get stitched together into one complete timeline, visualized as a **waterfall diagram** — showing exactly how long each service/step took, in order.

#### Tools

- **Zipkin** — simpler, lightweight, easier to set up.
- **Jaeger** — more feature-rich, built for larger-scale systems, originally from Uber.
- **Spring Cloud Sleuth** (older) / **Micrometer Tracing** (current, modern replacement) — automatically generates trace IDs and spans for Spring Boot apps, and propagates them across service calls without much manual code.
- **OpenTelemetry** — the newer, vendor-neutral industry standard for collecting traces (and metrics/logs too) — worth naming as the direction the industry is moving toward, so you're not locked into one specific vendor's format.

yaml

```yaml
management:
  tracing:
    sampling:
      probability: 1.0  # trace 100% of requests (lower this in high-traffic prod, e.g., 0.1)
  zipkin:
    tracing:
      endpoint: http://zipkin-server:9411/api/v2/spans
```

With Micrometer Tracing + Sleuth-style auto-instrumentation, trace IDs get added to your logs automatically too — this is exactly what ties tracing and logging together (search logs by trace ID, then view the full trace in Zipkin/Jaeger for the visual timeline).

**Important production point:** tracing **every single request** (100% sampling) can be expensive at high traffic volumes — most production systems sample a percentage (e.g., 10%) instead, while still tracing 100% of errors/slow requests specifically (this is called "tail-based sampling" in more advanced setups).

#### How the Three Pillars Work Together (1 minute)

A realistic debugging flow, good to describe in an interview:

1. **Metrics** show you something's wrong — Grafana dashboard shows Payment Service's error rate spiked.
2. **Traces** show you where in the request chain it's happening — you see Payment Service calls to a bank API are timing out.
3. **Logs** show you the specific detail — you search by trace ID and find the exact error message/stack trace from that failing call.

None of the three pillars alone gives you the complete picture — metrics tell you *something's* wrong, traces tell you *where*, logs tell you *exactly what*.

#### Additional Observability Practices Worth Mentioning (1 minute)

- **Health checks** (`/actuator/health`) — not just "is the app running," but real dependency checks (can I reach my database? my message broker?) — Kubernetes uses this for liveness/readiness probes.
- **Alerting** — metrics/logs are useless if nobody's watching; set up alerts (via Prometheus Alertmanager, PagerDuty, etc.) on key thresholds — error rate spikes, latency spikes, circuit breakers tripping — so humans get notified before customers start complaining.
- **Dashboards per service, plus a system-wide overview** — individual service dashboards for deep dives, plus a higher-level dashboard showing overall system health across all services.
- **Business metrics, not just technical ones** — tracking "orders placed per minute" or "payment success rate" alongside CPU/memory gives the business side visibility too, not just engineering.

#### Trade-offs / Challenges Worth Mentioning (45 seconds)

- **Observability tooling has real infrastructure cost** — running Prometheus, Grafana, ELK, and Zipkin/Jaeger is itself operational overhead; many teams use managed alternatives (Datadog, New Relic, Grafana Cloud) to offload this.
- **Too much logging can be as bad as too little** — excessive, noisy logging increases storage cost and makes it harder to find what actually matters; log at the right level (INFO for normal flow, ERROR for real problems, DEBUG for detailed troubleshooting, not left on in production by default).
- **100% trace sampling doesn't scale at high volume** — needs a deliberate sampling strategy, not an all-or-nothing decision.
- **Correlation IDs must be propagated consistently everywhere** — if even one service forgets to pass the header along to its own downstream calls, the trace/log chain breaks at that point, and you lose the full picture right where you need it most.

#### Simple Closing Line

> "In practice, I've relied on all three pillars together — Micrometer and Prometheus/Grafana for metrics, ELK for centralized structured logs with correlation IDs, and Zipkin or Jaeger for distributed tracing across services. The real value comes from combining them — metrics tell me something's wrong, tracing shows me where in the chain, and logs give me the exact detail, all tied together by the same correlation/trace ID."





## Walk me through how you've set up logging, metrics, and tracing (Actuator, Micrometer, distributed tracing) for a Spring Boot microservice — and a time it helped you catch something.

### Interview Answer (4–5 Minutes)

**In production, I believe observability has three pillars: Logging, Metrics, and Distributed Tracing. Logs tell me *what happened*, metrics tell me *how the system is behaving over time*, and distributed traces tell me *where the time is spent across services*. In our Spring Boot microservices, we implemented all three using Spring Boot Actuator, Micrometer, Prometheus, Grafana, OpenTelemetry, and centralized logging.**

---

# 1. Logging

We used **structured logging** with Logback and centralized logs using the **ELK stack (Elasticsearch, Logstash, Kibana)**.

Every log included:

- Timestamp
- Log level
- Service name
- Request ID / Trace ID
- User ID (where appropriate)
- Exception details

Example:

```
INFO
traceId=abc123
requestId=req-456
User authorization completed in 42ms
```

The trace ID allowed us to correlate logs across multiple microservices.

We also followed good logging practices:

-  INFO → Business events 
-  WARN → Recoverable issues 
-  ERROR → Exceptions 
-  DEBUG → Development troubleshooting only 

---

# 2. Metrics

We exposed application metrics using **Spring Boot Actuator**.

Micrometer automatically collected JVM and application metrics.

We exported them to **Prometheus**, and Grafana dashboards visualized them.

The metrics we monitored most frequently were:

### JVM

-  Heap usage 
-  GC pause time 
-  Thread count 
-  CPU utilization 

### Application

-  Request rate (RPS) 
-  P95/P99 latency 
-  HTTP error rate 
-  Active requests 
-  Tomcat thread utilization 

### Database

-  HikariCP active/idle connections 
-  Connection wait time 
-  Slow query count 

### Kafka

-  Consumer lag 
-  Processing rate 
-  Retry count 

These dashboards gave us a real-time view of system health.

---

# 3. Distributed Tracing

For request tracing across services, we used **OpenTelemetry** with **Jaeger**.

Each incoming request received a Trace ID.

Example flow:

```

```

```
Client
   ↓
API Gateway
   ↓
Authorization Service
   ↓
IAM Service
   ↓
PostgreSQL
```

Using a single trace, we could see exactly how much time each component consumed.

Instead of guessing where latency originated, the trace showed whether the bottleneck was:

-  Business logic 
-  Database 
-  External API 
-  Kafka 
-  Network 

---

# Production Incident

One incident where observability was extremely valuable involved an authorization API that suddenly slowed down.

Users reported response times increasing from around **100 ms** to nearly **2 seconds**.

### Step 1: Metrics

Grafana showed:

-  CPU usage was normal. 
-  Memory usage was stable. 
-  Error rate was low. 
-  P95 latency had increased significantly. 

This suggested the application wasn't overloaded.

---

### Step 2: Distributed Trace

Using Jaeger, we examined a slow request.

The trace showed:

```

```

```
Authentication      20 ms

Authorization Logic 40 ms

Database Query    1800 ms
```

Immediately we knew the delay wasn't in Spring Boot—it was the database query.

---

### Step 3: Logs

The logs confirmed that the same repository method was executing repeatedly.

Enabling Hibernate SQL logging revealed an **N+1 query problem** introduced after a recent code change.

---

### Resolution

We:

-  Optimized the query using `JOIN FETCH` 
-  Added an index where needed 
-  Reduced unnecessary database calls 
-  Added Redis caching for frequently accessed authorization data 

After deployment:

-  Response time dropped back below **100 ms** 
-  Database CPU utilization decreased 
-  The Grafana latency graph returned to normal 

Without metrics and tracing, identifying the bottleneck would have taken much longer.

---

# Why All Three Matter

I think of observability like this:


| Tool        | Answers                  |
| ----------- | ------------------------ |
| **Logs**    | What happened?           |
| **Metrics** | Is the system healthy?   |
| **Tracing** | Where is the time spent? |


Using only one of them leaves blind spots. Together, they significantly reduce the time needed to diagnose production issues.

---

# Interview Closing

> **In our Spring Boot microservices, we implemented structured logging with centralized log aggregation, Actuator and Micrometer for metrics exported to Prometheus and visualized in Grafana, and OpenTelemetry-based distributed tracing with Jaeger. One production issue involved a sudden increase in API latency. Metrics showed the application itself was healthy, tracing pinpointed a slow database call, and logs confirmed an N+1 query introduced by a recent change. We optimized the query and added caching, restoring latency to normal. That experience reinforced the importance of combining logs, metrics, and traces to diagnose production problems quickly and accurately.**

