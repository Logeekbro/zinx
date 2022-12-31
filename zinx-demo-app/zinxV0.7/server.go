package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (router *PingRouter) Handle(r ziface.IRequest) {
	fmt.Println("-->Recv client data:", string(r.GetData()))
	err := r.GetConnection().SendMsg(200, []byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("Callback Handle error:", err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (router *HelloRouter) Handle(r ziface.IRequest) {
	err := r.GetConnection().SendMsg(201, []byte("hello...hello...hello...\n"))
	if err != nil {
		fmt.Println("Callback Handle error:", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}
