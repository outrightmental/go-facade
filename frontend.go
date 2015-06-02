// Facade memorizes one static index.html to use as a minimal site template.
package facade

import (
  "net/http"
  "fmt"
  "io"
  "io/ioutil"
  "regexp"
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
	templateHtml []byte
  distFilePath string
  replacements map[string]FrontendReplacement
}

// FrontendReplacement stores one operation.
// A Frontend has a chain of FrontendReplacement that
// will all be run on the Casing provided at Write time
type FrontendReplacement struct {
  regex *regexp.Regexp
  repl string
}

// A Casing is a set of values from which
// to compile a one-off html output
type Casing map[string]interface{}

// New creates a new frontend,
// memorizing the given distfile.
// New serves no other function,
// template is only manipulated during a call to Render().
func New(d string) *Frontend {
	t, err := ioutil.ReadFile(d)
	if (err != nil) {
		panic(err)
	}
	return &Frontend{
    templateHtml: t,
    distFilePath: d,
    replacements: make(map[string]FrontendReplacement,0),
  }
}

// RegexReplaceAll modifies the Template HTML, replacing matches of the Regexp
// with the replacement text repl.  Inside repl, $ signs are interpreted as
// in Expand, so for instance $1 represents the text of the first submatch.
func (f *Frontend) PreReplaceAll(exp string, repl string) error {
  r, err := regexp.Compile(exp)
  if (err != nil) {
    return err
  }
  t := r.ReplaceAllString(string(f.templateHtml), repl)
  f.templateHtml = []byte(t)
  return nil
}

//
func (f *Frontend) WillReplaceAll(key string, exp string, repl string) error {
  r, err := regexp.Compile(exp)
  if (err != nil) {
    return err
  }
  f.replacements[key] = FrontendReplacement{
      regex: r,
      repl: repl,
    }
  return nil
}

// Render generates output HTML
// using the memorized Template
func (f *Frontend) Render(w io.Writer, c Casing) error {
  o := f.templateHtml
  for key, r := range f.replacements {
    o = []byte( r.regex.ReplaceAllString( string(o), fmt.Sprintf(r.repl, c[key])))
  }
  w.Write(o)
  return nil
}

func (f *Frontend) GetTemplateLength() int {
  return len(f.templateHtml)
}
