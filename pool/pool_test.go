package pool

import (
	"fmt"
	"testing"
)

func TestStart(t *testing.T) {
	var p *Pool
	fmt.Println("Hola", p)
}

/*
func main() {

		var errc int
		rejectedChannel := make(chan Job)
		var waitFor time.Duration
		p := new(Pool)
		p.maxWorkers = 10
		p.maxElementsPerQueue = 100

		//Iniciando...
		go p.start()

		for {
			if p.len() == p.maxWorkers {
				break
			}
		}

		for i := 0; i < 100; i++ {
			waitFor = 0
			if i == 10 {
				waitFor = 12
			}
			MyJob := Job{Payload: "Mensaje " + strconv.Itoa(i), wait: waitFor}

			//p.addJob(MyJob)
			start := time.Now()
			log.Printf("iniciando.. %d", i)
			j, err := p.addJob(MyJob)
			elapsed := time.Since(start)
			log.Printf("Transcurrió %s en el mensaje %d", elapsed, i)

			if err != nil {
				errc++
				rejectedChannel <- j
			}

		}
		fmt.Println("Colas antes de stop():", p.len())
		p.stop()

		//time.Sleep(time.Second * 20)

		fmt.Println("Rejected: ", errc)
		fmt.Println("Colas después de stop():", p.len())

}
*/
