package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"

	"github.com/labstack/echo/v4"
)

// TestRouter is for testing echo APIs the easy way
type TestRouter struct {
	Echo *echo.Echo
	req  *http.Request
}

// NewTestRouter returns a pointer to a new TestRouter instance
func NewTestRouter() *TestRouter {
	return &TestRouter{
		Echo: echo.New(),
	}
}

// Request creates a new http.Request and returns it, so you can add everything you want
func (r *TestRouter) Request(method, route string, postData interface{}) *http.Request {
	var data []byte
	contentType := ""
	if postData != nil {
		t := reflect.TypeOf(postData)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		switch t.Kind() {
		case reflect.Struct, reflect.Map:
			data, _ = json.Marshal(postData)
			contentType = "application/json"
		default:
			switch d := postData.(type) {
			case string:
				data = []byte(d)
			case []byte:
				data = d
			default:
				panic("data makes no sense")
			}
		}
	}
	r.req = httptest.NewRequest(method, route, bytes.NewReader(data))
	if contentType != "" {
		r.req.Header.Add("Content-Type", contentType)
	} else {
		r.req.Header.Add("Content-Type", http.DetectContentType(data))
	}
	return r.req
}

// Start starts the router, thus makes the requests
func (r *TestRouter) Start() *TestResponse {
	rec := httptest.NewRecorder()
	r.Echo.ServeHTTP(rec, r.req)
	return newTestResponse(rec)
}

// TestResponse is the response to a test request
type TestResponse struct {
	Code   int
	Header http.Header
	Body   *bytes.Buffer
	rec    *httptest.ResponseRecorder
}

func newTestResponse(rec *httptest.ResponseRecorder) *TestResponse {
	return &TestResponse{
		Code:   rec.Code,
		Header: rec.HeaderMap,
		Body:   rec.Body,
		rec:    rec,
	}
}

// Response returns the http.Response object
func (r *TestResponse) Response() *http.Response {
	return r.rec.Result()
}

// Bind binds the response data to given object
func (r *TestResponse) Bind(obj interface{}) error {
	if r.Body == nil {
		return fmt.Errorf("no data")
	}
	return json.Unmarshal(r.Body.Bytes(), obj)
}

// String is the response to your request as a string
func (r *TestResponse) String() string {
	if r.Body == nil {
		return ""
	}
	return r.Body.String()
}
