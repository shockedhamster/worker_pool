package workerpool

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/worker_pool/worker"
)

// Структура воркер-пула
type WorkerPool struct {
	taskCh  chan string
	workers []*worker.Worker
	ctx     context.Context
	cancel  context.CancelFunc
	mu      sync.Mutex
	handle  func(task string, id int)
}

// Создание воркер-пула
func NewWorkerPool(ctx context.Context, handle func(task string, id int)) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	return &WorkerPool{
		taskCh: make(chan string),
		ctx:    ctx,
		cancel: cancel,
		handle: handle,
	}
}

// Добавление воркеров в воркер-пул
func (wp *WorkerPool) AddWorkers(countWorkers int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	for i := 0; i < countWorkers; i++ {
		wp.addWorkerUnsafe()
	}
}

func (wp *WorkerPool) addWorkerUnsafe() {
	id := len(wp.workers) + 1

	worker := worker.NewWorker(wp.ctx, id, wp.taskCh, wp.handle)

	wp.workers = append(wp.workers, worker)
	go worker.Start()
	fmt.Printf("Worker %d added\n", id)
}

// Удаление воркеров из воркер-пула
func (wp *WorkerPool) RemoveWorkers(countWorkers int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for i := 0; i < countWorkers; i++ {
		wp.removeWorkerUnsafe()
	}
}

func (wp *WorkerPool) removeWorkerUnsafe() {
	if len(wp.workers) == 0 {
		fmt.Println("No workers to remove")
		return
	}

	worker := wp.workers[len(wp.workers)-1]
	worker.Stop()
	wp.workers = wp.workers[:len(wp.workers)-1]
	fmt.Printf("Worker %d removed\n", worker.Id)
}

// Отправка задачи в канал
func (wp *WorkerPool) SendTask(task string) error {
	select {
	case <-wp.ctx.Done():
		return errors.New("pool is shutting down, task not accepted")
	default:
		wp.taskCh <- task
	}

	return nil
}

// Завершение работы воркер-пула и всех воркеров
func (wp *WorkerPool) Shutdown() {
	wp.cancel()
	wp.mu.Lock()
	defer wp.mu.Unlock()
	for _, w := range wp.workers {
		w.Stop()
	}
	close(wp.taskCh)
	fmt.Println("All workers have been shut down")
}
