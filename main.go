package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sirupsen/logrus"
	"github.com/solo-io/gloo-platform-lambda-demo/pkg/demo"
	"github.com/solo-io/gloo-platform-lambda-demo/pkg/helpers"
)

type GlooHttpRequest struct {
	Body        string            `json:"body"`
	Headers     map[string]string `json:"headers"`
	HttpMethod  string            `json:"httpMethod"`
	Path        string            `json:"path"`
	QueryString string            `json:"queryString"`
}

type GlooHttpResponse struct {
	Body       interface{}       `json:"body"`
	Headers    map[string]string `json:"headers"`
	StatusCode int               `json:"statusCode"`
}

var response404 GlooHttpResponse = GlooHttpResponse{
	StatusCode: 404,
	Body:       `{"message": "page not found"}`,
	Headers: map[string]string{
		"content-type": "application/json",
		"x-solo":       "test",
	},
}

var response400 GlooHttpResponse = GlooHttpResponse{
	StatusCode: 400,
	Body:       `{"message": "bad request"}`,
	Headers: map[string]string{
		"content-type": "application/json",
		"x-solo":       "test",
	},
}

func handleLambdaEvent(event GlooHttpRequest) (interface{}, error) {
	theUrl, err := url.Parse("?" + event.QueryString)
	if err != nil {
		logrus.Errorf("failed to parse url and querystring %s: %s", event.QueryString, err)
		return response400, nil
	}

	switch {
	case matchesRoute(event, "^/lambda/dump"):
		logrus.Debug("handling route /dump")
		body := map[string]interface{}{}
		result, _ := json.Marshal(event)
		body["requestEvent"] = string(result)
		return GlooHttpResponse{
			StatusCode: 200,
			Body:       body,
			Headers: map[string]string{
				"x-solo": "test",
			},
		}, nil
	case matchesRoute(event, "^/lambda/strings/reverse"):
		logrus.Debug("handling route /strings/reverse")
		return GlooHttpResponse{
			StatusCode: 200,
			Body:       map[string]interface{}{"output": demo.ReverseString(theUrl.Query().Get("input"))},
			Headers: map[string]string{
				"x-solo": "test",
			},
		}, nil
	case matchesRoute(event, "^/lambda/string"):
		logrus.Debug("handling route /lambda/string")
		return GlooHttpResponse{
			StatusCode: 200,
			Body:       "example string response",
			Headers: map[string]string{
				"x-solo": "test",
			},
		}, nil
	case matchesRoute(event, "^/lambda/headerlen"):
		logrus.Debug("handling route /lambda/headerlen")
		headers := map[string]string{}
		for i := 1; i <= 50; i++ {
			headers[fmt.Sprintf("x-solo-%d", i)] = strings.Repeat("x", i)
		}
		return GlooHttpResponse{
			StatusCode: 200,
			Body:       "sending headers of len 1-50",
			Headers:    headers,
		}, nil
	case matchesRoute(event, "^/lambda/echo"):
		logrus.Debug("handling route /echo")
		return GlooHttpResponse{
			StatusCode: 200,
			Body:       map[string]interface{}{"output": theUrl.Query().Get("input")},
			Headers: map[string]string{
				"x-solo": "test",
			},
		}, nil
	case matchesRoute(event, "^/lambda/json"):
		logrus.Debug("handling route /lambda/json")
		return GlooHttpResponse{
			StatusCode: 200,
			Body:       map[string]interface{}{"output": "example JSON response"},
			Headers: map[string]string{
				"x-solo": "test",
			},
		}, nil
	}

	// default respond 404
	return response404, nil
}

func matchesRoute(event GlooHttpRequest, route string) bool {
	if ok, _ := regexp.MatchString(route, event.Path); ok {
		return true
	}
	return false
}

func main() {
	lvl := logrus.InfoLevel
	if os.Getenv("LOG_LEVEL") != "" {
		lvl, _ = logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	}
	logrus.SetLevel(lvl)

	// when deployed to AWS Lambda invoke the handler directly
	if os.Getenv("AWS_LAMBDA_RUNTIME_API") != "" {
		logrus.Info("detected AWS lambda runtime, starting lambda handler")
		lambda.Start(handleLambdaEvent)
		return
	}

	// when not deployed to AWS Lambda, pass in local file for local development
	eventFile := flag.String("event", "./test/event-gloo-request-apigw-response.json", "an event file to load for local testing")
	flag.Parse()
	logrus.Infof("AWS lambda runtime not detected, invoking with test file %s", *eventFile)
	helpers.InvokeLambdaFromEventFileFor(*eventFile, handleLambdaEvent)
}
