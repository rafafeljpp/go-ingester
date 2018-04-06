package pool_test

import (
	"fmt"
	"go-ingester/pool"
	"testing"
	"time"
)

type TestJob struct {
	payload string
	wait    time.Duration
}

// Método para serializar el payload
func (j TestJob) Serialize() bool {
	fmt.Println("Serializando...")
	return true
}

// Publicar a donde yo quiera
func (j TestJob) Publish() bool {
	fmt.Println("Publicando..." + j.GetPayload())
	return true
}

// Retornar el Payload
func (j TestJob) GetPayload() string {
	return j.payload
}

// Rechazados
func (j TestJob) Rejected() {
	fmt.Println("Rechazados")
}

/**************************************************
 * TESTS!
 **************************************************/

func TestPoolInstanciation(t *testing.T) {
	myPool := new(pool.Manager)

	t.Log(myPool)
}

func TestGetMaxWorkers(t *testing.T) {
	myPool := new(pool.Manager)
	if myPool.GetWorkersQuantity() != 0 {
		t.Fail()
	}

}

func TestWorkFlow(t *testing.T) {
	var j TestJob

	myPool := new(pool.Manager)

	j = TestJob{"MiJob", 0}

	go myPool.Start(10, 10)

	time.Sleep(time.Second * 3)

	if myPool.Length() != 10 {
		t.Fatal("No se crearon los workers correctamente. Verifique el método Start()")
	}

	time.Sleep(time.Second * 3)

	myPool.AddJob(j)

	if myPool.CountJobs() > 0 {
		//fmt.Println(myPool.CountJobs())
		t.Fatal("La cantidad de trabajos no es la esperada se estén recibiendo los mensajes en el canal")
	}

	if myPool.Length() != 10 {
		t.Fatal("La cantidad de workers no es la esperada.")
	}

	myPool.Stop()

	time.Sleep(time.Second * 5)

	if myPool.Length() > 0 {
		t.Fatal("La cantidad de trabajos, luego del Stop() no es la esperada.")
	}

}
