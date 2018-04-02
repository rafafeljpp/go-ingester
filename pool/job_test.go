package pool

import "testing"

func TestInit(t *testing.T) {
	var x Job
	x.Payload = "S"
	x.wait = 1

}
