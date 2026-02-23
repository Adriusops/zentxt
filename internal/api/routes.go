package api

import (
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app fiber.Router) {
	app.Post("/files", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "create file"})
	})

	app.Post("/files/:id/versions", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "save version"})
	})

	app.Get("/files/:id/versions", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "list versions"})
	})
}
