package sqssrv

import (
	"fmt"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

type SQSServiceCollector struct {
	messageCalls         *prometheus.CounterVec
	messageDuration      *prometheus.CounterVec
	messageSuccess       *prometheus.CounterVec
	messageFailures      *prometheus.CounterVec
	messageTrafficAmount *prometheus.CounterVec
	messageTrafficSize   *prometheus.CounterVec
}

type SQSServiceCollectorOpts struct {
	Prefix string
}

var (
	messageMetricVectorLabels = []string{"queue", "method"}
)

const (
	MessageMetricMethodSendMessage        string = "SendMessage"
	MessageMetricMethodSendMessageBatch   string = "SendMessageBatch"
	MessageMetricMethodDeleteMessage      string = "DeleteMessage"
	MessageMetricMethodDeleteMessageBatch string = "DeleteMessageBatch"
	MessageMetricMethodReceiveMessage     string = "ReceiveMessage"
)

func NewSQSServiceCollector(opts *SQSServiceCollectorOpts) *SQSServiceCollector {
	prefix := opts.Prefix
	if prefix != "" && !strings.HasSuffix(opts.Prefix, "_") {
		prefix += "_"
	}
	return &SQSServiceCollector{
		messageCalls: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("sqs_%smessage_calls", prefix),
				Help: "The total number of method called",
			},
			messageMetricVectorLabels,
		),
		messageDuration: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("sqs_%smessage_duration", prefix),
				Help: "The total duration (in seconds) of method called",
			},
			messageMetricVectorLabels,
		),
		messageSuccess: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("sqs_%smessage_success", prefix),
				Help: "The number of methods executed with success",
			},
			messageMetricVectorLabels,
		),
		messageFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("sqs_%smessage_failures", prefix),
				Help: "The number of methods executed with failures",
			},
			messageMetricVectorLabels,
		),
		messageTrafficAmount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("sqs_%smessage_traffic_amount", prefix),
				Help: "The total number of messages trafficked",
			},
			messageMetricVectorLabels,
		),
		messageTrafficSize: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: fmt.Sprintf("sqs_%smessage_traffic_size", prefix),
				Help: "The total size (number of characters) of messages trafficked",
			},
			messageMetricVectorLabels,
		),
	}
}

func (collector *SQSServiceCollector) Describe(descs chan<- *prometheus.Desc) {
	collector.messageCalls.Describe(descs)
	collector.messageDuration.Describe(descs)
	collector.messageSuccess.Describe(descs)
	collector.messageFailures.Describe(descs)
	collector.messageTrafficAmount.Describe(descs)
	collector.messageTrafficSize.Describe(descs)
}

func (collector *SQSServiceCollector) Collect(metrics chan<- prometheus.Metric) {
	collector.messageCalls.Collect(metrics)
	collector.messageDuration.Collect(metrics)
	collector.messageSuccess.Collect(metrics)
	collector.messageFailures.Collect(metrics)
	collector.messageTrafficAmount.Collect(metrics)
	collector.messageTrafficSize.Collect(metrics)
}
