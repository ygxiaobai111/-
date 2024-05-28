package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
)

var SkillService = &skillService{}

type skillService struct {
}

func (s *skillService) GetSkills(rid int) ([]model.Skill, error) {
	mrs := make([]data.Skill, 0)
	mr := &data.Skill{}
	err := db.Engine.Table(mr).
		Where("rid=? ", rid).Find(&mrs)
	if err != nil {
		log.Println("技能查询出错", err)
		return nil, common.New(constant.DBError, "技能查询出错")
	}
	modelMrs := make([]model.Skill, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.Skill))
	}
	return modelMrs, nil
}
