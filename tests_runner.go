package main

import (
	"errors"
	"github.com/viktorminko/monitor/authorization"
	cerror "github.com/viktorminko/monitor/error"
	"github.com/viktorminko/monitor/request"
	"log"
	"time"
)

// TestsRunner runs test suite based on execution period
// and sends statistics data to provided channels
type TestsRunner struct {
	Suite             *request.Suite
	ExecutionsPeriod  time.Duration
	TestsStatsChannel chan<- request.ExecutionData
	ErrorChannel      chan<- error
}

// Run executes test suite periodically based on provided execution period
func (t *TestsRunner) Run(authHandler *authorization.Handler, caller *request.Runner) {
	ticker := time.NewTicker(t.ExecutionsPeriod)
	for {

		log.Println("Running request suite")

		var token *authorization.Token
		var authError error
		if authHandler != nil {

			timeStart := time.Now()

			token, authError = authHandler.Authorizer.GetToken()

			authHandler.StatsChan <- authorization.RequestData{
				time.Now(),
				authError,
				time.Since(timeStart),
			}

			if authError != nil {
				t.ErrorChannel <- cerror.NonFatal{"error occurred while retrieving authorization token. Current request round aborted", authError}
			}
		}

		if authError == nil {
			isAllTestsPassed := t.Suite.Run(token, caller)

			if !isAllTestsPassed {
				t.ErrorChannel <- cerror.NonFatal{"error occurred while running tests suite. One or more API request(s) returned unexpected result.", errors.New("")}
			}
		}

		<-ticker.C
	}
}
