# Senior Golang Interview Question Handbook - Part 1

> **Target Audience:** Experienced Golang Developers (5–12 years)  
> **Purpose:** Prepare for senior backend/product company interviews  
> **Difficulty Legend:** 🟢 Basic | 🟡 Intermediate | 🔴 Advanced | ⚫ Expert

---

## Table of Contents - Part 1

1. [Go Fundamentals (Senior Level Refresh)](#1-go-fundamentals-senior-level-refresh)
2. [Memory Management & Internals](#2-memory-management--internals)
3. [Goroutines & Concurrency](#3-goroutines--concurrency-critical)
4. [Go Data Structures & Performance](#4-go-data-structures--performance)

---

## 1. Go Fundamentals (Senior Level Refresh)

### 1.1 Go Philosophy & Design Goals

#### 🟢 Q1: What are Go's core design principles?

**Answer:**
- **Simplicity over cleverness** — minimal syntax, no implicit magic
- **Composition over inheritance** — embedding, interfaces
- **Explicitness** — no hidden control flow, explicit error handling
- **Fast compilation** — dependency management designed for speed
- **Built-in concurrency** — goroutines and channels as first-class citizens
- **Single binary deployment** — static linking by default

---

#### 🟡 Q2: Why does Go have no inheritance?

**Answer:**
Go deliberately avoids inheritance to prevent:
- **Fragile base class problem** — changes in parent break children
- **Diamond problem** — ambiguity in multiple inheritance
- **Tight coupling** — inheritance creates rigid hierarchies

Go uses **composition** and **interface satisfaction** instead:

```go
// Composition via embedding
type Logger struct{}
func (l Logger) Log(msg string) { fmt.Println(msg) }

type Service struct {
    Logger // Embedded - Service "has-a" Logger, not "is-a"
}

// Interface satisfaction is implicit
type Writer interface {
    Write([]byte) (int, error)
}
// Any type with Write method satisfies Writer - no "implements" keyword
```

---

#### 🟡 Q3: Why are errors values in Go?

**Answer:**
Errors as values provide:
- **Explicit error handling** — forces developers to handle errors at call site
- **Composability** — errors can be wrapped, compared, type-asserted
- **No hidden control flow** — unlike exceptions, errors don't jump stack frames
- **Performance** — no stack unwinding overhead

```go
// Errors are just interface implementations
type error interface {
    Error() string
}

// Custom errors with context
type QueryError struct {
    Query string
    Err   error
}

func (e *QueryError) Error() string {
    return fmt.Sprintf("query %q failed: %v", e.Query, e.Err)
}

func (e *QueryError) Unwrap() error { return e.Err }
```

---

#### 🔴 Q4: Go vs Java/Python — When would you NOT choose Go?

**Answer:**

| Scenario | Preferred Over Go | Reason |
|----------|-------------------|--------|
| Heavy OOP with deep hierarchies | Java | Go lacks inheritance, generics were limited until 1.18 |
| Rapid prototyping/scripting | Python | Go requires explicit types, compilation step |
| GUI applications | Java/Python | Go's GUI ecosystem is immature |
| Heavy numerical computing | Python (NumPy) | Go lacks mature scientific libraries |
| Enterprise ecosystems | Java | Spring ecosystem, established patterns |

**Choose Go for:**
- High-concurrency network services
- CLI tools
- Microservices
- Infrastructure tooling (Docker, K8s written in Go)

---

### 1.2 Compilation Model

#### 🟡 Q5: Explain Go's compilation model and why it's fast.

**Answer:**

**Key factors for fast compilation:**

1. **Dependency management** — imports are explicit; compiler only processes what's needed
2. **No header files** — package object files contain all export info
3. **No circular dependencies** — enforced by compiler, enables parallel compilation
4. **Unused import = error** — no dead code in compilation
5. **Single pass compilation** — declarations must precede usage

```bash
# Compilation produces single binary
go build -o myapp main.go

# Cross-compilation in one line
GOOS=linux GOARCH=amd64 go build -o myapp-linux main.go
```

---

#### 🟡 Q6: Explain cross-compilation in Go.

**Answer:**

Go supports cross-compilation natively via `GOOS` and `GOARCH`:

```bash
# Build for Linux on Windows
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o app-linux

# Build for ARM (Raspberry Pi)
$env:GOOS="linux"; $env:GOARCH="arm"; $env:GOARM="7"; go build

# Build for macOS
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build
```

**Caveats:**
- CGO disabled by default for cross-compilation (`CGO_ENABLED=0`)
- Pure Go code cross-compiles seamlessly
- C dependencies require target platform toolchain

---

### 1.3 Go Runtime Overview

#### 🔴 Q7: What does the Go runtime provide?

**Answer:**

The Go runtime is **not a VM** — it's compiled into every Go binary:

| Component | Responsibility |
|-----------|----------------|
| **Scheduler** | M:N scheduling of goroutines onto OS threads |
| **Garbage Collector** | Concurrent, tri-color mark-sweep GC |
| **Memory Allocator** | Per-P mcache, mcentral, mheap hierarchy |
| **Stack Management** | Growable stacks, stack copying |
| **Channel Operations** | Send/receive synchronization |
| **Reflection** | Runtime type information |

**Runtime overhead:** ~2MB base memory, but provides automatic memory management and scheduling.

---

#### 🔴 Q8: What is the difference between `runtime.GOMAXPROCS` and actual parallelism?

**Answer:**

```go
runtime.GOMAXPROCS(4) // Sets max number of OS threads executing Go code
```

- **GOMAXPROCS** = number of **P** (logical processors) in Go scheduler
- Default = number of CPU cores
- More goroutines ≠ more parallelism
- Only `GOMAXPROCS` goroutines run **truly in parallel**

```go
// Example: 1000 goroutines, GOMAXPROCS=4
// Only 4 execute simultaneously; rest are scheduled cooperatively
```

**When to tune:**
- CPU-bound work: GOMAXPROCS = CPU cores (default)
- I/O-bound work: can exceed cores (goroutines block on I/O)

---

## 2. Memory Management & Internals

### 2.1 Stack vs Heap Allocation

#### 🟡 Q9: How does Go decide stack vs heap allocation?

**Answer:**

**Escape Analysis** determines allocation:

| Allocated On | Condition |
|--------------|-----------|
| **Stack** | Variable doesn't outlive function, compiler proves it won't escape |
| **Heap** | Variable escapes to heap (pointer returned, stored in interface, etc.) |

```go
// Stack allocation - s doesn't escape
func stackAlloc() int {
    s := 42
    return s // Value copied, s stays on stack
}

// Heap allocation - s escapes
func heapAlloc() *int {
    s := 42
    return &s // Pointer returned, s escapes to heap
}
```

**Check escape analysis:**
```bash
go build -gcflags="-m" main.go
# Output: ./main.go:10:2: moved to heap: s
```

---

#### 🔴 Q10: What is escape analysis and how do you use it for optimization?

**Answer:**

Escape analysis is a compile-time analysis determining if variables can be stack-allocated.

**Common escape scenarios:**
```go
// 1. Returning pointer
func escape1() *int { x := 1; return &x } // x escapes

// 2. Storing in interface
func escape2() interface{} { x := 1; return x } // x escapes

// 3. Closure capturing variable
func escape3() func() int {
    x := 1
    return func() int { return x } // x escapes
}

// 4. Slice/map with unknown size
func escape4(n int) []int { return make([]int, n) } // escapes

// 5. Sending pointer to channel
func escape5(ch chan *int) { x := 1; ch <- &x } // x escapes
```

**Optimization techniques:**
```go
// BAD: Unnecessary heap allocation
func processBad(data []byte) *Result {
    r := &Result{} // escapes
    r.Parse(data)
    return r
}

// GOOD: Let caller decide allocation
func processGood(data []byte, r *Result) {
    r.Parse(data) // r provided by caller, may be stack-allocated
}

// GOOD: Return value instead of pointer
func processBetter(data []byte) Result {
    var r Result // stack-allocated
    r.Parse(data)
    return r // copied, but no heap allocation
}
```

---

### 2.2 Garbage Collector

#### 🔴 Q11: Explain Go's garbage collector architecture.

**Answer:**

**Go GC: Concurrent, Tri-color, Mark-Sweep**

**Phases:**
1. **Mark Setup** (STW) — Enable write barrier, prepare for marking (~10-30μs)
2. **Marking** (Concurrent) — Traverse object graph, mark reachable objects
3. **Mark Termination** (STW) — Drain work queues, disable write barrier (~60-90μs)
4. **Sweeping** (Concurrent) — Reclaim unmarked memory

**Tri-color marking:**
- **White** — Not yet visited (potentially garbage)
- **Grey** — Visited, but children not yet scanned
- **Black** — Visited, all children scanned (definitely reachable)

**Write barrier:** Ensures no black object points to white object during concurrent marking.

```go
// GC behavior configuration
debug.SetGCPercent(100)  // Trigger GC when heap doubles (default)
debug.SetMemoryLimit(1 << 30)  // Go 1.19+: soft memory limit
```

---

#### 🔴 Q12: Why is Go's GC non-generational?

**Answer:**

**Generational hypothesis:** Most objects die young, so separate young/old generations with different collection frequencies.

**Why Go skipped generational GC:**

1. **Write barrier cost** — Generational GC needs write barriers even during mutator
2. **Go's allocation patterns** — Many small, short-lived allocations (goroutine stacks)
3. **Concurrent collector** — Already achieves low latency without generations
4. **Complexity** — Generational adds complexity for marginal gains in Go's use cases

**Go's approach instead:**
- Very fast allocator (per-P caches)
- Concurrent marking reduces STW
- GOGC tuning for latency vs throughput trade-off

---

#### ⚫ Q13: How to reduce GC pressure in high-throughput services?

**Answer:**

**1. Object pooling:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func handler(w http.ResponseWriter, r *http.Request) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    // Use buf...
}
```

**2. Reduce allocations:**
```go
// BAD: Allocates on every call
func concat(a, b string) string {
    return a + b // Allocates new string
}

// GOOD: Pre-allocate
func concatMany(strs []string) string {
    var b strings.Builder
    b.Grow(totalLen(strs)) // Pre-allocate
    for _, s := range strs {
        b.WriteString(s)
    }
    return b.String()
}
```

**3. Avoid pointer-heavy structures:**
```go
// BAD: Slice of pointers - GC scans all pointers
type Bad struct {
    items []*Item
}

// GOOD: Slice of values - GC only scans one pointer
type Good struct {
    items []Item
}
```

**4. Tune GOGC:**
```go
// Reduce GC frequency (trade memory for CPU)
debug.SetGCPercent(200) // GC when heap is 3x live data

// Go 1.19+: Set memory limit
debug.SetMemoryLimit(4 << 30) // 4GB limit
```

**5. Use arena (Go 1.20+ experimental):**
```go
import "arena"

func processRequest() {
    a := arena.NewArena()
    defer a.Free() // Free all at once
    
    data := arena.MakeSlice[byte](a, 1024, 1024)
    // Use data...
}
```

---

### 2.3 Memory Profiling

#### 🔴 Q14: How do you identify memory leaks in Go?

**Answer:**

**1. pprof heap profile:**
```go
import _ "net/http/pprof"

func main() {
    go func() {
        http.ListenAndServe(":6060", nil)
    }()
    // Application code...
}
```

```bash
# Capture heap profile
go tool pprof http://localhost:6060/debug/pprof/heap

# Commands in pprof
(pprof) top10           # Top memory consumers
(pprof) list funcName   # Source-level breakdown
(pprof) web             # Visual graph
```

**2. Compare profiles over time:**
```bash
# Take baseline
curl http://localhost:6060/debug/pprof/heap > base.pprof

# After some time
curl http://localhost:6060/debug/pprof/heap > current.pprof

# Compare
go tool pprof -base=base.pprof current.pprof
```

**3. Common leak patterns:**
```go
// Leak: Goroutine leak
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // Blocked forever if ch never receives
        fmt.Println(val)
    }()
    // ch never closed, goroutine leaks
}

// Leak: Slice capacity not released
func sliceLeak(data []byte) []byte {
    return data[:10] // Still holds reference to full backing array
}

// Fix: Copy to release original
func sliceFix(data []byte) []byte {
    result := make([]byte, 10)
    copy(result, data[:10])
    return result
}
```

---

## 3. Goroutines & Concurrency (CRITICAL)

> **This is the most important section for senior Go interviews.**

### 3.1 Goroutine Fundamentals

#### 🟢 Q15: What is a goroutine and how is it different from OS threads?

**Answer:**

| Aspect | Goroutine | OS Thread |
|--------|-----------|-----------|
| **Stack size** | 2KB initial (growable) | 1-8MB fixed |
| **Creation cost** | ~300 bytes, ~1μs | ~1MB, ~1ms |
| **Scheduling** | Go runtime (cooperative) | OS kernel (preemptive) |
| **Context switch** | ~tens of ns | ~μs |
| **Max count** | Millions feasible | Thousands practical |

```go
// Creating goroutine - trivially cheap
go func() {
    // Concurrent execution
}()
```

---

#### 🟡 Q16: Explain the M:N scheduling model in Go.

**Answer:**

**Components:**
- **G** (Goroutine) — User-space thread
- **M** (Machine) — OS thread
- **P** (Processor) — Logical processor, scheduling context

**How it works:**
```
[G] [G] [G]     <- Goroutines
    |
   [P]          <- Logical Processor (holds runqueue)
    |
   [M]          <- OS Thread
    |
  [CPU]         <- Hardware
```

**Key behaviors:**
- Each **P** has a local run queue of **G**s
- **M** must acquire a **P** to run **G**s
- When **G** blocks (syscall), **M** releases **P** for other **M**s
- **Work stealing:** idle **P** steals from busy **P**'s queue

```go
runtime.GOMAXPROCS(4) // 4 Ps, meaning 4 Gs can run in parallel
```

---

#### 🔴 Q17: What happens when a goroutine makes a blocking syscall?

**Answer:**

**Scenario: Goroutine blocks on syscall**

1. **G1** running on **M1** with **P1** makes blocking syscall
2. Runtime detects block, **M1** releases **P1**
3. **P1** is handed to another **M** (or new **M** is created)
4. **P1** continues running other goroutines
5. When syscall completes, **G1** is re-queued

```
Before:               After syscall:
G1 ─ M1 ─ P1         G1 ─ M1 (blocked on syscall)
                     G2 ─ M2 ─ P1 (P1 reused)
```

**Network poller (optimization):**
- Network I/O uses non-blocking syscalls + epoll/kqueue
- Goroutine doesn't block **M** on network I/O

---

#### 🔴 Q18: What is goroutine preemption and how did it change in Go 1.14?

**Answer:**

**Pre-Go 1.14:** Cooperative preemption only
- Goroutines yield at function calls
- Tight loops without function calls could starve other goroutines

```go
// This could starve other goroutines pre-1.14
for {
    sum += i // No function call, no preemption point
}
```

**Go 1.14+:** Asynchronous preemption
- Runtime uses signals (SIGURG on Unix) to preempt
- Any goroutine can be preempted at almost any safe point
- Prevents starvation from tight loops

```go
// Now preemptible even without function calls
for {
    sum += i // Runtime can still preempt via signal
}
```

---

### 3.2 Channels Deep Dive

#### 🟡 Q19: Buffered vs Unbuffered channels — when to use which?

**Answer:**

| Channel Type | Behavior | Use Case |
|--------------|----------|----------|
| **Unbuffered** | Synchronous, send blocks until receive | Signaling, synchronization |
| **Buffered** | Asynchronous until full | Decoupling, batching, rate limiting |

```go
// Unbuffered: Handshake semantics
done := make(chan struct{})
go func() {
    work()
    done <- struct{}{} // Blocks until main receives
}()
<-done // Sync point

// Buffered: Rate limiting
semaphore := make(chan struct{}, 10) // Max 10 concurrent
for _, item := range items {
    semaphore <- struct{}{} // Blocks if 10 in flight
    go func(item Item) {
        defer func() { <-semaphore }()
        process(item)
    }(item)
}
```

---

#### 🔴 Q20: Explain channel closing rules and best practices.

**Answer:**

**Rules:**
1. Only **sender** should close channel
2. Closing already-closed channel **panics**
3. Sending on closed channel **panics**
4. Receiving from closed channel returns zero value immediately

```go
// Pattern: Signal completion
func producer(ch chan<- int) {
    defer close(ch) // Sender closes
    for i := 0; i < 10; i++ {
        ch <- i
    }
}

func consumer(ch <-chan int) {
    for v := range ch { // Exits when channel closed
        process(v)
    }
}

// Detect closed channel
v, ok := <-ch
if !ok {
    // Channel closed
}
```

**Anti-pattern: Closing from receiver**
```go
// WRONG: Multiple senders, receiver closes
func bad() {
    ch := make(chan int)
    go sender1(ch)
    go sender2(ch)
    go func() {
        // DON'T DO THIS - senders will panic
        close(ch)
    }()
}
```

**Safe multi-sender close:**
```go
func safeClose() {
    ch := make(chan int)
    done := make(chan struct{})
    
    // Multiple senders
    for i := 0; i < 3; i++ {
        go func() {
            for {
                select {
                case <-done:
                    return // Exit on done signal
                case ch <- produce():
                }
            }
        }()
    }
    
    // Signal shutdown (instead of closing ch)
    close(done)
}
```

---

#### 🔴 Q21: Implement fan-in and fan-out patterns.

**Answer:**

**Fan-out: One channel → Multiple workers**
```go
func fanOut(input <-chan Job, workers int) {
    for i := 0; i < workers; i++ {
        go func(workerID int) {
            for job := range input {
                result := process(job)
                fmt.Printf("Worker %d processed: %v\n", workerID, result)
            }
        }(i)
    }
}
```

**Fan-in: Multiple channels → One channel**
```go
func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}

