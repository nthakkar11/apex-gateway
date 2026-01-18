# Apex Gateway
**High-Throughput Financial Middleware for Distributed System Integrity**

[![Go Report Card](https://goreportcard.com/badge/github.com/nthakkar11/apex-gateway)](https://goreportcard.com/report/github.com/nthakkar11/apex-gateway)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Apex Gateway is a performance-critical API middleware built in Go, designed to solve the dual challenges of **transactional idempotency** and **distributed rate limiting** in high-frequency environments. By leveraging atomic Lua scripting within Redis, the gateway eliminates race conditions and ensures "exactly-once" processing for sensitive financial operations.

[Live Instance](https://apex-gateway.onrender.com) | [Technical Architecture](./PRODUCTION.md)

---

## The Engineering Problem
In distributed systems, network partitions and client-side retries often lead to two critical failures:
1. **The Double-Spend:** Multiple concurrent requests resulting in duplicate state changes.
2. **Resource Exhaustion:** Uncontrolled retry storms that overwhelm downstream services.

Apex Gateway centralizes the protection layer using a **Stateless Backend + Centralized State** architecture.

---

## Core Capabilities

* **Atomic Transaction Guard:** Executes rate-limiting and idempotency checks in a single Redis RTT using Lua scripting to prevent "Time-of-Check to Time-of-Use" (TOCTOU) bugs.
* **Distributed Rate Limiting:** Implements a sliding window/fixed-window hybrid to protect service availability.
* **Idempotency Shield:** Caches successful transaction outcomes for a 24-hour TTL, providing consistent responses to redundant requests.
* **Concurrency-First Design:** Engineered with Go’s `net/http` standard library for minimal overhead and maximum throughput.

---

## 🏗️ System Architecture

![Apex Gateway Sequence Diagram](./apex-gateway.png)

---

## Performance Benchmarks
Verified via local and remote stress testing:
- **Throughput:** 10,000+ RPM.
- **Latency (p99):** < 20ms (Remote Redis Cluster).
- **Reliability:** 100% success rate on 110-request concurrent bursts.

---

## 🚥 Interactive Testing Suite

The [Live Workbench](https://apex-gateway.onrender.com/) allows you to audit the system's logic in real-time:

![Apex Gateway Workbench](./workbench-screenshot.png)

### 1. Idempotency Validation (The "Exactly-Once" Test)
* **Action**: Click **IDEMPOTENCY_VALIDATION**.
* **Behavior**: Fires 5 requests with the **same key**.
* **Result**: 1 success (201) and 4 cache hits (200), proving the system remembers the first transaction.

### 2. Concurrent Stress Test (The "Hammer" Test)
* **Action**: Click **CONCURRENT_STRESS**.
* **Behavior**: Fires **110 requests** simultaneously.
* **Result**: Exactly 100 Allowed (201) and 10 Blocked (429), proving atomic threshold enforcement.

### 3. Manual Multi-User Test
* **Action**: Input unique **USER_ID**s and custom **IDEMPOTENCY_KEY**s.
* **Result**: Demonstrates per-user isolation and manual override of the state cache.

---

## Technical Stack
- **Runtime:** Go 1.21+
- **State Store:** Redis (with Lua 5.1 scripting)
- **Observability:** Structured logging for audit trails
- **Deployment:** Render (PaaS) with Upstash (Serverless Redis)

---

## Getting Started

### Local Development
1. **Clone the repository:**
   ```bash
   git clone [https://github.com/nthakkar11/apex-gateway.git](https://github.com/nthakkar11/apex-gateway.git)
   cd apex-gateway

Configure Environment:

Bash

export REDIS_URL="redis://default:password@host:port"
Run Suite:

Bash

go test -v gateway_test.go
go run main.go


Verification
Bash

# Initial Request (201 Created)
curl -i -H "X-Idempotency-Key: txn_01" "http://localhost:8080/v1/transaction?user_id=u1"

# Replay Attack/Retry (200 OK via Cache)
curl -i -H "X-Idempotency-Key: txn_01" "http://localhost:8080/v1/transaction?user_id=u1"

