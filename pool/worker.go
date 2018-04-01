package pool

import (
	"fmt"
	"sync"
	"time"
)

// Worker type
type Worker struct {
	id        int
	startedAt time.Time
	messages  chan Job
	signals   chan string
	status    string
}

func init() {
	fmt.Println("worker inicializado")
}

func (w *Worker) listen(wg *sync.WaitGroup) {

	defer wg.Done()
	for {
		select {
		case msg := <-w.messages:
			fmt.Println("Mensaje recibido por el worker: ", w.id, msg.Payload)
			time.Sleep(time.Second * msg.wait)
			break
		case <-w.signals:
			fmt.Println("finalizando el worker ", w.id)
			w.status = Locked
			close(w.messages)

			return
		}
	}
}
