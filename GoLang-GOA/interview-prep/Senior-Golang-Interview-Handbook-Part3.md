# Senior Golang Interview Question Handbook - Part 3

> **Target Audience:** Experienced Golang Developers (5–12 years)  
> **Purpose:** Prepare for senior backend/product company interviews  
> **Difficulty Legend:** 🟢 Basic | 🟡 Intermediate | 🔴 Advanced | ⚫ Expert

---

## Table of Contents - Part 3

10. [Performance & Profiling](#10-performance--profiling)
11. [Distributed Systems & Microservices](#11-distributed-systems--microservices)
12. [System Design Questions (CRITICAL)](#12-system-design-questions-critical)
13. [Rapid Revision Section](#13-rapid-revision-section)

---

## 10. Performance & Profiling

### 10.1 pprof Fundamentals

#### 🟡 Q74: How do you use pprof for profiling Go applications?

**Answer:**

**Setup for HTTP server:**
```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe(":6060", nil))
    }()
    // Application code...
}
```

**Available profiles:**
```bash
# CPU profile (30 seconds default)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap profile (current memory)
go tool pprof http://localhost:6060/debug/pprof/heap

# Goroutine profile
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Block profile (blocking operations)
go tool pprof http://localhost:6060/debug/pprof/block

# Mutex profile (lock contention)
go tool pprof http://localhost:6060/debug/pprof/mutex

# Threadcreate profile
go tool pprof http://localhost:6060/debug/pprof/threadcreate

# Trace (execution tracer)
curl -o trace.out http://localhost:6060/debug/pprof/trace?seconds=5
go tool trace trace.out
```

**pprof interactive commands:**
```bash
(pprof) top10              # Top 10 functions by CPU/memory
(pprof) top -cum           # Top by cumulative time
(pprof) list funcName      # Source code view
(pprof) web                # Generate graph (requires graphviz)
(pprof) svg > profile.svg  # Export as SVG
(pprof) peek funcName      # Show callers and callees
```

---

#### 🔴 Q75: What's the difference between CPU and memory profiling?

**Answer:**

| Aspect | CPU Profile | Memory Profile |
|--------|-------------|----------------|
| **Measures** | Time spent in functions | Memory allocations |
| **How** | Samples stack traces (100Hz) | Records allocations |
| **Use for** | Finding slow functions | Finding allocation-heavy code |
| **Key metric** | % of total time | Bytes/objects allocated |

**CPU profiling:**
```go
import "runtime/pprof"

func main() {
    f, _ := os.Create("cpu.pprof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Code to profile
}
```

**Memory profiling:**
```go
func main() {
    // Run your code
    doWork()
    
    // Write heap profile
    f, _ := os.Create("mem.pprof")
    runtime.GC() // Get accurate picture
    pprof.WriteHeapProfile(f)
    f.Close()
}
```

**Heap profile types:**
```bash
# inuse_space: Currently allocated bytes (default)
go tool pprof -inuse_space http://localhost:6060/debug/pprof/heap

# inuse_objects: Currently allocated objects
go tool pprof -inuse_objects http://localhost:6060/debug/pprof/heap

# alloc_space: Total bytes allocated (cumulative)
go tool pprof -alloc_space http://localhost:6060/debug/pprof/heap

# alloc_objects: Total objects allocated (cumulative)
go tool pprof -alloc_objects http://localhost:6060/debug/pprof/heap
```

---

### 10.2 Benchmarking Techniques

#### 🔴 Q76: How do you write effective benchmarks?

**Answer:**

**Basic benchmark:**
```go
func BenchmarkStringConcat(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _ = "Hello" + " " + "World"
    }
}
```

**Benchmark with setup (excluded from timing):**
```go
func BenchmarkSort(b *testing.B) {
    data := make([]int, 10000)
    for i := range data {
        data[i] = rand.Intn(10000)
    }
    
    b.ResetTimer() // Exclude setup from timing
    
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        testData := make([]int, len(data))
        copy(testData, data)
        b.StartTimer()
        
        sort.Ints(testData)
    }
}
```

**Compare implementations:**
```go
func BenchmarkConcat(b *testing.B) {
    strs := []string{"Hello", " ", "World", "!", " ", "Go"}
    
    b.Run("Plus", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            result := ""
            for _, s := range strs {
                result += s
            }
            _ = result
        }
    })
    
    b.Run("Builder", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            var builder strings.Builder
            for _, s := range strs {
                builder.WriteString(s)
            }
            _ = builder.String()
        }
    })
    
    b.Run("Join", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            _ = strings.Join(strs, "")
        }
    })
}
```

**Memory allocation benchmarks:**
```go
func BenchmarkWithAllocs(b *testing.B) {
    b.ReportAllocs() // Report allocations
    
    for i := 0; i < b.N; i++ {
        m := make(map[string]int)
        m["key"] = 42
        _ = m
    }
}
```

**Running benchmarks:**
```bash
go test -bench=.                        # Run all benchmarks
go test -bench=BenchmarkConcat         # Run specific benchmark
go test -bench=. -benchmem             # Include memory stats
go test -bench=. -benchtime=5s         # Run for 5 seconds
go test -bench=. -count=10             # Run 10 times
go test -bench=. -cpuprofile=cpu.out   # Generate CPU profile
go test -bench=. -memprofile=mem.out   # Generate memory profile
```

**Output interpretation:**
```
BenchmarkConcat/Plus-8       1000000    1234 ns/op    256 B/op    8 allocs/op
BenchmarkConcat/Builder-8    5000000     234 ns/op     64 B/op    2 allocs/op
BenchmarkConcat/Join-8       5000000     189 ns/op     48 B/op    1 allocs/op
```
- `1000000` = iterations run
- `1234 ns/op` = nanoseconds per operation
- `256 B/op` = bytes allocated per operation
- `8 allocs/op` = allocations per operation

---

### 10.3 Optimization Techniques

#### 🔴 Q77: What are common Go optimization techniques?

**Answer:**

**1. Reduce allocations:**
```go
// BAD: Allocates on every iteration
func processBad(items []Item) {
    for _, item := range items {
        result := make([]byte, 1024) // Allocation
        process(item, result)
    }
}

// GOOD: Reuse buffer
func processGood(items []Item) {
    result := make([]byte, 1024) // Single allocation
    for _, item := range items {
        result = result[:0] // Reset slice
        process(item, result)
    }
}
```

**2. Use sync.Pool:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 4096)
    },
}

func processWithPool(data []byte) {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf)
    
    // Use buf...
}
```

**3. Pre-allocate slices and maps:**
```go
// BAD: Multiple reallocations
var result []int
for i := 0; i < n; i++ {
    result = append(result, i)
}

