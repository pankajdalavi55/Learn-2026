# Phase 3: Scalability Patterns

**Navigation:** [← Previous: Distributed Systems](02-distributed-systems.md) | [Next: Case Studies (Mandatory) →](04-case-studies-mandatory.md)

---

This phase covers architectural patterns for building resilient, scalable systems. These patterns are critical for Staff+ level interviews where interviewers expect you to discuss system evolution, organizational boundaries, and operational excellence.

---

# Section 1: Service Architecture Evolution

## 1.1 Monolith → Modular Monolith → Microservices

### Concept Overview (What & Why)

**Monolith:**
- Single deployable unit
- Shared database
- All code in one repository
- Simple to develop and deploy initially

**Modular Monolith:**
- Single deployable unit
- Internal module boundaries (packages/namespaces)
- Modules communicate through defined interfaces
- Can be split later if needed

**Microservices:**
- Independent deployable services
- Each service owns its data
- Network communication (REST, gRPC)
- Independent technology choices

**Why This Matters:**
- Every large system starts as a monolith
- Premature microservices is a common mistake
- The evolution path matters more than the end state

### Key Design Principles

**When to Choose:**

| Stage | Architecture | Characteristics |
|-------|--------------|-----------------|
| Early startup | Monolith | <10 engineers, finding product-market fit |
| Growing startup | Modular Monolith | 10-50 engineers, clear domain boundaries |
| Scale-up | Microservices (selective) | 50+ engineers, independent team scaling |
| Enterprise | Full microservices | 100+ engineers, organizational boundaries |

**The Monolith Trap:**
- Too often: "We're building for scale, so microservices from day one"
- Reality: Microservices add complexity before they add value
- Better: Start monolith, extract services as needed

### Trade-offs & Decision Matrix

| Factor | Monolith | Modular Monolith | Microservices |
|--------|----------|------------------|---------------|
| Deployment | Simple | Simple | Complex |
| Debugging | Easy | Easy | Hard (distributed tracing) |
| Data consistency | ACID | ACID | Eventual |
| Team autonomy | Low | Medium | High |
| Operational overhead | Low | Low | High |
| Technology diversity | Single stack | Single stack | Polyglot |
| Scaling | All or nothing | All or nothing | Granular |

### Real-World Examples

**Successful Monolith:**
- Shopify ran as a monolith with 2,000+ developers
- Carefully modularized with strict module boundaries
- Only split services for specific scaling needs

**Microservices Done Right:**
- Amazon: Services aligned to teams (two-pizza rule)
- Netflix: Started monolith, migrated during cloud transition
- Uber: Thousands of services, but started with monolith

**Microservices Done Wrong:**
- Premature extraction before understanding boundaries
- Too fine-grained services (one per entity)
- Distributed monolith (services coupled, deploy together)

### Interview Perspective

**What interviewers look for:**
- Pragmatic thinking about architecture evolution
- Understanding of organizational factors
- Awareness of complexity costs

**Common traps:**
- ❌ "We need microservices for scalability" (modular monolith scales too)
- ❌ "Microservices from day one" (premature complexity)
- ❌ "One service per entity" (too granular)

**Strong signals:**
- ✅ "Start with modular monolith, extract when needed"
- ✅ "Microservices boundaries should align with team boundaries"
- ✅ "Each service owns its data; no shared databases"

**Follow-up questions:**
- "How do you decide when to extract a service?"
- "How do you handle a service calling 10 other services?"
- "What about shared libraries vs duplicate code?"

---

## 1.2 Service Decomposition Strategies

### Concept Overview (What & Why)

**Service Decomposition:** Breaking a system into services with clear boundaries.

**Bad Decomposition Signs:**
- Services that must deploy together
- Circular dependencies
- Shared mutable state
- Change in one service cascades to others

**Good Decomposition Signs:**
- Services can be developed independently
- Services have stable APIs
- Most changes are local to one service
- Teams can release independently

### Key Design Principles

**Decomposition Strategies:**

| Strategy | Description | Example |
|----------|-------------|---------|
| By Business Capability | Align to what business does | Orders, Payments, Shipping |
| By Subdomain (DDD) | Bounded contexts from Domain-Driven Design | Customer, Catalog, Fulfillment |
| By Team | One team, one or few services | Team owns the whole slice |
| By Data | Who owns this data? | User Service owns users table |

**Domain-Driven Design (DDD) Concepts:**

**Bounded Context:**
- A boundary within which a term has specific meaning
- "Customer" in Sales vs "Customer" in Shipping = different models
- Each bounded context can be a service

**Aggregate:**
- Cluster of domain objects treated as a unit
- Transaction boundary
- Good candidate for service boundary

### Trade-offs & Decision Matrix

| Decomposition | Pros | Cons |
|---------------|------|------|
| By capability | Business-aligned, intuitive | May not match data ownership |
| By subdomain | Clear boundaries, DDD benefits | Requires domain expertise |
| By team | Organizational alignment | Services may be arbitrary |
| By data | Clear ownership | May not match use cases |

**Rules of Thumb:**
- Service should be ownable by one team (5-10 people)
- Service should be deployable independently
- Service should have one reason to change
- When in doubt, keep together (can split later)

### Interview Perspective

**Strong signals:**
- ✅ "Bounded context from DDD gives natural service boundaries"
- ✅ "Services aligned with teams and business capabilities"
- ✅ "Each service owns its data, no shared databases"

**Common trap:**
- ❌ Decomposing by technical layer (API service, DB service)

---

## 1.3 Database per Service

### Concept Overview (What & Why)

**Pattern:** Each microservice owns its data and database. No direct database sharing.

```
Order Service → Order DB
User Service → User DB
Payment Service → Payment DB

Order Service needs user data?
  → Call User Service API (not User DB)
```

**Why This Matters:**
- Loose coupling (can change schema without affecting others)
- Independent deployment
- Technology freedom (SQL vs NoSQL per service)
- Clear ownership

**Challenges:**
- No cross-service transactions
- Data consistency is eventual
- Queries across services are complex
- Data duplication may be needed

### Key Design Principles

**Handling Cross-Service Data:**

| Need | Solution |
|------|----------|
| Reference data | Cache local copy, refresh periodically |
| Real-time data | API call (sync) or subscribe to events (async) |
| Reports/analytics | Data warehouse (ETL from all services) |
| Cross-service queries | API composition or CQRS read model |

