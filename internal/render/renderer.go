package render

import (
	"html/template"
	"net/http"
)

type Renderer struct {
	tmpls map[string]*template.Template
	entry string
}

var renderer *Renderer

func Init(t map[string]*template.Template, entry string) {
	renderer = &Renderer{
		tmpls: t,
		entry: entry,
	}
}

// Execute executes a template with the given data and writes the result to the response writer.
func Execute(w http.ResponseWriter, name string, data interface{}) {
	if renderer == nil {
		http.Error(w, "template not initialized", http.StatusInternalServerError)
	}
	tmpl, ok := renderer.tmpls[name]
	if !ok {
		http.Error(w, "template not found", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// err := tmpl.Execute(w, data)
	err := tmpl.ExecuteTemplate(w, renderer.entry, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
