package tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"ozon-omp-api/internal/model"
	"ozon-omp-api/internal/producer"
	"ozon-omp-api/internal/tests/mocks"
	"strconv"
	"sync"
	"testing"
)

func TestProducer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventsChan := make(chan model.CourseEvent, 3)
	sender := mocks.NewMockEventSender(ctrl)
	ctx := context.Background()
	wp := mocks.NewMockWorkerPool(ctrl)

	prodConfig := producer.Config{
		ParallelCount: 5,
		EventsChan:    eventsChan,
		Context:       ctx,
		Sender:        sender,
		WorkerPool:    wp,
	}
	p := producer.NewSenderProducer(prodConfig)

	var calls []string
	var callsSend []string

	var wg sync.WaitGroup
	wg.Add(6)

	wp.EXPECT().Remove(gomock.Any()).DoAndReturn(func(id uint64) error {
		defer wg.Done()
		calls = append(calls, strconv.FormatUint(id, 10))
		return nil
	}).AnyTimes()

	sender.EXPECT().Send(gomock.Any()).DoAndReturn(func(event *model.CourseEvent) error {
		defer wg.Done()
		callsSend = append(callsSend, strconv.FormatUint(event.ID, 10))
		return nil
	}).AnyTimes()

	eventsChan <- model.CourseEvent{
		ID:     1,
		Type:   model.Created,
		Status: model.Deferrer,
		Entity: nil,
	}
	eventsChan <- model.CourseEvent{
		ID:     2,
		Type:   model.Created,
		Status: model.Deferrer,
		Entity: nil,
	}
	eventsChan <- model.CourseEvent{
		ID:     3,
		Type:   model.Created,
		Status: model.Deferrer,
		Entity: nil,
	}

	p.Start()
	wg.Wait()
	p.Close()

	if len(calls) != 3 {
		t.Errorf("wrong number of calls: %v", calls)
	}
	if len(callsSend) != 3 {
		t.Errorf("wrong number of callsSend: %v", callsSend)
	}

}