**Data Ownership Questions:**
- Who creates this data?
- Who is the source of truth?
- Who needs to query it?
- How fresh does it need to be?

### Trade-offs & Decision Matrix

| Challenge | Solution | Trade-off |
|-----------|----------|-----------|
| Cross-service joins | API calls + client-side join | Latency, complexity |
| Consistency | Saga pattern | Eventual, complex |
| Reporting | Data warehouse | Staleness, extra infrastructure |
| Duplicate data | Event-driven sync | Consistency window |

### Interview Perspective

**Strong signals:**
- ✅ "No shared databases; each service owns its data"
- ✅ "Cross-service data via APIs or events, not database"
- ✅ "Data warehouse for cross-cutting analytics"

**Follow-up questions:**
- "How do you handle a report that needs data from 5 services?"
- "What if Service A needs to guarantee consistency with Service B?"

---

### Section 1 Cheat Sheet

```
SERVICE ARCHITECTURE

EVOLUTION PATH:
Monolith → Modular Monolith → Microservices (if needed)

DON'T START WITH MICROSERVICES:
• <50 engineers: Modular monolith likely sufficient
• Microservices add complexity before value
• Distributed transactions, debugging, operations

DECOMPOSITION STRATEGIES:
• By business capability (Orders, Payments)
• By subdomain (DDD bounded contexts)
• By team ownership
• By data ownership

GOOD DECOMPOSITION SIGNS:
• Independent deployment
• Stable APIs
• Changes local to one service
• Team autonomy

DATABASE PER SERVICE:
• No shared databases
• Cross-service data via API/events
• Accept eventual consistency
• Data warehouse for analytics

WHEN TO EXTRACT A SERVICE:
• Different scaling requirements
• Different team ownership
• Different release cadence
• Clear bounded context
```

---

# Section 2: Distributed Transaction Patterns

## 2.1 Saga Pattern

### Concept Overview (What & Why)

**Problem:** Microservices can't use distributed transactions (2PC) at scale.
**Solution:** Saga - a sequence of local transactions with compensating actions.

```
Order Saga:
1. Order Service: Create order (PENDING)
2. Payment Service: Charge customer
3. Inventory Service: Reserve items
4. Order Service: Confirm order (CONFIRMED)

If step 3 fails:
  Compensate step 2: Refund customer
  Compensate step 1: Cancel order
```

**Why This Matters:**
- Distributed transactions (2PC) don't scale
- Business processes span multiple services
- Need to handle partial failures

### Key Design Principles

**Two Saga Types:**

| Type | Coordination | Pros | Cons |
|------|--------------|------|------|
| Choreography | Events, no central coordinator | Loose coupling, simple | Hard to track, complex flows |
| Orchestration | Central saga orchestrator | Clear flow, easier to manage | Single point of coordination |

**Choreography Example:**
```
Order Service → publishes OrderCreated
                    ↓
Payment Service ← subscribes, charges, publishes PaymentCompleted
                    ↓
Inventory Service ← subscribes, reserves, publishes InventoryReserved
                    ↓
Order Service ← subscribes, confirms order
```

**Orchestration Example:**
```
Order Saga Orchestrator:
  1. Call Order Service: CreateOrder
  2. Call Payment Service: ChargeCustomer
  3. Call Inventory Service: ReserveItems
  4. Call Order Service: ConfirmOrder
  
  On failure: Execute compensating transactions
```

### Trade-offs & Decision Matrix

| Factor | Choreography | Orchestration |
|--------|--------------|---------------|
| Coupling | Looser | Tighter (to orchestrator) |
| Visibility | Harder to trace | Easy to see in orchestrator |
| Complexity | In events | In orchestrator |
| Testing | Harder | Easier |
| Best for | Simple flows, 2-3 steps | Complex flows, many steps |

### Real-World Examples

**E-commerce Order:**
1. Reserve inventory
2. Charge payment
3. Create shipment
4. Send confirmation

Compensations:
- Release inventory
- Refund payment
- Cancel shipment

**Travel Booking:**
1. Book flight
2. Book hotel
3. Book car

If hotel booking fails, compensate flight booking.

### Failure Scenarios & Edge Cases

| Scenario | Handling |
|----------|----------|
| Step fails | Execute compensating transactions backwards |
| Compensation fails | Retry compensation, alert if stuck |
| Saga stuck | Timeout + compensate or manual intervention |
| Idempotency | Each step must be idempotent (retries) |

### Interview Perspective

**Strong signals:**
- ✅ "Saga pattern for cross-service transactions"
- ✅ "Each step has a compensating action"
- ✅ "Orchestration for complex flows, choreography for simple"

**Follow-up questions:**
- "What if a compensation fails?"
- "How do you test a saga?"
- "How do you handle a saga that's stuck?"

---

### Section 2 Cheat Sheet

```
SAGA PATTERN

PROBLEM: No distributed transactions across services
SOLUTION: Local transactions + compensating actions

CHOREOGRAPHY:
• Event-driven
• No central coordinator
• Simple flows, loose coupling
• Harder to trace/debug

ORCHESTRATION:
• Central coordinator
• Explicit flow definition
• Complex flows, visibility
• Easier testing

COMPENSATION:
• Each step has reverse action
• Execute on failure
• Must handle compensation failures

IDEMPOTENCY:
• Every step must be idempotent
• Retries are expected
• Use unique transaction IDs

EXAMPLE (Order):
1. CreateOrder → compensation: CancelOrder
2. ChargePayment → compensation: Refund
3. ReserveInventory → compensation: Release
4. Ship → compensation: CancelShipment
```

---

# Section 3: Resilience Patterns

## 3.1 Circuit Breaker

### Concept Overview (What & Why)

**Problem:** When a downstream service is failing, callers keep trying, wasting resources and cascading failures.

**Solution:** Circuit breaker - stop calling a failing service temporarily.

```
States:
CLOSED → Normal operation, calls go through
OPEN → Service failing, calls fail fast
HALF-OPEN → Testing if service recovered

CLOSED: Calls succeed
  ↓ (failure threshold exceeded)
OPEN: Calls fail immediately (don't try)
  ↓ (timeout expires)
HALF-OPEN: Allow one test call
  ↓ (success)
CLOSED (or back to OPEN if test fails)
```

