package main

import (
	"linn221/Requester/utils"
)

func main() {
	// baseDir := config.GetBaseDir()
	// db := config.ConnectDB()
	secret := utils.GenerateRandomString(20)

	app := App{
		// db:             db,
		port:      "8080",
		secret:    secret,
	}

	app.Serve()
}
