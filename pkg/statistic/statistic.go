package statistic

import "time"

type Statistic struct {
	AmountOfExecutions  int
	AverageResponseTime time.Duration
	AmountOfErrors      int
}

func (t *Statistic) UpdateExecutionData(newResponseTime time.Duration, err error) error {
	if err != nil {
		t.AmountOfErrors++
	}

	//Average execution is calculated before increasing execution amount
	current := int64(t.AmountOfExecutions) * int64(t.AverageResponseTime)
	t.AverageResponseTime = time.Duration((current + int64(newResponseTime)) / (int64(t.AmountOfExecutions) + 1))

	t.AmountOfExecutions++

	return nil
}
