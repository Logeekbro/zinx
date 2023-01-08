package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

/**
连接管理模块
*/

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 保护连接集合的读写锁
}

// NewConnManager 创建连接管理模块
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (c *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源map, 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	// 将连接添加到连接集合中
	c.connections[conn.GetConnID()] = conn
	fmt.Printf("Add Conn(Id:%d) success, current conn num:%d\n", conn.GetConnID(), c.Size())
}

func (c *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源map, 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	// 根据Id删除连接
	delete(c.connections, conn.GetConnID())
	fmt.Printf("Remove Conn(Id:%d) success, current conn num:%d\n", conn.GetConnID(), c.Size())
}

func (c *ConnManager) Get(connId uint32) (ziface.IConnection, error) {
	// 保护共享资源map, 加读锁
	c.connLock.RLock()
	defer c.connLock.RUnlock()
	if conn, ok := c.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New(fmt.Sprintf("Conn(Id:%d) not exist!", connId))
	}
}

func (c *ConnManager) Size() int {
	return len(c.connections)
}

func (c *ConnManager) CloseAll() {
	// 保护共享资源map, 加写锁
	c.connLock.Lock()
	defer c.connLock.Unlock()
	// 关闭并删除连接
	for connId, conn := range c.connections {
		//TODO 关闭, 存在死锁, 需要解决
		conn.Stop()
		// 删除
		delete(c.connections, connId)
	}
	fmt.Println("All Connection is closed, connNum=", c.Size())
}
