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

func (a *App) HandleMin(h func(v *views.MyTemplates, w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	v := a.templates
	return func(w http.ResponseWriter, r *http.Request) {

		err := h(v, w, r)
		if err != nil {
			w.Header().Add("HX-Reswap", "outerHTML")
			w.Header().Add("HX-Retarget", "#flash")
			v.ShowErrorBox(w, err.Error())
			return
		}
	}
}
