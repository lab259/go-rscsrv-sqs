package sqssrv

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	rscsrv "github.com/lab259/go-rscsrv"
)

// SQSServiceConfiguration is the configuration for the `SQS`
type SQSServiceConfiguration struct {
	QUrl     string `yaml:"q_url"`
	Region   string `yaml:"region"`
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
}

// CredentialsFromStruct define credentials from sqs configuration
type CredentialsFromStruct struct {
	credentials *SQSServiceConfiguration
}

// NewCredentialsFromStruct is the configuration for the `SQS`
func NewCredentialsFromStruct(credentials *SQSServiceConfiguration) *CredentialsFromStruct {
	return &CredentialsFromStruct{
		credentials: credentials,
	}
}

// Retrieve return the AWS credentials
func (c *CredentialsFromStruct) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     c.credentials.Key,
		SecretAccessKey: c.credentials.Secret,
	}, nil
}

// IsExpired return if credentials is expired
func (*CredentialsFromStruct) IsExpired() bool {
	return false
}

// SQSService is the service which manages a service queue on the AWS.
type SQSService struct {
	m             sync.RWMutex
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
	return rscsrv.ErrWrongConfigurationInformed
}

// Restart stops and then starts the service again.
func (service *SQSService) Restart() error {
	if err := service.Stop(); err != nil {
		if err != nil {
			return err
		}
	}
	return service.Start()
}

// Start starts the service pool.
func (service *SQSService) Start() error {
	if !service.isRunning() {
		service.m.Lock()
		defer service.m.Unlock()

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

		confQURLParsed, err := url.Parse(service.Configuration.QUrl)
		if err != nil {
			return fmt.Errorf("could not parse the qurl: %s (%s)", service.Configuration.QUrl, err.Error())
		}

		listQueuesOutput, err := service.awsSQS.ListQueues(&sqs.ListQueuesInput{
			QueueNamePrefix: aws.String(path.Base(confQURLParsed.Path)),
		})
		if err != nil {
			return err
		}
		err = func() error {
			for _, q := range listQueuesOutput.QueueUrls {
				qURLParsed, err := url.Parse(aws.StringValue(q))
				if err != nil {
					return fmt.Errorf("could not parse the qurl: %s (%s)", aws.StringValue(q), err.Error())
				}
				if path.Base(qURLParsed.Path) == path.Base(confQURLParsed.Path) {
					return nil
				}
			}
			return fmt.Errorf("queue %s not found", service.Configuration.QUrl)
		}()
		if err != nil {
			return err
		}

	}

	return nil
}

func (service *SQSService) isRunning() bool {
	service.m.RLock()
	defer service.m.RUnlock()
	return service.awsSQS != nil
}

func (service *SQSService) getSQS() *sqs.SQS {
	service.m.RLock()
	defer service.m.RUnlock()
	return service.awsSQS
}

// Stop erases the aws client reference.
func (service *SQSService) Stop() error {
	if service.isRunning() {
		service.m.Lock()
		service.awsSQS = nil
		service.m.Unlock()
	}
	return nil
}

// RunWithSQS runs a handler passing the reference of a `sqs.SQS` client.
func (service *SQSService) RunWithSQS(handler func(client *sqs.SQS) error) error {
	if service.isRunning() {
		return handler(service.getSQS())
	}
	return rscsrv.ErrServiceNotRunning
}

// SendMessage is a wrapper for the `sqs.SQS.SendMessage`.
func (service *SQSService) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().SendMessage(input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// SendMessageWithContext is a wrapper for the `sqs.SQS.SendMessage`.
func (service *SQSService) SendMessageWithContext(ctx context.Context, input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().SendMessageWithContext(ctx, input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// SendMessageBatch is a wrapper for the `sqs.SQS.SendMessageBatch`.
func (service *SQSService) SendMessageBatch(input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().SendMessageBatch(input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// SendMessageBatchWithContext is a wrapper for the `sqs.SQS.SendMessageBatchWithContext`.
func (service *SQSService) SendMessageBatchWithContext(ctx context.Context, input *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().SendMessageBatchWithContext(ctx, input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// ReceiveMessage is a wrapper for the `sqs.SQS.ReceiveMessage`.
func (service *SQSService) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().ReceiveMessage(input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// ReceiveMessageWithContext is a wrapper for the `sqs.SQS.ReceiveMessageWithContext`.
func (service *SQSService) ReceiveMessageWithContext(ctx context.Context, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().ReceiveMessageWithContext(ctx, input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// DeleteMessage is a wrapper for the `sqs.SQS.DeleteMessage`.
func (service *SQSService) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().DeleteMessage(input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// DeleteMessageWithContext is a wrapper for the `sqs.SQS.DeleteMessageWithContext`.
func (service *SQSService) DeleteMessageWithContext(ctx context.Context, input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().DeleteMessageWithContext(ctx, input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// DeleteMessageBatch is a wrapper for the `sqs.SQS.DeleteMessageBatch`.
func (service *SQSService) DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().DeleteMessageBatch(input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}

// DeleteMessageBatchWithContext is a wrapper for the `sqs.SQS.DeleteMessageBatchWithContext`.
func (service *SQSService) DeleteMessageBatchWithContext(ctx context.Context, input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	if service.isRunning() {
		input.QueueUrl = aws.String(service.Configuration.QUrl)
		return service.getSQS().DeleteMessageBatchWithContext(ctx, input)
	}
	return nil, rscsrv.ErrServiceNotRunning
}
