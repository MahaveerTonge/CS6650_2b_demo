package main

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()

	// Support both direct paths and /v1 prefixed paths from api.yaml
	r.GET("/products/:productId", getProduct)
	r.POST("/products/:productId/details", addProductDetails)
	r.GET("/v1/products/:productId", getProduct)
	r.POST("/v1/products/:productId/details", addProductDetails)

	// Matches default container_port (8080) in terraform/variables.tf
	r.Run(":8080")
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
