package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type MyTemplates struct {
	home       *template.Template
	importFile *template.Template
	jobStatus  *template.Template
	errorBox   *template.Template
}

func NewTemplates(templateDir string) *MyTemplates {

	return &MyTemplates{
		home:       template.Must(template.ParseFiles(filepath.Join(templateDir, "index.html"))),
		importFile: template.Must(template.ParseFiles(filepath.Join(templateDir, "importForm.html"))),
		jobStatus:  template.Must(template.ParseFiles(filepath.Join(templateDir, "job.html"))),
		errorBox:   template.Must(template.ParseFiles(filepath.Join(templateDir, "errorbox.html"))),
	}
}

func (t *MyTemplates) HomePage(w http.ResponseWriter) error {
	return t.home.Execute(w, nil)
}

func (t *MyTemplates) ShowImportForm(w http.ResponseWriter) error {
	return t.home.Execute(w, nil)
}

func (t *MyTemplates) ShowErrorBox(w http.ResponseWriter, s string) {
	t.errorBox.Execute(w, map[string]any{
		"Message": s,
	})
}
