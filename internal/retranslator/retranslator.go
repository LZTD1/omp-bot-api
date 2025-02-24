package retranslator

import (
	"ozon-omp-api/internal/consumer"
	"ozon-omp-api/internal/model"
	"ozon-omp-api/internal/producer"
)

type Retranslator interface {
	Start()
	Close()
}

type retranslator struct {
	events   chan model.CourseEvent
	consumer consumer.Consumer
	producer producer.Producer
}

func NewRetranslator(consumerConfig consumer.Config, producerConfig producer.Config) Retranslator {
	c := consumer.NewRepoConsumer(consumerConfig)
	p := producer.NewSenderProducer(producerConfig)

	events := make(chan model.CourseEvent)

	rts := retranslator{
		events:   events,
		consumer: c,
		producer: p,
	}
	return rts
}

func (r retranslator) Start() {
	r.producer.Start()
	r.consumer.Start()
}

func (r retranslator) Close() {
	r.producer.Close()
	r.consumer.Close()
}
