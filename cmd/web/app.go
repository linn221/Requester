package web

import (
	"linn221/Requester/views"
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	db             *gorm.DB
	port           string
	templates      *views.MyTemplates
	authMiddleware func(http.Handler) http.Handler
}
