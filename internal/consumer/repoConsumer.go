package consumer

import (
	"context"
	"fmt"
	"ozon-omp-api/internal/model"
	"ozon-omp-api/internal/repo"
	"sync"
	"time"
)

type Consumer interface {
	Start()
	Close()
}

type Config struct {
	ParallelCount  int
	EventsChan     chan<- model.CourseEvent
	Repo           repo.EventRepo
	BatchSize      int
	UpdateInterval time.Duration
	Context        context.Context
}

type consumer struct {
	n int
	e chan<- model.CourseEvent

	r repo.EventRepo

	b int
	i time.Duration

	d   chan bool
	wg  *sync.WaitGroup
	ctx context.Context
}

func NewRepoConsumer(config Config) Consumer {
	wg := new(sync.WaitGroup)
	d := make(chan bool)

	return &consumer{
		n:   config.ParallelCount,
		e:   config.EventsChan,
		r:   config.Repo,
		b:   config.BatchSize,
		i:   config.UpdateInterval,
		ctx: config.Context,
		d:   d,
		wg:  wg,
	}
}

func (c *consumer) Start() {
	for i := 0; i < c.n; i++ {
		c.wg.Add(1)

		go func() {
			defer c.wg.Done()
			ticker := time.NewTicker(c.i)

			for {
				select {
				case <-ticker.C:
					events, err := c.r.Lock(uint64(c.b))
					if err != nil {
						fmt.Printf("%v\n", fmt.Errorf("ошибка при получении записей : %w", err))
						continue
					}
					for _, event := range events {
						c.e <- event
					}
				case <-c.d:
					return
				case <-c.ctx.Done():
					return
				}
			}
		}()
	}
}

func (c *consumer) Close() {
	close(c.d)
	c.wg.Wait()
}
