package main

import (
	"fmt"
	"go-ingester/pool"
	"log"
	"strconv"
	"time"
)

func main() {

	var errc int
	rejectedChannel := make(chan pool.Job)
	var waitFor time.Duration
	pm := new(pool.Manager)

	//Iniciando...
	go pm.Start(10, 100)

	for {
		if pm.Length() == pm.GetMaxWorkers() {
			break
		}
	}

	for i := 0; i < 100; i++ {
		waitFor = 0
		if i == 10 {
			waitFor = 12
		}
		MyJob := pool.Job{Payload: "Mensaje " + strconv.Itoa(i), Wait: waitFor}

		//p.addJob(MyJob)
		start := time.Now()
		log.Printf("iniciando.. %d", i)
		j, err := pm.AddJob(MyJob)

		elapsed := time.Since(start)
		log.Printf("Transcurrió %s en el mensaje %d", elapsed, i)

		if err != nil {
			errc++
			rejectedChannel <- j
		}

	}
	//	fmt.Println("Colas antes de stop():", pm.Length())
	pm.Stop()

	time.Sleep(time.Second * 20)

	fmt.Println("Rejected: ", errc)
	fmt.Println("Colas después de stop():", pm.Length())
}
