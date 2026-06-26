package main

type SaveArgs struct {
	Key   string
	Value string
}

type Reply struct {
	Success bool
}

type GetArgs struct {
	Key string
}

type GetReply struct {
	Value    string
	Successs bool
}

// type RaftNode struct {
// 	mu            sync.Mutex
// 	Id            int
// 	Port          string
// 	Peers         []string
// 	State         State
// 	LastHeartBeat time.Time
// }
