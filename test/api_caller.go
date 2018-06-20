package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"github.com/viktorminko/monitor/helper"
	chttp "github.com/viktorminko/monitor/http"
	"github.com/viktorminko/monitor/authorization"
	cerror "github.com/viktorminko/monitor/error"
	"github.com/viktorminko/monitor/method"
)

type APICaller struct {
	ErrorChannel chan<- error
	APIClient    *chttp.Client
}

func (a *APICaller) RunApiMethod(apiMethod *method.Data, domain string, token *authorization.Token) error {
	req, err := http.NewRequest(apiMethod.HTTPMethod, domain+apiMethod.URL, strings.NewReader(apiMethod.Payload))
	if err != nil {
		a.ErrorChannel <- cerror.NonFatal{"error occurred while calling http", err}
		return err
	}

	if token != nil {
		req.Header.Add("authorization", "Bearer "+token.Token)
	}

	res, err := a.APIClient.Call(
		req,
		time.Duration(apiMethod.TimeOut)*time.Second,
	)
	if err != nil {
		a.ErrorChannel <- cerror.NonFatal{"error occurred while calling http", err}
		return err
	}

	//Check expected response code
	if apiMethod.ResponseCode != res.StatusCode {
		err = cerror.NonFatal{
			"unexpected HTTP response code",
			fmt.Errorf("expected: %d, received: %d", apiMethod.ResponseCode, res.StatusCode)}

		a.ErrorChannel <- err
		return err
	}

	//If sample is false, no need to check it
	if false == apiMethod.Sample {
		return nil
	}

	body, _ := ioutil.ReadAll(res.Body)

	isExpectedResponse, err := helper.AreEqualJSON(string(body), apiMethod.Sample)
	if err != nil {
		a.ErrorChannel <- cerror.NonFatal{"error occurred while comparing API response and sample", err}
		return err
	}

	if !isExpectedResponse {
		err := cerror.Test{
			"unexpected API response",
			*apiMethod,
			res.StatusCode,
			string(body),
		}
		a.ErrorChannel <- err
		return err
	}

	return nil
}
