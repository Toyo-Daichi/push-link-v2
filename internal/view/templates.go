package view

import (
	"html/template"
	"io"
	"path/filepath"
	"time"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewTemplateRenderer(pattern string) (*TemplateRenderer, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	templates, err := template.New("root").Funcs(template.FuncMap{
		"formatDate": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}
			return t.Format("2006-01-02")
		},
		"add1": func(value int) int {
			return value + 1
		},
	}).ParseFiles(matches...)
	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{templates: templates}, nil
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data any) error {
	return r.templates.ExecuteTemplate(w, name, data)
}
