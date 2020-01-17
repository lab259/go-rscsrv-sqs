package services

import (
	"github.com/lab259/go-rscsrv"
	sqssrv "github.com/lab259/go-rscsrv-sqs"
)

var DefaultSQSService MessageSQSService

type MessageSQSService struct {
	sqssrv.SQSService
}

func (*MessageSQSService) Name() string {
	return "SQS Service"
}

func (service *MessageSQSService) LoadConfiguration() (interface{}, error) {
	var configuration sqssrv.SQSServiceConfiguration

	configurationLoader := rscsrv.NewFileConfigurationLoader("./")
	configurationUnmarshaler := &rscsrv.ConfigurationUnmarshalerYaml{}

	config, err := configurationLoader.Load("conf.yml")
	if err != nil {
		return nil, err
	}
	err = configurationUnmarshaler.Unmarshal(config, &configuration)
	if err != nil {
		return nil, err
	}
	return configuration, nil
}

// Start starts the service. If successful nil will be returned, otherwise
// the error.
func (service *MessageSQSService) Start() error {
	err := service.SQSService.Start()
	if err != nil {
		return err
	}
	return DefaultPromService.Register(service.Collector)
}
