# Spring Boot — Complete Learning Guide

> A comprehensive, interview-ready reference covering Spring Boot from fundamentals to production deployment.  
> **Scope:** Spring Boot core features, auto-configuration, starters, Actuator, testing, and production readiness.  
> **Prerequisites:** Familiarity with Spring Framework fundamentals (see companion Spring Framework guide).

---

## Table of Contents

1. [Introduction & Philosophy](#1-introduction--philosophy)
2. [Auto-Configuration Deep Dive](#2-auto-configuration-deep-dive)
3. [Starters & Dependency Management](#3-starters--dependency-management)
4. [Externalized Configuration](#4-externalized-configuration)
5. [Spring Boot Web & REST](#5-spring-boot-web--rest)
6. [Data Access with Spring Boot](#6-data-access-with-spring-boot)
7. [Spring Boot Actuator](#7-spring-boot-actuator)
8. [Testing in Spring Boot](#8-testing-in-spring-boot)
9. [Logging & Observability](#9-logging--observability)
10. [Security in Spring Boot](#10-security-in-spring-boot)
11. [Production Readiness](#11-production-readiness)
12. [Docker & Kubernetes Deployment](#12-docker--kubernetes-deployment)
13. [Performance Tuning](#13-performance-tuning)
14. [Best Practices & Anti-Patterns](#14-best-practices--anti-patterns)
15. [Interview Questions by Experience Level](#15-interview-questions-by-experience-level)

---

## 1. Introduction & Philosophy

### 1.1 What is Spring Boot?

Spring Boot is an **opinionated, convention-over-configuration** framework built on top of Spring Framework. It dramatically simplifies Spring application development by providing:

- **Auto-configuration** — Automatically configures beans based on classpath
- **Starter dependencies** — Curated dependency sets that "just work" together
- **Embedded servers** — No need for external application servers
- **Production-ready features** — Health checks, metrics, externalized config out of the box

### 1.1.1 The Philosophy: Convention Over Configuration

**Convention over Configuration (CoC)** is a software design paradigm that reduces the number of decisions developers need to make without losing flexibility.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    Traditional Spring (Pre-Boot)                        │
│                                                                         │
│  Developer must decide:                                                 │
│  • Which versions of 40+ dependencies are compatible?                   │
│  • How to configure DataSource, EntityManager, TransactionManager?      │
│  • Where to place configuration files?                                  │
│  • How to package and deploy the application?                           │
│  • Which server to install and configure?                               │
│                                                                         │
│  Result: Days of setup before writing business logic                    │
└─────────────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────────────┐
│                    Spring Boot (Convention over Config)                 │
│                                                                         │
│  Boot decides (but allows override):                                    │
│  • Curated dependency versions via starters                             │
│  • Default configurations for common scenarios                          │
│  • Standard locations for config files                                  │
│  • Executable JAR with embedded server                                  │
│  • Sensible production defaults                                         │
│                                                                         │
│  Result: Working app in minutes, customize as needed                    │
└─────────────────────────────────────────────────────────────────────────┘
```

**Core Principles of CoC:**

| Principle | Spring Boot Implementation |
|---|---|
| **Sensible Defaults** | HikariCP is default connection pool (fastest) |
| **Zero Configuration** | Add `spring-boot-starter-web` → web server works |
| **Override When Needed** | Any default can be customized via properties |
| **Fail Fast** | Misconfiguration detected at startup, not runtime |
| **Explicit Over Implicit** | Conditions clearly documented |

### 1.1.2 Opinionated Defaults: The "Right" Choices

Spring Boot makes opinionated choices based on production experience:

```
┌───────────────────────────────────────────────────────────────────┐
│                    Spring Boot's Opinions                         │
├───────────────────────────────────────────────────────────────────┤
│  Connection Pool:    HikariCP (fastest, most reliable)            │
│  JSON Processor:     Jackson (industry standard)                  │
│  Web Server:         Tomcat (balanced, well-documented)           │
│  Logging:            Logback + SLF4J (proven combination)         │
│  Testing:            JUnit 5 + Mockito + AssertJ                  │
│  Build:              Maven/Gradle with wrapper scripts            │
│  Packaging:          Executable JAR (single artifact)             │
│  Config Format:      YAML/Properties (simple, portable)           │
└───────────────────────────────────────────────────────────────────┘
```

**Why Opinions Matter:**

1. **Reduced Decision Fatigue** — Fewer choices = faster development
2. **Community Alignment** — Everyone uses similar patterns
3. **Better Documentation** — Common setup = more resources available
4. **Production Tested** — Defaults proven in real-world scenarios
5. **Upgrade Safety** — Version combinations tested together

### 1.1.3 The "Bootiful" Development Experience

```plaintext
Traditional Java Web Development Timeline:
├── Day 1-3:    Setup build, download dependencies, resolve conflicts
├── Day 4-5:    Configure app server (Tomcat/WebSphere/JBoss)
├── Day 6-7:    Wire up Spring XML configurations
├── Day 8-9:    Configure database, transactions, security
├── Day 10:     Finally write first business feature
└── Total:      10 days before productive work

Spring Boot Development Timeline:
├── Minute 1-5:  Generate project at start.spring.io
├── Minute 5-10: Write first REST endpoint and entity
├── Minute 10-15: Run `mvn spring-boot:run`, test with curl
└── Total:       15 minutes to first working feature
```

**Historical Context:**
| Year | Event |
|---|---|
| 2012 | Phil Webb starts the project after a request to support containerless deployments |
| 2014 | Spring Boot 1.0 released |
| 2018 | Spring Boot 2.0 with WebFlux, Kotlin support |
| 2022 | Spring Boot 3.0 with Jakarta EE 9+, GraalVM native support |
| 2023 | Spring Boot 3.1+ with virtual threads, improved observability |

### 1.2 Spring Boot vs Spring Framework

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Spring Boot                                  │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │              Auto-Configuration Layer                       │    │
│  │  • @EnableAutoConfiguration                                │    │
│  │  • Conditional bean registration                            │    │
│  │  • Opinionated defaults                                     │    │
│  └────────────────────────────────────────────────────────────┘    │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │              Starter Dependencies                           │    │
│  │  • spring-boot-starter-web                                  │    │
│  │  • spring-boot-starter-data-jpa                             │    │
│  │  • Compatible version management                            │    │
│  └────────────────────────────────────────────────────────────┘    │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │              Production Features                            │    │
│  │  • Actuator, Metrics, Health checks                         │    │
│  │  • Embedded servers                                         │    │
│  └────────────────────────────────────────────────────────────┘    │
│  ┌────────────────────────────────────────────────────────────┐    │
│  │                 Spring Framework                            │    │
│  │  • IoC Container, DI, AOP                                   │    │
│  │  • Spring MVC, Spring Data, Spring Security                 │    │
│  │  • Transaction Management                                   │    │
│  └────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────┘
```

| Aspect | Spring Framework | Spring Boot |
|---|---|---|
| **Configuration** | Manual (XML/Java) | Auto-configured |
| **Dependencies** | Individual artifacts | Curated starters |
| **Server** | Deploy to external | Embedded (Tomcat/Jetty/Undertow) |
| **Startup** | Requires extensive setup | `main()` method with `@SpringBootApplication` |
| **Packaging** | WAR file | Executable JAR (fat JAR) |
| **Monitoring** | Manual setup | Actuator out of the box |

### 1.3 The Twelve-Factor App & Spring Boot

**What is the Twelve-Factor App?**

The Twelve-Factor App methodology was created by Heroku developers as a set of best practices for building modern, cloud-native applications. These factors address common problems in distributed systems and container deployments.

**Why It Matters:**
- Applications built following these principles are **portable** across cloud providers
- They **scale horizontally** without architectural changes
- They are **resilient** to infrastructure failures
- They support **continuous deployment** workflows

Spring Boot naturally aligns with [12-Factor App](https://12factor.net) methodology:

| Factor | Spring Boot Support |
|---|---|
| **I. Codebase** | Single deployable artifact per service |
| **II. Dependencies** | Explicit via Maven/Gradle with starters |
| **III. Config** | Externalized via `application.properties/yml`, env vars |
| **IV. Backing Services** | Connection strings externalized |
| **V. Build, Release, Run** | Fat JAR, distinct stages |
| **VI. Processes** | Stateless; use external stores |
| **VII. Port Binding** | Embedded server, `server.port` |
| **VIII. Concurrency** | Scale out via multiple instances |
| **IX. Disposability** | Fast startup, graceful shutdown |
| **X. Dev/Prod Parity** | Profiles for environment-specific config |
| **XI. Logs** | Stdout/stderr logging (Logback default) |
| **XII. Admin Processes** | Spring Batch, scheduled tasks |

**Deep Dive on Critical Factors:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│              Factor III: Config — Store Config in Environment           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ❌ WRONG: Hardcoded values                                            │
│  String dbUrl = "jdbc:postgresql://localhost:5432/mydb";               │
│                                                                         │
│  ❌ WRONG: Config in code (even in constants)                          │
│  public static final String DB_URL = "...";                            │
│                                                                         │
│  ✅ RIGHT: Externalized, environment-specific                          │
│  spring.datasource.url=${DATABASE_URL}                                  │
│                                                                         │
│  Benefits:                                                              │
│  • Same artifact runs in dev, staging, prod                             │
│  • Secrets never in version control                                     │
│  • Config can change without recompilation                              │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│              Factor VI: Processes — Execute App as Stateless            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ❌ WRONG: Storing state in memory                                     │
│  private Map<String, ShoppingCart> carts = new HashMap<>();             │
│  → Lost on restart, not shared across instances                         │
│                                                                         │
│  ✅ RIGHT: External state stores                                       │
│  @Autowired RedisTemplate<String, ShoppingCart> cartStore;              │
│  → Survives restarts, shared across instances                           │
│                                                                         │
│  Why it matters:                                                        │
│  • Scale horizontally by adding instances                               │
│  • Any instance can handle any request                                  │
│  • Instance failure doesn't lose data                                   │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│              Factor IX: Disposability — Fast Startup & Shutdown         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Startup:                                                               │
│  • Spring Boot 3.x: ~1-2 seconds typical                                │
│  • With GraalVM native: ~50-100 milliseconds                            │
│  • Enables rapid scaling during traffic spikes                          │
│                                                                         │
│  Shutdown:                                                              │
│  • Graceful: Complete in-flight requests                                │
│  • Clean: Release DB connections, close files                           │
│  • Fast: Don't block deployment pipelines                               │
│                                                                         │
│  Spring Boot support:                                                   │
│  server.shutdown=graceful                                               │
│  spring.lifecycle.timeout-per-shutdown-phase=30s                        │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.4 Creating a Spring Boot Application

#### 1.4.1 Using Spring Initializr

The fastest way: [start.spring.io](https://start.spring.io)

```bash
# CLI with HTTPie
curl https://start.spring.io/starter.tgz \
  -d type=maven-project \
  -d language=java \
  -d bootVersion=3.2.0 \
  -d dependencies=web,data-jpa,postgresql,actuator \
  -d groupId=com.example \
  -d artifactId=demo \
  | tar -xzvf -
```

#### 1.4.2 Minimal Application Structure

```
my-app/
├── src/
│   ├── main/
│   │   ├── java/
│   │   │   └── com/example/demo/
│   │   │       ├── DemoApplication.java
│   │   │       ├── controller/
│   │   │       ├── service/
│   │   │       ├── repository/
│   │   │       └── model/
│   │   └── resources/
│   │       ├── application.properties
│   │       ├── application-dev.properties
│   │       ├── application-prod.properties
│   │       └── static/
│   │       └── templates/
│   └── test/
│       └── java/
│           └── com/example/demo/
│               └── DemoApplicationTests.java
├── pom.xml (or build.gradle)
└── mvnw / mvnw.cmd
```

#### 1.4.3 The Main Class

```java
package com.example.demo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication  // = @Configuration + @EnableAutoConfiguration + @ComponentScan
public class DemoApplication {

    public static void main(String[] args) {
        SpringApplication.run(DemoApplication.class, args);
    }
}
```

**What `@SpringBootApplication` Does:**

```java
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@SpringBootConfiguration      // → @Configuration
@EnableAutoConfiguration      // → Trigger auto-configuration
@ComponentScan                // → Scan from this package down
public @interface SpringBootApplication { }
```

---

## 2. Auto-Configuration Deep Dive

### 2.0 Understanding Auto-Configuration Theory

**What Problem Does Auto-Configuration Solve?**

In traditional Spring, configuring a simple web application with database access required extensive boilerplate:

```java
// Traditional Spring — Manual DataSource configuration
@Configuration
public class DataSourceConfig {
    
    @Bean
    public DataSource dataSource() {
        HikariConfig config = new HikariConfig();
        config.setJdbcUrl("jdbc:postgresql://localhost/mydb");
        config.setUsername("user");
        config.setPassword("pass");
        config.setMaximumPoolSize(10);
        config.setMinimumIdle(5);
        // ... 20+ more settings
        return new HikariDataSource(config);
    }
    
    @Bean
    public LocalContainerEntityManagerFactoryBean entityManagerFactory(
            DataSource dataSource) {
        LocalContainerEntityManagerFactoryBean em = new LocalContainerEntityManagerFactoryBean();
        em.setDataSource(dataSource);
        em.setPackagesToScan("com.example.entity");
        // ... more configuration
        return em;
    }
    
    @Bean
    public PlatformTransactionManager transactionManager(
            EntityManagerFactory emf) {
        return new JpaTransactionManager(emf);
    }
}
```

**With Spring Boot Auto-Configuration:**

```yaml
# application.yml — That's it!
spring:
  datasource:
    url: jdbc:postgresql://localhost/mydb
    username: user
    password: pass
```

Spring Boot detects:
1. HikariCP on classpath → creates `DataSource`
2. JPA on classpath → creates `EntityManagerFactory`
3. Both exist → creates `TransactionManager`

### 2.0.1 The Magic Behind Auto-Configuration

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AUTO-CONFIGURATION MENTAL MODEL                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│    Auto-Config = Pattern Matching + Conditional Creation                 │
│                                                                          │
│    IF (HikariCP.class is on classpath)                                  │
│       AND (no DataSource bean exists)                                   │
│       AND (datasource.url property is set)                              │
│    THEN                                                                 │
│       → Create HikariDataSource bean with properties                    │
│                                                                          │
│    The "magic" is just well-organized IF-THEN rules!                   │
│                                                                          │
└─────────────────────────────────────────────────────────────────────────┘
```

**Key Insight:** Auto-configuration is NOT magic — it's a sophisticated set of conditional rules that:
- Inspect the classpath for available libraries
- Check if beans already exist
- Read configuration properties
- Apply sensible defaults

### 2.0.2 Classpath Scanning vs Auto-Configuration

| Mechanism | Purpose | How It Works |
|---|---|---|
| **Component Scanning** | Find YOUR beans | Scans packages for `@Component`, `@Service`, etc. |
| **Auto-Configuration** | Configure FRAMEWORK beans | Loads pre-defined `@Configuration` classes conditionally |

```
┌─────────────────────────────────────────────────────────────────────────┐
│                        APPLICATION STARTUP                               │
│                              │                                           │
│        ┌─────────────────────┴─────────────────────┐                    │
│        │                                           │                    │
│        ▼                                           ▼                    │
│   Component Scan                            Auto-Configuration          │
│   (Your Code)                               (Framework Code)            │
│        │                                           │                    │
│        ▼                                           ▼                    │
│   @Service, @Repository                    @AutoConfiguration           │
│   @Controller, @Component                  classes from JARs            │
│        │                                           │                    │
│        └─────────────────────┬─────────────────────┘                    │
│                              │                                           │
│                              ▼                                           │
│                    ApplicationContext                                    │
│                    (All beans merged)                                    │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.1 How Auto-Configuration Works

Auto-configuration is Spring Boot's "magic" — **automatic bean registration based on classpath contents and existing beans**.

```
┌─────────────────────────────────────────────────────────────────┐
│                    Application Startup                          │
│                           │                                     │
│        ┌──────────────────▼──────────────────┐                  │
│        │  @EnableAutoConfiguration           │                  │
│        │  triggers AutoConfigurationImport   │                  │
│        └──────────────────┬──────────────────┘                  │
│                           │                                     │
│        ┌──────────────────▼──────────────────┐                  │
│        │  Load META-INF/spring/              │                  │
│        │  org.springframework.boot.autoconfigure│               │
│        │  .AutoConfiguration.imports         │                  │
│        └──────────────────┬──────────────────┘                  │
│                           │                                     │
│        ┌──────────────────▼──────────────────┐                  │
│        │  For each auto-configuration class: │                  │
│        │  • Check @Conditional annotations   │                  │
│        │  • If conditions match → register   │                  │
│        └──────────────────┬──────────────────┘                  │
│                           │                                     │
│        ┌──────────────────▼──────────────────┐                  │
│        │  Beans registered in context        │                  │
│        └─────────────────────────────────────┘                  │
└─────────────────────────────────────────────────────────────────┘
```

### 2.2 The @Conditional Family

Auto-configuration relies heavily on conditional annotations:

| Annotation | Condition |
|---|---|
| `@ConditionalOnClass` | Class is on classpath |
| `@ConditionalOnMissingClass` | Class is NOT on classpath |
| `@ConditionalOnBean` | Bean of type exists in context |
| `@ConditionalOnMissingBean` | Bean of type does NOT exist |
| `@ConditionalOnProperty` | Property has specific value |
| `@ConditionalOnResource` | Resource exists (e.g., file) |
| `@ConditionalOnWebApplication` | Is a web application |
| `@ConditionalOnNotWebApplication` | Is NOT a web application |
| `@ConditionalOnExpression` | SpEL expression evaluates to true |
| `@ConditionalOnJava` | Specific Java version |
| `@ConditionalOnSingleCandidate` | Single bean or @Primary exists |
| `@ConditionalOnCloudPlatform` | Running on specific cloud (Kubernetes, CloudFoundry) |

### 2.3 Example: DataSource Auto-Configuration

```java
// Simplified version of actual Spring Boot code
@AutoConfiguration(after = DataSourceAutoConfiguration.class)
@ConditionalOnClass({ DataSource.class, EmbeddedDatabaseType.class })
@ConditionalOnMissingBean(type = "io.r2dbc.spi.ConnectionFactory")
@EnableConfigurationProperties(DataSourceProperties.class)
public class DataSourceAutoConfiguration {

    @Configuration(proxyBeanMethods = false)
    @Conditional(EmbeddedDatabaseCondition.class)
    @ConditionalOnMissingBean({ DataSource.class, XADataSource.class })
    @Import(EmbeddedDataSourceConfiguration.class)
    protected static class EmbeddedDatabaseConfiguration { }

    @Configuration(proxyBeanMethods = false)
    @Conditional(PooledDataSourceCondition.class)
    @ConditionalOnMissingBean({ DataSource.class, XADataSource.class })
    @Import({ DataSourceConfiguration.Hikari.class, 
              DataSourceConfiguration.Tomcat.class,
              DataSourceConfiguration.Dbcp2.class })
    protected static class PooledDataSourceConfiguration { }
}
```

**Reading this:**
1. Only activates if `DataSource` and `EmbeddedDatabaseType` classes are on classpath
2. Skips if R2DBC (reactive) is being used
3. If no `DataSource` bean exists AND embedded DB condition matches → create embedded
4. If no `DataSource` bean exists AND pooled condition matches → create HikariCP/Tomcat/DBCP2

### 2.4 Understanding Auto-Configuration Order

```java
@AutoConfiguration(
    before = SecurityAutoConfiguration.class,    // Run BEFORE this
    after = DataSourceAutoConfiguration.class    // Run AFTER this
)
public class MyAutoConfiguration { }
```

This ensures dependencies are configured in correct order.

### 2.5 Debugging Auto-Configuration

**1. Debug Mode:**
```properties
# application.properties
debug=true
```

Produces a CONDITIONS EVALUATION REPORT in logs showing:
- Positive matches (conditions met)
- Negative matches (why something wasn't configured)
- Exclusions
- Unconditional classes

**2. Actuator Endpoint:**
```bash
GET /actuator/conditions
```

**3. Common Debug Output:**

```
============================
CONDITIONS EVALUATION REPORT
============================

Positive matches:
-----------------
   DataSourceAutoConfiguration matched:
      - @ConditionalOnClass found required classes 'javax.sql.DataSource', 
        'org.springframework.jdbc.datasource.embedded.EmbeddedDatabaseType'
      
   DataSourceAutoConfiguration.PooledDataSourceConfiguration matched:
      - AnyNestedCondition 1 matched 1 did not; NestedCondition on 
        DataSourceAutoConfiguration.PooledDataSourceCondition found
        
Negative matches:
-----------------
   MongoAutoConfiguration:
      Did not match:
         - @ConditionalOnClass did not find required class 
           'com.mongodb.client.MongoClient'
```

### 2.6 Excluding Auto-Configurations

```java
// Method 1: In annotation
@SpringBootApplication(exclude = {
    DataSourceAutoConfiguration.class,
    SecurityAutoConfiguration.class
})
public class MyApp { }

// Method 2: In properties
spring.autoconfigure.exclude=\
  org.springframework.boot.autoconfigure.jdbc.DataSourceAutoConfiguration,\
  org.springframework.boot.autoconfigure.security.servlet.SecurityAutoConfiguration
```

### 2.7 Creating Custom Auto-Configuration

```java
// 1. Create the configuration class
@AutoConfiguration
@ConditionalOnClass(MyService.class)
@ConditionalOnMissingBean(MyService.class)
@EnableConfigurationProperties(MyServiceProperties.class)
public class MyServiceAutoConfiguration {

    @Bean
    @ConditionalOnProperty(prefix = "myservice", name = "enabled", havingValue = "true", matchIfMissing = true)
    public MyService myService(MyServiceProperties properties) {
        return new MyService(properties.getEndpoint(), properties.getTimeout());
    }
}

// 2. Properties class
@ConfigurationProperties(prefix = "myservice")
public class MyServiceProperties {
    private String endpoint = "http://localhost:8080";
    private Duration timeout = Duration.ofSeconds(30);
    
    // getters and setters
}

// 3. Register in META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports
com.example.MyServiceAutoConfiguration
```

### 2.8 Auto-Configuration Best Practices

| Practice | Why |
|---|---|
| Use `@ConditionalOnMissingBean` | Allow users to override |
| Use `@ConditionalOnClass` | Don't fail if library absent |
| Use `@ConfigurationProperties` | Type-safe configuration |
| Place after relevant configurations | Ensure dependencies exist |
| Provide sensible defaults | Zero-config should work |
| Document properties | Help users customize |

---

## 3. Starters & Dependency Management

### 3.1 What are Starters?

Starters are **curated sets of dependencies** that provide everything needed for a specific functionality. They follow the naming convention `spring-boot-starter-*`.

**Why Starters?**
- No hunting for compatible versions
- Transitive dependencies included
- Opinionated but overridable
- Single dependency for complex features

### 3.2 Common Starters

| Starter | Provides |
|---|---|
| `spring-boot-starter` | Core (logging, auto-config, YAML) |
| `spring-boot-starter-web` | Web MVC, embedded Tomcat, Jackson |
| `spring-boot-starter-webflux` | Reactive web with Netty |
| `spring-boot-starter-data-jpa` | JPA with Hibernate, Spring Data |
| `spring-boot-starter-data-mongodb` | MongoDB support |
| `spring-boot-starter-data-redis` | Redis support |
| `spring-boot-starter-security` | Spring Security |
| `spring-boot-starter-actuator` | Production monitoring |
| `spring-boot-starter-test` | JUnit 5, Mockito, AssertJ, Testcontainers |
| `spring-boot-starter-validation` | Bean Validation (Hibernate Validator) |
| `spring-boot-starter-cache` | Spring Cache abstraction |
| `spring-boot-starter-mail` | JavaMail |
| `spring-boot-starter-amqp` | RabbitMQ |
| `spring-boot-starter-oauth2-client` | OAuth2 client |
| `spring-boot-starter-oauth2-resource-server` | OAuth2 resource server |
| `spring-boot-starter-aop` | AspectJ AOP |

### 3.3 Starter Internals: What's in `spring-boot-starter-web`?

```xml
<!-- spring-boot-starter-web dependencies (simplified) -->
spring-boot-starter
├── spring-boot
├── spring-boot-autoconfigure
├── logback-classic
└── slf4j-api

spring-boot-starter-json
├── jackson-databind
├── jackson-datatype-jsr310
└── jackson-module-parameter-names

spring-boot-starter-tomcat
├── tomcat-embed-core
├── tomcat-embed-el
└── tomcat-embed-websocket

spring-web
spring-webmvc
```

### 3.4 Dependency Management with BOM

Spring Boot uses a **Bill of Materials (BOM)** to manage versions:

```xml
<!-- pom.xml -->
<parent>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-parent</artifactId>
    <version>3.2.0</version>
</parent>

<dependencies>
    <!-- No version needed — managed by parent -->
    <dependency>
        <groupId>org.springframework.boot</groupId>
        <artifactId>spring-boot-starter-web</artifactId>
    </dependency>
</dependencies>
```

Without parent (using BOM import):

```xml
<dependencyManagement>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-dependencies</artifactId>
            <version>3.2.0</version>
            <type>pom</type>
            <scope>import</scope>
        </dependency>
    </dependencies>
</dependencyManagement>
```

### 3.5 Overriding Dependency Versions

```xml
<!-- Override managed version -->
<properties>
    <postgresql.version>42.6.0</postgresql.version>
    <jackson-bom.version>2.15.0</jackson-bom.version>
</properties>

<!-- Or explicitly in dependency -->
<dependency>
    <groupId>org.postgresql</groupId>
    <artifactId>postgresql</artifactId>
    <version>42.6.0</version>  <!-- Explicit version overrides managed -->
</dependency>
```

### 3.6 Switching Embedded Server

```xml
<!-- Default: Tomcat. Switch to Jetty: -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-web</artifactId>
    <exclusions>
        <exclusion>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-tomcat</artifactId>
        </exclusion>
    </exclusions>
</dependency>
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-jetty</artifactId>
</dependency>

<!-- Or switch to Undertow -->
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-undertow</artifactId>
</dependency>
```

### 3.7 Creating Custom Starters

```
my-custom-starter/
├── my-service-spring-boot-autoconfigure/
│   ├── src/main/java/
│   │   └── com/example/autoconfigure/
│   │       ├── MyServiceAutoConfiguration.java
│   │       └── MyServiceProperties.java
│   ├── src/main/resources/
│   │   └── META-INF/
│   │       └── spring/
│   │           └── org.springframework.boot.autoconfigure.AutoConfiguration.imports
│   └── pom.xml
│
└── my-service-spring-boot-starter/
    └── pom.xml  <!-- Just depends on autoconfigure + actual library -->
```

**Naming Convention:**
- Official: `spring-boot-starter-{name}`
- Custom: `{name}-spring-boot-starter`

---

## 4. Externalized Configuration

### 4.1 Property Sources Hierarchy

Spring Boot loads configuration from multiple sources in this **priority order** (highest to lowest):

```
1.  Command line arguments (--server.port=9090)
2.  Java System properties (-Dserver.port=9090)
3.  OS environment variables (SERVER_PORT=9090)
4.  @PropertySource annotations
5.  application-{profile}.properties/yml
6.  application.properties/yml
7.  @ConfigurationProperties defaults
8.  SpringApplication.setDefaultProperties()
```

Higher priority sources override lower ones.

### 4.2 Property File Locations

Spring Boot searches for `application.properties` or `application.yml` in:

```
1. ./config/  (current directory subdirectory)
2. ./         (current directory)
3. classpath:/config/
4. classpath:/
```

### 4.3 YAML vs Properties

**Properties format:**
```properties
server.port=8080
spring.datasource.url=jdbc:postgresql://localhost/mydb
spring.datasource.username=user
spring.datasource.password=secret
app.features[0]=feature1
app.features[1]=feature2
```

**YAML format:**
```yaml
server:
  port: 8080

spring:
  datasource:
    url: jdbc:postgresql://localhost/mydb
    username: user
    password: secret

app:
  features:
    - feature1
    - feature2
```

**YAML Advantages:**
- Hierarchical structure (no repetition)
- Lists are cleaner
- Multi-document support (`---`)
- Comments supported

### 4.4 Profile-Specific Configuration

```yaml
# application.yml — shared config
spring:
  application:
    name: my-service

---
# application-dev.yml
spring:
  config:
    activate:
      on-profile: dev
  datasource:
    url: jdbc:h2:mem:devdb

---
# application-prod.yml  
spring:
  config:
    activate:
      on-profile: prod
  datasource:
    url: jdbc:postgresql://prod-server/mydb
```

**Activating Profiles:**
```bash
# Command line
java -jar app.jar --spring.profiles.active=prod

# Environment variable
export SPRING_PROFILES_ACTIVE=prod

# In application.properties
spring.profiles.active=dev

# Programmatically
SpringApplication app = new SpringApplication(MyApp.class);
app.setAdditionalProfiles("dev");
app.run(args);
```

### 4.5 @ConfigurationProperties — Type-Safe Configuration

```java
// 1. Define properties class
@ConfigurationProperties(prefix = "app.mail")
@Validated  // Enable validation
public class MailProperties {

    @NotNull
    private String host;
    
    @Min(1) @Max(65535)
    private int port = 25;
    
    private String username;
    
    private String password;
    
    private Duration timeout = Duration.ofSeconds(30);
    
    private final List<String> recipients = new ArrayList<>();
    
    private final Map<String, String> headers = new HashMap<>();
    
    private final Security security = new Security();
    
    // Getters and setters...
    
    public static class Security {
        private boolean enabled = false;
        private String protocol = "TLS";
        // Getters and setters...
    }
}

// 2. Enable in main class or @Configuration
@SpringBootApplication
@EnableConfigurationProperties(MailProperties.class)
// OR use @ConfigurationPropertiesScan
public class MyApp { }

// 3. Use in your service
@Service
public class MailService {
    private final MailProperties props;
    
    public MailService(MailProperties props) {
        this.props = props;
    }
    
    public void send() {
        // props.getHost(), props.getSecurity().isEnabled(), etc.
    }
}
```

**Corresponding YAML:**
```yaml
app:
  mail:
    host: smtp.example.com
    port: 587
    username: user
    password: secret
    timeout: 60s
    recipients:
      - admin@example.com
      - ops@example.com
    headers:
      X-Priority: "1"
      X-Mailer: "MyApp"
    security:
      enabled: true
      protocol: TLSv1.3
```

### 4.6 Property Binding Rules

| Property Format | Java Field Mapping |
|---|---|
| `app.myProp` | `myProp` |
| `app.my-prop` | `myProp` (kebab-case) |
| `app.my_prop` | `myProp` (underscore) |
| `APP_MYPROP` | `myProp` (env var) |
| `app.myProp[0]` | `myProp` (List) |
| `app.myProp.key` | `myProp` (Map) |

**Relaxed binding:** `myProp`, `my-prop`, `my_prop`, `MY_PROP` all map to `myProp`.

### 4.7 @Value vs @ConfigurationProperties

| Feature | @Value | @ConfigurationProperties |
|---|---|---|
| **Type** | Simple injection | Structured binding |
| **Validation** | Limited | Full JSR-303 |
| **IDE Support** | Basic | Metadata generation |
| **Nested Objects** | No | Yes |
| **Collections** | Comma-separated | Native lists/maps |
| **Use Case** | Single values | Groups of config |

```java
// @Value — for simple values
@Value("${server.port:8080}")
private int port;

// @ConfigurationProperties — for related config groups
@ConfigurationProperties(prefix = "server")
public class ServerProperties {
    private int port;
    private Ssl ssl;
    // ...
}
```

### 4.8 Configuration Metadata (IDE Hints)

Generate metadata for IDE auto-completion:

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-configuration-processor</artifactId>
    <optional>true</optional>
</dependency>
```

Creates `META-INF/spring-configuration-metadata.json` with property hints.

### 4.9 Secrets Management

**DON'T commit secrets to version control!**

Options for secrets:

| Method | Use Case |
|---|---|
| **Environment Variables** | Simple, cloud-native |
| **Spring Cloud Config** | Centralized, encrypted |
| **HashiCorp Vault** | Enterprise secret management |
| **AWS Secrets Manager** | AWS deployments |
| **Kubernetes Secrets** | K8s deployments |

```yaml
# Reference environment variable
spring:
  datasource:
    password: ${DB_PASSWORD}

# With default
spring:
  datasource:
    password: ${DB_PASSWORD:default_for_dev}
```

### 4.10 Config Import (Spring Boot 2.4+)

```yaml
# Import additional config files
spring:
  config:
    import:
      - optional:file:./config/extra.yml
      - configserver:http://config-server:8888
      - vault://
```

---

## 5. Spring Boot Web & REST

### 5.0 Web Application Architecture Theory

**Understanding the Request-Response Lifecycle:**

Every HTTP request in a Spring Boot web application follows a precise path through multiple layers:

```
┌─────────────────────────────────────────────────────────────────────────┐
│                   HTTP REQUEST LIFECYCLE                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Client                    Server                                       │
│    │                                                                    │
│    │  GET /api/users/123                                               │
│    ├──────────────────────────────────────────────────────────▶│                    │
│    │                         │  1. Servlet Container (Tomcat)          │
│    │                         ▼  2. Filter Chain                         │
│    │                    ┌─────┴─────┐                                    │
│    │                    │  Filters  │ Auth, CORS, Logging              │
│    │                    └─────┬─────┘                                    │
│    │                          ▼  3. DispatcherServlet                    │
│    │                    ┌─────┴─────┐                                    │
│    │                    │ Dispatch │ HandlerMapping lookup             │
│    │                    └─────┬─────┘                                    │
│    │                          ▼  4. HandlerInterceptors                  │
│    │                    ┌─────┴─────┐                                    │
│    │                    │Intercept│ preHandle, postHandle             │
│    │                    └─────┬─────┘                                    │
│    │                          ▼  5. Controller Method                    │
│    │                    ┌─────┴─────┐                                    │
│    │                    │@GetMap │ Your business logic               │
│    │                    └─────┬─────┘                                    │
│    │                          ▼  6. Response Processing                  │
│    │                    ┌─────┴─────┐                                    │
│    │                    │HttpMsg │ Object → JSON (Jackson)           │
│    │                    └─────┬─────┘                                    │
│    │  200 OK {"id":123}       │                                        │
│    ◀─────────────────────────────────┴──────────────────────────┘                    │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.0.1 Embedded Server vs External Server

**Traditional Deployment Model (Pre-Spring Boot):**
```
┌─────────────────┐        ┌─────────────────┐        ┌─────────────────┐
│   Build WAR    │ deploy │ Install & Config │ deploy │    Production   │
│    artifact    │───────▶│  Tomcat/JBoss    │───────▶│     Server      │
└─────────────────┘        └─────────────────┘        └─────────────────┘
                         Separate install          Version conflicts
                         Admin required            Complex upgrades
```

**Spring Boot Embedded Model:**
```
┌─────────────────┐        ┌─────────────────┐
│  Build FAT JAR │        │    Production   │
│  (Server +App) │───────▶│     Server      │
└─────────────────┘        └─────────────────┘
  java -jar app.jar       No separate install
  Single artifact         Same everywhere
```

| Aspect | External Server | Embedded Server |
|---|---|---|
| **Deployment** | WAR to server | `java -jar` |
| **Configuration** | Server admin | Application controls |
| **Scaling** | Complex | Docker/K8s friendly |
| **Dev-Prod Parity** | Hard to achieve | Same artifact everywhere |
| **Startup Time** | Slow (server init) | Fast (~2 seconds) |
| **Containerization** | Difficult | Natural fit |

### 5.0.2 REST Architectural Constraints (Theory)

**REST (Representational State Transfer)** is an architectural style, not a protocol. Roy Fielding defined 6 constraints:

| Constraint | Meaning | Spring Boot Implementation |
|---|---|---|
| **Client-Server** | Separation of concerns | Controllers = API boundary |
| **Stateless** | No session state on server | Token-based auth (JWT) |
| **Cacheable** | Responses must be cacheable | HTTP cache headers |
| **Uniform Interface** | Consistent resource URLs | `@RequestMapping` patterns |
| **Layered System** | Client doesn't know intermediaries | Works behind proxies/LBs |
| **Code on Demand** | Optional executable code | HATEOAS links |

**Richardson Maturity Model — REST Maturity Levels:**

```
Level 3: Hypermedia Controls (HATEOAS)
        ↑  Links guide client through API
        │  GET /orders/123 → {"id":123, "_links":{"pay":"/orders/123/pay"}}
        │
Level 2: HTTP Verbs  ←─── Most APIs stop here
        ↑  GET, POST, PUT, DELETE properly used
        │  GET /orders (read), POST /orders (create)
        │
Level 1: Resources
        ↑  Individual URIs per resource
        │  /orders/123, /users/456
        │
Level 0: The Swamp of POX
        Single endpoint, RPC-style
        POST /api {"action":"getOrder", "id":123}
```

### 5.1 Embedded Server Architecture

```
┌────────────────────────────────────────────────────────────┐
│                    Spring Boot Web                         │
├────────────────────────────────────────────────────────────┤
│  ┌──────────────────────────────────────────────────────┐  │
│  │            Embedded Servlet Container                │  │
│  │         (Tomcat / Jetty / Undertow)                  │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │              Servlet / Filter                  │  │  │
│  │  │  ┌──────────────────────────────────────────┐  │  │  │
│  │  │  │           DispatcherServlet              │  │  │  │
│  │  │  │  ┌────────────────────────────────────┐  │  │  │  │
│  │  │  │  │     Your @RestController           │  │  │  │  │
│  │  │  │  │       @Service, @Repository        │  │  │  │  │
│  │  │  │  └────────────────────────────────────┘  │  │  │  │
│  │  │  └──────────────────────────────────────────┘  │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘
```

### 5.2 Web Server Configuration

```yaml
server:
  port: 8080
  servlet:
    context-path: /api
  
  # Tomcat specific
  tomcat:
    max-threads: 200
    accept-count: 100
    max-connections: 10000
    connection-timeout: 20000
    
  # SSL/TLS
  ssl:
    enabled: true
    key-store: classpath:keystore.p12
    key-store-password: secret
    key-store-type: PKCS12
    
  # Compression
  compression:
    enabled: true
    mime-types: application/json,application/xml,text/html
    min-response-size: 1024
    
  # HTTP/2
  http2:
    enabled: true
```

### 5.3 Building REST APIs

```java
@RestController
@RequestMapping("/api/v1/users")
@Validated
@Slf4j
public class UserController {

    private final UserService userService;

    public UserController(UserService userService) {
        this.userService = userService;
    }

    @GetMapping
    public ResponseEntity<Page<UserDto>> getAllUsers(
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size,
            @RequestParam(defaultValue = "createdAt") String sortBy) {
        
        Pageable pageable = PageRequest.of(page, size, Sort.by(sortBy).descending());
        return ResponseEntity.ok(userService.findAll(pageable));
    }

    @GetMapping("/{id}")
    public ResponseEntity<UserDto> getUser(@PathVariable Long id) {
        return userService.findById(id)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    @PostMapping
    public ResponseEntity<UserDto> createUser(
            @Valid @RequestBody CreateUserRequest request,
            UriComponentsBuilder uriBuilder) {
        
        UserDto created = userService.create(request);
        URI location = uriBuilder.path("/api/v1/users/{id}")
            .buildAndExpand(created.getId())
            .toUri();
        return ResponseEntity.created(location).body(created);
    }

    @PutMapping("/{id}")
    public ResponseEntity<UserDto> updateUser(
            @PathVariable Long id,
            @Valid @RequestBody UpdateUserRequest request) {
        
        return userService.update(id, request)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void deleteUser(@PathVariable Long id) {
        userService.delete(id);
    }
    
    @PatchMapping("/{id}")
    public ResponseEntity<UserDto> partialUpdate(
            @PathVariable Long id,
            @RequestBody Map<String, Object> updates) {
        
        return userService.partialUpdate(id, updates)
            .map(ResponseEntity::ok)
            .orElse(ResponseEntity.notFound().build());
    }
}
```

### 5.4 Request/Response DTOs with Validation

```java
// Request DTO
public record CreateUserRequest(
    @NotBlank(message = "Username is required")
    @Size(min = 3, max = 50)
    String username,
    
    @NotBlank
    @Email(message = "Invalid email format")
    String email,
    
    @NotBlank
    @Pattern(regexp = "^(?=.*[A-Z])(?=.*[a-z])(?=.*\\d).{8,}$",
             message = "Password must be 8+ chars with uppercase, lowercase, and digit")
    String password,
    
    @Past
    LocalDate birthDate,
    
    @Valid  // Validate nested object
    AddressRequest address
) {}

public record AddressRequest(
    @NotBlank String street,
    @NotBlank String city,
    @NotBlank @Size(min = 2, max = 2) String country
) {}

// Response DTO
public record UserDto(
    Long id,
    String username,
    String email,
    LocalDate birthDate,
    Instant createdAt,
    Instant updatedAt
) {}
```

### 5.5 Global Exception Handling

```java
@RestControllerAdvice
@Slf4j
public class GlobalExceptionHandler {

    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ErrorResponse> handleValidation(MethodArgumentNotValidException ex) {
        List<FieldError> errors = ex.getBindingResult().getFieldErrors().stream()
            .map(fe -> new FieldError(fe.getField(), fe.getDefaultMessage()))
            .toList();
        
        return ResponseEntity.badRequest()
            .body(new ErrorResponse("VALIDATION_FAILED", "Request validation failed", errors));
    }

    @ExceptionHandler(ResourceNotFoundException.class)
    public ResponseEntity<ErrorResponse> handleNotFound(ResourceNotFoundException ex) {
        return ResponseEntity.status(HttpStatus.NOT_FOUND)
            .body(new ErrorResponse("NOT_FOUND", ex.getMessage(), List.of()));
    }

    @ExceptionHandler(DataIntegrityViolationException.class)
    public ResponseEntity<ErrorResponse> handleDataIntegrity(DataIntegrityViolationException ex) {
        return ResponseEntity.status(HttpStatus.CONFLICT)
            .body(new ErrorResponse("CONFLICT", "Data integrity violation", List.of()));
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorResponse> handleGeneral(Exception ex) {
        log.error("Unhandled exception", ex);
        return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
            .body(new ErrorResponse("INTERNAL_ERROR", "An unexpected error occurred", List.of()));
    }
}

public record ErrorResponse(String code, String message, List<FieldError> fieldErrors) {}
public record FieldError(String field, String message) {}
```

### 5.6 JSON Customization (Jackson)

```yaml
# application.yml
spring:
  jackson:
    serialization:
      write-dates-as-timestamps: false
      indent-output: true
    deserialization:
      fail-on-unknown-properties: false
    default-property-inclusion: non_null
    date-format: yyyy-MM-dd HH:mm:ss
    time-zone: UTC
```

**Programmatic Configuration:**
```java
@Configuration
public class JacksonConfig {

    @Bean
    public ObjectMapper objectMapper() {
        return JsonMapper.builder()
            .addModule(new JavaTimeModule())
            .configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false)
            .configure(DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false)
            .serializationInclusion(JsonInclude.Include.NON_NULL)
            .build();
    }
}
```

### 5.7 CORS Configuration

```java
@Configuration
public class CorsConfig {

    @Bean
    public WebMvcConfigurer corsConfigurer() {
        return new WebMvcConfigurer() {
            @Override
            public void addCorsMappings(CorsRegistry registry) {
                registry.addMapping("/api/**")
                    .allowedOrigins("https://frontend.example.com")
                    .allowedMethods("GET", "POST", "PUT", "DELETE", "PATCH")
                    .allowedHeaders("*")
                    .exposedHeaders("X-Custom-Header")
                    .allowCredentials(true)
                    .maxAge(3600);
            }
        };
    }
}
```

### 5.8 API Versioning Strategies

```java
// 1. URL Path versioning (recommended)
@RequestMapping("/api/v1/users")
@RequestMapping("/api/v2/users")

// 2. Header versioning
@GetMapping(headers = "X-API-VERSION=1")
@GetMapping(headers = "X-API-VERSION=2")

// 3. Media type versioning
@GetMapping(produces = "application/vnd.company.v1+json")
@GetMapping(produces = "application/vnd.company.v2+json")

// 4. Query parameter versioning
@GetMapping(params = "version=1")
@GetMapping(params = "version=2")
```

### 5.9 File Upload/Download

```java
@RestController
@RequestMapping("/api/files")
public class FileController {

    @Value("${file.upload-dir}")
    private String uploadDir;

    @PostMapping("/upload")
    public ResponseEntity<String> uploadFile(@RequestParam("file") MultipartFile file) {
        if (file.isEmpty()) {
            return ResponseEntity.badRequest().body("File is empty");
        }
        
        String fileName = StringUtils.cleanPath(file.getOriginalFilename());
        Path targetPath = Paths.get(uploadDir).resolve(fileName);
        
        try {
            Files.copy(file.getInputStream(), targetPath, StandardCopyOption.REPLACE_EXISTING);
            return ResponseEntity.ok("File uploaded: " + fileName);
        } catch (IOException e) {
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body("Failed to upload");
        }
    }

    @GetMapping("/download/{filename}")
    public ResponseEntity<Resource> downloadFile(@PathVariable String filename) {
        try {
            Path filePath = Paths.get(uploadDir).resolve(filename);
            Resource resource = new UrlResource(filePath.toUri());
            
            if (resource.exists()) {
                return ResponseEntity.ok()
                    .contentType(MediaType.APPLICATION_OCTET_STREAM)
                    .header(HttpHeaders.CONTENT_DISPOSITION,
                            "attachment; filename=\"" + resource.getFilename() + "\"")
                    .body(resource);
            }
            return ResponseEntity.notFound().build();
        } catch (MalformedURLException e) {
            return ResponseEntity.badRequest().build();
        }
    }
}
```

**Configuration:**
```yaml
spring:
  servlet:
    multipart:
      enabled: true
      max-file-size: 10MB
      max-request-size: 10MB
```

### 5.10 REST Client — RestTemplate vs WebClient

**RestTemplate (Blocking):**
```java
@Configuration
public class RestTemplateConfig {
    
    @Bean
    public RestTemplate restTemplate(RestTemplateBuilder builder) {
        return builder
            .setConnectTimeout(Duration.ofSeconds(5))
            .setReadTimeout(Duration.ofSeconds(30))
            .build();
    }
}

@Service
public class ExternalApiService {
    
    private final RestTemplate restTemplate;
    
    public UserDto fetchUser(Long id) {
        return restTemplate.getForObject(
            "https://api.example.com/users/{id}",
            UserDto.class,
            id
        );
    }
}
```

**WebClient (Reactive/Non-blocking) — Preferred in Spring Boot 3+:**
```java
@Configuration
public class WebClientConfig {
    
    @Bean
    public WebClient webClient() {
        return WebClient.builder()
            .baseUrl("https://api.example.com")
            .defaultHeader(HttpHeaders.CONTENT_TYPE, MediaType.APPLICATION_JSON_VALUE)
            .filter(logRequest())
            .build();
    }
    
    private ExchangeFilterFunction logRequest() {
        return (request, next) -> {
            log.debug("Request: {} {}", request.method(), request.url());
            return next.exchange(request);
        };
    }
}

@Service
public class ExternalApiService {
    
    private final WebClient webClient;
    
    public Mono<UserDto> fetchUserReactive(Long id) {
        return webClient.get()
            .uri("/users/{id}", id)
            .retrieve()
            .bodyToMono(UserDto.class);
    }
    
    // Blocking call (when you need sync in non-reactive app)
    public UserDto fetchUserBlocking(Long id) {
        return fetchUserReactive(id).block();
    }
}
```

---

## 6. Data Access with Spring Boot

### 6.1 Spring Data JPA Auto-Configuration

With `spring-boot-starter-data-jpa`, Spring Boot auto-configures:
- DataSource (HikariCP by default)
- EntityManagerFactory
- TransactionManager
- Spring Data JPA repositories

```yaml
spring:
  datasource:
    url: jdbc:postgresql://localhost:5432/mydb
    username: user
    password: secret
    driver-class-name: org.postgresql.Driver
    hikari:
      maximum-pool-size: 10
      minimum-idle: 5
      connection-timeout: 30000
      idle-timeout: 600000
      max-lifetime: 1800000

  jpa:
    database-platform: org.hibernate.dialect.PostgreSQLDialect
    hibernate:
      ddl-auto: validate  # none, validate, update, create, create-drop
    show-sql: false
    properties:
      hibernate:
        format_sql: true
        jdbc:
          batch_size: 50
        order_inserts: true
        order_updates: true
    open-in-view: false  # Disable OSIV anti-pattern
```

### 6.2 Entity Modeling

```java
@Entity
@Table(name = "users", indexes = {
    @Index(name = "idx_user_email", columnList = "email", unique = true)
})
@EntityListeners(AuditingEntityListener.class)
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, length = 100)
    private String username;

    @Column(nullable = false, unique = true)
    private String email;

    @Column(nullable = false)
    private String passwordHash;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private UserStatus status = UserStatus.PENDING;

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private Instant createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private Instant updatedAt;

    @Version
    private Long version;  // Optimistic locking

    @OneToMany(mappedBy = "user", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<Order> orders = new ArrayList<>();

    @ManyToMany
    @JoinTable(
        name = "user_roles",
        joinColumns = @JoinColumn(name = "user_id"),
        inverseJoinColumns = @JoinColumn(name = "role_id")
    )
    private Set<Role> roles = new HashSet<>();

    // Constructors, getters, setters, equals/hashCode...
}

public enum UserStatus {
    PENDING, ACTIVE, SUSPENDED, DELETED
}
```

### 6.3 Spring Data JPA Repositories

```java
public interface UserRepository extends JpaRepository<User, Long> {

    // Derived query methods
    Optional<User> findByEmail(String email);
    
    List<User> findByStatusAndCreatedAtAfter(UserStatus status, Instant date);
    
    boolean existsByEmail(String email);
    
    long countByStatus(UserStatus status);
    
    // @Query with JPQL
    @Query("SELECT u FROM User u WHERE u.status = :status ORDER BY u.createdAt DESC")
    List<User> findActiveUsersSorted(@Param("status") UserStatus status);
    
    // @Query with native SQL
    @Query(value = "SELECT * FROM users WHERE email LIKE %:domain", nativeQuery = true)
    List<User> findByEmailDomain(@Param("domain") String domain);
    
    // Pagination
    Page<User> findByStatus(UserStatus status, Pageable pageable);
    
    // Projections
    @Query("SELECT u.id as id, u.username as username FROM User u WHERE u.status = :status")
    List<UserSummary> findSummariesByStatus(@Param("status") UserStatus status);
    
    // Modifying queries
    @Modifying
    @Query("UPDATE User u SET u.status = :status WHERE u.id = :id")
    int updateStatus(@Param("id") Long id, @Param("status") UserStatus status);
    
    // Entity Graph (N+1 solution)
    @EntityGraph(attributePaths = {"roles", "orders"})
    Optional<User> findWithRolesAndOrdersById(Long id);
}

// Projection interface
public interface UserSummary {
    Long getId();
    String getUsername();
}
```

### 6.4 Query Specification for Dynamic Queries

```java
public class UserSpecifications {

    public static Specification<User> hasStatus(UserStatus status) {
        return (root, query, cb) -> 
            status == null ? null : cb.equal(root.get("status"), status);
    }

    public static Specification<User> emailContains(String email) {
        return (root, query, cb) -> 
            email == null ? null : cb.like(cb.lower(root.get("email")), 
                                           "%" + email.toLowerCase() + "%");
    }

    public static Specification<User> createdAfter(Instant date) {
        return (root, query, cb) -> 
            date == null ? null : cb.greaterThan(root.get("createdAt"), date);
    }
}

// Repository extends JpaSpecificationExecutor
public interface UserRepository extends JpaRepository<User, Long>, 
                                        JpaSpecificationExecutor<User> { }

// Usage in service
@Service
public class UserService {
    
    public Page<User> search(UserStatus status, String email, Instant since, Pageable pageable) {
        Specification<User> spec = Specification
            .where(UserSpecifications.hasStatus(status))
            .and(UserSpecifications.emailContains(email))
            .and(UserSpecifications.createdAfter(since));
        
        return userRepository.findAll(spec, pageable);
    }
}
```

### 6.5 Transaction Management

```java
@Service
@Transactional(readOnly = true)  // Default for all methods
public class OrderService {

    @Transactional  // Override: writable transaction
    public Order createOrder(CreateOrderRequest request) {
        // Validates, creates entities, saves
        User user = userRepository.findById(request.userId())
            .orElseThrow(() -> new ResourceNotFoundException("User not found"));
        
        Order order = new Order();
        order.setUser(user);
        order.setItems(mapItems(request.items()));
        order.setTotal(calculateTotal(order.getItems()));
        
        return orderRepository.save(order);
    }

    @Transactional(
        propagation = Propagation.REQUIRES_NEW,
        isolation = Isolation.READ_COMMITTED,
        timeout = 30,
        rollbackFor = BusinessException.class
    )
    public void processPayment(Long orderId, PaymentDetails payment) {
        // Separate transaction for payment
    }
    
    // Read-only (class default) - optimized for queries
    public Page<Order> findOrdersByUser(Long userId, Pageable pageable) {
        return orderRepository.findByUserId(userId, pageable);
    }
}
```

### 6.6 Auditing

```java
@Configuration
@EnableJpaAuditing
public class JpaConfig {

    @Bean
    public AuditorAware<String> auditorProvider() {
        return () -> Optional.ofNullable(SecurityContextHolder.getContext())
            .map(SecurityContext::getAuthentication)
            .filter(Authentication::isAuthenticated)
            .map(Authentication::getName);
    }
}

@MappedSuperclass
@EntityListeners(AuditingEntityListener.class)
public abstract class Auditable {

    @CreatedDate
    @Column(nullable = false, updatable = false)
    private Instant createdAt;

    @LastModifiedDate
    @Column(nullable = false)
    private Instant updatedAt;

    @CreatedBy
    @Column(updatable = false)
    private String createdBy;

    @LastModifiedBy
    private String updatedBy;
    
    // Getters...
}

@Entity
public class Order extends Auditable {
    // Inherits audit fields
}
```

### 6.7 Database Migrations with Flyway

```yaml
spring:
  flyway:
    enabled: true
    locations: classpath:db/migration
    baseline-on-migrate: true
```

**Migration Files:**
```
src/main/resources/db/migration/
├── V1__create_users_table.sql
├── V2__create_orders_table.sql
├── V3__add_user_status.sql
└── V4__add_indexes.sql
```

```sql
-- V1__create_users_table.sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    version BIGINT DEFAULT 0
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
```

### 6.8 Multiple DataSources

```java
@Configuration
public class DataSourceConfig {

    @Primary
    @Bean
    @ConfigurationProperties("spring.datasource.primary")
    public DataSource primaryDataSource() {
        return DataSourceBuilder.create().build();
    }

    @Bean
    @ConfigurationProperties("spring.datasource.secondary")
    public DataSource secondaryDataSource() {
        return DataSourceBuilder.create().build();
    }
    
    @Primary
    @Bean
    public LocalContainerEntityManagerFactoryBean primaryEntityManager(
            @Qualifier("primaryDataSource") DataSource dataSource,
            EntityManagerFactoryBuilder builder) {
        return builder
            .dataSource(dataSource)
            .packages("com.example.primary.entity")
            .persistenceUnit("primary")
            .build();
    }
    
    // Similar for secondary...
}
```

---

## 7. Spring Boot Actuator

### 7.0 Observability Theory: The Three Pillars

**What is Observability?**

Observability is the ability to understand the internal state of a system by examining its external outputs. In distributed systems, observability is CRITICAL because:
- Services communicate asynchronously
- Failures cascade unpredictably
- Traditional debugging (attach debugger) is impossible

**The Three Pillars of Observability:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    THREE PILLARS OF OBSERVABILITY                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌───────────────────┐  ┌───────────────────┐  ┌───────────────────┐   │
│  │      LOGS         │  │     METRICS      │  │     TRACES       │   │
│  ├───────────────────┤  ├───────────────────┤  ├───────────────────┤   │
│  │ Discrete events   │  │ Numeric time-    │  │ Request flow     │   │
│  │ with context      │  │ series data      │  │ across services  │   │
│  │                   │  │                   │  │                   │   │
│  │ "Order 123        │  │ request_count=N  │  │ A→B→C→D          │   │
│  │  created for      │  │ latency_p99=50ms │  │ with timing      │   │
│  │  user 456"        │  │ error_rate=0.1%  │  │ at each hop      │   │
│  ├───────────────────┤  ├───────────────────┤  ├───────────────────┤   │
│  │ Tool: ELK,Loki    │  │ Tool: Prometheus │  │ Tool: Jaeger,    │   │
│  │                   │  │       + Grafana  │  │       Zipkin     │   │
│  └───────────────────┘  └───────────────────┘  └───────────────────┘   │
│                                                                         │
│  Use Cases:                                                             │
│  LOGS    → Debug specific errors, audit trails                         │
│  METRICS → Dashboards, alerting, capacity planning                     │
│  TRACES  → Find bottlenecks, understand request flow                   │
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.0.1 Metrics Theory: RED and USE Methods

**RED Method (Request-driven)** — For user-facing services:
- **R**ate: Requests per second
- **E**rrors: Failed requests per second
- **D**uration: Time per request (latency)

**USE Method (Resource-oriented)** — For infrastructure:
- **U**tilization: % time resource is busy
- **S**aturation: Queue length / waiting work
- **E**rrors: Error count

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    WHEN TO USE EACH METHOD                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  RED Method (Services)          USE Method (Resources)                  │
│  ──────────────────────          ──────────────────────                  │
│  • API endpoints                 • CPU, Memory                            │
│  • Microservices                 • Disk I/O                               │
│  • Web applications              • Network bandwidth                      │
│  • User-facing systems           • Connection pools                       │
│                                  • Thread pools                           │
│                                                                         │
│  Spring Boot Actuator provides metrics for BOTH methods!                │
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.0.2 Health Checks: Liveness vs Readiness

**Understanding Kubernetes Probes:**

| Probe | Question | Failure Action |
|---|---|---|
| **Liveness** | "Is the app alive?" | Restart container |
| **Readiness** | "Can it handle traffic?" | Remove from load balancer |
| **Startup** | "Has it finished starting?" | Wait longer |

```
┌─────────────────────────────────────────────────────────────────────────┐
│              PROBE DECISION GUIDE                                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  LIVENESS: Include ONLY internal app health                            │
│  • JVM running?                                                         │
│  • App not deadlocked?                                                  │
│  • Critical threads alive?                                              │
│                                                                         │
│  DO NOT include external dependencies!                                  │
│  ❌ Database down → liveness fails → app restarts → still down        │
│     = Restart loop!                                                     │
│                                                                         │
│  READINESS: Include dependencies needed to serve traffic               │
│  • Database reachable?                                                  │
│  • Cache warmed up?                                                     │
│  • Initial data loaded?                                                 │
│                                                                         │
│  ✅ Database down → readiness fails → no traffic → wait for recovery   │
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.1 What is Actuator?

Actuator provides **production-ready features** for monitoring and managing your application:
- Health checks
- Metrics
- Application info
- Environment inspection
- HTTP request tracing
- Scheduled task info

### 7.2 Enabling Actuator

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-actuator</artifactId>
</dependency>
```

### 7.3 Endpoint Configuration

```yaml
management:
  endpoints:
    web:
      exposure:
        include: health,info,metrics,prometheus,env,beans,mappings
        # include: "*"  # Expose all (NOT recommended for production)
      base-path: /actuator
      
  endpoint:
    health:
      show-details: when_authorized  # never, when_authorized, always
      show-components: when_authorized
    shutdown:
      enabled: true  # POST /actuator/shutdown to stop app
      
  server:
    port: 8081  # Separate port for management endpoints
```

### 7.4 Available Endpoints

| Endpoint | Description |
|---|---|
| `/actuator/health` | Application health status |
| `/actuator/info` | Application info (build, git, custom) |
| `/actuator/metrics` | Application metrics |
| `/actuator/metrics/{name}` | Specific metric details |
| `/actuator/prometheus` | Prometheus-format metrics |
| `/actuator/env` | Environment properties |
| `/actuator/beans` | All beans in ApplicationContext |
| `/actuator/mappings` | All @RequestMapping paths |
| `/actuator/conditions` | Auto-configuration conditions |
| `/actuator/configprops` | @ConfigurationProperties |
| `/actuator/loggers` | View and modify log levels |
| `/actuator/threaddump` | Thread dump |
| `/actuator/heapdump` | Heap dump (downloads .hprof) |
| `/actuator/scheduledtasks` | Scheduled tasks |
| `/actuator/caches` | Cache info |
| `/actuator/flyway` | Flyway migrations |

### 7.5 Health Indicators

```yaml
# application.yml
management:
  endpoint:
    health:
      show-details: always
```

**Response:**
```json
{
  "status": "UP",
  "components": {
    "db": {
      "status": "UP",
      "details": {
        "database": "PostgreSQL",
        "validationQuery": "isValid()"
      }
    },
    "diskSpace": {
      "status": "UP",
      "details": {
        "total": 499963174912,
        "free": 250000000000,
        "threshold": 10485760
      }
    },
    "redis": {
      "status": "UP",
      "details": {
        "version": "7.0.0"
      }
    }
  }
}
```

### 7.6 Custom Health Indicator

```java
@Component
public class ExternalServiceHealthIndicator implements HealthIndicator {

    private final ExternalServiceClient client;

    @Override
    public Health health() {
        try {
            boolean isUp = client.ping();
            if (isUp) {
                return Health.up()
                    .withDetail("service", "external-api")
                    .withDetail("response_time_ms", 45)
                    .build();
            }
            return Health.down()
                .withDetail("error", "Service not responding")
                .build();
        } catch (Exception e) {
            return Health.down()
                .withException(e)
                .build();
        }
    }
}
```

### 7.7 Application Info

```yaml
# Build info — auto-generated
info:
  app:
    name: ${spring.application.name}
    version: @project.version@
    build-time: @maven.build.timestamp@
  java:
    version: ${java.version}
  
management:
  info:
    env:
      enabled: true
    git:
      enabled: true
      mode: full  # simple or full
    build:
      enabled: true
```

Add Git info plugin:
```xml
<plugin>
    <groupId>pl.project13.maven</groupId>
    <artifactId>git-commit-id-maven-plugin</artifactId>
</plugin>
```

### 7.8 Metrics with Micrometer

Spring Boot uses **Micrometer** as the metrics facade (like SLF4J for logging).

```java
@Service
public class OrderService {

    private final MeterRegistry meterRegistry;
    private final Counter orderCounter;
    private final Timer orderProcessingTimer;

    public OrderService(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
        this.orderCounter = Counter.builder("orders.created")
            .description("Number of orders created")
            .tag("type", "total")
            .register(meterRegistry);
        this.orderProcessingTimer = Timer.builder("orders.processing.time")
            .description("Time to process an order")
            .register(meterRegistry);
    }

    public Order createOrder(CreateOrderRequest request) {
        return orderProcessingTimer.record(() -> {
            Order order = processOrder(request);
            orderCounter.increment();
            
            // Gauge for current pending orders
            meterRegistry.gauge("orders.pending", 
                orderRepository.countByStatus(OrderStatus.PENDING));
            
            return order;
        });
    }
}
```

### 7.9 Prometheus Integration

```xml
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-registry-prometheus</artifactId>
</dependency>
```

```yaml
management:
  endpoints:
    web:
      exposure:
        include: prometheus,health,info,metrics
  prometheus:
    metrics:
      export:
        enabled: true
```

Access metrics at: `GET /actuator/prometheus`

### 7.10 Security for Actuator Endpoints

```java
@Configuration
@EnableWebSecurity
public class ActuatorSecurityConfig {

    @Bean
    public SecurityFilterChain actuatorSecurityFilterChain(HttpSecurity http) throws Exception {
        return http
            .securityMatcher(EndpointRequest.toAnyEndpoint())
            .authorizeHttpRequests(auth -> auth
                .requestMatchers(EndpointRequest.to("health", "info")).permitAll()
                .requestMatchers(EndpointRequest.toAnyEndpoint()).hasRole("ACTUATOR")
            )
            .httpBasic(Customizer.withDefaults())
            .build();
    }
}
```

---

## 8. Testing in Spring Boot

### 8.0 Testing Theory: The Testing Pyramid

**Why Test?**

Tests are an investment with compounding returns:
- Catch bugs before production
- Enable fearless refactoring
- Document expected behavior
- Speed up development (catch issues early)

**The Testing Pyramid:**

```
                              ▲
                             ╱ ╲
                            ╱   ╲
                           ╱ E2E ╲          Few, expensive, slow
                          ╱───────╲         (Selenium, Cypress)
                         ╱         ╲
                        ╱Integration╲        More, moderate cost
                       ╱─────────────╲       (@SpringBootTest)
                      ╱               ╲
                     ╱    Unit Tests   ╲      Many, cheap, fast
                    ╱───────────────────╲     (JUnit + Mockito)
                   ╱                     ╲
                  ─────────────────────────
```

| Layer | Count | Speed | Cost | Spring Boot |
|---|---|---|---|---|
| **Unit** | Many (70%) | Fast (<100ms) | Low | JUnit 5, Mockito |
| **Integration** | Some (20%) | Medium (~1s) | Medium | @WebMvcTest, @DataJpaTest |
| **E2E** | Few (10%) | Slow (~10s) | High | @SpringBootTest, Testcontainers |

### 8.0.1 Test Doubles: Mocks, Stubs, Fakes, Spies

**Understanding Test Doubles:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      TEST DOUBLES EXPLAINED                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  DUMMY: Passed but never used                                          │
│  ─────                                                                   │
│  new Service(dummyLogger)  // Logger required but not tested            │
│                                                                         │
│  STUB: Returns predefined values                                       │
│  ────                                                                    │
│  when(userRepo.findById(1L)).thenReturn(Optional.of(testUser));        │
│                                                                         │
│  MOCK: Verifies interactions                                           │
│  ────                                                                    │
│  verify(emailService, times(1)).sendWelcomeEmail(user);                │
│                                                                         │
│  SPY: Real object with selective overrides                             │
│  ───                                                                     │
│  UserService spy = spy(realService);                                    │
│  doReturn(cachedUser).when(spy).expensiveOperation();                  │
│                                                                         │
│  FAKE: Working implementation (simplified)                             │
│  ────                                                                    │
│  class FakeUserRepository implements UserRepository {                   │
│      private Map<Long, User> db = new HashMap<>();                      │
│  }                                                                      │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.0.2 The AAA Pattern (Arrange-Act-Assert)

**Structure every test consistently:**

```java
@Test
void shouldCreateUserWhenValidRequest() {
    // ARRANGE: Set up test data and dependencies
    CreateUserRequest request = new CreateUserRequest("john", "john@test.com", "Pass123!");
    when(userRepository.save(any())).thenAnswer(inv -> {
        User u = inv.getArgument(0);
        u.setId(1L);
        return u;
    });
    
    // ACT: Execute the code under test
    UserDto result = userService.create(request);
    
    // ASSERT: Verify expected outcomes
    assertThat(result.id()).isEqualTo(1L);
    assertThat(result.username()).isEqualTo("john");
    verify(userRepository).save(any(User.class));
}
```

### 8.0.3 Test Isolation and the F.I.R.S.T. Principles

| Principle | Meaning |
|---|---|
| **F**ast | Tests run quickly (milliseconds) |
| **I**ndependent | No test depends on another |
| **R**epeatable | Same result every time |
| **S**elf-validating | Pass/fail without manual inspection |
| **T**imely | Written at the same time as code (or before in TDD) |

### 8.0.4 Why Test Slices?

```
┌─────────────────────────────────────────────────────────────────────────┐
│                 @SpringBootTest vs Test Slices                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  @SpringBootTest                    @WebMvcTest                         │
│  ────────────────                    ────────────                         │
│  Loads EVERYTHING:                  Loads ONLY web layer:               │
│  • All @Controller                  • @Controller under test             │
│  • All @Service                     • MockMvc                            │
│  • All @Repository                  • Web-specific beans                 │
│  • Database connection                                                  │
│  • Message queues                   Services → @MockBean                │
│  • External services                                                    │
│                                                                         │
│  Startup: ~5-10 seconds             Startup: <1 second                  │
│                                                                         │
│  100 tests × 5s = 8+ minutes        100 tests × 0.5s = 50 seconds       │
└─────────────────────────────────────────────────────────────────────────┘
```

**Rule of Thumb:**
- Use `@WebMvcTest` for controller logic
- Use `@DataJpaTest` for repository queries
- Use `@SpringBootTest` ONLY for end-to-end flows

### 8.1 Test Dependencies

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-test</artifactId>
    <scope>test</scope>
</dependency>
<!-- Includes: JUnit 5, Mockito, AssertJ, Hamcrest, JSONassert, JsonPath -->

<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-testcontainers</artifactId>
    <scope>test</scope>
</dependency>
<dependency>
    <groupId>org.testcontainers</groupId>
    <artifactId>postgresql</artifactId>
    <scope>test</scope>
</dependency>
```

### 8.2 Test Slices Overview

| Annotation | What It Loads | Use Case |
|---|---|---|
| `@SpringBootTest` | Full application context | Integration tests |
| `@WebMvcTest` | Web layer only | Controller tests |
| `@DataJpaTest` | JPA components only | Repository tests |
| `@JsonTest` | JSON serialization | DTO/JSON tests |
| `@RestClientTest` | REST client components | External API client tests |
| `@WebFluxTest` | WebFlux components | Reactive controller tests |

### 8.3 @SpringBootTest — Full Integration Test

```java
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ActiveProfiles("test")
class UserControllerIntegrationTest {

    @Autowired
    private TestRestTemplate restTemplate;
    
    @Autowired
    private UserRepository userRepository;
    
    @LocalServerPort
    private int port;

    @BeforeEach
    void setup() {
        userRepository.deleteAll();
    }

    @Test
    void createUser_shouldReturnCreatedUser() {
        CreateUserRequest request = new CreateUserRequest(
            "john", "john@example.com", "Password123", null, null
        );
        
        ResponseEntity<UserDto> response = restTemplate.postForEntity(
            "/api/v1/users", request, UserDto.class
        );
        
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.CREATED);
        assertThat(response.getBody()).isNotNull();
        assertThat(response.getBody().username()).isEqualTo("john");
        assertThat(response.getHeaders().getLocation()).isNotNull();
    }

    @Test
    void getUser_whenNotFound_shouldReturn404() {
        ResponseEntity<ErrorResponse> response = restTemplate.getForEntity(
            "/api/v1/users/999", ErrorResponse.class
        );
        
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NOT_FOUND);
    }
}
```

### 8.4 @WebMvcTest — Controller Layer Test

```java
@WebMvcTest(UserController.class)
class UserControllerTest {

    @Autowired
    private MockMvc mockMvc;
    
    @MockBean
    private UserService userService;

    @Test
    void getAllUsers_shouldReturnPagedResult() throws Exception {
        Page<UserDto> page = new PageImpl<>(
            List.of(new UserDto(1L, "john", "john@example.com", null, null, null)),
            PageRequest.of(0, 20),
            1
        );
        when(userService.findAll(any(Pageable.class))).thenReturn(page);
        
        mockMvc.perform(get("/api/v1/users")
                .param("page", "0")
                .param("size", "20"))
            .andExpect(status().isOk())
            .andExpect(jsonPath("$.content").isArray())
            .andExpect(jsonPath("$.content[0].username").value("john"))
            .andExpect(jsonPath("$.totalElements").value(1));
    }

    @Test
    void createUser_withInvalidEmail_shouldReturn400() throws Exception {
        String invalidRequest = """
            {
                "username": "john",
                "email": "not-an-email",
                "password": "Password123"
            }
            """;
        
        mockMvc.perform(post("/api/v1/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content(invalidRequest))
            .andExpect(status().isBadRequest())
            .andExpect(jsonPath("$.code").value("VALIDATION_FAILED"));
    }

    @Test
    void createUser_shouldValidateAndCreate() throws Exception {
        CreateUserRequest request = new CreateUserRequest(
            "john", "john@example.com", "Password123", null, null
        );
        UserDto response = new UserDto(1L, "john", "john@example.com", null, null, null);
        
        when(userService.create(any())).thenReturn(response);
        
        mockMvc.perform(post("/api/v1/users")
                .contentType(MediaType.APPLICATION_JSON)
                .content(new ObjectMapper().writeValueAsString(request)))
            .andExpect(status().isCreated())
            .andExpect(header().exists("Location"))
            .andExpect(jsonPath("$.id").value(1))
            .andExpect(jsonPath("$.username").value("john"));
    }
}
```

### 8.5 @DataJpaTest — Repository Layer Test

```java
@DataJpaTest
@AutoConfigureTestDatabase(replace = AutoConfigureTestDatabase.Replace.NONE)
@Testcontainers
class UserRepositoryTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15");

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private TestEntityManager entityManager;

    @Test
    void findByEmail_shouldReturnUser() {
        User user = new User();
        user.setUsername("john");
        user.setEmail("john@example.com");
        user.setPasswordHash("hash");
        user.setStatus(UserStatus.ACTIVE);
        entityManager.persistAndFlush(user);
        
        Optional<User> found = userRepository.findByEmail("john@example.com");
        
        assertThat(found).isPresent();
        assertThat(found.get().getUsername()).isEqualTo("john");
    }

    @Test
    void findByStatus_shouldReturnPagedResults() {
        // Setup test data...
        
        Page<User> activeUsers = userRepository.findByStatus(
            UserStatus.ACTIVE, 
            PageRequest.of(0, 10)
        );
        
        assertThat(activeUsers.getContent()).hasSize(expectedCount);
    }
}
```

### 8.6 Testcontainers Integration

```java
@SpringBootTest
@Testcontainers
class OrderServiceIntegrationTest {

    @Container
    @ServiceConnection
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15");

    @Container
    @ServiceConnection
    static GenericContainer<?> redis = new GenericContainer<>("redis:7")
        .withExposedPorts(6379);

    @Autowired
    private OrderService orderService;

    @Test
    void createOrder_shouldPersistAndCache() {
        // Test with real PostgreSQL and Redis containers
    }
}
```

### 8.7 Testing with @MockBean vs @SpyBean

```java
@WebMvcTest(OrderController.class)
class OrderControllerTest {

    @MockBean  // Complete replacement
    private OrderService orderService;
    
    @SpyBean   // Partial mock — calls real methods unless stubbed
    private OrderMapper orderMapper;
    
    @Test
    void test() {
        // orderService is fully mocked
        // orderMapper calls real implementation unless specifically stubbed
        doReturn(new OrderDto()).when(orderMapper).toDto(any());
    }
}
```

### 8.8 Testing Configuration Properties

```java
@SpringBootTest
@TestPropertySource(properties = {
    "app.feature.enabled=true",
    "app.timeout=5s"
})
class FeatureToggleTest {
    
    @Autowired
    private AppProperties properties;
    
    @Test
    void propertiesShouldBeOverridden() {
        assertThat(properties.getFeature().isEnabled()).isTrue();
    }
}
```

### 8.9 @DirtiesContext — When to Use

```java
@SpringBootTest
@DirtiesContext(classMode = DirtiesContext.ClassMode.AFTER_EACH_TEST_METHOD)
class StatefulComponentTest {
    // Context is reset after each test
    // Use sparingly — slow!
}
```

---

## 9. Logging & Observability

### 9.1 Logging Configuration

Spring Boot uses **Logback** by default.

```yaml
# application.yml
logging:
  level:
    root: INFO
    com.example: DEBUG
    org.springframework.web: INFO
    org.hibernate.SQL: DEBUG
    org.hibernate.type.descriptor.sql.BasicBinder: TRACE  # Show parameter values
    
  pattern:
    console: "%d{yyyy-MM-dd HH:mm:ss} [%thread] %-5level %logger{36} - %msg%n"
    file: "%d{yyyy-MM-dd HH:mm:ss} [%thread] %-5level %logger{36} - %msg%n"
    
  file:
    name: logs/application.log
    max-size: 10MB
    max-history: 30
    total-size-cap: 1GB
```

### 9.2 Logback XML Configuration

For advanced configuration, create `logback-spring.xml`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<configuration>
    <include resource="org/springframework/boot/logging/logback/defaults.xml"/>
    
    <springProfile name="dev">
        <appender name="CONSOLE" class="ch.qos.logback.core.ConsoleAppender">
            <encoder>
                <pattern>%clr(%d{HH:mm:ss.SSS}){faint} %clr(%-5level) %clr([%15.15thread]){faint} %clr(%-40.40logger{39}){cyan} %clr(:){faint} %msg%n</pattern>
            </encoder>
        </appender>
        <root level="DEBUG">
            <appender-ref ref="CONSOLE"/>
        </root>
    </springProfile>
    
    <springProfile name="prod">
        <appender name="FILE" class="ch.qos.logback.core.rolling.RollingFileAppender">
            <file>logs/application.log</file>
            <rollingPolicy class="ch.qos.logback.core.rolling.TimeBasedRollingPolicy">
                <fileNamePattern>logs/application.%d{yyyy-MM-dd}.%i.log.gz</fileNamePattern>
                <maxFileSize>100MB</maxFileSize>
                <maxHistory>30</maxHistory>
                <totalSizeCap>3GB</totalSizeCap>
            </rollingPolicy>
            <encoder>
                <pattern>%d{ISO8601} [%thread] %-5level %logger{36} - %msg%n</pattern>
            </encoder>
        </appender>
        
        <appender name="JSON" class="ch.qos.logback.core.ConsoleAppender">
            <encoder class="net.logstash.logback.encoder.LogstashEncoder"/>
        </appender>
        
        <root level="INFO">
            <appender-ref ref="FILE"/>
            <appender-ref ref="JSON"/>
        </root>
    </springProfile>
</configuration>
```

### 9.3 Structured Logging (JSON)

```xml
<dependency>
    <groupId>net.logstash.logback</groupId>
    <artifactId>logstash-logback-encoder</artifactId>
    <version>7.4</version>
</dependency>
```

**Output:**
```json
{
  "@timestamp": "2024-01-15T10:30:45.123Z",
  "level": "INFO",
  "logger_name": "com.example.OrderService",
  "thread_name": "http-nio-8080-exec-1",
  "message": "Order created",
  "orderId": "12345",
  "userId": "67890",
  "traceId": "abc123def456"
}
```

### 9.4 MDC (Mapped Diagnostic Context)

```java
@Component
public class RequestIdFilter implements Filter {

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain) 
            throws IOException, ServletException {
        String requestId = UUID.randomUUID().toString();
        MDC.put("requestId", requestId);
        MDC.put("userId", getUserId(request));
        try {
            chain.doFilter(request, response);
        } finally {
            MDC.clear();
        }
    }
}

// In logback pattern
<pattern>%d [%X{requestId}] [%X{userId}] %-5level %logger{36} - %msg%n</pattern>
```

### 9.5 Distributed Tracing with Micrometer Tracing

```xml
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-tracing-bridge-brave</artifactId>
</dependency>
<dependency>
    <groupId>io.zipkin.reporter2</groupId>
    <artifactId>zipkin-reporter-brave</artifactId>
</dependency>
```

```yaml
management:
  tracing:
    sampling:
      probability: 1.0  # 100% sampling for dev, lower for prod
  zipkin:
    tracing:
      endpoint: http://zipkin:9411/api/v2/spans
```

---

## 10. Security in Spring Boot

### 10.1 Spring Security Auto-Configuration

With `spring-boot-starter-security`, Spring Boot automatically:
- Secures all endpoints with HTTP Basic
- Generates a random password (logged at startup)
- Enables CSRF protection for form-based auth
- Adds security headers

### 10.2 Basic Security Configuration

```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .csrf(csrf -> csrf.disable())  // Disable for APIs
            .sessionManagement(session -> 
                session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/public/**", "/actuator/health").permitAll()
                .requestMatchers("/api/admin/**").hasRole("ADMIN")
                .requestMatchers(HttpMethod.DELETE).hasAuthority("DELETE")
                .anyRequest().authenticated()
            )
            .httpBasic(Customizer.withDefaults())
            .build();
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder(12);
    }

    @Bean
    public UserDetailsService userDetailsService(UserRepository userRepo) {
        return username -> userRepo.findByEmail(username)
            .map(user -> User.builder()
                .username(user.getEmail())
                .password(user.getPasswordHash())
                .authorities(user.getRoles().stream()
                    .map(role -> new SimpleGrantedAuthority("ROLE_" + role.getName()))
                    .toList())
                .build())
            .orElseThrow(() -> new UsernameNotFoundException("User not found"));
    }
}
```

### 10.3 JWT Authentication

```java
@Configuration
@EnableWebSecurity
public class JwtSecurityConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http, JwtAuthFilter jwtAuthFilter) 
            throws Exception {
        return http
            .csrf(csrf -> csrf.disable())
            .sessionManagement(session -> 
                session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/auth/**").permitAll()
                .anyRequest().authenticated()
            )
            .addFilterBefore(jwtAuthFilter, UsernamePasswordAuthenticationFilter.class)
            .build();
    }
}

@Component
public class JwtAuthFilter extends OncePerRequestFilter {

    private final JwtService jwtService;
    private final UserDetailsService userDetailsService;

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) 
            throws ServletException, IOException {
        
        String authHeader = request.getHeader("Authorization");
        
        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            filterChain.doFilter(request, response);
            return;
        }
        
        String jwt = authHeader.substring(7);
        String username = jwtService.extractUsername(jwt);
        
        if (username != null && SecurityContextHolder.getContext().getAuthentication() == null) {
            UserDetails userDetails = userDetailsService.loadUserByUsername(username);
            
            if (jwtService.isTokenValid(jwt, userDetails)) {
                UsernamePasswordAuthenticationToken authToken = 
                    new UsernamePasswordAuthenticationToken(
                        userDetails, null, userDetails.getAuthorities()
                    );
                authToken.setDetails(new WebAuthenticationDetailsSource().buildDetails(request));
                SecurityContextHolder.getContext().setAuthentication(authToken);
            }
        }
        
        filterChain.doFilter(request, response);
    }
}

@Service
public class JwtService {

    @Value("${jwt.secret}")
    private String secretKey;
    
    @Value("${jwt.expiration}")
    private Duration expiration;

    public String generateToken(UserDetails userDetails) {
        return Jwts.builder()
            .setSubject(userDetails.getUsername())
            .setIssuedAt(new Date())
            .setExpiration(Date.from(Instant.now().plus(expiration)))
            .signWith(getSigningKey(), SignatureAlgorithm.HS256)
            .compact();
    }

    public String extractUsername(String token) {
        return extractClaim(token, Claims::getSubject);
    }

    public boolean isTokenValid(String token, UserDetails userDetails) {
        String username = extractUsername(token);
        return username.equals(userDetails.getUsername()) && !isTokenExpired(token);
    }

    private boolean isTokenExpired(String token) {
        return extractClaim(token, Claims::getExpiration).before(new Date());
    }

    private <T> T extractClaim(String token, Function<Claims, T> claimsResolver) {
        Claims claims = Jwts.parserBuilder()
            .setSigningKey(getSigningKey())
            .build()
            .parseClaimsJws(token)
            .getBody();
        return claimsResolver.apply(claims);
    }

    private Key getSigningKey() {
        return Keys.hmacShaKeyFor(Decoders.BASE64.decode(secretKey));
    }
}
```

### 10.4 OAuth2 Resource Server

```yaml
spring:
  security:
    oauth2:
      resourceserver:
        jwt:
          issuer-uri: https://auth.example.com/
          # OR
          jwk-set-uri: https://auth.example.com/.well-known/jwks.json
```

```java
@Configuration
@EnableWebSecurity
public class OAuth2ResourceServerConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/public/**").permitAll()
                .anyRequest().authenticated()
            )
            .oauth2ResourceServer(oauth2 -> oauth2
                .jwt(jwt -> jwt
                    .jwtAuthenticationConverter(jwtAuthenticationConverter())
                )
            )
            .build();
    }

    private JwtAuthenticationConverter jwtAuthenticationConverter() {
        JwtGrantedAuthoritiesConverter authoritiesConverter = new JwtGrantedAuthoritiesConverter();
        authoritiesConverter.setAuthoritiesClaimName("roles");
        authoritiesConverter.setAuthorityPrefix("ROLE_");

        JwtAuthenticationConverter converter = new JwtAuthenticationConverter();
        converter.setJwtGrantedAuthoritiesConverter(authoritiesConverter);
        return converter;
    }
}
```

---

## 11. Production Readiness

### 11.0 Production Theory: The Fallacies of Distributed Computing

**Understanding why production is hard:**

Peter Deutsch identified 8 fallacies that developers often believe about networks:

```
┌─────────────────────────────────────────────────────────────────────────┐
│          THE 8 FALLACIES OF DISTRIBUTED COMPUTING                      │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. The network is reliable        → Use retries, circuit breakers     │
│  2. Latency is zero                → Set timeouts, async operations    │
│  3. Bandwidth is infinite          → Paginate, compress, cache         │
│  4. The network is secure          → TLS everywhere, zero trust        │
│  5. Topology doesn't change        → Service discovery, DNS            │
│  6. There is one administrator     → Automation, GitOps                │
│  7. Transport cost is zero         → Minimize serialization            │
│  8. The network is homogeneous     → Standard protocols (HTTP, gRPC)   │
│                                                                         │
│  Spring Boot helps address these via:                                   │
│  • Resilience4j (circuit breakers, retries)                             │
│  • WebClient with timeouts                                              │
│  • Spring Security (TLS, auth)                                          │
│  • Spring Cloud (service discovery)                                     │
└─────────────────────────────────────────────────────────────────────────┘
```

### 11.0.1 CAP Theorem — The Fundamental Trade-off

**Every distributed system must choose two of three:**

```
                         C (Consistency)
                              ▲
                             ╱ ╲
                            ╱   ╲
                           ╱     ╲
                          ╱       ╲
                         ╱ CP      ╲
                        ╱ Systems   ╲
                       ╱             ╲
                      ╱───────────────╲
                     ╱                 ╲
                    ╱    CA Systems    ╲
                   ╱  (single node only) ╲
                  ╱                       ╲
                 ╱                         ╲
               A ◄───────────────────────────▸ P
         (Availability)    AP Systems   (Partition Tolerance)

CP: Strong consistency, sacrifice availability during partitions
    Examples: Traditional RDBMS with sync replication, ZooKeeper
    
AP: Always available, may return stale data during partitions
    Examples: Cassandra, DynamoDB, DNS
    
CA: Not possible in distributed systems (partitions WILL happen)
```

**In Practice — PACELC:**

> "If there's a Partition, choose A or C; Else, choose L (Latency) or C (Consistency)"

| System | During Partition | Normal Operation |
|---|---|---|
| RDBMS (single) | N/A | Consistent |
| Cassandra | Available | Tunable (eventual) |
| MongoDB | Consistent | Consistent |
| Redis Cluster | Available | Consistent (leader) |

### 11.0.2 Resilience Patterns

**Circuit Breaker Pattern:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CIRCUIT BREAKER STATES                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│         ┌────────────┐                                                  │
│         │   CLOSED   │ ◀── Normal operation                            │
│         │ (Healthy)  │     Requests pass through                       │
│         └─────┬──────┘                                                  │
│               │ Failures exceed threshold                              │
│               ▼                                                         │
│         ┌────────────┐                                                  │
│         │    OPEN    │ ◀── Circuit tripped                              │
│         │  (Failing) │     Requests fail fast                          │
│         └─────┬──────┘     (no calls to failing service)               │
│               │ Wait timeout                                           │
│               ▼                                                         │
│         ┌────────────┐                                                  │
│         │ HALF-OPEN  │ ◀── Testing recovery                            │
│         │ (Testing)  │     Allow limited requests                      │
│         └─────┬──────┘                                                  │
│               │                                                        │
│               ├─── Success → Back to CLOSED                            │
│               └─── Failure → Back to OPEN                              │
└─────────────────────────────────────────────────────────────────────────┘
```

**Other Resilience Patterns:**

| Pattern | Purpose | Use Case |
|---|---|---|
| **Retry** | Recover from transient failures | Network glitches |
| **Timeout** | Prevent indefinite waiting | Slow dependencies |
| **Bulkhead** | Isolate failures | Prevent cascade |
| **Rate Limiter** | Control traffic | Protect resources |
| **Fallback** | Graceful degradation | Feature unavailable |

### 11.1 Graceful Shutdown

```yaml
server:
  shutdown: graceful

spring:
  lifecycle:
    timeout-per-shutdown-phase: 30s
```

This ensures:
- New requests are rejected (503)
- In-flight requests complete
- Resources are cleaned up

### 11.2 Liveness vs Readiness Probes

```yaml
management:
  endpoint:
    health:
      probes:
        enabled: true
      group:
        liveness:
          include: livenessState
        readiness:
          include: readinessState,db,redis
```

| Probe | Purpose | When to Fail |
|---|---|---|
| **Liveness** | "Is the app alive?" | App is stuck, needs restart |
| **Readiness** | "Is the app ready for traffic?" | Dependencies unavailable |

```java
@Component
public class CustomReadinessIndicator implements HealthIndicator {
    
    @Override
    public Health health() {
        if (!warmupComplete) {
            return Health.outOfService()
                .withDetail("reason", "Cache warming in progress")
                .build();
        }
        return Health.up().build();
    }
}
```

### 11.3 Connection Pool Tuning

```yaml
spring:
  datasource:
    hikari:
      maximum-pool-size: 10           # Max connections
      minimum-idle: 5                  # Min idle connections
      connection-timeout: 30000        # Wait for connection (ms)
      idle-timeout: 600000             # Max idle time (ms)
      max-lifetime: 1800000            # Max connection lifetime (ms)
      leak-detection-threshold: 60000  # Leak detection (ms)
      pool-name: MyHikariPool
```

**Sizing formula:**
```
pool_size = Tn × (Cm - 1) + 1
Where:
  Tn = Number of threads
  Cm = Number of simultaneous connections per thread
```

**Conservative starting point:** `pool_size = (core_count * 2) + number_of_disks`

### 11.4 JVM Tuning

```bash
java -jar app.jar \
  -Xms512m \
  -Xmx2g \
  -XX:+UseG1GC \
  -XX:MaxGCPauseMillis=200 \
  -XX:+HeapDumpOnOutOfMemoryError \
  -XX:HeapDumpPath=/logs/heapdump.hprof \
  -Djava.security.egd=file:/dev/./urandom
```

### 11.5 Configuration for Different Environments

```yaml
# application.yml (shared)
spring:
  application:
    name: my-service
    
---
# application-dev.yml
spring:
  config:
    activate:
      on-profile: dev
  datasource:
    url: jdbc:h2:mem:devdb
  jpa:
    show-sql: true
logging:
  level:
    com.example: DEBUG
    
---
# application-prod.yml
spring:
  config:
    activate:
      on-profile: prod
  datasource:
    url: jdbc:postgresql://${DB_HOST}/${DB_NAME}
    username: ${DB_USER}
    password: ${DB_PASSWORD}
  jpa:
    show-sql: false
logging:
  level:
    root: WARN
    com.example: INFO
```

---

## 12. Docker & Kubernetes Deployment

### 12.0 Container Theory: Why Containers?

**The "It Works on My Machine" Problem:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│             THE DEPLOYMENT PROBLEM (Before Containers)                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Developer Laptop        Staging Server        Production              │
│  ────────────────        ──────────────        ──────────              │
│  Java 21                 Java 17               Java 11                  │
│  macOS                   Ubuntu 22             RHEL 8                   │
│  OpenSSL 3.0             OpenSSL 1.1           OpenSSL 1.0              │
│  glibc 2.35              glibc 2.31            glibc 2.28               │
│                                                                         │
│  App works ✅             Subtle bugs ⚠         Crash! ❌                 │
└─────────────────────────────────────────────────────────────────────────┘
```

**Containers Solve This:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                     CONTAINER = PORTABLE UNIT                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Container Image = App + Runtime + Dependencies + Config               │
│                                                                         │
│  ┌───────────────────────────────────────────────────────────────┐   │
│  │  Docker Image: myapp:1.0.0                                       │   │
│  │  ┌─────────────────────────────────────────────────────────┐    │   │
│  │  │  Your Application (app.jar)                                │    │   │
│  │  ├─────────────────────────────────────────────────────────┤    │   │
│  │  │  JDK 21 (exact version)                                    │    │   │
│  │  ├─────────────────────────────────────────────────────────┤    │   │
│  │  │  Linux Libraries (glibc, OpenSSL, etc.)                    │    │   │
│  │  ├─────────────────────────────────────────────────────────┤    │   │
│  │  │  Base OS (Debian Slim, Alpine)                             │    │   │
│  │  └─────────────────────────────────────────────────────────┘    │   │
│  └───────────────────────────────────────────────────────────────┘   │
│                                                                         │
│  RUNS IDENTICALLY ON: Dev laptop, CI server, Staging, Production       │
└─────────────────────────────────────────────────────────────────────────┘
```

### 12.0.1 Containers vs VMs

```
┌─────────────────────────────────┐   ┌───────────────────────────────────┐
│     Virtual Machines             │   │        Containers                 │
├─────────────────────────────────┤   ├───────────────────────────────────┤
│  ┌─────┐ ┌─────┐ ┌─────┐        │   │  ┌───┐ ┌───┐ ┌───┐ ┌───┐ ┌───┐  │
│  │ App │ │ App │ │ App │        │   │  │App│ │App│ │App│ │App│ │App│  │
│  ├─────┤ ├─────┤ ├─────┤        │   │  └───┘ └───┘ └───┘ └───┘ └───┘  │
│  │Bins │ │Bins │ │Bins │        │   │  ─────────────────────────────  │
│  ├─────┤ ├─────┤ ├─────┤        │   │       Container Runtime (Docker)  │
│  │Guest│ │Guest│ │Guest│        │   │  ─────────────────────────────  │
│  │ OS  │ │ OS  │ │ OS  │        │   │         Host Operating System    │
│  └─────┘ └─────┘ └─────┘        │   │  ─────────────────────────────  │
│  ────────────────────────────   │   │             Hardware             │
│         Hypervisor (VMware, etc)  │   └───────────────────────────────────┘
│  ────────────────────────────   │
│       Host Operating System        │   Size:    ~MBs
│  ────────────────────────────   │   Startup: ~seconds
│           Hardware                 │   Density: Dozens per host
└─────────────────────────────────┘

 Size: ~GBs
 Startup: ~minutes
 Density: Few per host
```

| Aspect | VM | Container |
|---|---|---|
| **Isolation** | Hardware-level | Process-level |
| **Size** | GBs | MBs |
| **Startup** | Minutes | Seconds |
| **Overhead** | High | Low |
| **Density** | ~10s per host | ~100s per host |
| **Portability** | Hypervisor-specific | Any container runtime |

### 12.0.2 Kubernetes Concepts for Developers

```
┌─────────────────────────────────────────────────────────────────────────┐
│               KUBERNETES KEY CONCEPTS                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  POD: Smallest deployable unit                                         │
│       One or more containers with shared network/storage               │
│       Usually 1 container per pod for microservices                    │
│                                                                         │
│  DEPLOYMENT: Manages ReplicaSets of Pods                               │
│              Handles rollouts, scaling, rollbacks                      │
│              "I want 3 replicas of my-app:v2"                          │
│                                                                         │
│  SERVICE: Stable network endpoint for Pods                             │
│           Pods are ephemeral; Services provide DNS name                │
│           Load balances across pod replicas                            │
│                                                                         │
│  CONFIGMAP: Non-sensitive configuration data                           │
│             Injected as environment variables or files                 │
│                                                                         │
│  SECRET: Sensitive data (passwords, tokens)                            │
│          Base64 encoded, can be encrypted at rest                      │
│                                                                         │
│  INGRESS: HTTP routing rules to Services                               │
│           SSL termination, path-based routing                          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 12.1 Dockerfile (Multi-Stage Build)

```dockerfile
# Build stage
FROM eclipse-temurin:21-jdk as builder
WORKDIR /app
COPY mvnw pom.xml ./
COPY .mvn .mvn
RUN ./mvnw dependency:go-offline
COPY src src
RUN ./mvnw package -DskipTests

# Extract layers for better caching
FROM eclipse-temurin:21-jdk as extractor
WORKDIR /app
COPY --from=builder /app/target/*.jar app.jar
RUN java -Djarmode=layertools -jar app.jar extract

# Runtime stage
FROM eclipse-temurin:21-jre
WORKDIR /app

# Create non-root user
RUN groupadd -r spring && useradd -r -g spring spring
USER spring:spring

# Copy layers in order of change frequency
COPY --from=extractor /app/dependencies/ ./
COPY --from=extractor /app/spring-boot-loader/ ./
COPY --from=extractor /app/snapshot-dependencies/ ./
COPY --from=extractor /app/application/ ./

EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=60s \
  CMD curl -f http://localhost:8080/actuator/health || exit 1

ENTRYPOINT ["java", "org.springframework.boot.loader.launch.JarLauncher"]
```

### 12.2 Building with Cloud Native Buildpacks

```bash
# Maven
./mvnw spring-boot:build-image -Dspring-boot.build-image.imageName=myapp:latest

# Gradle
./gradlew bootBuildImage --imageName=myapp:latest
```

No Dockerfile needed!

### 12.3 Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
  labels:
    app: my-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
      - name: my-app
        image: myregistry/my-app:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: SPRING_PROFILES_ACTIVE
          value: "prod"
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: password
        resources:
          requests:
            memory: "512Mi"
            cpu: "250m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
        livenessProbe:
          httpGet:
            path: /actuator/health/liveness
            port: 8080
          initialDelaySeconds: 60
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /actuator/health/readiness
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 5
        volumeMounts:
        - name: config
          mountPath: /config
          readOnly: true
      volumes:
      - name: config
        configMap:
          name: my-app-config
---
apiVersion: v1
kind: Service
metadata:
  name: my-app
spec:
  selector:
    app: my-app
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-app-config
data:
  application.yml: |
    server:
      port: 8080
    logging:
      level:
        root: INFO
```

### 12.4 Kubernetes ConfigMaps and Secrets

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-secret
type: Opaque
stringData:
  password: mysecretpassword
---
# In deployment:
env:
- name: SPRING_DATASOURCE_PASSWORD
  valueFrom:
    secretKeyRef:
      name: db-secret
      key: password
```

---

## 13. Performance Tuning

### 13.0 Performance Theory: Amdahl's Law and Scalability

**Amdahl's Law — The Limits of Parallelism:**

```
Speedup = 1 / ((1 - P) + P/N)

Where:
  P = Fraction of code that can be parallelized
  N = Number of processors

┌─────────────────────────────────────────────────────────────────────────┐
│                     AMDAHL'S LAW IN PRACTICE                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  If 95% of your code is parallelizable (P=0.95):                        │
│                                                                         │
│    N=2  CPU:   Speedup = 1 / (0.05 + 0.475) = 1.90x                     │
│    N=4  CPU:   Speedup = 1 / (0.05 + 0.238) = 3.48x                     │
│    N=8  CPU:   Speedup = 1 / (0.05 + 0.119) = 5.93x                     │
│    N=16 CPU:   Speedup = 1 / (0.05 + 0.059) = 9.14x                     │
│    N=∞  CPU:   Speedup = 1 / 0.05 = 20x (MAX)                           │
│                                                                         │
│  Takeaway: The serial portion is the bottleneck!                       │
│  Adding more threads won't help if code isn't parallel.                 │
└─────────────────────────────────────────────────────────────────────────┘
```

### 13.0.1 Little's Law — Understanding Throughput

```
L = λ × W

Where:
  L = Number of requests in system (concurrency)
  λ = Arrival rate (requests/second)
  W = Average time in system (latency)

Rearranged for throughput:
  λ = L / W

┌─────────────────────────────────────────────────────────────────────────┐
│                     LITTLE'S LAW EXAMPLE                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  If you have:                                                           │
│    • 200 threads (max concurrency)                                      │
│    • 100ms average latency                                              │
│                                                                         │
│  Max throughput = 200 / 0.1 = 2000 requests/second                      │
│                                                                         │
│  To increase throughput, you can either:                                │
│    • Increase concurrency (more threads, async)                         │
│    • Reduce latency (optimize code, caching)                            │
└─────────────────────────────────────────────────────────────────────────┘
```

### 13.0.2 Threading Models: Traditional vs Virtual Threads

```
┌─────────────────────────────────────────────────────────────────────────┐
│                 THREADING MODEL COMPARISON                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  PLATFORM THREADS (Traditional)                                         │
│  ─────────────────────────────                                         │
│  • 1 platform thread = 1 OS thread                                      │
│  • Stack size: ~1MB each                                                │
│  • Max practical: ~200-500 per JVM                                      │
│  • Blocking I/O blocks the entire thread                                │
│                                                                         │
│  Problem: 10,000 concurrent connections = 10GB RAM just for stacks!    │
│                                                                         │
│  VIRTUAL THREADS (Java 21+)                                             │
│  ───────────────────────                                                 │
│  • Millions of virtual threads per JVM                                  │
│  • Stack: ~KBs, grows as needed                                         │
│  • Scheduled on carrier threads (platform threads)                      │
│  • Blocking I/O releases carrier thread                                 │
│                                                                         │
│  Benefit: Write blocking code, get async performance!                  │
│                                                                         │
│  Platform Thread:    [Thread A: request handling.........]             │
│                                     ^ waiting for DB                    │
│                                                                         │
│  Virtual Thread:     [VT: start]--[park]--[continue]                   │
│  Carrier Thread:     [VT1][VT2][VT1][VT3][VT1]   (multiplexed)          │
└─────────────────────────────────────────────────────────────────────────┘
```

### 13.0.3 JIT vs AOT Compilation

| Aspect | JIT (Just-In-Time) | AOT (Ahead-of-Time) |
|---|---|---|
| **When** | At runtime | At build time |
| **Startup** | Slow (interpret then compile) | Fast (native binary) |
| **Peak Performance** | Excellent (runtime optimization) | Good (no runtime data) |
| **Memory** | Higher (JVM overhead) | Lower |
| **Use Case** | Long-running services | Serverless, CLI tools |
| **Spring Boot** | Default | GraalVM native image |

### 13.1 Startup Time Optimization

```yaml
spring:
  main:
    lazy-initialization: true   # Defer bean creation
    
  jpa:
    defer-datasource-initialization: true
    open-in-view: false         # Disable OSIV
    
  data:
    jpa:
      repositories:
        bootstrap-mode: lazy    # Lazy repo init
```

### 13.2 Virtual Threads (Java 21+)

```yaml
spring:
  threads:
    virtual:
      enabled: true
```

```java
@Bean
public TomcatProtocolHandlerCustomizer<?> virtualThreadsCustomizer() {
    return protocolHandler -> {
        protocolHandler.setExecutor(Executors.newVirtualThreadPerTaskExecutor());
    };
}
```

### 13.3 Native Image with GraalVM

```xml
<plugin>
    <groupId>org.graalvm.buildtools</groupId>
    <artifactId>native-maven-plugin</artifactId>
</plugin>
```

```bash
./mvnw -Pnative native:compile
```

Benefits:
- Instant startup (~50ms vs ~2s)
- Lower memory footprint
- Smaller container images

Trade-offs:
- Longer build time
- Limited reflection support
- Some libraries need configuration

### 13.4 Caching

```yaml
spring:
  cache:
    type: redis  # or caffeine, ehcache
    redis:
      time-to-live: 600000

  data:
    redis:
      host: localhost
      port: 6379
```

```java
@EnableCaching
@Configuration
public class CacheConfig {

    @Bean
    public CacheManager cacheManager(RedisConnectionFactory connectionFactory) {
        RedisCacheConfiguration config = RedisCacheConfiguration.defaultCacheConfig()
            .entryTtl(Duration.ofMinutes(10))
            .serializeValuesWith(
                RedisSerializationContext.SerializationPair.fromSerializer(
                    new GenericJackson2JsonRedisSerializer()
                )
            );
        
        return RedisCacheManager.builder(connectionFactory)
            .cacheDefaults(config)
            .withCacheConfiguration("users", 
                config.entryTtl(Duration.ofMinutes(5)))
            .build();
    }
}

@Service
public class UserService {

    @Cacheable(value = "users", key = "#id")
    public UserDto findById(Long id) {
        return userRepository.findById(id)
            .map(this::toDto)
            .orElseThrow();
    }

    @CacheEvict(value = "users", key = "#id")
    public void deleteUser(Long id) {
        userRepository.deleteById(id);
    }

    @CachePut(value = "users", key = "#result.id")
    public UserDto updateUser(Long id, UpdateUserRequest request) {
        // Update and return new value — updates cache
    }
}
```

---

## 14. Best Practices & Anti-Patterns

### 14.1 Best Practices

| # | Practice | Description |
|---|---|---|
| 1 | **Use starters wisely** | Don't add starters you don't need |
| 2 | **Externalize all config** | Never hardcode URLs, credentials, feature flags |
| 3 | **Profile-specific config** | Separate dev/staging/prod configurations |
| 4 | **Health endpoints** | Implement meaningful health checks |
| 5 | **Structured logging** | JSON logs with correlation IDs |
| 6 | **Graceful shutdown** | Allow in-flight requests to complete |
| 7 | **Connection pool tuning** | Size based on workload, not defaults |
| 8 | **Test slices** | Use `@WebMvcTest`, `@DataJpaTest` for focused tests |
| 9 | **Testcontainers** | Real databases in integration tests |
| 10 | **Layer Docker images** | Leverage Spring Boot's layer extraction |

### 14.2 Common Anti-Patterns

| Anti-Pattern | Problem | Solution |
|---|---|---|
| **OSIV enabled** | Long-running DB connections | Set `spring.jpa.open-in-view=false` |
| **N+1 queries** | Performance degradation | Use `@EntityGraph`, fetch joins, or projections |
| **Blocking in async** | Thread starvation | Use proper async patterns |
| **Secrets in properties** | Security risk | Use env vars, Vault, or K8s secrets |
| **@SpringBootTest everywhere** | Slow tests | Use test slices |
| **No connection pool tuning** | Connection exhaustion | Size pool appropriately |
| **Missing health checks** | Undetected failures | Add custom health indicators |
| **Including all Actuator endpoints** | Security exposure | Expose only needed endpoints |
| **No graceful shutdown** | Dropped requests on deploy | Enable graceful shutdown |
| **Relying on auto-restart** | Masks real issues | Investigate root causes |

### 14.3 Production Deployment Checklist

```markdown
□ Externalized configuration (no hardcoded values)
□ Secrets managed securely (not in version control)
□ Health endpoints configured
□ Metrics and tracing enabled
□ Logging configured (JSON, appropriate levels)
□ Graceful shutdown enabled
□ Connection pools tuned
□ Resource limits set (JVM heap, container limits)
□ Security hardened (HTTPS, auth, input validation)
□ Error handling comprehensive
□ API versioned
□ Database migrations in place (Flyway/Liquibase)
□ Caching strategy implemented
□ Rate limiting considered
□ Circuit breakers for external calls
□ Backup and recovery plan
□ Monitoring and alerting configured
□ Documentation updated (API docs, runbooks)
```

---

## 15. Interview Questions by Experience Level

### 15.1 Junior (0–2 Years)

**Q1: What is Spring Boot and how is it different from Spring Framework?**
> Spring Boot is an opinionated framework built on Spring that provides auto-configuration, starter dependencies, and embedded servers. Spring Framework requires manual configuration; Spring Boot applies sensible defaults automatically while allowing customization.

**Q2: What does `@SpringBootApplication` do?**
> It's a combination of `@Configuration` (Java config class), `@EnableAutoConfiguration` (trigger auto-config), and `@ComponentScan` (scan for components from current package down).

**Q3: What are Spring Boot starters?**
> Curated sets of dependencies that provide all you need for a specific functionality. e.g., `spring-boot-starter-web` includes Spring MVC, embedded Tomcat, Jackson.

**Q4: How do you configure properties in Spring Boot?**
> Via `application.properties` or `application.yml`. Properties can be overridden by environment variables, command-line arguments, or profile-specific files.

**Q5: What is the purpose of `application.properties`?**
> It's the default configuration file where you define properties like server port, database URL, logging levels. It supports placeholders (`${VAR}`) and defaults (`${VAR:default}`).

**Q6: How do you run a Spring Boot application?**
> Run the `main()` method of the class annotated with `@SpringBootApplication`, or use `./mvnw spring-boot:run`, or `java -jar app.jar` for packaged JAR.

**Q7: What is the default embedded server in Spring Boot?**
> Tomcat. You can switch to Jetty or Undertow by excluding Tomcat and adding the alternative starter.

---

### 15.2 Mid-Level (2–5 Years)

**Q8: Explain auto-configuration. How does Spring Boot decide which beans to create?**
> Auto-configuration uses `@Conditional` annotations to check conditions like "is class X on classpath?", "is bean Y missing?", "is property Z set?". Classes are registered in `META-INF/spring/org.springframework.boot.autoconfigure.AutoConfiguration.imports`.

**Q9: How would you create a custom starter?**
> Create two modules: 1) `*-autoconfigure` with `@AutoConfiguration` class, `@ConfigurationProperties`, and registration in imports file. 2) `*-starter` that depends on autoconfigure module and actual library.

**Q10: Explain `@ConfigurationProperties` vs `@Value`.**
> `@ConfigurationProperties` provides type-safe binding for groups of related properties with validation support. `@Value` is for single values with SpEL support. Use `@ConfigurationProperties` for structured config, `@Value` for simple one-off values.

**Q11: What is Spring Boot Actuator? Name 5 important endpoints.**
> Actuator provides production-ready features. Important endpoints: `/health` (health status), `/metrics` (application metrics), `/info` (app info), `/env` (environment properties), `/loggers` (view/modify log levels).

**Q12: How do you secure Actuator endpoints?**
> Use Spring Security to protect endpoints. Expose only needed endpoints via `management.endpoints.web.exposure.include`. Run management on separate port. Use roles like `ACTUATOR_ADMIN`.

**Q13: Explain the difference between `@MockBean` and `@Mock`.**
> `@Mock` is Mockito annotation for plain unit tests. `@MockBean` is Spring Boot annotation that creates a mock AND replaces the bean in the ApplicationContext, useful in integration tests.

**Q14: What are test slices? Give examples.**
> Test slices load only relevant parts of the context. `@WebMvcTest` loads web layer only, `@DataJpaTest` loads JPA components only. Faster than `@SpringBootTest` and better isolation.

**Q15: How do you handle different configurations for dev/prod?**
> Use profiles: `application-dev.yml`, `application-prod.yml`. Activate with `SPRING_PROFILES_ACTIVE` env var or `--spring.profiles.active` argument.

---

### 15.3 Senior / Lead (5+ Years)

**Q16: Explain the auto-configuration ordering mechanism.**
> Auto-configurations run in a defined order. Use `@AutoConfiguration(before=X.class, after=Y.class)` to control order. This ensures dependent configurations (like DataSource before JPA) execute correctly.

**Q17: How would you troubleshoot a bean that's not being created by auto-configuration?**
> 1) Enable `debug=true` to see CONDITIONS EVALUATION REPORT. 2) Check `/actuator/conditions` endpoint. 3) Verify class is on classpath. 4) Check if another bean of same type exists (defeating `@ConditionalOnMissingBean`). 5) Verify property conditions.

**Q18: Design a multi-tenant Spring Boot application.**
> Use a `TenantContext` (ThreadLocal) set by filter/interceptor. For DB isolation: `AbstractRoutingDataSource` to switch DataSource per tenant, or Hibernate filters for row-level isolation. Consider separate schemas or databases per tenant based on isolation requirements.

**Q19: How do you optimize Spring Boot startup time?**
> 1) `spring.main.lazy-initialization=true`. 2) Narrow `@ComponentScan`. 3) Use `@Import` instead of scanning. 4) Enable Spring AOT for native images. 5) Profile with `ApplicationStartup`. 6) Defer non-critical init to `ApplicationReadyEvent`.

**Q20: Explain Spring Boot's layered JAR and how it improves Docker builds.**
> Spring Boot 2.3+ packages JARs in layers: dependencies, spring-boot-loader, snapshot-dependencies, application. Docker can cache unchanged layers. Extract with `java -Djarmode=layertools -jar app.jar extract`, copy each layer separately.

**Q21: How do you implement graceful shutdown in a production environment?**
> Set `server.shutdown=graceful` and `spring.lifecycle.timeout-per-shutdown-phase=30s`. On SIGTERM: new requests get 503, in-flight requests complete within timeout, resources cleanup. Configure load balancer to drain connections before killing pods.

**Q22: Describe your approach to managing secrets in Spring Boot applications.**
> Never commit secrets. Options: 1) Environment variables (simple). 2) Spring Cloud Config with encryption. 3) HashiCorp Vault integration. 4) Cloud provider secrets (AWS Secrets Manager, Azure Key Vault). 5) Kubernetes Secrets. Reference via `${SECRET_NAME}` placeholders.

