# 🚀 Go-Based High-Performance Data Center Core & Security Monitor

This repository contains a high-performance backend core designed in **Go (Golang)** to simulate, monitor, and protect data center infrastructure metrics (CPU load, temperature, energy consumption) in real-time. Inspired by industry-leading infrastructure management platforms, this project solves two critical data center challenges: cost optimization and edge cybersecurity.

---

## 🛠️ Key Features

* **Real-Time Core Metrics:** Utilizes native operating system calls to fetch and track precise CPU utilization and dynamically calculated thermal metrics.
* **Smart Energy Logging:** Automatically benchmarks core operational states and writes highly structural energy telemetry to secure logs (`datacenter_energy_log.txt`).
* **Built-in Cybersecurity Countermeasures:** Features an active DDoS/Brute-Force mitigation framework that dynamically throttles simulated malicious requests to secure data center endpoints.
* **Enterprise Architecture:** Structured completely with Go modules, leveraging highly concurrent goroutines for non-blocking I/O operations.

---

## 📊 Technical Architecture

* **Language:** Go (Golang) 
* **Core Libraries:** `gopsutil` (Native hardware monitoring)
* **Design Pattern:** Structural concurrency and event-driven monitoring loop.

---

## 🚀 How to Run It Globally

1. Ensure you have **Go** installed on your machine.
2. Clone this production-ready repository:
   ```bash
   git clone [https://github.com/Aliakber123123/Aliakber.git](https://github.com/Aliakber123123/Aliakber.git)
