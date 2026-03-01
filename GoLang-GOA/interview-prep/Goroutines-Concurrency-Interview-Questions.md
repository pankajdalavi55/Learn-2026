# Goroutines & Concurrency Interview Questions by Experience Level

> **Complete Concurrency Deep-Dive for Go Developers**  
> **Organized by Experience Level: Junior → Staff/Principal**

---

## Table of Contents

1. [Junior Level (0-2 Years)](#1-junior-level-0-2-years)
2. [Mid-Level (2-5 Years)](#2-mid-level-2-5-years)
3. [Senior Level (5-8 Years)](#3-senior-level-5-8-years)
4. [Staff/Principal Level (8+ Years)](#4-staffprincipal-level-8-years)
5. [Concurrency Coding Challenges](#5-concurrency-coding-challenges)
6. [Quick Reference & Cheat Sheet](#6-quick-reference--cheat-sheet)

---

## 1. Junior Level (0-2 Years)

### 1.1 Goroutine Basics

#### Q1: What is a goroutine?

**Answer:**
A goroutine is a lightweight thread managed by the Go runtime. It's a function that runs concurrently with other functions.

```go
// Creating a goroutine
go func() {
    fmt.Println("Hello from goroutine")
}()

// Or with a named function
go myFunction()
```

**Key points:**
- Starts with `go` keyword
- Very cheap to create (~2KB stack)
- Managed by Go runtime, not OS
- Can have millions of goroutines

---

#### Q2: What happens if main() exits before goroutine completes?

**Answer:**
The program terminates immediately, and all goroutines are killed without completing their work.

```go
func main() {
    go func() {
        time.Sleep(time.Second)
        fmt.Println("This may never print!")
    }()
    // main exits immediately, goroutine is killed
}
```

**Fix: Wait for goroutine to complete:**
```go
func main() {
    done := make(chan bool)
    
    go func() {
        time.Sleep(time.Second)
        fmt.Println("This will print!")
        done <- true
    }()
    
    <-done // Wait for goroutine
}
```

---

#### Q3: How do you create a channel?

**Answer:**
Channels are created using `make`:

```go
// Unbuffered channel
ch := make(chan int)

// Buffered channel with capacity 10
bufferedCh := make(chan int, 10)

// Receive-only channel (type)
var recvOnly <-chan int

// Send-only channel (type)
var sendOnly chan<- int
```

---

#### Q4: What is the difference between buffered and unbuffered channels?

**Answer:**

| Unbuffered | Buffered |
|------------|----------|
| `make(chan int)` | `make(chan int, 5)` |
| Blocks until receiver ready | Blocks only when buffer full |
| Synchronous communication | Asynchronous until full |
| Guarantees delivery before send returns | May queue messages |

```go
// Unbuffered - blocks until someone receives
ch := make(chan int)
go func() { ch <- 42 }() // Would block without receiver
val := <-ch              // Receives 42

// Buffered - can send without immediate receiver
buffered := make(chan int, 2)
buffered <- 1 // Doesn't block
buffered <- 2 // Doesn't block
buffered <- 3 // BLOCKS - buffer full
```

---

#### Q5: How do you iterate over a channel?

**Answer:**
Use `range` to iterate until channel is closed:

```go
func main() {
    ch := make(chan int)
    
    go func() {
        for i := 0; i < 5; i++ {
            ch <- i
        }
        close(ch) // Must close for range to exit
    }()
    
    for val := range ch {
        fmt.Println(val) // Prints 0, 1, 2, 3, 4
    }
}
```

---

#### Q6: What is sync.WaitGroup and how do you use it?

**Answer:**
`sync.WaitGroup` waits for a collection of goroutines to finish.

```go
func main() {
    var wg sync.WaitGroup
    
    for i := 0; i < 5; i++ {
        wg.Add(1) // Increment counter
        go func(id int) {
            defer wg.Done() // Decrement counter when done
            fmt.Printf("Worker %d done\n", id)
        }(i)
    }
    
    wg.Wait() // Block until counter is 0
    fmt.Println("All workers completed")
}
```

**Rules:**
- Call `Add()` before starting goroutine
- Call `Done()` when goroutine completes (usually with `defer`)
- Call `Wait()` to block until all done

---

#### Q7: What happens when you send to a closed channel?

**Answer:**
**Panic!** Sending to a closed channel causes a runtime panic.

```go
ch := make(chan int)
close(ch)
ch <- 1 // PANIC: send on closed channel
```

**Receiving from closed channel is safe:**
```go
ch := make(chan int)
close(ch)
val, ok := <-ch
fmt.Println(val, ok) // 0, false (zero value, not ok)
```

---

#### Q8: How do you check if a channel is closed?

**Answer:**
Use the two-value receive form:

```go
val, ok := <-ch
if !ok {
    fmt.Println("Channel is closed")
} else {
    fmt.Println("Received:", val)
}
```

---

### 1.2 Basic Synchronization

#### Q9: What is a mutex and when do you use it?

**Answer:**
A mutex (mutual exclusion) protects shared data from concurrent access.

```go
type Counter struct {
    mu    sync.Mutex
    count int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

**Use mutex when:**
- Multiple goroutines access shared data
- At least one goroutine modifies the data

---

#### Q10: What is a data race?

**Answer:**
A data race occurs when two or more goroutines access the same memory location concurrently, and at least one access is a write.

```go
// DATA RACE!
var counter int

func main() {
    for i := 0; i < 1000; i++ {
        go func() {
            counter++ // Multiple goroutines reading and writing
        }()
    }
    time.Sleep(time.Second)
    fmt.Println(counter) // Unpredictable result!
}
```

**Detect with race detector:**
```bash
go run -race main.go
```

---

## 2. Mid-Level (2-5 Years)

### 2.1 Channel Patterns

#### Q11: Explain the select statement.

**Answer:**
`select` lets you wait on multiple channel operations:

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case ch3 <- value:
    fmt.Println("Sent to ch3")
default:
    fmt.Println("No channel ready")
}
```

**Key behaviors:**
- Blocks until one case is ready
- If multiple ready, chooses randomly
- `default` makes it non-blocking

---

#### Q12: How do you implement a timeout with channels?

**Answer:**

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

---

#### Q13: What is the fan-out pattern?

**Answer:**
Fan-out distributes work from one channel to multiple workers:

```go
func fanOut(input <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    
    for i := 0; i < numWorkers; i++ {
        outputs[i] = worker(input)
    }
    
    return outputs
}

func worker(input <-chan int) <-chan int {
    output := make(chan int)
    go func() {
        defer close(output)
        for n := range input {
            output <- process(n)
        }
    }()
    return output
}
```

---

#### Q14: What is the fan-in pattern?

**Answer:**
Fan-in merges multiple channels into one:

```go
func fanIn(channels ...<-chan int) <-chan int {
    merged := make(chan int)
    var wg sync.WaitGroup
    
    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                merged <- val
            }
        }(ch)
    }
    
    go func() {
        wg.Wait()
        close(merged)
    }()
    
    return merged
}
```

---

#### Q15: How do you implement a worker pool?

**Answer:**

```go
func workerPool(numWorkers int, jobs <-chan Job, results chan<- Result) {
    var wg sync.WaitGroup
    
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                result := processJob(job)
                results <- result
            }
        }(i)
    }
    
    wg.Wait()
    close(results)
}

// Usage
func main() {
    jobs := make(chan Job, 100)
    results := make(chan Result, 100)
    
    // Start workers
    go workerPool(5, jobs, results)
    
    // Send jobs
    for i := 0; i < 50; i++ {
        jobs <- Job{ID: i}
    }
    close(jobs)
    
    // Collect results
    for result := range results {
        fmt.Println(result)
    }
}
```

---

### 2.2 Context Package

#### Q16: What is context.Context and why is it important?

**Answer:**
`context.Context` carries deadlines, cancellation signals, and request-scoped values across API boundaries.

```go
func handleRequest(ctx context.Context) error {
    // Check if already cancelled
    if ctx.Err() != nil {
        return ctx.Err()
    }
    
    // Pass context to downstream calls
    result, err := queryDatabase(ctx)
    if err != nil {
        return err
    }
    
    return processResult(ctx, result)
}
```

**Why important:**
- Propagates cancellation across goroutines
- Prevents resource leaks
- Sets deadlines for operations
- Standard way to handle request lifecycle

---

#### Q17: What are the different context types?

**Answer:**

```go
// Root contexts (never cancelled)
ctx := context.Background() // Main function, initialization
ctx := context.TODO()       // Placeholder when unsure

// Derived contexts
ctx, cancel := context.WithCancel(parent)
defer cancel() // Always call cancel

ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

ctx, cancel := context.WithDeadline(parent, time.Now().Add(time.Hour))
defer cancel()

ctx := context.WithValue(parent, key, value)
```

---

#### Q18: How do you properly use context cancellation?

**Answer:**

```go
func processItems(ctx context.Context, items []Item) error {
    for _, item := range items {
        // Check cancellation before each item
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }
        
        if err := processItem(ctx, item); err != nil {
            return err
        }
    }
    return nil
}

// Long-running operation with cancellation
func longOperation(ctx context.Context) error {
    ticker := time.NewTicker(time.Second)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            // Do periodic work
            if done := doWork(); done {
                return nil
            }
        }
    }
}
```

---

#### Q19: What is the difference between context.WithTimeout and context.WithDeadline?

**Answer:**

```go
// WithTimeout - relative time from now
ctx, cancel := context.WithTimeout(parent, 30*time.Second)
// Cancels after 30 seconds from now

// WithDeadline - absolute time
deadline := time.Now().Add(30 * time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)
// Cancels at the specific deadline

// They're functionally equivalent:
// WithTimeout(parent, d) == WithDeadline(parent, time.Now().Add(d))
```

**Use WithTimeout** for most cases (clearer intent)
**Use WithDeadline** when you have a specific deadline

---

### 2.3 Synchronization Primitives

#### Q20: What is the difference between sync.Mutex and sync.RWMutex?

**Answer:**

| `sync.Mutex` | `sync.RWMutex` |
|--------------|----------------|
| Exclusive lock only | Read lock + Write lock |
| One goroutine at a time | Multiple readers OR one writer |
| Simpler, slightly faster | Better for read-heavy workloads |

```go
// Mutex - exclusive access
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}

// RWMutex - multiple readers
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()           // Multiple readers allowed
    defer c.mu.RUnlock()
    return c.data[key]
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()            // Exclusive write
    defer c.mu.Unlock()
    c.data[key] = value
}
```

---

#### Q21: What is sync.Once and when do you use it?

**Answer:**
`sync.Once` ensures a function is executed exactly once, regardless of how many goroutines call it.

```go
var (
    instance *Database
    once     sync.Once
)

func GetDatabase() *Database {
    once.Do(func() {
        // This runs exactly once
        instance = connectToDatabase()
    })
    return instance
}

// Safe to call from multiple goroutines
go GetDatabase()
go GetDatabase()
go GetDatabase()
// connectToDatabase() is called only once
```

**Use cases:**
- Singleton initialization
- One-time configuration
- Lazy initialization

---

#### Q22: What is sync.Pool and when should you use it?

**Answer:**
`sync.Pool` is a cache of temporary objects that can be reused to reduce allocations.

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

func processRequest() {
    // Get buffer from pool (or create new)
    buf := bufferPool.Get().([]byte)
    
    // Use buffer...
    
    // Return to pool for reuse
    bufferPool.Put(buf)
}
```

**When to use:**
- High-frequency allocations of same type
- Objects are expensive to create
- Short-lived objects

**Caveats:**
- Pool may be cleared on GC
- Don't rely on pool for caching

---

#### Q23: Explain atomic operations.

**Answer:**
Atomic operations are thread-safe without locks for simple operations.

```go
import "sync/atomic"

var counter int64

func increment() {
    atomic.AddInt64(&counter, 1)
}

func get() int64 {
    return atomic.LoadInt64(&counter)
}

func set(val int64) {
    atomic.StoreInt64(&counter, val)
}

func compareAndSwap(old, new int64) bool {
    return atomic.CompareAndSwapInt64(&counter, old, new)
}
```

**Available operations:**
- `Add` - Add to value
- `Load` - Read value
- `Store` - Write value
- `Swap` - Exchange value
- `CompareAndSwap` - CAS operation

**Use atomics for:**
- Simple counters
- Flags/booleans
- When mutex overhead is too high

---

## 3. Senior Level (5-8 Years)

### 3.1 Go Scheduler Internals

#### Q24: Explain the Go scheduler's M:N model.

**Answer:**

**Components:**
- **G** (Goroutine) - Lightweight user-space thread
- **M** (Machine) - OS thread
- **P** (Processor) - Scheduling context

```
Goroutines (G):    [G1] [G2] [G3] [G4] [G5] [G6] [G7] [G8]
                      \   |   /       \   |   /
Processors (P):       [  P1  ]       [  P2  ]
                          |              |
OS Threads (M):        [ M1 ]         [ M2 ]
                          |              |
CPU Cores:             [Core1]       [Core2]
```

**How it works:**
1. Each **P** has a local run queue (LRQ) of goroutines
2. **M** must acquire a **P** to run goroutines
3. Global run queue (GRQ) exists for overflow
4. Work stealing: idle P steals from busy P's LRQ

---

#### Q25: What is work stealing in the Go scheduler?

**Answer:**
When a processor's local run queue is empty, it steals goroutines from other processors.

```
Before stealing:
P1: [G1] [G2] [G3] [G4]    P2: [ empty ]
         ^work              ^idle

After stealing:
P1: [G1] [G2]              P2: [G3] [G4]
                                ^stolen
```

**Stealing order:**
1. Check local run queue
2. Check global run queue
3. Check network poller
4. Steal from other P's local queue (steal half)

---

#### Q26: When does the Go scheduler preempt goroutines?

**Answer:**

**Cooperative preemption points:**
- Function calls
- Channel operations
- System calls
- `runtime.Gosched()` explicit yield

**Asynchronous preemption (Go 1.14+):**
- Uses signals (SIGURG on Unix)
- Can preempt even tight loops without function calls

```go
// Pre-Go 1.14: Could starve other goroutines
for {
    compute() // If no function calls, never yields
}

// Go 1.14+: Runtime sends signal to preempt
for {
    x++ // Even this can be preempted
}
```

---

#### Q27: What is GOMAXPROCS and how does it affect performance?

**Answer:**

```go
runtime.GOMAXPROCS(n) // Set max number of Ps
n := runtime.GOMAXPROCS(0) // Query current value (0 = don't change)
```

**Default:** Number of CPU cores

**Impact:**
- Controls true parallelism (goroutines running simultaneously)
- More GOMAXPROCS ≠ always better
- Diminishing returns beyond CPU cores

**Tuning guidelines:**
| Workload | GOMAXPROCS |
|----------|------------|
| CPU-bound | = CPU cores (default) |
| I/O-bound | Can exceed cores |
| Mixed | Profile and tune |

```go
// Check current value
fmt.Println(runtime.GOMAXPROCS(0))  // e.g., 8
fmt.Println(runtime.NumCPU())       // e.g., 8
```

---

### 3.2 Advanced Channel Patterns

#### Q28: How do you implement a semaphore with channels?

**Answer:**

```go
type Semaphore struct {
    sem chan struct{}
}

func NewSemaphore(max int) *Semaphore {
    return &Semaphore{
        sem: make(chan struct{}, max),
    }
}

func (s *Semaphore) Acquire() {
    s.sem <- struct{}{}
}

func (s *Semaphore) TryAcquire() bool {
    select {
    case s.sem <- struct{}{}:
        return true
    default:
        return false
    }
}

func (s *Semaphore) Release() {
    <-s.sem
}

// Usage: Limit concurrent operations
sem := NewSemaphore(10) // Max 10 concurrent

for _, task := range tasks {
    sem.Acquire()
    go func(t Task) {
        defer sem.Release()
        process(t)
    }(task)
}
```

---

#### Q29: How do you implement a rate limiter with channels?

**Answer:**

```go
type RateLimiter struct {
    tokens   chan struct{}
    interval time.Duration
    stop     chan struct{}
}

func NewRateLimiter(rate int, per time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens:   make(chan struct{}, rate),
        interval: per / time.Duration(rate),
        stop:     make(chan struct{}),
    }
    
    // Fill initial bucket
    for i := 0; i < rate; i++ {
        rl.tokens <- struct{}{}
    }
    
    // Refill tokens
    go rl.refill()
    
    return rl
}

func (rl *RateLimiter) refill() {
    ticker := time.NewTicker(rl.interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-rl.stop:
            return
        case <-ticker.C:
            select {
            case rl.tokens <- struct{}{}:
            default: // Bucket full
            }
        }
    }
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
```

---

#### Q30: How do you implement a pub/sub system with channels?

**Answer:**

```go
type PubSub struct {
    mu     sync.RWMutex
    subs   map[string][]chan interface{}
    closed bool
}

func NewPubSub() *PubSub {
    return &PubSub{
        subs: make(map[string][]chan interface{}),
    }
}

func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    ch := make(chan interface{}, 10)
    ps.subs[topic] = append(ps.subs[topic], ch)
    return ch
}

func (ps *PubSub) Publish(topic string, msg interface{}) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()
    
    for _, ch := range ps.subs[topic] {
        select {
        case ch <- msg:
        default:
            // Subscriber too slow, drop message
        }
    }
}

func (ps *PubSub) Unsubscribe(topic string, ch <-chan interface{}) {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    
    subs := ps.subs[topic]
    for i, sub := range subs {
        if sub == ch {
            ps.subs[topic] = append(subs[:i], subs[i+1:]...)
            close(sub)
            return
        }
    }
}
```

---

#### Q31: Explain the nil channel trick in select.

**Answer:**
A nil channel blocks forever in select. Use this to dynamically enable/disable cases.

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

**Use cases:**
- Merge channels that may close at different times
- Dynamically enable/disable select cases
- Implement priority channels

---

### 3.3 Debugging Concurrency Issues

#### Q32: How do you debug goroutine leaks?

**Answer:**

**Detection:**
```go
// Monitor goroutine count
func monitorGoroutines() {
    ticker := time.NewTicker(10 * time.Second)
    for range ticker.C {
        fmt.Printf("Goroutines: %d\n", runtime.NumGoroutine())
    }
}

// Or expose via HTTP
import _ "net/http/pprof"
// curl http://localhost:6060/debug/pprof/goroutine?debug=1
```

**Common causes:**
```go
// 1. Blocked channel send
go func() {
    ch <- value // Blocks forever if no receiver
}()

// 2. Blocked channel receive
go func() {
    <-ch // Blocks forever if channel never receives/closes
}()

// 3. Missing cancellation
go func() {
    for {
        doWork() // No way to exit
    }
}()
```

**Fixes:**
```go
// Use context for cancellation
func worker(ctx context.Context, ch <-chan int) {
    for {
        select {
        case <-ctx.Done():
            return // Clean exit
        case v, ok := <-ch:
            if !ok {
                return // Channel closed
            }
            process(v)
        }
    }
}

// Use buffered channel or select with default
go func() {
    select {
    case ch <- value:
    case <-time.After(time.Second):
        // Timeout, give up
    }
}()
```

---

#### Q33: How do you use the race detector effectively?

**Answer:**

```bash
# Enable race detector
go test -race ./...
go build -race
go run -race main.go
```

**Example race:**
```go
var count int

func main() {
    go func() { count++ }()
    go func() { count++ }()
    time.Sleep(time.Second)
    fmt.Println(count)
}
```

**Race detector output:**
```
WARNING: DATA RACE
Write at 0x... by goroutine 7:
  main.main.func1()
      main.go:8 +0x3a

Previous write at 0x... by goroutine 6:
  main.main.func2()
      main.go:9 +0x3a
```

**Best practices:**
- Run race detector in CI/CD
- Test with `-race` regularly
- Race detector has 2-20x slowdown
- Not suitable for production

---

#### Q34: What causes deadlocks and how do you prevent them?

**Answer:**

**Common deadlock patterns:**

**1. Self-deadlock:**
```go
ch := make(chan int)
ch <- 1 // Deadlock! No one to receive
```

**2. Circular wait:**
```go
ch1, ch2 := make(chan int), make(chan int)

go func() {
    <-ch1
    ch2 <- 1
}()

go func() {
    <-ch2
    ch1 <- 1
}()
// Both waiting for each other
```

**3. Mutex deadlock:**
```go
var mu1, mu2 sync.Mutex

// Goroutine 1
go func() {
    mu1.Lock()
    mu2.Lock() // Waits for mu2
    mu2.Unlock()
    mu1.Unlock()
}()

// Goroutine 2
go func() {
    mu2.Lock()
    mu1.Lock() // Waits for mu1 - DEADLOCK
    mu1.Unlock()
    mu2.Unlock()
}()
```

**Prevention:**
```go
// 1. Use buffered channels
ch := make(chan int, 1)
ch <- 1 // Doesn't block

// 2. Use select with timeout
select {
case ch <- 1:
case <-time.After(time.Second):
    log.Println("timeout")
}

// 3. Always lock mutexes in same order
// 4. Use defer for unlocking
// 5. Use context with timeout
```

---

### 3.4 Performance Optimization

#### Q35: How do you reduce lock contention?

**Answer:**

**1. Reduce critical section size:**
```go
// BAD: Lock held too long
func (c *Cache) Get(key string) string {
    c.mu.Lock()
    defer c.mu.Unlock()
    val := c.data[key]
    processValue(val) // Don't do work under lock!
    return val
}

// GOOD: Minimal critical section
func (c *Cache) Get(key string) string {
    c.mu.Lock()
    val := c.data[key]
    c.mu.Unlock()
    
    processValue(val)
    return val
}
```

**2. Use RWMutex for read-heavy workloads:**
```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) string {
    c.mu.RLock()         // Multiple readers
    defer c.mu.RUnlock()
    return c.data[key]
}
```

**3. Shard data:**
```go
const numShards = 32

type ShardedMap struct {
    shards [numShards]struct {
        mu   sync.RWMutex
        data map[string]string
    }
}

func (m *ShardedMap) getShard(key string) int {
    return int(fnv32(key) % numShards)
}

func (m *ShardedMap) Get(key string) string {
    shard := &m.shards[m.getShard(key)]
    shard.mu.RLock()
    defer shard.mu.RUnlock()
    return shard.data[key]
}
```

**4. Use atomic operations:**
```go
// Instead of mutex for simple counter
var counter int64
atomic.AddInt64(&counter, 1)
```

**5. Use sync.Pool:**
```go
var pool = sync.Pool{
    New: func() interface{} { return new(Buffer) },
}
```

---

#### Q36: When should you use channels vs mutexes?

**Answer:**

| Use Channels | Use Mutex |
|--------------|-----------|
| Passing ownership of data | Protecting shared state |
| Coordinating goroutines | Simple counter/flag |
| Implementing patterns (worker pool) | Cache with concurrent access |
| Signaling events | When ownership isn't transferred |

**Rule of thumb:**
- Share memory by communicating (channels)
- Don't communicate by sharing memory (mutex)

```go
// Channel: Transfer ownership
func producer(out chan<- *Data) {
    data := &Data{}
    // ... populate data
    out <- data // Transfer ownership
    // Don't use data after sending!
}

// Mutex: Protect shared access
type Counter struct {
    mu  sync.Mutex
    val int
}
func (c *Counter) Inc() {
    c.mu.Lock()
    c.val++
    c.mu.Unlock()
}
```

---

## 4. Staff/Principal Level (8+ Years)

### 4.1 Advanced Scheduler Topics

#### Q37: Explain goroutine stack growth and shrinking.

**Answer:**

**Initial stack:** 2KB (tiny compared to OS thread's 1-8MB)

**Growth mechanism:**
1. Compiler inserts stack checks at function entry
2. If stack too small, runtime allocates larger stack (2x)
3. Copies entire stack to new location
4. Updates all pointers

```go
// Stack check at function entry (simplified)
func someFunc() {
    if sp < stackGuard {
        runtime.morestack()
    }
    // ... function body
}
```

**Stack shrinking:**
- Happens during GC
- If stack is <25% used, shrinks by half
- Prevents memory waste from goroutines that had deep call stacks

**Implications:**
- Pointers to stack variables can change
- Can't pass pointer to stack variable to C code (use `runtime.LockOSThread()`)
- Stack trace shows accurate call chain

---

#### Q38: Explain the network poller in Go.

**Answer:**

**Traditional blocking I/O:**
```
Goroutine → Syscall → OS Thread blocks → Wastes resources
```

**Go's netpoller:**
```
Goroutine → Non-blocking syscall → Register with epoll/kqueue
            → Goroutine parked → OS Thread freed
            → I/O ready → Goroutine woken → Continues
```

**How it works:**
1. Network operations use non-blocking sockets
2. Goroutine is parked (removed from run queue)
3. OS thread continues running other goroutines
4. When I/O ready (via epoll/kqueue/IOCP), goroutine re-queued

```go
// This doesn't block the OS thread:
conn, err := net.Dial("tcp", "example.com:80")
```

**Benefits:**
- Millions of concurrent connections
- Efficient use of limited OS threads
- No callback hell (looks synchronous)

---

#### Q39: How do you implement a lock-free data structure?

**Answer:**

**Lock-free stack using CAS:**
```go
type Stack struct {
    head atomic.Pointer[node]
}

type node struct {
    value interface{}
    next  *node
}

func (s *Stack) Push(v interface{}) {
    newNode := &node{value: v}
    for {
        oldHead := s.head.Load()
        newNode.next = oldHead
        if s.head.CompareAndSwap(oldHead, newNode) {
            return
        }
        // CAS failed, retry
    }
}

func (s *Stack) Pop() (interface{}, bool) {
    for {
        oldHead := s.head.Load()
        if oldHead == nil {
            return nil, false
        }
        newHead := oldHead.next
        if s.head.CompareAndSwap(oldHead, newHead) {
            return oldHead.value, true
        }
        // CAS failed, retry
    }
}
```

**When to use lock-free:**
- Extremely high contention
- Real-time requirements
- When lock overhead is bottleneck

**Challenges:**
- Complex to implement correctly
- ABA problem
- Memory ordering issues
- Often not worth it vs well-designed locks

---

### 4.2 System Design with Concurrency

#### Q40: How do you design a high-throughput concurrent system?

**Answer:**

**Principles:**

**1. Minimize shared state:**
```go
// BAD: Global shared state
var globalCache map[string]*Item

// GOOD: Pass data explicitly
func process(ctx context.Context, cache *Cache, item *Item) error
```

**2. Use bounded concurrency:**
```go
type WorkerPool struct {
    workers  int
    jobQueue chan Job
    results  chan Result
    ctx      context.Context
}

func (p *WorkerPool) Submit(job Job) error {
    select {
    case p.jobQueue <- job:
        return nil
    case <-p.ctx.Done():
        return p.ctx.Err()
    default:
        return ErrQueueFull // Back-pressure
    }
}
```

**3. Design for cancellation:**
```go
func longOperation(ctx context.Context) error {
    resultCh := make(chan result, 1)
    errCh := make(chan error, 1)
    
    go func() {
        r, err := doExpensiveWork()
        if err != nil {
            errCh <- err
            return
        }
        resultCh <- r
    }()
    
    select {
    case r := <-resultCh:
        return processResult(r)
    case err := <-errCh:
        return err
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

**4. Use proper back-pressure:**
```go
// Rate limit incoming requests
limiter := rate.NewLimiter(rate.Limit(1000), 100) // 1000 RPS, burst 100

func handleRequest(w http.ResponseWriter, r *http.Request) {
    if !limiter.Allow() {
        http.Error(w, "rate limited", http.StatusTooManyRequests)
        return
    }
    // Process request
}
```

---

#### Q41: How do you implement graceful shutdown with concurrent workers?

**Answer:**

```go
type Server struct {
    httpServer *http.Server
    workers    *WorkerPool
    wg         sync.WaitGroup
    shutdown   chan struct{}
}

func (s *Server) Start() {
    // Start HTTP server
    go func() {
        if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatal(err)
        }
    }()
    
    // Start workers
    s.workers.Start()
}

func (s *Server) Shutdown(ctx context.Context) error {
    // Signal shutdown
    close(s.shutdown)
    
    // Stop accepting new requests
    if err := s.httpServer.Shutdown(ctx); err != nil {
        return err
    }
    
    // Wait for workers to finish current jobs
    done := make(chan struct{})
    go func() {
        s.workers.Stop()
        s.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// Main function
func main() {
    server := NewServer()
    server.Start()
    
    // Wait for interrupt
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    // Graceful shutdown with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("forced shutdown:", err)
    }
    
    log.Println("server stopped gracefully")
}
```

---

#### Q42: How do you handle panic recovery in goroutines?

**Answer:**

**Individual goroutine recovery:**
```go
func safeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("panic recovered: %v\n%s", r, debug.Stack())
            }
        }()
        fn()
    }()
}
```

**Worker pool with recovery:**
```go
func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()
    
    for job := range p.jobs {
        func() {
            defer func() {
                if r := recover(); r != nil {
                    log.Printf("worker %d panic: %v", id, r)
                    p.errors <- fmt.Errorf("worker panic: %v", r)
                }
            }()
            
            result := process(job)
            p.results <- result
        }()
    }
}
```

**HTTP handler recovery middleware:**
```go
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("panic: %v\n%s", err, debug.Stack())
                http.Error(w, "Internal Server Error", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

---

### 4.3 Advanced Patterns

#### Q43: Implement a pipeline with cancellation.

**Answer:**

```go
func pipeline(ctx context.Context, input <-chan int) <-chan int {
    stage1 := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                select {
                case <-ctx.Done():
                    return
                case out <- n * 2:
                }
            }
        }()
        return out
    }
    
    stage2 := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                select {
                case <-ctx.Done():
                    return
                case out <- n + 1:
                }
            }
        }()
        return out
    }
    
    stage3 := func(in <-chan int) <-chan int {
        out := make(chan int)
        go func() {
            defer close(out)
            for n := range in {
                select {
                case <-ctx.Done():
                    return
                case out <- n * n:
                }
            }
        }()
        return out
    }
    
    return stage3(stage2(stage1(input)))
}

