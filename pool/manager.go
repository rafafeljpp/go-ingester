// Package pool implements a pool of workers funcionality
// Author: rafael.pellicer@gmail.com
package pool

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	isOk   string = "OK"
	locked string = "Locked"
)

// Manager of the pool
type Manager struct {
	workers             []*Worker
	maxWorkers          int
	maxElementsPerQueue int
	cIndex              int
	mutex               sync.Mutex
	rejected            chan iJob
}

// iJob Interface Contains methods for serialize and publish data.
type iJob interface {
	GetPayload() string
	Serialize() bool
	Publish() bool
}

// Worker
type Worker struct {
	id        int
	startedAt time.Time
	messages  chan iJob
	signals   chan string
	status    string
}

// Listen worker.
func (w *Worker) Listen(wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		select {
		case msg := <-w.messages:
			if msg.Serialize() {
				msg.Publish()
			}

			//time.Sleep(time.Second * msg.wait)
			break
		case <-w.signals:
			fmt.Println("finalizando el worker ", w.id)
			w.status = locked
			close(w.messages)

			return
		}
	}
}

// GetMaxWorkwers returns the maximum number of workers in the pool
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

// AddJob to the pool
func (p *Manager) AddJob(j iJob) (iJob, error) {
	var w *Worker

	// FIX:
	// Falla cuando hay mucha concurrencia.
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
		if w.status == isOk {
			AliveCounter++
		}
	}
	return AliveCounter
}

// createWorker Create a worker in the pool.
func (p *Manager) createWorker(id int) *Worker {

	w := new(Worker)

	w.id = id
	w.status = isOk
	w.startedAt = time.Now()
	w.messages = make(chan iJob)
	w.signals = make(chan string)

	p.workers = append(p.workers, w)

	return w
}
