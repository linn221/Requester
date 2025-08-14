package main

import (
	"linn221/Requester/views"
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	db             *gorm.DB
	port           string
	secret         string
	templates      *views.MyTemplates
	authMiddleware func(http.Handler) http.Handler
}
