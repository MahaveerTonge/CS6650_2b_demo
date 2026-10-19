package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{Name: "http_requests_total", Help: "Total number of HTTP requests"},
		[]string{"path", "method", "status"},
	)
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{Name: "http_request_duration_seconds", Help: "Histogram of response latencies"},
		[]string{"path", "method"},
	)
)

type Product struct {
	ProductID    int    `json:"product_id" binding:"required,min=1"`
	SKU          string `json:"sku" binding:"required,min=1,max=100"`
	Manufacturer string `json:"manufacturer" binding:"required,min=1,max=200"`
	CategoryID   int    `json:"category_id" binding:"required,min=1"`
	Weight       int    `json:"weight" binding:"min=0"`
	SomeOtherID  int    `json:"some_other_id" binding:"required,min=1"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

var (
	products = make(map[int]Product)
	mu       sync.RWMutex
)

func ObservabilityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}

		// Process request
		c.Next()

		duration := time.Since(start)
		status := strconv.Itoa(c.Writer.Status())

		// Update Prometheus Metrics
		httpRequestsTotal.WithLabelValues(path, c.Request.Method, status).Inc()
		httpRequestDuration.WithLabelValues(path, c.Request.Method).Observe(duration.Seconds())

		// Emit Structured JSON Log
		logData := map[string]interface{}{
			"endpoint": path,
			"method":   c.Request.Method,
			"status":   c.Writer.Status(),
			"duration": duration.Milliseconds(),
		}
		jsonLog, _ := json.Marshal(logData)
		log.Println(string(jsonLog))
	}
}

func init() {
	products[1] = Product{
		ProductID:    1,
		SKU:          "ABC-123-XYZ",
		Manufacturer: "Acme Corporation",
		CategoryID:   456,
		Weight:       1250,
		SomeOtherID:  789,
	}
}

func main() {

	router := gin.Default()
	router.Use(ObservabilityMiddleware())
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/products/:productId", getProduct)
	router.POST("/products/:productId/details", addProductDetails)
	router.GET("/v1/products/:productId", getProduct)
	router.POST("/v1/products/:productId/details", addProductDetails)
	router.Run(":8080")
}

func getProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("productId"))
	if err != nil || id < 1 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "Product ID must be a positive integer",
		})
		return
	}

	if id == 500 {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
			Details: "Simulated server failure",
		})
		return
	}

	mu.RLock()
	product, exists := products[id]
	mu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func addProductDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("productId"))
	if err != nil || id < 1 {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Product not found",
			Details: "Product ID must be a positive integer",
		})
		return
	}

	if id == 500 {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
			Details: "Simulated server failure",
		})
		return
	}

	var input Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: err.Error(),
		})
		return
	}

	if input.ProductID != id {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_INPUT",
			Message: "Invalid input data",
			Details: "Path productId must match body product_id",
		})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := products[id]; !exists {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "NOT_FOUND",
			Message: "Product not found",
		})
		return
	}

	products[id] = input
	c.Status(http.StatusNoContent)
}
