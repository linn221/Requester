package main

import (
	"fmt"
	"linn221/Requester/views"
	"log"
	"net/http"
	"runtime/debug"
)

func MakeAuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookiesSecret, err := r.Cookie("secret")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Redirect(w, r, "/secret-required", http.StatusTemporaryRedirect)
					return
				}
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if cookiesSecret.Value == secret {
				next.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, "/secret-required", http.StatusTemporaryRedirect)
			}
		})
	}
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				// Log the panic message and stack trace
				log.Printf("[PANIC RECOVERED] %v\n%s", rec, debug.Stack())

				// Optional: customize the error response
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func makeDashboardRoutes(app *App, mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		app.templates.HomePage(w)
	})
	mux.HandleFunc("GET /import", app.HandleMin(func(v *views.MyTemplates, w http.ResponseWriter, r *http.Request) error {
		return v.ShowImportForm(w)
	}))
}

func (a *App) Serve() {
	mux := http.NewServeMux()
	authMux := http.NewServeMux()
	authMiddleware := MakeAuthMiddleware(a.secret)
	makeDashboardRoutes(a, authMux)

	mux.HandleFunc("GET /start-session", func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("secret")
		if s != "" {
			// set cookies
			http.SetCookie(w, &http.Cookie{
				Name:   "secret",
				Value:  s,
				MaxAge: 10 * 60 * 60,
				Path:   "/", Domain: "",
				Secure: false, HttpOnly: true,
			})

			http.Redirect(w, r, "/dashboard/", http.StatusTemporaryRedirect) // the final trail is important for some reason
			return
		}
	})

	mux.HandleFunc("GET /secret-required", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Secrets required. Open the link from running the app"))
	})

	mux.Handle("/dashboard/", http.StripPrefix("/dashboard", authMiddleware(authMux))) // the final trail is important

	fmt.Printf("http://localhost:%s/start-session?secret=%s\n", a.port, a.secret)
	err := http.ListenAndServe(":"+a.port, RecoveryMiddleware(mux))
	if err != nil {
		log.Fatal(err)
	}
}
