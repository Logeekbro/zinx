package znet

import (
	"zinx/ziface"
)

// BaseRouter 实现Router时，先嵌入这个BaseRouter，然后根据需要重写BaseRouter的函数实现自定义Router
type BaseRouter struct{}

func (router *BaseRouter) PreHandle(r ziface.IRequest) {

}

func (router *BaseRouter) Handle(r ziface.IRequest) {

}

func (router *BaseRouter) PostHandle(r ziface.IRequest) {

}