### Key Design Principles

**Configuration:**
- **Failure threshold:** How many failures trigger open state (e.g., 5 failures in 10 seconds)
- **Open timeout:** How long to wait before testing (e.g., 30 seconds)
- **Half-open success threshold:** Successes needed to close (e.g., 3)

**What Counts as Failure:**
- HTTP 5xx responses
- Timeouts
- Connection failures
- NOT 4xx (client errors)

### Real-World Examples

**Netflix Hystrix (now Resilience4j):**
```java
CircuitBreaker circuitBreaker = CircuitBreaker.ofDefaults("paymentService");

Supplier<String> decoratedSupplier = CircuitBreaker
    .decorateSupplier(circuitBreaker, () -> paymentService.charge(amount));

Try.ofSupplier(decoratedSupplier)
    .recover(throwable -> "Fallback response");
```

**Fallback Strategies:**
- Return cached data
- Return default value
- Return degraded experience
- Fail gracefully with error message

### Interview Perspective

**Strong signals:**
- ✅ "Circuit breaker to prevent cascading failures"
- ✅ "Fail fast when downstream is unhealthy"
- ✅ "Fallback for graceful degradation"

---

## 3.2 Retries with Exponential Backoff

### Concept Overview (What & Why)

**Problem:** Transient failures happen; immediate retries may not help.

**Solution:** Retry with increasing delays.

```
Attempt 1: Immediate
Attempt 2: Wait 1 second
Attempt 3: Wait 2 seconds
Attempt 4: Wait 4 seconds
Attempt 5: Wait 8 seconds
Give up after max attempts
```

### Key Design Principles

**Exponential Backoff Formula:**
```
delay = base_delay * (2 ^ attempt) + random_jitter
```

**Jitter (Critical):**
- Without jitter: All clients retry at same time → thundering herd
- With jitter: Retries spread out over time

```
Attempt 1: 1s + random(0, 0.5s)
Attempt 2: 2s + random(0, 1s)
Attempt 3: 4s + random(0, 2s)
```

**What to Retry:**
- Transient errors (503, timeouts)
- NOT permanent errors (400, 404, 401)
- Idempotent operations only (or with idempotency key)

### Trade-offs & Decision Matrix

| Config | Impact |
|--------|--------|
| More retries | Higher success rate, longer latency |
| Shorter delays | Faster recovery, more load on failing service |
| Longer delays | Less load, slower recovery |
| No jitter | Thundering herd |

### Interview Perspective

**Strong signals:**
- ✅ "Exponential backoff with jitter"
- ✅ "Only retry transient errors"
- ✅ "Ensure idempotency for retries"

---

## 3.3 Bulkhead

### Concept Overview (What & Why)

**Problem:** One slow dependency uses up all resources, affecting all requests.

**Solution:** Isolate resources so one failure doesn't affect others.

**Analogy:** Ship compartments - if one floods, others stay dry.

```
Without Bulkhead:
Thread Pool: [A] [A] [A] [B] [B] [B] [C] [C] [C] [C]
If C is slow → all threads blocked → A and B fail too

With Bulkhead:
Pool A: [A] [A] [A]
Pool B: [B] [B] [B]
Pool C: [C] [C] [C] [C]
If C is slow → only Pool C affected
```

### Key Design Principles

**Bulkhead Types:**

| Type | Mechanism | Use Case |
|------|-----------|----------|
| Thread pool isolation | Separate thread pool per dependency | JVM-based services |
| Semaphore isolation | Limit concurrent calls | Lighter weight |
| Connection pool | Separate DB connection pools | Database access |

**Sizing Considerations:**
- Too small: Unnecessary throttling
- Too large: No protection
- Based on: Expected concurrency, dependency latency

### Interview Perspective

**Strong signals:**
- ✅ "Bulkhead to isolate failures"
- ✅ "Separate thread pools for critical dependencies"
- ✅ "Prevents one slow service from affecting others"

---

## 3.4 Timeouts

### Concept Overview (What & Why)

**Problem:** Without timeouts, slow services can hold resources indefinitely.

**Solution:** Set timeouts for all external calls.

**Timeout Types:**
- **Connection timeout:** Time to establish connection (typically 1-5 seconds)
- **Read timeout:** Time to receive response (varies by operation)
- **Total timeout:** End-to-end for the request

### Key Design Principles

**Timeout Guidelines:**

| Operation | Typical Timeout |
|-----------|----------------|
| Database query (simple) | 1-5 seconds |
| Database query (complex) | 30 seconds |
| Cache lookup | 100-500 ms |
| Internal service call | 1-5 seconds |
| External API call | 5-30 seconds |
| File upload | Minutes (depends on size) |

**Timeout Propagation:**
```
Client → API Gateway (5s) → Service A (3s) → Service B (1s)

Each layer should have shorter timeout than its caller.
Otherwise: Caller times out, but downstream continues processing.
```

**Deadline Propagation (gRPC):**
```
Client sets deadline: 5 seconds from now
Propagated to all downstream calls
Each service checks remaining time
```

### Interview Perspective

**Strong signals:**
- ✅ "Timeouts on all external calls"
- ✅ "Timeout budget: Each hop has less than caller"
- ✅ "Combine with circuit breaker for resilience"

---

### Section 3 Cheat Sheet

```
RESILIENCE PATTERNS

CIRCUIT BREAKER:
• States: Closed → Open → Half-Open
• Fail fast when service is down
• Fallback for graceful degradation

RETRY WITH BACKOFF:
delay = base * 2^attempt + jitter
• Jitter prevents thundering herd
• Only retry transient errors
• Require idempotency

BULKHEAD:
• Isolate resources per dependency
• Thread pools, semaphores, connection pools
• One failure doesn't affect others

TIMEOUTS:
• Set on ALL external calls
• Shorter than caller's timeout
• Connection timeout + read timeout

COMBINING PATTERNS:
Request → Timeout → Retry (with backoff) → Circuit Breaker → Fallback

TYPICAL CONFIGURATION:
• Timeout: 3s
• Retries: 3 with exponential backoff
• Circuit: Open after 5 failures
• Circuit: Half-open after 30s
```

---

# Section 4: Observability

## 4.1 The Three Pillars: Metrics, Logging, Tracing

### Concept Overview (What & Why)

