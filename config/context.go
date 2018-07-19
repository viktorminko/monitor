package config

import (
	"github.com/viktorminko/monitor/helper"
)

// Context contains parameters to be used in test suite
type Context map[string]interface{}

// InitFromFile inits object from JSON file
func (c *Context) InitFromFile (filePath string) (error) {
	return helper.InitObjectFromJsonFile(filePath, c)
}