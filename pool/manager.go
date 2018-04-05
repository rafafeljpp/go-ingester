// Package pool Implementa la funcionalidad de un pool de hilos.
// @author: rafael.pellicer@gmail.com
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
	workers         []*Worker
	workersQuantity int
	queueSize       int
	cIndex          int
	mutex           sync.Mutex
}

// IJob Interface Contains methods for serialize and publish data.
type IJob interface {
	GetPayload() string
	Serialize() bool
	Publish() bool
	Rejected()
}

// Worker Estructura que simula un hilo y sus atributos.
type Worker struct {
	id        int
	startedAt time.Time
	messages  chan IJob
	signals   chan bool
	status    string
}

// Listen Escucha los mensajes recibidos en el canal de mensaje (cola).
func (w *Worker) Listen(wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		select {
		case msg := <-w.messages:
			if msg.Serialize() {
				if !msg.Publish() {
					msg.Rejected()
				}

			}

			break
		case <-w.signals:

			w.status = locked
			close(w.messages)

			return
		}
	}
}

// GetWorkersQuantity Retorna la cantidad de workers establecidos.
func (p *Manager) GetWorkersQuantity() int {
	return p.workersQuantity
}

// Start inicia los "hilos" del pool.
func (p *Manager) Start(WorkersQuantity int, QueueSize int) {

	var w *Worker
	var wg sync.WaitGroup
	p.queueSize = QueueSize
	p.workersQuantity = WorkersQuantity
	p.cIndex = -1

	// Creando workers
	for i := 0; i < p.workersQuantity; i++ {
		w = p.createWorker(i)
		wg.Add(1)
		go w.Listen(&wg)
	}

	wg.Wait()
}

// Stop Detiene la recepción de mensajes en las colas.
func (p *Manager) Stop() {
	clean := func(p *Manager) {
		fmt.Println("Finalizando el stop")
		for _, wk := range p.workers {
			wk.messages = nil
		}
	}
	defer clean(p)

	for _, wk := range p.workers {
		wk.signals <- true

	}
	for {

		if p.Length() == 0 {
			break
		}

	}
	return
}

// AddJob Método que añade trabajos a la cola de los workers.
func (p *Manager) AddJob(j IJob) (IJob, error) {
	var w *Worker

	p.mutex.Lock()
	if p.cIndex+1 >= len(p.workers) {
		p.cIndex = 0
	} else {
		p.cIndex++
	}
	p.mutex.Unlock()

	w = p.workers[p.cIndex]

	if len(w.messages)+1 < p.queueSize {
		w.messages <- j
		return j, nil
	}

	return j, errors.New("Rechazado, Máximo número de elementos en cola o mensaje bloqueado")

}

// Length Cantidad de workers activos.
func (p *Manager) Length() int {
	var AliveCounter int

	for _, w := range p.workers {
		if w.status == isOk {
			AliveCounter++
		}
	}
	return AliveCounter
}

// CountJobs retorna la cantidad de mensajes que se encuentran en las colas.
func (p *Manager) CountJobs() int {
	var counter int

	for _, w := range p.workers {
		counter += len(w.messages)
	}
	return counter
}

// createWorker Crea un hilo nuevo para el pool
func (p *Manager) createWorker(id int) *Worker {

	w := new(Worker)

	w.id = id
	w.status = isOk
	w.startedAt = time.Now()
	w.messages = make(chan IJob, p.queueSize)
	w.signals = make(chan bool)

	p.workers = append(p.workers, w)

	return w
}
