package statistic

import (
  "github.com/viktorminko/monitor/request"
)

type ForTest struct {
	Statistic
}

type Suite struct {
	Tests map[*request.Request]*ForTest
}

func (s *Suite) Update(newData request.ExecutionData) error {

	if s.Tests == nil {
		s.Tests = make(map[*request.Request]*ForTest)
	}

	if _, ok := s.Tests[newData.Test]; !ok {
		s.Tests[newData.Test] = &ForTest{
			Statistic{0, 0, 0},
		}
	}

	s.Tests[newData.Test].UpdateExecutionData(newData.ResponseTime, newData.Err)

	return nil
}

func (s *Suite) Reset() error {
	s.Tests = nil

	return nil
}
