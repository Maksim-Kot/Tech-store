package http

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	catalogmodel "github.com/Maksim-Kot/Tech-store-catalog/pkg/model"
	"github.com/Maksim-Kot/Tech-store-web/internal/model"
	"github.com/Maksim-Kot/Tech-store-web/ui"
)

type templateData struct {
	CurrentYear     int
	Categories      []*catalogmodel.Category
	Products        []*catalogmodel.Product
	Product         *processedProduct
	Form            any
	Flash           string
	IsAuthenticated bool
	User            *model.User
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
