package statistic

import (
	"errors"
	"testing"
	"time"
)

func checkStats(t *testing.T, s Statistic, err error, ex int, er int, re time.Duration) {

	if err != nil {
		t.Errorf("Unexpected error return %s", err)
	}

	if s.AmountOfExecutions != ex {
		t.Errorf(
			"invalid amount of executions, expected: %v, got: %v",
			ex,
			s.AmountOfExecutions,
		)
	}

	if s.AmountOfErrors != er {
		t.Errorf(
			"invalid number of errors, expected: %v, got: %v",
			er,
			s.AmountOfErrors,
		)
	}

	if s.AverageResponseTime != re {
		t.Errorf(
			"invalid AverageResponseTime, expected: %v, got: %v",
			re,
			s.AverageResponseTime,
		)
	}
}

func TestStatistic_UpdateExecutionData(t *testing.T) {
	s := Statistic{}

	testData := []struct {
		D   string
		Err error
		Ex  int
		Er  int
		Ere string
	}{
		{"2s", nil, 1, 0, "2s"},
		{"4s", nil, 2, 0, "3s"},
		{"6s", nil, 3, 0, "4s"},
		{"8s", nil, 4, 0, "5s"},
		{"5s", nil, 5, 0, "5s"},
		{"2s", nil, 6, 0, "4.5s"},
		{"8s", errors.New(""), 7, 1, "5s"},
		{"5s", errors.New(""), 8, 2, "5s"},
	}

	for _, test := range testData {
		d, _ := time.ParseDuration(test.D)
		expAverage, _ := time.ParseDuration(test.Ere)
		err := s.UpdateExecutionData(d, test.Err)
		checkStats(t, s, err, test.Ex, test.Er, expAverage)
	}
}
