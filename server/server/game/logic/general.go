package logic

import (
	"encoding/json"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig/general"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
	"time"
)

var GeneralService = &generalService{}

type generalService struct {
}

// GetGenerals 获取武将
func (g *generalService) GetGenerals(rid int) ([]model.General, error) {
	mrs := make([]*data.General, 0)
	mr := &data.General{}
	err := db.Engine.Table(mr).Where("rid=?", rid).Find(&mrs)
	if err != nil {
		log.Println("武将查询出错", err)
		return nil, common.New(constant.DBError, "武将查询出错")
	}
	if len(mrs) <= 0 {
		//随机3个武将
		var count = 0
		for {
			if count >= 3 {
				break
			}
			cfgId := general.General.Rand()
			gen, err := g.NewGeneral(cfgId, rid, 0)
			if err != nil {
				log.Println(err)
				continue
			}
			mrs = append(mrs, gen)
			count++
		}
	}
	modelMrs := make([]model.General, 0)
	for _, v := range mrs {
		modelMrs = append(modelMrs, v.ToModel().(model.General))
	}
	return modelMrs, nil
}

const (
	GeneralNormal      = 0 //正常
	GeneralComposeStar = 1 //星级合成
	GeneralConvert     = 2 //转换
)

func (g *generalService) NewGeneral(cfgId int, rid int, level int8) (*data.General, error) {
	cfg := general.General.GMap[cfgId]
	//初始 武将 无技能 但是有三个技能槽
	sa := make([]*model.GSkill, 3)
	ss, _ := json.Marshal(sa)
	gen := &data.General{
		PhysicalPower: gameConfig.Base.General.PhysicalPowerLimit,
		RId:           rid,
		CfgId:         cfg.CfgId,
		Order:         0,
		CityId:        0,
		Level:         level,
		CreatedAt:     time.Now(),
		CurArms:       cfg.Arms[0],
		HasPrPoint:    0,
		UsePrPoint:    0,
		AttackDis:     0,
		ForceAdded:    0,
		StrategyAdded: 0,
		DefenseAdded:  0,
		SpeedAdded:    0,
		DestroyAdded:  0,
		Star:          cfg.Star,
		StarLv:        0,
		ParentId:      0,
		SkillsArray:   sa,
		Skills:        string(ss),
		State:         GeneralNormal,
	}

	_, err := db.Engine.Table(gen).Insert(gen)
	if err != nil {
		log.Println("GetGenerals插入", err)
		return nil, err
	}
	return gen, nil
}
