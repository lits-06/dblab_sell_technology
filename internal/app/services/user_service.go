package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/lits-06/sell_technology/internal/app/models"
	"github.com/lits-06/sell_technology/pkg/db"
	"github.com/lits-06/sell_technology/pkg/utils"
	"gorm.io/gorm"
)

func Register(registerUser models.User) error {
	_, err := queryUser(registerUser.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("error when query user register %v", err)
	}
	if err == nil {
		return fmt.Errorf("email already exists")
	}

	hashPassword, err := utils.HassPasword(registerUser.Password)
	if err != nil {
		return fmt.Errorf("can not hash password")
	}

	user := models.User{
		Email: registerUser.Email,
		Password: hashPassword,
		Name: registerUser.Name,
		Role: "Customer",
	}

	if err := db.DB.Create(&user).Error; err != nil {
		return fmt.Errorf("error could not create user %v", err)
	}

	return nil
}

func Login(loginUser models.User) (string, error) {
	user, err := queryUser(loginUser.Email)
	if err != nil {
		return "", fmt.Errorf("email not found")
	}

	ok := utils.CheckPasswordHash(loginUser.Password, user.Password)
	if !ok {
		return "", fmt.Errorf("wrong password")
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func UpdateUser(user models.User) error {
	updates := map[string]interface{}{
        "name":    user.Name,
        "phone":   user.Phone,
        "address": user.Address,
		"avatar": user.Avatar,
    }

	err := db.DB.Model(&models.User{}).
		Where("email = ?", user.Email).
		Updates(updates).Error

	if err != nil {
		return fmt.Errorf("failed to update user info %v", err)
	}

	return nil
}

func GetUserByEmail(email string) (models.User, error) {
	return queryUser(email)
}

func queryUser(email string) (models.User, error) {
	var user models.User
	err := db.DB.Where("email = ? AND role = ?", email, "Customer").First(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func queryUserID(email string) (string, error) {
	var userID string

	result := db.DB.Table("users").
		Select("id").
		Where("email = ?", email).
		Scan(&userID)

	if result.RowsAffected == 0 {
		return "", fmt.Errorf("user not found with email %s", email)
	}

	if result.Error != nil {
		return "", fmt.Errorf("error querying user ID")
	}

	return userID, nil
}

func GetCart(email string) ([]map[string]interface{}, error) {
	return queryCart(email)
}

func GetCartTotalPrice(email string) (int, error) {
	userID, err := queryUserID(email)
	if err != nil {
		return 0, err
	}

	return queryCartTotalPrice(userID)
}

func AddToCart(email string, productID string, quantity int) error {
	stock, err := queryProductQuantity(productID)
	if err != nil {
		return err
	}

	if stock < quantity {
		return fmt.Errorf("not enough stock")
	}

	userID, err := queryUserID(email)
	if err != nil {
		return err
	}

	currentQuantity, err := queryCartByProduct(userID, productID)
	if err != nil {
		return err
	}

	newQuantity := currentQuantity + quantity

	if stock < newQuantity {
		return fmt.Errorf("not enough stock")
	}

	if currentQuantity == 0 {
		err = insertCart(userID, productID, quantity)
		if err != nil {
			return err
		}
	} else {
		err = updateCartQuantity(userID, productID, newQuantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateCartItem(email string, productID string, quantity int) error {
	userID, err := queryUserID(email)
	if err != nil {
		return err
	}

	stock, err := queryProductQuantity(productID)
	if err != nil {
		return err
	}

	if stock < quantity {
		return fmt.Errorf("not enough stock")
	}

	if quantity == 0 {
		err = removeProductFromCart(userID, productID)
		if err != nil {
			return err
		}
	} else {
		err = updateCartQuantity(userID, productID, quantity)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveFromCart(email string, productID string) error {
	userID, err := queryUserID(email)
	if err != nil {
		return err
	}

	return removeProductFromCart(userID, productID)
}

func queryCart(email string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := db.DB.Table("cart").
		Select("cart.product_id, products.name, cart.quantity, products.price, products.image_url, products.category_id").
		Joins("JOIN users ON cart.user_id = users.id").
		Joins("JOIN products ON cart.product_id = products.id").
		Where("users.email = ?", email).
		Order("cart.updated_at DESC").
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to query cart products: %v", err)
	}

	if results == nil {
		return []map[string]interface{}{}, nil
	}

	return results, nil
}

func queryCartByProduct(userID string, productID string) (int, error) {
	var quantity int

	result := db.DB.Table("cart").
		Select("quantity").
		Where("user_id = ? AND product_id = ?", userID, productID).
		Scan(&quantity)

	if result.RowsAffected == 0 {
		return 0, nil
	}

	if result.Error != nil {
		return -1, result.Error
	}

	return quantity, nil
}

func queryCartTotalPrice(userID string) (int, error) {
	var totalPrice int

	err := db.DB.Table("cart").
		Select("COALESCE(SUM(cart.quantity * products.price), 0) AS total_price").
		Joins("JOIN products ON cart.product_id = products.id").
		Where("cart.user_id = ?", userID).
		Scan(&totalPrice).Error

	if err != nil {
		return 0, fmt.Errorf("failed to calculate total price")
	}

	return totalPrice, nil
}
 
func insertCart(userID string, productID string, quantity int) error {
	cart := models.Cart{
		UserID: uuid.MustParse(userID),
		ProductID: uuid.MustParse(productID),
		Quantity: quantity,
	}
	
	if err := db.DB.Create(&cart).Error; err != nil {
		return fmt.Errorf("failed to insert into cart: %v", err)
	}
	return nil
}

func updateCartQuantity(userID string, productID string, quantity int) error {
	if err := db.DB.Model(&models.Cart{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Update("quantity", quantity).Error; err != nil {
		return fmt.Errorf("failed to update cart %v", err)
	}

	return nil
}

func removeProductFromCart(userID string, productID string) error {
	result := db.DB.Where("user_id = ? AND product_id = ?", userID, productID).Delete(&models.Cart{})
	if result.Error != nil {
		return fmt.Errorf("failed to remove product from cart %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found in cart")
	}
	return nil
}

func queryProductQuantity(productID string) (int, error) {
	var quantity int

	result := db.DB.Table("products").
		Select("quantity").
		Where("id = ?", productID).
		Scan(&quantity)

	if result.RowsAffected == 0 {
		return -1, fmt.Errorf("product not found with id %s", productID)
	}

	if result.Error != nil {
		return -1, fmt.Errorf("error querying product %v", result.Error)
	}

	return quantity, nil
}

