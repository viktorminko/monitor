package authorization

import "time"

//Authorization token used in API requests
type Token struct {
	Token     string `json:"access_token"`
	Expires   int    `json:"expires_in"`
	Scope     string `json:"scope"`
	TokenType string `json:"token_type"`
}

type RequestData struct {
	Time         time.Time
	Err          error
	ResponseTime time.Duration
}

type Handler struct {
	Authorizer Authorizer
	StatsChan  chan<- RequestData
}

type Authorizer interface {
	GetToken() (*Token, error)
}
