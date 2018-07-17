package error

import (
	"encoding/json"
	"fmt"
	"log"
	"errors"
	"github.com/viktorminko/monitor/notifiers"
	"github.com/viktorminko/monitor/helper"
	"github.com/viktorminko/monitor/config"
)

type Fatal struct {
	Msg string
	Err error
}

func (e Fatal) Error() string {
	return e.Msg + "\r\n" + e.Err.Error()
}

type NonFatal struct {
	Msg string
	Err error
}

func (e NonFatal) Error() string {
	return e.Msg + "\r\n" + e.Err.Error()
}

type Test struct {
	Msg      string
	Request  config.Definition
	Code     int
	Response string
}

func (e Test) Error() string {
	return fmt.Sprintf("error: %v, Request %v, Code: %v, Response: %s", e.Msg, e.Request, e.Code, e.Response)
}

func Report(e error, senders *notifiers.Senders) error {

	if len(*senders) == 0 {
		return errors.New("senders list is empty, error will not be reported")
	}

	var mID string
	var mBody map[string]interface{}

	switch e.(type) {
	default:
		mID = "unexpected_error"
		mBody = map[string]interface{}{"err" : e.Error()}
	case Fatal:
		mID = "fatal_error"
		mBody = map[string]interface{}{"err" : e.(Fatal).Error()}
	case NonFatal:
		mID = "non_fatal_error"
		mBody = map[string]interface{}{"err" : e.(NonFatal).Error()}
	case Test:
		mID = "test_error"

		sample, _ := json.Marshal(e.(Test).Request.Sample)
		mBody = map[string]interface{}{
			"request":  e.(Test).Request.URL,
			"response": string(helper.FormatJSON([]byte(e.(Test).Response))),
			"sample":   string(helper.FormatJSON(sample)),
		}
	}


	senders.SendToAll(mID, mBody)

	return nil
}

func Check(err error) {
	if _, ok := err.(Fatal); ok {
		log.Fatalf("Fatal error: %v", err)
	}

	log.Printf("Non fatal error: %v", err)
	return

}
