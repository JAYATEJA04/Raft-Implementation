package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

type Args struct {
	A, B, C int
}

type Reply struct {
	Result int
}

type Storage struct {
	value int
}

type Calculator struct{}

func (arith *Calculator) Add(args *Args, reply *Reply) error {
	reply.Result = args.A + args.B
	return nil
}

func (arith *Calculator) Multiply(args *Args, reply *Reply) error {
	reply.Result = args.A * args.B
	return nil
}

func (arith *Calculator) Divide(args *Args, reply *Reply) error {
	if args.B == 0 {
		return errors.New("cannot divide by zero")
	}

	reply.Result = args.A / args.B
	return nil
}

func main() {
	calc := new(Calculator)
	rpc.Register(calc)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening: ", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port :8080...")
	rpc.Accept(listener)
}