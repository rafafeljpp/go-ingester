package pool

import (
	"fmt"
	. "go-ingester/pool/job"
	"sync"
	"time"
)

const (
	IsOk   string = "OK"
	Locked string = "Locked"
)

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

func init() {
	fmt.Println("Worker del paquete pool inicializado")
}
