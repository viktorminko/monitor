package error

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/viktorminko/monitor/config"
	"github.com/viktorminko/monitor/helper"
	"github.com/viktorminko/monitor/notifiers"
	"log"
)

// Fatal is a fatal error, should stop immediately
type Fatal struct {
	Msg string
	Err error
}

// Error returns error message with internal error data
func (e Fatal) Error() string {
	return e.Msg + "\r\n" + e.Err.Error()
}

// NonFatal is a non fatal error, should not stop execution
type NonFatal struct {
	Msg string
	Err error
}

// Error returns error message with internal error data
func (e NonFatal) Error() string {
	return e.Msg + "\r\n" + e.Err.Error()
}

// Test is a test error, contains data of test execution
type Test struct {
	Msg      string
	Request  config.Definition
	Code     int
	Response string
}

// Error returns error message using test information
func (e Test) Error() string {
	return fmt.Sprintf("error: %v, Request %v, Code: %v, Response: %s", e.Msg, e.Request, e.Code, e.Response)
}

// Report passes error message to notifiers
func Report(e error, senders *notifiers.Senders) error {

	if len(*senders) == 0 {
		return errors.New("senders list is empty, error will not be reported")
	}

	var mID string
	var mBody map[string]interface{}

	switch e.(type) {
	default:
		mID = "unexpected_error"
		mBody = map[string]interface{}{"err": e.Error()}
	case Fatal:
		mID = "fatal_error"
		mBody = map[string]interface{}{"err": e.(Fatal).Error()}
	case NonFatal:
		mID = "non_fatal_error"
		mBody = map[string]interface{}{"err": e.(NonFatal).Error()}
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

// Check performs required actions based on error type
func Check(err error) {
	if _, ok := err.(Fatal); ok {
		log.Fatalf("Fatal error: %v", err)
	}

	log.Printf("Non fatal error: %v", err)
	return

}