// GOOD: Single allocation
result := make([]int, 0, n)
for i := 0; i < n; i++ {
    result = append(result, i)
}

// GOOD: Maps too
m := make(map[string]int, expectedSize)
```

**4. Avoid interface{} when possible:**
```go
// BAD: Boxing causes allocation
func sum(nums []interface{}) int64 {
    var total int64
    for _, n := range nums {
        total += n.(int64)
    }
    return total
}

// GOOD: Use generics (Go 1.18+)
func sum[T constraints.Integer](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

**5. String building:**
```go
// BAD: O(n²) allocations
func concatBad(strs []string) string {
    result := ""
    for _, s := range strs {
        result += s // New allocation each time
    }
    return result
}

// GOOD: O(n) single allocation
func concatGood(strs []string) string {
    var total int
    for _, s := range strs {
        total += len(s)
    }
    
    var b strings.Builder
    b.Grow(total)
    for _, s := range strs {
        b.WriteString(s)
    }
    return b.String()
}
```

**6. Avoid defer in hot loops:**
```go
// BAD: defer overhead in loop
func processBad(files []string) {
    for _, file := range files {
        f, _ := os.Open(file)
        defer f.Close() // Deferred until function returns
        // Process...
    }
}

// GOOD: Close immediately or use helper
func processGood(files []string) {
    for _, file := range files {
        processFile(file)
    }
}

func processFile(file string) {
    f, _ := os.Open(file)
    defer f.Close() // Deferred correctly per file
    // Process...
}
```

**7. Use appropriate data structures:**
```go
// For membership checks: map > slice
set := make(map[string]struct{})
set["key"] = struct{}{}
if _, ok := set["key"]; ok {
    // Found - O(1)
}

// For ordered iteration: slice > map
// For frequency counts: map[T]int
// For concurrent access: sync.Map (specific cases)
```

---

#### ⚫ Q78: How do you identify and fix memory leaks in Go?

**Answer:**

**Common leak patterns:**

**1. Goroutine leaks:**
```go
// LEAK: Goroutine blocked forever
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // Blocks forever
        fmt.Println(val)
    }()
    // ch never written to, goroutine never exits
}

// FIX: Use context or close channel
func fixed(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done():
            return // Clean exit
        }
    }()
}
```

**2. Unclosed resources:**
```go
// LEAK: Response body not closed
func leak(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    // resp.Body never closed!
    return io.ReadAll(resp.Body)
}

// FIX: Always close response body
func fixed(url string) ([]byte, error) {
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}
```

**3. Subscriber not unsubscribed:**
```go
// LEAK: Channel subscription never cancelled
type EventBus struct {
    subscribers map[string][]chan Event
    mu          sync.RWMutex
}

func (e *EventBus) Subscribe(topic string) <-chan Event {
    ch := make(chan Event, 10)
    e.mu.Lock()
    e.subscribers[topic] = append(e.subscribers[topic], ch)
    e.mu.Unlock()
    return ch
    // No way to unsubscribe!
}

// FIX: Return unsubscribe function
func (e *EventBus) Subscribe(topic string) (<-chan Event, func()) {
    ch := make(chan Event, 10)
    e.mu.Lock()
    e.subscribers[topic] = append(e.subscribers[topic], ch)
    e.mu.Unlock()
    
    unsubscribe := func() {
        e.mu.Lock()
        defer e.mu.Unlock()
        // Remove ch from subscribers[topic]
    }
    return ch, unsubscribe
}
```

**4. Slice capacity retention:**
```go
// LEAK: Holds reference to large backing array
func leak(large []byte) []byte {
    return large[:10] // Still references large array
}

// FIX: Copy to new slice
func fixed(large []byte) []byte {
    result := make([]byte, 10)
    copy(result, large[:10])
    return result
}
```

**Detecting leaks:**
```bash
# Monitor goroutine count
curl http://localhost:6060/debug/pprof/goroutine?debug=1

# Compare heap profiles over time
curl http://localhost:6060/debug/pprof/heap > t1.pprof
# Wait...
curl http://localhost:6060/debug/pprof/heap > t2.pprof
go tool pprof -base=t1.pprof t2.pprof

# Check goroutine growth
watch -n 5 'curl -s http://localhost:6060/debug/pprof/goroutine?debug=1 | head -1'
```

**Runtime metrics:**
```go
import "runtime"

func logMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    log.Printf("Alloc: %d MB", m.Alloc/1024/1024)
    log.Printf("HeapAlloc: %d MB", m.HeapAlloc/1024/1024)
    log.Printf("NumGC: %d", m.NumGC)
    log.Printf("Goroutines: %d", runtime.NumGoroutine())
}
```

---

### 10.4 Execution Tracer

#### 🔴 Q79: How do you use the Go execution tracer?

**Answer:**

**Capture trace:**
```go
import "runtime/trace"

func main() {
    f, _ := os.Create("trace.out")
    trace.Start(f)
    defer trace.Stop()
    
    // Code to trace
}
```

**Or via HTTP:**
```bash
curl -o trace.out 'http://localhost:6060/debug/pprof/trace?seconds=5'
go tool trace trace.out
```

**What trace shows:**
- Goroutine scheduling
- Network/syscall blocking
- GC events
- Heap size over time
- Processor utilization

**Trace views:**
```bash
go tool trace trace.out
# Opens browser with:
# - View trace (timeline)
# - Goroutine analysis
# - Network blocking profile
# - Sync blocking profile
# - Syscall blocking profile
# - Scheduler latency profile
```

**User-defined tasks and regions:**
```go
import "runtime/trace"

func processOrder(ctx context.Context, orderID string) {
    // Create a task
    ctx, task := trace.NewTask(ctx, "processOrder")
    defer task.End()
    
    // Mark regions within the task
    trace.WithRegion(ctx, "validation", func() {
        validateOrder(orderID)
    })
    
    trace.WithRegion(ctx, "payment", func() {
        chargePayment(orderID)
    })
    
    trace.WithRegion(ctx, "fulfillment", func() {
        fulfillOrder(orderID)
    })
}
```

**When to use tracer vs profiler:**
| Tool | Use For |
|------|---------|
| CPU Profiler | Finding which functions are slow |
| Memory Profiler | Finding allocation-heavy code |
| Tracer | Understanding scheduling, latency, concurrency |

---

## 11. Distributed Systems & Microservices

### 11.1 Service Communication

#### 🟡 Q80: REST vs gRPC — when to use which?

**Answer:**

| Aspect | REST | gRPC |
|--------|------|------|
| **Protocol** | HTTP/1.1 (JSON) | HTTP/2 (Protobuf) |
| **Performance** | Slower (~10x) | Faster (binary, multiplexed) |
| **Browser support** | Native | Requires proxy/gRPC-web |
| **Tooling** | Ubiquitous | Requires protoc |
| **Schema** | Optional (OpenAPI) | Required (protobuf) |
| **Streaming** | Limited | Bidirectional native |

**Use REST for:**
- Public APIs
- Browser clients
- Simple CRUD operations
- When compatibility matters

**Use gRPC for:**
- Internal microservice communication
- High-performance requirements
- Bidirectional streaming
- Strongly-typed contracts

**gRPC in Go:**
```protobuf
// user.proto
syntax = "proto3";
package user;
option go_package = "github.com/myorg/api/user";

service UserService {
    rpc GetUser(GetUserRequest) returns (User);
    rpc ListUsers(ListUsersRequest) returns (stream User);
}

message GetUserRequest {
    string id = 1;
}

message User {
    string id = 1;
    string name = 2;
    string email = 3;
}
```

```go
// Server implementation
type userServer struct {
    pb.UnimplementedUserServiceServer
    repo UserRepository
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    user, err := s.repo.GetByID(ctx, req.Id)
    if err != nil {
        return nil, status.Errorf(codes.NotFound, "user not found")
    }
    return &pb.User{Id: user.ID, Name: user.Name, Email: user.Email}, nil
}

// Server setup
lis, _ := net.Listen("tcp", ":50051")
grpcServer := grpc.NewServer()
pb.RegisterUserServiceServer(grpcServer, &userServer{repo: repo})
grpcServer.Serve(lis)
```

---

### 11.2 Idempotency

#### 🔴 Q81: How do you implement idempotent operations?

**Answer:**

**What is idempotency?** Same request executed multiple times produces same result.

**Pattern 1: Idempotency key**
```go
type IdempotencyStore interface {
    Get(ctx context.Context, key string) (*Response, bool)
    Set(ctx context.Context, key string, resp *Response, ttl time.Duration) error
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
    // Check for existing result
    if cached, ok := s.idempotencyStore.Get(ctx, req.IdempotencyKey); ok {
        return cached, nil
    }
    
    // Process payment
    result, err := s.processPaymentInternal(ctx, req)
    if err != nil {
        return nil, err
    }
    
    // Store result
    s.idempotencyStore.Set(ctx, req.IdempotencyKey, result, 24*time.Hour)
    
    return result, nil
}
```

**Pattern 2: Database constraints**
```go
func (r *OrderRepo) CreateOrder(ctx context.Context, order Order) error {
    _, err := r.db.ExecContext(ctx, `
        INSERT INTO orders (id, user_id, amount, idempotency_key)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (idempotency_key) DO NOTHING
    `, order.ID, order.UserID, order.Amount, order.IdempotencyKey)
    
    if err != nil {
        return err
    }
    return nil
}
```

**Pattern 3: Conditional updates**
```go
func (r *InventoryRepo) DecrementStock(ctx context.Context, productID string, version int, quantity int) error {
    result, err := r.db.ExecContext(ctx, `
        UPDATE inventory 
        SET quantity = quantity - $1, version = version + 1
        WHERE product_id = $2 AND version = $3 AND quantity >= $1
    `, quantity, productID, version)
    
    if err != nil {
        return err
    }
    
    rows, _ := result.RowsAffected()
    if rows == 0 {
        return ErrConcurrentModification
    }
    return nil
}
```

**HTTP methods idempotency:**
| Method | Idempotent | Safe |
|--------|------------|------|
| GET | Yes | Yes |
| HEAD | Yes | Yes |
| PUT | Yes | No |
| DELETE | Yes | No |
| POST | No | No |
| PATCH | No | No |

---

### 11.3 Rate Limiting

#### 🔴 Q82: How do you implement rate limiting in a distributed system?

**Answer:**

**Algorithm 1: Token Bucket**
```go
type TokenBucket struct {
    rate       float64   // Tokens per second
    capacity   float64   // Max tokens
    tokens     float64   // Current tokens
    lastUpdate time.Time
    mu         sync.Mutex
}

func (tb *TokenBucket) Allow() bool {
    tb.mu.Lock()
    defer tb.mu.Unlock()
    
    now := time.Now()
    elapsed := now.Sub(tb.lastUpdate).Seconds()
    tb.tokens = math.Min(tb.capacity, tb.tokens+elapsed*tb.rate)
    tb.lastUpdate = now
    
    if tb.tokens >= 1 {
        tb.tokens--
        return true
    }
    return false
}
```

**Algorithm 2: Sliding Window (Redis)**
```go
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    now := time.Now().UnixNano()
    windowStart := now - window.Nanoseconds()
    
    pipe := r.redis.Pipeline()
    
    // Remove old entries
    pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart, 10))
    
    // Count current window
    countCmd := pipe.ZCard(ctx, key)
    
    // Add current request
    pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
    
    // Set expiry
    pipe.Expire(ctx, key, window)
    
    _, err := pipe.Exec(ctx)
    if err != nil {
        return false, err
    }
    
    return countCmd.Val() < int64(limit), nil
}
```

**Algorithm 3: Fixed Window Counter**
```go
func (r *RateLimiter) AllowFixedWindow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
    windowKey := fmt.Sprintf("%s:%d", key, time.Now().Unix()/int64(window.Seconds()))
    
    count, err := r.redis.Incr(ctx, windowKey).Result()
    if err != nil {
        return false, err
    }
    
    if count == 1 {
        r.redis.Expire(ctx, windowKey, window)
    }
    
    return count <= int64(limit), nil
}
```

**Distributed rate limiting considerations:**
- Use Redis for shared state
- Handle Redis failures gracefully (fail open vs fail closed)
- Consider per-user vs per-endpoint limits
- Return rate limit headers (`X-RateLimit-Limit`, `X-RateLimit-Remaining`)

```go
func RateLimitMiddleware(limiter *RateLimiter) Middleware {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userID := getUserID(r)
            
            allowed, err := limiter.Allow(r.Context(), userID, 100, time.Minute)
            if err != nil {
                // Fail open on error
                next.ServeHTTP(w, r)
                return
            }
            
            if !allowed {
                w.Header().Set("Retry-After", "60")
                http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

---

### 11.4 Circuit Breaker

#### 🔴 Q83: How do you implement a circuit breaker in Go?

**Answer:**

**States:** Closed → Open → Half-Open → Closed

```go
type CircuitBreaker struct {
    name          string
    maxFailures   int
    timeout       time.Duration
    halfOpenMax   int
    
    state         State
    failures      int
    successes     int
    lastFailure   time.Time
    mu            sync.RWMutex
}

type State int

const (
    StateClosed State = iota
    StateOpen
    StateHalfOpen
)

func (cb *CircuitBreaker) Execute(fn func() error) error {
    if !cb.allowRequest() {
        return ErrCircuitOpen
    }
    
    err := fn()
    cb.recordResult(err)
    return err
}

func (cb *CircuitBreaker) allowRequest() bool {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    switch cb.state {
    case StateClosed:
        return true
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = StateHalfOpen
            cb.successes = 0
            return true
        }
        return false
    case StateHalfOpen:
        return cb.successes < cb.halfOpenMax
    }
    return false
}

