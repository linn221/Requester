package config

import (
	"linn221/Requester/requests"

	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	// Disable foreign key checks temporarily
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	// Drop existing foreign key constraints if they exist
	db.Exec("ALTER TABLE endpoints DROP FOREIGN KEY IF EXISTS fk_programs_endpoints")
	db.Exec("ALTER TABLE import_jobs DROP FOREIGN KEY IF EXISTS fk_programs_import_jobs")
	db.Exec("ALTER TABLE my_requests DROP FOREIGN KEY IF EXISTS fk_programs_my_requests")

	// First, create the Program table
	err := db.AutoMigrate(&requests.Program{})
	if err != nil {
		panic("Error migrating Program table: " + err.Error())
	}

	// Then migrate the other tables
	err = db.AutoMigrate(&requests.Endpoint{}, &requests.ImportJob{}, &requests.MyRequest{})
	if err != nil {
		panic("Error migrating other tables: " + err.Error())
	}
}
