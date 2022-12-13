package znet

import (
	"bytes"
	"encoding/binary"
	"zinx/ziface"
)

// DataPack 封包、拆包的具体模块
type DataPack struct {
}

func (d *DataPack) GetHeadLen() uint32 {
	// DataLen uint32(4 bytes) + DataId uint32(4 bytes)
	return 8
}

func (d *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	// 创建一个bytes缓冲
	bytesData := bytes.NewBuffer([]byte{})
	// 写入内容长度 (binary.LittleEndian: 涉及大小端的问题, 知识盲区, 需要去了解一下)
	if err := binary.Write(bytesData, binary.LittleEndian, message.GetDataLen()); err != nil {
		return nil, err
	}
	// 写入消息类型(Id)
	if err := binary.Write(bytesData, binary.LittleEndian, message.GetId()); err != nil {
		return nil, err
	}
	// 写入数据
	if err := binary.Write(bytesData, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}
	return bytesData.Bytes(), nil
}

func (d *DataPack) UnPack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个读取二进制数据的ioReader
	reader := bytes.NewReader(binaryData)

	// Message信息
	message := &Message{}

	// 读DataLen
	if err := binary.Read(reader, binary.LittleEndian, &message.DataLen); err != nil {
		return nil, err
	}

	// 读MsgId
	if err := binary.Read(reader, binary.LittleEndian, &message.Id); err != nil {
		return nil, err
	}

	return message, nil

}
