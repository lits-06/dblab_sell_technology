package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/lits-06/sell_technology/internal/app/routes"
	"github.com/lits-06/sell_technology/pkg/db"
	"github.com/lits-06/sell_technology/pkg/utils"
)

func main() {
	utils.InitLogger()
	
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	db.Connect()
	db.Migrate()

	router := routes.SetupRoute()
	log.Println("Server is running on http://localhost:8080")
	router.Run(":8080")
}