# E-Commerce Product API Service (HW5)

A RESTful microservice built with Go and the Gin framework implementing the Product specification from `api.yaml`, containerized with Docker, and deployed to AWS ECS Fargate using Terraform.

## Prerequisites
* [Go](https://go.dev/dl/) (1.21+)
* [Docker Desktop](https://www.docker.com/products/docker-desktop/) (running locally)
* [Terraform](https://developer.hashicorp.com/terraform/downloads)
* [AWS CLI](https://aws.amazon.com/cli/) configured with valid credentials (`~/.aws/credentials`)

---

## Deploying Infrastructure on Any Machine (Terraform)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/MahaveerTonge/CS6650_2b_demo
   cd CS6650_2b_demo/terraform
Initialize Terraform providers and modules:

Bash
terraform init
Build the Docker image, push to ECR, and provision ECS Fargate:

Bash
terraform apply -auto-approve
Retrieve the live server address:
After deployment completes, note the public IP output displayed in your terminal (or run terraform output).

Tear down infrastructure when finished:

Bash
terraform destroy -auto-approve
API Documentation & Response Code Examples
Replace BASE_URL below with http://localhost:8080 (local testing) or http://<ECS_PUBLIC_IP>:8080 (AWS deployment).

1. GET /products/{productId}
200 OK – Product found successfully:

Bash
curl -i http://BASE_URL/products/1
404 Not Found – Product does not exist:

Bash
curl -i http://BASE_URL/products/999
500 Internal Server Error – Simulated server failure:

Bash
curl -i http://BASE_URL/products/500
2. POST /products/{productId}/details
204 No Content – Product details updated successfully:

Bash
curl -i -X POST http://BASE_URL/products/1/details \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 1,
    "sku": "NEW-SKU-999",
    "manufacturer": "Updated Corp",
    "category_id": 10,
    "weight": 500,
    "some_other_id": 20
  }'
400 Bad Request – Invalid input data (missing required fields or ID mismatch):

Bash
curl -i -X POST http://BASE_URL/products/1/details \
  -H "Content-Type: application/json" \
  -d '{"product_id": 1, "sku": ""}'
404 Not Found – Target product ID does not exist:

Bash
curl -i -X POST http://BASE_URL/products/999/details \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 999,
    "sku": "GHOST-SKU",
    "manufacturer": "Acme",
    "category_id": 1,
    "weight": 100,
    "some_other_id": 1
  }'
500 Internal Server Error – Simulated server failure:

Bash
curl -i -X POST http://BASE_URL/products/500/details \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": 500,
    "sku": "ERR-SKU",
    "manufacturer": "Acme",
    "category_id": 1,
    "weight": 100,
    "some_other_id": 1
  }'
