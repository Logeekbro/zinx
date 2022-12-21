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
	err := r.GetConnection().SendMsg(r.GetMsgId(), []byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("Callback Handle error:", err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