// Usage
ch1, ch2, ch3 := producer1(), producer2(), producer3()
merged := fanIn(ch1, ch2, ch3)
for v := range merged {
    process(v)
}
```

**Worker Pool Pattern:**
```go
func workerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
    close(results)
}
```

---

### 3.3 Select Statement

#### 🟡 Q22: How does select work with multiple ready channels?

**Answer:**

When multiple cases are ready, **select chooses one at random** (pseudo-random).

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)
ch1 <- 1
ch2 <- 2

select {
case v := <-ch1:
    fmt.Println("ch1:", v) // May print
case v := <-ch2:
    fmt.Println("ch2:", v) // May print
}
// Order is NOT deterministic
```

**Why random?** Prevents starvation — no channel is favored.

---

#### 🔴 Q23: Implement timeout and non-blocking channel patterns.

**Answer:**

**Timeout pattern:**
```go
func fetchWithTimeout(url string, timeout time.Duration) ([]byte, error) {
    result := make(chan []byte, 1)
    errCh := make(chan error, 1)
    
    go func() {
        data, err := fetch(url)
        if err != nil {
            errCh <- err
            return
        }
        result <- data
    }()
    
    select {
    case data := <-result:
        return data, nil
    case err := <-errCh:
        return nil, err
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout after %v", timeout)
    }
}
```

