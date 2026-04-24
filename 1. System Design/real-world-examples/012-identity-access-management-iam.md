# Complete System Design: Identity & Access Management (IAM) (Production-Ready)

> **Complexity Level:** Advanced  
> **Estimated Time:** 45-60 minutes in interview  
> **Real-World Examples:** AWS IAM, Auth0, Okta, Keycloak, Google Cloud IAM

---

## Table of Contents
1. [Problem Statement](#1-problem-statement)
2. [Requirements Clarification](#2-requirements-clarification)
3. [Scale Estimation](#3-scale-estimation)
4. [High-Level Design](#4-high-level-design)
5. [Deep Dive: Core Components](#5-deep-dive-core-components)
6. [Deep Dive: Database Design](#6-deep-dive-database-design)
7. [Deep Dive: Authentication & Authorization Engine](#7-deep-dive-authentication--authorization-engine)
8. [Scaling Strategies](#8-scaling-strategies)
9. [Failure Scenarios & Mitigation](#9-failure-scenarios--mitigation)
10. [Monitoring & Observability](#10-monitoring--observability)
11. [Advanced Features](#11-advanced-features)
12. [Interview Q&A](#12-interview-qa)
13. [Production Checklist](#13-production-checklist)

---

## 1. Problem Statement

**Initial Question:**  
"Design an Identity and Access Management (IAM) system like AWS IAM that handles authentication, authorization, and access control for a cloud platform with millions of users."

**Interviewer's Perspective:**  
This problem assesses:
- Security design principles (defense in depth)
- Token management (JWT, OAuth 2.0, OIDC)
- Authorization models (RBAC, ABAC)
- Policy evaluation engine design
- Scalability of auth systems in the critical path
- Compliance and audit requirements

---

## 2. Requirements Clarification

### Interview Dialog

**Candidate:** "Let me clarify the requirements before designing."

### 2.1 Functional Requirements

**Candidate:** "For functional requirements:
1. What authentication methods do we need — password, social login, SSO?
2. Do we need multi-factor authentication (MFA)?
3. What authorization model — role-based (RBAC), attribute-based (ABAC), or both?
4. Do we need OAuth 2.0/OpenID Connect provider capabilities?
5. Do we need API key management for service-to-service auth?
6. Is audit logging required for compliance?"

**Interviewer:** "Let's support:
- Username/password authentication
- Social login (Google, GitHub)
- MFA (TOTP-based, like Google Authenticator)
- RBAC as the primary model, with ABAC for fine-grained policies
- OAuth 2.0/OIDC for third-party app authorization
- API key management
- Full audit logging"

**Candidate:** "Core features:
1. ✅ User registration and login (password + social login)
2. ✅ Multi-factor authentication (TOTP)
3. ✅ Role-Based Access Control (RBAC) with role hierarchy
4. ✅ Attribute-Based Access Control (ABAC) for fine-grained policies
5. ✅ OAuth 2.0 / OpenID Connect provider
6. ✅ API key management for programmatic access
7. ✅ Service-to-service authentication
8. ✅ Comprehensive audit logging"

### 2.2 Non-Functional Requirements

**Candidate:** "For non-functional requirements:
1. What latency is acceptable for authentication and authorization?
2. What's the scale — users, API calls needing authorization?
3. Availability requirements? Auth is in the critical path.
4. Multi-tenant support?"

**Interviewer:**
- Authentication: <100ms for login flow
- Authorization: <10ms per API call (evaluated on every request)
- Scale: 100M users, 1M authorization checks/sec
- Availability: 99.999% (auth is critical infrastructure)
- Multi-tenant: yes, each organization manages their own users/roles

**Candidate:** "Summary:
- **Scale:** 100M users, 1M auth checks/sec
- **Auth latency:** <100ms for login, <10ms for authorization
- **Availability:** 99.999% (5.26 minutes downtime/year)
- **Security:** Zero tolerance for unauthorized access
- **Compliance:** Full audit trail, GDPR/SOC2 compatible
- **Multi-tenant:** Organization-level isolation"

---

## 3. Scale Estimation

### 3.1 Traffic Estimation

**Candidate:** "Let me estimate the traffic:

**Authentication (Login/Token):**
- Logins per day: assume 10% of users log in daily = 10M logins/day
- Logins per second: 10M / 86,400 ≈ **115 logins/sec**
- Peak (3x): **~350 logins/sec**
- Token refresh: 10M active sessions × refresh every 15 min = **11K refreshes/sec**

**Authorization (Policy Evaluation):**
- 1M authorization checks/sec (every API call in the platform)
- This is the hot path — must be <10ms
- Peak (3x): **3M checks/sec**

**Audit Logging:**
- 1M events/sec (one per authorization check + login events)
- 1M × 200 bytes = **200 MB/sec** write throughput"

### 3.2 Storage Estimation

**Candidate:** "For storage:

**User Data:**
- 100M users × 2 KB per user = **200 GB**
- Includes: credentials, profile, MFA secrets

**Roles & Policies:**
- 10K roles × 1 KB = 10 MB
- 100K policies × 5 KB = 500 MB
- Very small, easily cached

**Sessions/Tokens:**
- 10M active sessions × 1 KB = **10 GB** (Redis)

**Audit Logs:**
- 200 MB/sec × 86,400 sec = **17 TB/day**
- 30-day retention: **510 TB**
- Archived to cold storage (S3) after 30 days"

### 3.3 Cache Estimation

**Candidate:** "Authorization must be <10ms, so caching is critical:

**Policy Cache:**
- 100K policies × 5 KB = 500 MB — easily fits in Redis
- Cache per-role permission set: 10K roles × 10 KB = 100 MB

**User-Role Mapping Cache:**
- 100M users × 50 bytes (list of role IDs) = 5 GB
- Partition across Redis cluster

**Token Validation Cache:**
- Cache validated JWT claims (avoid re-parsing): 10M × 500 bytes = 5 GB"

---

## 4. High-Level Design

### 4.1 Architecture Diagram

```
┌──────────────┐
│   Clients    │
│(Web/Mobile/  │
│ Services)    │
└──────┬───────┘
       │ HTTPS
       ▼
┌──────────────────────────────────────────┐
│           API Gateway                     │
│  (Token validation on every request)      │
│  ┌──────────────────────────────┐        │
│  │ JWT Validation (local, no    │        │
│  │ round-trip to auth service)  │        │
│  │ → Extract claims → Check     │        │
│  │   permissions via Policy     │        │
│  │   Engine SDK                 │        │
│  └──────────────────────────────┘        │
└──────────────┬───────────────────────────┘
               │ Authorized requests only
               ▼
┌──────────────────────────┐
│     Resource Servers     │
│   (Microservices)        │
└──────────────────────────┘

  Auth Infrastructure:

┌─────────────────────────────────────────────────────┐
│                                                     │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │  Auth        │  │  Token       │  │  Policy   │ │
│  │  Service     │  │  Service     │  │  Engine   │ │
│  │              │  │              │  │           │ │
│  │  - Login     │  │  - JWT issue │  │  - RBAC   │ │
│  │  - Register  │  │  - Refresh   │  │  - ABAC   │ │
│  │  - MFA       │  │  - Revoke    │  │  - Eval   │ │
│  │  - Social    │  │  - JWKS      │  │  - Cache  │ │
│  └──────┬───────┘  └──────┬───────┘  └─────┬─────┘ │
│         │                 │                 │       │
│         ▼                 ▼                 ▼       │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │  User        │  │  Redis       │  │  Policy   │ │
│  │  Directory   │  │  (Sessions,  │  │  Store    │ │
│  │  (PostgreSQL)│  │   Tokens,    │  │  (PgSQL)  │ │
│  │              │  │   Cache)     │  │           │ │
│  └──────────────┘  └──────────────┘  └───────────┘ │
│                                                     │
│  ┌──────────────┐  ┌──────────────┐                 │
│  │   Audit      │  │   Key        │                 │
│  │   Service    │  │   Management │                 │
│  │  (Kafka →    │  │   Service    │                 │
│  │   S3/ES)     │  │  (HSM/Vault) │                 │
│  └──────────────┘  └──────────────┘                 │
└─────────────────────────────────────────────────────┘
```

### 4.2 API Design

**Candidate:** "Core APIs:

**Authentication APIs:**
```http
POST /auth/register
{
  "email": "user@example.com",
  "password": "secureP@ss123",
  "displayName": "John Doe"
}
→ 201 Created { "userId": "usr_123", "email": "..." }

POST /auth/login
{
  "email": "user@example.com",
  "password": "secureP@ss123"
}
→ 200 OK {
  "accessToken": "eyJhbG...",
  "refreshToken": "dGhpcyBpcyBh...",
  "expiresIn": 900,
  "tokenType": "Bearer",
  "mfaRequired": true,
  "mfaToken": "mfa_temp_456"  // if MFA enabled
}

POST /auth/mfa/verify
{
  "mfaToken": "mfa_temp_456",
  "code": "123456"
}
→ 200 OK { "accessToken": "...", "refreshToken": "..." }

POST /auth/token/refresh
{
  "refreshToken": "dGhpcyBpcyBh..."
}
→ 200 OK { "accessToken": "...", "expiresIn": 900 }
```

**Authorization APIs:**
```http
POST /auth/authorize
{
  "principal": "usr_123",
  "action": "s3:GetObject",
  "resource": "arn:mycloud:s3:::my-bucket/private/*",
  "context": { "ip": "10.0.0.1", "time": "2026-01-15T10:00:00Z" }
}
→ 200 OK { "decision": "ALLOW", "matchedPolicies": ["policy_456"] }
```

**IAM Management APIs:**
```http
POST /iam/roles
{ "name": "developer", "description": "...", "permissions": ["read:repos", "write:repos"] }

POST /iam/policies
{
  "name": "s3-readonly",
  "effect": "ALLOW",
  "actions": ["s3:GetObject", "s3:ListBucket"],
  "resources": ["arn:mycloud:s3:::*"],
  "conditions": { "IpAddress": { "sourceIp": "10.0.0.0/8" } }
}

PUT /iam/users/{userId}/roles
{ "roles": ["developer", "viewer"] }

POST /iam/api-keys
{ "name": "CI/CD Pipeline", "permissions": ["deploy:*"], "expiresAt": "2027-01-01" }
```
"

### 4.3 Data Flow

**Candidate:** "Two critical flows:

**Flow 1: User Login**
1. Client sends email/password to Auth Service
2. Auth Service looks up user in User Directory (PostgreSQL)
3. Verify password hash (bcrypt/argon2)
4. If MFA enabled → return temporary MFA token, await TOTP code
5. On MFA success → Token Service generates JWT access token + refresh token
6. Store refresh token in Redis (server-side, for revocation)
7. Return tokens to client
8. Audit Service logs: LOGIN_SUCCESS event

**Flow 2: Authorization Check (on every API call)**
1. API Gateway receives request with `Authorization: Bearer <JWT>`
2. Gateway validates JWT signature (using cached JWKS public key — no network call)
3. Extract claims: userId, roles, orgId, permissions
4. Policy Engine evaluates: can this user perform this action on this resource?
5. Policy evaluation uses cached policies (Redis) — <10ms
6. If ALLOW → forward to resource server
7. If DENY → return 403 Forbidden
8. Audit Service logs: ACCESS_DECISION event"

---

## 5. Deep Dive: Core Components

### 5.1 Auth Service

**Candidate:** "Handles all authentication flows:

**Password Authentication:**
```javascript
async function login(email, password) {
    const user = await userDb.findByEmail(email);
    if (!user) {
        // Timing-safe: don't reveal if email exists
        await bcrypt.hash(password, 10); // dummy hash to prevent timing attack
        throw new AuthError('Invalid credentials');
    }

    // Check account lockout
    if (user.lockoutUntil && user.lockoutUntil > Date.now()) {
        throw new AuthError('Account locked. Try again later.');
    }

    const valid = await bcrypt.compare(password, user.passwordHash);
    if (!valid) {
        await incrementFailedAttempts(user.id);
        throw new AuthError('Invalid credentials');
    }

    // Reset failed attempts on success
    await resetFailedAttempts(user.id);

    // Check MFA
    if (user.mfaEnabled) {
        const mfaToken = generateMfaToken(user.id);
        return { mfaRequired: true, mfaToken };
    }

    return tokenService.issueTokens(user);
}

async function incrementFailedAttempts(userId) {
    const attempts = await redis.incr(`failed_attempts:${userId}`);
    await redis.expire(`failed_attempts:${userId}`, 900); // 15 min window

    if (attempts >= 5) {
        await userDb.setLockout(userId, Date.now() + 15 * 60 * 1000);
        await auditService.log('ACCOUNT_LOCKED', { userId });
    }
}
```

**Social Login (OAuth 2.0 Client):**
```javascript
async function handleGoogleCallback(code) {
    // Exchange authorization code for tokens
    const tokens = await googleOAuth.exchangeCode(code);
    const googleUser = await googleOAuth.getUserInfo(tokens.accessToken);

    // Find or create local user
    let user = await userDb.findByExternalId('google', googleUser.id);
    if (!user) {
        user = await userDb.create({
            email: googleUser.email,
            displayName: googleUser.name,
            externalProvider: 'google',
            externalId: googleUser.id,
            emailVerified: true
        });
    }

    return tokenService.issueTokens(user);
}
```
"

### 5.2 Token Service

**Candidate:** "Manages JWT lifecycle:

**JWT Structure:**
```
Header: { "alg": "RS256", "kid": "key-2026-01" }
Payload: {
  "sub": "usr_123",
  "iss": "https://auth.mycloud.com",
  "aud": "https://api.mycloud.com",
  "iat": 1704067200,
  "exp": 1704068100,    // 15 minutes
  "orgId": "org_456",
  "roles": ["developer", "viewer"],
  "permissions": ["read:repos", "write:repos", "read:deployments"],
  "scope": "openid profile email"
}
Signature: RS256(header + payload, privateKey)
```

**Token Issuance:**
```javascript
const jwt = require('jsonwebtoken');
const { v4: uuidv4 } = require('uuid');

async function issueTokens(user) {
    const roles = await getRolesForUser(user.id);
    const permissions = await getPermissionsForRoles(roles);

    const accessToken = jwt.sign({
        sub: user.id,
        orgId: user.orgId,
        roles: roles.map(r => r.name),
        permissions
    }, privateKey, {
        algorithm: 'RS256',
        expiresIn: '15m',
        issuer: 'https://auth.mycloud.com',
        audience: 'https://api.mycloud.com',
        keyid: currentKeyId
    });

    const refreshToken = uuidv4();
    await redis.setex(
        `refresh:${refreshToken}`,
        7 * 24 * 3600,  // 7 days
        JSON.stringify({ userId: user.id, orgId: user.orgId })
    );

    return { accessToken, refreshToken, expiresIn: 900 };
}
```

**JWKS Endpoint (for token validation by API Gateway):**
```http
GET /.well-known/jwks.json
{
  "keys": [
    {
      "kty": "RSA",
      "kid": "key-2026-01",
      "use": "sig",
      "alg": "RS256",
      "n": "0vx7agoebGcQ...",
      "e": "AQAB"
    },
    {
      "kty": "RSA",
      "kid": "key-2025-12",   // previous key, still valid
      "use": "sig",
      "alg": "RS256",
      "n": "3p8ljsdf...",
      "e": "AQAB"
    }
  ]
}
```
"

### 5.3 Policy Engine

**Candidate:** "The policy engine evaluates authorization decisions:

```javascript
class PolicyEngine {
    constructor(policyStore) {
        this.policyStore = policyStore;
        this.cache = new Map();
    }

    async evaluate(principal, action, resource, context) {
        // Get all applicable policies
        const policies = await this.getApplicablePolicies(principal, action, resource);

        let decision = 'DEFAULT_DENY';

        for (const policy of policies) {
            const match = this.matchPolicy(policy, action, resource, context);

            if (match && policy.effect === 'DENY') {
                // Explicit DENY always wins
                return { decision: 'DENY', matchedPolicy: policy.id };
            }

            if (match && policy.effect === 'ALLOW') {
                decision = 'ALLOW';
            }
        }

        return { decision, matchedPolicy: null };
    }

    matchPolicy(policy, action, resource, context) {
        // Match action (supports wildcards)
        if (!this.matchPattern(policy.actions, action)) return false;

        // Match resource (supports wildcards, e.g., arn:mycloud:s3:::my-bucket/*)
        if (!this.matchPattern(policy.resources, resource)) return false;

        // Match conditions (IP range, time, etc.)
        if (policy.conditions && !this.evaluateConditions(policy.conditions, context)) {
            return false;
        }

        return true;
    }

    matchPattern(patterns, value) {
        return patterns.some(pattern => {
            const regex = pattern.replace(/\*/g, '.*');
            return new RegExp(`^${regex}$`).test(value);
        });
    }

    evaluateConditions(conditions, context) {
        for (const [condType, condValue] of Object.entries(conditions)) {
            switch (condType) {
                case 'IpAddress':
                    if (!isIpInRange(context.ip, condValue.sourceIp)) return false;
                    break;
                case 'DateGreaterThan':
                    if (new Date(context.time) <= new Date(condValue.currentTime)) return false;
                    break;
                case 'StringEquals':
                    for (const [key, expected] of Object.entries(condValue)) {
                        if (context[key] !== expected) return false;
                    }
                    break;
            }
        }
        return true;
    }
}
```
"

### 5.4 Key Management Service

**Candidate:** "Manages cryptographic keys for JWT signing:

- Store private keys in HSM (Hardware Security Module) or HashiCorp Vault
- Key rotation: generate new key pair every 90 days
- During rotation: sign new tokens with new key, validate old tokens with old key
- JWKS endpoint serves both current and previous public keys
- Automatic retirement of old keys after 2× token lifetime"

---

## 6. Deep Dive: Database Design

### 6.1 User Directory (PostgreSQL)

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),          -- NULL for social-login-only users
    display_name VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active', -- active, suspended, deleted
    mfa_enabled BOOLEAN DEFAULT FALSE,
    mfa_secret VARCHAR(64),              -- encrypted TOTP secret
    email_verified BOOLEAN DEFAULT FALSE,
    failed_login_attempts INT DEFAULT 0,
    lockout_until TIMESTAMP,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_external_identities (
    user_id UUID REFERENCES users(id),
    provider VARCHAR(50) NOT NULL,       -- 'google', 'github', 'saml_okta'
    external_id VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    metadata JSONB,
    PRIMARY KEY (provider, external_id)
);

CREATE INDEX idx_users_org ON users(org_id);
CREATE INDEX idx_users_email ON users(email);
```

### 6.2 Organizations (Multi-Tenant)

```sql
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    plan VARCHAR(50) DEFAULT 'free',     -- free, pro, enterprise
    sso_enabled BOOLEAN DEFAULT FALSE,
    sso_config JSONB,                    -- SAML/OIDC provider config
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 6.3 Roles & Permissions (RBAC)

```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT FALSE,    -- system roles can't be modified
    parent_role_id UUID REFERENCES roles(id),  -- role hierarchy
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(org_id, name)
);

CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL UNIQUE,   -- e.g., 'repos:read', 'deploy:create'
    description TEXT,
    resource_type VARCHAR(100)           -- 'repository', 'deployment', etc.
);

CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id),
    permission_id UUID REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_role_assignments (
    user_id UUID REFERENCES users(id),
    role_id UUID REFERENCES roles(id),
    org_id UUID REFERENCES organizations(id),
    scope VARCHAR(255),                  -- optional resource scope: 'project:proj_123'
    assigned_by UUID REFERENCES users(id),
    assigned_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,               -- temporary role assignment
    PRIMARY KEY (user_id, role_id, org_id)
);

CREATE INDEX idx_user_roles ON user_role_assignments(user_id, org_id);
```

### 6.4 Policies (ABAC)

```sql
CREATE TABLE policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    effect VARCHAR(10) NOT NULL CHECK (effect IN ('ALLOW', 'DENY')),
    actions TEXT[] NOT NULL,             -- ['s3:GetObject', 's3:PutObject']
    resources TEXT[] NOT NULL,           -- ['arn:mycloud:s3:::my-bucket/*']
    conditions JSONB,                    -- {"IpAddress": {"sourceIp": "10.0.0.0/8"}}
    priority INT DEFAULT 0,             -- higher = evaluated first
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE policy_attachments (
    policy_id UUID REFERENCES policies(id),
    principal_type VARCHAR(20) NOT NULL, -- 'user', 'role', 'group'
    principal_id UUID NOT NULL,
    PRIMARY KEY (policy_id, principal_type, principal_id)
);
```

### 6.5 API Keys

```sql
CREATE TABLE api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    org_id UUID REFERENCES organizations(id),
    key_prefix VARCHAR(10) NOT NULL,     -- first 10 chars (for identification)
    key_hash VARCHAR(255) NOT NULL,      -- bcrypt hash of full key
    name VARCHAR(200) NOT NULL,
    permissions TEXT[],
    last_used TIMESTAMP,
    expires_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT NOW()
);
```

### 6.6 Audit Log

```sql
-- This is the schema; actual storage is in Kafka → S3/Elasticsearch
CREATE TABLE audit_log (
    id BIGSERIAL,
    timestamp TIMESTAMP NOT NULL,
    org_id UUID,
    actor_id UUID,                       -- who performed the action
    actor_type VARCHAR(20),              -- 'user', 'service', 'api_key'
    action VARCHAR(100) NOT NULL,        -- 'LOGIN', 'AUTHORIZE', 'ROLE_ASSIGN'
    resource VARCHAR(255),
    decision VARCHAR(10),                -- 'ALLOW', 'DENY' (for auth events)
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,                      -- additional context
    PRIMARY KEY (timestamp, id)
) PARTITION BY RANGE (timestamp);
```

---

## 7. Deep Dive: Authentication & Authorization Engine

### 7.1 Authentication Flows

**Candidate:** "Let me cover each auth flow:

**A. Username/Password with MFA:**
```
Client                    Auth Service              User DB          Redis
  │                           │                        │              │
  │── POST /auth/login ──────▶│                        │              │
  │   {email, password}       │── Find user ──────────▶│              │
  │                           │◀── user record ────────│              │
  │                           │                        │              │
  │                           │── bcrypt.compare ─────▶│              │
  │                           │◀── match ──────────────│              │
  │                           │                        │              │
  │                           │── MFA enabled? Yes     │              │
  │                           │── Generate mfaToken ──────────────────▶│
  │◀── {mfaRequired, token} ──│                        │              │
  │                           │                        │              │
  │── POST /auth/mfa/verify ─▶│                        │              │
  │   {mfaToken, code}        │── Verify TOTP ────────▶│              │
  │                           │── Issue JWT ───────────────────────────▶│
  │◀── {accessToken, refresh} │                        │              │
```

**B. OAuth 2.0 Authorization Code Flow (with PKCE):**
```
Client              Auth Server              Resource Server
  │                      │                         │
  │ 1. Generate code_verifier, code_challenge      │
  │                      │                         │
  │ 2. GET /oauth/authorize?                       │
  │    response_type=code&                         │
  │    client_id=app_123&                          │
  │    redirect_uri=https://app.com/callback&      │
  │    code_challenge=abc123&                      │
  │    code_challenge_method=S256&                 │
  │    scope=openid+profile                        │
  │ ────────────────────▶│                         │
  │                      │                         │
  │ 3. User authenticates, grants consent          │
  │                      │                         │
  │ 4. Redirect to:      │                         │
  │    https://app.com/callback?code=AUTH_CODE      │
  │ ◀────────────────────│                         │
  │                      │                         │
  │ 5. POST /oauth/token │                         │
  │    grant_type=authorization_code&              │
  │    code=AUTH_CODE&                             │
  │    code_verifier=original_verifier             │
  │ ────────────────────▶│                         │
  │                      │ Verify:                 │
  │                      │ SHA256(verifier)==challenge│
  │                      │                         │
  │ 6. {access_token, id_token, refresh_token}     │
  │ ◀────────────────────│                         │
  │                      │                         │
  │ 7. GET /api/resource │                         │
  │    Authorization: Bearer access_token          │
  │ ──────────────────────────────────────────────▶│
```

**C. TOTP Multi-Factor Authentication:**
```javascript
const speakeasy = require('speakeasy');
const qrcode = require('qrcode');

// Setup MFA
async function enableMFA(userId) {
    const secret = speakeasy.generateSecret({
        name: `MyCloud (${user.email})`,
        issuer: 'MyCloud',
        length: 20
    });

    // Store encrypted secret
    await userDb.update(userId, {
        mfaSecret: encrypt(secret.base32),
        mfaPending: true
    });

    // Generate QR code for authenticator app
    const qrCodeUrl = await qrcode.toDataURL(secret.otpauth_url);
    return { qrCode: qrCodeUrl, backupCodes: generateBackupCodes() };
}

// Verify TOTP code
function verifyTOTP(secret, code) {
    return speakeasy.totp.verify({
        secret: secret,
        encoding: 'base32',
        token: code,
        window: 1  // allow 1 step tolerance (30 sec before/after)
    });
}
```
"

### 7.2 Authorization Models

**Candidate:** "We support two complementary models:

**RBAC (Role-Based Access Control):**
```
User → Role(s) → Permission(s)

Example:
  User: alice
  Roles: [developer, viewer]
  Permissions: [repos:read, repos:write, deployments:read]

Role Hierarchy:
  admin → developer → viewer
  (admin inherits all developer permissions,
   developer inherits all viewer permissions)
```

**RBAC Resolution:**
```javascript
async function getEffectivePermissions(userId, orgId) {
    const roles = await getUserRoles(userId, orgId);
    const allRoles = new Set();

    // Resolve role hierarchy
    for (const role of roles) {
        allRoles.add(role);
        const ancestors = await getRoleAncestors(role);
        ancestors.forEach(a => allRoles.add(a));
    }

    // Collect permissions from all roles
    const permissions = new Set();
    for (const role of allRoles) {
        const rolePerms = await getRolePermissions(role);
        rolePerms.forEach(p => permissions.add(p));
    }

    return Array.from(permissions);
}
```

**ABAC (Attribute-Based Access Control):**
```json
{
  "name": "Allow S3 read from office IP only",
  "effect": "ALLOW",
  "actions": ["s3:GetObject", "s3:ListBucket"],
  "resources": ["arn:mycloud:s3:::confidential-bucket/*"],
  "conditions": {
    "IpAddress": { "sourceIp": "10.0.0.0/8" },
    "DateGreaterThan": { "currentTime": "2026-01-01T00:00:00Z" },
    "StringEquals": { "principalOrgId": "org_456" }
  }
}
```

**Combined Evaluation (RBAC + ABAC):**
```javascript
async function authorize(principal, action, resource, context) {
    // Step 1: Check RBAC permissions (fast path)
    const permissions = await getEffectivePermissions(principal.userId, principal.orgId);
    const actionPermission = mapActionToPermission(action);

    if (!permissions.includes(actionPermission) && !permissions.includes('*')) {
        return { decision: 'DENY', reason: 'No RBAC permission' };
    }

    // Step 2: Evaluate ABAC policies (fine-grained)
    const policies = await policyEngine.getApplicablePolicies(principal, action, resource);
    const abacDecision = await policyEngine.evaluate(policies, action, resource, context);

    // Explicit DENY in ABAC overrides RBAC ALLOW
    if (abacDecision.decision === 'DENY') {
        return { decision: 'DENY', reason: abacDecision.matchedPolicy };
    }

    return { decision: 'ALLOW' };
}
```
"

### 7.3 Token Management

**Candidate:** "Token lifecycle management:

**Access Token (Short-Lived, 15 min):**
- JWT format, self-contained claims
- Validated locally at API Gateway (no auth service round-trip)
- Contains: userId, orgId, roles, permissions
- Signed with RS256 (asymmetric key)

**Refresh Token (Long-Lived, 7 days):**
- Opaque string (UUID), stored in Redis
- Used to get new access tokens without re-login
- Rotation: each refresh issues a new refresh token (old one invalidated)
- Stored server-side for revocation capability

**Token Refresh:**
```javascript
async function refreshAccessToken(refreshToken) {
    const storedData = await redis.get(`refresh:${refreshToken}`);
    if (!storedData) throw new AuthError('Invalid refresh token');

    const { userId, orgId } = JSON.parse(storedData);
    const user = await userDb.findById(userId);

    if (user.status !== 'active') throw new AuthError('Account suspended');

    // Rotate refresh token
    await redis.del(`refresh:${refreshToken}`);
    const newRefreshToken = uuidv4();
    await redis.setex(`refresh:${newRefreshToken}`, 7 * 24 * 3600,
        JSON.stringify({ userId, orgId }));

    // Issue new access token
    const accessToken = await issueAccessToken(user);

    return { accessToken, refreshToken: newRefreshToken };
}
```

**Token Revocation:**
```javascript
async function revokeAllTokens(userId) {
    // Revoke all refresh tokens
    const pattern = `refresh:*`;
    // In practice, maintain a set of refresh tokens per user
    const userRefreshTokens = await redis.smembers(`user_refresh:${userId}`);
    for (const token of userRefreshTokens) {
        await redis.del(`refresh:${token}`);
    }

    // For access tokens: add to blacklist (checked on validation)
    await redis.setex(`revoked_user:${userId}`, 900, '1'); // 15 min (token lifetime)
}
```

**Key Rotation:**
```javascript
async function rotateSigningKey() {
    const newKeyPair = await generateRSAKeyPair();
    const newKeyId = `key-${Date.now()}`;

    // Store new private key in HSM/Vault
    await vault.store(newKeyId, newKeyPair.privateKey);

    // Update JWKS endpoint (serve both old and new public keys)
    const jwks = await getJWKS();
    jwks.keys.unshift({
        kid: newKeyId,
        ...publicKeyToJWK(newKeyPair.publicKey)
    });

    // Remove keys older than 2× max token lifetime
    jwks.keys = jwks.keys.filter(k => !isExpired(k, 30 * 60)); // 30 min

    await updateJWKS(jwks);

    // New tokens signed with new key
    currentKeyId = newKeyId;
}
```
"

### 7.4 Session Management

**Candidate:** "Two approaches:

**Stateless (JWT-Based) — Recommended for API access:**
- No server-side session state
- Token contains all necessary claims
- Validated locally using public key
- Trade-off: can't instantly revoke (wait for expiry or use blacklist)

**Stateful (Server-Side Sessions) — For web dashboard:**
```javascript
async function createSession(userId, deviceInfo) {
    const sessionId = generateSecureRandom(32);
    const session = {
        userId,
        orgId: user.orgId,
        device: deviceInfo,
        createdAt: Date.now(),
        lastActive: Date.now(),
        ip: deviceInfo.ip
    };

    await redis.setex(`session:${sessionId}`, 24 * 3600, JSON.stringify(session));
    await redis.sadd(`user_sessions:${userId}`, sessionId);

    return sessionId;
}
```

**Concurrent Session Control:**
```javascript
async function enforceSessionLimit(userId, maxSessions = 5) {
    const sessions = await redis.smembers(`user_sessions:${userId}`);
    if (sessions.length >= maxSessions) {
        // Revoke oldest session
        const oldest = await findOldestSession(sessions);
        await redis.del(`session:${oldest}`);
        await redis.srem(`user_sessions:${userId}`, oldest);
    }
}
```
"

### 7.5 Security Hardening

**Candidate:** "Critical security measures:

**Brute Force Protection:**
- Progressive delay: 1st fail = 0s, 2nd = 1s, 3rd = 2s, 4th = 4s, 5th = lock 15 min
- Per-IP rate limiting: 20 login attempts per IP per 15 minutes
- Per-account rate limiting: 5 failed attempts = lock 15 minutes

**Credential Stuffing Defense:**
```javascript
const rateLimit = require('express-rate-limit');

const loginLimiter = rateLimit({
    windowMs: 15 * 60 * 1000,  // 15 minutes
    max: 20,                    // 20 attempts per IP
    keyGenerator: (req) => req.ip,
    handler: (req, res) => {
        res.status(429).json({ error: 'Too many login attempts' });
    }
});

app.post('/auth/login', loginLimiter, handleLogin);
```

**Password Security:**
- Argon2id hashing (preferred) or bcrypt with cost factor 12
- Minimum 8 characters, check against breached password database (HaveIBeenPwned API)
- No password hints or security questions

**Token Security:**
- HttpOnly, Secure, SameSite=Strict cookies for web
- Short-lived access tokens (15 min)
- Refresh token rotation (one-time use)
- Token binding to device fingerprint (optional, high-security)
"

---

## 8. Scaling Strategies

### 8.1 Current Bottlenecks

**Candidate:** "At our scale:

1. **Authorization (1M checks/sec):** This is the hot path. Must be sub-10ms.
   - Solution: Policy evaluation cached at API Gateway level
   - Pre-compute user permissions on login, embed in JWT
   
2. **Authentication (115 logins/sec):** Not a bottleneck.
   
3. **Audit Logging (1M events/sec):** Write-heavy.
   - Solution: Kafka buffering → async write to S3/Elasticsearch

4. **User Directory:** 100M users, 115 login lookups/sec — PostgreSQL handles easily."

### 8.2 Making Authorization Sub-10ms

**Candidate:** "The key scaling challenge is 1M authorization checks/sec at <10ms:

**Strategy 1: Embed permissions in JWT**
- At login, compute all permissions and put in JWT
- API Gateway reads permissions from token — zero network calls
- Trade-off: JWT gets large for users with many permissions
- Solution: use permission sets/groups instead of individual permissions

**Strategy 2: Edge-deployed policy engine**
```
API Gateway has Policy Engine SDK embedded:
  1. Cache JWKS public key (refresh every 5 min)
  2. Cache policies per role (refresh every 1 min)
  3. Validate JWT locally (no network call)
  4. Evaluate policies locally (no network call)
  5. Total auth overhead: <5ms
```

**Strategy 3: Policy materialization**
```javascript
// Pre-compute: for each role, what resources can it access?
// Store in Redis as hash map
await redis.hset('role_permissions:developer', {
    'repos:read': '1',
    'repos:write': '1',
    'deployments:read': '1'
});

// Authorization check: O(1) lookup
const allowed = await redis.hget(`role_permissions:${role}`, action);
```
"

### 8.3 Scaling to 10x (10M auth checks/sec)

**Candidate:** "
1. **Horizontally scale API Gateways** — each does local token validation
2. **Redis cluster** for session/policy cache — shard by org_id
3. **PostgreSQL read replicas** for User Directory
4. **Kafka scaling** for audit log throughput
5. **Regional deployment** — auth infrastructure in each region"

---

## 9. Failure Scenarios & Mitigation

### 9.1 Auth Service Failure

**Scenario:** Auth Service (login/register) goes down.

**Impact:**
- New logins fail
- Existing sessions still work (JWT validated locally)

**Mitigation:**
- Multiple Auth Service replicas behind load balancer
- Health checks with automatic failover
- Cached tokens continue working until expiry
- Refresh tokens still served by Redis (separate from Auth Service)
- **Impact on existing users: NONE** (JWT validation is local)
- **Impact on new logins: Degraded** until Auth Service recovers

### 9.2 Key Rotation Failure

**Scenario:** New signing key deployed but JWKS endpoint not updated.

**Impact:**
- New tokens signed with new key can't be validated
- All new logins produce unusable tokens

**Mitigation:**
- Always maintain at least 2 keys in JWKS
- Validate that JWKS endpoint returns new key BEFORE issuing tokens with it
- Rollback: revert to old key if JWKS update fails
- Monitor: alert if JWKS endpoint returns fewer than 2 keys

### 9.3 MFA Service Timeout

**Scenario:** TOTP validation service is slow or down.

**Impact:**
- Users with MFA can't complete login

**Mitigation:**
- TOTP validation is a local computation (no external service needed)
- If using SMS-based MFA: fall back to email-based OTP
- Allow backup codes (pre-generated) for emergency access
- Adaptive: skip MFA for trusted devices/IPs (risk-based)

### 9.4 Policy Cache Inconsistency

**Scenario:** Policy updated but cache not refreshed; stale permissions used.

**Impact:**
- User may have access after permission removal (brief window)

**Mitigation:**
- Cache TTL: 60 seconds maximum for policy cache
- On critical policy changes (permission removal): push invalidation via Kafka
- Emergency: endpoint to force-flush cache across all gateways
```javascript
async function updatePolicy(policyId, newPolicy) {
    await policyDb.update(policyId, newPolicy);
    // Push invalidation to all API Gateways
    await kafka.publish('policy-invalidation', {
        policyId, action: 'update', timestamp: Date.now()
    });
}
```

### 9.5 Redis Failure (Session Store)

**Scenario:** Redis cluster goes down.

**Impact:**
- Refresh tokens can't be validated → users can't refresh sessions
- But existing JWT access tokens still work (validated locally)

**Mitigation:**
- Redis Cluster with failover replicas
- Redis Sentinel for automatic failover
- During outage: extend JWT acceptance window (accept slightly expired tokens)
- Users may need to re-login once Redis recovers

### 9.6 Audit Log Pipeline Failure

**Scenario:** Kafka or S3 write fails; audit events lost.

**Impact:**
- Compliance violation (missing audit trail)

**Mitigation:**
- Kafka replication factor 3 (survive 2 broker failures)
- Local buffer on API Gateway: if Kafka is down, buffer events locally
- Replay from gateway logs when Kafka recovers
- Alert immediately: compliance-critical

---

## 10. Monitoring & Observability

### 10.1 Key Metrics

**Candidate:** "For an IAM system:

**Security Metrics:**
- Login success/failure rate
- Failed login attempts by IP (detect brute force)
- MFA adoption rate
- Suspicious activity score
- Privilege escalation attempts

**Performance Metrics (RED):**
1. **Rate:** Logins/sec, token refreshes/sec, auth checks/sec
2. **Errors:** Auth failures, token validation errors
3. **Duration:** Login latency, authorization check latency

**Business Metrics:**
- Active users (DAU/MAU)
- API key usage
- Role distribution (how many admins vs developers)

**Example Dashboard (Grafana):**
```
Row 1: Authentication
- [Graph] Login success vs failure rate
- [Heatmap] Login latency distribution
- [Graph] MFA verification rate

Row 2: Authorization
- [Graph] Authorization checks/sec (ALLOW vs DENY)
- [Heatmap] Authorization latency (must be <10ms)
- [Graph] Policy evaluation cache hit rate

Row 3: Security
- [Graph] Failed login attempts by IP (top 10)
- [Graph] Account lockout events
- [Alert Panel] Privilege escalation attempts

Row 4: Infrastructure
- [Graph] Redis memory usage and ops/sec
- [Graph] JWT signing key age
- [Graph] Audit log pipeline lag
```
"

### 10.2 Alerting Rules

```yaml
alert: HighLoginFailureRate
expr: rate(login_failures_total[5m]) / rate(login_attempts_total[5m]) > 0.3
for: 5m
severity: critical
message: "Login failure rate above 30% - possible credential stuffing attack"

alert: AuthorizationLatencyHigh
expr: histogram_quantile(0.99, auth_check_latency_seconds) > 0.01
for: 2m
severity: critical
message: "Authorization p99 latency above 10ms - SLA violation"

alert: SuspiciousBruteForce
expr: rate(login_failures_total[5m]) > 100
for: 2m
severity: critical
message: "More than 100 failed logins per second - potential attack"

alert: JWTSigningKeyExpiringSoon
expr: jwt_signing_key_age_days > 80
for: 1h
severity: warning
message: "JWT signing key is 80+ days old - rotation due"

alert: AuditLogLag
expr: kafka_consumer_lag{topic="audit-events"} > 1000000
for: 5m
severity: critical
message: "Audit log pipeline lag over 1M events - compliance risk"
```

### 10.3 Audit Logging

**Candidate:** "Every security-relevant event must be logged:

```javascript
const auditEvents = {
    LOGIN_SUCCESS: (user, context) => ({
        action: 'LOGIN_SUCCESS',
        actorId: user.id,
        ip: context.ip,
        userAgent: context.userAgent,
        mfaUsed: context.mfaUsed,
        method: context.authMethod  // 'password', 'google', 'saml'
    }),

    LOGIN_FAILURE: (email, context) => ({
        action: 'LOGIN_FAILURE',
        email: hashEmail(email),  // don't log raw email
        ip: context.ip,
        reason: context.reason    // 'invalid_password', 'account_locked'
    }),

    AUTHORIZE_DENY: (principal, action, resource) => ({
        action: 'AUTHORIZE_DENY',
        actorId: principal.userId,
        attemptedAction: action,
        resource: resource,
        reason: 'insufficient_permissions'
    }),

    ROLE_ASSIGNED: (admin, targetUser, role) => ({
        action: 'ROLE_ASSIGNED',
        actorId: admin.id,
        targetUserId: targetUser.id,
        roleName: role.name
    }),

    POLICY_UPDATED: (admin, policy) => ({
        action: 'POLICY_UPDATED',
        actorId: admin.id,
        policyId: policy.id,
        changes: policy.diff
    })
};
```

**Storage:** Kafka → Elasticsearch (searchable, 30 days) + S3 (archive, 7 years for compliance)"

---

## 11. Advanced Features

### 11.1 Passwordless Authentication

**Candidate:** "Modern alternatives to passwords:

**Magic Links (Email):**
```javascript
async function sendMagicLink(email) {
    const token = generateSecureRandom(32);
    await redis.setex(`magic:${token}`, 600, email); // 10 min TTL

    await emailService.send(email, {
        subject: 'Sign in to MyCloud',
        body: `Click to sign in: https://auth.mycloud.com/magic/${token}`
    });
}

async function verifyMagicLink(token) {
    const email = await redis.get(`magic:${token}`);
    if (!email) throw new AuthError('Invalid or expired link');

    await redis.del(`magic:${token}`); // one-time use
    const user = await userDb.findByEmail(email);
    return tokenService.issueTokens(user);
}
```

**WebAuthn / FIDO2 (Biometric/Hardware Key):**
- Browser API for hardware-backed authentication
- Supports fingerprint, face ID, YubiKey
- Phishing-resistant (bound to origin domain)
- Strongest authentication factor available"

### 11.2 Adaptive / Risk-Based Authentication

**Candidate:** "Dynamically adjust auth requirements based on risk:

```javascript
async function calculateRiskScore(loginAttempt) {
    let riskScore = 0;

    // Known device?
    if (!isKnownDevice(loginAttempt.deviceFingerprint)) riskScore += 30;

    // Known IP?
    if (!isKnownIP(loginAttempt.ip)) riskScore += 20;

    // Unusual location?
    const location = geolocate(loginAttempt.ip);
    if (isImpossibleTravel(loginAttempt.userId, location)) riskScore += 40;

    // Unusual time?
    if (isUnusualLoginTime(loginAttempt.userId)) riskScore += 10;

    return riskScore;
}

async function adaptiveAuth(user, loginAttempt) {
    const risk = await calculateRiskScore(loginAttempt);

    if (risk < 20) {
        return { requiredFactors: ['password'] };           // low risk
    } else if (risk < 50) {
        return { requiredFactors: ['password', 'totp'] };   // medium risk
    } else {
        return { requiredFactors: ['password', 'totp', 'email_otp'] }; // high risk
    }
}
```
"

### 11.3 Service Mesh Integration (mTLS + SPIFFE)

**Candidate:** "For service-to-service authentication:

- Each service gets a SPIFFE ID: `spiffe://mycloud.com/service/payment-service`
- Mutual TLS (mTLS) between services — both sides present certificates
- SPIRE (SPIFFE Runtime Environment) manages certificate issuance and rotation
- No static API keys between services — certificates auto-rotate every hour
- IAM policies can reference SPIFFE IDs for service-level authorization"

### 11.4 Just-In-Time (JIT) Provisioning

```javascript
async function handleSAMLAssertion(samlAssertion) {
    const externalUser = parseSAMLAssertion(samlAssertion);

    let localUser = await userDb.findByExternalId('saml', externalUser.id);

    if (!localUser) {
        // JIT: create user on first login via SSO
        localUser = await userDb.create({
            email: externalUser.email,
            displayName: externalUser.name,
            orgId: resolveOrgFromSAMLIssuer(samlAssertion.issuer),
            externalProvider: 'saml',
            externalId: externalUser.id
        });

        // Auto-assign default role
        await assignRole(localUser.id, 'viewer');
    }

    // Sync group memberships from SAML attributes
    await syncGroupMemberships(localUser.id, externalUser.groups);

    return tokenService.issueTokens(localUser);
}
```

### 11.5 Consent Management (GDPR)

```javascript
async function recordConsent(userId, consentType, granted) {
    await db.query(`
        INSERT INTO user_consents (user_id, consent_type, granted, recorded_at, ip_address)
        VALUES ($1, $2, $3, NOW(), $4)
    `, [userId, consentType, granted, context.ip]);

    if (!granted) {
        // If user withdraws consent, take action
        if (consentType === 'data_processing') {
            await scheduleAccountDeletion(userId);
        }
    }
}

// GDPR data export
async function exportUserData(userId) {
    const userData = await userDb.findById(userId);
    const auditLog = await auditDb.findByUser(userId);
    const sessions = await redis.smembers(`user_sessions:${userId}`);

    return {
        profile: sanitize(userData),
        loginHistory: auditLog,
        activeSessions: sessions.length,
        exportedAt: new Date().toISOString()
    };
}
```

---

## 12. Interview Q&A

### Q1: How do you make authorization checks fast enough for every API call (<10ms)?

**Answer:**
Three-layer approach:

1. **Embed permissions in JWT:** At login time, compute the user's effective permissions (from roles + policies) and embed them in the JWT claims. The API Gateway reads permissions from the token — zero network calls.

2. **Local policy evaluation:** Deploy the Policy Engine as an SDK/sidecar at the API Gateway. Cache policies in memory (refresh every 60 seconds). Evaluation is a local in-memory operation: match action against permission list.

3. **Pre-materialized permissions:** For common checks, pre-compute and cache `role → permission set` mappings in Redis. Authorization becomes a single hash lookup.

The result: JWT validation (crypto verification ~1ms) + permission check (memory lookup ~0.1ms) = ~1-2ms total. Well under 10ms.

### Q2: JWT vs opaque tokens — trade-offs?

**Answer:**

| Aspect | JWT | Opaque Token |
|--------|-----|-------------|
| **Validation** | Local (no network call) | Requires round-trip to auth server |
| **Revocation** | Hard (need blacklist or short TTL) | Easy (delete from store) |
| **Size** | Large (500+ bytes) | Small (36 bytes UUID) |
| **Scalability** | Excellent (no central bottleneck) | Limited by auth server capacity |
| **Information** | Self-contained claims | Requires introspection endpoint |
| **Security** | Token theft = access until expiry | Token theft = access until revoked |

**Recommendation:** Use JWT for API access (scalability + no round-trips) with short TTL (15 min) and refresh token rotation. Use opaque tokens for refresh tokens (server-side, revocable).

### Q3: How do you handle token revocation with stateless JWTs?

**Answer:**
JWTs are inherently hard to revoke since they're self-contained. Strategies:

1. **Short TTL (primary approach):** Access tokens expire in 15 minutes. Even if compromised, the damage window is small.

2. **Token blacklist:** On revocation (logout, password change), add the JWT's `jti` (JWT ID) to a Redis blacklist. API Gateway checks blacklist on each request. Blacklist entries expire after the token's original TTL.

3. **Refresh token revocation:** Revoking the refresh token prevents new access tokens from being issued. Combined with short TTL, the user is effectively logged out within 15 minutes.

4. **Emergency: user-level revocation:** Add `revoked_user:{userId}` to Redis. All tokens for that user are rejected. Useful for account compromise.

### Q4: RBAC vs ABAC — when would you use each?

**Answer:**

**RBAC (Role-Based):**
- Best for: organizational structure, job-function-based access
- Example: "Developers can deploy to staging, only SREs can deploy to production"
- Pros: Simple to understand, audit, and manage
- Cons: Role explosion (too many roles), no contextual rules

**ABAC (Attribute-Based):**
- Best for: fine-grained, context-dependent access control
- Example: "Allow read access to patient records only if doctor is assigned to patient AND request comes from hospital network AND it's during work hours"
- Pros: Extremely flexible, handles complex scenarios
- Cons: Complex to manage, harder to audit, policy explosion

**In practice:** Use both. RBAC for coarse-grained access (role membership), ABAC for fine-grained conditions (IP restrictions, time-based access, resource-level permissions).

### Q5: How would you implement OAuth 2.0 Authorization Code flow with PKCE?

**Answer:**
PKCE (Proof Key for Code Exchange) prevents authorization code interception:

1. **Client generates:** `code_verifier` (random 43-128 chars) and `code_challenge` = SHA256(code_verifier), base64url-encoded

2. **Authorization request:** Client sends `code_challenge` and `code_challenge_method=S256` to `/oauth/authorize`

3. **User authenticates** and grants consent. Auth server stores `code_challenge` with the authorization code.

4. **Token exchange:** Client sends `code_verifier` (not challenge) with the authorization code to `/oauth/token`

5. **Server verifies:** SHA256(code_verifier) === stored code_challenge. If match, issue tokens.

This prevents attackers who intercept the authorization code from exchanging it, since they don't have the original `code_verifier`. PKCE is now required for all public clients (SPAs, mobile apps).

### Q6: How do you handle key rotation for JWT signing without invalidating existing tokens?

**Answer:**
Overlap period approach:

1. **Generate new key pair** with a new `kid` (key ID)
2. **Publish both keys** in JWKS endpoint (new + old)
3. **Start signing** new tokens with the new key
4. **Wait** for all old tokens to expire (15 minutes for access tokens)
5. **Remove old key** from JWKS after 2× token lifetime
6. **API Gateways** use the `kid` in the JWT header to select the correct public key for verification

Timeline example:
```
T=0:   Generate key-2026-04, add to JWKS
T=0:   Start signing with key-2026-04
T=0:   JWKS serves both key-2026-01 (old) and key-2026-04 (new)
T=30m: All old tokens have expired
T=30m: Remove key-2026-01 from JWKS
```

Automation: key rotation runs as a cron job every 90 days.

### Q7: How would you design multi-tenant IAM (each tenant has their own users/roles)?

**Answer:**
Organization-level isolation:

1. **Data isolation:** Every entity (user, role, policy) has an `org_id` foreign key. All queries filter by `org_id`.

2. **Role scoping:** Roles are per-organization. The "admin" role in Org A has no access to Org B's resources.

3. **Policy evaluation:** Policies are scoped to the organization. Cross-org access requires explicit federation policies.

4. **Database:** Single database with `org_id` column (shared tenancy) vs separate database per org (dedicated tenancy). Start with shared tenancy, offer dedicated for enterprise customers.

5. **Token claims:** JWT includes `orgId` — API Gateway enforces that users can only access resources in their org.

```javascript
// Middleware: enforce org isolation
function orgIsolation(req, res, next) {
    const userOrgId = req.auth.orgId;
    const resourceOrgId = getResourceOrgId(req);

    if (userOrgId !== resourceOrgId) {
        return res.status(403).json({ error: 'Cross-organization access denied' });
    }
    next();
}
```

### Q8: How do you prevent privilege escalation attacks?

**Answer:**
Multiple layers of defense:

1. **Role assignment validation:** Only users with `admin` role can assign roles. Admins cannot assign roles higher than their own level.

```javascript
async function assignRole(adminUserId, targetUserId, roleId) {
    const adminRoles = await getUserRoles(adminUserId);
    const targetRole = await getRole(roleId);

    // Admin can't assign higher-privilege role
    if (targetRole.level > getMaxLevel(adminRoles)) {
        throw new AuthError('Cannot assign role with higher privilege');
    }
}
```

2. **Policy validation:** When creating policies, verify the creator has all permissions they're granting. You can't create a policy granting permissions you don't have.

3. **Audit every change:** All role/policy modifications are logged with before/after state. Regular audit reviews catch unauthorized changes.

4. **Separation of duties:** Require two admins to approve privilege changes (4-eyes principle for critical operations).

5. **Time-limited elevated access:** Instead of permanent admin roles, grant temporary elevated access (e.g., 4-hour admin session) with automatic expiration.

### Q9: How would you implement zero-trust architecture with this IAM system?

**Answer:**
Zero-trust principles applied to our IAM:

1. **Never trust, always verify:** Every API call validates JWT and checks permissions, even for internal service-to-service calls.

2. **Least privilege:** Default deny. Users get minimum permissions for their role. Temporary escalation with approval workflow.

3. **Assume breach:** Short-lived tokens (15 min), continuous validation, anomaly detection on access patterns.

4. **Micro-segmentation:** Each microservice has its own SPIFFE identity. Service mesh enforces mTLS between services. IAM policies define which services can communicate.

5. **Device trust:** Token binding to device fingerprint. Unknown devices trigger step-up authentication (MFA).

6. **Continuous authorization:** Re-evaluate permissions periodically during long sessions, not just at login. If user's risk score changes (new IP, unusual behavior), force re-authentication.

---

## 13. Production Checklist

### 13.1 Pre-Launch

- [ ] **Security Audit:** Third-party penetration testing
- [ ] **Compliance:** SOC2 Type II readiness, GDPR data handling
- [ ] **Password Policy:** Enforce minimum strength, check against breached databases
- [ ] **Key Management:** HSM or Vault setup, key rotation automated
- [ ] **Load Testing:** Simulate 1M auth checks/sec
- [ ] **Failover Testing:** Kill auth servers, verify JWT validation continues
- [ ] **Rate Limiting:** Login attempt limits, API key rate limits
- [ ] **Audit Logging:** Verify all events captured, retention policy set
- [ ] **Backup:** User directory backup, key backup in separate HSM

### 13.2 Day-1 Operations

- [ ] Monitor login success/failure rate
- [ ] Monitor authorization latency (must be <10ms p99)
- [ ] Check audit log pipeline health
- [ ] Verify no 500 errors on auth endpoints
- [ ] Monitor Redis session store health

### 13.3 Week-1 Optimization

- [ ] Analyze login patterns (peak hours, failure reasons)
- [ ] Tune rate limiting thresholds
- [ ] Review and tune policy cache TTL
- [ ] Check MFA adoption rate
- [ ] Review audit logs for suspicious patterns

### 13.4 Month-1 Scaling

- [ ] Review capacity (plan for 3x growth)
- [ ] Implement adaptive authentication if not in MVP
- [ ] Add social login providers based on user demand
- [ ] Implement SSO/SAML for enterprise customers
- [ ] Cost optimization (reserved instances, cache sizing)

---

## Summary: Key Takeaways

### Technical Decisions

| Component | Choice | Rationale |
|-----------|--------|-----------|
| **Access Token** | JWT (RS256) | Local validation, no round-trips, scalable |
| **Refresh Token** | Opaque (Redis) | Server-side, revocable, secure |
| **User Directory** | PostgreSQL | ACID, relational data, mature |
| **Session/Cache** | Redis Cluster | Sub-ms lookups, TTL, Pub/Sub for invalidation |
| **Audit Logs** | Kafka → S3/Elasticsearch | High-throughput writes, searchable, archivable |
| **Key Management** | HSM / HashiCorp Vault | Hardware-backed security, key rotation |
| **Authorization** | RBAC + ABAC | Coarse-grained roles + fine-grained policies |
| **Password Hashing** | Argon2id / bcrypt | Memory-hard, resistant to GPU attacks |

### Scalability Path

1. **Current (1M auth checks/sec):** JWT local validation, Redis policy cache, PostgreSQL
2. **10x (10M checks/sec):** Edge-deployed policy engine, sharded Redis, read replicas
3. **100x (100M checks/sec):** Multi-region IAM, federated auth, per-tenant isolation

### Interview Performance Tips

1. ✅ Start with the distinction: authentication (who are you?) vs authorization (what can you do?)
2. ✅ Explain JWT structure and why local validation matters for scale
3. ✅ Discuss RBAC vs ABAC with concrete examples
4. ✅ Address token revocation (the classic JWT dilemma)
5. ✅ Deep dive into OAuth 2.0 + PKCE flow
6. ✅ Mention security hardening (brute force, credential stuffing)
7. ✅ Discuss audit logging as a compliance requirement

---

**End of Identity & Access Management (IAM) System Design**  
[← Back to Main Index](../README.md)
