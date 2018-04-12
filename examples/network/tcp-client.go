package main

import (
	"fmt"
	"net"
	"time"
)

func mains() {

	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("fail")
	}
	start := time.Now()
	for i := 0; i < 5; i++ {
		conn.Write([]byte("Test message" + "\n"))
	}
	elapsed := time.Since(start)
	fmt.Println("Tiempo transcurrido: ", elapsed)
	conn.Close()
}
