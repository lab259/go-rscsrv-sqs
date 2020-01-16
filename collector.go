package sqssrv

import "github.com/prometheus/client_golang/prometheus"

type SQSServiceCollector struct {
	// prometheus counters
	sendMessageCalls     *prometheus.CounterVec
	sendMessageDuration  *prometheus.CounterVec
	sendMessageSuccess   *prometheus.CounterVec
	sendMessageFailures  *prometheus.CounterVec
	messageTrafficAmount *prometheus.CounterVec
	messageTrafficSize   *prometheus.CounterVec
}

var sendMessageVectorLabels = []string{"queue", "method"}

func NewSQSServiceCollector() *SQSServiceCollector {
	return &SQSServiceCollector{
		sendMessageCalls: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			sendMessageVectorLabels,
		),
		sendMessageDuration: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			sendMessageVectorLabels,
		),
		sendMessageSuccess: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			sendMessageVectorLabels,
		),
		sendMessageFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			sendMessageVectorLabels,
		),
		messageTrafficAmount: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			[]string{"queue", "direction"},
		),
		messageTrafficSize: prometheus.NewCounterVec(
			prometheus.CounterOpts{},
			[]string{"queue", "direction"},
		),
	}
}

func (collector *SQSServiceCollector) Describe(descs chan<- *prometheus.Desc) {
	// TODO: Add describe
}

func (collector *SQSServiceCollector) Collect(metrics chan<- prometheus.Metric) {
	collector.sendMessageDuration.Collect(metrics)
}
