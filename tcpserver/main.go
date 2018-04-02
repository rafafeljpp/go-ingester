package main

import (
	"fmt"
	p "go-ingester/pool"
)

func main() {

	p := new(p.Pool)
	fmt.Println(p.start())

}