// Usage
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    input := make(chan int)
    go func() {
        defer close(input)
        for i := 0; i < 100; i++ {
            select {
            case <-ctx.Done():
                return
            case input <- i:
            }
        }
    }()
    
    for result := range pipeline(ctx, input) {
        fmt.Println(result)
    }
}
```

---

#### Q44: Implement singleflight pattern.

**Answer:**

`singleflight` ensures only one execution for duplicate requests.

```go
import "golang.org/x/sync/singleflight"

type Cache struct {
    sf    singleflight.Group
    data  map[string]string
    mu    sync.RWMutex
}

func (c *Cache) Get(key string) (string, error) {
    // Check cache first
    c.mu.RLock()
    if val, ok := c.data[key]; ok {
        c.mu.RUnlock()
        return val, nil
    }
    c.mu.RUnlock()
    
    // Singleflight: only one goroutine fetches
    result, err, _ := c.sf.Do(key, func() (interface{}, error) {
        // This runs only once even if 1000 goroutines call Get(key)
        val, err := fetchFromDatabase(key)
        if err != nil {
            return "", err
        }
        
        // Update cache
        c.mu.Lock()
        c.data[key] = val
        c.mu.Unlock()
        
        return val, nil
    })
    
    if err != nil {
        return "", err
    }
    return result.(string), nil
}
```

**Use cases:**
- Prevent cache stampede
- Deduplicate expensive operations
- Rate limit backend calls

---

#### Q45: Implement errgroup pattern.

**Answer:**

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([]Response, error) {
    g, ctx := errgroup.WithContext(ctx)
    responses := make([]Response, len(urls))
    
    for i, url := range urls {
        i, url := i, url // Capture loop variables
        g.Go(func() error {
            resp, err := fetch(ctx, url)
            if err != nil {
                return err
            }
            responses[i] = resp
            return nil
        })
    }
    
    if err := g.Wait(); err != nil {
        return nil, err // First error cancels all
    }
    
    return responses, nil
}
```

