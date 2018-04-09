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
	init    time.Time
}

// Método para serializar el payload
func (j TestJob) IsValid() bool {
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
	myPool := pool.NewManager(10, 100)
	myPool.Start()

	if myPool.Length() != 10 {
		t.Fatal("There are ", myPool.Length())

	}

	t.Log(myPool)

}

func TestWorkFlow(t *testing.T) {
	var j TestJob

	myPool := pool.NewManager(10, 10)

	j = TestJob{"MiJob", 0, time.Now()}

	myPool.Start()

	time.Sleep(time.Second * 3)

	if myPool.Length() != 10 {
		t.Fatal("No se crearon los workers correctamente. Verifique el método Start()")
	}

	time.Sleep(time.Second * 3)

	myPool.AddJob(j)

	if myPool.CountJobs() > 0 {

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

/*
func TestProfiling(t *testing.T) {
	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// ... rest of the program ...

	myPool := pool.NewManager(10, 10)

	myPool.Start()
	time.Sleep(time.Millisecond * 10)

	for i := 0; i < 1000000; i++ {

		payload := "Mensaje Job " + strconv.Itoa(i)
		mj := TestJob{payload, 0, time.Now()}
		myPool.AddJob(mj)

	}

	myPool.Stop()

	// end

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f.Close()
	}
}
*/
