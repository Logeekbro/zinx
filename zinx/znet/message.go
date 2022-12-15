package znet

type Message struct {
	Id      uint32 // 消息Id
	DataLen uint32 // 消息长度
	Data    []byte // 消息数据
}

// NewMessage 创建一个新的消息对象
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:   id,
		Data: data,
	}
}

func (m *Message) GetId() uint32 {
	return m.Id
}

func (m *Message) SerId(id uint32) {
	m.Id = id
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
