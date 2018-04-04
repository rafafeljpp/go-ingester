package main

import (
	"fmt"
	"go-ingester/pool"
	"log"
	"strconv"
	"time"
)

// MyJob Estructura
type MyJob struct {
	payload string
	wait    time.Duration
}

func main() {

	var waitFor time.Duration

	//rejectedChannel := make(chan MyJob)
	pm := new(pool.Manager)

	//Iniciando...
	go pm.Start(10, 10)

	for {
		if pm.Length() == pm.GetWorkersQuantity() {
			break
		}
	}

	for i := 0; i < 10; i++ {
		waitFor = 0

		payload := "Mensaje Job " + strconv.Itoa(i)
		mj := MyJob{payload, waitFor}

		// FIX: problema en el Método AddJob. Solucionar la selección del worker.
		// Falla cuando hay mucha concurrencia. no Eliminar el "log.Printf" aquí abajo
		log.Printf("iniciando.. %d", i)
		pm.AddJob(mj)

	}

	fmt.Println("Cantidad de jobs", pm.CountJobs())
	pm.Stop()

	time.Sleep(time.Second * 20)

}

// Serialize Método para serializar el payload
func (j MyJob) Serialize() bool {
	//fmt.Println("Serializando...")
	return true
}

// Publish Publicar a donde yo quiera
func (j MyJob) Publish() bool {
	fmt.Println("Publicando..." + j.GetPayload())
	return true
}

// GetPayload Retornar el Payload
func (j MyJob) GetPayload() string {
	return j.payload
}

// Rejected Recibir los mensajes rechazados
func (j MyJob) Rejected() {
	fmt.Println("Mensaje Rechazado", j)
}
