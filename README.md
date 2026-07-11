# ⚙️ Enterprise Cloud Cluster Optimizer & FinOps Engine

An advanced, high-performance cluster management and cost-optimization infrastructure simulator written in Go. This system models a large-scale enterprise environment (simulating a 100,000-node core matrix) to dynamically balance incoming high-burst transaction traffic while aggressively eliminating server resource underutilization to reduce cloud infrastructure costs in real-time.

---

## 💎 Core Business Value & Strategy (FinOps)

In massive cloud environments, running underutilized or idle servers costs organizations millions of dollars annually. This engine addresses this issue through automated, programmatic financial operations:
* **Resource Waste Elimination (Cost Cutter):** Automatically scales down and puts redundant servers into deep sleep modes when the average cluster load drops below 25%, dropping idle power usage from 150W to 12W.
* **Proactive Auto-Scaling (SLA Protection):** Continuously monitors the high-speed task queue backlog. If a traffic burst is detected or the cluster average load passes 65%, sleeping nodes are instantly awoken to maintain Service Level Agreements (SLA) and prevent dropped requests.
* **Real-Time Financial Analytics:** Simulates real-world electricity data center tariffs ($0.00015 per kWh multiplied across a global scale) to output continuous, exact dollar savings directly to the system console and live telemetry dashboard.

---

## 🏗️ Architecture Topography

The system shifts away from basic synchronous API routing by implementing robust, thread-safe concurrency primitives:
1. **Asynchronous Task Queue:** A bounded channel buffer capable of absorbing sudden high-velocity traffic bursts (up to 5,000 concurrent packets) without choking the main gateway.
2. **Worker Pool Pattern:** Multiplexed background goroutines acting as independent core processors that safely dequeue tasks and route them to the most optimal, lowest-loaded active server instance.
3. **Thread-Safe Telemetry:** Employs Read-Write Mutexes (`sync.RWMutex`) and atomic operations (`sync/atomic`) to guarantee precise metrics evaluation under extreme load conditions.

---

## 🚀 Live Endpoint Telemetry

The core manager exposes production-ready REST API endpoints to stream the state of the cluster matrix:

### 1. Cluster Status Dashboard
* **Endpoint:** `GET /api/cluster/dashboard`
* **Returns:** A live JSON map tracking operational health, active/inactive node counts, precise CPU loads, internal temperatures (°C), power capping metrics, and cumulative enterprise money saved in USD.

### 2. High-Velocity Traffic Burst Simulator
* **Endpoint:** `GET /api/cluster/simulate-burst`
* **Trigger:** Instantly injects a sudden surge of 1,000 heavy transactional packets into the cluster queue to stress-test the pro-active auto-scaling and worker load allocation response.

---

## ⚡ How to Run & Test

1. Execute the core engine:
   ```bash
   go run main.go