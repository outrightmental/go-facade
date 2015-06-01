// Package facade implements a middleware frontending solution.
package facade

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

// A constructor for middleware
// that writes its own "tag" into the RW and does nothing else.
// Useful in checking if a frontend is behaving in the right order.
func tagMiddleware(tag string) Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(tag))
			h.ServeHTTP(w, r)
		})
	}
}

var testApp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("app\n"))
})

// Tests creating a new frontend
func TestNew(t *testing.T) {
	frontend := New("./frontend_test.html")
	assert.NotNil(t, frontend.template)
}

// TODO: Test new frontend with incorrect HTML file panics

// TODO: Render returns correct output html
