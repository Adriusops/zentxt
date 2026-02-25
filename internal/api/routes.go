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

type SaveVersionRequest struct {
	Path    string `json:"path"`
	Author  string `json:"author"`
	Message string `json:"message"`
	Content string `json:"content"`
}

func SetupRoutes(app fiber.Router, db *sql.DB) {

	app.Get("/files", func(c fiber.Ctx) error {
		files, err := versioning.ListFiles(db)
		if err != nil {
			return err
		}
		return c.JSON(files)
	})

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

	app.Get("/files/:id", func(c fiber.Ctx) error {
		id := c.Params("id")
		file, err := versioning.GetFile(db, id)
		if err != nil {
			return err
		}
		return c.JSON(file)
	})

	app.Post("/files/:id/versions", func(c fiber.Ctx) error {
		var req SaveVersionRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}
		// 2. Appeler CreateFile
		version, err := versioning.SaveVersion(db, c.Params("id"), req.Path, req.Author, req.Message, req.Content)
		if err != nil {
			return err
		}
		// 3. Retourner le résultat en JSON
		return c.JSON(version)
	})

	app.Get("/files/:id/versions", func(c fiber.Ctx) error {
		id := c.Params("id")
		versions, err := versioning.ListVersions(db, id)
		if err != nil {
			return err
		}
		return c.JSON(versions)
	})

	app.Get("/files/:id/versions/:version_id", func(c fiber.Ctx) error {
		versionID := c.Params("version_id")
		version, err := versioning.GetVersion(db, versionID)
		if err != nil {
			return err
		}
		return c.JSON(version)
	})

	app.Get("/files/:id/diff", func(c fiber.Ctx) error {
		v1 := c.Query("v1")
		v2 := c.Query("v2")

		version1, err := versioning.GetVersion(db, v1)
		if err != nil {
			return err
		}
		version2, err := versioning.GetVersion(db, v2)
		if err != nil {
			return err
		}

		diff := versioning.GenerateDiff(version1.Content, version2.Content)

		return c.JSON(diff)
	})
}
