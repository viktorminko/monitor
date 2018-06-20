package main

import (
	"errors"
	"log"
	"time"
	"github.com/viktorminko/monitor/test"
	"github.com/viktorminko/monitor/authorization"
	cerror "github.com/viktorminko/monitor/error"

)

type TestsRunner struct {
	Suite             *test.Suite
	ExecutionsPeriod  time.Duration
	TestsStatsChannel chan<- test.ExecutionData
	ErrorChannel      chan<- error
}

func (t *TestsRunner) Run(authHandler *authorization.Handler, caller *test.APICaller) {
	ticker := time.NewTicker(t.ExecutionsPeriod)
	for {

		log.Println("Running test suite")

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
				t.ErrorChannel <- cerror.NonFatal{"error occurred while retrieving authorization token. Current test round aborted", authError}
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
