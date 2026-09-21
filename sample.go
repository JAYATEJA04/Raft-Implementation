package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	// 1. Define the flag (returns a string pointer)
	portMatesRaw := flag.String("peers", "", "comma-separated list of port mates")
	flag.Parse()

	// 2. Handle empty input safely
	if *portMatesRaw == "" {
		fmt.Println("No peers provided.")
		return
	}

	// 3. Split the single string into a slice (array) of strings
	portMates := strings.Split(*portMatesRaw, ",")

	// Now portMates is a []string slice: ["8081", "8082", "8083"]
	fmt.Printf("Parsed slice: %v (Type: %T)\n", portMates, portMates)
}
