package controller

import (
	"github.com/gofiber/fiber/v2"
)

func GetPosts(c *fiber.Ctx) error {
	return c.SendString("Hello Posts!")
}
