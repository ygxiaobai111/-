package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
)

var RoleResService = &roleResService{}

type roleResService struct {
}

func (r *roleResService) GetRoleRes(rid int) *data.RoleRes {
	roleRes := &data.RoleRes{}
	ok, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return nil
	}
	if ok {
		return roleRes
	}
	return nil
}

func (r *roleResService) GetYield(rid int) data.Yield {
	//基础产量+ 建筑产量 + 城池设施的产量
	rbYield := RoleBuildService.GetYield(rid)
	cfYield := CityFacilityService.GetYield(rid)
	var yield data.Yield
	yield.Gold = rbYield.Gold + cfYield.Gold + gameConfig.Base.Role.GoldYield
	yield.Stone = rbYield.Stone + cfYield.Stone + gameConfig.Base.Role.StoneYield
	yield.Iron = rbYield.Iron + cfYield.Iron + gameConfig.Base.Role.IronYield
	yield.Grain = rbYield.Grain + cfYield.Grain + gameConfig.Base.Role.GrainYield
	yield.Wood = rbYield.Wood + cfYield.Wood + gameConfig.Base.Role.WoodYield
	return yield
}

func (r *roleResService) IsEnoughGold(rid int, cost int) bool {
	rr := r.GetRoleRes(rid)
	return rr.Gold >= cost
}

func (r *roleResService) CostGold(rid int, cost int) {
	rr := r.GetRoleRes(rid)
	if rr.Gold >= cost {
		rr.Gold -= cost
		rr.SyncExecute()
	}
}
