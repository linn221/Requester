package services

import (
	"gorm.io/gorm"
)

// ServiceContainer holds all services
type ServiceContainer struct {
	ImportService    *ImportService
	RequestService   *RequestService
	ImportJobService *ImportJobService
	EndpointService  *EndpointService
	FormParser       *FormParser
}

// NewServiceContainer creates a new service container with all services
func NewServiceContainer(db *gorm.DB) *ServiceContainer {
	database := NewGormDatabaseAdapter(db)
	endpointService := NewEndpointService(database)

	return &ServiceContainer{
		ImportService:    NewImportService(database, endpointService),
		RequestService:   NewRequestService(database),
		ImportJobService: NewImportJobService(database),
		EndpointService:  endpointService,
		FormParser:       NewFormParser(),
	}
}
