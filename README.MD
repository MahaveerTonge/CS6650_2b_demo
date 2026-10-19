# Product API: Infrastructure, Load Testing & Observability

This repository contains a highly scalable, containerized Go microservice deployed on AWS ECS Fargate using Terraform. It includes comprehensive load testing results (Locust) and a local observability stack (Prometheus & Grafana) tracking RED (Rate, Errors, Duration) metrics.

##  Directory Structure

*   **`src/`**: Contains the application source code (`product-service.go`), the `Dockerfile` for building the container, and the local observability configuration (`docker-compose.yml`, `prometheus.yml`).
*   **`terraform/`**: Contains all declarative Infrastructure-as-Code (IaC) configuration files used to provision the AWS ECS cluster, task definitions, and networking resources.
*   **`screenshots/`**: Contains all evaluation graphs and Locust results from load testing.

---

## Deployment Instructions

To deploy this infrastructure and run the API on a new machine, follow these steps:

**1. Infrastructure Provisioning (AWS)**
Ensure you have the AWS CLI configured with appropriate credentials and Terraform installed.
```bash
cd terraform
terraform init
terraform apply -auto-approve
```

**2. Local Observability Stack**
To run Prometheus and Grafana locally for monitoring:
```bash
cd src
docker-compose up -d
```
Grafana will be available at `http://localhost:3000`.

---

## API Endpoints & Usage Examples

Here are examples of how to interact with the server and the expected response codes.

**1. GET Product (200 OK)**
```bash
curl -i http://<your-server-url>:8080/products/1
```

**2. POST Update Product Details (200 OK)**
```bash
curl -i -X POST http://<your-server-url>:8080/products/1/details \
-H "Content-Type: application/json" \
-d '{"sku": "SKU-123", "manufacturer": "Acme", "weight": 1250}'
```

**3. Error Handling (404 Not Found)**
```bash
curl -i http://<your-server-url>:8080/invalid-endpoint
```
![Baseline 50 Users](/Users/mahaveer/Desktop/Screenshot 2026-10-18 at 5.07.19 PM.png)
---

## Part 4: Load Testing & Performance Report

The system was aggressively load-tested using Locust to evaluate steady-state performance and failure thresholds under high concurrency.

### Baseline Test (50 Users)
At a low concurrency of 50 users, the system performed flawlessly.
*   **Failure Rate:** 0%
*   **Latency:** Median response times hovered around ~20ms.
*   **Conclusion:** Both standard `HttpUser` and `FastHttpUser` frameworks handled the baseline load identically with negligible overhead.

![Baseline 50 Users](/Users/mahaveer/Desktop/normal50.png)
![Baseline 50 Users](/Users/mahaveer/Desktop/Fast50.png)

### Stress Test & Capacity Limits (5,000 Users)
To test architectural resilience, concurrency was scaled up to 5,000 active users.
*   **Resilience:** The Go backend and AWS ECS infrastructure successfully processed hundreds of thousands of requests with a **0% failure rate**.
*   **Throughput:** Standard `HttpUser` achieved an impressive **~2,032 Requests Per Second (RPS)** with a 99th percentile latency of 1.8 seconds.

![Stress Test 5000 Users](/Users/mahaveer/Desktop/normal5000.png)

### Comparative Analysis: `HttpUser` vs. `FastHttpUser`
A comparative stress test was run using `FastHttpUser` (which uses a C-based HTTP parser) to see if throughput could be improved further.
*   **Finding:** Counterintuitively, `FastHttpUser` performed *worse* at 5,000 users, dropping to ~1,369 RPS and experiencing massive tail latency spikes (99th percentile of 7.5 seconds).
*   **Architectural Explanation:** Because `FastHttpUser` fires non-blocking requests so aggressively over the internet, it saturated the client machine's local TCP socket buffers and port allocations. This caused local network congestion on the testing machine before the AWS server could even be overwhelmed. The standard `HttpUser` paced requests more cleanly, resulting in higher sustained throughput.

![FastHttpUser 5000 Users](/Users/mahaveer/Desktop/fast 5000.png)

---

##  Part 5: Observability & RED Signals

To ensure system health can be monitored in real-time and debugged post-incident, the API is instrumented to capture RED signals (Rate, Errors, Duration) using Prometheus metrics and structured JSON logging.

### Grafana Dashboard (Macro-Level Monitoring)
Below is the Grafana dashboard tracking the RED signals, clearly demonstrating the transition from an "idle" state to a heavily "loaded" state during the Locust stress test.

![Grafana Dashboard](/Users/mahaveer/Desktop/Screenshot 2026-10-18 at 10.04.04 PM.png)

### Structured JSON Logs (Micro-Level Debugging)
Every request generates a structured log entry for deep debugging. Here is an excerpt from the test run:

```json
2026/10/18 22:00:11 {"duration":0,"endpoint":"/products/:productId","method":"GET","status":200}[cite: 7]
2026/10/18 22:00:11 {"duration":0,"endpoint":"/products/:productId/details","method":"POST","status":204}[cite: 7]
2026/10/18 22:00:11 {"duration":0,"endpoint":"/products/:productId","method":"GET","status":200}[cite: 7]
2026/10/18 22:00:11 {"duration":0,"endpoint":"/products/:productId/details","method":"POST","status":204}[cite: 7]
```

### Dashboard vs. Logs Strategy
I use the **Grafana dashboard** to monitor macro-level system health and identify real-time anomalies or trends (such as a sudden, unexpected latency spike). I use the **structured JSON request logs** for micro-level debugging to investigate exactly *what* went wrong with specific failed requests after an anomaly is detected.