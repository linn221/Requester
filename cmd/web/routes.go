package web

import (
	"log"
	"net/http"
)

func (a *App) Serve() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {

	})

	err := http.ListenAndServe(":"+a.port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
