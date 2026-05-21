package main

import (
	"fmt"
	"net/rpc"
)

type Args struct {
	A, B int
}

type Reply struct {
	Result int
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}

	defer client.Close()

	args := &Args{3, 4}
	var result Reply
	
	err = client.Call("Calculator.Multiply", args, &result)
	if err != nil {
		fmt.Println("Error calling calculator.Multiply: ", err)
		return
	}

	mainClientResult := result

	fmt.Println("Result of multiplication: ", result)
	fmt.Printf("\nand the type is: %T", mainClientResult)
}