package main

import (
	"fmt"
	"net/rpc"
)

func main() {
	client, _ := rpc.Dial("tcp", "localhost:1234")
	args := &SaveArgs{"username", "golang_user"}
	var result Reply

	err := client.Call("DataStore.SaveData", args, &result)
	fmt.Println(result.Success)
	if err == nil && result.Success {
		fmt.Println("Data stored succesfully")
	}
}