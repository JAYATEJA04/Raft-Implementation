echo "Building Raft binary..."
go build -o raftnode node2.go

PORTS=(8001 8002 8003 8004 8005)

cleanup() {
    echo "Stopping all Raft nodes.."
    pkill -f ./raftnode
    exit
}

trap cleanup SIGINT SIGTERM

echo "Starting 5 Raft nodes in the background..."

for i in "${!PORTS[@]}"; do
    current_port=${PORTS[$i]}

    peers=""
    for peer_port in "${PORTS[@]}"; do
        if [ "$peer_port" != "$current_port" ]; then
            if [ -z "$peers" ]; then
                peers="$peer_port"
            else
                peers="$peers, $peer_port"
            fi
        fi
    done

    echo "Launching Node on port $current_port with peers: [$peers]"
    ./raftnode -port=$current_port -peers="$peers" > "node_$current_port.log" 2>&1 &
done

echo "Cluster is running! Press Ctrl+C to stop all nodes."
echo "Check logs using: tail -f node_8001.log"

wait

