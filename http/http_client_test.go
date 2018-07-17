package http

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"log"
	"io/ioutil"
	"os"
)

func TestHTTPClient_Call(t *testing.T) {

	t.Run("invalid proxy", func(t *testing.T) {
		proxyURL, _ := url.Parse("invalid proxy")

		client := &Client{
			proxyURL,
		}

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))

		req, _ := http.NewRequest("GET", server.URL, nil)
		_, err := client.Call(
			req,
			0,
		)

		if err == nil {
			t.Fatal("error expected but not returned")
		}
	})

	//Test that request to server is passed via proxy and back.
	//The server will return one status code to proxy and proxy will change it and return to client.
	//So we are sure that request is handled via proxy
	t.Run("valid proxy request", func(t *testing.T) {
		codeServerProxy := http.StatusAlreadyReported
		codeProxyClient := http.StatusAccepted

		//Handler in proxy server
		proxyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//Just forward request to target server in r.URL
			req, _ := http.NewRequest("GET", r.URL.String(), nil)
			res, err := (&http.Client{}).Do(req)

			if err != nil {
				t.Fatalf("unexpected error returned, %v", err)
			}

			//Check that target server returned correct status code
			if res.StatusCode != codeServerProxy {
				t.Fatalf("unexpected status code returned to proxy from server, expected %v, got %v", codeServerProxy, res.StatusCode)
			}

			//Change response status code to client
			w.WriteHeader(codeProxyClient)
		}))

		//Handler in target server
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(codeServerProxy)
		}))

		//Get already created proxyServer as proxy
		proxyURL, _ := url.Parse(proxyServer.URL)

		//Define client with proxy
		client := &Client{
			proxyURL,
		}

		//Send request to target server
		req, _ := http.NewRequest("GET", server.URL, nil)
		res, err := client.Call(
			req,
			0,
		)

		//Request should work correctly
		if err != nil {
			t.Fatalf("unexpected error returned, %v", err)
		}

		//Check that response code is the one from proxy server
		if res.StatusCode != codeProxyClient {
			t.Fatalf("unexpected status code returned to client from proxy, expected %v, got %v", codeProxyClient, res.StatusCode)
		}
	})

	//Test that request to server is passed to proxy and keeps same URL
	t.Run("valid request", func(t *testing.T) {
		timeout := time.Duration(2) * time.Second

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			time.Sleep(timeout - time.Duration(1)*time.Second)
		}))

		req, _ := http.NewRequest("GET", server.URL, nil)
		res, err := (&Client{}).Call(
			req,
			timeout,
		)

		//Request should work correctly
		if err != nil {
			t.Fatalf("unexpected error returned, %v", err)
		}

		if res.StatusCode != http.StatusCreated {
			t.Errorf("invalid response status, expected %v, got %v", http.StatusCreated, res.StatusCode)
		}
	})

}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