**Q23: How would you implement feature flags in Spring Boot?**
> Options: 1) `@ConditionalOnProperty`-based beans. 2) Custom `FeatureToggle` service reading from properties/database. 3) Libraries like FF4j, Togglz, or LaunchDarkly. 4) Spring Cloud Config for centralized control. Consider A/B testing, gradual rollout, and kill switches.

**Q24: Explain how you'd set up observability for a Spring Boot microservices architecture.**
> 1) Structured logging (JSON with correlation IDs). 2) Metrics with Micrometer + Prometheus. 3) Distributed tracing with Micrometer Tracing + Zipkin/Jaeger. 4) Centralized log aggregation (ELK, Loki). 5) Dashboards (Grafana). 6) Alerting on SLOs. Propagate trace context across services.

---

### 15.4 Scenario-Based Questions

**Q25: Your Spring Boot app is running out of database connections. How do you diagnose and fix it?**
> 1) Check HikariCP metrics (`hikaricp.connections.*`). 2) Enable leak detection (`leak-detection-threshold`). 3) Look for unclosed connections (missing `@Transactional` or try-with-resources). 4) Check for long-running transactions. 5) Tune pool size based on workload. 6) Consider connection validation settings.

**Q26: Auto-configuration is not working for a library you added. Walk through your debugging process.**
> 1) Verify library JAR is in classpath. 2) Enable `debug=true`, check CONDITIONS EVALUATION REPORT. 3) Look for failed conditions (missing class, property). 4) Check if you defined a bean that satisfies `@ConditionalOnMissingBean`. 5) Verify auto-configuration is registered in META-INF files. 6) Check for conflicting auto-configurations.

