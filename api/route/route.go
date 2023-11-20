package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Ethea2/nat-dev/api/controller"
	"github.com/Ethea2/nat-dev/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	user := api.Group("/user")
	user.Post("/login", controller.Login)
	user.Post("/signup", controller.SignUp)

	project := api.Group("/project")
	project.Get("/", middleware.Protected(), controller.GetPosts)
}
