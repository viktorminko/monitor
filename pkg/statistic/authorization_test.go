package statistic

import (
	"github.com/viktorminko/monitor/pkg/authorization"
	"testing"
	"time"
)

func checkAuth(t *testing.T, auth Authorization, err error, ex int, er int, re time.Duration) {

	if err != nil {
		t.Errorf("Unexpected error return %s", err)
	}

	if auth.AmountOfExecutions != ex {
		t.Errorf(
			"invalid number of executions, expected: %v, got: %v",
			ex,
			auth.AmountOfExecutions,
		)
	}

	if auth.AmountOfErrors != er {
		t.Errorf(
			"invalid number of errors, expected: %v, got: %v",
			er,
			auth.AmountOfErrors,
		)
	}

	if auth.AverageResponseTime != re {
		t.Errorf(
			"invalid AverageResponseTime, expected: %v, got: %v",
			re,
			auth.AverageResponseTime,
		)
	}
}

func TestAuthorization_Update(t *testing.T) {

	auth := Authorization{
		Statistic{
			0,
			0,
			0,
		},
		[]*authorization.RequestData{},
		0,
	}

	err := auth.Update(
		authorization.RequestData{
			time.Now(),
			nil,
			time.Duration(2),
		})

	checkAuth(t, auth, err, 1, 0, time.Duration(2))
}

func TestAuthorization_Reset(t *testing.T) {
	auth := Authorization{
		Statistic{
			1,

			1,
			1,
		},
		[]*authorization.RequestData{},
		0,
	}

	err := auth.Reset()

	checkAuth(t, auth, err, 0, 0, time.Duration(0))
}
