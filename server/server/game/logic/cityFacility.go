package logic

import (
	"encoding/json"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
)

var CityFacilityService = &cityFacilityService{}

type cityFacilityService struct {
}

func (c *cityFacilityService) TryCreate(cid, rid int) error {
	cf := &data.CityFacility{}
	ok, err := db.Engine.Table(cf).Where("cityId=?", cid).Get(cf)
	if err != nil {
		log.Println("查询城市设施出错", err)
		return common.New(constant.DBError, "数据库错误")
	}
	if ok {
		return nil
	}
	cf.RId = rid
	cf.CityId = cid
	list := gameConfig.FacilityConf.List
	facs := make([]data.Facility, len(list))
	for index, v := range list {
		fac := data.Facility{
			Name:         v.Name,
			Type:         v.Type,
			PrivateLevel: 0,
			UpTime:       0,
		}
		facs[index] = fac
	}
	dataJson, _ := json.Marshal(facs)
	cf.Facilities = string(dataJson)
	_, err = db.Engine.Table(cf).Insert(cf)
	if err != nil {
		log.Println("插入城市设施出错", err)
		return common.New(constant.DBError, "数据库错误")
	}
	return nil
}
