FROM --platform=linux/amd64 golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY product-service.go .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o product-service product-service.go

FROM --platform=linux/amd64 alpine:latest
WORKDIR /app
COPY --from=builder /app/product-service .

EXPOSE 8080
CMD ["./product-service"]