package statistic

import (
	"github.com/viktorminko/monitor/pkg/notifiers"
	"github.com/viktorminko/monitor/pkg/request"
	"reflect"
	"sync"
	"testing"
	"time"
)

type TestSender struct {
	TestStats map[*request.Request]*Statistic
	AuthStats Authorization
	wg        *sync.WaitGroup
}

func (s *TestSender) SendMessage(mID string, mBody map[string]interface{}) error {
	s.TestStats = mBody["Tests"].(map[*request.Request]*Statistic)
	s.AuthStats = mBody["AuthStats"].(Authorization)

	s.wg.Done()

	return nil
}

//1. Init custom sender and wait until its executed with WaitGroup
//2. Run statistics reporter in goroutine
//3. Send statistics message to the channel for reporter
//4. Check if sender received valid statistics data to send
func TestReporter_StatisticIsSent(t *testing.T) {
	errch := make(chan<- error)
	reporter := Reporter{errch}

	statsReceiver := make(chan *Monitor)

	var wg sync.WaitGroup

	wg.Add(1)
	testSender := &TestSender{nil, Authorization{}, &wg}
	senders := notifiers.Senders{}
	senders = append(senders, testSender)

	statisticReportingInterval := time.Millisecond * 100

	go func() {
		reporter.Run(
			statisticReportingInterval,
			statsReceiver,
			&senders,
		)
	}()

	testStats := map[*request.Request]*Statistic{
		&request.Request{}: &Statistic{
			3,
			time.Millisecond * 2,
			1,
		},
	}

	testStatsSuite := Suite{testStats}

	authStats := Authorization{
		Statistic{
			4,
			time.Millisecond * 3,
			2,
		},
		nil,
		1,
	}
	go func() {
		statsReceiver <- &Monitor{testStatsSuite, authStats}
	}()

	//wait until sender receives statistics
	wg.Wait()

	if !reflect.DeepEqual(testSender.TestStats, testStats) {
		t.Fatalf(
			"unexpected tests statistics, expected: %v, got %v",
			testStats,
			testSender.TestStats,
		)
	}

	if !reflect.DeepEqual(testSender.AuthStats, authStats) {
		t.Fatalf(
			"unexpected authorization statistics, expected: %v, got %v",
			authStats,
			testSender.AuthStats,
		)
	}
}
