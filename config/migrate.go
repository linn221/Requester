package config

import (
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {

	// err := db.AutoMigrate(&models.User{}, &models.Image{},
	// 	&models.Label{},
	// 	&models.Note{})
	// if err != nil {
	// 	panic("Error migrating: " + err.Error())
	// }
}