**Observability:** Ability to understand system behavior from external outputs.

| Pillar | What | For |
|--------|------|-----|
| Metrics | Numerical measurements over time | Trends, alerts, dashboards |
| Logs | Discrete events with details | Debugging, auditing |
| Traces | Request flow across services | Understanding distributed calls |

**Why This Matters:**
- Microservices are hard to debug without observability
- Production issues need fast diagnosis
- Capacity planning requires metrics

### Key Design Principles

**Metrics (What to Measure):**

| Metric Type | Examples |
|-------------|----------|
| Counters | Request count, error count |
| Gauges | Current connections, queue depth |
| Histograms | Latency distribution, request sizes |
| Summaries | Similar to histograms (different trade-offs) |

**Logging Best Practices:**
- Structured logs (JSON, not plain text)
- Correlation ID across services
- Appropriate log levels (ERROR, WARN, INFO, DEBUG)
- Avoid PII in logs
- Centralized log aggregation (ELK, Splunk)

**Distributed Tracing:**
```
Request enters system:
  TraceID: abc123
  
Service A (SpanID: 1):
  → Calls Service B (SpanID: 2, ParentSpanID: 1)
      → Calls Database (SpanID: 3, ParentSpanID: 2)
  → Calls Service C (SpanID: 4, ParentSpanID: 1)
```

Tools: Jaeger, Zipkin, AWS X-Ray, Datadog APM

---

## 4.2 RED and USE Metrics

### Concept Overview (What & Why)

**RED Method (For Services):**
- **R**ate: Requests per second
- **E**rrors: Failed requests per second
- **D**uration: Latency distribution (p50, p95, p99)

**USE Method (For Resources):**
- **U**tilization: % of resource used (CPU %, memory %)
- **S**aturation: Queue depth (work waiting)
- **E**rrors: Error count

### Key Design Principles

**RED Metrics Example:**
```
Service: Order API
Rate: 1000 req/s
Errors: 5 req/s (0.5%)
Duration: p50=50ms, p95=200ms, p99=500ms
```

**USE Metrics Example:**
```
Resource: Database
Utilization: CPU 70%, Memory 85%
Saturation: Connection queue depth = 50
Errors: Connection timeout = 10/min
```

**When to Use:**
- RED: Request-driven services (APIs)
- USE: Infrastructure resources (CPU, memory, DB, queues)

---

## 4.3 SLIs, SLOs, and SLAs

### Concept Overview (What & Why)

| Term | Definition | Example |
|------|------------|---------|
| SLI (Indicator) | Quantitative measure of service | p99 latency, error rate |
| SLO (Objective) | Target value for an SLI | p99 latency < 200ms |
| SLA (Agreement) | Contract with consequences | 99.9% uptime or refund |

**Relationship:**
```
SLI: What you measure
SLO: What you aim for (internal target, tighter)
SLA: What you promise (external contract, looser)
```

### Key Design Principles

**Good SLIs:**
- User-centric (what users experience)
- Measurable and actionable
- Reflect actual quality

**Common SLIs:**
| Category | SLI |
|----------|-----|
| Availability | % of successful requests |
| Latency | p95 or p99 response time |
| Throughput | Requests per second at full load |
| Error rate | % of failed requests |
| Freshness | Age of data (for async systems) |

**Error Budget:**
```
SLO: 99.9% availability
Budget: 0.1% downtime per month = 43 minutes

If you've used 30 minutes this month:
  → 13 minutes remaining
  → Slow down risky changes
```

### Interview Perspective

**Strong signals:**
- ✅ "SLO of p99 < 200ms, SLA at p99 < 500ms (buffer)"
- ✅ "Error budget approach to balance velocity and reliability"
- ✅ "Alert on SLO burn rate, not just thresholds"

---

### Section 4 Cheat Sheet

```
OBSERVABILITY

THREE PILLARS:
Metrics: Numbers over time (Prometheus)
Logs: Events with details (ELK, Splunk)
Traces: Request flow (Jaeger, Zipkin)

RED (Services):
Rate: Requests per second
Errors: Error rate
Duration: Latency (p50, p95, p99)

USE (Resources):
Utilization: % used
Saturation: Queue depth
Errors: Error count

SLI/SLO/SLA:
SLI: What you measure
SLO: Internal target (stricter)
SLA: External promise (with consequences)

ERROR BUDGET:
99.9% SLO = 43 min downtime/month allowed
Track budget consumption
Slow down if burning too fast

LOGGING BEST PRACTICES:
• Structured (JSON)
• Correlation IDs
• No PII
• Appropriate levels

ALERTING:
• Alert on symptoms, not causes
• SLO burn rate alerts
• Avoid alert fatigue
```

---

# Phase 3 Summary: Building Resilient Systems

These patterns are what separate production-ready systems from prototypes. In Staff+ interviews, you're expected to discuss these naturally.

**Key Takeaways:**

1. **Architecture Evolution:** Start simple, extract services when needed
2. **Service Boundaries:** Align with teams and domains, not technical layers
3. **Distributed Transactions:** Saga pattern with compensating actions
4. **Resilience:** Circuit breakers, retries, bulkheads, timeouts
5. **Observability:** Metrics, logs, traces, and SLOs

**Interviewer Expectations:**

| Level | Expectation |
|-------|-------------|
| Senior (L5) | Know patterns, apply when prompted |
| Staff (L6) | Proactively bring up resilience patterns |
| Principal (L7) | Discuss organizational impact, build vs buy |

**Questions That Probe Depth:**
- "How would you handle partial failures?"
- "What's your alerting strategy?"
- "How do you know if the system is healthy?"
- "What happens when Service X goes down?"

**Red Flags:**
- No mention of observability
- No failure handling discussion
- "We'll use microservices" without justification
- No compensation strategy for distributed operations

**Green Flags:**
- "Circuit breaker to prevent cascade"
- "Saga with compensating transactions"
- "SLO-based alerting with error budget"
- "Start monolith, extract when team boundaries form"

---

## Common Interview Questions & Model Answers

This section provides realistic interview questions on scalability patterns, with ideal answers and follow-up questions.

---

### Q1: When would you break a monolith into microservices? What are the trade-offs?

**Ideal Answer:**

"Breaking a monolith into microservices is a major decision with significant trade-offs. Here's when and why:

