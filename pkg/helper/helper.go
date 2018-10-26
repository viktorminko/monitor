package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"reflect"
)

// InitObjectFromJsonFile inits object from JSON file
func InitObjectFromJsonFile(dir string, fileName string, obj interface{}) error {
	file, err := os.Open(path.Join(dir, fileName))
	if err != nil {
		return err
	}

	return InitObjectFromJsonReader(file, obj)
}

// InitObjectFromJsonFile inits object from JSON reader
func InitObjectFromJsonReader(reader io.Reader, obj interface{}) error {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(obj)

	return err
}

// PrepareLog sets log output to provided file
func PrepareLog(logfile string) error {
	f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(f)

	return nil
}

// AreEqualJSON performs deep comparison of 2 JSONs
func AreEqualJSON(s1 string, jsonData interface{}) (bool, error) {
	var o1 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("error mashalling string 1 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, jsonData), nil
}

// FormatJSON makes JSON look nice
func FormatJSON(s []byte) []byte {
	var sPretty bytes.Buffer
	err := json.Indent(&sPretty, s, "", "    ")
	if err != nil {
		return s
	}
	return sPretty.Bytes()
}
