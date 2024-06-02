package logic

import (
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/constant"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/db"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/common"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model"
	"github.com/ygxiaobai111/Three_Kingdoms_of_Longning/server/server/game/model/data"
	"log"
	"sync"
)

var CoalitionService = &coalitionService{
	unions: make(map[int]*data.Coalition),
}

type coalitionService struct {
	mutex  sync.RWMutex
	unions map[int]*data.Coalition
}

// members [1,2,3,4,5]  json的字符串
func (c *coalitionService) Load() {
	rr := make([]*data.Coalition, 0)
	err := db.Engine.Table(new(data.Coalition)).Where("state=?", data.UnionRunning).Find(&rr)
	if err != nil {
		log.Println("coalitionService load error", err)
	}
	for _, v := range rr {
		c.unions[v.Id] = v
	}
}

func (c *coalitionService) List() ([]model.Union, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	uns := make([]model.Union, 0)
	for _, v := range c.unions {
		mas := make([]model.Major, 0)
		//盟主
		if role := RoleService.Get(v.Chairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}
		if role := RoleService.Get(v.ViceChairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}

		union := v.ToModel().(model.Union)
		union.Major = mas
		uns = append(uns, union)
	}
	return uns, nil
}

func (c *coalitionService) ListCoalition() []*data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	uns := make([]*data.Coalition, 0)
	for _, v := range c.unions {
		uns = append(uns, v)
	}
	return uns
}

func (c *coalitionService) Get(id int) model.Union {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coa, ok := c.unions[id]
	if ok {
		union := coa.ToModel().(model.Union)
		mas := make([]model.Major, 0)
		if role := RoleService.Get(coa.Chairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}
		if role := RoleService.Get(coa.ViceChairman); role != nil {
			ma := model.Major{
				RId:   role.RId,
				Name:  role.NickName,
				Title: model.UnionChairman,
			}
			mas = append(mas, ma)
		}
		union.Major = mas
		return union
	}
	return model.Union{}
}

func (c *coalitionService) GetCoalition(id int) *data.Coalition {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	coa, ok := c.unions[id]
	if ok {
		return coa
	}
	return nil
}

func (c *coalitionService) GetListApply(unionId int, state int) ([]model.ApplyItem, error) {
	applys := make([]data.CoalitionApply, 0)
	err := db.Engine.Table(new(data.CoalitionApply)).
		Where("union_id=? and state=?", unionId, state).
		Find(&applys)
	if err != nil {
		log.Println("coalitionService GetListApply find error", err)
		return nil, common.New(constant.DBError, "数据库错误")
	}
	ais := make([]model.ApplyItem, 0)
	for _, v := range applys {
		var ai model.ApplyItem
		ai.Id = v.Id
		role := RoleService.Get(v.RId)
		ai.NickName = role.NickName
		ai.RId = role.RId
		ais = append(ais, ai)
	}
	return ais, nil
}
