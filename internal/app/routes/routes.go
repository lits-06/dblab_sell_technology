package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lits-06/sell_technology/internal/app/controllers"
	"github.com/lits-06/sell_technology/internal/app/middleware"
)

func SetupRoute() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:           12 * 3600,
	}))

	users := router.Group("/api/users")
	{
		users.POST("/register", controllers.Register)
		users.POST("/login", controllers.Login)
	}

	usersProtected := router.Group("/api/users")
	usersProtected.Use(middleware.Auth())
	{
		usersProtected.GET("/info", controllers.GetUserInfo)
		usersProtected.PUT("/info", controllers.UpdateUser)
		usersProtected.GET("/cart", controllers.GetUserCart)
		usersProtected.PUT("/cart", controllers.UpdateCartItem)
		usersProtected.GET("/cart/totalprice", controllers.GetUserCartTotalPrice)
		usersProtected.POST("/cart", controllers.AddToCart)
		usersProtected.DELETE("/cart/:product_id", controllers.RemoveFromCart)
		usersProtected.GET("/order", controllers.GetOrder)
		usersProtected.POST("/order", controllers.CreateOrder)
		usersProtected.PUT("/order/cancel", controllers.CancelOrder)
	}

	products := router.Group("/api/products")
	{
		products.GET("/", controllers.GetProducts)
		products.GET("/:id", controllers.GetProductsByID)
	}

	category := router.Group("/api/category")
	{
		category.GET("/", controllers.GetCategories)
		category.GET("/:name", controllers.GetCategoryByName)
	}

	return router
}