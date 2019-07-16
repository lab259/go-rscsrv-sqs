[![CircleCI](https://circleci.com/gh/lab259/go-rscsrv-sqs.svg?style=shield)](https://circleci.com/gh/lab259/go-rscsrv-sqs)
[![codecov](https://codecov.io/gh/lab259/go-rscsrv-sqs/branch/master/graph/badge.svg)](https://codecov.io/gh/lab259/go-rscsrv-sqs)
[![GoDoc](https://godoc.org/github.com/lab259/go-rscsrv-sqs?status.svg)](http://godoc.org/github.com/lab259/go-rscsrv-sqs)
[![Go Report Card](https://goreportcard.com/badge/github.com/lab259/go-rscsrv-sqs)](https://goreportcard.com/report/github.com/lab259/go-rscsrv-sqs)

# go-rscsrv-sqs

The go-rscsrv-sqs is the [lab259/go-rscsrv](//github.com/lab259/go-rscsrv) service for
the Amazon Simple Queue Service (SQS).

It wraps the logic of dealing with credentials and keeping the session.

## Dependencies

It depends on the [lab259/go-rscsrv](//github.com/lab259/go-rscsrv) (and its dependencies,
of course) itself and the [aws/aws-sdk-go](//github.com/aws/aws-sdk-go) library.

## Installation

First, fetch the library to the repository.

    go get github.com/lab259/go-rscsrv-sqs

## Usage

The service is designed to be "extended" and not used directly.

**srv.go**

```Go
package mail

import (
	"github.com/lab259/go-rscsrv"
	"github.com/lab259/go-rscsrv-sqs"
)

type MessageSQSService struct {
	sqssrv.SQSService
}

func (service *MessageSQSService) LoadConfiguration() (interface{}, error) {
	var configuration sqssrv.SQSServiceConfiguration

	configurationLoader := rscsrv.NewFileConfigurationLoader("/etc/mail")
	configurationUnmarshaler := &rscsrv.ConfigurationUnmarshalerYaml{}

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

```Go
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

## Development

```bash
git clone git@github.com:lab259/go-rscsrv-sqs.git # clone the project
cd go-rscsrv-sqs                                  # enter the directory
make dcup                                         # start the ElasticMQ
go mod download                                   # download the dependencies
make test                                         # run the tests
```
