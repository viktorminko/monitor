package error

import (
	"log"
	"github.com/viktorminko/monitor/notifiers"
)

type ErrorHandler struct {
}

func (eh *ErrorHandler) Run(senders *notifiers.Senders) chan<- error {

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
