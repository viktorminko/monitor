package config

import (
	"os"
	"github.com/viktorminko/monitor/helper"
)

type Configuration struct {
    Domain                       string
	RunPeriod                    int
	StatisticRunPeriod           int
	LogFile                      string
	Proxy                        string
}

func (c *Configuration) InitFromFile (filePath string) (error) {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	return helper.InitObjectFromJsonReader(file, c)
}