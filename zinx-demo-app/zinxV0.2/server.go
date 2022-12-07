package main

import "zinx/znet"

func main() {
	s := znet.NewServer("[Zinx0.2]")
	s.Serve()
}
