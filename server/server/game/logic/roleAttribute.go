package logic

import (
	"encoding/json"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/net"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
	"sync"
	"xorm.io/xorm"
)

var RoleAttrService = &roleAttrService{
	attrs: make(map[int]*data.RoleAttribute),
}

type roleAttrService struct {
	mutex sync.RWMutex
	attrs map[int]*data.RoleAttribute
}

func (r *roleAttrService) Load() {
	ras := make([]*data.RoleAttribute, 0)
	err := db.Engine.Table(new(data.RoleAttribute)).Find(&ras)
	if err != nil {
		log.Println("roleAttrService Load err", err)
	}
	for _, v := range ras {
		r.attrs[v.RId] = v
	}

}
func (r *roleAttrService) TryCreate(rid int, req *net.WsMsgReq) error {

	role := &data.RoleAttribute{}
	ok, err := db.Engine.Table(role).Where("rid=?", rid).Get(role)
	if err != nil {
		log.Println("查询角色属性出错", err)
		return common.New(constant.DBError, "数据库出错")
	}
	if ok {
		//缓存
		r.mutex.Lock()
		r.attrs[rid] = role
		defer r.mutex.Unlock()
		return nil
	} else {
		//初始化
		role.RId = rid
		role.UnionId = 0
		role.ParentId = 0
		role.PosTags = ""
		if session := req.Context.Get("dbSession"); session != nil {
			_, err = session.(*xorm.Session).Table(role).Insert(role)
		} else {
			_, err = db.Engine.Table(role).Insert(role)
		}

		if err != nil {
			log.Println("插入角色属性出错", err)
			return common.New(constant.DBError, "数据库出错")
		}
		r.mutex.Lock()
		r.attrs[rid] = role
		defer r.mutex.Unlock()
	}
	return nil
}
func (r *roleAttrService) GetTagList(rid int) ([]model.PosTag, error) {
	ra, ok := r.attrs[rid]
	if !ok {
		ra = &data.RoleAttribute{}
		var err error
		ok, err = db.Engine.Table(ra).Where("rid=?", rid).Get(ra)
		if err != nil {
			log.Println("GetTagList", err)
			return nil, common.New(constant.DBError, "数据库错误")
		}
	}

	posTags := make([]model.PosTag, 0)
	if ok {
		tags := ra.PosTags
		if tags != "" {
			err := json.Unmarshal([]byte(tags), &posTags)
			if err != nil {
				return nil, common.New(constant.DBError, "数据库错误")
			}
		}
	}
	return posTags, nil
}
func (r *roleAttrService) Get(rid int) *data.RoleAttribute {

	r.mutex.RLock()
	defer r.mutex.RUnlock()
	ra, ok := r.attrs[rid]
	if ok {
		return ra
	}
	return nil
}
