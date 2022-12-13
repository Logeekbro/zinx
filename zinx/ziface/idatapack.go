package ziface

/**
封包、拆包模块
直接面向TCP连接中的数据流，用于处理TCP粘包问题
*/

type IDataPack interface {
	// GetHeadLen 获取包的头的长度
	GetHeadLen() uint32
	// Pack 封包
	Pack(IMessage) ([]byte, error)
	// UnPack 拆包
	UnPack([]byte) (IMessage, error)
}
