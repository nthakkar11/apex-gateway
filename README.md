# Apex Gateway: High-Throughput Financial Infrastructure

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![Redis](https://img.shields.io/badge/Redis-FF4438?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

**Apex Gateway** is a distributed API middleware engineered for high-frequency environments where **data integrity** and **system uptime** are critical. It solves the "Double-Spend" and "Retry-Storm" problems by integrating atomic idempotency shields and distributed rate limiting into a single, high-performance layer.

[Live Demo]([YOUR_DEPLOYED_URL]) • [Production Notes](./PRODUCTION.md) • [Report Bug](https://github.com/[YOUR_GITHUB_USERNAME]/apex-gateway/issues)

---

## ⚡ The Problem: Why This Exists
In distributed systems (like Betting or Payment platforms), network partitions are inevitable. 
1. **The Double-Spend:** A user clicks 'Pay' twice during a lag; without **Idempotency**, they are charged twice.
2. **The Retry Storm:** Thousands of clients automatically retry failed requests, creating a self-inflicted DDoS. Without **Rate Limiting**, the database crashes.

**Apex Gateway solves both.**

---

## 🏗️ System Architecture
![alt text](apex-gateway.png)


### How it Works (The "Crazy" Scale)
Unlike standard middleware that performs multiple database round-trips, Apex Gateway uses **Redis Lua Scripting**. This ensures that the Rate Limit check and the Idempotency lookup happen **atomically** in a single trip, allowing the system to handle **10,000+ Requests Per Minute** with sub-10ms latency.

---

## 🚀 Core Features

- **Atomic Gatekeeping:** Uses Lua scripts to eliminate race conditions in high-concurrency environments.
- **Idempotency Shield:** Implements an "Exactly-Once" processing guarantee for financial transactions.
- **Distributed State:** Designed to scale horizontally across multiple regions while maintaining a unified state in Redis.
- **Load-Tested Performance:** Sustained **p99 < 15ms** under heavy load (Tested via Artillery).

---

## 🛠️ Tech Stack

- **Backend:** Golang (Standard Library + Gin for Routing)
- **Data Store:** Redis (Atomic Lua Scripting, TTL Caching)
- **Testing:** Artillery (Load Testing), Go Subtests (Race Condition Testing)
- **Ops:** Render/Railway (PaaS), GitHub Actions

---

## 🚦 Getting Started

### Installation
1. **Clone & Install**
   ```bash
   git clone [https://github.com/](https://github.com/)[YOUR_GITHUB_USERNAME]/apex-gateway.git
   cd apex-gateway
   go mod tidy