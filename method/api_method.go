package method

import (
	"bytes"
	"net/http"
	"text/template"
	"os"
	"github.com/viktorminko/monitor/helper"
	"github.com/viktorminko/monitor/config"
)

type APITests []Data

func (a *APITests) InitFromFile (filePath string) (error) {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	return helper.InitObjectFromJsonReader(file, a)
}

type Data struct {
	ID           string
	URL          string
	Sample       interface{}
	HTTPMethod   string
	Header       http.Header
	Payload      string
	RunPeriod    int
	TimeOut      int
	ResponseCode int
}

func (m *Data) Prepare(data *config.Environment) error {
	tmpl, err := template.New("tpl").Parse(m.URL)
	if err != nil {
		return err
	}

	var t bytes.Buffer
	err = tmpl.Execute(&t, &data)
	if err != nil {
		return err
	}
	m.URL = t.String()

	return nil
}
