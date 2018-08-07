package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/movio/kasper"
	"github.com/nirnanaaa/cloudive-mailer/services/kafka/event"
	"github.com/nirnanaaa/cloudive-mailer/services/smtp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var errRecordsEmpty = fmt.Errorf("records array was empty so we coudln't continue with processing")

// S3Processor is message processor that enriches messages from s3 with info from the head request metadata
type S3Processor struct {
	OutputTopicName string
	Logger          *logrus.Entry
	SMTP            *smtp.Service
}

// SetLogOutput sets a new log output for this module
func (processor *S3Processor) SetLogOutput(log *logrus.Logger) {
	processor.Logger = log.WithField("prefix", "s3Processor")
}

// Process starts processing messages
func (processor *S3Processor) Process(msgs []*sarama.ConsumerMessage, sender kasper.Sender) error {
	logger := processor.Logger
	logger.Debugf("Started processing a batch of %d messages", len(msgs))
	for idx, msg := range msgs {
		totalProcessedCount.Inc()
		logger.Debugf("[%d/%d] Processing started", idx+1, len(msgs))
		timer := prometheus.NewTimer(processingTime)
		defer timer.ObserveDuration()
		if err := processor.ProcessMessage(msg, sender); err != nil {
			logger.Errorf("[%d/%d] Processing errored: %s", idx+1, len(msgs), err.Error())
			errorCount.Inc()
		}
		logger.Debugf("[%d/%d] Processing done", idx+1, len(msgs))

	}
	return nil
}

type MailMessageMetadata struct {
	Tries int `json:"tries"`
}

// ProcessMessage processes an incomming message
func (processor *S3Processor) ProcessMessage(msg *sarama.ConsumerMessage, sender kasper.Sender) error {
	// l := processor.Logger
	var decoded event.InboundEmailEvent
	if err := json.Unmarshal(msg.Value, &decoded); err != nil {
		return err
	}

	// re-queue
	if err := processor.SMTP.Deliver(&decoded); err != nil {
		outgoingMessage := &sarama.ProducerMessage{
			Topic:     processor.OutputTopicName,
			Partition: 0,
			Key:       sarama.ByteEncoder(msg.Key),
			Value:     sarama.ByteEncoder(msg.Value),
		}
		sender.Send(outgoingMessage)
		return err
	}
	// processor.SMTP.
	return nil
}
