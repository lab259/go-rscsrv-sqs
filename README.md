[![CircleCI](https://circleci.com/gh/lab259/http-sqs-service.svg?style=shield)](https://circleci.com/gh/lab259/http-sqs-service)
[![codecov](https://codecov.io/gh/lab259/http-sqs-service/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/http-sqs-service)
[![GoDoc](https://godoc.org/github.com/lab259/http-sqs-service?status.svg)](http://godoc.org/github.com/lab259/http-sqs-service)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab259/http-sqs-service)](https://goreportcard.com/report/github.com/lab259/http-sqs-service)

# http-sqs-service

The http-sqs-service is the [lab259/http](//github.com/lab259/http) service for
the Amazon Simple Queue Service (SQS).

It wraps the logic of dealing with credentials and keeping the session.

## Dependencies

It depends on the [lab259/http](//github.com/lab259/http) (and its dependencies,
of course) itself and the [aws/aws-sdk-go](//github.com/aws/aws-sdk-go) library.

## Installation

First, fetch the library to the repository.

	go get github.com/lab259/http-sqs-service

## Usage

The service is designed to be "extended" and not used directly.

**srv.go**

```Go
package mail

import (
	"github.com/lab259/http"
	"github.com/lab259/http-sqs-service"
)

type MessageSQSService struct {
	sqssrv.SQSService
}

func (service *MessageSQSService) LoadConfiguration() (interface{}, error) {
	var configuration sqssrv.SQSServiceConfiguration

	configurationLoader := http.NewFileConfigurationLoader("/etc/mail")
	configurationUnmarshaler := &http.ConfigurationUnmarshalerYaml{}

	config, err := configurationLoader.Load(file)
	if err != nil {
		return err
	}
	return configurationUnmarshaler.Unmarshal(config, dst)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

```

**example.go**
```
// ...

var mq rscsrsv.MessageSQSService

func init() {
	configuration, err := mq.LoadConfiguration()
	if err != nil {
		panic(err)
	}
	err = mq.ApplyConfiguration(configuration)
	if err != nil {
		panic(err)
	}
	err = mq.Start()
	if err != nil {
		panic(err)
	}
}

// enqueueMessage enqueues a message
func enqueueMessage(message string) {
	output, err := service.SendMessage(&sqs.SendMessageInput{
		MessageBody: aws.String("this is the content of the message"),
	})
	if err != nil {
		panic(err)
	}
	// ... handles the return of the SendMessage
}

// ...
```