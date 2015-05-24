// Package facade implements a middleware frontending solution.
package facade

import (
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

// TODO: test Frontend.New ingest distFilePath and save beforeContent and afterContent


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

	/*
	// TODO: figure out what to do with this now that we are using frontend_test.html AND move this into on func that all tests can use
	testBeforeContent := "<html><head><title>Test</title></head><body facade>"
	testAfterContent := "</body></html>"
	testTemplate := testBeforeContent + testAfterContent
	*/

	frontend := New("./frontend_test.html")
	assert.NotNil(t, frontend.Template)
}
