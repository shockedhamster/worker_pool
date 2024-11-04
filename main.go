package main

import (
	"context"
	"fmt"
	"time"

	workerpool "github.com/worker_pool/workerPool"
)

const (
	initialWorkers       = 3                // Начальное количество воркеров
	taskInterval         = 1 * time.Second  // Интервал, с которым поступают новые задачи
	addWorkerInterval    = 4 * time.Second  // Интервал, с которым добавляются воркеры
	removeWorkerInterval = 9 * time.Second  // Интервал, с которым воркеры удаляются
	workDuration         = 20 * time.Second // Время работы программы
)

func main() {
	ctx, _ := context.WithCancel(context.Background())

	pool := workerpool.NewWorkerPool(ctx, func(task string, id int) {
		fmt.Printf("Worker %d processing: %s\n", id, task)
	})

	pool.AddWorkers(initialWorkers)

	// Добавляем воркеров в фоне
	go backAddWorkers(pool)

	// Удаляем воркеров в фоне
	go backRemoveWorkers(pool)

	// Отправляем задачи в фоне
	go backSendTask(pool)

	// Время работы программы
	time.Sleep(workDuration)

	pool.Shutdown()
}

// Генерация задач с временным промежутком
func backSendTask(pool *workerpool.WorkerPool) {
	i := 0
	for {
		i++
		pool.SendTask(fmt.Sprintf("Task %d", i))
		time.Sleep(taskInterval)
	}
}

// Добавление воркеров с временным промежутком
func backAddWorkers(pool *workerpool.WorkerPool) {
	for {
		for {
			time.Sleep(addWorkerInterval)
			pool.AddWorkers(1)
		}
	}
}

// Удаление воркеров с временным промежутком
func backRemoveWorkers(pool *workerpool.WorkerPool) {
	for {
		for {
			time.Sleep(removeWorkerInterval)
			pool.RemoveWorkers(1)
		}
	}
}
