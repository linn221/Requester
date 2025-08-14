package main

import (
	"linn221/Requester/utils"
	"linn221/Requester/views"
	"path/filepath"
)

func main() {
	// baseDir := config.GetBaseDir()
	// db := config.ConnectDB()
	baseDir := utils.GetBaseDir()
	templates := views.NewTemplates(filepath.Join(baseDir, "../../views/templates"))
	secret := utils.GenerateRandomString(20)

	app := App{
		// db:             db,
		templates: templates,
		port:      "8080",
		secret:    secret,
	}

	app.Serve()
}
