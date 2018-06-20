package authorization

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	chttp "github.com/viktorminko/monitor/http"
)

func TestHttpAuthorizer_GetToken(t *testing.T) {

	appID := "my_app_id"
	appSecret := "my_secret"
	tokenStr := "my_token"
	timeout := 2

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check authorization data
		passedAppID, passedAppSecret, ok := r.BasicAuth()

		if !ok {
			t.Fatal("basic auth parsing failed")
		}

		if appID != passedAppID {
			t.Errorf("unexpected appID passed, expected %v, got %v", appID, passedAppID)
		}

		if appSecret != passedAppSecret {
			t.Errorf("unexpected appSecret passed, expected %v, got %v", appSecret, passedAppSecret)
		}

		io.WriteString(w, `{"access_token": "`+tokenStr+`"}`)

		//Wait less then expected timeout and there should not be any error
		time.Sleep(time.Duration(timeout-1) * time.Second)

	}))
	defer server.Close()

	token, err := (&HttpAuthorizer{
		server.URL,
		timeout,
		appID,
		appSecret,
		&chttp.Client{},
	}).GetToken()

	if err != nil {
		t.Fatalf("unexpected error returned: %v", err)
	}

	if token.Token != tokenStr {
		t.Errorf("unexpected token returned, expected %v, got %v", tokenStr, token)
	}
}

func TestHttpAuthorizer_GetTokenInvalidResponseCode(t *testing.T) {
	if _, err := (&HttpAuthorizer{
		URL: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusAccepted)

		})).URL,
		Client: &chttp.Client{},
	}).GetToken(); err == nil {
		t.Fatalf("error expected, but not returned")
	}
}
