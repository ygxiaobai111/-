package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
)

var RoleAttrService = &roleAttrService{}

type roleAttrService struct {
}

func (r *roleAttrService) TryCreate(rid int, conn net.WSConn) error {
	role := &data.RoleAttribute{}
	ok, err := db.Engine.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色属性出错", err)
		return common.New(constant.DBError, "数据库出错")
	}
	if ok {
		return nil
	} else {
		//初始化
		role.RId = rid
		role.UnionId = 0
		role.ParentId = 0
		_, err = db.Engine.Table(role).Insert(role)
		if err != nil {
			log.Println("插入角色属性出错", err)
			return common.New(constant.DBError, "数据库出错")
		}
	}
	return nil
}
