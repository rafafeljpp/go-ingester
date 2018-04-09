package main

import (
	"fmt"
	"strconv"
	"time"

	"go-ingester/pool"
)

// MyJob Estructura
type MyJob struct {
	payload string
	wait    time.Duration
	init    time.Time
}

func main() {
	var waitFor time.Duration

	pm := pool.NewManager(10, 10)

	// Iniciando...
	pm.Start()

	for i := 0; i < 10000; i++ {
		start := time.Now()
		waitFor = 0

		payload := "Mensaje Job " + strconv.Itoa(i)
		mj := MyJob{payload, waitFor, start}
		pm.AddJob(mj)

	}

	pm.Stop()
}

// Serialize Método para validar
func (j MyJob) IsValid() bool {
	//fmt.Println("Serializando...")
	return true
}

// Publish Publicar a donde yo quiera
func (j MyJob) Publish() bool {
	var elapse time.Duration
	elapse = time.Since(j.init)
	fmt.Println("Publicando..." + elapse.String())
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
