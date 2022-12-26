package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	ExitChan   chan bool
	MsgHandler ziface.IMsgHandle
}

var dp = NewDataPack()

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	return &Connection{
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
	}
}

func (c *Connection) startReader() {
	fmt.Println("Starting read data from:", c.RemoteAddr())
	defer fmt.Printf("Connection(ID:%d) is closed\n", c.ConnID)
	defer c.Stop()
	for {
		headData := make([]byte, dp.GetHeadLen())
		_, err := c.GetTCPConnection().Read(headData)
		if err != nil {
			if err == io.EOF {

			} else {
				fmt.Println("Read headData error:", err)
			}
			return
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Unpack headData error:", err)
			return
		}
		data := make([]byte, msg.GetDataLen())
		_, err = c.GetTCPConnection().Read(data)
		if err != nil {
			if err == io.EOF {

			} else {
				fmt.Println("Read headData error:", err)
			}
		}
		msg.SetData(data)
		request := Request{
			Conn: c,
			msg:  msg,
		}
		// 根据MsgId找到对应的处理API 执行
		go c.MsgHandler.DoMsgHandler(&request)
	}
}

func (c *Connection) Start() {
	fmt.Printf("Connection starting(ID:%d)...\n", c.ConnID)
	// 启动从当前连接读数据的业务
	c.startReader()

	// TODO 启动从当前连接写数据的业务

}

func (c *Connection) Stop() {
	fmt.Printf("Close connection(ID:%d)...\n", c.ConnID)
	if c.isClosed {
		return
	}
	c.Conn.Close()
	// 回收资源
	close(c.ExitChan)

}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(id uint32, data []byte) error {
	binaryMsg, err := dp.Pack(NewMessage(id, data))
	if err != nil {
		return errors.New(fmt.Sprintln("Pack msg error:", err))
	}
	_, err = c.GetTCPConnection().Write(binaryMsg)
	if err != nil {
		return errors.New(fmt.Sprintln("Send msg to client error:", err))
	}
	return nil
}
