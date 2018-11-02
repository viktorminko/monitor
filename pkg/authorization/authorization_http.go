package authorization

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	chttp "github.com/viktorminko/monitor/pkg/http"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// HTTPAuthorizer gets authorization token via HTTP request to URL
type HTTPAuthorizer struct {
	URL       string
	Timeout   time.Duration
	AppID     string
	AppSecret string
	Client    *chttp.Client
}

// GetToken retrieves authorization token from URL via HTTP request
func (a *HTTPAuthorizer) GetToken() (*Token, error) {

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(
		"POST",
		a.URL,
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte(a.AppID+":"+a.AppSecret))},
		"Cache-Control": {"no-cache"},
		"Content-Type":  {"application/x-www-form-urlencoded"},
	}

	log.Println("sending authorization token request: POST ", a.URL)
	res, err := a.Client.Call(req, a.Timeout)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := errors.New("invalid http response status: " + http.StatusText(res.StatusCode))
		return nil, err
	}

	body, _ := ioutil.ReadAll(res.Body)

	var token Token
	err = json.Unmarshal([]byte(body), &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
