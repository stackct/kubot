package api

import (
	"bytes"
	"fmt"
	"kubot/command"
	"kubot/config"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func init() {
	command.SlackCommandPrefix = "!"
}

func TestExecuteWithInvalidCommand(t *testing.T) {
	rr := Request{T: t,
		Handler: Execute,
		Body:    "foo",
	}.Post()
	AssertResponseCode(t, rr, http.StatusInternalServerError)
}

func TestExecuteWithValidCommand(t *testing.T) {
	config.Conf = config.NewMockConfig()
	rr := Request{T: t,
		Handler: Execute,
		Body:    "!help",
	}.Post()
	AssertResponseCode(t, rr, http.StatusOK)
	assert.Equal(t, "[\"available commands: [cmd]\"]\n", rr.Body.String())
}

// Request struct
type Request struct {
	T           *testing.T
	Handler     http.HandlerFunc
	PathParams  map[string]string
	QueryParams map[string]string
	Headers     map[string]string
	Cookies     map[string]string
	Body        string
}

func (r Request) Post() *httptest.ResponseRecorder {
	return r.send("POST")
}

func (r Request) send(method string) *httptest.ResponseRecorder {
	path := "/mock"

	if nil != r.QueryParams {
		params := make(url.Values)
		for k, v := range r.QueryParams {
			params.Add(k, v)
		}
		path += fmt.Sprintf("%s?%s", path, params.Encode())
	}

	body := new(bytes.Buffer)
	if "" != r.Body {
		body = bytes.NewBuffer([]byte(r.Body))
	}

	req, err := http.NewRequest(method, path, body)
	if nil != err {
		r.T.Fatal(err)
	}

	if nil != r.PathParams {
		req = mux.SetURLVars(req, r.PathParams)
	}

	if nil != r.Headers {
		for k, v := range r.Headers {
			req.Header.Add(k, v)
		}
	}

	if nil != r.Cookies {
		for k, v := range r.Cookies {
			req.AddCookie(&http.Cookie{Name: k, Value: v})
		}
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(r.Handler)
	handler.ServeHTTP(rr, req)

	return rr
}

func AssertResponseCode(t *testing.T, rr *httptest.ResponseRecorder, expectedStatusCode int) {
	if status := rr.Code; status != expectedStatusCode {
		t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatusCode)
	}
}