**Q27: Your application startup takes 60 seconds. How would you reduce it?**
> 1) Profile with `ApplicationStartup` to identify slow beans. 2) Enable lazy init for dev. 3) Check for slow `@PostConstruct` methods. 4) Defer Hibernate schema validation. 5) Reduce classpath scanning scope. 6) Identify slow external calls during init. 7) Consider native image for extreme cases.

**Q28: How would you migrate from Spring Boot 2.x to 3.x?**
> 1) Update Java to 17+. 2) Update dependencies (Jakarta EE namespace: `javax.*` → `jakarta.*`). 3) Update Spring Security (new DSL). 4) Update Hibernate 6 queries if needed. 5) Migrate property names (check migration guide). 6) Update test annotations. 7) Address removed/deprecated features.

**Q29: Design a health check strategy for a Spring Boot app with Redis, PostgreSQL, and external API dependencies.**
> Separate liveness from readiness. Liveness: basic app health (not dependencies). Readiness: include `db` (PostgreSQL), `redis`, custom indicator for external API. External API health should be tolerant of occasional failures (circuit breaker state). Consider synthetic transactions vs. ping checks.

**Q30: Your team is deploying Spring Boot apps to Kubernetes. What production configurations would you mandate?**
> 1) Externalized config via ConfigMaps/Secrets. 2) Liveness/readiness probes on Actuator endpoints. 3) Resource limits (CPU, memory). 4) Graceful shutdown enabled. 5) Non-root container user. 6) Structured JSON logging. 7) Prometheus metrics endpoint. 8) Connection pool tuning. 9) JVM settings optimized for containers (`-XX:+UseContainerSupport`).