**Benefits:**
- Waits for all goroutines
- Returns first error
- Cancels remaining on error
- Cleaner than manual WaitGroup + error handling

---

## 5. Concurrency Coding Challenges

### Challenge 1: Implement a Bounded Worker Pool

```go
// Requirements:
// - Fixed number of workers
// - Bounded job queue (back-pressure)
// - Graceful shutdown
// - Error collection

type Job func() error

type WorkerPool struct {
    numWorkers int
    jobQueue   chan Job
    errors     chan error
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

func NewWorkerPool(numWorkers, queueSize int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    return &WorkerPool{
        numWorkers: numWorkers,
        jobQueue:   make(chan Job, queueSize),
        errors:     make(chan error, queueSize),
        ctx:        ctx,
        cancel:     cancel,
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.numWorkers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
}

func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()
    
    for {
        select {
        case <-p.ctx.Done():
            return
        case job, ok := <-p.jobQueue:
            if !ok {
                return
            }
            if err := p.safeExecute(job); err != nil {
                select {
                case p.errors <- err:
                default:
                    log.Printf("error dropped: %v", err)
                }
            }
        }
    }
}

func (p *WorkerPool) safeExecute(job Job) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    return job()
}

func (p *WorkerPool) Submit(job Job) error {
    select {
    case p.jobQueue <- job:
        return nil
    case <-p.ctx.Done():
        return errors.New("pool is shutting down")
    default:
        return errors.New("job queue is full")
    }
}

func (p *WorkerPool) Shutdown(timeout time.Duration) error {
    p.cancel()
    close(p.jobQueue)
    
    done := make(chan struct{})
    go func() {
        p.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        close(p.errors)
        return nil
    case <-time.After(timeout):
        return errors.New("shutdown timeout")
    }
}

func (p *WorkerPool) Errors() <-chan error {
    return p.errors
}
```

