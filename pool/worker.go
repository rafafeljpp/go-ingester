package pool

import (
	"fmt"
	"sync"
	"time"
)

type Worker struct {
	id        int
	startedAt time.Time
	messages  chan Job
	signals   chan string
	status    string
}

// Listen worker.
func Listen(w *Worker, wg *sync.WaitGroup) {

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

func init() {
	fmt.Println("worker inicializado")
}
