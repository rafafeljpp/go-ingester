package pool

import (
	"fmt"
	"time"
)

// Job type.
type Job struct {
	// Payload
	Payload string

	// Wait
	wait time.Duration
}

func init() {
	fmt.Println("job package inicializado")
}
