# The Complete Docker Guide

## Table of Contents

1. [Introduction](#1-introduction)
2. [Architecture & Core Concepts](#2-architecture--core-concepts)
3. [Installation & Setup](#3-installation--setup)
4. [Docker CLI — Essential Commands](#4-docker-cli--essential-commands)
5. [Docker Images](#5-docker-images)
6. [Dockerfile — Building Images](#6-dockerfile--building-images)
7. [Docker Containers](#7-docker-containers)
8. [Docker Networking](#8-docker-networking)
9. [Docker Volumes & Storage](#9-docker-volumes--storage)
10. [Docker Compose](#10-docker-compose)
11. [Docker Registry](#11-docker-registry)
12. [Multi-Stage Builds](#12-multi-stage-builds)
13. [Docker Build Cache & Optimization](#13-docker-build-cache--optimization)
14. [Docker Security](#14-docker-security)
15. [Docker Logging & Monitoring](#15-docker-logging--monitoring)
16. [Docker in CI/CD](#16-docker-in-cicd)
17. [Docker Troubleshooting](#17-docker-troubleshooting)
18. [Docker with Spring Boot — Complete Guide](#18-docker-with-spring-boot--complete-guide)
19. [Interview Questions](#19-interview-questions)
20. [Quick Reference Cheat Sheet](#20-quick-reference-cheat-sheet)

---

## 1. Introduction

### What is Docker?

Docker is an open-source platform for building, shipping, and running applications in isolated environments called **containers**. Containers package an application with all its dependencies, libraries, and configuration so it runs consistently across any environment.

### Containers vs Virtual Machines

```
┌────────────────────────────────────────────────────────────────┐
│          Virtual Machines              Containers               │
│                                                                │
│  ┌──────┐ ┌──────┐ ┌──────┐   ┌──────┐ ┌──────┐ ┌──────┐    │
│  │ App A│ │ App B│ │ App C│   │ App A│ │ App B│ │ App C│    │
│  ├──────┤ ├──────┤ ├──────┤   ├──────┤ ├──────┤ ├──────┤    │
│  │ Bins │ │ Bins │ │ Bins │   │ Bins │ │ Bins │ │ Bins │    │
│  │ Libs │ │ Libs │ │ Libs │   │ Libs │ │ Libs │ │ Libs │    │
│  ├──────┤ ├──────┤ ├──────┤   └──┬───┘ └──┬───┘ └──┬───┘    │
│  │Guest │ │Guest │ │Guest │      │        │        │         │
│  │  OS  │ │  OS  │ │  OS  │   ┌──▼────────▼────────▼───┐     │
│  └──┬───┘ └──┬───┘ └──┬───┘   │    Container Runtime   │     │
│  ┌──▼────────▼────────▼───┐   │      (Docker Engine)    │     │
│  │      Hypervisor        │   └────────────┬────────────┘     │
│  └────────────┬───────────┘                │                  │
│  ┌────────────▼───────────┐   ┌────────────▼────────────┐     │
│  │       Host OS          │   │       Host OS            │     │
│  └────────────────────────┘   └──────────────────────────┘     │
│  ┌────────────────────────┐   ┌──────────────────────────┐     │
│  │      Hardware          │   │       Hardware           │     │
│  └────────────────────────┘   └──────────────────────────┘     │
└────────────────────────────────────────────────────────────────┘
```

| Feature | Container | Virtual Machine |
|---------|-----------|-----------------|
| Boot time | Seconds | Minutes |
| Size | MBs | GBs |
| Performance | Near-native | Overhead from hypervisor |
| OS | Shares host kernel | Full guest OS |
| Isolation | Process-level | Hardware-level |
| Portability | Highly portable | Less portable |
| Density | 100s per host | 10s per host |
| Use case | Microservices, CI/CD | Legacy apps, OS-level isolation |

### Why Docker?

| Problem | Docker Solution |
|---------|----------------|
| "Works on my machine" | Identical environment everywhere via images |
| Dependency conflicts | Isolated containers with their own dependencies |
| Slow environment setup | `docker run` starts in seconds |
| Resource waste | Lightweight containers share host kernel |
| Complex deployments | Build once, deploy anywhere |
| Inconsistent environments | Same image in dev, staging, and production |

---

## 2. Architecture & Core Concepts

### Docker Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                      Docker Architecture                      │
│                                                              │
│  ┌─── Docker Client ───┐     ┌─── Docker Host (Daemon) ───┐ │
│  │                      │     │                             │ │
│  │  docker build         │────▶│  Docker Daemon (dockerd)   │ │
│  │  docker pull          │     │       │                    │ │
│  │  docker run           │     │  ┌────▼────┐  ┌────────┐  │ │
│  │  docker push          │     │  │ Images  │  │Containers│ │ │
│  │  docker compose       │     │  └─────────┘  └────────┘  │ │
│  │                      │     │                             │ │
│  └──────────────────────┘     │  ┌─────────┐  ┌────────┐  │ │
│                               │  │ Volumes │  │Networks│  │ │
│                               │  └─────────┘  └────────┘  │ │
│                               └──────────┬──────────────────┘ │
│                                          │                    │
│                               ┌──────────▼──────────────┐    │
│                               │    Docker Registry       │    │
│                               │  (Docker Hub, ECR, GCR)  │    │
│                               └─────────────────────────┘    │
└──────────────────────────────────────────────────────────────┘
```

### Core Components

| Component | Description |
|-----------|-------------|
| **Docker Daemon** (`dockerd`) | Background service that manages images, containers, networks, and volumes |
| **Docker Client** (`docker`) | CLI tool that sends commands to the daemon via REST API |
| **Docker Registry** | Stores Docker images (Docker Hub, Amazon ECR, Google GCR, GitHub GHCR) |
| **Docker Image** | Read-only template with instructions to create a container (layered filesystem) |
| **Docker Container** | Runnable instance of an image — an isolated process with its own filesystem |
| **Dockerfile** | Text file with instructions to build a Docker image |
| **Docker Compose** | Tool for defining and running multi-container applications |

### Image Layers

```
┌──────────────────────────┐
│   Writable Container     │  ← Container layer (R/W)
│         Layer            │
├──────────────────────────┤
│   COPY app.jar /app/     │  ← Layer 5 (R/O)
├──────────────────────────┤
│   RUN apt-get install    │  ← Layer 4 (R/O)
├──────────────────────────┤
│   RUN apt-get update     │  ← Layer 3 (R/O)
├──────────────────────────┤
│   ENV JAVA_HOME=/usr/lib │  ← Layer 2 (R/O)
├──────────────────────────┤
│   Base Image (Ubuntu)    │  ← Layer 1 (R/O)
└──────────────────────────┘
```

Each Dockerfile instruction creates a new layer. Layers are cached and shared between images, saving disk space and build time.

---

## 3. Installation & Setup

### Install Docker

#### macOS

```bash
# Docker Desktop
brew install --cask docker

# Or download from https://www.docker.com/products/docker-desktop
```

#### Linux (Ubuntu/Debian)

```bash
# Remove old versions
sudo apt-get remove docker docker-engine docker.io containerd runc

# Install prerequisites
sudo apt-get update
sudo apt-get install ca-certificates curl gnupg lsb-release

# Add Docker's official GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg

# Set up the repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Run Docker without sudo
sudo usermod -aG docker $USER
newgrp docker
```

#### Windows

Download Docker Desktop from https://www.docker.com/products/docker-desktop. Requires WSL 2 backend.

### Verify Installation

```bash
docker --version
docker compose version
docker run hello-world
docker info
```

### Docker Desktop Settings (Recommended)

```
Resources:
  CPUs: 4 (minimum for development)
  Memory: 8 GB (minimum for Spring Boot + DB)
  Disk: 60 GB

Features:
  ✓ Use Docker Compose V2
  ✓ Enable Kubernetes (optional)
```

---

## 4. Docker CLI — Essential Commands

### Image Commands

```bash
# Pull an image
docker pull nginx
docker pull nginx:1.25              # Specific tag
docker pull ubuntu:22.04

# List images
docker images
docker images -a                     # Include intermediate images
docker images --filter "dangling=true"

# Remove images
docker rmi nginx:1.25
docker rmi $(docker images -q)      # Remove all images
docker image prune                   # Remove dangling images
docker image prune -a                # Remove all unused images

# Image info
docker inspect nginx
docker history nginx                 # Show layer history

# Tag an image
docker tag my-app:latest my-app:1.0
docker tag my-app:latest registry.example.com/my-app:1.0

# Save/Load images
docker save -o my-app.tar my-app:1.0
docker load -i my-app.tar

# Search Docker Hub
docker search nginx
```

### Container Commands

```bash
# Run a container
docker run nginx                            # Foreground
docker run -d nginx                         # Detached (background)
docker run -d --name my-nginx nginx         # Named container
docker run -d -p 8080:80 nginx              # Port mapping host:container
docker run -d -P nginx                      # Map all exposed ports to random host ports
docker run --rm nginx                       # Auto-remove when stopped
docker run -it ubuntu /bin/bash             # Interactive terminal

# Environment variables
docker run -e "DB_HOST=localhost" -e "DB_PORT=5432" my-app
docker run --env-file .env my-app

# Resource limits
docker run -d --memory="512m" --cpus="1.0" my-app

# List containers
docker ps                                   # Running containers
docker ps -a                                # All containers (including stopped)
docker ps -q                                # Only IDs
docker ps --filter "status=exited"

# Container operations
docker stop <container>                     # Graceful stop (SIGTERM → SIGKILL)
docker start <container>                    # Start stopped container
docker restart <container>
docker kill <container>                     # Force stop (SIGKILL)
docker pause <container>                    # Freeze processes
docker unpause <container>

# Remove containers
docker rm <container>
docker rm -f <container>                    # Force remove running container
docker container prune                      # Remove all stopped containers

# Container info
docker inspect <container>
docker stats                                # Live resource usage
docker stats <container>
docker top <container>                      # Running processes
docker port <container>                     # Port mappings
docker diff <container>                     # Filesystem changes

# Copy files
docker cp file.txt <container>:/path/       # Host → container
docker cp <container>:/path/file.txt .      # Container → host

# Execute commands
docker exec -it <container> /bin/bash       # Shell into running container
docker exec <container> ls /app             # Run a command
docker exec -u root <container> whoami      # Run as specific user

# Logs
docker logs <container>
docker logs -f <container>                  # Follow
docker logs --tail 100 <container>          # Last 100 lines
docker logs --since 1h <container>          # Last hour
docker logs --timestamps <container>

# Create image from container
docker commit <container> my-new-image:1.0
```

### System Commands

```bash
# System-wide info
docker system info
docker system df                             # Disk usage

# Clean up everything
docker system prune                          # Remove stopped containers, dangling images, unused networks
docker system prune -a                       # Also remove all unused images
docker system prune -a --volumes             # Also remove volumes (CAUTION: data loss)
```

---

## 5. Docker Images

### Understanding Image Names

```
registry.example.com/namespace/repository:tag
│                     │          │          │
│                     │          │          └─ Version tag (default: latest)
│                     │          └─ Image name
│                     └─ Organization/user
└─ Registry (default: docker.io)

Examples:
  nginx                          → docker.io/library/nginx:latest
  myuser/my-app:1.0              → docker.io/myuser/my-app:1.0
  ghcr.io/org/app:v2.1           → GitHub Container Registry
  123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:latest  → AWS ECR
```

### Image Tags Best Practices

| Practice | Example | Why |
|----------|---------|-----|
| Use specific versions | `nginx:1.25.3` | Reproducible builds |
| Avoid `latest` in prod | `my-app:1.2.3` | `latest` is mutable |
| Use semantic versioning | `my-app:2.1.0` | Clear version tracking |
| Use Git SHA for CI/CD | `my-app:a1b2c3d` | Traceable to commit |
| Multi-tag releases | `my-app:1.2.3`, `my-app:1.2`, `my-app:1` | Flexible pinning |

### Base Image Selection

| Base Image | Size | Use Case |
|------------|------|----------|
| `scratch` | 0 MB | Static binaries (Go, Rust) |
| `alpine` | ~5 MB | Minimal Linux |
| `distroless` | ~20 MB | Security-focused, no shell |
| `debian-slim` | ~80 MB | When you need apt |
| `ubuntu` | ~78 MB | Full Ubuntu |
| `eclipse-temurin:21-jre-alpine` | ~120 MB | Java/Spring Boot apps |
| `eclipse-temurin:21-jre-jammy` | ~230 MB | Java apps needing glibc |

---

## 6. Dockerfile — Building Images

### Dockerfile Instructions Reference

| Instruction | Description | Example |
|------------|-------------|---------|
| `FROM` | Base image | `FROM eclipse-temurin:21-jre-alpine` |
| `WORKDIR` | Set working directory | `WORKDIR /app` |
| `COPY` | Copy files from build context | `COPY target/app.jar app.jar` |
| `ADD` | Copy files (supports URLs, auto-extract tar) | `ADD https://example.com/file.tar.gz /tmp/` |
| `RUN` | Execute command during build | `RUN apt-get update && apt-get install -y curl` |
| `CMD` | Default command when container starts | `CMD ["java", "-jar", "app.jar"]` |
| `ENTRYPOINT` | Main executable (not easily overridden) | `ENTRYPOINT ["java", "-jar"]` |
| `ENV` | Set environment variable | `ENV JAVA_OPTS="-Xmx512m"` |
| `ARG` | Build-time variable | `ARG JAR_FILE=app.jar` |
| `EXPOSE` | Document which port the app listens on | `EXPOSE 8080` |
| `VOLUME` | Create mount point | `VOLUME /data` |
| `USER` | Set user for subsequent instructions | `USER 1001` |
| `LABEL` | Add metadata | `LABEL version="1.0"` |
| `HEALTHCHECK` | Container health check | `HEALTHCHECK CMD curl -f http://localhost:8080/health` |
| `STOPSIGNAL` | Set system call signal for stop | `STOPSIGNAL SIGTERM` |

### CMD vs ENTRYPOINT

```dockerfile
# CMD — defines default command (can be overridden at runtime)
CMD ["nginx", "-g", "daemon off;"]
# docker run my-image                    → runs: nginx -g daemon off;
# docker run my-image /bin/sh            → runs: /bin/sh (CMD overridden)

# ENTRYPOINT — defines the main executable (not easily overridden)
ENTRYPOINT ["java", "-jar", "app.jar"]
# docker run my-image                    → runs: java -jar app.jar
# docker run my-image --server.port=9090 → runs: java -jar app.jar --server.port=9090

# Combined — ENTRYPOINT + CMD (CMD provides default args)
ENTRYPOINT ["java", "-jar"]
CMD ["app.jar"]
# docker run my-image                    → runs: java -jar app.jar
# docker run my-image other.jar          → runs: java -jar other.jar
```

### Shell Form vs Exec Form

```dockerfile
# Exec form (preferred) — runs directly, signals handled properly
CMD ["java", "-jar", "app.jar"]
ENTRYPOINT ["python", "app.py"]
RUN ["apt-get", "install", "-y", "curl"]

# Shell form — runs through /bin/sh -c, variable substitution works
CMD java -jar app.jar
RUN echo "Hello $NAME"
```

> Always use **exec form** for `CMD` and `ENTRYPOINT` so the process receives OS signals (SIGTERM for graceful shutdown).

### Basic Dockerfile Examples

#### Java / Spring Boot

```dockerfile
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY target/my-app-1.0.0.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

#### Node.js

```dockerfile
FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE 3000
CMD ["node", "server.js"]
```

#### Python

```dockerfile
FROM python:3.12-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE 8000
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
```

#### Go

```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server

FROM scratch
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### .dockerignore

Prevent unnecessary files from entering the build context.

```
# .dockerignore
.git
.gitignore
.dockerignore
Dockerfile
docker-compose*.yml
*.md
LICENSE

# IDE
.idea/
.vscode/
*.iml

# Build artifacts
target/
build/
node_modules/
dist/
__pycache__/
*.pyc

# Environment
.env
.env.local
*.log

# OS
.DS_Store
Thumbs.db
```

### Build Commands

```bash
# Basic build
docker build -t my-app:1.0 .

# Build with specific Dockerfile
docker build -f Dockerfile.prod -t my-app:1.0 .

# Build with build args
docker build --build-arg JAR_FILE=app.jar -t my-app:1.0 .

# Build with no cache
docker build --no-cache -t my-app:1.0 .

# Build for specific platform
docker build --platform linux/amd64 -t my-app:1.0 .

# Multi-platform build
docker buildx build --platform linux/amd64,linux/arm64 -t my-app:1.0 --push .

# Build and show output
docker build --progress=plain -t my-app:1.0 .
```

---

## 7. Docker Containers

### Container Lifecycle

```
Created → Running → Paused → Running → Stopped → Removed
   │         │                   │          │
   │         └── docker pause ───┘          │
   │         └── docker stop ───────────────┘
   └── docker start ────────────────────────┘
```

### Restart Policies

```bash
docker run -d --restart=no         my-app    # Default: never restart
docker run -d --restart=always     my-app    # Always restart (even on daemon restart)
docker run -d --restart=unless-stopped my-app # Restart unless explicitly stopped
docker run -d --restart=on-failure my-app     # Restart only on non-zero exit code
docker run -d --restart=on-failure:5 my-app   # Max 5 restart attempts
```

| Policy | On Crash | On `docker stop` | On Daemon Restart |
|--------|----------|-------------------|-------------------|
| `no` | No restart | — | No restart |
| `always` | Restart | — | Restart |
| `unless-stopped` | Restart | — | No restart |
| `on-failure` | Restart | — | No restart |

### Resource Limits

```bash
# Memory
docker run -d --memory="512m" my-app                # Hard limit
docker run -d --memory="512m" --memory-swap="1g" my-app  # Memory + swap
docker run -d --memory="512m" --oom-kill-disable my-app   # Prevent OOM kill

# CPU
docker run -d --cpus="1.5" my-app                   # 1.5 CPUs
docker run -d --cpu-shares=512 my-app                # Relative weight (default 1024)
docker run -d --cpuset-cpus="0,1" my-app             # Pin to specific CPUs

# Combined
docker run -d --memory="512m" --cpus="1.0" \
  --pids-limit=100 my-app
```

### HEALTHCHECK

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --retries=3 --start-period=60s \
  CMD curl -f http://localhost:8080/actuator/health || exit 1
```

```bash
# Check container health
docker inspect --format='{{.State.Health.Status}}' <container>
docker ps   # Shows health status in STATUS column
```

---

## 8. Docker Networking

### Network Drivers

| Driver | Description | Use Case |
|--------|-------------|----------|
| `bridge` | Default. Isolated network on the host. | Single-host container communication |
| `host` | Container shares host's network stack | Performance-sensitive apps |
| `none` | No networking | Isolated containers |
| `overlay` | Multi-host networking (Swarm) | Distributed applications |
| `macvlan` | Assign MAC address to container | Legacy apps needing direct network access |

### Network Commands

```bash
# List networks
docker network ls

# Create network
docker network create my-network
docker network create --driver bridge --subnet 172.20.0.0/16 my-network

# Connect/disconnect containers
docker network connect my-network <container>
docker network disconnect my-network <container>

# Inspect network
docker network inspect my-network

# Remove network
docker network rm my-network
docker network prune                # Remove unused networks
```

### Bridge Network (Default)

Containers on the same bridge network can communicate by container name.

```bash
# Create a custom bridge network
docker network create app-network

# Run containers on the same network
docker run -d --name db --network app-network \
  -e POSTGRES_PASSWORD=secret postgres:16

docker run -d --name api --network app-network \
  -e DB_HOST=db -e DB_PORT=5432 my-api:1.0

# "api" container can reach "db" by name: postgres://db:5432
```

> **Important:** The default `bridge` network does NOT support container name DNS resolution. Always create a custom bridge network.

### Port Mapping

```bash
# Map specific port
docker run -d -p 8080:80 nginx                # host:container
docker run -d -p 127.0.0.1:8080:80 nginx      # Bind to specific interface

# Map multiple ports
docker run -d -p 8080:80 -p 8443:443 nginx

# Map all exposed ports to random host ports
docker run -d -P nginx

# Map UDP port
docker run -d -p 5000:5000/udp my-app

# Map port range
docker run -d -p 8080-8090:8080-8090 my-app
```

### Host Network

Container uses the host's network directly. No port mapping needed.

```bash
docker run -d --network host nginx
# Nginx is accessible at host's port 80 directly
```

### DNS Resolution Inside Containers

```bash
# Containers on the same custom network resolve each other by name
docker exec api ping db           # Works on custom bridge network

# Custom DNS
docker run -d --dns 8.8.8.8 --dns-search example.com my-app
```

---

## 9. Docker Volumes & Storage

### Storage Types

```
┌─────────────────────────────────────────────────────────────┐
│                    Docker Storage Types                       │
│                                                              │
│  Named Volumes        Bind Mounts           tmpfs Mounts     │
│  ┌──────────┐         ┌──────────┐          ┌──────────┐    │
│  │ Docker   │         │ Host     │          │  Memory  │    │
│  │ managed  │         │ filesystem│          │  (RAM)   │    │
│  │ /var/lib/│         │ any path │          │          │    │
│  │ docker/  │         │          │          │          │    │
│  │ volumes/ │         │          │          │          │    │
│  └──────────┘         └──────────┘          └──────────┘    │
│  Persistent           Development           Sensitive data   │
│  Production-ready     Live reload           Non-persistent   │
└─────────────────────────────────────────────────────────────┘
```

### Named Volumes (Recommended for Production)

```bash
# Create volume
docker volume create my-data

# Run with volume
docker run -d --name db \
  -v my-data:/var/lib/postgresql/data \
  postgres:16

# List volumes
docker volume ls

# Inspect volume
docker volume inspect my-data

# Remove volume
docker volume rm my-data
docker volume prune              # Remove all unused volumes

# Backup a volume
docker run --rm -v my-data:/data -v $(pwd):/backup \
  busybox tar czf /backup/my-data-backup.tar.gz /data

# Restore a volume
docker run --rm -v my-data:/data -v $(pwd):/backup \
  busybox tar xzf /backup/my-data-backup.tar.gz -C /
```

### Bind Mounts (Development)

```bash
# Mount a host directory into the container
docker run -d -v /host/path:/container/path my-app
docker run -d -v $(pwd)/src:/app/src my-app          # Current directory

# Read-only mount
docker run -d -v $(pwd)/config:/app/config:ro my-app

# Using --mount (more explicit)
docker run -d \
  --mount type=bind,source=$(pwd)/src,target=/app/src \
  my-app
```

### tmpfs Mounts

```bash
# In-memory filesystem (Linux only)
docker run -d --tmpfs /tmp:rw,size=100m my-app

docker run -d \
  --mount type=tmpfs,destination=/tmp,tmpfs-size=100m \
  my-app
```

### Volume Driver Plugins

```bash
# Use a specific volume driver (e.g., for NFS, cloud storage)
docker volume create --driver local \
  --opt type=nfs \
  --opt o=addr=192.168.1.100,rw \
  --opt device=:/shared \
  nfs-volume
```

### Volume vs Bind Mount Summary

| Feature | Named Volume | Bind Mount |
|---------|-------------|------------|
| Managed by Docker | Yes | No |
| Location | `/var/lib/docker/volumes/` | Anywhere on host |
| Pre-populated | Yes (from container) | No |
| Backup support | Via Docker CLI | Direct filesystem access |
| Best for | Production data persistence | Development live reload |
| Docker Compose | `volumes:` top-level key | Direct path mapping |

---

## 10. Docker Compose

### What is Docker Compose?

Docker Compose defines and manages multi-container applications using a YAML file. It simplifies running complex applications with multiple services, networks, and volumes.

### Compose File Structure

```yaml
# docker-compose.yml (or compose.yaml)

services:          # Container definitions
  web:
    ...
  db:
    ...
  cache:
    ...

volumes:           # Named volumes
  db-data:
  cache-data:

networks:          # Custom networks
  frontend:
  backend:

configs:           # Configuration files
  nginx-config:

secrets:           # Sensitive data
  db-password:
```

### Complete Docker Compose Example

```yaml
version: "3.9"

services:
  # Spring Boot Application
  api:
    build:
      context: .
      dockerfile: Dockerfile
      args:
        JAR_FILE: target/app.jar
    image: my-api:latest
    container_name: api
    ports:
      - "8080:8080"
    environment:
      SPRING_PROFILES_ACTIVE: docker
      SPRING_DATASOURCE_URL: jdbc:postgresql://db:5432/myapp
      SPRING_DATASOURCE_USERNAME: appuser
      SPRING_DATASOURCE_PASSWORD_FILE: /run/secrets/db_password
      SPRING_REDIS_HOST: cache
    depends_on:
      db:
        condition: service_healthy
      cache:
        condition: service_started
    networks:
      - backend
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/actuator/health"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 60s
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 512M
        reservations:
          cpus: "0.5"
          memory: 256M
    secrets:
      - db_password

  # PostgreSQL Database
  db:
    image: postgres:16-alpine
    container_name: postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: appuser
      POSTGRES_PASSWORD_FILE: /run/secrets/db_password
    volumes:
      - db-data:/var/lib/postgresql/data
      - ./init-scripts:/docker-entrypoint-initdb.d:ro
    networks:
      - backend
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U appuser -d myapp"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s
    secrets:
      - db_password

  # Redis Cache
  cache:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    command: redis-server --appendonly yes --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - cache-data:/data
    networks:
      - backend
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3

  # Nginx Reverse Proxy
  nginx:
    image: nginx:1.25-alpine
    container_name: nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/certs:/etc/nginx/certs:ro
    depends_on:
      api:
        condition: service_healthy
    networks:
      - frontend
      - backend
    restart: unless-stopped

  # pgAdmin (Development only)
  pgadmin:
    image: dpage/pgadmin4:latest
    container_name: pgadmin
    ports:
      - "5050:80"
    environment:
      PGADMIN_DEFAULT_EMAIL: admin@example.com
      PGADMIN_DEFAULT_PASSWORD: admin
    depends_on:
      - db
    networks:
      - backend
    profiles:
      - dev                         # Only starts with --profile dev

volumes:
  db-data:
    driver: local
  cache-data:
    driver: local

networks:
  frontend:
    driver: bridge
  backend:
    driver: bridge
    internal: false

secrets:
  db_password:
    file: ./secrets/db_password.txt
```

### Docker Compose Commands

```bash
# Start services
docker compose up                      # Foreground
docker compose up -d                   # Detached
docker compose up -d --build           # Rebuild images before starting
docker compose up -d api db            # Start specific services

# Stop services
docker compose stop                    # Stop (keep containers)
docker compose down                    # Stop and remove containers, networks
docker compose down -v                 # Also remove volumes (CAUTION: data loss)
docker compose down --rmi all          # Also remove images

# Manage services
docker compose ps                      # List services
docker compose logs                    # View logs
docker compose logs -f api             # Follow logs for specific service
docker compose exec api /bin/sh        # Shell into running service
docker compose run api bash            # Run one-off command

# Scale services
docker compose up -d --scale api=3

# Build
docker compose build                   # Build all services
docker compose build api               # Build specific service
docker compose build --no-cache        # Build without cache

# Config validation
docker compose config                  # Validate and show resolved config

# Profiles
docker compose --profile dev up -d     # Start with dev profile
```

### Compose Override Files

```bash
# Docker Compose automatically merges:
# docker-compose.yml + docker-compose.override.yml

# Or specify explicitly:
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

```yaml
# docker-compose.override.yml (development overrides)
services:
  api:
    build:
      context: .
      target: development
    volumes:
      - ./src:/app/src                  # Hot reload
    environment:
      SPRING_PROFILES_ACTIVE: dev
      JAVA_TOOL_OPTIONS: "-agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005"
    ports:
      - "8080:8080"
      - "5005:5005"                     # Debug port
```

```yaml
# docker-compose.prod.yml
services:
  api:
    image: registry.example.com/my-api:${VERSION}
    environment:
      SPRING_PROFILES_ACTIVE: prod
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: "2.0"
          memory: 1G
```

### Environment Variables in Compose

```yaml
services:
  api:
    # Method 1: Inline
    environment:
      DB_HOST: postgres
      DB_PORT: 5432

    # Method 2: From .env file
    env_file:
      - .env
      - .env.local

    # Method 3: Variable substitution (from shell or .env)
    image: my-app:${VERSION:-latest}
    environment:
      DB_HOST: ${DB_HOST:?DB_HOST is required}
```

```bash
# .env file (auto-loaded by Compose)
VERSION=1.2.3
DB_HOST=postgres
DB_PORT=5432
POSTGRES_PASSWORD=secret
```

---

## 11. Docker Registry

### Docker Hub

```bash
# Login
docker login
docker login -u username

# Push
docker tag my-app:1.0 username/my-app:1.0
docker push username/my-app:1.0

# Pull
docker pull username/my-app:1.0
```

### Private Registry

```bash
# Run a local registry
docker run -d -p 5000:5000 --name registry \
  -v registry-data:/var/lib/registry \
  registry:2

# Push to local registry
docker tag my-app:1.0 localhost:5000/my-app:1.0
docker push localhost:5000/my-app:1.0

# Pull from local registry
docker pull localhost:5000/my-app:1.0

# List images in local registry
curl http://localhost:5000/v2/_catalog
curl http://localhost:5000/v2/my-app/tags/list
```

### Cloud Registries

```bash
# AWS ECR
aws ecr get-login-password --region us-east-1 | \
  docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com
docker tag my-app:1.0 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:1.0
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/my-app:1.0

# Google GCR / Artifact Registry
gcloud auth configure-docker us-central1-docker.pkg.dev
docker tag my-app:1.0 us-central1-docker.pkg.dev/my-project/my-repo/my-app:1.0
docker push us-central1-docker.pkg.dev/my-project/my-repo/my-app:1.0

# Azure ACR
az acr login --name myregistry
docker tag my-app:1.0 myregistry.azurecr.io/my-app:1.0
docker push myregistry.azurecr.io/my-app:1.0

# GitHub Container Registry
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
docker tag my-app:1.0 ghcr.io/username/my-app:1.0
docker push ghcr.io/username/my-app:1.0
```

---

## 12. Multi-Stage Builds

### Why Multi-Stage Builds?

Multi-stage builds separate the build environment from the runtime environment, resulting in smaller, more secure production images.

```
┌─── Stage 1: Build ──────────┐     ┌─── Stage 2: Runtime ──────┐
│                              │     │                            │
│  Full JDK + Maven/Gradle     │     │  JRE only (no build tools)│
│  Source code                 │     │  Application JAR only      │
│  Dependencies                │     │                            │
│  Compiled artifacts   ───────┼────▶│  Final image: ~150 MB      │
│                              │     │                            │
│  Image size: ~800 MB         │     │                            │
│  (discarded)                 │     │                            │
└──────────────────────────────┘     └────────────────────────────┘
```

### Java / Spring Boot Multi-Stage

```dockerfile
# Stage 1: Build
FROM eclipse-temurin:21-jdk-alpine AS build
WORKDIR /app

COPY pom.xml .
COPY .mvn .mvn
COPY mvnw .
RUN chmod +x mvnw && ./mvnw dependency:go-offline -B

COPY src ./src
RUN ./mvnw package -DskipTests -B

# Stage 2: Runtime
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=build /app/target/*.jar app.jar

USER appuser
EXPOSE 8080

ENTRYPOINT ["java", "-jar", "app.jar"]
```

### Gradle Variant

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS build
WORKDIR /app

COPY gradle gradle
COPY gradlew build.gradle settings.gradle ./
RUN chmod +x gradlew && ./gradlew dependencies --no-daemon

COPY src ./src
RUN ./gradlew bootJar --no-daemon -x test

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=build /app/build/libs/*.jar app.jar
USER appuser
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

### Node.js Multi-Stage

```dockerfile
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY --from=build /app/dist ./dist
USER node
EXPOSE 3000
CMD ["node", "dist/server.js"]
```

### React Frontend Multi-Stage

```dockerfile
FROM node:20-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:1.25-alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### Multi-Stage with Test Stage

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS build
WORKDIR /app
COPY pom.xml .
COPY .mvn .mvn
COPY mvnw .
RUN chmod +x mvnw && ./mvnw dependency:go-offline -B
COPY src ./src
RUN ./mvnw package -DskipTests -B

FROM build AS test
RUN ./mvnw test

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=build /app/target/*.jar app.jar
USER 1001
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

```bash
# Build only up to the test stage
docker build --target test -t my-app:test .

# Build production image (skips test stage unless targeted)
docker build -t my-app:prod .
```

---

## 13. Docker Build Cache & Optimization

### Cache Rules

1. Each instruction creates a layer
2. Docker caches layers and reuses them if the instruction and its context haven't changed
3. If one layer's cache is invalidated, all subsequent layers are rebuilt

### Optimization Best Practices

#### 1. Order Instructions by Change Frequency

```dockerfile
# BAD — Copying source code first invalidates everything below
FROM eclipse-temurin:21-jre-alpine
COPY . /app
RUN ./mvnw dependency:go-offline
RUN ./mvnw package

# GOOD — Dependencies change less often than source code
FROM eclipse-temurin:21-jdk-alpine
WORKDIR /app
COPY pom.xml .                           # Changes rarely
RUN ./mvnw dependency:go-offline -B      # Cached when pom.xml unchanged
COPY src ./src                           # Changes frequently
RUN ./mvnw package -DskipTests -B
```

#### 2. Minimize Layer Count

```dockerfile
# BAD — Each RUN creates a layer
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y wget
RUN apt-get clean

# GOOD — Single layer
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
```

#### 3. Use Specific Base Image Tags

```dockerfile
# BAD
FROM eclipse-temurin:latest

# GOOD
FROM eclipse-temurin:21.0.2_13-jre-alpine
```

#### 4. Remove Unnecessary Files

```dockerfile
RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential && \
    make install && \
    apt-get purge -y build-essential && \    # Remove build tools
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/*              # Remove apt cache
```

#### 5. Use --no-cache-dir for pip

```dockerfile
RUN pip install --no-cache-dir -r requirements.txt
```

#### 6. Use .dockerignore

Reduces build context size and prevents cache invalidation from irrelevant files.

### Image Size Comparison

| Approach | Image Size |
|----------|-----------|
| Ubuntu + JDK + source | ~800 MB |
| Alpine + JDK + source | ~400 MB |
| Alpine + JRE (multi-stage) | ~150 MB |
| Distroless + JRE (multi-stage) | ~130 MB |
| Layered Spring Boot JAR | ~150 MB (better cache) |

### Analyzing Image Size

```bash
# View image layers and sizes
docker history my-app:1.0

# Detailed inspection
docker inspect my-app:1.0

# Use dive for interactive layer analysis
docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  wagoodman/dive my-app:1.0
```

---

## 14. Docker Security

### Security Best Practices

#### 1. Don't Run as Root

```dockerfile
# Create and use non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Or use numeric UID (more portable)
USER 1001
```

#### 2. Use Read-Only Filesystem

```bash
docker run --read-only --tmpfs /tmp my-app
```

#### 3. Drop Capabilities

```bash
docker run --cap-drop ALL --cap-add NET_BIND_SERVICE my-app
```

#### 4. Use Minimal Base Images

```dockerfile
# Prefer distroless or Alpine
FROM gcr.io/distroless/java21-debian12
# No shell, no package manager, no unnecessary utilities
```

#### 5. Scan Images for Vulnerabilities

```bash
# Docker Scout (built-in)
docker scout cves my-app:1.0
docker scout quickview my-app:1.0

# Trivy
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock \
  aquasec/trivy image my-app:1.0

# Snyk
docker scan my-app:1.0
```

#### 6. Don't Store Secrets in Images

```dockerfile
# BAD — Secret embedded in image layer
COPY credentials.json /app/
ENV DATABASE_PASSWORD=secret123

# GOOD — Use build secrets (BuildKit)
RUN --mount=type=secret,id=db_password \
    cat /run/secrets/db_password > /dev/null

# GOOD — Pass secrets at runtime
docker run -e DATABASE_PASSWORD_FILE=/run/secrets/db_pass \
  -v ./secrets/db_pass:/run/secrets/db_pass:ro my-app
```

#### 7. Use Content Trust

```bash
# Enable Docker Content Trust (image signing)
export DOCKER_CONTENT_TRUST=1
docker push my-app:1.0    # Requires signing
docker pull my-app:1.0    # Verifies signature
```

### Security Scanning in CI/CD

```yaml
# GitHub Actions example
- name: Build image
  run: docker build -t my-app:${{ github.sha }} .

- name: Scan with Trivy
  uses: aquasecurity/trivy-action@master
  with:
    image-ref: my-app:${{ github.sha }}
    severity: CRITICAL,HIGH
    exit-code: 1
```

### Docker Security Checklist

- [ ] Use minimal base images (Alpine, distroless)
- [ ] Run as non-root user
- [ ] Don't store secrets in images or environment variables
- [ ] Scan images for vulnerabilities regularly
- [ ] Use specific image tags (not `latest`)
- [ ] Enable Docker Content Trust for image signing
- [ ] Drop all capabilities and add only needed ones
- [ ] Use read-only filesystem where possible
- [ ] Limit container resources (memory, CPU)
- [ ] Keep Docker Engine and images updated
- [ ] Use multi-stage builds (no build tools in production image)
- [ ] Set up proper `.dockerignore`
- [ ] Don't expose Docker socket to containers (unless absolutely needed)
- [ ] Use network segmentation (custom networks)

---

## 15. Docker Logging & Monitoring

### Logging Drivers

| Driver | Description |
|--------|-------------|
| `json-file` | Default. Writes JSON to local files. |
| `syslog` | Sends logs to syslog daemon. |
| `journald` | Sends logs to systemd journal. |
| `fluentd` | Sends logs to Fluentd. |
| `awslogs` | Sends logs to AWS CloudWatch. |
| `gcplogs` | Sends logs to Google Cloud Logging. |
| `splunk` | Sends logs to Splunk. |
| `none` | Disables logging. |

```bash
# Set logging driver per container
docker run -d --log-driver=json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  my-app

# Set default logging driver in daemon.json
# /etc/docker/daemon.json
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "5"
  }
}
```

### Log Commands

```bash
docker logs <container>
docker logs -f <container>                     # Follow
docker logs --since 2024-01-01T00:00:00 <container>
docker logs --until 1h <container>
docker logs --tail 200 <container>
docker logs -t <container>                     # With timestamps
```

### Monitoring Commands

```bash
# Live resource usage
docker stats
docker stats --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}"

# Container processes
docker top <container>

# System resource usage
docker system df
docker system df -v                            # Verbose
```

### Prometheus + Grafana Stack (Docker Compose)

```yaml
services:
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus-data:/prometheus

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
    volumes:
      - grafana-data:/var/lib/grafana

  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    ports:
      - "8081:8080"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:rw
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro

volumes:
  prometheus-data:
  grafana-data:
```

---

## 16. Docker in CI/CD

### GitHub Actions

```yaml
name: Build and Push Docker Image

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v3

    - name: Login to Docker Hub
      uses: docker/login-action@v3
      with:
        username: ${{ secrets.DOCKERHUB_USERNAME }}
        password: ${{ secrets.DOCKERHUB_TOKEN }}

    - name: Build and test
      run: docker build --target test -t my-app:test .

    - name: Build and push
      uses: docker/build-push-action@v5
      with:
        context: .
        platforms: linux/amd64,linux/arm64
        push: ${{ github.event_name != 'pull_request' }}
        tags: |
          username/my-app:latest
          username/my-app:${{ github.sha }}
        cache-from: type=gha
        cache-to: type=gha,mode=max

    - name: Scan for vulnerabilities
      uses: aquasecurity/trivy-action@master
      with:
        image-ref: username/my-app:${{ github.sha }}
        severity: CRITICAL,HIGH
```

### GitLab CI

```yaml
stages:
  - build
  - test
  - push

variables:
  IMAGE: $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

build:
  stage: build
  image: docker:24
  services:
    - docker:24-dind
  script:
    - docker build -t $IMAGE .
    - docker save $IMAGE > image.tar
  artifacts:
    paths:
      - image.tar

test:
  stage: test
  image: docker:24
  services:
    - docker:24-dind
  script:
    - docker load < image.tar
    - docker run --rm $IMAGE npm test

push:
  stage: push
  image: docker:24
  services:
    - docker:24-dind
  script:
    - docker load < image.tar
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker push $IMAGE
  only:
    - main
```

### Jenkins Pipeline

```groovy
pipeline {
    agent any
    environment {
        REGISTRY = 'registry.example.com'
        IMAGE = "${REGISTRY}/my-app"
    }
    stages {
        stage('Build') {
            steps {
                sh "docker build -t ${IMAGE}:${BUILD_NUMBER} ."
            }
        }
        stage('Test') {
            steps {
                sh "docker run --rm ${IMAGE}:${BUILD_NUMBER} ./run-tests.sh"
            }
        }
        stage('Scan') {
            steps {
                sh "trivy image --exit-code 1 --severity HIGH,CRITICAL ${IMAGE}:${BUILD_NUMBER}"
            }
        }
        stage('Push') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'registry-creds',
                    usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                    sh "docker login -u $USER -p $PASS ${REGISTRY}"
                    sh "docker push ${IMAGE}:${BUILD_NUMBER}"
                    sh "docker tag ${IMAGE}:${BUILD_NUMBER} ${IMAGE}:latest"
                    sh "docker push ${IMAGE}:latest"
                }
            }
        }
    }
    post {
        always {
            sh "docker rmi ${IMAGE}:${BUILD_NUMBER} || true"
        }
    }
}
```

---

## 17. Docker Troubleshooting

### Common Issues & Solutions

| Symptom | Likely Cause | Solution |
|---------|-------------|----------|
| Container exits immediately | CMD/ENTRYPOINT fails or process goes to background | Check logs: `docker logs <c>`, use exec form, keep process in foreground |
| `COPY failed: file not found` | File not in build context or in `.dockerignore` | Check path relative to Dockerfile, check `.dockerignore` |
| `port already allocated` | Port in use by another process/container | `docker ps` to find conflict, use different host port |
| `no space left on device` | Docker disk full | `docker system prune -a --volumes` |
| Build is slow | Bad cache usage, large context | Optimize layer order, use `.dockerignore`, use multi-stage |
| `permission denied` | File ownership issues in container | `chown` in Dockerfile or match UID/GID |
| `network ... not found` | Network removed or misspelled | `docker network ls`, recreate network |
| Image too large | Unnecessary files, no multi-stage | Use multi-stage builds, Alpine base, `.dockerignore` |
| Container can't reach another | Different networks | Put on same network: `docker network connect` |
| DNS resolution fails | Docker DNS issues | Restart Docker daemon, use `--dns` flag |

### Debugging Commands

```bash
# View container details
docker inspect <container>
docker inspect --format='{{.State.ExitCode}}' <container>
docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <container>

# View logs
docker logs <container> 2>&1 | grep -i error

# Shell into a running container
docker exec -it <container> /bin/sh
docker exec -it <container> /bin/bash

# Shell into a stopped container (create new container from image)
docker run -it --entrypoint /bin/sh <image>

# Check container filesystem changes
docker diff <container>

# View running processes
docker top <container>

# Check resource usage
docker stats <container>

# Check events
docker events
docker events --filter 'container=<name>'
docker events --filter 'event=die'

# Networking debug
docker exec <container> cat /etc/hosts
docker exec <container> cat /etc/resolv.conf
docker network inspect bridge

# Disk usage
docker system df -v

# Build debugging
docker build --progress=plain --no-cache -t my-app .
docker history my-app:1.0
```

### Debugging Flowchart

```
Container won't start?
├── docker logs <container>
│   ├── Application error → Fix application code
│   ├── Missing dependency → Update Dockerfile
│   └── Permission denied → Fix file permissions / USER
│
├── docker inspect <container>
│   ├── ExitCode: 137 → OOM killed (increase memory)
│   ├── ExitCode: 1 → Application error
│   └── ExitCode: 126/127 → Command not found
│
Image build fails?
├── COPY/ADD fails → Check .dockerignore and build context
├── RUN fails → Run interactively from previous layer
│   └── docker run -it <previous-layer-hash> /bin/sh
├── Out of disk → docker system prune
└── Permission denied → Check file permissions in context

Network issues?
├── docker network inspect <network>
├── Containers on same network? → docker network connect
├── DNS not resolving? → Use custom bridge network (not default)
└── Port not accessible? → Check -p mapping and EXPOSE
```

---

## 18. Docker with Spring Boot — Complete Guide

### Project Structure

```
my-spring-app/
├── src/
│   └── main/
│       ├── java/
│       │   └── com/example/app/
│       │       ├── Application.java
│       │       ├── config/
│       │       ├── controller/
│       │       ├── service/
│       │       └── repository/
│       └── resources/
│           ├── application.yml
│           ├── application-docker.yml
│           └── application-prod.yml
├── Dockerfile
├── .dockerignore
├── docker-compose.yml
├── docker-compose.override.yml
├── pom.xml (or build.gradle)
└── .env
```

### Spring Boot Profile for Docker

```yaml
# application-docker.yml
server:
  port: 8080
  shutdown: graceful

spring:
  lifecycle:
    timeout-per-shutdown-phase: 30s

  datasource:
    url: jdbc:postgresql://${DB_HOST:localhost}:${DB_PORT:5432}/${DB_NAME:myapp}
    username: ${DB_USERNAME:appuser}
    password: ${DB_PASSWORD:secret}
    hikari:
      minimum-idle: 5
      maximum-pool-size: 20
      idle-timeout: 30000
      connection-timeout: 20000
      max-lifetime: 1800000

  redis:
    host: ${REDIS_HOST:localhost}
    port: ${REDIS_PORT:6379}

  jpa:
    hibernate:
      ddl-auto: validate
    properties:
      hibernate:
        dialect: org.hibernate.dialect.PostgreSQLDialect
    open-in-view: false

management:
  endpoints:
    web:
      exposure:
        include: health,info,metrics,prometheus
  endpoint:
    health:
      show-details: always
      probes:
        enabled: true
  health:
    db:
      enabled: true
    redis:
      enabled: true

logging:
  level:
    root: INFO
    com.example: INFO
  pattern:
    console: "%d{yyyy-MM-dd HH:mm:ss.SSS} [%thread] %-5level %logger{36} - %msg%n"
```

### Production-Grade Dockerfile (Maven)

```dockerfile
# ============================================
# Stage 1: Build
# ============================================
FROM eclipse-temurin:21-jdk-alpine AS build
WORKDIR /app

# Copy Maven wrapper and pom.xml first (dependency caching)
COPY pom.xml .
COPY .mvn .mvn
COPY mvnw .
RUN chmod +x mvnw

# Download dependencies (cached unless pom.xml changes)
RUN ./mvnw dependency:go-offline -B

# Copy source and build
COPY src ./src
RUN ./mvnw package -DskipTests -B && \
    mv target/*.jar target/app.jar

# ============================================
# Stage 2: Extract layers (Spring Boot layered JAR)
# ============================================
FROM eclipse-temurin:21-jdk-alpine AS extract
WORKDIR /app
COPY --from=build /app/target/app.jar app.jar
RUN java -Djarmode=layertools -jar app.jar extract

# ============================================
# Stage 3: Runtime
# ============================================
FROM eclipse-temurin:21-jre-alpine

# Security: install only necessary packages, remove caches
RUN apk add --no-cache curl tini && \
    addgroup -S spring && adduser -S spring -G spring

WORKDIR /app

# Copy layers in order of change frequency (best cache utilization)
COPY --from=extract /app/dependencies/ ./
COPY --from=extract /app/spring-boot-loader/ ./
COPY --from=extract /app/snapshot-dependencies/ ./
COPY --from=extract /app/application/ ./

# Security: non-root user
RUN chown -R spring:spring /app
USER spring

EXPOSE 8080

# Use tini as init process for proper signal handling
ENTRYPOINT ["tini", "--"]

CMD ["java", \
     "-XX:+UseContainerSupport", \
     "-XX:MaxRAMPercentage=75.0", \
     "-XX:InitialRAMPercentage=50.0", \
     "-Djava.security.egd=file:/dev/./urandom", \
     "-Dspring.profiles.active=docker", \
     "org.springframework.boot.loader.launch.JarLauncher"]
```

### Production-Grade Dockerfile (Gradle)

```dockerfile
FROM eclipse-temurin:21-jdk-alpine AS build
WORKDIR /app

COPY gradle gradle
COPY gradlew build.gradle settings.gradle ./
RUN chmod +x gradlew && ./gradlew dependencies --no-daemon

COPY src ./src
RUN ./gradlew bootJar --no-daemon -x test && \
    mv build/libs/*.jar build/libs/app.jar

FROM eclipse-temurin:21-jdk-alpine AS extract
WORKDIR /app
COPY --from=build /app/build/libs/app.jar app.jar
RUN java -Djarmode=layertools -jar app.jar extract

FROM eclipse-temurin:21-jre-alpine
RUN apk add --no-cache curl tini && \
    addgroup -S spring && adduser -S spring -G spring
WORKDIR /app

COPY --from=extract /app/dependencies/ ./
COPY --from=extract /app/spring-boot-loader/ ./
COPY --from=extract /app/snapshot-dependencies/ ./
COPY --from=extract /app/application/ ./

RUN chown -R spring:spring /app
USER spring
EXPOSE 8080

ENTRYPOINT ["tini", "--"]
CMD ["java", \
     "-XX:+UseContainerSupport", \
     "-XX:MaxRAMPercentage=75.0", \
     "-Dspring.profiles.active=docker", \
     "org.springframework.boot.loader.launch.JarLauncher"]
```

### Spring Boot Layered JAR

Spring Boot 3.x supports layered JARs that split dependencies into layers for better Docker cache performance.

```xml
<!-- pom.xml — Enable layered JAR (enabled by default in Spring Boot 3.x) -->
<build>
    <plugins>
        <plugin>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-maven-plugin</artifactId>
            <configuration>
                <layers>
                    <enabled>true</enabled>
                </layers>
            </configuration>
        </plugin>
    </plugins>
</build>
```

Layer order (least → most likely to change):

| Layer | Contents | Changes |
|-------|----------|---------|
| `dependencies` | Third-party JARs from Maven/Gradle | Rarely |
| `spring-boot-loader` | Spring Boot loader classes | Rarely |
| `snapshot-dependencies` | SNAPSHOT dependencies | Sometimes |
| `application` | Your compiled code and resources | Every build |

### Spring Boot Cloud Native Buildpacks (No Dockerfile Needed)

```bash
# Maven
./mvnw spring-boot:build-image -Dspring-boot.build-image.imageName=my-app:1.0

# Gradle
./gradlew bootBuildImage --imageName=my-app:1.0
```

```xml
<!-- pom.xml — customize buildpack image -->
<plugin>
    <groupId>org.springframework.boot</groupId>
    <artifactId>spring-boot-maven-plugin</artifactId>
    <configuration>
        <image>
            <name>registry.example.com/my-app:${project.version}</name>
            <env>
                <BP_JVM_VERSION>21</BP_JVM_VERSION>
                <BPE_JAVA_TOOL_OPTIONS>-XX:MaxRAMPercentage=75.0</BPE_JAVA_TOOL_OPTIONS>
            </env>
        </image>
    </configuration>
</plugin>
```

### .dockerignore for Spring Boot

```
.git
.gitignore
.dockerignore
Dockerfile*
docker-compose*
*.md
LICENSE

# IDE
.idea/
.vscode/
*.iml
.project
.classpath
.settings/

# Build output
target/
build/
!target/*.jar
!build/libs/*.jar
out/

# Environment
.env
.env.local
*.log

# OS
.DS_Store
Thumbs.db

# Test & docs
src/test/
docs/
```

### Complete Docker Compose for Spring Boot

```yaml
version: "3.9"

services:
  # ───────────── Spring Boot Application ─────────────
  app:
    build:
      context: .
      dockerfile: Dockerfile
    image: my-spring-app:latest
    container_name: spring-app
    ports:
      - "8080:8080"
    environment:
      SPRING_PROFILES_ACTIVE: docker
      DB_HOST: postgres
      DB_PORT: 5432
      DB_NAME: myapp
      DB_USERNAME: appuser
      DB_PASSWORD: ${DB_PASSWORD:-secret}
      REDIS_HOST: redis
      REDIS_PORT: 6379
      JAVA_TOOL_OPTIONS: >-
        -XX:+UseContainerSupport
        -XX:MaxRAMPercentage=75.0
        -XX:InitialRAMPercentage=50.0
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/actuator/health"]
      interval: 30s
      timeout: 10s
      retries: 5
      start_period: 60s
    deploy:
      resources:
        limits:
          memory: 512M
        reservations:
          memory: 256M

  # ───────────── PostgreSQL ─────────────
  postgres:
    image: postgres:16-alpine
    container_name: postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: myapp
      POSTGRES_USER: appuser
      POSTGRES_PASSWORD: ${DB_PASSWORD:-secret}
    volumes:
      - postgres-data:/var/lib/postgresql/data
      - ./docker/init-db:/docker-entrypoint-initdb.d:ro
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U appuser -d myapp"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 20s

  # ───────────── Redis ─────────────
  redis:
    image: redis:7-alpine
    container_name: redis
    ports:
      - "6379:6379"
    command: >
      redis-server
      --appendonly yes
      --maxmemory 128mb
      --maxmemory-policy allkeys-lru
    volumes:
      - redis-data:/data
    networks:
      - app-network
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3

  # ───────────── Kafka (Optional) ─────────────
  zookeeper:
    image: confluentinc/cp-zookeeper:7.6.0
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
    networks:
      - app-network
    profiles:
      - kafka

  kafka:
    image: confluentinc/cp-kafka:7.6.0
    container_name: kafka
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
    depends_on:
      - zookeeper
    networks:
      - app-network
    profiles:
      - kafka

volumes:
  postgres-data:
  redis-data:

networks:
  app-network:
    driver: bridge
```

### Development Override (docker-compose.override.yml)

```yaml
services:
  app:
    build:
      context: .
      target: build
    volumes:
      - ./src:/app/src:ro
      - ./target:/app/target
    environment:
      SPRING_PROFILES_ACTIVE: docker,dev
      SPRING_DEVTOOLS_RESTART_ENABLED: "true"
      JAVA_TOOL_OPTIONS: >-
        -XX:+UseContainerSupport
        -agentlib:jdwp=transport=dt_socket,server=y,suspend=n,address=*:5005
    ports:
      - "8080:8080"
      - "5005:5005"
```

### JVM Settings for Containers

```bash
# Essential JVM flags for containerized Java
java \
  -XX:+UseContainerSupport \        # Respect container memory limits (default in JDK 11+)
  -XX:MaxRAMPercentage=75.0 \       # Use 75% of container memory for heap
  -XX:InitialRAMPercentage=50.0 \   # Start with 50% heap
  -XX:+UseG1GC \                    # G1 garbage collector (good for containers)
  -XX:+UseStringDeduplication \     # Reduce memory for duplicate strings
  -Djava.security.egd=file:/dev/./urandom \  # Faster startup (non-blocking entropy)
  -jar app.jar
```

| JVM Flag | Purpose |
|----------|---------|
| `-XX:+UseContainerSupport` | JVM respects cgroup memory limits |
| `-XX:MaxRAMPercentage=75.0` | Max heap = 75% of container memory |
| `-XX:InitialRAMPercentage=50.0` | Initial heap size |
| `-XX:+UseG1GC` | G1 GC (recommended for most workloads) |
| `-XX:+UseZGC` | ZGC for low-latency (JDK 21+) |
| `-XX:+ExitOnOutOfMemoryError` | Exit on OOM (let container restart) |
| `-Xss512k` | Reduce thread stack size |

### Graceful Shutdown in Spring Boot

```yaml
# application.yml
server:
  shutdown: graceful

spring:
  lifecycle:
    timeout-per-shutdown-phase: 30s
```

```dockerfile
# Dockerfile — use tini for proper signal forwarding
RUN apk add --no-cache tini
ENTRYPOINT ["tini", "--"]
CMD ["java", "-jar", "app.jar"]

# Alternative: use exec form (signals forwarded to PID 1)
ENTRYPOINT ["java", "-jar", "app.jar"]
```

### Spring Boot Actuator Health Checks

```yaml
# application.yml
management:
  endpoints:
    web:
      exposure:
        include: health,info,metrics,prometheus
  endpoint:
    health:
      show-details: always
      probes:
        enabled: true        # Enables /actuator/health/liveness and /actuator/health/readiness
  health:
    livenessState:
      enabled: true
    readinessState:
      enabled: true
```

```yaml
# docker-compose.yml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8080/actuator/health"]
  interval: 30s
  timeout: 10s
  retries: 5
  start_period: 60s
```

```yaml
# Kubernetes deployment
livenessProbe:
  httpGet:
    path: /actuator/health/liveness
    port: 8080
  initialDelaySeconds: 60
  periodSeconds: 15
readinessProbe:
  httpGet:
    path: /actuator/health/readiness
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10
```

### Spring Boot Docker Best Practices Summary

| Practice | Why |
|----------|-----|
| Use multi-stage builds | Smaller images, no build tools in production |
| Use layered JAR extraction | Better Docker cache — only `application` layer rebuilds |
| Use `eclipse-temurin` JRE Alpine base | Small, well-maintained Java runtime |
| Run as non-root user | Security |
| Use `tini` as init process | Proper signal handling and zombie process reaping |
| Set `-XX:MaxRAMPercentage=75.0` | Let JVM respect container memory limits |
| Use `-XX:+UseContainerSupport` | JVM aware of container cgroup limits |
| Enable `server.shutdown=graceful` | Clean shutdown of in-flight requests |
| Use Spring profiles (`docker`, `prod`) | Environment-specific configuration |
| Externalize config via env vars | 12-factor app compliance |
| Use Actuator health endpoints | Docker/K8s health checks |
| Use `.dockerignore` | Smaller build context, faster builds |
| Scan images for vulnerabilities | Security |
| Use Cloud Native Buildpacks as alternative | No Dockerfile maintenance needed |
| Log to stdout/stderr | Docker captures logs automatically |
| Pin dependency versions | Reproducible builds |
| Use Docker Compose for local development | Easy multi-service setup |

### Full Example: Spring Boot Microservice Workflow

```bash
# 1. Build the application
./mvnw clean package -DskipTests

# 2. Build Docker image
docker build -t my-spring-app:1.0 .

# 3. Start all dependencies
docker compose up -d postgres redis

# 4. Start the application
docker compose up -d app

# 5. Check health
docker compose ps
curl http://localhost:8080/actuator/health

# 6. View logs
docker compose logs -f app

# 7. Run tests against running services
./mvnw test -Dspring.profiles.active=docker-test

# 8. Tag and push
docker tag my-spring-app:1.0 registry.example.com/my-spring-app:1.0
docker push registry.example.com/my-spring-app:1.0

# 9. Stop everything
docker compose down

# 10. Stop and clean volumes (fresh start)
docker compose down -v
```

---

## 19. Interview Questions

### Beginner Level

**Q: What is Docker and why is it used?**
A: Docker is a containerization platform that packages applications and their dependencies into lightweight, portable containers. It solves the "works on my machine" problem by ensuring consistent environments across development, testing, and production.

**Q: What is the difference between an image and a container?**
A: An image is a read-only template (like a class in OOP) containing the application and its dependencies. A container is a running instance of an image (like an object). You can create multiple containers from a single image.

**Q: What is the difference between `CMD` and `ENTRYPOINT`?**
A:
- `CMD` provides default arguments that can be completely overridden at runtime
- `ENTRYPOINT` defines the main executable and cannot be easily overridden (arguments are appended)
- Combined: `ENTRYPOINT` defines the command, `CMD` provides default arguments

**Q: What is a Docker volume and why would you use one?**
A: A Docker volume is a mechanism for persisting data generated by and used by Docker containers. Container filesystems are ephemeral — data is lost when a container is removed. Volumes persist data independently of the container lifecycle and can be shared between containers.

**Q: What happens when you run `docker run nginx`?**
A:
1. Docker checks for the `nginx` image locally
2. If not found, pulls it from Docker Hub
3. Creates a new container from the image
4. Allocates a read-write filesystem layer
5. Creates a network interface and assigns an IP
6. Starts the container and runs the default command (`CMD`)

### Intermediate Level

**Q: Explain Docker networking. How do containers communicate?**
A:
- **Bridge network (default)**: Containers on the same custom bridge network communicate by name via DNS. The default bridge doesn't support DNS.
- **Host network**: Container shares the host's network — no isolation but better performance.
- **Overlay network**: Enables communication across Docker hosts (Swarm).
- Port mapping (`-p 8080:80`) exposes container ports on the host.

**Q: How do you optimize a Dockerfile for build speed?**
A:
1. Order instructions by change frequency (dependencies before source code)
2. Use `.dockerignore` to reduce build context
3. Combine `RUN` instructions to minimize layers
4. Use multi-stage builds
5. Leverage BuildKit cache mounts
6. Use specific base image tags for consistent caching

**Q: What is a multi-stage build and when would you use it?**
A: A multi-stage build uses multiple `FROM` instructions in a single Dockerfile. Each stage can use a different base image. Only the final stage produces the output image. Use it to separate build tools (compilers, SDKs) from the runtime, resulting in smaller, more secure production images. Example: build with JDK, run with JRE.

**Q: How do you pass secrets to Docker containers securely?**
A:
- Runtime: Use Docker secrets (Compose/Swarm), mount secret files as volumes, or use environment variable files
- Build time: Use BuildKit's `--mount=type=secret` (never `ARG`/`ENV` for secrets)
- Never embed secrets in the image (they persist in layers)
- Use external secret managers (Vault, AWS Secrets Manager) in production

**Q: Explain the difference between `COPY` and `ADD`.**
A:
- `COPY`: Simply copies files from build context to image. Predictable behavior.
- `ADD`: Same as COPY but also supports URLs and auto-extracts tar archives.
- Best practice: Always use `COPY` unless you specifically need `ADD`'s extra features.

### Advanced Level

**Q: How would you reduce a Docker image from 1GB to under 200MB?**
A:
1. Use multi-stage builds (separate build and runtime)
2. Switch to Alpine or distroless base images
3. For Java: use JRE instead of JDK in runtime stage
4. Remove unnecessary packages and caches in the same `RUN` layer
5. Use Spring Boot layered JARs for better caching
6. Use `.dockerignore` to exclude unnecessary files
7. Combine `RUN` instructions to avoid intermediate layer bloat
8. Use `--no-install-recommends` with apt

**Q: How do you handle logging in containerized applications?**
A:
- Applications should log to stdout/stderr (not files)
- Docker captures stdout/stderr via logging drivers
- Configure log rotation: `--log-opt max-size=10m --log-opt max-file=5`
- For centralized logging: use Fluentd/Filebeat → Elasticsearch → Kibana
- Spring Boot: configure `logging.pattern.console` for structured JSON logs

**Q: How does Docker manage container resources on a Linux host?**
A: Docker uses Linux kernel features:
- **cgroups** (Control Groups): Limit and monitor CPU, memory, disk I/O, network
- **namespaces**: Isolate process tree (PID), network, mount points, users, hostname
- **Union filesystem** (OverlayFS): Layer-based image storage
- **seccomp**: Restrict system calls
- **AppArmor/SELinux**: Mandatory access control

**Q: Explain the Docker build cache. When is it invalidated?**
A: Docker caches each layer. Cache is invalidated when:
- The Dockerfile instruction changes
- Files referenced by `COPY`/`ADD` have changed (checksum comparison)
- Any parent layer was invalidated (all subsequent layers are rebuilt)
- `--no-cache` flag is used

Best strategy: Put rarely-changing instructions first (dependencies), frequently-changing instructions last (source code).

---

## 20. Quick Reference Cheat Sheet

```bash
# ====================== IMAGES ======================
docker pull <image>:<tag>              # Download image
docker images                          # List images
docker build -t <name>:<tag> .         # Build from Dockerfile
docker tag <src> <dest>                # Tag image
docker push <image>:<tag>              # Push to registry
docker rmi <image>                     # Remove image
docker image prune -a                  # Remove unused images
docker save -o file.tar <image>        # Export image
docker load -i file.tar                # Import image

# ===================== CONTAINERS ====================
docker run -d --name <n> -p 8080:80 <image>   # Run container
docker ps                              # List running containers
docker ps -a                           # List all containers
docker stop <container>                # Stop container
docker start <container>               # Start container
docker restart <container>             # Restart container
docker rm <container>                  # Remove container
docker rm -f $(docker ps -aq)          # Remove all containers
docker logs -f <container>             # Follow logs
docker exec -it <container> /bin/sh    # Shell into container
docker inspect <container>             # Detailed info
docker stats                           # Resource usage
docker cp <src> <container>:<dest>     # Copy files

# ===================== VOLUMES =======================
docker volume create <name>            # Create volume
docker volume ls                       # List volumes
docker volume rm <name>                # Remove volume
docker volume prune                    # Remove unused volumes
docker run -v <vol>:/data <image>      # Named volume
docker run -v $(pwd):/app <image>      # Bind mount

# ===================== NETWORKS ======================
docker network create <name>           # Create network
docker network ls                      # List networks
docker network inspect <name>          # Network details
docker network connect <net> <cont>    # Connect container
docker network rm <name>               # Remove network

# =================== COMPOSE ========================
docker compose up -d                   # Start services
docker compose down                    # Stop and remove
docker compose down -v                 # Also remove volumes
docker compose ps                      # List services
docker compose logs -f <service>       # Follow service logs
docker compose exec <svc> /bin/sh      # Shell into service
docker compose build                   # Build images
docker compose pull                    # Pull images
docker compose restart <service>       # Restart service
docker compose --profile dev up -d     # Start with profile

# ==================== CLEANUP ========================
docker system prune                    # Remove unused data
docker system prune -a --volumes       # Remove everything unused
docker system df                       # Show disk usage

# ================= QUICK GENERATORS ==================
# Generate Compose YAML from running container
docker inspect <container> | docker compose convert

# Create image from running container
docker commit <container> <image>:<tag>
```

---

*Last updated: March 2026*
