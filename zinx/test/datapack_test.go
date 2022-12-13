package test

import (
	"fmt"
	"net"
	"testing"
	"zinx/znet"
)

// 封包、拆包的单元测试
func TestDataPack(t *testing.T) {
	/**
	模拟服务器
	*/
	// 1 创建TCP socket
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("listen addr error:", err)
	}

	fmt.Println("listening...")

	// 等待客户端连接
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("Connect to client error:", err)
				return
			}
			go func(conn net.Conn) {
				// 2 从客户端连接中读取数据，拆包处理
				dp := znet.DataPack{}
				for {
					// 第一次把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					cnt, err := conn.Read(headData)
					if err != nil {
						fmt.Println("Read headData error:", err)
						break
					} else if cnt != int(dp.GetHeadLen()) {
						fmt.Println("Invalid headLen !")
						break
					}
					// 先将headData解包, 获取Message对象
					msg, err := dp.UnPack(headData)
					/**
					新知识: message = msg.(*Message) 可将接口类型的IMessage 转为它的实现类 Message (类型断言)
					*/
					// 第二次把内容读出来
					if msg.GetDataLen() > 0 {
						data := make([]byte, msg.GetDataLen())
						_, err = conn.Read(data)
						if err != nil {
							fmt.Println("Read data error:", err)
							break
						}
						msg.SetData(data)
						// 此时一个完整的消息已经读取完毕
						fmt.Println("Recv data info: Id:", msg.GetId(), "DataLen:", msg.GetDataLen(), "Data:", string(msg.GetData()))
					}

				}

			}(conn)
		}
	}()

	/**
	模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("Connect to server error:", err)
	}
	dp := znet.DataPack{}

	// 第一个包
	message1 := znet.Message{
		Id:      1,
		DataLen: 5,
		Data:    []byte("hello"),
	}
	sendData1, err := dp.Pack(&message1)
	if err != nil {
		fmt.Println("Pack data error:", err)
	}
	// 第二个包
	message2 := znet.Message{
		Id:      2,
		DataLen: 8,
		Data:    []byte("theworld"),
	}
	sendData2, err := dp.Pack(&message2)
	if err != nil {
		fmt.Println("Pack data error:", err)
	}
	// 拼接两个包 (需使用... 打散数据)
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
