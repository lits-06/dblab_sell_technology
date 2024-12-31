package services

import (
	"github.com/lits-06/sell_technology/internal/app/models"
	"github.com/lits-06/sell_technology/pkg/db"
	"gorm.io/gorm"
)

func GetCategories() ([]models.Category, error) {
	return queryCategory("")
}

func GetCategoryByName(name string) (models.Category, error) {
	categories, err := queryCategory(name)
	if err != nil {
        return models.Category{}, err
    }

	if len(categories) == 0 {
        return models.Category{}, gorm.ErrRecordNotFound
    }

	return categories[0], nil
}

func queryCategory(name string) ([]models.Category, error) {
	var categories []models.Category

	if name == "" {
		if err := db.DB.Find(&categories).Error; err != nil {
			return []models.Category{}, err
		}
		return categories, nil
	} else {
		if err := db.DB.Where("name = ?", name).Find(&categories).Error; err != nil {
			return []models.Category{}, err
		}
		return categories, nil
	}
}