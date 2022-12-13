package ziface

/**
将请求的消息封装到一个Message中, 定义一个抽象的接口
*/

type IMessage interface {
	GetId() uint32
	SerId(uint32)
	GetDataLen() uint32
	SetDataLen(uint32)
	GetData() []byte
	SetData([]byte)
}