**Non-blocking send/receive:**
```go
// Non-blocking receive
select {
case msg := <-ch:
    process(msg)
default:
    // Channel empty, don't block
}

// Non-blocking send
select {
case ch <- msg:
    // Sent
default:
    // Channel full, drop or log
    log.Println("channel full, dropping message")
}
```

**Priority select (handle urgent first):**
```go
func prioritySelect(urgent, normal <-chan Event) {
    for {
        // First, drain all urgent
        select {
        case e := <-urgent:
            handleUrgent(e)
            continue
        default:
        }
        
        // Then handle any
        select {
        case e := <-urgent:
            handleUrgent(e)
        case e := <-normal:
            handleNormal(e)
        }
    }
}
```

---

### 3.4 Synchronization Primitives

#### 🟡 Q24: When to use Mutex vs RWMutex?

**Answer:**

| Type | Use Case | Behavior |
|------|----------|----------|
| `sync.Mutex` | Exclusive access | One goroutine at a time |
| `sync.RWMutex` | Read-heavy workloads | Multiple readers OR one writer |

```go
// Mutex: Simple exclusive access
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

// RWMutex: Many readers, few writers
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()         // Multiple readers allowed
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()          // Exclusive write
    defer c.mu.Unlock()
    c.data[key] = value
}
```

