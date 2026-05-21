package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"time"
)

// type Args struct {
//     A, B int
// }

// type Quotient struct {
//     Quo, Rem int
// }

// type Arith int

func (t *Arith) Divide(args *Args, reply *Quotient) error {
    if args.B == 0 {
        return errors.New("Divide by zero")
    }

    reply.Quo = args.A / args.B
    reply.Rem = args.A % args.B
    return nil
}

func (t *Arith) Addition(args *Args, reply *Result) error {
    reply.Res = args.A + args.B
    return nil
}

func main() {
    arith := new(Arith)
    rpc.Register(arith)

    listen, error := net.Listen("tcp", ":1234")

    if error != nil {
        log.Fatal("listen error: ", error)
    }

    log.Println("Server started on :1234")
    time.Sleep(100 * time.Millisecond)
    for {
        conn, err := listen.Accept()
        if err != nil {
            continue
        }
        go rpc.ServeConn(conn)
    }
}