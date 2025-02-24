package repo

import (
	"context"
	"fmt"
	"github.com/gammazero/workerpool"
	"sync"
	"time"
)

type WorkerPool interface {
	Remove(ids uint64)
	Close()
}

type workerPool struct {
	repo  EventRepo
	wp    *workerpool.WorkerPool
	mu    sync.Mutex
	batch []uint64
	bs    int
	timer *time.Timer
	time  time.Duration
	ctx   context.Context
}

type Config struct {
	Repo      EventRepo
	PoolSize  int
	BatchSize int
	Timer     time.Duration
	Context   context.Context
}

// Что бы при случайном переполнении бача горутинами, не аллоцировался новый слайс
const offsetBatch = 10

func NewWorkerPool(cfg Config) WorkerPool {
	wp := workerpool.New(cfg.PoolSize)

	pool := &workerPool{
		repo:  cfg.Repo,
		wp:    wp,
		ctx:   cfg.Context,
		time:  cfg.Timer,
		timer: time.NewTimer(cfg.Timer),
		batch: make([]uint64, 0, cfg.BatchSize+offsetBatch),
		bs:    cfg.BatchSize,
	}

	go pool.flushOnTimeout()

	return pool
}

func (w *workerPool) Remove(id uint64) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.batch = append(w.batch, id)

	if len(w.batch) >= w.bs {
		w.flush()
	}
}

func (w *workerPool) flushOnTimeout() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for {
		select {
		case <-w.timer.C:
			w.flush()
		case <-w.ctx.Done():
			return
		}
	}
}

func (w *workerPool) flush() {
	if len(w.batch) == 0 {
		return
	}

	batchCopy := append([]uint64(nil), w.batch...)
	w.wp.Submit(func() {
		err := w.repo.Remove(batchCopy)

		if err != nil {
			fmt.Printf("repo remove err:%v\n", err)
		}
	})
	w.wp.StopWait()
	w.batch = w.batch[:0:0]
	if !w.timer.Stop() {
		<-w.timer.C
	}
	w.timer.Reset(w.time)
}

func (w *workerPool) Close() {
	w.wp.StopWait()
}
