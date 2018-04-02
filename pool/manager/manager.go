package pool

import (
	"fmt"
	"sync"

	. "go-ingester/pool/job"
	. "go-ingester/pool/worker"
)

type Pool struct {
	workers             []*Worker
	maxWorkers          int
	maxElementsPerQueue int
	cIndex              int
	mutex               sync.Mutex
	rejected            chan Job
}

func start() {
	fmt.Println("Hola")
}

/*
// Start workers of the pool
func (p *Pool) Start() {

	var w *Worker
	var wg sync.WaitGroup

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
func (p *Pool) Stop() {

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

// Length of the pool
func (p *Pool) Length() int {
	var AliveCounter int

	for _, w := range p.workers {
		if w.status == IsOk {
			AliveCounter++
		}
	}
	return AliveCounter
}

func (p *Pool) createWorker(id int) *Worker {

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
func (p *Pool) addJob(j Job) (Job, error) {
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
*/
