package logic

import "github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"

// BeforeInit 避免循环依赖

func BeforeInit() {
	data.GetYield = RoleResService.GetYield
	data.GetUnion = RoleAttrService.GetUnion
	data.GetParentId = RoleAttrService.GetParentId
	data.MapResTypeLevel = RoleBuildService.MapResTypeLevel
	data.GetMainMembers = CoalitionService.GetMainMembers
	data.GetRoleNickName = RoleService.GetRoleNickName
}
