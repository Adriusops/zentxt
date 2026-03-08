package main

import (
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"

	"github.com/Adriusops/zentxt"
	"github.com/Adriusops/zentxt/internal/storage"
)

func main() {

	// Initialize the database connection
	db, err := storage.InitDB(zentxt.Migrations)
	if err != nil {
		log.Fatal(err)
	}

	app := zentxt.NewApp(db)

	err = wails.Run(&options.App{
		Title:  "zentxt",
		Width:  1024,
		Height: 768,
		Bind: []interface{}{
			app,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

}
