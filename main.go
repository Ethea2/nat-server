package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/Ethea2/nat-dev/api/route"
	"github.com/Ethea2/nat-dev/database"
)

func main() {
	path, oerr := os.Getwd()
	if oerr != nil {
		log.Println(oerr)
	}

	entries, readerr := os.ReadDir(path)
	if readerr != nil {
		log.Fatal(readerr)
	}

	for _, e := range entries {
		fmt.Println(e.Name())
	}

	p, perr := os.Executable()
	if perr != nil {
		log.Fatal(perr)
	}

	godotenv.Load(p, ".env")

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
