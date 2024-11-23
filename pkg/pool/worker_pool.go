package pool

import (
    "fmt"
    "sync"
)

var GlobalWorkerPool *WorkerPool

type Task struct {
    Execute func() error
}

type WorkerPool struct {
    workerCount int
    taskChan    chan Task
    wg          sync.WaitGroup
    isRunning   bool
    mu          sync.Mutex
}

func NewWorkerPool(workerCount int) *WorkerPool {
    return &WorkerPool{
        workerCount: workerCount,
        taskChan:    make(chan Task, workerCount*2),
        isRunning:   false,
    }
}

func (p *WorkerPool) Start() {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.isRunning {
        return
    }

    p.isRunning = true
    p.wg.Add(p.workerCount)
    
    for i := 0; i < p.workerCount; i++ {
        go p.worker()
    }
}

func (p *WorkerPool) Stop() {
    p.mu.Lock()
    if !p.isRunning {
        p.mu.Unlock()
        return
    }
    p.isRunning = false
    p.mu.Unlock()

    close(p.taskChan)
    p.wg.Wait()
}

func (p *WorkerPool) worker() {
    defer p.wg.Done()
    
    for task := range p.taskChan {
        if err := task.Execute(); err != nil {
            fmt.Printf("Task execution error: %v\n", err)
        }
    }
}

func (p *WorkerPool) Submit(task Task) error {
    p.mu.Lock()
    if !p.isRunning {
        p.mu.Unlock()
        return fmt.Errorf("worker pool is not running")
    }
    p.mu.Unlock()

    p.taskChan <- task
    return nil
}