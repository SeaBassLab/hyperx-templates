package templates

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type Renderer struct {
	viewsDir    string
	layout      string
	cache       map[string]*template.Template
	mu          sync.RWMutex
	isProduction bool
}

func NewRenderer(viewsDir string, isProduction bool) *Renderer {
	return &Renderer{
		viewsDir:    viewsDir,
		layout:      "_layout.html",
		cache:       make(map[string]*template.Template),
		isProduction: isProduction,
	}
}

func (r *Renderer) Render(w http.ResponseWriter, page string, data interface{}) {
	if r.isProduction {
		r.mu.RLock()
		tmpl, ok := r.cache[page]
		r.mu.RUnlock()

		if !ok {
			// No está en cache, parsear y guardar
			layoutPath := filepath.Join(r.viewsDir, r.layout)
			pagePath := filepath.Join(r.viewsDir, page)

			var err error
			tmpl, err = template.ParseFiles(layoutPath, pagePath)
			if err != nil {
				http.Error(w, "Error parsing templates: "+err.Error(), http.StatusInternalServerError)
				return
			}

			r.mu.Lock()
			r.cache[page] = tmpl
			r.mu.Unlock()
		}

		// Ejecutar template cacheado
		err := tmpl.ExecuteTemplate(w, r.layout, map[string]interface{}{
			"Page": page,
			"Data": data,
		})
		if err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}

	} else {
		// Dev mode: parsear siempre para ver cambios sin reiniciar
		layoutPath := filepath.Join(r.viewsDir, r.layout)
		pagePath := filepath.Join(r.viewsDir, page)

		tmpl, err := template.ParseFiles(layoutPath, pagePath)
		if err != nil {
			http.Error(w, "Error parsing templates: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.ExecuteTemplate(w, r.layout, map[string]interface{}{
			"Page": page,
			"Data": data,
		})
		if err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
	}
}
