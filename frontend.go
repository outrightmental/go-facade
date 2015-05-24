// Package facade provides a convenient way to frontend http handlers.
package facade

import "net/http"

// A constructor for a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New
type Constructor func(http.Handler) http.Handler

// Frontend acts as a very-simple template engine.
// Frontend is effectively immutable:
// once created, it will always hold
// the same template.
type Frontend struct {
	beforeContent string
	afterContent string
}

// New creates a new frontend,
// memorizing the given distfile.
// New serves no other function,
// output is only built upon a call to Write().
func New(distFilePath string) Frontend {
	c := Frontend{}
//	c.distFilePath = distFilePath
	// TODO: ingest distFilePath and save beforeContent and afterContent
	return c
}
