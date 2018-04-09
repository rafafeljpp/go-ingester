package main

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"go-ingester/pool"
)

// MyJob Estructura
type MyJob struct {
	payload string
	init    time.Time
	id      int
}

func main() {

	pm := pool.NewManager(1000, 10)

	fmt.Println("Iniciando")
	pm.Start()

	for i := 0; i < 10000; i++ {
		start := time.Now()
		payload := "Job " + strconv.Itoa(i)
		mj := MyJob{payload, start, i}
		pm.AddJob(mj)

	}

	pm.Stop()
	mj := MyJob{"A", time.Now(), 0}
	pm.AddJob(mj)
	fmt.Println("End")
}

// IsValid Método para validar el Job
func (j MyJob) IsValid() bool {

	x := float64(j.id)

	if math.Mod(x, 2) == 0 {
		return true
	}

	return false
}

// Publish  donde yo quiera
func (j MyJob) Publish() (bool, error) {
	var elapse time.Duration
	retVal := true
	elapse = time.Since(j.init)
	fmt.Println("Mensaje Publicado: " + elapse.String())

	return retVal, nil
}

// Rejected Recibir los mensajes rechazados
func (j MyJob) Rejected(e error) {
	var elapse time.Duration
	elapse = time.Since(j.init)

	fmt.Println("Mensaje Rechazado: "+strconv.Itoa(j.id)+" -> "+e.Error(), elapse)
}
