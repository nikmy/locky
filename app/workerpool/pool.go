package workerpool

import (
	"context"

	"github.com/sourcegraph/conc/pool"
)

func New[T any](nWorkers int) *WorkerPool[T] {
	return &WorkerPool[T]{
		workers: nWorkers,
	}
}

type WorkerPool[T any] struct {
	context context.Context
	handler func(T)
	workers int
}

func (wp *WorkerPool[T]) WithContext(ctx context.Context) *WorkerPool[T] {
	wp.context = ctx
	return wp
}

func (wp *WorkerPool[T]) WithHandler(h func(T)) *WorkerPool[T] {
	wp.handler = h
	return wp
}

func (wp *WorkerPool[T]) Range(tasks <-chan T) {
	p := pool.New().WithMaxGoroutines(8)
	for i := 0; i < wp.workers; i++ {
		p.Go(func() {
			for {
				select {
				case t := <-tasks:
					wp.handler(t)
				case <-wp.context.Done():
					return
				}
			}
		})
	}
	p.Wait()
}
