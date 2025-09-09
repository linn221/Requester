package services

import (
	"gorm.io/gorm"
)

// ServiceContainer holds all services
type ServiceContainer struct {
	ImportService    *ImportService
	RequestService   *RequestService
	ImportJobService *ImportJobService
	FormParser       *FormParser
}

// NewServiceContainer creates a new service container with all services
func NewServiceContainer(db *gorm.DB) *ServiceContainer {
	database := NewGormDatabaseAdapter(db)
	
	return &ServiceContainer{
		ImportService:    NewImportService(database),
		RequestService:   NewRequestService(database),
		ImportJobService: NewImportJobService(database),
		FormParser:       NewFormParser(),
	}
}
