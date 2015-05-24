// Package facade provides a convenient way to frontend http handlers.
package facade

import "net/http"

// A constructor for a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New
type Constructor func(http.Handler) http.Handler

// Frontend acts as a list of http.Handler constructors.
// Frontend is effectively immutable:
// once created, it will always hold
// the same set of constructors in the same order.
type Frontend struct {
	constructors []Constructor
}

// New creates a new frontend,
// memorizing the given list of middleware constructors.
// New serves no other function,
// constructors are only called upon a call to Then().
func New(constructors ...Constructor) Frontend {
	c := Frontend{}
	c.constructors = append(c.constructors, constructors...)

	return c
}

// Then chains the middleware and returns the final http.Handler.
//     New(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
//
// A frontend can be safely reused by calling Then() several times.
//     stdStack := facade.New(ratelimitHandler, csrfHandler)
//     indexPipe = stdStack.Then(indexHandler)
//     authPipe = stdStack.Then(authHandler)
// Note that constructors are called on every call to Then()
// and thus several instances of the same middleware will be created
// when a frontend is reused in this way.
// For proper middleware, this should cause no problems.
//
// Then() treats nil as http.DefaultServeMux.
func (c Frontend) Then(h http.Handler) http.Handler {
	var final http.Handler
	if h != nil {
		final = h
	} else {
		final = http.DefaultServeMux
	}

	for i := len(c.constructors) - 1; i >= 0; i-- {
		final = c.constructors[i](final)
	}

	return final
}

// ThenFunc works identically to Then, but takes
// a HandlerFunc instead of a Handler.
//
// The following two statements are equivalent:
//     c.Then(http.HandlerFunc(fn))
//     c.ThenFunc(fn)
//
// ThenFunc provides all the guarantees of Then.
func (c Frontend) ThenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(http.HandlerFunc(fn))
}

// Append extends a frontend, adding the specified constructors
// as the last ones in the request flow.
//
// Append returns a new frontend, leaving the original one untouched.
//
//     stdFrontend := facade.New(m1, m2)
//     extFrontend := stdFrontend.Append(m3, m4)
//     // requests in stdFrontend go m1 -> m2
//     // requests in extFrontend go m1 -> m2 -> m3 -> m4
func (c Frontend) Append(constructors ...Constructor) Frontend {
	newCons := make([]Constructor, len(c.constructors)+len(constructors))
	copy(newCons, c.constructors)
	copy(newCons[len(c.constructors):], constructors)

	newFrontend := New(newCons...)
	return newFrontend
}