**When to keep the monolith:**
- Team < 10 engineers
- Product still finding product-market fit
- Simple domain with few bounded contexts
- Deployment complexity not justified yet

**When to consider microservices:**

**1. Team scaling (organizational driver)**
- Team > 50 engineers
- Multiple teams stepping on each other's code
- Deployment coordination becoming painful
- Example: Can't deploy Team A's feature without Team B's changes

**2. Technical scaling (performance driver)**
- Different components have different scaling needs
- Example: Image processing needs GPUs, API needs CPUs
- Monolith forces scaling everything together (expensive)

**3. Technology diversity**
- Different services benefit from different tech stacks
- Example: ML model in Python, API in Go
- Monolith locks you into one language

**4. Deployment independence**
- Need to deploy services independently
- Different release cycles (daily vs weekly)
- Reduce blast radius (deploy one service, not entire app)

**Trade-offs:**

| Aspect | Monolith | Microservices |
|--------|----------|---------------|
| Complexity | Simple | High (distributed systems) |
| Deployment | One artifact | Many (orchestration needed) |
| Debugging | Easy (single process) | Hard (distributed tracing) |
| Performance | Fast (in-process calls) | Slower (network calls) |
| Transactions | ACID (database) | Saga pattern (eventual) |
| Team autonomy | Low | High |
| Infrastructure cost | Low | High (more services) |

**Migration strategy (not all-or-nothing):**

**Phase 1: Modular monolith**
- Create clear module boundaries within monolith
- Enforce module interfaces (no cross-module DB access)
- Test deployment independently (feature flags)

**Phase 2: Extract by bounded context**
- Identify high-value candidates for extraction:
  - Independent business capability (payment processing)
  - Different scaling needs (image processing)
  - Different team ownership (identity service)
- Start with stateless services (easier)

**Phase 3: Data separation**
- Each service owns its database
- Use API calls instead of shared DB
- Eventual consistency where possible

**Phase 4: Handle distributed concerns**
- Service discovery (Kubernetes, Consul)
- Distributed tracing (Jaeger, Zipkin)
- Circuit breakers, retries
- Saga pattern for transactions

**Red flags (premature microservices):**
- Startup with 3 engineers building 20 microservices
- No clear service boundaries
- "We need to be like Netflix" (they have 1000s of engineers)

**My recommendation:**
- Start with modular monolith (80% of systems should stay here)
- Extract selectively when organizational or technical pressure is clear
- Use strangler fig pattern (gradually extract services)

**Real example:** Shopify stayed monolithic for years, only extracted specific services (payments, fulfillment) when they hit clear scaling limits."

**Follow-up Q:** "How do you handle transactions across microservices?"

**Ideal Answer:**

"Distributed transactions are one of the hardest problems in microservices. Here are the approaches:

**1. Saga Pattern (most common)**
- Break transaction into multiple local transactions
- Each service commits independently
- If one fails, compensate (undo) previous steps

**Example: E-commerce order**
```
1. Order Service: Create order (status=PENDING)
2. Inventory Service: Reserve items
3. Payment Service: Charge customer
4. Order Service: Mark order CONFIRMED

If Payment fails:
  - Compensate: Release inventory
  - Compensate: Mark order CANCELLED
```

**Two types of sagas:**

**A. Choreography (event-driven)**
```
Order Service: OrderCreated event → 
Inventory Service: InventoryReserved event →
Payment Service: PaymentSuccessful event →
Order Service: OrderConfirmed

If PaymentFailed event:
  Inventory Service: Release reservation
  Order Service: Cancel order
```
- **Pros:** Decoupled, no orchestrator
- **Cons:** Hard to track, debug (event soup)

**B. Orchestration (coordinator)**
```
Order Saga Orchestrator:
  1. Call Inventory Service (reserve)
  2. Call Payment Service (charge)
  3. If both succeed → confirm order
  4. If Payment fails → tell Inventory to release
```
- **Pros:** Clear flow, easier to debug
- **Cons:** Orchestrator is a single point of coordination
- **My preference:** Orchestration for business-critical flows

**2. Two-Phase Commit (2PC) - Avoid if possible**
```
Phase 1 (Prepare):
  Coordinator asks: "Can you commit?"
  Each service: "Yes, I can" (locks resources)

Phase 2 (Commit):
  Coordinator: "Everyone commit"
  Each service: Commits
```
- **Pros:** ACID guarantees
- **Cons:** Blocking (locks held during network calls), coordinator is single point of failure
- **When to use:** Never in microservices (too slow, fragile)

**3. Eventually Consistent (best for most cases)**
- Accept temporary inconsistency
- Example: Payment succeeded, but order confirmation email delayed
- User sees "Processing order" for a few seconds
- **Pros:** High availability, resilient
- **Cons:** Complex error handling, user experience

**4. Outbox Pattern (reliable event publishing)**
```
with transaction:
  update_order(order_id, status='CONFIRMED')
  insert_into_outbox({
    'event': 'OrderConfirmed',
    'payload': {...}
  })

Background job:
  Read from outbox table
  Publish events to message queue
  Mark as published
```
- Guarantees: DB update and event publishing are atomic
- Prevents: Order confirmed in DB but event lost

**Best practices:**

**1. Design for idempotency**
- All operations must be safely retryable
- Use idempotency keys

