package workerpool

import (
	"context"

	"github.com/sourcegraph/conc/pool"
)

func New[T any]() *WorkerPool[T] {
	return &WorkerPool[T]{
		context: context.Background(),
		workers: pool.New(),
	}
}

type WorkerPool[T any] struct {
	context context.Context
	handler func(T)
	workers *pool.Pool
}

func (wp *WorkerPool[T]) WithMaxGoroutines(n int) *WorkerPool[T] {
	wp.workers.WithMaxGoroutines(n)
	return wp
}

func (wp *WorkerPool[T]) WithContext(ctx context.Context) *WorkerPool[T] {
	wp.context = ctx
	return wp
}

func (wp *WorkerPool[T]) WithHandler(h func(T)) *WorkerPool[T] {
	wp.handler = h
	return wp
}

func (wp *WorkerPool[T]) Range(tasks <-chan T) *WorkerPool[T] {
	for i := 0; i < wp.workers.MaxGoroutines(); i++ {
		wp.workers.Go(func() {
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
	return wp
}

func (wp *WorkerPool[T]) Await() <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		wp.workers.Wait()
		ch <- struct{}{}
	}()
	return ch
}
