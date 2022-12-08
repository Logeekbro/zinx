package ziface

type IRouter interface {
	// PreHandle 处理请求业务前的Hook函数
	PreHandle(r IRequest)
	// Handle 处理请求的主函数
	Handle(r IRequest)
	// PostHandle 处理完请求之后的Hook函数
	PostHandle(r IRequest)
}
