package config

import (
	"bytes"
	"net/http"
	"text/template"
	"github.com/viktorminko/monitor/helper"
)

// Definition contains properties to single monitor test
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

// Definitions is a slice of monitor tests
type Definitions []Definition

// InitFromFile inits object from JSON file
func (ds *Definitions) InitFromFile (filePath string) (error) {
	return helper.InitObjectFromJsonFile(filePath, ds)
}

// Prepare modifies definition based on provided context
func (d *Definition) Prepare(data *Context) error {
	tmpl, err := template.New("tpl").Parse(d.URL)
	if err != nil {
		return err
	}

	var t bytes.Buffer
	err = tmpl.Execute(&t, &data)
	if err != nil {
		return err
	}
	d.URL = t.String()

	return nil
}
