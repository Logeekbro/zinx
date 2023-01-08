package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	BindServer ziface.IServer // 当前连接绑定的服务器
	Conn       *net.TCPConn
	ConnID     uint32
	isClosed   bool
	ExitChan   chan bool
	msgChan    chan []byte // 无缓冲的管道，用于读、写协程间的通信
	MsgHandler ziface.IMsgHandle
}

var dp = NewDataPack()

func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandle) *Connection {
	return &Connection{
		BindServer: server,
		Conn:       conn,
		ConnID:     connID,
		MsgHandler: msgHandler,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
	}
}

func (c *Connection) startWriter() {
	fmt.Println("Starting write data to:", c.RemoteAddr())
	defer fmt.Printf("Connection(ID:%d) Writer is closed\n", c.ConnID)
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.GetTCPConnection().Write(data); err != nil {
				fmt.Println("Write data error:", err)
				return
			}
		case <-c.ExitChan:
			// close(ExitChan)时会执行到此处
			fmt.Println("Writer exit...")
			return
		}
	}
}

func (c *Connection) startReader() {
	fmt.Println("Starting read data from:", c.RemoteAddr())
	defer fmt.Printf("Connection(ID:%d) Reader is closed\n", c.ConnID)
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
			return
		}
		msg.SetData(data)
		request := Request{
			Conn: c,
			msg:  msg,
		}
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 用户配置了工作池worker数则使用工作池处理request
			c.MsgHandler.SendRequestToTaskQueue(&request)
		} else {
			// 否则根据MsgId找到对应的处理API 执行(不使用工作池)
			go c.MsgHandler.DoMsgHandler(&request)
		}

	}
}

func (c *Connection) Start() {
	fmt.Printf("Connection starting(ID:%d)...\n", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.startReader()

	// 启动从当前连接写数据的业务
	go c.startWriter()
}

func (c *Connection) Stop() {
	fmt.Printf("Close connection(ID:%d)...\n", c.ConnID)
	if c.isClosed {
		return
	}
	c.Conn.Close()
	// 回收资源
	c.BindServer.GetConnMgr().Remove(c)
	close(c.ExitChan)
	close(c.msgChan)
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
	c.msgChan <- binaryMsg
	return nil
}
