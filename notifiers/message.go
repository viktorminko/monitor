package notifiers

import (
	"text/template"
	"bytes"
	"github.com/viktorminko/monitor/helper"
	"path"
)

type Message struct {
	ID           string
	Type         string
	Name         string
	Subject      string
	BodyTemplate string
	Body         string
}

func (m *Message) GetRFCMessageString() (string, error) {

	messageType := "text/plain"
	if len(m.Type) > 0 {
		messageType = m.Type
	}

	msg := "MIME-version: 1.0;\nContent-Type: " + messageType + "; charset=\"UTF-8\";\r\n"
	msg += "Subject: API monitor: " + m.Subject + "\r\n"
	msg += "\n"
	msg += m.Body + "\n"

	return msg, nil
}

func (m *Message) InsertDataInBody(data map[string]interface{}) error {
	//Prepare message body
	tmpl, err := template.New("tpl").Parse(m.BodyTemplate)
	if err != nil {
		return err
	}
	var t bytes.Buffer

	err = tmpl.Execute(&t, data)
	if err != nil {
		return err
	}
	m.Body = t.String()

	return nil
}

func InitMessage(messageFile string) (*Message, error) {
	message := &Message{}
	err := helper.InitObjectFromJsonFile(messageFile, message)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func BuildMessage(workDir string, messageID string, bodyData map[string]interface{}) (*Message, error) {

	dir := path.Join(workDir, "messages")
	message, err := InitMessage(
		path.Join(dir ,messageID+".json"))

	if err != nil {
		return nil, err
	}

	if len(message.BodyTemplate) > 0 {
		tmpl, err := template.ParseFiles(path.Join(dir,message.BodyTemplate))

		if err != nil {
			return nil, err
		}

		var t bytes.Buffer
		err = tmpl.Execute(&t, bodyData)
		if err != nil {
			return nil, err
		}

		message.Body = t.String()
	}

	return message, nil
}
