package pool

import (
	"fmt"
	"time"
)

// Job type.
type Job struct {
	Payload string
	wait    time.Duration
}

func init() {
	fmt.Println("job inicializado")
}
