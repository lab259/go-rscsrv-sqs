package sqssrv

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/lab259/http"
	"net/url"
	"path"
)

// RedigoServiceConfiguration is the configuration for the `RedigoService`
type SQSServiceConfiguration struct {
	QUrl     string `yaml:"q_url"`
	Region   string `yaml:"region"`
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
}

type CredentialsFromStruct struct {
	credentials *SQSServiceConfiguration
}

func NewCredentialsFromStruct(credentials *SQSServiceConfiguration) *CredentialsFromStruct {
	return &CredentialsFromStruct{
		credentials: credentials,
	}
}

func (c *CredentialsFromStruct) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     c.credentials.Key,
		SecretAccessKey: c.credentials.Secret,
	}, nil
}

func (*CredentialsFromStruct) IsExpired() bool {
	return false
}

// SQSService is the service which manages a service queue on the AWS.
type SQSService struct {
	running       bool
	awsSQS        *sqs.SQS
	Configuration SQSServiceConfiguration
}

// LoadConfiguration returns
func (service *SQSService) LoadConfiguration() (interface{}, error) {
	return nil, errors.New("not implemented")
}

// ApplyConfiguration applies a given configuration to the service.
func (service *SQSService) ApplyConfiguration(configuration interface{}) error {
	switch c := configuration.(type) {
	case SQSServiceConfiguration:
		service.Configuration = c
		return nil
	case *SQSServiceConfiguration:
		service.Configuration = *c
		return nil
	}
	return http.ErrWrongConfigurationInformed
}

// Restart stops and then starts the service again.
func (service *SQSService) Restart() error {
	if service.running {
		err := service.Stop()
		if err != nil {
			return err
		}
	}
	return service.Start()
}

// Start starts the redis pool.
func (service *SQSService) Start() error {
	if !service.running {
		conf := aws.Config{
			Credentials: credentials.NewCredentials(NewCredentialsFromStruct(&service.Configuration)),
		}

		if service.Configuration.Endpoint != "" {
			conf.Endpoint = aws.String(service.Configuration.Endpoint)
		}

		if service.Configuration.Region == "" {
			conf.Region = aws.String("sa-east-1")
		} else {
			conf.Region = aws.String(service.Configuration.Region)
		}

		sess, err := session.NewSessionWithOptions(session.Options{
			Config: conf,
		})
		if err != nil {
			return err
		}
		service.awsSQS = sqs.New(sess)

		qurl, err := url.Parse(service.Configuration.QUrl)
		if err != nil {
			return fmt.Errorf("could not parse QUrl: %s", err.Error())
		}

		listQueuesOutput, err := service.awsSQS.ListQueues(&sqs.ListQueuesInput{
			QueueNamePrefix: aws.String(path.Base(qurl.Path)),
		})
		if err != nil {
			return err
		}
		err = func() error {
			for _, q := range listQueuesOutput.QueueUrls {
				if aws.StringValue(q) == service.Configuration.QUrl {
					return nil
				}
			}
			return fmt.Errorf("queue %s not found", service.Configuration.QUrl)
		}()
		if err != nil {
			return err
		}

		service.running = true
	}
	return nil
}

// Stop erases the aws client reference.
func (service *SQSService) Stop() error {
	if service.running {
		service.awsSQS = nil
		service.running = false
	}
	return nil
}

// RunWithSQS runs a handler passing the reference of a `sqs.SQS` client.
func (service *SQSService) RunWithSQS(handler func(client *sqs.SQS) error) error {
	if service.running {
		return handler(service.awsSQS)
	}
	return http.ErrServiceNotRunning
}

// SendMessage is a wrapper for the `sqs.SQS.SendMessage`.
func (service *SQSService) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.SendMessage(input)
	}
	return nil, http.ErrServiceNotRunning
}

// SendMessageWithContext is a wrapper for the `sqs.SQS.SendMessage`.
func (service *SQSService) SendMessageWithContext(ctx context.Context, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.SendMessageWithContext(ctx, input)
	}
	return nil, http.ErrServiceNotRunning
}

// SendMessageBatch is a wrapper for the `sqs.SQS.SendMessageBatch`.
func (service *SQSService) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.SendMessageBatch(input)
	}
	return nil, http.ErrServiceNotRunning
}

// SendMessageBatchWithContext is a wrapper for the `sqs.SQS.SendMessageBatchWithContext`.
func (service *SQSService) SendMessageBatchWithContext(ctx context.Context, input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.SendMessageBatchWithContext(ctx, input)
	}
	return nil, http.ErrServiceNotRunning
}

// ReceiveMessage is a wrapper for the `sqs.SQS.ReceiveMessage`.
func (service *SQSService) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.ReceiveMessage(input)
	}
	return nil, http.ErrServiceNotRunning
}

// ReceiveMessageWithContext is a wrapper for the `sqs.SQS.ReceiveMessageWithContext`.
func (service *SQSService) ReceiveMessageWithContext(ctx context.Context, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.ReceiveMessageWithContext(ctx, input)
	}
	return nil, http.ErrServiceNotRunning
}

// DeleteMessage is a wrapper for the `sqs.SQS.DeleteMessage`.
func (service *SQSService) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.DeleteMessage(input)
	}
	return nil, http.ErrServiceNotRunning
}

// DeleteMessageWithContext is a wrapper for the `sqs.SQS.DeleteMessageWithContext`.
func (service *SQSService) DeleteMessageWithContext(ctx context.Context, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.DeleteMessageWithContext(ctx, input)
	}
	return nil, http.ErrServiceNotRunning
}

// DeleteMessageBatch is a wrapper for the `sqs.SQS.DeleteMessageBatch`.
func (service *SQSService) DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.DeleteMessageBatch(input)
	}
	return nil, http.ErrServiceNotRunning
}

// DeleteMessageBatchWithContext is a wrapper for the `sqs.SQS.DeleteMessageBatchWithContext`.
func (service *SQSService) DeleteMessageBatchWithContext(ctx context.Context, input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	if service.running {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.awsSQS.DeleteMessageBatchWithContext(ctx, input)
	}
	return nil, http.ErrServiceNotRunning
}
