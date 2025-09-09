package config

import (
	"linn221/Requester/requests"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&requests.Endpoint{}, &requests.ImportJob{}, &requests.MyRequest{})
	if err != nil {
		panic("Error migrating: " + err.Error())
	}
}
