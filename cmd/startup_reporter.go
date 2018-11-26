package main

import (
	"encoding/json"
	"github.com/viktorminko/monitor/pkg/config"
	cerror "github.com/viktorminko/monitor/pkg/error"
	"github.com/viktorminko/monitor/pkg/helper"
	"github.com/viktorminko/monitor/pkg/notifiers"
	"log"
	"time"
)

// StartupReporter sends monitor startup report to all specified senders
type StartupReporter struct {
	ErrorChannel chan<- error
}

// Send sends message to all senders with current startup data
func (s *StartupReporter) Send(
	config *config.Configuration,
	auth *config.AuthorizationConfiguration,
	Requests []config.Definition,
	senders *notifiers.Senders) {

	log.Println("Startup reporter started")

	go func() {

		testsJSON, err := json.Marshal(Requests)
		if err != nil {
			s.ErrorChannel <- cerror.NonFatal{"error occurred decoding tests data", err}
		}

		senders.SendToAll("startup_report", map[string]interface{}{
			"startup_date": time.Now().Local().Format("Mon Jan 2 15:04:05 2006"),
			"url":          config.Domain,
			"auth_url":     auth.AuthorizationURL,
			"auth_timeout": auth.GetAuthorizationTokenTimeout,
			"app_id":       auth.AppID,
			"exec_period":  config.RunPeriod,
			"tests":        string(helper.FormatJSON(testsJSON)),
			"proxy":        config.Proxy,
			"stats_period": config.StatisticRunPeriod,
		})

	}()

	return
}
