package main

import (
	"context"
	"ozon-omp-api/internal/consumer"
	"ozon-omp-api/internal/model"
	"ozon-omp-api/internal/producer"
	"ozon-omp-api/internal/repo"
	"ozon-omp-api/internal/retranslator"
	"ozon-omp-api/internal/tests/mocks"
	"time"
)

func main() {
	eventsChan := make(chan model.CourseEvent)
	repos := mocks.MockEventRepo{}
	ctx := context.Background()
	sender := mocks.MockEventSender{}

	workerPool := repo.NewWorkerPool(repo.Config{
		Repo:      &repos,
		PoolSize:  2,
		BatchSize: 10,
		Timer:     5 * time.Second,
		Context:   ctx,
	})

	consumerConfig := consumer.Config{
		ParallelCount:  1,
		EventsChan:     eventsChan,
		Repo:           &repos,
		BatchSize:      5,
		UpdateInterval: 5 * time.Second,
		Context:        ctx,
	}
	producerConfig := producer.Config{
		ParallelCount: 1,
		EventsChan:    eventsChan,
		Context:       ctx,
		Sender:        &sender,
		WorkerPool:    workerPool,
	}

	r := retranslator.NewRetranslator(consumerConfig, producerConfig)
	r.Start()
	i := 0
	for {
		i += 1
	}
}
