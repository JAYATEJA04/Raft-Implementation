package main

import (
	"flag"
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
	"sync"
	"time"
)

//	type Node struct {
//		ID       int
//		Interval time.Duration
//		ctx      context.Context
//		cancel   context.CancelFunc
//	}

type Node struct {
	mu          sync.Mutex
	ID          int
	currentTerm int
	votedFor    []int
	state       string
	peers       []int
}

type RequestVoteArgs struct {
	Term        int
	CandidateID int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// func NewNode(nodeID int, interval time.Duration) *Node {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	return &Node{
// 		ID:       nodeID,
// 		Interval: interval,
// 		ctx:      ctx,
// 		cancel:   cancel,
// 	}
// }

func NewNode(nodeID int, state string, peerNodes []int) *Node {
	return &Node{
		ID:          nodeID,
		currentTerm: 0,
		votedFor:    nil,
		state:       state,
		peers:       peerNodes,
	}
}

func (n *Node) sendRequestVoteRPC(peerIdx int, args *RequestVoteArgs, Reply *RequestVoteReply) bool {
	n.mu.Lock()
	node := n.peers[peerIdx]
	defer n.mu.Unlock()

	if node == 0 {
		return false
	}

	fmt.Println("node: ", node, peerIdx, n.ID)

	return true
}

func (n *Node) callElectionLeader() {
	n.mu.Lock()
	n.state = "Candidate"
	n.votedFor = append(n.votedFor, n.ID)
	n.currentTerm++
	// fmt.Println("here, here: ", n.ID, n.currentTerm, n.votedFor, n.state)
	defer n.mu.Unlock()

	for idx, peer := range n.peers {
		if peer == n.ID {
			fmt.Println("matched, peer = n.Id!", peer, n.currentTerm, n.votedFor, n.state)
			continue
		}
		// fmt.Println("peer: ", peer)

		go func(peerIdx int) {
			args := RequestVoteArgs{
				Term:        n.currentTerm,
				CandidateID: n.ID,
			}
			var reply RequestVoteReply
			call := n.sendRequestVoteRPC(peerIdx, &args, &reply)
			fmt.Println("call: ", call)
		}(idx)
	}
}

func (n *Node) RequestVote(CandidateID *Node, Term *Node) error {
	fmt.Println("hey hey 1")
	return nil
}

func (n *Node) Start(wg *sync.WaitGroup) {
	// fmt.Print("hi hi hi hi")
	defer wg.Done()
	timeout := time.Duration(150+rand.IntN(150)) * time.Millisecond
	timer := time.NewTimer(timeout)

	heartBeatTimer := time.Duration(50+rand.IntN(50)) * time.Millisecond
	// heartbeatTicker := time.NewTicker(50 * time.Millisecond)
	// heartBeatTicker := time.NewTicker(heartBeatTimer)
	heartBeatTicker := time.NewTicker(heartBeatTimer)

	for {
		select {
		case msg := <-heartBeatTicker.C:
			fmt.Println("Received message", msg.Format("15:04:05.000"), n.ID)
		case <-timer.C:
			fmt.Println("Timeout!", timeout, n.ID)
			n.callElectionLeader()
			time.Sleep(15 * time.Millisecond)
			return
			// fmt.Println("hello")
			// timer.Stop()
		}
	}
}

func main() {
	port := flag.Int("port", 8001, "port-number")
	portMates := flag.String("peers", "", "port-mates")
	flag.Parse()

	var wg sync.WaitGroup

	convertedPortMates := strings.Split(*portMates, ",")
	intSlice := make([]int, 0, len(convertedPortMates))

	for _, peer := range convertedPortMates {
		if num, err := strconv.Atoi(strings.TrimSpace(peer)); err == nil {
			intSlice = append(intSlice, num)
		}
	}

	// fmt.Printf("Type of portMates: %T & Type of port: %T\n", portMates, port)
	// fmt.Println("port:", *port, ", converted portMates:", convertedPortMates, ", integer slice: ", intSlice)

	// for _, p := range RaftNode.peers {
	// 	wg.Add(1)
	// 	node := NewNode(p, "Follower")
	// 	nodes = append(nodes, node)
	// 	// count++
	// 	// time.Sleep(1500 * time.Millisecond)
	// 	go node.Start(&wg)
	// }

	// RaftNode := &Node{
	// 	ID:    *port,
	// 	state: "Follower",
	// 	peers: intSlice,
	// }

	registeredNode := NewNode(*port, "Follower", intSlice)

	wg.Add(1)
	go registeredNode.Start(&wg)
	fmt.Println("oi oi oi")

	// time.Sleep(1000 * time.Millisecond)
	wg.Wait()

	// shutdownSig := make(chan os.Signal, 1)
	// signal.Notify(shutdownSig, os.Interrupt, syscall.SIGTERM)
	// <-shutdownSig
	// fmt.Println("shutdown signal received, stopping nodes...")
	// for _, n := range nodes {
	// 	n.Stop()
	// }
}
