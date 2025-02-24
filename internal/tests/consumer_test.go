package tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"ozon-omp-api/internal/consumer"
	"ozon-omp-api/internal/model"
	"ozon-omp-api/internal/tests/mocks"
	"reflect"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventsChan := make(chan model.CourseEvent)
	repos := mocks.NewMockEventRepo(ctrl)
	ctx := context.Background()

	consumerConfig := consumer.Config{
		ParallelCount:  5,
		EventsChan:     eventsChan,
		Repo:           repos,
		BatchSize:      5,
		UpdateInterval: time.Millisecond,
		Context:        ctx,
	}

	event := model.CourseEvent{
		ID:     1,
		Type:   model.Created,
		Status: model.Deferrer,
		Entity: &model.Course{
			ID:          1,
			Title:       "title",
			Description: "description",
		},
	}
	repos.EXPECT().Lock(uint64(5)).Return([]model.CourseEvent{event}, nil).AnyTimes()
	c := consumer.NewRepoConsumer(consumerConfig)
	c.Start()

	got := <-eventsChan
	if !reflect.DeepEqual(event, got) {
		t.Fatalf("Wanted %v, got %v", event, got)
	}
}