**When NOT to use RWMutex:**
- Write frequency > 10-20% — RWMutex overhead not worth it
- Critical section is very short — Mutex may be faster

---

#### 🔴 Q25: Explain sync.Once and its use cases.

**Answer:**

`sync.Once` ensures a function executes **exactly once**, even across goroutines.

```go
var (
    instance *Database
    once     sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        instance = connectDatabase() // Runs exactly once
    })
    return instance
}
```

**Internals:**
- Uses atomic operations + mutex
- First call executes function
- Subsequent calls return immediately (no lock contention)

**Common use cases:**
- Singleton initialization
- One-time configuration loading
- Lazy initialization

**Gotcha: Deadlock**
```go
var once sync.Once

func init() {
    once.Do(func() {
        once.Do(func() { // DEADLOCK: once is already held
            fmt.Println("never prints")
        })
    })
}
```

---

#### 🔴 Q26: How does sync.WaitGroup work internally?

**Answer:**

```go
var wg sync.WaitGroup

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        work(id)
    }(i)
}

wg.Wait() // Blocks until counter == 0
```

**Internal state:**
- Counter (int64) — tracks outstanding work
- Wait semaphore — blocks Wait() callers

**Rules:**
1. `Add()` must happen before `Wait()` starts
2. `Add()` with negative delta should not go below zero
3. Reuse: counter must be zero before reuse

