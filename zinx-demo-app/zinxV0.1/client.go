package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	// 1、直接连接远程服务器，得到一个conn连接
	fmt.Println("Client starting...")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Client start error:", err)
		return
	}

	for {
		// 2、连接调用Write写数据
		_, err := conn.Write([]byte("ZinxV0.1 Hello"))
		if err != nil {
			fmt.Println("Write error:", err)
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read to buf error；", err)
		}
		fmt.Printf("Server callback: %s, cnt=%d\n", buf, cnt)
		// 阻塞一会
		time.Sleep(1 * time.Second)
	}
}
