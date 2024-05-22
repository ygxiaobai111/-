package login

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/login/controller"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/login/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/models"
)

var Router = net.NewRouter()

func Init() {
	//测试数据库，并且初始化数据库
	db.TestDB().Sync(new(models.User), new(model.LoginHistory), new(model.LoginLast))
	//还有别的初始化方法
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
