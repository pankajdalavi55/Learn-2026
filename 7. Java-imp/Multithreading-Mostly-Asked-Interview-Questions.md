# Multithreading — Mostly Asked Interview Questions

> Interview-ready answers sized for ~4 minutes of speaking each.  
> Focus: clear definition → why it matters → how it works → when to use → common pitfalls.

---

## Table of Contents

1. [Process vs Thread](#q1-what-is-the-difference-between-a-process-and-a-thread)
2. [Why Multithreading?](#q2-why-do-we-need-multithreading)
3. [Thread Lifecycle](#q3-explain-the-java-thread-lifecycle)
4. [Creating Threads](#q4-how-do-you-create-a-thread-in-java)
5. [Runnable vs Callable](#q5-difference-between-runnable-and-callable)
6. [Race Condition](#q6-what-is-a-race-condition)
7. [synchronized Keyword](#q7-explain-the-synchronized-keyword)
8. [Object vs Class Lock](#q8-object-level-vs-class-level-locking)
9. [wait vs sleep](#q9-difference-between-wait-and-sleep)
10. [notify vs notifyAll](#q10-difference-between-notify-and-notifyall)
11. [Deadlock](#q11-what-is-deadlock-how-do-you-prevent-it)
12. [Livelock and Starvation](#q12-livelock-vs-starvation)
13. [volatile Keyword](#q13-explain-the-volatile-keyword)
14. [volatile vs synchronized](#q14-volatile-vs-synchronized)
15. [Atomic Variables & CAS](#q15-what-are-atomic-variables-and-cas)
16. [ExecutorService & Thread Pools](#q16-explain-executorservice-and-thread-pools)
17. [Future and CompletableFuture](#q17-future-vs-completablefuture)
18. [ReentrantLock vs synchronized](#q18-reentrantlock-vs-synchronized)
19. [CountDownLatch vs CyclicBarrier vs Semaphore](#q19-countdownlatch-vs-cyclicbarrier-vs-semaphore)
20. [ConcurrentHashMap](#q20-how-does-concurrenthashmap-work)
21. [Producer-Consumer](#q21-explain-the-producer-consumer-problem)
22. [ThreadLocal](#q22-what-is-threadlocal-and-when-do-you-use-it)
23. [Thread Safety Strategies](#q23-how-do-you-make-a-class-thread-safe)
24. [Daemon Threads](#q24-what-is-a-daemon-thread)
25. [Common Pitfalls & Best Practices](#q25-common-multithreading-pitfalls-in-production)

---

## Q1: What is the difference between a Process and a Thread?

**~4 min answer**

A **process** is an independent program in execution with its own memory space, file handles, and OS resources. A **thread** is the smallest unit of execution inside a process. Multiple threads in the same process share heap, code, and data, but each thread has its own stack, program counter, and registers.

Key differences I emphasize in interviews:

| Aspect | Process | Thread |
|--------|---------|--------|
| Memory | Isolated address space | Shares process heap |
| Creation cost | Heavy | Lightweight |
| Communication | IPC (pipes, sockets) — slower | Shared memory — faster |
| Context switch | Expensive | Cheaper |
| Failure | One crash usually does not kill others | Bad thread can crash the whole JVM process |

**Why this matters:** In Java, one JVM process typically runs many threads — HTTP request threads, GC threads, scheduler threads, async workers. Because they share heap, we get speed, but we also get race conditions, visibility issues, and deadlocks. That trade-off is the entire reason concurrency APIs exist.

**Real-world framing:** Chrome may use multiple processes for tabs (isolation/security). A Spring Boot order service usually uses one process with a Tomcat thread pool — many concurrent requests, shared caches and DB pools, so thread safety becomes a design concern.

**One-liner closer:** Process = isolation + cost. Thread = shared memory + speed + synchronization responsibility.

---

## Q2: Why do we need Multithreading?

**~4 min answer**

We use multithreading to improve **throughput**, **responsiveness**, and **resource utilization** — not “because concurrency is cool.”

Three practical reasons:

1. **CPU utilization** — On multi-core machines, CPU-bound work (image processing, aggregation, encryption) can run in parallel across cores.
2. **I/O concurrency** — Most backend apps are I/O-bound (DB, HTTP, Kafka). While one thread waits on network, others continue serving requests. A single-threaded server would block and waste time.
3. **Responsiveness** — UI or API gateways stay responsive while background work (email, report generation, cache warming) runs on worker threads.

**Where it shows up in Spring Boot:** Tomcat/Jetty request pool, `@Async`, `CompletableFuture`, Kafka consumer threads, scheduled jobs, parallel stream processing.

**Important caveat interviewers like:** More threads ≠ always faster. Too many threads cause context-switch overhead, lock contention, and memory pressure (each thread stack costs memory). For I/O-bound work, tune pool size around waiting; for CPU-bound work, roughly `cores` or `cores + 1` is a starting point.

**When not to multithread:** Tiny tasks where thread overhead dominates, simple sequential scripts, or cases where correctness is harder than the performance gain. Prefer async non-blocking or batching when appropriate.

**Closer:** Multithreading is a tool for overlapping wait and parallelizing compute — with correctness and pool sizing as first-class concerns.

---

## Q3: Explain the Java Thread Lifecycle.

**~4 min answer**

A Java thread moves through these states (`Thread.State`):

```
NEW → RUNNABLE → (RUNNING) → BLOCKED / WAITING / TIMED_WAITING → TERMINATED
```

- **NEW** — Thread object created, `start()` not called yet.
- **RUNNABLE** — Eligible to run. JVM/OS scheduler may be running it or ready to run it. Java does not expose a separate RUNNING state in the enum.
- **BLOCKED** — Waiting to acquire a monitor lock (`synchronized`).
- **WAITING** — Waiting indefinitely (`Object.wait()`, `LockSupport.park()`, `Thread.join()` without timeout).
- **TIMED_WAITING** — Waiting with timeout (`sleep`, `wait(timeout)`, `join(timeout)`, `Lock.tryLock(timeout)`).
- **TERMINATED** — `run()` finished or died with an exception.

**Common misconceptions I correct:**
- `sleep()` does **not** release the monitor lock; `wait()` does.
- Calling `run()` directly does **not** start a new thread — it runs on the current thread. Always use `start()` (or better, submit to an executor).
- A terminated thread cannot be restarted.

**Debugging tip:** In production dumps (`jstack` / thrreaddump), I look for many threads in BLOCKED (lock contention) or WAITING on the same lock (possible deadlock or slow critical section).

**Closer:** Lifecycle knowledge helps you reason about dumps, timeouts, and whether a thread is stuck on a lock, a condition, or just sleeping.

---

## Q4: How do you create a Thread in Java?

**~4 min answer**

Classic ways:

1. **Extend `Thread`** and override `run()`
2. **Implement `Runnable`** and pass to `Thread` or executor
3. **Implement `Callable`** and submit to `ExecutorService` (returns result / can throw checked exception)
4. **Prefer executors** in real code: `Executors` / `ThreadPoolExecutor` / `@Async`

```java
// Preferred in production: submit work to a pool
ExecutorService pool = Executors.newFixedThreadPool(8);
pool.submit(() -> processOrder(orderId));
```

**Why Runnable over extending Thread?** Composition over inheritance. Your class may already extend something else. Separating task from execution mechanism lets the same task run on different executors.

**Modern practice:** Almost never create raw `new Thread()` in business code. Use a bounded thread pool so you control concurrency, queueing, rejection policy, and naming (important for dumps). Unbounded thread creation under load can OOM the JVM.

**Naming & MDC:** In production I set a `ThreadFactory` with meaningful names (`order-worker-1`) and propagate request IDs/MDC so logs remain traceable across async boundaries.

**Closer:** Know the classic APIs for interviews; in production, submit `Runnable`/`Callable`/`CompletableFuture` tasks to a sized pool.

---

## Q5: Difference between Runnable and Callable?

**~4 min answer**

| Point | `Runnable` | `Callable<V>` |
|-------|------------|---------------|
| Method | `void run()` | `V call()` |
| Return value | None | Returns result |
| Checked exceptions | Cannot throw | Can throw |
| Introduced | Java 1.0 | Java 5 (`java.util.concurrent`) |
| Typical use | Fire-and-forget tasks | Tasks that compute a value |

```java
Callable<Integer> task = () -> {
    return inventoryService.reserve(itemId);
};
Future<Integer> future = executor.submit(task);
Integer reserved = future.get(2, TimeUnit.SECONDS);
```

**Interview depth:** `ExecutorService.submit(Callable)` returns `Future`. `invokeAll` / `invokeAny` work with callables for fan-out patterns. With `CompletableFuture`, we often use `supplyAsync` (like callable) vs `runAsync` (like runnable).

**Production angle:** Prefer returning values via `Callable`/`CompletableFuture` over mutating shared variables from worker threads — clearer data flow, fewer races.

**Closer:** Runnable = side-effect task. Callable = compute-and-return task with proper exception surface.

---

## Q6: What is a Race Condition?

**~4 min answer**

A **race condition** happens when the correctness of a program depends on the unpredictable timing/interleaving of threads accessing shared mutable state.

Classic example: non-atomic `count++`.

```java
// NOT thread-safe
count++; // read → increment → write (3 steps)
```

Two threads can both read `10`, both write `11`, and one increment is lost.

**Three ingredients:**
1. Shared data
2. Mutation
3. No (or insufficient) synchronization / atomicity / immutability

**Related terms:**
- **Critical section** — code that must execute atomically w.r.t. shared state
- **Data race** — concurrent conflicting accesses without happens-before (JMM term)
- **Check-then-act bug** — `if (!map.containsKey(k)) map.put(k, v)` without atomicity

**How I fix them:**
- Eliminate sharing (thread confinement)
- Make state immutable
- Use `synchronized` / locks
- Use concurrent collections / atomics
- Use atomic APIs: `ConcurrentHashMap.computeIfAbsent`, `AtomicInteger.incrementAndGet`

**Detection:** Unit tests with many threads help but are flaky; better to design for safety and use stress tests, race detectors, and code review for check-then-act patterns.

**Closer:** Race conditions are about timing-dependent correctness on shared mutable state — fix by removing the race window, not by hoping the scheduler is kind.

---

## Q7: Explain the `synchronized` keyword.

**~4 min answer**

`synchronized` provides **mutual exclusion** and **visibility** using intrinsic object monitors.

Forms:
1. Synchronized instance method → locks `this`
2. Synchronized static method → locks `Class` object
3. Synchronized block → locks a specific object (preferred for smaller critical sections)

```java
public void transfer(Account to, int amount) {
    synchronized (this) {
        // critical section only
        this.balance -= amount;
        to.balance += amount; // may need finer locking design
    }
}
```

**What it guarantees:**
- Only one thread holds a given monitor at a time
- Unlocking flushes writes; acquiring sees the latest writes (happens-before)

**Design tips:**
- Keep critical sections short
- Lock on private `final` objects, not on public APIs / `this` if callers might lock on your instance
- Avoid nested locks with inconsistent ordering (deadlock risk)

**Limitations vs explicit locks:** No try-lock, no interruptible lock wait, no fair locking option, single condition queue (`wait/notify`). For advanced needs, use `ReentrantLock`.

**Closer:** `synchronized` is the simplest correctness tool for mutual exclusion + visibility — still valid and often enough; don’t over-engineer with locks unless you need the extra features.

---

## Q8: Object-level vs Class-level locking?

**~4 min answer**

- **Object-level lock:** `synchronized` instance method / `synchronized(this)` / `synchronized(instanceLock)` — protects **per-instance** state.
- **Class-level lock:** `synchronized static` method / `synchronized(MyClass.class)` — protects **static / shared-across-instances** state.

```java
public synchronized void instanceMethod() { }      // lock = this
public static synchronized void staticMethod() { } // lock = MyClass.class
```

These are **different monitors**. An instance method and a static method can run concurrently because they don’t compete for the same lock.

**Interview trap:** Synchronizing instance methods does **not** protect static fields. Synchronizing static methods does **not** protect instance fields of multiple objects unless all access goes through that class lock (usually the wrong design).

**Practical rule:**
- Instance state → instance lock (or better: confine / concurrent structures)
- Static mutable cache/counters → class lock or `Atomic*` / concurrent map

**Closer:** Always ask “what exact memory am I protecting, and which monitor guards it?” Mismatched lock scope is a common bug.

---

## Q9: Difference between `wait()` and `sleep()`?

**~4 min answer**

| Point | `wait()` | `sleep()` |
|-------|----------|-----------|
| Belongs to | `Object` | `Thread` |
| Releases monitor lock | **Yes** | **No** |
| Needs synchronized | Yes (must own monitor) | No |
| Woken by | `notify` / `notifyAll` / interrupt / spurious wakeups | Time elapse / interrupt |
| Typical use | Condition-based coordination | Pause / pacing / retry backoff |

```java
synchronized (lock) {
    while (!condition) {
        lock.wait(); // releases lock, waits for signal
    }
    // proceed when condition is true
}
```

**Must use `while`, not `if`, with wait:** spurious wakeups and multiple waiters mean you re-check the condition.

**Production note:** Prefer higher-level constructs (`BlockingQueue`, `CountDownLatch`, `Condition`, `CompletableFuture`) over hand-rolled `wait/notify` — fewer bugs.

**Closer:** `sleep` = “pause me but keep the lock if I hold one.” `wait` = “release the lock and wait until someone signals that state may have changed.”

---

## Q10: Difference between `notify()` and `notifyAll()`?

**~4 min answer**

Both wake threads waiting on the same monitor:

- **`notify()`** — wakes **one** arbitrary waiter
- **`notifyAll()`** — wakes **all** waiters; they compete to re-acquire the lock and re-check conditions

**When `notify()` is risky:** Multiple waiters waiting for **different** conditions on the same lock (e.g., producers waiting for “not full” and consumers waiting for “not empty”). `notify()` might wake the wrong kind of thread → missed signal / stuck system. `notifyAll()` is safer.

**When `notify()` can be fine:** All waiters are waiting for the **same** condition and any one can proceed.

**Modern preference:** `BlockingQueue` for producer-consumer; `ReentrantLock` + multiple `Condition` objects (`notFull`, `notEmpty`) for precise signaling without waking everyone.

**Closer:** Default to `notifyAll` (or better abstractions) unless you have a proven single-condition, single-purpose wait set.

---

## Q11: What is Deadlock? How do you prevent it?

**~4 min answer**

**Deadlock** is a permanent hold-and-wait cycle: Thread A holds lock L1 and waits for L2; Thread B holds L2 and waits for L1. Nobody progresses.

**Four Coffman conditions** (all required):
1. Mutual exclusion
2. Hold and wait
3. No preemption of locks
4. Circular wait

**Prevention strategies (break a condition):**
1. **Lock ordering** — always acquire locks in a global order (e.g., by account ID ascending)
2. **Try-lock with timeout** — `ReentrantLock.tryLock(timeout)` and back off
3. **Avoid nested locks** — shrink critical sections; use concurrent structures
4. **Single lock / lock striping carefully** — reduce lock graph complexity
5. **Lock-free / atomic designs** where appropriate

```java
// Consistent ordering example
Object first = System.identityHashCode(a) < System.identityHashCode(b) ? a : b;
Object second = (first == a) ? b : a;
synchronized (first) {
    synchronized (second) {
        // transfer
    }
}
```

**Detection:** Thread dump shows threads blocked on each other’s locks; some JVMs / tools can detect cycles. In prod, I’d capture `jstack`, identify the cycle, fix ordering, add timeouts/metrics for lock wait time.

**Closer:** Deadlocks are design bugs in lock acquisition graphs — prevent with ordering and timeouts; diagnose with dumps.

---

## Q12: Livelock vs Starvation?

**~4 min answer**

- **Deadlock:** Threads are stuck waiting; no progress.
- **Livelock:** Threads keep changing state in response to each other but still make no useful progress (like two people politely stepping aside forever).
- **Starvation:** A thread never gets CPU/lock access because others keep getting preferred (e.g., unfair lock + continuous high-priority work).

**Livelock example:** Two transactions detect conflict, both rollback and immediately retry with the same timing, colliding forever. Fix with randomized backoff / jitter.

**Starvation example:** Readers keep entering a read lock and writers never run — or unfair locks always let newcomers cut in. Fix with fair locks (costly), writer preference policies, or bounded admission.

**Closer:** Deadlock = frozen. Livelock = busy but useless. Starvation = some thread forever ignored. All are liveness failures, not just safety bugs.

---

## Q13: Explain the `volatile` keyword.

**~4 min answer**

`volatile` guarantees **visibility** and **ordering for that variable’s reads/writes**, but **not atomicity of compound actions**.

If Thread A writes `volatile flag = true`, Thread B reading `flag` is guaranteed to see the write (happens-before). Without volatile/synchronization, B may see a stale cached value indefinitely.

```java
private volatile boolean shutdown = false;

public void stop() { shutdown = true; }

public void run() {
    while (!shutdown) {
        doWork();
    }
}
```

**What volatile does NOT do:** Make `count++` thread-safe. That is read-modify-write. Use `AtomicInteger` or synchronization.

**volatile + double-checked locking:** Classic singleton pattern uses volatile on the instance reference so publication is safe under JMM.

**Closer:** volatile = “always read/write main memory visibility for this field.” It is not a general replacement for locks.

---

## Q14: `volatile` vs `synchronized`?

**~4 min answer**

| Need | Prefer |
|------|--------|
| Visibility of a single flag/status | `volatile` |
| Atomic compound update (`check-then-act`, increment) | `synchronized` / lock / atomic |
| Mutual exclusion of a critical section | `synchronized` / lock |
| Multiple related fields updated together | `synchronized` / lock (or immutable snapshot publication) |

`synchronized` gives exclusion + visibility. `volatile` gives visibility/ordering for one variable without exclusion.

**Rule of thumb I use:** If the question is “can two threads be in this code at once and still be correct?” and the answer is no → need a lock/atomic, not just volatile.

**Closer:** volatile for simple publication/flags; synchronized/atomics for atomic state transitions.

---

## Q15: What are Atomic variables and CAS?

**~4 min answer**

`java.util.concurrent.atomic` (`AtomicInteger`, `AtomicLong`, `AtomicReference`, etc.) provides lock-free thread-safe updates using **CAS (Compare-And-Swap)** at the hardware/JVM level.

CAS algorithm:
1. Read current value
2. Compute new value
3. Atomically: if value is still expected, set new; else fail and retry

```java
AtomicInteger count = new AtomicInteger(0);
count.incrementAndGet(); // lock-free atomic ++
```

**Why useful:** For hot counters/flags, atomics often beat coarse `synchronized` because they avoid lock contention and context switches — under very high contention, though, CAS retries can burn CPU (then consider LongAdder / striping).

**ABA problem:** Value changes A→B→A; CAS succeeds even though intermediate change happened. Sometimes use `AtomicStampedReference` / `AtomicMarkableReference`.

**Closer:** Atomics give atomic read-modify-write without explicit locks via CAS; great for counters and simple state machines.

---

## Q16: Explain ExecutorService and Thread Pools.

**~4 min answer**

Creating a thread per task is expensive and dangerous under load. An **`ExecutorService`** manages a pool of worker threads and a queue of tasks.

Common pools:
- **Fixed** — constant thread count; bounded concurrency
- **Cached** — creates threads on demand, reuses idle ones (can explode under load — use carefully)
- **Single-thread** — sequential async execution
- **Scheduled** — delayed/periodic tasks
- **Custom `ThreadPoolExecutor`** — production default choice

Key tuning knobs for `ThreadPoolExecutor`:
- `corePoolSize`, `maximumPoolSize`
- `workQueue` (bounded `ArrayBlockingQueue` preferred over unbounded `LinkedBlockingQueue` for overload protection)
- `RejectedExecutionHandler` (abort, caller-runs, discard, custom)
- `ThreadFactory` (naming, daemon flag, UncaughtExceptionHandler)
- keep-alive time

```java
ThreadPoolExecutor executor = new ThreadPoolExecutor(
    8, 16,
    60, TimeUnit.SECONDS,
    new ArrayBlockingQueue<>(1000),
    new ThreadFactoryBuilder().setNameFormat("pay-%d").build(),
    new ThreadPoolExecutor.CallerRunsPolicy()
);
```

**Shutdown:** `shutdown()` (orderly) vs `shutdownNow()` (interrupt). Always plan graceful shutdown in Spring via lifecycle hooks.

**Closer:** Thread pools convert unbounded thread creation into controlled concurrency + queueing + rejection policy — essential for stable services.

---

## Q17: Future vs CompletableFuture?

**~4 min answer**

**`Future`** represents a result of async computation. You `get()` (blocking), cancel, or check `isDone`. Limited composition — hard to chain without blocking.

**`CompletableFuture`** (Java 8+) is a reactive-style Future:
- Non-blocking callbacks: `thenApply`, `thenAccept`, `thenCompose`, `thenCombine`
- Combine many: `allOf`, `anyOf`
- Async variants with custom executors
- Exception handling: `exceptionally`, `handle`, `whenComplete`

```java
CompletableFuture<User> user = supplyAsync(() -> userClient.get(id), pool);
CompletableFuture<Credit> credit = supplyAsync(() -> creditClient.get(id), pool);

CompletableFuture<Decision> decision =
    user.thenCombine(credit, this::approve)
        .orTimeout(2, TimeUnit.SECONDS)
        .exceptionally(ex -> Decision.REJECT);
```

**Production warnings:**
- Default `ForkJoinPool.commonPool()` may be wrong for blocking I/O — pass your own executor
- Always handle exceptions; otherwise failures can be swallowed until `join/get`
- Avoid blocking inside `thenApply` stages; use `thenCompose` + async I/O minded design

**Closer:** Future = placeholder result. CompletableFuture = composable async pipelines — powerful if you control executors and errors.

---

## Q18: ReentrantLock vs synchronized?

**~4 min answer**

Both provide reentrant mutual exclusion. Prefer `synchronized` for simple cases; use `ReentrantLock` when you need features.

| Feature | synchronized | ReentrantLock |
|---------|--------------|---------------|
| Lock release | Automatic (end of block) | Must `unlock()` in `finally` |
| Try lock | No | `tryLock()` |
| Interruptible wait | No (block on monitor) | `lockInterruptibly()` |
| Fairness option | No | Yes (`new ReentrantLock(true)`) |
| Multiple conditions | One wait set | Multiple `Condition`s |
| Structured try-with | Language-level | Manual / careful |

```java
lock.lock();
try {
    // critical section
} finally {
    lock.unlock();
}
```

**Reentrancy:** Same thread can re-acquire the same lock (needed for nested calls). Both support this.

**Closer:** Start with synchronized; graduate to ReentrantLock for try-lock, fairness, interruptibility, or multiple conditions.

---

## Q19: CountDownLatch vs CyclicBarrier vs Semaphore?

**~4 min answer**

All are synchronizers, different coordination patterns:

**CountDownLatch**
- One-shot gate
- Threads wait until count reaches 0
- Typical: “start workers after init” or “wait until N tasks finish”
- Cannot reset

**CyclicBarrier**
- Parties wait for each other at a barrier point
- Reusable (cyclic)
- Optional barrier action runs when last party arrives
- Typical: multi-phase parallel algorithms

**Semaphore**
- Permits controlling access to a limited resource
- `acquire()` / `release()`
- Typical: DB connection throttling, rate limiting concurrent outbound calls

| Use case | Tool |
|----------|------|
| Main thread waits for N workers once | CountDownLatch |
| N threads repeatedly sync at phases | CyclicBarrier |
| Limit concurrent access to K resources | Semaphore |

```java
Semaphore dbPermits = new Semaphore(20);
dbPermits.acquire();
try { callDatabase(); }
finally { dbPermits.release(); }
```

**Closer:** Latch = wait for events/countdown. Barrier = wait for peers. Semaphore = throttle concurrency.

---

## Q20: How does ConcurrentHashMap work?

**~4 min answer**

`ConcurrentHashMap` is a thread-safe map designed for high concurrent reads/writes without locking the entire map (unlike `Collections.synchronizedMap` / `Hashtable`).

**High-level behavior (Java 8+):**
- Uses CAS and synchronized locking at **bin/node** level for updates
- Reads are usually non-blocking
- Does **not** allow `null` keys/values (unlike `HashMap`)
- Iterators are weakly consistent (may reflect later updates; won’t throw `ConcurrentModificationException`)

**Important API choices:**
- Prefer `putIfAbsent`, `computeIfAbsent`, `compute`, `merge` for atomic check-then-act
- Size is an estimate under concurrent modification (not a strongly consistent snapshot)

```java
cache.computeIfAbsent(userId, id -> loadFromDb(id));
```

**Common interview misconception:** “ConcurrentHashMap makes any compound action safe.” False.  
`if (!map.containsKey(k)) map.put(k,v)` is still racy. Use atomic methods.

**When I use it:** shared caches, metrics maps, request-scoped registries across threads. For mostly-read maps with rare writes, it shines.

**Closer:** ConcurrentHashMap gives concurrent map operations safely and scalably — but you must still use its atomic APIs for compound logic.

---

## Q21: Explain the Producer-Consumer problem.

**~4 min answer**

Producers generate data; consumers process it. They share a bounded buffer. Challenges: don’t overfill, don’t consume from empty, keep thread-safe.

**Classic solution today:** `BlockingQueue` (`ArrayBlockingQueue`, `LinkedBlockingQueue`, `LinkedTransferQueue`).

```java
BlockingQueue<Order> queue = new ArrayBlockingQueue<>(1000);

// producer
queue.put(order); // waits if full

// consumer
Order order = queue.take(); // waits if empty
```

**Why this is interview gold:** It demonstrates wait/notify concepts without reinventing them. Bounded queues provide **backpressure** — critical in production so producers don’t OOM the JVM.

**Design knobs:**
- Queue capacity (backpressure vs latency)
- Number of producers/consumers
- Rejection / timeout policies (`offer` with timeout vs blocking `put`)
- Poison-pill or shutdown signals for graceful stop

**Spring/Kafka angle:** Same pattern appears in async processing pipelines and consumer groups — buffer + workers + backpressure.

**Closer:** Producer-consumer is coordination around a bounded buffer; prefer `BlockingQueue` and treat capacity as a stability control, not an afterthought.

---

## Q22: What is ThreadLocal and when do you use it?

**~4 min answer**

`ThreadLocal` gives each thread its own independent variable copy — **thread confinement** without explicit synchronization.

Common uses:
- Request context / correlation ID / security principal in web apps
- DateFormat / SimpleDateFormat legacy thread-safety (better: `DateTimeFormatter` which is immutable)
- Transaction/session context binding

```java
private static final ThreadLocal<String> REQUEST_ID = new ThreadLocal<>();

REQUEST_ID.set(id);
try {
    // downstream code can read REQUEST_ID.get()
} finally {
    REQUEST_ID.remove(); // CRITICAL in thread pools
}
```

**Biggest production bug:** Forgetting `remove()` in thread pools. Threads are reused; leftover values leak memory and leak data across requests (security/privacy incident waiting to happen).

**Modern alternatives:** Prefer passing context explicitly, Micrometer/Observation context, or framework request scopes. Use ThreadLocal sparingly.

**Closer:** ThreadLocal = per-thread state. Powerful for context propagation; dangerous in pools without reliable cleanup.

---

## Q23: How do you make a Class Thread-safe?

**~4 min answer**

I answer with a strategy ladder:

1. **Immutability** — final fields, no setters, defensive copies. Best default if possible.
2. **Thread confinement** — don’t share the instance across threads.
3. **Synchronization / locks** — guard all accesses to mutable shared state with the same lock policy.
4. **Concurrent collections / atomics** — replace coarse locks where contention matters.
5. **Actor-style / single-writer** — queue mutations to one thread; others only read snapshots.
6. **Safe publication** — ensure references are visible after construction (`volatile`, final, synchronized, concurrent structures).

**Checklist interviewers like:**
- Are all mutable fields protected?
- Is lock scope consistent?
- Are check-then-act sequences atomic?
- Are invariants preserved under concurrency?
- Are we documenting the thread-safety contract?

**Example:** A counter class — prefer `LongAdder`/`AtomicLong` over sync methods if it’s just a hot counter. A complex aggregate with multi-field invariants — use one lock or immutable updates.

**Closer:** Thread safety is a design property of state + access protocol, not a keyword you sprinkle randomly.

---

## Q24: What is a Daemon Thread?

**~4 min answer**

A **daemon thread** is a background thread. The JVM exits when only daemon threads remain; non-daemon (user) threads keep the JVM alive.

Examples: garbage collector, some housekeeping timers. You can mark a thread daemon via `thread.setDaemon(true)` **before** `start()`.

**Interview caution:** Don’t run critical business work as daemon — JVM may kill it mid-write on shutdown. For background workers that must finish cleanly, use non-daemon pool threads + orderly shutdown hooks.

**In executors:** Default thread factory creates non-daemon threads. That’s usually what you want for application workers.

**Closer:** Daemon = “don’t keep process alive for me.” Use for optional background helpers, not for must-complete work.

---

## Q25: Common Multithreading pitfalls in Production?

**~4 min answer**

These are the ones I’ve seen cause real outages:

1. **Unbounded thread creation / cached pools under spike** → native memory / OOM
2. **Unbounded queues** → delayed failure; heap fills; latency explodes before rejection
3. **Shared mutable state without sync** → rare Heisenbugs in prod only
4. **Deadlocks from inconsistent lock order** → stuck requests, cascading timeouts
5. **ThreadLocal leaks in Tomcat pools** → memory growth + cross-request data bleed
6. **Blocking the common ForkJoinPool** with I/O in `parallelStream` / default `CompletableFuture` → app-wide stalls
7. **Ignoring interrupts** → shutdown hangs
8. **No timeouts on `future.get()`** → threads stuck forever on downstream dependency
9. **Synchronizing on external/untrusted objects** → accidental deadlock with callers
10. **Logging/MDC not propagated to async threads** → untraceable incidents

**What good looks like:**
- Bounded pools + bounded queues + explicit rejection policy
- Timeouts everywhere for remote calls and waits
- Meaningful thread names
- Metrics: queue depth, active threads, reject count, lock wait
- Prefer immutability and concurrent utilities over custom wait/notify

**Closer:** Most concurrency incidents are operational — pool sizing, timeouts, leaks, and blocking on shared pools — not exotic algorithm bugs.

---

## Quick Revision Cheat Sheet

| Topic | Remember |
|-------|----------|
| Process vs Thread | Isolation vs shared heap |
| Race | Shared mutable + no atomicity |
| synchronized | Exclusion + visibility |
| wait vs sleep | Releases lock vs does not |
| volatile | Visibility, not compound atomicity |
| Atomic/CAS | Lock-free RMW |
| Executor | Bound concurrency + queue + rejection |
| CHM | Concurrent map ops; use atomic APIs |
| Latch/Barrier/Semaphore | Countdown / peer sync / permits |
| ThreadLocal | Always `remove()` in pools |
| Deadlock | Fix ordering + tryLock timeouts |

---

## How to Use This File

1. Speak each answer out loud in ~4 minutes.
2. For each question, add one personal production example.
3. Be ready for follow-ups: “How would you debug this in prod?” and “What metrics would you add?”
4. Pair with `01-Multithreading.md` for deeper code walkthroughs.

---

*Related guides:* `01-Multithreading.md` · `02-springboot-multithreading.md` · `Java-Collections-Internals-Performance-Interview-Guide.md`
