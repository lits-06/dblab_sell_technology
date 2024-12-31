package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lits-06/sell_technology/internal/app/services"
	"github.com/lits-06/sell_technology/pkg/utils"
)

func GetCategories(c *gin.Context) {
	categories, err := services.GetCategories()
	if err != nil {
		utils.Logger.Error(
			"Failed to fetch categories", 
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategoryByName(c *gin.Context) {
	name := c.Param("name")

	category, err := services.GetCategoryByName(name)
	if err != nil {
		utils.Logger.Warn(
			"Category not found",
			slog.String("name", name),
			slog.String("error", err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}