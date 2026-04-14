# Spring Transaction Management — Complete Production Guide

## Table of Contents
1. [Core Concepts](#1-core-concepts-must-know)
2. [@Transactional Deep Dive](#2-transactional-deep-dive)
3. [Propagation](#3-propagation-top-interview-topic)
4. [Isolation Levels](#4-isolation-levels-critical-for-concurrency)
5. [Real Production Scenarios](#5-real-production-scenarios)
6. [Spring + JPA Specific](#6-spring--jpa-specific)
7. [Common Bugs](#7-common-bugs-very-important-for-interview)
8. [Advanced / Architecture Level](#8-advanced--architecture-level)
9. [Scenario-Based Questions](#9-scenario-based-questions-interview-gold)
10. [Bonus — Senior Level Internals](#10-bonus--senior-level-internals)

---

## 1. Core Concepts (Must Know)

### 1.1 What is a Transaction?

A **transaction** is a logical unit of work that groups one or more database operations into a single indivisible operation. Either **all** operations succeed, or **none** of them take effect.

```
BEGIN TRANSACTION
    ├── Debit ₹5000 from Account A
    ├── Credit ₹5000 to Account B
    └── Log the transfer record
COMMIT  (all succeed)
   or
ROLLBACK (all reverted)
```

**Real-world analogy:** Think of a bank transfer — you never want money debited from one account without being credited to another. The transaction guarantees both happen or neither does.

---

### 1.2 ACID Properties

| Property | Meaning | What Happens Without It | Example |
|----------|---------|-------------------------|---------|
| **Atomicity** | All-or-nothing execution | Partial updates corrupt data | Debit succeeds but credit fails → money disappears |
| **Consistency** | DB moves from one valid state to another | Constraint violations, orphan records | Account balance goes negative violating business rule |
| **Isolation** | Concurrent transactions don't interfere | Dirty reads, lost updates | Two users buy last ticket simultaneously |
| **Durability** | Committed data survives crashes | Data loss on server restart | Payment recorded but lost after power failure |

```
          ┌──────────────────────────────────────────────────┐
          │                  ACID Properties                  │
          ├────────────┬────────────┬────────────┬───────────┤
          │  Atomicity │ Consistency│ Isolation  │ Durability│
          │            │            │            │           │
          │  All or    │  Valid     │ Concurrent │ Survives  │
          │  Nothing   │  State     │ Safety     │ Crashes   │
          │            │  Always    │            │           │
          │  Undo Log  │ Constraints│ MVCC/Locks │ WAL/Redo  │
          │  (Rollback)│  + Rules   │            │   Log     │
          └────────────┴────────────┴────────────┴───────────┘
```

**Production insight:** In distributed systems, full ACID across services is impractical. We often relax to **BASE** (Basically Available, Soft state, Eventually consistent) and use patterns like Saga.

---

### 1.3 What is Transaction Management in Spring?

Spring provides a **consistent abstraction** over different transaction APIs (JDBC, JPA, JTA, Hibernate) so you write transaction logic once regardless of the underlying technology.

```
┌─────────────────────────────────────────────────────┐
│                  Your Service Code                   │
│              @Transactional / TX Template            │
├─────────────────────────────────────────────────────┤
│          Spring Transaction Abstraction              │
│          (PlatformTransactionManager)                │
├──────────┬──────────┬──────────┬────────────────────┤
│   JDBC   │   JPA    │   JTA    │   Hibernate        │
│   TX     │   TX     │   TX     │   TX               │
├──────────┴──────────┴──────────┴────────────────────┤
│               Database / Message Broker              │
└─────────────────────────────────────────────────────┘
```

**Why this matters in production:**
- Switch from JDBC to JPA without rewriting transaction code
- Consistent rollback/commit behavior across all data access technologies
- Unified configuration through Spring Boot auto-configuration

---

### 1.4 Programmatic vs Declarative Transaction Management

| Aspect | Programmatic | Declarative |
|--------|-------------|-------------|
| **How** | Manual code using `TransactionTemplate` or `PlatformTransactionManager` | Annotations (`@Transactional`) or XML |
| **Coupling** | Transaction logic mixed with business logic | Clean separation of concerns |
| **Flexibility** | Fine-grained control over commit/rollback points | Applies to method boundaries |
| **Verbosity** | More boilerplate | Minimal code |
| **Use when** | Need partial commits, complex TX control | 95% of cases — standard CRUD operations |

#### Programmatic Approach

```java
@Service
public class PaymentService {

    private final TransactionTemplate transactionTemplate;

    public PaymentService(PlatformTransactionManager txManager) {
        this.transactionTemplate = new TransactionTemplate(txManager);
    }

    public PaymentResult processPayment(PaymentRequest request) {
        return transactionTemplate.execute(status -> {
            try {
                debitAccount(request.getFrom(), request.getAmount());
                creditAccount(request.getTo(), request.getAmount());
                return new PaymentResult(SUCCESS);
            } catch (InsufficientFundsException e) {
                status.setRollbackOnly();
                return new PaymentResult(FAILED, e.getMessage());
            }
        });
    }
}
```

#### Declarative Approach

```java
@Service
public class PaymentService {

    @Transactional(rollbackFor = Exception.class)
    public PaymentResult processPayment(PaymentRequest request) {
        debitAccount(request.getFrom(), request.getAmount());
        creditAccount(request.getTo(), request.getAmount());
        return new PaymentResult(SUCCESS);
    }
}
```

**Production recommendation:** Use declarative (`@Transactional`) for 95% of cases. Reserve programmatic for scenarios like:
- Partial transaction commits within a single method
- Dynamic transaction configuration at runtime
- Callback-based async processing where you need explicit TX boundaries

---

### 1.5 What Does @Transactional Actually Do Internally?

When Spring encounters `@Transactional`, it creates a **proxy** around your bean. Here's the internal flow:

```
Client Code calls service.save(entity)
          │
          ▼
┌─────────────────────────────────────┐
│       Transaction Proxy (AOP)       │
│                                     │
│  1. Get TransactionManager          │
│  2. Check existing TX (propagation) │
│  3. Begin new TX if needed          │
│  4. Set isolation level             │
│  5. Bind Connection to ThreadLocal  │
│                                     │
│  ┌───────────────────────────────┐  │
│  │   Your Actual Method          │  │
│  │   service.save(entity)        │  │
│  │   (business logic runs here)  │  │
│  └───────────────────────────────┘  │
│                                     │
│  6. If no exception → COMMIT        │
│  7. If RuntimeException → ROLLBACK  │
│  8. Unbind Connection               │
│  9. Return result to caller         │
└─────────────────────────────────────┘
```

**Step-by-step internal flow:**

```
1. TransactionInterceptor.invoke()
      │
2. TransactionAspectSupport.invokeWithinTransaction()
      │
3. PlatformTransactionManager.getTransaction(definition)
      │
      ├── Check propagation behavior
      ├── Create/reuse transaction
      └── Bind DataSource Connection to ThreadLocal
              (via TransactionSynchronizationManager)
      │
4. Execute target method (your business logic)
      │
5. On success → txManager.commit(status)
   On failure → txManager.rollback(status)
      │
6. Unbind resources, clean up ThreadLocal
```

**Key insight:** The connection is bound to a `ThreadLocal`, which is why:
- All DB operations in the same thread share the same transaction
- Async operations (`@Async`) break the transaction boundary
- Thread pool tasks won't participate in the caller's transaction

---

### 1.6 Proxy-Based Transaction Management

Spring uses **proxies** to intercept method calls and wrap them with transaction logic. Two proxy mechanisms exist:

```
┌──────────────────────────────────────────────────────┐
│              JDK Dynamic Proxy                       │
│  (When bean implements an interface)                 │
│                                                      │
│  Interface ←── Proxy ──→ Target Bean                 │
│                  │                                    │
│          TransactionInterceptor                      │
└──────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────┐
│              CGLIB Proxy (Default in Spring Boot)    │
│  (When bean has no interface, or proxyTargetClass)   │
│                                                      │
│  Subclass Proxy extends Target Bean                  │
│         │                                            │
│  TransactionInterceptor                              │
└──────────────────────────────────────────────────────┘
```

| Proxy Type | Mechanism | Limitation |
|-----------|-----------|------------|
| **JDK Dynamic Proxy** | Creates a proxy implementing the same interface | Only works with interface-based beans |
| **CGLIB Proxy** | Creates a subclass of your bean class | Cannot proxy `final` classes or `final` methods |

**Spring Boot default:** CGLIB proxy (`spring.aop.proxy-target-class=true`).

---

### 1.7 Why Does @Transactional Not Work on Private Methods?

Because the proxy **cannot override** private methods.

```
┌─────────────────────────────────────────────┐
│  CGLIB Proxy (subclass of your bean)        │
│                                             │
│  public save() { ← CAN override, CAN add   │
│      // TX begin    transaction logic       │
│      super.save()                           │
│      // TX commit                           │
│  }                                          │
│                                             │
│  private helper() { ← CANNOT override,     │
│      // No TX logic   INVISIBLE to proxy    │
│  }                                          │
└─────────────────────────────────────────────┘
```

**Rules of visibility:**

| Method Modifier | Proxy Can Intercept? | @Transactional Works? |
|----------------|---------------------|----------------------|
| `public` | Yes | Yes |
| `protected` | CGLIB: Yes, JDK: No | Technically yes with CGLIB, but discouraged |
| `package-private` | CGLIB: Yes, JDK: No | Technically yes with CGLIB, but discouraged |
| `private` | No | **Never** |

**Production tip:** Always put `@Transactional` on `public` methods only. Spring will silently ignore it on private methods — no error, no warning, just no transaction.

---

### 1.8 Self-Invocation Problem

When a `@Transactional` method calls another `@Transactional` method **in the same class**, the second annotation is ignored.

```java
@Service
public class OrderService {

    @Transactional
    public void createOrder(Order order) {
        saveOrder(order);
        sendNotification(order); // THIS CALL BYPASSES THE PROXY
    }

    @Transactional(propagation = Propagation.REQUIRES_NEW)
    public void sendNotification(Order order) {
        // Intended: run in a NEW transaction
        // Reality: runs in the SAME transaction as createOrder
        notificationRepo.save(new Notification(order));
    }
}
```

```
External call → Proxy → createOrder()
                              │
                              │ (internal call — bypasses proxy)
                              ▼
                        sendNotification()  ← NO new transaction!
```

**Why?** The internal call uses `this.sendNotification()` which is the actual object, not the proxy.

**Solutions (ranked by preference):**

```
1. Extract to a separate @Service class (RECOMMENDED)
   OrderService → NotificationService.sendNotification()

2. Inject self reference (works but awkward)
   @Autowired @Lazy private OrderService self;
   self.sendNotification(order);

3. Use AopContext.currentProxy() (fragile, not recommended)
   ((OrderService) AopContext.currentProxy()).sendNotification(order);

4. Use AspectJ weaving instead of proxies (complex setup)
   @EnableTransactionManagement(mode = AdviceMode.ASPECTJ)
```

---

## 2. @Transactional Deep Dive

### 2.1 Key Attributes of @Transactional

```java
@Transactional(
    propagation = Propagation.REQUIRED,       // TX boundary behavior
    isolation = Isolation.DEFAULT,            // Concurrency control
    timeout = -1,                             // Seconds before TX timeout
    readOnly = false,                         // Optimization hint
    rollbackFor = {Exception.class},          // Force rollback for checked exceptions
    noRollbackFor = {MailException.class},    // Prevent rollback for specific exceptions
    transactionManager = "primaryTxManager"   // Which TX manager to use (multi-DB)
)
```

#### Detailed Attribute Breakdown

| Attribute | Default | Purpose | Production Usage |
|-----------|---------|---------|-----------------|
| `propagation` | `REQUIRED` | How TX boundaries interact | Critical for nested service calls |
| `isolation` | `DEFAULT` (DB default) | Concurrency control level | Rarely changed from default |
| `timeout` | `-1` (no timeout) | Max seconds for TX | **Always set in production** (e.g., 30s) |
| `readOnly` | `false` | Optimization hint for DB | Use for all read operations — enables replica routing |
| `rollbackFor` | `{RuntimeException, Error}` | Which exceptions trigger rollback | **Always set `Exception.class`** for safety |
| `noRollbackFor` | `{}` | Exceptions that should NOT rollback | Rarely used — specific business cases |
| `transactionManager` | Default bean | Which TX manager | Multi-datasource setups |

---

### 2.2 Propagation Attribute

Controls what happens when a transactional method is called while a transaction already exists. Covered in detail in [Section 3](#3-propagation-top-interview-topic).

---

### 2.3 Isolation Attribute

Controls visibility of data changes between concurrent transactions. Covered in detail in [Section 4](#4-isolation-levels-critical-for-concurrency).

---

### 2.4 Timeout Attribute

```java
@Transactional(timeout = 30) // 30 seconds
public void processLargeOrder(Order order) {
    // If this takes > 30 seconds, TransactionTimedOutException is thrown
}
```

**Production practice:** ALWAYS set a timeout. A missing timeout can cause:
- Connection pool exhaustion (HikariCP default max pool = 10)
- Database lock escalation
- Cascading failures across services

```
Recommended timeouts:
  Simple CRUD:           5-10 seconds
  Complex aggregation:   30 seconds
  Batch processing:      Use programmatic TX with chunking instead
  Never:                 timeout = -1 in production
```

---

### 2.5 readOnly Attribute

```java
@Transactional(readOnly = true)
public List<Order> getOrdersByUser(Long userId) {
    return orderRepository.findByUserId(userId);
}
```

**What `readOnly = true` actually does:**

| Layer | Effect |
|-------|--------|
| **Spring** | Sets `Connection.setReadOnly(true)` |
| **Hibernate/JPA** | Disables dirty checking → performance boost |
| **Hibernate** | Sets FlushMode to `MANUAL` → no automatic flushes |
| **MySQL** | May route to read replica (with appropriate DataSource config) |
| **PostgreSQL** | Sets `SET TRANSACTION READ ONLY` → DB rejects writes |

**Production tip:** Use `readOnly = true` on every read operation. It:
- Saves 20-30% memory (no dirty checking snapshots)
- Enables read-replica routing with libraries like `AbstractRoutingDataSource`
- Prevents accidental writes in query methods

---

### 2.6 Default Behavior of @Transactional

```java
@Transactional  // equivalent to:
@Transactional(
    propagation = Propagation.REQUIRED,
    isolation = Isolation.DEFAULT,
    timeout = -1,
    readOnly = false,
    rollbackFor = {RuntimeException.class, Error.class}
)
```

**Critical default to remember:** Only `RuntimeException` (unchecked) and `Error` trigger rollback by default. **Checked exceptions do NOT trigger rollback.**

---

### 2.7 When Does Spring Rollback a Transaction?

```
Method completes
      │
      ├── No exception → COMMIT
      │
      └── Exception thrown
              │
              ├── RuntimeException (unchecked) → ROLLBACK ✓
              ├── Error → ROLLBACK ✓
              └── Checked Exception (IOException, etc.) → COMMIT (!)
                    │
                    └── This is the #1 source of transaction bugs
```

---

### 2.8 Checked vs Unchecked Exceptions in Rollback

| Exception Type | Examples | Default Rollback? | Why? |
|---------------|----------|-------------------|------|
| **Unchecked** (`RuntimeException`) | `NullPointerException`, `IllegalArgumentException`, custom runtime exceptions | **Yes** | Assumed to be unrecoverable programming errors |
| **Checked** (`Exception`) | `IOException`, `SQLException`, `CustomBusinessException` | **No** | Assumed to be recoverable business scenarios |
| **Error** | `OutOfMemoryError`, `StackOverflowError` | **Yes** | JVM-level failures |

**This is a dangerous default.** Consider this:

```java
@Transactional
public void transferMoney(TransferRequest req) throws InsufficientFundsException {
    accountRepo.debit(req.getFrom(), req.getAmount());
    accountRepo.credit(req.getTo(), req.getAmount());

    if (balance < 0) {
        throw new InsufficientFundsException("Not enough funds");
        // If InsufficientFundsException extends Exception (checked)
        // → Transaction COMMITS → money debited but not credited!
    }
}
```

---

### 2.9 How to Force Rollback for Checked Exceptions

```java
// Option 1: rollbackFor on the method (RECOMMENDED)
@Transactional(rollbackFor = Exception.class)
public void transferMoney(TransferRequest req) throws InsufficientFundsException {
    // Now ALL exceptions trigger rollback
}

// Option 2: rollbackFor for specific exceptions
@Transactional(rollbackFor = {InsufficientFundsException.class, PaymentGatewayException.class})
public void transferMoney(TransferRequest req) throws Exception {
    // Only specified checked exceptions trigger rollback
}

// Option 3: Make your business exception extend RuntimeException
public class InsufficientFundsException extends RuntimeException {
    // Now automatically triggers rollback
}
```

**Production standard:** Always use `rollbackFor = Exception.class` unless you have a specific reason not to. Many teams enforce this via a shared base `@Transactional` configuration or a custom composed annotation:

```java
@Target({ElementType.METHOD, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
@Transactional(rollbackFor = Exception.class, timeout = 30)
public @interface SafeTransaction {
}
```

---

## 3. Propagation (Top Interview Topic)

### 3.1 What is Transaction Propagation?

Propagation defines **how transactions relate to each other** when a transactional method calls another transactional method.

```
ServiceA.methodA()  ──calls──→  ServiceB.methodB()
      │                                │
   Has TX?                      What should B do?
      │                                │
      └── Propagation rules determine the answer
```

---

### 3.2 All Propagation Types Explained

#### Visual Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                    PROPAGATION TYPES                              │
├──────────────┬────────────────────────┬──────────────────────────┤
│   Type       │  Existing TX Present   │  No Existing TX          │
├──────────────┼────────────────────────┼──────────────────────────┤
│ REQUIRED     │  Join existing TX      │  Create new TX           │
│ REQUIRES_NEW │  Suspend + Create new  │  Create new TX           │
│ NESTED       │  Create savepoint      │  Create new TX           │
│ SUPPORTS     │  Join existing TX      │  Run without TX          │
│ NOT_SUPPORTED│  Suspend existing TX   │  Run without TX          │
│ NEVER        │  Throw exception!      │  Run without TX          │
│ MANDATORY    │  Join existing TX      │  Throw exception!        │
└──────────────┴────────────────────────┴──────────────────────────┘
```

---

#### REQUIRED (Default)

```java
@Transactional(propagation = Propagation.REQUIRED)
public void placeOrder(Order order) { ... }
```

```
Case 1: No existing TX              Case 2: Existing TX present
┌─────────────────────┐             ┌─────────────────────────────┐
│  Creates NEW TX     │             │  Outer TX                    │
│  ┌───────────────┐  │             │  ┌────────────────────────┐ │
│  │ placeOrder()  │  │             │  │ checkout()             │ │
│  └───────────────┘  │             │  │   ┌──────────────────┐ │ │
│  Commits on success │             │  │   │ placeOrder()     │ │ │
└─────────────────────┘             │  │   │ JOINS outer TX   │ │ │
                                    │  │   └──────────────────┘ │ │
                                    │  └────────────────────────┘ │
                                    │  Single commit/rollback     │
                                    └─────────────────────────────┘
```

**Use when:** Default choice. Most service methods should use this.

---

#### REQUIRES_NEW

```java
@Transactional(propagation = Propagation.REQUIRES_NEW)
public void logAuditEvent(AuditEvent event) { ... }
```

```
┌────────────────────────────────────────────────┐
│  Outer TX (SUSPENDED)                          │
│  ┌──────────────────────────────────────────┐  │
│  │ processPayment()                         │  │
│  │    │                                     │  │
│  │    │ SUSPEND outer TX                    │  │
│  │    ▼                                     │  │
│  │  ┌──────────────────────────┐            │  │
│  │  │ NEW TX (independent)     │            │  │
│  │  │ logAuditEvent()          │            │  │
│  │  │ Commits independently    │            │  │
│  │  └──────────────────────────┘            │  │
│  │    │                                     │  │
│  │    │ RESUME outer TX                     │  │
│  │    ▼                                     │  │
│  │  (continue processPayment)               │  │
│  └──────────────────────────────────────────┘  │
│  Outer TX commit/rollback independent          │
└────────────────────────────────────────────────┘
```

**Use when:**
- Audit logging (must persist even if main TX fails)
- Generating sequence numbers
- Sending notifications that shouldn't roll back
- Any operation that must commit regardless of the outer transaction's outcome

---

#### NESTED

```java
@Transactional(propagation = Propagation.NESTED)
public void applyDiscount(Order order) { ... }
```

```
┌──────────────────────────────────────────────────┐
│  Outer TX                                         │
│  ┌────────────────────────────────────────────┐  │
│  │ processOrder()                              │  │
│  │    │                                        │  │
│  │    │ SAVEPOINT sp1                          │  │
│  │    ▼                                        │  │
│  │  ┌──────────────────────────┐               │  │
│  │  │ Nested TX (savepoint)    │               │  │
│  │  │ applyDiscount()          │               │  │
│  │  │                          │               │  │
│  │  │ Fail? → ROLLBACK TO sp1 │               │  │
│  │  │ (outer TX continues!)    │               │  │
│  │  │                          │               │  │
│  │  │ Success? → Release sp1   │               │  │
│  │  │ (final commit with outer)│               │  │
│  │  └──────────────────────────┘               │  │
│  └────────────────────────────────────────────┘  │
│  Outer TX commit includes nested changes          │
│  Outer TX rollback rolls back nested too          │
└──────────────────────────────────────────────────┘
```

**Use when:**
- Optional steps that can fail without killing the whole transaction
- Applying coupons/discounts (fail gracefully, continue order)
- Batch processing where individual items can fail

**Important:** Not all databases/drivers support savepoints. JDBC savepoints are required. Works with MySQL, PostgreSQL. Does NOT work with JTA transactions.

---

#### SUPPORTS

```java
@Transactional(propagation = Propagation.SUPPORTS)
public User findUser(Long id) { ... }
```

```
Case 1: Existing TX → Join it       Case 2: No TX → Run without TX
┌───────────────────────┐            ┌───────────────────────┐
│  Outer TX             │            │  No TX context        │
│  ┌─────────────────┐  │            │  ┌─────────────────┐  │
│  │ findUser()      │  │            │  │ findUser()      │  │
│  │ Runs in TX      │  │            │  │ Runs without TX │  │
│  └─────────────────┘  │            │  └─────────────────┘  │
└───────────────────────┘            └───────────────────────┘
```

**Use when:** Read-only methods that can work either way.

---

#### NOT_SUPPORTED

```java
@Transactional(propagation = Propagation.NOT_SUPPORTED)
public void generateReport() { ... }
```

```
┌──────────────────────────────────────────┐
│  Outer TX (SUSPENDED)                    │
│  ┌────────────────────────────────────┐  │
│  │ generateReport()                   │  │
│  │ Runs WITHOUT any TX                │  │
│  │ (outer TX suspended during this)   │  │
│  └────────────────────────────────────┘  │
│  Outer TX RESUMES after                  │
└──────────────────────────────────────────┘
```

**Use when:**
- Long-running read operations (reports) that would hold locks too long
- Operations that explicitly don't need transactional guarantees

---

#### NEVER

```java
@Transactional(propagation = Propagation.NEVER)
public void sendEmail(String to, String body) { ... }
```

```
Case 1: Existing TX → EXCEPTION!     Case 2: No TX → Run normally
┌───────────────────────────────┐    ┌───────────────────────┐
│  IllegalTransactionState      │    │  No TX context        │
│  Exception thrown!             │    │  ┌─────────────────┐  │
│                               │    │  │ sendEmail()     │  │
│  "Existing transaction found  │    │  │ Runs without TX │  │
│   for transaction marked      │    │  └─────────────────┘  │
│   with propagation 'never'"   │    └───────────────────────┘
└───────────────────────────────┘
```

**Use when:** Enforcing that a method must never run inside a transaction (e.g., external API calls, email sending).

---

#### MANDATORY

```java
@Transactional(propagation = Propagation.MANDATORY)
public void debitAccount(Long accountId, BigDecimal amount) { ... }
```

```
Case 1: Existing TX → Join it       Case 2: No TX → EXCEPTION!
┌───────────────────────┐            ┌───────────────────────────┐
│  Outer TX             │            │  IllegalTransactionState   │
│  ┌─────────────────┐  │            │  Exception thrown!          │
│  │ debitAccount()  │  │            │                            │
│  │ Joins outer TX  │  │            │  "No existing transaction  │
│  └─────────────────┘  │            │   found for transaction    │
└───────────────────────┘            │   marked with propagation  │
                                     │   'mandatory'"             │
                                     └───────────────────────────┘
```

**Use when:** Operations that should NEVER run standalone — they must always be part of a larger transaction (e.g., `debitAccount` must always be called within a transfer transaction).

---

### 3.3 REQUIRED vs REQUIRES_NEW — Detailed Comparison

| Aspect | REQUIRED | REQUIRES_NEW |
|--------|----------|--------------|
| **Existing TX?** | Join it | Suspend it, create new |
| **Rollback scope** | One rollback kills all | Each TX independent |
| **DB connections** | Shares connection | Needs **separate** connection |
| **Connection pool impact** | Low | **Higher** — can exhaust pool |
| **Commit timing** | All at once (outer) | Inner commits immediately |
| **Use case** | Standard business logic | Audit logs, sequence gen |
| **Deadlock risk** | Lower | **Higher** — two connections may lock each other |

**Connection pool warning with REQUIRES_NEW:**

```
Thread-1: processOrder() → TX1 (conn1)
    └── logAudit() → TX2 (conn2)    ← Needs SECOND connection!

If pool size = 10 and 10 threads hit processOrder():
  10 connections for outer TX + 10 needed for inner TX = 20
  Pool exhausted! → Deadlock (outer waiting for inner, pool has no connections)

Fix: Pool size >= 2 × max concurrent REQUIRES_NEW chains
```

---

### 3.4 When to Use NESTED?

**Scenario:** Order processing with optional discount application.

```java
@Service
public class OrderService {

    @Transactional
    public OrderResult processOrder(OrderRequest request) {
        Order order = createOrder(request);
        addItems(order, request.getItems());

        try {
            discountService.applyDiscount(order); // NESTED
        } catch (DiscountException e) {
            log.warn("Discount failed, continuing without discount", e);
            // Order still processes! Savepoint rolled back, but outer TX continues.
        }

        return completeOrder(order);
    }
}

@Service
public class DiscountService {

    @Transactional(propagation = Propagation.NESTED)
    public void applyDiscount(Order order) {
        // If this fails, only the discount part rolls back
        // The main order creation continues
    }
}
```

**NESTED vs REQUIRES_NEW for this scenario:**
- `NESTED`: Discount failure doesn't affect order. But if order fails, discount also rolls back. Single connection.
- `REQUIRES_NEW`: Discount commits independently. If order fails later, discount stays committed (orphaned discount!). Two connections.

---

### 3.5 What Happens When REQUIRES_NEW is Used Inside REQUIRED?

```java
@Service
public class OrderService {

    @Transactional // REQUIRED (default)
    public void processOrder(Order order) {
        orderRepo.save(order);               // Part of TX-A
        auditService.logEvent("ORDER_CREATED"); // Creates TX-B (REQUIRES_NEW)
        paymentService.charge(order);         // Back in TX-A

        // If charge() throws exception:
        // TX-A rolls back (order NOT saved)
        // TX-B already committed (audit log IS saved) ✓
    }
}
```

```
Timeline:
─────────────────────────────────────────────────────
TX-A: [BEGIN]──save()──[SUSPEND]──────────[RESUME]──charge()──[ROLLBACK]
TX-B:                   [BEGIN]──log()──[COMMIT]
─────────────────────────────────────────────────────
Result: Order not saved. Audit log saved.
```

---

### 3.6 What is Transaction Suspension?

When `REQUIRES_NEW` or `NOT_SUPPORTED` is encountered, the **current transaction is suspended**:

1. Current transaction's resources (connection, synchronizations) are **unbound** from the thread
2. A new transaction (or no transaction) begins
3. After the inner method completes, the outer transaction is **resumed** (resources rebound)

```
Thread's ThreadLocal state:

Before suspension:  TX-A resources bound
During suspension:  TX-B resources bound (TX-A saved aside)
After resumption:   TX-A resources restored
```

**Production concern:** Suspended transactions still hold their database connection. With `REQUIRES_NEW`, the thread holds **two connections simultaneously**. This can cause connection pool exhaustion under load.

---

### 3.7 Production Example: Wrong Propagation Bug

**The bug:** Notification records disappearing in production.

```java
@Service
public class OrderService {

    @Transactional
    public void placeOrder(OrderRequest request) {
        Order order = orderRepo.save(new Order(request));
        notificationService.createNotification(order); // REQUIRED (default)
        inventoryService.reserveStock(order);           // Throws exception!
    }
}

@Service
public class NotificationService {

    @Transactional // REQUIRED — joins the outer TX
    public void createNotification(Order order) {
        notificationRepo.save(new Notification(order));
        // This JOINS OrderService's transaction
    }
}
```

**What happened:** When `inventoryService.reserveStock()` threw an exception, the entire transaction rolled back — including the notification. Users never got order failure notifications.

**Fix:** Changed notification to `REQUIRES_NEW`:

```java
@Transactional(propagation = Propagation.REQUIRES_NEW)
public void createNotification(Order order) {
    notificationRepo.save(new Notification(order));
    // Now commits independently — survives outer TX rollback
}
```

---

### 3.8 Designing Transaction Boundaries in Microservices

In microservices, **transactions don't span service boundaries**. Each service owns its database and transaction.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│ Order Service│     │Payment Svc  │     │Inventory Svc│
│  ┌────────┐  │HTTP │  ┌────────┐ │HTTP │  ┌────────┐ │
│  │ TX-1   │──┼────→│  │ TX-2   │─┼────→│  │ TX-3   │ │
│  │ Local  │  │     │  │ Local  │ │     │  │ Local  │ │
│  └────────┘  │     │  └────────┘ │     │  └────────┘ │
│  OrderDB     │     │  PaymentDB  │     │  InventoryDB│
└─────────────┘     └─────────────┘     └─────────────┘

NO distributed transaction! Each has independent TX.
Use Saga pattern for coordination (see Section 8).
```

**Design principles:**
1. **Keep transactions local** — never try to span HTTP calls
2. **Use Saga for coordination** — choreography or orchestration
3. **Design for idempotency** — every operation must be safely retryable
4. **Use outbox pattern** — reliably publish events with local transactions

---

## 4. Isolation Levels (Critical for Concurrency)

### 4.1 What is Isolation in Transactions?

Isolation defines **how much one transaction can see of another's uncommitted changes**. Higher isolation = more correctness, but lower concurrency and performance.

```
     Low Isolation                          High Isolation
     ◄──────────────────────────────────────────────────►
     READ_UNCOMMITTED → READ_COMMITTED → REPEATABLE_READ → SERIALIZABLE

     More concurrency                       Less concurrency
     More anomalies                         Fewer anomalies
     Better performance                     Worse performance
```

---

### 4.2 All Isolation Levels

#### READ_UNCOMMITTED

```
TX-A: [BEGIN]── UPDATE balance=500 ───────────────── [ROLLBACK]
TX-B:                    [BEGIN]── READ balance ── [uses 500!] ── [COMMIT]
                                        │
                                   Reads 500 (uncommitted!)
                                   This is a DIRTY READ
```

- Can see **uncommitted** changes from other transactions
- Almost never used in production
- Only useful for approximate aggregates where accuracy doesn't matter

---

#### READ_COMMITTED

```
TX-A: [BEGIN]── UPDATE balance=500 ─── [COMMIT]
TX-B:                    [BEGIN]── READ(1000) ──────── READ(500) ── [COMMIT]
                                    │                    │
                               Sees old value       Sees new value
                               (committed: 1000)    (committed: 500)
                                    │
                               Non-repeatable read!
```

- Only sees **committed** data — no dirty reads
- Same query can return different results within one transaction
- **Default for PostgreSQL, Oracle, SQL Server**

---

#### REPEATABLE_READ

```
TX-A: [BEGIN]── UPDATE balance=500 ─── [COMMIT]
TX-B: [BEGIN]── READ(1000) ──────────────────── READ(1000) ── [COMMIT]
                    │                               │
               Sees 1000                       Still sees 1000
               (snapshot at TX start)          (consistent snapshot)
```

- Same query returns same results throughout the transaction
- No dirty reads, no non-repeatable reads
- **Phantom reads possible** (new rows inserted by other TX)
- **Default for MySQL (InnoDB)**

---

#### SERIALIZABLE

```
TX-A: [BEGIN]── SELECT * WHERE status='ACTIVE' ─── INSERT(new row) ─── [COMMIT]
TX-B: [BEGIN]── SELECT * WHERE status='ACTIVE' ─── [BLOCKED until TX-A commits]
                                                        │
                                                   No phantoms!
                                                   But significant performance cost
```

- Transactions execute as if they were serial (one after another)
- No anomalies at all — strongest guarantee
- **Significant performance penalty** — heavy locking or MVCC conflict detection
- Used for critical financial calculations, account reconciliation

---

### 4.3 Read Phenomena Explained

#### Dirty Read

```
Time  TX-A                        TX-B
─────────────────────────────────────────────────
 T1   UPDATE salary = 8000
      (not committed yet)
 T2                                SELECT salary
                                   → Returns 8000 (DIRTY!)
 T3   ROLLBACK
      (salary back to 5000)
 T4                                Uses 8000 for calculation
                                   → WRONG DATA USED!
```

**Impact:** Financial calculations based on data that never actually existed.

---

#### Non-Repeatable Read

```
Time  TX-A                        TX-B
─────────────────────────────────────────────────
 T1                                SELECT salary → 5000
 T2   UPDATE salary = 8000
      COMMIT
 T3                                SELECT salary → 8000
                                   (Different from T1!)
```

**Impact:** Report shows inconsistent totals when summing across related queries.

---

#### Phantom Read

```
Time  TX-A                          TX-B
─────────────────────────────────────────────────
 T1                                  SELECT COUNT(*) WHERE dept='SALES'
                                     → Returns 5
 T2   INSERT INTO employees
      (dept='SALES') COMMIT
 T3                                  SELECT COUNT(*) WHERE dept='SALES'
                                     → Returns 6 (phantom row!)
```

**Impact:** Pagination breaks, aggregate counts become inconsistent within a report.

---

### 4.4 Isolation Level vs Read Phenomena Matrix

| Isolation Level | Dirty Read | Non-Repeatable Read | Phantom Read | Performance |
|----------------|------------|---------------------|--------------|-------------|
| **READ_UNCOMMITTED** | Possible | Possible | Possible | Fastest |
| **READ_COMMITTED** | Prevented | Possible | Possible | Fast |
| **REPEATABLE_READ** | Prevented | Prevented | Possible* | Moderate |
| **SERIALIZABLE** | Prevented | Prevented | Prevented | Slowest |

*MySQL InnoDB's REPEATABLE_READ also prevents phantom reads via gap locking (next-key locking).

---

### 4.5 Default Isolation Levels

| Database | Default Isolation | Why |
|----------|------------------|-----|
| **MySQL (InnoDB)** | REPEATABLE_READ | Strong consistency with MVCC, gap locks prevent phantoms |
| **PostgreSQL** | READ_COMMITTED | Good balance of consistency and performance with MVCC |
| **Oracle** | READ_COMMITTED | Optimized for high concurrency OLTP workloads |
| **SQL Server** | READ_COMMITTED | (with READ_COMMITTED_SNAPSHOT option for MVCC) |
| **Spring** | DEFAULT (delegates to DB) | Uses whatever the DB defaults to |

**Production advice:**
- Rarely change the default isolation level
- If you need SERIALIZABLE for a specific operation, set it only on that method:

```java
@Transactional(isolation = Isolation.SERIALIZABLE)
public void reconcileAccounts() { ... }
```

- Changing isolation globally impacts every query in the application

---

### 4.6 How Isolation Affects Performance

```
                    Isolation Level
                          │
            ┌─────────────┼──────────────┐
            ▼             ▼              ▼
       Low (RC)     Medium (RR)     High (SER)
            │             │              │
    Fewer locks     More locks      Most locks
    Less memory     MVCC snapshots  Range locks
    Higher TPS      Moderate TPS    Lowest TPS
            │             │              │
            ▼             ▼              ▼
    Best for:       Best for:       Best for:
    High-traffic    Standard OLTP   Financial
    read-heavy      applications    reconciliation
    applications
```

| Isolation | Lock Overhead | Memory | Throughput | Deadlock Risk |
|-----------|--------------|--------|------------|---------------|
| READ_UNCOMMITTED | Minimal | Low | Highest | Very Low |
| READ_COMMITTED | Row-level | Moderate | High | Low |
| REPEATABLE_READ | Row + gap | Higher (snapshots) | Moderate | Medium |
| SERIALIZABLE | Range locks | Highest | Lowest | High |

---

## 5. Real Production Scenarios

### 5.1 What Happens If DB Is Slow but Transaction Is Open?

```
┌──────────────────────────────────────────────────────────┐
│  Thread-1: @Transactional processOrder()                  │
│                                                          │
│  [BEGIN TX] ─── [DB query: 50ms] ─── [Slow query: 30s]  │
│       │                                        │         │
│       │         Connection HELD for 30+ seconds│         │
│       │                                        │         │
│       └── Locks held ── Other threads BLOCKED ─┘         │
│                                                          │
│  HikariCP Connection Pool:                               │
│  [busy][busy][busy][busy][busy]...[waiting][waiting]     │
│        All 10 connections tied up!                        │
│                                                          │
│  Result: Application UNRESPONSIVE                        │
└──────────────────────────────────────────────────────────┘
```

**Consequences:**
1. **Connection pool exhaustion** — new requests queue up
2. **DB locks held longer** — other transactions wait/timeout
3. **Cascading timeouts** — upstream services timeout → retry → amplify load
4. **Memory pressure** — pending transactions hold Hibernate dirty-checking state

**Mitigation:**
```java
@Transactional(timeout = 10) // Fail fast
public void processOrder(Order order) { ... }
```
Plus: set HikariCP `connectionTimeout` and `maxLifetime`, add circuit breakers.

---

### 5.2 Why Should Transactions Be Kept Short?

```
Long Transaction:
[BEGIN]────────────────────── 30 seconds ──────────────────────[COMMIT]
   │                                                              │
   ├── Holds DB connection for 30s                                │
   ├── Holds row/table locks for 30s                              │
   ├── Prevents other TX from accessing same rows                 │
   ├── Increases deadlock probability                             │
   ├── Consumes undo/redo log space                               │
   └── Hibernate dirty check state grows in memory                │

Short Transaction:
[BEGIN]── 50ms ──[COMMIT]
   │                │
   └── Minimal impact on system
```

**Rule of thumb:** A transaction should complete in **< 1 second** for OLTP systems. If longer, redesign.

**Strategies for shortening:**
1. Move non-DB work (API calls, file I/O, computation) outside the transaction
2. Use batch processing with chunk-sized transactions
3. Pre-compute or cache data before entering the transaction
4. Use async processing for non-critical post-commit work

---

### 5.3 What Happens If External API Call Is Inside a Transaction?

```java
// DANGEROUS PATTERN
@Transactional
public void processOrder(Order order) {
    orderRepo.save(order);

    // External API call — blocks for 2-15 seconds
    PaymentResult result = paymentGateway.charge(order);  // HTTP call!

    if (result.isSuccess()) {
        order.setStatus(PAID);
        orderRepo.save(order);
    }
}
```

**Problems:**

```
DB Connection held during HTTP call:
[BEGIN TX]──save()──[HTTP call: 2-15s]──save()──[COMMIT]
                        │
                   DB connection IDLE
                   but still HELD
                   Locks still ACTIVE
                   Pool connection WASTED
```

1. **Connection held for seconds** during network I/O — waste
2. **Network timeout/failure** can leave transaction in uncertain state
3. **Retry on HTTP failure** may cause duplicate DB writes
4. **DB locks held** while waiting for external system

**Correct pattern:**

```java
public void processOrder(Order order) {
    // Step 1: Save order (short TX)
    orderService.saveOrder(order); // @Transactional - completes in ms

    // Step 2: Call external API (no TX)
    PaymentResult result = paymentGateway.charge(order);

    // Step 3: Update status (short TX)
    if (result.isSuccess()) {
        orderService.markAsPaid(order.getId()); // @Transactional - completes in ms
    } else {
        orderService.markAsFailed(order.getId());
    }
}
```

Or even better — use the **Outbox Pattern** for guaranteed delivery.

---

### 5.4 How to Handle Transactions in Distributed Systems?

```
┌─────────────────────────────────────────────────────────────┐
│         Distributed Transaction Strategies                   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. Saga Pattern (PREFERRED)                                │
│     ┌──────┐   ┌──────┐   ┌──────┐   ┌──────┐            │
│     │Order │──→│Payment│──→│Invent│──→│Notif │            │
│     │Create│   │Charge │   │Reserve│   │Send  │            │
│     └──┬───┘   └──┬───┘   └──┬───┘   └──────┘            │
│        │          │          │                              │
│     Compensate Compensate Compensate                        │
│     (Cancel)   (Refund)   (Release)                        │
│                                                             │
│  2. Outbox Pattern (for event reliability)                  │
│     ┌──────────────────┐                                    │
│     │ Service DB        │    ┌─────────┐    ┌──────────┐  │
│     │ ┌──────┐ ┌─────┐ │───→│ Debezium│───→│  Kafka   │  │
│     │ │Orders│ │Outbox│ │    │  CDC    │    │          │  │
│     │ └──────┘ └─────┘ │    └─────────┘    └──────────┘  │
│     └──────────────────┘                                    │
│     Single local TX writes both order + event               │
│                                                             │
│  3. Try-Confirm/Cancel (TCC)                                │
│     Try: Reserve resources (soft lock)                      │
│     Confirm: Finalize reservation                           │
│     Cancel: Release reservation                             │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

See [Section 8](#8-advanced--architecture-level) for detailed Saga and distributed transaction patterns.

---

### 5.5 Why @Transactional Should Not Be Used on Controller Layer?

```java
// BAD — Transaction in controller
@RestController
public class OrderController {

    @Transactional  // DON'T DO THIS
    @PostMapping("/orders")
    public ResponseEntity<Order> createOrder(@RequestBody OrderRequest req) {
        Order order = orderService.createOrder(req);
        return ResponseEntity.ok(order);
    }
}
```

**Why this is problematic:**

| Issue | Explanation |
|-------|------------|
| **Transaction too long** | Includes HTTP request parsing, serialization, response writing |
| **View rendering inside TX** | JSON serialization happens while TX is open → lazy loading works but connection held |
| **Separation of concerns** | Controller handles HTTP concerns, not business transaction boundaries |
| **Error handling confusion** | HTTP error handling mixed with TX rollback logic |
| **Testing difficulty** | Can't test business logic transaction behavior independently |
| **Connection held during response** | DB connection held while writing HTTP response to client |

**Correct layering:**

```
Controller Layer    → HTTP concerns, validation, response mapping
        │
Service Layer       → @Transactional here — business transaction boundaries
        │
Repository Layer    → Data access (participates in service's TX)
```

---

## 6. Spring + JPA Specific

### 6.1 What is PlatformTransactionManager?

`PlatformTransactionManager` is the **central interface** in Spring's transaction abstraction. All transaction managers implement it.

```java
public interface PlatformTransactionManager {
    TransactionStatus getTransaction(TransactionDefinition definition);
    void commit(TransactionStatus status);
    void rollback(TransactionStatus status);
}
```

```
PlatformTransactionManager (interface)
        │
        ├── DataSourceTransactionManager (JDBC)
        ├── JpaTransactionManager (JPA/Hibernate)
        ├── JtaTransactionManager (distributed/XA)
        ├── HibernateTransactionManager (native Hibernate)
        └── ChainedTransactionManager (multiple resources, best-effort)
```

Spring Boot **auto-configures** the appropriate transaction manager:
- Has `spring-boot-starter-data-jpa`? → `JpaTransactionManager`
- Has only `spring-boot-starter-jdbc`? → `DataSourceTransactionManager`
- Has JTA provider (Atomikos, Narayana)? → `JtaTransactionManager`

---

### 6.2 JpaTransactionManager vs DataSourceTransactionManager

| Aspect | JpaTransactionManager | DataSourceTransactionManager |
|--------|----------------------|------------------------------|
| **Works with** | JPA (`EntityManager`) | Plain JDBC (`DataSource`) |
| **Manages** | EntityManager + Connection | Connection only |
| **Hibernate features** | Full support (dirty check, cache, lazy loading) | No Hibernate awareness |
| **Persistence context** | Creates and manages it | Not applicable |
| **Flush behavior** | Controls Hibernate flush | Not applicable |
| **Use when** | Using Spring Data JPA / Hibernate | Using `JdbcTemplate` only |
| **Can use JdbcTemplate?** | Yes — shares same connection | Yes |

**Important:** `JpaTransactionManager` can also manage plain JDBC operations if they use the same `DataSource`. You can mix JPA and `JdbcTemplate` in the same transaction.

```java
@Transactional
public void mixedOperations() {
    // JPA operation — uses EntityManager
    entityManager.persist(new Order());

    // JDBC operation — uses same underlying Connection
    jdbcTemplate.update("INSERT INTO audit_log ...");

    // Both participate in the SAME transaction ✓
}
```

---

### 6.3 What is Persistence Context?

The **persistence context** is Hibernate's first-level cache — a map of managed entities for the current session/transaction.

```
┌───────────────────────────────────────────────────────────┐
│                  Persistence Context                       │
│                  (per Transaction)                         │
│                                                           │
│   Entity Map:                                             │
│   ┌──────────────────────────────────────────────┐       │
│   │ Key: (Entity Class, ID)  → Value: Entity     │       │
│   │                                              │       │
│   │ (User.class, 1)         → User{id=1, "Ali"} │       │
│   │ (Order.class, 42)       → Order{id=42, ...}  │       │
│   │ (Product.class, 7)      → Product{id=7, ...} │       │
│   └──────────────────────────────────────────────┘       │
│                                                           │
│   Snapshot Map (for dirty checking):                      │
│   ┌──────────────────────────────────────────────┐       │
│   │ (User.class, 1) → {id=1, name="Ali"} (orig) │       │
│   └──────────────────────────────────────────────┘       │
│                                                           │
│   Entity States:                                          │
│   ┌──────────┬───────────┬──────────┬──────────┐         │
│   │ Transient│ Managed   │ Detached │ Removed  │         │
│   │ (new)    │ (in PC)   │ (TX done)│ (deleted)│         │
│   └──────────┴───────────┴──────────┴──────────┘         │
└───────────────────────────────────────────────────────────┘
```

**Entity lifecycle states:**

```
         new Entity()              persist()
Transient ─────────────────→ Managed (in Persistence Context)
                                    │
                              ┌─────┼──────┐
                              │     │      │
                           flush  remove  TX ends
                           (SQL)   │    (detach)
                              │     │      │
                              ▼     ▼      ▼
                           Still  Removed  Detached
                           Managed         (no longer tracked)
```

**Key behaviors:**
- **Identity guarantee:** `em.find(User.class, 1) == em.find(User.class, 1)` (same object)
- **Repeatable reads:** Within a TX, same query returns same managed objects
- **Write-behind:** Changes are not immediately sent to DB — collected and flushed

---

### 6.4 What is Dirty Checking?

Dirty checking is Hibernate's mechanism to **automatically detect which entities have been modified** and generate the appropriate UPDATE SQL at flush time.

```java
@Transactional
public void updateUserName(Long userId, String newName) {
    User user = userRepository.findById(userId).orElseThrow();
    user.setName(newName);
    // NO explicit save() call needed!
    // Hibernate detects the change via dirty checking
}
```

**How it works internally:**

```
1. Load entity → Hibernate takes SNAPSHOT of all field values
   user = {id=1, name="Ali", email="ali@x.com"}
   snapshot = {id=1, name="Ali", email="ali@x.com"}

2. You modify:
   user.setName("Ahmed")

3. At flush time, Hibernate compares:
   current:  {id=1, name="Ahmed", email="ali@x.com"}
   snapshot: {id=1, name="Ali",   email="ali@x.com"}
                         ↑
                    DIFFERENT! → Generate UPDATE

4. SQL generated:
   UPDATE users SET name='Ahmed' WHERE id=1
```

**Performance implications:**
- Every managed entity has a snapshot copy in memory → **doubles memory** per entity
- Comparison happens field-by-field at flush → CPU cost with many entities
- `readOnly = true` disables dirty checking → significant memory savings for read operations

**Optimization for bulk reads:**

```java
@Transactional(readOnly = true)
public List<Report> generateReport() {
    // readOnly = true → FlushMode.MANUAL → NO dirty checking
    // 10,000 entities loaded without snapshot copies
    return reportRepository.findAll();
}
```

---

### 6.5 When Does Hibernate Flush Data?

Flushing = synchronizing persistence context state with the database (executing pending SQL).

```
Flush Triggers:
─────────────────────────────────────────────────────
1. Transaction COMMIT
   @Transactional method returns → auto-flush → commit

2. Before JPQL/HQL query execution
   If you query the same table that has pending changes,
   Hibernate flushes first to ensure query sees latest data

3. Explicit flush
   entityManager.flush()

4. Native SQL query (depends on FlushMode)
   Native queries may NOT trigger auto-flush — be careful!
```

**FlushMode options:**

| FlushMode | Behavior | When Used |
|-----------|----------|-----------|
| `AUTO` (default) | Flush before queries + at commit | Normal operations |
| `COMMIT` | Flush only at commit | Performance optimization — skip pre-query flush |
| `MANUAL` | Never auto-flush | `readOnly = true` transactions |
| `ALWAYS` | Flush before every query | Rarely used, expensive |

**Gotcha with native queries:**

```java
@Transactional
public void riskyOperation() {
    User user = userRepo.findById(1L).get();
    user.setName("Updated");
    // Change is in persistence context, NOT in DB yet

    // Native query — does NOT trigger auto-flush!
    List<User> result = em.createNativeQuery("SELECT * FROM users WHERE id = 1", User.class)
                          .getResultList();
    // result may show OLD name!

    // Fix: manually flush before native query
    em.flush();
}
```

---

## 7. Common Bugs (Very Important for Interview)

### 7.1 Self-Invocation Issue in Spring Transactions

Already covered in [Section 1.8](#18-self-invocation-problem). Here's a quick production detection checklist:

```
Symptom: @Transactional(propagation = REQUIRES_NEW) not creating new TX

Debug checklist:
 ✓ Is the method being called from SAME class? → Self-invocation!
 ✓ Is the method public?
 ✓ Is the bean Spring-managed (@Service, @Component)?
 ✓ Is @EnableTransactionManagement present? (auto in Spring Boot)

Quick test:
  Add logging in TransactionSynchronizationManager:
  log.info("TX active: {}", TransactionSynchronizationManager.isActualTransactionActive());
  log.info("TX name: {}", TransactionSynchronizationManager.getCurrentTransactionName());
```

---

### 7.2 Transaction Not Working Due to Missing Proxy

**Scenario 1: Bean not managed by Spring**

```java
// NOT a Spring bean — no proxy created
public class OrderService {

    @Transactional
    public void createOrder(Order order) {
        // @Transactional is IGNORED — no proxy wrapping this class
    }
}

// Fix: Add @Service or @Component
@Service
public class OrderService { ... }
```

**Scenario 2: Creating instance manually**

```java
@Configuration
public class AppConfig {

    @Bean
    public OrderService orderService() {
        return new OrderService(); // Spring wraps with proxy ✓
    }
}

// But if somewhere you do:
OrderService svc = new OrderService(); // NO proxy — @Transactional won't work!
```

**Scenario 3: Final class or method**

```java
@Service
public final class OrderService { // CGLIB cannot subclass final class!

    @Transactional
    public void createOrder(Order order) { ... } // TX won't work
}

// Fix: Remove 'final' or use interface-based proxy
```

---

### 7.3 LazyInitializationException Due to Transaction Boundary

```java
@Service
public class OrderService {

    @Transactional
    public Order getOrder(Long id) {
        return orderRepo.findById(id).orElseThrow();
        // Transaction ENDS here
        // Persistence context CLOSED
    }
}

@RestController
public class OrderController {

    @GetMapping("/orders/{id}")
    public OrderDTO getOrder(@PathVariable Long id) {
        Order order = orderService.getOrder(id);
        order.getItems().size(); // LazyInitializationException!
        // Accessing lazy collection AFTER TX closed
        return mapper.toDTO(order);
    }
}
```

```
Timeline:
Service method          Controller
[TX BEGIN]──findById()──[TX END/PC CLOSED]──getItems()──BOOM!
                                                │
                                    LazyInitializationException
                                    "could not initialize proxy
                                     - no Session"
```

**Solutions (ranked by preference):**

| Solution | Approach | Trade-off |
|----------|----------|-----------|
| **1. Fetch in query** | `@Query("SELECT o FROM Order o JOIN FETCH o.items")` | Best — explicit, no N+1 |
| **2. @EntityGraph** | `@EntityGraph(attributePaths = {"items"})` | Clean, declarative |
| **3. DTO projection** | Return DTO from service, not entity | Cleanest separation |
| **4. Open Session in View** | `spring.jpa.open-in-view=true` (default!) | **Avoid in production** — keeps TX open during view rendering |

**Production recommendation:** Disable Open Session in View and use explicit fetching:

```yaml
spring:
  jpa:
    open-in-view: false  # Disable OSIV
```

---

### 7.4 Long-Running Transaction Causing DB Locks

```
Thread-1: @Transactional (takes 30 seconds)
─────────────────────────────────────────────────────
[BEGIN]──SELECT FOR UPDATE (row id=1)──[processing 30s]──[COMMIT]
              │
              │  ROW LOCKED for 30 seconds
              │
Thread-2: @Transactional
─────────────────────────────────────────────────────
         [BEGIN]──UPDATE (row id=1)──[WAITING...]──[TIMEOUT!]
                                         │
                              Lock wait timeout exceeded
```

**Detection:**

```sql
-- MySQL: Check current locks
SELECT * FROM information_schema.INNODB_LOCKS;
SELECT * FROM information_schema.INNODB_LOCK_WAITS;

-- PostgreSQL: Check blocking queries
SELECT pid, query, state, wait_event_type
FROM pg_stat_activity
WHERE state = 'active';

-- Find blocking locks
SELECT blocked.pid AS blocked_pid,
       blocking.pid AS blocking_pid,
       blocked.query AS blocked_query
FROM pg_locks blocked
JOIN pg_locks blocking ON blocked.locktype = blocking.locktype
WHERE NOT blocked.granted;
```

**Prevention strategies:**
1. Keep transactions under 1 second
2. Set `@Transactional(timeout = 10)`
3. Use optimistic locking (`@Version`) instead of `SELECT FOR UPDATE`
4. Process in batches with separate transactions per batch
5. Move non-DB work outside transactions

---

### 7.5 Deadlock in Transactions — How to Handle

```
Thread-1:                          Thread-2:
[BEGIN]                            [BEGIN]
Lock Row A ✓                       Lock Row B ✓
   │                                  │
   ├── Wait for Row B (held by T2)    ├── Wait for Row A (held by T1)
   │       ↓                          │       ↓
   │   BLOCKED                        │   BLOCKED
   │                                  │
   └──────── DEADLOCK! ──────────────┘
             │
     DB detects and kills one TX
     (DeadlockLoserDataAccessException)
```

**Prevention strategies:**

```
1. Consistent lock ordering
   Always lock resources in the same order (e.g., by ID ascending)

   // BAD
   Thread-1: lock(accountA) → lock(accountB)
   Thread-2: lock(accountB) → lock(accountA)

   // GOOD
   Both:     lock(min(A,B)) → lock(max(A,B))
```

```java
@Transactional
public void transfer(Long fromId, Long toId, BigDecimal amount) {
    Long first = Math.min(fromId, toId);
    Long second = Math.max(fromId, toId);

    Account firstAccount = accountRepo.findByIdWithLock(first);
    Account secondAccount = accountRepo.findByIdWithLock(second);

    // Now proceed with debit/credit
}
```

```
2. Use optimistic locking (@Version) to avoid explicit locks

3. Keep transactions short → smaller lock window

4. Set lock timeout:
   @QueryHint(name = "javax.persistence.lock.timeout", value = "3000")

5. Retry on deadlock:
   @Retryable(value = DeadlockLoserDataAccessException.class, maxAttempts = 3)
```

---

## 8. Advanced / Architecture Level

### 8.1 How Transactions Work in Microservices

In a monolith, a single `@Transactional` covers everything. In microservices, each service has its own database and local transaction.

```
Monolith (single TX):
┌──────────────────────────────────────────┐
│  @Transactional                          │
│  createOrder() + chargePayment() +       │
│  reserveInventory() + sendNotification() │
│  All in ONE transaction, ONE database    │
└──────────────────────────────────────────┘

Microservices (local TX per service):
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Order Svc│     │Payment Svc│    │Inventory │
│ Local TX │────→│ Local TX  │───→│ Local TX │
│ OrderDB  │     │ PaymentDB │    │ InvDB    │
└──────────┘     └──────────┘     └──────────┘
     │                │                │
     └── No shared transaction boundary ──┘
```

**Key principles:**
1. **Database per service** — no shared DB
2. **Local transactions only** — `@Transactional` within each service
3. **Eventual consistency** — accept that cross-service consistency is async
4. **Compensating actions** — undo completed steps on failure

---

### 8.2 What is Distributed Transaction?

A distributed transaction spans multiple databases or services, trying to maintain ACID across all of them.

```
Distributed Transaction (2PC):
┌────────────────────────────────────────────────────┐
│  Transaction Coordinator                            │
│       │                                            │
│       ├── PREPARE → Service A (DB-1) → VOTE YES   │
│       ├── PREPARE → Service B (DB-2) → VOTE YES   │
│       ├── PREPARE → Service C (DB-3) → VOTE YES   │
│       │                                            │
│       │   All voted YES?                           │
│       │                                            │
│       ├── COMMIT → Service A                       │
│       ├── COMMIT → Service B                       │
│       └── COMMIT → Service C                       │
└────────────────────────────────────────────────────┘
```

---

### 8.3 Why 2PC (Two Phase Commit) Is Avoided

| Problem | Impact |
|---------|--------|
| **Synchronous blocking** | All participants lock resources until coordinator decides |
| **Single point of failure** | Coordinator crash → all participants stuck in PREPARED state |
| **High latency** | Multiple round trips across network |
| **Reduced availability** | If any participant is down, entire TX fails |
| **Not scalable** | Holding locks across services kills throughput |
| **Network partitions** | What if coordinator can't reach a participant after PREPARE? |

```
2PC Failure Scenario:
─────────────────────────────────────────────────
Coordinator: PREPARE → all vote YES
Coordinator: COMMIT → sent to A ✓, sent to B ✓
Coordinator: COMMIT → network failure to C ✗
             │
             C is stuck in PREPARED state
             Locks held indefinitely!
             Manual intervention required
```

**When 2PC IS acceptable:**
- Within a single application using XA transactions across 2 databases (e.g., DB + JMS)
- Low-throughput back-office systems where correctness > performance

---

### 8.4 Saga Pattern

The Saga pattern replaces distributed transactions with a sequence of local transactions + compensating actions.

#### Choreography-based Saga (Event-driven)

```
┌──────────┐  OrderCreated  ┌──────────┐  PaymentDone  ┌──────────┐
│  Order   │ ──────────────→│ Payment  │──────────────→│Inventory │
│  Service │                │ Service  │               │ Service  │
└────┬─────┘                └────┬─────┘               └────┬─────┘
     │                           │                          │
     │  PaymentFailed            │  InventoryFailed         │
     │◄──────────────────────────│◄─────────────────────────│
     │                           │                          │
  Cancel                      Refund                     Release
  Order                       Payment                    Stock
(compensate)               (compensate)              (compensate)
```

#### Orchestration-based Saga (Central coordinator)

```
                    ┌────────────────────┐
                    │   Saga Orchestrator │
                    │   (Order Saga)      │
                    └──┬──────┬──────┬───┘
                       │      │      │
              Step 1   │      │      │  Step 3
           Create Order│      │      │  Reserve Inventory
                       ▼      │      ▼
                  ┌────────┐  │  ┌──────────┐
                  │ Order  │  │  │Inventory │
                  │ Service│  │  │ Service  │
                  └────────┘  │  └──────────┘
                              │
                     Step 2   │
                  Charge      │
                  Payment     │
                              ▼
                         ┌──────────┐
                         │ Payment  │
                         │ Service  │
                         └──────────┘

On any step failure → Orchestrator triggers compensations for completed steps
```

| Aspect | Choreography | Orchestration |
|--------|-------------|---------------|
| **Coupling** | Loose (event-driven) | Tighter (orchestrator knows all steps) |
| **Complexity** | Simple for 2-3 services | Better for 4+ services |
| **Debugging** | Hard (distributed flow) | Easier (central coordinator has full state) |
| **Single point of failure** | None | Orchestrator (but stateless, so scalable) |
| **Recommended for** | Simple workflows | Complex business processes |

---

### 8.5 Eventual Consistency vs Strong Consistency

| Aspect | Strong Consistency | Eventual Consistency |
|--------|-------------------|---------------------|
| **Guarantee** | Read always returns latest write | Read may return stale data temporarily |
| **Availability** | Lower (blocks during sync) | Higher (always responds) |
| **Latency** | Higher (synchronous) | Lower (asynchronous) |
| **Implementation** | 2PC, distributed locks | Events, Saga, CQRS |
| **Use when** | Financial transfers, inventory count | Social media feeds, analytics, notifications |

```
Strong Consistency:
Write → [Sync to all replicas] → Read (guaranteed latest)
         ↑ Blocks until all confirm

Eventual Consistency:
Write → [Async replication] → Read (might be stale)
         ↑ Returns immediately
         Replicas catch up eventually (ms to seconds)
```

**Production reality:** Most microservice systems use **eventual consistency with compensation**. The key is making the "eventual" window as short as possible and handling inconsistency gracefully in the UI.

---

## 9. Scenario-Based Questions (Interview Gold)

### 9.1 If Payment Succeeds but Order Fails — How to Fix?

```
Problem:
  1. Create Order      → SUCCESS
  2. Charge Payment    → SUCCESS
  3. Update Order Status → FAILS (DB error)

  Money charged but order not confirmed!
```

**Solution architecture:**

```
┌─────────────────────────────────────────────────────────┐
│  Approach 1: Saga with Compensation                      │
│                                                         │
│  Step 1: Create Order (PENDING)   → Local TX            │
│  Step 2: Charge Payment           → API call            │
│  Step 3: Confirm Order (CONFIRMED)→ Local TX            │
│                                                         │
│  If Step 3 fails:                                       │
│    → Compensate Step 2: Refund payment                  │
│    → Compensate Step 1: Cancel order                    │
└─────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────┐
│  Approach 2: Outbox Pattern (PREFERRED)                  │
│                                                         │
│  @Transactional                                         │
│  void createOrder(request):                             │
│    order = orderRepo.save(new Order(PENDING))           │
│    outboxRepo.save(new OutboxEvent("PAYMENT_REQUESTED"))│
│    // Single local TX — atomic!                         │
│                                                         │
│  Outbox Poller/CDC reads event → triggers Payment Service│
│  Payment Service responds → another event → confirm order│
└─────────────────────────────────────────────────────────┘
```

```java
@Service
public class OrderSagaOrchestrator {

    @Transactional
    public OrderResult processOrder(OrderRequest request) {
        Order order = orderService.createOrder(request); // PENDING

        try {
            PaymentResult payment = paymentService.charge(request);

            if (payment.isSuccess()) {
                orderService.confirmOrder(order.getId());
                return OrderResult.success(order);
            } else {
                orderService.cancelOrder(order.getId());
                return OrderResult.paymentFailed();
            }
        } catch (Exception e) {
            // Payment status unknown — need reconciliation
            orderService.markForReconciliation(order.getId());
            reconciliationService.scheduleCheck(order.getId());
            return OrderResult.pendingReconciliation();
        }
    }
}
```

---

### 9.2 How Do You Ensure Idempotency in Transactions?

**Idempotency** = executing the same operation multiple times produces the same result as executing it once.

```
Without idempotency:
  Request: Transfer $100 (retry due to timeout)
  Execution 1: Debit $100 ✓
  Execution 2: Debit $100 ✓ (DOUBLE CHARGE!)

With idempotency:
  Request: Transfer $100 (idempotencyKey = "txn-abc-123")
  Execution 1: Debit $100, store key ✓
  Execution 2: Key exists → return previous result ✓ (NO double charge)
```

**Implementation pattern:**

```java
@Service
public class PaymentService {

    @Transactional
    public PaymentResult processPayment(String idempotencyKey, PaymentRequest request) {
        // Check if already processed
        Optional<Payment> existing = paymentRepo.findByIdempotencyKey(idempotencyKey);
        if (existing.isPresent()) {
            return existing.get().toResult(); // Return previous result
        }

        // Process new payment
        Payment payment = new Payment(idempotencyKey, request);
        payment.execute();
        paymentRepo.save(payment);

        return payment.toResult();
    }
}
```

```sql
-- Idempotency key table
CREATE TABLE idempotency_keys (
    idempotency_key VARCHAR(255) PRIMARY KEY,
    response_payload JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '24 hours'
);

-- UNIQUE constraint prevents race condition
CREATE UNIQUE INDEX idx_idempotency_key ON idempotency_keys(idempotency_key);
```

---

### 9.3 How to Handle Retry Safely in Transactions?

```java
@Service
public class OrderService {

    @Retryable(
        value = {OptimisticLockException.class, DeadlockLoserDataAccessException.class},
        maxAttempts = 3,
        backoff = @Backoff(delay = 100, multiplier = 2)
    )
    @Transactional
    public Order placeOrder(OrderRequest request) {
        // Each retry gets a FRESH transaction
        // Previous failed TX was already rolled back
        return createAndSaveOrder(request);
    }

    @Recover
    public Order recoverPlaceOrder(Exception e, OrderRequest request) {
        log.error("Order placement failed after retries", e);
        throw new OrderProcessingException("Unable to place order", e);
    }
}
```

**Critical rules for safe retry:**

```
1. @Retryable MUST be on the OUTER method (not inside @Transactional)
   Why? If TX fails, the entire TX must be rolled back and retried fresh

2. Ensure idempotency — retried operations must be safe to repeat

3. Use exponential backoff — don't hammer a struggling database

4. Set max attempts — infinite retry = infinite resource consumption

5. Only retry transient errors:
   ✓ OptimisticLockException (someone else updated)
   ✓ DeadlockLoserDataAccessException (deadlock victim)
   ✓ CannotAcquireLockException (lock timeout)
   ✗ ConstraintViolationException (data issue, won't fix on retry)
   ✗ DataIntegrityViolationException (business logic error)
```

```
Retry flow:
Attempt 1: [BEGIN TX]──execute()──[EXCEPTION]──[ROLLBACK]──wait 100ms
Attempt 2: [BEGIN TX]──execute()──[EXCEPTION]──[ROLLBACK]──wait 200ms
Attempt 3: [BEGIN TX]──execute()──[SUCCESS]──[COMMIT] ✓
```

---

### 9.4 How to Avoid Double Booking Problem?

**Problem:** Two users try to book the last seat simultaneously.

```
User A: SELECT seats WHERE available=true → Seat 1A available
User B: SELECT seats WHERE available=true → Seat 1A available
User A: UPDATE seat 1A SET available=false → SUCCESS
User B: UPDATE seat 1A SET available=false → SUCCESS (DOUBLE BOOKING!)
```

**Solutions:**

#### Solution 1: Optimistic Locking (Preferred for low contention)

```java
@Entity
public class Seat {
    @Id
    private Long id;

    @Version
    private Long version; // Optimistic lock

    private boolean available;
}

@Transactional
public BookingResult bookSeat(Long seatId, Long userId) {
    Seat seat = seatRepo.findById(seatId).orElseThrow();
    if (!seat.isAvailable()) {
        return BookingResult.notAvailable();
    }
    seat.setAvailable(false);
    seat.setBookedBy(userId);
    seatRepo.save(seat); // Throws OptimisticLockException if version mismatch
    return BookingResult.success();
}
```

```
User A: SELECT seat (version=1) → UPDATE SET available=false, version=2 ✓
User B: SELECT seat (version=1) → UPDATE WHERE version=1 → 0 rows updated!
         → OptimisticLockException → retry or show "seat taken"
```

#### Solution 2: Pessimistic Locking (For high contention)

```java
@Repository
public interface SeatRepository extends JpaRepository<Seat, Long> {

    @Lock(LockModeType.PESSIMISTIC_WRITE)
    @Query("SELECT s FROM Seat s WHERE s.id = :id")
    Optional<Seat> findByIdWithLock(@Param("id") Long id);
}
```

```
User A: SELECT ... FOR UPDATE (locks row)
User B: SELECT ... FOR UPDATE (WAITS for A to release)
User A: UPDATE → COMMIT → lock released
User B: Now reads updated row → sees seat taken → handles gracefully
```

#### Solution 3: Database Unique Constraint (Simplest)

```sql
CREATE UNIQUE INDEX idx_booking_seat ON bookings(seat_id, event_id)
    WHERE status = 'CONFIRMED';
```

```java
@Transactional
public BookingResult bookSeat(Long seatId, Long eventId, Long userId) {
    try {
        bookingRepo.save(new Booking(seatId, eventId, userId, CONFIRMED));
        return BookingResult.success();
    } catch (DataIntegrityViolationException e) {
        return BookingResult.alreadyBooked();
    }
}
```

---

### 9.5 Design a Transaction Strategy for Booking System (IRCTC/Airline)

```
┌─────────────────────────────────────────────────────────────────┐
│                    Booking System Architecture                    │
│                                                                 │
│  Phase 1: Search & Selection (No TX needed)                     │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Read-only queries against read replicas                  │   │
│  │ Cache flight/seat availability in Redis                  │   │
│  │ @Transactional(readOnly=true) on read-replica DataSource │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│  Phase 2: Seat Hold (Short TX + TTL)                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ @Transactional(timeout=5)                                │   │
│  │ SELECT seat FOR UPDATE SKIP LOCKED                       │   │
│  │ Mark seat as HELD (with expiry = now + 10 minutes)       │   │
│  │ Return hold_token to user                                │   │
│  │ Background job releases expired holds                    │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│  Phase 3: Payment (Outside TX — Saga)                           │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Call Payment Gateway (NOT inside DB transaction)         │   │
│  │ Store payment_intent_id for reconciliation               │   │
│  │ On success → trigger confirmation                        │   │
│  │ On failure → release hold                                │   │
│  │ On timeout → schedule reconciliation check               │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│  Phase 4: Confirm Booking (Short TX)                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ @Transactional(timeout=5, isolation=SERIALIZABLE)        │   │
│  │ Verify hold is still valid + payment confirmed           │   │
│  │ Update seat: HELD → BOOKED                               │   │
│  │ Create booking record                                    │   │
│  │ Write to outbox: BOOKING_CONFIRMED event                 │   │
│  └─────────────────────────────────────────────────────────┘   │
│                           │                                     │
│  Phase 5: Post-Booking (Async, Eventually Consistent)           │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Outbox event → Kafka → Multiple consumers:               │   │
│  │   • Email confirmation                                   │   │
│  │   • SMS notification                                     │   │
│  │   • Loyalty points credit                                │   │
│  │   • Analytics update                                     │   │
│  │ Each consumer has its own local transaction               │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**Key design decisions:**

```
┌───────────────────────┬────────────────────────────────────────────┐
│ Decision              │ Rationale                                  │
├───────────────────────┼────────────────────────────────────────────┤
│ SKIP LOCKED for holds │ No blocking — users don't wait for others │
│ TTL on holds          │ Prevents permanent seat locks              │
│ Payment outside TX    │ Don't hold DB connection during HTTP call  │
│ Outbox for events     │ Guaranteed event delivery with local TX    │
│ SERIALIZABLE for      │ Prevent double-booking at confirmation     │
│ confirmation          │                                            │
│ Idempotency keys      │ Safe retry for payment + confirmation      │
│ Reconciliation job    │ Handle payment timeouts / edge cases       │
└───────────────────────┴────────────────────────────────────────────┘
```

```java
@Service
public class BookingService {

    @Transactional(timeout = 5)
    public HoldResult holdSeat(Long flightId, Long seatId, Long userId) {
        Seat seat = seatRepo.findAvailableWithSkipLock(seatId)
            .orElseThrow(() -> new SeatUnavailableException());

        seat.setStatus(SeatStatus.HELD);
        seat.setHeldBy(userId);
        seat.setHoldExpiry(Instant.now().plusSeconds(600));
        seatRepo.save(seat);

        String holdToken = UUID.randomUUID().toString();
        holdRepo.save(new Hold(holdToken, flightId, seatId, userId));

        return new HoldResult(holdToken, seat.getHoldExpiry());
    }

    public BookingResult confirmBooking(String holdToken, String paymentId) {
        PaymentStatus status = paymentGateway.verify(paymentId);
        if (status != PaymentStatus.SUCCESS) {
            releaseHold(holdToken);
            return BookingResult.paymentFailed();
        }
        return finalizeBooking(holdToken, paymentId);
    }

    @Transactional(timeout = 5, isolation = Isolation.SERIALIZABLE)
    public BookingResult finalizeBooking(String holdToken, String paymentId) {
        Hold hold = holdRepo.findByToken(holdToken)
            .orElseThrow(() -> new HoldExpiredException());

        if (hold.isExpired()) {
            throw new HoldExpiredException();
        }

        Seat seat = seatRepo.findById(hold.getSeatId()).orElseThrow();
        seat.setStatus(SeatStatus.BOOKED);

        Booking booking = new Booking(hold, paymentId);
        bookingRepo.save(booking);

        outboxRepo.save(OutboxEvent.bookingConfirmed(booking));

        return BookingResult.success(booking);
    }
}
```

---

## 10. Bonus — Senior Level Internals

### 10.1 How Spring Handles Nested Transactions Internally

```
@Transactional(propagation = NESTED) triggers:

1. AbstractPlatformTransactionManager.handleExistingTransaction()
      │
2. Check if nested TX supported (useSavepointForNestedTransaction)
      │
3. Create Savepoint:
      connection.setSavepoint("SAVEPOINT_1")
      │
4. Execute nested method
      │
5a. Success → connection.releaseSavepoint("SAVEPOINT_1")
      │
5b. Exception → connection.rollback(savepoint)
              → Outer TX continues (not rolled back)
```

```
Database-level view:

BEGIN;                          -- Outer TX
INSERT INTO orders ...;         -- Outer operation
SAVEPOINT sp1;                  -- Nested TX start
  INSERT INTO discounts ...;    -- Nested operation
  -- If fails:
  ROLLBACK TO sp1;              -- Only nested rolled back
  -- If succeeds:
  RELEASE SAVEPOINT sp1;        -- Nested merged with outer
INSERT INTO audit_log ...;      -- Outer continues
COMMIT;                         -- Everything committed together
```

**Key internal detail:** The nested transaction does NOT get its own database connection. It shares the outer transaction's connection and uses database savepoints for partial rollback capability.

---

### 10.2 What is Transaction Synchronization?

Transaction synchronization allows you to **register callbacks** that execute at specific points in the transaction lifecycle.

```java
@Transactional
public void createOrder(OrderRequest request) {
    Order order = orderRepo.save(new Order(request));

    // Register callback: runs AFTER commit
    TransactionSynchronizationManager.registerSynchronization(
        new TransactionSynchronization() {

            @Override
            public void afterCommit() {
                // Safe to send notifications here
                // TX is already committed — data is persisted
                eventPublisher.publish(new OrderCreatedEvent(order.getId()));
            }

            @Override
            public void afterCompletion(int status) {
                if (status == STATUS_ROLLED_BACK) {
                    log.warn("Order creation rolled back: {}", order.getId());
                }
            }
        }
    );
}
```

**Spring 4.2+ shortcut with `@TransactionalEventListener`:**

```java
@Service
public class OrderService {

    @Transactional
    public void createOrder(OrderRequest request) {
        Order order = orderRepo.save(new Order(request));
        applicationEventPublisher.publishEvent(new OrderCreatedEvent(order));
        // Event is NOT published immediately — held until TX commits
    }
}

@Component
public class OrderEventListener {

    @TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
    public void handleOrderCreated(OrderCreatedEvent event) {
        // Runs only if TX committed successfully
        emailService.sendConfirmation(event.getOrder());
    }

    @TransactionalEventListener(phase = TransactionPhase.AFTER_ROLLBACK)
    public void handleOrderRollback(OrderCreatedEvent event) {
        log.warn("Order TX rolled back: {}", event.getOrder().getId());
    }
}
```

**Synchronization lifecycle:**

```
TX Lifecycle:
[BEGIN] → beforeCommit() → beforeCompletion() → [COMMIT/ROLLBACK]
         → afterCommit() (only on commit)
         → afterCompletion(status)
```

---

### 10.3 What Happens at Commit Phase Internally?

```
@Transactional method returns successfully
          │
          ▼
AbstractPlatformTransactionManager.commit(status)
          │
          ├── 1. triggerBeforeCommit(status)
          │       → TransactionSynchronization.beforeCommit()
          │       → Last chance to add work to this TX
          │
          ├── 2. triggerBeforeCompletion(status)
          │       → TransactionSynchronization.beforeCompletion()
          │
          ├── 3. doCommit(status)
          │       │
          │       ├── JPA: entityManager.flush()
          │       │       (generate SQL, send to DB)
          │       │
          │       ├── JDBC: connection.commit()
          │       │       (tell DB to persist changes)
          │       │
          │       └── DB: Write WAL/Redo log → fsync → confirm
          │
          ├── 4. triggerAfterCommit()
          │       → TransactionSynchronization.afterCommit()
          │       → Safe to send events/notifications
          │
          ├── 5. triggerAfterCompletion(STATUS_COMMITTED)
          │       → Cleanup registered resources
          │
          └── 6. cleanupAfterCompletion(status)
                  → Remove TX from ThreadLocal
                  → Release DB connection back to pool
                  → Clear persistence context
```

**What the database does on commit:**

```
Application: connection.commit()
     │
     ▼
Database Engine:
  1. Flush dirty pages from buffer pool to disk (if needed)
  2. Write COMMIT record to WAL (Write-Ahead Log)
  3. fsync WAL to disk (durability guarantee)
  4. Release all locks held by this TX
  5. Make changes visible to other transactions
  6. Return success to application
```

---

### 10.4 How Rollback Is Triggered Internally

```
@Transactional method throws exception
          │
          ▼
TransactionAspectSupport.completeTransactionAfterThrowing()
          │
          ├── Check: Should this exception cause rollback?
          │     │
          │     ├── Is it RuntimeException or Error?
          │     │     → Yes: ROLLBACK (default rule)
          │     │
          │     ├── Is it in rollbackFor list?
          │     │     → Yes: ROLLBACK
          │     │
          │     ├── Is it in noRollbackFor list?
          │     │     → Yes: COMMIT (despite exception)
          │     │
          │     └── Checked exception not in any list?
          │           → COMMIT (default for checked exceptions!)
          │
          ▼ (if rollback decided)
AbstractPlatformTransactionManager.rollback(status)
          │
          ├── 1. triggerBeforeCompletion(status)
          │
          ├── 2. doRollback(status)
          │       │
          │       ├── JPA: Clear persistence context (discard changes)
          │       │
          │       ├── JDBC: connection.rollback()
          │       │
          │       └── DB: Apply undo log → restore original state
          │
          ├── 3. triggerAfterCompletion(STATUS_ROLLED_BACK)
          │
          └── 4. cleanupAfterCompletion(status)
                  → Same cleanup as commit
```

**Special case — rollback-only marking:**

```java
@Transactional
public void outerMethod() {
    try {
        innerService.riskyOperation(); // Throws, TX marked rollback-only
    } catch (Exception e) {
        // Caught the exception, but TX is ALREADY marked rollback-only!
        // outerMethod's commit attempt will fail with:
        // UnexpectedRollbackException
    }
}
```

```
What happens:
1. innerService TX (REQUIRED) joins outer TX
2. Exception in inner → Spring marks TX as rollback-only
3. You catch exception in outer
4. Outer method returns normally → Spring tries to COMMIT
5. TX is rollback-only → Spring rolls back + throws UnexpectedRollbackException!
```

This is one of the most confusing Spring transaction behaviors. **If an inner method fails and the TX is marked rollback-only, catching the exception does NOT save the transaction.**

---

### 10.5 How to Debug Transaction Issues in Production

#### 1. Enable Transaction Logging

```yaml
# application.yml
logging:
  level:
    org.springframework.transaction: DEBUG
    org.springframework.orm.jpa: DEBUG
    org.hibernate.SQL: DEBUG
    org.hibernate.type.descriptor.sql: TRACE  # Log bind parameters
```

**Sample output:**

```
DEBUG o.s.t.i.TransactionInterceptor - Getting transaction for [OrderService.createOrder]
DEBUG o.s.o.j.JpaTransactionManager - Creating new transaction with name [OrderService.createOrder]
DEBUG o.s.o.j.JpaTransactionManager - Opened new EntityManager [SessionImpl(123)] for JPA transaction
DEBUG o.s.o.j.JpaTransactionManager - Exposing JPA transaction as JDBC [HibernateJpaDialect$HibernateConnectionHandle]
...
DEBUG o.s.o.j.JpaTransactionManager - Initiating transaction commit
DEBUG o.s.o.j.JpaTransactionManager - Committing JPA transaction on EntityManager [SessionImpl(123)]
DEBUG o.s.o.j.JpaTransactionManager - Closing JPA EntityManager [SessionImpl(123)] after transaction
```

#### 2. Runtime Transaction Inspection

```java
@Aspect
@Component
@Slf4j
public class TransactionMonitorAspect {

    @Around("@annotation(transactional)")
    public Object monitor(ProceedingJoinPoint pjp, Transactional transactional) throws Throwable {
        String method = pjp.getSignature().toShortString();
        boolean isActive = TransactionSynchronizationManager.isActualTransactionActive();
        String txName = TransactionSynchronizationManager.getCurrentTransactionName();

        log.info("TX Monitor | Method: {} | Active: {} | Name: {} | Propagation: {}",
            method, isActive, txName, transactional.propagation());

        long start = System.currentTimeMillis();
        try {
            Object result = pjp.proceed();
            log.info("TX Monitor | Method: {} | Duration: {}ms | Status: COMMIT",
                method, System.currentTimeMillis() - start);
            return result;
        } catch (Exception e) {
            log.warn("TX Monitor | Method: {} | Duration: {}ms | Status: ROLLBACK | Error: {}",
                method, System.currentTimeMillis() - start, e.getMessage());
            throw e;
        }
    }
}
```

#### 3. Connection Pool Monitoring

```yaml
# HikariCP metrics
spring:
  datasource:
    hikari:
      register-mbeans: true  # JMX monitoring
      metrics-tracker-factory: com.zaxxer.hikari.metrics.prometheus.PrometheusMetricsTrackerFactory
```

**Key metrics to watch:**

```
┌──────────────────────────────────┬─────────────────────────────────┐
│ Metric                           │ Alert Threshold                 │
├──────────────────────────────────┼─────────────────────────────────┤
│ hikari_connections_active        │ > 80% of max pool size          │
│ hikari_connections_pending       │ > 0 sustained for > 5s          │
│ hikari_connections_timeout_total │ Any increment                   │
│ hikari_connections_usage_ms      │ p99 > 5 seconds                 │
└──────────────────────────────────┴─────────────────────────────────┘
```

#### 4. Database-Side Monitoring

```sql
-- MySQL: Long-running transactions
SELECT * FROM information_schema.INNODB_TRX
WHERE TIME_TO_SEC(TIMEDIFF(NOW(), trx_started)) > 10;

-- PostgreSQL: Long-running transactions
SELECT pid, now() - xact_start AS duration, query, state
FROM pg_stat_activity
WHERE xact_start IS NOT NULL
  AND state != 'idle'
ORDER BY xact_start;

-- PostgreSQL: Lock contention
SELECT relation::regclass, mode, granted, pid
FROM pg_locks
WHERE NOT granted;
```

#### 5. Production Debugging Checklist

```
Transaction not working?
─────────────────────────
□ Is the class a Spring bean? (@Service, @Component)
□ Is the method public?
□ Is it a self-invocation? (same-class call)
□ Is the exception checked? (won't rollback by default)
□ Is @EnableTransactionManagement present? (auto in Boot)
□ Is the correct TransactionManager being used? (multi-DB)
□ Is the method called from a proxy? (not a manual 'new' instance)

Transaction too slow?
─────────────────────────
□ Check for external API calls inside TX
□ Check for N+1 query issues (Hibernate)
□ Check for missing indexes (EXPLAIN ANALYZE)
□ Check for lock contention (DB monitoring)
□ Check connection pool saturation (HikariCP metrics)
□ Check for unnecessary flush operations

Data inconsistency?
─────────────────────────
□ Check propagation behavior (REQUIRED vs REQUIRES_NEW)
□ Check isolation level for concurrency issues
□ Check for missing rollbackFor on checked exceptions
□ Check for UnexpectedRollbackException (inner TX failure)
□ Check for race conditions (optimistic/pessimistic locking)
```

---

## Quick Reference Cheat Sheet

### @Transactional Defaults

```
propagation  = REQUIRED
isolation    = DEFAULT (DB-specific)
timeout      = -1 (no timeout)
readOnly     = false
rollbackFor  = RuntimeException, Error
noRollbackFor = (nothing)
```

### Propagation Quick Reference

```
REQUIRED      → Join or create       (DEFAULT — use for most cases)
REQUIRES_NEW  → Always create new    (audit logs, sequence generation)
NESTED        → Savepoint in current (optional steps that can fail)
SUPPORTS      → Join or run without  (read methods)
NOT_SUPPORTED → Suspend and run bare (long reads, reports)
NEVER         → Throw if TX exists   (external API calls)
MANDATORY     → Throw if no TX       (must be part of larger TX)
```

### Isolation Quick Reference

```
READ_UNCOMMITTED → See uncommitted data     (almost never use)
READ_COMMITTED   → See only committed data  (PostgreSQL default)
REPEATABLE_READ  → Consistent snapshot      (MySQL default)
SERIALIZABLE     → Full serial execution    (financial calculations)
```

### Golden Rules for Production

```
1. Always set timeout on @Transactional
2. Always use rollbackFor = Exception.class
3. Keep transactions SHORT (< 1 second for OLTP)
4. Never put external API calls inside transactions
5. Never use @Transactional on controllers
6. Disable Open Session in View (spring.jpa.open-in-view=false)
7. Use readOnly = true for all read operations
8. Watch for self-invocation (same-class method calls)
9. Use optimistic locking (@Version) over pessimistic by default
10. Monitor connection pool — most TX issues show there first
```

---