func (cb *CircuitBreaker) recordResult(err error) {
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        
        if cb.state == StateHalfOpen || cb.failures >= cb.maxFailures {
            cb.state = StateOpen
        }
    } else {
        if cb.state == StateHalfOpen {
            cb.successes++
            if cb.successes >= cb.halfOpenMax {
                cb.state = StateClosed
                cb.failures = 0
            }
        } else {
            cb.failures = 0
        }
    }
}
```

**Using sony/gobreaker:**
```go
import "github.com/sony/gobreaker"

cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
    Name:        "api",
    MaxRequests: 3,                // Max requests in half-open
    Interval:    10 * time.Second, // Clear counts after interval
    Timeout:     30 * time.Second, // Time in open before half-open
    ReadyToTrip: func(counts gobreaker.Counts) bool {
        return counts.ConsecutiveFailures > 5
    },
    OnStateChange: func(name string, from, to gobreaker.State) {
        log.Printf("circuit %s: %s -> %s", name, from, to)
    },
})

result, err := cb.Execute(func() (interface{}, error) {
    return callExternalAPI()
})
```

---

### 11.5 Message Queues

#### 🔴 Q84: How do you work with Kafka in Go?

**Answer:**

**Producer:**
```go
import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

func NewProducer(brokers string) (*kafka.Producer, error) {
    return kafka.NewProducer(&kafka.ConfigMap{
        "bootstrap.servers":  brokers,
        "acks":               "all",
        "retries":            3,
        "enable.idempotence": true,
    })
}

