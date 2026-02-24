package api

import (
	"database/sql"

	"github.com/Adriusops/zentxt/internal/versioning"

	"github.com/gofiber/fiber/v3"
)

type CreateFileRequest struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func SetupRoutes(app fiber.Router, db *sql.DB) {
	app.Post("/files", func(c fiber.Ctx) error {
		// 1. Définir et lire le body
		var req CreateFileRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		// 2. Appeler CreateFile
		file, err := versioning.CreateFile(db, req.Name, req.Path, nil)
		if err != nil {
			return err
		}
		// 3. Retourner le résultat en JSON
		return c.JSON(file)
	})

	app.Post("/files/:id/versions", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "save version"})
	})

	app.Get("/files/:id/versions", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "list versions"})
	})
}
