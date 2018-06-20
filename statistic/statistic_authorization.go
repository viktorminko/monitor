package statistic

import (
	"net"
	"github.com/viktorminko/monitor/authorization"
	)

type Authorization struct {
	Statistic
	Errors           []*authorization.RequestData
	AmountOfTimeouts int
}

func (a *Authorization) Update(newData authorization.RequestData) error {
	a.UpdateExecutionData(newData.ResponseTime, newData.Err)

	if newData.Err != nil {
		a.Errors = append(a.Errors, &newData)

		if err, ok := newData.Err.(net.Error); ok && err.Timeout() {
			a.AmountOfTimeouts++
		}
	}

	return nil
}

func (a *Authorization) Reset() error {

	a.Statistic = Statistic{0, 0, 0}
	a.AmountOfTimeouts = 0
	a.Errors = nil

	return nil
}
