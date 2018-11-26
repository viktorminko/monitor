package statistic

import (
	cerror "github.com/viktorminko/monitor/pkg/error"
	"github.com/viktorminko/monitor/pkg/notifiers"
	"log"
	"time"
)

type Reporter struct {
	ErrorChannel chan<- error
}

func (s *Reporter) Run(
	ExecutionsPeriod time.Duration,
	statsReceiver <-chan *Monitor,
	senders *notifiers.Senders) {

	log.Println("Statistics reporter started")

	go func() {

		t := time.NewTicker(ExecutionsPeriod)
		for {
			<-t.C
			stats := <-statsReceiver
			err := s.SendReport(stats, senders)

			//if we didn't send report, then don't reset statistics
			if err != nil {
				s.ErrorChannel <- cerror.NonFatal{"error occurred while sending statistics report", err}
				continue
			}

			log.Println("resetting statistics")
			stats.Reset()

		}
	}()

	return
}

func (s *Reporter) SendReport(stats *Monitor, senders *notifiers.Senders) error {
	log.Println("sending statistics")

	return senders.SendToAll(
		"statistic_report",
		map[string]interface{}{
			"Date":      time.Now().Local().Format("Mon Jan 2 15:04:05 2006"),
			"Tests":     stats.Suite.Tests,
			"AuthStats": stats.Authorization,
		})
}
