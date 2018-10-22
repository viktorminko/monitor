package config

import "time"

// AuthorizationConfiguration contains data for authorization
type AuthorizationConfiguration struct {
	AppID                        string
	AppSecret                    string
	AuthorizationURL             string
	GetAuthorizationTokenTimeout time.Duration
}
