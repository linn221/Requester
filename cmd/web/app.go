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
	authMiddleware func(http.Handler) http.Handler
}

func (a *App) HandleMin(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err != nil {
			w.Header().Add("HX-Reswap", "outerHTML")
			w.Header().Add("HX-Retarget", "#flash")
			views.ErrorBox(err.Error()).Render(r.Context(), w)
			return
		}
	}
}
