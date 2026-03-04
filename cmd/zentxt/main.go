package main

import (
	"io/fs"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/Adriusops/zentxt"
	"github.com/Adriusops/zentxt/internal/api"
	"github.com/Adriusops/zentxt/internal/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/gofiber/template/html/v2"
)

func main() {

	// Initialize the database connection
	db, err := storage.InitDB(zentxt.Migrations)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a new Fiber app
	templatesFS, err := fs.Sub(zentxt.Templates, "templates")
	if err != nil {
		log.Fatal(err)
	}
	engine := html.NewFileSystem(http.FS(templatesFS), ".html")
	engine.Reload(true)
	engine.AddFunc("add", func(a, b int) int { return a + b })
	engine.AddFunc("sub", func(a, b int) int { return a - b })
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	staticFS, err := fs.Sub(zentxt.Static, "static")
	if err != nil {
		log.Fatal(err)
	}
	app.Use("/static", static.New("", static.Config{FS: staticFS}))

	api.SetupRoutes(app, db)

	time.Sleep(500 * time.Millisecond)
	switch runtime.GOOS {
	case "darwin":
		exec.Command("open", "http://localhost:3000").Start()
	case "windows":
		exec.Command("cmd", "/c", "start", "http://localhost:3000").Start()
	case "linux":
		exec.Command("xdg-open", "http://localhost:3000").Start()
	}

	// Start the server on port 3000
	log.Fatal(app.Listen(":3000"))

}
