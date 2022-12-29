package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

func InvokeLambdaFromEventFileFor[e any](eventFile string, handlerFunc func(event e) (interface{}, error)) {
	data, err := ioutil.ReadFile(eventFile)
	if err != nil {
		logrus.Fatalf("failed to read %s: %s", "", err)
	}

	var payload e
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.Decode(&payload)

	result, err := handlerFunc(payload)
	if err != nil {
		logrus.Fatalf("failed to handle request: %s", err)
	}

	resultString, _ := json.Marshal(&result)
	fmt.Println(string(resultString))
}
