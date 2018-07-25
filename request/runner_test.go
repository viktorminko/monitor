package request

import (
	"fmt"
	"github.com/viktorminko/monitor/authorization"
	"github.com/viktorminko/monitor/config"
	cerror "github.com/viktorminko/monitor/error"
	chttp "github.com/viktorminko/monitor/http"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type TestData struct {
	Name       string
	Definition *config.Definition
	Token      *authorization.Token
	handler    func(w http.ResponseWriter, r *http.Request)
	tester     func(err error, t *testing.T)
}

var testData = []*TestData{
	//Request Timeout
	{
		"Timeout",
		&config.Definition{
			TimeOut:      config.Duration{1},
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(time.Duration(2) * time.Second)
			fmt.Fprintln(w, "Response from request server")
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
		&config.Definition{
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
		&config.Definition{
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
		&config.Definition{
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
		&config.Definition{
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
		&config.Definition{
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
		&config.Definition{
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
		&config.Definition{
			Sample:       false,
			URL:          "/request/1/post?q=1&j=2",
			ResponseCode: http.StatusOK,
		},
		nil,
		func(w http.ResponseWriter, r *http.Request) {
			if "/request/1/post?q=1&j=2" == r.URL.RequestURI() {
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
		&config.Definition{
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
		&config.Definition{
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

	for _, test := range testData {
		server := httptest.NewServer(http.HandlerFunc(test.handler))

		t.Run(test.Name, func(t *testing.T) {
			test.tester(
				(&Runner{
					errChan,
					&chttp.Client{},
				}).RunTest(
					test.Definition,
					server.URL,
					test.Token,
				),
				t,
			)
		})
		server.Close()
	}
}
