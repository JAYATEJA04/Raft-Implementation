package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Node struct {
	ID       string
	Interval time.Duration
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewNode(nodeID string, interval time.Duration) *Node {
	ctx, cancel := context.WithCancel(context.Background())
	return &Node{
		ID:       nodeID,
		Interval: interval,
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (n *Node) Start() {
	fmt.Printf("[Node %s] Initialising and starting node services...\n", n.ID)
	go n.runHeartBeatLoop()
}

func (n *Node) runHeartBeatLoop() {
	ticker := time.NewTicker(n.Interval)
	defer ticker.Stop()

	fmt.Printf("[Node %s] Heartbeat timer loop has started.\n", n.ID)

	for {
		select {
		case <-n.ctx.Done():
			fmt.Printf("[Node %s] Stopping heartbeat gracefully", n.ID)
			return
		case t := <-ticker.C:
			n.sendHeartBeat(t)
		}
	}
}

func (n *Node) sendHeartBeat(t time.Time) {
	fmt.Printf("[Node %s] has sent Heartbeat to cluster at %s\n", n.ID, t.Format("15:04:05"))
}

func (n *Node) Stop() {
	fmt.Printf("[Node %s] shutting down node...\n", n.ID)
	n.cancel()
}

func main() {
	// node := NewNode("node-alpha", 3*time.Second)
	// node.Start()

	peers := []string{"8001", "8002", "8003", "8004", "8005"}
	var node Node
	var count int

	for range peers {
		node := NewNode("node-alpha", 3*time.Second)
		count++
		go node.Start()
	}

	shutdownSig := make(chan os.Signal, 1)
	signal.Notify(shutdownSig, os.Interrupt, syscall.SIGTERM)

	<-shutdownSig
	node.Stop()

	time.Sleep(500 * time.Millisecond)
}
