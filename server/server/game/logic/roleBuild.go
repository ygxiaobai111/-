package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
)

var RoleBuildService = &roleBuildService{}

type roleBuildService struct {
}

func (r *roleBuildService) GetBuilds(rid int) ([]model.MapRoleBuild, error) {
	mrs := make([]data.MapRoleBuild, 0)
	mr := &data.MapRoleBuild{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("建筑查询出错", err)
		return nil, common.New(constant.DBError, "建筑查询出错")
	}
	modelMrs := make([]model.MapRoleBuild, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.MapRoleBuild))
	}
	return modelMrs, nil
}
