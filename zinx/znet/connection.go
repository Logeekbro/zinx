package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn     *net.TCPConn
	ConnID   uint32
	Callback ziface.HandleFunc
	isClosed bool
	ExitChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callback ziface.HandleFunc) *Connection {
	return &Connection{
		Conn:     conn,
		ConnID:   connID,
		Callback: callback,
		isClosed: false,
		ExitChan: make(chan bool, 1),
	}
}

func (c *Connection) startReader() {
	fmt.Println("Starting read data from:", c.RemoteAddr())
	defer fmt.Println("Read over, exit...")
	defer c.Stop()
	for {
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("Read data error:%s, ConnID=%d\n", err, c.ConnID)
				continue
			}
		}
		if err = c.Callback(c.Conn, buf, cnt); err != nil {
			fmt.Println("Callback error:", err)
		}
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

func (c *Connection) Send(data []byte) error {
	//TODO implement me
	panic("implement me")
}
