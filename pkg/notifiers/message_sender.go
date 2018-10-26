package notifiers

import "log"

type Sender interface {
	SendMessage(mID string, mBody map[string]interface{}) error
}

type Senders []Sender

func (m *Senders) SendToAll(mID string, mBody map[string]interface{}) {
	for i := range *m {
		s := (*m)[i]
		go func() {
			err := s.SendMessage(mID, mBody)
			if err != nil {
				log.Printf("error while reporting message: %v", err)
			}
		}()
	}
}

func (m *Senders) Add(s Sender) {
	*m = append(*m, s)
}
