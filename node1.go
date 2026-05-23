package main

import (
	"fmt"
	"net"
	"net/rpc"
	"sync"
)

type DataStore struct {
	mu sync.Mutex
	Items map[string]string
}

func (ds *DataStore) SaveData(args *SaveArgs, reply *Reply) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.Items[args.Key] = args.Value
	fmt.Printf("Node stored: [%s] = %s\n", args.Key, args.Value)
	reply.Success = true
	return nil
}

func (ds *DataStore) GetData(args *GetArgs, reply *GetReply) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	value, exists := ds.Items[args.Key]

	if exists {
		reply.Value = value
		reply.Successs = true
		fmt.Printf("Node : Providing data for the Key %s\n", args.Key)
	} else {
		reply.Successs = false
		fmt.Printf("Node : Key not found %s\n", args.Key)
	}
	return nil
}

func main() {
	store := &DataStore{Items: make(map[string]string)}
	rpc.Register(store)

	listener, err := net.Listen("tcp", ":1234")

	if err != nil {
		fmt.Println("Error listening: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port :1234....")
	rpc.Accept(listener)
}