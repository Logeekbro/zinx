package main

import (
	"fmt"
	"net"
	"time"
	"zinx/znet"
)

func main() {
	// 1、直接连接远程服务器，得到一个conn连接
	fmt.Println("Client0 starting...")
	time.Sleep(1 * time.Second)
	conn, err := net.Dial("tcp", "127.0.0.1:6666")
	if err != nil {
		fmt.Println("Client start error:", err)
		return
	}
	// 获取用于封包、拆包的对象
	dp := znet.NewDataPack()
	for {
		// 将消息打包
		msgPack, err := dp.Pack(znet.NewMessage(0, []byte("ping...")))
		if err != nil {
			fmt.Println("Pack message error:", err)
			return
		}
		// 将消息发送给服务端
		_, err = conn.Write(msgPack)
		if err != nil {
			fmt.Println("Send message to server error:", err)
			return
		}
		// 将服务端返回的消息进行拆包
		// 1、获取消息头
		binaryHeadData := make([]byte, dp.GetHeadLen())
		_, err = conn.Read(binaryHeadData)
		if err != nil {
			fmt.Println("Read headData error:", err)
			return
		}
		message, err := dp.UnPack(binaryHeadData)
		if err != nil {
			fmt.Println("UnPack headData error:", err)
			return
		}
		// 根据消息头提供的信息获取具体数据内容
		binaryData := make([]byte, message.GetDataLen())
		_, err = conn.Read(binaryData)
		if err != nil {
			fmt.Println("Read data error:", err)
			return
		}
		// 将获取到的内容放入message中
		message.SetData(binaryData)
		// 打印获取到的信息
		fmt.Printf("---->[Client0]Recv server message: MsgId=%d , MsgLen=%d, MsgData=%s\n", message.GetId(), message.GetDataLen(), string(message.GetData()))
		// 阻塞一会
		time.Sleep(1 * time.Second)
	}
}
