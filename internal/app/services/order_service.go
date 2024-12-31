package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lits-06/sell_technology/internal/app/models"
	"github.com/lits-06/sell_technology/pkg/db"
)

func CreateOrder(email, name, address, phone string) error {
	cart, err := queryCart(email)
    if err != nil {
        return err
    }
	if len(cart) == 0 {
		return fmt.Errorf("empty cart!!!")
	}

	userID, err := queryUserID(email)
	if err != nil {
		return err
	}

	totalPrice, err := queryCartTotalPrice(userID)
	if err != nil {
		return err
	}

	return processOrder(userID, name, address, phone, totalPrice, cart)
}

func processOrder(userID, name, address, phone string, totalPrice int, cart []map[string]interface{}) error {
	tx := db.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("can not begin transaction %v", tx.Error)
	}

	order := models.Order{
		UserID: uuid.MustParse(userID),
		Name: name,
		Address: address,
		Phone: phone,
		TotalPrice: totalPrice,
		Status: "pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error when create order")
	}

	for _, item := range cart {
		productIDStr, ok := item["product_id"].(string)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("product_id is not a string")
		}

		productID, err := uuid.Parse(productIDStr)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("invalid product_id format: %v", err)
		}

		orderDetail := models.OrderDetail{
			OrderID: order.ID,
			ProductID: productID,
			Quantity: item["quantity"].(int32),
			Price: item["price"].(int32),
		}

		if err := tx.Create(&orderDetail).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error when create order detail")
		}
	}

	if err := tx.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("error when remove user cart")
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("can not commit transaction")
	}

	return nil
}

func CancelOrder(email, orderID string) error {
	order, err := findOrderByID(orderID)
    if err != nil {
        return err
    }

	if !isCancelable(order.Status) {
        return fmt.Errorf("can not cancel order")
    }

	return updateOrderStatus(orderID, "canceled")
}

func findOrderByID(orderID string) (*models.Order, error) {
    var order models.Order
    if err := db.DB.First(&order, "id = ?", orderID).Error; err != nil {
        return nil, fmt.Errorf("order not found")
    }
    return &order, nil
}

func isCancelable(status string) bool {
    return status == "pending" || status == "completed"
}

func updateOrderStatus(orderID string, status string) error {
    return db.DB.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func GetOrders(email string) ([]models.Order, error) {
	userID, err := queryUserID(email)
	if err != nil {
		return nil, err
	}

	var orders []models.Order
    if err := db.DB.Where("user_id = ?", userID).Order("updated_at DESC").Find(&orders).Error; err != nil {
        return nil, err
    }

	return orders, nil
}

func GetItemFromOrder(orderID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := db.DB.Table("order_detail").
		Select("order_detail.id, order_detail.order_id, order_detail.product_id, products.price, order_detail.quantity, order_detail.price, products.image_url").
		Joins("JOIN products ON order_detail.product_id = products.id").
		Joins("JOIN orders ON order_detail.order_id = orders.id").
		Where("order_detail.order_id = ?", orderID).
		Order("orders.updated_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query order detail")
	}

	if results == nil {
		return []map[string]interface{}{}, nil
	}

	return results, nil
}