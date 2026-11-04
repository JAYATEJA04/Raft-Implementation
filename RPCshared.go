package main

import "errors"

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, reply *int) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	*reply = args.A / args.B
	return nil
}
