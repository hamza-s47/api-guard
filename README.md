# API Guard

A **production-ready API Gateway** built in **Go** that provides reverse proxying, middleware-based rate limiting, logging, health checks, and graceful shutdown. Designed with clean architecture and extensibility in mind.

---

## Features

* 🔁 **Reverse Proxy** – Routes API requests to backend services
* 🚦 **Rate Limiting** – Per-client (IP-based) token bucket limiter
* 🧩 **Middleware Architecture** – Clean, composable HTTP middleware
* 🌐 **Real Client IP Detection** – Supports `X-Forwarded-For` and `X-Real-IP`
* 📊 **Rate Limit Headers** – Exposes limit and remaining quota
* ❤️ **Health Check Endpoint** – For monitoring and load balancers
* 🛑 **Graceful Shutdown** – Safe termination on SIGINT / SIGTERM
* 🧠 **Extensible Design** – Redis-based distributed rate limiting ready (v2)

---

## Architecture Overview

```
Client
  ↓
API Gateway (:8080)
  ├─ Rate Limit (real IP)
  ├─ Headers
  ├─ Logging
  └─ Reverse Proxy
        ↓
Backend Services (:9000+)
```

### Request Flow (Step-by-Step)

1. Client sends request to **API Gateway (:8080)**
2. Request passes through **Logging Middleware**
3. Request passes through **Rate Limiting Middleware**
4. Gateway router (`ServeMux`) decides:

   * `/health` → local handler
   * `/api/*` → reverse proxy
5. Reverse proxy forwards request to backend service (:9000)
6. Response flows back through middleware to client

---

<!--## 📁 Project Structure

```
api-guard/
├── cmd/
│   └── main.go                # Application entry point
├── handler/
│   └── health.go              # Health check handler
├── middleware/
│   ├── logging.go             # Request logging middleware
│   ├── rate_limit.go          # In-memory rate limiting middleware
│   └── redis_rate_limit.go    # (Optional) Redis-based limiter
├── internal/
│   ├── limiter/
│   │   └── token_bucket.go    # Token bucket implementation
│   ├── proxy/
│   │   └── reverse_proxy.go   # Reverse proxy logic
│   └── store/
│       ├── memory_store.go    # In-memory store for rate limits
│       └── redis_store.go     # Redis store (optional / v2)
├── go.mod
├── go.sum
└── README.md
```

--- -->

## Getting Started

### Prerequisites

* Go 1.21+

---

### Run Backend Service (Dummy)

```bash
go run main.go   # Runs on :9000
```

Or any HTTP service on port `9000`.

---

### Run API Gateway

```bash
go run cmd/main.go
```

Gateway will start on:

```
http://localhost:8080
```

---

## 🧪 Testing

### Health Check

```bash
curl http://localhost:8080/health
```

### Rate Limiting Test

```bash
for i in {1..10}; do curl http://localhost:8080/api/test; done
```

After the limit is reached:

```
HTTP/1.1 429 Too Many Requests
```

Response headers:

```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 0
```

---

## 🧠 Rate Limiting Strategy

* **Algorithm:** Token Bucket (in-memory)
* **Key:** Client IP address
* **Scope:** Per-IP, per-gateway instance

### Note

The current implementation is intentionally kept in-memory for simplicity and clarity, a Redis-backed distributed rate limiter can be prepared.

---

## Graceful Shutdown

The gateway listens for:

* `SIGINT` (Ctrl+C)
* `SIGTERM` (Docker / Kubernetes)

Active requests are allowed to complete before shutdown.

---

## Tech Stack

* **Language:** Go
* **HTTP:** net/http
* **Architecture:** Middleware-based

---

## Possible Enhancements

* Per-route rate limiting
* Config-driven rate limits (YAML / Dynamic Config)
* Redis-backed distributed rate limiting
* Admin / Control Plane API
* Authentication & Authorization middleware (JWT/API key validation at the gateway level)

---