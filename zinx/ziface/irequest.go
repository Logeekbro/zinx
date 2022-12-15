package ziface

/**
IRequest接口
实际上是把来自客户端的请求连接信息和请求的数据包装成一个Request
*/

type IRequest interface {
	// GetConnection 得到当前连接
	GetConnection() IConnection
	// GetData 得到数据
	GetData() []byte
	// GetMsgId 获取消息Id
	GetMsgId() uint32
}
