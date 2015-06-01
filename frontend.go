// Facade memorizes one static index.html to use as a minimal site template.
package facade

import (
  "io"
	"net/http"
	"html/template"
)

// A constructor for a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New
type Constructor func(http.Handler) http.Handler

// Frontend acts as a very-simple template engine.
// Frontend is effectively immutable:
// once created, it will always hold
// the same template.
type Frontend struct {
	template *template.Template
  distFilePath string
}

// New creates a new frontend,
// memorizing the given distfile.
// New serves no other function,
// template is only manipulated during a call to Render().
func New(distFilePath string) Frontend {
	t, err := template.ParseFiles(distFilePath)
	if (err != nil) {
		panic(err)
	}
	c := Frontend{}
	c.template = t
	c.distFilePath = distFilePath
	return c
}

// Render generates output HTML
// using the memorized Template
func (f Frontend) Render(w io.Writer, innerHtml string) error {
  // TODO: find the HTML element with `facade` inside the memorized template
  // TODO: inject the innerHtml and save the outputHtml
  // TODO: panic any errors
  return f.template.Execute(w, innerHtml)
}
