package test

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
	"github.com/viktorminko/monitor/config"
	"github.com/viktorminko/monitor/method"
	chttp "github.com/viktorminko/monitor/http"
	"github.com/viktorminko/monitor/authorization"
)

func TestTest_IsNeedToRun(t *testing.T) {

	lastExecuted := time.Now()

	runPeriod := 2

	apiTest := Test{
		LastExecutedAt: lastExecuted,
		APIMethod: &method.Data{
			RunPeriod: runPeriod,
		},
	}

	//Exactly in one period
	if apiTest.IsNeedToRun(lastExecuted.Add(time.Duration(runPeriod) * time.Second)) {
		t.Error("Test should not be executed exactly after one period")
	}

	//Later then one period
	if !apiTest.IsNeedToRun(lastExecuted.Add(time.Duration(runPeriod+1) * time.Second)) {
		t.Error("Test should be executed later then 1 period")
	}

	//Earlier then one period
	if apiTest.IsNeedToRun(lastExecuted.Add(time.Duration(runPeriod-1) * time.Second)) {
		t.Error("Test should not be executed earlier then 1 period")
	}

	//Should be executed if LastExecutedAt is 0
	apiTest.LastExecutedAt = time.Time{}
	if !apiTest.IsNeedToRun(time.Time{}.Add(time.Duration(runPeriod-1) * time.Second)) {
		t.Error("Test should be executed if it was never executed yet")
	}
}

func TestTest_Run(t *testing.T) {

	//Create error channel to handle errors and put in separate goroutine
	//We need this, so API caller will not halt
	errChan := make(chan error)
	go func() {
		for {
			select {
			case <-errChan:
			}
		}
	}()
	apiCaller := &APICaller{
		errChan,
		&chttp.Client{},
	}

	//Random auth token
	//Server will check if it is presented
	rand.Seed(time.Now().Unix())
	token := strconv.Itoa(rand.Intn(1000))

	//Server will delay its response
	//Then we test if correct response value is sent to statistics
	responseDelay := time.Duration(1) * time.Second

	//How server handles request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "Bearer "+token != r.Header.Get("Authorization") {
			t.Error("Valid authorization token was not provided")
		}

		time.Sleep(responseDelay)
	}))
	defer server.Close()

	//APIMethod we are about to run
	apiTest := Test{
		Domain: server.URL,
		APIMethod: &method.Data{
			ResponseCode: 200,
			Sample:       false,
		},
	}

	//Function "Run" will send results to statistics channel
	statsChan := make(chan ExecutionData)

	//Avoid races
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		res := <-statsChan

		if res.Err != nil {
			t.Errorf("Unexpected error return %s", res.Err.Error())
		}

		if res.ResponseTime < responseDelay {
			t.Errorf("Expected response time more then %s, got %s", responseDelay, res.ResponseTime)
		}

		wg.Done()

	}()

	apiTest.Run(&authorization.Token{Token: token}, statsChan, apiCaller)
	wg.Wait()
}

func TestPrepareTests(t *testing.T) {

	domain := "my_domain"
	methods := []method.Data{{}, {}}

	p, err := Prepare(
		methods,
		&config.Environment{}, domain)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if p[0].Domain != domain || p[1].Domain != domain {
		t.Error("Domain is not inserted while test preparation")
	}

	if p[0].APIMethod != &methods[0] || p[1].APIMethod != &methods[1] {
		t.Error("Methods were not inserted in test correctly")
	}

	if !p[0].LastExecutedAt.IsZero() || !p[0].LastExecutedAt.IsZero() {
		t.Error("Tests time was not set correctly")
	}

}
