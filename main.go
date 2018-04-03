package main

import (
	"fmt"
	"go-ingester/pool"
	"log"
	"strconv"
	"time"
)

type MyJob struct {
	payload string
	wait    time.Duration
}

func main() {
	var waitFor time.Duration

	pm := new(pool.Manager)

	//Iniciando...
	go pm.Start(10, 10)

	for {
		if pm.Length() == pm.GetMaxWorkers() {
			break
		}
	}

	for i := 0; i < 50; i++ {
		waitFor = 0
		if i == 10 {
			waitFor = 12
		}

		payload := "Mensaje Job " + strconv.Itoa(i)
		MyJob := MyJob{payload, waitFor}

		// FIX: problema en el Método AddJob. Solucionar la selección del worker.
		// Falla cuando hay mucha concurrencia. no Eliminar el "log.Printf" aquí abajo
		log.Printf("iniciando.. %d", i)
		pm.AddJob(MyJob)

	}
	pm.Stop()

	time.Sleep(time.Second * 20)

}
// Método para serializar el payload
func (j MyJob) Serialize() bool {
	fmt.Println("Serializando...")
	return true
}

// Publicar a donde yo quiera
func (j MyJob) Publish() bool {
	fmt.Println("Publicando..." + j.GetPayload())
	return true
}
/
func (j MyJob) GetPayload() string {
	return j.payload
}