**Common mistake:**
```go
// WRONG: Add inside goroutine
for i := 0; i < 10; i++ {
    go func() {
        wg.Add(1) // Race: Wait() might start before Add()
        defer wg.Done()
        work()
    }()
}
wg.Wait()
```

---

#### 🔴 Q27: When would you use sync.Cond?

**Answer:**

`sync.Cond` enables goroutines to wait for a condition to become true.

```go
type Queue struct {
    items []int
    cond  *sync.Cond
    mu    sync.Mutex
}

func NewQueue() *Queue {
    q := &Queue{}
    q.cond = sync.NewCond(&q.mu)
    return q
}

func (q *Queue) Enqueue(item int) {
    q.mu.Lock()
    defer q.mu.Unlock()
    q.items = append(q.items, item)
    q.cond.Signal() // Wake one waiter
}

func (q *Queue) Dequeue() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    for len(q.items) == 0 {
        q.cond.Wait() // Releases lock, waits, reacquires lock
    }
    item := q.items[0]
    q.items = q.items[1:]
    return item
}
```

**Use cases:**
- Producer-consumer with blocking wait
- Waiting for specific state change
- Broadcasting state updates (`Broadcast()`)

**Prefer channels when:**
- Simple signaling
- You need select semantics
- Timeout handling

---

### 3.5 Context Package (VERY IMPORTANT)

#### 🟡 Q28: Explain the purpose and usage of context.Context.

**Answer:**

**Context provides:**
1. **Cancellation propagation** — Cancel downstream operations
2. **Timeouts** — Deadline-based cancellation
3. **Request-scoped values** — Pass metadata (use sparingly)

```go
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // Request context
    
    // Pass context to downstream calls
    result, err := fetchData(ctx)
    if err != nil {
        if ctx.Err() == context.Canceled {
            return // Client disconnected
        }
        http.Error(w, err.Error(), 500)
        return
    }
    
    json.NewEncoder(w).Encode(result)
}

func fetchData(ctx context.Context) (*Data, error) {
    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
    resp, err := http.DefaultClient.Do(req)
    // If ctx is cancelled, request aborts
    return parseResponse(resp)
}
```

---

#### 🔴 Q29: Implement proper context cancellation in a service.

**Answer:**

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    // Create child context with timeout
    ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
    defer cancel() // Always call cancel to release resources
    
    // Run steps with context checking
    if err := validateOrder(ctx, orderID); err != nil {
        return err
    }
    
    if err := chargePayment(ctx, orderID); err != nil {
        return err
    }
    
    if err := fulfillOrder(ctx, orderID); err != nil {
        return err
    }
    
    return nil
}

func validateOrder(ctx context.Context, orderID string) error {
    // Check context before expensive operation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    
    // Do validation...
    return nil
}

// Goroutine respecting cancellation
func pollUntilReady(ctx context.Context) error {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err() // Cancelled or deadline exceeded
        case <-ticker.C:
            if ready, _ := checkStatus(); ready {
                return nil
            }
        }
    }
}
```

---

#### 🔴 Q30: What are the context package antipatterns?

**Answer:**

**❌ Antipattern 1: Storing large objects in context**
```go
// BAD: Context is not for data storage
ctx = context.WithValue(ctx, "user", largeUserObject)

// GOOD: Store minimal identifiers, fetch data when needed
ctx = context.WithValue(ctx, userIDKey, "user-123")
```

**❌ Antipattern 2: Using string keys**
```go
// BAD: Collision risk
ctx = context.WithValue(ctx, "requestID", id)

// GOOD: Use typed unexported keys
type contextKey string
const requestIDKey contextKey = "requestID"
ctx = context.WithValue(ctx, requestIDKey, id)
```

**❌ Antipattern 3: Passing nil context**
```go
// BAD
fetch(nil, url)

// GOOD: Use context.Background() or context.TODO()
fetch(context.Background(), url)
```

**❌ Antipattern 4: Not propagating context**
```go
// BAD: Creates orphan context
func handler(ctx context.Context) {
    go func() {
        doWork(context.Background()) // Ignores parent cancellation
    }()
}

// GOOD: Propagate context
func handler(ctx context.Context) {
    go func() {
        doWork(ctx) // Respects parent cancellation
    }()
}
```

**❌ Antipattern 5: Not calling cancel**
```go
// BAD: Resource leak
ctx, _ := context.WithTimeout(parent, 10*time.Second)

