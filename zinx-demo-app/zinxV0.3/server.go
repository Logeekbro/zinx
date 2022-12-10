package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (router *PingRouter) PreHandle(r ziface.IRequest) {
	_, err := r.GetConnection().Send([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("callback PreHandle error:", err)
	}
}

func (router *PingRouter) Handle(r ziface.IRequest) {
	_, err := r.GetConnection().Send([]byte("ping...\n"))
	if err != nil {
		fmt.Println("callback Handle error:", err)
	}
}

func (router *PingRouter) PostHandle(r ziface.IRequest) {
	_, err := r.GetConnection().Send([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("callback PostHandle error:", err)
	}
}

func main() {
	s := znet.NewServer("[Zinx0.3]")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
