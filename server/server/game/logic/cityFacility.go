package logic

import (
	"encoding/json"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
	"xorm.io/xorm"
)

var CityFacilityService = &cityFacilityService{}

type cityFacilityService struct {
}

func (c *cityFacilityService) TryCreate(cid, rid int, req *net.WsMsgReq) error {

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
	if session := req.Context.Get("dbsession"); session != nil {
		_, err = session.(*xorm.Session).Table(cf).Insert(cf)
	} else {
		_, err = db.Engine.Table(cf).Insert(cf)
	}

	if err != nil {
		log.Println("插入城市设施出错", err)
		return common.New(constant.DBError, "数据库错误")
	}
	return nil
}
func (c *cityFacilityService) GetByRId(rid int) ([]*data.CityFacility, error) {
	cf := make([]*data.CityFacility, 0)
	err := db.Engine.Table(new(data.CityFacility)).Where("rid=?", rid).Find(&cf)
	if err != nil {
		log.Println(err)
		return cf, common.New(constant.DBError, "数据库错误")
	}
	return cf, nil
}

func (c *cityFacilityService) GetYield(rid int) data.Yield {
	//查询 把表中的设施 获取到
	//设施的不同类型 去设施配置中查询匹配，匹配到增加产量的设施 木头 金钱 产量的计算
	//设施的等级不同 产量也不同
	cfs, err := c.GetByRId(rid)
	var y data.Yield
	if err == nil {
		for _, v := range cfs {
			facilities := v.Facility()
			for _, fa := range facilities {
				//计算等级 资源的产出是不同的
				if fa.GetLevel() > 0 {
					values := gameConfig.FacilityConf.GetValues(fa.Type, fa.GetLevel())
					adds := gameConfig.FacilityConf.GetAdditions(fa.Type)
					for i, aType := range adds {
						if aType == gameConfig.TypeWood {
							y.Wood += values[i]
						} else if aType == gameConfig.TypeGrain {
							y.Grain += values[i]
						} else if aType == gameConfig.TypeIron {
							y.Iron += values[i]
						} else if aType == gameConfig.TypeStone {
							y.Stone += values[i]
						} else if aType == gameConfig.TypeTax {
							y.Gold += values[i]
						}
					}
				}

			}
		}
	}
	return y
}
