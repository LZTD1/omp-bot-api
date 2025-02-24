package tests

import (
	"context"
	"github.com/golang/mock/gomock"
	"ozon-omp-api/internal/repo"
	"ozon-omp-api/internal/tests/mocks"
	"sync"
	"testing"
	"time"
)

func TestWorkpool(t *testing.T) {
	t.Run("Flush on timeout", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		repos := mocks.NewMockEventRepo(ctrl)

		var ids []uint64
		repos.EXPECT().Remove(gomock.Any()).DoAndReturn(func(eventIDs []uint64) error {
			for _, d := range eventIDs {
				ids = append(ids, d)
			}

			return nil
		}).AnyTimes()

		conf := repo.Config{
			Repo:      repos,
			PoolSize:  10,
			BatchSize: 2,
			Timer:     50 * time.Millisecond,
			Context:   ctx,
		}
		wp := repo.NewWorkerPool(conf)

		wp.Remove(1)
		time.Sleep(100 * time.Millisecond)
		if len(ids) != 1 {
			t.Errorf("ids length should be 1, but got %d", len(ids))
		}
	})

	t.Run("Flush on size", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		repos := mocks.NewMockEventRepo(ctrl)

		var ids []uint64
		var wg sync.WaitGroup
		wg.Add(3)

		repos.EXPECT().Remove(gomock.Any()).DoAndReturn(func(eventIDs []uint64) error {

			for _, d := range eventIDs {
				ids = append(ids, d)
			}

			return nil
		}).AnyTimes()

		conf := repo.Config{
			Repo:      repos,
			PoolSize:  10,
			BatchSize: 2,
			Timer:     50 * time.Minute,
			Context:   ctx,
		}
		wp := repo.NewWorkerPool(conf)

		wp.Remove(1)
		wp.Remove(2)
		wp.Remove(3)

		wp.Close()

		if len(ids) != 2 {
			t.Errorf("ids length should be 1, but got %d", len(ids))
		}
	})
}
