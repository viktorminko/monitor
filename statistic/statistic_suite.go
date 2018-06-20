package statistic

import (
  "github.com/viktorminko/monitor/test"
)

type ForTest struct {
	Statistic
}

type Suite struct {
	Tests map[*test.Test]*ForTest
}

func (s *Suite) Update(newData test.ExecutionData) error {

	if s.Tests == nil {
		s.Tests = make(map[*test.Test]*ForTest)
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
