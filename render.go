package templates

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

type Renderer struct {
	viewsDir     string
	layout       string
	cache        map[string]*template.Template
	mu           sync.RWMutex
	isProduction bool
}

func NewRenderer(viewsDir string, isProduction bool) *Renderer {
	return &Renderer{
		viewsDir:     viewsDir,
		layout:       "_layout.html",
		cache:        make(map[string]*template.Template),
		isProduction: isProduction,
	}
}

func (r *Renderer) Render(w http.ResponseWriter, page string, data interface{}) {
	if r.isProduction {
		r.mu.RLock()
		tmpl, ok := r.cache[page]
		r.mu.RUnlock()

		if !ok {
			tmpl, err := r.parseTemplateFiles(page)
			if err != nil {
				http.Error(w, "Error parsing templates: "+err.Error(), http.StatusInternalServerError)
				return
			}

			r.mu.Lock()
			r.cache[page] = tmpl
			r.mu.Unlock()
		}

		r.mu.RLock()
		tmpl = r.cache[page]
		r.mu.RUnlock()

		err := tmpl.ExecuteTemplate(w, r.layout, map[string]interface{}{
			"Page": page,
			"Data": data,
		})
		if err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		}
	} else {
		tmpl, err := r.parseTemplateFiles(page)
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

// 🔥 Acá se parsea todo: partials, layout y page específica
func (r *Renderer) parseTemplateFiles(page string) (*template.Template, error) {
	layoutPath := filepath.Join(r.viewsDir, r.layout)
	pagePath := filepath.Join(r.viewsDir, page)
	fmt.Printf("partials/*.html:%v\n", pagePath)
	partialsGlob := "partials/*.html"
	partials, err := filepath.Glob(partialsGlob)
	fmt.Printf("partials/*.html:%v\n", partials)
	if err != nil {
		return nil, err
	}

	// Si no hay partials, igualamos a slice vacío para no fallar
	if partials == nil {
		partials = []string{}
	}

	// Parsear todos juntos
	allFiles := append(partials, layoutPath, pagePath)
	return template.ParseFiles(allFiles...)
}
