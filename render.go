package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	templates *template.Template
}

func NewRenderer(viewsDir string) *Renderer {
	tmpl := template.Must(template.ParseGlob(filepath.Join(viewsDir, "*.html")))
	return &Renderer{templates: tmpl}
}

func (r *Renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	err := r.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
