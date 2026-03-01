# JPA and Hibernate Complete Guide

## Table of Contents
1. [Introduction](#introduction)
2. [ORM Fundamentals](#orm-fundamentals)
3. [JPA Fundamentals](#jpa-fundamentals)
4. [Hibernate Architecture](#hibernate-architecture)
5. [Hibernate Internals](#hibernate-internals)
6. [Entity Mappings](#entity-mappings)
7. [Relationships](#relationships)
8. [JPQL and Criteria API](#jpql-and-criteria-api)
9. [Spring Data JPA Advanced](#spring-data-jpa-advanced)
10. [Transactions](#transactions)
11. [Caching](#caching)
12. [Performance Optimization](#performance-optimization)
13. [Best Practices](#best-practices)
14. [Interview Questions](#interview-questions)

---

## Introduction

### What is JPA?
**Java Persistence API (JPA)** is a specification for Object-Relational Mapping (ORM) in Java. It provides a standard way to map Java objects to relational database tables.

```
┌─────────────────────────────────────────────────────────────┐
│                    Java Application                          │
├─────────────────────────────────────────────────────────────┤
│                         JPA API                              │
│              (EntityManager, Query, etc.)                    │
├─────────────────────────────────────────────────────────────┤
│                   JPA Implementation                         │
│           (Hibernate, EclipseLink, OpenJPA)                  │
├─────────────────────────────────────────────────────────────┤
│                      JDBC Driver                             │
├─────────────────────────────────────────────────────────────┤
│                   Relational Database                        │
│              (MySQL, PostgreSQL, Oracle)                     │
└─────────────────────────────────────────────────────────────┘
```

### What is Hibernate?
**Hibernate** is the most popular implementation of JPA specification. It provides additional features beyond the JPA standard.

### JPA vs Hibernate

| Aspect | JPA | Hibernate |
|--------|-----|-----------|
| Type | Specification | Implementation |
| Standardization | Java EE Standard | Proprietary + JPA |
| Portability | High (switch implementations) | Tied to Hibernate |
| Features | Core ORM features | Extended features |
| Annotations | `javax.persistence.*` | `org.hibernate.annotations.*` |

---

## ORM Fundamentals

### What is Object-Relational Mapping (ORM)?

ORM is a programming technique that creates a "virtual object database" by mapping objects in code to tables in a relational database. It acts as a bridge between two fundamentally different paradigms.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         ORM Conceptual Model                                 │
│                                                                              │
│   Object-Oriented World              ORM Layer              Relational World │
│   ──────────────────────            ─────────             ─────────────────  │
│                                                                              │
│   ┌─────────────────┐         ┌─────────────────┐       ┌────────────────┐  │
│   │    Classes      │ ◄─────► │   Metadata      │ ◄───► │    Tables      │  │
│   └─────────────────┘         │   Mappings      │       └────────────────┘  │
│                               └─────────────────┘                            │
│   ┌─────────────────┐         ┌─────────────────┐       ┌────────────────┐  │
│   │   Properties    │ ◄─────► │   Type          │ ◄───► │    Columns     │  │
│   │   (Fields)      │         │   Converters    │       │    (Fields)    │  │
│   └─────────────────┘         └─────────────────┘       └────────────────┘  │
│                                                                              │
│   ┌─────────────────┐         ┌─────────────────┐       ┌────────────────┐  │
│   │   References    │ ◄─────► │   Relationship  │ ◄───► │  Foreign Keys  │  │
│   │   (Associations)│         │   Strategies    │       │    (Joins)     │  │
│   └─────────────────┘         └─────────────────┘       └────────────────┘  │
│                                                                              │
│   ┌─────────────────┐         ┌─────────────────┐       ┌────────────────┐  │
│   │   Inheritance   │ ◄─────► │   Inheritance   │ ◄───► │  Table         │  │
│   │   Hierarchy     │         │   Mapping       │       │  Strategies    │  │
│   └─────────────────┘         └─────────────────┘       └────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Object-Relational Impedance Mismatch (Deep Dive)

The **Object-Relational Impedance Mismatch** describes the fundamental incompatibilities between object-oriented programming and relational databases. This is the core problem that ORM frameworks attempt to solve.

#### The Five Major Mismatches

```
┌─────────────────────────────────────────────────────────────────────────────┐
│              Object-Relational Impedance Mismatch Problems                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  1. GRANULARITY MISMATCH                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Object World: Multiple fine-grained objects                        │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                  │    │
│  │  │   Person    │  │   Address   │  │   Phone     │                  │    │
│  │  └─────────────┘  └─────────────┘  └─────────────┘                  │    │
│  │                        vs                                            │    │
│  │  Relational World: Fewer coarse-grained tables                       │    │
│  │  ┌─────────────────────────────────────────────────────────────┐    │    │
│  │  │                    PERSON_DATA                               │    │    │
│  │  │  (all data flattened into single table)                      │    │    │
│  │  └─────────────────────────────────────────────────────────────┘    │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  2. INHERITANCE MISMATCH                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Object World: Natural inheritance hierarchy                         │    │
│  │           ┌─────────┐                                                │    │
│  │           │ Vehicle │                                                │    │
│  │           └────┬────┘                                                │    │
│  │          ┌─────┴─────┐                                               │    │
│  │     ┌────┴────┐ ┌────┴────┐                                          │    │
│  │     │   Car   │ │  Bike   │                                          │    │
│  │     └─────────┘ └─────────┘                                          │    │
│  │                        vs                                            │    │
│  │  Relational World: No native inheritance support                     │    │
│  │  Must use: Single Table / Joined Tables / Table-per-Class            │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  3. IDENTITY MISMATCH                                                        │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Object World: Two types of identity                                 │    │
│  │  • Reference Identity: obj1 == obj2 (same memory address)            │    │
│  │  • Value Equality: obj1.equals(obj2) (same content)                  │    │
│  │                        vs                                            │    │
│  │  Relational World: Single identity concept                           │    │
│  │  • Primary Key equality only                                         │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  4. ASSOCIATION MISMATCH                                                     │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Object World: References are directional                            │    │
│  │  Employee ───────────► Department (unidirectional)                   │    │
│  │  Employee ◄────────────► Department (bidirectional, 2 references)    │    │
│  │                        vs                                            │    │
│  │  Relational World: Foreign keys are inherently bidirectional         │    │
│  │  EMPLOYEE.dept_id ────FK──── DEPARTMENT.id                           │    │
│  │  (Can join from either table)                                        │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
│  5. DATA NAVIGATION MISMATCH                                                 │
│  ┌─────────────────────────────────────────────────────────────────────┐    │
│  │  Object World: Navigate object graph                                 │    │
│  │  employee.getDepartment().getManager().getAddress().getCity()        │    │
│  │  (Multiple hops through memory references)                           │    │
│  │                        vs                                            │    │
│  │  Relational World: JOIN-based access                                 │    │
│  │  SELECT ... FROM emp e                                               │    │
│  │  JOIN department d ON e.dept_id = d.id                               │    │
│  │  JOIN employee m ON d.manager_id = m.id                              │    │
│  │  JOIN address a ON m.address_id = a.id                               │    │
│  │  (Single query with multiple JOINs is efficient)                     │    │
│  └─────────────────────────────────────────────────────────────────────┘    │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### 1. Granularity Mismatch (Detailed)

**Problem**: Objects can have varying granularity, but tables are typically coarser.

```java
// Object-Oriented Design: Fine-grained objects
public class Person {
    private Long id;
    private String name;
    private Address homeAddress;    // Separate object
    private Address workAddress;    // Separate object  
    private List<Phone> phones;     // Collection of objects
}

public class Address {
    private String street;
    private String city;
    private String state;
    private String zipCode;
    private GeoLocation location;   // Even finer granularity
}

public class GeoLocation {
    private double latitude;
    private double longitude;
}
```

**ORM Solutions**:

```java
// Solution 1: @Embedded (Component mapping)
@Entity
public class Person {
    @Id
    private Long id;
    
    @Embedded
    @AttributeOverrides({
        @AttributeOverride(name = "street", column = @Column(name = "home_street")),
        @AttributeOverride(name = "city", column = @Column(name = "home_city"))
    })
    private Address homeAddress;  // Stored in same table
    
    @Embedded
    @AttributeOverrides({
        @AttributeOverride(name = "street", column = @Column(name = "work_street")),
        @AttributeOverride(name = "city", column = @Column(name = "work_city"))
    })
    private Address workAddress;  // Stored in same table
}

// Solution 2: Separate entity with @OneToOne
@Entity
public class Person {
    @Id
    private Long id;
    
    @OneToOne(cascade = CascadeType.ALL)
    @JoinColumn(name = "home_address_id")
    private Address homeAddress;  // Separate table
}
```

#### 2. Inheritance Mismatch (Detailed)

**Problem**: SQL has no concept of inheritance.

```java
// Natural OO inheritance
public abstract class BillingDetails {
    protected Long id;
    protected String owner;
}

public class CreditCard extends BillingDetails {
    private String cardNumber;
    private String expMonth;
    private String expYear;
}

public class BankAccount extends BillingDetails {
    private String accountNumber;
    private String bankName;
    private String routingNumber;
}
```

**ORM Strategy Trade-offs**:

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                     Inheritance Mapping Strategies                            │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                               │
│  SINGLE_TABLE                                                                 │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │  billing_details                                                        │ │
│  │  ├── id (PK)                                                            │ │
│  │  ├── dtype (discriminator: 'CC' or 'BA')                                │ │
│  │  ├── owner                                                              │ │
│  │  ├── card_number (NULL for bank accounts)                               │ │
│  │  ├── exp_month (NULL for bank accounts)                                 │ │
│  │  ├── account_number (NULL for credit cards)                             │ │
│  │  └── bank_name (NULL for credit cards)                                  │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│  ✓ Best polymorphic query performance (single table scan)                    │
│  ✓ Simple schema                                                             │
│  ✗ NULL columns waste space                                                  │
│  ✗ Cannot have NOT NULL constraints on subclass columns                      │
│                                                                               │
│  JOINED                                                                       │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │  billing_details          credit_card           bank_account            │ │
│  │  ├── id (PK)              ├── id (PK,FK)        ├── id (PK,FK)          │ │
│  │  └── owner                ├── card_number       ├── account_number      │ │
│  │                           └── exp_month         └── bank_name           │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│  ✓ Normalized, no NULL columns                                               │
│  ✓ Can have NOT NULL constraints                                             │
│  ✗ Polymorphic queries require JOINs                                         │
│  ✗ Deep hierarchies = many JOINs = slow                                      │
│                                                                               │
│  TABLE_PER_CLASS                                                              │
│  ┌─────────────────────────────────────────────────────────────────────────┐ │
│  │  credit_card (complete)        bank_account (complete)                  │ │
│  │  ├── id (PK)                   ├── id (PK)                              │ │
│  │  ├── owner                     ├── owner                                │ │
│  │  ├── card_number               ├── account_number                       │ │
│  │  └── exp_month                 └── bank_name                            │ │
│  └─────────────────────────────────────────────────────────────────────────┘ │
│  ✓ No JOINs for concrete class queries                                       │
│  ✗ Polymorphic queries require UNION                                         │
│  ✗ Difficult to maintain (schema changes affect multiple tables)             │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘
```

#### 3. Identity Mismatch (Detailed)

**Problem**: Java has two notions of "same", databases have one.

```java
// Object identity confusion
Employee emp1 = repository.findById(1L);  // Load from DB
Employee emp2 = repository.findById(1L);  // Load again

// Within same persistence context: emp1 == emp2 (same object)
// Different persistence contexts: emp1 != emp2 (different objects)
// But: emp1.equals(emp2) should be true (same business identity)

// The mismatch:
// 1. Database: row with id=1 is THE identity
// 2. Java: which equality matters? == or equals()?
```

**ORM Solution - Proper equals/hashCode**:

```java
@Entity
public class Employee {
    @Id
    @GeneratedValue
    private Long id;  // Database identity
    
    @NaturalId  // Business identity (Hibernate-specific)
    @Column(unique = true, updatable = false)
    private String employeeNumber;
    
    // WRONG: Using database ID (fails for transient objects)
    // @Override
    // public boolean equals(Object o) {
    //     return Objects.equals(id, ((Employee) o).id);  // id is null before persist!
    // }
    
    // CORRECT: Use business key
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Employee)) return false;
        Employee that = (Employee) o;
        return Objects.equals(employeeNumber, that.employeeNumber);
    }
    
    @Override
    public int hashCode() {
        return Objects.hash(employeeNumber);  // Immutable field!
    }
}
```

**Persistence Context Identity Guarantee**:

```java
@Transactional
public void demonstrateIdentity() {
    // Within SAME persistence context
    Employee emp1 = entityManager.find(Employee.class, 1L);
    Employee emp2 = entityManager.find(Employee.class, 1L);
    
    assert emp1 == emp2;  // TRUE - same object from first-level cache
    
    // This is the ORM's identity map pattern at work
}
```

#### 4. Association Mismatch (Detailed)

**Problem**: Object references vs. Foreign keys

```java
// Objects: Directional references
public class Order {
    private Customer customer;  // Order knows about Customer
}

public class Customer {
    // Does Customer know about Orders? Depends on your design!
    // private List<Order> orders;  // Optional bidirectional
}

// SQL: Foreign keys are inherently bidirectional
// SELECT * FROM orders WHERE customer_id = 123;  (Order -> Customer)
// SELECT * FROM orders o JOIN customers c ON o.customer_id = c.id;  (Either direction)
```

**ORM Challenge - Bidirectional Sync**:

```java
@Entity
public class Department {
    @Id
    private Long id;
    
    // Inverse side (mappedBy) - doesn't own the FK
    @OneToMany(mappedBy = "department")
    private List<Employee> employees = new ArrayList<>();
    
    // CRITICAL: Helper methods to maintain consistency
    public void addEmployee(Employee emp) {
        employees.add(emp);
        emp.setDepartment(this);  // Must set both sides!
    }
}

@Entity
public class Employee {
    @Id
    private Long id;
    
    // Owning side - has the FK column
    @ManyToOne
    @JoinColumn(name = "department_id")
    private Department department;  // Only this side is persisted!
}

// Common bug:
department.getEmployees().add(newEmployee);  // WRONG! FK not set
newEmployee.setDepartment(department);        // CORRECT! FK will be set
```

#### 5. Data Navigation Mismatch (Detailed)

**Problem**: Object traversal vs. Set-based operations

```java
// Object-oriented navigation (inefficient in SQL world)
for (Order order : customer.getOrders()) {
    for (LineItem item : order.getLineItems()) {
        Product product = item.getProduct();
        Category category = product.getCategory();
        // Multiple round trips to database!
    }
}

// SQL approach (efficient: one query)
// SELECT c.*, o.*, li.*, p.*, cat.*
// FROM customers c
// JOIN orders o ON c.id = o.customer_id
// JOIN line_items li ON o.id = li.order_id
// JOIN products p ON li.product_id = p.id
// JOIN categories cat ON p.category_id = cat.id
// WHERE c.id = ?
```

**ORM Solutions**:

```java
// Solution 1: Eager fetching (careful with cartesian products!)
@ManyToMany(fetch = FetchType.EAGER)
private Set<Category> categories;

// Solution 2: JOIN FETCH in queries
@Query("SELECT o FROM Order o " +
       "JOIN FETCH o.lineItems li " +
       "JOIN FETCH li.product p " +
       "WHERE o.customer.id = :customerId")
List<Order> findOrdersWithDetails(@Param("customerId") Long customerId);

// Solution 3: Entity Graph
@NamedEntityGraph(
    name = "Order.withDetails",
    attributeNodes = {
        @NamedAttributeNode(value = "lineItems", subgraph = "lineItems-subgraph")
    },
    subgraphs = {
        @NamedSubgraph(
            name = "lineItems-subgraph",
            attributeNodes = {@NamedAttributeNode("product")}
        )
    }
)
@Entity
public class Order { ... }

// Solution 4: DTO Projection (best for read-only)
@Query("SELECT new com.example.OrderSummaryDTO(o.id, o.total, c.name) " +
       "FROM Order o JOIN o.customer c")
List<OrderSummaryDTO> findOrderSummaries();
```

### Why ORM Despite The Mismatch?

| Without ORM | With ORM |
|-------------|----------|
| Write SQL for every operation | Automatic SQL generation |
| Manual result set mapping | Automatic object hydration |
| No dirty tracking | Automatic change detection |
| Manual transaction handling | Declarative transactions |
| Type-unsafe string queries | Type-safe criteria/JPQL |
| Database-specific SQL | Database portability |
| Manual caching | Built-in caching |

---

## JPA Fundamentals

### Core Components

#### 1. Entity
An entity represents a table in the database.

```java
import javax.persistence.*;

@Entity
@Table(name = "employees")
public class Employee {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(name = "first_name", nullable = false, length = 50)
    private String firstName;
    
    @Column(name = "last_name", nullable = false, length = 50)
    private String lastName;
    
    @Column(unique = true)
    private String email;
    
    @Temporal(TemporalType.DATE)
    private Date hireDate;
    
    @Enumerated(EnumType.STRING)
    private EmployeeStatus status;
    
    // Constructors, getters, setters
}
```

#### 2. EntityManager
The central interface for persistence operations.

```java
@Service
public class EmployeeService {
    
    @PersistenceContext
    private EntityManager entityManager;
    
    // Persist - Insert new entity
    public void createEmployee(Employee employee) {
        entityManager.persist(employee);
    }
    
    // Find - Retrieve by primary key
    public Employee findEmployee(Long id) {
        return entityManager.find(Employee.class, id);
    }
    
    // Merge - Update detached entity
    public Employee updateEmployee(Employee employee) {
        return entityManager.merge(employee);
    }
    
    // Remove - Delete entity
    public void deleteEmployee(Long id) {
        Employee employee = entityManager.find(Employee.class, id);
        if (employee != null) {
            entityManager.remove(employee);
        }
    }
    
    // Refresh - Reload from database
    public void refreshEmployee(Employee employee) {
        entityManager.refresh(employee);
    }
    
    // Detach - Remove from persistence context
    public void detachEmployee(Employee employee) {
        entityManager.detach(employee);
    }
}
```

#### 3. Persistence Context
The persistence context is a set of managed entity instances. Think of it as a first-level cache.

```
┌─────────────────────────────────────────────────────────────────┐
│                     Persistence Context                          │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                  Managed Entities                        │    │
│  │   Employee@1 ──────────────────────────────────────────┐│    │
│  │   Employee@2 ─────────────────────────────────────────┐││    │
│  │   Department@1 ──────────────────────────────────────┐│││    │
│  └─────────────────────────────────────────────────────┘││││    │
│                                                         │││┘    │
│                                                         ││┘     │
│                                                         │┘      │
├─────────────────────────────────────────────────────────────────┤
│                        Database                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Entity Lifecycle States

```
┌──────────────────────────────────────────────────────────────────────┐
│                        Entity Lifecycle                               │
│                                                                       │
│    ┌─────────┐     persist()      ┌─────────┐                        │
│    │   NEW   │ ─────────────────> │ MANAGED │                        │
│    │(Transient)│                   │         │                        │
│    └─────────┘                     └────┬────┘                        │
│         ^                               │                             │
│         │                               │ remove()                    │
│         │ new                           v                             │
│         │                          ┌─────────┐                        │
│    ┌────┴────┐    merge()          │ REMOVED │                        │
│    │ DETACHED│ <─────────────────  └─────────┘                        │
│    │         │                          │                             │
│    └─────────┘                          │ flush()/commit()            │
│         ^                               v                             │
│         │ detach()/clear()         ┌─────────┐                        │
│         │ close()/serialize        │ DATABASE│                        │
│         └───────────────────────── └─────────┘                        │
│                                                                       │
└──────────────────────────────────────────────────────────────────────┘
```

| State | Description |
|-------|-------------|
| **New/Transient** | Entity created with `new`, not associated with persistence context |
| **Managed** | Entity associated with persistence context, changes auto-synced |
| **Detached** | Entity was managed but persistence context closed/cleared |
| **Removed** | Entity marked for deletion, will be deleted on flush |

### ID Generation Strategies

```java
// 1. IDENTITY - Database auto-increment
@Id
@GeneratedValue(strategy = GenerationType.IDENTITY)
private Long id;

// 2. SEQUENCE - Database sequence (preferred for batch inserts)
@Id
@GeneratedValue(strategy = GenerationType.SEQUENCE, generator = "emp_seq")
@SequenceGenerator(name = "emp_seq", sequenceName = "employee_sequence", allocationSize = 50)
private Long id;

// 3. TABLE - Simulated sequence using a table
@Id
@GeneratedValue(strategy = GenerationType.TABLE, generator = "emp_gen")
@TableGenerator(name = "emp_gen", table = "id_generator", pkColumnName = "gen_name",
                valueColumnName = "gen_value", allocationSize = 50)
private Long id;

// 4. UUID - Universally unique identifier
@Id
@GeneratedValue(generator = "uuid2")
@GenericGenerator(name = "uuid2", strategy = "uuid2")
@Column(columnDefinition = "BINARY(16)")
private UUID id;

// 5. AUTO - Let Hibernate decide
@Id
@GeneratedValue(strategy = GenerationType.AUTO)
private Long id;
```

**Strategy Comparison:**

| Strategy | Pros | Cons | Use Case |
|----------|------|------|----------|
| IDENTITY | Simple, native DB support | Breaks batch inserts | Single inserts |
| SEQUENCE | Batch-friendly, performant | Not all DBs support | High-volume apps |
| TABLE | Portable across DBs | Performance overhead | Legacy systems |
| UUID | Globally unique, no DB call | Larger storage size | Distributed systems |

---

## Hibernate Architecture

```
┌───────────────────────────────────────────────────────────────────────────┐
│                         Hibernate Architecture                             │
│                                                                            │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │                        Application Layer                              │ │
│  │    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐                 │ │
│  │    │   Domain    │  │   Service   │  │    DAO      │                 │ │
│  │    │   Objects   │  │    Layer    │  │    Layer    │                 │ │
│  │    └─────────────┘  └─────────────┘  └─────────────┘                 │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
│                                    │                                       │
│                                    v                                       │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │                        Hibernate Core                                 │ │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐   │ │
│  │  │ SessionFactory  │  │     Session     │  │    Transaction      │   │ │
│  │  │   (Singleton)   │  │  (Per Request)  │  │    Management       │   │ │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────┘   │ │
│  │                                                                       │ │
│  │  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────────┐   │ │
│  │  │   First-Level   │  │  Second-Level   │  │   Query Cache       │   │ │
│  │  │     Cache       │  │     Cache       │  │                     │   │ │
│  │  └─────────────────┘  └─────────────────┘  └─────────────────────┘   │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
│                                    │                                       │
│                                    v                                       │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │                           JDBC Layer                                  │ │
│  │      Connection Pool  │  Statement Cache  │  Result Set Handling     │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
│                                    │                                       │
│                                    v                                       │
│  ┌──────────────────────────────────────────────────────────────────────┐ │
│  │                          Database Layer                               │ │
│  │             MySQL  │  PostgreSQL  │  Oracle  │  SQL Server            │ │
│  └──────────────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────────────┘
```

### SessionFactory vs EntityManagerFactory

```java
// JPA Standard - EntityManagerFactory
@Configuration
public class JpaConfig {
    
    @Bean
    public LocalContainerEntityManagerFactoryBean entityManagerFactory(
            DataSource dataSource) {
        LocalContainerEntityManagerFactoryBean em = 
            new LocalContainerEntityManagerFactoryBean();
        em.setDataSource(dataSource);
        em.setPackagesToScan("com.example.entity");
        em.setJpaVendorAdapter(new HibernateJpaVendorAdapter());
        em.setJpaProperties(hibernateProperties());
        return em;
    }
    
    private Properties hibernateProperties() {
        Properties properties = new Properties();
        properties.put("hibernate.dialect", "org.hibernate.dialect.MySQL8Dialect");
        properties.put("hibernate.show_sql", "true");
        properties.put("hibernate.format_sql", "true");
        properties.put("hibernate.hbm2ddl.auto", "update");
        return properties;
    }
}

// Hibernate Native - SessionFactory
@Configuration
public class HibernateConfig {
    
    @Bean
    public LocalSessionFactoryBean sessionFactory(DataSource dataSource) {
        LocalSessionFactoryBean sessionFactory = new LocalSessionFactoryBean();
        sessionFactory.setDataSource(dataSource);
        sessionFactory.setPackagesToScan("com.example.entity");
        sessionFactory.setHibernateProperties(hibernateProperties());
        return sessionFactory;
    }
}
```

---

## Hibernate Internals

### Session Internals Deep Dive

The Hibernate Session (JPA EntityManager) is the primary interface for persistence operations. Understanding its internals is crucial for performance optimization.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                         Session Internal Structure                           │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │                     Persistence Context                                  ││
│  │  ┌──────────────────────────────────────────────────────────────────┐   ││
│  │  │                  ENTITY SNAPSHOT MAP                              │   ││
│  │  │   Key: EntityKey(Class, ID)  │  Value: Object[] (loaded state)   │   ││
│  │  │   ───────────────────────────┼───────────────────────────────    │   ││
│  │  │   (Employee, 1)              │  ["John", "Doe", 50000, ...]      │   ││
│  │  │   (Employee, 2)              │  ["Jane", "Smith", 60000, ...]    │   ││
│  │  │   (Department, 1)            │  ["Engineering", "NYC", ...]      │   ││
│  │  └──────────────────────────────────────────────────────────────────┘   ││
│  │                                                                          ││
│  │  ┌──────────────────────────────────────────────────────────────────┐   ││
│  │  │                  ENTITY INSTANCE MAP                              │   ││
│  │  │   Key: EntityKey(Class, ID)  │  Value: Entity Instance           │   ││
│  │  │   ───────────────────────────┼───────────────────────────────    │   ││
│  │  │   (Employee, 1)              │  Employee@a1b2c3                  │   ││
│  │  │   (Employee, 2)              │  Employee@d4e5f6                  │   ││
│  │  └──────────────────────────────────────────────────────────────────┘   ││
│  │                                                                          ││
│  │  ┌──────────────────────────────────────────────────────────────────┐   ││
│  │  │                     ACTION QUEUE                                  │   ││
│  │  │   ┌─────────────┐ ┌─────────────┐ ┌─────────────┐                │   ││
│  │  │   │  Insertions │ │   Updates   │ │  Deletions  │                │   ││
│  │  │   │  (ordered)  │ │  (ordered)  │ │  (ordered)  │                │   ││
│  │  │   └─────────────┘ └─────────────┘ └─────────────┘                │   ││
│  │  │   ┌─────────────┐ ┌─────────────┐                                │   ││
│  │  │   │ Collection  │ │ Collection  │                                │   ││
│  │  │   │   Updates   │ │  Removals   │                                │   ││
│  │  │   └─────────────┘ └─────────────┘                                │   ││
│  │  └──────────────────────────────────────────────────────────────────┘   ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │                        JDBC Connection                                   ││
│  │   • Connection acquired lazily (on first SQL)                            ││
│  │   • Released on transaction commit/rollback                              ││
│  │   • Connection pool integration                                          ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Dirty Checking Mechanism

Dirty checking is Hibernate's automatic detection of entity changes. Understanding this helps optimize performance.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                        Dirty Checking Process                                │
│                                                                              │
│  1. ENTITY LOAD                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Employee emp = em.find(Employee.class, 1L);                           │ │
│  │                                                                         │ │
│  │  ┌──────────────────────┐   ┌──────────────────────────────────┐       │ │
│  │  │   Entity Instance    │   │        Snapshot (Copy)            │       │ │
│  │  │   emp.name = "John"  │   │   snapshot[0] = "John"            │       │ │
│  │  │   emp.salary = 50000 │   │   snapshot[1] = 50000             │       │ │
│  │  └──────────────────────┘   └──────────────────────────────────┘       │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  2. MODIFY ENTITY                                                            │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  emp.setSalary(55000);  // Just a setter call                          │ │
│  │                                                                         │ │
│  │  ┌──────────────────────┐   ┌──────────────────────────────────┐       │ │
│  │  │   Entity Instance    │   │        Snapshot (unchanged)       │       │ │
│  │  │   emp.name = "John"  │   │   snapshot[0] = "John"            │       │ │
│  │  │   emp.salary = 55000 │   │   snapshot[1] = 50000  ← original │       │ │
│  │  │              ↑ new   │   └──────────────────────────────────┘       │ │
│  │  └──────────────────────┘                                               │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  3. FLUSH (Dirty Check)                                                      │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  // At flush time, Hibernate compares:                                  │ │
│  │  for (each managed entity) {                                            │ │
│  │      Object[] currentState = getPropertyValues(entity);                 │ │
│  │      Object[] loadedState = snapshot.get(entity);                       │ │
│  │                                                                         │ │
│  │      for (int i = 0; i < properties.length; i++) {                      │ │
│  │          if (!equals(currentState[i], loadedState[i])) {                │ │
│  │              // Entity is DIRTY - schedule UPDATE                       │ │
│  │              actionQueue.addUpdate(entity, changedProperties);          │ │
│  │          }                                                              │ │
│  │      }                                                                  │ │
│  │  }                                                                      │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
│  4. EXECUTE SQL                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  UPDATE employees SET salary = 55000 WHERE id = 1                       │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Dirty Checking Strategies

```java
// 1. Default: Field-by-field comparison
@Entity
public class Employee {
    // All fields compared during dirty check
}

// 2. @DynamicUpdate: Only update changed columns
@Entity
@DynamicUpdate
public class Employee {
    // Generates: UPDATE emp SET salary=? WHERE id=?
    // Instead of: UPDATE emp SET name=?, salary=?, dept=?... WHERE id=?
}

// 3. @SelectBeforeUpdate: Check DB before update (for detached entities)
@Entity
@SelectBeforeUpdate
public class Employee {
    // Issues SELECT before UPDATE to verify current state
}

// 4. @Immutable: Skip dirty checking entirely
@Entity
@Immutable
public class AuditLog {
    // Never checked for changes, never updated
}

// 5. Manual dirty flag (custom interceptor)
@Entity
public class Employee implements SelfDirtinessTracker {
    @Transient
    private Set<String> dirtyFields = new HashSet<>();
    
    public void setSalary(BigDecimal salary) {
        if (!Objects.equals(this.salary, salary)) {
            dirtyFields.add("salary");
        }
        this.salary = salary;
    }
    
    @Override
    public boolean $$_hibernate_hasDirtyAttributes() {
        return !dirtyFields.isEmpty();
    }
}
```

#### Bytecode Enhancement for Dirty Tracking

```xml
<!-- Enable bytecode enhancement in build -->
<plugin>
    <groupId>org.hibernate.orm.tooling</groupId>
    <artifactId>hibernate-enhance-maven-plugin</artifactId>
    <executions>
        <execution>
            <configuration>
                <enableDirtyTracking>true</enableDirtyTracking>
                <enableLazyInitialization>true</enableLazyInitialization>
            </configuration>
        </execution>
    </executions>
</plugin>
```

```java
// With bytecode enhancement, entity automatically tracks changes:
@Entity
public class Employee {
    private String name;
    
    // After enhancement, setter becomes:
    public void setName(String name) {
        if (!Objects.equals(this.name, name)) {
            $$_hibernate_trackChange("name");  // Auto-injected!
        }
        this.name = name;
    }
}
```

### Flush Modes

Flush mode determines WHEN Hibernate synchronizes persistence context with database.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Flush Mode Comparison                              │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  FlushMode.AUTO (Default)                                                    │
│  ─────────────────────────                                                   │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  persist(emp) ──► Query employees ──► Flush triggered ──► Execute SQL  │ │
│  │       │                   │                │                            │ │
│  │       │                   │                └── Before query that might  │ │
│  │       │                   │                    see stale data           │ │
│  │       │                   │                                             │ │
│  │       ▼                   ▼                                             │ │
│  │  Action Queue    "SELECT from employees"                                │ │
│  │  has pending     would miss the new emp                                 │ │
│  │  insert          without flush                                          │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│  Flushes: before query, before commit                                        │
│                                                                              │
│  FlushMode.COMMIT                                                            │
│  ─────────────────                                                           │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  persist(emp) ──► Query employees ──► NO flush ──► Commit ──► Flush    │ │
│  │       │                   │              │            │                 │ │
│  │       │                   │              │            └── Only now!     │ │
│  │       │                   │              │                              │ │
│  │       │                   │              └── Query returns stale data   │ │
│  │       ▼                   ▼                                             │ │
│  │  Action Queue      May not see                                          │ │
│  │  accumulates       pending changes                                      │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│  Flushes: before commit only                                                 │
│  Use case: Batch processing (avoid intermediate flushes)                     │
│                                                                              │
│  FlushMode.MANUAL                                                            │
│  ────────────────                                                            │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Changes accumulate ──► Must call em.flush() explicitly                │ │
│  │                                                                         │ │
│  │  Even commit() won't flush unless you call flush()!                    │ │
│  │  (Spring's @Transactional may auto-flush on commit depending on impl)  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│  Use case: Maximum control, read-heavy operations                            │
│                                                                              │
│  FlushMode.ALWAYS (Hibernate-specific)                                       │
│  ─────────────────                                                           │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │  Flushes before EVERY query, even native SQL                           │ │
│  │  Performance impact: significant overhead                              │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│  Use case: Rarely needed, when native SQL must see all pending changes       │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

#### Setting Flush Mode

```java
// 1. Per EntityManager
@PersistenceContext
private EntityManager em;

public void batchProcess() {
    em.setFlushMode(FlushModeType.COMMIT);  // JPA standard
    // or
    em.unwrap(Session.class).setHibernateFlushMode(FlushMode.MANUAL);  // Hibernate
}

// 2. Per Query
@Query("SELECT e FROM Employee e WHERE e.status = :status")
@QueryHints(@QueryHint(name = "org.hibernate.flushMode", value = "COMMIT"))
List<Employee> findByStatus(@Param("status") String status);

// 3. Global configuration
spring.jpa.properties.org.hibernate.flushMode=COMMIT
```

### Flush Order and Action Queue

Hibernate executes SQL in a specific order to maintain referential integrity:

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          Flush Execution Order                               │
│                                                                              │
│  1. Insert entities (in order of dependency)                                 │
│     └── Parents before children (to satisfy FK constraints)                  │
│                                                                              │
│  2. Update entities                                                          │
│     └── Changed entities in any order                                        │
│                                                                              │
│  3. Delete collection elements                                               │
│     └── Remove from join tables, orphaned children                           │
│                                                                              │
│  4. Insert collection elements                                               │
│     └── Add to join tables                                                   │
│                                                                              │
│  5. Delete entities                                                          │
│     └── Children before parents (to satisfy FK constraints)                  │
│                                                                              │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │ Example: Delete Department with Employees                                ││
│  │                                                                          ││
│  │  // In persistence context:                                              ││
│  │  em.remove(department);                                                  ││
│  │                                                                          ││
│  │  // At flush time, Hibernate generates:                                  ││
│  │  DELETE FROM employee WHERE department_id = ?  // Children first         ││
│  │  DELETE FROM department WHERE id = ?            // Then parent           ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Persistence Context Modes

```java
// 1. Transaction-scoped (Default with Spring)
@PersistenceContext
private EntityManager em;  // New context per transaction

// 2. Extended persistence context
@PersistenceContext(type = PersistenceContextType.EXTENDED)
private EntityManager em;  // Spans multiple transactions (stateful beans)

// Usage difference:
@Transactional
public void transactionScoped() {
    Employee emp = em.find(Employee.class, 1L);  // Managed
}  // Transaction ends, emp becomes DETACHED

public void laterAccess() {
    emp.getDepartment();  // LazyInitializationException!
}

// With EXTENDED:
@Stateful  // EJB or similar
public class EmployeeEditor {
    @PersistenceContext(type = PersistenceContextType.EXTENDED)
    private EntityManager em;
    
    private Employee emp;
    
    public void load(Long id) {
        emp = em.find(Employee.class, id);  // Managed
    }
    
    public void modifyName(String name) {
        emp.setName(name);  // Still managed! No LazyInit issues
    }
    
    @TransactionAttribute(REQUIRES_NEW)
    public void save() {
        // Changes automatically saved on next flush
    }
}
```

### Session Statistics and Monitoring

```java
// Enable statistics
spring.jpa.properties.hibernate.generate_statistics=true

// Access statistics
@Service
public class HibernateMonitor {
    
    @PersistenceUnit
    private EntityManagerFactory emf;
    
    public void logStatistics() {
        Statistics stats = emf.unwrap(SessionFactory.class).getStatistics();
        
        log.info("=================== Hibernate Statistics ===================");
        log.info("Queries executed: {}", stats.getQueryExecutionCount());
        log.info("Second-level cache hit: {}", stats.getSecondLevelCacheHitCount());
        log.info("Second-level cache miss: {}", stats.getSecondLevelCacheMissCount());
        log.info("Entities loaded: {}", stats.getEntityLoadCount());
        log.info("Entities inserted: {}", stats.getEntityInsertCount());
        log.info("Entities updated: {}", stats.getEntityUpdateCount());
        log.info("Collections loaded: {}", stats.getCollectionLoadCount());
        log.info("Flush count: {}", stats.getFlushCount());
        log.info("Session open count: {}", stats.getSessionOpenCount());
        log.info("Session close count: {}", stats.getSessionCloseCount());
        
        // Slowest query
        log.info("Slowest query: {} ({}ms)", 
            stats.getQueryStatistics(stats.getQueries()[0]).getExecutionMaxTime());
        
        stats.clear();  // Reset for next interval
    }
}
```

### Understanding Proxies and Lazy Loading Internals

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                      Hibernate Proxy Mechanism                               │
│                                                                              │
│  em.getReference(Employee.class, 1L)                                         │
│                    │                                                         │
│                    ▼                                                         │
│  ┌─────────────────────────────────────────────────────────────────────────┐│
│  │                    Employee$HibernateProxy$abc123                        ││
│  │  ┌───────────────────────────────────────────────────────────────────┐  ││
│  │  │  - target: null (not yet loaded)                                  │  ││
│  │  │  - id: 1L (only this is known)                                    │  ││
│  │  │  - session: Session@xyz (reference to load when needed)           │  ││
│  │  │                                                                   │  ││
│  │  │  Methods:                                                         │  ││
│  │  │  getId() → returns 1L (no DB hit!)                                │  ││
│  │  │  getName() → triggers initialization → loads from DB             │  ││
│  │  │                                                                   │  ││
│  │  │  Initialization:                                                  │  ││
│  │  │  1. Check if target is null                                       │  ││
│  │  │  2. If null and session is open → SELECT * FROM employee WHERE id=1 ││
│  │  │  3. If null and session closed → LazyInitializationException     │  ││
│  │  │  4. Set target to loaded Employee                                 │  ││
│  │  │  5. Delegate method call to target                                │  ││
│  │  └───────────────────────────────────────────────────────────────────┘  ││
│  └─────────────────────────────────────────────────────────────────────────┘│
│                                                                              │
│  LazyInitializationException Prevention:                                     │
│  ──────────────────────────────────────                                      │
│  1. Access within transaction (session still open)                           │
│  2. Use JOIN FETCH in query                                                  │
│  3. Use @EntityGraph                                                         │
│  4. Hibernate.initialize(proxy) within transaction                           │
│  5. Open Session in View (OSIV) - NOT recommended for APIs                   │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

```java
// Check if proxy is initialized
if (Hibernate.isInitialized(employee.getDepartment())) {
    // Safe to access
}

// Force initialization
@Transactional(readOnly = true)
public Employee getEmployeeWithDepartment(Long id) {
    Employee emp = em.find(Employee.class, id);
    Hibernate.initialize(emp.getDepartment());  // Force load
    return emp;
}

// Unproxy to get real entity
Employee realEmployee = Hibernate.unproxy(proxyEmployee, Employee.class);
```

---

## Entity Mappings

### Basic Column Mappings

```java
@Entity
@Table(name = "products", 
       uniqueConstraints = @UniqueConstraint(columnNames = {"sku", "warehouse_id"}),
       indexes = @Index(name = "idx_product_name", columnList = "name"))
public class Product {
    
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @Column(name = "product_name", nullable = false, length = 100)
    private String name;
    
    @Column(precision = 10, scale = 2)
    private BigDecimal price;
    
    @Column(columnDefinition = "TEXT")
    private String description;
    
    @Lob
    private byte[] image;
    
    @Column(insertable = false, updatable = false)
    private LocalDateTime createdAt;
    
    @Transient  // Not persisted
    private String calculatedField;
    
    @Formula("(SELECT COUNT(*) FROM orders o WHERE o.product_id = id)")
    private int orderCount;
}
```

### Embedded Types

```java
@Embeddable
public class Address {
    private String street;
    private String city;
    private String state;
    
    @Column(name = "zip_code")
    private String zipCode;
    private String country;
}

@Entity
public class Customer {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String name;
    
    @Embedded
    @AttributeOverrides({
        @AttributeOverride(name = "street", column = @Column(name = "billing_street")),
        @AttributeOverride(name = "city", column = @Column(name = "billing_city")),
        @AttributeOverride(name = "state", column = @Column(name = "billing_state")),
        @AttributeOverride(name = "zipCode", column = @Column(name = "billing_zip")),
        @AttributeOverride(name = "country", column = @Column(name = "billing_country"))
    })
    private Address billingAddress;
    
    @Embedded
    @AttributeOverrides({
        @AttributeOverride(name = "street", column = @Column(name = "shipping_street")),
        @AttributeOverride(name = "city", column = @Column(name = "shipping_city")),
        @AttributeOverride(name = "state", column = @Column(name = "shipping_state")),
        @AttributeOverride(name = "zipCode", column = @Column(name = "shipping_zip")),
        @AttributeOverride(name = "country", column = @Column(name = "shipping_country"))
    })
    private Address shippingAddress;
}
```

### Inheritance Mapping Strategies

#### 1. Single Table Strategy (Default)

```java
@Entity
@Inheritance(strategy = InheritanceType.SINGLE_TABLE)
@DiscriminatorColumn(name = "payment_type", discriminatorType = DiscriminatorType.STRING)
public abstract class Payment {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private BigDecimal amount;
    private LocalDateTime paymentDate;
}

@Entity
@DiscriminatorValue("CREDIT_CARD")
public class CreditCardPayment extends Payment {
    private String cardNumber;
    private String cardHolderName;
    private String expiryDate;
}

@Entity
@DiscriminatorValue("BANK_TRANSFER")
public class BankTransferPayment extends Payment {
    private String bankName;
    private String accountNumber;
    private String routingNumber;
}
```

**Database Table:**
```
payment
├── id (PK)
├── amount
├── payment_date
├── payment_type (Discriminator: CREDIT_CARD, BANK_TRANSFER)
├── card_number (NULL for bank transfers)
├── card_holder_name
├── expiry_date
├── bank_name (NULL for credit cards)
├── account_number
└── routing_number
```

#### 2. Joined Table Strategy

```java
@Entity
@Inheritance(strategy = InheritanceType.JOINED)
public abstract class Vehicle {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    private String manufacturer;
    private String model;
    private int year;
}

@Entity
@PrimaryKeyJoinColumn(name = "vehicle_id")
public class Car extends Vehicle {
    private int numberOfDoors;
    private String fuelType;
}

@Entity
@PrimaryKeyJoinColumn(name = "vehicle_id")
public class Motorcycle extends Vehicle {
    private int engineCC;
    private boolean hasSidecar;
}
```

**Database Tables:**
```
vehicle                    car                        motorcycle
├── id (PK)               ├── vehicle_id (PK, FK)   ├── vehicle_id (PK, FK)
├── manufacturer          ├── number_of_doors       ├── engine_cc
├── model                 └── fuel_type             └── has_sidecar
└── year
```

#### 3. Table Per Class Strategy

```java
@Entity
@Inheritance(strategy = InheritanceType.TABLE_PER_CLASS)
public abstract class Notification {
    @Id
    @GeneratedValue(strategy = GenerationType.TABLE)
    private Long id;
    private String recipient;
    private String message;
    private LocalDateTime sentAt;
}

@Entity
public class EmailNotification extends Notification {
    private String subject;
    private String fromAddress;
}

@Entity
public class SmsNotification extends Notification {
    private String phoneNumber;
    private String provider;
}
```

**Inheritance Strategy Comparison:**

| Strategy | Pros | Cons | Use Case |
|----------|------|------|----------|
| SINGLE_TABLE | Best performance, no joins | Nullable columns, large table | Few subclasses, simple hierarchy |
| JOINED | Normalized, no null columns | Complex queries, joins | Complex hierarchy, data integrity |
| TABLE_PER_CLASS | No joins for concrete class | Poor polymorphic queries | Rare polymorphic queries |

---

## Relationships

### One-to-One

```java
@Entity
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String username;
    
    // Unidirectional
    @OneToOne(cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    @JoinColumn(name = "profile_id", referencedColumnName = "id")
    private UserProfile profile;
}

@Entity
public class UserProfile {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String bio;
    private String avatarUrl;
    
    // Bidirectional (optional - mappedBy indicates non-owning side)
    @OneToOne(mappedBy = "profile")
    private User user;
}
```

**Shared Primary Key (More efficient):**

```java
@Entity
public class User {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    @OneToOne(mappedBy = "user", cascade = CascadeType.ALL, fetch = FetchType.LAZY)
    @PrimaryKeyJoinColumn
    private UserProfile profile;
}

@Entity
public class UserProfile {
    @Id
    private Long id;  // Same as User's ID
    
    @OneToOne
    @MapsId
    @JoinColumn(name = "id")
    private User user;
    
    private String bio;
}
```

### One-to-Many / Many-to-One

```java
@Entity
public class Department {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String name;
    
    // One department has many employees
    @OneToMany(mappedBy = "department", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<Employee> employees = new ArrayList<>();
    
    // Helper methods for bidirectional sync
    public void addEmployee(Employee employee) {
        employees.add(employee);
        employee.setDepartment(this);
    }
    
    public void removeEmployee(Employee employee) {
        employees.remove(employee);
        employee.setDepartment(null);
    }
}

@Entity
public class Employee {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String name;
    
    // Many employees belong to one department (OWNING SIDE - has FK)
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "department_id")
    private Department department;
}
```

### Many-to-Many

```java
@Entity
public class Student {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String name;
    
    @ManyToMany(cascade = {CascadeType.PERSIST, CascadeType.MERGE})
    @JoinTable(
        name = "student_course",
        joinColumns = @JoinColumn(name = "student_id"),
        inverseJoinColumns = @JoinColumn(name = "course_id")
    )
    private Set<Course> courses = new HashSet<>();
    
    public void enrollInCourse(Course course) {
        courses.add(course);
        course.getStudents().add(this);
    }
    
    public void dropCourse(Course course) {
        courses.remove(course);
        course.getStudents().remove(this);
    }
}

@Entity
public class Course {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String title;
    
    @ManyToMany(mappedBy = "courses")
    private Set<Student> students = new HashSet<>();
}
```

**Many-to-Many with Extra Columns (Join Entity):**

```java
@Entity
public class Enrollment {
    @EmbeddedId
    private EnrollmentId id;
    
    @ManyToOne(fetch = FetchType.LAZY)
    @MapsId("studentId")
    private Student student;
    
    @ManyToOne(fetch = FetchType.LAZY)
    @MapsId("courseId")
    private Course course;
    
    private LocalDate enrollmentDate;
    private String grade;
    private boolean completed;
}

@Embeddable
public class EnrollmentId implements Serializable {
    private Long studentId;
    private Long courseId;
    
    // equals() and hashCode() required
}
```

### Relationship Fetch Types

```
┌────────────────────────────────────────────────────────────────────┐
│                    Fetch Type Comparison                            │
├────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  EAGER Loading                    LAZY Loading                      │
│  ─────────────                    ────────────                      │
│                                                                     │
│  Load immediately                 Load on access                    │
│  More initial queries             Fewer initial queries             │
│  Higher memory usage              Lower memory usage                │
│  Risk of N+1 with collections     Risk of LazyInitException         │
│                                                                     │
│  Default for:                     Default for:                      │
│  • @ManyToOne                     • @OneToMany                      │
│  • @OneToOne                      • @ManyToMany                     │
│                                                                     │
└────────────────────────────────────────────────────────────────────┘
```

**Best Practice: Always use LAZY and fetch explicitly when needed**

```java
// Entity definition - LAZY by default
@ManyToOne(fetch = FetchType.LAZY)
@JoinColumn(name = "department_id")
private Department department;

// Repository - Fetch when needed
@Query("SELECT e FROM Employee e JOIN FETCH e.department WHERE e.id = :id")
Optional<Employee> findByIdWithDepartment(@Param("id") Long id);

// Or use EntityGraph
@EntityGraph(attributePaths = {"department", "projects"})
Optional<Employee> findById(Long id);
```

---

## JPQL and Criteria API

### JPQL (Java Persistence Query Language)

```java
@Repository
public class EmployeeRepository {
    
    @PersistenceContext
    private EntityManager em;
    
    // Basic SELECT
    public List<Employee> findAll() {
        return em.createQuery("SELECT e FROM Employee e", Employee.class)
                 .getResultList();
    }
    
    // WHERE clause with parameters
    public List<Employee> findByDepartment(String deptName) {
        return em.createQuery(
            "SELECT e FROM Employee e WHERE e.department.name = :deptName", 
            Employee.class)
            .setParameter("deptName", deptName)
            .getResultList();
    }
    
    // JOIN FETCH to avoid N+1
    public List<Employee> findAllWithDepartments() {
        return em.createQuery(
            "SELECT DISTINCT e FROM Employee e " +
            "JOIN FETCH e.department " +
            "JOIN FETCH e.projects", 
            Employee.class)
            .getResultList();
    }
    
    // Aggregation
    public List<Object[]> getEmployeeCountByDepartment() {
        return em.createQuery(
            "SELECT d.name, COUNT(e) FROM Department d " +
            "LEFT JOIN d.employees e " +
            "GROUP BY d.name " +
            "HAVING COUNT(e) > 5 " +
            "ORDER BY COUNT(e) DESC")
            .getResultList();
    }
    
    // Projection with DTO
    public List<EmployeeDTO> findEmployeeSummaries() {
        return em.createQuery(
            "SELECT new com.example.dto.EmployeeDTO(e.id, e.name, e.department.name) " +
            "FROM Employee e", 
            EmployeeDTO.class)
            .getResultList();
    }
    
    // Subquery
    public List<Employee> findEmployeesAboveAvgSalary() {
        return em.createQuery(
            "SELECT e FROM Employee e " +
            "WHERE e.salary > (SELECT AVG(e2.salary) FROM Employee e2)", 
            Employee.class)
            .getResultList();
    }
    
    // Pagination
    public List<Employee> findPaginated(int page, int size) {
        return em.createQuery("SELECT e FROM Employee e ORDER BY e.name", Employee.class)
                 .setFirstResult(page * size)
                 .setMaxResults(size)
                 .getResultList();
    }
    
    // Named Query (defined on entity)
    public List<Employee> findActiveEmployees() {
        return em.createNamedQuery("Employee.findActive", Employee.class)
                 .getResultList();
    }
    
    // UPDATE query
    public int updateSalaryByDepartment(String dept, BigDecimal increase) {
        return em.createQuery(
            "UPDATE Employee e SET e.salary = e.salary + :increase " +
            "WHERE e.department.name = :dept")
            .setParameter("increase", increase)
            .setParameter("dept", dept)
            .executeUpdate();
    }
    
    // DELETE query
    public int deleteInactiveEmployees() {
        return em.createQuery(
            "DELETE FROM Employee e WHERE e.status = 'INACTIVE' " +
            "AND e.lastLoginDate < :cutoffDate")
            .setParameter("cutoffDate", LocalDate.now().minusYears(1))
            .executeUpdate();
    }
}

// Named Query on Entity
@Entity
@NamedQueries({
    @NamedQuery(name = "Employee.findActive", 
                query = "SELECT e FROM Employee e WHERE e.status = 'ACTIVE'"),
    @NamedQuery(name = "Employee.findByDepartment",
                query = "SELECT e FROM Employee e WHERE e.department.id = :deptId")
})
public class Employee { ... }
```

### Criteria API

```java
@Repository
public class EmployeeCriteriaRepository {
    
    @PersistenceContext
    private EntityManager em;
    
    // Basic query
    public List<Employee> findAll() {
        CriteriaBuilder cb = em.getCriteriaBuilder();
        CriteriaQuery<Employee> cq = cb.createQuery(Employee.class);
        Root<Employee> root = cq.from(Employee.class);
        
        cq.select(root);
        
        return em.createQuery(cq).getResultList();
    }
    
    // Dynamic query with optional filters
    public List<Employee> findByFilters(String name, String department, 
                                         BigDecimal minSalary, BigDecimal maxSalary) {
        CriteriaBuilder cb = em.getCriteriaBuilder();
        CriteriaQuery<Employee> cq = cb.createQuery(Employee.class);
        Root<Employee> root = cq.from(Employee.class);
        
        List<Predicate> predicates = new ArrayList<>();
        
        if (name != null && !name.isEmpty()) {
            predicates.add(cb.like(cb.lower(root.get("name")), 
                          "%" + name.toLowerCase() + "%"));
        }
        
        if (department != null) {
            Join<Employee, Department> deptJoin = root.join("department");
            predicates.add(cb.equal(deptJoin.get("name"), department));
        }
        
        if (minSalary != null) {
            predicates.add(cb.greaterThanOrEqualTo(root.get("salary"), minSalary));
        }
        
        if (maxSalary != null) {
            predicates.add(cb.lessThanOrEqualTo(root.get("salary"), maxSalary));
        }
        
        cq.where(predicates.toArray(new Predicate[0]));
        cq.orderBy(cb.asc(root.get("name")));
        
        return em.createQuery(cq).getResultList();
    }
    
    // Aggregation with grouping
    public List<DepartmentStats> getDepartmentStatistics() {
        CriteriaBuilder cb = em.getCriteriaBuilder();
        CriteriaQuery<DepartmentStats> cq = cb.createQuery(DepartmentStats.class);
        Root<Employee> root = cq.from(Employee.class);
        Join<Employee, Department> dept = root.join("department");
        
        cq.multiselect(
            dept.get("name"),
            cb.count(root),
            cb.avg(root.get("salary")),
            cb.max(root.get("salary")),
            cb.min(root.get("salary"))
        );
        
        cq.groupBy(dept.get("name"));
        cq.having(cb.gt(cb.count(root), 0L));
        cq.orderBy(cb.desc(cb.count(root)));
        
        return em.createQuery(cq).getResultList();
    }
    
    // Subquery example
    public List<Employee> findEmployeesWithHigherThanAvgSalary() {
        CriteriaBuilder cb = em.getCriteriaBuilder();
        CriteriaQuery<Employee> cq = cb.createQuery(Employee.class);
        Root<Employee> root = cq.from(Employee.class);
        
        // Subquery for average salary
        Subquery<Double> subquery = cq.subquery(Double.class);
        Root<Employee> subRoot = subquery.from(Employee.class);
        subquery.select(cb.avg(subRoot.get("salary")));
        
        cq.where(cb.gt(root.get("salary"), subquery));
        
        return em.createQuery(cq).getResultList();
    }
}
```

### Spring Data JPA Query Methods

```java
public interface EmployeeRepository extends JpaRepository<Employee, Long> {
    
    // Derived query methods
    List<Employee> findByName(String name);
    List<Employee> findByNameContainingIgnoreCase(String name);
    List<Employee> findByDepartmentName(String departmentName);
    List<Employee> findByStatusAndDepartmentId(Status status, Long deptId);
    List<Employee> findBySalaryBetween(BigDecimal min, BigDecimal max);
    List<Employee> findByHireDateAfter(LocalDate date);
    List<Employee> findTop5ByOrderBySalaryDesc();
    long countByDepartmentName(String name);
    boolean existsByEmail(String email);
    void deleteByStatus(Status status);
    
    // @Query with JPQL
    @Query("SELECT e FROM Employee e WHERE e.department.name = :dept")
    List<Employee> findEmployeesInDepartment(@Param("dept") String department);
    
    // @Query with native SQL
    @Query(value = "SELECT * FROM employees WHERE YEAR(hire_date) = :year", 
           nativeQuery = true)
    List<Employee> findHiredInYear(@Param("year") int year);
    
    // Modifying queries
    @Modifying
    @Query("UPDATE Employee e SET e.salary = e.salary * :multiplier WHERE e.department.id = :deptId")
    int updateSalaryByDepartment(@Param("deptId") Long deptId, @Param("multiplier") BigDecimal multiplier);
    
    // Projection
    @Query("SELECT e.name as name, e.salary as salary FROM Employee e")
    List<EmployeeProjection> findAllProjected();
    
    // EntityGraph for eager loading
    @EntityGraph(attributePaths = {"department", "projects"})
    List<Employee> findByStatus(Status status);
    
    // Pagination and Sorting
    Page<Employee> findByDepartmentId(Long deptId, Pageable pageable);
    Slice<Employee> findByStatus(Status status, Pageable pageable);
}

// Projection interface
public interface EmployeeProjection {
    String getName();
    BigDecimal getSalary();
}
```

---

## Spring Data JPA Advanced

### Custom Repository Implementation

```java
// 1. Define custom interface
public interface EmployeeRepositoryCustom {
    List<Employee> findByComplexCriteria(EmployeeSearchCriteria criteria);
    void batchUpdate(List<Long> ids, String status);
}

// 2. Implement custom interface
public class EmployeeRepositoryCustomImpl implements EmployeeRepositoryCustom {
    
    @PersistenceContext
    private EntityManager em;
    
    @Override
    public List<Employee> findByComplexCriteria(EmployeeSearchCriteria criteria) {
        CriteriaBuilder cb = em.getCriteriaBuilder();
        CriteriaQuery<Employee> cq = cb.createQuery(Employee.class);
        Root<Employee> root = cq.from(Employee.class);
        
        List<Predicate> predicates = new ArrayList<>();
        
        if (criteria.getName() != null) {
            predicates.add(cb.like(root.get("name"), "%" + criteria.getName() + "%"));
        }
        
        if (criteria.getMinSalary() != null) {
            predicates.add(cb.ge(root.get("salary"), criteria.getMinSalary()));
        }
        
        if (criteria.getDepartmentIds() != null && !criteria.getDepartmentIds().isEmpty()) {
            predicates.add(root.get("department").get("id").in(criteria.getDepartmentIds()));
        }
        
        cq.where(predicates.toArray(new Predicate[0]));
        return em.createQuery(cq).getResultList();
    }
    
    @Override
    @Modifying
    public void batchUpdate(List<Long> ids, String status) {
        em.createQuery("UPDATE Employee e SET e.status = :status WHERE e.id IN :ids")
          .setParameter("status", status)
          .setParameter("ids", ids)
          .executeUpdate();
    }
}

// 3. Extend both interfaces in main repository
public interface EmployeeRepository extends 
        JpaRepository<Employee, Long>, 
        JpaSpecificationExecutor<Employee>,
        EmployeeRepositoryCustom {
    // Derived queries here
}
```

### Specifications for Dynamic Queries

```java
// Specification builder
public class EmployeeSpecifications {
    
    public static Specification<Employee> hasName(String name) {
        return (root, query, cb) -> 
            name == null ? null : cb.like(cb.lower(root.get("name")), "%" + name.toLowerCase() + "%");
    }
    
    public static Specification<Employee> hasSalaryGreaterThan(BigDecimal minSalary) {
        return (root, query, cb) -> 
            minSalary == null ? null : cb.greaterThanOrEqualTo(root.get("salary"), minSalary);
    }
    
    public static Specification<Employee> isInDepartment(Long deptId) {
        return (root, query, cb) -> 
            deptId == null ? null : cb.equal(root.get("department").get("id"), deptId);
    }
    
    public static Specification<Employee> hasStatus(EmployeeStatus status) {
        return (root, query, cb) -> 
            status == null ? null : cb.equal(root.get("status"), status);
    }
    
    // Fetch join for avoiding N+1
    public static Specification<Employee> fetchDepartment() {
        return (root, query, cb) -> {
            if (Long.class != query.getResultType()) {  // Avoid for count queries
                root.fetch("department", JoinType.LEFT);
            }
            return null;
        };
    }
}

// Usage in service
@Service
public class EmployeeService {
    
    @Autowired
    private EmployeeRepository repository;
    
    public Page<Employee> search(EmployeeSearchRequest request, Pageable pageable) {
        Specification<Employee> spec = Specification
            .where(EmployeeSpecifications.hasName(request.getName()))
            .and(EmployeeSpecifications.hasSalaryGreaterThan(request.getMinSalary()))
            .and(EmployeeSpecifications.isInDepartment(request.getDepartmentId()))
            .and(EmployeeSpecifications.hasStatus(request.getStatus()))
            .and(EmployeeSpecifications.fetchDepartment());
        
        return repository.findAll(spec, pageable);
    }
}
```

### QueryDSL Integration

```xml
<!-- Maven dependencies -->
<dependency>
    <groupId>com.querydsl</groupId>
    <artifactId>querydsl-apt</artifactId>
    <scope>provided</scope>
</dependency>
<dependency>
    <groupId>com.querydsl</groupId>
    <artifactId>querydsl-jpa</artifactId>
</dependency>

<!-- APT plugin to generate Q classes -->
<plugin>
    <groupId>com.mysema.maven</groupId>
    <artifactId>apt-maven-plugin</artifactId>
    <executions>
        <execution>
            <goals><goal>process</goal></goals>
            <configuration>
                <outputDirectory>target/generated-sources/java</outputDirectory>
                <processor>com.querydsl.apt.jpa.JPAAnnotationProcessor</processor>
            </configuration>
        </execution>
    </executions>
</plugin>
```

```java
// Generated Q class (automatic)
// QEmployee.java will be generated with type-safe accessors

// Repository with QueryDSL support
public interface EmployeeRepository extends 
        JpaRepository<Employee, Long>,
        QuerydslPredicateExecutor<Employee> {
}

// Type-safe queries
@Service
public class EmployeeQueryService {
    
    @Autowired
    private EmployeeRepository repository;
    
    @PersistenceContext
    private EntityManager em;
    
    // Using QuerydslPredicateExecutor
    public List<Employee> findActiveInDepartment(Long deptId) {
        QEmployee emp = QEmployee.employee;
        
        BooleanExpression predicate = emp.status.eq(EmployeeStatus.ACTIVE)
            .and(emp.department.id.eq(deptId))
            .and(emp.salary.goe(new BigDecimal("50000")));
        
        return (List<Employee>) repository.findAll(predicate);
    }
    
    // Using JPAQueryFactory for complex queries
    public List<EmployeeDTO> findWithProjection(String departmentName) {
        JPAQueryFactory queryFactory = new JPAQueryFactory(em);
        QEmployee emp = QEmployee.employee;
        QDepartment dept = QDepartment.department;
        
        return queryFactory
            .select(Projections.constructor(EmployeeDTO.class,
                emp.id,
                emp.name,
                emp.salary,
                dept.name))
            .from(emp)
            .join(emp.department, dept)
            .where(dept.name.eq(departmentName))
            .orderBy(emp.name.asc())
            .fetch();
    }
    
    // Subqueries
    public List<Employee> findAboveAverageSalary() {
        JPAQueryFactory queryFactory = new JPAQueryFactory(em);
        QEmployee emp = QEmployee.employee;
        QEmployee empSub = new QEmployee("empSub");
        
        return queryFactory
            .selectFrom(emp)
            .where(emp.salary.gt(
                JPAExpressions
                    .select(empSub.salary.avg())
                    .from(empSub)))
            .fetch();
    }
}
```

### Projections Deep Dive

```java
// 1. Interface-based Projection (Recommended for simplicity)
public interface EmployeeSummary {
    Long getId();
    String getName();
    
    @Value("#{target.firstName + ' ' + target.lastName}")
    String getFullName();  // SpEL expression
    
    DepartmentInfo getDepartment();  // Nested projection
    
    interface DepartmentInfo {
        String getName();
        String getLocation();
    }
}

// Repository method
public interface EmployeeRepository extends JpaRepository<Employee, Long> {
    List<EmployeeSummary> findByStatus(EmployeeStatus status);  // Returns projections
    
    <T> List<T> findByDepartmentId(Long deptId, Class<T> type);  // Dynamic projection
}

// Usage
List<EmployeeSummary> summaries = repository.findByStatus(ACTIVE);
List<EmployeeFullDTO> full = repository.findByDepartmentId(1L, EmployeeFullDTO.class);
List<EmployeeSummary> summary = repository.findByDepartmentId(1L, EmployeeSummary.class);

// 2. Class-based Projection (DTO)
@Value  // Lombok immutable
public class EmployeeDTO {
    Long id;
    String name;
    BigDecimal salary;
    String departmentName;
}

@Query("SELECT new com.example.dto.EmployeeDTO(e.id, e.name, e.salary, d.name) " +
       "FROM Employee e JOIN e.department d")
List<EmployeeDTO> findAllAsDTO();

// 3. Tuple Projection
@Query("SELECT e.id, e.name, e.salary FROM Employee e")
List<Tuple> findAllTuples();

public void processTuples() {
    repository.findAllTuples().forEach(tuple -> {
        Long id = tuple.get(0, Long.class);
        String name = tuple.get(1, String.class);
        BigDecimal salary = tuple.get(2, BigDecimal.class);
    });
}
```

### Entity Graphs Advanced

```java
// 1. Named Entity Graph on Entity
@Entity
@NamedEntityGraphs({
    @NamedEntityGraph(
        name = "Employee.withDepartment",
        attributeNodes = @NamedAttributeNode("department")
    ),
    @NamedEntityGraph(
        name = "Employee.withAll",
        attributeNodes = {
            @NamedAttributeNode(value = "department", subgraph = "dept-subgraph"),
            @NamedAttributeNode("projects")
        },
        subgraphs = {
            @NamedSubgraph(
                name = "dept-subgraph",
                attributeNodes = {
                    @NamedAttributeNode("manager"),
                    @NamedAttributeNode("location")
                }
            )
        }
    )
})
public class Employee {
    @ManyToOne(fetch = FetchType.LAZY)
    private Department department;
    
    @ManyToMany(fetch = FetchType.LAZY)
    private Set<Project> projects;
}

// 2. Repository with EntityGraph
public interface EmployeeRepository extends JpaRepository<Employee, Long> {
    
    @EntityGraph(value = "Employee.withDepartment", type = EntityGraphType.FETCH)
    Optional<Employee> findWithDepartmentById(Long id);
    
    @EntityGraph(attributePaths = {"department", "projects"})
    List<Employee> findByStatus(EmployeeStatus status);
}

// 3. Programmatic EntityGraph
@Repository
public class EmployeeCustomRepository {
    
    @PersistenceContext
    private EntityManager em;
    
    public Employee findWithGraph(Long id, String... attributePaths) {
        EntityGraph<Employee> graph = em.createEntityGraph(Employee.class);
        for (String path : attributePaths) {
            graph.addAttributeNodes(path);
        }
        
        Map<String, Object> hints = new HashMap<>();
        hints.put("javax.persistence.fetchgraph", graph);
        
        return em.find(Employee.class, id, hints);
    }
    
    public List<Employee> findAllWithDynamicGraph(List<String> fetchPaths) {
        EntityGraph<Employee> graph = em.createEntityGraph(Employee.class);
        
        for (String path : fetchPaths) {
            if (path.contains(".")) {
                // Handle nested paths like "department.manager"
                String[] parts = path.split("\\.");
                Subgraph<?> subgraph = graph.addSubgraph(parts[0]);
                for (int i = 1; i < parts.length; i++) {
                    if (i == parts.length - 1) {
                        subgraph.addAttributeNodes(parts[i]);
                    } else {
                        subgraph = subgraph.addSubgraph(parts[i]);
                    }
                }
            } else {
                graph.addAttributeNodes(path);
            }
        }
        
        return em.createQuery("SELECT e FROM Employee e", Employee.class)
            .setHint("javax.persistence.fetchgraph", graph)
            .getResultList();
    }
}
```

### Auditing with Spring Data JPA

```java
// 1. Enable auditing
@Configuration
@EnableJpaAuditing(auditorAwareRef = "auditorProvider")
public class JpaAuditingConfig {
    
    @Bean
    public AuditorAware<String> auditorProvider() {
        return () -> Optional.ofNullable(SecurityContextHolder.getContext())
            .map(SecurityContext::getAuthentication)
            .filter(Authentication::isAuthenticated)
            .map(Authentication::getName);
    }
}

// 2. Base auditable entity
@MappedSuperclass
@EntityListeners(AuditingEntityListener.class)
public abstract class AuditableEntity {
    
    @CreatedDate
    @Column(updatable = false)
    private LocalDateTime createdAt;
    
    @LastModifiedDate
    private LocalDateTime updatedAt;
    
    @CreatedBy
    @Column(updatable = false)
    private String createdBy;
    
    @LastModifiedBy
    private String updatedBy;
    
    @Version
    private Long version;
}

// 3. Entity extending auditable
@Entity
public class Employee extends AuditableEntity {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;
    
    private String name;
    // ...
}

// 4. Custom revision entity (for Hibernate Envers integration)
@Entity
@RevisionEntity(CustomRevisionListener.class)
public class CustomRevisionEntity extends DefaultRevisionEntity {
    
    private String username;
    private String ipAddress;
    // ...
}

public class CustomRevisionListener implements RevisionListener {
    @Override
    public void newRevision(Object revisionEntity) {
        CustomRevisionEntity rev = (CustomRevisionEntity) revisionEntity;
        rev.setUsername(getCurrentUser());
        rev.setIpAddress(getClientIpAddress());
    }
}
```

### Events and Callbacks

```java
// 1. JPA Lifecycle Callbacks
@Entity
public class Employee {
    
    @PrePersist
    protected void onCreate() {
        this.createdAt = LocalDateTime.now();
        this.status = EmployeeStatus.PENDING;
    }
    
    @PreUpdate
    protected void onUpdate() {
        this.updatedAt = LocalDateTime.now();
    }
    
    @PostLoad
    protected void onLoad() {
        // Initialize transient fields
        this.fullName = firstName + " " + lastName;
    }
    
    @PreRemove
    protected void onDelete() {
        // Cleanup or validation
        if (this.status == EmployeeStatus.ACTIVE) {
            throw new IllegalStateException("Cannot delete active employee");
        }
    }
}

// 2. External EntityListener
@Component
public class EmployeeEntityListener {
    
    @PostPersist
    public void afterCreate(Employee employee) {
        // Publish event, send notification, etc.
        applicationEventPublisher.publishEvent(
            new EmployeeCreatedEvent(employee));
    }
    
    @PostUpdate
    public void afterUpdate(Employee employee) {
        // Audit logging, cache invalidation, etc.
    }
}

@Entity
@EntityListeners(EmployeeEntityListener.class)
public class Employee { ... }

// 3. Spring Data JPA Domain Events
@Entity
public class Order extends AbstractAggregateRoot<Order> {
    
    @Id
    @GeneratedValue
    private Long id;
    
    public Order complete() {
        this.status = OrderStatus.COMPLETED;
        registerEvent(new OrderCompletedEvent(this));  // Queue event
        return this;
    }
}

// Events published automatically when repository.save() is called
@TransactionalEventListener(phase = TransactionPhase.AFTER_COMMIT)
public void handleOrderCompleted(OrderCompletedEvent event) {
    // Send confirmation email, update inventory, etc.
}
```

### Repository Configuration and Customization

```java
// 1. Custom base repository
@NoRepositoryBean
public interface BaseRepository<T, ID> extends JpaRepository<T, ID> {
    
    Optional<T> findByIdAndNotDeleted(ID id);
    
    @Query("SELECT e FROM #{#entityName} e WHERE e.deleted = false")
    List<T> findAllActive();
    
    @Modifying
    @Query("UPDATE #{#entityName} e SET e.deleted = true WHERE e.id = :id")
    void softDelete(@Param("id") ID id);
}

// Implementation
public class BaseRepositoryImpl<T, ID extends Serializable> 
        extends SimpleJpaRepository<T, ID> 
        implements BaseRepository<T, ID> {
    
    private final EntityManager em;
    private final JpaEntityInformation<T, ?> entityInformation;
    
    public BaseRepositoryImpl(JpaEntityInformation<T, ?> entityInformation,
                               EntityManager entityManager) {
        super(entityInformation, entityManager);
        this.em = entityManager;
        this.entityInformation = entityInformation;
    }
    
    @Override
    public Optional<T> findByIdAndNotDeleted(ID id) {
        // Custom implementation
    }
}

// Configure Spring Data to use custom base
@EnableJpaRepositories(
    basePackages = "com.example.repository",
    repositoryBaseClass = BaseRepositoryImpl.class
)
public class JpaConfig { }

// 2. Repository Fragments (Composition)
// Fragment interface
public interface EmployeeFragmentRepository {
    List<Employee> customFind();
}

// Fragment implementation
public class EmployeeFragmentRepositoryImpl implements EmployeeFragmentRepository {
    @Override
    public List<Employee> customFind() { ... }
}

// Compose in main repository
public interface EmployeeRepository extends 
        JpaRepository<Employee, Long>,
        EmployeeFragmentRepository,
        AnotherFragment {
}
```

### Pagination and Sorting Advanced

```java
// 1. Slice vs Page
// Page: includes total count (extra query)
Page<Employee> page = repository.findAll(PageRequest.of(0, 20));
page.getTotalElements();  // SELECT COUNT(*)
page.getTotalPages();

// Slice: no total count (more efficient for infinite scroll)
Slice<Employee> slice = repository.findByStatus(status, PageRequest.of(0, 20));
slice.hasNext();  // Based on fetching size+1 elements

// 2. Keyset / Cursor Pagination (more efficient for large offsets)
public interface EmployeeRepository extends JpaRepository<Employee, Long> {
    
    @Query("SELECT e FROM Employee e WHERE e.id > :lastId ORDER BY e.id")
    List<Employee> findNextPage(@Param("lastId") Long lastId, Pageable pageable);
    
    // With multiple sort columns
    @Query("SELECT e FROM Employee e WHERE " +
           "(e.lastName > :lastName) OR " +
           "(e.lastName = :lastName AND e.id > :id) " +
           "ORDER BY e.lastName, e.id")
    List<Employee> findNextPageByCursor(
        @Param("lastName") String lastName,
        @Param("id") Long id,
        Pageable pageable);
}

// 3. Dynamic sorting
public Page<Employee> search(EmployeeSearchRequest request) {
    Sort sort = Sort.by(
        Sort.Order.asc("lastName").ignoreCase(),
        Sort.Order.desc("salary").nullsLast()
    );
    
    // Or from request
    Sort dynamicSort = Sort.by(
        request.getSortFields().stream()
            .map(f -> f.isAsc() ? Sort.Order.asc(f.getName()) : Sort.Order.desc(f.getName()))
            .collect(Collectors.toList())
    );
    
    Pageable pageable = PageRequest.of(request.getPage(), request.getSize(), dynamicSort);
    return repository.findAll(spec, pageable);
}

// 4. Window functions for ranking
@Query(value = "SELECT e.*, " +
               "RANK() OVER (PARTITION BY e.department_id ORDER BY e.salary DESC) as salary_rank " +
               "FROM employees e WHERE e.department_id IN :deptIds", 
       nativeQuery = true)
List<Object[]> findWithSalaryRank(@Param("deptIds") List<Long> deptIds);
```

---

## Transactions

### Transaction Configuration

```java
@Configuration
@EnableTransactionManagement
public class TransactionConfig {
    
    @Bean
    public PlatformTransactionManager transactionManager(EntityManagerFactory emf) {
        return new JpaTransactionManager(emf);
    }
}
```

### @Transactional Annotation

```java
@Service
public class OrderService {
    
    @Autowired
    private OrderRepository orderRepository;
    
    @Autowired
    private InventoryService inventoryService;
    
    @Autowired
    private PaymentService paymentService;
    
    // Basic transaction
    @Transactional
    public Order createOrder(OrderRequest request) {
        Order order = new Order();
        order.setCustomerId(request.getCustomerId());
        order.setItems(request.getItems());
        
        // All operations in single transaction
        inventoryService.reserveItems(request.getItems());
        paymentService.processPayment(request.getPaymentInfo());
        
        return orderRepository.save(order);
    }
    
    // Read-only transaction (optimized)
    @Transactional(readOnly = true)
    public List<Order> getOrderHistory(Long customerId) {
        return orderRepository.findByCustomerId(customerId);
    }
    
    // Transaction with specific propagation
    @Transactional(propagation = Propagation.REQUIRES_NEW)
    public void logOrderEvent(Long orderId, String event) {
        // Always runs in new transaction
    }
    
    // Transaction with rollback rules
    @Transactional(
        rollbackFor = {PaymentException.class, InventoryException.class},
        noRollbackFor = {NotificationException.class}
    )
    public Order processOrder(OrderRequest request) {
        // ...
    }
    
    // Transaction with timeout and isolation
    @Transactional(timeout = 30, isolation = Isolation.READ_COMMITTED)
    public void longRunningProcess() {
        // ...
    }
}
```

### Propagation Types

```
┌────────────────────────────────────────────────────────────────────────────┐
│                        Transaction Propagation                              │
├────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│  REQUIRED (Default)                                                         │
│  ┌─────────────────────────────────────────────────────────┐               │
│  │ Outer TX  ──────────────────────────────────────────────│               │
│  │           │ Inner (joins existing)                      │               │
│  └───────────┴─────────────────────────────────────────────┘               │
│  Creates new TX if none exists, otherwise joins existing                   │
│                                                                             │
│  REQUIRES_NEW                                                               │
│  ┌─────────────────────────────────────────────────────────┐               │
│  │ Outer TX (suspended)                                    │               │
│  │           ┌──────────────────────┐                      │               │
│  │           │ Inner TX (new)       │                      │               │
│  │           └──────────────────────┘                      │               │
│  └─────────────────────────────────────────────────────────┘               │
│  Always creates new TX, suspends existing                                  │
│                                                                             │
│  MANDATORY                                                                  │
│  Must have existing TX, throws exception if none                           │
│                                                                             │
│  SUPPORTS                                                                   │
│  Uses TX if exists, runs non-transactional if not                          │
│                                                                             │
│  NOT_SUPPORTED                                                              │
│  Suspends existing TX, runs non-transactional                              │
│                                                                             │
│  NEVER                                                                      │
│  Throws exception if TX exists                                              │
│                                                                             │
│  NESTED                                                                     │
│  Creates savepoint within existing TX (requires JDBC 3.0)                  │
│                                                                             │
└────────────────────────────────────────────────────────────────────────────┘
```

### Isolation Levels

```java
@Transactional(isolation = Isolation.READ_COMMITTED)
public void process() { }
```

| Isolation Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|-----------------|------------|---------------------|--------------|
| READ_UNCOMMITTED | Yes | Yes | Yes |
| READ_COMMITTED | No | Yes | Yes |
| REPEATABLE_READ | No | No | Yes |
| SERIALIZABLE | No | No | No |

### Programmatic Transactions

```java
@Service
public class OrderService {
    
    @Autowired
    private TransactionTemplate transactionTemplate;
    
    @Autowired
    private PlatformTransactionManager transactionManager;
    
    // Using TransactionTemplate
    public Order createOrderProgrammatic(OrderRequest request) {
        return transactionTemplate.execute(status -> {
            try {
                Order order = new Order();
                // ... process order
                return orderRepository.save(order);
            } catch (Exception e) {
                status.setRollbackOnly();
                throw e;
            }
        });
    }
    
    // Using TransactionManager directly
    public void complexTransaction() {
        TransactionDefinition def = new DefaultTransactionDefinition();
        TransactionStatus status = transactionManager.getTransaction(def);
        
        try {
            // ... operations
            transactionManager.commit(status);
        } catch (Exception e) {
            transactionManager.rollback(status);
            throw e;
        }
    }
}
```

---

## Caching

### First-Level Cache (Session Cache)

The first-level cache is the persistence context (session). It's enabled by default and cannot be disabled.

```java
@Transactional
public void demonstrateFirstLevelCache() {
    // First query - hits database
    Employee emp1 = entityManager.find(Employee.class, 1L);
    
    // Second query - returns from cache (same object)
    Employee emp2 = entityManager.find(Employee.class, 1L);
    
    System.out.println(emp1 == emp2);  // true - same object from cache
    
    // Clear cache
    entityManager.clear();
    
    // Now hits database again
    Employee emp3 = entityManager.find(Employee.class, 1L);
    System.out.println(emp1 == emp3);  // false - different objects
}
```

### Second-Level Cache

```java
// 1. Add dependency
// implementation 'org.hibernate:hibernate-ehcache:5.6.x'
// or
// implementation 'org.hibernate:hibernate-jcache:5.6.x'

// 2. Configure in application.properties
spring.jpa.properties.hibernate.cache.use_second_level_cache=true
spring.jpa.properties.hibernate.cache.region.factory_class=org.hibernate.cache.jcache.JCacheRegionFactory
spring.jpa.properties.hibernate.javax.cache.provider=org.ehcache.jsr107.EhcacheCachingProvider
spring.jpa.properties.hibernate.cache.use_query_cache=true

// 3. Enable caching on entity
@Entity
@Cacheable
@org.hibernate.annotations.Cache(usage = CacheConcurrencyStrategy.READ_WRITE)
public class Employee {
    @Id
    private Long id;
    private String name;
    
    @OneToMany(mappedBy = "employee")
    @org.hibernate.annotations.Cache(usage = CacheConcurrencyStrategy.READ_WRITE)
    private List<Project> projects;
}
```

**Cache Concurrency Strategies:**

| Strategy | Description | Use Case |
|----------|-------------|----------|
| READ_ONLY | Never modified after creation | Reference data |
| NONSTRICT_READ_WRITE | Occasional updates, no strict consistency | Low-contention data |
| READ_WRITE | Strict consistency with soft locks | Frequently read, occasionally updated |
| TRANSACTIONAL | Full transactional cache (JTA required) | Critical data integrity |

### Query Cache

```java
// Enable query cache
@Query("SELECT e FROM Employee e WHERE e.department.id = :deptId")
@QueryHints(@QueryHint(name = "org.hibernate.cacheable", value = "true"))
List<Employee> findByDepartmentIdCached(@Param("deptId") Long deptId);

// Or programmatically
List<Employee> employees = entityManager
    .createQuery("SELECT e FROM Employee e WHERE e.status = :status")
    .setParameter("status", "ACTIVE")
    .setHint("org.hibernate.cacheable", true)
    .getResultList();
```

### Cache Architecture

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           Hibernate Caching                                  │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                     First-Level Cache (Session)                        │ │
│  │  • Per session/transaction                                             │ │
│  │  • Automatic, cannot disable                                           │ │
│  │  • Stores entities within transaction                                  │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                    │                                         │
│                                    ▼                                         │
│  ┌────────────────────────────────────────────────────────────────────────┐ │
│  │                    Second-Level Cache (SessionFactory)                 │ │
│  │  ┌──────────────────────┐  ┌──────────────────────┐                    │ │
│  │  │    Entity Cache      │  │   Collection Cache   │                    │ │
│  │  │  Employee: {1->data} │  │  Employee.projects   │                    │ │
│  │  │  Employee: {2->data} │  │     {1->[1,2,3]}     │                    │ │
│  │  └──────────────────────┘  └──────────────────────┘                    │ │
│  │                                                                         │ │
│  │  ┌──────────────────────┐  ┌──────────────────────┐                    │ │
│  │  │    Query Cache       │  │   Update Timestamps  │                    │ │
│  │  │ "SELECT e WHERE..."  │  │   Employee: 12345    │                    │ │
│  │  │    -> [1, 2, 3]      │  │   Project: 12340     │                    │ │
│  │  └──────────────────────┘  └──────────────────────┘                    │ │
│  └────────────────────────────────────────────────────────────────────────┘ │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## Performance Optimization

### N+1 Problem

```java
// Problem: This generates N+1 queries
@Transactional(readOnly = true)
public void nPlusOneProblem() {
    List<Department> departments = departmentRepository.findAll();  // 1 query
    
    for (Department dept : departments) {
        // N additional queries - one for each department's employees
        System.out.println(dept.getName() + ": " + dept.getEmployees().size());
    }
}

// Solution 1: JOIN FETCH
@Query("SELECT d FROM Department d JOIN FETCH d.employees")
List<Department> findAllWithEmployees();

// Solution 2: EntityGraph
@EntityGraph(attributePaths = {"employees"})
List<Department> findAll();

// Solution 3: Batch fetching
@Entity
public class Department {
    @OneToMany(mappedBy = "department")
    @BatchSize(size = 20)  // Fetch employees in batches of 20
    private List<Employee> employees;
}

// Solution 4: Subselect fetching
@Entity
public class Department {
    @OneToMany(mappedBy = "department")
    @Fetch(FetchMode.SUBSELECT)  // Single subselect for all employees
    private List<Employee> employees;
}
```

### Batch Processing

```java
@Service
public class BatchService {
    
    @PersistenceContext
    private EntityManager em;
    
    // Batch insert with periodic flush and clear
    @Transactional
    public void batchInsert(List<Employee> employees) {
        int batchSize = 50;
        
        for (int i = 0; i < employees.size(); i++) {
            em.persist(employees.get(i));
            
            if (i > 0 && i % batchSize == 0) {
                em.flush();
                em.clear();
            }
        }
        
        em.flush();
        em.clear();
    }
    
    // Stateless session for pure batch operations (Hibernate-specific)
    public void batchInsertStateless(List<Employee> employees) {
        StatelessSession session = sessionFactory.openStatelessSession();
        Transaction tx = session.beginTransaction();
        
        try {
            for (Employee emp : employees) {
                session.insert(emp);
            }
            tx.commit();
        } catch (Exception e) {
            tx.rollback();
            throw e;
        } finally {
            session.close();
        }
    }
}

// Configuration for batching
spring.jpa.properties.hibernate.jdbc.batch_size=50
spring.jpa.properties.hibernate.order_inserts=true
spring.jpa.properties.hibernate.order_updates=true
spring.jpa.properties.hibernate.batch_versioned_data=true
```

### Projection for Read Operations

```java
// Instead of loading full entity
@Query("SELECT new com.example.dto.EmployeeSummary(e.id, e.name, e.department.name) " +
       "FROM Employee e WHERE e.status = 'ACTIVE'")
List<EmployeeSummary> findActiveEmployeeSummaries();

// Interface projection
public interface EmployeeSummary {
    Long getId();
    String getName();
    String getDepartmentName();
}

@Query("SELECT e.id as id, e.name as name, d.name as departmentName " +
       "FROM Employee e JOIN e.department d")
List<EmployeeSummary> findAllSummaries();

// Tuple projection for dynamic columns
@Query("SELECT e.id, e.name, e.salary FROM Employee e")
List<Tuple> findEmployeeTuples();
```

### Read-Only Transactions

```java
@Service
public class ReportService {
    
    // Read-only optimization - no dirty checking
    @Transactional(readOnly = true)
    public List<EmployeeReport> generateReport() {
        return employeeRepository.findAllForReport();
    }
}
```

### Pagination for Large Results

```java
// Always paginate large result sets
@Query("SELECT e FROM Employee e")
Page<Employee> findAllPaginated(Pageable pageable);

// Keyset pagination (more efficient for large offsets)
@Query("SELECT e FROM Employee e WHERE e.id > :lastId ORDER BY e.id")
List<Employee> findNextPage(@Param("lastId") Long lastId, Pageable pageable);

// Stream for processing large datasets
@Query("SELECT e FROM Employee e")
@QueryHints(@QueryHint(name = HINT_FETCH_SIZE, value = "50"))
Stream<Employee> streamAllEmployees();

@Transactional(readOnly = true)
public void processAllEmployees() {
    try (Stream<Employee> stream = employeeRepository.streamAllEmployees()) {
        stream.forEach(this::processEmployee);
    }
}
```

### Index Optimization

```java
@Entity
@Table(name = "employees", indexes = {
    @Index(name = "idx_emp_email", columnList = "email", unique = true),
    @Index(name = "idx_emp_dept_status", columnList = "department_id, status"),
    @Index(name = "idx_emp_name", columnList = "last_name, first_name")
})
public class Employee {
    // ...
}
```

---

## Best Practices

### 1. Entity Design

```java
// DO: Use wrapper types for nullable columns
@Column(nullable = true)
private Integer age;  // Can be null

// DON'T: Use primitives for nullable columns
private int age;  // Can't represent null, defaults to 0

// DO: Implement equals/hashCode properly
@Entity
public class Employee {
    @Id
    @GeneratedValue
    private Long id;
    
    @NaturalId
    private String employeeNumber;  // Business key
    
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (!(o instanceof Employee)) return false;
        Employee employee = (Employee) o;
        return Objects.equals(employeeNumber, employee.employeeNumber);
    }
    
    @Override
    public int hashCode() {
        return Objects.hash(employeeNumber);
    }
}

// DO: Use Set for ManyToMany (better performance than List)
@ManyToMany
private Set<Role> roles = new HashSet<>();

// DO: Initialize collections
@OneToMany(mappedBy = "department")
private List<Employee> employees = new ArrayList<>();
```

### 2. Relationship Management

```java
// DO: Implement helper methods for bidirectional relationships
@Entity
public class Department {
    
    @OneToMany(mappedBy = "department", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<Employee> employees = new ArrayList<>();
    
    public void addEmployee(Employee employee) {
        employees.add(employee);
        employee.setDepartment(this);
    }
    
    public void removeEmployee(Employee employee) {
        employees.remove(employee);
        employee.setDepartment(null);
    }
}

// DO: Use FetchType.LAZY and fetch explicitly
@ManyToOne(fetch = FetchType.LAZY)
private Department department;

// DON'T: Use CascadeType.ALL on ManyToMany
@ManyToMany(cascade = CascadeType.ALL)  // Dangerous!
private Set<Tag> tags;

// DO: Be selective with cascade types
@ManyToMany(cascade = {CascadeType.PERSIST, CascadeType.MERGE})
private Set<Tag> tags;
```

### 3. Query Optimization

```java
// DO: Use DTO projections for read-only queries
@Query("SELECT new com.example.dto.EmployeeDTO(e.id, e.name) FROM Employee e")
List<EmployeeDTO> findAllDTO();

// DO: Use pagination
Page<Employee> findAll(Pageable pageable);

// DO: Use specific fetching strategies
@EntityGraph(attributePaths = {"department"})
Optional<Employee> findById(Long id);

// DON'T: Fetch all data when you only need specific fields
List<Employee> findAll();  // Loads everything

// DO: Use EXISTS instead of COUNT when checking existence
boolean existsByEmail(String email);
```

### 4. Transaction Management

```java
// DO: Keep transactions short
@Transactional
public Order createOrder(OrderDTO dto) {
    Order order = mapper.toEntity(dto);
    return orderRepository.save(order);
}

// DON'T: Do external calls within transaction
@Transactional
public Order createOrder(OrderDTO dto) {
    Order order = orderRepository.save(new Order(dto));
    emailService.sendEmail(order);  // External call - bad!
    return order;
}

// DO: Separate transaction from external calls
public Order createOrderAndNotify(OrderDTO dto) {
    Order order = createOrder(dto);  // Transaction
    emailService.sendEmail(order);   // Outside transaction
    return order;
}

// DO: Use readOnly for queries
@Transactional(readOnly = true)
public List<Order> getOrders() {
    return orderRepository.findAll();
}
```

### 5. Common Pitfalls to Avoid

```java
// PITFALL 1: LazyInitializationException
// ❌ Wrong - session closed
public Employee getEmployee(Long id) {
    return employeeRepository.findById(id).orElse(null);
}

public void useEmployee() {
    Employee emp = getEmployee(1L);
    emp.getDepartment().getName();  // LazyInitializationException!
}

// ✅ Correct - fetch eagerly or keep session open
@Transactional(readOnly = true)
public Employee getEmployeeWithDepartment(Long id) {
    return employeeRepository.findByIdWithDepartment(id).orElse(null);
}

// PITFALL 2: Detached entity passed to persist
// ❌ Wrong
@Transactional
public void updateEmployee(Employee employee) {
    entityManager.persist(employee);  // Error if entity was detached
}

// ✅ Correct
@Transactional
public Employee updateEmployee(Employee employee) {
    return entityManager.merge(employee);
}

// PITFALL 3: Losing updates with merge
// ❌ Wrong - partial update loses data
public void updateName(Long id, String name) {
    Employee emp = new Employee();
    emp.setId(id);
    emp.setName(name);
    employeeRepository.save(emp);  // Other fields become null!
}

// ✅ Correct - fetch first, then update
@Transactional
public void updateName(Long id, String name) {
    Employee emp = employeeRepository.findById(id).orElseThrow();
    emp.setName(name);
    // Auto-saved due to dirty checking
}
```

---

## Interview Questions

### Basic Level

**Q1: What is the difference between JPA and Hibernate?**
> JPA is a specification (interface) that defines the standard for ORM in Java. Hibernate is the most popular implementation of this specification. JPA provides portability across implementations, while Hibernate offers additional features beyond the standard.

**Q2: Explain the entity lifecycle states.**
> - **New/Transient**: Entity created but not managed
> - **Managed**: Entity in persistence context, changes tracked
> - **Detached**: Was managed, persistence context closed
> - **Removed**: Marked for deletion

**Q3: What is the difference between `persist()` and `merge()`?**
> - `persist()`: Makes a new transient entity managed. Throws error if entity already has ID.
> - `merge()`: Copies state of detached entity to a managed entity. Returns the managed instance.

**Q4: What is lazy loading?**
> Lazy loading delays the loading of an association until it's accessed. It's the default for `@OneToMany` and `@ManyToMany` relationships. Improves performance by not loading unnecessary data.

**Q5: What is the N+1 problem?**
> When fetching N entities that each have a lazy-loaded collection, accessing those collections results in N additional queries (1 for parent + N for children). Solution: Use JOIN FETCH, EntityGraph, or batch fetching.

### Intermediate Level

**Q6: Explain different ID generation strategies.**
```java
// IDENTITY - Database auto-increment (breaks batching)
// SEQUENCE - Database sequence (best for batching)
// TABLE - Simulates sequence with table
// UUID - For distributed systems
```

**Q7: Compare inheritance mapping strategies.**
> - **SINGLE_TABLE**: Best performance, nullable columns
> - **JOINED**: Normalized, complex queries with joins
> - **TABLE_PER_CLASS**: No joins for concrete, poor polymorphic queries

**Q8: What is the persistence context?**
> The persistence context is a set of managed entities. It acts as a first-level cache and ensures entity identity (same ID = same object). It tracks changes for automatic dirty checking.

**Q9: Explain optimistic vs pessimistic locking.**
```java
// Optimistic - uses version column
@Version
private Long version;

// Pessimistic - database locks
entityManager.find(Employee.class, id, LockModeType.PESSIMISTIC_WRITE);
```

**Q10: What is dirty checking?**
> Hibernate automatically tracks changes to managed entities. During flush, it compares current state with snapshot taken when entity was loaded and generates UPDATE statements for modified entities.

### Advanced Level

**Q11: How do you handle batch processing efficiently?**
```java
int batchSize = 50;
for (int i = 0; i < entities.size(); i++) {
    entityManager.persist(entities.get(i));
    if (i % batchSize == 0) {
        entityManager.flush();
        entityManager.clear();
    }
}
```

**Q12: Explain second-level cache architecture.**
> The second-level cache is shared across sessions. It caches:
> - Entities (by ID)
> - Collections (by owner ID)
> - Query results (by query + parameters)
> 
> Uses cache concurrency strategies: READ_ONLY, NONSTRICT_READ_WRITE, READ_WRITE, TRANSACTIONAL.

**Q13: How do you optimize read-only operations?**
```java
@Transactional(readOnly = true)  // Skips dirty checking
public List<Employee> getAll() {
    return repository.findAll();
}

// Use projections
@Query("SELECT new DTO(e.id, e.name) FROM Employee e")
List<DTO> findAllDTO();
```

**Q14: What is the difference between `@JoinColumn` and `@JoinTable`?**
> - `@JoinColumn`: Foreign key in owning entity's table (One-to-One, Many-to-One)
> - `@JoinTable`: Separate join table for the relationship (Many-to-Many, optionally One-to-Many)

**Q15: How does `orphanRemoval` differ from `CascadeType.REMOVE`?**
> - `CascadeType.REMOVE`: Deletes child when parent is deleted
> - `orphanRemoval=true`: Also deletes child when removed from parent's collection

```java
@OneToMany(mappedBy = "parent", cascade = CascadeType.ALL, orphanRemoval = true)
private List<Child> children;

// This triggers removal:
parent.getChildren().remove(child);
```

**Q16: Explain transaction propagation types with examples.**
```java
// REQUIRED (default): Join existing or create new
// REQUIRES_NEW: Always create new, suspend existing
// MANDATORY: Must have existing transaction
// SUPPORTS: Use if available, else non-transactional
// NOT_SUPPORTED: Suspend existing, run non-transactional
// NEVER: Error if transaction exists
// NESTED: Create savepoint within existing
```

**Q17: How do you handle LazyInitializationException?**
> 1. Join fetch in query
> 2. Use `@EntityGraph`
> 3. Open Session in View pattern (not recommended for APIs)
> 4. Initialize within transaction
> 5. Use DTO projections

**Q18: What is the difference between `@Embeddable` and `@Entity`?**
> - `@Entity`: Has its own table, lifecycle, and identity
> - `@Embeddable`: Value object, stored in parent's table, no identity, lifecycle tied to parent

**Q19: Explain `@MapsId` and when to use it.**
```java
// For shared primary key in One-to-One
@Entity
public class UserProfile {
    @Id
    private Long id;  // Same as User's ID
    
    @OneToOne
    @MapsId  // User's ID becomes this entity's ID
    private User user;
}
```

**Q20: How do you implement soft deletes?**
```java
@Entity
@SQLDelete(sql = "UPDATE employee SET deleted = true WHERE id = ?")
@Where(clause = "deleted = false")
public class Employee {
    private boolean deleted = false;
}
```

---

## Quick Reference

### Common Annotations

| Annotation | Purpose |
|------------|---------|
| `@Entity` | Mark class as JPA entity |
| `@Table` | Specify table details |
| `@Id` | Primary key field |
| `@GeneratedValue` | Auto-generate ID |
| `@Column` | Column mapping details |
| `@Transient` | Exclude field from persistence |
| `@Temporal` | Date/time mapping |
| `@Enumerated` | Enum mapping |
| `@Lob` | Large object (BLOB/CLOB) |
| `@Embedded` / `@Embeddable` | Composite value |
| `@OneToOne` | One-to-one relationship |
| `@OneToMany` | One-to-many relationship |
| `@ManyToOne` | Many-to-one relationship |
| `@ManyToMany` | Many-to-many relationship |
| `@JoinColumn` | Foreign key column |
| `@JoinTable` | Join table for relationship |
| `@Version` | Optimistic locking |
| `@Cacheable` | Enable second-level cache |

### Configuration Properties

```properties
# Basic Configuration
spring.jpa.show-sql=true
spring.jpa.properties.hibernate.format_sql=true

# DDL Generation
spring.jpa.hibernate.ddl-auto=update  # none, validate, update, create, create-drop

# Performance
spring.jpa.properties.hibernate.jdbc.batch_size=50
spring.jpa.properties.hibernate.order_inserts=true
spring.jpa.properties.hibernate.order_updates=true

# Caching
spring.jpa.properties.hibernate.cache.use_second_level_cache=true
spring.jpa.properties.hibernate.cache.use_query_cache=true

# Statistics (for debugging)
spring.jpa.properties.hibernate.generate_statistics=true
```

---

*Last Updated: February 2026*
