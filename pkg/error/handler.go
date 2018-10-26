package error

import (
	"github.com/viktorminko/monitor/pkg/notifiers"
	"log"
)

// Handler starts error handler goroutine
type Handler struct {
}

// Run starts handler goroutine and passes notifiers to error reporter
func (eh *Handler) Run(senders *notifiers.Senders) chan<- error {

	log.Println("Error handler started")

	c := make(chan error)

	go func() {
		defer close(c)
		for v := range c {
			err := Report(v, senders)
			if err != nil {
				Check(err)
			}

			Check(v)
		}
	}()

	return c
}