// GOOD: Always defer cancel
ctx, cancel := context.WithTimeout(parent, 10*time.Second)
defer cancel()
```

---

### 3.6 Concurrency Code Challenges

#### ⚫ Q31: Implement a rate limiter using channels.

**Answer:**

```go
type RateLimiter struct {
    tokens chan struct{}
    done   chan struct{}
}

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens: make(chan struct{}, rate),
        done:   make(chan struct{}),
    }
    
    // Fill initial tokens
    for i := 0; i < rate; i++ {
        rl.tokens <- struct{}{}
    }
    
    // Refill tokens
    interval := per / time.Duration(rate)
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-rl.done:
                return
            case <-ticker.C:
                select {
                case rl.tokens <- struct{}{}:
                default: // Bucket full
                }
            }
        }
    }()
    
    return rl
}

func (rl *RateLimiter) Allow() bool {
    select {
    case <-rl.tokens:
        return true
    default:
        return false
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

func (rl *RateLimiter) Stop() {
    close(rl.done)
}
```

---

#### ⚫ Q32: Implement a worker pool with graceful shutdown.

**Answer:**

```go
type WorkerPool struct {
    jobs    chan Job
    results chan Result
    done    chan struct{}
    wg      sync.WaitGroup
}

func NewWorkerPool(numWorkers, jobBuffer int) *WorkerPool {
    pool := &WorkerPool{
        jobs:    make(chan Job, jobBuffer),
        results: make(chan Result, jobBuffer),
        done:    make(chan struct{}),
    }
    
    for i := 0; i < numWorkers; i++ {
        pool.wg.Add(1)
        go pool.worker(i)
    }
    
    return pool
}

func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()
    for {
        select {
        case <-p.done:
            return
        case job, ok := <-p.jobs:
            if !ok {
                return
            }
            result := process(job)
            select {
            case p.results <- result:
            case <-p.done:
                return
            }
        }
    }
}

func (p *WorkerPool) Submit(job Job) error {
    select {
    case p.jobs <- job:
        return nil
    case <-p.done:
        return errors.New("pool is shutting down")
    }
}

func (p *WorkerPool) Results() <-chan Result {
    return p.results
}

func (p *WorkerPool) Shutdown(ctx context.Context) error {
    close(p.done)    // Signal workers to stop
    close(p.jobs)    // No more jobs
    
    // Wait with timeout
    done := make(chan struct{})
    go func() {
        p.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        close(p.results)
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

#### ⚫ Q33: Implement a concurrent-safe LRU cache.

**Answer:**

```go
type LRUCache struct {
    capacity int
    items    map[string]*list.Element
    order    *list.List
    mu       sync.RWMutex
}

type entry struct {
    key   string
    value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        items:    make(map[string]*list.Element),
        order:    list.New(),
    }
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.items[key]; ok {
        c.order.MoveToFront(elem)
        return elem.Value.(*entry).value, true
    }
    return nil, false
}

func (c *LRUCache) Set(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.items[key]; ok {
        c.order.MoveToFront(elem)
        elem.Value.(*entry).value = value
        return
    }
    
    // Evict if at capacity
    if c.order.Len() >= c.capacity {
        oldest := c.order.Back()
        if oldest != nil {
            c.order.Remove(oldest)
            delete(c.items, oldest.Value.(*entry).key)
        }
    }
    
    elem := c.order.PushFront(&entry{key: key, value: value})
    c.items[key] = elem
}

func (c *LRUCache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.items[key]; ok {
        c.order.Remove(elem)
        delete(c.items, key)
    }
}
```

---

## 4. Go Data Structures & Performance

### 4.1 Slices Internals

#### 🟡 Q34: Explain slice internals and capacity growth.

**Answer:**

**Slice header (24 bytes on 64-bit):**
```go
type slice struct {
    array unsafe.Pointer // Pointer to underlying array
    len   int            // Number of elements
    cap   int            // Capacity of underlying array
}
```

**Capacity growth strategy (Go 1.18+):**
- If cap < 256: double
- If cap >= 256: grow by 25% + 192 (smoothed growth)

```go
s := make([]int, 0)
for i := 0; i < 10; i++ {
    s = append(s, i)
    fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
}
// Output shows capacity: 1, 2, 4, 8, 8, 8, 8, 8, 16, 16
```

**Pre-allocation for performance:**
```go
// BAD: Multiple reallocations
var s []int
for i := 0; i < 10000; i++ {
    s = append(s, i)
}

// GOOD: Single allocation
s := make([]int, 0, 10000)
for i := 0; i < 10000; i++ {
    s = append(s, i)
}
```

---

#### 🔴 Q35: Why does nil slice vs empty slice matter?

**Answer:**

```go
var nilSlice []int          // nil slice: nil pointer, len=0, cap=0
emptySlice := []int{}       // empty slice: non-nil pointer, len=0, cap=0
emptyMake := make([]int, 0) // empty slice: non-nil pointer, len=0, cap=0

fmt.Println(nilSlice == nil)  // true
fmt.Println(emptySlice == nil) // false
```

**When it matters:**

1. **JSON marshaling:**
```go
type Response struct {
    Items []string `json:"items"`
}

// nil slice
r1 := Response{}
json.Marshal(r1) // {"items":null}

// empty slice
r2 := Response{Items: []string{}}
json.Marshal(r2) // {"items":[]}
```

2. **Reflection and deep equality:**
```go
reflect.DeepEqual(nilSlice, emptySlice) // false
```

3. **API contracts:** Some APIs expect non-nil slice.

**Best practice:** Initialize slices if you'll return them in JSON APIs.

---

#### 🔴 Q36: What is the slice gotcha with append?

**Answer:**

```go
original := []int{1, 2, 3, 4, 5}
slice1 := original[:3]  // [1, 2, 3], shares backing array
slice2 := append(slice1, 100)

fmt.Println(original) // [1, 2, 3, 100, 5] - MODIFIED!
fmt.Println(slice1)   // [1, 2, 3]
fmt.Println(slice2)   // [1, 2, 3, 100]
```

**Fix: Use full slice expression:**
```go
original := []int{1, 2, 3, 4, 5}
slice1 := original[:3:3]  // [1, 2, 3], cap=3 (third index limits capacity)
slice2 := append(slice1, 100) // Forces new allocation

fmt.Println(original) // [1, 2, 3, 4, 5] - UNCHANGED
fmt.Println(slice2)   // [1, 2, 3, 100]
```

**Or copy explicitly:**
```go
slice1 := make([]int, 3)
copy(slice1, original[:3])
```

---

### 4.2 Maps Internals

#### 🔴 Q37: How do Go maps work internally?

**Answer:**

**Structure:**
- Hash table with buckets
- Each bucket holds 8 key-value pairs
- Buckets linked for overflow

```go
// Simplified internal structure
type hmap struct {
    count     int            // Number of elements
    B         uint8          // log2 of number of buckets
    buckets   unsafe.Pointer // Array of 2^B buckets
    oldbuckets unsafe.Pointer // For incremental growth
    // ...
}

type bmap struct {
    tophash [8]uint8  // Top 8 bits of hash for quick comparison
    keys    [8]keyType
    values  [8]valueType
    overflow *bmap
}
```

**Operations:**
1. **Lookup:** Hash key → find bucket → compare tophash → compare full key
2. **Insert:** Same as lookup, insert if not found
3. **Growth:** When load factor > 6.5, double buckets (incremental evacuation)

---

#### 🔴 Q38: Why is map not thread-safe and how to handle it?

**Answer:**

**Why not thread-safe:**
- Performance: synchronization has overhead
- Not all use cases need concurrency
- Go philosophy: don't pay for what you don't use

**Option 1: sync.RWMutex**
```go
type SafeMap struct {
    mu   sync.RWMutex
    data map[string]int
}

func (m *SafeMap) Get(key string) int {
    m.mu.RLock()
    defer m.mu.RUnlock()
    return m.data[key]
}

func (m *SafeMap) Set(key string, val int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = val
}
```

**Option 2: sync.Map (specialized cases)**
```go
var m sync.Map

m.Store("key", "value")
val, ok := m.Load("key")
m.Delete("key")
m.Range(func(k, v interface{}) bool {
    fmt.Println(k, v)
    return true // continue iteration
})
```

**When to use sync.Map:**
- Key set is stable (few writes, many reads)
- Disjoint key sets across goroutines
- NOT for general-purpose concurrent map

---

### 4.3 Struct Alignment & Memory Padding

#### 🔴 Q39: Explain struct memory alignment and how to optimize it.

**Answer:**

**Alignment rules:**
- Fields aligned to their size (int64 on 8-byte boundary)
- Struct size rounded to largest field alignment

```go
// BAD: 24 bytes due to padding
type BadOrder struct {
    a bool   // 1 byte + 7 padding
    b int64  // 8 bytes
    c bool   // 1 byte + 7 padding
}

// GOOD: 16 bytes
type GoodOrder struct {
    b int64 // 8 bytes
    a bool  // 1 byte
    c bool  // 1 byte + 6 padding
}
```

**Check struct size:**
```go
fmt.Println(unsafe.Sizeof(BadOrder{}))  // 24
fmt.Println(unsafe.Sizeof(GoodOrder{})) // 16
```

**Tool: fieldalignment**
```bash
go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
fieldalignment -fix ./...
```

---

### 4.4 Value vs Pointer Receivers

#### 🟡 Q40: When to use value vs pointer receivers?

**Answer:**

| Use Pointer Receiver | Use Value Receiver |
|---------------------|-------------------|
| Method modifies receiver | Method doesn't modify state |
| Large struct (avoid copy) | Small struct (≤ 2-3 words) |
| Struct contains sync.Mutex | Consistency with other methods |
| Any method has pointer receiver | |

```go
// Pointer receiver: Modifies state
func (u *User) SetName(name string) {
    u.Name = name
}

// Value receiver: No modification, small struct
func (p Point) Distance(other Point) float64 {
    return math.Sqrt(math.Pow(p.X-other.X, 2) + math.Pow(p.Y-other.Y, 2))
}
```

**Consistency rule:** If one method needs pointer receiver, use pointer for all methods.

---

### 4.5 Interface Internals

#### 🔴 Q41: Explain interface representation in Go.

**Answer:**

**Interface value = (type, value) pair**

```go
// Empty interface
type eface struct {
    _type *_type        // Type information
    data  unsafe.Pointer // Pointer to actual data
}

// Non-empty interface
type iface struct {
    tab  *itab          // Type + method table
    data unsafe.Pointer // Pointer to actual data
}
```

**itab structure:**
```go
type itab struct {
    inter *interfacetype // Interface type
    _type *_type         // Concrete type
    fun   [1]uintptr     // Method pointers (variable size)
}
```

**Interface nil gotcha:**
```go
type MyError struct{}
func (e *MyError) Error() string { return "error" }

func returnsError() error {
    var err *MyError = nil
    return err // Returns (type=*MyError, value=nil)
}

func main() {
    err := returnsError()
    fmt.Println(err == nil) // false! Interface is not nil
}
```

**Fix:**
```go
func returnsError() error {
    var err *MyError = nil
    if err == nil {
        return nil // Return actual nil interface
    }
    return err
}
```

---

#### 🔴 Q42: What is the cost of interface{} and type assertions?

**Answer:**

**Cost of interface{}:**
1. **Boxing:** Value types must be allocated on heap and wrapped
2. **Memory:** 16 bytes for interface value (on 64-bit)
3. **Indirection:** Extra pointer dereference

```go
// Interface causes heap allocation
func process(v interface{}) {
    // v is boxed
}

x := 42
process(x) // x escapes to heap
```

**Type assertion cost:**
```go
var v interface{} = "hello"

// Type assertion (runtime check)
s, ok := v.(string) // O(1) but has cost

// Type switch (multiple checks)
switch v := v.(type) {
case string:   // Check 1
case int:      // Check 2
}
```

**Optimization: Use generics (Go 1.18+):**
```go
// Before: interface{} + type assertion
func maxInt(a, b interface{}) interface{} {
    if a.(int) > b.(int) {
        return a
    }
    return b
}

// After: Generic, no boxing
func max[T constraints.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}
```

---

### 4.6 Tricky Questions

#### 🔴 Q43: What is the zero value of a channel and what happens when you use it?

**Answer:**

```go
var ch chan int // Zero value: nil

// Operations on nil channel:
ch <- 1    // Blocks forever (deadlock if only goroutine)
<-ch       // Blocks forever
close(ch)  // Panic!
```

**Use case for nil channel: Disable select case**
```go
func merge(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for ch1 != nil || ch2 != nil {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil // Disable this case
                    continue
                }
                out <- v
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil // Disable this case
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}
```

---

#### 🔴 Q44: What is the difference between make and new?

**Answer:**

| `make` | `new` |
|--------|-------|
| Only for slice, map, channel | Any type |
| Returns initialized value | Returns pointer to zeroed memory |
| `make([]int, 5)` → `[]int` | `new([]int)` → `*[]int` (nil slice pointer) |

```go
// make: Initialize internal structure
s := make([]int, 5)     // []int with len=5, cap=5
m := make(map[string]int) // Initialized map
ch := make(chan int, 10)  // Buffered channel

// new: Allocate zeroed memory, return pointer
p := new(int)  // *int pointing to 0
s := new([]int) // *[]int pointing to nil slice
```

**Equivalent without make/new:**
```go
// These are equivalent:
p := new(int)
var v int; p := &v

// These are NOT equivalent:
m := make(map[string]int) // Usable map
m := new(map[string]int)  // *map, points to nil map (unusable)
```

---

#### ⚫ Q45: Explain what happens in this code:

```go
func main() {
    m := make(map[int]*int)
    for i := 0; i < 5; i++ {
        m[i] = &i
    }
    for k, v := range m {
        fmt.Println(k, *v)
    }
}
```

**Answer:**

All map values point to the same variable `i`, which has final value 5.

**Output (order varies):**
```
0 5
1 5
2 5
3 5
4 5
```

**Explanation:**
- `i` is a single variable, reused each iteration
- `&i` stores address of same variable
- After loop, `i == 5`

**Fix:**
```go
for i := 0; i < 5; i++ {
    i := i // Shadow with new variable
    m[i] = &i
}

// Or in Go 1.22+: Loop variables are per-iteration by default
```

---

## Part 1 Summary Checklist

### Go Fundamentals
- [ ] Go design philosophy: simplicity, composition, explicitness
- [ ] Why no inheritance (fragile base class, diamond problem)
- [ ] Errors as values (explicit handling, no hidden control flow)
- [ ] Compilation model (fast, single pass, no circular deps)
- [ ] Cross-compilation (GOOS, GOARCH)
- [ ] Runtime components (scheduler, GC, allocator)

### Memory Management
- [ ] Stack vs heap allocation decision (escape analysis)
- [ ] Escape analysis flags (`go build -gcflags="-m"`)
- [ ] GC: Concurrent, tri-color, mark-sweep
- [ ] Why non-generational GC
- [ ] Reducing GC pressure (pools, pre-allocation, GOGC)
- [ ] Memory profiling with pprof

### Concurrency
- [ ] Goroutines vs threads (stack size, cost, scheduling)
- [ ] M:N scheduling (G, M, P)
- [ ] Blocking syscalls handling
- [ ] Preemption (Go 1.14+ async preemption)
- [ ] Buffered vs unbuffered channels
- [ ] Channel closing rules (sender closes)
- [ ] Fan-in/fan-out patterns
- [ ] Select statement (random selection, non-blocking)
- [ ] Mutex vs RWMutex (use cases)
- [ ] sync.Once, sync.WaitGroup, sync.Cond
- [ ] Context: cancellation, timeout, values
- [ ] Context antipatterns

### Data Structures
- [ ] Slice internals (header, capacity growth)
- [ ] nil slice vs empty slice (JSON marshaling)
- [ ] Slice append gotcha (shared backing array)
- [ ] Map internals (buckets, growth)
- [ ] Maps not thread-safe (sync.RWMutex, sync.Map)
- [ ] Struct alignment and padding optimization
- [ ] Value vs pointer receivers
- [ ] Interface internals (iface, eface)
- [ ] Interface nil gotcha

---

> **Continue to Part 2:** Error Handling, Modules, Testing, Production APIs, Database & Caching
