# E-Commerce Product API Service (HW5)

A RESTful microservice built with Go and the Gin framework implementing the Product specification from `api.yaml`, containerized with Docker, and deployed to AWS ECS Fargate using Terraform.

## Prerequisites
* [Go](https://go.dev/dl/) (1.21+)
* [Docker Desktop](https://www.docker.com/products/docker-desktop/) (running locally)
* [Terraform](https://developer.hashicorp.com/terraform/downloads)
* [AWS CLI](https://aws.amazon.com/cli/) configured with valid credentials (`~/.aws/credentials`)

---
## 📂 Directory Structure

*   **`src/`**: Contains the application source code (`product-service.go`), the `Dockerfile` for building the container, and the local observability configuration (`docker-compose.yml`, `prometheus.yml`).
*   **`terraform/`**: Contains all declarative Infrastructure-as-Code (IaC) configuration files used to provision the AWS ECS cluster, task definitions, and networking resources.
*   **`screenshots/`**: Contains all evaluation graphs and Locust results from load testing.

---

## 🚀 Deployment Instructions

To deploy this infrastructure and run the API on a new machine, follow these steps:

**1. Infrastructure Provisioning (AWS)**

Ensure you have the AWS CLI configured with appropriate credentials and Terraform installed.


cd terraform
terraform init
terraform apply -auto-approve

**2. Local Observability Stack**

   To run Prometheus and Grafana locally for monitoring:

cd src
docker-compose up -d
Grafana will be available at http://localhost:3000.

 GET Product (200 OK)

curl -i http://<your-server-url>:8080/products/1 

POST Update Product Details (200 OK)

curl -i -X POST http://<your-server-url>:8080/products/1/details \
-H "Content-Type: application/json" \
-d '{"sku": "SKU-123", "manufacturer": "Acme", "weight": 1250}'
3. Error Handling (404 Not Found)

Bash
curl -i http://<your-server-url>:8080/invalid-endpoint
Request Logs

2026/10/18 22:00:11 {"duration":0,"endpoint":"/products/:productId/details","method":"POST","status":204}
[GIN] 2026/10/18 - 22:00:11 | 204 | 57.792µs |             ::1 | POST     "/products/1/details"

2026/10/18 22:00:11 {"duration":0,"endpoint":"/products/:productId","method":"GET","status":200}
[GIN] 2026/10/18 - 22:00:11 | 200 | 22.791µs |             ::1 | GET      "/products/1"

2026/10/18 22:00:11 {"duration":0,"endpoint":"/products/:productId","method":"GET","status":200}
[GIN] 2026/10/18 - 22:00:11 | 200 | 25.834µs |             ::1 | GET      "/products/1"
