package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing: ", err)
	}

	args := &Args{7, 5}
	var reply Quotient
	err = client.Call("Arith.Divide", args, &reply)
	if err != nil {
		log.Fatal("arith error: ", err)
	}

	fmt.Printf("Arith: %d / %d = %d and remainder = %d\n", args.A, args.B, reply.Quo, reply.Rem)
}