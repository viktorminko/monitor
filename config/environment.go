package config

import (
	"os"
	"github.com/viktorminko/monitor/helper"
)

type Environment map[string]interface{}

func (e *Environment) InitFromFile (filePath string) (error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}

	return helper.InitObjectFromJsonReader(file, e)
}