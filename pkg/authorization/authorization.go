package authorization

import "time"

// Token is an authorization token used in HTTP requests
type Token struct {
	Token     string `json:"access_token"`
	Expires   int    `json:"expires_in"`
	Scope     string `json:"scope"`
	TokenType string `json:"token_type"`
}

// RequestData contains information about request execution
type RequestData struct {
	Time         time.Time
	Err          error
	ResponseTime time.Duration
}

// Handler gets authorization tokens and sends statistics
type Handler struct {
	Authorizer Authorizer
	StatsChan  chan<- RequestData
}

// Authorizer interface to get authorization tokens
type Authorizer interface {
	GetToken() (*Token, error)
}
