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

type Pool struct {
	workers             []*Worker
	maxWorkers          int
	maxElementsPerQueue int
	cIndex              int
	mutex               sync.Mutex
	rejected            chan Job
}

// Start workers of the pool
func Start(p *Pool) {
	var w *Worker
	var wg sync.WaitGroup

	p.cIndex = -1

	// Creando workers

	for i := 0; i < p.maxWorkers; i++ {
		w = createWorker(p, i)
		wg.Add(1)
		go Listen(w, &wg)
	}

	wg.Wait()
}

// Stop the pool
func Stop(p *Pool) {

	for _, wk := range p.workers {
		wk.signals <- "stop"
	}
	for {
		if Length(p) == 0 {
			break
		}

	}
	return
}

// Length of the pool
func Length(p *Pool) int {
	var AliveCounter int

	for _, worker := range p.workers {
		if worker.status == IsOk {
			AliveCounter++
		}
	}
	return AliveCounter
}

func createWorker(p *Pool, id int) *Worker {

	w := new(Worker)
	w.id = id
	w.status = IsOk
	w.startedAt = time.Now()
	w.messages = make(chan Job)
	w.signals = make(chan string)

	p.workers = append(p.workers, w)

	return w
}

// AddJob to de pool
func addJob(p *Pool, j Job) (Job, error) {
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

func init() {
	fmt.Println("Pool package inicializado")
}