---

### Challenge 2: Implement a Concurrent-Safe LRU Cache

```go
type LRUCache struct {
    capacity int
    cache    map[string]*list.Element
    list     *list.List
    mu       sync.RWMutex
}

type entry struct {
    key   string
    value interface{}
}

func NewLRUCache(capacity int) *LRUCache {
    return &LRUCache{
        capacity: capacity,
        cache:    make(map[string]*list.Element),
        list:     list.New(),
    }
}

func (c *LRUCache) Get(key string) (interface{}, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.cache[key]; ok {
        c.list.MoveToFront(elem)
        return elem.Value.(*entry).value, true
    }
    return nil, false
}

func (c *LRUCache) Put(key string, value interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.cache[key]; ok {
        c.list.MoveToFront(elem)
        elem.Value.(*entry).value = value
        return
    }
    
    if c.list.Len() >= c.capacity {
        oldest := c.list.Back()
        if oldest != nil {
            c.list.Remove(oldest)
            delete(c.cache, oldest.Value.(*entry).key)
        }
    }
    
    elem := c.list.PushFront(&entry{key: key, value: value})
    c.cache[key] = elem
}

func (c *LRUCache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    if elem, ok := c.cache[key]; ok {
        c.list.Remove(elem)
        delete(c.cache, key)
    }
}
```

