package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/global"
	"log"
	"math/rand"
	"time"
)

var RoleCityService = &roleCityService{}

type roleCityService struct {
}

func (r *roleCityService) InitCity(rid int, name string, conn net.WSConn) error {
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
		roleCity.X = rand.Intn(global.MapWith)
		roleCity.Y = rand.Intn(global.MapHeight)
		//这个城池 能不能在这个坐标点创建 需要判断 五格之内 不能有玩家的城池
		//TODO
		roleCity.RId = rid
		roleCity.Name = name
		roleCity.CurDurable = gameConfig.Base.City.Durable
		roleCity.CreatedAt = time.Now()
		roleCity.IsMain = 1
		_, err = db.Engine.Table(roleCity).Insert(roleCity)
		if err != nil {
			log.Println("插入角色城池出错", err)
			return common.New(constant.DBError, "数据库出错")
		}
		//初始化城池的设施
		//TODO
	}
	return nil
}

func (r *roleCityService) GetRoleCitys(rid int) ([]model.MapRoleCity, error) {
	citys := make([]data.MapRoleCity, 0)
	city := &data.MapRoleCity{}
	err := db.Engine.Table(city).Where("rid=?", rid).Find(&citys)
	modelCitys := make([]model.MapRoleCity, len(citys))
	if err != nil {
		log.Println("查询角色城池出错", err)
		return modelCitys, err
	}
	for _, v := range citys {
		modelCitys = append(modelCitys, v.ToModel().(model.MapRoleCity))
	}
	return modelCitys, nil

}
