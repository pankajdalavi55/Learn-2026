# The Complete Kubernetes Guide

## Table of Contents

1. [Introduction](#1-introduction)
2. [Core Concepts & Architecture](#2-core-concepts--architecture)
3. [Installation & Setup](#3-installation--setup)
4. [kubectl — The Kubernetes CLI](#4-kubectl--the-kubernetes-cli)
5. [Pods](#5-pods)
6. [ReplicaSets](#6-replicasets)
7. [Deployments](#7-deployments)
8. [Services](#8-services)
9. [Namespaces](#9-namespaces)
10. [ConfigMaps & Secrets](#10-configmaps--secrets)
11. [Volumes & Persistent Storage](#11-volumes--persistent-storage)
12. [StatefulSets](#12-statefulsets)
13. [DaemonSets](#13-daemonsets)
14. [Jobs & CronJobs](#14-jobs--cronjobs)
15. [Ingress](#15-ingress)
16. [Network Policies](#16-network-policies)
17. [RBAC — Role-Based Access Control](#17-rbac--role-based-access-control)
18. [Resource Management](#18-resource-management)
19. [Horizontal Pod Autoscaler (HPA)](#19-horizontal-pod-autoscaler-hpa)
20. [Helm — Package Manager](#20-helm--package-manager)
21. [Probes — Health Checks](#21-probes--health-checks)
22. [Taints, Tolerations & Affinity](#22-taints-tolerations--affinity)
23. [Logging & Monitoring](#23-logging--monitoring)
24. [Kubernetes on Cloud (EKS, AKS, GKE)](#24-kubernetes-on-cloud-eks-aks-gke)
25. [Security Best Practices](#25-security-best-practices)
26. [Troubleshooting](#26-troubleshooting)
27. [Real-World Architecture Patterns](#27-real-world-architecture-patterns)
28. [Interview Questions](#28-interview-questions)

---

## 1. Introduction

### What is Kubernetes?

Kubernetes (K8s) is an open-source container orchestration platform originally developed by Google, now maintained by the Cloud Native Computing Foundation (CNCF). It automates deployment, scaling, and management of containerized applications.

### Why Kubernetes?

| Challenge | How Kubernetes Solves It |
|-----------|--------------------------|
| Manual scaling | Auto-scaling based on CPU/memory/custom metrics |
| Downtime during deployments | Rolling updates with zero downtime |
| Container failures | Self-healing — restarts failed containers automatically |
| Service discovery | Built-in DNS and service discovery |
| Configuration management | ConfigMaps and Secrets |
| Load balancing | Service abstraction with built-in load balancing |
| Storage orchestration | Persistent Volumes with pluggable storage backends |

### Kubernetes vs Docker Swarm vs Nomad

| Feature | Kubernetes | Docker Swarm | Nomad |
|---------|-----------|--------------|-------|
| Auto-scaling | Yes (HPA, VPA, Cluster Autoscaler) | Limited | Yes |
| Service mesh | Istio, Linkerd | No native support | Consul Connect |
| Learning curve | Steep | Easy | Moderate |
| Community | Massive | Declining | Growing |
| Production readiness | Battle-tested | Limited adoption | Growing adoption |
| Multi-cloud | Excellent | Limited | Good |

---

## 2. Core Concepts & Architecture

### Cluster Architecture

```
┌──────────────────────────────────────────────────────────────────┐
│                        KUBERNETES CLUSTER                        │
│                                                                  │
│  ┌────────────────────────── Control Plane ───────────────────┐  │
│  │                                                            │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌────────────────┐   │  │
│  │  │  API Server   │  │  Scheduler   │  │   Controller   │   │  │
│  │  │  (kube-api)   │  │              │  │    Manager     │   │  │
│  │  └──────┬───────┘  └──────────────┘  └────────────────┘   │  │
│  │         │                                                  │  │
│  │  ┌──────▼───────┐  ┌──────────────┐                       │  │
│  │  │     etcd      │  │   Cloud      │                       │  │
│  │  │  (key-value   │  │  Controller  │                       │  │
│  │  │   store)      │  │   Manager    │                       │  │
│  │  └──────────────┘  └──────────────┘                       │  │
│  └────────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌──────── Worker Node 1 ───────┐  ┌──── Worker Node 2 ──────┐  │
│  │                              │  │                          │  │
│  │  ┌────────┐  ┌────────────┐  │  │  ┌────────┐  ┌───────┐  │  │
│  │  │ kubelet│  │ kube-proxy │  │  │  │ kubelet│  │ kube- │  │  │
│  │  └────────┘  └────────────┘  │  │  └────────┘  │ proxy │  │  │
│  │                              │  │               └───────┘  │  │
│  │  ┌──────────────────────┐    │  │  ┌───────────────────┐   │  │
│  │  │  Container Runtime   │    │  │  │ Container Runtime │   │  │
│  │  │  (containerd/CRI-O)  │    │  │  │ (containerd)      │   │  │
│  │  └──────────────────────┘    │  │  └───────────────────┘   │  │
│  │                              │  │                          │  │
│  │  ┌──────┐ ┌──────┐ ┌─────┐  │  │  ┌──────┐ ┌──────┐      │  │
│  │  │ Pod  │ │ Pod  │ │ Pod │  │  │  │ Pod  │ │ Pod  │      │  │
│  │  └──────┘ └──────┘ └─────┘  │  │  └──────┘ └──────┘      │  │
│  └──────────────────────────────┘  └──────────────────────────┘  │
└──────────────────────────────────────────────────────────────────┘
```

### Control Plane Components

| Component | Purpose |
|-----------|---------|
| **kube-apiserver** | Front-end for the Kubernetes control plane. All communication goes through the API server. |
| **etcd** | Consistent and highly-available key-value store used as Kubernetes' backing store for all cluster data. |
| **kube-scheduler** | Watches for newly created Pods with no assigned node and selects a node for them to run on. |
| **kube-controller-manager** | Runs controller processes (Node Controller, ReplicaSet Controller, Endpoint Controller, etc.). |
| **cloud-controller-manager** | Embeds cloud-specific control logic (load balancers, storage volumes, routes). |

### Worker Node Components

| Component | Purpose |
|-----------|---------|
| **kubelet** | Agent that runs on each node. Ensures containers are running in a Pod. |
| **kube-proxy** | Network proxy that maintains network rules on nodes for Service communication. |
| **Container runtime** | Software responsible for running containers (containerd, CRI-O). |

### Key Kubernetes Objects

```
Workloads:       Pod → ReplicaSet → Deployment / StatefulSet / DaemonSet / Job
Networking:      Service → Ingress → NetworkPolicy
Configuration:   ConfigMap, Secret
Storage:         PersistentVolume (PV), PersistentVolumeClaim (PVC), StorageClass
Access Control:  ServiceAccount, Role, ClusterRole, RoleBinding, ClusterRoleBinding
```

---

## 3. Installation & Setup

### Local Development Options

#### Minikube

```bash
# Install minikube
brew install minikube          # macOS
choco install minikube         # Windows
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 \
  && sudo install minikube-linux-amd64 /usr/local/bin/minikube   # Linux

# Start cluster
minikube start --driver=docker --cpus=4 --memory=8192

# Check status
minikube status

# Access Kubernetes dashboard
minikube dashboard

# Stop/Delete cluster
minikube stop
minikube delete
```

#### kind (Kubernetes IN Docker)

```bash
# Install kind
brew install kind              # macOS
go install sigs.k8s.io/kind@latest   # Go

# Create cluster
kind create cluster --name my-cluster

# Create multi-node cluster
cat <<EOF | kind create cluster --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
EOF

# Delete cluster
kind delete cluster --name my-cluster
```

#### Docker Desktop

Enable Kubernetes in Docker Desktop Settings → Kubernetes → Enable Kubernetes.

### Install kubectl

```bash
# macOS
brew install kubectl

# Linux
curl -LO "https://dl.k8s.io/release/$(curl -Ls https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Verify
kubectl version --client
```

### kubeconfig

The kubeconfig file (`~/.kube/config`) stores cluster access information.

```yaml
apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://127.0.0.1:6443
    certificate-authority: /path/to/ca.crt
  name: my-cluster
contexts:
- context:
    cluster: my-cluster
    user: admin
    namespace: default
  name: my-context
current-context: my-context
users:
- name: admin
  user:
    client-certificate: /path/to/admin.crt
    client-key: /path/to/admin.key
```

```bash
# View current config
kubectl config view

# Switch context
kubectl config use-context my-context

# Set default namespace for a context
kubectl config set-context --current --namespace=my-namespace
```

---

## 4. kubectl — The Kubernetes CLI

### Syntax

```
kubectl [command] [TYPE] [NAME] [flags]
```

### Essential Commands

```bash
# Cluster info
kubectl cluster-info
kubectl get nodes
kubectl get componentstatuses

# Creating resources
kubectl apply -f manifest.yaml          # Declarative (preferred)
kubectl create -f manifest.yaml         # Imperative
kubectl run nginx --image=nginx         # Quick run

# Reading resources
kubectl get pods                        # List pods
kubectl get pods -o wide                # Wide output with node info
kubectl get pods -o yaml                # Full YAML output
kubectl get pods -o json                # JSON output
kubectl get all                         # All resources in namespace
kubectl get pods --all-namespaces       # Across all namespaces (or -A)
kubectl get pods -l app=nginx           # Filter by label
kubectl get pods --field-selector status.phase=Running
kubectl get pods --sort-by='.metadata.creationTimestamp'
kubectl get pods -w                     # Watch for changes

# Detailed info
kubectl describe pod <pod-name>
kubectl describe node <node-name>

# Updating resources
kubectl edit deployment <name>          # Edit in-place
kubectl set image deployment/nginx nginx=nginx:1.25
kubectl scale deployment nginx --replicas=5
kubectl patch deployment nginx -p '{"spec":{"replicas":3}}'

# Deleting resources
kubectl delete pod <pod-name>
kubectl delete -f manifest.yaml
kubectl delete pods --all
kubectl delete pods -l app=nginx

# Debugging
kubectl logs <pod-name>                 # Container logs
kubectl logs <pod-name> -c <container>  # Specific container
kubectl logs <pod-name> --previous      # Previous container instance
kubectl logs -f <pod-name>              # Follow/stream logs
kubectl exec -it <pod-name> -- /bin/sh  # Shell into container
kubectl port-forward pod/<name> 8080:80 # Port forward
kubectl top pods                        # Resource usage
kubectl top nodes

# Dry-run & diff
kubectl apply -f manifest.yaml --dry-run=client
kubectl diff -f manifest.yaml
```

### Output Formatting with JSONPath

```bash
# Get pod IPs
kubectl get pods -o jsonpath='{.items[*].status.podIP}'

# Custom columns
kubectl get pods -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,NODE:.spec.nodeName

# Get specific field
kubectl get pod nginx -o jsonpath='{.spec.containers[0].image}'
```

---

## 5. Pods

### What is a Pod?

A Pod is the smallest deployable unit in Kubernetes. It represents a single instance of a running process and can contain one or more containers that share:
- Network namespace (same IP address and port space)
- Storage volumes
- Process namespace (optional)

### Pod Lifecycle

```
Pending → Running → Succeeded/Failed
              ↓
          Unknown (node communication lost)
```

| Phase | Description |
|-------|-------------|
| **Pending** | Pod accepted but containers not yet created. Includes image pulling. |
| **Running** | At least one container is running. |
| **Succeeded** | All containers terminated successfully (exit code 0). |
| **Failed** | At least one container terminated with failure. |
| **Unknown** | Pod state cannot be determined (usually node communication failure). |

### Basic Pod Manifest

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
    environment: dev
  annotations:
    description: "Simple nginx web server"
spec:
  containers:
  - name: nginx
    image: nginx:1.25
    ports:
    - containerPort: 80
    resources:
      requests:
        memory: "64Mi"
        cpu: "250m"
      limits:
        memory: "128Mi"
        cpu: "500m"
```

```bash
kubectl apply -f nginx-pod.yaml
kubectl get pod nginx-pod
kubectl describe pod nginx-pod
kubectl delete pod nginx-pod
```

### Multi-Container Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: multi-container-pod
spec:
  containers:
  - name: app
    image: my-app:1.0
    ports:
    - containerPort: 8080
    volumeMounts:
    - name: shared-logs
      mountPath: /var/log/app

  - name: log-sidecar
    image: busybox
    command: ["sh", "-c", "tail -f /var/log/app/app.log"]
    volumeMounts:
    - name: shared-logs
      mountPath: /var/log/app

  volumes:
  - name: shared-logs
    emptyDir: {}
```

### Multi-Container Pod Patterns

| Pattern | Description | Example |
|---------|-------------|---------|
| **Sidecar** | Helper container that enhances the main container | Log shipper, config syncer |
| **Ambassador** | Proxy container that handles external communication | Local proxy for database connection |
| **Adapter** | Container that transforms output of the main container | Log format normalizer |

### Init Containers

Init containers run **before** app containers start. They run to completion sequentially.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app-with-init
spec:
  initContainers:
  - name: wait-for-db
    image: busybox
    command: ['sh', '-c', 'until nc -z mysql-service 3306; do echo waiting for db; sleep 2; done']

  - name: init-schema
    image: mysql:8
    command: ['sh', '-c', 'mysql -h mysql-service -u root -p$MYSQL_ROOT_PASSWORD < /schema/init.sql']
    env:
    - name: MYSQL_ROOT_PASSWORD
      valueFrom:
        secretKeyRef:
          name: mysql-secret
          key: password
    volumeMounts:
    - name: schema
      mountPath: /schema

  containers:
  - name: app
    image: my-app:1.0
    ports:
    - containerPort: 8080

  volumes:
  - name: schema
    configMap:
      name: db-schema
```

### Static Pods

Static Pods are managed directly by the kubelet on a specific node, without the API server observing them. They are defined by placing manifest files in the kubelet's configured static pod path (usually `/etc/kubernetes/manifests/`).

```bash
# Find static pod path
ps aux | grep kubelet   # look for --pod-manifest-path or --config
# In kubelet config: staticPodPath: /etc/kubernetes/manifests
```

---

## 6. ReplicaSets

A ReplicaSet ensures a specified number of pod replicas are running at any given time.

### ReplicaSet Manifest

```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: nginx-replicaset
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.25
        ports:
        - containerPort: 80
        resources:
          requests:
            cpu: "100m"
            memory: "64Mi"
```

```bash
kubectl apply -f replicaset.yaml
kubectl get rs
kubectl scale rs nginx-replicaset --replicas=5
kubectl delete rs nginx-replicaset
```

> **Note:** You should almost always use a Deployment instead of directly creating ReplicaSets. Deployments manage ReplicaSets and provide declarative updates.

---

## 7. Deployments

### What is a Deployment?

A Deployment provides declarative updates for Pods and ReplicaSets. It manages the lifecycle of your application: creating, updating, scaling, and rolling back.

```
Deployment → manages → ReplicaSet → manages → Pods
```

### Deployment Manifest

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1           # Max pods above desired count during update
      maxUnavailable: 0     # Max pods unavailable during update (zero-downtime)
  template:
    metadata:
      labels:
        app: nginx
        version: "1.25"
    spec:
      containers:
      - name: nginx
        image: nginx:1.25
        ports:
        - containerPort: 80
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "250m"
            memory: "256Mi"
        readinessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 5
          periodSeconds: 10
        livenessProbe:
          httpGet:
            path: /
            port: 80
          initialDelaySeconds: 15
          periodSeconds: 20
```

### Deployment Strategies

#### RollingUpdate (default)

Gradually replaces old pods with new ones. Zero-downtime deployment.

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 25%          # Can go 25% over desired replica count
    maxUnavailable: 25%    # 25% of pods can be unavailable during update
```

#### Recreate

Terminates all existing pods before creating new ones. Causes downtime.

```yaml
strategy:
  type: Recreate
```

### Deployment Operations

```bash
# Create
kubectl apply -f deployment.yaml

# Update image
kubectl set image deployment/nginx-deployment nginx=nginx:1.26

# Check rollout status
kubectl rollout status deployment/nginx-deployment

# View rollout history
kubectl rollout history deployment/nginx-deployment
kubectl rollout history deployment/nginx-deployment --revision=2

# Rollback
kubectl rollout undo deployment/nginx-deployment                # Previous version
kubectl rollout undo deployment/nginx-deployment --to-revision=1  # Specific revision

# Pause/Resume rollout
kubectl rollout pause deployment/nginx-deployment
kubectl rollout resume deployment/nginx-deployment

# Scale
kubectl scale deployment/nginx-deployment --replicas=5

# Restart all pods (rolling restart)
kubectl rollout restart deployment/nginx-deployment
```

### Blue-Green Deployment Pattern

```yaml
# blue-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-blue
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
      version: blue
  template:
    metadata:
      labels:
        app: my-app
        version: blue
    spec:
      containers:
      - name: app
        image: my-app:1.0
---
# green-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-green
spec:
  replicas: 3
  selector:
    matchLabels:
      app: my-app
      version: green
  template:
    metadata:
      labels:
        app: my-app
        version: green
    spec:
      containers:
      - name: app
        image: my-app:2.0
---
# Switch traffic by updating the Service selector
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
spec:
  selector:
    app: my-app
    version: green     # Switch from blue → green
  ports:
  - port: 80
    targetPort: 8080
```

### Canary Deployment Pattern

```yaml
# Stable deployment (90% traffic)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-stable
spec:
  replicas: 9
  selector:
    matchLabels:
      app: my-app
      track: stable
  template:
    metadata:
      labels:
        app: my-app
        track: stable
    spec:
      containers:
      - name: app
        image: my-app:1.0
---
# Canary deployment (10% traffic)
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-canary
spec:
  replicas: 1
  selector:
    matchLabels:
      app: my-app
      track: canary
  template:
    metadata:
      labels:
        app: my-app
        track: canary
    spec:
      containers:
      - name: app
        image: my-app:2.0
---
# Service selects both (by shared label)
apiVersion: v1
kind: Service
metadata:
  name: my-app-service
spec:
  selector:
    app: my-app          # Matches both stable and canary
  ports:
  - port: 80
    targetPort: 8080
```

---

## 8. Services

### What is a Service?

A Service is an abstraction that defines a logical set of Pods and a policy to access them. Services provide stable networking for Pods (which have ephemeral IPs).

### Service Types

```
┌──────────────────────────────────────────────────────────────┐
│                      Service Types                           │
│                                                              │
│  ClusterIP ──► Internal-only access within the cluster       │
│  NodePort  ──► Exposes on each Node's IP at a static port    │
│  LoadBalancer ► Provisions external load balancer (cloud)     │
│  ExternalName ► Maps to an external DNS name (CNAME)         │
└──────────────────────────────────────────────────────────────┘
```

### ClusterIP (default)

Only accessible within the cluster.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-service
spec:
  type: ClusterIP
  selector:
    app: backend
  ports:
  - name: http
    port: 80             # Service port
    targetPort: 8080     # Container port
    protocol: TCP
```

### NodePort

Exposes the service on each node's IP at a static port (30000-32767).

```yaml
apiVersion: v1
kind: Service
metadata:
  name: frontend-service
spec:
  type: NodePort
  selector:
    app: frontend
  ports:
  - port: 80
    targetPort: 8080
    nodePort: 30080      # Optional: auto-assigned if omitted
```

Access via: `http://<node-ip>:30080`

### LoadBalancer

Provisions an external load balancer (works on cloud providers).

```yaml
apiVersion: v1
kind: Service
metadata:
  name: web-service
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: "nlb"
spec:
  type: LoadBalancer
  selector:
    app: web
  ports:
  - port: 80
    targetPort: 8080
```

### ExternalName

Maps a Service to a DNS name (no proxying).

```yaml
apiVersion: v1
kind: Service
metadata:
  name: external-db
spec:
  type: ExternalName
  externalName: db.example.com
```

### Headless Service

A headless service (clusterIP: None) doesn't allocate a cluster IP. DNS returns the individual Pod IPs directly. Used with StatefulSets.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: db-headless
spec:
  clusterIP: None
  selector:
    app: database
  ports:
  - port: 5432
    targetPort: 5432
```

### Service DNS Resolution

```
<service-name>.<namespace>.svc.cluster.local

# Examples:
backend-service.default.svc.cluster.local
redis.cache.svc.cluster.local
```

### Service Commands

```bash
kubectl get svc
kubectl describe svc backend-service
kubectl get endpoints backend-service    # See which pods back the service
```

---

## 9. Namespaces

Namespaces provide a mechanism for isolating groups of resources within a single cluster.

### Default Namespaces

| Namespace | Purpose |
|-----------|---------|
| **default** | For resources with no other namespace |
| **kube-system** | For objects created by the Kubernetes system |
| **kube-public** | Readable by all users, reserved for cluster usage |
| **kube-node-lease** | Holds Lease objects for node heartbeats |

### Managing Namespaces

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: development
  labels:
    environment: dev
```

```bash
# Create
kubectl create namespace staging
kubectl apply -f namespace.yaml

# List
kubectl get namespaces

# Set default namespace
kubectl config set-context --current --namespace=development

# Deploy to specific namespace
kubectl apply -f deployment.yaml -n staging

# Get resources in specific namespace
kubectl get pods -n kube-system

# Delete namespace (deletes ALL resources within it)
kubectl delete namespace development
```

### Resource Quotas

Limit total resource consumption per namespace.

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: dev-quota
  namespace: development
spec:
  hard:
    pods: "20"
    requests.cpu: "4"
    requests.memory: "8Gi"
    limits.cpu: "8"
    limits.memory: "16Gi"
    persistentvolumeclaims: "10"
    services.loadbalancers: "2"
    services.nodeports: "5"
```

### Limit Ranges

Set default resource constraints for individual pods/containers in a namespace.

```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: default-limits
  namespace: development
spec:
  limits:
  - default:              # Default limits
      cpu: "500m"
      memory: "256Mi"
    defaultRequest:       # Default requests
      cpu: "100m"
      memory: "128Mi"
    max:
      cpu: "2"
      memory: "1Gi"
    min:
      cpu: "50m"
      memory: "64Mi"
    type: Container
```

---

## 10. ConfigMaps & Secrets

### ConfigMaps

ConfigMaps store non-confidential configuration data as key-value pairs.

#### Creating ConfigMaps

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  # Simple key-value pairs
  DATABASE_HOST: "mysql-service"
  DATABASE_PORT: "3306"
  LOG_LEVEL: "info"

  # File-like keys
  application.properties: |
    server.port=8080
    spring.datasource.url=jdbc:mysql://mysql-service:3306/mydb
    spring.jpa.hibernate.ddl-auto=update

  nginx.conf: |
    server {
      listen 80;
      server_name localhost;
      location / {
        proxy_pass http://backend:8080;
      }
    }
```

```bash
# Create from literal values
kubectl create configmap app-config \
  --from-literal=DATABASE_HOST=mysql-service \
  --from-literal=DATABASE_PORT=3306

# Create from file
kubectl create configmap nginx-config --from-file=nginx.conf

# Create from directory
kubectl create configmap app-configs --from-file=./config-dir/
```

#### Using ConfigMaps in Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
  - name: app
    image: my-app:1.0

    # Method 1: Environment variables from specific keys
    env:
    - name: DB_HOST
      valueFrom:
        configMapKeyRef:
          name: app-config
          key: DATABASE_HOST

    # Method 2: All keys as environment variables
    envFrom:
    - configMapRef:
        name: app-config

    # Method 3: Mount as volume (files)
    volumeMounts:
    - name: config-volume
      mountPath: /etc/config
      readOnly: true

  volumes:
  - name: config-volume
    configMap:
      name: app-config
      items:                        # Optional: mount specific keys
      - key: nginx.conf
        path: nginx.conf
```

### Secrets

Secrets store sensitive data (passwords, tokens, keys). Values are base64-encoded (not encrypted by default).

#### Creating Secrets

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-credentials
type: Opaque
data:
  username: YWRtaW4=          # base64 encoded "admin"
  password: cEBzc3cwcmQ=      # base64 encoded "p@ssw0rd"
---
# Using stringData (plain text, encoded automatically)
apiVersion: v1
kind: Secret
metadata:
  name: db-credentials
type: Opaque
stringData:
  username: admin
  password: p@ssw0rd
```

```bash
# Create from literals
kubectl create secret generic db-credentials \
  --from-literal=username=admin \
  --from-literal=password='p@ssw0rd'

# Create TLS secret
kubectl create secret tls tls-secret \
  --cert=path/to/tls.crt \
  --key=path/to/tls.key

# Create Docker registry secret
kubectl create secret docker-registry regcred \
  --docker-server=https://index.docker.io/v1/ \
  --docker-username=user \
  --docker-password=pass \
  --docker-email=user@example.com

# Base64 encode/decode
echo -n 'admin' | base64          # YWRtaW4=
echo 'YWRtaW4=' | base64 -d       # admin
```

#### Secret Types

| Type | Description |
|------|-------------|
| `Opaque` | Arbitrary user-defined data (default) |
| `kubernetes.io/tls` | TLS certificate and key |
| `kubernetes.io/dockerconfigjson` | Docker registry credentials |
| `kubernetes.io/basic-auth` | Basic authentication credentials |
| `kubernetes.io/ssh-auth` | SSH authentication credentials |
| `kubernetes.io/service-account-token` | Service account token |

#### Using Secrets in Pods

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
  - name: app
    image: my-app:1.0

    env:
    - name: DB_USERNAME
      valueFrom:
        secretKeyRef:
          name: db-credentials
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: db-credentials
          key: password

    volumeMounts:
    - name: secret-volume
      mountPath: /etc/secrets
      readOnly: true

  volumes:
  - name: secret-volume
    secret:
      secretName: db-credentials

  imagePullSecrets:              # For private registries
  - name: regcred
```

> **Security Note:** Enable encryption at rest for Secrets in etcd using `EncryptionConfiguration`. Consider external secret managers (Vault, AWS Secrets Manager) for production.

---

## 11. Volumes & Persistent Storage

### Volume Types

| Type | Lifetime | Description |
|------|----------|-------------|
| `emptyDir` | Pod | Temporary directory, deleted when pod is removed |
| `hostPath` | Node | Mounts a file/directory from the host node |
| `persistentVolumeClaim` | Cluster | Uses PersistentVolume storage |
| `configMap` / `secret` | Cluster | Mounts ConfigMap/Secret as files |
| `nfs` | External | Network File System mount |
| `csi` | External | Container Storage Interface driver |

### emptyDir

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: shared-data
spec:
  containers:
  - name: writer
    image: busybox
    command: ["sh", "-c", "echo hello > /data/message; sleep 3600"]
    volumeMounts:
    - name: shared
      mountPath: /data

  - name: reader
    image: busybox
    command: ["sh", "-c", "cat /data/message; sleep 3600"]
    volumeMounts:
    - name: shared
      mountPath: /data

  volumes:
  - name: shared
    emptyDir: {}
    # emptyDir:
    #   medium: Memory        # Use RAM-backed tmpfs
    #   sizeLimit: 100Mi
```

### hostPath

```yaml
volumes:
- name: host-data
  hostPath:
    path: /data/app
    type: DirectoryOrCreate    # Creates if doesn't exist
```

> **Warning:** `hostPath` should only be used for development/testing or DaemonSets. It ties pods to specific nodes.

### PersistentVolume (PV) & PersistentVolumeClaim (PVC)

```
StorageClass ──► PersistentVolume (PV) ◄── PersistentVolumeClaim (PVC) ◄── Pod
  (provisioner)     (actual storage)          (request for storage)
```

#### StorageClass

```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast-ssd
provisioner: kubernetes.io/aws-ebs    # Or pd.csi.storage.gke.io, etc.
parameters:
  type: gp3
  iopsPerGB: "10"
reclaimPolicy: Delete         # Delete or Retain
allowVolumeExpansion: true
volumeBindingMode: WaitForFirstConsumer
```

#### PersistentVolume (static provisioning)

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-data
spec:
  capacity:
    storage: 10Gi
  accessModes:
  - ReadWriteOnce           # RWO: single node read-write
  # - ReadOnlyMany          # ROX: multi-node read-only
  # - ReadWriteMany         # RWX: multi-node read-write
  persistentVolumeReclaimPolicy: Retain
  storageClassName: fast-ssd
  hostPath:                  # For local testing
    path: /mnt/data
```

#### PersistentVolumeClaim

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: app-data-pvc
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: fast-ssd   # Match the StorageClass
```

#### Using PVC in a Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  containers:
  - name: app
    image: my-app:1.0
    volumeMounts:
    - name: data
      mountPath: /app/data
  volumes:
  - name: data
    persistentVolumeClaim:
      claimName: app-data-pvc
```

```bash
kubectl get pv
kubectl get pvc
kubectl describe pv pv-data
```

### Access Modes Explained

| Mode | Abbreviation | Description |
|------|-------------|-------------|
| ReadWriteOnce | RWO | Mounted as read-write by a single node |
| ReadOnlyMany | ROX | Mounted as read-only by many nodes |
| ReadWriteMany | RWX | Mounted as read-write by many nodes |
| ReadWriteOncePod | RWOP | Mounted as read-write by a single pod |

---

## 12. StatefulSets

### What is a StatefulSet?

StatefulSets are used for applications that require one or more of:
- Stable, unique network identifiers
- Stable, persistent storage
- Ordered, graceful deployment and scaling
- Ordered, automated rolling updates

Common use cases: databases (MySQL, PostgreSQL), message queues (Kafka), distributed systems (ZooKeeper, Elasticsearch).

### StatefulSet Manifest

```yaml
apiVersion: v1
kind: Service
metadata:
  name: mysql-headless
spec:
  clusterIP: None             # Headless service required
  selector:
    app: mysql
  ports:
  - port: 3306
    targetPort: 3306
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: mysql
spec:
  serviceName: mysql-headless   # Must reference headless service
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql:8.0
        ports:
        - containerPort: 3306
        env:
        - name: MYSQL_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-secret
              key: root-password
        volumeMounts:
        - name: mysql-data
          mountPath: /var/lib/mysql
        resources:
          requests:
            cpu: "500m"
            memory: "512Mi"
          limits:
            cpu: "1"
            memory: "1Gi"

  volumeClaimTemplates:          # Each pod gets its own PVC
  - metadata:
      name: mysql-data
    spec:
      accessModes: ["ReadWriteOnce"]
      storageClassName: fast-ssd
      resources:
        requests:
          storage: 10Gi
```

### StatefulSet Behavior

| Feature | Behavior |
|---------|----------|
| **Pod names** | `mysql-0`, `mysql-1`, `mysql-2` (ordinal index) |
| **DNS** | `mysql-0.mysql-headless.default.svc.cluster.local` |
| **Creation order** | Sequential: 0 → 1 → 2 |
| **Deletion order** | Reverse: 2 → 1 → 0 |
| **Storage** | Each pod gets its own PVC that persists across restarts |
| **Scaling down** | PVCs are NOT deleted (data preserved) |

### Deployment vs StatefulSet

| Feature | Deployment | StatefulSet |
|---------|-----------|-------------|
| Pod identity | Random names | Ordered, stable names |
| Storage | Shared PVC | Per-pod PVC |
| Scaling | Parallel | Sequential (ordered) |
| DNS | Via Service | Individual pod DNS |
| Use case | Stateless apps | Stateful apps (databases) |

---

## 13. DaemonSets

### What is a DaemonSet?

A DaemonSet ensures that all (or some) nodes run a copy of a Pod. As nodes are added to the cluster, Pods are added to them.

Common use cases:
- Log collection (Fluentd, Filebeat)
- Node monitoring (Prometheus Node Exporter)
- Network plugins (Calico, Cilium)
- Storage daemons (Ceph, GlusterFS)

### DaemonSet Manifest

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: kube-system
  labels:
    app: fluentd
spec:
  selector:
    matchLabels:
      app: fluentd
  template:
    metadata:
      labels:
        app: fluentd
    spec:
      tolerations:
      - key: node-role.kubernetes.io/control-plane
        effect: NoSchedule       # Run on control-plane nodes too
      containers:
      - name: fluentd
        image: fluentd:v1.16
        resources:
          requests:
            cpu: "100m"
            memory: "200Mi"
          limits:
            cpu: "200m"
            memory: "400Mi"
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: containers
          mountPath: /var/lib/docker/containers
          readOnly: true
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: containers
        hostPath:
          path: /var/lib/docker/containers
```

```bash
kubectl get daemonsets -n kube-system
kubectl describe daemonset fluentd -n kube-system
```

---

## 14. Jobs & CronJobs

### Job

A Job creates one or more Pods and ensures a specified number of them successfully terminate.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: data-migration
spec:
  completions: 1          # Number of successful completions needed
  parallelism: 1          # Number of pods running in parallel
  backoffLimit: 3          # Retries before marking as failed
  activeDeadlineSeconds: 600  # Max runtime
  ttlSecondsAfterFinished: 100  # Auto-cleanup after completion
  template:
    spec:
      restartPolicy: Never    # Must be Never or OnFailure
      containers:
      - name: migrate
        image: my-app:1.0
        command: ["python", "migrate.py"]
        env:
        - name: DB_URL
          valueFrom:
            secretKeyRef:
              name: db-credentials
              key: url
```

### Parallel Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: batch-process
spec:
  completions: 10        # Need 10 completions total
  parallelism: 3         # Run 3 pods at a time
  template:
    spec:
      restartPolicy: OnFailure
      containers:
      - name: worker
        image: batch-worker:1.0
```

### CronJob

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: backup-db
spec:
  schedule: "0 2 * * *"           # Every day at 2 AM
  concurrencyPolicy: Forbid        # Don't run concurrent jobs
  successfulJobsHistoryLimit: 3
  failedJobsHistoryLimit: 1
  startingDeadlineSeconds: 200
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
          - name: backup
            image: backup-tool:1.0
            command: ["/bin/sh", "-c", "pg_dump -h $DB_HOST mydb > /backup/dump.sql"]
            env:
            - name: DB_HOST
              value: "postgres-service"
            volumeMounts:
            - name: backup-storage
              mountPath: /backup
          volumes:
          - name: backup-storage
            persistentVolumeClaim:
              claimName: backup-pvc
```

### Cron Schedule Syntax

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of the month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of the week (0 - 6) (Sunday = 0)
│ │ │ │ │
* * * * *
```

| Expression | Description |
|-----------|-------------|
| `*/5 * * * *` | Every 5 minutes |
| `0 * * * *` | Every hour |
| `0 0 * * *` | Every day at midnight |
| `0 2 * * 1` | Every Monday at 2 AM |
| `0 0 1 * *` | First day of every month |

```bash
kubectl get jobs
kubectl get cronjobs
kubectl describe job data-migration
kubectl logs job/data-migration
```

---

## 15. Ingress

### What is Ingress?

Ingress manages external access to services in a cluster (typically HTTP/HTTPS). It provides load balancing, SSL termination, and name-based virtual hosting.

```
Internet → Ingress Controller → Ingress Rules → Services → Pods
```

> **Prerequisite:** An Ingress Controller must be installed (e.g., NGINX Ingress Controller, Traefik, HAProxy).

### Install NGINX Ingress Controller

```bash
# Using Helm
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm install ingress-nginx ingress-nginx/ingress-nginx

# Or using manifest
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.0/deploy/static/provider/cloud/deploy.yaml
```

### Simple Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: simple-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

### Path-Based Routing

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: path-based-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$2
spec:
  ingressClassName: nginx
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /api(/|$)(.*)
        pathType: ImplementationSpecific
        backend:
          service:
            name: api-service
            port:
              number: 8080
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend-service
            port:
              number: 80
```

### Host-Based Routing (Virtual Hosting)

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: virtual-host-ingress
spec:
  ingressClassName: nginx
  rules:
  - host: app.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api-service
            port:
              number: 8080
```

### TLS / HTTPS

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: tls-ingress
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - myapp.example.com
    secretName: tls-secret          # Must contain tls.crt and tls.key
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: app-service
            port:
              number: 80
```

### Common Ingress Annotations

```yaml
annotations:
  # Rate limiting
  nginx.ingress.kubernetes.io/limit-rps: "10"

  # Body size
  nginx.ingress.kubernetes.io/proxy-body-size: "10m"

  # Timeouts
  nginx.ingress.kubernetes.io/proxy-connect-timeout: "60"
  nginx.ingress.kubernetes.io/proxy-read-timeout: "60"

  # CORS
  nginx.ingress.kubernetes.io/enable-cors: "true"
  nginx.ingress.kubernetes.io/cors-allow-origin: "https://myapp.com"

  # Authentication
  nginx.ingress.kubernetes.io/auth-type: basic
  nginx.ingress.kubernetes.io/auth-secret: basic-auth

  # Redirect
  nginx.ingress.kubernetes.io/permanent-redirect: "https://new-site.com"

  # Canary
  nginx.ingress.kubernetes.io/canary: "true"
  nginx.ingress.kubernetes.io/canary-weight: "20"
```

---

## 16. Network Policies

### What are Network Policies?

Network Policies control traffic flow between Pods, namespaces, and external endpoints. By default, all pods can communicate with all other pods — Network Policies restrict this.

> **Prerequisite:** A CNI plugin that supports Network Policies (Calico, Cilium, Weave Net).

### Default Deny All Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-all-ingress
  namespace: production
spec:
  podSelector: {}           # Applies to all pods in namespace
  policyTypes:
  - Ingress                 # No ingress rules = deny all incoming
```

### Allow Specific Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: allow-frontend-to-backend
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Ingress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    - namespaceSelector:
        matchLabels:
          environment: production
    ports:
    - protocol: TCP
      port: 8080
```

### Allow Egress to Specific Destinations

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: backend-egress
spec:
  podSelector:
    matchLabels:
      app: backend
  policyTypes:
  - Egress
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: database
    ports:
    - protocol: TCP
      port: 5432
  - to:                      # Allow DNS
    - namespaceSelector: {}
    ports:
    - protocol: UDP
      port: 53
    - protocol: TCP
      port: 53
```

### Complete Network Policy Example

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-network-policy
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: api
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    - ipBlock:
        cidr: 10.0.0.0/8
        except:
        - 10.0.1.0/24
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: database
    ports:
    - protocol: TCP
      port: 5432
  - to:
    ports:
    - protocol: UDP
      port: 53
```

---

## 17. RBAC — Role-Based Access Control

### RBAC Components

```
User/Group/ServiceAccount
         │
         ▼
    RoleBinding / ClusterRoleBinding
         │
         ▼
    Role / ClusterRole
         │
         ▼
    Permissions (verbs on resources)
```

| Object | Scope | Description |
|--------|-------|-------------|
| **Role** | Namespace | Defines permissions within a namespace |
| **ClusterRole** | Cluster-wide | Defines permissions cluster-wide |
| **RoleBinding** | Namespace | Binds a Role/ClusterRole to subjects in a namespace |
| **ClusterRoleBinding** | Cluster-wide | Binds a ClusterRole to subjects cluster-wide |

### ServiceAccount

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-service-account
  namespace: production
```

### Role & RoleBinding (Namespace-Scoped)

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: pod-reader
  namespace: development
rules:
- apiGroups: [""]                # "" = core API group
  resources: ["pods", "pods/log"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["apps"]
  resources: ["deployments"]
  verbs: ["get", "list", "watch", "create", "update", "patch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: read-pods
  namespace: development
subjects:
- kind: ServiceAccount
  name: app-service-account
  namespace: development
- kind: User
  name: jane
  apiGroup: rbac.authorization.k8s.io
- kind: Group
  name: developers
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

### ClusterRole & ClusterRoleBinding

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: secret-reader
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "list", "watch"]
- apiGroups: [""]
  resources: ["namespaces"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: read-secrets-global
subjects:
- kind: ServiceAccount
  name: monitoring-sa
  namespace: monitoring
roleRef:
  kind: ClusterRole
  name: secret-reader
  apiGroup: rbac.authorization.k8s.io
```

### RBAC Verbs

| Verb | Description |
|------|-------------|
| `get` | Read a specific resource |
| `list` | List resources |
| `watch` | Watch for changes |
| `create` | Create a resource |
| `update` | Update an existing resource |
| `patch` | Partially update a resource |
| `delete` | Delete a resource |
| `deletecollection` | Delete a collection of resources |

```bash
# Check permissions
kubectl auth can-i create pods --namespace=development
kubectl auth can-i '*' '*'    # Check if you are a cluster admin
kubectl auth can-i get pods --as=system:serviceaccount:default:app-sa
```

---

## 18. Resource Management

### Resource Requests and Limits

```yaml
containers:
- name: app
  image: my-app:1.0
  resources:
    requests:              # Minimum guaranteed resources (used for scheduling)
      cpu: "250m"          # 250 millicores = 0.25 CPU
      memory: "256Mi"      # 256 MiB
    limits:                # Maximum allowed resources
      cpu: "500m"          # Throttled if exceeded
      memory: "512Mi"      # OOMKilled if exceeded
```

### CPU vs Memory Behavior

| Resource | Exceeds Limit | Behavior |
|----------|--------------|----------|
| **CPU** | Pod uses more than limit | Throttled (slowed down) |
| **Memory** | Pod uses more than limit | OOMKilled (terminated) |

### CPU Units

| Value | Meaning |
|-------|---------|
| `1` | 1 vCPU/Core |
| `500m` | 0.5 vCPU (500 millicores) |
| `100m` | 0.1 vCPU |

### Memory Units

| Value | Meaning |
|-------|---------|
| `128Mi` | 128 mebibytes (128 × 1024² bytes) |
| `1Gi` | 1 gibibyte |
| `128M` | 128 megabytes (128 × 1000² bytes) |

### Quality of Service (QoS) Classes

| QoS Class | Condition | Eviction Priority |
|-----------|-----------|-------------------|
| **Guaranteed** | requests == limits for all containers | Last to be evicted |
| **Burstable** | At least one container has request or limit set | Middle priority |
| **BestEffort** | No requests or limits set | First to be evicted |

```yaml
# Guaranteed QoS
resources:
  requests:
    cpu: "500m"
    memory: "256Mi"
  limits:
    cpu: "500m"       # Same as request
    memory: "256Mi"   # Same as request
```

---

## 19. Horizontal Pod Autoscaler (HPA)

### What is HPA?

HPA automatically scales the number of pod replicas based on observed CPU/memory utilization or custom metrics.

### Basic HPA (CPU-based)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: app-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70     # Scale when avg CPU > 70%
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 300    # Wait 5 min before scaling down
      policies:
      - type: Percent
        value: 10
        periodSeconds: 60                # Scale down max 10% per minute
    scaleUp:
      stabilizationWindowSeconds: 0
      policies:
      - type: Percent
        value: 100
        periodSeconds: 15                # Can double capacity every 15s
      - type: Pods
        value: 4
        periodSeconds: 15                # Or add 4 pods every 15s
      selectPolicy: Max
```

```bash
# Imperative creation
kubectl autoscale deployment my-app --min=2 --max=10 --cpu-percent=70

# Monitor HPA
kubectl get hpa
kubectl describe hpa app-hpa
kubectl get hpa app-hpa -w    # Watch
```

### Vertical Pod Autoscaler (VPA)

VPA automatically adjusts CPU and memory requests/limits for pods.

```yaml
apiVersion: autoscaling.k8s.io/v1
kind: VerticalPodAutoscaler
metadata:
  name: my-app-vpa
spec:
  targetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: my-app
  updatePolicy:
    updateMode: "Auto"    # Off, Initial, Recreate, Auto
  resourcePolicy:
    containerPolicies:
    - containerName: app
      minAllowed:
        cpu: "100m"
        memory: "128Mi"
      maxAllowed:
        cpu: "2"
        memory: "2Gi"
```

### Cluster Autoscaler

Automatically adjusts the number of nodes in the cluster.

```yaml
# Typically configured via cloud provider
# AWS example using eksctl:
# eksctl create nodegroup --cluster=my-cluster \
#   --name=workers --nodes-min=2 --nodes-max=10 \
#   --asg-access
```

---

## 20. Helm — Package Manager

### What is Helm?

Helm is the package manager for Kubernetes. It packages Kubernetes manifests into reusable "charts."

### Install Helm

```bash
# macOS
brew install helm

# Linux
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
```

### Helm Concepts

| Concept | Description |
|---------|-------------|
| **Chart** | A package of pre-configured Kubernetes resources |
| **Release** | An instance of a chart running in a cluster |
| **Repository** | A place to store and share charts |
| **Values** | Configuration for a chart |

### Essential Helm Commands

```bash
# Add repositories
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add stable https://charts.helm.sh/stable
helm repo update

# Search charts
helm search repo nginx
helm search hub wordpress        # Search Artifact Hub

# Install a chart
helm install my-release bitnami/nginx
helm install my-release bitnami/nginx --namespace production --create-namespace
helm install my-release bitnami/nginx -f custom-values.yaml
helm install my-release bitnami/nginx --set service.type=NodePort

# List releases
helm list
helm list -A                     # All namespaces

# Upgrade
helm upgrade my-release bitnami/nginx --set replicaCount=3
helm upgrade --install my-release bitnami/nginx    # Install or upgrade

# Rollback
helm rollback my-release 1      # Rollback to revision 1
helm history my-release

# Uninstall
helm uninstall my-release

# Show chart info
helm show values bitnami/nginx   # Default values
helm show chart bitnami/nginx    # Chart metadata

# Template rendering (debug)
helm template my-release bitnami/nginx -f values.yaml
```

### Creating a Helm Chart

```bash
helm create my-app
```

This generates:

```
my-app/
├── Chart.yaml           # Chart metadata
├── values.yaml          # Default configuration values
├── charts/              # Dependent charts
├── templates/           # Kubernetes manifest templates
│   ├── deployment.yaml
│   ├── service.yaml
│   ├── ingress.yaml
│   ├── hpa.yaml
│   ├── serviceaccount.yaml
│   ├── _helpers.tpl     # Template helpers
│   ├── NOTES.txt        # Post-install notes
│   └── tests/
│       └── test-connection.yaml
└── .helmignore
```

### Chart.yaml

```yaml
apiVersion: v2
name: my-app
description: A Helm chart for my application
type: application
version: 0.1.0            # Chart version
appVersion: "1.0.0"       # App version
dependencies:
- name: postgresql
  version: "12.x.x"
  repository: "https://charts.bitnami.com/bitnami"
  condition: postgresql.enabled
```

### values.yaml

```yaml
replicaCount: 3

image:
  repository: my-app
  tag: "1.0.0"
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: true
  hosts:
  - host: myapp.example.com
    paths:
    - path: /
      pathType: Prefix

resources:
  requests:
    cpu: "100m"
    memory: "128Mi"
  limits:
    cpu: "500m"
    memory: "256Mi"

autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

postgresql:
  enabled: true
  auth:
    database: myapp
```

### Template Example

```yaml
# templates/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "my-app.fullname" . }}
  labels:
    {{- include "my-app.labels" . | nindent 4 }}
spec:
  replicas: {{ .Values.replicaCount }}
  selector:
    matchLabels:
      {{- include "my-app.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "my-app.selectorLabels" . | nindent 8 }}
    spec:
      containers:
      - name: {{ .Chart.Name }}
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        imagePullPolicy: {{ .Values.image.pullPolicy }}
        ports:
        - containerPort: {{ .Values.service.port }}
        resources:
          {{- toYaml .Values.resources | nindent 12 }}
```

---

## 21. Probes — Health Checks

### Probe Types

| Probe | Purpose | When it Fails |
|-------|---------|---------------|
| **Liveness** | Is the container alive? | Container is restarted |
| **Readiness** | Is the container ready to serve traffic? | Removed from Service endpoints |
| **Startup** | Has the container started? | Liveness/readiness probes are disabled until it succeeds |

### Probe Mechanisms

| Mechanism | Description |
|-----------|-------------|
| `httpGet` | HTTP GET request, success = 200-399 status code |
| `tcpSocket` | TCP connection check |
| `exec` | Execute a command, success = exit code 0 |
| `grpc` | gRPC health check |

### Complete Probe Configuration

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      containers:
      - name: app
        image: my-app:1.0
        ports:
        - containerPort: 8080

        # Startup probe: for slow-starting containers
        startupProbe:
          httpGet:
            path: /healthz
            port: 8080
          failureThreshold: 30       # 30 × 10s = 5 min max startup time
          periodSeconds: 10

        # Liveness probe: restart if unhealthy
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
            httpHeaders:
            - name: Custom-Header
              value: LivenessCheck
          initialDelaySeconds: 0     # Startup probe handles initial delay
          periodSeconds: 15
          timeoutSeconds: 5
          failureThreshold: 3
          successThreshold: 1

        # Readiness probe: remove from service if not ready
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3
          successThreshold: 1
```

### TCP Socket Probe

```yaml
livenessProbe:
  tcpSocket:
    port: 3306
  initialDelaySeconds: 15
  periodSeconds: 20
```

### Exec Probe

```yaml
livenessProbe:
  exec:
    command:
    - sh
    - -c
    - pg_isready -U postgres
  initialDelaySeconds: 30
  periodSeconds: 10
```

### gRPC Probe

```yaml
livenessProbe:
  grpc:
    port: 50051
    service: my.health.Service    # Optional
  initialDelaySeconds: 10
  periodSeconds: 10
```

### Probe Parameter Reference

| Parameter | Default | Description |
|-----------|---------|-------------|
| `initialDelaySeconds` | 0 | Seconds to wait before first probe |
| `periodSeconds` | 10 | Interval between probes |
| `timeoutSeconds` | 1 | Seconds before probe times out |
| `failureThreshold` | 3 | Consecutive failures to mark unhealthy |
| `successThreshold` | 1 | Consecutive successes to mark healthy |

---

## 22. Taints, Tolerations & Affinity

### Taints & Tolerations

Taints are set on **nodes** to repel pods. Tolerations are set on **pods** to allow scheduling on tainted nodes.

```bash
# Add taint
kubectl taint nodes node1 gpu=true:NoSchedule
kubectl taint nodes node1 dedicated=database:NoExecute

# Remove taint
kubectl taint nodes node1 gpu=true:NoSchedule-
```

| Effect | Description |
|--------|-------------|
| `NoSchedule` | Don't schedule new pods (existing pods stay) |
| `PreferNoSchedule` | Try not to schedule (soft) |
| `NoExecute` | Evict existing pods + don't schedule new ones |

```yaml
# Pod tolerating the taint
apiVersion: v1
kind: Pod
metadata:
  name: gpu-pod
spec:
  tolerations:
  - key: "gpu"
    operator: "Equal"
    value: "true"
    effect: "NoSchedule"
  - key: "dedicated"
    operator: "Exists"      # Matches any value
    effect: "NoExecute"
    tolerationSeconds: 3600  # Tolerate for 1 hour then evict
  containers:
  - name: gpu-app
    image: gpu-app:1.0
```

### Node Affinity

Controls which nodes a pod can be scheduled on based on node labels.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: app
spec:
  affinity:
    nodeAffinity:
      # Hard requirement
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: topology.kubernetes.io/zone
            operator: In
            values:
            - us-east-1a
            - us-east-1b
          - key: node.kubernetes.io/instance-type
            operator: In
            values:
            - m5.xlarge
            - m5.2xlarge

      # Soft preference
      preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 80
        preference:
          matchExpressions:
          - key: disk-type
            operator: In
            values:
            - ssd
  containers:
  - name: app
    image: my-app:1.0
```

### Pod Affinity & Anti-Affinity

Controls pod placement relative to other pods.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web
  template:
    metadata:
      labels:
        app: web
    spec:
      affinity:
        # Run near cache pods
        podAffinity:
          requiredDuringSchedulingIgnoredDuringExecution:
          - labelSelector:
              matchExpressions:
              - key: app
                operator: In
                values:
                - cache
            topologyKey: kubernetes.io/hostname

        # Spread web pods across nodes
        podAntiAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
          - weight: 100
            podAffinityTerm:
              labelSelector:
                matchExpressions:
                - key: app
                  operator: In
                  values:
                  - web
              topologyKey: kubernetes.io/hostname
      containers:
      - name: web
        image: web-app:1.0
```

### Topology Spread Constraints

Distribute pods evenly across failure domains.

```yaml
spec:
  topologySpreadConstraints:
  - maxSkew: 1                    # Max difference between zones
    topologyKey: topology.kubernetes.io/zone
    whenUnsatisfiable: DoNotSchedule
    labelSelector:
      matchLabels:
        app: web
  - maxSkew: 1
    topologyKey: kubernetes.io/hostname
    whenUnsatisfiable: ScheduleAnyway
    labelSelector:
      matchLabels:
        app: web
```

---

## 23. Logging & Monitoring

### Logging

#### Pod Logs

```bash
kubectl logs <pod-name>
kubectl logs <pod-name> -c <container-name>   # Multi-container pod
kubectl logs <pod-name> --previous             # Previous container
kubectl logs -f <pod-name>                     # Stream logs
kubectl logs -l app=nginx                      # By label
kubectl logs --since=1h <pod-name>             # Last hour
kubectl logs --tail=100 <pod-name>             # Last 100 lines
```

#### Centralized Logging Stack (EFK/ELK)

```
Pods → Fluentd/Filebeat (DaemonSet) → Elasticsearch → Kibana
                                            ↓
                                     OpenSearch Dashboard
```

#### Fluentd DaemonSet Example

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd
  namespace: logging
spec:
  selector:
    matchLabels:
      app: fluentd
  template:
    metadata:
      labels:
        app: fluentd
    spec:
      serviceAccountName: fluentd
      tolerations:
      - key: node-role.kubernetes.io/control-plane
        effect: NoSchedule
      containers:
      - name: fluentd
        image: fluent/fluentd-kubernetes-daemonset:v1.16
        env:
        - name: FLUENT_ELASTICSEARCH_HOST
          value: "elasticsearch.logging.svc.cluster.local"
        - name: FLUENT_ELASTICSEARCH_PORT
          value: "9200"
        volumeMounts:
        - name: varlog
          mountPath: /var/log
        - name: dockercontainers
          mountPath: /var/lib/docker/containers
          readOnly: true
      volumes:
      - name: varlog
        hostPath:
          path: /var/log
      - name: dockercontainers
        hostPath:
          path: /var/lib/docker/containers
```

### Monitoring

#### Prometheus + Grafana Stack

```bash
# Install using Helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace
```

This installs:
- **Prometheus** — metrics collection and alerting
- **Grafana** — visualization dashboards
- **Alertmanager** — alert routing and notifications
- **Node Exporter** — hardware/OS metrics
- **kube-state-metrics** — Kubernetes object metrics

#### ServiceMonitor (Prometheus Operator)

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: app-monitor
  namespace: monitoring
spec:
  selector:
    matchLabels:
      app: my-app
  namespaceSelector:
    matchNames:
    - production
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

#### PrometheusRule (Alerting)

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: app-alerts
  namespace: monitoring
spec:
  groups:
  - name: app.rules
    rules:
    - alert: HighErrorRate
      expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "High error rate on {{ $labels.instance }}"
        description: "Error rate is {{ $value }} req/s"

    - alert: PodCrashLooping
      expr: rate(kube_pod_container_status_restarts_total[15m]) > 0
      for: 15m
      labels:
        severity: warning
      annotations:
        summary: "Pod {{ $labels.pod }} is crash looping"
```

### Useful Monitoring Commands

```bash
# Resource usage
kubectl top nodes
kubectl top pods
kubectl top pods --sort-by=cpu
kubectl top pods --sort-by=memory

# Events
kubectl get events --sort-by='.lastTimestamp'
kubectl get events -w     # Watch events
```

---

## 24. Kubernetes on Cloud (EKS, AKS, GKE)

### Amazon EKS (Elastic Kubernetes Service)

```bash
# Install eksctl
brew install eksctl

# Create cluster
eksctl create cluster \
  --name my-cluster \
  --region us-east-1 \
  --version 1.29 \
  --nodegroup-name workers \
  --node-type t3.medium \
  --nodes 3 \
  --nodes-min 2 \
  --nodes-max 5 \
  --managed

# Update kubeconfig
aws eks update-kubeconfig --name my-cluster --region us-east-1

# Add IAM OIDC provider (for IAM roles for service accounts)
eksctl utils associate-iam-oidc-provider --cluster my-cluster --approve

# Install AWS Load Balancer Controller
helm install aws-load-balancer-controller eks/aws-load-balancer-controller \
  -n kube-system --set clusterName=my-cluster

# Install EBS CSI Driver (for persistent volumes)
eksctl create addon --name aws-ebs-csi-driver --cluster my-cluster

# Delete cluster
eksctl delete cluster --name my-cluster
```

### Azure AKS (Azure Kubernetes Service)

```bash
# Create resource group
az group create --name myResourceGroup --location eastus

# Create cluster
az aks create \
  --resource-group myResourceGroup \
  --name myAKSCluster \
  --node-count 3 \
  --node-vm-size Standard_DS2_v2 \
  --enable-managed-identity \
  --generate-ssh-keys

# Get credentials
az aks get-credentials --resource-group myResourceGroup --name myAKSCluster

# Scale
az aks scale --resource-group myResourceGroup --name myAKSCluster --node-count 5

# Delete
az aks delete --resource-group myResourceGroup --name myAKSCluster
```

### Google GKE (Google Kubernetes Engine)

```bash
# Create cluster
gcloud container clusters create my-cluster \
  --zone us-central1-a \
  --num-nodes 3 \
  --machine-type e2-standard-4 \
  --enable-autoscaling --min-nodes 2 --max-nodes 10

# Get credentials
gcloud container clusters get-credentials my-cluster --zone us-central1-a

# Delete
gcloud container clusters delete my-cluster --zone us-central1-a
```

### Cloud Comparison

| Feature | EKS | AKS | GKE |
|---------|-----|-----|-----|
| Control plane cost | ~$73/month | Free | Free (Standard), ~$73/month (Autopilot) |
| Default CNI | VPC-CNI (AWS) | Azure CNI / Kubenet | GKE Dataplane V2 (Cilium) |
| Auto-scaling | Karpenter / Cluster Autoscaler | KEDA / Cluster Autoscaler | GKE Autopilot / Cluster Autoscaler |
| Managed node groups | Yes | Yes (VMSS) | Yes (Node pools) |
| Serverless option | Fargate | Virtual Nodes (ACI) | Autopilot |
| CLI tool | eksctl | az aks | gcloud |

---

## 25. Security Best Practices

### Pod Security Standards

```yaml
# Restricted pod (production recommended)
apiVersion: v1
kind: Pod
metadata:
  name: secure-pod
spec:
  securityContext:
    runAsNonRoot: true
    runAsUser: 1000
    runAsGroup: 3000
    fsGroup: 2000
    seccompProfile:
      type: RuntimeDefault
  containers:
  - name: app
    image: my-app:1.0
    securityContext:
      allowPrivilegeEscalation: false
      readOnlyRootFilesystem: true
      capabilities:
        drop:
        - ALL
    volumeMounts:
    - name: tmp
      mountPath: /tmp
  volumes:
  - name: tmp
    emptyDir: {}
```

### Pod Security Admission (PSA)

```yaml
# Namespace-level enforcement
apiVersion: v1
kind: Namespace
metadata:
  name: production
  labels:
    pod-security.kubernetes.io/enforce: restricted
    pod-security.kubernetes.io/audit: restricted
    pod-security.kubernetes.io/warn: restricted
```

| Level | Description |
|-------|-------------|
| `privileged` | No restrictions |
| `baseline` | Prevents known privilege escalations |
| `restricted` | Heavily restricted, follows security best practices |

### Security Checklist

- [ ] Enable RBAC and follow least-privilege principle
- [ ] Use Network Policies to restrict pod-to-pod traffic
- [ ] Enable Pod Security Admission (restricted mode for production)
- [ ] Scan container images for vulnerabilities (Trivy, Snyk)
- [ ] Don't run containers as root
- [ ] Use read-only root filesystem
- [ ] Drop all capabilities, add only needed ones
- [ ] Enable audit logging
- [ ] Encrypt Secrets at rest (EncryptionConfiguration)
- [ ] Use external secret managers (Vault, AWS Secrets Manager)
- [ ] Regularly rotate credentials and certificates
- [ ] Limit API server access (firewall, private endpoint)
- [ ] Use namespaces to isolate workloads
- [ ] Set resource limits to prevent DoS
- [ ] Keep Kubernetes version up to date
- [ ] Use image pull policies and signed images
- [ ] Enable seccomp and AppArmor profiles

### Encryption at Rest for Secrets

```yaml
# /etc/kubernetes/enc/encryption-config.yaml
apiVersion: apiserver.config.k8s.io/v1
kind: EncryptionConfiguration
resources:
- resources:
  - secrets
  providers:
  - aescbc:
      keys:
      - name: key1
        secret: <base64-encoded-32-byte-key>
  - identity: {}
```

---

## 26. Troubleshooting

### Debugging Flowchart

```
Pod not starting?
├── kubectl describe pod <name>
│   ├── ImagePullBackOff → Wrong image name/tag, no registry access
│   ├── CrashLoopBackOff → App crashes on startup
│   │   └── kubectl logs <pod> --previous
│   ├── Pending → No node has enough resources / node selector mismatch
│   │   └── kubectl describe node <name>
│   ├── CreateContainerConfigError → Bad ConfigMap/Secret reference
│   └── OOMKilled → Memory limit too low
│       └── Increase memory limits
│
Service not working?
├── kubectl get endpoints <svc-name>
│   ├── No endpoints → Labels don't match between Service and Pod
│   └── Has endpoints → Check kube-proxy, NetworkPolicy
├── kubectl exec -it <pod> -- curl <service>:<port>
└── DNS issues → kubectl exec -it <pod> -- nslookup <service>
```

### Common Issues & Solutions

| Symptom | Likely Cause | Solution |
|---------|-------------|---------|
| `ImagePullBackOff` | Wrong image name, private registry | Fix image name, add `imagePullSecrets` |
| `CrashLoopBackOff` | App exits or crashes | Check logs: `kubectl logs <pod> --previous` |
| `Pending` (unschedulable) | Insufficient resources | Scale cluster or adjust resource requests |
| `OOMKilled` | Container exceeds memory limit | Increase memory limit |
| `CreateContainerConfigError` | Missing ConfigMap/Secret | Create the referenced ConfigMap/Secret |
| `Evicted` | Node under resource pressure | Add resource limits, check node capacity |
| Pod can't reach Service | Label mismatch | Verify selector labels match pod labels |
| DNS resolution fails | CoreDNS issues | Check CoreDNS pods in `kube-system` |

### Debugging Commands

```bash
# Pod debugging
kubectl describe pod <pod>
kubectl logs <pod> --previous
kubectl logs <pod> --all-containers
kubectl exec -it <pod> -- /bin/sh
kubectl get pod <pod> -o yaml

# Ephemeral debug container
kubectl debug -it <pod> --image=busybox --target=<container>

# Node debugging
kubectl describe node <node>
kubectl get node <node> -o yaml
kubectl debug node/<node> -it --image=ubuntu

# Service debugging
kubectl get endpoints <service>
kubectl describe service <service>

# DNS debugging
kubectl run dns-test --image=busybox:1.28 --rm -it -- nslookup kubernetes.default

# Network debugging
kubectl run net-test --image=nicolaka/netshoot --rm -it -- bash

# Cluster-wide issues
kubectl get events --sort-by='.lastTimestamp' -A
kubectl get componentstatuses
kubectl cluster-info dump
```

### Resource Debugging

```bash
# Check resource usage
kubectl top nodes
kubectl top pods --containers

# Check resource quotas
kubectl describe resourcequota -n <namespace>

# Check limit ranges
kubectl describe limitrange -n <namespace>

# Find pods without resource limits
kubectl get pods -A -o json | jq '.items[] | select(.spec.containers[].resources.limits == null) | .metadata.name'
```

---

## 27. Real-World Architecture Patterns

### Microservices Architecture

```yaml
# Complete microservices deployment example
---
# Frontend
apiVersion: apps/v1
kind: Deployment
metadata:
  name: frontend
  namespace: production
spec:
  replicas: 3
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      labels:
        app: frontend
    spec:
      containers:
      - name: frontend
        image: frontend:1.0
        ports:
        - containerPort: 3000
        env:
        - name: API_URL
          value: "http://api-gateway:8080"
        resources:
          requests:
            cpu: "100m"
            memory: "128Mi"
          limits:
            cpu: "250m"
            memory: "256Mi"
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: frontend
  namespace: production
spec:
  selector:
    app: frontend
  ports:
  - port: 80
    targetPort: 3000
---
# API Gateway
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: production
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api-gateway
  template:
    metadata:
      labels:
        app: api-gateway
    spec:
      containers:
      - name: gateway
        image: api-gateway:1.0
        ports:
        - containerPort: 8080
        env:
        - name: USER_SERVICE_URL
          value: "http://user-service:8081"
        - name: ORDER_SERVICE_URL
          value: "http://order-service:8082"
        resources:
          requests:
            cpu: "200m"
            memory: "256Mi"
          limits:
            cpu: "500m"
            memory: "512Mi"
---
apiVersion: v1
kind: Service
metadata:
  name: api-gateway
  namespace: production
spec:
  selector:
    app: api-gateway
  ports:
  - port: 8080
    targetPort: 8080
---
# User Service
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
  namespace: production
spec:
  replicas: 2
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: user-service:1.0
        ports:
        - containerPort: 8081
        env:
        - name: DB_HOST
          valueFrom:
            configMapKeyRef:
              name: user-db-config
              key: host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: user-db-secret
              key: password
---
apiVersion: v1
kind: Service
metadata:
  name: user-service
  namespace: production
spec:
  selector:
    app: user-service
  ports:
  - port: 8081
    targetPort: 8081
---
# Ingress
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: main-ingress
  namespace: production
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - myapp.example.com
    secretName: tls-secret
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /api
        pathType: Prefix
        backend:
          service:
            name: api-gateway
            port:
              number: 8080
      - path: /
        pathType: Prefix
        backend:
          service:
            name: frontend
            port:
              number: 80
---
# HPA
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-gateway-hpa
  namespace: production
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-gateway
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
---
# Network Policy
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: api-gateway-policy
  namespace: production
spec:
  podSelector:
    matchLabels:
      app: api-gateway
  policyTypes:
  - Ingress
  - Egress
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - protocol: TCP
      port: 8080
  egress:
  - to:
    - podSelector:
        matchLabels:
          app: user-service
    - podSelector:
        matchLabels:
          app: order-service
    ports:
    - protocol: TCP
      port: 8081
    - protocol: TCP
      port: 8082
```

### CI/CD Pipeline Pattern

```
Developer → Git Push → CI Pipeline → Build Image → Push to Registry
                                                         ↓
Production ← Argo CD / Flux ← Update Manifests ← Image Tag Update
```

### GitOps with ArgoCD

```bash
# Install ArgoCD
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

```yaml
# ArgoCD Application
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: my-app
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/my-org/k8s-manifests.git
    targetRevision: main
    path: production
  destination:
    server: https://kubernetes.default.svc
    namespace: production
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
    - CreateNamespace=true
```

---

## 28. Interview Questions

### Beginner Level

**Q: What is the difference between a Pod and a Container?**
A: A container is a single running process (e.g., Docker container). A Pod is the smallest deployable unit in Kubernetes that can contain one or more containers sharing the same network namespace, IP address, and storage volumes. Pods are the unit of scheduling in Kubernetes.

**Q: What happens when a node goes down?**
A: The Node Controller (part of kube-controller-manager) detects the failure. After a configurable timeout (`--pod-eviction-timeout`, default 5 minutes), pods on the failed node are marked for eviction. If managed by a Deployment/ReplicaSet, new pods are scheduled on healthy nodes to maintain the desired replica count.

**Q: Explain the difference between ClusterIP, NodePort, and LoadBalancer services.**
A:
- **ClusterIP**: Internal only. Accessible within the cluster. Default type.
- **NodePort**: Exposes service on a static port (30000-32767) on every node's IP. Includes ClusterIP.
- **LoadBalancer**: Provisions an external cloud load balancer. Includes NodePort and ClusterIP.

**Q: What is a Namespace and when would you use one?**
A: Namespaces are virtual clusters within a physical cluster. Use them to isolate environments (dev/staging/prod), separate teams, or apply resource quotas and RBAC policies to different groups of resources.

### Intermediate Level

**Q: How does a Deployment perform a rolling update?**
A: A Deployment creates a new ReplicaSet with the updated pod template. It gradually scales up the new ReplicaSet and scales down the old one, respecting `maxSurge` (extra pods during update) and `maxUnavailable` (pods that can be down during update). This ensures zero-downtime updates.

**Q: Explain the difference between ConfigMaps and Secrets.**
A: Both store configuration data, but:
- ConfigMaps store non-sensitive data in plain text
- Secrets store sensitive data base64-encoded (and can be encrypted at rest)
- Secrets are stored in tmpfs on nodes (not written to disk)
- RBAC can restrict Secret access more tightly
- Both can be consumed as environment variables or mounted as files

**Q: What is the difference between a Deployment and a StatefulSet?**
A:
- **Deployment**: For stateless apps. Pods are interchangeable, random names, shared storage.
- **StatefulSet**: For stateful apps. Pods have stable identities (ordered names like `mysql-0`, `mysql-1`), stable network identifiers, and each pod gets its own persistent storage that follows it across rescheduling.

**Q: How does Kubernetes handle resource limits?**
A:
- **CPU limit exceeded**: Pod is throttled (slowed down) but not killed
- **Memory limit exceeded**: Pod is OOMKilled and restarted based on restartPolicy
- **No limits set**: Pod can consume unlimited resources (BestEffort QoS, first to be evicted under pressure)

### Advanced Level

**Q: How would you design a zero-downtime deployment strategy?**
A:
1. Use Rolling Update strategy with `maxUnavailable: 0`
2. Implement readiness probes so traffic only routes to ready pods
3. Use PodDisruptionBudgets to prevent too many pods from being down
4. Add preStop lifecycle hooks for graceful shutdown
5. Configure proper `terminationGracePeriodSeconds`

```yaml
spec:
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    spec:
      terminationGracePeriodSeconds: 60
      containers:
      - name: app
        lifecycle:
          preStop:
            exec:
              command: ["sh", "-c", "sleep 10"]
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
---
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: app-pdb
spec:
  minAvailable: "50%"
  selector:
    matchLabels:
      app: my-app
```

**Q: Explain how Kubernetes networking works.**
A:
1. **Pod-to-Pod**: Every pod gets a unique IP. Pods can communicate directly using pod IPs (flat network, no NAT). Implemented by CNI plugins (Calico, Cilium, Flannel).
2. **Pod-to-Service**: Services provide stable virtual IPs. kube-proxy manages iptables/IPVS rules to route traffic to backend pods.
3. **External-to-Service**: Via NodePort, LoadBalancer, or Ingress.
4. **DNS**: CoreDNS resolves service names to ClusterIPs. Format: `<service>.<namespace>.svc.cluster.local`

**Q: How would you secure a Kubernetes cluster?**
A:
1. **RBAC**: Least-privilege access for users and service accounts
2. **Network Policies**: Segment traffic between namespaces and pods
3. **Pod Security**: Run as non-root, drop capabilities, read-only filesystem
4. **Secrets Management**: Encrypt at rest, use external secret managers
5. **Image Security**: Scan images, use signed images, restrict registries
6. **Audit Logging**: Enable API server audit logs
7. **Updates**: Keep K8s and node OS up to date
8. **API Server**: Restrict access via firewall rules, use private endpoints
9. **etcd**: Encrypt data, restrict access, regular backups

**Q: What is a Service Mesh and when would you use one?**
A: A service mesh (e.g., Istio, Linkerd) adds a sidecar proxy to each pod to handle service-to-service communication. It provides:
- **mTLS**: Automatic encryption between services
- **Traffic management**: Canary deployments, traffic splitting, retries, circuit breaking
- **Observability**: Distributed tracing, metrics, access logs
- **Policy enforcement**: Rate limiting, access control

Use when: you have many microservices, need fine-grained traffic control, require mTLS everywhere, or need advanced observability without modifying application code.

---

## Quick Reference Cheat Sheet

```bash
# Cluster
kubectl cluster-info
kubectl get nodes -o wide
kubectl get all -A

# Pods
kubectl run nginx --image=nginx --port=80
kubectl get pods -o wide
kubectl describe pod <name>
kubectl logs <pod> [-f] [--previous]
kubectl exec -it <pod> -- /bin/sh
kubectl delete pod <name>

# Deployments
kubectl create deployment nginx --image=nginx --replicas=3
kubectl scale deployment nginx --replicas=5
kubectl set image deployment/nginx nginx=nginx:1.26
kubectl rollout status deployment/nginx
kubectl rollout undo deployment/nginx
kubectl rollout restart deployment/nginx

# Services
kubectl expose deployment nginx --port=80 --type=NodePort
kubectl get svc
kubectl get endpoints

# ConfigMaps & Secrets
kubectl create configmap my-config --from-literal=key=value
kubectl create secret generic my-secret --from-literal=password=s3cret

# Namespaces
kubectl create namespace dev
kubectl get pods -n dev
kubectl config set-context --current --namespace=dev

# Debugging
kubectl describe <resource> <name>
kubectl get events --sort-by='.lastTimestamp'
kubectl top pods
kubectl top nodes
kubectl auth can-i <verb> <resource>

# Dry run + YAML generation
kubectl run nginx --image=nginx --dry-run=client -o yaml > pod.yaml
kubectl create deployment nginx --image=nginx --dry-run=client -o yaml > deploy.yaml
kubectl expose deployment nginx --port=80 --dry-run=client -o yaml > svc.yaml
```

---

*Last updated: March 2026*
