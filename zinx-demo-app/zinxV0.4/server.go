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
	err := r.GetConnection().SendMsg(r.GetMsgId(), []byte("before ping...\n"))
	if err != nil {
		fmt.Println("callback PreHandle error:", err)
	}
}

func (router *PingRouter) Handle(r ziface.IRequest) {
	err := r.GetConnection().SendMsg(r.GetMsgId(), []byte("ping...\n"))
	if err != nil {
		fmt.Println("callback Handle error:", err)
	}
}

func (router *PingRouter) PostHandle(r ziface.IRequest) {
	err := r.GetConnection().SendMsg(r.GetMsgId(), []byte("after ping...\n"))
	if err != nil {
		fmt.Println("callback PostHandle error:", err)
	}
}

func main() {
	s := znet.NewServer("[Zinx0.4]")
	s.AddRouter(1, &PingRouter{})
	s.Serve()
}
