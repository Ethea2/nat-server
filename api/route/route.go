package route

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Ethea2/nat-server/api/controller"
	"github.com/Ethea2/nat-server/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	user := api.Group("/user")
	user.Post("/login", controller.Login)
	user.Post("/signup", controller.SignUp)

	exp := api.Group("/exp")
	exp.Post("/:id", controller.CreateExperience)
	exp.Get("/", controller.GetExperiences)

	project := api.Group("/project")
	project.Get("/", controller.GetProjects)
	project.Get("/:id", controller.GetSinglePost)
	project.Post("/create/:userID", middleware.Protected(), controller.CreateProject)
	project.Patch("/:id", middleware.Protected(), controller.UpdateProject)
	project.Patch("/image_edit/:projectID", middleware.Protected(), controller.UpdateProjectImage)
}