---

### Challenge 3: Implement a Concurrent Barrier

```go
type Barrier struct {
    n        int
    count    int
    epoch    int
    mu       sync.Mutex
    cond     *sync.Cond
}

func NewBarrier(n int) *Barrier {
    b := &Barrier{n: n}
    b.cond = sync.NewCond(&b.mu)
    return b
}

func (b *Barrier) Wait() {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    epoch := b.epoch
    b.count++
    
    if b.count == b.n {
        b.epoch++
        b.count = 0
        b.cond.Broadcast()
        return
    }
    
    for epoch == b.epoch {
        b.cond.Wait()
    }
}

// Usage
func main() {
    barrier := NewBarrier(3)
    
    for i := 0; i < 3; i++ {
        go func(id int) {
            fmt.Printf("Worker %d starting\n", id)
            time.Sleep(time.Duration(id) * time.Second)
            
            fmt.Printf("Worker %d waiting at barrier\n", id)
            barrier.Wait()
            
            fmt.Printf("Worker %d passed barrier\n", id)
        }(i)
    }
    
    time.Sleep(5 * time.Second)
}
```

---

### Challenge 4: Implement a Ring Buffer

```go
type RingBuffer struct {
    data     []interface{}
    size     int
    head     int
    tail     int
    count    int
    mu       sync.Mutex
    notEmpty *sync.Cond
    notFull  *sync.Cond
}

func NewRingBuffer(size int) *RingBuffer {
    rb := &RingBuffer{
        data: make([]interface{}, size),
        size: size,
    }
    rb.notEmpty = sync.NewCond(&rb.mu)
    rb.notFull = sync.NewCond(&rb.mu)
    return rb
}

func (rb *RingBuffer) Put(item interface{}) {
    rb.mu.Lock()
    defer rb.mu.Unlock()
    
    for rb.count == rb.size {
        rb.notFull.Wait()
    }
    
    rb.data[rb.tail] = item
    rb.tail = (rb.tail + 1) % rb.size
    rb.count++
    
    rb.notEmpty.Signal()
}

func (rb *RingBuffer) Get() interface{} {
    rb.mu.Lock()
    defer rb.mu.Unlock()
    
    for rb.count == 0 {
        rb.notEmpty.Wait()
    }
    
    item := rb.data[rb.head]
    rb.head = (rb.head + 1) % rb.size
    rb.count--
    
    rb.notFull.Signal()
    return item
}

func (rb *RingBuffer) TryPut(item interface{}) bool {
    rb.mu.Lock()
    defer rb.mu.Unlock()
    
    if rb.count == rb.size {
        return false
    }
    
    rb.data[rb.tail] = item
    rb.tail = (rb.tail + 1) % rb.size
    rb.count++
    
    rb.notEmpty.Signal()
    return true
}

func (rb *RingBuffer) TryGet() (interface{}, bool) {
    rb.mu.Lock()
    defer rb.mu.Unlock()
    
    if rb.count == 0 {
        return nil, false
    }
    
    item := rb.data[rb.head]
    rb.head = (rb.head + 1) % rb.size
    rb.count--
    
    rb.notFull.Signal()
    return item, true
}
```

