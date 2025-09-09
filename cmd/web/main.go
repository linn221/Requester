package main

import (
	"linn221/Requester/config"
	"linn221/Requester/utils"
)

func main() {
	// Connect to database
	db := config.ConnectDB()
	secret := utils.GenerateRandomString(20)

	app := NewApp(db, "8080", secret)
	app.Serve()
}
