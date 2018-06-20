package test

import (
	"time"
	"github.com/viktorminko/monitor/authorization"
)

type Suite struct {
	Tests            []Test
	StatisticChannel chan<- ExecutionData
	ErrorChannel     chan<- error
}

func (t *Suite) Run(token *authorization.Token, caller *APICaller) bool {
	isAllTestsPassed := true

	ch := make(chan bool)

	counter := 0

	for i := range t.Tests {
		if t.Tests[i].IsNeedToRun(time.Now()) {
			counter++
			go func(test *Test, c chan<- bool) {
				c <- test.Run(token, t.StatisticChannel, caller)
			}(&t.Tests[i], ch)
		}
	}

	for i := 0; i < counter; i++ {
		if !<-ch {
			isAllTestsPassed = false
		}
	}

	return isAllTestsPassed
}
