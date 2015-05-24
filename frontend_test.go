// Package facade implements a middleware frontending solution.
package facade

import (
	"net/http"
	"net/http/httptest"
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
	c1 := func(h http.Handler) http.Handler {
		return nil
	}
	c2 := func(h http.Handler) http.Handler {
		return http.StripPrefix("potato", nil)
	}

	slice := []Constructor{c1, c2}

	frontend := New(slice...)
	assert.Equal(t, frontend.constructors[0], slice[0])
	assert.Equal(t, frontend.constructors[1], slice[1])
}

func TestThenWorksWithNoMiddleware(t *testing.T) {
	assert.NotPanics(t, func() {
		frontend := New()
		final := frontend.Then(testApp)

		assert.Equal(t, final, testApp)
	})
}

func TestThenTreatsNilAsDefaultServeMux(t *testing.T) {
	chained := New().Then(nil)
	assert.Equal(t, chained, http.DefaultServeMux)
}

func TestThenFuncTreatsNilAsDefaultServeMux(t *testing.T) {
	chained := New().ThenFunc(nil)
	assert.Equal(t, chained, http.DefaultServeMux)
}

func TestThenOrdersHandlersRight(t *testing.T) {
	t1 := tagMiddleware("t1\n")
	t2 := tagMiddleware("t2\n")
	t3 := tagMiddleware("t3\n")

	chained := New(t1, t2, t3).Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTP(w, r)

	assert.Equal(t, w.Body.String(), "t1\nt2\nt3\napp\n")
}

func TestAppendAddsHandlersCorrectly(t *testing.T) {
	frontend := New(tagMiddleware("t1\n"), tagMiddleware("t2\n"))
	newFrontend := frontend.Append(tagMiddleware("t3\n"), tagMiddleware("t4\n"))

	assert.Equal(t, len(frontend.constructors), 2)
	assert.Equal(t, len(newFrontend.constructors), 4)

	chained := newFrontend.Then(testApp)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	chained.ServeHTTP(w, r)

	assert.Equal(t, w.Body.String(), "t1\nt2\nt3\nt4\napp\n")
}

func TestAppendRespectsImmutability(t *testing.T) {
	frontend := New(tagMiddleware(""))
	newFrontend := frontend.Append(tagMiddleware(""))

	assert.NotEqual(t, &frontend.constructors[0], &newFrontend.constructors[0])
}
