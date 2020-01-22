package sqssrv

import (
	"context"
	"log"
	"os"
	"path"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jamillosantos/macchiato"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

func TestServiceCollector(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)

	description := "SQS Service Collector Test Suite"
	if os.Getenv("CI") == "" {
		macchiato.RunSpecs(t, description)
	} else {
		reporterOutputDir := path.Join("./test-results/go-rscsrv-sqs")
		os.MkdirAll(reporterOutputDir, os.ModePerm)
		junitReporter := reporters.NewJUnitReporter(path.Join(reporterOutputDir, "results.xml"))
		macchiatoReporter := macchiato.NewReporter()
		RunSpecsWithCustomReporters(t, description, []Reporter{macchiatoReporter, junitReporter})
	}
}

var _ = Describe("SQSServiceCollector", func() {
	Context("testing prometheus metrics", func() {
		InitForTesting()

		When("using SendMessage", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing this body 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in size", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing message size 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficSize.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))

				expectedSize := 22 // len("testing message size x")

				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(expectedSize))
			})
		})

		When("using SendMessageWithContext", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
					MessageBody: aws.String("testing this body 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in size", func() {
				output, err := sqsService.SendMessageWithContext(context.Background(), &sqs.SendMessageInput{
					MessageBody: aws.String("testing message size 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficSize.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessage,
				}).Write(&metric)).To((Succeed()))

				expectedSize := 22 // len("testing message size x")

				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(expectedSize))
			})
		})

		When("using SendMessageBatch", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(2))
			})

			It("should increase traffic in size", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
					Entries: []*sqs.SendMessageBatchRequestEntry{
						{
							Id:          aws.String("message1"),
							MessageBody: aws.String("testing message size 1"),
						},
						{
							Id:          aws.String("message2"),
							MessageBody: aws.String("testing message size 2"),
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficSize.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))

				expectedSize := 44 // len("testing message size x") * 2

				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(expectedSize))
			})
		})

		When("using SendMessageBatchWithContext", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
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
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(2))
			})

			It("should increase traffic in size", func() {
				output, err := sqsService.SendMessageBatchWithContext(context.Background(), &sqs.SendMessageBatchInput{
					Entries: []*sqs.SendMessageBatchRequestEntry{
						{
							Id:          aws.String("message1"),
							MessageBody: aws.String("testing message size 1"),
						},
						{
							Id:          aws.String("message2"),
							MessageBody: aws.String("testing message size 2"),
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficSize.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodSendMessageBatch,
				}).Write(&metric)).To((Succeed()))

				expectedSize := 44 // len("testing message size x") * 2

				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(expectedSize))
			})
		})

		When("using ReceiveMessage", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				_, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())

				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err = sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))

			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing this body 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in size", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing message size 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficSize.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))

				expectedSize := 22 // len("testing message size x")

				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(expectedSize))
			})
		})

		When("using ReceiveMessageWithContext", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				_, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())

				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err = sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))

			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing this body 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in size", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing message size 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessageWithContext(context.Background(), &sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficSize.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodReceiveMessage,
				}).Write(&metric)).To((Succeed()))

				expectedSize := 22 // len("testing message size x")

				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(expectedSize))
			})
		})

		When("using DeleteMessage", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				_, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())

				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))

			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing this body 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessage(&sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})
		})

		When("using DeleteMessageWithContext", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).NotTo(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessageWithContext(context.Background(), &sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				_, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())

				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err = sqsService.DeleteMessageWithContext(context.Background(), &sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessageWithContext(context.Background(), &sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))

			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("this is the body of the message"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessageWithContext(context.Background(), &sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessage(&sqs.SendMessageInput{
					MessageBody: aws.String("testing this body 1"),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(aws.StringValue(output.MessageId)).ToNot(BeEmpty())

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds: aws.Int64(1),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(1))

				_, err = sqsService.DeleteMessageWithContext(context.Background(), &sqs.DeleteMessageInput{
					ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
				})
				Expect(err).ToNot(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessage,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})
		})

		When("using DeleteMessageBatch", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err = sqsService.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(2))
			})
		})

		When("using DeleteMessageBatchWithContext", func() {
			It("should increase duration", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatchWithContext(context.Background(), &sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageDuration.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeNumerically(">", 0))
			})

			It("should increase error amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				sqsService.Configuration.QUrl = "fake-url-to-return-error"
				_, err = sqsService.DeleteMessageBatchWithContext(context.Background(), &sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).To(HaveOccurred())

				var metric dto.Metric
				Expect(sqsService.Collector.messageFailures.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase success amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatchWithContext(context.Background(), &sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageSuccess.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase calls amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatchWithContext(context.Background(), &sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageCalls.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(1))
			})

			It("should increase traffic in amount", func() {
				output, err := sqsService.SendMessageBatch(&sqs.SendMessageBatchInput{
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
				Expect(err).ToNot(HaveOccurred())
				Expect(output.Successful).To(HaveLen(2))

				rcvOut, err := sqsService.ReceiveMessage(&sqs.ReceiveMessageInput{
					WaitTimeSeconds:     aws.Int64(1),
					MaxNumberOfMessages: aws.Int64(2),
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(rcvOut.Messages).To(HaveLen(2))

				delOut, err := sqsService.DeleteMessageBatchWithContext(context.Background(), &sqs.DeleteMessageBatchInput{
					Entries: []*sqs.DeleteMessageBatchRequestEntry{
						{
							Id:            rcvOut.Messages[0].MessageId,
							ReceiptHandle: rcvOut.Messages[0].ReceiptHandle,
						},
						{
							Id:            rcvOut.Messages[1].MessageId,
							ReceiptHandle: rcvOut.Messages[1].ReceiptHandle,
						},
					},
				})
				Expect(err).ToNot(HaveOccurred())
				Expect(delOut.Successful).To(HaveLen(2))

				var metric dto.Metric
				Expect(sqsService.Collector.messageTrafficAmount.With(prometheus.Labels{
					"queue":  *aws.String(sqsService.Configuration.QUrl),
					"method": MessageMetricMethodDeleteMessageBatch,
				}).Write(&metric)).To((Succeed()))
				Expect(metric.GetCounter().GetValue()).To(BeEquivalentTo(2))
			})
		})

	})

})
