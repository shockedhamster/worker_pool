package worker

import (
	"context"
	"fmt"
)

// Структура воркера
type Worker struct {
	Id     int
	taskCh <-chan string
	ctx    context.Context
	handle func(task string, id int)
	cancel context.CancelFunc

	stoppedCh chan struct{}
	stopped   bool
}

// Создание воркера
func NewWorker(ctx context.Context, id int, taskCh <-chan string, handle func(task string, id int)) *Worker {
	workerCtx, workerCancel := context.WithCancel(ctx)
	return &Worker{
		Id:        id,
		taskCh:    taskCh,
		ctx:       workerCtx,
		handle:    handle,
		cancel:    workerCancel,
		stoppedCh: make(chan struct{}, 1),
	}
}

// Запуск воркера
func (w *Worker) Start() {
LOOP:
	for {
		select {
		case task, ok := <-w.taskCh:
			if !ok {
				break LOOP
			}
			w.handle(task, w.Id)

		case <-w.ctx.Done():
			break LOOP
		}
	}

	w.stopped = true
	w.stoppedCh <- struct{}{}
	fmt.Printf("worker %d stopped\n", w.Id)
}

// Остановка воркера
func (w *Worker) Stop() {
	if w.stopped {
		return
	}
	w.cancel()
	<-w.stoppedCh
}
