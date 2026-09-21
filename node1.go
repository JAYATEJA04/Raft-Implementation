package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"strings"
	"sync"
)

// type State int
//
//	type Node struct {
//		id            string
//		peers         []string
//		state         State
//		lastHeartBeat time.Time
//	}
type DataStore struct {
	mu    sync.Mutex
	Items map[string]string
	// NodeAddresses ["localhost:8001", "localhost:8002", "localhost:8003", "localhost:8004"]string
	NodeAddresses []string
	Port          string
}

func (ds *DataStore) SaveData(args *SaveArgs, reply *Reply) error {
	ds.mu.Lock()
	defer ds.mu.Unlock()

	ds.Items[args.Key] = args.Value
	fmt.Printf("Node stored: [%s] = %s\n", args.Key, args.Value)

	for _, addr := range ds.NodeAddresses {
		go func(target string) {
			client, err := rpc.Dial("tcp", target)
			if err != nil {
				fmt.Println("failed to reach %s\n", target)
				return
			}
			defer client.Close()
			var secondaryReply Reply
			client.Call("DataStore.SaveData", args, &secondaryReply)
		}(addr)
	}
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
	port := flag.String("port", "8001", "Port to listen on")
	secondary := flag.String("addresses", "", "comma-separated secondary addresses")
	flag.Parse()

	node := &DataStore{Items: make(map[string]string), Port: *port}
	if *secondary != "" {
		node.NodeAddresses = strings.Split(*secondary, ",")
	}

	fmt.Println("node.nodeaddresses: ", node.NodeAddresses)

	rpc.Register(node)
	listener, _ := net.Listen("tcp", ":"+*port)
	fmt.Printf("Node started on %s. Replicating to: %v\n", *port, node.NodeAddresses)
	rpc.Accept(listener)
}
