# AWS Enterprise Architecture & Operations

> **Target Audience**: Senior Software Engineers transitioning to Platform Engineer / Architect roles
> **Prerequisites**: Completed Part 1 & Part 2, hands-on AWS experience with multiple services

---

## Table of Contents

1. [Multi-Account Strategy](#1-multi-account-strategy)
2. [Disaster Recovery & Business Continuity](#2-disaster-recovery--business-continuity)
3. [Security Deep Dive](#3-security-deep-dive)
4. [Cost Optimization at Scale](#4-cost-optimization-at-scale)
5. [DevOps & CI/CD on AWS](#5-devops--cicd-on-aws)
6. [Real-World Architecture Case Studies](#6-real-world-architecture-case-studies)
7. [Interview Questions & Scenarios](#7-interview-questions--scenarios)

---

## 1. Multi-Account Strategy

### 1.1 Why Multi-Account?

```
Single Account Problems:
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Single AWS Account                                │
│                                                                      │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐           │
│  │   Dev    │  │ Staging  │  │   Prod   │  │ Security │           │
│  │Resources │  │Resources │  │Resources │  │  Tools   │           │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘           │
│                                                                      │
│  Problems:                                                           │
│  ✗ Blast radius - mistake affects everything                       │
│  ✗ Permission complexity - hard to isolate access                  │
│  ✗ Service quotas shared across all workloads                      │
│  ✗ Cost allocation difficult                                       │
│  ✗ Compliance challenges                                           │
│  ✗ Audit complexity                                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

```
Multi-Account Benefits:
──────────────────────────────────────────────────────────────────────

✓ BLAST RADIUS REDUCTION
  └── Issues in Dev don't affect Prod

✓ SECURITY ISOLATION
  └── Separate IAM boundaries per account

✓ SERVICE QUOTAS
  └── Each account has its own limits

✓ COST ALLOCATION
  └── Clear billing per account/workload

✓ COMPLIANCE
  └── Different compliance controls per account

✓ TEAM AUTONOMY
  └── Teams can self-manage within guardrails
```

### 1.2 AWS Organizations

**Organization Structure**:

```
AWS Organizations Hierarchy
──────────────────────────────────────────────────────────────────────

                    ┌─────────────────────┐
                    │        Root         │
                    │  (Organization)     │
                    └──────────┬──────────┘
                               │
         ┌─────────────────────┼─────────────────────┐
         │                     │                     │
         ▼                     ▼                     ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Security OU   │  │  Infrastructure │  │   Workloads OU  │
│                 │  │       OU        │  │                 │
│ • Log Archive   │  │ • Network Hub   │  │  ┌───────────┐  │
│ • Security      │  │ • Shared Svcs   │  │  │Production │  │
│   Tooling       │  │ • DNS           │  │  │    OU     │  │
│ • Audit         │  │                 │  │  │ • App1    │  │
│                 │  │                 │  │  │ • App2    │  │
└─────────────────┘  └─────────────────┘  │  └───────────┘  │
                                          │  ┌───────────┐  │
                                          │  │Non-Prod OU│  │
                                          │  │ • Dev     │  │
                                          │  │ • Staging │  │
                                          │  └───────────┘  │
                                          └─────────────────┘

Key Concepts:
• Organization: Collection of AWS accounts
• Root: Parent container for all accounts
• OU (Organizational Unit): Logical grouping of accounts
• Account: Individual AWS account
• SCP: Service Control Policy (guardrails)
```

**AWS Control Tower**:

```
Control Tower vs Manual Setup
──────────────────────────────────────────────────────────────────────

MANUAL SETUP:                    CONTROL TOWER:
• Build everything yourself      • Pre-built best practices
• Custom SCPs                    • Managed guardrails (400+)
• Manual account provisioning    • Account Factory automation
• DIY logging/auditing           • Built-in logging (CloudTrail, Config)
• Weeks to set up                • Hours to set up

Control Tower Landing Zone:
┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│  ┌────────────────┐  ┌────────────────┐  ┌────────────────┐        │
│  │ Management     │  │ Log Archive    │  │ Audit          │        │
│  │ Account        │  │ Account        │  │ Account        │        │
│  │                │  │                │  │                │        │
│  │ • Organizations│  │ • CloudTrail   │  │ • Config Rules │        │
│  │ • Control Tower│  │ • Config Logs  │  │ • Security Hub │        │
│  │ • SSO          │  │ • VPC Flow Logs│  │ • GuardDuty    │        │
│  │ • SCPs         │  │ • S3 Access    │  │                │        │
│  └────────────────┘  └────────────────┘  └────────────────┘        │
│                                                                      │
│  Guardrails:                                                        │
│  ├── Mandatory (Detective & Preventive)                            │
│  │   └── Cannot be disabled                                        │
│  ├── Strongly Recommended                                           │
│  │   └── Best practice defaults                                    │
│  └── Elective                                                       │
│      └── Optional controls                                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 1.3 Service Control Policies (SCPs)

**SCP Fundamentals**:

```
SCP Evaluation
──────────────────────────────────────────────────────────────────────

SCPs are Permission BOUNDARIES, not grants!

                    Account's Effective Permissions
                    ═══════════════════════════════
                    
    ┌─────────────────┐         ┌─────────────────┐
    │  SCP Allows     │         │  IAM Allows     │
    │  (from all      │    ∩    │  (User/Role     │
    │   parent OUs)   │         │   policies)     │
    └────────┬────────┘         └────────┬────────┘
             │                           │
             └───────────┬───────────────┘
                         │
                         ▼
              ┌─────────────────────┐
              │ EFFECTIVE PERMISSIONS │
              │  (Intersection only) │
              └─────────────────────┘

Key Rules:
• SCPs don't grant permissions
• SCPs restrict maximum permissions
• Root account is exempt (use sparingly!)
• SCPs are inherited down the hierarchy
```

**Common SCP Patterns**:

```json
// 1. DENY LEAVING ORGANIZATION
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DenyLeaveOrganization",
      "Effect": "Deny",
      "Action": "organizations:LeaveOrganization",
      "Resource": "*"
    }
  ]
}

// 2. RESTRICT REGIONS
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DenyOutsideRegions",
      "Effect": "Deny",
      "NotAction": [
        "cloudfront:*",
        "iam:*",
        "route53:*",
        "support:*",
        "budgets:*",
        "organizations:*"
      ],
      "Resource": "*",
      "Condition": {
        "StringNotEquals": {
          "aws:RequestedRegion": [
            "us-east-1",
            "us-west-2",
            "eu-west-1"
          ]
        }
      }
    }
  ]
}

// 3. REQUIRE IMDSV2 FOR EC2
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "RequireIMDSv2",
      "Effect": "Deny",
      "Action": "ec2:RunInstances",
      "Resource": "arn:aws:ec2:*:*:instance/*",
      "Condition": {
        "StringNotEquals": {
          "ec2:MetadataHttpTokens": "required"
        }
      }
    }
  ]
}

// 4. PROTECT SECURITY RESOURCES
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DenyCloudTrailModification",
      "Effect": "Deny",
      "Action": [
        "cloudtrail:DeleteTrail",
        "cloudtrail:StopLogging",
        "cloudtrail:UpdateTrail"
      ],
      "Resource": "*",
      "Condition": {
        "StringNotLike": {
          "aws:PrincipalArn": [
            "arn:aws:iam::*:role/SecurityAdmin"
          ]
        }
      }
    }
  ]
}

// 5. ENFORCE TAGGING
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "DenyUntaggedResources",
      "Effect": "Deny",
      "Action": [
        "ec2:RunInstances",
        "ec2:CreateVolume",
        "rds:CreateDBInstance"
      ],
      "Resource": "*",
      "Condition": {
        "Null": {
          "aws:RequestTag/Environment": "true",
          "aws:RequestTag/Owner": "true"
        }
      }
    }
  ]
}
```

### 1.4 Account Vending (Account Factory)

```
Account Factory Pipeline
──────────────────────────────────────────────────────────────────────

Request                         Provisioning                 Ready
───────                         ────────────                 ─────

┌──────────────┐               ┌──────────────┐         ┌──────────────┐
│  ServiceNow  │               │   Control    │         │    Ready     │
│  or Ticket   │──────────────►│   Tower      │────────►│   Account    │
│  Request     │               │   Account    │         │              │
│              │               │   Factory    │         │ • VPC        │
│ • Name       │               │              │         │ • IAM Roles  │
│ • Owner      │               │   OR         │         │ • Security   │
│ • Budget     │               │              │         │ • Networking │
│ • OU         │               │   Custom     │         │ • Baseline   │
│              │               │   Terraform  │         │              │
└──────────────┘               └──────────────┘         └──────────────┘

Custom Account Factory (Terraform):
───────────────────────────────────────────────────────────────────────

# account-factory/main.tf
resource "aws_organizations_account" "workload" {
  name      = var.account_name
  email     = var.account_email
  parent_id = var.organizational_unit_id
  role_name = "OrganizationAccountAccessRole"

  lifecycle {
    ignore_changes = [role_name]
  }
}

# Baseline deployment via StackSets
resource "aws_cloudformation_stack_set_instance" "baseline" {
  account_id     = aws_organizations_account.workload.id
  region         = var.home_region
  stack_set_name = aws_cloudformation_stack_set.baseline.name
}

# VPC provisioning
module "vpc" {
  source = "../modules/vpc"
  providers = {
    aws = aws.workload_account
  }
  
  vpc_cidr         = var.vpc_cidr
  environment      = var.environment
  transit_gateway_id = data.aws_ec2_transit_gateway.main.id
}

# Security baseline
module "security_baseline" {
  source = "../modules/security-baseline"
  providers = {
    aws = aws.workload_account
  }
  
  log_archive_bucket = var.log_archive_bucket
  security_account_id = var.security_account_id
}
```

### 1.5 Cross-Account Access Patterns

```
Cross-Account Access Methods
──────────────────────────────────────────────────────────────────────

1. IAM ROLE ASSUMPTION (Most Common)
   ┌─────────────────┐         ┌─────────────────┐
   │   Account A     │         │   Account B     │
   │                 │         │                 │
   │  ┌───────────┐  │ assume  │  ┌───────────┐  │
   │  │ Developer │──┼─────────┼──│ CrossAcct │  │
   │  │   Role    │  │  role   │  │   Role    │  │
   │  └───────────┘  │         │  └───────────┘  │
   │                 │         │                 │
   └─────────────────┘         └─────────────────┘

   Trust Policy (Account B):
   {
     "Principal": {
       "AWS": "arn:aws:iam::ACCOUNT_A:root"
     },
     "Action": "sts:AssumeRole",
     "Condition": {
       "StringEquals": {
         "sts:ExternalId": "unique-secret-id"
       }
     }
   }

2. RESOURCE-BASED POLICIES
   ┌─────────────────┐         ┌─────────────────┐
   │   Account A     │         │   Account B     │
   │                 │         │                 │
   │  ┌───────────┐  │ direct  │  ┌───────────┐  │
   │  │   Lambda  │──┼─────────┼──│  S3 Bucket│  │
   │  └───────────┘  │ access  │  │  (Policy) │  │
   │                 │         │  └───────────┘  │
   └─────────────────┘         └─────────────────┘

   Bucket Policy:
   {
     "Principal": {
       "AWS": "arn:aws:iam::ACCOUNT_A:role/LambdaRole"
     },
     "Action": "s3:GetObject",
     "Resource": "arn:aws:s3:::bucket/*"
   }

3. AWS RESOURCE ACCESS MANAGER (RAM)
   Share resources: Transit Gateway, Subnets, License Manager
   
   ┌─────────────────────────────────────────────────────────────┐
   │ Management Account                                          │
   │                                                             │
   │  ┌─────────────────────────────────────────────────────┐   │
   │  │            Transit Gateway                           │   │
   │  └───────────────────────┬─────────────────────────────┘   │
   │                          │ RAM Share                        │
   │           ┌──────────────┼──────────────┐                  │
   └───────────┼──────────────┼──────────────┼──────────────────┘
               │              │              │
               ▼              ▼              ▼
        ┌──────────┐   ┌──────────┐   ┌──────────┐
        │ Account A│   │ Account B│   │ Account C│
        │  (Dev)   │   │ (Staging)│   │  (Prod)  │
        └──────────┘   └──────────┘   └──────────┘
```

### 1.6 AWS Identity Center (SSO)

```
Identity Center Architecture
──────────────────────────────────────────────────────────────────────

                    ┌────────────────────────┐
                    │   Identity Provider    │
                    │  (Okta, Azure AD,      │
                    │   AWS Directory)       │
                    └───────────┬────────────┘
                                │ SAML/SCIM
                                ▼
                    ┌────────────────────────┐
                    │   AWS Identity Center  │
                    │   (Management Account) │
                    │                        │
                    │  ┌──────────────────┐  │
                    │  │ Permission Sets  │  │
                    │  │ • Admin          │  │
                    │  │ • Developer      │  │
                    │  │ • ReadOnly       │  │
                    │  │ • SecurityAudit  │  │
                    │  └──────────────────┘  │
                    │                        │
                    │  ┌──────────────────┐  │
                    │  │ Account          │  │
                    │  │ Assignments      │  │
                    │  │ User/Group →     │  │
                    │  │ Account +        │  │
                    │  │ Permission Set   │  │
                    │  └──────────────────┘  │
                    └───────────┬────────────┘
                                │
         ┌──────────────────────┼──────────────────────┐
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│   Dev Account   │  │ Staging Account │  │  Prod Account   │
│                 │  │                 │  │                 │
│ Role: Developer │  │ Role: Developer │  │ Role: ReadOnly  │
│ Role: Admin     │  │ Role: Admin     │  │ Role: Admin     │
└─────────────────┘  └─────────────────┘  └─────────────────┘

Permission Set Example:
{
  "Name": "DeveloperAccess",
  "SessionDuration": "PT8H",
  "ManagedPolicies": [
    "arn:aws:iam::aws:policy/PowerUserAccess"
  ],
  "InlinePolicy": {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Sid": "DenyIAMChanges",
        "Effect": "Deny",
        "Action": [
          "iam:CreateUser",
          "iam:DeleteUser",
          "iam:CreateRole",
          "iam:DeleteRole"
        ],
        "Resource": "*"
      }
    ]
  }
}
```

### 1.7 Multi-Account Network Architecture

```
Hub-and-Spoke Network
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Network Hub Account                               │
│                                                                      │
│  ┌───────────────────────────────────────────────────────────────┐ │
│  │                     Transit Gateway                            │ │
│  │                                                                │ │
│  │  Route Tables:                                                 │ │
│  │  ├── Production (isolated)                                    │ │
│  │  ├── Non-Production (can access shared)                       │ │
│  │  ├── Shared Services (accessible by all)                     │ │
│  │  └── Inspection (routes through firewall)                     │ │
│  │                                                                │ │
│  └───────────────────────────────────────────────────────────────┘ │
│                              │                                       │
│  ┌───────────────────────────┼───────────────────────────────────┐ │
│  │                           │                                    │ │
│  │  Egress VPC              │  Ingress VPC                       │ │
│  │  ├── NAT Gateways        │  ├── ALB/NLB                      │ │
│  │  ├── Network Firewall    │  ├── WAF                          │ │
│  │  └── Internet Gateway    │  └── CloudFront Origin            │ │
│  │                           │                                    │ │
│  └───────────────────────────┼───────────────────────────────────┘ │
│                              │                                       │
│  ┌───────────────────────────┼───────────────────────────────────┐ │
│  │  On-Premise Connectivity  │                                    │ │
│  │  ├── Direct Connect      │                                    │ │
│  │  └── Site-to-Site VPN    │                                    │ │
│  └───────────────────────────┘                                     │ │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
         │                    │                    │
         │ TGW Attachment     │ TGW Attachment     │ TGW Attachment
         ▼                    ▼                    ▼
┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐
│  Prod Account   │  │  Dev Account    │  │ Shared Services │
│                 │  │                 │  │                 │
│  VPC: 10.1.0.0  │  │  VPC: 10.2.0.0  │  │  VPC: 10.100.0  │
│  ├── App Subnet │  │  ├── App Subnet │  │  ├── DNS        │
│  └── DB Subnet  │  │  └── DB Subnet  │  │  ├── AD         │
│                 │  │                 │  │  └── Tools      │
└─────────────────┘  └─────────────────┘  └─────────────────┘

Traffic Flow Examples:
─────────────────────────────────────────────────────────────────────
Dev → Internet:      Dev VPC → TGW → Egress VPC → NAT → IGW
Prod ↔ Shared:       Prod VPC → TGW → Shared VPC
Internet → Prod:     CloudFront → Ingress VPC → TGW → Prod VPC
Prod → On-Prem:      Prod VPC → TGW → Direct Connect
```

---

## 2. Disaster Recovery & Business Continuity

### 2.1 DR Fundamentals

**Key Metrics**:

```
DR Metrics
──────────────────────────────────────────────────────────────────────

RTO (Recovery Time Objective)
└── Maximum acceptable downtime
└── "How long can we be down?"
└── Example: 4 hours RTO = must recover within 4 hours

RPO (Recovery Point Objective)
└── Maximum acceptable data loss
└── "How much data can we lose?"
└── Example: 1 hour RPO = max 1 hour of data loss

              Data Loss          Downtime
    ◄────────────────────────────────────────────►
                    │
    RPO             │              RTO
    ◄───────────────│──────────────────────────►
                    │
              Disaster
              Occurs
```

### 2.2 DR Strategies

```
DR Strategies Spectrum
──────────────────────────────────────────────────────────────────────

Cost ←─────────────────────────────────────────────────────→ RTO/RPO
Low                                                           Fast

┌─────────────┬─────────────┬─────────────┬─────────────────────────┐
│   Backup &  │  Pilot      │   Warm      │   Multi-Site            │
│   Restore   │  Light      │   Standby   │   Active-Active         │
├─────────────┼─────────────┼─────────────┼─────────────────────────┤
│             │             │             │                         │
│ RTO: Hours  │ RTO: 10min- │ RTO: Minutes│ RTO: Real-time         │
│      to Days│      Hours  │             │      (near-zero)       │
│             │             │             │                         │
│ RPO: Hours  │ RPO: Minutes│ RPO: Seconds│ RPO: Near-zero         │
│             │      to     │      to     │                         │
│             │      Hours  │      Minutes│                         │
│             │             │             │                         │
│ Cost: $     │ Cost: $$    │ Cost: $$$   │ Cost: $$$$             │
│             │             │             │                         │
└─────────────┴─────────────┴─────────────┴─────────────────────────┘
```

**Strategy 1: Backup & Restore**

```
Backup & Restore Architecture
──────────────────────────────────────────────────────────────────────

Primary Region (us-east-1)              DR Region (us-west-2)
─────────────────────────              ──────────────────────

┌─────────────────────────┐            ┌─────────────────────────┐
│                         │            │                         │
│  ┌─────┐  ┌─────┐      │            │    (Empty - No infra)   │
│  │ EC2 │  │ RDS │      │            │                         │
│  └──┬──┘  └──┬──┘      │            │                         │
│     │        │         │            │                         │
│     ▼        ▼         │            │                         │
│  ┌──────────────────┐  │  Cross-    │  ┌──────────────────┐   │
│  │   S3 Bucket      │──┼──Region───►│  │   S3 Bucket      │   │
│  │   (Backups)      │  │  Replication│  │   (DR Backups)   │   │
│  │   • AMIs         │  │            │  │                  │   │
│  │   • DB Snapshots │  │            │  │                  │   │
│  │   • App Data     │  │            │  │                  │   │
│  └──────────────────┘  │            │  └──────────────────┘   │
│                         │            │                         │
└─────────────────────────┘            └─────────────────────────┘

Recovery Process:
1. Restore AMIs from S3/shared
2. Restore RDS from snapshot
3. Launch infrastructure (CloudFormation/Terraform)
4. Restore application data
5. Update DNS

Pros: Lowest cost
Cons: Longest recovery time (hours to days)
```

**Strategy 2: Pilot Light**

```
Pilot Light Architecture
──────────────────────────────────────────────────────────────────────

Primary Region (us-east-1)              DR Region (us-west-2)
─────────────────────────              ──────────────────────

┌─────────────────────────┐            ┌─────────────────────────┐
│                         │            │                         │
│  ┌─────────────────┐    │            │  ┌─────────────────┐    │
│  │   ALB (Active)  │    │            │  │   ALB (Stopped) │    │
│  └────────┬────────┘    │            │  └────────┬────────┘    │
│           │             │            │           │             │
│  ┌────────▼────────┐    │            │  ┌────────▼────────┐    │
│  │  App Servers    │    │            │  │  App Servers    │    │
│  │  (Running)      │    │            │  │  (Stopped/Min)  │    │
│  │  ASG: 4         │    │            │  │  ASG: 0         │    │
│  └────────┬────────┘    │            │  └────────┬────────┘    │
│           │             │            │           │             │
│  ┌────────▼────────┐    │  Async     │  ┌────────▼────────┐    │
│  │   Aurora        │────┼──Replica──►│  │   Aurora        │    │
│  │   (Primary)     │    │            │  │   (Read Replica)│    │
│  │   Writer        │    │            │  │   Promotable    │    │
│  └─────────────────┘    │            │  └─────────────────┘    │
│                         │            │                         │
└─────────────────────────┘            └─────────────────────────┘

Core Infrastructure Running:
• Database replica (always on, syncing)
• Network infrastructure (VPC, subnets)
• Security groups, IAM roles

Scaled Down/Off:
• Application servers
• Load balancers
• Caches

Recovery Process (10 min - 1 hour):
1. Promote DB replica to primary
2. Scale up ASG (0 → 4)
3. Start/configure ALB
4. Update Route 53 (health check failover)
```

**Strategy 3: Warm Standby**

```
Warm Standby Architecture
──────────────────────────────────────────────────────────────────────

Primary Region (us-east-1)              DR Region (us-west-2)
─────────────────────────              ──────────────────────

┌─────────────────────────┐            ┌─────────────────────────┐
│        100% Traffic     │            │     0% Traffic (Ready)  │
│                         │            │                         │
│  ┌─────────────────┐    │            │  ┌─────────────────┐    │
│  │   ALB (Active)  │    │            │  │   ALB (Active)  │    │
│  └────────┬────────┘    │            │  └────────┬────────┘    │
│           │             │            │           │             │
│  ┌────────▼────────┐    │            │  ┌────────▼────────┐    │
│  │  App Servers    │    │            │  │  App Servers    │    │
│  │  ASG: 10        │    │            │  │  ASG: 2 (min)   │    │
│  │  (Full Scale)   │    │            │  │  (Reduced)      │    │
│  └────────┬────────┘    │            │  └────────┬────────┘    │
│           │             │            │           │             │
│  ┌────────▼────────┐    │  Sync     │  ┌────────▼────────┐    │
│  │   Aurora        │────┼──Global───►│  │   Aurora        │    │
│  │   Global DB     │    │  Database  │  │   Global DB     │    │
│  │   (Writer)      │    │            │  │   (Reader)      │    │
│  └─────────────────┘    │            │  └─────────────────┘    │
│                         │            │                         │
│  ┌─────────────────┐    │  Cross-    │  ┌─────────────────┐    │
│  │  ElastiCache    │────┼──Region───►│  │  ElastiCache    │    │
│  │  Global Store   │    │  Replication│  │  Global Store   │    │
│  └─────────────────┘    │            │  └─────────────────┘    │
│                         │            │                         │
└─────────────────────────┘            └─────────────────────────┘

All Infrastructure Running (scaled down):
• Application servers (minimal capacity)
• Database (read replica, promotable)
• Caches (replicated)
• Load balancers (active)

Recovery Process (minutes):
1. Promote Aurora secondary to primary (<1 min)
2. Scale up ASG (2 → 10)
3. Route 53 failover (automatic via health checks)
```

**Strategy 4: Multi-Site Active-Active**

```
Multi-Site Active-Active Architecture
──────────────────────────────────────────────────────────────────────

                    ┌─────────────────────────┐
                    │      Route 53           │
                    │  Latency/Geolocation    │
                    │      Routing            │
                    └───────────┬─────────────┘
                                │
              ┌─────────────────┴─────────────────┐
              │                                   │
              ▼                                   ▼
Primary Region (us-east-1)              Secondary Region (eu-west-1)
       50% Traffic                            50% Traffic

┌─────────────────────────┐            ┌─────────────────────────┐
│                         │            │                         │
│  ┌─────────────────┐    │            │  ┌─────────────────┐    │
│  │   CloudFront    │    │            │  │   CloudFront    │    │
│  └────────┬────────┘    │            │  └────────┬────────┘    │
│           │             │            │           │             │
│  ┌────────▼────────┐    │            │  ┌────────▼────────┐    │
│  │   ALB           │    │            │  │   ALB           │    │
│  └────────┬────────┘    │            │  └────────┬────────┘    │
│           │             │            │           │             │
│  ┌────────▼────────┐    │            │  ┌────────▼────────┐    │
│  │  App Servers    │    │            │  │  App Servers    │    │
│  │  ASG: 10        │    │            │  │  ASG: 10        │    │
│  └────────┬────────┘    │            │  └────────┬────────┘    │
│           │             │            │           │             │
│  ┌────────▼────────┐    │            │  ┌────────▼────────┐    │
│  │   Aurora        │◄───┼────────────┼───►│   Aurora      │    │
│  │   Global DB     │    │   Bi-dir   │  │   Global DB   │    │
│  │   (Writer)      │    │   Sync     │  │   (Writer*)   │    │
│  └─────────────────┘    │            │  └─────────────────┘    │
│                         │            │                         │
│  ┌─────────────────┐    │            │  ┌─────────────────┐    │
│  │   DynamoDB      │◄───┼────────────┼───►│  DynamoDB     │    │
│  │   Global Table  │    │   Async    │  │  Global Table │    │
│  └─────────────────┘    │   Sync     │  └─────────────────┘    │
│                         │            │                         │
└─────────────────────────┘            └─────────────────────────┘

*Write forwarding or application-level routing

Considerations:
• Conflict resolution (last writer wins for DynamoDB)
• Data consistency (eventual vs strong)
• Application must be region-aware
• Cost: 2x infrastructure
```

### 2.3 AWS DR Services

**Aurora Global Database**:

```
Aurora Global Database
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                       Aurora Global Database                         │
│                                                                      │
│  Primary Region                    Secondary Region(s)               │
│  (Read/Write)                      (Read Only)                       │
│                                                                      │
│  ┌─────────────────┐              ┌─────────────────┐               │
│  │                 │              │                 │               │
│  │   ┌─────────┐   │   Storage    │   ┌─────────┐   │               │
│  │   │  Writer │   │   Replication│   │  Reader │   │               │
│  │   │ Instance│   │   <1 second  │   │ Instance│   │               │
│  │   └─────────┘   │──────────────│   └─────────┘   │               │
│  │                 │              │                 │               │
│  │   ┌─────────┐   │              │   ┌─────────┐   │               │
│  │   │  Reader │   │              │   │  Reader │   │               │
│  │   │ Instance│   │              │   │ Instance│   │               │
│  │   └─────────┘   │              │   └─────────┘   │               │
│  │                 │              │                 │               │
│  └─────────────────┘              └─────────────────┘               │
│                                                                      │
│  Failover: Managed or Manual                                        │
│  RTO: <1 minute (managed failover)                                  │
│  RPO: <1 second (typical lag)                                       │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

# Terraform
resource "aws_rds_global_cluster" "main" {
  global_cluster_identifier = "global-aurora-cluster"
  engine                    = "aurora-postgresql"
  engine_version            = "15.4"
  database_name             = "myapp"
}

resource "aws_rds_cluster" "primary" {
  provider                  = aws.primary
  cluster_identifier        = "aurora-primary"
  global_cluster_identifier = aws_rds_global_cluster.main.id
  engine                    = aws_rds_global_cluster.main.engine
  engine_version            = aws_rds_global_cluster.main.engine_version
  master_username           = "admin"
  master_password           = var.db_password
  
  # Primary cluster settings
}

resource "aws_rds_cluster" "secondary" {
  provider                  = aws.secondary
  cluster_identifier        = "aurora-secondary"
  global_cluster_identifier = aws_rds_global_cluster.main.id
  engine                    = aws_rds_global_cluster.main.engine
  engine_version            = aws_rds_global_cluster.main.engine_version
  
  # No credentials needed - pulls from global cluster
  depends_on = [aws_rds_cluster.primary]
}
```

**S3 Cross-Region Replication**:

```
S3 CRR Configuration
──────────────────────────────────────────────────────────────────────

# Terraform
resource "aws_s3_bucket" "source" {
  bucket = "my-app-data-primary"
}

resource "aws_s3_bucket_versioning" "source" {
  bucket = aws_s3_bucket.source.id
  versioning_configuration {
    status = "Enabled"  # Required for replication
  }
}

resource "aws_s3_bucket_replication_configuration" "replication" {
  bucket = aws_s3_bucket.source.id
  role   = aws_iam_role.replication.arn

  rule {
    id     = "replicate-all"
    status = "Enabled"

    filter {
      prefix = ""  # Replicate everything
    }

    destination {
      bucket        = aws_s3_bucket.destination.arn
      storage_class = "STANDARD"
      
      # Optional: Different account
      account = var.dr_account_id
      
      access_control_translation {
        owner = "Destination"
      }
      
      # Replication Time Control (RTC) - 15 min SLA
      replication_time {
        status = "Enabled"
        time {
          minutes = 15
        }
      }
      
      metrics {
        status = "Enabled"
        event_threshold {
          minutes = 15
        }
      }
    }

    delete_marker_replication {
      status = "Enabled"
    }
  }
}
```

**Route 53 Health Checks & Failover**:

```
Route 53 Failover Configuration
──────────────────────────────────────────────────────────────────────

                    ┌─────────────────────────┐
                    │      Route 53           │
                    │   Health Check          │
                    │   ┌─────────────────┐   │
                    │   │ Check Primary   │   │
                    │   │ /health endpoint│   │
                    │   │ every 10 sec    │   │
                    │   └────────┬────────┘   │
                    │            │            │
                    │   Healthy? │            │
                    │   ┌────────┴────────┐   │
                    │   │                 │   │
                    │  YES               NO   │
                    │   │                 │   │
                    └───┼─────────────────┼───┘
                        ▼                 ▼
               ┌─────────────────┐ ┌─────────────────┐
               │ Primary (Active)│ │Secondary (Standby)│
               │ us-east-1       │ │ us-west-2         │
               └─────────────────┘ └─────────────────┘

# Terraform
resource "aws_route53_health_check" "primary" {
  fqdn              = "api-primary.example.com"
  port              = 443
  type              = "HTTPS"
  resource_path     = "/health"
  failure_threshold = 3
  request_interval  = 10

  tags = {
    Name = "primary-health-check"
  }
}

resource "aws_route53_record" "failover_primary" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "api.example.com"
  type    = "A"

  failover_routing_policy {
    type = "PRIMARY"
  }

  set_identifier  = "primary"
  health_check_id = aws_route53_health_check.primary.id

  alias {
    name                   = aws_lb.primary.dns_name
    zone_id                = aws_lb.primary.zone_id
    evaluate_target_health = true
  }
}

resource "aws_route53_record" "failover_secondary" {
  zone_id = aws_route53_zone.main.zone_id
  name    = "api.example.com"
  type    = "A"

  failover_routing_policy {
    type = "SECONDARY"
  }

  set_identifier = "secondary"

  alias {
    name                   = aws_lb.secondary.dns_name
    zone_id                = aws_lb.secondary.zone_id
    evaluate_target_health = true
  }
}
```

### 2.4 DR Testing & Runbooks

```
DR Testing Strategy
──────────────────────────────────────────────────────────────────────

Testing Frequency:
┌─────────────────────────────────────────────────────────────────────┐
│ Test Type           │ Frequency    │ Scope                         │
├─────────────────────┼──────────────┼───────────────────────────────┤
│ Backup Verification │ Daily        │ Automated restore test        │
│ Component Failover  │ Monthly      │ Single component (DB, cache)  │
│ Full DR Drill       │ Quarterly    │ Complete environment failover │
│ Chaos Engineering   │ Continuous   │ Random failure injection      │
└─────────────────────┴──────────────┴───────────────────────────────┘

DR Runbook Template:
──────────────────────────────────────────────────────────────────────

## Runbook: Database Failover to DR Region

### Pre-requisites
- [ ] Verify DR region infrastructure is running
- [ ] Confirm replication lag < 1 second
- [ ] Notify stakeholders (Slack: #incidents)

### Failover Steps

**Step 1: Verify Current State (2 min)**
```bash
aws rds describe-global-clusters \
  --global-cluster-identifier my-global-cluster \
  --query 'GlobalClusters[0].GlobalClusterMembers'
```
Expected: Primary in us-east-1, Secondary in us-west-2

**Step 2: Initiate Failover (5 min)**
```bash
aws rds failover-global-cluster \
  --global-cluster-identifier my-global-cluster \
  --target-db-cluster-identifier arn:aws:rds:us-west-2:ACCOUNT:cluster:aurora-secondary
```

**Step 3: Verify Failover Complete**
```bash
# Check cluster status
aws rds describe-db-clusters \
  --db-cluster-identifier aurora-secondary \
  --query 'DBClusters[0].Status'
```
Expected: "available"

**Step 4: Update Application Configuration**
- Application should auto-discover new endpoint via Global DB
- If not, update connection strings

**Step 5: Verify Application Health**
```bash
curl https://api.example.com/health
```

### Rollback (if needed)
Repeat steps with original primary as target

### Post-Failover
- [ ] Update incident ticket
- [ ] Schedule post-mortem
- [ ] Plan failback window
```

### 2.5 AWS Elastic Disaster Recovery (DRS)

```
AWS DRS Architecture
──────────────────────────────────────────────────────────────────────

Source Environment                    AWS DR Region
(On-Premise or Other Cloud)           (Recovery Target)

┌─────────────────────────┐          ┌─────────────────────────┐
│                         │          │                         │
│  ┌─────────────────┐    │          │  ┌─────────────────┐    │
│  │   Source        │    │  Continuous│  │   Staging Area  │    │
│  │   Server 1      │────┼──Replication│  │   (Low-cost)    │    │
│  │   (Windows)     │    │          │  │   EBS Snapshots │    │
│  └─────────────────┘    │          │  └─────────────────┘    │
│                         │          │           │             │
│  ┌─────────────────┐    │          │           │ Launch      │
│  │   Source        │    │          │           │ (on demand) │
│  │   Server 2      │────┼──────────│           ▼             │
│  │   (Linux)       │    │          │  ┌─────────────────┐    │
│  └─────────────────┘    │          │  │   Recovery      │    │
│                         │          │  │   Instances     │    │
│  ┌─────────────────┐    │          │  │   (Full size)   │    │
│  │   DRS           │    │          │  └─────────────────┘    │
│  │   Agent         │    │          │                         │
│  │   (Installed)   │    │          │                         │
│  └─────────────────┘    │          │                         │
│                         │          │                         │
└─────────────────────────┘          └─────────────────────────┘

Key Features:
• Continuous block-level replication
• Sub-second RPO
• Minutes RTO (launch recovery instances)
• Non-disruptive DR drills
• Support for Windows, Linux, databases
• Cost-effective (staging uses minimal resources)

Use Cases:
• On-premise to AWS DR
• AWS region to region
• Other cloud to AWS
```

---

## 3. Security Deep Dive

### 3.1 AWS Security Services Ecosystem

```
Security Services Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                        Security Hub (Central Dashboard)             │
│   Aggregates findings from all security services                    │
└─────────────────────────────────────────────────────────────────────┘
          ▲           ▲           ▲           ▲           ▲
          │           │           │           │           │
┌─────────┴─────────┐ │  ┌────────┴────────┐ │  ┌────────┴────────┐
│    GuardDuty      │ │  │    Inspector    │ │  │     Macie       │
│                   │ │  │                 │ │  │                 │
│  Threat Detection │ │  │  Vulnerability  │ │  │  Data Discovery │
│  • Network        │ │  │  Assessment     │ │  │  • PII/PHI      │
│  • Account        │ │  │  • EC2          │ │  │  • Financial    │
│  • S3             │ │  │  • ECR          │ │  │  • Credentials  │
│  • EKS            │ │  │  • Lambda       │ │  │                 │
└───────────────────┘ │  └─────────────────┘ │  └─────────────────┘
                      │                      │
          ┌───────────┴───────────┐        │
          │   IAM Access Analyzer │        │
          │                       │        │
          │  External Access      │        │
          │  • S3 buckets         │        │
          │  • IAM roles          │        │
          │  • KMS keys           │        │
          │  • Lambda functions   │        │
          │  • SQS queues         │        │
          └───────────────────────┘        │
                                           │
          ┌────────────────────────────────┴───────────────────────┐
          │                    AWS Config                          │
          │                                                        │
          │  Configuration Compliance & Change Tracking           │
          │  • 300+ managed rules                                 │
          │  • Custom rules (Lambda)                              │
          │  • Remediation actions                                │
          └────────────────────────────────────────────────────────┘
```

### 3.2 Amazon GuardDuty

```
GuardDuty Architecture
──────────────────────────────────────────────────────────────────────

Data Sources:                         Analysis Engine:
                                      
┌─────────────────┐                   ┌─────────────────────────────┐
│  VPC Flow Logs  │──────────────────►│                             │
└─────────────────┘                   │    Machine Learning         │
                                      │    Anomaly Detection        │
┌─────────────────┐                   │                             │
│  DNS Logs       │──────────────────►│    Threat Intelligence      │
└─────────────────┘                   │    • AWS threat feeds       │
                                      │    • Partner feeds          │
┌─────────────────┐                   │    • CrowdStrike, Proofpoint│
│  CloudTrail     │──────────────────►│                             │
│  Events         │                   │    Behavioral Models        │
└─────────────────┘                   │    • Baseline normal        │
                                      │    • Detect deviations      │
┌─────────────────┐                   │                             │
│  S3 Data Events │──────────────────►└──────────────┬──────────────┘
└─────────────────┘                                  │
                                                     ▼
┌─────────────────┐                   ┌─────────────────────────────┐
│  EKS Audit Logs │──────────────────►│         Findings            │
└─────────────────┘                   │    • Severity (Low/Med/High)│
                                      │    • Finding Type           │
┌─────────────────┐                   │    • Affected Resource      │
│  Lambda Network │──────────────────►│    • Remediation            │
│  Activity       │                   └─────────────────────────────┘
└─────────────────┘

Finding Types:
──────────────────────────────────────────────────────────────────────

EC2 Findings:
• Backdoor:EC2/DenialOfService.Tcp
• CryptoCurrency:EC2/BitcoinTool.B
• Trojan:EC2/BlackholeTraffic
• UnauthorizedAccess:EC2/SSHBruteForce

IAM Findings:
• CredentialAccess:IAMUser/AnomalousBehavior
• UnauthorizedAccess:IAMUser/InstanceCredentialExfiltration

S3 Findings:
• Policy:S3/BucketBlockPublicAccessDisabled
• Stealth:S3/ServerAccessLoggingDisabled
• UnauthorizedAccess:S3/TorIPCaller

# Enable GuardDuty (Terraform)
resource "aws_guardduty_detector" "main" {
  enable = true
  
  datasources {
    s3_logs {
      enable = true
    }
    kubernetes {
      audit_logs {
        enable = true
      }
    }
    malware_protection {
      scan_ec2_instance_with_findings {
        ebs_volumes {
          enable = true
        }
      }
    }
  }

  finding_publishing_frequency = "FIFTEEN_MINUTES"
  
  tags = {
    Environment = "production"
  }
}

# Multi-account with Organizations
resource "aws_guardduty_organization_admin_account" "main" {
  admin_account_id = var.security_account_id
}
```

### 3.3 AWS Security Hub

```
Security Hub Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                         Security Hub                                 │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │                    Standards & Controls                          ││
│  │                                                                  ││
│  │  ┌──────────────┐ ┌──────────────┐ ┌──────────────────────────┐ ││
│  │  │AWS Foundation│ │    CIS AWS   │ │   PCI DSS v3.2.1        │ ││
│  │  │Security Best │ │  Benchmark   │ │                         │ ││
│  │  │Practices v1.0│ │  v1.4.0      │ │   (For payment card)    │ ││
│  │  │              │ │              │ │                         │ ││
│  │  │ 33 controls  │ │ 49 controls  │ │   31 controls          │ ││
│  │  └──────────────┘ └──────────────┘ └──────────────────────────┘ ││
│  │                                                                  ││
│  │  ┌──────────────┐ ┌──────────────┐                             ││
│  │  │   NIST 800   │ │   Custom     │                             ││
│  │  │   -53 Rev 5  │ │  Standards   │                             ││
│  │  │              │ │              │                             ││
│  │  │ 163 controls │ │  Your rules  │                             ││
│  │  └──────────────┘ └──────────────┘                             ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │                   Integrated Services                            ││
│  │                                                                  ││
│  │  GuardDuty → Inspector → Macie → Firewall Manager → IAM Access ││
│  │              Analyzer → Config → Detective → Audit Manager      ││
│  └─────────────────────────────────────────────────────────────────┘│
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────────┐│
│  │                   Dashboard & Automation                         ││
│  │                                                                  ││
│  │  • Security Score (0-100%)                                      ││
│  │  • Finding aggregation across accounts/regions                   ││
│  │  • Custom insights                                              ││
│  │  • Automated response (EventBridge → Lambda)                    ││
│  └─────────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────────┘

# Terraform Configuration
resource "aws_securityhub_account" "main" {}

resource "aws_securityhub_standards_subscription" "cis" {
  depends_on    = [aws_securityhub_account.main]
  standards_arn = "arn:aws:securityhub:::ruleset/cis-aws-foundations-benchmark/v/1.4.0"
}

resource "aws_securityhub_standards_subscription" "aws_foundational" {
  depends_on    = [aws_securityhub_account.main]
  standards_arn = "arn:aws:securityhub:${var.region}::standards/aws-foundational-security-best-practices/v/1.0.0"
}

# Auto-remediation with EventBridge
resource "aws_cloudwatch_event_rule" "security_hub_findings" {
  name        = "security-hub-high-findings"
  description = "Capture high severity Security Hub findings"

  event_pattern = jsonencode({
    source      = ["aws.securityhub"]
    detail-type = ["Security Hub Findings - Imported"]
    detail = {
      findings = {
        Severity = {
          Label = ["HIGH", "CRITICAL"]
        }
        Workflow = {
          Status = ["NEW"]
        }
      }
    }
  })
}

resource "aws_cloudwatch_event_target" "remediation" {
  rule      = aws_cloudwatch_event_rule.security_hub_findings.name
  target_id = "remediation-lambda"
  arn       = aws_lambda_function.remediation.arn
}
```

### 3.4 IAM Access Analyzer

```
IAM Access Analyzer
──────────────────────────────────────────────────────────────────────

Analysis Types:

1. External Access Analysis (Organization/Account scope)
   └── Finds resources shared outside your zone of trust
   
2. Unused Access Analysis (New)
   └── Finds unused IAM roles, access keys, permissions

┌─────────────────────────────────────────────────────────────────────┐
│                    Supported Resource Types                          │
├─────────────────────────────────────────────────────────────────────┤
│ • S3 buckets          • Lambda functions                            │
│ • IAM roles           • SQS queues                                  │
│ • KMS keys            • Secrets Manager secrets                     │
│ • SNS topics          • EBS volume snapshots                        │
│ • ECR repositories    • RDS DB snapshots                            │
│ • EFS file systems    • DynamoDB streams/tables                     │
└─────────────────────────────────────────────────────────────────────┘

Finding Example:
──────────────────────────────────────────────────────────────────────

{
  "analyzedAt": "2024-01-15T10:30:00Z",
  "condition": {},
  "createdAt": "2024-01-15T10:30:00Z",
  "id": "finding-123",
  "isPublic": true,
  "principal": {
    "AWS": "*"
  },
  "resource": "arn:aws:s3:::my-bucket",
  "resourceOwnerAccount": "123456789012",
  "resourceType": "AWS::S3::Bucket",
  "status": "ACTIVE"
}

# Terraform
resource "aws_accessanalyzer_analyzer" "org_analyzer" {
  analyzer_name = "organization-analyzer"
  type          = "ORGANIZATION"  # or "ACCOUNT"
  
  tags = {
    Environment = "production"
  }
}

# Archive rule (suppress known findings)
resource "aws_accessanalyzer_archive_rule" "known_partner" {
  analyzer_name = aws_accessanalyzer_analyzer.org_analyzer.analyzer_name
  rule_name     = "known-partner-account"

  filter {
    criteria = "principal.AWS"
    eq       = ["arn:aws:iam::${var.partner_account_id}:root"]
  }

  filter {
    criteria = "resourceType"
    eq       = ["AWS::S3::Bucket"]
  }
}
```

### 3.5 AWS Config Rules

```
AWS Config Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                           AWS Config                                 │
│                                                                      │
│  ┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐ │
│  │   Configuration │    │   Configuration │    │    Rules        │ │
│  │   Recorder      │───►│    History      │───►│    Evaluation   │ │
│  │                 │    │    Timeline     │    │                 │ │
│  └─────────────────┘    └─────────────────┘    └─────────────────┘ │
│                                                        │            │
│                                                        ▼            │
│                                              ┌─────────────────┐   │
│                                              │   Compliance    │   │
│                                              │   Dashboard     │   │
│                                              └─────────────────┘   │
│                                                        │            │
│                                                        ▼            │
│                                              ┌─────────────────┐   │
│                                              │  Auto-Remediate │   │
│                                              │  (SSM Documents)│   │
│                                              └─────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘

Essential Config Rules for Platform Engineers:
──────────────────────────────────────────────────────────────────────

# Terraform - Common Security Rules

resource "aws_config_config_rule" "s3_bucket_public_read_prohibited" {
  name = "s3-bucket-public-read-prohibited"

  source {
    owner             = "AWS"
    source_identifier = "S3_BUCKET_PUBLIC_READ_PROHIBITED"
  }
}

resource "aws_config_config_rule" "ec2_instance_no_public_ip" {
  name = "ec2-instance-no-public-ip"

  source {
    owner             = "AWS"
    source_identifier = "EC2_INSTANCE_NO_PUBLIC_IP"
  }

  scope {
    compliance_resource_types = ["AWS::EC2::Instance"]
  }
}

resource "aws_config_config_rule" "rds_storage_encrypted" {
  name = "rds-storage-encrypted"

  source {
    owner             = "AWS"
    source_identifier = "RDS_STORAGE_ENCRYPTED"
  }
}

resource "aws_config_config_rule" "iam_password_policy" {
  name = "iam-password-policy"

  source {
    owner             = "AWS"
    source_identifier = "IAM_PASSWORD_POLICY"
  }

  input_parameters = jsonencode({
    RequireUppercaseCharacters = "true"
    RequireLowercaseCharacters = "true"
    RequireSymbols             = "true"
    RequireNumbers             = "true"
    MinimumPasswordLength      = "14"
    PasswordReusePrevention    = "24"
    MaxPasswordAge             = "90"
  })
}

resource "aws_config_config_rule" "encrypted_volumes" {
  name = "encrypted-volumes"

  source {
    owner             = "AWS"
    source_identifier = "ENCRYPTED_VOLUMES"
  }
}

resource "aws_config_config_rule" "vpc_flow_logs_enabled" {
  name = "vpc-flow-logs-enabled"

  source {
    owner             = "AWS"
    source_identifier = "VPC_FLOW_LOGS_ENABLED"
  }
}

# Auto-Remediation
resource "aws_config_remediation_configuration" "s3_remediation" {
  config_rule_name = aws_config_config_rule.s3_bucket_public_read_prohibited.name

  resource_type    = "AWS::S3::Bucket"
  target_type      = "SSM_DOCUMENT"
  target_id        = "AWS-DisableS3BucketPublicReadWrite"

  parameter {
    name         = "S3BucketName"
    resource_value = "RESOURCE_ID"
  }

  automatic                  = true
  maximum_automatic_attempts = 5
  retry_attempt_seconds      = 60
}
```

### 3.6 Amazon Detective

```
Amazon Detective Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                        Amazon Detective                              │
│                                                                      │
│  Data Sources:                    Analysis:                          │
│  ┌─────────────────┐             ┌─────────────────────────────────┐│
│  │  VPC Flow Logs  │────────────►│                                 ││
│  └─────────────────┘             │   Behavior Graph                ││
│                                  │                                 ││
│  ┌─────────────────┐             │   • Entity relationships        ││
│  │  CloudTrail     │────────────►│   • Activity timelines          ││
│  └─────────────────┘             │   • Unusual patterns            ││
│                                  │                                 ││
│  ┌─────────────────┐             │   Machine Learning:             ││
│  │  GuardDuty      │────────────►│   • Anomaly detection           ││
│  │  Findings       │             │   • Baseline establishment      ││
│  └─────────────────┘             │                                 ││
│                                  └─────────────────────────────────┘│
│  ┌─────────────────┐                                                │
│  │  EKS Audit Logs │─────────────────────────────────────────────►  │
│  └─────────────────┘                                                │
│                                                                      │
│  Investigation Workflow:                                             │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐     │
│  │ Finding  │───►│ Entity   │───►│ Timeline │───►│ Root     │     │
│  │ Triage   │    │ Profile  │    │ Analysis │    │ Cause    │     │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Use Case: GuardDuty Finding Investigation
──────────────────────────────────────────────────────────────────────

1. GuardDuty detects: UnauthorizedAccess:IAMUser/InstanceCredentialExfiltration

2. Open in Detective:
   • See EC2 instance activity timeline
   • Track API calls made with instance credentials
   • Identify external IP addresses involved
   • Visualize relationships between entities

3. Investigation questions Detective answers:
   • What other API calls were made?
   • What resources were accessed?
   • Is this normal behavior for this entity?
   • What was the sequence of events?

# Enable Detective
resource "aws_detective_graph" "main" {
  tags = {
    Environment = "production"
  }
}

# Multi-account invitation
resource "aws_detective_invitation_accepter" "member" {
  provider  = aws.member
  graph_arn = aws_detective_graph.main.graph_arn
}
```

### 3.7 Security Architecture Pattern

```
Enterprise Security Architecture
──────────────────────────────────────────────────────────────────────

                    ┌─────────────────────────────────┐
                    │        Security Account          │
                    │                                  │
                    │  ┌─────────────────────────────┐│
                    │  │       Security Hub          ││
                    │  │  (Delegated Admin)          ││
                    │  │                             ││
                    │  │  Aggregates all findings:   ││
                    │  │  • GuardDuty                ││
                    │  │  • Inspector                ││
                    │  │  • Macie                    ││
                    │  │  • Access Analyzer          ││
                    │  │  • Config                   ││
                    │  └─────────────────────────────┘│
                    │                                  │
                    │  ┌─────────────────────────────┐│
                    │  │    Detective                ││
                    │  │  (Investigation hub)        ││
                    │  └─────────────────────────────┘│
                    │                                  │
                    │  ┌─────────────────────────────┐│
                    │  │  SIEM Integration           ││
                    │  │  (Splunk, Datadog, etc.)    ││
                    │  └─────────────────────────────┘│
                    │                                  │
                    └───────────────┬─────────────────┘
                                    │
          ┌─────────────────────────┼─────────────────────────┐
          │                         │                         │
          ▼                         ▼                         ▼
┌─────────────────┐      ┌─────────────────┐      ┌─────────────────┐
│   Workload      │      │   Workload      │      │   Workload      │
│   Account 1     │      │   Account 2     │      │   Account N     │
│                 │      │                 │      │                 │
│ ┌─────────────┐ │      │ ┌─────────────┐ │      │ ┌─────────────┐ │
│ │ GuardDuty   │ │      │ │ GuardDuty   │ │      │ │ GuardDuty   │ │
│ │ (Member)    │ │      │ │ (Member)    │ │      │ │ (Member)    │ │
│ └─────────────┘ │      │ └─────────────┘ │      │ └─────────────┘ │
│ ┌─────────────┐ │      │ ┌─────────────┐ │      │ ┌─────────────┐ │
│ │ Config      │ │      │ │ Config      │ │      │ │ Config      │ │
│ │ (Member)    │ │      │ │ (Member)    │ │      │ │ (Member)    │ │
│ └─────────────┘ │      │ └─────────────┘ │      │ └─────────────┘ │
│                 │      │                 │      │                 │
└─────────────────┘      └─────────────────┘      └─────────────────┘

Security Automation Flow:
────────────────────────────────────────────────────────────────

Finding       EventBridge      Lambda           Action
Detected  ───► Rule        ───► Function   ───► Remediate
                                                  │
                                                  ▼
                                           ┌──────────────┐
                                           │ • Block IP   │
                                           │ • Isolate EC2│
                                           │ • Revoke IAM │
                                           │ • Notify     │
                                           └──────────────┘
```

---

## 4. Cost Optimization

### 4.1 Cost Management Services Overview

```
AWS Cost Management Stack
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Cost Intelligence Dashboard                       │
│           (QuickSight dashboards from CUR data)                     │
└─────────────────────────────────────────────────────────────────────┘
                                ▲
                                │
┌─────────────────────────────────────────────────────────────────────┐
│                Cost & Usage Report (CUR)                             │
│         Most detailed billing data (hourly, resource-level)         │
│                    → S3 → Athena/QuickSight                         │
└─────────────────────────────────────────────────────────────────────┘
                                ▲
          ┌─────────────────────┼─────────────────────┐
          │                     │                     │
┌─────────┴────────┐  ┌────────┴────────┐  ┌────────┴────────┐
│   Cost Explorer  │  │   AWS Budgets   │  │    Compute      │
│                  │  │                 │  │    Optimizer    │
│ • Visualize costs│  │ • Set budgets   │  │                 │
│ • Filter/group   │  │ • Alert on      │  │ • Right-sizing  │
│ • Forecast       │  │   thresholds    │  │ • EC2, Lambda   │
│ • Recommendations│  │ • Take actions  │  │ • EBS, ECS      │
└──────────────────┘  └─────────────────┘  └─────────────────┘

Pricing Models:
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                                                                      │
│   On-Demand ◄────────────────────────────────────────► Spot         │
│   (Full price)                                        (Up to 90% off)│
│        │                                                   │         │
│        │    ┌─────────────────────────────────┐           │         │
│        │    │                                 │           │         │
│        ├───►│   Reserved Instances            │◄──────────┤         │
│        │    │   • 1 or 3 year commitment      │           │         │
│        │    │   • Up to 72% savings           │           │         │
│        │    │   • Standard or Convertible     │           │         │
│        │    │                                 │           │         │
│        │    └─────────────────────────────────┘           │         │
│        │                                                   │         │
│        │    ┌─────────────────────────────────┐           │         │
│        │    │                                 │           │         │
│        └───►│   Savings Plans                 │◄──────────┘         │
│             │   • Compute SP (flexible)       │                      │
│             │   • EC2 Instance SP (specific)  │                      │
│             │   • 1 or 3 year, up to 66%     │                      │
│             │                                 │                      │
│             └─────────────────────────────────┘                      │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

### 4.2 Reserved Instances vs Savings Plans

```
Comparison: Reserved Instances vs Savings Plans
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│ Attribute          │ Reserved Instances  │ Savings Plans           │
├────────────────────┼─────────────────────┼─────────────────────────┤
│ Commitment Type    │ Capacity (instances)│ $/hour expenditure      │
│                    │                     │                         │
│ Flexibility        │ Standard: Low       │ Compute SP: High        │
│                    │ Convertible: Medium │ EC2 SP: Medium          │
│                    │                     │                         │
│ Applies to         │ EC2 only           │ EC2, Fargate, Lambda    │
│                    │                     │ SageMaker               │
│                    │                     │                         │
│ Regional/Zonal     │ Both options        │ Regional only           │
│                    │                     │                         │
│ Instance Size      │ Flexible within     │ Fully flexible          │
│ Flexibility        │ family (Regional)   │ (Compute SP)            │
│                    │                     │                         │
│ Capacity Reserv.   │ Yes (Zonal)        │ No                      │
│                    │                     │                         │
│ Max Savings        │ 72%                │ 66%                     │
│                    │                     │                         │
│ Recommendation     │ Stable, predictable │ Variable workloads      │
│                    │ EC2-only workloads  │ Multi-service usage     │
└─────────────────────┴─────────────────────┴─────────────────────────┘

Savings Plan Types:
──────────────────────────────────────────────────────────────────────

1. Compute Savings Plan (Most Flexible)
   • Applies to: EC2, Fargate, Lambda
   • Change: Instance family, size, OS, tenancy, region
   • Best for: Variable workloads, multi-region, containers

2. EC2 Instance Savings Plan (Better Discount)
   • Applies to: EC2 only
   • Locked: Instance family, region
   • Change: Size, OS, tenancy
   • Best for: Stable EC2 workloads in specific region

3. SageMaker Savings Plan
   • Applies to: SageMaker only
   • Best for: ML workloads

# Purchase via CLI
aws savingsplans create-savings-plan \
  --savings-plan-offering-id "offering-123" \
  --commitment "10.00" \
  --savings-plan-type "ComputeSavingsPlans"
```

### 4.3 Spot Instances Strategies

```
Spot Instance Strategies
──────────────────────────────────────────────────────────────────────

Strategy 1: Diversification (Recommended)
─────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    ASG with Mixed Instances                          │
│                                                                      │
│  Allocation Strategy: capacity-optimized-prioritized                │
│                                                                      │
│  Instance Types:                                                     │
│  ┌─────────┬─────────┬─────────┬─────────┬─────────┐               │
│  │m5.large │m5.xlarge│m5a.large│m4.large │c5.large │               │
│  └─────────┴─────────┴─────────┴─────────┴─────────┘               │
│                                                                      │
│  Multiple AZs: us-east-1a, us-east-1b, us-east-1c                 │
│                                                                      │
│  Result: Up to 15 capacity pools                                   │
│          Lower interruption rate                                    │
└─────────────────────────────────────────────────────────────────────┘

# Terraform - ASG Mixed Instances
resource "aws_autoscaling_group" "spot_diverse" {
  name                = "spot-diverse-asg"
  vpc_zone_identifier = var.subnet_ids
  min_size            = 2
  max_size            = 20
  desired_capacity    = 10

  mixed_instances_policy {
    instances_distribution {
      on_demand_base_capacity                  = 2  # Baseline on-demand
      on_demand_percentage_above_base_capacity = 20 # 20% on-demand above base
      spot_allocation_strategy                 = "capacity-optimized-prioritized"
      spot_instance_pools                      = 0  # Use all pools
    }

    launch_template {
      launch_template_specification {
        launch_template_id = aws_launch_template.main.id
        version            = "$Latest"
      }

      override {
        instance_type     = "m5.large"
        weighted_capacity = "1"
      }
      override {
        instance_type     = "m5.xlarge"
        weighted_capacity = "2"
      }
      override {
        instance_type     = "m5a.large"
        weighted_capacity = "1"
      }
      override {
        instance_type     = "m4.large"
        weighted_capacity = "1"
      }
      override {
        instance_type     = "c5.large"
        weighted_capacity = "1"
      }
    }
  }

  tag {
    key                 = "Environment"
    value               = "production"
    propagate_at_launch = true
  }
}

Strategy 2: Spot Fleet for Batch Processing
────────────────────────────────────────────

resource "aws_spot_fleet_request" "batch_processing" {
  iam_fleet_role                      = aws_iam_role.spot_fleet.arn
  target_capacity                     = 50
  allocation_strategy                 = "capacityOptimized"
  terminate_instances_with_expiration = true
  valid_until                        = timeadd(timestamp(), "6h")

  launch_template_config {
    launch_template_specification {
      id      = aws_launch_template.batch.id
      version = "$Latest"
    }

    overrides {
      instance_type     = "c5.2xlarge"
      availability_zone = "us-east-1a"
    }
    overrides {
      instance_type     = "c5.2xlarge"
      availability_zone = "us-east-1b"
    }
    overrides {
      instance_type     = "c5a.2xlarge"
      availability_zone = "us-east-1a"
    }
  }
}

Spot Interruption Handling:
────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                  Spot Interruption Flow                              │
│                                                                      │
│  2-min Warning ──► EventBridge ──► Lambda ──► Actions               │
│                                                                      │
│  Actions:                                                            │
│  • Drain connections from ELB                                        │
│  • Checkpoint processing state                                       │
│  • Push work to SQS                                                  │
│  • Deregister from service discovery                                │
└─────────────────────────────────────────────────────────────────────┘

# EventBridge rule for spot interruption
resource "aws_cloudwatch_event_rule" "spot_interruption" {
  name = "spot-interruption-warning"

  event_pattern = jsonencode({
    source      = ["aws.ec2"]
    detail-type = ["EC2 Spot Instance Interruption Warning"]
  })
}

resource "aws_cloudwatch_event_target" "drain_handler" {
  rule      = aws_cloudwatch_event_rule.spot_interruption.name
  target_id = "spot-drain"
  arn       = aws_lambda_function.spot_handler.arn
}
```

### 4.4 Right-Sizing with Compute Optimizer

```
AWS Compute Optimizer
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Compute Optimizer Flow                            │
│                                                                      │
│  CloudWatch Metrics ──► ML Analysis ──► Recommendations             │
│  (14 days min)                                                      │
│                                                                      │
│  Analyzes:                                                          │
│  ┌──────────────┬────────────────────────────────────────────────┐ │
│  │ EC2          │ CPU, Memory, Network, Storage                   │ │
│  ├──────────────┼────────────────────────────────────────────────┤ │
│  │ EBS          │ IOPS, Throughput                               │ │
│  ├──────────────┼────────────────────────────────────────────────┤ │
│  │ Lambda       │ Memory, Duration                               │ │
│  ├──────────────┼────────────────────────────────────────────────┤ │
│  │ ECS on       │ CPU, Memory                                    │ │
│  │ Fargate      │                                                │ │
│  ├──────────────┼────────────────────────────────────────────────┤ │
│  │ Auto Scaling │ Group configuration                            │ │
│  │ Groups       │                                                │ │
│  └──────────────┴────────────────────────────────────────────────┘ │
│                                                                      │
│  Recommendation Types:                                               │
│  • Under-provisioned (upgrade for performance)                      │
│  • Over-provisioned (downgrade for savings)                         │
│  • Optimized (no change needed)                                     │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Example Recommendation:
──────────────────────────────────────────────────────────────────────

Current Instance: m5.2xlarge ($0.384/hr)
  CPU Utilization: 12% avg, 25% max
  Memory: 18% avg
  Network: 2% of capacity

Recommended: m5.large ($0.096/hr)
  Projected CPU: 48% avg
  Monthly Savings: ~$207
  Risk: Low

# Enable Compute Optimizer via CLI
aws compute-optimizer update-enrollment-status \
  --status Active \
  --include-member-accounts
```

### 4.5 AWS Budgets & Alerts

```
Budget Configuration
──────────────────────────────────────────────────────────────────────

Budget Types:
┌─────────────────────────────────────────────────────────────────────┐
│ Type              │ Use Case                        │ Alerts        │
├───────────────────┼─────────────────────────────────┼───────────────┤
│ Cost Budget       │ Track spending against limit    │ 50%, 80%, 100%│
│ Usage Budget      │ Track resource usage            │ Thresholds    │
│ RI Utilization    │ Track RI coverage               │ Below target  │
│ RI Coverage       │ Track RI usage %                │ Below target  │
│ Savings Plans     │ Track SP utilization            │ Below target  │
│ Utilization       │                                 │               │
└───────────────────┴─────────────────────────────────┴───────────────┘

# Terraform - Cost Budget with Alerts
resource "aws_budgets_budget" "monthly_cost" {
  name         = "monthly-cost-budget"
  budget_type  = "COST"
  limit_amount = "10000"
  limit_unit   = "USD"
  time_unit    = "MONTHLY"

  cost_filter {
    name   = "TagKeyValue"
    values = ["user:Environment$production"]
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = 80
    threshold_type            = "PERCENTAGE"
    notification_type         = "FORECASTED"
    subscriber_email_addresses = ["platform-team@company.com"]
  }

  notification {
    comparison_operator        = "GREATER_THAN"
    threshold                  = 100
    threshold_type            = "PERCENTAGE"
    notification_type         = "ACTUAL"
    subscriber_email_addresses = ["platform-team@company.com"]
    subscriber_sns_topic_arns = [aws_sns_topic.budget_alerts.arn]
  }
}

# Budget Action - Automated Response
resource "aws_budgets_budget_action" "stop_instances" {
  budget_name        = aws_budgets_budget.monthly_cost.name
  action_type        = "RUN_SSM_DOCUMENTS"
  approval_model     = "AUTOMATIC"
  notification_type  = "ACTUAL"
  action_threshold {
    action_threshold_type  = "PERCENTAGE"
    action_threshold_value = 120
  }

  definition {
    ssm_action_definition {
      action_sub_type = "STOP_EC2_INSTANCES"
      region          = var.region
      instance_ids    = var.non_critical_instance_ids
    }
  }

  execution_role_arn = aws_iam_role.budget_action.arn
}
```

### 4.6 FinOps Practices

```
FinOps Framework for Platform Engineers
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                         FinOps Lifecycle                             │
│                                                                      │
│     Inform ──────────► Optimize ──────────► Operate                 │
│        │                   │                    │                    │
│  • Visibility         • Right-sizing       • Governance            │
│  • Allocation         • Commitments        • Automation            │
│  • Benchmarking       • Spot usage         • Continuous            │
│  • Forecasting        • Waste elimination    improvement           │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Cost Allocation Strategy:
──────────────────────────────────────────────────────────────────────

1. Tagging Strategy
┌─────────────────────────────────────────────────────────────────────┐
│ Tag Key          │ Example Value       │ Purpose                    │
├──────────────────┼─────────────────────┼────────────────────────────┤
│ Environment      │ prod, staging, dev  │ Environment costs          │
│ Team             │ platform, data, ml  │ Team chargeback            │
│ Project          │ checkout, search    │ Project allocation         │
│ CostCenter       │ CC-12345           │ Financial accounting        │
│ Owner            │ team@company.com    │ Accountability             │
│ Application      │ api-gateway         │ Application costs          │
└──────────────────┴─────────────────────┴────────────────────────────┘

# AWS Tag Policy (Organizations)
{
  "tags": {
    "Environment": {
      "tag_key": { "@@assign": "Environment" },
      "tag_value": {
        "@@assign": ["prod", "staging", "dev", "sandbox"]
      },
      "enforced_for": {
        "@@assign": [
          "ec2:instance",
          "rds:db",
          "s3:bucket"
        ]
      }
    },
    "Team": {
      "tag_key": { "@@assign": "Team" },
      "enforced_for": {
        "@@assign": ["ec2:instance", "rds:db"]
      }
    }
  }
}

2. Showback/Chargeback Dashboard

cost_by_team = """
SELECT 
  line_item_usage_account_id,
  resource_tags_user_team AS team,
  SUM(line_item_blended_cost) AS cost
FROM cur_database.cur_table
WHERE month = '2024-01'
GROUP BY 1, 2
ORDER BY 3 DESC
"""

3. Cost Anomaly Detection
resource "aws_ce_anomaly_monitor" "service_monitor" {
  name              = "service-anomaly-monitor"
  monitor_type      = "DIMENSIONAL"
  monitor_dimension = "SERVICE"
}

resource "aws_ce_anomaly_subscription" "alert" {
  name      = "cost-anomaly-alert"
  frequency = "IMMEDIATE"

  monitor_arn_list = [
    aws_ce_anomaly_monitor.service_monitor.arn
  ]

  subscriber {
    type    = "EMAIL"
    address = "finops@company.com"
  }

  threshold_expression {
    dimension {
      key           = "ANOMALY_TOTAL_IMPACT_PERCENTAGE"
      match_options = ["GREATER_THAN_OR_EQUAL"]
      values        = ["10"]  # 10% increase
    }
  }
}

Cost Optimization Checklist:
──────────────────────────────────────────────────────────────────────

□ Quick Wins (Week 1)
  ├── Delete unused EBS volumes
  ├── Release unattached Elastic IPs
  ├── Delete old snapshots
  ├── Stop/terminate unused instances
  └── Delete unused load balancers

□ Medium Term (Month 1)
  ├── Implement auto-scaling
  ├── Use Spot for non-critical workloads
  ├── Right-size based on Compute Optimizer
  ├── Review and optimize data transfer
  └── Implement S3 lifecycle policies

□ Long Term (Quarter 1)
  ├── Purchase Savings Plans
  ├── Implement FinOps processes
  ├── Set up showback/chargeback
  ├── Architect for cost efficiency
  └── Regular cost review meetings
```

---

## 5. DevOps & CI/CD on AWS

### 5.1 AWS CI/CD Services Overview

```
AWS CI/CD Ecosystem
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                        CodePipeline                                  │
│              (Orchestrates the entire CI/CD workflow)               │
│                                                                      │
│  Source      ─►  Build      ─►  Test      ─►  Deploy               │
│  Stage          Stage          Stage          Stage                 │
└─────────────────────────────────────────────────────────────────────┘
      │               │              │              │
      ▼               ▼              ▼              ▼
┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────────────┐
│CodeCommit│   │CodeBuild │   │CodeBuild │   │    CodeDeploy    │
│  GitHub  │   │          │   │CodeBuild │   │    ECS/EKS       │
│GitLab    │   │ Maven    │   │with tests│   │    Elastic       │
│BitBucket │   │ Gradle   │   │          │   │    Beanstalk     │
│   S3     │   │ npm      │   │ Unit     │   │    Lambda        │
│          │   │ Docker   │   │ Integr.  │   │    CloudFormation│
└──────────┘   └──────────┘   └──────────┘   └──────────────────┘

Related Services:
┌──────────────────────────────────────────────────────────────────────┐
│ • ECR (Container Registry)  - Docker image storage                   │
│ • Secrets Manager           - Secure credential storage              │
│ • Parameter Store           - Configuration management               │
│ • S3                        - Artifact storage                       │
│ • CloudWatch                - Logging, metrics, alarms               │
└──────────────────────────────────────────────────────────────────────┘
```

### 5.2 CodePipeline Architecture

```
CodePipeline Deep Dive
──────────────────────────────────────────────────────────────────────

Pipeline Structure:
┌─────────────────────────────────────────────────────────────────────┐
│ Pipeline                                                             │
│  │                                                                   │
│  ├── Stage: Source                                                  │
│  │    └── Action: GitHub (Branch: main)                             │
│  │                                                                   │
│  ├── Stage: Build                                                   │
│  │    ├── Action: CodeBuild (Build App)        ─┐                  │
│  │    └── Action: CodeBuild (Build Container)  ─┤ Parallel         │
│  │                                               │                   │
│  ├── Stage: Test                                                    │
│  │    └── Action: CodeBuild (Integration Tests)                     │
│  │                                                                   │
│  ├── Stage: Deploy-Staging                                          │
│  │    └── Action: ECS Deploy (Staging Cluster)                      │
│  │                                                                   │
│  ├── Stage: Approval                                                │
│  │    └── Action: Manual Approval (SNS notification)                │
│  │                                                                   │
│  └── Stage: Deploy-Production                                       │
│       └── Action: ECS Deploy (Production Cluster)                   │
└─────────────────────────────────────────────────────────────────────┘

# Terraform - Complete Pipeline
resource "aws_codepipeline" "main" {
  name     = "app-pipeline"
  role_arn = aws_iam_role.codepipeline.arn

  artifact_store {
    location = aws_s3_bucket.artifacts.bucket
    type     = "S3"
    
    encryption_key {
      id   = aws_kms_key.artifacts.arn
      type = "KMS"
    }
  }

  stage {
    name = "Source"

    action {
      name             = "Source"
      category         = "Source"
      owner            = "AWS"
      provider         = "CodeStarSourceConnection"
      version          = "1"
      output_artifacts = ["source_output"]

      configuration = {
        ConnectionArn    = aws_codestarconnections_connection.github.arn
        FullRepositoryId = "org/repo"
        BranchName       = "main"
      }
    }
  }

  stage {
    name = "Build"

    action {
      name             = "Build"
      category         = "Build"
      owner            = "AWS"
      provider         = "CodeBuild"
      input_artifacts  = ["source_output"]
      output_artifacts = ["build_output"]
      version          = "1"

      configuration = {
        ProjectName = aws_codebuild_project.build.name
      }
    }
  }

  stage {
    name = "Deploy-Staging"

    action {
      name            = "Deploy"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "ECS"
      input_artifacts = ["build_output"]
      version         = "1"

      configuration = {
        ClusterName = aws_ecs_cluster.staging.name
        ServiceName = aws_ecs_service.staging.name
        FileName    = "imagedefinitions.json"
      }
    }
  }

  stage {
    name = "Approval"

    action {
      name     = "Approval"
      category = "Approval"
      owner    = "AWS"
      provider = "Manual"
      version  = "1"

      configuration = {
        NotificationArn = aws_sns_topic.approval.arn
        CustomData      = "Please review staging deployment and approve for production"
      }
    }
  }

  stage {
    name = "Deploy-Production"

    action {
      name            = "Deploy"
      category        = "Deploy"
      owner           = "AWS"
      provider        = "ECS"
      input_artifacts = ["build_output"]
      version         = "1"

      configuration = {
        ClusterName = aws_ecs_cluster.production.name
        ServiceName = aws_ecs_service.production.name
        FileName    = "imagedefinitions.json"
      }
    }
  }
}
```

### 5.3 CodeBuild Configuration

```
CodeBuild Project
──────────────────────────────────────────────────────────────────────

# Terraform
resource "aws_codebuild_project" "build" {
  name          = "app-build"
  description   = "Build and test application"
  service_role  = aws_iam_role.codebuild.arn
  build_timeout = 30

  source {
    type      = "CODEPIPELINE"
    buildspec = "buildspec.yml"
  }

  artifacts {
    type = "CODEPIPELINE"
  }

  environment {
    compute_type                = "BUILD_GENERAL1_MEDIUM"
    image                       = "aws/codebuild/amazonlinux2-x86_64-standard:4.0"
    type                        = "LINUX_CONTAINER"
    image_pull_credentials_type = "CODEBUILD"
    privileged_mode             = true  # Required for Docker builds

    environment_variable {
      name  = "AWS_ACCOUNT_ID"
      value = data.aws_caller_identity.current.account_id
    }

    environment_variable {
      name  = "ECR_REPO"
      value = aws_ecr_repository.app.repository_url
    }

    environment_variable {
      name  = "DB_PASSWORD"
      value = "/app/db-password"
      type  = "SECRETS_MANAGER"
    }
  }

  vpc_config {
    vpc_id             = aws_vpc.main.id
    subnets            = aws_subnet.private[*].id
    security_group_ids = [aws_security_group.codebuild.id]
  }

  cache {
    type     = "S3"
    location = "${aws_s3_bucket.cache.bucket}/build-cache"
  }

  logs_config {
    cloudwatch_logs {
      group_name  = aws_cloudwatch_log_group.codebuild.name
      stream_name = "build"
    }

    s3_logs {
      status   = "ENABLED"
      location = "${aws_s3_bucket.logs.bucket}/codebuild-logs"
    }
  }
}

buildspec.yml:
──────────────────────────────────────────────────────────────────────

version: 0.2

env:
  variables:
    APP_ENV: "production"
  secrets-manager:
    DB_PASSWORD: "prod/app/db:password"
  exported-variables:
    - IMAGE_TAG

phases:
  install:
    runtime-versions:
      java: corretto17
      docker: 20
    commands:
      - echo "Installing dependencies..."

  pre_build:
    commands:
      - echo "Logging into ECR..."
      - aws ecr get-login-password --region $AWS_DEFAULT_REGION | docker login --username AWS --password-stdin $ECR_REPO
      - IMAGE_TAG=$(echo $CODEBUILD_RESOLVED_SOURCE_VERSION | cut -c 1-7)

  build:
    commands:
      - echo "Building application..."
      - ./gradlew clean build -x test
      - echo "Running tests..."
      - ./gradlew test
      - echo "Building Docker image..."
      - docker build -t $ECR_REPO:$IMAGE_TAG .
      - docker tag $ECR_REPO:$IMAGE_TAG $ECR_REPO:latest

  post_build:
    commands:
      - echo "Pushing Docker image..."
      - docker push $ECR_REPO:$IMAGE_TAG
      - docker push $ECR_REPO:latest
      - echo "Creating artifacts..."
      - printf '[{"name":"app","imageUri":"%s"}]' $ECR_REPO:$IMAGE_TAG > imagedefinitions.json

artifacts:
  files:
    - imagedefinitions.json
    - appspec.yml
    - taskdef.json
  discard-paths: yes

cache:
  paths:
    - '/root/.gradle/caches/**/*'
    - '/root/.gradle/wrapper/**/*'

reports:
  junit-reports:
    files:
      - 'build/test-results/test/*.xml'
    file-format: 'JUNITXML'
```

### 5.4 Deployment Strategies

```
Deployment Strategies Comparison
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│ Strategy           │ Risk   │ Rollback │ Duration │ Resource Cost  │
├────────────────────┼────────┼──────────┼──────────┼────────────────┤
│ In-Place           │ High   │ Redeploy │ Fast     │ Low            │
│ Rolling            │ Medium │ Redeploy │ Medium   │ Low            │
│ Blue/Green         │ Low    │ Instant  │ Medium   │ 2x during      │
│ Canary             │ Low    │ Instant  │ Slow     │ Low-Medium     │
│ A/B Testing        │ Low    │ Instant  │ Varies   │ Medium         │
└────────────────────┴────────┴──────────┴──────────┴────────────────┘

Strategy 1: Rolling Deployment (ECS)
──────────────────────────────────────────────────────────────────────

        ┌───────────────────────────────────────┐
        │         ECS Service                    │
        │    Desired: 4, Running: 4             │
        └───────────────────────────────────────┘
                        │
         ┌──────────────┴──────────────┐
         │                             │
         ▼                             ▼
  ┌──────────────┐            ┌──────────────┐
  │  Task v1     │            │  Task v1     │
  │  (running)   │            │  (running)   │
  └──────────────┘            └──────────────┘
  ┌──────────────┐            ┌──────────────┐
  │  Task v2     │            │  Task v1     │
  │  (starting)  │            │  (draining)  │
  └──────────────┘            └──────────────┘

# Terraform - ECS Service Rolling Update
resource "aws_ecs_service" "app" {
  name            = "app-service"
  cluster         = aws_ecs_cluster.main.id
  task_definition = aws_ecs_task_definition.app.arn
  desired_count   = 4

  deployment_configuration {
    maximum_percent         = 200  # Can scale to 8 tasks during deploy
    minimum_healthy_percent = 50   # Keep at least 2 healthy
  }

  deployment_circuit_breaker {
    enable   = true
    rollback = true  # Auto-rollback on failure
  }
}

Strategy 2: Blue/Green Deployment (CodeDeploy + ECS)
──────────────────────────────────────────────────────────────────────

                     ALB
                      │
         ┌───────────┴───────────┐
         │                       │
    ┌────▼────┐            ┌────▼────┐
    │ Target  │            │ Target  │
    │ Group   │            │ Group   │
    │ (Blue)  │            │ (Green) │
    │  100%   │            │   0%    │
    └────┬────┘            └────┬────┘
         │                       │
    ┌────▼────┐            ┌────▼────┐
    │ ECS     │            │ ECS     │
    │ Tasks   │            │ Tasks   │
    │ v1      │            │ v2      │
    └─────────┘            └─────────┘

    After shift:
    Blue: 0%  ←──────────→  Green: 100%

# appspec.yml for ECS Blue/Green
version: 0.0
Resources:
  - TargetService:
      Type: AWS::ECS::Service
      Properties:
        TaskDefinition: "<TASK_DEFINITION>"
        LoadBalancerInfo:
          ContainerName: "app"
          ContainerPort: 8080
        PlatformVersion: "LATEST"

Hooks:
  - BeforeInstall: "LambdaFunctionToValidateBeforeTrafficShift"
  - AfterInstall: "LambdaFunctionToValidateAfterTrafficShift"
  - AfterAllowTestTraffic: "LambdaFunctionToValidateTestTraffic"
  - BeforeAllowTraffic: "LambdaFunctionToValidateBeforeAllowingProdTraffic"
  - AfterAllowTraffic: "LambdaFunctionToValidateAfterAllowingProdTraffic"

# Terraform - CodeDeploy for ECS
resource "aws_codedeploy_deployment_group" "ecs" {
  app_name               = aws_codedeploy_app.ecs.name
  deployment_group_name  = "production"
  service_role_arn       = aws_iam_role.codedeploy.arn
  deployment_config_name = "CodeDeployDefault.ECSAllAtOnce"

  ecs_service {
    cluster_name = aws_ecs_cluster.main.name
    service_name = aws_ecs_service.app.name
  }

  blue_green_deployment_config {
    deployment_ready_option {
      action_on_timeout = "CONTINUE_DEPLOYMENT"
    }

    terminate_blue_instances_on_deployment_success {
      action                           = "TERMINATE"
      termination_wait_time_in_minutes = 5
    }
  }

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  load_balancer_info {
    target_group_pair_info {
      prod_traffic_route {
        listener_arns = [aws_lb_listener.https.arn]
      }

      target_group {
        name = aws_lb_target_group.blue.name
      }

      target_group {
        name = aws_lb_target_group.green.name
      }
    }
  }

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE", "DEPLOYMENT_STOP_ON_ALARM"]
  }

  alarm_configuration {
    alarms  = [aws_cloudwatch_metric_alarm.error_rate.name]
    enabled = true
  }
}

Strategy 3: Canary Deployment (Lambda)
──────────────────────────────────────────────────────────────────────

Traffic Shift:  10% ─────(10 min)────► 100%

               Incoming Traffic
                      │
                      ▼
              ┌───────────────┐
              │  Lambda Alias │
              │   "live"      │
              └───────┬───────┘
                      │
         ┌────────────┴────────────┐
         │                         │
    ┌────▼────┐              ┌────▼────┐
    │ Version │              │ Version │
    │    5    │              │    6    │
    │  (90%)  │              │  (10%)  │
    └─────────┘              └─────────┘

# Terraform - Lambda with Canary Deployment
resource "aws_lambda_alias" "live" {
  name             = "live"
  function_name    = aws_lambda_function.app.function_name
  function_version = aws_lambda_function.app.version

  routing_config {
    additional_version_weights = {
      (aws_lambda_function.app_new.version) = 0.1  # 10% to new version
    }
  }
}

resource "aws_codedeploy_deployment_group" "lambda" {
  app_name               = aws_codedeploy_app.lambda.name
  deployment_group_name  = "production"
  service_role_arn       = aws_iam_role.codedeploy.arn
  deployment_config_name = "CodeDeployDefault.LambdaCanary10Percent10Minutes"

  deployment_style {
    deployment_option = "WITH_TRAFFIC_CONTROL"
    deployment_type   = "BLUE_GREEN"
  }

  auto_rollback_configuration {
    enabled = true
    events  = ["DEPLOYMENT_FAILURE", "DEPLOYMENT_STOP_ON_ALARM"]
  }
}
```

### 5.5 GitOps with ArgoCD on EKS

```
GitOps Architecture
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                          Git Repository                              │
│                                                                      │
│  ├── apps/                                                          │
│  │   ├── app-a/                                                     │
│  │   │   ├── base/                                                  │
│  │   │   │   ├── deployment.yaml                                   │
│  │   │   │   ├── service.yaml                                      │
│  │   │   │   └── kustomization.yaml                                │
│  │   │   └── overlays/                                              │
│  │   │       ├── staging/                                           │
│  │   │       │   └── kustomization.yaml                            │
│  │   │       └── production/                                        │
│  │   │           └── kustomization.yaml                            │
│  │   └── app-b/                                                     │
│  └── argocd/                                                        │
│      └── applications.yaml                                          │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                │
                │ Watch
                ▼
        ┌───────────────┐
        │   ArgoCD      │
        │  Controller   │
        └───────┬───────┘
                │
                │ Sync
                ▼
        ┌───────────────┐
        │    EKS        │
        │   Cluster     │
        └───────────────┘

# ArgoCD Application
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/org/gitops-repo.git
    targetRevision: HEAD
    path: apps/my-app/overlays/production
  destination:
    server: https://kubernetes.default.svc
    namespace: my-app
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
    retry:
      limit: 5
      backoff:
        duration: 5s
        factor: 2
        maxDuration: 3m

# Terraform - Install ArgoCD on EKS
resource "helm_release" "argocd" {
  name             = "argocd"
  repository       = "https://argoproj.github.io/argo-helm"
  chart            = "argo-cd"
  namespace        = "argocd"
  create_namespace = true
  version          = "5.46.7"

  values = [
    yamlencode({
      server = {
        service = {
          type = "LoadBalancer"
          annotations = {
            "service.beta.kubernetes.io/aws-load-balancer-type" = "nlb"
          }
        }
      }
      configs = {
        repositories = {
          gitops-repo = {
            url  = "https://github.com/org/gitops-repo.git"
            type = "git"
          }
        }
      }
    })
  ]
}
```

---

## 6. Case Studies & Reference Architectures

### 6.1 Case Study: Multi-Region E-Commerce Platform

```
E-Commerce Platform Architecture
──────────────────────────────────────────────────────────────────────

Business Requirements:
• 99.99% availability SLA
• Handle 100K concurrent users
• Sub-100ms response time globally
• PCI-DSS compliance for payments
• Disaster recovery with <5 min RTO

┌─────────────────────────────────────────────────────────────────────┐
│                         Global Layer                                 │
│                                                                      │
│        Route 53 (Latency-based routing)                             │
│                    │                                                 │
│    ┌───────────────┴───────────────┐                                │
│    │                               │                                 │
│    ▼                               ▼                                 │
│ CloudFront              CloudFront Edge                             │
│ Distribution            Locations (~450)                            │
│    │                                                                 │
│    ├── S3 Static Assets (Images, JS, CSS)                          │
│    ├── Lambda@Edge (Auth, A/B Testing)                              │
│    └── API Gateway (Origin for APIs)                                │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
                              │
          ┌───────────────────┴───────────────────┐
          │                                       │
          ▼                                       ▼
┌─────────────────────────────────┐  ┌─────────────────────────────────┐
│     US-EAST-1 (Primary)         │  │     EU-WEST-1 (Primary)         │
│                                 │  │                                 │
│  ┌───────────────────────────┐  │  │  ┌───────────────────────────┐  │
│  │      API Gateway          │  │  │  │      API Gateway          │  │
│  │      (Regional)           │  │  │  │      (Regional)           │  │
│  └─────────────┬─────────────┘  │  │  └─────────────┬─────────────┘  │
│                │                │  │                │                │
│  ┌─────────────▼─────────────┐  │  │  ┌─────────────▼─────────────┐  │
│  │    Application Layer      │  │  │  │    Application Layer      │  │
│  │    ┌─────────────────┐    │  │  │  │    ┌─────────────────┐    │  │
│  │    │ EKS Cluster     │    │  │  │  │    │ EKS Cluster     │    │  │
│  │    │ ├── Catalog API │    │  │  │  │    │ ├── Catalog API │    │  │
│  │    │ ├── Cart API    │    │  │  │  │    │ ├── Cart API    │    │  │
│  │    │ ├── Order API   │    │  │  │  │    │ ├── Order API   │    │  │
│  │    │ └── Payment API │    │  │  │  │    │ └── Payment API │    │  │
│  │    └─────────────────┘    │  │  │  │    └─────────────────┘    │  │
│  └─────────────┬─────────────┘  │  │  └─────────────┬─────────────┘  │
│                │                │  │                │                │
│  ┌─────────────▼─────────────┐  │  │  ┌─────────────▼─────────────┐  │
│  │    Data Layer             │  │  │  │    Data Layer             │  │
│  │                           │  │  │  │                           │  │
│  │  ┌──────────────────────┐ │  │  │  │  ┌──────────────────────┐ │  │
│  │  │ Aurora Global        │◄┼──┼──┼──┼─►│ Aurora Global        │ │  │
│  │  │ (Primary Writer)     │ │  │  │  │  │ (Secondary Reader)   │ │  │
│  │  └──────────────────────┘ │  │  │  │  └──────────────────────┘ │  │
│  │                           │  │  │  │                           │  │
│  │  ┌──────────────────────┐ │  │  │  │  ┌──────────────────────┐ │  │
│  │  │ DynamoDB Global      │◄┼──┼──┼──┼─►│ DynamoDB Global      │ │  │
│  │  │ (Cart, Sessions)     │ │  │  │  │  │ (Cart, Sessions)     │ │  │
│  │  └──────────────────────┘ │  │  │  │  └──────────────────────┘ │  │
│  │                           │  │  │  │                           │  │
│  │  ┌──────────────────────┐ │  │  │  │  ┌──────────────────────┐ │  │
│  │  │ ElastiCache Global   │◄┼──┼──┼──┼─►│ ElastiCache Global   │ │  │
│  │  │ (Product Cache)      │ │  │  │  │  │ (Product Cache)      │ │  │
│  │  └──────────────────────┘ │  │  │  │  └──────────────────────┘ │  │
│  │                           │  │  │  │                           │  │
│  └───────────────────────────┘  │  │  └───────────────────────────┘  │
│                                 │  │                                 │
└─────────────────────────────────┘  └─────────────────────────────────┘

Key Design Decisions:
──────────────────────────────────────────────────────────────────────

1. Database Selection:
   • Aurora Global - Product catalog, Orders (ACID transactions)
   • DynamoDB Global - Cart, Sessions (low latency, auto-scale)
   • ElastiCache Global - Product cache (sub-ms reads)

2. Multi-Region Strategy:
   • Active-Active for reads (both regions serve traffic)
   • Active-Passive for writes (us-east-1 is primary writer)
   • Route 53 latency-based routing

3. PCI Compliance:
   • Payment API in isolated PCI subnet
   • AWS Payment Cryptography for tokenization
   • Secrets Manager for card data encryption keys
   • Dedicated Security Hub standards for PCI-DSS

4. Cost Optimization:
   • Reserved capacity for baseline (60% of traffic)
   • Spot instances for batch processing
   • CloudFront caching reduces origin requests by 70%
```

### 6.2 Case Study: Data Analytics Platform

```
Data Analytics Platform Architecture
──────────────────────────────────────────────────────────────────────

Business Requirements:
• Process 10TB of data daily
• Real-time dashboards with <5 second latency
• Machine learning predictions
• Cost-efficient storage (<$0.01/GB/month for cold data)

┌─────────────────────────────────────────────────────────────────────┐
│                        Data Ingestion Layer                          │
│                                                                      │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────┐               │
│  │ Kinesis     │   │ API Gateway │   │ Database    │               │
│  │ Data Streams│   │ (Events)    │   │ CDC (DMS)   │               │
│  └──────┬──────┘   └──────┬──────┘   └──────┬──────┘               │
│         │                 │                 │                        │
│         └─────────────────┴─────────────────┘                        │
│                           │                                          │
│                           ▼                                          │
│                  ┌─────────────────┐                                │
│                  │ Kinesis Data    │                                │
│                  │ Firehose        │                                │
│                  │ (Batching)      │                                │
│                  └────────┬────────┘                                │
│                           │                                          │
└───────────────────────────┼──────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Data Lake (S3)                                │
│                                                                      │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                     S3 Buckets                               │   │
│  │                                                              │   │
│  │  ┌─────────┐    ┌─────────┐    ┌─────────┐    ┌─────────┐  │   │
│  │  │ Raw     │───►│ Cleaned │───►│ Curated │───►│ Aggregated│ │   │
│  │  │ (Bronze)│    │ (Silver)│    │ (Gold)  │    │ (Platinum)│ │   │
│  │  │         │    │         │    │         │    │           │ │   │
│  │  │ Parquet │    │ Parquet │    │ Parquet │    │ Parquet   │ │   │
│  │  │ Part.   │    │ Part.   │    │ Part.   │    │ Part.     │ │   │
│  │  └─────────┘    └─────────┘    └─────────┘    └─────────── │   │
│  │                                                              │   │
│  │  Intelligent Tiering: Hot → Warm → Cold → Glacier           │   │
│  │                                                              │   │
│  └─────────────────────────────────────────────────────────────┘   │
│                                                                      │
└───────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Processing Layer                                 │
│                                                                      │
│  ┌──────────────────────┐    ┌──────────────────────┐              │
│  │    Batch Processing  │    │    Stream Processing │              │
│  │                      │    │                      │              │
│  │  ┌────────────────┐  │    │  ┌────────────────┐  │              │
│  │  │   AWS Glue     │  │    │  │   Kinesis      │  │              │
│  │  │   ETL Jobs     │  │    │  │   Analytics    │  │              │
│  │  └────────────────┘  │    │  └────────────────┘  │              │
│  │                      │    │                      │              │
│  │  ┌────────────────┐  │    │  ┌────────────────┐  │              │
│  │  │   EMR          │  │    │  │   Lambda       │  │              │
│  │  │   (Spark)      │  │    │  │   (Transforms) │  │              │
│  │  └────────────────┘  │    │  └────────────────┘  │              │
│  │                      │    │                      │              │
│  └──────────────────────┘    └──────────────────────┘              │
│                                                                      │
└───────────────────────────────────────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────────┐
│                     Analytics & ML Layer                             │
│                                                                      │
│  ┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐  │
│  │   Athena         │  │   Redshift       │  │   SageMaker      │  │
│  │                  │  │   Serverless     │  │                  │  │
│  │   SQL Queries    │  │   Data Warehouse │  │   ML Training &  │  │
│  │   on S3          │  │                  │  │   Inference      │  │
│  └────────┬─────────┘  └────────┬─────────┘  └────────┬─────────┘  │
│           │                     │                     │            │
│           └─────────────────────┴─────────────────────┘            │
│                                 │                                   │
│                                 ▼                                   │
│                        ┌─────────────────┐                         │
│                        │   QuickSight    │                         │
│                        │   Dashboards    │                         │
│                        └─────────────────┘                         │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Data Lake Organization (Medallion Architecture):
──────────────────────────────────────────────────────────────────────

s3://data-lake-bucket/
├── bronze/           # Raw data (immutable)
│   ├── source=orders/
│   │   └── year=2024/month=01/day=15/
│   │       └── data.parquet
│   └── source=clickstream/
│
├── silver/           # Cleaned & validated
│   ├── orders/
│   │   └── year=2024/month=01/
│   └── clickstream/
│
├── gold/             # Business-ready
│   ├── customer_360/
│   ├── product_performance/
│   └── revenue_metrics/
│
└── platinum/         # Aggregated for reporting
    └── daily_kpis/

Key Design Decisions:
──────────────────────────────────────────────────────────────────────

1. Storage Strategy:
   • S3 Intelligent-Tiering for automatic cost optimization
   • Parquet format for columnar analytics (10x compression)
   • Partition by date for efficient querying
   • Lifecycle policies: Hot (30 days) → Glacier (90 days)

2. Processing Choice:
   • Glue for ETL (serverless, pay-per-use)
   • EMR for complex Spark jobs (better for heavy processing)
   • Kinesis Analytics for real-time aggregations
   • Lambda for lightweight transformations

3. Query Engine Selection:
   • Athena: Ad-hoc queries, infrequent access ($5/TB scanned)
   • Redshift Serverless: Complex analytics, joins, aggregations
   • OpenSearch: Full-text search, log analytics

4. Cost Optimization:
   • Athena: Use partitions, columnar formats, compression
   • Redshift: Auto-pause, right-size, spectrum for cold data
   • EMR: Spot instances for task nodes (70% savings)
   • Data lifecycle: Auto-archive to Glacier after 90 days

5. Governance:
   • Lake Formation for fine-grained access control
   • Glue Data Catalog as central metadata store
   • AWS Macie for PII detection in S3
```

### 6.3 Case Study: SaaS Multi-Tenant Architecture

```
Multi-Tenant SaaS Architecture
──────────────────────────────────────────────────────────────────────

Tenancy Models Comparison:
┌─────────────────────────────────────────────────────────────────────┐
│ Model        │ Isolation │ Cost     │ Complexity │ Use Case        │
├──────────────┼───────────┼──────────┼────────────┼─────────────────┤
│ Silo         │ High      │ High     │ Low        │ Enterprise/     │
│ (Per tenant) │           │          │            │ Compliance      │
├──────────────┼───────────┼──────────┼────────────┼─────────────────┤
│ Pool         │ Low       │ Low      │ Medium     │ SMB / Freemium  │
│ (Shared)     │           │          │            │                 │
├──────────────┼───────────┼──────────┼────────────┼─────────────────┤
│ Bridge       │ Medium    │ Medium   │ High       │ Mixed customer  │
│ (Hybrid)     │           │          │            │ base            │
└──────────────┴───────────┴──────────┴────────────┴─────────────────┘

Hybrid Multi-Tenant Architecture:
──────────────────────────────────────────────────────────────────────

                    ┌─────────────────────────┐
                    │      API Gateway        │
                    │   (Tenant Routing)      │
                    └───────────┬─────────────┘
                                │
                    ┌───────────▼─────────────┐
                    │    Lambda Authorizer    │
                    │    (Tenant Context)     │
                    └───────────┬─────────────┘
                                │
         ┌──────────────────────┼──────────────────────┐
         │                      │                      │
         ▼                      ▼                      ▼
┌─────────────────┐   ┌─────────────────┐   ┌─────────────────┐
│ Premium Tier    │   │ Standard Tier   │   │ Free Tier       │
│ (Silo)          │   │ (Bridge)        │   │ (Pool)          │
│                 │   │                 │   │                 │
│ ┌─────────────┐ │   │ ┌─────────────┐ │   │ ┌─────────────┐ │
│ │ Dedicated   │ │   │ │ Shared EKS  │ │   │ │ Shared EKS  │ │
│ │ EKS Cluster │ │   │ │ (Namespace  │ │   │ │ (Shared     │ │
│ │             │ │   │ │  per tenant)│ │   │ │  Namespace) │ │
│ └─────────────┘ │   │ └─────────────┘ │   │ └─────────────┘ │
│                 │   │                 │   │                 │
│ ┌─────────────┐ │   │ ┌─────────────┐ │   │ ┌─────────────┐ │
│ │ Dedicated   │ │   │ │ Shared      │ │   │ │ Shared      │ │
│ │ Aurora DB   │ │   │ │ Aurora      │ │   │ │ Aurora      │ │
│ │             │ │   │ │ (Schema/    │ │   │ │ (Row-level) │ │
│ │             │ │   │ │  tenant)    │ │   │ │             │ │
│ └─────────────┘ │   │ └─────────────┘ │   │ └─────────────┘ │
│                 │   │                 │   │                 │
│ SLA: 99.99%     │   │ SLA: 99.9%      │   │ SLA: Best Effort│
│ Price: $$$$$    │   │ Price: $$$      │   │ Price: Free     │
└─────────────────┘   └─────────────────┘   └─────────────────┘

Tenant Isolation Patterns:
──────────────────────────────────────────────────────────────────────

1. Database Isolation (Row-Level)
   
   # All tenants in same table
   SELECT * FROM orders 
   WHERE tenant_id = :tenant_id  -- Always filtered
   
   # PostgreSQL Row-Level Security
   CREATE POLICY tenant_isolation ON orders
     USING (tenant_id = current_setting('app.tenant_id')::uuid);

2. Database Isolation (Schema-Level)
   
   # Per-tenant schema
   SET search_path TO tenant_acme;
   SELECT * FROM orders;

3. Compute Isolation (Kubernetes)
   
   # Namespace per tenant
   apiVersion: v1
   kind: Namespace
   metadata:
     name: tenant-acme
     labels:
       tenant: acme
   ---
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: deny-cross-tenant
     namespace: tenant-acme
   spec:
     podSelector: {}
     policyTypes:
       - Ingress
       - Egress
     ingress:
       - from:
         - namespaceSelector:
             matchLabels:
               tenant: acme
     egress:
       - to:
         - namespaceSelector:
             matchLabels:
               tenant: acme

4. IAM Isolation (Per Tenant Roles)
   
   {
     "Version": "2012-10-17",
     "Statement": [{
       "Effect": "Allow",
       "Action": ["s3:GetObject", "s3:PutObject"],
       "Resource": "arn:aws:s3:::saas-data/${aws:PrincipalTag/tenant_id}/*"
     }]
   }

Key Design Decisions:
──────────────────────────────────────────────────────────────────────

1. Tenant Onboarding:
   • Account Factory (Control Tower) for premium silo tenants
   • Automated namespace/schema creation for standard tenants
   • Self-service portal with API Gateway + Lambda

2. Tenant Context Propagation:
   • JWT tokens with tenant_id claim
   • Lambda authorizer extracts and validates
   • Context passed via headers through service mesh
   
   # Example tenant context middleware
   def extract_tenant(event):
       token = event['headers']['Authorization']
       claims = decode_jwt(token)
       return claims['tenant_id']

3. Metering & Billing:
   • CloudWatch custom metrics per tenant
   • API Gateway usage plans for rate limiting
   • AWS Marketplace integration for billing
   
   # CloudWatch metric for tenant usage
   cloudwatch.put_metric_data(
       Namespace='SaaS/Usage',
       MetricData=[{
           'MetricName': 'APIRequests',
           'Dimensions': [{'Name': 'TenantId', 'Value': tenant_id}],
           'Value': 1,
           'Unit': 'Count'
       }]
   )

4. Noisy Neighbor Prevention:
   • API Gateway throttling per tenant
   • Kubernetes ResourceQuotas per namespace
   • Database connection pooling limits
   • SQS per-tenant queues for isolation

5. Tenant Data Backup & Recovery:
   • Automated backups tagged with tenant_id
   • Point-in-time recovery per tenant (Aurora)
   • Cross-region replication for premium tenants
```

### 6.4 Case Study: Event-Driven Microservices

```
Event-Driven Architecture
──────────────────────────────────────────────────────────────────────

Business Requirements:
• Loose coupling between services
• Async processing for scalability
• Exactly-once delivery semantics
• Event replay for debugging

┌─────────────────────────────────────────────────────────────────────┐
│                      Event-Driven Architecture                       │
│                                                                      │
│   ┌─────────────────────────────────────────────────────────────┐  │
│   │                    EventBridge (Event Bus)                     │  │
│   │                                                                │  │
│   │   Rules:                                                       │  │
│   │   • order.created → OrderProcessor                            │  │
│   │   • order.completed → NotificationService                     │  │
│   │   • payment.processed → OrderService                          │  │
│   │   • inventory.reserved → ShippingService                      │  │
│   └─────────────────────────────────────────────────────────────┘  │
│                              │                                      │
│              ┌───────────────┼───────────────┐                     │
│              │               │               │                      │
│              ▼               ▼               ▼                      │
│   ┌──────────────┐ ┌──────────────┐ ┌──────────────┐              │
│   │   SQS Queue  │ │   SQS Queue  │ │   SQS Queue  │              │
│   │   (Orders)   │ │  (Payments)  │ │ (Inventory)  │              │
│   └──────┬───────┘ └──────┬───────┘ └──────┬───────┘              │
│          │                │                │                        │
│          ▼                ▼                ▼                        │
│   ┌──────────────┐ ┌──────────────┐ ┌──────────────┐              │
│   │   Order      │ │   Payment    │ │  Inventory   │              │
│   │   Service    │ │   Service    │ │   Service    │              │
│   │   (ECS)      │ │   (Lambda)   │ │   (ECS)      │              │
│   └──────┬───────┘ └──────┬───────┘ └──────┬───────┘              │
│          │                │                │                        │
│          ▼                ▼                ▼                        │
│   ┌──────────────┐ ┌──────────────┐ ┌──────────────┐              │
│   │   DynamoDB   │ │   Aurora     │ │   DynamoDB   │              │
│   │   (Orders)   │ │  (Payments)  │ │  (Inventory) │              │
│   └──────────────┘ └──────────────┘ └──────────────┘              │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘

Event Flow: Order Creation
──────────────────────────────────────────────────────────────────────

1. Order API → OrderCreated event
2. EventBridge → routes to SQS queues
3. Parallel Processing:
   • Payment Service: Reserve funds
   • Inventory Service: Check & reserve stock
4. Both complete → OrderConfirmed event
5. Notification Service: Send confirmation email

# EventBridge Rule
resource "aws_cloudwatch_event_rule" "order_created" {
  name        = "order-created-rule"
  event_bus_name = aws_cloudwatch_event_bus.main.name

  event_pattern = jsonencode({
    source      = ["order-service"]
    detail-type = ["OrderCreated"]
  })
}

resource "aws_cloudwatch_event_target" "payment_queue" {
  rule           = aws_cloudwatch_event_rule.order_created.name
  event_bus_name = aws_cloudwatch_event_bus.main.name
  target_id      = "payment-queue"
  arn            = aws_sqs_queue.payment.arn

  dead_letter_config {
    arn = aws_sqs_queue.dlq.arn
  }

  retry_policy {
    maximum_event_age_in_seconds = 86400
    maximum_retry_attempts       = 3
  }
}

# Event Archive for Replay
resource "aws_cloudwatch_event_archive" "order_events" {
  name             = "order-events-archive"
  event_source_arn = aws_cloudwatch_event_bus.main.arn
  retention_days   = 365

  event_pattern = jsonencode({
    source = ["order-service"]
  })
}

Saga Pattern Implementation:
──────────────────────────────────────────────────────────────────────

┌─────────────────────────────────────────────────────────────────────┐
│                    Step Functions Saga                               │
│                                                                      │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐                   │
│   │ Reserve  │────►│ Process  │────►│  Ship    │                   │
│   │ Inventory│     │ Payment  │     │  Order   │                   │
│   └────┬─────┘     └────┬─────┘     └────┬─────┘                   │
│        │                │                │                          │
│        │ Failed?        │ Failed?        │ Failed?                  │
│        ▼                ▼                ▼                          │
│   ┌──────────┐     ┌──────────┐     ┌──────────┐                   │
│   │ Release  │◄────│ Refund   │◄────│ Cancel   │                   │
│   │ Inventory│     │ Payment  │     │ Shipment │                   │
│   └──────────┘     └──────────┘     └──────────┘                   │
│                                                                      │
│   Compensating Transactions (Rollback on failure)                   │
└─────────────────────────────────────────────────────────────────────┘

# Step Functions State Machine (ASL)
{
  "Comment": "Order Processing Saga",
  "StartAt": "ReserveInventory",
  "States": {
    "ReserveInventory": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:reserve-inventory",
      "ResultPath": "$.inventoryResult",
      "Catch": [{
        "ErrorEquals": ["States.ALL"],
        "ResultPath": "$.error",
        "Next": "OrderFailed"
      }],
      "Next": "ProcessPayment"
    },
    "ProcessPayment": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:process-payment",
      "ResultPath": "$.paymentResult",
      "Catch": [{
        "ErrorEquals": ["States.ALL"],
        "ResultPath": "$.error",
        "Next": "ReleaseInventory"
      }],
      "Next": "ShipOrder"
    },
    "ShipOrder": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:ship-order",
      "ResultPath": "$.shipmentResult",
      "Catch": [{
        "ErrorEquals": ["States.ALL"],
        "ResultPath": "$.error",
        "Next": "RefundPayment"
      }],
      "Next": "OrderSucceeded"
    },
    "RefundPayment": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:refund-payment",
      "ResultPath": "$.refundResult",
      "Next": "ReleaseInventory"
    },
    "ReleaseInventory": {
      "Type": "Task",
      "Resource": "arn:aws:lambda:us-east-1:123456789012:function:release-inventory",
      "ResultPath": "$.releaseResult",
      "Next": "OrderFailed"
    },
    "OrderSucceeded": {
      "Type": "Succeed"
    },
    "OrderFailed": {
      "Type": "Fail",
      "Error": "OrderProcessingFailed",
      "Cause": "One or more steps in the saga failed"
    }
  }
}

# Terraform - Step Functions State Machine
resource "aws_sfn_state_machine" "order_saga" {
  name     = "order-processing-saga"
  role_arn = aws_iam_role.step_functions.arn
  type     = "STANDARD"

  definition = jsonencode({
    Comment = "Order Processing Saga with Compensating Transactions"
    StartAt = "ReserveInventory"
    States = {
      ReserveInventory = {
        Type     = "Task"
        Resource = aws_lambda_function.reserve_inventory.arn
        ResultPath = "$.inventoryResult"
        Retry = [{
          ErrorEquals     = ["Lambda.ServiceException"]
          IntervalSeconds = 2
          MaxAttempts     = 3
          BackoffRate     = 2
        }]
        Catch = [{
          ErrorEquals = ["States.ALL"]
          ResultPath  = "$.error"
          Next        = "OrderFailed"
        }]
        Next = "ProcessPayment"
      }
      ProcessPayment = {
        Type     = "Task"
        Resource = aws_lambda_function.process_payment.arn
        ResultPath = "$.paymentResult"
        Catch = [{
          ErrorEquals = ["States.ALL"]
          ResultPath  = "$.error"
          Next        = "CompensateInventory"
        }]
        Next = "ShipOrder"
      }
      ShipOrder = {
        Type     = "Task"
        Resource = aws_lambda_function.ship_order.arn
        ResultPath = "$.shipmentResult"
        Catch = [{
          ErrorEquals = ["States.ALL"]
          ResultPath  = "$.error"
          Next        = "CompensatePayment"
        }]
        Next = "OrderSucceeded"
      }
      CompensatePayment = {
        Type     = "Task"
        Resource = aws_lambda_function.refund_payment.arn
        ResultPath = "$.refundResult"
        Next = "CompensateInventory"
      }
      CompensateInventory = {
        Type     = "Task"
        Resource = aws_lambda_function.release_inventory.arn
        ResultPath = "$.releaseResult"
        Next = "OrderFailed"
      }
      OrderSucceeded = {
        Type = "Succeed"
      }
      OrderFailed = {
        Type  = "Fail"
        Error = "OrderProcessingFailed"
        Cause = "Saga rolled back due to failure"
      }
    }
  })

  logging_configuration {
    log_destination        = "${aws_cloudwatch_log_group.saga.arn}:*"
    include_execution_data = true
    level                  = "ALL"
  }

  tracing_configuration {
    enabled = true
  }
}

# Lambda Function Example - Reserve Inventory
resource "aws_lambda_function" "reserve_inventory" {
  function_name = "reserve-inventory"
  handler       = "index.handler"
  runtime       = "python3.11"
  role          = aws_iam_role.lambda.arn
  filename      = "reserve_inventory.zip"
  
  environment {
    variables = {
      INVENTORY_TABLE = aws_dynamodb_table.inventory.name
    }
  }
}

# Python Lambda - Reserve Inventory
"""
import boto3
import json

dynamodb = boto3.resource('dynamodb')
table = dynamodb.Table(os.environ['INVENTORY_TABLE'])

def handler(event, context):
    order_id = event['orderId']
    items = event['items']
    
    # Reserve inventory with conditional update
    for item in items:
        try:
            table.update_item(
                Key={'productId': item['productId']},
                UpdateExpression='SET reserved = reserved + :qty, available = available - :qty',
                ConditionExpression='available >= :qty',
                ExpressionAttributeValues={':qty': item['quantity']},
                ReturnValues='UPDATED_NEW'
            )
        except ClientError as e:
            if e.response['Error']['Code'] == 'ConditionalCheckFailedException':
                raise Exception(f"Insufficient inventory for {item['productId']}")
            raise
    
    return {
        'orderId': order_id,
        'reservationId': str(uuid.uuid4()),
        'status': 'RESERVED',
        'items': items
    }
"""

# Python Lambda - Release Inventory (Compensating Transaction)
"""
def handler(event, context):
    # Extract reservation from previous step
    inventory_result = event.get('inventoryResult', {})
    items = inventory_result.get('items', [])
    
    # Release reserved inventory
    for item in items:
        table.update_item(
            Key={'productId': item['productId']},
            UpdateExpression='SET reserved = reserved - :qty, available = available + :qty',
            ExpressionAttributeValues={':qty': item['quantity']}
        )
    
    return {
        'status': 'RELEASED',
        'items': items
    }
"""

# Trigger Saga from API Gateway
resource "aws_apigatewayv2_integration" "saga" {
  api_id             = aws_apigatewayv2_api.main.id
  integration_type   = "AWS_PROXY"
  integration_subtype = "StepFunctions-StartExecution"
  credentials_arn    = aws_iam_role.api_gateway.arn

  request_parameters = {
    "StateMachineArn" = aws_sfn_state_machine.order_saga.arn
    "Input"           = "$request.body"
  }
}

Key Design Decisions:
──────────────────────────────────────────────────────────────────────

1. Event Bus Selection:
   • EventBridge: Native AWS integration, schema registry, archiving
   • SNS + SQS: Simpler, fan-out pattern, proven reliability
   • MSK (Kafka): High throughput, replay, complex event processing
   
   Decision: EventBridge for AWS-native; MSK for high-volume/complex

2. Exactly-Once Processing:
   • SQS FIFO queues for ordering guarantees
   • DynamoDB conditional writes for idempotency
   • Idempotency keys in event payload
   
   # Idempotent handler pattern
   def handle_event(event):
       idempotency_key = event['detail']['idempotency_key']
       try:
           table.put_item(
               Item={'pk': idempotency_key, 'processed': True},
               ConditionExpression='attribute_not_exists(pk)'
           )
           # Process event...
       except ConditionalCheckFailedException:
           logger.info(f"Duplicate event: {idempotency_key}")

3. Error Handling:
   • Dead Letter Queues (DLQ) for failed messages
   • CloudWatch alarms on DLQ depth
   • Automated retry with exponential backoff
   • Manual reprocessing UI for DLQ items

4. Event Schema Evolution:
   • EventBridge Schema Registry for versioning
   • Backward compatible changes (additive only)
   • Consumer-driven contract testing

5. Observability:
   • X-Ray tracing across services
   • Correlation ID in all events
   • Service Map for dependency visualization
   • Custom metrics for business events
```

---

## 7. Interview Questions & Answers

### 7.1 Multi-Account & Organizations

**Q1: How would you design a multi-account strategy for a company moving from single account to enterprise scale?**

**A:** I would approach this systematically:

1. **Assessment Phase:**
   - Identify workloads and compliance requirements
   - Document current resource inventory
   - Define isolation boundaries (dev/staging/prod, teams, business units)

2. **Account Structure:**
```
Root (Management Account)
├── Security OU
│   ├── Security-Prod (GuardDuty, Security Hub)
│   └── Log-Archive (CloudTrail, Config logs)
├── Infrastructure OU
│   ├── Network-Hub (Transit Gateway, DNS)
│   └── Shared-Services (CI/CD, Container Registry)
├── Workloads OU
│   ├── Development OU
│   │   └── Dev accounts per team
│   ├── Staging OU
│   └── Production OU
└── Sandbox OU
```

3. **Key Components:**
   - Control Tower for account provisioning
   - SCPs for guardrails (deny actions, region restrictions)
   - Identity Center for centralized access
   - Transit Gateway for network hub

4. **Migration Strategy:**
   - Start with non-production workloads
   - Use AWS Resource Access Manager for shared resources
   - Implement StackSets for consistent configuration

---

**Q2: Explain how SCPs, IAM policies, and Permission Boundaries work together. What's the effective permission?**

**A:** 
```
Effective Permissions:
─────────────────────

                    ┌─────────────────┐
                    │      SCP        │
                    │  (Account max)  │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   Permission    │
                    │   Boundary      │
                    │   (User max)    │
                    └────────┬────────┘
                             │
                    ┌────────▼────────┐
                    │   IAM Policy    │
                    │   (Granted)     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │    Effective    │
                    │   Permission    │
                    │  (Intersection) │
                    └─────────────────┘
```

**Effective Permission = SCP ∩ Permission Boundary ∩ IAM Policy**

- **SCP:** Sets maximum permissions for all principals in an account
- **Permission Boundary:** Sets maximum permissions for a specific IAM entity
- **IAM Policy:** Grants actual permissions within those boundaries

Example: If SCP allows `s3:*`, Permission Boundary allows `s3:GetObject`, `s3:PutObject`, and IAM Policy grants `s3:*`, the effective permission is only `s3:GetObject`, `s3:PutObject`.

---

### 7.2 Disaster Recovery

**Q3: Compare DR strategies and when would you use each?**

**A:**

| Strategy | RTO | RPO | Cost | Use Case |
|----------|-----|-----|------|----------|
| **Backup/Restore** | Hours-Days | Hours | $ | Non-critical, batch systems |
| **Pilot Light** | 10min-Hours | Minutes | $$ | Web apps, databases with read replicas |
| **Warm Standby** | Minutes | Seconds | $$$ | Critical business apps |
| **Active-Active** | Near-zero | Near-zero | $$$$ | Customer-facing, revenue-generating |

**Decision Criteria:**
- **Backup/Restore:** RTO > 24 hours acceptable, cost-sensitive
- **Pilot Light:** RTO 10-60 minutes, core infrastructure pre-deployed
- **Warm Standby:** RTO < 15 minutes, can't afford extended downtime
- **Active-Active:** Zero downtime tolerance, global user base

---

**Q4: How would you implement Aurora Global Database failover?**

**A:**

1. **Planned Failover (Maintenance):**
```bash
aws rds failover-global-cluster \
  --global-cluster-identifier my-global-cluster \
  --target-db-cluster-identifier arn:aws:rds:us-west-2:123456789012:cluster:secondary
```
- Ensures no data loss (waits for replication)
- Takes 1-2 minutes
- Application needs connection string update (or use Global Database endpoint)

2. **Unplanned Failover (Disaster):**
   - Detach secondary and promote to standalone
   - May have minimal data loss (< 1 second typically)
   - Application must switch endpoints

3. **Application Configuration:**
```java
// Use reader endpoint for reads, fail back to cluster endpoint
String writerEndpoint = "global-cluster-endpoint.cluster-xxx.us-east-1.rds.amazonaws.com";
String readerEndpoint = "global-cluster-endpoint.cluster-ro-xxx.us-east-1.rds.amazonaws.com";
```

4. **Automation:**
   - Route 53 health checks trigger Lambda
   - Lambda initiates failover via SDK
   - SNS notification to ops team

---

### 7.3 Security

**Q5: Walk me through how you would investigate a GuardDuty finding for suspicious IAM activity.**

**A:**

1. **Initial Assessment:**
```
Finding: UnauthorizedAccess:IAMUser/InstanceCredentialExfiltration
Severity: High
Resource: EC2 instance i-1234567890abcdef0
```

2. **Investigation with Detective:**
   - Open finding in Detective
   - Review 30-day behavior baseline
   - Identify anomalous API calls
   - Track source IPs and geolocation

3. **CloudTrail Analysis:**
```bash
aws logs filter-log-events \
  --log-group-name /aws/cloudtrail \
  --filter-pattern '{ $.userIdentity.principalId = "*i-1234567890abcdef0*" }' \
  --start-time 1705000000000
```

4. **Containment (if malicious):**
   - Isolate instance: modify security group to deny all traffic
   - Rotate credentials: invalidate IAM role session
   - Preserve for forensics: create EBS snapshot

5. **Root Cause & Remediation:**
   - Identify entry point (misconfigured SG, vulnerable app)
   - Patch vulnerability
   - Implement IMDSv2 to prevent SSRF credential theft

---

**Q6: How do you implement defense-in-depth for a web application on AWS?**

**A:**

```
Defense in Depth Layers:
────────────────────────

┌─────────────────────────────────────────────────────────┐
│ 1. Edge Protection                                       │
│    • CloudFront + WAF (SQL injection, XSS, rate limit)  │
│    • AWS Shield (DDoS protection)                       │
│    • Route 53 (Health checks, geoblocking)              │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│ 2. Network Security                                      │
│    • VPC (Private subnets for app/data)                 │
│    • NACLs (Stateless subnet protection)                │
│    • Security Groups (Instance-level firewall)          │
│    • VPC Flow Logs (Network monitoring)                 │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│ 3. Application Security                                  │
│    • ALB (TLS termination, request validation)          │
│    • Cognito (Authentication)                           │
│    • Secrets Manager (No hardcoded credentials)         │
│    • IMDSv2 (Prevent SSRF credential theft)             │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│ 4. Data Protection                                       │
│    • KMS (Encryption at rest)                           │
│    • TLS (Encryption in transit)                        │
│    • S3 bucket policies (Principle of least privilege)  │
│    • RDS/Aurora encryption                              │
└─────────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────────┐
│ 5. Monitoring & Detection                                │
│    • GuardDuty (Threat detection)                       │
│    • Security Hub (Compliance)                          │
│    • CloudTrail (Audit logs)                            │
│    • Config (Configuration compliance)                   │
└─────────────────────────────────────────────────────────┘
```

---

### 7.4 Cost Optimization

**Q7: You notice AWS costs have increased 40% month-over-month. How do you investigate and reduce costs?**

**A:**

1. **Immediate Analysis:**
   - Cost Explorer: Group by Service, then by Usage Type
   - Check for cost anomaly alerts
   - Compare regions and accounts

2. **Common Culprits:**
   - Data transfer (inter-region, internet egress)
   - Orphaned resources (unused EBS, unattached EIPs)
   - Over-provisioned instances
   - Missing Savings Plans coverage

3. **Investigation Queries (CUR):**
```sql
-- Top cost drivers
SELECT product_product_name, 
       SUM(line_item_blended_cost) as cost
FROM cur_table
WHERE month = '2024-01'
GROUP BY 1 
ORDER BY 2 DESC LIMIT 10;

-- Data transfer analysis
SELECT product_product_name,
       line_item_usage_type,
       SUM(line_item_blended_cost) as cost
FROM cur_table
WHERE line_item_usage_type LIKE '%DataTransfer%'
GROUP BY 1, 2
ORDER BY 3 DESC;
```

4. **Optimization Actions:**
   - **Quick wins:** Delete unused resources, release EIPs
   - **Right-sizing:** Apply Compute Optimizer recommendations
   - **Commitments:** Purchase Savings Plans for stable workloads
   - **Architecture:** CloudFront for content delivery, VPC endpoints for AWS services

---

**Q8: Compare Reserved Instances, Savings Plans, and Spot Instances. When would you use each?**

**A:**

| Aspect | Reserved Instances | Savings Plans | Spot Instances |
|--------|-------------------|---------------|----------------|
| **Commitment** | Capacity | $/hour spend | None |
| **Flexibility** | Low (Standard) / Med (Convertible) | High (Compute SP) | Very High |
| **Savings** | Up to 72% | Up to 66% | Up to 90% |
| **Risk** | Lock-in | Lower lock-in | Interruption |
| **Best For** | Stable EC2-only | Variable workloads | Fault-tolerant |

**My Strategy:**
```
Baseline (always on):    Reserved or Savings Plans (60% of capacity)
Variable (elastic):      On-Demand with Auto Scaling (20%)
Batch/Non-critical:      Spot with diversification (20%)
```

**Key Considerations:**
- Use Savings Plans for new workloads (flexibility)
- Use Spot for stateless workers, batch jobs, dev environments
- Maintain On-Demand for bursty, unpredictable workloads

---

### 7.5 Architecture Design

**Q9: Design a highly available, scalable architecture for a real-time chat application.**

**A:**

```
Chat Application Architecture:
──────────────────────────────

Client (Web/Mobile)
        │
        ▼
┌─────────────────────────────────────────────────────────┐
│     CloudFront (WebSocket support)                       │
└─────────────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────┐
│     API Gateway (WebSocket API)                          │
│     • $connect, $disconnect, sendMessage routes         │
└─────────────────────────────────────────────────────────┘
        │
        ▼
┌─────────────────────────────────────────────────────────┐
│     Lambda Functions                                     │
│     • Connect: Store connection ID                       │
│     • Message: Fan-out to recipients                    │
│     • Disconnect: Remove connection ID                  │
└─────────────────────────────────────────────────────────┘
        │
        ├─────────────────┐
        ▼                 ▼
┌─────────────────┐ ┌─────────────────┐
│   DynamoDB      │ │   DynamoDB      │
│   (Connections) │ │   (Messages)    │
│                 │ │   (DAX Cache)   │
│   PK: room_id   │ │                 │
│   SK: conn_id   │ │   PK: room_id   │
│                 │ │   SK: timestamp │
└─────────────────┘ └─────────────────┘
```

**Key Design Decisions:**

1. **Why WebSocket over polling?**
   - Real-time bidirectional communication
   - Lower latency, reduced server load
   - API Gateway manages connections automatically

2. **Scalability:**
   - DynamoDB scales automatically
   - Lambda concurrent executions scale to 1000s
   - Connection state in DynamoDB (serverless)

3. **Message Fan-out:**
   - Query connections by room_id
   - Parallel Lambda invocations
   - API Gateway Management API to push messages

4. **High Availability:**
   - Multi-AZ by default (Lambda, DynamoDB)
   - Global Tables for multi-region

---

**Q10: How would you migrate a monolithic application to microservices on AWS?**

**A:**

**Phase 1: Assessment (Weeks 1-4)**
- Identify bounded contexts (domain-driven design)
- Map dependencies between components
- Prioritize services for extraction (strangler pattern)

**Phase 2: Foundation (Weeks 5-8)**
- Set up EKS or ECS clusters
- Implement service mesh (App Mesh) or API Gateway
- Establish CI/CD pipelines
- Set up observability (CloudWatch, X-Ray)

**Phase 3: Incremental Migration**
```
Strangler Fig Pattern:
──────────────────────

    Before:                          After:
    
    ┌─────────────┐              ┌─────────────┐
    │  Monolith   │              │  API        │──► Microservice A
    │             │              │  Gateway    │──► Microservice B
    │ A + B + C   │     →        │             │──► Monolith (C only)
    │             │              └─────────────┘
    └─────────────┘              
```

**Extraction Order:**
1. Stateless services first (easy to scale, test)
2. Services with clear boundaries
3. High-change-frequency components
4. Core domain services last

**Phase 4: Data Migration**
- Database per service pattern
- Event-driven sync during transition
- Eventually consistent where possible

**Anti-patterns to Avoid:**
- Big bang migration
- Distributed monolith (tight coupling)
- Shared database between services

---

### 7.6 Scenario-Based Questions

**Q11: Your application is experiencing intermittent latency spikes. How do you troubleshoot?**

**A:**

1. **Data Collection:**
```
CloudWatch Metrics:
• ELB: RequestCount, TargetResponseTime, HealthyHostCount
• EC2/ECS: CPUUtilization, MemoryUtilization
• RDS: DatabaseConnections, ReadLatency, WriteLatency
• Lambda: Duration, Throttles, ConcurrentExecutions
```

2. **X-Ray Analysis:**
   - Identify slow segments (which service?)
   - Look for outliers (P99 vs P50)
   - Check for downstream dependencies

3. **Common Causes & Solutions:**

| Symptom | Cause | Solution |
|---------|-------|----------|
| Periodic spikes | Cron jobs, batch processing | Spread workload, separate cluster |
| Gradual degradation | Memory leak, connection exhaustion | Fix leak, implement connection pooling |
| Random spikes | Noisy neighbor | Dedicated hosts, enhanced networking |
| After deployment | Code regression | Rollback, profiling |

4. **Database-Specific:**
   - Check for missing indexes (explain plan)
   - Connection pool exhaustion
   - Lock contention

5. **Resolution:**
   - Implement caching (ElastiCache)
   - Add read replicas for read-heavy workloads
   - Use async processing (SQS) for non-critical paths

---

**Q12: Design disaster recovery for a financial trading platform with RPO < 1 second and RTO < 1 minute.**

**A:**

**Architecture: Multi-Region Active-Active**

```
Trading Platform DR:
──────────────────

                    ┌─────────────────────────┐
                    │      Route 53           │
                    │   (Active-Active)       │
                    │   Latency Routing       │
                    └───────────┬─────────────┘
                                │
          ┌─────────────────────┴─────────────────────┐
          │                                           │
          ▼                                           ▼
┌─────────────────────────────────┐   ┌─────────────────────────────────┐
│     US-EAST-1 (Primary)          │   │     US-WEST-2 (Secondary)       │
│                                  │   │                                  │
│  ┌────────────────────────────┐  │   │  ┌────────────────────────────┐  │
│  │   Trading Engine (EKS)    │  │   │  │   Trading Engine (EKS)    │  │
│  │   Active-Active           │  │   │  │   Active-Active           │  │
│  └────────────┬───────────────┘  │   │  └────────────┬───────────────┘  │
│               │                  │   │               │                  │
│  ┌────────────▼───────────────┐  │   │  ┌────────────▼───────────────┐  │
│  │   Aurora Global Database  │  │   │  │   Aurora Global Database  │  │
│  │   (Primary Writer)        │  │   │  │   (Read w/ Write Forward) │  │
│  │   Sync Replication <1sec  │◄─┼───┼──►│   Promotable in <1 min   │  │
│  └────────────────────────────┘  │   │  └────────────────────────────┘  │
│                                  │   │                                  │
│  ┌────────────────────────────┐  │   │  ┌────────────────────────────┐  │
│  │   ElastiCache Global       │◄─┼───┼──►│  ElastiCache Global       │  │
│  │   (Order Book Cache)       │  │   │  │  (Order Book Cache)       │  │
│  └────────────────────────────┘  │   │  └────────────────────────────┘  │
│                                  │   │                                  │
└──────────────────────────────────┘   └──────────────────────────────────┘
```

**Key Decisions:**

1. **Data Replication:**
   - Aurora Global: < 1 second lag
   - Write forwarding to primary region
   - Synchronous replication for critical data

2. **Trade Order Consistency:**
   - Order IDs generated with region prefix
   - Conflict resolution: last-write-wins or timestamp-based
   - Financial transactions use 2PC or saga pattern

3. **Failover Automation:**
```python
# Lambda triggered by CloudWatch alarm
def failover_handler(event, context):
    # 1. Promote Aurora secondary
    rds.failover_global_cluster(
        GlobalClusterIdentifier='trading-db',
        TargetDbClusterIdentifier='secondary-cluster'
    )
    
    # 2. Update Route 53 (already automatic with health checks)
    
    # 3. Notify operations
    sns.publish(TopicArn=ops_topic, Message='Failover initiated')
```

4. **Testing:**
   - Monthly DR drills
   - Chaos engineering (randomly fail primary)
   - Measure actual RTO/RPO

---

## 8. Quick Reference

### 8.1 Service Comparison Matrix

| Use Case | Service | When to Use |
|----------|---------|-------------|
| Multi-account governance | AWS Organizations + Control Tower | Enterprise, compliance |
| Cost allocation | Cost Explorer, CUR | Any organization |
| Threat detection | GuardDuty | Always enable |
| Compliance | Security Hub + Config | Regulated industries |
| DR - Database | Aurora Global, DynamoDB Global | Critical data |
| DR - Compute | Pilot Light, Warm Standby | Based on RTO requirements |
| CI/CD | CodePipeline + CodeBuild | AWS-native |
| GitOps | ArgoCD on EKS | Kubernetes-centric |
| Secrets | Secrets Manager | Rotatable secrets |
| Config | Parameter Store | Static configuration |

### 8.2 Architecture Decision Guide

```
When to use Multi-Region:
────────────────────────
• RTO < 1 hour required
• Global user base (latency optimization)
• Compliance requires geo-redundancy
• Critical revenue-generating applications

When to use Multi-Account:
─────────────────────────
• Security boundaries between teams/environments
• Compliance isolation (PCI, HIPAA)
• Cost allocation by business unit
• Blast radius containment

When to use Serverless:
──────────────────────
• Variable/unpredictable traffic
• Event-driven processing
• Minimal operational overhead desired
• Sub-second scaling required
```

---

## 9. Summary

This guide covered enterprise-level AWS architecture and operations:

1. **Multi-Account Strategy:** AWS Organizations, Control Tower, SCPs, Identity Center
2. **Disaster Recovery:** Strategies from Backup/Restore to Active-Active
3. **Security Deep Dive:** GuardDuty, Security Hub, Config, Detective
4. **Cost Optimization:** Savings Plans, Spot strategies, FinOps practices
5. **DevOps & CI/CD:** CodePipeline, deployment strategies, GitOps
6. **Case Studies:** E-commerce, analytics, SaaS, event-driven architectures

**Key Takeaways for Platform Engineers:**

- Design for failure from day one
- Automate everything (IaC, CI/CD, remediation)
- Implement defense in depth
- Balance cost with reliability requirements
- Use managed services to reduce operational burden
- Continuously monitor and optimize

**Next Steps:**
- Practice building multi-account setups in a sandbox
- Implement a complete CI/CD pipeline
- Run DR drills regularly
- Stay updated with AWS re:Invent announcements

---

*Part 3 of 3 - For foundational concepts, see Part 1 (01-AWS-Foundations-Architecture-Perspective.md). For advanced services, see Part 2 (02-AWS-Advanced-Services-Platform-Engineering.md).*