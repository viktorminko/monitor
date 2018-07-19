package config

import (
	"github.com/viktorminko/monitor/helper"
)

// AuthorizationConfiguration contains data for authorization
type AuthorizationConfiguration struct {
	AppID                        string
	AppSecret                    string
	AuthorizationURL             string
	GetAuthorizationTokenTimeout int
}

// InitFromFile inits object from JSON file
func (a *AuthorizationConfiguration) InitFromFile (filePath string) (error) {
	return helper.InitObjectFromJsonFile(filePath, a)
}
