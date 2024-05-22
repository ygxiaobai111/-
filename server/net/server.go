package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type server struct {
	addr       string
	router     *Router
	needSecret bool
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) NeedSecret(needSecret bool) {
	s.needSecret = needSecret
}

func (s *server) Router(router *Router) {
	s.router = router
}

// 启动服务
func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	//思考 websocket
	//1. http协议升级为websocket协议
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		//打印日志
		log.Println("websocket服务连接出错", err)
	}

	//发消息的时候 把消息当做路由 来去处理 消息是有格式的，先定义消息的格式
	//客户端 发消息的时候 {Name:"account.login"} 收到之后 进行解析 认为想要处理登录逻辑

	wsServer := NewWsServer(wsConn, s.needSecret)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()
}
