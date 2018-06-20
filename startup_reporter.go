package main

import (
	"encoding/json"
	"log"
	"time"
	"github.com/viktorminko/monitor/config"
	"github.com/viktorminko/monitor/notifiers"
	"github.com/viktorminko/monitor/helper"
	"github.com/viktorminko/monitor/method"
	cerror "github.com/viktorminko/monitor/error"
)

type StartupReporter struct {
	ErrorChannel chan<- error
}

func (s *StartupReporter) Send(
	config *config.Configuration,
	auth *config.AuthorizationConfiguration,
	APIMethods []method.Data,
	senders *notifiers.Senders) {

	log.Println("Startup reporter started")

	go func() {

		testsJSON, err := json.Marshal(APIMethods)
		if err != nil {
			s.ErrorChannel <- cerror.NonFatal{"error occurred decoding tests data", err}
		}

		senders.SendToAll("startup_report", map[string]interface{}{
			"startup_date": time.Now().Local().Format("Mon Jan 2 15:04:05 2006"),
			"api_url":      config.Domain,
			"auth_url":     auth.AuthorizationURL,
			"auth_timeout": auth.GetAuthorizationTokenTimeout,
			"app_id":       auth.AppID,
			"exec_period":  config.RunPeriod,
			"stats_period": config.StatisticRunPeriod,
			"tests":        string(helper.FormatJSON(testsJSON)),
			"proxy":        config.Proxy,
		})

	}()

	return
}
