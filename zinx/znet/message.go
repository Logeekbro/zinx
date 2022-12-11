package znet

type Message struct {
	Id      uint32 // 消息Id
	DataLen uint32 // 消息长度
	Data    []byte // 消息数据
}
