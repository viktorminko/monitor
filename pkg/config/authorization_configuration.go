package config

// AuthorizationConfiguration contains data for authorization
type AuthorizationConfiguration struct {
	AppID                        string
	AppSecret                    string
	AuthorizationURL             string
	GetAuthorizationTokenTimeout Duration
}
