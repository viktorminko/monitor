package notifiers

import ()
import (
	"os/exec"
	"log"
	"strings"
	"fmt"
	"io/ioutil"
)

type SendmailEmailSender struct {
	Account *EmailAccount
}

func (e *SendmailEmailSender) SendMessage(message *Message) error {
	log.Println("Sending message: " + message.Subject + " " + message.Body)
	recipients, err := e.Account.GetRecipients()
	if err != nil {
		return err
	}

	recipientsString := strings.Join(recipients, ", ")

	msg := "From: " + e.Account.From + "\n"
	msg += "To: " + recipientsString + "\n"
	msg += "Subject: " + message.Subject + "\n\n"
	msg += message.Body + "\n"

	sendmail := exec.Command("/usr/sbin/sendmail", "-t")
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := sendmail.StdoutPipe()
	if err != nil {
		panic(err)
	}

	sendmail.Start()
	stdin.Write([]byte(msg))
	stdin.Close()
	sentBytes, _ := ioutil.ReadAll(stdout)
	sendmail.Wait()

	fmt.Println("Send Command Output")
	fmt.Println(string(sentBytes))

	return nil
}
