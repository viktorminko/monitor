package authorization

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
	"net/url"
	chttp "github.com/viktorminko/monitor/http"
)

type HttpAuthorizer struct {
	URL       string
	Timeout   int
	AppID     string
	AppSecret string
	Client    *chttp.Client
}

func (a *HttpAuthorizer) GetToken() (*Token, error) {

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
		"Cache-Control" : {"no-cache"},
		"Content-Type": {"application/x-www-form-urlencoded"},

	}

	log.Println("Sending authorization token request: POST ", a.URL)
	res, err := a.Client.Call(req, time.Duration(a.Timeout)*time.Second)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		err := errors.New("Invalid http response status: " + http.StatusText(res.StatusCode))
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
