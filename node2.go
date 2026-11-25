package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"net/rpc"
	"strconv"
	"strings"
	"sync"
	"time"
)

//	type Node struct {
//		NodeID       int
//		Interval time.Duration
//		ctx      context.Context
//		cancel   context.CancelFunc
//	}

type LogEntry struct {
	Command interface{}
	Term    int
}
type Node struct {
	mu          sync.Mutex
	NodeID      int
	port        int
	currentTerm int
	votedFor    int
	Log         []LogEntry
	state       string
	peers       []int
	Heartbeat   chan bool
}

type RequestVoteArgs struct {
	Term         int
	CandidateID  int
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteReply struct {
	Term        int
	VoteGranted bool
}

// func NewNode(nodeID int, interval time.Duration) *Node {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	return &Node{
// 		NodeID:       nodeID,
// 		Interval: interval,
// 		ctx:      ctx,
// 		cancel:   cancel,
// 	}
// }

func NewNode(nodeID int, port int, state string, peerNodes []int) *Node {
	return &Node{
		NodeID:      nodeID,
		port:        port,
		currentTerm: 0,
		state:       state,
		peers:       peerNodes,
		Log:         make([]LogEntry, 1),
	}
}

func printHandleFunction(conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	fmt.Println("hey hey 1")
}

func (n *Node) StartServer(wg *sync.WaitGroup) {
	defer wg.Done()

	rpcServer := rpc.NewServer()
	registrationErr := rpcServer.Register(n)
	if registrationErr != nil {
		log.Fatalf("Node%d failed to register RPC: %v", n.NodeID, registrationErr)
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(n.port))

	if err != nil {
		fmt.Println("Error starting listener: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("inside startserver function", n.port)

	for {
		fmt.Println("inside for loop")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err)
			continue
		}
		wg.Add(1)
		go func(c net.Conn) {
			defer wg.Done()
			rpcServer.ServeConn(conn)
		}(conn)
	}
}

func (n *Node) RequestVoteRPC(args *RequestVoteArgs, reply *RequestVoteReply) error {
	fmt.Println("inside request vote RPC")
	n.mu.Lock()
	defer n.mu.Unlock()

	if args.Term > n.currentTerm {
		n.currentTerm = args.Term
		n.state = "Follower"
		n.votedFor = -1
	}

	if args.Term < n.currentTerm {
		reply.Term = n.currentTerm
		reply.VoteGranted = false
		return nil
	}

	localLastLogIndex := len(n.Log) - 1
	localLastLogTerm := n.Log[localLastLogIndex].Term

	logUptoDate := false
	if args.LastLogTerm > localLastLogTerm {
		logUptoDate = true
	} else if args.LastLogTerm == localLastLogTerm && args.LastLogIndex >= localLastLogIndex {
		logUptoDate = true
	}

	if (n.votedFor == -1 || n.votedFor == args.CandidateID) && logUptoDate {
		n.votedFor = args.CandidateID
		reply.VoteGranted = true
	} else {
		reply.VoteGranted = false
	}
	reply.Term = n.currentTerm
	return nil
}

func (n *Node) SendRequestVoteRPC(peerPort int, args *RequestVoteArgs, Reply *RequestVoteReply) bool {
	listener, err := rpc.Dial("tcp", "localhost:"+strconv.Itoa(peerPort))
	if err != nil {
		fmt.Println("Failed to connect %s", peerPort)
		return false
	}

	defer listener.Close()

	err1 := listener.Call("Node.RequestVoteRPC", args, Reply)
	// time.Sleep(1000 * time.Millisecond)
	fmt.Println("after request vote RPC")
	if err1 != nil {
		fmt.Println("callOk is not ok", n.port, err1)
		return false
	}
	return true
}

func (n *Node) getLastLogTerm() int {
	lastLogIndex := len(n.Log) - 1
	lastLogTerm := n.Log[lastLogIndex].Term

	return lastLogTerm
}

func (n *Node) getLastLogIndex() int {
	lastLongIndex := len(n.Log) - 1

	return lastLongIndex
}

func (n *Node) callElectionLeader() {
	n.mu.Lock()
	n.state = "Candidate"
	n.votedFor = n.NodeID
	n.currentTerm++
	// fmt.Println("here, here: ", n.NodeID, n.currentTerm, n.votedFor, n.state)
	defer n.mu.Unlock()

	for _, peer := range n.peers {
		// if peer == n.port {
		// 	fmt.Println("matched, peer = n.Id!", n.NodeID, peer, n.currentTerm, n.votedFor, n.state)
		// 	continue
		// }
		// fmt.Println("peer: ", peer)

		go func(peer int) {
			args := RequestVoteArgs{
				Term:         n.currentTerm,
				CandidateID:  n.NodeID,
				LastLogIndex: n.getLastLogIndex(),
				LastLogTerm:  n.getLastLogTerm(),
			}
			fmt.Println("current term is: ", args.Term, " last log index: ", args.LastLogIndex, " and last log term is: ", args.LastLogTerm)
			var reply RequestVoteReply
			call := n.SendRequestVoteRPC(peer, &args, &reply)
			fmt.Println("call: ", call)
		}(peer)
	}
}

func getNewtimer() time.Duration {
	timeoutPeriod := time.Duration(150+rand.IntN(150)) * time.Millisecond
	return timeoutPeriod
}

func (n *Node) RunElectionLoop() {
	newTimeoutPeriod := getNewtimer()
	timer := time.NewTimer(newTimeoutPeriod)

	// heartbeatTicker := time.Duration(50+rand.IntN(20)) * time.Millisecond
	// // heartbeatTicker := time.NewTicker(50 * time.Millisecond)
	// // heartBeatTicker := time.NewTicker(heartBeatTimer)
	// heartBeatTicker := time.NewTicker(heartbeatTicker)

	for {
		select {
		// case msg := <-heartBeatTicker.C:
		// 	fmt.Println("ticker", msg.Format("15:04:05.000"), n.port)
		case <-n.Heartbeat:
			fmt.Println("Node received a message. Resetting timer......")
			if !timer.Stop() {
				<-timer.C
			}
			timer.Reset(time.Duration(150+rand.IntN(150)) * time.Millisecond)
		case <-timer.C:
			fmt.Println("Timeout!", newTimeoutPeriod, n.port)
			n.callElectionLeader()
			newTimeoutPeriod = getNewtimer()
			timer.Reset(newTimeoutPeriod)
		}
	}
}

func main() {
	port := flag.Int("port", 8001, "port-number")
	portMates := flag.String("peers", "", "port-mates")
	flag.Parse()

	var wg sync.WaitGroup

	splittedPortMates := strings.Split(*portMates, ",")
	convertedPortMates := make([]int, 0, len(splittedPortMates))

	for _, peer := range splittedPortMates {
		if num, err := strconv.Atoi(strings.TrimSpace(peer)); err == nil {
			convertedPortMates = append(convertedPortMates, num)
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
	// 	NodeID:    *port,
	// 	state: "Follower",
	// 	peers: intSlice,
	// }

	registeredNode := NewNode(1, *port, "Follower", convertedPortMates)
	// fmt.Println(convertedPortMates, *port)

	wg.Add(1)
	// time.Sleep(500 * time.Millisecond)
	// registeredNode.Start(&wg)
	go registeredNode.StartServer(&wg)
	registeredNode.RunElectionLoop()

	// time.Sleep(1000 * time.Millisecond)
	// select {}
	// time.Sleep(5000 * time.Millisecond)
	fmt.Println("oi oi oi")
	wg.Wait()

	// shutdownSig := make(chan os.Signal, 1)
	// signal.Notify(shutdownSig, os.Interrupt, syscall.SIGTERM)
	// <-shutdownSig
	// fmt.Println("shutdown signal received, stopping nodes...")
	// for _, n := range nodes {
	// 	n.Stop()
	// }
}
