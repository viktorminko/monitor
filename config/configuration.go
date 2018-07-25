package config

import (
	"encoding/json"
	"errors"
	"github.com/viktorminko/monitor/helper"
	"time"
)

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// Configuration contains main monitor settings
type Configuration struct {
	Domain             string
	RunPeriod          Duration
	StatisticRunPeriod Duration
	LogFile            string
	Proxy              string
}

// InitFromFile inits object from JSON file
func (c *Configuration) InitFromFile(filePath string) error {
	return helper.InitObjectFromJsonFile(filePath, c)
}
