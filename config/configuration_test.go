package config

import (
	"bytes"
	"github.com/viktorminko/monitor/helper"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

type testData struct {
	sample          string
	expected        Configuration
	isErrorExpected bool
}

var sample = `{"Domain": "domain", "RunPeriod": "2s", "StatisticRunPeriod": 20, "LogFile": "my.log", "Proxy": "my_proxy"}`

var expectedConf = Configuration{
	"domain",
	Duration{time.Second * 2},
	Duration{time.Nanosecond * 20},
	"my.log",
	"my_proxy",
}

func checkConfig(t *testing.T, test testData) {
	conf := Configuration{}

	err := helper.InitObjectFromJsonReader(bytes.NewReader([]byte(test.sample)), &conf)

	if test.isErrorExpected && err == nil {
		t.Fatal("error expected but not returned")
	}

	if !test.isErrorExpected && err != nil {
		t.Fatalf("unexpected error returned: %s", err.Error())
	}

	if conf != test.expected {
		t.Fatalf(
			"unexpected configuration parsed. Expected %v, got %v",
			test.expected,
			conf,
		)
	}
}

func TestConfigurationReadFromFile(t *testing.T) {
	sample := []byte(sample)

	tmpfile, err := ioutil.TempFile("", "config_test")
	defer os.Remove(tmpfile.Name())

	if err != nil {
		t.Fatalf("unexpected error while creting temporary file: %s, error: %v", tmpfile.Name(), err)
	}

	if _, err := tmpfile.Write(sample); err != nil {
		t.Fatalf("unexpected error while writing config sample to temporary file: %v", err)
	}

	conf := Configuration{}

	dir, file := path.Split(tmpfile.Name())
	err = helper.InitObjectFromJsonFile(dir, file, &conf)

	if err != nil {
		t.Fatalf("unexpected error while reading configuration from temporary file: %s", err)
	}

	if conf != expectedConf {
		t.Fatalf(
			"unexpected configuration parsed. Expected %v, got %v",
			expectedConf,
			conf,
		)
	}

}

func TestConfiguration(t *testing.T) {
	tests := []testData{
		{
			`{: "invalid json"}`,
			Configuration{},
			true,
		},
		{
			`{"Domain": "test"}`,
			Configuration{Domain: "test"},
			false,
		},
		{
			`{"RunPeriod": "2s"}`,
			Configuration{RunPeriod: Duration{time.Second * 2}},
			false,
		},
		{
			`{"RunPeriod": 20}`,
			Configuration{RunPeriod: Duration{time.Nanosecond * 20}},
			false,
		},
		{
			`{"StatisticRunPeriod": "2s"}`,
			Configuration{StatisticRunPeriod: Duration{time.Second * 2}},
			false,
		},
		{
			//Invalid duration
			`{"StatisticRunPeriod": "ds"}`,
			Configuration{},
			true,
		},
		{
			//Invalid duration
			`{"StatisticRunPeriod": {}}`,
			Configuration{},
			true,
		},
		{
			//Invalid duration integer as string
			`{"StatisticRunPeriod": "20"}`,
			Configuration{},
			true,
		},
		{
			`{"LogFile": "my.log"}`,
			Configuration{LogFile: "my.log"},
			false,
		},
		{
			`{"Proxy": "my_proxy"}`,
			Configuration{Proxy: "my_proxy"},
			false,
		},
		{
			sample,
			expectedConf,
			false,
		},
	}

	for _, test := range tests {
		checkConfig(t, test)
	}
}
