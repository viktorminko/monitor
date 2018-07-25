package request

import (
	"github.com/viktorminko/monitor/authorization"
	"time"
)

type Suite struct {
	Requests         []Request
	StatisticChannel chan<- ExecutionData
	ErrorChannel     chan<- error
}

func (t *Suite) Run(token *authorization.Token, caller *Runner) bool {
	isAllTestsPassed := true

	ch := make(chan bool)

	counter := 0

	for i := range t.Requests {
		if t.Requests[i].IsNeedToRun(time.Now()) {
			counter++
			go func(test *Request, c chan<- bool) {
				c <- test.Run(token, t.StatisticChannel, caller)
			}(&t.Requests[i], ch)
		}
	}

	for i := 0; i < counter; i++ {
		if !<-ch {
			isAllTestsPassed = false
		}
	}

	return isAllTestsPassed
}
