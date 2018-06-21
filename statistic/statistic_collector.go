package statistic

import (
	"log"
	"github.com/viktorminko/monitor/request"
	"github.com/viktorminko/monitor/authorization"

)


type Collector struct {
	Statistics *Monitor
}

func (s *Collector) Run() (
	chan<- request.ExecutionData,
	chan<- authorization.RequestData,
	<-chan *Monitor) {

	log.Println("Statistics collector started")

	c1 := make(chan request.ExecutionData)
	c2 := make(chan authorization.RequestData)
	c3 := make(chan *Monitor)
	go func() {
		for {
			select {
			case s1 := <-c1:
				s.Statistics.Suite.Update(s1)
			case s2 := <-c2:
				s.Statistics.Authorization.Update(s2)
			case c3 <- s.Statistics:
				//Send stats to some handler
			}
		}
	}()

	return c1, c2, c3
}