**2. Timeouts and retries**
- Set aggressive timeouts (don't wait forever)
- Exponential backoff on retries
- Circuit breakers to prevent cascade

**3. Compensating actions must be idempotent**
- Releasing inventory twice should be safe
- Refunding twice should be safe

**4. Monitor saga state**
- Track which step each saga is on
- Alert on stuck sagas (>5 minutes)
- Dashboard showing active, completed, failed sagas

**When NOT to use microservices:**
- If you need ACID transactions everywhere
- Domain doesn't have clear boundaries
- Team isn't ready for operational complexity

**Real example:** Uber uses saga pattern for trip booking:
1. Driver service: Assign driver
2. Rider service: Update rider
3. Payment service: Pre-authorize card
4. If any fails: Compensate (unassign driver, cancel ride)"

---

### Q2: Explain the Circuit Breaker pattern. How would you implement it?

**Ideal Answer:**

"The **Circuit Breaker** pattern prevents cascading failures by stopping requests to a failing service, giving it time to recover.

**Analogy:** Like an electrical circuit breaker—if too much current flows (too many failures), it trips (opens) to prevent fire (system crash).

**Three states:**

**1. CLOSED (normal operation)**
- Requests flow through normally
- Track failures
- If failure rate > threshold → Open circuit

**2. OPEN (failure mode)**
- Reject requests immediately (fail fast)
- Don't even try the failing service
- After timeout (e.g., 30 seconds) → Half-Open

**3. HALF-OPEN (testing recovery)**
- Allow limited requests through (e.g., 1 request)
- If succeeds → Close circuit (back to normal)
- If fails → Open circuit (back to failure mode)

**State diagram:**
```
CLOSED --[too many failures]--> OPEN
OPEN --[timeout expires]--> HALF-OPEN
HALF-OPEN --[success]--> CLOSED
HALF-OPEN --[failure]--> OPEN
```

**Implementation:**

```python
class CircuitBreaker:
    def __init__(self, failure_threshold=5, timeout=60, half_open_attempts=1):
        self.failure_threshold = failure_threshold
        self.timeout = timeout  # seconds to wait before trying again
        self.half_open_attempts = half_open_attempts
        
        self.failure_count = 0
        self.last_failure_time = None
        self.state = 'CLOSED'
        self.half_open_count = 0
    
    def call(self, func, *args, **kwargs):
        if self.state == 'OPEN':
            # Check if timeout has passed
            if time.now() - self.last_failure_time > self.timeout:
                self.state = 'HALF-OPEN'
                self.half_open_count = 0
            else:
                raise CircuitBreakerOpenError('Service unavailable')
        
        try:
            result = func(*args, **kwargs)
            self.on_success()
            return result
        except Exception as e:
            self.on_failure()
            raise
    
    def on_success(self):
        if self.state == 'HALF-OPEN':
            self.half_open_count += 1
            if self.half_open_count >= self.half_open_attempts:
                self.state = 'CLOSED'
                self.failure_count = 0
        elif self.state == 'CLOSED':
            self.failure_count = 0  # Reset on success
    
    def on_failure(self):
        self.failure_count += 1
        self.last_failure_time = time.now()
        
        if self.state == 'HALF-OPEN':
            self.state = 'OPEN'  # Back to open on any failure
        elif self.state == 'CLOSED':
            if self.failure_count >= self.failure_threshold:
                self.state = 'OPEN'

# Usage
breaker = CircuitBreaker(failure_threshold=5, timeout=30)

def call_payment_service():
    return breaker.call(payment_api.charge, user_id, amount)
```

**Configuration tuning:**

| Parameter | Typical Value | Purpose |
|-----------|---------------|---------|
| Failure threshold | 5-10 failures | How many failures before opening |
| Timeout | 30-60 seconds | How long to wait before retrying |
| Half-open attempts | 1-3 requests | How many test requests before closing |

**When to use:**

**1. External service calls**
- Payment gateway, third-party APIs
- Protect your system from their failures

**2. Database queries**
- If DB is overloaded, stop sending more queries
- Give it time to recover

**3. Microservice communication**
- Service A calls Service B
- If B is down, don't keep trying

**Benefits:**

1. **Fail fast:** Return error immediately instead of waiting for timeout
2. **Prevent cascade:** Don't overload failing service with retries
3. **Give time to recover:** Failing service gets breathing room
4. **Improve UX:** Fast failure = faster response to user (show cached data or error)

**Combined with other patterns:**

**+ Fallback:**
```python
try:
    return breaker.call(payment_service.charge)
except CircuitBreakerOpenError:
    return queue_for_retry()  # Process payment later
```

**+ Timeout:**
```python
try:
    return breaker.call(timeout(payment_service.charge, 3_seconds))
except TimeoutError:
    # Counts as failure, increments circuit breaker count
```

**+ Bulkhead:**
- Separate circuit breakers for different dependencies
- Payment service failure doesn't affect inventory service

**Monitoring:**

```python
metrics.gauge('circuit_breaker.state', state)  # 0=closed, 1=half-open, 2=open
metrics.counter('circuit_breaker.failures')
metrics.counter('circuit_breaker.open_count')
metrics.timer('circuit_breaker.recovery_time')

# Alert if circuit has been open for > 5 minutes
if state == 'OPEN' and duration > 300:
    alert('Circuit breaker stuck open: payment_service')
```

**Real examples:**
- Netflix Hystrix (Java circuit breaker library)
- AWS Lambda: Automatic throttling (similar concept)
- Stripe API: Rate limiting + circuit breaking for failing merchants

**My recommendation:**
- Use library (resilience4j, polly) instead of building from scratch
- Set conservative thresholds initially (high threshold, long timeout)
- Tune based on metrics and alerts
- Always have fallback strategy (cache, queue, graceful degradation)"

**Follow-up Q:** "What's the difference between Circuit Breaker and Retry?"

**Ideal Answer:**

"Both handle failures, but they're complementary and serve different purposes:

**Retry:**
- **Purpose:** Handle transient failures (temporary network glitch)
- **Assumption:** Problem will likely resolve quickly
- **Action:** Try again immediately or after brief delay
- **Risk:** Can make problem worse (overload failing service)

**Circuit Breaker:**
- **Purpose:** Handle sustained failures (service is down)
- **Assumption:** Problem will NOT resolve in milliseconds
- **Action:** Stop trying, fail fast
- **Risk:** Might give up too early on intermittent issues

**When to use each:**

| Scenario | Pattern | Reason |
|----------|---------|--------|
| Network blip | Retry | Likely transient, retry helps |
| Service overloaded | Circuit Breaker | Retries make it worse |
| Database deadlock | Retry | Retrying after delay may succeed |
| Database down | Circuit Breaker | Retries won't help, fail fast |

**Best practice: Use both together**

```python
# Outer: Circuit breaker (protects against sustained failure)
# Inner: Retry (handles transient failures)

circuit_breaker = CircuitBreaker(failure_threshold=5, timeout=30)
retry_policy = Retry(max_attempts=3, backoff=ExponentialBackoff())

def call_service():
    return circuit_breaker.call(
        lambda: retry_policy.execute(actual_service_call)
    )

Flow:
1. Circuit breaker checks state
2. If CLOSED, allow through
3. Retry policy tries up to 3 times
4. If all 3 fail, circuit breaker increments failure count
5. After 5 such failures, circuit breaker opens
```

**Configuration guidelines:**

```
Retry:
- Max attempts: 2-3 (aggressive retries harmful)
- Backoff: Exponential (100ms, 200ms, 400ms)
- Jitter: Add randomness to prevent thundering herd

Circuit Breaker:
- Failure threshold: 5-10 (tolerate some failures)
- Timeout: 30-60 seconds (time to recover)
- Half-open attempts: 1-3 (test recovery conservatively)
```

**Real example:** AWS SDK uses both:
- Retries: 3 attempts with exponential backoff for transient errors (500, 503)
- Circuit breaking: If 10 consecutive failures (service down), stop trying for 60 seconds"

---

### Q3: What is the Saga pattern? Give an example of choreography vs orchestration.

**Ideal Answer:**

"The **Saga pattern** manages distributed transactions as a sequence of local transactions with compensating actions if something fails.

**Problem it solves:**
- In microservices, you can't use database ACID transactions across services
- Need to maintain consistency without distributed locks

**Two types:**

**1. Choreography (event-driven, decentralized)**

Each service listens for events and publishes new events. No central coordinator.

**Example: E-commerce order**

```
┌─────────────┐
│Order Service│ OrderCreated event
└──────┬──────┘
       │
       v
┌─────────────────┐
│Inventory Service│ InventoryReserved event
└───────┬─────────┘
        │
        v
┌─────────────┐
│Payment Svc  │ PaymentSuccessful event OR PaymentFailed event
└──────┬──────┘
       │
       v (if successful)
┌─────────────┐
│Order Service│ OrderConfirmed
└─────────────┘

If PaymentFailed:
  InventoryReleased event ← Inventory Service
  OrderCancelled event ← Order Service
```

**Implementation:**
```python
# Order Service
def create_order(order):
    db.insert_order(order, status='PENDING')
    event_bus.publish('OrderCreated', order)

def on_payment_successful(event):
    db.update_order(event.order_id, status='CONFIRMED')
    event_bus.publish('OrderConfirmed', event)

def on_payment_failed(event):
    db.update_order(event.order_id, status='CANCELLED')
    event_bus.publish('OrderCancelled', event)

# Inventory Service
def on_order_created(event):
    if can_reserve(event.items):
        reserve_inventory(event.items)
        event_bus.publish('InventoryReserved', event)
    else:
        event_bus.publish('InventoryFailed', event)

def on_order_cancelled(event):
    release_inventory(event.items)  # Compensating action
    event_bus.publish('InventoryReleased', event)

# Payment Service
def on_inventory_reserved(event):
    if charge_customer(event.user_id, event.amount):
        event_bus.publish('PaymentSuccessful', event)
    else:
        event_bus.publish('PaymentFailed', event)
```

**Pros:**
- Loose coupling (services don't know about each other)
- High availability (no single coordinator)
- Scales well

**Cons:**
- Hard to track saga state (distributed across services)
- Complex debugging (follow event trail)
- Risk of event loops
- Eventual consistency (order might be PENDING for seconds)

**2. Orchestration (centralized coordinator)**

A saga orchestrator coordinates all steps. Services don't publish events, just respond to orchestrator.

**Example: Same order flow**

```
┌──────────────────┐
│ Order Saga       │
│ Orchestrator     │
└────┬─────────────┘
     │
     ├─1─> Inventory.reserve()
     │     └─> Success/Failure
     │
     ├─2─> Payment.charge()
     │     └─> Success/Failure
     │
     └─3─> Order.confirm()
           OR
       Compensate:
       ├─> Payment.refund()
       └─> Inventory.release()
```

**Implementation:**
```python
class OrderSagaOrchestrator:
    def execute(self, order):
        saga_state = {
            'order_id': order.id,
            'status': 'STARTED',
            'steps_completed': []
        }
        
        try:
            # Step 1: Reserve inventory
            inventory_service.reserve(order.items)
            saga_state['steps_completed'].append('INVENTORY_RESERVED')
            
            # Step 2: Charge payment
            payment_service.charge(order.user_id, order.amount)
            saga_state['steps_completed'].append('PAYMENT_CHARGED')
            
            # Step 3: Confirm order
            order_service.confirm(order.id)
            saga_state['status'] = 'COMPLETED'
            
        except InventoryError as e:
            # No compensation needed (nothing committed yet)
            order_service.cancel(order.id, reason='No inventory')
            saga_state['status'] = 'FAILED'
            
        except PaymentError as e:
            # Compensate: Release inventory
            self.compensate(saga_state)
            order_service.cancel(order.id, reason='Payment failed')
            saga_state['status'] = 'COMPENSATED'
        
        return saga_state
    
    def compensate(self, saga_state):
        # Execute compensating actions in reverse order
        if 'INVENTORY_RESERVED' in saga_state['steps_completed']:
            inventory_service.release(order.items)
        
        if 'PAYMENT_CHARGED' in saga_state['steps_completed']:
            payment_service.refund(order.user_id, order.amount)
```

**Pros:**
- Clear flow (easy to understand)
- Centralized monitoring (saga state in one place)
- Easier to debug (orchestrator logs show everything)
- Better for complex workflows (many steps, conditional logic)

**Cons:**
- Orchestrator is a dependency (coupling)
- Single point of coordination (but not failure—can be replicated)
- Orchestrator needs to be highly available

**Comparison:**

| Aspect | Choreography | Orchestration |
|--------|--------------|---------------|
| Coupling | Loose | Medium |
| Complexity | Distributed | Centralized |
| Debugging | Hard | Easy |
| Scalability | High | Medium (orchestrator bottleneck) |
| Use case | Simple flows | Complex flows |

**My recommendation:**
- **Choreography:** Simple flows (2-3 steps), eventual consistency OK
- **Orchestration:** Business-critical flows (orders, payments), complex logic

**Real-world hybrid:**
- High-level: Orchestration (order saga orchestrator)
- Low-level: Choreography (within each service, event-driven)

**Patterns to combine:**

**1. Outbox pattern** (reliable event publishing)
```python
with transaction:
    update_order(status='CONFIRMED')
    insert_into_outbox({'event': 'OrderConfirmed', 'payload': order})

# Background job publishes events from outbox
```

**2. Idempotency** (safe retries)
```python
def reserve_inventory(idempotency_key, items):
    if already_reserved(idempotency_key):
        return  # Idempotent
    
    actually_reserve(items)
    save_idempotency_key(idempotency_key)
```

**3. Timeouts** (don't wait forever)
```python
try:
    inventory_service.reserve(order.items, timeout=3_seconds)
except TimeoutError:
    compensate_and_fail()
```

**Real examples:**
- Uber: Orchestration for trip booking saga
- Netflix: Choreography for content delivery pipeline
- Amazon: Hybrid (orchestration for checkout, choreography for recommendations)"

---

### Q4: Explain SLOs, SLIs, and SLAs. How do you set them?

**Ideal Answer:**

"These are related but distinct concepts for measuring reliability:

**SLI (Service Level Indicator):**
- **What:** A metric that measures service health
- **Examples:** 
  - Request latency (p99 < 200ms)
  - Error rate (< 0.1%)
  - Availability (% of successful requests)

**SLO (Service Level Objective):**
- **What:** Target value for an SLI
- **Examples:**
  - \"99.9% of requests succeed\"
  - \"p99 latency < 200ms\"
  - \"99.95% uptime per month\"

**SLA (Service Level Agreement):**
- **What:** Contract with consequences if SLO violated
- **Examples:**
  - \"If uptime < 99.9%, customer gets 10% credit\"
  - Legal/financial commitment
  - Always less strict than internal SLO (buffer)

**Relationship:**
```
SLI: How we measure (latency, errors, uptime)
  ↓
SLO: Our internal target (99.9% success rate)
  ↓
SLA: Promise to customers (99% success rate + penalty)
```

**Example:**

| Layer | Description |
|-------|-------------|
| **SLI** | Request success rate: `successful_requests / total_requests` |
| **SLO** | 99.95% of requests succeed in any 30-day window |
| **SLA** | Guarantee 99.9% success rate; if breached, 10% refund |

**How to set SLOs:**

**1. Start with current performance**
```
Look at last 90 days:
- p50 latency: 50ms
- p99 latency: 180ms
- Success rate: 99.97%

Set SLO slightly below current performance:
- p99 latency: 200ms (gives 20ms buffer)
- Success rate: 99.95% (gives 0.02% error budget)
```

**2. Consider user expectations**
- What do users notice? (>500ms latency = slow)
- What's acceptable for the product? (Social media vs banking)

**3. Error budget**
```
SLO: 99.9% uptime per month
Error budget: 0.1% downtime allowed
= 30 days * 24 hours * 60 min * 0.001
= 43.2 minutes of downtime per month

Use error budget for:
- Deployments
- Experimentation
- Infrastructure maintenance
```

**4. Multiple SLOs**

Don't use a single number. Track multiple dimensions:

```
API SLOs:
1. Availability: 99.95% of requests return 2xx/3xx
2. Latency: p99 < 200ms
3. Correctness: 99.99% of payments succeed

Why multiple?
- Can have low error rate but high latency (bad UX)
- Can have low latency but high error rate (unreliable)
```

**Setting thresholds:**

| Tier | Uptime SLO | Downtime/month | Use Case |
|------|------------|----------------|----------|
| 99% | Two 9s | 7.2 hours | Non-critical, internal tools |
| 99.9% | Three 9s | 43.2 minutes | Standard SaaS |
| 99.95% | | 21.6 minutes | High-value SaaS |
| 99.99% | Four 9s | 4.3 minutes | Financial, healthcare |
| 99.999% | Five 9s | 26 seconds | Critical infrastructure |

**SLO-based alerting:**

Instead of alerting on every error, alert on SLO burn rate.

```python
# Bad: Alert on every error
if error_rate > 0:
    alert()  # Too noisy

# Good: Alert on SLO violation trajectory
def slo_burn_rate_alert():
    error_rate_1h = errors_last_hour / requests_last_hour
    monthly_budget = 0.001  # 99.9% SLO = 0.1% error budget
    
    # If current rate continues, will we violate SLO?
    projected_monthly_errors = error_rate_1h * 30 * 24
    
    if projected_monthly_errors > monthly_budget:
        alert('Burning SLO budget at 14x rate')
```

**Multi-window alerting:**

```
Alert levels:
1. Critical: 5% of error budget consumed in 1 hour (36x burn rate)
2. Warning: 10% of budget consumed in 6 hours (20x burn rate)
3. Notice: 50% of budget consumed in 3 days (normal rate)

This prevents:
- Alert fatigue (not every blip)
- Slow burns (gradually degrading, noticed too late)
```

**Monitoring dashboard:**

```
┌─────────────────────────────┐
│ API Health - Last 30 Days   │
├─────────────────────────────┤
│ Success Rate:  99.96% ✓     │
│   SLO: 99.95%               │
│   Error Budget: 12% used    │
├─────────────────────────────┤
│ p99 Latency:  185ms ✓       │
│   SLO: 200ms                │
│   Budget: 8% used           │
├─────────────────────────────┤
│ p50 Latency:  45ms ✓        │
└─────────────────────────────┘
```

**When you've violated SLO:**

```
1. Stop new deployments (preserve error budget)
2. Root cause analysis (what changed?)
3. Implement fix or rollback
4. Post-mortem (blameless, focus on systems)
5. Update runbooks
6. If needed, adjust SLO (maybe too aggressive)
```

**Real examples:**

**Google Search:**
- SLI: Query success rate, latency
- SLO: 99.9% success, p99 < 200ms
- SLA: No public SLA (free product)

**AWS S3:**
- SLI: Request success rate, durability
- SLO: Internal (likely 99.99%+)
- SLA: 99.9% availability or get service credits

**Stripe API:**
- SLI: API success rate, latency
- SLO: 99.99% uptime (internal)
- SLA: 99.95% uptime with credits

**My recommendations:**
1. Start with 99.9% (achievable without heroics)
2. Measure for 90 days before setting SLO
3. Set SLA lower than SLO (give yourself buffer)
4. Alert on SLO burn rate, not individual errors
5. Review SLOs quarterly (adjust based on reality)"

---

**Navigation:** [← Previous: Distributed Systems](02-distributed-systems.md) | [Next: Case Studies (Mandatory) →](04-case-studies-mandatory.md)
- "Saga with compensating transactions"
- "SLO-based alerting with error budget"
- "Start monolith, extract when team boundaries form"
