package kafka

import "github.com/prometheus/client_golang/prometheus"

var (
	totalProcessedCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kafka_filter_processed",
		Help: "Number of processed messages",
	})
	processingTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "kafka_filter_duration",
		Help: "Time it takes for messages to be filtered",
	})
	errorCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "kafka_filter_error",
		Help: "Number of unsuccessfully processed messages",
	})
)
