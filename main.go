package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/Ethea2/nat-dev/api/route"
	"github.com/Ethea2/nat-dev/database"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	err := database.ConnectDB()
	if err != nil {
		log.Fatal("Error!", err)
	}

	route.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	app.Listen("0.0.0.0:" + port)
	database.CloseDB()
}
