package config

import (
	"os"
	"github.com/viktorminko/monitor/helper"
)

type AuthorizationConfiguration struct {
	AppID                        string
	AppSecret                    string
	AuthorizationURL             string
	GetAuthorizationTokenTimeout int
}

func (a *AuthorizationConfiguration) InitFromFile (filePath string) (error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}

	return helper.InitObjectFromJsonReader(file, a)
}
