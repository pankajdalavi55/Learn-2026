# Spring Security — Complete Learning Guide

> A comprehensive, interview-ready reference covering Spring Security from fundamentals to advanced enterprise patterns.  
> **Scope:** Authentication, Authorization, OAuth2, JWT, CSRF/CORS, Method Security, and production hardening.  
> **Prerequisites:** Familiarity with Spring Boot fundamentals (see companion guides).

---

## Table of Contents

1. [Introduction & Security Fundamentals](#1-introduction--security-fundamentals)
2. [Spring Security Architecture](#2-spring-security-architecture)
3. [Authentication Deep Dive](#3-authentication-deep-dive)
4. [Authorization & Access Control](#4-authorization--access-control)
5. [Password Storage & Encoding](#5-password-storage--encoding)
6. [Session Management](#6-session-management)
7. [JWT Authentication](#7-jwt-authentication)
8. [OAuth2 & OpenID Connect](#8-oauth2--openid-connect)
9. [Method-Level Security](#9-method-level-security)
10. [CORS, CSRF & Security Headers](#10-cors-csrf--security-headers)
11. [Security Testing](#11-security-testing)
12. [Common Vulnerabilities & Prevention](#12-common-vulnerabilities--prevention)
13. [Production Security Checklist](#13-production-security-checklist)
14. [Interview Questions by Experience Level](#14-interview-questions-by-experience-level)

---

## 1. Introduction & Security Fundamentals

### 1.1 What is Application Security?

**Security** is the practice of protecting systems, networks, and data from unauthorized access, attacks, and damage. In web applications, security encompasses:

- **Confidentiality** — Only authorized users can access data
- **Integrity** — Data cannot be tampered with undetected
- **Availability** — System remains accessible to legitimate users
- **Non-repudiation** — Actions can be traced to their source

### 1.1.1 The CIA Triad

```
                         ┌─────────────────┐
                         │ CONFIDENTIALITY │
                         │                 │
                         │  "Only right    │
                         │   people see    │
                         │   the data"     │
                         └────────┬────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
              │            ┌──────┴──────┐            │
              │            │             │            │
              │            │   SECURE    │            │
              │            │   SYSTEM    │            │
              │            │             │            │
              │            └─────────────┘            │
              │                                       │
     ┌────────┴────────┐                   ┌─────────┴───────┐
     │   INTEGRITY     │                   │  AVAILABILITY   │
     │                 │                   │                 │
     │ "Data hasn't    │                   │ "System is      │
     │  been altered"  │                   │  accessible"    │
     └─────────────────┘                   └─────────────────┘
```

### 1.1.2 Authentication vs Authorization

**Critical Distinction — Interview Favorite!**

```
┌─────────────────────────────────────────────────────────────────────────┐
│              AUTHENTICATION vs AUTHORIZATION                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  AUTHENTICATION (AuthN)              AUTHORIZATION (AuthZ)              │
│  ─────────────────────               ──────────────────────             │
│                                                                         │
│  "WHO are you?"                      "WHAT can you do?"                 │
│                                                                         │
│  Verifies IDENTITY                   Verifies PERMISSIONS               │
│                                                                         │
│  Happens FIRST                       Happens AFTER authentication       │
│                                                                         │
│  Methods:                            Methods:                           │
│  • Username/Password                 • Role-based (RBAC)                │
│  • Biometrics                        • Attribute-based (ABAC)           │
│  • Certificates                      • Permission-based                 │
│  • OAuth2 tokens                     • ACL (Access Control Lists)       │
│                                                                         │
│  Result:                             Result:                            │
│  Principal (user identity)           Granted/Denied access              │
│                                                                         │
│  Analogy:                            Analogy:                           │
│  Showing ID at airport               Boarding pass for specific flight  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.1.3 Security Principles

| Principle | Description | Spring Security Implementation |
|---|---|---|
| **Defense in Depth** | Multiple layers of security | Filter chain + method security + view security |
| **Least Privilege** | Minimum access needed | Role hierarchies, fine-grained permissions |
| **Fail Secure** | Deny by default | `denyAll()` default, explicit allows |
| **Separation of Duties** | No single point of control | Different roles for different operations |
| **Open Design** | Security shouldn't rely on secrecy | Well-documented, auditable mechanisms |

### 1.1.4 Security Mental Model — Thinking Like an Attacker

**To build secure systems, you must understand how attackers think.**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    ATTACKER'S PERSPECTIVE                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  RECONNAISSANCE          "What does the target look like?"              │
│  ─────────────           • Technology stack detection                   │
│                          • Version fingerprinting                       │
│                          • Endpoint enumeration                         │
│                          • Public info gathering (OSINT)                │
│                                                                         │
│  ATTACK SURFACE          "Where can I get in?"                          │
│  ─────────────           • Authentication endpoints                     │
│                          • File upload functionality                    │
│                          • Input fields (SQL, XSS vectors)              │
│                          • Third-party integrations                     │
│                          • API endpoints                                │
│                                                                         │
│  EXPLOITATION            "How do I break in?"                           │
│  ────────────            • Credential stuffing/brute force              │
│                          • Injection attacks                            │
│                          • Session hijacking                            │
│                          • Token theft                                  │
│                                                                         │
│  PERSISTENCE             "How do I stay in?"                            │
│  ───────────             • Backdoor creation                            │
│                          • Privilege escalation                         │
│                          • Data exfiltration                            │
│                                                                         │
│  COVERING TRACKS         "How do I avoid detection?"                    │
│  ──────────────          • Log manipulation                             │
│                          • Timestamp modification                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Defender's Response** — Map each attack phase to security controls:

| Attack Phase | Spring Security Control |
|---|---|
| Reconnaissance | Hide version headers, custom error pages |
| Attack Surface | Minimize exposed endpoints, input validation |
| Exploitation | Strong auth, CSRF/XSS protection, rate limiting |
| Persistence | Session management, token revocation |
| Covering Tracks | Comprehensive audit logging |

### 1.1.5 Threat Modeling — STRIDE Framework

**STRIDE** is a threat modeling framework developed by Microsoft to systematically identify security threats.

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    STRIDE THREAT MODEL                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  S — SPOOFING IDENTITY                                                  │
│      ─────────────────                                                  │
│      Threat: Pretending to be someone else                              │
│      Example: Using stolen credentials                                  │
│      Mitigation: Strong authentication (MFA, certificates)             │
│      Spring: AuthenticationManager, OAuth2, JWT validation              │
│                                                                         │
│  T — TAMPERING WITH DATA                                                │
│      ──────────────────                                                 │
│      Threat: Modifying data in transit or at rest                       │
│      Example: Man-in-the-middle attack, SQL injection                   │
│      Mitigation: Integrity checks, input validation, HTTPS              │
│      Spring: CSRF tokens, signed JWTs, TLS                              │
│                                                                         │
│  R — REPUDIATION                                                        │
│      ───────────                                                        │
│      Threat: Denying performed actions                                  │
│      Example: "I never made that transaction"                           │
│      Mitigation: Audit logs, digital signatures                         │
│      Spring: Security event logging, signed tokens                      │
│                                                                         │
│  I — INFORMATION DISCLOSURE                                             │
│      ───────────────────────                                            │
│      Threat: Exposing data to unauthorized parties                      │
│      Example: Stack traces in errors, verbose logging                   │
│      Mitigation: Encryption, access control, error handling             │
│      Spring: Authorization rules, custom error handlers                 │
│                                                                         │
│  D — DENIAL OF SERVICE                                                  │
│      ─────────────────                                                  │
│      Threat: Making system unavailable                                  │
│      Example: Resource exhaustion, infinite loops                       │
│      Mitigation: Rate limiting, resource quotas, timeouts               │
│      Spring: Filters, circuit breakers (Resilience4j)                   │
│                                                                         │
│  E — ELEVATION OF PRIVILEGE                                             │
│      ────────────────────────                                           │
│      Threat: Gaining unauthorized access levels                         │
│      Example: Accessing admin functions as regular user                 │
│      Mitigation: Proper authorization, least privilege                  │
│      Spring: @PreAuthorize, role hierarchies, method security           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**How to Apply STRIDE:**

1. **Identify assets** — What are you protecting? (user data, API keys, etc.)
2. **Create data flow diagrams** — How does data move through the system?
3. **For each component/data flow** — Apply STRIDE categories
4. **Rate threats** — Use DREAD (Damage, Reproducibility, Exploitability, Affected users, Discoverability)
5. **Plan mitigations** — Map to security controls

### 1.2 What is Spring Security?

Spring Security is a **powerful, highly customizable authentication and access-control framework** for Java applications. It is the de-facto standard for securing Spring-based applications.

**Key Features:**
- Comprehensive authentication support
- Protection against common attacks (CSRF, session fixation, clickjacking)
- Servlet API integration
- Optional integration with Spring MVC
- Extensible and customizable

### 1.2.1 Spring Security History

| Version | Year | Key Features |
|---|---|---|
| Acegi Security | 2003 | Original project, complex XML config |
| Spring Security 2.0 | 2008 | Rebranded, simplified namespace config |
| Spring Security 3.0 | 2010 | Expression-based access control, SpEL |
| Spring Security 4.0 | 2015 | Java config, WebSocket security |
| Spring Security 5.0 | 2017 | OAuth2/OIDC, reactive support, password encoding |
| Spring Security 6.0 | 2022 | Jakarta EE 9+, new SecurityFilterChain DSL |

### 1.2.2 Zero Trust Architecture — Theory

**"Never trust, always verify"** — The fundamental shift in security thinking.

```
┌─────────────────────────────────────────────────────────────────────────┐
│            PERIMETER SECURITY vs ZERO TRUST                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  TRADITIONAL (Castle & Moat)         ZERO TRUST                         │
│  ───────────────────────────         ──────────                         │
│                                                                         │
│       ┌─────────────────┐            ┌─────────────────┐                │
│       │   ╔═══════════╗ │            │  🔒  🔒  🔒  🔒  │                │
│   🏰  │   ║  TRUSTED  ║ │            │  🔒  🔒  🔒  🔒  │                │
│  ═════│   ║  NETWORK  ║ │            │  🔒  🔒  🔒  🔒  │                │
│ MOAT  │   ╚═══════════╝ │            │  🔒  🔒  🔒  🔒  │                │
│       └─────────────────┘            └─────────────────┘                │
│                                                                         │
│  • Trust inside firewall             • Trust nothing by default         │
│  • Verify at perimeter               • Verify every request             │
│  • Static access rules               • Dynamic, context-aware           │
│  • VPN = trusted                     • Identity = perimeter             │
│                                                                         │
│  Problem: Once inside,               Solution: Every resource           │
│  attacker has free reign             protected individually             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Zero Trust Principles:**

| Principle | Description | Implementation |
|---|---|---|
| **Verify Explicitly** | Always authenticate and authorize | JWT/OAuth2 on every request |
| **Least Privilege** | Just-in-time, just-enough access | Fine-grained roles, time-limited tokens |
| **Assume Breach** | Design as if attacker is already inside | Micro-segmentation, encryption everywhere |
| **Never Trust Network** | Internal ≠ trusted | mTLS between services |
| **Continuous Validation** | Re-verify throughout session | Token refresh, step-up auth |

**Zero Trust in Spring Security:**

```java
// Traditional: Trust internal requests
.requestMatchers("/internal/**").permitAll()  // ❌ DANGEROUS

// Zero Trust: Verify everything
.requestMatchers("/internal/**").authenticated()  // ✅ Better
.requestMatchers("/internal/**").access(internalServiceAuth())  // ✅ Best

// Implement service-to-service authentication
@Bean
public AuthorizationManager<RequestAuthorizationContext> internalServiceAuth() {
    return (authentication, context) -> {
        // Verify service identity (mTLS cert, service token)
        // Check allowed services list
        // Validate request context
        return new AuthorizationDecision(isValidServiceCall(context));
    };
}
```

### 1.2.3 Defense in Depth — Layered Security

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    DEFENSE IN DEPTH LAYERS                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│                    ┌─────────────────────────┐                          │
│                    │       PERIMETER         │ WAF, DDoS protection     │
│                    │  ┌───────────────────┐  │                          │
│                    │  │     NETWORK       │  │ Firewalls, VPNs, TLS     │
│                    │  │  ┌─────────────┐  │  │                          │
│                    │  │  │   HOST      │  │  │ OS hardening, patches    │
│                    │  │  │  ┌───────┐  │  │  │                          │
│                    │  │  │  │ APP   │  │  │  │ Spring Security          │
│                    │  │  │  │┌─────┐│  │  │  │                          │
│                    │  │  │  ││DATA ││  │  │  │ Encryption, ACLs         │
│                    │  │  │  │└─────┘│  │  │  │                          │
│                    │  │  │  └───────┘  │  │  │                          │
│                    │  │  └─────────────┘  │  │                          │
│                    │  └───────────────────┘  │                          │
│                    └─────────────────────────┘                          │
│                                                                         │
│  Spring Security operates at APPLICATION layer:                         │
│  • Filter Chain         → Request interception                          │
│  • Authentication       → Identity verification                         │
│  • Authorization        → Access control                                │
│  • Method Security      → Fine-grained control                          │
│  • Exception Handling   → Secure error responses                        │
│                                                                         │
│  But security requires ALL layers working together!                     │
└─────────────────────────────────────────────────────────────────────────┘
```

### 1.3 Getting Started

```xml
<dependency>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-starter-security</artifactId>
</dependency>
```

**What Happens Immediately:**
1. All endpoints require authentication
2. Default login form at `/login`
3. Default logout at `/logout`
4. Random password generated (check console logs)
5. CSRF protection enabled
6. Security headers added

```
// Console output on startup
Using generated security password: 8a7b3f2e-1234-5678-abcd-ef1234567890
```

---

## 2. Spring Security Architecture

### 2.1 The Big Picture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    SPRING SECURITY ARCHITECTURE                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   HTTP Request                                                          │
│        │                                                                │
│        ▼                                                                │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                  DelegatingFilterProxy                          │   │
│   │            (Servlet Filter → Spring Bean bridge)                │   │
│   └─────────────────────────────┬───────────────────────────────────┘   │
│                                 │                                       │
│                                 ▼                                       │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │                 FilterChainProxy                                │   │
│   │         (Manages multiple SecurityFilterChains)                 │   │
│   └─────────────────────────────┬───────────────────────────────────┘   │
│                                 │                                       │
│                                 ▼                                       │
│   ┌─────────────────────────────────────────────────────────────────┐   │
│   │               SecurityFilterChain                               │   │
│   │  ┌──────────────────────────────────────────────────────────┐   │   │
│   │  │ SecurityContextPersistenceFilter                         │   │   │
│   │  │ HeaderWriterFilter                                       │   │   │
│   │  │ CsrfFilter                                               │   │   │
│   │  │ LogoutFilter                                             │   │   │
│   │  │ UsernamePasswordAuthenticationFilter                     │   │   │
│   │  │ BasicAuthenticationFilter                                │   │   │
│   │  │ BearerTokenAuthenticationFilter (JWT)                    │   │   │
│   │  │ RequestCacheAwareFilter                                  │   │   │
│   │  │ SecurityContextHolderAwareRequestFilter                  │   │   │
│   │  │ AnonymousAuthenticationFilter                            │   │   │
│   │  │ SessionManagementFilter                                  │   │   │
│   │  │ ExceptionTranslationFilter                               │   │   │
│   │  │ AuthorizationFilter (was FilterSecurityInterceptor)      │   │   │
│   │  └──────────────────────────────────────────────────────────┘   │   │
│   └─────────────────────────────┬───────────────────────────────────┘   │
│                                 │                                       │
│                                 ▼                                       │
│                          Your Controller                                │
└─────────────────────────────────────────────────────────────────────────┘
```

### 2.2 Key Components Explained

#### 2.2.1 DelegatingFilterProxy

**Bridge between Servlet container and Spring's ApplicationContext.**

```java
// The servlet container knows about DelegatingFilterProxy
// DelegatingFilterProxy delegates to a Spring bean named "springSecurityFilterChain"

// In web.xml (traditional)
<filter>
    <filter-name>springSecurityFilterChain</filter-name>
    <filter-class>org.springframework.web.filter.DelegatingFilterProxy</filter-class>
</filter>

// In Spring Boot — Auto-configured!
```

#### 2.2.2 FilterChainProxy

**Routes requests to the appropriate SecurityFilterChain.**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                      FilterChainProxy                                   │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Request: /api/users/123                                                │
│                    │                                                    │
│                    ▼                                                    │
│  ┌─────────────────────────────────────────┐                            │
│  │ Match: /api/** → SecurityFilterChain #1 │ ← JWT/Bearer Token         │
│  └─────────────────────────────────────────┘                            │
│                                                                         │
│  Request: /admin/dashboard                                              │
│                    │                                                    │
│                    ▼                                                    │
│  ┌─────────────────────────────────────────┐                            │
│  │ Match: /admin/** → SecurityFilterChain #2│ ← Form login, CSRF        │
│  └─────────────────────────────────────────┘                            │
│                                                                         │
│  Multiple filter chains for different URL patterns!                     │
└─────────────────────────────────────────────────────────────────────────┘
```

#### 2.2.3 SecurityFilterChain

```java
@Configuration
@EnableWebSecurity
public class SecurityConfig {

    @Bean
    @Order(1)
    public SecurityFilterChain apiFilterChain(HttpSecurity http) throws Exception {
        return http
            .securityMatcher("/api/**")
            .authorizeHttpRequests(auth -> auth.anyRequest().authenticated())
            .sessionManagement(session -> session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .oauth2ResourceServer(oauth2 -> oauth2.jwt(Customizer.withDefaults()))
            .build();
    }

    @Bean
    @Order(2)
    public SecurityFilterChain webFilterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/", "/public/**").permitAll()
                .requestMatchers("/admin/**").hasRole("ADMIN")
                .anyRequest().authenticated()
            )
            .formLogin(Customizer.withDefaults())
            .build();
    }
}
```

### 2.3 Security Filter Order

| Order | Filter | Purpose |
|---|---|---|
| 1 | DisableEncodeUrlFilter | Prevents session ID in URLs |
| 2 | WebAsyncManagerIntegrationFilter | Async request support |
| 3 | SecurityContextHolderFilter | Establishes SecurityContext |
| 4 | HeaderWriterFilter | Adds security headers |
| 5 | CorsFilter | CORS handling |
| 6 | CsrfFilter | CSRF protection |
| 7 | LogoutFilter | Logout handling |
| 8 | OAuth2AuthorizationRequestRedirectFilter | OAuth2 redirect |
| 9 | UsernamePasswordAuthenticationFilter | Form login |
| 10 | BasicAuthenticationFilter | HTTP Basic auth |
| 11 | BearerTokenAuthenticationFilter | JWT/Bearer token |
| 12 | RequestCacheAwareFilter | Caches unauthenticated requests |
| 13 | SecurityContextHolderAwareRequestFilter | Servlet API integration |
| 14 | AnonymousAuthenticationFilter | Creates anonymous principal |
| 15 | ExceptionTranslationFilter | Translates exceptions |
| 16 | AuthorizationFilter | Authorization checks |

### 2.4 SecurityContext and SecurityContextHolder

```java
// SecurityContext holds the Authentication object
public interface SecurityContext {
    Authentication getAuthentication();
    void setAuthentication(Authentication authentication);
}

// SecurityContextHolder provides access to SecurityContext
// Default strategy: ThreadLocal (per-thread storage)

// Getting current user
Authentication auth = SecurityContextHolder.getContext().getAuthentication();
String username = auth.getName();
Collection<? extends GrantedAuthority> authorities = auth.getAuthorities();

// In a controller — Spring MVC integration
@GetMapping("/me")
public UserDto getCurrentUser(@AuthenticationPrincipal UserDetails user) {
    return new UserDto(user.getUsername(), user.getAuthorities());
}

// Or via Principal
@GetMapping("/me")
public String getCurrentUser(Principal principal) {
    return principal.getName();
}
```

**SecurityContextHolder Strategies:**

| Strategy | Use Case |
|---|---|
| `MODE_THREADLOCAL` | Default, single-threaded requests |
| `MODE_INHERITABLETHREADLOCAL` | Spawn child threads with same context |
| `MODE_GLOBAL` | Standalone apps, single security context |

---

## 3. Authentication Deep Dive

### 3.1 Authentication Flow

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AUTHENTICATION FLOW                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. User submits credentials                                            │
│     │                                                                   │
│     ▼                                                                   │
│  ┌─────────────────────────────────────┐                                │
│  │ AuthenticationFilter                │ Creates Authentication token   │
│  │ (e.g., UsernamePasswordAuthFilter)  │ (unauthenticated)              │
│  └─────────────────┬───────────────────┘                                │
│                    │                                                    │
│     2. Delegates   ▼                                                    │
│  ┌─────────────────────────────────────┐                                │
│  │ AuthenticationManager               │ Usually ProviderManager        │
│  └─────────────────┬───────────────────┘                                │
│                    │                                                    │
│     3. Iterates    ▼                                                    │
│  ┌─────────────────────────────────────┐                                │
│  │ AuthenticationProvider              │ e.g., DaoAuthenticationProvider│
│  │ (one or more)                       │                                │
│  └─────────────────┬───────────────────┘                                │
│                    │                                                    │
│     4. Loads       ▼                                                    │
│  ┌─────────────────────────────────────┐                                │
│  │ UserDetailsService                  │ Your custom implementation     │
│  │ → loadUserByUsername(username)      │                                │
│  └─────────────────┬───────────────────┘                                │
│                    │                                                    │
│     5. Returns     ▼                                                    │
│  ┌─────────────────────────────────────┐                                │
│  │ UserDetails                         │ (username, password,           │
│  │                                     │  authorities, enabled, etc.)   │
│  └─────────────────┬───────────────────┘                                │
│                    │                                                    │
│     6. Verifies    ▼                                                    │
│  ┌─────────────────────────────────────┐                                │
│  │ PasswordEncoder.matches()           │ Compares submitted vs stored   │
│  └─────────────────┬───────────────────┘                                │
│                    │                                                    │
│     7. Success!    ▼                                                    │
│  ┌─────────────────────────────────────┐                                │
│  │ Authentication (authenticated=true) │ Stored in SecurityContext     │
│  │ Contains: Principal, Credentials,   │                                │
│  │           Authorities               │                                │
│  └─────────────────────────────────────┘                                │
└─────────────────────────────────────────────────────────────────────────┘
```

### 3.2 Core Authentication Interfaces

```java
// The Authentication object
public interface Authentication extends Principal {
    Collection<? extends GrantedAuthority> getAuthorities();  // Roles/permissions
    Object getCredentials();                                    // Password (cleared after auth)
    Object getDetails();                                        // Extra info (IP, session)
    Object getPrincipal();                                      // User identity
    boolean isAuthenticated();                                  // Auth status
    void setAuthenticated(boolean isAuthenticated);
}

// UserDetails — your user representation
public interface UserDetails {
    Collection<? extends GrantedAuthority> getAuthorities();
    String getPassword();
    String getUsername();
    boolean isAccountNonExpired();
    boolean isAccountNonLocked();
    boolean isCredentialsNonExpired();
    boolean isEnabled();
}

// UserDetailsService — loads users
public interface UserDetailsService {
    UserDetails loadUserByUsername(String username) throws UsernameNotFoundException;
}

// AuthenticationProvider — performs authentication
public interface AuthenticationProvider {
    Authentication authenticate(Authentication authentication) throws AuthenticationException;
    boolean supports(Class<?> authentication);
}
```

### 3.3 Implementing UserDetailsService

```java
@Service
public class CustomUserDetailsService implements UserDetailsService {

    private final UserRepository userRepository;

    public CustomUserDetailsService(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    @Override
    public UserDetails loadUserByUsername(String username) throws UsernameNotFoundException {
        User user = userRepository.findByEmail(username)
            .orElseThrow(() -> new UsernameNotFoundException("User not found: " + username));

        return org.springframework.security.core.userdetails.User.builder()
            .username(user.getEmail())
            .password(user.getPasswordHash())
            .authorities(user.getRoles().stream()
                .map(role -> new SimpleGrantedAuthority("ROLE_" + role.getName()))
                .toList())
            .accountExpired(false)
            .accountLocked(user.isLocked())
            .credentialsExpired(false)
            .disabled(!user.isActive())
            .build();
    }
}
```

### 3.4 Custom UserDetails Implementation

```java
public class CustomUserDetails implements UserDetails {

    private final User user;
    private final Collection<? extends GrantedAuthority> authorities;

    public CustomUserDetails(User user) {
        this.user = user;
        this.authorities = user.getRoles().stream()
            .flatMap(role -> {
                List<GrantedAuthority> auths = new ArrayList<>();
                auths.add(new SimpleGrantedAuthority("ROLE_" + role.getName()));
                role.getPermissions().forEach(perm ->
                    auths.add(new SimpleGrantedAuthority(perm.getName()))
                );
                return auths.stream();
            })
            .toList();
    }

    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() {
        return authorities;
    }

    @Override
    public String getPassword() { return user.getPasswordHash(); }

    @Override
    public String getUsername() { return user.getEmail(); }

    @Override
    public boolean isAccountNonExpired() { return true; }

    @Override
    public boolean isAccountNonLocked() { return !user.isLocked(); }

    @Override
    public boolean isCredentialsNonExpired() { return true; }

    @Override
    public boolean isEnabled() { return user.isActive(); }

    // Extra methods to access your domain user
    public Long getUserId() { return user.getId(); }
    public String getFullName() { return user.getFullName(); }
    public User getUser() { return user; }
}
```

### 3.5 Authentication Providers

#### 3.5.1 DaoAuthenticationProvider (Default)

```java
@Configuration
public class AuthenticationConfig {

    @Bean
    public AuthenticationManager authenticationManager(
            UserDetailsService userDetailsService,
            PasswordEncoder passwordEncoder) {
        
        DaoAuthenticationProvider provider = new DaoAuthenticationProvider();
        provider.setUserDetailsService(userDetailsService);
        provider.setPasswordEncoder(passwordEncoder);
        
        return new ProviderManager(provider);
    }
}
```

#### 3.5.2 Multiple Authentication Providers

```java
@Configuration
public class MultiAuthConfig {

    @Bean
    public AuthenticationManager authenticationManager(
            DaoAuthenticationProvider daoProvider,
            LdapAuthenticationProvider ldapProvider,
            JwtAuthenticationProvider jwtProvider) {
        
        // Tries each provider in order until one succeeds
        return new ProviderManager(
            daoProvider,    // Try database first
            ldapProvider,   // Then LDAP
            jwtProvider     // Then JWT
        );
    }
}
```

### 3.6 Custom AuthenticationProvider

```java
@Component
public class CustomAuthenticationProvider implements AuthenticationProvider {

    private final UserDetailsService userDetailsService;
    private final PasswordEncoder passwordEncoder;
    private final TwoFactorService twoFactorService;

    @Override
    public Authentication authenticate(Authentication authentication) 
            throws AuthenticationException {
        
        String username = authentication.getName();
        String password = authentication.getCredentials().toString();
        
        // Get 2FA code from details (custom implementation)
        String twoFactorCode = ((CustomAuthDetails) authentication.getDetails())
            .getTwoFactorCode();

        // Load user
        UserDetails user = userDetailsService.loadUserByUsername(username);

        // Verify password
        if (!passwordEncoder.matches(password, user.getPassword())) {
            throw new BadCredentialsException("Invalid password");
        }

        // Verify 2FA
        if (!twoFactorService.verifyCode(username, twoFactorCode)) {
            throw new BadCredentialsException("Invalid 2FA code");
        }

        // Return authenticated token
        return new UsernamePasswordAuthenticationToken(
            user, 
            null,  // Clear credentials
            user.getAuthorities()
        );
    }

    @Override
    public boolean supports(Class<?> authentication) {
        return UsernamePasswordAuthenticationToken.class.isAssignableFrom(authentication);
    }
}
```

### 3.7 Form Login Configuration

```java
@Configuration
@EnableWebSecurity
public class FormLoginConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/login", "/register", "/css/**", "/js/**").permitAll()
                .anyRequest().authenticated()
            )
            .formLogin(form -> form
                .loginPage("/login")                    // Custom login page
                .loginProcessingUrl("/perform_login")   // Form action URL
                .usernameParameter("email")             // Custom username field
                .passwordParameter("pass")              // Custom password field
                .defaultSuccessUrl("/dashboard", true)  // Redirect after success
                .failureUrl("/login?error=true")        // Redirect after failure
                .successHandler(customSuccessHandler()) // Custom success handling
                .failureHandler(customFailureHandler()) // Custom failure handling
                .permitAll()
            )
            .logout(logout -> logout
                .logoutUrl("/perform_logout")
                .logoutSuccessUrl("/login?logout=true")
                .deleteCookies("JSESSIONID")
                .invalidateHttpSession(true)
                .clearAuthentication(true)
            )
            .build();
    }

    @Bean
    public AuthenticationSuccessHandler customSuccessHandler() {
        return (request, response, authentication) -> {
            // Log successful login
            log.info("User {} logged in from {}", 
                authentication.getName(), 
                request.getRemoteAddr());
            
            // Redirect based on role
            if (authentication.getAuthorities().stream()
                    .anyMatch(a -> a.getAuthority().equals("ROLE_ADMIN"))) {
                response.sendRedirect("/admin/dashboard");
            } else {
                response.sendRedirect("/dashboard");
            }
        };
    }

    @Bean
    public AuthenticationFailureHandler customFailureHandler() {
        return (request, response, exception) -> {
            String errorMessage = "Invalid credentials";
            
            if (exception instanceof LockedException) {
                errorMessage = "Account is locked";
            } else if (exception instanceof DisabledException) {
                errorMessage = "Account is disabled";
            } else if (exception instanceof BadCredentialsException) {
                errorMessage = "Invalid username or password";
            }
            
            request.getSession().setAttribute("errorMessage", errorMessage);
            response.sendRedirect("/login?error=true");
        };
    }
}
```

### 3.8 HTTP Basic Authentication

```java
@Bean
public SecurityFilterChain apiFilterChain(HttpSecurity http) throws Exception {
    return http
        .securityMatcher("/api/**")
        .csrf(csrf -> csrf.disable())
        .sessionManagement(session -> 
            session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
        .authorizeHttpRequests(auth -> auth.anyRequest().authenticated())
        .httpBasic(basic -> basic
            .realmName("My API")
            .authenticationEntryPoint(customEntryPoint())
        )
        .build();
}

@Bean
public AuthenticationEntryPoint customEntryPoint() {
    return (request, response, authException) -> {
        response.setContentType("application/json");
        response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
        response.getWriter().write("""
            {"error": "Unauthorized", "message": "Authentication required"}
            """);
    };
}
```

---

## 4. Authorization & Access Control

### 4.1 Authorization Concepts

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    AUTHORIZATION MODELS                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  RBAC (Role-Based Access Control)                                       │
│  ─────────────────────────────────                                      │
│  User → Role → Permissions                                              │
│                                                                         │
│  User "Alice" → ROLE_ADMIN → [read, write, delete, admin]               │
│  User "Bob"   → ROLE_USER  → [read]                                     │
│                                                                         │
│  Simple, widely used, can become rigid                                  │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  ABAC (Attribute-Based Access Control)                                  │
│  ─────────────────────────────────────                                  │
│  Access based on attributes of:                                         │
│  • Subject (user): department, clearance, role                          │
│  • Resource: classification, owner, type                                │
│  • Action: read, write, approve                                         │
│  • Environment: time, location, device                                  │
│                                                                         │
│  Example: "Managers in Finance can approve expenses under $10K          │
│            during business hours from company network"                  │
│                                                                         │
│  Flexible, complex to implement                                         │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  ACL (Access Control List)                                              │
│  ─────────────────────────                                              │
│  Per-object permissions                                                 │
│                                                                         │
│  Document 123: Alice=read,write; Bob=read; Carol=admin                  │
│                                                                         │
│  Fine-grained, can be hard to manage at scale                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 4.2 URL-Based Authorization

```java
@Configuration
@EnableWebSecurity
public class AuthorizationConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                // Public endpoints
                .requestMatchers("/", "/public/**", "/health").permitAll()
                
                // Static resources
                .requestMatchers("/css/**", "/js/**", "/images/**").permitAll()
                
                // API versioning
                .requestMatchers("/api/v1/**").authenticated()
                
                // Role-based
                .requestMatchers("/admin/**").hasRole("ADMIN")
                .requestMatchers("/manager/**").hasAnyRole("ADMIN", "MANAGER")
                
                // Authority-based (without ROLE_ prefix)
                .requestMatchers("/reports/**").hasAuthority("VIEW_REPORTS")
                .requestMatchers(HttpMethod.DELETE, "/api/**").hasAuthority("DELETE")
                
                // HTTP method specific
                .requestMatchers(HttpMethod.GET, "/api/users/**").authenticated()
                .requestMatchers(HttpMethod.POST, "/api/users/**").hasRole("ADMIN")
                
                // Pattern matching
                .requestMatchers("/api/users/{id}").access(
                    (authentication, context) -> {
                        Long userId = Long.parseLong(context.getVariables().get("id"));
                        return new AuthorizationDecision(
                            userId.equals(getCurrentUserId(authentication.get()))
                        );
                    }
                )
                
                // Default: deny all
                .anyRequest().denyAll()
            )
            .build();
    }
}
```

### 4.3 GrantedAuthority and Roles

```java
// GrantedAuthority is just a String wrapper
public interface GrantedAuthority {
    String getAuthority();
}

// ROLE_ prefix convention
// hasRole("ADMIN") checks for ROLE_ADMIN
// hasAuthority("ROLE_ADMIN") also checks for ROLE_ADMIN
// hasAuthority("DELETE") checks for DELETE (no prefix)

// Creating authorities
List<GrantedAuthority> authorities = List.of(
    new SimpleGrantedAuthority("ROLE_USER"),
    new SimpleGrantedAuthority("ROLE_ADMIN"),
    new SimpleGrantedAuthority("READ"),
    new SimpleGrantedAuthority("WRITE"),
    new SimpleGrantedAuthority("DELETE")
);
```

### 4.4 Role Hierarchy

```java
@Configuration
public class RoleHierarchyConfig {

    @Bean
    public RoleHierarchy roleHierarchy() {
        RoleHierarchyImpl hierarchy = new RoleHierarchyImpl();
        hierarchy.setHierarchy("""
            ROLE_ADMIN > ROLE_MANAGER
            ROLE_MANAGER > ROLE_USER
            ROLE_USER > ROLE_GUEST
            """);
        return hierarchy;
    }

    // Use in expressions with custom MethodSecurityExpressionHandler
    @Bean
    public DefaultMethodSecurityExpressionHandler methodSecurityExpressionHandler(
            RoleHierarchy roleHierarchy) {
        DefaultMethodSecurityExpressionHandler handler = 
            new DefaultMethodSecurityExpressionHandler();
        handler.setRoleHierarchy(roleHierarchy);
        return handler;
    }
}
```

**Role Hierarchy Visualization:**
```
ROLE_ADMIN
    │
    ├── (inherits all from)
    │
ROLE_MANAGER
    │
    ├── (inherits all from)
    │
ROLE_USER
    │
    ├── (inherits all from)
    │
ROLE_GUEST
```

### 4.5 SpEL Security Expressions

| Expression | Description |
|---|---|
| `hasRole('ADMIN')` | Has ROLE_ADMIN |
| `hasAnyRole('ADMIN', 'MANAGER')` | Has any of these roles |
| `hasAuthority('DELETE')` | Has DELETE authority |
| `hasAnyAuthority('READ', 'WRITE')` | Has any authority |
| `permitAll` | Always allows |
| `denyAll` | Always denies |
| `isAnonymous()` | Is anonymous user |
| `isAuthenticated()` | Is authenticated |
| `isRememberMe()` | Is remember-me auth |
| `isFullyAuthenticated()` | Is NOT anonymous or remember-me |
| `principal` | The Principal object |
| `authentication` | The Authentication object |
| `#paramName` | Method parameter reference |

**Complex Expressions:**

```java
@PreAuthorize("""
    hasRole('ADMIN') or 
    (hasRole('MANAGER') and #department == authentication.principal.department)
    """)
public Report generateReport(String department) { }

@PreAuthorize("@customSecurityService.canAccess(#id, authentication)")
public Resource getResource(Long id) { }
```

---

## 5. Password Storage & Encoding

### 5.1 Password Security Theory

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PASSWORD STORAGE EVOLUTION                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ❌ PLAIN TEXT                                                          │
│     password = "secret123"                                              │
│     Problem: Anyone with DB access sees all passwords                   │
│                                                                         │
│  ❌ ENCRYPTED                                                           │
│     password = AES.encrypt("secret123", key)                            │
│     Problem: Decryptable if key is compromised                          │
│                                                                         │
│  ❌ SIMPLE HASH (MD5, SHA-1)                                            │
│     password = MD5("secret123") = "f1d3ff8443297732862df21dc4e57262"   │
│     Problem: Rainbow tables, fast brute force                           │
│                                                                         │
│  ❌ HASH + SALT                                                         │
│     password = SHA256(salt + "secret123")                               │
│     Problem: Fast hashes enable brute force at scale                    │
│                                                                         │
│  ✅ ADAPTIVE HASH (bcrypt, scrypt, Argon2)                              │
│     password = bcrypt("secret123", cost=12)                             │
│     Benefits:                                                           │
│     • Intentionally slow (configurable work factor)                     │
│     • Built-in salt                                                     │
│     • Resistant to GPU/ASIC attacks (memory-hard)                       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.1.1 Cryptography Fundamentals — Theory Deep Dive

**Understanding the building blocks of secure password storage and authentication.**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CRYPTOGRAPHIC PRIMITIVES                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  HASHING (One-Way Function)                                             │
│  ──────────────────────────                                             │
│                                                                         │
│  Input (any size) ─────► [ HASH FUNCTION ] ─────► Output (fixed size)   │
│  "password123"                SHA-256              "ef92b778..." (256 bits)│
│                                                                         │
│  Properties:                                                            │
│  • Deterministic: Same input → same output                              │
│  • One-way: Cannot reverse hash to get input                            │
│  • Collision-resistant: Hard to find two inputs with same hash          │
│  • Avalanche effect: Small input change → completely different hash     │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  SYMMETRIC ENCRYPTION (Same Key)                                        │
│  ───────────────────────────────                                        │
│                                                                         │
│  Plaintext ──► [ ENCRYPT + KEY ] ──► Ciphertext                         │
│  Ciphertext ─► [ DECRYPT + KEY ] ──► Plaintext                          │
│                                                                         │
│  Examples: AES-256, ChaCha20                                            │
│  Use case: Encrypting data at rest, session data                        │
│  Problem: Key distribution — how to share key securely?                 │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  ASYMMETRIC ENCRYPTION (Key Pair)                                       │
│  ────────────────────────────────                                       │
│                                                                         │
│  ┌──────────────┐    ┌──────────────┐                                   │
│  │ PUBLIC KEY   │    │ PRIVATE KEY  │                                   │
│  │ (share with  │    │ (keep secret)│                                   │
│  │  everyone)   │    │              │                                   │
│  └──────┬───────┘    └──────┬───────┘                                   │
│         │                    │                                          │
│         ▼                    ▼                                          │
│    Encrypt data         Decrypt data                                    │
│    Verify signature     Create signature                                │
│                                                                         │
│  Examples: RSA, ECDSA (Elliptic Curve)                                  │
│  Use case: JWT signing (RS256), TLS handshake, OAuth2                   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Why Fast Hashes Are BAD for Passwords:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    HASH SPEED COMPARISON                                │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Algorithm        Speed (hashes/sec)    Time to crack 8-char password   │
│  ─────────        ──────────────────    ─────────────────────────────   │
│  MD5              10 billion            Seconds                         │
│  SHA-256          1 billion             Minutes                         │
│  bcrypt (cost=10) 100                   Centuries                       │
│  Argon2           10                    Millennia                       │
│                                                                         │
│  Modern GPUs can compute BILLIONS of SHA-256 hashes per second!         │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  RAINBOW TABLES                                                         │
│  ──────────────                                                         │
│  Pre-computed hash → password mappings                                  │
│                                                                         │
│  5d41402abc4b2a76 → "hello"                                             │
│  098f6bcd4621d373 → "test"                                              │
│  ... millions of entries ...                                            │
│                                                                         │
│  Defense: SALT (random data added before hashing)                       │
│  hash(salt + password) → even common passwords produce unique hashes    │
│                                                                         │
│  bcrypt, Argon2, scrypt include salt automatically!                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Key Stretching — Making Brute Force Impractical:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    KEY STRETCHING                                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Idea: Intentionally slow down hashing                                  │
│                                                                         │
│  Single SHA-256:                                                        │
│  password ──► SHA-256 ──► hash                                          │
│  Time: 0.000001 seconds                                                 │
│                                                                         │
│  PBKDF2 (10,000 iterations):                                            │
│  password ──► SHA-256 ──► SHA-256 ──► ... (10,000x) ──► hash            │
│  Time: 0.01 seconds                                                     │
│                                                                         │
│  bcrypt (cost=12):                                                      │
│  password ──► Blowfish setup ──► EksBlowfish (4096 iterations) ──► hash │
│  Time: 0.3 seconds                                                      │
│                                                                         │
│  Argon2 (memory-hard):                                                  │
│  password ──► Fill memory ──► Multiple passes ──► hash                  │
│  Requires: Time + Memory (defeats GPU parallelism)                      │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  Work Factor Tuning:                                                    │
│  • Too low: Fast brute force                                            │
│  • Too high: Slow login, DoS risk                                       │
│  • Sweet spot: 100-500ms per hash (adjust as hardware improves)         │
│                                                                         │
│  // Benchmark on YOUR hardware                                          │
│  BCryptPasswordEncoder encoder = new BCryptPasswordEncoder(12);         │
│  long start = System.currentTimeMillis();                               │
│  encoder.encode("test");                                                │
│  System.out.println("Time: " + (System.currentTimeMillis() - start));   │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Algorithm Comparison:**

| Algorithm | Type | Memory-Hard | GPU Resistant | Standard |
|---|---|---|---|---|
| bcrypt | Adaptive | ❌ | Moderate | De facto |
| scrypt | Adaptive | ✅ | Good | RFC 7914 |
| Argon2id | Adaptive | ✅ | Excellent | PHC Winner |
| PBKDF2 | Iterative | ❌ | Poor | NIST SP 800-132 |

### 5.2 PasswordEncoder Interface

```java
public interface PasswordEncoder {
    String encode(CharSequence rawPassword);
    boolean matches(CharSequence rawPassword, String encodedPassword);
    default boolean upgradeEncoding(String encodedPassword) {
        return false;  // Override to trigger re-encoding
    }
}
```

### 5.3 Available Password Encoders

| Encoder | Algorithm | Recommended |
|---|---|---|
| `BCryptPasswordEncoder` | bcrypt | ✅ Yes (default choice) |
| `Argon2PasswordEncoder` | Argon2 | ✅ Yes (memory-hard) |
| `SCryptPasswordEncoder` | scrypt | ✅ Yes (memory-hard) |
| `Pbkdf2PasswordEncoder` | PBKDF2 | ✅ Yes (NIST approved) |
| `NoOpPasswordEncoder` | Plain text | ❌ Never in production |
| `StandardPasswordEncoder` | SHA-256 | ❌ Deprecated, weak |
| `Md5PasswordEncoder` | MD5 | ❌ Never use |

### 5.4 BCrypt Configuration

```java
@Configuration
public class PasswordConfig {

    @Bean
    public PasswordEncoder passwordEncoder() {
        // Strength 4-31 (default 10)
        // Higher = slower = more secure but impacts login performance
        // 10 ≈ 100ms, 12 ≈ 300ms, 14 ≈ 1s
        return new BCryptPasswordEncoder(12);
    }
}
```

**BCrypt Hash Format:**
```
$2a$12$N9qo8uLOickgx2ZMRZoMyeIjZRGT.XuR.FmS5MH3L7T1Xq1B2Zc0G
│  │  │                                                        │
│  │  │                                                        └─ Hash (31 chars)
│  │  └─ Salt (22 chars, random, different each time)
│  └─ Cost factor (10^12 iterations)
└─ Algorithm version ($2a$, $2b$, $2y$)
```

### 5.5 DelegatingPasswordEncoder (Recommended)

Supports **password migration** between algorithms:

```java
@Bean
public PasswordEncoder passwordEncoder() {
    // Supports multiple encodings, default is bcrypt
    return PasswordEncoderFactories.createDelegatingPasswordEncoder();
}

// Storage format: {encoderId}encodedPassword
// {bcrypt}$2a$10$dXJ3SW6G7P50lGmMkkmwe.20cQQubK3.HZWzG3YB1tlRy.fqvM/BG
// {noop}plaintext     (for testing only!)
// {pbkdf2}5d923b44a6d129f3...
// {scrypt}$e0801$8bWJaSu2IKSn9Z9kM+TPXfOc/9bdYSrN1oD9qfVThWEwdR...
```

### 5.6 Password Migration Strategy

```java
@Service
public class UserService {

    private final PasswordEncoder passwordEncoder;
    private final UserRepository userRepository;

    public UserDetails loadUserByUsername(String username) {
        User user = userRepository.findByUsername(username).orElseThrow();
        
        // Check if password needs re-encoding
        if (passwordEncoder.upgradeEncoding(user.getPassword())) {
            // Will be re-encoded on next login
            user.setPasswordNeedsUpgrade(true);
        }
        
        return new CustomUserDetails(user);
    }

    @Transactional
    public void onSuccessfulLogin(String username, String rawPassword) {
        User user = userRepository.findByUsername(username).orElseThrow();
        
        if (user.isPasswordNeedsUpgrade()) {
            // Re-encode with current algorithm
            user.setPassword(passwordEncoder.encode(rawPassword));
            user.setPasswordNeedsUpgrade(false);
            userRepository.save(user);
        }
    }
}
```

---

## 6. Session Management

### 6.1 Session Security Concepts

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    SESSION SECURITY                                     │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Session Fixation Attack:                                               │
│  ─────────────────────────                                              │
│  1. Attacker obtains a session ID (e.g., visits site)                   │
│  2. Attacker tricks victim into using same session ID                   │
│  3. Victim authenticates, attacker now has authenticated session!       │
│                                                                         │
│  Prevention: Change session ID after authentication                     │
│              Spring Security does this by default!                      │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  Concurrent Session Attack:                                             │
│  ──────────────────────────                                             │
│  Multiple logins from different locations                               │
│  Could indicate credential theft                                        │
│                                                                         │
│  Prevention: Limit concurrent sessions per user                         │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 6.2 Session Management Configuration

```java
@Configuration
@EnableWebSecurity
public class SessionConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .sessionManagement(session -> session
                // Session fixation protection (default: changeSessionId)
                .sessionFixation(fixation -> fixation.changeSessionId())
                // .sessionFixation(fixation -> fixation.migrateSession())  // Copies attributes
                // .sessionFixation(fixation -> fixation.newSession())      // New empty session
                // .sessionFixation(fixation -> fixation.none())            // Disable (NOT recommended!)
                
                // Session creation policy
                .sessionCreationPolicy(SessionCreationPolicy.IF_REQUIRED)
                // IF_REQUIRED (default): Creates session when needed
                // ALWAYS: Always creates session
                // NEVER: Never creates, but uses existing
                // STATELESS: Never creates or uses (for JWT APIs)
                
                // Invalid session handling
                .invalidSessionUrl("/login?invalid-session")
                
                // Maximum sessions per user
                .maximumSessions(1)
                    .maxSessionsPreventsLogin(false)  // false = kick previous session
                                                       // true = prevent new login
                    .expiredUrl("/login?session-expired")
            )
            .build();
    }

    // Required for concurrent session control
    @Bean
    public HttpSessionEventPublisher httpSessionEventPublisher() {
        return new HttpSessionEventPublisher();
    }
}
```

### 6.3 Session Storage Options

| Storage | Use Case | Spring Config |
|---|---|---|
| **In-memory** | Single server, development | Default |
| **JDBC** | Multiple servers, simple | `spring-session-jdbc` |
| **Redis** | High performance, scalable | `spring-session-data-redis` |
| **Hazelcast** | Distributed, in-memory grid | `spring-session-hazelcast` |
| **MongoDB** | Document-based | `spring-session-data-mongodb` |

**Redis Session Configuration:**

```yaml
spring:
  session:
    store-type: redis
    timeout: 30m
  data:
    redis:
      host: localhost
      port: 6379
```

```java
@Configuration
@EnableRedisHttpSession(maxInactiveIntervalInSeconds = 1800)
public class RedisSessionConfig {
    
    @Bean
    public RedisConnectionFactory connectionFactory() {
        return new LettuceConnectionFactory();
    }
}
```

### 6.4 Remember-Me Authentication

```java
@Bean
public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
    return http
        .rememberMe(remember -> remember
            .key("uniqueAndSecretKey")  // For token signing
            .tokenValiditySeconds(86400 * 30)  // 30 days
            .userDetailsService(userDetailsService)
            .rememberMeParameter("remember-me")  // Form checkbox name
            
            // Persistent tokens (more secure)
            .tokenRepository(persistentTokenRepository())
            
            // Cookie config
            .rememberMeCookieName("remember-me")
            .rememberMeCookieDomain(".example.com")
            .useSecureCookie(true)
        )
        .build();
}

@Bean
public PersistentTokenRepository persistentTokenRepository() {
    JdbcTokenRepositoryImpl repository = new JdbcTokenRepositoryImpl();
    repository.setDataSource(dataSource);
    // Create table: persistent_logins
    // repository.setCreateTableOnStartup(true);
    return repository;
}
```

**Remember-Me Token Table:**
```sql
CREATE TABLE persistent_logins (
    username VARCHAR(64) NOT NULL,
    series VARCHAR(64) PRIMARY KEY,
    token VARCHAR(64) NOT NULL,
    last_used TIMESTAMP NOT NULL
);
```

---

## 7. JWT Authentication

### 7.1 JWT Theory

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    JWT (JSON Web Token) STRUCTURE                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.                                   │
│  eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4iLCJpYXQiOjE1MTYyMzkwMjJ9.  │
│  SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c                            │
│  │                                                                      │
│  └─ Three parts separated by dots (.)                                   │
│                                                                         │
│  ┌─────────────────┐                                                    │
│  │     HEADER      │  Base64Url encoded JSON                            │
│  │  {              │  • alg: Signing algorithm (HS256, RS256)           │
│  │    "alg":"HS256"│  • typ: Token type (JWT)                           │
│  │    "typ":"JWT"  │                                                    │
│  │  }              │                                                    │
│  └─────────────────┘                                                    │
│          .                                                              │
│  ┌─────────────────┐                                                    │
│  │     PAYLOAD     │  Base64Url encoded JSON                            │
│  │  {              │  • Claims (statements about user)                  │
│  │    "sub":"123", │  • Registered: iss, sub, aud, exp, nbf, iat, jti   │
│  │    "name":"John"│  • Public: name, email, roles (your data)          │
│  │    "iat":15162..│  • Private: custom claims                          │
│  │    "exp":15163..│                                                    │
│  │  }              │                                                    │
│  └─────────────────┘                                                    │
│          .                                                              │
│  ┌─────────────────┐                                                    │
│  │    SIGNATURE    │  HMACSHA256(                                       │
│  │  SflKxwRJSMe... │    base64UrlEncode(header) + "." +                 │
│  │                 │    base64UrlEncode(payload),                       │
│  │                 │    secret                                          │
│  │                 │  )                                                 │
│  └─────────────────┘                                                    │
│                                                                         │
│  NOTE: Payload is NOT encrypted, only signed!                           │
│        Don't put sensitive data in JWTs!                                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 7.2 JWT vs Sessions Comparison

| Aspect | Session | JWT |
|---|---|---|
| **Storage** | Server-side | Client-side (token) |
| **Scalability** | Requires sticky sessions or shared store | Stateless, any server can verify |
| **Revocation** | Easy (delete from store) | Hard (requires blacklist) |
| **Size** | Small cookie (session ID) | Larger (contains claims) |
| **Security** | Session fixation risks | Token theft risks |
| **Use Case** | Traditional web apps | APIs, microservices, SPAs |

### 7.2.1 Stateless vs Stateful Authentication — Theory Deep Dive

**Understanding the fundamental architectural trade-offs.**

```
┌─────────────────────────────────────────────────────────────────────────┐
│             STATEFUL AUTHENTICATION (Sessions)                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────┐        ┌──────────────────┐        ┌──────────────────┐   │
│  │  Client  │───────►│     Server       │───────►│  Session Store   │   │
│  │          │        │                  │        │  (Memory/Redis)  │   │
│  │ Cookie:  │        │ SecurityContext  │        │                  │   │
│  │ SESSION= │◄───────│ loaded from      │◄───────│ {user, roles,    │   │
│  │ abc123   │        │ session on each  │        │  lastAccess...}  │   │
│  └──────────┘        │ request          │        └──────────────────┘   │
│                      └──────────────────┘                               │
│                                                                         │
│  Flow:                                                                  │
│  1. User logs in → Server creates session → Returns session ID cookie   │
│  2. Each request → Cookie sent → Server looks up session → Gets user    │
│  3. Logout → Server deletes session → Cookie invalidated                │
│                                                                         │
│  Pros:                                Cons:                             │
│  ✅ Easy revocation (delete session)  ❌ Server stores state            │
│  ✅ Small cookie size                  ❌ Horizontal scaling complex    │
│  ✅ Server controls session lifetime   ❌ Requires sticky sessions OR   │
│  ✅ Can update user mid-session            distributed session store    │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────┐
│             STATELESS AUTHENTICATION (Tokens/JWT)                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  ┌──────────┐        ┌──────────────────┐                               │
│  │  Client  │───────►│     Server       │     No session store needed!  │
│  │          │        │                  │                               │
│  │ Header:  │        │ Validates token  │     All user info is IN       │
│  │ Bearer   │◄───────│ signature        │     the token itself          │
│  │ eyJhb... │        │ Extracts claims  │                               │
│  └──────────┘        └──────────────────┘                               │
│                                                                         │
│  Flow:                                                                  │
│  1. User logs in → Server creates JWT (signed) → Returns token          │
│  2. Each request → Token in header → Server verifies signature,         │
│                    extracts claims → User authenticated                 │
│  3. Logout → Client deletes token (server stateless - can't invalidate) │
│                                                                         │
│  Pros:                                Cons:                             │
│  ✅ Truly stateless                    ❌ Hard to revoke (need blacklist)│
│  ✅ Easy horizontal scaling            ❌ Larger payload size            │
│  ✅ Works across domains/services      ❌ Token theft = full access      │
│  ✅ Perfect for microservices          ❌ Can't update mid-session       │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**The Token Revocation Problem:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    JWT REVOCATION STRATEGIES                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Problem: JWTs are self-contained. Server can't "delete" them.          │
│           A stolen token is valid until expiration!                     │
│                                                                         │
│  Solutions:                                                             │
│                                                                         │
│  1. SHORT EXPIRATION + REFRESH TOKENS                                   │
│     ────────────────────────────────                                    │
│     Access token: 15 minutes (short window if stolen)                   │
│     Refresh token: 7 days (stored server-side, revocable)               │
│                                                                         │
│     Trade-off: More refresh requests, but limited damage window         │
│                                                                         │
│  2. TOKEN BLACKLIST                                                     │
│     ───────────────                                                     │
│     Store revoked token IDs (jti) in Redis with TTL                     │
│     Check blacklist on every request                                    │
│                                                                         │
│     Trade-off: Adds state, but only for revoked tokens                  │
│                                                                         │
│  3. TOKEN VERSIONING                                                    │
│     ─────────────────                                                   │
│     Store "token_version" per user in DB                                │
│     Include version in JWT claims                                       │
│     Password change → increment version → all old tokens invalid        │
│                                                                         │
│     Trade-off: DB lookup per request, but enables mass revocation       │
│                                                                         │
│  4. REFERENCE TOKENS (Opaque)                                           │
│     ─────────────────────────                                           │
│     Return random string, store actual data server-side                 │
│     Basically sessions with "token" name                                │
│                                                                         │
│     Trade-off: Not truly stateless, but easy revocation                 │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Choosing the Right Approach:**

| Scenario | Recommendation |
|---|---|
| Monolith, server-rendered pages | Sessions (simple, easy revocation) |
| SPA with same-domain API | Sessions with SameSite cookies |
| Mobile app | JWT with refresh tokens |
| Microservices | JWT (stateless, cross-service) |
| High-security (banking) | Short JWT + refresh + token versioning |
| Third-party API access | OAuth2 with JWT |

### 7.3 JWT Service Implementation

```java
@Service
public class JwtService {

    @Value("${jwt.secret}")
    private String secretKey;

    @Value("${jwt.expiration}")
    private Duration accessTokenExpiration;

    @Value("${jwt.refresh-expiration}")
    private Duration refreshTokenExpiration;

    private SecretKey getSigningKey() {
        byte[] keyBytes = Decoders.BASE64.decode(secretKey);
        return Keys.hmacShaKeyFor(keyBytes);
    }

    public String generateAccessToken(UserDetails userDetails) {
        return generateToken(userDetails, accessTokenExpiration, Map.of(
            "type", "access",
            "roles", userDetails.getAuthorities().stream()
                .map(GrantedAuthority::getAuthority)
                .toList()
        ));
    }

    public String generateRefreshToken(UserDetails userDetails) {
        return generateToken(userDetails, refreshTokenExpiration, Map.of(
            "type", "refresh"
        ));
    }

    private String generateToken(UserDetails userDetails, Duration expiration, 
                                 Map<String, Object> extraClaims) {
        Instant now = Instant.now();
        return Jwts.builder()
            .setClaims(extraClaims)
            .setSubject(userDetails.getUsername())
            .setIssuedAt(Date.from(now))
            .setExpiration(Date.from(now.plus(expiration)))
            .setId(UUID.randomUUID().toString())  // jti - for revocation
            .signWith(getSigningKey(), SignatureAlgorithm.HS256)
            .compact();
    }

    public String extractUsername(String token) {
        return extractClaim(token, Claims::getSubject);
    }

    public boolean isTokenValid(String token, UserDetails userDetails) {
        final String username = extractUsername(token);
        return username.equals(userDetails.getUsername()) && !isTokenExpired(token);
    }

    private boolean isTokenExpired(String token) {
        return extractClaim(token, Claims::getExpiration).before(new Date());
    }

    private <T> T extractClaim(String token, Function<Claims, T> resolver) {
        Claims claims = Jwts.parserBuilder()
            .setSigningKey(getSigningKey())
            .build()
            .parseClaimsJws(token)
            .getBody();
        return resolver.apply(claims);
    }

    public Claims extractAllClaims(String token) {
        return Jwts.parserBuilder()
            .setSigningKey(getSigningKey())
            .build()
            .parseClaimsJws(token)
            .getBody();
    }
}
```

### 7.4 JWT Authentication Filter

```java
@Component
@RequiredArgsConstructor
public class JwtAuthenticationFilter extends OncePerRequestFilter {

    private final JwtService jwtService;
    private final UserDetailsService userDetailsService;

    @Override
    protected void doFilterInternal(
            HttpServletRequest request,
            HttpServletResponse response,
            FilterChain filterChain) throws ServletException, IOException {

        // 1. Extract Authorization header
        final String authHeader = request.getHeader("Authorization");

        if (authHeader == null || !authHeader.startsWith("Bearer ")) {
            filterChain.doFilter(request, response);
            return;
        }

        // 2. Extract token
        final String jwt = authHeader.substring(7);

        try {
            // 3. Extract username from token
            final String username = jwtService.extractUsername(jwt);

            // 4. If not already authenticated
            if (username != null && 
                SecurityContextHolder.getContext().getAuthentication() == null) {

                // 5. Load user details
                UserDetails userDetails = userDetailsService.loadUserByUsername(username);

                // 6. Validate token
                if (jwtService.isTokenValid(jwt, userDetails)) {
                    
                    // 7. Create authentication object
                    UsernamePasswordAuthenticationToken authToken =
                        new UsernamePasswordAuthenticationToken(
                            userDetails,
                            null,
                            userDetails.getAuthorities()
                        );

                    // 8. Set details
                    authToken.setDetails(
                        new WebAuthenticationDetailsSource().buildDetails(request)
                    );

                    // 9. Update SecurityContext
                    SecurityContextHolder.getContext().setAuthentication(authToken);
                }
            }
        } catch (ExpiredJwtException e) {
            // Token expired — let it through, will be handled by entry point
            request.setAttribute("jwt-expired", true);
        } catch (JwtException e) {
            // Invalid token
            request.setAttribute("jwt-invalid", true);
        }

        filterChain.doFilter(request, response);
    }

    @Override
    protected boolean shouldNotFilter(HttpServletRequest request) {
        return request.getServletPath().startsWith("/api/auth/");
    }
}
```

### 7.5 JWT Security Configuration

```java
@Configuration
@EnableWebSecurity
@RequiredArgsConstructor
public class JwtSecurityConfig {

    private final JwtAuthenticationFilter jwtAuthFilter;
    private final AuthenticationProvider authenticationProvider;

    @Bean
    public SecurityFilterChain securityFilterChain(HttpSecurity http) throws Exception {
        return http
            .csrf(csrf -> csrf.disable())
            .sessionManagement(session -> 
                session.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/api/auth/**").permitAll()
                .requestMatchers("/api/public/**").permitAll()
                .anyRequest().authenticated()
            )
            .authenticationProvider(authenticationProvider)
            .addFilterBefore(jwtAuthFilter, UsernamePasswordAuthenticationFilter.class)
            .exceptionHandling(ex -> ex
                .authenticationEntryPoint(jwtAuthenticationEntryPoint())
                .accessDeniedHandler(jwtAccessDeniedHandler())
            )
            .build();
    }

    @Bean
    public AuthenticationEntryPoint jwtAuthenticationEntryPoint() {
        return (request, response, authException) -> {
            response.setContentType(MediaType.APPLICATION_JSON_VALUE);
            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
            
            String message = "Unauthorized";
            if (Boolean.TRUE.equals(request.getAttribute("jwt-expired"))) {
                message = "Token expired";
            } else if (Boolean.TRUE.equals(request.getAttribute("jwt-invalid"))) {
                message = "Invalid token";
            }
            
            response.getWriter().write("""
                {"error": "Unauthorized", "message": "%s"}
                """.formatted(message));
        };
    }

    @Bean
    public AccessDeniedHandler jwtAccessDeniedHandler() {
        return (request, response, accessDeniedException) -> {
            response.setContentType(MediaType.APPLICATION_JSON_VALUE);
            response.setStatus(HttpServletResponse.SC_FORBIDDEN);
            response.getWriter().write("""
                {"error": "Forbidden", "message": "Access denied"}
                """);
        };
    }
}
```

### 7.6 Auth Controller (Login/Refresh)

```java
@RestController
@RequestMapping("/api/auth")
@RequiredArgsConstructor
public class AuthController {

    private final AuthenticationManager authenticationManager;
    private final JwtService jwtService;
    private final UserDetailsService userDetailsService;
    private final RefreshTokenService refreshTokenService;

    @PostMapping("/login")
    public ResponseEntity<AuthResponse> login(@Valid @RequestBody LoginRequest request) {
        // Authenticate
        Authentication authentication = authenticationManager.authenticate(
            new UsernamePasswordAuthenticationToken(
                request.email(),
                request.password()
            )
        );

        UserDetails user = (UserDetails) authentication.getPrincipal();
        
        // Generate tokens
        String accessToken = jwtService.generateAccessToken(user);
        String refreshToken = jwtService.generateRefreshToken(user);
        
        // Store refresh token (for revocation)
        refreshTokenService.saveRefreshToken(user.getUsername(), refreshToken);

        return ResponseEntity.ok(new AuthResponse(accessToken, refreshToken));
    }

    @PostMapping("/refresh")
    public ResponseEntity<AuthResponse> refresh(@Valid @RequestBody RefreshRequest request) {
        String refreshToken = request.refreshToken();
        
        // Validate refresh token
        if (!refreshTokenService.isValid(refreshToken)) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
        }
        
        String username = jwtService.extractUsername(refreshToken);
        UserDetails user = userDetailsService.loadUserByUsername(username);
        
        // Generate new access token
        String newAccessToken = jwtService.generateAccessToken(user);
        
        // Optionally rotate refresh token
        String newRefreshToken = jwtService.generateRefreshToken(user);
        refreshTokenService.rotateRefreshToken(refreshToken, newRefreshToken);

        return ResponseEntity.ok(new AuthResponse(newAccessToken, newRefreshToken));
    }

    @PostMapping("/logout")
    public ResponseEntity<Void> logout(@RequestHeader("Authorization") String authHeader) {
        if (authHeader != null && authHeader.startsWith("Bearer ")) {
            String token = authHeader.substring(7);
            refreshTokenService.revokeByUsername(jwtService.extractUsername(token));
        }
        return ResponseEntity.ok().build();
    }
}

public record LoginRequest(
    @NotBlank @Email String email,
    @NotBlank String password
) {}

public record RefreshRequest(@NotBlank String refreshToken) {}

public record AuthResponse(String accessToken, String refreshToken) {}
```

### 7.7 Token Revocation Strategy

```java
@Service
@RequiredArgsConstructor
public class RefreshTokenService {

    private final RefreshTokenRepository repository;  // Redis or DB

    public void saveRefreshToken(String username, String token) {
        Claims claims = jwtService.extractAllClaims(token);
        RefreshToken entity = new RefreshToken(
            claims.getId(),  // jti
            username,
            token,
            claims.getExpiration().toInstant()
        );
        repository.save(entity);
    }

    public boolean isValid(String token) {
        try {
            Claims claims = jwtService.extractAllClaims(token);
            return repository.existsByJtiAndRevokedFalse(claims.getId());
        } catch (JwtException e) {
            return false;
        }
    }

    public void revoke(String token) {
        Claims claims = jwtService.extractAllClaims(token);
        repository.updateRevokedByJti(claims.getId(), true);
    }

    public void revokeByUsername(String username) {
        repository.updateRevokedByUsername(username, true);
    }

    public void rotateRefreshToken(String oldToken, String newToken) {
        revoke(oldToken);
        saveRefreshToken(jwtService.extractUsername(newToken), newToken);
    }
}
```

---

## 8. OAuth2 & OpenID Connect

### 8.1 OAuth2 Concepts

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    OAUTH2 ROLES                                         │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  RESOURCE OWNER:  The user who owns the data                            │
│                   (e.g., you and your Google account)                   │
│                                                                         │
│  CLIENT:          Application requesting access                         │
│                   (e.g., Spotify wants your Google profile)             │
│                                                                         │
│  AUTHORIZATION    Server that authenticates resource owner              │
│  SERVER:          and issues tokens (e.g., Google's auth server)        │
│                                                                         │
│  RESOURCE         Server hosting protected resources                    │
│  SERVER:          (e.g., Google's API server)                           │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.2 OAuth2 Grant Types

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    OAUTH2 FLOWS (GRANT TYPES)                           │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  AUTHORIZATION CODE (Most Common, Most Secure)                          │
│  ─────────────────────────────────────────────                          │
│  For: Server-side web apps                                              │
│                                                                         │
│  User ──(1)──> Client ──(2)──> Auth Server ──(3)──> User Login          │
│                                                     │                   │
│  User <──(4)── Auth Server (redirect with code)    │                   │
│       │                                             │                   │
│       └──(5)──> Client ──(6)──> Auth Server        ▼                   │
│                  (exchange code for token)    (consent)                 │
│                         │                                               │
│  Client <──(7)── Auth Server (access_token, refresh_token)              │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  AUTHORIZATION CODE + PKCE (SPAs, Mobile)                               │
│  ─────────────────────────────────────────                              │
│  Same as above but with code_verifier/code_challenge                    │
│  Prevents authorization code interception attacks                       │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  CLIENT CREDENTIALS (Machine-to-Machine)                                │
│  ─────────────────────────────────────────                              │
│  For: Service-to-service authentication                                 │
│                                                                         │
│  Client ──(client_id, client_secret)──> Auth Server                     │
│  Client <──(access_token)────────────── Auth Server                     │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  ❌ IMPLICIT (Deprecated)                                               │
│  ❌ PASSWORD (Deprecated) — Only for migration scenarios                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 8.3 OAuth2 vs OpenID Connect

| OAuth2 | OpenID Connect (OIDC) |
|---|---|
| Authorization protocol | Authentication layer ON TOP of OAuth2 |
| "Can app access my data?" | "Who is the user?" |
| Access tokens | ID tokens (JWT with user info) |
| Scopes: custom | Scopes: `openid`, `profile`, `email` |
| No standard user info | Standardized claims (sub, name, email) |

### 8.3.1 PKCE — Proof Key for Code Exchange (Theory)

**Why PKCE exists and how it works.**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    THE PROBLEM PKCE SOLVES                              │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Authorization Code Flow WITHOUT PKCE (vulnerable):                     │
│                                                                         │
│  1. SPA redirects to auth server                                        │
│  2. User authenticates                                                  │
│  3. Auth server redirects back with CODE in URL                         │
│       https://spa.example.com/callback?code=abc123                      │
│                                     │                                   │
│                                     ▼                                   │
│  ┌─────────────────────────────────────────────────────────────────┐    │
│  │  🚨 ATTACKER CAN INTERCEPT THIS CODE!                          │    │
│  │  • Malicious browser extension sees URL                         │    │
│  │  • Compromised redirect URI                                     │    │
│  │  • Man-in-the-middle on mobile (deep links)                     │    │
│  └─────────────────────────────────────────────────────────────────┘    │
│                                                                         │
│  4. Attacker exchanges code for tokens BEFORE legitimate client!        │
│                                                                         │
│  Traditional fix: Client secret (but SPAs can't keep secrets!)          │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**PKCE Solution:**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PKCE FLOW                                            │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  STEP 1: Generate code_verifier (client)                                │
│  ────────────────────────────────────────                               │
│  code_verifier = random_string(43-128 chars)                            │
│  // Example: "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"              │
│                                                                         │
│  STEP 2: Create code_challenge (client)                                 │
│  ───────────────────────────────────────                                │
│  code_challenge = BASE64URL(SHA256(code_verifier))                      │
│  // Example: "E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM"              │
│                                                                         │
│  STEP 3: Authorization request (client → auth server)                   │
│  ────────────────────────────────────────────────────                   │
│  GET /authorize?                                                        │
│    response_type=code&                                                  │
│    client_id=spa-app&                                                   │
│    redirect_uri=https://spa.example.com/callback&                       │
│    code_challenge=E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM&          │
│    code_challenge_method=S256                                           │
│                                                                         │
│  (Auth server stores code_challenge with the authorization code)        │
│                                                                         │
│  STEP 4: User authenticates, gets redirected with code                  │
│  ─────────────────────────────────────────────────────                  │
│  (Same as before — code in URL, potentially interceptable)              │
│                                                                         │
│  STEP 5: Token exchange (client → auth server)                          │
│  ─────────────────────────────────────────────                          │
│  POST /token                                                            │
│    grant_type=authorization_code&                                       │
│    code=abc123&                                                         │
│    redirect_uri=https://spa.example.com/callback&                       │
│    code_verifier=dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk  ← SECRET!│
│                                                                         │
│  STEP 6: Auth server validates                                          │
│  ─────────────────────────────────                                      │
│  calculated_challenge = BASE64URL(SHA256(code_verifier))                │
│  if (calculated_challenge == stored_code_challenge) {                   │
│    // ✅ Return tokens                                                  │
│  } else {                                                               │
│    // ❌ Reject - attacker doesn't know code_verifier!                  │
│  }                                                                      │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Why It Works:**
- Attacker intercepts `code`, but doesn't know `code_verifier`
- `code_verifier` never sent over network except in final token exchange
- Even if attacker guesses `code_challenge`, can't reverse SHA256 to get `code_verifier`

### 8.3.2 Identity Federation — Theory

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    IDENTITY FEDERATION PROTOCOLS                        │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  SAML 2.0 (Security Assertion Markup Language)                          │
│  ─────────────────────────────────────────────                          │
│  • XML-based, enterprise-focused                                        │
│  • Browser redirects with XML assertions                                │
│  • Complex, heavyweight                                                 │
│  • Still common in enterprise SSO (Okta, ADFS)                          │
│                                                                         │
│  OAuth 2.0 (Open Authorization)                                         │
│  ─────────────────────────────                                          │
│  • JSON-based, authorization-focused                                    │
│  • "Can this app access my data?"                                       │
│  • Access tokens, refresh tokens                                        │
│  • NOT for authentication (no user identity standard)                   │
│                                                                         │
│  OpenID Connect (OIDC)                                                  │
│  ────────────────────                                                   │
│  • Authentication layer ON TOP of OAuth 2.0                             │
│  • "Who is this user?"                                                  │
│  • ID Token (JWT with user claims)                                      │
│  • Standardized: Google, Microsoft, Auth0, Keycloak                     │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  Protocol Comparison:                                                   │
│                                                                         │
│                  SAML          OAuth2         OIDC                      │
│  Year            2005          2012           2014                      │
│  Format          XML           JSON           JSON                      │
│  Token           Assertion     Access Token   ID Token (JWT)            │
│  Purpose         AuthN+AuthZ   AuthZ only     AuthN (+ OAuth2 AuthZ)    │
│  Complexity      High          Medium         Medium                    │
│  Mobile-friendly No            Yes            Yes                       │
│  Use Today       Enterprise    APIs           Modern apps               │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**OIDC Token Types:**

| Token | Purpose | Format | Contains |
|---|---|---|---|
| **ID Token** | Prove user identity | JWT | sub, name, email, iat, exp |
| **Access Token** | Access APIs | JWT or opaque | scopes, permissions |
| **Refresh Token** | Get new tokens | Opaque | Only meaningful to auth server |

### 8.4 OAuth2 Client Configuration

```yaml
spring:
  security:
    oauth2:
      client:
        registration:
          google:
            client-id: ${GOOGLE_CLIENT_ID}
            client-secret: ${GOOGLE_CLIENT_SECRET}
            scope: openid, profile, email
            
          github:
            client-id: ${GITHUB_CLIENT_ID}
            client-secret: ${GITHUB_CLIENT_SECRET}
            scope: read:user, user:email
            
          custom-provider:
            client-id: ${CUSTOM_CLIENT_ID}
            client-secret: ${CUSTOM_CLIENT_SECRET}
            authorization-grant-type: authorization_code
            redirect-uri: "{baseUrl}/login/oauth2/code/{registrationId}"
            scope: openid, profile
            
        provider:
          custom-provider:
            authorization-uri: https://auth.example.com/oauth2/authorize
            token-uri: https://auth.example.com/oauth2/token
            user-info-uri: https://auth.example.com/userinfo
            jwk-set-uri: https://auth.example.com/.well-known/jwks.json
            user-name-attribute: sub
```

### 8.5 OAuth2 Login Security Configuration

```java
@Configuration
@EnableWebSecurity
public class OAuth2LoginConfig {

    @Bean
    public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
        return http
            .authorizeHttpRequests(auth -> auth
                .requestMatchers("/", "/login/**", "/error").permitAll()
                .anyRequest().authenticated()
            )
            .oauth2Login(oauth2 -> oauth2
                .loginPage("/login")
                .defaultSuccessUrl("/dashboard")
                .failureUrl("/login?error=true")
                .userInfoEndpoint(userInfo -> userInfo
                    .userService(customOAuth2UserService())
                    .oidcUserService(customOidcUserService())
                )
                .successHandler(oAuth2SuccessHandler())
            )
            .build();
    }

    @Bean
    public OAuth2UserService<OAuth2UserRequest, OAuth2User> customOAuth2UserService() {
        DefaultOAuth2UserService delegate = new DefaultOAuth2UserService();
        
        return userRequest -> {
            OAuth2User oAuth2User = delegate.loadUser(userRequest);
            
            String registrationId = userRequest.getClientRegistration().getRegistrationId();
            
            // Extract user info based on provider
            String email = extractEmail(oAuth2User, registrationId);
            String name = extractName(oAuth2User, registrationId);
            
            // Find or create user in your database
            User user = userService.findOrCreateOAuth2User(email, name, registrationId);
            
            // Return custom OAuth2User with your authorities
            return new CustomOAuth2User(oAuth2User, user);
        };
    }
}
```

### 8.6 OAuth2 Resource Server (JWT)

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
                .requestMatchers("/api/admin/**").hasAuthority("SCOPE_admin")
                .anyRequest().authenticated()
            )
            .oauth2ResourceServer(oauth2 -> oauth2
                .jwt(jwt -> jwt
                    .jwtAuthenticationConverter(jwtAuthenticationConverter())
                )
            )
            .build();
    }

    @Bean
    public JwtAuthenticationConverter jwtAuthenticationConverter() {
        JwtGrantedAuthoritiesConverter authoritiesConverter = 
            new JwtGrantedAuthoritiesConverter();
        // Map "roles" claim to authorities
        authoritiesConverter.setAuthoritiesClaimName("roles");
        authoritiesConverter.setAuthorityPrefix("ROLE_");

        JwtAuthenticationConverter converter = new JwtAuthenticationConverter();
        converter.setJwtGrantedAuthoritiesConverter(authoritiesConverter);
        return converter;
    }
}
```

---

## 9. Method-Level Security

### 9.1 Enabling Method Security

```java
@Configuration
@EnableMethodSecurity(
    prePostEnabled = true,   // @PreAuthorize, @PostAuthorize
    securedEnabled = true,   // @Secured
    jsr250Enabled = true     // @RolesAllowed
)
public class MethodSecurityConfig {
    // Configuration if needed
}
```

### 9.2 @PreAuthorize and @PostAuthorize

```java
@Service
public class DocumentService {

    // Checked BEFORE method execution
    @PreAuthorize("hasRole('ADMIN')")
    public void deleteAll() { }

    @PreAuthorize("hasRole('USER') and #userId == authentication.principal.id")
    public List<Document> getUserDocuments(Long userId) { }

    @PreAuthorize("@documentSecurity.canAccess(#id, authentication)")
    public Document getDocument(Long id) { }

    // Checked AFTER method execution
    @PostAuthorize("returnObject.owner == authentication.name")
    public Document getDocumentPostCheck(Long id) {
        return documentRepository.findById(id).orElseThrow();
    }

    // Complex expressions
    @PreAuthorize("""
        hasRole('ADMIN') or 
        (hasRole('MANAGER') and @teamService.isInSameTeam(#employeeId, authentication.principal.id))
        """)
    public Employee getEmployee(Long employeeId) { }
}

// Custom security component
@Component("documentSecurity")
public class DocumentSecurityService {
    
    public boolean canAccess(Long documentId, Authentication auth) {
        Document doc = documentRepository.findById(documentId).orElse(null);
        if (doc == null) return false;
        
        CustomUserDetails user = (CustomUserDetails) auth.getPrincipal();
        return doc.getOwnerId().equals(user.getUserId()) ||
               auth.getAuthorities().stream()
                   .anyMatch(a -> a.getAuthority().equals("ROLE_ADMIN"));
    }
}
```

### 9.3 @PreFilter and @PostFilter

```java
@Service
public class DataService {

    // Filter input collection BEFORE processing
    @PreFilter("filterObject.owner == authentication.name")
    public void deleteDocuments(List<Document> documents) {
        // Only documents owned by current user will be passed
        documentRepository.deleteAll(documents);
    }

    // Filter output collection AFTER execution
    @PostFilter("filterObject.visibility == 'PUBLIC' or filterObject.owner == authentication.name")
    public List<Document> getAllDocuments() {
        return documentRepository.findAll();
        // Returns only public docs or docs owned by user
    }
    
    // Multiple filters
    @PreFilter(value = "filterObject.status == 'PENDING'", filterTarget = "orders")
    public void processOrders(List<Order> orders) { }
}
```

### 9.4 @Secured and @RolesAllowed

```java
@Service
public class LegacyService {

    // Old Spring Security annotation
    @Secured("ROLE_ADMIN")
    public void adminOnly() { }

    @Secured({"ROLE_ADMIN", "ROLE_MANAGER"})
    public void managersAndAdmin() { }

    // JSR-250 annotation (portable)
    @RolesAllowed("ADMIN")  // Automatically adds ROLE_ prefix
    public void anotherAdminOnly() { }

    @RolesAllowed({"ADMIN", "MANAGER"})
    public void multiRole() { }
}
```

### 9.5 Custom Security Expressions

```java
@Component
public class CustomSecurityExpressionRoot extends SecurityExpressionRoot 
        implements MethodSecurityExpressionOperations {

    private Object filterObject;
    private Object returnObject;

    public CustomSecurityExpressionRoot(Authentication authentication) {
        super(authentication);
    }

    // Custom method available in SpEL
    public boolean isOwner(Long resourceId) {
        CustomUserDetails user = (CustomUserDetails) getAuthentication().getPrincipal();
        Resource resource = resourceService.findById(resourceId);
        return resource != null && resource.getOwnerId().equals(user.getUserId());
    }

    public boolean hasSubscription(String tier) {
        CustomUserDetails user = (CustomUserDetails) getAuthentication().getPrincipal();
        return user.getSubscriptionTier().equals(tier);
    }

    // Required interface methods...
    @Override public void setFilterObject(Object o) { this.filterObject = o; }
    @Override public Object getFilterObject() { return filterObject; }
    @Override public void setReturnObject(Object o) { this.returnObject = o; }
    @Override public Object getReturnObject() { return returnObject; }
    @Override public Object getThis() { return this; }
}

// Usage
@PreAuthorize("isOwner(#id) or hasSubscription('PREMIUM')")
public Resource getResource(Long id) { }
```

---

## 10. CORS, CSRF & Security Headers

### 10.1 CORS (Cross-Origin Resource Sharing)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CORS EXPLAINED                                       │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  Same-Origin Policy:                                                    │
│  Browser prevents JavaScript from making requests to different origin   │
│                                                                         │
│  Origin = Protocol + Domain + Port                                      │
│  https://example.com:443 ≠ https://api.example.com:443                  │
│                                                                         │
│  CORS: Server explicitly allows cross-origin requests                   │
│                                                                         │
│  ┌──────────────┐    Preflight (OPTIONS)    ┌──────────────┐            │
│  │   Browser    │ ────────────────────────> │    Server    │            │
│  │ (frontend)   │                           │    (API)     │            │
│  └──────────────┘                           └──────────────┘            │
│                                              Access-Control-Allow-Origin│
│                                              Access-Control-Allow-Methods│
│                                             Access-Control-Allow-Headers│
│  ┌──────────────┐    Actual Request         ┌──────────────┐            │
│  │   Browser    │ ────────────────────────> │    Server    │            │
│  │ (frontend)   │ <──────────────────────── │    (API)     │            │
│  └──────────────┘    Response + CORS headers└──────────────┘            │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 10.2 CORS Configuration

```java
@Configuration
public class CorsConfig {

    @Bean
    public CorsConfigurationSource corsConfigurationSource() {
        CorsConfiguration config = new CorsConfiguration();
        
        // Allowed origins
        config.setAllowedOrigins(List.of(
            "https://frontend.example.com",
            "https://admin.example.com"
        ));
        // Or use patterns
        config.setAllowedOriginPatterns(List.of("https://*.example.com"));
        
        // Allowed methods
        config.setAllowedMethods(List.of("GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"));
        
        // Allowed headers
        config.setAllowedHeaders(List.of("*"));
        // Or specific: List.of("Authorization", "Content-Type", "X-Requested-With")
        
        // Exposed headers (client can read these)
        config.setExposedHeaders(List.of("X-Custom-Header", "Authorization"));
        
        // Allow credentials (cookies, auth headers)
        config.setAllowCredentials(true);
        
        // Preflight cache duration
        config.setMaxAge(3600L);

        UrlBasedCorsConfigurationSource source = new UrlBasedCorsConfigurationSource();
        source.registerCorsConfiguration("/api/**", config);
        return source;
    }
}

// In SecurityFilterChain
@Bean
public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
    return http
        .cors(cors -> cors.configurationSource(corsConfigurationSource()))
        // ... other config
        .build();
}
```

### 10.3 CSRF (Cross-Site Request Forgery)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    CSRF ATTACK                                          │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. User logs into bank.com (session cookie set)                        │
│  2. User visits evil.com (while still logged in)                        │
│  3. evil.com has: <img src="https://bank.com/transfer?to=hacker&amt=1000">│
│  4. Browser sends request WITH bank.com cookies!                        │
│  5. Bank processes transfer (user authenticated via cookie)             │
│                                                                         │
│  Prevention:                                                            │
│  • CSRF Token: Server generates, client must include in requests        │
│  • SameSite Cookies: Browser doesn't send cookie cross-site             │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 10.4 CSRF Configuration

```java
@Bean
public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
    return http
        // Disable for stateless APIs (JWT)
        .csrf(csrf -> csrf.disable())
        
        // OR configure properly for traditional web apps
        .csrf(csrf -> csrf
            .csrfTokenRepository(CookieCsrfTokenRepository.withHttpOnlyFalse())
            .csrfTokenRequestHandler(new SpaCsrfTokenRequestHandler())
            .ignoringRequestMatchers("/api/public/**", "/webhooks/**")
        )
        .build();
}

// For SPAs (Angular, React)
public class SpaCsrfTokenRequestHandler extends CsrfTokenRequestAttributeHandler {
    
    private final CsrfTokenRequestHandler delegate = new XorCsrfTokenRequestAttributeHandler();

    @Override
    public void handle(HttpServletRequest request, HttpServletResponse response, 
                       Supplier<CsrfToken> csrfToken) {
        delegate.handle(request, response, csrfToken);
    }

    @Override
    public String resolveCsrfTokenValue(HttpServletRequest request, CsrfToken csrfToken) {
        // Try header first (SPAs), then parameter (forms)
        String headerValue = request.getHeader(csrfToken.getHeaderName());
        return (headerValue != null) ? super.resolveCsrfTokenValue(request, csrfToken)
                                     : delegate.resolveCsrfTokenValue(request, csrfToken);
    }
}
```

### 10.5 Security Headers

```java
@Bean
public SecurityFilterChain filterChain(HttpSecurity http) throws Exception {
    return http
        .headers(headers -> headers
            // Content Security Policy
            .contentSecurityPolicy(csp -> csp
                .policyDirectives("default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'")
            )
            
            // X-Frame-Options (clickjacking protection)
            .frameOptions(frame -> frame.deny())
            // or .sameOrigin() if you need iframes from same origin
            
            // X-Content-Type-Options
            .contentTypeOptions(Customizer.withDefaults())  // nosniff
            
            // X-XSS-Protection (legacy, CSP is better)
            .xssProtection(xss -> xss.headerValue(XXssProtectionHeaderWriter.HeaderValue.ENABLED_MODE_BLOCK))
            
            // Referrer Policy
            .referrerPolicy(referrer -> referrer.policy(ReferrerPolicy.STRICT_ORIGIN_WHEN_CROSS_ORIGIN))
            
            // Permissions Policy (Feature Policy)
            .permissionsPolicy(permissions -> permissions
                .policy("geolocation=(), microphone=(), camera=()")
            )
            
            // HSTS (HTTPS only)
            .httpStrictTransportSecurity(hsts -> hsts
                .includeSubDomains(true)
                .maxAgeInSeconds(31536000)
                .preload(true)
            )
        )
        .build();
}
```

### 10.6 Security Headers Reference

| Header | Purpose | Value |
|---|---|---|
| `X-Frame-Options` | Prevent clickjacking | `DENY`, `SAMEORIGIN` |
| `X-Content-Type-Options` | Prevent MIME sniffing | `nosniff` |
| `X-XSS-Protection` | XSS filter (legacy) | `1; mode=block` |
| `Content-Security-Policy` | Control resource loading | Various directives |
| `Strict-Transport-Security` | Force HTTPS | `max-age=31536000; includeSubDomains` |
| `Referrer-Policy` | Control referer header | `strict-origin-when-cross-origin` |
| `Permissions-Policy` | Control browser features | Feature list |

---

## 11. Security Testing

### 11.1 Testing with Spring Security

```java
@WebMvcTest(UserController.class)
@Import(SecurityConfig.class)
class UserControllerSecurityTest {

    @Autowired
    private MockMvc mockMvc;

    @MockBean
    private UserService userService;

    // Test unauthenticated access
    @Test
    void getUsers_unauthenticated_returns401() throws Exception {
        mockMvc.perform(get("/api/users"))
            .andExpect(status().isUnauthorized());
    }

    // Test with mock user
    @Test
    @WithMockUser(username = "john", roles = "USER")
    void getUsers_authenticated_returns200() throws Exception {
        when(userService.findAll(any())).thenReturn(Page.empty());
        
        mockMvc.perform(get("/api/users"))
            .andExpect(status().isOk());
    }

    // Test with specific authorities
    @Test
    @WithMockUser(authorities = "VIEW_REPORTS")
    void getReport_withAuthority_returns200() throws Exception {
        mockMvc.perform(get("/api/reports/1"))
            .andExpect(status().isOk());
    }

    // Test admin-only endpoint
    @Test
    @WithMockUser(roles = "USER")
    void deleteUser_asUser_returns403() throws Exception {
        mockMvc.perform(delete("/api/users/1"))
            .andExpect(status().isForbidden());
    }

    @Test
    @WithMockUser(roles = "ADMIN")
    void deleteUser_asAdmin_returns204() throws Exception {
        mockMvc.perform(delete("/api/users/1"))
            .andExpect(status().isNoContent());
    }
}
```

### 11.2 Custom User for Tests

```java
// Custom annotation
@Retention(RetentionPolicy.RUNTIME)
@WithSecurityContext(factory = WithMockCustomUserSecurityContextFactory.class)
public @interface WithMockCustomUser {
    String username() default "testuser";
    String email() default "test@example.com";
    long id() default 1L;
    String[] roles() default {"USER"};
}

// Security context factory
public class WithMockCustomUserSecurityContextFactory 
        implements WithSecurityContextFactory<WithMockCustomUser> {

    @Override
    public SecurityContext createSecurityContext(WithMockCustomUser annotation) {
        SecurityContext context = SecurityContextHolder.createEmptyContext();

        User user = new User();
        user.setId(annotation.id());
        user.setUsername(annotation.username());
        user.setEmail(annotation.email());

        CustomUserDetails principal = new CustomUserDetails(user);
        
        List<GrantedAuthority> authorities = Arrays.stream(annotation.roles())
            .map(role -> new SimpleGrantedAuthority("ROLE_" + role))
            .collect(Collectors.toList());

        Authentication auth = new UsernamePasswordAuthenticationToken(
            principal, null, authorities
        );
        
        context.setAuthentication(auth);
        return context;
    }
}

// Usage
@Test
@WithMockCustomUser(id = 123, username = "john", roles = {"ADMIN", "MANAGER"})
void testWithCustomUser() throws Exception {
    // Test with custom user details
}
```

### 11.3 Testing JWT Authentication

```java
@SpringBootTest
@AutoConfigureMockMvc
class JwtAuthenticationTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private JwtService jwtService;

    private String validToken;

    @BeforeEach
    void setup() {
        UserDetails user = User.builder()
            .username("test@example.com")
            .password("password")
            .authorities("ROLE_USER")
            .build();
        validToken = jwtService.generateAccessToken(user);
    }

    @Test
    void protectedEndpoint_withValidToken_returns200() throws Exception {
        mockMvc.perform(get("/api/protected")
                .header("Authorization", "Bearer " + validToken))
            .andExpect(status().isOk());
    }

    @Test
    void protectedEndpoint_withExpiredToken_returns401() throws Exception {
        String expiredToken = createExpiredToken();
        
        mockMvc.perform(get("/api/protected")
                .header("Authorization", "Bearer " + expiredToken))
            .andExpect(status().isUnauthorized())
            .andExpect(jsonPath("$.message").value("Token expired"));
    }

    @Test
    void protectedEndpoint_withInvalidToken_returns401() throws Exception {
        mockMvc.perform(get("/api/protected")
                .header("Authorization", "Bearer invalid.token.here"))
            .andExpect(status().isUnauthorized());
    }

    @Test
    void protectedEndpoint_withoutToken_returns401() throws Exception {
        mockMvc.perform(get("/api/protected"))
            .andExpect(status().isUnauthorized());
    }
}
```

### 11.4 Testing Method Security

```java
@SpringBootTest
class MethodSecurityTest {

    @Autowired
    private DocumentService documentService;

    @Test
    @WithMockUser(roles = "USER")
    void deleteAll_asUser_throwsAccessDenied() {
        assertThrows(AccessDeniedException.class, () -> {
            documentService.deleteAll();
        });
    }

    @Test
    @WithMockUser(roles = "ADMIN")
    void deleteAll_asAdmin_succeeds() {
        assertDoesNotThrow(() -> documentService.deleteAll());
    }

    @Test
    @WithMockCustomUser(id = 123)
    void getUserDocuments_ownDocuments_succeeds() {
        List<Document> docs = documentService.getUserDocuments(123L);
        assertThat(docs).isNotNull();
    }

    @Test
    @WithMockCustomUser(id = 123)
    void getUserDocuments_otherUserDocuments_throwsAccessDenied() {
        assertThrows(AccessDeniedException.class, () -> {
            documentService.getUserDocuments(456L);  // Different user
        });
    }
}
```

---

## 12. Common Vulnerabilities & Prevention

### 12.1 OWASP Top 10 and Spring Security

| Vulnerability | Description | Spring Security Mitigation |
|---|---|---|
| **A01: Broken Access Control** | Missing authorization checks | Method security, URL authorization |
| **A02: Cryptographic Failures** | Weak encryption | Strong PasswordEncoder, HTTPS |
| **A03: Injection** | SQL, LDAP injection | Not directly, use parameterized queries |
| **A04: Insecure Design** | Missing security controls | Defense in depth |
| **A05: Security Misconfiguration** | Default/weak settings | Secure defaults, explicit config |
| **A06: Vulnerable Components** | Outdated libraries | Regular updates |
| **A07: Auth Failures** | Broken authentication | Proper auth configuration |
| **A08: Integrity Failures** | Unsigned data | CSRF protection, signed tokens |
| **A09: Logging Failures** | Missing audit logs | Security event logging |
| **A10: SSRF** | Server-side request forgery | Not directly, validate URLs |

### 12.1.1 Rate Limiting — Theory and Algorithms

**Protecting your application from abuse and denial of service.**

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    RATE LIMITING ALGORITHMS                             │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  1. FIXED WINDOW                                                        │
│  ───────────────                                                        │
│  Count requests per time window (e.g., 100 requests/minute)             │
│                                                                         │
│  │ Minute 1: 100 allowed  │ Minute 2: 100 allowed  │                    │
│  │ ████████████████████   │ ████████████████████   │                    │
│  └────────────────────────┴────────────────────────┘                    │
│                                                                         │
│  Problem: Burst at window boundary (200 requests in 1 second span)      │
│           └─ 100 at :59 ─┘└─ 100 at :00 ─┘                              │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  2. SLIDING WINDOW LOG                                                  │
│  ─────────────────────                                                  │
│  Store timestamp of each request, count in rolling window               │
│                                                                         │
│  Requests: [12:00:05, 12:00:23, 12:00:45, 12:01:02]                      │
│  Now: 12:01:30                                                          │
│  Window: 12:00:30 to 12:01:30 → Count: 2 (12:00:45, 12:01:02)           │
│                                                                         │
│  Pros: Accurate                                                         │
│  Cons: Memory-intensive (stores all timestamps)                         │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  3. SLIDING WINDOW COUNTER                                              │
│  ─────────────────────────                                              │
│  Weighted average of current and previous window                        │
│                                                                         │
│  Previous window: 80 requests                                           │
│  Current window: 40 requests (60% elapsed)                              │
│  Estimate: 80 * 0.4 + 40 = 72 requests                                  │
│                                                                         │
│  Pros: Memory-efficient, smooth                                         │
│  Cons: Approximate                                                      │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  4. TOKEN BUCKET                                                        │
│  ────────────────                                                       │
│                                                                         │
│      ┌─────────────┐                                                    │
│      │  ● ● ● ● ●  │  Bucket (capacity: 10 tokens)                      │
│      │  ● ● ● ● ●  │                                                    │
│      └──────┬──────┘                                                    │
│             │ Take 1 token per request                                  │
│             ▼                                                           │
│         Request allowed if tokens > 0                                   │
│                                                                         │
│      Refill: 10 tokens per second (rate)                                │
│                                                                         │
│  Pros: Allows controlled bursts (use accumulated tokens)                │
│  Cons: Slightly more complex to implement                               │
│                                                                         │
│  ─────────────────────────────────────────────────────────────────────  │
│                                                                         │
│  5. LEAKY BUCKET                                                        │
│  ───────────────                                                        │
│                                                                         │
│      ┌─────────────┐                                                    │
│      │  ↓ ↓ ↓ ↓ ↓  │  Requests queue up                                 │
│      │  ● ● ● ● ●  │                                                    │
│      └──────┬──────┘                                                    │
│             │ Processes at constant rate                                │
│             ● → ● → ● → (output)                                        │
│                                                                         │
│  Pros: Smooth output rate, no bursts                                    │
│  Cons: Can add latency, requests may queue                              │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

**Rate Limiting in Spring Security:**

```java
@Component
public class RateLimitingFilter extends OncePerRequestFilter {

    private final Map<String, Bucket> buckets = new ConcurrentHashMap<>();

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain filterChain) 
            throws ServletException, IOException {
        
        String key = getClientKey(request);  // IP or user ID
        Bucket bucket = buckets.computeIfAbsent(key, k -> createBucket());

        if (bucket.tryConsume(1)) {
            filterChain.doFilter(request, response);
        } else {
            response.setStatus(429);  // Too Many Requests
            response.setHeader("Retry-After", "60");
            response.getWriter().write("Rate limit exceeded");
        }
    }

    private Bucket createBucket() {
        // 100 requests per minute, burst of 20
        return Bucket.builder()
            .addLimit(Bandwidth.classic(20, Refill.greedy(100, Duration.ofMinutes(1))))
            .build();
    }

    private String getClientKey(HttpServletRequest request) {
        // Prefer user ID if authenticated
        Authentication auth = SecurityContextHolder.getContext().getAuthentication();
        if (auth != null && auth.isAuthenticated() && !"anonymousUser".equals(auth.getPrincipal())) {
            return "user:" + auth.getName();
        }
        // Fall back to IP (consider X-Forwarded-For behind proxy)
        return "ip:" + request.getRemoteAddr();
    }
}
```

**Rate Limiting Strategies:**

| Strategy | When to Use |
|---|---|
| Per IP | Anonymous endpoints, public APIs |
| Per User | Authenticated endpoints |
| Per API Key | Third-party integrations |
| Per Endpoint | Different limits for different operations |
| Adaptive | Reduce limits during high load |

### 12.1.2 API Security Best Practices

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    API SECURITY CHECKLIST                               │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  TRANSPORT                                                              │
│  ─────────                                                              │
│  ✅ TLS 1.2+ for all endpoints                                          │
│  ✅ HSTS header (Strict-Transport-Security)                             │
│  ✅ Certificate pinning for mobile apps                                 │
│                                                                         │
│  AUTHENTICATION                                                         │
│  ──────────────                                                         │
│  ✅ Use OAuth2/JWT for APIs (not sessions)                              │
│  ✅ Short token expiration (15 minutes)                                 │
│  ✅ Secure token storage (HttpOnly cookies or secure storage)           │
│  ✅ Implement token refresh mechanism                                   │
│                                                                         │
│  AUTHORIZATION                                                          │
│  ─────────────                                                          │
│  ✅ Validate permissions on every request                               │
│  ✅ Use principle of least privilege                                    │
│  ✅ Implement object-level authorization (IDOR prevention)              │
│                                                                         │
│  INPUT VALIDATION                                                       │
│  ────────────────                                                       │
│  ✅ Validate all input (type, length, format, range)                    │
│  ✅ Use allowlists, not blocklists                                      │
│  ✅ Parameterized queries (prevent injection)                           │
│  ✅ Sanitize output (prevent XSS)                                       │
│                                                                         │
│  RATE LIMITING                                                          │
│  ─────────────                                                          │
│  ✅ Implement per-user and per-IP limits                                │
│  ✅ Return 429 Too Many Requests                                        │
│  ✅ Include Retry-After header                                          │
│                                                                         │
│  ERROR HANDLING                                                         │
│  ──────────────                                                         │
│  ✅ Generic error messages to clients                                   │
│  ✅ Log detailed errors internally                                      │
│  ✅ Don't leak stack traces, versions, or internal IPs                  │
│                                                                         │
│  LOGGING & MONITORING                                                   │
│  ────────────────────                                                   │
│  ✅ Log all authentication events                                       │
│  ✅ Alert on anomalies (failed logins, rate limit hits)                 │
│  ✅ Audit trail for sensitive operations                                │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 12.2 Preventing Brute Force Attacks

```java
@Component
public class LoginAttemptService {

    private final LoadingCache<String, Integer> attemptsCache;

    public LoginAttemptService() {
        this.attemptsCache = CacheBuilder.newBuilder()
            .expireAfterWrite(15, TimeUnit.MINUTES)
            .build(new CacheLoader<>() {
                @Override
                public Integer load(String key) {
                    return 0;
                }
            });
    }

    public void loginSucceeded(String key) {
        attemptsCache.invalidate(key);
    }

    public void loginFailed(String key) {
        int attempts = attemptsCache.getUnchecked(key);
        attemptsCache.put(key, attempts + 1);
    }

    public boolean isBlocked(String key) {
        return attemptsCache.getUnchecked(key) >= 5;
    }
}

// Authentication listener
@Component
public class AuthenticationEventListener {

    @Autowired
    private LoginAttemptService loginAttemptService;

    @EventListener
    public void onSuccess(AuthenticationSuccessEvent event) {
        WebAuthenticationDetails details = 
            (WebAuthenticationDetails) event.getAuthentication().getDetails();
        loginAttemptService.loginSucceeded(details.getRemoteAddress());
    }

    @EventListener
    public void onFailure(AuthenticationFailureBadCredentialsEvent event) {
        WebAuthenticationDetails details = 
            (WebAuthenticationDetails) event.getAuthentication().getDetails();
        loginAttemptService.loginFailed(details.getRemoteAddress());
    }
}

// In UserDetailsService
@Override
public UserDetails loadUserByUsername(String username) {
    String ip = getClientIP();
    if (loginAttemptService.isBlocked(ip)) {
        throw new LockedException("IP blocked due to too many failed attempts");
    }
    // ... normal loading
}
```

### 12.3 Secure Headers Implementation

```java
@Component
public class SecurityHeadersFilter extends OncePerRequestFilter {

    @Override
    protected void doFilterInternal(HttpServletRequest request, 
                                    HttpServletResponse response, 
                                    FilterChain filterChain) 
            throws ServletException, IOException {
        
        // Content Security Policy
        response.setHeader("Content-Security-Policy", 
            "default-src 'self'; " +
            "script-src 'self' 'unsafe-inline' https://cdn.example.com; " +
            "style-src 'self' 'unsafe-inline'; " +
            "img-src 'self' data: https:; " +
            "font-src 'self' https://fonts.gstatic.com; " +
            "connect-src 'self' https://api.example.com; " +
            "frame-ancestors 'none'; " +
            "form-action 'self'");
        
        // Other security headers
        response.setHeader("X-Frame-Options", "DENY");
        response.setHeader("X-Content-Type-Options", "nosniff");
        response.setHeader("X-XSS-Protection", "1; mode=block");
        response.setHeader("Referrer-Policy", "strict-origin-when-cross-origin");
        response.setHeader("Permissions-Policy", 
            "geolocation=(), microphone=(), camera=()");
        
        filterChain.doFilter(request, response);
    }
}
```

### 12.4 Secrets Management

```java
// ❌ DON'T: Hardcode secrets
private static final String JWT_SECRET = "my-super-secret-key";

// ✅ DO: Use environment variables
@Value("${jwt.secret}")
private String jwtSecret;

// ✅ DO: Use Spring Cloud Vault
@Configuration
public class VaultConfig {
    // Secrets loaded automatically from Vault
}

// ✅ DO: Use AWS Secrets Manager
@Bean
public String jwtSecret(SecretsManagerClient client) {
    GetSecretValueRequest request = GetSecretValueRequest.builder()
        .secretId("production/jwt-secret")
        .build();
    return client.getSecretValue(request).secretString();
}
```

---

## 13. Production Security Checklist

### 13.1 Security Configuration Checklist

```markdown
## Authentication
□ Use strong password encoder (BCrypt with cost ≥ 12)
□ Implement account lockout after failed attempts
□ Force password change on first login
□ Implement proper logout (invalidate tokens/sessions)
□ Secure remember-me with persistent tokens

## Authorization
□ Apply least privilege principle
□ Use method-level security for sensitive operations
□ Audit authorization decisions
□ Implement proper role hierarchy

## Session Management
□ Enable session fixation protection
□ Limit concurrent sessions
□ Use secure, HttpOnly cookies
□ Implement session timeout
□ Use HTTPS only (secure cookies)

## API Security
□ Implement rate limiting
□ Use JWT with short expiration
□ Implement token refresh mechanism
□ Validate all input
□ Return generic error messages

## Headers & CORS
□ Configure Content-Security-Policy
□ Enable HSTS
□ Set X-Frame-Options
□ Configure CORS properly
□ Set Referrer-Policy

## Monitoring & Logging
□ Log authentication events
□ Monitor failed login attempts
□ Set up security alerts
□ Regular security audits
□ Penetration testing

## Dependencies
□ Keep Spring Security updated
□ Scan for vulnerabilities (OWASP, Snyk)
□ Remove unused dependencies
□ Update all transitive dependencies
```

### 13.2 Secure Application Properties

```yaml
# production profile
spring:
  profiles: prod

server:
  ssl:
    enabled: true
    key-store: ${SSL_KEYSTORE_PATH}
    key-store-password: ${SSL_KEYSTORE_PASSWORD}
  servlet:
    session:
      cookie:
        secure: true
        http-only: true
        same-site: strict
      timeout: 30m

# Security properties
security:
  password-encoder:
    strength: 12
  jwt:
    secret: ${JWT_SECRET}
    access-token-expiration: 15m
    refresh-token-expiration: 7d
  rate-limit:
    enabled: true
    max-requests: 100
    duration: 1m
  login:
    max-attempts: 5
    lockout-duration: 15m

# Never log sensitive data
logging:
  level:
    org.springframework.security: WARN
```

---

## 14. Interview Questions by Experience Level

### 14.1 Junior (0–2 Years)

**Q1: What is the difference between authentication and authorization?**
> Authentication verifies WHO you are (identity), typically via username/password. Authorization determines WHAT you can do (permissions), checked after authentication. Example: Logging in is authentication; accessing /admin is authorization.

**Q2: What happens when you add `spring-boot-starter-security` to your project?**
> Spring Security auto-configures: all endpoints require authentication, provides default login/logout pages, generates random password (logged on startup), enables CSRF protection, and adds security headers.

**Q3: What is `UserDetailsService` and why is it important?**
> `UserDetailsService` is an interface with one method: `loadUserByUsername(String username)`. It's how Spring Security loads user data during authentication. You implement it to connect to your user database.

**Q4: What is the purpose of `PasswordEncoder`?**
> `PasswordEncoder` hashes passwords for storage and verifies passwords during login. Never store plain text passwords. `BCryptPasswordEncoder` is the recommended default.

**Q5: How do you permit certain URLs without authentication?**
> Use `.requestMatchers("/public/**").permitAll()` in your `SecurityFilterChain` configuration.

**Q6: What are roles in Spring Security?**
> Roles are authorities with `ROLE_` prefix. `hasRole("ADMIN")` checks for `ROLE_ADMIN`. They represent coarse-grained permissions.

---

### 14.2 Mid-Level (2–5 Years)

**Q7: Explain the Spring Security filter chain.**
> Requests pass through a chain of filters in order. Each filter handles specific security concerns: `SecurityContextPersistenceFilter` (loads/saves security context), `UsernamePasswordAuthenticationFilter` (processes login forms), `ExceptionTranslationFilter` (handles security exceptions), `AuthorizationFilter` (enforces access rules).

**Q8: What is the `SecurityContext` and how is it stored?**
> `SecurityContext` holds the `Authentication` object for the current user. By default, it's stored in `ThreadLocal` via `SecurityContextHolder`, meaning each thread has its own context. For async operations, use `MODE_INHERITABLETHREADLOCAL`.

**Q9: How would you implement JWT authentication in Spring Security?**
> Create a filter extending `OncePerRequestFilter` that: extracts JWT from Authorization header, validates signature and expiration, loads UserDetails, creates `UsernamePasswordAuthenticationToken`, and sets it in `SecurityContextHolder`. Register filter before `UsernamePasswordAuthenticationFilter`.

**Q10: Explain the difference between `@PreAuthorize` and `@Secured`.**
> `@Secured` is simpler, only accepts role names: `@Secured("ROLE_ADMIN")`. `@PreAuthorize` supports SpEL expressions: `@PreAuthorize("hasRole('ADMIN') and #userId == authentication.principal.id")`. Use `@PreAuthorize` for complex authorization logic.

**Q11: What is CSRF and how does Spring Security protect against it?**
> CSRF is when an attacker tricks a user's browser into making unwanted requests using their session. Spring Security generates a CSRF token that must be included in state-changing requests (POST, PUT, DELETE). The server validates the token matches.

**Q12: How do you configure multiple authentication mechanisms?**
> Create multiple `AuthenticationProvider` beans and register them with `ProviderManager`. The manager tries each provider in order until one succeeds. Example: try DB auth first, then LDAP.

**Q13: Explain OAuth2 Authorization Code flow.**
> 1) User clicks "Login with Google", 2) Redirect to Google auth server with client_id, 3) User logs in and consents, 4) Google redirects back with authorization code, 5) Backend exchanges code for access/refresh tokens, 6) Backend uses access token to get user info.

---

### 14.3 Senior / Lead (5+ Years)

**Q14: Design a multi-tenant security architecture with Spring Security.**
> Implement a `TenantContext` stored in `ThreadLocal`, populated by filter from JWT claim or subdomain. Create custom `AuthenticationProvider` that validates tenant membership. Use `@PreAuthorize` with custom method: `@PreAuthorize("@tenantSecurity.belongsToTenant(#tenantId)")`. Consider row-level security in database layer.

**Q15: How would you implement token revocation for JWTs?**
> JWTs are stateless, so revocation requires additional infrastructure. Options: 1) Short expiration + refresh tokens (revoke refresh token), 2) Token blacklist in Redis with TTL, 3) Token versioning (store version per user, invalidate all tokens by incrementing), 4) JTI claim with database lookup.

**Q16: Explain how to secure a microservices architecture.**
> 1) API Gateway handles authentication, validates JWT, 2) Services trust the gateway (internal network), 3) Service-to-service: use OAuth2 Client Credentials flow or mutual TLS, 4) Propagate security context via headers (X-User-Id, X-Roles), 5) Each service still validates authorization.

**Q17: How do you handle security in reactive applications?**
> Use `@EnableWebFluxSecurity` instead of `@EnableWebSecurity`. Security context stored in Reactor Context, not ThreadLocal. Return `Mono<SecurityFilterChain>`. Use `ReactiveAuthenticationManager`, `ReactiveUserDetailsService`. For propagation: use `ReactorContextWebFilter`.

**Q18: Design a zero-trust security model for a Spring Boot application.**
> 1) Never trust network location — authenticate every request, 2) Short-lived tokens with frequent rotation, 3) Mutual TLS between services, 4) Encrypt all traffic (TLS 1.3), 5) Fine-grained authorization at every layer (API, method, data), 6) Continuous validation — re-check permissions periodically, 7) Immutable audit log of all access.

**Q19: How do you implement fine-grained authorization like ABAC?**
> Create custom `PermissionEvaluator` implementing complex attribute checks. Use `@PreAuthorize("hasPermission(#resource, 'read')")`. Evaluator checks user attributes, resource attributes, action, and environment (time, IP). Consider policy engines like OPA or Spring Authorization Server for complex scenarios.

**Q20: Explain the security implications of different session management strategies.**
> `SessionCreationPolicy.ALWAYS`: Always creates session, risk of session fixation, requires session affinity in LB. `STATELESS`: No session for security (JWT), but can't use CSRF token, need alternative protection. `IF_REQUIRED`: Creates on demand, default choice for web apps. Consider distributed session store (Redis) for horizontal scaling.

---

### 14.4 Scenario-Based Questions

**Q21: Your application is getting brute force attacks. How do you mitigate?**
> 1) Implement account lockout after N failed attempts (with LoginAttemptService), 2) Add rate limiting per IP (Spring Cloud Gateway or filter), 3) CAPTCHA after failed attempts, 4) Exponential backoff, 5) Monitor and alert on patterns, 6) Consider WAF rules.

**Q22: Users are complaining about being logged out unexpectedly. How do you debug?**
> 1) Check session timeout configuration, 2) Look for session fixation protection triggering, 3) Check if cookies are lost (secure cookie on HTTP?), 4) Verify Redis/session store connectivity, 5) Check for concurrent session limit, 6) Review logs for session invalidation events.

**Q23: How would you migrate from session-based to JWT authentication?**
> 1) Add JWT infrastructure alongside session auth, 2) Return JWT on login while keeping session, 3) Accept both auth methods in filter chain, 4) Migrate clients to use JWT, 5) Monitor adoption via metrics, 6) Remove session auth when migration complete, 7) Plan for token refresh and revocation.

**Q24: A security audit found that tokens don't expire for 30 days. How do you fix?**
> 1) Reduce access token expiration to 15 minutes, 2) Implement refresh token rotation (7 days), 3) Add token versioning in user table, 4) Implement logout to revoke refresh tokens, 5) Create admin endpoint to revoke all user tokens, 6) Add monitoring for token age.

**Q25: How would you implement row-level security in Spring Security?**
> 1) Use `@PostFilter` for collection filtering, 2) Implement custom `Specification` that adds user-based predicates, 3) Use Hibernate `@Filter` with security context, 4) Create repository methods that implicitly filter by owner, 5) Consider database-level RLS (PostgreSQL RLS) for defense in depth.

---

## Quick Reference — Security Annotations

| Annotation | Purpose | Example |
|---|---|---|
| `@EnableWebSecurity` | Enable Spring Security | Class level |
| `@EnableMethodSecurity` | Enable method security | Class level |
| `@PreAuthorize` | Check BEFORE method | `@PreAuthorize("hasRole('ADMIN')")` |
| `@PostAuthorize` | Check AFTER method | `@PostAuthorize("returnObject.owner == principal.username")` |
| `@PreFilter` | Filter input collection | `@PreFilter("filterObject.owner == principal.username")` |
| `@PostFilter` | Filter output collection | `@PostFilter("filterObject.public")` |
| `@Secured` | Simple role check | `@Secured("ROLE_ADMIN")` |
| `@RolesAllowed` | JSR-250 role check | `@RolesAllowed("ADMIN")` |
| `@WithMockUser` | Test with mock user | `@WithMockUser(roles="ADMIN")` |
| `@AuthenticationPrincipal` | Inject current user | `@AuthenticationPrincipal UserDetails user` |

---

## Quick Reference — Common Configurations

```java
// Stateless API (JWT)
http.csrf(csrf -> csrf.disable())
    .sessionManagement(s -> s.sessionCreationPolicy(SessionCreationPolicy.STATELESS))
    .authorizeHttpRequests(auth -> auth.anyRequest().authenticated())
    .oauth2ResourceServer(oauth2 -> oauth2.jwt(Customizer.withDefaults()));

// Traditional Web App
http.authorizeHttpRequests(auth -> auth
        .requestMatchers("/login", "/css/**").permitAll()
        .anyRequest().authenticated())
    .formLogin(Customizer.withDefaults())
    .logout(Customizer.withDefaults());

// OAuth2 Login
http.oauth2Login(oauth2 -> oauth2
        .loginPage("/login")
        .defaultSuccessUrl("/dashboard"))
    .oauth2Client(Customizer.withDefaults());

// Multiple Filter Chains
@Order(1) // API
http.securityMatcher("/api/**")
    .csrf(csrf -> csrf.disable())
    .sessionManagement(s -> s.sessionCreationPolicy(SessionCreationPolicy.STATELESS));

@Order(2) // Web
http.authorizeHttpRequests(auth -> auth.anyRequest().authenticated())
    .formLogin(Customizer.withDefaults());
```

---

> **Related Guides:**  
> - [Spring Framework Complete Guide](./Spring-Framework-Complete-Guide.md)  
> - [Spring Boot Complete Guide](./Spring-Boot-Complete-Guide.md)
