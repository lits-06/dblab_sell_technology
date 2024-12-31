package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lits-06/sell_technology/internal/app/models"
	"github.com/lits-06/sell_technology/internal/app/services"
	"github.com/lits-06/sell_technology/pkg/utils"
)

// chưa có log
func Register(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
		return
	}

	err := services.Register(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	token, err := services.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	foundUser, _ := services.GetUserByEmail(user.Email)

	c.JSON(http.StatusCreated, gin.H{
		"Email": foundUser.Email,
		"Name": foundUser.Name,
		"Token": token,
	})
}

// chưa có log
func Login(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
		return
	}

	token, err := services.Login(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	foundUser, _ := services.GetUserByEmail(user.Email)

	c.JSON(http.StatusOK, gin.H{
		"Email": foundUser.Email,
		"Name": foundUser.Name,
		"Token": token,
	})
}

func UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid input"})
		return
	}

	err := services.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
        return
	}

	c.JSON(http.StatusOK, "User info updated successfully")
}

// chưa có log
func GetUserInfo(c *gin.Context) {
	e, emailExists := c.Get("email")

	if !emailExists {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		return
	}

	email, _ := e.(string)

	user, err := services.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email":  user.Email,
		"name":   user.Name,
		"phone":  user.Phone,
		"address": user.Address,
	})
}

func GetUserCart(c *gin.Context) {
	e, emailExists := c.Get("email")

	if !emailExists {
		utils.Logger.Warn(
			"Unauthorized access attempt",
			slog.String("endpoint", c.Request.URL.Path),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "Unauthorized"})
		return
	}

	email, _ := e.(string)

	products, err := services.GetCart(email)
	if err != nil {
		utils.Logger.Error(
			"Failed to retrieve user cart",
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to retrieve cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": len(products),
		"products": products,
	})
}

func GetUserCartTotalPrice(c *gin.Context) {
	e, _ := c.Get("email")
	email, _ := e.(string)

	totalPrice, err := services.GetCartTotalPrice(email)
	if err != nil {
		utils.Logger.Error(
			"Failed to calculate user cart total price",
			slog.String("email", email),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Failed to calculate user cart total price"})
		return
	}

	c.JSON(http.StatusOK, totalPrice)
}

func AddToCart(c *gin.Context) {
	e, _ := c.Get("email")
	email, _ := e.(string)

	var request struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	productID := request.ProductID
	quantity := request.Quantity

	err := services.AddToCart(email, productID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Product added to cart successfully")
}

func UpdateCartItem(c *gin.Context) {
	e, _ := c.Get("email")
	email, _ := e.(string)

	var request struct {
		ProductID string `json:"product_id" binding:"required"`
		Quantity  int    `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	productID := request.ProductID
	quantity := request.Quantity

	err := services.UpdateCartItem(email, productID, quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Cart updated successfully")
}

func RemoveFromCart(c *gin.Context) {
	e, _ := c.Get("email")
	email, _ := e.(string)
	productID := c.Param("product_id")

	err := services.RemoveFromCart(email, productID)
	if err != nil {
		if err.Error() == "product not found in cart" {
			c.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, "Product removed from cart successfully")
}

func CancelOrder(c *gin.Context) {
	e, _ := c.Get("email")
	email, _ := e.(string)

	var request struct {
		OrderID string `json:"order_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid request body"})
		return
	}

	orderID := request.OrderID
	err := services.CancelOrder(email, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Can not cancel order"})
        return
	}

	c.JSON(http.StatusOK, "Order has been successfully cancelled")
}