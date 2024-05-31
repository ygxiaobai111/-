package game

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/controller"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig/general"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/logic"
)

var Router = &net.Router{}

func Init() {
	db.TestDB()
	//加载基础配置
	gameConfig.Base.Load()
	//加载地图的资源配置
	gameConfig.MapBuildConf.Load()
	//加载地图单元格配置
	gameConfig.MapRes.Load()
	//加载城池设施配置
	gameConfig.FacilityConf.Load()
	//加载武将配置
	general.General.Load()
	//加载技能配置
	gameConfig.Skill.Load()

	logic.BeforeInit()

	//加载所有的建筑信息
	logic.RoleBuildService.Load()
	//加载所有的城池信息
	logic.RoleCityService.Load()
	//加载所有的角色属性
	logic.RoleAttrService.Load()
	initRouter()
}

func initRouter() {
	controller.DefaultRoleController.Router(Router)
	controller.DefaultNationMapController.Router(Router)
	controller.DefaultGeneralController.Router(Router)
	controller.DefaultArmyController.Router(Router)
	controller.WarController.Router(Router)
	controller.SkillController.Router(Router)
	controller.InteriorController.Router(Router)
}
