package helper

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"io"
	"bytes"
)

func InitObjectFromJsonFile(filename string, obj interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	return InitObjectFromJsonReader(file, obj)
}

func InitObjectFromJsonReader(reader io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(obj)

	return err
}

func PrepareLog(logfile string) error {
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(f)

	return nil
}

/**
 * Deep compare 2 jsons
 */
func AreEqualJSON(s1 string, jsonData interface{}) (bool, error) {
	var o1 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, jsonData), nil
}

func FormatJSON(s []byte) []byte {
	var sPretty bytes.Buffer
	err := json.Indent(&sPretty, s, "", "    ")
	if err != nil {
		return s
	}
	return sPretty.Bytes()
}
