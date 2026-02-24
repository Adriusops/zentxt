package main

import (
	"log"

	"github.com/Adriusops/zentxt/internal/api"
	"github.com/Adriusops/zentxt/internal/storage"
	"github.com/gofiber/fiber/v3"
)

func main() {

	// Initialize the database connection
	db, err := storage.InitDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a new Fiber app
	app := fiber.New()

	// Define a route for the GET method on the root path '/'
	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	api.SetupRoutes(app, db)

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))
}
