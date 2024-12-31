package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lits-06/sell_technology/internal/app/services"
	"gorm.io/gorm"
)

// chưa có log
func GetProducts(c *gin.Context) {
	products, err := services.QueryProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// chưa có log
func GetProductsByID(c *gin.Context) {
	id := c.Param("id")

	product, err := services.QueryProductByID(id)
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"Error": "Product not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": fmt.Sprintf("Failed to retrieve product %v", err.Error())})
		return
	}

	c.JSON(http.StatusOK, product)
}