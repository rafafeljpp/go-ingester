package pool

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	IsOk   string = "OK"
	Locked string = "Locked"
)

// Manager of the pool
type Manager struct {
	workers             []*Worker
	maxWorkers          int
	maxElementsPerQueue int
	cIndex              int
	mutex               sync.Mutex
	rejected            chan Job
}

// Job type.
type Job struct {
	// Payload
	Payload string

	// Wait
	Wait time.Duration
}

// Worker
type Worker struct {
	id        int
	startedAt time.Time
	messages  chan Job
	signals   chan string
	status    string
}

// Listen worker.
func (w *Worker) Listen(wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		select {
		case msg := <-w.messages:
			fmt.Println("Mensaje recibido por el worker: ", w.id, msg.Payload)

			//time.Sleep(time.Second * msg.wait)
			break
		case <-w.signals:
			fmt.Println("finalizando el worker ", w.id)
			w.status = Locked
			close(w.messages)

			return
		}
	}
}

// GetMaxWorkwers var
func (p *Manager) GetMaxWorkers() int {
	return p.maxWorkers
}

// Start workers of the pool
func (p *Manager) Start(mw int, mepq int) {

	var w *Worker
	var wg sync.WaitGroup
	p.maxElementsPerQueue = mepq
	p.maxWorkers = mw
	p.cIndex = -1

	// Creando workers

	for i := 0; i < p.maxWorkers; i++ {
		w = p.createWorker(i)
		wg.Add(1)
		go w.Listen(&wg)
	}

	wg.Wait()
}

// Stop the pool
func (p *Manager) Stop() {

	for _, wk := range p.workers {
		wk.signals <- "stop"

	}
	for {
		if p.Length() == 0 {
			break
		}

	}
	return
}

// AddJob to de pool
func (p *Manager) AddJob(j Job) (Job, error) {
	var w *Worker

	// Selección del índice
	p.mutex.Lock()
	if p.cIndex+1 >= len(p.workers) {
		p.cIndex = 0
	} else {
		p.cIndex++
	}
	p.mutex.Unlock()

	w = p.workers[p.cIndex]

	if len(w.messages)+1 < p.maxElementsPerQueue {
		w.messages <- j
		return j, nil
	}

	return j, errors.New("Rechazado, Máximo número de elementos en cola o mensaje bloqueado")

}

// Length of the pool
func (p *Manager) Length() int {
	var AliveCounter int

	for _, w := range p.workers {
		if w.status == IsOk {
			AliveCounter++
		}
	}
	return AliveCounter
}

func (p *Manager) createWorker(id int) *Worker {

	w := new(Worker)

	w.id = id
	w.status = IsOk
	w.startedAt = time.Now()
	w.messages = make(chan Job)
	w.signals = make(chan string)

	p.workers = append(p.workers, w)

	return w
}
