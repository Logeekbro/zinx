package ziface

/**
消息管理抽象层
*/

type IMsgHandle interface {
	// DoMsgHandler 调度/执行 对应的Router消息处理方法
	DoMsgHandler(request IRequest)
	// AddRouter 为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
	// StartWorkerPool 启动Worker工作池
	StartWorkerPool()
	// SendRequestToTaskQueue 发送request到消息队列进行处理
	SendRequestToTaskQueue(request IRequest)
}
