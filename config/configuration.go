package config

import (
	"github.com/viktorminko/monitor/helper"
)

// Configuration contains main monitor settings
type Configuration struct {
    Domain                       string
	RunPeriod                    int
	StatisticRunPeriod           int
	LogFile                      string
	Proxy                        string
}

// InitFromFile inits object from JSON file
func (c *Configuration) InitFromFile (filePath string) (error) {
	return helper.InitObjectFromJsonFile(filePath, c)
}