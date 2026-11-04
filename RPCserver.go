package main

import (
	"log"
	"net"
	"net/rpc"
)

func main() {
	arith := new(main.Arith)

	rpc.Register(arith)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Listen error: %v", err)
	}
	defer listener.Close()
	log.Println("RPC server listening on port 1234....")
}
