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

	// Send 发送数据给客户端
	Send(data []byte) error
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
