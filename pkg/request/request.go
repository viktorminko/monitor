package request

import (
	"github.com/viktorminko/monitor/pkg/authorization"
	"github.com/viktorminko/monitor/pkg/config"
	"time"
)

type Request struct {
	Domain         string
	Definition     *config.Definition
	LastExecutedAt time.Time
}

type ExecutionData struct {
	Test         *Request
	ResponseTime time.Duration
	Err          error
}

//Check if current request should be executed at provided time
//Uses Definition RunPeriod
func (t *Request) IsNeedToRun(at time.Time) bool {
	return t.LastExecutedAt.IsZero() || at.Sub(t.LastExecutedAt) > t.Definition.RunPeriod.Duration
}

func (t *Request) Run(token *authorization.Token, statisticsChannel chan<- ExecutionData, caller *Runner) bool {
	timeStart := time.Now()

	err := caller.RunTest(t.Definition, t.Domain, token)

	t.LastExecutedAt = time.Now()

	statisticsChannel <- ExecutionData{
		t,
		time.Since(timeStart),
		err,
	}

	return true
}

func Prepare(definitions []config.Definition, environment *config.Context, domain string) ([]Request, error) {
	var tests []Request
	for i := range definitions {
		err := definitions[i].Prepare(environment)
		if err != nil {
			return nil, err
		}
		tests = append(tests, Request{
			domain,
			&definitions[i],
			time.Time{},
		})
	}

	return tests, nil
}
