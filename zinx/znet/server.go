package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

type Server struct {
	// 服务器名称
	Name string
	// 服务器IP版本
	IPVersion string
	// 服务器IP
	IP string
	// 服务器监听的端口
	Port uint16
	// 当前server的消息管理模块，用来绑定MsgId和对应的业务处理业务API关系
	MsgHandler ziface.IMsgHandle
	// 当前server的连接管理模块
	ConnMgr ziface.IConnManager
}

func (s *Server) Start() {
	fmt.Println("[Zinx]Starting server:", s.Name)
	fmt.Println("[Zinx]Version:", utils.GlobalObject.Version)
	fmt.Printf("[Zinx]MaxConn: %d, MaxPackageSize: %d\n", utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPackageSize)
	// 服务器启动步骤
	// 0、启动消息队列和worker工作池
	if utils.GlobalObject.WorkerPoolSize > 0 {
		s.MsgHandler.StartWorkerPool()
	}
	// 1、获取一个TCPAddr
	addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("Resolve tcp addr error:", err)
		return
	}
	// 2、监听服务器地址
	listener, err := net.ListenTCP(s.IPVersion, addr)
	if err != nil {
		fmt.Println("Listen", s.IP, "error:", err)
		return
	}
	fmt.Println("[Zinx]Start listening:", fmt.Sprintf("%s:%d", s.IP, s.Port))
	var connID uint32 = 0
	// 3、阻塞地等待客户端链接，处理客户端业务(读写)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept tcp error:", err)
			continue
		}
		// 检查是否超过最大连接数
		if s.ConnMgr.Size() >= utils.GlobalObject.MaxConn {
			fmt.Println("[WARNING]Conn num over the limit, new conn will be close")
			conn.Close()
			continue
		}
		// 将连接加入连接管理模块中
		newConn := NewConnection(s, conn, connID, s.MsgHandler)
		s.ConnMgr.Add(newConn)
		connID++
		go newConn.Start()
	}
}

func (s *Server) Stop() {
	// 将一些资源进行停止或回收
	fmt.Println("Close all conn...")
	s.ConnMgr.CloseAll()
	fmt.Println("All conn is closed!")
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// NewServer 初始化Server函数
func NewServer(name ...string) ziface.IServer {
	// 没在代码中声明服务器名称时使用配置文件中的名称
	if len(name) == 0 {
		name = make([]string, 1)
		name[0] = utils.GlobalObject.Name
	}
	return &Server{
		Name:       name[0],
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandler(),
		ConnMgr:    NewConnManager(),
	}
}
