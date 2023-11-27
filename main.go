package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/Ethea2/nat-server/api/route"
	"github.com/Ethea2/nat-server/database"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	err := database.ConnectDB()
	if err != nil {
		log.Fatal("Error!", err)
	}

	route.SetupRoutes(app)

	app.Listen(":4000")
	database.CloseDB()
}
