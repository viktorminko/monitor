package notifiers

import "testing"

func TestEmailMessage_GetRFCMessageString(t *testing.T) {
	t.Run("valid string", func(t *testing.T) {
		msg := &Message{
			Type:    "my_type",
			Subject: "my_subj",
			Body:    "my_body",
		}

		from := "from@domain.com"
		res, err := msg.GetRFCMessageString(from)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "MIME-version: 1.0;\nContent-Type: " +
			msg.Type +
			"; charset=\"UTF-8\";\r\n" +
			"From: " + from + "\r\n" +
			"Subject: " + msg.Subject + "\r\n" +
			"\n" +
			msg.Body +
			"\n"
		if res != expected {
			t.Fatalf("unexpected result, expected %v, got %v", expected, res)
		}
	})

	t.Run("default message type", func(t *testing.T) {
		msg := &Message{}

		from := "from@domain.com"
		res, err := msg.GetRFCMessageString(from)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\r\n" +
			"From: " + from + "\r\n" +
			"Subject: " + "\r\n" +
			"\n" +
			"\n"
		if res != expected {
			t.Fatalf("unexpected result, expected %v, got %v", expected, res)
		}
	})
}
