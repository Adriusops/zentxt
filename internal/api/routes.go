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

	app.Get("/", func(c fiber.Ctx) error {
		files, err := versioning.ListFiles(db)
		if err != nil {
			return err
		}
		return c.Render("home", fiber.Map{
			"Title": "ZenTxt - Your Files",
			"files": files,
		})
	})

	app.Get("api/files", func(c fiber.Ctx) error {
		files, err := versioning.ListFiles(db)
		if err != nil {
			return err
		}
		return c.JSON(files)
	})

	app.Post("api/files", func(c fiber.Ctx) error {
		// 1. Définir et lire le body
		var req CreateFileRequest
		if err := c.Bind().Body(&req); err != nil {
			return err
		}

		if req.Name == "" || req.Path == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name and path are required"})
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
		fileID := c.Params("id")
		file, err := versioning.GetFile(db, fileID)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).SendString("File not found")
			}
			return err
		}
		versions, err := versioning.ListVersions(db, fileID)
		if err != nil {
			return err
		}

		// Create a map to store previous version IDs
		type VersionWithPrev struct {
			*versioning.Version
			PrevVersionID string
		}

		versionsWithPrev := make([]*VersionWithPrev, len(versions))
		for i, v := range versions {
			prevID := ""
			if i > 0 {
				prevID = versions[i-1].ID
			}
			versionsWithPrev[i] = &VersionWithPrev{
				Version:       v,
				PrevVersionID: prevID,
			}
		}

		return c.Render("timeline", fiber.Map{
			"Title":    file.Name + " - Timeline",
			"file":     file,
			"versions": versionsWithPrev,
		})
	})

	app.Post("api/files/:id/versions", func(c fiber.Ctx) error {
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

	app.Get("api/files/:id/versions", func(c fiber.Ctx) error {
		id := c.Params("id")
		versions, err := versioning.ListVersions(db, id)
		if err != nil {
			return err
		}
		return c.JSON(versions)
	})

	app.Get("api/files/:id/versions/:version_id", func(c fiber.Ctx) error {
		versionID := c.Params("version_id")
		version, err := versioning.GetVersion(db, versionID)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		return c.JSON(version)
	})

	app.Get("api/files/:id/diff", func(c fiber.Ctx) error {
		v1 := c.Query("v1")
		v2 := c.Query("v2")

		version1, err := versioning.GetVersion(db, v1)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		version2, err := versioning.GetVersion(db, v2)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		diff := versioning.GenerateDiff(version1.Content, version2.Content)

		return c.JSON(diff)
	})

	app.Patch("api/files/:id/restore/:versionID", func(c fiber.Ctx) error {
		fileID := c.Params("id")
		versionID := c.Params("versionID")

		restoredVersion, err := versioning.RestoreVersion(db, fileID, versionID)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}

		return c.JSON(restoredVersion)
	})

	// File timeline page
	app.Get("/files/:id", func(c fiber.Ctx) error {
		fileID := c.Params("id")
		file, err := versioning.GetFile(db, fileID)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).SendString("File not found")
			}
			return err
		}
		versions, err := versioning.ListVersions(db, fileID)
		if err != nil {
			return err
		}
		return c.Render("timeline", fiber.Map{
			"Title":    file.Name + " - Timeline",
			"file":     file,
			"versions": versions,
		})
	})

	app.Get("/files/:id/diff", func(c fiber.Ctx) error {
		fileID := c.Params("id")
		v1 := c.Query("v1")
		v2 := c.Query("v2")

		file, err := versioning.GetFile(db, fileID)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).SendString("File not found")
			}
			return err
		}

		version1, err := versioning.GetVersion(db, v1)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).SendString("Version not found")
			}
			return err
		}

		version2, err := versioning.GetVersion(db, v2)
		if err != nil {
			if err == versioning.ErrNotFound {
				return c.Status(fiber.StatusNotFound).SendString("Version not found")
			}
			return err
		}

		diff := versioning.GenerateDiff(version1.Content, version2.Content)

		return c.Render("diff", fiber.Map{
			"Title":    "Compare versions" + file.Name,
			"file":     file,
			"version1": version1,
			"version2": version2,
			"diff":     diff,
		})
	})

}
