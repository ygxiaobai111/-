package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/global"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
	"math/rand"
	"time"
	"xorm.io/xorm"
)

var RoleCityService = &roleCityService{}

type roleCityService struct {
}

func (r *roleCityService) InitCity(rid int, name string, req *net.WsMsgReq) error {

	roleCity := &data.MapRoleCity{}
	ok, err := db.Engine.Table(roleCity).Where("rid=?", rid).Get(roleCity)
	if err != nil {
		log.Println("查询角色城池出错", err)
		return common.New(constant.DBError, "数据库出错")
	}
	if ok {
		return nil
	} else {
		//初始化
		for {
			roleCity.X = rand.Intn(global.MapWith)
			roleCity.Y = rand.Intn(global.MapHeight)
			//这个城池 能不能在这个坐标点创建 需要判断 系统城池五格之内 不能有玩家的城池
			if IsCanBuild(roleCity.X, roleCity.Y) {
				roleCity.RId = rid
				roleCity.Name = name
				roleCity.CurDurable = gameConfig.Base.City.Durable
				roleCity.CreatedAt = time.Now()
				roleCity.IsMain = 1
				if session := req.Context.Get("dbsession"); session != nil {
					_, err = session.(*xorm.Session).Table(roleCity).Insert(roleCity)
				} else {
					_, err = db.Engine.Table(roleCity).Insert(roleCity)
				}
				if err != nil {
					log.Println("插入角色城池出错", err)
					return common.New(constant.DBError, "数据库出错")
				}
				//初始化城池的设施
				if err := CityFacilityService.TryCreate(roleCity.CityId, rid, req); err != nil {
					log.Println("城池设施出错", err)
					return common.New(err.(*common.MyError).Code(), err.Error())
				}
				break
			}
		}

	}
	return nil
}

func IsCanBuild(x int, y int) bool {
	confs := gameConfig.MapRes.Confs
	pIndex := global.ToPosition(x, y)
	_, ok := confs[pIndex]
	if !ok {
		return false
	}
	sysBuild := gameConfig.MapRes.SysBuild
	//系统城池的5格内 不能创建玩家城池
	//此处逻辑为遍历所有的系统城池，然后与玩家城池坐标进行判断 todo:可优化
	for _, v := range sysBuild {
		if v.Type == gameConfig.MapBuildSysCity {
			if x >= v.X-5 &&
				x <= v.X+5 &&
				y >= v.Y-5 &&
				y <= v.Y+5 {
				return false
			}
		}
	}
	return true
}

func (r *roleCityService) GetRoleCitys(rid int) ([]model.MapRoleCity, error) {
	citys := make([]data.MapRoleCity, 0)
	city := &data.MapRoleCity{}
	err := db.Engine.Table(city).Where("rid=?", rid).Find(&citys)
	modelCitys := make([]model.MapRoleCity, 0)
	if err != nil {
		log.Println("查询角色城池出错", err)
		return modelCitys, err
	}
	for _, v := range citys {
		modelCitys = append(modelCitys, v.ToModel().(model.MapRoleCity))
	}
	return modelCitys, nil

}
