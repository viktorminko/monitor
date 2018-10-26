package statistic

import (
	"github.com/viktorminko/monitor/pkg/request"
)

type Suite struct {
	Tests map[*request.Request]*Statistic
}

func (s *Suite) Update(newData request.ExecutionData) error {

	if s.Tests == nil {
		s.Tests = make(map[*request.Request]*Statistic)
	}

	if _, ok := s.Tests[newData.Test]; !ok {
		s.Tests[newData.Test] = &Statistic{0, 0, 0}

	}

	s.Tests[newData.Test].UpdateExecutionData(newData.ResponseTime, newData.Err)

	return nil
}

func (s *Suite) Reset() error {
	s.Tests = nil

	return nil
}
