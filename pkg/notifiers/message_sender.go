package notifiers

import (
	"errors"
	"log"
	"sync"
)

type Sender interface {
	SendMessage(mID string, mBody map[string]interface{}) error
}

type Senders []Sender

//if sending to all senders failed then we return error
func (m *Senders) SendToAll(mID string, mBody map[string]interface{}) error {

	isAllFailed := true

	var wg sync.WaitGroup
	for i := range *m {
		s := (*m)[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.SendMessage(mID, mBody)
			if err != nil {
				log.Printf("error while reporting message: %v", err)
				return
			}
			isAllFailed = false
		}()
	}

	wg.Wait()

	if isAllFailed {
		return errors.New("all senders failed to send messages")
	}

	return nil
}

func (m *Senders) Add(s Sender) {
	*m = append(*m, s)
}
