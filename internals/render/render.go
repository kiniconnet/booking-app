package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/justinas/nosurf"
	"github.com/kiniconnect/booking-app/internals/config"
	"github.com/kiniconnect/booking-app/internals/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template renders a template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	// Get the absolute path to the templates directory
	absPath, err := filepath.Abs("./templates")
	if err != nil {
		return myCache, fmt.Errorf("error getting absolute path: %v", err)
	}

	pages, err := filepath.Glob(filepath.Join(absPath, "*.page.html"))
	if err != nil {
		return myCache, fmt.Errorf("error finding page templates: %v", err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, fmt.Errorf("error parsing template %s: %v", name, err)
		}

		matches, err := filepath.Glob(filepath.Join(absPath, "*.layout.html"))
		if err != nil {
			return myCache, fmt.Errorf("error finding layout templates: %v", err)
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(filepath.Join(absPath, "*.layout.html"))
			if err != nil {
				return myCache, fmt.Errorf("error parsing layout templates: %v", err)
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
