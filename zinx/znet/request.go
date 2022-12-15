package znet

import "zinx/ziface"

type Request struct {
	// 和客户端建立的连接
	Conn ziface.IConnection
	// 客户端发送的消息
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgId() uint32 {
	return r.msg.GetId()
}
