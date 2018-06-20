package test

import (
	"time"
	"github.com/viktorminko/monitor/config"
	"github.com/viktorminko/monitor/authorization"
	"github.com/viktorminko/monitor/method"

)

type Test struct {
	Domain         string
	APIMethod      *method.Data
	LastExecutedAt time.Time
}

type ExecutionData struct {
	Test         *Test
	ResponseTime time.Duration
	Err          error
}

//Check if current test should be executed at provided time
//Uses APIMethod RunPeriod
func (t *Test) IsNeedToRun(at time.Time) bool {
	return t.LastExecutedAt.IsZero() || at.Sub(t.LastExecutedAt).Seconds() > float64(t.APIMethod.RunPeriod)
}

func (t *Test) Run(token *authorization.Token, statisticsChannel chan<- ExecutionData, caller *APICaller) bool {
	timeStart := time.Now()

	err := caller.RunApiMethod(t.APIMethod, t.Domain, token)

	t.LastExecutedAt = time.Now()

	statisticsChannel <- ExecutionData{
		t,
		time.Since(timeStart),
		err,
	}

	return true
}

func Prepare(apiMethods []method.Data, environment *config.Environment, domain string) ([]Test, error) {
	var tests []Test
	for i := range apiMethods {
		err := apiMethods[i].Prepare(environment)
		if err != nil {
			return nil, err
		} else {
			tests = append(tests, Test{
				domain,
				&apiMethods[i],
				time.Time{},
			})
		}
	}

	return tests, nil
}