---

## 6. Quick Reference & Cheat Sheet

### Channel Operations Quick Reference

| Operation | nil channel | closed channel | open channel |
|-----------|-------------|----------------|--------------|
| `ch <- v` | blocks forever | **panic** | sends or blocks |
| `<-ch` | blocks forever | returns zero, false | receives or blocks |
| `close(ch)` | **panic** | **panic** | closes |
| `len(ch)` | 0 | 0 or remaining items | buffered items |
| `cap(ch)` | 0 | capacity | capacity |

### Synchronization Primitives

| Primitive | Use Case | Key Methods |
|-----------|----------|-------------|
| `sync.Mutex` | Exclusive access | `Lock()`, `Unlock()` |
| `sync.RWMutex` | Read-heavy access | `RLock()`, `RUnlock()`, `Lock()`, `Unlock()` |
| `sync.WaitGroup` | Wait for goroutines | `Add()`, `Done()`, `Wait()` |
| `sync.Once` | One-time init | `Do(func())` |
| `sync.Pool` | Object reuse | `Get()`, `Put()` |
| `sync.Cond` | Wait for condition | `Wait()`, `Signal()`, `Broadcast()` |
| `sync.Map` | Concurrent map | `Load()`, `Store()`, `Delete()`, `Range()` |

### Context Methods

