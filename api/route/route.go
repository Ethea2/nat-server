package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Ethea2/nat-dev/api/controller"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	user := api.Group("/user")
	user.Post("/login", controller.Login)

	post := api.Group("/post")
	post.Get("/", controller.GetPosts)
}
