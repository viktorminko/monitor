package main

import (
	"flag"
	"github.com/viktorminko/monitor/authorization"
	"github.com/viktorminko/monitor/config"
	cerror "github.com/viktorminko/monitor/error"
	"github.com/viktorminko/monitor/helper"
	chttp "github.com/viktorminko/monitor/http"
	"github.com/viktorminko/monitor/notifiers"
	"github.com/viktorminko/monitor/request"
	"github.com/viktorminko/monitor/statistic"
	"log"
	"net/url"
	"os"
	"path"
)

const configFolder = "config"
const configFile = "config.json"
const testsFile = "tests.json"
const environmentFile = "environment.json"
const authorizationConfigurationFile = "authorization.json"

// Init senders fro notifications
// Currently supported: email, telegram
func initSenders(workDir string) *notifiers.Senders {
	//Init message senders
	var s notifiers.Senders

	//Init email sender
	emailSender, err := notifiers.InitEmailSender(workDir)
	if err != nil {
		cerror.Check(cerror.NonFatal{"unable to init email sender", err})
	} else {
		s.Add(emailSender)
	}

	//Init telegram sender
	telegram, err := notifiers.InitTelegramSender(workDir)
	if err != nil {
		cerror.Check(cerror.NonFatal{"unable to init telegram sender", err})
	} else {
		s.Add(telegram)
	}

	return &s
}

func main() {

	log.Println("API monitor started")

	var workDir string
	flag.StringVar(&workDir, "workdir", configFolder, "working directory to load configuration files from")

	flag.Parse()

	if _, err := os.Stat("workDir"); err == nil {
		cerror.Check(cerror.Fatal{"incorrect working directory", err})
	}

	configuration := &config.Configuration{}
	if err := configuration.InitFromFile(path.Join(path.Dir(workDir), configFile)); err != nil {
		cerror.Check(cerror.Fatal{"error loading configuration", err})
	}

	//Configure log
	if len(configuration.LogFile) > 0 {
		if err := helper.PrepareLog(configuration.LogFile); err != nil {
			cerror.Check(cerror.NonFatal{"error configuring log", err})
		}
	}

	//Init tests
	rawTests := config.Definitions{}
	if err := rawTests.InitFromFile(path.Join(path.Dir(workDir), testsFile)); err != nil {
		cerror.Check(cerror.Fatal{"error loading tests", err})
	}

	//Init environment
	env := &config.Context{}
	if err := env.InitFromFile(path.Join(path.Dir(workDir), environmentFile)); err != nil {
		cerror.Check(cerror.Fatal{"error loading environment", err})
	}

	//Update tests based on data from environment file
	tests, err := request.Prepare(rawTests, env, configuration.Domain)
	if err != nil {
		cerror.Check(cerror.NonFatal{"error occurred while preparing raw tests", err})
	}

	//Init authorization configuration
	authConf := &config.AuthorizationConfiguration{}
	if err := authConf.InitFromFile(path.Join(path.Dir(workDir), authorizationConfigurationFile)); err != nil {
		cerror.Check(cerror.NonFatal{"error loading authorization configuration", err})
	}

	//Init message senders
	senders := initSenders(workDir)

	//Run Error handler
	errorChannel := (&cerror.Handler{}).Run(senders)

	//Send startup message
	(&StartupReporter{errorChannel}).Send(
		configuration,
		authConf,
		rawTests,
		senders,
	)

	testStatsChan,
		authStatsChan,
		statsRequester := (&statistic.Collector{&statistic.Monitor{
		&statistic.Suite{
			nil,
		},
		&statistic.Authorization{
			statistic.Statistic{0, 0, 0},
			nil,
			0,
		},
	}}).Run()

	//Run statistics reporter
	(&statistic.Reporter{errorChannel}).
		Run(
			configuration.StatisticRunPeriod.Duration,
			statsRequester,
			senders,
		)

	//Set proxy for http request if necessary
	client := &chttp.Client{}
	if len(configuration.Proxy) > 0 {
		proxyURL, err := url.Parse(configuration.Proxy)
		if err == nil {
			client.Proxy = proxyURL
		}
	}

	var authHandler *authorization.Handler

	if len(authConf.AuthorizationURL) != 0 {
		authHandler = &authorization.Handler{
			&authorization.HTTPAuthorizer{
				configuration.Domain + authConf.AuthorizationURL,
				authConf.GetAuthorizationTokenTimeout,
				authConf.AppID,
				authConf.AppSecret,
				client,
			},
			authStatsChan,
		}
	}
	//Run tests runner
	(&TestsRunner{
		&request.Suite{
			tests,
			testStatsChan,
			errorChannel,
		},
		configuration.RunPeriod.Duration,
		testStatsChan,
		errorChannel,
	}).Run(
		authHandler,
		&request.Runner{
			errorChannel,
			client,
		},
	)

}
