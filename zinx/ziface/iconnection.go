package ziface

import "net"

// 定义连接模块的抽象层

type IConnection interface {
	// Start 让当前连接开始工作
	Start()

	// Stop 让当前连接停止工作
	Stop()

	// GetTCPConnection 获取当前连接绑定的 socket conn
	GetTCPConnection() *net.TCPConn

	// GetConnID 获取当前连接的连接ID
	GetConnID() uint32

	// RemoteAddr 获取客户端的TCP状态 IP Port
	RemoteAddr() net.Addr

	// SendMsg 发送封装好的数据包给客户端
	SendMsg(id uint32, data []byte) error
}
