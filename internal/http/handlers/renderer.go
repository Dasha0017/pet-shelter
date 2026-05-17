package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Renderer struct {
	templatesDir string
}

func NewRenderer(templatesDir string) *Renderer {
	return &Renderer{templatesDir: templatesDir}
}

func (r *Renderer) Render(w http.ResponseWriter, page string, data any) error {
	files := []string{
		filepath.Join(r.templatesDir, "layout.html"),
		filepath.Join(r.templatesDir, page+".html"),
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "layout", data)
}
