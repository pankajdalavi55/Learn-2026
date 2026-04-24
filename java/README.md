# Java — Basic to Advanced (FAANG-Level) Guide

A comprehensive collection of Java guides covering core fundamentals through advanced topics expected at Google, Meta, Amazon, Apple, Netflix, and other top-tier engineering organizations. Written from the perspective of a senior Java engineer with production experience at scale.

---

## Who Is This For?

- Backend engineers with 3+ years of Java experience preparing for FAANG interviews
- Senior developers looking to solidify their understanding of JVM internals, concurrency, and performance
- Anyone who wants a single, well-organized reference covering Java end-to-end (excluding DSA and System Design)

---

## Guides

| # | Guide | Key Topics |
|---|-------|------------|
| 01 | [Core Java & OOP](01-Java-Core-and-OOP-Guide.md) | JVM/JDK/JRE, data types, OOP pillars, SOLID principles, immutability, enums, keywords |
| 02 | [Generics & Type System](02-Java-Generics-and-Type-System.md) | Generics, wildcards, PECS, type erasure, bounded types, recursive bounds |
| 03 | [Exception Handling](03-Java-Exception-Handling-Guide.md) | Exception hierarchy, checked vs unchecked, try-with-resources, best practices, anti-patterns |
| 04 | [Collections Deep Dive](04-Java-Collections-Deep-Dive.md) | HashMap internals, ConcurrentHashMap, List/Set/Map/Queue, fail-fast iterators, choosing collections |
| 05 | [Streams & Functional Programming](05-Java-Streams-and-Functional-Programming.md) | Lambdas, Stream API, Optional, collectors, parallel streams, method references |
| 06 | [Multithreading & Concurrency](06-Java-Multithreading-and-Concurrency.md) | Threads, synchronized, ExecutorService, CompletableFuture, locks, atomics, virtual threads |
| 07 | [Memory Model & JVM Internals](07-Java-Memory-Model-and-JVM-Internals.md) | Heap/stack, GC algorithms (G1, ZGC), JMM, class loading, JIT, monitoring tools |
| 08 | [Modern Features (Java 8-21)](08-Java-Modern-Features-8-to-21.md) | Records, sealed classes, pattern matching, text blocks, var, modules, version-by-version changes |
| 09 | [Design Patterns](09-Java-Design-Patterns-Guide.md) | All 23 GoF patterns, mermaid UML diagrams, JDK/Spring real-world usage, anti-patterns |
| 10 | [Testing](10-Java-Testing-Guide.md) | JUnit 5, Mockito, Spring Boot testing, TDD, TestContainers, mutation testing |
| 11 | [Performance Tuning](11-Java-Performance-Tuning-Guide.md) | JMH benchmarking, profiling, string/collection performance, HikariCP, flame graphs |
| 12 | [Serialization, I/O & Networking](12-Java-Serialization-IO-and-Networking.md) | java.io, NIO, Jackson/Gson, HTTP Client, socket programming |

---

## Reading Order

The guides are numbered in a recommended reading order, progressing from fundamentals to advanced topics. Each guide is self-contained, so you can jump to any topic directly.

```text
Fundamentals          Advanced Language        JVM & Runtime           Engineering Practice
─────────────         ──────────────────       ─────────────           ────────────────────
01 Core & OOP    ──►  05 Streams & FP     ──►  07 JVM Internals   ──►  10 Testing
02 Generics      ──►  06 Concurrency      ──►  08 Modern Features ──►  11 Performance
03 Exceptions    ──►  09 Design Patterns                              12 Serialization & I/O
04 Collections
```

---

## What Each Guide Includes

- Production-grade Java code examples
- Mermaid diagrams for architecture, hierarchies, and flows
- Comparison tables for quick reference
- Google/FAANG perspective callouts and best practices
- Interview-focused summary with rapid-fire Q&A tables

---
