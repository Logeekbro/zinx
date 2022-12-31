package utils

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"zinx/ziface"
)

/**
存储所有Zinx框架的全局参数，供其它模块使用
所有的参数都可通过 zinx.json 文件由用户进行配置
*/

type GlobalObj struct {
	/**
	Server
	*/
	TcpServer ziface.IServer // Zinx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   uint16         //当前服务器主机监听的端口
	Name      string         //当前服务器的名称

	/**
	Zinx
	*/
	Version          string //当前Zinx的版本号
	MaxConn          int    //当前服务器主机允许的最大连接数
	MaxPackageSize   uint32 //当前Zinx数据包的最大值
	WorkerPoolSize   uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskLen uint32 //Zinx框架允许每个Worker对应的消息队列任务的最大值
}

// GlobalObject 全局配置对象
var GlobalObject *GlobalObj

// LoadConfig 从zinx.json中加载用户自定义参数
func (g *GlobalObj) LoadConfig() {
	data, err := os.ReadFile("conf/zinx.json")
	if err != nil {
		// 找不到配置文件时使用默认配置
		if _, ok := err.(*fs.PathError); ok {
			fmt.Println("zinx.json not found, using default config...")
			return
		} else {
			panic(fmt.Sprintln("Read zinx.json failed:", err))
		}

	}
	// 将json文件数据解析到struct中
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(fmt.Sprintln("Load config failed:", err))
	}
}

// init函数: 在导包时会调用此函数进行初始化, 此函数只会执行一次, 此函数在main函数之前执行
func init() {
	// 配置文件没有加载时的默认值
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		TcpPort:          9999,
		Name:             "ZinxServerApp",
		Version:          "V0.8",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   0,
		MaxWorkerTaskLen: 1024,
	}
	// 应该尝试从 conf/zinx.json 中加载一些用户自定义的参数
	GlobalObject.LoadConfig()
}
