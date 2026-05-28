package main

import (
	"fmt"
	"net/rpc"
)

func fetch(key string) {
	client, _ := rpc.Dial("tcp", "localhost:8001")

	args1 := &GetArgs{key}
	var response GetReply

	err0 := client.Call("DataStore.GetData", args1, &response)
	if err0 != nil {
		fmt.Println("RPC error", err0)
	} else if response.Successs {
		fmt.Printf("Client: Received Value -> %s\n", response.Value)
	} else {
		fmt.Println("Client: Key does not exist")
	}
}

func main() {
	client, err := rpc.Dial("tcp", "localhost:8001")

	if err != nil {
		fmt.Printf("Connection error: %v. Is the server runnong on port 8001?\n", err)
		return
	}
	args := &SaveArgs{"username", "golang_user"}
	var result Reply

	err = client.Call("DataStore.SaveData", args, &result)
	fmt.Println(result.Success)
	if err == nil && result.Success {
		fmt.Println("Data stored succesfully")
		fetch("username")
	}
}
