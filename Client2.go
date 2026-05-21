package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing: ", err)
	}

	args := &Args{4, 4}
	var result Result
	err = client.Call("Arith.Addition", args, &result)
	if err != nil {
		log.Fatal("Arith error: ", err)
	}
	fmt.Printf("Result: %d + %d = %d\n", args.A, args.B, result.Res)
}