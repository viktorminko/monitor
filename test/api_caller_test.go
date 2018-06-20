package test

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/viktorminko/monitor/method"
	"github.com/viktorminko/monitor/authorization"
	cerror "github.com/viktorminko/monitor/error"
	chttp "github.com/viktorminko/monitor/http"
)

type TestData struct {
	Name    string
	Method  *method.Data
	Token   *authorization.Token
	handler func(w http.ResponseWriter, r *http.Request)
	tester  func(err error, t *testing.T)
}

var testData = []*TestData{
	//Test Timeout
	{
		"Timeout",
		&method.Data{
			TimeOut:      1,
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(2) * time.Second)
			fmt.Fprintln(w, "Response from test server")
		},
		func(err error, t *testing.T) {
			if nil == err {
				t.Fatal("Error expected, but not returned")
			}

			if err, ok := err.(net.Error); !ok || !err.Timeout() {
				t.Error("Timeout error expected, but another returned")
			}
		},
	},

	{
		"Invalid path",
		&method.Data{
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
		func(err error, t *testing.T) {
			if err == nil {
				t.Error("Error expected but not returned")
			}
		},
	},

	{
		"Invalid response code",
		&method.Data{
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			//Return any non 200 code
			w.WriteHeader(http.StatusCreated)
		},
		func(err error, t *testing.T) {
			if _, ok := err.(cerror.NonFatal); !ok {
				t.Error("Error expected, but not returned")
			}
		},
	},

	{
		"Don't check body",
		&method.Data{
			Sample:       false,
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, `{"status":"error"}`)
		},
		func(err error, t *testing.T) {
			if nil != err {
				t.Errorf("Unexpected error returned: %s", err.Error())
			}
		},
	},

	{
		"Invalid response body",
		&method.Data{
			Sample:       map[string]interface{}{"status": "ok"},
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `{"status":"fail"}`)
		},
		func(err error, t *testing.T) {
			if err == nil {
				t.Fatal("Error expected, but not returned")
			}

			if _, ok := err.(cerror.Test); !ok {
				t.Errorf("Unexpected error returned: %s", err.Error())
			}
		},
	},

	{
		"Valid response body",
		&method.Data{
			Sample:       map[string]interface{}{"status": "ok"},
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `{"status": "ok"}`)
		},
		func(err error, t *testing.T) {
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			}
		},
	},

	{
		"Valid token provided",
		&method.Data{
			Sample:       false,
			ResponseCode: http.StatusOK,
		},
		&authorization.Token{
			Token: "abcd123",
		},
		func(w http.ResponseWriter, r *http.Request) {
			if "Bearer abcd123" == r.Header.Get("Authorization") {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		},
		func(err error, t *testing.T) {
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			}
		},
	},

	{
		"Valid URL requested",
		&method.Data{
			Sample:       false,
			URL:          "/test/1/post?q=1&j=2",
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			if "/test/1/post?q=1&j=2" == r.URL.RequestURI() {
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		},
		func(err error, t *testing.T) {
			if err != nil {
				t.Errorf("Unexpected error: %s", err.Error())
			}
		},
	},

	{
		"Non JSON server response",
		&method.Data{
			Sample:       map[string]interface{}{"status": "ok"},
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `Some non json response`)
		},
		func(err error, t *testing.T) {
			if err == nil {
				t.Fatal("Error expected, but not returned")
			}
		},
	},

	{
		"Invalid URL in request",
		&method.Data{
			Sample:       false,
			ResponseCode: http.StatusOK,
			URL:          "invalid part of URL",
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {

		},
		func(err error, t *testing.T) {
			if err == nil {
				t.Fatal("Error expected, but not returned")
			}
		},
	},
}

func TestAPICaller_RunApiMethod(t *testing.T) {

	errChan := make(chan error)
	go func() {
		for {
			select {
			case <-errChan:
			}
		}
	}()

	server := httptest.NewServer(http.HandlerFunc(nil))
	defer server.Close()

	for _, test := range testData {
		server.Config.Handler = http.HandlerFunc(test.handler)

		t.Run(test.Name, func(t *testing.T) {
			test.tester(
				(&APICaller{
					errChan,
					&chttp.Client{},
				}).RunApiMethod(
					test.Method,
					server.URL,
					test.Token,
				),
				t,
			)
		})
	}
}
