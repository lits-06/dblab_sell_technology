package services

import (
	"github.com/lits-06/sell_technology/internal/app/models"
	"github.com/lits-06/sell_technology/pkg/db"
)

func QueryProducts() ([]models.Product, error) {
	var products []models.Product

	if err := db.DB.Find(&products).Error; err != nil {
		return []models.Product{}, err
	}
	return products, nil
}

func QueryProductByID(id string) (models.Product, error) {
	var product models.Product
	err := db.DB.Where("id = ?", id).First(&product).Error
	if err != nil {
		return models.Product{}, err
	}
	return product, nil
}