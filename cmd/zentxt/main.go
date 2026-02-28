package main

import (
	"log"

	"github.com/Adriusops/zentxt/internal/api"
	"github.com/Adriusops/zentxt/internal/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/gofiber/template/html/v2"
)

func main() {

	// Initialize the database connection
	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a new Fiber app
	engine := html.New("./templates", ".html")
	engine.Reload(true)
	engine.AddFunc("add", func(a, b int) int { return a + b })
	engine.AddFunc("sub", func(a, b int) int { return a - b })
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(static.New("./static"))

	api.SetupRoutes(app, db)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