func (p *Producer) Send(ctx context.Context, topic string, key string, value []byte) error {
    deliveryChan := make(chan kafka.Event)
    
    err := p.producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Key:            []byte(key),
        Value:          value,
    }, deliveryChan)
    
    if err != nil {
        return err
    }
    
    select {
    case e := <-deliveryChan:
        m := e.(*kafka.Message)
        if m.TopicPartition.Error != nil {
            return m.TopicPartition.Error
        }
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

**Consumer:**
```go
func NewConsumer(brokers, groupID string) (*kafka.Consumer, error) {
    return kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  brokers,
        "group.id":           groupID,
        "auto.offset.reset":  "earliest",
        "enable.auto.commit": false,
    })
}

func (c *Consumer) Consume(ctx context.Context, topics []string, handler func(*kafka.Message) error) error {
    c.consumer.SubscribeTopics(topics, nil)
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
            if err != nil {
                if err.(kafka.Error).Code() == kafka.ErrTimedOut {
                    continue
                }
                return err
            }
            
            if err := handler(msg); err != nil {
                // Handle error (retry, DLQ, etc.)
                continue
            }
            
            // Commit offset after successful processing
            c.consumer.CommitMessage(msg)
        }
    }
}
```

**Consumer group with rebalancing:**
```go
func (c *Consumer) ConsumeWithRebalance(ctx context.Context, topics []string, handler func(*kafka.Message) error) error {
    c.consumer.SubscribeTopics(topics, func(c *kafka.Consumer, event kafka.Event) error {
        switch e := event.(type) {
        case kafka.AssignedPartitions:
            log.Printf("Assigned: %v", e.Partitions)
            c.Assign(e.Partitions)
        case kafka.RevokedPartitions:
            log.Printf("Revoked: %v", e.Partitions)
            c.Unassign()
        }
        return nil
    })
    
    // ... consume loop
}
```

---

#### 🔴 Q85: How do you ensure exactly-once processing with message queues?

**Answer:**

**Exactly-once is hard. Approaches:**

**1. Idempotent processing:**
```go
func (h *Handler) ProcessMessage(ctx context.Context, msg *kafka.Message) error {
    messageID := string(msg.Key)
    
    // Check if already processed
    processed, err := h.store.IsProcessed(ctx, messageID)
    if err != nil {
        return err
    }
    if processed {
        return nil // Already handled
    }
    
    // Process in transaction
    tx, _ := h.db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    // Business logic
    if err := h.processOrder(ctx, tx, msg.Value); err != nil {
        return err
    }
    
    // Mark as processed
    if err := h.store.MarkProcessed(ctx, tx, messageID); err != nil {
        return err
    }
    
    return tx.Commit()
}
```

**2. Transactional outbox pattern:**
```go
// Instead of publishing directly, write to outbox table
func (s *Service) CreateOrder(ctx context.Context, order Order) error {
    tx, _ := s.db.BeginTx(ctx, nil)
    defer tx.Rollback()
    
    // Insert order
    _, err := tx.ExecContext(ctx, "INSERT INTO orders ...", order)
    if err != nil {
        return err
    }
    
    // Insert to outbox (same transaction)
    event := OrderCreatedEvent{OrderID: order.ID}
    _, err = tx.ExecContext(ctx, `
        INSERT INTO outbox (id, aggregate_type, event_type, payload)
        VALUES ($1, $2, $3, $4)
    `, uuid.New(), "order", "created", event)
    if err != nil {
        return err
    }
    
    return tx.Commit()
}

// Separate process polls outbox and publishes
func (r *OutboxRelay) Run(ctx context.Context) {
    ticker := time.NewTicker(100 * time.Millisecond)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            r.publishPendingEvents(ctx)
        }
    }
}
```

**3. Kafka transactions:**
```go
producer.InitTransactions(ctx)

for msg := range messages {
    producer.BeginTransaction()
    
    producer.Produce(outputMsg, nil)
    producer.SendOffsetsToTransaction(offsets, consumerMeta)
    
    if err := producer.CommitTransaction(ctx); err != nil {
        producer.AbortTransaction(ctx)
    }
}
```

---

### 11.6 Observability

#### 🔴 Q86: How do you implement distributed tracing in Go?

**Answer:**

**Using OpenTelemetry:**
```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
    "go.opentelemetry.io/otel/sdk/trace"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func initTracer() func() {
    exporter, _ := otlptracegrpc.New(context.Background(),
        otlptracegrpc.WithEndpoint("localhost:4317"),
        otlptracegrpc.WithInsecure(),
    )
    
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceNameKey.String("my-service"),
        )),
    )
    
    otel.SetTracerProvider(tp)
    
    return func() { tp.Shutdown(context.Background()) }
}

// HTTP handler with tracing
func main() {
    cleanup := initTracer()
    defer cleanup()
    
    handler := http.HandlerFunc(myHandler)
    wrappedHandler := otelhttp.NewHandler(handler, "my-handler")
    http.ListenAndServe(":8080", wrappedHandler)
}

// Manual span creation
func processOrder(ctx context.Context, orderID string) error {
    tracer := otel.Tracer("order-service")
    ctx, span := tracer.Start(ctx, "processOrder")
    defer span.End()
    
    span.SetAttributes(attribute.String("order.id", orderID))
    
    // Child span
    ctx, childSpan := tracer.Start(ctx, "validateOrder")
    err := validateOrder(ctx, orderID)
    if err != nil {
        childSpan.RecordError(err)
        childSpan.SetStatus(codes.Error, err.Error())
    }
    childSpan.End()
    
    return err
}
```

**Propagating context across services:**
```go
// HTTP client with trace propagation
client := &http.Client{
    Transport: otelhttp.NewTransport(http.DefaultTransport),
}

req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := client.Do(req) // Trace context automatically propagated
```

**gRPC with tracing:**
```go
import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

// Server
grpcServer := grpc.NewServer(
    grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
    grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
)

// Client
conn, _ := grpc.Dial(address,
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
    grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
)
```

---

## 12. System Design Questions (CRITICAL)

### 15 Real System Design Interview Questions

---

#### ⚫ SD1: Design a URL Shortener in Go

**Requirements:**
- Shorten long URLs to short codes
- Redirect short codes to original URLs
- Track click analytics
- Handle 1000 QPS reads, 100 QPS writes

**Key Design Points:**

```
┌─────────┐     ┌──────────────┐     ┌───────────┐
│ Client  │────▶│ Load Balancer│────▶│ API Server│
└─────────┘     └──────────────┘     └─────┬─────┘
                                           │
                          ┌────────────────┴────────────────┐
                          │                                 │
                    ┌─────▼─────┐                    ┌──────▼──────┐
                    │   Redis   │                    │  PostgreSQL │
                    │  (Cache)  │                    │  (Storage)  │
                    └───────────┘                    └─────────────┘
```

**Short code generation:**
```go
// Option 1: Base62 encoding of auto-increment ID
func encode(id int64) string {
    const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var result []byte
    for id > 0 {
        result = append([]byte{charset[id%62]}, result...)
        id /= 62
    }
    return string(result)
}

// Option 2: Hash-based (collision handling needed)
func hashURL(url string) string {
    h := sha256.Sum256([]byte(url))
    return base64.URLEncoding.EncodeToString(h[:])[:8]
}
```

**Database schema:**
```sql
CREATE TABLE urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id UUID,
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);

CREATE INDEX idx_short_code ON urls(short_code);
```

**Scalability considerations:**
- Cache hot URLs in Redis (90%+ read traffic)
- Horizontal scaling of API servers
- Database sharding by short_code hash
- Async analytics processing via Kafka

---

#### ⚫ SD2: Design a Rate Limiter

**Requirements:**
- Per-user and per-endpoint rate limiting
- Distributed (multiple API servers)
- Support multiple algorithms (token bucket, sliding window)
- Return rate limit headers

**Key Design Points:**

```go
type RateLimitConfig struct {
    RequestsPerMinute int
    BurstSize         int
    Algorithm         string // "token_bucket", "sliding_window"
}

type RateLimiter interface {
    Allow(ctx context.Context, key string) (allowed bool, remaining int, resetAt time.Time, err error)
}
```

**Redis-based sliding window:**
```go
func (r *RedisRateLimiter) Allow(ctx context.Context, key string) (bool, int, time.Time, error) {
    now := time.Now()
    windowStart := now.Add(-r.window).UnixMicro()
    
    script := redis.NewScript(`
        redis.call('ZREMRANGEBYSCORE', KEYS[1], 0, ARGV[1])
        local count = redis.call('ZCARD', KEYS[1])
        if count < tonumber(ARGV[2]) then
            redis.call('ZADD', KEYS[1], ARGV[3], ARGV[3])
            redis.call('EXPIRE', KEYS[1], ARGV[4])
            return {1, tonumber(ARGV[2]) - count - 1}
        end
        return {0, 0}
    `)
    
    result, err := script.Run(ctx, r.client, []string{key},
        windowStart, r.limit, now.UnixMicro(), int(r.window.Seconds())).Slice()
    
    allowed := result[0].(int64) == 1
    remaining := int(result[1].(int64))
    resetAt := now.Add(r.window)
    
    return allowed, remaining, resetAt, err
}
```

**Architecture:**
```
                    ┌─────────────────┐
                    │  Rate Limiter   │
                    │   Middleware    │
                    └────────┬────────┘
                             │
              ┌──────────────┴──────────────┐
              │                             │
      ┌───────▼───────┐             ┌───────▼───────┐
      │ Local Cache   │             │    Redis      │
      │ (sync.Map)    │             │   Cluster     │
      │ Short-circuit │             │   (Source)    │
      └───────────────┘             └───────────────┘
```

---

#### ⚫ SD3: Design a Distributed Worker System

**Requirements:**
- Job queue with priorities
- Distributed workers (horizontal scaling)
- Job retry with exponential backoff
- Dead letter queue for failed jobs

**Key Design Points:**

```go
type Job struct {
    ID        string
    Type      string
    Payload   json.RawMessage
    Priority  int
    Attempts  int
    MaxRetry  int
    CreatedAt time.Time
    RunAt     time.Time
}

type Worker interface {
    Process(ctx context.Context, job *Job) error
}

type WorkerPool struct {
    redis     *redis.Client
    handlers  map[string]Worker
    queueName string
    workers   int
}
```

**Job processing with retry:**
```go
func (p *WorkerPool) processJob(ctx context.Context, job *Job) {
    handler, ok := p.handlers[job.Type]
    if !ok {
        p.moveToDLQ(ctx, job, "unknown job type")
        return
    }
    
    err := handler.Process(ctx, job)
    if err != nil {
        job.Attempts++
        if job.Attempts >= job.MaxRetry {
            p.moveToDLQ(ctx, job, err.Error())
            return
        }
        
        // Exponential backoff
        delay := time.Duration(math.Pow(2, float64(job.Attempts))) * time.Second
        job.RunAt = time.Now().Add(delay)
        p.enqueue(ctx, job)
    }
}
```

**Architecture:**
```
┌─────────┐     ┌─────────────┐     ┌─────────────┐
│ Client  │────▶│ Job Queue   │────▶│  Workers    │
│   API   │     │  (Redis)    │     │  (N pods)   │
└─────────┘     └─────────────┘     └──────┬──────┘
                      │                     │
               ┌──────┴──────┐       ┌──────▼──────┐
               │ Priority    │       │   Handler   │
               │  Queues     │       │  Registry   │
               └─────────────┘       └─────────────┘
```

---

#### ⚫ SD4: Design a Real-time Notification Service

**Requirements:**
- Push notifications (WebSocket, SSE, Mobile push)
- 1M+ concurrent connections
- Message persistence
- Delivery guarantees

**Key Design Points:**

```go
type NotificationService struct {
    hub      *WebSocketHub
    producer *kafka.Producer
    store    NotificationStore
}

type Notification struct {
    ID        string
    UserID    string
    Type      string
    Payload   json.RawMessage
    Channels  []string // "websocket", "push", "email"
    CreatedAt time.Time
}
```

**WebSocket hub:**
```go
type WebSocketHub struct {
    clients    map[string]map[*Client]bool // userID -> clients
    register   chan *Client
    unregister chan *Client
    broadcast  chan *Notification
    mu         sync.RWMutex
}

func (h *WebSocketHub) Run() {
    for {
        select {
        case client := <-h.register:
            h.addClient(client)
        case client := <-h.unregister:
            h.removeClient(client)
        case notification := <-h.broadcast:
            h.sendToUser(notification.UserID, notification)
        }
    }
}
```

**Architecture:**
```
                                 ┌─────────────────┐
                                 │   Notification  │
                                 │     Service     │
                                 └────────┬────────┘
                                          │
              ┌───────────────────────────┼───────────────────────────┐
              │                           │                           │
      ┌───────▼───────┐           ┌───────▼───────┐           ┌───────▼───────┐
      │   WebSocket   │           │    Push       │           │    Email      │
      │    Servers    │           │   Gateway     │           │   Service     │
      │  (Stateful)   │           │   (FCM/APNs)  │           │   (Async)     │
      └───────────────┘           └───────────────┘           └───────────────┘
              │
      ┌───────▼───────┐
      │    Redis      │
      │   Pub/Sub     │
      │ (Fan-out)     │
      └───────────────┘
```

---

#### ⚫ SD5: Design a Log Ingestion Pipeline

**Requirements:**
- Handle 100K logs/second
- Searchable within 30 seconds
- Retention for 30 days
- Support structured and unstructured logs

**Key Design Points:**

```go
type LogEntry struct {
    Timestamp   time.Time
    Service     string
    Level       string
    Message     string
    TraceID     string
    Attributes  map[string]interface{}
}

type LogPipeline struct {
    collector  *Collector
    processor  *BatchProcessor
    indexer    *ElasticIndexer
}
```

**Batch processing:**
```go
type BatchProcessor struct {
    input    chan LogEntry
    batchSize int
    flushInterval time.Duration
}

func (p *BatchProcessor) Run(ctx context.Context) {
    batch := make([]LogEntry, 0, p.batchSize)
    ticker := time.NewTicker(p.flushInterval)
    
    for {
        select {
        case entry := <-p.input:
            batch = append(batch, entry)
            if len(batch) >= p.batchSize {
                p.flush(batch)
                batch = batch[:0]
            }
        case <-ticker.C:
            if len(batch) > 0 {
                p.flush(batch)
                batch = batch[:0]
            }
        case <-ctx.Done():
            p.flush(batch)
            return
        }
    }
}
```

**Architecture:**
```
┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│Services │────▶│   Kafka     │────▶│  Processor  │────▶│Elasticsearch│
│(Fluentd)│     │  (Buffer)   │     │  (Flink/Go) │     │  (Search)   │
└─────────┘     └─────────────┘     └─────────────┘     └──────┬──────┘
                                                               │
                                                        ┌──────▼──────┐
                                                        │   Kibana    │
                                                        │   (UI)      │
                                                        └─────────────┘
```

---

#### ⚫ SD6: Design a Distributed Cache

**Requirements:**
- Sub-millisecond latency
- Horizontal scaling
- Cache invalidation
- Support for TTL

**Key Design Points:**
- Consistent hashing for key distribution
- Replication for availability
- Write-through vs write-behind caching
- Cache stampede prevention (singleflight)

```go
import "golang.org/x/sync/singleflight"

type Cache struct {
    local  *sync.Map
    remote *redis.Client
    sf     singleflight.Group
}

func (c *Cache) Get(ctx context.Context, key string) (interface{}, error) {
    // Check local cache
    if val, ok := c.local.Load(key); ok {
        return val, nil
    }
    
    // Singleflight prevents cache stampede
    val, err, _ := c.sf.Do(key, func() (interface{}, error) {
        // Check remote cache
        data, err := c.remote.Get(ctx, key).Result()
        if err == nil {
            c.local.Store(key, data)
            return data, nil
        }
        
        // Fetch from source
        data, err = fetchFromDB(ctx, key)
        if err != nil {
            return nil, err
        }
        
        // Store in caches
        c.remote.Set(ctx, key, data, 15*time.Minute)
        c.local.Store(key, data)
        return data, nil
    })
    
    return val, err
}
```

---

#### ⚫ SD7: Design an API Gateway

**Requirements:**
- Routing to microservices
- Authentication/Authorization
- Rate limiting
- Request/Response transformation
- Circuit breaking

**Key Design Points:**

```go
type Gateway struct {
    router       *Router
    authService  AuthService
    rateLimiter  RateLimiter
    circuitBreakers map[string]*CircuitBreaker
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // 1. Rate limiting
    if !g.rateLimiter.Allow(r.Context(), getClientID(r)) {
        http.Error(w, "rate limit exceeded", 429)
        return
    }
    
    // 2. Authentication
    claims, err := g.authService.Validate(r.Header.Get("Authorization"))
    if err != nil {
        http.Error(w, "unauthorized", 401)
        return
    }
    ctx := context.WithValue(r.Context(), claimsKey, claims)
    
    // 3. Route to backend
    backend := g.router.Match(r.URL.Path, r.Method)
    
    // 4. Circuit breaker
    cb := g.circuitBreakers[backend.Name]
    resp, err := cb.Execute(func() (*http.Response, error) {
        return g.proxy(ctx, backend, r)
    })
    
    // 5. Response transformation
    g.transformResponse(w, resp)
}
```

---

#### ⚫ SD8: Design a Distributed Lock Service

**Requirements:**
- Mutual exclusion across processes
- Fault tolerant
- Deadlock prevention (TTL)
- Lock extension (lease renewal)

**Key Design Points:**

```go
type DistributedLock struct {
    redis   *redis.Client
    key     string
    value   string // Unique identifier
    ttl     time.Duration
    stopCh  chan struct{}
}

func (l *DistributedLock) Acquire(ctx context.Context) error {
    acquired, err := l.redis.SetNX(ctx, l.key, l.value, l.ttl).Result()
    if err != nil {
        return err
    }
    if !acquired {
        return ErrLockNotAcquired
    }
    
    // Start renewal goroutine
    l.stopCh = make(chan struct{})
    go l.renewLease(ctx)
    return nil
}

func (l *DistributedLock) renewLease(ctx context.Context) {
    ticker := time.NewTicker(l.ttl / 3)
    defer ticker.Stop()
    
    for {
        select {
        case <-l.stopCh:
            return
        case <-ctx.Done():
            return
        case <-ticker.C:
            l.redis.Expire(ctx, l.key, l.ttl)
        }
    }
}

func (l *DistributedLock) Release(ctx context.Context) error {
    close(l.stopCh)
    
    // Only delete if we own the lock (Lua script for atomicity)
    script := redis.NewScript(`
        if redis.call("get", KEYS[1]) == ARGV[1] then
            return redis.call("del", KEYS[1])
        end
        return 0
    `)
    return script.Run(ctx, l.redis, []string{l.key}, l.value).Err()
}
```

---

#### ⚫ SD9: Design a Metrics Aggregation System

**Requirements:**
- Collect metrics from 1000s of services
- Support counters, gauges, histograms
- Query last 24 hours at second granularity
- Alerting on thresholds

**Architecture:**
```
┌─────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│Services │────▶│  Prometheus │────▶│   Thanos    │────▶│   Grafana   │
│(Metrics)│     │  (Scrape)   │     │ (Long-term) │     │   (Query)   │
└─────────┘     └─────────────┘     └─────────────┘     └─────────────┘
                      │
              ┌───────▼───────┐
              │ Alertmanager  │
              │  (Alerts)     │
              └───────────────┘
```

---

#### ⚫ SD10: Design a Feature Flag System

**Requirements:**
- Toggle features without deployment
- Gradual rollout (percentage, user segments)
- A/B testing support
- Low latency evaluation

**Key Design Points:**

```go
type FeatureFlag struct {
    Key         string
    Enabled     bool
    Percentage  int // 0-100 for gradual rollout
    UserIDs     []string // Specific users
    Rules       []Rule // Complex targeting
}

type FlagService struct {
    cache    *sync.Map
    redis    *redis.Client
    listener *redis.PubSub
}

func (s *FlagService) IsEnabled(ctx context.Context, flagKey string, userID string) bool {
    flag := s.getFlag(flagKey)
    if flag == nil {
        return false
    }
    
    if !flag.Enabled {
        return false
    }
    
    // Check specific user list
    for _, id := range flag.UserIDs {
        if id == userID {
            return true
        }
    }
    
    // Percentage rollout (consistent per user)
    hash := xxhash.Sum64String(flagKey + userID)
    return int(hash%100) < flag.Percentage
}
```

---

#### ⚫ SD11: Design a Search Autocomplete System

**Requirements:**
- Latency < 100ms
- Support 10K QPS
- Personalization
- Typo tolerance

**Key Design Points:**
- Trie data structure for prefix matching
- Pre-computed suggestions
- Redis sorted sets for ranking
- Bloom filter for spell checking

```go
type Autocomplete struct {
    trie  *Trie
    redis *redis.Client
}

func (a *Autocomplete) Suggest(ctx context.Context, prefix string, userID string, limit int) []string {
    // Get global suggestions
    global, _ := a.redis.ZRevRange(ctx, "autocomplete:"+prefix, 0, int64(limit)).Result()
    
    // Get personalized suggestions
    personal, _ := a.redis.ZRevRange(ctx, "autocomplete:"+userID+":"+prefix, 0, int64(limit/2)).Result()
    
    // Merge and rank
    return mergeAndRank(personal, global, limit)
}
```

---

#### ⚫ SD12: Design a Distributed Configuration Service

**Requirements:**
- Centralized config management
- Real-time config updates
- Version control
- Environment-specific configs

**Architecture:**
```
┌─────────┐     ┌─────────────┐     ┌─────────────┐
│  Admin  │────▶│ Config API  │────▶│   etcd      │
│   UI    │     │  (CRUD)     │     │  (Storage)  │
└─────────┘     └─────────────┘     └──────┬──────┘
                                           │ Watch
                                    ┌──────▼──────┐
                                    │  Services   │
                                    │ (Consumers) │
                                    └─────────────┘
```

---

#### ⚫ SD13: Design a Job Scheduler (like cron)

**Requirements:**
- Schedule jobs with cron expressions
- Distributed (no single point of failure)
- Exactly-once execution
- Job history and monitoring

**Key Design Points:**

```go
type ScheduledJob struct {
    ID           string
    Name         string
    CronExpr     string
    Handler      string
    Payload      json.RawMessage
    NextRunAt    time.Time
    LastRunAt    time.Time
    Status       string
}

type Scheduler struct {
    db     *sql.DB
    locker DistributedLock
}

func (s *Scheduler) Run(ctx context.Context) {
    ticker := time.NewTicker(time.Second)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            s.processDueJobs(ctx)
        }
    }
}

func (s *Scheduler) processDueJobs(ctx context.Context) {
    jobs := s.getDueJobs(ctx)
    for _, job := range jobs {
        // Acquire lock to prevent duplicate execution
        lockKey := "scheduler:job:" + job.ID
        if err := s.locker.Acquire(ctx, lockKey); err != nil {
            continue // Another instance is processing
        }
        
        go func(j ScheduledJob) {
            defer s.locker.Release(ctx, lockKey)
            s.executeJob(ctx, j)
        }(job)
    }
}
```

---

#### ⚫ SD14: Design a Content Delivery Network (CDN) Edge Cache

**Requirements:**
- Cache static content at edge
- Cache invalidation
- Origin shield
- Geo routing

**Key Design Points:**
- Consistent hashing for cache distribution
- Two-tier caching (edge + origin shield)
- Cache-Control header parsing
- Purge API

---

#### ⚫ SD15: Design an Event Sourcing System

**Requirements:**
- Immutable event log
- Rebuild state from events
- Event versioning
- Snapshots for performance

**Key Design Points:**

```go
type Event struct {
    ID            string
    AggregateID   string
    AggregateType string
    EventType     string
    Version       int
    Payload       json.RawMessage
    Timestamp     time.Time
}

type EventStore interface {
    Append(ctx context.Context, aggregateID string, events []Event, expectedVersion int) error
    Load(ctx context.Context, aggregateID string, fromVersion int) ([]Event, error)
}

type Aggregate interface {
    Apply(event Event)
    GetUncommittedEvents() []Event
}

// Rebuild aggregate from events
func (s *EventStore) LoadAggregate(ctx context.Context, aggregate Aggregate, id string) error {
    events, err := s.Load(ctx, id, 0)
    if err != nil {
        return err
    }
    for _, event := range events {
        aggregate.Apply(event)
    }
    return nil
}
```

---

## 13. Rapid Revision Section

### Last Day Before Interview — Go Quick Revision Notes

---

### One-Liners

| Topic | Key Point |
|-------|-----------|
| **Goroutine size** | 2KB initial stack (growable) |
| **GOMAXPROCS** | Default = CPU cores; controls parallelism |
| **Channel zero value** | nil (blocks forever on send/receive) |
| **Slice zero value** | nil (len=0, cap=0, can append) |
| **Map zero value** | nil (read ok, write panics) |
| **Interface nil check** | Interface is nil only if both type AND value are nil |
| **defer order** | LIFO (last in, first out) |
| **panic/recover** | recover() only works in deferred function |
| **context.Background()** | Root context, never cancelled |
| **context.TODO()** | Placeholder when unsure which context to use |
| **error is nil** | Always check immediately after function call |
| **sync.Pool** | Temporary object reuse, cleared on GC |
| **sync.Once** | Execute function exactly once (thread-safe) |
| **sync.Map** | Use for write-once-read-many or disjoint keys only |

---

### Tricky Interview Questions

**Q: What prints?**
```go
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }()
}
time.Sleep(time.Second)
```
**A:** `3 3 3` (closure captures variable, not value) — Fix: `go func(i int) { fmt.Println(i) }(i)`

---

**Q: Is this a memory leak?**
```go
func getFirst(s []byte) []byte {
    return s[:1]
}
```
**A:** Potentially yes. Returned slice references entire backing array. Fix: copy to new slice.

---

**Q: What happens?**
```go
var m map[string]int
m["key"] = 1
```
**A:** Panic! Map is nil. Fix: `m := make(map[string]int)`

---

**Q: Will this compile?**
```go
type I interface { M() }
type S struct{}
func (s S) M() {}

var i I = S{}
```
**A:** Yes. S implements I (implicit interface satisfaction).

---

**Q: What's wrong?**
```go
func returnsError() error {
    var err *MyError = nil
    return err
}
if returnsError() == nil { ... }
```
**A:** Condition is false! Interface has type (*MyError) with nil value ≠ nil interface.

---

**Q: Are these equal?**
```go
var a []int
b := []int{}
```
**A:** `a == nil` is true, `b == nil` is false. But `len(a) == len(b)` is true.

---

### Common Mistakes

| Mistake | Fix |
|---------|-----|
| Not closing response body | `defer resp.Body.Close()` |
| Ignoring context cancellation | Check `ctx.Done()` in loops |
| Goroutine without exit condition | Use context or done channel |
| Range variable capture in goroutine | Shadow: `i := i` or pass as argument |
| Not handling partial writes | Check `n` returned by `Write()` |
| Using `==` for slice comparison | Use `reflect.DeepEqual` or loop |
| Concurrent map write | Use `sync.Mutex` or `sync.Map` |
| defer in loop | Extract to function or close manually |
| Ignoring error from `Close()` | Log or handle file/connection close errors |
| Using `time.Sleep` for synchronization | Use channels or sync primitives |

---

### Memory & Performance Checklist

- [ ] Pre-allocate slices/maps when size is known
- [ ] Use `sync.Pool` for frequently allocated objects
- [ ] Avoid `interface{}` in hot paths
- [ ] Use `strings.Builder` for string concatenation
- [ ] Profile before optimizing (`pprof`)
- [ ] Check escape analysis (`go build -gcflags="-m"`)
- [ ] Avoid defer in tight loops
- [ ] Use value receivers for small structs
- [ ] Align struct fields (largest first)
- [ ] Use buffered channels for async communication

---

### Concurrency Patterns Cheat Sheet

**1. Done channel:**
```go
done := make(chan struct{})
go func() {
    defer close(done)
    work()
}()
<-done
```

**2. Worker pool:**
```go
jobs := make(chan Job, 100)
for i := 0; i < workers; i++ {
    go func() {
        for job := range jobs {
            process(job)
        }
    }()
}
```

**3. Timeout:**
```go
select {
case result := <-ch:
    use(result)
case <-time.After(timeout):
    handleTimeout()
}
```

**4. First response wins:**
```go
result := make(chan int, n)
for _, server := range servers {
    go func(s string) { result <- query(s) }(server)
}
return <-result
```

**5. Graceful shutdown:**
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
server.Shutdown(ctx)
```

---

### Interview Day Reminders

1. **Before coding:** Clarify requirements, discuss trade-offs
2. **While coding:** Think aloud, explain choices
3. **Error handling:** Always handle errors, don't ignore
4. **Concurrency:** Consider race conditions, use appropriate sync
5. **Testing:** Mention how you'd test the code
6. **Optimization:** Mention it, but don't premature optimize
7. **Trade-offs:** There's no perfect solution, discuss alternatives

---

### Quick API Checklist

```go
// HTTP Server Setup
server := &http.Server{
    Addr:         ":8080",
    Handler:      router,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}

// Database Setup
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)

// Context everywhere
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
```

---

## Complete Handbook Summary

| Part | Sections | Questions |
|------|----------|-----------|
| **Part 1** | Go Fundamentals, Memory, Concurrency, Data Structures | Q1-Q45 |
| **Part 2** | Error Handling, Modules, Testing, APIs, Database | Q46-Q73 |
| **Part 3** | Performance, Distributed Systems, System Design, Revision | Q74-Q86 + 15 SD |

**Total:** 86 detailed questions + 15 system design problems

---

> **Good luck with your interview!** 🚀
