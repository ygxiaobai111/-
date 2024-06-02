package chat

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/chat/controller"
)

var Router = &net.Router{}

func Init() {
	initRouter()
}

func initRouter() {

	controller.ChatController.Router(Router)

}
