package ziface

/**
连接管理模块抽象层
*/

type IConnManager interface {
	// Add 添加一个连接
	Add(conn IConnection)
	// Remove 删除一个连接
	Remove(conn IConnection)
	// Get 根据Id获取一个连接
	Get(connId uint32) (IConnection, error)
	// Size 获取当前连接数
	Size() int
	// CloseAll 关闭所有连接
	CloseAll()
}