```go
context.Background()                    // Root context
context.TODO()                          // Placeholder
context.WithCancel(parent)              // Cancellable
context.WithTimeout(parent, duration)   // Timeout
context.WithDeadline(parent, time)      // Deadline
context.WithValue(parent, key, val)     // With value
ctx.Done()                              // Cancellation channel
ctx.Err()                               // context.Canceled or DeadlineExceeded
ctx.Value(key)                          // Get value
```

### Common Patterns

**Done channel:**
```go
done := make(chan struct{})
go func() {
    defer close(done)
    // work
}()
<-done
```

**Timeout:**
```go
select {
case res := <-ch:
    use(res)
case <-time.After(timeout):
    handle timeout
}
```

**Non-blocking:**
```go
select {
case msg := <-ch:
    process(msg)
default:
    // channel empty
}
```

**First wins:**
```go
result := make(chan int, n)
for _, s := range servers {
    go func(s string) { result <- query(s) }(s)
}
return <-result
```

### Goroutine Lifecycle Checklist

- [ ] Has clear exit condition
- [ ] Respects context cancellation
- [ ] Handles channel closure
- [ ] Has panic recovery (if needed)
- [ ] Resources cleaned up on exit
- [ ] Documented ownership of channels

### Race Condition Checklist

- [ ] All shared data protected by mutex OR channels
- [ ] Mutex locked before access, unlocked after
- [ ] No holding locks while calling external code
- [ ] No copying sync types (Mutex, WaitGroup, etc.)
- [ ] Run with `-race` flag regularly

---

## Summary by Experience Level

| Level | Focus Areas |
|-------|-------------|
| **Junior** | Goroutines, channels, WaitGroup, basic mutex |
| **Mid** | Context, patterns (fan-in/out, worker pool), RWMutex, channel idioms |
| **Senior** | Scheduler internals, debugging, performance, advanced patterns |
| **Staff+** | Lock-free, system design, architecture, distributed patterns |

**Total Questions:** 45 questions across all levels + 4 coding challenges

---

> **Good luck with your concurrency interviews!** 🚀
