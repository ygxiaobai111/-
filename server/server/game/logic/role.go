package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/gameConfig"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	utils "github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/util"
	"log"

	"time"
)

var RoleService = &roleService{}

type roleService struct {
}

func (r *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, conn net.WSConn) error {
	role := &data.Role{}
	ok, err := db.Engine.Table(role).Where("uid=?", uid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return common.New(constant.DBError, "数据库出错")
	}
	if ok {
		rid := role.RId
		roleRes := &data.RoleRes{}
		ok, err = db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
		if err != nil {
			log.Println("查询角色资源出错", err)
			return common.New(constant.DBError, "数据库出错")
		}
		if !ok {
			roleRes.RId = rid
			roleRes.Gold = gameConfig.Base.Role.Gold
			roleRes.Decree = gameConfig.Base.Role.Decree
			roleRes.Grain = gameConfig.Base.Role.Grain
			roleRes.Iron = gameConfig.Base.Role.Iron
			roleRes.Stone = gameConfig.Base.Role.Stone
			roleRes.Wood = gameConfig.Base.Role.Wood
			_, err := db.Engine.Table(roleRes).Insert(roleRes)
			if err != nil {
				log.Println("插入角色资源出错", err)
				return common.New(constant.DBError, "数据库出错")
			}
			log.Println("插入角色资源成功")
		}
		//模型转换
		rsp.RoleRes = roleRes.ToModel().(model.RoleRes)
		rsp.Role = role.ToModel().(model.Role)
		rsp.Time = time.Now().UnixNano() / 1e6
		token, _ := utils.Award(rid)
		rsp.Token = token
		conn.SetProperty("role", role)
		// 初始化玩家属性
		if err := RoleAttrService.TryCreate(rid, conn); err != nil {
			return common.New(constant.DBError, "数据库错误")
		}
		//初始化城池
		if err := RoleCityService.InitCity(rid, role.NickName, conn); err != nil {
			return common.New(constant.DBError, "数据库错误")
		}

	} else {
		return common.New(constant.RoleNotExist, "角色不存在")
	}
	return nil
}

func (r *roleService) GetRoleRes(rid int) (model.RoleRes, error) {
	roleRes := &data.RoleRes{}
	ok, err := db.Engine.Table(roleRes).Where("rid=?", rid).Get(roleRes)
	if err != nil {
		log.Println("查询角色资源出错", err)
		return model.RoleRes{}, common.New(constant.DBError, "数据库出错")
	}
	if ok {
		return roleRes.ToModel().(model.RoleRes), nil
	}
	return model.RoleRes{}, common.New(constant.DBError, "角色资源不存在")
}
