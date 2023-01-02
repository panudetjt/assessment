package util

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

type HttpResponse struct {
	*http.Response
	Error error
}

func (r *HttpResponse) Decode(v interface{}) error {
	if r.Error != nil {
		return r.Error
	}

	return json.NewDecoder(r.Body).Decode(v)
}

type Response struct {
	Context  echo.Context
	Recorder *httptest.ResponseRecorder
}

func (r *Response) Decode(v any) error {
	return json.NewDecoder(r.Recorder.Body).Decode(v)
}

func Uri(paths ...string) string {
	var host string
	if h := os.Getenv("HOST"); h != "" {
		host = h
	} else {
		host = "http://localhost:2565"
	}

	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func Request(method, url string, body io.Reader) *HttpResponse {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &HttpResponse{res, err}
}

func RequestE(method, url string, body io.Reader) *Response {
	req := httptest.NewRequest(method, url, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rr := httptest.NewRecorder()
	c := echo.New().NewContext(req, rr)
	return &Response{Context: c, Recorder: rr}
}
