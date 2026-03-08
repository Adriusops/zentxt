package main

import (
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"

	"github.com/Adriusops/zentxt"
	"github.com/Adriusops/zentxt/internal/storage"
)

func main() {

	// Initialize the database connection
	db, err := storage.InitDB(zentxt.Migrations)
	if err != nil {
		log.Fatal(err)
	}

	app := application.New(application.Options{
		Name: "zentxt",
		Services: []application.Service{
			application.NewService(zentxt.NewApp(db)),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(zentxt.Assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "zentxt",
		URL:   "/",
	})

	err = app.Run()
}
