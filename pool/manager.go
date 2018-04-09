// Package pool Implementa la funcionalidad de un pool de hilos.
// @author: rafael.pellicer@gmail.com
package pool

import (
	"container/ring"
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
	mutex           sync.Mutex
	indexSelector   *ring.Ring
}

// IJob Interface Contains methods for serialize and publish data.
type IJob interface {
	GetPayload() string
	IsValid() bool
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

// NewManager create a new Manager instance
func NewManager(WorkersQuantity int, QueueSize int) *Manager {
	return &Manager{
		queueSize:       QueueSize,
		workersQuantity: WorkersQuantity,
		indexSelector:   ring.New(WorkersQuantity),
	}
}

// NewWorker create a new workder instance
func newWorker(id int, QueueSize int) *Worker {
	return &Worker{
		id:        id,
		status:    isOk,
		startedAt: time.Now(),
		messages:  make(chan IJob, QueueSize),
		signals:   make(chan bool),
	}
}

// GetWorkersQuantity Retorna la cantidad de workers establecidos.
/*func (p *Manager) GetWorkersQuantity() int {
	return p.workersQuantity
}*/

// Start inicia los "hilos" del pool.
func (p *Manager) Start() {

	// Creando workers
	for i := 0; i < p.workersQuantity; i++ {
		p.indexSelector.Value = i

		w := newWorker(i, p.queueSize)
		go w.listen()

		p.workers = append(p.workers, w)
		p.indexSelector = p.indexSelector.Next()
	}

}

// Stop Detiene la recepción de mensajes en las colas.
func (p *Manager) Stop() {

	defer func(p *Manager) {
		fmt.Println("Finalizando el stop")
		for _, wk := range p.workers {
			wk.messages = nil
		}
	}(p)

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

	p.mutex.Lock()
	w := p.workers[p.indexSelector.Next().Value.(int)]
	p.mutex.Unlock()

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

// Listen Escucha los mensajes recibidos en el canal de mensajes (cola).
func (w *Worker) listen() {

	for {
		select {
		case msg := <-w.messages:
			if msg.IsValid() {
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
