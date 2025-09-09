package main

import (
	"linn221/Requester/services"
	"linn221/Requester/views/templates"
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	db       *gorm.DB
	port     string
	secret   string
	services *services.ServiceContainer
}

func NewApp(db *gorm.DB, port, secret string) *App {
	return &App{
		db:       db,
		port:     port,
		secret:   secret,
		services: services.NewServiceContainer(db),
	}
}

func (a *App) HandleMin(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			w.Header().Add("HX-Reswap", "outerHTML")
			w.Header().Add("HX-Retarget", "#flash")
			templates.ErrorBox(err.Error()).Render(r.Context(), w)
			return
		}
	}
}
