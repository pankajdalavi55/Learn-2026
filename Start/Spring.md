
Spring framework:

Spring is an open-source, lightweight, enterprise-grade Java application framework that provides comprehensive infrastructure support. It was created by Rod Johnson in 2003 as a response to the complexity of J2EE (now Jakarta EE).

Core Philosophy:

- Favour plain Java objects (POJOs) over heavy enterprise components
- Promote loose coupling through IoC/DI
- Provide declarative programming via AOP
- Reduce boilerplate code
- Embrace convention over configuration
- Support non-invasive programming (your code doesn't depend on Spring APIs)


IoC & DI:

What is Inversion of Control?
IoC is a design principle (not a pattern) where the control of object creation and lifecycle is transferred from the application code to a framework or container. The "inversion" refers to reversing the traditional flow where objects directly instantiate their dependencies.

traditional
Application Code → Creates Dependencies → Uses Dependencies

Inverted Control flow:
Container Creates Dependencies → Injects into Application Code → Code Uses Dependencies

IoC vs DI — Clearing the Confusion
Concept	Definition	Relationship
IoC	                Design principle	        The broader concept
DI	                Implementation technique	One way to achieve IoC
Service Locator	    Alternative technique	    Another way to achieve IoC

IoC can be achieved through:
Dependency Injection (Spring's approach) — Dependencies pushed to the object
Service Locator — Object pulls dependencies from a registry
Factory Pattern — Factory creates and returns dependencies
Template Method — Superclass controls algorithm, subclass provides specifics


# The Spring IoC Container
The container is responsible for:

Instantiating beans
Configuring them (injecting dependencies)
Managing their lifecycle (init → use → destroy)


wo main container implementations:

Container	           Interface	        Use Case
BeanFactory	           BeanFactory	        Lightweight, lazy init, low memory
ApplicationContext	   ApplicationContext	Full-featured: events, i18n, AOP, eager init

In practice, always use ApplicationContext; BeanFactory is legacy.

// Bootstrapping ApplicationContext
// 1. Annotation-based
ApplicationContext ctx = new AnnotationConfigApplicationContext(AppConfig.class);

// 2. XML-based
ApplicationContext ctx = new ClassPathXmlApplicationContext("applicationContext.xml");

// 3. Web-based
// Configured via web.xml ContextLoaderListener or WebApplicationInitializer


Types of DI
1. Constuctor Injection :

@Component
public class NotificationService {
    private final EmailSender emailSender;
    private final SmsSender smsSender;

    @Autowired // optional on single constructor since Spring 4.3
    public NotificationService(EmailSender emailSender, SmsSender smsSender) {
        this.emailSender = emailSender;
        this.smsSender = smsSender;
    }
}

Why prefer constructor injection?

Ensures immutability (final fields)
Makes dependencies explicit — visible in constructor signature
Object is never in an incomplete state
Easier to unit test (just call new with mocks)

2. Setter injection
@Component
public class ReportGenerator {
    private DataSource dataSource;

    @Autowired
    public void setDataSource(DataSource dataSource) {
        this.dataSource = dataSource;
    }
}
Use setter injection for optional dependencies or when you need to allow reconfiguration after construction.

3. Field Injection:
@Component
public class UserService {
    @Autowired
    private UserRepository userRepository; // no setter, no constructor param
}

Why avoid?
Cannot make fields final
Hides dependencies
Impossible to instantiate without reflection (hard to unit test)
Violates single-responsibility detection (class can silently accumulate many injected fields)



@Autowired Resolution
Spring resolves @Autowired in this order:

1. Match by TYPE  →  Only one bean of requested type?  ✅ Inject it.
2. Multiple candidates?
   a. @Primary bean present?   ✅ Use it.
   b. @Qualifier specified?    ✅ Match by qualifier name.
   c. Match by FIELD/PARAM NAME as bean name?  ✅ Use it.
   d. None matched?           ❌ Throw NoUniqueBeanDefinitionException.




Bean Lifecycyle & Scope

Why Bean Lifecycle Matters
Understanding bean lifecycle is crucial for:

1. Resource Management — Acquire resources (DB connections, file handles, thread pools) at the right time and release them properly
2. Initialization Order — Ensure dependencies are ready before a bean is used
3. Framework Integration — Hook into Spring's infrastructure (AOP proxies, transaction management)
4. Debugging — Diagnose issues like "bean not found" or "bean not fully initialized"
5. Performance Optimization — Defer expensive initialization until necessary


The Container's Responsibilities
Spring's IoC container acts as a sophisticated object factory with these responsibilities:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Spring IoC Container                         │
├─────────────────────────────────────────────────────────────────┤
│  1. Read Configuration    │  XML, annotations, Java config      │
│  2. Create Bean Graph     │  Resolve dependencies, detect cycles│
│  3. Instantiate Beans     │  Call constructors in correct order │
│  4. Inject Dependencies   │  Wire beans together                │
│  5. Apply Post-Processing │  AOP proxies, validation            │
│  6. Initialize            │  Call @PostConstruct, init methods  │
│  7. Make Available        │  Beans ready for use                │
│  8. Destroy               │  Clean shutdown, release resources  │
└─────────────────────────────────────────────────────────────────┘
```

# Bean Scopes

| Scope | Annotation / XML | Behaviour |
|---|---|---|
| **singleton** (default) | `@Scope("singleton")` | One instance per Spring container |
| **prototype** | `@Scope("prototype")` | New instance every time requested |
| **request** | `@RequestScope` | One per HTTP request |
| **session** | `@SessionScope` | One per HTTP session |
| **application** | `@ApplicationScope` | One per ServletContext |
| **websocket** | `@Scope("websocket")` | One per WebSocket session |

### 3.6 The Singleton-Prototype Injection Problem

```java
@Component // singleton by default
public class CheckoutService {
    @Autowired
    private ShoppingCart cart; // prototype

    // ⚠ BUG: same cart instance every time — prototype scope "lost"
}
```

**Solutions:**

**1. Method Injection with `@Lookup` (Preferred):**
```java
@Component
public abstract class CheckoutService {

    public void checkout() {
        ShoppingCart cart = getCart(); // fresh prototype each time
        cart.addItem(...);
    }

    @Lookup
    protected abstract ShoppingCart getCart(); // Spring overrides at runtime
}
```

**2. `ObjectProvider`:**
```java
@Component
public class CheckoutService {
    private final ObjectProvider<ShoppingCart> cartProvider;

    public CheckoutService(ObjectProvider<ShoppingCart> cartProvider) {
        this.cartProvider = cartProvider;
    }

    public void checkout() {
        ShoppingCart cart = cartProvider.getObject(); // new instance
    }
}
```

**3. Scoped Proxy:**
```java
@Component
@Scope(value = "prototype", proxyMode = ScopedProxyMode.TARGET_CLASS)
public class ShoppingCart { }
// Spring injects a CGLIB proxy; each method call delegates to a fresh instance
```



Aspect Oriented Programming (AOP)
In object-oriented programming, we organize code into classes representing business entities. However, some concerns cut across multiple classes:
for eg. Logging, Security, Transactions, Caching

This leads to code tangling (mixing concerns in one class) and code scattering (same concern duplicated across classes).


AOP Paradigm
AOP introduces a new dimension of modularity by allowing you to define:

What to do (the advice — logging code, security check, etc.)
Where to do it (the pointcut — which methods/classes)
When to do it (advice type — before, after, around)


Spring MVC & REST:

MVC is an architectural pattern that separates an application into three interconnected components:

1. Model : Business logic, states, and logic :: Domain objects, @service, Model attribute
2. View : presentation logic, rendering :: JSP, thymeleaf, JSON/XML Serializers
3. Controller : request handling, Flow control :: @Controller, @Restcontroller

Front Controller pattern:
A single controller (DispatcherServlet) handles all incoming requests, then delegates to appropriate handlers.

Benefits:

Centralized control — Common processing (security, logging) in one place
Consistent handling — All requests follow the same pipeline
Decoupling — Controllers don't need to know Servlet API details.

Request-Response Lifecycle Theory

1. Request arrives at server
2. Servlet container routes to DispatcherServlet
3. DispatcherServlet consults HandlerMapping(s)
4. Appropriate handler (controller method) identified
5. HandlerAdapter invokes the handler
6. Handler returns ModelAndView (or @ResponseBody data)
7. ViewResolver resolves logical view name to actual View
8. View renders the response
9. Response sent to client

   Client Request
         │
         ▼
┌─────────────────────────┐
│    DispatcherServlet    │  ← Front Controller
│    (single entry point) │
└────────┬────────────────┘
         │
    ① Consult HandlerMapping
         │
         ▼
    ② Find Handler (Controller method)
         │
    ③ Call HandlerAdapter
         │
         ▼
┌──────────────────────────┐
│   @Controller / @Rest    │
│   method executes        │
│   returns ModelAndView   │
│   or @ResponseBody       │
└────────┬─────────────────┘
         │
    ④ ViewResolver (if needed)
         │
         ▼
    ⑤ Render View / Write JSON
         │
         ▼
  HTTP Response → Client