package sqssrv

import "github.com/prometheus/client_golang/prometheus"

type SQSServiceCollector struct {

	// prometheus counters
	messageCalls         *prometheus.CounterVec
	messageDuration      *prometheus.CounterVec
	messageSuccess       *prometheus.CounterVec
	messageFailures      *prometheus.CounterVec
	messageTrafficAmount *prometheus.CounterVec
	messageTrafficSize   *prometheus.CounterVec
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

func NewSQSServiceCollector() *SQSServiceCollector {
	return &SQSServiceCollector{
		messageCalls: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			messageMetricVectorLabels,
		),
		messageDuration: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			messageMetricVectorLabels,
		),
		messageSuccess: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			messageMetricVectorLabels,
		),
		messageFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			messageMetricVectorLabels,
		),
		messageTrafficAmount: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			messageMetricVectorLabels,
		),
		messageTrafficSize: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			messageMetricVectorLabels,
		),
	}
}

func (collector *SQSServiceCollector) Describe(descs chan<- *prometheus.Desc) {
	// TODO: Add describe
}

func (collector *SQSServiceCollector) Collect(metrics chan<- prometheus.Metric) {
	collector.messageDuration.Collect(metrics)
}
