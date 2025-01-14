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

func (r *roleService) EnterServer(uid int, rsp *model.EnterServerRsp, req *net.WsMsgReq) error {
	role := &data.Role{}
	session := db.Engine.NewSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		log.Println("数据库出错", err)
		return common.New(constant.DBError, "数据库出错")
	}
	req.Context.Set("dbSession", session)
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
			_, err := session.Table(roleRes).Insert(roleRes)
			if err != nil {
				session.Rollback()
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
		req.Conn.SetProperty("role", role)
		// 初始化玩家属性
		if err := RoleAttrService.TryCreate(rid, req); err != nil {
			session.Rollback()
			return common.New(constant.DBError, "数据库错误")
		}
		//初始化城池
		if err := RoleCityService.InitCity(rid, role.NickName, req); err != nil {
			session.Rollback()
			return common.New(constant.DBError, "数据库错误")
		}

	} else {
		return common.New(constant.RoleNotExist, "角色不存在")
	}
	if err := session.Commit(); err != nil {
		log.Println("数据库出错", err)
		return common.New(constant.DBError, "数据库出错")
	}
	log.Println("“user", rsp.Role)

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

func (r *roleService) Get(rid int) *data.Role {
	role := &data.Role{}
	ok, err := db.Engine.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return nil
	}
	if ok {
		return role
	}
	return nil
}

func (r *roleService) GetRoleNickName(rid int) string {
	role := &data.Role{}
	ok, err := db.Engine.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色出错", err)
		return ""
	}
	if ok {
		return role.NickName
	}
	return ""
}
