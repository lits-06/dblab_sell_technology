package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lits-06/sell_technology/internal/app/services"
)

func CreateOrder(c *gin.Context) {
	e, _ := c.Get("email")
	email, _ := e.(string)

	var request struct {
		Name    string `json:"name" binding:"required"`
        Address string `json:"address" binding:"required"`
        Phone   string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	err := services.CreateOrder(email, request.Name, request.Address, request.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Create order success")
}

func GetOrder(c *gin.Context) {
	e, _ := c.Get("email")
	email, _ := e.(string)

	orders, err := services.GetOrders(email)
	if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch orders"})
        return
    }

	var orderDetails = []map[string]interface{}{}
	if (len(orders) != 0) {
		for _, order := range orders {
			items, err := services.GetItemFromOrder(order.ID.String())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to fetch order items"})
				return
			}

			orderDetails = append(orderDetails, items...)
		}
	}
    
	c.JSON(http.StatusOK, gin.H{
        "orders":       orders,
        "order_detail": orderDetails,
    })
}