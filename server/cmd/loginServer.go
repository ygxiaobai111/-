package main

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/config"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/login"
)

func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")
	s := net.NewServer(host + ":" + port)
	login.Init()
	s.Router(login.Router)
	s.Start()
}
