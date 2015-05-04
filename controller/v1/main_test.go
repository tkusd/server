package v1

import (
	"io"

	"github.com/tommy351/app-studio-server/model/types"

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

var fixtureElements = []struct {
	Name       string            `json:"name"`
	Type       types.ElementType `json:"type"`
	Attributes types.JSONObject  `json:"attributes"`
}{
	{
		Name: "Text",
		Type: types.ElementTypeScreen,
		Attributes: map[string]interface{}{
			"foo": "bar",
		},
	},
}

func init() {
	router = Router()
}

func request(options *requestOptions) *httptest.ResponseRecorder {
	var body io.Reader

	if options.Body != nil {
		if b, err := json.Marshal(options.Body); err != nil {
			log.Fatal(err)
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
