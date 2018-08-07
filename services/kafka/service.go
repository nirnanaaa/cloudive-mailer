package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/movio/kasper"
	"github.com/nirnanaaa/cloudive-mailer/services/kafka/event"
	"github.com/nirnanaaa/cloudive-mailer/services/smtp"
	"github.com/prometheus/client_golang/prometheus"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

// Service manages the listener and handler for an HTTP endpoint.
type Service struct {
	Logger            *logrus.Logger
	Config            *Config
	topicProcessor    *kasper.TopicProcessor
	MessageProcessors map[int]kasper.MessageProcessor
	KafkaClient       sarama.Client
	SMTP              *smtp.Service

	Producer sarama.AsyncProducer
}

// NewService returns a new instance of Service.
func NewService(c *Config) *Service {

	s := &Service{
		Config: c,
	}
	return s
}

// Connect estabiles a connection without Starting to work.
func (s *Service) Connect() error {

	client, err := sarama.NewClient(s.Config.Brokers, sarama.NewConfig())
	if err != nil {
		return err
	}
	s.KafkaClient = client
	if len(s.MessageProcessors) < 1 {
		return nil
	}
	config := kasper.Config{
		TopicProcessorName: s.Config.GroupName,
		Client:             client,
		InputTopics:        []string{s.Config.InboundQueueName},
		InputPartitions:    []int{0},
		Logger:             s.Logger.WithField("prefix", "kafka"),
	}

	s.topicProcessor = kasper.NewTopicProcessor(&config, s.MessageProcessors)
	return nil

}

// QueueMail queues a message into our internal kafka queue
func (s *Service) QueueMail(msg *event.InboundEmailEvent) error {
	s.Logger.Debugf("Delivering email with Trace ID %s", msg.TraceID)
	encoded, err := event.EncodeOutgoingEvent(msg)
	if err != nil {
		return err
	}
	key := uuid.NewV4().String()
	outgoingMessage := &sarama.ProducerMessage{
		Topic:     s.Config.OutboundQueueName,
		Partition: 0,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(encoded),
	}
	s.Producer.Input() <- outgoingMessage
	return nil
}

// ConnectProducer connects a kafka producer
func (s *Service) ConnectProducer() error {
	cConfig := s.KafkaClient.Config()
	cConfig.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	cConfig.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	cConfig.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	producer, err := sarama.NewAsyncProducer(s.Config.Brokers, cConfig)
	if err != nil {
		return err
	}
	go func() {
		for err := range producer.Errors() {
			s.Logger.WithError(err).Error("Producing a message failed")
			return
		}
	}()
	s.Producer = producer
	return nil
}

// Start starts the service
func (s *Service) Start() error {
	if err := s.Connect(); err != nil {
		return err
	}

	if err := s.ConnectProducer(); err != nil {
		return err
	}
	// metrics
	if err := s.registerMetrics(); err != nil {
		return err
	}
	if len(s.MessageProcessors) < 1 {
		return nil
	}
	return s.topicProcessor.RunLoop()
}

// SetDefaultMessageProcessor applies the default custom message processor
func (s *Service) SetDefaultMessageProcessor(smtp *smtp.Service) {
	processor := S3Processor{
		OutputTopicName: s.Config.OutboundQueueName,
		SMTP:            smtp,
	}
	processor.SetLogOutput(s.Logger)
	s.MessageProcessors = map[int]kasper.MessageProcessor{0: &processor}
}

// SetProcessor applies a custom message processor
func (s *Service) SetProcessor(processor kasper.MessageProcessor) {
	s.MessageProcessors = map[int]kasper.MessageProcessor{0: processor}
}

func (s *Service) registerMetrics() error {
	if err := prometheus.Register(processingTime); err != nil {
		return err
	}
	if err := prometheus.Register(totalProcessedCount); err != nil {
		return err
	}
	return prometheus.Register(errorCount)

}

// Stop closes the underlying listener.
func (s *Service) Stop() error {
	s.Producer.Close()
	if len(s.MessageProcessors) < 1 {
		return nil
	}
	s.topicProcessor.Close()
	return nil
}

// SetLogOutput sets the writer to which all logs are written. It must not be
// called after Open is called.
func (s *Service) SetLogOutput(log *logrus.Logger, prefix string) {
	s.Logger = log
}
