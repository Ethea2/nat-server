package controller

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	fmt.Print("Hello!")

	return c.SendString("Hello!")
}