---

## Quick Reference — Common Properties

```yaml
# Server
server.port=8080
server.servlet.context-path=/api
server.shutdown=graceful

# DataSource
spring.datasource.url=jdbc:postgresql://localhost/db
spring.datasource.hikari.maximum-pool-size=10

# JPA
spring.jpa.hibernate.ddl-auto=validate
spring.jpa.open-in-view=false
spring.jpa.show-sql=false

# Actuator
management.endpoints.web.exposure.include=health,info,metrics,prometheus
management.endpoint.health.show-details=when_authorized

# Logging
logging.level.root=INFO
logging.level.com.example=DEBUG

# Jackson
spring.jackson.serialization.write-dates-as-timestamps=false
spring.jackson.default-property-inclusion=non_null

# Security
spring.security.user.name=admin
spring.security.user.password=secret

# Cache
spring.cache.type=redis
spring.data.redis.host=localhost

# Flyway
spring.flyway.enabled=true
spring.flyway.locations=classpath:db/migration
```

---

## Annotation Quick Reference

| Annotation | Purpose |
|---|---|
| `@SpringBootApplication` | Main class marker (config + auto-config + scan) |
| `@EnableAutoConfiguration` | Enable auto-configuration |
| `@ConfigurationProperties` | Type-safe configuration binding |
| `@ConditionalOnClass` | Bean creation if class on classpath |
| `@ConditionalOnMissingBean` | Bean creation if no bean of type exists |
| `@ConditionalOnProperty` | Bean creation based on property value |
| `@MockBean` | Create mock and replace in context |
| `@SpyBean` | Create spy over existing bean |
| `@SpringBootTest` | Full integration test |
| `@WebMvcTest` | Web layer test slice |
| `@DataJpaTest` | JPA test slice |
| `@LocalServerPort` | Inject random test server port |
| `@TestPropertySource` | Override properties in tests |
| `@ActiveProfiles` | Activate profiles in tests |

---

> **Related Guides:**  
> - [Spring Framework Complete Guide](./Spring-Framework-Complete-Guide.md)  
> - Microservices with Spring Cloud (coming soon)
