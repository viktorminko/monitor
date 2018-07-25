package request

import (
	"fmt"
	"github.com/viktorminko/monitor/authorization"
	"github.com/viktorminko/monitor/config"
	cerror "github.com/viktorminko/monitor/error"
	"github.com/viktorminko/monitor/helper"
	chttp "github.com/viktorminko/monitor/http"
	"io/ioutil"
	"net/http"
	"strings"
)

type Runner struct {
	ErrorChannel chan<- error
	Client       *chttp.Client
}

func (a *Runner) RunTest(definition *config.Definition, domain string, token *authorization.Token) error {
	req, err := http.NewRequest(definition.HTTPMethod, domain+definition.URL, strings.NewReader(definition.Payload))
	if err != nil {
		a.ErrorChannel <- cerror.NonFatal{"error occurred while calling http", err}
		return err
	}

	if token != nil {
		req.Header.Add("authorization", "Bearer "+token.Token)
	}

	res, err := a.Client.Call(
		req,
		definition.TimeOut.Duration,
	)
	if err != nil {
		a.ErrorChannel <- cerror.NonFatal{"error occurred while calling http", err}
		return err
	}

	//Check expected response code
	if definition.ResponseCode != res.StatusCode {
		err = cerror.NonFatal{
			"unexpected HTTP response code",
			fmt.Errorf("expected: %d, received: %d", definition.ResponseCode, res.StatusCode)}

		a.ErrorChannel <- err
		return err
	}

	//If sample is false, no need to check it
	if false == definition.Sample {
		return nil
	}

	body, _ := ioutil.ReadAll(res.Body)

	isExpectedResponse, err := helper.AreEqualJSON(string(body), definition.Sample)
	if err != nil {
		a.ErrorChannel <- cerror.NonFatal{"error occurred while comparing API response and sample", err}
		return err
	}

	if !isExpectedResponse {
		err := cerror.Test{
			"unexpected API response",
			*definition,
			res.StatusCode,
			string(body),
		}
		a.ErrorChannel <- err
		return err
	}

	return nil
}
