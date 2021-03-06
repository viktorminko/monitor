package request

import (
	"github.com/viktorminko/monitor/pkg/authorization"
	"github.com/viktorminko/monitor/pkg/config"
	chttp "github.com/viktorminko/monitor/pkg/http"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTest_IsNeedToRun(t *testing.T) {

	lastExecuted := time.Now()

	runPeriod := config.Duration{2}

	urlTest := Request{
		LastExecutedAt: lastExecuted,
		Definition: &config.Definition{
			RunPeriod: runPeriod,
		},
	}

	//Exactly in one period
	if urlTest.IsNeedToRun(lastExecuted.Add(runPeriod.Duration)) {
		t.Error("Request should not be executed exactly after one period")
	}

	//Later then one period
	if !urlTest.IsNeedToRun(lastExecuted.Add(runPeriod.Duration + 1)) {
		t.Error("Request should be executed later then 1 period")
	}

	//Earlier then one period
	if urlTest.IsNeedToRun(lastExecuted.Add(runPeriod.Duration - 1)) {
		t.Error("Request should not be executed earlier then 1 period")
	}

	//Should be executed if LastExecutedAt is 0
	urlTest.LastExecutedAt = time.Time{}
	if !urlTest.IsNeedToRun(time.Time{}.Add(runPeriod.Duration - 1)) {
		t.Error("Request should be executed if it was never executed yet")
	}
}

func TestTest_Run(t *testing.T) {

	//Create error channel to handle errors and put in separate goroutine
	//We need this, so caller will not halt
	errChan := make(chan error)
	go func() {
		for {
			select {
			case <-errChan:
			}
		}
	}()
	urlCaller := &Runner{
		errChan,
		&chttp.Client{},
	}

	//Random auth token
	//Server will check if it is presented
	rand.Seed(time.Now().Unix())
	token := strconv.Itoa(rand.Intn(1000))

	//Server will delay its response
	//Then we request if correct response value is sent to statistics
	responseDelay := time.Duration(1) * time.Second

	//How server handles request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "Bearer "+token != r.Header.Get("Authorization") {
			t.Error("Valid authorization token was not provided")
		}

		time.Sleep(responseDelay)
	}))
	defer server.Close()

	//Definition we are about to run
	urlTest := Request{
		Domain: server.URL,
		Definition: &config.Definition{
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
			t.Errorf("unexpected error return %s", res.Err.Error())
		}

		if res.ResponseTime < responseDelay {
			t.Errorf("expected response time more then %s, got %s", responseDelay, res.ResponseTime)
		}

		wg.Done()

	}()

	urlTest.Run(&authorization.Token{Token: token}, statsChan, urlCaller)
	wg.Wait()
}

func TestPrepareTests(t *testing.T) {

	domain := "my_domain"
	methods := []config.Definition{{}, {}}

	p, err := Prepare(
		methods,
		&config.Context{}, domain)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if p[0].Domain != domain || p[1].Domain != domain {
		t.Error("Domain is not inserted while request preparation")
	}

	if p[0].Definition != &methods[0] || p[1].Definition != &methods[1] {
		t.Error("Methods were not inserted in request correctly")
	}

	if !p[0].LastExecutedAt.IsZero() || !p[0].LastExecutedAt.IsZero() {
		t.Error("Definitions time was not set correctly")
	}
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}
