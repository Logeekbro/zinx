package znet

import (
	"errors"
	"fmt"
	"net"
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
	Port int
	// 服务器路由方法
	Router ziface.IRouter
}

// 回调方法，暂时写死
func callback(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("Callback running...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		return errors.New("CallBack error")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Println("Starting server:", s.Name)
	// 服务器启动步骤
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
	fmt.Println("Start listening:", fmt.Sprintf("%s:%d", s.IP, s.Port))
	var connID uint32 = 0
	// 3、阻塞地等待客户端链接，处理客户端业务(读写)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Accept tcp error:", err)
			continue
		}
		newConn := NewConnection(conn, connID, s.Router)
		connID++
		go newConn.Start()
	}
}

func (s *Server) Stop() {
	// TODO 将一些资源进行停止或回收
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务器之后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add Router success!")
}

// NewServer 初始化Server函数
func NewServer(name string) ziface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
}
