package config

import (
	"bytes"
	"net/http"
	"text/template"
	"os"
	"github.com/viktorminko/monitor/helper"
)

type Definition struct {
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

type Definitions []Definition

func (a *Definitions) InitFromFile (filePath string) (error) {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	return helper.InitObjectFromJsonReader(file, a)
}

func (m *Definition) Prepare(data *Environment) error {
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
