package znet

import (
	"fmt"
	"zinx/ziface"
)

type MsgHandle struct {
	// 存放每个msgId所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// NewMsgHandler 创建/初始化 MsgHandle的方法
func NewMsgHandler() *MsgHandle {
	return &MsgHandle{Apis: make(map[uint32]ziface.IRouter)}
}

func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 根据request中的msgId获取router
	router, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		// 没获取到对应的router
		fmt.Println("Router NOT FOUND! msgId:", request.GetMsgId(), ", please register!")
	} else {
		router.PreHandle(request)
		router.Handle(request)
		router.PostHandle(request)
	}
}

func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	// 先判断msgId是否已经绑定了路由
	if _, ok := mh.Apis[msgId]; ok {
		// 已经绑定了
		panic(fmt.Sprintf("MsgId(%d) has bound a Router\n", msgId))
	}
	// 绑定路由
	mh.Apis[msgId] = router
	fmt.Printf("MsgId(%d) add router success!\n", msgId)
}
