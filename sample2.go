package main

import (
	"fmt"
	"sync"
)

func main() {
	nodePorts := map[int]int{
		1: 8001,
		2: 8002,
		3: 8003,
	}
	fmt.Println("hello")
	var wg sync.WaitGroup

	for id, port := range nodePorts {
		var peerPorts []int
		for peerId, peerPort := range nodePorts {
			if peerId != id {
				peerPorts = append(peerPorts, peerPort)
			}
		}

		fmt.Println("here after the iteration, ", peerPorts)
		wg.Add(1)

		go func(nodeId int, nodePort int, peerPort []int) {
			defer wg.Done()
			fmt.Printf("Node %d running on port %d. Peers: %v\n", nodeId, nodePort, peerPort)
		}(id, port, peerPorts)
	}

	wg.Wait()
}
