// Package facade implements a middleware frontending solution.
package facade

import (
	"net/http"
	"testing"
  "io/ioutil"
  "regexp"
	"github.com/stretchr/testify/assert"
  "github.com/stretchr/testify/mock"
)

//
// Tests
//

func TestNewFrontend(t *testing.T) {
  frontend := New(pathToTestTemplateHtml)
  assert.NotNil(t, frontend.GetTemplateLength())
}

func TestNewFrontendFailsWithIncorrectPath(t *testing.T) {
  assert.Panics( t, panicNewFrontendWithIncorrectPath)
}

func TestRenderCorrectOutputHtml(t *testing.T) {
  frontend := New(pathToTestTemplateHtml)
  w := new(mockIoWriter)
  h, _ := readTemplateHtml()
  w.On("Write", h).Return(len(h),nil)
  frontend.Render(w, []byte(""))
  w.AssertExpectations(t)
  assert.Equal(t, w.Written, h)
}

func TestRegexReplaceAllString(t *testing.T) {
  frontend := New(pathToTestTemplateHtml)
  frontend.RegexReplaceAllString("<html>", "<booty>")
  frontend.RegexReplaceAllString("</html>", "</booty>")
  w := new(mockIoWriter)
  h, _ := readTemplateHtml()
  h = bytesReplaceAllString(h, "<html>", "<booty>")
  h = bytesReplaceAllString(h, "</html>", "</booty>")
  w.On("Write", h).Return(len(h),nil)
  frontend.Render(w, []byte(""))
  w.AssertExpectations(t)
  assert.Equal(t, w.Written, h)
}

// TODO: Test malformed regex replace all returns error

//
// Components (to support Testing)
//

const pathToTestTemplateHtml = "./frontend_test.html"

var testApp = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("app\n"))
})

func readTemplateHtml() (t []byte, err error) {
  t, err = ioutil.ReadFile(pathToTestTemplateHtml)
  return
}

func panicNewFrontendWithIncorrectPath() {
  New("./frontend_test_incorrect_path.html")
}

type mockIoWriter struct{
  mock.Mock
  Written []byte
}

func (m *mockIoWriter) Write(p []byte) (n int, err error) {
  m.Called(p)
  m.Written = p
  n = len(p)
  return
}

func bytesReplaceAllString(i []byte, exp string, repl string) []byte {
  r, _ := regexp.Compile(exp)
  t := r.ReplaceAllString(string(i), repl)
  return []byte(t)
}
