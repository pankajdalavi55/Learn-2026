# Java Multithreading Complete Guide

## Table of Contents
- [1. Introduction to Multithreading](#1-introduction-to-multithreading)
  - [1.1 What is a Process?](#11-what-is-a-process)
  - [1.2 What is a Thread?](#12-what-is-a-thread)
  - [1.3 Difference Between Process and Thread](#13-difference-between-process-and-thread)
  - [1.4 Why Multithreading is Important](#14-why-multithreading-is-important)
  - [1.5 Real-World Use Cases](#15-real-world-use-cases)
- [2. Java Thread Basics](#2-java-thread-basics)
  - [2.1 Thread Lifecycle](#21-thread-lifecycle)
  - [2.2 Creating Threads](#22-creating-threads)
  - [2.3 Thread Methods](#23-thread-methods)
- [3. Thread Synchronization](#3-thread-synchronization)
  - [3.1 Race Conditions](#31-race-conditions)
  - [3.2 Critical Section](#32-critical-section)
  - [3.3 synchronized Keyword](#33-synchronized-keyword)
  - [3.4 Object-Level vs Class-Level Locking](#34-object-level-vs-class-level-locking)
  - [3.5 synchronized Blocks](#35-synchronized-blocks)
  - [3.6 Deadlock](#36-deadlock)
  - [3.7 Livelock and Starvation](#37-livelock-and-starvation)
- [4. Inter-Thread Communication](#4-inter-thread-communication)
  - [4.1 wait() Method](#41-wait-method)
  - [4.2 notify() Method](#42-notify-method)
  - [4.3 notifyAll() Method](#43-notifyall-method)
  - [4.4 Producer-Consumer Problem](#44-producer-consumer-problem)
- [5. Advanced Concurrency (java.util.concurrent)](#5-advanced-concurrency-javautilconcurrent)
  - [5.1 ExecutorService](#51-executorservice)
  - [5.2 Thread Pools](#52-thread-pools)
  - [5.3 Callable vs Runnable](#53-callable-vs-runnable)
  - [5.4 Future](#54-future)
  - [5.5 CountDownLatch](#55-countdownlatch)
  - [5.6 CyclicBarrier](#56-cyclicbarrier)
  - [5.7 Semaphore](#57-semaphore)
  - [5.8 ReentrantLock](#58-reentrantlock)
  - [5.9 ReadWriteLock](#59-readwritelock)
  - [5.10 ConcurrentHashMap](#510-concurrenthashmap)
- [6. Atomic Variables and Volatile](#6-atomic-variables-and-volatile)
  - [6.1 volatile Keyword](#61-volatile-keyword)
  - [6.2 AtomicInteger](#62-atomicinteger)
  - [6.3 Compare and Swap (CAS)](#63-compare-and-swap-cas)
- [7. Fork/Join Framework](#7-forkjoin-framework)
  - [7.1 Work-Stealing Algorithm](#71-work-stealing-algorithm)
  - [7.2 RecursiveTask and RecursiveAction](#72-recursivetask-and-recursiveaction)
- [8. CompletableFuture](#8-completablefuture)
  - [8.1 Asynchronous Programming](#81-asynchronous-programming)
  - [8.2 Chaining Tasks](#82-chaining-tasks)
  - [8.3 Exception Handling](#83-exception-handling)
- [9. Best Practices](#9-best-practices)
  - [9.1 Avoiding Deadlocks](#91-avoiding-deadlocks)
  - [9.2 Minimizing Synchronization](#92-minimizing-synchronization)
  - [9.3 Thread Safety Strategies](#93-thread-safety-strategies)
  - [9.4 Immutable Objects](#94-immutable-objects)
  - [9.5 When NOT to Use Multithreading](#95-when-not-to-use-multithreading)
- [10. Practical Projects](#10-practical-projects)
  - [10.1 Multithreaded File Downloader](#101-multithreaded-file-downloader)
  - [10.2 Parallel Data Processing](#102-parallel-data-processing)
  - [10.3 Simple Web Server Simulation](#103-simple-web-server-simulation)
  - [10.4 Producer-Consumer using BlockingQueue](#104-producer-consumer-using-blockingqueue)

---

## 1. Introduction to Multithreading

### 1.1 What is a Process?

A **process** is an independent program in execution. It is a self-contained execution environment with its own memory space, resources, and state.

#### Key Characteristics of a Process:
- **Isolation**: Each process runs in its own memory space
- **Resources**: Has its own set of resources (memory, file handles, etc.)
- **Independent**: Processes are independent of each other
- **Heavyweight**: Creating and managing processes is resource-intensive
- **Inter-Process Communication (IPC)**: Processes communicate through IPC mechanisms (pipes, sockets, shared memory)

```
┌─────────────────────────────────────────────────────────────┐
│                         PROCESS                              │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐       │
│  │     Code     │  │     Data     │  │     Heap     │       │
│  │   Segment    │  │   Segment    │  │              │       │
│  └──────────────┘  └──────────────┘  └──────────────┘       │
│                                                              │
│  ┌─────────────────────────────────────────────────────┐    │
│  │                      Stack                           │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                              │
│  Process ID (PID) | Program Counter | Registers | State     │
└─────────────────────────────────────────────────────────────┘
```

#### Example: Multiple Processes
When you open Chrome, Word, and Spotify simultaneously, each application runs as a separate process with its own memory space.

---

### 1.2 What is a Thread?

A **thread** is the smallest unit of execution within a process. It is also called a **lightweight process** because it shares the process's resources but has its own execution path.

#### Key Characteristics of a Thread:
- **Shared Memory**: Threads within a process share the same memory space (heap, code, data)
- **Own Stack**: Each thread has its own stack and program counter
- **Lightweight**: Creating threads is less resource-intensive than processes
- **Concurrent Execution**: Multiple threads can execute concurrently
- **Communication**: Threads communicate through shared memory (faster than IPC)

```
┌─────────────────────────────────────────────────────────────────────┐
│                              PROCESS                                 │
├─────────────────────────────────────────────────────────────────────┤
│      SHARED RESOURCES                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────────┐   │
│  │     Code     │  │     Data     │  │         Heap             │   │
│  │   Segment    │  │   Segment    │  │   (Shared Objects)       │   │
│  └──────────────┘  └──────────────┘  └──────────────────────────┘   │
│                                                                      │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │
│  │  Thread 1   │  │  Thread 2   │  │  Thread 3   │                  │
│  ├─────────────┤  ├─────────────┤  ├─────────────┤                  │
│  │ Own Stack   │  │ Own Stack   │  │ Own Stack   │                  │
│  │ Own PC      │  │ Own PC      │  │ Own PC      │                  │
│  │ Own Regs    │  │ Own Regs    │  │ Own Regs    │                  │
│  └─────────────┘  └─────────────┘  └─────────────┘                  │
└─────────────────────────────────────────────────────────────────────┘

PC = Program Counter, Regs = Registers
```

#### Types of Threads:
1. **User Threads**: Managed by the application (JVM in Java's case)
2. **Kernel Threads**: Managed by the operating system
3. **Daemon Threads**: Background threads that don't prevent JVM termination (e.g., Garbage Collector)

---

### 1.3 Difference Between Process and Thread

| Aspect | Process | Thread |
|--------|---------|--------|
| **Definition** | Independent program in execution | Smallest unit of execution within a process |
| **Memory** | Own memory space | Shares memory with other threads in same process |
| **Creation Cost** | High (heavyweight) | Low (lightweight) |
| **Communication** | IPC mechanisms (slow) | Shared memory (fast) |
| **Context Switch** | Expensive | Cheaper |
| **Isolation** | Fully isolated | Not isolated (shares heap) |
| **Failure Impact** | One process crash doesn't affect others | Thread crash can crash the entire process |
| **Resource Overhead** | High | Low |
| **Creation Time** | More time | Less time |
| **Data Sharing** | Requires explicit IPC | Direct access to shared data |

```
┌─────────────────────────────────────────────────────────────────────┐
│                        OPERATING SYSTEM                              │
├──────────────────────┬──────────────────────┬───────────────────────┤
│      PROCESS 1       │      PROCESS 2       │      PROCESS 3        │
│  ┌────────────────┐  │  ┌────────────────┐  │  ┌────────────────┐   │
│  │ Memory Space A │  │  │ Memory Space B │  │  │ Memory Space C │   │
│  │                │  │  │                │  │  │                │   │
│  │ ┌───┐ ┌───┐   │  │  │ ┌───┐ ┌───┐   │  │  │    ┌───┐       │   │
│  │ │T1 │ │T2 │   │  │  │ │T1 │ │T2 │   │  │  │    │T1 │       │   │
│  │ └───┘ └───┘   │  │  │ └───┘ └───┘   │  │  │    └───┘       │   │
│  │    ┌───┐      │  │  │               │  │  │                │   │
│  │    │T3 │      │  │  │               │  │  │                │   │
│  │    └───┘      │  │  │               │  │  │                │   │
│  └────────────────┘  │  └────────────────┘  │  └────────────────┘   │
│   (3 threads)        │   (2 threads)        │   (1 thread)          │
└──────────────────────┴──────────────────────┴───────────────────────┘

T = Thread, Each process has isolated memory
```

---

### 1.4 Why Multithreading is Important

#### 1.4.1 Benefits of Multithreading

**1. Better Resource Utilization**
```
Single-Threaded:
┌─────┬─────┬─────┬─────┬─────┬─────┬─────┬─────┐
│ CPU │ I/O │WAIT │ CPU │ I/O │WAIT │ CPU │ ... │  (CPU often idle)
└─────┴─────┴─────┴─────┴─────┴─────┴─────┴─────┘

Multi-Threaded:
Thread 1: ┌─────┬─────┬─────┬─────────────────────┐
          │ CPU │ I/O │WAIT │       CPU           │
          └─────┴─────┴─────┴─────────────────────┘
Thread 2:       ┌─────┬─────┬─────────────────────┐
                │ CPU │ I/O │       CPU           │  (CPU utilized better)
                └─────┴─────┴─────────────────────┘
```

**2. Improved Responsiveness**
```java
// Without multithreading - UI freezes during heavy operation
public class SingleThreadedApp {
    public void handleButtonClick() {
        downloadLargeFile();  // UI freezes for minutes
        processData();        // User cannot interact
        updateUI();
    }
}

// With multithreading - UI stays responsive
public class MultiThreadedApp {
    public void handleButtonClick() {
        new Thread(() -> {
            downloadLargeFile();  // Runs in background
            processData();
            SwingUtilities.invokeLater(() -> updateUI());
        }).start();
        // UI remains responsive, user can continue interacting
    }
}
```

**3. Parallel Processing on Multi-Core CPUs**
```
Single-Threaded on 4-Core CPU:
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│ Core 1 │ │ Core 2 │ │ Core 3 │ │ Core 4 │
│  BUSY  │ │  IDLE  │ │  IDLE  │ │  IDLE  │
└────────┘ └────────┘ └────────┘ └────────┘
           CPU Utilization: 25%

Multi-Threaded on 4-Core CPU:
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│ Core 1 │ │ Core 2 │ │ Core 3 │ │ Core 4 │
│Thread 1│ │Thread 2│ │Thread 3│ │Thread 4│
│  BUSY  │ │  BUSY  │ │  BUSY  │ │  BUSY  │
└────────┘ └────────┘ └────────┘ └────────┘
           CPU Utilization: 100%
```

**4. Simplified Modeling of Real-World Scenarios**
Many real-world problems are naturally parallel:
- Web server handling multiple client requests
- Chat application managing multiple conversations
- Game engine updating physics, AI, and rendering simultaneously

**5. Cost Efficiency**
- Threads share process resources, reducing memory overhead
- Faster context switching compared to processes
- More scalable for concurrent operations

#### 1.4.2 Challenges of Multithreading

| Challenge | Description |
|-----------|-------------|
| **Race Conditions** | Multiple threads accessing shared data simultaneously |
| **Deadlocks** | Threads waiting for each other indefinitely |
| **Thread Starvation** | Some threads never get CPU time |
| **Complexity** | Harder to debug and test |
| **Non-Determinism** | Execution order varies between runs |

---

### 1.5 Real-World Use Cases

#### 1.5.1 Web Servers

```
                    ┌─────────────────────────────────────┐
                    │          WEB SERVER                  │
                    │    (Apache Tomcat / Jetty)           │
                    │                                      │
   Client 1 ───────►│  ┌─────────────────────────────┐    │
                    │  │      Thread Pool             │    │
   Client 2 ───────►│  │  ┌───┐ ┌───┐ ┌───┐ ┌───┐   │    │
                    │  │  │T1 │ │T2 │ │T3 │ │T4 │   │    │
   Client 3 ───────►│  │  └───┘ └───┘ └───┘ └───┘   │    │
                    │  │         ...                 │    │
   Client N ───────►│  │  ┌───┐ ┌───┐              │    │
                    │  │  │Tn │ │Tn+1│              │    │
                    │  │  └───┘ └───┘               │    │
                    │  └─────────────────────────────┘    │
                    └─────────────────────────────────────┘

Each client request is handled by a separate thread from the pool
```

```java
// Simplified Web Server Example
public class SimpleWebServer {
    private final ExecutorService threadPool = Executors.newFixedThreadPool(100);
    
    public void start(int port) throws IOException {
        ServerSocket serverSocket = new ServerSocket(port);
        System.out.println("Server started on port " + port);
        
        while (true) {
            Socket clientSocket = serverSocket.accept();
            // Each client request handled by separate thread
            threadPool.execute(new ClientHandler(clientSocket));
        }
    }
}

class ClientHandler implements Runnable {
    private final Socket clientSocket;
    
    public ClientHandler(Socket socket) {
        this.clientSocket = socket;
    }
    
    @Override
    public void run() {
        try {
            // Handle HTTP request
            BufferedReader in = new BufferedReader(
                new InputStreamReader(clientSocket.getInputStream()));
            PrintWriter out = new PrintWriter(clientSocket.getOutputStream(), true);
            
            String request = in.readLine();
            System.out.println("Thread " + Thread.currentThread().getName() + 
                             " handling: " + request);
            
            // Process request and send response
            out.println("HTTP/1.1 200 OK");
            out.println("Content-Type: text/html");
            out.println();
            out.println("<h1>Hello from Thread " + Thread.currentThread().getName() + "</h1>");
            
            clientSocket.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
```

#### 1.5.2 Database Connection Pooling

```
┌─────────────────────────────────────────────────────────────┐
│                    APPLICATION SERVER                        │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │                   Connection Pool                       │ │
│  │  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐      │ │
│  │  │Conn1│ │Conn2│ │Conn3│ │Conn4│ │Conn5│ │Conn6│      │ │
│  │  └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘ └──┬──┘      │ │
│  └─────┼───────┼───────┼───────┼───────┼───────┼──────────┘ │
└────────┼───────┼───────┼───────┼───────┼───────┼────────────┘
         │       │       │       │       │       │
         └───────┴───────┴───┬───┴───────┴───────┘
                             │
                    ┌────────▼────────┐
                    │    DATABASE     │
                    └─────────────────┘

Multiple threads reuse connections from the pool instead of 
creating new connections for each request
```

#### 1.5.3 File Download Manager

```java
public class FileDownloadManager {
    private final ExecutorService downloadPool = Executors.newFixedThreadPool(5);
    
    public void downloadFiles(List<String> urls) {
        List<Future<String>> futures = new ArrayList<>();
        
        for (String url : urls) {
            Future<String> future = downloadPool.submit(() -> {
                System.out.println(Thread.currentThread().getName() + 
                                 " downloading: " + url);
                // Simulate download
                Thread.sleep(2000);
                return "Downloaded: " + url;
            });
            futures.add(future);
        }
        
        // Wait for all downloads to complete
        for (Future<String> future : futures) {
            try {
                System.out.println(future.get());
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        
        downloadPool.shutdown();
    }
    
    public static void main(String[] args) {
        FileDownloadManager manager = new FileDownloadManager();
        List<String> urls = Arrays.asList(
            "http://example.com/file1.pdf",
            "http://example.com/file2.pdf",
            "http://example.com/file3.pdf",
            "http://example.com/file4.pdf",
            "http://example.com/file5.pdf"
        );
        
        long start = System.currentTimeMillis();
        manager.downloadFiles(urls);
        long end = System.currentTimeMillis();
        
        // Sequential: 5 files × 2 seconds = 10 seconds
        // Parallel (5 threads): ~2 seconds
        System.out.println("Total time: " + (end - start) + "ms");
    }
}
```

#### 1.5.4 Real-Time Data Processing

```
┌─────────────────────────────────────────────────────────────────┐
│              REAL-TIME STOCK TRADING SYSTEM                      │
│                                                                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │ Data Feed   │    │  Strategy   │    │   Order     │          │
│  │   Thread    │───►│   Thread    │───►│   Thread    │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
│         │                  │                  │                  │
│         ▼                  ▼                  ▼                  │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐          │
│  │   Market    │    │    Risk     │    │  Logging    │          │
│  │   Data      │    │  Management │    │   Thread    │          │
│  │   Queue     │    │   Thread    │    │             │          │
│  └─────────────┘    └─────────────┘    └─────────────┘          │
│                                                                  │
│  Each component runs independently, communicating via queues     │
└─────────────────────────────────────────────────────────────────┘
```

#### 1.5.5 Video Game Engine

```
┌──────────────────────────────────────────────────────────────┐
│                    GAME ENGINE ARCHITECTURE                   │
│                                                               │
│  ┌──────────────────────────────────────────────────────┐    │
│  │                    MAIN GAME LOOP                     │    │
│  └──────────────────────────────────────────────────────┘    │
│           │              │              │              │      │
│           ▼              ▼              ▼              ▼      │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌────────┐ │
│  │   Input     │ │   Physics   │ │     AI      │ │ Audio  │ │
│  │   Thread    │ │   Thread    │ │   Thread    │ │ Thread │ │
│  │             │ │             │ │             │ │        │ │
│  │ - Keyboard  │ │ - Collision │ │ - Pathfind  │ │ - SFX  │ │
│  │ - Mouse     │ │ - Gravity   │ │ - Behavior  │ │ - Music│ │
│  │ - Gamepad   │ │ - Movement  │ │ - Decisions │ │        │ │
│  └─────────────┘ └─────────────┘ └─────────────┘ └────────┘ │
│           │              │              │              │      │
│           └──────────────┴──────┬───────┴──────────────┘      │
│                                 ▼                             │
│                    ┌─────────────────────┐                   │
│                    │   Render Thread     │                   │
│                    │   (60 FPS target)   │                   │
│                    └─────────────────────┘                   │
└──────────────────────────────────────────────────────────────┘
```

#### 1.5.6 Chat Application

```java
public class ChatServer {
    private final List<ClientHandler> clients = Collections.synchronizedList(new ArrayList<>());
    private final ExecutorService pool = Executors.newCachedThreadPool();
    
    public void start(int port) throws IOException {
        ServerSocket serverSocket = new ServerSocket(port);
        System.out.println("Chat Server started on port " + port);
        
        while (true) {
            Socket clientSocket = serverSocket.accept();
            ClientHandler handler = new ClientHandler(clientSocket, this);
            clients.add(handler);
            pool.execute(handler);
        }
    }
    
    public void broadcastMessage(String message, ClientHandler sender) {
        for (ClientHandler client : clients) {
            if (client != sender) {
                client.sendMessage(message);
            }
        }
    }
    
    public void removeClient(ClientHandler client) {
        clients.remove(client);
    }
}
```

---

## 2. Java Thread Basics

### 2.1 Thread Lifecycle

A thread in Java goes through various states during its lifetime. Understanding these states is crucial for writing correct multithreaded programs.

#### Thread States

```
                              ┌─────────────────┐
                              │      NEW        │
                              │ (Thread created │
                              │  but not started)│
                              └────────┬────────┘
                                       │ start()
                                       ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           RUNNABLE                                   │
│  ┌─────────────────────┐         ┌─────────────────────┐            │
│  │       READY         │ ◄─────► │      RUNNING        │            │
│  │ (Waiting for CPU)   │scheduler│ (Executing on CPU)  │            │
│  └─────────────────────┘         └─────────────────────┘            │
└────────────────────────────────────┬────────────────────────────────┘
              │                      │                    │
              │ wait()               │ sleep()/join()     │ waiting for
              │ wait(timeout)        │ yield()            │ lock
              ▼                      ▼                    ▼
     ┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
     │    WAITING      │   │  TIMED_WAITING  │   │    BLOCKED      │
     │                 │   │                 │   │                 │
     │ - wait()        │   │ - sleep(ms)     │   │ Waiting to      │
     │ - join()        │   │ - wait(ms)      │   │ enter           │
     │ - park()        │   │ - join(ms)      │   │ synchronized    │
     │                 │   │ - parkNanos()   │   │ block/method    │
     └────────┬────────┘   └────────┬────────┘   └────────┬────────┘
              │ notify()            │ timeout             │ lock acquired
              │ notifyAll()         │ expires             │
              │ interrupt()         │                     │
              └─────────────────────┴─────────────────────┘
                                    │
                                    ▼
                         ┌─────────────────┐
                         │   TERMINATED    │
                         │                 │
                         │ - run() returns │
                         │ - Exception     │
                         │   thrown        │
                         └─────────────────┘
```

#### Java Thread.State Enum

```java
public class ThreadStateDemo {
    public static void main(String[] args) throws InterruptedException {
        // 1. NEW - Thread created but not started
        Thread newThread = new Thread(() -> {});
        System.out.println("New Thread State: " + newThread.getState()); // NEW
        
        // 2. RUNNABLE - Thread is executing or ready to execute
        Thread runnableThread = new Thread(() -> {
            while (true) {} // Keep running
        });
        runnableThread.start();
        Thread.sleep(100);
        System.out.println("Runnable Thread State: " + runnableThread.getState()); // RUNNABLE
        
        // 3. BLOCKED - Waiting for monitor lock
        Object lock = new Object();
        Thread blockedThread = new Thread(() -> {
            synchronized (lock) {
                // Won't reach here until main releases lock
            }
        });
        synchronized (lock) {
            blockedThread.start();
            Thread.sleep(100);
            System.out.println("Blocked Thread State: " + blockedThread.getState()); // BLOCKED
        }
        
        // 4. WAITING - Waiting indefinitely
        Thread waitingThread = new Thread(() -> {
            synchronized (lock) {
                try {
                    lock.wait();
                } catch (InterruptedException e) {}
            }
        });
        waitingThread.start();
        Thread.sleep(100);
        System.out.println("Waiting Thread State: " + waitingThread.getState()); // WAITING
        
        // 5. TIMED_WAITING - Waiting for specific time
        Thread timedWaitingThread = new Thread(() -> {
            try {
                Thread.sleep(5000);
            } catch (InterruptedException e) {}
        });
        timedWaitingThread.start();
        Thread.sleep(100);
        System.out.println("Timed Waiting Thread State: " + timedWaitingThread.getState()); // TIMED_WAITING
        
        // 6. TERMINATED - Thread has completed execution
        Thread terminatedThread = new Thread(() -> 
            System.out.println("Thread executing"));
        terminatedThread.start();
        terminatedThread.join();
        System.out.println("Terminated Thread State: " + terminatedThread.getState()); // TERMINATED
        
        // Cleanup
        runnableThread.interrupt();
        synchronized (lock) { lock.notifyAll(); }
    }
}
```

#### State Transition Summary Table

| Current State | Trigger | New State |
|---------------|---------|-----------|
| NEW | `start()` | RUNNABLE |
| RUNNABLE | `wait()`, `join()` | WAITING |
| RUNNABLE | `sleep(ms)`, `wait(ms)`, `join(ms)` | TIMED_WAITING |
| RUNNABLE | Waiting for synchronized block | BLOCKED |
| RUNNABLE | `run()` completes or exception | TERMINATED |
| WAITING | `notify()`, `notifyAll()`, `interrupt()` | RUNNABLE |
| TIMED_WAITING | Timeout expires or `interrupt()` | RUNNABLE |
| BLOCKED | Lock acquired | RUNNABLE |

---

### 2.2 Creating Threads

Java provides three main ways to create threads:

#### 2.2.1 Extending the Thread Class

The simplest way to create a thread is by extending the `Thread` class and overriding the `run()` method.

```java
// Method 1: Extending Thread class
class MyThread extends Thread {
    private String threadName;
    
    public MyThread(String name) {
        this.threadName = name;
    }
    
    @Override
    public void run() {
        for (int i = 1; i <= 5; i++) {
            System.out.println(threadName + " - Count: " + i);
            try {
                Thread.sleep(500); // Sleep for 500ms
            } catch (InterruptedException e) {
                System.out.println(threadName + " interrupted");
            }
        }
        System.out.println(threadName + " finished execution");
    }
}

public class ThreadExtendDemo {
    public static void main(String[] args) {
        System.out.println("Main thread started");
        
        // Create thread instances
        MyThread thread1 = new MyThread("Thread-1");
        MyThread thread2 = new MyThread("Thread-2");
        
        // Start threads (calls run() internally)
        thread1.start();
        thread2.start();
        
        // Wait for threads to complete
        try {
            thread1.join();
            thread2.join();
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        
        System.out.println("Main thread finished");
    }
}
```

**Output (interleaved execution):**
```
Main thread started
Thread-1 - Count: 1
Thread-2 - Count: 1
Thread-1 - Count: 2
Thread-2 - Count: 2
Thread-1 - Count: 3
Thread-2 - Count: 3
Thread-1 - Count: 4
Thread-2 - Count: 4
Thread-1 - Count: 5
Thread-2 - Count: 5
Thread-1 finished execution
Thread-2 finished execution
Main thread finished
```

**Pros and Cons:**

| Pros | Cons |
|------|------|
| Simple and straightforward | Cannot extend another class (Java single inheritance) |
| Direct access to Thread methods | Tight coupling between thread and task logic |
| Easy to understand for beginners | Less flexible for reuse |

---

#### 2.2.2 Implementing the Runnable Interface

The preferred approach is to implement the `Runnable` interface. This separates the task from the thread mechanism.

```java
// Method 2: Implementing Runnable interface
class MyTask implements Runnable {
    private String taskName;
    
    public MyTask(String name) {
        this.taskName = name;
    }
    
    @Override
    public void run() {
        for (int i = 1; i <= 5; i++) {
            System.out.println(taskName + " - Count: " + i + 
                             " [Thread: " + Thread.currentThread().getName() + "]");
            try {
                Thread.sleep(500);
            } catch (InterruptedException e) {
                System.out.println(taskName + " interrupted");
                return;
            }
        }
        System.out.println(taskName + " completed");
    }
}

public class RunnableDemo {
    public static void main(String[] args) {
        System.out.println("Main thread: " + Thread.currentThread().getName());
        
        // Create Runnable tasks
        Runnable task1 = new MyTask("Task-A");
        Runnable task2 = new MyTask("Task-B");
        
        // Create threads with Runnable tasks
        Thread thread1 = new Thread(task1, "Worker-1");
        Thread thread2 = new Thread(task2, "Worker-2");
        
        // Start execution
        thread1.start();
        thread2.start();
        
        // Wait for completion
        try {
            thread1.join();
            thread2.join();
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        
        System.out.println("All tasks completed");
    }
}
```

**Using Lambda Expressions (Java 8+):**

```java
public class RunnableLambdaDemo {
    public static void main(String[] args) {
        // Lambda expression - concise way to create Runnable
        Runnable task1 = () -> {
            for (int i = 0; i < 5; i++) {
                System.out.println("Task 1 - " + i);
            }
        };
        
        // Even more concise with method reference
        Runnable task2 = () -> System.out.println("Task 2 executed");
        
        // Anonymous class (traditional way)
        Thread thread1 = new Thread(new Runnable() {
            @Override
            public void run() {
                System.out.println("Anonymous Runnable executed");
            }
        });
        
        // Lambda (modern way)
        Thread thread2 = new Thread(() -> {
            System.out.println("Lambda Runnable executed");
        });
        
        // Inline lambda
        new Thread(() -> System.out.println("Inline Thread")).start();
        
        thread1.start();
        thread2.start();
    }
}
```

**Pros and Cons:**

| Pros | Cons |
|------|------|
| Class can extend another class | Cannot return a result |
| Separates task from thread | Cannot throw checked exceptions |
| More flexible and reusable | Need to create Thread to execute |
| Can share single Runnable among multiple threads | |
| Works with ExecutorService | |

---

#### 2.2.3 Using Callable and Future

`Callable` is similar to `Runnable` but can return a result and throw checked exceptions. It works with `ExecutorService` and returns a `Future` object.

```java
import java.util.concurrent.*;

// Method 3: Using Callable and Future
class CalculationTask implements Callable<Integer> {
    private int number;
    
    public CalculationTask(int number) {
        this.number = number;
    }
    
    @Override
    public Integer call() throws Exception {
        System.out.println(Thread.currentThread().getName() + 
                         " calculating factorial of " + number);
        
        // Simulate time-consuming calculation
        Thread.sleep(1000);
        
        int result = factorial(number);
        
        System.out.println(Thread.currentThread().getName() + 
                         " completed: " + number + "! = " + result);
        return result;
    }
    
    private int factorial(int n) {
        if (n <= 1) return 1;
        return n * factorial(n - 1);
    }
}

public class CallableDemo {
    public static void main(String[] args) {
        // Create ExecutorService with fixed thread pool
        ExecutorService executor = Executors.newFixedThreadPool(3);
        
        try {
            // Submit Callable tasks
            Future<Integer> future1 = executor.submit(new CalculationTask(5));
            Future<Integer> future2 = executor.submit(new CalculationTask(7));
            Future<Integer> future3 = executor.submit(new CalculationTask(10));
            
            // Check if tasks are done (non-blocking)
            System.out.println("Task 1 done? " + future1.isDone());
            System.out.println("Task 2 done? " + future2.isDone());
            System.out.println("Task 3 done? " + future3.isDone());
            
            // Get results (blocking until completion)
            System.out.println("\nWaiting for results...\n");
            
            Integer result1 = future1.get(); // Blocks until result available
            Integer result2 = future2.get();
            Integer result3 = future3.get();
            
            System.out.println("\n--- Results ---");
            System.out.println("5! = " + result1);
            System.out.println("7! = " + result2);
            System.out.println("10! = " + result3);
            
        } catch (InterruptedException | ExecutionException e) {
            e.printStackTrace();
        } finally {
            executor.shutdown();
        }
    }
}
```

**Future API Methods:**

```java
public class FutureMethodsDemo {
    public static void main(String[] args) {
        ExecutorService executor = Executors.newSingleThreadExecutor();
        
        Future<String> future = executor.submit(() -> {
            Thread.sleep(2000);
            return "Task completed!";
        });
        
        // isDone() - Check if task is completed
        System.out.println("Is done: " + future.isDone()); // false
        
        // isCancelled() - Check if task was cancelled
        System.out.println("Is cancelled: " + future.isCancelled()); // false
        
        // cancel(mayInterruptIfRunning) - Attempt to cancel task
        // future.cancel(true);
        
        try {
            // get() - Wait indefinitely for result
            // String result = future.get();
            
            // get(timeout, unit) - Wait with timeout
            String result = future.get(5, TimeUnit.SECONDS);
            System.out.println("Result: " + result);
            
        } catch (TimeoutException e) {
            System.out.println("Task took too long!");
            future.cancel(true);
        } catch (InterruptedException | ExecutionException e) {
            e.printStackTrace();
        }
        
        System.out.println("Is done: " + future.isDone()); // true
        
        executor.shutdown();
    }
}
```

**Callable with Exception Handling:**

```java
public class CallableExceptionDemo {
    public static void main(String[] args) {
        ExecutorService executor = Executors.newSingleThreadExecutor();
        
        // Callable that throws an exception
        Callable<Integer> riskyTask = () -> {
            int a = 10, b = 0;
            if (b == 0) {
                throw new ArithmeticException("Cannot divide by zero!");
            }
            return a / b;
        };
        
        Future<Integer> future = executor.submit(riskyTask);
        
        try {
            Integer result = future.get();
            System.out.println("Result: " + result);
        } catch (ExecutionException e) {
            // ExecutionException wraps the actual exception
            System.out.println("Task threw an exception: " + e.getCause().getMessage());
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
        
        executor.shutdown();
    }
}
```

**Comparison: Runnable vs Callable**

| Feature | Runnable | Callable |
|---------|----------|----------|
| Return value | `void` | Generic type `V` |
| Exception | Cannot throw checked exceptions | Can throw Exception |
| Method | `run()` | `call()` |
| Introduced | Java 1.0 | Java 5 |
| Use with | `Thread`, `ExecutorService` | `ExecutorService` only |
| Result access | Not applicable | Through `Future` |

---

### 2.3 Thread Methods

#### 2.3.1 start() Method

The `start()` method creates a new thread of execution and calls the `run()` method in that new thread.

```java
public class StartMethodDemo {
    public static void main(String[] args) {
        Thread thread = new Thread(() -> {
            System.out.println("Running in: " + Thread.currentThread().getName());
        });
        
        System.out.println("Before start: " + thread.getState()); // NEW
        
        thread.start(); // Creates new thread, calls run()
        
        System.out.println("After start: " + thread.getState()); // RUNNABLE
        
        // IMPORTANT: start() can only be called once
        try {
            thread.start(); // IllegalThreadStateException
        } catch (IllegalThreadStateException e) {
            System.out.println("Cannot start thread twice!");
        }
    }
}
```

**start() vs run() - Critical Difference:**

```java
public class StartVsRunDemo {
    public static void main(String[] args) {
        Runnable task = () -> {
            System.out.println("Task running in: " + Thread.currentThread().getName());
        };
        
        Thread thread = new Thread(task);
        
        // WRONG: Calling run() directly - runs in main thread
        System.out.println("\n--- Calling run() directly ---");
        thread.run(); // Output: Task running in: main
        
        // CORRECT: Calling start() - runs in new thread
        System.out.println("\n--- Calling start() ---");
        Thread newThread = new Thread(task);
        newThread.start(); // Output: Task running in: Thread-1
        
        System.out.println("Main thread: " + Thread.currentThread().getName());
    }
}
```

```
┌─────────────────────────────────────────────────────────────────┐
│                     start() vs run()                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  thread.run()                    thread.start()                  │
│  ─────────────                   ──────────────                  │
│                                                                  │
│  Main Thread                     Main Thread    New Thread       │
│      │                               │              │            │
│      ▼                               ▼              │            │
│  ┌───────┐                      ┌─────────┐        │            │
│  │ run() │                      │ start() │────────┤            │
│  │       │                      └────┬────┘        │            │
│  │  ...  │                           │         ┌───▼───┐        │
│  │       │                           │         │ run() │        │
│  └───────┘                           │         │       │        │
│      │                               │         │  ...  │        │
│      ▼                               │         │       │        │
│  continues                           │         └───────┘        │
│  (same thread)                       ▼              │            │
│                                  continues      finishes         │
│                               (parallel execution)               │
│                                                                  │
│  NO MULTITHREADING!              TRUE MULTITHREADING!            │
└─────────────────────────────────────────────────────────────────┘
```

---

#### 2.3.2 run() Method

The `run()` method contains the code that will be executed by the thread. It's either overridden from `Thread` class or defined in `Runnable`.

```java
public class RunMethodDemo {
    public static void main(String[] args) {
        // Custom run() by extending Thread
        Thread customThread = new Thread() {
            @Override
            public void run() {
                for (int i = 0; i < 3; i++) {
                    System.out.println("Custom Thread: " + i);
                }
            }
        };
        
        // run() via Runnable
        Thread runnableThread = new Thread(() -> {
            for (int i = 0; i < 3; i++) {
                System.out.println("Runnable Thread: " + i);
            }
        });
        
        customThread.start();
        runnableThread.start();
    }
}
```

---

#### 2.3.3 sleep() Method

The `sleep()` method pauses the current thread for a specified duration. The thread goes into `TIMED_WAITING` state.

```java
public class SleepMethodDemo {
    public static void main(String[] args) {
        Thread countdownThread = new Thread(() -> {
            for (int i = 5; i >= 1; i--) {
                System.out.println("Countdown: " + i);
                try {
                    Thread.sleep(1000); // Sleep for 1 second
                } catch (InterruptedException e) {
                    System.out.println("Countdown interrupted!");
                    return;
                }
            }
            System.out.println("🚀 Launch!");
        });
        
        countdownThread.start();
    }
}
```

**sleep() with nanoseconds precision:**

```java
public class SleepNanosecondsDemo {
    public static void main(String[] args) throws InterruptedException {
        // sleep(milliseconds)
        Thread.sleep(1000); // 1 second
        
        // sleep(milliseconds, nanoseconds)
        Thread.sleep(1000, 500000); // 1 second + 500,000 nanoseconds
        
        // Using TimeUnit (more readable)
        TimeUnit.SECONDS.sleep(2);      // 2 seconds
        TimeUnit.MILLISECONDS.sleep(500); // 500 milliseconds
        TimeUnit.MINUTES.sleep(1);       // 1 minute
    }
}
```

**Important Points about sleep():**

```java
public class SleepBehaviorDemo {
    public static void main(String[] args) throws InterruptedException {
        Object lock = new Object();
        
        Thread thread = new Thread(() -> {
            synchronized (lock) {
                System.out.println("Thread acquired lock");
                try {
                    // IMPORTANT: sleep() does NOT release the lock
                    Thread.sleep(3000);
                } catch (InterruptedException e) {
                    e.printStackTrace();
                }
                System.out.println("Thread releasing lock");
            }
        });
        
        thread.start();
        Thread.sleep(100); // Give thread time to acquire lock
        
        System.out.println("Main trying to acquire lock...");
        synchronized (lock) {
            // This will wait until thread releases lock (after sleep)
            System.out.println("Main acquired lock");
        }
    }
}
```

```
┌─────────────────────────────────────────────────────────────────┐
│                   sleep() Key Points                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  1. Pauses current thread for specified time                     │
│  2. Thread goes to TIMED_WAITING state                           │
│  3. Does NOT release any locks held by the thread               │
│  4. Throws InterruptedException if thread is interrupted         │
│  5. Minimum sleep time not guaranteed (depends on OS scheduler)  │
│  6. Static method - works on current thread only                 │
│                                                                  │
│  Thread State During sleep():                                    │
│  RUNNABLE ──sleep()──► TIMED_WAITING ──time expires──► RUNNABLE │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

#### 2.3.4 join() Method

The `join()` method allows one thread to wait for another thread to complete. The calling thread goes into `WAITING` or `TIMED_WAITING` state.

```java
public class JoinMethodDemo {
    public static void main(String[] args) {
        Thread thread1 = new Thread(() -> {
            System.out.println("Thread 1 started");
            try {
                Thread.sleep(2000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            System.out.println("Thread 1 completed");
        });
        
        Thread thread2 = new Thread(() -> {
            System.out.println("Thread 2 started");
            try {
                Thread.sleep(1000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            System.out.println("Thread 2 completed");
        });
        
        thread1.start();
        thread2.start();
        
        System.out.println("Main waiting for threads...");
        
        try {
            thread1.join(); // Wait for thread1 to complete
            thread2.join(); // Wait for thread2 to complete
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        
        System.out.println("All threads completed. Main continues.");
    }
}
```

**Output:**
```
Thread 1 started
Thread 2 started
Main waiting for threads...
Thread 2 completed
Thread 1 completed
All threads completed. Main continues.
```

**join() with Timeout:**

```java
public class JoinTimeoutDemo {
    public static void main(String[] args) {
        Thread slowThread = new Thread(() -> {
            System.out.println("Slow thread started");
            try {
                Thread.sleep(10000); // 10 seconds
            } catch (InterruptedException e) {
                System.out.println("Slow thread interrupted");
            }
            System.out.println("Slow thread completed");
        });
        
        slowThread.start();
        
        System.out.println("Main waiting for max 3 seconds...");
        
        try {
            slowThread.join(3000); // Wait max 3 seconds
        } catch (InterruptedException e) {
            e.printStackTrace();
        }
        
        if (slowThread.isAlive()) {
            System.out.println("Thread still running. Main continuing without waiting.");
            slowThread.interrupt(); // Optionally interrupt
        } else {
            System.out.println("Thread completed within timeout.");
        }
        
        System.out.println("Main finished");
    }
}
```

```
┌─────────────────────────────────────────────────────────────────┐
│                        join() Diagrams                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Without join():                With join():                     │
│  ────────────────               ─────────────                    │
│                                                                  │
│  Main      Thread1              Main      Thread1                │
│    │         │                    │         │                    │
│    │    ┌────┴────┐               │    ┌────┴────┐               │
│    │    │ Working │               │    │ Working │               │
│    │    └────┬────┘               │    └────┬────┘               │
│    │         │                    │         │                    │
│    ▼         │               join()────────►│                    │
│  Main        │             (wait) │         │                    │
│  ends        │                    │         │                    │
│              ▼                    │         ▼                    │
│           Thread1                 │     Thread1                  │
│           ends                    │      ends                    │
│                                   ◄─────────┘                    │
│                                   │                              │
│                                   ▼                              │
│                                 Main                             │
│                                 ends                             │
│                                                                  │
│  Problem: Main might end        Solution: Main waits for         │
│  before Thread1 finishes        Thread1 to complete              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**Practical Example: Parallel Processing with join():**

```java
public class ParallelProcessingDemo {
    public static void main(String[] args) throws InterruptedException {
        int[] numbers = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10};
        int[] partialSums = new int[2];
        
        // Thread 1: Sum first half
        Thread thread1 = new Thread(() -> {
            int sum = 0;
            for (int i = 0; i < 5; i++) {
                sum += numbers[i];
            }
            partialSums[0] = sum;
            System.out.println("Thread 1 sum: " + sum);
        });
        
        // Thread 2: Sum second half
        Thread thread2 = new Thread(() -> {
            int sum = 0;
            for (int i = 5; i < 10; i++) {
                sum += numbers[i];
            }
            partialSums[1] = sum;
            System.out.println("Thread 2 sum: " + sum);
        });
        
        thread1.start();
        thread2.start();
        
        // Wait for both threads to complete
        thread1.join();
        thread2.join();
        
        // Combine results
        int totalSum = partialSums[0] + partialSums[1];
        System.out.println("Total sum: " + totalSum); // 55
    }
}
```

---

#### 2.3.5 yield() Method

The `yield()` method hints to the scheduler that the current thread is willing to give up its current use of CPU. It's just a hint - the scheduler may ignore it.

```java
public class YieldMethodDemo {
    public static void main(String[] args) {
        Runnable task = () -> {
            for (int i = 1; i <= 5; i++) {
                System.out.println(Thread.currentThread().getName() + " - " + i);
                
                // Give other threads a chance to run
                if (i % 2 == 0) {
                    System.out.println(Thread.currentThread().getName() + " yielding...");
                    Thread.yield();
                }
            }
        };
        
        Thread thread1 = new Thread(task, "Low-Priority");
        Thread thread2 = new Thread(task, "High-Priority");
        
        thread1.setPriority(Thread.MIN_PRIORITY);
        thread2.setPriority(Thread.MAX_PRIORITY);
        
        thread1.start();
        thread2.start();
    }
}
```

```
┌─────────────────────────────────────────────────────────────────┐
│                     yield() vs sleep()                           │
├────────────────────────────────┬────────────────────────────────┤
│           yield()              │            sleep()             │
├────────────────────────────────┼────────────────────────────────┤
│ Hint to scheduler              │ Guaranteed pause               │
│ May or may not pause           │ Will pause for specified time  │
│ Stays in RUNNABLE state        │ Goes to TIMED_WAITING state    │
│ No time specified              │ Time must be specified         │
│ Does not throw exception       │ Throws InterruptedException    │
│ Gives chance to same/higher    │ Gives chance to all threads    │
│ priority threads               │                                │
└────────────────────────────────┴────────────────────────────────┘
```

**When to use yield():**
- In compute-intensive loops to prevent thread starvation
- When implementing cooperative multitasking
- Testing concurrent code behavior

**Important:** `yield()` is rarely used in production code because:
1. Its behavior is platform-dependent
2. Modern JVMs and OS schedulers handle thread scheduling efficiently
3. It's just a hint that can be ignored

---

#### 2.3.6 interrupt() Method

The `interrupt()` method is used to interrupt a thread. It sets the thread's interrupt flag or throws `InterruptedException` if the thread is in a blocking state.

```java
public class InterruptMethodDemo {
    public static void main(String[] args) throws InterruptedException {
        // Scenario 1: Interrupting a sleeping thread
        Thread sleepingThread = new Thread(() -> {
            try {
                System.out.println("Thread going to sleep...");
                Thread.sleep(10000); // Sleep for 10 seconds
                System.out.println("Thread woke up normally");
            } catch (InterruptedException e) {
                System.out.println("Thread was interrupted during sleep!");
                // Good practice: restore interrupt status
                Thread.currentThread().interrupt();
            }
        });
        
        sleepingThread.start();
        Thread.sleep(1000); // Let it start sleeping
        
        System.out.println("Main interrupting sleeping thread...");
        sleepingThread.interrupt();
        
        sleepingThread.join();
        
        // Scenario 2: Interrupting a running thread (checking flag)
        Thread runningThread = new Thread(() -> {
            int count = 0;
            while (!Thread.currentThread().isInterrupted()) {
                count++;
                // Simulate work
                if (count % 1000000 == 0) {
                    System.out.println("Running thread count: " + count);
                }
            }
            System.out.println("Running thread stopped at count: " + count);
        });
        
        runningThread.start();
        Thread.sleep(100); // Let it run for a bit
        
        System.out.println("Main interrupting running thread...");
        runningThread.interrupt();
        
        runningThread.join();
        System.out.println("Main finished");
    }
}
```

**Interrupt Handling Best Practices:**

```java
public class InterruptBestPractices {
    
    // Pattern 1: Using InterruptedException
    public void blockingOperation() {
        try {
            while (true) {
                // Do some work
                Thread.sleep(100);
            }
        } catch (InterruptedException e) {
            // Option A: Restore interrupt status and return
            Thread.currentThread().interrupt();
            // Cleanup and exit
            System.out.println("Interrupted, cleaning up...");
        }
    }
    
    // Pattern 2: Checking interrupt flag
    public void nonBlockingOperation() {
        while (!Thread.currentThread().isInterrupted()) {
            // Do work
            performTask();
        }
        System.out.println("Thread interrupted, exiting...");
    }
    
    // Pattern 3: Handling in loops with blocking calls
    public void mixedOperation() {
        try {
            while (!Thread.currentThread().isInterrupted()) {
                performTask();
                Thread.sleep(100); // Blocking call
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } finally {
            cleanup();
        }
    }
    
    private void performTask() { /* ... */ }
    private void cleanup() { /* ... */ }
}
```

```
┌─────────────────────────────────────────────────────────────────┐
│                    Interrupt Mechanism                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Thread State when interrupt() called:                          │
│                                                                  │
│  ┌─────────────────┐                                            │
│  │ RUNNABLE        │ → Sets interrupt flag = true               │
│  │ (not blocking)  │   Thread should check isInterrupted()      │
│  └─────────────────┘                                            │
│                                                                  │
│  ┌─────────────────┐                                            │
│  │ WAITING /       │ → Throws InterruptedException              │
│  │ TIMED_WAITING   │   Clears interrupt flag                    │
│  │ (sleep/wait/    │   Thread wakes up immediately              │
│  │  join)          │                                            │
│  └─────────────────┘                                            │
│                                                                  │
│  ┌─────────────────┐                                            │
│  │ BLOCKED         │ → Sets interrupt flag                      │
│  │(waiting for lock)│   InterruptedException thrown when        │
│  │                 │   lock is acquired                         │
│  └─────────────────┘                                            │
│                                                                  │
│  Key Methods:                                                    │
│  • interrupt()        - Set interrupt flag / throw exception    │
│  • isInterrupted()    - Check flag (doesn't clear)              │
│  • Thread.interrupted() - Check flag (clears it)                │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

**isInterrupted() vs Thread.interrupted():**

```java
public class InterruptStatusDemo {
    public static void main(String[] args) {
        Thread thread = new Thread(() -> {
            // Thread.interrupted() - static, clears flag
            System.out.println("First check (interrupted()): " + Thread.interrupted());  // true
            System.out.println("Second check (interrupted()): " + Thread.interrupted()); // false (cleared)
        });
        
        thread.start();
        thread.interrupt();
        
        // Give thread time to run
        try { Thread.sleep(100); } catch (InterruptedException e) {}
        
        // Another thread demonstrating isInterrupted()
        Thread thread2 = new Thread(() -> {
            // isInterrupted() - instance method, does NOT clear flag
            System.out.println("\nFirst check (isInterrupted()): " + 
                             Thread.currentThread().isInterrupted());  // true
            System.out.println("Second check (isInterrupted()): " + 
                             Thread.currentThread().isInterrupted()); // true (not cleared)
        });
        
        thread2.start();
        thread2.interrupt();
    }
}
```

---

### 2.4 Thread Methods Summary Table

| Method | Description | Static | Throws InterruptedException |
|--------|-------------|--------|----------------------------|
| `start()` | Starts thread execution, calls run() | No | No |
| `run()` | Contains code executed by thread | No | No |
| `sleep(ms)` | Pauses current thread for specified time | Yes | Yes |
| `join()` | Waits for thread to complete | No | Yes |
| `join(ms)` | Waits for thread with timeout | No | Yes |
| `yield()` | Hints scheduler to pause current thread | Yes | No |
| `interrupt()` | Interrupts a thread | No | No |
| `isInterrupted()` | Checks if thread is interrupted | No | No |
| `interrupted()` | Checks and clears interrupt status | Yes | No |
| `isAlive()` | Checks if thread is running | No | No |
| `setName()/getName()` | Set/get thread name | No | No |
| `setPriority()/getPriority()` | Set/get thread priority (1-10) | No | No |
| `setDaemon()/isDaemon()` | Set/check daemon status | No | No |
| `getState()` | Returns thread state | No | No |
| `getId()` | Returns thread ID | No | No |

---

### 2.5 Daemon Threads

Daemon threads are background threads that don't prevent the JVM from exiting. When all non-daemon threads finish, the JVM terminates, killing any remaining daemon threads.

```java
public class DaemonThreadDemo {
    public static void main(String[] args) throws InterruptedException {
        Thread daemonThread = new Thread(() -> {
            while (true) {
                System.out.println("Daemon thread running...");
                try {
                    Thread.sleep(500);
                } catch (InterruptedException e) {
                    break;
                }
            }
        });
        
        // Must set daemon BEFORE starting thread
        daemonThread.setDaemon(true);
        daemonThread.start();
        
        System.out.println("Is daemon: " + daemonThread.isDaemon()); // true
        
        // Main thread sleeps for 2 seconds then exits
        Thread.sleep(2000);
        System.out.println("Main thread ending...");
        
        // When main() ends, daemon thread is also terminated
        // No need to explicitly stop it
    }
}
```

**Common Daemon Thread Use Cases:**
- Garbage Collection
- Background monitoring
- Auto-save features
- Logging services

```
┌─────────────────────────────────────────────────────────────────┐
│                   Daemon vs User Threads                         │
├────────────────────────────────┬────────────────────────────────┤
│        User Thread             │        Daemon Thread           │
├────────────────────────────────┼────────────────────────────────┤
│ Created by default             │ Must explicitly setDaemon(true)│
│ JVM waits for completion       │ JVM doesn't wait               │
│ High priority tasks            │ Background/support tasks       │
│ Examples: Main thread,         │ Examples: GC thread,           │
│ UI thread, Worker threads      │ Signal handler, Finalizer      │
└────────────────────────────────┴────────────────────────────────┘
```

---

### 2.6 Complete Example: Thread Methods in Action

```java
import java.util.concurrent.*;

public class ComprehensiveThreadDemo {
    public static void main(String[] args) {
        System.out.println("=== Comprehensive Thread Demo ===\n");
        
        // 1. Thread using Runnable
        Thread downloadThread = new Thread(() -> {
            System.out.println("[Download] Starting download...");
            for (int i = 0; i <= 100; i += 20) {
                System.out.println("[Download] Progress: " + i + "%");
                try {
                    Thread.sleep(500);
                } catch (InterruptedException e) {
                    System.out.println("[Download] Interrupted!");
                    return;
                }
            }
            System.out.println("[Download] Complete!");
        }, "DownloadThread");
        
        // 2. Thread using Callable with Future
        ExecutorService executor = Executors.newSingleThreadExecutor();
        Future<Integer> calculationFuture = executor.submit(() -> {
            System.out.println("[Calculation] Starting complex calculation...");
            Thread.sleep(2000);
            return 42; // The answer to everything
        });
        
        // 3. Daemon thread for monitoring
        Thread monitorThread = new Thread(() -> {
            while (!Thread.currentThread().isInterrupted()) {
                System.out.println("[Monitor] System OK - " + 
                    Runtime.getRuntime().freeMemory() / 1024 + " KB free");
                try {
                    Thread.sleep(1000);
                } catch (InterruptedException e) {
                    break;
                }
            }
        }, "MonitorDaemon");
        monitorThread.setDaemon(true);
        
        // Start all threads
        downloadThread.start();
        monitorThread.start();
        
        System.out.println("\n--- Thread Info ---");
        System.out.println("Download Thread - Name: " + downloadThread.getName() + 
                          ", State: " + downloadThread.getState() + 
                          ", ID: " + downloadThread.getId());
        System.out.println("Monitor Thread - Daemon: " + monitorThread.isDaemon());
        
        // Wait for download to complete
        try {
            downloadThread.join();
            System.out.println("\n[Main] Download thread finished");
            
            // Get calculation result with timeout
            Integer result = calculationFuture.get(5, TimeUnit.SECONDS);
            System.out.println("[Main] Calculation result: " + result);
            
        } catch (InterruptedException | ExecutionException | TimeoutException e) {
            e.printStackTrace();
        }
        
        executor.shutdown();
        System.out.println("\n[Main] All tasks completed!");
    }
}
```

**Expected Output:**
```
=== Comprehensive Thread Demo ===

[Download] Starting download...
[Download] Progress: 0%
[Monitor] System OK - 252416 KB free

--- Thread Info ---
Download Thread - Name: DownloadThread, State: RUNNABLE, ID: 14
Monitor Thread - Daemon: true
[Calculation] Starting complex calculation...
[Download] Progress: 20%
[Monitor] System OK - 252416 KB free
[Download] Progress: 40%
[Monitor] System OK - 252416 KB free
[Download] Progress: 60%
[Download] Progress: 80%
[Monitor] System OK - 252416 KB free
[Download] Progress: 100%
[Download] Complete!

[Main] Download thread finished
[Monitor] System OK - 251392 KB free
[Main] Calculation result: 42

[Main] All tasks completed!
```

---

## 3. Thread Synchronization

When multiple threads access shared resources concurrently, we need mechanisms to ensure data consistency and prevent corruption. Thread synchronization is the coordination of multiple threads to ensure safe access to shared resources.

### 3.1 Race Conditions

A **race condition** occurs when two or more threads access shared data simultaneously, and the final result depends on the timing/order of thread execution.

#### The Problem: Unsynchronized Counter

```java
public class RaceConditionDemo {
    private int counter = 0;
    
    public void increment() {
        counter++; // NOT atomic! (read -> modify -> write)
    }
    
    public int getCounter() {
        return counter;
    }
    
    public static void main(String[] args) throws InterruptedException {
        RaceConditionDemo demo = new RaceConditionDemo();
        
        // Create 1000 threads, each incrementing counter 1000 times
        Thread[] threads = new Thread[1000];
        
        for (int i = 0; i < 1000; i++) {
            threads[i] = new Thread(() -> {
                for (int j = 0; j < 1000; j++) {
                    demo.increment();
                }
            });
            threads[i].start();
        }
        
        // Wait for all threads to complete
        for (Thread t : threads) {
            t.join();
        }
        
        // Expected: 1,000,000 but actual will be less!
        System.out.println("Expected: 1000000");
        System.out.println("Actual: " + demo.getCounter());
    }
}
```

**Output (varies each run):**
```
Expected: 1000000
Actual: 987654  (or some other number less than 1000000)
```

#### Why Does This Happen?

```
┌─────────────────────────────────────────────────────────────────────┐
│               counter++ is NOT ATOMIC                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  counter++ actually involves 3 CPU operations:                       │
│                                                                      │
│  1. READ:   Load value from memory to register                       │
│  2. MODIFY: Increment the value                                      │
│  3. WRITE:  Store value back to memory                               │
│                                                                      │
│  Race Condition Example (counter = 5):                               │
│                                                                      │
│  Thread A             Memory              Thread B                   │
│      │                counter=5               │                      │
│      │                   │                    │                      │
│   1. READ 5 ◄────────────┤                    │                      │
│      │                   │            1. READ 5 ◄─────               │
│   2. 5+1=6               │               2. 5+1=6                    │
│      │                   │                    │                      │
│   3. WRITE 6 ────────────►                    │                      │
│      │                counter=6               │                      │
│      │                   │            3. WRITE 6 ─────►              │
│      │                counter=6 (should be 7!)│                      │
│                                                                      │
│  Lost Update! One increment was lost due to race condition.          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

### 3.2 Critical Section

A **critical section** is a portion of code that accesses shared resources and must be executed by only one thread at a time to maintain data consistency.

```
┌─────────────────────────────────────────────────────────────────────┐
│                      CRITICAL SECTION                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  public void transferMoney(Account from, Account to, int amount) {   │
│      // Non-critical section                                         │
│      validateAccounts(from, to);                                     │
│      logTransferAttempt(from, to, amount);                           │
│                                                                      │
│      ╔═══════════════════════════════════════════════════════╗      │
│      ║           CRITICAL SECTION START                       ║      │
│      ╠═══════════════════════════════════════════════════════╣      │
│      ║  if (from.getBalance() >= amount) {                    ║      │
│      ║      from.withdraw(amount);                            ║      │
│      ║      to.deposit(amount);                               ║      │
│      ║  }                                                     ║      │
│      ╠═══════════════════════════════════════════════════════╣      │
│      ║           CRITICAL SECTION END                         ║      │
│      ╚═══════════════════════════════════════════════════════╝      │
│                                                                      │
│      // Non-critical section                                         │
│      logTransferComplete(from, to, amount);                          │
│  }                                                                   │
│                                                                      │
│  Only ONE thread should execute the critical section at a time!      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Properties of Critical Section:
- **Mutual Exclusion**: Only one thread can execute at a time
- **Progress**: If no thread is in critical section, waiting threads should be able to enter
- **Bounded Waiting**: Threads should not wait indefinitely

---

### 3.3 synchronized Keyword

Java provides the `synchronized` keyword to protect critical sections. It ensures that only one thread can execute a synchronized block/method at a time on the same object.

#### 3.3.1 Synchronized Methods

```java
public class SynchronizedCounter {
    private int counter = 0;
    
    // Synchronized instance method
    public synchronized void increment() {
        counter++; // Now thread-safe!
    }
    
    public synchronized void decrement() {
        counter--;
    }
    
    public synchronized int getCounter() {
        return counter;
    }
}
```

**Fixed Race Condition Example:**

```java
public class SynchronizedDemo {
    private int counter = 0;
    
    public synchronized void increment() {
        counter++;
    }
    
    public int getCounter() {
        return counter;
    }
    
    public static void main(String[] args) throws InterruptedException {
        SynchronizedDemo demo = new SynchronizedDemo();
        
        Thread[] threads = new Thread[1000];
        
        for (int i = 0; i < 1000; i++) {
            threads[i] = new Thread(() -> {
                for (int j = 0; j < 1000; j++) {
                    demo.increment();
                }
            });
            threads[i].start();
        }
        
        for (Thread t : threads) {
            t.join();
        }
        
        // Now always: 1,000,000
        System.out.println("Expected: 1000000");
        System.out.println("Actual: " + demo.getCounter());
    }
}
```

**Output:**
```
Expected: 1000000
Actual: 1000000  ✓ Correct every time!
```

#### How synchronized Works

```
┌─────────────────────────────────────────────────────────────────────┐
│                 MONITOR / INTRINSIC LOCK                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Every Java object has an intrinsic lock (monitor).                  │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    Java Object                               │    │
│  │  ┌───────────────────────────────────────────────────────┐  │    │
│  │  │                  Object Header                         │  │    │
│  │  │  ┌─────────────┐  ┌──────────────────────────────┐    │  │    │
│  │  │  │  Mark Word  │  │  Class Pointer               │    │  │    │
│  │  │  │  (Lock Info)│  │                              │    │  │    │
│  │  │  └─────────────┘  └──────────────────────────────┘    │  │    │
│  │  └───────────────────────────────────────────────────────┘  │    │
│  │  ┌───────────────────────────────────────────────────────┐  │    │
│  │  │              Instance Data (fields)                    │  │    │
│  │  └───────────────────────────────────────────────────────┘  │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  synchronized method/block acquires this intrinsic lock:            │
│                                                                      │
│  Thread A                         Thread B                          │
│     │                                │                              │
│     ▼                                │                              │
│  ┌──────────────┐                    │                              │
│  │ Acquire Lock │◄──── Lock ────►    │                              │
│  └──────┬───────┘                    │                              │
│         │                            ▼                              │
│         ▼                     ┌──────────────┐                      │
│  ┌──────────────┐             │  BLOCKED     │                      │
│  │ Execute      │             │  (waiting    │                      │
│  │ Critical     │             │   for lock)  │                      │
│  │ Section      │             └──────────────┘                      │
│  └──────┬───────┘                    │                              │
│         │                            │                              │
│         ▼                            │                              │
│  ┌──────────────┐                    │                              │
│  │ Release Lock │────► Lock ─────────┤                              │
│  └──────────────┘                    ▼                              │
│                               ┌──────────────┐                      │
│                               │ Acquire Lock │                      │
│                               └──────┬───────┘                      │
│                                      │                              │
│                                      ▼                              │
│                               ┌──────────────┐                      │
│                               │ Execute      │                      │
│                               │ Critical     │                      │
│                               │ Section      │                      │
│                               └──────────────┘                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

### 3.4 Object-Level vs Class-Level Locking

Java supports two types of locking:

#### 3.4.1 Object-Level Locking (Instance Lock)

Lock is on the object instance. Different instances can be accessed by different threads simultaneously.

```java
public class ObjectLevelLocking {
    private int count = 0;
    
    // Synchronized on 'this' object
    public synchronized void increment() {
        count++;
    }
    
    // Equivalent to above
    public void incrementExplicit() {
        synchronized (this) {
            count++;
        }
    }
}
```

```java
public class ObjectLevelDemo {
    public static void main(String[] args) {
        ObjectLevelLocking obj1 = new ObjectLevelLocking();
        ObjectLevelLocking obj2 = new ObjectLevelLocking();
        
        // Thread 1 locks on obj1
        Thread t1 = new Thread(() -> {
            obj1.increment(); // Locks obj1
        });
        
        // Thread 2 locks on obj2
        Thread t2 = new Thread(() -> {
            obj2.increment(); // Locks obj2 (different lock)
        });
        
        // Both can execute simultaneously!
        t1.start();
        t2.start();
    }
}
```

```
┌─────────────────────────────────────────────────────────────────────┐
│                   OBJECT-LEVEL LOCKING                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│    ┌─────────────┐                    ┌─────────────┐               │
│    │   Object 1  │                    │   Object 2  │               │
│    │   (Lock 1)  │                    │   (Lock 2)  │               │
│    └──────┬──────┘                    └──────┬──────┘               │
│           │                                  │                       │
│           ▼                                  ▼                       │
│    ┌─────────────┐                    ┌─────────────┐               │
│    │  Thread A   │                    │  Thread B   │               │
│    │  (holds     │                    │  (holds     │               │
│    │   Lock 1)   │                    │   Lock 2)   │               │
│    └─────────────┘                    └─────────────┘               │
│                                                                      │
│    Both threads can run simultaneously (different locks)             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### 3.4.2 Class-Level Locking (Static Lock)

Lock is on the Class object. Only one thread can execute across ALL instances.

```java
public class ClassLevelLocking {
    private static int count = 0;
    
    // Synchronized on Class object
    public static synchronized void increment() {
        count++;
    }
    
    // Equivalent to above
    public static void incrementExplicit() {
        synchronized (ClassLevelLocking.class) {
            count++;
        }
    }
    
    // Instance method with class-level lock
    public void instanceMethodWithClassLock() {
        synchronized (ClassLevelLocking.class) {
            count++;
        }
    }
}
```

```java
public class ClassLevelDemo {
    public static void main(String[] args) {
        ClassLevelLocking obj1 = new ClassLevelLocking();
        ClassLevelLocking obj2 = new ClassLevelLocking();
        
        // Thread 1 locks on class
        Thread t1 = new Thread(() -> {
            ClassLevelLocking.increment(); // Locks ClassLevelLocking.class
        });
        
        // Thread 2 also locks on same class
        Thread t2 = new Thread(() -> {
            ClassLevelLocking.increment(); // Same lock!
        });
        
        // Threads must wait for each other
        t1.start();
        t2.start();
    }
}
```

```
┌─────────────────────────────────────────────────────────────────────┐
│                   CLASS-LEVEL LOCKING                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│                   ┌─────────────────────┐                            │
│                   │ ClassLevelLocking   │                            │
│                   │      .class         │                            │
│                   │    (Single Lock)    │                            │
│                   └──────────┬──────────┘                            │
│                              │                                       │
│              ┌───────────────┼───────────────┐                       │
│              │               │               │                       │
│              ▼               ▼               ▼                       │
│       ┌─────────────┐ ┌─────────────┐ ┌─────────────┐               │
│       │  Object 1   │ │  Object 2   │ │  Object 3   │               │
│       └─────────────┘ └─────────────┘ └─────────────┘               │
│              │               │               │                       │
│              └───────────────┼───────────────┘                       │
│                              │                                       │
│                        All share the                                │
│                        SAME class lock                              │
│                                                                      │
│   Thread A ───► BLOCKED (waiting)                                   │
│   Thread B ───► EXECUTING (holds class lock)                        │
│   Thread C ───► BLOCKED (waiting)                                   │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Comparison Table

| Aspect | Object-Level Lock | Class-Level Lock |
|--------|-------------------|------------------|
| Lock on | Instance (`this`) | Class (`ClassName.class`) |
| Keyword | `synchronized` on instance method | `synchronized` on static method |
| Scope | Per object instance | All instances share same lock |
| Use case | Instance-specific data | Static/shared data |
| Concurrency | High (multiple objects) | Low (single lock) |

---

### 3.5 synchronized Blocks

Synchronized blocks allow finer-grained control over synchronization, enabling you to lock only the critical section rather than the entire method.

#### Syntax

```java
synchronized (lockObject) {
    // Critical section
}
```

#### Advantages Over Synchronized Methods

```java
public class SynchronizedBlockDemo {
    private final Object lock1 = new Object();
    private final Object lock2 = new Object();
    
    private List<String> list1 = new ArrayList<>();
    private List<String> list2 = new ArrayList<>();
    
    // Problem: Entire method synchronized - excessive locking
    public synchronized void addToListsBad(String item1, String item2) {
        // Some expensive non-critical operation
        String processed1 = processItem(item1); // No lock needed!
        String processed2 = processItem(item2); // No lock needed!
        
        // Only these need synchronization
        list1.add(processed1);
        list2.add(processed2);
    }
    
    // Solution: Use synchronized blocks - lock only what's necessary
    public void addToListsGood(String item1, String item2) {
        // Non-critical operations - no lock held
        String processed1 = processItem(item1);
        String processed2 = processItem(item2);
        
        // Only lock when accessing shared data
        synchronized (lock1) {
            list1.add(processed1);
        }
        
        synchronized (lock2) {
            list2.add(processed2);
        }
    }
    
    // Even better: Different locks for independent resources
    public void addToListsParallel() {
        Thread t1 = new Thread(() -> {
            synchronized (lock1) {
                list1.add("item1");
            }
        });
        
        Thread t2 = new Thread(() -> {
            synchronized (lock2) {
                list2.add("item2"); // Can run in parallel with t1!
            }
        });
        
        t1.start();
        t2.start();
    }
    
    private String processItem(String item) {
        // Expensive operation
        return item.toUpperCase();
    }
}
```

#### Practical Example: Bank Account

```java
public class BankAccount {
    private double balance;
    private final Object lock = new Object();
    private final List<String> transactionHistory = new ArrayList<>();
    
    public BankAccount(double initialBalance) {
        this.balance = initialBalance;
    }
    
    public void deposit(double amount) {
        // Validate outside synchronized block
        if (amount <= 0) {
            throw new IllegalArgumentException("Amount must be positive");
        }
        
        synchronized (lock) {
            balance += amount;
            transactionHistory.add("Deposited: " + amount);
        }
        
        // Logging outside synchronized block
        System.out.println("Deposited: " + amount);
    }
    
    public boolean withdraw(double amount) {
        if (amount <= 0) {
            throw new IllegalArgumentException("Amount must be positive");
        }
        
        synchronized (lock) {
            if (balance >= amount) {
                balance -= amount;
                transactionHistory.add("Withdrew: " + amount);
                return true;
            }
            return false;
        }
    }
    
    public double getBalance() {
        synchronized (lock) {
            return balance;
        }
    }
    
    // Transfer between accounts - be careful of deadlock!
    public static void transfer(BankAccount from, BankAccount to, double amount) {
        // Always lock in consistent order to prevent deadlock
        Object firstLock = from.hashCode() < to.hashCode() ? from.lock : to.lock;
        Object secondLock = from.hashCode() < to.hashCode() ? to.lock : from.lock;
        
        synchronized (firstLock) {
            synchronized (secondLock) {
                if (from.balance >= amount) {
                    from.balance -= amount;
                    to.balance += amount;
                    from.transactionHistory.add("Transferred out: " + amount);
                    to.transactionHistory.add("Transferred in: " + amount);
                }
            }
        }
    }
}
```

---

### 3.6 Deadlock

A **deadlock** occurs when two or more threads are blocked forever, each waiting for a resource held by the other.

```
┌─────────────────────────────────────────────────────────────────────┐
│                         DEADLOCK                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│                    ┌─────────────┐                                   │
│                    │   Lock A    │                                   │
│                    └──────┬──────┘                                   │
│                           │                                          │
│           holds ──────────┴────────── waits for                      │
│              │                            │                          │
│              ▼                            ▼                          │
│       ┌─────────────┐            ┌─────────────┐                    │
│       │  Thread 1   │            │  Thread 2   │                    │
│       └─────────────┘            └─────────────┘                    │
│              │                            │                          │
│           waits for ─────────┬────── holds                          │
│                              │                                       │
│                    ┌─────────▼─────┐                                 │
│                    │    Lock B     │                                 │
│                    └───────────────┘                                 │
│                                                                      │
│   Thread 1: Holds Lock A, waiting for Lock B                        │
│   Thread 2: Holds Lock B, waiting for Lock A                        │
│                                                                      │
│   RESULT: Both threads wait forever! (Deadlock)                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Deadlock Example

```java
public class DeadlockDemo {
    private static final Object LOCK_A = new Object();
    private static final Object LOCK_B = new Object();
    
    public static void main(String[] args) {
        Thread thread1 = new Thread(() -> {
            System.out.println("Thread 1: Trying to acquire Lock A");
            synchronized (LOCK_A) {
                System.out.println("Thread 1: Acquired Lock A");
                
                // Simulate some work
                try { Thread.sleep(100); } catch (InterruptedException e) {}
                
                System.out.println("Thread 1: Trying to acquire Lock B");
                synchronized (LOCK_B) {
                    System.out.println("Thread 1: Acquired Lock B");
                }
            }
        }, "Thread-1");
        
        Thread thread2 = new Thread(() -> {
            System.out.println("Thread 2: Trying to acquire Lock B");
            synchronized (LOCK_B) {
                System.out.println("Thread 2: Acquired Lock B");
                
                // Simulate some work
                try { Thread.sleep(100); } catch (InterruptedException e) {}
                
                System.out.println("Thread 2: Trying to acquire Lock A");
                synchronized (LOCK_A) {
                    System.out.println("Thread 2: Acquired Lock A");
                }
            }
        }, "Thread-2");
        
        thread1.start();
        thread2.start();
        
        // Program will hang here due to deadlock!
    }
}
```

**Output:**
```
Thread 1: Trying to acquire Lock A
Thread 1: Acquired Lock A
Thread 2: Trying to acquire Lock B
Thread 2: Acquired Lock B
Thread 1: Trying to acquire Lock B
Thread 2: Trying to acquire Lock A
... (program hangs - DEADLOCK!)
```

#### Four Conditions for Deadlock (Coffman Conditions)

All four must be present for deadlock to occur:

| Condition | Description |
|-----------|-------------|
| **Mutual Exclusion** | Resources cannot be shared (exclusive access) |
| **Hold and Wait** | Thread holds resource while waiting for another |
| **No Preemption** | Resources cannot be forcibly taken from thread |
| **Circular Wait** | Circular chain of threads waiting for resources |

#### Preventing Deadlock

```java
public class DeadlockPrevention {
    private static final Object LOCK_A = new Object();
    private static final Object LOCK_B = new Object();
    
    // Solution 1: Lock Ordering - Always acquire locks in the same order
    public static void method1() {
        synchronized (LOCK_A) {      // Always A first
            synchronized (LOCK_B) {  // Then B
                // Work
            }
        }
    }
    
    public static void method2() {
        synchronized (LOCK_A) {      // Always A first (same order)
            synchronized (LOCK_B) {  // Then B
                // Work
            }
        }
    }
    
    // Solution 2: Lock with Timeout using ReentrantLock
    private static final ReentrantLock lockA = new ReentrantLock();
    private static final ReentrantLock lockB = new ReentrantLock();
    
    public static void safeMethod() {
        boolean acquiredA = false;
        boolean acquiredB = false;
        
        try {
            acquiredA = lockA.tryLock(1, TimeUnit.SECONDS);
            acquiredB = lockB.tryLock(1, TimeUnit.SECONDS);
            
            if (acquiredA && acquiredB) {
                // Do work
            } else {
                // Could not acquire both locks - handle gracefully
                System.out.println("Could not acquire locks, will retry...");
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        } finally {
            if (acquiredB) lockB.unlock();
            if (acquiredA) lockA.unlock();
        }
    }
    
    // Solution 3: Single Lock (reduce complexity)
    private static final Object SINGLE_LOCK = new Object();
    
    public static void singleLockMethod() {
        synchronized (SINGLE_LOCK) {
            // Access both resources
        }
    }
}
```

#### Deadlock Detection

```java
import java.lang.management.*;

public class DeadlockDetector {
    public static void detectDeadlock() {
        ThreadMXBean threadMBean = ManagementFactory.getThreadMXBean();
        long[] deadlockedThreads = threadMBean.findDeadlockedThreads();
        
        if (deadlockedThreads != null) {
            System.out.println("Deadlock detected!");
            ThreadInfo[] threadInfos = threadMBean.getThreadInfo(deadlockedThreads);
            
            for (ThreadInfo info : threadInfos) {
                System.out.println("Deadlocked thread: " + info.getThreadName());
                System.out.println("  Blocked on: " + info.getLockName());
                System.out.println("  Owned by: " + info.getLockOwnerName());
            }
        } else {
            System.out.println("No deadlock detected.");
        }
    }
}
```

---

### 3.7 Livelock and Starvation

#### 3.7.1 Livelock

A **livelock** occurs when threads are not blocked but are unable to make progress because they keep responding to each other's actions.

```
┌─────────────────────────────────────────────────────────────────────┐
│                          LIVELOCK                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Two people meeting in a corridor:                                   │
│                                                                      │
│  Step 1:               Step 2:               Step 3:                │
│  ┌───────────┐         ┌───────────┐         ┌───────────┐          │
│  │ A → → → B │         │ A ← ← ← B │         │ A → → → B │          │
│  │     ↕     │         │     ↕     │         │     ↕     │          │
│  │   Both    │         │   Both    │         │   Both    │          │
│  │  move     │         │  move     │         │  move     │          │
│  │  right    │         │  left     │         │  right    │          │
│  └───────────┘         └───────────┘         └───────────┘          │
│                                                                      │
│  Both keep yielding to each other - no progress made!               │
│                                                                      │
│  Unlike deadlock: Threads are ACTIVE but still can't proceed        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
public class LivelockDemo {
    static class Spoon {
        private Diner owner;
        
        public Spoon(Diner owner) {
            this.owner = owner;
        }
        
        public Diner getOwner() { return owner; }
        public void setOwner(Diner owner) { this.owner = owner; }
        
        public synchronized void use() {
            System.out.println(owner.getName() + " is eating!");
        }
    }
    
    static class Diner {
        private String name;
        private boolean isHungry;
        
        public Diner(String name) {
            this.name = name;
            this.isHungry = true;
        }
        
        public String getName() { return name; }
        public boolean isHungry() { return isHungry; }
        
        public void eatWith(Spoon spoon, Diner otherDiner) {
            while (isHungry) {
                if (spoon.getOwner() != this) {
                    try { Thread.sleep(100); } 
                    catch (InterruptedException e) {}
                    continue;
                }
                
                // "Politely" give spoon to hungry partner
                if (otherDiner.isHungry()) {
                    System.out.println(name + ": Oh, you're hungry " + 
                                     otherDiner.getName() + ", please take the spoon!");
                    spoon.setOwner(otherDiner);
                    continue; // LIVELOCK: Both keep giving spoon to each other!
                }
                
                // Eat
                spoon.use();
                isHungry = false;
                System.out.println(name + ": I'm done eating.");
                spoon.setOwner(otherDiner);
            }
        }
    }
    
    public static void main(String[] args) {
        Diner husband = new Diner("Husband");
        Diner wife = new Diner("Wife");
        Spoon spoon = new Spoon(husband);
        
        Thread t1 = new Thread(() -> husband.eatWith(spoon, wife));
        Thread t2 = new Thread(() -> wife.eatWith(spoon, husband));
        
        t1.start();
        t2.start();
        
        // Both keep yielding to each other - LIVELOCK!
    }
}
```

#### Livelock Solutions

```java
public class LivelockSolution {
    
    // Solution 1: Add randomness to break symmetry
    public void eatWithRandomDelay(Spoon spoon, Diner otherDiner) {
        Random random = new Random();
        while (isHungry) {
            if (spoon.getOwner() != this) {
                try { 
                    // Random delay breaks symmetry
                    Thread.sleep(random.nextInt(100)); 
                } catch (InterruptedException e) {}
                continue;
            }
            
            if (otherDiner.isHungry()) {
                spoon.setOwner(otherDiner);
                continue;
            }
            
            spoon.use();
            isHungry = false;
            spoon.setOwner(otherDiner);
        }
    }
    
    // Solution 2: Priority-based resolution
    public void eatWithPriority(Spoon spoon, Diner otherDiner) {
        while (isHungry) {
            if (spoon.getOwner() != this) {
                try { Thread.sleep(100); } catch (InterruptedException e) {}
                continue;
            }
            
            // Only yield if other has higher priority
            if (otherDiner.isHungry() && otherDiner.getPriority() > this.getPriority()) {
                spoon.setOwner(otherDiner);
                continue;
            }
            
            spoon.use();
            isHungry = false;
            spoon.setOwner(otherDiner);
        }
    }
}
```

#### 3.7.2 Starvation

**Starvation** occurs when a thread is perpetually denied access to resources because other threads are constantly acquiring them first.

```
┌─────────────────────────────────────────────────────────────────────┐
│                        STARVATION                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│                    ┌───────────────────┐                            │
│                    │   Shared Resource │                            │
│                    └─────────┬─────────┘                            │
│                              │                                       │
│        ┌─────────────────────┼─────────────────────┐                │
│        │                     │                     │                │
│        ▼                     ▼                     ▼                │
│  ┌───────────┐         ┌───────────┐         ┌───────────┐          │
│  │ Thread A  │         │ Thread B  │         │ Thread C  │          │
│  │ Priority:9│         │ Priority:9│         │ Priority:1│          │
│  │ Accessing │         │ Accessing │         │ STARVING! │          │
│  │ resource  │         │ resource  │         │  Waiting  │          │
│  │ regularly │         │ regularly │         │  forever  │          │
│  └───────────┘         └───────────┘         └───────────┘          │
│        │                     │                     ▲                │
│        └─────────────────────┴─────────────────────┘                │
│                              │                                       │
│          Low-priority Thread C never gets CPU time                  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
public class StarvationDemo {
    private static final Object lock = new Object();
    
    public static void main(String[] args) {
        // Create high-priority threads
        for (int i = 0; i < 5; i++) {
            Thread highPriority = new Thread(() -> {
                while (true) {
                    synchronized (lock) {
                        System.out.println(Thread.currentThread().getName() + 
                                         " acquired lock");
                        // Hold lock for a while
                        try { Thread.sleep(100); } 
                        catch (InterruptedException e) {}
                    }
                }
            }, "High-Priority-" + i);
            highPriority.setPriority(Thread.MAX_PRIORITY);
            highPriority.start();
        }
        
        // Create low-priority thread - may starve!
        Thread lowPriority = new Thread(() -> {
            while (true) {
                synchronized (lock) {
                    System.out.println("*** LOW PRIORITY THREAD GOT THE LOCK! ***");
                    try { Thread.sleep(100); } 
                    catch (InterruptedException e) {}
                }
            }
        }, "Low-Priority");
        lowPriority.setPriority(Thread.MIN_PRIORITY);
        lowPriority.start();
    }
}
```

#### Causes of Starvation

| Cause | Description |
|-------|-------------|
| **Thread Priority** | High-priority threads monopolize CPU |
| **Long-Running Threads** | Some threads hold locks too long |
| **Unfair Scheduler** | Scheduler favors certain threads |
| **Synchronized Methods** | Long waits in synchronized queue |

#### Preventing Starvation

```java
import java.util.concurrent.locks.ReentrantLock;

public class StarvationPrevention {
    
    // Solution 1: Use fair locks
    private final ReentrantLock fairLock = new ReentrantLock(true); // fairness = true
    
    public void fairMethod() {
        fairLock.lock();
        try {
            // Critical section
            // Threads acquire lock in FIFO order
        } finally {
            fairLock.unlock();
        }
    }
    
    // Solution 2: Avoid priority-based scheduling for critical resources
    // Don't rely on thread priorities for correctness
    
    // Solution 3: Use bounded waiting (timeout)
    public void methodWithTimeout() {
        try {
            if (fairLock.tryLock(5, TimeUnit.SECONDS)) {
                try {
                    // Critical section
                } finally {
                    fairLock.unlock();
                }
            } else {
                // Handle timeout - log, retry later, etc.
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }
}
```

#### Comparison: Deadlock vs Livelock vs Starvation

| Aspect | Deadlock | Livelock | Starvation |
|--------|----------|----------|------------|
| **Thread State** | BLOCKED | RUNNABLE | BLOCKED/RUNNABLE |
| **Progress** | No progress | No progress (active) | Limited progress |
| **Cause** | Circular wait | Responding to each other | Resource monopolization |
| **Detection** | Easy (thread dump) | Difficult | Moderate |
| **Solution** | Lock ordering, timeout | Randomness, priority | Fair locks, timeout |

---

## 4. Inter-Thread Communication

Inter-thread communication allows threads to coordinate their activities by signaling each other when certain conditions are met.

### 4.1 wait() Method

The `wait()` method causes the current thread to release the lock and wait until another thread calls `notify()` or `notifyAll()` on the same object.

```java
// Signature
public final void wait() throws InterruptedException
public final void wait(long timeout) throws InterruptedException
public final void wait(long timeout, int nanos) throws InterruptedException
```

#### Key Points about wait()

```
┌─────────────────────────────────────────────────────────────────────┐
│                        wait() Behavior                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. MUST be called from within synchronized block/method            │
│     (Thread must own the monitor/lock)                              │
│                                                                      │
│  2. When wait() is called:                                          │
│     • Thread releases the lock immediately                          │
│     • Thread goes to WAITING state                                  │
│     • Thread is added to object's wait set                          │
│                                                                      │
│  3. Thread wakes up when:                                           │
│     • Another thread calls notify() and this thread is selected    │
│     • Another thread calls notifyAll()                              │
│     • Another thread interrupts this thread                         │
│     • Timeout expires (if wait(timeout) was used)                   │
│                                                                      │
│  4. After waking up:                                                │
│     • Thread must re-acquire the lock before continuing             │
│     • Thread should always re-check the condition (spurious wakeup) │
│                                                                      │
│  Thread State: RUNNABLE ──wait()──► WAITING ──notify()──► RUNNABLE  │
│                                              (re-acquire lock)       │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
public class WaitDemo {
    private static final Object lock = new Object();
    private static boolean dataAvailable = false;
    
    public static void main(String[] args) {
        // Consumer thread - waits for data
        Thread consumer = new Thread(() -> {
            synchronized (lock) {
                System.out.println("Consumer: Waiting for data...");
                
                // Always use while loop - handles spurious wakeups
                while (!dataAvailable) {
                    try {
                        lock.wait(); // Releases lock and waits
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        return;
                    }
                }
                
                System.out.println("Consumer: Data received! Processing...");
            }
        }, "Consumer");
        
        // Producer thread - produces data
        Thread producer = new Thread(() -> {
            try {
                Thread.sleep(2000); // Simulate data preparation
            } catch (InterruptedException e) {}
            
            synchronized (lock) {
                System.out.println("Producer: Data ready! Notifying consumer...");
                dataAvailable = true;
                lock.notify(); // Wake up waiting thread
            }
        }, "Producer");
        
        consumer.start();
        producer.start();
    }
}
```

**Output:**
```
Consumer: Waiting for data...
Producer: Data ready! Notifying consumer...
Consumer: Data received! Processing...
```

#### Why Use While Loop Instead of If?

```java
// WRONG - Using if (vulnerable to spurious wakeups)
synchronized (lock) {
    if (!condition) {
        lock.wait();
    }
    // May proceed even if condition is still false!
}

// CORRECT - Using while (safe from spurious wakeups)
synchronized (lock) {
    while (!condition) {
        lock.wait();
    }
    // Condition is guaranteed to be true here
}
```

**Spurious Wakeup**: A thread can wake up from wait() without being notified, interrupted, or timing out. This is allowed by JVM specification for performance reasons.

---

### 4.2 notify() Method

The `notify()` method wakes up a single thread that is waiting on the object's monitor.

```java
// Signature
public final void notify()
```

```
┌─────────────────────────────────────────────────────────────────────┐
│                        notify() Behavior                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. MUST be called from within synchronized block/method            │
│                                                                      │
│  2. Wakes up ONE thread from the object's wait set                  │
│     • Which thread is implementation-dependent (not guaranteed)     │
│     • If multiple threads waiting, one is chosen arbitrarily        │
│                                                                      │
│  3. Notifying thread does NOT release lock immediately              │
│     • Lock is released when synchronized block exits                │
│                                                                      │
│  4. Notified thread cannot proceed until it re-acquires lock        │
│                                                                      │
│     Object's Wait Set:                                              │
│     ┌─────────────────────────────────────────┐                     │
│     │ Thread A │ Thread B │ Thread C │        │                     │
│     └────┬─────┴──────────┴──────────┘        │                     │
│          │                                     │                     │
│          ▼                                     │                     │
│     notify() selects ONE (arbitrary)           │                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
public class NotifyDemo {
    private static final Object lock = new Object();
    
    public static void main(String[] args) throws InterruptedException {
        // Create multiple waiting threads
        for (int i = 1; i <= 3; i++) {
            final int threadNum = i;
            new Thread(() -> {
                synchronized (lock) {
                    try {
                        System.out.println("Thread " + threadNum + ": Waiting...");
                        lock.wait();
                        System.out.println("Thread " + threadNum + ": Woken up!");
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    }
                }
            }, "Thread-" + i).start();
        }
        
        Thread.sleep(1000); // Let all threads start waiting
        
        // Notify ONE thread at a time
        for (int i = 1; i <= 3; i++) {
            synchronized (lock) {
                System.out.println("\nMain: Calling notify() #" + i);
                lock.notify(); // Wakes up only ONE thread
            }
            Thread.sleep(500);
        }
    }
}
```

**Output:**
```
Thread 1: Waiting...
Thread 2: Waiting...
Thread 3: Waiting...

Main: Calling notify() #1
Thread 1: Woken up!

Main: Calling notify() #2
Thread 2: Woken up!

Main: Calling notify() #3
Thread 3: Woken up!
```

---

### 4.3 notifyAll() Method

The `notifyAll()` method wakes up ALL threads waiting on the object's monitor.

```java
// Signature
public final void notifyAll()
```

```
┌─────────────────────────────────────────────────────────────────────┐
│                      notifyAll() Behavior                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Object's Wait Set:                                                 │
│  ┌────────────────────────────────────────────┐                     │
│  │ Thread A │ Thread B │ Thread C │ Thread D │                     │
│  └────┬─────┴────┬─────┴────┬─────┴────┬─────┘                     │
│       │          │          │          │                            │
│       └──────────┴─────┬────┴──────────┘                            │
│                        │                                             │
│                   notifyAll()                                        │
│                        │                                             │
│       ┌──────────┬─────┴────┬──────────┐                            │
│       ▼          ▼          ▼          ▼                            │
│  All threads wake up and compete for the lock                       │
│  Only ONE acquires the lock at a time                               │
│  Others wait in BLOCKED state until lock is released               │
│                                                                      │
│  Use when:                                                          │
│  • Multiple threads may need to respond to notification             │
│  • You don't know which thread should wake up                       │
│  • Safer than notify() (avoids lost wakeup problem)                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
public class NotifyAllDemo {
    private static final Object lock = new Object();
    private static int resource = 0;
    
    public static void main(String[] args) throws InterruptedException {
        // Create multiple waiting threads with different conditions
        for (int i = 1; i <= 5; i++) {
            final int requiredValue = i;
            new Thread(() -> {
                synchronized (lock) {
                    // Each thread waits for a different condition
                    while (resource < requiredValue) {
                        try {
                            System.out.println(Thread.currentThread().getName() + 
                                             ": Waiting for resource >= " + requiredValue);
                            lock.wait();
                        } catch (InterruptedException e) {
                            Thread.currentThread().interrupt();
                            return;
                        }
                    }
                    System.out.println(Thread.currentThread().getName() + 
                                     ": Proceeding! Resource = " + resource);
                }
            }, "Thread-" + i).start();
        }
        
        Thread.sleep(1000); // Let all threads start waiting
        
        // Update resource and notify all
        for (int value = 1; value <= 5; value++) {
            synchronized (lock) {
                resource = value;
                System.out.println("\n=== Main: Resource set to " + value + 
                                 ", calling notifyAll() ===");
                lock.notifyAll(); // Wake up ALL threads
            }
            Thread.sleep(500);
        }
    }
}
```

#### notify() vs notifyAll()

| Aspect | notify() | notifyAll() |
|--------|----------|-------------|
| **Threads Woken** | One (arbitrary) | All waiting threads |
| **Performance** | Better (less overhead) | More overhead |
| **Safety** | Can cause lost wakeup | Safer, all threads check condition |
| **Use When** | Single waiter or all wait for same condition | Multiple waiters with different conditions |
| **Risk** | Wrong thread might wake up | None (but more contention) |

---

### 4.4 Producer-Consumer Problem

The Producer-Consumer problem is a classic synchronization problem where producers generate data and put it into a buffer, while consumers take data from the buffer.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    PRODUCER-CONSUMER PATTERN                         │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   PRODUCER                   BUFFER                    CONSUMER      │
│  ┌─────────┐            ┌────────────────┐            ┌─────────┐   │
│  │         │───put()───►│ [1][2][3][ ][ ]│───get()───►│         │   │
│  │ Creates │            │    Bounded     │            │ Consumes│   │
│  │  Data   │            │     Queue      │            │  Data   │   │
│  │         │            └────────────────┘            │         │   │
│  └─────────┘                   │  │                   └─────────┘   │
│       │                        │  │                        │        │
│       │                        │  │                        │        │
│       ▼                        ▼  ▼                        ▼        │
│   Waits if                 Waits if                   Waits if      │
│   buffer full              thread safe               buffer empty   │
│                                                                      │
│  Challenges:                                                         │
│  1. Producer must wait when buffer is FULL                          │
│  2. Consumer must wait when buffer is EMPTY                         │
│  3. Access to buffer must be thread-safe                            │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Complete Producer-Consumer Implementation

```java
import java.util.LinkedList;
import java.util.Queue;

// Bounded Buffer (Thread-safe)
class BoundedBuffer<T> {
    private final Queue<T> buffer = new LinkedList<>();
    private final int capacity;
    
    public BoundedBuffer(int capacity) {
        this.capacity = capacity;
    }
    
    // Producer calls this method
    public synchronized void put(T item) throws InterruptedException {
        // Wait if buffer is full
        while (buffer.size() == capacity) {
            System.out.println(Thread.currentThread().getName() + 
                             ": Buffer full, waiting...");
            wait();
        }
        
        buffer.add(item);
        System.out.println(Thread.currentThread().getName() + 
                         ": Produced " + item + ", Buffer size: " + buffer.size());
        
        // Notify consumers that data is available
        notifyAll();
    }
    
    // Consumer calls this method
    public synchronized T get() throws InterruptedException {
        // Wait if buffer is empty
        while (buffer.isEmpty()) {
            System.out.println(Thread.currentThread().getName() + 
                             ": Buffer empty, waiting...");
            wait();
        }
        
        T item = buffer.poll();
        System.out.println(Thread.currentThread().getName() + 
                         ": Consumed " + item + ", Buffer size: " + buffer.size());
        
        // Notify producers that space is available
        notifyAll();
        
        return item;
    }
    
    public synchronized int size() {
        return buffer.size();
    }
}

// Producer Thread
class Producer implements Runnable {
    private final BoundedBuffer<Integer> buffer;
    private final int itemsToProduce;
    
    public Producer(BoundedBuffer<Integer> buffer, int itemsToProduce) {
        this.buffer = buffer;
        this.itemsToProduce = itemsToProduce;
    }
    
    @Override
    public void run() {
        try {
            for (int i = 1; i <= itemsToProduce; i++) {
                buffer.put(i);
                // Simulate production time
                Thread.sleep((int) (Math.random() * 500));
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.out.println(Thread.currentThread().getName() + " interrupted!");
        }
    }
}

// Consumer Thread
class Consumer implements Runnable {
    private final BoundedBuffer<Integer> buffer;
    private final int itemsToConsume;
    
    public Consumer(BoundedBuffer<Integer> buffer, int itemsToConsume) {
        this.buffer = buffer;
        this.itemsToConsume = itemsToConsume;
    }
    
    @Override
    public void run() {
        try {
            for (int i = 0; i < itemsToConsume; i++) {
                Integer item = buffer.get();
                // Simulate consumption/processing time
                Thread.sleep((int) (Math.random() * 1000));
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
            System.out.println(Thread.currentThread().getName() + " interrupted!");
        }
    }
}

// Main class to run the demo
public class ProducerConsumerDemo {
    public static void main(String[] args) throws InterruptedException {
        // Create buffer with capacity 5
        BoundedBuffer<Integer> buffer = new BoundedBuffer<>(5);
        
        // Total items to produce and consume
        int totalItems = 10;
        
        // Create producer and consumer threads
        Thread producer1 = new Thread(new Producer(buffer, totalItems / 2), "Producer-1");
        Thread producer2 = new Thread(new Producer(buffer, totalItems / 2), "Producer-2");
        Thread consumer1 = new Thread(new Consumer(buffer, totalItems / 2), "Consumer-1");
        Thread consumer2 = new Thread(new Consumer(buffer, totalItems / 2), "Consumer-2");
        
        System.out.println("=== Producer-Consumer Demo Started ===");
        System.out.println("Buffer Capacity: 5");
        System.out.println("Producers: 2, Consumers: 2");
        System.out.println("Total Items: " + totalItems);
        System.out.println("========================================\n");
        
        // Start all threads
        producer1.start();
        producer2.start();
        consumer1.start();
        consumer2.start();
        
        // Wait for all threads to complete
        producer1.join();
        producer2.join();
        consumer1.join();
        consumer2.join();
        
        System.out.println("\n========================================");
        System.out.println("=== All threads completed! ===");
        System.out.println("Final buffer size: " + buffer.size());
    }
}
```

**Sample Output:**
```
=== Producer-Consumer Demo Started ===
Buffer Capacity: 5
Producers: 2, Consumers: 2
Total Items: 10
========================================

Producer-1: Produced 1, Buffer size: 1
Producer-2: Produced 1, Buffer size: 2
Consumer-1: Consumed 1, Buffer size: 1
Producer-1: Produced 2, Buffer size: 2
Consumer-2: Consumed 1, Buffer size: 1
Producer-2: Produced 2, Buffer size: 2
Producer-1: Produced 3, Buffer size: 3
Producer-2: Produced 3, Buffer size: 4
Producer-1: Produced 4, Buffer size: 5
Producer-2: Buffer full, waiting...
Consumer-1: Consumed 2, Buffer size: 4
Producer-2: Produced 4, Buffer size: 5
Consumer-2: Consumed 2, Buffer size: 4
Producer-1: Produced 5, Buffer size: 5
Consumer-1: Consumed 3, Buffer size: 4
Producer-2: Produced 5, Buffer size: 5
Consumer-2: Consumed 3, Buffer size: 4
Consumer-1: Consumed 4, Buffer size: 3
Consumer-2: Consumed 4, Buffer size: 2
Consumer-1: Consumed 5, Buffer size: 1
Consumer-2: Consumed 5, Buffer size: 0

========================================
=== All threads completed! ===
Final buffer size: 0
```

#### Producer-Consumer Using BlockingQueue (Modern Approach)

Java's `BlockingQueue` is the preferred way to implement producer-consumer:

```java
import java.util.concurrent.*;

public class ProducerConsumerBlockingQueue {
    public static void main(String[] args) throws InterruptedException {
        // BlockingQueue handles synchronization internally!
        BlockingQueue<Integer> queue = new ArrayBlockingQueue<>(5);
        
        // Producer
        Thread producer = new Thread(() -> {
            try {
                for (int i = 1; i <= 10; i++) {
                    System.out.println("Producing: " + i);
                    queue.put(i); // Blocks if queue is full
                    Thread.sleep(100);
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        });
        
        // Consumer
        Thread consumer = new Thread(() -> {
            try {
                for (int i = 0; i < 10; i++) {
                    Integer item = queue.take(); // Blocks if queue is empty
                    System.out.println("Consumed: " + item);
                    Thread.sleep(200);
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        });
        
        producer.start();
        consumer.start();
        
        producer.join();
        consumer.join();
        
        System.out.println("Done!");
    }
}
```

#### BlockingQueue Types

| Type | Description |
|------|-------------|
| `ArrayBlockingQueue` | Fixed-size, backed by array, fair/non-fair |
| `LinkedBlockingQueue` | Optionally bounded, backed by linked nodes |
| `PriorityBlockingQueue` | Unbounded, elements ordered by priority |
| `DelayQueue` | Elements available only after delay expires |
| `SynchronousQueue` | Zero capacity, direct handoff between threads |

---

### 4.5 wait/notify vs sleep Comparison

| Aspect | wait() | sleep() |
|--------|--------|--------|
| **Class** | Object | Thread |
| **Lock Release** | Releases lock | Does NOT release lock |
| **Must Own Lock** | Yes (IllegalMonitorStateException otherwise) | No |
| **Wake Condition** | notify()/notifyAll()/interrupt | Time expires/interrupt |
| **Use Case** | Inter-thread communication | Pausing execution |
| **State** | WAITING | TIMED_WAITING |

```java
public class WaitVsSleepDemo {
    private static final Object lock = new Object();
    
    public static void main(String[] args) {
        // wait() releases lock
        Thread waitThread = new Thread(() -> {
            synchronized (lock) {
                System.out.println("Wait thread: Acquired lock");
                try {
                    lock.wait(2000); // Releases lock!
                } catch (InterruptedException e) {}
                System.out.println("Wait thread: Woke up");
            }
        });
        
        // sleep() does NOT release lock
        Thread sleepThread = new Thread(() -> {
            synchronized (lock) {
                System.out.println("Sleep thread: Acquired lock");
                try {
                    Thread.sleep(2000); // Does NOT release lock!
                } catch (InterruptedException e) {}
                System.out.println("Sleep thread: Woke up");
            }
        });
        
        waitThread.start();
        // sleepThread.start();
    }
}
```

---

### 4.6 Common Patterns and Best Practices

#### Guarded Block Pattern

```java
public class GuardedBlock {
    private Object result = null;
    
    public synchronized void setResult(Object result) {
        this.result = result;
        notifyAll();
    }
    
    public synchronized Object getResult() {
        // Standard guarded block pattern
        while (result == null) {
            try {
                wait();
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                return null;
            }
        }
        return result;
    }
}
```

#### Best Practices Summary

```
┌─────────────────────────────────────────────────────────────────────┐
│              INTER-THREAD COMMUNICATION BEST PRACTICES              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. Always call wait() in a while loop, not if                      │
│     while (condition) { wait(); } // CORRECT                        │
│     if (condition) { wait(); }    // WRONG                          │
│                                                                      │
│  2. Always call wait()/notify() inside synchronized block           │
│     synchronized (obj) {                                            │
│         obj.wait();   // OK                                         │
│         obj.notify(); // OK                                         │
│     }                                                                │
│                                                                      │
│  3. Prefer notifyAll() over notify()                                │
│     • Safer, avoids lost wakeup problem                             │
│     • Use notify() only when:                                       │
│       - All waiters wait for same condition                         │
│       - Only one waiter needs to proceed                            │
│                                                                      │
│  4. Hold locks for shortest time possible                           │
│     • Don't do I/O or expensive operations in synchronized blocks   │
│                                                                      │
│  5. Prefer higher-level concurrency utilities                       │
│     • BlockingQueue for producer-consumer                           │
│     • CountDownLatch for waiting on events                          │
│     • Semaphore for resource pooling                                │
│     • CyclicBarrier for phased computations                         │
│                                                                      │
│  6. Document synchronization policy                                 │
│     • Which lock protects which state                               │
│     • Thread safety guarantees                                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 5. Advanced Concurrency (java.util.concurrent)

The `java.util.concurrent` package, introduced in Java 5, provides high-level concurrency utilities that are more powerful and easier to use than low-level primitives like `synchronized` and `wait/notify`.

### 5.1 ExecutorService

`ExecutorService` is a higher-level replacement for working with threads directly. It manages a pool of threads and provides methods to submit tasks for execution.

```
┌─────────────────────────────────────────────────────────────────────┐
│                     EXECUTOR FRAMEWORK HIERARCHY                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│                        ┌─────────────────┐                          │
│                        │    Executor     │ (Interface)              │
│                        │   execute(r)    │                          │
│                        └────────┬────────┘                          │
│                                 │                                    │
│                                 ▼                                    │
│                     ┌─────────────────────┐                         │
│                     │  ExecutorService    │ (Interface)             │
│                     │  submit(), shutdown │                         │
│                     │  invokeAll/Any()    │                         │
│                     └──────────┬──────────┘                         │
│                                │                                     │
│              ┌─────────────────┼─────────────────┐                  │
│              │                 │                 │                  │
│              ▼                 ▼                 ▼                  │
│    ┌─────────────────┐ ┌─────────────┐ ┌─────────────────────┐     │
│    │ThreadPoolExecutor│ │ScheduledExec│ │ ForkJoinPool        │     │
│    │                 │ │utorService  │ │                     │     │
│    └─────────────────┘ └─────────────┘ └─────────────────────┘     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Basic ExecutorService Usage

```java
import java.util.concurrent.*;

public class ExecutorServiceDemo {
    public static void main(String[] args) {
        // Create ExecutorService with fixed thread pool
        ExecutorService executor = Executors.newFixedThreadPool(3);
        
        try {
            // Submit Runnable tasks
            for (int i = 1; i <= 5; i++) {
                final int taskId = i;
                executor.execute(() -> {
                    System.out.println("Task " + taskId + " running on " + 
                                     Thread.currentThread().getName());
                    try {
                        Thread.sleep(1000);
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                    }
                });
            }
            
            // Submit Callable task with return value
            Future<String> future = executor.submit(() -> {
                Thread.sleep(500);
                return "Task completed!";
            });
            
            System.out.println("Future result: " + future.get());
            
        } catch (Exception e) {
            e.printStackTrace();
        } finally {
            // Always shutdown the executor
            executor.shutdown();
            try {
                if (!executor.awaitTermination(60, TimeUnit.SECONDS)) {
                    executor.shutdownNow();
                }
            } catch (InterruptedException e) {
                executor.shutdownNow();
            }
        }
    }
}
```

#### ExecutorService Methods

| Method | Description |
|--------|-------------|
| `execute(Runnable)` | Executes task, no return value |
| `submit(Runnable)` | Submits task, returns Future<?> |
| `submit(Callable<T>)` | Submits task, returns Future<T> |
| `invokeAll(Collection<Callable>)` | Executes all tasks, waits for completion |
| `invokeAny(Collection<Callable>)` | Returns result of first completed task |
| `shutdown()` | Initiates orderly shutdown |
| `shutdownNow()` | Attempts to stop all tasks immediately |
| `awaitTermination(timeout, unit)` | Blocks until all tasks complete |
| `isShutdown()` | Returns true if shutdown has been called |
| `isTerminated()` | Returns true if all tasks have completed |

---

### 5.2 Thread Pools

Thread pools manage a pool of worker threads, reducing the overhead of creating new threads for each task.

```
┌─────────────────────────────────────────────────────────────────────┐
│                         THREAD POOL CONCEPT                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   Task Queue                        Thread Pool                      │
│   ┌────────────┐                   ┌─────────────────────┐          │
│   │ Task 1     │ ───────────────►  │ ┌───┐ ┌───┐ ┌───┐  │          │
│   │ Task 2     │                   │ │ T1│ │ T2│ │ T3│  │          │
│   │ Task 3     │     Workers       │ └─┬─┘ └─┬─┘ └─┬─┘  │          │
│   │ Task 4     │    pick tasks     │   │     │     │    │          │
│   │ Task 5     │                   │   ▼     ▼     ▼    │          │
│   │    ...     │                   │ Execute tasks      │          │
│   └────────────┘                   │ concurrently       │          │
│                                    └─────────────────────┘          │
│                                                                      │
│   Benefits:                                                         │
│   • Reuse threads (avoid creation overhead)                         │
│   • Control max concurrency                                         │
│   • Manage resources efficiently                                    │
│   • Queue excess tasks                                              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Types of Thread Pools

```java
import java.util.concurrent.*;

public class ThreadPoolTypesDemo {
    
    public static void main(String[] args) {
        
        // 1. Fixed Thread Pool - Fixed number of threads
        ExecutorService fixedPool = Executors.newFixedThreadPool(4);
        // Use when: known number of concurrent tasks
        // Threads: exactly n threads
        // Queue: unbounded LinkedBlockingQueue
        
        // 2. Cached Thread Pool - Creates threads as needed
        ExecutorService cachedPool = Executors.newCachedThreadPool();
        // Use when: many short-lived tasks
        // Threads: 0 to Integer.MAX_VALUE
        // Queue: SynchronousQueue (direct handoff)
        // Idle threads kept for 60 seconds
        
        // 3. Single Thread Executor - One thread
        ExecutorService singlePool = Executors.newSingleThreadExecutor();
        // Use when: sequential execution required
        // Threads: exactly 1 thread
        // Queue: unbounded LinkedBlockingQueue
        
        // 4. Scheduled Thread Pool - For scheduled/periodic tasks
        ScheduledExecutorService scheduledPool = Executors.newScheduledThreadPool(2);
        
        // Schedule task to run after 5 seconds
        scheduledPool.schedule(() -> {
            System.out.println("Delayed task executed!");
        }, 5, TimeUnit.SECONDS);
        
        // Schedule task to run every 3 seconds
        scheduledPool.scheduleAtFixedRate(() -> {
            System.out.println("Periodic task: " + System.currentTimeMillis());
        }, 0, 3, TimeUnit.SECONDS);
        
        // Schedule with fixed delay between executions
        scheduledPool.scheduleWithFixedDelay(() -> {
            System.out.println("Fixed delay task");
        }, 0, 2, TimeUnit.SECONDS);
        
        // 5. Work Stealing Pool (Java 8+) - For parallel processing
        ExecutorService workStealingPool = Executors.newWorkStealingPool();
        // Use when: CPU-intensive tasks with varying execution times
        // Threads: number of available processors
        // Uses ForkJoinPool internally
        
        // Don't forget to shutdown!
        fixedPool.shutdown();
        cachedPool.shutdown();
        singlePool.shutdown();
        // scheduledPool.shutdown(); // Keep for demo
        workStealingPool.shutdown();
    }
}
```

#### Custom ThreadPoolExecutor

```java
import java.util.concurrent.*;

public class CustomThreadPoolDemo {
    public static void main(String[] args) {
        // Create custom ThreadPoolExecutor
        ThreadPoolExecutor executor = new ThreadPoolExecutor(
            2,                      // corePoolSize
            4,                      // maximumPoolSize
            60L,                    // keepAliveTime
            TimeUnit.SECONDS,       // time unit
            new ArrayBlockingQueue<>(10),  // work queue
            new ThreadFactory() {   // thread factory
                private int count = 0;
                @Override
                public Thread newThread(Runnable r) {
                    return new Thread(r, "CustomThread-" + count++);
                }
            },
            new ThreadPoolExecutor.CallerRunsPolicy()  // rejection policy
        );
        
        // Submit tasks
        for (int i = 0; i < 20; i++) {
            final int taskId = i;
            executor.execute(() -> {
                System.out.println("Task " + taskId + " on " + 
                                 Thread.currentThread().getName());
                try { Thread.sleep(500); } catch (InterruptedException e) {}
            });
        }
        
        // Monitor pool status
        System.out.println("Pool Size: " + executor.getPoolSize());
        System.out.println("Active Count: " + executor.getActiveCount());
        System.out.println("Completed Tasks: " + executor.getCompletedTaskCount());
        System.out.println("Queue Size: " + executor.getQueue().size());
        
        executor.shutdown();
    }
}
```

#### Rejection Policies

| Policy | Behavior |
|--------|----------|
| `AbortPolicy` (default) | Throws RejectedExecutionException |
| `CallerRunsPolicy` | Task runs in the submitting thread |
| `DiscardPolicy` | Silently discards the task |
| `DiscardOldestPolicy` | Discards oldest task in queue |

```
┌─────────────────────────────────────────────────────────────────────┐
│                    THREAD POOL SIZING GUIDE                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  CPU-Bound Tasks (calculations, processing):                        │
│  └── Pool Size = Number of CPUs + 1                                 │
│                                                                      │
│  I/O-Bound Tasks (network, file, database):                         │
│  └── Pool Size = Number of CPUs × (1 + Wait Time / Service Time)    │
│                                                                      │
│  Example for 4-core CPU:                                            │
│  • CPU-bound: 4 + 1 = 5 threads                                     │
│  • I/O-bound (80% wait): 4 × (1 + 0.8/0.2) = 4 × 5 = 20 threads    │
│                                                                      │
│  Get CPU count: Runtime.getRuntime().availableProcessors()          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

### 5.3 Callable vs Runnable

```java
import java.util.concurrent.*;

public class CallableVsRunnableDemo {
    public static void main(String[] args) throws Exception {
        ExecutorService executor = Executors.newFixedThreadPool(2);
        
        // Runnable - No return value, no checked exceptions
        Runnable runnable = () -> {
            System.out.println("Runnable executing");
            // Cannot return value
            // Cannot throw checked exception
        };
        
        // Callable - Returns value, can throw exceptions
        Callable<Integer> callable = () -> {
            System.out.println("Callable executing");
            if (Math.random() < 0.3) {
                throw new Exception("Random failure!"); // Can throw!
            }
            return 42; // Returns value!
        };
        
        // Execute Runnable
        executor.execute(runnable);
        Future<?> runnableFuture = executor.submit(runnable);
        
        // Execute Callable
        Future<Integer> callableFuture = executor.submit(callable);
        
        try {
            Integer result = callableFuture.get(); // Blocks until complete
            System.out.println("Callable result: " + result);
        } catch (ExecutionException e) {
            System.out.println("Callable threw: " + e.getCause().getMessage());
        }
        
        executor.shutdown();
    }
}
```

| Feature | Runnable | Callable |
|---------|----------|----------|
| **Return Type** | `void` | Generic `V` |
| **Method** | `run()` | `call()` |
| **Exceptions** | Cannot throw checked | Can throw Exception |
| **Introduced** | Java 1.0 | Java 5 |
| **Use With** | Thread, ExecutorService | ExecutorService only |

---

### 5.4 Future

`Future` represents the result of an asynchronous computation.

```java
import java.util.concurrent.*;
import java.util.*;

public class FutureDemo {
    public static void main(String[] args) {
        ExecutorService executor = Executors.newFixedThreadPool(3);
        
        // Submit task and get Future
        Future<String> future = executor.submit(() -> {
            Thread.sleep(2000);
            return "Hello from Future!";
        });
        
        // Check if done (non-blocking)
        System.out.println("Is done? " + future.isDone());     // false
        System.out.println("Is cancelled? " + future.isCancelled()); // false
        
        try {
            // Get with timeout
            String result = future.get(5, TimeUnit.SECONDS);
            System.out.println("Result: " + result);
            
        } catch (TimeoutException e) {
            System.out.println("Task took too long!");
            future.cancel(true); // Try to cancel
        } catch (InterruptedException | ExecutionException e) {
            e.printStackTrace();
        }
        
        System.out.println("Is done? " + future.isDone());     // true
        
        // Multiple Futures
        List<Future<Integer>> futures = new ArrayList<>();
        for (int i = 0; i < 5; i++) {
            final int num = i;
            futures.add(executor.submit(() -> {
                Thread.sleep((int)(Math.random() * 1000));
                return num * num;
            }));
        }
        
        // Collect results
        for (int i = 0; i < futures.size(); i++) {
            try {
                System.out.println("Future " + i + " result: " + futures.get(i).get());
            } catch (Exception e) {
                e.printStackTrace();
            }
        }
        
        executor.shutdown();
    }
}
```

#### Future Methods

| Method | Description |
|--------|-------------|
| `get()` | Blocks until result available, returns result |
| `get(timeout, unit)` | Blocks with timeout |
| `isDone()` | Returns true if task completed |
| `isCancelled()` | Returns true if task was cancelled |
| `cancel(mayInterrupt)` | Attempts to cancel task |

---

### 5.5 CountDownLatch

`CountDownLatch` allows one or more threads to wait until a set of operations being performed by other threads completes.

```
┌─────────────────────────────────────────────────────────────────────┐
│                       COUNTDOWNLATCH                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Initial: CountDownLatch(3)         count = 3                       │
│                                                                      │
│  Thread A: countDown() ─────────►   count = 2                       │
│  Thread B: countDown() ─────────►   count = 1                       │
│  Thread C: countDown() ─────────►   count = 0                       │
│                                        │                             │
│  Main Thread: await() ◄────────────────┘ (unblocked when count=0)   │
│                                                                      │
│  Key Points:                                                        │
│  • Count can only decrease (no reset)                               │
│  • Single-use synchronizer                                          │
│  • When count reaches 0, all waiting threads released              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.*;

public class CountDownLatchDemo {
    public static void main(String[] args) throws InterruptedException {
        int numberOfServices = 3;
        CountDownLatch latch = new CountDownLatch(numberOfServices);
        
        // Simulate starting multiple services
        ExecutorService executor = Executors.newFixedThreadPool(numberOfServices);
        
        String[] services = {"Database", "Cache", "MessageQueue"};
        
        for (String service : services) {
            executor.submit(() -> {
                try {
                    System.out.println(service + " starting...");
                    Thread.sleep((int)(Math.random() * 3000)); // Simulate startup
                    System.out.println(service + " started!");
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    latch.countDown(); // Decrease count
                    System.out.println("Latch count: " + latch.getCount());
                }
            });
        }
        
        System.out.println("\nMain: Waiting for all services to start...");
        latch.await(); // Block until count reaches 0
        System.out.println("\nMain: All services started! Application ready.");
        
        executor.shutdown();
    }
}
```

**Output:**
```
Main: Waiting for all services to start...
Database starting...
Cache starting...
MessageQueue starting...
Cache started!
Latch count: 2
Database started!
Latch count: 1
MessageQueue started!
Latch count: 0

Main: All services started! Application ready.
```

#### Use Cases for CountDownLatch

1. **Application startup**: Wait for all services to initialize
2. **Parallel processing**: Wait for all workers to complete
3. **Testing**: Ensure all threads start simultaneously

```java
// Testing use case: Start all threads at same time
public class SimultaneousStartDemo {
    public static void main(String[] args) throws InterruptedException {
        int threadCount = 5;
        CountDownLatch readyLatch = new CountDownLatch(threadCount);
        CountDownLatch startLatch = new CountDownLatch(1);
        CountDownLatch doneLatch = new CountDownLatch(threadCount);
        
        for (int i = 0; i < threadCount; i++) {
            final int id = i;
            new Thread(() -> {
                try {
                    readyLatch.countDown(); // Signal ready
                    startLatch.await();     // Wait for start signal
                    
                    System.out.println("Thread " + id + " started at " + 
                                     System.currentTimeMillis());
                    // Do work...
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    doneLatch.countDown();
                }
            }).start();
        }
        
        readyLatch.await();  // Wait for all threads to be ready
        System.out.println("All threads ready, starting!");
        startLatch.countDown(); // Release all threads simultaneously
        
        doneLatch.await();   // Wait for all to complete
        System.out.println("All threads completed!");
    }
}
```

---

### 5.6 CyclicBarrier

`CyclicBarrier` allows a set of threads to all wait for each other to reach a common barrier point. Unlike CountDownLatch, it can be reused.

```
┌─────────────────────────────────────────────────────────────────────┐
│                        CYCLICBARRIER                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Phase 1:                                                           │
│  Thread A ───────► await() ─┐                                       │
│  Thread B ───────► await() ─┼───► Barrier ───► All proceed         │
│  Thread C ───────► await() ─┘     (released)                        │
│                                                                      │
│  Phase 2 (barrier resets automatically):                            │
│  Thread A ───────► await() ─┐                                       │
│  Thread B ───────► await() ─┼───► Barrier ───► All proceed         │
│  Thread C ───────► await() ─┘     (released)                        │
│                                                                      │
│  Key Features:                                                      │
│  • Reusable (cyclic) - resets after all threads arrive             │
│  • Optional barrier action runs when barrier trips                  │
│  • Threads wait for each other                                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.*;

public class CyclicBarrierDemo {
    public static void main(String[] args) {
        int parties = 3;
        
        // Barrier action runs when all threads reach barrier
        Runnable barrierAction = () -> {
            System.out.println("\n>>> Barrier reached! All threads synchronized. <<<\n");
        };
        
        CyclicBarrier barrier = new CyclicBarrier(parties, barrierAction);
        
        // Simulate phases of computation
        for (int i = 0; i < parties; i++) {
            final int threadId = i;
            new Thread(() -> {
                try {
                    for (int phase = 1; phase <= 3; phase++) {
                        // Do phase work
                        System.out.println("Thread " + threadId + 
                                         " completed phase " + phase);
                        Thread.sleep((int)(Math.random() * 1000));
                        
                        // Wait for other threads
                        int arrival = barrier.await();
                        System.out.println("Thread " + threadId + 
                                         " was arrival #" + arrival);
                    }
                } catch (InterruptedException | BrokenBarrierException e) {
                    e.printStackTrace();
                }
            }, "Thread-" + i).start();
        }
    }
}
```

**Output:**
```
Thread 0 completed phase 1
Thread 1 completed phase 1
Thread 2 completed phase 1

>>> Barrier reached! All threads synchronized. <<<

Thread 2 was arrival #0
Thread 0 was arrival #2
Thread 1 was arrival #1
Thread 0 completed phase 2
Thread 2 completed phase 2
Thread 1 completed phase 2

>>> Barrier reached! All threads synchronized. <<<

...(continues for phase 3)
```

#### CountDownLatch vs CyclicBarrier

| Feature | CountDownLatch | CyclicBarrier |
|---------|----------------|---------------|
| **Reusable** | No (one-time use) | Yes (resets automatically) |
| **Count Direction** | Counts down to zero | Counts up to parties |
| **Wait Mechanism** | One or more threads wait | All threads wait for each other |
| **Barrier Action** | None | Optional runnable when barrier trips |
| **Use Case** | Wait for events | Synchronize phases of computation |
| **Reset** | Not possible | Automatic or manual reset() |

---

### 5.7 Semaphore

`Semaphore` maintains a set of permits. Threads can acquire permits (blocking if none available) and release permits.

```
┌─────────────────────────────────────────────────────────────────────┐
│                          SEMAPHORE                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Semaphore(3) - 3 permits available                                 │
│                                                                      │
│  Permits: [■][■][■]                                                 │
│                                                                      │
│  Thread A: acquire() → [□][■][■] (got permit)                       │
│  Thread B: acquire() → [□][□][■] (got permit)                       │
│  Thread C: acquire() → [□][□][□] (got permit)                       │
│  Thread D: acquire() → BLOCKED (no permits!)                        │
│                                                                      │
│  Thread A: release() → [■][□][□] (returned permit)                  │
│  Thread D: (unblocked) → [□][□][□] (got permit)                     │
│                                                                      │
│  Use Cases:                                                         │
│  • Connection pool limiting                                         │
│  • Rate limiting                                                    │
│  • Resource access control                                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.*;

public class SemaphoreDemo {
    public static void main(String[] args) {
        // Only 3 threads can access resource simultaneously
        Semaphore semaphore = new Semaphore(3);
        
        // Simulate 10 threads trying to access limited resource
        ExecutorService executor = Executors.newFixedThreadPool(10);
        
        for (int i = 1; i <= 10; i++) {
            final int userId = i;
            executor.submit(() -> {
                try {
                    System.out.println("User " + userId + " waiting for permit...");
                    semaphore.acquire();
                    
                    System.out.println("User " + userId + " ACQUIRED permit. " +
                                     "Available permits: " + semaphore.availablePermits());
                    
                    // Simulate using resource
                    Thread.sleep(2000);
                    
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                } finally {
                    semaphore.release();
                    System.out.println("User " + userId + " RELEASED permit. " +
                                     "Available permits: " + semaphore.availablePermits());
                }
            });
        }
        
        executor.shutdown();
    }
}
```

#### Practical Example: Connection Pool

```java
import java.util.concurrent.*;
import java.util.*;

public class ConnectionPool {
    private final Semaphore semaphore;
    private final Queue<Connection> pool;
    
    public ConnectionPool(int poolSize) {
        this.semaphore = new Semaphore(poolSize);
        this.pool = new LinkedList<>();
        
        // Initialize connections
        for (int i = 0; i < poolSize; i++) {
            pool.add(new Connection("Connection-" + i));
        }
    }
    
    public Connection acquire() throws InterruptedException {
        semaphore.acquire(); // Block if no connections available
        synchronized (pool) {
            return pool.poll();
        }
    }
    
    public Connection tryAcquire(long timeout, TimeUnit unit) 
            throws InterruptedException {
        if (semaphore.tryAcquire(timeout, unit)) {
            synchronized (pool) {
                return pool.poll();
            }
        }
        return null; // Timeout - no connection available
    }
    
    public void release(Connection conn) {
        synchronized (pool) {
            pool.offer(conn);
        }
        semaphore.release();
    }
    
    public int availableConnections() {
        return semaphore.availablePermits();
    }
    
    // Dummy Connection class
    static class Connection {
        private final String name;
        Connection(String name) { this.name = name; }
        @Override public String toString() { return name; }
    }
}
```

#### Semaphore Methods

| Method | Description |
|--------|-------------|
| `acquire()` | Acquires permit, blocking if necessary |
| `acquire(n)` | Acquires n permits |
| `tryAcquire()` | Acquires permit if available, non-blocking |
| `tryAcquire(timeout, unit)` | Acquires with timeout |
| `release()` | Releases a permit |
| `release(n)` | Releases n permits |
| `availablePermits()` | Returns available permits |
| `drainPermits()` | Acquires all available permits |

---

### 5.8 ReentrantLock

`ReentrantLock` is an alternative to `synchronized` with additional features like fairness, interruptible lock acquisition, and try-lock.

```
┌─────────────────────────────────────────────────────────────────────┐
│                       REENTRANTLOCK                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  synchronized                      ReentrantLock                     │
│  ─────────────                     ──────────────                    │
│  synchronized (obj) {              lock.lock();                      │
│      // critical section           try {                             │
│  }                                     // critical section           │
│                                    } finally {                       │
│                                        lock.unlock();                │
│                                    }                                 │
│                                                                      │
│  Additional Features:                                                │
│  • Fairness policy (FIFO ordering)                                  │
│  • Interruptible lock acquisition                                   │
│  • Non-blocking tryLock()                                           │
│  • Multiple Condition objects                                       │
│  • Lock query methods                                               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.locks.*;

public class ReentrantLockDemo {
    private final ReentrantLock lock = new ReentrantLock();
    private int count = 0;
    
    // Basic usage
    public void increment() {
        lock.lock();
        try {
            count++;
        } finally {
            lock.unlock(); // ALWAYS unlock in finally!
        }
    }
    
    // Try lock - non-blocking
    public boolean tryIncrement() {
        if (lock.tryLock()) {
            try {
                count++;
                return true;
            } finally {
                lock.unlock();
            }
        }
        return false; // Could not acquire lock
    }
    
    // Try lock with timeout
    public boolean tryIncrementWithTimeout() throws InterruptedException {
        if (lock.tryLock(1, TimeUnit.SECONDS)) {
            try {
                count++;
                return true;
            } finally {
                lock.unlock();
            }
        }
        return false;
    }
    
    // Interruptible lock acquisition
    public void interruptibleIncrement() throws InterruptedException {
        lock.lockInterruptibly(); // Can be interrupted while waiting
        try {
            count++;
        } finally {
            lock.unlock();
        }
    }
    
    // Query methods
    public void printLockInfo() {
        System.out.println("Is locked: " + lock.isLocked());
        System.out.println("Is held by current: " + lock.isHeldByCurrentThread());
        System.out.println("Hold count: " + lock.getHoldCount());
        System.out.println("Has queued threads: " + lock.hasQueuedThreads());
        System.out.println("Queue length: " + lock.getQueueLength());
    }
}
```

#### Fair vs Non-Fair Lock

```java
import java.util.concurrent.locks.*;

public class FairLockDemo {
    public static void main(String[] args) {
        // Non-fair lock (default) - may allow barging
        ReentrantLock nonFairLock = new ReentrantLock();
        
        // Fair lock - FIFO ordering
        ReentrantLock fairLock = new ReentrantLock(true);
        
        // Fair lock ensures threads acquire lock in order they requested
        // But has lower throughput due to overhead
    }
}
```

#### Condition Variables

```java
import java.util.concurrent.locks.*;
import java.util.*;

public class BoundedBufferWithLock<T> {
    private final Lock lock = new ReentrantLock();
    private final Condition notFull = lock.newCondition();
    private final Condition notEmpty = lock.newCondition();
    
    private final Queue<T> buffer = new LinkedList<>();
    private final int capacity;
    
    public BoundedBufferWithLock(int capacity) {
        this.capacity = capacity;
    }
    
    public void put(T item) throws InterruptedException {
        lock.lock();
        try {
            while (buffer.size() == capacity) {
                notFull.await(); // Wait until not full
            }
            buffer.add(item);
            notEmpty.signal(); // Signal one waiting consumer
        } finally {
            lock.unlock();
        }
    }
    
    public T take() throws InterruptedException {
        lock.lock();
        try {
            while (buffer.isEmpty()) {
                notEmpty.await(); // Wait until not empty
            }
            T item = buffer.poll();
            notFull.signal(); // Signal one waiting producer
            return item;
        } finally {
            lock.unlock();
        }
    }
}
```

#### synchronized vs ReentrantLock

| Feature | synchronized | ReentrantLock |
|---------|--------------|---------------|
| **Ease of use** | Simpler | More verbose |
| **Release** | Automatic | Must call unlock() |
| **Fairness** | No | Optional |
| **Try lock** | No | Yes |
| **Interruptible** | No | Yes |
| **Multiple conditions** | No (single wait set) | Yes |
| **Lock timeout** | No | Yes |
| **Performance** | Similar | Similar (slightly faster) |

---

### 5.9 ReadWriteLock

`ReadWriteLock` maintains a pair of locks: one for read-only operations and one for writes. Multiple readers can access simultaneously, but writers have exclusive access.

```
┌─────────────────────────────────────────────────────────────────────┐
│                       READWRITELOCK                                  │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Read Lock (Shared):                                                │
│  • Multiple readers can hold simultaneously                         │
│  • Blocked if write lock is held                                    │
│                                                                      │
│  Write Lock (Exclusive):                                            │
│  • Only one writer at a time                                        │
│  • Blocked if any read/write lock is held                          │
│                                                                      │
│  Scenario: Cache with mostly reads                                  │
│                                                                      │
│  Time ─────►                                                        │
│  Reader 1: ████████                                                 │
│  Reader 2:    ████████████                                          │
│  Reader 3:      ██████                                              │
│  Writer 1:              ░░░░ (waits) ████                           │
│  Reader 4:                            ░░░░ (waits) ████             │
│                                                                      │
│  ████ = holding lock, ░░░░ = waiting                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.locks.*;
import java.util.*;

public class ReadWriteLockDemo {
    private final ReadWriteLock rwLock = new ReentrantReadWriteLock();
    private final Lock readLock = rwLock.readLock();
    private final Lock writeLock = rwLock.writeLock();
    
    private Map<String, String> cache = new HashMap<>();
    
    // Multiple threads can read simultaneously
    public String get(String key) {
        readLock.lock();
        try {
            System.out.println(Thread.currentThread().getName() + 
                             " reading key: " + key);
            return cache.get(key);
        } finally {
            readLock.unlock();
        }
    }
    
    // Only one thread can write
    public void put(String key, String value) {
        writeLock.lock();
        try {
            System.out.println(Thread.currentThread().getName() + 
                             " writing key: " + key);
            cache.put(key, value);
        } finally {
            writeLock.unlock();
        }
    }
    
    // Read-then-write pattern (lock upgrading not supported!)
    public String getOrCompute(String key, java.util.function.Function<String, String> compute) {
        readLock.lock();
        try {
            String value = cache.get(key);
            if (value != null) {
                return value;
            }
        } finally {
            readLock.unlock();
        }
        
        // Must release read lock before acquiring write lock
        writeLock.lock();
        try {
            // Double-check (another thread might have computed)
            String value = cache.get(key);
            if (value == null) {
                value = compute.apply(key);
                cache.put(key, value);
            }
            return value;
        } finally {
            writeLock.unlock();
        }
    }
    
    public static void main(String[] args) {
        ReadWriteLockDemo demo = new ReadWriteLockDemo();
        
        // Pre-populate cache
        demo.put("key1", "value1");
        demo.put("key2", "value2");
        
        // Multiple readers
        for (int i = 0; i < 5; i++) {
            new Thread(() -> {
                System.out.println("Read: " + demo.get("key1"));
            }, "Reader-" + i).start();
        }
        
        // Single writer
        new Thread(() -> {
            demo.put("key3", "value3");
        }, "Writer-1").start();
    }
}
```

#### StampedLock (Java 8+)

`StampedLock` is a more sophisticated alternative with optimistic read support:

```java
import java.util.concurrent.locks.*;

public class StampedLockDemo {
    private final StampedLock lock = new StampedLock();
    private double x, y;
    
    // Exclusive write
    public void move(double deltaX, double deltaY) {
        long stamp = lock.writeLock();
        try {
            x += deltaX;
            y += deltaY;
        } finally {
            lock.unlockWrite(stamp);
        }
    }
    
    // Optimistic read - no locking overhead!
    public double distanceFromOrigin() {
        long stamp = lock.tryOptimisticRead(); // Non-blocking!
        double currentX = x, currentY = y;
        
        if (!lock.validate(stamp)) {
            // Optimistic read failed, fall back to read lock
            stamp = lock.readLock();
            try {
                currentX = x;
                currentY = y;
            } finally {
                lock.unlockRead(stamp);
            }
        }
        
        return Math.sqrt(currentX * currentX + currentY * currentY);
    }
}
```

---

### 5.10 ConcurrentHashMap

`ConcurrentHashMap` is a thread-safe HashMap with high concurrency support.

```
┌─────────────────────────────────────────────────────────────────────┐
│                     CONCURRENTHASHMAP                                │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Java 7 (Segment-based):                                            │
│  ┌──────────┬──────────┬──────────┬──────────┐                      │
│  │ Segment 0│ Segment 1│ Segment 2│ Segment 3│                      │
│  │  [Lock]  │  [Lock]  │  [Lock]  │  [Lock]  │                      │
│  │  bucket  │  bucket  │  bucket  │  bucket  │                      │
│  │  bucket  │  bucket  │  bucket  │  bucket  │                      │
│  └──────────┴──────────┴──────────┴──────────┘                      │
│  Each segment has its own lock (16 by default)                      │
│                                                                      │
│  Java 8+ (Node-based with CAS):                                     │
│  ┌──────────────────────────────────────────┐                       │
│  │ [Node] [Node] [Node] [Node] [Node] ...   │                       │
│  │   │       │                    │         │                       │
│  │   ↓       ↓                    ↓         │                       │
│  │ [Node] [Tree]               [Node]       │                       │
│  │   │     / \                   │          │                       │
│  │   ↓    /   \                  ↓          │                       │
│  │        ...                              │                       │
│  └──────────────────────────────────────────┘                       │
│  Uses CAS for head, synchronized for mutations                      │
│  Converts to tree when bucket > 8 entries                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.*;
import java.util.*;

public class ConcurrentHashMapDemo {
    public static void main(String[] args) throws InterruptedException {
        ConcurrentHashMap<String, Integer> map = new ConcurrentHashMap<>();
        
        // Basic operations (thread-safe)
        map.put("a", 1);
        map.put("b", 2);
        Integer value = map.get("a");
        
        // Atomic operations
        map.putIfAbsent("c", 3);  // Only if not present
        map.remove("a", 1);       // Only if value matches
        map.replace("b", 2, 20);  // Only if old value matches
        
        // Compute operations (atomic)
        map.compute("d", (k, v) -> (v == null) ? 1 : v + 1);
        map.computeIfAbsent("e", k -> k.length());
        map.computeIfPresent("b", (k, v) -> v * 2);
        map.merge("f", 1, Integer::sum);
        
        // Concurrent-safe iteration
        map.forEach((k, v) -> System.out.println(k + "=" + v));
        
        // Bulk operations (Java 8+)
        map.forEach(2, (k, v) -> System.out.println(k + "=" + v)); // parallelism threshold
        
        long sum = map.reduceValues(2, Integer::sum);
        System.out.println("Sum of values: " + sum);
        
        String keys = map.reduceKeys(2, (k1, k2) -> k1 + "," + k2);
        System.out.println("All keys: " + keys);
        
        // Search (returns first match)
        String found = map.search(2, (k, v) -> v > 10 ? k : null);
        System.out.println("Key with value > 10: " + found);
    }
}
```

#### Thread-Safe Counter Pattern

```java
import java.util.concurrent.*;

public class WordCounter {
    private final ConcurrentHashMap<String, Long> wordCounts = new ConcurrentHashMap<>();
    
    // Thread-safe increment
    public void addWord(String word) {
        wordCounts.merge(word, 1L, Long::sum);
        // Equivalent to:
        // wordCounts.compute(word, (k, v) -> (v == null) ? 1L : v + 1L);
    }
    
    public long getCount(String word) {
        return wordCounts.getOrDefault(word, 0L);
    }
    
    public static void main(String[] args) throws InterruptedException {
        WordCounter counter = new WordCounter();
        String text = "hello world hello java world java java";
        
        ExecutorService executor = Executors.newFixedThreadPool(4);
        
        for (String word : text.split(" ")) {
            executor.submit(() -> counter.addWord(word));
        }
        
        executor.shutdown();
        executor.awaitTermination(1, TimeUnit.SECONDS);
        
        System.out.println("hello: " + counter.getCount("hello")); // 2
        System.out.println("java: " + counter.getCount("java"));   // 3
        System.out.println("world: " + counter.getCount("world")); // 2
    }
}
```

#### Thread-Safe Collections Comparison

| Collection | Thread-Safety Mechanism | Performance |
|------------|-------------------------|-------------|
| `Hashtable` | Synchronized (whole map) | Poor |
| `Collections.synchronizedMap()` | Synchronized wrapper | Poor |
| `ConcurrentHashMap` | CAS + fine-grained locks | Excellent |
| `ConcurrentSkipListMap` | Lock-free (CAS) | Good (sorted) |

```
┌─────────────────────────────────────────────────────────────────────┐
│              CONCURRENT COLLECTIONS CHEAT SHEET                      │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Maps:                                                              │
│  • ConcurrentHashMap      - High-concurrency hash map               │
│  • ConcurrentSkipListMap  - Sorted, navigable map                   │
│                                                                      │
│  Sets:                                                              │
│  • ConcurrentSkipListSet  - Sorted concurrent set                   │
│  • ConcurrentHashMap.newKeySet() - Concurrent hash set              │
│                                                                      │
│  Queues:                                                            │
│  • ConcurrentLinkedQueue  - Unbounded, non-blocking                 │
│  • ConcurrentLinkedDeque  - Unbounded deque, non-blocking           │
│  • ArrayBlockingQueue     - Bounded, blocking                       │
│  • LinkedBlockingQueue    - Optionally bounded, blocking            │
│  • PriorityBlockingQueue  - Priority-based, blocking                │
│  • DelayQueue             - Delayed elements                        │
│  • SynchronousQueue       - Zero-capacity, direct handoff           │
│  • LinkedTransferQueue    - Transfer queue                          │
│                                                                      │
│  Copy-on-Write (for mostly-read scenarios):                        │
│  • CopyOnWriteArrayList   - Thread-safe ArrayList                   │
│  • CopyOnWriteArraySet    - Thread-safe Set                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 6. Atomic Variables and Volatile

### 6.1 volatile Keyword

The `volatile` keyword ensures visibility of changes across threads and prevents instruction reordering.

```
┌─────────────────────────────────────────────────────────────────────┐
│                        VOLATILE KEYWORD                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Without volatile:                                                  │
│  ┌─────────────────┐         ┌─────────────────┐                    │
│  │    Thread A     │         │    Thread B     │                    │
│  │  ┌───────────┐  │         │  ┌───────────┐  │                    │
│  │  │ CPU Cache │  │         │  │ CPU Cache │  │                    │
│  │  │ flag=true │  │         │  │ flag=false│  │ ← Stale value!    │
│  │  └───────────┘  │         │  └───────────┘  │                    │
│  └────────┬────────┘         └────────┬────────┘                    │
│           │                           │                              │
│           ▼                           ▼                              │
│  ┌─────────────────────────────────────────────┐                    │
│  │           Main Memory: flag=true            │                    │
│  └─────────────────────────────────────────────┘                    │
│                                                                      │
│  With volatile:                                                     │
│  ┌─────────────────┐         ┌─────────────────┐                    │
│  │    Thread A     │         │    Thread B     │                    │
│  │  writes flag    │────────►│  reads flag     │                    │
│  │  (flush to main)│         │  (read from main)│                   │
│  └─────────────────┘         └─────────────────┘                    │
│                                                                      │
│  volatile guarantees:                                               │
│  1. Visibility - changes visible to all threads immediately        │
│  2. Ordering - prevents instruction reordering                      │
│  3. NOT atomicity - compound operations still unsafe               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
public class VolatileDemo {
    // Without volatile - may never stop!
    // private boolean running = true;
    
    // With volatile - visibility guaranteed
    private volatile boolean running = true;
    
    public void startWorker() {
        new Thread(() -> {
            System.out.println("Worker started");
            while (running) {
                // Do work
            }
            System.out.println("Worker stopped");
        }).start();
    }
    
    public void stopWorker() {
        running = false; // Visible to worker thread immediately
    }
    
    public static void main(String[] args) throws InterruptedException {
        VolatileDemo demo = new VolatileDemo();
        demo.startWorker();
        
        Thread.sleep(1000);
        System.out.println("Main: Stopping worker...");
        demo.stopWorker();
    }
}
```

#### When to Use volatile

```java
public class VolatileUseCases {
    
    // USE CASE 1: Status flags
    private volatile boolean initialized = false;
    
    public void initialize() {
        // Do initialization...
        initialized = true; // Signal completion
    }
    
    public void doWork() {
        while (!initialized) {
            Thread.yield(); // Wait for initialization
        }
        // Proceed with work
    }
    
    // USE CASE 2: Double-checked locking (Singleton)
    private static volatile VolatileUseCases instance;
    
    public static VolatileUseCases getInstance() {
        if (instance == null) {
            synchronized (VolatileUseCases.class) {
                if (instance == null) {
                    instance = new VolatileUseCases();
                }
            }
        }
        return instance;
    }
    
    // WRONG USE: Counter (NOT atomic!)
    private volatile int counter = 0;
    
    public void increment() {
        counter++; // NOT thread-safe! (read-modify-write)
    }
}
```

#### volatile vs synchronized

| Feature | volatile | synchronized |
|---------|----------|-------------|
| **Visibility** | Yes | Yes |
| **Atomicity** | No (reads/writes only) | Yes |
| **Blocking** | No | Yes |
| **Use Case** | Status flags, simple state | Compound operations |
| **Performance** | Better | More overhead |

---

### 6.2 AtomicInteger

`AtomicInteger` provides atomic operations on integers without locking.

```java
import java.util.concurrent.atomic.*;
import java.util.concurrent.*;

public class AtomicIntegerDemo {
    private AtomicInteger counter = new AtomicInteger(0);
    
    // Atomic increment
    public void increment() {
        counter.incrementAndGet(); // returns new value
        // or
        counter.getAndIncrement(); // returns old value
    }
    
    // Atomic add
    public void add(int delta) {
        counter.addAndGet(delta);
    }
    
    // Atomic compare and set
    public boolean tryUpdate(int expected, int newValue) {
        return counter.compareAndSet(expected, newValue);
    }
    
    // Atomic update with function
    public void updateWithFunction() {
        counter.updateAndGet(x -> x * 2);      // double the value
        counter.accumulateAndGet(5, Integer::sum); // add 5
    }
    
    public int get() {
        return counter.get();
    }
    
    public static void main(String[] args) throws InterruptedException {
        AtomicIntegerDemo demo = new AtomicIntegerDemo();
        
        ExecutorService executor = Executors.newFixedThreadPool(10);
        
        // 1000 threads incrementing counter
        for (int i = 0; i < 1000; i++) {
            executor.submit(demo::increment);
        }
        
        executor.shutdown();
        executor.awaitTermination(5, TimeUnit.SECONDS);
        
        System.out.println("Final count: " + demo.get()); // Always 1000!
    }
}
```

#### Atomic Classes

| Class | Description |
|-------|-------------|
| `AtomicInteger` | Atomic int operations |
| `AtomicLong` | Atomic long operations |
| `AtomicBoolean` | Atomic boolean operations |
| `AtomicReference<V>` | Atomic reference operations |
| `AtomicIntegerArray` | Atomic array of ints |
| `AtomicLongArray` | Atomic array of longs |
| `AtomicReferenceArray` | Atomic array of references |
| `AtomicStampedReference` | Atomic reference with stamp (ABA) |
| `AtomicMarkableReference` | Atomic reference with mark |
| `LongAdder` | High-contention counter |
| `LongAccumulator` | Custom accumulation |

#### LongAdder for High Contention

```java
import java.util.concurrent.atomic.*;

public class LongAdderDemo {
    // AtomicLong - all threads compete for single variable
    private AtomicLong atomicCounter = new AtomicLong(0);
    
    // LongAdder - distributed cells, better for high contention
    private LongAdder adderCounter = new LongAdder();
    
    public void incrementAtomic() {
        atomicCounter.incrementAndGet();
    }
    
    public void incrementAdder() {
        adderCounter.increment(); // Much faster under high contention!
    }
    
    public long getAtomicValue() {
        return atomicCounter.get();
    }
    
    public long getAdderValue() {
        return adderCounter.sum(); // Sums all cells
    }
}
```

```
┌─────────────────────────────────────────────────────────────────────┐
│              ATOMICLONG vs LONGADDER                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  AtomicLong:                      LongAdder:                        │
│  ─────────────                    ───────────                        │
│  ┌─────────────┐                 ┌─────────────────────────┐        │
│  │   counter   │◄── all threads  │ base │cell1│cell2│cell3│        │
│  │    (CAS)    │   compete       └──┬───┴──┬──┴──┬──┴──┬──┘        │
│  └─────────────┘                    │      │     │     │            │
│                                     ▲      ▲     ▲     ▲            │
│                                     T1     T2    T3    T4           │
│                                                                      │
│  Best for: Low/moderate contention  Best for: High contention       │
│  Exact value always available       Value computed on sum()         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

### 6.3 Compare and Swap (CAS)

CAS is a CPU instruction that atomically compares a memory location's value and replaces it if it matches the expected value.

```
┌─────────────────────────────────────────────────────────────────────┐
│                    COMPARE AND SWAP (CAS)                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  CAS(memory_location, expected_value, new_value):                   │
│                                                                      │
│  if (memory_location == expected_value) {                           │
│      memory_location = new_value;                                   │
│      return true;  // Success                                       │
│  } else {                                                           │
│      return false; // Failure - retry                               │
│  }                                                                  │
│                                                                      │
│  This entire operation is ATOMIC at hardware level!                 │
│                                                                      │
│  Example: Increment counter from 5 to 6                             │
│                                                                      │
│  Thread A: CAS(counter, 5, 6) → true (counter is now 6)            │
│  Thread B: CAS(counter, 5, 6) → false (counter is 6, not 5!)       │
│  Thread B: retry... read 6, CAS(counter, 6, 7) → true              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.atomic.*;

public class CASDemo {
    private AtomicInteger value = new AtomicInteger(0);
    
    // Implementing custom atomic operation using CAS
    public int incrementAndDoubleIfOdd() {
        while (true) {
            int current = value.get();
            int next = (current % 2 == 1) ? (current + 1) * 2 : current + 1;
            
            if (value.compareAndSet(current, next)) {
                return next; // Success!
            }
            // CAS failed, another thread modified value, retry...
        }
    }
    
    // Custom lock-free stack using CAS
    static class LockFreeStack<E> {
        private AtomicReference<Node<E>> top = new AtomicReference<>();
        
        static class Node<E> {
            final E item;
            Node<E> next;
            Node(E item) { this.item = item; }
        }
        
        public void push(E item) {
            Node<E> newNode = new Node<>(item);
            while (true) {
                Node<E> currentTop = top.get();
                newNode.next = currentTop;
                if (top.compareAndSet(currentTop, newNode)) {
                    return; // Success!
                }
                // Retry...
            }
        }
        
        public E pop() {
            while (true) {
                Node<E> currentTop = top.get();
                if (currentTop == null) {
                    return null; // Stack empty
                }
                if (top.compareAndSet(currentTop, currentTop.next)) {
                    return currentTop.item; // Success!
                }
                // Retry...
            }
        }
    }
}
```

#### ABA Problem

```
┌─────────────────────────────────────────────────────────────────────┐
│                        ABA PROBLEM                                   │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Thread 1: Read value A                                             │
│  Thread 2: Change A → B → A                                         │
│  Thread 1: CAS(A, A, C) → Success! (but value was modified!)       │
│                                                                      │
│  Timeline:                                                          │
│  Thread 1         Memory          Thread 2                          │
│      │              A                 │                              │
│  read A ◄──────────┤                  │                              │
│      │              │                 │                              │
│      │              │             change to B                        │
│      │              B ◄───────────────┤                              │
│      │              │             change to A                        │
│      │              A ◄───────────────┤                              │
│      │              │                 │                              │
│  CAS(A,A,C) ───────►C                 │  ← Problem: Thread 1        │
│                     │                 │    doesn't know A changed!  │
│                                                                      │
│  Solution: Use AtomicStampedReference (adds version stamp)         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.atomic.*;

public class ABADemo {
    // AtomicStampedReference solves ABA problem
    private AtomicStampedReference<String> ref = 
        new AtomicStampedReference<>("A", 0);
    
    public void updateWithStamp() {
        int[] stampHolder = new int[1];
        String current = ref.get(stampHolder);
        int currentStamp = stampHolder[0];
        
        // CAS with stamp - detects ABA!
        boolean success = ref.compareAndSet(
            current,           // expected reference
            "C",              // new reference
            currentStamp,      // expected stamp
            currentStamp + 1   // new stamp
        );
        
        System.out.println("Update success: " + success);
    }
}
```

---

## 7. Fork/Join Framework

The Fork/Join framework (Java 7+) is designed for parallel processing of tasks that can be recursively broken down into smaller subtasks.

### 7.1 Work-Stealing Algorithm

```
┌─────────────────────────────────────────────────────────────────────┐
│                    WORK-STEALING ALGORITHM                           │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Each worker thread has its own deque (double-ended queue)          │
│                                                                      │
│  Worker 1 Deque        Worker 2 Deque        Worker 3 Deque         │
│  ┌──────────────┐      ┌──────────────┐      ┌──────────────┐       │
│  │ Task A       │      │ Task D       │      │ (empty)      │       │
│  │ Task B       │      │ Task E       │      │              │       │
│  │ Task C       │      │              │      │              │       │
│  └──────────────┘      └──────────────┘      └──────────────┘       │
│        │                     │                     │                 │
│        ▼ pop (LIFO)          ▼                     │                 │
│   Process own tasks      Process own          Worker 3 STEALS       │
│   from bottom            tasks                from TOP of Worker 1  │
│                                                     │                │
│                                               ┌─────▼─────┐          │
│                                               │ Steals    │          │
│                                               │ Task A    │          │
│                                               └───────────┘          │
│                                                                      │
│  Benefits:                                                          │
│  • Keeps all cores busy                                             │
│  • Minimizes contention (steal from opposite end)                   │
│  • Automatic load balancing                                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 7.2 RecursiveTask and RecursiveAction

```
┌─────────────────────────────────────────────────────────────────────┐
│                   FORK/JOIN TASK TYPES                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ForkJoinTask<V>                                                    │
│       │                                                              │
│       ├── RecursiveTask<V>    - Returns a result                    │
│       │         compute() → V                                       │
│       │                                                              │
│       └── RecursiveAction     - No result (void)                    │
│                 compute() → void                                    │
│                                                                      │
│  Pattern:                                                           │
│  if (task is small enough) {                                        │
│      compute directly                                               │
│  } else {                                                           │
│      split into subtasks                                            │
│      fork subtasks                                                  │
│      join results                                                   │
│  }                                                                  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### RecursiveTask Example: Parallel Sum

```java
import java.util.concurrent.*;

public class ParallelSumTask extends RecursiveTask<Long> {
    private static final int THRESHOLD = 10_000; // Sequential threshold
    
    private final long[] array;
    private final int start;
    private final int end;
    
    public ParallelSumTask(long[] array, int start, int end) {
        this.array = array;
        this.start = start;
        this.end = end;
    }
    
    @Override
    protected Long compute() {
        int length = end - start;
        
        // Base case: Small enough to compute directly
        if (length <= THRESHOLD) {
            return computeDirectly();
        }
        
        // Recursive case: Split into two subtasks
        int mid = start + length / 2;
        
        // Fork left subtask (executes asynchronously)
        ParallelSumTask leftTask = new ParallelSumTask(array, start, mid);
        leftTask.fork();
        
        // Compute right subtask in current thread
        ParallelSumTask rightTask = new ParallelSumTask(array, mid, end);
        Long rightResult = rightTask.compute();
        
        // Join left result (waits if necessary)
        Long leftResult = leftTask.join();
        
        return leftResult + rightResult;
    }
    
    private long computeDirectly() {
        long sum = 0;
        for (int i = start; i < end; i++) {
            sum += array[i];
        }
        return sum;
    }
    
    public static void main(String[] args) {
        // Create array with 100 million elements
        int size = 100_000_000;
        long[] array = new long[size];
        for (int i = 0; i < size; i++) {
            array[i] = i + 1;
        }
        
        // Sequential sum
        long startTime = System.currentTimeMillis();
        long sequentialSum = 0;
        for (long num : array) {
            sequentialSum += num;
        }
        System.out.println("Sequential sum: " + sequentialSum);
        System.out.println("Sequential time: " + (System.currentTimeMillis() - startTime) + "ms");
        
        // Parallel sum using Fork/Join
        ForkJoinPool pool = ForkJoinPool.commonPool();
        startTime = System.currentTimeMillis();
        
        Long parallelSum = pool.invoke(new ParallelSumTask(array, 0, array.length));
        
        System.out.println("Parallel sum: " + parallelSum);
        System.out.println("Parallel time: " + (System.currentTimeMillis() - startTime) + "ms");
        System.out.println("Parallelism: " + pool.getParallelism());
    }
}
```

**Output (on 8-core machine):**
```
Sequential sum: 5000000050000000
Sequential time: 85ms
Parallel sum: 5000000050000000
Parallel time: 25ms
Parallelism: 7
```

#### RecursiveAction Example: Parallel Array Update

```java
import java.util.concurrent.*;

public class ParallelArrayUpdate extends RecursiveAction {
    private static final int THRESHOLD = 10_000;
    
    private final int[] array;
    private final int start;
    private final int end;
    
    public ParallelArrayUpdate(int[] array, int start, int end) {
        this.array = array;
        this.start = start;
        this.end = end;
    }
    
    @Override
    protected void compute() {
        if (end - start <= THRESHOLD) {
            // Base case: update directly
            for (int i = start; i < end; i++) {
                array[i] = array[i] * 2;
            }
        } else {
            // Split and fork
            int mid = start + (end - start) / 2;
            invokeAll(
                new ParallelArrayUpdate(array, start, mid),
                new ParallelArrayUpdate(array, mid, end)
            );
        }
    }
    
    public static void main(String[] args) {
        int[] array = new int[1_000_000];
        for (int i = 0; i < array.length; i++) {
            array[i] = i;
        }
        
        ForkJoinPool pool = ForkJoinPool.commonPool();
        pool.invoke(new ParallelArrayUpdate(array, 0, array.length));
        
        // Verify
        System.out.println("array[100] = " + array[100]); // 200
        System.out.println("array[999999] = " + array[999999]); // 1999998
    }
}
```

#### ForkJoinPool Methods

| Method | Description |
|--------|-------------|
| `invoke(task)` | Executes task and waits for result |
| `execute(task)` | Submits task for async execution |
| `submit(task)` | Submits task, returns Future |
| `commonPool()` | Returns shared common pool |
| `getParallelism()` | Returns parallelism level |
| `getPoolSize()` | Returns number of worker threads |
| `getActiveThreadCount()` | Returns active threads |
| `getStealCount()` | Returns total steals |

---

## 8. CompletableFuture

`CompletableFuture` (Java 8+) is a powerful class for asynchronous programming with support for composition, combining, and exception handling.

### 8.1 Asynchronous Programming

```
┌─────────────────────────────────────────────────────────────────────┐
│                    COMPLETABLEFUTURE                                 │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Future Limitations:                                                │
│  • Cannot be manually completed                                     │
│  • No way to chain multiple futures                                 │
│  • No exception handling                                            │
│  • Cannot combine multiple futures                                  │
│                                                                      │
│  CompletableFuture Advantages:                                      │
│  ✓ Can be manually completed                                        │
│  ✓ Chain dependent operations (thenApply, thenAccept)              │
│  ✓ Combine multiple futures (thenCombine, allOf, anyOf)            │
│  ✓ Exception handling (exceptionally, handle)                      │
│  ✓ Non-blocking callbacks                                          │
│                                                                      │
│  Execution Model:                                                   │
│  ┌────────────┐    ┌────────────┐    ┌────────────┐                │
│  │  Stage 1   │───►│  Stage 2   │───►│  Stage 3   │                │
│  │ (async)    │    │ (depends)  │    │ (depends)  │                │
│  └────────────┘    └────────────┘    └────────────┘                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.concurrent.*;

public class CompletableFutureBasics {
    public static void main(String[] args) throws Exception {
        
        // 1. Create completed future
        CompletableFuture<String> completed = 
            CompletableFuture.completedFuture("Already done!");
        System.out.println(completed.get());
        
        // 2. Run async without return value
        CompletableFuture<Void> runAsync = CompletableFuture.runAsync(() -> {
            System.out.println("Running in: " + Thread.currentThread().getName());
        });
        runAsync.join();
        
        // 3. Supply async with return value
        CompletableFuture<String> supplyAsync = CompletableFuture.supplyAsync(() -> {
            try {
                Thread.sleep(1000);
            } catch (InterruptedException e) {}
            return "Hello from async!";
        });
        System.out.println(supplyAsync.get()); // Blocks until complete
        
        // 4. Manual completion
        CompletableFuture<String> future = new CompletableFuture<>();
        new Thread(() -> {
            try {
                Thread.sleep(500);
                future.complete("Manually completed!");
            } catch (InterruptedException e) {
                future.completeExceptionally(e);
            }
        }).start();
        System.out.println("Result: " + future.get());
        
        // 5. With custom executor
        ExecutorService executor = Executors.newFixedThreadPool(4);
        CompletableFuture<String> withExecutor = 
            CompletableFuture.supplyAsync(() -> "Custom executor!", executor);
        System.out.println(withExecutor.get());
        executor.shutdown();
    }
}
```

---

### 8.2 Chaining Tasks

```java
import java.util.concurrent.*;

public class CompletableFutureChaining {
    public static void main(String[] args) throws Exception {
        
        // thenApply - Transform result (like map)
        CompletableFuture<String> future1 = CompletableFuture
            .supplyAsync(() -> "hello")
            .thenApply(s -> s.toUpperCase())     // HELLO
            .thenApply(s -> s + " WORLD");       // HELLO WORLD
        System.out.println("thenApply: " + future1.get());
        
        // thenAccept - Consume result (no return)
        CompletableFuture<Void> future2 = CompletableFuture
            .supplyAsync(() -> "Message")
            .thenAccept(s -> System.out.println("thenAccept: " + s));
        future2.join();
        
        // thenRun - Run action after completion (no access to result)
        CompletableFuture<Void> future3 = CompletableFuture
            .supplyAsync(() -> "Ignored result")
            .thenRun(() -> System.out.println("thenRun: Task completed!"));
        future3.join();
        
        // thenCompose - Flatten nested CompletableFuture (like flatMap)
        CompletableFuture<String> future4 = CompletableFuture
            .supplyAsync(() -> 1)
            .thenCompose(num -> getUserById(num)); // Returns CompletableFuture
        System.out.println("thenCompose: " + future4.get());
        
        // thenCombine - Combine two independent futures
        CompletableFuture<String> nameFuture = 
            CompletableFuture.supplyAsync(() -> "John");
        CompletableFuture<Integer> ageFuture = 
            CompletableFuture.supplyAsync(() -> 30);
        
        CompletableFuture<String> combined = nameFuture
            .thenCombine(ageFuture, (name, age) -> name + " is " + age);
        System.out.println("thenCombine: " + combined.get());
        
        // Async variants (executes in different thread)
        CompletableFuture<String> asyncVariant = CompletableFuture
            .supplyAsync(() -> "data")
            .thenApplyAsync(s -> s.toUpperCase()); // Async transformation
    }
    
    static CompletableFuture<String> getUserById(int id) {
        return CompletableFuture.supplyAsync(() -> "User-" + id);
    }
}
```

#### Combining Multiple Futures

```java
import java.util.concurrent.*;
import java.util.*;
import java.util.stream.*;

public class CombiningFutures {
    public static void main(String[] args) throws Exception {
        
        // allOf - Wait for ALL futures to complete
        CompletableFuture<String> f1 = CompletableFuture.supplyAsync(() -> {
            sleep(1000); return "Result 1";
        });
        CompletableFuture<String> f2 = CompletableFuture.supplyAsync(() -> {
            sleep(2000); return "Result 2";
        });
        CompletableFuture<String> f3 = CompletableFuture.supplyAsync(() -> {
            sleep(1500); return "Result 3";
        });
        
        long start = System.currentTimeMillis();
        
        CompletableFuture<Void> allOf = CompletableFuture.allOf(f1, f2, f3);
        allOf.join(); // Wait for all
        
        System.out.println("All completed in: " + 
                         (System.currentTimeMillis() - start) + "ms"); // ~2000ms
        
        // Get results (futures are complete)
        List<String> results = Stream.of(f1, f2, f3)
            .map(CompletableFuture::join)
            .collect(Collectors.toList());
        System.out.println("Results: " + results);
        
        // anyOf - Complete when FIRST future completes
        CompletableFuture<String> fast = CompletableFuture.supplyAsync(() -> {
            sleep(500); return "Fast!";
        });
        CompletableFuture<String> slow = CompletableFuture.supplyAsync(() -> {
            sleep(2000); return "Slow...";
        });
        
        start = System.currentTimeMillis();
        CompletableFuture<Object> anyOf = CompletableFuture.anyOf(fast, slow);
        System.out.println("First result: " + anyOf.get()); // Fast!
        System.out.println("Time: " + (System.currentTimeMillis() - start) + "ms"); // ~500ms
    }
    
    static void sleep(long ms) {
        try { Thread.sleep(ms); } catch (InterruptedException e) {}
    }
}
```

```
┌─────────────────────────────────────────────────────────────────────┐
│              COMPLETABLEFUTURE METHOD CATEGORIES                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Creation:                                                          │
│  • completedFuture(value)     - Already completed                   │
│  • supplyAsync(supplier)      - Async with return value             │
│  • runAsync(runnable)         - Async without return value          │
│                                                                      │
│  Transformation (returns new CompletableFuture):                    │
│  • thenApply(fn)              - Transform result                    │
│  • thenCompose(fn)            - Flatten nested future              │
│                                                                      │
│  Consumption (returns Void):                                        │
│  • thenAccept(consumer)       - Consume result                      │
│  • thenRun(runnable)          - Run action after completion        │
│                                                                      │
│  Combination:                                                       │
│  • thenCombine(other, fn)     - Combine two futures                │
│  • thenAcceptBoth(other, fn)  - Consume both results               │
│  • allOf(futures...)          - Wait for all                        │
│  • anyOf(futures...)          - Wait for any                        │
│                                                                      │
│  Exception Handling:                                                │
│  • exceptionally(fn)          - Handle exception                    │
│  • handle(fn)                 - Handle result or exception          │
│  • whenComplete(action)       - Action on completion                │
│                                                                      │
│  Async Variants (add 'Async' suffix):                              │
│  • thenApplyAsync             - Run in different thread             │
│  • thenAcceptAsync            - etc.                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

### 8.3 Exception Handling

```java
import java.util.concurrent.*;

public class CompletableFutureExceptionHandling {
    public static void main(String[] args) {
        
        // 1. exceptionally - Handle exception, return fallback
        CompletableFuture<String> future1 = CompletableFuture
            .supplyAsync(() -> {
                if (true) throw new RuntimeException("Error!");
                return "Success";
            })
            .exceptionally(ex -> {
                System.out.println("Exception: " + ex.getMessage());
                return "Fallback value";
            });
        System.out.println("exceptionally result: " + future1.join());
        
        // 2. handle - Handle both success and failure
        CompletableFuture<String> future2 = CompletableFuture
            .supplyAsync(() -> {
                if (Math.random() > 0.5) {
                    throw new RuntimeException("Random failure!");
                }
                return "Success!";
            })
            .handle((result, ex) -> {
                if (ex != null) {
                    return "Handled error: " + ex.getMessage();
                }
                return "Result: " + result;
            });
        System.out.println("handle result: " + future2.join());
        
        // 3. whenComplete - Side effect on completion (doesn't transform result)
        CompletableFuture<String> future3 = CompletableFuture
            .supplyAsync(() -> "Data")
            .whenComplete((result, ex) -> {
                if (ex != null) {
                    System.out.println("Failed: " + ex.getMessage());
                } else {
                    System.out.println("Completed with: " + result);
                }
            });
        future3.join();
        
        // 4. Chain with exception handling in middle
        CompletableFuture<Integer> future4 = CompletableFuture
            .supplyAsync(() -> "10")
            .thenApply(s -> Integer.parseInt(s))   // Could throw
            .exceptionally(ex -> 0)                 // Recover
            .thenApply(n -> n * 2);                 // Continue
        System.out.println("Chained: " + future4.join());
        
        // 5. completeExceptionally - Manual exceptional completion
        CompletableFuture<String> manual = new CompletableFuture<>();
        manual.completeExceptionally(new RuntimeException("Manual error"));
        try {
            manual.join();
        } catch (CompletionException e) {
            System.out.println("Caught: " + e.getCause().getMessage());
        }
    }
}
```

#### Practical Example: Async API Calls

```java
import java.util.concurrent.*;
import java.util.*;

public class AsyncApiExample {
    private static ExecutorService executor = Executors.newFixedThreadPool(10);
    
    public static void main(String[] args) {
        long start = System.currentTimeMillis();
        
        // Simulate fetching user profile with multiple API calls
        CompletableFuture<UserProfile> profileFuture = fetchUserAsync(1)
            .thenCombine(fetchOrdersAsync(1), (user, orders) -> {
                user.setOrders(orders);
                return user;
            })
            .thenCombine(fetchRecommendationsAsync(1), (user, recs) -> {
                user.setRecommendations(recs);
                return user;
            })
            .exceptionally(ex -> {
                System.err.println("Error: " + ex.getMessage());
                return new UserProfile("Error User");
            });
        
        UserProfile profile = profileFuture.join();
        
        System.out.println("Profile: " + profile);
        System.out.println("Total time: " + (System.currentTimeMillis() - start) + "ms");
        // ~1000ms (parallel) instead of ~3000ms (sequential)
        
        executor.shutdown();
    }
    
    // Simulated async API calls
    static CompletableFuture<UserProfile> fetchUserAsync(int userId) {
        return CompletableFuture.supplyAsync(() -> {
            sleep(1000);
            return new UserProfile("User-" + userId);
        }, executor);
    }
    
    static CompletableFuture<List<String>> fetchOrdersAsync(int userId) {
        return CompletableFuture.supplyAsync(() -> {
            sleep(800);
            return Arrays.asList("Order-1", "Order-2");
        }, executor);
    }
    
    static CompletableFuture<List<String>> fetchRecommendationsAsync(int userId) {
        return CompletableFuture.supplyAsync(() -> {
            sleep(600);
            return Arrays.asList("Product-A", "Product-B");
        }, executor);
    }
    
    static void sleep(long ms) {
        try { Thread.sleep(ms); } catch (InterruptedException e) {}
    }
    
    static class UserProfile {
        String name;
        List<String> orders = new ArrayList<>();
        List<String> recommendations = new ArrayList<>();
        
        UserProfile(String name) { this.name = name; }
        void setOrders(List<String> orders) { this.orders = orders; }
        void setRecommendations(List<String> recs) { this.recommendations = recs; }
        
        @Override
        public String toString() {
            return "UserProfile{name='" + name + "', orders=" + orders + 
                   ", recommendations=" + recommendations + "}";
        }
    }
}
```

#### CompletableFuture Best Practices

```
┌─────────────────────────────────────────────────────────────────────┐
│              COMPLETABLEFUTURE BEST PRACTICES                        │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. Always handle exceptions                                        │
│     • Use exceptionally(), handle(), or whenComplete()              │
│     • Unhandled exceptions will be swallowed!                       │
│                                                                      │
│  2. Use custom executor for blocking operations                     │
│     • Common pool has limited threads                               │
│     • Blocking can starve other tasks                               │
│                                                                      │
│  3. Avoid blocking calls like get()                                 │
│     • Use join() if you must block                                  │
│     • Prefer chaining with thenApply, thenAccept                    │
│                                                                      │
│  4. Use thenCompose() instead of thenApply() for nested futures    │
│     • thenApply returns CompletableFuture<CompletableFuture<T>>    │
│     • thenCompose flattens to CompletableFuture<T>                 │
│                                                                      │
│  5. Use appropriate async variant                                   │
│     • thenApply - runs in same thread as previous stage            │
│     • thenApplyAsync - runs in executor                             │
│                                                                      │
│  6. Prefer allOf/anyOf for multiple independent futures            │
│     • Enables parallel execution                                    │
│     • Better than sequential waits                                  │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 9. Best Practices

Writing correct and efficient multithreaded code requires careful consideration of various factors. This section covers essential best practices for concurrent programming.

### 9.1 Avoiding Deadlocks

Deadlocks can bring your application to a complete halt. Here are proven strategies to prevent them:

```
┌─────────────────────────────────────────────────────────────────────┐
│                  DEADLOCK PREVENTION STRATEGIES                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. LOCK ORDERING                                                   │
│     Always acquire locks in a consistent, global order              │
│                                                                      │
│  2. LOCK TIMEOUT                                                    │
│     Use tryLock() with timeout instead of lock()                    │
│                                                                      │
│  3. AVOID NESTED LOCKS                                              │
│     Minimize holding multiple locks simultaneously                  │
│                                                                      │
│  4. AVOID CALLING UNKNOWN CODE WHILE HOLDING LOCK                   │
│     Don't call callbacks/listeners in synchronized blocks           │
│                                                                      │
│  5. USE HIGHER-LEVEL CONCURRENCY UTILITIES                          │
│     Prefer java.util.concurrent over manual synchronization         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Strategy 1: Lock Ordering

```java
public class LockOrderingExample {
    // Always acquire locks in the same order based on object identity
    private static final Object tieLock = new Object();
    
    public void transferMoney(Account from, Account to, int amount) {
        // Establish consistent lock ordering
        int fromHash = System.identityHashCode(from);
        int toHash = System.identityHashCode(to);
        
        Object firstLock, secondLock;
        if (fromHash < toHash) {
            firstLock = from;
            secondLock = to;
        } else if (fromHash > toHash) {
            firstLock = to;
            secondLock = from;
        } else {
            // Rare hash collision - use tie-breaker lock
            synchronized (tieLock) {
                synchronized (from) {
                    synchronized (to) {
                        doTransfer(from, to, amount);
                    }
                }
            }
            return;
        }
        
        synchronized (firstLock) {
            synchronized (secondLock) {
                doTransfer(from, to, amount);
            }
        }
    }
    
    private void doTransfer(Account from, Account to, int amount) {
        if (from.getBalance() >= amount) {
            from.debit(amount);
            to.credit(amount);
        }
    }
}
```

#### Strategy 2: Try Lock with Timeout

```java
import java.util.concurrent.locks.*;
import java.util.concurrent.*;

public class TryLockExample {
    private final Lock lock1 = new ReentrantLock();
    private final Lock lock2 = new ReentrantLock();
    
    public boolean transferWithTimeout(int amount) {
        long deadline = System.nanoTime() + TimeUnit.SECONDS.toNanos(10);
        
        while (true) {
            if (lock1.tryLock()) {
                try {
                    if (lock2.tryLock()) {
                        try {
                            // Both locks acquired - do transfer
                            performTransfer(amount);
                            return true;
                        } finally {
                            lock2.unlock();
                        }
                    }
                } finally {
                    lock1.unlock();
                }
            }
            
            // Check if we've exceeded deadline
            if (System.nanoTime() >= deadline) {
                return false; // Could not acquire locks in time
            }
            
            // Back off and retry
            try {
                Thread.sleep(10 + (int)(Math.random() * 50));
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                return false;
            }
        }
    }
    
    private void performTransfer(int amount) {
        // Transfer logic
    }
}
```

#### Strategy 3: Open Calls (Avoid Alien Method Calls)

```java
public class OpenCallExample {
    private final List<Listener> listeners = new CopyOnWriteArrayList<>();
    private int value;
    
    // BAD: Calling unknown code while holding lock
    public synchronized void setValueBad(int newValue) {
        this.value = newValue;
        // DANGEROUS: listener.onValueChanged() is alien code!
        // It could try to acquire other locks -> potential deadlock
        for (Listener listener : listeners) {
            listener.onValueChanged(newValue);
        }
    }
    
    // GOOD: Open call - release lock before calling listeners
    public void setValueGood(int newValue) {
        synchronized (this) {
            this.value = newValue;
        }
        // Notify listeners outside synchronized block
        for (Listener listener : listeners) {
            listener.onValueChanged(newValue);
        }
    }
    
    interface Listener {
        void onValueChanged(int value);
    }
}
```

---

### 9.2 Minimizing Synchronization

Synchronization is expensive. Minimize its scope and duration for better performance.

```
┌─────────────────────────────────────────────────────────────────────┐
│               MINIMIZING SYNCHRONIZATION                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  PRINCIPLE: Hold locks for the shortest time possible               │
│                                                                      │
│  1. Narrow the synchronized block                                   │
│  2. Don't do I/O, expensive computations while holding lock         │
│  3. Use fine-grained locking (multiple locks for independent data)  │
│  4. Use lock-free data structures when possible                     │
│  5. Prefer read-write locks for read-heavy workloads                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Narrowing Synchronized Blocks

```java
public class NarrowSynchronizationExample {
    private final Object lock = new Object();
    private Map<String, Data> cache = new HashMap<>();
    
    // BAD: Entire method synchronized
    public synchronized Data getDataBad(String key) {
        Data data = cache.get(key);
        if (data == null) {
            data = fetchFromDatabase(key);  // Expensive I/O while holding lock!
            cache.put(key, data);
        }
        return data;
    }
    
    // GOOD: Minimal synchronized scope
    public Data getDataGood(String key) {
        // First, try to get from cache with minimal locking
        synchronized (lock) {
            Data data = cache.get(key);
            if (data != null) {
                return data;
            }
        }
        
        // Fetch outside synchronized block
        Data data = fetchFromDatabase(key);
        
        // Put in cache with minimal locking
        synchronized (lock) {
            // Double-check: another thread might have populated
            Data existing = cache.get(key);
            if (existing != null) {
                return existing;
            }
            cache.put(key, data);
        }
        
        return data;
    }
    
    private Data fetchFromDatabase(String key) {
        // Expensive database call
        return new Data();
    }
    
    static class Data {}
}
```

#### Fine-Grained Locking (Lock Striping)

```java
public class StripedLockExample<K, V> {
    private static final int NUM_LOCKS = 16;
    private final Object[] locks;
    private final Map<K, V>[] buckets;
    
    @SuppressWarnings("unchecked")
    public StripedLockExample() {
        locks = new Object[NUM_LOCKS];
        buckets = new Map[NUM_LOCKS];
        
        for (int i = 0; i < NUM_LOCKS; i++) {
            locks[i] = new Object();
            buckets[i] = new HashMap<>();
        }
    }
    
    private int getLockIndex(K key) {
        return Math.abs(key.hashCode() % NUM_LOCKS);
    }
    
    public void put(K key, V value) {
        int index = getLockIndex(key);
        synchronized (locks[index]) {  // Only lock one stripe
            buckets[index].put(key, value);
        }
    }
    
    public V get(K key) {
        int index = getLockIndex(key);
        synchronized (locks[index]) {  // Only lock one stripe
            return buckets[index].get(key);
        }
    }
    
    // Different keys can be accessed concurrently if they hash to different stripes!
}
```

#### Using Read-Write Locks for Read-Heavy Workloads

```java
import java.util.concurrent.locks.*;

public class ReadWriteOptimization {
    private final ReadWriteLock rwLock = new ReentrantReadWriteLock();
    private final Lock readLock = rwLock.readLock();
    private final Lock writeLock = rwLock.writeLock();
    
    private Map<String, Object> data = new HashMap<>();
    
    // Multiple readers can access simultaneously
    public Object read(String key) {
        readLock.lock();
        try {
            return data.get(key);
        } finally {
            readLock.unlock();
        }
    }
    
    // Writers have exclusive access
    public void write(String key, Object value) {
        writeLock.lock();
        try {
            data.put(key, value);
        } finally {
            writeLock.unlock();
        }
    }
    
    // Great for configurations, caches with infrequent updates
}
```

---

### 9.3 Thread Safety Strategies

There are multiple approaches to achieving thread safety. Choose the right one based on your use case.

```
┌─────────────────────────────────────────────────────────────────────┐
│                  THREAD SAFETY STRATEGIES                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Strategy               Best For                                    │
│  ────────────────────────────────────────────────────────────────    │
│                                                                      │
│  1. Immutability        Data that doesn't change                    │
│     (Preferred)         Simple, no synchronization needed           │
│                                                                      │
│  2. Thread Confinement  Data accessed by single thread              │
│     (ThreadLocal)       Servlet requests, UI threads                │
│                                                                      │
│  3. Lock-Based Sync     Mutable shared state                        │
│     (synchronized/Lock) Full control needed                         │
│                                                                      │
│  4. Lock-Free (Atomic)  Simple counters, flags                      │
│     (CAS operations)    High contention scenarios                   │
│                                                                      │
│  5. Concurrent Colls    Shared collections                          │
│     (ConcurrentHashMap) High-concurrency access                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Strategy 1: Thread Confinement with ThreadLocal

```java
public class ThreadLocalExample {
    // Each thread has its own copy - no synchronization needed!
    private static final ThreadLocal<SimpleDateFormat> dateFormat = 
        ThreadLocal.withInitial(() -> new SimpleDateFormat("yyyy-MM-dd"));
    
    private static final ThreadLocal<DatabaseConnection> connection = 
        ThreadLocal.withInitial(() -> ConnectionPool.getConnection());
    
    public String formatDate(Date date) {
        // Each thread uses its own SimpleDateFormat instance
        return dateFormat.get().format(date);
    }
    
    public void executeQuery(String sql) {
        // Each thread uses its own connection
        connection.get().execute(sql);
    }
    
    // IMPORTANT: Clean up to prevent memory leaks!
    public void cleanup() {
        dateFormat.remove();
        connection.remove();
    }
}
```

#### Strategy 2: Copy-On-Write for Read-Heavy Scenarios

```java
import java.util.concurrent.*;

public class CopyOnWriteExample {
    // Great for listeners/observers - reads vastly outnumber writes
    private final CopyOnWriteArrayList<EventListener> listeners = 
        new CopyOnWriteArrayList<>();
    
    public void addListener(EventListener listener) {
        listeners.add(listener);  // Creates new copy - expensive!
    }
    
    public void removeListener(EventListener listener) {
        listeners.remove(listener);  // Creates new copy - expensive!
    }
    
    public void fireEvent(Event event) {
        // Iteration is fast - no locking needed!
        // Sees snapshot at start of iteration
        for (EventListener listener : listeners) {
            listener.onEvent(event);
        }
    }
    
    interface EventListener {
        void onEvent(Event event);
    }
    
    static class Event {}
}
```

#### Strategy 3: Thread-Safe Singleton Patterns

```java
public class SingletonPatterns {
    
    // 1. Enum Singleton (Recommended - simplest and safest)
    public enum ConfigManager {
        INSTANCE;
        
        private final Map<String, String> config = new HashMap<>();
        
        public String get(String key) {
            return config.get(key);
        }
    }
    
    // 2. Static Holder Pattern (Lazy, thread-safe)
    public static class ServiceLocator {
        private ServiceLocator() {}
        
        private static class Holder {
            static final ServiceLocator INSTANCE = new ServiceLocator();
        }
        
        public static ServiceLocator getInstance() {
            return Holder.INSTANCE;  // Class loaded lazily
        }
    }
    
    // 3. Double-Checked Locking (if you need lazy init with volatile)
    public static class ResourceManager {
        private static volatile ResourceManager instance;
        
        private ResourceManager() {}
        
        public static ResourceManager getInstance() {
            if (instance == null) {
                synchronized (ResourceManager.class) {
                    if (instance == null) {
                        instance = new ResourceManager();
                    }
                }
            }
            return instance;
        }
    }
}
```

---

### 9.4 Immutable Objects

Immutable objects are inherently thread-safe - they require no synchronization.

```
┌─────────────────────────────────────────────────────────────────────┐
│                     IMMUTABILITY RULES                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. Declare class as final (prevent subclassing)                    │
│  2. Make all fields private and final                               │
│  3. Don't provide setters                                           │
│  4. Initialize all fields in constructor                            │
│  5. Make defensive copies of mutable objects                        │
│  6. Don't expose mutable objects                                    │
│                                                                      │
│  Benefits:                                                          │
│  ✓ No synchronization needed                                        │
│  ✓ Can be freely shared between threads                             │
│  ✓ Safe to use as Map keys and Set elements                         │
│  ✓ Easier to reason about                                           │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Complete Immutable Class Example

```java
import java.util.*;

// Rule 1: Final class
public final class ImmutablePerson {
    // Rule 2: Private final fields
    private final String name;
    private final int age;
    private final List<String> hobbies;
    private final Date birthDate;
    
    // Rule 4: Initialize in constructor with defensive copies
    public ImmutablePerson(String name, int age, List<String> hobbies, Date birthDate) {
        this.name = name;
        this.age = age;
        // Rule 5: Defensive copy of mutable input
        this.hobbies = new ArrayList<>(hobbies);
        this.birthDate = new Date(birthDate.getTime());
    }
    
    // Rule 3: No setters - only getters
    public String getName() {
        return name;
    }
    
    public int getAge() {
        return age;
    }
    
    // Rule 6: Return defensive copy of mutable object
    public List<String> getHobbies() {
        return new ArrayList<>(hobbies);  // Return copy!
    }
    
    public Date getBirthDate() {
        return new Date(birthDate.getTime());  // Return copy!
    }
    
    // Alternative: Return unmodifiable view
    public List<String> getHobbiesUnmodifiable() {
        return Collections.unmodifiableList(hobbies);
    }
    
    // Create modified copy (functional style)
    public ImmutablePerson withAge(int newAge) {
        return new ImmutablePerson(this.name, newAge, this.hobbies, this.birthDate);
    }
    
    @Override
    public String toString() {
        return "ImmutablePerson{name='" + name + "', age=" + age + 
               ", hobbies=" + hobbies + ", birthDate=" + birthDate + "}";
    }
}
```

#### Java Records (Java 14+)

```java
// Records are immutable by default - much simpler!
public record Person(String name, int age, List<String> hobbies) {
    // Compact constructor for defensive copying
    public Person {
        hobbies = List.copyOf(hobbies);  // Unmodifiable copy
    }
    
    // Accessor automatically generated: name(), age(), hobbies()
}

// Usage:
public class RecordExample {
    public static void main(String[] args) {
        Person person = new Person("Alice", 30, List.of("Reading", "Hiking"));
        
        System.out.println(person.name());    // Alice
        System.out.println(person.age());     // 30
        System.out.println(person.hobbies()); // [Reading, Hiking]
        
        // Thread-safe! Can be shared freely
    }
}
```

---

### 9.5 When NOT to Use Multithreading

Multithreading adds complexity and overhead. It's not always the right solution.

```
┌─────────────────────────────────────────────────────────────────────┐
│            WHEN NOT TO USE MULTITHREADING                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ✘ Simple Sequential Tasks                                          │
│    If work is inherently sequential, threading adds overhead         │
│                                                                      │
│  ✘ Tasks Faster Than Context Switch                                 │
│    Thread overhead > computation time = net loss                     │
│                                                                      │
│  ✘ Single CPU Core with CPU-Bound Work                              │
│    No parallelism benefit, only overhead                            │
│                                                                      │
│  ✘ Heavy Shared State                                               │
│    Extensive synchronization negates parallelism benefits            │
│                                                                      │
│  ✘ Simple Scripts or One-Off Tools                                  │
│    Complexity not worth it for simple programs                       │
│                                                                      │
│  ✘ When Async I/O Suffices                                          │
│    Non-blocking I/O may be simpler and more efficient               │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Cost-Benefit Analysis

```java
public class MultithreadingCostAnalysis {
    
    public static void main(String[] args) {
        int[] data = new int[1000];
        for (int i = 0; i < data.length; i++) data[i] = i;
        
        // Scenario 1: Trivial computation - DON'T use threads
        // Thread overhead > computation time
        long start = System.nanoTime();
        int sum = 0;
        for (int n : data) sum += n;  // Nanoseconds
        long sequential = System.nanoTime() - start;
        System.out.println("Sequential: " + sequential + "ns");
        
        // Parallel version would be SLOWER due to:
        // - Thread creation/destruction
        // - Context switching
        // - Synchronization for combining results
        
        // Scenario 2: Heavy computation - DO use threads
        // Thread overhead << computation time
        int[] bigData = new int[10_000_000];
        // ... parallel processing makes sense here
    }
}
```

#### Decision Flowchart

```
┌─────────────────────────────────────────────────────────────────────┐
│           SHOULD I USE MULTITHREADING?                               │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│                  Is task I/O-bound?                                 │
│                         │                                            │
│                ┌────────┴────────┐                                   │
│               YES                NO                                 │
│                │                  │                                  │
│                ▼                  ▼                                  │
│    Consider async I/O      Is task CPU-bound?                       │
│    or threads for                │                                  │
│    parallel I/O         ┌────────┴────────┐                          │
│                        YES                NO                        │
│                         │                  │                         │
│                         ▼                  ▼                         │
│              Multiple cores?      Probably don't                    │
│                    │              need threading                    │
│            ┌───────┴──────┐                                          │
│           YES           NO                                          │
│            │             │                                          │
│            ▼             ▼                                          │
│    ✓ Use parallel     Single thread                                 │
│      processing       is fine                                       │
│                                                                      │
│  Additional checks:                                                 │
│  • Is the task large enough to offset thread overhead?             │
│  • Can the task be parallelized (limited shared state)?            │
│  • Is the added complexity worth the performance gain?             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Alternatives to Multithreading

| Instead of... | Consider... |
|---------------|-------------|
| Multiple threads for I/O | Async I/O (NIO, CompletableFuture) |
| Thread per request | Virtual Threads (Java 21+) |
| Manual thread management | ExecutorService with appropriate pool |
| Threads for batch processing | Parallel Streams (for simple cases) |
| Threads for scheduling | ScheduledExecutorService |
| Complex thread coordination | Actor model (Akka), reactive streams |

#### Virtual Threads (Java 21+) - Game Changer

```java
// Java 21+ Virtual Threads - lightweight threads managed by JVM
public class VirtualThreadExample {
    public static void main(String[] args) throws Exception {
        // Create 100,000 virtual threads easily!
        try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
            for (int i = 0; i < 100_000; i++) {
                executor.submit(() -> {
                    // I/O-bound work - virtual threads excel here
                    Thread.sleep(1000);
                    return "Done";
                });
            }
        }
        
        // Virtual threads are cheap:
        // - 1000s of bytes vs 1MB+ for platform threads
        // - Minimal context switch overhead
        // - Great for I/O-bound work with blocking calls
    }
}
```

---

### Summary: Thread Safety Checklist

```
┌─────────────────────────────────────────────────────────────────────┐
│                  THREAD SAFETY CHECKLIST                             │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  Before writing concurrent code, ask:                               │
│                                                                      │
│  ☐ Can I avoid shared mutable state?                                │
│    → Use immutable objects, local variables, ThreadLocal           │
│                                                                      │
│  ☐ Can I use existing thread-safe classes?                          │
│    → ConcurrentHashMap, AtomicInteger, BlockingQueue               │
│                                                                      │
│  ☐ Have I documented the synchronization policy?                    │
│    → Which lock protects which state?                              │
│                                                                      │
│  ☐ Am I holding locks for minimum time?                             │
│    → No I/O, no callbacks while holding locks                      │
│                                                                      │
│  ☐ Have I checked for deadlock potential?                           │
│    → Lock ordering, tryLock, avoid nested locks                    │
│                                                                      │
│  ☐ Am I using the right concurrency tool?                           │
│    → ExecutorService vs raw threads                                │
│    → CompletableFuture vs callbacks                                │
│    → Atomic vs synchronized                                        │
│                                                                      │
│  ☐ Have I tested with concurrency tools?                            │
│    → Thread sanitizers, stress tests, jcstress                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 10. Practical Projects

This section contains complete, working examples that demonstrate real-world applications of Java multithreading concepts.

### 10.1 Multithreaded File Downloader

A file downloader that downloads multiple files concurrently, with progress tracking and error handling.

```
┌─────────────────────────────────────────────────────────────────────┐
│              MULTITHREADED FILE DOWNLOADER ARCHITECTURE              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│    Main Thread                                                       │
│        │                                                             │
│        ▼                                                             │
│  ┌──────────────┐                                                   │
│  │DownloadManager│                                                   │
│  └──────┬───────┘                                                   │
│         │                                                            │
│         │  submit downloads                                          │
│         ▼                                                            │
│  ┌──────────────────────────────────────────────┐                   │
│  │         ExecutorService (Thread Pool)         │                   │
│  │  ┌────────┐  ┌────────┐  ┌────────┐         │                   │
│  │  │Worker 1│  │Worker 2│  │Worker 3│  ...    │                   │
│  │  └───┬────┘  └───┬────┘  └───┬────┘         │                   │
│  │      │           │           │               │                   │
│  └──────┼───────────┼───────────┼───────────────┘                   │
│         │           │           │                                    │
│         ▼           ▼           ▼                                    │
│    File1.zip    File2.pdf   File3.mp4                               │
│                                                                      │
│  Progress tracked via CompletableFuture + callbacks                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.io.*;
import java.net.*;
import java.nio.file.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;

public class MultithreadedFileDownloader {
    
    private final ExecutorService executor;
    private final int maxConcurrentDownloads;
    private final Map<String, DownloadTask> activeTasks = new ConcurrentHashMap<>();
    
    public MultithreadedFileDownloader(int maxConcurrentDownloads) {
        this.maxConcurrentDownloads = maxConcurrentDownloads;
        this.executor = Executors.newFixedThreadPool(maxConcurrentDownloads);
    }
    
    // Download task with progress tracking
    public static class DownloadTask {
        private final String url;
        private final Path destination;
        private final AtomicLong bytesDownloaded = new AtomicLong(0);
        private volatile long totalBytes = -1;
        private volatile DownloadStatus status = DownloadStatus.PENDING;
        
        public DownloadTask(String url, Path destination) {
            this.url = url;
            this.destination = destination;
        }
        
        public double getProgress() {
            if (totalBytes <= 0) return 0;
            return (double) bytesDownloaded.get() / totalBytes * 100;
        }
        
        public DownloadStatus getStatus() { return status; }
        public String getUrl() { return url; }
        public Path getDestination() { return destination; }
    }
    
    public enum DownloadStatus {
        PENDING, DOWNLOADING, COMPLETED, FAILED, CANCELLED
    }
    
    // Submit a download and get a CompletableFuture for tracking
    public CompletableFuture<DownloadTask> download(String url, Path destination) {
        DownloadTask task = new DownloadTask(url, destination);
        activeTasks.put(url, task);
        
        return CompletableFuture.supplyAsync(() -> {
            try {
                task.status = DownloadStatus.DOWNLOADING;
                performDownload(task);
                task.status = DownloadStatus.COMPLETED;
            } catch (Exception e) {
                task.status = DownloadStatus.FAILED;
                throw new CompletionException(e);
            }
            return task;
        }, executor);
    }
    
    private void performDownload(DownloadTask task) throws IOException {
        URL url = new URL(task.url);
        HttpURLConnection connection = (HttpURLConnection) url.openConnection();
        connection.setRequestMethod("GET");
        connection.setConnectTimeout(10000);
        connection.setReadTimeout(30000);
        
        task.totalBytes = connection.getContentLengthLong();
        
        // Create parent directories if needed
        Files.createDirectories(task.destination.getParent());
        
        try (InputStream in = new BufferedInputStream(connection.getInputStream());
             OutputStream out = new BufferedOutputStream(
                 Files.newOutputStream(task.destination))) {
            
            byte[] buffer = new byte[8192];
            int bytesRead;
            
            while ((bytesRead = in.read(buffer)) != -1) {
                out.write(buffer, 0, bytesRead);
                task.bytesDownloaded.addAndGet(bytesRead);
                
                // Check for cancellation
                if (Thread.currentThread().isInterrupted()) {
                    task.status = DownloadStatus.CANCELLED;
                    throw new InterruptedIOException("Download cancelled");
                }
            }
        } finally {
            connection.disconnect();
        }
    }
    
    // Download multiple files with progress reporting
    public void downloadAll(List<String> urls, Path baseDir, 
                           ProgressCallback callback) {
        List<CompletableFuture<DownloadTask>> futures = new ArrayList<>();
        
        for (String url : urls) {
            String fileName = url.substring(url.lastIndexOf('/') + 1);
            Path destination = baseDir.resolve(fileName);
            
            CompletableFuture<DownloadTask> future = download(url, destination)
                .whenComplete((task, error) -> {
                    if (error != null) {
                        callback.onError(url, error);
                    } else {
                        callback.onComplete(task);
                    }
                });
            
            futures.add(future);
        }
        
        // Progress monitoring in separate thread
        CompletableFuture.runAsync(() -> {
            while (!allDone(futures)) {
                for (DownloadTask task : activeTasks.values()) {
                    if (task.getStatus() == DownloadStatus.DOWNLOADING) {
                        callback.onProgress(task.getUrl(), task.getProgress());
                    }
                }
                try {
                    Thread.sleep(500);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    break;
                }
            }
        });
        
        // Wait for all downloads to complete
        CompletableFuture.allOf(futures.toArray(new CompletableFuture[0])).join();
    }
    
    private boolean allDone(List<CompletableFuture<DownloadTask>> futures) {
        return futures.stream().allMatch(CompletableFuture::isDone);
    }
    
    public void shutdown() {
        executor.shutdown();
        try {
            if (!executor.awaitTermination(60, TimeUnit.SECONDS)) {
                executor.shutdownNow();
            }
        } catch (InterruptedException e) {
            executor.shutdownNow();
            Thread.currentThread().interrupt();
        }
    }
    
    // Callback interface for progress updates
    public interface ProgressCallback {
        void onProgress(String url, double percentComplete);
        void onComplete(DownloadTask task);
        void onError(String url, Throwable error);
    }
    
    // Example usage
    public static void main(String[] args) {
        MultithreadedFileDownloader downloader = new MultithreadedFileDownloader(4);
        
        List<String> urls = Arrays.asList(
            "https://example.com/file1.zip",
            "https://example.com/file2.pdf",
            "https://example.com/file3.mp4"
        );
        
        Path downloadDir = Paths.get("downloads");
        
        downloader.downloadAll(urls, downloadDir, new ProgressCallback() {
            @Override
            public void onProgress(String url, double percent) {
                System.out.printf("%s: %.1f%%\n", url, percent);
            }
            
            @Override
            public void onComplete(DownloadTask task) {
                System.out.println("Completed: " + task.getDestination());
            }
            
            @Override
            public void onError(String url, Throwable error) {
                System.err.println("Failed: " + url + " - " + error.getMessage());
            }
        });
        
        downloader.shutdown();
    }
}
```

---

### 10.2 Parallel Data Processing

A data processing pipeline that processes large datasets in parallel using Fork/Join and parallel streams.

```
┌─────────────────────────────────────────────────────────────────────┐
│              PARALLEL DATA PROCESSING PIPELINE                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌───────────┐    ┌───────────┐    ┌───────────┐    ┌───────────┐  │
│  │  Source   │───▶│  Parse    │───▶│ Transform │───▶│  Collect  │  │
│  │  (File)   │    │  (Split)  │    │ (Parallel)│    │ (Combine) │  │
│  └───────────┘    └───────────┘    └───────────┘    └───────────┘  │
│                                                                      │
│                   Fork/Join Framework                                │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    Main Task                                 │    │
│  │              ┌──────────┴──────────┐                        │    │
│  │              │                     │                        │    │
│  │         Subtask 1            Subtask 2                      │    │
│  │        ┌────┴────┐          ┌────┴────┐                    │    │
│  │        │         │          │         │                    │    │
│  │    Sub 1.1   Sub 1.2    Sub 2.1   Sub 2.2                 │    │
│  │        │         │          │         │                    │    │
│  │        ▼         ▼          ▼         ▼                    │    │
│  │     Process   Process    Process   Process                 │    │
│  │        │         │          │         │                    │    │
│  │        └────┬────┘          └────┬────┘                    │    │
│  │             │                    │                          │    │
│  │        Combine              Combine                         │    │
│  │             └────────┬───────────┘                          │    │
│  │                      │                                      │    │
│  │                Final Result                                 │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.io.*;
import java.nio.file.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.function.*;
import java.util.stream.*;

public class ParallelDataProcessor {
    
    // Generic data processing result
    public static class ProcessingResult<T> {
        private final T result;
        private final long recordsProcessed;
        private final long processingTimeMs;
        private final List<String> errors;
        
        public ProcessingResult(T result, long records, long timeMs, List<String> errors) {
            this.result = result;
            this.recordsProcessed = records;
            this.processingTimeMs = timeMs;
            this.errors = errors;
        }
        
        // Getters
        public T getResult() { return result; }
        public long getRecordsProcessed() { return recordsProcessed; }
        public long getProcessingTimeMs() { return processingTimeMs; }
        public List<String> getErrors() { return errors; }
    }
    
    // Sales Record for example
    public record SalesRecord(
        String productId,
        String region,
        double amount,
        int quantity,
        String date
    ) {}
    
    // Process CSV file using parallel streams
    public ProcessingResult<Map<String, Double>> processSalesDataParallel(
            Path csvFile) throws IOException {
        
        long startTime = System.currentTimeMillis();
        List<String> errors = Collections.synchronizedList(new ArrayList<>());
        
        // Read all lines and process in parallel
        List<SalesRecord> records;
        try (Stream<String> lines = Files.lines(csvFile).skip(1)) { // Skip header
            records = lines
                .parallel()
                .map(line -> parseRecord(line, errors))
                .filter(Objects::nonNull)
                .toList();
        }
        
        // Aggregate by region in parallel
        Map<String, Double> salesByRegion = records.parallelStream()
            .collect(Collectors.groupingByConcurrent(
                SalesRecord::region,
                Collectors.summingDouble(SalesRecord::amount)
            ));
        
        long duration = System.currentTimeMillis() - startTime;
        return new ProcessingResult<>(salesByRegion, records.size(), duration, errors);
    }
    
    private SalesRecord parseRecord(String line, List<String> errors) {
        try {
            String[] parts = line.split(",");
            return new SalesRecord(
                parts[0].trim(),
                parts[1].trim(),
                Double.parseDouble(parts[2].trim()),
                Integer.parseInt(parts[3].trim()),
                parts[4].trim()
            );
        } catch (Exception e) {
            errors.add("Failed to parse: " + line + " - " + e.getMessage());
            return null;
        }
    }
    
    // Fork/Join implementation for complex processing
    public static class DataProcessingTask extends RecursiveTask<Map<String, Statistics>> {
        private static final int THRESHOLD = 1000;
        private final List<SalesRecord> data;
        private final int start;
        private final int end;
        
        public DataProcessingTask(List<SalesRecord> data, int start, int end) {
            this.data = data;
            this.start = start;
            this.end = end;
        }
        
        @Override
        protected Map<String, Statistics> compute() {
            if (end - start <= THRESHOLD) {
                // Process directly
                return processSequentially();
            }
            
            // Split into subtasks
            int mid = (start + end) / 2;
            DataProcessingTask leftTask = new DataProcessingTask(data, start, mid);
            DataProcessingTask rightTask = new DataProcessingTask(data, mid, end);
            
            // Fork left task
            leftTask.fork();
            
            // Compute right task directly
            Map<String, Statistics> rightResult = rightTask.compute();
            
            // Join left task
            Map<String, Statistics> leftResult = leftTask.join();
            
            // Merge results
            return mergeResults(leftResult, rightResult);
        }
        
        private Map<String, Statistics> processSequentially() {
            Map<String, Statistics> result = new HashMap<>();
            
            for (int i = start; i < end; i++) {
                SalesRecord record = data.get(i);
                result.merge(
                    record.region(),
                    new Statistics(record.amount(), record.amount(), record.amount(), 1),
                    Statistics::merge
                );
            }
            
            return result;
        }
        
        private Map<String, Statistics> mergeResults(
                Map<String, Statistics> left, Map<String, Statistics> right) {
            Map<String, Statistics> merged = new HashMap<>(left);
            right.forEach((key, value) -> 
                merged.merge(key, value, Statistics::merge)
            );
            return merged;
        }
    }
    
    // Statistics aggregation class
    public record Statistics(double sum, double min, double max, long count) {
        public double average() {
            return count > 0 ? sum / count : 0;
        }
        
        public static Statistics merge(Statistics a, Statistics b) {
            return new Statistics(
                a.sum + b.sum,
                Math.min(a.min, b.min),
                Math.max(a.max, b.max),
                a.count + b.count
            );
        }
    }
    
    // Process using Fork/Join pool
    public ProcessingResult<Map<String, Statistics>> processWithForkJoin(
            List<SalesRecord> data) {
        
        long startTime = System.currentTimeMillis();
        
        ForkJoinPool pool = new ForkJoinPool(
            Runtime.getRuntime().availableProcessors()
        );
        
        try {
            DataProcessingTask task = new DataProcessingTask(data, 0, data.size());
            Map<String, Statistics> result = pool.invoke(task);
            
            long duration = System.currentTimeMillis() - startTime;
            return new ProcessingResult<>(result, data.size(), duration, List.of());
        } finally {
            pool.shutdown();
        }
    }
    
    // Batch processing with ExecutorService
    public <T, R> List<R> processBatches(
            List<T> data,
            int batchSize,
            Function<List<T>, R> batchProcessor) throws InterruptedException {
        
        ExecutorService executor = Executors.newFixedThreadPool(
            Runtime.getRuntime().availableProcessors()
        );
        
        List<Future<R>> futures = new ArrayList<>();
        
        // Submit batches
        for (int i = 0; i < data.size(); i += batchSize) {
            int end = Math.min(i + batchSize, data.size());
            List<T> batch = data.subList(i, end);
            
            futures.add(executor.submit(() -> batchProcessor.apply(batch)));
        }
        
        // Collect results
        List<R> results = new ArrayList<>();
        for (Future<R> future : futures) {
            try {
                results.add(future.get());
            } catch (ExecutionException e) {
                throw new RuntimeException("Batch processing failed", e.getCause());
            }
        }
        
        executor.shutdown();
        return results;
    }
    
    // Example usage
    public static void main(String[] args) throws Exception {
        ParallelDataProcessor processor = new ParallelDataProcessor();
        
        // Generate sample data
        List<SalesRecord> sampleData = generateSampleData(100_000);
        
        // Process with Fork/Join
        System.out.println("Processing with Fork/Join...");
        ProcessingResult<Map<String, Statistics>> result = 
            processor.processWithForkJoin(sampleData);
        
        System.out.println("Records processed: " + result.getRecordsProcessed());
        System.out.println("Processing time: " + result.getProcessingTimeMs() + "ms");
        System.out.println("\nResults by region:");
        result.getResult().forEach((region, stats) -> 
            System.out.printf("  %s: Total=%.2f, Avg=%.2f, Min=%.2f, Max=%.2f, Count=%d\n",
                region, stats.sum(), stats.average(), stats.min(), stats.max(), stats.count())
        );
        
        // Process in batches
        System.out.println("\nProcessing in batches...");
        List<Long> batchCounts = processor.processBatches(
            sampleData,
            10_000,
            batch -> (long) batch.size()
        );
        System.out.println("Batch counts: " + batchCounts);
    }
    
    private static List<SalesRecord> generateSampleData(int count) {
        String[] regions = {"North", "South", "East", "West", "Central"};
        Random random = new Random();
        List<SalesRecord> data = new ArrayList<>(count);
        
        for (int i = 0; i < count; i++) {
            data.add(new SalesRecord(
                "PROD-" + (i % 100),
                regions[random.nextInt(regions.length)],
                random.nextDouble() * 1000,
                random.nextInt(100) + 1,
                "2024-01-" + String.format("%02d", random.nextInt(28) + 1)
            ));
        }
        return data;
    }
}
```

---

### 10.3 Simple Web Server Simulation

A multithreaded web server that handles multiple client connections concurrently.

```
┌─────────────────────────────────────────────────────────────────────┐
│              SIMPLE WEB SERVER ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                    Main Server Thread                         │    │
│  │                          │                                    │    │
│  │               ServerSocket.accept()                           │    │
│  │                          │                                    │    │
│  │        ┌─────────────────┼─────────────────┐                 │    │
│  │        │                 │                 │                 │    │
│  │        ▼                 ▼                 ▼                 │    │
│  │   Connection 1      Connection 2      Connection 3           │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐    │
│  │                   Thread Pool (Workers)                       │    │
│  │  ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐  │    │
│  │  │ Worker 1 │   │ Worker 2 │   │ Worker 3 │   │ Worker N │  │    │
│  │  │          │   │          │   │          │   │          │  │    │
│  │  │ Process  │   │ Process  │   │ Process  │   │ Process  │  │    │
│  │  │ Request  │   │ Request  │   │ Request  │   │ Request  │  │    │
│  │  └──────────┘   └──────────┘   └──────────┘   └──────────┘  │    │
│  └─────────────────────────────────────────────────────────────┘    │
│                                                                      │
│  Request Queue (BlockingQueue) manages pending requests             │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.io.*;
import java.net.*;
import java.nio.file.*;
import java.time.*;
import java.time.format.*;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;

public class SimpleWebServer {
    
    private final int port;
    private final Path webRoot;
    private final ExecutorService threadPool;
    private final AtomicBoolean running = new AtomicBoolean(false);
    private ServerSocket serverSocket;
    
    // Statistics
    private final AtomicLong requestCount = new AtomicLong(0);
    private final AtomicLong totalResponseTime = new AtomicLong(0);
    private final ConcurrentHashMap<String, AtomicLong> statusCounts = new ConcurrentHashMap<>();
    
    public SimpleWebServer(int port, Path webRoot, int maxThreads) {
        this.port = port;
        this.webRoot = webRoot;
        this.threadPool = Executors.newFixedThreadPool(maxThreads, r -> {
            Thread t = new Thread(r);
            t.setName("WebServer-Worker-" + t.getId());
            return t;
        });
    }
    
    public void start() throws IOException {
        serverSocket = new ServerSocket(port);
        running.set(true);
        
        System.out.println("Server started on port " + port);
        System.out.println("Web root: " + webRoot.toAbsolutePath());
        
        // Accept connections in a loop
        while (running.get()) {
            try {
                Socket clientSocket = serverSocket.accept();
                threadPool.submit(new RequestHandler(clientSocket));
            } catch (SocketException e) {
                if (running.get()) {
                    System.err.println("Socket error: " + e.getMessage());
                }
            }
        }
    }
    
    public void stop() {
        running.set(false);
        try {
            if (serverSocket != null) {
                serverSocket.close();
            }
        } catch (IOException e) {
            System.err.println("Error closing server socket: " + e.getMessage());
        }
        
        threadPool.shutdown();
        try {
            if (!threadPool.awaitTermination(30, TimeUnit.SECONDS)) {
                threadPool.shutdownNow();
            }
        } catch (InterruptedException e) {
            threadPool.shutdownNow();
            Thread.currentThread().interrupt();
        }
        
        printStatistics();
    }
    
    // HTTP Request representation
    private static class HttpRequest {
        String method;
        String path;
        String version;
        Map<String, String> headers = new HashMap<>();
        
        static HttpRequest parse(BufferedReader reader) throws IOException {
            HttpRequest request = new HttpRequest();
            
            // Parse request line
            String requestLine = reader.readLine();
            if (requestLine == null || requestLine.isEmpty()) {
                return null;
            }
            
            String[] parts = requestLine.split(" ");
            if (parts.length >= 3) {
                request.method = parts[0];
                request.path = parts[1];
                request.version = parts[2];
            }
            
            // Parse headers
            String headerLine;
            while ((headerLine = reader.readLine()) != null && !headerLine.isEmpty()) {
                int colonIndex = headerLine.indexOf(':');
                if (colonIndex > 0) {
                    String name = headerLine.substring(0, colonIndex).trim();
                    String value = headerLine.substring(colonIndex + 1).trim();
                    request.headers.put(name.toLowerCase(), value);
                }
            }
            
            return request;
        }
    }
    
    // HTTP Response builder
    private static class HttpResponse {
        int statusCode;
        String statusText;
        Map<String, String> headers = new LinkedHashMap<>();
        byte[] body;
        
        HttpResponse(int code, String text) {
            this.statusCode = code;
            this.statusText = text;
            headers.put("Server", "SimpleJavaServer/1.0");
            headers.put("Date", DateTimeFormatter.RFC_1123_DATE_TIME
                .format(ZonedDateTime.now(ZoneOffset.UTC)));
        }
        
        void setBody(byte[] body, String contentType) {
            this.body = body;
            headers.put("Content-Type", contentType);
            headers.put("Content-Length", String.valueOf(body.length));
        }
        
        void write(OutputStream out) throws IOException {
            PrintWriter writer = new PrintWriter(out, false);
            
            // Status line
            writer.print("HTTP/1.1 " + statusCode + " " + statusText + "\r\n");
            
            // Headers
            for (Map.Entry<String, String> header : headers.entrySet()) {
                writer.print(header.getKey() + ": " + header.getValue() + "\r\n");
            }
            writer.print("\r\n");
            writer.flush();
            
            // Body
            if (body != null) {
                out.write(body);
            }
            out.flush();
        }
    }
    
    // Request handler (runs in thread pool)
    private class RequestHandler implements Runnable {
        private final Socket clientSocket;
        
        RequestHandler(Socket socket) {
            this.clientSocket = socket;
        }
        
        @Override
        public void run() {
            long startTime = System.currentTimeMillis();
            String statusKey = "unknown";
            
            try (clientSocket;
                 BufferedReader reader = new BufferedReader(
                     new InputStreamReader(clientSocket.getInputStream()));
                 OutputStream out = clientSocket.getOutputStream()) {
                
                HttpRequest request = HttpRequest.parse(reader);
                HttpResponse response;
                
                if (request == null || request.method == null) {
                    response = createErrorResponse(400, "Bad Request");
                } else if (!"GET".equals(request.method)) {
                    response = createErrorResponse(405, "Method Not Allowed");
                } else {
                    response = handleGetRequest(request);
                }
                
                statusKey = String.valueOf(response.statusCode);
                response.write(out);
                
                // Log request
                logRequest(clientSocket.getInetAddress().getHostAddress(),
                          request, response, System.currentTimeMillis() - startTime);
                
            } catch (IOException e) {
                System.err.println("Error handling request: " + e.getMessage());
            } finally {
                requestCount.incrementAndGet();
                totalResponseTime.addAndGet(System.currentTimeMillis() - startTime);
                statusCounts.computeIfAbsent(statusKey, k -> new AtomicLong(0))
                           .incrementAndGet();
            }
        }
        
        private HttpResponse handleGetRequest(HttpRequest request) {
            // Normalize path
            String path = request.path;
            if (path.equals("/")) {
                path = "/index.html";
            }
            
            // Security: prevent directory traversal
            if (path.contains("..")) {
                return createErrorResponse(403, "Forbidden");
            }
            
            // Resolve file path
            Path filePath = webRoot.resolve(path.substring(1)).normalize();
            
            // Ensure file is within web root
            if (!filePath.startsWith(webRoot)) {
                return createErrorResponse(403, "Forbidden");
            }
            
            // Check if file exists
            if (!Files.exists(filePath) || Files.isDirectory(filePath)) {
                return createErrorResponse(404, "Not Found");
            }
            
            // Read and serve file
            try {
                byte[] content = Files.readAllBytes(filePath);
                String contentType = getContentType(filePath);
                
                HttpResponse response = new HttpResponse(200, "OK");
                response.setBody(content, contentType);
                return response;
                
            } catch (IOException e) {
                return createErrorResponse(500, "Internal Server Error");
            }
        }
        
        private HttpResponse createErrorResponse(int code, String message) {
            HttpResponse response = new HttpResponse(code, message);
            String body = String.format(
                "<html><head><title>%d %s</title></head>" +
                "<body><h1>%d %s</h1></body></html>",
                code, message, code, message
            );
            response.setBody(body.getBytes(), "text/html");
            return response;
        }
        
        private String getContentType(Path path) {
            String fileName = path.getFileName().toString().toLowerCase();
            if (fileName.endsWith(".html") || fileName.endsWith(".htm")) {
                return "text/html";
            } else if (fileName.endsWith(".css")) {
                return "text/css";
            } else if (fileName.endsWith(".js")) {
                return "application/javascript";
            } else if (fileName.endsWith(".json")) {
                return "application/json";
            } else if (fileName.endsWith(".png")) {
                return "image/png";
            } else if (fileName.endsWith(".jpg") || fileName.endsWith(".jpeg")) {
                return "image/jpeg";
            } else if (fileName.endsWith(".gif")) {
                return "image/gif";
            } else if (fileName.endsWith(".txt")) {
                return "text/plain";
            }
            return "application/octet-stream";
        }
        
        private void logRequest(String clientIp, HttpRequest request, 
                               HttpResponse response, long duration) {
            String logLine = String.format(
                "%s - [%s] \"%s %s %s\" %d %d %dms",
                clientIp,
                LocalDateTime.now().format(DateTimeFormatter.ISO_LOCAL_DATE_TIME),
                request != null ? request.method : "-",
                request != null ? request.path : "-",
                request != null ? request.version : "-",
                response.statusCode,
                response.body != null ? response.body.length : 0,
                duration
            );
            System.out.println(logLine);
        }
    }
    
    private void printStatistics() {
        System.out.println("\n=== Server Statistics ===");
        System.out.println("Total requests: " + requestCount.get());
        if (requestCount.get() > 0) {
            System.out.println("Average response time: " + 
                (totalResponseTime.get() / requestCount.get()) + "ms");
        }
        System.out.println("Status codes:");
        statusCounts.forEach((code, count) -> 
            System.out.println("  " + code + ": " + count.get())
        );
    }
    
    // Example usage
    public static void main(String[] args) throws Exception {
        Path webRoot = Paths.get("www");
        Files.createDirectories(webRoot);
        
        // Create sample index.html
        Path indexFile = webRoot.resolve("index.html");
        if (!Files.exists(indexFile)) {
            Files.writeString(indexFile, 
                "<html><body><h1>Welcome to Simple Java Web Server!</h1></body></html>");
        }
        
        SimpleWebServer server = new SimpleWebServer(8080, webRoot, 10);
        
        // Add shutdown hook
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            System.out.println("\nShutting down server...");
            server.stop();
        }));
        
        System.out.println("Press Ctrl+C to stop the server");
        server.start();
    }
}
```

---

### 10.4 Producer-Consumer using BlockingQueue

A complete producer-consumer implementation with multiple producers, consumers, and graceful shutdown.

```
┌─────────────────────────────────────────────────────────────────────┐
│              PRODUCER-CONSUMER ARCHITECTURE                          │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  PRODUCERS                    QUEUE                    CONSUMERS     │
│  ─────────                    ─────                    ─────────     │
│                                                                      │
│  ┌──────────┐            ┌───────────────┐          ┌──────────┐    │
│  │Producer 1│────────────│               │──────────│Consumer 1│    │
│  └──────────┘            │               │          └──────────┘    │
│                          │  Blocking     │                          │
│  ┌──────────┐      put() │    Queue      │ take()   ┌──────────┐    │
│  │Producer 2│────────────│               │──────────│Consumer 2│    │
│  └──────────┘            │  (bounded)    │          └──────────┘    │
│                          │               │                          │
│  ┌──────────┐            │  ┌─┬─┬─┬─┬─┐ │          ┌──────────┐    │
│  │Producer 3│────────────│  │ │ │ │ │ │ │──────────│Consumer 3│    │
│  └──────────┘            │  └─┴─┴─┴─┴─┘ │          └──────────┘    │
│                          │               │                          │
│                          └───────────────┘                          │
│                                                                      │
│  • Queue blocks producers when full (backpressure)                  │
│  • Queue blocks consumers when empty (wait for work)                │
│  • Poison pill pattern for graceful shutdown                        │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```java
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.atomic.*;
import java.time.*;

public class ProducerConsumerSystem<T> {
    
    // Configuration
    private final int queueCapacity;
    private final int numProducers;
    private final int numConsumers;
    
    // Components
    private final BlockingQueue<WorkItem<T>> queue;
    private final List<Producer<T>> producers = new ArrayList<>();
    private final List<Consumer<T>> consumers = new ArrayList<>();
    private final ExecutorService producerExecutor;
    private final ExecutorService consumerExecutor;
    
    // State
    private final AtomicBoolean running = new AtomicBoolean(false);
    private final CountDownLatch producersDone;
    private final CountDownLatch consumersDone;
    
    // Statistics
    private final AtomicLong itemsProduced = new AtomicLong(0);
    private final AtomicLong itemsConsumed = new AtomicLong(0);
    private final AtomicLong processingErrors = new AtomicLong(0);
    
    // Work item wrapper
    public static class WorkItem<T> {
        private final T item;
        private final boolean isPoisonPill;
        private final Instant createdAt;
        
        private WorkItem(T item, boolean isPoisonPill) {
            this.item = item;
            this.isPoisonPill = isPoisonPill;
            this.createdAt = Instant.now();
        }
        
        public static <T> WorkItem<T> of(T item) {
            return new WorkItem<>(item, false);
        }
        
        public static <T> WorkItem<T> poisonPill() {
            return new WorkItem<>(null, true);
        }
        
        public T getItem() { return item; }
        public boolean isPoisonPill() { return isPoisonPill; }
        public Instant getCreatedAt() { return createdAt; }
    }
    
    // Functional interfaces
    @FunctionalInterface
    public interface ItemProducer<T> {
        T produce() throws Exception;
    }
    
    @FunctionalInterface
    public interface ItemConsumer<T> {
        void consume(T item) throws Exception;
    }
    
    // Producer class
    public class Producer<T> implements Runnable {
        private final int id;
        private final ItemProducer<T> producer;
        private final AtomicLong itemCount = new AtomicLong(0);
        
        public Producer(int id, ItemProducer<T> producer) {
            this.id = id;
            this.producer = producer;
        }
        
        @Override
        public void run() {
            Thread.currentThread().setName("Producer-" + id);
            System.out.println("Producer-" + id + " started");
            
            try {
                while (running.get() && !Thread.currentThread().isInterrupted()) {
                    try {
                        T item = producer.produce();
                        if (item != null) {
                            // Use offer with timeout to check running flag periodically
                            boolean added = queue.offer(
                                WorkItem.of(item), 
                                100, 
                                TimeUnit.MILLISECONDS
                            );
                            if (added) {
                                itemCount.incrementAndGet();
                                itemsProduced.incrementAndGet();
                            }
                        }
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        break;
                    } catch (Exception e) {
                        System.err.println("Producer-" + id + " error: " + e.getMessage());
                        processingErrors.incrementAndGet();
                    }
                }
            } finally {
                System.out.println("Producer-" + id + " stopped. Items produced: " + itemCount.get());
                producersDone.countDown();
            }
        }
        
        public long getItemCount() { return itemCount.get(); }
    }
    
    // Consumer class
    public class Consumer<T> implements Runnable {
        private final int id;
        private final ItemConsumer<T> consumer;
        private final AtomicLong itemCount = new AtomicLong(0);
        
        public Consumer(int id, ItemConsumer<T> consumer) {
            this.id = id;
            this.consumer = consumer;
        }
        
        @Override
        public void run() {
            Thread.currentThread().setName("Consumer-" + id);
            System.out.println("Consumer-" + id + " started");
            
            try {
                while (true) {
                    WorkItem<T> workItem = queue.poll(100, TimeUnit.MILLISECONDS);
                    
                    if (workItem == null) {
                        // Check if we should stop (producers done and queue empty)
                        if (!running.get() && queue.isEmpty()) {
                            break;
                        }
                        continue;
                    }
                    
                    if (workItem.isPoisonPill()) {
                        // Re-add poison pill for other consumers
                        queue.put(workItem);
                        break;
                    }
                    
                    try {
                        consumer.consume(workItem.getItem());
                        itemCount.incrementAndGet();
                        itemsConsumed.incrementAndGet();
                    } catch (Exception e) {
                        System.err.println("Consumer-" + id + " error: " + e.getMessage());
                        processingErrors.incrementAndGet();
                    }
                }
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            } finally {
                System.out.println("Consumer-" + id + " stopped. Items consumed: " + itemCount.get());
                consumersDone.countDown();
            }
        }
        
        public long getItemCount() { return itemCount.get(); }
    }
    
    public ProducerConsumerSystem(int queueCapacity, int numProducers, int numConsumers) {
        this.queueCapacity = queueCapacity;
        this.numProducers = numProducers;
        this.numConsumers = numConsumers;
        
        this.queue = new ArrayBlockingQueue<>(queueCapacity);
        this.producerExecutor = Executors.newFixedThreadPool(numProducers);
        this.consumerExecutor = Executors.newFixedThreadPool(numConsumers);
        this.producersDone = new CountDownLatch(numProducers);
        this.consumersDone = new CountDownLatch(numConsumers);
    }
    
    public void start(ItemProducer<T> producerLogic, ItemConsumer<T> consumerLogic) {
        running.set(true);
        
        // Start consumers first
        for (int i = 0; i < numConsumers; i++) {
            Consumer<T> consumer = new Consumer<>(i, consumerLogic);
            consumers.add(consumer);
            consumerExecutor.submit(consumer);
        }
        
        // Start producers
        for (int i = 0; i < numProducers; i++) {
            Producer<T> producer = new Producer<>(i, producerLogic);
            producers.add(producer);
            producerExecutor.submit(producer);
        }
        
        // Start monitoring thread
        startMonitoring();
    }
    
    public void stop() throws InterruptedException {
        System.out.println("\nInitiating shutdown...");
        
        // Signal producers to stop
        running.set(false);
        
        // Wait for producers to finish
        System.out.println("Waiting for producers to finish...");
        producersDone.await(30, TimeUnit.SECONDS);
        producerExecutor.shutdown();
        
        // Add poison pill to signal consumers
        System.out.println("Sending poison pill to consumers...");
        queue.put(WorkItem.poisonPill());
        
        // Wait for consumers to finish
        System.out.println("Waiting for consumers to finish...");
        consumersDone.await(30, TimeUnit.SECONDS);
        consumerExecutor.shutdown();
        
        printStatistics();
    }
    
    private void startMonitoring() {
        Thread monitor = new Thread(() -> {
            while (running.get()) {
                try {
                    Thread.sleep(2000);
                    System.out.printf(
                        "[Monitor] Queue size: %d/%d, Produced: %d, Consumed: %d, Errors: %d\n",
                        queue.size(), queueCapacity,
                        itemsProduced.get(), itemsConsumed.get(),
                        processingErrors.get()
                    );
                } catch (InterruptedException e) {
                    break;
                }
            }
        }, "Monitor");
        monitor.setDaemon(true);
        monitor.start();
    }
    
    private void printStatistics() {
        System.out.println("\n=== Final Statistics ===");
        System.out.println("Total items produced: " + itemsProduced.get());
        System.out.println("Total items consumed: " + itemsConsumed.get());
        System.out.println("Processing errors: " + processingErrors.get());
        System.out.println("Items remaining in queue: " + queue.size());
    }
    
    // Example: Order Processing System
    public static void main(String[] args) throws InterruptedException {
        // Create an order processing system
        ProducerConsumerSystem<Order> system = new ProducerConsumerSystem<>(
            100,  // Queue capacity
            3,    // Number of producers
            5     // Number of consumers
        );
        
        AtomicLong orderIdGenerator = new AtomicLong(0);
        Random random = new Random();
        
        // Producer logic: Generate random orders
        ItemProducer<Order> orderGenerator = () -> {
            Thread.sleep(random.nextInt(100));  // Simulate varying production rate
            return new Order(
                orderIdGenerator.incrementAndGet(),
                "Customer-" + random.nextInt(1000),
                random.nextDouble() * 500 + 10,
                List.of("Item-" + random.nextInt(100))
            );
        };
        
        // Consumer logic: Process orders
        ItemConsumer<Order> orderProcessor = order -> {
            Thread.sleep(random.nextInt(200));  // Simulate processing time
            System.out.printf("Processed order #%d for %s, amount: $%.2f\n",
                order.id(), order.customerId(), order.amount());
        };
        
        // Start the system
        system.start(orderGenerator, orderProcessor);
        
        // Let it run for 10 seconds
        System.out.println("System running for 10 seconds...");
        Thread.sleep(10000);
        
        // Graceful shutdown
        system.stop();
    }
    
    // Order record for the example
    public record Order(long id, String customerId, double amount, List<String> items) {}
}
```

#### Key Concepts Demonstrated

```
┌─────────────────────────────────────────────────────────────────────┐
│            PRODUCER-CONSUMER KEY CONCEPTS                            │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  1. BLOCKING QUEUE OPERATIONS                                       │
│     • put() - blocks when queue is full                             │
│     • take() - blocks when queue is empty                           │
│     • offer(timeout) - bounded wait for space                       │
│     • poll(timeout) - bounded wait for item                         │
│                                                                      │
│  2. GRACEFUL SHUTDOWN PATTERN                                       │
│     • Signal producers to stop (running flag)                       │
│     • Wait for producers to complete                                │
│     • Poison pill to signal consumers                               │
│     • Wait for consumers to complete                                │
│                                                                      │
│  3. BACKPRESSURE                                                    │
│     • Bounded queue prevents memory exhaustion                      │
│     • Producers naturally slow down when queue fills                │
│     • System automatically balances production/consumption          │
│                                                                      │
│  4. MULTIPLE PRODUCERS/CONSUMERS                                    │
│     • BlockingQueue handles all synchronization                     │
│     • No explicit locks needed                                      │
│     • Scales with thread count                                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

#### Choosing the Right BlockingQueue

| Queue Type | Bounded | Ordering | Best For |
|------------|---------|----------|----------|
| `ArrayBlockingQueue` | Yes | FIFO | General purpose, bounded capacity |
| `LinkedBlockingQueue` | Optional | FIFO | Higher throughput, optional bounds |
| `PriorityBlockingQueue` | No | Priority | Processing by priority |
| `SynchronousQueue` | No (0) | N/A | Direct handoff, no buffering |
| `DelayQueue` | No | Delay | Scheduled tasks |
| `LinkedTransferQueue` | No | FIFO | Producer waits for consumer |

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────────┐
│                    JAVA THREAD QUICK REFERENCE                       │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  CREATING THREADS:                                                   │
│  ─────────────────                                                   │
│  // Extend Thread                                                    │
│  class MyThread extends Thread { public void run() {...} }           │
│                                                                      │
│  // Implement Runnable                                               │
│  Thread t = new Thread(() -> { ... });                               │
│                                                                      │
│  // Callable with Future                                             │
│  Future<T> f = executor.submit(() -> { return value; });             │
│                                                                      │
│  LIFECYCLE:                                                          │
│  ──────────                                                          │
│  NEW → RUNNABLE ⟷ BLOCKED/WAITING/TIMED_WAITING → TERMINATED        │
│                                                                      │
│  KEY METHODS:                                                        │
│  ────────────                                                        │
│  t.start()        - Start thread (calls run())                       │
│  t.join()         - Wait for thread to finish                        │
│  Thread.sleep(ms) - Pause current thread                             │
│  Thread.yield()   - Hint to pause current thread                     │
│  t.interrupt()    - Interrupt thread                                 │
│  t.isAlive()      - Check if thread is running                       │
│  t.setDaemon(b)   - Set as daemon thread                             │
│                                                                      │
│  THREAD STATES:                                                      │
│  ──────────────                                                      │
│  Thread.State.NEW, RUNNABLE, BLOCKED, WAITING,                       │
│  TIMED_WAITING, TERMINATED                                           │
│                                                                      │
│  BEST PRACTICES:                                                     │
│  ───────────────                                                     │
│  ✓ Prefer Runnable/Callable over extending Thread                   │
│  ✓ Use ExecutorService for thread pools                             │
│  ✓ Always handle InterruptedException properly                       │
│  ✓ Use join() to wait for thread completion                         │
│  ✓ Set daemon status before calling start()                         │
│  ✗ Don't call run() directly - use start()                          │
│  ✗ Don't use stop()/suspend()/resume() - deprecated                 │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

---

## Interview Questions

### Basic Level

**Q1: What is the difference between `start()` and `run()` methods?**
- `start()`: Creates a new thread and executes `run()` in that new thread
- `run()`: Just a normal method call, executes in the current thread
- Calling `run()` directly does NOT create a new thread

**Q2: Can we start a thread twice?**
No. Calling `start()` on an already started thread throws `IllegalThreadStateException`.

**Q3: What is a daemon thread?**
A daemon thread is a background thread that doesn't prevent JVM termination. When all user threads finish, the JVM exits, killing daemon threads.

**Q4: What's the difference between `Runnable` and `Callable`?**
- `Runnable.run()` returns void and cannot throw checked exceptions
- `Callable.call()` returns a value and can throw checked exceptions
- `Callable` is used with `ExecutorService` and returns `Future`

### Intermediate Level

**Q5: What happens when a thread is interrupted while sleeping?**
`InterruptedException` is thrown, and the interrupt flag is cleared. The thread should handle this exception appropriately.

**Q6: Explain thread states in Java.**
NEW → RUNNABLE → (BLOCKED/WAITING/TIMED_WAITING) → TERMINATED

**Q7: What is the difference between `isInterrupted()` and `Thread.interrupted()`?**
- `isInterrupted()`: Instance method, checks flag without clearing
- `Thread.interrupted()`: Static method, checks AND clears the flag

### Synchronization Questions

**Q8: What is a race condition?**
A race condition occurs when multiple threads access shared data simultaneously, and the final result depends on the timing/order of thread execution, leading to unpredictable behavior.

**Q9: What is the difference between `synchronized` method and `synchronized` block?**
- `synchronized` method: Locks the entire method on `this` (instance) or `ClassName.class` (static)
- `synchronized` block: Allows finer-grained control, can lock on any object, and synchronize only critical sections

**Q10: What is a deadlock? How can you prevent it?**
Deadlock occurs when two or more threads are blocked forever, waiting for each other. Prevention strategies:
- Lock ordering: Always acquire locks in the same order
- Lock timeout: Use `tryLock()` with timeout
- Avoid nested locks when possible

**Q11: What is the difference between `notify()` and `notifyAll()`?**
- `notify()`: Wakes up ONE arbitrary thread from the wait set
- `notifyAll()`: Wakes up ALL threads in the wait set
- `notifyAll()` is safer but has more overhead

**Q12: Why must `wait()` be called inside a `synchronized` block?**
Because `wait()` releases the object's monitor lock. If not in a synchronized block, there's no lock to release, causing `IllegalMonitorStateException`.

**Q13: Why should `wait()` be called in a while loop instead of an if statement?**
To handle spurious wakeups. A thread can wake up without being notified 0A while loop ensures the condition is re-checked after waking up.

**Q14: What is the difference between `wait()` and `sleep()`?**
| wait() | sleep() |
|--------|--------|
| Releases lock | Does NOT release lock |
| Must be in synchronized block | Can be called anywhere |
| Object method | Thread static method |
| Woken by notify() | Woken by timeout |

**Q15: Explain livelock vs deadlock.**
- Deadlock: Threads are BLOCKED, waiting for each other forever
- Livelock: Threads are ACTIVE but keep responding to each other without making progress

### Advanced Concurrency Questions

**Q16: What is the difference between `execute()` and `submit()` in ExecutorService?**
- `execute()`: Takes Runnable, returns void, exceptions are sent to UncaughtExceptionHandler
- `submit()`: Takes Runnable/Callable, returns Future, exceptions are wrapped in ExecutionException

**Q17: What are the different types of thread pools?**
- `FixedThreadPool`: Fixed number of threads, unbounded queue
- `CachedThreadPool`: Creates threads as needed, removes idle threads after 60s
- `SingleThreadExecutor`: One thread, sequential execution
- `ScheduledThreadPool`: For scheduled/periodic tasks
- `WorkStealingPool`: Uses ForkJoinPool, work-stealing algorithm

**Q18: What is the purpose of `volatile` keyword?**
- Ensures visibility of changes across threads
- Prevents instruction reordering
- Does NOT provide atomicity for compound operations (like `i++`)

**Q19: What is Compare-And-Swap (CAS)?**
CAS is an atomic CPU instruction that compares a memory location's value with an expected value and, if they match, replaces it with a new value. It's the foundation of lock-free programming and used by Atomic classes.

**Q20: What is the ABA problem?**
In CAS, a value changes from A to B and back to A. A thread using CAS may incorrectly think nothing changed. Solution: Use `AtomicStampedReference` with version stamps.

**Q21: What is the difference between `CountDownLatch` and `CyclicBarrier`?**
| CountDownLatch | CyclicBarrier |
|----------------|---------------|
| One-time use | Reusable |
| Counts down to zero | Counts up to parties |
| Threads wait for events | Threads wait for each other |
| No barrier action | Optional barrier action |

**Q22: When would you use `ReentrantLock` over `synchronized`?**
- Need fairness (FIFO ordering)
- Need try-lock with timeout
- Need interruptible lock acquisition
- Need multiple condition variables
- Need to query lock status

**Q23: What is the difference between `thenApply()` and `thenCompose()` in CompletableFuture?**
- `thenApply()`: Transforms result, like `map()`. If function returns `CompletableFuture<T>`, result is `CompletableFuture<CompletableFuture<T>>`
- `thenCompose()`: Flattens nested futures, like `flatMap()`. Returns `CompletableFuture<T>`

**Q24: How does Fork/Join framework work?**
- Uses work-stealing algorithm where idle threads steal tasks from busy threads
- Tasks are split recursively until small enough to compute directly
- Uses `RecursiveTask` (with result) or `RecursiveAction` (void)
- Worker threads have their own deques - push/pop from bottom, steal from top

**Q25: What is the difference between `LongAdder` and `AtomicLong`?**
- `AtomicLong`: Single variable, all threads compete via CAS
- `LongAdder`: Distributed cells, reduces contention under high load
- Use `LongAdder` for counters with high contention
- Use `AtomicLong` when you need the exact current value frequently

---
