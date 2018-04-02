package pool

import (
	"testing"
	"time"
)

func TestInstance(t *testing.T) {
	var w *Worker
	//var wg *sync.WaitGroup

	w = new(Worker)
	w.id = 1
	w.startedAt = time.Now()
	w.status = IsOk

}
