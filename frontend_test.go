// Package facade implements a middleware frontending solution.
package facade

import (
  "github.com/stretchr/testify/assert"
  "net/http"
  "testing"
  "fmt"
  "github.com/stretchr/testify/mock"
  "io/ioutil"
  "regexp"
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
  frontend.Render(w, Casing{})
  w.AssertExpectations(t)
  assert.Equal(t, w.Written, h)
}

func TestPreReplaceAll(t *testing.T) {
  frontend := New(pathToTestTemplateHtml)
  frontend.PreReplaceAll("<html>", "<booty>")
  frontend.PreReplaceAll("</html>", "</booty>")
  w := new(mockIoWriter)
  h, _ := readTemplateHtml()
  h = bytesReplaceAllString(h, "<html>", "<booty>")
  h = bytesReplaceAllString(h, "</html>", "</booty>")
  w.On("Write", h).Return(len(h),nil)
  frontend.Render(w, Casing{})
  w.AssertExpectations(t)
  assert.Equal(t, w.Written, h)
}

// TODO: Test malformed regex replace all returns error

func TestWillReplaceAll(t *testing.T) {
  title := "best seats in the house"
  title_exp := "<title>([^<]*)</title>"
  title_repl := "<title>%s</title>"
  content := "ipsum dolor dolor bill yall"
  content_exp := "<facade/>"
  content_repl := "%s"
  frontend := New(pathToTestTemplateHtml)
  frontend.WillReplaceAll( "title", title_exp, title_repl)
  frontend.WillReplaceAll( "content", content_exp, content_repl)
  w := new(mockIoWriter)
  h, _ := readTemplateHtml()
  h = bytesReplaceAllString(h, title_exp, fmt.Sprintf(title_repl, title))
  h = bytesReplaceAllString(h, content_exp, fmt.Sprintf(content_repl, content))
  w.On("Write", h).Return(len(h),nil)
  frontend.Render(w, Casing{
    "title":title,
    "content":content,
    })
  w.AssertExpectations(t)
  assert.Equal(t, w.Written, h)
}

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
