# Microservices & Spring Cloud — Complete Guide (Part 2)

> A comprehensive guide covering advanced Microservices patterns with Spring Boot and Spring Cloud.  
> **Part 2 Scope:** Distributed Tracing, Security, Event-Driven Architecture, Saga Pattern, Docker/Kubernetes Deployment  
> **Prerequisites:** Part 1 (Fundamentals, Communication, Discovery, Gateway, Config, Resilience)

---

## Table of Contents

1. [Distributed Tracing & Observability](#1-distributed-tracing--observability)
2. [Security in Microservices](#2-security-in-microservices)
3. [Event-Driven Architecture](#3-event-driven-architecture)
4. [Saga Pattern & Distributed Transactions](#4-saga-pattern--distributed-transactions)
5. [Docker & Container Orchestration](#5-docker--container-orchestration)
6. [Kubernetes Deployment](#6-kubernetes-deployment)
7. [Testing Microservices](#7-testing-microservices)
8. [Interview Questions - Part 2](#8-interview-questions---part-2)

---

## 1. Distributed Tracing & Observability

### 1.1 The Observability Challenge

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Distributed Tracing Problem                                   │
│                                                                                  │
│  Request Flow: Client → Gateway → Order → User → Inventory → Payment            │
│                                                                                  │
│  WITHOUT TRACING:                                                                │
│  ────────────────                                                                │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐       │
│  │ Gateway │ ─▶ │  Order  │ ─▶ │  User   │ ─▶ │Inventory│ ─▶ │ Payment │       │
│  │  Log A  │    │  Log B  │    │  Log C  │    │  Log D  │    │  Log E  │       │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘       │
│                                                                                  │
│  Question: Request failed. Which service? Which step? How long at each?         │
│  Answer: 🤷 Check 5 different log systems, correlate timestamps manually        │
│                                                                                  │
│  WITH DISTRIBUTED TRACING:                                                       │
│  ─────────────────────────                                                       │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐       │
│  │ Gateway │ ─▶ │  Order  │ ─▶ │  User   │ ─▶ │Inventory│ ─▶ │ Payment │       │
│  │TraceID:A│    │TraceID:A│    │TraceID:A│    │TraceID:A│    │TraceID:A│       │
│  │SpanID:1 │    │SpanID:2 │    │SpanID:3 │    │SpanID:4 │    │SpanID:5 │       │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘    └─────────┘       │
│       │              │              │              │              │             │
│       └──────────────┴──────────────┴──────────────┴──────────────┘             │
│                                    │                                             │
│                                    ▼                                             │
│                         ┌─────────────────────┐                                  │
│                         │   Tracing Backend   │                                  │
│                         │ (Zipkin / Jaeger)   │                                  │
│                         │                     │                                  │
│                         │ TraceID: A          │                                  │
│                         │ ├─ Span 1: Gateway  │  25ms                           │
│                         │ ├─ Span 2: Order    │  150ms                          │
│                         │ ├─ Span 3: User     │  45ms                           │
│                         │ ├─ Span 4: Inventory│  30ms                           │
│                         │ └─ Span 5: Payment  │  200ms ← Slow!                  │
│                         └─────────────────────┘                                  │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 1.2 Tracing Concepts

| Concept | Description | Example |
|---------|-------------|---------|
| **Trace** | End-to-end journey of a request | Complete order flow |
| **Span** | Single unit of work within a trace | "Call user service" |
| **Trace ID** | Unique identifier for entire request | `abc123def456` |
| **Span ID** | Unique identifier for each operation | `span001`, `span002` |
| **Parent Span ID** | Links child span to parent | Creates hierarchy |
| **Tags** | Key-value metadata | `http.method=GET` |
| **Logs** | Time-stamped events within span | "User found", "DB query" |
| **Baggage** | Data propagated across services | `user-id`, `tenant-id` |

### 1.3 Micrometer Tracing (Spring Boot 3.x)

```xml
<!-- pom.xml dependencies -->
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-tracing-bridge-brave</artifactId>
</dependency>
<dependency>
    <groupId>io.zipkin.reporter2</groupId>
    <artifactId>zipkin-reporter-brave</artifactId>
</dependency>

<!-- For Jaeger instead -->
<!--
<dependency>
    <groupId>io.micrometer</groupId>
    <artifactId>micrometer-tracing-bridge-otel</artifactId>
</dependency>
<dependency>
    <groupId>io.opentelemetry</groupId>
    <artifactId>opentelemetry-exporter-otlp</artifactId>
</dependency>
-->
```

```yaml
# application.yml - Tracing Configuration
spring:
  application:
    name: order-service

management:
  tracing:
    sampling:
      probability: 1.0  # 100% sampling (use 0.1 for production)
  zipkin:
    tracing:
      endpoint: http://localhost:9411/api/v2/spans

logging:
  pattern:
    level: "%5p [${spring.application.name:},%X{traceId:-},%X{spanId:-}]"
```

### 1.4 Custom Spans and Annotations

```java
import io.micrometer.observation.annotation.Observed;
import io.micrometer.tracing.Tracer;
import io.micrometer.tracing.Span;

@Service
@Slf4j
public class OrderService {
    
    private final Tracer tracer;
    private final UserClient userClient;
    private final InventoryClient inventoryClient;
    
    // Automatic span with @Observed
    @Observed(name = "order.create", 
              contextualName = "creating-order",
              lowCardinalityKeyValues = {"order.type", "standard"})
    public Order createOrder(OrderRequest request) {
        log.info("Creating order for user: {}", request.getUserId());
        
        // Validate user - automatic child span via Feign
        User user = userClient.getUser(request.getUserId());
        
        // Check inventory with manual span
        checkInventory(request.getItems());
        
        // Create order
        Order order = orderRepository.save(mapToOrder(request));
        
        return order;
    }
    
    // Manual span creation
    private void checkInventory(List<OrderItem> items) {
        Span inventorySpan = tracer.nextSpan().name("check-inventory");
        
        try (Tracer.SpanInScope ws = tracer.withSpan(inventorySpan.start())) {
            inventorySpan.tag("items.count", String.valueOf(items.size()));
            
            for (OrderItem item : items) {
                Span itemSpan = tracer.nextSpan().name("check-item-" + item.getSku());
                try (Tracer.SpanInScope itemScope = tracer.withSpan(itemSpan.start())) {
                    boolean available = inventoryClient.checkAvailability(
                        item.getSku(), item.getQuantity());
                    
                    itemSpan.tag("sku", item.getSku());
                    itemSpan.tag("available", String.valueOf(available));
                    
                    if (!available) {
                        itemSpan.event("Item out of stock");
                        throw new InsufficientStockException(item.getSku());
                    }
                } finally {
                    itemSpan.end();
                }
            }
            
            inventorySpan.event("All items available");
        } catch (Exception e) {
            inventorySpan.error(e);
            throw e;
        } finally {
            inventorySpan.end();
        }
    }
    
    // Adding baggage (propagated data)
    public void processWithBaggage(String tenantId) {
        BaggageInScope baggage = tracer.createBaggageInScope("tenant-id", tenantId);
        try {
            // All downstream calls will have tenant-id in baggage
            processOrder();
        } finally {
            baggage.close();
        }
    }
}
```

### 1.5 Enabling Observation

```java
@Configuration
public class ObservationConfig {
    
    @Bean
    public ObservationRegistry observationRegistry() {
        return ObservationRegistry.create();
    }
    
    // Custom observation handler
    @Bean
    public ObservationHandler<Observation.Context> customHandler() {
        return new ObservationHandler<>() {
            @Override
            public void onStart(Observation.Context context) {
                log.info("Observation started: {}", context.getName());
            }
            
            @Override
            public void onStop(Observation.Context context) {
                log.info("Observation stopped: {}", context.getName());
            }
            
            @Override
            public boolean supportsContext(Observation.Context context) {
                return true;
            }
        };
    }
}

// Enable observation for WebClient
@Bean
public WebClient.Builder webClientBuilder(ObservationRegistry registry) {
    return WebClient.builder()
        .observationRegistry(registry);
}

// Enable observation for RestTemplate
@Bean
public RestTemplate restTemplate(ObservationRegistry registry) {
    RestTemplate restTemplate = new RestTemplate();
    restTemplate.setObservationRegistry(registry);
    return restTemplate;
}
```

### 1.6 Zipkin Setup

```yaml
# docker-compose.yml
services:
  zipkin:
    image: openzipkin/zipkin:latest
    ports:
      - "9411:9411"
    environment:
      - STORAGE_TYPE=elasticsearch
      - ES_HOSTS=elasticsearch:9200
    depends_on:
      - elasticsearch
  
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.11.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"
```

### 1.7 Jaeger Setup with OpenTelemetry

```yaml
# docker-compose.yml
services:
  jaeger:
    image: jaegertracing/all-in-one:latest
    ports:
      - "16686:16686"  # UI
      - "4317:4317"    # OTLP gRPC
      - "4318:4318"    # OTLP HTTP
    environment:
      - COLLECTOR_OTLP_ENABLED=true
```

```yaml
# application.yml for Jaeger/OTLP
management:
  otlp:
    tracing:
      endpoint: http://localhost:4318/v1/traces
  tracing:
    sampling:
      probability: 1.0
```

### 1.8 Logging Correlation

```java
// Logs automatically include trace and span IDs
// Log pattern: %5p [service-name,traceId,spanId]

@RestController
@Slf4j
public class OrderController {
    
    @PostMapping("/orders")
    public Order createOrder(@RequestBody OrderRequest request) {
        // Log will show: INFO [order-service,abc123,span456] Creating order...
        log.info("Creating order for user: {}", request.getUserId());
        return orderService.createOrder(request);
    }
}

// Custom MDC propagation
@Component
public class TracingFilter implements Filter {
    
    private final Tracer tracer;
    
    @Override
    public void doFilter(ServletRequest request, ServletResponse response, 
                         FilterChain chain) throws IOException, ServletException {
        Span span = tracer.currentSpan();
        if (span != null) {
            MDC.put("traceId", span.context().traceId());
            MDC.put("spanId", span.context().spanId());
            MDC.put("userId", extractUserId((HttpServletRequest) request));
        }
        try {
            chain.doFilter(request, response);
        } finally {
            MDC.clear();
        }
    }
}
```

### 1.9 Metrics with Micrometer

```java
@Service
public class OrderService {
    
    private final MeterRegistry meterRegistry;
    private final Counter orderCounter;
    private final Timer orderTimer;
    
    public OrderService(MeterRegistry meterRegistry) {
        this.meterRegistry = meterRegistry;
        
        // Counter for order counts
        this.orderCounter = Counter.builder("orders.created")
            .tag("service", "order-service")
            .description("Total orders created")
            .register(meterRegistry);
        
        // Timer for order processing time
        this.orderTimer = Timer.builder("orders.processing.time")
            .tag("service", "order-service")
            .description("Order processing time")
            .register(meterRegistry);
    }
    
    public Order createOrder(OrderRequest request) {
        return orderTimer.record(() -> {
            Order order = processOrder(request);
            
            // Increment counter with tags
            orderCounter.increment();
            
            // Custom gauge for pending orders
            meterRegistry.gauge("orders.pending", 
                orderRepository.countByStatus("PENDING"));
            
            return order;
        });
    }
    
    // Distribution summary for order amounts
    public void recordOrderAmount(BigDecimal amount) {
        DistributionSummary.builder("order.amount")
            .tag("currency", "USD")
            .publishPercentiles(0.5, 0.95, 0.99)
            .register(meterRegistry)
            .record(amount.doubleValue());
    }
}
```

```yaml
# Prometheus endpoint
management:
  endpoints:
    web:
      exposure:
        include: health,info,prometheus,metrics
  metrics:
    export:
      prometheus:
        enabled: true
    tags:
      application: ${spring.application.name}
```

---

## 2. Security in Microservices

### 2.1 Security Challenges in Microservices

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Microservices Security Challenges                             │
│                                                                                  │
│  ┌───────────────────────────────────────────────────────────────────────────┐  │
│  │ MONOLITH SECURITY                    MICROSERVICES SECURITY              │  │
│  │ ────────────────────                 ────────────────────────            │  │
│  │                                                                          │  │
│  │  ┌─────────────────┐                 ┌─────┐  ┌─────┐  ┌─────┐          │  │
│  │  │   Application   │                 │ Svc │  │ Svc │  │ Svc │          │  │
│  │  │  ┌───────────┐  │                 │  A  │  │  B  │  │  C  │          │  │
│  │  │  │ Security  │  │                 └──┬──┘  └──┬──┘  └──┬──┘          │  │
│  │  │  │  Layer    │  │                    │        │        │              │  │
│  │  │  └───────────┘  │                    └────────┼────────┘              │  │
│  │  │   One place     │                             │                       │  │
│  │  └─────────────────┘                     ┌───────┴───────┐               │  │
│  │                                          │ How to secure │               │  │
│  │  Simple: One login,                      │ all of them?  │               │  │
│  │  one session                             └───────────────┘               │  │
│  │                                                                          │  │
│  └───────────────────────────────────────────────────────────────────────────┘  │
│                                                                                  │
│  CHALLENGES:                                                                     │
│  ───────────                                                                     │
│  • Authentication at each service?                                               │
│  • How to pass user identity between services?                                   │
│  • Service-to-service authentication                                             │
│  • Token validation overhead                                                     │
│  • Secret management across services                                             │
│  • Network security between services                                             │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.2 OAuth2 / OpenID Connect Architecture

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         OAuth2 + OIDC Flow                                       │
│                                                                                  │
│  ┌─────────┐                      ┌─────────────────────────┐                   │
│  │  User   │                      │   Authorization Server  │                   │
│  │(Browser)│                      │   (Keycloak/Auth0/Okta) │                   │
│  └────┬────┘                      └────────────┬────────────┘                   │
│       │                                        │                                 │
│       │ 1. Access /orders                      │                                 │
│       ▼                                        │                                 │
│  ┌─────────────┐                               │                                 │
│  │ API Gateway │                               │                                 │
│  └──────┬──────┘                               │                                 │
│         │                                      │                                 │
│         │ 2. No token? Redirect to login       │                                 │
│         │────────────────────────────────────▶│                                 │
│         │                                      │                                 │
│         │ 3. User logs in                      │                                 │
│         │◀────────────────────────────────────│                                 │
│         │    (Authorization Code)              │                                 │
│         │                                      │                                 │
│         │ 4. Exchange code for tokens          │                                 │
│         │────────────────────────────────────▶│                                 │
│         │                                      │                                 │
│         │ 5. Receive tokens                    │                                 │
│         │◀────────────────────────────────────│                                 │
│         │    - Access Token (JWT)              │                                 │
│         │    - Refresh Token                   │                                 │
│         │    - ID Token (OIDC)                 │                                 │
│         │                                      │                                 │
│         │ 6. Forward request + Access Token    │                                 │
│         ▼                                      │                                 │
│  ┌─────────────┐                               │                                 │
│  │Order Service│ 7. Validate token             │                                 │
│  │             │────────────────────────────▶│                                 │
│  │             │◀───────────────────────────── │                                 │
│  │             │   (or validate locally with   │                                 │
│  │             │    public key)                │                                 │
│  └─────────────┘                               │                                 │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 2.3 Keycloak Setup

```yaml
# docker-compose.yml
services:
  keycloak:
    image: quay.io/keycloak/keycloak:23.0
    command: start-dev
    environment:
      - KEYCLOAK_ADMIN=admin
      - KEYCLOAK_ADMIN_PASSWORD=admin
      - KC_DB=postgres
      - KC_DB_URL=jdbc:postgresql://postgres:5432/keycloak
      - KC_DB_USERNAME=keycloak
      - KC_DB_PASSWORD=keycloak
    ports:
      - "8180:8080"
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=keycloak
      - POSTGRES_USER=keycloak
      - POSTGRES_PASSWORD=keycloak
    volumes:
      - keycloak_data:/var/lib/postgresql/data

volumes:
  keycloak_data:
```

### 2.4 Resource Server Configuration

```java
// Order Service - Resource Server
// pom.xml: spring-boot-starter-oauth2-resource-server

@Configuration
@EnableWebSecurity
@EnableMethodSecurity
public class SecurityConfig {
    
    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/actuator/**").permitAll()
                .requestMatchers("/api/public/**").permitAll()
                .requestMatchers("/api/orders/**").hasRole("USER")
                .requestMatchers("/api/admin/**").hasRole("ADMIN")
                .anyRequest().authenticated()
            )
            .oauth2ResourceServer(oauth2 -> oauth2
                .jwt(jwt -> jwt
                    .jwtAuthenticationConverter(jwtAuthenticationConverter())
                )
            )
            .sessionManagement(session -> 
                session.sessionCreationPolicy(SessionCreationPolicy.STATELESS)
            )
            .csrf(csrf -> csrf.disable());
        
        return http.build();
    }
    
    // Convert Keycloak roles to Spring Security authorities
    @Bean
    public JwtAuthenticationConverter jwtAuthenticationConverter() {
        JwtGrantedAuthoritiesConverter grantedAuthoritiesConverter = 
            new JwtGrantedAuthoritiesConverter();
        grantedAuthoritiesConverter.setAuthorityPrefix("ROLE_");
        grantedAuthoritiesConverter.setAuthoritiesClaimName("realm_access.roles");
        
        JwtAuthenticationConverter converter = new JwtAuthenticationConverter();
        converter.setJwtGrantedAuthoritiesConverter(jwt -> {
            Collection<GrantedAuthority> authorities = new ArrayList<>();
            
            // Extract realm roles
            Map<String, Object> realmAccess = jwt.getClaimAsMap("realm_access");
            if (realmAccess != null) {
                List<String> roles = (List<String>) realmAccess.get("roles");
                if (roles != null) {
                    authorities.addAll(roles.stream()
                        .map(role -> new SimpleGrantedAuthority("ROLE_" + role.toUpperCase()))
                        .collect(Collectors.toList()));
                }
            }
            
            // Extract resource/client roles
            Map<String, Object> resourceAccess = jwt.getClaimAsMap("resource_access");
            if (resourceAccess != null) {
                Map<String, Object> clientAccess = 
                    (Map<String, Object>) resourceAccess.get("order-service");
                if (clientAccess != null) {
                    List<String> clientRoles = (List<String>) clientAccess.get("roles");
                    if (clientRoles != null) {
                        authorities.addAll(clientRoles.stream()
                            .map(role -> new SimpleGrantedAuthority("ROLE_" + role.toUpperCase()))
                            .collect(Collectors.toList()));
                    }
                }
            }
            
            return authorities;
        });
        
        return converter;
    }
}
```

```yaml
# application.yml - Resource Server
spring:
  security:
    oauth2:
      resourceserver:
        jwt:
          issuer-uri: http://localhost:8180/realms/microservices
          jwk-set-uri: http://localhost:8180/realms/microservices/protocol/openid-connect/certs
```

### 2.5 Accessing User Information

```java
@RestController
@RequestMapping("/api/orders")
public class OrderController {
    
    private final OrderService orderService;
    
    // Get current user from JWT
    @GetMapping("/my-orders")
    public List<Order> getMyOrders(@AuthenticationPrincipal Jwt jwt) {
        String userId = jwt.getSubject();
        String email = jwt.getClaimAsString("email");
        String name = jwt.getClaimAsString("preferred_username");
        
        log.info("Fetching orders for user: {} ({})", name, userId);
        return orderService.getOrdersByUserId(userId);
    }
    
    // Using SecurityContextHolder
    @PostMapping
    public Order createOrder(@RequestBody OrderRequest request) {
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        Jwt jwt = (Jwt) auth.getPrincipal();
        
        request.setUserId(jwt.getSubject());
        return orderService.createOrder(request);
    }
    
    // Method-level security
    @PreAuthorize("hasRole('ADMIN') or #userId == authentication.principal.subject")
    @GetMapping("/user/{userId}")
    public List<Order> getOrdersByUser(@PathVariable String userId) {
        return orderService.getOrdersByUserId(userId);
    }
    
    @PreAuthorize("hasAuthority('SCOPE_orders:write')")
    @DeleteMapping("/{orderId}")
    public void deleteOrder(@PathVariable Long orderId) {
        orderService.deleteOrder(orderId);
    }
}
```

### 2.6 Service-to-Service Authentication

```java
// Option 1: Token Relay (propagate user token)
@Configuration
public class WebClientConfig {
    
    @Bean
    public WebClient webClient(ReactiveClientRegistrationRepository clientRegistrations,
                                ServerOAuth2AuthorizedClientRepository authorizedClients) {
        ServerOAuth2AuthorizedClientExchangeFilterFunction oauth2 =
            new ServerOAuth2AuthorizedClientExchangeFilterFunction(
                clientRegistrations, authorizedClients);
        oauth2.setDefaultClientRegistrationId("keycloak");
        
        return WebClient.builder()
            .filter(oauth2)
            .build();
    }
}

// Option 2: Client Credentials (service account)
@Configuration
public class ServiceAccountConfig {
    
    @Bean
    public OAuth2AuthorizedClientManager authorizedClientManager(
            ClientRegistrationRepository clientRegistrationRepository,
            OAuth2AuthorizedClientService authorizedClientService) {
        
        OAuth2AuthorizedClientProvider authorizedClientProvider =
            OAuth2AuthorizedClientProviderBuilder.builder()
                .clientCredentials()
                .build();
        
        AuthorizedClientServiceOAuth2AuthorizedClientManager authorizedClientManager =
            new AuthorizedClientServiceOAuth2AuthorizedClientManager(
                clientRegistrationRepository, authorizedClientService);
        authorizedClientManager.setAuthorizedClientProvider(authorizedClientProvider);
        
        return authorizedClientManager;
    }
}

@Service
public class UserServiceClient {
    
    private final WebClient webClient;
    private final OAuth2AuthorizedClientManager authorizedClientManager;
    
    public User getUser(Long userId) {
        OAuth2AuthorizedClient client = authorizedClientManager.authorize(
            OAuth2AuthorizeRequest.withClientRegistrationId("user-service-client")
                .principal("order-service")
                .build());
        
        return webClient.get()
            .uri("http://user-service/api/users/{id}", userId)
            .headers(headers -> headers.setBearerAuth(client.getAccessToken().getTokenValue()))
            .retrieve()
            .bodyToMono(User.class)
            .block();
    }
}
```

```yaml
# application.yml - Client Credentials
spring:
  security:
    oauth2:
      client:
        registration:
          user-service-client:
            provider: keycloak
            client-id: order-service
            client-secret: ${ORDER_SERVICE_SECRET}
            authorization-grant-type: client_credentials
            scope: openid,profile
        provider:
          keycloak:
            issuer-uri: http://localhost:8180/realms/microservices
```

### 2.7 API Gateway Security

```java
// Gateway with OAuth2 Login
@Configuration
@EnableWebFluxSecurity
public class GatewaySecurityConfig {
    
    @Bean
    public SecurityWebFilterChain securityWebFilterChain(ServerHttpSecurity http) {
        http
            .authorizeExchange(exchanges -> exchanges
                .pathMatchers("/actuator/**").permitAll()
                .pathMatchers("/api/public/**").permitAll()
                .anyExchange().authenticated()
            )
            .oauth2Login(Customizer.withDefaults())
            .oauth2ResourceServer(oauth2 -> oauth2.jwt(Customizer.withDefaults()));
        
        return http.build();
    }
    
    // Token relay filter - passes token to downstream services
    @Bean
    public TokenRelayGatewayFilterFactory tokenRelayFilterFactory(
            ReactiveClientRegistrationRepository clientRegistrationRepository) {
        return new TokenRelayGatewayFilterFactory(clientRegistrationRepository);
    }
}
```

```yaml
# Gateway application.yml
spring:
  cloud:
    gateway:
      routes:
        - id: order-service
          uri: lb://order-service
          predicates:
            - Path=/api/orders/**
          filters:
            - TokenRelay=
            - RemoveRequestHeader=Cookie
```

### 2.8 JWT Token Structure

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              JWT Token Structure                                 │
│                                                                                  │
│  eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.signature                     │
│  ───────────────────── ─────────────────────────── ─────────                    │
│        Header              Payload (Claims)         Signature                    │
│                                                                                  │
│  HEADER:                                                                         │
│  {                                                                               │
│    "alg": "RS256",           // Algorithm                                        │
│    "typ": "JWT",             // Token type                                       │
│    "kid": "abc123"           // Key ID for validation                            │
│  }                                                                               │
│                                                                                  │
│  PAYLOAD (CLAIMS):                                                               │
│  {                                                                               │
│    // Standard claims                                                            │
│    "iss": "http://keycloak/realms/microservices",  // Issuer                    │
│    "sub": "user-uuid-123",                          // Subject (user ID)        │
│    "aud": ["order-service", "user-service"],       // Audience                  │
│    "exp": 1709123456,                              // Expiration                │
│    "iat": 1709120000,                              // Issued at                 │
│    "nbf": 1709120000,                              // Not before                │
│    "jti": "token-uuid",                            // Token ID                  │
│                                                                                  │
│    // Keycloak specific                                                          │
│    "preferred_username": "john.doe",                                             │
│    "email": "john@example.com",                                                  │
│    "realm_access": {                                                             │
│      "roles": ["USER", "ADMIN"]                                                  │
│    },                                                                            │
│    "resource_access": {                                                          │
│      "order-service": {                                                          │
│        "roles": ["order-manager"]                                                │
│      }                                                                           │
│    },                                                                            │
│    "scope": "openid profile email"                                               │
│  }                                                                               │
│                                                                                  │
│  SIGNATURE:                                                                      │
│  RSASHA256(base64UrlEncode(header) + "." + base64UrlEncode(payload), privateKey) │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Event-Driven Architecture

### 3.1 Event-Driven Patterns

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Event-Driven Architecture Patterns                            │
│                                                                                  │
│  1. EVENT NOTIFICATION                                                           │
│  ─────────────────────                                                           │
│  ┌─────────┐   OrderCreated    ┌─────────────┐                                  │
│  │  Order  │ ────────────────▶ │  Inventory  │  "Something happened"            │
│  │ Service │   (minimal data)  │   Service   │  Consumer fetches details        │
│  └─────────┘                   └─────────────┘  if needed                       │
│                                                                                  │
│  2. EVENT-CARRIED STATE TRANSFER                                                 │
│  ───────────────────────────────                                                 │
│  ┌─────────┐   OrderCreated    ┌─────────────┐                                  │
│  │  Order  │ ────────────────▶ │  Inventory  │  "Here's all you need"           │
│  │ Service │   (full data)     │   Service   │  No callback required            │
│  └─────────┘                   └─────────────┘                                  │
│                                                                                  │
│  3. EVENT SOURCING                                                               │
│  ─────────────────                                                               │
│  ┌─────────┐                   ┌─────────────┐                                  │
│  │  Order  │ ────────────────▶ │   Event     │  Store events, not state         │
│  │ Service │   OrderCreated    │   Store     │  Rebuild state from events       │
│  └─────────┘   ItemAdded       └──────┬──────┘                                  │
│                OrderCompleted         │                                          │
│                                       ▼                                          │
│                               ┌─────────────┐                                    │
│                               │ Projections │  Materialized views               │
│                               └─────────────┘                                    │
│                                                                                  │
│  4. CQRS (Command Query Responsibility Segregation)                              │
│  ───────────────────────────────────────────────                                 │
│          Commands                      Queries                                   │
│              │                            │                                      │
│              ▼                            ▼                                      │
│       ┌─────────────┐            ┌─────────────┐                                │
│       │   Write     │   Events   │    Read     │                                │
│       │   Model     │ ─────────▶ │   Model     │                                │
│       │ (normalized)│            │(denormalized)│                                │
│       └─────────────┘            └─────────────┘                                │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Spring Cloud Stream

```java
// Spring Cloud Stream with Kafka/RabbitMQ
// pom.xml: spring-cloud-starter-stream-kafka (or stream-rabbit)

@SpringBootApplication
public class OrderServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(OrderServiceApplication.class, args);
    }
    
    // Functional style (Spring Cloud Stream 3.x+)
    @Bean
    public Supplier<OrderCreatedEvent> orderSource() {
        return () -> new OrderCreatedEvent();  // Produces messages
    }
    
    @Bean
    public Consumer<OrderCreatedEvent> orderSink() {
        return event -> processEvent(event);   // Consumes messages
    }
    
    @Bean
    public Function<OrderCreatedEvent, InventoryReservedEvent> orderProcessor() {
        return event -> {                       // Transforms messages
            // Process and return new event
            return new InventoryReservedEvent(event.getOrderId());
        };
    }
}
```

```yaml
# application.yml - Spring Cloud Stream
spring:
  cloud:
    stream:
      bindings:
        # Supplier binding
        orderSource-out-0:
          destination: orders
          content-type: application/json
        
        # Consumer binding
        orderSink-in-0:
          destination: orders
          group: inventory-service
          content-type: application/json
        
        # Function binding (in and out)
        orderProcessor-in-0:
          destination: orders
          group: order-processor
        orderProcessor-out-0:
          destination: inventory-events
      
      kafka:
        binder:
          brokers: localhost:9092
          auto-create-topics: true
        bindings:
          orderSink-in-0:
            consumer:
              start-offset: earliest
              enable-dlq: true
              dlq-name: orders-dlq
```

### 3.3 Domain Events

```java
// Domain Event
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class OrderCreatedEvent {
    private String eventId;
    private String eventType;
    private LocalDateTime timestamp;
    private String aggregateId;
    
    // Event payload
    private Long orderId;
    private String userId;
    private List<OrderItemDto> items;
    private BigDecimal totalAmount;
    private String status;
    
    public static OrderCreatedEvent from(Order order) {
        return OrderCreatedEvent.builder()
            .eventId(UUID.randomUUID().toString())
            .eventType("ORDER_CREATED")
            .timestamp(LocalDateTime.now())
            .aggregateId(order.getId().toString())
            .orderId(order.getId())
            .userId(order.getUserId())
            .items(order.getItems().stream()
                .map(OrderItemDto::from)
                .collect(Collectors.toList()))
            .totalAmount(order.getTotalAmount())
            .status(order.getStatus().name())
            .build();
    }
}

// Event Publisher
@Service
@Slf4j
public class OrderEventPublisher {
    
    private final StreamBridge streamBridge;
    
    public void publish(OrderCreatedEvent event) {
        log.info("Publishing OrderCreatedEvent: {}", event.getEventId());
        
        boolean sent = streamBridge.send("orders-out-0", event);
        
        if (!sent) {
            log.error("Failed to publish event: {}", event.getEventId());
            // Handle failure - maybe save to outbox table
        }
    }
}

// Order Service with Event Publishing
@Service
@Transactional
public class OrderService {
    
    private final OrderRepository orderRepository;
    private final OrderEventPublisher eventPublisher;
    
    public Order createOrder(CreateOrderRequest request) {
        // Create order
        Order order = Order.builder()
            .userId(request.getUserId())
            .items(request.getItems())
            .status(OrderStatus.CREATED)
            .build();
        
        order = orderRepository.save(order);
        
        // Publish event
        eventPublisher.publish(OrderCreatedEvent.from(order));
        
        return order;
    }
}
```

### 3.4 Event Consumer

```java
@Service
@Slf4j
public class OrderEventConsumer {
    
    private final InventoryService inventoryService;
    
    @Bean
    public Consumer<OrderCreatedEvent> processOrder() {
        return event -> {
            log.info("Received OrderCreatedEvent: {}", event.getEventId());
            
            try {
                // Reserve inventory
                for (OrderItemDto item : event.getItems()) {
                    inventoryService.reserveStock(
                        item.getSku(), 
                        item.getQuantity()
                    );
                }
                
                log.info("Inventory reserved for order: {}", event.getOrderId());
                
            } catch (InsufficientStockException e) {
                log.error("Insufficient stock for order: {}", event.getOrderId());
                // Publish compensation event
                throw e;  // Will trigger DLQ
            }
        };
    }
}
```

### 3.5 Transactional Outbox Pattern

```java
// Outbox table entity
@Entity
@Table(name = "outbox_events")
@Data
public class OutboxEvent {
    @Id
    private String id;
    
    private String aggregateType;
    private String aggregateId;
    private String eventType;
    
    @Column(columnDefinition = "TEXT")
    private String payload;
    
    private LocalDateTime createdAt;
    private boolean published;
}

// Order Service with Outbox
@Service
@Transactional
public class OrderService {
    
    private final OrderRepository orderRepository;
    private final OutboxRepository outboxRepository;
    private final ObjectMapper objectMapper;
    
    public Order createOrder(CreateOrderRequest request) {
        // Create order
        Order order = Order.builder()
            .userId(request.getUserId())
            .items(request.getItems())
            .status(OrderStatus.CREATED)
            .build();
        
        order = orderRepository.save(order);
        
        // Save event to outbox (same transaction)
        OutboxEvent outboxEvent = new OutboxEvent();
        outboxEvent.setId(UUID.randomUUID().toString());
        outboxEvent.setAggregateType("Order");
        outboxEvent.setAggregateId(order.getId().toString());
        outboxEvent.setEventType("OrderCreated");
        outboxEvent.setPayload(objectMapper.writeValueAsString(
            OrderCreatedEvent.from(order)));
        outboxEvent.setCreatedAt(LocalDateTime.now());
        outboxEvent.setPublished(false);
        
        outboxRepository.save(outboxEvent);
        
        return order;
    }
}

// Outbox Poller (separate process)
@Service
@Slf4j
public class OutboxPoller {
    
    private final OutboxRepository outboxRepository;
    private final StreamBridge streamBridge;
    
    @Scheduled(fixedDelay = 1000)
    @Transactional
    public void pollAndPublish() {
        List<OutboxEvent> events = outboxRepository
            .findByPublishedFalseOrderByCreatedAtAsc();
        
        for (OutboxEvent event : events) {
            try {
                boolean sent = streamBridge.send(
                    getDestination(event.getEventType()),
                    event.getPayload()
                );
                
                if (sent) {
                    event.setPublished(true);
                    outboxRepository.save(event);
                    log.info("Published outbox event: {}", event.getId());
                }
            } catch (Exception e) {
                log.error("Failed to publish event: {}", event.getId(), e);
            }
        }
    }
    
    private String getDestination(String eventType) {
        return switch (eventType) {
            case "OrderCreated" -> "orders-out-0";
            case "OrderCancelled" -> "order-cancellations-out-0";
            default -> "events-out-0";
        };
    }
}
```

### 3.6 Debezium CDC (Change Data Capture)

```yaml
# docker-compose.yml - Debezium setup
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
  
  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092
    ports:
      - "9092:9092"
  
  connect:
    image: debezium/connect:latest
    depends_on:
      - kafka
    ports:
      - "8083:8083"
    environment:
      BOOTSTRAP_SERVERS: kafka:9092
      GROUP_ID: 1
      CONFIG_STORAGE_TOPIC: connect_configs
      OFFSET_STORAGE_TOPIC: connect_offsets
      STATUS_STORAGE_TOPIC: connect_statuses
```

```json
// Create Debezium connector
// POST http://localhost:8083/connectors
{
  "name": "order-outbox-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "database.hostname": "postgres",
    "database.port": "5432",
    "database.user": "order_service",
    "database.password": "password",
    "database.dbname": "orders",
    "database.server.name": "order-service",
    "table.include.list": "public.outbox_events",
    "transforms": "outbox",
    "transforms.outbox.type": "io.debezium.transforms.outbox.EventRouter",
    "transforms.outbox.table.field.event.key": "aggregate_id",
    "transforms.outbox.table.field.event.type": "event_type",
    "transforms.outbox.table.field.event.payload": "payload",
    "transforms.outbox.route.by.field": "aggregate_type"
  }
}
```

---

## 4. Saga Pattern & Distributed Transactions

### 4.1 The Distributed Transaction Problem

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                    Distributed Transaction Problem                               │
│                                                                                  │
│  ORDER CREATION FLOW:                                                            │
│                                                                                  │
│  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐                      │
│  │  Order  │ ─▶ │Inventory│ ─▶ │ Payment │ ─▶ │Shipping │                      │
│  │ Service │    │ Service │    │ Service │    │ Service │                      │
│  └─────────┘    └─────────┘    └─────────┘    └─────────┘                      │
│       │              │              │              │                            │
│  Create Order  Reserve Stock  Charge Card   Create Shipment                    │
│       ✓              ✓              ✗                                           │
│                                  FAILED!                                        │
│                                                                                  │
│  PROBLEM: How to rollback Order and Inventory?                                  │
│  • Each service has its own database                                            │
│  • No distributed ACID transactions                                             │
│  • Network can fail between any steps                                           │
│                                                                                  │
│  SOLUTION: SAGA PATTERN                                                          │
│  • Sequence of local transactions                                               │
│  • Each step has a compensating transaction                                     │
│  • Eventually consistent                                                         │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 4.2 Saga Pattern Types

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         Saga Pattern Types                                       │
│                                                                                  │
│  1. CHOREOGRAPHY-BASED SAGA                                                      │
│  ───────────────────────────                                                     │
│  No central coordinator. Services react to events.                               │
│                                                                                  │
│  ┌─────────┐  OrderCreated   ┌─────────┐ StockReserved  ┌─────────┐            │
│  │  Order  │ ───────────────▶│Inventory│ ──────────────▶│ Payment │            │
│  │ Service │                 │ Service │                │ Service │            │
│  └─────────┘                 └─────────┘                └─────────┘            │
│       ▲                           │                          │                  │
│       │    OrderCancelled         │   StockReserveFailed    │                  │
│       └───────────────────────────┴──────────────────────────┘                  │
│                              (Compensation)                                      │
│                                                                                  │
│  Pros: Simple, decoupled                                                         │
│  Cons: Hard to track, cyclic dependencies possible                              │
│                                                                                  │
│  2. ORCHESTRATION-BASED SAGA                                                     │
│  ────────────────────────────                                                    │
│  Central orchestrator coordinates the saga.                                      │
│                                                                                  │
│                      ┌─────────────────────┐                                    │
│                      │   Saga Orchestrator │                                    │
│                      │   (Order Saga)      │                                    │
│                      └─────────┬───────────┘                                    │
│                                │                                                 │
│            ┌───────────────────┼───────────────────┐                            │
│            │                   │                   │                            │
│            ▼                   ▼                   ▼                            │
│       ┌─────────┐        ┌─────────┐        ┌─────────┐                        │
│       │Inventory│        │ Payment │        │Shipping │                        │
│       │ Service │        │ Service │        │ Service │                        │
│       └─────────┘        └─────────┘        └─────────┘                        │
│                                                                                  │
│  Pros: Easy to understand, centralized logic                                    │
│  Cons: Single point of failure, orchestrator complexity                         │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 4.3 Choreography Saga Implementation

```java
// Order Service - Initiates Saga
@Service
@Slf4j
public class OrderService {
    
    private final OrderRepository orderRepository;
    private final StreamBridge streamBridge;
    
    @Transactional
    public Order createOrder(CreateOrderRequest request) {
        Order order = Order.builder()
            .userId(request.getUserId())
            .items(request.getItems())
            .status(OrderStatus.PENDING)
            .build();
        
        order = orderRepository.save(order);
        
        // Publish event to start saga
        OrderCreatedEvent event = OrderCreatedEvent.from(order);
        streamBridge.send("order-created", event);
        
        log.info("Order created and saga initiated: {}", order.getId());
        return order;
    }
    
    // Handle compensation
    @Bean
    public Consumer<PaymentFailedEvent> handlePaymentFailed() {
        return event -> {
            log.info("Payment failed, cancelling order: {}", event.getOrderId());
            
            Order order = orderRepository.findById(event.getOrderId())
                .orElseThrow();
            order.setStatus(OrderStatus.CANCELLED);
            order.setFailureReason(event.getReason());
            orderRepository.save(order);
            
            // Publish order cancelled event
            streamBridge.send("order-cancelled", 
                new OrderCancelledEvent(order.getId(), event.getReason()));
        };
    }
}

// Inventory Service - Step 2
@Service
@Slf4j
public class InventoryService {
    
    private final InventoryRepository inventoryRepository;
    private final StreamBridge streamBridge;
    
    @Bean
    public Consumer<OrderCreatedEvent> reserveStock() {
        return event -> {
            log.info("Reserving stock for order: {}", event.getOrderId());
            
            try {
                for (OrderItemDto item : event.getItems()) {
                    Inventory inventory = inventoryRepository
                        .findBySku(item.getSku())
                        .orElseThrow(() -> new NotFoundException("SKU not found"));
                    
                    if (inventory.getAvailable() < item.getQuantity()) {
                        throw new InsufficientStockException(item.getSku());
                    }
                    
                    inventory.setAvailable(inventory.getAvailable() - item.getQuantity());
                    inventory.setReserved(inventory.getReserved() + item.getQuantity());
                    inventoryRepository.save(inventory);
                }
                
                // Success - notify payment service
                streamBridge.send("stock-reserved", StockReservedEvent.builder()
                    .orderId(event.getOrderId())
                    .userId(event.getUserId())
                    .items(event.getItems())
                    .totalAmount(event.getTotalAmount())
                    .build());
                
            } catch (Exception e) {
                log.error("Stock reservation failed: {}", e.getMessage());
                
                // Publish failure event
                streamBridge.send("stock-reserve-failed", StockReserveFailedEvent.builder()
                    .orderId(event.getOrderId())
                    .reason(e.getMessage())
                    .build());
            }
        };
    }
    
    // Compensation handler
    @Bean
    public Consumer<OrderCancelledEvent> releaseStock() {
        return event -> {
            log.info("Releasing stock for cancelled order: {}", event.getOrderId());
            // Release reserved stock
            releaseReservedStock(event.getOrderId());
        };
    }
}

// Payment Service - Step 3
@Service
@Slf4j
public class PaymentService {
    
    private final PaymentRepository paymentRepository;
    private final PaymentGateway paymentGateway;
    private final StreamBridge streamBridge;
    
    @Bean
    public Consumer<StockReservedEvent> processPayment() {
        return event -> {
            log.info("Processing payment for order: {}", event.getOrderId());
            
            try {
                PaymentResult result = paymentGateway.charge(
                    event.getUserId(),
                    event.getTotalAmount()
                );
                
                Payment payment = Payment.builder()
                    .orderId(event.getOrderId())
                    .amount(event.getTotalAmount())
                    .status(PaymentStatus.COMPLETED)
                    .transactionId(result.getTransactionId())
                    .build();
                
                paymentRepository.save(payment);
                
                // Success - notify shipping
                streamBridge.send("payment-completed", PaymentCompletedEvent.builder()
                    .orderId(event.getOrderId())
                    .transactionId(result.getTransactionId())
                    .build());
                
            } catch (PaymentException e) {
                log.error("Payment failed: {}", e.getMessage());
                
                // Trigger compensation
                streamBridge.send("payment-failed", PaymentFailedEvent.builder()
                    .orderId(event.getOrderId())
                    .reason(e.getMessage())
                    .build());
            }
        };
    }
}
```

### 4.4 Orchestration Saga Implementation

```java
// Saga Orchestrator
@Service
@Slf4j
public class CreateOrderSaga {
    
    private final OrderRepository orderRepository;
    private final InventoryClient inventoryClient;
    private final PaymentClient paymentClient;
    private final ShippingClient shippingClient;
    private final SagaStateRepository sagaStateRepository;
    
    @Transactional
    public Order execute(CreateOrderRequest request) {
        String sagaId = UUID.randomUUID().toString();
        SagaState state = new SagaState(sagaId, SagaStep.STARTED);
        
        try {
            // Step 1: Create Order
            state.setStep(SagaStep.ORDER_CREATING);
            Order order = createOrder(request);
            state.setOrderId(order.getId());
            sagaStateRepository.save(state);
            
            // Step 2: Reserve Inventory
            state.setStep(SagaStep.INVENTORY_RESERVING);
            ReservationResult reservation = inventoryClient.reserve(
                new ReserveRequest(order.getId(), request.getItems())
            );
            state.setReservationId(reservation.getId());
            sagaStateRepository.save(state);
            
            // Step 3: Process Payment
            state.setStep(SagaStep.PAYMENT_PROCESSING);
            PaymentResult payment = paymentClient.charge(
                new ChargeRequest(order.getUserId(), order.getTotalAmount())
            );
            state.setPaymentId(payment.getTransactionId());
            sagaStateRepository.save(state);
            
            // Step 4: Create Shipment
            state.setStep(SagaStep.SHIPMENT_CREATING);
            ShipmentResult shipment = shippingClient.createShipment(
                new ShipmentRequest(order.getId(), request.getShippingAddress())
            );
            
            // Success - Complete saga
            state.setStep(SagaStep.COMPLETED);
            order.setStatus(OrderStatus.CONFIRMED);
            orderRepository.save(order);
            sagaStateRepository.save(state);
            
            log.info("Saga completed successfully: {}", sagaId);
            return order;
            
        } catch (Exception e) {
            log.error("Saga failed at step {}: {}", state.getStep(), e.getMessage());
            compensate(state);
            throw new SagaException("Order creation failed", e);
        }
    }
    
    private void compensate(SagaState state) {
        log.info("Starting compensation for saga: {}", state.getSagaId());
        
        try {
            switch (state.getStep()) {
                case SHIPMENT_CREATING:
                    // No compensation needed - shipment wasn't created
                    
                case PAYMENT_PROCESSING:
                    if (state.getPaymentId() != null) {
                        paymentClient.refund(state.getPaymentId());
                        log.info("Payment refunded: {}", state.getPaymentId());
                    }
                    
                case INVENTORY_RESERVING:
                    if (state.getReservationId() != null) {
                        inventoryClient.release(state.getReservationId());
                        log.info("Inventory released: {}", state.getReservationId());
                    }
                    
                case ORDER_CREATING:
                    if (state.getOrderId() != null) {
                        Order order = orderRepository.findById(state.getOrderId())
                            .orElse(null);
                        if (order != null) {
                            order.setStatus(OrderStatus.CANCELLED);
                            orderRepository.save(order);
                            log.info("Order cancelled: {}", state.getOrderId());
                        }
                    }
                    break;
                    
                default:
                    break;
            }
            
            state.setStep(SagaStep.COMPENSATED);
            sagaStateRepository.save(state);
            
        } catch (Exception e) {
            log.error("Compensation failed for saga: {}", state.getSagaId(), e);
            state.setStep(SagaStep.COMPENSATION_FAILED);
            sagaStateRepository.save(state);
            // Alert operations team
        }
    }
}

// Saga State Entity
@Entity
@Data
public class SagaState {
    @Id
    private String sagaId;
    
    @Enumerated(EnumType.STRING)
    private SagaStep step;
    
    private Long orderId;
    private String reservationId;
    private String paymentId;
    
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;
}

public enum SagaStep {
    STARTED,
    ORDER_CREATING,
    INVENTORY_RESERVING,
    PAYMENT_PROCESSING,
    SHIPMENT_CREATING,
    COMPLETED,
    COMPENSATING,
    COMPENSATED,
    COMPENSATION_FAILED
}
```

### 4.5 Saga State Machine (Async Orchestration)

```java
// Using Spring State Machine for complex sagas
@Configuration
@EnableStateMachine
public class OrderSagaStateMachineConfig 
        extends StateMachineConfigurerAdapter<SagaStep, SagaEvent> {
    
    @Override
    public void configure(StateMachineStateConfigurer<SagaStep, SagaEvent> states) 
            throws Exception {
        states
            .withStates()
            .initial(SagaStep.STARTED)
            .state(SagaStep.ORDER_CREATING)
            .state(SagaStep.INVENTORY_RESERVING)
            .state(SagaStep.PAYMENT_PROCESSING)
            .state(SagaStep.SHIPMENT_CREATING)
            .end(SagaStep.COMPLETED)
            .end(SagaStep.COMPENSATED);
    }
    
    @Override
    public void configure(StateMachineTransitionConfigurer<SagaStep, SagaEvent> transitions) 
            throws Exception {
        transitions
            .withExternal()
                .source(SagaStep.STARTED).target(SagaStep.ORDER_CREATING)
                .event(SagaEvent.START)
            .and()
            .withExternal()
                .source(SagaStep.ORDER_CREATING).target(SagaStep.INVENTORY_RESERVING)
                .event(SagaEvent.ORDER_CREATED)
            .and()
            .withExternal()
                .source(SagaStep.INVENTORY_RESERVING).target(SagaStep.PAYMENT_PROCESSING)
                .event(SagaEvent.STOCK_RESERVED)
            .and()
            .withExternal()
                .source(SagaStep.PAYMENT_PROCESSING).target(SagaStep.SHIPMENT_CREATING)
                .event(SagaEvent.PAYMENT_COMPLETED)
            .and()
            .withExternal()
                .source(SagaStep.SHIPMENT_CREATING).target(SagaStep.COMPLETED)
                .event(SagaEvent.SHIPMENT_CREATED)
            // Compensation transitions
            .and()
            .withExternal()
                .source(SagaStep.PAYMENT_PROCESSING).target(SagaStep.COMPENSATING)
                .event(SagaEvent.PAYMENT_FAILED)
            .and()
            .withExternal()
                .source(SagaStep.INVENTORY_RESERVING).target(SagaStep.COMPENSATING)
                .event(SagaEvent.RESERVE_FAILED);
    }
}
```

---

## 5. Docker & Container Orchestration

### 5.1 Dockerizing Spring Boot

```dockerfile
# Dockerfile - Multi-stage build
FROM eclipse-temurin:21-jdk-alpine AS builder
WORKDIR /app

# Copy maven wrapper and pom
COPY mvnw .
COPY .mvn .mvn
COPY pom.xml .

# Download dependencies (cached layer)
RUN ./mvnw dependency:go-offline -B

# Copy source and build
COPY src src
RUN ./mvnw package -DskipTests

# Extract layers for better caching
RUN java -Djarmode=layertools -jar target/*.jar extract

# Runtime image
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app

# Add non-root user
RUN addgroup -S spring && adduser -S spring -G spring
USER spring:spring

# Copy layers in order of change frequency
COPY --from=builder /app/dependencies/ ./
COPY --from=builder /app/spring-boot-loader/ ./
COPY --from=builder /app/snapshot-dependencies/ ./
COPY --from=builder /app/application/ ./

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s \
  CMD wget -q --spider http://localhost:8080/actuator/health || exit 1

ENTRYPOINT ["java", "org.springframework.boot.loader.launch.JarLauncher"]
```

```yaml
# application.yml for containerized deployment
server:
  port: 8080
  shutdown: graceful

spring:
  lifecycle:
    timeout-per-shutdown-phase: 30s
  datasource:
    url: jdbc:postgresql://${DB_HOST:localhost}:${DB_PORT:5432}/${DB_NAME:orders}
    username: ${DB_USERNAME}
    password: ${DB_PASSWORD}

management:
  endpoints:
    web:
      exposure:
        include: health,info,prometheus
  endpoint:
    health:
      probes:
        enabled: true
      show-details: always
```

### 5.2 Docker Compose for Local Development

```yaml
# docker-compose.yml
version: '3.8'

services:
  # Infrastructure Services
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./init-scripts:/docker-entrypoint-initdb.d
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  kafka:
    image: confluentinc/cp-kafka:7.5.0
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1

  zookeeper:
    image: confluentinc/cp-zookeeper:7.5.0
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181

  # Discovery & Config
  config-server:
    build: ./config-server
    ports:
      - "8888:8888"
    environment:
      SPRING_PROFILES_ACTIVE: native
    volumes:
      - ./config-repo:/config-repo

  discovery-server:
    build: ./discovery-server
    ports:
      - "8761:8761"
    depends_on:
      - config-server
    environment:
      SPRING_CONFIG_IMPORT: configserver:http://config-server:8888

  # API Gateway
  api-gateway:
    build: ./api-gateway
    ports:
      - "8080:8080"
    depends_on:
      - discovery-server
      - config-server
    environment:
      SPRING_CONFIG_IMPORT: configserver:http://config-server:8888
      EUREKA_CLIENT_SERVICEURL_DEFAULTZONE: http://discovery-server:8761/eureka/

  # Business Services
  user-service:
    build: ./user-service
    ports:
      - "8081:8081"
    depends_on:
      postgres:
        condition: service_healthy
      discovery-server:
        condition: service_started
    environment:
      SPRING_CONFIG_IMPORT: configserver:http://config-server:8888
      EUREKA_CLIENT_SERVICEURL_DEFAULTZONE: http://discovery-server:8761/eureka/
      DB_HOST: postgres
      DB_NAME: users

  order-service:
    build: ./order-service
    ports:
      - "8082:8082"
    depends_on:
      postgres:
        condition: service_healthy
      kafka:
        condition: service_started
      discovery-server:
        condition: service_started
    environment:
      SPRING_CONFIG_IMPORT: configserver:http://config-server:8888
      EUREKA_CLIENT_SERVICEURL_DEFAULTZONE: http://discovery-server:8761/eureka/
      DB_HOST: postgres
      DB_NAME: orders
      SPRING_KAFKA_BOOTSTRAP_SERVERS: kafka:29092

  inventory-service:
    build: ./inventory-service
    ports:
      - "8083:8083"
    depends_on:
      - postgres
      - kafka
      - discovery-server
    environment:
      SPRING_CONFIG_IMPORT: configserver:http://config-server:8888
      EUREKA_CLIENT_SERVICEURL_DEFAULTZONE: http://discovery-server:8761/eureka/
      DB_HOST: postgres
      DB_NAME: inventory
      SPRING_KAFKA_BOOTSTRAP_SERVERS: kafka:29092

  # Monitoring
  zipkin:
    image: openzipkin/zipkin
    ports:
      - "9411:9411"

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana_data:/var/lib/grafana

volumes:
  postgres_data:
  grafana_data:

networks:
  default:
    name: microservices-network
```

---

## 6. Kubernetes Deployment

### 6.1 Kubernetes Architecture for Microservices

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         Kubernetes Deployment Architecture                       │
│                                                                                  │
│  ┌─────────────────────────────────────────────────────────────────────────┐    │
│  │                          Kubernetes Cluster                              │    │
│  │                                                                          │    │
│  │  ┌────────────────────────────────────────────────────────────────────┐ │    │
│  │  │                        Ingress Controller                          │ │    │
│  │  │                    (nginx / traefik / istio)                       │ │    │
│  │  └─────────────────────────────┬──────────────────────────────────────┘ │    │
│  │                                │                                         │    │
│  │  ┌─────────────────────────────┼─────────────────────────────────────┐  │    │
│  │  │                  microservices namespace                           │  │    │
│  │  │                                                                    │  │    │
│  │  │   ┌───────────┐   ┌───────────┐   ┌───────────┐   ┌───────────┐   │  │    │
│  │  │   │   API     │   │   User    │   │   Order   │   │ Inventory │   │  │    │
│  │  │   │  Gateway  │   │  Service  │   │  Service  │   │  Service  │   │  │    │
│  │  │   │ (3 pods)  │   │ (3 pods)  │   │ (3 pods)  │   │ (3 pods)  │   │  │    │
│  │  │   └─────┬─────┘   └─────┬─────┘   └─────┬─────┘   └─────┬─────┘   │  │    │
│  │  │         │               │               │               │         │  │    │
│  │  │   ┌─────┴─────┐   ┌─────┴─────┐   ┌─────┴─────┐   ┌─────┴─────┐   │  │    │
│  │  │   │  Service  │   │  Service  │   │  Service  │   │  Service  │   │  │    │
│  │  │   │ (ClusterIP│   │ (ClusterIP│   │ (ClusterIP│   │ (ClusterIP│   │  │    │
│  │  │   └───────────┘   └───────────┘   └───────────┘   └───────────┘   │  │    │
│  │  │                                                                    │  │    │
│  │  │   ┌─────────────────────────────────────────────────────────────┐ │  │    │
│  │  │   │                    ConfigMap / Secrets                      │ │  │    │
│  │  │   └─────────────────────────────────────────────────────────────┘ │  │    │
│  │  └────────────────────────────────────────────────────────────────────┘  │    │
│  │                                                                          │    │
│  │  ┌────────────────────────────────────────────────────────────────────┐ │    │
│  │  │                    infrastructure namespace                        │ │    │
│  │  │   ┌───────────┐   ┌───────────┐   ┌───────────┐   ┌───────────┐   │ │    │
│  │  │   │PostgreSQL │   │   Kafka   │   │   Redis   │   │  Zipkin   │   │ │    │
│  │  │   │StatefulSet│   │StatefulSet│   │StatefulSet│   │Deployment │   │ │    │
│  │  │   └───────────┘   └───────────┘   └───────────┘   └───────────┘   │ │    │
│  │  └────────────────────────────────────────────────────────────────────┘ │    │
│  │                                                                          │    │
│  └──────────────────────────────────────────────────────────────────────────┘    │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Kubernetes Manifests

```yaml
# deployment.yaml - Order Service
apiVersion: apps/v1
kind: Deployment
metadata:
  name: order-service
  namespace: microservices
  labels:
    app: order-service
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: order-service
  template:
    metadata:
      labels:
        app: order-service
        version: v1
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/actuator/prometheus"
    spec:
      serviceAccountName: order-service
      containers:
        - name: order-service
          image: myregistry/order-service:1.0.0
          ports:
            - containerPort: 8080
              name: http
          env:
            - name: SPRING_PROFILES_ACTIVE
              value: "kubernetes"
            - name: DB_HOST
              valueFrom:
                configMapKeyRef:
                  name: order-service-config
                  key: db.host
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: order-service-secrets
                  key: db-password
          resources:
            requests:
              memory: "512Mi"
              cpu: "250m"
            limits:
              memory: "1Gi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /actuator/health/liveness
              port: 8080
            initialDelaySeconds: 60
            periodSeconds: 10
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /actuator/health/readiness
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 5
            failureThreshold: 3
          lifecycle:
            preStop:
              exec:
                command: ["sh", "-c", "sleep 10"]
      terminationGracePeriodSeconds: 60
---
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: order-service
  namespace: microservices
spec:
  selector:
    app: order-service
  ports:
    - port: 80
      targetPort: 8080
      name: http
  type: ClusterIP
---
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: order-service-config
  namespace: microservices
data:
  db.host: "postgres.infrastructure.svc.cluster.local"
  db.name: "orders"
  kafka.bootstrap-servers: "kafka.infrastructure.svc.cluster.local:9092"
  spring.profiles.active: "kubernetes"
---
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: order-service-secrets
  namespace: microservices
type: Opaque
data:
  db-password: cGFzc3dvcmQxMjM=  # base64 encoded
  jwt-secret: c2VjcmV0LWtleS0xMjM=
---
# hpa.yaml - Horizontal Pod Autoscaler
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: order-service-hpa
  namespace: microservices
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: order-service
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
```

### 6.3 Ingress Configuration

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: microservices-ingress
  namespace: microservices
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$2
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - api.example.com
      secretName: api-tls-secret
  rules:
    - host: api.example.com
      http:
        paths:
          - path: /api/users(/|$)(.*)
            pathType: Prefix
            backend:
              service:
                name: user-service
                port:
                  number: 80
          - path: /api/orders(/|$)(.*)
            pathType: Prefix
            backend:
              service:
                name: order-service
                port:
                  number: 80
          - path: /api/inventory(/|$)(.*)
            pathType: Prefix
            backend:
              service:
                name: inventory-service
                port:
                  number: 80
```

### 6.4 Spring Cloud Kubernetes

```yaml
# application-kubernetes.yml
spring:
  application:
    name: order-service
  cloud:
    kubernetes:
      discovery:
        enabled: true
        all-namespaces: false
      config:
        enabled: true
        name: order-service
        namespace: microservices
      reload:
        enabled: true
        mode: polling
        period: 15000
  config:
    import: kubernetes:

management:
  endpoint:
    health:
      probes:
        enabled: true
  health:
    livenessState:
      enabled: true
    readinessState:
      enabled: true
```

```java
// Using Kubernetes discovery instead of Eureka
@SpringBootApplication
@EnableDiscoveryClient
public class OrderServiceApplication {
    public static void main(String[] args) {
        SpringApplication.run(OrderServiceApplication.class, args);
    }
}
```

---

## 7. Testing Microservices

### 7.1 Testing Pyramid for Microservices

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                     Microservices Testing Pyramid                                │
│                                                                                  │
│                              ┌─────────┐                                        │
│                              │   E2E   │  Few, slow, expensive                  │
│                              │  Tests  │  Full system tests                     │
│                             ┌┴─────────┴┐                                       │
│                             │ Component │  Service in isolation                 │
│                             │   Tests   │  with mocked dependencies             │
│                            ┌┴───────────┴┐                                      │
│                            │  Contract   │  API contracts between               │
│                            │   Tests     │  services (Pact, Spring Cloud        │
│                           ┌┴─────────────┴┐ Contract)                           │
│                           │  Integration  │  Database, messaging,               │
│                           │    Tests      │  external services                  │
│                          ┌┴───────────────┴┐                                    │
│                          │    Unit Tests    │  Business logic, fast,           │
│                          │                  │  many tests                       │
│                          └──────────────────┘                                   │
│                                                                                  │
│  Recommended Distribution:                                                       │
│  • Unit Tests: 70%                                                              │
│  • Integration Tests: 20%                                                        │
│  • Contract Tests: 5%                                                           │
│  • Component/E2E: 5%                                                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

### 7.2 Unit Testing

```java
@ExtendWith(MockitoExtension.class)
class OrderServiceTest {
    
    @Mock
    private OrderRepository orderRepository;
    
    @Mock
    private UserClient userClient;
    
    @Mock
    private InventoryClient inventoryClient;
    
    @InjectMocks
    private OrderService orderService;
    
    @Test
    void createOrder_Success() {
        // Given
        CreateOrderRequest request = CreateOrderRequest.builder()
            .userId("user-123")
            .items(List.of(new OrderItem("SKU001", 2)))
            .build();
        
        when(userClient.getUser("user-123"))
            .thenReturn(new User("user-123", "John", "ACTIVE"));
        when(inventoryClient.checkAvailability("SKU001", 2))
            .thenReturn(true);
        when(orderRepository.save(any(Order.class)))
            .thenAnswer(inv -> {
                Order order = inv.getArgument(0);
                order.setId(1L);
                return order;
            });
        
        // When
        Order result = orderService.createOrder(request);
        
        // Then
        assertThat(result.getId()).isEqualTo(1L);
        assertThat(result.getStatus()).isEqualTo(OrderStatus.CREATED);
        verify(orderRepository).save(any(Order.class));
    }
    
    @Test
    void createOrder_UserNotFound_ThrowsException() {
        // Given
        CreateOrderRequest request = CreateOrderRequest.builder()
            .userId("user-123")
            .build();
        
        when(userClient.getUser("user-123"))
            .thenThrow(new UserNotFoundException("User not found"));
        
        // When/Then
        assertThatThrownBy(() -> orderService.createOrder(request))
            .isInstanceOf(UserNotFoundException.class);
        
        verify(orderRepository, never()).save(any());
    }
}
```

### 7.3 Integration Testing

```java
@SpringBootTest
@AutoConfigureTestDatabase(replace = Replace.NONE)
@Testcontainers
class OrderRepositoryIntegrationTest {
    
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:15")
        .withDatabaseName("test")
        .withUsername("test")
        .withPassword("test");
    
    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }
    
    @Autowired
    private OrderRepository orderRepository;
    
    @BeforeEach
    void setUp() {
        orderRepository.deleteAll();
    }
    
    @Test
    void findByUserId_ReturnsUserOrders() {
        // Given
        Order order1 = Order.builder()
            .userId("user-123")
            .status(OrderStatus.CREATED)
            .build();
        Order order2 = Order.builder()
            .userId("user-123")
            .status(OrderStatus.COMPLETED)
            .build();
        orderRepository.saveAll(List.of(order1, order2));
        
        // When
        List<Order> orders = orderRepository.findByUserId("user-123");
        
        // Then
        assertThat(orders).hasSize(2);
    }
}

// Kafka Integration Test
@SpringBootTest
@EmbeddedKafka(partitions = 1, topics = {"orders"})
class OrderEventIntegrationTest {
    
    @Autowired
    private EmbeddedKafkaBroker embeddedKafka;
    
    @Autowired
    private KafkaTemplate<String, OrderCreatedEvent> kafkaTemplate;
    
    @Autowired
    private OrderEventConsumer orderEventConsumer;
    
    @Test
    void consumeOrderCreatedEvent() throws Exception {
        // Given
        OrderCreatedEvent event = OrderCreatedEvent.builder()
            .orderId(1L)
            .userId("user-123")
            .build();
        
        // When
        kafkaTemplate.send("orders", event).get();
        
        // Then
        // Verify consumer processed the event
        await().atMost(10, TimeUnit.SECONDS)
            .untilAsserted(() -> 
                verify(inventoryService).reserveStock(anyList()));
    }
}
```

### 7.4 Contract Testing with Spring Cloud Contract

```groovy
// contract/shouldCreateOrder.groovy (Producer side)
Contract.make {
    description "should create order"
    
    request {
        method POST()
        url "/api/orders"
        headers {
            contentType applicationJson()
            header("Authorization", "Bearer token")
        }
        body([
            userId: "user-123",
            items: [
                [sku: "SKU001", quantity: 2]
            ]
        ])
    }
    
    response {
        status CREATED()
        headers {
            contentType applicationJson()
        }
        body([
            id: $(producer(regex('[0-9]+')), consumer(1)),
            userId: "user-123",
            status: "CREATED",
            items: [
                [sku: "SKU001", quantity: 2]
            ]
        ])
    }
}
```

```java
// Consumer side test
@SpringBootTest
@AutoConfigureStubRunner(
    ids = "com.example:order-service:+:stubs:8082",
    stubsMode = StubsMode.LOCAL
)
class OrderClientContractTest {
    
    @Autowired
    private OrderClient orderClient;
    
    @Test
    void createOrder_ContractVerified() {
        // Given
        CreateOrderRequest request = new CreateOrderRequest(
            "user-123",
            List.of(new OrderItem("SKU001", 2))
        );
        
        // When
        Order result = orderClient.createOrder(request);
        
        // Then
        assertThat(result.getId()).isNotNull();
        assertThat(result.getStatus()).isEqualTo("CREATED");
    }
}
```

### 7.5 Component Testing

```java
@SpringBootTest(webEnvironment = WebEnvironment.RANDOM_PORT)
@AutoConfigureWireMock(port = 0)
class OrderServiceComponentTest {
    
    @Autowired
    private TestRestTemplate restTemplate;
    
    @LocalServerPort
    private int port;
    
    @BeforeEach
    void setUp() {
        // Mock external services
        stubFor(get(urlEqualTo("/api/users/user-123"))
            .willReturn(aResponse()
                .withStatus(200)
                .withHeader("Content-Type", "application/json")
                .withBody("""
                    {"id": "user-123", "name": "John", "status": "ACTIVE"}
                    """)));
        
        stubFor(post(urlEqualTo("/api/inventory/reserve"))
            .willReturn(aResponse()
                .withStatus(200)));
    }
    
    @Test
    void createOrder_EndToEnd() {
        // Given
        String requestBody = """
            {
                "userId": "user-123",
                "items": [{"sku": "SKU001", "quantity": 2}]
            }
            """;
        
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_JSON);
        headers.setBearerAuth("test-token");
        
        HttpEntity<String> request = new HttpEntity<>(requestBody, headers);
        
        // When
        ResponseEntity<Order> response = restTemplate.postForEntity(
            "/api/orders", request, Order.class);
        
        // Then
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.CREATED);
        assertThat(response.getBody().getStatus()).isEqualTo(OrderStatus.CREATED);
        
        // Verify external calls
        verify(1, getRequestedFor(urlEqualTo("/api/users/user-123")));
    }
}
```

---

## 8. Interview Questions - Part 2

### Distributed Systems

**Q1: How do you implement distributed tracing across microservices?**
> Use Micrometer Tracing with a backend like Zipkin or Jaeger. Each request gets a unique Trace ID that propagates through all services. Spans capture individual operations. Configure sampling rate for production (typically 10-20%). Include trace IDs in logs for correlation.

**Q2: Explain the Circuit Breaker pattern vs Retry pattern. When to use each?**
> - **Retry**: For transient failures (network blips). Retries immediately or with backoff.
> - **Circuit Breaker**: For persistent failures. Fails fast when a service is down, preventing resource exhaustion.
> Use them together: Retry first, Circuit Breaker when retries exceed threshold.

**Q3: How do you handle distributed transactions in microservices?**
> Avoid 2PC (two-phase commit) in microservices. Use:
> - **Saga Pattern**: Orchestration (central coordinator) or Choreography (event-driven)
> - Compensating transactions for rollback
> - Eventual consistency with idempotent operations
> - Outbox pattern for reliable event publishing

### Security

**Q4: How do you secure service-to-service communication?**
> Options:
> - **Token Relay**: Propagate user's JWT token
> - **Client Credentials**: Service account tokens (OAuth2)
> - **mTLS**: Mutual TLS for transport security
> - **Service Mesh**: Istio/Linkerd handles mTLS automatically

**Q5: Explain JWT token validation in a microservices architecture.**
> - Gateway validates token signature using public key from auth server
> - Token claims (roles, scopes) extracted for authorization
> - Token can be validated locally (no auth server call) using JWK
> - Consider token caching for performance

### Event-Driven

**Q6: What is the Outbox pattern and why is it needed?**
> The Outbox pattern ensures atomicity between database updates and event publishing. Save events to an "outbox" table in the same transaction as business data. A separate process (or Debezium CDC) reads and publishes events. This prevents "dual-write" problems where DB succeeds but message publish fails.

**Q7: Compare Kafka vs RabbitMQ for microservices.**

| Aspect | Kafka | RabbitMQ |
|--------|-------|----------|
| Model | Log-based, pull | Message queue, push |
| Ordering | Per-partition | Per-queue |
| Retention | Configurable (days/weeks) | Until consumed |
| Throughput | Very high (millions/sec) | High (thousands/sec) |
| Use Case | Event streaming, analytics | Task queues, RPC |

### Kubernetes

**Q8: How do you handle configuration in Kubernetes microservices?**
> - **ConfigMaps**: Non-sensitive configuration
> - **Secrets**: Sensitive data (encrypted at rest)
> - **Spring Cloud Kubernetes Config**: Auto-reload on changes
> - External tools: Vault, AWS Secrets Manager
> - Avoid hardcoding in images

**Q9: Explain liveness vs readiness probes.**
> - **Liveness**: Is the container alive? Failure triggers restart.
> - **Readiness**: Can the container serve traffic? Failure removes from service endpoints.
> Configure: Liveness with longer initial delay, Readiness more aggressive.

**Q10: How do you implement zero-downtime deployments?**
> - Use Rolling Updates or Blue-Green deployments
> - Configure proper readiness probes
> - Implement graceful shutdown (preStop hook, `spring.lifecycle.timeout-per-shutdown-phase`)
> - Use PodDisruptionBudget to ensure minimum replicas
> - Database migrations: backward-compatible changes first

### Advanced Patterns

**Q11: Design a saga for order processing that handles all failure scenarios.**
```
Order Saga Steps:
1. Create Order (pending)
2. Reserve Inventory → Compensate: Release Inventory
3. Process Payment → Compensate: Refund Payment
4. Create Shipment → Compensate: Cancel Shipment
5. Confirm Order

Failure Handling:
- Each step stores state for compensation
- Timeout handling with dead-letter queues
- Idempotency keys prevent duplicate processing
- Saga state machine tracks progress
```

**Q12: How would you implement rate limiting at scale?**
> - **Centralized**: Redis-based token bucket (Spring Cloud Gateway)
> - **Distributed**: Consistent hashing for rate limit state
> - **Strategies**: Per-user, per-IP, per-API key
> - Consider: Sliding window for smoother limits, burst capacity

**Q13: Explain the strangler fig pattern for microservices migration.**
> Gradually replace monolith with microservices:
> 1. Add facade/proxy in front of monolith
> 2. Extract one bounded context to new service
> 3. Route specific traffic to new service
> 4. Repeat until monolith is hollow
> 5. Decommission monolith
> 
> Key: Never big-bang rewrite. Always incremental.

**Q14: How do you implement idempotency across distributed services?**
```java
// Strategy 1: Idempotency Key
@PostMapping("/orders")
public Order create(
    @RequestHeader("Idempotency-Key") String key,
    @RequestBody OrderRequest req) {
    
    // Check cache/DB for existing result
    Order existing = idempotencyStore.get(key);
    if (existing != null) return existing;
    
    // Process and store result
    Order order = orderService.create(req);
    idempotencyStore.put(key, order, TTL_24_HOURS);
    return order;
}

// Strategy 2: Event Deduplication
@KafkaListener(topics = "orders")
public void handle(OrderEvent event) {
    if (processedEvents.contains(event.getEventId())) {
        return; // Already processed
    }
    // Process event
    processedEvents.add(event.getEventId());
}
```

**Q15: How do you debug issues in a distributed system?**
> 1. **Distributed Tracing**: Follow request across services
> 2. **Centralized Logging**: ELK/Loki with correlation IDs
> 3. **Metrics**: Prometheus + Grafana for anomaly detection
> 4. **Service Mesh**: Kiali for traffic visualization
> 5. **Chaos Engineering**: Identify weaknesses before production issues

---

## Quick Reference - Part 2

### Key Annotations

| Annotation | Purpose |
|------------|---------|
| `@Observed` | Add tracing span to method |
| `@PreAuthorize` | Method-level authorization |
| `@Transactional` | Local transaction boundary |
| `@KafkaListener` | Consume Kafka messages |
| `@StreamListener` | Spring Cloud Stream consumer |

### Essential Configuration

```yaml
# Tracing
management.tracing.sampling.probability: 0.1  # 10% sampling

# Security
spring.security.oauth2.resourceserver.jwt.issuer-uri: ${AUTH_SERVER}

# Kafka
spring.kafka.consumer.group-id: ${spring.application.name}
spring.kafka.consumer.auto-offset-reset: earliest

# Kubernetes probes
management.endpoint.health.probes.enabled: true
```

### Saga State Transitions

```
STARTED → ORDER_CREATED → INVENTORY_RESERVED → PAYMENT_COMPLETED → COMPLETED
    ↓            ↓                ↓                    ↓
COMPENSATING ← COMPENSATING ← COMPENSATING ← COMPENSATING → COMPENSATED
```

---

*This completes the Microservices & Spring Cloud guide. Combined with Part 1, you have comprehensive coverage for senior-level interviews.*

*Last Updated: February 2026*
