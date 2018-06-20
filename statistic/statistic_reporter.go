package statistic

import (
	"log"
	"time"
	"github.com/viktorminko/monitor/notifiers"
	cerror "github.com/viktorminko/monitor/error"
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

			if err != nil {
				s.ErrorChannel <- cerror.NonFatal{"error occurred while sending statistics report", err}
			} else {
				//Reset statistics when report is sent
				stats.Suite.Reset()
				stats.Authorization.Reset()
			}
		}
	}()

	return
}

func (s *Reporter) SendReport(stats *Monitor, senders *notifiers.Senders) error {
	log.Println("Sending statistics")

	senders.SendToAll(
		"statistic_report",
		 map[string]interface{}{
		"Date":      time.Now().Local().Format("Mon Jan 2 15:04:05 2006"),
		"Tests":     stats.Suite.Tests,
		"AuthStats": stats.Authorization,
	})

	return nil
}
