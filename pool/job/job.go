package job

import "fmt"

// Job type.
type Job struct {
	Payload string
}

func init() {
	fmt.Println("job inicializado")
}
