package sqssrv

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jamillosantos/macchiato"
	rscsrv "github.com/lab259/go-rscsrv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	macchiato.RunSpecs(t, "SQS Service Test Suite")
}

var _ = Describe("SQSService", func() {
	It("should fail loading a configuration", func() {
		var service SQSService
		configuration, err := service.LoadConfiguration()
		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("not implemented"))
		Expect(configuration).To(BeNil())
	})

	It("should fail applying configuration", func() {
		var service SQSService
		err := service.ApplyConfiguration(map[string]interface{}{
			"address": "localhost",
		})
		Expect(err).To(Equal(rscsrv.ErrWrongConfigurationInformed))
	})

	It("should apply the configuration using a pointer", func() {
		var service SQSService
		err := service.ApplyConfiguration(&SQSServiceConfiguration{
			Region:   "region",
			Endpoint: "endpoint",
			Secret:   "secret",
			QUrl:     "qurl",
			Key:      "key",
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Region).To(Equal("region"))
		Expect(service.Configuration.Endpoint).To(Equal("endpoint"))
		Expect(service.Configuration.Secret).To(Equal("secret"))
		Expect(service.Configuration.QUrl).To(Equal("qurl"))
		Expect(service.Configuration.Key).To(Equal("key"))
	})

	It("should apply the configuration using a copy", func() {
		var service SQSService
		err := service.ApplyConfiguration(SQSServiceConfiguration{
			Region:   "region",
			Endpoint: "endpoint",
			Secret:   "secret",
			QUrl:     "qurl",
			Key:      "key",
		})
		Expect(err).To(BeNil())
		Expect(service.Configuration.Region).To(Equal("region"))
		Expect(service.Configuration.Endpoint).To(Equal("endpoint"))
		Expect(service.Configuration.Secret).To(Equal("secret"))
		Expect(service.Configuration.QUrl).To(Equal("qurl"))
		Expect(service.Configuration.Key).To(Equal("key"))
	})

	validConfiguration := SQSServiceConfiguration{
		Endpoint: "http://localhost:9324",
		QUrl:     "http://localhost:9324/queue/queue-test",
	}

	It("should start the service", func() {
		var service SQSService
		Expect(service.ApplyConfiguration(&validConfiguration)).To(BeNil())
		Expect(service.Start()).To(BeNil())
		defer service.Stop()
		output, err := service.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String("this is the body of the message"),
		})
		Expect(err).To(BeNil())
		Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())
	})

	It("should start the service with a different host", func() {
		var service SQSService
		v := validConfiguration
		v.QUrl = "http://differenthost:9324/queue/queue-test"
		Expect(service.ApplyConfiguration(&v)).To(BeNil())
		Expect(service.Start()).To(BeNil())
		defer service.Stop()
		output, err := service.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String("this is the body of the message"),
		})
		Expect(err).To(BeNil())
		Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())
	})

	It("should fail starting the service with a non existing queue", func() {
		var service SQSService
		v := validConfiguration
		v.QUrl += "-nonexistent"
		Expect(service.ApplyConfiguration(&v)).To(BeNil())
		err := service.Start()
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("queue"))
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("should stop the service", func() {
		var service SQSService
		Expect(service.ApplyConfiguration(&validConfiguration)).To(BeNil())
		Expect(service.Start()).To(BeNil())
		Expect(service.Stop()).To(BeNil())
		Expect(service.RunWithSQS(func(client *sqs.SQS) error {
			return nil
		})).To(Equal(rscsrv.ErrServiceNotRunning))
	})

	It("should restart the service", func() {
		var service SQSService
		Expect(service.ApplyConfiguration(&validConfiguration)).To(BeNil())
		Expect(service.Start()).To(BeNil())
		Expect(service.Restart()).To(BeNil())
		output, err := service.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String("this is the body of the message"),
		})
		Expect(err).To(BeNil())
		Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())
	})

	When("not running the service", func() {
		sqsService := &SQSService{}

		It("should fail sending a message", func() {
			_, err := sqsService.SendMessage(&sqs.SendMessageInput{
				MessageBody: aws.String("testing data"),
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail sending a message with context", func() {
			_, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
				MessageBody: aws.String("testing data"),
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail sending a message in batch", func() {
			_, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
				Entries: []*sqs.SendMessageBatchRequestEntry{
					{
						MessageBody: aws.String("this is a message 1"),
					},
					{
						MessageBody: aws.String("this is a message 2"),
					},
				},
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail sending a message in batch with context", func() {
			_, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
				Entries: []*sqs.SendMessageBatchRequestEntry{
					{
						MessageBody: aws.String("this is a message 1"),
					},
					{
						MessageBody: aws.String("this is a message 2"),
					},
				},
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail receiving a message", func() {
			_, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(1),
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail receiving a message with context", func() {
			_, err := sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(1),
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail deleting a message", func() {
			_, err := sqsService.DeleteMessage(&sqs.DeleteMessageInput{
				ReceiptHandle: aws.String("fake message"),
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail deleting a message with context", func() {
			_, err := sqsService.DeleteMessageWithContext(context.Background(), &sqs.DeleteMessageInput{
				ReceiptHandle: aws.String("fake message"),
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail deleting a message batch", func() {
			_, err := sqsService.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
				Entries: []*sqs.DeleteMessageBatchRequestEntry{
					{
						ReceiptHandle: aws.String("fake message 1"),
					},
					{
						ReceiptHandle: aws.String("fake message 2"),
					},
				},
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})

		It("should fail deleting a message batch with context", func() {
			_, err := sqsService.DeleteMessageBatchWithContext(context.Background(), &sqs.DeleteMessageBatchInput{
				Entries: []*sqs.DeleteMessageBatchRequestEntry{
					{
						ReceiptHandle: aws.String("fake message 1"),
					},
					{
						ReceiptHandle: aws.String("fake message 2"),
					},
				},
			})
			Expect(err).To(Equal(rscsrv.ErrServiceNotRunning))
		})
	})

	Context("sending and receiving messages", func() {
		var sqsService *SQSService

		BeforeEach(func() {
			sqsService = &SQSService{}
			Expect(sqsService.ApplyConfiguration(&validConfiguration)).To(BeNil())
			Expect(sqsService.Start()).To(BeNil())
			for {
				messages, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					MaxNumberOfMessages: aws.Int64(1),
					WaitTimeSeconds:     aws.Int64(1),
				})
				Expect(err).To(BeNil())
				if len(messages.Messages) == 0 {
					break
				}
				_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: messages.Messages[0].ReceiptHandle,
				})
				Expect(err).To(BeNil())
			}
			time.Sleep(time.Millisecond * 100)
		})

		AfterEach(func() {
			Expect(sqsService.Stop()).To(BeNil())
		})

		It("should send and receive a message", func() {
			sendOut, err := sqsService.SendMessage(&sqs.SendMessageInput{
				MessageBody: aws.String("testing this body"),
			})
			Expect(err).To(BeNil())
			rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(1),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(1))
			Expect(aws.StringValue(sendOut.MessageId)).To(Equal(aws.StringValue(rcvOut.Messages[0].MessageId)))
			Expect(aws.StringValue(rcvOut.Messages[0].Body)).To(Equal("testing this body"))
		})

		It("should send and receive a message with context", func() {
			sendOut, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
				MessageBody: aws.String("testing this body"),
			})
			Expect(err).To(BeNil())
			rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(1),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(1))
			Expect(aws.StringValue(sendOut.MessageId)).To(Equal(aws.StringValue(rcvOut.Messages[0].MessageId)))
			Expect(aws.StringValue(rcvOut.Messages[0].Body)).To(Equal("testing this body"))
		})

		It("should delete a message", func() {
			_, err := sqsService.SendMessage(&sqs.SendMessageInput{
				MessageBody: aws.String("testing this body"),
			})
			Expect(err).To(BeNil())
			rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(1),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(1))
			_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
				ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
			})
			Expect(err).To(BeNil())
			rcvOut, err = sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(0),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(BeEmpty())
		})

		It("should delete a message with context", func() {
			_, err := sqsService.SendMessage(&sqs.SendMessageInput{
				MessageBody: aws.String("testing this body"),
			})
			Expect(err).To(BeNil())
			rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(1),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(1))
			_, err = sqsService.DeleteMessageWithContext(context.Background(), &sqs.DeleteMessageInput{
				ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
			})
			Expect(err).To(BeNil())
			rcvOut, err = sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds: aws.Int64(0),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(BeEmpty())
		})

		It("should send a message batch", func() {
			sendOut, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
				Entries: []*sqs.SendMessageBatchRequestEntry{
					{
						Id:          aws.String("message1"),
						MessageBody: aws.String("testing this body 1"),
					},
					{
						Id:          aws.String("message2"),
						MessageBody: aws.String("testing this body 2"),
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(sendOut.Successful).To(HaveLen(2))
			rcvOut, err := sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
				WaitTimeSeconds:     aws.Int64(1),
				MaxNumberOfMessages: aws.Int64(2),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(2))
			Expect([]string{aws.StringValue(rcvOut.Messages[0].MessageId), aws.StringValue(rcvOut.Messages[1].MessageId)}).To(ConsistOf(aws.StringValue(sendOut.Successful[1].MessageId), aws.StringValue(sendOut.Successful[0].MessageId)))
		})

		It("should send a message batch with context", func() {
			sendOut, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
				Entries: []*sqs.SendMessageBatchRequestEntry{
					{
						Id:          aws.String("message1"),
						MessageBody: aws.String("testing this body 1"),
					},
					{
						Id:          aws.String("message2"),
						MessageBody: aws.String("testing this body 2"),
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(sendOut.Successful).To(HaveLen(2))
			rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds:     aws.Int64(1),
				MaxNumberOfMessages: aws.Int64(2),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(2))
			Expect([]string{aws.StringValue(rcvOut.Messages[0].MessageId), aws.StringValue(rcvOut.Messages[1].MessageId)}).To(ConsistOf(aws.StringValue(sendOut.Successful[1].MessageId), aws.StringValue(sendOut.Successful[0].MessageId)))
		})

		It("should delete a message in batch", func() {
			sendOut, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
				Entries: []*sqs.SendMessageBatchRequestEntry{
					{
						Id:          aws.String("message1"),
						MessageBody: aws.String("testing this body 1"),
					},
					{
						Id:          aws.String("message2"),
						MessageBody: aws.String("testing this body 2"),
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(sendOut.Successful).To(HaveLen(2))
			rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds:     aws.Int64(1),
				MaxNumberOfMessages: aws.Int64(2),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(2))
			Expect([]string{aws.StringValue(rcvOut.Messages[0].MessageId), aws.StringValue(rcvOut.Messages[1].MessageId)}).To(ConsistOf(aws.StringValue(sendOut.Successful[1].MessageId), aws.StringValue(sendOut.Successful[0].MessageId)))

			sqsService.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
				Entries: []*sqs.DeleteMessageBatchRequestEntry{
					{
						ReceiptHandle: sendOut.Successful[0].MessageId,
					},
					{
						ReceiptHandle: sendOut.Successful[1].MessageId,
					},
				},
			})

			rcvOut, err = sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds:     aws.Int64(1),
				MaxNumberOfMessages: aws.Int64(2),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(BeEmpty())
		})

		It("should delete a message in batch with context", func() {
			sendOut, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
				Entries: []*sqs.SendMessageBatchRequestEntry{
					{
						Id:          aws.String("message1"),
						MessageBody: aws.String("testing this body 1"),
					},
					{
						Id:          aws.String("message2"),
						MessageBody: aws.String("testing this body 2"),
					},
				},
			})
			Expect(err).To(BeNil())
			Expect(sendOut.Successful).To(HaveLen(2))
			rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds:     aws.Int64(1),
				MaxNumberOfMessages: aws.Int64(2),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(HaveLen(2))
			Expect([]string{aws.StringValue(rcvOut.Messages[0].MessageId), aws.StringValue(rcvOut.Messages[1].MessageId)}).To(ConsistOf(aws.StringValue(sendOut.Successful[1].MessageId), aws.StringValue(sendOut.Successful[0].MessageId)))

			sqsService.DeleteMessageBatchWithContext(context.Background(), &sqs.DeleteMessageBatchInput{
				Entries: []*sqs.DeleteMessageBatchRequestEntry{
					{
						ReceiptHandle: sendOut.Successful[0].MessageId,
					},
					{
						ReceiptHandle: sendOut.Successful[1].MessageId,
					},
				},
			})

			rcvOut, err = sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
				WaitTimeSeconds:     aws.Int64(1),
				MaxNumberOfMessages: aws.Int64(2),
			})
			Expect(err).To(BeNil())
			Expect(rcvOut.Messages).To(BeEmpty())
		})
	})
})
