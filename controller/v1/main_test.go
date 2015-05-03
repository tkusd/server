package v1

import (
	"io"

	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
)

type requestOptions struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
}

var router http.Handler

var fixtureUsers = []struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}{
	{Name: "John", Email: "john@abc.com", Password: "123456"},
	{Name: "Mary", Email: "mary@abc.com", Password: "234567"},
}

var fixtureProjects = []struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	IsPrivate   bool   `json:"is_private"`
}{
	{Title: "Hello", Description: "Test"},
	{Title: "World", IsPrivate: true},
}

func init() {
	router = Router()
}

func request(options *requestOptions) *httptest.ResponseRecorder {
	var body io.Reader

	if options.Body != nil {
		if b, err := json.Marshal(options.Body); err != nil {
			panic(err)
		} else {
			body = bytes.NewReader(b)
		}
	}

	req, err := http.NewRequest(options.Method, options.URL, body)
	w := httptest.NewRecorder()

	req.Header.Set("Content-Type", "application/json")

	if options.Headers != nil {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}

	router.ServeHTTP(w, req)

	if err != nil {
		log.Fatal(err)
	}

	return w
}

func parseJSON(body *bytes.Buffer, data interface{}) error {
	return json.NewDecoder(body).Decode(data)
}
