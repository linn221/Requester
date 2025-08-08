package web

import (
	"linn221/Requester/config"
	"linn221/Requester/utils"
	"linn221/Requester/views"
	"net/http"
	"path/filepath"
)

func main() {
	baseDir := config.GetBaseDir()
	db := config.ConnectDB()
	templates := views.NewTemplates(filepath.Join(baseDir, "../../views/templates"))
	secretKey := utils.GenerateRandomString(20)
	auth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		})
	}
	app := App{
		db:        db,
		templates: templates,
		port:      "8080",
	}

}
