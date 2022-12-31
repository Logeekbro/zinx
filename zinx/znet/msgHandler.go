package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandle struct {
	// 存放每个msgId所对应的处理方法
	Apis map[uint32]ziface.IRouter
	// 负责worker取任务的消息队列
	TaskQueue []chan ziface.IRequest
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// NewMsgHandler 创建/初始化 MsgHandle的方法
func NewMsgHandler() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]ziface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
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

// StartWorkerPool 启动一个Worker工作池 (开启工作池的动作只能发生一次，框架中只能存在一个worker工作池)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据WorkPoolSize来开启Worker, 每个Worker用go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 开辟一个任务管道
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动Worker
		go mh.startOneWorker(i)
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) startOneWorker(workerId int) {
	fmt.Printf("Worker(Id:%d) is started\n", workerId)
	taskChannel := mh.TaskQueue[workerId]
	for {
		select {
		case request := <-taskChannel:
			fmt.Printf("Start handle Connection(Id:%d), workId:%d\n", request.GetConnection().GetConnID(), workerId)
			mh.DoMsgHandler(request)
		}
	}
}

// SendRequestToTaskQueue 发送Request到消息队列
func (mh *MsgHandle) SendRequestToTaskQueue(request ziface.IRequest) {
	// 平均分配
	workerId := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerId] <- request
}
