package game

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/controller"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
)

var Router = &net.Router{}

func Init() {
	db.TestDB()
	//加载基础配置
	gameConfig.Base.Load()
	//加载地图的资源配置
	gameConfig.MapBuildConf.Load()
	initRouter()
}

func initRouter() {
	controller.DefaultRoleController.Router(Router)
	controller.DefaultNationMapController.Router(Router)
}
