package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	isClosed bool
	ExitChan chan bool
	Router   ziface.IRouter
}

var dp = NewDataPack()

func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) startReader() {
	fmt.Println("Starting read data from:", c.RemoteAddr())
	defer fmt.Printf("Connection(ID:%d) is closed\n", c.ConnID)
	defer c.Stop()
	for {
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	if err == io.EOF {
		//		break
		//	} else {
		//		fmt.Printf("Read data error:%s, ConnID=%d\n", err, c.ConnID)
		//		continue
		//	}
		//}
		headData := make([]byte, dp.GetHeadLen())
		_, err := c.GetTCPConnection().Read(headData)
		if err != nil {
			fmt.Println("Read headData error:", err)
		}
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("Unpack headData error:", err)
		}
		data := make([]byte, msg.GetDataLen())
		_, err = c.GetTCPConnection().Read(data)
		if err != nil {
			fmt.Println("Read data error:", err)
		}
		msg.SetData(data)
		request := Request{
			Conn: c,
			msg:  msg,
		}

		go func(req ziface.IRequest) {
			c.Router.PreHandle(req)
			c.Router.Handle(req)
			c.Router.PostHandle(req)
		}(&request)
	}
}

func (c *Connection) Start() {
	fmt.Printf("Connection starting(ID:%d)...\n", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.startReader()

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
