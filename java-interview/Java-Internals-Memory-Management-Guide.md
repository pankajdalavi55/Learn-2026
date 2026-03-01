# Java Internals & Memory Management Guide

> **Target Audience:** Staff-level and Senior-level Java Engineers (8+ years experience)
> **Purpose:** Deep-dive reference for interview preparation covering JVM internals, memory management, garbage collection, and advanced interview questions.

---

## Table of Contents

1. [Memory Management in Java (In-Depth)](#1-memory-management-in-java-in-depth)
   - 1.1 JVM Memory Structure
   - 1.2 Object Allocation Process
   - 1.3 Escape Analysis
   - 1.4 Stack vs Heap Memory
   - 1.5 Garbage Collection (Minor, Major, Full GC)
   - 1.6 GC Algorithms (Serial, Parallel, CMS, G1, ZGC, Shenandoah)
   - 1.7 GC Tuning Parameters
   - 1.8 Memory Leaks in Java
   - 1.9 OutOfMemoryError Types & Root Causes
   - 1.10 Monitoring Tools
   - 1.11 Best Practices for Memory Optimization

2. [How a Java Program Works Internally](#2-how-a-java-program-works-internally)
   - 2.1 From .java to Execution
   - 2.2 ClassLoader Mechanism
   - 2.3 JVM Architecture
   - 2.4 Execution Engine (Interpreter & JIT)
   - 2.5 Class Loading Lifecycle
   - 2.6 JIT Optimizations & AOT
   - 2.7 Java Memory Model (JMM)
   - 2.8 Thread Lifecycle & Concurrency Overview

3. [Staff-Level / Senior-Level Interview Questions](#3-staff-level--senior-level-java-interview-questions)
   - 3A. Advanced JVM & Memory Questions
   - 3B. Concurrency & Performance Questions
   - 3C. System Design & Architecture (Java Focused)

---

## 1. Memory Management in Java (In-Depth)

### 1.1 JVM Memory Structure

The JVM divides memory into several distinct runtime data areas. Understanding each region's purpose, sizing, and failure modes is critical for production debugging and performance tuning.

```
┌─────────────────────────────────────────────────────────────────────┐
│                        JVM MEMORY LAYOUT                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  ┌──────────────────────── HEAP (Shared) ───────────────────────┐  │
│  │                                                               │  │
│  │  ┌─────────── Young Generation ───────────┐  ┌────────────┐ │  │
│  │  │ ┌───────┐ ┌─────┐ ┌─────┐             │  │            │ │  │
│  │  │ │ Eden  │ │ S0  │ │ S1  │             │  │    Old     │ │  │
│  │  │ │       │ │(From│ │(To) │             │  │ Generation │ │  │
│  │  │ │  ~80% │ │~10%)│ │~10% │             │  │ (Tenured)  │ │  │
│  │  │ └───────┘ └─────┘ └─────┘             │  │            │ │  │
│  │  └────────────────────────────────────────┘  └────────────┘ │  │
│  └───────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  ┌─────────────────── Non-Heap ──────────────────────────────────┐ │
│  │                                                                │ │
│  │  ┌────────────┐  ┌──────────────┐  ┌───────────────────────┐  │ │
│  │  │ Metaspace  │  │ Code Cache   │  │ Compressed Class Space│  │ │
│  │  │ (Classes,  │  │ (JIT compiled│  │ (Class pointers)      │  │ │
│  │  │  metadata) │  │  native code)│  │                       │  │ │
│  │  └────────────┘  └──────────────┘  └───────────────────────┘  │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  ┌──────────── Per-Thread (Private) ─────────────────────────────┐ │
│  │                                                                │ │
│  │  ┌──────────┐  ┌────────────┐  ┌──────────────────────────┐   │ │
│  │  │  Stack   │  │ PC Register│  │ Native Method Stack      │   │ │
│  │  │ (Frames, │  │ (Current   │  │ (JNI calls, native libs) │   │ │
│  │  │  locals, │  │  bytecode  │  │                          │   │ │
│  │  │  operand │  │  address)  │  │                          │   │ │
│  │  │  stack)  │  │            │  │                          │   │ │
│  │  └──────────┘  └────────────┘  └──────────────────────────┘   │ │
│  └────────────────────────────────────────────────────────────────┘ │
│                                                                     │
│  ┌──────────── Direct / Off-Heap Memory ─────────────────────────┐ │
│  │  ByteBuffer.allocateDirect(), mapped files, Unsafe            │ │
│  └────────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

#### Heap Memory (Shared across all threads)

The heap is where **all object instances and arrays** are allocated. It is the primary region managed by the garbage collector.

**Young Generation (Eden + Survivor Spaces)**

| Region | Purpose | Typical Size |
|--------|---------|-------------|
| **Eden Space** | Where new objects are first allocated via TLAB (Thread-Local Allocation Buffer) | ~80% of Young Gen |
| **Survivor S0 (From)** | Holds objects that survived at least one minor GC | ~10% of Young Gen |
| **Survivor S1 (To)** | Target for copying during minor GC; S0 and S1 swap roles | ~10% of Young Gen |

- Objects are born in Eden (unless they are too large → go directly to Old Gen)
- After surviving a configurable number of GC cycles (`-XX:MaxTenuringThreshold`, default 15 for most collectors), objects are **promoted** (tenured) to Old Generation
- The `-XX:NewRatio=N` flag sets Old:Young ratio (e.g., `NewRatio=2` means Old is 2x Young)
- `-XX:SurvivorRatio=N` sets Eden:Survivor ratio

**Old Generation (Tenured)**

- Stores long-lived objects that survived multiple minor GC cycles
- Collected during **Major GC** or **Full GC** — these are significantly more expensive
- If Old Gen fills up and cannot be collected → `java.lang.OutOfMemoryError: Java heap space`

#### Metaspace (Java 8+, replaces PermGen)

- Stores **class metadata**: class structures, method metadata, constant pool, annotations
- Allocated in **native memory** (not on the JVM heap)
- Grows dynamically by default (no fixed upper bound unless `-XX:MaxMetaspaceSize` is set)
- Each classloader gets its own chunk; when a classloader is GC'd, its entire metaspace chunk is freed
- Common issue: ClassLoader leaks in application servers (redeployment leaks)

```
Key JVM Flags:
  -XX:MetaspaceSize=256m          # Initial threshold for triggering GC of metaspace
  -XX:MaxMetaspaceSize=512m       # Hard upper limit
  -XX:CompressedClassSpaceSize=1g # Limit for compressed class pointer space
```

#### Stack (Per-Thread, Private)

- Each thread gets its own stack at creation time
- Stores **stack frames**: one per method invocation
- Each frame contains: local variables array, operand stack, frame data (constant pool reference, exception table)
- Size controlled by `-Xss` (e.g., `-Xss512k`). Default varies by platform (~512k-1m)
- Overflow → `java.lang.StackOverflowError`
- Deep recursion or excessive thread creation are common causes

```java
// Each call adds a frame to the thread's stack
public int factorial(int n) {
    if (n <= 1) return 1;
    return n * factorial(n - 1); // New stack frame per recursive call
}
// factorial(50000) → StackOverflowError with default -Xss
```

#### PC (Program Counter) Register

- Each thread has its own PC register
- Points to the **address of the current bytecode instruction** being executed
- If executing a **native method**, the PC register is undefined
- Negligible memory footprint; no developer control needed

#### Native Method Stack

- Separate stack for **JNI (Java Native Interface)** calls
- Used when Java invokes C/C++ code via `native` methods
- Implementation-dependent (HotSpot merges it with the Java stack internally)
- Overflow → `java.lang.StackOverflowError` (same error type as Java stack)

#### Direct (Off-Heap) Memory

- Allocated via `ByteBuffer.allocateDirect()`, `Unsafe.allocateMemory()`, or memory-mapped files
- Not managed by GC (except the `DirectByteBuffer` wrapper object on the heap)
- Controlled by `-XX:MaxDirectMemorySize`
- Used heavily by NIO, Netty, and database drivers for zero-copy I/O
- Leak-prone: if wrapper objects are not GC'd, native memory is not freed

```java
// Direct buffer allocation — memory is off-heap
ByteBuffer directBuf = ByteBuffer.allocateDirect(1024 * 1024); // 1 MB off-heap

// Heap buffer — backed by byte[] on the heap
ByteBuffer heapBuf = ByteBuffer.allocate(1024 * 1024); // 1 MB on-heap
```

---

### 1.2 Object Allocation Process

Understanding how the JVM allocates objects is essential for reasoning about performance at scale.

```
Object Allocation Flow:
                                                    
  new Object()                                      
       │                                            
       ▼                                            
  ┌─────────────┐    YES    ┌──────────────────┐   
  │ Escape       ├─────────►│ Stack Allocation  │   
  │ Analysis:    │          │ (no GC needed)    │   
  │ Does object  │          └──────────────────┘   
  │ escape?      │                                  
  └──────┬──────┘                                   
         │ NO (escapes)                             
         ▼                                          
  ┌─────────────┐    YES    ┌──────────────────┐   
  │ Scalar       ├─────────►│ Decomposed into  │   
  │ Replacement  │          │ primitive fields  │   
  │ possible?    │          │ on stack/registers│   
  └──────┬──────┘          └──────────────────┘   
         │ NO                                       
         ▼                                          
  ┌─────────────┐    YES    ┌──────────────────┐   
  │ TLAB has     ├─────────►│ Bump-pointer      │   
  │ space?       │          │ allocation in TLAB│   
  └──────┬──────┘          │ (lock-free, fast) │   
         │ NO              └──────────────────┘   
         ▼                                          
  ┌─────────────┐    YES    ┌──────────────────┐   
  │ New TLAB     ├─────────►│ Allocate new TLAB │   
  │ available?   │          │ in Eden           │   
  └──────┬──────┘          └──────────────────┘   
         │ NO                                       
         ▼                                          
  ┌─────────────┐    YES    ┌──────────────────┐   
  │ Object >     ├─────────►│ Allocate directly │   
  │ threshold?   │          │ in Old Gen        │   
  └──────┬──────┘          └──────────────────┘   
         │ NO                                       
         ▼                                          
  ┌──────────────┐                                  
  │ Trigger       │                                  
  │ Minor GC      │                                  
  │ → retry alloc │                                  
  └──────────────┘                                  
```

**TLAB (Thread-Local Allocation Buffer)**
- Each thread gets a private buffer inside Eden space
- Allocation = bump a pointer (no synchronization required) → extremely fast
- When TLAB is exhausted, a new one is obtained from Eden (requires CAS)
- TLAB sizing: `-XX:TLABSize`, `-XX:+ResizeTLAB` (adaptive by default)
- TLABs eliminate contention on Eden's allocation pointer — critical for high-allocation-rate apps

```java
// Internally, allocation in TLAB is essentially:
// 1. threadLocalTop += objectSize;
// 2. return pointer to old threadLocalTop;
// No locking, no CAS — just a pointer bump within the thread's local buffer
```

**Large Object Allocation**
- Objects exceeding `-XX:PretenureSizeThreshold` (default 0 = disabled for most GCs) go straight to Old Gen
- G1 GC: objects larger than half a region (`-XX:G1HeapRegionSize / 2`) become **humongous objects** allocated in contiguous humongous regions
- Humongous allocations are expensive; they trigger special GC paths

### 1.3 Escape Analysis

Escape analysis is a JIT compiler optimization (C2 compiler) that determines whether an object's reference **escapes** the method or thread scope.

**Three escape states:**

| State | Description | Optimization Possible |
|-------|-------------|----------------------|
| **NoEscape** | Object is only used within the method, never stored to heap or passed elsewhere | Stack allocation, Scalar replacement, Lock elision |
| **ArgEscape** | Object is passed as argument to a method but doesn't escape that call chain | Partial optimizations |
| **GlobalEscape** | Object is stored in a static field, returned from method, or assigned to a heap object | No optimization — must heap-allocate |

**Optimizations enabled by Escape Analysis:**

1. **Stack Allocation** — Object allocated on the stack frame; freed automatically when method returns (no GC)
2. **Scalar Replacement** — Object is decomposed; its fields become local variables (registers/stack)
3. **Lock Elision** — Synchronization on a non-escaping object is removed entirely

```java
// BEFORE Escape Analysis
public long computeSum(int a, int b) {
    // Point object does NOT escape this method
    Point p = new Point(a, b);
    return p.x + p.y;
}

// AFTER Escape Analysis + Scalar Replacement (what JIT effectively does)
public long computeSum(int a, int b) {
    int p_x = a;  // Scalar replacement — no object allocated
    int p_y = b;
    return p_x + p_y;
}
```

```java
// Lock Elision Example
public void process() {
    Object lock = new Object();  // lock doesn't escape
    synchronized (lock) {         // JIT removes this synchronization entirely
        doWork();
    }
}
```

**JVM Flags:**
```
-XX:+DoEscapeAnalysis          # Enabled by default since Java 6u23
-XX:+EliminateAllocations      # Enable scalar replacement (default: on)
-XX:+EliminateLocks            # Enable lock elision (default: on)
-XX:+PrintEscapeAnalysis       # Debug: print escape analysis results (debug build)
```

**Interview Insight:** Escape analysis only works with the **C2 JIT compiler** (not the interpreter or C1). Short-lived methods that never get JIT-compiled won't benefit. Use `-XX:+PrintCompilation` to verify.

### 1.4 Stack vs Heap Memory

| Aspect | Stack | Heap |
|--------|-------|------|
| **Scope** | Per-thread, private | Shared across all threads |
| **Stores** | Primitives, local references, stack frames | Object instances, arrays |
| **Allocation** | Automatic (push on method entry) | Explicit (`new`, reflection, deserialization) |
| **Deallocation** | Automatic (pop on method exit) | Garbage collector |
| **Speed** | Extremely fast (pointer adjustment) | Slower (GC overhead, cache misses) |
| **Size** | Small (default ~512k-1m per thread) | Large (can be multi-GB) |
| **Fragmentation** | Never (LIFO discipline) | Possible (depends on GC algorithm) |
| **Overflow Error** | `StackOverflowError` | `OutOfMemoryError` |
| **Thread Safety** | Inherently thread-safe (private) | Requires synchronization |

**Key Production Consideration:** With 1000 threads at 1MB stack each = 1GB just for thread stacks. Virtual threads (Java 21+) use ~kilobytes of stack, enabling millions of concurrent tasks.

---

### 1.5 Garbage Collection — Minor GC, Major GC, Full GC

#### GC Roots and Reachability

The GC determines liveness through **reachability analysis** from GC roots (not reference counting).

**GC Roots include:**
- Local variables and parameters on active thread stacks
- Active threads themselves
- Static fields of loaded classes
- JNI references
- Internal JVM references (e.g., class objects, system classloader)
- Monitor objects (synchronized locks held)

```
         GC Roots
         ┌──┐ ┌──┐ ┌──┐
         │  │ │  │ │  │   (stack refs, static fields, JNI)
         └─┬┘ └┬─┘ └┬─┘
           │   │    │
           ▼   ▼    ▼
         ┌──┐ ┌──┐ ┌──┐
         │ A│→│ B│ │ C│→──┐   Reachable objects (LIVE)
         └──┘ └─┬┘ └──┘   │
                │          ▼
              ┌─▼┐       ┌──┐
              │ D│       │ E│   Reachable (LIVE)
              └──┘       └──┘

         ┌──┐ ┌──┐
         │ X│→│ Y│   Unreachable from any root → GARBAGE
         └──┘ └──┘
```

#### Types of GC Events

**Minor GC (Young Generation Collection)**
- Triggered when **Eden space is full**
- Scans Eden + active Survivor space
- Copies live objects to the other Survivor space (From ↔ To swap)
- Objects exceeding tenuring threshold → promoted to Old Gen
- Typically takes **5-50ms** for multi-GB heaps
- **Stop-the-world** pause (all application threads halted), but usually short

```
Minor GC Process:
                                                    
  BEFORE                           AFTER            
  ┌───────┬─────┬─────┐          ┌───────┬─────┬─────┐
  │ Eden  │ S0  │ S1  │          │ Eden  │ S0  │ S1  │
  │ FULL  │(has │(empt│          │ EMPTY │(empt│(live│
  │objects│live)│ y)  │          │       │ y)  │objs)│
  └───────┴─────┴─────┘          └───────┴─────┴─────┘
                                          ↑       ↑
                                     cleared   survivors
                                              copied here
                                              
  Objects exceeding age threshold ──────────► Old Gen
```

**Major GC (Old Generation Collection)**
- Collects **Old Generation** only
- Triggered when Old Gen occupancy exceeds a threshold
- Much slower than Minor GC (seconds for large heaps)
- Not all collectors have a distinct "Major GC" — CMS and G1 do concurrent old gen collection
- Stop-the-world duration depends on collector and heap size

**Full GC**
- Collects the **entire heap** (Young + Old) **and** Metaspace
- Most expensive GC event — can cause **multi-second pauses**
- Triggers:
  - `System.gc()` call (avoid in production; use `-XX:+DisableExplicitGC`)
  - Old Gen is full and promotion fails
  - Metaspace exceeds threshold
  - Concurrent GC fails to keep up (Concurrent Mode Failure in CMS)
  - Heap dump requested
- **Production alert**: Full GCs lasting > 1 second are a red flag

#### Stop-the-World (STW) Events

- During STW, all application threads are **halted at safepoints**
- The GC thread(s) execute exclusively
- **Safepoint**: a point in code where all object references are in a known state (e.g., method calls, loop back-edges, allocation points)
- Time-to-safepoint (TTSP) can itself be a latency issue — counted loops without safepoint checks (`-XX:+UseCountedLoopSafepoints` in newer JVMs)

```java
// This loop may delay reaching a safepoint (pre-Java 17 default):
for (int i = 0; i < 1_000_000_000; i++) {
    // JIT may not insert safepoint in counted int loops
    sum += array[i % array.length];
}
// Fix: use long loop variable, or JVM 17+ has -XX:+UseCountedLoopSafepoints by default
```

**Measuring GC impact:**
```
-Xlog:gc*:file=gc.log:time,uptime,level,tags  # Java 9+ unified logging
-XX:+PrintGCDetails -XX:+PrintGCDateStamps     # Java 8
-XX:+PrintGCApplicationStoppedTime             # Total STW time INCLUDING safepoint time
```

---

### 1.6 GC Algorithms — Deep Dive

#### Comparison Matrix

| Collector | Young Gen Strategy | Old Gen Strategy | Pause Target | Best For | Java Version |
|-----------|-------------------|-----------------|-------------|----------|-------------|
| **Serial** | Copy (single-thread) | Mark-Sweep-Compact (single) | None | Small heaps, client apps | All |
| **Parallel** | Copy (multi-thread) | Mark-Sweep-Compact (multi) | `-XX:MaxGCPauseMillis` | Throughput-critical, batch | All (default Java 8) |
| **CMS** | ParNew (copy, multi) | Concurrent Mark-Sweep | Low pause | Legacy low-latency | Deprecated Java 9, Removed 14 |
| **G1** | Evacuate (multi) | Concurrent + mixed GC | `-XX:MaxGCPauseMillis=200` | General purpose, large heaps | Default Java 9+ |
| **ZGC** | Concurrent, colored pointers | Concurrent, load barriers | <1ms (sub-ms) | Ultra-low latency, huge heaps | Production Java 15+ |
| **Shenandoah** | Concurrent, Brooks pointers | Concurrent compaction | <10ms | Low latency (RedHat) | Java 12+ (not in Oracle JDK) |

---

#### Serial GC (`-XX:+UseSerialGC`)

```
Single GC Thread:
  App Threads: ═══════╤══════════════╤═══════════
                      │  STW Pause   │
  GC Thread:          │▓▓▓▓▓▓▓▓▓▓▓▓▓│
                      Mark  Copy/Compact
```

- **Single-threaded** for both young and old generation
- Simple mark-copy (young) and mark-sweep-compact (old)
- Suitable for single-core machines or containers with `<2 CPUs`
- Predictable but long pauses for large heaps
- Use case: small microservices in containers with 256MB-512MB heap

#### Parallel GC (`-XX:+UseParallelGC`) — Default in Java 8

```
Multiple GC Threads (parallel):
  App Threads: ═══════╤══════════════════╤═══════════
                      │    STW Pause     │
  GC Thread 1:       │▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓│
  GC Thread 2:       │▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓│
  GC Thread N:       │▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓│
```

- Also called **Throughput Collector**
- Uses multiple threads for both minor and major GC
- **Still STW** — all app threads paused, but GC completes faster due to parallelism
- Ergonomic goals: `-XX:MaxGCPauseMillis=N` and `-XX:GCTimeRatio=N` (throughput target)
- Thread count: `-XX:ParallelGCThreads=N` (default: based on CPU count)

```
Key Flags:
  -XX:+UseParallelGC
  -XX:ParallelGCThreads=8
  -XX:MaxGCPauseMillis=100           # Soft goal
  -XX:GCTimeRatio=99                 # 99% throughput target (1% GC time)
  -XX:+UseAdaptiveSizePolicy         # JVM auto-tunes gen sizes (default: on)
```

#### CMS (Concurrent Mark-Sweep) — Deprecated, Removed in Java 14

```
CMS Old Gen Collection Phases:
  App Threads: ═══╤══╤═══════════════╤══╤═══════════
                  │  │ Concurrent    │  │
  GC Threads:     │1 │  2     3      │4 │
                  │  │               │  │
  1 = Initial Mark (STW — mark GC roots, very fast)
  2 = Concurrent Mark (runs WITH app threads — traces heap)
  3 = Concurrent Preclean + Abortable Preclean
  4 = Final Remark (STW — handles mutated references)
      → Concurrent Sweep (runs WITH app threads)
```

- Designed for **low-latency** applications
- Does NOT compact → leads to **heap fragmentation** over time
- **Concurrent Mode Failure**: if Old Gen fills before CMS finishes → falls back to Serial Full GC (disaster)
- Required more heap headroom (objects allocated during concurrent marking)
- Replaced by G1; removed in Java 14

#### G1 GC (`-XX:+UseG1GC`) — Default since Java 9

```
G1 Heap Layout (Region-Based):

  ┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
  │ E │ E │ S │ O │ O │ H │ H │ E │ O │ E │ S │ F │
  └───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
  
  E = Eden    S = Survivor    O = Old    H = Humongous    F = Free
  
  Each region: 1-32 MB (power of 2, auto-tuned by JVM)
  Total regions: ~2048 (target)
```

**How G1 works:**

1. **Young-Only Phase**
   - Collects Eden + Survivor regions (STW, parallel, copying collector)
   - Pause-time target driven: adjusts number of regions to collect

2. **Concurrent Marking** (triggered when heap occupancy exceeds `-XX:InitiatingHeapOccupancyPercent`, default 45%)
   - Similar phases to CMS: initial mark → concurrent mark → remark → cleanup
   - Identifies regions with the most garbage ("garbage-first" — hence the name)

3. **Mixed GC Phase**
   - Collects young regions + selected old regions with most garbage
   - Incrementally cleans old gen over multiple mixed GC cycles
   - Controls how many old regions per mix: `-XX:G1MixedGCCountTarget=8`

4. **Full GC** (fallback)
   - Only if mixed GC can't keep up — single-threaded (Java 8) or parallel (Java 10+)
   - **Must be avoided in production** — tune IHOP, heap size, or marking threshold

```
Key G1 Flags:
  -XX:+UseG1GC
  -XX:MaxGCPauseMillis=200            # Default target pause (200ms)
  -XX:G1HeapRegionSize=16m            # Region size (1m-32m, power of 2)
  -XX:InitiatingHeapOccupancyPercent=45  # When to start concurrent marking
  -XX:G1MixedGCCountTarget=8          # Mixed GC cycles to fully clean old regions
  -XX:G1ReservePercent=10             # Reserve heap for promotion during evacuation
  -XX:ConcGCThreads=4                 # Threads for concurrent marking
  -XX:ParallelGCThreads=8             # Threads for STW phases
```

**Real-World G1 Tuning Scenario:**
```bash
# Production microservice: 8GB heap, latency-sensitive API
java -Xms8g -Xmx8g \
     -XX:+UseG1GC \
     -XX:MaxGCPauseMillis=100 \
     -XX:InitiatingHeapOccupancyPercent=35 \
     -XX:G1HeapRegionSize=16m \
     -XX:+ParallelRefProcEnabled \
     -Xlog:gc*:file=/var/log/gc.log:time,uptime,level,tags:filecount=10,filesize=50m \
     -jar my-service.jar
```

#### ZGC (`-XX:+UseZGC`) — Production-Ready since Java 15

```
ZGC Approach:
  App Threads: ═══════════════════════════════════════
  GC Threads:  ▓░▓░▓░▓░▓░▓░▓░▓░▓░▓░  (concurrent)
                                        
  ▓ = Concurrent GC work    
  ░ = Brief STW pauses (<1ms each)    
  
  Almost everything is concurrent — marking, relocation, reference processing
```

**Key Innovations:**
- **Colored Pointers:** Uses unused bits in 64-bit object pointers to store GC metadata (marked, remapped, finalizable)
- **Load Barriers:** Injected by JIT at every object reference load — checks pointer color and fixes if needed
- **No generational division** (until Generational ZGC in Java 21)
- **Concurrent compaction** — relocates objects while app runs
- Sub-millisecond STW pauses regardless of heap size (tested up to **16TB**)

```
ZGC Colored Pointers (64-bit pointer layout):
  
  Bit 63         Bit 47  Bit 46  Bit 45  Bit 44  Bit 43-0
  ┌──────────────┬───────┬───────┬───────┬───────┬──────────────┐
  │   unused     │Marked0│Marked1│Remapped│Final. │ Object Addr  │
  │              │       │       │       │       │ (44 bits =   │
  │              │       │       │       │       │  16TB addr)  │
  └──────────────┴───────┴───────┴───────┴───────┴──────────────┘
```

```
Key ZGC Flags:
  -XX:+UseZGC
  -XX:+ZGenerational                  # Java 21+: Generational ZGC (much better)
  -XX:SoftMaxHeapSize=4g             # Soft limit (GC tries to stay under this)
  -XX:ZCollectionInterval=5          # Force GC every N seconds (proactive)
  -XX:ZAllocationSpikeTolerance=2    # Tolerance for allocation rate spikes
```

**When to use ZGC:**
- Sub-millisecond latency requirements (trading systems, real-time analytics)
- Very large heaps (> 8GB) where G1 pauses become problematic
- Java 21+ with Generational ZGC is the recommended default for new projects

#### Shenandoah GC (`-XX:+UseShenandoahGC`)

- Developed by Red Hat, available in OpenJDK (not Oracle JDK)
- **Concurrent compaction** — similar goals to ZGC but different mechanism
- Uses **Brooks forwarding pointers** (extra word per object) instead of colored pointers
- Pause times: typically <10ms, often <1ms
- Better suited for moderate heap sizes (2-32 GB) vs ZGC's extreme scalability

```
Shenandoah Phases:
  1. Init Mark         (STW, brief)
  2. Concurrent Mark   (with app)
  3. Final Mark        (STW, brief)
  4. Concurrent Cleanup
  5. Concurrent Evacuation (with app — this is the key innovation)
  6. Init Update Refs  (STW, brief)
  7. Concurrent Update Refs
  8. Final Update Refs (STW, brief)
  9. Concurrent Cleanup
```

**ZGC vs Shenandoah:**

| Aspect | ZGC | Shenandoah |
|--------|-----|------------|
| Pointer mechanism | Colored pointers (metadata in pointer bits) | Brooks forwarding pointer (extra word per object) |
| Memory overhead | Pointer metadata (no per-object overhead) | +8 bytes per object |
| Max heap tested | 16 TB | Multi-TB |
| Availability | All OpenJDK & Oracle JDK | OpenJDK only (not Oracle JDK) |
| Best for | Ultra-low latency, huge heaps | Low latency, moderate heaps |

---

### 1.7 GC Tuning Parameters — Production Reference

#### Essential JVM Flags Cheatsheet

```bash
# ──────────────────── Heap Sizing ────────────────────
-Xms4g                          # Initial heap size (set = Xmx to avoid resize pauses)
-Xmx4g                          # Maximum heap size
-XX:NewRatio=2                   # Old:Young ratio (Old = 2x Young → Young = 1/3 heap)
-XX:SurvivorRatio=8              # Eden:Survivor ratio (Eden = 8x one Survivor)
-XX:MaxTenuringThreshold=15      # Promotions needed to move to Old Gen

# ──────────────────── GC Selection ────────────────────
-XX:+UseG1GC                     # G1 (default Java 9+)
-XX:+UseZGC -XX:+ZGenerational   # Generational ZGC (Java 21+)
-XX:+UseParallelGC               # Parallel/Throughput collector
-XX:+UseShenandoahGC             # Shenandoah (OpenJDK only)

# ──────────────────── G1 Tuning ────────────────────
-XX:MaxGCPauseMillis=200         # Target max pause
-XX:InitiatingHeapOccupancyPercent=45  # Start concurrent marking
-XX:G1HeapRegionSize=16m         # Region size (1m-32m)
-XX:G1ReservePercent=10          # Reserve for evacuation failures
-XX:G1MixedGCLiveThresholdPercent=85   # Skip regions above this liveness

# ──────────────────── Metaspace ────────────────────
-XX:MetaspaceSize=256m           # Initial metaspace threshold
-XX:MaxMetaspaceSize=512m        # Cap metaspace growth

# ──────────────────── Thread Stacks ────────────────────
-Xss512k                        # Stack size per thread

# ──────────────────── GC Logging (Java 11+) ────────────────────
-Xlog:gc*:file=gc.log:time,uptime,level,tags:filecount=10,filesize=50m
-Xlog:gc+heap=debug              # Heap details before/after GC
-Xlog:gc+age=trace               # Object age distribution in survivors
-Xlog:safepoint:file=safepoint.log  # Safepoint timing

# ──────────────────── Diagnostics ────────────────────
-XX:+HeapDumpOnOutOfMemoryError  # Auto heap dump on OOM
-XX:HeapDumpPath=/var/dumps/     # Heap dump location
-XX:OnOutOfMemoryError="kill -9 %p"  # Action on OOM
-XX:+ExitOnOutOfMemoryError      # Exit JVM on OOM (for container restarts)
-XX:NativeMemoryTracking=summary # Track native memory usage
```

#### Tuning Methodology (Step-by-Step)

1. **Set `-Xms` = `-Xmx`** — eliminates heap resizing pauses
2. **Enable GC logging** — always, even in production (minimal overhead)
3. **Measure baseline** — GC frequency, pause durations, throughput %
4. **Right-size heap** — Old Gen occupancy after Full GC should be 30-50% of Old Gen
5. **Tune incrementally** — change ONE parameter at a time, measure impact
6. **Set pause target** — `-XX:MaxGCPauseMillis` for G1/Parallel
7. **Monitor promotion rate** — high promotion = objects living just long enough to tenure
8. **Watch for Full GCs** — if frequent, heap may be too small or there's a leak

### 1.8 Memory Leaks in Java — How They Still Happen

Despite garbage collection, memory leaks are common in production Java applications. A "leak" in Java = objects that are **reachable but no longer needed**.

#### Common Memory Leak Patterns

**1. Unintentional Object Retention in Collections**
```java
// LEAK: Map grows indefinitely — keys are never removed
public class SessionCache {
    private static final Map<String, UserSession> sessions = new HashMap<>();
    
    public void addSession(String id, UserSession session) {
        sessions.put(id, session);  // Added but never removed on logout/timeout
    }
    // Fix: Use ConcurrentHashMap with expiration, or WeakHashMap, 
    //       or a bounded cache (Caffeine, Guava Cache)
}
```

**2. Inner Class / Anonymous Class Holding Outer Reference**
```java
// LEAK: Non-static inner class holds implicit reference to outer instance
public class DataProcessor {
    private byte[] largeBuffer = new byte[10_000_000]; // 10 MB
    
    public Runnable createTask() {
        return new Runnable() {  // Anonymous inner class → holds ref to DataProcessor
            @Override
            public void run() { /* ... */ }
        };
    }
    // Fix: Use static inner class or lambda (lambdas don't capture `this` unless used)
}
```

**3. ThreadLocal Not Cleaned Up**
```java
// LEAK: ThreadLocal in thread pools — thread is reused, value persists
private static final ThreadLocal<List<byte[]>> threadCache = new ThreadLocal<>();

public void processRequest() {
    threadCache.set(new ArrayList<>());
    try {
        // ... processing
    } finally {
        threadCache.remove(); // CRITICAL — must remove in finally block!
    }
}
```

**4. ClassLoader Leaks (Application Server Redeployment)**
```
Redeploy #1: ClassLoader A loads classes → objects created
Redeploy #2: ClassLoader B loads new classes
BUT: if any reference (ThreadLocal, static field, JMX, shutdown hook)
     still points to a class loaded by ClassLoader A
     → ClassLoader A cannot be GC'd
     → ALL classes it loaded stay in Metaspace
     → Eventually: OutOfMemoryError: Metaspace
```

**5. Unclosed Resources**
```java
// LEAK: InputStream holds native memory / file descriptors
public String readFile(String path) throws IOException {
    FileInputStream fis = new FileInputStream(path);
    // If exception occurs before close → resource leaked
    // Fix: try-with-resources
    try (FileInputStream fis2 = new FileInputStream(path)) {
        return new String(fis2.readAllBytes());
    }
}
```

**6. String.intern() Abuse**
```java
// LEAK: Interned strings live in the string pool (heap since Java 7)
// Unbounded intern() calls = unbounded memory growth
for (String line : millionsOfLines) {
    processedStrings.add(line.intern()); // DON'T do this with unbounded data
}
```

**7. Listeners and Callbacks Not Deregistered**
```java
// LEAK: Observer registered but never removed
eventBus.register(myListener);
// When myListener should be eligible for GC, eventBus still holds a strong reference
// Fix: Use WeakReference-based observer pattern, or explicitly deregister
```

---

### 1.9 OutOfMemoryError Types and Root Causes

| Error | Root Cause | Investigation |
|-------|-----------|---------------|
| `Java heap space` | Heap exhausted — leak or heap too small | Heap dump analysis (MAT), increase `-Xmx`, check allocation rate |
| `GC overhead limit exceeded` | >98% time in GC, <2% heap recovered | Usually a memory leak; analyze dominators in heap dump |
| `Metaspace` | Too many classes loaded or classloader leak | Check classloader hierarchy, dynamic proxies, bytecode generation (CGLIB, Javassist) |
| `Unable to create new native thread` | OS thread limit hit | Reduce `-Xss`, increase OS `ulimit -u`, use virtual threads |
| `Direct buffer memory` | Off-heap buffers exhausted | Increase `-XX:MaxDirectMemorySize`, check for `ByteBuffer` leaks |
| `Map failed` | Memory-mapped file failure | OS virtual memory limits, check `vm.max_map_count` |
| `Requested array size exceeds VM limit` | Single array > `Integer.MAX_VALUE - 8` elements | Redesign data structure; use chunked arrays |
| `Kill process or sacrifice child` (Linux OOM Killer) | OS-level, not JVM | Check `oom_score_adj`, container memory limits, use `-XX:+ExitOnOutOfMemoryError` |

**Production OOM Response Playbook:**
```bash
# 1. Ensure heap dump is captured automatically
-XX:+HeapDumpOnOutOfMemoryError -XX:HeapDumpPath=/var/dumps/

# 2. Analyze heap dump with Eclipse MAT (Memory Analyzer Tool)
#    - Look at Dominator Tree → largest retained objects
#    - Check Leak Suspects Report (auto-generated)
#    - Examine GC Roots path to suspicious objects

# 3. For Metaspace OOM:
jcmd <pid> VM.classloader_stats       # Classloader hierarchy
jcmd <pid> GC.class_stats             # Class-level memory breakdown

# 4. For native memory:
-XX:NativeMemoryTracking=detail
jcmd <pid> VM.native_memory detail    # Full native memory map
```

### 1.10 Monitoring Tools

#### JVisualVM (bundled until Java 8, standalone afterwards)
- Real-time CPU, memory, threads monitoring
- Heap dump analysis and thread dump capture
- Plugin support (Visual GC for generation-level monitoring)
- **Visual GC plugin**: real-time Eden/Survivor/Old Gen fill visualization
- Limitation: requires JMX connection; not suitable for short-lived processes

#### JConsole
- JMX-based monitoring tool (bundled with JDK)
- Memory tab: real-time heap/non-heap usage with GC overlay
- Thread tab: thread count, deadlock detection
- MBeans tab: access any registered MBean
- Lighter weight than VisualVM; good for quick checks

#### Java Flight Recorder (JFR) + JDK Mission Control (JMC)

```bash
# Enable Flight Recorder (zero/near-zero overhead in production)
-XX:+FlightRecorder
-XX:StartFlightRecording=duration=60s,filename=recording.jfr

# Continuous recording with dump-on-demand
-XX:StartFlightRecording=disk=true,maxage=24h,maxsize=1g,dumponexit=true,filename=continuous.jfr

# Dump from running process
jcmd <pid> JFR.dump name=1 filename=dump.jfr
```

**JFR captures:**
- GC events with detailed pause breakdown
- Object allocation profiling (which methods allocate the most)
- Thread contention (lock wait times)
- IO events (file, socket)
- Method profiling (CPU sampling)
- JIT compilation events
- **This is THE go-to tool for production JVM diagnostics**

#### Command-Line Tools

```bash
# ──────────── jmap: Heap inspection ────────────
jmap -heap <pid>                         # Heap configuration and usage summary
jmap -histo:live <pid>                   # Class histogram (forces GC first)
jmap -dump:live,format=b,file=heap.hprof <pid>  # Generate heap dump

# ──────────── jstack: Thread dumps ────────────
jstack <pid>                             # Thread dump (stack traces)
jstack -l <pid>                          # Include lock info (ownable synchronizers)
# Tip: Take 3-5 thread dumps 5 seconds apart to identify stuck threads

# ──────────── jstat: GC statistics ────────────
jstat -gcutil <pid> 1000 10              # GC stats every 1s, 10 samples
# Output: S0% S1% E% O% M% CCS% YGC YGCT FGC FGCT CGC CGCT GCT

jstat -gc <pid> 1000                    # Raw capacity and usage numbers

# ──────────── jcmd: Swiss army knife (Java 9+) ────────────
jcmd <pid> VM.flags                      # All active JVM flags
jcmd <pid> GC.heap_info                  # Heap summary
jcmd <pid> Thread.print                  # Thread dump
jcmd <pid> VM.native_memory summary      # Native memory (requires NMT)
jcmd <pid> VM.system_properties          # System properties
jcmd <pid> Compiler.queue                # JIT compilation queue

# ──────────── jinfo: Runtime flag inspection ────────────
jinfo -flags <pid>                       # All non-default flags
jinfo -flag MaxHeapSize <pid>            # Check specific flag value
```

**Real-World Investigation Flow:**
```
Alert: High p99 latency spike
  │
  ├─► Check GC logs (gc.log) → Long STW pause?
  │     └─ YES → Check which GC type (Minor/Mixed/Full)
  │              Full GC → Heap too small? Leak? System.gc()?
  │
  ├─► jstat -gcutil → Old Gen growing steadily? (leak indicator)
  │
  ├─► jstack × 3 → Threads BLOCKED on same lock? (contention)
  │
  ├─► JFR recording → Allocation profiling → Which methods allocate most?
  │
  └─► If OOM, analyze heap dump in Eclipse MAT
        └─ Dominator Tree → Find retained heap giants
```

---

### 1.11 Best Practices for Memory Optimization

#### Object Allocation Reduction
1. **Reuse objects where safe** — object pools for expensive objects (e.g., `StringBuilder` reuse in tight loops)
2. **Avoid autoboxing in hot paths** — `Integer` allocation for every `int` → use primitive collections (Eclipse Collections, HPPC)
3. **Use `StringBuilder` instead of `String` concatenation in loops**
4. **Prefer `List.of()`, `Map.of()` (Java 9+)** — compact, unmodifiable, lower memory footprint
5. **Right-size collections** — `new ArrayList<>(expectedSize)` avoids internal array copies

#### Memory-Efficient Data Structures
```java
// BAD: HashMap<Integer, String> — Integer boxing + Node overhead (~48 bytes per entry)
// GOOD: IntObjectHashMap<String> (Eclipse Collections) — ~16 bytes per entry

// BAD: List<Boolean> — each Boolean object = 16 bytes on heap
// GOOD: BitSet — 1 bit per boolean value

// BAD: Enum.values() — creates a new array on every call
// GOOD: Cache it: private static final MyEnum[] VALUES = values();
```

#### GC-Friendly Coding
1. **Short-lived objects are cheap** — allocate, use, discard (collected in Minor GC, fast)
2. **Mid-life objects are expensive** — survive to Old Gen, then die → worst GC pressure
3. **Nullify large references when done** in long-lived objects (helps GC find them sooner)
4. **Avoid finalizers and `finalize()`** — deprecated; causes objects to survive extra GC cycles (phantom reachable queue)
5. **Use `Cleaner` (Java 9+)** instead of finalizers for native resource cleanup
6. **Never call `System.gc()`** in production — use `-XX:+DisableExplicitGC` as safety net

#### Container-Specific Considerations
```bash
# JVM in Docker/K8s: ensure JVM respects container memory limits
# Java 10+: container-aware by default
-XX:+UseContainerSupport              # Default: on (Java 10+)
-XX:MaxRAMPercentage=75.0             # Use 75% of container memory for heap
-XX:InitialRAMPercentage=75.0         # Start with same to avoid resize
# Leave 25% for: Metaspace, thread stacks, direct buffers, OS overhead, code cache

# Example: 4GB container
# MaxRAMPercentage=75 → Xmx = 3GB
# Remaining 1GB: metaspace (~256m), threads (~200m), codeCachee (~240m), OS (~300m)
```

#### Production Memory Configuration Template
```bash
# Recommended production JVM flags for a G1-based microservice in Kubernetes
java \
  -XX:+UseG1GC \
  -XX:MaxRAMPercentage=75.0 \
  -XX:InitialRAMPercentage=75.0 \
  -XX:MaxGCPauseMillis=150 \
  -XX:+ParallelRefProcEnabled \
  -XX:+UseStringDeduplication \
  -XX:+HeapDumpOnOutOfMemoryError \
  -XX:HeapDumpPath=/var/dumps/ \
  -XX:+ExitOnOutOfMemoryError \
  -Xlog:gc*:file=/var/log/gc.log:time,uptime,level,tags:filecount=5,filesize=20m \
  -XX:+FlightRecorder \
  -XX:StartFlightRecording=disk=true,maxage=6h,maxsize=256m,dumponexit=true \
  -jar my-service.jar
```

---

## 2. How a Java Program Works Internally

### 2.1 From .java to Execution — Complete Flow

```
┌──────────────────────────────────────────────────────────────────────────┐
│                    JAVA PROGRAM EXECUTION PIPELINE                       │
└──────────────────────────────────────────────────────────────────────────┘

  ┌──────────┐    javac     ┌──────────┐   ClassLoader   ┌───────────────┐
  │ .java    ├─────────────►│ .class   ├────────────────►│ Runtime Data  │
  │ Source   │  Compilation │ Bytecode │   Loading +     │ Areas (JVM    │
  │ File     │              │ File(s)  │   Linking +     │ Memory)       │
  └──────────┘              └──────────┘   Init          └───────┬───────┘
                                                                  │
                                                                  ▼
                                                          ┌───────────────┐
                                                          │  Execution    │
                                                          │  Engine       │
                                                          │               │
                                                          │ ┌───────────┐ │
                                                          │ │Interpreter│ │
                                                          │ │(bytecode  │ │
                                                          │ │ by line)  │ │
                                                          │ └─────┬─────┘ │
                                                          │       │       │
                                                          │    Hot Code   │
                                                          │    Detected   │
                                                          │       │       │
                                                          │ ┌─────▼─────┐ │
                                                          │ │    JIT    │ │
                                                          │ │ Compiler  │ │
                                                          │ │(C1 → C2) │ │
                                                          │ └─────┬─────┘ │
                                                          │       │       │
                                                          │  Native Code  │
                                                          │  in Code Cache│
                                                          └───────┬───────┘
                                                                  │
                                                                  ▼
                                                           ┌────────────┐
                                                           │  CPU       │
                                                           │  Execution │
                                                           └────────────┘
```

#### Step 1: Writing .java Source File
```java
// HelloWorld.java
public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}
```

#### Step 2: Compilation — `javac`

The Java compiler (`javac`) performs:
1. **Lexical Analysis** — tokenization of source into keywords, identifiers, literals
2. **Syntax Analysis (Parsing)** — builds Abstract Syntax Tree (AST)
3. **Semantic Analysis** — type checking, name resolution, flow analysis
4. **Annotation Processing** — runs annotation processors (e.g., Lombok, MapStruct)
5. **Desugaring** — transforms syntactic sugar:
   - Lambda → invokedynamic + generated method
   - Try-with-resources → nested try-finally
   - Enhanced for-loop → iterator pattern
   - String concatenation → `StringBuilder` or `invokedynamic` (Java 9+)
6. **Bytecode Generation** — produces `.class` files containing platform-independent bytecode

```bash
javac HelloWorld.java          # Produces HelloWorld.class
javap -c HelloWorld            # Disassemble bytecode
javap -v HelloWorld            # Verbose: constant pool, flags, everything
```

#### Step 3: Bytecode

Bytecode is the **intermediate representation** — platform-independent instructions executed by the JVM.

```
// javap -c output for main method:
public static void main(java.lang.String[]);
  Code:
     0: getstatic     #2    // Field java/lang/System.out:Ljava/io/PrintStream;
     3: ldc           #3    // String Hello, World!
     5: invokevirtual #4    // Method java/io/PrintStream.println:(Ljava/lang/String;)V
     8: return

// Bytecode is stack-based (operand stack), NOT register-based
// Each instruction: opcode (1 byte) + operands
// ~200 opcodes defined in JVM spec
```

**Key bytecode categories:**
| Category | Examples | Description |
|----------|---------|-------------|
| Load/Store | `aload`, `istore`, `ldc` | Move values between local vars and operand stack |
| Arithmetic | `iadd`, `lmul`, `dneg` | Type-prefixed arithmetic on operand stack |
| Type conversion | `i2l`, `d2f`, `checkcast` | Widening/narrowing conversions |
| Object manipulation | `new`, `getfield`, `putfield` | Create objects, access fields |
| Stack management | `pop`, `dup`, `swap` | Manipulate operand stack directly |
| Control flow | `ifeq`, `goto`, `tableswitch` | Branching and jumping |
| Method invocation | `invokevirtual`, `invokestatic`, `invokedynamic` | Call methods |
| Exception | `athrow`, exception table entries | Throw/catch exceptions |

---

### 2.2 ClassLoader Mechanism

ClassLoaders are responsible for **finding, loading, and defining** classes at runtime. They follow the **delegation model** (parent-first).

```
ClassLoader Hierarchy & Delegation:

  ┌─────────────────────────────────────┐
  │       Bootstrap ClassLoader         │  ← Written in C/C++ (native)
  │  Loads: java.base module            │  ← rt.jar (Java 8), jrt:/ (Java 9+)
  │         (java.lang.*, java.util.*)  │
  └──────────────────┬──────────────────┘
                     │ parent
  ┌──────────────────▼──────────────────┐
  │     Platform ClassLoader            │  ← (was "Extension ClassLoader" in Java 8)
  │  Loads: java.sql, java.xml,         │  ← ext/ directory (Java 8)
  │         javax.* platform modules    │  ← Platform modules (Java 9+)
  └──────────────────┬──────────────────┘
                     │ parent
  ┌──────────────────▼──────────────────┐
  │     Application ClassLoader         │  ← (System ClassLoader)
  │  Loads: classpath / module-path     │  ← Your application classes
  │         (-cp, -classpath, JAR files)│
  └──────────────────┬──────────────────┘
                     │ parent
  ┌──────────────────▼──────────────────┐
  │     Custom ClassLoaders             │  ← Web app servers (Tomcat),
  │     (User-defined)                  │     OSGi, plugin systems
  └─────────────────────────────────────┘

  Delegation: Child asks Parent first → Parent asks its Parent
              → If none find it, child attempts to load it itself
              → ClassNotFoundException if nobody can load it
```

**Parent-First Delegation:**
```java
// Simplified ClassLoader.loadClass() logic:
protected Class<?> loadClass(String name, boolean resolve) throws ClassNotFoundException {
    // 1. Check if already loaded
    Class<?> c = findLoadedClass(name);
    if (c == null) {
        try {
            // 2. Delegate to PARENT first
            c = parent.loadClass(name, false);
        } catch (ClassNotFoundException e) {
            // 3. Parent couldn't find it → try loading ourselves
            c = findClass(name);
        }
    }
    return c;
}
```

**Key ClassLoader Behaviors:**
- **Uniqueness**: A class is uniquely identified by its **fully qualified name + ClassLoader**
  - Same class loaded by two different ClassLoaders = two different `Class` objects
  - `ClassCastException` even between same-named classes from different loaders
- **Visibility**: Child can see parent's classes, but parent cannot see child's classes
- **Thread context ClassLoader**: `Thread.currentThread().getContextClassLoader()` — used by SPI, JNDI, JAXP to break parent-first delegation

```java
// ClassLoader edge case: Two different "versions" of the same class
ClassLoader loader1 = new URLClassLoader(new URL[]{jarURL}, null);
ClassLoader loader2 = new URLClassLoader(new URL[]{jarURL}, null);

Class<?> c1 = loader1.loadClass("com.example.MyClass");
Class<?> c2 = loader2.loadClass("com.example.MyClass");

System.out.println(c1 == c2);           // false! Different ClassLoader = different Class
System.out.println(c1.equals(c2));       // false!
// c1.cast(c2Instance) → ClassCastException
```

### 2.3 JVM Architecture — Complete Picture

```
┌────────────────────────────────────────────────────────────────────────┐
│                           JVM ARCHITECTURE                              │
├────────────────────────────────────────────────────────────────────────┤
│                                                                        │
│  ┌─────────────── Class Loading Subsystem ──────────────────────────┐ │
│  │                                                                   │ │
│  │  Loading ──► Linking ──► Initialization                          │ │
│  │              │                                                    │ │
│  │              ├─ Verification  (bytecode validity)                 │ │
│  │              ├─ Preparation   (allocate static fields, defaults)  │ │
│  │              └─ Resolution    (symbolic → direct references)      │ │
│  │                                                                   │ │
│  │  ClassLoaders: Bootstrap → Platform → Application → Custom       │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                               │                                        │
│                               ▼                                        │
│  ┌─────────────── Runtime Data Areas ───────────────────────────────┐ │
│  │                                                                   │ │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐                          │ │
│  │  │  Heap   │  │Metaspace│  │  Code   │  ← Shared                │ │
│  │  │(objects)│  │(classes)│  │  Cache  │                           │ │
│  │  └─────────┘  └─────────┘  └─────────┘                          │ │
│  │                                                                   │ │
│  │  Per Thread:                                                      │ │
│  │  ┌───────┐  ┌──────────┐  ┌──────────────────┐                  │ │
│  │  │ Stack │  │PC Register│  │Native Method Stack│                  │ │
│  │  └───────┘  └──────────┘  └──────────────────┘                  │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                               │                                        │
│                               ▼                                        │
│  ┌─────────────── Execution Engine ─────────────────────────────────┐ │
│  │                                                                   │ │
│  │  ┌──────────┐    ┌───────────┐    ┌───────────┐                  │ │
│  │  │Interpreter│    │ JIT       │    │ GC        │                  │ │
│  │  │           │    │ Compiler  │    │           │                  │ │
│  │  │(bytecode  │    │ C1 (fast) │    │(Serial,   │                  │ │
│  │  │ by line)  │    │ C2 (opt)  │    │ G1, ZGC..)│                  │ │
│  │  └──────────┘    └───────────┘    └───────────┘                  │ │
│  └───────────────────────────────────────────────────────────────────┘ │
│                               │                                        │
│                               ▼                                        │
│  ┌─────────────── Native Interface ─────────────────────────────────┐ │
│  │  JNI (Java Native Interface) ←→ Native Libraries (.dll/.so)     │ │
│  └───────────────────────────────────────────────────────────────────┘ │
└────────────────────────────────────────────────────────────────────────┘
```

### 2.4 Class Loading Lifecycle

```
Class Loading Lifecycle:

  ┌──────────┐    ┌─────────────────────────────────┐    ┌──────────────┐
  │ Loading  │───►│           Linking                │───►│Initialization│
  └──────────┘    │                                  │    └──────────────┘
                  │ ┌────────────┐ ┌───────────┐     │
                  │ │Verification│►│Preparation│     │
                  │ └────────────┘ └─────┬─────┘     │
                  │                ┌─────▼──────┐    │
                  │                │ Resolution  │    │
                  │                │ (lazy)      │    │
                  │                └────────────┘    │
                  └─────────────────────────────────┘
```

**Phase 1: Loading**
- Read the `.class` file bytes (from file system, JAR, network, generated at runtime)
- Create a `java.lang.Class` object representing the class in Metaspace
- Store class metadata: method table, field info, constant pool, access flags

**Phase 2: Linking**

- **Verification**: Ensures bytecode is structurally correct and safe
  - Format check (magic number `0xCAFEBABE`, version)
  - Semantic check (type safety, stack consistency, final classes not extended)
  - Bytecode verification (data-flow analysis of each method)
  - Most expensive part of class loading; can be skipped with `-Xverify:none` (dangerous, not recommended)

- **Preparation**: Allocates memory for static fields and sets them to **default values** (0, null, false)
  ```java
  static int count = 42;  // After Preparation: count = 0 (not 42 yet)
  static String name;     // After Preparation: name = null
  ```

- **Resolution**: Converts symbolic references in the constant pool to direct references
  - Method calls: symbolic method name → pointer to method in method table
  - Field access: symbolic field name → offset in object layout
  - Can be **lazy** (resolved on first use) — most JVMs do this

**Phase 3: Initialization**
- Executes the class's **static initializer block** (`<clinit>` method)
- Sets static fields to their **declared values**
- Triggers initialization of parent class first (if not already initialized)
- Thread-safe: JVM guarantees exactly one thread initializes each class (others block)

```java
public class InitOrder {
    // <clinit> executes in textual order:
    static int a = 10;                    // 1. a = 10
    static int b = compute();             // 2. b = compute() called
    static { System.out.println("init"); } // 3. static block executes
    
    static int compute() { return a * 2; } // returns 20
}
```

**When classes are initialized (JLS §12.4.1):**
- `new` instance creation
- Accessing a static field (except `final` compile-time constants)
- Invoking a static method
- Reflection (`Class.forName()`)
- Initializing a subclass triggers parent initialization
- JVM startup (main class)

---

### 2.5 Execution Engine — Interpreter & JIT Compiler

#### Interpreter
- Reads bytecode instruction by instruction and executes each one
- Uses a **dispatch table** mapping opcodes → native implementation routines
- **Fast startup** (no compilation overhead) but **slow sustained throughput**
- Template interpreter (HotSpot): each bytecode has a pre-generated native template

#### JIT Compiler — Tiered Compilation

```
Tiered Compilation (default since Java 8):

  Method Called
       │
       ▼
  ┌────────────┐ invocations < threshold
  │ Level 0    ├─────────────────────────► Interpreted (baseline)
  │ Interpreter│
  └─────┬──────┘
        │ invocations increase
        ▼
  ┌────────────┐
  │ Level 1-3  │ ───► C1 Compiler (Client)
  │ C1 Compiled│      - Fast compilation
  │            │      - Basic optimizations
  │            │      - Inserts PROFILING COUNTERS
  └─────┬──────┘
        │ Hot method detected (invocation + backedge counters)
        ▼
  ┌────────────┐
  │ Level 4    │ ───► C2 Compiler (Server / Opto)
  │ C2 Compiled│      - Aggressive optimizations
  │            │      - Escape analysis
  │            │      - Loop unrolling, vectorization
  │            │      - Method inlining
  │            │      - Dead code elimination
  └────────────┘      - Null check elimination
                      - Range check elimination
                      
  Compilation Thresholds (default):
    C1: ~1,500 invocations (Level 3 with profiling)
    C2: ~10,000 invocations (Level 4 optimized)
```

**Key JIT Optimizations:**

| Optimization | Description | Impact |
|-------------|-------------|--------|
| **Method Inlining** | Replace call site with method body (up to `-XX:MaxInlineSize=35` bytes) | Eliminates call overhead, enables further opts |
| **Escape Analysis** | Stack-allocate non-escaping objects | Eliminates heap allocation + GC pressure |
| **Loop Unrolling** | Replicate loop body to reduce branch count | 2-4x speedup for tight loops |
| **Vectorization (SIMD)** | Auto-vectorize loops using CPU vector instructions (SSE/AVX) | Massive throughput for array ops |
| **Dead Code Elimination** | Remove unreachable or unused code | Smaller code footprint |
| **Constant Folding** | Compute constant expressions at compile time | Zero runtime cost for constants |
| **Devirtualization** | Convert virtual calls to static calls when only one implementation exists | Enables inlining of virtual methods |
| **On-Stack Replacement (OSR)** | JIT-compile a method while it's running (mid-loop) | Long-running loops get optimized without restarting |

```java
// Method inlining example — JIT sees:
public int getX() { return this.x; }
public int compute() { return getX() + 1; }

// After inlining, effectively becomes:
public int compute() { return this.x + 1; }  // No method call overhead
```

```bash
# Useful JIT diagnostic flags:
-XX:+PrintCompilation            # Log each compiled method
-XX:+UnlockDiagnosticVMOptions -XX:+PrintInlining  # Show inlining decisions
-XX:+TraceClassLoading           # Log each class loaded
-XX:CompileThreshold=10000       # Legacy: compilation invocation threshold (non-tiered)
```

**Deoptimization:** JIT can **undo** optimizations when assumptions are invalidated:
```java
// JIT devirtualized call to Animal.speak() assuming only Dog exists
// Then Cat class is loaded → deoptimization: back to interpreted/C1, re-profile, recompile
abstract class Animal { abstract void speak(); }
class Dog extends Animal { void speak() { /* bark */ } }
// Later loaded: class Cat extends Animal { void speak() { /* meow */ } }
// → JIT must revert devirtualization of speak()
```

### 2.6 JIT Optimizations & AOT (Ahead-of-Time Compilation)

#### AOT Compilation
- **GraalVM Native Image**: Compiles Java to native executable at build time
  - No JVM, no warmup, instant startup (milliseconds)
  - Lower memory footprint (no JIT compiler, no class metadata for unused classes)
  - Trade-off: peak throughput may be lower than JIT-optimized code (no profile-guided optimization at runtime)
  - Closed-world assumption: all classes must be known at build time (reflection requires configuration)

- **JEP 295 (jaotc)**: Experimental AOT in OpenJDK (removed in Java 17)
- **CRaC (Coordinated Restore at Checkpoint)**: Checkpoint a running JVM, restore it later (near-instant startup with JIT warmup preserved)

```bash
# GraalVM Native Image
native-image -jar myapp.jar        # Produces native executable
./myapp                             # Starts in ~20ms vs ~2s for JVM

# CRaC (Project CRaC)
java -XX:CRaCCheckpointTo=checkpoint-dir -jar myapp.jar
# ... app warms up, then trigger checkpoint
jcmd <pid> JDK.checkpoint
# Later: restore
java -XX:CRaCRestoreFrom=checkpoint-dir
```

### 2.7 Reflection & Dynamic Class Loading

```java
// Reflection bypasses compile-time type checking
Class<?> clazz = Class.forName("com.example.MyService");     // Dynamic loading
Object instance = clazz.getDeclaredConstructor().newInstance(); // Reflective instantiation
Method method = clazz.getMethod("process", String.class);
Object result = method.invoke(instance, "data");               // Reflective invocation

// Performance: Reflection calls are ~10-100x slower than direct calls
// JIT can optimize reflective calls after seeing them enough times (inflation)
// -Dsun.reflect.inflationThreshold=15 (default): after 15 calls, generates bytecode accessor

// Reflection + Modules (Java 9+):
// Must open packages for deep reflection:
// --add-opens java.base/java.lang=ALL-UNNAMED
```

**Method Invocation Internals:**
```
Method Call Types (bytecode instructions):

  invokestatic    → Static methods (resolved at compile time)
  invokevirtual   → Instance methods (vtable dispatch — virtual method table)
  invokeinterface → Interface methods (itable dispatch — interface method table)
  invokespecial   → Constructors, super calls, private methods (no dispatch)
  invokedynamic   → Lambda expressions, string concat (Java 9+), 
                     dynamic languages (bootstrap method resolves target)
                     
  vtable:
  ┌─────────┐     ┌────────────────────┐
  │ Object  │     │ toString() → addr1 │
  │ header  │────►│ equals()   → addr2 │
  │ (klass  │     │ hashCode() → addr3 │
  │  pointer│     │ myMethod() → addr4 │ ← overridden = different addr
  └─────────┘     └────────────────────┘
```

---

### 2.8 Java Memory Model (JMM) & Happens-Before

The JMM (JSR-133, Java 5+) defines how threads interact through memory — what values a thread is **guaranteed** to see when reading shared variables.

#### Why JMM Matters
```
Without JMM guarantees:

  Thread 1 (CPU Core 1)         Thread 2 (CPU Core 2)
  ┌──────────────────┐         ┌──────────────────┐
  │ x = 42;          │         │ while (!ready) {} │
  │ ready = true;    │         │ print(x);         │
  └──────────────────┘         └──────────────────┘
  
  L1 Cache: x=42, ready=true   L1 Cache: may see ready=true, x=0
  
  Problem: CPU/compiler may REORDER instructions or cache stale values
  → Thread 2 might see ready=true but x=0 (or even never see ready=true)
  
  JMM defines rules to prevent such visibility and ordering issues.
```

#### Happens-Before Rules (Formal Guarantees)

| Rule | Description |
|------|-------------|
| **Program Order** | Within a single thread, each statement happens-before the next |
| **Monitor Lock** | An unlock on a monitor HB a subsequent lock on that same monitor |
| **Volatile Variable** | A write to `volatile` HB a subsequent read of that same field |
| **Thread Start** | `thread.start()` HB any action in the started thread |
| **Thread Termination** | Any action in a thread HB `thread.join()` returning |
| **Interruption** | `thread.interrupt()` HB the interrupted thread detecting it |
| **Finalizer** | Constructor completion HB start of `finalize()` |
| **Transitivity** | If A HB B and B HB C, then A HB C |

```java
// Volatile establishes happens-before
volatile boolean ready = false;
int value = 0;

// Thread 1:
value = 42;          // (1)
ready = true;        // (2) volatile write — publishes ALL preceding writes

// Thread 2:
if (ready) {         // (3) volatile read — sees (2) and everything before it
    print(value);    // (4) GUARANTEED to print 42 (not 0)
}
// Happens-before chain: (1) →program→ (2) →volatile→ (3) →program→ (4)
```

#### volatile vs synchronized vs Atomic

```java
// ── volatile ──
// Guarantees: visibility (no caching) + happens-before ordering
// Does NOT provide: atomicity for compound operations (i++ is NOT atomic)
volatile int counter = 0;
counter++;  // NOT thread-safe! Read + increment + write = 3 operations

// ── synchronized ──
// Guarantees: mutual exclusion + visibility + happens-before
// Provides: atomicity (only one thread in critical section)
synchronized (lock) {
    counter++;  // Thread-safe: exclusive access
}

// ── Atomic classes ──
// Guarantees: atomic compound operations via CAS (Compare-And-Swap)
AtomicInteger counter = new AtomicInteger(0);
counter.incrementAndGet();  // Thread-safe, lock-free, uses CPU CAS instruction
// compareAndSet(expected, new) → hardware-level atomic operation
```

### 2.9 Thread Lifecycle & Concurrency Overview

```
Thread State Machine:

       new Thread()
            │
            ▼
      ┌──────────┐    start()    ┌──────────┐
      │   NEW    ├──────────────►│ RUNNABLE │◄─────────────────┐
      └──────────┘               └────┬─────┘                  │
                                      │                         │
                    ┌─────────────────┼─────────────────┐      │
                    │                 │                  │      │
                    ▼                 ▼                  ▼      │
              ┌──────────┐    ┌────────────┐    ┌─────────────┐│
              │ BLOCKED  │    │  WAITING   │    │TIMED_WAITING││
              │(monitor  │    │(.wait(),   │    │(.sleep(n),  ││
              │ contention│   │ .join(),   │    │ .wait(n),   ││
              │          │    │ LockSupport│    │ .join(n))   ││
              │          │    │ .park())   │    │             ││
              └────┬─────┘    └─────┬──────┘    └──────┬──────┘│
                   │                │                   │       │
                   └────────────────┴───────────────────┘───────┘
                                                   (notify/timeout/
                                                    signal/complete)
                                      │
                                      ▼
                                ┌──────────┐
                                │TERMINATED│  (run() completed or exception)
                                └──────────┘
```

**Thread Pool Internals (ThreadPoolExecutor)**

```java
// Core parameters that define pool behavior:
ThreadPoolExecutor executor = new ThreadPoolExecutor(
    corePoolSize,      // Threads kept alive even when idle
    maximumPoolSize,   // Max threads created under load
    keepAliveTime,     // Idle time before excess threads are terminated
    TimeUnit.SECONDS,
    workQueue,         // Queue for tasks when all core threads are busy
    threadFactory,     // Custom thread naming, daemon flag
    rejectionHandler   // What to do when queue AND pool are full
);

// Execution flow:
// 1. task arrives → core threads available? → execute immediately
// 2. core threads busy → queue has space? → enqueue task
// 3. queue full → can create thread (< maxPoolSize)? → create new thread
// 4. at maxPoolSize AND queue full → rejection handler fires
```

**Common Thread Pool Mistakes:**
- Using `Executors.newFixedThreadPool(N)` → unbounded `LinkedBlockingQueue` → OOM from queued tasks
- Using `Executors.newCachedThreadPool()` → unbounded threads → thread explosion under load
- **Production best practice**: Always use `ThreadPoolExecutor` directly with bounded queue + explicit rejection policy

```java
// Production-grade thread pool:
ThreadPoolExecutor executor = new ThreadPoolExecutor(
    10,                                    // core
    50,                                    // max
    60, TimeUnit.SECONDS,                  // keep-alive
    new ArrayBlockingQueue<>(1000),        // BOUNDED queue
    new ThreadFactoryBuilder()
        .setNameFormat("api-worker-%d")
        .setDaemon(true)
        .build(),
    new ThreadPoolExecutor.CallerRunsPolicy() // Back-pressure: caller thread runs task
);
```

---

## 3. Staff-Level / Senior-Level Java Interview Questions

### 3A. Advanced JVM & Memory Questions

---

#### Q1: Your production microservice experiences increasing GC pause times over several days, eventually reaching 5-second Full GCs. Walk me through your investigation process.

**Expected Senior-Level Answer:**

1. **Check GC logs**: look for gradual increase in Old Gen occupancy post-GC → memory leak indicator
2. **Run `jstat -gcutil`**: confirm Old Gen usage trending upward across multiple GC cycles
3. **Capture heap dump**: `-XX:+HeapDumpOnOutOfMemoryError` or `jmap -dump:live,format=b,file=heap.hprof <pid>`
4. **Analyze with Eclipse MAT**: 
   - Dominator Tree → find largest retained objects
   - Leak Suspects Report → auto-detected accumulation points
   - Path to GC Roots → understand why objects aren't collected
5. **Common findings**: unbounded cache, ThreadLocal not cleaned in thread pool, listener not deregistered, ClassLoader leak on redeploy
6. **Fix**: bound the cache (TTL + max size), add `finally { threadLocal.remove(); }`, fix resource lifecycle

**Common Mistakes:**
- Jumping to increasing `-Xmx` without investigating the root cause (just delays the problem)
- Not enabling GC logging in production (missed diagnostic data)
- Using `-Xverify:none` or `-XX:+DisableExplicitGC` as "fixes"

**Follow-Up Probing Questions:**
- How would you differentiate a memory leak from an under-provisioned heap?
- What if the leak is in Metaspace instead of Heap?
- How does the investigation differ in a containerized (K8s) environment?

---

#### Q2: Explain the difference between G1's concurrent marking cycle and a Full GC. When does G1 fall back to Full GC?

**Expected Senior-Level Answer:**

- **Concurrent marking**: runs concurrently with application threads (Init Mark STW → Concurrent Mark → Remark STW → Cleanup). Identifies garbage regions. Followed by mixed GC which incrementally collects old regions alongside young collections.
- **Full GC**: single-threaded (Java 8) or parallel (Java 10+) STW collection of the entire heap. Compacts everything.
- **Fallback triggers**:
  - Evacuation failure: not enough free regions to copy surviving objects during young/mixed collection
  - Humongous allocation failure: can't find contiguous regions
  - "To-space exhausted": survivor objects can't be evacuated
- **Prevention**: increase heap, lower `-XX:InitiatingHeapOccupancyPercent` (start marking earlier), increase `-XX:G1ReservePercent`

**Common Mistakes:**
- Confusing "mixed GC" with "Full GC" — mixed is incremental and pause-time-targeted
- Not knowing that G1 Full GC was single-threaded before Java 10

**Follow-Up Probing Questions:**
- How does G1 determine which regions to include in a mixed collection?
- What is the relationship between IHOP and concurrent marking start?
- How would you diagnose "to-space exhausted" in GC logs?

---

#### Q3: A Java application creates 50,000 threads. What JVM memory areas are impacted, and what errors can occur?

**Expected Senior-Level Answer:**

- **Stack memory**: 50,000 × `-Xss` (default ~1MB) = ~50GB native memory for thread stacks alone
- **Native memory**: each thread requires OS-level resources (pthread, kernel stack, file descriptors)
- **Errors**:
  - `OutOfMemoryError: unable to create new native thread` — OS limits (`ulimit -u`, `pid_max`)
  - `OutOfMemoryError: Java heap space` — if thread objects and their Runnables consume heap
  - OS OOM killer may terminate the process
- **Solutions**:
  - Reduce `-Xss` (512k or 256k if call stacks are shallow)
  - Use thread pool with bounded concurrency
  - Java 21+ Virtual Threads: ~kilobytes per virtual thread, can handle millions

**Common Mistakes:**
- Forgetting that thread stack memory is off-heap (not counted in `-Xmx`)
- Not considering OS-level thread limits (`/proc/sys/kernel/threads-max`)

**Follow-Up Probing Questions:**
- How do virtual threads (Project Loom) solve this differently from thread pools?
- What's the relationship between `-Xss` and `StackOverflowError` vs the total memory impact?

---

#### Q4: Explain ClassLoader leaks. How do they happen during application redeployment, and how would you prevent them?

**Expected Senior-Level Answer:**

- When a web app is redeployed, the app server creates a **new ClassLoader** and tries to GC the old one
- The old ClassLoader can't be GC'd if **any reference chain** from a GC root reaches any class loaded by it
- **Common leak sources**: 
  - `ThreadLocal` values holding references to app classes
  - JDBC drivers registered in `DriverManager` (static, loaded by parent classloader)
  - Shutdown hooks referencing app classes
  - JMX MBeans registered with platform MBeanServer
  - Logging frameworks holding references to app ClassLoaders
  - Static collections in library classes that hold app-level objects
- **Impact**: Metaspace grows with each redeploy → `OutOfMemoryError: Metaspace`
- **Prevention**:
  - Clean up ThreadLocals on context shutdown
  - Deregister JDBC drivers in `ServletContextListener.contextDestroyed()`
  - Remove MBeans on undeploy
  - Use leak detection: Tomcat's `JreMemoryLeakPreventionListener`

**Common Mistakes:**
- Thinking `static` fields in app classes are automatically cleaned up (they're not if the ClassLoader isn't GC'd)
- Blaming the app server without investigating application-level root causes

**Follow-Up Probing Questions:**
- How would you detect which ClassLoader cannot be GC'd from a heap dump?
- What's the difference between PermGen leak (Java 7) vs Metaspace leak (Java 8+)?

---

#### Q5: What happens internally when you call `new Object()` at the bytecode and JVM level?

**Expected Senior-Level Answer:**

1. **Bytecode**: `new` instruction → allocates memory and pushes reference on operand stack; `invokespecial <init>` → runs constructor
2. **Memory allocation**: 
   - JVM checks TLAB (Thread-Local Allocation Buffer) → bump-pointer allocation if space exists
   - If TLAB exhausted → allocate new TLAB in Eden (CAS on Eden allocation pointer)
   - If Eden full → trigger Minor GC → retry
3. **Object header setup**: 
   - Mark Word (8 bytes): hash code, GC age, lock state, biased thread ID
   - Klass Pointer (4 bytes with compressed oops): pointer to class metadata in Metaspace
   - Array length (4 bytes, only for arrays)
4. **Zero-initialization**: all instance fields set to default values (0/null/false)
5. **Constructor execution**: `<init>` method runs (field initializers + constructor body)

```
Object Memory Layout (64-bit JVM, compressed oops):

  ┌────────────────────────────────────┐
  │ Mark Word          (8 bytes)       │  hash | age | lock | GC bits
  ├────────────────────────────────────┤
  │ Klass Pointer      (4 bytes)       │  → Class metadata in Metaspace
  ├────────────────────────────────────┤
  │ Instance Fields    (variable)      │  Packed by JVM (field reordering)
  ├────────────────────────────────────┤
  │ Padding            (to align 8B)   │  JVM aligns objects to 8-byte boundary
  └────────────────────────────────────┘

  Minimum object size: 16 bytes (header + padding)
  new Object(): 16 bytes on heap
  new int[0]:   16 bytes (header + array length field)
```

**Common Mistakes:**
- Not knowing about TLAB or thinking all allocation requires synchronization
- Forgetting object alignment padding (can significantly impact memory for small objects)

**Follow-Up Probing Questions:**
- What is the overhead of object headers and how does compressed oops reduce it?
- How does the JVM reorder fields to minimize padding waste?

---

### 3B. Concurrency & Performance Questions

---

#### Q6: Explain the Java Memory Model's happens-before relationship. Why is `volatile` not sufficient for compound operations?

**Expected Senior-Level Answer:**

- JMM defines a **partial ordering** of actions (reads, writes, locks, etc.) that guarantees memory visibility between threads
- `volatile` guarantees:
  - **Visibility**: write is immediately visible to all threads (no CPU cache staleness)
  - **Ordering**: acts as a memory barrier — instructions before volatile write can't be reordered after it
- `volatile` does **NOT** guarantee **atomicity** for compound operations:
  ```java
  volatile int count = 0;
  count++;  // = read count → increment → write count (3 separate operations)
  // Two threads can both read 5, increment to 6, write 6 → lost update
  ```
- For compound atomics, use: `AtomicInteger`, `synchronized`, or `VarHandle` (Java 9+)
- `VarHandle` provides finer-grained memory ordering modes: `plain`, `opaque`, `acquire/release`, `volatile`

**Common Mistakes:**
- Believing `volatile` makes all operations on that variable thread-safe
- Not understanding that the happens-before guarantee extends to ALL writes before the volatile write (not just the volatile field)

**Follow-Up Probing Questions:**
- What's the difference between `acquire/release` semantics and `volatile` (sequential consistency)?
- When would you use `VarHandle` over `AtomicReference`?

---

#### Q7: What is false sharing, and how do you prevent it in Java?

**Expected Senior-Level Answer:**

- **False sharing** occurs when two threads modify different variables that happen to reside on the **same CPU cache line** (typically 64 bytes)
- CPU cache coherence protocol (MESI) forces the entire cache line to be invalidated and reloaded, even though the variables are independent → severe performance degradation

```
Cache Line (64 bytes):
  ┌──────────────────────────────────────────────────────────────────┐
  │ Thread1's counter (8 bytes) │ Thread2's counter (8 bytes) │ ... │
  └──────────────────────────────────────────────────────────────────┘
  
  Thread 1 writes its counter → invalidates ENTIRE cache line
  → Thread 2 must reload cache line from L3/main memory to read its own counter
  → Even though they access DIFFERENT fields, they contend on the SAME cache line
```

**Prevention in Java:**
```java
// Java 8+: @Contended annotation (requires -XX:-RestrictContended for non-JDK classes)
@jdk.internal.vm.annotation.Contended
public class PaddedCounter {
    volatile long count;  // Padded to occupy its own cache line
}

// Manual padding (pre-Java 8):
public class PaddedCounter {
    volatile long p1, p2, p3, p4, p5, p6, p7; // 56 bytes padding
    volatile long count;                         // Target field
    volatile long q1, q2, q3, q4, q5, q6, q7; // 56 bytes padding
}

// java.util.concurrent LongAdder uses @Contended internally
// to pad its "Cell" array entries → much faster than AtomicLong under contention
```

**Common Mistakes:**
- Micro-optimizing for false sharing when the real bottleneck is algorithmic
- Not verifying with benchmarks (JMH) — padding adds memory overhead

**Follow-Up Probing Questions:**
- How do you detect false sharing using hardware performance counters?
- Why is `LongAdder` faster than `AtomicLong` under high contention?

---

#### Q8: Design an efficient `CompletableFuture`-based pipeline for an API that aggregates data from 3 downstream services with timeout and fallback.

**Expected Senior-Level Answer:**

```java
public CompletableFuture<AggregatedResponse> aggregateData(String requestId) {
    
    ExecutorService executor = customBoundedExecutor; // Never use ForkJoinPool.commonPool() for I/O
    
    // Fan-out: 3 parallel async calls
    CompletableFuture<ServiceAData> futureA = CompletableFuture
        .supplyAsync(() -> serviceA.fetch(requestId), executor)
        .orTimeout(2, TimeUnit.SECONDS)                    // Hard timeout (Java 9+)
        .exceptionally(ex -> {
            log.warn("ServiceA failed: {}", ex.getMessage());
            return ServiceAData.fallbackDefault();          // Graceful degradation
        });
    
    CompletableFuture<ServiceBData> futureB = CompletableFuture
        .supplyAsync(() -> serviceB.fetch(requestId), executor)
        .orTimeout(2, TimeUnit.SECONDS)
        .exceptionally(ex -> ServiceBData.fallbackDefault());
    
    CompletableFuture<ServiceCData> futureC = CompletableFuture
        .supplyAsync(() -> serviceC.fetch(requestId), executor)
        .orTimeout(3, TimeUnit.SECONDS)  // ServiceC gets longer timeout (less critical)
        .exceptionally(ex -> ServiceCData.fallbackDefault());
    
    // Fan-in: combine all results
    return futureA.thenCombine(futureB, (a, b) -> new PartialResult(a, b))
                  .thenCombine(futureC, (partial, c) -> new AggregatedResponse(partial, c));
}

// Key design decisions:
// 1. Custom executor (not ForkJoinPool.commonPool) — bounded threads for I/O
// 2. Per-service timeout — different SLAs per dependency
// 3. Fallback on failure — partial data is better than total failure
// 4. No .join()/.get() in async chain — fully non-blocking pipeline
```

**Common Mistakes:**
- Using `ForkJoinPool.commonPool()` for blocking I/O → starves CPU-bound work
- Calling `.get()` without timeout → potential infinite blocking
- Not handling exceptions → `CompletableFuture` silently swallows them

**Follow-Up Probing Questions:**
- How would you add circuit breaker pattern to this?
- What happens if the executor's queue is full when `supplyAsync` is called?
- How would you trace a request across these async calls (distributed tracing)?

---

#### Q9: How does ForkJoinPool's work-stealing algorithm work, and when should you prefer it over ThreadPoolExecutor?

**Expected Senior-Level Answer:**

- **ForkJoinPool** is designed for **recursive, divide-and-conquer** tasks
- Each worker thread has a **deque (double-ended queue)**
  - Worker pushes/pops tasks from the **bottom** (LIFO — hot cache locality)
  - **Stealing** happens from the **top** (FIFO — oldest/largest tasks)
- When a thread's deque is empty, it **steals** from other threads' deques
- Minimizes thread idling → high CPU utilization for parallelizable work

```
Work-Stealing:

  Thread 1 Deque:    Thread 2 Deque:    Thread 3 Deque:
  ┌─────────┐       ┌─────────┐       ┌─────────┐
  │ Task A  │←steal │         │       │ Task F  │
  ├─────────┤       │ (empty  │       ├─────────┤
  │ Task B  │       │  deque, │       │ Task G  │
  ├─────────┤       │  steals │       ├─────────┤
  │ Task C  │       │  from   │       │ Task H  │←push/pop
  └─────────┘       │  others)│       └─────────┘
        ↑           └─────────┘
     push/pop
```

- **Prefer ForkJoinPool when**: recursive parallelism (sort, tree processing, parallel streams)
- **Prefer ThreadPoolExecutor when**: independent I/O-bound tasks, producer-consumer patterns, backpressure needed

**Common Mistakes:**
- Using `ForkJoinPool.commonPool()` for blocking I/O (starves the shared pool)
- Not using `ManagedBlocker` when blocking operations are unavoidable in ForkJoinPool

**Follow-Up Probing Questions:**
- How does `parallelStream()` use ForkJoinPool internally?
- What happens if a ForkJoinTask does blocking I/O without `ManagedBlocker`?

---

#### Q10: How would you debug a deadlock in production? Walk through the full process.

**Expected Senior-Level Answer:**

1. **Detect**: Application becomes unresponsive; health checks fail; thread count plateaus but no throughput
2. **Capture** multiple thread dumps:
   ```bash
   jstack -l <pid> > thread_dump_1.txt
   sleep 5
   jstack -l <pid> > thread_dump_2.txt
   sleep 5
   jstack -l <pid> > thread_dump_3.txt
   # Compare: same threads BLOCKED in all 3 → deadlock
   ```
3. **Identify** in thread dump:
   ```
   "Thread-1" BLOCKED on 0x00000007abc12340 (a java.lang.Object)
     waiting to lock 0x00000007abc12360 (a java.lang.Object)
     locked 0x00000007abc12340
   
   "Thread-2" BLOCKED on 0x00000007abc12360 (a java.lang.Object)
     waiting to lock 0x00000007abc12340 (a java.lang.Object)
     locked 0x00000007abc12360
   
   Found one Java-level deadlock:  ← JVM auto-detects and reports
   ```
4. **Analyze**: determine lock ordering violation
5. **Fix patterns**:
   - Consistent lock ordering (always acquire locks in same order)
   - Use `tryLock(timeout)` with `ReentrantLock` (allows timeout-based recovery)
   - Reduce lock granularity (lock striping, concurrent collections)
   - Consider lock-free algorithms (`ConcurrentHashMap`, `AtomicReference` + CAS)

**Common Mistakes:**
- Only taking one thread dump (can't distinguish deadlock from slow I/O)
- Not using `-l` flag with jstack (misses ownable synchronizer info like `ReentrantLock`)

**Follow-Up Probing Questions:**
- How would you detect a livelock (threads running but not progressing)?
- Can Java detect deadlocks involving `ReentrantLock`? (Yes — `ThreadMXBean.findDeadlockedThreads()`)
- What about distributed deadlocks across microservice calls?

---

### 3C. System Design & Architecture (Java Focused)

---

#### Q11: How would you design a Java REST API to handle 100k+ concurrent requests?

**Expected Senior-Level Answer:**

**Architecture layers:**

```
                    ┌────────────────┐
                    │  Load Balancer │  (AWS ALB, Nginx)
                    └───────┬────────┘
                            │
              ┌─────────────┼─────────────┐
              │             │             │
        ┌─────▼────┐ ┌─────▼────┐ ┌─────▼────┐
        │  Java    │ │  Java    │ │  Java    │  Stateless instances
        │  App #1  │ │  App #2  │ │  App #N  │  (auto-scaled)
        └─────┬────┘ └─────┬────┘ └─────┬────┘
              │             │             │
        ┌─────▼─────────────▼─────────────▼────┐
        │          Connection Pool              │  HikariCP
        │          (bounded, per-instance)       │
        └────────────────┬─────────────────────┘
                         │
              ┌──────────▼──────────┐
              │    Database / Cache │  PostgreSQL + Redis
              └─────────────────────┘
```

**Key design decisions:**

1. **Non-blocking I/O**: Use Spring WebFlux (Reactor/Netty) or virtual threads (Java 21+)
   - Traditional thread-per-request with 100k requests → 100k threads → OOM
   - WebFlux: event loop with small thread pool handles all connections
   - Virtual threads: millions of lightweight threads, blocking I/O is OK

2. **Connection pooling**: HikariCP with bounded pool (50-200 connections per instance)
   - Database can't handle 100k concurrent connections anyway

3. **Caching**:
   - L1: In-process cache (Caffeine) — 10,000 entries, 50ms TTL for hot data
   - L2: Distributed cache (Redis) — centralized, shared across instances
   - Cache-aside pattern: check cache → miss → DB → populate cache

4. **Backpressure**:
   - Bounded request queues (Tomcat: `server.tomcat.max-threads`, `accept-count`)
   - Rate limiting per client (Resilience4j, bucket4j)
   - Circuit breaker for downstream services

5. **GC optimization**: 
   - G1 or ZGC (generational) for low-latency requirements
   - Size heap appropriately for per-request allocation rate
   - Minimize per-request object creation (reuse serializers, avoid unnecessary copies)

6. **Stateless design**: No HTTP session state; use JWT or centralized session store

**Common Mistakes:**
- Using `synchronized` or `ReentrantLock` extensively in request handlers (serializes throughput)
- Not implementing backpressure (unbounded queues → OOM under spike)
- Choosing thread-per-request Tomcat and opening 100k threads

**Follow-Up Probing Questions:**
- How would you benchmark the p99 latency at 100k RPS?
- What metrics would you monitor to detect degradation before failure?
- How would you handle graceful shutdown without dropping in-flight requests?

---

#### Q12: Design a caching strategy for a high-traffic Spring Boot service. Address consistency, eviction, and memory pressure.

**Expected Senior-Level Answer:**

```
Multi-Level Caching Architecture:

  Request → ┌──────────────┐  HIT   ┌──────────────────┐
            │ L1: Caffeine ├───────►│ Return immediately │
            │ (in-process) │        └──────────────────┘
            └──────┬───────┘
                   │ MISS
                   ▼
            ┌──────────────┐  HIT   ┌──────────────────┐
            │ L2: Redis    ├───────►│ Populate L1,      │
            │ (distributed)│        │ return             │
            └──────┬───────┘        └──────────────────┘
                   │ MISS
                   ▼
            ┌──────────────┐
            │ Database     ├───────► Populate L2 + L1, return
            └──────────────┘
```

**Implementation:**
```java
@Configuration
public class CacheConfig {
    @Bean
    public CacheManager caffeineCacheManager() {
        CaffeineCacheManager manager = new CaffeineCacheManager();
        manager.setCaffeine(Caffeine.newBuilder()
            .maximumSize(10_000)           // Bounded — prevents OOM
            .expireAfterWrite(5, TimeUnit.MINUTES) // TTL for consistency
            .recordStats());                // Enable metrics
        return manager;
    }
}

// Cache-aside with read-through:
@Cacheable(value = "users", key = "#userId", unless = "#result == null")
public User getUser(String userId) {
    return userRepository.findById(userId).orElse(null);
}

// Cache invalidation on write:
@CacheEvict(value = "users", key = "#user.id")
public User updateUser(User user) {
    return userRepository.save(user);
}
```

**Consistency strategies:**
| Strategy | Consistency | Performance | Complexity |
|----------|------------|-------------|------------|
| TTL-based expiry | Eventual (bounded staleness) | Best | Low |
| Write-through | Strong (write cache + DB atomically) | Good | Medium |
| Write-behind (async) | Eventual (async DB update) | Best write perf | High |
| Cache invalidation events (Pub/Sub) | Near real-time across instances | Good | Medium |

**Memory pressure management:**
- Caffeine: weight-based eviction (`maximumWeight`) — evict based on actual byte size, not just entry count
- Monitor cache hit ratio: < 80% → cache may be too small or TTL too short
- GC impact: large in-process caches → long-lived objects in Old Gen → longer GC pauses
- Solution: use off-heap cache (EhCache with off-heap tier, or Redis entirely) for >1GB cached data

**Common Mistakes:**
- Unbounded in-process cache → OOM under diverse key space
- Cache without TTL → stale data indefinitely
- Not invalidating across instances (each instance sees different data)

**Follow-Up Probing Questions:**
- How do you handle cache stampede (thundering herd) on cache miss?
- How would you implement distributed cache invalidation across 50 instances?
- What's the GC impact of caching 1M objects in-process?

---

#### Q13: How does GC behavior impact microservice latency, and how would you mitigate it?

**Expected Senior-Level Answer:**

**Impact:**
- GC STW pauses directly add to request latency (p99, p99.9)
- A 200ms G1 pause means some requests get +200ms latency
- In a microservice chain (A → B → C), GC pauses **compound**: each hop can independently add GC latency
- Kubernetes health check during GC pause → container killed → cascading failure

**Mitigation strategies:**

1. **Choose the right collector:**
   - Latency-critical: ZGC (sub-ms pauses) or Shenandoah
   - Throughput-critical (batch): Parallel GC
   - General purpose: G1 with tuned pause target

2. **Right-size the heap:**
   - Too small → frequent GC; too large → longer pauses (more to scan)
   - Rule of thumb: live data set × 3-4 for G1
   - Profile under realistic load (not just functional tests)

3. **Reduce allocation rate:**
   - Profile with JFR allocation profiling → identify top allocating methods
   - Reduce unnecessary object creation in hot paths
   - Use object pooling for expensive objects (but increases Old Gen pressure if misused)

4. **Container considerations:**
   ```yaml
   # Kubernetes: generous liveness probe timeouts
   livenessProbe:
     httpGet:
       path: /health
     initialDelaySeconds: 30
     periodSeconds: 10
     timeoutSeconds: 5      # Must be > max expected GC pause
     failureThreshold: 3    # Allow multiple failures before kill
   ```

5. **Monitor and alert:**
   - GC pause p99 > SLA threshold → alert
   - Full GC count > 0 in last hour → investigate
   - Allocation rate trending up → potential leak

```
GC Impact in Microservice Chain:

  Client → Service A → Service B → Service C → DB
           (20ms)      (15ms)      (10ms)
           
  Without GC: total = 45ms
  
  With GC: Service B has 200ms GC pause
           total = 20ms + (15ms + 200ms) + 10ms = 245ms
           
  p99 latency → driven by worst-case GC across ALL services in chain
```

**Common Mistakes:**
- Setting very large heaps (32GB+) without considering GC pause implications
- Not accounting for GC pauses in SLA calculations
- Using `-XX:+DisableExplicitGC` without checking if libraries legitimately need `System.gc()` (NIO direct buffer cleanup)

**Follow-Up Probing Questions:**
- How do you correlate GC pauses with request latency spikes in production?
- What's the impact of running a Java microservice in a container with swap enabled?

---

#### Q14: Design a rate limiter for a Java microservice. Address distributed scenarios.

**Expected Senior-Level Answer:**

**Single-instance (in-memory):**
```java
// Token Bucket Algorithm using Caffeine + AtomicLong
public class RateLimiter {
    private final ConcurrentHashMap<String, TokenBucket> buckets = new ConcurrentHashMap<>();
    
    public boolean allowRequest(String clientId) {
        TokenBucket bucket = buckets.computeIfAbsent(clientId,
            k -> new TokenBucket(100, 10)); // 100 max tokens, 10 tokens/sec refill
        return bucket.tryConsume();
    }
}

class TokenBucket {
    private final long maxTokens;
    private final double refillRate; // tokens per nanosecond
    private long availableTokens;
    private long lastRefillTimestamp;
    
    synchronized boolean tryConsume() {
        refill();
        if (availableTokens > 0) {
            availableTokens--;
            return true;
        }
        return false;
    }
    
    private void refill() {
        long now = System.nanoTime();
        long tokensToAdd = (long)((now - lastRefillTimestamp) * refillRate);
        availableTokens = Math.min(maxTokens, availableTokens + tokensToAdd);
        lastRefillTimestamp = now;
    }
}
```

**Distributed (Redis-based sliding window):**
```java
// Redis Lua script for atomic sliding window rate limiting
String luaScript = """
    local key = KEYS[1]
    local window = tonumber(ARGV[1])
    local limit = tonumber(ARGV[2])
    local now = tonumber(ARGV[3])
    
    -- Remove expired entries
    redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
    
    -- Count current window
    local count = redis.call('ZCARD', key)
    
    if count < limit then
        redis.call('ZADD', key, now, now .. '-' .. math.random())
        redis.call('EXPIRE', key, window / 1000)
        return 1  -- allowed
    end
    return 0  -- rejected
    """;

// Spring Boot integration:
@Component
public class DistributedRateLimiter {
    @Autowired private StringRedisTemplate redis;
    
    public boolean isAllowed(String clientId, int limit, int windowMs) {
        Long result = redis.execute(script, 
            List.of("rate:" + clientId),
            String.valueOf(windowMs),
            String.valueOf(limit),
            String.valueOf(System.currentTimeMillis()));
        return result != null && result == 1;
    }
}
```

**Common Mistakes:**
- Using `System.currentTimeMillis()` for precision timing (varies across hosts)
- Not handling Redis failures gracefully (fail-open or fail-closed decision)
- Implementing rate limiting per-instance instead of globally (each instance allows full rate)

**Follow-Up Probing Questions:**
- How would you implement rate limiting at the API gateway level vs application level?
- What happens when Redis is unavailable — fail-open (allow) or fail-closed (reject)?
- How would you handle rate limit burst allowance vs sustained rate?

---

#### Q15: How would you implement distributed locking in a Java microservices environment?

**Expected Senior-Level Answer:**

**Redis-based (Redisson):**
```java
// Redisson distributed lock — wraps Redis SET NX PX pattern
RLock lock = redissonClient.getLock("order:process:" + orderId);

try {
    // Wait up to 5s to acquire, auto-release after 30s (watchdog extends if still held)
    boolean acquired = lock.tryLock(5, 30, TimeUnit.SECONDS);
    if (acquired) {
        processOrder(orderId);
    } else {
        throw new LockAcquisitionException("Could not acquire lock for " + orderId);
    }
} finally {
    if (lock.isHeldByCurrentThread()) {
        lock.unlock();
    }
}
```

**Key considerations:**
- **Lease time / TTL**: Lock auto-expires to prevent deadlocks if holder crashes
- **Watchdog**: Redisson extends lock TTL while holder is alive
- **Fencing token**: Monotonically increasing token ensures stale lock holders can't perform writes after lock expires
- **Redis failover issues**: Lock on master → master fails → replica promotes → lock lost (addressed by RedLock algorithm, but controversial)

**Alternative: Database-based locking:**
```sql
-- Pessimistic lock via SELECT FOR UPDATE
SELECT * FROM orders WHERE id = ? FOR UPDATE NOWAIT;
-- Simple, ACID, no external dependency, but limits scalability
```

**Alternative: ZooKeeper/etcd:**
- Stronger consistency guarantees (CP system)
- Ephemeral nodes for automatic lock release
- Higher latency than Redis

**Common Mistakes:**
- Not handling lock expiration (process takes longer than TTL → two holders)
- Using distributed locks where an idempotent design would eliminate the need
- Not implementing fencing tokens for correctness

**Follow-Up Probing Questions:**
- What are the limitations of the RedLock algorithm?
- When should you prefer database-level locking over distributed locks?
- How would you monitor lock contention and hold times in production?

---

#### Q16: Describe your approach to observability and profiling for a Java microservice in production.

**Expected Senior-Level Answer:**

**Three Pillars:**

```
Observability Stack:

  ┌──────────┐  ┌──────────┐  ┌──────────┐
  │ Metrics  │  │  Logs    │  │  Traces  │
  │(Micrometer│  │(Structured│  │(OpenTelemetry│
  │→Prometheus│  │ JSON,    │  │ → Jaeger/   │
  │→Grafana) │  │ ELK/Loki)│  │   Tempo)    │
  └──────────┘  └──────────┘  └──────────┘
       │              │              │
       └──────────────┼──────────────┘
                      ▼
              ┌──────────────┐
              │  Correlation │  (trace-id in all metrics + logs)
              └──────────────┘
```

**Metrics (Micrometer + Prometheus):**
```java
@Component
public class OrderMetrics {
    private final Counter orderCreated;
    private final Timer orderProcessingTime;
    private final Gauge activeOrders;
    
    public OrderMetrics(MeterRegistry registry) {
        orderCreated = Counter.builder("orders.created")
            .tag("status", "success")
            .register(registry);
        orderProcessingTime = Timer.builder("orders.processing.time")
            .publishPercentiles(0.5, 0.95, 0.99, 0.999) // p50, p95, p99, p99.9
            .register(registry);
    }
}

// JVM metrics auto-registered by Micrometer:
// jvm_gc_pause_seconds (histogram by GC cause and action)
// jvm_memory_used_bytes (by area: heap, non-heap)
// jvm_threads_states (by state: runnable, blocked, waiting)
// jvm_buffer_memory_used_bytes (direct buffers)
```

**Profiling in production:**
```bash
# Java Flight Recorder — always-on, near-zero overhead
-XX:StartFlightRecording=disk=true,maxage=6h,maxsize=256m,dumponexit=true

# async-profiler — CPU + allocation profiling with minimal overhead
./asprof -d 30 -f profile.html <pid>     # 30-second CPU profile
./asprof -e alloc -d 30 -f alloc.html <pid>  # Allocation profile

# jcmd for on-demand diagnostics
jcmd <pid> JFR.start duration=60s filename=recording.jfr
jcmd <pid> Thread.print                    # Thread dump
jcmd <pid> GC.heap_info                    # Heap summary
```

**Common Mistakes:**
- Not including trace-id in logs (impossible to correlate across services)
- Only monitoring averages (masks p99 issues)
- Not profiling allocation rate (often the root cause of GC issues)

**Follow-Up Probing Questions:**
- How would you set up alerts for GC pause times exceeding your SLA?
- What's the overhead of running JFR continuously in production?
- How would you trace a slow request across 5 microservices?

---

## Quick Reference — JVM Flags Cheatsheet

```bash
# ════════════════════════════════════════════════════════════════
# PRODUCTION JVM FLAGS TEMPLATE (Java 21+, Container, G1 GC)
# ════════════════════════════════════════════════════════════════

java \
  # ── Memory ──
  -XX:MaxRAMPercentage=75.0 \
  -XX:InitialRAMPercentage=75.0 \
  -XX:MaxMetaspaceSize=512m \
  -Xss512k \
  
  # ── GC ──
  -XX:+UseG1GC \
  -XX:MaxGCPauseMillis=150 \
  -XX:+ParallelRefProcEnabled \
  -XX:+UseStringDeduplication \
  
  # ── Resilience ──
  -XX:+HeapDumpOnOutOfMemoryError \
  -XX:HeapDumpPath=/var/dumps/ \
  -XX:+ExitOnOutOfMemoryError \
  
  # ── Observability ──
  -Xlog:gc*:file=/var/log/gc.log:time,uptime,level,tags:filecount=5,filesize=20m \
  -XX:+FlightRecorder \
  -XX:StartFlightRecording=disk=true,maxage=6h,maxsize=256m,dumponexit=true \
  -XX:NativeMemoryTracking=summary \
  
  -jar my-service.jar

# ════════════════════════════════════════════════════════════════
# ULTRA-LOW LATENCY TEMPLATE (Java 21+, ZGC)
# ════════════════════════════════════════════════════════════════

java \
  -XX:+UseZGC -XX:+ZGenerational \
  -Xms8g -Xmx8g \
  -XX:SoftMaxHeapSize=6g \
  -XX:+HeapDumpOnOutOfMemoryError \
  -XX:+ExitOnOutOfMemoryError \
  -Xlog:gc*:file=/var/log/gc.log:time,uptime,level,tags:filecount=5,filesize=20m \
  -XX:+FlightRecorder \
  -jar my-trading-engine.jar
```

---

> **Document Version:** 1.0 | **Last Updated:** February 2026
> **Author:** Interview Preparation Guide for Staff/Senior Java Engineers
