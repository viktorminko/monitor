package notifiers

import (
	"log"
	"net/smtp"
	"strconv"
)

type SmtpEmailSender struct {
	Account *EmailAccount
	WorkDir string
}

func (e *SmtpEmailSender) SendMessage(mID string, mBody map[string]interface{}) error {

	//Init message
	message, err := BuildMessage(e.WorkDir, mID, mBody)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("",
		e.Account.Username,
		e.Account.Password,
		e.Account.EmailServer,
	)

	recepients, err := e.Account.GetRecipients()
	if err != nil {
		return err
	}

	from := e.Account.From

	messageString, err := message.GetRFCMessageString(from)
	if err != nil {
		return err
	}

	if e.Account.IsDebugMode {
		log.Println("Email debug mode: " + message.Subject + " " + message.Body)
		return nil
	}

	err = smtp.SendMail(
		e.Account.EmailServer+":"+strconv.Itoa(e.Account.Port),
		auth,
		from,
		recepients,
		[]byte(messageString))
	if err != nil {
		return err
	}

	return nil
}