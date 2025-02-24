package producer

import (
	"context"
	"fmt"
	"ozon-omp-api/internal/model"
	"ozon-omp-api/internal/repo"
	"ozon-omp-api/internal/sender"
	"sync"
)

type Producer interface {
	Start()
	Close()
}

type Config struct {
	ParallelCount int
	EventsChan    <-chan model.CourseEvent
	Context       context.Context
	Sender        sender.EventSender
	WorkerPool    repo.WorkerPool
}

type eventProducer struct {
	n int
	e <-chan model.CourseEvent

	s sender.EventSender

	d  chan bool
	wg *sync.WaitGroup

	ctx context.Context
	wp  repo.WorkerPool
}

func NewSenderProducer(cfg Config) Producer {
	wg := new(sync.WaitGroup)
	done := make(chan bool)

	return &eventProducer{
		n:   cfg.ParallelCount,
		e:   cfg.EventsChan,
		s:   cfg.Sender,
		ctx: cfg.Context,
		wp:  cfg.WorkerPool,
		d:   done,
		wg:  wg,
	}
}

func (s *eventProducer) Start() {
	for i := 0; i < s.n; i++ {
		s.wg.Add(1)

		go func() {
			defer s.wg.Done()

			for {
				select {
				case e := <-s.e:
					if err := s.s.Send(&e); err != nil {
						fmt.Printf("%v\n", fmt.Errorf("ошибка при отправке записи : %w", err))
						continue
					}
					s.wp.Remove(e.ID)
				case <-s.d:
					return
				case <-s.ctx.Done():
					return
				}
			}
		}()
	}
}

func (s *eventProducer) Close() {
	close(s.d)
	s.wg.Wait()
}
